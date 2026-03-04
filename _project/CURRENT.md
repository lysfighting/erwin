# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-05
- **Last Session**: [2026-03-05-Y](sessions/2026-03-05-Y.md) (Session Y — Design Coding Skill 讨论 + /session-record)

## Current Phase

**Phase**: Design Coding Skill 多轮设计讨论中

## What Was Just Completed

**Session Y** — Design Coding Skill 讨论 + /session-record skill:

- 8 轮深度讨论，确立 design coding skill 核心理念
- **灵魂确立**: Taste > Rules > Trends（品味 > 规则 > 趋势）
- 结构方案 v3（9 模块）：Manifesto, Output Standard, Reference Protocol, Visual Feedback Loop, Design Language, Pantheon, Quality Protocol, Anti-patterns, Execution
- 参考范围扩展至 5 域：Nature / Art & Architecture / Mathematics / Digital / Philosophy
- 艺术哲学体系：黄金比例、侘寂、間、Sprezzatura、格式塔、日神酒神
- 新 skill: `/session-record` v1.0.0（完整对话记录模式）
- 全程 session 记录：`_project/sessions/2026-03-05-Y.md`（580 行完整对话）

## Active Task / What's In Progress

Design Coding Skill — 讨论阶段（未进入实现）

## Next Session Should Start With

1. Read this file
2. **继续 Design Coding Skill 讨论**（按优先级排序）：
   - Manifesto 具体信条（taste-first 已定为第一条，还需哪几条？）
   - 独立 skill (A) vs 可与 /drive 组合的品味层 (B)？
   - 众神殿人选——需分领域 web research
   - Browser MCP 方案选型（已记录，不急）
   - 9 模块是否太多？合并讨论
3. **闭环设计讨论仍待继续**：自动化测试、IaC 归属、deploy Action、vi-agent docs
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
