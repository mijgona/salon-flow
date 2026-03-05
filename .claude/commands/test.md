# /test — Write Test Plan (@tester)

Generate a complete test plan with test cases for a feature or aggregate.

---

## How to use

- `/test BookAppointment` — test plan for a specific command handler
- `/test Appointment aggregate` — all tests for an aggregate
- `/test US-1 from plan` — tests for a User Story
- `/test TC-004` — expand a specific test case to full implementation

---

## Output Structure

### 1. Test Scope

```
🧪 Test object: [feature / aggregate / handler]
📦 Bounded context: [Client / Scheduling / Loyalty / Certificate]
🔗 Dependencies: [which repos, services, events]
```

### 2. Test Cases Table

```
| ID | Name | Layer | Priority | Type |
|----|------|-------|----------|------|
| TC-001 | Happy path — [description] | Integration | Critical | Functional |
| TC-002 | Invariant — [description] | Unit | Critical | Negative |
| TC-003 | Edge case — [description] | Unit | High | Boundary |
| TC-004 | Race condition — [description] | Integration | High | Concurrency |
| TC-005 | Event raised — [description] | Unit | High | Event |
| TC-006 | tenant isolation | Integration | High | Security |
```

### 3. Gherkin for Critical TCs

```gherkin
Scenario: TC-001 [name]
  Given [precondition]
  When [action]
  Then [expected outcome]
    And [side effect / event raised]

Scenario: TC-002 [invariant name]
  Given [setup that will violate invariant]
  When [action]
  Then error is returned containing "[invariant message]"
    And no domain events are raised
```

### 4. Go Test Stubs

For Critical tests, produce the Go test stub:

```go
func Test[Feature]_[Scenario](t *testing.T) {
    // Arrange
    repo := inmemory.NewClientRepository()
    handler := commands.New[Feature]Handler(repo, ...)
    cmd := commands.[Feature]Command{
        // fields
    }

    // Act
    err := handler.Handle(context.Background(), cmd)

    // Assert
    require.NoError(t, err)
    // verify state
    // verify events
}
```

### 5. Definition of Done

```
✅ Definition of Done for this feature:
- [ ] All Critical TCs pass
- [ ] All High TCs pass
- [ ] No Blocker or Critical bugs open
- [ ] go vet ./... clean
- [ ] make test-domain green
- [ ] Regression test for any fixed bug
```
