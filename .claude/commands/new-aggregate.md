# /new-aggregate — Scaffold a New DDD Aggregate (@coder)

Scaffold a complete, compilable DDD aggregate following salon-flow patterns.

---

## How to use

`/new-aggregate [name] in [bounded-context]`

Examples:
- `/new-aggregate Subscription in certificate`
- `/new-aggregate MasterProfile in scheduling`

---

## Step 1: Gather Info

Ask the user:
1. **Aggregate name**: (e.g. `Subscription`)
2. **Bounded context / package**: (e.g. `certificate`, `scheduling`, `loyalty`)
3. **Key fields**: (list the main data this aggregate owns)
4. **Key methods**: (what actions can be performed? e.g. `Activate`, `Cancel`, `Renew`)
5. **Domain events it raises**: (e.g. `SubscriptionRenewed`)
6. **Invariants**: (business rules that must always hold)

---

## Step 2: Generate Files

Create these files (replace `[Name]` and `[pkg]` with actual values):

### `internal/core/domain/model/[pkg]/[name].go`
```go
package [pkg]

import (
    "fmt"
    "github.com/google/uuid"
    "github.com/mijgona/salon-crm/internal/core/domain/model"
    "github.com/mijgona/salon-crm/internal/pkg/ddd"
)

type [Name] struct {
    *ddd.BaseAggregate[uuid.UUID]
    tenantID model.TenantID
    // ... fields
}

// New[Name] creates a new [Name] aggregate, validates invariants, raises domain event.
func New[Name](id uuid.UUID, tenantID model.TenantID /* ... */) (*[Name], error) {
    if /* invariant */ {
        return nil, fmt.Errorf("[name]: [invariant description]")
    }
    a := &[Name]{
        BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),
        tenantID:      tenantID,
    }
    a.RaiseDomainEvent(New[Name]Created(/* fields */))
    return a, nil
}

// Restore[Name] rehydrates from DB — no validation, no events.
func Restore[Name](id uuid.UUID, tenantID model.TenantID /* ... */) *[Name] {
    return &[Name]{
        BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),
        tenantID:      tenantID,
    }
}

// Getters
func (a *[Name]) TenantID() model.TenantID { return a.tenantID }

// [Action] performs [description] and enforces [invariant].
func (a *[Name]) [Action]( /* args */ ) error {
    if /* invariant check */ {
        return fmt.Errorf("[name].[action]: [reason]")
    }
    // mutate state
    a.RaiseDomainEvent(New[Name][Action]ed(/* fields */))
    return nil
}
```

### `internal/core/domain/model/[pkg]/[name]_[event].go`
```go
package [pkg]

import "github.com/google/uuid"

type [Name][Action]ed struct {
    EventID  uuid.UUID
    [Name]ID uuid.UUID
    TenantID uuid.UUID
    // event-specific fields
}

func New[Name][Action]ed(/* args */) [Name][Action]ed {
    return [Name][Action]ed{
        EventID:  uuid.New(),
        // ...
    }
}

func (e [Name][Action]ed) GetID() uuid.UUID { return e.EventID }
func (e [Name][Action]ed) GetName() string  { return "[pkg].[Name][Action]ed" }
```

### `internal/core/domain/model/[pkg]/[name]_test.go`
```go
package [pkg]_test

import (
    "testing"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "[module]/internal/core/domain/model/[pkg]"
)

func Test[Name]_New_HappyPath(t *testing.T) {
    a, err := [pkg].New[Name](uuid.New(), /* valid args */)
    require.NoError(t, err)
    assert.NotNil(t, a)
    // assert domain event raised
    events := a.GetDomainEvents()
    require.Len(t, events, 1)
    assert.Equal(t, "[pkg].[Name]Created", events[0].GetName())
}

func Test[Name]_New_InvariantViolation(t *testing.T) {
    _, err := [pkg].New[Name](uuid.New(), /* invalid arg */)
    require.Error(t, err)
    assert.Contains(t, err.Error(), "[expected message]")
}

func Test[Name]_[Action]_UpdatesState(t *testing.T) {
    a, _ := [pkg].New[Name](uuid.New(), /* valid args */)
    a.ClearDomainEvents()

    err := a.[Action](/* args */)

    require.NoError(t, err)
    events := a.GetDomainEvents()
    require.Len(t, events, 1)
    assert.Equal(t, "[pkg].[Name][Action]ed", events[0].GetName())
}
```

### `internal/core/ports/[name]_repository.go`
```go
package ports

import (
    "context"
    "github.com/google/uuid"
    "[module]/internal/core/domain/model/[pkg]"
)

type [Name]Repository interface {
    Add(ctx context.Context, tx Tx, entity *[pkg].[Name]) error
    Update(ctx context.Context, tx Tx, entity *[pkg].[Name]) error
    Get(ctx context.Context, tx Tx, id uuid.UUID) (*[pkg].[Name], error)
}
```

---

## Step 3: Wire in Composition Root

Add to `cmd/composition_root.go`:
```go
// [Name] repository
var [name]Repo ports.[Name]Repository
var [name]RepoOnce sync.Once
func (r *Root) [Name]Repository() ports.[Name]Repository {
    [name]RepoOnce.Do(func() {
        [name]Repo = [name]repoImpl.NewRepository(r.DB())
    })
    return [name]Repo
}
```

---

## Step 4: Verify

```bash
make vet
make test-domain
```
All tests must pass before marking feature complete.