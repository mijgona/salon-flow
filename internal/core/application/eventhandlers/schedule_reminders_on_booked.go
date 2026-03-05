package eventhandlers

import (
	"context"
	"time"

	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"
	"github.com/mijgona/salon-crm/internal/core/ports"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"
)

// ScheduleRemindersOnBookedHandler schedules notification reminders when an appointment is booked.
type ScheduleRemindersOnBookedHandler struct {
	notificationSender ports.NotificationSender
}

// NewScheduleRemindersOnBookedHandler creates a new handler.
func NewScheduleRemindersOnBookedHandler(notificationSender ports.NotificationSender) *ScheduleRemindersOnBookedHandler {
	return &ScheduleRemindersOnBookedHandler{notificationSender: notificationSender}
}

// Handle processes the AppointmentBooked event.
func (h *ScheduleRemindersOnBookedHandler) Handle(ctx context.Context, event ddd.DomainEvent) error {
	booked, ok := event.(scheduling.AppointmentBooked)
	if !ok {
		return nil
	}

	startTime, err := time.Parse(time.RFC3339, booked.StartTime())
	if err != nil {
		return err
	}

	// Schedule 24h reminder
	reminder24h := startTime.Add(-24 * time.Hour)
	message24h := "Напоминаем о записи на " + booked.ServiceName() + " завтра в " + startTime.Format("15:04")
	_ = h.notificationSender.ScheduleNotification(ctx, "", message24h, reminder24h)

	// Schedule 2h reminder
	reminder2h := startTime.Add(-2 * time.Hour)
	message2h := "Через 2 часа у вас " + booked.ServiceName() + " в " + startTime.Format("15:04")
	_ = h.notificationSender.ScheduleNotification(ctx, "", message2h, reminder2h)

	return nil
}
