package outbox

import "github.com/mijgona/salon-crm/internal/pkg/ddd"

// EventRegistry maps event names to their types for deserialization.
type EventRegistry struct {
	events map[string]func() ddd.DomainEvent
}

// NewEventRegistry creates a new EventRegistry.
func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		events: make(map[string]func() ddd.DomainEvent),
	}
}

// Register registers an event constructor by name.
func (r *EventRegistry) Register(name string, constructor func() ddd.DomainEvent) {
	r.events[name] = constructor
}

// Get returns the constructor for the given event name.
func (r *EventRegistry) Get(name string) (func() ddd.DomainEvent, bool) {
	constructor, ok := r.events[name]
	return constructor, ok
}

// EventNames returns all registered event names.
func (r *EventRegistry) EventNames() []string {
	names := make([]string, 0, len(r.events))
	for name := range r.events {
		names = append(names, name)
	}
	return names
}
