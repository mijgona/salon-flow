package inmemory

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/client"
	"sync"

	"github.com/google/uuid"
)

// InMemoryClientRepository is an in-memory implementation for testing.
type InMemoryClientRepository struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]*client.Client
}

// NewInMemoryClientRepository creates a new in-memory client repository.
func NewInMemoryClientRepository() *InMemoryClientRepository {
	return &InMemoryClientRepository{
		clients: make(map[uuid.UUID]*client.Client),
	}
}

func (r *InMemoryClientRepository) Add(_ context.Context, _ interface{}, c *client.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c.ID()] = c
	return nil
}

func (r *InMemoryClientRepository) Update(_ context.Context, _ interface{}, c *client.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c.ID()] = c
	return nil
}

func (r *InMemoryClientRepository) Get(_ context.Context, _ interface{}, id uuid.UUID) (*client.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.clients[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *InMemoryClientRepository) FindByPhone(_ context.Context, _ interface{}, tenantID model.TenantID, phone model.PhoneNumber) (*client.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.clients {
		if c.TenantID().Equal(tenantID) && c.ContactInfo().Phone().Equal(phone) {
			return c, nil
		}
	}
	return nil, nil
}

func (r *InMemoryClientRepository) FindByTenant(_ context.Context, _ interface{}, tenantID model.TenantID, limit, offset int) ([]*client.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*client.Client
	for _, c := range r.clients {
		if c.TenantID().Equal(tenantID) {
			result = append(result, c)
		}
	}
	// Apply offset and limit
	if offset >= len(result) {
		return []*client.Client{}, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}
