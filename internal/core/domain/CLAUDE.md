# Domain Layer Rules

This is the heart of salon-flow. **No imports from application/ or adapters/ are allowed here.**

## Aggregate Checklist

Every aggregate root **must**:
- [ ] Embed `*ddd.BaseAggregate[uuid.UUID]` (not `*ddd.BaseEntity`)
- [ ] Have all fields **unexported** with getter methods
- [ ] Have `New*()` — validates invariants, raises domain event, returns `(*T, error)`
- [ ] Have `Must*()` — calls `New*()`, panics on error (for tests only)
- [ ] Have `Restore*()` — rehydrates from DB, no validation, no events
- [ ] Enforce every invariant inside the aggregate method, never in handlers

## Value Object Checklist

Every value object **must**:
- [ ] Have `New*()` returning `(VO, error)` — validate in constructor
- [ ] Have `Must*()` panic variant
- [ ] Have unexported `value` field with getter
- [ ] Have `Equal()` for comparison
- [ ] Be immutable — no setters

## Domain Event Checklist

Every domain event **must**:
- [ ] Implement `ddd.DomainEvent` interface: `GetID() uuid.UUID`, `GetName() string`
- [ ] Have name format: `"[package].[EventName]"` e.g. `"scheduling.AppointmentBooked"`
- [ ] Be raised inside the aggregate method that changes state (not in handler)
- [ ] Contain all data consumers need (denormalized — no DB lookups in handlers)

## Invariants Reference

| Aggregate | Key Invariants |
|-----------|---------------|
| Client | Phone required + unique per tenant; allergy deduplication by substance |
| Appointment | Cannot book in past; status transitions: Requested→Confirmed→InProgress→Completed; cancel only if not InProgress/Completed |
| MasterSchedule | Slots cannot overlap; must be within working hours; not during break |
| LoyaltyAccount | Cannot redeem more than balance; tier only changes upward; referral bonus once per referrer |
| Certificate | Cannot activate expired cert; balance cannot go negative |

## Error Messages

Use descriptive errors that include the aggregate name:
```go
// Good
fmt.Errorf("appointment: cannot cancel appointment in status %s", a.status)
fmt.Errorf("time_slot: end time %v must be after start time %v", end, start)

// Bad
fmt.Errorf("invalid status")
errors.New("error")
```