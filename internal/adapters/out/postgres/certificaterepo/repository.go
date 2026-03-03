package certificaterepo

import (
	"context"
	"fmt"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/certificate"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// PostgresCertificateRepository is a pgx + Squirrel implementation of CertificateRepository.
type PostgresCertificateRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresCertificateRepository creates a new repository.
func NewPostgresCertificateRepository(pool *pgxpool.Pool) *PostgresCertificateRepository {
	return &PostgresCertificateRepository{pool: pool}
}

func (r *PostgresCertificateRepository) Add(ctx context.Context, tx interface{}, c *certificate.Certificate) error {
	query, args, err := psql.Insert("certificates").
		Columns(
			"id", "tenant_id", "purchased_by", "activated_by",
			"balance_amount", "balance_currency", "status",
			"activated_at", "expires_at", "created_at",
		).
		Values(
			c.ID(), c.TenantID().UUID(), c.PurchasedBy(),
			nullableUUID(c.ActivatedBy()),
			c.Balance().Amount(), c.Balance().Currency(),
			string(c.Status()),
			nullableTimestamp(c.ActivatedAt()),
			c.ExpiresAt(), c.CreatedAt(),
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	return r.exec(ctx, tx, query, args...)
}

func (r *PostgresCertificateRepository) Update(ctx context.Context, tx interface{}, c *certificate.Certificate) error {
	query, args, err := psql.Update("certificates").
		Set("activated_by", nullableUUID(c.ActivatedBy())).
		Set("balance_amount", c.Balance().Amount()).
		Set("status", string(c.Status())).
		Set("activated_at", nullableTimestamp(c.ActivatedAt())).
		Where(sq.Eq{"id": c.ID()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	return r.exec(ctx, tx, query, args...)
}

func (r *PostgresCertificateRepository) Get(ctx context.Context, tx interface{}, id uuid.UUID) (*certificate.Certificate, error) {
	query, args, err := psql.Select(
		"id", "tenant_id", "purchased_by", "activated_by",
		"balance_amount", "balance_currency", "status",
		"activated_at", "expires_at", "created_at",
	).From("certificates").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	rows, err := r.query(ctx, tx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	return scanCertificateRow(rows)
}

// --- internal helpers ---

func (r *PostgresCertificateRepository) exec(ctx context.Context, tx interface{}, sql string, args ...any) error {
	if t := postgres.ExtractTx(tx); t != nil {
		_, err := t.Exec(ctx, sql, args...)
		return err
	}
	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *PostgresCertificateRepository) query(ctx context.Context, tx interface{}, sql string, args ...any) (pgx.Rows, error) {
	if t := postgres.ExtractTx(tx); t != nil {
		return t.Query(ctx, sql, args...)
	}
	return r.pool.Query(ctx, sql, args...)
}

func scanCertificateRow(rows pgx.Rows) (*certificate.Certificate, error) {
	var (
		id              uuid.UUID
		tenantID        uuid.UUID
		purchasedBy     uuid.UUID
		activatedBy     *uuid.UUID
		balanceAmount   decimal.Decimal
		balanceCurrency string
		status          string
		activatedAt     *time.Time
		expiresAt       time.Time
		createdAt       time.Time
	)

	err := rows.Scan(
		&id, &tenantID, &purchasedBy, &activatedBy,
		&balanceAmount, &balanceCurrency, &status,
		&activatedAt, &expiresAt, &createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan certificate: %w", err)
	}

	tid := model.MustNewTenantID(tenantID)
	balance := model.MustNewMoney(balanceAmount, balanceCurrency)

	var actBy uuid.UUID
	if activatedBy != nil {
		actBy = *activatedBy
	}

	var actAt time.Time
	if activatedAt != nil {
		actAt = *activatedAt
	}

	return certificate.RestoreCertificate(
		id, tid, purchasedBy, actBy,
		balance, certificate.CertificateStatus(status),
		actAt, expiresAt, createdAt,
	), nil
}

func nullableUUID(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

func nullableTimestamp(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
