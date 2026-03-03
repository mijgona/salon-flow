package certificate

import (
	"errors"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"
	"time"

	"github.com/google/uuid"
)

// CertificateStatus represents the lifecycle status of a certificate.
type CertificateStatus string

const (
	CertificateStatusCreated   CertificateStatus = "created"
	CertificateStatusActivated CertificateStatus = "activated"
	CertificateStatusUsed      CertificateStatus = "used"
	CertificateStatusExpired   CertificateStatus = "expired"
)

// Certificate is the aggregate root for gift cards and subscriptions.
type Certificate struct {
	*ddd.BaseAggregate[uuid.UUID]
	tenantID    model.TenantID
	purchasedBy uuid.UUID
	activatedBy uuid.UUID
	balance     model.Money
	status      CertificateStatus
	activatedAt time.Time
	expiresAt   time.Time
	createdAt   time.Time
}

// NewCertificate creates a new Certificate aggregate.
func NewCertificate(
	tenantID model.TenantID,
	purchasedBy uuid.UUID,
	balance model.Money,
	expiresAt time.Time,
) (*Certificate, error) {
	if balance.IsZero() {
		return nil, errors.New("certificate balance must be positive")
	}
	if expiresAt.Before(time.Now()) {
		return nil, errors.New("expiration date must be in the future")
	}

	return &Certificate{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](uuid.New()),
		tenantID:      tenantID,
		purchasedBy:   purchasedBy,
		balance:       balance,
		status:        CertificateStatusCreated,
		expiresAt:     expiresAt,
		createdAt:     time.Now(),
	}, nil
}

// RestoreCertificate rehydrates a Certificate from the database.
func RestoreCertificate(
	id uuid.UUID,
	tenantID model.TenantID,
	purchasedBy, activatedBy uuid.UUID,
	balance model.Money,
	status CertificateStatus,
	activatedAt, expiresAt, createdAt time.Time,
) *Certificate {
	return &Certificate{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),
		tenantID:      tenantID,
		purchasedBy:   purchasedBy,
		activatedBy:   activatedBy,
		balance:       balance,
		status:        status,
		activatedAt:   activatedAt,
		expiresAt:     expiresAt,
		createdAt:     createdAt,
	}
}

// Activate activates the certificate for a client.
func (c *Certificate) Activate(clientID uuid.UUID) error {
	if c.status != CertificateStatusCreated {
		return errors.New("certificate can only be activated from 'created' status")
	}
	if c.IsExpired() {
		return errors.New("certificate has expired")
	}

	c.activatedBy = clientID
	c.status = CertificateStatusActivated
	c.activatedAt = time.Now()

	c.RaiseDomainEvent(NewCertificateActivated(
		c.ID(), clientID, c.purchasedBy,
		c.balance, c.expiresAt,
	))

	return nil
}

// Deduct deducts an amount from the certificate balance.
func (c *Certificate) Deduct(amount model.Money) error {
	if c.status != CertificateStatusActivated {
		return errors.New("certificate must be activated to deduct")
	}
	if c.IsExpired() {
		return errors.New("certificate has expired")
	}

	newBalance, err := c.balance.Subtract(amount)
	if err != nil {
		return errors.New("insufficient certificate balance")
	}
	c.balance = newBalance

	if c.balance.IsZero() {
		c.status = CertificateStatusUsed
	}

	return nil
}

// IsExpired checks if the certificate has expired.
func (c *Certificate) IsExpired() bool {
	return time.Now().After(c.expiresAt)
}

// Getters
func (c *Certificate) TenantID() model.TenantID  { return c.tenantID }
func (c *Certificate) PurchasedBy() uuid.UUID    { return c.purchasedBy }
func (c *Certificate) ActivatedBy() uuid.UUID    { return c.activatedBy }
func (c *Certificate) Balance() model.Money      { return c.balance }
func (c *Certificate) Status() CertificateStatus { return c.status }
func (c *Certificate) ActivatedAt() time.Time    { return c.activatedAt }
func (c *Certificate) ExpiresAt() time.Time      { return c.expiresAt }
func (c *Certificate) CreatedAt() time.Time      { return c.createdAt }
