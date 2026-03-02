import { motion, useInView, AnimatePresence } from 'framer-motion'
import { useRef, useState } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

interface PhaseDetail {
  title: string
  content: string[]
}

interface Phase {
  id: string
  name: string
  chinese: string
  tagline: string
  color: string
  steps: string[]
  keyInsight: string
  without: string
  details?: PhaseDetail[]
}

const phases: Phase[] = [
  {
    id: '0',
    name: 'Deep Briefing',
    chinese: '深度摸底',
    tagline: 'Measure twice, cut once',
    color: '#06b6d4',
    steps: [
      'Launch Research Swarm (3-5 parallel scouts)',
      'The Interrogation — 8 dimensions, multiple rounds',
      'Lock the Mission Contract',
      'IMMEDIATELY proceed (do NOT stop here)',
    ],
    keyInsight: 'The Psychiatrist Test: Can you write a paragraph describing the user\'s vision AS IF YOU WERE THEM?',
    without: 'You build the wrong thing at full speed. A 10-minute briefing prevents hours of rework. Every unasked question is a future surprise.',
    details: [
      {
        title: 'The 8 Interrogation Dimensions',
        content: [
          '1. Success Definition — "What does done look like?" (MOST CRITICAL)',
          '2. User Mental Model — "What\'s the shape in your head?"',
          '3. Scope Boundaries — "Where do we stop?"',
          '4. Constraints — "What are the red lines?"',
          '5. Decision Points — "Which choices need your input?"',
          '6. Priority — "What\'s non-negotiable?"',
          '7. Edge Cases — "What happens when things go wrong?"',
          '8. Context & Motivation — "Why are we doing this?"',
        ],
      },
      {
        title: 'The Psychiatrist Test',
        content: [
          'After all questioning rounds, can you write a paragraph describing the user\'s vision AS IF YOU WERE THEM?',
          'If not — you haven\'t asked enough questions. Go back and ask more.',
          'This is your litmus test for understanding depth.',
        ],
      },
      {
        title: 'Multi-Round Questioning Strategy',
        content: [
          'Round 1: Broad discovery — cover all 8 dimensions',
          'Round 2: Depth on gaps — drill into vague answers',
          'Round 3: Confirmation — play back your understanding',
          'Only lock the Mission Contract when you pass the Psychiatrist Test.',
        ],
      },
    ],
  },
  {
    id: 'T',
    name: 'Team Assembly',
    chinese: '团队组建',
    tagline: '一个人走得快，一群人走得远',
    color: '#a855f7',
    steps: [
      'Create team via TeamCreate',
      'Assess team composition — size by parallelism width',
      'Spawn teammates (STLs for ≥3 task domains, leaves for ≤2)',
      'Emit 🗺️ TEAM TOPOLOGY visualization',
    ],
    keyInsight: 'direct_teammates = min(parallelism_width, 5). STL multiplier: 1 STL + 3 inner = 4x throughput.',
    without: 'One agent, one context window, serial execution. A 30-task mission becomes a sequential slog. The swarm advantage evaporates.',
    details: [
      {
        title: 'The STL Multiplier Formula',
        content: [
          'effective_parallelism = teammates + Σ(inner_subagents_per_STL)',
          '',
          'Example: 2 STLs + 1 leaf + tester',
          '  STL-frontend: 3 inner subagents',
          '  STL-backend:  2 inner subagents',
          '  leaf-docs:    0 inner (leaf)',
          '  tester:       0 inner (leaf)',
          '  ─────────────────────────────',
          '  effective_parallelism = 4 + 3 + 2 = 9 streams',
          '',
          'Key: STLs are Agent tool calls (visible panes).',
          'Inner subagents are Task tool calls (invisible but fast).',
        ],
      },
    ],
  },
  {
    id: '1',
    name: 'Battle Plan',
    chinese: '战斗计划',
    tagline: 'Think in waves, not lists',
    color: '#3b82f6',
    steps: [
      'Task decomposition with parallelism analysis',
      'Wave decomposition — map tasks to parallel waves',
      'Critical path analysis — shorten the longest chain',
      'Assign tasks to teammates, launch Wave 0',
    ],
    keyInsight: 'Wall-clock time = Critical Path Length / Parallelism Width. Your job: minimize path, maximize width.',
    without: 'Tasks get assigned ad-hoc, dependencies get missed, critical path is unknown. Parallelism is accidental, not engineered.',
    details: [
      {
        title: 'Wave Decomposition Algorithm',
        content: [
          '1. List all tasks with their dependencies',
          '2. Tasks with ZERO dependencies → Wave 0',
          '3. For each remaining task:',
          '   wave(task) = max(wave(dep) for dep in task.deps) + 1',
          '4. Tasks in the same wave can run in PARALLEL',
          '5. Tasks in different waves must run SEQUENTIALLY',
          '',
          'Example:',
          '  Wave 0: [setup-db, init-project, design-api]  → parallel',
          '  Wave 1: [build-models, build-routes]          → after Wave 0',
          '  Wave 2: [integration-tests, deploy-staging]   → after Wave 1',
          '',
          'Critical Path = longest chain through waves.',
          'Optimization: break large tasks to widen waves.',
        ],
      },
    ],
  },
  {
    id: 'L',
    name: 'The Loop',
    chinese: '蜂群执行',
    tagline: 'The loop does not break',
    color: '#22c55e',
    steps: [
      'Wave management — dispatch, track, launch next wave',
      'Execute your own critical-path tasks',
      'Orchestrate swarm — respond to STL messages, unblock',
      'Dynamic scaling — promote, split, reassign as needed',
      'Checkpoint every 3 tasks — direction health check',
    ],
    keyInsight: 'while (mission != COMPLETE) { wave_check → execute → verify → orchestrate → scaling_check }',
    without: 'Without the loop, you sprint 80% and stop. Error handling, edge cases, integration — the "last 20%" that actually makes software work — gets skipped.',
    details: [
      {
        title: 'The Execution Loop (Pseudocode)',
        content: [
          'while (mission != COMPLETE) {',
          '  // 1. Wave Management',
          '  if (current_wave.all_complete) {',
          '    verify_wave_outputs();',
          '    dispatch_next_wave();',
          '  }',
          '',
          '  // 2. Execute your own tasks',
          '  task = pick_critical_path_task();',
          '  execute(task);',
          '',
          '  // 3. Orchestrate',
          '  for (msg in pending_messages) {',
          '    if (msg.is_blocker) unblock(msg);',
          '    if (msg.is_completion) verify_and_accept(msg);',
          '  }',
          '',
          '  // 4. Scaling check',
          '  if (next_wave.width > team.capacity)',
          '    scale_up();',
          '  if (teammate.idle && no_future_tasks)',
          '    scale_down(teammate);',
          '',
          '  // 5. Checkpoint (every 3 tasks)',
          '  if (completed_since_last_checkpoint >= 3)',
          '    direction_health_check();',
          '}',
        ],
      },
      {
        title: 'The Relentless Rule',
        content: [
          '"If you can think of it, you can do it.',
          ' If you can do it, you MUST do it."',
          '',
          '- If you think "I should also handle X edge case"',
          '  --> handle it right now',
          '- If you think "this could be more robust if I added Y"',
          '  --> add Y right now',
          '- If you think "a follow-up task could address Z"',
          '  --> Z is not a follow-up, Z is NOW',
          '- If you think "the user might also want W" and W is in scope',
          '  --> do W right now',
          '',
          'ONLY things that go in "follow-up" are things GENUINELY',
          'OUT OF SCOPE per Mission Contract.',
          '',
          'You do not stop until you hit the theoretical ceiling.',
        ],
      },
      {
        title: 'Checkpoint Protocol (every 3 tasks)',
        content: [
          '1. Re-read the Mission Contract tracking task',
          '2. List which Success Criteria are now satisfied vs. remaining',
          '3. Direction health check:',
          '   - "Can I state in one sentence how my current work',
          '      connects to the Mission?"',
          '   - "Am I fighting the code or flowing with it?"',
          '     (fighting = possible direction error)',
          '   - "Has anything I\'ve learned invalidated an',
          '      earlier assumption?"',
          '4. If any answer is concerning',
          '   --> trigger Direction Validation Protocol',
          '5. Check for scope drift',
        ],
      },
      {
        title: 'Staying Effective in Long Sessions',
        content: [
          'STL teammates (in-process, visible)',
          '  --> for sustained multi-task domains',
          'Leaf teammates (in-process, visible)',
          '  --> for individual tasks',
          'Task(Explore) subagents (background)',
          '  --> for quick fire-and-forget lookups',
          '',
          'Small steps, fast feedback',
          '  --> don\'t write 500 lines then test',
          '',
          'When confused, re-read TaskList and Mission Contract.',
        ],
      },
      {
        title: 'The Escalation Ladder',
        content: [
          'When blocked, try in this order:',
          '',
          '1. Try a different approach',
          '   --> at least 3 alternatives to the same sub-problem',
          '2. Research',
          '   --> WebSearch, docs, source code.',
          '       The answer is usually findable.',
          '3. Reduce scope locally',
          '   --> simplify the sub-problem, not the mission',
          '4. Skip and continue',
          '   --> document as blocked task, move to next task',
          '5. Escalate -- LAST RESORT',
          '   --> only if 3+ approaches failed, or missing',
          '       credentials, or contradicting requirements,',
          '       or irreversible high-stakes action',
          '',
          'When escalating: say what you tried (specifically),',
          'what failed (with evidence), what you\'d do next',
          'given a choice. Resume IMMEDIATELY when unblocked.',
        ],
      },
    ],
  },
  {
    id: '✓',
    name: 'Mission Complete',
    chinese: '任务完成',
    tagline: 'Not done until verified',
    color: '#f59e0b',
    steps: [
      'Graceful team shutdown (shutdown_request to all)',
      'Self-Adversarial Review of entire deliverable',
      'Verify EVERY Success Criterion with evidence',
      'Pre-Completion Exhaustion Check',
      'MISSION_COMPLETE + structured debrief',
    ],
    keyInsight: 'If you can think of a single improvement within scope that you haven\'t made — you\'re not done.',
    without: 'Tasks are "done" when you stop working on them, not when they\'re verified. The 20% you didn\'t check is the 20% that fails in production.',
    details: [
      {
        title: '8 Conditions for MISSION_COMPLETE',
        content: [
          '1. Every task is marked complete (check with TaskList)',
          '2. All teammates shut down gracefully',
          '   (shutdown_request confirmed)',
          '3. Team resources cleaned up (TeamDelete)',
          '4. Final verification passes (build, tests)',
          '5. You have checked EACH Success Criterion --',
          '   literally verified by running/testing',
          '6. You can demonstrate concrete evidence',
          '   for each criterion',
          '7. Final Self-Adversarial Review of ENTIRE',
          '   deliverable (not just individual tasks)',
          '8. You cannot think of a single improvement',
          '   within scope you haven\'t already made',
        ],
      },
      {
        title: 'Pre-Completion Exhaustion Check',
        content: [
          'Before MISSION_COMPLETE, answer ALL of these:',
          '',
          '• Is there ANY edge case I thought about but didn\'t handle?',
          '• Is there ANY test I thought about writing but didn\'t?',
          '• Is there ANY code I\'m not 100% confident in?',
          '• Would I confidently ship this to production right now?',
          '• If the user tried the most unexpected way possible,',
          '  would it still work?',
          '',
          'If ANY answer is "yes" or "no" (respectively) --',
          'you are NOT done. Go back and fix it.',
        ],
      },
      {
        title: 'Team Shutdown Protocol',
        content: [
          '1. Check TaskList -- all tasks must be completed',
          '2. BROADCAST: "ALL WAVES COMPLETE --',
          '   Initiating team shutdown"',
          '3. For each teammate:',
          '   SendMessage type: "shutdown_request"',
          '4. BROADCAST: "SHUTDOWN: [agent-name] --',
          '   tasks delivered"',
          '5. Wait for shutdown confirmations',
          '6. Review all teammate work products one final time',
          '7. TeamDelete to clean up team infrastructure',
          '8. Proceed to final verification (solo)',
        ],
      },
    ],
  },
]

function PhaseDeepDive({ details, color, isOpen }: { details: PhaseDetail[]; color: string; isOpen: boolean }) {
  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ height: 0, opacity: 0 }}
          animate={{ height: 'auto', opacity: 1 }}
          exit={{ height: 0, opacity: 0 }}
          transition={{ duration: 0.3, ease: 'easeInOut' }}
          className="overflow-hidden"
        >
          <div className="mt-4 pt-4 border-t border-border/50 space-y-4">
            {details.map((detail, i) => (
              <motion.div
                key={detail.title}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.1 }}
                className="bg-bg-dark rounded-lg p-4 border border-border/30"
              >
                <h4
                  className="font-mono text-sm font-bold mb-3"
                  style={{ color }}
                >
                  {detail.title}
                </h4>
                <div className="font-mono text-xs leading-relaxed space-y-0.5">
                  {detail.content.map((line, li) => (
                    <div
                      key={li}
                      className={line === '' ? 'h-2' : 'text-text-secondary'}
                    >
                      {line}
                    </div>
                  ))}
                </div>
              </motion.div>
            ))}
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}

export function ExecutionPhasesSection() {
  const { lang } = useLang()
  const ref = useRef(null)
  const isInView = useInView(ref, { once: true, margin: '-50px' })
  const [expandedPhases, setExpandedPhases] = useState<Set<string>>(new Set())

  const togglePhase = (id: string) => {
    setExpandedPhases(prev => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      return next
    })
  }

  return (
    <section id="phases" className="py-24 px-4">
      <div className="max-w-6xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 40 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
        >
          <h2 className="text-4xl font-bold mb-2">{t('phases.title', lang)}</h2>
          <p className="text-text-secondary text-lg mb-12">
            {t('phases.subtitle', lang)}
          </p>
        </motion.div>

        <div ref={ref} className="relative">
          {/* Vertical timeline line */}
          <motion.div
            initial={{ height: 0 }}
            animate={isInView ? { height: '100%' } : {}}
            transition={{ duration: 1.5, ease: 'easeOut' }}
            className="absolute left-6 md:left-8 top-0 w-0.5 bg-border"
          />

          <div className="space-y-8">
            {phases.map((phase, i) => (
              <motion.div
                key={phase.id}
                initial={{ opacity: 0, x: 30 }}
                animate={isInView ? { opacity: 1, x: 0 } : {}}
                transition={{ delay: i * 0.2, duration: 0.5 }}
                className="relative pl-16 md:pl-20"
              >
                {/* Phase marker */}
                <div
                  className="absolute left-3 md:left-5 w-7 h-7 rounded-full flex items-center justify-center text-xs font-mono font-bold border-2"
                  style={{
                    backgroundColor: phase.color + '20',
                    borderColor: phase.color,
                    color: phase.color,
                  }}
                >
                  {phase.id}
                </div>

                {/* Phase content */}
                <div className="bg-bg-card border border-border rounded-lg p-6 hover:bg-bg-card-hover transition-colors"
                  style={{ borderLeftColor: phase.color, borderLeftWidth: 3 }}
                >
                  <div className="flex flex-col md:flex-row md:items-center gap-2 mb-3">
                    <h3 className="font-mono font-bold text-lg" style={{ color: phase.color }}>
                      Phase {phase.id}: {phase.name}
                    </h3>
                    <span className="text-text-muted text-xs">({phase.chinese})</span>
                    <span className="text-text-muted text-xs italic ml-auto">"{phase.tagline}"</span>
                  </div>

                  <div className="grid md:grid-cols-2 gap-4">
                    <div>
                      <ul className="space-y-1.5">
                        {phase.steps.map((step, si) => (
                          <motion.li
                            key={si}
                            initial={{ opacity: 0 }}
                            animate={isInView ? { opacity: 1 } : {}}
                            transition={{ delay: i * 0.2 + si * 0.05 + 0.3 }}
                            className="flex items-start gap-2 text-sm text-text-secondary"
                          >
                            <span className="text-text-muted mt-1">•</span>
                            {step}
                          </motion.li>
                        ))}
                      </ul>
                    </div>
                    <div>
                      <div className="bg-bg-dark rounded p-3 border border-border/50">
                        <p className="text-xs text-text-muted font-mono mb-1">Key Insight:</p>
                        <p className="text-xs text-gold/80 font-mono">{phase.keyInsight}</p>
                      </div>
                      <div className="bg-neon-red/5 border border-neon-red/10 rounded p-3 mt-2">
                        <p className="text-xs text-text-muted font-mono mb-1">WITHOUT THIS PHASE:</p>
                        <p className="text-xs text-neon-red/80 font-mono">{phase.without}</p>
                      </div>
                    </div>
                  </div>

                  {/* Deep dive toggle */}
                  {phase.details && phase.details.length > 0 && (
                    <>
                      <button
                        onClick={() => togglePhase(phase.id)}
                        className="mt-4 flex items-center gap-2 text-xs font-mono transition-colors hover:brightness-125"
                        style={{ color: phase.color }}
                      >
                        <motion.span
                          animate={{ rotate: expandedPhases.has(phase.id) ? 90 : 0 }}
                          transition={{ duration: 0.2 }}
                          className="inline-block"
                        >
                          ▶
                        </motion.span>
                        {expandedPhases.has(phase.id) ? 'Hide details' : 'Show details'}
                      </button>
                      <PhaseDeepDive
                        details={phase.details}
                        color={phase.color}
                        isOpen={expandedPhases.has(phase.id)}
                      />
                    </>
                  )}
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
