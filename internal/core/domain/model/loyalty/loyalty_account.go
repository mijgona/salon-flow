package loyalty

import (
	"errors"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"

	"github.com/google/uuid"
)

// LoyaltyAccount is the aggregate root for the loyalty program.
type LoyaltyAccount struct {
	*ddd.BaseAggregate[uuid.UUID]
	clientID       uuid.UUID
	tenantID       model.TenantID
	tier           LoyaltyTier
	balance        Points
	lifetimePoints Points
	transactions   []PointsTransaction
	referrals      []Referral
}

// NewLoyaltyAccount creates a new LoyaltyAccount with Bronze tier and 0 points.
func NewLoyaltyAccount(clientID uuid.UUID, tenantID model.TenantID) (*LoyaltyAccount, error) {
	if clientID == uuid.Nil {
		return nil, errors.New("client ID is required")
	}

	return &LoyaltyAccount{
		BaseAggregate:  ddd.NewBaseAggregate[uuid.UUID](uuid.New()),
		clientID:       clientID,
		tenantID:       tenantID,
		tier:           TierBronze,
		balance:        ZeroPoints(),
		lifetimePoints: ZeroPoints(),
		transactions:   make([]PointsTransaction, 0),
		referrals:      make([]Referral, 0),
	}, nil
}

// RestoreLoyaltyAccount rehydrates a LoyaltyAccount from the database.
func RestoreLoyaltyAccount(
	id, clientID uuid.UUID,
	tenantID model.TenantID,
	tier LoyaltyTier,
	balance, lifetimePoints Points,
	transactions []PointsTransaction,
	referrals []Referral,
) *LoyaltyAccount {
	return &LoyaltyAccount{
		BaseAggregate:  ddd.NewBaseAggregate[uuid.UUID](id),
		clientID:       clientID,
		tenantID:       tenantID,
		tier:           tier,
		balance:        balance,
		lifetimePoints: lifetimePoints,
		transactions:   transactions,
		referrals:      referrals,
	}
}

// EarnPoints adds points to the account.
func (la *LoyaltyAccount) EarnPoints(amount int, reason string, relatedEntityID uuid.UUID) {
	la.balance = la.balance.Add(MustNewPoints(amount))
	la.lifetimePoints = la.lifetimePoints.Add(MustNewPoints(amount))

	tx := NewPointsTransaction(amount, TransactionTypeEarn, reason, relatedEntityID)
	la.transactions = append(la.transactions, tx)

	la.RaiseDomainEvent(NewPointsEarned(
		la.ID(), la.clientID,
		amount, la.tier.PointsMultiplier(),
		reason, relatedEntityID,
		la.balance.Value(), la.lifetimePoints.Value(),
	))
}

// RedeemPoints subtracts points from the balance.
func (la *LoyaltyAccount) RedeemPoints(amount int, reason string, relatedEntityID uuid.UUID) error {
	pts := MustNewPoints(amount)
	newBalance, err := la.balance.Subtract(pts)
	if err != nil {
		return errors.New("insufficient points balance")
	}
	la.balance = newBalance

	tx := NewPointsTransaction(amount, TransactionTypeRedeem, reason, relatedEntityID)
	la.transactions = append(la.transactions, tx)

	return nil
}

// RecalculateTier recalculates the tier based on lifetime points. Tier only changes upward.
func (la *LoyaltyAccount) RecalculateTier() {
	newTier := DetermineNewTier(la.lifetimePoints)
	if newTier.IsHigherThan(la.tier) {
		previousTier := la.tier
		la.tier = newTier

		la.RaiseDomainEvent(NewTierChanged(
			la.ID(), la.clientID,
			previousTier, newTier,
			la.lifetimePoints.Value(),
			newTier.DiscountPercent(),
		))
	}
}

// AddReferral adds a referral. Returns error if the referred client already has a referral.
func (la *LoyaltyAccount) AddReferral(referredClientID uuid.UUID) error {
	for _, existing := range la.referrals {
		if existing.ReferredClientID() == referredClientID {
			return errors.New("referral for this client already exists")
		}
	}
	referral := NewReferral(referredClientID)
	la.referrals = append(la.referrals, referral)
	return nil
}

// GetPersonalDiscount returns the discount based on the current tier.
func (la *LoyaltyAccount) GetPersonalDiscount() model.Discount {
	return model.MustNewDiscount(la.tier.DiscountPercent())
}

// Getters
func (la *LoyaltyAccount) ClientID() uuid.UUID               { return la.clientID }
func (la *LoyaltyAccount) TenantID() model.TenantID          { return la.tenantID }
func (la *LoyaltyAccount) Tier() LoyaltyTier                 { return la.tier }
func (la *LoyaltyAccount) Balance() Points                   { return la.balance }
func (la *LoyaltyAccount) LifetimePoints() Points            { return la.lifetimePoints }
func (la *LoyaltyAccount) Transactions() []PointsTransaction { return la.transactions }
func (la *LoyaltyAccount) Referrals() []Referral             { return la.referrals }
