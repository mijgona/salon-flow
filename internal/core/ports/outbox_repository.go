package ports

import (
	"context"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"

	"github.com/google/uuid"
)

// OutboxRepository stores domain events for reliable publishing.
type OutboxRepository interface {
	Save(ctx context.Context, tx interface{}, event ddd.DomainEvent) error
	GetPending(ctx context.Context, limit int) ([]OutboxEntry, error)
	MarkProcessed(ctx context.Context, id uuid.UUID) error
}

// OutboxEntry represents a stored outbox event.
type OutboxEntry struct {
	ID        uuid.UUID
	EventType string
	Payload   []byte
}
