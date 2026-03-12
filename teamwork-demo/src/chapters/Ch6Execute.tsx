import { useState } from 'react'
import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import TerminalPlayground from '../components/TerminalPlayground'
import CommandBadge from '../components/CommandBadge'
import ITermTabBar, { type ITermTab } from '../components/ITermTabBar'
import ITermNotification from '../components/ITermNotification'
import { ch6ExecuteMocks } from '../data/mockData'
import { motion } from 'framer-motion'

const descriptions = {
  en: [
    {
      heading: 'Full drive loop',
      body: '/team-drive reads the contract, decomposes into sub-tasks, and executes each with a verify loop. If a test fails, it auto-fixes. Watch a full authentication feature built from scratch.',
    },
    {
      heading: 'Resume after interruption',
      body: 'Connection dropped? /team-drive detects saved state and resumes exactly where it left off — no duplicate work. Progress is saved via checkboxes in the Contract.',
    },
  ],
  zh: [
    {
      heading: '完整驱动循环',
      body: '/team-drive 读取合约，分解为子任务，通过验证循环逐个执行。测试失败自动修复。完整观看一个认证功能从零构建。',
    },
    {
      heading: '中断后恢复',
      body: '连接断了？/team-drive 检测保存的状态，从中断处精确恢复 — 不重复工作。进度通过合约中的复选框保存。',
    },
  ],
} as const

// Multi-agent parallel execution demo with interactive tabs
function ParallelDriveDemo({ lang }: { lang: 'en' | 'zh' }) {
  const tabs: ITermTab[] = [
    { label: 'pre-launch (claude)', color: '#3b82f6' },
    { label: 'mission/042-add-auth-bob (claude)', color: '#22c55e', active: true },
    { label: 'mission/048-fix-payment-alice (claude)', color: '#22c55e', hasNotification: true },
  ]

  const prelaunchContent = (
    <div className="p-3 font-mono text-xs space-y-0.5">
      <div className="text-[#e6edf3] font-bold">MY DASHBOARD — acme-corp/taskflow</div>
      <div className="text-[#8b949e]">Teamwork v3.8.3 · V1.0  <span className="text-[rgb(177,185,249)]">████░░░░░░</span>  25%  <span className="text-[rgb(177,185,249)]">1</span>/<span className="text-[rgb(177,185,249)]">4</span></div>
      <div className="text-[#8b949e] mt-1">👥 <span className="text-[#e6edf3] font-bold">TEAM</span> (3 members)</div>
      <div className="text-[#8b949e]">  alice (tech-lead)</div>
      <div className="text-[#8b949e]">    🟢 <span className="text-blue-400 font-bold">#048</span> Fix payment crash 🟠 P0 → <span className="italic">wip</span></div>
      <div className="text-[#8b949e]">  bob (backend)</div>
      <div className="text-[#8b949e]">    🟢 <span className="text-blue-400 font-bold">#042</span> Add user auth 🟠 P1 → <span className="italic">wip</span></div>
      <div className="text-[#8b949e]">  charlie (frontend)</div>
      <div className="text-[#8b949e]">    ⚪ <span className="italic">(idle)</span></div>
    </div>
  )

  const mission042Content = (
    <div className="p-3 font-mono text-xs space-y-0.5">
      <div className="text-[#e6edf3] font-bold">🎯 MISSION ── <span className="text-blue-400">#042</span> Add user authentication</div>
      <div className="text-[#8b949e]">🟠 P1  🔀 <span className="italic">mission/042-add-auth-bob</span></div>
      <div className="text-[#8b949e]">Progress:  <span className="text-[rgb(177,185,249)]">██████░░░░</span>  <span className="text-[rgb(177,185,249)]">3</span>/<span className="text-[rgb(177,185,249)]">5</span></div>
      <div className="text-[#8b949e] mt-1">📋 <span className="text-[#e6edf3] font-bold">SUB-TASKS</span></div>
      <div className="text-[#8b949e]">  ✅ Install jsonwebtoken + bcrypt</div>
      <div className="text-[#8b949e]">  ✅ Create auth service</div>
      <div className="text-[#8b949e]">  ✅ Implement login/register routes</div>
      <div className="text-[#f59e0b]">  ▸ Add JWT middleware  ← working</div>
      <div className="text-[#8b949e]">  ○ Write unit tests</div>
    </div>
  )

  const mission048Content = (
    <div className="p-3 font-mono text-xs space-y-0.5">
      <div className="text-[#e6edf3] font-bold">🎯 MISSION ── <span className="text-blue-400">#048</span> Fix payment crash</div>
      <div className="text-[#8b949e]">🟠 P0  🔀 <span className="italic">mission/048-fix-payment-alice</span></div>
      <div className="text-[#8b949e]">Progress:  <span className="text-[rgb(177,185,249)]">████████░░</span>  <span className="text-[rgb(177,185,249)]">2</span>/<span className="text-[rgb(177,185,249)]">2</span></div>
      <div className="text-[#8b949e] mt-1">📋 <span className="text-[#e6edf3] font-bold">SUB-TASKS</span></div>
      <div className="text-[#8b949e]">  ✅ Add null check for card.expiry</div>
      <div className="text-[#8b949e]">  ✅ Write regression test</div>
      <div className="text-[#8b949e] mt-1">────────────────────────────────────────────</div>
      <div className="text-[#e6edf3] font-bold">🏁 COMPLETE ── <span className="text-blue-400">#048</span> Fix payment crash</div>
      <div className="text-[#8b949e]"><span className="text-[rgb(177,185,249)]">/team-ship</span> to create PR</div>
    </div>
  )

  return (
    <div>
      <p className="text-text-muted text-sm mb-2">
        {lang === 'en'
          ? '↓ Multiple agents execute in parallel, each in their own worktree tab. Click tabs to see each agent\'s progress:'
          : '↓ 多个 agent 并行执行，各自在独立的 worktree tab 中。点击 tab 查看每个 agent 的进度：'}
      </p>
      <div className="flex flex-col md:flex-row gap-4">
        <div className="flex-1">
          <ITermTabBar
            tabs={tabs}
            tabContents={[prelaunchContent, mission042Content, mission048Content]}
            terminalTitle="mission/042-add-auth-bob (claude)"
          />
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
              {
                title: 'Alert',
                subtitle: 'Session pre-launch (claude) #3:',
                body: 'Claude is waiting for your input',
                time: '20m ago',
                icon: 'iterm',
              },
            ]}
            autoPlayDelay={1500}
          />
        </div>
      </div>
    </div>
  )
}

export default function Ch6Execute() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[6]
  const desc = descriptions[lang]
  const [activeCase, setActiveCase] = useState(0)

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      <p className="text-text-secondary">
        {lang === 'en'
          ? <><CommandBadge>/team-drive</CommandBadge> is the autonomous execution engine. It reads your Contract, validates the plan against actual code, classifies tasks (code vs manual), and executes with a verify loop.</>
          : <><CommandBadge>/team-drive</CommandBadge> 是自主执行引擎。读取合约，对照实际代码验证计划，分类任务（代码 vs 手动），通过验证循环执行。</>}
      </p>

      {/* What drive automates */}
      <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
        {(lang === 'en' ? [
          { icon: '📖', label: 'Reads Contract', detail: 'Validates plan vs codebase' },
          { icon: '🔍', label: 'Plan Assessment', detail: 'Confirms, revises, or adds tasks' },
          { icon: '⚡', label: 'Auto-scales', detail: '1-3: sequential, 4-7: team, 8+: STL hierarchy' },
          { icon: '🧪', label: 'Test-fix loop', detail: '3 auto-fix attempts per failure' },
          { icon: '💾', label: 'Checkpoint save', detail: 'Checkboxes preserve progress' },
          { icon: '🔧', label: 'Manual handoff', detail: 'AI skips manual tasks, hands off to you' },
        ] : [
          { icon: '📖', label: '读取合约', detail: '对照代码验证计划' },
          { icon: '🔍', label: '计划评估', detail: '确认、修订或增加任务' },
          { icon: '⚡', label: '自动扩展', detail: '1-3: 串行, 4-7: 团队, 8+: STL 层级' },
          { icon: '🧪', label: '测试-修复循环', detail: '每次失败 3 次自动修复' },
          { icon: '💾', label: '断点保存', detail: '复选框保存进度' },
          { icon: '🔧', label: '手动交接', detail: 'AI 跳过手动任务，交给你' },
        ]).map((item, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.06, type: 'spring' }}
            className="flex items-start gap-2 p-2 rounded-lg bg-bg-card border border-border text-xs"
          >
            <span>{item.icon}</span>
            <div>
              <div className="font-medium text-text-primary">{item.label}</div>
              <div className="text-[10px] text-text-muted">{item.detail}</div>
            </div>
          </motion.div>
        ))}
      </div>

      {/* Parallel execution demo */}
      <ParallelDriveDemo lang={lang} />

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
        steps={ch6ExecuteMocks[activeCase]}
        title="teamwork-drive"
      />
    </ChapterLayout>
  )
}
