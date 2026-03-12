import { useState, useEffect } from 'react'
import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import CommandBadge from '../components/CommandBadge'
import { motion, AnimatePresence } from 'framer-motion'

// ─── Animated counter for stats ───
function AnimatedNumber({ target, duration = 2 }: { target: number; duration?: number }) {
  const [count, setCount] = useState(0)
  useEffect(() => {
    let start = 0
    const step = target / (duration * 60)
    const timer = setInterval(() => {
      start += step
      if (start >= target) {
        setCount(target)
        clearInterval(timer)
      } else {
        setCount(Math.floor(start))
      }
    }, 1000 / 60)
    return () => clearInterval(timer)
  }, [target, duration])
  return <>{count}</>
}

// ─── 1. Staging Paradox — SVG git diagram ───
function StagingParadoxViz({ lang }: { lang: 'en' | 'zh' }) {
  const [showSolution, setShowSolution] = useState(false)

  useEffect(() => {
    const t = setTimeout(() => setShowSolution(true), 2500)
    return () => clearTimeout(t)
  }, [])

  return (
    <div className="rounded-lg border border-border bg-bg-dark overflow-hidden">
      <div className="p-3 border-b border-border flex items-center gap-2">
        <span className="text-xs font-mono text-text-muted">
          {lang === 'en' ? 'The Staging Paradox' : '暂存悖论'}
        </span>
        <div className="flex gap-1 ml-auto">
          <motion.button
            onClick={() => setShowSolution(false)}
            className={`px-2 py-0.5 text-xs rounded font-mono ${!showSolution ? 'bg-neon-red/20 text-neon-red border border-neon-red/30' : 'text-text-muted border border-border'}`}
          >
            {lang === 'en' ? 'Problem' : '问题'}
          </motion.button>
          <motion.button
            onClick={() => setShowSolution(true)}
            className={`px-2 py-0.5 text-xs rounded font-mono ${showSolution ? 'bg-neon-green/20 text-neon-green border border-neon-green/30' : 'text-text-muted border border-border'}`}
          >
            {lang === 'en' ? 'Solution' : '解决方案'}
          </motion.button>
        </div>
      </div>
      <AnimatePresence mode="wait">
        {!showSolution ? (
          <motion.div
            key="problem"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="p-4"
          >
            <svg viewBox="0 0 700 160" className="w-full h-auto">
              <rect width="700" height="160" fill="transparent" />
              {/* main line with random merges */}
              <line x1="60" y1="40" x2="640" y2="40" stroke="#ef4444" strokeWidth="2" strokeOpacity="0.6" />
              <text x="10" y="44" fill="#ef4444" fontSize="11" fontFamily="monospace">main</text>
              {/* Feature branches merging chaotically */}
              {['A', 'B', 'C', 'D', 'E'].map((label, i) => {
                const x = 120 + i * 110
                return (
                  <motion.g key={label}
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: i * 0.15 }}
                  >
                    <line x1={x} y1="40" x2={x} y2="80" stroke="#888" strokeWidth="1.5" strokeDasharray="4 2" />
                    <text x={x - 3} y="95" fill="#888" fontSize="9" fontFamily="monospace">feat/{label}</text>
                    <circle cx={x} cy="40" r="4" fill="#ef4444" />
                    <text x={x - 12} y="35" fill="#888" fontSize="8" fontFamily="monospace">(merge)</text>
                  </motion.g>
                )
              })}
              {/* Staging attempt annotation */}
              <motion.g
                initial={{ opacity: 0 }}
                animate={{ opacity: [0, 1, 0.5, 1] }}
                transition={{ delay: 0.8, duration: 2, repeat: Infinity, repeatDelay: 2 }}
              >
                <line x1="175" y1="40" x2="175" y2="130" stroke="#f59e0b" strokeWidth="1" strokeDasharray="3 3" />
                <text x="180" y="120" fill="#f59e0b" fontSize="9" fontFamily="monospace">
                  {lang === 'en' ? '↑ deploy staging here' : '↑ 在此部署暂存'}
                </text>
                <text x="180" y="132" fill="#f59e0b" fontSize="8" fontFamily="monospace" opacity="0.7">
                  {lang === 'en' ? 'D and E merge in → unstable!' : 'D 和 E 合并进来 → 不稳定！'}
                </text>
                <text x="180" y="152" fill="#ef4444" fontSize="9" fontFamily="monospace" fontWeight="bold">
                  {lang === 'en' ? '✗ can never verify a stable snapshot' : '✗ 永远无法验证稳定快照'}
                </text>
              </motion.g>
            </svg>
          </motion.div>
        ) : (
          <motion.div
            key="solution"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="p-4"
          >
            <svg viewBox="0 0 700 180" className="w-full h-auto">
              <rect width="700" height="180" fill="transparent" />
              {/* develop line */}
              <motion.line x1="60" y1="40" x2="580" y2="40" stroke="#3b82f6" strokeWidth="2"
                initial={{ pathLength: 0 }} animate={{ pathLength: 1 }} transition={{ duration: 0.8 }} />
              <text x="10" y="44" fill="#3b82f6" fontSize="11" fontFamily="monospace">develop</text>
              {/* Commits on develop */}
              {['A', 'B', 'C', 'D', 'E'].map((label, i) => (
                <motion.g key={label}
                  initial={{ opacity: 0, scale: 0 }}
                  animate={{ opacity: 1, scale: 1 }}
                  transition={{ delay: 0.3 + i * 0.1 }}
                >
                  <circle cx={120 + i * 90} cy={40} r="4" fill="#3b82f6" />
                  <text x={117 + i * 90} y="30" fill="#3b82f6" fontSize="9" fontFamily="monospace">{label}</text>
                </motion.g>
              ))}
              <motion.text x="440" y="30" fill="#8b949e" fontSize="9" fontFamily="monospace"
                initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 0.8 }}>
                {lang === 'en' ? '← AI noise goes here' : '← AI 噪音在这里'}
              </motion.text>
              {/* RC fork */}
              <motion.g initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 1.2 }}>
                <path d="M 300 40 C 340 40 340 100 360 100" stroke="#f59e0b" strokeWidth="2" fill="none" />
                <line x1="360" y1="100" x2="520" y2="100" stroke="#f59e0b" strokeWidth="2" />
                <text x="370" y="95" fill="#f59e0b" fontSize="11" fontFamily="monospace">rc/V0.1.0</text>
                <circle cx={360} cy={100} r="4" fill="#f59e0b" />
                <circle cx={440} cy={100} r="4" fill="#f59e0b" />
                <text x="325" y="70" fill="#f59e0b" fontSize="8" fontFamily="monospace" opacity="0.8">
                  {lang === 'en' ? 'frozen snapshot' : '冻结快照'}
                </text>
              </motion.g>
              {/* main line */}
              <motion.g initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 1.6 }}>
                <line x1="60" y1="150" x2="580" y2="150" stroke="#ef4444" strokeWidth="2" />
                <text x="10" y="154" fill="#ef4444" fontSize="11" fontFamily="monospace">main</text>
                {/* Merge from RC to main */}
                <path d="M 440 100 C 480 100 480 150 520 150" stroke="#f59e0b" strokeWidth="2" fill="none" strokeDasharray="6 3" />
                <circle cx={520} cy={150} r="5" fill="#ef4444" stroke="#f59e0b" strokeWidth="2" />
                <text x="525" y="145" fill="#ef4444" fontSize="10" fontFamily="monospace" fontWeight="bold">[V0.1.0]</text>
                <text x="450" y="170" fill="#22c55e" fontSize="9" fontFamily="monospace" fontWeight="bold">
                  {lang === 'en' ? '← only verified code' : '← 只有验证过的代码'}
                </text>
              </motion.g>
            </svg>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}

// ─── 2. Squash Comparison — animated ───
function SquashViz({ lang }: { lang: 'en' | 'zh' }) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {/* Without squash */}
      <motion.div
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
        className="rounded-lg border border-neon-red/30 bg-neon-red/5 p-4"
      >
        <div className="text-xs font-mono text-neon-red font-bold mb-2">
          {lang === 'en' ? 'WITHOUT SQUASH' : '不使用 SQUASH'}
        </div>
        <div className="space-y-1">
          {['fix typo', 'WIP', 'try approach', 'revert', 'fix', 'forgot save', '...'].map((msg, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.08 }}
              className="flex items-center gap-2 text-xs font-mono text-text-muted"
            >
              <span className="w-1.5 h-1.5 rounded-full bg-text-muted/30 shrink-0" />
              {msg}
            </motion.div>
          ))}
        </div>
        <div className="mt-2 text-xs font-mono text-neon-red">
          47 commits. {lang === 'en' ? 'None meaningful.' : '全无意义。'}
        </div>
      </motion.div>
      {/* With squash */}
      <motion.div
        initial={{ opacity: 0, x: 20 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ delay: 0.3 }}
        className="rounded-lg border border-neon-green/30 bg-neon-green/5 p-4"
      >
        <div className="text-xs font-mono text-neon-green font-bold mb-2">
          {lang === 'en' ? 'WITH SQUASH' : '使用 SQUASH'}
        </div>
        <div className="space-y-1.5">
          {[
            'feat: user login system (#42)',
            'feat: permissions (#45)',
            'fix: API timeout (#48)',
            'release: V0.1.0',
          ].map((msg, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.5 + i * 0.12 }}
              className="flex items-center gap-2 text-xs font-mono text-neon-green"
            >
              <span className="w-2 h-2 rounded-full bg-neon-green shrink-0" />
              {msg}
            </motion.div>
          ))}
        </div>
        <div className="mt-2 text-xs font-mono text-neon-green">
          4 commits. {lang === 'en' ? 'Each one meaningful.' : '每一个都有意义。'}
        </div>
      </motion.div>
    </div>
  )
}

// ─── 3. Automation Stats — what Teamwork automates ───
function AutomationShowcase({ lang }: { lang: 'en' | 'zh' }) {
  const automations = lang === 'en' ? [
    { icon: '📝', label: 'Contract generation', detail: 'AI scans codebase, discovers files, structures sub-tasks' },
    { icon: '🏷️', label: 'Label management', detail: 'mission, status:*, priority:*, size:* — auto-synced' },
    { icon: '🔀', label: 'Branch lifecycle', detail: 'Create → squash merge → delete — zero manual git' },
    { icon: '🔍', label: '12-check diagnostics', detail: 'Orphan branches, stale stashes, cross-contamination...' },
    { icon: '📦', label: 'Post-merge cleanup', detail: 'Close issue, update labels, check milestones — auto' },
    { icon: '🔧', label: 'Auto-fix & repair', detail: '/team doctor fix — execute safe actions automatically' },
    { icon: '🏗️', label: 'CI/hooks scaffolding', detail: 'ci.yml, pre-commit, pre-push — generated on /team init' },
    { icon: '💾', label: 'Resume from interruption', detail: 'Checkboxes in Contract preserve progress — zero replay' },
  ] : [
    { icon: '📝', label: '合约自动生成', detail: 'AI 扫描代码库，发现文件，结构化子任务' },
    { icon: '🏷️', label: '标签自动管理', detail: 'mission, status:*, priority:*, size:* — 自动同步' },
    { icon: '🔀', label: '分支全生命周期', detail: '创建 → squash 合并 → 删除 — 零手动 git' },
    { icon: '🔍', label: '12 项自动诊断', detail: '孤儿分支、过期暂存、交叉污染...' },
    { icon: '📦', label: '合并后自动清理', detail: '关闭 issue、更新标签、检查里程碑 — 全自动' },
    { icon: '🔧', label: '自动修复', detail: '/team doctor fix — 自动执行安全操作' },
    { icon: '🏗️', label: 'CI/hooks 脚手架', detail: 'ci.yml, pre-commit, pre-push — /team init 自动生成' },
    { icon: '💾', label: '中断恢复', detail: '合约中的复选框保存进度 — 零重复工作' },
  ]

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
      {automations.map((a, i) => (
        <motion.div
          key={i}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: i * 0.06, type: 'spring', stiffness: 300 }}
          whileHover={{ scale: 1.02, x: 4 }}
          className="flex items-start gap-2.5 p-2.5 rounded-lg bg-bg-card border border-border hover:border-neon-blue/30 transition-colors"
        >
          <span className="text-base mt-0.5">{a.icon}</span>
          <div>
            <div className="text-xs font-medium text-text-primary">{a.label}</div>
            <div className="text-[10px] text-text-muted leading-tight">{a.detail}</div>
          </div>
        </motion.div>
      ))}
    </div>
  )
}

// ─── 4. Before/After efficiency comparison ───
function EfficiencyComparison({ lang }: { lang: 'en' | 'zh' }) {
  const [showAfter, setShowAfter] = useState(false)

  useEffect(() => {
    const t = setTimeout(() => setShowAfter(true), 2000)
    return () => clearTimeout(t)
  }, [])

  const beforeItems = lang === 'en' ? [
    '"What branch am I on?"',
    '"What\'s the Issue number?"',
    '"Did I push?"',
    '"Is CI passing?"',
    '"Did I close the Issue?"',
    '"Did I update the labels?"',
    '"What version is this?"',
    '"Who\'s working on what?"',
  ] : [
    '"我在哪个分支？"',
    '"Issue 编号是什么？"',
    '"我推送了吗？"',
    '"CI 通过了吗？"',
    '"我关闭 Issue 了吗？"',
    '"我更新标签了吗？"',
    '"这是什么版本？"',
    '"谁在做什么？"',
  ]

  const afterItems = [
    { cmd: '/team-claim #42', desc: lang === 'en' ? 'Branch + Contract' : '分支 + 合约' },
    { cmd: '/team-drive', desc: lang === 'en' ? 'Code + test loop' : '编码 + 测试循环' },
    { cmd: '/team-ship', desc: lang === 'en' ? 'PR + CI + review' : 'PR + CI + 审查' },
    { cmd: '/team-rc', desc: lang === 'en' ? 'Release lifecycle' : '发布生命周期' },
    { cmd: '/team', desc: lang === 'en' ? 'Who does what' : '谁在做什么' },
  ]

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {/* Before */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: showAfter ? 0.5 : 1 }}
        className="rounded-lg border border-neon-red/30 bg-neon-red/5 p-4"
      >
        <div className="text-xs font-mono text-neon-red font-bold mb-3">
          {lang === 'en' ? 'BEFORE — Mental overhead' : '之前 — 心智负担'}
        </div>
        <div className="space-y-1.5">
          {beforeItems.map((item, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.06 }}
              className="text-xs font-mono text-text-muted flex items-center gap-2"
            >
              <motion.span
                animate={{ opacity: showAfter ? [1, 0.3] : 1 }}
                transition={{ delay: i * 0.05 }}
              >
                ❓
              </motion.span>
              {showAfter ? <span className="line-through">{item}</span> : <span>{item}</span>}
            </motion.div>
          ))}
        </div>
      </motion.div>
      {/* After */}
      <motion.div
        initial={{ opacity: 0, x: 20 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ delay: 0.5 }}
        className="rounded-lg border border-neon-green/30 bg-neon-green/5 p-4"
      >
        <div className="text-xs font-mono text-neon-green font-bold mb-3">
          {lang === 'en' ? 'AFTER — Six commands' : '之后 — 六个命令'}
        </div>
        <div className="space-y-2">
          {afterItems.map((item, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: 10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.8 + i * 0.1, type: 'spring' }}
              className="flex items-center gap-2"
            >
              <CommandBadge size="sm">{item.cmd}</CommandBadge>
              <span className="text-xs text-text-muted">{item.desc}</span>
            </motion.div>
          ))}
        </div>
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1.5 }}
          className="mt-3 pt-3 border-t border-neon-green/20 text-xs font-mono text-neon-green"
        >
          {lang === 'en'
            ? 'Not six new things to learn — six moments that used to require context-switching.'
            : '不是六个新东西 — 而是六个原本需要手动切换上下文的时刻。'}
        </motion.div>
      </motion.div>
    </div>
  )
}

// ─── 5. Worktree parallel execution visual ───
function WorktreeValueViz({ lang }: { lang: 'en' | 'zh' }) {
  const tabs = [
    { branch: 'pre-launch', color: '#3b82f6', icon: '🟦', task: lang === 'en' ? 'Dashboard + orchestrate' : '仪表盘 + 协调' },
    { branch: 'mission/042-auth', color: '#22c55e', icon: '🟩', task: lang === 'en' ? 'Bob: JWT auth (3/5 tasks)' : 'Bob: JWT 认证 (3/5 任务)' },
    { branch: 'mission/048-payment', color: '#22c55e', icon: '🟩', task: lang === 'en' ? 'Alice: Fix payment (DONE)' : 'Alice: 修复支付 (完成)' },
  ]

  return (
    <div className="rounded-lg border border-border bg-bg-dark overflow-hidden">
      {/* iTerm-style tab bar */}
      <div className="flex bg-[#1c2128] border-b border-border">
        {tabs.map((tab, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.15 }}
            className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-mono border-r border-border"
            style={{ borderBottom: `2px solid ${tab.color}` }}
          >
            <span className="text-[10px]">{tab.icon}</span>
            <span style={{ color: tab.color }}>{tab.branch}</span>
          </motion.div>
        ))}
      </div>
      {/* Content showing parallel work */}
      <div className="p-4 space-y-3">
        {tabs.map((tab, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.5 + i * 0.15 }}
            className="flex items-center gap-3"
          >
            <div className="w-2 h-2 rounded-full shrink-0" style={{ backgroundColor: tab.color }} />
            <div className="flex-1 font-mono text-xs">
              <span style={{ color: tab.color }}>{tab.branch}</span>
              <span className="text-text-muted ml-2">→ {tab.task}</span>
            </div>
            {i === 2 && (
              <motion.span
                animate={{ scale: [1, 1.2, 1] }}
                transition={{ repeat: Infinity, duration: 2 }}
                className="text-neon-green text-xs"
              >
                ✅
              </motion.span>
            )}
            {i === 1 && (
              <span className="text-xs font-mono text-[rgb(177,185,249)]">
                ██████░░░░ 60%
              </span>
            )}
          </motion.div>
        ))}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1.2 }}
          className="pt-3 border-t border-border"
        >
          <div className="flex items-center gap-2 text-xs font-mono">
            <span className="text-text-muted">
              {lang === 'en' ? 'Each agent has its own directory — zero conflicts, true parallelism' : '每个 agent 独立目录 — 零冲突，真正并行'}
            </span>
          </div>
          <div className="mt-2 space-y-0.5 text-[10px] font-mono text-text-muted">
            <div>~/project/ <span className="text-[#3b82f6]">← pre-launch (main repo)</span></div>
            <div>~/project.worktrees/042-auth/ <span className="text-[#22c55e]">← mission/042 (isolated)</span></div>
            <div>~/project.worktrees/048-payment/ <span className="text-[#22c55e]">← mission/048 (isolated)</span></div>
          </div>
        </motion.div>
      </div>
    </div>
  )
}

// ─── 6. Six Skills Lifecycle Flow — animated pipeline ───
function LifecycleFlow({ lang }: { lang: 'en' | 'zh' }) {
  const steps = [
    { cmd: '/team-issue', label: lang === 'en' ? 'Issue\ncreated' : 'Issue\n创建', color: '#3b82f6' },
    { cmd: '/team-claim', label: lang === 'en' ? 'Contract +\nbranch' : '合约 +\n分支', color: '#22c55e' },
    { cmd: '/team-drive', label: lang === 'en' ? 'Code + test\nloop' : '编码 + 测试\n循环', color: '#f59e0b' },
    { cmd: '/team-ship', label: lang === 'en' ? 'PR created\n(squash)' : 'PR 创建\n(squash)', color: '#a855f7' },
    { cmd: '/team-ship done', label: lang === 'en' ? 'Cleanup\nbranch deleted' : '清理\n分支删除', color: '#ef4444' },
  ]

  return (
    <div className="rounded-lg border border-border bg-bg-dark p-4 overflow-x-auto">
      <div className="flex items-center gap-1 min-w-[600px]">
        {steps.map((step, i) => (
          <motion.div key={i} className="flex items-center gap-1 flex-1">
            <motion.div
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: i * 0.15, type: 'spring' }}
              className="flex flex-col items-center gap-1.5 flex-1"
            >
              <div className="px-2 py-1 rounded text-[10px] font-mono font-bold text-center"
                style={{ color: step.color, backgroundColor: step.color + '15', border: `1px solid ${step.color}40` }}>
                {step.cmd}
              </div>
              <motion.div
                className="w-3 h-3 rounded-full"
                style={{ backgroundColor: step.color }}
                animate={{ scale: [1, 1.3, 1] }}
                transition={{ delay: 1 + i * 0.3, duration: 0.5 }}
              />
              <div className="text-[9px] font-mono text-text-muted text-center whitespace-pre-line leading-tight">
                {step.label}
              </div>
            </motion.div>
            {i < steps.length - 1 && (
              <motion.div
                initial={{ opacity: 0, scaleX: 0 }}
                animate={{ opacity: 1, scaleX: 1 }}
                transition={{ delay: 0.3 + i * 0.15 }}
                className="flex items-center shrink-0 mt-[-16px]"
              >
                <div className="w-6 h-px bg-text-muted/30" />
                <motion.span
                  animate={{ x: [0, 4, 0] }}
                  transition={{ repeat: Infinity, duration: 1.5, delay: i * 0.2 }}
                  className="text-text-muted text-xs"
                >
                  →
                </motion.span>
              </motion.div>
            )}
          </motion.div>
        ))}
      </div>
    </div>
  )
}

// ─── 7. Quality Gates — Four layers ───
function QualityGatesViz({ lang }: { lang: 'en' | 'zh' }) {
  const gates = lang === 'en' ? [
    { layer: 'L1', name: 'Local Hooks', speed: 'seconds', detail: 'pre-commit lint, pre-push test', color: '#22c55e' },
    { layer: 'L2', name: 'GitHub CI', speed: 'minutes', detail: 'lint → test → build on every PR', color: '#3b82f6' },
    { layer: 'L3', name: 'Branch Protection', speed: 'instant', detail: 'Require PR, CI pass, review', color: '#f59e0b' },
    { layer: 'L4', name: 'Post-merge Auto', speed: 'auto', detail: 'Labels, checkboxes, milestone check', color: '#a855f7' },
  ] : [
    { layer: 'L1', name: '本地 Hooks', speed: '秒级', detail: 'pre-commit lint, pre-push test', color: '#22c55e' },
    { layer: 'L2', name: 'GitHub CI', speed: '分钟级', detail: 'lint → test → build (每个 PR)', color: '#3b82f6' },
    { layer: 'L3', name: '分支保护', speed: '即时', detail: '要求 PR、CI 通过、代码审查', color: '#f59e0b' },
    { layer: 'L4', name: '合并后自动化', speed: '自动', detail: '标签、复选框、里程碑检查', color: '#a855f7' },
  ]

  return (
    <div className="space-y-1.5">
      {gates.map((g, i) => (
        <motion.div
          key={i}
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: i * 0.1 }}
          whileHover={{ x: 6 }}
          className="flex items-center gap-3 p-2 rounded-lg bg-bg-card border border-border"
        >
          <span className="text-xs font-mono font-bold w-6 shrink-0" style={{ color: g.color }}>{g.layer}</span>
          <div className="w-1 h-6 rounded-full shrink-0" style={{ backgroundColor: g.color }} />
          <div className="flex-1 min-w-0">
            <span className="text-xs font-medium text-text-primary">{g.name}</span>
            <span className="text-[10px] text-text-muted ml-2">({g.speed})</span>
          </div>
          <span className="text-[10px] font-mono text-text-muted hidden sm:block">{g.detail}</span>
        </motion.div>
      ))}
    </div>
  )
}

// ─── Main Ch0 Component ───
const commands = [
  { cmd: '/team', desc: { en: 'Initialize & dashboard', zh: '初始化 & 仪表盘' } },
  { cmd: '/team-issue', desc: { en: 'AI-enriched issue creation', zh: 'AI 增强的任务创建' } },
  { cmd: '/team-claim', desc: { en: 'Contract + branch + worktree', zh: '合约 + 分支 + worktree' } },
  { cmd: '/team-drive', desc: { en: 'Autonomous execution loop', zh: '自主执行循环' } },
  { cmd: '/team-ship', desc: { en: 'PR + CI + review + cleanup', zh: 'PR + CI + 审查 + 清理' } },
  { cmd: '/team-rc', desc: { en: 'Release candidate lifecycle', zh: '发布候选生命周期' } },
]

export default function Ch0Why() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[0]

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      {/* Hero intro */}
      <p className="text-lg text-text-secondary">
        {lang === 'en'
          ? 'AI coding assistants generate 20-30x more commits than human-only teams. Without structure, this noise drowns your project.'
          : 'AI 编程助手产生比纯人工团队多 20-30 倍的提交。没有结构，这些噪音会淹没你的项目。'}
      </p>

      {/* Stats bar */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex flex-wrap gap-4 justify-center p-4 rounded-lg bg-bg-card border border-border"
      >
        {[
          { value: 6, unit: lang === 'en' ? 'commands' : '个命令', detail: lang === 'en' ? 'total' : '总计' },
          { value: 12, unit: lang === 'en' ? 'health checks' : '项健康检查', detail: '/team doctor' },
          { value: 4, unit: lang === 'en' ? 'quality layers' : '层质量门禁', detail: lang === 'en' ? 'hooks → CI → protection → auto' : 'hooks → CI → 保护 → 自动' },
          { value: 0, unit: lang === 'en' ? 'infrastructure' : '基础设施', detail: lang === 'en' ? 'GitHub-native' : 'GitHub 原生' },
        ].map((stat, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.2 + i * 0.1, type: 'spring' }}
            className="text-center px-4"
          >
            <div className="text-2xl font-bold text-neon-blue font-mono">
              <AnimatedNumber target={stat.value} />
            </div>
            <div className="text-xs text-text-primary font-medium">{stat.unit}</div>
            <div className="text-[10px] text-text-muted">{stat.detail}</div>
          </motion.div>
        ))}
      </motion.div>

      {/* Section 1: The Staging Paradox — git-style SVG */}
      <div className="space-y-3">
        <h2 className="text-lg font-bold text-text-primary">
          {lang === 'en' ? 'The Problem: Staging Paradox' : '问题：暂存悖论'}
        </h2>
        <p className="text-sm text-text-secondary">
          {lang === 'en'
            ? 'With feature branches merging directly to main, you can never deploy a stable snapshot. New merges invalidate staging while you verify.'
            : '当 feature 分支直接合并到 main 时，你永远无法部署稳定快照。新的合并在你验证时就让暂存失效。'}
        </p>
        <StagingParadoxViz lang={lang} />
      </div>

      {/* Section 2: Squash — Clean History */}
      <div className="space-y-3">
        <h2 className="text-lg font-bold text-text-primary">
          {lang === 'en' ? 'AI Noise → Clean History' : 'AI 噪音 → 清晰历史'}
        </h2>
        <p className="text-sm text-text-secondary">
          {lang === 'en'
            ? 'AI generates 17+ commits per feature. Squash merge: 1 Issue = 1 meaningful commit. Cherry-pickable, bisectable, readable.'
            : 'AI 每个功能生成 17+ 个提交。Squash 合并：1 个 Issue = 1 个有意义的提交。可 cherry-pick、可 bisect、可读。'}
        </p>
        <SquashViz lang={lang} />
      </div>

      {/* Section 3: Before/After */}
      <div className="space-y-3">
        <h2 className="text-lg font-bold text-text-primary">
          {lang === 'en' ? 'Before & After' : '前后对比'}
        </h2>
        <p className="text-sm text-text-secondary">
          {lang === 'en'
            ? 'Six moments in a developer\'s day that used to require context-switching and manual bookkeeping. Now they don\'t.'
            : '开发者一天中六个原本需要手动切换上下文和记账的时刻。现在不需要了。'}
        </p>
        <EfficiencyComparison lang={lang} />
      </div>

      {/* Section 4: Six Skills Lifecycle */}
      <div className="space-y-3">
        <h2 className="text-lg font-bold text-text-primary">
          {lang === 'en' ? 'The Six-Skill Lifecycle' : '六技能生命周期'}
        </h2>
        <LifecycleFlow lang={lang} />
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-2">
          {commands.map((c, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.3 + i * 0.06 }}
              className="flex items-center gap-2 p-2 rounded-lg bg-bg-card border border-border"
            >
              <CommandBadge size="sm">{c.cmd}</CommandBadge>
              <span className="text-text-muted text-xs">{c.desc[lang]}</span>
            </motion.div>
          ))}
        </div>
      </div>

      {/* Section 5: What Teamwork Automates */}
      <div className="space-y-3">
        <h2 className="text-lg font-bold text-text-primary">
          {lang === 'en' ? 'What Teamwork Automates' : 'Teamwork 自动化了什么'}
        </h2>
        <p className="text-sm text-text-secondary">
          {lang === 'en'
            ? 'Scripts lock deterministic operations. LLM provides judgment. 21 bugs found in LLM-interpreted ops, 0 in scripts.'
            : '脚本锁定确定性操作，LLM 提供判断。21 个 bug 出现在 LLM 解释的操作中，脚本中为 0。'}
        </p>
        <AutomationShowcase lang={lang} />
      </div>

      {/* Section 6: Worktree — True Parallelism */}
      <div className="space-y-3">
        <h2 className="text-lg font-bold text-text-primary">
          {lang === 'en' ? 'Worktree: True Parallel Execution' : 'Worktree：真正的并行执行'}
        </h2>
        <p className="text-sm text-text-secondary">
          {lang === 'en'
            ? 'Each agent gets its own directory. No branch switching, no stashing, no conflicts. Multiple missions execute simultaneously in separate iTerm tabs.'
            : '每个 agent 有独立的目录。不切分支、不暂存、不冲突。多个任务在独立的 iTerm tab 中同时执行。'}
        </p>
        <WorktreeValueViz lang={lang} />
      </div>

      {/* Section 7: Quality Gates */}
      <div className="space-y-3">
        <h2 className="text-lg font-bold text-text-primary">
          {lang === 'en' ? 'Four-Layer Quality Gates' : '四层质量门禁'}
        </h2>
        <p className="text-sm text-text-secondary">
          {lang === 'en'
            ? 'From milliseconds (local hooks) to automatic (post-merge cleanup). Layer 4 catches human forgetfulness — the #1 cause of stale boards.'
            : '从毫秒级（本地 hooks）到全自动（合并后清理）。第 4 层捕获人类的遗忘 — 看板过期的头号原因。'}
        </p>
        <QualityGatesViz lang={lang} />
      </div>

      {/* iTerm2 recommendation */}
      <motion.div
        initial={{ opacity: 0, y: 12 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.5 }}
        className="p-5 rounded-lg border border-neon-blue/30 bg-neon-blue/5"
      >
        <div className="flex items-start gap-4">
          <div className="text-3xl shrink-0">🖥️</div>
          <div>
            <h3 className="font-bold text-neon-blue mb-1">
              {lang === 'en' ? 'Best with iTerm2' : '推荐使用 iTerm2'}
            </h3>
            <p className="text-text-secondary text-sm mb-3">
              {lang === 'en'
                ? 'Teamwork integrates with iTerm2 for tab colors (🟥🟦🟧🟩 by branch type), desktop notifications (click → jump to tab), and worktree isolation. Strongly recommended.'
                : 'Teamwork 与 iTerm2 深度集成：tab 颜色 (🟥🟦🟧🟩 按分支类型)、桌面通知 (点击 → 跳转到 tab)、worktree 隔离。强烈推荐。'}
            </p>
            <motion.a
              href="https://iterm2.com"
              target="_blank"
              rel="noopener noreferrer"
              whileHover={{ scale: 1.03 }}
              whileTap={{ scale: 0.97 }}
              className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-neon-blue text-white text-sm font-medium hover:brightness-110 transition-all"
            >
              {lang === 'en' ? 'Download iTerm2' : '下载 iTerm2'}
              <span>↗</span>
            </motion.a>
          </div>
        </div>
      </motion.div>
    </ChapterLayout>
  )
}
