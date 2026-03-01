# ADR-0001: Monorepo Structure for erwin

- **Date**: 2026-03-02
- **Status**: Accepted
- **Session**: docs/sessions/2026-03-02.md

## Context

The erwin project contains:
1. `drive-showcase` — a React web app
2. `reference-skills/drive-skill` — source material / reference implementation
3. `docs` — documentation and architecture artifacts

These could be managed as separate repositories or as a monorepo.

## Decision

Use a single git repository (monorepo) rooted at `/erwin/`.

## Consequences

**Positive**:
- Single `git log` shows the full project history
- CHANGELOG.md at root covers everything
- CLAUDE.md at root is loaded for all work in this project
- Cross-cutting concerns (docs + code changes) land in one commit
- No need to coordinate version bumps across repos

**Negative**:
- `drive-showcase/node_modules` is large (133MB); handled via .gitignore
- Reference skills are "owned" elsewhere (KC's GitHub); tracked here as read-only copy

## Alternatives Considered

- **Separate repos**: Rejected. Cross-cutting changes (e.g., architecture doc + code stub) would require multi-repo commits. Overhead not worth it for a 1-2 person project.
- **Git submodule for drive-skill**: Rejected. The drive-skill source needs to be readable by Claude at session start; submodule complexity not justified.
