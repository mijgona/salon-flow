package loyalty

import "github.com/google/uuid"

// PointsEarned is a domain event raised when loyalty points are earned.
type PointsEarned struct {
	eventID          uuid.UUID
	loyaltyAccountID uuid.UUID
	clientID         uuid.UUID
	pointsEarned     int
	multiplier       float64
	reason           string
	relatedEntityID  uuid.UUID
	newBalance       int
	lifetimePoints   int
}

// NewPointsEarned creates a new PointsEarned event.
func NewPointsEarned(
	loyaltyAccountID, clientID uuid.UUID,
	pointsEarned int,
	multiplier float64,
	reason string,
	relatedEntityID uuid.UUID,
	newBalance, lifetimePoints int,
) PointsEarned {
	return PointsEarned{
		eventID:          uuid.New(),
		loyaltyAccountID: loyaltyAccountID,
		clientID:         clientID,
		pointsEarned:     pointsEarned,
		multiplier:       multiplier,
		reason:           reason,
		relatedEntityID:  relatedEntityID,
		newBalance:       newBalance,
		lifetimePoints:   lifetimePoints,
	}
}

func (e PointsEarned) GetID() uuid.UUID            { return e.eventID }
func (e PointsEarned) GetName() string             { return "loyalty.points_earned" }
func (e PointsEarned) LoyaltyAccountID() uuid.UUID { return e.loyaltyAccountID }
func (e PointsEarned) ClientID() uuid.UUID         { return e.clientID }
func (e PointsEarned) PointsAmount() int           { return e.pointsEarned }
func (e PointsEarned) Multiplier() float64         { return e.multiplier }
func (e PointsEarned) Reason() string              { return e.reason }
func (e PointsEarned) RelatedEntityID() uuid.UUID  { return e.relatedEntityID }
func (e PointsEarned) NewBalance() int             { return e.newBalance }
func (e PointsEarned) LifetimePoints() int         { return e.lifetimePoints }
