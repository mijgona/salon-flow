package loyalty

import "github.com/google/uuid"

// TierChanged is a domain event raised when a client's loyalty tier changes.
type TierChanged struct {
	eventID            uuid.UUID
	loyaltyAccountID   uuid.UUID
	clientID           uuid.UUID
	previousTier       LoyaltyTier
	newTier            LoyaltyTier
	lifetimePoints     int
	newDiscountPercent int
}

// NewTierChanged creates a new TierChanged event.
func NewTierChanged(
	loyaltyAccountID, clientID uuid.UUID,
	previousTier, newTier LoyaltyTier,
	lifetimePoints, newDiscountPercent int,
) TierChanged {
	return TierChanged{
		eventID:            uuid.New(),
		loyaltyAccountID:   loyaltyAccountID,
		clientID:           clientID,
		previousTier:       previousTier,
		newTier:            newTier,
		lifetimePoints:     lifetimePoints,
		newDiscountPercent: newDiscountPercent,
	}
}

func (e TierChanged) GetID() uuid.UUID            { return e.eventID }
func (e TierChanged) GetName() string             { return "loyalty.tier_changed" }
func (e TierChanged) LoyaltyAccountID() uuid.UUID { return e.loyaltyAccountID }
func (e TierChanged) ClientID() uuid.UUID         { return e.clientID }
func (e TierChanged) PreviousTier() LoyaltyTier   { return e.previousTier }
func (e TierChanged) NewTier() LoyaltyTier        { return e.newTier }
func (e TierChanged) LifetimePoints() int         { return e.lifetimePoints }
func (e TierChanged) NewDiscountPercent() int     { return e.newDiscountPercent }
