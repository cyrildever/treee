package exception

import (
	"fmt"
)

//--- ERRORS

// InvalidHashStringError ...
type InvalidHashStringError struct {
	message string
}

func (e InvalidHashStringError) Error() string {
	return e.message
}

// NewInvalidHashStringError ...
func NewInvalidHashStringError(str string) *InvalidHashStringError {
	return &InvalidHashStringError{
		message: fmt.Sprintf("invalid hash string: %s", str),
	}
}

// InvalidUUIDError ...
type InvalidUUIDError struct {
	message string
}

func (e InvalidUUIDError) Error() string {
	return e.message
}

// NewInvalidUUIDError ...
func NewInvalidUUIDError(str string) *InvalidUUIDError {
	return &InvalidUUIDError{
		message: fmt.Sprintf("invalid UUID string: %s", str),
	}
}
