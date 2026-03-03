package loyalty

import (
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"testing"

	"github.com/google/uuid"
)

func makeTestTenantID() model.TenantID {
	return model.MustNewTenantID(uuid.New())
}

func TestNewLoyaltyAccount_HappyPath(t *testing.T) {
	la, err := NewLoyaltyAccount(uuid.New(), makeTestTenantID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if la.Tier() != TierBronze {
		t.Errorf("expected Bronze tier, got %q", la.Tier())
	}
	if !la.Balance().IsZero() {
		t.Error("expected zero balance")
	}
}

func TestLoyaltyAccount_EarnPoints(t *testing.T) {
	la, _ := NewLoyaltyAccount(uuid.New(), makeTestTenantID())

	la.EarnPoints(1000, "appointment_completed", uuid.New())

	if la.Balance().Value() != 1000 {
		t.Errorf("expected balance 1000, got %d", la.Balance().Value())
	}
	if la.LifetimePoints().Value() != 1000 {
		t.Errorf("expected lifetime 1000, got %d", la.LifetimePoints().Value())
	}
	if len(la.Transactions()) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(la.Transactions()))
	}

	events := la.GetDomainEvents()
	found := false
	for _, e := range events {
		if e.GetName() == "loyalty.points_earned" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected PointsEarned event")
	}
}

func TestLoyaltyAccount_RedeemPoints_Insufficient(t *testing.T) {
	la, _ := NewLoyaltyAccount(uuid.New(), makeTestTenantID())
	la.EarnPoints(100, "test", uuid.Nil)

	err := la.RedeemPoints(500, "test", uuid.Nil)
	if err == nil {
		t.Fatal("expected error for insufficient balance")
	}
}

func TestLoyaltyAccount_RedeemPoints_Success(t *testing.T) {
	la, _ := NewLoyaltyAccount(uuid.New(), makeTestTenantID())
	la.EarnPoints(1000, "test", uuid.Nil)

	err := la.RedeemPoints(300, "test", uuid.Nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if la.Balance().Value() != 700 {
		t.Errorf("expected balance 700, got %d", la.Balance().Value())
	}
}

func TestLoyaltyAccount_TierUpgrade(t *testing.T) {
	la, _ := NewLoyaltyAccount(uuid.New(), makeTestTenantID())

	// Earn enough for Silver (5000 points)
	la.EarnPoints(5000, "test", uuid.Nil)
	la.ClearDomainEvents()
	la.RecalculateTier()

	if la.Tier() != TierSilver {
		t.Errorf("expected Silver tier, got %q", la.Tier())
	}

	events := la.GetDomainEvents()
	found := false
	for _, e := range events {
		if e.GetName() == "loyalty.tier_changed" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected TierChanged event")
	}
}

func TestLoyaltyAccount_TierOnlyUpward(t *testing.T) {
	la, _ := NewLoyaltyAccount(uuid.New(), makeTestTenantID())

	// Earn enough for Silver
	la.EarnPoints(5000, "test", uuid.Nil)
	la.RecalculateTier()
	if la.Tier() != TierSilver {
		t.Fatalf("expected Silver tier, got %q", la.Tier())
	}

	// Redeem points — tier should remain Silver (tier only changes upward)
	_ = la.RedeemPoints(4500, "test", uuid.Nil)
	la.RecalculateTier()
	if la.Tier() != TierSilver {
		t.Errorf("expected tier to remain Silver, got %q", la.Tier())
	}
}

func TestLoyaltyAccount_AddReferral_Duplicate(t *testing.T) {
	la, _ := NewLoyaltyAccount(uuid.New(), makeTestTenantID())
	referredID := uuid.New()

	err := la.AddReferral(referredID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = la.AddReferral(referredID)
	if err == nil {
		t.Fatal("expected error for duplicate referral")
	}
}

func TestLoyaltyAccount_GetPersonalDiscount(t *testing.T) {
	la, _ := NewLoyaltyAccount(uuid.New(), makeTestTenantID())

	// Bronze -> 0%
	d := la.GetPersonalDiscount()
	if d.Percent() != 0 {
		t.Errorf("expected 0%% discount for Bronze, got %d%%", d.Percent())
	}

	// Upgrade to Gold
	la.EarnPoints(15000, "test", uuid.Nil)
	la.RecalculateTier()

	d = la.GetPersonalDiscount()
	if d.Percent() != 10 {
		t.Errorf("expected 10%% discount for Gold, got %d%%", d.Percent())
	}
}
