package certificate

import (
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func makeTestTenantID() model.TenantID {
	return model.MustNewTenantID(uuid.New())
}

func TestNewCertificate_HappyPath(t *testing.T) {
	balance := model.MustNewMoney(decimal.NewFromInt(5000), "RUB")
	expiresAt := time.Now().Add(365 * 24 * time.Hour)

	cert, err := NewCertificate(makeTestTenantID(), uuid.New(), balance, expiresAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cert.Status() != CertificateStatusCreated {
		t.Errorf("expected status 'created', got %q", cert.Status())
	}
}

func TestCertificate_Activate(t *testing.T) {
	balance := model.MustNewMoney(decimal.NewFromInt(5000), "RUB")
	expiresAt := time.Now().Add(365 * 24 * time.Hour)

	cert, _ := NewCertificate(makeTestTenantID(), uuid.New(), balance, expiresAt)

	clientID := uuid.New()
	if err := cert.Activate(clientID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cert.Status() != CertificateStatusActivated {
		t.Errorf("expected status 'activated', got %q", cert.Status())
	}
	if cert.ActivatedBy() != clientID {
		t.Error("activated by client ID mismatch")
	}

	events := cert.GetDomainEvents()
	found := false
	for _, e := range events {
		if e.GetName() == "certificate.activated" {
			found = true
		}
	}
	if !found {
		t.Error("expected CertificateActivated event")
	}
}

func TestCertificate_CannotActivateTwice(t *testing.T) {
	balance := model.MustNewMoney(decimal.NewFromInt(5000), "RUB")
	expiresAt := time.Now().Add(365 * 24 * time.Hour)

	cert, _ := NewCertificate(makeTestTenantID(), uuid.New(), balance, expiresAt)
	_ = cert.Activate(uuid.New())

	if err := cert.Activate(uuid.New()); err == nil {
		t.Fatal("expected error for double activation")
	}
}

func TestCertificate_Deduct(t *testing.T) {
	balance := model.MustNewMoney(decimal.NewFromInt(5000), "RUB")
	expiresAt := time.Now().Add(365 * 24 * time.Hour)

	cert, _ := NewCertificate(makeTestTenantID(), uuid.New(), balance, expiresAt)
	_ = cert.Activate(uuid.New())

	deductAmount := model.MustNewMoney(decimal.NewFromInt(2000), "RUB")
	if err := cert.Deduct(deductAmount); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := decimal.NewFromInt(3000)
	if !cert.Balance().Amount().Equal(expected) {
		t.Errorf("expected balance 3000, got %s", cert.Balance().Amount().String())
	}
}

func TestCertificate_DeductFull_UsesUp(t *testing.T) {
	balance := model.MustNewMoney(decimal.NewFromInt(5000), "RUB")
	expiresAt := time.Now().Add(365 * 24 * time.Hour)

	cert, _ := NewCertificate(makeTestTenantID(), uuid.New(), balance, expiresAt)
	_ = cert.Activate(uuid.New())

	if err := cert.Deduct(balance); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cert.Status() != CertificateStatusUsed {
		t.Errorf("expected status 'used', got %q", cert.Status())
	}
}

func TestCertificate_DeductInsufficientBalance(t *testing.T) {
	balance := model.MustNewMoney(decimal.NewFromInt(1000), "RUB")
	expiresAt := time.Now().Add(365 * 24 * time.Hour)

	cert, _ := NewCertificate(makeTestTenantID(), uuid.New(), balance, expiresAt)
	_ = cert.Activate(uuid.New())

	overAmount := model.MustNewMoney(decimal.NewFromInt(5000), "RUB")
	if err := cert.Deduct(overAmount); err == nil {
		t.Fatal("expected error for insufficient balance")
	}
}
