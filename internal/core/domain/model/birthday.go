package model

import (
	"github.com/mijgona/salon-crm/internal/pkg/errs"
	"time"
)

// Birthday represents a date of birth.
type Birthday struct {
	value time.Time
}

// NewBirthday creates a Birthday value object. The date must be in the past.
func NewBirthday(value time.Time) (Birthday, error) {
	if value.IsZero() {
		return Birthday{}, errs.NewErrValueRequired("birthday")
	}
	if value.After(time.Now()) {
		return Birthday{}, errs.NewErrValueMustBe("birthday", "in the past")
	}
	return Birthday{value: value}, nil
}

// MustNewBirthday creates a Birthday or panics.
func MustNewBirthday(value time.Time) Birthday {
	b, err := NewBirthday(value)
	if err != nil {
		panic(err)
	}
	return b
}

// Time returns the underlying time.Time.
func (b Birthday) Time() time.Time { return b.value }

// IsZero returns true if the birthday is unset.
func (b Birthday) IsZero() bool { return b.value.IsZero() }

// Equal checks value equality (date only).
func (b Birthday) Equal(other Birthday) bool {
	return b.value.Year() == other.value.Year() &&
		b.value.Month() == other.value.Month() &&
		b.value.Day() == other.value.Day()
}
