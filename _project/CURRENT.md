# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-04
- **Last Session**: [2026-03-04-W](sessions/2026-03-04-W.md) (Session W — Teamwork v2.7.1 Script Testing Audit)

## Current Phase

**Phase**: Teamwork v2.7.1 完成 — 脚本测试审计

## What Was Just Completed

**Session W** — Script Testing Audit (21 bugs found, 14 fixed):

- 4 parallel test agents + 1 cross-reference auditor tested all 32 subcommands
- 2 CRITICAL: tag rev-parse matched branches, merge-to on wrong branch
- 6 HIGH: silent ensure-base, sed injection, fetch kills script, grep metachar, commit-type false positives, toggle-task false success
- 6 MEDIUM: grep anchor, stderr, detached HEAD, label transition, PR base param, body validation
- Cross-reference audit: .md-to-script interfaces clean, 5 orphan subcommands (not bugs)

## Active Task / What's In Progress

None — v2.7.1 testing audit complete.

## Next Session Should Start With

1. Read this file
2. **vi_agent 实践验证** — run `./install.sh teamwork` then test on vi_agent
3. **继续闭环设计讨论**（按优先级排序）：
   - 自动化测试流程：AI QA 具体实现、工具选型、覆盖范围
   - IaC 归属：项目内 vs gcp-compute-infra，权限模型
   - deploy Action 具体实现（依赖上面两项）
   - vi-agent docs → milestone/issues（操作型，讨论量少）
4. **Showcase 任务仍待完成**：DesignPatternLibrary, CaseStudies — blocking v1.0.0

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
- [x] Teamwork v2.6.0 — dual-branch model + release safety gates + post-merge Action + config-driven refs
- [x] Teamwork v2.6.1 — falsification audit: 6 critical + 15 high bug fixes across all files
- [x] Teamwork v2.7.0 — script extraction: 4 scripts (32 subcommands) + 6 skill files refactored
- [x] Teamwork v2.7.1 — script testing audit: 21 bugs found, 14 fixed across 4 scripts

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Teamwork: version_file (package.json bump) not yet in team-release — deferred to future version
- Teamwork: /team-drive still unimplemented (thin wrapper over /drive)

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
