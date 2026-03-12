import { useState } from 'react'
import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import TerminalPlayground from '../components/TerminalPlayground'
import CommandBadge from '../components/CommandBadge'
import { GitHubIssueMock } from '../components/GitHubMock'
import { ch4CreateMocks } from '../data/mockData'
import { motion } from 'framer-motion'

const descriptions = {
  en: [
    {
      heading: 'Create issue',
      body: 'Describe what you need in plain English. AI analyzes, generates labels, acceptance criteria, sub-tasks, and publishes a fully enriched GitHub issue.',
    },
    {
      heading: 'Batch create',
      body: 'Create an entire milestone\'s worth of issues in a single command. Each with AI-generated structure and automatic prioritization.',
    },
    {
      heading: 'View issue detail',
      body: 'Check any issue\'s live status with /team #N — see sub-task progress, assignee, branch, and contract location.',
    },
  ],
  zh: [
    {
      heading: '创建单个任务',
      body: '用自然语言描述需求。AI 分析后自动生成标签、验收标准、子任务，并发布一个结构完整的 GitHub Issue。',
    },
    {
      heading: '批量创建',
      body: '一条命令创建整个里程碑的所有 Issue，每个都有 AI 生成的结构和自动优先级。',
    },
    {
      heading: '查看任务详情',
      body: '用 /team #N 查看任务实时状态 — 子任务进度、负责人、分支、合约位置。',
    },
  ],
} as const

export default function Ch4Create() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[4]
  const desc = descriptions[lang]
  const [activeCase, setActiveCase] = useState(0)

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      <p className="text-text-secondary">
        {lang === 'en'
          ? <><CommandBadge>/team-issue</CommandBadge> creates AI-enriched GitHub issues with contracts, labels, and acceptance criteria — ready for agents to claim.</>
          : <><CommandBadge>/team-issue</CommandBadge> 创建 AI 增强的 GitHub Issue，包含合约、标签和验收标准 — 等待 agent 认领。</>}
      </p>

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

      {/* Terminal */}
      <TerminalPlayground
        key={`term-${activeCase}`}
        steps={ch4CreateMocks[activeCase]}
        title="teamwork-issue"
      />

      {/* GitHub Issue mock — shows what the created issue looks like on GitHub */}
      {activeCase === 0 && (
        <motion.div
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <p className="text-text-muted text-sm mb-2">
            {lang === 'en' ? '↓ How it looks on GitHub:' : '↓ 在 GitHub 上的效果：'}
          </p>
          <GitHubIssueMock
            number={42}
            title="Add user authentication with JWT"
            state="open"
            author="bob"
            assignees={['bob']}
            labels={[
              { name: 'mission', color: '0075ca' },
              { name: 'status:queued', color: 'e4e669' },
              { name: 'priority:P1', color: 'e11d48' },
              { name: 'domain:backend', color: '7c3aed' },
              { name: 'size:M', color: '0ea5e9' },
            ]}
            milestone="V1.0"
            body={`## Mission Contract

**Objective:** Implement JWT-based authentication for the TaskFlow API.

**Sub-tasks:**
- [ ] Install jsonwebtoken + bcrypt dependencies
- [ ] Create auth service (generate, verify, refresh tokens)
- [ ] Implement POST /auth/login and POST /auth/register
- [ ] Add JWT validation middleware to protected routes
- [ ] Write unit tests (≥80% coverage for auth module)

**Acceptance Criteria:**
- POST /auth/login returns valid JWT with 15m expiry
- POST /auth/register hashes password with bcrypt
- Protected routes return 401 without valid token
- Token refresh endpoint with rotation

**Context Files:** src/routes/index.ts · src/middleware/ · package.json
**Test Command:** npm test -- --grep "auth"`}
          />
        </motion.div>
      )}
    </ChapterLayout>
  )
}
