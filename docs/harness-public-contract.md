# Codex-Native Harness Public Contract

> Status: Active
>
> Date: 2026-03-29
>
> Scope: define the durable public contract for the Codex-native orchestration harness package.

---

## Position

This package is the current public contract for the Codex-native orchestration harness stable line.

Its public package name is:

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

That manifest may evolve, but the stable package line should continue to present itself as `Harness` unless a later ADR explicitly supersedes that positioning.

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
- leader-facing diagnosis via `home` and `doctor`, with `home` mirroring scheduler decision-first guidance when bundle-governor decisions outrank simple foreground selection, surfacing scheduler maintenance backlog when stale state must be cleaned, and exposing the unified primary control-plane packet when governance repair outranks fresh foreground work
- leader-facing scheduler summaries and durable scheduler snapshots under `.harness/scheduler/snapshots/`, including mission classification, execution-mode recommendation, bundle-aware deferral/escalation, bundle-role tallies, bundle-scoped human-decision packets for escalation-required bundles, ranked `candidate_governor_ids` plus `candidate_governor_packets` when a bundle still needs an explicit governor, prioritized human-decision portfolios with `primary_required_human_decision` plus scope/surface tallies, mission-level queue visibility through `human_decision_queue_position` and `primary_human_decision`, explicit `defer_pending_primary_human_decision` actions for lower-priority governed decisions, decision-first foreground holds via `hold_for_primary_human_decision` when unresolved human decisions should block new low-urgency foreground work, sticky primary-decision retention via `decision_selection_policy` and `previous_primary_human_decision_key`, sticky foreground retention via `selection_policy` and `previous_foreground_mission_id`, sticky bounded-parallel retention via `previous_selected_mission_ids`, primary maintenance portfolios with `primary_maintenance_action`, deferred-maintenance count, `maintenance_selection_policy`, and `previous_primary_maintenance_key`, unified leader-facing control-plane priority through `primary_control_plane_kind`, `primary_control_plane_packet`, `primary_control_plane_selection_policy`, and deferred control-plane count, explicit maintenance-first decision deferral through mission-level `deferred_by_primary_control_plane_kind` and `defer_pending_primary_maintenance`, governance-first foreground holds via `hold_for_primary_maintenance` when critical maintenance should block low-urgency foreground work, bounded parallel admission from profile policy, explicit parallel-justification gating, P0 foreground serialization, profile-backed high-risk scope inference, profile-backed council recommendation from meeting triggers, explicit `recommended_decision_protocol` guidance for human decisions, fresh follow-up-signal rescheduling across normalized meeting/review/priority/risk/rebase signal kinds, projected `fresh_follow_up_signal_kinds` plus `follow_up_signal_urgency_rank` on follow-up decision packets, strongest-pending-signal selection when multiple fresh follow-up kinds coexist on one mission, pre-adoption `effective_risk_level: high` escalation for fresh `risk_raised` signals, direct leader adoption and clearing of explicit bundle governors, and next-action prioritization when human bundle decisions or governance maintenance are still pending
- structured bundle-governor council adoption, where `meetings lock-decision` can record a governing mission and `drive meeting-close` can emit the corresponding durable `bundle_governor_established` signal for bundle-scoped scheduler councils while superseding older explicit governors for the same bundle
- automated bundle-governor council locking, where auto-synthesis can derive a structured `governing_mission_id` lock from accepted findings on a bundle-scoped `establish_bundle_governor` council, allowing `meeting-close` to adopt durable bundle authority without a separate manual lock step when the winning mission is already clear
- provenance-bearing governed bundle state, where scheduler and mission-graph packets expose `governance_source_kind` and `adopted_from_meeting_id` when a durable governor came from a closed council rather than a direct leader-only explicit adoption
- leader-facing scheduler summaries and durable scheduler snapshots under `.harness/scheduler/snapshots/`, including mission classification, execution-mode recommendation, bundle-aware deferral/escalation, bundle-role tallies, bundle-scoped human-decision packets for escalation-required bundles, ranked `candidate_governor_ids` plus `candidate_governor_packets` when a bundle still needs an explicit governor, stable `decision_id` values on required human-decision packets, prioritized human-decision portfolios with `primary_required_human_decision` plus scope/surface tallies, mission-level queue visibility through `human_decision_queue_position`, `primary_human_decision`, and `human_decision_id`, explicit `defer_pending_primary_human_decision` actions for lower-priority governed decisions, decision-first foreground holds via `hold_for_primary_human_decision` when unresolved human decisions should block new low-urgency foreground work, sticky primary-decision retention via `decision_selection_policy` and `previous_primary_human_decision_key`, sticky foreground retention via `selection_policy` and `previous_foreground_mission_id`, sticky bounded-parallel retention via `previous_selected_mission_ids`, primary maintenance portfolios with `primary_maintenance_action`, deferred-maintenance count, `maintenance_selection_policy`, and `previous_primary_maintenance_key`, unified leader-facing control-plane priority through `primary_control_plane_kind`, `primary_control_plane_packet`, `primary_control_plane_selection_policy`, and deferred control-plane count, explicit maintenance-first decision deferral through mission-level `deferred_by_primary_control_plane_kind` and `defer_pending_primary_maintenance`, governance-first foreground holds via `hold_for_primary_maintenance` when critical maintenance should block low-urgency foreground work, bounded parallel admission from profile policy, explicit parallel-justification gating, P0 foreground serialization, profile-backed high-risk scope inference, profile-backed council recommendation from meeting triggers, explicit `recommended_decision_protocol` guidance for human decisions, fresh follow-up-signal rescheduling across normalized meeting/review/priority/risk/rebase signal kinds, projected `fresh_follow_up_signal_kinds` plus `follow_up_signal_urgency_rank` on follow-up decision packets, strongest-pending-signal selection when multiple fresh follow-up kinds coexist on one mission, pre-adoption `effective_risk_level: high` escalation for fresh `risk_raised` signals, direct leader adoption and clearing of explicit bundle governors, and next-action prioritization when human bundle decisions or governance maintenance are still pending
- decision-aware meeting materialization through `meetings materialize-primary-decision` and `meetings materialize-decision DECISION_ID`, including durable `source_decision_*` metadata, scheduler snapshot refs, and persisted source decision packets on the resulting meeting object
- idempotent decision materialization, where re-materializing an already-active `decision_id` reuses the live meeting instead of creating a duplicate council for the same scheduler object
- durable decision re-addressing, where an already-active scheduler `decision_id` remains materializable through the live meeting's persisted `source_decision_packet` even after the scheduler has promoted that packet into `active-council`
- bundle-scoped active-council portfolio retention, where a materialized bundle-governor council remains visible as the same bundle-scoped `decision_id` inside scheduler `required_human_decisions` instead of degrading into only mission-level review packets
- live council progress projection, where active-council decision packets now expose current meeting phase, protocol, next requirement, and debate counters instead of only preserving pre-materialization decision prompts
- current live council work-packet projection, where active-council decision packets and `meetings summary` expose the current work packet id/phase/action count and batch context, deriving that packet from either pending independent-review wave obligations or the active debate batch
- current live council action-packet projection, where those same leader-facing surfaces also expose the concrete current-wave action packets plus blocked debate-batch summaries so active and queued council work stays legible at the control-plane layer
- progress-aware active-council prioritization, where active-council decision packets expose `active_decision_progress_rank` and scheduler uses that rank to prioritize a more advanced live council ahead of a same-class peer that is still earlier in its governance flow
- bundle-governor meeting materialization that auto-seeds host plus candidate-governor participants, sibling mission refs, and candidate-governor summaries into the meeting packet when the source decision is `establish_bundle_governor`
- scheduler recognition of materialized decision packets as `active-council`, including stable `decision_id` values across surface changes and `active_meeting_ids` projected back into the portfolio packet
- bundle-aware recognition of materialized decision packets as live governance, including `bundle_active_decision_meeting_ids`, `bundle_active_decision_meeting_count`, `bundle_control_state`, `active_council_bundle_count`, and explicit `hold_for_active_bundle_council` reschedule actions when a bundle is already under active council control
- live active-council packets now drive leader-facing `next_action` and `home` summaries from the council's current `active_decision_next_requirement`, so in-flight governance exposes its current bounded ask instead of generic review wording
- active-council packets now also expose freshness through `active_decision_last_updated_at`, `active_decision_age_hours`, and `active_decision_stale`; stale in-flight councils are promoted into “refresh or close the stale active council” guidance instead of remaining passive bundle holds
- conflict detection between explicit bundle governors and active bundle councils, escalating mismatched anchors into `clear_invalid_bundle_governor` maintenance instead of silently trusting both governance sources at once
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
- protocol-aware live council guidance, where `incident-command` and `decision-record` now carry protocol-specific independent-review requirements and protocol-specific output requirements instead of reusing evidence-review wording once the meeting is already running
- protocol-aware council output packets, where synthesis/output can now materialize a structured provisional `output_packet` with lane-level synthesis plus structured `follow_up_packets`, and active-council leader surfaces expose both the top-level packet and its per-lane dispositions/follow-up before adoption or close
- closed-meeting output adoption, where `meeting-close` and `adopt-meeting` now write structured output packet metadata, follow-up packet counts, and the primary follow-up packet back into run state so later authority audit can detect `adoption_pending` drift beyond plain summary text
- leader-facing adopted-output projection, where scheduler mission packets and `home` now surface the adopted packet kind, adopted follow-up counts, and the primary adopted follow-up packet, letting the control plane keep a bounded execution object visible after a council closes
- durable adopted follow-up guidance, where an open primary adopted follow-up packet now also becomes the run's human-readable `next_expected_action`, and `runs authority-audit` treats missing promotion as adoption drift instead of assuming generic council-close wording is still good enough
- governing closed-meeting adoption anchors, where run-level adoption checks now prefer `adopted_from_meeting_id` as the governing closed-meeting reference for `latest_summary` and `next_expected_action`, preventing later non-governing closed councils from being mistaken for the run's durable adopted truth
- adopted follow-up lifecycle sync, where resolving an action on a closed adopted meeting now refreshes the governing run's adopted follow-up packets, primary adopted packet, and open-follow-up counts so leader-facing execution surfaces track real follow-up completion instead of a stale close-time snapshot
- teammate-native adopted follow-up handoff, where `teammates assign --adopted-follow-up` can bind the primary or named adopted follow-up packet from the governing run into a bounded teammate assignment, preserving packet id, execution kind, and governing meeting provenance while refusing duplicate active delegation of the same follow-up packet
- delegated teammate follow-up resolution, where a teammate run bound to an adopted follow-up packet can now call `drive resolve-follow-up` itself; runtime resolves the governing run's packet, updates the teammate into explicit post-follow-up reporting guidance, and keeps parent-run lifecycle advancement aligned with the same closure
- governing-to-delegated follow-up sync, where updating an adopted follow-up packet on the governing run now also refreshes any active delegated teammate run bound to that packet, including reopening a previously resolved delegated run when the governing packet returns to `open`, so delegated runtime state keeps following the governing packet truth even when the leader re-addresses the bounded work object directly
- delegated adopted-follow-up projection, where scheduler mission packets and `home` now expose delegated run id, delegated owner, delegated state, delegated freshness (`last_reported_at`, age, stale), and delegated monitor guidance for the primary adopted follow-up packet, making bounded teammate ownership visible without reopening the raw run record
- delegated-follow-up maintenance escalation, where a stale teammate run that still owns an open adopted follow-up packet now becomes `refresh_stale_delegated_follow_up` in the maintenance portfolio, allowing scheduler and `home` to prioritize bounded delegated-follow-up recovery over generic stale-run cleanup
- delegated-follow-up diagnostic alignment, where `doctor` now emits `stale-delegated-follow-up` for the same stale delegated-follow-up condition instead of collapsing it into generic stale-run drift
- delegated-follow-up authority-audit alignment, where `runs authority-audit` now emits a `delegated_follow_up_stale` adoption gap for that same stale delegated-follow-up condition, including delegated run and governing follow-up packet context so bounded drift stays explicit in run-level audit output
- delegated-follow-up home guidance gating, where `home.next_action` now follows delegated adopted follow-up wording only while the governing follow-up packet is still open, preventing a closed historical packet from continuing to own top-level execution guidance after resolution
- cancelled-follow-up realignment guidance, where a governing run now keeps explicit realignment wording when the last open adopted follow-up packet closes as `cancelled`, preventing cancelled bounded work from being treated as fully complete
- cancelled-follow-up scheduler escalation, where missions whose adopted follow-up packets are fully closed but include at least one `cancelled` packet now surface `decision_kind: realign_cancelled_follow_up` and `classification: needs-human-decision`, letting scheduler and `home` treat cancelled follow-up realignment as explicit leader work instead of only run-local guidance
- cancelled-follow-up audit alignment, where closed-meeting adopted-next-action semantics now expect that same realignment wording for fully cancelled follow-up packets, preventing `runs authority-audit` and `doctor` from falsely flagging cancelled realignment as missing adopted guidance
- completed-follow-up audit alignment, where closed-meeting adopted-next-action semantics now expect generic post-follow-up guidance once adopted follow-up packets exist but none remain open, preventing `runs authority-audit` and `doctor` from falsely flagging completed follow-up closure as missing adopted guidance
- open-only adopted-follow-up handoff, where `teammates assign --adopted-follow-up=<packet>` now rejects explicitly named adopted follow-up packets unless they are still `open`, preventing closed packet history from being rebound as new bounded teammate execution
- delegated completion gating, where a teammate run bound to an open adopted follow-up packet may not transition to `completed` until that governing packet is resolved, keeping leader-run and teammate-run closure semantics aligned around the same bounded work object
- run-native adopted follow-up completion, where `drive resolve-follow-up` can now resolve the primary or named adopted follow-up packet directly from the governing run, clearing the primary adopted packet when no open follow-up remains and advancing `next_expected_action` out of stale follow-up wording
- completion gating on adopted follow-up, where a run may not transition to `completed` while adopted follow-up packets remain open, ensuring that post-council execution obligations are cleared before the run can claim durable closure
- adopted-follow-up lifecycle ownership, where a closed meeting's open action items stop contributing to `open_meeting_action_count` after that meeting's structured output has been adopted into the governing run, preventing adopted follow-up work from being double-counted as lingering meeting backlog
- runtime-enforced council governance gates, so missing authority hierarchy and rejected findings without reasons block council progression instead of remaining doctor-only warnings
- authority audit across signals, meetings, runs, and live git state
- doctor-level validation that stale explicit bundle-governor signals are invalid governance rather than healthy active control-plane state
- mission-level authority audit across run aggregates, leader/teammate control boundaries, and multi-run mission conflicts
- graph-level mission authority can adjudicate a pairwise conflicted bundle into a governed bundle when one mission has a uniquely stronger explicit authority tuple
- graph-level mission authority audit across cross-mission ownership conflicts, dependency-governed execution boundaries, dependency-governed ownership adjudication, selector-specificity ownership adjudication, council-governed ownership adjudication, bundle-aware mission-graph grouping, and explicit bundle governors when a governed bundle has either a single root, one clear authority winner, one active council anchor, one leader-adopted explicit governor signal, or one uniquely strongest full candidate-governor tuple; if multiple active council anchors coexist on the same overlap bundle without a durable winner, runtime now escalates that state explicitly as `multi_council_bundle_governor_required`
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

The explicit release/freeze anchor for the first stable line is:

- [`docs/harness-release-criteria.md`](/Users/liyasong/Project/collov/erwin/docs/harness-release-criteria.md)

Compatibility expectations:

- patch/minor stable-line evolution may extend the manifest or runtime behavior compatibly
- lifecycle surfaces should remain stable even if implementation details are replaced
- bootstrap/sync should remain manifest-driven
- old runtime state may require normalization; maintenance tooling is allowed to upgrade it forward

The package is now in its first stable-line maturity.
That means protocol shape, package boundaries, and the committed release criteria are the compatibility anchor for subsequent evolution.

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
