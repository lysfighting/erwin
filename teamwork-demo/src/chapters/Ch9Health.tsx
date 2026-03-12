import { useState } from 'react'
import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import TerminalPlayground from '../components/TerminalPlayground'
import CommandBadge from '../components/CommandBadge'
import { ch9HealthMocks } from '../data/mockData'
import { motion } from 'framer-motion'

const descriptions = {
  en: [
    {
      heading: 'Run diagnostics',
      body: '/team doctor scans branches, PRs, stashes, worktrees, and issues — then produces a diagnosis with a prioritized action plan.',
    },
    {
      heading: 'Auto-fix',
      body: '/team doctor fix executes safe actions automatically (delete orphan branches, drop stale stashes, prune worktrees) and skips those needing manual review.',
    },
  ],
  zh: [
    {
      heading: '运行诊断',
      body: '/team doctor 扫描分支、PR、暂存、工作树和 Issue — 生成诊断报告和优先级行动计划。',
    },
    {
      heading: '自动修复',
      body: '/team doctor fix 自动执行安全操作（删除孤儿分支、清除过期暂存、修剪工作树），跳过需要人工审查的项目。',
    },
  ],
} as const

// The real 12 doctor checks
const checks = [
  { label: 'Orphan branches', pass: false },
  { label: 'Cross-contam.', pass: true },
  { label: 'Branch staleness', pass: false },
  { label: 'Stale stashes', pass: false },
  { label: 'Worktree health', pass: true },
  { label: 'Shippable PRs', pass: true },
  { label: 'Issue health', pass: false },
  { label: 'Remote sync', pass: true },
  { label: 'GitHub auth', pass: true },
  { label: 'Contract fresh.', pass: true },
  { label: 'Test suite', pass: true },
  { label: 'PR CI status', pass: true },
]

export default function Ch9Health() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[9]
  const desc = descriptions[lang]
  const [activeCase, setActiveCase] = useState(0)

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      <p className="text-text-secondary">
        {lang === 'en'
          ? <><CommandBadge>/team doctor</CommandBadge> diagnoses your project health with 12 automated checks. <CommandBadge>/team doctor fix</CommandBadge> auto-repairs issues.</>
          : <><CommandBadge>/team doctor</CommandBadge> 用 12 项自动检查诊断项目健康状态。<CommandBadge>/team doctor fix</CommandBadge> 自动修复问题。</>}
      </p>

      {/* Visual check grid */}
      <div className="grid grid-cols-4 md:grid-cols-6 gap-2">
        {checks.map((c, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: i * 0.05, type: 'spring', stiffness: 300 }}
            className={`p-2 rounded border text-center text-xs ${
              c.pass
                ? 'border-neon-green/30 bg-neon-green/5 text-neon-green'
                : 'border-neon-red/30 bg-neon-red/5 text-neon-red'
            }`}
          >
            <div className="text-lg mb-0.5">{c.pass ? '✓' : '✗'}</div>
            {c.label}
          </motion.div>
        ))}
      </div>

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

      <motion.p
        key={activeCase}
        initial={{ opacity: 0, y: 6 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-text-primary leading-relaxed"
      >
        {desc[activeCase].body}
      </motion.p>

      <TerminalPlayground
        key={`term-${activeCase}`}
        steps={ch9HealthMocks[activeCase]}
        title="teamwork-doctor"
      />
    </ChapterLayout>
  )
}
