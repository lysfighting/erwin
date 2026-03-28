# Wing Rewrite — Executive Summary

> Date: 2026-03-28
>
> Decision frame: how to get `wing` to full `teamwork` capability without re-importing `teamwork`'s prompt bloat, coupling, or historical baggage.

---

## Core Judgment

`wing` should be treated as the future product surface.

Not because `teamwork` is obsolete as a protocol, but because `wing` is the correct architectural container for that protocol:

- official Anthropic skill layout
- drastically smaller runtime prompt surface
- shared reference files instead of duplicated rules
- skill-local path resolution via `${CLAUDE_SKILL_DIR}`
- migration-aware config compatibility (`.wing` -> `.teamwork` -> `.teamspace`)

The strategic goal should therefore be:

**Do not "port teamwork into wing" file-for-file.**

Instead:

**Reconstruct full teamwork capability inside wing's modular architecture, keeping protocol semantics while deleting monolithic packaging and legacy prompt debt.**

---

## What Must Be Preserved

These are the non-negotiable protocol assets that made `teamwork` valuable:

1. GitHub-first system model
   Issues, PRs, Actions, Milestones, Branch Protection remain the source of truth.

2. Six-phase lifecycle
   `issue -> claim -> drive -> ship -> rc`

3. Mission Contract model
   Public Issue for collaboration, local Contract for AI execution.

4. Deterministic mechanism layer
   Scripts own git, labels, PRs, config, contract manipulation, notification.

5. Cognitive posture
   Anti-compliance, blast-radius awareness, standalone test, ripple check, relentless delivery.

6. Release discipline
   RC lifecycle remains first-class.

If these survive, the protocol survives.

---

## What Must Not Be Carried Forward

These are the parts of `teamwork` that should be intentionally left behind:

1. Monolithic skill files
   `team.md` at 4,605 lines is exactly the failure mode the rewrite was meant to fix.

2. Repeated VE and doctrine blocks inside every skill
   Shared rules belong in `reference/` or dedicated docs, not duplicated runtime context.

3. Hardcoded install-path assumptions
   `~/.claude/commands/scripts/...` should not return.

4. Overgrown help and onboarding embedded in runtime skill files
   User education should live in suite docs, not balloon the invocation surface.

5. Legacy naming baggage as architecture
   Compatibility may remain, but `.wing/` should stay the canonical home.

6. Prompt-surface fixes that are really documentation problems
   If something can be taught outside the hot path, do not push it back into `SKILL.md`.

---

## Current State in One Table

| Area | Status in Wing | Judgment |
|---|---|---|
| Core lifecycle coverage | Present across 6 skills | Strong |
| Operation modularization | Present | Strong |
| Mechanism-layer scripts | Mostly preserved | Strong |
| Config compatibility | Present | Strong |
| Install model | Modernized | Strong |
| VE single source | Present | Strong |
| Suite docs | Thin | Gap |
| End-to-end parity proof | Incomplete | Gap |
| Migration narrative | Partial | Gap |
| Explicit replacement plan for leftover teamwork-only assets | Missing | Gap |

---

## Main Finding

The rewrite is already successful at the architecture level.

What is still missing is not the core protocol. It is the **completion layer**:

- packaged documentation
- parity proof
- migration checklists
- explicit decisions about what to preserve, simplify, or intentionally delete

That means the next phase should not be "more rewriting."

It should be:

**controlled parity completion**

---

## Recommended Goal Statement

Use this as the working north star:

> `wing` reaches full `teamwork` functional parity when every real workflow that a user can complete with `teamwork` can also be completed with `wing`, while `wing` remains modular, skill-local, reference-driven, and free of monolithic prompt regressions.

This definition is stronger than "feature parity":

- it measures workflow completion, not line-by-line sameness
- it protects against reintroducing old architectural debt

---

## Recommended Workstreams

### 1. Parity audit

Create a capability map:

- what `teamwork` can do
- whether `wing` already does it
- whether it should be preserved, redesigned, or deleted

### 2. Completion docs

Add missing suite-level assets to `skills/wing/`:

- `README.md`
- `ONBOARDING.md`
- `CHANGELOG.md`

These should explain the system without inflating `SKILL.md`.

### 3. Validation hardening

Bring `wing` to explicit confidence parity with `teamwork` by adding:

- e2e smoke coverage
- install validation
- migration validation
- documentation of known behavioral differences

### 4. Migration cleanup

Decide what legacy compatibility remains transitional versus permanent:

- `.teamwork` fallback
- `.teamspace` fallback
- old command installs

Make the exit criteria explicit.

---

## Recommended Priorities

### P0

- produce a formal parity-gap matrix for all teamwork capabilities
- add an explicit wing migration roadmap
- add e2e parity validation criteria

### P1

- add `skills/wing/README.md`
- add `skills/wing/ONBOARDING.md`
- add `skills/wing/CHANGELOG.md`
- document what was intentionally simplified vs still missing

### P2

- remove residual ambiguity around legacy support windows
- collapse any remaining duplicated operational guidance into shared docs or operations files

---

## Final Recommendation

Do not attempt a literal 1:1 transplant of `teamwork` text into `wing`.

The correct move is:

1. preserve protocol behavior
2. preserve deterministic scripts where still valid
3. add missing docs and tests outside the hot path
4. prove workflow parity
5. delete legacy surfaces only after parity is demonstrated

That is how `wing` becomes "all of teamwork's power, none of its prompt-era monolith debt."
