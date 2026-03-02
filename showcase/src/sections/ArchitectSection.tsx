import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import { SectionWrapper } from '../components/SectionWrapper'
import { Annotatable } from '../components/Annotatable'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

// ─── Phase Data ───────────────────────────────────────────────────────────────
const phases = [
  {
    num: '0',
    name: 'Understand the Problem',
    icon: '🔭',
    colorBorder: 'border-blue-500',
    colorBg: 'bg-blue-500/10',
    textColor: 'text-blue-400',
    goal: 'Deeply understand before designing',
    steps: [
      'Restate the problem in your own words',
      'Identify the core tension — what makes this non-trivial?',
      'Ask 3-5 clarifying questions via AskUserQuestion',
      'Pin down: primary mission, users, hard constraints, success metrics, integration requirements',
    ],
    output: 'Clear problem statement + acceptance criteria',
    rule: 'DO NOT proceed until you have clear answers. Ambiguity here multiplies cost downstream.',
  },
  {
    num: '1',
    name: 'Exhaust the Design Space',
    icon: '🗺️',
    colorBorder: 'border-purple-500',
    colorBg: 'bg-purple-500/10',
    textColor: 'text-purple-400',
    goal: 'Map ALL possible approaches. Resist premature convergence.',
    steps: [
      'Broad research: industry leaders, unconventional approaches, adjacent domains',
      'Module decomposition: identify all abstract modules, list ≥3 implementation approaches per module',
      'Think: data flow, control flow, state management, failure modes, scaling axes',
      'Combinatorial exploration: enumerate 10-15 meaningfully distinct architecture paths',
      'Name each path vividly: "Event-Sourced Microservices", "Monolith-First", etc.',
    ],
    output: 'Architecture comparison table with 10-15 named paths',
    rule: 'Launch multiple Task agents in parallel for research. Read broadly: distributed systems, biology (coordination), information theory.',
  },
  {
    num: '2',
    name: 'Define Optimization Goals',
    icon: '🎯',
    colorBorder: 'border-gold',
    colorBg: 'bg-gold/10',
    textColor: 'text-gold',
    goal: 'Quantify what "good" means before comparing paths',
    steps: [
      'Primary metric: the ONE thing to optimize (e.g., p99 latency < 100ms)',
      'Secondary metrics: important but can trade off (e.g., cost < $X/month)',
      'Constraints: hard limits, non-negotiable (e.g., must run on AWS)',
      'Anti-goals: explicitly what we\'re NOT optimizing (e.g., no multi-region needed)',
    ],
    output: 'OPTIMIZATION GOALS spec: PRIMARY / SECONDARY / CONSTRAINTS / ANTI-GOALS',
    rule: 'Without quantified goals, "comparison" is just opinion. The spec IS the decision framework.',
  },
  {
    num: '3',
    name: 'Debate & Compare',
    icon: '⚔️',
    colorBorder: 'border-red-500',
    colorBg: 'bg-red-500/10',
    textColor: 'text-red-400',
    goal: 'Rigorously compare remaining paths against goals',
    steps: [
      'Elimination round: remove paths that violate hard constraints (with WHY)',
      'Deep comparison table: 12 dimensions (primary metric, complexity, team skills, time-to-MVP, scaling ceiling, key risk...)',
      'Devil\'s Advocate for each path: strongest FOR argument, strongest AGAINST, fatal flaw scenario',
      'User selects 2-3 paths for deep-dive via AskUserQuestion',
    ],
    output: '12-dimension comparison table + Devil\'s Advocate analysis per path',
    rule: 'Never just argue — find the fatal flaw scenario for each path. One counterexample kills a universal claim.',
  },
  {
    num: '4',
    name: 'Detail & Stress-Test',
    icon: '🔬',
    colorBorder: 'border-cyan-500',
    colorBg: 'bg-cyan-500/10',
    textColor: 'text-cyan-400',
    goal: 'Go deep on selected paths. Surface every hidden decision.',
    steps: [
      'Full component diagram (Mermaid/ASCII): data flow, interfaces, protocols, error handling, observability',
      'Answer 20+ execution questions: deployment, zero-downtime, failure recovery, data migration, schema evolution...',
      'Authentication model, rate limiting, caching strategy, distributed transactions, testing strategy',
      'Capacity planning numbers, CI/CD pipeline, local development workflow',
      'User narrows to 1 final path via AskUserQuestion',
    ],
    output: 'Full architecture design for each path + 20+ answered execution questions',
    rule: 'Each question is a future crisis averted. Missing an execution question = a production incident waiting to happen.',
  },
  {
    num: '5',
    name: 'System Review (Multi-Round)',
    icon: '🛡️',
    colorBorder: 'border-green-500',
    colorBg: 'bg-green-500/10',
    textColor: 'text-green-400',
    goal: 'Bulletproof the chosen design through adversarial review',
    steps: [
      'At least 3 rounds of review. Each round: Challenge → Resolve → Adjust → Simplify',
      'Challenge Phase: ask 10 pointed questions about failure modes, blast radius, single points of failure',
      'Resolve Phase: document if handled well, propose adjustment if weak, flag if fundamental issue',
      'Simplify Phase: "What can we remove? What can we merge? Where are we over-engineering?"',
      'Continue until no new significant issues, and design is as simple as possible but no simpler',
    ],
    output: 'Battle-hardened design + list of all resolved issues + simplifications applied',
    rule: '"Simplicity is the ultimate sophistication." If two designs are equal, choose the simpler one.',
  },
  {
    num: '6',
    name: 'Final Presentation',
    icon: '📜',
    colorBorder: 'border-border',
    colorBg: 'bg-bg-card',
    textColor: 'text-text-secondary',
    goal: 'Present the final architecture clearly and beautifully',
    steps: [
      'Executive summary (3-5 sentences: what, why, how)',
      'System diagram: high-level overview, data flow, deployment topology (Mermaid)',
      'Key Decisions Log: what chosen, what rejected, and WHY',
      'Component Specification: responsibility, interfaces, dependencies, scaling strategy',
      'Operational Runbook: deployment procedure, key alerts, top-5 incident responses',
      'Risk Register: top 5 risks with likelihood, impact, mitigation',
      'Implementation Roadmap: Phase 1 MVP, Phase 2, Future deferred',
      'Appendix: Rejected Alternatives — brief summary of paths not taken and why',
    ],
    output: 'Comprehensive architecture document written to file in project directory',
    rule: 'The document is the deliverable. Code comes later. A great architecture doc prevents months of wrong-direction coding.',
  },
]

// ─── Comparison Table Dimensions ──────────────────────────────────────────────
const comparisonDimensions = [
  'Primary metric fit',
  'Secondary metrics',
  'Implementation complexity',
  'Operational complexity',
  'Team skill requirements',
  'Time to MVP',
  'Scaling ceiling',
  'Migration difficulty',
  'Key risk',
  'Key advantage',
]

// ─── Core Principles ─────────────────────────────────────────────────────────
const corePrinciples = [
  {
    icon: '🪨',
    title: 'Simplicity is the ultimate sophistication',
    desc: 'Exhaust the design space — then converge on the simplest solution that meets all goals. Complexity is a signal of unclear thinking.',
  },
  {
    icon: '🔬',
    title: 'First Principles over conventions',
    desc: 'Draw inspiration from physics, biology (evolutionary design), philosophy (Occam\'s razor), mathematics (elegance), and information theory (minimum description length).',
  },
  {
    icon: '⚖️',
    title: 'Quantify before comparing',
    desc: 'Without defined optimization goals, comparison is opinion. The OPTIMIZATION GOALS spec IS the decision framework — build it before Phase 3.',
  },
  {
    icon: '🗡️',
    title: 'Find the fatal flaw',
    desc: 'For every path, find the one scenario that makes it catastrophically wrong. One counterexample kills a universal claim. Rejection needs only one.',
  },
]

export function ArchitectSection() {
  const { lang } = useLang()
  const [activePhase, setActivePhase] = useState<number | null>(0)

  return (
    <SectionWrapper
      id="architect"
      title={t('architect.title', lang)}
      subtitle={t('architect.subtitle', lang)}
    >
      {/* Core Philosophy */}
      <div className="mb-12">
        <h3 className="text-lg font-mono text-gold mb-6">Core Philosophy</h3>
        <div className="grid md:grid-cols-2 gap-4">
          {corePrinciples.map((principle) => (
            <div key={principle.title} className="flex gap-3 bg-bg-card border border-border rounded-lg p-4">
              <span className="text-2xl flex-shrink-0">{principle.icon}</span>
              <div>
                <div className="font-mono text-sm font-bold text-text-primary mb-1">{principle.title}</div>
                <div className="text-xs text-text-secondary leading-relaxed">{principle.desc}</div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* /drive vs /architect Comparison */}
      <div className="mb-12">
        <h3 className="text-lg font-mono text-gold mb-4">/drive vs /architect — When to Use Each</h3>
        <div className="grid md:grid-cols-2 gap-4">
          <div className="border border-gold rounded-lg p-4 bg-gold/5 text-gold">
            <div className="font-mono font-bold text-lg mb-3">/drive</div>
            <div className="text-xs font-mono mb-2 opacity-80">USE WHEN:</div>
            <ul className="space-y-1 mb-4">
              {[
                'You know WHAT to build, not just HOW',
                'Task is implementation-focused',
                'Requirements are reasonably clear',
                'Speed of execution is the bottleneck',
              ].map((w) => (
                <li key={w} className="text-xs flex gap-1.5">
                  <span className="opacity-60 flex-shrink-0">✓</span>
                  <span className="opacity-90">{w}</span>
                </li>
              ))}
            </ul>
            <div className="border-t border-gold/20 pt-3">
              <div className="text-xs font-mono opacity-70 mb-1">OUTPUT:</div>
              <div className="text-xs opacity-90">Working code + tests + verified delivery</div>
            </div>
          </div>
          <div className="border border-purple-400 rounded-lg p-4 bg-purple-500/10 text-purple-400">
            <div className="font-mono font-bold text-lg mb-3">/architect</div>
            <div className="text-xs font-mono mb-2 opacity-80">USE WHEN:</div>
            <ul className="space-y-1 mb-4">
              {[
                'You need to decide WHAT to build first',
                'Multiple valid architectures exist',
                'Wrong choice = months of rework',
                'Stakeholder alignment needed before code',
              ].map((w) => (
                <li key={w} className="text-xs flex gap-1.5">
                  <span className="opacity-60 flex-shrink-0">✓</span>
                  <span className="opacity-90">{w}</span>
                </li>
              ))}
            </ul>
            <div className="border-t border-purple-400/20 pt-3">
              <div className="text-xs font-mono opacity-70 mb-1">OUTPUT:</div>
              <div className="text-xs opacity-90">Architecture document + decision log + implementation roadmap</div>
            </div>
          </div>
        </div>
        <div className="mt-3 bg-bg-card border border-border rounded p-3 text-xs text-text-muted">
          <span className="text-text-secondary">Pro tip: </span>Start with /architect when requirements are fuzzy. Once architecture is decided, execute with /drive. They're designed to chain together.
        </div>
      </div>

      {/* 7-Phase Process */}
      <div className="mb-12">
        <Annotatable
          id="architect-phases"
          label="/architect 7-Phase Design Process"
          sectionId="architect"
          sectionLabel="/architect"
        >
          <div>
            <h3 className="text-lg font-mono text-gold mb-2">7-Phase Design Process</h3>
            <p className="text-text-muted text-sm mb-6">Click a phase to see details. Each phase has a clear goal, steps, and output.</p>

            {/* Phase list */}
            <div className="space-y-2">
              {phases.map((phase, i) => (
                <div key={phase.num}>
                  <motion.div
                    className={`flex items-center gap-3 cursor-pointer border rounded-lg p-3 transition-all ${
                      activePhase === i
                        ? `${phase.colorBorder} ${phase.colorBg} border-2`
                        : 'border-border bg-bg-card hover:border-border'
                    }`}
                    onClick={() => setActivePhase(activePhase === i ? null : i)}
                    whileHover={{ x: 2 }}
                  >
                    {/* Phase number */}
                    <div className={`flex-shrink-0 w-8 h-8 rounded-full border-2 flex items-center justify-center font-mono font-bold text-sm ${phase.colorBorder} ${phase.colorBg} ${phase.textColor}`}>
                      {phase.num}
                    </div>

                    {/* Icon + name */}
                    <div className="flex items-center gap-2 flex-1">
                      <span className="text-lg">{phase.icon}</span>
                      <div>
                        <div className="font-mono text-sm font-bold text-text-primary">Phase {phase.num}: {phase.name}</div>
                        <div className={`text-xs ${phase.textColor}`}>{phase.goal}</div>
                      </div>
                    </div>

                    {/* Output badge */}
                    <div className="hidden md:block text-xs text-text-muted font-mono bg-bg-surface px-2 py-1 rounded border border-border max-w-48 truncate">
                      → {phase.output.substring(0, 40)}{phase.output.length > 40 ? '...' : ''}
                    </div>

                    {/* Expand icon */}
                    <div className={`text-xs font-mono transition-transform ${activePhase === i ? 'rotate-90' : ''} text-text-muted`}>›</div>
                  </motion.div>

                  {/* Expanded detail */}
                  <AnimatePresence>
                    {activePhase === i && (
                      <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="overflow-hidden"
                      >
                        <div className="ml-11 mt-1 mb-2 bg-bg-card border border-border rounded-lg p-4">
                          <div className="grid md:grid-cols-3 gap-4">
                            {/* Steps */}
                            <div className="md:col-span-2">
                              <div className="text-xs font-mono text-text-muted mb-2">STEPS:</div>
                              <ul className="space-y-1.5">
                                {phase.steps.map((step, j) => (
                                  <li key={j} className="flex gap-2 text-xs">
                                    <span className={`flex-shrink-0 font-mono ${phase.textColor}`}>{j + 1}.</span>
                                    <span className="text-text-secondary leading-relaxed">{step}</span>
                                  </li>
                                ))}
                              </ul>
                            </div>

                            {/* Output + Rule */}
                            <div className="space-y-3">
                              <div>
                                <div className="text-xs font-mono text-text-muted mb-1">OUTPUT:</div>
                                <div className={`text-xs ${phase.textColor} ${phase.colorBg} border ${phase.colorBorder}/30 rounded p-2`}>
                                  {phase.output}
                                </div>
                              </div>
                              <div>
                                <div className="text-xs font-mono text-text-muted mb-1">RULE:</div>
                                <div className="text-xs text-text-secondary italic leading-relaxed">
                                  "{phase.rule}"
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </motion.div>
                    )}
                  </AnimatePresence>
                </div>
              ))}
            </div>
          </div>
        </Annotatable>
      </div>

      {/* Phase 3 Comparison Table Preview */}
      <div className="mb-12">
        <h3 className="text-lg font-mono text-gold mb-2">Phase 3: The 12-Dimension Comparison</h3>
        <p className="text-text-muted text-sm mb-4">After elimination, remaining 3-5 paths get rigorously compared across these dimensions:</p>
        <div className="bg-bg-card border border-border rounded-lg overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-xs font-mono">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left p-2 pl-3 text-text-muted bg-bg-surface">Dimension</th>
                  <th className="text-center p-2 text-gold/80">Path A</th>
                  <th className="text-center p-2 text-blue-400/80">Path B</th>
                  <th className="text-center p-2 text-purple-400/80">Path C</th>
                </tr>
              </thead>
              <tbody>
                {comparisonDimensions.map((dim, i) => (
                  <tr key={dim} className={`border-b border-border/50 ${i % 2 === 0 ? '' : 'bg-bg-surface/30'}`}>
                    <td className="p-2 pl-3 text-text-secondary">{dim}</td>
                    <td className="p-2 text-center text-text-muted">—</td>
                    <td className="p-2 text-center text-text-muted">—</td>
                    <td className="p-2 text-center text-text-muted">—</td>
                  </tr>
                ))}
                <tr className="bg-red-500/5">
                  <td className="p-2 pl-3 text-red-400 font-bold">Fatal flaw scenario</td>
                  <td className="p-2 text-center text-text-muted">—</td>
                  <td className="p-2 text-center text-text-muted">—</td>
                  <td className="p-2 text-center text-text-muted">—</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div className="p-3 border-t border-border text-xs text-text-muted">
            + Devil's Advocate for each path: strongest FOR, strongest AGAINST, fatal flaw scenario. User selects 2-3 for Phase 4 deep-dive.
          </div>
        </div>
      </div>

      {/* Phase 5 — Multi-Round Review Protocol */}
      <div className="mb-12">
        <h3 className="text-lg font-mono text-gold mb-2">Phase 5: Adversarial Review Loop</h3>
        <p className="text-text-muted text-sm mb-4">At least 3 review rounds. Each round follows a strict 4-step protocol:</p>
        <div className="grid md:grid-cols-4 gap-3">
          {[
            {
              step: '1',
              name: 'Challenge',
              icon: '🎯',
              desc: 'Ask 10 pointed questions: failure modes, blast radius, single points of failure, can this be simplified?',
              colorBorder: 'border-red-500/40',
              colorText: 'text-red-400',
            },
            {
              step: '2',
              name: 'Resolve',
              icon: '🔧',
              desc: 'For each question: document if handled well, propose adjustment if weak, flag if fundamental issue.',
              colorBorder: 'border-yellow-500/40',
              colorText: 'text-yellow-400',
            },
            {
              step: '3',
              name: 'Adjust',
              icon: '🔄',
              desc: 'Update the design based on findings. No issue left unresolved before next round begins.',
              colorBorder: 'border-blue-500/40',
              colorText: 'text-blue-400',
            },
            {
              step: '4',
              name: 'Simplify',
              icon: '✂️',
              desc: '"What can we remove? What can we merge? Where are we over-engineering?" Always ask this last.',
              colorBorder: 'border-green-500/40',
              colorText: 'text-green-400',
            },
          ].map((step) => (
            <div key={step.step} className={`border ${step.colorBorder} bg-bg-card rounded-lg p-3`}>
              <div className="flex items-center gap-2 mb-2">
                <span className="text-xl">{step.icon}</span>
                <div>
                  <div className="text-xs font-mono text-text-muted">Step {step.step}</div>
                  <div className={`font-mono text-sm font-bold ${step.colorText}`}>{step.name}</div>
                </div>
              </div>
              <div className="text-xs text-text-secondary leading-relaxed">{step.desc}</div>
            </div>
          ))}
        </div>
        <div className="mt-3 bg-bg-card border border-border rounded p-3 text-xs text-text-muted">
          <span className="text-green-400">STOP CONDITION: </span>No new significant issues found AND design is as simple as possible but no simpler. Usually 3-5 rounds.
        </div>
      </div>

      {/* Phase 6 — Output Document Structure */}
      <div className="mb-8">
        <h3 className="text-lg font-mono text-gold mb-2">Phase 6: Architecture Document</h3>
        <p className="text-text-muted text-sm mb-4">Written to a file in the project directory. Contains 8 sections:</p>
        <div className="grid md:grid-cols-2 gap-2">
          {[
            { num: '1', title: 'Executive Summary', desc: '3-5 sentences: what, why, how', icon: '📋' },
            { num: '2', title: 'System Diagram', desc: 'High-level overview + data flow + deployment topology (Mermaid)', icon: '🗺️' },
            { num: '3', title: 'Key Decisions Log', desc: 'What chosen, what rejected, and WHY for every major decision', icon: '📖' },
            { num: '4', title: 'Component Specification', desc: 'Responsibility, interfaces, dependencies, scaling strategy per component', icon: '🧩' },
            { num: '5', title: 'Operational Runbook', desc: 'Deployment procedure, key alerts, top-5 incident response', icon: '🚨' },
            { num: '6', title: 'Risk Register', desc: 'Top 5 risks: likelihood, impact, mitigation strategy', icon: '⚠️' },
            { num: '7', title: 'Implementation Roadmap', desc: 'Phase 1 MVP, Phase 2 extensions, Future deferrals', icon: '🛣️' },
            { num: '8', title: 'Appendix: Rejected Alternatives', desc: 'Brief summary of paths not taken and WHY — prevents re-litigating decisions', icon: '🗑️' },
          ].map((section) => (
            <div key={section.num} className="flex gap-3 bg-bg-card border border-border rounded-lg p-3">
              <span className="text-xl flex-shrink-0">{section.icon}</span>
              <div>
                <div className="flex items-center gap-2">
                  <span className="font-mono text-xs text-text-muted">{section.num}.</span>
                  <span className="font-mono text-sm font-bold text-text-primary">{section.title}</span>
                </div>
                <div className="text-xs text-text-secondary mt-0.5">{section.desc}</div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Process Rules */}
      <div className="bg-bg-card border border-gold/20 rounded-lg p-5">
        <div className="text-xs font-mono text-gold mb-3">7 PROCESS RULES</div>
        <div className="grid md:grid-cols-2 gap-2">
          {[
            'Use TodoWrite to track progress through all phases',
            'Use Task agents liberally for parallel research in Phase 1',
            'Use AskUserQuestion at EVERY decision point — never assume',
            'Write outputs to files in the project directory for persistence',
            'Be brutally honest about trade-offs. No hand-waving.',
            'Simplicity wins. If two designs are equal, choose the simpler one.',
            'Challenge yourself: "Is there a simpler way?" before every recommendation',
          ].map((rule, i) => (
            <div key={i} className="flex gap-2 text-xs">
              <span className="text-gold font-mono flex-shrink-0">{i + 1}.</span>
              <span className="text-text-secondary">{rule}</span>
            </div>
          ))}
        </div>
      </div>
    </SectionWrapper>
  )
}
