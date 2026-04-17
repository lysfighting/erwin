# Checkpoint Schema: `risk-inventory`

> Use when the mission is "identify what could break / what's unsafe / what needs mitigation **before** any change is made." Typical missions: pre-migration audits, security-adjacent work, breaking-change proposals, cross-cutting refactor impact surveys.

Not used in session 1 or 2, but seeded in the Sprout vocabulary (KQ4 option hybrid) for future multi-mission work that involves invasive change planning.

## Required Sections

### 1. Mission scope restated

One paragraph: what the Contract intends to change, and what surface area that change touches. Include explicit "this mission does NOT change X" call-outs based on the Contract's Non-Goals.

### 2. Risk table

One row per identified risk:

| # | Risk description | Surface (file:line or system) | Trigger condition | Blast radius | Mitigation | Owner |
|---|---|---|---|---|---|---|
| R-1 | ... | ... | ... | ... | ... | ... |

- **Risk description**: what breaks, not what you'd do. E.g., "Active loadtest sessions see empty pool" not "need to coordinate with loadtest"
- **Surface**: concrete file:line if code-level, or "ops system X" if external. Must be greppable or referenceable.
- **Trigger condition**: what action or state causes this risk to materialize? "Any deploy that reruns migration", "Any teammate pushing without pulling first"
- **Blast radius**: how much breaks? Single user / single service / prod region / all users / silent data corruption (highest)
- **Mitigation**: proposed specific action (not "be careful"). If no mitigation known, write "NONE IDENTIFIED" — that's itself a valuable finding.
- **Owner**: who handles the mitigation — teammate / team-lead / ops / defer to follow-up mission

### 3. Risk prioritization

Rank the top 3 risks by **blast radius × trigger probability**. Explain why these three dominate.

### 4. Decision requests back to team-lead

For each risk where you propose multiple mitigations, or where "NONE IDENTIFIED", the team-lead must choose:
- Accept the risk (document + proceed)
- Pick a mitigation from your list
- Pivot the mission shape to avoid the risk entirely
- Escalate to user (out of fleet scope)

### 5. Non-Goals re-check

Confirm the risks you've inventoried are within Contract scope. If any risk points at a Non-Goal violation the Contract forced on the mission, surface that contradiction via `pushback`, not through this checkpoint.

---

## Quality Bar

- Each risk is testable — the "trigger condition" is specific enough that someone could cause it on purpose to reproduce.
- No risk is described in abstract safety-theater language ("we should think about X").
- Mitigations are concrete actions, not attitudes.
- Missing mitigations are explicitly named "NONE IDENTIFIED" — this is not a failure, it is honest output.
- Zero write ops during research phase. Observation only.
