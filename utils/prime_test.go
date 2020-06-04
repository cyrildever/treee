package utils_test

import (
	"testing"

	"github.com/cyrildever/treee/utils"
	"gotest.tools/assert"
)

// TestIsPrime ...
func TestIsPrime(t *testing.T) {
	var number uint64 = 11
	assert.Assert(t, utils.IsPrime(number))

	number = 54
	assert.Assert(t, utils.IsPrime(number) == false)

	number = 999331
	assert.Assert(t, utils.IsPrime(number) == true) // Which is actually true
	number = 999332
	assert.Assert(t, utils.IsPrime(number) == true) // Which is obviously false
}

// TestNextPrime ...
func TestNextPrime(t *testing.T) {
	var number uint64 = 9
	prime, err := utils.NextPrime(number)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, prime, uint64(11))

	number = 11
	prime, _ = utils.NextPrime(number)
	assert.Equal(t, prime, uint64(13))

	number = 8000
	_, err = utils.NextPrime(number)
	assert.Error(t, err, "number is too high")
}
