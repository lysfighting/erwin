import { useState } from 'react'
import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import TerminalPlayground from '../components/TerminalPlayground'
import CommandBadge from '../components/CommandBadge'
import { ch3SetupMocks } from '../data/mockData'
import { motion } from 'framer-motion'

const descriptions = {
  en: [
    {
      heading: '/team init',
      body: 'Run /team in any repo. Teamwork scans the project — detects language, discovers contributors from git history, creates .teamwork/config.yml, syncs GitHub labels, creates a milestone. Zero questions.',
    },
    {
      heading: '/team config',
      body: 'View the full configuration: shared config (committed) and local settings (per-user). Includes base branch, production branch, team roster, labels, milestone, and worktree mode.',
    },
    {
      heading: 'Dashboard',
      body: 'After issues are created, /team shows the live dashboard: milestone progress bar, member assignments with status, unassigned queue, merge stats, and context-specific next action.',
    },
  ],
  zh: [
    {
      heading: '/team 初始化',
      body: '在任何仓库中运行 /team，Teamwork 扫描项目 — 检测语言、从 git 历史发现贡献者、创建 .teamwork/config.yml、同步 GitHub 标签、创建里程碑。零提问，即装即用。',
    },
    {
      heading: '/team config',
      body: '查看完整配置：共享配置（已提交）和本地设置（每用户）。包括基础分支、生产分支、团队名册、标签、里程碑和 worktree 模式。',
    },
    {
      heading: '仪表盘',
      body: '创建任务后，/team 展示实时仪表盘：里程碑进度条、成员分配与状态、未认领队列、合并统计和上下文相关的下一步操作。',
    },
  ],
} as const

// Init flow visualization
function InitFlowViz({ lang }: { lang: 'en' | 'zh' }) {
  const steps = lang === 'en'
    ? [
        { icon: '🔍', label: 'Scan project', detail: 'language, package manager, remote' },
        { icon: '👥', label: 'Detect team', detail: 'git log → contributors + roles' },
        { icon: '📝', label: 'Write config', detail: '.teamwork/config.yml' },
        { icon: '🏷️', label: 'Sync labels', detail: 'mission, status:*, priority:*, size:*' },
        { icon: '🎯', label: 'Create milestone', detail: 'from versions.current config' },
        { icon: '📊', label: 'Show dashboard', detail: 'role-based view' },
      ]
    : [
        { icon: '🔍', label: '扫描项目', detail: '语言、包管理器、远程仓库' },
        { icon: '👥', label: '检测团队', detail: 'git log → 贡献者 + 角色' },
        { icon: '📝', label: '写入配置', detail: '.teamwork/config.yml' },
        { icon: '🏷️', label: '同步标签', detail: 'mission, status:*, priority:*, size:*' },
        { icon: '🎯', label: '创建里程碑', detail: '从 versions.current 配置' },
        { icon: '📊', label: '显示仪表盘', detail: '基于角色的视图' },
      ]

  return (
    <div className="flex items-center gap-1 flex-wrap p-3 rounded-lg bg-bg-card border border-border">
      {steps.map((s, i) => (
        <motion.div key={i} className="flex items-center gap-1">
          <motion.div
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: i * 0.1, type: 'spring', stiffness: 300 }}
            className="flex items-center gap-1.5 px-2 py-1 rounded bg-bg-dark border border-border"
          >
            <span className="text-sm">{s.icon}</span>
            <div>
              <span className="text-xs text-text-primary font-medium">{s.label}</span>
              <span className="text-[10px] text-text-muted block">{s.detail}</span>
            </div>
          </motion.div>
          {i < steps.length - 1 && (
            <motion.span
              animate={{ x: [0, 3, 0] }}
              transition={{ repeat: Infinity, duration: 1.2, delay: i * 0.15 }}
              className="text-neon-green text-xs"
            >
              →
            </motion.span>
          )}
        </motion.div>
      ))}
    </div>
  )
}

// What /team init auto-generates
function InitGeneratesViz({ lang }: { lang: 'en' | 'zh' }) {
  const files = lang === 'en'
    ? [
        { path: '.teamwork/config.yml', detail: 'Team roster, conventions, test commands', icon: '⚙️' },
        { path: '.github/ISSUE_TEMPLATE/mission.yml', detail: 'Structured mission issue template', icon: '📋' },
        { path: '.github/workflows/ci.yml', detail: 'lint → test → build on every PR', icon: '🔄' },
        { path: '.githooks/pre-commit', detail: 'Runs linter before each commit', icon: '🔒' },
        { path: '.githooks/pre-push', detail: 'Runs tests before each push', icon: '🧪' },
      ]
    : [
        { path: '.teamwork/config.yml', detail: '团队名册、约定、测试命令', icon: '⚙️' },
        { path: '.github/ISSUE_TEMPLATE/mission.yml', detail: '结构化任务 Issue 模板', icon: '📋' },
        { path: '.github/workflows/ci.yml', detail: 'lint → test → build (每个 PR)', icon: '🔄' },
        { path: '.githooks/pre-commit', detail: '每次提交前运行 linter', icon: '🔒' },
        { path: '.githooks/pre-push', detail: '每次推送前运行测试', icon: '🧪' },
      ]

  const labels = ['mission', 'status:queued', 'status:wip', 'status:review', 'status:done', 'priority:P0', 'priority:P1', 'priority:P2', 'priority:P3', 'size:S', 'size:M', 'size:L', 'size:XL']

  return (
    <div className="space-y-3">
      <div className="text-xs font-medium text-text-secondary">
        {lang === 'en' ? '↓ Files auto-generated by /team init:' : '↓ /team init 自动生成的文件：'}
      </div>
      <div className="space-y-1.5">
        {files.map((f, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, x: -15 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: i * 0.08, type: 'spring', stiffness: 300 }}
            whileHover={{ x: 4 }}
            className="flex items-center gap-2.5 p-2 rounded-lg bg-bg-card border border-border"
          >
            <span className="text-sm">{f.icon}</span>
            <code className="text-xs text-[rgb(177,185,249)] font-mono">{f.path}</code>
            <span className="text-[10px] text-text-muted ml-auto hidden sm:block">{f.detail}</span>
          </motion.div>
        ))}
      </div>
      <div className="text-xs font-medium text-text-secondary mt-2">
        {lang === 'en' ? '↓ GitHub labels auto-synced:' : '↓ GitHub 标签自动同步：'}
      </div>
      <div className="flex flex-wrap gap-1.5">
        {labels.map((label, i) => (
          <motion.span
            key={i}
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.4 + i * 0.03, type: 'spring' }}
            className="px-2 py-0.5 rounded-full text-[10px] font-mono border"
            style={{
              color: label.startsWith('priority') ? '#f59e0b' : label.startsWith('status') ? '#3b82f6' : label.startsWith('size') ? '#a855f7' : '#22c55e',
              borderColor: label.startsWith('priority') ? '#f59e0b40' : label.startsWith('status') ? '#3b82f640' : label.startsWith('size') ? '#a855f740' : '#22c55e40',
              backgroundColor: label.startsWith('priority') ? '#f59e0b10' : label.startsWith('status') ? '#3b82f610' : label.startsWith('size') ? '#a855f710' : '#22c55e10',
            }}
          >
            {label}
          </motion.span>
        ))}
      </div>
    </div>
  )
}

export default function Ch3Setup() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[3]
  const desc = descriptions[lang]
  const [activeCase, setActiveCase] = useState(0)

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      {/* Intro */}
      <p className="text-text-secondary">
        {lang === 'en'
          ? <>The <CommandBadge>/team</CommandBadge> command is your entry point. Run it once to initialize, then anytime to see your dashboard.</>
          : <><CommandBadge>/team</CommandBadge> 是你的入口。运行一次初始化，之后随时查看仪表盘。</>}
      </p>

      {/* Init flow visualization */}
      <div>
        <p className="text-text-muted text-sm mb-2">
          {lang === 'en' ? '↓ What happens when you run /team for the first time:' : '↓ 首次运行 /team 时的流程：'}
        </p>
        <InitFlowViz lang={lang} />
      </div>

      {/* What init generates */}
      <InitGeneratesViz lang={lang} />

      {/* Case selector */}
      <div className="flex gap-2 flex-wrap">
        {desc.map((d, i) => (
          <motion.button
            key={i}
            onClick={() => setActiveCase(i)}
            whileHover={{ scale: 1.03 }}
            whileTap={{ scale: 0.97 }}
            className={`px-3 py-1.5 text-sm font-mono rounded border transition-colors ${
              activeCase === i
                ? 'border-neon-blue bg-neon-blue/10 text-neon-blue'
                : 'border-border bg-bg-card text-text-secondary hover:bg-bg-card-hover'
            }`}
          >
            {d.heading}
          </motion.button>
        ))}
      </div>

      {/* Description */}
      <motion.p
        key={activeCase}
        initial={{ opacity: 0, y: 6 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-text-primary leading-relaxed"
      >
        {desc[activeCase].body}
      </motion.p>

      {/* Terminal */}
      <TerminalPlayground
        key={activeCase}
        steps={ch3SetupMocks[activeCase]}
        title="teamwork-setup"
      />
    </ChapterLayout>
  )
}
