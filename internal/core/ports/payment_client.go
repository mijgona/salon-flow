package ports

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model"

	"github.com/google/uuid"
)

// PaymentClient provides payment operations via ACL.
type PaymentClient interface {
	ProcessPayment(ctx context.Context, amount model.Money, clientID uuid.UUID, description string) error
	RefundPayment(ctx context.Context, paymentID uuid.UUID) error
}
