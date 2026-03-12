# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-13
- **Last Session**: 2026-03-13 (Session AK — self-drive v3.9.3–v3.9.6)

## Current Phase

**Phase**: Teamwork v3.x hardening + quality sweep

## What Was Just Completed

**Session AK** (2026-03-12~13) — Teamwork v3.8.7 → v3.9.6:

- **v3.8.7**: docs completeness audit
- **v3.9.0**: `/team-rc fix` + `/team-rc hotfix` — complete RC lifecycle
- **v3.9.1**: Self-drive audit fixes (VE, formatting, counts)
- **v3.9.2**: Tips `shuf` crash on macOS
- **v3.9.3**: tw-bump-version.sh path bug, version desync, VE accuracy, docs staleness
- **v3.9.4**: Worktree config bug fix, dashboard RC section, tips 15→29, 48 unit tests
- **v3.9.5**: Script hardening (12 fixes: rebase ref, jq injection, git add safety, bash arrays, ship guard, hotfix re-entry). Docs staleness (skills/README, CLAUDE.md, idea pipeline). Ecosystem (install.sh robustness, erwin-init version).
- **v3.9.6**: Cross-platform sed portability, tw-pr.sh undefined var, cancel loop, stash inconsistency
- **VE color research**: Dead end — code block colors can't mix with inline VE formatting

## Active Task / What's In Progress

Self-drive audit complete. All known bugs fixed through v3.9.6.

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
- [x] Teamwork v2.0–v3.9.7 — 40+ versions, 6 skills, 10 scripts, 48 unit tests
- [ ] Deploy workflow (deploy.yml, staging gate)
- [ ] Issue version tracking (found-in labels, version labels)

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Teamwork: version_file (package.json bump) not yet in team-rc — deferred

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
