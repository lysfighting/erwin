# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-03
- **Last Session**: [2026-03-03-N](sessions/2026-03-03-N.md) (Session N)

## Current Phase

**Phase**: Teamwork v2.3 — 4 issues addressed in self-drive session, v2.3.3

## What Was Just Completed

**Self-Drive Session N** — 4 teamwork issues fixed across 3 commits:

1. **v2.3.1** — `/team-ship` status verification: label update now verified via `gh issue view`, output shows actual labels (not assumed). Dashboard gets UNDER REVIEW section for `status:review` issues.
2. **v2.3.2** — Version management + issue format: milestone-based version tracking (`versions.current` → GitHub Milestone assignment on create), canonical field names (`Success Criteria`, `Context & References`, `### Objective`), dashboard shows milestone progress.
3. **v2.3.3** — Script-based skill pilot: `scripts/setup-github-labels.sh` (idempotent, reads config), `install.sh` now installs `scripts/*.sh` to `$TARGET_DIR/scripts/`, `team.md` calls script first with fallback.

## Active Task / What's In Progress

None. Session N complete.

## Next Session Should Start With

1. Read this file
2. **E2E validate** teamwork v2.3 on `../vi_agent` — run full flow: /team → /team-issue → /team-claim → /team-drive → /team-ship
3. OR continue V2 Phase 3 showcase tasks (DesignPatternLibrary, CaseStudies)
4. OR: Issue #3 still has room to grow — add more scripts (e.g., `tw-dashboard.sh`)

## Implementation Roadmap

### v0.5.0 — V2 Phase 2: Deep Interaction (in progress)
- [x] OriginTrace (d3-force philosophy graph, ACT II)
- [x] Bundle code splitting (627KB → 221KB, React.lazy + vendor chunks)
- [x] Code health sweep (ESLint clean, dead code removal)
- [x] ExecutionSimulator: message flow animation enhancement
- [x] SwarmTopologySection: add interactive message flow

### v1.0.0 — V2 Phase 3: Creation Path
- [x] SkillBuilder (form → skill.md download, ACT IV)
- [ ] DesignPatternLibrary (3-5 reusable pattern cards)
- [ ] CaseStudies (2-3 skill variant examples)

### Skills Factory
- [x] Teamwork v2.0 — 4 GitHub-first skills
- [x] Teamwork v2.1 — 6 vi-agent features + 13 bug fixes
- [x] Teamwork v2.3 — status fix + version management + format standardization + scripts
- [ ] Teamwork E2E validation on vi_agent

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Placeholder naming inconsistent across teamwork skill files ({N} vs {ISSUE_NUMBER} vs {issue})
- Teamwork v2.3: milestone creation (if it doesn't exist) needs manual step before /team-issue

## File Locations

| What | Where |
|------|-------|
| Project goals | `_project/GOALS.md` |
| Idea garden | `_project/ideas/` |
| ADRs | `_project/decisions/` |
| V2 architecture | `docs/drive-showcase-v2-architecture.md` |
| Teamwork architecture | `docs/teamwork-v2-architecture.md` |
| Teamwork skills | `skills/teamwork/` |
| Teamwork scripts | `skills/teamwork/scripts/` |
| App root | `showcase/src/App.tsx` |
