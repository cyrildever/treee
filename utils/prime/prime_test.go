package prime_test

import (
	"testing"

	"github.com/cyrildever/treee/utils/prime"
	"gotest.tools/assert"
)

// TestIsPrime ...
func TestIsPrime(t *testing.T) {
	var number uint64 = 11
	assert.Assert(t, prime.IsPrime(number))

	number = 54
	assert.Assert(t, prime.IsPrime(number) == false)

	number = 999331
	assert.Assert(t, prime.IsPrime(number) == false) // Which is actually not right: we have this result only because it's over our current testing range
}

// TestNext ...
func TestNext(t *testing.T) {
	var number uint64 = 9
	p, err := prime.Next(number)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, p, uint64(11))

	number = 11
	p, _ = prime.Next(number)
	assert.Equal(t, p, uint64(13))

	number = 8000
	_, err = prime.Next(number)
	assert.Error(t, err, "number is too high")
}
