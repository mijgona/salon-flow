package queries

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"
	"github.com/mijgona/salon-crm/internal/core/ports"
)

// GetCalendarQuery represents a calendar query with filters.
type GetCalendarQuery struct {
	TenantID uuid.UUID
	MasterID uuid.UUID // optional — filter by master
	SalonID  uuid.UUID // optional — filter by salon
	From     time.Time
	To       time.Time
}

// CalendarDay groups appointments by date.
type CalendarDay struct {
	Date         string                `json:"date"`
	Appointments []CalendarAppointment `json:"appointments"`
	TotalCount   int                   `json:"total_count"`
	StatusCounts map[string]int        `json:"status_counts"`
}

// CalendarAppointment is a lightweight appointment DTO for calendar views.
type CalendarAppointment struct {
	ID          uuid.UUID `json:"id"`
	ClientID    uuid.UUID `json:"client_id"`
	MasterID    uuid.UUID `json:"master_id"`
	SalonID     uuid.UUID `json:"salon_id"`
	ServiceName string    `json:"service_name"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Status      string    `json:"status"`
	PriceAmount string    `json:"price_amount"`
	Comment     string    `json:"comment,omitempty"`
}

// CalendarResult is the result of a calendar query.
type CalendarResult struct {
	Days []CalendarDay `json:"days"`
	From string        `json:"from"`
	To   string        `json:"to"`
}

// GetCalendarHandler handles calendar queries.
type GetCalendarHandler struct {
	appointmentRepo ports.AppointmentRepository
}

// NewGetCalendarHandler creates a new GetCalendarHandler.
func NewGetCalendarHandler(appointmentRepo ports.AppointmentRepository) *GetCalendarHandler {
	return &GetCalendarHandler{appointmentRepo: appointmentRepo}
}

// Handle executes the calendar query.
func (h *GetCalendarHandler) Handle(ctx context.Context, q GetCalendarQuery) (*CalendarResult, error) {
	var appointments []*scheduling.Appointment
	var err error

	switch {
	case q.MasterID != uuid.Nil:
		appointments, err = h.appointmentRepo.FindByMasterDateRange(ctx, nil, q.MasterID, q.From, q.To)
	case q.SalonID != uuid.Nil:
		appointments, err = h.appointmentRepo.FindBySalonDateRange(ctx, nil, q.SalonID, q.From, q.To)
	default:
		appointments, err = h.appointmentRepo.FindByDateRange(ctx, nil, q.TenantID, q.From, q.To)
	}
	if err != nil {
		return nil, fmt.Errorf("get calendar: %w", err)
	}

	// Group by date
	dayMap := make(map[string][]CalendarAppointment)
	statusMap := make(map[string]map[string]int)

	for _, a := range appointments {
		dateKey := a.TimeSlot().StartTime().Format("2006-01-02")

		dto := CalendarAppointment{
			ID:          a.ID(),
			ClientID:    a.ClientID(),
			MasterID:    a.MasterID(),
			SalonID:     a.SalonID(),
			ServiceName: a.ServiceInfo().Name(),
			StartTime:   a.TimeSlot().StartTime().Format(time.RFC3339),
			EndTime:     a.TimeSlot().EndTime().Format(time.RFC3339),
			Status:      a.Status().String(),
			PriceAmount: a.Price().Amount().String(),
			Comment:     a.Comment(),
		}

		dayMap[dateKey] = append(dayMap[dateKey], dto)
		if statusMap[dateKey] == nil {
			statusMap[dateKey] = make(map[string]int)
		}
		statusMap[dateKey][a.Status().String()]++
	}

	// Build sorted days
	days := make([]CalendarDay, 0, len(dayMap))
	current := q.From
	for current.Before(q.To) {
		dateKey := current.Format("2006-01-02")
		appts := dayMap[dateKey]
		if appts == nil {
			appts = []CalendarAppointment{}
		}
		counts := statusMap[dateKey]
		if counts == nil {
			counts = map[string]int{}
		}
		days = append(days, CalendarDay{
			Date:         dateKey,
			Appointments: appts,
			TotalCount:   len(appts),
			StatusCounts: counts,
		})
		current = current.AddDate(0, 0, 1)
	}

	return &CalendarResult{
		Days: days,
		From: q.From.Format("2006-01-02"),
		To:   q.To.Format("2006-01-02"),
	}, nil
}
