# Harness Completion Context Pack

> Date: 2026-03-28
>
> Purpose: define the minimum complete context a future session needs in order to finish the harness mission, not merely continue local productization work.

---

## Why This Exists

Recent sessions proved two things at once:

1. the repository now contains enough durable state for a new session to resume implementation safely;
2. a thin context pack will bias work toward visible package/productization tasks rather than the larger orchestration mission.

That second point matters.

If a future session only reads the current adapter/package state, it will naturally continue along lines such as:

- manifest hardening
- smoke expansion
- bootstrap/sync polish
- packaging cleanup

Those are useful, but they are not the whole mission.

The larger mission is to finish the new harness as:

- leader-led
- teammate-capable
- council-corrected
- cybernetically governed
- git-grounded
- Codex-native

This document exists to prevent future sessions from mistaking "adapter productization" for "mission completion."

---

## Reality Check

The current harness already has strong progress:

- protocol loop implemented: `issue -> claim -> drive -> ship -> rc`
- deterministic substrate implemented in the harness runtime core (now the Go CLI behind `scripts/harness-runtime.sh`)
- ship now includes review/check gates and deterministic merge
- rc now includes `fix`, `cancel`, `hotfix`, and `backport`
- manifest-driven bootstrap/sync/package validation now exist

But the following are still not fully realized:

1. fully automated council-style decision protocol modeled after `arch-review`
2. a final answer to product naming vs neutral adapter naming
3. a future Codex-native mapping from teammate runtime to real subagent execution, if that path is chosen

So the harness is already a strong single-brain control plane and reusable adapter package.
It is not yet the full leader-centric orchestration runtime originally envisioned.

---

## What "Complete" Means

A future session should not call this mission complete merely because bootstrap, sync, manifest, and smoke are polished.

Mission completion means all of the following are true:

1. The Codex adapter preserves the autonomous development lifecycle protocol.
2. The harness has a true leader runtime, not just a leader-flavored surface.
3. Teammates are first-class execution units with explicit contracts.
4. Teammates can exchange bounded signals under leader governance.
5. High-blast-radius or high-authority conflicts escalate into a real meeting/council protocol.
6. The resulting system still remains Codex-native, repo-decoupled, and profile-driven.
7. Product naming and public package boundaries are settled deliberately rather than accidentally by prototype drift.

---

## The Critical Interpretation

Do not interpret the harness as:

- a thin migration of `teamwork`
- a blind Codex port of `wing`
- just a bigger `SKILL.md`
- only a packaging/installer problem

Interpret it as:

**a Codex-native runtime adapter for the autonomous development lifecycle protocol, evolving toward a leader-led orchestration harness.**

---

## Required Context Tiers

### Tier 0 — Current Implementation Reality

Read these first in every fresh session:

- `AGENTS.md`
- `_project/CURRENT.md`
- `_project/sessions/2026-03-28-A.md`
- `_project/sessions/2026-03-28-B.md`
- `_project/sessions/2026-03-28-C.md`
- `_project/sessions/2026-03-29-A.md`
- `harness/README.md`
- `harness/adapter-manifest.yml`
- `skills/harness/SKILL.md`
- `scripts/harness-runtime.sh`
- `cmd/harness-state/main.go`
- `internal/harness/runtime.go`

Also inspect recent commits on `codex/harness-release-lifecycle`.

This tier answers:

- what exists now
- what was recently changed
- what the package/runtime boundary currently looks like

### Tier 1 — Protocol and Adapter North Star

Read these before making architecture-level decisions:

- `docs/research/protocol-core-spec-2026-03-28.md`
- `docs/research/codex-runtime-adapter-spec-2026-03-28.md`
- `docs/research/codex-native-orchestration-harness-analysis-2026-03-28.md`
- `docs/research/parity-matrix-teamwork-wing-protocol-2026-03-28.md`

This tier answers:

- what should survive from `teamwork`
- what should survive from `wing`
- what belongs to protocol vs adapter vs product surface

### Tier 2 — Orchestration / Governance Vision

Read these before touching leader/teammate/council design:

- `docs/research/cybernetics-kubernetes-arch-review-synthesis-2026-03-28.md`
- `docs/research/leader-spec-2026-03-28.md`
- `docs/research/teammate-spec-2026-03-28.md`
- `docs/research/meeting-protocol-spec-2026-03-28.md`
- `docs/research/state-authority-spec-2026-03-28.md`
- `docs/research/control-plane-spec-2026-03-28.md`
- `docs/research/signal-schema-2026-03-28.md`
- `docs/research/run-registry-schema-2026-03-28.md`

This tier answers:

- how the system should behave as a cybernetic control system
- what the leader is
- what a teammate is
- how meetings work
- what counts as truth

### Tier 3 — Historical Design Intent

Read these when a session risks drifting away from the original deeper mission:

- `_project/sessions/2026-03-13-AM.md`
- `_project/sessions/2026-03-18-AM.md`
- `_project/ideas/skills/tw2-architecture.md`
- `skills/arch-review/arch-review/SKILL.md`

This tier answers:

- why V4 was never just "parallel agents"
- why `wing` was meant to carry more than a teamwork rewrite
- why debate/council and data authority matter

---

## The Most Important Missing Semantics

A future session must explicitly remember these:

### 1. V4 is leader-centric, not flat multi-agent

There is one central brain/leader.

- leader assigns work
- teammates execute bounded tasks
- teammates report upward
- teammates may communicate sideways, but under governance

### 2. Meetings are not decorative

Meetings are not just a prettier note object.

They exist for:

- high ambiguity
- high blast radius
- authority conflicts
- decision forks

And they should increasingly resemble the strong parts of `arch-review`:

- data authority hierarchy
- explicit roles
- independent findings
- cross-examination
- synthesis by a chair/leader

### 3. Git/GitHub remain the order-bearing substrate

The system should not drift toward a shadow task manager.

Issues, PRs, branches, worktrees, and checks remain load-bearing.

### 4. Packaging progress is not the same as orchestration completion

Manifest, bootstrap, sync, smoke, install path:

all good, all necessary, none sufficient.

---

## Recommended Fresh-Session Prompt

When starting a new session intended to continue this mission, the opener should say something like:

> Continue the Codex-native harness mission. The current repo already has a strong adapter/package and a leader-centric control plane, but the mission is not complete until leader runtime, teammate runtime, governed teammate communication, and council-grade meeting behavior are fully realized. Read the Tier 0 and Tier 1 files first, then use Tier 2/3 as needed before making architecture changes.

That prompt is much safer than:

- "continue harness"
- "polish the package"
- "keep going"

because it keeps the session from collapsing back into surface-level productization.

---

## Bottom Line

A future session can continue safely from repository state alone.

But to **finish the mission**, not merely continue implementation, it must load both:

- the current adapter/package reality
- the deeper orchestration vision

Without both, the session will likely strengthen the adapter while still missing the full harness.
