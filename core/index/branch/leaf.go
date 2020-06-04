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

// Print ...
func (l *Leaf) Print() string {
	bytes, _ := json.Marshal(l)
	return string(bytes)
}
