# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-07
- **Last Session**: [2026-03-07-AC](sessions/2026-03-07-AC.md) (Session AC — 视觉统一 + milestone 回流 + 闭环讨论)

## Current Phase

**Phase**: vi_agent 闭环完善 + Design Coding Skill 讨论并行

## What Was Just Completed

**Session AC** — Teamwork v3.2.0 视觉统一 + milestone 回流 + 闭环讨论:

- 6 个 skill 文件统一 emoji visual encoding（priority dots, section icons, CI icons, progress bar）
- Dashboard 性能优化（7-8 串行 API → 4-5 并行）
- vi_agent milestone 前缀匹配 hotfix 回流 erwin（4 处，含 1 处两边都漏的 batch gap）
- 闭环讨论：`/team-wrap`（反向闭环）三条路径确定 + 正向闭环延后
- Session log: `_project/sessions/2026-03-07-AC.md`

## Active Task / What's In Progress

1. **vi_agent 分支问题**：工作目录在 `system_v5`，skill 更新在 `main`（已 push）
2. Design Coding Skill — 讨论阶段（未进入实现）

## Next Session Should Start With

1. Read this file
2. **vi_agent 分支处理**：`system_v5` → merge main 或 checkout main
3. **Teamwork 18 项需求实现**（按优先级）：
   - A. Release 分支完整生命周期（prep, hotfix routing, backport, cleanup）
   - B. Deploy workflow（deploy.yml, staging gate）
   - C. Issue 版本追踪（found-in labels, version labels）
   - D. Skill 适配（team-ship/team-claim 感知 release 分支）
   - E. Config schema（release_branch 从静态值改为模式）
   - 详见 `_project/sessions/2026-03-05-Z.md` 需求汇总
4. **Teamwork 可视化增强**（Session AB 结论）：
   - `/team init` 集成 GitHub Projects 看板（`gh project create`）
   - `/team` board ASCII 汇总增强
   - Review 路径加速（/team-ship 内嵌 diff feedback）
   - 详见 `_project/sessions/2026-03-05-AB.md`
5. **继续 Design Coding Skill 讨论**
6. **Showcase 任务仍待完成**：DesignPatternLibrary, CaseStudies — blocking v1.0.0

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
- [x] Teamwork v2.0 — 4 GitHub-first skills
- [x] Teamwork v2.1 — 6 vi-agent features + 13 bug fixes
- [x] Teamwork v2.3 — status fix + version management + format standardization + scripts
- [x] Teamwork v2.4 — stability hardening + bidirectional sync + tw-config.sh
- [x] Teamwork v2.4.1 — E2E validated on vi_agent + mc_label dual-schema + inline comment fix
- [x] Teamwork v2.5.0 — team-release skill + init idempotency fix + doc sweep
- [x] Teamwork v2.5.1 — team-issue duplicate check (gh search + LLM semantic match)
- [x] Teamwork v2.6.0 — dual-branch model + release safety gates + post-merge Action + config-driven refs
- [x] Teamwork v2.6.1 — falsification audit: 6 critical + 15 high bug fixes across all files
- [x] Teamwork v2.7.0 — script extraction: 4 scripts (32 subcommands) + 6 skill files refactored
- [x] Teamwork v2.7.1 — script testing audit: 21 bugs found, 14+4 fixed across 4 scripts
- [x] Teamwork v2.7.2 — vi_agent E2E validation: 3 bugs found, 3 fixed
- [x] Teamwork v2.7.3 — init audit: dynamic branch protection + post-merge workflow + vi_agent deploy
- [x] Teamwork v3.2.0 — unified emoji visual encoding + milestone prefix matching backport

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Teamwork: version_file (package.json bump) not yet in team-release — deferred to future version
- Teamwork: `/team wrap` + `/team auto` 已实现 v3.3.0，经 3-agent 并行审计修复 16 bugs（3 CRITICAL + 8 HIGH）
- vi_agent: working directory on `system_v5` branch, skill updates on `main`

## File Locations

| What | Where |
|------|-------|
| Project goals | `_project/GOALS.md` |
| Idea garden | `_project/ideas/` |
| ADRs | `_project/decisions/` |
| V2 architecture | `docs/drive-showcase-v2-architecture.md` |
| Teamwork architecture | `docs/teamwork-ai-manual.md` |
| Teamwork skills | `skills/teamwork/` |
| Teamwork scripts | `skills/teamwork/scripts/` |
| App root | `showcase/src/App.tsx` |
