# Harness Session Handoff

> Status: Active
>
> Date: 2026-03-29
>
> Purpose: provide a durable, file-backed starting point for fresh sessions so `Harness` work can resume from repository state instead of relying on chat-context memory.

---

## Why This Exists

`Harness` is now large enough that chat-context alone is not a safe source of truth.

Fresh sessions should prefer:

1. committed files
2. machine-readable package/runtime metadata
3. recent session logs
4. recent commits

and only then any surviving conversation context.

This is intentional state persistence, not incidental documentation.

---

## Current Identity

`Harness` is:

- the public prototype/package name
- a Codex-native adapter for the Autonomous Development Lifecycle Protocol
- a leader-led orchestration harness

It is **not**:

- a renamed `teamwork`
- a renamed `wing`
- a generic standalone cloud agent platform

The identity split is:

- `protocol`: Autonomous Development Lifecycle Protocol
- `adapter`: `Harness`
- `machine-readable adapter id`: `codex-native-harness`

Authoritative references:

- [`AGENTS.md`](/Users/liyasong/Project/collov/erwin/AGENTS.md)
- [`_project/CURRENT.md`](/Users/liyasong/Project/collov/erwin/_project/CURRENT.md)
- [`docs/harness-documentation-map.md`](/Users/liyasong/Project/collov/erwin/docs/harness-documentation-map.md)
- [`docs/harness-architecture.md`](/Users/liyasong/Project/collov/erwin/docs/harness-architecture.md)
- [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
- [`_project/decisions/0013-harness-prototype-positioning.md`](/Users/liyasong/Project/collov/erwin/_project/decisions/0013-harness-prototype-positioning.md)

---

## Non-Negotiable Architecture

Never collapse these layers back together:

1. `protocol`
   Lifecycle semantics. The core loop is still:
   `issue -> claim -> drive -> ship -> rc`

2. `adapter`
   Codex-facing runtime, shell surfaces, bootstrap/sync flow, manifest, validation, packaging.

3. `profile`
   Repo-specific policy. Current proving-ground profile:
   [`harness/profiles/erwin.yml`](/Users/liyasong/Project/collov/erwin/harness/profiles/erwin.yml)

Also keep this boundary explicit:

- [`harness/`](/Users/liyasong/Project/collov/erwin/harness/README.md) = committed definition layer
- `.harness/` = local runtime state

That boundary is architectural, not accidental.

---

## Runtime Reality

The runtime is now Go-based.

Execution path:

- [`cmd/harness-state/main.go`](/Users/liyasong/Project/collov/erwin/cmd/harness-state/main.go)
- [`internal/harness/runtime.go`](/Users/liyasong/Project/collov/erwin/internal/harness/runtime.go)

Shell entrypoints are intentionally thin:

- [`scripts/harness-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-runtime.sh)
- [`scripts/harness-claim.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-claim.sh)
- [`scripts/harness-drive.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-drive.sh)
- [`scripts/harness-teammates.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-teammates.sh)
- [`scripts/harness-meetings.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-meetings.sh)
- [`scripts/harness-ship.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-ship.sh)
- [`scripts/harness-rc.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-rc.sh)

Runtime selection order:

1. `harness/bin/<os>-<arch>/harness-state`
2. `.harness/bin/harness-state`
3. local `go build`

Meaning:

- normal harness usage no longer depends on `python3`
- packaged binary is the preferred end-user path
- local Go toolchain is a fallback, not the desired steady state

---

## Machine-Readable Truth

The main machine-readable package boundary is:

- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)

Current important facts:

- `adapter_version`: `0.58.0-prototype`
- `public_name`: `Harness`
- `adapter_id`: `codex-native-harness`
- `public_contract_path`: [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)

This manifest is the best place to confirm:

- package identity
- package file list
- bootstrap contract
- profile contract
- capability surfaces
- smoke coverage

---

## Most Important Runtime Areas

When resuming implementation, the most important regions are in:

- [`internal/harness/runtime.go`](/Users/liyasong/Project/collov/erwin/internal/harness/runtime.go)

High-value areas:

- scheduler summary / snapshot generation / execution-mode selection / rescheduling hints
- selector / overlap / transfer model
- run authority audit
- mission-level and mission-graph authority
- meeting governance / wave plan / route plan / round plan / auto-synthesis / authority audit
- CLI subcommands

Do not treat `Harness` as “a bunch of shell scripts.”
The shell layer is intentionally thin; the Go runtime is the real control plane.

---

## What The System Already Does

Do not underestimate the current control-plane maturity.

`Harness` already has:

- leader-facing scheduler baseline with mission classification, execution-mode selection, bundle-aware deferral/escalation, bundle-role tallies, decision-first next-action prioritization, rescheduling hints, and durable `.harness/scheduler/snapshots/`
- `home` now mirrors that decision-first scheduler guidance instead of collapsing escalation-required bundles back into generic cross-mission conflict messaging
- scheduler now consumes profile-backed bounded parallel capacity, so overflow runnable missions can be deferred behind the current admitted batch instead of always widening all runnable work
- scheduler now consumes the execution policy `require_explicit_parallel_justification`, so parallel widening only happens when the current admitted batch carries explicit `parallel_justification`
- scheduler now enforces `P0` foreground serialization, so a runnable `P0` mission prevents widening into a parallel batch even when other isolated missions are also runnable
- scheduler now consumes profile-backed high-risk scopes, so configured boundaries such as `install.sh` or `skills/wing/` can escalate a mission into high-risk foreground serialization even when no manual `risk_level` was provided
- scheduler now consumes profile-backed meeting triggers, so some required human decisions become explicitly council-worthy with durable `recommended_decision_surface` and `meeting_trigger_refs` instead of only generic decision records
- escalation-required mission bundles now collapse into one bundle-scoped scheduler decision packet, so the leader sees one governing-bundle decision instead of one duplicated decision per mission
- scheduler now exposes a prioritized human-decision portfolio, including a `primary_required_human_decision`, deferred decision count, and scope/surface tallies, so leader-facing surfaces can reason about more than one pending governed decision at once
- mission classifications now mirror that human-decision queue through `human_decision_queue_position` and `primary_human_decision`, so leader-facing mission surfaces can see which pending decision is currently first-class
- lower-priority governed decisions now also emit explicit `defer_pending_primary_human_decision` reschedule actions, so the queue influences runtime behavior instead of remaining a read-only annotation
- scheduler now retains a sticky foreground mission across otherwise-tied scheduling passes, surfacing `selection_policy` and `previous_foreground_mission_id` so the leader can see when the control plane is intentionally preferring stability over freshness-only churn
- scheduler now treats fresh coordination/dependency/blocker signals as real rescheduling triggers instead of only including them in the snapshot payload
- scheduler now escalates stale runs, expired signals, and stale meetings into explicit `maintenance_actions` instead of leaving them only to `doctor`
- `home` now mirrors that maintenance backlog so leader-facing surfaces stay aligned with scheduler reality
- first-class mission objects
- separation between `latest_run_id` and `authoritative_run_id`
- graph-level authority audit with bundle grouping, explicit bundle governors, and pairwise governor adjudication from clear authority tuples
- ownership selectors across `file` / `dir` / `module` / `glob` / `scope` / `interface`
- path-family overlap plus same-boundary and hierarchical semantic overlap
- ownership transfer modes:
  - `move`
  - `copy`
  - `split`
- exclusion-aware parent/child delegation for split transfer
- role-bounded, authority-aware, claim/dimension-aware meetings
- wave / route / round obligations in councils
- deterministic current/blocked debate-batch execution state
- same-pass closure from synthesis to output when a council record is already ready
- direct synthesis when a council lane is aligned or already resolved by higher-authority evidence
- auto-synthesis of provisional findings
- latest-response verdict adoption when a governed debate lane resolves
- runtime-enforced council governance gaps, so missing authority hierarchy and rejected findings without reasons block council progression until repaired
- leader/teammate split with governed sibling communication

This is no longer “just a prototype shell workflow.”

---

## What Still Matters Most

Two statements are both true:

1. `Harness` already delivers the first complete prototype line.
2. The ultimate control-plane end state is still not fully complete.

These do not conflict.

The first means:

- runtime exists
- adapter exists
- packaging exists
- public contract exists
- major control-plane primitives exist

The second means higher-order control-plane semantics are still open.

The remaining high-value gaps are:

- fuller scheduler / rescheduling / execution-mode-selection behavior beyond the current baseline snapshot-and-hint control plane
- broader authority adjudication for larger mission bundles when no obvious governing dependency exists
- fuller `arch-review`-style automated debate execution beyond today’s route/wave/round/verdict-adoption baseline
- optional future mapping between Codex subagents and teammate runtime

---

## How To Read Research Correctly

These two audits remain important:

- [`docs/research/harness-gap-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-gap-audit-2026-03-29.md)
- [`docs/research/harness-research-drift-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-research-drift-audit-2026-03-29.md)

Use them for:

- control-plane maturity
- absorbed vs partial vs missing semantics
- avoiding design drift

Do **not** rely on them as the final truth for:

- naming
- distribution
- latest product-positioning details

Those areas were partially overtaken by later 2026-03-29 sessions and by:

- [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)
- recent session logs

---

## Fresh Session Start Order

For a new session that needs enough context to work safely, start here:

1. [`AGENTS.md`](/Users/liyasong/Project/collov/erwin/AGENTS.md)
2. [`_project/CURRENT.md`](/Users/liyasong/Project/collov/erwin/_project/CURRENT.md)
3. [`docs/harness-documentation-map.md`](/Users/liyasong/Project/collov/erwin/docs/harness-documentation-map.md)
4. [`docs/harness-architecture.md`](/Users/liyasong/Project/collov/erwin/docs/harness-architecture.md)
5. [`docs/harness-operating-model.md`](/Users/liyasong/Project/collov/erwin/docs/harness-operating-model.md)
6. [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
7. [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)
8. this file
9. [`docs/research/harness-gap-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-gap-audit-2026-03-29.md)
10. [`docs/research/harness-research-drift-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-research-drift-audit-2026-03-29.md)
11. recent session logs
12. recent commits on `codex/harness-release-lifecycle`

Recommended recent session logs:

- [`_project/sessions/2026-03-29-AL.md`](/Users/liyasong/Project/collov/erwin/_project/sessions/2026-03-29-AL.md)
- [`_project/sessions/2026-03-29-AM.md`](/Users/liyasong/Project/collov/erwin/_project/sessions/2026-03-29-AM.md)

---

## Fresh Session Rules

A fresh session should avoid these mistakes:

- do not reframe `Harness` as a `teamwork`/`wing` rename
- do not flatten `protocol / adapter / profile`
- do not treat shell wrappers as the real implementation layer
- do not confuse “first complete prototype line” with “control-plane end state fully solved”
- do not trust stale chat summaries over committed files

The repository should be treated as the persistent memory substrate.

That is the point of this file.
