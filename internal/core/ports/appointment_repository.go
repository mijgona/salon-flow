package ports

import (
	"context"
	"time"

	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"

	"github.com/google/uuid"
)

// AppointmentRepository defines operations for persisting Appointment aggregates.
type AppointmentRepository interface {
	Add(ctx context.Context, tx interface{}, a *scheduling.Appointment) error
	Update(ctx context.Context, tx interface{}, a *scheduling.Appointment) error
	Get(ctx context.Context, tx interface{}, id uuid.UUID) (*scheduling.Appointment, error)
	FindByClientID(ctx context.Context, tx interface{}, clientID uuid.UUID) ([]*scheduling.Appointment, error)
	FindByMasterAndDate(ctx context.Context, tx interface{}, masterID uuid.UUID, date time.Time) ([]*scheduling.Appointment, error)
	FindByDateRange(ctx context.Context, tx interface{}, tenantID uuid.UUID, from, to time.Time) ([]*scheduling.Appointment, error)
	FindByMasterDateRange(ctx context.Context, tx interface{}, masterID uuid.UUID, from, to time.Time) ([]*scheduling.Appointment, error)
	FindBySalonDateRange(ctx context.Context, tx interface{}, salonID uuid.UUID, from, to time.Time) ([]*scheduling.Appointment, error)
}
