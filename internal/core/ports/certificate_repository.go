package ports

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model/certificate"

	"github.com/google/uuid"
)

// CertificateRepository defines operations for persisting Certificate aggregates.
type CertificateRepository interface {
	Add(ctx context.Context, tx interface{}, c *certificate.Certificate) error
	Update(ctx context.Context, tx interface{}, c *certificate.Certificate) error
	Get(ctx context.Context, tx interface{}, id uuid.UUID) (*certificate.Certificate, error)
}
