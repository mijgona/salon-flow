package scheduling

import (
	"github.com/mijgona/salon-crm/internal/pkg/errs"
	"time"
)

// TimeSlot represents a time window with start and end times.
type TimeSlot struct {
	startTime time.Time
	endTime   time.Time
}

// NewTimeSlot creates a TimeSlot value object with validation.
func NewTimeSlot(startTime, endTime time.Time) (TimeSlot, error) {
	if startTime.IsZero() {
		return TimeSlot{}, errs.NewErrValueRequired("start time")
	}
	if endTime.IsZero() {
		return TimeSlot{}, errs.NewErrValueRequired("end time")
	}
	if !endTime.After(startTime) {
		return TimeSlot{}, errs.NewErrValueMustBe("end time", "after start time")
	}
	return TimeSlot{startTime: startTime, endTime: endTime}, nil
}

// MustNewTimeSlot creates a TimeSlot or panics.
func MustNewTimeSlot(startTime, endTime time.Time) TimeSlot {
	ts, err := NewTimeSlot(startTime, endTime)
	if err != nil {
		panic(err)
	}
	return ts
}

// StartTime returns the slot start time.
func (ts TimeSlot) StartTime() time.Time { return ts.startTime }

// EndTime returns the slot end time.
func (ts TimeSlot) EndTime() time.Time { return ts.endTime }

// Duration returns the duration of the time slot.
func (ts TimeSlot) Duration() time.Duration { return ts.endTime.Sub(ts.startTime) }

// OverlapsWith checks if this time slot overlaps with another.
func (ts TimeSlot) OverlapsWith(other TimeSlot) bool {
	return ts.startTime.Before(other.endTime) && other.startTime.Before(ts.endTime)
}

// Contains checks if a point in time falls within the slot.
func (ts TimeSlot) Contains(t time.Time) bool {
	return !t.Before(ts.startTime) && t.Before(ts.endTime)
}

// Equal checks value equality.
func (ts TimeSlot) Equal(other TimeSlot) bool {
	return ts.startTime.Equal(other.startTime) && ts.endTime.Equal(other.endTime)
}
