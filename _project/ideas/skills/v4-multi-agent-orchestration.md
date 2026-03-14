# V4: Multi-Agent Orchestration — 中控调度 + MC 间通信

> Status: **Design** (from Session 2026-03-13-AM, Phases 4-6)
> Prerequisite: v3.10 信息同步层 (已交付 v3.10.4)
> Difficulty: **Very High** — 前沿形态，当前 Claude Code 能力边界

---

## 1. Problem

当前 teamwork 是**单任务串行模型**：一次只能 claim→drive→ship 一个 Issue。

现实中一个 milestone 有 5-10 个 Issues，每个需要独立的 worktree + context。人工逐个执行效率极低，但如果并行执行，MC 之间完全信息孤立：
- 冲突修改（两个 MC 改同一个文件）
- 重复发现（两个 MC 各自 debug 同一个 root cause）
- 接口变更不感知（MC-A 改了函数签名，MC-B 还在用旧的）

---

## 2. Vision — 三个递进层次

| 层次 | 名称 | 性质 | 复杂度 |
|------|------|------|--------|
| L1 | Ship 时写 comments 通知 | 具体功能 | Low |
| L2 | MC 间通信协议 | 架构问题 | Medium |
| L3 | 中控调度器 (Scheduler) | 范式转变 | **Very High** |

L1 和 L2 的基础已在 v3.10 (信息同步层) 交付。V4 的核心是 **L3: 中控调度器**。

---

## 3. 通信协议 (Two-Layer)

### Layer 1 — 静态（claim 时注入）

Contract 中包含 `## Related MCs` section：
- 列出所有 sibling MC + objectives + overlapping files
- 由 `/team auto-all` coordinator 在 batch claim 时注入
- 每个 MC 从开始就知道"谁在旁边干什么"

### Layer 2 — 动态（执行中推送）

基于 v3.10 的 comment read/write 机制：

- `/team share` 主动推送发现（已实现，v3.8.5）
- drive Step 0b 读取 Cross-MC comments（已实现，v3.10）
- **V4 新增**：ship 完成后自动 post 标准化 completion signal：

```markdown
🔗 **Cross-MC** from #042 — SHIPPED (PR #150)
Modified: src/auth.py (lines 42-78), tests/test_auth.py
Interface change: authenticate(user, token) → authenticate(user, token, scope)
```

Anti-plague 规则不变：
1. Human gateway（AskUserQuestion before post）
2. No transitive propagation
3. Source attribution + confidence

---

## 4. 调度策略 (Scheduler Algorithm)

```
1. FETCH all assigned open Issues with mission label
2. For each: scan Context Files → overlapping files
3. BUILD overlap matrix (Issue_i × Issue_j → shared files)
4. BUILD dependency graph:
   - Explicit: "depends on #N" in Issue body
   - Implicit: write-before-read file overlap
5. TOPOLOGICAL SORT by dependency, then priority
6. APPLY parallelism rules:
   - P0 (critical): 永不并行，串行执行
   - P1 (high):     最多 2 并行（零文件重叠时）
   - P2/P3 (normal): 最多 3 并行（零文件重叠时）
   - Global cap = 3 concurrent MCs
7. OUTPUT JSON execution plan:
   [{"batch":1, "issues":[42], "mode":"sequential"},
    {"batch":2, "issues":[45,47], "mode":"parallel"}]
```

### 四种执行模式对比

| | A: Sequential | B: Parallel claim + seq drive | C: Parallel selective | D: Multi-terminal |
|---|---|---|---|---|
| 实现 | 简单 | 中等 | 复杂 | 最复杂 |
| 通信 | 不需要 | claim 时 | claim + drive | 实时 |
| 冲突 | 零 | 零 | 低（独立 MC） | 需要协调 |
| 质量 | 最高 | 高 | 高 | 最高 |
| 完成时间 | sum(MC) | ~sum(MC) | max(MC) | max(MC) |

**推荐路径**：
- MVP: **Option C** — parallel claim + sequential drive + optional parallel for independent MCs
- 终极: **Option D** — 多终端多 session，coordinator 只做 setup

---

## 5. v3.10 已交付的基础设施

| # | 能力 | v3.10 状态 | V4 用途 |
|---|------|-----------|---------|
| 1 | Comment 读写 | ✅ 3 个 script 子命令 | 跨 MC 通信 |
| 2 | Cross-MC comment format | ✅ `🔗 **Cross-MC**` 标准化 | completion signal |
| 3 | `/team share` | ✅ human gateway + anti-plague | 主动推送发现 |
| 4 | Drive Step 0b | ✅ 读取 Cross-MC comments | 感知 sibling 状态 |

### V4 还需要的基础设施（未实现）

| # | 能力 | 说明 |
|---|------|------|
| 1 | File overlap detection | `tw-scheduler.sh overlap-check` — 检测 MC 间文件冲突 |
| 2 | Batch claim + worktree | `/team auto-all` 一口气 claim 多个 Issue |
| 3 | Dependency graph builder | 从 Issue body + Context Files 构建依赖关系 |
| 4 | Scheduler script | `tw-scheduler.sh plan` → JSON 执行计划 |
| 5 | Completion signal 自动化 | ship 后自动 post Cross-MC shipped comment |
| 6 | Coordinator orchestrator | Team Lead 角色：分配 → 监控 → 汇总 |

---

## 6. 为什么这很难 — 前沿挑战

### 6.1 Claude Code 当前限制

- **context:fork** 不支持跨 agent 通信 — fork 出去的 agent 是信息孤岛
- **SendMessage** (Teams) 是文件邮箱模式，不是实时通信
- 没有共享内存 — 所有通信必须通过 GitHub comments 或文件系统
- 单个 session 的 context window 有限 — coordinator 无法同时 hold 5 个 MC 的完整上下文

### 6.2 根本难点

1. **冲突检测 timing**: 静态分析（claim 时）vs 动态检测（drive 中）。静态太粗，动态太贵。
2. **失败恢复**: MC-A ship 后 MC-B 因为 MC-A 的修改而冲突。谁 rebase？谁通知？
3. **优先级动态变化**: MC-A 的发现可能改变 MC-B 的优先级。协调成本 > 并行收益？
4. **测试**: 如何测试多 agent 并行行为？当前连单 agent 的行为测试都还没做（→ skill-tester）。

### 6.3 与 tw2 的关系

V4 是**在 tw2 架构之上**的能力层：
- tw2 解决的是**单 skill 的 prompt 效率**（4,594 → ≤500 lines）
- V4 解决的是**多 MC 的协调效率**（串行 → 并行）
- tw2 是 V4 的前提：如果单个 skill 都因为 prompt 过大而不可靠，多 agent 协调更不可能可靠

---

## 7. 实施路径

```
tw2 (架构重构) ──→ skill-tester (质量保证) ──→ V4 (多 agent 协调)
     ↑                    ↑                          ↑
  当前阶段              当前阶段                   未来目标
```

1. **先完成 tw2** — 把 teamwork 迁移到 Anthropic 官方 skill 架构
2. **再建 skill-tester** — 确保单 skill 行为可测试、VE 可验证
3. **最后做 V4** — 在可靠的单 skill 基础上构建多 agent 协调

---

## 8. Source

- Session: `_project/sessions/2026-03-13-AM.md` (Phases 4-6)
- Cross-MC idea: `_project/ideas/skills/cross-mc-knowledge-sharing.md` (v3.8.5 已实现)
- Session: `_project/sessions/2026-03-04-R.md` (跨 session 通信设计)
- Session: `_project/sessions/2026-03-12-AJ.md` (Cross-MC knowledge sharing)
- v3.10 release: `skills/teamwork/CHANGELOG.md`
