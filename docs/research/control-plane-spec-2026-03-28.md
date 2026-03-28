# Control Plane Spec

> Date: 2026-03-28
>
> Purpose: define the control plane for the new Codex-native orchestration harness so V4 orchestration ideas can mature into a stable, project-agnostic architecture rather than remain scattered across prompts, scripts, and intuition.

---

## Why This Exists

The protocol core defines the lifecycle.
The runtime adapter defines how the lifecycle lands in Codex.

What was still missing is the layer that decides:

- what should run
- in what order
- with what isolation
- under what risk rules
- with what degree of parallelism
- how runs inform each other
- when to escalate
- what to persist
- what to clean up

That missing layer is the **control plane**.

Without it, we have:

- useful lifecycle semantics
- useful deterministic mechanisms
- useful skill surfaces

but not a coherent orchestration harness.

---

## Scope

This control plane is:

- project-agnostic
- runtime-aware
- Codex-first in current implementation
- compatible with single-agent and multi-agent execution

It is not:

- a replacement for the protocol core
- a replacement for Codex runtime primitives
- a mandatory always-on scheduler
- a team-collaboration ontology

---

## Core Role

The control plane exists to convert a set of possible work items into a governed execution system.

It answers:

1. What can run now?
2. What should wait?
3. What must run in isolation?
4. What can run in parallel safely?
5. What signals should propagate?
6. When should a human be asked?
7. What background maintenance should continue even when no active run is foregrounded?

---

## Canonical Objects

### 1. Mission

A mission is one execution-ready work unit derived from a GitHub Issue and represented locally by a Contract.

Required fields:

- issue id
- title
- objective
- success criteria
- branch/worktree target
- verification commands
- state
- freshness metadata
- risk / blast-radius marker

### 2. Run

A run is one concrete execution attempt of a mission within a runtime context.

Why it matters:

- a mission is conceptual
- a run is operational

Required fields:

- run id
- mission id
- runtime
- mode
- owner
- started at
- current state
- execution surface
- active worktree/session reference
- latest summary

### 3. Sibling

A sibling is another mission/run considered relevant to the current one because of dependency, overlap, shared milestone, or shared release surface.

Relationship types:

- depends-on
- blocks
- overlaps
- adjacent
- release-coupled

### 4. Dependency

A dependency is a directed relationship where one mission/run should not advance past a gate before another reaches a required state.

Kinds:

- explicit dependency
- inferred file/system dependency
- release dependency
- policy dependency

### 5. Overlap

Overlap is a signal that two missions/runs may contend or interfere.

Kinds:

- file overlap
- module overlap
- interface overlap
- test overlap
- deploy/release overlap

Overlap is not automatically a block, but it changes scheduling and review rules.

### 6. Signal

A signal is a structured piece of state or evidence emitted from one run for the control plane or sibling runs to consume.

Kinds:

- shipped
- interface changed
- blocked
- risk raised
- dependency satisfied
- rebase required
- release-ready
- review requested

### 7. Escalation

An escalation is a control-plane decision that a run or mission requires higher-level attention.

Targets:

- human
- coordinator skill
- maintenance automation
- manager thread

### 8. Policy

A policy is a reusable rule that constrains scheduling or state transition behavior.

Examples:

- P0 missions never run in parallel
- overlap above threshold requires sequential execution
- RC-affecting work requires stronger human checkpoints
- production hotfix runs have stricter cleanup

### 9. Garbage

Garbage is stale or residual orchestration state that no longer reflects valid active work.

Examples:

- stale contracts
- stale worktrees
- orphaned run state
- unresolved stale sibling metadata
- obsolete overlap caches

This is a first-class object category because automation should manage it.

---

## State Models

### Mission State

Canonical states:

- `queued`
- `claimed`
- `wip`
- `review`
- `done`
- `blocked`
- `canceled`

Mission state is a durable lifecycle view.

### Run State

Canonical states:

- `draft`
- `ready`
- `executing`
- `paused`
- `waiting`
- `blocked`
- `completed`
- `abandoned`

Run state is a runtime execution view.

Important distinction:

- one mission may have multiple historical runs
- only one run should normally be the active foreground run for a mission

### Scheduler State

For a set of missions/runs, the scheduler can classify each as:

- `runnable`
- `deferred`
- `blocked-by-dependency`
- `blocked-by-policy`
- `blocked-by-overlap`
- `completed`
- `needs-human-decision`

---

## Execution Modes

The control plane must support multiple modes without changing protocol semantics.

### Mode A: Solo

One foreground run.
No sibling coordination required.

Use when:

- work is bounded
- overlap risk is low
- orchestration cost exceeds benefit

### Mode B: Managed Delegation

One primary run with bounded helper runs or agents.

Use when:

- there are separable subproblems
- synthesis should stay centralized

### Mode C: Scheduled Parallel Missions

Multiple missions progress under a scheduler using dependency and overlap policies.

Use when:

- mission graph is explicit enough
- isolation exists
- signaling exists
- recovery logic exists

This is the mature V4 mode.

---

## Scheduler Responsibilities

The scheduler is not "another agent personality."
It is a control-plane function.

Responsibilities:

- collect candidate missions
- derive or ingest dependencies
- derive or ingest overlap signals
- classify mission risk
- choose execution mode
- produce an execution plan
- emit scheduling decisions
- detect need for rescheduling

Minimal scheduler output should include:

- selected batch or next mission
- blocked missions and why
- required human decisions
- required follow-up signals

---

## Scheduling Heuristics

These heuristics are defaults, not hardcoded truth.

### 1. Dependency before priority

Priority matters only among runnable missions.
Dependencies gate runnability first.

### 2. Risk before throughput

Higher blast radius lowers acceptable parallelism.

### 3. Isolation before convenience

If two missions may interfere, isolation or sequencing wins over speed.

### 4. Release adjacency is special

RC- or production-adjacent work gets stricter policies even if code overlap looks low.

### 5. Signal freshness matters

Old sibling signals should not drive new decisions without freshness checks.

### 6. Default to less orchestration

If the control plane lacks enough information, it should reduce parallelism rather than guess.

---

## Signal Model

Signals are the bridge between runs without requiring hidden shared memory.

### Required properties

- source attribution
- timestamp
- mission/run identity
- confidence or certainty class
- affected scope
- freshness window

### Recommended signal types

- `completion`
- `interface_change`
- `blocking_issue`
- `test_regression`
- `review_feedback`
- `priority_change`
- `dependency_cleared`
- `requires_rebase`

### Signal rules

1. No transitive rumor propagation.
2. Signals must remain attributable.
3. Control-plane use of signals must respect freshness windows.
4. High-impact signals should trigger summary surfacing, not silent state mutation.

This keeps the old anti-plague wisdom while modernizing the model.

---

## Escalation Model

Escalation should be explicit, not accidental.

### Escalate to human when

- destructive action is about to occur
- ambiguity changes business or release risk
- dependency interpretation is uncertain and high-impact
- policy conflict cannot be resolved mechanically

### Escalate to coordinator when

- mission graph needs recomputation
- overlap/risk profile changed materially
- a blocked run may unblock siblings

### Escalate to maintenance automation when

- stale state should be cleaned
- background validation is due
- queued diagnostics should run

---

## Persistence Model

The control plane must decide what deserves durable persistence.

### Persist durably

- mission contract
- run registry / summary
- structured signals
- scheduler decisions that affect future work
- stale-state markers for maintenance

### Keep ephemeral

- transient reasoning
- shell scratch output
- temporary exploration branches of thought

### Principle

Persist anything another agent, another session, or an automation would need to continue safely.

---

## Control Plane vs Runtime Adapter

The control plane defines **what must be decided**.
The runtime adapter defines **how those decisions are executed in Codex**.

Examples:

- control plane says "mission B is blocked by mission A"
- adapter decides how that appears in a skill, state file, or automation

- control plane says "emit interface-change signal"
- adapter decides whether that becomes issue comment, local state, or run summary

The control plane must never collapse into Codex-specific UI details.

---

## Control Plane vs Skill Pack

The skill pack is one user-facing surface into the control plane.

Skills may expose:

- home / doctor
- claim
- drive
- ship
- rc
- scheduler / coordinate

But the control plane should remain a deeper model than any single skill file.

This protects the architecture from becoming prompt-shaped.

---

## Required Invariants

1. The control plane must remain project-agnostic.
2. The control plane must not require multi-agent execution.
3. The control plane must work with incomplete information by reducing aggressiveness.
4. The control plane must keep attribution and freshness explicit.
5. The control plane must treat cleanup and stale-state management as first-class.
6. The control plane must integrate with protocol semantics rather than override them.

---

## Phase-Based Maturity Model

### Phase 1: Single-run governance

Capabilities:

- mission/run distinction
- basic run state
- explicit next-step output
- stale-state cleanup

### Phase 2: Sibling awareness

Capabilities:

- sibling registry
- overlap metadata
- signal emission
- light dependency tracking

### Phase 3: Scheduler

Capabilities:

- runnable classification
- batch planning
- simple policy engine
- rescheduling triggers

### Phase 4: Mature orchestration

Capabilities:

- dynamic rescheduling
- richer signal taxonomy
- automation-driven maintenance
- release-aware coordination

This keeps ambition staged instead of overpromised.

---

## Immediate Design Consequences

If we take this spec seriously, the next implementation work should be:

1. define a minimal run registry format
2. define a signal schema
3. define scheduler input/output schema
4. decide where Codex persists each of those
5. expose the simplest useful control-plane surface before building full parallel orchestration

That is the correct next step after this document.
