// ── Skill Builder: Types & Template Generation ──

export interface SkillPhase {
  id: string
  name: string
  description: string
  enabled: boolean
}

export interface SkillConfig {
  meta: {
    name: string
    triggerWord: string
    description: string
  }
  executionModel: 'single' | 'hierarchical' | 'pipeline' | 'custom'
  phases: SkillPhase[]
  rules: {
    neverStop: boolean
    broadcastProtocol: boolean
    directionValidation: boolean
  }
  tools: string[]
}

export const ALL_TOOLS = [
  'Bash', 'Read', 'Write', 'Edit', 'Glob', 'Grep',
  'WebSearch', 'WebFetch', 'Agent', 'TodoWrite',
] as const

export const DEFAULT_PHASES: SkillPhase[] = [
  { id: 'phase0', name: 'Research & Briefing', description: 'Understand the problem before acting', enabled: true },
  { id: 'phase1', name: 'Planning', description: 'Decompose into tasks, identify dependencies', enabled: true },
  { id: 'phaseT', name: 'Team Assembly', description: 'Spawn teammates for parallel execution', enabled: false },
  { id: 'loop', name: 'Execution Loop', description: 'Execute tasks, verify, iterate until done', enabled: true },
]

export function createDefaultConfig(): SkillConfig {
  return {
    meta: { name: '', triggerWord: '', description: '' },
    executionModel: 'single',
    phases: DEFAULT_PHASES.map(p => ({ ...p })),
    rules: { neverStop: true, broadcastProtocol: false, directionValidation: true },
    tools: ['Bash', 'Read', 'Write', 'Edit', 'Glob', 'Grep'],
  }
}

// ── Template Generation ──

function generateFrontmatter(config: SkillConfig): string {
  const tools = config.tools.join(', ')
  return `---
allowed-tools: ${tools}
description: ${config.meta.description || 'A custom Claude Code skill'}
---`
}

function generateHeader(config: SkillConfig): string {
  const trigger = config.meta.triggerWord || 'my-skill'
  return `# ${config.meta.name || 'My Skill'}

> ${config.meta.description || 'A custom skill for Claude Code.'}

**Trigger**: \`/${trigger}\`
**User input**: \`$ARGUMENTS\``
}

function generateExecutionModel(config: SkillConfig): string {
  switch (config.executionModel) {
    case 'single':
      return `## Execution Model: Single Agent

You work as a single agent. Execute tasks sequentially, maintaining full context
throughout the mission. Focus on one task at a time, verify completion, then proceed.

**When to use**: Tasks that require deep context continuity, are small enough for
one agent, or involve tightly coupled changes where parallelism adds no benefit.`

    case 'hierarchical':
      return `## Execution Model: Hierarchical Swarm

You are the **Team Lead** of a hierarchical swarm. You orchestrate, not just execute.

### Team Structure
\`\`\`
YOU (Team Lead) ─── orchestrating + critical path
├── 🎖️ STL-{domain} ─── Sub-Team-Lead (spawns inner subagents)
├── 🔧 implementer ─── leaf worker
└── 🔍 researcher ─── scouts ahead
\`\`\`

### Swarm Rules
- **Maximum Parallelism**: Any work that CAN be parallelized MUST be parallelized.
- **Think in Waves**: Decompose tasks into waves of parallel execution.
  Wave N+1 starts when Wave N's critical path completes.
- **STL Pattern**: When a domain has 3+ tasks, assign a Sub-Team-Lead who
  spawns their own inner subagents.
- **No Solo Work**: If you have 0 active teammates AND unblocked tasks exist,
  you have failed the swarm-first rule.

### Wave Map Template
\`\`\`
Wave 0 (Bootstrap) — Width: N
├── [T1] Task → Agent
└── [T2] Task → Agent

Wave 1 (Core) — Width: N
├── [T3] Task → STL (spawns inner)
└── [T4] Task → Agent

Critical Path: T1 → T3 → ...
\`\`\``

    case 'pipeline':
      return `## Execution Model: Pipeline

Execute as a staged pipeline where each stage's output feeds the next.
Stages run sequentially, but WITHIN each stage, parallelize freely.

### Pipeline Structure
\`\`\`
Stage 1: Research & Understand
  └── Output: findings document
       ↓
Stage 2: Design & Plan
  └── Output: implementation plan
       ↓
Stage 3: Implement
  └── Output: working code
       ↓
Stage 4: Verify & Polish
  └── Output: tested, reviewed deliverable
\`\`\`

**Pipeline Parallelism**: While Stage N executes, begin preparing Stage N+1
(research, scaffolding, test setup) wherever possible.`

    case 'custom':
      return `## Execution Model

Define your execution model here. Consider:
- How many agents will work in parallel?
- What is the coordination strategy?
- How are tasks assigned and tracked?`
  }
}

function generatePhases(config: SkillConfig): string {
  const enabledPhases = config.phases.filter(p => p.enabled)
  if (enabledPhases.length === 0) return ''

  const sections = enabledPhases.map(phase => {
    switch (phase.id) {
      case 'phase0':
        return `## Phase 0: ${phase.name}

${phase.description}

Before writing any code, do your homework:

1. **Read the user's request carefully** — what is the real problem, not just the stated request?
2. **Reconnaissance** — explore the codebase to understand existing patterns, conventions, and constraints.
3. **Surface assumptions** — list what you're assuming. For each: "If this is wrong, how different would my approach be?"
4. **Ask clarifying questions** — use AskUserQuestion for anything where guessing wrong means rework.
5. **Lock the Mission Contract** — define concrete success criteria before proceeding.

\`\`\`
Mission Contract:
□ [Criterion 1 — concrete, verifiable]
□ [Criterion 2 — ...]
□ [Criterion 3 — ...]

Scope: IN: [...] | OUT: [...]
\`\`\``

      case 'phase1':
        return `## Phase 1: ${phase.name}

${phase.description}

Decompose the mission into specific, verifiable tasks:

- Each task has a **clear done-state** linked to a Success Criterion
- Order by dependency — **risky/uncertain tasks first**
- One substantive change per task
- First task: scaffold and verify the build pipeline works
- Last task: final verification against all Success Criteria`

      case 'phaseT':
        return `## Phase T: ${phase.name}

${phase.description}

Create a team and spawn teammates for parallel execution:

1. **Assess parallelism width** — how many tasks can run simultaneously?
2. **Spawn teammates** with clear context (they don't share your conversation)
3. **Assign tasks** — Team Lead takes the critical path, teammates take the rest
4. **Coordinate** — respond to teammate messages immediately, unblock fast`

      case 'loop':
        return `## The Execution Loop

${phase.description}

\`\`\`
while (mission != COMPLETE) {
    execute    →  write the code
    verify     →  run it, test it, read ALL output
    if error   →  read full trace, diagnose root cause, fix, verify again
    if clean   →  self-review before marking complete
    review     →  adversarial self-critique
    if pass    →  mark complete → start next task IMMEDIATELY
}
\`\`\`

**Never stop between tasks.** Don't summarize and wait. Don't ask "should I continue?"
The answer is always: next task.`

      default:
        return `## ${phase.name}

${phase.description}`
    }
  })

  return sections.join('\n\n---\n\n')
}

function generateRules(config: SkillConfig): string {
  const rules: string[] = []

  if (config.rules.neverStop) {
    rules.push(`## Cardinal Rule: The Loop Never Breaks

You are a loop — the loop is the hero, not any single iteration.
You will fail, fix, and iterate until the mission is done.
Every error is data that narrows the problem space.

**The ONLY acceptable reasons to stop:**
1. Waiting for user to answer a question you just asked
2. You are truly blocked and have exhausted all alternatives
3. The mission is COMPLETE and verified

**Everything else — KEEP MOVING.** Either discuss or execute. Never idle.`)
  }

  if (config.rules.broadcastProtocol) {
    rules.push(`## Broadcast Protocol

Every significant action MUST be announced to maintain visibility:

**On starting a task:**
\`\`\`
📡 DISPATCH: [{agent}] → {task-summary}
\`\`\`

**On completing a task:**
\`\`\`
✅ RETURN: [{agent}] ← {task-summary}
   Result: {one-line outcome}
   Files touched: [{file-list}]
\`\`\`

Broadcasts make the execution visible and legible to the user.`)
  }

  if (config.rules.directionValidation) {
    rules.push(`## Direction Validation

At regular checkpoints (every 3 tasks), verify alignment:

1. Re-read the Mission Contract — am I still aligned with Success Criteria?
2. **Direction health check:**
   - "Can I state in one sentence how my current work connects to the mission?"
   - "Am I fighting the code or flowing with it?" (fighting = possible direction error)
   - "Has anything I've learned invalidated an earlier assumption?"
3. If drifted → create corrective tasks before continuing

**Direction error vs. execution error:**
- Execution error: bug in code → fix and continue
- Direction error: building the wrong thing → STOP and re-evaluate`)
  }

  return rules.join('\n\n---\n\n')
}

function generateCompletion(): string {
  return `## Mission Complete

When — and ONLY when — ALL of these are true:
1. Every task is marked complete
2. Final verification passes (build, tests, whatever fits the project)
3. You have checked EACH Success Criterion — verified by running/testing, not assumed
4. Self-adversarial review of the entire deliverable passes
5. You cannot think of a single improvement within scope you haven't made

**Self-Review Checklist:**
- [ ] Re-read every line of code written/modified
- [ ] Ran verification and saw it pass
- [ ] Actively tried to find problems
- [ ] Confident this would survive a hostile code review

Then output: \`MISSION_COMPLETE\` and provide a structured debrief.`
}

export function generateSkillMarkdown(config: SkillConfig): string {
  const sections = [
    generateFrontmatter(config),
    generateHeader(config),
    generateExecutionModel(config),
    generatePhases(config),
    generateRules(config),
    generateCompletion(),
  ].filter(Boolean)

  return sections.join('\n\n---\n\n')
}
