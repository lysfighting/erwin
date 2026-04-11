# Harness Session Handoff

> Status: Active
>
> Date: 2026-04-10
>
> Purpose: provide a durable, file-backed starting point for fresh sessions so `Harness` work can resume from repository state instead of relying on chat-context memory.

---

## Why This Exists

`Harness` is now large enough that chat-context alone is not a safe source of truth.

Fresh sessions should prefer:

1. committed files
2. machine-readable package/runtime metadata
3. recent session logs
4. recent commits

and only then any surviving conversation context.

This is intentional state persistence, not incidental documentation.

External field intelligence and real usage feedback should enter `Harness` through the proposal-first learning loop and proposal ledger, not by directly rewriting runtime behavior from fresh observations alone.

When a fresh session needs the fastest human-readable orientation before diving into architecture and control-plane details, start with [`docs/harness-explainer.md`](/Users/liyasong/Project/collov/erwin/docs/harness-explainer.md).

Current human-direct entry stack:

- `bash scripts/harness-home.sh --brief`
- `bash scripts/harness-missions.sh schedule --brief`
- `bash scripts/harness-missions.sh graph --brief`
- `bash scripts/harness-missions.sh get MISSION_ID --brief`
- `bash scripts/harness-missions.sh authority-audit MISSION_ID --brief`
- `bash scripts/harness-missions.sh authority-audit --all --brief`
- `bash scripts/harness-runs.sh get RUN_ID --brief`
- `bash scripts/harness-runs.sh authority-audit RUN_ID --brief`
- `bash scripts/harness-signals.sh active [MISSION_ID] --brief`
- `bash scripts/harness-teammates.sh get TEAMMATE_RUN_ID --brief`
- `bash scripts/harness-meetings.sh summary MEETING_ID --brief`
- `bash scripts/harness-meetings.sh authority-audit MEETING_ID --brief`

---

## Current Identity

`Harness` is:

- the public package name
- a Codex-native adapter for the Autonomous Development Lifecycle Protocol
- a leader-led orchestration harness

It is **not**:

- a renamed `teamwork`
- a renamed `wing`
- a generic standalone cloud agent platform

The identity split is:

- `protocol`: Autonomous Development Lifecycle Protocol
- `adapter`: `Harness`
- `machine-readable adapter id`: `codex-native-harness`

Authoritative references:

- [`AGENTS.md`](/Users/liyasong/Project/collov/erwin/AGENTS.md)
- [`_project/CURRENT.md`](/Users/liyasong/Project/collov/erwin/_project/CURRENT.md)
- [`docs/harness-explainer.md`](/Users/liyasong/Project/collov/erwin/docs/harness-explainer.md)
- [`docs/harness-documentation-map.md`](/Users/liyasong/Project/collov/erwin/docs/harness-documentation-map.md)
- [`docs/harness-architecture.md`](/Users/liyasong/Project/collov/erwin/docs/harness-architecture.md)
- [`docs/harness-learning-loop.md`](/Users/liyasong/Project/collov/erwin/docs/harness-learning-loop.md)
- [`docs/harness-proposal-ledger.md`](/Users/liyasong/Project/collov/erwin/docs/harness-proposal-ledger.md)
- [`docs/harness-release-criteria.md`](/Users/liyasong/Project/collov/erwin/docs/harness-release-criteria.md)
- [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
- [`_project/decisions/0013-harness-prototype-positioning.md`](/Users/liyasong/Project/collov/erwin/_project/decisions/0013-harness-prototype-positioning.md)

---

## Non-Negotiable Architecture

Never collapse these layers back together:

1. `protocol`
   Lifecycle semantics. The core loop is still:
   `issue -> claim -> drive -> ship -> rc`

2. `adapter`
   Codex-facing runtime, shell surfaces, bootstrap/sync flow, manifest, validation, packaging.

3. `profile`
   Repo-specific policy. Current proving-ground profile:
   [`harness/profiles/erwin.yml`](/Users/liyasong/Project/collov/erwin/harness/profiles/erwin.yml)

Also keep this boundary explicit:

- [`harness/`](/Users/liyasong/Project/collov/erwin/harness/README.md) = committed definition layer
- `.harness/` = local runtime state

That boundary is architectural, not accidental.

---

## Runtime Reality

The runtime is now Go-based.

Execution path:

- [`cmd/harness-state/main.go`](/Users/liyasong/Project/collov/erwin/cmd/harness-state/main.go)
- [`internal/harness/runtime.go`](/Users/liyasong/Project/collov/erwin/internal/harness/runtime.go)

Shell entrypoints are intentionally thin:

- [`scripts/harness-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-runtime.sh)
- [`scripts/harness-claim.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-claim.sh)
- [`scripts/harness-drive.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-drive.sh)
- [`scripts/harness-teammates.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-teammates.sh)
- [`scripts/harness-meetings.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-meetings.sh)
- [`scripts/harness-ship.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-ship.sh)
- [`scripts/harness-rc.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-rc.sh)

Human-facing runtime note:

- `home` and `doctor` now support `--brief` for leader-readable summaries derived directly from the same control-plane truth, and `missions schedule --brief`, `missions graph --brief`, `missions get MISSION_ID --brief`, `missions authority-audit MISSION_ID --brief`, `missions authority-audit --all --brief`, `runs get RUN_ID --brief`, `runs authority-audit RUN_ID --brief`, `signals active [MISSION_ID] --brief`, `teammates get TEAMMATE_RUN_ID --brief`, `meetings summary MEETING_ID --brief`, and `meetings authority-audit MEETING_ID --brief` now do the same across scheduler, mission, graph, run, signal, teammate, and council object layers; the default output remains machine-readable JSON for audits, tooling, and exact state inspection

Runtime selection order:

1. `harness/bin/<os>-<arch>/harness-state`
2. `.harness/bin/harness-state`
3. local `go build`

Meaning:

- normal harness usage no longer depends on `python3`
- packaged binary is the preferred end-user path
- local Go toolchain is a fallback, not the desired steady state
- adapter validation now also treats `public_contract_path` as a required machine-readable package anchor and recognizes `release_criteria_path` as the machine-addressable `1.0.0` freeze bar

Current scheduler/council reality:

- active bundle-governor councils now remain durable bundle-scoped decision objects in the portfolio
- scheduler now normalizes fresh signal kinds before rescheduling, so aliases such as `meeting_required` land on the same control-plane path as `requires_meeting`
- fresh `review_feedback`, `priority_change`, `risk_raised`, and `requires_rebase` signals now become real scheduler work instead of passive telemetry, pausing foreground execution or reclassifying the mission until the signal is adopted
- follow-up decision packets now project `fresh_follow_up_signal_kinds` and `follow_up_signal_urgency_rank`, and a fresh `risk_raised` signal lifts `effective_risk_level` to `high` before adoption so urgent follow-up work can outrank lower-urgency peers in the scheduler portfolio
- when one mission carries multiple fresh follow-up signals at once, scheduler now lets the strongest pending signal govern both urgency ranking and leader guidance instead of following the first lower-urgency kind that happened to arrive
- scheduler and `home` follow an in-flight council's live `next_requirement` instead of generic review wording
- scheduler active-council packets and `meetings summary` now also project the current live council work packet, deriving it from either pending independent-review wave obligations or the current debate batch so leader-facing surfaces can see the current in-flight ask more concretely
- those same leader-facing surfaces now also project the concrete current-wave action packets plus blocked debate-batch summaries, so the leader can see both the exact open council work items and which bounded batches are still queued behind the current packet
- once a council reaches synthesis/output, orchestration now also materializes a protocol-aware `output_packet` with lane-level synthesis plus structured `follow_up_packets`, active-council packets plus `meetings summary` keep that provisional top-level plus per-lane output artifact visible until the meeting is adopted or closed, and `meeting-close` / `adopt-meeting` now write the adopted output packet plus primary follow-up packet back into run state
- scheduler mission packets plus `home` now project that adopted output state back into leader-facing execution, including adopted packet kind, follow-up counts, and the primary adopted follow-up packet, and they now let a selected mission's primary adopted follow-up packet drive `next_action` when no higher-priority control-plane work outranks it
- run adoption now also promotes an open primary adopted follow-up packet into human-readable `next_expected_action`, and `runs authority-audit` expects that promoted follow-up guidance instead of generic council-close wording once the meeting is already closed
- run-level adoption validation now also prefers `adopted_from_meeting_id` as the governing closed-meeting anchor for summary and next-action checks, so doctor and authority-audit do not compare run state against a newer but non-governing closed council on the same run
- when a closed adopted meeting later resolves one of its own action items, runtime now refreshes the governing run's adopted follow-up packet state as well, so scheduler/home adopted follow-up counts and primary packet selection move with actual completion instead of remaining frozen at adoption time
- the leader can now advance that same lifecycle directly through `drive resolve-follow-up`, which resolves the primary or named adopted follow-up packet on the governing run and clears stale follow-up guidance when no open adopted packet remains
- that same adopted follow-up object can now also hand off directly into teammate runtime through `teammates assign --adopted-follow-up`, which derives the delegated objective from the follow-up packet, records governing meeting / execution-kind provenance on the teammate run, and rejects duplicate active handoff of the same packet
- once that handoff exists, the delegated teammate can now also close the same adopted follow-up directly through `drive resolve-follow-up` on its own run; runtime routes the closure back into the governing leader run and advances both sides out of stale follow-up handling
- when the governing leader run re-addresses that same adopted follow-up packet directly, runtime now also resyncs the delegated teammate's `delegated_follow_up_state` from the governing packet truth, including reopening a previously resolved delegated run when the parent packet returns to `open`, so delegated execution does not drift just because the parent touched the bounded object first
- scheduler mission packets plus `home` now also project that delegated adopted-follow-up ownership back into the control plane, including delegated run id, delegated owner, delegated state, delegated freshness, and delegated monitor guidance, so leader-facing surfaces stop talking as though the governing run still owns that task directly and can escalate stale delegated follow-up handling before it silently drifts
- once that delegated follow-up owner is actually stale, scheduler now also elevates the condition into a dedicated critical maintenance action (`refresh_stale_delegated_follow_up`), so fresh sessions should expect top-level `next_action` to follow bounded delegated-follow-up recovery instead of generic stale-run cleanup
- `doctor` now names that same condition `stale-delegated-follow-up`, so fresh sessions should treat stale delegated adopted follow-up as a first-class bounded-control issue rather than generic teammate staleness
- `runs authority-audit` now names that same condition `delegated_follow_up_stale` inside adoption gaps, so fresh sessions can inspect governing run audit output and still see which delegated run plus follow-up packet is blocking trustworthy bounded execution
- once that delegated adopted follow-up has actually been resolved, `home` no longer lets the closed delegated follow-up wording own top-level `next_action`, so fresh sessions should treat resolved delegated packets as visible history rather than still-live control-plane guidance
- if that last open adopted follow-up closes as `cancelled`, the governing run now stays on explicit realignment guidance, so fresh sessions should treat cancelled follow-up as a bounded correction path rather than the same thing as completed work
- scheduler now also promotes that fully cancelled adopted follow-up state into explicit decision work (`realign_cancelled_follow_up`), so fresh sessions should expect leader-facing control-plane packets and `home` to treat cancelled follow-up realignment as first-class human-decision work when no higher-priority surface outranks it
- that decision portfolio is now also kind-aware, so fresh sessions should expect `realign_cancelled_follow_up` to stay ahead of ordinary `adopt_meeting_actions` work when both are open in the same mixed human-decision portfolio
- closed-meeting adopted-next-action semantics now follow that same cancelled realignment path, so fresh sessions should not expect `doctor` or run-level authority audit to complain when the governing run is correctly holding cancelled-follow-up realignment guidance
- closed-meeting adopted-next-action semantics now also follow generic post-follow-up guidance once no adopted follow-up packet remains open, so fresh sessions should not expect `doctor` or run-level authority audit to complain after a correctly completed follow-up lifecycle either
- `teammates assign --adopted-follow-up=<packet>` now only accepts still-open adopted follow-up packets, so fresh sessions should not expect closed `done` or `cancelled` packet history to be delegatable back into bounded teammate execution
- that delegated teammate still cannot be marked `completed` until the same governing adopted follow-up packet is closed, so fresh sessions should treat delegated completion as gated by the same bounded follow-up lifecycle rather than as an independent shortcut
- that same run still cannot be marked `completed` until those adopted follow-up packets are cleared, so fresh sessions should treat open adopted follow-up as a real completion gate rather than optional residue on an otherwise-finished run
- once a closed meeting's structured output has already been adopted, its remaining open action items no longer count as generic meeting backlog in scheduler/home surfaces, so leader-facing `open_meeting_action_count` reflects only action items that still belong to live or not-yet-adopted meeting control
- bundle authority now distinguishes one active council anchor from several: one active anchor can still govern an overlap bundle, but multiple active anchors on the same bundle now escalate explicitly into `multi_council_bundle_governor_required` so the leader consolidates governance instead of trusting split live councils
- stale active councils are now promoted into explicit “refresh or close the stale active council” guidance instead of remaining passive hold state
- scheduler now also exposes one leader-facing `primary_control_plane_packet`, so critical maintenance can deliberately outrank lower-priority human decisions without hiding those deferred decisions from the portfolio
- when that happens, affected governed decision missions now also carry explicit maintenance-first defer state instead of only silently staying visible in the queue

---

## Machine-Readable Truth

The main machine-readable package boundary is:

- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)

Current important facts:

- `adapter_version`: `1.0.0`
- `release_criteria_path`: `docs/harness-release-criteria.md`
- `public_name`: `Harness`
- `adapter_id`: `codex-native-harness`
- `public_contract_path`: [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)

This manifest is the best place to confirm:

- package identity
- package file list
- bootstrap contract
- profile contract
- capability surfaces
- smoke coverage

---

## Most Important Runtime Areas

When resuming implementation, the most important regions are in:

- [`internal/harness/runtime.go`](/Users/liyasong/Project/collov/erwin/internal/harness/runtime.go)

High-value areas:

- scheduler summary / snapshot generation / execution-mode selection / rescheduling hints
- selector / overlap / transfer model
- run authority audit
- mission-level and mission-graph authority
- meeting governance / wave plan / route plan / round plan / auto-synthesis / authority audit
- CLI subcommands

Do not treat `Harness` as “a bunch of shell scripts.”
The shell layer is intentionally thin; the Go runtime is the real control plane.

---

## What The System Already Does

Do not underestimate the current control-plane maturity.

`Harness` already has:

- leader-facing scheduler baseline with mission classification, execution-mode selection, bundle-aware deferral/escalation, bundle-role tallies, decision-first next-action prioritization, rescheduling hints, and durable `.harness/scheduler/snapshots/`
- `home` now mirrors that decision-first scheduler guidance instead of collapsing escalation-required bundles back into generic cross-mission conflict messaging
- scheduler now consumes profile-backed bounded parallel capacity, so overflow runnable missions can be deferred behind the current admitted batch instead of always widening all runnable work
- scheduler now consumes the execution policy `require_explicit_parallel_justification`, so parallel widening only happens when the current admitted batch carries explicit `parallel_justification`
- scheduler now enforces `P0` foreground serialization, so a runnable `P0` mission prevents widening into a parallel batch even when other isolated missions are also runnable
- scheduler now consumes profile-backed high-risk scopes, so configured boundaries such as `install.sh` or `skills/wing/` can escalate a mission into high-risk foreground serialization even when no manual `risk_level` was provided
- scheduler now consumes profile-backed meeting triggers, so some required human decisions become explicitly council-worthy with durable `recommended_decision_surface` and `meeting_trigger_refs` instead of only generic decision records
- escalation-required mission bundles now collapse into one bundle-scoped scheduler decision packet, so the leader sees one governing-bundle decision instead of one duplicated decision per mission
- scheduler now exposes a prioritized human-decision portfolio, including a `primary_required_human_decision`, deferred decision count, and scope/surface tallies, so leader-facing surfaces can reason about more than one pending governed decision at once
- scheduler human-decision packets now carry stable `decision_id` values, and mission classifications plus defer actions now project those ids back into runtime state instead of leaving the portfolio as anonymous packet blobs
- those stable decision packets can now be materialized directly into active meetings through `meetings materialize-primary-decision` or `meetings materialize-decision DECISION_ID`, and the resulting meeting persists `source_decision_*` metadata plus scheduler snapshot linkage for later audit
- re-materializing a still-active `decision_id` now reuses the existing governed meeting rather than silently creating a duplicate active council for the same scheduler object
- once a decision has already promoted into `active-council`, that same `decision_id` still remains re-addressable because runtime can recover the durable source packet from the active meeting itself instead of depending on the latest scheduler snapshot to keep the original packet shape alive forever
- when that active decision is bundle-scoped, scheduler now keeps it visible as the same bundle-scoped portfolio packet instead of letting it degrade into only per-mission `review_active_council` noise after materialization
- active-council packets now also reflect the live meeting's progress, so portfolio decisions can expose current phase and next requirement rather than only replaying the stale pre-materialization action text
- materialized bundle-governor councils now auto-seed the host plus candidate-governor participants, sibling mission refs, and candidate tuple summaries into the meeting packet, so the leader starts from bounded scheduler evidence instead of a blank participant roster
- once materialized, the same decision id now survives the shift from `council` to `active-council`, and scheduler projects `active_meeting_ids` back into the packet so the portfolio can see that the decision is already being governed live
- once a bundle-scoped decision is materialized, scheduler now also projects that live council back into every affected mission through `bundle_active_decision_meeting_*` plus `bundle_control_state`, counts `active_council_bundle_count`, and can emit `hold_for_active_bundle_council` instead of treating the bundle as generic unresolved governance
- if a leader later adopts an explicit governor that conflicts with the live bundle-council host, scheduler now escalates that mismatch into `clear_invalid_bundle_governor` maintenance instead of allowing the bundle to carry split governance anchors
- escalation-required bundle decisions now carry ranked `candidate_governor_ids` and `candidate_governor_packets`, so the leader sees plausible governors and their authority tuple before deciding whether to establish a governor directly or open a council
- leaders can now explicitly establish and clear a bundle governor through mission surfaces, turning an escalation-required bundle into governed authority deliberately instead of leaving that decision as a perpetual open packet
- bundle-scoped `establish_bundle_governor` councils can now close with a structured locked decision and automatically adopt that verdict into durable explicit bundle authority, so council outcomes do not stop at archived meeting notes
- those same bundle-governor councils can now auto-generate the structured governing lock from accepted council findings when the winner is already clear, reducing the common path from `materialize -> manual lock -> close` to `materialize -> orchestrate -> close`
- that adopted governor now keeps provenance in graph/scheduler state through `governance_source_kind` and `adopted_from_meeting_id`, so later sessions can tell whether a governed bundle is anchored by a direct leader decision or by a closed council adoption
- stale explicit bundle-governor signals are now surfaced as scheduler maintenance and doctor warnings, so orphaned bundle governance is treated as cleanup work instead of a silently trusted control-plane fact
- larger conflict bundles can now adjudicate from the fuller candidate-governor tuple rather than only raw priority/risk, so a unique release-adjacent or governance-active candidate can govern even when the simpler tuple still ties
- scheduler can now emit `hold_for_primary_human_decision`, refusing to start new low-urgency foreground work when an unresolved primary human decision should be cleared first and no existing foreground batch is being continued
- mission classifications now mirror that human-decision queue through `human_decision_queue_position` and `primary_human_decision`, so leader-facing mission surfaces can see which pending decision is currently first-class
- lower-priority governed decisions now also emit explicit `defer_pending_primary_human_decision` reschedule actions, so the queue influences runtime behavior instead of remaining a read-only annotation
- scheduler now retains a sticky primary human-decision packet across otherwise-tied portfolio passes, surfacing `decision_selection_policy` and `previous_primary_human_decision_key` so the leader can see when the control plane is intentionally holding decision focus steady
- when multiple same-class active councils coexist, scheduler now also uses `active_decision_progress_rank` to prioritize the more advanced live council instead of falling back to freshness-only ordering
- scheduler now retains a sticky foreground mission across otherwise-tied scheduling passes, surfacing `selection_policy` and `previous_foreground_mission_id` so the leader can see when the control plane is intentionally preferring stability over freshness-only churn
- scheduler now retains a sticky bounded parallel batch across otherwise-tied portfolio passes, surfacing `previous_selected_mission_ids` so the leader can see when the admitted set is being preserved deliberately instead of recomputed from raw freshness
- scheduler now sorts maintenance into a first-class portfolio, surfacing `primary_maintenance_action`, deferred-maintenance count, `maintenance_selection_policy`, and `previous_primary_maintenance_key` instead of leaving maintenance as an unordered side list
- governance-critical maintenance such as stale explicit bundle governors can now emit `hold_for_primary_maintenance`, deferring low-urgency foreground work until the primary maintenance action is repaired or explicitly deferred
- meeting objects now carry explicit `decision_protocol` and `decision_protocol_reason`, so the control plane no longer needs to pretend every council is implicitly one inherited `arch-review` debate
- scheduler human-decision packets now carry `recommended_decision_protocol`, letting the leader distinguish evidence-review, incident-command, and decision-record handling before materializing a meeting
- protocol-aware council automation is now partially real: `incident-command` and `decision-record` no longer auto-expand every same-authority conflict into route/round work, instead deferring bounded disagreement into synthesis/record output while `evidence-review` keeps the richest live debate controller
- protocol-aware live council guidance is also now real: once those meetings are in flight, `incident-command` and `decision-record` surface protocol-specific independent-review and output requirements instead of collapsing back to generic council wording
- debate-batch control is now stateful: when multiple open batches exist, the controller exposes one deterministic current batch, keeps the rest blocked, and now retains that current batch when a new same-phase lane appears instead of reordering on freshness alone
- debate-batch identity is now stable and hashed, so distinct batches cannot collide through truncated slugs and accidentally present multiple current batches
- scheduler now treats fresh coordination/dependency/blocker/review/rebase signals as real rescheduling triggers instead of only including them in the snapshot payload
- scheduler now escalates stale runs, expired signals, and stale meetings into explicit `maintenance_actions` instead of leaving them only to `doctor`
- `home` now mirrors that maintenance backlog and the current primary maintenance packet so leader-facing surfaces stay aligned with scheduler reality
- first-class mission objects
- separation between `latest_run_id` and `authoritative_run_id`
- graph-level authority audit with bundle grouping, explicit bundle governors, explicit-governor signal adoption, and pairwise governor adjudication from clear authority tuples
- ownership selectors across `file` / `dir` / `module` / `glob` / `scope` / `interface`
- path-family overlap plus same-boundary and hierarchical semantic overlap
- ownership transfer modes:
  - `move`
  - `copy`
  - `split`
- exclusion-aware parent/child delegation for split transfer
- role-bounded, authority-aware, claim/dimension-aware meetings
- wave / route / round obligations in councils
- deterministic current/blocked debate-batch execution state
- same-pass closure from synthesis to output when a council record is already ready
- direct synthesis when a council lane is aligned or already resolved by higher-authority evidence
- auto-synthesis of provisional findings
- latest-response verdict adoption when a governed debate lane resolves
- runtime-enforced council governance gaps, so missing authority hierarchy and rejected findings without reasons block council progression until repaired
- leader/teammate split with governed sibling communication

This is no longer “just a prototype shell workflow.”

Important interpretation rule:

- `arch-review` is still an important influence on bounded evidence-review councils
- it is no longer the sole semantic template for all governed decisions in `Harness`

---

## What Still Matters Most

Two statements are both true:

1. `Harness` already delivers the first stable line.
2. The ultimate control-plane end state is still not fully complete.

These do not conflict.

The first means:

- runtime exists
- adapter exists
- packaging exists
- public contract exists
- major control-plane primitives exist

The second means higher-order control-plane semantics are still open.

The remaining high-value gaps are:

- fuller scheduler / rescheduling / execution-mode-selection behavior beyond the current baseline snapshot-and-hint control plane
- broader authority adjudication for larger mission bundles when no obvious governing dependency exists
- fuller `arch-review`-style automated debate execution beyond today’s route/wave/round/verdict-adoption baseline
- optional future mapping between Codex subagents and teammate runtime

---

## How To Read Research Correctly

These two audits remain important:

- [`docs/research/harness-gap-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-gap-audit-2026-03-29.md)
- [`docs/research/harness-research-drift-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-research-drift-audit-2026-03-29.md)

Use them for:

- control-plane maturity
- absorbed vs partial vs missing semantics
- avoiding design drift

Do **not** rely on them as the final truth for:

- naming
- distribution
- latest product-positioning details

Those areas were partially overtaken by later 2026-03-29 sessions and by:

- [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)
- recent session logs

---

## Fresh Session Start Order

For a new session that needs enough context to work safely, start here:

1. [`AGENTS.md`](/Users/liyasong/Project/collov/erwin/AGENTS.md)
2. [`_project/CURRENT.md`](/Users/liyasong/Project/collov/erwin/_project/CURRENT.md)
3. [`docs/harness-explainer.md`](/Users/liyasong/Project/collov/erwin/docs/harness-explainer.md)
4. [`docs/harness-documentation-map.md`](/Users/liyasong/Project/collov/erwin/docs/harness-documentation-map.md)
5. [`docs/harness-architecture.md`](/Users/liyasong/Project/collov/erwin/docs/harness-architecture.md)
6. [`docs/harness-operating-model.md`](/Users/liyasong/Project/collov/erwin/docs/harness-operating-model.md)
7. [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
8. [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)
9. this file
10. [`docs/research/harness-gap-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-gap-audit-2026-03-29.md)
11. [`docs/research/harness-research-drift-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-research-drift-audit-2026-03-29.md)
12. recent session logs
13. recent Harness commits on `main` or the active task branch

Recommended recent session logs:

- [`_project/sessions/2026-03-29-AL.md`](/Users/liyasong/Project/collov/erwin/_project/sessions/2026-03-29-AL.md)
- [`_project/sessions/2026-03-29-AM.md`](/Users/liyasong/Project/collov/erwin/_project/sessions/2026-03-29-AM.md)

---

## Fresh Session Rules

A fresh session should avoid these mistakes:

- do not reframe `Harness` as a `teamwork`/`wing` rename
- do not flatten `protocol / adapter / profile`
- do not treat shell wrappers as the real implementation layer
- do not confuse “first stable line” with “control-plane end state fully solved”
- do not trust stale chat summaries over committed files

The repository should be treated as the persistent memory substrate.

That is the point of this file.
