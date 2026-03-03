package scheduling

import (
	"github.com/mijgona/salon-crm/internal/pkg/errs"
	"time"
)

// WorkingHours represents a master's working hours for a day.
type WorkingHours struct {
	startTime  time.Time
	endTime    time.Time
	breakStart time.Time
	breakEnd   time.Time
}

// NewWorkingHours creates a WorkingHours value object with validation.
func NewWorkingHours(startTime, endTime, breakStart, breakEnd time.Time) (WorkingHours, error) {
	if startTime.IsZero() || endTime.IsZero() {
		return WorkingHours{}, errs.NewErrValueRequired("working hours start/end")
	}
	if !endTime.After(startTime) {
		return WorkingHours{}, errs.NewErrValueMustBe("end time", "after start time")
	}
	if !breakStart.IsZero() && !breakEnd.IsZero() {
		if !breakEnd.After(breakStart) {
			return WorkingHours{}, errs.NewErrValueMustBe("break end", "after break start")
		}
	}
	return WorkingHours{
		startTime:  startTime,
		endTime:    endTime,
		breakStart: breakStart,
		breakEnd:   breakEnd,
	}, nil
}

// MustNewWorkingHours creates a WorkingHours or panics.
func MustNewWorkingHours(startTime, endTime, breakStart, breakEnd time.Time) WorkingHours {
	wh, err := NewWorkingHours(startTime, endTime, breakStart, breakEnd)
	if err != nil {
		panic(err)
	}
	return wh
}

// StartTime returns the work start time.
func (wh WorkingHours) StartTime() time.Time { return wh.startTime }

// EndTime returns the work end time.
func (wh WorkingHours) EndTime() time.Time { return wh.endTime }

// BreakStart returns the break start time.
func (wh WorkingHours) BreakStart() time.Time { return wh.breakStart }

// BreakEnd returns the break end time.
func (wh WorkingHours) BreakEnd() time.Time { return wh.breakEnd }

// HasBreak returns true if a break is defined.
func (wh WorkingHours) HasBreak() bool {
	return !wh.breakStart.IsZero() && !wh.breakEnd.IsZero()
}

// IsWithinWorkingHours checks if a time slot is within working hours and not during break.
func (wh WorkingHours) IsWithinWorkingHours(slot TimeSlot) bool {
	// Check if the slot is within the general working hours
	workStart := wh.startTime
	workEnd := wh.endTime

	if slot.StartTime().Before(workStart) || slot.EndTime().After(workEnd) {
		return false
	}

	// Check if the slot overlaps with the break
	if wh.HasBreak() {
		breakSlot, err := NewTimeSlot(wh.breakStart, wh.breakEnd)
		if err == nil && slot.OverlapsWith(breakSlot) {
			return false
		}
	}

	return true
}
