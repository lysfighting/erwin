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

- `adapter_version`: `0.34.0-prototype`
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

- first-class mission objects
- separation between `latest_run_id` and `authoritative_run_id`
- graph-level authority audit with bundle grouping and explicit bundle governors
- ownership selectors across `file` / `dir` / `module` / `glob` / `scope` / `interface`
- path-family overlap plus same-boundary and hierarchical semantic overlap
- ownership transfer modes:
  - `move`
  - `copy`
  - `split`
- exclusion-aware parent/child delegation for split transfer
- role-bounded, authority-aware, claim/dimension-aware meetings
- wave / route / round obligations in councils
- auto-synthesis of provisional findings
- latest-response verdict adoption when a governed debate lane resolves
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

- fuller scheduler / rescheduling / execution-mode-selection behavior beyond the current strong leader control plane
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
3. [`docs/harness-architecture.md`](/Users/liyasong/Project/collov/erwin/docs/harness-architecture.md)
4. [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
5. [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)
6. this file
7. [`docs/research/harness-gap-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-gap-audit-2026-03-29.md)
8. [`docs/research/harness-research-drift-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-research-drift-audit-2026-03-29.md)
9. recent session logs
10. recent commits on `codex/harness-release-lifecycle`

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
