# Harness Operating Model

> Status: Active
>
> Date: 2026-03-29
>
> Purpose: define how `Harness` operates as a leader-led orchestration system, especially the relationship between leader, teammates, councils, issue flow, authority, and governed communication.

---

## Why This Exists

`Harness` is not only a protocol and not only a runtime.

It also encodes a specific operating model:

- one leader/control plane
- multiple bounded teammates
- governed communication
- authority-aware reconciliation
- councils for ambiguous or high-blast-radius decisions

That operating model is one of the most important ideas in the system.
It should not be left implied across code, research, and chat history.

---

## Executive Model

The simplest summary is:

**leader-led, council-corrected, git-grounded execution**

This means:

- the leader owns global mission control
- teammates own bounded local execution
- git/GitHub objects provide the ordered state backbone
- authority rules determine what counts as truth
- councils correct the system when unilateral synthesis would be unsafe

It is not:

- a flat swarm
- a dictator-leader
- a set of independent workers with no common model
- a prompt ritual for role-playing collaboration

---

## Core Roles

### Leader

The leader is the control-plane regulator.

The leader is responsible for:

- issue and mission intake
- mission decomposition
- assignment of work to runs and teammates
- dependency and overlap management
- scheduling judgment
- signal interpretation
- authority-aware synthesis
- deciding when a council is required
- updating the system’s governed truth

The leader is **not**:

- the only one allowed to think
- the only one allowed to inspect evidence
- a solo oracle making unchecked design decrees

The leader is closer to:

- chair
- regulator
- synthesizer
- controller

than to “boss who writes every line.”

### Teammate

A teammate is a bounded execution unit.

A teammate is responsible for:

- executing a specific mission slice
- staying within its assignment boundary
- reporting discoveries, blockers, and changes
- emitting governed signals instead of relying on implicit memory
- escalating when local evidence challenges the current plan

A teammate is **not**:

- a second leader
- a free-roaming reviewer of everything
- a silent script runner

Teammates have real local reasoning power, but not global authority by default.

### Council

A council is a bounded governance mechanism for high-ambiguity or high-risk decisions.

A council is responsible for:

- independent review
- cross-examination
- structured synthesis
- recording accepted and rejected findings
- producing durable action-driving truth

A council is **not**:

- a ceremony for validating a pre-made answer
- a free-form chat room
- a substitute for routine execution

---

## Issue To Mission Control Flow

The operating flow begins with GitHub-native work and becomes locally governable through missions and runs.

### 1. Issue

The issue is the public work item.

It should contain enough information to be executable and auditable.
The issue remains part of the system of record.

### 2. Mission

The mission is the control-plane work object derived from the issue.

The mission exists to:

- anchor orchestration state
- separate conceptual work from run attempts
- hold mission-level authority
- preserve leader-facing continuity across multiple runs

Important rule:

- `latest_run_id` is not automatically the governing truth
- `authoritative_run_id` exists so fresh teammate activity does not silently seize mission control

### 3. Run

A run is one concrete execution attempt.

Runs are where:

- ownership is applied
- desired state is tracked
- signals are emitted
- execution status changes happen
- ship/rc state advances

### 4. Teammate Run

A teammate run is a run with explicit bounded delegation.

It carries:

- scope
- ownership selectors
- authority refs
- reporting expectations
- parent/leader linkage

This keeps delegated work visible to the control plane.

---

## Leader Responsibilities In Practice

The leader should manage the complete work graph, not just the current foreground task.

That includes:

### Issue / Mission Intake

- identify the real problem
- reject weak or ambiguous framing
- detect duplicates or coupling
- identify blast radius
- decide whether the work is standalone or coupled

### Mission Decomposition

- break work into bounded slices
- keep overlap explicit
- define dependencies
- preserve clear ownership

### Assignment

- decide whether work stays with the leader or is delegated
- assign teammate runs with bounded ownership
- keep forbidden sources and evidence expectations explicit where needed

### Monitoring

- read reports and signals
- compare desired state vs actual state
- check whether work has stalled, drifted, or collided

### Reconciliation

- absorb valid local discoveries into mission truth
- reject lower-authority noise
- reschedule when the old plan is wrong
- escalate when a council is needed

### Completion

- ensure completion is behavioral, not textual
- require verification and real state transition
- carry work through ship/rc when relevant

---

## Teammate Contract

Each teammate should operate under a clear contract.

That contract should make these things explicit:

- objective
- scope
- ownership selectors
- desired state
- allowed or preferred sources when needed
- forbidden sources when needed
- reporting expectations
- escalation triggers

The core rule is:

**local execution is free inside the bounded assignment, but global truth remains governed.**

That preserves initiative without collapsing into uncontrolled swarm behavior.

---

## Governed Communication

Teammates are allowed to communicate, but not in an uncontrolled gossip network.

The intended model is:

- teammate -> leader reports
- teammate -> teammate governed signals
- leader-mediated or explicitly permitted cross-run linkage
- councils when disagreement becomes important

This is why `Harness` treats communication as structured state, not just chat.

Important properties:

- communication should be attributable
- communication should be scoped
- communication should be reconcilable
- communication should be supersedable when fresher/higher-authority evidence appears

The goal is not to suppress communication.
The goal is to make communication governable.

---

## Ownership Model

Ownership is one of the main control-plane safeguards.

### Why Ownership Exists

Without explicit ownership:

- teammates collide silently
- global state becomes ambiguous
- reviews happen too late
- overlap is discovered only after damage

### Ownership Selectors

`Harness` currently supports selectors such as:

- `file`
- `dir`
- `module`
- `glob`
- `scope`
- `interface`

The important point is not the names.
The important point is that ownership can be:

- path-based
- semantic
- hierarchical

### Transfer Modes

Ownership transfer is governed, not implied.

Current modes include:

- `move`
- `copy`
- `split`

The key operating insight is:

- parallel work is allowed when governance makes it legible
- not every overlap is a hard block
- but overlap must never be invisible

---

## Desired State And Reconciliation

`Harness` follows a cybernetic control model.

The core pattern is:

- desired state exists
- actual state exists
- disturbances happen
- the leader reconciles them

Desired state may include:

- mission intent
- run intent
- ownership arrangement
- ship/rc target state
- meeting/council expected output

Actual state may include:

- live git state
- run records
- signals
- meeting packets
- verification outputs
- teammate reports

The operating rule is:

**the system should regulate against real current state, not stale summaries.**

That is why authority and freshness matter so much.

---

## Authority Model

Not all evidence is equally trustworthy.

The operating model depends on explicit authority ranking.

The general ordering is:

- live system-of-record state
- active orchestration state
- current design/planning artifacts
- derived or auxiliary artifacts
- legacy artifacts

Authority is also shaped by:

- freshness
- directness
- specificity
- accepted synthesis vs raw opinion

That means:

- a fresh git fact can override stale orchestration memory
- a governed council output can override an unreviewed teammate opinion
- a legacy document cannot silently override live state

This keeps the leader from regulating against the wrong model.

---

## Councils And Meetings

Councils exist because some decisions are too ambiguous or too risky for unilateral synthesis.

They should be used when:

- authority-level evidence conflicts
- multiple teammates produce materially different models
- design quality is in question
- blast radius is high
- release or architecture risk is high

They should not be used for:

- trivial routing
- ordinary local execution choices
- situations where higher-authority evidence already settles the issue

### Council Packet

A real council should carry:

- decision question
- scope
- known context
- locked decisions
- authority hierarchy
- participant roles
- expected output

### Council Phases

The intended flow is:

1. setup
2. independent review
3. cross-examination when conflict remains materially unresolved
4. synthesis
5. output / action

Important clarification:

- cross-examination is the normal path when a credible contested lane remains
- cross-examination is not mandatory ceremony when there is no conflicting lane
- cross-examination is also not required when higher-authority evidence already settles the lane during independent review

In those cases, the leader may move directly into synthesis with an explicit authority-aware record of why further debate is unnecessary.

### `arch-review` Influence

The specific lessons absorbed from [`arch-review`](/Users/liyasong/Project/collov/erwin/skills/arch-review/arch-review/SKILL.md) are:

- lead as synthesizer, not oracle
- explicit authority hierarchy
- bounded role missions
- allowed/forbidden sources
- independent review before debate
- rejected findings recorded with reasons

This is the core reason `Harness` meetings are governance objects rather than loose notes.
It is also why `Harness` treats debate as a bounded governance tool rather than a mandatory ritual.

---

## Debate Model

When councils reach disagreement, the goal is not endless discussion.

The goal is bounded, governed improvement of judgment.

`Harness` therefore treats debate as structured control-plane work:

- claims and dimensions can be grouped
- routes can be planned
- waves can be obligated
- rounds can be capped
- verdicts can be adopted when rules allow

This also means:

- empty debate should be skipped
- aligned lanes may go directly to synthesis
- authority-settled lanes may go directly to synthesis with an explicit verdict record
- only unresolved contested lanes should consume cross-examination rounds

This is already partly implemented in runtime and is one of the clearest places where `Harness` differs from simple workflow tooling.

The long-term target is a fuller `arch-review`-style debate controller, not a free-form argument engine.

---

## Relationship To Git / GitHub

The operating model is intentionally git-grounded.

GitHub and git provide:

- ordered objects
- durable transitions
- auditable evidence
- ownership visibility
- release grouping

That is why `Harness` treats them as the primary state backbone instead of just “one integration among many.”

This also explains why the protocol remains:

`issue -> claim -> drive -> ship -> rc`

Changing the runtime does not change that backbone.

---

## Relationship To Runtime Implementation

The operating model is broader than the current code, but it is not disconnected from code.

Today it lands through:

- Go runtime control plane
- thin shell surfaces
- manifest-driven packaging
- file-backed runtime state
- smoke-testable behavior

This matters because `Harness` is not intended to remain a conceptual diagram.
Its operating model is supposed to be executable.

---

## What Is Already Real

The following parts of the operating model are already concretely present:

- leader-facing `home` and `doctor`
- first-class mission objects
- mission authority separated from freshest run activity
- bounded teammate runs
- governed sibling communication
- ownership selectors and governed transfer
- authority audit across multiple state families
- meeting packets with role/authority discipline
- wave/route/round council orchestration
- governed verdict adoption in some resolved debate lanes

This is enough to call the current system a true prototype line, not just a research stub.

---

## What Is Still Expansion

The operating model still has open future layers:

- fuller scheduler behavior across more than a strong leader control plane
- broader authority adjudication for larger mission bundles without obvious governing dependency
- fuller automated council execution across multiple conflicting lanes
- optional direct mapping between real Codex sub-agents and teammate runtime

These are important, but they sit above a working operating foundation.

---

## Bottom Line

`Harness` is designed to run software work as a governed system:

- issues become missions
- missions become runs
- runs may be delegated to teammates
- teammates communicate under governance
- important disagreements become councils
- councils refine governed truth
- the leader keeps desired state and actual state aligned

That is the operating model the rest of the architecture is trying to serve.
