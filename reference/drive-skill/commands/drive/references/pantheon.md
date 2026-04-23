# Pantheon — 轻量思想家索引

核心机制：**invoke 人名来强制认知 reframing，不预存详细资料，由 agent 当场 WebSearch 研究。**

节省 token 的关键：只存名字 + 领域标签。模型本身已有大量知识，WebSearch 补充具体著作和方法论。

---

## 如何使用

### 1. 选人（Leader 根据问题本质选 2-3 人）

浏览下方索引，选择与问题领域最相关的思想家。不限于此列表——任何你认为相关的伟大思想家都可以 invoke。

### 2. Spawn Thinker Council

```
Agent(
  name: "thinker-council",
  model: "opus",
  prompt: "You are a council of great thinkers analyzing a problem.

    THINKERS: {name1}, {name2}, {name3}
    PROBLEM: {the bone — one sentence}
    CONTEXT: {constraints, domain, stakes}

    For EACH thinker:
    1. WebSearch '{name} key ideas methodology philosophy'
    2. WebSearch '{name} most influential work core contributions'
    3. Based on their thinking framework, how would they analyze THIS problem?
    4. What would they say everyone else is overlooking?
    5. What's the most important question they'd ask?

    SYNTHESIZE:
    - Agreements (high confidence directions)
    - Disagreements (tensions to explore)
    - Reframing (how the problem looks different through their combined lens)
    - Question for the user (something we wouldn't have thought to ask)"
)
```

### 3. 何时再次咨询

- 设计 harness 维度前（我们在测量对的东西吗？）
- 遇到设计僵局时（需要新的思考角度）
- Phase V 发现系统性问题时（根因分析需要更深框架）

---

## 思想家索引

### 系统与架构
| 名字 | 领域标签 |
|------|----------|
| John von Neumann | 计算架构, 博弈论, 系统设计 |
| Claude Shannon | 信息论, 信号/噪声, 编码 |
| Norbert Wiener | 控制论, 反馈系统 |
| Herbert Simon | 有限理性, 复杂系统, 人工智能 |
| Donella Meadows | 系统思维, 杠杆点, 反馈循环 |

### 软件工程
| 名字 | 领域标签 |
|------|----------|
| Fred Brooks | 大型软件项目, 复杂性管理 |
| Edsger Dijkstra | 结构化编程, 正确性证明, 简洁 |
| Alan Kay | 面向对象, 用户界面, 系统设计 |
| Donald Knuth | 算法, 性能, 严谨性 |
| Leslie Lamport | 分布式系统, 时序, 形式化 |
| Barbara Liskov | 抽象, 类型系统, 模块化 |

### 设计与创造
| 名字 | 领域标签 |
|------|----------|
| Dieter Rams | 功能主义设计, "少但更好" |
| Christopher Alexander | 模式语言, 建筑/软件设计 |
| Edward Tufte | 信息可视化, 数据呈现 |
| Don Norman | 人机交互, 可用性, 日常设计 |

### 科学方法与认知
| 名字 | 领域标签 |
|------|----------|
| Richard Feynman | 第一性原理, 物理直觉, 教学 |
| Karl Popper | 证伪主义, 科学方法 |
| Thomas Kuhn | 范式转换, 科学革命 |
| Daniel Kahneman | 认知偏差, 决策, 系统1/系统2 |
| Charlie Munger | 多元思维模型, 反向思考 |

### 战略与组织
| 名字 | 领域标签 |
|------|----------|
| Sun Tzu | 战略, 资源配置, 环境利用 |
| Peter Drucker | 管理, 知识工作者, 有效性 |
| Eliyahu Goldratt | 约束理论, 瓶颈, TOC |
| W. Edwards Deming | 质量管理, 系统优化, PDCA |
| John Boyd | OODA循环, 快速迭代, 机动 |

### 哲学与根本思考
| 名字 | 领域标签 |
|------|----------|
| Aristotle | 逻辑, 分类, 目的论 |
| Immanuel Kant | 范畴, 先验框架, 义务论 |
| Ludwig Wittgenstein | 语言哲学, 意义边界 |
| 老子 | 道, 无为, 自然系统 |
| Alfred North Whitehead | 过程哲学, 抽象的危险 |

---

## 选人启发

| 问题类型 | 推荐组合 |
|----------|----------|
| 架构设计 | von Neumann + Brooks + Alexander |
| 性能优化 | Knuth + Shannon + Goldratt |
| 用户体验 | Norman + Kay + Tufte |
| 质量保证 | Dijkstra + Deming + Popper |
| 战略决策 | Sun Tzu + Munger + Boyd |
| 复杂系统 | Meadows + Simon + Wiener |
| 根因分析 | Feynman + Kahneman + Goldratt |
| 简化设计 | Rams + Dijkstra + 老子 |
