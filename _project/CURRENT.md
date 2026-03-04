# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-04
- **Last Session**: [2026-03-04-R](sessions/2026-03-04-R.md) (Session R — Teamwork 闭环设计讨论)

## Current Phase

**Phase**: Teamwork v2.6 设计 — "需求自动领取交付"闭环 + 环境策略

## What Was Just Completed

**Session R** — Teamwork 闭环设计讨论（纯设计，无代码）:

全量讨论记录见 `sessions/2026-03-04-R.md`。达成 19 项设计共识，提炼 6 条设计原则。
核心决策：squash merge + release tag only + Skill/Deploy 职责分离 + workflow_dispatch 部署 + 模式 A 环境策略。

## Active Task / What's In Progress

**设计讨论进行中**，以下话题已完成，以下待续：

已完成：合并机制、tag 策略、changelog、milestone 集成、环境策略、部署触发、风险分离
待讨论：IaC 归属、自动循环上下文管理、vi-agent docs→issues、自动化测试

## Next Session Should Start With

1. Read this file + `sessions/2026-03-04-R.md`（全量讨论记录，19 项共识）
2. **继续闭环设计讨论**（未完成话题）：
   - IaC 归属：项目内 vs gcp-compute-infra，权限模型
   - 自动循环上下文管理：leader context 积累、中断恢复
   - vi-agent docs (product.md + system_v4.md) → 批量创建 milestone/issues
   - 自动化测试流程（AI QA 具体实现）
3. **设计讨论完成后 → 实施优先级**：
   - `/team-ship done` → GitHub Action 化
   - `/team-issue` 加 `--milestone` 参数
   - `/team-release` 加 CHANGELOG.md 生成 + staging Action 检查
   - `/team-drive` 加 "不在 main" 安全检查
   - deploy.yml (workflow_dispatch) 骨架
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
- [ ] Teamwork v2.6 — 闭环设计（Session R 设计中）

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Teamwork: version_file (package.json bump) not yet in team-release — deferred to v2.6
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
