# /new-project

Bootstrap any project with the complete AI-assisted project intelligence layer in ~5 minutes.

## What it does

Scaffolds a new project with the **erwin project pattern**: version management, session protocols, ADR system, idea garden, and Claude auto-memory. Derived from erwin's own `_project/` structure.

## Usage

```
/new-project
```

Run from the project's root directory. The skill asks 4 questions interactively, then creates everything.

## What gets created

```
_project/
├── GOALS.md           ← Vision + OKRs template
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
./install.sh new-project
```

Or to install all erwin skills:
```bash
./install.sh
```

## Origin

Derived from erwin's `_project/` structure, formalized as a reusable skill in Session F (2026-03-02).
See: [ADR-0006](../../_project/decisions/0006-new-project-skill.md)
