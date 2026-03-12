import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import BranchDiagram from '../components/BranchDiagram'
import CommandBadge from '../components/CommandBadge'
import ITermTabBar, { type ITermTab } from '../components/ITermTabBar'
import { motion } from 'framer-motion'

const content = {
  en: {
    intro: 'Teamwork uses a four-branch model that scales from solo to team without configuration overhead.',
    branchTitle: 'The Four Branches',
    branches: [
      { name: 'main', icon: '🟥', color: '#ef4444', rule: 'Production only. Receives squash merges from RC branches. Every commit is a tagged release.' },
      { name: 'pre-launch', icon: '🟦', color: '#3b82f6', rule: 'Integration hub. All mission branches merge here via PR. Base for RC cuts.' },
      { name: 'mission/*', icon: '🟩', color: '#22c55e', rule: 'One branch per issue. Created by /team-claim, named mission/{issue}-{slug}-{user}. Short-lived.' },
      { name: 'rc/*', icon: '🟧', color: '#f59e0b', rule: 'Release candidates. Cut from pre-launch by /team-rc. Frozen for QA, then promoted to main.' },
    ],
    flowTitle: 'The Flow',
    flow: [
      { text: 'Issue created → sits in queue on pre-launch', cmd: '/team-issue' },
      { text: 'Agent claims → creates mission branch from pre-launch', cmd: '/team-claim' },
      { text: 'Work happens on mission branch (each agent gets its own iTerm tab)', cmd: '/team-drive' },
      { text: 'PR from mission branch into pre-launch', cmd: '/team-ship' },
      { text: 'Enough features → cut RC from pre-launch', cmd: '/team-rc' },
      { text: 'QA passes → promote RC to main + tag release', cmd: '/team-rc promote' },
    ],
    tabDemo: 'Each mission runs in its own iTerm2 tab with color-coded branch indicator:',
  },
  zh: {
    intro: 'Teamwork 使用四分支模型，从单人到团队无缝扩展，零配置开销。',
    branchTitle: '四条分支',
    branches: [
      { name: 'main', icon: '🟥', color: '#ef4444', rule: '仅限生产。接收 RC 分支的 squash 合并。每个提交都是带标签的发布。' },
      { name: 'pre-launch', icon: '🟦', color: '#3b82f6', rule: '集成中心。所有 mission 分支通过 PR 合并到此。RC 分支的起点。' },
      { name: 'mission/*', icon: '🟩', color: '#22c55e', rule: '每个 issue 一个分支。由 /team-claim 创建，命名 mission/{编号}-{简称}-{用户}。生命周期短。' },
      { name: 'rc/*', icon: '🟧', color: '#f59e0b', rule: '发布候选。由 /team-rc 从 pre-launch 切出。冻结用于 QA，然后合并到 main。' },
    ],
    flowTitle: '工作流程',
    flow: [
      { text: 'Issue 创建 → 在 pre-launch 队列中等待', cmd: '/team-issue' },
      { text: 'Agent 认领 → 从 pre-launch 创建 mission 分支', cmd: '/team-claim' },
      { text: '在 mission 分支上开发（每个 agent 有独立的 iTerm tab）', cmd: '/team-drive' },
      { text: '从 mission 分支创建 PR 到 pre-launch', cmd: '/team-ship' },
      { text: '足够多功能 → 从 pre-launch 切出 RC', cmd: '/team-rc' },
      { text: 'QA 通过 → 将 RC 合并到 main + 打标签', cmd: '/team-rc promote' },
    ],
    tabDemo: '每个任务在独立的 iTerm2 tab 中运行，带有颜色分支标识：',
  },
}

const demoTabs: ITermTab[] = [
  { label: 'pre-launch (claude)', color: '#3b82f6' },
  { label: 'mission/042-add-auth-bob (claude)', color: '#22c55e', active: true },
  { label: 'mission/043-rate-limiting-charlie (claude)', color: '#22c55e', hasNotification: true },
  { label: 'main (claude)', color: '#ef4444' },
]

export default function Ch1GitModel() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[1]
  const t = content[lang]

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      <p className="text-text-secondary text-lg">{t.intro}</p>

      {/* Interactive branch diagram */}
      <div className="mt-6">
        <BranchDiagram lang={lang} />
      </div>

      {/* Branch rules */}
      <div className="mt-8">
        <h2 className="text-xl font-bold text-text-primary mb-4">{t.branchTitle}</h2>
        <div className="space-y-3">
          {t.branches.map((b, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: -16 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.1, type: 'spring', stiffness: 200 }}
              whileHover={{ x: 4 }}
              className="flex items-start gap-3 p-3 rounded-lg bg-bg-card border border-border"
            >
              <span className="text-sm mt-0.5">{b.icon}</span>
              <div>
                <code className="font-mono text-sm font-bold" style={{ color: b.color }}>{b.name}</code>
                <p className="text-sm text-text-secondary mt-1">{b.rule}</p>
              </div>
            </motion.div>
          ))}
        </div>
      </div>

      {/* iTerm2 tab demo */}
      <div className="mt-8">
        <p className="text-text-secondary text-sm mb-3">{t.tabDemo}</p>
        <ITermTabBar tabs={demoTabs} terminalTitle="mission/042-add-auth-bob (claude)">
          <div className="p-4 font-mono text-sm text-text-secondary">
            <div className="flex gap-2">
              <span className="text-neon-green">$</span>
              <span className="text-text-primary">claude</span>
            </div>
            <div className="mt-1 text-text-muted">
              {lang === 'en'
                ? '# Each agent works in its own tab + worktree'
                : '# 每个 agent 在独立的 tab + worktree 中工作'}
            </div>
          </div>
        </ITermTabBar>
      </div>

      {/* Flow */}
      <div className="mt-8">
        <h2 className="text-xl font-bold text-text-primary mb-4">{t.flowTitle}</h2>
        <div className="space-y-3 pl-4 border-l-2 border-neon-blue/30">
          {t.flow.map((step, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.5 + i * 0.1 }}
              className="flex items-center gap-3 text-sm"
            >
              <motion.span
                className="text-neon-blue font-mono font-bold w-5 shrink-0"
                animate={{ scale: [1, 1.1, 1] }}
                transition={{ delay: 1 + i * 0.2, duration: 0.3 }}
              >
                {i + 1}.
              </motion.span>
              <span className="text-text-secondary flex-1">{step.text}</span>
              <CommandBadge size="sm">{step.cmd}</CommandBadge>
            </motion.div>
          ))}
        </div>
      </div>
    </ChapterLayout>
  )
}
