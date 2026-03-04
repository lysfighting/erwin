# Teamwork Onboarding Guide

> From zero to first PR in 10 minutes.

---

## What You'll Learn

This guide walks you through the full teamwork cycle:

```
/team (setup) → /team-claim (pick task) → /team-drive (execute) → /team-ship (deliver PR)
```

By the end, you'll have:
- A configured teamwork system in your project
- GitHub labels, Issue template, CI workflow, and git hooks
- One complete mission cycle: Issue → Contract → Code → PR

---

## Prerequisites

Before starting, make sure you have:

| Requirement | Check command | Install |
|-------------|---------------|---------|
| Claude Code | `claude --version` | [claude.ai/code](https://claude.ai/code) |
| GitHub CLI | `gh --version` | [cli.github.com](https://cli.github.com/) |
| Git repo with GitHub remote | `git remote get-url origin` | `git remote add origin <url>` |
| gh authenticated | `gh auth status` | `gh auth login` |
| erwin skills installed | `ls ~/.claude/commands/team.md` | `cd erwin && ./install.sh teamwork` |

---

## Step 1: Initialize — `/team`

Open Claude Code in your project root:

```bash
cd your-project
claude
```

Then run:

```
/team
```

**What happens:**

1. Claude detects your project language (Node, Python, Go, etc.)
2. Asks you about team members and quality preferences
3. Generates:

| File | What it does |
|------|-------------|
| `.teamwork/config.yml` | Team roster, test/lint/build commands, conventions |
| `.github/ISSUE_TEMPLATE/mission.yml` | Structured template for mission Issues |
| `.github/workflows/ci.yml` | CI pipeline: lint → test → build on every PR |
| `.githooks/pre-commit` | Runs linter before each commit |
| `.githooks/pre-push` | Runs tests before each push |

4. Creates GitHub labels: `mission`, `priority:P0`-`P3`, `status:queued`/`wip`/`review`/`done`/`blocked`
5. Commits everything and pushes

**After setup, running `/team` again shows the team dashboard:**

```
TEAM DASHBOARD — your-org/your-project
Milestone: v1.0 (3/10 = 30%)
════════════════════════════════════════════

  alice (tech-lead)
     Working on: #42 Add user authentication [P1]
     PR: #45 — CI passing ✅

  bob (backend)
     Idle (last merged: #38)

QUEUED:
     #47 Add rate limiting [P1]
     #49 Write API docs [P2]

MERGE COUNT:
  alice: 12 | bob: 8

════════════════════════════════════════════
Commands: /team-claim #47 | /team-drive | /team-ship
```

---

## Step 2: Create a Mission Issue — `/team-issue`

Use `/team-issue` to turn a natural language description into a structured GitHub Issue:

```
/team-issue Add user authentication with email + password
```

Claude will:
1. Ask clarifying questions (priority, scope, domain)
2. AI-enrich the Issue body: objective, sub-tasks, success criteria, relevant files
3. Preview the Issue before publishing
4. Create it on GitHub with labels (`mission-contract`, `priority:P1`, `status:queued`)

Or create Issues manually via the GitHub "Mission Contract" template that `/team` installed.

> **Tip**: `/team-issue` requires `versions.current` to be set in config.yml if you want
> the Issue auto-assigned to the current milestone.

---

## Step 3: Claim the Mission — `/team-claim`

```
/team-claim #42
```

Or let Claude auto-select the highest priority unassigned Issue:

```
/team-claim
```

Or browse available missions:

```
/team-claim list
```

**What happens:**

1. Fetches Issue #42 from GitHub
2. Generates a **Mission Contract** at `.teamwork/active/MISSION-42.md`
   - This is NOT just a copy of the Issue — Claude enriches it with:
   - **Context Files**: scans your project to find files relevant to this mission
   - **Test Command**: from config or Issue body
   - **Structured sub-tasks**: broken down from the Issue body
3. Creates branch: `mission/42-add-user-auth-alice`
4. Assigns you on GitHub, adds `status:wip` label, posts claim comment

**Output:**

```
MISSION CLAIMED
═══════════════════════════════════════
Issue:    #42 — Add user authentication
Priority: P1
Branch:   mission/42-add-user-auth-alice
Contract: .teamwork/active/MISSION-42.md

Objective:
  Implement JWT-based authentication for the API

Sub-tasks:
  - [ ] Create auth middleware
  - [ ] Add login/register endpoints
  - [ ] Write tests for auth flow

Context Files:
  src/middleware/index.ts
  src/routes/api.ts
  tests/api.test.ts
═══════════════════════════════════════
Next: /team-drive to start execution
```

---

## Step 4: Execute — `/team-drive`

```
/team-drive
```

**What happens:**

Claude reads your Mission Contract and enters execution mode:

1. **Displays mission briefing** with progress tracker
2. **For each sub-task:**
   - Announces: "Working on: Create auth middleware"
   - Reads the context files from the Contract
   - Implements the change
   - Runs tests to verify
   - Checks off the sub-task in the Contract with a timestamp
   - Commits: `feat(auth): add JWT middleware | Mission: #42`
3. **After all sub-tasks:**
   - Verifies acceptance criteria
   - Runs full test suite
   - Self-reviews all changes
   - Adds execution notes to the Contract

**The Contract tracks progress** — if you interrupt mid-task, the next `/team-drive` picks up where you left off (checkboxes preserve state).

**Output when done:**

```
MISSION EXECUTION COMPLETE
═══════════════════════════════════════
Issue:    #42 — Add user authentication
Status:   All sub-tasks done
Commits:  3 commits on branch mission/42-add-user-auth-alice

Sub-tasks completed:
  [x] Create auth middleware — 14:23
  [x] Add login/register endpoints — 14:35
  [x] Write tests for auth flow — 14:42

Acceptance Criteria:
  JWT tokens work — verified by test suite
  Passwords hashed — verified by bcrypt in code

Tests: passing
═══════════════════════════════════════
Next: /team-ship to create PR and deliver
```

---

## Step 5: Deliver — `/team-ship`

```
/team-ship
```

**What happens:**

1. **Pre-flight checks:**
   - Verifies you're on the mission branch
   - Confirms all sub-tasks checked off
   - Runs tests one more time
2. **Pushes** branch to GitHub
3. **Creates PR** with:
   - Title: `feat: Add user authentication`
   - Body includes: objective, changes list, test results
   - Magic line: `Closes #42` (auto-closes Issue on merge)
4. **Watches CI** (if configured) — reports pass/fail
5. **Requests review** (if configured) — pings the tech-lead
6. **Cleans up** — removes local Contract, updates Issue labels

**Output:**

```
MISSION SHIPPED
═══════════════════════════════════════
Issue:  #42 — Add user authentication
Branch: mission/42-add-user-auth-alice
PR:     https://github.com/your-org/your-project/pull/45
CI:     passing ✅
Review: requested from alice

Contract: cleaned up ✅
Issue:    labeled "status:review" ✅
═══════════════════════════════════════
Next steps:
  - After merge: /team-ship done (close Issue, update labels)
  - To AI-review PR: /team-ship review
  - To claim next mission: /team-claim
  - To see team status: /team
```

---

## The Daily Workflow

Once set up, the daily cycle is:

```
Morning:    git pull → claude → /team (see dashboard) → /team-claim (pick task)
Working:    /team-drive (execute with AI assistance)
Done:       /team-ship (deliver PR)
After merge:/team-ship done (close Issue, clean up)
Repeat:     /team-claim (next task)
Release:    /team-release (when milestone complete → tag + publish + next milestone)
```

---

## How It All Connects

```
GitHub Issue (#42)
    │
    ├── /team-claim ──→ Mission Contract (.teamwork/active/MISSION-42.md)
    │                     AI-enriched: context files, test commands, structured tasks
    │
    ├── /team-drive ──→ Code changes on branch mission/42-*
    │                     Each sub-task: implement → verify → commit
    │
    ├── /team-ship ───→ Pull Request (Closes #42)
    │                     Push → PR → CI check → review request → cleanup
    │
    ├── /team-ship done → Issue #42 closed, labels updated, worktree cleaned
    │
    └── Merge PR ─────→ Issue #42 auto-closed (if not already)
                         Branch auto-deleted (GitHub setting)
```

**Source of truth**: GitHub (Issues, PRs, CI)
**AI execution layer**: Mission Contract (ephemeral, enriched, local)
**Quality gates**: git hooks (local) + GitHub Actions (CI) + branch protection (review)

---

## Team Setup (Multiple People)

When running `/team` for the first time, select "2-3 people" (or more) and enter GitHub usernames.

Each team member then:
1. `cd your-project && git pull`
2. `claude` → `/team` (sees the dashboard)
3. `/team-claim` (picks an unassigned Issue)
4. `/team-drive` → `/team-ship`

**Rules:**
- One active mission per person (enforced by `/team-claim`)
- Issues are claimed via GitHub assignment (visible to the whole team)
- PRs link to Issues via `Closes #N` (GitHub tracks delivery)
- Dashboard shows who's doing what, in real time from GitHub data

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| "Not authenticated" | Run `gh auth login` |
| "No git remote" | `git remote add origin <github-url>` |
| "Teamwork not initialized" | Run `/team` first |
| "Already have an active mission" | Finish current: `/team-ship`, or abandon: `rm .teamwork/active/MISSION-*.md` |
| "Tests failing" | Fix tests before `/team-ship` — quality gate is intentional |
| "Branch already exists" | Previous attempt — `/team-claim` will switch to it |
| "CI failing after PR" | Push fixes to the same branch — PR updates automatically |
| Label creation fails | You may not have admin rights — non-fatal, continue |
| Branch protection fails | Free GitHub plans may not support this — non-fatal |

---

## Config Reference

After `/team` init, your `.teamwork/config.yml` looks like:

```yaml
schema_version: 1
skill_version: 2.6.0

team:
  - github: alice
    role: tech-lead
  - github: bob
    role: backend

project:
  language: typescript
  test_command: "npm test"
  lint_command: "npm run lint"
  build_command: "npm run build"

conventions:
  branch_pattern: "mission/{issue}-{slug}-{user}"
  base_branch: main
  # release_branch: main  # optional: set when using dual-branch mode (base ≠ release)
  commit_format: "type(scope): description | Mission: #{issue}"

label_prefix:
  status: "status:"
  priority: "priority:"
  size: "size:"

quality:
  hooks: true
  ci: true
  review_required: true
  branch_protection: true

# Optional features:
# worktree:
#   enabled: true
#   path_pattern: "../{repo}-wt-{slug}"
# versions:
#   current: "V1.0"
#   spec_path: ".claude/drive/v1.0-definition/spec.md"  # enables AI audit in /team-release
#   lifecycle: [dev, qa, released]
```

Edit this file to change team members, commands, or quality preferences. Changes take effect immediately (no re-init needed).

---

## Step 6: Release a Version — `/team-release`

When all missions in the current milestone are done:

```
/team-release
```

Flow:
1. Checks milestone `V1.0` — shows open/closed Issue counts
2. If `versions.spec_path` is configured, AI audits spec goals vs merged PRs (✅/⚠️/❌ per criterion)
3. User confirmation
4. `git tag V1.0` + push
5. `gh release create V1.0 --generate-notes` — auto-generates release notes from PRs
6. Closes GitHub Milestone `V1.0`
7. Prompts for next version → creates Milestone `V1.1` → updates `versions.current` in config

**Prerequisites:**
- `versions.current: "V1.0"` in config.yml
- GitHub Milestone `V1.0` must exist (`gh api repos/{REPO}/milestones --method POST --field title="V1.0"`)

**With `spec_path` configured:**

```yaml
versions:
  current: "V1.0"
  spec_path: ".claude/drive/v1.0-definition/spec.md"
```

The spec file should have a `Success Criteria` section listing verifiable goals. `/team-release` reads these and checks each against merged PRs — giving you an audit before tagging.
