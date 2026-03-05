# Salon-Flow CRM — Project Guide

Beauty Salon Client Management & Loyalty CRM.
Go 1.25 · DDD · Hexagonal Architecture · CQRS · Outbox Pattern

---

## Architecture

```
cmd/                          ← entry point + composition_root.go (DI wiring)
internal/
  pkg/ddd/                    ← BaseEntity[ID], BaseAggregate[ID], DomainEvent, Mediatr
  pkg/errs/                   ← ValueRequired, ValueMustBe sentinel errors
  core/
    domain/model/             ← Aggregates, Value Objects, Domain Events
      client/                 ← Client AR · ContactInfo · VisitRecord · ClientRegistered
      scheduling/             ← Appointment AR · MasterSchedule AR · TimeSlot · AppointmentBooked/Completed/Cancelled
      loyalty/                ← LoyaltyAccount AR · Points · Tier · PointsEarned · TierChanged
      certificate/            ← Certificate AR · CertificateActivated
    domain/services/          ← LoyaltyPolicy · AvailabilityService (domain services)
    application/
      commands/               ← Write side: command structs + handlers (use TxManager)
      queries/                ← Read side: query structs + handlers (return DTOs, no aggregates)
      eventhandlers/          ← Cross-context: subscribe to domain events, coordinate aggregates
    ports/                    ← Interfaces: repositories, TxManager, NotificationSender, etc.
  adapters/
    in/http/                  ← Echo handlers (parse → command/query → JSON response)
    out/postgres/             ← pgx + squirrel repository implementations
    out/inmemory/             ← in-memory repos for tests
  jobs/
    outbox_job.go             ← polls outbox table, publishes events via Mediatr
migrations/                   ← 001_clients.sql … 005_calendar_indexes.sql
configs/config.yaml
```

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.25 |
| HTTP | Echo v4 (`github.com/labstack/echo/v4`) |
| Database | PostgreSQL via `pgx/v5` + `pgxpool` |
| Query builder | Squirrel (`github.com/Masterminds/squirrel`) |
| IDs | `github.com/google/uuid` |
| Money | `github.com/shopspring/decimal` |
| Multi-tenancy | `tenant_id UUID` column on every table (row-level) |

## Domain Events (Outbox Pattern)

Events flow: Aggregate raises event → TxManager saves to `outbox` table (same TX) → OutboxJob polls → Mediatr dispatches to handlers.

| Event | Raised by | Handled by |
|-------|-----------|------------|
| `ClientRegistered` | `Client.New*` | CreateLoyaltyOnRegistered, welcome notification |
| `AppointmentBooked` | `Appointment.New*` | ScheduleRemindersOnBooked |
| `AppointmentCompleted` | `Appointment.Complete()` | AccruePointsOnCompleted, AddVisitRecordOnCompleted |
| `AppointmentCancelledByClient` | `Appointment.Cancel()` | release slot, notify master |
| `LoyaltyPointsEarned` | `LoyaltyAccount.EarnPoints()` | tier recalculation |
| `ClientTierChanged` | `LoyaltyAccount.RecalculateTier()` | congratulation notification |
| `CertificateActivated` | `Certificate.Activate()` | payment balance update |

## Coding Rules

### Domain Layer (`internal/core/domain/`)
- Every aggregate embeds `*ddd.BaseAggregate[uuid.UUID]`
- All struct fields **unexported** — access via getter methods only
- `New*()` — constructor with validation, returns `error`; `Must*()` — panic variant for tests
- `Restore*()` — rehydration from DB: no validation, no events raised
- Business invariants enforced inside aggregate methods, never in handlers
- Domain events raised via `RaiseDomainEvent(event)` inside aggregate methods

### Application Layer (`internal/core/application/`)
- Command handlers: use `TxManager.Execute()` for all writes
- Query handlers: return plain DTOs (structs with exported fields), never domain objects
- Event handlers: subscribe via `mediatr.Subscribe()`, keep focused on one concern

### Adapter Layer (`internal/adapters/`)
- HTTP handlers: bind request → call handler → return JSON; no business logic
- Postgres repos: always include `tenant_id` filter in every query
- Use `squirrel` for query building — no raw string concatenation

### General
- `go vet ./...` and `go test ./...` must pass before any PR
- Unit tests live next to source (`*_test.go` in same package)
- Test aggregates with `inmemory` repos, not real DB

## API Endpoints

| Method | Path | Handler |
|--------|------|---------|
| POST | `/api/v1/clients` | RegisterClient |
| GET | `/api/v1/clients/:id` | GetClient |
| GET | `/api/v1/clients/:id/history` | GetClientHistory |
| POST | `/api/v1/appointments` | BookAppointment |
| POST | `/api/v1/appointments/:id/cancel` | CancelAppointment |
| POST | `/api/v1/appointments/:id/complete` | CompleteAppointment |
| GET | `/api/v1/appointments/available-slots` | GetAvailableSlots |
| GET | `/api/v1/loyalty/:client_id` | GetLoyaltyAccount |
| GET | `/api/v1/calendar` | GetCalendar |

## Development Commands

```bash
make build          # compile to bin/salon-crm
make run            # go run ./cmd/app
make test           # go test ./... -v
make test-domain    # go test ./internal/core/domain/... -v
make vet            # go vet ./...
make lint           # golangci-lint run
```

## Team Agents

Use slash commands to invoke the right agent for the task:

| Command | Agent | Use when |
|---------|-------|----------|
| `/feature` | @coder + @tester | Implementing a new feature end-to-end |
| `/plan` | @planner | Turning an idea into a User Story + AC |
| `/bug` | @tester + @coder | Investigating and fixing a bug |
| `/review` | @coder | Code review of a PR or file |
| `/test` | @tester | Writing test plan + test cases |
| `/new-aggregate` | @coder | Scaffolding a new DDD aggregate |