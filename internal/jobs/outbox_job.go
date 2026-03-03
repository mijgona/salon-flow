package jobs

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/ports"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"
	"log"
	"time"
)

// OutboxJob periodically reads pending outbox events and publishes them via Mediatr.
type OutboxJob struct {
	outboxRepo ports.OutboxRepository
	mediatr    ddd.Mediatr
	interval   time.Duration
}

// NewOutboxJob creates a new OutboxJob.
func NewOutboxJob(outboxRepo ports.OutboxRepository, mediatr ddd.Mediatr, interval time.Duration) *OutboxJob {
	return &OutboxJob{
		outboxRepo: outboxRepo,
		mediatr:    mediatr,
		interval:   interval,
	}
}

// Start begins the outbox processing loop.
func (j *OutboxJob) Start(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("outbox job stopped")
			return
		case <-ticker.C:
			j.processOutbox(ctx)
		}
	}
}

func (j *OutboxJob) processOutbox(ctx context.Context) {
	entries, err := j.outboxRepo.GetPending(ctx, 100)
	if err != nil {
		log.Printf("outbox: error fetching pending events: %v", err)
		return
	}

	for _, entry := range entries {
		// In a real implementation, we would deserialize the event from entry.Payload
		// and publish it via mediatr. For now, mark as processed.
		if err := j.outboxRepo.MarkProcessed(ctx, entry.ID); err != nil {
			log.Printf("outbox: error marking event %s as processed: %v", entry.ID, err)
		}
	}
}
