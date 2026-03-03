package commands

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/client"
	"github.com/mijgona/salon-crm/internal/core/ports"

	"github.com/google/uuid"
)

// RegisterClientCommand holds data for registering a new client.
type RegisterClientCommand struct {
	TenantID           uuid.UUID
	Phone              string
	Email              string
	FirstName          string
	LastName           string
	Source             string
	ReferredByClientID uuid.UUID
}

// RegisterClientHandler handles client registration.
type RegisterClientHandler struct {
	clientRepo ports.ClientRepository
	txManager  ports.TxManager
}

// NewRegisterClientHandler creates a new handler.
func NewRegisterClientHandler(clientRepo ports.ClientRepository, txManager ports.TxManager) *RegisterClientHandler {
	return &RegisterClientHandler{clientRepo: clientRepo, txManager: txManager}
}

// Handle executes the register client command.
func (h *RegisterClientHandler) Handle(ctx context.Context, cmd RegisterClientCommand) (uuid.UUID, error) {
	tenantID, err := model.NewTenantID(cmd.TenantID)
	if err != nil {
		return uuid.Nil, err
	}

	phone, err := model.NewPhoneNumber(cmd.Phone)
	if err != nil {
		return uuid.Nil, err
	}

	contactInfo, err := client.NewContactInfo(phone, cmd.Email, cmd.FirstName, cmd.LastName)
	if err != nil {
		return uuid.Nil, err
	}

	source := client.ClientSource(cmd.Source)

	var clientID uuid.UUID
	err = h.txManager.Execute(ctx, func(tx interface{}) error {
		// Check for existing client with same phone
		existing, _ := h.clientRepo.FindByPhone(ctx, tx, tenantID, phone)
		if existing != nil {
			return &DuplicateError{Field: "phone", Value: cmd.Phone}
		}

		c, err := client.NewClient(tenantID, contactInfo, source, cmd.ReferredByClientID)
		if err != nil {
			return err
		}

		if err := h.clientRepo.Add(ctx, tx, c); err != nil {
			return err
		}

		clientID = c.ID()
		return nil
	})

	return clientID, err
}

// DuplicateError indicates a duplicate entity.
type DuplicateError struct {
	Field string
	Value string
}

func (e *DuplicateError) Error() string {
	return e.Field + " already exists: " + e.Value
}
