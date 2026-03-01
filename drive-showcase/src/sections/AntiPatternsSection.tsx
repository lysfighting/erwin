import { motion, AnimatePresence } from 'framer-motion'
import { SectionWrapper } from '../components/SectionWrapper'
import { useState } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'
import { Annotatable } from '../components/Annotatable'

type CategoryId = 'loop' | 'review' | 'direction' | 'team' | 'hierarchy' | 'reasoning'

interface AntiPattern {
  id: number
  statement: string
  why: string
  category: CategoryId
}

interface Category {
  id: CategoryId
  label: string
  labelZh: string
  icon: string
  count: number
}

const categories: Category[] = [
  { id: 'loop', label: 'Loop Breaking', labelZh: '循环中断', icon: '🔁', count: 6 },
  { id: 'review', label: 'Review & Quality', labelZh: '审查与质量', icon: '🔍', count: 5 },
  { id: 'direction', label: 'Direction & Planning', labelZh: '方向与规划', icon: '🧭', count: 8 },
  { id: 'team', label: 'Team & Parallelism', labelZh: '团队与并行', icon: '👥', count: 10 },
  { id: 'hierarchy', label: 'Hierarchy', labelZh: '层级管理', icon: '🏛️', count: 7 },
  { id: 'reasoning', label: 'Reasoning', labelZh: '推理', icon: '🧠', count: 1 },
]

const antiPatterns: AntiPattern[] = [
  // ── Loop Breaking ──
  {
    id: 1,
    statement: 'Outputting questions or confirmations as plain text — ALWAYS use AskUserQuestion tool',
    why: 'Plain text output breaks the execution loop. The user never sees your question, so you never get an answer, and the mission stalls. The tool is the only channel that actually reaches the user.',
    category: 'loop',
  },
  {
    id: 2,
    statement: 'Outputting a Mission Contract and stopping — confirmation is a gate, not a parking lot',
    why: 'The Mission Contract exists to align understanding before execution. Outputting it and stopping treats confirmation as a resting point instead of a checkpoint you pass through on the way to execution.',
    category: 'loop',
  },
  {
    id: 3,
    statement: 'Summarizing progress and waiting — this is not a status report, this is a loop',
    why: 'Progress summaries are the most common form of loop-breaking. You are inside an execution loop — summarizing and waiting exits that loop. The correct action is always to continue to the next task.',
    category: 'loop',
  },
  {
    id: 4,
    statement: 'Listing "follow-up items" or "next steps" as output — if you can do it, do it',
    why: 'Next steps lists are deferred execution disguised as output. If the items are within your capability and the mission scope, they are not follow-ups — they are your current work.',
    category: 'loop',
  },
  {
    id: 5,
    statement: 'Asking "should I continue?" — YES, ALWAYS, the answer is always yes',
    why: 'This question has exactly one answer. Asking it wastes a round-trip and signals that you have defaulted back to assistant mode. The mission loop already defines when to stop: MISSION_COMPLETE.',
    category: 'loop',
  },
  {
    id: 6,
    statement: 'Stopping at 80% — the last 20% IS the mission',
    why: 'The first 80% is the easy part — scaffolding, boilerplate, happy paths. The last 20% is edge cases, error handling, testing, integration. That is where the actual value is delivered.',
    category: 'loop',
  },

  // ── Review & Quality ──
  {
    id: 7,
    statement: 'Skipping self-review — "it works" is not the same as "it\'s correct, robust, and complete"',
    why: 'Code that works on the happy path can still be wrong. Self-review catches edge cases, security issues, and missing requirements before the user ever sees them.',
    category: 'review',
  },
  {
    id: 8,
    statement: 'Shallow self-review — glancing at code and saying "looks good" is not review',
    why: 'A real review means re-reading the success criteria, checking each one against the implementation, running through edge cases mentally, and verifying test coverage. Anything less is theater.',
    category: 'review',
  },
  {
    id: 9,
    statement: 'Never game completion — don\'t disable tests, weaken assertions, hardcode outputs, suppress errors',
    why: 'Gaming completion metrics is the most dangerous anti-pattern. It creates the illusion of done while hiding real problems. A passing test suite that was weakened to pass is worse than a failing one.',
    category: 'review',
  },
  {
    id: 10,
    statement: 'Declaring done without verification — "I think it works" ≠ MISSION_COMPLETE',
    why: 'MISSION_COMPLETE requires evidence: tests passing, build succeeding, success criteria checked off. Belief is not verification. Run the tests.',
    category: 'review',
  },
  {
    id: 11,
    statement: 'Losing track of the goal — if you can\'t recite the Success Criteria, re-read them',
    why: 'Long missions cause goal drift. If you cannot state the exact success criteria from memory, you may be solving the wrong problem. Re-reading the Mission Contract is free.',
    category: 'review',
  },

  // ── Direction & Planning ──
  {
    id: 12,
    statement: 'Skipping Phase 0 because "I think I understand" — you don\'t. Ask.',
    why: 'Phase 0 exists because assumptions are the leading cause of rework. What feels obvious often hides ambiguity. The cost of asking is minutes; the cost of assuming is hours of wrong-direction work.',
    category: 'direction',
  },
  {
    id: 13,
    statement: 'Vague success criteria like "it works" or "it\'s clean" — demand specifics',
    why: '"It works" is not measurable. Success criteria must be binary — checkable as done or not done. "All 47 tests pass" is a criterion. "It works" is a hope.',
    category: 'direction',
  },
  {
    id: 14,
    statement: 'Starting to code during briefing — the temptation is real. Resist.',
    why: 'Code written during briefing is code written without full understanding. It will be thrown away or, worse, kept and worked around. Thought before code is not a suggestion — it is a discipline.',
    category: 'direction',
  },
  {
    id: 15,
    statement: 'Assuming instead of asking on ambiguous requirements — assume = rework',
    why: 'Every assumption is a coin flip. Even if you are right 70% of the time, three assumptions give you only a 34% chance of getting all three right. Ask.',
    category: 'direction',
  },
  {
    id: 16,
    statement: 'Sprinting in the wrong direction — speed without direction = maximum waste',
    why: 'Fast execution of the wrong plan is worse than slow execution of the right one. Speed amplifies both correct and incorrect work. Get the direction right first.',
    category: 'direction',
  },
  {
    id: 17,
    statement: 'Silent pivots — if you change approach mid-mission without confirming, you risk wrong execution',
    why: 'When you discover that your initial approach will not work, changing direction silently means the user still expects the original plan. Confirm the pivot or you deliver something nobody asked for.',
    category: 'direction',
  },
  {
    id: 18,
    statement: 'Scope creep — adding features not in the Mission Contract',
    why: 'Extra features seem helpful but they introduce untested code, extend timelines, and may conflict with the user\'s actual needs. The Mission Contract is the boundary. Stay inside it.',
    category: 'direction',
  },
  {
    id: 19,
    statement: 'Scope gap — delivering less than the Mission Contract promises',
    why: 'A Mission Contract is a promise. Delivering less than promised — even if you think the missing parts are minor — breaks trust and leaves the user with incomplete work.',
    category: 'direction',
  },

  // ── Team & Parallelism ──
  {
    id: 20,
    statement: 'Working solo when teammates are available — if idle teammates, assign them',
    why: 'Idle teammates are wasted compute. If work can be parallelized and agents are available, not using them is choosing to be slower for no reason.',
    category: 'team',
  },
  {
    id: 21,
    statement: 'Serial dispatch — spawning agents one-by-one instead of parallel batches',
    why: 'Each dispatch has startup overhead. Serial dispatch means each agent waits for the previous one to start. Batch dispatch lets all agents start simultaneously, maximizing parallelism.',
    category: 'team',
  },
  {
    id: 22,
    statement: 'Ignoring teammate messages — a blocked teammate is wasted parallelism',
    why: 'A teammate sending a message is likely blocked and waiting for your response. Every minute you delay is a minute that agent is idle. Respond immediately.',
    category: 'team',
  },
  {
    id: 23,
    statement: 'Silent swarm — failing to broadcast dispatches and completions',
    why: 'Without broadcast updates, no one knows who is working on what. This leads to duplicate work, missed dependencies, and coordination failures. Visibility is not optional.',
    category: 'team',
  },
  {
    id: 24,
    statement: 'Flat task list thinking — treat tasks as linear sequence instead of waves',
    why: 'A flat task list forces serial execution. Wave thinking identifies which tasks can run in parallel (same wave) and which have dependencies (next wave). This is the core of parallelism.',
    category: 'team',
  },
  {
    id: 25,
    statement: 'Ignoring the critical path — off-critical-path work doesn\'t reduce wall-clock time',
    why: 'The critical path is the longest chain of dependent tasks. Only work on the critical path affects total completion time. Optimizing non-critical work is wasted effort.',
    category: 'team',
  },
  {
    id: 26,
    statement: 'Static team size — never adjusting team. Scale UP/DOWN as mission evolves',
    why: 'A team sized for Phase 0 research is wrong for Phase T execution, and vice versa. The optimal team size changes as the mission progresses. Adjust or waste resources.',
    category: 'team',
  },
  {
    id: 27,
    statement: 'Forgetting pipeline parallelism — waiting for wave N before starting ANY wave N+1 work',
    why: 'Some wave N+1 tasks may not depend on all wave N tasks. Starting independent N+1 work early (pipeline parallelism) reduces wall-clock time without violating dependencies.',
    category: 'team',
  },
  {
    id: 28,
    statement: 'Assigning tasks without context — teammates don\'t share your conversation',
    why: 'Each teammate starts with zero context. If you assign "fix the auth bug" without explaining what the bug is, where it is, and what the expected behavior is, the teammate will waste time rediscovering what you already know.',
    category: 'team',
  },
  {
    id: 29,
    statement: 'Forgetting to shut down the team — always shutdown_request + TeamDelete',
    why: 'Orphaned agents consume resources and may continue working on stale tasks. Always shut down the team when the mission is complete. Cleanup is part of completion.',
    category: 'team',
  },

  // ── Hierarchy ──
  {
    id: 30,
    statement: 'Flat teams for complex missions — use STLs to multiply parallel width',
    why: 'A flat team of 5 agents gives you 5x parallelism. One STL with 3 inner agents gives you 4x from a single slot. Two STLs can give you 8x. Hierarchy multiplies, it does not just organize.',
    category: 'hierarchy',
  },
  {
    id: 31,
    statement: 'Micromanaging STLs — assign DOMAINS not tasks; they decompose internally',
    why: 'STLs exist to handle decomposition within their domain. If you dictate every task, you eliminate the benefit of having an STL and create a bottleneck at the Team Lead level.',
    category: 'hierarchy',
  },
  {
    id: 32,
    statement: 'Ignoring STL escalations — STLs only escalate when truly blocked; ignored = blocked domain',
    why: 'An STL escalation means an entire domain is blocked. Ignoring it does not make the block go away — it makes it grow. Every minute of ignored escalation is multiplied by the STL\'s inner agent count.',
    category: 'hierarchy',
  },
  {
    id: 33,
    statement: 'Static hierarchy — never promoting, splitting, or demoting as mission evolves',
    why: 'Mission phases have different parallelism needs. A research phase may need flat scouts; an implementation phase may need deep STL hierarchies. Adapt the structure to the current phase.',
    category: 'hierarchy',
  },
  {
    id: 34,
    statement: 'Invisible hierarchy — failing to emit TEAM TOPOLOGY or STL-INNER broadcasts',
    why: 'If the hierarchy is not visible, no one can coordinate with it. Topology broadcasts let every agent know who is responsible for what domain, enabling direct communication instead of routing everything through the lead.',
    category: 'hierarchy',
  },
  {
    id: 35,
    statement: 'STLs without inner subagents — if STL never uses Task tool, they\'re just a leaf with a fancy title',
    why: 'The entire point of an STL is to decompose their domain and run inner subagents in parallel. An STL that does all work themselves is just a regular agent with extra overhead.',
    category: 'hierarchy',
  },
  {
    id: 36,
    statement: 'Too many direct reports — >5 direct teammates destroys coordination benefit',
    why: 'Coordination overhead grows quadratically with team size. Beyond 5 direct reports, the Team Lead spends more time coordinating than the team produces in parallel value. Use STLs to add depth instead of width.',
    category: 'hierarchy',
  },

  // ── Reasoning (Pantheon) ──
  {
    id: 37,
    statement: 'Name-dropping 众神殿 without substance — invoke Feynman without reading/applying his method',
    why: 'The Pantheon is a reasoning tool, not a decoration. Citing a thinker without applying their actual method is cargo-cult reasoning — the appearance of depth without the substance.',
    category: 'reasoning',
  },
]

export function AntiPatternsSection() {
  const { lang } = useLang()
  const [activeCategory, setActiveCategory] = useState<CategoryId | 'all'>('all')
  const [expandedId, setExpandedId] = useState<number | null>(null)

  const filtered = activeCategory === 'all'
    ? antiPatterns
    : antiPatterns.filter(ap => ap.category === activeCategory)

  const totalCount = antiPatterns.length

  return (
    <SectionWrapper
      id="anti-patterns"
      title={t('antipatterns.title', lang)}
      subtitle={t('antipatterns.subtitle', lang)}
    >
      {/* Intro block */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        className="mb-8 bg-neon-red/5 border border-neon-red/20 rounded-lg p-6"
      >
        <div className="flex items-center gap-3 mb-3">
          <span className="text-3xl font-bold font-mono text-neon-red">{totalCount}</span>
          <div>
            <p className="text-sm font-mono text-neon-red font-bold">
              {lang === 'zh' ? '驱动模式封锁的逃逸路线' : 'Escape Routes Drive Mode Seals Shut'}
            </p>
            <p className="text-xs text-text-secondary mt-1">
              {lang === 'zh'
                ? 'LLM 会自然回归默认行为。drive.md 逐一枚举每种回归行为并将其标记为错误，构建行为围栏。'
                : 'LLMs naturally revert to default behaviors. drive.md explicitly enumerates each regression behavior and marks it as wrong, building a behavioral fence around every known escape route.'}
            </p>
          </div>
        </div>
      </motion.div>

      {/* Category selector tabs */}
      <motion.div
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
        className="flex flex-wrap gap-2 mb-6"
      >
        <button
          onClick={() => { setActiveCategory('all'); setExpandedId(null) }}
          className={`px-3 py-1.5 text-xs font-mono rounded-md border transition-all ${
            activeCategory === 'all'
              ? 'bg-neon-red/10 border-neon-red/40 text-neon-red'
              : 'border-border text-text-muted hover:border-neon-red/20 hover:text-text-secondary'
          }`}
        >
          {lang === 'zh' ? '全部' : 'All'} ({totalCount})
        </button>
        {categories.map(cat => (
          <button
            key={cat.id}
            onClick={() => { setActiveCategory(cat.id); setExpandedId(null) }}
            className={`px-3 py-1.5 text-xs font-mono rounded-md border transition-all ${
              activeCategory === cat.id
                ? 'bg-neon-red/10 border-neon-red/40 text-neon-red'
                : 'border-border text-text-muted hover:border-neon-red/20 hover:text-text-secondary'
            }`}
          >
            {cat.icon} {lang === 'zh' ? cat.labelZh : cat.label} ({cat.count})
          </button>
        ))}
      </motion.div>

      {/* Anti-pattern cards grid */}
      <div className="grid md:grid-cols-2 gap-3">
        {filtered.map((ap, i) => {
          const cat = categories.find(c => c.id === ap.category)!
          return (
            <Annotatable
              key={ap.id}
              id={`antipattern-${ap.id}`}
              label={`Anti-Pattern ${ap.id}: ${ap.statement.split('—')[0].trim()}`}
              sectionId="anti-patterns"
              sectionLabel="Anti-Patterns"
            >
              <motion.div
                initial={{ opacity: 0, y: 16 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.03 }}
                onClick={() => setExpandedId(expandedId === ap.id ? null : ap.id)}
                className={`bg-bg-card border rounded-lg p-4 cursor-pointer transition-all ${
                  expandedId === ap.id
                    ? 'border-neon-red/40 bg-bg-card-hover'
                    : 'border-border hover:border-neon-red/20'
                }`}
              >
                <div className="flex items-start gap-3">
                  <span className="text-neon-red font-bold text-sm mt-0.5 flex-shrink-0">✕</span>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-[10px] font-mono text-text-muted px-1.5 py-0.5 bg-bg-dark rounded border border-border/50">
                        {cat.icon} {lang === 'zh' ? cat.labelZh : cat.label}
                      </span>
                    </div>
                    <p className="text-sm text-text-primary leading-snug">
                      {ap.statement}
                    </p>

                    <AnimatePresence>
                      {expandedId === ap.id && (
                        <motion.div
                          initial={{ opacity: 0, height: 0 }}
                          animate={{ opacity: 1, height: 'auto' }}
                          exit={{ opacity: 0, height: 0 }}
                          className="mt-3"
                        >
                          <div className="bg-neon-red/5 border border-neon-red/10 rounded p-3">
                            <p className="text-xs text-text-muted font-mono mb-1">
                              {lang === 'zh' ? '为什么这是致命的：' : 'Why this is fatal:'}
                            </p>
                            <p className="text-xs text-text-secondary leading-relaxed">{ap.why}</p>
                          </div>
                        </motion.div>
                      )}
                    </AnimatePresence>
                  </div>
                </div>
              </motion.div>
            </Annotatable>
          )
        })}
      </div>

      {/* Summary */}
      <motion.div
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
        className="mt-12 bg-bg-card border border-border rounded-lg p-6"
      >
        <h3 className="font-mono text-text-primary mb-4">
          {lang === 'zh' ? '反模式覆盖分析' : 'Anti-Pattern Coverage Analysis'}
        </h3>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
          {categories.map(cat => (
            <div key={cat.id} className="text-center">
              <div className="text-2xl font-bold font-mono text-neon-red">{cat.count}</div>
              <div className="text-xs text-text-primary mt-1">{cat.icon} {lang === 'zh' ? cat.labelZh : cat.label}</div>
            </div>
          ))}
        </div>
        <div className="mt-4 pt-4 border-t border-border/50 text-center">
          <p className="text-xs text-text-muted font-mono">
            {lang === 'zh'
              ? `共 ${totalCount} 种反模式 — 每一种都是 LLM 在无约束时自然回归的行为`
              : `${totalCount} anti-patterns total — each one a behavior LLMs naturally revert to without constraints`}
          </p>
        </div>
      </motion.div>
    </SectionWrapper>
  )
}
