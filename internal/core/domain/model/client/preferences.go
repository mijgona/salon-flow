package client

import "github.com/google/uuid"

// CommunicationChannel represents the client's preferred communication channel.
type CommunicationChannel string

const (
	ChannelSMS      CommunicationChannel = "sms"
	ChannelWhatsApp CommunicationChannel = "whatsapp"
	ChannelEmail    CommunicationChannel = "email"
)

// Preferences holds the client's preferences.
type Preferences struct {
	preferredMasterID uuid.UUID
	favoriteServices  []uuid.UUID
	channel           CommunicationChannel
}

// NewPreferences creates a Preferences value object.
func NewPreferences(preferredMasterID uuid.UUID, favoriteServices []uuid.UUID, channel CommunicationChannel) Preferences {
	if favoriteServices == nil {
		favoriteServices = make([]uuid.UUID, 0)
	}
	if channel == "" {
		channel = ChannelSMS
	}
	return Preferences{
		preferredMasterID: preferredMasterID,
		favoriteServices:  favoriteServices,
		channel:           channel,
	}
}

// EmptyPreferences returns default preferences.
func EmptyPreferences() Preferences {
	return Preferences{
		favoriteServices: make([]uuid.UUID, 0),
		channel:          ChannelSMS,
	}
}

// PreferredMasterID returns the preferred master UUID.
func (p Preferences) PreferredMasterID() uuid.UUID { return p.preferredMasterID }

// FavoriteServices returns the list of favorite service IDs.
func (p Preferences) FavoriteServices() []uuid.UUID { return p.favoriteServices }

// Channel returns the preferred communication channel.
func (p Preferences) Channel() CommunicationChannel { return p.channel }
