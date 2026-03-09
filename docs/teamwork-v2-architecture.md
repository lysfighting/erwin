# Teamwork v2 Architecture — GitHub-First Design

> ADR companion: `_project/decisions/0010-teamwork-v2-github-first.md`
> Status: ⚠️ **Archived design document** — see `docs/teamwork-ai-manual.md` for current v3.5+ implementation. This document describes the original v2.0 design phase.
> Validates on: `vi-agent` project

---

## Executive Summary

Teamwork v2 replaces the repo-centric markdown task system (v1) with a **GitHub-first architecture** where GitHub Issues are the source of truth, PRs are delivery gates, and Actions provide quality enforcement. Local "Mission Contracts" serve as **AI-optimized execution views** — auto-generated from Issues, enhanced with project context, and discarded after PR merge. The result is a universal skill that works with any GitHub-hosted project.

---

## Design Philosophy

```
GitHub = infrastructure layer (Issues, PRs, Actions, Branch Protection)
Contract = AI experience layer (auto-generated, context-enriched, ephemeral)
gh CLI = bridge between worlds
```

**Three design principles:**

1. **GitHub IS the system** — don't replicate what GitHub already does
2. **Enhance, don't duplicate** — Contracts add what Issue body can't (file paths, test commands, AI notes)
3. **Ephemeral over persistent** — local state is disposable; GitHub state is durable

---

## System Architecture

```
                    ┌──────────────────────────────────────┐
                    │           GitHub (Source of Truth)     │
                    │                                        │
                    │  Issues ──── Milestones ──── Projects  │
                    │    │              │              │      │
                    │  Labels     Progress %      Kanban     │
                    │  Assignees                    View     │
                    │    │                                    │
                    │    ▼                                    │
                    │  PRs ──── Actions CI ──── Branch Rules  │
                    │    │         │                │         │
                    │  Reviews   lint/test/build  Require ✓  │
                    └────┬──────────────────────────┬────────┘
                         │ gh CLI                   │
                    ┌────▼──────────────────────────▼────────┐
                    │            Local Repository             │
                    │                                         │
                    │  .teamwork/                              │
                    │    config.yml        (team config)       │
                    │    active/           (.gitignore'd)      │
                    │      MISSION-042.md  (AI execution view) │
                    │                                         │
                    │  .github/                                │
                    │    ISSUE_TEMPLATE/mission.yml            │
                    │    workflows/ci.yml                      │
                    │                                         │
                    │  .githooks/                              │
                    │    pre-commit        (lint + format)     │
                    │    pre-push          (tests)             │
                    └─────────────────────────────────────────┘
```

---

## Lifecycle Flow

```
/team init
  │  Create config.yml, Issue template, Actions CI, git hooks
  │  Configure branch protection rules
  │  Detect project language/framework
  ▼
/team board
  │  gh issue list → formatted dashboard
  │  gh pr list → in-progress work
  │  Merge count leaderboard
  ▼
/team claim #42
  │  gh issue view #42 --json → parse structured content
  │  Generate .teamwork/active/MISSION-042.md (AI execution view)
  │  Create branch: mission/042-{slug}-{username}
  │  gh issue edit #42 --add-assignee @me --add-label in-progress
  ▼
/team drive
  │  Read Contract → enter drive execution loop
  │  Progress: update Contract checkboxes + commit incrementally
  │  Commits: feat(scope): description | Mission: #42
  ▼
/team ship
  │  Pre-flight: sub-tasks done? local tests pass?
  │  git push origin {branch}
  │  gh pr create --title "..." --body "Closes #42\n\n..."
  │  Wait for Actions CI (optional)
  │  Cleanup: remove local Contract
  ▼
  PR merged → Issue #42 auto-closed → done
```

---

## Skill Specifications

### 1. `/team` — Init + Dashboard

**First run (init):**

1. Detect project:
   - Language: package.json → Node, pyproject.toml → Python, go.mod → Go, etc.
   - Test command: `npm test`, `pytest`, `go test ./...`
   - Lint command: `npm run lint`, `ruff check .`, `golangci-lint run`
   - Build command: `npm run build`, `python -m build`, `go build ./...`

2. Interactive setup (AskUserQuestion):
   - Team members (GitHub usernames + roles)
   - Base branch (default: main)
   - Quality preferences (hooks, CI, review, branch protection)

3. Generate files:
   - `.teamwork/config.yml` — team configuration
   - `.github/ISSUE_TEMPLATE/mission.yml` — structured Issue template
   - `.github/workflows/ci.yml` — language-specific CI workflow
   - `.githooks/pre-commit` — lint + format
   - `.githooks/pre-push` — test suite
   - `.gitignore` entry: `.teamwork/active/`

4. Configure GitHub (via `gh` CLI):
   - Branch protection on base branch (require CI, require review)
   - Labels: `mission`, `P0`, `P1`, `P2`, `P3`, `in-progress`, `review`

5. Commit + push

**Subsequent runs (dashboard):**

```
┌─────────────────────────────────────────────┐
│  TEAM DASHBOARD — my-project                │
│  Milestone: v0.1 MVP (4/10 = 40%)          │
├─────────────────────────────────────────────┤
│                                             │
│  👤 liyasong (tech-lead)                    │
│     🔧 #42 Add user authentication [P1]    │
│     Branch: mission/042-add-auth-liyasong   │
│     PR: —                                   │
│                                             │
│  👤 qianhua (backend)                       │
│     ✅ Idle (last merged: #38)              │
│                                             │
│  QUEUED:                                    │
│     #45 Add rate limiting [P1]              │
│     #47 Write API docs [P2]                 │
│                                             │
│  MERGE COUNT:                               │
│     liyasong: 6 | qianhua: 4               │
└─────────────────────────────────────────────┘
```

Data sources:
- `gh issue list --label mission` → all mission issues
- `gh pr list` → in-progress PRs
- `gh api` → merge count per assignee
- `.teamwork/config.yml` → team roster

---

### 2. `/team-claim` — Claim Issue → Contract → Branch

**Args routing:**

| Input | Action |
|-------|--------|
| `#42` | Claim specific Issue |
| `list` | Browse available Issues |
| `create "title"` | Create new Issue from template |
| (empty) | Auto-select next by priority (P0 > P1 > P2 > P3) |

**Claim flow:**

1. Validate: user in team? already have active mission? (enforce one-at-a-time)
2. `gh issue view #42 --json title,body,labels,assignees,milestone`
3. Parse Issue body (structured via template)
4. Generate Contract: `.teamwork/active/MISSION-042.md`
   - Copy structured fields from Issue
   - AI-enhance: scan project, add relevant file paths, detect test commands
5. Create branch: `mission/042-{slug}-{username}`
6. Assign + label: `gh issue edit #42 --add-assignee @me --add-label in-progress`
7. Output mission briefing

---

### 3. `/team-drive` — Execute Mission

Wraps `/drive` methodology with Mission Contract context.

1. Read Contract from `.teamwork/active/MISSION-{N}.md`
2. Display briefing (objective, sub-tasks, acceptance criteria)
3. Enter drive execution loop:
   - For each sub-task: Read → Implement → Verify → Mark done
   - Update Contract: check off sub-task with timestamp
   - Commit: `feat(scope): description | Mission: #42`
   - Do NOT push individual commits (save for /team-ship)
4. After all sub-tasks: verify acceptance criteria
5. Output completion summary

**Drive integration:**
- Contract = drive task input (replaces Phase Ω problem definition)
- Sub-tasks = drive execution plan
- Acceptance criteria = drive success gates

---

### 4. `/team-ship` — Deliver Mission

1. **Pre-flight checks:**
   - On mission branch? All sub-tasks in Contract done?
   - Run local tests: `{config.project.test_command}`
   - Clean working tree?

2. **Push:**
   ```bash
   git push -u origin {branch}
   ```

3. **Create PR:**
   ```bash
   gh pr create \
     --title "feat: {mission title}" \
     --body "Closes #{issue}\n\n## Changes\n{summary}\n\n## Test\n{test results}"
   ```

4. **CI check (optional):**
   ```bash
   gh pr checks {pr-number} --watch
   ```

5. **Review:**
   - If `config.quality.review_required`: request review from tech-lead
   - Optional: AI-powered review (read diff, check for issues)

6. **Cleanup:**
   - Remove `.teamwork/active/MISSION-{N}.md`
   - Remove `in-progress` label: `gh issue edit #{N} --remove-label in-progress`
   - Add `review` label: `gh issue edit #{N} --add-label review`

7. **Output:** PR URL + status

**Post-merge (manual or hook):**
- Issue auto-closed by "Closes #N" in PR body
- Branch auto-deleted (GitHub setting)

---

## File Specifications

### .teamwork/config.yml

```yaml
schema_version: 1
skill_version: 2.0.0

team:
  - github: liyasong
    role: tech-lead
  - github: qianhua
    role: backend

project:
  language: typescript       # auto-detected
  test_command: "npm test"
  lint_command: "npm run lint"
  build_command: "npm run build"

conventions:
  branch_pattern: "mission/{issue}-{slug}-{user}"
  base_branch: main
  commit_format: "type(scope): description | Mission: #{issue}"

quality:
  hooks: true
  ci: true
  review_required: true
  branch_protection: true
```

### .teamwork/active/MISSION-042.md (AI Execution View)

```markdown
---
issue: 42
url: https://github.com/owner/repo/issues/42
title: "Add user authentication"
assignee: liyasong
priority: P1
labels: [mission, backend, P1]
branch: mission/042-add-auth-liyasong
milestone: "v0.1 MVP"
claimed: 2026-03-02T10:00:00Z
---

# MISSION-042: Add user authentication

## Objective
Implement JWT-based authentication for the API server.

## Sub-tasks
- [ ] Define auth middleware
- [ ] Add JWT token generation
- [ ] Write login/logout endpoints
- [ ] Add tests for auth flow

## Acceptance Criteria
- All auth endpoints return correct status codes
- JWT tokens expire after 24 hours
- Tests cover happy path + invalid token + expired token

## Context Files
- src/middleware/index.ts — existing middleware stack
- src/routes/ — route definitions
- tests/ — test directory
- package.json — need to add jsonwebtoken dependency

## Test Command
npm test -- --grep auth

## AI Notes
(populated during execution)
```

### .github/ISSUE_TEMPLATE/mission.yml

```yaml
name: Mission Contract
description: Structured task for team execution
labels: ["mission"]
body:
  - type: dropdown
    id: priority
    attributes:
      label: Priority
      options:
        - P0 (critical)
        - P1 (high)
        - P2 (medium)
        - P3 (low)
    validations:
      required: true

  - type: textarea
    id: objective
    attributes:
      label: Objective
      description: What needs to be accomplished?
    validations:
      required: true

  - type: textarea
    id: subtasks
    attributes:
      label: Sub-tasks
      description: Break down into checkable items
      value: |
        - [ ] Task 1
        - [ ] Task 2
    validations:
      required: true

  - type: textarea
    id: acceptance
    attributes:
      label: Acceptance Criteria
      description: How do we know this is done?
    validations:
      required: true

  - type: textarea
    id: context
    attributes:
      label: Context & References
      description: Relevant files, docs, links

  - type: input
    id: test_command
    attributes:
      label: Test Command
      description: Command to verify this mission
      placeholder: "npm test -- --grep auth"
```

### .github/workflows/ci.yml (Node.js example)

```yaml
name: CI

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: npm
      - run: npm ci
      - run: npm run lint
      - run: npm test
      - run: npm run build
```

### .githooks/pre-commit

```bash
#!/bin/bash
# Teamwork v2 — local quality gate (fast)
# Runs: lint + format check

set -e

CONFIG=".teamwork/config.yml"
if [ ! -f "$CONFIG" ]; then
  exit 0  # not a teamwork project, skip
fi

# Extract lint command from config (simple grep, no yq dependency)
LINT_CMD=$(grep 'lint_command:' "$CONFIG" | sed 's/.*lint_command: *//' | tr -d '"')

if [ -n "$LINT_CMD" ]; then
  echo "🔍 Running lint..."
  eval "$LINT_CMD"
fi
```

### .githooks/pre-push

```bash
#!/bin/bash
# Teamwork v2 — local quality gate (thorough)
# Runs: test suite

set -e

CONFIG=".teamwork/config.yml"
if [ ! -f "$CONFIG" ]; then
  exit 0
fi

TEST_CMD=$(grep 'test_command:' "$CONFIG" | sed 's/.*test_command: *//' | tr -d '"')

if [ -n "$TEST_CMD" ]; then
  echo "🧪 Running tests..."
  eval "$TEST_CMD"
fi
```

---

## Version Management

### Three Levels

| Level | What | How |
|-------|------|-----|
| **Skill version** | The teamwork skill itself | `skill_version` in config.yml + CHANGELOG.md |
| **Schema version** | config.yml format | `schema_version` in config.yml (migration on mismatch) |
| **Project milestones** | Team's work versions | GitHub Milestones (`gh api` to query progress) |

### Skill Distribution

```
erwin/skills/teamwork/
  ├── team.md              ← /team skill
  ├── team-claim.md        ← /team-claim skill
  ├── team-drive.md        ← /team-drive skill
  ├── team-ship.md         ← /team-ship skill
  ├── templates/           ← CI, hooks, Issue template sources
  ├── CHANGELOG.md
  └── README.md

install.sh → symlink *.md → ~/.claude/commands/
git pull   → auto-update (symlinks point to source)
```

### Freshness Check (built into /team)

```
On every /team invocation:
  1. Read skill_version from installed skill header
  2. Compare with config.yml's skill_version
  3. If mismatch → warn + suggest: cd erwin && git pull && ./install.sh
```

---

## Quality Gate Architecture

```
        Developer Machine                    GitHub
        ─────────────────                    ──────

        pre-commit hook
        ├── lint          ◀─── 秒级 ───▶  (not needed on CI
        └── format check                    if hooks work)

        pre-push hook
        └── test suite    ◀─── 分钟级 ──▶  Actions CI
                                            ├── lint
                                            ├── test
                                            └── build

                                            Branch Protection
                                            ├── ✓ Require CI pass
                                            ├── ✓ Require PR review
                                            └── ✓ Require linear history
```

**Defense in depth:**
- Hooks catch obvious errors fast (developer experience)
- Actions catch everything (team safety net)
- Branch protection enforces process (policy)

---

## Key Differences from v1

| Aspect | v1 (current) | v2 (this design) |
|--------|-------------|-------------------|
| Task source | `_teamwork/missions/*.md` | GitHub Issues |
| Board | Read N member files | `gh` API → formatted output |
| Status tracking | `_teamwork/members/{name}.md` | GitHub (assignee, labels, PR status) |
| PR creation | Optional | Mandatory |
| Code review | None | AI-assisted + human |
| CI | None | Auto-generated Actions |
| Local hooks | None | Auto-generated pre-commit/pre-push |
| Version tracking | None | GitHub Milestones + skill versioning |
| Merge count | None | `gh api` query |
| Offline support | Full (all in repo) | Degraded (Contract cached locally) |
| Local files | ~10+ markdown files | 1 config + ephemeral contracts |

---

## Borrowed from vi-agent

| Pattern | What we took | What we changed |
|---------|-------------|-----------------|
| Mission lifecycle | claim → execute → deliver loop | GitHub Issues as source, not board.md |
| `gh pr create` | PR as delivery mechanism | Auto-close Issue via "Closes #N" |
| Merge count metric | Throughput measurement | Query via `gh api`, not board.md |
| Skill freshness | Version check on startup | Lighter: version string compare, not MD5 |
| Structured config | YAML configuration | Simpler schema, auto-detect project |
| Code review | Automated PR review | Integrated into /team-ship, not separate agent |

**What we explicitly did NOT take:**
- `board.md` as source of truth (GitHub IS the board)
- Port isolation (too project-specific for universal skill)
- Forced worktree (Claude Code already supports worktrees)
- Separate code-reviewer agent (folded into ship flow)
- Slack/Telegram messaging (out of scope for v2.0)

---

## Implementation Roadmap

### v2.0.0 — Full Flow MVP

**Phase 1: Foundation (1 session)**
- [ ] Write `/team` skill (init + dashboard)
- [ ] Write config.yml schema + auto-detection logic
- [ ] Write CI template generation (Node, Python, Go)
- [ ] Write Issue template (mission.yml)
- [ ] Write git hooks templates

**Phase 2: Claim + Contract (1 session)**
- [ ] Write `/team-claim` skill
- [ ] Implement Issue → Contract generation
- [ ] Branch creation + assignment logic
- [ ] Test on vi-agent: create Issue, claim it, verify Contract

**Phase 3: Drive Integration (1 session)**
- [ ] Write `/team-drive` skill (wraps /drive with Contract context)
- [ ] Contract progress tracking (checkbox updates)
- [ ] Commit message format with Mission reference

**Phase 4: Ship + Close Loop (1 session)**
- [ ] Write `/team-ship` skill
- [ ] PR creation with Issue linking
- [ ] CI check integration
- [ ] Contract cleanup
- [ ] Test full loop on vi-agent

**Phase 5: Polish + Release**
- [ ] install.sh integration
- [ ] README.md
- [ ] CHANGELOG.md
- [ ] Test on fresh project (not vi-agent)

### v2.1.0 — Enhancements (backlog)
- [ ] AI-powered PR review (integrated code-reviewer)
- [ ] GitHub Projects board auto-setup
- [ ] Skill freshness checking (check-skills.sh equivalent)
- [ ] Milestone progress in dashboard
- [ ] Team velocity metrics

### v2.2.0 — Advanced (future)
- [ ] Worktree auto-provisioning
- [ ] Notification hooks (Slack/Telegram)
- [ ] Multi-repo support
- [ ] Custom CI template support
- [ ] Mission templates beyond the default

---

## Risk Register

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| `gh` CLI not installed | Medium | Blocks everything | Detect + guide installation on /team init |
| No GitHub repo | Medium | Blocks everything | Detect + error with clear message |
| Issue template not used | Medium | Contract generation fails | Fallback: parse freeform Issue body with AI |
| Actions quota exceeded | Low | CI stops | Local hooks still work as fallback |
| Offline development | Low | Can't read Issues | Contract cached locally, warn on stale |
| Config schema changes | Low | Breaking update | schema_version migration logic |

---

*Design by: architect session, 2026-03-02*
*Validates on: vi-agent project*
*Target: erwin skills/teamwork/ v2.0.0*
