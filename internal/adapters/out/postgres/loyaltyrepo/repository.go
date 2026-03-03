package loyaltyrepo

import (
	"context"
	"fmt"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/loyalty"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// PostgresLoyaltyRepository is a pgx + Squirrel implementation of LoyaltyRepository.
type PostgresLoyaltyRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresLoyaltyRepository creates a new repository.
func NewPostgresLoyaltyRepository(pool *pgxpool.Pool) *PostgresLoyaltyRepository {
	return &PostgresLoyaltyRepository{pool: pool}
}

func (r *PostgresLoyaltyRepository) Add(ctx context.Context, tx interface{}, la *loyalty.LoyaltyAccount) error {
	query, args, err := psql.Insert("loyalty_accounts").
		Columns("id", "client_id", "tenant_id", "tier", "balance", "lifetime_points").
		Values(
			la.ID(), la.ClientID(), la.TenantID().UUID(),
			la.Tier().String(), la.Balance().Value(), la.LifetimePoints().Value(),
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	if err := r.exec(ctx, tx, query, args...); err != nil {
		return err
	}

	// Persist transactions
	for _, t := range la.Transactions() {
		if err := r.insertTransaction(ctx, tx, la.ID(), t); err != nil {
			return err
		}
	}

	// Persist referrals
	for _, ref := range la.Referrals() {
		if err := r.insertReferral(ctx, tx, la.ID(), ref); err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresLoyaltyRepository) Update(ctx context.Context, tx interface{}, la *loyalty.LoyaltyAccount) error {
	query, args, err := psql.Update("loyalty_accounts").
		Set("tier", la.Tier().String()).
		Set("balance", la.Balance().Value()).
		Set("lifetime_points", la.LifetimePoints().Value()).
		Where(sq.Eq{"id": la.ID()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	if err := r.exec(ctx, tx, query, args...); err != nil {
		return err
	}

	// Upsert new transactions (insert only — transactions are append-only)
	for _, t := range la.Transactions() {
		// Use INSERT ... ON CONFLICT DO NOTHING to avoid duplicates
		tQuery, tArgs, err := psql.Insert("points_transactions").
			Columns("id", "loyalty_account_id", "amount", "type", "reason", "related_entity_id", "created_at").
			Values(t.ID(), la.ID(), t.Amount(), string(t.TransactionType()), t.Reason(), t.RelatedEntityID(), t.CreatedAt()).
			Suffix("ON CONFLICT (id) DO NOTHING").
			ToSql()
		if err != nil {
			return fmt.Errorf("build transaction upsert: %w", err)
		}
		if err := r.exec(ctx, tx, tQuery, tArgs...); err != nil {
			return err
		}
	}

	// Upsert referrals
	for _, ref := range la.Referrals() {
		rQuery, rArgs, err := psql.Insert("referrals").
			Columns("id", "loyalty_account_id", "referred_client_id", "status", "bonus_earned", "created_at").
			Values(ref.ID(), la.ID(), ref.ReferredClientID(), string(ref.Status()), ref.BonusEarned(), ref.CreatedAt()).
			Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, bonus_earned = EXCLUDED.bonus_earned").
			ToSql()
		if err != nil {
			return fmt.Errorf("build referral upsert: %w", err)
		}
		if err := r.exec(ctx, tx, rQuery, rArgs...); err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresLoyaltyRepository) GetByClientID(ctx context.Context, tx interface{}, clientID uuid.UUID) (*loyalty.LoyaltyAccount, error) {
	// 1. Load account
	query, args, err := psql.Select("id", "client_id", "tenant_id", "tier", "balance", "lifetime_points").
		From("loyalty_accounts").
		Where(sq.Eq{"client_id": clientID}).
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

	var (
		id             uuid.UUID
		loadedClientID uuid.UUID
		tenantID       uuid.UUID
		tier           string
		balance        int
		lifetimePoints int
	)
	if err := rows.Scan(&id, &loadedClientID, &tenantID, &tier, &balance, &lifetimePoints); err != nil {
		return nil, fmt.Errorf("scan loyalty_account: %w", err)
	}
	rows.Close()

	// 2. Load transactions
	transactions, err := r.loadTransactions(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	// 3. Load referrals
	referrals, err := r.loadReferrals(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return loyalty.RestoreLoyaltyAccount(
		id, loadedClientID,
		model.MustNewTenantID(tenantID),
		loyalty.LoyaltyTier(tier),
		loyalty.MustNewPoints(balance),
		loyalty.MustNewPoints(lifetimePoints),
		transactions,
		referrals,
	), nil
}

// --- internal helpers ---

func (r *PostgresLoyaltyRepository) exec(ctx context.Context, tx interface{}, sql string, args ...any) error {
	if t := postgres.ExtractTx(tx); t != nil {
		_, err := t.Exec(ctx, sql, args...)
		return err
	}
	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *PostgresLoyaltyRepository) query(ctx context.Context, tx interface{}, sql string, args ...any) (pgx.Rows, error) {
	if t := postgres.ExtractTx(tx); t != nil {
		return t.Query(ctx, sql, args...)
	}
	return r.pool.Query(ctx, sql, args...)
}

func (r *PostgresLoyaltyRepository) insertTransaction(ctx context.Context, tx interface{}, accountID uuid.UUID, t loyalty.PointsTransaction) error {
	q, a, err := psql.Insert("points_transactions").
		Columns("id", "loyalty_account_id", "amount", "type", "reason", "related_entity_id", "created_at").
		Values(t.ID(), accountID, t.Amount(), string(t.TransactionType()), t.Reason(), t.RelatedEntityID(), t.CreatedAt()).
		ToSql()
	if err != nil {
		return err
	}
	return r.exec(ctx, tx, q, a...)
}

func (r *PostgresLoyaltyRepository) insertReferral(ctx context.Context, tx interface{}, accountID uuid.UUID, ref loyalty.Referral) error {
	q, a, err := psql.Insert("referrals").
		Columns("id", "loyalty_account_id", "referred_client_id", "status", "bonus_earned", "created_at").
		Values(ref.ID(), accountID, ref.ReferredClientID(), string(ref.Status()), ref.BonusEarned(), ref.CreatedAt()).
		ToSql()
	if err != nil {
		return err
	}
	return r.exec(ctx, tx, q, a...)
}

func (r *PostgresLoyaltyRepository) loadTransactions(ctx context.Context, tx interface{}, accountID uuid.UUID) ([]loyalty.PointsTransaction, error) {
	q, a, err := psql.Select("id", "amount", "type", "reason", "related_entity_id", "created_at").
		From("points_transactions").
		Where(sq.Eq{"loyalty_account_id": accountID}).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.query(ctx, tx, q, a...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []loyalty.PointsTransaction
	for rows.Next() {
		var (
			id              uuid.UUID
			amount          int
			txType          string
			reason          string
			relatedEntityID uuid.UUID
			createdAt       time.Time
		)
		if err := rows.Scan(&id, &amount, &txType, &reason, &relatedEntityID, &createdAt); err != nil {
			return nil, err
		}
		result = append(result, loyalty.RestorePointsTransaction(
			id, amount, loyalty.TransactionType(txType), reason, relatedEntityID, createdAt,
		))
	}
	return result, rows.Err()
}

func (r *PostgresLoyaltyRepository) loadReferrals(ctx context.Context, tx interface{}, accountID uuid.UUID) ([]loyalty.Referral, error) {
	q, a, err := psql.Select("id", "referred_client_id", "status", "bonus_earned", "created_at").
		From("referrals").
		Where(sq.Eq{"loyalty_account_id": accountID}).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.query(ctx, tx, q, a...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []loyalty.Referral
	for rows.Next() {
		var (
			id               uuid.UUID
			referredClientID uuid.UUID
			status           string
			bonusEarned      int
			createdAt        time.Time
		)
		if err := rows.Scan(&id, &referredClientID, &status, &bonusEarned, &createdAt); err != nil {
			return nil, err
		}
		result = append(result, loyalty.RestoreReferral(
			id, referredClientID, loyalty.ReferralStatus(status), bonusEarned, createdAt,
		))
	}
	return result, rows.Err()
}
