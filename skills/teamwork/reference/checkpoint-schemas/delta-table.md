# Checkpoint Schema: `delta-table`

> Use when the mission boils down to "current state X, need target state Y, identify the delta". Typical missions: capacity expansion, data backfill, UID/role/resource pool provisioning, any **quantitative gap** that has to be filled.

Precedent: session 1 mission `#1299` (prod loadtest user pool: 3 UIDs → ≥25 UIDs).

## Required Sections

### 1. Current state (publicly evidenced only)

Table column: `Verified today`. Values must be sourced from **public evidence** you can cite — comments on the Issue, README passages, non-secret config. **Do not read secret material** to populate this column. If you cannot verify without reading secrets, say "unknown — needs authoritative ops read" and stop there.

### 2. Target state

Table column: `Target`. From the Contract's Success Criteria. If the Contract's target is vague, flag it via `pushback` — do not guess.

### 3. Gap

Table column: `Gap`. Quantitative delta (e.g., "22 UIDs", "5 more regions", "3× throughput"). Simple arithmetic; no interpretation.

### 4. Candidate ops paths

For each candidate way to close the gap:
- **Path name + ≤1-sentence description**
- **Blast radius**: what breaks if the path misfires? (e.g., "touches prod users table", "external API rate limit", "downtime window")
- **Execution invariants**: what must hold for the path to be safe? (e.g., "must run during maintenance window", "requires ops pair review")
- **Pros / cons** in ≤3 bullets each

Do **not** mark a path as "recommended" unilaterally — the team-lead picks from the trade-off analysis.

### 5. Open questions back to team-lead

Before the team-lead can pick a path, what must they decide or confirm? Usually:
- Approved ops path (your table leaves this blank)
- Who performs which step (ops vs team-lead vs teammate)
- Artifact destination (off-repo path if secret material is generated)
- Whether this mission ships code or is pure-ops

### 6. Hard-rule self-audit

Explicit per-turn attestation:
- No secrets (tokens, real UIDs, credentials) read, echoed, or staged.
- No `gh issue comment` issued from this worktree.
- No branch push, no PR, no prod write.
- No operations outside the assigned worktree.

---

## Quality Bar

- The delta table row count ≤ number of rows needed to answer "how do we fill this gap?" — don't pad with context rows.
- Every candidate path's blast radius is stated in terms a non-specialist can act on.
- Every open question has a concrete answer shape (not "thoughts?").
- If no candidate path is actually safe, say so — escalate via `pushback` rather than invent a bad option.
