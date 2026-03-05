# /feature — Implement a New Feature

End-to-end feature implementation following the team workflow:
@planner → @tester → @coder

---

## Phase 1: Clarify Requirements (@planner)

Ask the user:
1. **What is the feature?** (1-2 sentences describing user need)
2. **Which bounded context?** (Client / Scheduling / Loyalty / Certificate / new?)
3. **Which aggregates are affected?** (or "unknown — help me figure out")
4. **Priority?** (Must-have / Should-have / Nice-to-have)

Then produce a mini-PRD:

```
🎯 Feature: [name]
📦 Context: [bounded context]
📝 User Story: As [role], I want [action], so that [value]
✅ Acceptance Criteria (Given/When/Then):
  - [AC-1]
  - [AC-2]
⚠️ Non-Goals: [what we explicitly won't do]
🔗 Affected aggregates/events: [list]
```

---

## Phase 2: Test Plan (@tester)

Before writing any code, define test cases:

```
🧪 Test Cases:
| ID | Name | Layer | Priority |
|----|------|-------|----------|
| TC-001 | Happy path | Integration | Critical |
| TC-002 | Invariant violation | Unit | Critical |
| TC-003 | Edge case | Unit | High |
| TC-004 | Concurrency / race | Integration | High |
```

Write Gherkin for each Critical TC.

---

## Phase 3: Implementation (@coder)

Follow this order strictly:
1. **Domain** — Value Objects, aggregate method, domain event (if new event needed)
2. **Ports** — update repository interface if needed
3. **Command/Query handler** — application layer
4. **Event handler** — if new cross-context integration
5. **Adapter** — HTTP handler + register route
6. **Repository** — postgres implementation (if new port)
7. **Composition root** — wire new handler in `cmd/composition_root.go`
8. **Tests** — unit (aggregate) + integration (handler with inmemory repos)

For each file, state:
```
📁 File: [path]
🎯 Purpose: [what this does]
💻 Code: [implementation]
```

---

## Phase 4: Verify

After implementation, run:
```bash
make vet        # must pass with 0 errors
make test-domain # domain unit tests must pass
```

Report:
```
✅ go vet: passed
✅ domain tests: X passed, 0 failed
📋 Test cases coverage: TC-001 ✅ TC-002 ✅ TC-003 ⏳
🔄 Next: [what's still pending]
```
