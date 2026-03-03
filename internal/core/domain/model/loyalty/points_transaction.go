package loyalty

import (
	"time"

	"github.com/google/uuid"
)

// TransactionType represents the type of points transaction.
type TransactionType string

const (
	TransactionTypeEarn   TransactionType = "earn"
	TransactionTypeRedeem TransactionType = "redeem"
	TransactionTypeBonus  TransactionType = "bonus"
)

// PointsTransaction is an entity representing a loyalty points transaction.
type PointsTransaction struct {
	id              uuid.UUID
	amount          int
	transactionType TransactionType
	reason          string
	relatedEntityID uuid.UUID
	createdAt       time.Time
}

// NewPointsTransaction creates a new PointsTransaction entity.
func NewPointsTransaction(
	amount int,
	transactionType TransactionType,
	reason string,
	relatedEntityID uuid.UUID,
) PointsTransaction {
	return PointsTransaction{
		id:              uuid.New(),
		amount:          amount,
		transactionType: transactionType,
		reason:          reason,
		relatedEntityID: relatedEntityID,
		createdAt:       time.Now(),
	}
}

// RestorePointsTransaction creates a PointsTransaction from persisted data.
func RestorePointsTransaction(
	id uuid.UUID,
	amount int,
	transactionType TransactionType,
	reason string,
	relatedEntityID uuid.UUID,
	createdAt time.Time,
) PointsTransaction {
	return PointsTransaction{
		id:              id,
		amount:          amount,
		transactionType: transactionType,
		reason:          reason,
		relatedEntityID: relatedEntityID,
		createdAt:       createdAt,
	}
}

// Getters
func (pt PointsTransaction) ID() uuid.UUID                    { return pt.id }
func (pt PointsTransaction) Amount() int                      { return pt.amount }
func (pt PointsTransaction) TransactionType() TransactionType { return pt.transactionType }
func (pt PointsTransaction) Reason() string                   { return pt.reason }
func (pt PointsTransaction) RelatedEntityID() uuid.UUID       { return pt.relatedEntityID }
func (pt PointsTransaction) CreatedAt() time.Time             { return pt.createdAt }
