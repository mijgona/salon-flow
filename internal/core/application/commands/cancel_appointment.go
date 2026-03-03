package commands

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/ports"

	"github.com/google/uuid"
)

// CancelAppointmentCommand holds data for cancelling an appointment.
type CancelAppointmentCommand struct {
	AppointmentID uuid.UUID
	Reason        string
}

// CancelAppointmentHandler handles appointment cancellation.
type CancelAppointmentHandler struct {
	appointmentRepo ports.AppointmentRepository
	scheduleRepo    ports.MasterScheduleRepository
	txManager       ports.TxManager
}

// NewCancelAppointmentHandler creates a new handler.
func NewCancelAppointmentHandler(
	appointmentRepo ports.AppointmentRepository,
	scheduleRepo ports.MasterScheduleRepository,
	txManager ports.TxManager,
) *CancelAppointmentHandler {
	return &CancelAppointmentHandler{
		appointmentRepo: appointmentRepo,
		scheduleRepo:    scheduleRepo,
		txManager:       txManager,
	}
}

// Handle executes the cancel appointment command.
func (h *CancelAppointmentHandler) Handle(ctx context.Context, cmd CancelAppointmentCommand) error {
	return h.txManager.Execute(ctx, func(tx interface{}) error {
		appointment, err := h.appointmentRepo.Get(ctx, tx, cmd.AppointmentID)
		if err != nil {
			return err
		}

		// Release the slot on the master schedule
		schedule, err := h.scheduleRepo.GetByMasterAndDate(ctx, tx, appointment.MasterID(), appointment.TimeSlot().StartTime())
		if err == nil && schedule != nil {
			schedule.ReleaseSlot(appointment.TimeSlot())
			_ = h.scheduleRepo.Update(ctx, tx, schedule)
		}

		if err := appointment.Cancel(cmd.Reason); err != nil {
			return err
		}

		return h.appointmentRepo.Update(ctx, tx, appointment)
	})
}
