package ports

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"time"

	"github.com/google/uuid"
)

// ServiceCatalogItem represents a service in the catalog.
type ServiceCatalogItem struct {
	ID       uuid.UUID
	Name     string
	Duration time.Duration
	Price    model.Money
}

// ServiceCatalogClient provides service catalog data.
type ServiceCatalogClient interface {
	GetService(ctx context.Context, serviceID uuid.UUID) (*ServiceCatalogItem, error)
	ListServices(ctx context.Context, tenantID model.TenantID) ([]ServiceCatalogItem, error)
}
