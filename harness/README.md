# Harness

This directory contains the committed definition layer for the emerging orchestration harness.

If you want the human-readable overview first, read:

- [`docs/harness-explainer.md`](/Users/liyasong/Project/collov/erwin/docs/harness-explainer.md)

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

The current durable public-facing contract for that split is documented in:

- [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)

## Implemented prototype surface

The first usable Codex-native harness prototype now exists in `scripts/`.

The runtime core is now a Go CLI resolved by `scripts/harness-runtime.sh`.
When a packaged binary exists under `harness/bin/<os>-<arch>/harness-state`, the shell surfaces use it directly.
If no packaged binary is available, the runtime falls back to building a local Go binary under `.harness/bin/`.
That means normal harness use no longer requires `python3`, and repos can opt into a no-Go runtime path by shipping packaged binaries.

- `harness-adapter.sh` — adapter manifest introspection and package validation
- `harness-missions.sh` — scheduler portfolio inspection, mission registry inspection, mission graph view, mission-level authority audit, and leader summaries for `missions schedule --brief`, `missions graph --brief`, `missions get --brief`, and `missions authority-audit --brief` including `--all` graph mode
- `harness-home.sh` — leader home view, including `--brief` for a human-readable leader briefing
- `harness-doctor.sh` — runtime diagnostics, including `--brief` for a human-readable health briefing
- `harness-issue.sh` — structured issue drafting and creation
- `harness-claim.sh` — mission claim + contract/run creation, including `claim --brief` human-readable claim handoff
- `harness-drive.sh` — run-state control, progress reporting, dependency waiting, meetings, reprovision, outcome adoption, adopted follow-up resolution, and `drive report/start/resume/wait/block/complete/meeting/meeting-close/adopt-meeting/resolve-follow-up --brief` action briefings
- `harness-runs.sh` — direct run inspection, `runs get/authority-audit --brief` leader summaries, desired-state mutation, bounded ownership transfer, and run-level authority audit
- `harness-signals.sh` — active signal visibility, signal emission/state changes, and `signals active/adjudicate --brief` for human-readable control-signal guidance and repair
- `harness-teammates.sh` — bounded teammate assignment, governed communication, meeting escalation, `teammates assign --brief` for a human-readable delegation handoff, `teammates link --brief` for a human-readable governed peer-link handoff, `teammates report --brief` for a human-readable teammate report handoff, `teammates message --brief` for a human-readable governed peer-communication handoff, `teammates escalate --brief` for a human-readable council escalation handoff, and `teammates get --brief` for a human-readable bounded execution object
- `harness-ship.sh` — verification, PR lifecycle, review/check gates, deterministic merge, merge readiness, and `ship status --brief` for human-readable release-gate state
- `harness-ship.sh` also supports branch push and PR draft generation
- `harness-ship.sh` can also open GitHub PRs directly when `gh` is authenticated
- `harness-ship.sh` now also handles post-merge cleanup of worktrees and local branches
- `harness-rc.sh` — release-candidate create/verify/promote/close plus surgical `fix`, `cancel`, `hotfix`, and `backport` flow, with `rc list/get --brief` for human-readable RC state
- `harness-meetings.sh` — meeting inspection, `meetings summary/lock-decision/add-position/add-finding/add-action/resolve-action/orchestrate/authority-audit --brief` leader summaries, council orchestration, and council output mutation
- `harness-maintenance.sh` — reconciliation + runtime normalization, including `maintenance reconcile/normalize --brief` repair and migration summaries
- `harness-bootstrap.sh` — copy the harness adapter into another git repository
- `harness-sync.sh` — resync a bootstrapped adapter from its source repository

## Runtime model

- `harness/` stores committed definitions:
  profiles, templates, policies
- `.harness/` stores runtime state:
  runs, signals, contracts, meetings, maintenance artifacts

## Core loop

1. `home` shows the current governed state and can emit a leader-readable briefing with `--brief`.
2. `doctor` checks runtime health and drift and can emit a human-readable repair briefing with `--brief`.
3. `missions graph --brief` explains the active mission topology, `missions get MISSION_ID --brief` explains one mission's current state, `missions authority-audit MISSION_ID --brief` explains whether that durable mission record still matches authoritative run truth, and `missions authority-audit --all --brief` explains which connected mission bundle currently needs graph-level authority handling.
4. `claim --brief` creates a run and contract from a mission and explains the post-claim leader boundary before raw JSON inspection.
5. `drive report/start/resume/wait/block/complete/meeting/meeting-close/adopt-meeting/resolve-follow-up --brief` refreshes progress, moves the run through execution/blocking/completion states, opens councils, adopts council outcomes, and resolves adopted follow-up while explaining each mutation, next action, and whether teammates or councils are involved.
6. `signals active [MISSION_ID] --brief` explains current active control signals before humans adopt, resolve, or materialize them into meetings; `signals adjudicate SIGNAL_ID|--all --brief` explains which signals were actually resolved as superseded and why.
7. `teammates assign PARENT_RUN_ID TEAMMATE_ID ROLE [OBJECTIVE] --brief` turns delegation creation into a readable handoff, `teammates report TEAMMATE_RUN_ID SUMMARY --brief` turns teammate status into a readable control-signal handoff, `teammates message SOURCE_RUN_ID TARGET_RUN_ID KIND SUMMARY --brief` turns governed peer communication into a readable signal handoff, `teammates escalate TEAMMATE_RUN_ID QUESTION --brief` turns teammate uncertainty into a readable council escalation handoff, and `teammates get TEAMMATE_RUN_ID --brief` turns existing delegated work into a readable bounded execution object instead of implicit chat or raw teammate JSON.
8. `meetings summary MEETING_ID --brief` explains the live council state, `meetings lock-decision MEETING_ID TEXT --brief` explains the bounded council decision scope just recorded, `meetings add-position MEETING_ID TEAMMATE_RUN_ID PHASE STANCE TEXT --brief` explains what a recorded council position changed, `meetings add-finding MEETING_ID accepted|rejected|unresolved TEXT --brief` explains what a recorded judgment changed, `meetings add-action MEETING_ID TEXT --brief` explains the recorded follow-up action and who must own it next, `meetings resolve-action MEETING_ID ACTION_ID [STATE] --brief` explains action closure and adopted follow-up synchronization, and `meetings orchestrate MEETING_ID --brief` explains what changed when council orchestration advances or holds.
9. `maintenance reconcile --brief` explains what repair changed or why no repair was needed; `maintenance normalize --brief` explains what legacy runtime state was upgraded.

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
- explicit ownership claims plus overlap guardrails for active exclusive runs
- richer selector semantics for ownership claims, including path-family overlap between `file`, `dir`, `module`, and `glob`, plus semantic-family overlap and hierarchy between `scope` and `interface`
- bounded ownership transfer with `move` semantics for exclusive handoff and `copy` semantics for explicitly shared selectors
- dependency-aware waiting and declared desired state on runs and contracts
- first-class mission objects and a durable mission graph above run attempts
- first-class mission authority with `latest_run_id` separated from `authoritative_run_id`, so teammate activity cannot silently take over leader-facing mission control
- `maintenance normalize` can upgrade older run-only runtime state into mission-aware schema by backfilling ownership claims, desired state, and mission registry linkage
- governed teammate-to-teammate communication through structured sibling signals
- meeting participants, phases, role specs, locked decisions, and recorded positions for arch-review-like councils
- role-bounded meeting source discipline via `allowed_sources` / `forbidden_sources`
- rejected findings that carry explicit reasons instead of disappearing silently
- automated meeting phase advancement via `meetings orchestrate`
- provisional council auto-synthesis that can derive findings and unresolved-action prompts from bounded positions without overwriting manual findings already on the record
- debate-wave planning that surfaces missing independent-review obligations before councils move into cross-examination
- debate-plan routing that groups conflicting claim lanes, assigns bounded challenger route obligations, and tracks whether cross-examination has actually been linked for each conflict before synthesis is trusted
- bounded multi-round debate planning that can keep a council in cross-examination across successive replies until a response prevails on authority/freshness or the configured debate-round cap is reached
- durable council `phase_history`, linked cross-examination targets, and synthesis audits with blocker/evidence reporting
- leader-facing `meetings authority-audit` output for same-claim finding conflicts and higher-authority winner explanations
- leader-facing `runs authority-audit` output for cross-surface run/signal/meeting authority conflicts and adoption gaps
- leader-facing `runs authority-audit` output for artifact-family winner/loser rules across contract snapshots, active run state, and live git state
- leader-facing `runs authority-audit` output for clean-intent-vs-dirty-worktree conflicts when live git state contradicts governed verification or handoff intent
- leader-facing `missions authority-audit` output for mission-record drift plus leader-over-teammate and parallel-root-run control adjudication
- leader-facing `missions authority-audit --all` output for cross-mission ownership conflicts, dependency-governed mission blocking, dependency-governed ownership adjudication across overlapping mission boundaries, bundle-aware mission-graph views, explicit bundle governors when derivable, and bundle-level `required_action` guidance
- authority-level validation for signals and meeting findings
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
./scripts/harness-runtime-dist-smoke-test.sh
./scripts/harness-all-smoke-test.sh
```

Bootstrap the harness into another git repository:

```bash
./scripts/harness-bootstrap.sh /path/to/target-repo
./scripts/harness-sync.sh --source /path/to/source-repo /path/to/target-repo
```

This copies the deterministic substrate, committed definition layer, smoke test, and a generated default profile into the target repository.
It also copies `skills/harness/` and `install-codex.sh`, so the target repo gets a Codex-native local skill entry point immediately.
It also copies the durable package contract at [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md), so the bootstrapped repo keeps the same protocol/adapter/profile boundary instead of inheriting only implementation files.
It also copies the Go runtime module (`go.mod`, `go.sum`, `cmd/`, `internal/`) plus the build runner, so bootstrapped repos can build the same harness CLI locally without a Python runtime dependency.
If a repo wants to avoid a local Go dependency at runtime too, it can package platform binaries under `harness/bin/` using `./scripts/harness-build-runtime.sh`.
The bootstrap now follows `harness/adapter-manifest.yml` for its package file list, bootstrap profile template, and metadata location.
The bootstrap writes local source metadata to `.harness/cache/adapter-source.json`, so later `./scripts/harness-sync.sh` runs can reuse the remembered source path.
`./scripts/harness-sync.sh status` now reports source/target adapter identity plus manifest version alignment when the source repo is still reachable.
