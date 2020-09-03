package exception

import (
	"fmt"
)

//--- ERRORS

// AlreadyExistsInIndexError ...
type AlreadyExistsInIndexError struct {
	message string
}

func (e AlreadyExistsInIndexError) Error() string {
	return e.message
}

// NewAlreadyExistsInIndexError ...
func NewAlreadyExistsInIndexError(id string) *AlreadyExistsInIndexError {
	return &AlreadyExistsInIndexError{
		message: fmt.Sprintf("item already exists in the index: %s", id),
	}
}

// EmptyItemError ...
type EmptyItemError struct {
	message string
}

func (e EmptyItemError) Error() string {
	return e.message
}

// NewEmptyItemError ...
func NewEmptyItemError() *EmptyItemError {
	return &EmptyItemError{
		message: "empty item",
	}
}

// IncoherentSizeError ...
type IncoherentSizeError struct {
	message string
}

func (e IncoherentSizeError) Error() string {
	return e.message
}

// NewIncoherentSizeError ...
func NewIncoherentSizeError(expected, found int) *IncoherentSizeError {
	return &IncoherentSizeError{
		message: fmt.Sprintf("declared size [%d] not equal to actual size [%d]", expected, found),
	}
}

// LoopError ...
type LoopError struct {
	message string
}

func (e LoopError) Error() string {
	return e.message
}

// NewLoopError ...
func NewLoopError(operation string) *LoopError {
	return &LoopError{
		message: fmt.Sprintf("looping without %s item", operation),
	}
}

// NotFoundError ...
type NotFoundError struct {
	message string
}

func (e NotFoundError) Error() string {
	return e.message
}

// NewNotFoundError ...
func NewNotFoundError(itemID string) *NotFoundError {
	return &NotFoundError{
		message: fmt.Sprintf("nothing found for ID: %s", itemID),
	}
}

// NotAValidTreeeError ...
type NotAValidTreeeError struct {
	message string
}

func (e NotAValidTreeeError) Error() string {
	return e.message
}

// NewNotAValidTreeeError ...
func NewNotAValidTreeeError() *NotAValidTreeeError {
	return &NotAValidTreeeError{
		message: "not a valid treee",
	}
}
