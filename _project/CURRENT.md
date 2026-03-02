# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-02
- **Last Session**: [2026-03-02](sessions/2026-03-02.md) (Session H)

## Current Phase

**Phase**: V2 Showcase Implementation — Phase 2 in progress + erwin framework identity solidifying

## What Was Just Completed

1. **ADR-0008: /new-project → /erwin-init rename + redesign** (Session H, Drive Mode)
   - Skill identity reframed: from scaffolding tool to cognitive infrastructure installer
   - Zero-question Phase 1, optional Phase 2 GOALS.md conversation
   - Naming convention: `erwin-` prefix for framework-core skills only
   - 18 files changed, all doc references updated, install.sh verified
   - ADR-0006 marked superseded, session logs preserved as history

2. **erwin as vibe coding framework** — identity discussion (Session H, /architect)
   - erwin = vibe coding framework with cognitive self-evolution
   - `/erwin-init` = the framework installer (not a scaffolding tool)
   - Hard boundary: never touches code, never asks tech stack

## Active Task / What's In Progress

ADR-0008 implementation done. Next: SkillBuilder (OKR KR3) or bundle optimization.

## Next Session Should Start With

1. Read this file
2. **V2 Phase 3 priority**: SkillBuilder (OKR KR3, currently 0%)
   - ACT IV builder: form → skill.md download
   - See `docs/drive-showcase-v2-architecture.md` → 组件四: SkillBuilder
3. **Tech health**: Bundle code splitting (604KB → ~250KB, React.lazy for ACT III)
4. **Optional**: Session-end automation via Claude Code hooks

## Implementation Roadmap

### v0.3.0 — Project Structure ✅
### v0.4.0 — V2 Phase 1: Four-Act Structure ✅
### v0.5.0 — V2 Phase 2: Deep Interaction (in progress)
- [x] OriginTrace (d3-force philosophy graph, ACT II)
- [ ] ExecutionSimulator: message flow animation enhancement
- [ ] SwarmTopologySection: add interactive message flow

### v1.0.0 — V2 Phase 3: Creation Path
- [ ] SkillBuilder (form → skill.md download, ACT IV)
- [ ] DesignPatternLibrary (3-5 reusable pattern cards)
- [ ] CaseStudies (2-3 skill variant examples)

### v1.1.0 — Multi-Skill Platform (post-v1.0 backlog)
- [ ] SkillGallery section (one demo card per released skill)
- [ ] Hero section rebrand ("erwin: AI skill incubator")
- [ ] Skills index (table of all skills with install commands)

## Known Issues / Blockers

- 604KB single-chunk bundle — needs code splitting before adding more components

## File Locations

| What | Where |
|------|-------|
| Project goals | `_project/GOALS.md` |
| Idea garden | `_project/ideas/` |
| ADRs | `_project/decisions/` |
| V2 architecture | `docs/drive-showcase-v2-architecture.md` |
| drive-skill source | `reference/drive-skill/commands/drive.md` |
| App root | `showcase/src/App.tsx` |
