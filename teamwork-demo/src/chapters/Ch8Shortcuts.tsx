import { useState } from 'react'
import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import TerminalPlayground from '../components/TerminalPlayground'
import CommandBadge from '../components/CommandBadge'
import { ch8ShortcutsMocks } from '../data/mockData'
import { motion } from 'framer-motion'

const descriptions = {
  en: [
    {
      heading: 'Auto mode',
      body: '/team auto #43 runs the full workflow in one shot: claim → drive → ship. For well-defined, straightforward issues.',
    },
    {
      heading: 'Retroactive wrap',
      body: 'Already wrote code without a mission? /team wrap detects uncommitted changes, infers the purpose, creates issue + branch + contract + PR retroactively.',
    },
  ],
  zh: [
    {
      heading: '自动模式',
      body: '/team auto #43 一键完整工作流：认领 → 驱动 → 交付。适合定义清晰的任务。',
    },
    {
      heading: '追溯包装',
      body: '写了代码但没创建任务？/team wrap 检测未提交更改，推断目的，追溯创建 Issue + 分支 + 合约 + PR。',
    },
  ],
} as const

export default function Ch8Shortcuts() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[8]
  const desc = descriptions[lang]
  const [activeCase, setActiveCase] = useState(0)

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      <p className="text-text-secondary">
        {lang === 'en'
          ? <>Power-user shortcuts: <CommandBadge>/team auto</CommandBadge> for one-shot execution and <CommandBadge>/team wrap</CommandBadge> for retroactive tracking.</>
          : <>高级快捷操作：<CommandBadge>/team auto</CommandBadge> 一键执行，<CommandBadge>/team wrap</CommandBadge> 追溯跟踪。</>}
      </p>

      {/* Flow diagram */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.2 }}
        className="flex items-center gap-2 flex-wrap text-sm p-3 rounded-lg bg-bg-card border border-border"
      >
        {activeCase === 0 ? (
          <>
            <span className="text-text-muted">{lang === 'en' ? 'auto =' : 'auto ='}</span>
            <CommandBadge size="sm">/team-claim</CommandBadge>
            <motion.span animate={{ x: [0, 4, 0] }} transition={{ repeat: Infinity, duration: 1.5 }} className="text-neon-green">→</motion.span>
            <CommandBadge size="sm">/team-drive</CommandBadge>
            <motion.span animate={{ x: [0, 4, 0] }} transition={{ repeat: Infinity, duration: 1.5, delay: 0.5 }} className="text-neon-green">→</motion.span>
            <CommandBadge size="sm">/team-ship</CommandBadge>
          </>
        ) : (
          <>
            <span className="text-text-muted">{lang === 'en' ? 'wrap =' : 'wrap ='}</span>
            <span className="text-text-secondary">{lang === 'en' ? 'detect changes' : '检测更改'}</span>
            <motion.span animate={{ x: [0, 4, 0] }} transition={{ repeat: Infinity, duration: 1.5 }} className="text-gold">→</motion.span>
            <CommandBadge size="sm">/team-issue</CommandBadge>
            <motion.span animate={{ x: [0, 4, 0] }} transition={{ repeat: Infinity, duration: 1.5, delay: 0.3 }} className="text-gold">→</motion.span>
            <CommandBadge size="sm">/team-claim</CommandBadge>
            <motion.span animate={{ x: [0, 4, 0] }} transition={{ repeat: Infinity, duration: 1.5, delay: 0.6 }} className="text-gold">→</motion.span>
            <CommandBadge size="sm">/team-ship</CommandBadge>
          </>
        )}
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
        steps={ch8ShortcutsMocks[activeCase]}
        title="teamwork-shortcuts"
      />
    </ChapterLayout>
  )
}
