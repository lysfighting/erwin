# Teamwork — The Complete Manual

> **Version 3.9.6** | Author: liyasong + casey | 2026-03-12
>
> AI-Native Team Coordination for GitHub.
> For humans and AI agents. Covers WHY, HOW, and the philosophy behind every decision.

---

## Table of Contents

1. [Design Philosophy](#design-philosophy)
2. [The Problem We Solve](#the-problem-we-solve)
3. [System Architecture](#system-architecture)
4. [The Four Branches](#the-four-branches)
5. [The Six Skills](#the-six-skills)
6. [Issue Lifecycle — State Machine](#issue-lifecycle--state-machine)
7. [Git Flow — Visual Guide](#git-flow--visual-guide)
8. [RC Lifecycle — The Key Innovation](#rc-lifecycle--the-key-innovation)
9. [GitHub Ecosystem Integration](#github-ecosystem-integration)
10. [The Mission Contract](#the-mission-contract)
11. [Quality Gates — Defense in Depth](#quality-gates--defense-in-depth)
12. [Versioning Model](#versioning-model)
13. [Scripts Architecture](#scripts-architecture)
14. [Config Reference](#config-reference)
15. [End-to-End Example](#end-to-end-example)
16. [Branch Protection](#branch-protection)
17. [Closing Thoughts](#closing-thoughts)

---

## Design Philosophy

### Nine Principles

These are not rules we invented. They emerged from **pain** — bugs found, rework done, production incidents caused by AI-generated noise. Each principle has a scar behind it.

```
┌─────────────────────────────────────────────────────────────────────┐
│                    DESIGN PRINCIPLES                                │
│                                                                     │
│  ① GitHub IS the system                                            │
│     Not a mirror. Not a sync target. THE system.                   │
│     Issues = tasks. PRs = delivery. Actions = gates.               │
│     Cross-agent, cross-session, cross-human communication          │
│     happens via GitHub — same as any human team.                   │
│                                                                     │
│  ② Separate the dirty from the clean                               │
│     develop absorbs AI noise. main stays pristine.                 │
│     Cleanliness is a RESULT (of verification), not a starting pt.  │
│                                                                     │
│  ③ Squash everything                                               │
│     AI commits are machine noise. 1 Issue = 1 squash commit.      │
│     Enables cherry-pick recovery. Enables readable history.        │
│                                                                     │
│  ④ Risk levels don't share trigger paths                           │
│     staging deploy ≠ prod deploy. Never in the same command.       │
│     Skills own git ops. Actions own deployment.                    │
│                                                                     │
│  ⑤ Tags are semantic markers, not deployment logs                  │
│     V0.1.0 means "this version exists." Not "deploy this now."    │
│     workflow_dispatch for deploy. Tags for history.                │
│                                                                     │
│  ⑥ Scripts lock mechanism, LLM provides judgment                  │
│     Deterministic ops (git, labels, PR) = bash scripts.            │
│     Adaptive ops (code analysis, review) = LLM.                   │
│     Evidence: 21 bugs found in LLM-interpreted ops. 0 in scripts. │
│                                                                     │
│  ⑦ Close the loop first, optimize second                          │
│     A complete loop with rough edges > a polished partial flow.    │
│     Automation fills gaps humans forget.                           │
│                                                                     │
│  ⑧ Config-driven, not hardcoded                                   │
│     Every project-specific value lives in config.yml.              │
│     The skill works on any repo. Zero coupling to vi_agent.        │
│                                                                     │
│  ⑨ 6 skills is the ceiling                                        │
│     Expand via subcommands, not new skills.                        │
│     30+ commands confused users (evidence: CCPM project).          │
│     6 commands = one hand + one thumb.                             │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### The Von Neumann Separation

```
    ┌──────────────────────────────────────────────────┐
    │              skill.md (LLM Layer)                 │
    │                                                    │
    │  "Read the Issue, understand context, decide       │
    │   what type of branch to create, generate          │
    │   a Mission Contract with AI-enriched context"     │
    │                                                    │
    │         ↓ calls ↓           ↑ reads result ↑       │
    ├──────────────────────────────────────────────────┤
    │              scripts/*.sh (Mechanism Layer)        │
    │                                                    │
    │  tw-git.sh create-branch 42 user-login lysfighting │
    │  tw-label.sh transition 42 queued wip              │
    │  tw-contract.sh hash 42                            │
    │  tw-pr.sh create --squash --base develop           │
    │  tw-notify.sh mc.created --issue 42                │
    │                                                    │
    │  Deterministic. Testable. Zero LLM interpretation. │
    └──────────────────────────────────────────────────┘
```

> *"Content changes; mechanism endures."*
> — Inspired by von Neumann's stored-program architecture

---

## The Problem We Solve

Traditional teams: 6 people, manageable commit frequency, GitHub Flow works fine.

**Our reality**: 6 people + AI assistants = **20-30x commit frequency**. This changes everything.

```
THE STAGING PARADOX (GitHub Flow + AI)
═══════════════════════════════════════════════════════════════

main ── feat/A ── feat/B ── feat/C ── feat/D ── feat/E ──►
        (merge)   (merge)   (merge)   (merge)   (merge)

        ↑ you deploy staging here
        │ while verifying, D and E merge in
        │ deploy again → F and G merge in
        └── you can NEVER verify a stable snapshot

Problem: the target moves faster than you can aim.
```

**The insight**: you need a **frozen snapshot** for verification. But you also need a place where AI can commit freely without polluting production.

```
THE SOLUTION: Separation of Concerns
═══════════════════════════════════════════════════════════════

develop ── A ── B ── C ── D ── E ──►     ← AI noise goes here
                        │                    (that's OK)
                   rc/V0.1.0                ← frozen snapshot
                        │                    for verification
main ────────────── [V0.1.0] ──►            ← only verified code
                                              forever clean
```

> *This is not Git Flow reinvented. It's Git Flow understood through the lens of AI-era commit frequency.*

---

## System Architecture

### The Complete Picture

```
╔══════════════════════════════════════════════════════════════════════╗
║                    TEAMWORK v3 SYSTEM ARCHITECTURE                   ║
╠══════════════════════════════════════════════════════════════════════╣
║                                                                      ║
║  ┌─────────────────────────────────────────────────────────────┐    ║
║  │                   6 SKILLS (LLM Layer)                       │    ║
║  │                                                               │    ║
║  │  /team ──── /team-issue ──── /team-claim ──── /team-drive    │    ║
║  │    │              │               │               │           │    ║
║  │    │              │               │               │           │    ║
║  │    │         /team-ship ──────────────────── /team-rc         │    ║
║  │    │              │                              │            │    ║
║  │  init           create(push)   claim(my list)  prepare       │    ║
║  │  dashboard      comment        branch          promote       │    ║
║  │  learn          update         contract                      │    ║
║  │  queue          fix(untracked)                               │    ║
║  │  #N detail      batch / #N                                   │    ║
║  └───────────────────────┬─────────────────────────────────────┘    ║
║                          │ calls                                     ║
║  ┌───────────────────────▼─────────────────────────────────────┐    ║
║  │                8 SCRIPTS (Mechanism Layer)                    │    ║
║  │                                                               │    ║
║  │  tw-git.sh ─── tw-label.sh ─── tw-contract.sh ─── tw-pr.sh  │    ║
║  │  (15 cmds)     (4 cmds)        (8 cmds)           (7 cmds)  │    ║
║  │                                                               │    ║
║  │  tw-config.sh (YAML parser)  tw-notify.sh (Slack/Feishu)    │    ║
║  │  setup-github-labels.sh      tw-e2e-test.sh (validation)    │    ║
║  └───────────────────────┬─────────────────────────────────────┘    ║
║                          │ operates on                               ║
║  ┌───────────────────────▼─────────────────────────────────────┐    ║
║  │                  GITHUB (Source of Truth)                      │    ║
║  │                                                               │    ║
║  │  Issues ─── Labels ─── Milestones ─── PRs ─── Actions        │    ║
║  │  (tasks)    (state)    (versions)     (delivery) (quality)    │    ║
║  │                                                               │    ║
║  │  Branch Protection ─── Repo Settings ─── Releases ─── Tags   │    ║
║  │  (safety)               (merge policy)    (history)   (anchor)│    ║
║  └──────────────────────────────────────────────────────────────┘    ║
║                                                                      ║
║  ┌──────────────────────────────────────────────────────────────┐    ║
║  │                    LOCAL (Ephemeral)                            │    ║
║  │                                                               │    ║
║  │  .teamwork/config.yml ─── .teamwork/active/MISSION-*.md      │    ║
║  │  (permanent config)        (ephemeral Contract, gitignored)   │    ║
║  │                                                               │    ║
║  │  .githooks/pre-commit ─── .githooks/pre-push                 │    ║
║  │  (lint gate)               (test gate)                        │    ║
║  └──────────────────────────────────────────────────────────────┘    ║
║                                                                      ║
╚══════════════════════════════════════════════════════════════════════╝
```

### Responsibility Boundary

```
         SKILL LAYER                    ACTION LAYER
     (Claude executes)              (GitHub executes)
    ─────────────────────       ─────────────────────────
    /team-issue → Issue          on PR merge → post-merge
      + assign + branch           cleanup (labels, boxes,
    /team-claim → Contract         milestone check)
    /team-drive → Code
    /team-ship  → PR             on workflow_dispatch →
    /team-rc    → Tag              deploy staging/prod
         │                              │
         │  checks before tagging       │ provides the
         └────── Action status ─────────┘ check basis
              success? → proceed
              failed?  → warn + confirm

    Skills own git operations.
    Actions own deployment.
    Never cross the boundary.
```

### Push Model (v3)

```
    v2 (Pull Model):                 v3 (Push Model):
    ─────────────────                ─────────────────────
    /team-issue creates              /team-issue creates
      status:queued                    + assigns @member
      unassigned                       + status:wip
                                       + creates branch
    Dev self-claims via
      /team-claim                    Assignee runs
      → assigns to self               /team-claim #N
      → status:wip                     → generates Contract
      → creates branch                 (branch already exists)

    "Pull": dev picks from queue     "Push": creator assigns directly
    Good for solo/small teams        Good for coordinated teams
```

Any team member can create Issues and assign to anyone — no role restriction.
Solo projects (1 member) auto-assign to self.
Roles (leader/member) are used only for **dashboard view** distinction, not permissions.

---

## The Four Branches

```
┌──────────┬───────────────┬─────────────────────────┬────────────────┐
│ Branch   │ Role          │ Who writes to it        │ Cleanliness    │
├──────────┼───────────────┼─────────────────────────┼────────────────┤
│ main     │ Production    │ Only rc/* via PR        │ ★★★ Pristine  │
│ develop  │ Dev trunk     │ All feature PRs         │ ★★  CI-clean  │
│ rc/*     │ RC / Staging  │ Hotfixes only           │ ★★½ Converging│
│ feat/*   │ Work          │ Individual developer    │ ★   Dirty     │
└──────────┴───────────────┴─────────────────────────┴────────────────┘
```

### Branch Flow Diagram

```
    feat/42-login ──(squash)──►
    feat/45-perms ──(squash)──►
    feat/48-api   ──(squash)──►
                               │
develop ──●──●──●──●──●──●─────●──●──●──●──●──►
                         │                │
                    rc/V0.1.0        rc/V0.1.1
                         │                │
                    (staging)         (staging)
                         │                │
                    (squash PR)      (squash PR)
                         │                │
main ────────────────[V0.1.0]────────[V0.1.1]──────►
                       tag             tag
                    (release)       (release)
```

### Why branch from develop, not main?

```
WRONG: Branch from main
═══════════════════════════════════════════

main:     V0.1.0 ──────────────────────────────────
develop:  V0.1.0 ── A ── B ── C ── D ──
feat/E:   V0.1.0 ── (your work) ──
                ↑                    ↓
                │               merge to develop
                │                    ↓
                └── CONFLICT! feat/E is based on old V0.1.0,
                    but develop now has A, B, C, D changes


RIGHT: Branch from develop
═══════════════════════════════════════════

develop:  V0.1.0 ── A ── B ── C ── D ──
feat/E:                          D ── (your work) ──
                                      ↓
                                 merge to develop
                                      ↓
                                 clean merge ✓ (same base)
```

> **Rule: always branch from develop.**
> Main's cleanliness is a *result* (of RC verification), not a *starting point*.

---

## The Six Skills

```
┌──────────────────────────────────────────────────────────────────┐
│                     THE SIX SKILLS (v3 — Push Model)              │
│                                                                    │
│     /team              Hub — Dashboard + Init + Doctor + Help      │
│       │                     + Learn + Queue + Wrap + Auto          │
│       ▼                                                            │
│     /team-issue ──► Create + assign (Push)                        │
│       │              Smart routing: create / comment / update /    │
│       │              fix / batch / detail / help                   │
│       ▼                                                            │
│     /team-claim ──► Generate Contract + branch from assigned Issue │
│       │              (creates mission branch + AI execution plan)  │
│       ▼                                                            │
│     /team-drive ──► Execute mission from Contract                 │
│       │              (sub-tasks → test → commit loop)              │
│       │              Service-aware testing (monorepo support)      │
│       ▼                                                            │
│     /team-ship  ──► Push → PR → notify                            │
│       │              (auto-closes Issue on merge)                  │
│       ▼                                                            │
│     /team-rc    ──► Prepare RC / Promote to production            │
│                      (cut rc branch / squash → main → tag)        │
│                                                                    │
│  SHORTCUTS:                                                        │
│     /team auto #N ── claim → drive → ship in one command          │
│     /team wrap    ── code-first → retroactive issue + PR          │
│                                                                    │
│  Each skill is a PHASE in the lifecycle.                          │
│  Together they form a complete loop.                              │
│  The loop never breaks.                                           │
└──────────────────────────────────────────────────────────────────┘
```

### Quick Reference

| Skill | Command | What it does |
|-------|---------|--------------|
| `/team` | `/team` | Dashboard — who's working on what (role-based view) |
| | `/team init` | Initialize teamwork in current repo |
| | `/team help` | Quick start guide |
| | `/team learn` | Design philosophy + diagrams (14 sections) |
| | `/team config` | View current config + edit instructions |
| | `/team queue` | All open missions sorted by priority |
| | `/team #42` | MC detail view (branch, commits, PRs, criteria) |
| | `/team doctor` | Local git + issue health diagnostics (12 checks, 3-layer report) |
| | `/team doctor fix` | Interactive fix mode — auto-execute cleanup, confirm warnings |
| | `/team doctor help` | Self-contained doctor usage guide |
| | `/team wrap` | Retroactive closure: code-first → auto-create issue + PR |
| | `/team auto #42` | One-shot: claim → drive → ship in one command |
| | `/team auto` | Pick from assigned queue, then auto |
| `/team-issue` | `/team-issue {desc}` | Create Issue + assign (Push) |
| | `/team-issue {desc} @user` | Create and assign to specific member |
| | `/team-issue #42` | View Issue details |
| | `/team-issue #42 {text}` | Add comment to Issue |
| | `/team-issue update #42` | AI-assisted edit of Issue body |
| | `/team-issue fix #42` | Fix untracked Issue (add teamwork labels) |
| | `/team-issue batch V0.2 -- ...` | Batch create milestone + MCs |
| `/team-claim` | `/team-claim #42` | Generate Contract + branch for assigned Issue |
| | `/team-claim` | List MY assigned Issues |
| `/team-drive` | `/team-drive` | Execute mission (code → test → commit) |
| `/team-ship` | `/team-ship` | Push + create PR + notify |
| | `/team-ship done` | Post-merge cleanup |
| | `/team-ship review` | AI code review on PR |
| | `/team-ship sync` | Rebase on latest develop + auto-resolve conflicts |
| `/team` | `/team sync` | Rebase on base branch (standalone, anytime) |
| | `/team share [#N]` | Cross-MC knowledge sharing via issue comment |
| `/team-rc` | `/team-rc` | Prepare: cut rc from develop |
| | `/team-rc cancel` | Cancel: delete RC, fix on base, re-cut |
| | `/team-rc fix` | Cherry-pick fix to active RC (keep test progress) |
| | `/team-rc hotfix` | Post-release: fix on main + backport to develop |
| | `/team-rc promote` | Promote: squash to main, tag, release, notify |
| | `/team-rc help` | Full RC lifecycle guide |

### Doctor — Health Diagnostics (v3.4+)

`/team doctor` runs 12 diagnostic checks and produces a 3-layer report:

| Check | What it detects | Severity |
|-------|----------------|----------|
| 1. Orphan Branches | Local branches with closed Issues / merged PRs | 🟢 cleanup |
| 2. Cross-Contamination | Uncommitted changes belonging to different MC | 🔴 critical |
| 3. Branch Staleness | Mission branches >30 commits behind base | 🟡 warning |
| 4. Stale Stashes | Stashes >7 days old or on deleted branches | 🟡 warning |
| 5. Worktree Health | Prunable, orphan, or broken worktrees | 🟡 warning |
| 6. Shippable PRs | PRs with CI passing + mergeable but idle | 🟡 warning |
| 7. Issue Health | Duplicates, unassigned, stale WIP issues | 🟡 warning |
| 8. Remote Sync | Local branch behind remote (needs pull/push) | 🟡 warning |
| 9. GitHub Auth | `gh` CLI auth status and token scopes | 🔴 critical |
| 10. Contract Freshness | Active Contract hash vs current Issue content | 🟡 warning |
| 11. Test Suite | Run project test command, report failures | 🟡 warning |
| 12. PR CI Status | CI check details on open PRs (extends Check 6) | 🟡 warning |

**Output**: Layer 1 (snapshot tables) → Layer 2 (diagnosis narrative) → Layer 3 (action plan with exact commands).

**`/team doctor fix`**: Interactive fix mode. Cleanup actions auto-execute; warnings confirm each via prompt; critical issues confirm with risk explanation.

**`/team doctor help`**: Self-contained usage guide showing all checks and typical workflow.

### Operational Details (v3.5+)

**Base branch guard** (`/team-ship`): Before shipping, verifies you're NOT on a protected branch (main/develop). Prevents accidental ship from base branch.

**CI-preserves-contract** (`/team-ship`): When CI fails after PR creation, the Mission Contract is preserved (not deleted). This allows iterating with `/team-drive` → `/team-ship` until CI passes.

**Contract-missing fallback** (`/team-drive`): If Contract is missing but you're on a mission branch with an open Issue, suggests `/team-claim #N` to regenerate the Contract instead of a bare error.

**NEEDS CLEANUP dashboard section**: Appears when local mission branches have merged PRs. Shows branch name + PR number. Fix: `/team-ship done`.

### Visual Encoding

> Full spec: `docs/visual-encoding-standard.md`

All teamwork output follows a **color design language** built on Claude Code's markdown renderer constraints.

**Markdown formatting** (4 semantic roles):

| Syntax | Color Result | Used For |
|--------|-------------|----------|
| `**bold**` | White bold | Section headers, error prefixes, labels |
| `` `backtick` `` | Purple-blue rgb(177,185,249) | Commands, paths, counts, config values |
| `*italic*` | Dim italic | Branch names, secondary info, timestamps |
| `**#NNN**` | Light-blue clickable (OSC 8) | Issue/PR numbers (bold wraps auto-linker) |

**Rules**: `────` for dividers (not `═══` or `---`). No `###` headers (causes 2-line gap). No code blocks around output (kills all color). Issue numbers zero-padded to 3+ digits (`#042` not `#42`).

**Emoji conventions**:

| Symbol | Meaning | Context |
|--------|---------|---------|
| 🔴 | P0 critical / CI fail / critical diagnostic | Priority, CI, doctor |
| 🟠 | P1 high | Priority |
| 🟡 | P2 medium / under review / warning | Priority, status, doctor |
| ⚪ | P3 low / idle | Priority, member status |
| 🟢 | Active / WIP / passing / cleanup | Status, CI, doctor |
| ✅ | Done / CI passing | Completion indicators |
| ❌ | Failed / CI failing | Error indicators |
| ⚡ | Issue updated since claim | Dashboard stale warning |
| 🔀 | PR / merge related | Pull requests |
| 💀 | Orphan branch | Doctor output |
| ⏳ | CI pending | CI status |
| 🔧 | Manual / human action | Tasks AI skips |

**Progress bar**: `█` filled, `░` empty (10 chars). E.g., 50% → `█████░░░░░`.

### Shortcuts — Wrap & Auto (v3.3+)

**`/team wrap`** — Retroactive closure. You already wrote code, now back-fill the teamwork flow.

Auto-detects your code state and picks the right strategy:

| State | What happened | Strategy |
|-------|--------------|----------|
| Uncommitted changes | You modified files but didn't commit | Stash → create Issue → claim branch → apply stash → drive → ship |
| Committed, not pushed | You committed on current branch | Create Issue → claim branch → cherry-pick → reset original → drive → ship |
| Already pushed | Code is on remote | Create Issue with `status:done` → close immediately (tracking only) |

```
/team wrap
# Auto-detects state, generates issue title/priority from diff analysis
# Previews plan → confirms → executes
```

**When to use**: You pair-programmed, prototyped, or fixed a bug without `/team-claim` first. Wrap retroactively gives it a proper Issue, branch, and PR.

**`/team auto`** — One-shot forward closure. Claim → drive → ship in a single command.

```
/team auto #42     # Specific issue
/team auto         # Pick from your assigned queue
```

**Flow**: Validates Issue → creates Contract + branch → executes all sub-tasks → runs tests → pushes → creates PR. Stops at PR creation — human merges.

**When to use**: Small, well-defined tasks where the standard 4-step cycle (`claim → drive → ship → done`) feels like ceremony. Auto compresses it to one command.

**Differences from standard flow**:
- No interactive confirmations (assumes defaults)
- Blocked sub-tasks are skipped with a note (asks "Ship partial?" at end)
- Test failures trigger 3 auto-fix attempts before stopping

---

## Issue Lifecycle — State Machine

### Label State Transitions (v3 Push Model)

```
                    ┌─────────────────────────┐
                    │      ISSUE CREATED       │
                    │    /team-issue creates    │
                    │    + assigns @member      │
                    │    + creates branch       │
                    └────────────┬────────────┘
                                 │
                                 ▼
                    ┌─────────────────────────┐
                    │    status:wip            │
                    │    ● Assigned to member  │
                    │    ● Branch exists       │
                    │    ● Ready for claim     │
                    └────────────┬────────────┘
                           │     │
                    /team-ship   │ (blocked)
                           │     │
                           ▼     ▼
    ┌──────────────────────┐   ┌──────────────────────┐
    │    status:review     │   │    status:blocked     │
    │    ● PR is open      │   │    ● Waiting on       │
    │    ● CI running      │   │      dependency       │
    │    ● Awaiting merge  │   │    ● Manual set       │
    └──────────┬───────────┘   └──────────────────────┘
               │
        PR merged
     (post-merge Action
      OR /team-ship done)
               │
               ▼
    ┌──────────────────────┐
    │    status:done        │
    │    ● All boxes checked│
    │    ● Issue closed     │
    │    ● Branch deleted   │
    └──────────────────────┘
```

> **v3 change**: `status:queued` is used for unassigned Issues (via `/team-issue fix` or unassigned creation). Assigned Issues start as `status:wip` directly.

### Untracked Issues

Issues created or assigned directly on GitHub (not through `/team-issue`) lack teamwork labels. The `/team` dashboard detects these by comparing "all assigned issues" vs "mission-labeled issues" — the difference is **untracked**.

```
  GitHub assign #94 → @lysfighting
  /team dashboard   → UNTRACKED section shows #94
  /team-issue fix #94 → adds mission + status:wip + priority + size labels
  /team dashboard   → #94 now appears in MY MISSIONS
```

This ensures no assigned work is invisible, regardless of how the issue was created.

### Label Colors (visual recognition)

```
    ┌──────────────┬────────────┬──────────┐
    │ Label        │ Color      │ Meaning  │
    ├──────────────┼────────────┼──────────┤
    │ status:wip   │ 🟡 #fbca04 │ Active   │
    │ status:review│ 🟣 #7057ff │ Waiting  │
    │ status:done  │ 🟢 #0e8a16 │ Complete │
    │ status:blocked│ 🔴 #d73a4a│ Stuck    │
    ├──────────────┼────────────┼──────────┤
    │ priority:P0  │ 🔴 #d73a4a │ Critical │
    │ priority:P1  │ 🟡 #e4e669 │ High     │
    │ priority:P2  │ 🟢 #0e8a16 │ Medium   │
    │ priority:P3  │ ⚪ #cfd3d7 │ Low      │
    └──────────────┴────────────┴──────────┘
```

---

## Git Flow — Visual Guide

### The Complete Git Lifecycle

```
══════════════════════════════════════════════════════════════════════
                        GIT FLOW — FULL PICTURE
══════════════════════════════════════════════════════════════════════

 feat/42 ─●─●─●─┐ (squash)
                 │
 feat/45 ─●─●───┤ (squash)
                 │
 feat/48 ─●─────┤ (squash)
                 │
develop ─────●───●───●───●───●───────────────cherry-pick(H)──●──►
                             │                     ▲
                        rc/V0.1.0                  │
                             │              hotfix/fix-bug
                             │                (squash)
                          staging                  │
                          verify ─── bug? ── fix ──┘
                             │
                          all good
                             │
                        squash PR ──────────────────┐
                             │                       │
main ────────────────────────┼───[release: V0.1.0]──●──────────►
                             │      ↑
                          delete    tag V0.1.0
                         rc branch  ↑
                                    gh release create
                                    --generate-notes

══════════════════════════════════════════════════════════════════════
```

### Why Squash Everything

AI generates dozens of commits per feature. These are machine noise:

```
WITHOUT SQUASH (main history):           WITH SQUASH (main history):

  a1b2c3 fix typo                         x1y2z3 feat: user login (#42)
  d4e5f6 WIP                              a2b3c4 feat: permissions (#45)
  g7h8i9 try different approach            d5e6f7 fix: API timeout (#48)
  j0k1l2 revert previous                  g8h9i0 release: V0.1.0
  m3n4o5 actually fix it
  p6q7r8 forgot to save                   4 commits. Each one meaningful.
  s9t0u1 final fix for real               Each one cherry-pickable.
  ...                                     Each one a complete unit.
  47 commits. None meaningful.
```

### Merge Method Summary

```
┌───────────────────────┬──────────┬───────────────────────────────┐
│ Operation             │ Method   │ Why                           │
├───────────────────────┼──────────┼───────────────────────────────┤
│ feature → develop     │ squash   │ Compress AI noise             │
│ hotfix → rc/*         │ squash   │ Same                          │
│ rc/* → main           │ squash   │ 1 release = 1 commit on main  │
├───────────────────────┼──────────┼───────────────────────────────┤
│ NEVER allowed:        │          │                               │
│ direct push to develop│ blocked  │ Branch protection             │
│ direct push to main   │ blocked  │ Branch protection             │
│ force push anywhere   │ blocked  │ Branch protection             │
└───────────────────────┴──────────┴───────────────────────────────┘
```

---

## RC Lifecycle — The Key Innovation

### Why "RC" and not "release branch"?

```
"release branch" implies it IS the release.
But it's NOT — it's a CANDIDATE being verified.

rc/V0.1.0 on staging  →  might pass, might fail
                          it's auditioning for the role

tag V0.1.0 on main     →  THIS is the release
                          irreversible, permanent
```

> *The name matters. Words shape behavior.*
> *"Release branch" makes people treat it as final.*
> *"Release Candidate" keeps the right level of skepticism.*

### The RC Bus Analogy

```
                    ┌───────────────────────────────────────┐
                    │          THE BUS ANALOGY               │
                    │                                        │
                    │  RC is a bus. It departs on schedule.  │
                    │  If your feature isn't on this bus,    │
                    │  it takes the next one.                │
                    │                                        │
                    │  Bus schedule:                         │
                    │                                        │
                    │  rc/V0.1.0 ══► staging ══► main       │
                    │       ▲                                │
                    │       │ features A, B, C on board      │
                    │       │ feature D missed it            │
                    │       │                                │
                    │  rc/V0.1.1 ══► staging ══► main       │
                    │       ▲                                │
                    │       │ features D, E, F on board      │
                    │                                        │
                    │  Buses come frequently.                │
                    │  Missing one is not a crisis.          │
                    │  It's Tuesday.                         │
                    └───────────────────────────────────────┘
```

### RC Lifecycle — Step by Step

```
══════════════════════════════════════════════════════════════════════
                    RC LIFECYCLE — COMPLETE FLOW
══════════════════════════════════════════════════════════════════════

  PHASE 1: PREPARE (/team-rc)
  ─────────────────────────────────────────────

  ① Check: no existing rc/* branch on remote
     (only one RC at a time — enforced)

  ② Derive next version from remote tags
     V0.1.0 exists → next is V0.1.1
     No tags yet   → V0.1.0

  ③ Verify develop CI is green

  ④ Cut branch:
     develop ─────────── cut ──── rc/V0.1.1
                                    │
  ⑤ Deploy to staging              │ (optional, if
     (workflow_dispatch)            │  staging_workflow
                                    │  is configured)

  PHASE 2: VERIFY (manual / AI QA)
  ─────────────────────────────────────────────

  ⑥ Test on staging environment

  ⑦ Bug found? Two options:
     OPTION A — UPSTREAM-FIRST (default, clean slate):
     └── fix on develop (normal mission flow)
     └── /team-rc cancel (delete stale RC)
     └── /team-rc (re-cut from latest develop)
     └── re-deploy staging, re-verify

     OPTION B — CHERRY-PICK (surgical, keep test progress):
     └── fix on develop first
     └── /team-rc fix <commit> (cherry-pick to RC)
     └── re-deploy staging, re-verify

  ⑧ All good → proceed to promote

  PHASE 3: PROMOTE (/team-rc promote)
  ─────────────────────────────────────────────

  ⑨  Create PR: rc/V0.1.1 → main

  ⑩  Squash merge via GitHub
      (tries --auto first, waits for CI)
      (falls back to manual merge if auto-merge disabled)

  ⑪  Tag V0.1.1 on main (the squash commit)

  ⑫  Create GitHub Release (--generate-notes)

  ⑬  Check milestone: all issues closed?
      → yes: close milestone automatically
      → no:  report "X open, Y closed"

  ⑭  Delete rc/V0.1.1 branch (local + remote)

  ⑮  Notify team: tw-notify.sh milestone.done

  PHASE 4: HOTFIX (/team-rc hotfix) — post-release
  ─────────────────────────────────────────────

  ⑯  Create hotfix/{version} from main
  ⑰  Make fix (interactive)
  ⑱  PR to main (squash merge)
  ⑲  Tag hotfix version on main
  ⑳  Cherry-pick to develop (backport)
  ㉑  Delete hotfix branch

  (Requires dual-branch mode: base ≠ production)

══════════════════════════════════════════════════════════════════════
```

### RC Freeze Rule — Upstream-First Model

```
  RC is a FROZEN SNAPSHOT. No direct git push to RC.
  Default: fix on develop first, cancel + re-cut (upstream-first).
  Surgical: /team-rc fix cherry-picks a specific commit to RC.

  develop ── A ── B ── C ── [fix D] ── E ──►
                        │                │
                   rc/V0.1.0        rc/V0.1.0  (re-cut with fix D)
                   (has bug)        (/team-rc cancel ↗ /team-rc)
                        ✗                │
  main ──────────────────────── [V0.1.0] ──►

  Why upstream-first?
  1. No merge-back problem — develop always has all fixes
  2. RC has no unique commits — squash merge to main is clean
  3. Re-cut is cheap (one command) — cherry-pick is error-prone
```

### Only One RC At A Time

```
/team-rc checks remote:
  git ls-remote --heads origin 'rc/*'

  ├── Found rc/V0.1.0
  │   → AskUserQuestion: Promote / Cancel / Exit
  │
  └── None found
      → proceed to cut new rc

Version derived from remote tags (not local config)
→ no conflicts between team members running /team-rc simultaneously
```

---

## GitHub Ecosystem Integration

### What We Use and Why

```
══════════════════════════════════════════════════════════════════════
                GITHUB ECOSYSTEM — COMPLETE MAP
══════════════════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────────────┐
│                        GITHUB FEATURES                           │
│                                                                   │
│  ┌─────────┐    ┌─────────┐    ┌────────────┐    ┌──────────┐  │
│  │ Issues   │    │ Labels  │    │ Milestones │    │   PRs    │  │
│  │          │    │         │    │            │    │          │  │
│  │ Source   │◄──►│ State   │    │ Version    │    │ Delivery │  │
│  │ of truth │    │ machine │    │ tracking   │    │ mechanism│  │
│  │ for tasks│    │ for     │    │ (V0.1 has  │    │ squash   │  │
│  │          │    │ workflow│    │  12 issues) │    │ merge    │  │
│  └────┬─────┘    └────┬────┘    └─────┬──────┘    └────┬─────┘  │
│       │               │              │                  │        │
│  ┌────▼─────┐    ┌────▼────┐    ┌────▼──────┐    ┌────▼─────┐  │
│  │ Issue    │    │ Branch  │    │ Releases  │    │ Actions  │  │
│  │ Template │    │ Protect │    │ + Tags    │    │          │  │
│  │          │    │         │    │           │    │ ci.yml   │  │
│  │ mission  │    │ develop │    │ V0.1.0    │    │ post-    │  │
│  │ form     │    │ main    │    │ V0.1.1    │    │ merge.yml│  │
│  │ with     │    │ no force│    │ generate  │    │ deploy   │  │
│  │ priority │    │ push    │    │ notes     │    │ .yml     │  │
│  └──────────┘    └─────────┘    └───────────┘    └──────────┘  │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │                    REPO SETTINGS                          │    │
│  │                                                            │    │
│  │  ● delete_branch_on_merge: true                           │    │
│  │  ● allow_squash_merge: true (ONLY squash allowed)         │    │
│  │  ● allow_merge_commit: false                              │    │
│  │  ● allow_rebase_merge: false                              │    │
│  │  ● allow_auto_merge: true                                 │    │
│  │  ● squash_merge_commit_message: PR_BODY                   │    │
│  │    (preserves "Closes #N" in squash commit)               │    │
│  └──────────────────────────────────────────────────────────┘    │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### GitHub Actions — The Three Workflows

```
┌──────────────────────────────────────────────────────────────┐
│ ci.yml                                                        │
│ ────────                                                      │
│ Triggers: push to develop/main, PR to develop/main            │
│ Does: lint → test → build                                     │
│ Purpose: quality gate for ALL code changes                    │
│                                                                │
│   on:                                                          │
│     push: [main, develop]                                     │
│     pull_request: [develop, main]  ← main for rc→main PRs    │
└──────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────┐
│ post-merge.yml                                                │
│ ──────────────                                                │
│ Triggers: pull_request closed (merged)                        │
│ Does:                                                          │
│   1. Find "Closes #N" in PR body                              │
│   2. Remove status:wip / status:review labels                 │
│   3. Add status:done label                                    │
│   4. Check all checkboxes in Issue body                       │
│   5. Check milestone completion → comment if 100%             │
│ Purpose: automate what humans forget                          │
└──────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────┐
│ deploy.yml (optional, team-provided)                          │
│ ──────────                                                    │
│ Triggers: workflow_dispatch ONLY (no auto-deploy)             │
│ Inputs:                                                        │
│   environment: staging | prod                                 │
│   ref: rc/* branch (staging) or V* tag (prod)                 │
│ Purpose: explicit, parameterized deployment                   │
│                                                                │
│ Why NOT auto-deploy on push to main?                          │
│   → Decouples version semantics from deployment               │
│   → Rollback = re-run with previous tag, not git revert      │
│   → No accidental prod deploys from merge button              │
└──────────────────────────────────────────────────────────────┘
```

---

## The Mission Contract

The Contract is the bridge between GitHub (permanent) and AI execution (ephemeral).

```
┌──────────────────────────────────────────────────────────────────┐
│                     MISSION CONTRACT                              │
│                                                                    │
│  Location: .teamwork/active/MISSION-42.md (gitignored)           │
│                                                                    │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │ YAML Frontmatter                                          │    │
│  │                                                            │    │
│  │ issue: 42                                                 │    │
│  │ title: "feat: add user authentication"                    │    │
│  │ branch: mission/42-add-user-auth                          │    │
│  │ assignee: lysfighting                                     │    │
│  │ priority: P1                                              │    │
│  │ milestone: "V0.1"                                         │    │
│  │ issue_content_hash: "sha256..."  ← freshness check       │    │
│  └──────────────────────────────────────────────────────────┘    │
│                                                                    │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │ Body (AI-enriched)                                        │    │
│  │                                                            │    │
│  │ ## Objective (from Issue)                                  │    │
│  │ ## Sub-tasks (checkboxes, ticked during /team-drive)      │    │
│  │ ## Success Criteria (from Issue)                          │    │
│  │ ## Context Files (AI-SCANNED — the key enrichment)        │    │
│  │    → src/auth/login.ts                                    │    │
│  │    → src/middleware/auth.ts                                │    │
│  │    → tests/auth.test.ts                                   │    │
│  │ ## Test Command (from config or Issue)                     │    │
│  │ ## AI Notes (populated during execution)                  │    │
│  └──────────────────────────────────────────────────────────┘    │
│                                                                    │
│  LIFECYCLE:                                                        │
│    Created by /team-claim                                         │
│    Updated by /team-drive (subtask checkoffs, AI notes)           │
│    Read by /team-ship (generates PR body)                         │
│    DELETED by /team-ship (ephemeral — PR body is the record)      │
│                                                                    │
│  Contract ≠ Issue copy.                                           │
│  Contract = Issue + AI-enriched codebase context.                 │
│                                                                    │
│  FRESHNESS CHECK:                                                  │
│    issue_content_hash = SHA256(title + body) at claim time        │
│    /team-drive re-checks: if hash changed → "Issue updated!"     │
│    Prevents working on stale requirements.                        │
│                                                                    │
└──────────────────────────────────────────────────────────────────┘
```

---

## Quality Gates — Defense in Depth

```
══════════════════════════════════════════════════════════════════════
              QUALITY GATES — FOUR LAYERS OF DEFENSE
══════════════════════════════════════════════════════════════════════

  LAYER 1: LOCAL (git hooks)             ← fastest feedback
  ─────────────────────────
  .githooks/pre-commit → lint
  .githooks/pre-push   → test

  Catches: syntax errors, style violations, test failures
  Speed: seconds
  Bypass: only with --no-verify (discouraged)


  LAYER 2: REMOTE (GitHub Actions CI)    ← authoritative
  ──────────────────────────────────
  .github/workflows/ci.yml → lint + test + build
  Triggers on: PR to develop/main, push to develop/main

  Catches: integration failures, build errors, env-specific bugs
  Speed: minutes
  Cannot bypass (branch protection enforces)


  LAYER 3: STRUCTURAL (Branch Protection) ← prevents accidents
  ───────────────────────────────────────
  develop + main:
    ● Require PR (no direct push)
    ● No force push
    ● No branch deletion
    ● CI must pass (if quality.ci: true)
    ● Review required (if quality.review_required: true)

  Catches: accidental pushes, force-push disasters, unreviewed code
  Bypass: admin only (emergency use)


  LAYER 4: AUTOMATION (post-merge Action) ← prevents forgetting
  ──────────────────────────────────────
  .github/workflows/post-merge.yml
  Triggers on: PR merged

  Does: label transition, checkbox checking, milestone tracking
  Catches: human forgetfulness (the #1 cause of stale boards)


  ┌─────────────────────────────────────────────────────────┐
  │                                                           │
  │    Developer writes code                                  │
  │         │                                                 │
  │    git commit → pre-commit (lint) ── FAIL? → fix locally │
  │         │                                                 │
  │    git push  → pre-push (test) ──── FAIL? → fix locally │
  │         │                                                 │
  │    PR created → CI (lint+test+build) FAIL? → fix in PR   │
  │         │                                                 │
  │    PR merged → post-merge Action ── auto-cleanup         │
  │         │                                                 │
  │    ✅ Done. Four layers passed.                           │
  │                                                           │
  └─────────────────────────────────────────────────────────┘

══════════════════════════════════════════════════════════════════════
```

### Config-Driven Protection

The protection rules adapt to your team's maturity:

```
┌────────────┬──────────────────┬────────────────────────────────┐
│ quality.ci │ quality.review   │ What happens                   │
├────────────┼──────────────────┼────────────────────────────────┤
│ false      │ false            │ PR required (blocks direct     │
│            │                  │ push + force push). No other   │
│            │                  │ gates. Good for starting out.  │
├────────────┼──────────────────┼────────────────────────────────┤
│ true       │ false            │ CI must pass. PR required.     │
│            │                  │ Self-merge OK. Most teams.     │
├────────────┼──────────────────┼────────────────────────────────┤
│ false      │ true             │ 1 reviewer required. No CI     │
│            │                  │ gate. Rare configuration.      │
├────────────┼──────────────────┼────────────────────────────────┤
│ true       │ true             │ CI + reviewer. Full ceremony.  │
│            │                  │ Recommended for production.    │
└────────────┴──────────────────┴────────────────────────────────┘
```

---

## Versioning Model

```
══════════════════════════════════════════════════════════════════════
                   VERSIONING — TWO LEVELS
══════════════════════════════════════════════════════════════════════

  MILESTONE = Product Goal (big, strategic)
  ═════════════════════════════════════════

  Milestone V0.1 "AI Camera Pipeline"
    ├── Release V0.1.0 (features A, B, C)     ← deployment unit
    ├── Release V0.1.1 (features D, E)         ← small, frequent
    ├── Release V0.1.2 (hotfix)                ← patch
    └── Release V0.1.3 (feature F)
        └── all issues done → milestone closes automatically

  Milestone V0.2 "V4 Integration"
    ├── Release V0.2.0
    └── ...


  RELEASE = Deployment Unit (small, frequent)
  ═════════════════════════════════════════

  Each release is:
    ● One squash commit on main
    ● One git tag (V0.1.0)
    ● One GitHub Release (with auto-generated notes)

  Patch version auto-derived:
    git tag -l "V0.1.*" → V0.1.0, V0.1.1 → next is V0.1.2

  Config stores only the milestone prefix:
    versions:
      current: "V0.1"    # NOT "V0.1.3" — patch is automatic


  WHY THIS SEPARATION?
  ═════════════════════════════════════════

  Milestone ≠ Release. They serve different purposes:

  ┌────────────┬───────────────────┬───────────────────────┐
  │            │ Milestone         │ Release               │
  ├────────────┼───────────────────┼───────────────────────┤
  │ Granularity│ Weeks to months   │ Days to week          │
  │ Purpose    │ Product direction │ Deployment unit       │
  │ Tracked by │ GitHub Milestones │ Git tags + Releases   │
  │ Closes when│ All issues done   │ RC promoted to main   │
  │ Example    │ "V0.1 — Camera"   │ V0.1.0, V0.1.1       │
  └────────────┴───────────────────┴───────────────────────┘

══════════════════════════════════════════════════════════════════════
```

---

## Scripts Architecture

### The Nine Scripts

```
┌──────────────────────────────────────────────────────────────────┐
│                      SCRIPTS LAYER (v3.9.6)                         │
│                                                                    │
│  tw-config.sh ─── Foundation ─── Config reader + team queries     │
│  │  Key lookup:   tw-config.sh <dotted.key> [default]            │
│  │  Subcommands:  detect-dir, resolve-labels, get-user-level,    │
│  │                team-mode, resolve-member, list-members,        │
│  │                migrate, has-section                             │
│  │  Python3 YAML parser (handles nested keys, quotes, comments)  │
│  │  Shared _PYTHON_MEMBER_PARSER for member/role queries         │
│  │                                                                │
│  ├── tw-git.sh ─── 21 subcommands ─── All git operations         │
│  │   Branch:    ensure-base, create-branch, build-branch-name,   │
│  │              extract-issue-number, current, protect-check      │
│  │   Core:      push, rebase, tag, commit, slugify               │
│  │   Worktree:  worktree-add, worktree-remove,                   │
│  │              worktree-main-repo, worktree-find-by-branch      │
│  │   Query:     log-since, list-merged-branches                  │
│  │   Release:   cut-release, find-rc, next-version,              │
│  │              milestone-resolve                                 │
│  │                                                                │
│  ├── tw-label.sh ─── 4 subcommands ─── Label state machine       │
│  │   transition (from → to), set (force), verify, pr-label       │
│  │                                                                │
│  ├── tw-contract.sh ─── 6 subcommands ─── Contract lifecycle     │
│  │   read-field, hash, check-freshness, toggle-task,             │
│  │   sync-all-checkboxes, delete                                  │
│  │                                                                │
│  ├── tw-pr.sh ─── 8 subcommands ─── PR lifecycle                 │
│  │   exists, commit-type, create, comment, watch (CI polling),   │
│  │   wait-merged (merge state polling), add-reviewer,             │
│  │   verify-merged                                                │
│  │                                                                │
│  ├── tw-notify.sh ─── Notification layer                         │
│  │   Events: mc.created, mc.completed, milestone.created,         │
│  │           milestone.done, review.requested,                     │
│  │           review.changes_requested, review.approved             │
│  │   Channels: Slack (block kit), Feishu (card), custom webhook  │
│  │   Always non-fatal (exit 0 on failure)                         │
│  │                                                                │
│  ├── setup-github-labels.sh ─── Create all labels on repo        │
│  │                                                                │
│  ├── tw-e2e-test.sh ─── Full end-to-end validation               │
│  │   8 phases: pre-check → config → script-unit → create →       │
│  │   claim → drive → ship → cleanup. Safe to run on any repo.    │
│  │                                                                │
│  └── tw-bump-version.sh ─── Version management                   │
│      Updates version strings across all skill files.              │
│      Handles frontmatter, inline, bash variable patterns.         │
│                                                                    │
└──────────────────────────────────────────────────────────────────┘
```

### Why Scripts?

Evidence from v2.6.1 falsification audit:

```
    228 total operations catalogued:
    ├── 145 deterministic (64%) → should be scripts
    ├──  48 LLM-judgment (21%) → should stay in skill.md
    └──  35 mixed (15%)        → split into script + skill

    Bug distribution:
    ├── 21 bugs found in LLM-interpreted operations
    └──  0 bugs found in existing scripts

    Conclusion: LLM re-interprets deterministic operations
    differently across calls. Scripts eliminate that class
    of bugs entirely.
```

---

## Config Reference

```yaml
# .teamwork/config.yml (or .teamspace/config.yml)
# ─────────────────────────────────────────────────
# Schema version 3

schema_version: 3

# Team metadata
team:
  name: "My Team"
  repo: org/repo-name

# Role definitions (used for dashboard view distinction, NOT permissions)
roles:
  - id: leader
    level: leader
    label: "Tech Lead"
    description: "Product direction, task creation, review, release"
  - id: engineer
    level: member
    label: "Engineer"
    description: "Feature development"

# Team roster
members:
  - github: lysfighting
    name: "Liya Song"
    role: leader
  - github: yuang-yang
    name: "yuang-yang"
    role: engineer

# Project detection (auto-generated by /team init)
# Single-repo:
project:
  language: python
  test_command: "pytest tests/ -v"
  lint_command: "ruff check ."
  build_command: "python -m build"

# Monorepo (service-aware testing):
# project:
#   type: monorepo
#   services:
#     - name: api-server
#       language: python
#       test_command: "cd api-server && pytest tests/ -v"
#       lint_command: "cd api-server && ruff check ."
#     - name: frontend
#       language: typescript
#       test_command: "cd frontend && npm test"
#       build_command: "cd frontend && npm run build"

# Git conventions
conventions:
  branch_pattern: "mission/{issue}-{slug}"
  base_branch: develop
  production_branch: main
  commit_format: "type(scope): description | Mission: #{issue}"

# Labels (two schemas supported)
# Schema v3:
labels:
  mission: "mission"
  status_prefix: "status:"
  priority_prefix: "priority:"
# Legacy (also supported):
# label_prefix:
#   status: "status:"
#   priority: "priority:"

# Quality gates
quality:
  hooks: true
  ci: true
  review_required: false
  branch_protection: true

# Notifications (v3 NEW)
notifications:
  enabled: true
  channels:
    - type: slack
      channel: "#team"
      webhook: ""
      enabled: false
    - type: feishu
      webhook: ""
      enabled: false
  events:
    mc.created: [slack, feishu]
    mc.completed: [slack, feishu]
    milestone.created: [slack, feishu]
    milestone.done: [slack, feishu]
    review.requested: [slack, feishu]

# Deployment
deploy:
  staging_workflow: "deploy.yml"

# Version tracking
versions:
  current: "V0.1"

# Worktree mode is LOCAL (per-user, not in config.yml):
#   git config --local teamwork.worktree true

# Terminal tab title is LOCAL (per-user, configured in ~/.zshrc or equivalent):
#   /team init asks during setup — auto-detects shell (zsh/bash/fish)
#   Shows current branch name in terminal tab for multi-window awareness

# Mission label (legacy, prefer labels.mission)
# mc_label: "mission"
```

### Smart Sync (v3.7.0+)

**Dashboard behind-detection**: When a PR's branch falls behind the base branch, the dashboard shows it inline:

```
🔀 My PRs
  #162 feat(infra): deploy Firebase Auth  ✅ CI passing  ⚠️ 8 commits behind main
```

The NEXT box also prioritizes the sync recommendation:

```
┌─ 💡 NEXT ────────────────────────────────────┐
│  PR #162 is 8 commits behind main             │
│  — /team-ship sync to rebase                  │
└───────────────────────────────────────────────┘
```

**AI conflict resolution**: `/team-ship sync` no longer stops when rebase hits conflicts. Instead, Claude reads both sides of each conflict, resolves them intelligently, and continues the rebase. Ambiguous conflicts (semantic contradictions where both sides changed the same logic differently) are resolved with best judgment and flagged in the summary for review. After rebase completes, the branch is force-pushed automatically.

### Tips System (v3.6.1+)

Every `/team-*` command displays a random tip after its output. Tips are sourced from `scripts/tw-tips.txt` (installed alongside the skills). This enables progressive feature discovery — users learn about lesser-known commands (`/team auto`, `/team wrap`, `/team doctor`, etc.) through natural usage.

### Config Schema Evolution

| Version | What changed |
|---------|-------------|
| v1 | `team:` as member array, flat `mc_label`, `label_prefix:` |
| v2 | Added `deploy:`, `versions:`, `worktree:` |
| v3 | Added `roles:`, `members:` (separate from `team:`), `notifications:`, `project.services[]`. `team:` = metadata. Push model default. Worktree moved to `git config --local` (per-user). |

MERGE MODE (`/team init` on existing config): preserves all existing content, updates `schema_version`, removes `skill_version` from config if present, appends only missing sections. Safe to run on any existing config.

---

## End-to-End Example

```
══════════════════════════════════════════════════════════════════════
              A WEEK WITH TEAMWORK v3 — END TO END
══════════════════════════════════════════════════════════════════════

  MONDAY MORNING
  ──────────────

  $ /team-issue "实现用户登录功能，支持 OAuth + 密码登录" @lysfighting
    → Issue #42 created
    → Assigned to @lysfighting
    → Labels: mission, status:wip, priority:P1, size:M
    → Milestone: V0.1
    → Branch: mission/42-user-login
    → Notification sent to: #team (Slack)

  $ /team-claim 42
    → Contract generated: .teamwork/active/MISSION-42.md
      (AI scans codebase → discovers auth files, test files)
    → Switched to branch mission/42-user-login

  $ /team-drive
    → Reading Contract... 5 sub-tasks found
    → Freshness check: Issue #42 unchanged ✓
    → [1/5] Set up auth middleware ✓ (committed)
    → [2/5] Implement password login ✓ (committed)
    → [3/5] Implement OAuth flow ✓ (committed)
    → [4/5] Add auth tests ✓ (committed)
    → [5/5] Update API docs ✓ (committed)
    → Service-aware test: api-server tests passing ✓
    → All sub-tasks complete.

  $ /team-ship
    → Pushed branch to origin
    → PR #15 created: "feat: user login system | Closes #42"
    → CI running... ✓ passed
    → Label: status:wip → status:review
    → Contract deleted (ephemeral)
    → Notification: mc.completed sent

  [REVIEWER MERGES PR #15 — squash merge]
    → post-merge Action fires:
       - Label: status:review → status:done
       - All checkboxes in Issue #42 checked
       - Issue #42 auto-closed (via "Closes #42")

  $ /team-ship done
    → PR #15 verified merged ✓
    → Returned to develop branch
    → Clean workspace


  MONDAY AFTERNOON — BATCH CREATION
  ──────────────────────────────────

  $ /team-issue batch V0.1 -- User system and campaign basics
      MC1: 权限管理系统 @yuang-yang
        - [ ] RBAC model implemented
        - [ ] Permission middleware
      MC2: Campaign UI @xxLe
        - [ ] CRUD for campaigns
      MC3: API rate limiting @lysfighting
        - [ ] Rate limiter middleware

    → Milestone V0.1 created (or reused)
    → Issue #45 → @yuang-yang → branch mission/45-rbac
    → Issue #46 → @xxLe → branch mission/46-campaign-ui
    → Issue #47 → @lysfighting → branch mission/47-rate-limiting
    → Notifications: milestone.created + 3x mc.created


  TUESDAY
  ───────

  $ /team-rc
    → No existing RC branch ✓
    → Next version: V0.1.0 (auto-derived, no tags yet)
    → Develop CI: green ✓
    → Branch rc/V0.1.0 cut from develop
    → Staging deploy triggered

  [VERIFY ON STAGING]
    → Found bug in OAuth redirect!

  [HOTFIX]
    → Branch from rc/V0.1.0 → fix → PR to rc/V0.1.0 (squash)
    → Cherry-pick fix to develop (immediately!)
    → Re-deploy staging → re-verify → all good ✓

  $ /team-rc promote
    → PR: rc/V0.1.0 → main (squash merge)
    → Tag V0.1.0 on main
    → GitHub Release created (auto-generated notes)
    → Milestone V0.1: 3/6 issues closed
    → Branch rc/V0.1.0 deleted
    → Notification: milestone.done sent


  WEDNESDAY - THURSDAY
  ────────────────────

  [More features: Issues #48, #51, #53]
  [Each: claim → drive → ship → PR merged]

  [PM creates Issue #54 directly on GitHub, assigns @lysfighting]

  $ /team                      ← role-based dashboard
    → Leader view: all team members + their active missions
    → Version V0.1: 5/6 tasks done (83%)
    → UNTRACKED (1):
        #54 https cert missing on Safari [@lysfighting]
        Fix: /team-issue fix #54

  $ /team-issue fix #54        ← fix untracked Issue
    → Issue #54 analyzed: priority P1, size M
    → Labels added: mission, status:wip, priority:P1, size:M
    → Branch: mission/54-https-cert-missing
    → #54 now visible on dashboard

  $ /team #48                  ← MC detail view
    → Issue #48 — fix: API timeout
    → Branch: mission/48-fix-api-timeout
    → Commits: 3 since develop
    → PRs: #22 (open, CI passing)


  FRIDAY
  ──────

  $ /team-rc
    → Next version: V0.1.1 (auto-derived from V0.1.0)
    → rc/V0.1.1 cut from develop
    → Staging verify → pass ✓

  $ /team-rc promote
    → Tag V0.1.1
    → Milestone V0.1: 6/6 issues closed
    → Milestone V0.1 auto-closed! 🎉
    → Notification: milestone.done → all channels


  NEXT WEEK: V0.2 begins...

══════════════════════════════════════════════════════════════════════
```

---

## Branch Protection

```
┌─────────┬────────────┬──────┬────────┬────────────┬──────────┐
│ Branch  │ Require PR │ CI   │ Review │ Force Push │ Delete   │
├─────────┼────────────┼──────┼────────┼────────────┼──────────┤
│ main    │ yes        │ yes  │ config │ blocked    │ blocked  │
│ develop │ yes        │ yes  │ config │ blocked    │ blocked  │
├─────────┼────────────┼──────┼────────┼────────────┼──────────┤
│ rc/*    │ no (but    │ —    │ —      │ allowed    │ allowed  │
│         │ protect-   │      │        │ (for       │ (after   │
│         │ check      │      │        │ hotfix)    │ promote) │
│         │ blocks     │      │        │            │          │
│         │ direct     │      │        │            │          │
│         │ commits)   │      │        │            │          │
├─────────┼────────────┼──────┼────────┼────────────┼──────────┤
│ feat/*  │ no         │ —    │ —      │ allowed    │ auto     │
│ mission/│            │      │        │            │ (merge)  │
└─────────┴────────────┴──────┴────────┴────────────┴──────────┘

Admin bypass on main + develop (for emergency use by project lead).
```

---

## Closing Thoughts

> *"The best process is invisible. It doesn't ask you to change how you work — it changes how work flows around you."*

Teamwork doesn't replace GitHub. It doesn't ask you to learn a new tool. It takes the GitHub workflow you already know — Issues, PRs, branches, tags — and closes the gaps that AI-era commit frequency creates.

The six skills are not six new things to learn. They are six moments in a developer's day that used to require context-switching, remembering commands, and manual bookkeeping. Now they don't.

```
Before:                              After:
  "What branch am I on?"              /team-claim #42
  "What's the Issue number?"          /team-drive
  "Did I push?"                       /team-ship
  "Is CI passing?"                    /team-rc
  "Did I close the Issue?"
  "Did I update the labels?"
  "What version is this?"
  "Who's working on what?"            /team
  "Who's assigned to what?"           /team-issue batch ...
  "Someone assigned me an issue       /team-issue fix #94
   but it doesn't show up"
```

The loop never breaks. The system remembers what you forget.

---

*Teamwork v3.9.6 — Built for teams where humans and AI write code together.*
