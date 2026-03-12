import { useState } from 'react'
import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import TerminalPlayground from '../components/TerminalPlayground'
import CommandBadge from '../components/CommandBadge'
import ITermTabBar, { type ITermTab } from '../components/ITermTabBar'
import ITermNotification from '../components/ITermNotification'
import { ch12WorktreeMocks } from '../data/mockData'
import { motion } from 'framer-motion'

const descriptions = {
  en: [
    {
      heading: 'Enable worktree',
      body: 'One git config command enables worktree mode. Each /team-claim creates a separate directory with its own checkout — no stashing, no branch-switching.',
    },
    {
      heading: 'Claim with worktree',
      body: '/team-claim in worktree mode creates a git worktree + new iTerm tab. The terminal warns you NOT to run /team-drive in the main repo.',
    },
    {
      heading: 'Multiple missions',
      body: 'With worktree mode, multiple agents work in parallel — each in their own directory. No conflicts, no interference.',
    },
  ],
  zh: [
    {
      heading: '启用 worktree',
      body: '一条 git config 命令启用 worktree 模式。每次 /team-claim 创建独立目录和 checkout — 无需 stash，无需切分支。',
    },
    {
      heading: 'Worktree 认领',
      body: 'worktree 模式下 /team-claim 创建 git worktree + 新 iTerm tab。终端会警告你不要在主仓库运行 /team-drive。',
    },
    {
      heading: '多任务并行',
      body: 'worktree 模式下多个 agent 并行工作 — 各自在独立目录中。无冲突，无干扰。',
    },
  ],
} as const

const worktreeTabs: ITermTab[] = [
  { label: 'pre-launch (main repo)', color: '#3b82f6' },
  { label: 'mission/042-add-auth-bob (worktree)', color: '#22c55e', active: true },
  { label: 'mission/048-fix-payment-alice (worktree)', color: '#22c55e', hasNotification: true },
]

const benefits = {
  en: [
    { icon: '🔀', title: 'No branch switching', desc: 'Each mission has its own directory. No git stash, no checkout conflicts.' },
    { icon: '⚡', title: 'True parallelism', desc: 'Multiple agents work simultaneously in separate worktrees.' },
    { icon: '🛡️', title: 'Isolation', desc: 'A crash in one mission doesn\'t affect others. Clean rollback per worktree.' },
    { icon: '🖥️', title: 'iTerm2 integration', desc: 'Each worktree gets its own tab with color-coded branch indicator.' },
  ],
  zh: [
    { icon: '🔀', title: '无需切分支', desc: '每个任务有独立目录。无需 git stash，无 checkout 冲突。' },
    { icon: '⚡', title: '真正并行', desc: '多个 agent 在独立 worktree 中同时工作。' },
    { icon: '🛡️', title: '隔离性', desc: '一个任务崩溃不影响其他。每个 worktree 独立回滚。' },
    { icon: '🖥️', title: 'iTerm2 集成', desc: '每个 worktree 有独立 tab，带颜色分支标识。' },
  ],
}

export default function Ch12Worktree() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[12]
  const desc = descriptions[lang]
  const bens = benefits[lang]
  const [activeCase, setActiveCase] = useState(0)

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      <p className="text-text-secondary">
        {lang === 'en'
          ? <>Worktree mode gives each mission its own git directory. Enable with <code className="text-[rgb(177,185,249)]">git config --local teamwork.worktree true</code>, then every <CommandBadge>/team-claim</CommandBadge> creates an isolated workspace.</>
          : <>Worktree 模式让每个任务拥有独立的 git 目录。用 <code className="text-[rgb(177,185,249)]">git config --local teamwork.worktree true</code> 启用，之后每次 <CommandBadge>/team-claim</CommandBadge> 都创建隔离工作区。</>}
      </p>

      {/* Branch-only vs Worktree comparison */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          className="rounded-lg border border-neon-red/30 bg-neon-red/5 p-4"
        >
          <div className="text-xs font-mono text-neon-red font-bold mb-3">
            {lang === 'en' ? 'BRANCH-ONLY MODE' : '仅分支模式'}
          </div>
          <div className="space-y-2 text-xs font-mono text-text-muted">
            <div className="flex items-center gap-2"><span className="text-neon-red">✗</span> {lang === 'en' ? 'git stash → checkout → stash pop' : 'git stash → checkout → stash pop'}</div>
            <div className="flex items-center gap-2"><span className="text-neon-red">✗</span> {lang === 'en' ? 'One branch at a time' : '一次只能一个分支'}</div>
            <div className="flex items-center gap-2"><span className="text-neon-red">✗</span> {lang === 'en' ? 'Context switch = lost flow' : '切换上下文 = 丢失状态'}</div>
            <div className="flex items-center gap-2"><span className="text-neon-red">✗</span> {lang === 'en' ? 'Risk of cross-contamination' : '存在交叉污染风险'}</div>
          </div>
        </motion.div>
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.2 }}
          className="rounded-lg border border-neon-green/30 bg-neon-green/5 p-4"
        >
          <div className="text-xs font-mono text-neon-green font-bold mb-3">
            {lang === 'en' ? 'WORKTREE MODE' : 'Worktree 模式'}
          </div>
          <div className="space-y-2 text-xs font-mono text-text-muted">
            <div className="flex items-center gap-2"><span className="text-neon-green">✓</span> {lang === 'en' ? 'Each mission = own directory' : '每个任务 = 独立目录'}</div>
            <div className="flex items-center gap-2"><span className="text-neon-green">✓</span> {lang === 'en' ? 'Multiple missions in parallel' : '多任务真正并行'}</div>
            <div className="flex items-center gap-2"><span className="text-neon-green">✓</span> {lang === 'en' ? 'Zero context switching' : '零上下文切换'}</div>
            <div className="flex items-center gap-2"><span className="text-neon-green">✓</span> {lang === 'en' ? 'Complete isolation per agent' : '每个 agent 完全隔离'}</div>
          </div>
        </motion.div>
      </div>

      {/* Benefits grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {bens.map((b, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.08, type: 'spring', stiffness: 200 }}
            whileHover={{ scale: 1.02, x: 4 }}
            className="p-3 rounded-lg border border-neon-green/20 bg-neon-green/5"
          >
            <div className="text-xl mb-1">{b.icon}</div>
            <h3 className="font-bold text-text-primary text-sm mb-1">{b.title}</h3>
            <p className="text-text-secondary text-xs">{b.desc}</p>
          </motion.div>
        ))}
      </div>

      {/* Tab color legend */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.3 }}
        className="p-3 rounded-lg border border-border bg-bg-card"
      >
        <div className="text-xs font-medium text-text-primary mb-2">
          {lang === 'en' ? 'iTerm2 Tab Colors by Branch Type:' : 'iTerm2 Tab 颜色（按分支类型）：'}
        </div>
        <div className="flex flex-wrap gap-3">
          {[
            { icon: '🟥', branch: 'main', color: '#ef4444' },
            { icon: '🟦', branch: 'pre-launch', color: '#3b82f6' },
            { icon: '🟧', branch: 'rc/*', color: '#f59e0b' },
            { icon: '🟩', branch: 'mission/*', color: '#22c55e' },
            { icon: '⬜', branch: 'other', color: '#8b949e' },
          ].map((item, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: 0.4 + i * 0.05, type: 'spring' }}
              className="flex items-center gap-1.5 text-xs font-mono"
            >
              <span>{item.icon}</span>
              <span style={{ color: item.color }}>{item.branch}</span>
            </motion.div>
          ))}
        </div>
      </motion.div>

      {/* Worktree directory structure visualization */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.3 }}
        className="p-4 rounded-lg border border-border bg-bg-card font-mono text-sm"
      >
        <div className="text-text-muted mb-2">{lang === 'en' ? '# Directory layout with worktree mode:' : '# Worktree 模式的目录结构：'}</div>
        <div className="text-text-secondary space-y-0.5">
          <div><span className="text-neon-blue">~/Projects/taskflow/</span> <span className="text-text-muted">← main repo (pre-launch)</span></div>
          <div><span className="text-neon-green">~/Projects/taskflow.worktrees/042-add-auth-bob/</span> <span className="text-text-muted">← worktree</span></div>
          <div><span className="text-neon-green">~/Projects/taskflow.worktrees/048-fix-payment-alice/</span> <span className="text-text-muted">← worktree</span></div>
        </div>
      </motion.div>

      {/* iTerm multi-tab demo */}
      <div>
        <p className="text-text-muted text-sm mb-2">
          {lang === 'en'
            ? '↓ Each worktree gets its own iTerm tab with branch color:'
            : '↓ 每个 worktree 有独立的 iTerm tab 和分支颜色：'}
        </p>
        <div className="flex flex-col md:flex-row gap-4">
          <div className="flex-1">
            <ITermTabBar tabs={worktreeTabs} terminalTitle="mission/042-add-auth-bob (claude)">
              <div className="p-3 font-mono text-xs space-y-0.5">
                <div className="text-neon-green">$ claude</div>
                <div className="text-text-secondary">{lang === 'en' ? '# Working in isolated worktree' : '# 在隔离的 worktree 中工作'}</div>
                <div className="text-gold mt-1">▸ /team-drive</div>
              </div>
            </ITermTabBar>
          </div>
          <div className="shrink-0">
            <ITermNotification
              notifications={[
                {
                  title: 'Alert',
                  subtitle: 'Session mission/048-fix-payment (claude) #4:',
                  body: 'Mission complete! Run /team-ship to deliver.',
                  time: '2m ago',
                  icon: 'iterm',
                },
              ]}
              autoPlayDelay={1000}
            />
          </div>
        </div>
      </div>

      {/* Case selector + terminal */}
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
        steps={ch12WorktreeMocks[activeCase]}
        title="teamwork-worktree"
      />
    </ChapterLayout>
  )
}
