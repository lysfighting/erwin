import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'

const sectionTopics = [
  'Architecture Overview',
  'Execution Phases',
  'Swarm Topology',
  'Wave Execution Model',
  'Six-Step Discipline',
  'Design Philosophy',
  'Deep Reasoning Protocol',
  'Direction Validation',
  'Pantheon of Minds',
  'Self-Adversarial Review',
  'Anti-Patterns',
  'Self-Drive Mode',
  'Architect Command',
  'Prompt Engineering',
  'References & Ecosystem',
]

export function ActIIIWrapper({ children }: { children: React.ReactNode }) {
  const [isExpanded, setIsExpanded] = useState(false)

  return (
    <div id="act-iii">
      {/* Act III Header */}
      <div className="max-w-6xl mx-auto px-4 pt-24 pb-8">
        <motion.div
          initial={{ opacity: 0, y: 40 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-100px' }}
          transition={{ duration: 0.6 }}
        >
          <h2 className="text-3xl md:text-4xl font-bold text-text-primary mb-2">
            <span className="text-gold font-mono">ACT III:</span>{' '}
            <span className="text-gold/80">深度</span>{' '}
            <span className="text-text-secondary text-2xl md:text-3xl">
              — Deep Reference ({sectionTopics.length} sections)
            </span>
          </h2>
        </motion.div>
      </div>

      {/* Collapsed summary / Expand toggle */}
      <AnimatePresence mode="wait">
        {!isExpanded && (
          <motion.div
            key="collapsed"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0, height: 0 }}
            transition={{ duration: 0.3 }}
            className="max-w-6xl mx-auto px-4 pb-12"
          >
            <div className="bg-bg-card border border-border rounded-xl p-6">
              <p className="text-text-secondary text-sm font-mono mb-4">
                This section covers the full technical depth of drive-skill:
              </p>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-2 mb-6">
                {sectionTopics.map((topic, i) => (
                  <div
                    key={topic}
                    className="text-xs font-mono text-text-muted px-3 py-1.5 bg-bg-dark rounded border border-border/50"
                  >
                    <span className="text-gold/60 mr-1.5">{String(i + 1).padStart(2, '0')}</span>
                    {topic}
                  </div>
                ))}
              </div>
              <button
                onClick={() => setIsExpanded(true)}
                className="w-full py-3 bg-gold/10 hover:bg-gold/20 border border-gold/20 rounded-lg text-gold font-mono text-sm transition-all"
              >
                Expand to read →
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Expanded content */}
      <AnimatePresence>
        {isExpanded && (
          <motion.div
            key="expanded"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.4 }}
          >
            {/* Collapse button */}
            <div className="max-w-6xl mx-auto px-4 pb-4">
              <button
                onClick={() => setIsExpanded(false)}
                className="px-4 py-2 bg-bg-card hover:bg-bg-card-hover border border-border rounded-lg text-text-secondary hover:text-text-primary font-mono text-xs transition-all"
              >
                ← Hide Deep Content
              </button>
            </div>
            {children}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}
