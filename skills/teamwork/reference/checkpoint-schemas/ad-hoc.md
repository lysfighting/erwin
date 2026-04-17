# Checkpoint Schema: `ad-hoc`

> Escape hatch when the mission doesn't fit `design-note` / `delta-table` / `research-result` / `risk-inventory`. Rare — most missions fit one of the four named shapes. If you reach for `ad-hoc`, take 30 seconds to verify no named shape works first.

## When to use

- Mission involves multi-modal deliverable (e.g., code stub + docs + config + question list, none of which dominates)
- Mission is meta (e.g., "review the existing X and propose a decision matrix")
- One-off investigations that are genuinely not repeat patterns — if you see a second `ad-hoc` mission in the same fleet, it's a signal to extract a new named schema

## Required Sections

### 1. Why `ad-hoc`, not a named schema

One paragraph. State which named schemas you considered and why none fit cleanly. If you cannot convincingly reject all four, **do not use `ad-hoc`** — use the closest named schema.

### 2. Deliverable structure

Propose your own section headings for this checkpoint. Include:
- What team-lead reads first (executive summary ≤3 sentences)
- Claims + evidence (always required, regardless of schema)
- Open questions back to team-lead (always required)

### 3. Evidence policy (inherited from all schemas)

Every factual claim backed by a grep / Read / gh output the team-lead can reproduce in the worktree. No "I believe" without cite.

### 4. Non-Goals re-check

Confirm you're within Contract scope. If scope is drifting, surface via `pushback`, not silently through `ad-hoc` flexibility.

### 5. Promotion note

If you think this `ad-hoc` shape will recur, end with a one-paragraph sketch of what a new named schema would look like. Team-lead may propose adding it to `reference/checkpoint-schemas/` post-session.

---

## Quality Bar

- Explicit rejection of the 4 named schemas with reasons. If you skip this section, team-lead rejects the checkpoint immediately and asks you to restructure.
- Executive summary reads cleanly on its own.
- Evidence density same as named schemas — the lack of a fixed shape does not lower the evidence bar.
- If the mission turns out to be a true one-off, `ad-hoc` is appropriate. If it's a recurring shape the fleet hadn't anticipated, propose the schema extraction explicitly.

## Warning

`ad-hoc` is NOT a license to be unstructured. It is a license to self-define structure when the named schemas don't fit. The team-lead's batch review gets harder when checkpoints have wildly different shapes — respect that cost and use `ad-hoc` sparingly.
