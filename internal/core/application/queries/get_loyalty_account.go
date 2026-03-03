package queries

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model/loyalty"
	"github.com/mijgona/salon-crm/internal/core/ports"

	"github.com/google/uuid"
)

// GetLoyaltyAccountQuery holds data for retrieving a loyalty account.
type GetLoyaltyAccountQuery struct {
	ClientID uuid.UUID
}

// GetLoyaltyAccountResult is the result of the GetLoyaltyAccount query.
type GetLoyaltyAccountResult struct {
	Account *loyalty.LoyaltyAccount
}

// GetLoyaltyAccountHandler handles the get loyalty account query.
type GetLoyaltyAccountHandler struct {
	loyaltyRepo ports.LoyaltyRepository
}

// NewGetLoyaltyAccountHandler creates a new handler.
func NewGetLoyaltyAccountHandler(loyaltyRepo ports.LoyaltyRepository) *GetLoyaltyAccountHandler {
	return &GetLoyaltyAccountHandler{loyaltyRepo: loyaltyRepo}
}

// Handle executes the get loyalty account query.
func (h *GetLoyaltyAccountHandler) Handle(ctx context.Context, q GetLoyaltyAccountQuery) (*GetLoyaltyAccountResult, error) {
	account, err := h.loyaltyRepo.GetByClientID(ctx, nil, q.ClientID)
	if err != nil {
		return nil, err
	}
	return &GetLoyaltyAccountResult{Account: account}, nil
}
