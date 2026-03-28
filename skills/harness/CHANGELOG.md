# Changelog

All notable changes to the `harness` skill will be documented in this file.

## [Unreleased]

### Added
- `scripts/harness-build-runtime.sh` for building packaged platform binaries under `harness/bin/`.
- `scripts/harness-runtime-dist-smoke-test.sh` to prove the harness can run without `go` in `PATH` when a packaged runtime is present.
- Meeting role specs, locked decisions, and known-context recording so councils can carry stronger `arch-review`-style governance.

### Changed
- `scripts/harness-runtime.sh` now prefers packaged binaries before falling back to local Go build.
- The adapter manifest and bootstrap smoke now include the packaged-runtime distribution path.
- `harness-doctor.sh` now validates missing council source boundaries, missing expected output/context, and rejected findings without reasons.

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
