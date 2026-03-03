package services

import (
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/loyalty"
	"math"
)

// LoyaltyPolicy is a domain service for loyalty calculations.
type LoyaltyPolicy struct{}

// NewLoyaltyPolicy creates a new LoyaltyPolicy.
func NewLoyaltyPolicy() *LoyaltyPolicy {
	return &LoyaltyPolicy{}
}

// CalculatePointsForVisit calculates points earned for a visit.
// 1 point per 10 RUB × tier multiplier.
func (lp *LoyaltyPolicy) CalculatePointsForVisit(amount model.Money, tier loyalty.LoyaltyTier) loyalty.Points {
	basePoints := amount.Amount().IntPart() / 10
	multiplied := float64(basePoints) * tier.PointsMultiplier()
	return loyalty.MustNewPoints(int(math.Round(multiplied)))
}

// DetermineNewTier determines the new tier based on lifetime points.
func (lp *LoyaltyPolicy) DetermineNewTier(lifetimePoints loyalty.Points) loyalty.LoyaltyTier {
	return loyalty.DetermineNewTier(lifetimePoints)
}

// GetReferralBonus returns the referral bonus (500 points).
func (lp *LoyaltyPolicy) GetReferralBonus() loyalty.Points {
	return loyalty.MustNewPoints(500)
}

// GetPersonalDiscount returns the personal discount for a tier.
func (lp *LoyaltyPolicy) GetPersonalDiscount(tier loyalty.LoyaltyTier) model.Discount {
	return model.MustNewDiscount(tier.DiscountPercent())
}
