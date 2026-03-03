package client

import (
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents payment status for a visit.
type PaymentStatus string

const (
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusUnpaid  PaymentStatus = "unpaid"
	PaymentStatusPartial PaymentStatus = "partial"
)

// VisitRecord is an entity representing a client's visit history entry.
type VisitRecord struct {
	appointmentID uuid.UUID
	masterID      uuid.UUID
	service       string
	price         model.Money
	discount      model.Discount
	paymentStatus PaymentStatus
	rating        int
	review        string
	visitedAt     time.Time
}

// NewVisitRecord creates a new VisitRecord entity.
func NewVisitRecord(
	appointmentID, masterID uuid.UUID,
	service string,
	price model.Money,
	discount model.Discount,
	paymentStatus PaymentStatus,
	visitedAt time.Time,
) VisitRecord {
	return VisitRecord{
		appointmentID: appointmentID,
		masterID:      masterID,
		service:       service,
		price:         price,
		discount:      discount,
		paymentStatus: paymentStatus,
		visitedAt:     visitedAt,
	}
}

// RestoreVisitRecord creates a VisitRecord from persisted data.
func RestoreVisitRecord(
	appointmentID, masterID uuid.UUID,
	service string,
	price model.Money,
	discount model.Discount,
	paymentStatus PaymentStatus,
	rating int,
	review string,
	visitedAt time.Time,
) VisitRecord {
	return VisitRecord{
		appointmentID: appointmentID,
		masterID:      masterID,
		service:       service,
		price:         price,
		discount:      discount,
		paymentStatus: paymentStatus,
		rating:        rating,
		review:        review,
		visitedAt:     visitedAt,
	}
}

// AppointmentID returns the associated appointment ID.
func (vr VisitRecord) AppointmentID() uuid.UUID { return vr.appointmentID }

// MasterID returns the master who performed the service.
func (vr VisitRecord) MasterID() uuid.UUID { return vr.masterID }

// Service returns the service name.
func (vr VisitRecord) Service() string { return vr.service }

// Price returns the price paid.
func (vr VisitRecord) Price() model.Money { return vr.price }

// Discount returns the discount applied.
func (vr VisitRecord) Discount() model.Discount { return vr.discount }

// PaymentStatus returns the payment status.
func (vr VisitRecord) PaymentStatus() PaymentStatus { return vr.paymentStatus }

// Rating returns the rating (0-5).
func (vr VisitRecord) Rating() int { return vr.rating }

// Review returns the review text.
func (vr VisitRecord) Review() string { return vr.review }

// VisitedAt returns the visit timestamp.
func (vr VisitRecord) VisitedAt() time.Time { return vr.visitedAt }

// SetRating sets a rating for the visit.
func (vr *VisitRecord) SetRating(rating int, review string) {
	if rating < 0 {
		rating = 0
	}
	if rating > 5 {
		rating = 5
	}
	vr.rating = rating
	vr.review = review
}
