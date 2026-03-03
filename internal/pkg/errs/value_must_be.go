package errs

import "fmt"

// ErrValueMustBe indicates a value does not satisfy a constraint.
type ErrValueMustBe struct {
	Field      string
	Constraint string
}

func NewErrValueMustBe(field, constraint string) *ErrValueMustBe {
	return &ErrValueMustBe{Field: field, Constraint: constraint}
}

func (e *ErrValueMustBe) Error() string {
	return fmt.Sprintf("%s must be %s", e.Field, e.Constraint)
}
