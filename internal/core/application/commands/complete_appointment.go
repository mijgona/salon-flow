package commands

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/ports"

	"github.com/google/uuid"
)

// CompleteAppointmentCommand holds data for completing an appointment.
type CompleteAppointmentCommand struct {
	AppointmentID uuid.UUID
}

// CompleteAppointmentHandler handles appointment completion.
type CompleteAppointmentHandler struct {
	appointmentRepo ports.AppointmentRepository
	txManager       ports.TxManager
}

// NewCompleteAppointmentHandler creates a new handler.
func NewCompleteAppointmentHandler(appointmentRepo ports.AppointmentRepository, txManager ports.TxManager) *CompleteAppointmentHandler {
	return &CompleteAppointmentHandler{appointmentRepo: appointmentRepo, txManager: txManager}
}

// Handle executes the complete appointment command.
func (h *CompleteAppointmentHandler) Handle(ctx context.Context, cmd CompleteAppointmentCommand) error {
	return h.txManager.Execute(ctx, func(tx interface{}) error {
		appointment, err := h.appointmentRepo.Get(ctx, tx, cmd.AppointmentID)
		if err != nil {
			return err
		}

		if err := appointment.Complete(); err != nil {
			return err
		}

		return h.appointmentRepo.Update(ctx, tx, appointment)
	})
}
