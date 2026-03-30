# Archive Manifest — Round 1

This manifest records the first safe slimming pass for `erwin`.

## Archived Now

### Project History

- `_project/archive/sessions/`
  - reason: historical session logs remain useful, but are too noisy for the default project surface
  - compatibility: `_project/sessions` remains as a symlink
  - current policy: keep only logs still referenced by active docs, skills, changelog, or live project handoff

### Historical Docs

- `docs/archive/teamwork-ai-manual.md`
- `docs/archive/teamwork-design-constitution.md`
- `docs/archive/team-doctor-design.md`
- `docs/archive/drive-showcase-v2-architecture.md`
- `docs/archive/git-model-v0.1.0.md`

Compatibility is preserved through symlinks at the original `docs/*.md` paths.

### Historical Docs Deleted In Round 2

- `docs/archive/teamwork-v2-architecture.md`
- `docs/archive/teamwork-architecture.md`
- `docs/archive/showcase-v3-architecture.md`
- `docs/archive/visual-enhancement-architecture.md`
- `docs/archive/project-management-architecture.md`

Reason: these files had no remaining practical reference value, or were only mentioned by historical changelog text.

### Curated Session Keep Set

The current session archive intentionally keeps only the following files:

- `2026-03-02.md`
- `2026-03-13-AL.md`
- `2026-03-13-AM.md`
- `2026-03-18-AM.md`
- `2026-03-28-A.md`
- `2026-03-28-B.md`
- `2026-03-28-C.md`
- `2026-03-29-A.md`
- `2026-03-29-AL.md`
- `2026-03-29-AM.md`
- `README.md`

Everything else from the previous session archive was deleted in round 2 because it was no longer referenced by the active repo surface.

## Archive Candidates

### Demo / Product Surfaces

- `showcase/`
  - current state: still present at repo top level
  - suggested next step: extract anything still useful, then delete or move physically into `archive/`
- `teamwork-demo/`
  - current state: deleted in round 2
  - reason: no active references beyond repo narrative and no current implementation value

### Legacy Skill Surfaces

- `skills/teamwork/`
  - current state: legacy but still referenced
  - suggested next step: keep in place until the remaining useful ideas are extracted or retired

## Keep Active

- `docs/radar/`
- `harness/`
- `scripts/`
- `cmd/`
- `internal/`
- `skills/harness/`
- selected active skills under `skills/`
- `_project/CURRENT.md`
- `_project/GOALS.md`
- `_project/decisions/`

## Rules For Round 2

- do not delete archive candidates before extracting anything still useful
- prefer moves with compatibility shims when cross-references remain high
- delete only after one pass of human review over the manifest

Round 2 review outcome:

- low-value session history was deleted aggressively
- `teamwork-demo/` was deleted
- low-value historical docs were deleted instead of preserved indefinitely
- `showcase/` and `skills/teamwork/` remain the main unresolved legacy surfaces
