package utils

import (
	"reflect"
)

//--- FUNCTIONS

// IsPointer returns `true` if the passed item is a pointer, not a value
func IsPointer(item interface{}) bool {
	return reflect.ValueOf(item).Kind() == reflect.Ptr
}

// IsValue returns `true` if the passed item is a value, not a pointer
func IsValue(item interface{}) bool {
	return !IsPointer(item)
}

//--- ERRORS

// NotAPointerError ...
type NotAPointerError struct {
	message string
}

func (e NotAPointerError) Error() string {
	return e.message
}

// NewNotAPointerError ...
func NewNotAPointerError() *NotAPointerError {
	return &NotAPointerError{
		message: "not a pointer",
	}
}
