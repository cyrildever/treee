package model

import (
	"regexp"
	"strings"

	"github.com/cyrildever/treee/core/exception"
	"github.com/gofrs/uuid"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

//--- TYPES

// UUID is the string representation of a UUID, eg 'e05572b3-230a-45fd-a779-604c2b8ceb24'
type UUID string

//--- METHODS

// Bytes ...
func (u UUID) Bytes() ([]byte, error) {
	if _, err := u.String(); err != nil {
		return nil, err
	}
	return []byte(u), nil
}

// String ...
func (u UUID) String() (str string, err error) {
	if string(u) == "" || !uuidRegex.MatchString(string(u)) {
		err = exception.NewInvalidUUIDError(string(u))
		return
	}
	return strings.ToLower(string(u)), nil
}

// IsEmpty ...
func (u UUID) IsEmpty() bool {
	return !u.NonEmpty()
}

// NonEmpty ...
func (u UUID) NonEmpty() bool {
	str, err := u.String()
	if err != nil {
		return false
	}
	return str != ""
}

//--- FUNCTIONS

// ToUUID ...
func ToUUID(bytes []byte) UUID {
	if bytes == nil {
		return UUID("")
	}
	return UUID(string(bytes))
}

// GenerateUUID ...
func GenerateUUID() (id UUID, err error) {
	u, err := uuid.NewV4()
	if err != nil {
		return
	}
	id = ToUUID([]byte(u.String()))
	return
}

// IsUUIDString ...
func IsUUIDString(input string) bool {
	return input != "" && uuidRegex.MatchString(input)
}
