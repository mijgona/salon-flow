package scheduling

import (
	"time"

	"github.com/google/uuid"
)

// AppointmentCancelled is a domain event raised when an appointment is cancelled.
type AppointmentCancelled struct {
	eventID           uuid.UUID
	appointmentID     uuid.UUID
	clientID          uuid.UUID
	masterID          uuid.UUID
	salonID           uuid.UUID
	originalStartTime time.Time
	cancelledAt       time.Time
	reason            string
}

// NewAppointmentCancelled creates a new AppointmentCancelled event.
func NewAppointmentCancelled(
	appointmentID, clientID, masterID, salonID uuid.UUID,
	originalStartTime time.Time,
	reason string,
) AppointmentCancelled {
	return AppointmentCancelled{
		eventID:           uuid.New(),
		appointmentID:     appointmentID,
		clientID:          clientID,
		masterID:          masterID,
		salonID:           salonID,
		originalStartTime: originalStartTime,
		cancelledAt:       time.Now(),
		reason:            reason,
	}
}

func (e AppointmentCancelled) GetID() uuid.UUID             { return e.eventID }
func (e AppointmentCancelled) GetName() string              { return "appointment.cancelled_by_client" }
func (e AppointmentCancelled) AppointmentID() uuid.UUID     { return e.appointmentID }
func (e AppointmentCancelled) ClientID() uuid.UUID          { return e.clientID }
func (e AppointmentCancelled) MasterID() uuid.UUID          { return e.masterID }
func (e AppointmentCancelled) SalonID() uuid.UUID           { return e.salonID }
func (e AppointmentCancelled) OriginalStartTime() time.Time { return e.originalStartTime }
func (e AppointmentCancelled) CancelledAt() time.Time       { return e.cancelledAt }
func (e AppointmentCancelled) Reason() string               { return e.reason }
