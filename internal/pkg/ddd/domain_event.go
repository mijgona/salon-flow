package ddd

import "github.com/google/uuid"

// DomainEvent represents an event that occurred in the domain.
type DomainEvent interface {
	GetID() uuid.UUID
	GetName() string
}
