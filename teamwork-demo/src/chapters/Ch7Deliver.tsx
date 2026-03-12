import { useState } from 'react'
import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import TerminalPlayground from '../components/TerminalPlayground'
import CommandBadge from '../components/CommandBadge'
import { GitHubPRMock } from '../components/GitHubMock'
import { ch7DeliverMocks } from '../data/mockData'
import { motion } from 'framer-motion'

const descriptions = {
  en: [
    {
      heading: 'Ship a PR',
      body: '/team-ship pushes your branch, creates a PR with AI-generated description, waits for CI, and requests review. One command from code to PR.',
    },
    {
      heading: 'AI code review',
      body: '/team-ship review runs AI-powered review on your PR — security, best practices, and improvements before human review.',
    },
    {
      heading: 'Post-merge cleanup',
      body: '/team-ship done handles everything after merge: delete branch, close issue, archive contract, update milestone.',
    },
  ],
  zh: [
    {
      heading: '发起 PR',
      body: '/team-ship 推送分支、创建 AI 生成描述的 PR、等待 CI、请求审查。一条命令从代码到 PR。',
    },
    {
      heading: 'AI 代码审查',
      body: '/team-ship review 在人工审查前进行 AI 代码审查 — 安全性、最佳实践和改进建议。',
    },
    {
      heading: '合并后清理',
      body: '/team-ship done 处理合并后的一切：删除分支、关闭 Issue、归档合约、更新里程碑。',
    },
  ],
} as const

export default function Ch7Deliver() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[7]
  const desc = descriptions[lang]
  const [activeCase, setActiveCase] = useState(0)

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      <p className="text-text-secondary">
        {lang === 'en'
          ? <><CommandBadge>/team-ship</CommandBadge> delivers your work as a GitHub PR — with AI-generated description, CI checks, and review requests.</>
          : <><CommandBadge>/team-ship</CommandBadge> 将你的工作以 GitHub PR 的形式交付 — 包含 AI 生成的描述、CI 检查和审查请求。</>}
      </p>

      {/* Ship automation pipeline */}
      <div className="rounded-lg border border-border bg-bg-dark p-4 overflow-x-auto">
        <div className="flex items-center gap-2 min-w-[500px]">
          {(lang === 'en' ? [
            { icon: '🧪', label: 'Pre-flight', detail: 'Tests pass?' },
            { icon: '⬆️', label: 'Push', detail: 'Branch → remote' },
            { icon: '📝', label: 'PR', detail: 'AI description' },
            { icon: '🔄', label: 'CI', detail: 'Watch status' },
            { icon: '👁️', label: 'Review', detail: 'Request reviewer' },
            { icon: '🏷️', label: 'Labels', detail: 'status:review' },
          ] : [
            { icon: '🧪', label: '预检', detail: '测试通过？' },
            { icon: '⬆️', label: '推送', detail: '分支 → 远端' },
            { icon: '📝', label: 'PR', detail: 'AI 生成描述' },
            { icon: '🔄', label: 'CI', detail: '监控状态' },
            { icon: '👁️', label: '审查', detail: '请求审查人' },
            { icon: '🏷️', label: '标签', detail: 'status:review' },
          ]).map((step, i) => (
            <motion.div key={i} className="flex items-center gap-2 flex-1">
              <motion.div
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: i * 0.1, type: 'spring' }}
                className="flex flex-col items-center gap-1 flex-1"
              >
                <span className="text-base">{step.icon}</span>
                <span className="text-[10px] font-mono text-text-primary font-medium">{step.label}</span>
                <span className="text-[9px] text-text-muted">{step.detail}</span>
              </motion.div>
              {i < 5 && (
                <motion.span
                  animate={{ x: [0, 3, 0] }}
                  transition={{ repeat: Infinity, duration: 1.5, delay: i * 0.15 }}
                  className="text-neon-green text-xs shrink-0"
                >→</motion.span>
              )}
            </motion.div>
          ))}
        </div>
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
        steps={ch7DeliverMocks[activeCase]}
        title="teamwork-ship"
      />

      {/* GitHub PR mock — shows what the PR looks like */}
      {activeCase === 0 && (
        <motion.div
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <p className="text-text-muted text-sm mb-2">
            {lang === 'en' ? '↓ The PR created on GitHub:' : '↓ 在 GitHub 上创建的 PR：'}
          </p>
          <GitHubPRMock
            number={87}
            title="feat(auth): Add user authentication with JWT"
            state="open"
            author="bob"
            baseBranch="pre-launch"
            headBranch="mission/042-add-auth-bob"
            labels={[
              { name: 'mission', color: '0075ca' },
              { name: 'status:review', color: 'fbca04' },
              { name: 'priority:P1', color: 'e11d48' },
              { name: 'size:M', color: '0ea5e9' },
            ]}
            reviewers={['alice']}
            checks={[
              { name: 'lint', status: 'pass' },
              { name: 'typecheck', status: 'pass' },
              { name: 'unit-tests (47 passed)', status: 'pass' },
              { name: 'integration (12 passed)', status: 'pass' },
              { name: 'coverage (87%)', status: 'pass' },
            ]}
            additions={487}
            deletions={12}
            files={[
              'src/services/auth.ts',
              'src/routes/auth.ts',
              'src/middleware/auth.ts',
              'tests/auth.test.ts',
              'package.json',
              'package-lock.json',
            ]}
            body={`## Summary
Implements JWT-based authentication for the TaskFlow API.

Closes #42

## Changes
- **Auth service**: Token generation, verification, refresh with rotation
- **Routes**: POST /auth/login and POST /auth/register
- **Middleware**: JWT validation on protected routes (401/403)
- **Tests**: 12 test cases, 87% coverage

## Test Plan
- \`npm test -- --grep "auth"\`
- Manual: POST /auth/register → POST /auth/login → access protected route

🤖 Generated with Teamwork /team-ship`}
          />
        </motion.div>
      )}
    </ChapterLayout>
  )
}
