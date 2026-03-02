import { motion, AnimatePresence } from 'framer-motion'
import { useState, useCallback } from 'react'

interface TaskBar {
  label: string
  duration: number // seconds in the simulation world
  color: string
}

const WITHOUT_TASKS: TaskBar[] = [
  { label: 'Task 1: Research', duration: 4, color: '#ef4444' },
  { label: 'Task 2: Implement', duration: 4, color: '#ef4444' },
  { label: 'Task 3: Test', duration: 4, color: '#ef4444' },
]

const WITH_TASKS: TaskBar[] = [
  { label: 'Task 1: Research', duration: 4, color: '#22c55e' },
  { label: 'Task 2: Implement', duration: 4, color: '#22c55e' },
  { label: 'Task 3: Test', duration: 4, color: '#22c55e' },
  { label: 'Task 4: Document', duration: 4, color: '#22c55e' },
]

// Animation durations in real seconds
const SEQ_TASK_DURATION = 1.2 // each sequential task takes 1.2s real time
const PAR_TASK_DURATION = 1.6 // parallel tasks all animate over 1.6s real time
const STAGGER_DELAY = 0.1 // slight stagger for parallel to feel organic

function SequentialSide({ playing, runKey }: { playing: boolean; runKey: number }) {
  return (
    <div className="flex-1 min-w-0">
      <div className="flex items-center gap-2 mb-4">
        <div className="w-2 h-2 rounded-full bg-neon-red" />
        <h3 className="font-mono text-sm text-text-secondary uppercase tracking-wider">
          Without drive-skill
        </h3>
      </div>

      <div className="space-y-3">
        {WITHOUT_TASKS.map((task, i) => (
          <div key={`${runKey}-${i}`} className="space-y-1">
            <span className="font-mono text-xs text-text-muted">{task.label}</span>
            <div className="h-8 bg-bg-dark rounded-md overflow-hidden border border-border">
              {playing && (
                <motion.div
                  className="h-full rounded-md"
                  style={{ backgroundColor: task.color }}
                  initial={{ width: '0%' }}
                  animate={{ width: '100%' }}
                  transition={{
                    duration: SEQ_TASK_DURATION,
                    delay: i * SEQ_TASK_DURATION,
                    ease: 'easeInOut',
                  }}
                />
              )}
            </div>
          </div>
        ))}
      </div>

      <AnimatePresence>
        {playing && (
          <motion.div
            className="mt-4 font-mono text-lg text-neon-red"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: WITHOUT_TASKS.length * SEQ_TASK_DURATION + 0.2 }}
          >
            Total: 12s
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}

function ParallelSide({ playing, runKey }: { playing: boolean; runKey: number }) {
  return (
    <div className="flex-1 min-w-0">
      <div className="flex items-center gap-2 mb-4">
        <div className="w-2 h-2 rounded-full bg-neon-green" />
        <h3 className="font-mono text-sm text-text-secondary uppercase tracking-wider">
          With drive-skill
        </h3>
      </div>

      <div className="space-y-3">
        {WITH_TASKS.map((task, i) => (
          <div key={`${runKey}-${i}`} className="space-y-1">
            <span className="font-mono text-xs text-text-muted">{task.label}</span>
            <div className="h-8 bg-bg-dark rounded-md overflow-hidden border border-border">
              {playing && (
                <motion.div
                  className="h-full rounded-md"
                  style={{ backgroundColor: task.color }}
                  initial={{ width: '0%' }}
                  animate={{ width: '100%' }}
                  transition={{
                    duration: PAR_TASK_DURATION,
                    delay: i * STAGGER_DELAY,
                    ease: 'easeInOut',
                  }}
                />
              )}
            </div>
          </div>
        ))}
      </div>

      <AnimatePresence>
        {playing && (
          <motion.div
            className="mt-4 font-mono text-lg text-neon-green flex items-baseline gap-2"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: PAR_TASK_DURATION + 0.3 }}
          >
            <span>Total: 4s</span>
            <motion.span
              className="text-sm text-gold"
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: PAR_TASK_DURATION + 0.6 }}
            >
              &larr; 3x faster
            </motion.span>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}

export function ExecutionSimulator() {
  const [playing, setPlaying] = useState(false)
  const [runKey, setRunKey] = useState(0)

  const handlePlay = useCallback(() => {
    setPlaying(false)
    // Force remount of animation elements with new key
    setRunKey((k) => k + 1)
    // Use rAF to ensure state flush before re-enabling
    requestAnimationFrame(() => {
      setPlaying(true)
    })
  }, [])

  return (
    <div className="w-full">
      {/* Header */}
      <div className="text-center mb-8">
        <h2 className="font-mono text-2xl md:text-3xl font-bold text-text-primary mb-2">
          Sequential vs Parallel Execution
        </h2>
        <p className="text-text-secondary text-sm">
          See how drive-skill orchestrates tasks in parallel instead of one-by-one
        </p>
      </div>

      {/* Simulation area */}
      <div className="bg-bg-card border border-border rounded-xl p-6 md:p-8">
        <div className="flex flex-col md:flex-row gap-8 md:gap-12">
          <SequentialSide playing={playing} runKey={runKey} />

          {/* Divider */}
          <div className="hidden md:flex flex-col items-center justify-center">
            <div className="w-px h-full bg-border" />
          </div>
          <div className="md:hidden border-t border-border" />

          <ParallelSide playing={playing} runKey={runKey} />
        </div>

        {/* Play / Replay button */}
        <div className="mt-8 flex justify-center">
          <motion.button
            onClick={handlePlay}
            className="px-6 py-2.5 bg-gold/10 border border-gold/30 rounded-lg font-mono text-sm text-gold hover:bg-gold/20 transition-colors"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            {playing ? 'Play Again' : 'Run Simulation'}
          </motion.button>
        </div>
      </div>
    </div>
  )
}
