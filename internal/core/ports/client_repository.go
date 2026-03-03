package ports

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/client"

	"github.com/google/uuid"
)

// ClientRepository defines operations for persisting Client aggregates.
type ClientRepository interface {
	Add(ctx context.Context, tx interface{}, c *client.Client) error
	Update(ctx context.Context, tx interface{}, c *client.Client) error
	Get(ctx context.Context, tx interface{}, id uuid.UUID) (*client.Client, error)
	FindByPhone(ctx context.Context, tx interface{}, tenantID model.TenantID, phone model.PhoneNumber) (*client.Client, error)
	FindByTenant(ctx context.Context, tx interface{}, tenantID model.TenantID, limit, offset int) ([]*client.Client, error)
}
