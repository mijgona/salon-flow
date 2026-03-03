package loyalty

// LoyaltyTier represents a loyalty program tier.
type LoyaltyTier string

const (
	TierBronze LoyaltyTier = "Bronze"
	TierSilver LoyaltyTier = "Silver"
	TierGold   LoyaltyTier = "Gold"
	TierVIP    LoyaltyTier = "VIP"
)

// IsValid checks if the tier is a known value.
func (t LoyaltyTier) IsValid() bool {
	switch t {
	case TierBronze, TierSilver, TierGold, TierVIP:
		return true
	}
	return false
}

// DiscountPercent returns the discount percentage for the tier.
func (t LoyaltyTier) DiscountPercent() int {
	switch t {
	case TierBronze:
		return 0
	case TierSilver:
		return 5
	case TierGold:
		return 10
	case TierVIP:
		return 15
	}
	return 0
}

// PointsMultiplier returns the points earning multiplier for the tier.
func (t LoyaltyTier) PointsMultiplier() float64 {
	switch t {
	case TierBronze:
		return 1.0
	case TierSilver:
		return 1.2
	case TierGold:
		return 1.5
	case TierVIP:
		return 2.0
	}
	return 1.0
}

// String returns the string representation.
func (t LoyaltyTier) String() string { return string(t) }

// IsHigherThan checks if this tier is higher than another.
func (t LoyaltyTier) IsHigherThan(other LoyaltyTier) bool {
	return tierOrder(t) > tierOrder(other)
}

func tierOrder(t LoyaltyTier) int {
	switch t {
	case TierBronze:
		return 0
	case TierSilver:
		return 1
	case TierGold:
		return 2
	case TierVIP:
		return 3
	}
	return -1
}
