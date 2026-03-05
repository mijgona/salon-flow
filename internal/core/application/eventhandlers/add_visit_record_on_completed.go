package eventhandlers

import (
	"context"
	"time"

	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/client"
	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"
	"github.com/mijgona/salon-crm/internal/core/ports"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"

	"github.com/shopspring/decimal"
)

// AddVisitRecordOnCompletedHandler adds a visit record to the client when an appointment is completed.
type AddVisitRecordOnCompletedHandler struct {
	clientRepo ports.ClientRepository
	txManager  ports.TxManager
}

// NewAddVisitRecordOnCompletedHandler creates a new handler.
func NewAddVisitRecordOnCompletedHandler(
	clientRepo ports.ClientRepository,
	txManager ports.TxManager,
) *AddVisitRecordOnCompletedHandler {
	return &AddVisitRecordOnCompletedHandler{
		clientRepo: clientRepo,
		txManager:  txManager,
	}
}

// Handle processes the AppointmentCompleted event.
func (h *AddVisitRecordOnCompletedHandler) Handle(ctx context.Context, event ddd.DomainEvent) error {
	completed, ok := event.(scheduling.AppointmentCompleted)
	if !ok {
		return nil
	}

	return h.txManager.Execute(ctx, func(tx interface{}) error {
		c, err := h.clientRepo.Get(ctx, tx, completed.ClientID())
		if err != nil {
			return err
		}

		price, _ := model.NewMoney(decimal.NewFromFloat(completed.FinalPrice().InexactFloat64()), "RUB")
		discount, _ := model.NewDiscount(completed.Discount())

		record := client.NewVisitRecord(
			completed.AppointmentID(),
			completed.MasterID(),
			completed.ServiceName(),
			price,
			discount,
			client.PaymentStatusPaid,
			time.Now(),
		)

		c.AddVisitRecord(record)
		return h.clientRepo.Update(ctx, tx, c)
	})
}
