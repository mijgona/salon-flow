package ddd

import "github.com/google/uuid"

// AggregateRoot is a type alias for BaseAggregate with UUID identity.
type AggregateRoot = BaseAggregate[uuid.UUID]

// NewAggregateRoot creates a new AggregateRoot with the given UUID.
func NewAggregateRoot(id uuid.UUID) *AggregateRoot {
	return NewBaseAggregate[uuid.UUID](id)
}
