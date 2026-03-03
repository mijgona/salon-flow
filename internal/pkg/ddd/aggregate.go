package ddd

// BaseAggregate provides domain event support for aggregate roots.
type BaseAggregate[ID comparable] struct {
	baseEntity   *BaseEntity[ID]
	domainEvents []DomainEvent
}

// NewBaseAggregate creates a new BaseAggregate with the given ID.
func NewBaseAggregate[ID comparable](id ID) *BaseAggregate[ID] {
	return &BaseAggregate[ID]{
		baseEntity:   NewBaseEntity[ID](id),
		domainEvents: make([]DomainEvent, 0),
	}
}

// ID returns the aggregate's identifier.
func (ba *BaseAggregate[ID]) ID() ID {
	return ba.baseEntity.ID()
}

// RaiseDomainEvent appends a domain event to the aggregate.
func (ba *BaseAggregate[ID]) RaiseDomainEvent(event DomainEvent) {
	ba.domainEvents = append(ba.domainEvents, event)
}

// GetDomainEvents returns all raised domain events.
func (ba *BaseAggregate[ID]) GetDomainEvents() []DomainEvent {
	return ba.domainEvents
}

// ClearDomainEvents removes all domain events from the aggregate.
func (ba *BaseAggregate[ID]) ClearDomainEvents() {
	ba.domainEvents = []DomainEvent{}
}
