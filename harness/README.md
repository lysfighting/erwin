# Harness

This directory contains the committed definition layer for the emerging orchestration harness.

It is distinct from `.harness/`, which is reserved for local runtime state.

## Intended structure

- `profiles/` — committed project profiles
- `templates/` — optional reusable templates for contracts, runs, signals, and meetings
- `policies/` — optional committed policy defaults
- `bootstrap/` — adapter-owned bootstrap assets, including profile templates
- `adapter-manifest.yml` — machine-readable protocol, adapter, bootstrap, package, and smoke contract

## Design rule

Commit definitions here.
Keep live execution state in `.harness/`.

That split is part of the architecture, not just an implementation convenience.

The current productization direction is explicit three-layer separation:

- protocol semantics in the manifest and research docs
- reusable Codex adapter package in manifest-declared files and scripts
- repo-specific profile policy in `harness/profiles/` or `.harness/profile.local.yml`

## Implemented prototype surface

The first usable Codex-native harness prototype now exists in `scripts/`.

- `harness-adapter.sh` — adapter manifest introspection and package validation
- `harness-home.sh` — leader home view
- `harness-doctor.sh` — runtime diagnostics
- `harness-issue.sh` — structured issue drafting and creation
- `harness-claim.sh` — mission claim + contract/run creation
- `harness-drive.sh` — run-state control, meetings, reprovision, outcome adoption
- `harness-teammates.sh` — bounded teammate assignment, governed communication, and meeting escalation
- `harness-ship.sh` — verification, PR lifecycle, review/check gates, deterministic merge, merge readiness
- `harness-ship.sh` also supports branch push and PR draft generation
- `harness-ship.sh` can also open GitHub PRs directly when `gh` is authenticated
- `harness-ship.sh` now also handles post-merge cleanup of worktrees and local branches
- `harness-rc.sh` — release-candidate create/verify/promote/close plus surgical `fix`, `cancel`, `hotfix`, and `backport` flow
- `harness-meetings.sh` — meeting inspection and council output mutation
- `harness-maintenance.sh` — reconciliation + runtime normalization
- `harness-bootstrap.sh` — copy the harness adapter into another git repository
- `harness-sync.sh` — resync a bootstrapped adapter from its source repository

## Runtime model

- `harness/` stores committed definitions:
  profiles, templates, policies
- `.harness/` stores runtime state:
  runs, signals, contracts, meetings, maintenance artifacts

## Core loop

1. `home` shows the current governed state.
2. `doctor` checks runtime health and drift.
3. `claim` creates a run and contract from a mission.
4. `drive` moves the run through execution, blocking, meeting, and completion states.
5. `teammates` turns delegated work into first-class runtime state instead of implicit chat.
6. `meetings` records council outputs.
7. `maintenance` reconciles or normalizes runtime state.

The release side is no longer limited to the happy path.
The harness now also covers the high-friction release scenarios that used to depend on legacy operator memory:
RC surgical fix, RC cancellation, production hotfix prep, and dual-branch backport handling.

## Current guarantees

- file-backed runtime state with locking
- manifest-driven adapter packaging instead of bootstrap copy-list hardcoding
- validated run-state transitions
- tracked PR review state and CI/check state before merge
- deterministic merge path via `ship merge`
- meeting outputs that can be adopted back into runs
- reprovisioning of existing runs into branch/worktree mode
- reconciliation and normalization utilities
- leader-facing home/doctor views
- first-class teammate runs with bounded assignment scope and reporting contracts
- governed teammate-to-teammate communication through structured sibling signals
- meeting participants, phases, and recorded positions for arch-review-like councils
- machine-readable separation between protocol contract, adapter package, and profile contract

## Verification

Run the harness smoke test:

```bash
./scripts/harness-adapter.sh validate
./scripts/harness-issue-smoke-test.sh
./scripts/harness-smoke-test.sh
./scripts/harness-bootstrap-smoke-test.sh
./scripts/harness-release-smoke-test.sh
./scripts/harness-teammate-smoke-test.sh
./scripts/harness-all-smoke-test.sh
```

Bootstrap the harness into another git repository:

```bash
./scripts/harness-bootstrap.sh /path/to/target-repo
./scripts/harness-sync.sh --source /path/to/source-repo /path/to/target-repo
```

This copies the deterministic substrate, committed definition layer, smoke test, and a generated default profile into the target repository.
It also copies `skills/harness/` and `install-codex.sh`, so the target repo gets a Codex-native local skill entry point immediately.
The bootstrap now follows `harness/adapter-manifest.yml` for its package file list, bootstrap profile template, and metadata location.
The bootstrap writes local source metadata to `.harness/cache/adapter-source.json`, so later `./scripts/harness-sync.sh` runs can reuse the remembered source path.
`./scripts/harness-sync.sh status` now reports source/target adapter identity plus manifest version alignment when the source repo is still reachable.
