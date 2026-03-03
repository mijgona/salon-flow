package eventhandlers

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"
	"github.com/mijgona/salon-crm/internal/core/ports"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"
	"math"
)

// AccruePointsOnCompletedHandler accrues loyalty points when an appointment is completed.
type AccruePointsOnCompletedHandler struct {
	loyaltyRepo ports.LoyaltyRepository
	txManager   ports.TxManager
}

// NewAccruePointsOnCompletedHandler creates a new handler.
func NewAccruePointsOnCompletedHandler(
	loyaltyRepo ports.LoyaltyRepository,
	txManager ports.TxManager,
) *AccruePointsOnCompletedHandler {
	return &AccruePointsOnCompletedHandler{
		loyaltyRepo: loyaltyRepo,
		txManager:   txManager,
	}
}

// Handle processes the AppointmentCompleted event.
func (h *AccruePointsOnCompletedHandler) Handle(ctx context.Context, event ddd.DomainEvent) error {
	completed, ok := event.(scheduling.AppointmentCompleted)
	if !ok {
		return nil
	}

	return h.txManager.Execute(ctx, func(tx interface{}) error {
		account, err := h.loyaltyRepo.GetByClientID(ctx, tx, completed.ClientID())
		if err != nil {
			return err
		}

		// 1 point per 10 RUB × tier multiplier
		basePoints := completed.FinalPrice().IntPart() / 10
		multiplied := int(math.Round(float64(basePoints) * account.Tier().PointsMultiplier()))

		account.EarnPoints(multiplied, "appointment_completed", completed.AppointmentID())
		account.RecalculateTier()

		return h.loyaltyRepo.Update(ctx, tx, account)
	})
}
