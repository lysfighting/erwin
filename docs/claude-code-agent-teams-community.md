# Claude Code Agent Teams 社区优秀项目汇总

> 基于 drive-skill 分析，深度搜索的社区中使用 Claude Code Agent Teams / 多 Agent 编排的优秀开源项目。
> 研究日期：2026-03-02

---

## 一、项目评级总览

| 项目 | ⭐ 推荐度 | 核心亮点 | 与 drive-skill 相似度 |
|------|----------|---------|---------------------|
| **Pied-Piper** | ⭐⭐⭐⭐⭐ | Playbook + BEADS + 多模型 Dream Team | 🔴 高（最相似） |
| **Claude Swarm** | ⭐⭐⭐⭐ | 原生 Agent Teams swarm 编排 | 🔴 高 |
| **wshobson/agents** | ⭐⭐⭐⭐ | 112 agents + 146 skills 插件生态 | 🟡 中 |
| **ruflo (claude-flow)** | ⭐⭐⭐⭐ | 企业级 swarm 平台 | 🟡 中 |
| **Claude Squad** | ⭐⭐⭐⭐ | 终端多 Agent 并行管理 | 🟡 中 |
| **Auto-Claude** | ⭐⭐⭐ | 全 SDLC 自主编码 + Kanban UI | 🟡 中 |
| **The Agentic Startup** | ⭐⭐⭐ | 生产级 agent collection + Output Styles | 🟢 低 |
| **TSK** | ⭐⭐⭐ | Rust + Docker 沙箱化 Agent | 🟢 低 |
| **Claude Task Master** | ⭐⭐⭐ | 任务管理 + Cursor AI 集成 | 🟢 低 |

---

## 二、Top 项目深度解析

### 2.1 🏆 Pied-Piper — 最值得研究的项目

**GitHub**: [sathish316/pied-piper](https://github.com/sathish316/pied-piper)

> **Pied-Piper 是目前与 drive-skill 最相似且互补性最强的项目。**

#### 核心特性

| 特性 | 详情 |
|------|------|
| **定位** | AI SubAgent 团队，自主/半自主执行长时间编码任务 |
| **任务管理** | 集成 [BEADS](https://github.com/steveyegge/beads)（Steve Yegge 的 Git-native 任务追踪） |
| **人工审批** | 全链路端到端追踪 + Human-in-the-loop 审批 |
| **运行环境** | Claude Code CLI / Docker / Cloud Desktop（AFK 也能跑） |
| **多模型支持** | 通过 [Claude-Code-Router](https://github.com/musistudio/claude-code-router) 支持多模型 |

#### Playbook 系统（核心创新）

Pied-Piper 的 **Playbook** 是预定义的可复用 SDLC 工作流：

| Playbook | 用途 | 状态 |
|----------|------|------|
| **Dream Team (多模型)** | GPT-5 Planner + Opus 4.5 Coder + Gemini 3.0 Reviewer + Haiku Orchestrator | ✅ 可用 |
| **Dream Team (Claude)** | Opus 4.5 Planner + Sonnet 4.5 Coder + Opus/Sonnet Reviewer + Haiku Orchestrator | ✅ 可用 |
| **Language Migration** | TypeScript → Python 全自动迁移 | ✅ 可用 |
| **Unit Test Coverage** | 自动提升单元测试覆盖率 | ✅ 可用 |
| **Microservices → Monolith** | 微服务整合 | 🔲 TODO |
| **Library Version Migration** | 库版本升级 | 🔲 TODO |
| **Tech Stack Migration** | 技术栈迁移 | 🔲 TODO |
| **Static Code Analysis** | 静态代码分析改进 | 🔲 TODO |

#### 与 drive-skill 的关键差异

| 维度 | drive-skill | Pied-Piper |
|------|-------------|------------|
| **任务管理** | 松散的 plan.md / feedback.md 文件 | **BEADS** 结构化任务追踪 |
| **工作流定义** | 112KB prompt 中硬编码 | **Playbook** 可配置、可复用 |
| **多模型** | 仅 Claude Opus | Dream Team 支持 GPT + Claude + Gemini |
| **运行方式** | Claude Code 会话内 | Claude Code + Docker + Cloud Desktop |
| **人工介入** | AskUserQuestion | Human-in-the-loop 审批门 |
| **持久化** | 会话内 TaskList | BEADS Git-native 跨会话 |

**🔑 最值得借鉴**：Playbook 系统 + BEADS 集成 + 多模型 Dream Team

---

### 2.2 Claude Swarm — 原生 Agent Teams Swarm

**GitHub**: [parruda/claude-swarm](https://github.com/parruda/claude-swarm)

| 特性 | 详情 |
|------|------|
| **定位** | 启动 Claude Code 会话并连接到 Agent Swarm |
| **架构** | 直接使用 Claude Code Agent Teams API |
| **通信** | 基于 Mailbox 的 Agent-to-Agent 通信 |
| **特色** | 语义路由，基于技能的 teammate 选择 |

**与 drive-skill 对比**：Claude Swarm 是 drive-skill 的**轻量替代** —— 没有 112KB 的方法论，但提供了更原生的 Agent Teams 编排。drive-skill 用 prompt engineering 重塑行为，Claude Swarm 用代码框架实现编排。

---

### 2.3 wshobson/agents — 最大的 Agent 插件生态

**GitHub**: [wshobson/agents](https://github.com/wshobson/agents)

| 指标 | 数量 |
|------|------|
| **Agent 角色** | 112 |
| **插件** | 72 |
| **Skills** | 146 |
| **工作流编排器** | 16 |

#### 架构亮点
- **Agent Teams Plugin**：7 种预定义团队配置（code review、debugging、feature dev 等）
- **Agent Skills**：遵循 Anthropic Agent Skills Specification，按需加载，token 高效
- **模块化**：所有能力打包为独立插件，mix-and-match

**与 drive-skill 对比**：wshobson/agents 是 **"广度优先"**（112 agent × 146 skills），drive-skill 是 **"深度优先"**（1 个 112KB 深度 prompt）。前者更像工具箱，后者更像方法论。

---

### 2.4 ruflo (claude-flow) — 企业级 Swarm 平台

**GitHub**: [ruvnet/ruflo](https://github.com/ruvnet/ruflo)（原 claude-flow）

| 特性 | 详情 |
|------|------|
| **定位** | "The leading agent orchestration platform for Claude" |
| **架构** | 层级化 / 网格化多 Agent swarm |
| **规模** | 1,456 releases，企业级架构 |
| **集成** | RAG、Claude Code / Codex 原生集成 |
| **通信** | Agent 间共享上下文、分而治之 |

**与 drive-skill 对比**：ruflo 是 **平台级产品**（Rust 实现、Docker 部署），drive-skill 是 **prompt 级产品**（Markdown 文件、symlink 安装）。两者解决同一问题但在不同抽象层。

---

### 2.5 Claude Squad — 终端级多 Agent 管理

**GitHub**: [smtg-ai/claude-squad](https://github.com/smtg-ai/claude-squad)

| 特性 | 详情 |
|------|------|
| **定位** | 终端 app，管理多个 Claude Code 实例 |
| **特色** | 每个 agent 在独立 workspace 工作 |
| **兼容** | Claude Code + Codex + Aider |
| **并行** | 同时处理多个任务 |

**适合场景**：不需要 agent 间通信的**独立并行任务**（如同时修多个 bug）。比 drive-skill 更简单但覆盖面更窄。

---

### 2.6 Auto-Claude — 全 SDLC 自主框架

**GitHub**: [AndyMik90/Auto-Claude](https://github.com/AndyMik90/Auto-Claude)

| 特性 | 详情 |
|------|------|
| **定位** | "Plans, builds, and validates software for you" |
| **UI** | Kanban-style 看板界面 |
| **特色** | 全 SDLC 覆盖 + 自验证 QA |
| **设计** | "Well-designed but not over-engineered" |

---

### 2.7 其他值得关注的项目

| 项目 | GitHub | 核心特色 |
|------|--------|---------|
| **TSK** | [dtormoen/tsk](https://github.com/dtormoen/tsk) | Rust CLI + Docker 沙箱化 agent，返回 git branch 供人审 |
| **The Agentic Startup** | [rsmdt/the-startup](https://github.com/rsmdt/the-startup) | 生产导向的 agent 集合，使用 Output Styles |
| **Claude Task Master** | [eyaltoledano/claude-task-master](https://github.com/eyaltoledano/claude-task-master) | AI 驱动的任务管理，与 Cursor AI 无缝集成 |
| **sudocode** | [sudocode-ai/sudocode](https://github.com/sudocode-ai/sudocode) | 轻量编排工具，存活在 repo 中，"It's giving Jira" |
| **Happy Coder** | [slopus/happy](https://github.com/slopus/happy) | 手机/桌面并行控制多 Claude，push 通知 |
| **Claude Task Runner** | [grahama1970/claude-task-runner](https://github.com/grahama1970/claude-task-runner) | 解决上下文长度限制，隔离任务执行 |

---

## 三、awesome-claude-code 生态全景

[hesreallyhim/awesome-claude-code](https://github.com/hesreallyhim/awesome-claude-code) 是 Claude Code 社区最完整的资源索引：

| 分类 | 包含 |
|------|------|
| **Agent Skills** | 通用 skills、语言特定、领域特定 |
| **Workflows & Knowledge Guides** | 工作流指南、Ralph Wiggum 循环 |
| **Orchestrators** | 10+ 个编排器（上文已列出） |
| **IDE Integrations** | VS Code 等集成 |
| **Usage Monitors** | 成本监控工具 |
| **Config Managers** | 配置管理器 |
| **Status Lines** | 状态栏插件 |
| **Hooks** | 生命周期钩子 |
| **Slash-Commands** | 版本控制、代码分析、测试、文档等 |
| **CLAUDE.md Files** | 按语言/领域的最佳实践模板 |
| **Alternative Clients** | 替代客户端 |

---

## 四、对 Erwin 的具体建议

### 最高优先级：研究 Pied-Piper

Pied-Piper 与我们的场景最匹配：
1. **BEADS 集成** —— 解决跨会话记忆问题
2. **Playbook 系统** —— 可复用的工作流定义（可以为 VPS 部署、快照构建等创建 Playbook）
3. **Dream Team 多模型** —— 利用不同模型的优势
4. **Human-in-the-loop** —— 适合我们的生产环境操作

### 值得尝试

| 工具 | 用途 |
|------|------|
| **Claude Squad** | 并行运行多个独立任务 |
| **TSK** | Docker 沙箱中运行高风险的 Agent 任务 |
| **Happy Coder** | 手机推送通知 + 远程控制（类似 drive-skill 的 Slack 集成） |

### 资源索引

始终关注 **awesome-claude-code** 仓库获取最新社区动态。
