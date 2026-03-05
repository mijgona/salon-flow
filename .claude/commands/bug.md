# /bug — Bug Investigation & Fix (@tester + @coder)

Systematic bug investigation, reporting, and fix workflow.

---

## Phase 1: Triage (@tester)

Ask the user:
1. **What happened?** (actual behavior)
2. **What was expected?** (expected behavior)
3. **Steps to reproduce?** (numbered list)
4. **Environment?** (local / staging / production)
5. **Frequency?** (always / sometimes / once)

Classify severity immediately:
- **Blocker**: system unusable, data loss risk → fix today
- **Critical**: key feature broken, no workaround → fix this sprint
- **Major**: feature degraded, workaround exists → next sprint
- **Minor**: cosmetic, edge case → backlog

---

## Phase 2: Bug Report (@tester)

Produce a structured bug report:

```markdown
## Bug Report [BA-XXX]

**Title**: [Verb + what + where]  e.g. "AppointmentBooked event not saved to outbox on concurrent booking"
**Severity**: [Blocker / Critical / Major / Minor]
**Layer**: [Domain / Application / Adapter / Infrastructure]
**Aggregate/Handler**: [e.g. BookAppointmentHandler, Appointment aggregate]

**Steps to Reproduce**:
1. [step]
2. [step]
3. [step]

**Expected**: [what should happen]
**Actual**: [what happens instead]
**Logs**: [relevant log lines or error messages]
**Environment**: [Go version, PostgreSQL version, OS]

**Root Cause Hypothesis**: [initial theory — which layer, which invariant]

**Regression Test**:
Given [setup]
When [action that triggers bug]
Then [correct behavior that must now pass]
```

---

## Phase 3: Investigation (@coder)

1. Read the affected file(s) — start from the layer identified in the bug report
2. Trace the execution path: HTTP handler → command handler → aggregate → repository → outbox
3. Identify the exact line where the invariant breaks
4. State root cause clearly before writing any fix

Investigation report:
```
🔍 Root cause: [exact file:line + explanation]
📁 Affected files: [list]
🔗 Related: [other places that may have same bug]
```

---

## Phase 4: Fix (@coder)

- Fix the root cause, not the symptom
- Add the regression test from Phase 2 **before** the fix (TDD)
- Verify fix doesn't break existing tests:
  ```bash
  make test-domain
  make vet
  ```
- Keep the fix minimal — don't refactor surrounding code

---

## Phase 5: Verify (@tester)

After fix:
```
✅ Regression test: PASS
✅ Existing tests: X passed, 0 failed
✅ go vet: clean
📋 Closed: BA-XXX
🔄 Monitor: [what to watch in production for 48h]
```
