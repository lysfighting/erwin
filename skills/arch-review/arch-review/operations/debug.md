# Debug Recording Protocol

> When `--debug` flag is active, the Lead maintains a structured log of the entire review process.
> This log is the raw material for `/arch-review revolution` (process self-optimization).

---

## Debug Log Location

Write to: `arch-review-debug-YYYY-MM-DD.md` in the same directory as the review target.
If reviewing `docs/architecture/`, write to `docs/architecture/arch-review-debug-YYYY-MM-DD.md`.

---

## What to Record

Append to the debug log at each phase transition. Use the structured format below — revolution needs parseable sections, not freeform prose.

### Section 1: Triage

```markdown
## [DEBUG] Triage
- **Target**: [document path(s)]
- **Scale assessment**: [Solo / Lightweight / Full debate]
- **User override?**: [yes: user chose X instead of recommended Y / no]
- **Rationale**: [why this scale was chosen]
```

### Section 2: Review Plan

```markdown
## [DEBUG] Review Plan
- **Data Authority Table**: [paste the table]
- **Roles assigned**: [list with missions and source restrictions]
- **Lessons injected**: [count — how many lessons from lessons.md were included]
- **User modifications**: [did user change anything in the plan?]
```

### Section 3: Teammate Findings (one per role)

For each teammate, record their FULL output as received via SendMessage:

```markdown
## [DEBUG] Findings: [role-name]
### Raw findings (as received)
[paste complete SendMessage content]

### Metadata
- **Finding count**: N
- **Data sources cited**: [which L-levels were actually used]
- **Forbidden sources cited?**: [yes/no — if yes, which ones]
```

### Section 4: Cross-Examination

Record each round of cross-examination:

```markdown
## [DEBUG] Cross-Examination Round [N]
### [role-A] challenges [role-B]:
[paste challenge content]

### [role-B] responds:
[paste response content]

### Lead intervention (if any):
[what the Lead said/did, or "none"]
```

### Section 5: Lead Synthesis Decisions

For each finding, record the Lead's decision:

```markdown
## [DEBUG] Synthesis
| ID | Finding | Decision | Reason |
|----|---------|----------|--------|
| CV-01 | [summary] | Confirmed P0 | Evidence from L1 verified |
| CC-03 | [summary] | Rejected | Used L4 data to contradict L3 |
| SS-02 | [summary] | Merged with CV-05 | Same root cause |
```

### Section 6: Outcome

```markdown
## [DEBUG] Outcome
- **Total findings**: N (confirmed: X, rejected: Y, merged: Z)
- **P0/P1/P2/Doc**: [counts]
- **Action chosen**: [fix inline / report only / selective]
- **Self-sufficiency re-check**: [did it happen? results?]
- **New lessons identified**: [count, brief list]
```

---

## Recording Discipline

- Record WHAT HAPPENED, not what should have happened
- Include failures and mistakes — they're the most valuable data for revolution
- If a teammate violated data source restrictions, record it explicitly
- If the Lead made a judgment call that could have gone either way, note the alternatives considered
- Keep raw teammate outputs verbatim — don't summarize or clean up
