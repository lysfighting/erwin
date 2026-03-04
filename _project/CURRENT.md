# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-04
- **Last Session**: [2026-03-04-Q](sessions/2026-03-04-Q.md) (Session Q — Teamwork v2.5.0)

## Current Phase

**Phase**: Teamwork v2.5.0 — team-release + init idempotency fix

## What Was Just Completed

**Session Q** — Teamwork v2.5.0 (full day of teamwork iteration):

1. **team.md v2.4.1 upgrade** — replaced all sed chains in dashboard with tw-config.sh calls (mc_label dual-schema, STATUS_PREFIX, PRIORITY_PREFIX, CURRENT_VERSION). skill_version template updated.
2. **/team-release skill** — new `team-release.md` v2.5.0. Full version lifecycle: milestone check → AI audit (conditional on spec_path) → user confirm → git tag → gh release create --generate-notes → close milestone → next version setup.
3. **/team init idempotency fix** — Step 4b now has MERGE MODE (config exists → preserve all sections, update skill_version, append missing) vs FRESH INSTALL (no config → full write). Eliminates destructive overwrite of custom configs.
4. **Documentation sweep** — team.md help added /team-issue + /team-release + VERSION LIFECYCLE section. ONBOARDING.md: skill_version 2.1.0→2.5.0, Step 2 rewritten for /team-issue, Step 6 added for /team-release, daily workflow + versions comment updated.

## Active Task / What's In Progress

None. Session Q complete.

## Next Session Should Start With

1. Read this file
2. Continue V2 Phase 3 showcase tasks (DesignPatternLibrary, CaseStudies) — still blocking v1.0.0
3. **Future teamwork missions**:
   - `/team-drive` wrapping full `/drive` (thin wrapper, skip Phase 0) — still unimplemented
   - version_file support in team-release (package.json version bump) — deferred to v2.6

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
- [x] Teamwork v2.5.0 — team-release skill + init idempotency fix + doc sweep
- [x] Teamwork v2.5.1 — team-issue duplicate check (gh search + LLM semantic match)

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Teamwork: initial milestone must be created manually (non-fatal warning in /team-issue)
- Teamwork: version_file (package.json bump) not yet in team-release — deferred to v2.6

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
