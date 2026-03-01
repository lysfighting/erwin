import { motion } from 'framer-motion'
import { SectionWrapper } from '../components/SectionWrapper'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

const references = [
  {
    category: 'Architecture & Methodology',
    items: [
      { name: 'ATAM', author: 'Carnegie Mellon SEI', desc: 'Architecture Tradeoff Analysis Method → Phase 3-5 structure', attributed: false },
      { name: 'C4 Model', author: 'Simon Brown', desc: 'Context/Container/Component/Code → Phase 6 diagram types', attributed: false },
      { name: 'ADR', author: 'Michael Nygard (2011)', desc: 'Architecture Decision Records → Key Decisions Log format', attributed: false },
      { name: 'Devil\'s Advocate', author: 'Classical rhetoric', desc: 'Adversarial review technique → Phase 3c, Phase 5', attributed: false },
    ],
  },
  {
    category: 'AI Research Philosophy',
    items: [
      { name: 'The Bitter Lesson', author: 'Richard Sutton (2019)', desc: 'General methods + compute > handcrafted rules → Scaling Mindset chapter', attributed: false },
      { name: 'Ralph Wiggum Technique', author: 'Geoffrey Huntley (2025)', desc: 'Autonomous loop execution → drive description triggers', attributed: false },
    ],
  },
  {
    category: 'Philosophical Foundations',
    items: [
      { name: 'Hegelian Dialectic', author: 'G.W.F. Hegel (19th century)', desc: 'Thesis→Antithesis→Synthesis → Self-Dialectic protocol', attributed: false },
      { name: 'First Principles', author: 'Aristotle / Elon Musk', desc: 'Strip to atoms, rebuild → First Principles Decomposition', attributed: false },
      { name: 'Occam\'s Razor', author: 'William of Ockham (14th century)', desc: 'Simplicity principle → architect.md core philosophy', attributed: false },
      { name: 'Feynman Technique', author: 'Richard Feynman', desc: 'Explain simply to verify understanding → Pantheon invocation', attributed: true },
      { name: '孙子兵法 (Art of War)', author: '孙武 (5th century BC)', desc: 'Military strategy metaphors → Team Lead, Battle Plan, Waves, Reconnaissance', attributed: true },
    ],
  },
  {
    category: 'Computer Science',
    items: [
      { name: 'Fork-Join Model', author: 'Melvin Conway (1963)', desc: 'Parallel computation pattern → Wave execution model', attributed: false },
      { name: 'Pipeline Parallelism', author: 'Computer architecture', desc: 'Stage overlap → Pipeline concurrency pattern', attributed: false },
      { name: 'Speculative Execution', author: 'CPU design', desc: 'Branch prediction analog → Speculative concurrency pattern', attributed: false },
      { name: 'Impact-Effort Matrix', author: 'Product management (1980s)', desc: '2x2 priority framework → self-drive proposal ranking', attributed: false },
    ],
  },
  {
    category: 'Platform & APIs',
    items: [
      { name: 'Claude Code Agent Teams', author: 'Anthropic', desc: 'TeamCreate, TeamDelete, SendMessage → Core platform dependency', attributed: true },
      { name: 'Slack Socket Mode', author: 'Slack', desc: 'WebSocket real-time communication → slack-listen.py', attributed: true },
      { name: 'Telegram Bot API', author: 'Telegram', desc: 'sendMessage, getUpdates → telegram-*.py', attributed: true },
      { name: 'Fish Audio TTS', author: 'Fish Audio', desc: 'Text-to-speech API → speak.sh (deprecated)', attributed: true },
    ],
  },
]

const competitors = [
  {
    name: 'Drive Skill',
    author: 'KC',
    model: 'Hierarchical STL + Leaves',
    quality: 'Self-adversarial review (6-axis)',
    unique: 'Broadcast protocol, Slack remote control, Pantheon, Desktop pet',
    complexity: 'Extreme (112KB prompt)',
    icon: '🐝',
  },
  {
    name: 'Metaswarm',
    author: 'dsifry',
    model: 'Recursive swarm-of-swarms',
    quality: 'Adversarial + 5-person review',
    unique: '18 specialized agent roles, JSONL knowledge base, Cross-model support',
    complexity: 'Extreme (18 roles × 9 stages)',
    icon: '🏭',
  },
  {
    name: 'Ralph',
    author: 'frankbria',
    model: 'Single agent serial loop',
    quality: 'Dual exit gates',
    unique: 'Extreme simplicity — one Bash script',
    complexity: 'Minimal',
    icon: '🔁',
  },
  {
    name: 'Swarm Orchestration',
    author: 'Kieran Klaassen',
    model: 'Hub-and-spoke peer coordination',
    quality: 'None specified',
    unique: 'Teaching-oriented Gist',
    complexity: 'Moderate',
    icon: '📘',
  },
]

export function ReferencesSection() {
  const { lang } = useLang()
  return (
    <SectionWrapper
      id="references"
      title={t('refs.title', lang)}
      subtitle={t('refs.subtitle', lang)}
    >
      {/* References */}
      <h3 className="font-mono text-lg text-gold mb-6">Reference Sources</h3>
      <div className="space-y-8 mb-16">
        {references.map((cat, ci) => (
          <motion.div
            key={cat.category}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: ci * 0.1 }}
          >
            <h4 className="font-mono text-sm text-text-secondary mb-3">{cat.category}</h4>
            <div className="grid md:grid-cols-2 gap-3">
              {cat.items.map((item) => (
                <div
                  key={item.name}
                  className="bg-bg-card border border-border rounded-lg p-4"
                >
                  <div className="flex items-center gap-2 mb-1">
                    <span className="font-mono text-sm text-text-primary font-bold">{item.name}</span>
                    {item.attributed ? (
                      <span className="text-[10px] bg-neon-green/10 text-neon-green px-1.5 py-0.5 rounded font-mono">ATTRIBUTED</span>
                    ) : (
                      <span className="text-[10px] bg-neon-red/10 text-neon-red px-1.5 py-0.5 rounded font-mono">UNATTRIBUTED</span>
                    )}
                  </div>
                  <p className="text-xs text-text-muted">{item.author}</p>
                  <p className="text-xs text-text-secondary mt-1">{item.desc}</p>
                </div>
              ))}
            </div>
          </motion.div>
        ))}
      </div>

      {/* Competitor comparison */}
      <h3 className="font-mono text-lg text-text-primary mb-6">Competitive Landscape</h3>
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border">
              <th className="text-left py-3 px-4 font-mono text-text-muted text-xs">Project</th>
              <th className="text-left py-3 px-4 font-mono text-text-muted text-xs">Execution Model</th>
              <th className="text-left py-3 px-4 font-mono text-text-muted text-xs">Quality Assurance</th>
              <th className="text-left py-3 px-4 font-mono text-text-muted text-xs">Unique Features</th>
              <th className="text-left py-3 px-4 font-mono text-text-muted text-xs">Complexity</th>
            </tr>
          </thead>
          <tbody>
            {competitors.map((c, i) => (
              <motion.tr
                key={c.name}
                initial={{ opacity: 0, x: -20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.1 }}
                className="border-b border-border/50 hover:bg-bg-card/50"
              >
                <td className="py-3 px-4">
                  <span className="mr-2">{c.icon}</span>
                  <span className={`font-mono font-bold ${i === 0 ? 'text-gold' : 'text-text-primary'}`}>{c.name}</span>
                  <div className="text-xs text-text-muted">{c.author}</div>
                </td>
                <td className="py-3 px-4 text-xs text-text-secondary">{c.model}</td>
                <td className="py-3 px-4 text-xs text-text-secondary">{c.quality}</td>
                <td className="py-3 px-4 text-xs text-text-secondary">{c.unique}</td>
                <td className="py-3 px-4 text-xs text-text-muted">{c.complexity}</td>
              </motion.tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Overall rating */}
      <motion.div
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
        className="mt-12 bg-bg-card border border-gold/20 rounded-lg p-8"
      >
        <h3 className="font-mono text-gold text-lg mb-6 text-center">Overall Evaluation — 8.4/10</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-4">
          {[
            { dim: 'Vision', score: 9 },
            { dim: 'Completeness', score: 10 },
            { dim: 'Originality', score: 8 },
            { dim: 'Prompt Quality', score: 9 },
            { dim: 'Practicality', score: 7 },
            { dim: 'Maintainability', score: 6 },
            { dim: 'Philosophy', score: 10 },
          ].map((item) => (
            <div key={item.dim} className="text-center">
              <div className="relative w-16 h-16 mx-auto mb-2">
                <svg viewBox="0 0 36 36" className="w-full h-full -rotate-90">
                  <circle cx={18} cy={18} r={15.9} fill="none" stroke="#30363d" strokeWidth={2} />
                  <motion.circle
                    cx={18}
                    cy={18}
                    r={15.9}
                    fill="none"
                    stroke="#f59e0b"
                    strokeWidth={2}
                    strokeDasharray={`${item.score * 10} 100`}
                    initial={{ strokeDasharray: '0 100' }}
                    whileInView={{ strokeDasharray: `${item.score * 10} 100` }}
                    viewport={{ once: true }}
                    transition={{ duration: 1, delay: 0.3 }}
                  />
                </svg>
                <span className="absolute inset-0 flex items-center justify-center text-sm font-mono text-gold font-bold">
                  {item.score}
                </span>
              </div>
              <span className="text-xs text-text-muted">{item.dim}</span>
            </div>
          ))}
        </div>
        <p className="text-sm text-text-secondary text-center mt-6">
          "An awe-inspiring prompt engineering work combining engineering rigor with philosophical depth."
        </p>
      </motion.div>
    </SectionWrapper>
  )
}
