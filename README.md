# salon-flow

Beauty Salon Client Management & Loyalty CRM.

**Go 1.25 · DDD · Hexagonal Architecture · CQRS · Outbox Pattern**

---

## Overview

salon-flow is a backend CRM for beauty salons. It manages clients, appointments, loyalty accounts, and gift certificates with a clean Domain-Driven Design architecture.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.25 |
| HTTP | Echo v4 |
| Database | PostgreSQL via pgx/v5 + pgxpool |
| Query builder | Squirrel |
| IDs | google/uuid |
| Money | shopspring/decimal |
| Multi-tenancy | `tenant_id` column on every table (row-level isolation) |

## Architecture

```
cmd/                          # entry point + DI wiring (composition_root.go)
internal/
  pkg/ddd/                    # BaseEntity, BaseAggregate, DomainEvent, Mediatr
  pkg/errs/                   # sentinel errors (ValueRequired, ValueMustBe)
  core/
    domain/model/             # Aggregates + Value Objects + Domain Events
      client/
      scheduling/
      loyalty/
      certificate/
    domain/services/          # LoyaltyPolicy, AvailabilityService
    application/
      commands/               # write side: command structs + handlers
      queries/                # read side: query structs + handlers (DTOs only)
      eventhandlers/          # cross-context: subscribe to domain events
    ports/                    # interfaces: repositories, TxManager, etc.
  adapters/
    in/http/                  # Echo handlers
    out/postgres/             # pgx + squirrel repository implementations
    out/inmemory/             # in-memory repos for tests
  jobs/
    outbox_job.go             # polls outbox table, publishes events via Mediatr
migrations/                   # 001_clients.sql ... 005_calendar_indexes.sql
configs/config.yaml
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/clients` | Register a new client |
| GET | `/api/v1/clients/:id` | Get client profile |
| GET | `/api/v1/clients/:id/history` | Get client visit history |
| POST | `/api/v1/appointments` | Book an appointment |
| POST | `/api/v1/appointments/:id/cancel` | Cancel an appointment |
| POST | `/api/v1/appointments/:id/complete` | Complete an appointment |
| GET | `/api/v1/appointments/available-slots` | Get available time slots |
| GET | `/api/v1/loyalty/:client_id` | Get loyalty account |
| GET | `/api/v1/calendar` | Get calendar view |
| GET | `/health` | Health check |

## Domain Events

Events flow: Aggregate raises event → TxManager saves to `outbox` table (same TX) → OutboxJob polls → Mediatr dispatches to handlers.

| Event | Raised by | Handled by |
|-------|-----------|------------|
| `ClientRegistered` | `Client.New*` | CreateLoyaltyOnRegistered |
| `AppointmentBooked` | `Appointment.New*` | ScheduleRemindersOnBooked |
| `AppointmentCompleted` | `Appointment.Complete()` | AccruePointsOnCompleted, AddVisitRecordOnCompleted |
| `AppointmentCancelledByClient` | `Appointment.Cancel()` | release slot, notify master |
| `LoyaltyPointsEarned` | `LoyaltyAccount.EarnPoints()` | tier recalculation |
| `ClientTierChanged` | `LoyaltyAccount.RecalculateTier()` | congratulation notification |
| `CertificateActivated` | `Certificate.Activate()` | payment balance update |

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL 14+

### Configuration

Set the `DATABASE_URL` environment variable or edit `configs/config.yaml`:

```yaml
server:
  port: 8080

database:
  host: localhost
  port: 5432
  user: salon_crm
  password: salon_crm_pass
  name: salon_crm_db
  sslmode: disable

outbox:
  interval: 5s
  batch_size: 100
```

Default DSN (used if `DATABASE_URL` is not set):
```
postgres://salon_crm:salon_crm_pass@localhost:5432/salon_crm_db?sslmode=disable
```

### Run Migrations

```bash
psql $DATABASE_URL -f migrations/001_clients.sql
psql $DATABASE_URL -f migrations/002_appointments.sql
psql $DATABASE_URL -f migrations/003_loyalty.sql
psql $DATABASE_URL -f migrations/004_certificates.sql
psql $DATABASE_URL -f migrations/005_calendar_indexes.sql
```

### Run

```bash
go run ./cmd/app
```

Or with a custom port:

```bash
PORT=9090 go run ./cmd/app
```

## Development

```bash
go build -o bin/salon-crm ./cmd/app   # compile
go run ./cmd/app                       # run
go test ./...                          # all tests
go test ./internal/core/domain/... -v # domain tests only
go vet ./...                           # vet
```

## Testing

Unit tests live next to source files (`*_test.go` in the same package). Domain aggregates are tested with in-memory repositories — no real DB required.

```bash
go test ./internal/core/domain/... -v
```