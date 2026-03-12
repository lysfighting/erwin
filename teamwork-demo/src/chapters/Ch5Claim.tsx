import { useState } from 'react'
import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import TerminalPlayground from '../components/TerminalPlayground'
import CommandBadge from '../components/CommandBadge'
import ITermTabBar, { type ITermTab } from '../components/ITermTabBar'
import { ch5ClaimMocks } from '../data/mockData'
import { motion } from 'framer-motion'

const descriptions = {
  en: [
    {
      heading: 'Claim a mission',
      body: '/team-claim #42 reads the issue, generates a Contract with context files, test commands, and acceptance criteria, then creates a mission branch — automatically.',
    },
    {
      heading: 'Claim (worktree)',
      body: 'With worktree mode enabled, /team-claim also creates a git worktree — an isolated directory for this mission. You get a warning: do NOT run /team-drive in the main repo.',
    },
    {
      heading: 'Browse missions',
      body: '/team-claim list shows all open issues with priority and assignment status.',
    },
  ],
  zh: [
    {
      heading: '认领任务',
      body: '/team-claim #42 读取 Issue，生成包含上下文文件、测试命令和验收标准的合约，然后创建 mission 分支 — 全自动。',
    },
    {
      heading: '认领 (worktree)',
      body: '启用 worktree 模式后，/team-claim 还会创建 git worktree — 该任务的隔离目录。终端会警告：不要在主仓库运行 /team-drive。',
    },
    {
      heading: '浏览任务',
      body: '/team-claim list 显示所有待认领的 Issue，包含优先级和分配状态。',
    },
  ],
} as const

// Worktree demo: before claim (1 tab) → after claim (2 tabs, new one slides in)
function WorktreeClaimDemo({ lang }: { lang: 'en' | 'zh' }) {
  const [claimed, setClaimed] = useState(false)

  const baseTabs: ITermTab[] = [
    { label: 'pre-launch (claude)', color: '#3b82f6', active: !claimed },
  ]

  const claimedTabs: ITermTab[] = [
    { label: 'pre-launch (claude)', color: '#3b82f6' },
    { label: 'mission/042-add-auth-bob (claude)', color: '#22c55e', active: true },
  ]

  const tabs = claimed ? claimedTabs : baseTabs

  // Tab content: pre-launch shows dashboard, mission shows claim output
  const prelaunchContent = (
    <div className="p-3 font-mono text-xs space-y-0.5">
      <div className="text-[#e6edf3] font-bold">MY DASHBOARD — acme-corp/taskflow</div>
      <div className="text-[#8b949e]">Teamwork v3.8.3 · V1.0  <span className="text-[rgb(177,185,249)]">████░░░░░░</span>  25%</div>
      <div className="text-[#8b949e] mt-1">👥 TEAM (3 members)</div>
      <div className="text-[#8b949e]">  bob (backend)</div>
      <div className="text-[#8b949e]">    ⚪ <span className="italic">(idle)</span></div>
      <div className="text-[#8b949e] mt-1">📋 <span className="text-[#e6edf3] font-bold">Unassigned</span> (2)</div>
      <div className="text-[#8b949e]">  🟠 P1 <span className="text-blue-400 font-bold">#042</span>  Add user authentication with JWT</div>
      <div className="text-[#8b949e]">  🟡 P2 <span className="text-blue-400 font-bold">#043</span>  Add rate limiting</div>
      <div className="text-[#8b949e] mt-1">💡 <span className="text-[#e6edf3] font-bold">Next:</span> <span className="text-[rgb(177,185,249)]">/team-claim</span> <span className="text-blue-400 font-bold">#042</span> to start</div>
      {!claimed && (
        <div className="mt-2 pt-2 border-t border-[#30363d]">
          <div className="flex gap-2">
            <span className="text-[#22c55e]">$</span>
            <motion.span
              className="text-[#e6edf3]"
              animate={{ opacity: [1, 0.5, 1] }}
              transition={{ repeat: Infinity, duration: 1.5 }}
            >
              {lang === 'en' ? '▊ type /team-claim #42...' : '▊ 输入 /team-claim #42...'}
            </motion.span>
          </div>
        </div>
      )}
    </div>
  )

  const missionContent = (
    <div className="p-3 font-mono text-xs space-y-0.5">
      <div className="text-[#e6edf3] font-bold">CLAIMED ──── <span className="text-blue-400">#042</span> Add user authentication with JWT</div>
      <div className="text-[#8b949e]">🟠  <span className="text-[rgb(177,185,249)]">P1</span>  Branch: <span className="italic text-[#8b949e]">mission/042-add-auth-bob</span></div>
      <div className="text-[#8b949e]">Contract:  <span className="text-[rgb(177,185,249)]">.teamwork/active/MISSION-42.md</span></div>
      <div className="text-[#8b949e]">Worktree:  <span className="text-[rgb(177,185,249)]">~/Projects/taskflow.worktrees/042-add-auth-bob</span></div>
      <div className="text-[#8b949e] mt-1"><span className="text-[#e6edf3] font-bold">OBJECTIVE</span></div>
      <div className="text-[#8b949e]">  Implement JWT-based authentication for the TaskFlow API</div>
      <div className="text-[#8b949e] mt-1"><span className="text-[#e6edf3] font-bold">SUB-TASKS</span></div>
      <div className="text-[#8b949e]">  [ ] Install jsonwebtoken + bcrypt</div>
      <div className="text-[#8b949e]">  [ ] Create auth service</div>
      <div className="text-[#8b949e]">  [ ] Implement login/register routes</div>
      <div className="text-[#8b949e]">  [ ] Add JWT middleware</div>
      <div className="text-[#8b949e]">  [ ] Write unit tests</div>
      <div className="text-[#8b949e] mt-1">────────────────────────────────────────────</div>
      <div className="text-[#8b949e]"><span className="text-[rgb(177,185,249)]">/team-drive</span> to execute │ <span className="text-[rgb(177,185,249)]">/team-ship</span> when done</div>
    </div>
  )

  const tabContents = claimed
    ? [prelaunchContent, missionContent]
    : [prelaunchContent]

  return (
    <div>
      <div className="flex items-center gap-3 mb-3">
        <p className="text-text-muted text-sm flex-1">
          {lang === 'en'
            ? claimed
              ? '↓ After claiming, a new tab appears. Click tabs to switch:'
              : '↓ Before claiming — only the pre-launch tab. Click "Claim" to see the workflow:'
            : claimed
              ? '↓ 认领后新 tab 出现。点击 tab 切换内容：'
              : '↓ 认领前 — 只有 pre-launch tab。点击"认领"查看工作流：'}
        </p>
        <motion.button
          onClick={() => setClaimed(!claimed)}
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
          className={`px-3 py-1.5 text-xs font-mono rounded border transition-colors ${
            claimed
              ? 'border-text-muted bg-bg-card text-text-secondary'
              : 'border-neon-green bg-neon-green text-white'
          }`}
        >
          {claimed
            ? (lang === 'en' ? '↺ Reset' : '↺ 重置')
            : (lang === 'en' ? '▶ Claim #42' : '▶ 认领 #42')}
        </motion.button>
      </div>
      <ITermTabBar
        tabs={tabs}
        tabContents={tabContents}
        terminalTitle={claimed ? 'mission/042-add-auth-bob (claude)' : 'pre-launch (claude)'}
      />
    </div>
  )
}

export default function Ch5Claim() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[5]
  const desc = descriptions[lang]
  const [activeCase, setActiveCase] = useState(0)

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      <p className="text-text-secondary">
        {lang === 'en'
          ? <><CommandBadge>/team-claim</CommandBadge> transforms a GitHub issue into an executable mission — Contract + branch + worktree, ready for <CommandBadge>/team-drive</CommandBadge>.</>
          : <><CommandBadge>/team-claim</CommandBadge> 将 GitHub Issue 转化为可执行任务 — 合约 + 分支 + worktree，准备好给 <CommandBadge>/team-drive</CommandBadge> 执行。</>}
      </p>

      {/* Interactive worktree claim demo */}
      <WorktreeClaimDemo lang={lang} />

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
        steps={ch5ClaimMocks[activeCase]}
        title="teamwork-claim"
      />
    </ChapterLayout>
  )
}
