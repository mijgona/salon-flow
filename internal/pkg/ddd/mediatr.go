package ddd

import (
	"context"
	"sync"
)

// EventHandler handles a domain event.
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
}

// Mediatr is an in-process event dispatcher.
type Mediatr interface {
	Subscribe(handler EventHandler, events ...DomainEvent)
	Publish(ctx context.Context, event DomainEvent) error
}

// InProcessMediatr is a synchronous in-process implementation of Mediatr.
type InProcessMediatr struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}

// NewInProcessMediatr creates a new InProcessMediatr.
func NewInProcessMediatr() *InProcessMediatr {
	return &InProcessMediatr{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe registers an event handler for the given event types.
func (m *InProcessMediatr) Subscribe(handler EventHandler, events ...DomainEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, event := range events {
		name := event.GetName()
		m.handlers[name] = append(m.handlers[name], handler)
	}
}

// Publish dispatches a domain event to all subscribed handlers.
func (m *InProcessMediatr) Publish(ctx context.Context, event DomainEvent) error {
	m.mu.RLock()
	handlers := m.handlers[event.GetName()]
	m.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
