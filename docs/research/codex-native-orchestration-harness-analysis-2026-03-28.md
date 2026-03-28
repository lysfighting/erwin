# Codex-Native Orchestration Harness Analysis

> Date: 2026-03-28
>
> Purpose: define the shape of a new Codex-native, project-agnostic system that synthesizes:
>
> - V4 orchestration ideas
> - harness engineering
> - the protocol depth of `teamwork`
> - the architectural cleanup of `wing`
>
> without inheriting their historical baggage.

---

## Executive Thesis

The next thing should not be:

- `teamwork` migrated again
- `wing` renamed for Codex
- a thin adapter that only maps old commands to new tools

It should be:

**a Codex-native orchestration harness**

This is a higher-order object than a skill.
It combines:

- a portable protocol kernel
- an orchestration control plane
- a Codex runtime adapter
- thin skill surfaces
- project profiles
- eval and maintenance loops

That is the right synthesis of V4 + harness engineering + `teamwork`/`wing` lessons.

---

## Why This Is The Right Reframe

`teamwork` solved a real problem, but in a historically specific shell:

- Claude-command packaging
- long prompt surfaces
- implicit multi-human team framing
- accumulated `vi_agent` shaping

`wing` solved another real problem:

- skill architecture
- modularization
- lower context cost
- cleaner operations split

But `wing` also started drifting into a broader future-facing concept:

- multi-agent direction
- orchestration identity
- architecture beyond the original teamwork domain

Codex migration creates a chance to stop treating all of this as one thing.

The new system should keep:

- the protocol seriousness of `teamwork`
- the modular packaging discipline of `wing`
- the orchestration ambition of V4
- the environment-design mindset of harness engineering

while removing:

- command-era baggage
- project coupling
- runtime confusion
- fake equivalence between "advanced" and "always multi-agent"

---

## What The New Thing Actually Is

The most precise name for the target architecture is:

**a portable autonomous software-work protocol with a Codex-native orchestration harness**

That phrasing matters because it separates two concerns:

1. The protocol
   What software-work lifecycle semantics should remain true across runtimes

2. The harness
   How a concrete runtime like Codex makes those semantics reliable in practice

The new system should therefore be designed as:

- protocol-first
- harness-enabled
- orchestration-capable
- Codex-native

---

## What “Advanced” Should Mean This Time

The word "advanced" should not mean:

- more prompt text
- more ritual
- more custom jargon
- mandatory multi-agent execution

It should mean:

### 1. Better separation of planes

- protocol plane
- control plane
- execution plane
- policy plane
- evaluation plane

### 2. Better use of the runtime

Use Codex for what Codex is already strong at:

- isolated worktrees
- multiple concurrent agent threads
- skills
- automations
- shell-native execution
- optional web research
- optional sub-agents

### 3. Better legibility

The system should not rely on the model "remembering how we work."
It should make the repo, workflow, and state legible.

### 4. Better escalation structure

Single-agent when enough.
Manager/scheduler when needed.
Background automation when appropriate.

That is more advanced than defaulting to "swarm everything."

---

## Official Codex Signals That Matter

The current Codex ecosystem strongly supports this direction.

### Signal 1: Codex is explicitly multi-agent and worktree-native

OpenAI describes the Codex app as a command center for multiple agents, with isolated worktrees and long-running parallel tasks.

Implication:

- your orchestration ideas are no longer speculative
- worktree isolation belongs in the runtime utilization layer, not as a weird add-on

Source:

- [Introducing the Codex app](https://openai.com/index/introducing-the-codex-app/)

### Signal 2: Skills are a first-class extension point

OpenAI explicitly positions skills as reusable bundles of instructions, resources, and scripts.

Implication:

- skills are part of the answer
- but skills are not the entire architecture

Source:

- [Introducing the Codex app](https://openai.com/index/introducing-the-codex-app/)

### Signal 3: Automations are part of the ecosystem

Codex now supports recurring background work through Automations.

Implication:

- background maintenance, triage, hygiene, and garbage collection should be designed in from the start
- not bolted on later

Source:

- [Introducing the Codex app](https://openai.com/index/introducing-the-codex-app/)

### Signal 4: Harness engineering is the actual leverage point

OpenAI's own writeup says the job shifts from writing code to designing environments, specifying intent, and building feedback loops.

Implication:

- your real product is not "the best prompt"
- it is the environment and control structure that make agent work reliable

Source:

- [Harness engineering: leveraging Codex in an agent-first world](https://openai.com/index/harness-engineering/)

### Signal 5: OpenAI recommends incremental orchestration

The practical guide explicitly recommends maximizing a single agent first, then moving to multi-agent only when complexity truly warrants it.

Implication:

- V4 should become an orchestration capability tier
- not the default execution mode for everything

Source:

- [A practical guide to building agents](https://cdn.openai.com/business-guides-and-resources/a-practical-guide-to-building-agents.pdf)

---

## V4 Reinterpreted for Codex

V4 should not be imported literally.
It should be reinterpreted.

The enduring V4 ideas are:

- sibling task awareness
- overlap detection
- dependency graphing
- scheduler / coordinator role
- standardized inter-agent signals
- safe parallelism

Those are still good.

But the execution context is now different.

In Codex, the runtime already gives us:

- multiple threads
- worktrees
- automations
- tool-rich execution

So the new design should not try to reinvent the runtime.
It should add the missing orchestration semantics on top of it.

---

## The New Architecture

The cleanest architecture is six layers.

### Layer 1: Protocol Kernel

Runtime-neutral.

Contains:

- issue -> claim -> drive -> ship -> rc lifecycle
- contract schema
- state machine
- freshness rules
- verification semantics
- blast-radius semantics
- anti-compliance cognitive rules

This is the most durable layer.

### Layer 2: Orchestration Control Plane

This is the real "new thing."

Contains:

- task graph / dependency graph
- overlap analysis
- run scheduler
- coordination policies
- inter-run signaling model
- escalation model
- recovery / requeue rules

This is where V4 matures into a modern harness concept.

### Layer 3: Runtime Adapter

This binds the abstract system to Codex.

Contains:

- mapping abstract operations to Codex primitives
- defining how contracts are persisted
- defining how orchestration state is surfaced
- defining which runtime surfaces are used for which jobs

This is runtime-specific, but should remain thin and intentional.

### Layer 4: Skill Pack

This is not the architecture itself.
It is the user-facing entry surface inside Codex.

Possible skills:

- home / doctor
- issue
- claim
- drive
- ship
- rc
- scheduler / coordinate

These should be thin interfaces into the harness, not giant brains unto themselves.

### Layer 5: Project Profile

Project-agnostic system means project-specific reality must move here.

Contains:

- branch conventions
- verification commands
- deploy / staging workflows
- codebase entry docs
- repo-specific quality rules

Each project should supply a profile.
The harness should not secretly be a `vi_agent` profile in disguise.

### Layer 6: Eval / Maintenance Plane

Contains:

- lint for skills and prompts
- behavior smoke tests
- parity tests
- quality grading
- recurring cleanup automations

This is where harness engineering becomes operational rather than aspirational.

---

## Adapter vs Codex Skill

This is the boundary that must stay clear.

### A Codex skill is:

- one extension unit inside the Codex ecosystem
- a bundle of instructions, scripts, and resources
- something the agent can invoke directly for a task surface

Examples:

- a claim skill
- a release skill
- a review skill

### A runtime adapter is:

- the whole binding layer between an abstract system and a concrete runtime
- the logic that says "how does this protocol/harness live inside Codex?"

It answers questions like:

- which actions are skills vs scripts vs automations?
- how are worktrees used?
- how is state stored?
- how is orchestration surfaced to the user?
- what uses sub-agents, and when?
- what runs in the foreground vs the background?

### Relationship

In Codex:

- a skill is one component of the adapter
- the adapter is larger than the skills

So:

**adapter != skill**

More precisely:

**Codex adapter = skill pack + deterministic mechanisms + state model + automations + project profile conventions**

That distinction is load-bearing.

If you treat "the adapter" as "just a skill," you will repeat the historical overpacking problem.

---

## What Belongs To Codex vs What You Should Build

### Let Codex provide

- agent threads
- worktree isolation
- shell execution
- skill invocation
- automations
- optional delegation
- general-purpose tool access

### Build on top

- protocol kernel
- orchestration policies
- task graph semantics
- contract semantics
- inter-run signaling
- background hygiene rules
- project profile model
- eval / parity system

This separation avoids rebuilding the runtime badly.

---

## Multi-Agent Strategy for the New System

The new system should be orchestration-capable, not orchestration-dependent.

That means three execution tiers.

### Tier 1: Single-agent execution

Default mode.

Use when:

- work is well-bounded
- overlap risk is low
- synthesis cost exceeds delegation benefit

### Tier 2: Manager-guided delegation

Use when:

- there are separable bounded subtasks
- one central thread should retain synthesis/control
- task graph is shallow

This matches the manager pattern in OpenAI's orchestration guidance.

### Tier 3: Scheduled parallel mission graph

Use when:

- multiple tasks have explicit dependencies
- overlap analysis is available
- worktree isolation and review loops are in place
- recovery logic exists

This is the mature form of the V4 vision.

The key decision:

Tier 3 should be earned, not assumed.

---

## The Missing Concept: Control Plane

Historically, `teamwork` had protocol depth.
V4 had orchestration ambition.
Harness engineering adds the missing frame:

**control plane**

The control plane is what decides:

- what should run
- in what order
- under what guardrails
- with what level of parallelism
- when to escalate
- what to persist
- what to garbage-collect

This is probably the deepest missing abstraction in the old design lineage.

Without a control plane, you end up with:

- smart commands
- useful scripts
- parallel ideas

but not a coherent orchestration harness.

---

## Design Principles for the New System

### 1. Protocol over packaging

The lifecycle semantics matter more than the folder layout.

### 2. Control plane over ritual

We want real scheduling and policy, not more scripted ceremony.

### 3. Codex-native over cross-runtime mimicry

Use the actual strengths of Codex instead of emulating Claude-era behavior.

### 4. Project-agnostic by default

Nothing should secretly assume one repo's structure.

### 5. Thin skill surfaces, strong substrate

The more intelligence that can live in reusable mechanisms and explicit state, the better.

### 6. Single-agent first, orchestrate when justified

This is both simpler and more aligned with current best practice.

### 7. Background hygiene is first-class

Recurring cleanup, drift detection, and quality maintenance should be designed in early.

### 8. Legibility beats hidden cleverness

A readable system scales better than an impressive but implicit one.

---

## What To Delete From History

Delete outright:

- command-era packaging assumptions
- project-coupled defaults
- implicit `vi_agent` worldview
- prompt bloat
- menu-heavy branching interaction
- treating "team" as the defining abstraction

Redesign rather than delete:

- contracts
- worktree support
- review flows
- RC flows
- cross-task signaling
- scheduler ideas

Preserve strongly:

- GitHub-first lifecycle
- contract freshness
- deterministic mechanisms
- anti-compliance execution stance
- quality and verification loops

---

## What The New Thing Is Not

It is not:

- only a protocol spec
- only a skill pack
- only a scheduler
- only a Codex customization bundle

It is:

**a control-plane-oriented orchestration harness that uses Codex as its runtime substrate**

That is the clearest current definition.

---

## Strategic Recommendation

Do not define the next system as "Codex Wing."

That is too small and too loaded.

Also do not define it as "Teamwork vNext."

That is too backward-looking.

Define it first as an architectural category:

**Codex-native orchestration harness**

Then later decide whether `wing` should be:

- the product name for that harness
- the orchestration layer above multiple adapters
- or a separate future branch

Naming should follow architecture, not lead it.

---

## Recommended Next Build Order

1. Define the control-plane concepts explicitly:
   run, mission, sibling, overlap, dependency, scheduler, signal, escalation.
2. Define the Codex runtime adapter surface:
   which concerns map to skills, automations, scripts, state files, and optional sub-agents.
3. Define the project profile schema:
   prove project-agnosticism structurally.
4. Build a thin Codex `home/doctor` skill.
5. Build a thin Codex `claim` skill.
6. Build the first `scheduler/coordinate` surface only after single-agent claim-drive behavior is solid.
7. Add recurring garbage-collection / hygiene automations early.

This order keeps the ambitious part real without repeating old complexity mistakes.

---

## Bottom Line

The thing you should incubate now is not a cleaned-up `teamwork`, and not a direct `wing` continuation.

It is a new category:

**a Codex-native orchestration harness**

Its soul comes from:

- `teamwork`'s lifecycle protocol
- `wing`'s modular packaging discipline
- V4's coordination ambition
- harness engineering's environment-first worldview

But its architecture should be new.
