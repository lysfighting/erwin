# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-22
- **Last Session**: 2026-03-22 (arch-review v0.1.0 — creation + 5-round self-review + debug/revolution)

## Current Phase

**Phase**: Wing v0.1.0 — Anthropic official skill architecture rewrite

## What Was Just Completed

**Session AM** (2026-03-18) — Wing v0.1.0: full rewrite from teamwork

- **Naming**: "wing" chosen (Wings of Freedom / AoT connection, short, V4-extensible)
- **6 SKILL.md files**: wing, wing-issue, wing-claim, wing-drive, wing-ship, wing-rc (all ≤500 lines)
- **10 operations files**: init, doctor, auto, batch, fix, update, review, promote, hotfix, team-execution
- **VE single source**: `reference/ve-standard.md` (ends 6x duplication)
- **8 scripts migrated**: tw-*.sh → wing-*.sh (config fallback: .wing/ → .teamwork/ → .teamspace/)
- **skill-lint.sh**: Layer 1 static analysis (12 checks, 60 PASS across all 6 skills)
- **install.sh**: dual-mode (legacy→commands/, modern→skills/)
- **Validation**: config chain, git ops, contract hash all pass against `.wing/config.yml`
- **Metrics**: 14,231 → 5,904 lines (-59%), max SKILL.md 4,592 → 441 (-90%)

## Active Task / What's In Progress

**Wing v0.1.0 delivered. Transition phases remaining:**

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 1 | ✅ | Write 6 SKILL.md + 10 operations + reference |
| Phase 2 | ✅ | Migrate scripts (tw-*.sh → wing-*.sh) |
| Phase 3 | 🔄 | Switch erwin to wing (.wing/ config created, CLAUDE.md updated) |
| Phase 4 | ⬜ | Switch vi_agent to wing (--copy install) |
| Phase 5 | ⬜ | Delete skills/teamwork/ (after vi_agent stable on wing) |

## Next Session Should Start With

1. Read this file
2. **Wing remaining work**:
   - A. End-to-end test: create an Issue with `/wing-issue`, claim with `/wing-claim`, verify full flow
   - B. vi_agent migration (Phase 4): `./install.sh --copy --target ../vi_agent wing`
   - C. Skill-tester Layer 2 (smoke tests with mock data)
3. **Showcase tasks still pending**: DesignPatternLibrary, CaseStudies — blocking v1.0.0

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
- [x] Teamwork v2.0–v3.11.2 — 40+ versions, 6 skills, 10 scripts, 48 unit tests
- [x] Wing v0.1.0 — Anthropic skill architecture rewrite (6 SKILL.md, 10 operations, 8 scripts)
- [x] skill-lint.sh — Layer 1 static analysis (12 checks)
- [x] arch-review v0.1.0 — Multi-role debate review (8 mechanisms, 27 bugs fixed in 5 rounds)
- [ ] Wing Phase 4: vi_agent migration
- [ ] Wing Phase 5: teamwork deprecation
- [ ] Skill-tester Layer 2: smoke tests

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Wing scripts reference `${CLAUDE_SKILL_DIR}` which requires Anthropic skill infrastructure
- skill-lint.sh ### header check has false positives (structural headings vs output templates)

## File Locations

| What | Where |
|------|-------|
| Project goals | `_project/GOALS.md` |
| Idea garden | `_project/ideas/` |
| ADRs | `_project/decisions/` |
| V2 architecture | `docs/drive-showcase-v2-architecture.md` |
| Wing skills | `skills/wing/` |
| Wing architecture | `_project/ideas/skills/tw2-architecture.md` |
| Teamwork skills (legacy) | `skills/teamwork/` |
| Docs index | `docs/README.md` |
| App root | `showcase/src/App.tsx` |
