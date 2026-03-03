package errs

import "fmt"

// ErrValueRequired indicates a required value is missing.
type ErrValueRequired struct {
	Field string
}

func NewErrValueRequired(field string) *ErrValueRequired {
	return &ErrValueRequired{Field: field}
}

func (e *ErrValueRequired) Error() string {
	return fmt.Sprintf("%s is required", e.Field)
}
