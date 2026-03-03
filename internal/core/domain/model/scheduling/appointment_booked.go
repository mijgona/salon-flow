package scheduling

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AppointmentBooked is a domain event raised when an appointment is booked.
type AppointmentBooked struct {
	eventID       uuid.UUID
	appointmentID uuid.UUID
	clientID      uuid.UUID
	masterID      uuid.UUID
	salonID       uuid.UUID
	serviceID     uuid.UUID
	serviceName   string
	startTime     string
	endTime       string
	price         decimal.Decimal
	source        BookingSource
}

// NewAppointmentBooked creates a new AppointmentBooked event.
func NewAppointmentBooked(
	appointmentID, clientID, masterID, salonID, serviceID uuid.UUID,
	serviceName, startTime, endTime string,
	price decimal.Decimal,
	source BookingSource,
) AppointmentBooked {
	return AppointmentBooked{
		eventID:       uuid.New(),
		appointmentID: appointmentID,
		clientID:      clientID,
		masterID:      masterID,
		salonID:       salonID,
		serviceID:     serviceID,
		serviceName:   serviceName,
		startTime:     startTime,
		endTime:       endTime,
		price:         price,
		source:        source,
	}
}

func (e AppointmentBooked) GetID() uuid.UUID         { return e.eventID }
func (e AppointmentBooked) GetName() string          { return "appointment.booked" }
func (e AppointmentBooked) AppointmentID() uuid.UUID { return e.appointmentID }
func (e AppointmentBooked) ClientID() uuid.UUID      { return e.clientID }
func (e AppointmentBooked) MasterID() uuid.UUID      { return e.masterID }
func (e AppointmentBooked) SalonID() uuid.UUID       { return e.salonID }
func (e AppointmentBooked) ServiceID() uuid.UUID     { return e.serviceID }
func (e AppointmentBooked) ServiceName() string      { return e.serviceName }
func (e AppointmentBooked) StartTime() string        { return e.startTime }
func (e AppointmentBooked) EndTime() string          { return e.endTime }
func (e AppointmentBooked) Price() decimal.Decimal   { return e.price }
func (e AppointmentBooked) Source() BookingSource    { return e.source }
