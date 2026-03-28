# Changelog

All notable changes to the `harness` skill will be documented in this file.

## [0.2.0] - 2026-03-28

### Added
- Bootstrapable Codex adapter path via `scripts/harness-bootstrap.sh`.
- Remembered-source resync via `scripts/harness-sync.sh` and adapter metadata.
- Structured issue drafting surface via `scripts/harness-issue.sh`.
- Release lifecycle surfaces: `ship` and `rc`.
- Ship now covers branch push and PR draft generation in addition to lifecycle state tracking.
- Ship can now open GitHub PRs directly via `gh pr create` when auth is available.
- Smoke coverage for issue, core, bootstrap, and release flows.

## [0.1.0] - 2026-03-28

### Added
- Initial Codex-native repository-local harness skill.
- Leader workflow documentation for `home`, `doctor`, `claim`, `drive`, `meetings`, and `maintenance`.
- Codex installation path via `install-codex.sh`.
- Smoke-test verification path via `bash scripts/harness-smoke-test.sh`.
