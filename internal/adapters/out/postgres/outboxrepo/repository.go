package outboxrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres"
	"github.com/mijgona/salon-crm/internal/core/ports"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// PostgresOutboxRepository is a pgx + Squirrel implementation of OutboxRepository.
type PostgresOutboxRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresOutboxRepository creates a new repository.
func NewPostgresOutboxRepository(pool *pgxpool.Pool) *PostgresOutboxRepository {
	return &PostgresOutboxRepository{pool: pool}
}

func (r *PostgresOutboxRepository) Save(ctx context.Context, tx interface{}, event ddd.DomainEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	query, args, err := psql.Insert("outbox").
		Columns("id", "event_type", "payload", "created_at").
		Values(uuid.New(), event.GetName(), payload, time.Now()).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	return r.exec(ctx, tx, query, args...)
}

func (r *PostgresOutboxRepository) GetPending(ctx context.Context, limit int) ([]ports.OutboxEntry, error) {
	query, args, err := psql.Select("id", "event_type", "payload").
		From("outbox").
		Where("processed_at IS NULL").
		OrderBy("created_at ASC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ports.OutboxEntry
	for rows.Next() {
		var entry ports.OutboxEntry
		if err := rows.Scan(&entry.ID, &entry.EventType, &entry.Payload); err != nil {
			return nil, fmt.Errorf("scan outbox entry: %w", err)
		}
		result = append(result, entry)
	}
	return result, rows.Err()
}

func (r *PostgresOutboxRepository) MarkProcessed(ctx context.Context, id uuid.UUID) error {
	query, args, err := psql.Update("outbox").
		Set("processed_at", time.Now()).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	return err
}

func (r *PostgresOutboxRepository) exec(ctx context.Context, tx interface{}, sql string, args ...any) error {
	if t := postgres.ExtractTx(tx); t != nil {
		_, err := t.Exec(ctx, sql, args...)
		return err
	}
	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}
