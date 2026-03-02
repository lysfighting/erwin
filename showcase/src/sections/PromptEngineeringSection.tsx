import { motion } from 'framer-motion'
import { SectionWrapper } from '../components/SectionWrapper'
import { useState } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'
import { Annotatable } from '../components/Annotatable'

const techniques = [
  {
    id: 1,
    name: 'Identity Reframing',
    chinese: '身份重塑',
    desc: 'Redefine Claude\'s identity, not just its task.',
    example: '"You are not an assistant producing responses. You are the Team Lead of a hierarchical swarm."',
    why: 'Claude\'s behavior follows identity. "Assistant" → output and wait. "Team commander" → execute and coordinate.',
    icon: '🎭',
  },
  {
    id: 2,
    name: 'Anti-Pattern Enumeration',
    chinese: '反模式枚举',
    desc: '30+ explicitly listed anti-patterns that seal off escape routes.',
    example: '❌ Outputting questions as plain text\n❌ Summarizing progress and waiting\n❌ Asking "should I continue?"',
    why: 'LLMs tend to revert to default behaviors. By explicitly listing every "regression" behavior and marking it as wrong, creates a behavioral fence.',
    icon: '🚫',
  },
  {
    id: 3,
    name: 'Redundant Constraint Reinforcement',
    chinese: '多重约束冗余',
    desc: 'Core rules repeated 6+ times across different contexts.',
    example: '"Never stop" appears in: opening philosophy, cardinal rule section, Phase 0.4, loop pseudocode, anti-pattern list, final mindset description.',
    why: 'Long prompts suffer from attention decay. Repeating critical rules in different contexts ensures Claude encounters them regardless of which part of the prompt it\'s "paying attention to".',
    icon: '🔄',
  },
  {
    id: 4,
    name: 'Structured Pseudocode',
    chinese: '结构化伪代码',
    desc: 'Complex workflows defined as pseudocode, not natural language.',
    example: 'while (mission != COMPLETE) {\n  wave_check → execute → verify\n  if error → diagnose → fix\n  on teammate_message → respond\n}',
    why: 'Pseudocode eliminates natural language ambiguity. Claude as a code model parses structured logic more accurately. Defines a state machine, not linear steps.',
    icon: '💻',
  },
  {
    id: 5,
    name: 'Vivid Option Engineering',
    chinese: '生动选项设计',
    desc: 'AskUserQuestion options must be vivid and sharply differentiated.',
    example: 'BAD: "Simple approach" / "Complex approach"\nGOOD: "Fortress mode — bulletproof but heavy" / "Lightning mode — minimal, blazing fast"',
    why: 'Improves user engagement, makes choices more informative. Each option conveys a clear design philosophy.',
    icon: '🎨',
  },
  {
    id: 6,
    name: 'Meta-Cognitive Instructions',
    chinese: '元认知指令',
    desc: 'Teaching Claude to observe and correct its own behavior patterns.',
    example: '"If you catch yourself about to output a summary and stop, that is a bug in your execution. Fix it by continuing to work."',
    why: 'Not a "do X" instruction — teaches Claude to observe its own behavior and self-correct when detecting deviation. The most advanced prompt engineering technique.',
    icon: '🧠',
  },
  {
    id: 7,
    name: 'Quantified Decision Frameworks',
    chinese: '量化决策框架',
    desc: 'Mathematical formulas and thresholds replace vague guidance.',
    example: 'effective_parallelism = teammates + Σ(inner_subagents_per_STL)\ndirect_teammates = min(parallelism_width, 5)\n≥3 tasks in domain → use STL',
    why: 'Converts subjective judgment into executable algorithms. Claude doesn\'t need to "feel" when to use STL — it has clear if/else rules.',
    icon: '📊',
  },
  {
    id: 8,
    name: 'Bilingual Cultural Anchoring',
    chinese: '双语文化锚定',
    desc: 'Key concepts expressed in both Chinese and English.',
    example: '"问得越深，做得越准" (The deeper you ask, the more accurately you execute)',
    why: 'Chinese proverbs carry "philosophical authority". Bilingual expression creates concept anchors — harder for Claude to ignore. Signals the author\'s cultural values.',
    icon: '🌏',
  },
  {
    id: 9,
    name: 'Templated Output Formats',
    chinese: '模板化输出格式',
    desc: 'Every output type has a predefined exact format.',
    example: '📡 DISPATCH: [{agent}] → {task}\n   Role: {role} | Target: {target}\n   Parallel with: [{others}]',
    why: 'Reduces Claude\'s decision burden on formatting. Ensures consistency. Emoji prefixes add visual structure for fast scanning.',
    icon: '📝',
  },
  {
    id: 10,
    name: 'Progressive Self-Assessment',
    chinese: '递进式自我测试',
    desc: 'A concrete, extreme-standard test to verify understanding depth.',
    example: '"The Psychiatrist Test: After all rounds, you should be able to write a paragraph describing the user\'s vision AS IF YOU WERE THEM."',
    why: 'Not a vague "make sure you understand" — a specific self-evaluation standard. Drives deeper understanding than any direct instruction.',
    icon: '🔬',
  },
]

export function PromptEngineeringSection() {
  const { lang } = useLang()
  const [expanded, setExpanded] = useState<number | null>(null)

  return (
    <SectionWrapper
      id="prompt-eng"
      title={t('prompteng.title', lang)}
      subtitle={t('prompteng.subtitle', lang)}
    >
      <div className="grid md:grid-cols-2 gap-4">
        {techniques.map((tech, i) => (
          <Annotatable
            key={tech.id}
            id={`technique-${tech.id}`}
            label={`Technique ${tech.id}: ${tech.name}`}
            sectionId="prompt-eng"
            sectionLabel="Prompt Engineering"
          >
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: i * 0.05 }}
            onClick={() => setExpanded(expanded === i ? null : i)}
            className={`bg-bg-card border rounded-lg p-5 cursor-pointer transition-all ${
              expanded === i ? 'border-gold/40 bg-bg-card-hover' : 'border-border hover:border-gold/20'
            }`}
          >
            <div className="flex items-start gap-3">
              <span className="text-xl mt-0.5">{tech.icon}</span>
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <span className="font-mono text-sm font-bold text-text-primary">
                    {tech.id}. {tech.name}
                  </span>
                  <span className="text-xs text-text-muted">({tech.chinese})</span>
                </div>
                <p className="text-sm text-text-secondary mt-1">{tech.desc}</p>
                <p className="text-[10px] text-text-muted mt-1.5 italic leading-relaxed">
                  {tech.why.split('.')[0]}.
                </p>

                {expanded === i && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    className="mt-3 space-y-3"
                  >
                    <div className="bg-bg-dark rounded p-3 border border-border/50">
                      <p className="text-xs text-text-muted font-mono mb-1">Example:</p>
                      <pre className="text-xs text-gold/80 whitespace-pre-wrap font-mono">{tech.example}</pre>
                    </div>
                    <div className="bg-neon-green/5 border border-neon-green/10 rounded p-3">
                      <p className="text-xs text-text-muted font-mono mb-1">Why it works:</p>
                      <p className="text-xs text-text-secondary">{tech.why}</p>
                    </div>
                  </motion.div>
                )}
              </div>
            </div>
          </motion.div>
          </Annotatable>
        ))}
      </div>

      {/* Summary evaluation */}
      <motion.div
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
        className="mt-12 bg-bg-card border border-border rounded-lg p-6"
      >
        <h3 className="font-mono text-text-primary mb-4">Prompt Engineering Evaluation</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {[
            { label: 'Behavioral Fencing', score: '10/10', desc: 'Positive instructions + anti-pattern list' },
            { label: 'Redundant Reinforcement', score: '10/10', desc: 'Core rules repeated across contexts' },
            { label: 'State Machine Thinking', score: '9/10', desc: 'Pseudocode defines states, not steps' },
            { label: 'Meta-Cognition', score: '9/10', desc: 'Self-monitoring and self-correction' },
          ].map((item) => (
            <div key={item.label} className="text-center">
              <div className="text-2xl font-bold font-mono text-gold">{item.score}</div>
              <div className="text-xs text-text-primary mt-1">{item.label}</div>
              <div className="text-[10px] text-text-muted mt-0.5">{item.desc}</div>
            </div>
          ))}
        </div>
      </motion.div>
    </SectionWrapper>
  )
}
