import { motion, AnimatePresence } from 'framer-motion'
import { SectionWrapper } from '../components/SectionWrapper'
import { useState } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'
import { Annotatable } from '../components/Annotatable'

interface Trigger {
  id: number
  title: string
  titleZh: string
  desc: string
  descZh: string
}

const triggers: Trigger[] = [
  {
    id: 1,
    title: 'Before Phase 1',
    titleZh: '阶段一之前',
    desc: 'Always validate the overall approach before creating tasks.',
    descZh: '在创建任务之前，始终验证整体方向。',
  },
  {
    id: 2,
    title: '"Am I building the right thing?"',
    titleZh: '"我在构建正确的东西吗？"',
    desc: 'At any checkpoint where this question surfaces — stop and validate.',
    descZh: '在任何出现这个疑问的检查点——停下来验证。',
  },
  {
    id: 3,
    title: 'Unexpected complexity',
    titleZh: '意外复杂度',
    desc: 'Complexity often signals wrong abstraction. When a task reveals unexpected complexity, validate direction.',
    descZh: '复杂度通常意味着错误的抽象。当任务暴露出意外复杂度时，验证方向。',
  },
  {
    id: 4,
    title: 'Working around something',
    titleZh: '绕行某些东西',
    desc: 'Workarounds suggest wrong framing. If you catch yourself working around something, the framing may be off.',
    descZh: '变通方案暗示框架有误。如果你发现自己在绕行某些东西，框架可能有问题。',
  },
  {
    id: 5,
    title: '2+ viable paths, no clear winner',
    titleZh: '2+ 可行路径，无明确赢家',
    desc: 'When multiple viable paths exist and you cannot reason to a clear winner, validate before committing.',
    descZh: '当存在多条可行路径且无法推理出明确赢家时，先验证再决定。',
  },
  {
    id: 6,
    title: 'You feel uncertain',
    titleZh: '你感到不确定',
    desc: 'Trust that instinct. It is cheaper to validate now than to discover you were wrong later.',
    descZh: '相信你的直觉。现在验证的成本远低于日后发现错误。',
  },
]

interface Step {
  id: number
  name: string
  nameZh: string
  desc: string
  descZh: string
  icon: string
}

const validationSteps: Step[] = [
  {
    id: 1,
    name: 'Visualize the Big Picture',
    nameZh: '可视化全局',
    desc: 'Mermaid/ASCII diagram of current state and heading. See where you are and where you are going.',
    descZh: '用 Mermaid/ASCII 图展示当前状态和前进方向。看清你在哪里以及要去哪里。',
    icon: '🗺️',
  },
  {
    id: 2,
    name: 'Exhaustive Path Enumeration',
    nameZh: '穷举路径列举',
    desc: 'ALL distinct paths. Name each vividly (not "Option A" but "Fortress Architecture"). Key advantage + key risk + fatal flaw scenario for each.',
    descZh: '列举所有不同路径。给每条路径起生动的名字（不是"选项A"而是"堡垒架构"）。每条路径的关键优势 + 关键风险 + 致命缺陷场景。',
    icon: '🔀',
  },
  {
    id: 3,
    name: 'Adversarial Stress Test',
    nameZh: '对抗性压力测试',
    desc: '"What\'s the strongest argument AGAINST this path?" / "What would make this catastrophically wrong?"',
    descZh: '"反对这条路径的最强论据是什么？" / "什么情况会导致灾难性错误？"',
    icon: '⚔️',
  },
  {
    id: 4,
    name: 'Confirm with User (when uncertain)',
    nameZh: '与用户确认（不确定时）',
    desc: 'AskUserQuestion with VISUAL markdown previews showing architecture diagrams. Don\'t describe — show.',
    descZh: '使用 AskUserQuestion 并附带可视化 markdown 预览展示架构图。不要描述——展示。',
    icon: '🤝',
  },
  {
    id: 5,
    name: 'Lock and Proceed',
    nameZh: '锁定并执行',
    desc: 'Once validated, execute with full conviction. No second-guessing. The direction is set.',
    descZh: '一旦验证通过，全力执行。不再犹豫。方向已定。',
    icon: '🔒',
  },
]

interface Heuristic {
  signal: string
  signalZh: string
  meaning: string
  meaningZh: string
}

const heuristics: Heuristic[] = [
  {
    signal: 'Keep needing workarounds',
    signalZh: '不断需要变通方案',
    meaning: 'Wrong abstraction',
    meaningZh: '错误的抽象',
  },
  {
    signal: 'Code getting more complex than expected',
    signalZh: '代码比预期更复杂',
    meaning: 'Wrong decomposition',
    meaningZh: '错误的分解',
  },
  {
    signal: 'Fighting the framework/library',
    signalZh: '与框架/库对抗',
    meaning: 'Wrong tool choice',
    meaningZh: '错误的工具选择',
  },
  {
    signal: 'Success criteria feel disconnected',
    signalZh: '成功标准感觉脱节',
    meaning: 'Scope drift',
    meaningZh: '范围偏移',
  },
  {
    signal: "Can't explain in one sentence why this task serves the mission",
    signalZh: '无法一句话解释此任务如何服务于使命',
    meaning: 'Lost the thread',
    meaningZh: '失去了主线',
  },
]

export function DirectionValidationSection() {
  const { lang } = useLang()
  const [expandedStep, setExpandedStep] = useState<number | null>(null)
  const [expandedTrigger, setExpandedTrigger] = useState<number | null>(null)

  return (
    <SectionWrapper
      id="direction-validation"
      title={t('directionvalidation.title', lang)}
      subtitle={t('directionvalidation.subtitle', lang)}
    >
      {/* WHY IT EXISTS block */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        className="mb-10 bg-gold/5 border border-gold/20 rounded-lg p-6"
      >
        <div className="flex items-start gap-3">
          <span className="text-2xl flex-shrink-0 mt-0.5">🧭</span>
          <div>
            <p className="text-sm font-mono text-gold font-bold mb-2">
              {lang === 'zh' ? '为什么存在' : 'WHY IT EXISTS'}
            </p>
            <p className="text-text-secondary text-sm leading-relaxed mb-3">
              {lang === 'zh'
                ? '方向错误全速前进 = 最大浪费。正确方向上的慢行胜过错误路径上的冲刺。'
                : 'Wrong direction at full speed = maximum waste. A slow walk in the right direction beats a sprint down the wrong path.'}
            </p>
            <div className="bg-bg-dark/60 border border-border/50 rounded px-4 py-3">
              <p className="text-gold font-mono text-sm text-center">
                "方向错误的努力，本质上是在浪费时间"
              </p>
              <p className="text-text-muted text-xs text-center mt-1 italic">
                {lang === 'zh'
                  ? '方向公理'
                  : 'Direction Axiom — Effort in wrong direction is fundamentally wasted time'}
              </p>
            </div>
          </div>
        </div>
      </motion.div>

      {/* 6 Triggers */}
      <Annotatable
        id="direction-triggers"
        label="6 Triggers for Direction Validation"
        sectionId="direction-validation"
        sectionLabel="Direction Validation"
      >
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mb-10"
        >
          <h3 className="text-xl font-bold text-text-primary mb-2 font-mono">
            {lang === 'zh' ? '6 个触发时机' : '6 Triggers'}
          </h3>
          <p className="text-xs text-text-muted mb-4 font-mono">
            {lang === 'zh' ? '何时启动方向验证' : 'When to invoke direction validation'}
          </p>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-3">
            {triggers.map((trigger, i) => (
              <motion.div
                key={trigger.id}
                initial={{ opacity: 0, y: 16 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.05 }}
                onClick={() => setExpandedTrigger(expandedTrigger === trigger.id ? null : trigger.id)}
                className={`bg-bg-card border rounded-lg p-4 cursor-pointer transition-all ${
                  expandedTrigger === trigger.id
                    ? 'border-gold/40 bg-bg-card-hover'
                    : 'border-border hover:border-gold/20'
                }`}
              >
                <div className="flex items-start gap-3">
                  <span className="text-gold font-bold text-lg font-mono flex-shrink-0">
                    {trigger.id}
                  </span>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-text-primary font-medium leading-snug">
                      {lang === 'zh' ? trigger.titleZh : trigger.title}
                    </p>
                    <AnimatePresence>
                      {expandedTrigger === trigger.id && (
                        <motion.div
                          initial={{ opacity: 0, height: 0 }}
                          animate={{ opacity: 1, height: 'auto' }}
                          exit={{ opacity: 0, height: 0 }}
                          className="mt-2"
                        >
                          <p className="text-xs text-text-secondary leading-relaxed">
                            {lang === 'zh' ? trigger.descZh : trigger.desc}
                          </p>
                        </motion.div>
                      )}
                    </AnimatePresence>
                  </div>
                </div>
              </motion.div>
            ))}
          </div>
        </motion.div>
      </Annotatable>

      {/* 5-Step Validation Process */}
      <Annotatable
        id="direction-5-steps"
        label="5-Step Validation Process"
        sectionId="direction-validation"
        sectionLabel="Direction Validation"
      >
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mb-10"
        >
          <h3 className="text-xl font-bold text-text-primary mb-2 font-mono">
            {lang === 'zh' ? '5 步验证流程' : '5-Step Validation Process'}
          </h3>
          <p className="text-xs text-text-muted mb-4 font-mono">
            {lang === 'zh' ? '从可视化到锁定执行' : 'From visualization to locked execution'}
          </p>

          {/* Vertical step flow */}
          <div className="space-y-3">
            {validationSteps.map((step, i) => (
              <motion.div
                key={step.id}
                initial={{ opacity: 0, x: -20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.08 }}
                onClick={() => setExpandedStep(expandedStep === step.id ? null : step.id)}
                className={`bg-bg-card border rounded-lg p-4 cursor-pointer transition-all ${
                  expandedStep === step.id
                    ? 'border-gold/40 bg-bg-card-hover'
                    : 'border-border hover:border-gold/20'
                }`}
              >
                <div className="flex items-start gap-4">
                  {/* Step number + connector */}
                  <div className="flex flex-col items-center flex-shrink-0">
                    <div className="w-10 h-10 rounded-full bg-gold/10 border border-gold/30 flex items-center justify-center">
                      <span className="text-gold font-bold font-mono text-sm">{step.id}</span>
                    </div>
                    {i < validationSteps.length - 1 && (
                      <div className="w-px h-4 bg-border/50 mt-1" />
                    )}
                  </div>

                  <div className="flex-1 min-w-0 pt-1">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-lg">{step.icon}</span>
                      <p className="text-sm text-text-primary font-bold">
                        {lang === 'zh' ? step.nameZh : step.name}
                      </p>
                    </div>

                    <AnimatePresence>
                      {expandedStep === step.id && (
                        <motion.div
                          initial={{ opacity: 0, height: 0 }}
                          animate={{ opacity: 1, height: 'auto' }}
                          exit={{ opacity: 0, height: 0 }}
                          className="mt-2"
                        >
                          <div className="bg-gold/5 border border-gold/10 rounded p-3">
                            <p className="text-xs text-text-secondary leading-relaxed">
                              {lang === 'zh' ? step.descZh : step.desc}
                            </p>
                          </div>
                        </motion.div>
                      )}
                    </AnimatePresence>
                  </div>
                </div>
              </motion.div>
            ))}
          </div>
        </motion.div>
      </Annotatable>

      {/* Direction Error vs Execution Error — prominent comparison block */}
      <Annotatable
        id="direction-vs-execution"
        label="Direction vs Execution Error"
        sectionId="direction-validation"
        sectionLabel="Direction Validation"
      >
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mb-10"
        >
          <h3 className="text-xl font-bold text-text-primary mb-4 font-mono">
            {lang === 'zh' ? '方向错误 vs 执行错误' : 'Direction Error vs Execution Error'}
          </h3>

          <div className="grid md:grid-cols-2 gap-4">
            {/* Execution Error */}
            <div className="bg-bg-card border border-neon-green/20 rounded-lg p-5">
              <div className="flex items-center gap-2 mb-3">
                <div className="w-3 h-3 rounded-full bg-neon-green/60" />
                <h4 className="text-sm font-bold font-mono text-neon-green">
                  {lang === 'zh' ? '执行错误' : 'Execution Error'}
                </h4>
              </div>
              <div className="bg-neon-green/5 border border-neon-green/10 rounded p-3 mb-3">
                <p className="text-xs text-text-secondary italic">
                  {lang === 'zh'
                    ? '"我在构建正确的东西，但这个函数有 bug"'
                    : '"I\'m building the right thing but this function has a bug"'}
                </p>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-neon-green text-xs font-mono font-bold px-2 py-1 bg-neon-green/10 rounded">
                  {lang === 'zh' ? '行动' : 'ACTION'}
                </span>
                <p className="text-xs text-text-secondary">
                  {lang === 'zh'
                    ? '修复并继续。循环会处理这个。'
                    : 'Fix and continue. The loop handles this.'}
                </p>
              </div>
            </div>

            {/* Direction Error */}
            <div className="bg-bg-card border border-neon-red/20 rounded-lg p-5">
              <div className="flex items-center gap-2 mb-3">
                <div className="w-3 h-3 rounded-full bg-neon-red/60" />
                <h4 className="text-sm font-bold font-mono text-neon-red">
                  {lang === 'zh' ? '方向错误' : 'Direction Error'}
                </h4>
              </div>
              <div className="bg-neon-red/5 border border-neon-red/10 rounded p-3 mb-3">
                <p className="text-xs text-text-secondary italic">
                  {lang === 'zh'
                    ? '"我在构建完全错误的东西"'
                    : '"I\'m building the wrong thing entirely"'}
                </p>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-neon-red text-xs font-mono font-bold px-2 py-1 bg-neon-red/10 rounded">
                  {lang === 'zh' ? '行动' : 'ACTION'}
                </span>
                <p className="text-xs text-text-secondary">
                  {lang === 'zh'
                    ? '停止。启动方向验证。不及早发现 = 浪费全部剩余预算。'
                    : 'STOP. Invoke Direction Validation. Not caught early = waste the entire remaining budget.'}
                </p>
              </div>
            </div>
          </div>
        </motion.div>
      </Annotatable>

      {/* Heuristics for detecting direction errors */}
      <Annotatable
        id="direction-heuristics"
        label="Heuristics for Detecting Direction Errors"
        sectionId="direction-validation"
        sectionLabel="Direction Validation"
      >
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="bg-bg-card border border-border rounded-lg p-6"
        >
          <h3 className="text-sm font-bold font-mono text-text-primary mb-4">
            {lang === 'zh' ? '方向错误检测启发式' : 'Heuristics for Detecting Direction Errors Mid-Execution'}
          </h3>
          <div className="space-y-2">
            {heuristics.map((h, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0, x: -12 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.05 }}
                className="flex items-center gap-3 py-2 border-b border-border/30 last:border-b-0"
              >
                <span className="text-neon-red text-xs flex-shrink-0">&#x25B6;</span>
                <p className="text-xs text-text-secondary flex-1">
                  {lang === 'zh' ? h.signalZh : h.signal}
                </p>
                <span className="text-xs font-mono text-gold px-2 py-0.5 bg-gold/10 border border-gold/20 rounded flex-shrink-0">
                  {lang === 'zh' ? h.meaningZh : h.meaning}
                </span>
              </motion.div>
            ))}
          </div>
        </motion.div>
      </Annotatable>
    </SectionWrapper>
  )
}
