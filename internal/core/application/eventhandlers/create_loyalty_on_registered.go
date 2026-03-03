package eventhandlers

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/client"
	"github.com/mijgona/salon-crm/internal/core/domain/model/loyalty"
	"github.com/mijgona/salon-crm/internal/core/ports"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"

	"github.com/google/uuid"
)

// CreateLoyaltyOnRegisteredHandler creates a loyalty account when a client is registered.
type CreateLoyaltyOnRegisteredHandler struct {
	loyaltyRepo ports.LoyaltyRepository
	txManager   ports.TxManager
}

// NewCreateLoyaltyOnRegisteredHandler creates a new handler.
func NewCreateLoyaltyOnRegisteredHandler(
	loyaltyRepo ports.LoyaltyRepository,
	txManager ports.TxManager,
) *CreateLoyaltyOnRegisteredHandler {
	return &CreateLoyaltyOnRegisteredHandler{
		loyaltyRepo: loyaltyRepo,
		txManager:   txManager,
	}
}

// Handle processes the ClientRegistered event.
func (h *CreateLoyaltyOnRegisteredHandler) Handle(ctx context.Context, event ddd.DomainEvent) error {
	registered, ok := event.(client.ClientRegistered)
	if !ok {
		return nil
	}

	return h.txManager.Execute(ctx, func(tx interface{}) error {
		tenantID := model.MustNewTenantID(registered.TenantID())

		account, err := loyalty.NewLoyaltyAccount(registered.ClientID(), tenantID)
		if err != nil {
			return err
		}

		// If client was referred, add referral bonus
		if registered.ReferredByClientID() != uuid.Nil {
			referrerAccount, err := h.loyaltyRepo.GetByClientID(ctx, tx, registered.ReferredByClientID())
			if err == nil && referrerAccount != nil {
				_ = referrerAccount.AddReferral(registered.ClientID())
				referrerAccount.EarnPoints(500, "referral_bonus", registered.ClientID())
				_ = h.loyaltyRepo.Update(ctx, tx, referrerAccount)
			}
		}

		return h.loyaltyRepo.Add(ctx, tx, account)
	})
}
