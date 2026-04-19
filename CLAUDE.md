# erwin — Project Context for Claude

> Read this first. Then read `_project/CURRENT.md` for what is live right now.

---

## What Is This Project?

`erwin` is an **AI/agent frontier exploration lab**.

The active operating loop has five layers:

1. **Radar** — authoritative AI/agent field intelligence intake
2. **Feedback** — normalized real Harness / Claude / Codex / repo usage pressure
3. **Knowledge** — durable synthesis, active threads, and interpretation
4. **Proposals** — reviewed change candidates
5. **Implementation** — concrete Harness / Skills / eval / docs / automation changes

`Harness` is the main implementation target and a major feedback source, not a separate shortcut around the knowledge loop.
`Skills` remain useful reusable capabilities that support the current agent-systems direction.

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
├── CLAUDE.md               ← YOU ARE HERE
├── AGENTS.md               ← Codex entrypoint; should stay aligned with this file
├── CHANGELOG.md            ← Project version history
├── install.sh              ← Skills installer
│
├── docs/
│   ├── README.md           ← Docs index
│   ├── agent-continuity-guide.md ← Cross-agent continuity contract
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
├── _project/               ← Live status, goals, decisions, issues, archived sessions
├── archive/                ← Top-level archive boundary
└── reference/              ← Source material (READ-ONLY)
```

---

## Working Model

- `README.md` explains the repo in 30 seconds
- `_project/CURRENT.md` is the live operational handoff
- `docs/agent-continuity-guide.md` is the vendor-neutral cross-agent restart contract
- `docs/erwin-operating-model.md` explains the five-layer knowledge-to-action loop
- `docs/README.md` is the documentation map
- `docs/knowledge/index.json` is the generated machine-readable knowledge network for cross-layer navigation
- `docs/radar/` is the field-intelligence layer
- `docs/feedback/`, `docs/knowledge/`, and `docs/proposals/` move signal toward reviewed action
- `harness/`, `scripts/`, `cmd/`, and `internal/` are the main implementation surface

---

## Session Protocol

**Start**: read `CLAUDE.md` → `_project/CURRENT.md` → `docs/agent-continuity-guide.md` → `docs/erwin-operating-model.md` → `docs/README.md`

If the task touches Harness, also read:

- `docs/harness-documentation-map.md`
- `docs/harness-session-handoff.md`
- `docs/harness-learning-loop.md`

**End**:

- update the durable artifact for the layer changed
- refresh generated knowledge surfaces when needed with `python3 scripts/build-knowledge-surface.py`
- update `_project/CURRENT.md`
- update `CHANGELOG.md`
- commit intentionally

Write session history only when it carries durable operational value beyond the changed docs, `CURRENT.md`, and `CHANGELOG.md`.

---

## Key Files

| What | Where |
|------|-------|
| Live status | `_project/CURRENT.md` |
| Mission / goals | `_project/GOALS.md` |
| Repo entrypoint | `README.md` |
| Cross-agent continuity | `docs/agent-continuity-guide.md` |
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
| Historical boundary | `archive/README.md` |
| Source material | `reference/drive-skill/commands/drive.md` |

---

## What Not To Do

- Don't edit `reference/`
- Don't rely on session memory
- Don't treat historical teamwork material or showcase-era plans as the default source of truth
- Don't let tooling-watch noise outrank first-party AI/agent sources
- Don't mutate Harness runtime from fresh radar signal without Feedback / Knowledge / Proposal routing
