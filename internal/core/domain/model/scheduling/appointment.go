package scheduling

import (
	"errors"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"
	"github.com/mijgona/salon-crm/internal/pkg/errs"
	"time"

	"github.com/google/uuid"
)

// Appointment is the aggregate root for scheduling appointments.
type Appointment struct {
	*ddd.BaseAggregate[uuid.UUID]
	tenantID    model.TenantID
	clientID    uuid.UUID
	masterID    uuid.UUID
	salonID     uuid.UUID
	serviceInfo ServiceInfo
	timeSlot    TimeSlot
	status      AppointmentStatus
	price       model.Money
	source      BookingSource
	comment     string
	createdAt   time.Time
}

// NewAppointment creates a new Appointment aggregate with validation.
func NewAppointment(
	tenantID model.TenantID,
	clientID, masterID, salonID uuid.UUID,
	serviceInfo ServiceInfo,
	timeSlot TimeSlot,
	price model.Money,
	source BookingSource,
	comment string,
) (*Appointment, error) {
	if clientID == uuid.Nil {
		return nil, errs.NewErrValueRequired("client ID")
	}
	if masterID == uuid.Nil {
		return nil, errs.NewErrValueRequired("master ID")
	}
	if salonID == uuid.Nil {
		return nil, errs.NewErrValueRequired("salon ID")
	}
	if timeSlot.StartTime().Before(time.Now()) {
		return nil, errors.New("cannot book an appointment in the past")
	}

	id := uuid.New()
	a := &Appointment{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),
		tenantID:      tenantID,
		clientID:      clientID,
		masterID:      masterID,
		salonID:       salonID,
		serviceInfo:   serviceInfo,
		timeSlot:      timeSlot,
		status:        StatusRequested,
		price:         price,
		source:        source,
		comment:       comment,
		createdAt:     time.Now(),
	}

	a.RaiseDomainEvent(NewAppointmentBooked(
		id, clientID, masterID, salonID, serviceInfo.ServiceID(),
		serviceInfo.Name(),
		timeSlot.StartTime().Format(time.RFC3339),
		timeSlot.EndTime().Format(time.RFC3339),
		price.Amount(),
		source,
	))

	return a, nil
}

// RestoreAppointment rehydrates an Appointment from the database.
func RestoreAppointment(
	id uuid.UUID,
	tenantID model.TenantID,
	clientID, masterID, salonID uuid.UUID,
	serviceInfo ServiceInfo,
	timeSlot TimeSlot,
	status AppointmentStatus,
	price model.Money,
	source BookingSource,
	comment string,
	createdAt time.Time,
) *Appointment {
	return &Appointment{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),
		tenantID:      tenantID,
		clientID:      clientID,
		masterID:      masterID,
		salonID:       salonID,
		serviceInfo:   serviceInfo,
		timeSlot:      timeSlot,
		status:        status,
		price:         price,
		source:        source,
		comment:       comment,
		createdAt:     createdAt,
	}
}

// Confirm transitions to Confirmed status.
func (a *Appointment) Confirm() error {
	if err := a.status.CanTransitionTo(StatusConfirmed); err != nil {
		return err
	}
	a.status = StatusConfirmed
	return nil
}

// Cancel transitions to CancelledByClient status.
func (a *Appointment) Cancel(reason string) error {
	if err := a.status.CanTransitionTo(StatusCancelledByClient); err != nil {
		return err
	}
	a.status = StatusCancelledByClient

	a.RaiseDomainEvent(NewAppointmentCancelled(
		a.ID(), a.clientID, a.masterID, a.salonID,
		a.timeSlot.StartTime(), reason,
	))

	return nil
}

// CancelBySalon transitions to CancelledBySalon status.
func (a *Appointment) CancelBySalon(reason string) error {
	if err := a.status.CanTransitionTo(StatusCancelledBySalon); err != nil {
		return err
	}
	a.status = StatusCancelledBySalon
	return nil
}

// Reschedule changes the time slot (only if not in progress or completed).
func (a *Appointment) Reschedule(newSlot TimeSlot) error {
	if a.status.IsTerminal() {
		return errors.New("cannot reschedule a terminal appointment")
	}
	if a.status == StatusInProgress {
		return errors.New("cannot reschedule an in-progress appointment")
	}
	if newSlot.StartTime().Before(time.Now()) {
		return errors.New("cannot reschedule to a past time")
	}
	a.timeSlot = newSlot
	return nil
}

// Complete transitions to Completed status.
func (a *Appointment) Complete() error {
	// Allow direct requested → confirmed → in_progress → completed fast-track
	if a.status == StatusRequested {
		a.status = StatusConfirmed
	}
	if a.status == StatusConfirmed {
		a.status = StatusInProgress
	}
	if err := a.status.CanTransitionTo(StatusCompleted); err != nil {
		return err
	}
	a.status = StatusCompleted

	a.RaiseDomainEvent(NewAppointmentCompleted(
		a.ID(), a.clientID, a.masterID, a.salonID,
		a.serviceInfo.Name(),
		a.price.Amount(),
		0,  // discount
		"", // payment method
	))

	return nil
}

// NoShow transitions to NoShow status.
func (a *Appointment) NoShow() error {
	if err := a.status.CanTransitionTo(StatusNoShow); err != nil {
		return err
	}
	a.status = StatusNoShow
	return nil
}

// Getters
func (a *Appointment) TenantID() model.TenantID  { return a.tenantID }
func (a *Appointment) ClientID() uuid.UUID       { return a.clientID }
func (a *Appointment) MasterID() uuid.UUID       { return a.masterID }
func (a *Appointment) SalonID() uuid.UUID        { return a.salonID }
func (a *Appointment) ServiceInfo() ServiceInfo  { return a.serviceInfo }
func (a *Appointment) TimeSlot() TimeSlot        { return a.timeSlot }
func (a *Appointment) Status() AppointmentStatus { return a.status }
func (a *Appointment) Price() model.Money        { return a.price }
func (a *Appointment) Source() BookingSource     { return a.source }
func (a *Appointment) Comment() string           { return a.comment }
func (a *Appointment) CreatedAt() time.Time      { return a.createdAt }
