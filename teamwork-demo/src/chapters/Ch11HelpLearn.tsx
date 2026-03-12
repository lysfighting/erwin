import { useState } from 'react'
import { useLang } from '../i18n/LangContext'
import translations from '../i18n/translations'
import ChapterLayout from '../components/ChapterLayout'
import TerminalPlayground from '../components/TerminalPlayground'
import CommandBadge from '../components/CommandBadge'
import { ch11HelpMocks } from '../data/mockData'
import { motion } from 'framer-motion'

const descriptions = {
  en: [
    {
      heading: '/team help',
      body: 'The complete command reference. Every command, every subcommand, every shortcut — at your fingertips.',
    },
    {
      heading: '/team learn',
      body: 'Interactive concept guide. Understand the philosophy — GitHub-first, Mission Contracts, four-branch model, worktree mode, push notifications.',
    },
  ],
  zh: [
    {
      heading: '/team help',
      body: '完整命令参考。每个命令、每个子命令、每个快捷操作 — 触手可及。',
    },
    {
      heading: '/team learn',
      body: '交互式概念指南。理解核心理念 — GitHub 优先、任务合约、四分支模型、worktree 模式、推送通知。',
    },
  ],
} as const

const commandGroups = {
  en: [
    {
      title: 'Core Workflow',
      cmds: [
        { cmd: '/team', desc: 'Dashboard + init' },
        { cmd: '/team-issue', desc: 'Create issues' },
        { cmd: '/team-claim', desc: 'Claim missions' },
        { cmd: '/team-drive', desc: 'Execute code' },
        { cmd: '/team-ship', desc: 'Deliver PR' },
        { cmd: '/team-rc', desc: 'Release mgmt' },
      ],
    },
    {
      title: 'Shortcuts',
      cmds: [
        { cmd: '/team auto', desc: 'One-shot execution' },
        { cmd: '/team wrap', desc: 'Retroactive tracking' },
        { cmd: '/team #N', desc: 'View issue details' },
      ],
    },
    {
      title: 'Utilities',
      cmds: [
        { cmd: '/team doctor', desc: '12-check health' },
        { cmd: '/team config', desc: 'View/edit config' },
        { cmd: '/team help', desc: 'Command reference' },
        { cmd: '/team learn', desc: 'Concepts guide' },
      ],
    },
  ],
  zh: [
    {
      title: '核心工作流',
      cmds: [
        { cmd: '/team', desc: '仪表盘 + 初始化' },
        { cmd: '/team-issue', desc: '创建任务' },
        { cmd: '/team-claim', desc: '认领任务' },
        { cmd: '/team-drive', desc: '执行代码' },
        { cmd: '/team-ship', desc: '交付 PR' },
        { cmd: '/team-rc', desc: '发布管理' },
      ],
    },
    {
      title: '快捷操作',
      cmds: [
        { cmd: '/team auto', desc: '一键执行' },
        { cmd: '/team wrap', desc: '追溯跟踪' },
        { cmd: '/team #N', desc: '查看任务详情' },
      ],
    },
    {
      title: '工具',
      cmds: [
        { cmd: '/team doctor', desc: '12 项健康检查' },
        { cmd: '/team config', desc: '查看/编辑配置' },
        { cmd: '/team help', desc: '命令参考' },
        { cmd: '/team learn', desc: '概念指南' },
      ],
    },
  ],
}

export default function Ch11HelpLearn() {
  const { lang } = useLang()
  const ch = translations[lang].chapters[11]
  const desc = descriptions[lang]
  const groups = commandGroups[lang]
  const [activeCase, setActiveCase] = useState(0)

  return (
    <ChapterLayout title={ch.title} subtitle={ch.subtitle}>
      <p className="text-text-secondary">
        {lang === 'en'
          ? <>New to Teamwork? <CommandBadge>/team help</CommandBadge> gives you the full command reference. <CommandBadge>/team learn</CommandBadge> explains the concepts behind each command.</>
          : <>刚接触 Teamwork？<CommandBadge>/team help</CommandBadge> 提供完整命令参考。<CommandBadge>/team learn</CommandBadge> 解释每个命令背后的理念。</>}
      </p>

      {/* Command map */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {groups.map((g, gi) => (
          <motion.div
            key={gi}
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: gi * 0.1, type: 'spring', stiffness: 200 }}
            className="p-4 rounded-lg border border-border bg-bg-card"
          >
            <h3 className="text-sm font-bold text-text-primary mb-3">{g.title}</h3>
            <div className="space-y-2">
              {g.cmds.map((c, ci) => (
                <div key={ci} className="flex items-center gap-2">
                  <CommandBadge size="sm">{c.cmd}</CommandBadge>
                  <span className="text-text-muted text-xs">{c.desc}</span>
                </div>
              ))}
            </div>
          </motion.div>
        ))}
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
        steps={ch11HelpMocks[activeCase]}
        title="teamwork-help"
      />

      {/* Design Principles from learn docs */}
      <div className="space-y-3">
        <h2 className="text-lg font-bold text-text-primary">
          {lang === 'en' ? 'Design Principles' : '设计原则'}
        </h2>
        <p className="text-xs text-text-muted italic">
          {lang === 'en'
            ? '"These are not rules we invented. They emerged from pain — bugs found, rework done, production incidents caused by AI-generated noise."'
            : '"这些不是我们发明的规则。它们源于痛苦 — 发现的 bug、重做的工作、AI 生成噪音导致的生产事故。"'}
        </p>
        <div className="space-y-1.5">
          {(lang === 'en' ? [
            { num: '1', title: 'GitHub IS the system', desc: 'Issues = tasks. PRs = delivery. Actions = gates. Not a mirror.' },
            { num: '2', title: 'Separate dirty from clean', desc: 'develop absorbs AI noise. main stays pristine.' },
            { num: '3', title: 'Squash everything', desc: '1 Issue = 1 commit. Cherry-pickable, bisectable, readable.' },
            { num: '4', title: 'Risk levels don\'t share paths', desc: 'Staging deploy ≠ prod deploy. Never the same command.' },
            { num: '5', title: 'Tags are semantic markers', desc: 'V0.1.0 means "this version exists." Not "deploy this now."' },
            { num: '6', title: 'Scripts lock, LLM judges', desc: '21 bugs in LLM ops, 0 in scripts. Content changes; mechanism endures.' },
            { num: '7', title: 'Close the loop first', desc: 'A complete loop with rough edges > a polished partial flow.' },
            { num: '8', title: 'Config-driven', desc: 'Every project-specific value in config.yml. Not hardcoded.' },
            { num: '9', title: '6 skills is the ceiling', desc: '30+ commands confused users. Expand via subcommands, not new skills.' },
          ] : [
            { num: '1', title: 'GitHub 就是系统', desc: 'Issues = 任务。PRs = 交付。Actions = 门禁。不是镜像。' },
            { num: '2', title: '分离脏和干净', desc: 'develop 吸收 AI 噪音。main 保持干净。' },
            { num: '3', title: 'Squash 一切', desc: '1 Issue = 1 提交。可 cherry-pick、可 bisect、可读。' },
            { num: '4', title: '风险等级不共享路径', desc: '暂存部署 ≠ 生产部署。永远不是同一条命令。' },
            { num: '5', title: '标签是语义标记', desc: 'V0.1.0 意味着"此版本存在"，不是"现在部署这个"。' },
            { num: '6', title: '脚本锁定，LLM 判断', desc: 'LLM 操作 21 个 bug，脚本 0 个。内容变化；机制永恒。' },
            { num: '7', title: '先闭环', desc: '粗糙但完整的循环 > 精致但部分的流程。' },
            { num: '8', title: '配置驱动', desc: '每个项目特有的值都在 config.yml 中。不硬编码。' },
            { num: '9', title: '6 个技能是上限', desc: '30+ 命令让用户困惑。通过子命令扩展，不增加新技能。' },
          ]).map((p, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.05 }}
              whileHover={{ x: 4 }}
              className="flex items-start gap-3 p-2 rounded-lg bg-bg-card border border-border"
            >
              <span className="text-xs font-mono text-neon-blue font-bold w-4 shrink-0 mt-0.5">{p.num}</span>
              <div>
                <span className="text-xs font-medium text-text-primary">{p.title}</span>
                <span className="text-[10px] text-text-muted ml-1.5">{p.desc}</span>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </ChapterLayout>
  )
}
