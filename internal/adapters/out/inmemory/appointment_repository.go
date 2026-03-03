package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"

	"github.com/google/uuid"
)

// InMemoryAppointmentRepository is an in-memory implementation for testing.
type InMemoryAppointmentRepository struct {
	mu           sync.RWMutex
	appointments map[uuid.UUID]*scheduling.Appointment
}

// NewInMemoryAppointmentRepository creates a new in-memory appointment repository.
func NewInMemoryAppointmentRepository() *InMemoryAppointmentRepository {
	return &InMemoryAppointmentRepository{
		appointments: make(map[uuid.UUID]*scheduling.Appointment),
	}
}

func (r *InMemoryAppointmentRepository) Add(_ context.Context, _ interface{}, a *scheduling.Appointment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.appointments[a.ID()] = a
	return nil
}

func (r *InMemoryAppointmentRepository) Update(_ context.Context, _ interface{}, a *scheduling.Appointment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.appointments[a.ID()] = a
	return nil
}

func (r *InMemoryAppointmentRepository) Get(_ context.Context, _ interface{}, id uuid.UUID) (*scheduling.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.appointments[id]
	if !ok {
		return nil, nil
	}
	return a, nil
}

func (r *InMemoryAppointmentRepository) FindByClientID(_ context.Context, _ interface{}, clientID uuid.UUID) ([]*scheduling.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*scheduling.Appointment
	for _, a := range r.appointments {
		if a.ClientID() == clientID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *InMemoryAppointmentRepository) FindByMasterAndDate(_ context.Context, _ interface{}, masterID uuid.UUID, date time.Time) ([]*scheduling.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*scheduling.Appointment
	for _, a := range r.appointments {
		if a.MasterID() == masterID {
			y1, m1, d1 := a.TimeSlot().StartTime().Date()
			y2, m2, d2 := date.Date()
			if y1 == y2 && m1 == m2 && d1 == d2 {
				result = append(result, a)
			}
		}
	}
	return result, nil
}

func (r *InMemoryAppointmentRepository) FindByDateRange(_ context.Context, _ interface{}, tenantID uuid.UUID, from, to time.Time) ([]*scheduling.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*scheduling.Appointment
	for _, a := range r.appointments {
		if a.TenantID().UUID() == tenantID && !a.TimeSlot().StartTime().Before(from) && a.TimeSlot().StartTime().Before(to) {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *InMemoryAppointmentRepository) FindByMasterDateRange(_ context.Context, _ interface{}, masterID uuid.UUID, from, to time.Time) ([]*scheduling.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*scheduling.Appointment
	for _, a := range r.appointments {
		if a.MasterID() == masterID && !a.TimeSlot().StartTime().Before(from) && a.TimeSlot().StartTime().Before(to) {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *InMemoryAppointmentRepository) FindBySalonDateRange(_ context.Context, _ interface{}, salonID uuid.UUID, from, to time.Time) ([]*scheduling.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*scheduling.Appointment
	for _, a := range r.appointments {
		if a.SalonID() == salonID && !a.TimeSlot().StartTime().Before(from) && a.TimeSlot().StartTime().Before(to) {
			result = append(result, a)
		}
	}
	return result, nil
}
