# Codex Surface Mapping Spec

> Date: 2026-03-28
>
> Purpose: define how the new orchestration harness should map its responsibilities onto actual Codex surfaces so we do not repeat the mistake of overloading `SKILL.md` with architecture, state, policy, and mechanism all at once.

---

## Why This Exists

The most dangerous architectural mistake right now would be:

**treating Codex skills as the whole system**

Codex gives multiple usable surfaces:

- skills
- shell/scripts
- automations
- optional web research
- optional sub-agents
- local files
- repo instructions

This document defines what each surface should and should not own.

---

## Surface Inventory

### 1. Skills

User-facing entry surfaces for bounded workflows.

Good for:

- intent capture
- routing
- lightweight diagnosis
- summarization
- orchestration decisions
- invoking deterministic mechanisms

Bad for:

- being the sole state store
- holding giant operational truth tables
- embedding project-specific policy
- acting as the only place where architecture lives

### 2. Scripts / Deterministic Mechanisms

Machine-executable substrate.

Good for:

- config resolution
- contract lifecycle
- git / worktree operations
- label state machine
- PR helpers
- RC helpers
- overlap analysis
- scheduler mechanics that are deterministic enough

Bad for:

- subjective judgment
- broad natural-language synthesis
- human-facing explanation as primary responsibility

### 3. Automations

Recurring or background work surfaces.

Good for:

- stale-state cleanup
- periodic doctor checks
- drift detection
- queued report generation
- release hygiene
- regular parity/eval runs

Bad for:

- one-off foreground workflows
- deeply interactive decision trees
- tasks whose value depends on immediate human steering

### 4. State Files

Durable local state that should survive session boundaries.

Good for:

- run registry
- signal log
- scheduler outputs
- mission state snapshots
- maintenance backlog

Bad for:

- storing trivial shell output
- replacing GitHub as source of truth
- hiding important state in obscure locations

### 5. Sub-Agents

Optional parallel execution helpers.

Good for:

- bounded side investigations
- disjoint implementation slices
- parallel review/verification tasks

Bad for:

- being required for the protocol to work
- replacing the control plane
- compensating for unclear architecture

### 6. Web Research

Freshness and source-validation surface.

Good for:

- current ecosystem facts
- documentation verification
- release notes / product behavior checks

Bad for:

- replacing repo-local understanding
- driving execution when local evidence is stronger

---

## Responsibility Map

| Responsibility | Primary Codex surface | Notes |
|---|---|---|
| protocol entrypoints | Skills | `home`, `issue`, `claim`, `drive`, `ship`, `rc`, `coordinate` |
| contract generation/update | Skill + script + state file | Skill decides, script normalizes, file persists |
| git/worktree mechanics | Scripts | Keep deterministic and testable |
| run registry | State files | Should survive chat/session boundaries |
| scheduler decisions | Skill + state file | Skill computes/presents, state persists |
| overlap analysis | Script | Better as mechanized analysis than prose |
| recurring cleanup | Automations | First-class harness hygiene |
| freshness checks | Script + skill | Script computes, skill interprets |
| background reporting | Automations | Useful for manager/coordinator views |
| optional parallel investigation | Sub-agents | Only when justified and allowed |
| current external grounding | Web research | Use official sources when freshness matters |

---

## Recommended Skill Pack

The new system should expose a thin skill pack.

### `home`

Responsibilities:

- show current harness state
- diagnose obvious mismatches
- show next valid step
- route to specialized skills

### `issue`

Responsibilities:

- create / refine executable issues
- batch issue generation
- quality lint for issue completeness

### `claim`

Responsibilities:

- choose or accept mission
- validate executability
- create contract
- provision branch/worktree
- register run

### `drive`

Responsibilities:

- load current mission/run
- evaluate plan against code
- execute or coordinate execution
- emit state and signals

### `ship`

Responsibilities:

- run pre-flight
- create/update PR
- sync labels/state
- handle cleanup transitions

### `rc`

Responsibilities:

- prepare RC
- promote
- apply fix
- run hotfix/backport flows

### `coordinate`

Responsibilities:

- show mission graph
- show blockers and overlaps
- select next runnable batch
- surface signals and escalations

This skill should appear only after single-run behavior is mature enough.

---

## What Skills Must Not Own

Skills must not become:

- the sole scheduler database
- the only representation of active runs
- the only repository of project policy
- the only way to understand system state

If that happens, the architecture collapses back into prompt-centric design.

---

## Recommended File-State Layout

Exact paths can change, but the categories should exist.

### Durable local state

- run registry
- signal log
- scheduler snapshots
- maintenance backlog

### Durable mission state

- active contracts
- mission notes
- freshness metadata

### Project profile state

- profile config
- routing defaults
- verification conventions

The adapter should define the path conventions.
The architecture should define the categories.

---

## Skill vs Script Boundary

Use a skill when the task needs:

- interpretation
- orchestration
- summarization
- decision framing

Use a script when the task needs:

- deterministic transformation
- repeatable repo operations
- state normalization
- machine-checkable output

Rule of thumb:

If two good agents should produce the same output every time, prefer a script.
If good agents may legitimately choose different but valid actions, keep that in the skill/control-plane layer.

---

## Skill vs Automation Boundary

Use a skill when the task is:

- foreground
- interactive
- intent-driven
- execution-sensitive

Use an automation when the task is:

- recurring
- hygiene-oriented
- report-oriented
- cleanup-oriented

This is especially important for harness engineering because maintenance work is continuous.

---

## Skill vs Sub-Agent Boundary

Use the primary skill thread when:

- the task is on the critical path
- synthesis is tightly coupled
- the question is architectural

Use a sub-agent when:

- the task is well-scoped
- the write surface is disjoint
- waiting would otherwise waste time

The control plane may recommend delegation.
The runtime rules still govern whether it happens.

---

## Codex-Native Mapping for the New Harness

### Layer 1: Protocol kernel

Not a Codex surface.
Represented in documents/specs.

### Layer 2: Control plane

Primarily:

- skill logic
- state files
- scripts for deterministic analysis

### Layer 3: Runtime adapter

Primarily:

- skill pack
- scripts
- file layout conventions
- automation definitions

### Layer 4: Project profile

Primarily:

- config files
- repo docs
- profile metadata

### Layer 5: Eval / maintenance

Primarily:

- scripts
- automations
- reports

---

## Early Automation Candidates

To fit the Codex ecosystem better, the new system should design automations early.

Good first candidates:

1. stale contract cleanup report
2. worktree drift / mismatch report
3. queued mission freshness check
4. RC readiness report
5. parity/eval regression report

These are exactly the kinds of tasks that should stop living only in foreground sessions.

---

## Anti-Patterns

Do not:

- put scheduler truth only in a skill prompt
- rely on chat memory as active mission registry
- encode one repo's conventions into the generic harness
- use automations as a substitute for clear state models
- use sub-agents to paper over weak control-plane design

---

## Immediate Implementation Consequences

Based on this mapping, the next implementation work should be:

1. define the local state schema
2. decide path conventions for state and profiles
3. identify which existing `teamwork`/`wing` scripts can be promoted into generic deterministic substrate
4. design the thin `home` and `claim` skills against that substrate
5. design at least one maintenance automation early instead of treating it as optional later

This is how the new harness becomes truly Codex-native rather than merely Codex-compatible.
