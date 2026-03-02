# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-02
- **Last Session**: [2026-03-02-M](sessions/2026-03-02-M.md) (Session M)

## Current Phase

**Phase**: Teamwork v2.1 — architect-reviewed, tested, ready for E2E validation

## What Was Just Completed

1. **Self-Drive Session M** — 2 missions completed:
   - **Phase 2 Animation Enhancements**: ExecutionSimulator topology + message flow dots, SwarmTopologySection interactive message flow (auto-play + click-to-focus), /visual-test verified
   - **Code Health Sweep**: Fixed ESLint error (SwarmTopology setState-in-effect), removed orphaned DemoSection.tsx (24KB), removed 9 dead demo.note* i18n keys. 0 ESLint errors, 0 TS errors.

2. **Previous sessions (I-L)**:
   - Teamwork v2.1 architect review + 13 bug fixes (Session L)
   - SkillBuilder MVP complete (Session I) — ACT IV form + download
   - Bundle splitting done (Session J) — 627KB → 221KB
   - Teamwork v2.0 design + implementation (Session K)

## Active Task / What's In Progress

Self-Drive loop active. Phase 2 complete. Analyzing project for next proposals.

## Next Session Should Start With

1. Read this file
2. Continue Self-Drive loop (if session continues) or start new proposal round
3. OR **E2E validate** teamwork on `../vi_agent`
4. OR continue V2 Phase 3 remaining tasks

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
- [ ] Teamwork E2E validation on vi_agent

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Placeholder naming inconsistent across teamwork skill files ({N} vs {ISSUE_NUMBER} vs {issue})

## File Locations

| What | Where |
|------|-------|
| Project goals | `_project/GOALS.md` |
| Idea garden | `_project/ideas/` |
| ADRs | `_project/decisions/` |
| V2 architecture | `docs/drive-showcase-v2-architecture.md` |
| Teamwork architecture | `docs/teamwork-v2-architecture.md` |
| Teamwork skills | `skills/teamwork/` |
| App root | `showcase/src/App.tsx` |
