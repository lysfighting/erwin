# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-03
- **Last Session**: [2026-03-03-P](sessions/2026-03-03-P.md) (Session P — E2E validation)

## Current Phase

**Phase**: Teamwork v2.4.1 — E2E validated on vi_agent

## What Was Just Completed

**Session P** — E2E validation of teamwork v2.4 on vi_agent. Two bugs found and fixed (v2.4.1):

1. **tw-config.sh inline comment bug** — Values with trailing inline comments (e.g., `"mission-contract"  # label...`) were returned with garbage. Fixed `strip_val()` to extract quoted content directly, ignoring comments.
2. **mc_label dual-schema support** — Skills called `mc_label` (flat) but vi_agent config has it at `github.mc_label` (nested). team-claim and team-issue now try both paths with fallback.
3. **Full E2E flow verified** — Issue #27 created, claimed, driven (file created + bidirectional sync), shipped (PR #28), all artifacts cleaned up. vi_agent left in original state.

## Active Task / What's In Progress

None. Sessions O+P complete.

## Next Session Should Start With

1. Read this file
2. Continue V2 Phase 3 showcase tasks (DesignPatternLibrary, CaseStudies)
3. **Future teamwork missions**:
   - `/team-release` skill (close Milestone → git tag → `gh release create --generate-notes`)
   - `/team-drive` wrapping full `/drive` (thin wrapper, skip Phase 0)

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
- [x] Teamwork v2.4 — stability hardening + bidirectional sync + tw-config.sh
- [x] Teamwork v2.4.1 — E2E validated on vi_agent + mc_label dual-schema + inline comment fix

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Teamwork: milestone creation (if it doesn't exist) needs manual step before /team-issue (non-fatal warning works)
- Teamwork: team.md (dashboard) still at v2.3.2 — team board feature deferred

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
