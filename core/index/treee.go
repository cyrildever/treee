package index

import (
	"encoding/json"
	"errors"
	"fmt"
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
	trunk *branch.Node
	size  uint64
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
		existingPrevious, err := t.Search(item.Previous)
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
			targetBranch.Assign(item)
			previous.Next = item.ID
			origin.Previous = item.ID
			t.size++
			return nil
		} else if targetBranch.IsLeaf() {
			existingLeaf := targetBranch.GetLeaf()
			nextPrime, err := utils.NextPrime(currentStage.Uint64())
			if err != nil {
				log.Crit("Unable to get next prime number", "error", err)
				return err
			}
			newNode := branch.NewNode(nextPrime)
			newNode.AddLeaf(existingLeaf)
			newNode.AddLeaf(&item)
			targetBranch.Assign(*newNode)
			previous.Next = item.ID
			origin.Previous = item.ID
			t.size++
			return nil
		} else if targetBranch.IsNode() {
			currentNode = targetBranch.GetNode()
		}
	}
}

// PrintAll ...
// Use with caution!
func (t *Treee) PrintAll(beautify bool) string {
	t.RLock()
	defer t.RUnlock()

	str := `{"initPrime":` + strconv.FormatUint(t.InitPrime, 10) + `,"trunk":` + t.trunk.Print() + `,"size":` + strconv.FormatUint(t.size, 10) + "}"
	if beautify {
		var js interface{}
		json.Unmarshal([]byte(str), &js)
		bytes, _ := json.MarshalIndent(js, "", "  ")
		str = string(bytes)
	}
	return str
}

// Save ...
func (t *Treee) Save() {
	log := logger.Init("index", "Save")
	conf, _ := config.GetConfig()

	if !saving && !conf.IsTestEnvironment() {
		saving = true
		t0 := time.Now().UnixNano()
		f, err := os.Create("./saved/treee.json")
		if err != nil {
			log.Crit("Unable to create file", "error", err)
			saving = false
			return
		}
		size := t.Size()
		n, err := f.WriteString(t.PrintAll(false))
		if err != nil {
			t1 := time.Now().UnixNano()
			log.Error("An error occurred while saving the index", "error", err, "after", t1-t0)
			saving = false
			return
		}
		t1 := time.Now().UnixNano()
		log.Info("Index saved", "size", size, "bytes", n, "duration", t1-t0)
		f.Close()
		saving = false
		return
	}
}

// Search ...
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
		err = errors.New("not a valid treee")
		return
	}
	if sp, ok := stagePrime.(float64); !ok || uint64(sp) != st.InitPrime || !utils.IsPrime(uint64(sp)) {
		err = fmt.Errorf("not a valid stage prime number: %f", sp)
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
		return &treee, fmt.Errorf("declared size [%d] not equal to actual size [%d]", st.Size, actualSize)
	}

	return &treee, nil
}

// New ...
func New(initPrime uint64) (t *Treee, err error) {
	if utils.IsPrime(initPrime) {
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
		err = errors.New("invalid prime number")
	}
	return
}

func parse(node *branch.Node, arr []interface{}, actualSize *int) error {
	for _, v := range arr {
		if val, ok := v.(map[string]interface{}); ok {
			for prime, child := range val {
				b := branch.Branch{}
				if value, ok := child.(map[string]interface{}); ok {
					if len(value) == 0 {
						idx, err := strconv.Atoi(prime)
						if err != nil {
							return err
						}
						node.AddBranch(&b, uint64(idx))
						continue
					}
					if n, ok := value["children"]; ok {
						stagePrime, ok := value["stagePrime"]
						if !ok {
							return errors.New("not a valid treee")
						}
						sp, ok := stagePrime.(float64)
						if !ok || !utils.IsPrime(uint64(sp)) {
							return fmt.Errorf("not a valid stage prime number: %f", sp)
						}
						newNode := branch.NewNode(uint64(sp))
						err := parse(newNode, n.([]interface{}), actualSize)
						if err != nil {
							return err
						}
						b.Assign(newNode)
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
							b.Assign(leaf)
						}
					}
				}
				if b.IsLeaf() {
					(*actualSize)++
					node.AddLeaf(b.GetLeaf())
				} else if b.IsNode() {
					p, _ := strconv.Atoi(prime)
					node.AddNode(b.GetNode(), uint64(p))
				} else {
					p, _ := strconv.Atoi(prime)
					node.AddBranch(&b, uint64(p))
				}
			}
		}
	}
	return nil
}
