# Teamwork System — Architecture Design (v1, DEPRECATED)

> **DEPRECATED**: This is the v1 architecture. Superseded by [teamwork-v2-architecture.md](teamwork-v2-architecture.md) (GitHub-first redesign, ADR-0010).
> v1 used repo-centric Markdown files. v2 uses GitHub Issues/PRs/Actions as source of truth.

---

> _Original description:_ Team coordination system built on Claude Code skills + GitHub + Markdown files.
> Enables multi-person parallel development with zero-conflict metadata management.

---

## 1. System Overview

**What**: A set of 4 Claude Code skills that coordinate a team of developers working on the same repository. Each person uses Claude Code locally, and the system manages mission assignment, execution, and delivery through Markdown files in git.

**Why**: Enable structured AI-assisted teamwork where each person gets a clear mission contract, executes with /drive methodology, and delivers through a disciplined merge process — all coordinated through the repo itself (no external tools).

**How**: Per-person status files eliminate metadata merge conflicts. The repo IS the system of record.

```
┌─────────────────────────────────────────────────────────────────┐
│                    TEAMWORK SYSTEM MODEL                        │
│                                                                 │
│  ┌─────────────┐         ┌──────────────────┐                  │
│  │   config.md  │────────→│  Identity Layer  │                  │
│  │ (team def)   │         │  gh auth status  │                  │
│  │ Members:     │         │  → lookup role    │                  │
│  │  GH → Role   │         └────────┬─────────┘                  │
│  └─────────────┘                  │                             │
│                                    ▼                             │
│  ┌─────────────┐         ┌──────────────────┐                  │
│  │  missions/   │────────→│ /team-get-mission │                 │
│  │ (contracts)  │         │  Find my mission  │                  │
│  │ Lead writes  │  claim  │  Create branch    │                  │
│  │ Member reads │←────────│  Update status    │                  │
│  └─────────────┘         └────────┬─────────┘                  │
│                                    │                             │
│  ┌─────────────┐         ┌────────▼─────────┐                  │
│  │  members/    │←───────→│   /team-drive     │                  │
│  │ (per-person) │  r/w    │  Load mission     │                  │
│  │ ONLY owner   │         │  → drive mode     │                  │
│  │ edits own    │         │  Update progress  │                  │
│  └─────────────┘         └────────┬─────────┘                  │
│                                    │                             │
│                           ┌────────▼─────────┐                  │
│                           │  /team-complete   │                  │
│                           │  Test + merge     │                  │
│                           │  Mark done        │                  │
│                           └──────────────────┘                  │
│                                                                 │
│  EDIT BOUNDARY (zero-conflict guarantee):                       │
│    Lead edits:   config.md, missions/*.md                       │
│    Member edits: members/{own-name}.md ONLY                     │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Directory Structure

Created in the TARGET PROJECT (not in erwin):

```
target-project/
├── _teamwork/                        ← Teamwork root
│   ├── config.md                     ← Team definition (Lead maintains)
│   ├── missions/                     ← Mission Contracts (Lead writes)
│   │   ├── MISSION-001-modeling.md
│   │   └── MISSION-002-api-auth.md
│   ├── members/                      ← Per-person status (each owns their own)
│   │   ├── qianhua.md
│   │   └── liya.md
│   └── archive/                      ← Completed missions (auto-moved)
│       └── MISSION-000-setup.md
│
├── _project/                         ← (if erwin-init used, coexists)
└── [source code]
```

**Coexistence**: `_teamwork/` is independent of `_project/`. Both can exist in the same repo.

---

## 3. File Format Specifications

### config.md

```markdown
# Teamwork Config

## Team

| GitHub | Role | Name |
|--------|------|------|
| liyasong | tech-lead | Liya |
| qianhua | backend-engineer | 牵华 |
| zhangwei | frontend-engineer | 张伟 |

## Settings

- Branch pattern: `mission/{id}-{slug}-{username}`
- Base branch: `main`
- Test command: (project-specific)
```

### Mission Contract (missions/MISSION-NNN-slug.md)

```markdown
# MISSION-001: Data Modeling Layer

- **Assignee**: qianhua
- **Priority**: P1
- **Created**: 2026-03-03

## Objective

Build ORM data modeling for core entities.

## Sub-tasks

- [ ] Define User, Product, Order models
- [ ] Write Prisma schema
- [ ] Run migrations
- [ ] Seed dev data

## Context

[Detailed context, links to design docs, relevant code paths]

## Acceptance Criteria

- All models have proper relationships defined
- Migrations run cleanly on fresh database
- Seed data populates correctly

## Test

[Agent-based test prompt or test command — project-specific]
```

### Member Status (members/{username}.md)

```markdown
# qianhua — Teamwork Status

- **GitHub**: qianhua
- **Role**: backend-engineer

## Current Mission

- **ID**: MISSION-001
- **Title**: Data Modeling Layer
- **Status**: In Progress
- **Branch**: mission/001-modeling-qianhua
- **Started**: 2026-03-03 09:00

## Progress

- [x] Define User, Product, Order models — 09:30
- [ ] Write Prisma schema ← current
- [ ] Run migrations
- [ ] Seed dev data

## Completed

| Mission | Done | PR |
|---------|------|----|
```

---

## 4. Skill Specifications

### /teamwork — Init + Dashboard

**Purpose**: First-time setup + ongoing team overview.

**Flow**:
1. Detect identity (`gh auth status` → GitHub username)
2. Check if `_teamwork/` exists
   - **No**: Interactive setup — ask for team members + roles, create structure
   - **Yes**: Show dashboard — aggregate all member files into overview
3. Guide to next action

**From erwin-init**: Zero-question directory scaffolding, auto-detection, safety check.

### /team-get-mission — Claim

**Purpose**: Claim a mission, create branch, update status.

**Flow**:
1. Auto-detect identity
2. `git pull` (ensure latest main)
3. Scan `missions/` for missions assigned to me with no "In Progress" status in my member file
4. Show available missions
5. Create branch: `mission/{id}-{slug}-{username}`
6. Update `members/{me}.md` with mission details
7. Commit + push status change
8. Output: next step guidance

### /team-drive — Execute (wraps /drive)

**Purpose**: Execute mission using drive methodology.

**Flow**:
1. Read `members/{me}.md` → current mission ID
2. Read `missions/MISSION-{id}.md` → full contract
3. Enter drive execution mode with mission as context
4. Execute sub-tasks, update progress in member file
5. When done → guide to /team-complete

**Drive integration**: Mission contract = drive's task. Simplified drive phases (skip Phase Ω/T, focus on execution loop).

### /team-complete — Deliver

**Purpose**: Push, rebase, resolve conflicts, merge, close.

**Flow**:
1. Verify all sub-tasks checked
2. Run test hook (if defined)
3. Commit final changes
4. Push branch
5. Rebase on main (`git fetch && git rebase origin/main`)
6. Resolve conflicts (guide user)
7. Re-run tests
8. Push rebased branch
9. Create PR or merge
10. Update status to Done
11. Commit status change

---

## 5. Standard Operating Procedure (SOP)

```
DAILY TEAMWORK SOP
══════════════════

MORNING (5 min):
  $ git checkout main && git pull
  $ claude
  > /teamwork                    ← view dashboard
  > /team-get-mission            ← claim today's mission

WORKING:
  > /team-drive                  ← enter drive mode for mission

COMPLETION:
  > /team-complete               ← push + rebase + merge

MERGE PROTOCOL (Linear):
  First to merge  → few/no conflicts
  Second to merge → resolve vs first
  Third to merge  → resolve vs first + second
  Each merger must verify ALL features work after rebase.
```

---

## 6. Integration with erwin-init

| erwin-init Concept | Teamwork Equivalent | Inherited? | Rationale |
|-------------------|---------------------|------------|-----------|
| Auto-detection (basename) | Identity (gh auth) | Philosophy, not implementation | Same zero-config principle |
| Zero-question design | /teamwork init | Yes | One command creates everything |
| CURRENT.md (status) | members/{name}.md | Pattern, multiplied | Same "living status" concept, per-person |
| Session protocol | SOP | Spirit | "Read before act" → "read mission before work" |
| Git conventions | Branch + commit format | Extended | Added mission/ prefix and team scopes |
| GOALS.md | Mission Contracts | No | Missions replace goals |
| ADRs | — | No | Too heavy for daily tasks |
| MEMORY.md | — | No | Managed by each project's erwin-init |

---

## 7. Key Design Decisions

| # | Decision | Choice | Rationale |
|---|----------|--------|-----------|
| 1 | Concurrency | Per-person files | Zero metadata merge conflicts |
| 2 | Identity | `gh auth status` | Zero-config, reuses existing login |
| 3 | Branching | Skill auto-creates | Uniform naming: `mission/{id}-{slug}-{user}` |
| 4 | Merge strategy | Linear rebase | First in wins, last in resolves all |
| 5 | Testing | Hook only | Project-specific, skill provides the slot |
| 6 | Drive integration | Wrap /drive | Mission contract = drive task input |
| 7 | Status aggregation | Dynamic from member files | No single point of contention |
| 8 | erwin-init integration | Inherit philosophy, not code | Adapted for multi-person context |
| 9 | Config generation | Interactive (AskUserQuestion) | Like erwin-init Phase 2 |
| 10 | Naming | No erwin- prefix | Independent tool, per ADR-0008 |

---

## 8. Risk Register

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Member file edited by wrong person | Low | Low | SOP + code review |
| gh auth not logged in | Medium | Medium | Skill detects + guides `gh auth login` |
| Mission contract lacks detail | Medium | High | /team-drive Phase 0 deep briefs |
| Complex rebase conflicts | Medium | Medium | /team-complete guides step by step |
| 5 people merge simultaneously | Low | High | SOP: sequential merge, not parallel |

---

## 9. Implementation Roadmap

### v1.0.0 — MVP (tonight)
- [x] Architecture document
- [ ] 4 skill files (teamwork, team-get-mission, team-drive, team-complete)
- [ ] install.sh update for multi-file skills
- [ ] skills/README.md update
- [ ] ADR-0009

### v1.1.0 — Enhancements (post-MVP)
- [ ] Dashboard view with progress bars
- [ ] Slack notifications on mission complete
- [ ] Auto-archive completed missions
- [ ] Mission templates for common task types

### v2.0.0 — Scale (future)
- [ ] Multi-repo support
- [ ] Role-based mission filtering
- [ ] Automated conflict detection (pre-merge)
- [ ] Team velocity metrics from archive data

---

*Architecture designed: 2026-03-02 | ADR: 0009 | Session: K*
