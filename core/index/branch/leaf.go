package branch

import (
	"encoding/json"

	"github.com/cyrildever/treee/core/model"
)

//--- TYPES

// Leaf ...
type Leaf struct {
	ID       model.Hash `json:"id"`
	Position int64      `json:"position"`
	Size     int64      `json:"size"`
	Origin   model.Hash `json:"origin"`
	Previous model.Hash `json:"previous"`
	Next     model.Hash `json:"next"`
}

//--- METHODS

// IsEmpty ...
func (l *Leaf) IsEmpty() bool {
	return l.ID.IsEmpty() || l.Position == -1 || l.Size == 0
}

// Print ...
func (l *Leaf) Print() string {
	bytes, _ := json.Marshal(l)
	return string(bytes)
}

//--- FUNCTIONS

// NewEmptyLeaf ...
func NewEmptyLeaf() *Leaf {
	return &Leaf{
		ID:       model.EmptyHash,
		Position: -1,
		Size:     0,
		Origin:   model.EmptyHash,
		Previous: model.EmptyHash,
		Next:     model.EmptyHash,
	}
}
