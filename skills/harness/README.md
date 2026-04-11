# Harness Skill

Codex-native leader surface for the repository-local orchestration harness stable line.

If you want the full human-readable explanation first, read:

- [`docs/harness-explainer.md`](/Users/liyasong/Project/collov/erwin/docs/harness-explainer.md)

## What it wraps

- `bash scripts/harness-adapter.sh`
- `bash scripts/harness-missions.sh schedule --brief`
- `bash scripts/harness-missions.sh`
- `bash scripts/harness-missions.sh get issue-1234 --brief`
- `bash scripts/harness-runs.sh get RUN_ID --brief`
- `bash scripts/harness-signals.sh adjudicate --all --brief`
- `bash scripts/harness-meetings.sh summary MEETING_ID --brief`
- `bash scripts/harness-meetings.sh orchestrate MEETING_ID --brief`
- `bash scripts/harness-maintenance.sh reconcile --brief`
- `bash scripts/harness-home.sh` (`--brief` for a human-readable leader briefing)
- `bash scripts/harness-doctor.sh` (`--brief` for a human-readable health briefing)
- `bash scripts/harness-issue.sh`
- `bash scripts/harness-claim.sh`
- `bash scripts/harness-drive.sh`
- `bash scripts/harness-teammates.sh`
- `bash scripts/harness-ship.sh`
- `bash scripts/harness-rc.sh`
- `bash scripts/harness-meetings.sh`
- `bash scripts/harness-maintenance.sh`

## Best use

Use this when the model should act as the harness leader:

- validate adapter/package integrity before bootstrap or sync work
- inspect governed state, including a human-readable leader briefing before drilling into raw JSON
- inspect the scheduler portfolio as a human-readable leader object through `missions schedule --brief` before drilling into raw scheduler JSON
- inspect active mission topology as a human-readable leader object through `missions graph --brief` before drilling into raw graph JSON
- inspect one mission as a human-readable leader object through `missions get MISSION_ID --brief` before drilling into raw mission JSON
- inspect mission control drift as a human-readable leader object through `missions authority-audit MISSION_ID --brief` before drilling into raw mission-audit JSON
- inspect mission graph drift as a human-readable leader object through `missions authority-audit --all --brief` before drilling into raw bundle-audit JSON
- inspect one run as a human-readable leader object through `runs get RUN_ID --brief` before drilling into raw run JSON
- inspect run drift as a human-readable leader object through `runs authority-audit RUN_ID --brief` before drilling into raw audit JSON
- inspect active control signals as a human-readable leader object through `signals active [MISSION_ID] --brief` before drilling into raw signal JSON
- inspect signal repair as a human-readable leader object through `signals adjudicate SIGNAL_ID|--all --brief` before drilling into raw adjudication JSON
- inspect one bounded teammate as a human-readable leader object through `teammates get TEAMMATE_RUN_ID --brief` before drilling into raw teammate JSON
- inspect one active or closed council as a human-readable leader object through `meetings summary MEETING_ID --brief` before drilling into raw meeting JSON
- advance or refresh council orchestration as a human-readable leader object through `meetings orchestrate MEETING_ID --brief` before drilling into raw orchestration JSON
- inspect same-claim council authority as a human-readable leader object through `meetings authority-audit MEETING_ID --brief` before drilling into raw council-audit JSON
- inspect runtime repair as a human-readable leader object through `maintenance reconcile --brief` before drilling into raw reconciliation JSON
- inspect mission graph and mission-level desired state
- inspect mission-level authority through `missions authority-audit` when the freshest teammate activity and the leader-governed control run should differ
- inspect graph-level mission conflicts through `missions authority-audit --all` when active missions may be colliding across ownership boundaries
- draft or create executable issues
- claim work into execution
- convene and close councils
- assign and govern bounded teammate execution
- transfer bounded ownership selectors formally when work changes hands
- copy bounded ownership selectors across explicitly shared runs when paired execution or review should retain the same governed boundary
- express ownership with richer selector families such as `file:`, `dir:`, `module:`, `glob:`, `scope:`, and `interface:`
- let semantic selector families govern the same boundary across `scope:` and `interface:` labels, including hierarchical prefixes like `scope:release` and `scope:release/notes`
- adjudicate active signals when fresher run or meeting state supersedes them
- record teammate reports and sibling communication as structured signals
- run councils with explicit cross-examination linkage and inspect synthesis readiness before adoption
- inspect missing independent-review wave obligations before assuming a council is really ready for cross-examination
- let `meetings orchestrate` auto-synthesize provisional findings and unresolved follow-up actions from bounded council positions
- inspect whether conflicting council claims have been turned into explicit route obligations and then actually routed into linked cross-examination before trusting synthesis
- inspect whether a debated claim lane still requires another bounded debate round before synthesis is allowed to advance
- inspect same-claim finding conflicts through `meetings authority-audit` when higher/lower authority evidence collides
- inspect cross-surface run conflicts through `runs authority-audit` when active signals, meeting outcomes, and current run state disagree
- inspect artifact-family authority conflicts through `runs authority-audit` when contract snapshots, active run state, and live git state disagree
- inspect clean-intent-vs-live-git conflicts through `runs authority-audit` when a run says it is ready or verified but the live worktree is still dirty
- inspect mission-record drift and leader-over-teammate control adjudication through `missions authority-audit`
- inspect cross-mission ownership collisions, dependency-governed mission blocking, dependency-governed ownership adjudication, bundle-aware mission-graph groupings, and explicit bundle governors through `missions authority-audit --all`
- inspect bundle-level `required_action`, `selector_refs`, and bundle `resolution_kind` through `missions authority-audit --all` when the leader needs to know how to untangle or govern a connected mission bundle
- adopt meeting outcomes
- reprovision runs safely
- verify ship readiness and record PR lifecycle state
- track review and checks state until a PR is merge-ready
- merge via a deterministic ship command instead of manual state-only promotion
- push mission branches and draft PR payloads
- open GitHub PRs directly when `gh` auth is available
- perform post-merge cleanup for worktrees and local branches
- create, surgically fix, cancel, hotfix, backport, verify, and promote release candidates
- repair or normalize runtime drift

## Status

Stable first line.
Public package name: `Harness`.

The current durable package boundary is documented in:

- [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)

## Verification

```bash
./scripts/harness-adapter.sh validate
./scripts/harness-issue-smoke-test.sh
./scripts/harness-smoke-test.sh
./scripts/harness-bootstrap-smoke-test.sh
./scripts/harness-release-smoke-test.sh
./scripts/harness-teammate-smoke-test.sh
./scripts/harness-runtime-dist-smoke-test.sh
./scripts/harness-runtime-package-smoke-test.sh
./scripts/harness-all-smoke-test.sh
```

The release smoke now covers:

- `claim -> ship -> rc promote`
- PR review/check gate before merge
- deterministic `ship merge`
- RC surgical `fix` via cherry-pick onto an active RC branch
- RC `cancel` with branch cleanup
- production-style `hotfix` branch preparation
- dual-branch `backport` after production hotfix

## Install To Codex

```bash
./install-codex.sh harness
```

## Reuse In Another Repository

Bootstrap the harness adapter into another git repository:

```bash
./scripts/harness-bootstrap.sh /path/to/target-repo
```

The target repo will receive:

- `scripts/harness_*`
- `docs/harness-public-contract.md`
- `harness/` definitions
- `skills/harness/`
- `install-codex.sh`

After bootstrap, resync from the remembered source:

```bash
./scripts/harness-sync.sh
./scripts/harness-sync.sh status
```

The file list now comes from `harness/adapter-manifest.yml`, and `harness-sync status` reports adapter identity plus source/target manifest versions when the source repo is still available.
