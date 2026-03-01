# Project Management Architecture

> The complete knowledge and management infrastructure for AI-assisted iterative development.
> Version 1.0 | Date: 2026-03-02

---

## Executive Summary

This document specifies the information architecture for erwin: a project managed through
multiple AI (Claude) conversations over weeks or months. The system solves three problems:

1. **Context loss**: Each AI session starts fresh → solved by CLAUDE.md + _project/CURRENT.md
2. **Decision loss**: "Why did we do it this way?" → solved by ADRs
3. **Vision drift**: Daily tasks obscure the north star → solved by GOALS.md + Idea Garden

The design follows a three-layer separation: **Product** (what + why) / **Process** (how we work) / **Artifact** (what we built).

---

## The Three-Layer Model

```
┌─────────────────────────────────────────────────────────────┐
│                    erwin/ (git root)                         │
│                                                             │
│  CLAUDE.md ──────── Bridge ──────── CHANGELOG.md           │
│  (AI bootstrap)      ↕↕              (version history)      │
│                      ↕↕                                     │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  _project/  ← PROJECT INTELLIGENCE LAYER            │   │
│  │                                                     │   │
│  │  GOALS.md        CURRENT.md       ROADMAP.md        │   │
│  │  (vision+OKR)   (living status)   (milestones)      │   │
│  │                                                     │   │
│  │  ideas/          decisions/        sessions/        │   │
│  │  (idea garden)   (ADRs)           (conv. logs)      │   │
│  │                                                     │   │
│  │  issues/         releases/                          │   │
│  │  (bug/feature)   (release notes)                    │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  docs/  ← TECHNICAL DOCUMENTATION LAYER             │   │
│  │  Architecture docs, analysis, research               │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌──────────────────┐  ┌──────────────────────────────┐   │
│  │  drive-showcase/ │  │  reference/                  │   │
│  │  CODE ARTIFACT   │  │  SOURCE MATERIAL (read-only) │   │
│  └──────────────────┘  └──────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## Full Directory Structure

```
erwin/
│
├── CLAUDE.md                  ← Session bootstrap (auto-loaded, <200 lines)
├── CHANGELOG.md               ← Project version history
│
├── .claude/                   ← Claude Code configuration
│   └── rules/                 ← Modular instruction files (topic-scoped)
│       ├── coding-standards.md
│       ├── git-conventions.md
│       └── ai-workflow.md
│
├── _project/                  ← PROJECT INTELLIGENCE LAYER
│   │
│   ├── GOALS.md               ← North Star + current OKRs (the WHY)
│   ├── CURRENT.md             ← Living status document (the NOW)
│   ├── ROADMAP.md             ← Milestone plan (the WHEN)
│   │
│   ├── ideas/                 ← Idea Garden
│   │   ├── README.md          ← Garden rules + stage definitions
│   │   ├── INBOX.md           ← Raw captures (process weekly)
│   │   ├── skills/            ← Drive-skill-like skill ideas
│   │   │   └── *.md           ← One file per idea
│   │   └── meta/              ← Universal/cross-project ideas
│   │       └── project-management-template.md
│   │
│   ├── decisions/             ← Architecture Decision Records (ADRs)
│   │   ├── README.md          ← ADR index + template
│   │   └── NNNN-title.md
│   │
│   ├── sessions/              ← Conversation artifacts (AI session logs)
│   │   ├── README.md
│   │   └── YYYY-MM-DD.md
│   │
│   ├── issues/                ← Lightweight issue tracking
│   │   ├── README.md
│   │   └── ISS-NNNN-title.md
│   │
│   └── releases/              ← Detailed release notes
│       └── vX.Y.Z.md
│
├── docs/                      ← TECHNICAL DOCUMENTATION LAYER
│   └── *.md                   ← Architecture, analysis, research docs
│
├── drive-showcase/            ← CODE ARTIFACT (React web app)
│   ├── CHANGELOG.md           ← App-level changelog
│   └── src/
│
└── reference/                 ← SOURCE MATERIAL (read-only)
    └── drive-skill/           ← The reference implementation
```

---

## Component Specifications

### `_project/GOALS.md` — North Star

**Purpose**: The single source of truth for WHERE this project is going. Read at session start to
prevent vision drift during daily work.

**Structure**:
```markdown
# Project Goals

## North Star (永不变)
[One paragraph: what success looks like in 2 years]

## Current OKRs (季度滚动)

### Objective 1: [...]
- KR1: [measurable] → Progress: X%
- KR2: [measurable] → Progress: X%

### Objective 2: [...]
...

## What This Project IS
[3 bullet points]

## What This Project IS NOT
[3 bullet points — explicit anti-goals]
```

**Update frequency**: OKRs updated quarterly. North Star updated only on major pivots.

---

### `_project/ideas/` — The Idea Garden

**Metaphor**: A garden, not a todo list. Ideas have lifecycle stages and need tending.

**Four stages**:

```
Seed 🌱    →    Sprout 🌿    →    Plant 🌱    →    Fruit 🍎
(raw capture)   (explored)      (specified)      (built/shipped)
INBOX.md        skills/*.md     with ADR link    with issue link
```

**Stage definitions**:

| Stage | Criteria | Location | Template |
|-------|----------|----------|---------|
| Seed 🌱 | 1 sentence, unvalidated | INBOX.md | `- 💡 idea title` |
| Sprout 🌿 | Has problem + hypothesis | ideas/skills/*.md | Full idea file |
| Plant 🌱 | Has ADR, ready to build | ideas/skills/*.md | + ADR link |
| Fruit 🍎 | Built and shipped | Archived | + issue/release link |

**Weekly process**: Review INBOX.md, promote Seed→Sprout or discard.

---

### `_project/issues/` — Lightweight Issue Tracking

**Purpose**: Track bugs, features, and tasks WITHOUT GitHub dependency. Pure markdown.

**Naming**: `ISS-NNNN-short-title.md`

**Template**:
```markdown
# ISS-NNNN: Title

- **Date**: YYYY-MM-DD
- **Status**: Open | In Progress | Done | Won't Fix
- **Type**: Bug | Feature | Task | Research
- **Priority**: P0 (urgent) | P1 (soon) | P2 (backlog)
- **Session**: [link to session that created this]
- **Commit**: [git SHA or version when closed]

## Description
[What's wrong or what's needed]

## Acceptance Criteria
- [ ] criterion 1
- [ ] criterion 2

## Notes
[Investigation notes, links, context]
```

---

### `_project/releases/` — Detailed Release Notes

**Purpose**: More narrative than CHANGELOG.md. Tells the story of each release.

**When to create**: After each MINOR or MAJOR version bump.

**Template**: `vX.Y.Z.md`
```markdown
# Release v1.0.0

- **Date**: YYYY-MM-DD
- **Session**: [link]
- **Issues Closed**: ISS-0001, ISS-0002, ISS-0005

## What Shipped
[Narrative description, not just bullet list]

## Key Decisions Made
[ADR references]

## What We Learned
[For future reference]

## What's Next
[Preview of v1.1.0]
```

---

## Cross-Artifact Connection System

**The "打通" mechanism**: Every artifact type has a unique ID and can reference others.

```
┌──────────┐   triggers   ┌────────────┐   closes   ┌──────────────┐
│ Session  │ ──────────→  │   Issue    │ ─────────→ │   Release    │
│ YYYY-MM  │              │  ISS-NNNN  │            │   vX.Y.Z     │
└──────────┘              └────────────┘            └──────────────┘
     │                          │                          │
     │ inspires                 │ references               │ links to
     ▼                          ▼                          ▼
┌──────────┐            ┌─────────────┐           ┌──────────────┐
│   Idea   │ → ADR →    │  Decision   │           │  git commit  │
│ ISS-NNNN │            │  ADR-NNNN   │           │  sha/tag     │
└──────────┘            └─────────────┘           └──────────────┘
```

**Linking convention**:
- Session references: `Session: YYYY-MM-DD` (in YAML header or inline)
- Issue references: `ISS-NNNN` (inline link)
- ADR references: `ADR-NNNN` (inline link)
- Release references: `vX.Y.Z` (inline link)
- Commit messages: `Closes ISS-NNNN` (git convention)

---

## The Universal Capability: Project Template

**Key insight from the user**: The version management system should be reusable across projects.

This entire `_project/` structure is a **portable template**. It should live at:
`_project/ideas/meta/project-management-template.md` (as a design idea)
and eventually become:
`reference/project-template/` (an installable scaffold)

**What makes it universal**:
- No project-specific content in the structure
- All project-specific content lives inside the files, not the structure
- Can be applied by running `cp -r reference/project-template/ new-project/_project/`

**Future skill idea** 🌱: `/new-project` — a Claude Code skill that scaffolds a new project with
this entire `_project/` structure pre-populated from a template.

---

## `.claude/rules/` Structure

Following Claude Code's best practice for modular instructions:

```
.claude/
└── rules/
    ├── coding-standards.md     ← Code style, patterns
    ├── git-conventions.md      ← Commit format, branch naming
    └── ai-workflow.md          ← Session protocol, update rules
```

Each file can be scoped to specific file patterns via YAML frontmatter:
```yaml
---
paths:
  - "drive-showcase/src/**/*.tsx"
---
```

CLAUDE.md imports these via `@` syntax if needed, or they're loaded automatically.

---

## Goals Management: OKR Structure

**Following Claude Code community best practices**: Store OKRs in markdown, track alongside code.

```markdown
# GOALS.md — erwin

## North Star
Create the world's best educational platform for understanding and creating
AI agent orchestration skills — so anyone can build their own "drive-skill."

## 2026 Q1 OKRs

### O1: Showcase that converts visitors to skill creators
- KR1: V2 showcase with 4-act structure shipped → 0%
- KR2: SkillBuilder generates downloadable skill.md → 0%
- KR3: OriginTrace maps 30+ philosophical sources → 0%

### O2: Erwin as a meta-skill factory
- KR1: One new skill idea at Plant 🌱 stage → 0%
- KR2: Project template extractable for other projects → 0%

## Anti-Goals (What We Are NOT Building)
- A SaaS product (no auth, no database, no users)
- A perfect product (MVP mentality, ship and iterate)
- A comprehensive AI framework (focused on skill design patterns only)
```

---

## Session Protocol (Updated)

### Starting a Session (3 files, ~5 min)

```
1. CLAUDE.md auto-loads
2. cat _project/CURRENT.md     ← What's happening now
3. cat _project/GOALS.md       ← Remember the north star
```

### Ending a Session (4 updates, ~10 min)

```
1. Update _project/CURRENT.md  ← What changed
2. Write _project/sessions/YYYY-MM-DD.md
3. Create/update _project/issues/ for discovered bugs/tasks
4. git add + git commit (conventional format)
```

---

## Migration Plan (from current state)

**Step 1**: Create `_project/` and move files
```bash
mkdir -p _project/{ideas/skills,ideas/meta,decisions,sessions,issues,releases}
git mv docs/decisions/ _project/decisions/
git mv docs/sessions/ _project/sessions/
git mv docs/CURRENT_STATE.md _project/CURRENT.md
```

**Step 2**: Create new files
```
_project/GOALS.md
_project/ROADMAP.md
_project/ideas/README.md
_project/ideas/INBOX.md
_project/issues/README.md
_project/releases/README.md
```

**Step 3**: Rename reference-skills/ → reference/
```bash
git mv reference-skills/ reference/
```

**Step 4**: Create .claude/rules/

**Step 5**: Update CLAUDE.md paths

**Step 6**: Write architecture doc (this file)

**Step 7**: Commit

---

## Stress Test Q&A

**Q: What if _project/ gets too large over time?**
A: Archive old sessions/issues. Sessions older than 3 months → `sessions/archive/YYYY/`. Active = last 30 days.

**Q: How do issues and ADRs connect to actual git commits?**
A: By convention in commit messages: `feat(showcase): add simulator. Closes ISS-0001. See ADR-0006.`

**Q: What's the difference between _project/CURRENT.md and docs/CURRENT_STATE.md?**
A: Same content, moved to _project/ for organizational clarity. CLAUDE.md updated to point to new path.

**Q: Should drive-showcase/ move to apps/?**
A: Not now. Moving would require updating Vite paths, CI configs, README links. Cost > benefit at current scale. Add `apps/` when a second code artifact exists.

**Q: What about the universal template — where does it live?**
A: First as an idea in `_project/ideas/meta/project-management-template.md`. When ready to extract, moves to `reference/project-template/`. Eventually becomes its own GitHub repo + install.sh.

---

*Document ends here.*
