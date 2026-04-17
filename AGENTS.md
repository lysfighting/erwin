# erwin — Project Context for Codex

> Read this first. Then read `_project/CURRENT.md` for what is live right now.

---

## What Is This Project?

`erwin` is an **AI/agent frontier exploration lab**.

The active operating loop has five layers:

1. **Radar** — authoritative AI/agent field intelligence intake
2. **Feedback** — normalized real Harness / Codex / repo usage pressure
3. **Knowledge** — durable synthesis, active threads, and interpretation
4. **Proposals** — reviewed change candidates
5. **Implementation** — concrete Harness / Skills / eval / docs / automation changes

`Harness` is the main implementation target and a major feedback source, not a separate shortcut around the knowledge loop. `Skills` remain useful reusable capabilities that support the current agent-systems direction.

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
│   ├── erwin-operating-model.md ← Five-layer operating model
│   ├── radar/              ← AI/agent radar knowledge base
│   ├── feedback/           ← Normalized usage feedback
│   ├── knowledge/          ← Durable synthesis and active threads
│   ├── proposals/          ← Reviewed change candidates
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
- `docs/erwin-operating-model.md` explains the five-layer knowledge-to-action loop
- `docs/README.md` is the documentation map
- `docs/knowledge/index.json` is the generated machine-readable knowledge network for cross-layer navigation
- `docs/radar/` is the field-intelligence layer
- `docs/feedback/`, `docs/knowledge/`, and `docs/proposals/` move signal toward reviewed action
- `harness/`, `scripts/`, `cmd/`, and `internal/` are the main implementation surface

---

## Session Protocol

**Start**: read `AGENTS.md` → `_project/CURRENT.md` → `docs/erwin-operating-model.md` → `docs/README.md`

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
| Five-layer operating model | `docs/erwin-operating-model.md` |
| Docs index | `docs/README.md` |
| Knowledge network index | `docs/knowledge/index.json` |
| Radar entrypoint | `docs/radar/README.md` |
| Radar sources | `docs/radar/sources.md` |
| Feedback entrypoint | `docs/feedback/README.md` |
| Knowledge entrypoint | `docs/knowledge/README.md` |
| Proposal entrypoint | `docs/proposals/README.md` |
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
