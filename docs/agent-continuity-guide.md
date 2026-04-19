# Agent Continuity Guide

> Status: Active
>
> Date: 2026-04-19
>
> Purpose: provide a vendor-neutral continuity contract so Codex, Claude, or another agent can restart work from repository files instead of chat memory.

---

## Why This Exists

`erwin` has several connected systems: Radar, Feedback, Knowledge, Proposals, Implementation, Harness, Skills, automations, and the knowledge surface.

No agent should need access to an old Codex thread, Claude conversation, inbox card, or automation run to understand the project.

This guide defines the durable read order, routing rules, and end-of-session writeback discipline that keep the project portable across agents and across closed sessions.

## Start Protocol

Every agent should start with this sequence:

1. Read the agent-specific entrypoint:
   - Codex: [`AGENTS.md`](../AGENTS.md)
   - Claude: [`CLAUDE.md`](../CLAUDE.md)
2. Read [`_project/CURRENT.md`](../_project/CURRENT.md) for live status.
3. Read this guide.
4. Read [`docs/erwin-operating-model.md`](erwin-operating-model.md) for the five-layer loop.
5. Read [`docs/README.md`](README.md) for the documentation map.
6. If the task crosses systems, inspect [`docs/knowledge/index.json`](knowledge/index.json) and [`docs/knowledge/active-threads.md`](knowledge/active-threads.md).
7. If the task touches Harness, read [`docs/harness-documentation-map.md`](harness-documentation-map.md) and [`docs/harness-session-handoff.md`](harness-session-handoff.md).

Do not begin cross-layer work after reading only a local implementation file.

## Durable System Map

| System | Role | Canonical Entry | Must Update When |
|---|---|---|---|
| `Radar` | Intake of authoritative AI/agent field signal. | [`docs/radar/README.md`](radar/README.md) | A daily, weekly, trending, or paper-watch note changes. |
| `Feedback` | Normalized real Harness, Codex, automation, or repo-process friction. | [`docs/feedback/README.md`](feedback/README.md) | Repeated usage pressure or failure mode becomes visible. |
| `Knowledge` | Durable interpretation and active threads. | [`docs/knowledge/README.md`](knowledge/README.md) | Evidence changes what the repo should remember. |
| `Proposals` | Reviewed change candidates. | [`docs/proposals/README.md`](proposals/README.md) | A mature thread justifies `watch`, `encode-as-eval`, `prototype`, or `adopt`. |
| `Implementation` | Actual changes to Harness, Skills, evals, docs, automation, or the knowledge surface. | repo code and skills | A reviewed proposal, explicit internal gap, or bugfix is implemented. |

`Harness` is an implementation target and a feedback source.
`Skills` are reusable implementation capabilities.
`Automations` are execution mechanisms.
`Knowledge Surface` is the human-facing reading/navigation layer.
None of these replace the five-layer loop.

## What To Trust

When artifacts disagree, use this order:

1. Machine-readable runtime or package truth for runtime shape.
   - [`harness/adapter-manifest.yml`](../harness/adapter-manifest.yml)
   - [`internal/harness/runtime.go`](../internal/harness/runtime.go)
2. Reviewed proposal artifacts for change-candidate status.
   - [`docs/proposals/harness-proposal-ledger.md`](proposals/harness-proposal-ledger.md)
   - [`docs/proposals/harness-review-gates.md`](proposals/harness-review-gates.md)
3. Active knowledge threads for durable interpretation.
   - [`docs/knowledge/active-threads.md`](knowledge/active-threads.md)
4. `_project/CURRENT.md` for live handoff, if it is fresh.
5. Changelog and recent commits for chronological history.
6. Archived sessions and research notes only as supporting context.

Chat memory is not a source of truth.

## Task Routing

Use the smallest correct layer:

| If You See | Route To | Do Not |
|---|---|---|
| New external AI/agent signal | `Radar` | Turn it directly into Harness runtime change. |
| Repeated local friction | `Feedback` | Hide it in a session summary only. |
| Stable interpretation across evidence | `Knowledge` | Leave it in raw research or daily notes. |
| Mature pressure for change | `Proposals` | Create a proposal from novelty alone. |
| Clear bug, explicit internal gap, or reviewed proposal | `Implementation` | Bypass handoff, changelog, or verification. |

If unsure, update Knowledge or Feedback first.

## Harness Evolution Rule

Harness should evolve through the five-layer loop:

`Radar / Feedback -> Knowledge -> Proposals -> Implementation -> Feedback / Knowledge`

Before strategic Harness work, read:

- [`docs/harness-explainer.md`](harness-explainer.md)
- [`docs/harness-documentation-map.md`](harness-documentation-map.md)
- [`docs/harness-session-handoff.md`](harness-session-handoff.md)
- [`docs/harness-learning-loop.md`](harness-learning-loop.md)
- [`docs/proposals/harness-proposal-ledger.md`](proposals/harness-proposal-ledger.md)
- [`docs/proposals/harness-eval-candidate-ledger.md`](proposals/harness-eval-candidate-ledger.md)

Bugfixes and already-decided internal gaps can go directly to implementation, but the session must still explain the evidence and update the durable handoff surfaces.

## Knowledge Network Contract

The generated knowledge network lives at [`docs/knowledge/index.json`](knowledge/index.json).

It contains:

- `knowledge_network.artifacts`
- `knowledge_network.links`
- `knowledge_network.threads`
- `knowledge_network.proposal_records`

Regenerate it whenever Radar, Feedback, Knowledge, Proposals, agent entrypoints, or canonical implementation anchors change:

```bash
python3 scripts/build-knowledge-surface.py
```

Do not edit generated surfaces by hand.

Generated outputs:

- [`docs/index.html`](index.html)
- [`docs/knowledge/briefing.md`](knowledge/briefing.md)
- [`docs/knowledge/index.json`](knowledge/index.json)

## End Protocol

Before closing a meaningful session:

1. Update the durable artifact for the layer you changed.
2. If any cross-layer relation changed, update active threads or proposal ledgers.
3. If real usage friction appeared, add or update a feedback case.
4. If generated knowledge surfaces are affected, run `python3 scripts/build-knowledge-surface.py`.
5. Update [`_project/CURRENT.md`](../_project/CURRENT.md).
6. Update [`CHANGELOG.md`](../CHANGELOG.md).
7. Run the smallest useful verification gate.
8. Commit intentionally.

Only write a session log when it carries durable operational value beyond `CURRENT.md`, `CHANGELOG.md`, and the changed docs.

## Portability Rules

- Keep decisions in repo files, not chat.
- Keep agent-specific entrypoints thin and aligned with this guide.
- If `AGENTS.md` and `CLAUDE.md` disagree, fix the stale file instead of relying on the current agent's memory.
- Prefer relative Markdown links inside docs so cloned worktrees remain portable.
- Do not depend on Codex inbox cards, Claude thread history, automation session logs, or local UI state for project continuity.
- If you discover stale guidance, update the misleading entrypoint before ending the session.

