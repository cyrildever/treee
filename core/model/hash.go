package model

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/cyrildever/treee/core/exception"
	"github.com/cyrildever/treee/utils"
)

var regexHash = regexp.MustCompile(`^(?:[0-9a-fA-F]{2})*$`)

//--- TYPES

// Hash is the hexadecimal string representation of a hash, ie. a string of even length only made out of 0 to 9 and a to f characters.
type Hash string

// EmptyHash ...
const EmptyHash = Hash("")

//--- METHODS

// Bytes ...
func (h Hash) Bytes() ([]byte, error) {
	return utils.FromHex(string(h))
}

// String ...
func (h Hash) String() (str string, err error) {
	if !regexHash.MatchString(string(h)) {
		err = exception.NewInvalidHashStringError(string(h))
		return
	}
	return strings.ToLower(string(h)), nil
}

// IsEmpty ...
func (h Hash) IsEmpty() bool {
	return !h.NonEmpty()
}

// NonEmpty ...
func (h Hash) NonEmpty() bool {
	str, err := h.String()
	if err != nil {
		return false
	}
	return str != ""
}

//--- FUNCTIONS

// ToHash ...
func ToHash(bytes []byte) Hash {
	if bytes == nil {
		return Hash("")
	}
	return Hash(utils.ToHex(bytes))
}

// IsHashedString ...
func IsHashedString(input string) bool {
	return regexHash.MatchString(input)
}

// Hashes is an array of Hash.
type Hashes []Hash

func (h Hashes) Len() int           { return len(h) }
func (h Hashes) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h Hashes) Less(i, j int) bool { return h[i] < h[j] }

// Equals ...
func (h *Hashes) Equals(to *Hashes) bool {
	length := h.Len()
	if length != to.Len() {
		return false
	}
	for i := 0; i < length; i++ {
		if !reflect.DeepEqual((*h)[i], (*to)[i]) {
			return false
		}
	}
	return true
}

// Contains ...
func (h Hashes) Contains(item Hash) bool {
	for _, hash := range h {
		if hash == item {
			return true
		}
	}
	return false
}
