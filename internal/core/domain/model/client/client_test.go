package client

import (
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func makeTestContactInfo() ContactInfo {
	phone := model.MustNewPhoneNumber("+79001234567")
	return MustNewContactInfo(phone, "test@example.com", "Алина", "Иванова")
}

func makeTestTenantID() model.TenantID {
	return model.MustNewTenantID(uuid.New())
}

func TestNewClient_HappyPath(t *testing.T) {
	tenantID := makeTestTenantID()
	contactInfo := makeTestContactInfo()

	c, err := NewClient(tenantID, contactInfo, ClientSourceOnlineBooking, uuid.Nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.ID() == uuid.Nil {
		t.Fatal("expected non-nil ID")
	}
	if c.ContactInfo().FirstName() != "Алина" {
		t.Errorf("expected first name 'Алина', got %q", c.ContactInfo().FirstName())
	}
	if c.Source() != ClientSourceOnlineBooking {
		t.Errorf("expected source 'online_booking', got %q", c.Source())
	}
	if c.TotalVisits() != 0 {
		t.Errorf("expected 0 visits, got %d", c.TotalVisits())
	}
}

func TestNewClient_RaisesClientRegisteredEvent(t *testing.T) {
	tenantID := makeTestTenantID()
	contactInfo := makeTestContactInfo()

	c, err := NewClient(tenantID, contactInfo, ClientSourceAdminEntry, uuid.Nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := c.GetDomainEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event, ok := events[0].(ClientRegistered)
	if !ok {
		t.Fatal("expected ClientRegistered event")
	}
	if event.GetName() != "client.registered" {
		t.Errorf("expected event name 'client.registered', got %q", event.GetName())
	}
	if event.ClientID() != c.ID() {
		t.Error("event clientID does not match aggregate ID")
	}
}

func TestNewClient_InvalidSource(t *testing.T) {
	tenantID := makeTestTenantID()
	contactInfo := makeTestContactInfo()

	_, err := NewClient(tenantID, contactInfo, "invalid", uuid.Nil)
	if err == nil {
		t.Fatal("expected error for invalid source")
	}
}

func TestClient_AddAllergy_Deduplication(t *testing.T) {
	tenantID := makeTestTenantID()
	contactInfo := makeTestContactInfo()
	c, _ := NewClient(tenantID, contactInfo, ClientSourceWalkIn, uuid.Nil)

	allergy1 := MustNewAllergy("парабены", SeverityHigh)
	allergy2 := MustNewAllergy("парабены", SeverityLow) // same substance

	if err := c.AddAllergy(allergy1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := c.AddAllergy(allergy2); err == nil {
		t.Fatal("expected error for duplicate allergy substance")
	}
	if len(c.Allergies()) != 1 {
		t.Errorf("expected 1 allergy, got %d", len(c.Allergies()))
	}
}

func TestClient_AddVisitRecord_TotalSpent(t *testing.T) {
	tenantID := makeTestTenantID()
	contactInfo := makeTestContactInfo()
	c, _ := NewClient(tenantID, contactInfo, ClientSourceWalkIn, uuid.Nil)

	price1 := model.MustNewMoney(decimal.NewFromInt(3000), "RUB")
	price2 := model.MustNewMoney(decimal.NewFromInt(5000), "RUB")

	vr1 := NewVisitRecord(uuid.New(), uuid.New(), "Стрижка", price1, model.ZeroDiscount(), PaymentStatusPaid, time.Now())
	vr2 := NewVisitRecord(uuid.New(), uuid.New(), "Окрашивание", price2, model.ZeroDiscount(), PaymentStatusPaid, time.Now())

	c.AddVisitRecord(vr1)
	c.AddVisitRecord(vr2)

	if c.TotalVisits() != 2 {
		t.Errorf("expected 2 visits, got %d", c.TotalVisits())
	}

	expectedTotal := decimal.NewFromInt(8000)
	if !c.TotalSpentDecimal().Equal(expectedTotal) {
		t.Errorf("expected total spent 8000, got %s", c.TotalSpentDecimal().String())
	}
}

func TestClient_AddNote(t *testing.T) {
	tenantID := makeTestTenantID()
	contactInfo := makeTestContactInfo()
	c, _ := NewClient(tenantID, contactInfo, ClientSourceWalkIn, uuid.Nil)

	note := MustNewNote("Предпочитает естественные тона", uuid.New())
	c.AddNote(note)

	if len(c.Notes()) != 1 {
		t.Errorf("expected 1 note, got %d", len(c.Notes()))
	}
}

func TestClient_UpdateProfile(t *testing.T) {
	tenantID := makeTestTenantID()
	contactInfo := makeTestContactInfo()
	c, _ := NewClient(tenantID, contactInfo, ClientSourceWalkIn, uuid.Nil)

	newPhone := model.MustNewPhoneNumber("+79009876543")
	newContactInfo := MustNewContactInfo(newPhone, "new@example.com", "Мария", "Петрова")
	birthday := model.MustNewBirthday(time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC))
	prefs := NewPreferences(uuid.New(), []uuid.UUID{uuid.New()}, ChannelWhatsApp)

	c.UpdateProfile(newContactInfo, birthday, prefs)

	if c.ContactInfo().FirstName() != "Мария" {
		t.Errorf("expected updated first name 'Мария', got %q", c.ContactInfo().FirstName())
	}
	if c.Preferences().Channel() != ChannelWhatsApp {
		t.Errorf("expected channel 'whatsapp', got %q", c.Preferences().Channel())
	}
}
