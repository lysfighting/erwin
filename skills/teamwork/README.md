# Teamwork Skills (v2 — GitHub-first)

GitHub-first team coordination. Issues are tasks, PRs are delivery, Actions are quality gates.

## Skills

| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `/team` | Init system + show dashboard | First time setup, daily overview |
| `/team-claim` | Claim a GitHub Issue as your mission | Pick your next task |
| `/team-drive` | Execute mission from Contract | During work: AI-assisted execution |
| `/team-ship` | Push, create PR, check CI, cleanup | When done: deliver via Pull Request |

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
Done:     /team-ship
```

### Step by Step

1. `/team` → see team dashboard (who's working on what, queued missions)
2. `/team-claim` → claim a GitHub Issue, generate AI-enriched Mission Contract, create branch
3. `/team-drive` → read Contract, execute sub-tasks systematically with verification
4. `/team-ship` → push branch, create PR (auto-closes Issue on merge), clean up Contract

## Requirements

- [Claude Code](https://claude.ai/code) installed
- [GitHub CLI](https://cli.github.com/) authenticated (`gh auth login`)
- Git repository with GitHub remote

## Architecture

See [docs/teamwork-v2-architecture.md](../../docs/teamwork-v2-architecture.md) for full design.

Key principles:
- **GitHub Issues = source of truth** — tasks live where developers already work
- **Mission Contract = AI execution view** — auto-generated, enriched with context files, ephemeral
- **Three-layer quality gates** — git hooks (local) + GitHub Actions (CI) + branch protection (review)
- **`gh` CLI as bridge** — all GitHub operations via authenticated CLI
- **Wraps /drive** — Contract sub-tasks = drive execution plan

## What Gets Generated

`/team` init creates:

| File | Purpose |
|------|---------|
| `.teamwork/config.yml` | Team roster, project settings, conventions |
| `.github/ISSUE_TEMPLATE/mission.yml` | Structured Issue template for missions |
| `.github/workflows/ci.yml` | CI workflow (lint + test + build) |
| `.githooks/pre-commit` | Local lint on commit |
| `.githooks/pre-push` | Local tests on push |

`/team-claim` generates:

| File | Purpose |
|------|---------|
| `.teamwork/active/MISSION-{N}.md` | AI-enriched execution Contract (ephemeral) |

## Version

- Skill version: 2.0.0
- Config schema: 1
