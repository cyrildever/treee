package core_test

import (
	"testing"

	"github.com/cyrildever/treee/models/core"
	"gotest.tools/assert"
)

// TestHash ...
func TestHash(t *testing.T) {
	nonHash := "123"
	_, err := core.Hash(nonHash).String()
	assert.Error(t, err, "invalid hash string")

	valid := "123a"
	str, err := core.Hash(valid).String()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, str, valid)

	emptyString := core.Hash("")
	assert.Assert(t, emptyString.NonEmpty() == false)
	emptyBytes := core.ToHash(nil)
	assert.Assert(t, emptyBytes.NonEmpty() == false)

	lowercase := "abcd"
	uppercase := "ABCD"
	found1, _ := core.Hash(lowercase).String()
	found2, _ := core.Hash(uppercase).String()
	assert.Equal(t, found1, found2)
	assert.Equal(t, found2, lowercase)
}
