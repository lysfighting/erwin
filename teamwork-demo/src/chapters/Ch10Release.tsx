import { useState } from 'react'
import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import TerminalPlayground from '../components/TerminalPlayground'
import CommandBadge from '../components/CommandBadge'
import { ch10ReleaseMocks } from '../data/mockData'
import { motion } from 'framer-motion'

const descriptions = {
  en: [
    {
      heading: 'Cut RC',
      body: '/team-rc creates an RC branch from pre-launch, runs all pre-release checks (lint, typecheck, tests, coverage), and prepares for QA.',
    },
    {
      heading: 'Promote to production',
      body: '/team-rc promote squash-merges RC to main, creates a git tag, publishes a GitHub Release with changelog, and cleans up.',
    },
  ],
  zh: [
    {
      heading: '创建 RC',
      body: '/team-rc 从 pre-launch 创建 RC 分支，运行所有预发布检查（lint、类型、测试、覆盖率），准备 QA。',
    },
    {
      heading: '发布到生产',
      body: '/team-rc promote 将 RC squash 合并到 main，创建 git 标签，发布 GitHub Release 并清理。',
    },
  ],
} as const

export default function Ch10Release() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[10]
  const desc = descriptions[lang]
  const [activeCase, setActiveCase] = useState(0)

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      <p className="text-text-secondary">
        {lang === 'en'
          ? <><CommandBadge>/team-rc</CommandBadge> manages the release lifecycle — cut from pre-launch, QA on staging, then <CommandBadge>/team-rc promote</CommandBadge> to ship.</>
          : <><CommandBadge>/team-rc</CommandBadge> 管理发布生命周期 — 从 pre-launch 切出 RC，在 staging 验证，然后 <CommandBadge>/team-rc promote</CommandBadge> 发布。</>}
      </p>

      {/* Release flow visualization */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.2 }}
        className="flex items-center gap-3 flex-wrap text-sm p-4 rounded-lg bg-bg-card border border-border"
      >
        <div className="flex items-center gap-1.5">
          <span className="text-neon-blue">🟦</span>
          <span className="font-mono text-neon-blue text-xs">pre-launch</span>
        </div>
        <motion.span animate={{ x: [0, 4, 0] }} transition={{ repeat: Infinity, duration: 1.5 }} className="text-gold">→</motion.span>
        <div className="flex items-center gap-1.5">
          <span className="text-gold">🟧</span>
          <span className="font-mono text-gold text-xs">rc/V1.0.0</span>
        </div>
        <span className="text-text-muted text-xs">(QA)</span>
        <motion.span animate={{ x: [0, 4, 0] }} transition={{ repeat: Infinity, duration: 1.5, delay: 0.5 }} className="text-neon-green">→</motion.span>
        <div className="flex items-center gap-1.5">
          <span className="text-neon-red">🟥</span>
          <span className="font-mono text-neon-red text-xs">main</span>
        </div>
        <motion.span
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1 }}
          className="px-2 py-0.5 rounded bg-neon-green/10 border border-neon-green/30 text-neon-green text-xs font-mono"
        >
          V1.0.0 🏷️
        </motion.span>
      </motion.div>

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
        steps={ch10ReleaseMocks[activeCase]}
        title="teamwork-release"
      />
    </ChapterLayout>
  )
}
