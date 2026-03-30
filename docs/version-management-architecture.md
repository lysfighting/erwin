# Version Management Architecture

> The knowledge infrastructure for AI-assisted iterative development on erwin.

---

## Executive Summary

erwin is developed through multiple AI-assisted conversations over weeks/months. Without
persistent memory infrastructure, each conversation starts from zero, loses context, and
cannot build on previous decisions. This document describes the knowledge management system
that makes iterative AI development reliable and auditable.

**Core principle**: Version control is our memory. Everything significant gets committed.

---

## The Problem

```
Traditional development:
  Developer → Code → Git → History

AI-assisted development WITHOUT this system:
  Session 1 → work → forgotten
  Session 2 → work → forgotten
  Session N → "I forgot what we decided..."
```

Claude Code's built-in memory (`~/.claude/projects/.../memory/`) is:
- Machine-local (not portable)
- Not version-controlled (not auditable)
- Limited to 200 lines in MEMORY.md
- Opaque (hard to manually review and correct)

We need a **human-readable, git-tracked, deterministic bootstrap** for every session.

---

## Solution Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              ERWIN KNOWLEDGE MANAGEMENT SYSTEM               │
│                                                             │
│  Git Repository                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │                                                       │ │
│  │  CLAUDE.md              ← Session bootstrap (auto-    │ │
│  │  (project root)            loaded, <200 lines)        │ │
│  │                                                       │ │
│  │  CHANGELOG.md           ← What changed, when, why    │ │
│  │  (project root)                                       │ │
│  │                                                       │ │
│  │  docs/                                                │ │
│  │  ├── CURRENT_STATE.md   ← Living status document      │ │
│  │  ├── decisions/         ← Architecture Decision       │ │
│  │  │   └── NNNN-*.md         Records (ADRs)            │ │
│  │  ├── sessions/          ← Per-conversation artifacts  │ │
│  │  │   └── YYYY-MM-DD.md                               │ │
│  │  └── *.md               ← Architecture, analysis     │ │
│  │                                                       │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## File Responsibilities

### 1. `CLAUDE.md` — The Session Bootstrap

**Purpose**: Loaded automatically at session start by Claude Code. Tells Claude everything
it needs to know about the project in <200 lines.

**Contents**:
- What is this project (1 paragraph)
- Repository structure (tree diagram)
- Tech stack
- Current architecture plan (summary + link to full doc)
- Workflow rules
- Key files table
- Anti-patterns (what NOT to do)

**Rule**: Keep under 200 lines. Link to detailed docs, don't inline them.

### 2. `CHANGELOG.md` — Version History

**Purpose**: Human-readable record of what changed, when, and why.

**Format**: Keep a Changelog + Semantic Versioning + Rationale sections

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added / Changed / Fixed / Removed / Security

- **Feature name** — description

### Rationale
Why this version exists, what problem it solves.
```

**Versioning policy**:
- `MAJOR.x.x` — Complete redesign or breaking user-facing change
- `x.MINOR.x` — New features (new sections, new components, new capabilities)
- `x.x.PATCH` — Bug fixes, copy edits, minor UI tweaks

**Current policy**:
- `/erwin/CHANGELOG.md` — the durable project-level changelog

Historical note: `erwin` used to keep a separate showcase-level changelog, but that app has been retired and removed from the active repository surface.

### 3. `_project/CURRENT.md` — Living Status

**Purpose**: The "what's happening right now" document. Read at session start, updated at session end.

**Contents**:
- Current version
- Last session link
- What was just completed
- What's in progress
- Next session's first task
- Implementation roadmap checkboxes
- Known issues / blockers

**Critical rule**: This file is useless if it's stale. Update it at the end of EVERY session.

### 4. `docs/decisions/` — Architecture Decision Records

**Purpose**: Record WHY decisions were made, including what was rejected.

**When to write an ADR**:
- Any decision that would be non-obvious to someone starting fresh
- Any decision where multiple reasonable alternatives existed
- Any decision that constrains future choices

**Format**: `NNNN-short-title.md`

```markdown
# ADR-NNNN: Title
- Date / Status / Session

## Context
## Decision
## Consequences
## Alternatives Considered
```

### 5. `docs/sessions/` — Conversation Artifacts

**Purpose**: Archive of each session's work, decisions, and outputs.

**When to write**: At the end of each session (or the start of the next).

**Format**: `YYYY-MM-DD.md` (or `YYYY-MM-DD-2.md` for multiple sessions per day)

```markdown
# Session: YYYY-MM-DD

## Context (What I knew coming in)
## Work Done
## Key Decisions Made
## Artifacts Created/Modified
## What's Next
## CURRENT_STATE.md Updates
```

---

## Session Protocol

### Starting a Session

```
1. AGENTS.md / CLAUDE.md provide the project bootstrap
2. Read `_project/CURRENT.md` (current status, in-progress work)
3. If working on a specific feature, read the relevant ADRs
4. Begin work
```

### During a Session

- Make decisions → note them (ADR if significant, inline in session log if minor)
- Create/modify files → note them in session log
- Update CHANGELOG.md when completing features

### Ending a Session

```
1. Update `_project/CURRENT.md`:
   - Move completed tasks to "What Was Just Completed"
   - Update "Active Task / What's In Progress"
   - Update "Next Session Should Start With"

2. Write/update `_project/sessions/YYYY-MM-DD*.md`

3. Update CHANGELOG.md if features were completed

4. git add + git commit
   Format: "type(scope): description"
   Example: "feat(harness): add scheduler admission policy"
```

---

## Commit Convention

```
feat(scope): add new feature
fix(scope): fix a bug
docs(scope): documentation only
refactor(scope): code refactor, no behavior change
chore(scope): maintenance (deps, gitignore, etc.)
style(scope): formatting only
```

**Scopes**: `harness`, `radar`, `docs`, `skills`, `archive`, `root`

**Examples**:
```
feat(harness): add scheduler admission policy
docs(decisions): add ADR-0006 for simulation engine design
chore(root): initialize git repository and version management
docs(radar): tighten source-priority policy after a live run
```

---

## What This System Replaces

| Before | After |
|--------|-------|
| Relying on Claude's session memory | `_project/CURRENT.md` + git history |
| "We discussed this last time" | ADR with full context |
| Forgetting why something was built | Rationale in CHANGELOG + ADRs |
| Starting each session from scratch | AGENTS/CLAUDE bootstrap + `_project/CURRENT.md` |
| No record of rejected alternatives | ADR "Alternatives Considered" sections |

---

## Stress Test Q&A

**Q: What if I start a session and CURRENT_STATE.md is stale?**
A: The git log is the fallback. `git log --oneline -20` shows recent commits. `git diff HEAD~5` shows recent changes.

**Q: What if there are multiple sessions in one day?**
A: Name them `2026-03-02-1.md`, `2026-03-02-2.md`. CURRENT_STATE.md points to the latest.

**Q: What if a decision was wrong and needs to be reversed?**
A: Mark the ADR as "Deprecated" or "Superseded by ADR-NNNN". Create a new ADR for the replacement. Never delete old ADRs.

**Q: How much detail in session logs?**
A: Enough that you could reconstruct the session's work in 5 minutes of reading. Not a transcript — a summary with key decisions and artifacts.

**Q: What about the Claude Code auto-memory at ~/.claude/projects/?**
A: Use it as a supplement, not a replacement. The auto-memory is good for short-term operational notes. The git-tracked system here is the authoritative record.

---

## File Size Targets

| File | Target | Max |
|------|--------|-----|
| `CLAUDE.md` | <150 lines | 200 lines |
| `docs/CURRENT_STATE.md` | <80 lines | 120 lines |
| ADR files | 50-100 lines | 150 lines |
| Session logs | 50-150 lines | 300 lines |
| `CHANGELOG.md` | grows with project | unbounded |

---

*End of document*
