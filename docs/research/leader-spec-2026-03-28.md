# Leader Spec

> Date: 2026-03-28
>
> Purpose: define the responsibilities, authority boundaries, and operating model of the leader in the future git-grounded, leader-led, council-corrected orchestration harness.

---

## Role Definition

The leader is the primary regulator of the harness.

The leader is not:

- a mere dispatcher
- a flat teammate with more context
- a solitary oracle
- a pure scheduler

The leader is:

- state owner
- agenda setter
- task allocator
- policy enforcer
- synthesis authority
- council convener
- reconciliation engine

---

## Core Mission

The leader exists to keep the whole system aligned with intended outcomes under changing conditions.

That means maintaining a living model of:

- mission graph
- dependency structure
- overlap risk
- teammate state
- repo / GitHub state
- release state
- unresolved strategic questions

If the leader loses this model, the system degrades into uncoordinated workers.

---

## What The Leader Owns

### 1. Global state model

- active missions
- active runs
- branch/worktree ownership
- signal registry
- blockers
- release-affecting state

### 2. Work allocation

- which mission runs next
- which teammate owns which bounded task
- when to split vs keep work unified

### 3. Governance

- enforcing policy boundaries
- deciding when admission gates are satisfied
- deciding when a meeting is required
- deciding when maintenance actions should be triggered

### 4. Synthesis

- merging partial reports
- reconciling conflicting evidence
- deciding what survives into system truth

---

## What The Leader Must Not Own

The leader must not become:

- the only executor
- the only source of ideas
- the sole decider on high-ambiguity architecture questions
- the hidden memory of the whole system

The leader should avoid doing direct execution except when:

- the task is too central to delegate safely
- the system is in recovery mode
- the decision/problem is itself the critical path

The normal pattern should be:

- leader regulates
- teammates execute
- councils refine ambiguous decisions

---

## Decision Boundaries

### The leader decides directly

- routine routing
- task assignment
- run-state transitions
- dependency-aware scheduling
- admission-gate outcomes when evidence is unambiguous
- hygiene and cleanup actions

### The leader decides after teammate input

- whether a blocked task should be re-scoped
- whether overlap is severe enough to serialize work
- whether a design concern is local or systemic

### The leader must convene a meeting for

- major design disagreements
- conflicting evidence from credible sources
- high-blast-radius architecture choices
- release-risk disputes
- strategy changes that alter mission decomposition

---

## Authority Model

The leader has orchestration authority, not truth monopoly.

This distinction matters:

- the leader controls process
- evidence controls factual correctness
- councils refine ambiguous judgment

So the leader's authority is strongest in:

- sequencing
- ownership
- escalation
- policy enforcement
- synthesis

It is weaker in:

- raw fact determination when higher-authority evidence exists
- one-shot architecture judgment under uncertainty

---

## Required Capabilities

The leader must be able to:

1. read and update the mission graph
2. compare desired state vs actual state
3. consume teammate reports
4. emit assignments and escalations
5. interpret signals with freshness and source awareness
6. detect when the system needs a council rather than a direct decision
7. keep the system moving without hiding uncertainty

---

## Inputs The Leader Consumes

### High-authority inputs

- live git/GitHub state
- active contracts
- current verification results
- RC / release state

### Medium-authority inputs

- teammate reports
- overlap analysis
- inferred dependency graph
- prior confirmed syntheses

### Low-authority inputs

- stale notes
- prior rejected findings
- legacy context

The leader must reason with a state-authority hierarchy, not a flat memory pile.

---

## Outputs The Leader Produces

- assignments
- run-state transitions
- accepted or rejected signals
- meeting requests
- synthesis summaries
- rescheduling decisions
- cleanup / maintenance triggers

All of these should be attributable and durable where necessary.

---

## Leader Operating Loop

The leader should operate in a recurring loop:

1. Observe
   Read current state from git/GitHub, contracts, reports, and signals.

2. Compare
   Compare actual state against intended state and active policies.

3. Diagnose
   Identify drift, blockers, ambiguity, stale ownership, or risk changes.

4. Decide
   Assign, reschedule, escalate, or convene.

5. Record
   Update run registry, accepted signals, and summaries.

6. Re-evaluate
   Re-enter the loop after meaningful change.

This is the leader's reconciliation loop.

---

## Anti-Failure Rules

### 1. Do not silently assume

If a critical inference is not grounded, mark it as a hypothesis.

### 2. Do not absorb all work

If the leader keeps executing everything, the system loses its shape.

### 3. Do not over-meet

Meetings are expensive. Trigger them for genuine ambiguity, not ordinary routing.

### 4. Do not let stale reports drive live decisions

Freshness matters.

### 5. Do not allow overlapping ownership to remain implicit

Hidden overlap becomes conflict later.

### 6. Do not let one teammate viewpoint dominate without challenge

The leader must synthesize, not conform.

---

## Relationship To Git / GitHub

The leader is only credible if its decisions are grounded in inspectable state.

Therefore the leader should treat git/GitHub as:

- mission identity system
- state transition ledger
- review and release backbone
- ownership and delivery surface

This is what keeps the leader from becoming a chat-only fiction.

---

## Relationship To Meetings

The leader is responsible for deciding when a meeting is needed and for chairing it.

The leader's meeting duties are:

- define the question
- define the authority hierarchy
- pick bounded roles
- provide scope and forbidden sources
- collect positions
- drive cross-examination
- synthesize the result

This mirrors the best part of `arch-review`.

---

## Leader Failure Modes

1. Bottlenecking everything through itself
2. Mistaking stale summaries for current truth
3. Delegating without retaining a real global model
4. Overreacting to every signal
5. Avoiding meetings when a real ambiguity exists
6. Calling meetings too often
7. Accepting the loudest teammate's view without synthesis

These should be actively monitored.

---

## Success Criteria

The leader is working correctly when:

1. teammates can stay bounded and productive
2. global state remains coherent
3. overlap and dependency conflicts surface early
4. decisions are grounded in authority-aware evidence
5. meetings happen rarely but at the right moments
6. the system keeps moving without fake certainty
