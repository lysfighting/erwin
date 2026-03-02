import { motion, AnimatePresence } from 'framer-motion'
import { SectionWrapper } from '../components/SectionWrapper'
import { useState } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'
import { Annotatable } from '../components/Annotatable'

const axioms = [
  {
    id: 1,
    name: 'Ownership Axiom',
    chinese: '所有权公理',
    statement: 'Claude owns the task\'s complete lifecycle, from understanding to delivery.',
    derivations: ['Phase 0 deep briefing', 'Mission Contract', 'Six-Step Discipline'],
  },
  {
    id: 2,
    name: 'Continuity Axiom',
    chinese: '持续性公理',
    statement: 'The execution loop does not break due to output.',
    derivations: ['Cardinal Rule (never stop)', 'AskUserQuestion instead of plain text'],
  },
  {
    id: 3,
    name: 'Parallelism Axiom',
    chinese: '并行性公理',
    statement: 'Work that CAN be parallelized MUST be parallelized.',
    derivations: ['Hierarchical topology', 'STL pattern', 'Wave decomposition', 'Fork-Join/Pipeline/Speculative'],
  },
  {
    id: 4,
    name: 'Visibility Axiom',
    chinese: '可见性公理',
    statement: 'The execution process must be transparent to the user.',
    derivations: ['Broadcast Protocol', 'Team Topology visualization', 'SWARM STATUS dashboard'],
  },
  {
    id: 5,
    name: 'Direction Axiom',
    chinese: '方向优先公理',
    statement: 'Direction correctness takes priority over execution speed.',
    derivations: ['Direction Validation Protocol', 'Self-Dialectic', 'Self-Adversarial Review'],
  },
]

const philosophies = [
  {
    quote: '循环是英雄',
    translation: 'The loop is the hero, not any single iteration',
    meaning: 'Continuous iteration beats single-pass perfection. You will fail, fix, and iterate.',
  },
  {
    quote: '永不停止',
    translation: 'Never stop moving',
    meaning: 'Outputting text and waiting for the user is FAILURE. The only acceptable stops are Mission Contract confirmation and genuine blockers.',
  },
  {
    quote: '蜂群优先',
    translation: 'Swarm first — any work that CAN be parallelized MUST be parallelized',
    meaning: 'Serial execution is the exception, not the default. Maximum parallelism is the operating principle.',
  },
  {
    quote: '量两次，切一次',
    translation: 'Measure twice, cut once',
    meaning: 'Phase 0 deep briefing before touching any code. A 10-minute briefing saves hours of rework.',
  },
  {
    quote: '错误方向全速前进 = 最大浪费',
    translation: 'Wrong direction at full speed = maximum waste',
    meaning: 'A slow walk in the right direction beats a sprint down the wrong path.',
  },
  {
    quote: '一个人走得快，一群人走得远',
    translation: 'One walks fast alone, many walk far together',
    meaning: 'Team assembly is mandatory. Drive Mode always operates as a team.',
  },
  {
    quote: '问得越深，做得越准',
    translation: 'The deeper you ask, the more accurately you execute',
    meaning: 'Every unasked question is a future surprise. Every future surprise is potential rework.',
  },
  {
    quote: '将帅有别，层层拆解',
    translation: 'Generals and captains differ — decompose layer by layer',
    meaning: 'The hierarchical topology: Team Lead doesn\'t manage every worker. Sub-Team-Leads own domains.',
  },
]

const bitterLesson = {
  quote: '"The biggest lesson that can be read from 70 years of AI research is that general methods that leverage computation are ultimately the most effective, and by a large margin."',
  author: 'Richard Sutton, "The Bitter Lesson" (2019)',
  paradox: 'Drive Skill teaches Claude "don\'t write rules, write general methods" — while itself using 112KB of handcrafted rules. This is the self-referential paradox: acceptable as a "cold start solution" when insufficient training data exists for models to learn team orchestration autonomously.',
}

const scalingTabs = ['core', 'questions', 'antipatterns', 'rules-fit'] as const
type ScalingTab = typeof scalingTabs[number]

const scalingTabLabels: Record<ScalingTab, string> = {
  core: 'Core Insight',
  questions: '4 Scaling Questions',
  antipatterns: '5 Anti-patterns',
  'rules-fit': 'When Rules Fit',
}

const scalingQuestions = [
  {
    number: 1,
    question: '"How does this scale with problem complexity?"',
    flag: 'if answer is "add more rules" --> RED FLAG O(2^C)',
  },
  {
    number: 2,
    question: '"Am I encoding WHAT or HOW?"',
    flag: 'encoding WHAT (objectives) = good; HOW (procedures) = suspicious',
  },
  {
    number: 3,
    question: '"Where does the intelligence live?"',
    flag: 'in code (rules) = fragile; in data/weights = robust',
  },
  {
    number: 4,
    question: '"What\'s the loss function?"',
    flag: 'every well-posed problem has one; if you can\'t state it, you don\'t understand the problem',
  },
]

const antiPatterns = [
  { bad: 'Enumerating cases', good: 'Define objective, let system find mapping' },
  { bad: 'Hardcoding heuristics', good: 'Make it learnable: threshold = f(context)' },
  { bad: 'Taxonomizing upfront', good: 'Define error metric, learn the handler' },
  { bad: 'Config-driven complexity (47 parameters)', good: 'Learning algorithm finds good config' },
  { bad: 'Template proliferation', good: 'Find generative model that produces templates' },
]

const rulesAppropriate = [
  'Problem space genuinely small (<=10 cases)',
  'Correctness provably necessary (safety constraints, invariants)',
  'The rule encodes a LAW (physics, math, logic) not a HEURISTIC',
  'No data to learn from (cold start)',
]

function ScalingMindsetSection() {
  const [activeTab, setActiveTab] = useState<ScalingTab>('core')

  return (
    <motion.div
      initial={{ opacity: 0 }}
      whileInView={{ opacity: 1 }}
      viewport={{ once: true }}
      className="bg-bg-card border border-neon-red/20 rounded-lg p-8"
    >
      <h3 className="font-mono text-neon-red text-lg mb-4">
        The Scaling Mindset — 苦涩教训 (The Bitter Lesson)
      </h3>

      {/* Original Bitter Lesson quote as intro */}
      <blockquote className="text-text-secondary italic border-l-2 border-neon-red/30 pl-4 mb-4">
        {bitterLesson.quote}
      </blockquote>
      <p className="text-xs text-text-muted mb-4">— {bitterLesson.author}</p>
      <div className="bg-neon-red/5 border border-neon-red/10 rounded p-4 mb-6">
        <p className="text-xs text-text-muted font-mono mb-1">THE SELF-REFERENTIAL PARADOX:</p>
        <p className="text-sm text-text-secondary">{bitterLesson.paradox}</p>
      </div>

      {/* Tab navigation */}
      <div className="flex flex-wrap gap-2 mb-6 border-b border-border/50 pb-3">
        {scalingTabs.map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`px-3 py-1.5 rounded-t font-mono text-xs transition-colors ${
              activeTab === tab
                ? 'bg-neon-red/20 text-neon-red border border-neon-red/30 border-b-0'
                : 'text-text-muted hover:text-text-secondary'
            }`}
          >
            {scalingTabLabels[tab]}
          </button>
        ))}
      </div>

      {/* Tab content */}
      <AnimatePresence mode="wait">
        {activeTab === 'core' && (
          <motion.div
            key="core"
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.2 }}
          >
            {/* Two-column comparison */}
            <div className="grid md:grid-cols-2 gap-4 mb-6">
              <div className="bg-neon-red/5 border border-neon-red/20 rounded-lg p-4">
                <p className="font-mono text-xs text-neon-red font-bold mb-2">APPROACH A (Rules)</p>
                <p className="font-mono text-xs text-text-secondary leading-relaxed">
                  Engineer enumerates cases &rarr; writes rules for each &rarr; system follows rules
                </p>
                <div className="mt-3 pt-3 border-t border-neon-red/10">
                  <p className="font-mono text-xs text-text-muted">Scaling behavior:</p>
                  <p className="font-mono text-sm text-neon-red mt-1">
                    performance = f(human_effort)
                  </p>
                  <p className="font-mono text-xs text-text-muted mt-1">
                    bounded by human ingenuity
                  </p>
                </div>
              </div>
              <div className="bg-green-500/5 border border-green-500/20 rounded-lg p-4">
                <p className="font-mono text-xs text-green-400 font-bold mb-2">APPROACH B (Learning)</p>
                <p className="font-mono text-xs text-text-secondary leading-relaxed">
                  Engineer defines objective + provides data &rarr; system discovers the rules
                </p>
                <div className="mt-3 pt-3 border-t border-green-500/10">
                  <p className="font-mono text-xs text-text-muted">Scaling behavior:</p>
                  <p className="font-mono text-sm text-green-400 mt-1">
                    performance = f(data x compute)
                  </p>
                  <p className="font-mono text-xs text-text-muted mt-1">
                    bounded by resources, which scale
                  </p>
                </div>
              </div>
            </div>

            {/* O-notation comparison card */}
            <div className="bg-bg-dark border border-border/50 rounded-lg p-5">
              <p className="font-mono text-xs text-text-muted mb-3">As problem complexity C grows:</p>
              <div className="grid md:grid-cols-2 gap-4">
                <div className="flex items-center gap-3">
                  <span className="font-mono text-2xl font-bold text-neon-red">O(2^C)</span>
                  <div>
                    <p className="font-mono text-xs text-neon-red">Rules</p>
                    <p className="font-mono text-xs text-text-muted">hand-written cases &mdash; combinatorial explosion</p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <span className="font-mono text-2xl font-bold text-green-400">O(C&middot;log&nbsp;C)</span>
                  <div>
                    <p className="font-mono text-xs text-green-400">Learning</p>
                    <p className="font-mono text-xs text-text-muted">data needed &mdash; tractable</p>
                  </div>
                </div>
              </div>
            </div>
          </motion.div>
        )}

        {activeTab === 'questions' && (
          <motion.div
            key="questions"
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.2 }}
            className="space-y-3"
          >
            <p className="font-mono text-xs text-text-muted mb-4">Ask before any design decision:</p>
            {scalingQuestions.map((q) => (
              <div key={q.number} className="bg-bg-dark border border-border/30 rounded-lg p-4">
                <div className="flex items-start gap-3">
                  <span className="font-mono text-lg font-bold text-gold">{q.number}</span>
                  <div>
                    <p className="font-mono text-sm text-text-primary mb-1">{q.question}</p>
                    <p className="font-mono text-xs text-text-muted">{q.flag}</p>
                  </div>
                </div>
              </div>
            ))}
          </motion.div>
        )}

        {activeTab === 'antipatterns' && (
          <motion.div
            key="antipatterns"
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.2 }}
            className="space-y-3"
          >
            <p className="font-mono text-xs text-text-muted mb-4">5 Anti-patterns of Rule Obsession (with scaling alternatives):</p>
            {antiPatterns.map((ap, i) => (
              <div key={i} className="bg-bg-dark border border-border/30 rounded-lg p-4">
                <div className="flex items-start gap-3 mb-2">
                  <span className="font-mono text-sm font-bold text-text-muted">{i + 1}.</span>
                  <div className="flex-1 space-y-2">
                    <div className="flex items-start gap-2">
                      <span className="text-neon-red font-mono text-sm shrink-0">X</span>
                      <p className="font-mono text-sm text-neon-red/80">{ap.bad}</p>
                    </div>
                    <div className="flex items-start gap-2">
                      <span className="text-green-400 font-mono text-sm shrink-0">&#10003;</span>
                      <p className="font-mono text-sm text-green-400/80">{ap.good}</p>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </motion.div>
        )}

        {activeTab === 'rules-fit' && (
          <motion.div
            key="rules-fit"
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.2 }}
          >
            <p className="font-mono text-xs text-text-muted mb-4">When rules ARE the right choice:</p>
            <div className="space-y-2 mb-6">
              {rulesAppropriate.map((rule, i) => (
                <div key={i} className="bg-bg-dark border border-border/30 rounded-lg p-3 flex items-start gap-3">
                  <span className="text-gold font-mono text-sm mt-0.5">&#8226;</span>
                  <p className="font-mono text-sm text-text-secondary">{rule}</p>
                </div>
              ))}
            </div>
            <div className="bg-gold/5 border border-gold/20 rounded-lg p-4">
              <p className="font-mono text-xs text-gold font-bold mb-2">THE LITMUS TEST:</p>
              <p className="font-mono text-sm text-text-secondary">
                "Will this rule still be correct in 1 year when the system has grown 10x?"
              </p>
              <div className="grid md:grid-cols-2 gap-3 mt-3">
                <div className="bg-green-500/5 border border-green-500/20 rounded p-2">
                  <p className="font-mono text-xs text-green-400">If yes = LAW (keep)</p>
                </div>
                <div className="bg-neon-red/5 border border-neon-red/20 rounded p-2">
                  <p className="font-mono text-xs text-neon-red">If no = HEURISTIC (make learnable)</p>
                </div>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  )
}

export function PhilosophySection() {
  const { lang } = useLang()
  const [expandedAxiom, setExpandedAxiom] = useState<number | null>(null)

  return (
    <SectionWrapper
      id="philosophy"
      title={t('philosophy.title', lang)}
      subtitle={t('philosophy.subtitle', lang)}
    >
      {/* Five Axioms */}
      <h3 className="font-mono text-lg text-gold mb-6">The Five Axioms</h3>
      <div className="space-y-3 mb-16">
        {axioms.map((axiom, i) => (
          <Annotatable
            key={axiom.id}
            id={`axiom-${axiom.id}`}
            label={`Axiom A${axiom.id}: ${axiom.name}`}
            sectionId="philosophy"
            sectionLabel="Philosophy"
          >
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            transition={{ delay: i * 0.08 }}
            onClick={() => setExpandedAxiom(expandedAxiom === i ? null : i)}
            className="bg-bg-card border border-border rounded-lg p-5 cursor-pointer hover:border-gold/30 transition-colors"
          >
            <div className="flex items-center gap-3">
              <span className="font-mono text-gold font-bold text-lg">A{axiom.id}</span>
              <span className="font-mono text-sm text-text-primary">{axiom.name}</span>
              <span className="text-text-muted text-xs ml-1">({axiom.chinese})</span>
            </div>
            <p className="text-sm text-text-secondary mt-2">{axiom.statement}</p>
            {expandedAxiom === i && (
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                className="mt-3 pt-3 border-t border-border/50"
              >
                <span className="text-xs text-text-muted font-mono">Derives →</span>
                <div className="flex flex-wrap gap-2 mt-1">
                  {axiom.derivations.map((d) => (
                    <span key={d} className="text-xs bg-gold/10 text-gold px-2 py-1 rounded font-mono">
                      {d}
                    </span>
                  ))}
                </div>
              </motion.div>
            )}
          </motion.div>
          </Annotatable>
        ))}
      </div>

      {/* Core Beliefs */}
      <h3 className="font-mono text-lg text-text-primary mb-6">Core Beliefs — 核心理念</h3>
      <div className="grid md:grid-cols-2 gap-4 mb-16">
        {philosophies.map((p, i) => (
          <Annotatable
            key={p.quote}
            id={`belief-${i}`}
            label={`Belief: ${p.translation}`}
            sectionId="philosophy"
            sectionLabel="Philosophy"
          >
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.05 }}
              className="bg-bg-card border border-border rounded-lg p-5"
            >
              <p className="text-lg font-bold text-gold mb-1">"{p.quote}"</p>
              <p className="text-sm text-text-secondary mb-2">{p.translation}</p>
              <p className="text-xs text-text-muted">{p.meaning}</p>
            </motion.div>
          </Annotatable>
        ))}
      </div>

      {/* The Scaling Mindset — expanded from The Bitter Lesson */}
      <ScalingMindsetSection />
    </SectionWrapper>
  )
}
