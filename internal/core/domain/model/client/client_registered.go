package client

import (
	"github.com/google/uuid"
)

// ClientRegistered is a domain event raised when a new client is registered.
type ClientRegistered struct {
	eventID            uuid.UUID
	clientID           uuid.UUID
	tenantID           uuid.UUID
	firstName          string
	lastName           string
	phone              string
	source             ClientSource
	referredByClientID uuid.UUID
}

// NewClientRegistered creates a new ClientRegistered domain event.
func NewClientRegistered(
	clientID, tenantID uuid.UUID,
	firstName, lastName, phone string,
	source ClientSource,
	referredByClientID uuid.UUID,
) ClientRegistered {
	return ClientRegistered{
		eventID:            uuid.New(),
		clientID:           clientID,
		tenantID:           tenantID,
		firstName:          firstName,
		lastName:           lastName,
		phone:              phone,
		source:             source,
		referredByClientID: referredByClientID,
	}
}

func (e ClientRegistered) GetID() uuid.UUID              { return e.eventID }
func (e ClientRegistered) GetName() string               { return "client.registered" }
func (e ClientRegistered) ClientID() uuid.UUID           { return e.clientID }
func (e ClientRegistered) TenantID() uuid.UUID           { return e.tenantID }
func (e ClientRegistered) FirstName() string             { return e.firstName }
func (e ClientRegistered) LastName() string              { return e.lastName }
func (e ClientRegistered) Phone() string                 { return e.phone }
func (e ClientRegistered) Source() ClientSource          { return e.source }
func (e ClientRegistered) ReferredByClientID() uuid.UUID { return e.referredByClientID }
