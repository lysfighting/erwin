# Changelog

All notable changes to the `harness` skill will be documented in this file.

## [Unreleased]

### Added
- Multi-mission authority adjudication via `bash scripts/harness-missions.sh authority-audit --all`, including cross-mission ownership conflicts and dependency-governed mission blocking.
- Mission-level authority adjudication via `bash scripts/harness-missions.sh authority-audit MISSION_ID`, including explicit `authoritative_run_id` handling so teammate activity cannot silently replace leader-governed mission control.
- Runtime bundle packaging via `scripts/harness-package-runtime.sh` and `scripts/harness-runtime-package-smoke-test.sh`.
- Accepted public prototype positioning via ADR `0013`, with `Harness` now treated as the package's public prototype name.
- Artifact-family authority adjudication via `runs authority-audit`, covering stale contract snapshots versus active run/live state.
- Bounded second-round debate control via meeting `round_actions` and debate-round gating in `meetings orchestrate`.
- Durable public package contract via `docs/harness-public-contract.md`.
- `scripts/harness-build-runtime.sh` for building packaged platform binaries under `harness/bin/`.
- `scripts/harness-runtime-dist-smoke-test.sh` to prove the harness can run without `go` in `PATH` when a packaged runtime is present.
- Meeting role specs, locked decisions, and known-context recording so councils can carry stronger `arch-review`-style governance.
- Ownership claims, desired state, dependency waiting, and meeting orchestration on the harness control plane.
- `scripts/harness-missions.sh` plus a first-class mission registry and mission graph surface.
- Selector-aware ownership transfer via `bash scripts/harness-runs.sh transfer-ownership ...`.
- Signal authority adjudication via `bash scripts/harness-signals.sh adjudicate ...` and `bash scripts/harness-maintenance.sh reconcile`.
- Council phase history, linked cross-examination targeting, and synthesis-readiness reporting via `meetings summary` and `meetings orchestrate`.
- Same-claim meeting finding authority audit via `meetings authority-audit`, plus `--claim=` / `--dimension=` support on `meetings add-finding`.
- Richer ownership selector families with semantic overlap for `file:` / `dir:` / `module:` plus generic bounded-selector transfer checks.
- Cross-surface run authority adjudication via `bash scripts/harness-runs.sh authority-audit RUN_ID`, including higher-authority active-signal conflict detection and meeting-outcome adoption gaps.
- Claim-aware council auto-synthesis via `meetings add-position --claim=... --dimension=...` and `meetings orchestrate`, including provisional findings and unresolved follow-up actions.
- Clean-intent-vs-live-git authority detection via `runs authority-audit` when a dirty live worktree contradicts governed verification or handoff intent.
- Debate-route planning via meeting `debate_plan`, plus routed/unrouted conflicting claim reporting in `meetings summary` and `doctor`.
- Debate route obligations via meeting `debate_plan.route_actions`, including assigned bounded challenger roles/runs plus resolved/open route counts in `meetings summary` and `doctor`.
- Missing governed-worktree detection in `runs authority-audit`, so leader-facing authority review now catches absent provisioned worktrees in addition to branch drift and dirty-worktree intent drift.
- `glob:` ownership selectors, so path-pattern ownership can participate in selector-aware overlap checks and teammate admission control.
- Debate wave obligations in `meetings orchestrate`, so leader-facing council summaries now distinguish missing independent-review wave work from later cross-examination route work.
- Ship-gate authority adjudication in `runs authority-audit`, so leader-facing run review can explain when verification/review/checks state overrides an overly optimistic delivery status.

### Changed
- `harness-home.sh` and `harness-doctor.sh` now surface graph-level mission authority conflicts so multi-mission control drift is visible above per-run overlap noise.
- `scripts/harness-maintenance.sh normalize` now regenerates older mission state when `authoritative_run_id` is missing, so upgraded repos can reconcile into the new mission-control contract without manual surgery.
- Bootstrap smoke now validates the runtime bundle packaging surfaces in addition to runtime build/dist surfaces.
- The public contract and manifest now treat `Harness` as the public prototype name instead of leaving naming unresolved.
- Core smoke now validates that a stale contract branch loses to active run state before live-state authority drift checks continue.
- `meetings summary` and `doctor` now surface planned/resolved/open second-round debate obligations before synthesis.
- The adapter manifest now points at the durable public contract and packages it into bootstrapped target repos.
- `scripts/harness-runtime.sh` now prefers packaged binaries before falling back to local Go build.
- The adapter manifest and bootstrap smoke now include the packaged-runtime distribution path.
- `harness-doctor.sh` now validates missing council source boundaries, missing expected output/context, and rejected findings without reasons.
- `harness-doctor.sh` now also flags missing phase history, unlinked cross-examination, incomplete synthesis, and unresolved council findings without follow-up action.
- `harness-home.sh` and `harness-doctor.sh` now surface dependency-blocked runs, ownership overlap, and desired-state drift more explicitly.
- `harness-doctor.sh` now also flags unresolved run-level authority conflicts when fresher or higher-authority active signals disagree with current run state.
- Mission state is now materialized above runs so leader surfaces can reconcile against a durable mission graph instead of only per-run metadata.
- `scripts/harness-maintenance.sh normalize` now upgrades older run-only runtime state into the mission-aware schema by backfilling ownership claims, desired state, and mission linkage.
- Legacy `ownership_claims` are now normalized into `ownership_selectors`, so overlap checks and transfer semantics operate on structured selectors instead of raw strings alone.
- `meetings summary` now surfaces auto-generated council finding/action counts so leader views can see when synthesis is still purely manual versus partially automated.
- Worktree authority snapshots now ignore the harness-owned `.mission` marker so live-dirty detection tracks execution drift instead of harness bookkeeping noise.
- `doctor` now treats unrouted conflicting claims as a governance problem once a meeting has entered cross-examination or beyond.
- `meetings orchestrate` now turns conflicting claim lanes into explicit route obligations before synthesis, and `doctor` flags open route obligations separately from generic debate blockers.

## [0.6.0] - 2026-03-28

### Added
- First-class teammate runtime via `scripts/harness-teammates.sh`.
- Governed teammate-to-teammate communication using structured sibling signals.
- Meeting participant roster, phase tracking, and recorded positions for council-grade review.
- Dedicated `scripts/harness-teammate-smoke-test.sh`.

### Changed
- `harness-home.sh` now surfaces active teammate and coordination signal counts.
- `harness-doctor.sh` now validates teammate assignment quality plus meeting phase/participant/position health.
- The reusable adapter line advanced to `0.6.0-prototype`.

## [0.5.0] - 2026-03-28

### Added
- `scripts/harness-adapter.sh validate` for machine-readable adapter/package validation.
- Protocol identity, bootstrap contract, and profile contract metadata in `harness/adapter-manifest.yml`.

### Changed
- Bootstrap, sync, and smoke flows now consume manifest-declared adapter files instead of shell hardcoding.
- `harness-doctor.sh` now includes adapter validation in its runtime health report.

## [0.4.0] - 2026-03-28

### Added
- Ship review/check gates and deterministic merge.
- RC hotfix/backport coverage for dual-branch release topologies.
- Adapter manifest and one-command full smoke suite.
- Richer `harness-sync status` output with manifest version comparison.

## [0.3.0] - 2026-03-28

### Added
- Adapter manifest at `harness/adapter-manifest.yml`.
- One-shot smoke aggregator via `scripts/harness-all-smoke-test.sh`.
- PR review and checks gates via `scripts/harness-ship.sh review` and `scripts/harness-ship.sh checks`.
- Deterministic merge via `scripts/harness-ship.sh merge`, with local git merge and optional GitHub merge modes.
- RC surgical fix flow via `scripts/harness-rc.sh fix`, including cherry-pick and optional push to the active RC branch.
- RC cancel flow via `scripts/harness-rc.sh cancel`, including optional local and remote RC branch cleanup.
- Production hotfix preparation via `scripts/harness-rc.sh hotfix`.
- Dual-branch backport flow via `scripts/harness-rc.sh backport`.
- Release smoke coverage for PR review/check gates, deterministic merge, plus RC fix/cancel, hotfix, and backport flows.

## [0.2.0] - 2026-03-28

### Added
- Bootstrapable Codex adapter path via `scripts/harness-bootstrap.sh`.
- Remembered-source resync via `scripts/harness-sync.sh` and adapter metadata.
- Structured issue drafting surface via `scripts/harness-issue.sh`.
- Release lifecycle surfaces: `ship` and `rc`.
- Ship now covers branch push and PR draft generation in addition to lifecycle state tracking.
- Ship can now open GitHub PRs directly via `gh pr create` when auth is available.
- Ship now supports post-merge cleanup for worktrees and local branches.
- Smoke coverage for issue, core, bootstrap, and release flows.

## [0.1.0] - 2026-03-28

### Added
- Initial Codex-native repository-local harness skill.
- Leader workflow documentation for `home`, `doctor`, `claim`, `drive`, `meetings`, and `maintenance`.
- Codex installation path via `install-codex.sh`.
- Smoke-test verification path via `bash scripts/harness-smoke-test.sh`.
