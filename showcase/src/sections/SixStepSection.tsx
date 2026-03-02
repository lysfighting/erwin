import { motion, useInView } from 'framer-motion'
import { useRef, useState } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

const steps = [
  {
    id: 1,
    name: 'Research',
    chinese: '研究',
    desc: 'Deep reading, externalized understanding. Launch a Research Swarm — parallel scouts explore the codebase from multiple dimensions.',
    file: 'research-{slug}.md',
    color: '#06b6d4',
    detail: 'Not a summary — a structured analysis with specific file:line references. Record surprises — anything that contradicts assumptions.',
    without: 'You build what you THINK the user wants. The gap between "what was asked" and "what was meant" only surfaces after hours of implementation. Rework is the tax on skipping research.',
  },
  {
    id: 2,
    name: 'Plan',
    chinese: '计划',
    desc: 'The written blueprint. Explains WHAT will be built, WHY each decision was made, and HOW pieces fit together.',
    file: 'plan-{slug}.md',
    color: '#3b82f6',
    detail: 'Every design decision must have a "because". Specific enough that another developer could implement without asking questions.',
    without: 'Code becomes the first draft of thought. There is nothing to review before execution, nothing to annotate, nothing to push back on. The user can\'t catch a bad plan that only exists in the AI\'s head.',
  },
  {
    id: 3,
    name: 'Annotate',
    chinese: '校准',
    desc: '1-6 calibration cycles between your understanding and the user\'s intent. Annotations are additive, not destructive.',
    file: 'plan-{slug}.md (annotated)',
    color: '#a855f7',
    detail: 'The most powerful insight: refine a plan by annotating it, not rewriting it. This is where creativity happens. Implementation should be boring.',
    without: 'The plan is your interpretation of what the user wants — not theirs. You execute perfectly on the wrong vision. This is the "I gave you exactly what you asked for, not what you needed" failure mode.',
  },
  {
    id: 4,
    name: 'Todo List',
    chinese: '任务分解',
    desc: 'The approved plan is decomposed into executable tasks. The mapping is mechanical, not creative.',
    file: 'TaskCreate calls',
    color: '#f59e0b',
    detail: 'Each section of the Implementation Sequence becomes one or more tasks. Tasks reference the plan document.',
    without: 'Implementation is one big untracked blob. No parallelism. No progress visibility. No ability to assign domains to teammates. The swarm cannot execute what hasn\'t been decomposed.',
  },
  {
    id: 5,
    name: 'Implement',
    chinese: '执行',
    desc: 'If the Six-Step Discipline was followed correctly, implementation should be straightforward execution.',
    file: 'code files',
    color: '#22c55e',
    detail: 'The boring part. If implementation is creative, the plan wasn\'t detailed enough. Strategic decisions go back to annotation.',
    without: 'Without the prior 4 steps, implementation becomes firefighting. Every ambiguity becomes a micro-decision, every architectural question becomes a context-switch. The execution is creative — and that\'s the warning sign.',
  },
  {
    id: 6,
    name: 'Feedback',
    chinese: '反馈',
    desc: 'Close the loop. What worked, what didn\'t, surprises, and lessons for future missions.',
    file: 'feedback-{slug}.md',
    color: '#ef4444',
    detail: 'This file feeds the next Research phase. When a similar mission comes up, past feedback accelerates research.',
    without: 'Each mission starts from zero. No compound learning, no pattern recognition across missions. The 10th similar feature takes as long as the 1st. The loop never improves.',
  },
]

export function SixStepSection() {
  const { lang } = useLang()
  const ref = useRef(null)
  const isInView = useInView(ref, { once: true, margin: '-100px' })
  const [activeStep, setActiveStep] = useState<number | null>(null)

  return (
    <section id="six-steps" className="min-h-screen py-24 px-4">
      <div className="max-w-6xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 40 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
        >
          <h2 className="text-4xl font-bold mb-2">{t('sixstep.title', lang)}</h2>
          <p className="text-text-secondary text-lg mb-2">
            {t('sixstep.subtitle', lang)}
          </p>
          <p className="text-text-muted text-sm mb-12">
            {t('sixstep.desc', lang)}
          </p>
        </motion.div>

        {/* Circular diagram */}
        <div ref={ref} className="flex flex-col lg:flex-row gap-12 items-center">
          <div className="relative w-80 h-80 flex-shrink-0">
            <svg viewBox="0 0 300 300" className="w-full h-full">
              {steps.map((step, i) => {
                const angle = (i * 60 - 90) * (Math.PI / 180)
                const r = 110
                const x = 150 + r * Math.cos(angle)
                const y = 150 + r * Math.sin(angle)

                // Arrow to next step
                const midAngle = ((i + 0.5) * 60 - 90) * (Math.PI / 180)
                const arrowR = 95
                const ax = 150 + arrowR * Math.cos(midAngle)
                const ay = 150 + arrowR * Math.sin(midAngle)

                return (
                  <motion.g
                    key={step.id}
                    initial={{ opacity: 0, scale: 0 }}
                    animate={isInView ? { opacity: 1, scale: 1 } : {}}
                    transition={{ delay: i * 0.15, duration: 0.4 }}
                    onClick={() => setActiveStep(activeStep === i ? null : i)}
                    style={{ cursor: 'pointer' }}
                  >
                    {/* Connection arc hint */}
                    <text
                      x={ax}
                      y={ay}
                      textAnchor="middle"
                      dominantBaseline="middle"
                      fill="#484f58"
                      fontSize={10}
                    >
                      →
                    </text>
                    {/* Node */}
                    <circle
                      cx={x}
                      cy={y}
                      r={activeStep === i ? 28 : 24}
                      fill={step.color + '20'}
                      stroke={step.color}
                      strokeWidth={activeStep === i ? 3 : 2}
                    />
                    <text
                      x={x}
                      y={y - 3}
                      textAnchor="middle"
                      fill={step.color}
                      fontSize={8}
                      fontWeight="bold"
                      fontFamily="monospace"
                    >
                      {step.id}
                    </text>
                    <text
                      x={x}
                      y={y + 8}
                      textAnchor="middle"
                      fill="#8b949e"
                      fontSize={7}
                      fontFamily="monospace"
                    >
                      {step.chinese}
                    </text>
                  </motion.g>
                )
              })}
              {/* Center text */}
              <text x={150} y={146} textAnchor="middle" fill="#f59e0b" fontSize={10} fontWeight="bold" fontFamily="monospace">
                THE LOOP
              </text>
              <text x={150} y={160} textAnchor="middle" fill="#484f58" fontSize={8} fontFamily="monospace">
                is the hero
              </text>
            </svg>
          </div>

          {/* Step details */}
          <div className="flex-1 space-y-3">
            {steps.map((step, i) => (
              <motion.div
                key={step.id}
                initial={{ opacity: 0, x: 20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.08 }}
                onClick={() => setActiveStep(activeStep === i ? null : i)}
                className={`p-4 rounded-lg border cursor-pointer transition-all ${
                  activeStep === i
                    ? 'bg-bg-card-hover border-gold/30'
                    : 'bg-bg-card border-border/50 hover:border-border hover:bg-bg-card-hover'
                }`}
                style={{ borderLeftColor: step.color, borderLeftWidth: 3 }}
              >
                <div className="flex items-center gap-3">
                  <span className="font-mono text-sm font-bold" style={{ color: step.color }}>
                    {step.id}. {step.name}
                  </span>
                  <span className="text-text-muted text-xs">{step.chinese}</span>
                  <span className="ml-auto text-xs font-mono text-text-muted">{step.file}</span>
                </div>
                <p className="text-sm text-text-secondary mt-1">{step.desc}</p>
                {activeStep === i && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    className="mt-2 pt-2 border-t border-border/50 space-y-2"
                  >
                    <p className="text-xs text-text-muted">{step.detail}</p>
                    <div className="bg-neon-red/5 border border-neon-red/10 rounded p-2">
                      <p className="text-[10px] font-mono text-neon-red mb-1">WITHOUT THIS STEP:</p>
                      <p className="text-xs text-text-secondary">{step.without}</p>
                    </div>
                  </motion.div>
                )}
              </motion.div>
            ))}
          </div>
        </div>

        {/* Key insight + Why 6 steps */}
        <motion.div
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          className="mt-12 grid md:grid-cols-2 gap-4"
        >
          <div className="bg-gold/5 border border-gold/20 rounded-lg p-6">
            <p className="text-lg font-mono text-gold mb-2">
              "The annotation cycle is where creativity happens."
            </p>
            <p className="text-sm text-text-secondary">
              Implementation should be boring. If implementation is creative, the plan wasn't detailed enough. That's a signal to go back to annotation — not to improvise forward.
            </p>
          </div>
          <div className="bg-bg-card border border-border rounded-lg p-6">
            <p className="text-xs font-mono text-text-muted mb-3">WHY 6 STEPS AND NOT 4 OR 8?</p>
            <p className="text-sm text-text-secondary mb-3">
              The 6 steps map to the 3 failure modes of AI-assisted development:
            </p>
            <div className="space-y-2 text-xs">
              <div className="flex items-start gap-2">
                <span className="text-neon-red font-mono mt-0.5">1-2</span>
                <p className="text-text-secondary"><span className="text-text-primary">Research+Plan</span> — prevent "building the wrong thing" (wrong direction)</p>
              </div>
              <div className="flex items-start gap-2">
                <span className="text-neon-red font-mono mt-0.5">3</span>
                <p className="text-text-secondary"><span className="text-text-primary">Annotate</span> — prevents "building MY interpretation of the right thing" (miscalibration)</p>
              </div>
              <div className="flex items-start gap-2">
                <span className="text-neon-red font-mono mt-0.5">4-5</span>
                <p className="text-text-secondary"><span className="text-text-primary">TodoList+Implement</span> — parallel execution at maximum throughput (speed)</p>
              </div>
              <div className="flex items-start gap-2">
                <span className="text-neon-red font-mono mt-0.5">6</span>
                <p className="text-text-secondary"><span className="text-text-primary">Feedback</span> — compound learning across missions (compounding)</p>
              </div>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  )
}
