# Changelog

All notable changes to the `harness` skill will be documented in this file.

## [0.3.0] - 2026-03-28

### Added
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
