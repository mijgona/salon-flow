package scheduling

import (
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func futureTime(hoursFromNow int) time.Time {
	return time.Now().Add(time.Duration(hoursFromNow) * time.Hour)
}

func makeTestServiceInfo() ServiceInfo {
	price := model.MustNewMoney(decimal.NewFromInt(3000), "RUB")
	return MustNewServiceInfo(uuid.New(), "Стрижка", 60*time.Minute, price)
}

func makeTestTenantID() model.TenantID {
	return model.MustNewTenantID(uuid.New())
}

func TestNewAppointment_HappyPath(t *testing.T) {
	tenantID := makeTestTenantID()
	serviceInfo := makeTestServiceInfo()
	slot := MustNewTimeSlot(futureTime(24), futureTime(25))
	price := model.MustNewMoney(decimal.NewFromInt(3000), "RUB")

	a, err := NewAppointment(tenantID, uuid.New(), uuid.New(), uuid.New(), serviceInfo, slot, price, BookingSourceOnline, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a.Status() != StatusRequested {
		t.Errorf("expected status 'requested', got %q", a.Status())
	}

	events := a.GetDomainEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].GetName() != "appointment.booked" {
		t.Errorf("expected 'appointment.booked' event, got %q", events[0].GetName())
	}
}

func TestNewAppointment_PastTimeSlot(t *testing.T) {
	tenantID := makeTestTenantID()
	serviceInfo := makeTestServiceInfo()
	pastSlot := MustNewTimeSlot(time.Now().Add(-2*time.Hour), time.Now().Add(-1*time.Hour))
	price := model.MustNewMoney(decimal.NewFromInt(3000), "RUB")

	_, err := NewAppointment(tenantID, uuid.New(), uuid.New(), uuid.New(), serviceInfo, pastSlot, price, BookingSourceOnline, "")
	if err == nil {
		t.Fatal("expected error for past time slot")
	}
}

func TestAppointment_StatusTransitions(t *testing.T) {
	tenantID := makeTestTenantID()
	serviceInfo := makeTestServiceInfo()
	slot := MustNewTimeSlot(futureTime(24), futureTime(25))
	price := model.MustNewMoney(decimal.NewFromInt(3000), "RUB")

	a, _ := NewAppointment(tenantID, uuid.New(), uuid.New(), uuid.New(), serviceInfo, slot, price, BookingSourceOnline, "")

	// Requested → Confirmed
	if err := a.Confirm(); err != nil {
		t.Fatalf("unexpected error on confirm: %v", err)
	}
	if a.Status() != StatusConfirmed {
		t.Errorf("expected status 'confirmed', got %q", a.Status())
	}

	// Confirmed → Cancel should work
	a2, _ := NewAppointment(tenantID, uuid.New(), uuid.New(), uuid.New(), serviceInfo, slot, price, BookingSourceOnline, "")
	_ = a2.Confirm()
	if err := a2.Cancel("changed mind"); err != nil {
		t.Fatalf("unexpected error on cancel: %v", err)
	}
	if a2.Status() != StatusCancelledByClient {
		t.Errorf("expected status 'cancelled_by_client', got %q", a2.Status())
	}
}

func TestAppointment_Complete(t *testing.T) {
	tenantID := makeTestTenantID()
	serviceInfo := makeTestServiceInfo()
	slot := MustNewTimeSlot(futureTime(24), futureTime(25))
	price := model.MustNewMoney(decimal.NewFromInt(3000), "RUB")

	a, _ := NewAppointment(tenantID, uuid.New(), uuid.New(), uuid.New(), serviceInfo, slot, price, BookingSourceOnline, "")
	a.ClearDomainEvents() // clear the booked event

	if err := a.Complete(); err != nil {
		t.Fatalf("unexpected error on complete: %v", err)
	}
	if a.Status() != StatusCompleted {
		t.Errorf("expected status 'completed', got %q", a.Status())
	}

	events := a.GetDomainEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].GetName() != "appointment.completed" {
		t.Errorf("expected 'appointment.completed' event, got %q", events[0].GetName())
	}
}

func TestAppointment_CannotCancelCompleted(t *testing.T) {
	tenantID := makeTestTenantID()
	serviceInfo := makeTestServiceInfo()
	slot := MustNewTimeSlot(futureTime(24), futureTime(25))
	price := model.MustNewMoney(decimal.NewFromInt(3000), "RUB")

	a, _ := NewAppointment(tenantID, uuid.New(), uuid.New(), uuid.New(), serviceInfo, slot, price, BookingSourceOnline, "")
	_ = a.Complete()

	if err := a.Cancel("reason"); err == nil {
		t.Fatal("expected error when cancelling a completed appointment")
	}
}

func TestMasterSchedule_BookSlot(t *testing.T) {
	baseTime := time.Date(2026, 3, 15, 9, 0, 0, 0, time.UTC)
	wh := MustNewWorkingHours(
		baseTime,
		baseTime.Add(9*time.Hour), // 9:00-18:00
		baseTime.Add(4*time.Hour), // break 13:00-14:00
		baseTime.Add(5*time.Hour),
	)

	ms, err := NewMasterSchedule(uuid.New(), uuid.New(), baseTime, wh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slot := MustNewTimeSlot(baseTime.Add(1*time.Hour), baseTime.Add(2*time.Hour)) // 10:00-11:00
	if err := ms.BookSlot(slot); err != nil {
		t.Fatalf("unexpected error booking slot: %v", err)
	}

	// Same slot should not be available
	if ms.IsAvailable(slot) {
		t.Error("expected slot to be unavailable after booking")
	}

	// Overlapping slot should fail
	overlap := MustNewTimeSlot(baseTime.Add(90*time.Minute), baseTime.Add(150*time.Minute)) // 10:30-11:30
	if err := ms.BookSlot(overlap); err == nil {
		t.Error("expected error for overlapping slot")
	}
}

func TestMasterSchedule_BreakSlot(t *testing.T) {
	baseTime := time.Date(2026, 3, 15, 9, 0, 0, 0, time.UTC)
	wh := MustNewWorkingHours(
		baseTime,
		baseTime.Add(9*time.Hour),
		baseTime.Add(4*time.Hour), // break 13:00-14:00
		baseTime.Add(5*time.Hour),
	)

	ms, _ := NewMasterSchedule(uuid.New(), uuid.New(), baseTime, wh)

	// Slot during break should not be available
	breakSlot := MustNewTimeSlot(baseTime.Add(4*time.Hour), baseTime.Add(5*time.Hour)) // 13:00-14:00
	if ms.IsAvailable(breakSlot) {
		t.Error("expected break slot to be unavailable")
	}
}

func TestMasterSchedule_ReleaseSlot(t *testing.T) {
	baseTime := time.Date(2026, 3, 15, 9, 0, 0, 0, time.UTC)
	wh := MustNewWorkingHours(baseTime, baseTime.Add(9*time.Hour), time.Time{}, time.Time{})

	ms, _ := NewMasterSchedule(uuid.New(), uuid.New(), baseTime, wh)
	slot := MustNewTimeSlot(baseTime.Add(1*time.Hour), baseTime.Add(2*time.Hour))

	_ = ms.BookSlot(slot)
	ms.ReleaseSlot(slot)

	if !ms.IsAvailable(slot) {
		t.Error("expected slot to be available after release")
	}
}

func TestMasterSchedule_GetAvailableSlots(t *testing.T) {
	baseTime := time.Date(2026, 3, 15, 9, 0, 0, 0, time.UTC)
	wh := MustNewWorkingHours(baseTime, baseTime.Add(9*time.Hour), time.Time{}, time.Time{})

	ms, _ := NewMasterSchedule(uuid.New(), uuid.New(), baseTime, wh)

	slots := ms.GetAvailableSlots(60 * time.Minute)
	if len(slots) == 0 {
		t.Error("expected available slots")
	}
}
