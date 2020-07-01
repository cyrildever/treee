package index

import (
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/cyrildever/treee/common/logger"
	"github.com/cyrildever/treee/config"
	"github.com/cyrildever/treee/core/exception"
	"github.com/cyrildever/treee/core/index/branch"
	"github.com/cyrildever/treee/core/model"
	"github.com/cyrildever/treee/utils"
	"github.com/cyrildever/treee/utils/prime"
)

// Current is the current Treee index used when running the executable app
var Current *Treee

var saving bool = false

const (
	// INIT_PRIME ...
	INIT_PRIME uint64 = 2 // TODO Shouldn't be used in production
)

//--- TYPES

// Treee is an instance of a database tree
type Treee struct {
	InitPrime uint64
	sync.RWMutex
	trunk               *branch.Node
	size                uint64
	persistence         bool
	overridePersistence bool
}

//--- METHODS

// Add ...
func (t *Treee) Add(item branch.Leaf) error {
	t.Lock()
	defer t.Unlock()

	log := logger.Init("index", "Add")

	if item.Size == 0 {
		return exception.NewEmptyItemError()
	}
	idStr, err := item.ID.String()
	if err != nil {
		return err
	}

	// 1- Prepare and check
	if found, err := t.search(item.ID); err == nil && found.ID == item.ID {
		return exception.NewAlreadyExistsInIndexError(idStr)
	}

	var previous *branch.Leaf
	previousID, e := item.Previous.String()
	if item.Previous.NonEmpty() && e == nil && previousID != idStr {
		existingPrevious, err := t.search(item.Previous)
		if err != nil {
			return err
		}
		previous = existingPrevious
	} else {
		previous = &item
	}
	item.Origin = previous.Origin
	item.Previous = previous.ID

	var origin *branch.Leaf
	if item.Origin.NonEmpty() && item.Origin != item.ID {
		existingOrigin, err := t.search(item.Origin)
		if err != nil {
			return err
		}
		origin = existingOrigin
	} else {
		origin = &item
	}
	item.Origin = origin.ID

	item.Next = model.EmptyHash

	// 2- Actually add it to the Treee index
	id := new(big.Int)
	id.SetString(idStr, 16)
	currentNode := t.trunk
	currentStage := new(big.Int)
	for {
		usedStage := new(big.Int)
		usedStage.Set(currentStage)
		currentStage.SetUint64(currentNode.StagePrime)
		if usedStage.Cmp(currentStage) == 0 {
			return exception.NewLoopError("adding")
		}
		modulo := new(big.Int)
		modulo = modulo.Mod(id, currentStage)
		idx := modulo.Uint64()
		targetBranch, exists := currentNode.ChildAt(idx)
		if !exists || targetBranch.IsEmpty() {
			if !targetBranch.Assign(&item) {
				// This shouldn't happen so we'd better log it
				log.Crit("Impossible to assign non-pointer", "leafPtr", &item)
				return utils.NewNotAPointerError()
			}
			previous.Next = item.ID
			origin.Previous = item.ID
			t.size++
			return nil
		} else if targetBranch.IsLeaf() {
			existingLeaf := targetBranch.GetLeaf()
			nextPrime, err := prime.Next(currentStage.Uint64())
			if err != nil {
				// This shouldn't happen so we'd better log it
				log.Crit("Unable to get next prime number", "error", err)
				return err
			}
			newNode := branch.NewNode(nextPrime)
			newNode.AddLeaf(existingLeaf)
			newNode.AddLeaf(&item)
			if !targetBranch.Assign(newNode) {
				// This shouldn't happen so we'd better log it
				log.Crit("Impossible to assign non-pointer", "nodePtr", newNode)
				return utils.NewNotAPointerError()
			}
			previous.Next = item.ID
			origin.Previous = item.ID
			t.size++
			return nil
		} else if targetBranch.IsNode() {
			currentNode = targetBranch.GetNode()
		}
	}
}

// Last finds the last item in a subchain from any ID of the subchain;
// it implements `search.Engine`
func (t *Treee) Last(id model.Hash) (lastInChain *branch.Leaf, err error) {
	t.RLock()
	defer t.RUnlock()

	found, err := t.search(id)
	if err != nil {
		return
	}
	if found.Previous == found.ID || found.Next.IsEmpty() {
		lastInChain = found
		return
	}
	origin, err := t.search(found.Origin)
	if err != nil {
		return
	}
	last, err := t.search(origin.Previous) // By definition of the circular linked list
	if err != nil {
		return
	}
	lastInChain = last
	return
}

// Line fetches the whole list of items in a subchain
func (t *Treee) Line(id model.Hash) (subchain []*branch.Leaf, err error) {
	t.RLock()
	defer t.RUnlock()

	found, err := t.search(id)
	if err != nil {
		return
	}
	if found.ID == found.Origin && found.Next.IsEmpty() {
		return []*branch.Leaf{found}, nil
	}
	origin, err := t.search(found.Origin)
	if err != nil {
		return
	}
	subchain = append(subchain, origin)
	current := origin.ID
	next := origin.Next
	for next.NonEmpty() && next != current {
		following, e := t.search(next)
		if e != nil {
			if _, ok := e.(*exception.NotFoundError); ok {
				err = e
			}
			return
		}
		subchain = append(subchain, following)
		current = following.ID
		next = following.Next
	}
	return
}

// PrintAll ...
// Use with caution!
func (t *Treee) PrintAll(beautify bool) string {
	t.RLock()
	defer t.RUnlock()

	str := `{"initPrime":` + strconv.FormatUint(t.InitPrime, 10) + `,"trunk":` + t.trunk.Print() + `,"size":` + strconv.FormatUint(t.size, 10) + "}"
	if beautify {
		var js interface{}
		_ = json.Unmarshal([]byte(str), &js)
		bytes, _ := json.MarshalIndent(js, "", "  ")
		str = string(bytes)
	}
	return str
}

// Save ...
func (t *Treee) Save() {
	log := logger.Init("index", "Save")
	conf, _ := config.GetConfig()

	if !saving && (t.persistence || (!t.overridePersistence && conf.UsePersistence)) && !conf.IsTestEnvironment() {
		saving = true
		t0 := time.Now().UnixNano()
		path := conf.IndexPath
		if path == "" {
			path = "saved" + string(os.PathSeparator) + "treee.json"
		}
		f, err := os.Create(path)
		if err != nil {
			log.Crit("Unable to create file", "error", err)
			saving = false
			return
		}
		size := t.Size()
		n, err := f.WriteString(t.PrintAll(false))
		if err != nil {
			t1 := time.Now().UnixNano()
			log.Error("An error occurred while saving the index", "error", err, "after", strconv.FormatInt((t1-t0)/int64(time.Millisecond), 10)+"ms")
			saving = false
			return
		}
		t1 := time.Now().UnixNano()
		log.Info("Index saved", "size", size, "bytes", n, "duration", strconv.FormatInt((t1-t0)/int64(time.Millisecond), 10)+"ms")
		f.Close()
		saving = false
		return
	}
}

// Search fetches a Leaf from the Treee index;
// it implements `search.Engine`
func (t *Treee) Search(ID model.Hash) (found *branch.Leaf, err error) {
	t.RLock()
	defer t.RUnlock()

	return t.search(ID)
}

func (t *Treee) search(ID model.Hash) (found *branch.Leaf, err error) {
	idStr, err := ID.String()
	if err != nil {
		return
	}
	id := new(big.Int)
	id.SetString(idStr, 16)
	currentNode := t.trunk
	currentStage := new(big.Int)
	for {
		usedStage := new(big.Int)
		usedStage.Set(currentStage)
		currentStage.SetUint64(currentNode.StagePrime)
		if usedStage.Cmp(currentStage) == 0 {
			err = exception.NewLoopError("finding")
			return
		}
		modulo := new(big.Int)
		modulo = modulo.Mod(id, currentStage)
		idx := modulo.Uint64()
		targetBranch, exists := currentNode.ChildAt(idx)
		if !exists || targetBranch.IsEmpty() {
			err = exception.NewNotFoundError(idStr)
			return
		} else if targetBranch.IsLeaf() {
			found = targetBranch.GetLeaf()
			return
		} else if targetBranch.IsNode() {
			currentNode = targetBranch.GetNode()
		}
	}
}

// Size ...
func (t *Treee) Size() uint64 {
	t.RLock()
	defer t.RUnlock()

	return t.size
}

// UsePersistence forces the index to be permanent or not, overriding the corresponding configuration parameter
func (t *Treee) UsePersistence(value bool) {
	conf, _ := config.GetConfig()
	conf.UsePersistence = value
	t.persistence = value
	t.overridePersistence = true
}

//--- FUNCTIONS

// Load ...
func Load(path string) (t *Treee, err error) {
	if path == "" {
		path = "saved" + string(os.PathSeparator) + "treee.json"
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}

	type savedTreee struct {
		InitPrime uint64                 `json:"initPrime"`
		Trunk     map[string]interface{} `json:"trunk"`
		Size      uint64                 `json:"size"`
	}

	st := savedTreee{}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&st)
	if err != nil {
		return
	}

	trunk := branch.NewNode(st.InitPrime)

	stagePrime, ok := st.Trunk["stagePrime"]
	if !ok {
		err = exception.NewNotAValidTreeeError()
		return
	}
	if sp, ok := stagePrime.(float64); !ok || uint64(sp) != st.InitPrime || !prime.IsPrime(uint64(sp)) {
		err = prime.NewNotAValidNumberError(uint64(sp))
		return
	}

	actualSize := 0
	err = parse(trunk, st.Trunk["children"].([]interface{}), &actualSize)
	if err != nil {
		return
	}
	var treee Treee
	if actualSize == int(st.Size) {
		treee = Treee{
			InitPrime: st.InitPrime,
			trunk:     trunk,
			size:      st.Size,
		}
	} else {
		return &treee, exception.NewIncoherentSizeError(int(st.Size), actualSize)
	}

	return &treee, nil
}

// New ...
func New(initPrime uint64) (t *Treee, err error) {
	if prime.IsPrime(initPrime) {
		t = &Treee{
			InitPrime: initPrime,
			trunk:     branch.NewNode(initPrime),
		}
	} else if initPrime == 0 {
		t = &Treee{
			InitPrime: INIT_PRIME,
			trunk:     branch.NewNode(INIT_PRIME),
		}
	} else {
		err = prime.NewNotAValidNumberError(initPrime)
	}
	return
}

// parse fills the passed node with the information in the array of interface 'arr', and increments the global actual size of the index;
// 'arr' is supposed to be in the form of map[0:...] where 0 is the remainder/index at the considered stage and ... either an empty map/branch,
// the children of a node, or a leaf as a map
func parse(node *branch.Node, arr []interface{}, actualSize *int) error {
	for _, v := range arr {
		if val, ok := v.(map[string]interface{}); ok {
			for remainder, child := range val {
				b := branch.Branch{}
				if value, ok := child.(map[string]interface{}); ok {
					if len(value) == 0 {
						idx, err := strconv.Atoi(remainder)
						if err != nil {
							return err
						}
						node.AddBranch(&b, uint64(idx))
						continue
					}
					if n, ok := value["children"]; ok {
						stagePrime, ok := value["stagePrime"]
						if !ok {
							return exception.NewNotAValidTreeeError()
						}
						sp, ok := stagePrime.(float64)
						if !ok || !prime.IsPrime(uint64(sp)) {
							return prime.NewNotAValidNumberError(uint64(sp))
						}
						newNode := branch.NewNode(uint64(sp))
						err := parse(newNode, n.([]interface{}), actualSize)
						if err != nil {
							return err
						}
						if !b.Assign(newNode) {
							return utils.NewNotAPointerError()
						}
					} else {
						if _, ok := value["id"]; ok {
							position, _ := value["position"].(float64)
							size, _ := value["size"].(float64)
							leaf := branch.Leaf{
								ID:       model.Hash(value["id"].(string)),
								Position: int64(position),
								Size:     int64(size),
								Origin:   model.Hash(value["origin"].(string)),
								Previous: model.Hash(value["previous"].(string)),
								Next:     model.Hash(value["next"].(string)),
							}
							if !b.Assign(&leaf) {
								return utils.NewNotAPointerError()
							}
						}
					}
				}
				if b.IsLeaf() {
					(*actualSize)++
					node.AddLeaf(b.GetLeaf())
				} else if b.IsNode() {
					r, _ := strconv.Atoi(remainder)
					node.AddNode(b.GetNode(), uint64(r))
				} else {
					r, _ := strconv.Atoi(remainder)
					node.AddBranch(&b, uint64(r))
				}
			}
		}
	}
	return nil
}
