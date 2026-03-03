package scheduling

import "errors"

// AppointmentStatus represents the lifecycle status of an appointment.
type AppointmentStatus string

const (
	StatusRequested         AppointmentStatus = "requested"
	StatusConfirmed         AppointmentStatus = "confirmed"
	StatusInProgress        AppointmentStatus = "in_progress"
	StatusCompleted         AppointmentStatus = "completed"
	StatusCancelledByClient AppointmentStatus = "cancelled_by_client"
	StatusCancelledBySalon  AppointmentStatus = "cancelled_by_salon"
	StatusNoShow            AppointmentStatus = "no_show"
)

// IsValid checks if the status is a known value.
func (s AppointmentStatus) IsValid() bool {
	switch s {
	case StatusRequested, StatusConfirmed, StatusInProgress,
		StatusCompleted, StatusCancelledByClient, StatusCancelledBySalon, StatusNoShow:
		return true
	}
	return false
}

// CanTransitionTo checks if a status transition is valid.
func (s AppointmentStatus) CanTransitionTo(target AppointmentStatus) error {
	valid := map[AppointmentStatus][]AppointmentStatus{
		StatusRequested:  {StatusConfirmed, StatusCancelledByClient, StatusCancelledBySalon},
		StatusConfirmed:  {StatusInProgress, StatusCancelledByClient, StatusCancelledBySalon, StatusNoShow},
		StatusInProgress: {StatusCompleted},
	}

	allowed, exists := valid[s]
	if !exists {
		return errors.New("appointment is in a terminal state")
	}

	for _, a := range allowed {
		if a == target {
			return nil
		}
	}
	return errors.New("invalid status transition from " + string(s) + " to " + string(target))
}

// String returns the string representation.
func (s AppointmentStatus) String() string { return string(s) }

// IsTerminal returns true if the appointment is in a final state.
func (s AppointmentStatus) IsTerminal() bool {
	switch s {
	case StatusCompleted, StatusCancelledByClient, StatusCancelledBySalon, StatusNoShow:
		return true
	}
	return false
}
