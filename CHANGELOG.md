# Changelog

All notable changes to the **erwin** project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added — Harness Stable Line
- **Harness v1.0.0 — Stale delegated follow-up now escalates into explicit monitor guidance**: scheduler mission packets and `home` now project delegated adopted follow-up freshness through delegated `last_reported_at`, age, and stale fields, and they upgrade the delegated monitor action when that teammate run goes stale. Teammate smoke now proves stale delegated follow-up ownership surfaces as stale and that a fresh teammate report returns the monitor path to normal.
- **Harness v1.0.0 — Leader-facing control plane now shows delegated follow-up state**: scheduler mission packets and `home` now expose the delegated adopted follow-up state alongside delegated run id, owner, and monitor guidance, so the leader can see whether the delegated bounded execution unit is currently `open` without reopening the teammate run directly. Teammate smoke now proves both surfaces project that state during handoff.
- **Harness v1.0.0 — Governing reopen now reopens delegated follow-up runs too**: if a delegated teammate already resolved an adopted follow-up packet and the governing run later reopens that packet, runtime now reopens the delegated teammate's bounded follow-up state and restores its follow-up guidance instead of leaving the delegated run stuck on `done`. Teammate smoke now proves the full resolve -> governing reopen -> delegated reopen -> resolve-again lifecycle.
- **Harness v1.0.0 — Governing follow-up refresh now syncs delegated teammate state**: when the governing run updates an adopted follow-up packet directly, runtime now also refreshes any active teammate run bound to that packet instead of leaving delegated follow-up state stale until the teammate itself resolves it. Teammate smoke now proves the leader can re-address a delegated adopted follow-up packet and the delegated run immediately reflects the refreshed open state.
- **Harness v1.0.0 — Delegated adopted follow-up now gates teammate completion**: a teammate run bound through `teammates assign --adopted-follow-up` can no longer transition to `completed` while the governing adopted follow-up packet is still open. After `drive resolve-follow-up` closes that packet, the delegated teammate run can complete normally, keeping completion semantics aligned across leader and delegated execution units.
- **Harness v1.0.0 — Delegated adopted follow-up now projects back into scheduler and home**: once an adopted follow-up packet is handed off to a teammate run, scheduler mission packets and `home` now expose the delegated run id, delegated owner, and delegated monitor guidance. Selected-mission guidance therefore tells the leader to monitor the delegated execution unit rather than replaying the raw follow-up text as though it were still directly owned by the governing run.
- **Harness v1.0.0 — Delegated teammates can now resolve adopted follow-up directly**: `drive resolve-follow-up` now also works on a teammate run that was created from `teammates assign --adopted-follow-up`, routing the packet resolution back into the governing run while updating the delegated teammate into explicit post-follow-up reporting guidance. Scheduler/home open-follow-up counts still clear from the governing run, and the parent run still advances into generic post-follow-up execution guidance.
- **Harness v1.0.0 — Adopted follow-up can now hand off directly into teammate runtime**: `teammates assign --adopted-follow-up` can now bind a governing run's primary or named adopted follow-up packet directly into a teammate assignment, deriving the delegated objective from the packet, preserving governing meeting plus execution-kind provenance, and refusing duplicate active handoff of the same follow-up packet. Follow-up resolution still advances the parent run out of adopted-follow-up handling even after that delegated monitoring handoff.
- **Harness v1.0.0 — Run completion now waits for adopted follow-up closure**: runs can no longer transition to `completed` while adopted follow-up packets remain open. The leader must first clear those packets, for example through `drive resolve-follow-up`, before the run can close with durable completion semantics.
- **Harness v1.0.0 — Adopted follow-up can now resolve directly from drive**: introduced `drive resolve-follow-up`, letting the governing run resolve the primary or named adopted follow-up packet directly. When the last open adopted follow-up closes, the runtime clears the primary adopted packet, reduces adopted open-follow-up counts in scheduler/home, and advances `next_expected_action` out of stale follow-up wording.
- **Harness v1.0.0 — Adopted follow-up lifecycle now syncs after action resolution**: resolving an action on a closed adopted meeting now refreshes the governing run's adopted follow-up packets, primary adopted packet, and open-follow-up counts. Scheduler and `home` therefore reduce adopted open-follow-up counts when the underlying follow-up is actually completed instead of staying pinned to the original `meeting-close` snapshot.
- **Harness v1.0.0 — Adoption audit now anchors to the governing closed meeting**: run-level adoption checks now prefer `adopted_from_meeting_id` as the governing closed-meeting reference for `latest_summary` and `next_expected_action` instead of blindly comparing against the newest closed meeting on the run. That keeps doctor and authority-audit aligned with adopted follow-up guidance and avoids false `run-missing-meeting-*` warnings when later non-governing councils also exist.
- **Harness v1.0.0 — Adopted council actions now leave meeting backlog cleanly**: once a closed council already has structured output and that output has been adopted into the governing run, its remaining open action items no longer contribute to `open_meeting_action_count`. Scheduler and `home` now treat that work as adopted follow-up execution instead of lingering generic meeting backlog, and core smoke proves the backlog clears immediately at `meeting-close`.
- **Harness v1.0.0 — Adopted council follow-up now becomes human-readable run guidance**: when a closed council already carries an open primary follow-up packet, `meeting-close` and `adopt-meeting` now promote that packet's human-readable summary into `next_expected_action` instead of leaving the run on generic council-close wording or an internal follow-up code. `runs authority-audit` still expects that promoted guidance as part of adoption completeness, and teammate smoke proves the readable run-guidance path end to end.
- **Harness v1.0.0 — Adopted council output now drives leader-facing execution surfaces**: scheduler mission packets and `home` now surface adopted output packet kinds, follow-up packet counts, and the primary adopted follow-up packet from run state instead of leaving that structure buried in raw run JSON. When no higher-priority decision or maintenance work outranks the selected mission, scheduler and `home` now let the primary adopted follow-up packet drive the next action. Teammate smoke now proves both `missions schedule` and `home` preserve protocol-aware adopted follow-up packet execution semantics after `meeting-close`.
- **Harness v1.0.0 — Closed council output now adopts back into run state**: `meeting-close` and `adopt-meeting` now copy structured `output_packet` metadata, follow-up packet counts, and the primary follow-up packet back into run state instead of only copying summary text. `runs authority-audit` now treats missing adopted output metadata as explicit `adoption_pending` drift, teammate smoke proves protocol-aware follow-up packet adoption, and meeting ids now include sub-second entropy so one mission can materialize multiple meetings in the same second without collision.
- **Harness v1.0.0 — Council output packets now synthesize structured follow-up packets**: provisional council `output_packet` artifacts now carry `follow_up_packets` derived from open action items and unresolved/open lane state, while `meetings summary` and active-council scheduler packets project follow-up packet counts and payloads back into leader-facing surfaces. Teammate smoke now proves protocol-specific incident and decision-record follow-up packets, and core smoke proves active-council output projection keeps follow-up packet counts visible.
- **Harness v1.0.0 — First stable line frozen and fully verified**: promoted the adapter line from the prototype series to `1.0.0`, set `release_stage: stable`, and completed the full committed verification matrix (`go test`, adapter validate, issue/core/release/teammate/runtime-dist/runtime-package/bootstrap smoke, and doctor) against the stable package boundary. Remaining deeper scheduler, bundle, and council work is now treated as post-1.0 expansion rather than a blocker for the stable line.
- **Harness v1.0.0 — Council output packets now preserve lane-level synthesis**: provisional council `output_packet` artifacts now carry synthesized lane packets for debated and directly synthesized claim lanes, while `meetings summary` and active-council scheduler packets project those lane counts and packets back into leader-facing surfaces. Teammate smoke now proves accepted versus unresolved multi-lane output explicitly, and core smoke proves scheduler projection of the lane-aware active-council output artifact.
- **Harness v1.0.0 — Multi-council overlap now escalates explicitly**: mission-graph authority now treats multiple active council anchors on the same overlap bundle as a distinct escalation state instead of generic conflict, emitting `multi_council_bundle_governor_required` plus a consolidation-first leader action. Core smoke now proves the bundle flips from single-anchor governed authority into explicit multi-council escalation as soon as the second overlapping council goes active.

### Added — Harness Learning Loop
- **Proposal-first Harness learning loop added**: introduced `docs/harness-learning-loop.md` as the canonical design for turning field signals, real usage cases, and internal harness gaps into reviewed proposals instead of direct runtime churn.
- **Harness proposal ledger added**: introduced `docs/harness-proposal-ledger.md` as the durable control-plane surface for reviewed `Harness` evolution proposals.

### Added — Repo Repositioning
- **AI/agent radar foundation added**: created `docs/radar/README.md`, `docs/radar/sources.md`, `docs/radar/themes.md`, and the kickoff weekly synthesis at `docs/radar/weekly/2026-03-30.md`, establishing a lightweight repo-native knowledge base for authoritative field signals.
- **Archive boundary added**: created `archive/README.md`, `docs/archive/README.md`, and `_project/archive/README.md` so historical material can be demoted safely before any delete round.
- **Radar automations activated**: the daily radar and weekly synthesis automations are now stored as active Codex automations for the `erwin` workspace, scheduled for weekday mornings and Friday evening in the local timezone.

### Changed — Harness Evolution Discipline
- **Harness evolution now has a documented proposal gate**: `docs/README.md`, `_project/CURRENT.md`, and `docs/harness-session-handoff.md` now point fresh sessions toward the learning loop and proposal ledger, making it explicit that external field intelligence and real usage should enter `Harness` through reviewed proposals rather than direct implementation.


### Changed — Repo Structure
- **erwin repositioned around Radar + Harness + Skills**: added a new top-level `README.md`, rewrote `AGENTS.md` and `CLAUDE.md`, and updated `_project/GOALS.md` plus `_project/CURRENT.md` so the repository no longer presents itself as a showcase-first skill incubator.
- **Session history moved behind an archive boundary**: `_project/sessions` now resolves to `_project/archive/sessions` through a compatibility symlink, reducing default surface noise while preserving older links and skill behavior.
- **Low-value historical material was cut aggressively**: archive now keeps only the small set of session logs still referenced by the active repo surface, and `teamwork-demo/` has been removed entirely because it no longer provided implementation or reference value.
- **Low-value historical docs were deleted instead of preserved indefinitely**: removed `showcase-v3-architecture.md`, `teamwork-architecture.md`, `teamwork-v2-architecture.md`, `visual-enhancement-architecture.md`, and `project-management-architecture.md` after confirming they had no remaining practical reference value beyond legacy changelog text.
- **Showcase was fully retired**: removed the last live harness profile and smoke-test references to `showcase/`, then deleted the directory because it no longer provided active implementation value.

### Added — Harness Prototype
- **Harness v0.95.0-prototype — Councils now auto-synthesize protocol-aware output packets**: orchestration now generates durable `output_packet` artifacts from council state instead of stopping at findings/action items, `meetings summary` and active-council scheduler packets now expose those packets, and meeting outcome adoption now prefers the generated output summary/next action when present. Teammate smoke now proves `incident-command` and `decision-record` output packets, and core smoke proves scheduler projection of the active evidence-review output packet.
- **Harness v0.94.0-prototype — 1.0.0 release criteria are now committed and packaged**: introduced `docs/harness-release-criteria.md` as the canonical promotion/freeze bar for the first stable adapter line, taught `harness/adapter-manifest.yml` to declare `release_criteria_path`, taught `adapter validate` to verify both the public contract and release criteria anchors, and widened runtime bundles plus package smoke so release criteria ship with the binary artifact.
- **Harness v0.93.0-prototype — Active councils now project live action packets and blocked batches**: scheduler active-council packets and `meetings summary` now expose the concrete current-wave action packets for the live council work packet, and debate-phase summaries now also expose blocked batch summaries so leader-facing surfaces can see both the active work items and the queued bounded batches without reopening raw meeting state. Core and teammate smoke now prove both the action-packet projection and blocked-batch visibility paths.
- **Harness v0.92.0-prototype — Active councils now project their current live work packet**: scheduler active-council packets and `meetings summary` now expose the current live council work packet instead of only a generic `next_requirement`, including current packet id/phase/action count and batch context. Runtime now derives that current work packet from either pending independent-review wave obligations or the active debate batch, and smoke coverage now proves both the scheduler-facing and meeting-summary-facing paths.
- **Harness v0.91.0-prototype — Strongest-signal follow-up urgency now wins on mixed-signal missions**: scheduler now computes `follow_up_signal_urgency_rank` from the strongest pending fresh follow-up signal instead of the first matching kind, and mission-level follow-up guidance now follows that same strongest pending signal. Core smoke now proves a mission carrying both fresh `priority_change` and fresh `risk_raised` still exposes urgency rank `4` and risk-raised guidance.
- **Harness v0.90.0-prototype — Fresh follow-up urgency is now first-class portfolio state**: scheduler now projects `fresh_follow_up_signal_kinds` and `follow_up_signal_urgency_rank` onto mission classifications and required human-decision packets, raises `effective_risk_level` to `high` when a fresh `risk_raised` signal is still pending adoption, and prioritizes urgent risk-raised follow-up work ahead of lower-urgency fresh `priority_change` peers in the human-decision portfolio. Core smoke now proves the urgency ranking, portfolio ordering, and risk escalation path end to end.
- **Harness smoke coverage now explicitly proves priority/risk fresh-signal follow-up**: core smoke now covers fresh `priority_change` and `risk_raised` signals directly, proving they route into governed follow-up decision work with signal-specific scheduler guidance instead of being covered only implicitly through the shared follow-up path.
- **Harness v0.89.0-prototype — Fresh signal rescheduling now covers review and rebase control-plane work**: runtime now normalizes signal aliases such as `meeting_required -> requires_meeting`, scheduler treats fresh `review_feedback` as governed follow-up decision work, and scheduler treats fresh `requires_rebase` as explicit blocking/reschedule work with a dedicated `pause_for_fresh_rebase_signal` action. Core smoke now proves alias normalization plus review-feedback and rebase-driven rescheduling end to end.
- **Harness v0.88.0-prototype — Scheduler now prioritizes more advanced live councils**: active-council decision packets now expose `active_decision_progress_rank`, and scheduler uses that rank to prioritize the more advanced live council when multiple same-class active councils coexist. Core smoke now proves a synthesis-phase council outranks a fresher setup-phase peer in the human-decision portfolio.
- **Harness v0.87.0-prototype — Maintenance-first decision deferral is now explicit**: when governance-critical maintenance wins the unified control-plane portfolio, scheduler now marks affected decision-bearing missions with `deferred_by_primary_control_plane_kind: maintenance`, records the winning maintenance action, and emits explicit `defer_pending_primary_maintenance` reschedule actions instead of only leaving those decisions visible behind the primary packet. Core smoke now proves maintenance-first packet selection and maintenance-first decision deferral together.
- **Harness v0.86.0-prototype — Unified leader-facing control-plane priority**: scheduler no longer leaves the leader to merge separate human-decision and maintenance queues by hand. It now exposes `primary_control_plane_kind`, `primary_control_plane_packet`, `primary_control_plane_selection_policy`, and deferred control-plane count, and it lets governance-critical maintenance such as `clear_invalid_bundle_governor` outrank lower-priority human decisions without hiding those deferred decisions from the portfolio. Core smoke now proves the unified priority surface and `home` mirrors it directly.
- **Harness v0.85.0-prototype — Council protocols now shape live execution guidance**: `incident-command` and `decision-record` no longer only change route planning; orchestration now emits protocol-specific independent-review requirements and protocol-specific output guidance, so in-flight councils ask for incident assessments or decision-record tradeoff capture instead of generic evidence-review wording. Teammate smoke now proves both protocol paths surface distinct live guidance.
- **Harness v0.84.0-prototype — Stale active councils now become explicit scheduler work**: active-council decision packets now expose `active_decision_last_updated_at`, `active_decision_age_hours`, and `active_decision_stale`. When a live bundle council goes stale under cleanup policy, scheduler and `home` now elevate “refresh or close the stale active council” into the primary next action instead of silently continuing bundle holds. Core smoke now proves stale active bundle councils become explicit leader-facing work.
- **Harness v0.83.0-prototype — Scheduler and home now follow live council requirements**: once a bundle-scoped decision is already under active council control, scheduler `next_action`, `primary_required_human_decision.recommended_action`, `hold_for_active_bundle_council.reason`, and `home.next_action` now follow the meeting's live `next_requirement` rather than generic review wording. Core smoke now proves leader-facing surfaces mirror live bundle-council requirements directly.
- **Harness v0.82.0-prototype — Scheduler now projects live council progress**: active-council decision packets now mirror the live meeting's phase, protocol, next requirement, and debate counters back into the scheduler portfolio, so leader-facing decision surfaces can follow governed council progress without dropping into raw meeting state. Core smoke now proves an active bundle council's phase and next requirement appear directly in scheduler output.
- **Harness v0.81.0-prototype — Active bundle councils now stay bundle-scoped in the portfolio**: once a bundle-governor decision has been materialized into a live council, scheduler now recovers that active bundle packet back into `required_human_decisions`, keeping the same bundle-scoped `decision_id` visible in the portfolio instead of degrading into only mission-level `review_active_council` packets. Core smoke now proves active bundle councils remain visible as bundle-scoped control-plane objects after materialization.
- **Harness v0.80.0-prototype — Active decision ids now survive active-council rollover**: once a scheduler decision has already materialized into a live council, `meetings materialize-decision DECISION_ID` now falls back to the active meeting's persisted `source_decision_packet`, refreshes that governed meeting from current scheduler context, and reuses the same live meeting instead of failing after the decision packet leaves the latest snapshot's pre-materialization shape. Core smoke now proves bundle-scoped active-council decisions remain durably re-addressable.
- **Harness v0.79.0-prototype — Scheduler decision materialization now reuses active councils**: materializing a scheduler `decision_id` now reuses the existing active meeting when that governed decision is already under live council control, instead of spawning duplicate active councils for the same control-plane object. Core smoke now proves repeated materialization keeps the same meeting id and preserves `active-council` packet stability.
- **Harness v0.78.0-prototype — Materialized bundle councils now start with seeded participants and candidate context**: bundle-scoped `establish_bundle_governor` scheduler decisions now materialize into meetings that auto-seed the host plus candidate-governor participants, carry sibling mission refs into packet context, and summarize ranked `candidate_governor_packets` in `known_context`. Core smoke now proves the council starts with bounded scheduler-derived bundle context before orchestration begins.
- **Harness v0.77.0-prototype — Bundle-governor councils now auto-lock durable authority**: bundle-scoped `establish_bundle_governor` councils no longer need a manual `meetings lock-decision` in the happy path. Auto-synthesis now preserves `source_run_id` on generated findings, derives one structured `governing_mission_id` lock when accepted evidence clearly favors a candidate mission, and lets `drive meeting-close` adopt that governor directly into durable bundle authority. Core smoke now proves the full materialized-council -> auto-synthesis -> auto-lock -> close/adopt path.
- **Harness v0.76.0-prototype — Bundle authority now keeps council-adoption provenance**: mission-graph authority and scheduler packets now preserve whether a durable governor came from a closed council adoption via `governance_source_kind: council-adopted-governor` plus `adopted_from_meeting_id`. Core smoke now proves governed bundle state keeps that provenance after council close, so leader-facing control surfaces can distinguish direct explicit adoption from council-adopted authority.
- **Harness v0.75.0-prototype — Closed bundle councils now adopt durable bundle authority**: `meetings lock-decision` can now record a structured `governing_mission_id`, and `drive meeting-close` now consumes that verdict for bundle-scoped `establish_bundle_governor` councils. Closing the council emits a durable `bundle_governor_established` signal, supersedes any prior active governor signals for the same bundle, and core smoke now proves the bundle becomes governed immediately after council close without polluting later runtime health.
- **Harness v0.74.0-prototype — Explicit governors must now agree with live bundle councils**: when a `bundle_governor_established` signal points at a different mission than the active bundle-council host, scheduler now escalates that mismatch into `clear_invalid_bundle_governor` maintenance. Core smoke proves the conflicting-governor path and its cleanup, so bundle governance can no longer silently split across two incompatible anchors.
- **Harness v0.73.0-prototype — Active bundle councils now become first-class scheduler state**: once a bundle-scoped decision is materialized into a live meeting, scheduler now projects that council back into affected missions through `bundle_active_decision_meeting_*` and `bundle_control_state`, counts `active_council_bundle_count` at the portfolio level, and emits `hold_for_active_bundle_council` so live bundle governance is treated as active control-plane state rather than generic overlap noise.
- **Harness v0.72.0-prototype — Materialized decisions now stay visible as active councils**: scheduler decision ids are now stable across surface changes, so once a human-decision packet is materialized into an active meeting the portfolio can promote it into `active-council` state without minting a new identity. Scheduler snapshots now project `active_meeting_ids` back onto the decision packet, which tightens the bridge between the decision portfolio and live council governance.
- **Harness v0.71.0-prototype — Scheduler decisions now materialize into durable councils**: leader surfaces can now convert the current primary human-decision packet or a specific stable `decision_id` into an active meeting via `meetings materialize-primary-decision` and `meetings materialize-decision`. Materialized meetings persist `source_decision_*` metadata, scheduler snapshot refs, and the source decision packet itself, so the scheduler portfolio is now bridged directly into durable council state instead of remaining a read-only queue.
- **Harness v0.70.0-prototype — Scheduler decision packets now have stable ids**: required human-decision packets now carry durable `decision_id` values, and mission classifications plus `defer_pending_primary_human_decision` actions now project those ids back into runtime state. This turns the decision portfolio into a more formal control-plane surface instead of anonymous JSON and prepares a cleaner bridge into future meeting materialization. Core smoke now proves stable decision ids on primary packets, mission classifications, and defer actions.
- **Harness v0.69.0-prototype — Bundle adjudication now consumes the fuller candidate-governor tuple**: larger mission bundles no longer stop at `priority` / `risk_level` when deciding whether a unique governor already exists. Graph authority now consumes the fuller candidate-governor tuple, including release-adjacent and governance-activity signals, so a uniquely stronger bundle candidate can govern even when the simpler tuple still ties. Core smoke now proves both the unresolved equal-tuple bundle and the newly adjudicated release-adjacent winner path.
- **Harness v0.68.0-prototype — Council controller now retains sticky current debate batches**: `meetingDebatePlan` now preserves the previous open `current_wave_batch_id` when a new same-phase lane appears, exposing `batch_selection_policy` and `previous_current_wave_batch_id` so councils keep one deterministic active batch instead of thrashing to the latest inserted lane. Debate-batch ids now use stable hashed identities, preventing truncated-slug collisions that could accidentally mark multiple batches as current. Teammate smoke now proves the sticky-batch path and the repaired blocked-batch counts.
- **Harness v0.67.0-prototype — Maintenance is now a first-class scheduler portfolio**: scheduler now sorts maintenance work into a real portfolio, exposing `primary_maintenance_action`, deferred-maintenance count, `maintenance_selection_policy`, and `previous_primary_maintenance_key`. Governance-critical maintenance such as stale explicit bundle governors can now emit `hold_for_primary_maintenance`, deferring low-urgency foreground execution until the primary maintenance task is repaired or explicitly deferred. `home` now mirrors the current primary maintenance packet, and core smoke proves stale bundle governance both surfaces as maintenance and blocks fresh foreground admission.
- **Harness v0.66.0-prototype — Scheduler now retains sticky bounded parallel batches**: when a bounded parallel batch is still valid and a fresher same-class peer appears, scheduler now keeps the admitted batch instead of churning the portfolio on freshness alone. Snapshots and `home` now expose `previous_selected_mission_ids`, while overflow peers still surface `defer_due_parallel_capacity` so admission pressure remains visible. Core smoke now proves a fresher peer cannot evict the established admitted batch.
- **Harness v0.65.0-prototype — Stale explicit governors now escalate into maintenance**: if a leader-adopted `bundle_governor_established` signal no longer points at a live active bundle or its governing mission has fallen out of that bundle, scheduler now emits `clear_invalid_bundle_governor` maintenance work and `doctor` reports `invalid-explicit-bundle-governor`. Core smoke now proves stale explicit governance is surfaced and then clears cleanly back to an escalation-required surviving bundle.
- **Harness v0.64.0-prototype — Leaders can now explicitly adopt bundle governors**: `missions establish-governor` and `missions clear-governor` now let the leader turn an escalation-required mission bundle into governed authority via durable `bundle_governor_established` signals. Mission-graph authority consumes that explicit governor, scheduler stops repeating `establish_bundle_governor` for the adopted bundle, and core smoke proves the bundle can move into and back out of explicit governance cleanly.
- **Harness v0.63.0-prototype — Council orchestration is now protocol-aware**: `incident-command` and `decision-record` meetings no longer auto-expand every same-authority conflict into live route/round obligations. Instead they defer bounded disagreement into synthesis/record output, expose `protocol_deferred_route_count` plus `protocol_bound_route_count`, and auto-synthesize protocol-specific follow-up actions while `evidence-review` retains the richest debate controller. Teammate smoke now proves both protocol paths can reach output without opening cross-examination work for same-authority disagreement.
- **Harness v0.62.0-prototype — Scheduler now enforces decision-first foreground admission**: when a primary governed human decision remains unresolved and no established foreground batch is being continued, scheduler now refuses to start new low-urgency foreground work, emitting `hold_for_primary_human_decision` instead of quietly widening execution. Core smoke now proves the runtime pauses an older foreground, surfaces the unresolved decision backlog, withholds a fresh peer mission, and then restores the baseline foreground once the backlog clears.
- **Harness v0.61.0-prototype — Escalation bundles now expose candidate governors**: mission-graph authority, scheduler decision packets, and mission classifications now carry ranked `candidate_governor_ids` plus richer `candidate_governor_packets` when a bundle still needs an explicit governor. This gives the leader concrete governor candidates and their current authority tuple before choosing a direct governance action or opening a council. Core smoke now proves those packets survive both graph authority output and scheduler portfolio output.
- **Harness v0.60.0-prototype — Decision protocols are now first-class control-plane metadata**: meetings now carry `decision_protocol` and `decision_protocol_reason`, materialized meetings seed protocol-specific roles, role specs, and expected output, and scheduler human-decision packets now expose `recommended_decision_protocol` so the leader can distinguish `evidence-review`, `incident-command`, and `decision-record` handling before materializing a meeting. Core and teammate smoke now prove default evidence-review materialization, explicit incident-command overrides, manual protocol retargeting, and scheduler protocol guidance.
- **Harness v0.59.0-prototype — Scheduler retains sticky primary human decisions**: scheduler now remembers the previous `primary_required_human_decision` and exposes `decision_selection_policy` plus `previous_primary_human_decision_key`, allowing the current primary decision packet to stay first when a same-class peer arrives later. `home` now mirrors this decision-selection metadata, and core smoke proves a first mission-scoped decision remains primary when a later mission-scoped peer arrives with the same decision class.
- **Harness v0.58.0-prototype — Scheduler retains sticky foreground selection**: scheduler snapshots now remember `previous_foreground_mission_id` and expose `selection_policy`, allowing the current foreground mission to stay selected when a newly arrived runnable mission is otherwise tied on priority, effective risk, and release posture. Core smoke proves the control plane keeps the original foreground mission instead of thrashing to the fresher timestamp while the parallel-justification gate is still serializing the batch.
- **Harness v0.57.0-prototype — Lower-priority human decisions now defer behind the primary queue slot**: scheduler now emits `defer_pending_primary_human_decision` actions when a mission-scoped governed decision sits behind a higher-priority bundle- or mission-scoped decision in the same portfolio. Core smoke proves the secondary decision is both annotated with queue position 2 and actively deferred behind the primary packet.
- **Harness v0.56.0-prototype — Mission classifications now mirror decision-queue order**: scheduler now projects its prioritized human-decision portfolio back into mission classifications through `human_decision_queue_position` and `primary_human_decision`. Core smoke proves escalation-bundle missions are marked as queue position 1 while a lower-priority mission-scoped decision is marked as queue position 2.
- **Harness v0.55.0-prototype — Scheduler exposes a prioritized human-decision portfolio**: scheduler now sorts pending human decisions into a real decision queue, exposing `primary_required_human_decision`, deferred-decision count, and scope/surface tallies. `home` now mirrors that portfolio metadata, and core smoke proves a bundle-scoped council decision outranks a concurrent mission-scoped follow-up decision under the current profile policy.
- **Harness v0.54.0-prototype — Scheduler collapses bundle decisions into bundle-scoped packets**: escalation-required mission bundles no longer emit one repeated `establish_bundle_governor` decision per mission. Scheduler now aggregates those cases into a single bundle-scoped human-decision packet with merged mission ids, trigger refs, and bundle-required action, and core smoke proves `council_decision_count` stays aligned with the real number of bundle decisions the leader must make.
- **Harness v0.53.0-prototype — Home mirrors scheduler maintenance backlog**: `home` now surfaces `scheduler_maintenance_action_count` and `scheduler_maintenance_action_tally`, so the leader-facing surface no longer hides maintenance work that already exists in scheduler snapshots. Core smoke now proves stale-state maintenance guidance is visible in both `schedule` and `home`.
- **Harness v0.52.0-prototype — Scheduler escalates stale control-plane state into maintenance work**: scheduler snapshots now emit `maintenance_actions`, `maintenance_action_count`, and `maintenance_action_tally`, turning stale runs, expired signals, and stale meetings into explicit leader-facing maintenance work instead of leaving them as doctor-only observations. Core smoke now proves stale run plus stale signal maintenance guidance and shows the scheduler keeps the foreground batch bounded while stale state is cleaned.
- **Harness v0.51.0-prototype — Council governance gates are now runtime-enforced**: `meetings orchestrate` and `meetings summary` now surface explicit `governance_gaps`, and councils no longer advance across setup/synthesis boundaries while required governance policy is violated. Missing authority hierarchy and rejected findings without reasons are now control-plane blockers rather than doctor-only warnings, and teammate smoke proves both the blocked and repaired paths.
- **Harness v0.50.0-prototype — Scheduler requires explicit parallel justification**: scheduler now consumes the execution policy `require_explicit_parallel_justification`, turning it from dormant policy metadata into a real admission gate. Parallel widening only occurs when the current admitted mission batch carries explicit `parallel_justification`; otherwise the scheduler falls back to serialized execution and emits `defer_due_missing_parallel_justification` actions. Core smoke now proves both the justified parallel path and the missing-justification fallback path.
- **Harness v0.49.0-prototype — Scheduler enforces P0 serialization**: scheduler now turns the control-plane rule “P0 missions never run in parallel” into real admission behavior. When a runnable `P0` mission exists, the foreground batch is serialized under `serialize_p0_priority`, peers are deferred with `defer_until_p0_clears`, and core smoke now proves that a `P0` mission suppresses otherwise-isolated runnable peers.
- **Harness v0.48.0-prototype — Scheduler recommends councils from profile meeting triggers**: scheduler now consumes profile `meeting_triggers`, deriving `recommended_decision_surface`, `meeting_trigger_refs`, and `council_decision_count` for human decisions that should escalate into council handling. `home` now mirrors that council-oriented scheduler guidance, and core smoke proves an escalation-required mission bundle is surfaced as council-worthy rather than generic human decision noise.
- **Harness v0.47.0-prototype — Scheduler infers profile-backed high-risk foreground work**: scheduler now consumes profile `high_risk_scopes`, matching them against mission/run ownership selectors plus release-adjacent semantics to derive `profile_high_risk`, `profile_high_risk_scope_refs`, and `effective_risk_level`. Missions touching configured high-risk boundaries such as `install.sh` now serialize the foreground batch even without an explicit manual `risk_level`, and core smoke now proves the low-risk peer is deferred behind the inferred high-risk foreground mission.
- **Harness v0.46.0-prototype — Scheduler reacts to fresh follow-up signals**: scheduler now distinguishes fresh vs stale follow-up signals and treats fresh `requires_meeting`, `interface_change`, `dependency_discovered`, and `blocked` signals as real rescheduling triggers. Fresh coordination signals now escalate missions into explicit human decisions, while fresh dependency/blocker signals can pause a foreground run until the signal is adopted into runtime state. Core smoke now proves a fresh coordination signal reclassifies a mission and emits a reschedule action before clearing cleanly.
- **Harness v0.45.0-prototype — Scheduler enforces bounded parallel admission**: scheduler now obeys profile-backed `max_parallel_runs` instead of widening every runnable mission into the same parallel batch. Overflow runnable missions are deferred with explicit `defer_due_parallel_capacity` reschedule actions, and snapshots now record `batch_admission_policy`, `max_parallel_runs`, and `capacity_deferred_*` metadata. Core smoke now proves a third isolated runnable mission is deferred behind the configured capacity ceiling.
- **Harness v0.44.0-prototype — Home mirrors scheduler bundle guidance**: `home` now prioritizes scheduler decision-first guidance when escalation-required bundle decisions outrank generic conflict messaging, and it surfaces `scheduler_bundle_role_tally` so the leader sees the same bundle portfolio shape in both places. Core smoke now proves the parity between `schedule` and `home`.
- **Harness v0.43.0-prototype — Scheduler prioritizes bundle decisions**: scheduler snapshots now emit `bundle_role_tally`, and `next_action` now prioritizes required human bundle decisions ahead of foreground batch messaging when both coexist. Core smoke now proves escalation-required bundles keep the leader focused on establishing a governor even while other bounded work remains runnable.
- **Harness v0.42.0-prototype — Scheduler consumes bundle authority**: `missions schedule` now consumes mission-bundle authority packets directly instead of only the raw conflict list. Governed losing missions defer behind the current bundle governor, escalation-required bundles surface explicit `establish_bundle_governor` decisions, and scheduler snapshots now expose bundle metadata such as `bundle_role`, `bundle_resolution_kind`, and `bundle_governing_mission_id`. Core smoke now proves both the defer-behind-governor and escalation-to-decision paths.
- **Harness v0.41.0-prototype — Pairwise bundle governor adjudication**: mission-graph authority can now elect a governing mission for a two-mission conflicted bundle when one root mission has a uniquely stronger priority/risk authority tuple, instead of reserving governor selection for only larger multi-mission components. Core smoke now proves a pairwise overlap becomes `bundle_governor_adjudicated` without dependency, council, or selector-specificity help.
- **Harness v0.40.0-prototype — Stateful debate batch execution**: `meetingDebatePlan` now computes deterministic `current_wave_batch_id`, `ready_wave_batch_count`, `blocked_wave_batch_count`, and per-batch `state/current` metadata instead of treating debate batches as flat open buckets. `meetings orchestrate` now refuses to jump from `synthesis` to `output` while a debate batch still remains open, `meetings summary` surfaces the new controller state, and `doctor` warns when no deterministic current batch exists or synthesis/output is reached with open debate batches. Teammate smoke now proves one ready plus one blocked batch and verifies that a manually pushed `synthesis` phase still cannot bypass the active batch.
- **Harness v0.39.0-prototype — Direct council output closure**: `meetings orchestrate` now advances straight from a newly generated synthesis record into `output` whenever the synthesis audit is already ready, so multi-lane councils, governed verdict lanes, authority-resolved lanes, and aligned lanes no longer need an extra orchestration pass just to leave `synthesis`. Teammate smoke now proves those one-pass council closures end-to-end.
- **Harness v0.38.0-prototype — Multi-lane debate batching**: `meetingDebatePlan` now groups open route and round obligations into per-participant, per-round `wave_batches`, and `meetings summary` surfaces `wave_batch_count`, `multi_lane_batch_count`, and `max_wave_lane_count`. Teammate smoke now proves one bounded challenger can receive a two-lane debate batch in the same council round.
- **Council-governed bundle authority**: `missions authority-audit --all` now treats a single active council anchor as a real ownership adjudication source, emitting `council_governed_mission_ownership` plus `council_governed_bundle` so overlapping missions can become governed even without an explicit dependency edge when one mission is already chairing the relevant council.
- **Harness v0.37.0-prototype — Selector-specific bundle authority**: `missions authority-audit --all` can now adjudicate some no-dependency mission overlaps by applying the authority rule `specific over vague` to ownership selectors. Narrower selectors can win against broader overlapping selectors, yielding `selector_specificity_governed_mission_ownership`, `single_winner_governed_bundle`, and a leader-facing `required_action` instead of collapsing every such case into `bundle_conflict`. Core smoke now proves `dir:docs/release` loses to `file:docs/release/CHANGELOG.md` while later scheduler/home checks still return to the original clean baseline.
- **Harness v0.36.0-prototype — Leader-facing scheduler baseline**: added `missions schedule` plus durable `.harness/scheduler/snapshots/` output, mission classification, execution-mode recommendation (`solo`, `managed_delegation`, `scheduled_parallel_missions`), reschedule hints, required-human-decision reporting, and `home` scheduler summary fields. Core and teammate smoke now prove solo scheduling, dependency-blocked rescheduling, governed parallel batches, and managed delegation.
- **Council semantics synced into canonical docs**: updated `docs/harness-architecture.md`, `docs/harness-operating-model.md`, and `docs/harness-public-contract.md` so they now explicitly reflect that cross-examination is the normal path for unresolved contested lanes, while aligned or already authority-settled lanes may go directly into synthesis with an explicit governed record.
- **Council direct-synthesis control**: `meetings orchestrate` can now skip empty debate when independent review is aligned, and can also materialize a governed `higher_authority_prevails` verdict when one side already outranks the other before cross-examination. `meetings summary` now surfaces the skipped-debate counts, and teammate smoke proves both direct-synthesis paths end-to-end.
- **Harness documentation system**: added `docs/harness-documentation-map.md` to define the canonical Harness doc set and read paths, plus `docs/harness-operating-model.md` to capture the leader/teammate/council execution model as a first-class durable design artifact instead of leaving that core behavior scattered across architecture prose, research, and chat history.
- **Single architecture source of truth**: added `docs/harness-architecture.md` as the complete design overview for `Harness`, unifying identity, lineage, preserved protocol semantics, intentionally dropped historical baggage, control-plane model, Go runtime rationale, current control-plane maturity, and remaining expansion gaps in one durable architecture document.
- **Durable session handoff doc**: added `docs/harness-session-handoff.md` as the active file-first entrypoint for fresh sessions, explicitly telling future work to resume from committed state, manifest metadata, recent session logs, and recent commits instead of relying on chat memory.
- **Debate verdict adoption**: auto-synthesis can now honor a governed `latest_response_prevails` route verdict, so a debated claim lane can materialize directly into an accepted or rejected finding instead of falling back to unresolved output when the latest bounded reply has already won.
- **Split ownership transfer**: `runs transfer-ownership` now supports `--mode=split`, so a run can keep a broader parent selector while delegating a governed child selector through explicit ownership exclusions instead of fake overlap or full handoff.
- **Hierarchical semantic selector overlap**: ownership governance now treats parent/child semantic boundaries like `scope:release` and `scope:release/notes` as the same governed family, including cross-kind `scope`/`interface` combinations. Teammate smoke now proves an exclusive `interface:release` assignment is rejected when `scope:release/notes` already owns the child boundary.
- **Bundle-level authority packets**: `missions authority-audit --all` now emits bundle-level `resolution_kind`, `selector_refs`, and `required_action`, so connected mission bundles explain not just that they exist, but why they are conflicted/governed and what the leader should do next.
- **Semantic selector-family overlap**: ownership governance now treats same-boundary `scope:` and `interface:` selectors as overlapping semantic claims instead of unrelated exact strings, and teammate smoke now proves that an exclusive `scope:release-contract` assignment is rejected when `interface:release-contract` is already active.
- **Bounded multi-round council controller**: `meetings orchestrate` now plans round 2, round 3, and later bounded replies from policy instead of stopping at a second-round special case; equal-authority debate lanes can keep exchanging governed responses until a fresher/stronger reply prevails or the debate-round cap is reached.
- **Multi-mission authority adjudication**: `missions authority-audit --all` now surfaces cross-mission ownership collisions and dependency-governed mission blocking as first-class control-plane verdicts, and `doctor` now reports `mission-graph-authority-conflict` when active missions collide.
- **Mission-level authority adjudication**: missions now track both `latest_run_id` and `authoritative_run_id`, and `missions authority-audit` explains when leader-facing mission control should stay anchored to a leader/root run even as fresher teammate activity lands.
- **Runtime bundle packaging**: added `scripts/harness-package-runtime.sh` plus `scripts/harness-runtime-package-smoke-test.sh`, so `Harness` can now emit per-target runtime tarballs that carry the binary, adapter manifest, and public contract together.
- **Harness prototype positioning ADR**: accepted ADR `0013-harness-prototype-positioning.md`, making `Harness` the public prototype name for the Codex-native adapter line while keeping protocol identity separate from package identity.
- **Artifact-family authority adjudication**: `runs authority-audit` now emits explicit artifact winner/loser records across contract snapshots, active run state, and live git state, so stale planning artifacts lose explicitly instead of only surfacing as generic drift.
- **Bounded debate-round control**: `meetings orchestrate` now plans governed follow-up debate obligations when a linked cross-examination does not yet outrank the opposing claim lane, keeping councils in `cross_examination` until a bounded multi-round controller says the lane has resolved or hit the round cap.
- **Public package contract**: Added `docs/harness-public-contract.md` as the durable public boundary for the Codex-native harness package, defining protocol/adapter/profile separation, required surfaces, runtime guarantees, distribution expectations, compatibility, and non-goals while final naming remains open.
- **Automated council synthesis**: `meetings orchestrate` can now derive provisional findings and unresolved follow-up actions from grouped council positions without overwriting manual findings already on the record.
- **Claim-aware council positions**: `meetings add-position` now accepts `--claim=` and `--dimension=` so debate inputs can be grouped into the same governed claim during synthesis.
- **Meeting summary automation metrics**: `meetings summary` now reports auto-generated finding/action counts alongside traced-finding and authority-conflict reporting.
- **Live intent authority detection**: `runs authority-audit` now surfaces dirty live worktrees that contradict governed verification or handoff intent, and worktree snapshots ignore the harness-owned `.mission` marker so authority drift tracks real execution state.
- **Debate-route governance**: meetings now materialize routed vs unrouted conflicting claim lanes, `meetings summary` reports that routing state, and `doctor` flags councils that reached advanced phases without routing their conflicting claims into linked cross-examination.
- **Debate route obligations**: meetings now assign bounded challenger roles/runs to each conflicting claim lane, track resolved vs open route obligations in `meetings summary`, and let `doctor` flag councils that still have unfulfilled debate routes in advanced phases.
- **Missing-worktree authority detection**: `runs authority-audit` now treats a missing provisioned worktree as first-class live-state drift, and the core smoke verifies that reprovision can restore the governed worktree back to health.
- **Glob ownership selectors**: the runtime now understands `glob:` ownership selectors, and teammate smoke verifies that glob-pattern claims participate in semantic overlap protection alongside `file:` and `dir:` selectors.
- **Debate wave obligations**: `meetings orchestrate` now materializes missing independent-review obligations, and teammate smoke verifies that councils surface and then clear those wave duties before moving into cross-examination.
- **Ship-gate authority winner/loser rules**: `runs authority-audit` now treats failing or non-approved verification/review/check gate state as stronger evidence than an inconsistent optimistic `ship_status`, and release smoke verifies that those gate claims win until the run is brought back into agreement.
- **Codex-native orchestration harness prototype**: Added repository-local `harness/` definition layer and `.harness/` runtime model with profile, templates, and policies. Implemented deterministic substrate in `scripts/harness_state.py` plus shell surfaces for `home`, `doctor`, `claim`, `drive`, `meetings`, and `maintenance`.
- **Leader control plane features**: file-backed runtime locking, run-state transition validation, meeting materialization and archival, meeting outcome adoption back into runs, reprovisioning of existing runs into branch/worktree modes, reconciliation and normalization utilities.
- **Codex-native skill surface**: Added `skills/harness/SKILL.md` and `skills/harness/README.md` as the first repository-local skill for operating the harness.
- **Codex installer + smoke test**: Added `install-codex.sh` for installing Codex-native `SKILL.md` skills into `~/.codex/skills/`, plus `scripts/harness-smoke-test.sh` for isolated end-to-end verification of the harness prototype.
- **Packaging verification**: Validated direct `./install-codex.sh` execution, made the installer and smoke-test script executable, and installed `/harness` into the global Codex skill directory.
- **Bootstrap adapter flow**: Added `scripts/harness-bootstrap.sh` plus `harness/bootstrap/profile-template.yml` so the harness substrate, definitions, generated default profile, local `skills/harness/`, and `install-codex.sh` can be seeded into another git repository.
- **Bootstrap smoke coverage**: Added `scripts/harness-bootstrap-smoke-test.sh` to verify bootstrap, target-repo leader surfaces, target-repo `install-codex.sh`, and overwrite safety boundaries end-to-end.
- **Adapter sync path**: Added `.harness/cache/adapter-source.json` metadata and `scripts/harness-sync.sh` so bootstrapped repositories can resync from a remembered or explicit source repository.
- **Release lifecycle surfaces**: Added `scripts/harness-ship.sh` and `scripts/harness-rc.sh` plus RC runtime artifacts so the harness now covers verification, PR state, merge state, RC creation, RC verification, promotion, and closure.
- **Release smoke coverage**: Added `scripts/harness-release-smoke-test.sh` to validate the `claim -> ship -> rc` path end-to-end.
- **Issue surface**: Added `scripts/harness-issue.sh` plus `harness/templates/issue.md` so the harness can draft or create structured executable issues before claim.
- **Issue smoke coverage**: Added `scripts/harness-issue-smoke-test.sh` for offline validation of structured issue drafting.
- **Protocol coverage milestone**: The Codex-native harness now has first-class surfaces for the full lifecycle: `issue -> claim -> drive -> ship -> rc`.
- **Ship delivery mechanics**: Added deterministic branch push and PR draft generation to the `ship` surface so delivery is less state-only and more action-oriented.
- **GitHub PR opening**: Added a `ship pr-open` path that uses `gh pr create` to open real GitHub PRs when authenticated.
- **Post-merge cleanup**: Added a `ship cleanup` path to remove worktrees and optionally delete local mission branches after delivery.
- **RC hard-case coverage**: Added `rc fix`, `rc cancel`, and `rc hotfix` flows so the harness can now handle surgical RC cherry-picks, RC abandonment with branch cleanup, and production hotfix branch preparation.
- **Backport flow**: Added `rc backport` so dual-branch repositories can explicitly carry production hotfixes back into the development branch via a governed branch flow.
- **PR gate tracking**: Added `ship review` and `ship checks` so review feedback and CI state are now first-class merge gates in the harness runtime.
- **Deterministic merge**: Added `ship merge` so the harness can execute a real merge path instead of only recording a merged state transition.
- **Adapter manifest + full smoke**: Added `harness/adapter-manifest.yml` plus `scripts/harness-all-smoke-test.sh`, and taught bootstrap/sync to carry/report richer adapter packaging metadata.
- **Adapter version tightening**: bumped the harness adapter packaging line to `0.4.0-prototype`, and extended bootstrap smoke to verify manifest/full-smoke propagation plus sync-status manifest reporting.
- **Expanded release smoke**: `scripts/harness-release-smoke-test.sh` now validates RC surgical fix, RC cancel, hotfix preparation, and dual-branch backport in addition to the happy-path release flow.
- **Manifest-driven adapter contract**: `harness/adapter-manifest.yml` now declares protocol identity, adapter package files, bootstrap paths, and profile contract metadata, and the adapter line advanced to `0.5.0-prototype`.
- **Adapter validation surface**: Added `scripts/harness-adapter.sh` plus `harness_state.py adapter ...` so the package can be introspected and validated deterministically.
- **Manifest-driven bootstrap and smoke**: bootstrap now copies manifest-declared files instead of a shell hardcoded list, `harness-sync status` now compares adapter id as well as manifest version, and all smoke paths now validate adapter integrity explicitly.
- **Leader-led teammate runtime**: Added `scripts/harness-teammates.sh` plus first-class delegated teammate runs with bounded scope, reporting contracts, authority refs, forbidden sources, and parent/leader linkage.
- **Governed sibling communication**: teammate-to-teammate communication is now structured through sibling-approved signals instead of implicit chat.
- **Council packet strengthening**: meetings now record phase, participant roster, and role-bounded positions so councils can reflect independent review and cross-examination, not only final findings.
- **Teammate smoke coverage**: Added `scripts/harness-teammate-smoke-test.sh`, integrated it into `scripts/harness-all-smoke-test.sh`, and extended bootstrap coverage to include the teammate surface.
- **Completion handoff pack**: Added `docs/research/harness-completion-context-pack-2026-03-28.md` so future sessions can distinguish adapter/package progress from the larger orchestration mission.
- **Packaged runtime preference**: `scripts/harness-runtime.sh` now prefers packaged `harness/bin/<os>-<arch>/harness-state` binaries before falling back to a local Go build in `.harness/bin/`.
- **Runtime builder + dist smoke**: Added `scripts/harness-build-runtime.sh`, `harness/bin/README.md`, and `scripts/harness-runtime-dist-smoke-test.sh` so the adapter can prove a no-Python, no-local-Go runtime path when packaged binaries are present.
- **Council governance hardening**: meetings now carry role specs with allowed/forbidden sources, known-context notes, locked decisions, and rejected-finding reasons; `doctor` validates those governance fields explicitly.
- **Harness gap audit**: Added `docs/research/harness-gap-audit-2026-03-29.md` to track which earlier design inputs are absorbed, partial, or still missing, so future sessions do not drift from the original mission.
- **Research drift audit**: Added `docs/research/harness-research-drift-audit-2026-03-29.md` to re-check the original research set against the implemented harness and identify the remaining highest-risk drift areas.
- **Control-plane hardening**: runs and contracts now carry `ownership_claims`, `desired_state`, and dependency policy; teammate assignment rejects exclusive overlap by default; `drive wait`, `home`, and `doctor` now surface dependency and ownership control-plane state directly.
- **Authority + council advancement**: invalid authority levels are now rejected or surfaced, and meetings can advance phases via `meetings orchestrate` instead of relying entirely on manual phase changes.
- **Run-level authority adjudication**: added `scripts/harness-runs.sh authority-audit` plus doctor coverage for unresolved cross-surface authority conflicts, so leader views can now explain when higher-authority active signals disagree with current run state and when closed meeting outcomes still need adoption.
- **Contract control-plane sync**: run control-plane mutations now refresh contract snapshots for desired state, ownership selectors, branch, and worktree fields, and `maintenance normalize` now repairs stale historical contracts so old runtime state can be brought forward cleanly.
- **Live-state drift detection**: core smoke now proves the harness can detect and recover from run-record branch drift against a real worktree, and `doctor` now surfaces `run-live-state-drift` when live git state disagrees with the run record.
- **Meeting synthesis trace**: findings can now point back to supporting/challenging positions, synthesis audit now counts traced vs untraced findings, and teammate smoke proves a fully traced council outcome.
- **Mission graph materialization**: added `scripts/harness-missions.sh`, `harness/templates/mission.yml`, and durable `.harness/missions/` state so leader surfaces can reconcile first-class mission objects above run attempts.
- **Legacy runtime migration**: `scripts/harness-maintenance.sh normalize` now backfills ownership claims, desired state, dependency policy, and mission linkage for older pre-mission runs so existing `.harness/` state upgrades cleanly.
- **Selector-aware ownership transfer**: normalized legacy `ownership_claims` into structured `ownership_selectors`, added formal `runs transfer-ownership` semantics for bounded selector moves inside a mission, and extended teammate smoke to guard that control-plane behavior.
- **Signal authority adjudication**: added supersession checks so `doctor` can flag active signals that have been overtaken by fresher higher-authority run or meeting state, and `maintenance reconcile` / `signals adjudicate` can now resolve them deterministically.
- **Council orchestration hardening**: meetings now preserve `phase_history`, explicit `orchestration` state, linked cross-examination metadata, and synthesis audits with blocker/evidence coverage instead of treating phase labels as sufficient proof of debate quality.
- **Council health visibility**: `meetings summary` now exposes debate-round and synthesis-readiness signals, while `doctor` flags missing phase history, unlinked cross-examination, incomplete synthesis, and unresolved findings without follow-up actions.
- **Meeting authority audit**: same-claim council findings can now be grouped by claim/dimension, higher-authority winners are explained explicitly, unresolved same-authority conflicts block synthesis readiness, and leaders can inspect the result through `meetings authority-audit`.
- **Richer selector semantics**: ownership selectors now carry family/match/bounded metadata, path-family selectors detect semantic overlap (`dir`/`file`/`module`), and overlap admission/reporting is no longer limited to exact raw-string equality.

### Changed — Project Intelligence
- **Current project focus updated**: `_project/CURRENT.md` now reflects the shift from Wing-only transition work to the new Codex-native orchestration harness prototype.
- **Harness documentation expanded**: `harness/README.md` now documents the implemented surfaces and operating loop; `skills/README.md` now lists `/harness` and notes the `SKILL.md`-based Codex-native prototype pattern.
- **Harness session durability expanded**: `_project/sessions/2026-03-28-B.md` now captures the adapter-contract tightening pass so future sessions can resume from the package boundary work without replaying chat history.
- **Harness completion context pack**: Added `docs/research/harness-completion-context-pack-2026-03-28.md` to distinguish "current adapter/package state" from the larger remaining orchestration mission, and to define the minimum full-context bundle required for future sessions to finish the harness rather than only continue productization.
- **Harness runtime migration**: Replaced the Python harness runtime with a local-build Go CLI (`go.mod`, `cmd/harness-state`, `internal/harness/runtime.go`, `scripts/harness-runtime.sh`), updated all first-class harness wrappers to call the Go binary path, removed `scripts/harness_state.py`, and moved bootstrap/sync/smoke flows onto the Go adapter contract.
- **Adapter package contract**: `harness/adapter-manifest.yml` now includes the runtime builder and runtime-dist smoke path, bootstrap smoke checks those artifacts explicitly, and the docs now describe packaged binary preference plus the local-build fallback.

### Added — Skills Factory
- **Harness v0.32.0-prototype — Hierarchical semantic selector overlap**: semantic ownership is no longer flat exact-match only. Parent/child boundaries like `scope:release` and `scope:release/notes` now collide, including across `scope:` / `interface:` families. Teammate smoke proves the hierarchy-aware overlap gate end-to-end.
- **Harness v0.31.0-prototype — Bundle-level authority packets**: `missions authority-audit --all` now gives connected mission bundles their own `resolution_kind`, shared `selector_refs`, and a leader-facing `required_action`. Core smoke proves both conflicted and governed bundles now expose these higher-level control-plane hints.
- **Harness v0.30.0-prototype — Semantic selector-family overlap**: ownership governance now treats same-boundary `scope:` and `interface:` selectors as overlapping semantic claims, not unrelated labels. Teammate smoke proves the runtime rejects an exclusive `scope:release-contract` assignment when `interface:release-contract` already owns that boundary.
- **Harness v0.29.0-prototype — Bounded multi-round council controller**: `meetings orchestrate` now generalizes debate-round planning past the old second-round special case. Councils can emit third-round obligations when a second reply still loses, and they stop only when a lane resolves on authority/freshness or the configured debate-round cap is reached. Teammate smoke now proves a three-round council path before synthesis.
- **Harness v0.28.0-prototype — Governing mission for governed bundles**: bundle-aware mission graph authority now identifies the `governing_mission_id` when a governed bundle has a single root mission. Core smoke proves the governed bundle now names `issue-20001` as its governor, so the leader no longer has to infer control from dependency edges by hand.
- **Harness v0.27.0-prototype — Bundle-aware mission graph authority**: `missions authority-audit --all` now groups connected active missions into bundle-aware graph views, distinguishing conflicted bundles from governed bundles. Core smoke now proves both a conflicted overlap bundle and a dependency-governed bundle, and `home` now surfaces bundle counts alongside the existing mission-graph counts.
- **Harness v0.26.0-prototype — Dependency-governed mission ownership adjudication**: `missions authority-audit --all` now distinguishes unresolved parallel mission overlap from explicit dependency-governed overlap. When one active mission depends on another and their bounded ownership selectors overlap, the upstream dependency target wins authority over the shared boundary instead of the graph staying in generic conflict. Core smoke now proves this adjudication path, and `doctor` no longer emits dependency-state noise for inactive dependent missions.
- **Harness v0.25.0-prototype — Shared ownership transfer mode**: `runs transfer-ownership` now supports `--mode=copy` in addition to `move`, so bounded selectors can be intentionally shared between active runs when both runs explicitly allow shared overlap. Teammate smoke now proves the source retains ownership, the target receives the selector, and both contracts stay aligned under the shared-control policy.
- **arch-review v0.1.0 — Multi-role debate review skill**: Distilled from 3 projects (gcp-compute-infra, vi_agent), 7+ rounds, 14 opus agents. 8 mechanisms: Data Authority Hierarchy (M1), Role Templates (M2), Lessons Injection (M3), Standard Output Format (M4), Debate Orchestration via TeamCreate + opus teammates (M5), Self-Sufficiency Test (M6), Debug Recording (M7), Revolution self-optimization (M8). 27 bugs found and fixed during 5-round self-review. Complementary to /architect (which designs new systems).
- **Wing v0.1.0 — Anthropic official skill architecture rewrite**: Complete rewrite of teamwork skills as "wing" — AI autonomous development lifecycle protocol. 6 SKILL.md files (all ≤500 lines, -82% from teamwork's 9,758), 10 operations files (extracted heavy ops), VE single source reference, 8 migrated scripts. New install.sh dual-mode (legacy→commands/, modern→skills/). skill-lint.sh Layer 1 static analysis (12 checks). Named after Wings of Freedom (AoT). Architecture: `~/.claude/skills/{name}/SKILL.md` per Anthropic convention.
- **Teamwork v3.11.1 — VE hyperlink enforcement**: Dashboard templates: abstract "check each line" instruction replaced with concrete ✅/❌ examples (`**#306**` vs `#306`). LLMs mimic examples >> follow meta-instructions.
- **Teamwork v3.10.9 — Challenge mechanism + dashboard audit**: (1) team-issue: factual claims = hypotheses. (2) team-claim: detail ≠ accuracy, falsify before AGREE. (3) team-drive: MANUAL tasks go through Plan Validation (not skipped), resume skips re-assessment, AskUserQuestion = ⛔ HARD STOP. (4) team.md: worktree paths use absolute paths (no abbreviation), My Missions has explicit worktree conditional.
- **Teamwork v3.10.5 — Challenge mechanism fix (3-layer defense failure)**: Real incident: Issue #305 created with 7 sub-tasks for staging+prod, but staging was already working — 3 tasks redundant. All 3 challenge layers (issue/claim/drive) failed to catch it. Fixes: (1) team-issue Step 2: "Factual claims are hypotheses, not axioms" — verify system state before locking into sub-tasks. (2) team-claim Step 4c: "Detail ≠ accuracy" — falsify concrete claims before AGREE, unverified claim = cannot AGREE. (3) team-drive Plan Assessment: scope changed from code-only to ALL types (code + MANUAL) — "MANUAL is a classification for execution, not an exemption from challenge." (4) team-drive Manual Ops Handoff: AskUserQuestion upgraded from descriptive suggestion to ⛔ HARD STOP with explicit options (引导执行/Done/Ship first/Skip). Root cause: Plan Assessment was coupled to Execution Mode selection — 0 code tasks → skipped challenge entirely. Also fixed: dashboard worktree path rendering — LLM incorrectly shortened `~/Project/repo/.claude/worktrees/` to `~/.claude/worktrees/`, making `cd` commands non-functional. Added explicit instruction: `$HOME` → `~` substitution ONLY at prefix.
- **Teamwork v3.9.11 — Environment Envelope (Side Effect Triage)**: New Step 2e in team-ship: auto-scans git diff for non-code file changes (config, Dockerfile, docker-compose, .env) + greps for new env var references not in .env.example. P0 (missing env vars) blocks ship until resolved. P1 (infra/ops needs) prompts user to create Issue or PR comment. P2 (info) auto-injected into PR `## Infra Impact` section. team-drive gains Side Effect Check in per-task loop — AI records `[side-effect]` entries to AI Notes for remote ops and implicit assumptions that no diff can detect. Design: mechanism captures (deterministic diff), LLM annotates (judgment). Gate at ship, not block at drive.
- **Teamwork v3.9.10 — Cognitive upgrade round 2 (independent brainstorm)**: 3 independent thinkers found 5 additional gaps. Added: opposite test ("what if we don't do this at all?"), trade-off honesty ("when proposing alternative, state its cost"), blast radius assessment, response taxonomy (AGREE+EXTEND / CHALLENGE+PROPOSE / REDIRECT), proactive triggers in upstream skills (mediocre consensus, implicit decisions). Total drive essences absorbed: 17/20.
- **Teamwork v3.9.9 — Cognitive upgrade (drive essence injection)**: Inject 5 cognitive sentences across team-issue, team-claim, team-drive — identity language at line 1, bone/flesh analysis, Standalone Test self-check, anti-compliance before confirm, evidence-based Plan Validation. New Step 4b in team-claim: Contract confirmation before branch creation. Design: +sentences not +templates (Anthropic research: general instructions > prescriptive steps).
- **Teamwork v3.9.8 — Von Neumann worktree fix**: Deterministic worktree routing moved from LLM-interpreted skill.md into `tw-git.sh` scripts. New `claim-setup` (atomic branch+worktree setup) and `worktree-guard` (boundary mismatch detection) subcommands. Simplified team-claim Step 5 (~80 lines → ~15 lines), updated team-drive Step 0c and team-ship Step 1b to use scripted guard. 5 new unit tests. Root cause: worktree routing was a deterministic operation left to LLM interpretation — violated Von Neumann separation principle.
- **Teamwork v3.9.2 — Tips `shuf` crash on macOS**: Tips instruction now requires Read tool instead of Bash. Prohibits `shuf`/`sort -R` (GNU-only, absent on Darwin).
- **Teamwork v3.9.1 — Self-drive audit fixes**: hotfix continue routing bug, VE compliance, manual footer v3.1.0→v3.9.1, RC Freeze clarification, team-ship RC hint
- **Teamwork v3.9.6 — Cross-platform sed + remaining fixes**: `_sed_i` helper for GNU/BSD portability, tw-pr.sh undefined var, team-rc cancel loop, team-ship stash -u inconsistency
- **Teamwork v3.9.5 — Script hardening (12 fixes) + docs staleness + ecosystem**: rebase ref fallback, jq injection x2, git add -A→-u, bash array safety, hotfix re-entry, ship subcommand guard, notify worktree, temp file trap. Docs: skills/README v2→v3.9.x, CLAUDE.md version, idea pipeline stages. install.sh copy-mode cleanup + --target guard.
- **Teamwork v3.9.4 — Worktree config bug + dashboard RC + 48 unit tests**: Fix `_detect_config()` worktree logic, add RC section to dashboard, expand tips 15→29, add unit test suite
- **Teamwork v3.9.3 — Tooling + VE accuracy + docs staleness**: Fix `tw-bump-version.sh` path bug (manual never updated), VE accuracy (Promote Help contradiction, Cancel Flow raw git), CURRENT.md + ROADMAP.md brought current
- **Teamwork v3.9.0 — RC complete lifecycle**: `/team-rc fix` (cherry-pick to RC) + `/team-rc hotfix` (post-release production fix + backport). Full docs sync.
- **Teamwork v3.8.7 — Docs completeness audit**: Documentation completeness fixes
- **Teamwork v3.8.6 — Documentation audit + version sync**: Fix 3 version pollution files, sync manual to 3.8.6, add sync/share to all docs
- **Teamwork v3.8.5 — Cross-MC knowledge sharing**: `/team sync` (standalone rebase), `/team share` (cross-MC comments), L1 intake in team-drive
- **Teamwork v3.8.4 — Dashboard formatting + worktree init fix + tab colors**: Bold issue numbers, worktree _detect_dir fix, tab color refresh, bump-version semver validation
- **Teamwork v3.8.3 — /team-rc cancel + upstream-first RC**: RC cancel command, upstream-first merge model, post-merge PR checkbox
- **Teamwork v3.8.2 — Script沉淀 + DRY + UX dead-end fixes**:
  - 6 new script subcommands: `build-branch-name`, `extract-issue-number`, `worktree-main-repo`, `worktree-find-by-branch` (tw-git.sh) + `resolve-member`, `list-members` (tw-config.sh)
  - Shared `_PYTHON_MEMBER_PARSER` in tw-config.sh — 4 functions DRY-ed (net -148 lines)
  - `tw-bump-version.sh` handles 3 more version string patterns (inline backtick, parens, bash var)
  - AskUserQuestion 4-option cap enforced across all skills (was exceeding limit)
  - 3 UX dead-end fixes: team-ship done error guide, team-rc polling progress, team-claim worktree recovery
  - Wrap mode checkbox gate: team-drive skips already-checked tasks
  - E2E test Phase 1b: 13 tests for new script subcommands
  - Manual scripts architecture section updated (15→21 tw-git, +6 tw-config subcommands)
- **Teamwork v3.8.1 — Per-user setup + sync-all-checkboxes + 13 recovery hints**:
  - `sync-all-checkboxes` in tw-contract.sh (fuzzy-match Contract↔Issue checkboxes)
  - Per-user setup extracted to `/team init` Step 5 (worktree, tab title, notifications)
  - 13 `/team doctor` recovery hints added across all skills (was: raw shell commands)
  - `team-mode` subcommand in tw-config.sh (solo/team detection)
- **Teamwork v3.8.0 — Desktop notifications in `/team init`**:
  - macOS notification hook (`~/.claude/notify.sh`) with iTerm2 bell + fallback chain
  - Registered on Notification, Stop, TeammateIdle events
  - Subagent filtering (skip noise from inner agents)
- **Teamwork v3.7.1 — Visual Encoding Standard + audit sweep + doctor expansion**:
  - `docs/visual-encoding-standard.md`: single source of truth for color design language (6 semantic roles, 3 blue-family colors, DO/DON'T rules, rendering pipeline reference)
  - 5-scout comprehensive audit: 11 issues found and fixed (visual encoding violations, raw `rm` suggestions, missing tips/next-step, stale docs)
  - ONBOARDING.md modernized: all code-block examples → visual-encoding-compliant formatted markdown + Advanced Features section
  - Visual Encoding reference embedded in all 5 sub-skill file headers (drift prevention)
  - Doctor expanded from 8 → 11 checks: GitHub auth, Contract freshness (hash comparison), test suite health
  - Teamwork manual: Visual Encoding section expanded with markdown formatting rules + link to spec doc
- **Teamwork v3.7.0 — Smart sync: dashboard behind-detection + AI conflict resolution**:
  - Dashboard PR lines show `⚠️ N behind {base}` when branch is behind base branch
  - NEXT box prioritizes sync recommendation when PR is behind
  - `/team-ship sync` auto-resolves rebase conflicts via AI (no more manual STOP)
  - Sync summary shows resolved files + flagged ambiguities
  - Auto force-push after successful rebase
- **Teamwork v3.6.1 — Terminal tab title + Tips system**:
  - `/team init`: new Step 5b asks user to configure terminal tab title (auto-shows branch name)
  - Auto-detects shell (zsh/bash/fish), appends appropriate hook, idempotent
  - Tips system: every `/team-*` command shows random tip after output (progressive feature discovery)
  - Tips pool in `scripts/tw-tips.txt` (15 tips)
  - `install.sh`: now installs `*.txt` files from scripts/ directory
- **Teamwork v3.6.0 — Manual Ops Handoff**:
  - `/team-drive`: classifies sub-tasks as code (AI) vs `🔧 MANUAL` (human)
  - Completion splits into "CODE COMPLETE" + "YOUR TURN" handoff
  - `/team-ship`: MANUAL tasks don't block shipping, added to PR description
- **Teamwork v3.5.0 — `/team config` + dashboard version + docs completeness**:
  - `/team config` subcommand: shows current config.yml content, worktree status, edit instructions
  - Dashboard footer shows `Teamwork v{version}` (read dynamically from config)
  - Help/Learn outputs updated with all subcommands and dynamic version
  - Learn Section 11: Config Management (shared vs local config, edit workflow)
- **Teamwork v3.4.0 — Doctor actionability upgrade**:
  - `/team doctor fix`: interactive fix mode with severity-based confirmations
  - `/team doctor help`: self-contained usage guide
  - Branch table: new "vs Remote" column showing ahead/behind remote tracking
  - Check 8: Remote Sync diagnostic
  - Cross-skill UX audit: 5 new help routes, 2 branch readiness checks, ~12 error→guidance fixes

### Changed — Skills Factory
- **Teammate smoke widened for verdict adoption**: `scripts/harness-teammate-smoke-test.sh` now proves a same-authority debate lane can route into a second-round reply, resolve as `latest_response_prevails`, and auto-synthesize into accepted output without leaving unresolved action debt.
- **Home overlap accounting now respects exclusions**: leader-facing overlap detection no longer flags a parent run after it splits out an explicitly excluded child boundary, so `home` stays aligned with the same exclusion-aware ownership semantics used by admission control, mission authority audit, and `doctor`.
- **Teammate smoke widened again**: `scripts/harness-teammate-smoke-test.sh` now proves `split` transfer, including retiring the temporary child owner, reconciling superseded signals, preserving the parent selector, recording exclusions, and keeping both contracts healthy after delegation.
- **Home/doctor graph awareness**: leader surfaces now include graph-level mission conflict/adjudication counts instead of hiding multi-mission control drift under generic overlap reporting.
- **Core smoke widened again**: `scripts/harness-smoke-test.sh` now proves a cross-mission overlap can be created, surfaced by mission-graph authority audit, reported by `doctor`, and then cleared cleanly.
- **Mission migration path widened**: `maintenance normalize` now regenerates older mission state when the new authoritative control field is missing, so upgraded repositories do not stay noisy in `doctor`.
- **Teammate smoke widened**: `scripts/harness-teammate-smoke-test.sh` now proves that the freshest teammate activity can differ from the authoritative mission-control run without silently handing mission truth away from the leader surface.
- **Distribution strategy formalized**: `docs/harness-distribution-strategy.md` now defines source-bootstrap, packaged-binary, and portable runtime-bundle tiers for the `Harness` prototype line.
- **Bootstrap/package coverage widened**: bootstrap smoke now checks the runtime bundle packaging surfaces in addition to build/dist/package identity surfaces.
- **Manifest/product identity alignment**: `harness/adapter-manifest.yml` now carries public prototype identity metadata, and the public contract now references the accepted naming/positioning ADR instead of treating naming as unresolved.
- **Core smoke authority coverage**: `scripts/harness-smoke-test.sh` now validates that a manually-staled contract snapshot loses to the active run before continuing into live-state authority drift checks.
- **Council control-plane accounting**: meeting templates, summaries, and doctor now track planned/resolved/open second-round debate obligations, and teammate smoke now validates the bounded two-round council path before synthesis.
- **Harness package portability**: `harness/adapter-manifest.yml` now points at the durable public contract and includes it in manifest-driven package files, while bootstrap smoke now verifies bootstrapped repos receive that contract too.
- **Teamwork v3.4.1 — Worktree mode: shared → local**:
  - Worktree mode moved from shared `config.yml` to per-user `git config --local teamwork.worktree true`
  - Dashboard always shows worktree map when enabled (★ current, MC association, ⚠️ mismatch)
  - Init merge mode: auto-migrates old `worktree:` section to `git config --local`
  - Doctor Check 5: detects idle worktrees, context mismatch
- **Teamwork v3.2.0 — Unified visual design language**:
  - Standardized emoji encoding across all 6 skill files (priority dots, CI icons, progress bar)
  - Dashboard performance: 7-8 serial API calls → 4-5 parallel
- **Teamwork v3.1.0 — Untracked issue detection + learn audit**:
  - Dashboard UNTRACKED section: detects GitHub Issues assigned but missing teamwork labels
  - `/team-issue fix #N`: adds missing teamwork labels to untracked Issues
  - `/team learn`: complete 10-section design philosophy + visual diagrams

### Fixed — Skills Factory
- **Teamwork v3.5.1 — Design review: logic gaps + docs alignment**:
  - team-ship: base branch guard prevents accidental ship from main/develop
  - team-ship: CI failure preserves Contract (enables drive→ship iteration)
  - team-drive: Contract-missing fallback detects mission branch + open Issue, suggests recovery
  - docs: teamwork-v2-architecture.md archived, team-doctor-design.md Check 8 added
- **Teamwork v3.4.3 — Merge mode logic gaps**:
  - Init merge mode asks worktree preference when never configured
  - Quality section uses safe defaults instead of undefined Step 3 variables
- **Teamwork v3.4.2 — Worktree dashboard render reliability**:
  - Worktree section restructured as mandatory pre-render (was skipped by LLM in conditional template)
- **Teamwork v3.2.1 — Script extraction + deep bug fix + stale cleanup**:
  - Milestone prefix matching backport from vi_agent (4 locations + batch gap)
  - Learn: direct embedded output, no file reading (fix exit code 1)

### Changed — Skills Factory (v2.x)
- **Teamwork v2.7.3 — Init audit + dynamic branch protection**:
  - team.md Step 5c: hardcoded branch protection → config-driven (`quality.ci` + `quality.review_required` → 4 combinations)
  - Added `allow_force_pushes: false` + `allow_deletions: false` to protection rules
  - vi_agent: created post-merge.yml workflow, updated skill_version 2.5.0 → 2.7.2, set merge settings (squash-only)
  - vi_agent: branch protection enabled (require PR, review=0, no CI, block force push)
  - 18 requirements documented for release branch lifecycle, deploy workflow, issue version tracking

### Fixed — Skills Factory
- **Teamwork v2.7.2 — vi_agent E2E Validation** (3 bugs found, 3 fixed):
  - 1 CRITICAL: `_branch_pattern()` `${p:-default}` with `{}` in default — bash closes parameter expansion at first unmatched `}`, garbling branch names on any project with custom branch_pattern
  - 1 MEDIUM: `check-freshness` with non-existent file exits 3 instead of 2 — `exit` in function kills `$(…)` subshell before `|| true` can fire; added pre-check + moved `|| true` outside subshell
  - 1 LOW: `protect-check` USAGE comment said exit 1, actual code exits 3
- **Teamwork v2.7.1 — Script Testing Audit** (21 bugs found, 14+4 fixed across 4 scripts):
  - 2 CRITICAL: tag check matched branches not tags, merge-to on wrong branch after checkout failure
  - 6 HIGH: silent ensure-base failure, sed injection, fetch kills script, grep metachar, commit-type false positives, toggle-task false success
  - 6 MEDIUM: missing grep anchor, stderr pollution, detached HEAD, label transition drop, unused PR base param, body validation message
  - 4 LOW (fixed in follow-up): label verify substring, worktree-remove return 0, log-since feedback, branch_pattern double lookup

### Changed — Skills Factory
- **Teamwork v2.7.0 — Script Extraction** (Von Neumann Mechanism-Content Separation):
  - 4 new scripts: tw-git.sh (13 subcommands), tw-label.sh (4), tw-contract.sh (8), tw-pr.sh (7)
  - 6 skill files refactored: deterministic operations → script calls, LLM judgment sections unchanged
  - Design: scripts handle execution, skill.md handles orchestration + judgment
  - ~64% of operations moved from LLM-interpreted prose to executable scripts

### Fixed — Skills Factory
- **Teamwork v2.6.1 — Falsification Audit** (Popper's method applied across all 6 skill files + scripts):
  - 6 CRITICAL: freshness check always false-positive, PyYAML dependency, branch protection 422, CHANGELOG shell injection, stale merge, wrong push branch
  - 15 HIGH: init force unimplemented, hardcoded label prefixes (5 files), macOS sed syntax, milestone quoting, undefined variables, missing --state flags, tag existence check
  - All skill files bumped to v2.6.1

### Changed — Skills Factory
- **Teamwork v2.6.0 — Dual-branch model + release safety gates**:
  - Dual-branch support: `conventions.release_branch` enables base/release separation (backward compat when same)
  - team-release: staging deploy gate, merge base→release, CHANGELOG auto-append
  - team-drive: branch safety check prevents code execution on protected branches
  - team-claim: explicit base_branch checkout before feature branch creation
  - team-issue: `--milestone` parameter for per-issue milestone override
  - team-ship: sync/done operations use config base_branch; done D3/D4 now fallback (Action-first)
  - team.md: post-merge Action (auto-label + checkbox + milestone check), repo merge settings, config-driven branch protection + CI templates
  - Zero project coupling in skill files — fully generic design
- **Teamwork v2.5.2 — Dashboard queue overflow fix**:
  - QUEUED section now sorted by priority (P0→P1→P2→P3→none) and truncated to top 5 with overflow count
  - New `/team queue` subcommand shows full queued issue list grouped by priority
  - Prevents dashboard from growing unbounded when many issues are queued

### Added — Skills Factory
- **Teamwork v2 implementation** — 4 GitHub-first skills replacing v1:
  - `/team` — init system (auto-detect project, generate config/CI/hooks/templates) + dashboard from live GitHub data
  - `/team-claim` — claim GitHub Issue → generate AI-enriched Mission Contract → create branch
  - `/team-drive` — execute mission from Contract with sub-task verification loop
  - `/team-ship` — push → create PR (auto-closes Issue) → watch CI → cleanup Contract
- **Teamwork v2 architecture** — GitHub-first redesign (ADR-0010). Issues as task source, ephemeral Mission Contracts as AI execution views, auto-generated Actions CI + git hooks + Issue templates. Full design at `docs/teamwork-v2-architecture.md`. Supersedes v1 (ADR-0009).
- **Teamwork v2.1 — vi-agent feature absorption** — 6 features from vi-agent's `.teamspace` system:
  - Prefixed labels (`status:wip`, `priority:P1`, `status:blocked`, `status:done`)
  - Claim comment visibility (team transparency via `gh issue comment`)
  - Optional worktree isolation (parallel development per mission)
  - Optional version management (version labels + dashboard progress)
  - `/team-ship done` (post-merge cleanup), `/team-ship review` (AI code review), `/team-ship sync` (rebase)
- **Teamwork v2.1 architect review** — 13 bug fixes:
  - 3 critical: worktree contract access, Issue template label mismatch, branch pattern hardcoding
  - 3 high: config format collision, missing config health check, branch name parsing
  - 5 medium: sed pattern consistency, git staging guidance, worktree detection logic
  - 2 flow bugs: Step 1b→1c routing, $BRANCH_PATTERN extraction

### Changed — Documentation
- **docs/ reorganization** — Created `docs/README.md` index (active/archived/research), moved 3 research docs to `docs/research/`, added archived status banners to 4 stale docs, fixed stale cross-references to teamwork-v2-architecture.md

### Removed — Skills Factory
- **Teamwork v1 skills** — `/teamwork`, `/team-get-mission`, `/team-complete` removed (replaced by v2)

### Added — Showcase
- **SkillBuilder MVP** (ACT IV) — 5-step form to create custom Claude Code skills: Meta, Execution Model, Phases, Rules, Tools. Real-time preview with collapsible panel. Download generates valid skill.md with substantive prompt engineering content. Full i18n (EN/ZH).
  - `showcase/src/lib/skillTemplate.ts` — types + template generation (~310 lines)
  - `showcase/src/components/SkillBuilder.tsx` — stepper form + preview + download (~527 lines)
  - ~80 new translation keys in `showcase/src/i18n/translations.ts`

### Changed — Showcase (Performance)
- **Bundle code splitting** — 627KB single chunk → 221KB main + vendor chunks + 20 lazy-loaded chunks. React.lazy() for 19 components (15 ACT III sections + 3 ACT II components + SkillBuilder). Vendor splitting (react, framer-motion) in vite.config.ts. Vite 500KB warning eliminated.

### Added — Showcase (Animation Enhancements)
- **ExecutionSimulator topology** — 4-node SVG topology diagram (Team Lead → STL → Worker-1/Worker-2) with animated message flow dots (green dispatch, blue return) synced to simulation playback.
- **SwarmTopologySection message flow** — 12 animated message dots flow along connection paths in 3 phases (dispatch → inner dispatch → return). Auto-plays on scroll into view. Click any node to focus its message paths (non-involved nodes/connections dim to 8% opacity). Click again or click background to reset.
- **Consistent color coding** — dispatch=#22c55e (green), return=#3b82f6 (blue) across both components.

### Fixed — Showcase (Code Health)
- **SwarmTopologySection** — fixed `react-hooks/set-state-in-effect` ESLint error by replacing `useEffect`+`setMessageKey` with `activeNodeId` in React key (same remount behavior, no effect).
- **HeroSection** — fixed `react-hooks/set-state-in-effect` error by extracting `HeroAnimation` inner component with `key={lang}` reset pattern. Fixed `exhaustive-deps` warning with `useMemo`.
- **Context files** — suppressed `react-refresh/only-export-components` false positives for hook exports.
- **Semantic HTML** — added `<header>`, `<main>`, `<section aria-label>` landmarks to App.tsx.

### Removed — Showcase (Dead Code)
- **DemoSection.tsx** — removed orphaned 24KB component (not imported in App.tsx).
- **9 dead i18n keys** — removed `demo.note1.*`, `demo.note2.*`, `demo.note3.*` from translations.ts (only used by removed DemoSection).

### Changed — Skills Factory
- **`/new-project` → `/erwin-init`** — Renamed and redesigned per ADR-0008: from scaffolding tool to cognitive infrastructure installer. Zero-question entry, optional GOALS.md conversation. `erwin-` prefix reserved for framework-core skills.

### Added — Skills Factory
- **`skills/`** — New: source of truth for all erwin Claude Code skills (ADR-0007)
- **`skills/erwin-init/`** — `/erwin-init` v2.0.0: initializes erwin cognitive framework for any project; now version-controlled in erwin repo with README + CHANGELOG
- **`skills/visual-test/`** — `/visual-test` v1.0.0: screenshot-driven visual bug detection; now version-controlled in erwin repo with README + CHANGELOG
- **`install.sh`** — Symlink-based skills installer: `./install.sh` → all skills globally available; `git pull` auto-updates installed skills
- **`.claude/commands/`** — Project-level incubation stage for skills under development
- **`docs/erwin-skills-incubator-architecture.md`** — Full architecture design for skills incubator pattern
- **`_project/decisions/0007-skills-incubator-structure.md`** — ADR-0007 documenting the incubator structure decision

### Changed — `_project/ideas/`
- `visual-test-skill.md` moved from `ideas/` root → `ideas/skills/` (correct pipeline location)

### Added — Meta / Cross-Project (previous sessions)
- **`_project/decisions/0006-new-project-skill.md`** — ADR-0006 documenting /new-project skill design decisions
- **`~/.claude/projects/-Users-liyasong-Project-collov-erwin/memory/MEMORY.md`** — Claude auto-memory initialized for erwin project

### Changed — `.claude/rules/`
- **`ai-workflow.md`** — 将 "Session Start 读 CURRENT.md" 改为原则驱动的触发条件。不再罗列场景，改为一个判断测试：「行动如果方向错了，撤销成本高吗？」高成本行动（写代码、做决策、动手前）必须先读；纯对话（解释概念、讨论方案）无需读。Added step 4 to Session End: update MEMORY.md.

---

## [0.5.0] - 2026-03-02 (V2 Phase 2 — OriginTrace)

### Added — `showcase/`
- **`OriginTrace`** — d3-force directed graph showing drive-skill's intellectual genealogy (ACT II, after PrincipleCards). 15 nodes: 5 drive-concept nodes (matching PrincipleCards group colors A-E) + 10 historical source nodes (OODA Loop, Cybernetics, Mission Command, Mythical Man-Month, CPM/PERT, Amdahl's Law, Information Theory, Agile Manifesto, 孙子兵法, Double Diamond). Click any node to open a detail panel showing the origin description and how it manifests in drive-skill. `motion.svg` whileInView fade-in, `motion.g` stagger entrance, edge highlighting on hover/select.

### Changed — `showcase/`
- **`App.tsx`** — OriginTrace deployed in ACT II after PrincipleCards (previously placeholder comment).

### Dependencies
- `d3-force` added to dependencies
- `@types/d3-force` added to devDependencies

### Rationale
OriginTrace completes the ACT II "理解" experience: PrincipleCards explain the WHAT of each principle, OriginTrace traces the WHY (intellectual lineage). Together they answer "what does drive-skill do?" and "where did these ideas come from?"

---

## [0.4.0] - 2026-03-02 (V2 Phase 1 — Four-Act Structure)

### Added — `showcase/`
- **`ExecutionSimulator`** — Before/after animated comparison component (ACT I). Sequential vs parallel execution shown with Framer Motion progress bars, "3x faster" callout, replay button.
- **`PrincipleCards`** — Five Core Principles component (ACT II). Each card shows Chinese/English name, summary, expandable detail, and historical "thought sources" (OODA Loop, Mission Command, Amdahl's Law, Shannon, Sun Tzu).
- **`ActIIIWrapper`** — Collapsible wrapper for 15 deep sections (ACT III). Collapsed by default, shows section index, expands with Framer Motion animation.

### Changed — `showcase/`
- **`App.tsx`** — Restructured from flat 18-section scroll into Four-Act layout (ACT I 震撼 / ACT II 理解 / ACT III 深度 / ACT IV 创造). DemoSection replaced by ExecutionSimulator. ClaudeMechanicsSection moved to ACT II.
- **`NavBar.tsx`** — Replaced 18 individual section links with 4 Act anchor links (ACT I-IV), with intersection-observer-based active tracking.
- **`i18n/translations.ts`** — Added nav translation keys for four acts (en + zh).

### Rationale
V2 Phase 1 implements the approved Four-Act architecture from `docs/drive-showcase-v2-architecture.md`. The new structure enables the "15-second shock → 5-minute understanding → create" journey.

---

## [0.3.0] - 2026-03-02

### Added
- **`_project/` — Project Intelligence Layer** — three-layer information architecture:
  - `GOALS.md` — Vision + quarterly OKRs (the north star, read at session start)
  - `CURRENT.md` — Living status document (moved from `docs/CURRENT_STATE.md`)
  - `ROADMAP.md` — Milestone plan with version targets
  - `ideas/` — Idea Garden with four maturity stages (Seed/Sprout/Plant/Fruit)
    - `INBOX.md` — Raw captures, processed weekly
    - `skills/` — Drive-skill-like skill ideas
    - `meta/` — Universal cross-project patterns
  - `decisions/` — ADRs (moved from `docs/decisions/`)
  - `sessions/` — Conversation logs (moved from `docs/sessions/`)
  - `issues/` — Lightweight issue tracking (ISS-NNNN naming)
  - `releases/` — Narrative release notes
- **`.claude/rules/`** — Modular Claude Code instruction files
  - `git-conventions.md` — Commit format, scope definitions
  - `ai-workflow.md` — Session protocol, what to do at start/end
- **`docs/project-management-architecture.md`** — Full design specification
- **`_project/ideas/meta/project-management-template.md`** — Seed idea for universal project template skill

### Changed
- **`reference-skills/` → `reference/`** — cleaner naming; signals read-only source material
- **`CLAUDE.md`** — Updated paths, added GOALS reference, cleaner structure
- **`_project/CURRENT.md`** — Updated to reflect new structure and next tasks

### Rationale
The docs/ directory was a mix of project management files (CURRENT_STATE.md, sessions/,
decisions/) and technical documentation (architecture docs, analysis). This creates confusion
when navigating as a project manager vs. as a developer. The `_project/` layer creates a
dedicated home for all project intelligence, while `docs/` becomes pure technical documentation.
Additionally, the user correctly identified that the entire `_project/` pattern is a universal
capability for AI-assisted projects — captured as a Sprout idea in `_project/ideas/meta/`.

---

## [0.2.0] - 2026-03-02

### Added
- **Version management system** — git init, CHANGELOG.md, CLAUDE.md, ADR structure, session logs
- **CLAUDE.md** — project bootstrap context for AI sessions (readable at conversation start)
- **docs/decisions/** — Architecture Decision Records (ADR) directory
- **docs/sessions/** — conversation artifact archive (designs, decisions, research per session)
- **docs/CURRENT_STATE.md** — living document of current project state for session bootstrapping
- **docs/drive-showcase-v2-architecture.md** — V2 interactive showcase architecture design

### Rationale
The project had no git history and relied entirely on Claude Code's session memory for context
continuity. This release establishes the knowledge infrastructure: every decision, artifact,
and design gets versioned and committed, so any new conversation can bootstrap full context
by reading CLAUDE.md + docs/CURRENT_STATE.md rather than depending on volatile AI memory.

---

## [0.1.0] - 2026-03-01

### Added
- **showcase** — React/Vite/TypeScript web app showcasing drive-skill
  - 18 sections covering all aspects of drive.md, self-drive.md, architect.md
  - Bilingual (EN/ZH) with LangContext
  - Framer Motion animations throughout
  - Annotation system
  - CHANGELOG.md for showcase itself
- **reference-skills/drive-skill/** — local copy of the drive-skill repository
  - commands/: drive.md, self-drive.md, architect.md
  - hooks/: notify.sh
  - sounds/: 5 bundled TTS mp3 files
  - messaging/: slack and telegram scripts
  - pantheon/: 1000+ greatest minds knowledge base
  - config/: settings templates
- **docs/** — project documentation directory
  - drive-skill-analysis.md
  - visual-enhancement-architecture.md

### Rationale
Initial state of the erwin project at the point of git initialization. The showcase
had reached v1.3.0 (full content parity with drive.md) and was ready for V2 visual upgrade.
