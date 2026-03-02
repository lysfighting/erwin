import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'

interface ThoughtSource {
  origin: string
  arrow: string // what it became in drive-skill
}

interface Principle {
  number: number
  zhName: string
  enName: string
  summary: string
  detail: string
  sources: ThoughtSource[]
  accentColor: string
}

const PRINCIPLES: Principle[] = [
  {
    number: 1,
    zhName: '永不停止的循环',
    enName: 'The Loop Never Breaks',
    summary:
      'The system never halts. Every error is a recovery opportunity, every completion triggers the next cycle. Autonomous momentum is the cardinal rule.',
    detail:
      'Inspired by cybernetic feedback loops, the OODA cycle, and lean manufacturing flow. drive-skill treats interruptions not as failures but as signals to re-orient and continue. The loop runs until the mission is complete — no human nudge required.',
    sources: [
      { origin: 'OODA Loop (Boyd, 1970s)', arrow: 'Cardinal Rule: never stop' },
      { origin: 'Cybernetics (Wiener, 1948)', arrow: 'Self-regulating feedback' },
    ],
    accentColor: '#f59e0b', // gold
  },
  {
    number: 2,
    zhName: '层级指挥结构',
    enName: 'Hierarchical Command Structure',
    summary:
      'A single STL (Sub-Team Lead) coordinates up to 4 agents. Intent flows down, results flow up. Execution is autonomous at every level.',
    detail:
      'Borrowed from military Mission Command (Auftragstaktik): leaders communicate intent, not step-by-step instructions. Each STL decomposes its mandate into parallel sub-tasks and dispatches them without waiting for central approval.',
    sources: [
      { origin: 'Mission Command (Prussian military)', arrow: 'STL hierarchy' },
      { origin: 'The Mythical Man-Month (Brooks)', arrow: 'Surgical team model' },
    ],
    accentColor: '#3b82f6', // blue
  },
  {
    number: 3,
    zhName: '最大并行度',
    enName: 'Maximum Parallelism',
    summary:
      'If tasks have no dependency, they run simultaneously. The critical path is minimized by default. Serial execution is treated as a bug.',
    detail:
      'Rooted in Critical Path Method (CPM) and Amdahl\'s Law. drive-skill analyzes task dependencies at planning time and maximizes the parallel width of each execution wave. The goal: wall-clock time approaches the longest single task, not the sum of all tasks.',
    sources: [
      { origin: 'CPM / PERT (1957)', arrow: 'Wave-based execution' },
      { origin: "Amdahl's Law (1967)", arrow: 'Minimize serial fraction' },
    ],
    accentColor: '#22c55e', // green
  },
  {
    number: 4,
    zhName: '透明广播协议',
    enName: 'Transparent Broadcast Protocol',
    summary:
      'Every dispatch and return is broadcast to all agents. No information silos. Visibility eliminates uncertainty and prevents duplicated work.',
    detail:
      'Drawn from Shannon\'s information theory (visibility reduces entropy) and agile stand-up rituals. Every agent sees every message in the system, enabling self-coordination without explicit synchronization points.',
    sources: [
      { origin: 'Information Theory (Shannon, 1948)', arrow: 'Broadcast protocol' },
      { origin: 'Agile Daily Stand-up', arrow: 'Continuous sync' },
    ],
    accentColor: '#a855f7', // purple
  },
  {
    number: 5,
    zhName: '阶段纪律',
    enName: 'Phase Discipline',
    summary:
      'Thinking before doing. Phase 0 (research) and Phase T (team formation) run before any code is written. Preparation prevents rework.',
    detail:
      'Echoes Sun Tzu\'s "plan before you move" and the Double Diamond model\'s diverge-then-converge rhythm. By enforcing distinct phases, drive-skill ensures the system understands the problem space before committing resources to execution.',
    sources: [
      { origin: 'Sun Tzu — The Art of War', arrow: 'Phase 0 research' },
      { origin: 'Double Diamond (Design Council, 2004)', arrow: 'Diverge then converge' },
    ],
    accentColor: '#06b6d4', // cyan
  },
]

function PrincipleCard({ principle }: { principle: Principle }) {
  const [expanded, setExpanded] = useState(false)

  return (
    <motion.div
      className="bg-bg-card border border-border rounded-xl overflow-hidden"
      initial={{ opacity: 0, y: 30 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: '-60px' }}
      transition={{ duration: 0.5, delay: principle.number * 0.08 }}
    >
      {/* Card header */}
      <div className="p-5 md:p-6">
        {/* Number badge + names */}
        <div className="flex items-start gap-3 mb-3">
          <span
            className="flex-shrink-0 w-8 h-8 rounded-lg flex items-center justify-center font-mono text-sm font-bold text-bg-dark"
            style={{ backgroundColor: principle.accentColor }}
          >
            {principle.number}
          </span>
          <div>
            <h3 className="font-mono text-lg font-bold text-text-primary leading-tight">
              {principle.zhName}
            </h3>
            <p
              className="font-mono text-xs mt-0.5"
              style={{ color: principle.accentColor }}
            >
              {principle.enName}
            </p>
          </div>
        </div>

        {/* Summary */}
        <p className="text-text-secondary text-sm leading-relaxed">
          {principle.summary}
        </p>

        {/* Expand toggle */}
        <button
          onClick={() => setExpanded(!expanded)}
          className="mt-3 font-mono text-xs text-text-muted hover:text-text-secondary transition-colors flex items-center gap-1"
        >
          <motion.span
            animate={{ rotate: expanded ? 90 : 0 }}
            transition={{ duration: 0.2 }}
            className="inline-block"
          >
            &rsaquo;
          </motion.span>
          {expanded ? 'Collapse' : 'Read more'}
        </button>
      </div>

      {/* Expandable detail */}
      <AnimatePresence initial={false}>
        {expanded && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.3, ease: 'easeInOut' }}
            className="overflow-hidden"
          >
            <div className="px-5 pb-5 md:px-6 md:pb-6 border-t border-border pt-4 space-y-4">
              {/* Detail text */}
              <p className="text-text-secondary text-sm leading-relaxed">
                {principle.detail}
              </p>

              {/* Thought sources */}
              <div>
                <span className="font-mono text-xs text-text-muted uppercase tracking-wider">
                  Thought Sources
                </span>
                <div className="mt-2 space-y-2">
                  {principle.sources.map((source, i) => (
                    <div
                      key={i}
                      className="flex items-center gap-2 text-xs font-mono"
                    >
                      <span className="text-text-secondary">{source.origin}</span>
                      <span style={{ color: principle.accentColor }}>&rarr;</span>
                      <span className="text-text-muted">{source.arrow}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  )
}

export function PrincipleCards() {
  return (
    <div className="w-full">
      {/* Header */}
      <div className="text-center mb-8">
        <h2 className="font-mono text-2xl md:text-3xl font-bold text-text-primary mb-2">
          Five Core Principles
        </h2>
        <p className="text-text-secondary text-sm">
          The design ideas behind drive-skill, traced to their origins
        </p>
      </div>

      {/* Cards grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6">
        {PRINCIPLES.map((p) => (
          <PrincipleCard key={p.number} principle={p} />
        ))}
      </div>
    </div>
  )
}
