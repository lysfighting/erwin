import { motion, useInView } from 'framer-motion'
import { useRef } from 'react'
import { SectionWrapper } from '../components/SectionWrapper'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

interface WaveData {
  id: number
  name: string
  tasks: { label: string; agent: string; color: string; width: number }[]
  criticalPath?: string
}

const waves: WaveData[] = [
  {
    id: 0,
    name: 'Bootstrap',
    tasks: [
      { label: 'Scaffold + verify', agent: 'Team Lead', color: '#f59e0b', width: 30 },
      { label: 'Research API X', agent: 'STL-backend', color: '#3b82f6', width: 25 },
      { label: 'Setup test infra', agent: 'tester', color: '#a855f7', width: 20 },
    ],
    criticalPath: 'Scaffold + verify',
  },
  {
    id: 1,
    name: 'Core Implementation',
    tasks: [
      { label: 'Frontend components', agent: 'STL-frontend (3 inner)', color: '#3b82f6', width: 60 },
      { label: 'API endpoints', agent: 'STL-backend (2 inner)', color: '#3b82f6', width: 55 },
      { label: 'Data model', agent: 'Team Lead', color: '#f59e0b', width: 40 },
    ],
    criticalPath: 'Data model',
  },
  {
    id: 2,
    name: 'Integration + Testing',
    tasks: [
      { label: 'Write tests', agent: 'tester', color: '#a855f7', width: 45 },
      { label: 'Integrate FE + BE', agent: 'Team Lead', color: '#f59e0b', width: 50 },
      { label: 'Frontend polish', agent: 'STL-frontend', color: '#3b82f6', width: 35 },
      { label: 'Error handling', agent: 'STL-backend', color: '#3b82f6', width: 30 },
    ],
    criticalPath: 'Integrate FE + BE',
  },
  {
    id: 3,
    name: 'Hardening',
    tasks: [
      { label: 'E2E tests', agent: 'tester', color: '#a855f7', width: 40 },
      { label: 'Final verification', agent: 'Team Lead', color: '#f59e0b', width: 35 },
    ],
    criticalPath: 'Final verification',
  },
]

const concurrencyPatterns = [
  {
    name: 'Fork-Join',
    desc: 'N independent tasks → all run in parallel → join when all complete',
    example: 'Wave 1: 3 agents start simultaneously, wait for all before Wave 2',
    whenToUse: 'Use when N tasks are fully independent and you need all results before proceeding',
    icon: '🔀',
    color: '#3b82f6',
  },
  {
    name: 'Pipeline',
    desc: 'Output of stage N feeds N+1, but stages can overlap in time',
    example: 'Tester writes test skeletons during Wave 1 before implementation completes',
    whenToUse: 'Use when a later stage can START with partial input from an earlier stage',
    icon: '⏩',
    color: '#a855f7',
  },
  {
    name: 'Speculative',
    desc: 'Start the likely path while validating; if wrong, cost is small',
    example: 'Begin implementing the obvious auth approach while exploring edge cases',
    whenToUse: 'Use when one path has >80% probability of being correct — the save outweighs the risk',
    icon: '🎯',
    color: '#22c55e',
  },
]

export function WaveExecutionSection() {
  const { lang } = useLang()
  const ref = useRef(null)
  const isInView = useInView(ref, { once: true, margin: '-50px' })

  return (
    <SectionWrapper
      id="waves"
      title={t('waves.title', lang)}
      subtitle={t('waves.subtitle', lang)}
    >
      {/* Why Waves? */}
      <div className="grid md:grid-cols-2 gap-4 mb-10">
        <div className="bg-neon-red/5 border border-neon-red/20 rounded-lg p-5">
          <p className="text-xs font-mono text-neon-red font-bold mb-2">THE PROBLEM WITH TASK LISTS</p>
          <p className="text-sm text-text-secondary mb-3">
            Writing "Task 1, Task 2, Task 3" implies serial execution. Nobody explicitly says these can run simultaneously, so they don't. Serial execution is the path of least resistance when structure is ambiguous.
          </p>
          <div className="flex gap-1">
            {Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="h-6 w-8 bg-neon-red/20 border border-neon-red/30 rounded flex items-center justify-center">
                <span className="text-[8px] text-neon-red font-mono">{i+1}</span>
              </div>
            ))}
          </div>
          <p className="text-xs text-text-muted mt-2 font-mono">{'\u2192'} 6 units wall-clock time</p>
        </div>
        <div className="bg-neon-green/5 border border-neon-green/20 rounded-lg p-5">
          <p className="text-xs font-mono text-neon-green font-bold mb-2">WHY WAVES FIX THIS</p>
          <p className="text-sm text-text-secondary mb-3">
            Waves make parallelism explicit and mandatory. Tasks in the same wave MUST run simultaneously — no ambiguity. The Wave Map is an execution contract, not a suggestion.
          </p>
          <div className="flex gap-2">
            {[[3],[3],[2]].map((col, wi) => (
              <div key={wi} className="flex flex-col gap-1">
                {col.map((_, i) => (
                  <div key={i} className="h-6 w-10 bg-neon-green/20 border border-neon-green/30 rounded flex items-center justify-center">
                    <span className="text-[8px] text-neon-green font-mono">W{wi}</span>
                  </div>
                ))}
              </div>
            ))}
          </div>
          <p className="text-xs text-text-muted mt-2 font-mono">{'\u2192'} 3 units wall-clock time (2x faster)</p>
        </div>
      </div>

      {/* Wave Gantt chart */}
      <div ref={ref} className="bg-bg-card border border-border rounded-lg p-6 mb-12 overflow-x-auto">
        <div className="min-w-[600px]">
          {waves.map((wave, wi) => (
            <div key={wave.id} className="mb-6 last:mb-0">
              <div className="flex items-center gap-3 mb-2">
                <span className="font-mono text-xs text-text-muted w-24">Wave {wave.id}</span>
                <span className="text-sm font-mono text-gold">{wave.name}</span>
                <span className="text-xs text-text-muted">Width: {wave.tasks.length}</span>
              </div>
              <div className="space-y-1.5 pl-24">
                {wave.tasks.map((task, ti) => (
                  <div key={task.label} className="flex items-center gap-2">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={isInView ? { width: `${task.width}%` } : {}}
                      transition={{ delay: wi * 0.6 + ti * 0.1, duration: 0.5, ease: 'easeOut' }}
                      className="h-7 rounded flex items-center px-3 relative overflow-hidden"
                      style={{
                        backgroundColor: task.color + '20',
                        borderLeft: `3px solid ${task.color}`,
                      }}
                    >
                      <span className="text-xs font-mono truncate" style={{ color: task.color }}>
                        {task.label}
                      </span>
                      {task.label === wave.criticalPath && (
                        <span className="ml-auto text-[10px] text-neon-red font-mono">CRITICAL</span>
                      )}
                    </motion.div>
                    <span className="text-[10px] text-text-muted font-mono whitespace-nowrap">{task.agent}</span>
                  </div>
                ))}
              </div>
            </div>
          ))}

          {/* Critical path line */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={isInView ? { opacity: 1 } : {}}
            transition={{ delay: 2.5 }}
            className="mt-4 pl-24 pt-3 border-t border-border"
          >
            <div className="flex items-center gap-2 text-xs font-mono">
              <span className="text-neon-red">Critical Path:</span>
              <span className="text-text-muted">T1 → T6 → T8 → T12</span>
              <span className="text-text-muted">|</span>
              <span className="text-text-muted">4 waves</span>
            </div>
          </motion.div>
        </div>
      </div>

      {/* Speedup comparison */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        className="grid md:grid-cols-2 gap-6 mb-12"
      >
        <div className="bg-bg-card border border-neon-red/20 rounded-lg p-6">
          <h3 className="font-mono text-neon-red text-sm mb-4">Serial Execution</h3>
          <div className="flex gap-1 mb-3">
            {Array.from({ length: 12 }).map((_, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0 }}
                whileInView={{ opacity: 1 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.05 }}
                className="h-6 w-8 bg-neon-red/20 border border-neon-red/30 rounded"
              />
            ))}
          </div>
          <p className="font-mono text-xs text-text-muted">12 tasks × 1 unit = <span className="text-neon-red">12 time units</span></p>
        </div>
        <div className="bg-bg-card border border-neon-green/20 rounded-lg p-6">
          <h3 className="font-mono text-neon-green text-sm mb-4">Parallel Swarm Execution</h3>
          <div className="flex gap-1 mb-3">
            {[3, 3, 4, 2].map((count, wi) => (
              <div key={wi} className="flex flex-col gap-0.5">
                {Array.from({ length: count }).map((_, i) => (
                  <motion.div
                    key={i}
                    initial={{ opacity: 0, scale: 0 }}
                    whileInView={{ opacity: 1, scale: 1 }}
                    viewport={{ once: true }}
                    transition={{ delay: wi * 0.2 + i * 0.05 }}
                    className="h-6 w-8 bg-neon-green/20 border border-neon-green/30 rounded"
                  />
                ))}
              </div>
            ))}
          </div>
          <p className="font-mono text-xs text-text-muted">4 waves = <span className="text-neon-green">4 time units</span> → <span className="text-gold font-bold">3x speedup</span></p>
        </div>
      </motion.div>

      {/* Three concurrency patterns */}
      <h3 className="font-mono text-lg text-text-primary mb-4">Three Concurrency Patterns</h3>
      <div className="grid md:grid-cols-3 gap-4">
        {concurrencyPatterns.map((pattern, i) => (
          <motion.div
            key={pattern.name}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: i * 0.1 }}
            className="bg-bg-card border border-border rounded-lg p-5 hover:border-gold/30 transition-colors"
            style={{ borderTopColor: pattern.color, borderTopWidth: 2 }}
          >
            <div className="flex items-center gap-2 mb-3">
              <span className="text-xl">{pattern.icon}</span>
              <span className="font-mono font-bold text-gold">{pattern.name}</span>
            </div>
            <p className="text-sm text-text-secondary mb-3">{pattern.desc}</p>
            <p className="text-xs text-text-muted italic">{pattern.example}</p>
            <div className="mt-3 pt-3 border-t border-border/30">
              <p className="text-[10px] font-mono text-text-muted mb-1">WHEN TO USE:</p>
              <p className="text-xs text-text-secondary">{pattern.whenToUse}</p>
            </div>
          </motion.div>
        ))}
      </div>

      {/* Source attribution */}
      <motion.p
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
        className="mt-8 text-xs text-text-muted font-mono"
      >
        Concurrency patterns adapted from: Fork-Join (Conway, 1963), Pipeline Parallelism (computer architecture), Speculative Execution (CPU branch prediction)
      </motion.p>
    </SectionWrapper>
  )
}
