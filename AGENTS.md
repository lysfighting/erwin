# erwin — Project Context for Codex

> Read this first. Then read `_project/CURRENT.md` for what is live right now.

---

## What Is This Project?

`erwin` is an **AI/agent frontier exploration lab**.

The active lanes are:

1. **Radar** — authoritative AI/agent field intelligence, with `Harness` as one key theme
2. **Harness** — Codex-native orchestration operating model and implementation
3. **Skills** — reusable skills that still support the current agent-systems direction

Historical teamwork artifacts still exist, but they are no longer the repo center.

Install all erwin skills:

```bash
./install.sh
```

See `_project/GOALS.md` for current mission and `_project/CURRENT.md` for live status.

---

## Repository Structure

```
erwin/
├── README.md               ← 30-second repo entrypoint
├── AGENTS.md               ← YOU ARE HERE
├── CHANGELOG.md            ← Project version history
├── install.sh              ← Skills installer
│
├── docs/
│   ├── README.md           ← Docs index
│   ├── radar/              ← AI/agent radar knowledge base
│   ├── archive/            ← Doc archive boundary
│   └── research/           ← One-time analysis and research reports
│
├── harness/                ← Harness definitions and package boundary
├── scripts/                ← Harness CLI surfaces and utilities
├── cmd/                    ← Go CLI entrypoints
├── internal/               ← Harness runtime implementation
│
├── skills/                 ← Reusable skills
│   ├── harness/            ← Active
│   ├── knowledge/          ← Active
│   ├── cwork/              ← Active
│   ├── arch-review/        ← Active
│   ├── visual-test/        ← Active
│   ├── session-record/     ← Active
│   └── teamwork/           ← Legacy
│
├── _project/
│   ├── CURRENT.md          ← Live status
│   ├── GOALS.md            ← Mission and OKRs
│   ├── decisions/          ← ADRs
│   ├── ideas/              ← Idea pipeline
│   ├── issues/             ← Issue tracking
│   ├── archive/sessions/   ← Historical session logs
│   └── sessions -> archive/sessions
│
├── archive/                ← Top-level archive boundary
└── reference/              ← Source material (READ-ONLY)
```

---

## Working Model

- `README.md` explains the repo in 30 seconds
- `_project/CURRENT.md` is the live operational handoff
- `docs/README.md` is the documentation map
- `docs/radar/` is the field-intelligence layer
- `harness/`, `scripts/`, `cmd/`, and `internal/` are the main implementation surface

---

## Session Protocol

**Start**: read `AGENTS.md` → `_project/CURRENT.md` → `docs/README.md`

**End**:

- update `_project/CURRENT.md`
- write session history if needed under `_project/sessions/` (compat path to archive)
- update `CHANGELOG.md`
- commit intentionally

---

## Key Files

| What | Where |
|------|-------|
| Live status | `_project/CURRENT.md` |
| Mission / goals | `_project/GOALS.md` |
| Repo entrypoint | `README.md` |
| Docs index | `docs/README.md` |
| Radar entrypoint | `docs/radar/README.md` |
| Radar sources | `docs/radar/sources.md` |
| Harness docs | `docs/harness-documentation-map.md` |
| Harness package boundary | `harness/`, `scripts/`, `cmd/`, `internal/` |
| Skills | `skills/` |
| Codex work loop skill | `skills/cwork/SKILL.md` |
| Historical boundary | `archive/README.md` |
| Source material | `reference/drive-skill/commands/drive.md` |

---

## What Not To Do

- Don't edit `reference/`
- Don't rely on session memory
- Don't treat historical teamwork material or showcase-era plans as the default source of truth
- Don't let tooling-watch noise outrank first-party AI/agent sources
