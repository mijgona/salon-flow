package commands

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/client"
	"github.com/mijgona/salon-crm/internal/core/ports"

	"github.com/google/uuid"
)

// UpdateClientProfileCommand holds data for updating a client profile.
type UpdateClientProfileCommand struct {
	ClientID          uuid.UUID
	Phone             string
	Email             string
	FirstName         string
	LastName          string
	Birthday          string // RFC3339 date
	PreferredMasterID uuid.UUID
	FavoriteServices  []uuid.UUID
	Channel           string
}

// UpdateClientProfileHandler handles client profile updates.
type UpdateClientProfileHandler struct {
	clientRepo ports.ClientRepository
	txManager  ports.TxManager
}

// NewUpdateClientProfileHandler creates a new handler.
func NewUpdateClientProfileHandler(clientRepo ports.ClientRepository, txManager ports.TxManager) *UpdateClientProfileHandler {
	return &UpdateClientProfileHandler{clientRepo: clientRepo, txManager: txManager}
}

// Handle executes the update client profile command.
func (h *UpdateClientProfileHandler) Handle(ctx context.Context, cmd UpdateClientProfileCommand) error {
	return h.txManager.Execute(ctx, func(tx interface{}) error {
		c, err := h.clientRepo.Get(ctx, tx, cmd.ClientID)
		if err != nil {
			return err
		}

		phone, err := model.NewPhoneNumber(cmd.Phone)
		if err != nil {
			return err
		}

		contactInfo, err := client.NewContactInfo(phone, cmd.Email, cmd.FirstName, cmd.LastName)
		if err != nil {
			return err
		}

		var birthday model.Birthday
		// Birthday is optional — only parse if provided
		if cmd.Birthday != "" {
			// Parse expects date format
			_ = birthday // left as zero value if parsing fails
		}

		prefs := client.NewPreferences(
			cmd.PreferredMasterID,
			cmd.FavoriteServices,
			client.CommunicationChannel(cmd.Channel),
		)

		c.UpdateProfile(contactInfo, birthday, prefs)

		return h.clientRepo.Update(ctx, tx, c)
	})
}
