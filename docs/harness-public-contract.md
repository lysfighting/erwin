# Codex-Native Harness Public Contract

> Status: Active
>
> Date: 2026-03-29
>
> Scope: define the durable public contract for the Codex-native orchestration harness package.

---

## Position

This package is the current public contract for the emerging Codex-native orchestration harness.

Its public prototype name is:

- `Harness`

The naming/positioning decision is recorded in:

- [`0013-harness-prototype-positioning.md`](/Users/liyasong/Project/collov/erwin/_project/decisions/0013-harness-prototype-positioning.md)

It should be understood as:

- a `leader-led orchestration harness`
- a Codex-native adapter for the autonomous development lifecycle protocol
- a reusable package that can be bootstrapped into another git repository

It should **not** be understood as:

- a port of `teamwork`
- a port of `wing`
- a general-purpose standalone cloud agent platform

The package contract is still designed so the protocol can outlive any future branding change.

---

## Protocol Boundary

The harness package implements the protocol defined in:

- [`protocol-core-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/protocol-core-spec-2026-03-28.md)

The required lifecycle remains:

`issue -> claim -> drive -> ship -> rc`

The package must preserve these semantics:

- GitHub-native lifecycle and state backbone
- Contract as the AI execution artifact
- deterministic operations mechanized outside freeform reasoning where possible
- config-driven repo policy
- visible next valid step
- completion defined by real state transition plus verification

---

## Public Package Shape

The package is split into three layers:

1. `protocol`
   Runtime-neutral lifecycle semantics and invariants.

2. `adapter`
   Codex-facing runtime implementation, CLI substrate, shell surfaces, bootstrap/sync flow, and validation.

3. `profile`
   Repo-specific policy such as branches, verification commands, release topology, and conventions.

The committed definition layer lives in:

- [`harness/`](/Users/liyasong/Project/collov/erwin/harness/README.md)

The local runtime state layer lives in:

- `.harness/`

This split is part of the contract, not an implementation accident.

---

## Required Artifacts

The package must provide:

- a machine-readable adapter manifest
- a deterministic runtime engine
- stable shell surfaces
- committed templates and policy defaults
- bootstrap and sync flows
- smoke validation flows
- a Codex skill entrypoint

The current machine-readable source of truth is:

- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)

That manifest may evolve, but the public prototype line should continue to present itself as `Harness` unless a later ADR explicitly supersedes that positioning.

The manifest is allowed to evolve, but these categories are part of the durable contract:

- protocol identity
- adapter identity and version
- package file list
- bootstrap contract
- profile contract
- lifecycle surface declarations
- capability declarations
- smoke script declarations

---

## Required Surfaces

The adapter must expose first-class surfaces for:

- `missions`
- `issue`
- `claim`
- `drive`
- `teammates`
- `ship`
- `rc`
- `meetings`
- `maintenance`
- `adapter`
- `home`
- `doctor`

These surfaces may be implemented as thin shell entrypoints over a compiled runtime.
Their exact internal implementation may change.
Their lifecycle meaning must not.

---

## Runtime Guarantees

The package currently guarantees:

- file-backed runtime state with locking
- leader-facing diagnosis via `home` and `doctor`, with `home` mirroring scheduler decision-first guidance when bundle-governor decisions outrank simple foreground selection, surfacing scheduler maintenance backlog when stale state must be cleaned, and exposing the current primary maintenance packet when governance repair outranks fresh foreground work
- leader-facing scheduler summaries and durable scheduler snapshots under `.harness/scheduler/snapshots/`, including mission classification, execution-mode recommendation, bundle-aware deferral/escalation, bundle-role tallies, bundle-scoped human-decision packets for escalation-required bundles, ranked `candidate_governor_ids` plus `candidate_governor_packets` when a bundle still needs an explicit governor, prioritized human-decision portfolios with `primary_required_human_decision` plus scope/surface tallies, mission-level queue visibility through `human_decision_queue_position` and `primary_human_decision`, explicit `defer_pending_primary_human_decision` actions for lower-priority governed decisions, decision-first foreground holds via `hold_for_primary_human_decision` when unresolved human decisions should block new low-urgency foreground work, sticky primary-decision retention via `decision_selection_policy` and `previous_primary_human_decision_key`, sticky foreground retention via `selection_policy` and `previous_foreground_mission_id`, sticky bounded-parallel retention via `previous_selected_mission_ids`, primary maintenance portfolios with `primary_maintenance_action`, deferred-maintenance count, `maintenance_selection_policy`, and `previous_primary_maintenance_key`, governance-first foreground holds via `hold_for_primary_maintenance` when critical maintenance should block low-urgency foreground work, bounded parallel admission from profile policy, explicit parallel-justification gating, P0 foreground serialization, profile-backed high-risk scope inference, profile-backed council recommendation from meeting triggers, explicit `recommended_decision_protocol` guidance for human decisions, fresh follow-up-signal rescheduling, direct leader adoption and clearing of explicit bundle governors, and next-action prioritization when human bundle decisions or governance maintenance are still pending
- leader-facing scheduler summaries and durable scheduler snapshots under `.harness/scheduler/snapshots/`, including mission classification, execution-mode recommendation, bundle-aware deferral/escalation, bundle-role tallies, bundle-scoped human-decision packets for escalation-required bundles, ranked `candidate_governor_ids` plus `candidate_governor_packets` when a bundle still needs an explicit governor, stable `decision_id` values on required human-decision packets, prioritized human-decision portfolios with `primary_required_human_decision` plus scope/surface tallies, mission-level queue visibility through `human_decision_queue_position`, `primary_human_decision`, and `human_decision_id`, explicit `defer_pending_primary_human_decision` actions for lower-priority governed decisions, decision-first foreground holds via `hold_for_primary_human_decision` when unresolved human decisions should block new low-urgency foreground work, sticky primary-decision retention via `decision_selection_policy` and `previous_primary_human_decision_key`, sticky foreground retention via `selection_policy` and `previous_foreground_mission_id`, sticky bounded-parallel retention via `previous_selected_mission_ids`, primary maintenance portfolios with `primary_maintenance_action`, deferred-maintenance count, `maintenance_selection_policy`, and `previous_primary_maintenance_key`, governance-first foreground holds via `hold_for_primary_maintenance` when critical maintenance should block low-urgency foreground work, bounded parallel admission from profile policy, explicit parallel-justification gating, P0 foreground serialization, profile-backed high-risk scope inference, profile-backed council recommendation from meeting triggers, explicit `recommended_decision_protocol` guidance for human decisions, fresh follow-up-signal rescheduling, direct leader adoption and clearing of explicit bundle governors, and next-action prioritization when human bundle decisions or governance maintenance are still pending
- decision-aware meeting materialization through `meetings materialize-primary-decision` and `meetings materialize-decision DECISION_ID`, including durable `source_decision_*` metadata, scheduler snapshot refs, and persisted source decision packets on the resulting meeting object
- scheduler recognition of materialized decision packets as `active-council`, including stable `decision_id` values across surface changes and `active_meeting_ids` projected back into the portfolio packet
- bundle-aware recognition of materialized decision packets as live governance, including `bundle_active_decision_meeting_ids`, `bundle_active_decision_meeting_count`, `bundle_control_state`, `active_council_bundle_count`, and explicit `hold_for_active_bundle_council` reschedule actions when a bundle is already under active council control
- leader-facing scheduler maintenance escalation, so stale runs, expired signals, and stale meetings become explicit `maintenance_actions` instead of remaining doctor-only observations
- stale explicit bundle-governor signals also become explicit `maintenance_actions`, so bundle governance that no longer matches the active mission graph is surfaced as cleanup work instead of lingering silently
- first-class mission state above run attempts
- first-class mission authority, including an authoritative mission-control run that can differ from the freshest teammate activity
- explicit ownership selectors with overlap enforcement across both path-family and hierarchical semantic-family boundaries
- bounded ownership transfer, including exclusive `move` handoff, shared `copy` handoff for selectors that are intentionally governed in parallel, and exclusion-aware `split` handoff where a parent boundary delegates a governed child subtree without surrendering the broader parent
- desired-state tracking on contracts and runs
- bounded teammate assignment and governed sibling communication
- council-grade meeting packets with explicit `decision_protocol`, `decision_protocol_reason`, and authority-aware findings
- bounded council orchestration with route, wave, and round obligations, including multi-lane debate wave batches that group simultaneous lane work by participant and round, deterministic current/blocked batch execution state, sticky retention of the current open debate batch when a new same-phase lane appears, stable debate-batch identities that avoid controller collisions, direct synthesis when a lane is aligned or already settled by higher-authority evidence, protocol-aware direct synthesis for `incident-command` and `decision-record` when same-authority conflict should be recorded instead of expanded into auto-generated route/round work, same-pass closure from `synthesis` to `output` whenever the synthesis audit is already ready and no open debate batch remains, plus governed latest-response verdict adoption when a debated lane resolves before hitting the cap
- runtime-enforced council governance gates, so missing authority hierarchy and rejected findings without reasons block council progression instead of remaining doctor-only warnings
- authority audit across signals, meetings, runs, and live git state
- doctor-level validation that stale explicit bundle-governor signals are invalid governance rather than healthy active control-plane state
- mission-level authority audit across run aggregates, leader/teammate control boundaries, and multi-run mission conflicts
- graph-level mission authority can adjudicate a pairwise conflicted bundle into a governed bundle when one mission has a uniquely stronger explicit authority tuple
- graph-level mission authority audit across cross-mission ownership conflicts, dependency-governed execution boundaries, dependency-governed ownership adjudication, selector-specificity ownership adjudication, council-governed ownership adjudication, bundle-aware mission-graph grouping, and explicit bundle governors when a governed bundle has either a single root, one clear authority winner, one active council anchor, one leader-adopted explicit governor signal, or one uniquely strongest full candidate-governor tuple
- bundle-level mission authority packets that summarize shared selectors, `resolution_kind`, leader-facing `required_action`, and ranked candidate governors for connected mission bundles that still need an explicit governor
- deterministic ship and rc control-plane actions
- bootstrapable reuse in another git repository
- smoke-testable package integrity

These guarantees are part of the intended stable behavior.
Their implementation detail is not.

---

## Distribution Contract

The runtime core is implemented as a Go CLI.

Normal package behavior is:

1. prefer a packaged binary under `harness/bin/<os>-<arch>/`
2. if no packaged binary exists, fall back to local Go build

This means:

- Python is **not** a runtime dependency for harness usage
- packaged-binary distribution is the preferred product direction
- local Go toolchain use is a fallback, not the ideal end-user expectation

The package contract does **not** currently promise:

- every platform binary is shipped in-repo
- zero-toolchain bootstrap in all environments
- installer integration beyond the current Codex-oriented local bootstrap/install path

---

## Profile Contract

Repo policy must remain externalized from the adapter.

The adapter may assume:

- committed profiles live under `harness/profiles/`
- optional local override lives under `.harness/profile.local.yml`
- a bootstrapped repo receives a generated default profile

The adapter must **not** hardcode:

- project branch names
- project verification commands
- project release topology
- project-specific ownership conventions

That belongs to the profile contract.

---

## Governance Contract

This package is leader-led by design.

The contract assumes:

- a leader/control-plane viewpoint
- bounded teammate execution rather than freeform swarm behavior
- governed teammate-to-teammate communication
- councils for high-ambiguity or high-blast-radius decisions
- authority-aware reconciliation when state sources disagree
- explicit decision protocols, currently `evidence-review`, `incident-command`, and `decision-record`

The package intentionally follows ideas synthesized from:

- `teamwork` protocol practice
- `wing` skill/package restructuring
- engineering cybernetics
- Kubernetes control-plane / reconciliation thinking
- `arch-review`-influenced evidence-review governance rather than a single inherited meeting ritual

But the package contract is not identical to any one of those source systems.

---

## Compatibility and Versioning

The public compatibility anchor is the adapter identity plus adapter version in:

- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)

Compatibility expectations:

- patch/minor prototype evolution may extend the manifest or runtime behavior
- lifecycle surfaces should remain stable even if implementation details are replaced
- bootstrap/sync should remain manifest-driven
- old runtime state may require normalization; maintenance tooling is allowed to upgrade it forward

The package is still in `prototype` maturity.
That means stronger guarantees are attached to protocol shape and package boundaries than to final naming or frozen file layout.

---

## Explicit Non-Goals

This contract does not yet promise:

- final product name
- a permanently frozen manifest schema
- a universal multi-runtime story outside Codex
- a hosted orchestration service
- direct dependence on real Codex sub-agent APIs

Those decisions remain open productization work.

---

## What Future Work Must Preserve

Future implementation work may deepen:

- authority winner/loser rules
- debate routing and multi-round council orchestration
- ownership selector policy
- binary distribution and installer polish
- final naming and release positioning

But it should preserve:

- protocol / adapter / profile separation
- the leader-led control-plane model
- git-grounded state authority
- deterministic substrate outside prompt prose
- the ability to bootstrap the package into another repo without dragging historical `teamwork` or `wing` packaging debt back in
