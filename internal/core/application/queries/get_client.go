package queries

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model/client"
	"github.com/mijgona/salon-crm/internal/core/ports"

	"github.com/google/uuid"
)

// GetClientQuery holds data for retrieving a client.
type GetClientQuery struct {
	ClientID uuid.UUID
}

// GetClientResult is the result of the GetClient query.
type GetClientResult struct {
	Client *client.Client
}

// GetClientHandler handles the get client query.
type GetClientHandler struct {
	clientRepo ports.ClientRepository
}

// NewGetClientHandler creates a new handler.
func NewGetClientHandler(clientRepo ports.ClientRepository) *GetClientHandler {
	return &GetClientHandler{clientRepo: clientRepo}
}

// Handle executes the get client query.
func (h *GetClientHandler) Handle(ctx context.Context, q GetClientQuery) (*GetClientResult, error) {
	c, err := h.clientRepo.Get(ctx, nil, q.ClientID)
	if err != nil {
		return nil, err
	}
	return &GetClientResult{Client: c}, nil
}
