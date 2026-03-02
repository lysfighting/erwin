# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-02
- **Last Session**: [2026-03-02-K](sessions/2026-03-02-K.md) (Session K)

## Current Phase

**Phase**: Teamwork v2 implemented, ready for validation

## What Was Just Completed

1. **Teamwork v2 Implementation** (Session K continued, /self-drive → /drive)
   - 4 skill files written: `/team`, `/team-claim`, `/team-drive`, `/team-ship`
   - v1 skills removed: `/teamwork`, `/team-get-mission`, `/team-complete`
   - README.md, CHANGELOG.md updated (skill-level + root-level)
   - install.sh verified — 4 symlinks correct, stale v1 symlinks cleaned
   - Architecture: `docs/teamwork-v2-architecture.md` (ADR-0010)

## Active Task / What's In Progress

Teamwork v2 skills complete. Next: validate on vi-agent project, or address stale docs.

## Next Session Should Start With

1. Read this file
2. **Validate teamwork v2** on `../vi-agent` project (`/team` init)
3. OR clean up stale docs (teamwork-architecture.md v1 deprecation, install.sh orphan cleanup)
4. OR continue **V2 Showcase** (Phase 2/3 remaining tasks below)

## Implementation Roadmap

### v0.5.0 — V2 Phase 2: Deep Interaction (in progress)
- [x] OriginTrace (d3-force philosophy graph, ACT II)
- [x] Bundle code splitting (627KB → 221KB, React.lazy + vendor chunks)
- [x] Code health sweep (ESLint clean, semantic HTML)
- [ ] ExecutionSimulator: message flow animation enhancement
- [ ] SwarmTopologySection: add interactive message flow

### v1.0.0 — V2 Phase 3: Creation Path
- [x] SkillBuilder (form → skill.md download, ACT IV)
- [ ] DesignPatternLibrary (3-5 reusable pattern cards)
- [ ] CaseStudies (2-3 skill variant examples)

### v1.1.0 — Multi-Skill Platform (post-v1.0 backlog)
- [ ] SkillGallery section (one demo card per released skill)
- [ ] Hero section rebrand ("erwin: AI skill incubator")
- [ ] Skills index (table of all skills with install commands)

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce

## File Locations

| What | Where |
|------|-------|
| Project goals | `_project/GOALS.md` |
| Idea garden | `_project/ideas/` |
| ADRs | `_project/decisions/` |
| V2 architecture | `docs/drive-showcase-v2-architecture.md` |
| drive-skill source | `reference/drive-skill/commands/drive.md` |
| App root | `showcase/src/App.tsx` |
