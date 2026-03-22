# Debug Recording Protocol

> When `--debug` flag is active, the Lead maintains a structured log of the entire review process.
> This log is the raw material for `/arch-review revolution` (process self-optimization).

---

## Debug Log Location

**Fixed directory**: `.claude/arch-review/debug/`
**Filename**: `debug-YYYY-MM-DD-HHMMSS.md` (timestamped, never overwritten)

Create the directory if it doesn't exist: `mkdir -p .claude/arch-review/debug/`

Revolution auto-reads from this directory — no path needed.

---

## What to Record

Append to the debug log at each phase transition. Not all sections apply to all triage levels:

| Section | Solo | Lightweight | Full |
|---------|------|-------------|------|
| 1. Triage | yes | yes | yes |
| 2. Review Plan | skip | minimal (L1+L3 only) | yes |
| 3. Teammate Findings | skip | yes | yes |
| 4. Cross-Examination | skip | yes | yes |
| 5. Lead Synthesis | yes (Lead's own findings) | yes | yes |
| 6. Outcome | yes | yes | yes |

### Section 1: Triage

```markdown
## [DEBUG] Triage
- **Start time**: [YYYY-MM-DD HH:MM]
- **Target**: [document path(s)]
- **Target size**: [line count, chapter count]
- **Scale assessment**: [Solo / Lightweight / Full debate]
- **User override?**: [yes: user chose X instead of recommended Y / no]
- **Rationale**: [why this scale was chosen]
```

### Section 2: Review Plan

```markdown
## [DEBUG] Review Plan
- **Data Authority Table**: [paste the table]
- **Roles assigned**: [list with missions and source restrictions]
- **Lessons injected**: [list lesson titles — not just count, revolution needs to track which ones were used]
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
- **End time**: [YYYY-MM-DD HH:MM]
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
