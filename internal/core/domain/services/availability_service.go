package services

import (
	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"
	"time"

	"github.com/google/uuid"
)

// AvailabilityService is a domain service for checking master availability.
type AvailabilityService struct {
	scheduleRepo ScheduleProvider
}

// ScheduleProvider provides master schedule data.
type ScheduleProvider interface {
	GetByMasterAndDate(masterID uuid.UUID, date time.Time) (*scheduling.MasterSchedule, error)
}

// NewAvailabilityService creates a new AvailabilityService.
func NewAvailabilityService(scheduleRepo ScheduleProvider) *AvailabilityService {
	return &AvailabilityService{scheduleRepo: scheduleRepo}
}

// GetAvailableSlots returns available time slots for a master on a given date.
func (as *AvailabilityService) GetAvailableSlots(
	masterID, salonID uuid.UUID,
	date time.Time,
	serviceDuration time.Duration,
) ([]scheduling.TimeSlot, error) {
	schedule, err := as.scheduleRepo.GetByMasterAndDate(masterID, date)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return []scheduling.TimeSlot{}, nil
	}
	return schedule.GetAvailableSlots(serviceDuration), nil
}

// IsSlotAvailable checks if a specific time slot is available for a master.
func (as *AvailabilityService) IsSlotAvailable(
	masterID uuid.UUID,
	date time.Time,
	slot scheduling.TimeSlot,
) (bool, error) {
	schedule, err := as.scheduleRepo.GetByMasterAndDate(masterID, date)
	if err != nil {
		return false, err
	}
	if schedule == nil {
		return false, nil
	}
	return schedule.IsAvailable(slot), nil
}
