# /erwin-init

Initialize the erwin cognitive framework for any project — new or existing.

## What it does

Installs the **erwin cognitive framework**: version management, session protocols, ADR system, idea garden, and Claude auto-memory. Not a scaffolding tool — this is a cognitive operating system for AI-assisted development.

## Usage

```
/erwin-init
```

Run from the project's root directory. Zero questions — infrastructure is created instantly.

## How it works

**Phase 1 (automatic, ~15 seconds):**
- Detects project name from directory name
- Creates `_project/`, CLAUDE.md, `.claude/rules/`, CHANGELOG.md
- Initializes Claude auto-memory (local)
- Git init + first commit

**Phase 2 (optional, prompted):**
- Asks: "Want to fill in GOALS.md now?"
- If yes → guided conversation: vision, success criteria, anti-goals
- If no → done, user fills GOALS.md later

Tech stack, project description, and language preferences are NOT asked — Claude discovers them naturally through working on the project.

## What gets created

```
_project/
├── GOALS.md           ← Vision + OKRs (north star)
├── CURRENT.md         ← Living status (read before any implementation)
├── ROADMAP.md         ← Milestone plan template
├── ideas/
│   ├── README.md      ← Idea garden rules (4 stages: Seed→Sprout→Plant→Fruit)
│   └── INBOX.md       ← Raw idea captures, processed weekly
├── decisions/         ← Architecture Decision Records
├── sessions/          ← AI conversation logs
├── issues/            ← Bug/feature tracking
└── releases/          ← Narrative release notes
CLAUDE.md              ← AI navigation handbook
CHANGELOG.md           ← Keep a Changelog format
.claude/
└── rules/
    ├── ai-workflow.md     ← Session protocol
    └── git-conventions.md ← Commit message format
~/.claude/projects/{hash}/memory/MEMORY.md   (local only)
```

## Requirements

- Claude Code
- Git (for the initial commit step)

## Install

Via erwin's install.sh:
```bash
./install.sh erwin-init
```

Or to install all erwin skills:
```bash
./install.sh
```

## Origin

Derived from erwin's `_project/` structure. Redesigned in ADR-0008 from scaffolding tool to cognitive infrastructure installer.

See: [ADR-0008](../../_project/decisions/0008-erwin-init-redesign.md) | [ADR-0006 (superseded)](../../_project/decisions/0006-new-project-skill.md)
