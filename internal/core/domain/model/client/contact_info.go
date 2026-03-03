package client

import (
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/pkg/errs"
)

// ContactInfo holds a client's contact details.
type ContactInfo struct {
	phone     model.PhoneNumber
	email     string
	firstName string
	lastName  string
}

// NewContactInfo creates a ContactInfo value object with validation.
func NewContactInfo(phone model.PhoneNumber, email, firstName, lastName string) (ContactInfo, error) {
	if firstName == "" {
		return ContactInfo{}, errs.NewErrValueRequired("first name")
	}
	return ContactInfo{
		phone:     phone,
		email:     email,
		firstName: firstName,
		lastName:  lastName,
	}, nil
}

// MustNewContactInfo creates a ContactInfo or panics.
func MustNewContactInfo(phone model.PhoneNumber, email, firstName, lastName string) ContactInfo {
	ci, err := NewContactInfo(phone, email, firstName, lastName)
	if err != nil {
		panic(err)
	}
	return ci
}

// Phone returns the phone number.
func (ci ContactInfo) Phone() model.PhoneNumber { return ci.phone }

// Email returns the email address.
func (ci ContactInfo) Email() string { return ci.email }

// FirstName returns the first name.
func (ci ContactInfo) FirstName() string { return ci.firstName }

// LastName returns the last name.
func (ci ContactInfo) LastName() string { return ci.lastName }
