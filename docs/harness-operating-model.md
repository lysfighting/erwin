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

It is also not:

- a system where every serious decision must impersonate one `arch-review` debate

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

Councils and other governed decisions should now be understood through explicit protocols rather than one default ritual.

---

## Decision Protocols

`Harness` currently recognizes three decision protocols:

### `evidence-review`

Use when:

- technical or factual disagreement is credible
- authority and evidence quality determine truth
- challenge is needed before synthesis

This is the protocol most directly influenced by `arch-review`.

### `incident-command`

Use when:

- the decision is release-adjacent or operationally urgent
- the leader needs clear action ownership and stop/go boundaries
- extended debate would delay bounded response

### `decision-record`

Use when:

- the decision is important but not urgent
- the main need is durable recording of context, alternatives, choice, and consequences

Current runtime note:

- the meeting engine still shares one common phased execution baseline
- protocol choice now affects default role packets, expected output, scheduler guidance, and whether same-authority conflict automatically becomes a live debate obligation
- `evidence-review` retains the richest route/round automation
- `incident-command` and `decision-record` now prefer bounded synthesis/record handling over auto-generated cross-examination when no authority winner already exists
- `incident-command` now also drives protocol-specific live requirements such as bounded incident assessments in independent review and explicit action/owner/stop-go output guidance at close
- `decision-record` now also drives protocol-specific live requirements such as bounded decision-record assessments in independent review and explicit decision/alternatives/consequences/follow-up output guidance at close
- once those councils reach synthesis/output, orchestration now also auto-synthesizes protocol-aware output packets so the leader can review a structured provisional output artifact before adoption/close instead of reconstructing it from raw findings and actions
- deeper protocol-specific automation is still not complete, but protocol choice is now part of real council control flow rather than metadata only

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
If that synthesis record is already ready for output, the same orchestration pass may close the council directly into `output` instead of requiring a second ceremonial pass just to advance the phase marker.
When more than one debate batch is open, the controller should expose one deterministic current batch and keep the others explicitly blocked behind it instead of pretending every open lane is equally active at once.

When a new same-phase lane appears while the current batch is still open, the controller should normally retain that current batch instead of thrashing to the newly inserted lane. `Harness` now makes this explicit through `batch_selection_policy` and `previous_current_wave_batch_id`.

### `arch-review` Influence

The specific lessons absorbed from [`arch-review`](/Users/liyasong/Project/collov/erwin/skills/arch-review/arch-review/SKILL.md) are:

- lead as synthesizer, not oracle
- explicit authority hierarchy
- bounded role missions
- allowed/forbidden sources
- independent review before debate
- rejected findings recorded with reasons
- governance policy that blocks council advancement when those bounded-review requirements are missing

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
- the same bounded participant can receive one multi-lane debate wave batch instead of several invisible scattered obligations
- rounds can be capped
- verdicts can be adopted when rules allow

This also means:

- empty debate should be skipped
- aligned lanes may go directly to synthesis
- authority-settled lanes may go directly to synthesis with an explicit verdict record
- only unresolved contested lanes should consume cross-examination rounds

This is already partly implemented in runtime and is one of the clearest places where `Harness` differs from simple workflow tooling.

The long-term target is a fuller `arch-review`-style debate controller, not a free-form argument engine.

But that does **not** mean every governed decision in `Harness` should use the same debate shape.
The control plane should choose among `evidence-review`, `incident-command`, and `decision-record` according to decision type.

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
- leader-facing scheduler baseline with durable snapshots, mission classification, execution-mode selection, bundle-aware deferral/escalation, bounded parallel admission from profile policy, profile-backed high-risk scope inference, and rescheduling hints
- scheduler now enforces explicit parallel-justification gating when the execution policy requires it, reducing parallelism rather than guessing
- scheduler now enforces the control-plane rule that runnable `P0` missions do not widen into parallel foreground work
- fresh follow-up signals can now trigger scheduler-driven reclassification and rescheduling instead of remaining display-only telemetry, including normalized `meeting_required` / `requires_meeting`, `review_feedback`, `priority_change`, `risk_raised`, and `requires_rebase` control-plane signals
- those follow-up packets now also project `fresh_follow_up_signal_kinds` plus `follow_up_signal_urgency_rank`, and a fresh `risk_raised` signal raises `effective_risk_level` to `high` before adoption so urgent follow-up decisions can outrank lower-urgency peers inside the same scheduler portfolio
- when more than one fresh follow-up signal is pending on the same mission, the operating model now expects the strongest pending signal to control urgency and leader guidance rather than whichever lower-urgency signal was seen first
- scheduler can now distinguish generic leader decisions from council-worthy decisions by consuming profile `meeting_triggers`, surfacing `recommended_decision_surface`, `meeting_trigger_refs`, and `council_decision_count`, and collapsing escalation-required mission bundles into one bundle-scoped decision packet instead of one duplicated decision per mission
- scheduler now orders human decisions into a portfolio queue, exposing a `primary_required_human_decision`, deferred-decision count, and scope/surface tallies so the leader can handle bundle-scoped council decisions before lower-scope follow-up decisions when both coexist
- scheduler now projects that queue back into mission classifications through `human_decision_queue_position` and `primary_human_decision`, so leader-facing mission views stay aligned with the control-plane decision portfolio
- scheduler now emits explicit `defer_pending_primary_human_decision` actions for lower-priority governed decisions, turning the decision portfolio into real rescheduling behavior instead of only display metadata
- scheduler human-decision packets now carry stable `decision_id` values, and mission classifications plus defer actions now project those ids back into the runtime so the decision portfolio can become a durable bridge into later council materialization
- leader surfaces can now materialize either the current primary human-decision packet or a specific `decision_id` into a governed meeting, and the resulting council records its source decision packet plus scheduler snapshot refs directly in durable meeting state
- materialization now reuses an already-active meeting for the same `decision_id`, so one governed scheduler decision does not silently fan out into multiple active councils
- that stable `decision_id` also remains re-addressable after the decision has already promoted into `active-council`, because the active meeting now acts as a durable anchor for the source packet rather than forcing future materialization to depend on whichever packet shape happens to be present in the latest scheduler snapshot
- when that materialized decision is a bundle-scoped `establish_bundle_governor`, the meeting now also auto-seeds the host plus candidate-governor participants, sibling mission refs, and candidate tuple summaries into meeting state so the leader can begin from bounded scheduler evidence rather than rebuild the obvious council packet by hand
- once that materialization happens, scheduler now recognizes the same stable decision as an `active-council` packet and projects `active_meeting_ids` back into the portfolio instead of treating the open council as detached meeting noise
- once a bundle-scoped decision is already materialized, scheduler now projects that council back into every affected mission through `bundle_active_decision_meeting_ids`, `bundle_active_decision_meeting_count`, and `bundle_control_state`, counts `active_council_bundle_count` at the portfolio level, and emits `hold_for_active_bundle_council` when execution should stay behind a live bundle council
- that active bundle council also remains represented as a bundle-scoped decision in the portfolio, so the leader keeps seeing the governed bundle as one durable in-flight decision instead of watching it collapse into mission-level review noise after materialization
- once a council is in flight, scheduler now mirrors its live progress back into the decision packet itself, including current phase and next requirement, so the portfolio becomes a true control surface for governed councils rather than just a launch point
- that live council packet now also carries the current work packet itself, deriving it from either pending independent-review wave obligations or the current debate batch, so leader-facing surfaces can see the concrete in-flight ask rather than only a generic requirement string
- leader-facing council surfaces now also carry the concrete current-wave action packets and blocked batch summaries, making the active bounded ask and queued follow-up batches visible without reopening raw meeting state
- that live council requirement now also drives the leader-facing `next_action`, `home` summary, and `hold_for_active_bundle_council` reason, so in-flight governance speaks in current bounded work rather than generic review prompts
- active councils now surface freshness directly in the same portfolio packet, and once an in-flight council goes stale the scheduler promotes “refresh or close the stale active council” ahead of ordinary bundle hold messaging
- explicit bundle-governor signals are no longer trusted blindly when a live bundle council exists; if the explicit governor and the active bundle-council host disagree, scheduler now escalates that mismatch into governance maintenance rather than carrying two incompatible authority anchors
- scheduler can now refuse to admit brand-new low-urgency foreground work while a higher-priority human decision remains unresolved, exposing `hold_for_primary_human_decision` when no established foreground batch should be widened past the current decision backlog
- escalation-required mission bundles now surface ranked `candidate_governor_ids` and `candidate_governor_packets`, so the leader can see plausible governors and their current authority tuple before deciding whether to establish a governor directly or materialize a council
- leaders can now explicitly adopt and retire bundle governors through mission surfaces, letting an escalation-required bundle become governed authority deliberately instead of waiting only for implicit adjudication or a future council outcome
- bundle-scoped `establish_bundle_governor` councils can now close with a structured locked decision naming the governing mission, and `meeting-close` adopts that verdict into durable bundle authority automatically
- when a bundle-scoped governance council already has a clear accepted winner, auto-synthesis can now generate that structured governing lock directly, reducing the normal path from `materialize -> manual lock -> close` to `materialize -> orchestrate -> close`
- governed bundle packets now preserve whether that authority came from a closed council adoption, exposing both `governance_source_kind` and the originating `adopted_from_meeting_id`
- stale explicit bundle-governor signals now become explicit maintenance work, and `doctor` reports them as invalid control-plane governance instead of letting them drift silently behind the active mission graph
- scheduler now retains a sticky primary human-decision packet when a fresher same-class peer arrives, exposing `decision_selection_policy` and `previous_primary_human_decision_key` so the leader can distinguish deliberate decision stability from raw freshness sorting
- active-council packets now also carry `active_decision_progress_rank`, and scheduler uses that progress signal to prioritize the more advanced live council when multiple same-class active councils coexist in the portfolio
- scheduler now retains a sticky foreground mission when a fresher runnable peer is otherwise tied on priority, effective risk, and release posture, reducing avoidable foreground thrash and exposing `selection_policy` plus `previous_foreground_mission_id` in durable snapshots
- scheduler now retains a sticky bounded parallel batch when the admitted portfolio is still valid and a fresher same-class peer arrives, exposing `previous_selected_mission_ids` so the leader can see the control plane is preserving the admitted set intentionally
- scheduler now escalates stale runs, expired signals, and stale meetings into explicit maintenance work instead of leaving them only to doctor output
- scheduler now sorts maintenance into a first-class portfolio, exposing `primary_maintenance_action`, deferred-maintenance count, `maintenance_selection_policy`, and `previous_primary_maintenance_key` instead of leaving maintenance as an unordered backlog
- scheduler now unifies those human-decision and maintenance queues into one leader-facing `primary_control_plane_packet`, exposing `primary_control_plane_kind`, `primary_control_plane_selection_policy`, and deferred control-plane count so the control plane itself states whether maintenance or a governed decision is currently first
- when maintenance wins that unified portfolio, scheduler now records the waiting governed decisions as explicitly deferred behind maintenance, including mission-level `deferred_by_primary_control_plane_kind` state and `defer_pending_primary_maintenance` reschedule actions
- governance-critical maintenance such as stale explicit bundle governors or stale active councils can now trigger `hold_for_primary_maintenance`, deferring low-urgency foreground execution until the primary maintenance action is repaired or explicitly deferred
- mission-bundle governor selection is no longer limited to raw `priority` / `risk_level`; larger bundles can now adjudicate from the fuller candidate-governor tuple, including release-adjacent and governance-activity signals when those create a unique winner
- first-class mission objects
- mission authority separated from freshest run activity
- bounded teammate runs
- governed sibling communication
- ownership selectors and governed transfer
- authority audit across multiple state families
- mission-bundle authority that can now let either a more specific ownership boundary or a single active council anchor govern a broader overlapping mission even without a dependency edge
- meeting packets with role/authority discipline
- wave/route/round council orchestration
- protocol-aware council control, where `incident-command` and `decision-record` can defer same-authority disagreement directly into synthesis/record handling instead of auto-generating live debate obligations
- governed verdict adoption in some resolved debate lanes
- runtime-enforced governance gaps that block councils lacking authority hierarchy or rejected-finding reasons

This is enough to call the current system a true prototype line, not just a research stub.

---

## What Is Still Expansion

The operating model still has open future layers:

- fuller scheduler behavior across larger mission portfolios than the current snapshot/reschedule/execution-mode baseline
- broader authority adjudication for larger mission bundles without obvious governing dependency, selector-specific winner, or council outcome
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
