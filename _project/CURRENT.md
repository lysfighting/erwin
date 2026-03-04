# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-04
- **Last Session**: [2026-03-04-U](sessions/2026-03-04-U.md) (Session U — Teamwork v2.6.1 Falsification Audit)

## Current Phase

**Phase**: Teamwork v2.6.1 完成 — 全文件审计 + 21 bug 修复

## What Was Just Completed

**Session U** — Falsification Audit (Popper's method):

- Applied Popper/Descartes/Bacon to critical self-review of all teamwork skill files
- 5 parallel audit agents reviewed 7 files → found ~80 bugs
- Triaged and fixed 21 verified bugs (6 CRITICAL + 15 HIGH)
- CRITICAL: freshness check, PyYAML dep, branch protection API, CHANGELOG injection, stale merge, wrong push branch
- HIGH: init force, hardcoded labels (5 files), sed syntax, milestone quoting, undefined vars, missing --state
- All skill files bumped to v2.6.1

## Active Task / What's In Progress

None — v2.6.1 audit complete.

## Next Session Should Start With

1. Read this file
2. **vi_agent 实践验证** — 在 vi_agent 上运行 `/team init` 验证 v2.6.1 fixes
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
