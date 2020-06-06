package branch

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/cyrildever/treee/utils"
)

//--- TYPES

// Node ...
type Node struct {
	StagePrime uint64
	children   map[uint64]*Branch
}

//--- METHODS

// AddBranch ...
func (n *Node) AddBranch(item *Branch, idx uint64) bool {
	if _, exists := n.children[idx]; !exists {
		n.children[idx] = item
		return true
	}
	return false
}

// AddLeaf ...
func (n *Node) AddLeaf(item *Leaf) bool {
	idStr, err := item.ID.String()
	if err != nil {
		return false
	}
	id := new(big.Int)
	id.SetString(idStr, 16)
	currentStage := new(big.Int)
	currentStage.SetUint64(n.StagePrime)
	modulo := new(big.Int)
	modulo = modulo.Mod(id, currentStage)
	idx := modulo.Uint64()
	existing, exists := n.children[idx]
	if !exists || existing.IsEmpty() {
		newBranch := Branch{}
		newBranch.Assign(*item)
		n.children[idx] = &newBranch
		return true
	} else if existing.IsLeaf() {
		nextPrime, err := utils.NextPrime(n.StagePrime)
		if err != nil {
			// TODO Log and warn because it could mean that we are beyond the first 1000 prime numbers
			return false
		}
		newNode := NewNode(nextPrime)
		newNode.AddLeaf(existing.GetLeaf())
		newNode.AddLeaf(item)
		newBranch := Branch{}
		newBranch.Assign(*newNode)
		n.children[idx] = &newBranch
	}
	return false
}

// AddNode ...
func (n *Node) AddNode(item *Node, idx uint64) bool {
	if _, exists := n.children[idx]; !exists {
		newBranch := Branch{}
		newBranch.Assign(*item)
		n.children[idx] = &newBranch
		return true
	}
	return false
}

// ChildAt ...
func (n *Node) ChildAt(idx uint64) (b *Branch, exists bool) {
	if br, ok := n.children[idx]; ok {
		return br, true
	}
	return
}

// Print ...
func (n *Node) Print() string {
	str := `{"stagePrime":` + strconv.FormatUint(n.StagePrime, 10) + ","
	if len(n.children) > 0 {
		str += `"children":[`
		children := []string{}
		for i, b := range n.children {
			children = append(children, `{"`+strconv.FormatUint(i, 10)+`": `+b.Print()+"}")
		}
		str += strings.Join(children, ",")
		str += "]"
	} else {
		str += `"children":[]`
	}
	str += "}"
	return str
}

//--- FUNCTIONS

// NewNode ...
func NewNode(stagePrime uint64) *Node {
	children := make(map[uint64]*Branch, int(stagePrime))
	for i := uint64(0); i < stagePrime; i++ {
		children[i] = &Branch{}
	}
	return &Node{
		StagePrime: stagePrime,
		children:   children,
	}
}
