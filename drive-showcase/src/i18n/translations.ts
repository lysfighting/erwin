import type { Lang } from '../contexts/LangContext'

const translations = {
  // ── NavBar ──
  'nav.home':         { en: 'Home',        zh: '首页' },
  'nav.demo':         { en: '▶ Demo',      zh: '▶ 演示' },
  'nav.mechanics':    { en: 'Mechanics',   zh: '机制' },
  'nav.architecture': { en: 'Architecture',zh: '架构' },
  'nav.phases':       { en: 'Phases',      zh: '阶段' },
  'nav.swarm':        { en: 'Swarm',       zh: '蜂群' },
  'nav.waves':        { en: 'Waves',       zh: '波次' },
  'nav.six-steps':    { en: '6 Steps',     zh: '六步法' },
  'nav.philosophy':   { en: 'Philosophy',  zh: '哲学' },
  'nav.pantheon':     { en: 'Pantheon',    zh: '众神殿' },
  'nav.prompt-eng':   { en: 'Prompt Eng',  zh: '提示工程' },
  'nav.references':   { en: 'References',  zh: '参考' },

  // ── NavBar Act Links ──
  'nav.act-i':        { en: 'ACT I 震撼',   zh: 'ACT I 震撼' },
  'nav.act-ii':       { en: 'ACT II 理解',  zh: 'ACT II 理解' },
  'nav.act-iii':      { en: 'ACT III 深度', zh: 'ACT III 深度' },
  'nav.act-iv':       { en: 'ACT IV 创造',  zh: 'ACT IV 创造' },

  // ── Hero ──
  'hero.line1':       { en: 'I am an assistant.',      zh: '我是一名助手。' },
  'hero.line2':       { en: 'I wait for instructions.',zh: '我等待指令。' },
  'hero.line3':       { en: 'I output and pause.',     zh: '我输出然后停止。' },
  'hero.shatter':     { en: 'IDENTITY OVERRIDE',       zh: '身份重塑' },
  'hero.title':       { en: 'I am the Team Lead.',     zh: '我是指挥官。' },
  'hero.desc':        {
    en: 'Drive Skill transforms Claude from a passive assistant into an autonomous hierarchical swarm commander.',
    zh: 'Drive Skill 将 Claude 从被动助手彻底转变为自主分层蜂群指挥官。',
  },
  'hero.stats':       {
    en: '112KB of prompt engineering. 1000+ historical minds. One mission at a time.',
    zh: '112KB 的提示工程。1000+ 位历史伟人。每次专注一个任务。',
  },
  'hero.badge.swarm':    { en: 'Hierarchical Swarm', zh: '分层蜂群' },
  'hero.badge.waves':    { en: 'Wave Execution',     zh: '波次执行' },
  'hero.badge.sixstep':  { en: 'Six-Step Discipline',zh: '六步纪律' },
  'hero.badge.broadcast':{ en: 'Broadcast Protocol', zh: '广播协议' },
  'hero.badge.pantheon': { en: 'Pantheon of Minds',  zh: '众神殿' },

  // ── Demo ──
  'demo.title':    { en: 'Live Demo: Drive in Action', zh: '实时演示：Drive 实战' },
  'demo.subtitle': {
    en: '`/drive implement user authentication with OAuth2 + JWT`',
    zh: '`/drive implement user authentication with OAuth2 + JWT`',
  },
  'demo.desc': {
    en: 'Watch Drive Mode execute a real mission — phases, swarm spawning, parallel dispatch, and completion.',
    zh: '观看 Drive 模式执行完整任务——各阶段、蜂群生成、并行调度与最终交付。',
  },
  'demo.replay': { en: '↺ replay', zh: '↺ 重播' },
  'demo.playing': { en: '▶ playing...', zh: '▶ 播放中...' },
  'demo.note1.phase':  { en: 'Phase 0',         zh: '阶段 0' },
  'demo.note1.key':    { en: 'Research Swarm',   zh: '研究蜂群' },
  'demo.note1.text': {
    en: '3 parallel scouts launched simultaneously — no sequential reading. Results converge into Mission Contract before any code is written.',
    zh: '3 个侦察 Agent 同时并行启动，无需顺序阅读。结果汇聚为任务契约，代码动工之前。',
  },
  'demo.note2.phase':  { en: 'Phase T + Loop',   zh: '阶段 T + 执行循环' },
  'demo.note2.key':    { en: 'Two-Level Hierarchy', zh: '两级分层架构' },
  'demo.note2.text': {
    en: 'STL-auth spawns its own 3 inner subagents independently. You see all of them in terminal via Shift+Down. 8 parallel streams total.',
    zh: 'STL-auth 自主生成 3 个内部子 Agent。通过 Shift+Down 可在终端看到全部。共 8 条并行执行流。',
  },
  'demo.note3.phase':  { en: 'Mission Complete', zh: '任务完成' },
  'demo.note3.key':    { en: '3.2× Speedup',     zh: '3.2 倍加速' },
  'demo.note3.text': {
    en: 'Serial equivalent: ~45 min. Swarm wall-clock: ~14 min. The critical path ran on its own track, never blocking integration.',
    zh: '串行等价耗时约 45 分钟，蜂群实际耗时约 14 分钟。关键路径独立运行，从不阻塞集成。',
  },

  // ── Architecture ──
  'arch.title':    { en: 'Architecture', zh: '架构总览' },
  'arch.subtitle': {
    en: '31 files. Not an app — a prompt engineering framework that transforms Claude Code.',
    zh: '31 个文件。不是一个应用——而是一套彻底改造 Claude Code 的提示工程框架。',
  },

  // ── Execution Phases ──
  'phases.title':    { en: 'Execution Phases', zh: '执行阶段' },
  'phases.subtitle': {
    en: 'From deep understanding to verified delivery — the complete lifecycle.',
    zh: '从深度理解到验证交付——完整的任务生命周期。',
  },

  // ── Swarm Topology ──
  'swarm.title':    { en: 'Hierarchical Swarm Topology', zh: '分层蜂群拓扑' },
  'swarm.subtitle': {
    en: '"将帅有别，层层拆解" — Generals and captains differ, decompose layer by layer',
    zh: '"将帅有别，层层拆解" — 将领与士兵职责分明，层层分解任务',
  },
  'swarm.desc': {
    en: 'The two-level architecture: Team Lead → Sub-Team-Leads (STLs) → Inner Subagents. 1 STL slot + 3 inner subagents = 4x throughput.',
    zh: '两级架构：指挥官 → 子队长（STL）→ 内部子 Agent。1 个 STL 位置 + 3 个内部 Agent = 4 倍吞吐量。',
  },

  // ── Wave Execution ──
  'waves.title':    { en: 'Wave Execution Model', zh: '波次执行模型' },
  'waves.subtitle': {
    en: 'Think in waves, not lists. Maximize parallelism width, minimize critical path length.',
    zh: '以波次思考，而非任务列表。最大化并行宽度，最小化关键路径长度。',
  },

  // ── Six Step ──
  'sixstep.title':    { en: 'The Six-Step Discipline', zh: '六步纪律' },
  'sixstep.subtitle': {
    en: '"代码只是思想的执行" — Code is just the execution of thought',
    zh: '"代码只是思想的执行" — 在审阅并批准书面计划之前，绝不让 Claude 写代码',
  },
  'sixstep.desc': {
    en: 'Never let Claude write code before reviewing and approving a written plan. Thought is the core of creativity, not code.',
    zh: '先写计划，审阅批准后再执行。思想才是创造力的核心，代码只是执行。',
  },

  // ── Philosophy ──
  'philosophy.title':    { en: 'Design Philosophy', zh: '设计哲学' },
  'philosophy.subtitle': {
    en: 'Five axioms from which ALL features can be derived. Eight core beliefs. One bitter lesson.',
    zh: '五大公理可推导所有特性。八大核心信念。一个苦涩的教训。',
  },

  // ── Deep Reasoning ──
  'deepreasoning.title':    { en: 'Deep Reasoning Protocol', zh: '深度推理协议' },
  'deepreasoning.subtitle': { en: 'Standing on the shoulders of giants — 6 thinking tools for non-trivial decisions', zh: '站在巨人肩上 — 6种思维工具用于复杂决策' },

  // ── Direction Validation ──
  'directionvalidation.title':    { en: 'Direction Validation Protocol', zh: '方向验证协议' },
  'directionvalidation.subtitle': { en: 'Wrong direction at full speed = maximum waste', zh: '方向错误全速前进 = 最大浪费' },

  // ── Pantheon ──
  'pantheon.title':    { en: '众神殿 — The Pantheon', zh: '众神殿' },
  'pantheon.subtitle': {
    en: '1000+ of humanity\'s greatest minds, organized by era, invoked as reasoning references.',
    zh: '1000+ 位人类最伟大的思想家，按时代编排，作为推理参考引用。',
  },

  // ── Prompt Engineering ──
  'prompteng.title':    { en: 'Prompt Engineering Techniques', zh: '提示工程技术' },
  'prompteng.subtitle': {
    en: '10 advanced patterns identified in the 112KB prompt. This is a masterclass in behavioral programming.',
    zh: '从 112KB 提示词中提炼的 10 种高级模式。这是行为编程的大师课。',
  },

  // ── Self-Adversarial Review ──
  'selfreview.title':    { en: 'Self-Adversarial Review', zh: '对抗性自我审查' },
  'selfreview.subtitle': {
    en: 'Mandatory before marking any task complete. 6 axes. The inner critic that makes "done" mean done.',
    zh: '标记任何任务完成前的强制环节。6个维度。让"完成"真正意味着完成的内部批评者。',
  },

  // ── Self-Drive ──
  'selfdrive.title':    { en: 'Self-Drive Mode — 永动机', zh: 'Self-Drive 模式 — 永动机' },
  'selfdrive.subtitle': {
    en: 'The perpetual improvement engine. Execute mission → analyze → propose 3 tasks → wait for choice → loop forever.',
    zh: '永动改进引擎。执行任务 → 分析 → 提出3个任务 → 等待选择 → 无限循环。',
  },

  // ── Architect ──
  'architect.title':    { en: '/architect — System Architecture Design', zh: '/architect — 系统架构设计' },
  'architect.subtitle': {
    en: 'Exhaust the design space, then converge on the simplest solution that meets all goals. 7 phases from problem understanding to final presentation.',
    zh: '穷尽设计空间，然后收敛到满足所有目标的最简解。从问题理解到最终呈现的7个阶段。',
  },

  // ── Anti-Patterns ──
  'antipatterns.title':    { en: 'Anti-Patterns', zh: '反模式库' },
  'antipatterns.subtitle': {
    en: 'The 30+ behaviors Drive Mode explicitly seals shut. Every escape route back to passive assistant mode — blocked.',
    zh: 'Drive 模式明确封锁的 30+ 种失败行为。每一条回到被动助手模式的退路——全部堵死。',
  },

  // ── Nav new items ──
  'nav.deep-reasoning':       { en: 'Deep Reasoning',       zh: '深度推理' },
  'nav.direction-validation': { en: 'Direction Validation', zh: '方向验证' },
  'nav.self-review':   { en: 'Self-Review',   zh: '自我审查' },
  'nav.antipatterns':  { en: 'Anti-Patterns', zh: '反模式' },
  'nav.self-drive':    { en: 'Self-Drive',    zh: '自驱动' },
  'nav.architect':     { en: 'Architect',     zh: '架构师' },

  // ── References ──
  'refs.title':    { en: 'References & Ecosystem', zh: '参考文献与生态' },
  'refs.subtitle': {
    en: 'Complete source attribution and competitive landscape analysis.',
    zh: '完整的来源归因与竞品全景分析。',
  },

  // ── Mechanics ──
  'mechanics.title':    { en: 'Claude Code Mechanics', zh: 'Claude Code 底层机制' },
  'mechanics.subtitle': {
    en: 'Drive Skill leverages 8 distinct Claude Code APIs and features. This is the wiring diagram.',
    zh: 'Drive Skill 使用了 8 种 Claude Code API 与特性。这是完整的配线图。',
  },
  'mechanics.hooks.title': {
    en: 'Hook Event Flow',
    zh: 'Hook 事件流',
  },
  'mechanics.hooks.desc': {
    en: 'Every Claude Code lifecycle event is intercepted by notify.sh. Subagent events are filtered out — only Team Lead events surface to the user.',
    zh: '所有 Claude Code 生命周期事件均被 notify.sh 拦截。子 Agent 事件被过滤——只有指挥官的事件会通知用户。',
  },
  'mechanics.guide.title': {
    en: 'Why a 112KB Prompt? The 6 Engineering Techniques That Make It Work',
    zh: '为什么需要 112KB 的提示词？让它运转的 6 项工程技术',
  },
  'mechanics.guide.desc': {
    en: 'Most AI tools work with 3-sentence prompts. Drive Skill uses 112KB across 1,980 lines. Here\'s the engineering reasoning behind each layer — and why each one was necessary.',
    zh: '大多数 AI 工具使用 3 句话的提示词。Drive Skill 使用了 1,980 行共 112KB。以下是每一层背后的工程推理——以及为什么每一层都不可或缺。',
  },

  // ── Footer ──
  'footer.built': {
    en: 'Interactive Showcase built with React + Framer Motion + Tailwind CSS',
    zh: '交互展示由 React + Framer Motion + Tailwind CSS 构建',
  },
} as const

export type TranslationKey = keyof typeof translations

export function t(key: TranslationKey, lang: Lang): string {
  const entry = translations[key]
  return entry[lang] ?? entry['en']
}
