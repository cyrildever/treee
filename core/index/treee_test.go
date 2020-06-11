package index_test

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cyrildever/treee/core/index"
	"github.com/cyrildever/treee/core/index/branch"
	"github.com/cyrildever/treee/core/index/search"
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

// TestLastOrSearch ...
func TestLastOrSearch(t *testing.T) {
	treee, _ := index.New(index.INIT_PRIME)
	id := model.Hash("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	firstLeaf := branch.Leaf{
		ID:       id,
		Position: 0,
		Size:     100,
	}
	treee.Add(firstLeaf)
	secondLeaf := branch.Leaf{
		ID:       model.Hash("fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321"),
		Position: 100,
		Size:     50,
		Previous: firstLeaf.ID,
	}
	treee.Add(secondLeaf)

	var engine search.Engine
	for i := 0; i < 2; i++ {
		if i%2 == 0 {
			engine = treee.Last
		} else {
			engine = treee.Search
		}
		found, err := engine(id)
		if err != nil {
			t.Fatal(err)
		}
		if i%2 == 0 {
			assert.Equal(t, found.ID, secondLeaf.ID, "Last should find second leaf")
		} else {
			assert.Equal(t, found.ID, firstLeaf.ID, "Search should find first leaf")
		}
	}
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

// TestScalability ...
func TestScalability(t *testing.T) {
	var wg sync.WaitGroup
	rounds := 10000
	t0 := time.Now().UnixNano()
	treee, _ := index.New(101)
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int64, treee *index.Treee) {
			defer wg.Done()
			leaf := branch.Leaf{
				ID:       model.Hash(fmt.Sprintf("%0x", strconv.FormatInt(i, 16))),
				Position: i,
				Size:     1,
			}
			err := treee.Add(leaf)
			if err != nil {
				fmt.Println(err)
			}
		}(&wg, int64(i), treee)
	}
	wg.Wait()
	t1 := time.Now().UnixNano()
	fmt.Printf("inserting %d leaves completed in %d ms\n", rounds, (t1-t0)/int64(time.Millisecond))
	assert.Equal(t, treee.Size(), uint64(rounds))

	var engine search.Engine
	var resp = make(chan branch.Leaf, rounds*100)
	t0 = time.Now().UnixNano()
	for i := 0; i < rounds*100; i++ {
		wg.Add(1)
		if i%2 == 0 {
			engine = treee.Last
		} else {
			engine = treee.Search
		}
		go func(wg *sync.WaitGroup, i int64, engine search.Engine, treee *index.Treee, resp chan branch.Leaf) {
			defer wg.Done()
			id := model.Hash(fmt.Sprintf("%0x", strconv.FormatInt(int64(rand.Intn(rounds)), 16)))
			found, err := engine(id)
			if err != nil {
				fmt.Println("i", i, "error", err)
			}
			resp <- *found
		}(&wg, int64(i), engine, treee, resp)
	}
	wg.Wait()
	t1 = time.Now().UnixNano()
	fmt.Printf("searching %d ids completed in %d ms\n", rounds*100, (t1-t0)/int64(time.Millisecond))
	results := 0
	close(resp)
	for range resp {
		results++
	}
	assert.Equal(t, results, rounds*100)

	// assert.Assert(t, false) // TODO Uncomment to get performance logs
}
