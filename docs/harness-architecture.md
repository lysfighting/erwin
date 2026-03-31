# Harness Architecture

> Status: Active
>
> Date: 2026-03-29
>
> Purpose: define the single authoritative architecture narrative for `Harness`, consolidating the protocol lineage, control-plane model, runtime implementation, design influences, and remaining expansion areas into one durable file.

---

## Why This Document Exists

`Harness` no longer fits safely inside chat-context memory or a scattered reading of research notes.

This document exists to provide one durable architectural source that explains:

- what `Harness` is
- what it is not
- what design inputs shaped it
- how its layers are separated
- how the control plane actually works
- why the Go runtime exists
- what is already implemented
- what remains future expansion rather than core incompleteness
- why stale control-plane state must become explicit maintenance work
- why debate-batch identity and stickiness matter for trustworthy automated council control

This document is descriptive and architectural.
It does not replace:

- the machine-readable package truth in [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)
- the durable external boundary in [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)

It is now also part of a managed documentation system described in:

- [`docs/harness-documentation-map.md`](/Users/liyasong/Project/collov/erwin/docs/harness-documentation-map.md)

---

## Identity

The system has three distinct identities that must not be collapsed:

1. `Protocol`
   `Autonomous Development Lifecycle Protocol`

2. `Adapter / public package line`
   `Harness`

3. `Machine-readable adapter id`
   `codex-native-harness`

This means:

- `Harness` is not a renamed `teamwork`
- `Harness` is not a renamed `wing`
- `Harness` is not merely a skill
- `Harness` is not a generic hosted agent platform

It is:

- a leader-led orchestration harness
- a Codex-native adapter for the protocol
- a reusable package that can be bootstrapped into another git repository

Authoritative positioning references:

- [`_project/decisions/0013-harness-prototype-positioning.md`](/Users/liyasong/Project/collov/erwin/_project/decisions/0013-harness-prototype-positioning.md)
- [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
- [`docs/harness-session-handoff.md`](/Users/liyasong/Project/collov/erwin/docs/harness-session-handoff.md)

---

## Canonical Document Set

`Harness` should no longer be understood through scattered files alone.

The canonical documentation set is:

1. this file for the complete architecture narrative
2. [`docs/harness-operating-model.md`](/Users/liyasong/Project/collov/erwin/docs/harness-operating-model.md) for the leader/teammate/council execution model
3. [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md) for the durable package boundary
4. [`docs/harness-session-handoff.md`](/Users/liyasong/Project/collov/erwin/docs/harness-session-handoff.md) for safe continuation rules
5. [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml) for machine-readable package truth
6. [`docs/harness-documentation-map.md`](/Users/liyasong/Project/collov/erwin/docs/harness-documentation-map.md) for how the whole document system is meant to be read and maintained

This is deliberate:

- one document explains the whole system
- one document explains how it operates
- one document defines the public contract
- one document explains how future sessions should resume
- one file anchors the machine-readable truth

---

## Core Design Goal

The goal is not to port old commands into a new syntax.

The goal is to build a Codex-native orchestration harness that preserves the serious protocol semantics learned from `teamwork`, the modular packaging discipline learned from `wing`, and the stronger governance ambitions from V4, while absorbing `arch-review` as only one input among several decision/governance methodologies instead of letting it become the whole council model.

The design target is:

**leader-led, council-corrected, git-grounded orchestration**

This implies:

- strong lifecycle semantics
- explicit control-plane state
- bounded teammate execution
- governed communication
- authority-aware reconciliation
- deterministic mechanisms outside freeform reasoning
- packaging and distribution discipline
- stable controller identity for active debate work, so councils do not thrash or alias two distinct open batches into one apparent control slot

Key source documents:

- [`docs/research/codex-native-orchestration-harness-analysis-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/codex-native-orchestration-harness-analysis-2026-03-28.md)
- [`docs/research/codex-runtime-adapter-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/codex-runtime-adapter-spec-2026-03-28.md)

---

## What Harness Is Not

The architecture explicitly rejects several tempting but wrong framings:

- not a direct `teamwork` migration
- not a direct `wing` migration
- not a flat multi-agent swarm
- not a dictator-leader system
- not a prompt-monolith disguised as architecture
- not a project-coupled `vi_agent` derivative
- not a Python-bound repo-local script pile

These rejections are important because most historical confusion came from layering multiple concerns into one name or one runtime shell.

---

## Primary Influences

### 1. `teamwork`

`teamwork` contributes the protocol seriousness:

- GitHub-first execution model
- canonical lifecycle `issue -> claim -> drive -> ship -> rc`
- contract as the AI execution artifact
- release governance depth
- deterministic mechanisms outside freeform reasoning
- Git/GitHub orderliness as the backbone of governed execution

Important references:

- [`docs/teamwork-ai-manual.md`](/Users/liyasong/Project/collov/erwin/docs/teamwork-ai-manual.md)
- [`docs/teamwork-design-constitution.md`](/Users/liyasong/Project/collov/erwin/docs/teamwork-design-constitution.md)
- [`_project/decisions/0010-teamwork-v2-github-first.md`](/Users/liyasong/Project/collov/erwin/_project/decisions/0010-teamwork-v2-github-first.md)
- [`_project/decisions/0011-teamwork-closed-loop-design.md`](/Users/liyasong/Project/collov/erwin/_project/decisions/0011-teamwork-closed-loop-design.md)

### 2. `wing`

`wing` contributes packaging and context discipline:

- official skill-style modular decomposition
- smaller invocation surface
- split between thin skill text and supporting operations
- better separation between runtime prompt surface and deterministic mechanisms

Important references:

- [`docs/research/teamwork-vs-wing-analysis-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/teamwork-vs-wing-analysis-2026-03-28.md)
- [`docs/research/parity-matrix-teamwork-wing-protocol-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/parity-matrix-teamwork-wing-protocol-2026-03-28.md)
- [`_project/ideas/skills/tw2-architecture.md`](/Users/liyasong/Project/collov/erwin/_project/ideas/skills/tw2-architecture.md)

### 3. V4 Orchestration Ideas

V4 contributes the orchestration reframe:

- one core brain / leader, not a flat swarm
- teammates as bounded execution units
- optional parallelism under governance, not parallelism by default
- overlap, dependency, and signal awareness
- teammate-to-teammate communication under leader governance
- key decisions escalated into councils rather than silently improvised

What survived from V4 is not “always multi-agent.”
What survived is **leader-centric orchestration capability**.

Primary references:

- [`docs/research/control-plane-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/control-plane-spec-2026-03-28.md)
- [`docs/research/codex-native-orchestration-harness-analysis-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/codex-native-orchestration-harness-analysis-2026-03-28.md)

### 4. `drive`

`drive` contributes the cognitive posture:

- own the mission
- challenge weak framing
- do not silently comply with wrong direction
- treat conversation as part of the work, not overhead
- use standalone-test thinking
- keep next valid steps visible
- do not fake completion

This influence matters because `Harness` is not only an orchestration shell.
It is also meant to preserve a stronger execution attitude than old workflow-only systems captured.

Primary references:

- [`reference/drive-skill/commands/drive.md`](/Users/liyasong/Project/collov/erwin/reference/drive-skill/commands/drive.md)
- [`docs/research/drive-skill-analysis.md`](/Users/liyasong/Project/collov/erwin/docs/research/drive-skill-analysis.md)
- [`docs/research/protocol-core-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/protocol-core-spec-2026-03-28.md)
- [`_project/sessions/2026-03-13-AL.md`](/Users/liyasong/Project/collov/erwin/_project/sessions/2026-03-13-AL.md)

### 5. Engineering Cybernetics

Cybernetics contributes the regulator mindset:

- desired state vs actual state
- disturbance, feedback, correction
- model quality as a prerequisite for good control
- the leader as a regulator, not just a dispatcher

`Harness` therefore treats stale state, drift, and authority conflict as control-plane concerns, not incidental implementation bugs.

Key reference:

- [`docs/research/cybernetics-kubernetes-arch-review-synthesis-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/cybernetics-kubernetes-arch-review-synthesis-2026-03-28.md)

### 6. Kubernetes Control-Plane Thinking

Kubernetes contributes structural discipline:

- control plane vs workers
- desired state vs actual state
- reconciliation loops
- admission and validation gates
- explicit ownership boundaries
- garbage collection as first-class maintenance

`Harness` borrows the control-plane mindset, not the literal platform architecture.

### 7. Harness Engineering

Harness engineering contributes the environment-design frame:

- the system is not “a prompt”
- the product is the environment, control structure, validation, packaging, and feedback loops around the model
- installability, smoke coverage, bootstrap/sync, and manifest-driven reuse are design requirements, not afterthoughts

Important references:

- [`docs/research/harness-incubation-strategy-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-incubation-strategy-2026-03-28.md)
- [`docs/research/erwin-harness-incubation-blueprint-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/erwin-harness-incubation-blueprint-2026-03-28.md)

### 8. `arch-review`

`arch-review` contributes bounded plural judgment:

- the leader as chair, not oracle
- explicit authority hierarchy
- bounded roles
- allowed and forbidden sources
- independent review before cross-examination when conflict remains
- synthesis that records rejected findings and reasons
- governance gates that stop councils from advancing while those bounded-review records are incomplete

This is why meetings in `Harness` are governance objects rather than loose notes.
It is also why `Harness` does not treat cross-examination as mandatory ceremony when a lane is aligned or already settled by higher-authority evidence.

Key references:

- [`skills/arch-review/arch-review/SKILL.md`](/Users/liyasong/Project/collov/erwin/skills/arch-review/arch-review/SKILL.md)
- [`docs/research/meeting-protocol-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/meeting-protocol-spec-2026-03-28.md)
- [`docs/research/state-authority-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/state-authority-spec-2026-03-28.md)

---

## Non-Negotiable Layer Separation

The architecture has three layers and they must remain separate.

### Protocol

Runtime-neutral lifecycle semantics and invariants.

Canonical loop:

`issue -> claim -> drive -> ship -> rc`

The protocol defines:

- GitHub as system of record
- contract as the AI execution artifact
- standalone executability
- deterministic mechanisms outside freeform reasoning where possible
- blast-radius awareness
- completion as state transition plus verification

Primary reference:

- [`docs/research/protocol-core-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/protocol-core-spec-2026-03-28.md)

### Adapter

The Codex-facing implementation layer.

This includes:

- the compiled runtime
- shell entrypoints
- bootstrap and sync flows
- packaging and distribution
- smoke validation
- Codex skill surfaces

The adapter should feel native to Codex, not like a compatibility shim for another runtime.

### Profile

Repo-specific policy.

Examples:

- branch names
- release topology
- verification commands
- ownership conventions
- worktree policy

Current proving-ground profile:

- [`harness/profiles/erwin.yml`](/Users/liyasong/Project/collov/erwin/harness/profiles/erwin.yml)

---

## Definition Layer vs Runtime State

This boundary is architectural, not accidental.

- [`harness/`](/Users/liyasong/Project/collov/erwin/harness/README.md) stores committed definitions
- `.harness/` stores local runtime state

Committed definition layer includes:

- profiles
- templates
- policies
- bootstrap assets
- adapter manifest

Runtime state includes:

- runs
- missions
- signals
- meetings
- contracts
- maintenance artifacts

This split exists to preserve:

- repo-visible architecture
- runtime mutability without polluting committed source
- upgrade and normalization paths
- bootstrapable reuse

---

## Runtime Implementation

The core runtime is implemented in Go.

Entry points:

- [`cmd/harness-state/main.go`](/Users/liyasong/Project/collov/erwin/cmd/harness-state/main.go)
- [`internal/harness/runtime.go`](/Users/liyasong/Project/collov/erwin/internal/harness/runtime.go)

The shell layer is intentionally thin:

- [`scripts/harness-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-runtime.sh)
- [`scripts/harness-*.sh`](/Users/liyasong/Project/collov/erwin/scripts)

Runtime selection order:

1. `harness/bin/<os>-<arch>/harness-state`
2. `.harness/bin/harness-state`
3. local `go build`

This design was chosen to avoid Python runtime coupling and to move toward a distributable binary-first control plane.

Consequences:

- normal `Harness` usage no longer depends on `python3`
- packaged binaries are the preferred user path
- local Go toolchain is a fallback, not the ideal steady state

Distribution references:

- [`docs/harness-distribution-strategy.md`](/Users/liyasong/Project/collov/erwin/docs/harness-distribution-strategy.md)
- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)

---

## Why Go, Not Python

The original prototype concentrated the deterministic state engine in Python because it was fast to iterate on structured runtime logic.

That was acceptable for incubation, but not for the intended long-term package shape.

The Go migration was driven by real architectural concerns:

- avoid depending on project-local Python environments
- avoid Python version drift across unrelated repos
- make the runtime feel like a real CLI/control plane
- move toward packaged-binary distribution
- keep deterministic runtime logic in a compiled, portable substrate

So the architecture now treats:

- Go runtime = core deterministic engine
- shell wrappers = stable human-facing entrypoints
- skills = thin guidance surfaces

This is deliberate, not incidental.

---

## Public Package Contract

The durable public boundary is defined in:

- [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)

The machine-readable package truth is:

- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)

The adapter currently guarantees:

- manifest-driven packaging
- reusable bootstrap and sync
- deterministic runtime engine
- stable shell surfaces
- smoke-testable integrity
- mission/run/meeting/state persistence
- leader-facing `home` and `doctor`

This package is already a reusable adapter line, not just a local experiment.

---

## Core Protocol Semantics

The protocol invariants that must survive any future runtime evolution are:

1. GitHub is the system of record
2. issue is the public collaborative artifact
3. contract is the AI execution artifact
4. standalone executability matters at issue and contract time
5. sub-task decomposition is challengeable, not sacred
6. deterministic mechanisms should be mechanized
7. quality is enforced by feedback loops
8. blast radius is first-class
9. completion is behavioral, not textual

Those invariants are inherited from `teamwork`, clarified during protocol extraction, and preserved through the adapter.

---

## Control-Plane Thesis

`Harness` is not just a lifecycle wrapper.
It is a control plane.

The control plane exists to answer:

- what can run now
- what should wait
- what must remain isolated
- what can run in parallel safely
- what evidence should change state
- what should be escalated
- what garbage should be cleaned

This is why the system has first-class objects above individual shell commands.

Primary reference:

- [`docs/research/control-plane-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/control-plane-spec-2026-03-28.md)

---

## Leader Model

The leader is the control-plane regulator and chair.

The leader is responsible for:

- maintaining the mission graph
- assigning bounded execution
- interpreting signals
- convening councils
- reconciling desired vs actual state
- deciding when system truth changed

The leader is not:

- a flat scheduler only
- an oracle who should decide alone
- a worker expected to do all local execution personally

This is a chairing regulator model, not dictator leadership and not flat swarm coordination.

References:

- [`docs/research/leader-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/leader-spec-2026-03-28.md)
- [`docs/research/cybernetics-kubernetes-arch-review-synthesis-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/cybernetics-kubernetes-arch-review-synthesis-2026-03-28.md)

---

## Teammate Model

Teammates are bounded execution units.

They exist to:

- own delegated local work
- report progress and evidence
- raise blockers
- participate in councils
- communicate through governed channels

They do not exist to replace the leader’s global model.

This means:

- teammates are first-class runtime objects
- teammate work must have explicit scope
- teammate communication must be governed
- teammate reports are evidence, not unquestioned truth

References:

- [`docs/research/teammate-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/teammate-spec-2026-03-28.md)

---

## Mission, Run, and Graph Model

`Harness` deliberately separates conceptual work objects from execution attempts.

### Mission

A first-class durable work object above runs.

Important fields include:

- mission id
- issue id
- state
- latest run
- authoritative run
- objective
- desired state
- ownership surface

### Run

A concrete execution attempt for a mission.

Runs carry:

- runtime state
- desired state
- worktree or branch linkage
- ownership selectors
- exclusions
- authority-audit surface

### Mission Graph

The graph connects missions by:

- dependency
- overlap
- release coupling
- ownership contention

This graph supports:

- bundle grouping
- bundle governors
- explicit bundle-governor adoption by leader decision
- structured bundle-governor adoption from closed bundle councils
- durable provenance for council-adopted bundle governors
- pairwise governor adjudication from clear authority tuples
- dependency-governed authority
- selector-specificity authority
- council-governed authority
- leader-facing required actions

This is one of the ways `Harness` moved beyond a plain run registry.

When two active missions overlap without a dependency edge, the graph does not have to fall straight to generic conflict.

If one mission holds the more specific ownership boundary and no competing higher-authority evidence contradicts that shape, the control plane may adjudicate that narrower mission as the current bundle governor.

If one mission is already chairing the relevant active council, the control plane may also keep bundle authority anchored to that council-host mission until the governed meeting outcome is adopted.

If more than one overlapping mission is simultaneously carrying an active council anchor on the same bundle and no single durable governor already exists, the control plane now escalates that state explicitly instead of pretending it is an ordinary overlap conflict.

The scheduler now consumes these bundle packets directly rather than leaving them as audit-only output.

That means:

- the governing mission can remain foreground-runnable
- governed losing missions can be deferred behind the current bundle governor
- escalation-required bundles can surface as explicit human decisions instead of generic overlap noise
- escalation-required bundles can now expose ranked candidate governors, so the leader sees likely control anchors before deciding whether to establish one directly or escalate into a council
- overlap bundles with multiple active council anchors now get their own escalation shape as well, surfacing `multi_council_bundle_governor_required` plus guidance to consolidate the competing councils or establish one governing mission before widening execution
- the leader can now explicitly adopt a governor for an escalation-required bundle, converting that bundle into governed authority through a durable control-plane signal instead of leaving it as a perpetual decision packet
- a bundle-scoped council can now also close with a structured governing-mission verdict, and `meeting-close` will adopt that verdict into the same durable explicit-governor signal path instead of leaving the outcome trapped in archived meeting notes
- when auto-synthesis already produces a clear accepted winner for a bundle-scoped `establish_bundle_governor` council, the controller can now auto-generate that structured governing-mission verdict instead of forcing the leader to insert a manual lock step before close
- once that adoption happens, the resulting governed bundle now keeps provenance about how the governor was established, so graph authority and scheduler surfaces can distinguish `council-adopted-governor` from a direct leader-established explicit governor
- stale explicit governors are not trusted indefinitely; when the adopted governor no longer belongs to the active bundle, the control plane escalates that signal into maintenance instead of silently carrying orphaned governance forward
- the scheduler can now retain a sticky primary human-decision packet across otherwise-tied portfolio passes, preferring decision stability over freshness-only churn when the decision class still matches
- when multiple same-class active councils coexist, the scheduler can now also prioritize by live council progress (`active_decision_progress_rank`) so a council already in synthesis/output outranks a fresher but less advanced peer
- the scheduler now gives governed human-decision packets stable `decision_id` values and projects them back into mission classifications and defer actions, creating a more durable bridge between scheduler output and future council materialization
- the scheduler-to-council bridge is now real rather than hypothetical: leader-facing meeting surfaces can materialize the current primary decision packet or a specific `decision_id` into an active meeting, and that meeting persists the source decision packet plus scheduler snapshot linkage for later audit
- scheduler decision materialization is now idempotent with respect to live governance: if the same `decision_id` is already under active council control, materialization reuses the existing meeting instead of minting a duplicate council for the same control-plane object
- that same `decision_id` now stays durably re-addressable even after the scheduler packet has promoted into `active-council`; when the latest snapshot no longer carries the original pre-materialization packet, runtime can still recover the decision through the live meeting's persisted `source_decision_packet` and refresh the governed council from current scheduler context
- when that materialized decision is a bundle-scoped `establish_bundle_governor`, the meeting now starts with an auto-seeded participant roster, sibling mission refs, and candidate-governor packet summaries, so the governed council begins from bounded scheduler evidence instead of a blank roster
- those decision ids are now stable across surface changes, so once a decision is materialized the scheduler can promote it into `active-council` state without inventing a new decision identity or breaking later reconciliation
- when the materialized decision is bundle-scoped, scheduler now projects that live council back into every affected mission through bundle-level active-council fields, counts active bundle councils at the portfolio level, and can hold ready/executing peers behind the live bundle council instead of treating the bundle as generic unresolved governance
- active bundle-governor councils now also remain bundle-scoped in the scheduler portfolio itself: runtime can recover the bundle decision packet from the live meeting when needed, so the leader continues to see one durable bundle decision object instead of a degraded set of mission-level `review_active_council` packets
- active-council decision packets now also reflect live meeting progress, projecting phase, protocol, next requirement, and current debate counters back into the scheduler portfolio so the leader sees what the council needs now, not only the stale pre-materialization decision prompt
- scheduler and `home` now follow that live council requirement directly, so once a bundle-scoped council is already in flight the primary next action and `hold_for_active_bundle_council` reason inherit the meeting's current bounded requirement rather than generic “review the council” wording
- active-council packets now also project staleness (`active_decision_last_updated_at`, `active_decision_age_hours`, `active_decision_stale`), and stale live councils now outrank routine bundle holds so the control plane explicitly says to refresh or close the stale council before trusting it as current governance
- active-council packets and meeting summaries now also project the concrete current-wave action packets plus blocked batch summaries, so live council execution is visible as real queued work rather than only as a batch id or generic requirement string
- protocol-aware council execution now reaches farther into the live meeting loop itself: `incident-command` councils ask for bounded incident assessments and close toward action/owner/stop-go output, while `decision-record` councils ask for bounded decision-record assessments and close toward explicit decision/alternatives/consequences/follow-up output rather than falling back to evidence-review phrasing
- that protocol-aware execution now also materializes provisional output artifacts directly in the control plane, including lane-level synthesis and structured `follow_up_packets` inside the bounded `output_packet`, and closed-meeting adoption now writes that artifact back into run state so councils do not stop automation at findings/action items when the runtime already has enough structure to preserve per-lane dispositions and bounded follow-up execution
- scheduler mission packets and `home` now also project that adopted follow-up state back into leader-facing execution surfaces, including the adopted packet kind, follow-up counts, and primary adopted follow-up packet, so closed councils remain visible as bounded execution state instead of only as raw run metadata
- run adoption now also promotes the primary adopted follow-up packet into `next_expected_action` when one is open, and authority audit treats missing promotion as adoption drift, so durable run guidance stays aligned with the closed council's bounded next task
- resolving an action on a closed adopted meeting now also refreshes the governing run's adopted follow-up packets and primary adopted packet, so the runtime can keep the adopted execution object synchronized with real completion instead of freezing packet state at the original `meeting-close`
- the same adopted execution object now also has a run-native completion path through `drive resolve-follow-up`, so post-council bounded execution can advance from the run control plane itself instead of remaining coupled to the closed meeting surface
- that bounded adopted execution object can now also project directly into teammate runtime through `teammates assign --adopted-follow-up`, preserving governing packet provenance while creating one bounded delegated execution unit and refusing duplicate active handoff of the same packet
- once projected into teammate runtime, that same bounded object can now also resolve back through the delegated teammate run itself; runtime keeps the governing run as the durable truth while still letting the bounded execution unit close the packet it owns
- that governing run remains the durable packet source even when it re-addresses the adopted follow-up directly, because delegated teammate state now resyncs from the governing packet truth instead of drifting on a stale open snapshot, including reopening the delegated bounded execution unit when the governing packet itself is reopened
- that delegated ownership is now also reflected back into scheduler/home leader surfaces, including the delegated follow-up state and freshness itself, so the control plane can distinguish “leader still owns this adopted follow-up directly” from “leader is now governing a delegated execution unit that owns it, it remains open, and it may now require stale-follow-up recovery rather than passive monitoring”
- when that delegated execution unit actually goes stale while still owning the bounded follow-up packet, the scheduler now promotes that condition into a dedicated critical maintenance object instead of a generic stale-run warning, keeping delegated adopted-follow-up recovery inside the main control-plane priority system
- the diagnostic layer now also reflects that same specialization: `doctor` names stale delegated adopted-follow-up ownership explicitly instead of flattening it into generic run staleness
- completion gating now follows that same delegation path as well, so a delegated teammate run cannot claim closure while the governing adopted follow-up packet it owns is still open
- that run control plane now also gates `completed` on the adopted follow-up lifecycle, preventing the run from claiming closure while any adopted post-council packet still remains open
- once a closed meeting's structured output has already been adopted, its remaining open action items no longer count as generic meeting backlog in scheduler/home mission tallies, so the runtime can hand control cleanly from council execution into the adopted follow-up lifecycle without double-counting the same work object
- explicit bundle-governor adoption is now checked against live bundle councils as well; when an adopted governor points at a different mission than the active bundle-council host, the control plane escalates that mismatch into maintenance so leader-governed authority cannot silently split across two anchors
- the scheduler can now retain a sticky bounded parallel batch across otherwise-tied portfolio passes, preferring admitted-batch stability over freshness-only churn when the batch class still matches
- the scheduler now treats maintenance as a first-class portfolio as well, surfacing `primary_maintenance_action`, deferred-maintenance count, and sticky maintenance-selection metadata instead of leaving stale governance cleanup as an unordered side list
- the scheduler now also collapses those decision and maintenance portfolios into one leader-facing `primary_control_plane_packet`, exposing `primary_control_plane_kind`, `primary_control_plane_selection_policy`, and deferred control-plane count so the leader can see which governance object truly comes first without mentally merging two separate queues
- when critical maintenance does outrank governed human decisions, the control plane now records that deferral explicitly on affected mission classifications and emits `defer_pending_primary_maintenance`, so “maintenance first” is runtime state rather than only portfolio prose
- governance-critical maintenance such as stale explicit bundle governors or stale active councils can now hold back low-urgency foreground execution through `hold_for_primary_maintenance`, so the control plane repairs broken governance before widening fresh work
- the scheduler can now hold back brand-new low-urgency foreground admission when unresolved human decisions already dominate the portfolio and no established foreground batch is being continued
- parallel widening now requires explicit justification when execution policy says the control plane should reduce parallelism rather than guess
- the scheduler can now retain a sticky foreground mission across otherwise-tied passes, preferring stability over timestamp churn when priority, effective risk, and release posture still match
- larger mission bundles can now adjudicate from a fuller candidate-governor tuple rather than only raw priority/risk, letting a uniquely release-adjacent or governance-active mission govern even when the simpler tuple still ties
- a runnable `P0` mission now serializes the foreground batch even when other runnable missions are otherwise isolated
- profile policy can now also escalate a foreground mission into high-risk serialization when its governed ownership boundary matches configured high-risk scopes, even if no one manually set `risk_level`
- profile meeting-trigger policy can now mark some human decisions as explicitly council-worthy, so the leader-facing scheduler does not flatten all escalations into the same generic decision surface

---

## Ownership Model

Ownership is explicit and governed.

Supported selector families include:

- `file`
- `dir`
- `module`
- `glob`
- `scope`
- `interface`

The system enforces:

- path-family overlap
- semantic same-boundary overlap
- hierarchical semantic overlap
- exclusion-aware delegation

Supported transfer modes include:

- `move`
- `copy`
- `split`

`split` is especially important:

it allows a parent boundary to stay with one run while an explicitly carved child subtree is delegated elsewhere.

This is one of the strongest Kubernetes-inspired design lines in the current runtime.

---

## Desired State and Reconciliation

The runtime explicitly models desired state instead of only recording what happened.

Desired state exists on:

- contracts
- runs
- missions

Reconciliation happens through:

- `home`
- `doctor`
- authority audits
- maintenance reconcile
- maintenance normalize

This is where the cybernetic and Kubernetes influence is most direct:

- compare intended state to observed state
- surface drift
- reconcile or escalate

The design does not assume one command makes the world stay correct.

---

## Authority Model

The system uses an explicit authority hierarchy because not all artifacts deserve equal trust.

Canonical levels:

- `L1` live system-of-record state
- `L2` active orchestration state
- `L3` current design/planning artifacts
- `L4` auxiliary or derived artifacts
- `L5` legacy or deprecated artifacts

Authority is resolved by:

1. higher authority over lower authority
2. fresher over staler
3. direct over derived
4. specific over vague
5. accepted synthesis over raw individual opinion

Authority is now applied across:

- signals
- meetings
- runs
- contracts
- live git state
- missions
- mission bundles

This is one of the major reasons `Harness` is a control plane instead of a workflow wrapper.

References:

- [`docs/research/state-authority-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/state-authority-spec-2026-03-28.md)

---

## Meetings and Council Governance

Meetings are not archival notes.
They are governance objects.

They exist when:

- evidence conflicts
- blast radius is high
- design quality matters
- mission decomposition may change
- release risk is meaningful

Every meeting should carry:

- question
- scope
- authority hierarchy
- runtime-enforced governance gaps for missing authority hierarchy and missing rejected-finding reasons
- roles
- known context
- locked decisions
- expected output

The council model includes:

- bounded participants
- role-specific missions
- allowed and forbidden sources
- phase history
- independent review
- cross-examination
- synthesis
- output
- rejected reasons
- findings with provenance
- direct output closure when synthesis is already ready

This is the clearest place where `arch-review` directly shaped `Harness`.

### 9. Incident Command / Incident Management

Incident-management practice contributes a different decision rhythm than technical debate:

- stabilize first
- assign bounded action ownership
- keep explicit checkpoints
- do not let broad architecture debate delay immediate operational control

This is especially relevant for:

- hotfix / rollback / RC decisions
- release-affecting uncertainty
- urgent operational coordination

### 10. Decision Records / ADRs

Decision records contribute durable decision memory:

- context
- alternatives
- chosen direction
- consequences
- follow-up obligations

This matters because not every important decision needs a live debate controller.
Some need bounded recording more than cross-examination.

---

## Decision Protocol Family

`Harness` now treats decision governance as a small family of protocols, not one inherited ritual.

The current protocol taxonomy is:

1. `evidence-review`
   For contested technical or factual judgment.

2. `incident-command`
   For urgent release-adjacent or operational decisions where bounded action ownership matters more than extended debate.

3. `decision-record`
   For lower-urgency design and tradeoff decisions where durable recording matters more than live rebuttal rounds.

Important rule:

- `arch-review` remains an important influence on `evidence-review`
- it is **not** the whole semantics of council/decision governance in `Harness`

The broader methodology synthesis is recorded in:

- [`docs/research/harness-decision-methodology-synthesis-2026-03-30.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-decision-methodology-synthesis-2026-03-30.md)

---

## Automated Council Behavior

Council automation is intentionally bounded.

Current capabilities include:

- phase orchestration
- wave obligations
- multi-lane debate wave batching
- deterministic current/blocked debate-batch execution control
- route obligations
- multi-round debate planning
- protocol-aware direct synthesis for `incident-command` and `decision-record` when same-authority conflict should be recorded rather than expanded into live route/round work
- same-pass closure from synthesis to output when no synthesis gaps remain
- summary and readiness metrics
- same-claim authority audit
- latest-response verdict adoption
- provisional auto-synthesis

That means the council system is no longer merely descriptive.
It actively helps govern debate.

But this is still not the fully automated end state imagined by the earliest control-plane research.

---

## Leader-Led, Council-Corrected Orchestration

One of the most important architectural corrections is:

the leader should not decide everything alone.

`Harness` is designed around:

- leader-led control
- teammate execution
- governed sibling communication
- council correction when ambiguity or blast radius is high

This avoids both failure modes:

- dictator leader failure
- flat swarm failure

That balance is central to the design.

---

## Git and GitHub as the State Backbone

Git and GitHub are not just integrations.
They are the ordered substrate that makes regulation possible.

Important control-plane objects include:

- issue
- contract
- branch
- worktree
- commit
- PR
- checks
- release candidate
- release state

These objects provide:

- inspectability
- partial ordering
- ownership markers
- review gates
- rollback points
- durable history

This is why the architecture remained GitHub-first even while packaging, runtime, and naming changed.

---

## Deterministic Mechanisms vs Judgment

One of the oldest and most stable lessons from `teamwork` is the Von Neumann-style separation:

- deterministic mechanisms belong in scripts/runtime
- adaptive judgment belongs in the model

In the current architecture:

- Go runtime holds deterministic state and control logic
- shell wrappers provide stable entrypoints
- skills provide compact operational guidance

This separation is essential to:

- testability
- portability
- packaging clarity
- resistance to prompt drift

---

## Skill Surface

Skills remain important, but they are not the whole system.

`Harness` uses Codex skills as:

- operator entrypoints
- control-plane guidance
- concise surface documentation

The skill should remain thin because the real system is:

- runtime
- manifest
- definition layer
- profile contract
- maintenance and smoke coverage

This is a direct lesson from the move away from monolithic prompt systems.

Relevant surface:

- [`skills/harness/SKILL.md`](/Users/liyasong/Project/collov/erwin/skills/harness/SKILL.md)

---

## Packaging, Bootstrap, and Distribution

`Harness` is designed to be reusable across repositories.

That is why it includes:

- bootstrap
- sync
- adapter validation
- smoke suites
- runtime packaging

Distribution is tiered:

1. source bootstrap
2. in-repo packaged binaries
3. portable runtime bundles

This is part of the architecture because harness engineering requires the environment to be distributable and verifiable, not merely described.

---

## Current Implemented Maturity

`Harness` already delivers the first stable line.

That includes:

- protocol extraction
- adapter identity
- manifest-driven packaging
- Go runtime
- binary-preferred runtime path
- leader control plane
- signal-driven scheduler rescheduling across normalized meeting/review/priority/risk/rebase follow-up kinds, including explicit urgency ranking for fresh follow-up portfolio work
- teammate runtime
- mission objects
- mission graph and bundle authority
- governed ownership and transfer
- authority audit
- council objects and orchestration
- reusable bootstrap/sync/distribution path
- smoke-testable package integrity

This is important because the remaining gaps should be understood as expansion above a working stable line, not evidence that the core architecture is still missing.

---

## Remaining Expansion Areas

The current architecture is intentionally honest about what is not fully closed.

High-value future expansion areas are:

1. fuller scheduler behavior
   deeper portfolio scheduling policy beyond the now-real baseline for execution-mode selection, mission classification, durable scheduler snapshots, prioritized human-decision portfolios, leader-facing rescheduling hints, and maintenance escalation for stale control-plane state

2. broader bundle adjudication
   larger multi-mission conflicts when no obvious governing dependency, selector-specific authority winner, signal, or council outcome exists

3. deeper automated council execution
   fuller `arch-review`-style automation across more than the current wave/route/round/verdict-adoption baseline

4. optional Codex sub-agent mapping
   whether and how real Codex sub-agent execution should map directly onto teammate runtime

These are future layers, not reasons to re-open the core identity or architecture split.

---

## Relationship To Research

This document consolidates the architecture, but research still matters.

The most important research files remain:

- [`docs/research/protocol-core-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/protocol-core-spec-2026-03-28.md)
- [`docs/research/control-plane-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/control-plane-spec-2026-03-28.md)
- [`docs/research/meeting-protocol-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/meeting-protocol-spec-2026-03-28.md)
- [`docs/research/state-authority-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/state-authority-spec-2026-03-28.md)
- [`docs/research/harness-gap-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-gap-audit-2026-03-29.md)
- [`docs/research/harness-research-drift-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-research-drift-audit-2026-03-29.md)

This document should reduce the need to stitch those files together during normal continuation work.
It does not make them irrelevant.

---

## Canonical Fresh-Session Reading Order

For a new session that needs enough context to work safely, the preferred read order is:

1. [`AGENTS.md`](/Users/liyasong/Project/collov/erwin/AGENTS.md)
2. [`_project/CURRENT.md`](/Users/liyasong/Project/collov/erwin/_project/CURRENT.md)
3. this file
4. [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
5. [`docs/harness-session-handoff.md`](/Users/liyasong/Project/collov/erwin/docs/harness-session-handoff.md)
6. [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)
7. targeted research and recent session logs as needed

This sequence preserves the design intent while keeping context loading disciplined.

---

## Bottom Line

`Harness` is the current Codex-native adapter line for the Autonomous Development Lifecycle Protocol.

Its architecture is defined by:

- GitHub-first lifecycle semantics
- protocol / adapter / profile separation
- committed definitions vs local runtime state
- Go-based deterministic runtime
- leader-led control plane
- bounded teammate execution
- authority-aware reconciliation
- council-corrected governance
- manifest-driven packaging and reuse

It carries forward the best lessons from `teamwork`, `wing`, cybernetics, Kubernetes control-plane thinking, harness engineering, and `arch-review`, while deliberately rejecting their historical baggage.
