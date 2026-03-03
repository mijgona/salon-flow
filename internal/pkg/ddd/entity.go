package ddd

// BaseEntity provides generic identity for domain entities.
type BaseEntity[ID comparable] struct {
	id ID
}

// NewBaseEntity creates a new BaseEntity with the given ID.
func NewBaseEntity[ID comparable](id ID) *BaseEntity[ID] {
	return &BaseEntity[ID]{id: id}
}

// ID returns the entity's identifier.
func (be *BaseEntity[ID]) ID() ID {
	return be.id
}

// Equal checks identity equality with another BaseEntity.
func (be *BaseEntity[ID]) Equal(other *BaseEntity[ID]) bool {
	return other != nil && be.id == other.id
}
