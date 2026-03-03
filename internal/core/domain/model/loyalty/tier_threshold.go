package loyalty

// TierThreshold defines the minimum lifetime points required for a tier.
type TierThreshold struct {
	tier      LoyaltyTier
	minPoints Points
}

// DefaultTierThresholds returns the standard tier thresholds.
func DefaultTierThresholds() []TierThreshold {
	return []TierThreshold{
		{tier: TierVIP, minPoints: MustNewPoints(50000)},
		{tier: TierGold, minPoints: MustNewPoints(15000)},
		{tier: TierSilver, minPoints: MustNewPoints(5000)},
		{tier: TierBronze, minPoints: MustNewPoints(0)},
	}
}

// DetermineNewTier returns the highest tier the given lifetime points qualify for.
func DetermineNewTier(lifetimePoints Points) LoyaltyTier {
	for _, threshold := range DefaultTierThresholds() {
		if lifetimePoints.GreaterThanOrEqual(threshold.minPoints) {
			return threshold.tier
		}
	}
	return TierBronze
}

// Tier returns the tier.
func (tt TierThreshold) Tier() LoyaltyTier { return tt.tier }

// MinPoints returns the minimum points required.
func (tt TierThreshold) MinPoints() Points { return tt.minPoints }
