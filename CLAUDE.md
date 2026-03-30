# erwin — Project Context for Claude

> Read this first. Then read `_project/CURRENT.md` for what is live right now.

---

## What Is This Project?

`erwin` is an **AI/agent frontier exploration lab**.

The active lanes are:

1. **Radar** — authoritative AI/agent field intelligence, with `Harness` as one key theme
2. **Harness** — Codex/Claude-native orchestration operating model and implementation
3. **Skills** — reusable skills that still support the current agent-systems direction

Historical teamwork artifacts still exist, but they are no longer the repo center.

Install all erwin skills:

```bash
./install.sh
```

---

## Repository Structure

```
erwin/
├── README.md               ← 30-second repo entrypoint
├── CLAUDE.md               ← YOU ARE HERE
├── CHANGELOG.md            ← Project version history
├── docs/                   ← docs map + radar + research + archive boundary
├── harness/                ← Harness definitions and package boundary
├── scripts/                ← Harness shell surfaces
├── cmd/                    ← Go CLI entrypoints
├── internal/               ← Harness runtime implementation
├── skills/                 ← reusable skills, including legacy material
├── _project/               ← live status + goals + ideas + ADRs + archived sessions
├── archive/                ← top-level archive boundary
└── reference/              ← source material (READ-ONLY)
```

---

## Start And End

**Start**: read `README.md` → `_project/CURRENT.md` → `docs/README.md`

**End**:

- update `_project/CURRENT.md`
- write session history if useful under `_project/sessions/` (compat path to archive)
- update `CHANGELOG.md`
- commit intentionally

---

## Key Files

| What | Where |
|------|-------|
| Live status | `_project/CURRENT.md` |
| Mission / goals | `_project/GOALS.md` |
| Docs index | `docs/README.md` |
| Radar entrypoint | `docs/radar/README.md` |
| Harness docs | `docs/harness-documentation-map.md` |
| Skills | `skills/` |
| Historical boundary | `archive/README.md` |
| Source material | `reference/drive-skill/commands/drive.md` |

---

## What Not To Do

- Don't edit `reference/`
- Don't rely on session memory
- Don't treat historical teamwork material or showcase-era plans as the default source of truth
