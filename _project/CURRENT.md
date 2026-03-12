# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-12
- **Last Session**: 2026-03-12 (Session AK — v3.8.7–v3.9.2 + VE color research)

## Current Phase

**Phase**: Teamwork v3.x 完善 + docs 整理

## What Was Just Completed

**Session AK** (2026-03-12) — Teamwork v3.8.7 → v3.9.2:

- **v3.8.7**: docs completeness audit
- **v3.9.0**: `/team-rc fix` (cherry-pick to RC) + `/team-rc hotfix` (post-release production fix + backport to develop). Full docs/help/learn/onboarding sync.
- **v3.9.1**: Self-drive audit fixes — VE variable case, FIX APPLIED formatting, help section count, team-ship RC hint, manual version headers
- **v3.9.2**: Tips `shuf` crash on macOS — Read tool instead of Bash for tips
- **VE color research**: Investigated Shiki syntax highlighting as VE enhancement. Concluded: dead end — code block colors cannot mix with inline VE formatting. Research artifact: `docs/research/shiki-ve-color-lab.html`
- **v3.9.3 (in progress)**: tw-bump-version.sh path bug fix, version desync cleanup, VE accuracy fixes, _project/ docs update

## Active Task / What's In Progress

Self-drive audit: bug/VE/docs staleness sweep → v3.9.3

## Next Session Should Start With

1. Read this file
2. **Teamwork remaining needs**（按优先级）：
   - A. ~~Release 分支完整生命周期~~ ✅ Done in v3.9.0
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
- [x] Teamwork v3.7.1 — visual encoding standard + audit sweep + doctor 8→12 checks
- [x] Teamwork v3.8.1 — comprehensive audit: worktree chain + sync-all-checkboxes + 13 hints
- [x] Teamwork v3.8.2 — script沉淀 + DRY + UX dead-ends + 21 tests
- [x] Teamwork v3.8.3 — /team-rc cancel + upstream-first RC + post-merge PR checkbox
- [x] Teamwork v3.8.4 — dashboard formatting + worktree init fix + tab colors + bump validation
- [x] Teamwork v3.8.5 — /team sync + /team share + cross-MC knowledge intake
- [x] Teamwork v3.8.6 — documentation audit + version sync
- [x] Teamwork v3.8.7 — docs completeness audit
- [x] Teamwork v3.9.0 — /team-rc fix + /team-rc hotfix + full docs sync
- [x] Teamwork v3.9.1 — self-drive audit fixes (VE, formatting, counts)
- [x] Teamwork v3.9.2 — tips shuf crash on macOS
- [ ] Teamwork v3.9.3 — tw-bump-version.sh bug fix + VE accuracy + docs staleness

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
