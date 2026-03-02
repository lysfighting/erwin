import { motion } from 'framer-motion'
import { SectionWrapper } from '../components/SectionWrapper'
import { useState } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

// WHY THIS EXISTS
// ────────────────
// The gap between "it works" and "it's correct, robust, and complete" is enormous.
// An AI's instinct after finishing is to declare done and move on.
// The Self-Adversarial Review forces the opposite: actively try to break what you built.
// Only when the inner critic finds nothing left to attack is the task truly complete.

const axes = [
  {
    id: 1,
    name: 'Correctness',
    nameCh: '正确性',
    icon: '🎯',
    color: '#ef4444',
    question: 'If I were trying to BREAK this, where would I attack?',
    details: [
      'What inputs would cause failure?',
      'What race conditions or timing issues exist?',
      'What null cases or boundary conditions did I miss?',
      'What happens with the minimum valid input? The maximum?',
      'What if two things happen simultaneously?',
    ],
    why: 'Code that works with expected inputs fails with unexpected ones. The attacker mindset catches what the builder mindset misses — because builders optimize for the happy path, attackers probe the edges.',
  },
  {
    id: 2,
    name: 'Completeness',
    nameCh: '完整性',
    icon: '📋',
    color: '#f59e0b',
    question: 'Did I implement EVERYTHING, or just the easy 80%?',
    details: [
      'Am I leaving implicit TODOs (unhandled cases I\'ll "fix later")?',
      'Is error handling complete, or just happy-path?',
      'Are loading states, empty states, and error states all handled?',
      'Did I skip the "hard parts" of the task?',
      'Would the user find a missing feature in 5 minutes of use?',
    ],
    why: 'The last 20% is the mission. Features without error handling, APIs without validation, UIs without loading states — these aren\'t "80% done." They\'re broken in ways that matter.',
  },
  {
    id: 3,
    name: 'Quality',
    nameCh: '质量',
    icon: '⚗️',
    color: '#a855f7',
    question: 'Would a senior engineer approve this in code review?',
    details: [
      'What would a hostile reviewer flag immediately?',
      'What\'s the weakest part of this implementation?',
      'Is this code readable by someone who didn\'t write it?',
      'Are there magic numbers or unexplained decisions?',
      'Is the naming clear, consistent, and accurate?',
    ],
    why: 'Not "does it work?" but "would it survive a hostile code review?" The standard isn\'t personal satisfaction — it\'s whether an external expert would approve it.',
  },
  {
    id: 4,
    name: 'Scalability',
    nameCh: '可扩展性',
    icon: '📈',
    color: '#3b82f6',
    question: 'How does this behave as the problem grows 10× or 100×?',
    details: [
      'Is there an O(2^n) lurking where O(n) is possible?',
      'Did I hard-code what should be parameterized?',
      'Did I enumerate cases that should be learned from data?',
      'Will this approach still be correct when the system has grown 10×?',
      'Is the intelligence in the code (fragile) or in learnable patterns?',
    ],
    why: 'The Bitter Lesson: hand-crafted rules scale as O(2^complexity). Scalable designs move intelligence from rules to data. If the answer to "how does this scale?" is "we add more rules," that\'s a red flag.',
  },
  {
    id: 5,
    name: 'Integration',
    nameCh: '集成性',
    icon: '🔗',
    color: '#22c55e',
    question: 'Does this work in isolation AND in the full system?',
    details: [
      'Did I test it the way a REAL USER would use it, not just the way I built it?',
      'Does it compose correctly with the other modules?',
      'Are there any implicit assumptions about call order or state?',
      'Does it still work after a page refresh, session reset, or reconnect?',
      'What does the user see, end-to-end, when this code runs?',
    ],
    why: 'Unit tests prove a component works in a lab. Integration tests prove it works in the real system. Most bugs live at the boundaries — between your code and everything else.',
  },
  {
    id: 6,
    name: 'Regression',
    nameCh: '回归性',
    icon: '🔄',
    color: '#06b6d4',
    question: 'Did I break anything that was working before?',
    details: [
      'Did I run the full test suite, not just the tests for this feature?',
      'Did I check the files I modified for side effects?',
      'Are there other callers of the functions I changed?',
      'Did I check related functionality, not just the direct task?',
      'Did the build pass in CI (or local equivalent)?',
    ],
    why: 'New code breaks old code in non-obvious ways. A change to a shared utility, a renamed export, a shifted type — any of these can silently break something that worked yesterday. Run the full suite, every time.',
  },
]

const reviewSteps = [
  {
    step: '1',
    label: 'Re-Read Your Code',
    desc: 'Don\'t just remember what you wrote — read it fresh, line by line. Use the Read tool. You\'ll catch things your writing brain missed.',
    color: '#f59e0b',
  },
  {
    step: '2',
    label: 'Run Tests Again',
    desc: 'Even if they passed before your changes. Especially if they passed before. Your changes may have broken something.',
    color: '#3b82f6',
  },
  {
    step: '3',
    label: 'Try to Break It',
    desc: 'Run edge cases. Pass bad inputs. Check error messages. Become the attacker, not the builder. One real failure invalidates all your assumptions.',
    color: '#ef4444',
  },
  {
    step: '4',
    label: 'Fix Before Marking Done',
    desc: 'If you find ANY issue — fix it now. No exceptions. No "I\'ll address this in a follow-up." If you can do it, you must do it.',
    color: '#22c55e',
  },
]

const qualityBar = [
  'I re-read every line of code I wrote or modified',
  'I ran verification and saw it pass with my own eyes',
  'I actively tried to find problems and either found none or fixed them all',
  'I am confident this would survive a hostile code review',
]

export function SelfAdversarialReviewSection() {
  const { lang } = useLang()
  const [activeAxis, setActiveAxis] = useState<number | null>(null)
  const [checkedItems, setCheckedItems] = useState<Set<number>>(new Set())

  const toggleCheck = (i: number) => {
    const next = new Set(checkedItems)
    if (next.has(i)) next.delete(i)
    else next.add(i)
    setCheckedItems(next)
  }

  return (
    <SectionWrapper
      id="self-review"
      title={t('selfreview.title', lang)}
      subtitle={t('selfreview.subtitle', lang)}
    >
      {/* WHY THIS EXISTS */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        className="mb-12 bg-neon-red/5 border border-neon-red/20 rounded-xl p-8"
      >
        <p className="text-[10px] font-mono text-neon-red mb-4 tracking-widest">WHY THIS EXISTS — MANDATORY BEFORE MARKING ANY TASK COMPLETE</p>
        <h3 className="text-xl font-bold text-text-primary mb-4">
          "It Works" ≠ Done
        </h3>
        <div className="grid md:grid-cols-2 gap-6">
          <div>
            <p className="text-sm text-text-secondary leading-relaxed mb-3">
              An AI's natural instinct after finishing a task: declare done, move on. The code passes the test you wrote for the happy path. The obvious cases work. You feel the momentum of completion.
            </p>
            <p className="text-sm text-text-secondary leading-relaxed">
              But the user will hit an edge case you didn't test. A concurrent operation you didn't consider. An error state that throws instead of failing gracefully. A regression in code you touched but didn't think you changed.
            </p>
          </div>
          <div className="bg-bg-dark rounded-lg p-5 border border-border">
            <p className="text-xs font-mono text-neon-red mb-3">THE GAP</p>
            <div className="space-y-2">
              {[
                '"It works" → happy path tested',
                '"It\'s correct" → edge cases handled',
                '"It\'s complete" → last 20% done',
                '"It\'s robust" → errors handled gracefully',
                '"It\'s integrated" → works in full system',
                '"It\'s verified" → evidence for each claim',
              ].map((line, i) => (
                <div key={i} className="flex items-center gap-2 text-xs">
                  <div
                    className="w-2 h-2 rounded-full flex-shrink-0"
                    style={{ backgroundColor: i < 1 ? '#22c55e' : i < 3 ? '#f59e0b' : '#ef4444' }}
                  />
                  <span className={i < 1 ? 'text-neon-green' : i < 3 ? 'text-gold/70' : 'text-text-muted'}>
                    {line}
                  </span>
                </div>
              ))}
            </div>
            <p className="text-[10px] text-text-muted mt-3">Most "done" claims stop at line 1. The Self-Adversarial Review closes the gap.</p>
          </div>
        </div>
      </motion.div>

      {/* The 6 Axes */}
      <div className="mb-12">
        <h3 className="font-mono text-lg text-text-primary mb-2">The 6 Review Axes</h3>
        <p className="text-sm text-text-muted mb-6">
          Before marking any task complete, you MUST answer the trigger question for each axis. Click any axis to see the full checklist.
        </p>
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
          {axes.map((axis, i) => (
            <motion.div
              key={axis.id}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.07 }}
              onClick={() => setActiveAxis(activeAxis === i ? null : i)}
              className="bg-bg-card border border-border rounded-lg p-5 cursor-pointer hover:bg-bg-card-hover transition-all"
              style={{ borderTopColor: axis.color, borderTopWidth: 2 }}
            >
              <div className="flex items-center gap-3 mb-2">
                <span className="text-xl">{axis.icon}</span>
                <div>
                  <span className="font-mono text-sm font-bold text-text-primary">{axis.id}. {axis.name}</span>
                  <span className="text-xs text-text-muted ml-2">({axis.nameCh})</span>
                </div>
              </div>
              <p className="text-xs italic text-text-muted mb-2 leading-relaxed">"{axis.question}"</p>

              {activeAxis === i && (
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: 'auto' }}
                  className="mt-3 pt-3 border-t border-border/50 space-y-3"
                >
                  <div>
                    <p className="text-[10px] font-mono mb-1.5" style={{ color: axis.color }}>CHECKLIST:</p>
                    <ul className="space-y-1">
                      {axis.details.map((d, di) => (
                        <li key={di} className="flex items-start gap-2 text-xs text-text-secondary">
                          <span className="text-text-muted mt-0.5 flex-shrink-0">•</span>
                          {d}
                        </li>
                      ))}
                    </ul>
                  </div>
                  <div className="bg-bg-dark rounded p-3 border border-border/50">
                    <p className="text-[10px] font-mono text-text-muted mb-1">WHY THIS AXIS:</p>
                    <p className="text-xs text-text-secondary leading-relaxed">{axis.why}</p>
                  </div>
                </motion.div>
              )}

              <div className="mt-2 text-right">
                <span className="text-[10px] text-text-muted">{activeAxis === i ? '▲ collapse' : '▼ expand'}</span>
              </div>
            </motion.div>
          ))}
        </div>
      </div>

      {/* How to Review */}
      <div className="mb-12">
        <h3 className="font-mono text-lg text-text-primary mb-6">How to Actually Review</h3>
        <div className="grid md:grid-cols-4 gap-4">
          {reviewSteps.map((step, i) => (
            <motion.div
              key={step.step}
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.1 }}
              className="bg-bg-card border border-border rounded-lg p-4 relative"
            >
              <div
                className="absolute -top-3 left-4 text-xs font-mono font-bold px-2 py-0.5 rounded"
                style={{ backgroundColor: step.color, color: '#0d1117' }}
              >
                {step.step}
              </div>
              <p className="text-sm font-bold text-text-primary mt-1 mb-2">{step.label}</p>
              <p className="text-xs text-text-secondary leading-relaxed">{step.desc}</p>
            </motion.div>
          ))}
        </div>
      </div>

      {/* Quality Bar — Interactive Checklist */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        className="bg-bg-card border border-border rounded-xl p-8"
      >
        <h3 className="font-mono text-lg text-text-primary mb-2">The Quality Bar</h3>
        <p className="text-sm text-text-muted mb-6">
          The review is not a checkbox. You must be able to answer YES to ALL of these before marking complete:
        </p>
        <div className="space-y-3 mb-6">
          {qualityBar.map((item, i) => (
            <motion.button
              key={i}
              onClick={() => toggleCheck(i)}
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.08 }}
              className="w-full text-left flex items-center gap-3 p-4 rounded-lg border transition-all hover:bg-bg-card-hover"
              style={{
                borderColor: checkedItems.has(i) ? '#22c55e50' : '#484f5850',
                backgroundColor: checkedItems.has(i) ? '#22c55e08' : undefined,
              }}
            >
              <div
                className="w-5 h-5 rounded flex items-center justify-center flex-shrink-0 border transition-all"
                style={{
                  borderColor: checkedItems.has(i) ? '#22c55e' : '#484f58',
                  backgroundColor: checkedItems.has(i) ? '#22c55e20' : undefined,
                }}
              >
                {checkedItems.has(i) && <span className="text-neon-green text-xs">✓</span>}
              </div>
              <span className={`text-sm transition-colors ${checkedItems.has(i) ? 'text-neon-green' : 'text-text-secondary'}`}>
                {item}
              </span>
            </motion.button>
          ))}
        </div>

        <div
          className="rounded-lg p-4 text-center transition-all"
          style={{
            backgroundColor: checkedItems.size === 4 ? '#22c55e10' : '#484f5810',
            borderColor: checkedItems.size === 4 ? '#22c55e50' : '#484f5840',
            border: '1px solid',
          }}
        >
          {checkedItems.size === 4 ? (
            <p className="font-mono text-neon-green font-bold">TASK_COMPLETE ✓ — All axes satisfied. Ship it.</p>
          ) : (
            <p className="font-mono text-text-muted">
              {checkedItems.size}/4 checks passed — {4 - checkedItems.size} remaining before TASK_COMPLETE
            </p>
          )}
        </div>
      </motion.div>

      {/* Pre-Completion Exhaustion Check */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        className="mt-8 p-6 bg-gold/5 border border-gold/20 rounded-xl"
      >
        <h3 className="font-mono text-gold mb-4">Pre-Completion Exhaustion Check</h3>
        <p className="text-sm text-text-muted mb-4">
          Before declaring MISSION_COMPLETE (the entire deliverable, not just a task), answer ALL of these:
        </p>
        <div className="grid md:grid-cols-2 gap-3">
          {[
            'Is there ANY edge case I thought about but didn\'t handle?',
            'Is there ANY test I thought about writing but didn\'t?',
            'Is there ANY code I\'m not 100% confident in?',
            'Would I confidently ship this to production right now?',
            'If the user tried the most unexpected path, would it still work?',
            'Can I think of a single improvement within scope I haven\'t made?',
          ].map((q, i) => (
            <div key={i} className="flex items-start gap-2 p-3 bg-bg-card border border-border rounded">
              <span className="text-gold text-xs mt-0.5 flex-shrink-0">?</span>
              <p className="text-xs text-text-secondary leading-relaxed">{q}</p>
            </div>
          ))}
        </div>
        <p className="text-xs text-text-muted mt-4 font-mono">
          Only when every answer is "yes / nothing remaining" → output MISSION_COMPLETE
        </p>
      </motion.div>
    </SectionWrapper>
  )
}
