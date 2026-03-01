import { motion } from 'framer-motion'
import { SectionWrapper } from '../components/SectionWrapper'
import { useState } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

// WHY THE PANTHEON EXISTS
// ────────────────────────
// Claude already knows Einstein, Turing, Sun Tzu, von Neumann.
// The knowledge isn't the point — activation mapping is.
// Without explicit guidance, Claude doesn't naturally think:
//   "I should apply Gödel's incompleteness to reason about this system's limits"
//   "I should apply Sun Tzu's center-of-gravity to this architecture decision"
// The Pantheon is a prompt that says: when facing problem type X, invoke thinker Y's method.
// It's not a reference library. It's an activation key.

const activationMap = [
  {
    type: 'Architecture & Structure',
    typeCh: '架构与结构',
    thinkers: ['Euclid', '亚里士多德', '冯·诺依曼'],
    question: 'What are the axioms? What\'s the minimal correct architecture?',
    why: 'These thinkers mastered axiomatic construction — building from irreducible first principles. Euclid proved geometry from 5 axioms. Von Neumann distilled computation to a minimal architecture. When Claude invokes them, it asks: what are the irreducible constraints here?',
    color: '#f59e0b',
    icon: '📐',
  },
  {
    type: 'Strategy & Trade-offs',
    typeCh: '战略与权衡',
    thinkers: ['孙子', '克劳塞维茨', '马基雅维利'],
    question: 'What\'s the center of gravity? What terrain am I fighting on?',
    why: 'Strategy is the art of selecting the right battle, not just winning the current one. 孙子\'s center-of-gravity thinking reframes decisions: not "what feature should I build?" but "what decision, if right, makes everything else easier?"',
    color: '#ef4444',
    icon: '⚔️',
  },
  {
    type: 'Creative Breakthroughs',
    typeCh: '创意突破',
    thinkers: ['达芬奇', '特斯拉', '庄子'],
    question: 'What if I invert this? What does another domain teach me?',
    why: '达芬奇\'s cross-domain transfer method — borrowing from anatomy to inform hydraulics, art to inform engineering. When stuck in one framing, invoke this to force inversion: what if the constraint is actually the solution?',
    color: '#a855f7',
    icon: '💡',
  },
  {
    type: 'Rigor & Correctness',
    typeCh: '严谨与正确性',
    thinkers: ['哥德尔', '图灵', '费曼'],
    question: 'What can\'t this system express? Can I explain it simply?',
    why: 'Gödel proved every sufficiently powerful formal system has statements it cannot prove from within. Applied: every architecture has blind spots it can\'t see by its own rules. Feynman\'s technique: if you can\'t explain it simply, you don\'t understand it — use this as a correctness check.',
    color: '#06b6d4',
    icon: '🔬',
  },
  {
    type: 'Human-Centered Design',
    typeCh: '以人为本的设计',
    thinkers: ['孔子', '乔布斯'],
    question: 'What does the user truly need vs. what they said?',
    why: '孔子\'s emphasis on 仁 (benevolence as practice, not concept): understanding what people actually need requires patient inquiry, not assumption. Jobs\' "people don\'t know what they want until you show them" — the user\'s stated request is the starting point, not the specification.',
    color: '#22c55e',
    icon: '🤝',
  },
  {
    type: 'Systematic Method',
    typeCh: '系统化方法',
    thinkers: ['笛卡尔', '培根'],
    question: 'How do I decompose this into solvable parts?',
    why: 'Descartes\' method: divide every problem into the smallest subproblems you can solve. Bacon\'s induction: collect observations before concluding. Together they form the research-then-plan discipline — never start from assumptions, always start from evidence.',
    color: '#3b82f6',
    icon: '🗂️',
  },
  {
    type: 'Evolution & Adaptation',
    typeCh: '进化与适应',
    thinkers: ['达尔文'],
    question: 'What survives? What\'s the selection pressure here?',
    why: 'Darwin\'s core insight applied to engineering: don\'t ask "what\'s the perfect design?" but "what survives the constraints?" Selection pressure thinking: given the real-world pressures (performance, maintainability, time), what design pattern survives longest?',
    color: '#f97316',
    icon: '🧬',
  },
]

const invocationProtocol = [
  {
    step: '1',
    label: 'Identify Problem Type',
    desc: 'What kind of thinking does this decision require? Strategy? Rigor? Creative inversion?',
    code: 'Architecture decision → invoke Euclid / von Neumann',
  },
  {
    step: '2',
    label: 'Read the Era File',
    desc: 'Use the Read tool to open the relevant era file from ~/.claude/pantheon/. Read the actual entry.',
    code: 'Read("~/.claude/pantheon/06_20世纪上半叶.md")',
  },
  {
    step: '3',
    label: 'Extract a Specific Principle',
    desc: 'Find the specific idea from their entry that maps to your problem. Not their name — their method.',
    code: '→ Turing: "Can I define a formal test for this?" (Turing completeness as architecture check)',
  },
  {
    step: '4',
    label: 'Apply Explicitly',
    desc: 'Show the analogy. State how the principle maps. Make the reasoning visible.',
    code: '"Applying Feynman: explain this system simply. If I can\'t, the design is too complex."',
  },
]

const eraFiles = [
  { file: '01_古代.md', label: '古代', range: '孔子, 孙子, 亚里士多德, Euclid, 庄子...', color: '#f59e0b', count: 200 },
  { file: '02_中世纪.md', label: '中世纪', range: '阿维森纳, 托马斯·阿奎那, 伊本·赫勒敦...', color: '#ef4444', count: 120 },
  { file: '03_文艺复兴与早期近代.md', label: '文艺复兴', range: '达芬奇, 伽利略, 牛顿, 马基雅维利...', color: '#a855f7', count: 115 },
  { file: '04_启蒙时代与18世纪.md', label: '启蒙时代', range: '笛卡尔, 培根, 休谟, 亚当·斯密...', color: '#3b82f6', count: 70 },
  { file: '05_19世纪.md', label: '19世纪', range: '达尔文, 马克思, 尼采, 麦克斯韦...', color: '#06b6d4', count: 165 },
  { file: '06_20世纪上半叶.md', label: '20世纪上', range: '爱因斯坦, 图灵, 冯·诺依曼, 哥德尔...', color: '#22c55e', count: 190 },
  { file: '07_20世纪下半叶至当代.md', label: '当代', range: '费曼, 乔布斯, 香农, 迪杰斯特拉...', color: '#f97316', count: 270 },
]

export function PantheonSection() {
  const { lang } = useLang()
  const [expandedThinker, setExpandedThinker] = useState<number | null>(null)

  return (
    <SectionWrapper
      id="pantheon"
      title={t('pantheon.title', lang)}
      subtitle={t('pantheon.subtitle', lang)}
    >
      {/* WHY THIS EXISTS — the most important block */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        className="mb-12 bg-gold/5 border border-gold/20 rounded-xl p-8"
      >
        <p className="text-[10px] font-mono text-gold mb-4 tracking-widest">WHY THIS EXISTS</p>
        <h3 className="text-xl font-bold text-text-primary mb-4">
          Claude Already Knows These People. That's Not the Point.
        </h3>
        <p className="text-sm text-text-secondary mb-4 leading-relaxed">
          Claude has read everything Einstein, Turing, and 孙子 ever wrote. Their ideas are already in the model weights. So why does a directory of 1000+ historical figures exist in the skill?
        </p>
        <div className="bg-bg-dark rounded-lg p-5 mb-4 border border-border">
          <p className="text-xs font-mono text-neon-red mb-2">THE ACTIVATION PROBLEM</p>
          <p className="text-sm text-text-secondary leading-relaxed">
            Without explicit guidance, Claude doesn't naturally think:{' '}
            <span className="text-gold italic">"I should apply Gödel's incompleteness to reason about this system's limits"</span>{' '}
            or{' '}
            <span className="text-gold italic">"I should use Sun Tzu's center-of-gravity to find the one architectural decision that makes everything else easier."</span>
          </p>
          <p className="text-sm text-text-secondary leading-relaxed mt-3">
            Claude defaults to technical reasoning. It reaches for the obvious tool. The Pantheon is an activation key — it maps <span className="text-text-primary font-semibold">problem type → thinking framework → specific method</span>. Not "who was smart?" but "when facing <em>this</em> type of decision, invoke <em>that</em> person's method."
          </p>
        </div>
        <div className="grid md:grid-cols-3 gap-3">
          {[
            { label: 'NOT this', icon: '❌', text: '"Include von Neumann for architecture" — vague attribution, no substance', color: 'neon-red' },
            { label: 'NOT this either', icon: '❌', text: 'Loading all 400KB upfront "just in case" — wastes context window on irrelevant noise', color: 'neon-red' },
            { label: 'THIS', icon: '✅', text: 'Read() the relevant era file → extract the specific principle → apply explicitly to the current decision', color: 'neon-green' },
          ].map((item) => (
            <div key={item.label} className={`rounded-lg p-3 bg-${item.color}/5 border border-${item.color}/20`}>
              <p className={`text-[10px] font-mono text-${item.color} mb-1`}>{item.icon} {item.label}</p>
              <p className="text-xs text-text-muted leading-relaxed">{item.text}</p>
            </div>
          ))}
        </div>
      </motion.div>

      {/* Invocation Protocol */}
      <div className="mb-12">
        <h3 className="font-mono text-lg text-text-primary mb-2">The Invocation Protocol</h3>
        <p className="text-sm text-text-muted mb-6">Four steps from "I have a hard decision" to "I applied Feynman's method"</p>
        <div className="grid md:grid-cols-4 gap-3">
          {invocationProtocol.map((step, i) => (
            <motion.div
              key={step.step}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.1 }}
              className="bg-bg-card border border-border rounded-lg p-4 relative"
            >
              <div className="absolute -top-3 left-4 bg-gold text-bg-dark text-xs font-mono font-bold px-2 py-0.5 rounded">
                {step.step}
              </div>
              <p className="text-sm font-bold text-text-primary mt-1 mb-2">{step.label}</p>
              <p className="text-xs text-text-muted mb-3 leading-relaxed">{step.desc}</p>
              <div className="bg-bg-dark rounded p-2 border border-border/50">
                <p className="text-[10px] font-mono text-gold/80 leading-relaxed">{step.code}</p>
              </div>
            </motion.div>
          ))}
        </div>
      </div>

      {/* Era Files — redesigned as on-demand loading explanation */}
      <div className="mb-12">
        <h3 className="font-mono text-lg text-text-primary mb-2">Seven Era Files — On-Demand Loading</h3>
        <p className="text-sm text-text-muted mb-2 leading-relaxed">
          400KB of structured wisdom. Never loaded upfront — injected via{' '}
          <span className="font-mono text-gold">Read()</span> at the moment of a hard decision.
          Why 7 files? Because architectural decisions span different domains of thought. Loading the era
          closest to your problem type focuses the context.
        </p>
        <div className="bg-neon-blue/5 border border-neon-blue/20 rounded p-3 mb-6 text-xs font-mono text-text-muted">
          <span className="text-neon-blue">Design principle:</span> Context-window cost is paid at invocation time, not session start.
          Loading "06_20世纪上半叶.md" when you face a computation theory question gives you Turing,
          von Neumann, and Shannon exactly when needed — not as background noise for the whole session.
        </div>
        <div className="space-y-2">
          {eraFiles.map((era, i) => (
            <motion.div
              key={era.file}
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.05 }}
              className="flex items-center gap-4 p-3 bg-bg-card border border-border rounded-lg hover:border-gold/20 transition-colors group"
            >
              {/* Bar */}
              <div className="w-24 flex-shrink-0">
                <div
                  className="h-2 rounded-full"
                  style={{ width: `${(era.count / 270) * 100}%`, backgroundColor: era.color + '60', border: `1px solid ${era.color}` }}
                />
              </div>
              <div className="flex-shrink-0 w-20">
                <span className="text-xs font-mono font-bold" style={{ color: era.color }}>{era.label}</span>
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-xs font-mono text-text-muted truncate">{era.range}</p>
              </div>
              <div className="flex-shrink-0">
                <span className="text-[10px] font-mono text-text-muted">{era.count} figures</span>
              </div>
              <div className="flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
                <span className="text-[10px] font-mono text-gold">{era.file}</span>
              </div>
            </motion.div>
          ))}
        </div>
        <p className="text-xs text-text-muted mt-3 font-mono text-right">Total: 1,130+ figures | ~400KB across 7 files</p>
      </div>

      {/* Activation Mapping — the core value */}
      <div className="mb-10">
        <h3 className="font-mono text-lg text-gold mb-2">Activation Mapping — When to Summon Which Thinker</h3>
        <p className="text-sm text-text-muted mb-6">
          For each type of decision, there's a class of thinkers whose method directly applies.
          Click to see WHY each pairing works — and the specific question it triggers.
        </p>
        <div className="space-y-3">
          {activationMap.map((item, i) => (
            <motion.div
              key={item.type}
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.05 }}
              onClick={() => setExpandedThinker(expandedThinker === i ? null : i)}
              className="bg-bg-card border border-border rounded-lg p-4 cursor-pointer hover:border-gold/20 transition-colors"
            >
              <div className="flex items-start gap-4">
                <span className="text-xl flex-shrink-0 mt-0.5">{item.icon}</span>
                <div className="flex-1 min-w-0">
                  <div className="flex flex-col md:flex-row md:items-center gap-2 mb-1">
                    <span className="font-mono text-sm font-bold text-text-primary">{item.type}</span>
                    <span className="text-text-muted text-xs">({item.typeCh})</span>
                  </div>
                  <div className="flex flex-wrap gap-2 mb-2">
                    {item.thinkers.map((thinker) => (
                      <span
                        key={thinker}
                        className="text-xs px-2 py-0.5 rounded font-mono"
                        style={{ backgroundColor: item.color + '15', color: item.color, border: `1px solid ${item.color}40` }}
                      >
                        {thinker}
                      </span>
                    ))}
                  </div>
                  <p className="text-xs text-text-muted italic">
                    Trigger question: "{item.question}"
                  </p>

                  {expandedThinker === i && (
                    <motion.div
                      initial={{ opacity: 0, height: 0 }}
                      animate={{ opacity: 1, height: 'auto' }}
                      className="mt-3 pt-3 border-t border-border/50"
                    >
                      <p className="text-[10px] font-mono mb-1" style={{ color: item.color }}>WHY THIS PAIRING WORKS:</p>
                      <p className="text-sm text-text-secondary leading-relaxed">{item.why}</p>
                    </motion.div>
                  )}
                </div>
                <span className="text-text-muted text-xs flex-shrink-0 mt-1">
                  {expandedThinker === i ? '▲' : '▼'}
                </span>
              </div>
            </motion.div>
          ))}
        </div>
      </div>

      {/* The Cardinal Rule */}
      <motion.div
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
        className="p-5 bg-neon-red/5 border border-neon-red/20 rounded-lg"
      >
        <p className="text-xs font-mono text-neon-red mb-2">THE CARDINAL RULE OF PANTHEON USE</p>
        <p className="text-sm text-text-secondary font-mono leading-relaxed">
          <span className="text-neon-red">Never name-drop without substance.</span>{' '}
          If you invoke 孙子, quote a specific principle and show how it maps.
          If you invoke 费曼, actually apply the Feynman Technique (explain simply, find where your explanation breaks down).
          If you invoke 图灵, actually apply his method (formalize the question, define what a solution looks like).{' '}
          <span className="text-text-primary">The Pantheon is a toolbox. Each tool has a specific use. Waving it around achieves nothing.</span>
        </p>
      </motion.div>
    </SectionWrapper>
  )
}
