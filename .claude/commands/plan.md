# /plan — Product Planning (@planner)

Turn an idea into a structured PRD ready for development.

---

## Step 1: Discovery

Ask the user these questions before writing anything:

1. **Problem**: "What user pain does this solve? Give me a real example."
2. **Frequency**: "How often does this happen? How many salons are affected?"
3. **Success**: "How will you know this feature is working? What metric changes?"
4. **Constraints**: "Any technical or business constraints I should know? (offline, budget, timeline)"

---

## Step 2: RICE Prioritization

If comparing multiple features:

| Feature | Reach | Impact (1-3) | Confidence % | Effort (sprints) | RICE Score |
|---------|-------|--------------|--------------|------------------|------------|
| [A] | % salons | | | | **R×I×C/E** |
| [B] | % salons | | | | **R×I×C/E** |

Recommend the highest RICE score with rationale.

---

## Step 3: PRD Output

```markdown
## Feature: [Name]

**JTBD**: When [situation], I want [action], so that [outcome]

### Bounded Context
[Which context: Client / Scheduling / Loyalty / Certificate]
[Which aggregates are touched]
[Which domain events are raised or consumed]

### User Stories

**US-1 (Must-have)**:
As [role], I want [feature], so that [value]

Acceptance Criteria:
Given [precondition]
When [action]
Then [outcome]
  AND [side effect]

**US-2 (Should-have)**:
...

### Non-Goals
- [Explicitly out of scope]
- [Deferred to next iteration]

### Trade-offs
- If we do X now → Y deferred to Q[n] because [reason]

### Metrics
- North Star impact: [how this moves the main metric]
- Leading indicator: [what to measure after release]
- Success threshold: [specific number]

### Technical Notes for @coder
- New aggregate method needed: `[Aggregate].[Method]()`
- New domain event: `[EventName]` with fields [list]
- New command handler: `[Name]CommandHandler`
- Affected ports: [repository interfaces]
```

---

## Step 4: Handoffs

After PRD is approved:
- → **@tester**: "Please write test cases for US-1 and US-2 AC"
- → **@coder**: "PRD ready, AC defined, tests written — please implement"
- → **@seller**: "New feature [name] available — update your demo script"
