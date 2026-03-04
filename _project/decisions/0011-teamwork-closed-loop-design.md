# ADR-0011: Teamwork "需求自动领取交付" 闭环设计

**Date**: 2026-03-04
**Status**: Accepted (设计阶段，待实施)
**Session**: [2026-03-04-R](../sessions/2026-03-04-R.md)

## Context

Teamwork v2.5.1 有 6 个 skill 覆盖 issue → claim → drive → ship → release，但存在缝隙：
1. PR merge 后无自动化 cleanup（需手动 `/team-ship done`）
2. 无环境策略（dev/staging/prod）
3. Skill 和部署耦合
4. 无 changelog 自动生成

用户目标：做出完整闭环 + loop = 只需创建 milestone + issues，AI 自动迭代。

## Decision

### 闭环设计

**Per-issue 循环**:
```
/team-issue → /team-claim → /team-drive → /team-ship → [人 merge squash]
→ GitHub Action auto: close issue, label done, check milestone %
```

**Per-milestone**:
```
milestone 100% → Action 通知
→ workflow_dispatch deploy staging → AI QA
→ /team-release (检查 staging Action success) → git tag → gh release
→ workflow_dispatch deploy prod
```

### 核心决策

| 决策 | 选择 | 理由 |
|------|------|------|
| Merge 策略 | Squash merge | 一 commit 一 issue，可 cherry-pick |
| Tag 策略 | Release tag only（无 per-merge tag，无 rc tag） | squash commit 已含完整代码 |
| Branch 删除 | GitHub auto-delete head branches | repo 设置 |
| /team-ship done | 改为 GitHub Action（merge 触发） | 真正的自动化 |
| Deploy 触发 | workflow_dispatch（非 tag-triggered） | 参数化更安全 |
| Skill vs Deploy | 完全分离 | 风险等级不同不混合 |
| 分支模式 | 模式 A（tag-only，main 短暂冻结） | AI QA 周期短 |
| Dev 环境 | push to main → Action 自动部署，去 tag 化 | tag 不是部署日志 |
| Changelog | `--generate-notes` + append CHANGELOG.md | 零外部依赖 |
| Milestone 集成 | 轻量：/team-issue 加 --milestone | 避免复杂度爆炸 |
| Skill 数量 | 保持 6 个 | 扩展子命令而非新 skill |

### 设计原则

1. 风险等级不同的操作不能混在同一个触发路径上
2. Skill 管 git 操作，Action 管部署
3. Action run status 是门禁依据
4. 闭环优先，最佳实践其次
5. 不加新 skill，扩展现有子命令
6. Tag 是版本语义标记，不是部署日志

## Consequences

**Good**:
- 真正的闭环：merge 后全自动到 issue closed
- 部署和 skill 解耦，风险隔离
- 6 个 skill 不增加，复杂度可控

**Bad**:
- GitHub Action 比 skill 更难调试
- workflow_dispatch 需要手动触发（staging/prod），不完全自动
- 模式 A 在 QA 期间 freeze main（短暂阻塞新开发）

## Related

- ADR-0010: Teamwork v2 GitHub-first（前序决策）
- Session R: 全量讨论过程
