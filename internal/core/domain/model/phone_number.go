package model

import (
	"github.com/mijgona/salon-crm/internal/pkg/errs"
	"regexp"
)

// russianPhoneRegex matches Russian phone numbers: +7XXXXXXXXXX
var russianPhoneRegex = regexp.MustCompile(`^\+7\d{10}$`)

// PhoneNumber represents a validated Russian phone number.
type PhoneNumber struct {
	value string
}

// NewPhoneNumber creates a PhoneNumber with Russian format validation.
func NewPhoneNumber(value string) (PhoneNumber, error) {
	if value == "" {
		return PhoneNumber{}, errs.NewErrValueRequired("phone number")
	}
	if !russianPhoneRegex.MatchString(value) {
		return PhoneNumber{}, errs.NewErrValueMustBe("phone number", "in Russian format (+7XXXXXXXXXX)")
	}
	return PhoneNumber{value: value}, nil
}

// MustNewPhoneNumber creates a PhoneNumber or panics.
func MustNewPhoneNumber(value string) PhoneNumber {
	p, err := NewPhoneNumber(value)
	if err != nil {
		panic(err)
	}
	return p
}

// String returns the phone number string.
func (p PhoneNumber) String() string { return p.value }

// Equal checks value equality.
func (p PhoneNumber) Equal(other PhoneNumber) bool { return p.value == other.value }
