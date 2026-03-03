package scheduling

import (
	"errors"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"
	"time"

	"github.com/google/uuid"
)

// MasterSchedule is the aggregate root for managing a master's daily schedule.
type MasterSchedule struct {
	*ddd.BaseAggregate[uuid.UUID]
	masterID     uuid.UUID
	salonID      uuid.UUID
	date         time.Time
	workingHours WorkingHours
	bookedSlots  []TimeSlot
	blockedSlots []TimeSlot
}

// NewMasterSchedule creates a new MasterSchedule aggregate.
func NewMasterSchedule(
	masterID, salonID uuid.UUID,
	date time.Time,
	workingHours WorkingHours,
) (*MasterSchedule, error) {
	if masterID == uuid.Nil {
		return nil, errors.New("master ID is required")
	}
	if salonID == uuid.Nil {
		return nil, errors.New("salon ID is required")
	}

	return &MasterSchedule{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](uuid.New()),
		masterID:      masterID,
		salonID:       salonID,
		date:          date,
		workingHours:  workingHours,
		bookedSlots:   make([]TimeSlot, 0),
		blockedSlots:  make([]TimeSlot, 0),
	}, nil
}

// RestoreMasterSchedule rehydrates a MasterSchedule from the database.
func RestoreMasterSchedule(
	id, masterID, salonID uuid.UUID,
	date time.Time,
	workingHours WorkingHours,
	bookedSlots, blockedSlots []TimeSlot,
) *MasterSchedule {
	return &MasterSchedule{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),
		masterID:      masterID,
		salonID:       salonID,
		date:          date,
		workingHours:  workingHours,
		bookedSlots:   bookedSlots,
		blockedSlots:  blockedSlots,
	}
}

// IsAvailable checks if a time slot is available for booking.
func (ms *MasterSchedule) IsAvailable(slot TimeSlot) bool {
	// Must be within working hours
	if !ms.workingHours.IsWithinWorkingHours(slot) {
		return false
	}

	// Must not overlap with any booked slots
	for _, booked := range ms.bookedSlots {
		if slot.OverlapsWith(booked) {
			return false
		}
	}

	// Must not overlap with any blocked slots
	for _, blocked := range ms.blockedSlots {
		if slot.OverlapsWith(blocked) {
			return false
		}
	}

	return true
}

// BookSlot reserves a time slot.
func (ms *MasterSchedule) BookSlot(slot TimeSlot) error {
	if !ms.IsAvailable(slot) {
		return errors.New("time slot is not available")
	}
	ms.bookedSlots = append(ms.bookedSlots, slot)
	return nil
}

// ReleaseSlot removes a booked time slot.
func (ms *MasterSchedule) ReleaseSlot(slot TimeSlot) {
	result := make([]TimeSlot, 0, len(ms.bookedSlots))
	for _, booked := range ms.bookedSlots {
		if !booked.Equal(slot) {
			result = append(result, booked)
		}
	}
	ms.bookedSlots = result
}

// BlockSlot blocks a time slot.
func (ms *MasterSchedule) BlockSlot(slot TimeSlot) error {
	for _, blocked := range ms.blockedSlots {
		if slot.OverlapsWith(blocked) {
			return errors.New("slot overlaps with an already blocked slot")
		}
	}
	ms.blockedSlots = append(ms.blockedSlots, slot)
	return nil
}

// GetAvailableSlots returns all available slots for a given duration within working hours.
func (ms *MasterSchedule) GetAvailableSlots(duration time.Duration) []TimeSlot {
	available := make([]TimeSlot, 0)
	step := 15 * time.Minute // 15-minute intervals

	current := ms.workingHours.StartTime()
	workEnd := ms.workingHours.EndTime()

	for current.Add(duration).Before(workEnd) || current.Add(duration).Equal(workEnd) {
		candidate, err := NewTimeSlot(current, current.Add(duration))
		if err == nil && ms.IsAvailable(candidate) {
			available = append(available, candidate)
		}
		current = current.Add(step)
	}

	return available
}

// Getters
func (ms *MasterSchedule) MasterID() uuid.UUID        { return ms.masterID }
func (ms *MasterSchedule) SalonID() uuid.UUID         { return ms.salonID }
func (ms *MasterSchedule) Date() time.Time            { return ms.date }
func (ms *MasterSchedule) WorkingHours() WorkingHours { return ms.workingHours }
func (ms *MasterSchedule) BookedSlots() []TimeSlot    { return ms.bookedSlots }
func (ms *MasterSchedule) BlockedSlots() []TimeSlot   { return ms.blockedSlots }
