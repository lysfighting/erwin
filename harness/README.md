# Harness

This directory contains the committed definition layer for the emerging orchestration harness.

It is distinct from `.harness/`, which is reserved for local runtime state.

## Intended structure

- `profiles/` — committed project profiles
- `templates/` — optional reusable templates for contracts, runs, signals, and meetings
- `policies/` — optional committed policy defaults

## Design rule

Commit definitions here.
Keep live execution state in `.harness/`.

That split is part of the architecture, not just an implementation convenience.

## Implemented prototype surface

The first usable Codex-native harness prototype now exists in `scripts/`.

- `harness-home.sh` — leader home view
- `harness-doctor.sh` — runtime diagnostics
- `harness-claim.sh` — mission claim + contract/run creation
- `harness-drive.sh` — run-state control, meetings, reprovision, outcome adoption
- `harness-ship.sh` — verification, PR lifecycle, merge readiness
- `harness-rc.sh` — release-candidate create/verify/promote/close flow
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
5. `meetings` records council outputs.
6. `maintenance` reconciles or normalizes runtime state.

## Current guarantees

- file-backed runtime state with locking
- validated run-state transitions
- meeting outputs that can be adopted back into runs
- reprovisioning of existing runs into branch/worktree mode
- reconciliation and normalization utilities
- leader-facing home/doctor views

## Verification

Run the harness smoke test:

```bash
./scripts/harness-smoke-test.sh
./scripts/harness-bootstrap-smoke-test.sh
./scripts/harness-release-smoke-test.sh
```

Bootstrap the harness into another git repository:

```bash
./scripts/harness-bootstrap.sh /path/to/target-repo
./scripts/harness-sync.sh --source /path/to/source-repo /path/to/target-repo
```

This copies the deterministic substrate, committed definition layer, smoke test, and a generated default profile into the target repository.
It also copies `skills/harness/` and `install-codex.sh`, so the target repo gets a Codex-native local skill entry point immediately.
The bootstrap writes local source metadata to `.harness/cache/adapter-source.json`, so later `./scripts/harness-sync.sh` runs can reuse the remembered source path.
