---
# AI session workflow rules for erwin project
---

# AI Session Workflow Protocol

## Session Start (mandatory, 5 min)

1. Read `_project/CURRENT.md` — what's in progress
2. Read `_project/GOALS.md` — verify we're aligned with north star
3. Read relevant ADR(s) if working on a specific component

## During Session

- Discovered bugs/tasks → create `_project/issues/ISS-NNNN-title.md`
- Major decisions → create `_project/decisions/NNNN-title.md`
- Good skill ideas → add to `_project/ideas/INBOX.md` immediately

## Session End (mandatory, 10 min)

1. Update `_project/CURRENT.md`
2. Write `_project/sessions/YYYY-MM-DD.md`
3. Update `CHANGELOG.md` if features completed
4. Run: `git add -A && git commit -m "type(scope): description"`

## File Size Limits

| File | Target | Max |
|------|--------|-----|
| CLAUDE.md | <150 lines | 200 lines |
| _project/CURRENT.md | <80 lines | 120 lines |
| ADR files | 50-100 lines | 150 lines |
| Session logs | 50-150 lines | 300 lines |

## What Not To Do

- Don't rely on session memory for continuity — always read CURRENT.md first
- Don't skip the session-end update — stale CURRENT.md = broken continuity
- Don't add large images or binary files to docs/ or _project/
- Don't edit reference/ directory files (read-only source material)
