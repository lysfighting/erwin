# Harness Skill

Codex-native leader surface for the repository-local orchestration harness prototype.

## What it wraps

- `bash scripts/harness-adapter.sh`
- `bash scripts/harness-missions.sh`
- `bash scripts/harness-home.sh`
- `bash scripts/harness-doctor.sh`
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
- inspect governed state
- inspect mission graph and mission-level desired state
- draft or create executable issues
- claim work into execution
- convene and close councils
- assign and govern bounded teammate execution
- transfer bounded ownership selectors formally when work changes hands
- express ownership with richer selector families such as `file:`, `dir:`, `module:`, `scope:`, and `interface:`
- adjudicate active signals when fresher run or meeting state supersedes them
- record teammate reports and sibling communication as structured signals
- run councils with explicit cross-examination linkage and inspect synthesis readiness before adoption
- inspect same-claim finding conflicts through `meetings authority-audit` when higher/lower authority evidence collides
- inspect cross-surface run conflicts through `runs authority-audit` when active signals, meeting outcomes, and current run state disagree
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

Prototype, but working.

## Verification

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
- `harness/` definitions
- `skills/harness/`
- `install-codex.sh`

After bootstrap, resync from the remembered source:

```bash
./scripts/harness-sync.sh
./scripts/harness-sync.sh status
```

The file list now comes from `harness/adapter-manifest.yml`, and `harness-sync status` reports adapter identity plus source/target manifest versions when the source repo is still available.
