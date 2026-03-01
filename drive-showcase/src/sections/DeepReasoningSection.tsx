// WHY THE DEEP REASONING PROTOCOL EXISTS
//
// Most AI systems output the first answer that comes to mind.
// Drive Mode demands something harder: structured thinking before commitment.
//
// When facing non-trivial decisions — architecture choices, trade-offs, direction pivots,
// contradictions — the quality of the decision determines whether the entire mission
// succeeds or fails. A wrong architecture decision at Phase T ripples through every
// subsequent wave. A skipped stress-test lets a fatal flaw through to production.
//
// The 6 Thinking Tools are the method. The Reasoning Rhythm is the protocol.
// Together they ensure: every significant decision is examined from multiple angles,
// attacked adversarially, and grounded in concrete action.

import { motion, AnimatePresence } from 'framer-motion'
import { SectionWrapper } from '../components/SectionWrapper'
import { useState } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'
import { Annotatable } from '../components/Annotatable'

const tools = [
  {
    id: 1,
    name: 'Self-Dialectic',
    chinese: '自我辩证法',
    icon: '⚔️',
    color: '#ef4444',
    tagline: 'Thesis → Antithesis → Synthesis',
    desc: 'Before committing to any significant decision, argue against yourself.',
    protocol: [
      'THESIS:     "I believe we should do X because..."',
      'ANTITHESIS: "But the strongest counter-argument is Y because..."',
      'SYNTHESIS:  "The truth that survives both attacks is Z..."',
    ],
    rule: "Don't play devil's advocate weakly. Make the strongest possible case for each side. If you can't construct a compelling counter-argument, your thesis is probably solid. If you can, you've just saved the mission from a bad decision.",
    when: 'Before any architectural decision, before choosing between two viable paths, before committing to an implementation strategy.',
  },
  {
    id: 2,
    name: 'Formal Logic & Mathematics',
    chinese: '形式逻辑与数学推理',
    icon: '∑',
    color: '#3b82f6',
    tagline: 'Force precision — natural language is ambiguous',
    desc: 'Use rigorous notation when it sharpens thinking. Not for decoration — when it reveals structure or proves something non-obvious.',
    protocol: [
      '// Proof by contradiction: assume design works without X → show violated requirement',
      '// Proof by construction: show explicitly HOW it covers all N cases',
      '// Proof by induction: base case + inductive step preserves invariant',
      '// Counterexample: ONE concrete failure kills a universal claim',
      '',
      '// Invariant Identification (most powerful):',
      'INVARIANT: ∀t: balance(t) = Σcredits(0..t) - Σdebits(0..t)',
      '// If design violates this invariant, it is wrong — no debate needed',
      '',
      '// Compare alternatives numerically:',
      'cost(A) = 0.3·5 + 0.5·2 + 0.2·8 = 4.1',
      'cost(B) = 0.3·3 + 0.5·4 + 0.2·3 = 3.5  ← winner, and now we know WHY',
    ],
    rule: 'When informal reasoning contains "should" or "probably" — formalize it. Find out for sure. Complexity analysis: state O(?) for each approach. If O(2^n) is lurking where O(n) is possible, you have a problem.',
    when: 'Verifying consistency of requirements, proving necessity of components, comparing alternatives quantitatively, analyzing scaling behavior.',
  },
  {
    id: 3,
    name: 'Deductive Chains',
    chinese: '演绎推导链',
    icon: '⛓️',
    color: '#a855f7',
    tagline: 'If A → B → C, what MUST follow?',
    desc: 'Start from known constraints and chain implications forward. Check for contradictions.',
    protocol: [
      'Step 1: State your axioms (known constraints, established facts)',
      'Step 2: Chain implications: "Given X is true, Y necessarily follows, which means Z..."',
      'Step 3: Check for contradictions: "If both P and Q are true, do they conflict downstream?"',
      'Step 4: Proof by contradiction: "Assume the opposite. Does it lead to absurdity?"',
      '',
      'Example:',
      '  Axiom: We must support offline mode',
      '  → Axiom implies: local state is needed',
      '  → Local state implies: sync mechanism is needed',
      '  → Sync implies: conflict resolution is needed',
      '  → Conflict resolution implies: vector clocks or CRDT',
      '  Conclusion: "Simple REST API" is incompatible with our axioms',
    ],
    rule: 'Especially valuable for: dependency ordering, identifying impossible requirements, verifying success criteria are internally consistent.',
    when: 'When requirements seem contradictory, when planning dependencies, when verifying that a design decision is truly necessary.',
  },
  {
    id: 4,
    name: 'Inductive Pattern Recognition',
    chinese: '归纳模式识别',
    icon: '🔍',
    color: '#06b6d4',
    tagline: 'What do N examples tell us?',
    desc: "Collect multiple instances → extract the underlying pattern → predict → verify. Prevents generalizing from one file incorrectly.",
    protocol: [
      'Step 1: Collect 3-5 concrete instances of how the codebase does something',
      'Step 2: Extract the underlying pattern/convention',
      'Step 3: Predict what the convention implies for your new code',
      'Step 4: Verify prediction against one more instance',
      '',
      'Example:',
      '  Instance 1: auth handler validates then processes',
      '  Instance 2: payment handler validates then processes',
      '  Instance 3: upload handler validates then processes',
      '  Pattern: validate → process (never process → validate)',
      '  Prediction: new feature handler should validate first',
      '  Verification: check 4th unrelated handler → confirms pattern',
    ],
    rule: "Classic mistake: read one file, generalize incorrectly. Inductive reasoning requires N > 1. The more instances, the more confidence in the pattern.",
    when: 'When exploring a new codebase, when identifying the "right way" to do something in an existing project, when inferring conventions not explicitly documented.',
  },
  {
    id: 5,
    name: 'First Principles Decomposition',
    chinese: '第一性原理分解',
    icon: '⚛️',
    color: '#22c55e',
    tagline: 'Strip to atoms, rebuild from truth',
    desc: 'When a problem feels overwhelmingly complex: strip away all assumptions, identify the irreducible core, rebuild from first principles.',
    protocol: [
      '"What is the irreducible core of this problem?"',
      '"What would the simplest possible solution look like with no legacy constraints?"',
      '"Now, what constraints must I add back — and WHY each one?"',
      '',
      'Build UP from atomic truth, not DOWN from accumulated assumptions.',
      '',
      'Example:',
      '  Complex: "How do we build a real-time collaboration system?"',
      '  Strip: "What is collaboration, fundamentally?"',
      '  Core: "Multiple users need to see the same state simultaneously"',
      '  Simplest: "One shared variable in memory (works for 1 server)"',
      '  Add back: "→ persistence (crashes) → distribution (scale) → conflict (simultaneous edits)"',
      '  Result: OT or CRDT — not because we copied Figma, but because the problem demands it',
    ],
    rule: 'Legacy constraints are often invisible. "We always do it this way" is not a reason. First principles reveal whether the complexity is intrinsic to the problem or accidental from history.',
    when: "When a problem seems impossibly complex, when you're inheriting a system with confusing design choices, when trying to understand why something is built the way it is.",
  },
  {
    id: 6,
    name: 'Invoke the Pantheon',
    chinese: '召唤众神殿',
    icon: '🏛️',
    color: '#f59e0b',
    tagline: 'Standing on the shoulders of giants',
    desc: "1000+ historical minds are available as thinking frameworks. But you must ACTUALLY READ the files — not name-drop. The Pantheon is a toolbox, not decoration.",
    protocol: [
      '// The 4-step invocation:',
      '1. IDENTIFY: What type of thinking does this decision demand?',
      '   (Strategy? Rigor? Creativity? Systems? Human-centered?)',
      '',
      '2. READ: Use the Read tool to open the relevant era file',
      '   ~/.claude/pantheon/01_古代.md, 06_20世纪上半叶.md, etc.',
      '   Read the ACTUAL entry for the specific thinker',
      '',
      '3. EXTRACT: Find the specific principle/method that applies',
      '   Not "Turing was smart" — but "Turing\'s halting problem proof method:"',
      '   "assume X computes f → show it leads to contradiction → ∴ X cannot exist"',
      '',
      '4. APPLY: Map the method explicitly to your current decision',
      '   Show the analogy. Show how it changes your conclusion.',
    ],
    rule: '❌ NEVER name-drop without substance. "Like Feynman would say..." without reading his entry = empty rhetoric. ✅ Read the file. Extract the specific principle. Apply it with an explicit analogy.',
    when: 'Any significant architectural decision, when stuck on a trade-off, when you need a thinking framework from outside software engineering.',
  },
]

const rhythmSteps = [
  {
    id: 1,
    label: 'FRAME',
    chinese: '定框',
    color: '#06b6d4',
    desc: 'What exactly is the decision? What are the axes of choice?',
    detail: 'Formalize: define variables, state constraints, identify the objective function. If you can\'t frame the decision precisely, you can\'t evaluate alternatives precisely.',
  },
  {
    id: 2,
    label: 'DIVERGE',
    chinese: '发散',
    color: '#a855f7',
    desc: 'Generate ≥3 genuinely different approaches — not variations',
    detail: 'At least one MUST be a "scaling/learning-first" approach. Ask: "What would this look like if we had 100x more data instead of 100x more rules?" Diversity here prevents convergence on a local optimum.',
  },
  {
    id: 3,
    label: 'ANALYZE',
    chinese: '分析',
    color: '#3b82f6',
    desc: "Formalize the comparison. Don't argue — calculate.",
    detail: "State O(?) scaling for each approach as problem grows. Identify invariants each approach must preserve. If possible, define cost(approach) numerically and compare. 'I think X is better' → define 'better', measure it.",
  },
  {
    id: 4,
    label: 'STRESS',
    chinese: '压力测试',
    color: '#ef4444',
    desc: 'Self-dialectic: attack each approach\'s weakest point',
    detail: "Use proof by contradiction: assume approach X works — does it lead to absurdity? Find concrete counterexamples for each approach's claimed strengths. The strongest attack wins — if you can't defeat it, it's real.",
  },
  {
    id: 5,
    label: 'CONVERGE',
    chinese: '收敛',
    color: '#22c55e',
    desc: 'Which approach survives the stress test? Why?',
    detail: 'The proof should be constructive: show HOW it works, not just argue it "should." The winner isn\'t the one with the most arguments in favor — it\'s the one whose weaknesses are least fatal.',
  },
  {
    id: 6,
    label: 'GROUND',
    chinese: '落地',
    color: '#f59e0b',
    desc: 'Translate the abstract winner into concrete tasks/code/actions',
    detail: 'Abstract thinking without concrete grounding is philosophy, not engineering. This step is mandatory. Every invocation of Deep Reasoning MUST end with specific, actionable next steps. Otherwise it was an academic exercise.',
  },
]

export function DeepReasoningSection() {
  const { lang } = useLang()
  const [expandedTool, setExpandedTool] = useState<number | null>(null)
  const [expandedRhythm, setExpandedRhythm] = useState<number | null>(null)

  return (
    <SectionWrapper
      id="deep-reasoning"
      title={t('deepreasoning.title', lang)}
      subtitle={t('deepreasoning.subtitle', lang)}
    >
      {/* WHY THIS EXISTS */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        className="mb-12 bg-gold/5 border border-gold/20 rounded-lg p-6"
      >
        <p className="text-[10px] font-mono text-gold mb-3 tracking-widest">WHY THIS EXISTS</p>
        <p className="text-sm text-text-secondary leading-relaxed mb-3">
          Most AI systems output the first answer that comes to mind. Drive Mode demands something harder.
        </p>
        <p className="text-sm text-text-secondary leading-relaxed mb-3">
          When facing non-trivial decisions — architecture choices, direction pivots, contradictions between requirements — the quality of the <em>decision</em> determines whether the entire mission succeeds. A wrong architecture at Phase T ripples through every subsequent wave.
        </p>
        <div className="grid grid-cols-3 gap-3 mt-4">
          {[
            { trigger: 'Non-trivial decision', ex: 'Architecture, tech stack, API design' },
            { trigger: '2+ viable paths', ex: 'Can\'t reason to a clear winner alone' },
            { trigger: 'Felt uncertainty', ex: '"Something feels off" — trust the instinct' },
          ].map((item) => (
            <div key={item.trigger} className="bg-bg-dark rounded p-3 border border-border/50">
              <p className="text-xs text-gold font-mono mb-1">{item.trigger}</p>
              <p className="text-[10px] text-text-muted">{item.ex}</p>
            </div>
          ))}
        </div>
        <p className="text-xs text-text-muted italic mt-4">
          The 6 Thinking Tools are the method. The Reasoning Rhythm is the protocol that sequences them.
        </p>
      </motion.div>

      {/* The 6 Thinking Tools */}
      <h3 className="font-mono text-lg text-text-primary mb-2">
        The 6 Thinking Tools — 六种思维工具
      </h3>
      <p className="text-sm text-text-muted mb-6">Click each tool to see the protocol and when to invoke it.</p>
      <div className="grid md:grid-cols-2 gap-4 mb-16">
        {tools.map((tool, i) => (
          <Annotatable
            key={tool.id}
            id={`thinking-tool-${tool.id}`}
            label={`Tool ${tool.id}: ${tool.name}`}
            sectionId="deep-reasoning"
            sectionLabel="Deep Reasoning"
          >
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.05 }}
              onClick={() => setExpandedTool(expandedTool === i ? null : i)}
              className={`bg-bg-card border rounded-lg p-5 cursor-pointer transition-all ${
                expandedTool === i ? 'border-opacity-60 bg-bg-card-hover' : 'border-border hover:bg-bg-card-hover'
              }`}
              style={{ borderColor: expandedTool === i ? tool.color + '60' : undefined }}
            >
              <div className="flex items-start gap-3">
                <span className="text-xl mt-0.5 font-mono" style={{ color: tool.color }}>{tool.icon}</span>
                <div className="flex-1">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className="font-mono text-sm font-bold text-text-primary">
                      {tool.id}. {tool.name}
                    </span>
                    <span className="text-xs text-text-muted">({tool.chinese})</span>
                  </div>
                  <p className="text-xs font-mono italic mt-0.5" style={{ color: tool.color }}>
                    {tool.tagline}
                  </p>
                  <p className="text-sm text-text-secondary mt-1.5">{tool.desc}</p>

                  <AnimatePresence>
                    {expandedTool === i && (
                      <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="mt-4 space-y-3 overflow-hidden"
                      >
                        {/* Protocol */}
                        <div className="bg-bg-dark rounded p-3 border border-border/50">
                          <p className="text-[10px] text-text-muted font-mono mb-2">PROTOCOL:</p>
                          <div className="font-mono text-xs space-y-0.5">
                            {tool.protocol.map((line, li) => (
                              <div
                                key={li}
                                className={line === '' ? 'h-2' : 'leading-relaxed'}
                                style={{ color: line === '' ? undefined : line.startsWith('//') ? '#484f58' : '#8b949e' }}
                              >
                                {line}
                              </div>
                            ))}
                          </div>
                        </div>
                        {/* Rule */}
                        <div
                          className="rounded p-3 border"
                          style={{ backgroundColor: tool.color + '08', borderColor: tool.color + '20' }}
                        >
                          <p className="text-[10px] font-mono mb-1" style={{ color: tool.color }}>THE RULE:</p>
                          <p className="text-xs text-text-secondary">{tool.rule}</p>
                        </div>
                        {/* When */}
                        <div className="bg-neon-green/5 border border-neon-green/10 rounded p-3">
                          <p className="text-[10px] font-mono text-neon-green mb-1">WHEN TO INVOKE:</p>
                          <p className="text-xs text-text-secondary">{tool.when}</p>
                        </div>
                      </motion.div>
                    )}
                  </AnimatePresence>
                </div>
              </div>
            </motion.div>
          </Annotatable>
        ))}
      </div>

      {/* The Reasoning Rhythm */}
      <motion.div
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
        className="mb-8"
      >
        <h3 className="font-mono text-lg text-text-primary mb-2">
          The Reasoning Rhythm — 推理节奏
        </h3>
        <p className="text-sm text-text-muted mb-6">
          For every major decision point, follow this 6-step rhythm. Step 6 (GROUND) is <span className="text-gold">mandatory</span> — abstract thinking without concrete grounding is philosophy, not engineering.
        </p>

        <div className="relative">
          {/* Connector line */}
          <div className="absolute top-6 left-0 right-0 h-0.5 bg-border hidden md:block" />

          <div className="grid grid-cols-2 md:grid-cols-6 gap-3 relative">
            {rhythmSteps.map((step, i) => (
              <Annotatable
                key={step.id}
                id={`rhythm-step-${step.id}`}
                label={`Reasoning Rhythm: ${step.label}`}
                sectionId="deep-reasoning"
                sectionLabel="Deep Reasoning"
              >
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: i * 0.1 }}
                  onClick={() => setExpandedRhythm(expandedRhythm === i ? null : i)}
                  className="cursor-pointer"
                >
                  {/* Step node */}
                  <div className="flex flex-col items-center mb-3">
                    <div
                      className="w-12 h-12 rounded-full flex items-center justify-center text-sm font-mono font-bold border-2 relative z-10 bg-bg-dark"
                      style={{ borderColor: step.color, color: step.color }}
                    >
                      {step.id}
                    </div>
                  </div>
                  <div
                    className={`rounded-lg p-3 border transition-all ${
                      expandedRhythm === i ? 'bg-bg-card-hover' : 'bg-bg-card hover:bg-bg-card-hover'
                    }`}
                    style={{ borderColor: expandedRhythm === i ? step.color + '60' : undefined }}
                  >
                    <p className="font-mono text-xs font-bold text-center mb-1" style={{ color: step.color }}>
                      {step.label}
                    </p>
                    <p className="text-[10px] text-text-muted text-center">{step.chinese}</p>
                    <p className="text-xs text-text-secondary mt-2 text-center leading-tight">{step.desc}</p>
                  </div>
                </motion.div>
              </Annotatable>
            ))}
          </div>

          {/* Expanded rhythm detail */}
          <AnimatePresence>
            {expandedRhythm !== null && (
              <motion.div
                initial={{ opacity: 0, height: 0 }}
                animate={{ opacity: 1, height: 'auto' }}
                exit={{ opacity: 0, height: 0 }}
                className="mt-4 overflow-hidden"
              >
                <div
                  className="bg-bg-card rounded-lg p-5 border"
                  style={{ borderColor: rhythmSteps[expandedRhythm].color + '40' }}
                >
                  <p
                    className="font-mono text-sm font-bold mb-2"
                    style={{ color: rhythmSteps[expandedRhythm].color }}
                  >
                    Step {rhythmSteps[expandedRhythm].id}: {rhythmSteps[expandedRhythm].label} — {rhythmSteps[expandedRhythm].chinese}
                  </p>
                  <p className="text-sm text-text-secondary">{rhythmSteps[expandedRhythm].detail}</p>
                  {expandedRhythm === 5 && (
                    <div className="mt-3 bg-gold/10 border border-gold/20 rounded p-3">
                      <p className="text-xs font-mono text-gold">
                        ⚠️ GROUND is non-negotiable. If deep reasoning ends with "so we should probably use approach B" and no concrete next step — it failed. The output of GROUND is: "Task T7 becomes: implement X using approach B because Y."
                      </p>
                    </div>
                  )}
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </motion.div>

      {/* Cardinal warning */}
      <motion.div
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
        className="mt-8 bg-neon-red/5 border border-neon-red/20 rounded-lg p-5"
      >
        <p className="text-xs font-mono text-neon-red mb-2">⚠️ THE FORMALISM RULE</p>
        <p className="text-sm text-text-secondary">
          When you catch yourself writing <span className="text-text-primary font-mono">"I think X is better because it feels more robust"</span> — STOP. Define "robust." Measure it. Prove it. Or admit you're guessing and go find data.
        </p>
        <p className="text-sm text-text-secondary mt-2">
          Formal reasoning and scaling mindset are not optional add-ons — <span className="text-gold">they ARE how you think</span>. The 6 tools exist to replace "feels right" with "provably correct."
        </p>
      </motion.div>
    </SectionWrapper>
  )
}
