package model

import (
	"errors"
	"regexp"
	"strings"

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
		err = errors.New("invalid hash string")
		return
	}
	return strings.ToLower(string(h)), nil
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
