package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Querier abstracts pgx.Tx and *pgxpool.Pool for query execution.
type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// TxManager provides pgx transactional boundaries.
type TxManager struct {
	pool *pgxpool.Pool
}

// NewTxManager creates a new TxManager backed by a pgx pool.
func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// Pool returns the underlying pool (for non-transactional reads).
func (tm *TxManager) Pool() *pgxpool.Pool {
	return tm.pool
}

// Execute runs fn inside a pgx transaction.
// The tx interface{} passed to fn is a pgx.Tx.
func (tm *TxManager) Execute(ctx context.Context, fn func(tx interface{}) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

// ExtractTx extracts a pgx.Tx from the interface{} passed through ports.
// If tx is nil, returns nil (caller should use the pool directly).
func ExtractTx(tx interface{}) pgx.Tx {
	if tx == nil {
		return nil
	}
	if pgxTx, ok := tx.(pgx.Tx); ok {
		return pgxTx
	}
	return nil
}
