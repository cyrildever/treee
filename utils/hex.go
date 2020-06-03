package utils

import (
	"encoding/hex"
)

// FromHex tries to convert an hexadecimal representation of a value to its corresponding byte array
func FromHex(input string) ([]byte, error) {
	return hex.DecodeString(input)
}

// ToHex converts a byte array to its string representation in hexadecimal
func ToHex(input []byte) string {
	return hex.EncodeToString(input)
}
