package search

import (
	"github.com/cyrildever/treee/core/index/branch"
	"github.com/cyrildever/treee/core/model"
)

// Engine ...
type Engine func(model.Hash) (*branch.Leaf, error)
