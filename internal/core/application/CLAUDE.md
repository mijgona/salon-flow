# Application Layer Rules

Orchestrates domain objects. Contains use cases (commands + queries + event handlers).
**No business logic here — delegate everything to aggregates and domain services.**

## Command Handler Pattern

```go
type [Feature]Handler struct {
    repo1 ports.XRepository
    repo2 ports.YRepository
    tx    ports.TxManager
}

func (h *[Feature]Handler) Handle(ctx context.Context, cmd [Feature]Command) error {
    // 1. Validate inputs that don't require DB (cheap, fast)
    if cmd.StartTime.Before(time.Now().UTC()) {
        return fmt.Errorf("[feature]: %w", ErrStartTimeInPast)
    }

    // 2. All writes inside a transaction
    return h.tx.Execute(ctx, func(tx ports.Tx) error {
        // 3. Load aggregate(s)
        entity, err := h.repo1.Get(ctx, tx, cmd.ID)
        if err != nil {
            return fmt.Errorf("[feature]: load [entity]: %w", err)
        }

        // 4. Call aggregate method (business logic lives there)
        if err := entity.DoSomething(cmd.Arg); err != nil {
            return err
        }

        // 5. Persist (TxManager saves domain events to outbox in same TX)
        return h.repo1.Update(ctx, tx, entity)
    })
}
```

## Query Handler Pattern

```go
// Query handlers return plain DTOs — never domain objects
type [Feature]DTO struct {
    ID        uuid.UUID `json:"id"`
    // exported fields only
}

func (h *[Feature]QueryHandler) Handle(ctx context.Context, q [Feature]Query) (*[Feature]DTO, error) {
    // Direct DB query — no aggregate loading needed
    // Return DTO, not domain object
}
```

## Event Handler Pattern

```go
// One concern per handler — subscribe to one event
type AccruePointsOnCompleted struct {
    loyalty ports.LoyaltyRepository
    policy  services.LoyaltyPolicy
    tx      ports.TxManager
}

func (h *AccruePointsOnCompleted) Handle(ctx context.Context, event ddd.DomainEvent) error {
    completed, ok := event.(scheduling.AppointmentCompleted)
    if !ok {
        return nil
    }
    return h.tx.Execute(ctx, func(tx ports.Tx) error {
        // load → mutate via domain service → persist
    })
}
```

## Rules

- **Never** import from `adapters/` — dependency flows inward only
- **Never** return domain objects from query handlers — use DTOs
- **Always** use `fmt.Errorf("handler_name: %w", err)` for error wrapping
- **Always** put all writes in `TxManager.Execute()` — no partial updates
- Register event handlers in `cmd/composition_root.go` via `mediatr.Subscribe()`
