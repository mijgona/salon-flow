package commands

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"
	"github.com/mijgona/salon-crm/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

// BookAppointmentCommand holds data for booking an appointment.
type BookAppointmentCommand struct {
	TenantID  uuid.UUID
	ClientID  uuid.UUID
	MasterID  uuid.UUID
	SalonID   uuid.UUID
	ServiceID uuid.UUID
	StartTime time.Time
	Comment   string
	Source    string
}

// BookAppointmentHandler handles appointment booking.
type BookAppointmentHandler struct {
	appointmentRepo ports.AppointmentRepository
	scheduleRepo    ports.MasterScheduleRepository
	serviceCatalog  ports.ServiceCatalogClient
	txManager       ports.TxManager
}

// NewBookAppointmentHandler creates a new handler.
func NewBookAppointmentHandler(
	appointmentRepo ports.AppointmentRepository,
	scheduleRepo ports.MasterScheduleRepository,
	serviceCatalog ports.ServiceCatalogClient,
	txManager ports.TxManager,
) *BookAppointmentHandler {
	return &BookAppointmentHandler{
		appointmentRepo: appointmentRepo,
		scheduleRepo:    scheduleRepo,
		serviceCatalog:  serviceCatalog,
		txManager:       txManager,
	}
}

// Handle executes the book appointment command.
func (h *BookAppointmentHandler) Handle(ctx context.Context, cmd BookAppointmentCommand) (uuid.UUID, error) {
	tenantID, err := model.NewTenantID(cmd.TenantID)
	if err != nil {
		return uuid.Nil, err
	}

	// Get service details
	service, err := h.serviceCatalog.GetService(ctx, cmd.ServiceID)
	if err != nil {
		return uuid.Nil, err
	}

	serviceInfo, err := scheduling.NewServiceInfo(
		service.ID, service.Name, service.Duration, service.Price,
	)
	if err != nil {
		return uuid.Nil, err
	}

	// Build time slot
	endTime := cmd.StartTime.Add(service.Duration)
	timeSlot, err := scheduling.NewTimeSlot(cmd.StartTime, endTime)
	if err != nil {
		return uuid.Nil, err
	}

	source := scheduling.BookingSource(cmd.Source)

	var appointmentID uuid.UUID
	err = h.txManager.Execute(ctx, func(tx interface{}) error {
		// Load master schedule
		schedule, err := h.scheduleRepo.GetByMasterAndDate(ctx, tx, cmd.MasterID, cmd.StartTime)
		if err != nil {
			return err
		}
		if schedule == nil {
			return &NotFoundError{Entity: "master schedule"}
		}

		// Check availability
		if !schedule.IsAvailable(timeSlot) {
			return &ConflictError{Message: "time slot is not available"}
		}

		// Book the slot
		if err := schedule.BookSlot(timeSlot); err != nil {
			return err
		}

		// Create appointment
		appointment, err := scheduling.NewAppointment(
			tenantID, cmd.ClientID, cmd.MasterID, cmd.SalonID,
			serviceInfo, timeSlot, service.Price, source, cmd.Comment,
		)
		if err != nil {
			return err
		}

		// Persist
		if err := h.scheduleRepo.Update(ctx, tx, schedule); err != nil {
			return err
		}
		if err := h.appointmentRepo.Add(ctx, tx, appointment); err != nil {
			return err
		}

		appointmentID = appointment.ID()
		return nil
	})

	return appointmentID, err
}

// NotFoundError indicates an entity was not found.
type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return e.Entity + " not found"
}

// ConflictError indicates a business rule conflict.
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}
