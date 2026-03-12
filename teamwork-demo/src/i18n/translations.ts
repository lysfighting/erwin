const translations = {
  en: {
    appTitle: 'Teamwork Tutorial',
    appSubtitle: 'Learn the AI-native GitHub workflow',
    prev: 'Prev',
    next: 'Next',
    comingSoon: 'Coming soon...',
    chapters: [
      { title: 'Why Teamwork?', subtitle: 'The collaboration challenge in the AI era' },
      { title: 'Git Branch Model', subtitle: 'Four branches, clear responsibilities' },
      { title: 'Issue Lifecycle', subtitle: 'From idea to deployment' },
      { title: 'Setup', subtitle: 'Initialize your project with /team' },
      { title: 'Create Issues', subtitle: 'AI-enriched issue creation with /team-issue' },
      { title: 'Claim Missions', subtitle: 'Contract-based claiming with /team-claim' },
      { title: 'Execute', subtitle: 'Drive development with /team-drive' },
      { title: 'Deliver', subtitle: 'Ship code with /team-ship' },
      { title: 'Shortcuts', subtitle: 'auto, wrap, and queue operations' },
      { title: 'Health Check', subtitle: 'Diagnose with /team doctor' },
      { title: 'Release', subtitle: 'Cut and promote with /team-rc' },
      { title: 'Help & Learn', subtitle: 'Command reference and concepts guide' },
      { title: 'Worktree Mode', subtitle: 'Isolated parallel execution with git worktrees' },
    ],
  },
  zh: {
    appTitle: 'Teamwork 教程',
    appSubtitle: '学习 AI 原生的 GitHub 工作流',
    prev: '上一章',
    next: '下一章',
    comingSoon: '即将推出...',
    chapters: [
      { title: '为什么需要 Teamwork？', subtitle: 'AI 时代的协作挑战' },
      { title: 'Git 分支模型', subtitle: '四条分支，职责清晰' },
      { title: 'Issue 生命周期', subtitle: '从想法到部署' },
      { title: '项目配置', subtitle: '用 /team 初始化项目' },
      { title: '创建任务', subtitle: '用 /team-issue 创建 AI 增强的任务' },
      { title: '认领任务', subtitle: '用 /team-claim 基于合约认领' },
      { title: '执行开发', subtitle: '用 /team-drive 驱动开发' },
      { title: '交付代码', subtitle: '用 /team-ship 交付代码' },
      { title: '快捷操作', subtitle: 'auto、wrap 和 queue 操作' },
      { title: '健康检查', subtitle: '用 /team doctor 诊断项目' },
      { title: '发布上线', subtitle: '用 /team-rc 切版和发布' },
      { title: '帮助与学习', subtitle: '命令参考和概念指南' },
      { title: 'Worktree 模式', subtitle: '用 git worktree 实现隔离并行执行' },
    ],
  },
} as const

export type Translations = typeof translations.en

export default translations
