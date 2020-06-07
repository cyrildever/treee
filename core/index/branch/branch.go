package branch

import "github.com/cyrildever/treee/utils"

//--- TYPES

// Branch ...
type Branch struct {
	nature interface{}
}

// Printable ...
type Printable interface {
	Print() string
}

//--- METHODS

// Assign ...
func (b *Branch) Assign(leafOrNodePtr interface{}) bool {
	if !utils.IsPointer(leafOrNodePtr) {
		return false
	}
	b.nature = leafOrNodePtr
	return true
}

// GetLeaf ...
func (b *Branch) GetLeaf() *Leaf {
	if l, ok := b.nature.(*Leaf); ok {
		return l
	}
	return &Leaf{}
}

// GetNode ...
func (b *Branch) GetNode() *Node {
	if n, ok := b.nature.(*Node); ok {
		return n
	}
	return &Node{}
}

// IsEmpty ...
func (b *Branch) IsEmpty() bool {
	return b.nature == nil || (!b.IsLeaf() && !b.IsNode())
}

// IsLeaf ...
func (b *Branch) IsLeaf() bool {
	nature := b.nature
	if _, ok := nature.(*Leaf); ok {
		return true
	}
	return false
}

// IsNode ...
func (b *Branch) IsNode() bool {
	nature := b.nature
	if _, ok := nature.(*Node); ok {
		return true
	}
	return false
}

// Print ...
func (b *Branch) Print() string {
	if b.IsLeaf() {
		l, _ := b.nature.(*Leaf)
		return l.Print()
	} else if b.IsNode() {
		n, _ := b.nature.(*Node)
		return n.Print()
	} else {
		return "{}"
	}
}
