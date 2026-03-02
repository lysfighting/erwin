# ADR-0008: /new-project → /erwin-init — 从脚手架到认知基础设施

- **Date**: 2026-03-02
- **Status**: Accepted
- **Supersedes**: ADR-0006 (new-project skill design)

---

## Context

`/new-project` v1 存在身份混乱。它问了 4 个问题（项目名、描述、技术栈、AI 语言），行为像一个项目脚手架工具。但它实际创建的 `_project/`、CLAUDE.md、session protocol、ADR 系统、auto-memory 全是 **stack-agnostic 的认知基础设施**。技术栈问题暴露了这个矛盾——整个 intelligence layer 跟你用 React 还是 Python 毫无关系。

更根本的问题：erwin 的定位是 **vibe coding 框架**——为 AI 辅助开发提供完整的项目版本管理和认知自我进化机制。`/new-project` 是这个框架的安装入口，不是代码生成器。

---

## Decision

### 1. 重命名：`/new-project` → `/erwin-init`

**理由**：
- `new-project` 暗示「创建新项目」（脚手架语义），但 skill 也适用于已有项目
- `erwin-init` 明确表达「初始化 erwin 框架」，类比 `git init`
- 框架核心 skill 统一用 `erwin-` 前缀；独立工具 skill 不加前缀（如 `/visual-test`）

### 2. 零问题入口

**去掉全部 4 个初始化问题**：
- 项目名 → `basename $(pwd)` 自动获取
- 描述 → GOALS.md 对话中自然产生
- 技术栈 → Claude 在后续 session 中工作时自然发现，写入 MEMORY.md → 更新 CLAUDE.md
- AI 语言 → Phase 2 对话中推断用户语言偏好，写入 MEMORY.md

### 3. 两阶段执行

**Phase 1：建设施（自动，零问题，< 30 秒）**
- 创建 `_project/` 全套结构
- 创建 `CLAUDE.md`（Tech Stack: TBD）
- 创建 `.claude/rules/`（ai-workflow.md, git-conventions.md）
- 创建 `MEMORY.md`（local）
- 创建 `CHANGELOG.md`
- `git init` + first commit

**Phase 2：项目理解对话（提示但可跳过）**
- 问一句：「要现在填写 GOALS.md 吗？」
- YES → 引导对话：愿景（做什么、为谁）、成功标准、anti-goals → 结构化写入 GOALS.md
- NO → 报告完成，Next steps 里提醒
- 硬边界：对话只涉及项目理解，绝不涉及技术方案或代码实现

### 4. 命名空间规则

| 类型 | 前缀 | 示例 |
|------|------|------|
| 框架核心 skill | `erwin-` | `/erwin-init` |
| 独立工具 skill | 无 | `/visual-test`、`/drive` |

判断标准：这个 skill 是否「构成 erwin 框架本身」？是 → 加前缀。否 → 不加。

---

## Consequences

**Positive**:
- 身份清晰：认知基础设施安装器，不是脚手架
- 零摩擦入口：不问任何问题，直接创建
- 适用于已有项目：名字和行为都不暗示「新建」
- GOALS.md 引导对话帮用户结构化思考（但不强制）
- 技术栈信息由 Claude 自然积累，更准确

**Negative**:
- 重命名需要更新 install.sh、skills/ 目录结构、文档引用
- `erwin-` 前缀增加输入长度（但频率低，可接受）
- Phase 2 的对话引导质量依赖 skill prompt 的设计

## Alternatives Rejected

- **保留 `/new-project` 名字**：语义不准确，暗示脚手架
- **`/init`（无前缀）**：太通用，无法区分 erwin 框架 skill 和独立 skill
- **全部 skill 加 `erwin-` 前缀**：`/erwin-visual-test` 不合理，visual-test 跟框架无关
- **保留技术栈问题**：intelligence layer 是 stack-agnostic 的，初始化时问无意义
