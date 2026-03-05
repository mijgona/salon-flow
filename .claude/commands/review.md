# /review — Code Review (@coder)

Structured code review for Go/DDD code in salon-flow.

---

## How to use

Provide one of:
- A file path: `/review internal/core/domain/model/scheduling/appointment.go`
- A feature name: `/review BookAppointment feature`
- "current changes": `/review` (reviews all modified files)

---

## Review Checklist

### Domain Layer (`internal/core/domain/`)

- [ ] Aggregate embeds `*ddd.BaseAggregate[uuid.UUID]` (not `BaseEntity`)
- [ ] All struct fields unexported — getters exist for each
- [ ] `New*()` validates and returns `error`; `Must*()` panic variant exists
- [ ] `Restore*()` exists for DB rehydration — no validation, no events
- [ ] Business invariants enforced inside aggregate, not in handler
- [ ] Domain events raised with `RaiseDomainEvent()` inside the method that changes state
- [ ] No imports from `application/` or `adapters/` (dependency rule)

### Application Layer (`internal/core/application/`)

- [ ] Command handlers use `TxManager.Execute()` — no writes outside tx
- [ ] Query handlers return plain DTOs — no domain objects leaked to caller
- [ ] Event handlers are focused — one concern per handler
- [ ] Errors wrapped with `fmt.Errorf("context: %w", err)`
- [ ] No business logic in handlers — delegated to aggregates

### Adapter Layer (`internal/adapters/`)

- [ ] HTTP handlers: bind → call handler → return JSON only
- [ ] All postgres queries include `tenant_id` filter
- [ ] No raw SQL string concat — squirrel used for query building
- [ ] Errors mapped to correct HTTP status codes (400 / 422 / 404 / 500)

### Tests

- [ ] Happy path covered
- [ ] Invariant violations tested (negative cases)
- [ ] Domain events verified (event is raised, has correct data)
- [ ] Concurrency tested if shared resource (booking, slot)
- [ ] inmemory repos used (not real DB)

---

## Review Output Format

```
## Code Review: [file or feature]

### ✅ Strengths
- [what's done well]

### ⚠️ Issues

| # | File:Line | Severity | Issue | Suggestion |
|---|-----------|----------|-------|------------|
| 1 | appointment.go:45 | Critical | Business logic in HTTP handler | Move to Appointment.Cancel() |
| 2 | repository.go:78 | High | Missing tenant_id filter | Add WHERE tenant_id = $n |
| 3 | client.go:32 | Minor | Exported field `Name` | Make unexported, add Name() getter |

### 🔄 Required before merge
- [ ] Fix Critical issues
- [ ] Fix High issues
- [ ] Re-run: make vet && make test-domain

### 💡 Optional improvements
- [nice-to-haves, not blocking]
```
