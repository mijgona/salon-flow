# LLM Prompt: Create a Production-Ready Beauty Salon CRM (Go + DDD)

You are a Principal Go Engineer. Create a **complete, compilable Go project** for a Beauty Salon Client Management & Loyalty CRM following the specifications below. The project must use **Domain-Driven Design**, **Hexagonal Architecture**, and **CQRS**. Follow the reference patterns exactly.

---

## Reference Architecture Patterns

The project follows these Go DDD infrastructure patterns. Use them in every aggregate, entity, value object, and event:

### BaseEntity (generic identity)
```go
package ddd

type BaseEntity[ID comparable] struct { id ID }
func NewBaseEntity[ID comparable](id ID) *BaseEntity[ID] { return &BaseEntity[ID]{id: id} }
func (be *BaseEntity[ID]) ID() ID { return be.id }
func (be *BaseEntity[ID]) Equal(other *BaseEntity[ID]) bool { return other != nil && be.id == other.id }
```

### BaseAggregate (domain event support)
```go
package ddd

type BaseAggregate[ID comparable] struct {
    baseEntity   *BaseEntity[ID]
    domainEvents []DomainEvent
}
func NewBaseAggregate[ID comparable](id ID) *BaseAggregate[ID] {
    return &BaseAggregate[ID]{baseEntity: NewBaseEntity[ID](id), domainEvents: make([]DomainEvent, 0)}
}
func (ba *BaseAggregate[ID]) ID() ID                          { return ba.baseEntity.ID() }
func (ba *BaseAggregate[ID]) RaiseDomainEvent(event DomainEvent) { ba.domainEvents = append(ba.domainEvents, event) }
func (ba *BaseAggregate[ID]) GetDomainEvents() []DomainEvent  { return ba.domainEvents }
func (ba *BaseAggregate[ID]) ClearDomainEvents()              { ba.domainEvents = []DomainEvent{} }
```

### DomainEvent interface
```go
type DomainEvent interface {
    GetID() uuid.UUID
    GetName() string
}
```

### Mediatr (in-process event dispatcher)
```go
type EventHandler interface { Handle(ctx context.Context, event DomainEvent) error }
type Mediatr interface {
    Subscribe(handler EventHandler, events ...DomainEvent)
    Publish(ctx context.Context, event DomainEvent) error
}
```

---

## Project Folder Structure

Create exactly this structure:
```
salon-crm/
├── go.mod                              # module: salon-crm
├── go.sum
├── cmd/
│   ├── app/main.go
│   ├── composition_root.go
│   └── config.go
├── internal/
│   ├── pkg/
│   │   ├── ddd/
│   │   │   ├── entity.go              # BaseEntity[ID]
│   │   │   ├── aggregate.go           # BaseAggregate[ID]
│   │   │   ├── aggregate_root.go      # AggregateRoot type alias
│   │   │   ├── domain_event.go        # DomainEvent interface
│   │   │   └── mediatr.go            # Mediatr implementation
│   │   ├── errs/
│   │   │   ├── value_required.go
│   │   │   └── value_must_be.go
│   │   └── outbox/
│   │       └── event_registry.go
│   ├── core/
│   │   ├── domain/
│   │   │   ├── model/
│   │   │   │   ├── money.go           # Money value object (decimal, currency=RUB)
│   │   │   │   ├── phone_number.go    # PhoneNumber VO with Russian format validation
│   │   │   │   ├── tenant_id.go       # TenantID VO (shared kernel)
│   │   │   │   ├── birthday.go        # Birthday VO
│   │   │   │   ├── discount.go        # Discount VO (percent-based)
│   │   │   │   │
│   │   │   │   ├── client/
│   │   │   │   │   ├── client.go              # Aggregate Root
│   │   │   │   │   ├── contact_info.go        # VO: phone, email, firstName, lastName
│   │   │   │   │   ├── preferences.go         # VO: preferredMasterID, favoriteServices, channel
│   │   │   │   │   ├── allergy.go             # VO: substance, severity
│   │   │   │   │   ├── visit_record.go        # Entity: appointmentID, masterID, service, price, review
│   │   │   │   │   ├── note.go                # VO: text, authorID, createdAt
│   │   │   │   │   ├── photo.go               # VO: url, type, uploadedAt
│   │   │   │   │   ├── client_source.go       # enum: online_booking, admin_entry, referral, walk_in
│   │   │   │   │   ├── client_registered.go   # Domain Event
│   │   │   │   │   └── client_test.go
│   │   │   │   │
│   │   │   │   ├── scheduling/
│   │   │   │   │   ├── appointment.go         # Aggregate Root
│   │   │   │   │   ├── master_schedule.go     # Aggregate Root
│   │   │   │   │   ├── time_slot.go           # VO: startTime, endTime, OverlapsWith()
│   │   │   │   │   ├── service_info.go        # VO: serviceID, name, duration, basePrice
│   │   │   │   │   ├── working_hours.go       # VO: startTime, endTime, breakStart, breakEnd
│   │   │   │   │   ├── booking_source.go      # enum: online, admin
│   │   │   │   │   ├── appointment_status.go  # enum: Requested,Confirmed,InProgress,Completed,CancelledByClient,CancelledBySalon,NoShow
│   │   │   │   │   ├── appointment_booked.go  # Domain Event
│   │   │   │   │   ├── appointment_completed.go
│   │   │   │   │   ├── appointment_cancelled.go
│   │   │   │   │   └── appointment_test.go
│   │   │   │   │
│   │   │   │   ├── loyalty/
│   │   │   │   │   ├── loyalty_account.go     # Aggregate Root
│   │   │   │   │   ├── points.go              # VO: int value, Add(), Subtract()
│   │   │   │   │   ├── tier.go                # VO enum: Bronze, Silver, Gold, VIP
│   │   │   │   │   ├── tier_threshold.go      # VO: tier + minLifetimePoints
│   │   │   │   │   ├── points_transaction.go  # Entity: amount, type, reason, relatedEntityID
│   │   │   │   │   ├── referral.go            # Entity: referredClientID, status, bonusEarned
│   │   │   │   │   ├── points_earned.go       # Domain Event
│   │   │   │   │   ├── tier_changed.go        # Domain Event
│   │   │   │   │   └── loyalty_test.go
│   │   │   │   │
│   │   │   │   └── certificate/
│   │   │   │       ├── certificate.go         # Aggregate Root: balance, expiresAt, status
│   │   │   │       ├── certificate_activated.go
│   │   │   │       └── certificate_test.go
│   │   │   │
│   │   │   └── services/
│   │   │       ├── loyalty_policy.go          # Domain Service
│   │   │       └── availability_service.go    # Domain Service
│   │   │
│   │   ├── application/
│   │   │   ├── commands/
│   │   │   │   ├── register_client.go
│   │   │   │   ├── update_client_profile.go
│   │   │   │   ├── book_appointment.go
│   │   │   │   ├── cancel_appointment.go
│   │   │   │   ├── complete_appointment.go
│   │   │   │   ├── earn_points.go
│   │   │   │   └── activate_certificate.go
│   │   │   ├── queries/
│   │   │   │   ├── get_client.go
│   │   │   │   ├── get_client_history.go
│   │   │   │   ├── get_available_slots.go
│   │   │   │   └── get_loyalty_account.go
│   │   │   └── eventhandlers/
│   │   │       ├── accrue_points_on_completed.go
│   │   │       ├── add_visit_record_on_completed.go
│   │   │       ├── create_loyalty_on_registered.go
│   │   │       └── schedule_reminders_on_booked.go
│   │   │
│   │   └── ports/
│   │       ├── client_repository.go
│   │       ├── appointment_repository.go
│   │       ├── master_schedule_repository.go
│   │       ├── loyalty_repository.go
│   │       ├── certificate_repository.go
│   │       ├── notification_sender.go
│   │       ├── payment_client.go
│   │       ├── service_catalog_client.go
│   │       ├── outbox_repository.go
│   │       └── tx_manager.go
│   │
│   ├── adapters/
│   │   ├── in/http/
│   │   │   ├── client_handler.go
│   │   │   ├── appointment_handler.go
│   │   │   └── loyalty_handler.go
│   │   └── out/
│   │       ├── postgres/
│   │       │   ├── clientrepo/repository.go
│   │       │   ├── appointmentrepo/repository.go
│   │       │   ├── loyaltyrepo/repository.go
│   │       │   ├── schedulerepo/repository.go
│   │       │   └── tx_manager.go
│   │       └── inmemory/
│   │           ├── client_repository.go
│   │           └── appointment_repository.go
│   │
│   └── jobs/
│       └── outbox_job.go
├── migrations/
│   ├── 001_clients.sql
│   ├── 002_appointments.sql
│   ├── 003_loyalty.sql
│   └── 004_certificates.sql
├── configs/config.yaml
├── Dockerfile
└── makefile
```

---

## Bounded Contexts & Subdomains

| Subdomain | Type | Bounded Context |
|-----------|------|-----------------|
| Client Management | **Core** | Client Context — profiles, contacts, preferences, allergies, photos, notes, visit history |
| Scheduling & Appointments | **Core** | Scheduling Context — booking, master schedules, time slots, availability, service catalog |
| Loyalty & Rewards | **Core** | Loyalty Context — points, tiers (Bronze/Silver/Gold/VIP), referrals, personal discounts |
| Subscriptions & Certificates | Supporting | Certificate Context — gift cards, subscriptions, activation, balance, expiration |
| Notifications | Supporting | Notification Context — SMS/WhatsApp/Email reminders, birthday, promos |
| Marketing | Supporting | Marketing Context — RFM segmentation, campaigns |
| Payments | Generic | Payment Context — ACL to Yandex.Kassa/Tinkoff/SBP |
| Analytics | Generic | Analytics Context — LTV, retention, avg check |
| Identity & Tenancy | Generic | Tenant Context — auth, multi-tenant (row-level via `tenant_id`) |

### Context Map Relationships
- **Scheduling → Client**: Customer-Supplier (AppointmentCompleted → VisitRecord)
- **Scheduling → Loyalty**: Customer-Supplier (AppointmentCompleted → EarnPoints)
- **Scheduling → Notification**: Published Language (AppointmentBooked → Reminders)
- **Client → Marketing**: Open Host Service (Client API for RFM queries)
- **Scheduling → Payment**: Anti-Corruption Layer (payment request abstraction)
- **Tenant → All**: Shared Kernel (TenantID value object)
- **Loyalty → Client**: Conformist (reads client data, conforms to Client model)

---

## Aggregate Designs

### Client Aggregate Root
- **Root**: `Client` — UUID id, TenantID, ContactInfo, Birthday, Preferences, []Allergy, []Note, []Photo, []VisitRecord, ClientSource, registeredAt
- **Entities**: `VisitRecord` (appointmentID, masterID, service, price, discount, paymentStatus, rating, review, visitedAt)
- **Value Objects**: ContactInfo(phone, email, firstName, lastName), Preferences(preferredMasterID, favoriteServices, channel), Allergy(substance, severity), Note(text, authorID, createdAt), Photo(url, type, uploadedAt)
- **Events**: `ClientRegistered` → triggers loyalty account creation + welcome notification
- **Invariants**: Phone required, valid format. Allergy deduplication by substance. Status guards.
- **Methods**: `NewClient()`, `UpdateProfile()`, `AddAllergy()`, `AddVisitRecord()`, `AddNote()`, `TotalVisits()`, `TotalSpent()`

### Appointment Aggregate Root
- **Root**: `Appointment` — UUID id, TenantID, clientID, masterID, salonID, ServiceInfo, TimeSlot, status, price, BookingSource, comment
- **Value Objects**: TimeSlot(startTime, endTime, Duration(), OverlapsWith()), ServiceInfo(serviceID, name, duration, basePrice), AppointmentStatus(enum)
- **Events**: `AppointmentBooked`, `AppointmentCompleted`, `AppointmentCancelledByClient`
- **Invariants**: Cannot book in the past. Must check master availability. Status transitions: Requested→Confirmed→InProgress→Completed. Cancel only if not InProgress/Completed.
- **Methods**: `NewAppointment()`, `Confirm()`, `Cancel(reason)`, `Reschedule(newSlot)`, `Complete()`, `NoShow()`

### MasterSchedule Aggregate Root
- **Root**: `MasterSchedule` — UUID id, masterID, salonID, date, WorkingHours, []bookedSlots, []blockedSlots
- **Value Objects**: WorkingHours(startTime, endTime, breakStart, breakEnd)
- **Invariants**: Slots cannot overlap. Must be within working hours. Not during break.
- **Methods**: `IsAvailable(timeSlot)`, `BookSlot(timeSlot)`, `ReleaseSlot(timeSlot)`, `GetAvailableSlots(duration)`

### LoyaltyAccount Aggregate Root
- **Root**: `LoyaltyAccount` — UUID id, clientID, TenantID, tier, Points balance, Points lifetimePoints, []PointsTransaction, []Referral
- **Entities**: PointsTransaction(id, amount, type, reason, relatedEntityID, createdAt), Referral(id, referredClientID, status, bonusEarned, createdAt)
- **Value Objects**: Points(int value, Add, Subtract, IsZero), LoyaltyTier(enum: Bronze/Silver/Gold/VIP, DiscountPercent(), PointsMultiplier()), TierThreshold(tier, minPoints)
- **Events**: `LoyaltyPointsEarned`, `ClientTierChanged`
- **Invariants**: Cannot redeem more points than balance. Tier only changes upward. Referral bonus once per referred client.
- **Methods**: `EarnPoints(amount, reason)`, `RedeemPoints(amount)`, `RecalculateTier()`, `AddReferral()`, `GetPersonalDiscount()`

### Tier Thresholds
| Tier | Min Lifetime Points | Discount % | Points Multiplier |
|------|---------------------|------------|-------------------|
| Bronze | 0 | 0% | 1.0x |
| Silver | 5,000 | 5% | 1.2x |
| Gold | 15,000 | 10% | 1.5x |
| VIP | 50,000 | 15% | 2.0x |

---

## Domain Events (7 key events)

### 1. ClientRegistered
- **Fields**: eventId, clientId, tenantId, firstName, lastName, phone, source, referredByClientId
- **Consumers**: Loyalty (create account), Notification (welcome), Marketing (add to segment)

### 2. AppointmentBooked
- **Fields**: eventId, appointmentId, clientId, masterId, salonId, serviceId, serviceName, startTime, endTime, price, source
- **Consumers**: Notification (schedule 24h + 2h reminders)

### 3. AppointmentCompleted
- **Fields**: eventId, appointmentId, clientId, masterId, salonId, serviceName, finalPrice, discount, paymentMethod
- **Consumers**: Loyalty (accrue points), Client (add visit record), Notification (request review)

### 4. LoyaltyPointsEarned
- **Fields**: eventId, loyaltyAccountId, clientId, pointsEarned, multiplier, reason, relatedEntityId, newBalance, lifetimePoints
- **Consumers**: Tier recalculation check

### 5. ClientTierChanged
- **Fields**: eventId, loyaltyAccountId, clientId, previousTier, newTier, lifetimePoints, newDiscountPercent
- **Consumers**: Notification (congratulations), Client profile update

### 6. CertificateActivated
- **Fields**: eventId, certificateId, activatedByClientId, purchasedByClientId, balance, expiresAt
- **Consumers**: Payment (available balance)

### 7. AppointmentCancelledByClient
- **Fields**: eventId, appointmentId, clientId, masterId, salonId, originalStartTime, cancelledAt, reason
- **Consumers**: Schedule (release slot), Notification (notify master)

---

## Domain Services

### LoyaltyPolicy
```
Interface: LoyaltyPolicy
- CalculatePointsForVisit(amount Money, tier LoyaltyTier) → Points   // 1pt per 10 RUB × tier multiplier
- DetermineNewTier(lifetimePoints Points) → LoyaltyTier              // highest tier where threshold ≤ points
- GetReferralBonus() → Points                                         // 500 points
- GetPersonalDiscount(tier LoyaltyTier) → Discount                   // tier-based %
```

### AvailabilityService
```
Interface: AvailabilityService
- GetAvailableSlots(masterID, salonID, date, serviceDuration) → []TimeSlot
- IsSlotAvailable(masterID, date, timeSlot) → bool
```

---

## Repository Interfaces

```
ClientRepository:
  Add(ctx, tx, *Client) error
  Update(ctx, tx, *Client) error
  Get(ctx, tx, id UUID) (*Client, error)
  FindByPhone(ctx, tx, tenantID, phone) (*Client, error)
  FindByTenant(ctx, tx, tenantID, limit, offset) ([]*Client, error)

AppointmentRepository:
  Add(ctx, tx, *Appointment) error
  Update(ctx, tx, *Appointment) error
  Get(ctx, tx, id UUID) (*Appointment, error)
  FindByClientID(ctx, tx, clientID) ([]*Appointment, error)
  FindByMasterAndDate(ctx, tx, masterID, date) ([]*Appointment, error)

MasterScheduleRepository:
  Add(ctx, tx, *MasterSchedule) error
  Update(ctx, tx, *MasterSchedule) error
  GetByMasterAndDate(ctx, tx, masterID, date) (*MasterSchedule, error)

LoyaltyRepository:
  Add(ctx, tx, *LoyaltyAccount) error
  Update(ctx, tx, *LoyaltyAccount) error
  GetByClientID(ctx, tx, clientID) (*LoyaltyAccount, error)

TxManager:
  Execute(ctx, func(tx Tx) error) error
```

---

## Command Handlers (Use Cases)

### BookAppointmentCommandHandler
1. Validate startTime is in the future
2. Get service details (duration, price) from ServiceCatalog
3. Build TimeSlot from startTime + duration
4. Inside transaction:
   a. Load MasterSchedule for master+date (with lock)
   b. Check `schedule.IsAvailable(timeSlot)` → error if not
   c. `schedule.BookSlot(timeSlot)` to reserve
   d. Create `NewAppointment(...)` aggregate
   e. Persist appointment + updated schedule
5. Outbox publishes `AppointmentBooked` domain event

### RegisterClientCommandHandler
1. Check no existing client with same phone+tenantID
2. Create `NewClient(tenantID, contactInfo, source)`
3. Persist → raises `ClientRegistered` event
4. Event handler creates LoyaltyAccount + sends welcome notification

### CompleteAppointmentCommandHandler
1. Load appointment by ID
2. Call `appointment.Complete()`
3. Persist → raises `AppointmentCompleted`
4. EventHandler: accrue loyalty points, add visit record to client

---

## Event Handlers (Cross-Context Integration)

### AccruePointsOnAppointmentCompleted
- Subscribes to: `AppointmentCompleted`
- Loads LoyaltyAccount by clientID
- Calls `loyaltyPolicy.CalculatePointsForVisit(finalPrice, account.Tier())`
- Calls `account.EarnPoints(points, "appointment_completed")`
- Calls `account.RecalculateTier()` using `loyaltyPolicy.DetermineNewTier()`
- Persists LoyaltyAccount

### CreateLoyaltyOnClientRegistered
- Subscribes to: `ClientRegistered`
- Creates new `LoyaltyAccount` with Bronze tier, 0 points
- If referredByClientId present: adds Referral, earns bonus for both

### AddVisitRecordOnAppointmentCompleted
- Subscribes to: `AppointmentCompleted`
- Creates VisitRecord from event data
- Loads Client, calls `client.AddVisitRecord(record)`, persists

### ScheduleRemindersOnAppointmentBooked
- Subscribes to: `AppointmentBooked`
- Schedules notification: 24h before startTime
- Schedules notification: 2h before startTime

---

## Architecture Decisions

- **CQRS**: Yes. Separate command/query handlers. Write side = rich domain model; read side = direct DB queries.
- **Event Sourcing**: No. PostgreSQL CRUD with domain events via Outbox pattern.
- **Database**: PostgreSQL. JSONB for preferences, notes. Row-level multi-tenancy via `tenant_id`.
- **Event Bus**: Outbox table → cron job → Mediatr (in-process). Kafka for cross-service if needed later.
- **API**: REST (oapi-codegen from OpenAPI specs). WebSocket for real-time calendar updates.
- **Multi-tenant**: `tenant_id` column on every table. TenantID value object shared across all contexts (Shared Kernel).
- **Dependencies**: github.com/google/uuid, github.com/shopspring/decimal, gorm.io/gorm, github.com/labstack/echo/v4.

---

## Migrations (PostgreSQL)

### 001_clients.sql
```sql
CREATE TABLE clients (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    birthday DATE,
    preferences JSONB DEFAULT '{}',
    allergies JSONB DEFAULT '[]',
    notes JSONB DEFAULT '[]',
    source VARCHAR(50) NOT NULL,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, phone)
);
CREATE INDEX idx_clients_tenant ON clients(tenant_id);
CREATE INDEX idx_clients_phone ON clients(tenant_id, phone);
```

### 002_appointments.sql
```sql
CREATE TABLE appointments (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    client_id UUID NOT NULL REFERENCES clients(id),
    master_id UUID NOT NULL,
    salon_id UUID NOT NULL,
    service_id UUID NOT NULL,
    service_name VARCHAR(200) NOT NULL,
    service_duration INTERVAL NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'requested',
    price_amount DECIMAL(12,2) NOT NULL,
    price_currency VARCHAR(3) DEFAULT 'RUB',
    source VARCHAR(20) NOT NULL,
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE master_schedules (
    id UUID PRIMARY KEY,
    master_id UUID NOT NULL,
    salon_id UUID NOT NULL,
    schedule_date DATE NOT NULL,
    work_start TIME NOT NULL,
    work_end TIME NOT NULL,
    break_start TIME,
    break_end TIME,
    booked_slots JSONB DEFAULT '[]',
    blocked_slots JSONB DEFAULT '[]',
    UNIQUE(master_id, schedule_date)
);
```

### 003_loyalty.sql
```sql
CREATE TABLE loyalty_accounts (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL UNIQUE REFERENCES clients(id),
    tenant_id UUID NOT NULL,
    tier VARCHAR(10) NOT NULL DEFAULT 'Bronze',
    balance INT NOT NULL DEFAULT 0,
    lifetime_points INT NOT NULL DEFAULT 0
);
CREATE TABLE points_transactions (
    id UUID PRIMARY KEY,
    loyalty_account_id UUID NOT NULL REFERENCES loyalty_accounts(id),
    amount INT NOT NULL,
    type VARCHAR(20) NOT NULL,
    reason VARCHAR(100),
    related_entity_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE referrals (
    id UUID PRIMARY KEY,
    loyalty_account_id UUID NOT NULL REFERENCES loyalty_accounts(id),
    referred_client_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    bonus_earned INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 004_certificates.sql
```sql
CREATE TABLE certificates (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    purchased_by UUID REFERENCES clients(id),
    activated_by UUID REFERENCES clients(id),
    balance_amount DECIMAL(12,2) NOT NULL,
    balance_currency VARCHAR(3) DEFAULT 'RUB',
    status VARCHAR(20) NOT NULL DEFAULT 'created',
    activated_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### outbox.sql
```sql
CREATE TABLE outbox (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);
```

---

## Instructions

1. Create all files following the folder structure exactly.
2. Every value object must validate in its constructor (`New*` returns error) and have a `Must*` panic variant.
3. Every aggregate root embeds `*ddd.BaseAggregate[uuid.UUID]` and delegates ID/events to it.
4. Every aggregate has a `Restore*` function for rehydration from the database (no validation, no events).
5. All struct fields are **unexported**; provide getter methods.
6. Business rules are enforced in aggregate methods, not in application/adapter layer.
7. Write unit tests for every aggregate covering: happy path, invariant violations, event raising.
8. Repository interfaces go in `ports/`. Implementations go in `adapters/out/`.
9. Command handlers use `TxManager.Execute()` for transactional boundaries.
10. The Outbox pattern: TxManager saves domain events to outbox table within the same transaction. A cron job reads pending events and publishes them via Mediatr.
11. In-memory repository implementations for testing.
12. `composition_root.go` wires everything together using lazy initialization with `sync.Once`.
