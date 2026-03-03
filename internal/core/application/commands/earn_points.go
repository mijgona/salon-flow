package commands

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/ports"

	"github.com/google/uuid"
)

// EarnPointsCommand holds data for earning loyalty points.
type EarnPointsCommand struct {
	ClientID        uuid.UUID
	Amount          int
	Reason          string
	RelatedEntityID uuid.UUID
}

// EarnPointsHandler handles the earn points command.
type EarnPointsHandler struct {
	loyaltyRepo ports.LoyaltyRepository
	txManager   ports.TxManager
}

// NewEarnPointsHandler creates a new handler.
func NewEarnPointsHandler(loyaltyRepo ports.LoyaltyRepository, txManager ports.TxManager) *EarnPointsHandler {
	return &EarnPointsHandler{loyaltyRepo: loyaltyRepo, txManager: txManager}
}

// Handle executes the earn points command.
func (h *EarnPointsHandler) Handle(ctx context.Context, cmd EarnPointsCommand) error {
	return h.txManager.Execute(ctx, func(tx interface{}) error {
		account, err := h.loyaltyRepo.GetByClientID(ctx, tx, cmd.ClientID)
		if err != nil {
			return err
		}

		account.EarnPoints(cmd.Amount, cmd.Reason, cmd.RelatedEntityID)
		account.RecalculateTier()

		return h.loyaltyRepo.Update(ctx, tx, account)
	})
}
