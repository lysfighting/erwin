import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import StateMachine from '../components/StateMachine'
import CommandBadge from '../components/CommandBadge'
import { type ReactNode } from 'react'
import { motion } from 'framer-motion'

interface StateItem {
  name: string
  desc: ReactNode
}

const text = {
  en: {
    intro: 'Every issue in Teamwork follows a clear lifecycle — from a raw idea to a production release. Each transition is triggered by a specific command or event.',
    statesTitle: 'States',
    states: [
      { name: 'Idea', desc: 'A raw thought or requirement. Not yet tracked in GitHub.' },
      { name: 'Queued', desc: <>Published as a GitHub Issue via <CommandBadge size="sm">/team-issue</CommandBadge>. Visible to the whole team.</> },
      { name: 'Claimed', desc: <>An agent runs <CommandBadge size="sm">/team-claim</CommandBadge> — a Contract is generated, a branch is created.</> },
      { name: 'In Progress', desc: <><CommandBadge size="sm">/team-drive</CommandBadge> starts execution. Code is being written against the Contract.</> },
      { name: 'In Review', desc: <><CommandBadge size="sm">/team-ship</CommandBadge> creates a PR. The code is ready for review.</> },
      { name: 'Done', desc: 'PR is merged into pre-launch. Feature is integrated.' },
      { name: 'Released', desc: <><CommandBadge size="sm">/team-rc promote</CommandBadge> pushes the release to production.</> },
    ] as StateItem[],
    backwardTitle: 'Backward Transitions',
    backward: 'If a PR review requests changes, the issue moves back from "In Review" to "In Progress". The agent picks up the feedback and continues driving.',
    clickHint: 'Click any node to jump to that state. Auto-plays every 2 seconds.',
  },
  zh: {
    intro: 'Teamwork 中的每个 Issue 都遵循清晰的生命周期——从原始想法到生产发布。每次状态转换都由特定命令或事件触发。',
    statesTitle: '状态说明',
    states: [
      { name: 'Idea', desc: '原始想法或需求。尚未在 GitHub 中跟踪。' },
      { name: 'Queued', desc: <>通过 <CommandBadge size="sm">/team-issue</CommandBadge> 发布为 GitHub Issue，整个团队可见。</> },
      { name: 'Claimed', desc: <>Agent 运行 <CommandBadge size="sm">/team-claim</CommandBadge> — 生成合约，创建分支。</> },
      { name: 'In Progress', desc: <><CommandBadge size="sm">/team-drive</CommandBadge> 开始执行，按合约编写代码。</> },
      { name: 'In Review', desc: <><CommandBadge size="sm">/team-ship</CommandBadge> 创建 PR，代码进入评审阶段。</> },
      { name: 'Done', desc: 'PR 合并到 pre-launch 分支，功能已集成。' },
      { name: 'Released', desc: <><CommandBadge size="sm">/team-rc promote</CommandBadge> 将发布推送到生产环境。</> },
    ] as StateItem[],
    backwardTitle: '反向转换',
    backward: '如果 PR 评审要求修改，Issue 会从「In Review」退回「In Progress」。Agent 接收反馈并继续推进。',
    clickHint: '点击任意节点跳转状态。每 2 秒自动播放。',
  },
}

export default function Ch2Lifecycle() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[2]
  const t = text[lang]

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      <p className="text-text-secondary leading-relaxed">{t.intro}</p>

      <div className="bg-bg-card border border-border rounded-lg p-4">
        <StateMachine autoPlay intervalMs={2000} />
        <p className="text-text-muted text-xs text-center mt-1">{t.clickHint}</p>
      </div>

      <div>
        <h2 className="text-xl font-semibold text-text-primary mb-3">{t.statesTitle}</h2>
        <div className="space-y-2">
          {t.states.map(s => (
            <div key={s.name} className="flex gap-3 items-baseline">
              <code className="text-neon-blue text-sm shrink-0">{s.name}</code>
              <span className="text-text-secondary text-sm">{s.desc}</span>
            </div>
          ))}
        </div>
      </div>

      <div>
        <h2 className="text-xl font-semibold text-text-primary mb-3">{t.backwardTitle}</h2>
        <p className="text-text-secondary leading-relaxed">{t.backward}</p>
      </div>

      {/* Issue management value props */}
      <div className="space-y-3">
        <h2 className="text-xl font-semibold text-text-primary">
          {lang === 'en' ? 'What Issue Lifecycle Automates' : 'Issue 生命周期自动化了什么'}
        </h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
          {(lang === 'en' ? [
            { icon: '🏷️', title: 'Label transitions', desc: 'queued → wip → review → done — automatic' },
            { icon: '📋', title: 'Milestone tracking', desc: 'Progress bar updates as issues close' },
            { icon: '🔗', title: 'PR auto-close', desc: '"Closes #42" in PR → Issue closes on merge' },
            { icon: '✅', title: 'Checkbox completion', desc: 'Post-merge Action checks all boxes' },
            { icon: '🔍', title: 'Duplicate detection', desc: 'AI scans existing issues before creating' },
            { icon: '📊', title: 'Team visibility', desc: 'Dashboard shows who does what, real-time' },
          ] : [
            { icon: '🏷️', title: '标签自动转换', desc: 'queued → wip → review → done — 自动完成' },
            { icon: '📋', title: '里程碑追踪', desc: '随 Issue 关闭自动更新进度条' },
            { icon: '🔗', title: 'PR 自动关闭', desc: 'PR 中的 "Closes #42" → 合并时自动关闭 Issue' },
            { icon: '✅', title: '复选框完成', desc: '合并后 Action 自动勾选所有选项' },
            { icon: '🔍', title: '重复检测', desc: 'AI 在创建前扫描已有 Issue' },
            { icon: '📊', title: '团队可见性', desc: '仪表盘实时展示谁在做什么' },
          ]).map((item, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.06 }}
              whileHover={{ x: 4 }}
              className="flex items-start gap-2.5 p-2.5 rounded-lg bg-bg-card border border-border"
            >
              <span className="text-base mt-0.5">{item.icon}</span>
              <div>
                <div className="text-xs font-medium text-text-primary">{item.title}</div>
                <div className="text-[10px] text-text-muted">{item.desc}</div>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </ChapterLayout>
  )
}
