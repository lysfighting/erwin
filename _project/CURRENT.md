# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-10
- **Last Session**: [2026-03-10-AE](sessions/2026-03-10-AE.md) (Session AE — Visual Encoding Standard + Audit + Self-Drive)

## Current Phase

**Phase**: Teamwork v3.x 完善 + docs 整理

## What Was Just Completed

**Session AE** (2 conversation windows) — Visual Encoding Standard + comprehensive audit:

- Visual unification mission: all 6 skill files → Visual Encoding Standard (Drive Swarm, 5 parallel agents)
- Created `docs/visual-encoding-standard.md` — single source of truth for color design language
- 5-scout audit: found 11 real issues (2 visual, 3 raw cmd, 2 tips, 3 next-step, 1 stale doc)
- ONBOARDING.md modernized: code block examples → formatted markdown + Advanced Features section
- Visual Encoding reference embedded in all 5 sub-skill file headers
- Doctor expanded from 8 → 11 checks (+ auth, freshness, tests, CI)
- Teamwork manual updated: doctor + visual encoding sections

## Active Task / What's In Progress

无活跃任务。

## Next Session Should Start With

1. Read this file
2. **Teamwork 18 项需求实现**（按优先级）：
   - A. Release 分支完整生命周期（prep, hotfix routing, backport, cleanup）
   - B. Deploy workflow（deploy.yml, staging gate）
   - C. Issue 版本追踪（found-in labels, version labels）
   - 详见 `_project/sessions/2026-03-05-Z.md` 需求汇总
3. **Showcase 任务仍待完成**：DesignPatternLibrary, CaseStudies — blocking v1.0.0

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
- [x] /session-record v1.0.0 — full verbatim session recording skill
- [x] Teamwork v2.0–v2.7.3 — GitHub-first skills, script extraction, E2E validation
- [x] Teamwork v3.1.0 — untracked issue detection + UNTRACKED dashboard + learn audit
- [x] Teamwork v3.2.0 — unified emoji visual encoding + milestone prefix matching backport
- [x] Teamwork v3.2.1 — script extraction + deep bug fix + stale cleanup
- [x] Teamwork v3.3.0 — /team wrap + /team auto + 16 bugs fixed (3 CRITICAL + 8 HIGH)
- [x] Teamwork v3.4.0 — doctor actionability: fix mode, help, vs Remote column, Check 8
- [x] Teamwork v3.4.1 — worktree local mode + dashboard map + doctor idle detection
- [x] Teamwork v3.4.2 — worktree dashboard mandatory render
- [x] Teamwork v3.4.3 — merge mode logic gaps + worktree audit
- [x] Teamwork v3.5.0 — /team config + dashboard version + docs completeness
- [x] Teamwork v3.5.1 — design review: logic gaps + docs alignment
- [x] Teamwork v3.7.1 — visual encoding standard + audit sweep + doctor 8→11 checks

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Teamwork: version_file (package.json bump) not yet in team-rc — deferred
- vi_agent: working directory on `system_v5` branch, skill updates on `main`

## File Locations

| What | Where |
|------|-------|
| Project goals | `_project/GOALS.md` |
| Idea garden | `_project/ideas/` |
| ADRs | `_project/decisions/` |
| V2 architecture | `docs/drive-showcase-v2-architecture.md` |
| Teamwork manual | `docs/teamwork-ai-manual.md` |
| Docs index | `docs/README.md` |
| Teamwork skills | `skills/teamwork/` |
| Teamwork scripts | `skills/teamwork/scripts/` |
| App root | `showcase/src/App.tsx` |
