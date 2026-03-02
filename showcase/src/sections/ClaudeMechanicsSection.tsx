import { motion } from 'framer-motion'
import { SectionWrapper } from '../components/SectionWrapper'
import { useState } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'
import { Annotatable } from '../components/Annotatable'

interface MechanismCard {
  icon: string
  title: string
  subtitle: string
  color: string
  badge: string
  problem: string
  problemHint: string
  howDriveUses: string
  code: string
  codeLabel: string
  detail: string
}

const mechanisms: MechanismCard[] = [
  {
    icon: '📂',
    title: 'Slash Commands',
    subtitle: '~/.claude/commands/*.md',
    color: '#f59e0b',
    badge: 'SYSTEM',
    problem: 'A 112KB prompt can\'t live in your chat window — you\'d paste it every session and lose it when the session ends. You can\'t version-control it, iterate on it across sessions, or share it with your team. There was no native way to persist a complex system prompt as a reusable command.',
    problemHint: 'No way to persist system prompts',
    howDriveUses: 'The entire 112KB drive prompt lives in commands/drive.md. Claude Code reads it at `/drive` invocation and injects it as the system prompt for that session.',
    codeLabel: 'YAML frontmatter (drive.md)',
    code: `---
allowed-tools: Bash, Read, Write, Edit, Glob,
  Grep, Task, TaskCreate, TaskUpdate, TaskList,
  AskUserQuestion, TeamCreate, TeamDelete,
  SendMessage
description: High-drive autonomous execution mode
---

# Drive Mode
You are a Team Lead...`,
    detail: 'The allowed-tools frontmatter acts as a whitelist — Claude Code strips ALL other tools. Drive Skill only gets the tools it explicitly declares.',
  },
  {
    icon: '🏭',
    title: 'Agent Teams API',
    subtitle: 'TeamCreate / TeamDelete / SendMessage',
    color: '#3b82f6',
    badge: 'API',
    problem: 'Claude has one context window. Running a complex mission solo means one agent holds all the context and does all the work — sequentially. There was no built-in way to create multiple independent Claude instances that could coordinate on a shared mission with real-time communication.',
    problemHint: 'Single agent = single context window',
    howDriveUses: 'Drive creates a named team at Phase T. All teammates register into this team. SendMessage delivers progress reports back to the Team Lead. TeamDelete cleans up when mission completes.',
    codeLabel: 'Team lifecycle in drive execution',
    code: `// Phase T: Team Assembly
TeamCreate({ team_name: "auth-mission-v1" })

// Spawn teammates into team
Agent({ subagent_type: "general-purpose",
  team_name: "auth-mission-v1",
  name: "stl-auth",
  prompt: "[STL Domain Brief]..." })

// Teammate reports back
SendMessage({ type: "message",
  recipient: "team-lead",
  content: "Auth layer complete" })

// Mission complete cleanup
TeamDelete()`,
    detail: 'Requires CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 env var. Set in settings.json by the Drive Skill installer.',
  },
  {
    icon: '⚡',
    title: 'Task Tool (Subagents)',
    subtitle: 'Parallel task spawning within STLs',
    color: '#a855f7',
    badge: 'TOOL',
    problem: 'STLs need their own inner parallelism. But they can\'t use the Agent Teams API (that would create unbounded nesting and coordination overhead). They need a lightweight way to fork N subtasks simultaneously without creating a new team — and get the results back synchronously.',
    problemHint: 'STLs need lightweight inner parallelism',
    howDriveUses: 'Sub-Team-Leads use the Task tool (not Agent) to spawn their inner subagents. Task runs synchronously within the STL\'s context, returning when done. This enables STLs to run their own inner Fork-Join patterns.',
    codeLabel: 'STL inner parallelism pattern',
    code: `// Inside STL-auth's execution:
const [oauth, jwt, session] = await Promise.all([
  Task({
    description: "Implement OAuth2 PKCE flow",
    prompt: "...domain context...",
    model: "opus"  // ← ALL tasks run Opus
  }),
  Task({ description: "JWT token service",
    model: "opus" }),
  Task({ description: "Redis session store",
    model: "opus" })
])
// All 3 inner subagents ran in parallel`,
    detail: 'Model: "opus" is mandatory on every Task call — maximum reasoning for all subagents, no cost-cutting.',
  },
  {
    icon: '🪝',
    title: 'Hooks System',
    subtitle: 'Notification / Stop / TeammateIdle',
    color: '#22c55e',
    badge: 'HOOKS',
    problem: 'You kick off a 15-minute autonomous swarm and walk away. How do you know when it needs attention? How do you know it finished? Polling the terminal every 30 seconds defeats the purpose of autonomous execution. You need the system to reach out to YOU, not the other way around.',
    problemHint: 'Autonomous = invisible without notification',
    howDriveUses: 'The installer registers notify.sh as a hook for 3 event types. When any event fires, the hook receives a JSON payload and can trigger macOS banners, sound effects, and Body Zero desktop pet animations.',
    codeLabel: 'settings.json hook registration',
    code: `// ~/.claude/settings.json
{
  "hooks": {
    "Notification": [{ "hooks": [{
      "type": "command",
      "command": "~/.claude/notify.sh"
    }]}],
    "Stop": [{ "hooks": [{
      "type": "command",
      "command": "~/.claude/notify.sh"
    }]}],
    "TeammateIdle": [{ "hooks": [{
      "type": "command",
      "command": "~/.claude/notify.sh"
    }]}]
  }
}`,
    detail: 'notify.sh filters out subagent events (is_subagent=true) — only Team Lead events trigger notifications. No notification spam from inner subagents.',
  },
  {
    icon: '🔐',
    title: 'bypassPermissions Mode',
    subtitle: 'defaultMode: "bypassPermissions"',
    color: '#ef4444',
    badge: 'CONFIG',
    problem: 'With 8 parallel agents each doing 10 actions, that\'s 80 permission dialogs. Each one pauses ALL work until you click. The user becomes the bottleneck in a system designed for autonomy. One unattended permission dialog stalls the entire swarm indefinitely.',
    problemHint: '80 dialogs = user is the bottleneck',
    howDriveUses: 'Drive Skill sets defaultMode to bypassPermissions in settings.json. This lets the swarm write files, run bash, install packages, and execute tool calls without prompting — essential for autonomous multi-agent execution.',
    codeLabel: 'settings.json permissions',
    code: `// ~/.claude/settings.json
{
  "env": {
    "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"
  },
  "model": "opus",
  "permissions": {
    "allow": [
      "WebFetch(domain:github.com)"
    ],
    "defaultMode": "bypassPermissions"
  },
  "skipDangerousModePermissionPrompt": true
}`,
    detail: 'Without bypassPermissions, each of the 8+ parallel agents would trigger permission dialogs — making autonomous execution impossible.',
  },
  {
    icon: '💬',
    title: 'AskUserQuestion Tool',
    subtitle: 'Vivid option engineering for key decisions',
    color: '#06b6d4',
    badge: 'TOOL',
    problem: 'Claude\'s default behavior when uncertain: output a question as plain text and stop. The loop breaks. Session context is lost. Hours of momentum gone. There was no structured way to pause execution, ask a question with options, and resume — without breaking the autonomous loop.',
    problemHint: 'Plain-text questions break the loop',
    howDriveUses: 'Drive uses AskUserQuestion at Phase 0 to present the Mission Contract for confirmation, and at direction decision points. Options are engineered to be vivid — "Fortress mode" vs "Lightning mode", not just "A vs B".',
    codeLabel: 'Phase 0 Mission Contract confirmation',
    code: `AskUserQuestion({
  questions: [{
    question: "Mission Contract ready. Proceed?",
    header: "Go/No-go",
    options: [
      {
        label: "Launch the swarm",
        description: "All 3 phases confirmed..."
      },
      {
        label: "Revise scope",
        description: "Go back to briefing..."
      },
      {
        label: "Abort mission",
        description: "Cancel drive mode entirely"
      }
    ]
  }]
})`,
    detail: 'AskUserQuestion is the ONLY acceptable way to stop and ask the user. Using plain text to ask a question is explicitly listed as an anti-pattern.',
  },
  {
    icon: '🧠',
    title: 'In-Process Teammate Mode',
    subtitle: 'Shift+Down to cycle between agents',
    color: '#f97316',
    badge: 'UI',
    problem: 'Multi-agent execution is a black box. Something\'s running but you have no idea what 8 agents are doing. Is one stuck? Is it writing to the wrong file? Is it done and waiting? Autonomous doesn\'t have to mean invisible — but by default, all agent activity was hidden.',
    problemHint: 'Multi-agent = black box by default',
    howDriveUses: 'When Drive spawns teammates via the Agent tool (not Task), they appear as switchable panes in the Claude Code terminal. You can Shift+Down to see STL-frontend actively spawning inner subagents in real time.',
    codeLabel: 'Agent tool for visible teammates',
    code: `// Team Lead spawns STL as visible pane:
Agent({
  subagent_type: "general-purpose",
  team_name: "auth-mission-v1",
  name: "stl-frontend",       // ← appears in Shift+Down
  prompt: "You are STL-frontend.
    Own all UI tasks. Spawn 2-3 Task
    subagents for inner parallelism..."
})

// vs. inner subagents (not directly visible):
Task({ description: "...", model: "opus" })
// Task runs within STL context, invisible to user`,
    detail: 'This is why Drive uses a two-level hierarchy: STLs are Agent tool calls (visible panes), inner subagents are Task calls (invisible but fast). The terminal becomes a live war room.',
  },
  {
    icon: '🏛️',
    title: 'Context Injection via Files',
    subtitle: '~/.claude/pantheon/*.md loaded on demand',
    color: '#8b5cf6',
    badge: 'CONTEXT',
    problem: 'Loading all 400KB of Pantheon wisdom upfront burns the context window. 90% would be irrelevant noise. But loading zero context means Claude has no philosophical depth for architectural decisions. The question: how do you get context when you need it without paying the cost always?',
    problemHint: '400KB upfront wastes context window',
    howDriveUses: 'The Pantheon files (7 era files, ~400KB) live in ~/.claude/pantheon/. When Drive or any STL needs a thinker\'s perspective, they Read the relevant era file. Context is injected on demand, not loaded upfront — protecting the context window.',
    codeLabel: 'Pantheon invocation in practice',
    code: `// When facing architectural tradeoffs:
Read("~/.claude/pantheon/06_20世纪上半叶.md")
// → Loads von Neumann, Turing, Shannon, etc.

// System prompt tells Claude:
// "If you invoke von Neumann, quote a specific
//  principle. If you invoke Feynman, actually
//  apply the Feynman Technique to simplify
//  the problem. No name-dropping."

// Result: Reasoning grounded in specific
// methodologies, not vague attribution`,
    detail: 'All Pantheon files are symlinked from the repo — edit once, all Claude Code sessions see the change immediately.',
  },
]

const hookEvents = [
  { event: 'Notification', types: ['permission_prompt', 'elicitation_dialog', 'idle_prompt'], sound: 'permission-required.mp3 / waiting-for-input.mp3', bodyZero: 'permission / question / idle', needs_voice: true },
  { event: 'Stop', types: ['(task complete)'], sound: 'task-complete.mp3', bodyZero: 'complete', needs_voice: true },
  { event: 'TeammateIdle', types: ['(teammate finished)'], sound: 'teammate-idle.mp3', bodyZero: 'idle', needs_voice: false },
]

export function ClaudeMechanicsSection() {
  const { lang } = useLang()
  const [active, setActive] = useState<number>(0)
  const mech = mechanisms[active]

  return (
    <SectionWrapper
      id="mechanics"
      title={t('mechanics.title', lang)}
      subtitle={t('mechanics.subtitle', lang)}
    >
      {/* Mechanism selector + detail view */}
      <div className="grid lg:grid-cols-5 gap-6 mb-16">
        {/* Left: mechanism list */}
        <div className="lg:col-span-2 space-y-1.5">
          {mechanisms.map((m, i) => (
            <motion.button
              key={m.title}
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.04 }}
              onClick={() => setActive(i)}
              className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg border transition-all ${
                active === i
                  ? 'border-current bg-bg-card-hover'
                  : 'border-border/50 bg-bg-card hover:border-border hover:bg-bg-card'
              }`}
              style={{ borderColor: active === i ? m.color + '50' : undefined }}
            >
              <span className="text-lg flex-shrink-0">{m.icon}</span>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <span className="text-xs font-mono font-bold" style={{ color: active === i ? m.color : undefined }}>
                    {m.title}
                  </span>
                  <span className="text-[9px] px-1.5 py-0.5 rounded font-mono flex-shrink-0"
                    style={{ backgroundColor: m.color + '20', color: m.color }}>
                    {m.badge}
                  </span>
                </div>
                <p className="text-[10px] text-text-muted truncate">{m.problemHint}</p>
              </div>
            </motion.button>
          ))}
        </div>

        {/* Right: detail panel */}
        <div className="lg:col-span-3">
          <Annotatable
            id={`mechanism-${active}`}
            label={`Mechanism: ${mechanisms[active].title}`}
            sectionId="mechanics"
            sectionLabel="Claude Mechanics"
          >
          <motion.div
            key={active}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
            className="bg-bg-card border border-border rounded-lg overflow-hidden h-full"
          >
            {/* Header */}
            <div className="px-5 py-4 border-b border-border" style={{ borderLeftColor: mech.color, borderLeftWidth: 3 }}>
              <div className="flex items-center gap-3 mb-1">
                <span className="text-2xl">{mech.icon}</span>
                <div>
                  <h3 className="font-mono font-bold text-text-primary">{mech.title}</h3>
                  <p className="text-xs text-text-muted font-mono">{mech.subtitle}</p>
                </div>
              </div>
            </div>

            {/* THE PROBLEM */}
            <div className="px-5 py-3 border-b border-border bg-neon-red/3">
              <div className="flex items-center gap-2 mb-2">
                <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-neon-red/15 text-neon-red font-bold">
                  THE PROBLEM
                </span>
                <span className="text-[10px] text-text-muted font-mono">what this mechanism was built to solve</span>
              </div>
              <p className="text-xs text-text-secondary leading-relaxed">{mech.problem}</p>
            </div>

            {/* How Drive Uses It */}
            <div className="px-5 py-3 border-b border-border">
              <p className="text-[10px] text-text-muted font-mono mb-2">HOW DRIVE SOLVES IT:</p>
              <p className="text-sm text-text-secondary">{mech.howDriveUses}</p>
            </div>

            {/* Code block */}
            <div className="p-4 border-b border-border/50">
              <p className="text-[10px] text-text-muted font-mono mb-2">{mech.codeLabel}</p>
              <pre className="text-xs font-mono leading-relaxed overflow-x-auto"
                style={{ color: mech.color + 'cc' }}>
                {mech.code}
              </pre>
            </div>

            {/* Insight */}
            <div className="px-5 py-3">
              <div className="flex items-start gap-2">
                <span className="text-[10px] font-mono px-2 py-0.5 rounded flex-shrink-0"
                  style={{ backgroundColor: mech.color + '15', color: mech.color }}>
                  WHY IT MATTERS
                </span>
                <p className="text-xs text-text-muted">{mech.detail}</p>
              </div>
            </div>
          </motion.div>
          </Annotatable>
        </div>
      </div>

      {/* Hooks event table */}
      <h3 className="font-mono text-lg text-text-primary mb-4">{t('mechanics.hooks.title', lang)}</h3>
      <p className="text-sm text-text-muted mb-6">
        {t('mechanics.hooks.desc', lang)}
      </p>
      <div className="overflow-x-auto mb-12">
        <table className="w-full text-xs">
          <thead>
            <tr className="border-b border-border">
              <th className="text-left py-2 px-4 font-mono text-text-muted">Hook Event</th>
              <th className="text-left py-2 px-4 font-mono text-text-muted">Notification Types</th>
              <th className="text-left py-2 px-4 font-mono text-text-muted">Sound File</th>
              <th className="text-left py-2 px-4 font-mono text-text-muted">Body Zero Type</th>
              <th className="text-left py-2 px-4 font-mono text-text-muted">Voice</th>
            </tr>
          </thead>
          <tbody>
            {hookEvents.map((h, i) => (
              <motion.tr
                key={h.event}
                initial={{ opacity: 0 }}
                whileInView={{ opacity: 1 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.1 }}
                className="border-b border-border/30 hover:bg-bg-card/50"
              >
                <td className="py-3 px-4 font-mono font-bold text-neon-cyan">{h.event}</td>
                <td className="py-3 px-4">
                  <div className="flex flex-wrap gap-1">
                    {h.types.map((t) => (
                      <span key={t} className="bg-bg-card px-2 py-0.5 rounded text-text-muted font-mono">{t}</span>
                    ))}
                  </div>
                </td>
                <td className="py-3 px-4 text-text-secondary font-mono">{h.sound}</td>
                <td className="py-3 px-4 text-neon-purple font-mono">{h.bodyZero}</td>
                <td className="py-3 px-4">
                  {h.needs_voice
                    ? <span className="text-neon-green font-mono">afplay ✓</span>
                    : <span className="text-text-muted font-mono">silent</span>}
                </td>
              </motion.tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* How the prompt guides Claude — flow diagram */}
      <h3 className="font-mono text-lg text-text-primary mb-4">{t('mechanics.guide.title', lang)}</h3>
      <p className="text-sm text-text-muted mb-6">
        {t('mechanics.guide.desc', lang)}
      </p>
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
        {[
          {
            step: '1',
            name: 'YAML Frontmatter',
            color: '#f59e0b',
            desc: 'allowed-tools whitelist restricts Claude to only the declared tools. No accidents, no scope creep.',
            snippet: `---
allowed-tools: Bash, Task,
  TeamCreate, SendMessage...
---`,
          },
          {
            step: '2',
            name: 'Identity Reframing',
            color: '#a855f7',
            desc: 'First 3 sentences override the default "assistant" identity. Claude becomes Team Lead before reading anything else.',
            snippet: `You are a Team Lead.
Drive Mode ALWAYS operates
in Agent Team mode.`,
          },
          {
            step: '3',
            name: 'Pseudocode State Machine',
            color: '#3b82f6',
            desc: 'Execution flow defined as pseudocode — not prose. Eliminates ambiguity. Defines states, not suggestions.',
            snippet: `while (mission != COMPLETE) {
  wave_check()
  execute_wave()
  if error → diagnose_fix()
  check_teammate_messages()
}`,
          },
          {
            step: '4',
            name: 'Anti-Pattern List',
            color: '#ef4444',
            desc: '30+ explicitly listed behaviors marked as FAILURE. Seals every escape route back to passive assistant mode.',
            snippet: `❌ Summarizing + stopping
❌ "Should I continue?"
❌ Outputting plain text Q
❌ Working serially
❌ Not using AskUserQuestion`,
          },
          {
            step: '5',
            name: 'Redundant Reinforcement',
            color: '#22c55e',
            desc: '"Never stop" appears 6+ times across different sections. Attention decay in long prompts is compensated by repetition.',
            snippet: `// Opening: "The loop doesn't break"
// Cardinal Rule section: "NEVER STOP"
// Phase 0.4: "do not stop here"
// Loop pseudocode: "continue always"
// Anti-pattern: "❌ stopping"
// Final mindset: "loop is hero"`,
          },
          {
            step: '6',
            name: 'Quantified Rules',
            color: '#06b6d4',
            desc: 'Subjective judgment replaced with math. When to use STL? ≥3 domain tasks. Max teammates? min(width, 5). No guessing.',
            snippet: `effective_parallelism =
  teammates +
  Σ(inner_subagents_per_STL)

direct_teammates =
  min(parallelism_width, 5)

use_STL_if: domain_tasks >= 3`,
          },
        ].map((item, i) => (
          <Annotatable
            key={item.step}
            id={`guide-step-${item.step}`}
            label={`Drive Guide #${item.step}: ${item.name}`}
            sectionId="mechanics"
            sectionLabel="Claude Mechanics"
          >
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.06 }}
              className="bg-bg-card border border-border rounded-lg p-4 hover:bg-bg-card-hover transition-colors"
              style={{ borderTopColor: item.color, borderTopWidth: 2 }}
            >
              <div className="flex items-center gap-2 mb-2">
                <span className="text-xs font-mono font-bold px-2 py-0.5 rounded"
                  style={{ backgroundColor: item.color + '20', color: item.color }}>
                  #{item.step}
                </span>
                <span className="text-sm font-mono font-bold text-text-primary">{item.name}</span>
              </div>
              <p className="text-xs text-text-secondary mb-3">{item.desc}</p>
              <pre className="text-[10px] font-mono text-text-muted bg-bg-dark rounded p-2 overflow-x-auto leading-relaxed">
                {item.snippet}
              </pre>
            </motion.div>
          </Annotatable>
        ))}
      </div>
    </SectionWrapper>
  )
}
