package exception

import (
	"fmt"
)

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
