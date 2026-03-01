# Session Logs

This directory contains per-session artifacts from AI-assisted development conversations.

## Purpose

Since Claude Code sessions are stateless, each session's context, decisions, and outputs
are preserved here as markdown files. This directory IS the "memory" of the project.

## Naming Convention

```
YYYY-MM-DD.md              # single session per day
YYYY-MM-DD-N.md            # multiple sessions in one day (N = 1, 2, 3...)
```

## Session Log Template

Each session log captures:

```markdown
# Session: YYYY-MM-DD

## Context (What I knew coming in)
- Read CURRENT_STATE.md: [what it said]
- Current version: vX.Y.Z
- In-progress: [what was in progress]

## Work Done
- [bullet list of what was accomplished]

## Key Decisions Made
- [decision] → see ADR-NNNN if formal enough to warrant one

## Artifacts Created/Modified
- [file path] — [what changed]

## What's Next
- [top 3 tasks for next session]

## CURRENT_STATE.md Updates
- [what was updated in CURRENT_STATE.md]
```

## Session Index

| Date | Summary | Key Decisions | Version |
|------|---------|---------------|---------|
| [2026-03-02](2026-03-02.md) | Project setup: git init, version management system, V2 architecture design | ADR 0001-0005 | 0.2.0 |
