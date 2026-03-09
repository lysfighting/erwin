# Teamwork Skills (v3 — GitHub-first)

GitHub-first team coordination. Issues are tasks, PRs are delivery, Actions are quality gates.

## Skills

| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `/team` | Init system + show dashboard | First time setup, daily overview |
| `/team config` | View config + edit instructions | Check or change settings |
| `/team doctor` | Health diagnostics + interactive fix | Check branch/worktree/issue health |
| `/team-issue` | Create/view/update/batch Issues | Create missions, manage Issues |
| `/team-claim` | Claim a GitHub Issue as your mission | Pick your next task |
| `/team-drive` | Execute mission from Contract | During work: AI-assisted execution |
| `/team-ship` | Push, create PR, check CI, cleanup | When done: deliver via Pull Request |
| `/team-ship done` | Post-merge cleanup | After PR merged: close Issue, clean worktree |
| `/team-ship review` | AI code review on PR | Review PR for correctness, security, quality |
| `/team-ship sync` | Rebase on base branch | Keep branch up to date |
| `/team-rc` | Release candidate lifecycle | Prepare staging, promote to production |

## Onboarding

**New to teamwork?** Read [ONBOARDING.md](ONBOARDING.md) — a step-by-step tutorial from zero to first PR, with examples and screenshots of every command's output.

## Quick Start

### First Time Setup

1. Install: `cd erwin && ./install.sh teamwork`
2. In your project: `claude` → `/team`
   - Auto-detects language, test/lint/build commands
   - Creates `.teamwork/config.yml`, GitHub Issue template, CI workflow, git hooks
   - Configures GitHub labels and branch protection

### Daily Workflow

```
Morning:  git pull && claude → /team → /team-claim
Working:  /team-drive
Done:     /team-ship → (after merge) /team-ship done
Review:   /team-ship review
Release:  /team-rc → verify staging → /team-rc promote
```

### Step by Step

1. `/team` → see team dashboard (who's working on what, queued missions)
2. `/team-issue` → create mission Issues from natural language
3. `/team-claim` → claim a GitHub Issue, generate AI-enriched Mission Contract, create branch
4. `/team-drive` → read Contract, execute sub-tasks systematically with verification
5. `/team-ship` → push branch, create PR (auto-closes Issue on merge), clean up Contract
6. `/team-rc` → cut RC branch, verify on staging, promote to production

## Requirements

- [Claude Code](https://claude.ai/code) installed
- [GitHub CLI](https://cli.github.com/) authenticated (`gh auth login`)
- Git repository with GitHub remote

## Architecture

See [docs/teamwork-ai-manual.md](../../docs/teamwork-ai-manual.md) for the complete manual.

Key principles:
- **GitHub Issues = source of truth** — tasks live where developers already work
- **Mission Contract = AI execution view** — auto-generated, enriched with context files, ephemeral
- **Three-layer quality gates** — git hooks (local) + GitHub Actions (CI) + branch protection (review)
- **`gh` CLI as bridge** — all GitHub operations via authenticated CLI
- **Wraps /drive** — Contract sub-tasks = drive execution plan
- **Prefixed labels** — `status:wip`, `priority:P1` — queryable, no collisions
- **Milestones for versions** — `versions.current` maps to GitHub Milestone
- **Full lifecycle** — issue → claim → drive → ship → done → rc

## Optional Features

All opt-in via `config.yml` (except worktree, which is per-user):

| Feature | Config Key | Default |
|---------|-----------|---------|
| Prefixed labels | `label_prefix` | Enabled (new installs) |
| Git worktree isolation | `git config --local teamwork.worktree` | `false` (per-user, not shared) |
| Version tracking | `versions.current` | Not set |
| Dual-branch mode | `conventions.production_branch` | Same as `base_branch` (single-branch) |

## What Gets Generated

`/team` init creates:

| File | Purpose |
|------|---------|
| `.teamwork/config.yml` | Team roster, project settings, conventions, optional features |
| `.github/ISSUE_TEMPLATE/mission.yml` | Structured Issue template for missions |
| `.github/workflows/ci.yml` | CI workflow (lint + test + build) |
| `.githooks/pre-commit` | Local lint on commit |
| `.githooks/pre-push` | Local tests on push |

`/team-claim` generates:

| File | Purpose |
|------|---------|
| `.teamwork/active/MISSION-{N}.md` | AI-enriched execution Contract (ephemeral) |

## Development Utilities

Scripts in `scripts/` not called by skills directly — used for maintenance:

| Script | Purpose | Usage |
|--------|---------|-------|
| `tw-bump-version.sh` | Bump version across all skill files, README, ONBOARDING, manual | `bash scripts/tw-bump-version.sh 3.6.0` |
| `tw-e2e-test.sh` | End-to-end validation harness (creates test Issue, claims, drives, ships, cleans up) | `bash scripts/tw-e2e-test.sh` |

## Version

- Skill version: 3.5.1
- Config schema: 3
