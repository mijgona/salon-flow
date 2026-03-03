package model

import (
	"github.com/mijgona/salon-crm/internal/pkg/errs"

	"github.com/google/uuid"
)

// TenantID represents a tenant identifier (shared kernel).
type TenantID struct {
	value uuid.UUID
}

// NewTenantID creates a TenantID from a UUID.
func NewTenantID(value uuid.UUID) (TenantID, error) {
	if value == uuid.Nil {
		return TenantID{}, errs.NewErrValueRequired("tenant ID")
	}
	return TenantID{value: value}, nil
}

// MustNewTenantID creates a TenantID or panics.
func MustNewTenantID(value uuid.UUID) TenantID {
	t, err := NewTenantID(value)
	if err != nil {
		panic(err)
	}
	return t
}

// UUID returns the underlying UUID.
func (t TenantID) UUID() uuid.UUID { return t.value }

// String returns the string representation.
func (t TenantID) String() string { return t.value.String() }

// Equal checks value equality.
func (t TenantID) Equal(other TenantID) bool { return t.value == other.value }
