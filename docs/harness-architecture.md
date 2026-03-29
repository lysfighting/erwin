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

2. `Adapter / public prototype line`
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
- pairwise governor adjudication from clear authority tuples
- dependency-governed authority
- selector-specificity authority
- council-governed authority
- leader-facing required actions

This is one of the ways `Harness` moved beyond a plain run registry.

When two active missions overlap without a dependency edge, the graph does not have to fall straight to generic conflict.

If one mission holds the more specific ownership boundary and no competing higher-authority evidence contradicts that shape, the control plane may adjudicate that narrower mission as the current bundle governor.

If one mission is already chairing the relevant active council, the control plane may also keep bundle authority anchored to that council-host mission until the governed meeting outcome is adopted.

The scheduler now consumes these bundle packets directly rather than leaving them as audit-only output.

That means:

- the governing mission can remain foreground-runnable
- governed losing missions can be deferred behind the current bundle governor
- escalation-required bundles can surface as explicit human decisions instead of generic overlap noise
- the scheduler can now retain a sticky primary human-decision packet across otherwise-tied portfolio passes, preferring decision stability over freshness-only churn when the decision class still matches
- parallel widening now requires explicit justification when execution policy says the control plane should reduce parallelism rather than guess
- the scheduler can now retain a sticky foreground mission across otherwise-tied passes, preferring stability over timestamp churn when priority, effective risk, and release posture still match
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

`Harness` already delivers the first complete prototype line.

That includes:

- protocol extraction
- adapter identity
- manifest-driven packaging
- Go runtime
- binary-preferred runtime path
- leader control plane
- teammate runtime
- mission objects
- mission graph and bundle authority
- governed ownership and transfer
- authority audit
- council objects and orchestration
- reusable bootstrap/sync/distribution path
- smoke-testable package integrity

This is important because the remaining gaps should be understood as expansion above a working prototype line, not evidence that the core architecture is still missing.

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
