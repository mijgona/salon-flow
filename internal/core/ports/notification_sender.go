package ports

import "context"

// NotificationSender sends notifications to clients.
type NotificationSender interface {
	SendSMS(ctx context.Context, phone, message string) error
	SendWhatsApp(ctx context.Context, phone, message string) error
	SendEmail(ctx context.Context, email, subject, body string) error
	ScheduleNotification(ctx context.Context, phone, message string, sendAt interface{}) error
}
