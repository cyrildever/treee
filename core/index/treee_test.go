package index_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/cyrildever/treee/core/index"
	"github.com/cyrildever/treee/core/index/branch"
	"github.com/cyrildever/treee/core/model"
	"gotest.tools/assert"
)

// TestTreee ...
func TestTreee(t *testing.T) {
	treee, err := index.New(index.INIT_PRIME)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, treee.Size(), uint64(0))

	firstLeaf := branch.Leaf{
		ID:       model.Hash("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
		Position: 0,
		Size:     100,
	}
	err = treee.Add(firstLeaf)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, treee.Size(), uint64(1))

	found, err := treee.Search(firstLeaf.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, found.ID, firstLeaf.ID)
	assert.Equal(t, found.Position, firstLeaf.Position)
	assert.Equal(t, found.Size, firstLeaf.Size)

	secondLeaf := branch.Leaf{
		ID:       model.Hash("fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321"),
		Position: 100,
		Size:     50,
	}
	treee.Add(secondLeaf)
	assert.Equal(t, treee.Size(), uint64(2))

	thirdLeaf := branch.Leaf{
		ID:       model.Hash("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"),
		Position: 150,
		Size:     10,
	}
	treee.Add(thirdLeaf)
	assert.Equal(t, treee.Size(), uint64(3))

	fmt.Println(treee.PrintAll(true)) // TODO Only prints if tests failed

	found, err = treee.Search(firstLeaf.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, found.ID, firstLeaf.ID)
	assert.Equal(t, found.Position, firstLeaf.Position)
	assert.Equal(t, found.Size, firstLeaf.Size)

	found, err = treee.Search(thirdLeaf.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, found.ID, thirdLeaf.ID)
	assert.Equal(t, found.Position, thirdLeaf.Position)
	assert.Equal(t, found.Size, thirdLeaf.Size)

	found, err = treee.Search(secondLeaf.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, found.ID, secondLeaf.ID)
	assert.Equal(t, found.Position, secondLeaf.Position)
	assert.Equal(t, found.Size, secondLeaf.Size)
}

// TestLoad ...
func TestLoad(t *testing.T) {
	pwd, _ := os.Getwd()
	path := strings.Replace(pwd, "core"+string(os.PathSeparator)+"index", "", 1) // TODO Change if path to test changes
	treee, err := index.Load(path + "saved" + string(os.PathSeparator) + "test-treee.json")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, treee.Size(), uint64(3))
}
