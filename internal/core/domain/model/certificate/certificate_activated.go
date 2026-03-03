package certificate

import (
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"time"

	"github.com/google/uuid"
)

// CertificateActivated is a domain event raised when a certificate is activated.
type CertificateActivated struct {
	eventID             uuid.UUID
	certificateID       uuid.UUID
	activatedByClientID uuid.UUID
	purchasedByClientID uuid.UUID
	balance             model.Money
	expiresAt           time.Time
}

// NewCertificateActivated creates a new CertificateActivated event.
func NewCertificateActivated(
	certificateID, activatedByClientID, purchasedByClientID uuid.UUID,
	balance model.Money,
	expiresAt time.Time,
) CertificateActivated {
	return CertificateActivated{
		eventID:             uuid.New(),
		certificateID:       certificateID,
		activatedByClientID: activatedByClientID,
		purchasedByClientID: purchasedByClientID,
		balance:             balance,
		expiresAt:           expiresAt,
	}
}

func (e CertificateActivated) GetID() uuid.UUID               { return e.eventID }
func (e CertificateActivated) GetName() string                { return "certificate.activated" }
func (e CertificateActivated) CertificateID() uuid.UUID       { return e.certificateID }
func (e CertificateActivated) ActivatedByClientID() uuid.UUID { return e.activatedByClientID }
func (e CertificateActivated) PurchasedByClientID() uuid.UUID { return e.purchasedByClientID }
func (e CertificateActivated) Balance() model.Money           { return e.balance }
func (e CertificateActivated) ExpiresAt() time.Time           { return e.expiresAt }
