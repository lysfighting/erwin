# Checkpoint Schema: `design-note`

> Use when a mission requires an implementation plan before code. Typical missions: adding a new gate / guard / feature wired into existing code paths where the exact insertion points, test plan, and coupling risks matter.

Precedent: session 1 mission `#1298` (observability step_evidence gate).

## Required Sections

### 1. Insertion point(s)

For each place you will touch, give **file:line** backed by a grep command the team-lead can reproduce in the worktree. Include the current code excerpt (3-5 lines of context) so the team-lead sees what you're replacing, not just the target location.

### 2. Name & interface choices

If introducing a new env var / flag / function / class / label:
- Exact name
- Rationale (parallel to existing patterns? new namespace? why this level?)
- Parsing / semantics (truthy rules, default, precedence)

### 3. Test plan

- **Extend vs create**: if a test file already exists for the surface you touch, extend it. Cite the existing file. Only create a new test file if the surface is genuinely new.
- **Cases**: enumerate each new test case with one-line description.
- **Coupling regression guards**: if the new code must NOT interact with some other system (e.g., another env var, another module), write an explicit test pinning that non-interaction.
- **Command**: the exact shell command to run the new + old tests together.

### 4. Docs-diff plan

What documentation needs updating, and in which files/sections. Include READMEs, `.env.example`, CHANGELOG entries, and any runbooks.

### 5. Scope boundary call-out

If there is code that looks superficially in-scope but the Contract (intentionally or implicitly) excludes, name it explicitly and propose a follow-up mission. This is how you surface hidden scope ambiguity before code.

### 6. Non-Goal pins

Reference the Contract's Non-Goals section. For each Non-Goal you might accidentally violate, propose a concrete test that pins it (e.g., "assert function X is unchanged by this mission").

---

## Quality Bar

- Every factual claim about existing code has a grep or Read the team-lead can reproduce.
- Every new name is justified (not "seemed reasonable").
- Every test case maps to a Success Criterion in the Contract.
- No "we'll figure it out in code" — if you cannot decide now, escalate via `pushback`, not silent deferral.
