# Knowledge Surface Research Cockpit v1

> Status: Ready for prototype
>
> Date: 2026-04-17
>
> Purpose: upgrade the current static `docs/index.html` briefing into a local knowledge-system entrypoint for `erwin`.

---

## Problem

`erwin` now has a real knowledge pipeline:

- `radar` captures authoritative AI/agent signals.
- `feedback` records lived Harness and Codex usage pressure.
- `knowledge` turns those inputs into durable understanding.
- `proposals` converts mature understanding into bounded change candidates.

The current generated `docs/index.html` made that pipeline visible, but it is still mostly a static briefing page. It does not yet feel like a system where a human can inspect sources, follow themes, search notes, compare timeline movement, or hand the current context back to Codex without opening many Markdown files manually.

The next step should not be "make the HTML prettier." It should be a small local research cockpit that makes the repo-owned knowledge system usable every day.

## Evidence

The design should borrow explicit capabilities from mature knowledge and productivity systems without copying their product models:

| Reference | Useful pattern | Local interpretation for `erwin` |
|---|---|---|
| [NotebookLM](https://blog.google/technology/google-labs/notebooklm-new-features-december-2024/) | Sources, chat, and studio live in one research workspace. | Keep source reading, questions, and synthesis outputs adjacent instead of scattering them across files and Codex threads. |
| [Readwise Reader](https://readwise.io/read) | Power-reader workflow, search, and later processing. | Treat radar notes as a reading queue with search, filters, and focused Markdown reading. |
| [Linear](https://linear.app/) | Dense operational pulse and project status. | Make freshness, current streams, proposal status, and "what changed" visible at a glance. |
| [Raycast](https://www.raycast.com/) | Command-first interaction. | Add a command palette for search, navigation, source opening, and Codex handoff prompts. |
| [Heptabase](https://heptabase.com/) / [Obsidian Canvas](https://obsidian.md/canvas) / [Capacities](https://capacities.io/) / [Tana](https://tana.inc/) | Knowledge objects, relationships, visual organization. | Model themes, streams, notes, and proposals as connected objects without building a full graph database in v1. |

The user-facing pressure is also concrete: weekly and daily outputs exist, but the user cannot reliably "look at inbox" or browse a living knowledge surface from the app. The repo needs a durable, local, obvious entrypoint.

## Decision

Adopt a `Research Cockpit` shape for v1:

- Use a Vite + React + TypeScript frontend for the human knowledge surface.
- Use a small Go local server for read-only APIs and static asset serving.
- Continue using repo-owned Markdown and `docs/knowledge/index.json` as the source of truth.
- Keep automation responsible for writing radar files and refreshing generated indexes.
- Keep the knowledge surface responsible for reading, searching, filtering, navigating, and preparing Codex handoff context.

This is a prototype-level proposal, not a commitment to a long-lived product architecture. The first implementation should prove that the interaction model is useful before adding persistence, accounts, or external services.

## Non-goals

- Do not introduce a database in v1.
- Do not depend on an external SaaS service.
- Do not call an external LLM API from the browser or backend in v1.
- Do not build a full in-browser AI chat surface in v1.
- Do not replace Markdown as the durable repo record.
- Do not turn every daily radar note into a proposal.
- Do not build an infinite canvas or graph database in v1.

## Implementation Shape

### Frontend

Create `apps/knowledge-surface/` with Vite, React, and TypeScript.

The first UI should include:

- `Pulse`: default landing view with latest daily radar, GitHub Trending, weekly synthesis, active themes, freshness, and watch-next signals.
- `Reader`: Markdown note browsing with stream/layer/theme filters and source metadata.
- `Timeline`: chronological view across daily, weekly, trending, proposal, research, and feedback notes.
- `Threads`: theme-oriented view for durable topics such as `harness`, `coding agents`, `context engineering`, `agent evals`, `tool use`, and `observability`.
- `Proposals`: compact view of active proposal artifacts and their upstream anchors.
- `Command palette`: `Cmd/Ctrl+K` surface for search, navigation, opening notes, and copying a Codex handoff prompt.
- `Codex handoff`: generate a copyable prompt containing the selected note, stream, theme, and relevant file paths.

The visual direction should be a research cockpit, not a generic SaaS dashboard:

- Dark operational shell with high-contrast reading panels.
- Dense but calm status hierarchy inspired by Linear.
- Reading comfort inspired by power-reader tools.
- Explicit source and path affordances so the user always knows where repo truth lives.
- Responsive layout that preserves `Pulse`, `Search`, and `Reader` on small screens.

### Backend

Create `cmd/knowledge-surface/` as a Go local server.

Initial read-only API:

- `GET /`: serve the built knowledge surface.
- `GET /api/health`: return index presence, generated timestamp when available, and note counts.
- `GET /api/index`: return `docs/knowledge/index.json`.
- `GET /api/note?path=...`: return Markdown content for allowlisted docs paths.
- `GET /api/search?q=...&stream=...&layer=...`: perform lightweight local search over indexed notes.

Security boundary:

- Only allow reads under explicit repo-owned knowledge directories such as `docs/radar/`, `docs/knowledge/`, `docs/feedback/`, `docs/proposals/`, and `docs/research/`.
- Reject path traversal and absolute paths outside the repo.
- Do not expose write endpoints in v1.

### Data generator

Extend `scripts/build-knowledge-surface.py` so `docs/knowledge/index.json` becomes a stronger application contract, not just a byproduct of static HTML generation.

The index should carry enough structured data for:

- streams
- layers
- themes
- note paths
- note titles
- dates
- freshness
- latest stream pointers
- proposal and knowledge anchors when discoverable from existing docs

Keep `docs/index.html` as a fallback / launch page that points to the local cockpit and still works as a simple generated briefing when the local app is not running.

### Script

Create `scripts/knowledge-surface.sh` as the human entrypoint.

Expected behavior:

- Rebuild the knowledge index before launch unless explicitly skipped.
- Build or run the frontend as appropriate for local development.
- Start the Go local server.
- Print the local URL.

## Acceptance Criteria

The prototype is acceptable when:

- The user can open one local URL and see the current `erwin` field pulse within 30 seconds.
- The homepage clearly shows latest daily radar, latest GitHub Trending, latest weekly synthesis, active themes, and proposal movement.
- The user can search across indexed notes.
- The user can open and read Markdown notes without manually navigating the filesystem.
- The user can filter by stream, layer, and theme.
- The user can inspect a timeline of recent knowledge-system movement.
- The user can open proposal artifacts from the UI.
- The user can copy a Codex handoff prompt for the currently selected note or theme.
- The API cannot read outside the allowlisted docs paths.
- The implementation does not require a database or external service.

## Verification Plan

Proposal-only landing:

- Confirm this file is linked from `docs/proposals/README.md`.
- Confirm `_project/CURRENT.md` points future sessions at prototype implementation instead of more design discussion.
- Confirm `CHANGELOG.md` records the proposal landing.

Prototype implementation:

- Run the knowledge index generator and inspect `docs/knowledge/index.json`.
- Run frontend install/build/typecheck commands from `apps/knowledge-surface/`.
- Run `go test ./...`.
- Add backend tests for allowlisted note reads and path traversal rejection.
- Start the local surface through `scripts/knowledge-surface.sh`.
- Manually validate desktop and mobile layouts.

## Recommended Next Move

Open a dedicated implementation worktree and build the smallest useful Research Cockpit prototype:

1. Strengthen `docs/knowledge/index.json`.
2. Add the Go read-only API.
3. Add the Vite/React UI around `Pulse`, `Reader`, `Timeline`, and `Command palette`.
4. Keep `docs/index.html` as fallback.
5. Verify with local build, Go tests, and manual UI review.
