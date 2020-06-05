package errors

import (
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

//--- TYPES

// WrongSchema ...
type WrongSchema struct {
	schemaErrors []gojsonschema.ResultError
}

//--- METHODS

// Error ...
func (err WrongSchema) Error() string {
	var errs []string
	for _, se := range err.schemaErrors {
		errs = append(errs, se.String())
	}
	return strings.Join(errs, ";")
}

//--- FUNCTIONS

// NewWrongSchemaError ...
func NewWrongSchemaError(errs []gojsonschema.ResultError) *WrongSchema {
	return &WrongSchema{
		schemaErrors: errs,
	}
}
