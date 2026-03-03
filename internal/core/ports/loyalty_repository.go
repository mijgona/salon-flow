package ports

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model/loyalty"

	"github.com/google/uuid"
)

// LoyaltyRepository defines operations for persisting LoyaltyAccount aggregates.
type LoyaltyRepository interface {
	Add(ctx context.Context, tx interface{}, la *loyalty.LoyaltyAccount) error
	Update(ctx context.Context, tx interface{}, la *loyalty.LoyaltyAccount) error
	GetByClientID(ctx context.Context, tx interface{}, clientID uuid.UUID) (*loyalty.LoyaltyAccount, error)
}
