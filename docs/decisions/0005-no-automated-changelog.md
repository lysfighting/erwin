# ADR-0005: Manual Changelog Over Automated Tooling

- **Date**: 2026-03-02
- **Status**: Accepted
- **Session**: docs/sessions/2026-03-02.md

## Context

Two approaches to changelog management:
1. **Manual**: Write CHANGELOG.md entries by hand, following Keep a Changelog format
2. **Automated**: Use conventional commits + semantic-release/standard-version to auto-generate

## Decision

Use **manual changelogs** following the drive-skill/drive-showcase CHANGELOG pattern.

Key elements:
- Keep a Changelog format (Unreleased → version → date)
- Semantic versioning (MAJOR.MINOR.PATCH)
- Include "Rationale" section (borrowed from drive-showcase CHANGELOG style)
- Conventional commit messages in git log for human readability

No automated tooling (no semantic-release, no commitizen, no husky).

## Consequences

**Positive**:
- Zero tool configuration overhead
- Changelog entries can be narrative and informative (not just auto-generated summaries)
- The "Rationale" section (unique to this project) can't be auto-generated
- Works immediately, no setup required

**Negative**:
- Manual discipline required to update CHANGELOG on every meaningful change
- No automatic version bump enforcement

## Rule

**Convention for this project**: Every PR/session that adds a feature or makes a notable change
MUST include CHANGELOG.md updates as part of the work, not as an afterthought.

Claude's session workflow:
1. Before starting: read CURRENT_STATE.md
2. During work: note changes
3. Before ending: update CHANGELOG.md + CURRENT_STATE.md + session log

## Alternatives Considered

- **semantic-release**: Rejected. Requires CI pipeline, npm publish flow, token setup. Overkill for a single-developer project with AI pair programming.
- **standard-version**: Rejected. Same overhead, deprecated in favor of release-please.
- **conventional commits enforced by husky**: Rejected. Adds friction with no proportional benefit at this team size.
