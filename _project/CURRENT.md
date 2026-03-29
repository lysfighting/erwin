# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.46.0 (in progress)
- **Last Updated**: 2026-03-29
- **Last Session**: 2026-03-29 (Scheduler now reacts to fresh follow-up signals)

## Current Phase

**Phase**: Codex-native orchestration harness prototype

## What Was Just Completed

**Session 2026-03-29** — Harness prototype line stayed complete while bundle authority started driving scheduler behavior directly

- **Bundle-aware scheduler consumption added**: `missions schedule` now consumes mission-bundle authority packets instead of only the raw conflict list, so governed losing missions are deferred behind the current bundle governor, escalation-required bundles become explicit human decisions, and governed-but-ambiguous bundles no longer masquerade as ordinary runnable work
- **Scheduler output widened**: mission classifications now surface `bundle_id`, `bundle_state`, `bundle_role`, `bundle_resolution_kind`, `bundle_governing_mission_id`, and `bundle_required_action`, making bundle-level control visible directly in durable scheduler snapshots
- **Leader-facing scheduler priority fixed**: when required human bundle decisions coexist with a still-valid foreground batch, scheduler `next_action` now prioritizes the human decisions instead of pretending foreground selection is the only thing that matters; snapshots also emit `bundle_role_tally` so the leader can see the portfolio composition directly
- **Home/scheduler parity restored**: `home` now reuses scheduler decision-first guidance when escalation-required bundle decisions outrank generic cross-mission conflict wording, and it surfaces `scheduler_bundle_role_tally` so the leader sees the same bundle portfolio shape in both places
- **Bounded portfolio admission added**: scheduler now consumes profile-backed `max_parallel_runs`, so parallel execution no longer widens to all runnable missions by default; overflow runnable missions are deferred with explicit `defer_due_parallel_capacity` actions and the snapshot records `batch_admission_policy`, `max_parallel_runs`, and `capacity_deferred_*` metadata
- **Fresh-signal rescheduling added**: scheduler now treats fresh `requires_meeting`, `interface_change`, `dependency_discovered`, and `blocked` signals as real control-plane triggers instead of display-only telemetry, reclassifying missions and emitting pause/reschedule actions when the signal has not yet been adopted into runtime state
- **Core smoke widened again**: `scripts/harness-smoke-test.sh` now proves pairwise governed bundles defer the losing mission behind the governor and that escalation-required bundles become explicit `establish_bundle_governor` scheduler decisions

- **Stateful debate batch execution added**: `meetingDebatePlan` now computes `current_wave_batch_id`, `ready_wave_batch_count`, `blocked_wave_batch_count`, `resolved_wave_batch_count`, and batch-local `state/current` so the control plane can treat one debate batch as active while explicitly blocking the rest behind it
- **Synthesis/output guard tightened**: `meetings orchestrate` now refuses to jump from `synthesis` to `output` while an open debate batch still exists, even if someone manually pushed the meeting phase forward before the controller finished the active batch
- **Leader visibility widened again**: `meetings summary` now surfaces current/ready/blocked/resolved batch state, and `doctor` now warns when open debate obligations lack a deterministic current batch or when a meeting reaches `synthesis`/`output` while a debate batch still remains open
- **Pairwise bundle governor adjudication added**: mission-graph authority no longer reserves priority/risk-based governor selection for only larger components; a two-mission conflicted bundle can now become governed when one mission clearly outranks the other on authority tuple
- **Bundle escalation semantics sharpened**: conflicted bundles now remain generic `bundle_conflict` only when no meaningful authority signal exists; when authority signals exist but still fail to produce a unique governor, the control plane escalates to `bundle_governor_required` instead
- **Core smoke widened again**: `scripts/harness-smoke-test.sh` now proves a pairwise overlapping mission bundle becomes `bundle_governor_adjudicated` without dependency, council, or selector-specificity help when one mission has a uniquely stronger priority/risk tuple
- **Direct output closure added**: `meetings orchestrate` now advances straight from a newly synthesized council record into `output` when no synthesis gaps remain, so multi-lane, authority-resolved, and aligned councils no longer need a second orchestration pass just to leave `synthesis`
- **Teammate smoke widened again**: `scripts/harness-teammate-smoke-test.sh` now proves auto-synthesized, verdict-adopted, authority-resolved, and aligned councils all close to `output` in the same orchestration call that materializes their governed findings
- **Multi-lane debate batching added**: `meetingDebatePlan` now groups open route/round obligations into `wave_batches`, so one participant can receive a visible batch of same-round work across multiple conflicting claim lanes instead of only a flat list of isolated obligations
- **Meeting summary widened again**: `meetings summary` now surfaces `wave_batch_count`, `open_wave_batch_count`, `multi_lane_batch_count`, and `max_wave_lane_count`, making multi-lane debate execution visible to the leader
- **Teammate smoke widened again**: `scripts/harness-teammate-smoke-test.sh` now proves a two-lane conflicting council produces one grouped debate batch for the same bounded challenger

- **Council-governed bundle authority added**: `missions authority-audit --all` now treats a single active council anchor as a first-class ownership adjudication source, so overlapping missions can become a governed bundle even without an explicit dependency edge when one mission is already chairing the relevant council
- **Bundle packets widened again**: governed mission bundles can now emit `council_governed_mission_ownership` and `council_governed_bundle`, with `governing_mission_id` and a leader-facing action that keeps authority on the council-host mission until adoption
- **Core smoke widened again**: `scripts/harness-smoke-test.sh` now proves an active council on `issue-20001` converts a later overlapping mission into governed bundle authority instead of unresolved graph conflict

- **Leader-facing scheduler baseline added**: `missions schedule` now emits durable `.harness/scheduler/snapshots/` records with per-mission classification, selected foreground batch, execution-mode recommendation (`solo`, `managed_delegation`, `scheduled_parallel_missions`), required human decisions, follow-up signals, and reschedule actions
- **Scheduler semantics hardened**: the runtime now distinguishes true cross-mission dependency blocking from mission-internal teammate dependencies, so delegated teammate work no longer poisons mission scheduling with fake self-dependency
- **Release-adjacent scheduling tightened**: scheduler serialization now keys off real release-gate evidence instead of template-default `unknown` review/check states, so isolated ready missions can widen into governed parallel batches instead of being accidentally serialized
- **Leader home widened**: `home` now surfaces scheduler execution mode, selected mission batch size, classification tally, and scheduler-derived next action so the leader can reason from durable control-plane state instead of ad hoc heuristics alone
- **Canonical docs synced with scheduler reality**: the operating-model, architecture, public-contract, and handoff docs now treat scheduler snapshots and execution-mode selection as present runtime behavior rather than purely future aspiration
- **Selector-specific bundle adjudication added**: mission-graph authority can now let a narrower ownership selector govern a broader overlapping mission without an explicit dependency edge, producing `selector_specificity_governed_mission_ownership`, `single_winner_governed_bundle`, and a concrete leader-facing `required_action` instead of dropping every such case into `bundle_conflict`
- **Smoke coverage widened again**: core smoke now proves solo scheduling, dependency-blocked reschedule hints, and governed parallel batches, while teammate smoke proves `managed_delegation` scheduling for a leader mission with bounded helper runs

- **Authority-resolved council lanes added**: `meetings orchestrate` now recognizes when a higher-authority support/challenge lane already outranks its opponent during independent review, records a governed `higher_authority_prevails` route verdict, and advances directly into synthesis instead of forcing unnecessary cross-examination
- **Aligned councils can skip empty debate**: meetings with enough independent review but no conflicting claim lane can now move straight into synthesis instead of getting stuck waiting for a synthetic cross-examination step that adds no new judgment
- **Canonical docs synced with council semantics**: the architecture, operating-model, and public-contract docs now explicitly state that cross-examination is the normal path for unresolved contested lanes, not a mandatory ritual for aligned or already authority-settled lanes
- **Council summaries widened for direct synthesis**: `meetings summary` now surfaces `resolved_without_cross_examination_count` and `authority_resolved_route_count`, so the leader can see when the controller deliberately skipped debate because no unresolved lane remained
- **Teammate smoke widened again**: `scripts/harness-teammate-smoke-test.sh` now proves both new control-plane paths end-to-end, covering higher-authority direct synthesis and zero-conflict direct synthesis before output

- **Go runtime core added**: introduced a real Go module plus `cmd/harness-state` / `internal/harness/runtime.go`, migrating the harness control-plane and state-engine behavior out of `scripts/harness_state.py`
- **Shell surfaces preserved**: all first-class harness entrypoints now route through `scripts/harness-runtime.sh`, which builds and executes the local Go CLI while keeping the existing `home`, `doctor`, `issue`, `claim`, `drive`, `teammates`, `ship`, `rc`, `meetings`, `maintenance`, and adapter/profile/path surfaces stable
- **Python runtime dependency removed for harness usage**: normal harness operation, package inspection, bootstrap, sync, and smoke execution no longer depend on `python3`
- **Packaged runtime preference added**: `scripts/harness-runtime.sh` now prefers `harness/bin/<os>-<arch>/harness-state` when present, with local Go build under `.harness/bin/` only as a fallback
- **Binary build and dist validation added**: `scripts/harness-build-runtime.sh` now builds packaged binaries, and `scripts/harness-runtime-dist-smoke-test.sh` proves the adapter can run without `python3` and without `go` in `PATH` when the packaged runtime is present
- **Council governance hardened**: meetings now carry first-class `role_specs` with allowed/forbidden sources, known-context notes, locked decisions, and rejected-finding reasons; `doctor` validates those governance requirements instead of treating councils as loose notes
- **Gap audit written down**: `docs/research/harness-gap-audit-2026-03-29.md` now records which original design inputs are absorbed, partial, or still missing, so future sessions can resume from an explicit audit instead of replaying chat history
- **Research drift audit added**: `docs/research/harness-research-drift-audit-2026-03-29.md` now re-checks the original design docs against the current harness and confirms the remaining highest-risk drift areas are desired-state control, explicit ownership/overlap, harness-wide authority enforcement, and automated council orchestration
- **Bounded multi-round council control added**: `meetings orchestrate` now plans governed round 2 / round 3 / later replies from policy instead of stopping at a second-round special case, and teammate smoke now proves a three-round debated lane before synthesis
- **Semantic selector-family overlap added**: ownership governance now treats same-boundary `scope:` and `interface:` selectors as overlapping semantic claims, and teammate smoke proves the runtime now rejects exclusive cross-kind semantic collisions instead of only path-family collisions
- **Bundle-level authority packets added**: `missions authority-audit --all` now emits bundle-level `resolution_kind`, shared `selector_refs`, and leader-facing `required_action`, so the control plane explains what to do with a conflicted or governed mission bundle instead of only listing pairwise adjudications
- **Hierarchical semantic selector overlap added**: ownership governance now treats parent/child semantic boundaries like `scope:release` and `scope:release/notes` as overlapping governed ownership, including across `scope:` / `interface:` families
- **Split ownership transfer added**: `runs transfer-ownership` now supports `--mode=split`, so a broader parent selector can delegate a governed child selector through explicit ownership exclusions without surrendering the parent boundary
- **Exclusion-aware overlap accounting added**: leader-facing overlap detection in `home` now respects ownership exclusions, so split parent/child delegation does not appear as a false overlap once the child subtree has been explicitly carved out
- **Teammate smoke widened for split delegation**: `scripts/harness-teammate-smoke-test.sh` now proves the full split-transfer path, including temporary child-owner retirement, superseded-signal reconciliation, source exclusion recording, target child acquisition, and clean contract alignment on both runs
- **Debate verdict adoption added**: auto-synthesis now honors a governed `latest_response_prevails` route verdict, so a same-authority debate lane can materialize directly into accepted output once the latest bounded reply wins
- **Single architecture doc added**: `docs/harness-architecture.md` now serves as the complete design overview for `Harness`, unifying identity, lineage, retained protocol semantics, dropped historical baggage, control-plane model, Go runtime rationale, current capability surface, and remaining expansion gaps in one place instead of forcing fresh sessions to reconstruct the architecture from scattered docs
- **Documentation system formalized**: `docs/harness-documentation-map.md` now manages the canonical Harness doc set, and `docs/harness-operating-model.md` now captures the leader/teammate/council execution model so future sessions do not need to reconstruct the collaboration design from architecture notes plus chat memory
- **Teammate smoke widened for verdict adoption**: `scripts/harness-teammate-smoke-test.sh` now proves a second-round same-authority debate lane resolves as `latest_response_prevails` and lands as accepted output without unresolved action debt
- **Debate route obligations added**: `meetings orchestrate` now materializes claim-lane route obligations with assigned bounded challenger roles/runs, while `meetings summary` and `doctor` surface resolved vs open debate routes directly instead of only reporting routed/unrouted claims after the fact
- **Live worktree absence promoted into authority audit**: `runs authority-audit` now treats a missing provisioned worktree as first-class live-state drift (`live_worktree_missing`) instead of relying only on `doctor` to notice the broken path
- **Glob ownership selectors added**: the ownership model now understands `glob:` path selectors, so overlap control can express bounded markdown/file-set style ownership instead of forcing everything into `file:` or coarse `dir:` claims
- **Debate wave obligations added**: `meetings orchestrate` now materializes missing independent-review wave obligations, so the leader can see who still owes input before debate can advance instead of relying only on phase labels and coverage notes
- **Ship-gate authority adjudication added**: `runs authority-audit` now lets `verification_status`, `review_status`, and `checks_status` explicitly override an inconsistent optimistic `ship_status`, so gate evidence can win over stale delivery state
- **Public package contract formalized**: `docs/harness-public-contract.md` now defines the durable public boundary for the Codex-native harness package, including protocol/adapter/profile separation, required surfaces, runtime guarantees, distribution expectations, and explicit non-goals while final naming remains open
- **Contract wired into package metadata**: `harness/adapter-manifest.yml` now points at the public contract, and the manifest package list/bootstrap flow now carry that contract into bootstrapped repositories instead of leaving the durable boundary only in the source repo
- **Second-round debate control added**: `meetings orchestrate` now plans follow-up round obligations when a first linked cross-examination does not yet outrank the opposing claim lane, keeps the council in `cross_examination` until that reply lands, and surfaces the result through summary/doctor instead of collapsing every debated claim directly into synthesis after one route
- **Debate-plan accounting widened**: meeting summaries and doctor now track `planned_round_count`, `resolved_round_count`, and `open_round_action_count`, while teammate smoke validates a bounded two-round council path before synthesis
- **Artifact-family authority rules widened**: `runs authority-audit` now treats stale contract snapshots and live git drift as explicit winner/loser conflicts, with `artifact_authority_conflicts` explaining when active run state beats contract state and when live git state beats stale run records
- **Core smoke widened again**: `scripts/harness-smoke-test.sh` now proves a manually-staled contract branch loses to active run state before continuing into the existing live-branch/live-dirty authority path
- **Public prototype naming accepted**: ADR [`0013-harness-prototype-positioning.md`](/Users/liyasong/Project/collov/erwin/_project/decisions/0013-harness-prototype-positioning.md) now formalizes `Harness` as the public prototype name for the Codex-native adapter line
- **Manifest/package identity aligned**: [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml) now carries `public_name`, `release_stage`, and `positioning_adr`, while the public contract now references that naming decision explicitly
- **Distribution strategy formalized**: [`docs/harness-distribution-strategy.md`](/Users/liyasong/Project/collov/erwin/docs/harness-distribution-strategy.md) now defines source bootstrap, packaged-binary, and portable runtime-bundle tiers for `Harness`
- **Runtime bundle packaging added**: [`scripts/harness-package-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-package-runtime.sh) now builds per-target tar.gz bundles containing `harness-state`, the adapter manifest, and the public contract, and [`scripts/harness-runtime-package-smoke-test.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-runtime-package-smoke-test.sh) verifies that path
- **Mission-level authority adjudication added**: missions now record both `latest_run_id` and `authoritative_run_id`, so the control plane can acknowledge fresh teammate activity without letting it silently replace the leader-governed mission run
- **Mission authority audit added**: [`scripts/harness-missions.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-missions.sh) now exposes `authority-audit`, explaining mission-record drift, leader-over-teammate mission control, and unresolved parallel root-run conflicts
- **Mission migration path widened**: `maintenance normalize` now regenerates mission state when older mission records are missing the new authoritative control field, so upgraded repos do not stay permanently noisy in `doctor`
- **Smoke coverage widened again**: core smoke now validates clean mission authority baselines, and teammate smoke now proves that a teammate can become the freshest run while the leader run still correctly governs mission control
- **Multi-mission authority adjudication added**: [`scripts/harness-missions.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-missions.sh) now exposes `authority-audit --all`, surfacing cross-mission ownership collisions and dependency-governed mission blocking as first-class control-plane verdicts
- **Doctor/home now consume graph authority**: `doctor` surfaces `mission-graph-authority-conflict`, and `home` now reports graph-level mission conflict and dependency-adjudication counts instead of hiding multi-mission control drift under per-run overlap noise
- **Core smoke widened again**: core smoke now proves a cross-mission overlap can be created under explicit allowance, surfaced by `missions authority-audit --all`, reported by `doctor`, and then cleared by closing the conflicting mission
- **Shared ownership transfer added**: `runs transfer-ownership` now supports `--mode=copy`, allowing bounded selectors to be intentionally shared between active runs when both sides explicitly allow `shared` overlap instead of forcing every transfer into exclusive handoff semantics
- **Teammate smoke widened again**: teammate smoke now proves shared-selector copy transfer end-to-end, including source retention, target acquisition, and zero contract drift on both runs
- **Dependency-governed mission ownership adjudication added**: `missions authority-audit --all` now distinguishes unresolved parallel mission overlap from overlap that is explicitly governed by a mission dependency, producing a real graph-level winner/loser verdict instead of leaving every bounded overlap as generic conflict
- **Doctor noise tightened**: inactive missions with historical dependencies no longer trigger `mission-dependency-state-unusual`, so graph-level authority work can close cleanly back to healthy runtime state
- **Bundle-aware mission graph authority added**: `missions authority-audit --all` now groups connected active missions into bundle-aware graph views, distinguishing conflicted bundles from governed bundles instead of leaving the leader to infer multi-mission structure from pairwise verdicts
- **Home widened again**: `home` now surfaces bundle-level mission graph counts alongside conflict and dependency counts, so the leader view can see whether the active graph is fragmented, conflicted, or already governed
- **Bundle governor added**: governed mission bundles now expose `governing_mission_id` when a single root mission exists, so the control plane can name who governs the bundle instead of making the leader reverse-engineer it from dependency edges
- **Ownership and desired state strengthened**: runs and contracts now carry explicit `ownership_claims`, `desired_state`, and dependency policy; teammate assignment now rejects exclusive overlap by default, and `drive wait` plus `home`/`doctor` surface dependency-blocked runs directly
- **Authority and council control advanced**: invalid authority levels are now rejected or surfaced, and meetings can advance phases with `meetings orchestrate` instead of relying entirely on manual phase jumps
- **Mission graph made first-class**: the harness now materializes durable mission objects under `.harness/missions/`, exposes them through `scripts/harness-missions.sh`, and drives `home`/`doctor` from mission-level state instead of only run-level state
- **Legacy runtime migration path added**: `scripts/harness-maintenance.sh normalize` now backfills ownership claims, desired state, dependency policy, and mission registry state for older pre-mission runs, bringing existing `.harness/` state forward without manual surgery
- **Selector-aware ownership transfer added**: runs now normalize legacy `ownership_claims` into structured `ownership_selectors`, `runs transfer-ownership` can formally move bounded selectors between active runs in the same mission, and transfer events now survive into mission-level state
- **Signal authority adjudication added**: superseded active signals are now detected against fresher run/meeting state, surfaced by `doctor`, and automatically resolved by `maintenance reconcile` or explicitly via `signals adjudicate`
- **Council orchestration hardened**: meetings now keep durable `phase_history`, structured `orchestration` state, linked cross-examination metadata, and synthesis audits with blocking gaps, coverage notes, and evidence gaps instead of relying on loose phase labels alone
- **Adapter line advanced**: `harness/adapter-manifest.yml` now declares `0.8.0-prototype` to reflect the stronger council-governance contract
- **Meeting summaries upgraded**: `meetings summary` now surfaces debate round, phase-history count, synthesis readiness, and blocker counts so leader views can see whether a council is genuinely ready for adoption
- **Doctor now audits council rigor**: active meetings are checked for missing phase history, unlinked cross-examination, missing synthesis audit, output reached before synthesis readiness, and unresolved findings with no follow-up action
- **Teammate smoke now covers governed debate linkage**: the council path now validates explicit cross-examination targeting plus synthesis-readiness reporting in addition to the previous participant/position lifecycle
- **Meeting authority audit added**: meetings now detect same-claim finding conflicts, explain when a higher-authority finding wins, and surface unresolved same-authority conflicts through `meetings summary`, `meetings authority-audit`, and synthesis blocking rules
- **Richer selector semantics added**: ownership selectors now carry family/match/bounded metadata, normalize aliases like `directory -> dir`, and recognize semantic path overlap such as `dir:docs/release` conflicting with `file:docs/release/CHANGELOG.md`
- **Semantic overlap enforcement added**: teammate assignment, claim-time admission, and doctor/home overlap accounting now operate on selector semantics instead of only exact string equality
- **Run-level authority adjudication added**: `scripts/harness-runs.sh authority-audit` now explains higher-authority active-signal conflicts, meeting-outcome adoption gaps, and desired-state drift in one leader-facing view
- **Doctor widened to run-level authority conflicts**: `harness-doctor.sh` now emits `run-authority-conflict` when fresher or higher-authority active signals disagree with the current run state
- **Core smoke widened**: `scripts/harness-smoke-test.sh` now validates a full higher-authority blocked-signal scenario through audit, doctor surfacing, resolution, meeting re-adoption, and health recovery
- **Contract authority sync added**: run control-plane mutations now refresh contract snapshots for desired state, ownership selectors, branch, and worktree fields so planning artifacts do not silently drift behind active runtime state
- **Normalize widened to contracts**: `scripts/harness-maintenance.sh normalize` now repairs stale historical contracts in addition to runs/signals, bringing legacy `.harness/contracts/` forward to the current control-plane contract
- **Smoke widened again**: core and teammate smoke now explicitly prove zero contract drift for claimed runs and ownership-transfer flows
- **Live-state drift detection proved**: core smoke now also demonstrates a real run-record vs worktree-branch mismatch, and `runs authority-audit` plus `doctor` both surface the live git drift until the run record is repaired
- **Meeting synthesis trace added**: findings can now explicitly reference supporting and challenging positions, and `meetings summary` surfaces traced vs untraced finding counts so the leader can see whether conclusions still lack provenance
- **Automated council synthesis added**: `meetings orchestrate` can now derive provisional findings and unresolved follow-up actions directly from bounded positions grouped by `claim` / `dimension`, while leaving manually recorded findings intact
- **Position claim grouping added**: `meetings add-position` now accepts `--claim=` and `--dimension=` so councils can treat multiple positions as speaking about the same governed claim
- **Meeting summary widened again**: `meetings summary` now reports auto-generated finding/action counts alongside traced-finding and authority-conflict reporting
- **Live intent authority added**: `runAuthorityAudit` now flags a dirty live worktree when run/desired state already claims the execution should be clean enough for governed verification or handoff
- **Harness-generated worktree noise excluded**: worktree authority snapshots now ignore the harness-owned `.mission` marker so live dirty detection tracks real execution drift instead of harness bookkeeping
- **Core smoke widened again**: the live-state suite now proves both `record drift` and `intent drift`, including the case where desired clean completion intent conflicts with a dirty live worktree
- **Debate routing plan added**: meetings now materialize a `debate_plan` that groups conflicting claim lanes, tracks routed vs unrouted cross-examination, and makes missing challenge paths explicit in the control plane
- **Meeting summary widened again**: `meetings summary` now reports routed and unrouted conflicting claim counts in addition to synthesis and automation metrics
- **Doctor widened again**: `doctor` now flags meetings whose conflicting claims have not yet been routed into linked cross-examination before synthesis/output is trusted
- **Adapter line advanced again**: `harness/adapter-manifest.yml` now declares `0.18.0-prototype` to reflect stronger run-level gate authority adjudication
- **Adapter/package flow updated**: `harness/adapter-manifest.yml` now packages the Go module files and build runner, bootstrapped repos receive the Go sources, and `harness-bootstrap.sh` / `harness-sync.sh` now query adapter metadata through the Go runtime path
- **Smoke suite de-Pythonized**: issue/core/release/teammate smoke scripts now query package files through the harness adapter surface and extract runtime ids without Python helpers
- **Validation**: `./scripts/harness-build-runtime.sh`, `./scripts/harness-runtime-dist-smoke-test.sh`, `./scripts/harness-adapter.sh validate`, `./scripts/harness-issue-smoke-test.sh`, `./scripts/harness-smoke-test.sh`, `./scripts/harness-teammate-smoke-test.sh`, `./scripts/harness-release-smoke-test.sh`, `./scripts/harness-bootstrap-smoke-test.sh`, and `./scripts/harness-all-smoke-test.sh` all pass cleanly on the Go runtime; `./scripts/harness-maintenance.sh normalize` also upgrades older runtime state into mission-aware schema cleanly

### Previous completion

**Session C** (2026-03-28) — Harness prototype crossed from leader-centric control plane into first-class teammate runtime

- **Teammate runtime added**: `scripts/harness-teammates.sh` plus the harness runtime `teammates ...` surface now create bounded delegated runs with scope, authority refs, forbidden sources, reporting expectations, and parent/leader linkage
- **Governed lateral communication added**: teammate-to-teammate communication is now explicit and structured instead of implicit chat, using sibling-approved signal flow and leader-approved linkage
- **Council protocol strengthened**: meetings now track `phase`, `participants`, and `positions`, so council packets can represent independent review and cross-examination instead of only final findings
- **Home/doctor upgraded**: `home` now surfaces active teammate and coordination signal counts; `doctor` now validates teammate assignment quality plus meeting participant/phase/position health
- **Smoke coverage expanded**: `scripts/harness-teammate-smoke-test.sh` now validates delegation, report emission, sibling communication, meeting escalation, participant rostering, and council positions end-to-end
- **Adapter/package contract advanced**: teammate surface and teammate smoke are now part of `harness/adapter-manifest.yml`, and the reusable adapter line advanced to `0.6.0-prototype`
- **Local runtime cleaned**: legacy smoke teammate runs/signals were retired and reconciled so `scripts/harness-doctor.sh` returns cleanly again in the real repository runtime
- **Validation**: `scripts/harness-smoke-test.sh`, `scripts/harness-teammate-smoke-test.sh`, `scripts/harness-release-smoke-test.sh`, `scripts/harness-bootstrap-smoke-test.sh`, `scripts/harness-all-smoke-test.sh`, and `scripts/harness-doctor.sh` all pass cleanly

- **Manifest became executable contract**: `harness/adapter-manifest.yml` now declares protocol identity, adapter package files, bootstrap contract, and profile contract instead of only loose packaging metadata
- **Adapter surface added**: `scripts/harness-adapter.sh` now exposes manifest introspection and validation as a first-class deterministic surface
- **Bootstrap/package flow hardened**: `scripts/harness-bootstrap.sh` now copies manifest-declared adapter files instead of relying on shell hardcoded copy lists, and still generates the target repo profile from the manifest-declared bootstrap template
- **Sync status strengthened**: `scripts/harness-sync.sh status` now reports adapter identity alignment in addition to manifest version alignment
- **Doctor upgraded**: `scripts/harness-doctor.sh` now includes adapter/package validation and also fixes the misplaced meeting validation logic that previously lived under the active-RC loop
- **Smoke contract upgraded**: core/release/bootstrap smoke now validate adapter integrity explicitly, and `scripts/harness-all-smoke-test.sh` now starts with adapter validation before broader lifecycle smoke
- **Adapter version advanced**: the reusable package line is now `0.5.0-prototype`
- **Validation**: `scripts/harness-adapter.sh validate`, `scripts/harness-doctor.sh`, and `scripts/harness-all-smoke-test.sh` all pass cleanly; full smoke remains green end-to-end

## Active Task / What's In Progress

**The first complete `Harness` prototype line is now delivered: runtime, control plane, public contract, authority model, council model, naming, distribution strategy, and a first real leader-facing scheduler baseline are all formalized. Remaining work is future expansion, not core prototype completion:**

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 1 | ✅ | Research protocol / control plane / authority / cybernetics synthesis |
| Phase 2 | ✅ | Build repository-local harness definition layer (`harness/`) |
| Phase 3 | ✅ | Implement runtime substrate (`scripts/harness_*`) |
| Phase 4 | ✅ | Validate `home` / `doctor` / `claim` / `drive` / `meetings` / `maintenance` |
| Phase 5 | ✅ | Package the harness formally as a reusable Codex-native adapter/system |
| Phase 6 | ✅ | Split reusable adapter/bootstrap flow from repo-local proving ground |
| Phase 7 | 🔄 | Future expansion beyond the first complete prototype line: broader authority/adaptation layers and optional Codex-native multi-agent integration |

## Next Session Should Start With

1. Read this file
2. **Harness next work**:
   - A. Continue authority adjudication beyond the current run/mission/mission-graph plus dependency-governed ownership, selector-specificity winners, council-governed bundle authority, and bundle-governor rules into broader multi-mission bundle handling when no obvious governing dependency or single active council anchor exists
   - B. Continue pushing council behavior from auto-synthesis, direct synthesis for aligned/authority-resolved lanes, stateful current/blocked debate batching, wave obligations, route obligations, and bounded multi-round challenge planning toward fuller `arch-review`-style automated debate execution across more than one conflicting lane
   - C. Continue scheduler beyond the current snapshot/execution-mode plus bundle-aware deferral baseline into larger-portfolio admission, rescheduling, and execution-portfolio policy
   - D. Decide whether future Codex-native subagent integration should map onto the teammate runtime directly or remain an optional adapter layer
   - H. If starting a fresh session, read `AGENTS.md`, `_project/CURRENT.md`, `docs/harness-documentation-map.md`, `docs/harness-architecture.md`, `docs/harness-operating-model.md`, `docs/harness-public-contract.md`, `docs/harness-session-handoff.md`, `harness/adapter-manifest.yml`, `docs/research/harness-gap-audit-2026-03-29.md`, `docs/research/harness-research-drift-audit-2026-03-29.md`, recent `_project/sessions/2026-03-29-*.md` logs, and inspect recent commits on `codex/harness-release-lifecycle`
3. **Showcase tasks still pending**: DesignPatternLibrary, CaseStudies — blocking v1.0.0

## Implementation Roadmap

### v0.5.0 — V2 Phase 2: Deep Interaction (in progress)
- [x] OriginTrace (d3-force philosophy graph, ACT II)
- [x] Bundle code splitting (627KB → 221KB, React.lazy + vendor chunks)
- [x] Code health sweep (ESLint clean, dead code removal)
- [x] ExecutionSimulator: message flow animation enhancement
- [x] SwarmTopologySection: add interactive message flow

### v1.0.0 — V2 Phase 3: Creation Path
- [x] SkillBuilder (form → skill.md download, ACT IV)
- [ ] DesignPatternLibrary (3-5 reusable pattern cards)
- [ ] CaseStudies (2-3 skill variant examples)

### Skills Factory
- [x] Teamwork v2.0–v3.11.2 — 40+ versions, 6 skills, 10 scripts, 48 unit tests
- [x] Wing v0.1.0 — Anthropic skill architecture rewrite (6 SKILL.md, 10 operations, 8 scripts)
- [x] skill-lint.sh — Layer 1 static analysis (12 checks)
- [x] arch-review v0.1.0 — Multi-role debate review (8 mechanisms, 27 bugs fixed in 5 rounds)
- [x] Harness prototype — Codex-native repository-local orchestration harness (leader control plane + runtime substrate)
- [ ] Wing Phase 4: vi_agent migration
- [ ] Wing Phase 5: teamwork deprecation
- [ ] Skill-tester Layer 2: smoke tests

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Wing scripts reference `${CLAUDE_SKILL_DIR}` which requires Anthropic skill infrastructure
- skill-lint.sh ### header check has false positives (structural headings vs output templates)
- Harness is now installed and validated in Codex via `install-codex.sh`, can bootstrap into another repo, has manifest-driven adapter validation/packaging, first-class mission state, a leader-facing scheduler baseline with durable snapshots, execution-mode selection, and direct bundle-aware deferral/escalation, mission-level and graph-level authority adjudication including dependency-governed ownership verdicts, selector-specificity winners, council-governed bundle authority, bundle-aware grouping, and explicit bundle governors, selector-aware ownership transfer with shared-copy and exclusion-aware split modes, artifact-aware authority adjudication, teammate runtime, council orchestration with phase history + synthesis audit semantics, and a formalized runtime distribution strategy; remaining work is expansion beyond the first complete prototype line

## File Locations

| What | Where |
|------|-------|
| Project goals | `_project/GOALS.md` |
| Idea garden | `_project/ideas/` |
| ADRs | `_project/decisions/` |
| V2 architecture | `docs/drive-showcase-v2-architecture.md` |
| Wing skills | `skills/wing/` |
| Harness prototype | `harness/`, `scripts/harness_*`, `skills/harness/` |
| Wing architecture | `_project/ideas/skills/tw2-architecture.md` |
| Teamwork skills (legacy) | `skills/teamwork/` |
| Docs index | `docs/README.md` |
| App root | `showcase/src/App.tsx` |
