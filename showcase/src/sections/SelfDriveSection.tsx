import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import { SectionWrapper } from '../components/SectionWrapper'
import { Annotatable } from '../components/Annotatable'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

// ─── Loop Step Component ─────────────────────────────────────────────────────
const loopSteps = [
  {
    icon: '⚡',
    label: 'EXECUTE',
    sublabel: 'Mission via full /drive',
    color: 'border-gold bg-gold/10 text-gold',
    desc: 'Every task runs through the complete Drive Mode protocol — Phase 0 recon, Phase T team, Phase 1 wave map, The Loop, and Mission Complete review. No shortcuts. Quality each iteration earns user trust to keep the loop running.',
  },
  {
    icon: '🔍',
    label: 'ANALYZE',
    sublabel: '3 parallel scouts',
    color: 'border-blue-500 bg-blue-500/10 text-blue-400',
    desc: '3 scouts launched in parallel: scout-quality (TODOs, lint, type errors, dead code), scout-architecture (refactoring, DX, performance), scout-features (incomplete features, UX gaps, missing error handling). Mandatory — never guess.',
  },
  {
    icon: '📋',
    label: 'PROPOSE',
    sublabel: 'Top 3 via Impact×Effort',
    color: 'border-purple-500 bg-purple-500/10 text-purple-400',
    desc: 'Apply the Impact × Effort matrix. Rank: High Impact + Low Effort first (quick wins), High Impact + High Effort second (strategic), Low Impact + Low Effort third (polish). Skip Low Impact + High Effort entirely.',
  },
  {
    icon: '📡',
    label: 'SEND',
    sublabel: 'Proposals via Slack',
    color: 'border-green-500 bg-green-500/10 text-green-400',
    desc: 'Proposals sent to Slack DM. Primary: Socket Mode listener (real-time). Fallback 1: API polling every 30s. Fallback 2: AskUserQuestion in terminal. No timeout — user may reply in minutes or hours.',
  },
  {
    icon: '⏳',
    label: 'WAIT',
    sublabel: 'No timeout',
    color: 'border-yellow-500 bg-yellow-500/10 text-yellow-400',
    desc: 'Reply "1", "2", or "3" to pick a proposal. Any other text becomes a custom task. "stop", "exit", "结束" ends the session. Unclear replies get a prompt to retry. The loop is patient — you can respond hours later.',
  },
]

// ─── Scout Type Card ─────────────────────────────────────────────────────────
const scouts = [
  {
    name: 'scout-quality',
    icon: '🔬',
    color: 'bg-red-500/10',
    border: 'border-red-500/40',
    dimensions: [
      'TODOs and FIXMEs in codebase',
      'Dead code & unused exports',
      'Test coverage gaps',
      'Lint errors & warnings',
      'TypeScript type errors',
      'Missing error handling',
    ],
  },
  {
    name: 'scout-architecture',
    icon: '🏗️',
    color: 'bg-blue-500/10',
    border: 'border-blue-500/40',
    dimensions: [
      'Module boundary health',
      'Refactoring opportunities',
      'Performance bottlenecks',
      'Developer experience gaps',
      'Build speed issues',
      'Abstraction quality',
    ],
  },
  {
    name: 'scout-features',
    icon: '✨',
    color: 'bg-green-500/10',
    border: 'border-green-500/40',
    dimensions: [
      'Incomplete feature implementations',
      'UX gaps & edge case handling',
      'Missing documentation',
      'Error state coverage',
      'Accessibility issues',
      'User flow completeness',
    ],
  },
]

// ─── Diversity Categories ─────────────────────────────────────────────────────
const diversityCategories = [
  { num: '①', label: 'Code Quality', desc: 'Bugs, tech debt, type safety, linting, dead code', color: 'bg-red-500/15 text-red-400 border-red-500/30' },
  { num: '②', label: 'Feature Completeness', desc: 'Missing features, incomplete implementations, edge cases', color: 'bg-blue-500/15 text-blue-400 border-blue-500/30' },
  { num: '③', label: 'Testing', desc: 'Coverage gaps, edge case tests, integration tests', color: 'bg-purple-500/15 text-purple-400 border-purple-500/30' },
  { num: '④', label: 'Performance', desc: 'Profiling, optimization, caching, bundle size', color: 'bg-yellow-500/15 text-yellow-400 border-yellow-500/30' },
  { num: '⑤', label: 'Developer Experience', desc: 'Tooling, build speed, documentation, onboarding', color: 'bg-cyan-500/15 text-cyan-400 border-cyan-500/30' },
  { num: '⑥', label: 'Architecture', desc: 'Refactoring, better abstractions, cleaner module boundaries', color: 'bg-green-500/15 text-green-400 border-green-500/30' },
]

// ─── Slack Architecture ───────────────────────────────────────────────────────
const slackLevels = [
  {
    level: 'Primary',
    method: 'Socket Mode',
    icon: '⚡',
    desc: 'WebSocket-based real-time listener. Background daemon started before waiting. Replies land instantly as file ~/.claude/slack-reply.txt.',
    colorBorder: 'border-gold',
    colorBg: 'bg-gold/5',
    colorText: 'text-gold',
  },
  {
    level: 'Fallback 1',
    method: 'API Polling',
    icon: '🔄',
    desc: 'If Socket Mode fails, falls back to polling the Slack API every 30 seconds. Slightly slower but reliable on networks that block WebSockets.',
    colorBorder: 'border-blue-500/40',
    colorBg: 'bg-blue-500/5',
    colorText: 'text-blue-400',
  },
  {
    level: 'Fallback 2',
    method: 'Terminal AskUserQuestion',
    icon: '🖥️',
    desc: 'If Slack is not configured or all comms fail, falls back to AskUserQuestion in terminal. Requires user to be present but fully functional.',
    colorBorder: 'border-border',
    colorBg: 'bg-bg-card',
    colorText: 'text-text-secondary',
  },
]

// ─── Anti-Patterns ────────────────────────────────────────────────────────────
const antiPatterns = [
  { id: 1, text: 'Proposing vague tasks — "improve performance" is useless. Name the file, the function, the measured latency.', icon: '🌫️' },
  { id: 2, text: 'Repeating rejected proposals — if user chose something else, the old proposal lost. Don\'t re-surface it next round.', icon: '🔁' },
  { id: 3, text: 'Three tasks in the same category — use the diversity rule. Vary across quality, features, testing, performance, DX, architecture.', icon: '🗂️' },
  { id: 4, text: 'Skipping project analysis — the 3 scouts are mandatory. Bad proposals waste the user\'s time and erode trust in the loop.', icon: '🚫' },
  { id: 5, text: 'Cutting corners on mission execution — every task runs through FULL /drive. No abbreviated versions. No skipped phases.', icon: '✂️' },
  { id: 6, text: 'Proposing mega-tasks — "Rewrite the entire API layer" is not one Drive session. Propose the first concrete step.', icon: '🐘' },
  { id: 7, text: 'Ignoring user patterns — if they consistently pick performance tasks, weight future proposals toward performance.', icon: '👁️' },
  { id: 8, text: 'Rushing proposals — take the time to actually analyze. Quality proposals > fast proposals.', icon: '⏩' },
  { id: 9, text: 'Polling aggressively — the user has a life. The loop is patient by design. Don\'t spam reminders.', icon: '📣' },
  { id: 10, text: 'Forgetting to send debrief — after each mission, send the condensed debrief to Slack BEFORE starting analysis.', icon: '📝' },
  { id: 11, text: 'Losing session state — track mission_count and which proposals were accepted/rejected across the session.', icon: '🧠' },
]

// ─── Proposal Message Format ──────────────────────────────────────────────────
const proposalFormat = `🤖 MISSION COMPLETE: {mission-name}

✅ {one-line summary of what was accomplished}

━━━━━━━━━━━━━━━━━━━━
📋 Next improvements:
━━━━━━━━━━━━━━━━━━━━

1️⃣ {Task 1 Title}
   {Why it matters}
   Complexity: {Low/Medium/High}

2️⃣ {Task 2 Title}
   {Why it matters}
   Complexity: {Low/Medium/High}

3️⃣ {Task 3 Title}
   {Why it matters}
   Complexity: {Low/Medium/High}

━━━━━━━━━━━━━━━━━━━━
Reply: 1, 2, 3, or describe a new task
Reply "stop" to end self-drive`

export function SelfDriveSection() {
  const { lang } = useLang()
  const [activeStep, setActiveStep] = useState<number | null>(null)
  const [showFormat, setShowFormat] = useState(false)

  return (
    <SectionWrapper
      id="self-drive"
      title={t('selfdrive.title', lang)}
      subtitle={t('selfdrive.subtitle', lang)}
    >
      {/* Perpetual Loop Flowchart */}
      <div className="mb-16">
        <h3 className="text-lg font-mono text-gold mb-6">The Perpetual Loop — 永动机</h3>
        <Annotatable
          id="selfdrive-loop"
          label="Self-Drive Perpetual Loop"
          sectionId="self-drive"
          sectionLabel="Self-Drive Mode"
        >
          <div className="relative">
            {/* Loop container */}
            <div className="flex flex-wrap items-center justify-center gap-2">
              {loopSteps.map((step, i) => (
                <div key={step.label} className="flex items-center gap-2">
                  <motion.div
                    className={`cursor-pointer rounded-lg border-2 p-3 transition-all w-36 ${step.color} ${
                      activeStep === i ? 'ring-2 ring-white/30 shadow-lg scale-105' : 'hover:scale-102'
                    }`}
                    whileHover={{ scale: 1.05 }}
                    onClick={() => setActiveStep(activeStep === i ? null : i)}
                  >
                    <div className="text-2xl mb-1 text-center">{step.icon}</div>
                    <div className="font-mono font-bold text-sm text-center">{step.label}</div>
                    <div className="text-xs text-center opacity-80 mt-0.5">{step.sublabel}</div>
                  </motion.div>

                  {/* Arrow between steps */}
                  {i < loopSteps.length - 1 && (
                    <div className="text-text-muted font-mono text-lg">→</div>
                  )}
                </div>
              ))}

              {/* Loop-back arrow indicator */}
              <div className="w-full flex items-center justify-center mt-3">
                <div className="flex items-center gap-2 text-text-muted text-sm font-mono">
                  <div className="w-16 border-t border-dashed border-text-muted"></div>
                  <span>↩ loops forever until "stop"</span>
                  <div className="w-16 border-t border-dashed border-text-muted"></div>
                </div>
              </div>
            </div>

            {/* Expanded step detail */}
            <AnimatePresence>
              {activeStep !== null && (
                <motion.div
                  initial={{ opacity: 0, y: -8 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -8 }}
                  className="mt-4 bg-bg-card border border-border rounded-lg p-4"
                >
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-xl">{loopSteps[activeStep].icon}</span>
                    <span className="font-mono font-bold text-text-primary">{loopSteps[activeStep].label}</span>
                    <span className="text-text-muted text-sm">— {loopSteps[activeStep].sublabel}</span>
                  </div>
                  <p className="text-text-secondary text-sm leading-relaxed">{loopSteps[activeStep].desc}</p>
                </motion.div>
              )}
            </AnimatePresence>
          </div>
        </Annotatable>

        {/* Entry behavior */}
        <div className="mt-6 bg-bg-card border border-border rounded-lg p-4">
          <div className="text-xs font-mono text-gold mb-2">ENTRY BEHAVIOR</div>
          <div className="grid md:grid-cols-2 gap-3 text-sm">
            <div className="flex gap-2">
              <span className="text-text-muted font-mono flex-shrink-0">/self-drive</span>
              <span className="text-text-secondary">→ skip to Project Analysis + Proposal immediately</span>
            </div>
            <div className="flex gap-2">
              <span className="text-text-muted font-mono flex-shrink-0">/self-drive {'{task}'}</span>
              <span className="text-text-secondary">→ execute task via /drive first, then enter loop</span>
            </div>
          </div>
        </div>
      </div>

      {/* Three Scout Types */}
      <div className="mb-16">
        <h3 className="text-lg font-mono text-gold mb-2">3 Parallel Analysis Scouts</h3>
        <p className="text-text-muted text-sm mb-6">Launched simultaneously — never serial. Gives 3× analysis coverage in the time of one scout.</p>
        <div className="grid md:grid-cols-3 gap-4">
          {scouts.map((scout) => (
            <motion.div
              key={scout.name}
              className={`${scout.color} border ${scout.border} rounded-lg p-4`}
              whileHover={{ scale: 1.02 }}
            >
              <div className="text-2xl mb-2">{scout.icon}</div>
              <div className="font-mono text-sm font-bold text-text-primary mb-3">{scout.name}</div>
              <div className="space-y-1">
                {scout.dimensions.map((dim) => (
                  <div key={dim} className="flex items-start gap-1.5">
                    <span className="text-text-muted text-xs mt-0.5">•</span>
                    <span className="text-text-secondary text-xs">{dim}</span>
                  </div>
                ))}
              </div>
            </motion.div>
          ))}
        </div>
      </div>

      {/* Impact × Effort Matrix + Diversity */}
      <div className="grid md:grid-cols-2 gap-8 mb-16">
        {/* Impact × Effort Matrix */}
        <div>
          <h3 className="text-lg font-mono text-gold mb-2">Impact × Effort Matrix</h3>
          <p className="text-text-muted text-sm mb-4">How to rank scout findings into proposals</p>
          <Annotatable
            id="selfdrive-matrix"
            label="Impact × Effort Ranking Matrix"
            sectionId="self-drive"
            sectionLabel="Self-Drive Mode"
          >
            <div className="border border-border rounded-lg overflow-hidden">
              {/* Header row */}
              <div className="grid grid-cols-3 bg-bg-card border-b border-border">
                <div className="p-2 text-xs text-text-muted font-mono border-r border-border"></div>
                <div className="p-2 text-xs text-text-muted font-mono text-center border-r border-border">Low Effort</div>
                <div className="p-2 text-xs text-text-muted font-mono text-center">High Effort</div>
              </div>

              {/* High Impact row */}
              <div className="grid grid-cols-3 border-b border-border">
                <div className="p-2 text-xs font-mono text-text-secondary bg-bg-card border-r border-border flex items-center">High Impact</div>
                <div className="p-3 border-r border-border">
                  <div className="bg-gold/20 border border-gold rounded p-2 text-center">
                    <div className="text-gold font-bold text-xs font-mono">⭐ FIRST</div>
                    <div className="text-gold/70 text-xs">Quick Wins</div>
                  </div>
                </div>
                <div className="p-3">
                  <div className="bg-blue-500/20 border border-blue-400 rounded p-2 text-center">
                    <div className="text-blue-300 font-bold text-xs font-mono">🎯 SECOND</div>
                    <div className="text-blue-300/70 text-xs">Strategic</div>
                  </div>
                </div>
              </div>

              {/* Low Impact row */}
              <div className="grid grid-cols-3">
                <div className="p-2 text-xs font-mono text-text-secondary bg-bg-card border-r border-border flex items-center">Low Impact</div>
                <div className="p-3 border-r border-border">
                  <div className="bg-green-500/20 border border-green-400 rounded p-2 text-center">
                    <div className="text-green-300 font-bold text-xs font-mono">💡 THIRD</div>
                    <div className="text-green-300/70 text-xs">Nice-to-have</div>
                  </div>
                </div>
                <div className="p-3">
                  <div className="bg-red-500/10 border border-red-500/40 rounded p-2 text-center opacity-60">
                    <div className="text-red-400 font-bold text-xs font-mono">✗ SKIP</div>
                    <div className="text-red-400/70 text-xs">Not worth it</div>
                  </div>
                </div>
              </div>
            </div>
          </Annotatable>
        </div>

        {/* 6 Diversity Categories */}
        <div>
          <h3 className="text-lg font-mono text-gold mb-2">6 Diversity Categories</h3>
          <p className="text-text-muted text-sm mb-4">Never propose 3 tasks from the same category. Rotate across rounds.</p>
          <div className="space-y-2">
            {diversityCategories.map((cat) => (
              <div key={cat.num} className={`border rounded-lg px-3 py-2 flex items-start gap-3 ${cat.color}`}>
                <span className="font-mono font-bold text-sm flex-shrink-0">{cat.num}</span>
                <div>
                  <div className="font-mono text-sm font-bold">{cat.label}</div>
                  <div className="text-xs opacity-80 mt-0.5">{cat.desc}</div>
                </div>
              </div>
            ))}
          </div>
          <div className="mt-3 bg-bg-card border border-border rounded p-3 text-xs text-text-muted">
            <span className="text-gold font-mono">RULE: </span>If last round covered ①②③ → this round pulls from ④⑤⑥ (or wherever scouts found highest impact)
          </div>
        </div>
      </div>

      {/* Slack Remote Control Architecture */}
      <div className="mb-16">
        <h3 className="text-lg font-mono text-gold mb-2">Slack Remote Control Architecture</h3>
        <p className="text-text-muted text-sm mb-6">Self-Drive bridges the gap between AI execution and human oversight asynchronously — reply from anywhere, any time.</p>
        <div className="grid md:grid-cols-3 gap-4 mb-6">
          {slackLevels.map((level, i) => (
            <div key={level.level} className="relative">
              <div className={`border rounded-lg p-4 h-full ${level.colorBorder} ${level.colorBg}`}>
                <div className="flex items-center gap-2 mb-2">
                  <span className="text-lg">{level.icon}</span>
                  <div>
                    <div className={`text-xs font-mono ${level.colorText}`}>{level.level}</div>
                    <div className="font-mono text-sm font-bold text-text-primary">{level.method}</div>
                  </div>
                </div>
                <p className="text-text-secondary text-xs leading-relaxed">{level.desc}</p>
              </div>
              {i < slackLevels.length - 1 && (
                <div className="hidden md:block absolute -right-2 top-1/2 -translate-y-1/2 z-10 text-text-muted text-sm">→</div>
              )}
            </div>
          ))}
        </div>

        {/* Config */}
        <div className="bg-bg-card border border-border rounded-lg p-4">
          <div className="text-xs font-mono text-gold mb-2">CONFIG: ~/.claude/self-drive.json</div>
          <pre className="text-xs text-text-secondary font-mono overflow-x-auto">{`{
  "slack": {
    "bot_token": "xoxb-...",   // Slack bot token
    "app_token": "xapp-...",   // App-level token (for Socket Mode)
    "channel_id": "D..."       // Your DM channel ID
  }
}`}</pre>
        </div>
      </div>

      {/* Proposal Message Format */}
      <div className="mb-16">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h3 className="text-lg font-mono text-gold">Proposal Message Format</h3>
            <p className="text-text-muted text-sm mt-1">Exact format sent after every completed mission</p>
          </div>
          <button
            onClick={() => setShowFormat(!showFormat)}
            className="text-xs font-mono text-gold border border-gold/40 rounded px-3 py-1.5 hover:bg-gold/10 transition-colors"
          >
            {showFormat ? '▲ hide' : '▼ show format'}
          </button>
        </div>

        <AnimatePresence>
          {showFormat && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
              className="overflow-hidden"
            >
              <pre className="bg-bg-card border border-border rounded-lg p-4 text-xs font-mono text-text-secondary overflow-x-auto whitespace-pre">
                {proposalFormat}
              </pre>
            </motion.div>
          )}
        </AnimatePresence>

        {/* Adaptive learning callout */}
        <div className="mt-4 bg-bg-card border border-border rounded-lg p-4">
          <div className="text-xs font-mono text-gold mb-2">ADAPTIVE LEARNING WITHIN SESSIONS</div>
          <div className="grid md:grid-cols-2 gap-4 text-sm">
            <div>
              <div className="text-text-secondary font-mono text-xs mb-2">TRACK:</div>
              <ul className="space-y-1 text-text-muted text-xs">
                <li>• Which proposal number they pick (positional bias?)</li>
                <li>• Which categories they prefer repeatedly</li>
                <li>• When they give custom tasks instead (what patterns?)</li>
                <li>• Response latency (engaged vs. busy)</li>
              </ul>
            </div>
            <div>
              <div className="text-text-secondary font-mono text-xs mb-2">APPLY:</div>
              <ul className="space-y-1 text-text-muted text-xs">
                <li>• Rank proposals to match observed preferences</li>
                <li>• Put "most likely chosen" as option 1</li>
                <li>• Weight toward user's demonstrated interests</li>
                <li>• Align custom-task patterns into future proposals</li>
              </ul>
            </div>
          </div>
          <div className="mt-3 text-xs text-text-muted border-t border-border pt-3">
            <span className="text-text-secondary">Note:</span> Session-scoped only. Each new self-drive session starts fresh — no cross-session memory.
          </div>
        </div>
      </div>

      {/* Anti-Patterns */}
      <div className="mb-8">
        <Annotatable
          id="selfdrive-antipatterns"
          label="11 Self-Drive Anti-Patterns"
          sectionId="self-drive"
          sectionLabel="Self-Drive Mode"
        >
          <div>
            <h3 className="text-lg font-mono text-gold mb-2">11 Self-Drive Anti-Patterns</h3>
            <p className="text-text-muted text-sm mb-6">Behaviors that break the loop's trust and effectiveness</p>
            <div className="grid md:grid-cols-2 gap-3">
              {antiPatterns.map((ap) => (
                <motion.div
                  key={ap.id}
                  className="flex gap-3 bg-bg-card border border-red-500/20 rounded-lg p-3"
                  initial={{ opacity: 0, x: -8 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: ap.id * 0.04 }}
                >
                  <div className="flex-shrink-0">
                    <span className="text-red-500 font-mono text-xs font-bold">❌</span>
                  </div>
                  <div className="text-xs text-text-secondary leading-relaxed">
                    <span className="text-lg mr-1">{ap.icon}</span>
                    {ap.text}
                  </div>
                </motion.div>
              ))}
            </div>
          </div>
        </Annotatable>
      </div>

      {/* Mindset */}
      <div className="bg-bg-card border border-gold/20 rounded-lg p-5">
        <div className="text-xs font-mono text-gold mb-3">THE MINDSET</div>
        <p className="text-text-secondary text-sm leading-relaxed mb-4">
          You are not just executing tasks — you are <strong className="text-text-primary">continuously improving a project</strong>. Each loop iteration makes the codebase healthier, more robust, more complete.
        </p>
        <div className="grid md:grid-cols-4 gap-3">
          {[
            { icon: '👁️', text: 'See what they can\'t', sub: 'Find issues they haven\'t noticed' },
            { icon: '⚖️', text: 'Prioritize wisely', sub: 'Propose what matters most' },
            { icon: '🎯', text: 'Execute with quality', sub: 'Every mission at Drive standard' },
            { icon: '📡', text: 'Async communication', sub: 'Slack proposals from anywhere' },
          ].map((item) => (
            <div key={item.icon} className="text-center">
              <div className="text-2xl mb-1">{item.icon}</div>
              <div className="text-xs font-mono text-text-primary font-bold">{item.text}</div>
              <div className="text-xs text-text-muted mt-0.5">{item.sub}</div>
            </div>
          ))}
        </div>
        <div className="mt-4 text-center text-xs text-text-muted font-mono border-t border-border pt-3">
          "The loop is the hero. Each mission is one heartbeat. The project gets better with every beat."
        </div>
      </div>
    </SectionWrapper>
  )
}
