package scheduling

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AppointmentCompleted is a domain event raised when an appointment is completed.
type AppointmentCompleted struct {
	eventID       uuid.UUID
	appointmentID uuid.UUID
	clientID      uuid.UUID
	masterID      uuid.UUID
	salonID       uuid.UUID
	serviceName   string
	finalPrice    decimal.Decimal
	discount      int
	paymentMethod string
}

// NewAppointmentCompleted creates a new AppointmentCompleted event.
func NewAppointmentCompleted(
	appointmentID, clientID, masterID, salonID uuid.UUID,
	serviceName string,
	finalPrice decimal.Decimal,
	discount int,
	paymentMethod string,
) AppointmentCompleted {
	return AppointmentCompleted{
		eventID:       uuid.New(),
		appointmentID: appointmentID,
		clientID:      clientID,
		masterID:      masterID,
		salonID:       salonID,
		serviceName:   serviceName,
		finalPrice:    finalPrice,
		discount:      discount,
		paymentMethod: paymentMethod,
	}
}

func (e AppointmentCompleted) GetID() uuid.UUID            { return e.eventID }
func (e AppointmentCompleted) GetName() string             { return "appointment.completed" }
func (e AppointmentCompleted) AppointmentID() uuid.UUID    { return e.appointmentID }
func (e AppointmentCompleted) ClientID() uuid.UUID         { return e.clientID }
func (e AppointmentCompleted) MasterID() uuid.UUID         { return e.masterID }
func (e AppointmentCompleted) SalonID() uuid.UUID          { return e.salonID }
func (e AppointmentCompleted) ServiceName() string         { return e.serviceName }
func (e AppointmentCompleted) FinalPrice() decimal.Decimal { return e.finalPrice }
func (e AppointmentCompleted) Discount() int               { return e.discount }
func (e AppointmentCompleted) PaymentMethod() string       { return e.paymentMethod }
