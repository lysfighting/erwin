# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-10
- **Last Session**: [2026-03-10-AF](sessions/2026-03-10-AF.md) (Session AF — VE 系统性治理 #001)

## Current Phase

**Phase**: Teamwork v3.x 完善 + docs 整理

## What Was Just Completed

**Session AF** — VE 系统性治理 (GitHub #1):

- Fixed ~15 remaining VE violations across 4 skill files (team.md, team-ship.md, team-rc.md, team-drive.md)
- team-claim.md and team-issue.md verified clean (0 violations)
- Key fixes: ERROR bold prefix, doctor hints on all STOPs, raw git→/team-* skills, Learn code blocks→VE markdown
- 14 files changed, 344 insertions, 82 deletions (includes session AE uncommitted work)
- GitHub issue #1 closed

## Active Task / What's In Progress

无活跃任务。VE 治理完成，全部 6 skill 文件 0 违规。

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
- [x] Teamwork v3.7.1 — VE 系统性治理: 全部 skill 输出 100% VE 合规 (#1)

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
