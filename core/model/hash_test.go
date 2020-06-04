package model_test

import (
	"testing"

	"github.com/cyrildever/treee/core/model"
	"gotest.tools/assert"
)

// TestHash ...
func TestHash(t *testing.T) {
	nonHash := "123"
	_, err := model.Hash(nonHash).String()
	assert.Error(t, err, "invalid hash string")

	valid := "123a"
	str, err := model.Hash(valid).String()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, str, valid)

	emptyHash := model.EmptyHash
	assert.Assert(t, emptyHash.NonEmpty() == false)
	emptyBytes := model.ToHash(nil)
	assert.Assert(t, emptyBytes.NonEmpty() == false)

	lowercase := "abcd"
	uppercase := "ABCD"
	found1, _ := model.Hash(lowercase).String()
	found2, _ := model.Hash(uppercase).String()
	assert.Equal(t, found1, found2)
	assert.Equal(t, found2, lowercase)
}
