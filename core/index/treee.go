package index

import (
	"encoding/json"
	"errors"
	"math/big"
	"strconv"
	"sync"

	"github.com/cyrildever/treee/common/logger"
	"github.com/cyrildever/treee/core/exception"
	"github.com/cyrildever/treee/core/index/branch"
	"github.com/cyrildever/treee/core/model"
	"github.com/cyrildever/treee/utils"
)

// Current is the current Treee index used when running the executable app
var Current *Treee

const (
	// INIT_PRIME ...
	INIT_PRIME uint64 = 2 // TODO Shouldn't be used in production
)

//--- TYPES

// Treee is an instance of a database tree
type Treee struct {
	initPrime uint64
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
		existingOrigin, err := t.Search(item.Origin)
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

	str := `{"initPrime":` + strconv.FormatUint(t.initPrime, 10) + `,"trunk":` + t.trunk.Print() + `,"size":` + strconv.FormatUint(t.size, 10) + "}"
	if beautify {
		var js interface{}
		json.Unmarshal([]byte(str), &js)
		bytes, _ := json.MarshalIndent(js, "", "  ")
		str = string(bytes)
	}
	return str
}

// Search ...
func (t *Treee) Search(ID model.Hash) (found *branch.Leaf, err error) {
	t.RLock()
	defer t.RUnlock()

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

// New ...
func New(initPrime uint64) (t Treee, err error) {
	if utils.IsPrime(initPrime) {
		t = Treee{
			initPrime: initPrime,
			trunk:     branch.NewNode(initPrime),
		}
	} else if initPrime == 0 {
		t = Treee{
			initPrime: INIT_PRIME,
			trunk:     branch.NewNode(INIT_PRIME),
		}
	} else {
		err = errors.New("invalid prime number")
	}
	return
}
