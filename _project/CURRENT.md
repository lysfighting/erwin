# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 1.0.0
- **Last Updated**: 2026-04-19
- **Last Session**: 2026-04-19-A (cross-agent continuity guide added)

## Current Phase

**Phase**: AI/agent frontier exploration lab with active Radar, Knowledge, and Codex-native Harness tracks

## What Was Just Completed

**Session 2026-04-19-A** — Cross-agent continuity guide added

- **Cross-agent restart is now a first-class contract**: `docs/agent-continuity-guide.md` defines the vendor-neutral start protocol, task routing, trust order, Harness evolution rule, knowledge-network contract, and end protocol for Codex, Claude, and future agents
- **Claude no longer starts from stale repo cognition**: `CLAUDE.md` now matches the five-layer `Radar / Feedback / Knowledge / Proposals / Implementation` model instead of the old three-lane summary
- **Entrypoints now converge on the same durable truth**: `README.md`, `AGENTS.md`, `docs/README.md`, `docs/erwin-operating-model.md`, and `docs/harness-session-handoff.md` all route fresh work through the continuity guide before local subsystem docs
- **The knowledge network now indexes continuity artifacts**: `scripts/build-knowledge-surface.py` includes `docs/agent-continuity-guide.md` and `CLAUDE.md` in the generated cross-layer graph
- **Next sessions should not depend on old Codex or Claude thread memory**: start from the agent-specific entrypoint, then `_project/CURRENT.md`, `docs/agent-continuity-guide.md`, `docs/erwin-operating-model.md`, and `docs/README.md`

**Session 2026-04-17-C** — Knowledge network index added

- **The five-layer system now has a machine-readable graph snapshot**: `docs/knowledge/index.json` includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`
- **The generated surface now shows network health**: `docs/index.html` and `docs/knowledge/briefing.md` display artifact, link, thread, and proposal-record counts so the system is visible without reading raw JSON first
- **The graph stays repo-owned and Markdown-backed**: `scripts/build-knowledge-surface.py` derives nodes from Radar, Feedback, Knowledge, Proposals, and selected Implementation entrypoints; Markdown remains the source of truth
- **Operating docs now name the contract**: `docs/erwin-operating-model.md`, `docs/knowledge/README.md`, and `docs/README.md` document the generated graph as the handoff interface for future agents and cockpit UI work
- **Next sessions should consume the index instead of reconstructing cross-layer state from chat memory**: Knowledge Surface Cockpit work can read `docs/knowledge/index.json` as its local graph source, while Radar and Feedback sessions should rebuild it after changing stream or case files

**Session 2026-04-17-B** — Five-layer operating model added

- **Repo-level cognition is now explicit**: `docs/erwin-operating-model.md` defines the five-layer loop `Radar + Feedback -> Knowledge -> Proposals -> Implementation -> Feedback / Knowledge`
- **Harness is now framed in the loop correctly**: Harness is the main implementation target and feedback source, not a separate shortcut around Radar, Feedback, Knowledge, and proposal review
- **Fresh sessions now have a global route before local docs**: `AGENTS.md`, `docs/README.md`, and `docs/harness-documentation-map.md` now point cross-lane work at the five-layer model before narrower Harness or Knowledge documents
- **The knowledge thread map now tracks coordination itself**: `docs/knowledge/active-threads.md` carries `erwin five-layer operating model` as an active thread so future multi-session drift can flow back into durable knowledge
- **Post-merge release gate was rerun on mainline truth**: after PR #10 landed the Harness leader-briefing line on `main`, `go test ./...` and `./scripts/harness-all-smoke-test.sh` both passed from a fresh `origin/main` worktree

**Session 2026-04-17-A** — Hardened release briefings and merge readiness

- **Release status briefings now fail safe on gate drift**: `ship status --brief` now lets verification, review, checks, and merge readiness override stale `next_expected_action` text when ship state is optimistic but gates disagree
- **Release verification now proves the preserved JSON contract**: release smoke asserts default JSON for `ship status`, `rc list`, and `rc get`, then proves the human briefing blocks merge guidance in an inconsistent gate state
- **Packaged-runtime coverage now includes active RC detail**: runtime-distribution smoke creates an RC without Go in `PATH`, checks active `rc list --brief`, and checks `rc get RC_ID --brief`
- **Canonical surfaces now match runtime**: `docs/harness-session-handoff.md` lists the full current human-direct `--brief` stack, and `docs/harness-public-contract.md` includes `runs` and `signals` among first-class surfaces
- **Branch integration was advanced**: `codex/harness-100` was merged into `codex/harness-100-leader-briefing`, with journal conflicts resolved toward the newer human-briefing truth while preserving the explainer session history from the base Harness line

**Session 2026-04-16-K** — Extended the human-facing leader layer into `ship status --brief` and `rc list/get --brief`

- **Release state now has real human-facing status entrypoints**: `bash scripts/harness-ship.sh status RUN_ID --brief`, `bash scripts/harness-rc.sh list --brief`, and `bash scripts/harness-rc.sh get RC_ID --brief` now turn ship gates and RC state into readable handoffs that name verification, review, checks, merge readiness, tracked PR/RC refs, RC branch/version/run membership, staged verification, and next release action
- **Release briefings stay on the same Go control plane**: the new layer is derived from existing run ship state and RC records; default ship status and RC JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary release status**: release smoke proves source-runtime ship/RC status briefings, and runtime-distribution smoke proves packaged binaries can render the same status layer without Go in `PATH`

**Session 2026-04-16-J** — Extended the human-facing leader layer into `teammates link --brief`

- **Governed teammate link creation now has a real human-facing leader entrypoint**: `bash scripts/harness-teammates.sh link RUN_A RUN_B [REASON] --brief` now turns peer-link approval into a readable handoff that names both teammates, connected scopes, updated communication peers, leader-approved reason, and the boundary that permits governed messages without transferring ownership or mission authority
- **Link briefings stay on the same Go control plane**: the new layer is derived after both run records are updated with sibling and communication peer state; default `teammates link` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary governed links**: teammate smoke proves source-runtime peer-link briefings, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-16-I** — Extended the human-facing leader layer into `meeting-close` / `adopt-meeting --brief`

- **Council adoption now has real human-facing leader entrypoints**: `bash scripts/harness-drive.sh meeting-close RUN_ID MEETING_ID [NOTE] --brief` and `bash scripts/harness-drive.sh adopt-meeting RUN_ID MEETING_ID --brief` now turn council closure or re-adoption into readable handoffs that name the council, run, outcome, output packet kind, follow-up counts, authority adoption, and next leader/council boundary
- **Adoption briefings stay on the same Go control plane**: the new layer is derived after closed meeting state, backlog updates, authority adoption, and run adoption fields are persisted; default `meeting-close` and `adopt-meeting` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary council adoption**: core smoke proves source-runtime close and re-adoption briefings, and runtime-distribution smoke proves the same surfaces work through the packaged binary without Go in `PATH`

**Session 2026-04-16-H** — Extended the human-facing leader layer into `drive meeting --brief`

- **Council creation now has a real human-facing leader entrypoint**: `bash scripts/harness-drive.sh meeting RUN_ID QUESTION --brief` now turns the same meeting materialization mutation into a readable handoff that names the created council, host run, mission, selected decision protocol, backlog ref, and exact `meetings summary/orchestrate --brief` follow-up path
- **Meeting-open briefings stay on the same Go control plane**: the new layer is derived after the active meeting, maintenance backlog entry, and waiting run state are persisted; default `drive meeting` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary council creation**: core smoke proves source-runtime meeting-open briefings with explicit protocol selection, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-16-G** — Extended the human-facing leader layer into `teammates escalate --brief`

- **Teammate escalation now has a real human-facing leader entrypoint**: `bash scripts/harness-teammates.sh escalate TEAMMATE_RUN_ID QUESTION --brief` now turns the same teammate escalation mutation into a readable handoff that names the escalating teammate, emitted `requires_meeting` signal, optional materialized council, selected decision protocol, question, and next leader/teammate/council boundary
- **Escalation briefings stay on the same Go control plane**: the new layer is derived after the escalation signal and optional council materialization are persisted; default `teammates escalate` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary teammate escalation**: teammate smoke proves source-runtime materialized escalation briefings with explicit protocol selection, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-16-F** — Extended the human-facing leader layer into `teammates message --brief`

- **Governed teammate messaging now has a real human-facing leader entrypoint**: `bash scripts/harness-teammates.sh message SOURCE_RUN_ID TARGET_RUN_ID KIND SUMMARY --brief` now turns the same governed peer-communication mutation into a readable handoff that names source and target teammates, emitted signal, summary, evidence, freshness, and next leader/source/target/council boundary
- **Message briefings stay on the same Go control plane**: the new layer is derived after source/target run signal linkage and the sibling-scope signal are persisted; default `teammates message` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary teammate messaging**: teammate smoke proves source-runtime governed communication briefings, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-16-E** — Extended the human-facing leader layer into `teammates report --brief`

- **Teammate reporting now has a real human-facing leader entrypoint**: `bash scripts/harness-teammates.sh report TEAMMATE_RUN_ID SUMMARY --brief` now turns the same teammate report mutation into a readable handoff that names the teammate run, emitted report signal, summary, evidence, freshness, and next leader/teammate/council boundary
- **Report briefings stay on the same Go control plane**: the new layer is derived after the teammate run freshness/summary and active signal linkage are persisted; default `teammates report` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary teammate reports**: teammate smoke proves source-runtime progress report briefings, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-16-D** — Extended the human-facing leader layer into `teammates assign --brief`

- **Teammate assignment now has a real human-facing leader entrypoint**: `bash scripts/harness-teammates.sh assign PARENT_RUN_ID TEAMMATE_ID ROLE [OBJECTIVE] --brief` now turns the same bounded delegation mutation into a readable handoff that names the created teammate run, parent run, contract, ownership selectors, dependencies, governed communication peers, assignment signal, and next leader/teammate/council boundary
- **Assignment briefings stay on the same Go control plane**: the new layer is derived after the teammate run, parent run update, contract, and assignment signal are persisted; default `teammates assign` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary teammate delegation creation**: teammate smoke proves ordinary assignment briefings from source runtime, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-16-C** — Extended the human-facing leader layer into `drive resolve-follow-up --brief`

- **Adopted follow-up resolution now has a real human-facing leader entrypoint**: `bash scripts/harness-drive.sh resolve-follow-up RUN_ID [FOLLOW_UP_PACKET_ID] [STATE] --brief` now turns the same follow-up mutation into a readable handoff that names the invoking run, governing run, packet, target state, delegated-resolution status, and next leader/teammate/council boundary
- **Follow-up resolution briefings stay on the same Go control plane**: the new layer is derived after the governing run, delegated teammate run, and meeting action state have been synchronized; default `drive resolve-follow-up` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary adopted follow-up closure**: teammate smoke proves governing-run reopen and delegated teammate closure briefings, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-16-B** — Extended the human-facing leader layer into `drive report --brief`

- **Progress reporting now has a real human-facing leader entrypoint**: `bash scripts/harness-drive.sh report RUN_ID SUMMARY --brief` now turns the same report mutation into a readable handoff that names the run, mission, refreshed state, reported progress, and next leader/teammate/council boundary
- **Report briefings stay on the same Go control plane**: the new layer is derived from the updated run record after `latest_summary`, `last_reported_at`, and `updated_at` are persisted; default `drive report` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary progress reporting**: core smoke proves report brief output on live runtime state, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-16-A** — Extended the human-facing leader layer into `claim --brief` and core `drive ... --brief` action handoffs

- **Claiming work now has a real human-facing leader entrypoint**: `bash scripts/harness-claim.sh ISSUE_NUMBER [TITLE_OVERRIDE] --brief` now turns the same claim mutation into a readable handoff that names the created mission/run/branch, provisioning result, contract path, and post-claim leader/teammate/council boundary
- **Core drive actions now explain their state mutations to humans**: `bash scripts/harness-drive.sh start|resume|wait|block|complete ... --brief` now renders the run state transition, current summary, next expected action, teammate boundary, council boundary, and emitted blocker signal when `block` creates one, while default drive JSON remains unchanged
- **Human-facing verification now covers source and packaged-binary action handoffs**: core smoke proves claim/start/wait/block/resume/complete briefings on live runtime state, and runtime-distribution smoke proves the same surfaces work through the packaged binary without Go in `PATH`

**Session 2026-04-15-B** — Extended the human-facing leader layer into `meetings resolve-action --brief`

- **Council action closure now has a real human-facing leader entrypoint**: `bash scripts/harness-meetings.sh resolve-action MEETING_ID ACTION_ID [STATE] --brief` now turns the same action-resolution mutation into a readable briefing that tells a human which action was closed, cancelled, or reopened, what the remaining action backlog looks like, and what leader/teammate/council action comes next
- **Meeting action-resolution briefings stay on the same Go control plane**: the new layer is derived from updated `action_items`, refreshed meeting output packets, and the same adopted follow-up lifecycle sync that already updates governing runs; default `meetings resolve-action` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary action lifecycle paths**: core smoke now proves done and cancelled action resolution explain adopted follow-up synchronization, while runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-15-A** — Extended the human-facing leader layer into `meetings add-action --brief`

- **Council action recording now has a real human-facing leader entrypoint**: `bash scripts/harness-meetings.sh add-action MEETING_ID TEXT --brief` now turns the same action-recording mutation into a readable briefing that tells a human which council action item was added, what the open action backlog looks like, who owns the next step, and whether teammate assignment or adoption is still needed
- **Meeting action briefings stay on the same Go control plane**: the new layer is derived from the existing `action_items` meeting state plus meeting summary and governance-gap calculations rather than inventing a separate follow-up model; default `meetings add-action` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary council-action paths**: teammate smoke now proves `meetings add-action --brief` explains that teammate execution remains leader-owned, while runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-14-B** — Extended the human-facing leader layer into `meetings add-finding --brief`

- **Council judgment recording now has a real human-facing leader entrypoint**: `bash scripts/harness-meetings.sh add-finding MEETING_ID accepted|rejected|unresolved TEXT --brief` now turns the same finding-recording mutation into a readable briefing that tells a human which council claim/dimension was judged, what kind of finding was recorded, what authority/evidence/position trace supports it, and what leader/teammate/council action comes next
- **Meeting finding briefings stay on the same Go control plane**: the new layer is derived from the existing `findings` meeting state plus same-claim authority and governance-gap calculations rather than inventing a separate council-judgment model; default `meetings add-finding` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary council-judgment paths**: teammate smoke now proves `meetings add-finding --brief` explains rejected-finding governance repair before output, while runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-14-A** — Extended the human-facing leader layer into `meetings add-position --brief`

- **Council position recording now has a real human-facing leader entrypoint**: `bash scripts/harness-meetings.sh add-position MEETING_ID RUN_ID PHASE STANCE TEXT --brief` now turns the same position-recording mutation into a readable briefing that tells a human which council claim/dimension was affected, who contributed it, whether more independent review or debate is still owed, and what leader/teammate/council action comes next
- **Meeting position briefings stay on the same Go control plane**: the new layer is derived from the existing `positions` meeting state plus live coverage/wave/debate calculations rather than inventing a separate council-progress model; default `meetings add-position` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary council-position paths**: teammate smoke now proves `meetings add-position --brief` explains remaining independent-review coverage for a live council, while runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-17** — Knowledge Surface Research Cockpit v1 proposal landed

- **Knowledge surface upgrade is now a formal proposal**: `docs/proposals/knowledge-surface-research-cockpit-v1.md` records the decision to evolve `docs/index.html` from a static briefing into a local Research Cockpit backed by Vite/React/TypeScript and a small Go read-only service
- **Design benchmark evidence is captured in the repo**: the proposal records the specific patterns to borrow from NotebookLM, Readwise Reader, Linear, Raycast, Heptabase, Obsidian Canvas, Capacities, and Tana without turning v1 into a database, SaaS, or browser AI chat product
- **Next implementation step is bounded**: future work should open a dedicated implementation worktree, strengthen `docs/knowledge/index.json`, add the Go API, add the frontend surface, and keep `docs/index.html` as a fallback launch/briefing page
- **No UI prototype was implemented in this session**: this session intentionally landed the proposal artifact and navigation only, so the next session can build from a stable repo-backed plan instead of transient thread context

**Session 2026-04-14** — Daily and trending radar stream entrypoints added

- **Daily radar now has a directory-level contract**: `docs/radar/daily/README.md` documents dated files, `latest.md`, 7-day windows, first-party source priority, secondary ecosystem evidence, and knowledge-surface refresh expectations
- **GitHub Trending now has a directory-level contract**: `docs/radar/trending/README.md` documents default global Trending capture, direct repo links, capture dates, board-shape preservation, and Tier 2 interpretation boundaries
- **Radar navigation is less implicit**: `docs/README.md` and `docs/radar/README.md` now link directly to both stream entrypoints instead of leaving daily/trending directories as unexplained archives
- **No new session-log file was created**: this update is recorded in `CURRENT.md` and `CHANGELOG.md` to keep project history useful without adding low-value session clutter

**Session 2026-04-14** — cwork skill harvested into repo without merging stale branch history

- **cwork is now a repo-owned Codex skill**: `skills/cwork/` was copied from the historical `codex/cwork-skill` branch into the current repo line, matching the already-installed `~/.codex/skills/cwork` content
- **Old branch was not merged wholesale**: only the `skills/cwork/` directory was harvested, avoiding the stale branch's large unrelated deletions across docs, knowledge, radar, runtime, and old session files
- **Skill index and installer examples now mention cwork**: `README.md`, `AGENTS.md`, and `skills/README.md` list `/cwork` next to `/harness` and `/knowledge`, and `install-codex.sh` now documents `./install-codex.sh cwork`
- **Offline claim smoke exposed and fixed a branch override bug**: `cwork-claim.sh --branch` now keeps the contract branch and actual Git branch aligned during claim setup
- **No new session-log file was created**: this update is recorded in `CURRENT.md` and `CHANGELOG.md` to avoid adding more low-value session clutter

**Session 2026-04-13-B** — Extended the human-facing leader layer into `meetings lock-decision --brief`

- **Council scope locking now has a real human-facing leader entrypoint**: `bash scripts/harness-meetings.sh lock-decision MEETING_ID TEXT --brief` now turns the same locked-decision mutation into a readable briefing that tells a human which council decision was bounded, why that scope is authoritative, and what leader/teammate/council action comes next
- **Meeting decision-lock briefings stay on the same Go control plane**: the new layer is derived from the existing `locked_decisions` meeting state and decision-lock arguments rather than inventing a separate council-scope model; default `meetings lock-decision` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers source and packaged-binary council-scope locks**: teammate smoke now proves `meetings lock-decision --brief` explains bounded decision scope for a live council, while runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-13-A** — Extended the human-facing leader layer into `maintenance normalize --brief`

- **Runtime migration now has a real human-facing leader entrypoint**: `bash scripts/harness-maintenance.sh normalize --brief` now turns the same normalization mutation result into a readable migration briefing that tells a human which legacy run, signal, and contract records were upgraded, why no-op normalization did not change state, and what health-confirmation step comes next
- **Maintenance normalization briefings stay on the same Go control plane**: the new layer is derived from the existing normalize result and run/signal/contract migration semantics rather than inventing a parallel migration model; default `maintenance normalize` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers migrated and packaged-binary normalization paths**: core smoke now proves `maintenance normalize --brief` explains legacy run/signal migration and still leaves default JSON available afterward, while runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-12-A** — Extended the human-facing leader layer into `maintenance reconcile --brief`

- **Runtime repair now has a real human-facing leader entrypoint**: `bash scripts/harness-maintenance.sh reconcile --brief` now turns the same reconciliation mutation result into a readable repair briefing that tells a human which orphaned signals, superseded signals, and orphaned meetings were repaired, why no-op reconciliation did not change state, and what health-confirmation step comes next
- **Maintenance reconciliation briefings stay on the same Go control plane**: the new layer is derived from the existing reconcile result and signal/meeting repair semantics rather than inventing a parallel repair model; default `maintenance reconcile` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers repaired and packaged-binary reconciliation paths**: core smoke now proves `maintenance reconcile --brief` explains a superseded-signal repair and still leaves default JSON available afterward, while runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-11-D** — Extended the human-facing leader layer into `meetings orchestrate --brief`

- **Council orchestration now has a real human-facing leader entrypoint**: `bash scripts/harness-meetings.sh orchestrate MEETING_ID --brief` now turns the same orchestration mutation result into a readable council-chair briefing that tells a human whether the meeting advanced or held, why, what current wave/debate/output packet is live, and what leader/council action comes next
- **Meeting orchestration briefings stay on the same Go control plane**: the new layer is derived from the existing meeting packet, orchestration state, wave/debate plans, synthesis audit, and auto-synthesis result rather than inventing a parallel council-progress model; default `meetings orchestrate` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers advancing and packaged-binary council orchestration paths**: core smoke now proves `meetings orchestrate --brief` explains setup-to-independent-review advancement and current independent-review work, while runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-11-C** — Extended the human-facing leader layer into `signals adjudicate --brief`

- **Signal repair now has a real human-facing leader entrypoint**: `bash scripts/harness-signals.sh adjudicate SIGNAL_ID|--all --brief` now turns the same Go adjudication result into a readable signal-repair briefing that tells a human which signals were resolved as superseded, what fresher run or meeting truth won, why no-op adjudication did not change state, and what leader/teammate/council action remains
- **Signal adjudication briefings stay on the same Go control plane**: the new layer is derived from the existing signal supersession logic and signal log mutation result rather than inventing a parallel repair model; default `signals adjudicate` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers no-op, superseded-signal, and packaged-binary adjudication paths**: core smoke now proves `signals adjudicate --brief` explains both a fresh active signal that is not superseded and a blocked signal superseded by fresher run state, while runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-11-B** — Extended the human-facing leader layer into `signals active --brief`

- **Active control signals now have a real human-facing leader entrypoint**: `bash scripts/harness-signals.sh active [MISSION_ID] --brief` now turns the same active signal log, signal priority, freshness, scope, summary, and leader/teammate/council next actions into a readable signal-layer briefing instead of forcing direct raw signal JSON inspection
- **Signal briefings stay on the same Go control plane**: the new signal layer is derived from the existing signal log, run registry, signal freshness windows, and teammate ownership truth rather than inventing a parallel signal state model; default `signals active` JSON remains unchanged for exact inspection and tooling, and signal mutation/adjudication semantics remain unchanged
- **Human-facing verification now covers fresh coordination and packaged-binary teammate-signal paths**: core smoke now proves `signals active --brief` explains a fresh meeting-required signal and its leader action, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-11-A** — Extended the human-facing leader layer into `teammates get --brief`

- **Teammate objects now have a real human-facing leader entrypoint**: `bash scripts/harness-teammates.sh get TEAMMATE_RUN_ID --brief` now turns the same bounded assignment scope, parent run linkage, dependency/report freshness state, delegated follow-up provenance, governed peer communication, and leader/teammate/council next actions into a readable teammate-layer briefing instead of forcing direct raw teammate JSON inspection
- **Teammate briefings stay on the same Go control plane**: the new teammate layer is derived from the existing run registry, signal log, run authority audit, dependency truth, meeting coverage, and delegated follow-up monitor packet rather than inventing a parallel teammate state model; default `teammates get` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers ordinary, dependency-blocked, delegated follow-up, and packaged-binary teammate paths**: teammate smoke now proves `teammates get --brief` on ordinary bounded work, dependency-blocked sibling work, and adopted-follow-up handoff work, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-10-H** — Extended the human-facing leader layer into `missions graph --brief`

- **Mission topology now has a real human-facing leader entrypoint**: `bash scripts/harness-missions.sh graph --brief` now turns the same active roots, dependency edges, graph-level next action, and teammate/council coverage into a readable topology-layer briefing instead of forcing direct raw mission-graph JSON inspection
- **Mission topology briefings stay on the same Go control plane**: the new graph layer is derived from the existing mission registry topology, active run linkage, and live meeting coverage rather than inventing a parallel mission-graph summary model; default `missions graph` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers standalone, dependency-anchored, and packaged-binary graph-topology paths**: core smoke now proves both a single-root graph and a dependency-anchored graph read correctly in `missions graph --brief`, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-10-G** — Extended the human-facing leader layer into `missions authority-audit --all --brief`

- **Mission graph bundles now have a real human-facing leader entrypoint**: `bash scripts/harness-missions.sh authority-audit --all --brief` now turns the same bundle graph, conflict/escalation/governed verdicts, candidate governors, and required action into a readable graph-layer briefing instead of forcing direct raw mission-graph JSON inspection
- **Mission graph authority briefings stay on the same Go control plane**: the new graph-audit layer is derived from the existing mission graph authority audit, bundle governance notes, and candidate governor packets rather than inventing a parallel bundle-repair model; default `missions authority-audit --all` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers healthy, conflicted, escalation-required, and packaged-binary graph-authority paths**: core smoke now proves both a clean graph and two non-healthy bundle states read correctly in `missions authority-audit --all --brief`, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-10-F** — Extended the human-facing leader layer into `meetings authority-audit --brief`

- **Council same-claim authority now has a real human-facing leader entrypoint**: `bash scripts/harness-meetings.sh authority-audit MEETING_ID --brief` now turns the same meeting authority audit, same-claim winner/tie verdicts, and next authority action into a readable council-repair briefing instead of forcing direct raw council-audit JSON inspection
- **Meeting authority briefings stay on the same Go control plane**: the new council-audit layer is derived from the existing meeting authority audit, meeting packet state, and meeting summary packet rather than inventing a parallel council-repair model; default `meetings authority-audit` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers adjudicated and packaged-binary meeting-authority paths**: teammate smoke now proves a higher-authority winner reads correctly in `meetings authority-audit --brief`, and runtime-distribution smoke proves the same surface works through the packaged binary without Go in `PATH`

**Session 2026-04-10-E** — Extended the human-facing leader layer into `missions authority-audit --brief`

- **Mission control drift now has a real human-facing leader entrypoint**: `bash scripts/harness-missions.sh authority-audit MISSION_ID --brief` now turns the same mission-record drift, leader-over-teammate control adjudication, and competing root-run authority truth into a readable mission-repair briefing instead of forcing direct raw mission-audit JSON inspection
- **Mission authority briefings stay on the same Go control plane**: the new mission-audit layer is derived from the existing mission authority audit, mission registry state, run registry truth, and mission-level meeting activity counts rather than inventing a parallel mission-repair model; default `missions authority-audit` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers aligned, leader-over-teammate, and packaged-binary mission-authority paths**: core smoke now proves aligned mission authority brief output, teammate smoke proves the governed leader-over-teammate adjudication path, and runtime-distribution smoke proves `missions authority-audit --brief` works through the packaged binary without Go in `PATH`

**Session 2026-04-10-D** — Extended the human-facing leader layer into `runs authority-audit --brief`

- **Run authority drift now has a real human-facing leader entrypoint**: `bash scripts/harness-runs.sh authority-audit RUN_ID --brief` now turns the same higher-authority conflicts, adoption drift, live-state drift, contract drift, and desired-state gaps into a readable audit-layer briefing instead of forcing direct raw authority JSON inspection
- **Run authority briefings stay on the same Go control plane**: the new authority layer is derived from the existing run authority audit, live meeting linkage, teammate lineage, and audit evidence truth rather than inventing a parallel repair-summary model; default `runs authority-audit` JSON remains unchanged for exact inspection, tooling, and automation
- **Human-facing verification now covers desired-state-progress, higher-authority-conflict, and packaged-binary authority paths**: core smoke now proves both an in-progress desired-state authority brief and a non-healthy higher-authority-conflict brief, and runtime-distribution smoke proves `runs authority-audit --brief` works through the packaged binary without Go in `PATH`

**Session 2026-04-10-C** — Extended the human-facing leader layer into `missions schedule --brief`

- **Scheduler portfolios now have a real human-facing leader entrypoint**: `bash scripts/harness-missions.sh schedule --brief` now turns the same scheduler snapshot, primary control-plane packet, decision/maintenance queue, foreground batch admission, and reschedule backlog into a readable portfolio-layer briefing instead of forcing direct raw scheduler JSON inspection
- **Scheduler briefings stay on the same Go control plane**: the new scheduler layer is derived from the existing scheduler summary and durable snapshot truth rather than inventing a parallel scheduler summary model; default `missions schedule` JSON remains unchanged for exact inspection, automation, and materialization
- **Human-facing verification now covers foreground, active-council, and packaged-binary scheduler paths**: core smoke now proves scheduler brief output for an ordinary foreground mission and for an active-council primary control-plane packet that blocks foreground admission, and runtime-distribution smoke proves `missions schedule --brief` works through the packaged binary without Go in `PATH`

**Session 2026-04-10-B** — Extended the human-facing leader layer down to `meetings summary --brief`

- **Meeting objects now have a real human-facing leader entrypoint**: `bash scripts/harness-meetings.sh summary MEETING_ID --brief` now turns the same live council state, phase, blocking reason, current work packet, protocol-aware output state, and leader/council next actions into a readable object-layer briefing instead of forcing direct raw meeting JSON inspection
- **Meeting briefings stay on the same Go control plane**: the new meeting layer is derived from the existing meeting packet, orchestration state, authority audit, synthesis audit, and output packet truth instead of inventing a parallel council summary model; default `meetings summary` JSON remains unchanged for exact inspection and tooling
- **Human-facing verification now covers active, blocked, protocol-output, and packaged-binary meeting paths**: core smoke now proves meeting brief output for an active independent-review council, teammate smoke proves both governance-blocked and incident-command output paths, and runtime-distribution smoke proves `meetings summary --brief` works through the packaged binary without Go in `PATH`

**Session 2026-04-10-A** — Extended the human-facing leader layer down to `runs get --brief`

- **Run objects now have a real human-facing leader entrypoint**: `bash scripts/harness-runs.sh get RUN_ID --brief` now turns the same Go run state, mission-control relationship, authority/adoption drift, dependency state, teammate ownership, and council context into a readable object-layer briefing instead of forcing direct raw run JSON inspection
- **Run briefings preserve execution truth rather than replacing it**: the new run layer is derived from the existing run registry, mission registry, signal log, live meeting state, and `runs authority-audit` truth, while the default JSON surface remains unchanged for exact inspection and tooling
- **Human-facing verification now covers foreground, active-council, teammate, and packaged-binary run paths**: core smoke now proves run brief output for a normal foreground run and a run under active council handling, teammate smoke proves the bounded teammate-run path, and runtime-distribution smoke proves `runs get --brief` works through the packaged binary without Go in `PATH`

**Session 2026-04-09-B** — Extended the human-facing leader layer down to `missions get --brief`

- **Mission objects now have a real human-facing leader entrypoint**: `bash scripts/harness-missions.sh get MISSION_ID --brief` now turns the same scheduler mission-classification truth into a readable object-level briefing that tells a human the mission state, why it is blocked/deferred/selected, what the next action is, and what leader, teammate, and council roles need to do next
- **Mission briefings stay on the same Go control plane**: the new mission layer is derived from the existing mission registry, run registry, signal log, and scheduler bundle/decision state instead of inventing a parallel summary store; JSON mission output remains the default unchanged public contract
- **Human-facing verification now covers both foreground and live-council mission paths**: core smoke now proves mission brief output for a normal foreground mission and for an actively governed bundle-council mission, and runtime-distribution smoke proves `missions get --brief` still works through the packaged binary without Go in `PATH`

**Session 2026-04-09-A** — Added leader-readable `home` / `doctor` briefings on top of the stable Harness control plane

- **Human-facing collaboration now has a real first slice**: `harness-home.sh --brief` now turns the existing scheduler / mission / maintenance / council runtime truth into a leader-readable briefing that tells a human what is primary now, why it is primary, what to do next, what teammates need, and whether council work is active or required, instead of forcing direct raw-JSON reading for the first control-plane pass
- **Diagnostics now have the same human-readable entry path**: `harness-doctor.sh --brief` now summarizes runtime health, next repair action, and the top drift issues in readable text while keeping the default machine-readable JSON contract intact for exact audit/tooling use
- **The human layer is still derived from the same control-plane truth**: this did not add a parallel state system or move logic out of Go; smoke coverage now proves the briefing layer follows idle state, stale-maintenance escalation, active bundle-council guidance, and packaged-binary runtime execution from the real runtime payloads

**Session 2026-04-08-B** — Aligned fresh-session and skills entrypoints with the new Harness explainer

- **Fresh-session start order now leads with the explainer**: `docs/harness-session-handoff.md` now treats `docs/harness-explainer.md` as the fastest human-readable orientation layer before the deeper architecture and package docs, and it no longer points readers at an obsolete release branch when asking them to inspect recent commits
- **Top-level skills index now matches current Harness truth**: `skills/README.md` now describes `Harness` as a stable `1.0.0` leader-led orchestration surface instead of an older prototype, and it points readers to the explainer directly from the skills entrypoint
- **Explainer discoverability is now cleaner across human entrypoints**: the repo now has a more consistent path from `README.md` -> `docs/harness-explainer.md` -> `docs/harness-session-handoff.md` / `skills/README.md`, so readers no longer need to infer whether they should start from a runtime-heavy handoff or an outdated skill index

**Session 2026-04-08-A** — Added a durable human-readable Harness explainer and linked it from the main entrypoints

- **Readable overview is now part of the canonical doc system**: `docs/harness-explainer.md` now explains what `Harness` is, what problem it solves, how it is designed, how the control plane works, and how to use it, instead of forcing readers to infer those answers from architecture and runtime docs alone
- **The explainer is genuinely visual and durable**: the new doc includes Mermaid diagrams plus simplified ASCII fallbacks for the identity/layer model, operating model, and usage flow, so it remains readable both in rich Markdown renderers and plain text contexts
- **Harness entrypoints now route people to the right place**: `README.md`, `docs/README.md`, `docs/harness-documentation-map.md`, `harness/README.md`, and `skills/harness/README.md` now all point to the explainer as the best first stop for humans who want to understand Harness before diving into the deeper canonical docs

**Session 2026-04-12** — Selective knowledge-refactor artifacts absorbed into the main knowledge system

- **Old knowledge-refactor work was harvested instead of merged wholesale**: useful feedback cases, active knowledge threads, paper-watch notes, and proposal review gates were copied into `main` without accepting the old branch's stale generated-surface model or deleting current daily/trending/latest files
- **Feedback cases are now first-class knowledge inputs**: `docs/feedback/cases/` now carries normalized `HU-001` through `HU-007` cases for research-directory drift, premature proposal pressure, session-memory dependence, productized-surface overreach, delegated follow-up lifecycle gaps, council visibility gaps, and benchmark eval-trust interpretation
- **Knowledge now has a live thread registry**: `docs/knowledge/active-threads.md` connects durable threads, evidence classes, proposal links, and update triggers so future sessions do not need to reconstruct knowledge state from archived Codex threads
- **Paper watch and proposal gates are now discoverable**: `docs/radar/papers/`, `docs/radar/weekly/README.md`, and `docs/proposals/` now expose paper-aware intake, weekly synthesis gates, benchmark-integrity checks, monitorability checks, worked proposal reviews, and eval-candidate follow-through
- **The coding-agent self-iteration note moved into the durable knowledge layer**: the former research note now lives at `docs/knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md`, matching its actual role as a cross-cutting active knowledge anchor
- **The knowledge skill now follows the same structure**: `skills/knowledge/` now points future sessions at `active-threads`, proposal-layer ledgers, and review gates instead of the older top-level proposal ledger path

**Session 2026-04-11** — Knowledge system now has aligned visual, text, and machine-readable reading surfaces

- **Visual front door added**: `docs/index.html` is now generated from the latest daily, weekly, and trending notes, giving the repo a human-facing briefing surface with archive browsing, layer navigation, and a visible operating model instead of relying on transient Codex automation cards
- **Text and machine outputs now come from the same generator**: `docs/knowledge/briefing.md` and `docs/knowledge/index.json` are now rebuilt alongside `docs/index.html`, so humans, scripts, and future agents all read from the same surface data rather than inventing separate ad hoc summaries
- **Automation prompts now know to refresh the surface**: the daily radar, weekly synthesis, and GitHub Trending automations now tell each run to rebuild the knowledge surface after writing its radar files
- **Canonical layered architecture restored on the main line**: `docs/feedback/`, `docs/knowledge/`, `docs/proposals/`, and `docs/research/README.md` now define the active boundary `radar -> feedback -> knowledge -> proposals -> implementation`, reusing the good ideas from the older knowledge-system thread without continuing that detached execution lane
- **Radar was explicitly demoted to intake, not the whole knowledge system**: `README.md`, `docs/README.md`, `docs/radar/README.md`, `docs/knowledge/README.md`, and `skills/knowledge/SKILL.md` now align on the same model, and the `knowledge` skill gives future sessions a durable Codex-native entrypoint for synthesis work instead of letting new signal disappear into markdown exhaust or archived threads

**Session 2026-04-08** — Daily radar captured Codex seat governance, Google's terminal-plus-edge agent push, and Anthropic surface silence

- **Daily radar note added for 2026-04-08**: `docs/radar/daily/2026-04-08.md` now records the 2026-04-01 to 2026-04-08 window with first-party signals from OpenAI, Google, and DeepSeek plus a separate GitHub Trending snapshot.
- **OpenAI's strongest in-window move stayed operational rather than model-centric**: the radar note records OpenAI's April 2 `Codex` seat and rate-card update as further evidence that coding-agent execution is being sold and governed as its own runtime surface.
- **Google drove the strongest product motion this week**: April 2 and April 7 posts were consolidated into a clearer read that Google is simultaneously lowering the cost of terminal agentic coding and making on-device agent workflows concrete through Gemma 4, AI Edge Gallery, and LiteRT-LM.
- **Anthropic coverage remained deliberate even without an in-window launch**: the run explicitly checked Anthropic's public technical surfaces and recorded that the newest visible Engineering posts remain March 24 and March 25, keeping official-source silence visible instead of silently dropping Anthropic from the radar.

**Session 2026-04-08-B** — Aligned fresh-session and skills entrypoints with the new Harness explainer

- **Fresh-session start order now leads with the explainer**: `docs/harness-session-handoff.md` now treats `docs/harness-explainer.md` as the fastest human-readable orientation layer before the deeper architecture and package docs, and it no longer points readers at an obsolete release branch when asking them to inspect recent commits
- **Top-level skills index now matches current Harness truth**: `skills/README.md` now describes `Harness` as a stable `1.0.0` leader-led orchestration surface instead of an older prototype, and it points readers to the explainer directly from the skills entrypoint
- **Explainer discoverability is now cleaner across human entrypoints**: the repo now has a more consistent path from `README.md` -> `docs/harness-explainer.md` -> `docs/harness-session-handoff.md` / `skills/README.md`, so readers no longer need to infer whether they should start from a runtime-heavy handoff or an outdated skill index

**Session 2026-04-08-A** — Added a durable human-readable Harness explainer and linked it from the main entrypoints

- **Readable overview is now part of the canonical doc system**: `docs/harness-explainer.md` now explains what `Harness` is, what problem it solves, how it is designed, how the control plane works, and how to use it, instead of forcing readers to infer those answers from architecture and runtime docs alone
- **The explainer is genuinely visual and durable**: the new doc includes Mermaid diagrams plus simplified ASCII fallbacks for the identity/layer model, operating model, and usage flow, so it remains readable both in rich Markdown renderers and plain text contexts
- **Harness entrypoints now route people to the right place**: `README.md`, `docs/README.md`, `docs/harness-documentation-map.md`, `harness/README.md`, and `skills/harness/README.md` now all point to the explainer as the best first stop for humans who want to understand Harness before diving into the deeper canonical docs

**Session 2026-04-03** — Daily radar captured Codex seat economics, Google agent surface convergence, and provider reliability pressure

- **Daily radar note added for 2026-04-03**: `docs/radar/daily/2026-04-03.md` now records the 2026-03-27 to 2026-04-03 signal window with first-party sources from OpenAI, Google, Qwen, and DeepSeek plus a separate GitHub Trending snapshot
- **Strongest market signal shifted toward governed agent runtime**: OpenAI's new Codex seat type and token-priced enterprise access were recorded as the clearest evidence this week that agent execution is being sold and controlled as a separate operational surface
- **Realtime + project-aware coding convergence stayed explicit**: Google's April 1 roundup was captured as a single signal because it bundles Gemini 3.1 Flash Live and the Antigravity coding agent into one build story instead of separate model and tooling announcements
- **Provider-health pressure stayed in frame**: DeepSeek's March 30 and March 31 Web/API incidents were logged as a durability reminder for long-running agent systems, keeping reliability pressure visible without turning it into a direct Harness change request

**Session 2026-04-03** — Daily GitHub Trending note captured a small mixed board with two Codex/AI-adjacent entries

- **A durable daily radar note now exists for the 2026-04-03 Trending run**: `docs/radar/trending/2026-04-03.md` records the default unfiltered global board exactly as captured in this run, including direct project links and a concise board-level read
- **The board was mixed rather than AI-only**: the rendered global board contained just four repositories in this run, spanning screen-recording tooling, Codex customization, system-prompt aggregation, and OSINT username search
- **Codex and model-behavior signal still showed up inside that smaller board**: `oh-my-codex` and `system_prompts_leaks` were explicitly called out as the most relevant AI/agent/developer-tooling entries without forcing the whole board into an AI narrative

**Session 2026-04-02** — GitHub Trending daily note captured a Claude Code and developer-tooling-heavy board

- **A durable daily trending note now exists**: `docs/radar/trending/2026-04-02.md` captures the default global GitHub Trending board for the day instead of forcing an AI-only slice
- **Coding-agent signal was real but not exclusive**: the board leaned hard toward Claude Code, orchestration, rules, and adjacent AI coding workflow repos, while still surfacing broader utility and infrastructure projects such as dashboards, presenters, shell tooling, OSINT, and bundlers
- **Radar now has a reusable daily board checkpoint**: future weekly synthesis work can compare board-level clustering, not just individual publisher posts

**Session 2026-04-02-B** — Daily radar note captured this week's strongest first-party agent/runtime signals

- **A durable first-party radar note now exists**: `docs/radar/daily/2026-04-02.md` captures the strongest official AI and agent developments in the 2026-03-27 to 2026-04-02 window and keeps GitHub Trending explicitly separate as a secondary board-level signal
- **Google and Anthropic drove the strongest technical movement**: Google ADK Go 1.0 and ADK Java 1.0 both foreground service-grade agent runtime concerns such as tracing, HITL, config, and memory/session services, while Anthropic's March 30 release notes widened long-output control and clarified long-context migration pressure onto current Claude lines
- **Quiet publisher surfaces were recorded explicitly**: the daily note makes it clear that OpenAI had no new official technical release inside the window and that Moonshot, MiniMax, Qwen, and DeepSeek did not expose strong tracked first-party technical posts during this run

**Session 2026-04-02** — Radar now treats papers as a first-class knowledge-intake gap

- **Self-iteration research was grounded in first-party sources**: a new research note now records how OpenAI frames Codex self-iteration through explicit self-improvement evals, while Anthropic frames Claude Code self-iteration more through internal dogfooding, autonomous loops, and workflow automation
- **Paper intake is now named as a real radar blind spot**: `docs/radar/sources.md` now adds primary research surfaces such as arXiv, OpenReview, PMLR, and ACL Anthology instead of treating first-party blogs/docs as sufficient for agent-method and benchmark tracking
- **Proposal-first learning loop stayed intact**: the proposal ledger now carries explicit entries for paper-aware radar intake and coding-agent self-iteration watch, so automation changes can be reviewed through the existing ledger instead of bypassing it

**Session 2026-04-01** — Harness now prioritizes cancelled follow-up realignment in mixed decision portfolios

- **Cancelled follow-up correction now outranks ordinary action-adoption work**: when cancelled adopted follow-up realignment coexists with lower-urgency meeting-action adoption residue, scheduler now keeps `realign_cancelled_follow_up` as the primary human-decision packet instead of flattening both decisions into the same generic class
- **Decision-portfolio ordering is now explicit**: required human-decision packets now carry a bounded `decision_priority_rank`, giving scheduler a real kind-aware tiebreaker beyond surface, scope, freshness, priority, and risk alone
- **Core smoke widened again**: core smoke now proves cancelled-follow-up realignment stays primary even when a separate `adopt_meeting_actions` packet is open in the same portfolio, and that `home` preserves the same top-level guidance

**Session 2026-04-01** — Harness now schedules cancelled adopted follow-up as explicit decision work

- **Cancelled follow-up is now a real scheduler decision lane**: if a mission carries adopted follow-up packets but none remain open and at least one has closed as `cancelled`, scheduler now classifies that mission as `needs-human-decision` instead of quietly treating it as runnable residue
- **Leader-facing control plane now says to realign, not just to notice**: scheduler and `home` now let cancelled-follow-up realignment own the top-level decision guidance when no higher-priority control-plane work outranks it
- **Core smoke widened again**: core smoke now proves a cancelled adopted follow-up becomes `decision_kind: realign_cancelled_follow_up`, `classification: needs-human-decision`, and top-level `next_action` guidance in both scheduler and `home`

**Session 2026-04-01** — Harness now refuses handoff of closed adopted follow-up packets

- **Closed follow-up packets can no longer be re-delegated**: `teammates assign --adopted-follow-up=<packet>` now rejects explicitly named adopted follow-up packets unless they are still `open`, so `done` or `cancelled` history cannot be silently rebound as fresh bounded execution
- **Open-only delegation semantics are now explicit**: the adopted follow-up handoff surface now matches the control-plane model that only current open follow-up obligations should produce new bounded teammate work
- **Teammate smoke widened again**: teammate smoke now proves an explicitly named adopted follow-up packet is rejected for handoff once the packet has already been closed

**Session 2026-04-01** — Harness now aligns completed adopted follow-up guidance with audit and doctor

- **Meeting-adoption guidance now understands fully completed follow-up too**: once adopted follow-up packets exist but none remain open, closed-meeting adopted-next-action now expects the same generic post-follow-up guidance that governing run state already uses
- **Audit and doctor stay quiet after real follow-up completion**: `runs authority-audit` and `doctor` no longer have any reason to flag `next_expected_action` drift once the adopted follow-up has actually been completed
- **Teammate smoke widened again**: teammate smoke now proves the fully resolved delegated-follow-up path stays free of `run-missing-meeting-next-action` noise in both `doctor` and run-level authority audit

**Session 2026-04-01** — Harness now aligns cancelled adopted follow-up guidance with audit and doctor

- **Meeting-adoption guidance now understands cancelled follow-up too**: closed meetings whose adopted follow-up packets are all closed as `cancelled` now expect the same explicit realignment guidance that governing run state already uses
- **Audit and doctor no longer misread cancelled follow-up realignment**: `runs authority-audit` and `doctor` now stay quiet when the run is correctly holding cancelled-follow-up realignment guidance instead of flagging false `next_expected_action` drift
- **Teammate smoke widened again**: teammate smoke now proves the cancelled delegated follow-up path stays free of `run-missing-meeting-next-action` noise in both `doctor` and run-level authority audit

**Session 2026-04-01** — Harness now keeps cancelled adopted follow-up on realignment guidance

- **Cancelled follow-up no longer masquerades as completed**: when the last open adopted follow-up packet closes as `cancelled`, the governing run now moves to explicit realignment guidance instead of falling back to generic “follow-up is complete” wording
- **Delegated cancel path now composes cleanly with reopen**: a delegated teammate can cancel a reopened adopted follow-up, the governing run will hold realignment guidance, and the leader can still reopen that packet into bounded execution before final resolution
- **Teammate smoke widened again**: teammate smoke now proves the delegated `cancelled -> reopen -> done` lifecycle and the governing run's cancelled-follow-up realignment guidance

**Session 2026-04-01** — Harness now keeps closed delegated follow-up out of `home.next_action`

- **Home guidance now follows only open adopted follow-up packets**: once a delegated adopted follow-up has been resolved and its governing packet is no longer open, `home` stops letting the old delegated follow-up wording own top-level `next_action`
- **Closed follow-up remains visible without hijacking control**: delegated adopted follow-up metadata can still remain inspectable on leader-facing surfaces after resolution, but it no longer outranks generic post-follow-up execution guidance just because the historical packet is still present
- **Teammate smoke widened again**: teammate smoke now proves `home.next_action` returns to generic active execution guidance as soon as the delegated adopted follow-up closes

**Session 2026-04-01** — Harness now names stale delegated follow-up in authority-audit too

- **Authority-audit now matches scheduler, home, and doctor**: when a teammate run still owns an open adopted follow-up packet but has gone stale, `runs authority-audit` now emits an explicit `delegated_follow_up_stale` adoption gap instead of staying silent about that bounded drift
- **Governing audit now keeps the delegated packet context attached**: the new adoption gap carries the delegated run id, owner, and governing follow-up packet id so the leader can see exactly which bounded adopted follow-up is blocking trustworthy execution
- **Teammate smoke widened again**: teammate smoke now proves run-level authority audit names the stale delegated adopted follow-up gap in the same staged scenario that already drives scheduler/home maintenance escalation and `doctor`

**Session 2026-04-01** — Harness now names stale delegated follow-up in doctor

- **Doctor now matches scheduler maintenance semantics**: a stale teammate run that still owns an open adopted follow-up packet now raises `stale-delegated-follow-up` instead of a generic `run-stale` warning
- **Diagnosis is now aligned across control-plane layers**: scheduler, `home`, and `doctor` now all agree that this is bounded delegated-follow-up drift, not just generic teammate staleness
- **Teammate smoke widened again**: teammate smoke now proves `doctor` names stale delegated adopted follow-up explicitly in the same staged scenario that drives scheduler/home maintenance escalation

**Session 2026-03-31** — Harness now promotes stale delegated follow-up into maintenance

- **Stale delegated follow-up is now a real maintenance object**: when a teammate run owning an open adopted follow-up packet goes stale, scheduler now emits `refresh_stale_delegated_follow_up` instead of folding it into generic `refresh_stale_run`
- **Top-level control-plane guidance now follows that specialized maintenance**: scheduler and `home` now let the stale delegated follow-up maintenance action drive `next_action`, so the leader is told to refresh or close that bounded delegated follow-up before widening execution
- **Teammate smoke widened again**: teammate smoke now proves stale delegated adopted follow-up becomes the primary maintenance action, that both scheduler and `home` follow it in `next_action`, and that a fresh teammate report clears the stale-monitor path

**Session 2026-03-31** — Harness now escalates stale delegated follow-up monitoring

- **Delegated follow-up monitoring now carries freshness too**: scheduler mission packets and `home` now project delegated follow-up freshness fields such as delegated `last_reported_at`, age, and stale status instead of only owner/state
- **Leader guidance now hardens when a delegated follow-up goes stale**: once the delegated teammate stops reporting past the stale-run threshold, scheduler and `home` stop saying to merely monitor it and instead emit explicit stale follow-up recovery guidance
- **Teammate smoke widened again**: teammate smoke now proves a delegated adopted follow-up can be made stale in runtime state, that scheduler and `home` both escalate the monitor guidance, and that a teammate `drive report` refresh returns the lifecycle to normal

**Session 2026-03-31** — Harness now projects delegated follow-up state to leader surfaces

- **Scheduler and home now expose delegated follow-up state directly**: leader-facing control-plane packets now project whether the delegated adopted follow-up is currently `open` instead of only exposing who owns it
- **Delegated ownership is now more monitorable without opening teammate state**: once a bounded follow-up is handed off, the leader can see run id, owner, state, and monitor guidance from scheduler/home alone
- **Teammate smoke widened again**: teammate smoke now proves both scheduler and `home` surface the delegated adopted follow-up state during the handoff path

**Session 2026-03-31** — Harness now reopens delegated follow-up from governing truth

- **Delegated follow-up can now reopen cleanly from the parent side too**: if a delegated teammate already resolved an adopted follow-up packet and the governing run later reopens that packet, runtime now reopens the delegated teammate's bounded follow-up state instead of leaving it stuck on `done`
- **Reopen restores bounded execution guidance**: when the governing packet returns to `open`, the delegated teammate run also regains the follow-up action text in `next_expected_action`, so post-council execution can resume without hand-editing teammate state
- **Teammate smoke widened again**: teammate smoke now proves the delegated run resolves the packet, the governing run reopens it, the delegated run reopens too, and the full lifecycle can still resolve and complete cleanly afterward

**Session 2026-03-31** — Harness now syncs delegated follow-up state from the governing run

- **Delegated teammate state now tracks parent-side follow-up refresh too**: when the governing run updates an adopted follow-up packet directly through `drive resolve-follow-up`, runtime now also refreshes any active delegated teammate run bound to that packet instead of leaving teammate-local delegated state stale
- **Governing truth stays central without hiding delegated lifecycle**: delegated teammate runs still derive their durable packet truth from the governing run, but now that parent-owned truth also reprojects `open` state back into the delegated execution unit before final closure
- **Teammate smoke widened again**: teammate smoke now proves the leader can still address a delegated adopted follow-up packet directly and that the delegated teammate run immediately reflects the refreshed open packet state

**Session 2026-03-31** — Harness now blocks delegated teammate completion until follow-up closes

- **Delegated follow-up ownership now gates teammate completion too**: a teammate run that was bound through `teammates assign --adopted-follow-up` can no longer transition to `completed` while its delegated adopted follow-up packet is still open on the governing run
- **Completion semantics are now symmetric across leader and delegated execution**: once the delegated follow-up is resolved, both the teammate run and the governing leader run can complete normally, so bounded execution closure no longer has a loophole on the delegated side
- **Teammate smoke widened again**: teammate smoke now proves the delegated teammate run is blocked from early completion, then succeeds after `drive resolve-follow-up` closes the governed packet

**Session 2026-03-31** — Harness now projects delegated adopted follow-up into leader-facing control plane

- **Scheduler now knows when adopted follow-up has already been handed off**: mission packets now surface `adopted_primary_follow_up_delegated_run_id`, delegated owner, and delegated monitor guidance instead of only carrying the raw adopted follow-up packet
- **Leader-facing follow-up guidance is now more truthful**: when a selected mission's adopted follow-up is already bound to a teammate run, scheduler next-action guidance now tells the leader to monitor that delegated execution unit instead of repeating the original task wording as if the leader still owned it directly
- **Home now preserves the same delegated execution projection**: `home` now surfaces the delegated teammate run/owner plus delegated monitor guidance for the primary adopted follow-up object, while still letting higher-priority councils or maintenance own the top-level `next_action` when appropriate

**Session 2026-03-31** — Harness now lets delegated teammates resolve adopted follow-up directly

- **Delegated teammate runs can now close the follow-up they own**: `drive resolve-follow-up` now works on a teammate run that was bound through `teammates assign --adopted-follow-up`, routing the resolution back into the governing run while preserving the teammate's bounded ownership
- **Parent and delegated guidance now both advance cleanly**: when that delegated resolution lands, the teammate run moves into post-follow-up reporting guidance while the governing leader run simultaneously advances out of adopted-follow-up handling into generic post-follow-up execution guidance
- **Teammate smoke widened again**: teammate smoke now proves the delegated teammate path can resolve the adopted follow-up directly, clear scheduler/home open-follow-up counts, and still leave the parent run eligible for normal completion

**Session 2026-03-31** — Harness now hands adopted follow-up directly to teammate runtime

- **Adopted follow-up can now become a bounded teammate assignment directly**: `teammates assign --adopted-follow-up` can bind the primary or named adopted follow-up packet from the governing run, derive the delegated objective from that packet, and carry governing meeting / execution-kind provenance into the teammate run and contract
- **Duplicate handoff is now guarded**: the runtime now refuses a second active teammate assignment for the same adopted follow-up packet, so one bounded post-council work object cannot silently fan out into two concurrent owners
- **Delegated follow-up still composes with run lifecycle closure**: if the parent run temporarily shifts into teammate-monitoring guidance while the delegated adopted follow-up is open, resolving that follow-up still advances the parent run out of stale adopted-follow-up handling and back into generic post-follow-up guidance

**Session 2026-03-31** — Harness now blocks completion while adopted follow-up remains open

- **Completion now respects adopted follow-up obligations**: a run can no longer transition to `completed` while adopted follow-up packets remain open, so post-council execution work cannot be silently bypassed by an early close
- **Run-native follow-up completion and run completion now compose cleanly**: after `drive resolve-follow-up` clears the last open adopted packet, the same run can complete normally, which keeps `completed` aligned with real bounded execution closure
- **Teammate smoke widened again**: teammate smoke now proves `drive complete` is rejected while adopted follow-up remains open, then succeeds immediately after the adopted follow-up is resolved

**Session 2026-03-31** — Harness now resolves adopted follow-up from drive

- **Adopted follow-up now has a run-native completion surface**: `drive resolve-follow-up` can now complete the primary or named adopted follow-up packet directly from the governing run instead of forcing the leader back through a closed meeting surface
- **Run guidance now advances when the adopted follow-up closes**: once the last open adopted follow-up packet is completed, the governing run clears the primary adopted packet and advances `next_expected_action` out of the stale follow-up text into generic post-follow-up execution guidance
- **Teammate smoke widened again**: teammate smoke now proves `drive resolve-follow-up` clears adopted open-follow-up counts in scheduler and `home`, while also advancing run guidance after the adopted follow-up completes

**Session 2026-03-31** — Harness now syncs adopted follow-up lifecycle after action resolution

- **Adopted follow-up packets now move with resolved action state**: when `meetings resolve-action` updates a closed meeting that has already been adopted into run state, the governing run now refreshes its adopted follow-up packets, primary adopted packet, and packet counts instead of keeping the `meeting-close` snapshot forever
- **Leader-facing execution surfaces now see follow-up completion**: scheduler and `home` now reduce adopted open-follow-up counts after the underlying meeting action is resolved, so the adopted execution object advances with real completion rather than staying artificially open
- **Smoke widened again**: teammate smoke now proves a decision-record council's adopted open follow-up packet count drops to zero in both scheduler and `home` after `resolve-action`, and core smoke now guards the same lifecycle reduction on the simpler close/resolve path

**Session 2026-03-31** — Harness now anchors adoption audit to the governing closed meeting

- **Adoption checks now follow the adopted council anchor**: run-level adoption audit and doctor now prefer `adopted_from_meeting_id` as the governing closed-meeting reference instead of blindly comparing against whichever closed meeting happened to be newest
- **Adopted follow-up guidance no longer triggers doctor false positives**: once a closed council has been adopted back into run state, doctor now validates `latest_summary` and `next_expected_action` against that governing meeting's summary and adopted follow-up guidance rather than against unrelated later closed councils on the same run
- **Teammate smoke widened again**: teammate smoke now proves doctor stays free of `run-missing-meeting-next-action` and `run-missing-meeting-summary` noise when the run is already aligned with its governing adopted council

**Session 2026-03-31** — Harness now excludes adopted council actions from meeting backlog

- **Adopted council actions now leave meeting backlog cleanly**: once a closed council already has structured output and that output has been adopted back into the governing run, the meeting's remaining open action items no longer count as generic open meeting backlog
- **Leader-facing surfaces now respect the lifecycle handoff**: scheduler mission packets and `home` stop counting those adopted closed-meeting action items inside `open_meeting_action_count`, so the control plane can treat the work as adopted follow-up execution instead of stale meeting residue
- **Core smoke widened again**: home smoke now proves a closed adopted council clears `open_meeting_action_count` immediately at `meeting-close` rather than waiting for every adopted follow-up task to be completed

**Session 2026-03-31** — Harness now promotes human-readable adopted follow-up guidance

- **Adopted follow-up guidance is now leader-readable**: when a closed council already produced an open primary follow-up packet, `meeting-close` and `adopt-meeting` now promote that packet's human-readable summary into `next_expected_action` instead of surfacing an internal follow-up code
- **Adoption audit still stays strict**: `runs authority-audit` now expects the promoted adopted follow-up guidance in `next_expected_action`, so readable execution guidance is part of adoption completeness rather than a cosmetic overlay
- **Teammate smoke widened again**: teammate smoke now proves `meeting-close` promotes human-readable adopted follow-up guidance into `next_expected_action`, while scheduler/home continue to keep that readable adopted guidance visible even after earlier fresh signals are cleared

**Session 2026-03-31** — Harness now projects adopted council output follow-up into leader-facing execution surfaces

- **Adopted council output is now leader-visible**: scheduler mission packets and `home` now surface adopted output packet kinds, adopted/open follow-up packet counts, and the primary adopted follow-up packet instead of leaving that structure discoverable only through raw run inspection
- **Selected execution now follows adopted follow-up packets directly**: when no higher-priority decision or maintenance work outranks the foreground mission, scheduler and `home` now let the primary adopted follow-up packet drive the next action instead of collapsing back to generic “selected mission” wording
- **Smoke widened again**: teammate smoke now proves `missions schedule` plus `home` both preserve protocol-aware adopted follow-up packet execution semantics after `meeting-close`

**Session 2026-03-31** — Harness now adopts structured council output into run state

- **Closed councils now leave durable execution packets behind**: `meeting-close` and `adopt-meeting` now write the generated `output_packet`, follow-up packet count, and primary follow-up packet back into run state instead of only copying summary text into `latest_summary` / `next_expected_action`
- **Adoption drift is now auditable**: `runs authority-audit` now treats missing adopted output metadata, missing adopted follow-up packet counts, and missing adopted primary follow-up packets as explicit `adoption_pending` drift when a closed meeting already produced structured output
- **Meeting identity is now safer under bursty control-plane work**: new meetings now include sub-second entropy in their ids, so the same mission can materialize multiple meetings in one second without accidental id collision or archive overwrite

**Session 2026-03-31** — Harness output packets now synthesize structured follow-up packets

- **Council output is now more executable**: provisional `output_packet` artifacts now synthesize structured `follow_up_packets` from open council action items and unresolved/open lane state instead of leaving follow-up only as freeform summary text
- **Leader-facing output projection widened again**: `meetings summary` and scheduler active-council packets now project follow-up packet counts plus packet payloads alongside lane-level synthesis, so unresolved council output keeps a bounded execution/control-plane shape
- **Protocol-specific execution semantics now survive output synthesis**: `incident-command` output follow-up packets surface `incident_checkpoint`, while `decision-record` output follow-up packets surface `decision_record_follow_up`, keeping protocol intent visible after the council reaches `output`

**Session 2026-03-31** — Harness bundle authority now escalates multi-council overlap explicitly

- **Multi-council bundle escalation is now first-class**: when the same overlapping mission bundle accumulates multiple active council anchors without a single durable governor, mission-graph authority now promotes that state into explicit bundle-governor escalation instead of flattening it into generic overlap conflict
- **Leader guidance is sharper**: multi-council overlap bundles now emit the distinct `multi_council_bundle_governor_required` resolution kind plus a leader-facing action that says to consolidate the competing council anchors or establish one governing mission before widening execution
- **Core smoke widened again**: `scripts/harness-smoke-test.sh` now proves a once-council-governed overlap bundle becomes escalation-required as soon as a second active council anchor appears on the same shared boundary

**Session 2026-03-31** — Harness council output packets now carry lane-level synthesis

- **Council output is now lane-aware**: provisional `output_packet` artifacts now carry synthesized lane packets instead of only one top-level summary, so debated and directly synthesized claim lanes preserve their per-lane disposition, follow-up, and bounded status inside the durable council artifact
- **Leader-facing projection widened again**: `meetings summary` and scheduler active-council packets now project those synthesized lane counts and lane packets alongside the existing `output_packet`, making post-synthesis multi-lane council state more legible without reopening raw findings
- **Smoke widened again**: teammate smoke now proves a multi-lane evidence-review council surfaces accepted and unresolved output lanes explicitly, and core smoke now proves scheduler projects lane-level output state for an active bundle council once it reaches output

**Session 2026-03-31** — Harness 1.0.0 stable line was frozen and fully verified

- **Stable line promotion completed**: the adapter line is now `1.0.0`, `release_stage` is now `stable`, and the current package boundary is no longer described as a prototype-only line
- **Release criteria were satisfied end to end**: the committed `1.0.0` release criteria, package freeze anchors, and full verification matrix were all met in the same round before the version flip
- **Remaining gaps are now explicitly post-1.0 expansion**: deeper scheduler policy, broader bundle adjudication, richer multi-lane debate automation, and future subagent mapping remain worthwhile, but no longer block the first stable Harness line
- **Canonical surfaces and skill entrypoint synchronized**: `_project/CURRENT.md`, `CHANGELOG.md`, session log, canonical Harness docs, manifest metadata, and `skills/harness/README.md` now all describe the stable `1.0.0` line consistently

**Session 2026-03-31** — Harness councils now auto-synthesize protocol-aware output packets

- **Council automation now reaches durable output artifacts**: orchestration now auto-synthesizes protocol-aware `output_packet` records instead of stopping at findings and action items, so active councils can surface a structured provisional output artifact before close
- **Protocol-specific output is now leader-visible**: `meetings summary`, scheduler active-council packets, and meeting outcome adoption now all see the generated output packet, letting `incident-command` expose chosen action/owner/stop-go shape while `decision-record` exposes proposed decision/alternatives/consequences/follow-up shape
- **Council smoke widened again**: teammate smoke now proves both `incident-command` and `decision-record` output packets, and core smoke now proves scheduler projects the active bundle-council output packet once the meeting reaches output
- **Adapter line advanced again**: `harness/adapter-manifest.yml`, `_project/CURRENT.md`, `CHANGELOG.md`, session log, and canonical runtime docs now reflect `0.95.0-prototype`

**Session 2026-03-31** — Harness now carries explicit 1.0.0 release criteria and package freeze anchors

- **1.0.0 release bar is now canonical**: `docs/harness-release-criteria.md` now defines the stable first-line promotion bar, including non-negotiable freeze constraints, required runtime maturity, required packaging state, and the minimum verification matrix for `1.0.0`
- **Package freeze is now machine-addressable**: `harness/adapter-manifest.yml` now points at both `public_contract_path` and `release_criteria_path`, and `adapter validate` now checks those anchors instead of treating release/freeze expectations as chat-only context
- **Runtime bundle contract widened**: packaged runtime archives now include `docs/harness-release-criteria.md`, and package smoke now proves the release criteria ship with the binary, manifest, public contract, and README
- **Adapter line advanced again**: `harness/adapter-manifest.yml`, `_project/CURRENT.md`, `CHANGELOG.md`, session log, documentation map, distribution docs, and handoff/contract docs now reflect `0.94.0-prototype`

**Session 2026-03-31** — Harness active councils now project live action packets and blocked batches into leader-facing surfaces

- **Live council packets are now directly actionable**: scheduler active-council packets and `meetings summary` now project the concrete current-wave action packets, so leader-facing surfaces can see not just the current batch id/phase/count but also the exact open work items inside the live council packet
- **Blocked batch context is now leader-visible**: debate-phase summaries now also project blocked batch summaries alongside the current live packet, making it explicit which bounded participant batches are queued behind the current one instead of leaving that sequencing hidden in raw meeting state
- **Council execution smoke widened again**: core smoke now proves scheduler projects active-council action packets and blocked-batch summaries, while teammate smoke now proves `meetings summary` surfaces current batch action packets plus blocked batch ownership in stateful batch execution
- **Adapter line advanced again**: `harness/adapter-manifest.yml`, `_project/CURRENT.md`, `CHANGELOG.md`, session log, and canonical runtime docs now reflect `0.93.0-prototype`

**Session 2026-03-31** — Harness active councils now project their current live work packet into leader-facing surfaces

- **Live council execution is now more explicit**: scheduler active-council packets and `meetings summary` now project the current live work packet instead of only exposing a generic `next_requirement`, so leader-facing surfaces can see the current packet id, phase, action count, and batch context while a council is in progress
- **Current work packet now spans both review and debate phases**: runtime now derives the current live work packet from either pending independent-review wave obligations or the current debate batch, so the control plane stays informative before and after the council enters cross-examination
- **Council execution smoke widened**: core smoke now proves scheduler projects the current active bundle-council work packet during independent review, and teammate smoke now proves `meetings summary` surfaces current batch summary, assignee, and action count for a multi-lane debate batch
- **Adapter line advanced again**: `harness/adapter-manifest.yml`, `_project/CURRENT.md`, `CHANGELOG.md`, session log, and canonical handoff/contract docs now reflect `0.92.0-prototype`

**Session 2026-03-31** — Harness learning loop was formalized as a proposal-first evolution surface

- **Learning-loop canon added**: `docs/harness-learning-loop.md` now defines how field signals, real usage cases, and internal harness gaps become reviewed proposals instead of direct runtime churn
- **Proposal ledger added**: `docs/harness-proposal-ledger.md` now serves as the durable control-plane surface for reviewed evolution proposals, with status and target-layer tracking
- **Automation intent was tightened**: the existing daily radar and weekly synthesis automations are now explicitly treated as intake/synthesis surfaces for proposal-first Harness evolution rather than direct implementation drivers
- **Docs entrypoints now expose the loop**: `docs/README.md` and `docs/harness-session-handoff.md` now point fresh sessions to the learning loop and proposal ledger instead of leaving that evolution path implicit
- **Current work was reclassified correctly**: the repository now treats this as `designing the learning loop`, not as a runtime/policy mutation path for Harness itself

**Session 2026-03-31** — Harness scheduler now lets the strongest fresh follow-up signal govern mixed-signal missions

- **Mixed-signal urgency is now correct by construction**: `follow_up_signal_urgency_rank` now takes the strongest pending fresh follow-up signal on a mission instead of returning the first matching kind, so a later `risk_raised` signal cannot be accidentally masked by an earlier lower-urgency `priority_change`
- **Leader guidance now follows the strongest pending signal**: when multiple fresh follow-up signals coexist on one mission, scheduler `recommended_action` now follows the highest-urgency pending signal rather than whichever lower-urgency signal happened to be encountered first
- **Mixed-signal smoke is now explicit**: core smoke now proves one mission can hold both fresh `priority_change` and fresh `risk_raised`, while the scheduler still exposes both signals, assigns urgency rank `4`, and emits risk-raised guidance
- **Adapter line advanced again**: `harness/adapter-manifest.yml`, `_project/CURRENT.md`, `CHANGELOG.md`, and the handoff/contract docs now reflect `0.91.0-prototype`

**Session 2026-03-31** — Harness scheduler now prioritizes urgent fresh follow-up decisions explicitly

- **Follow-up urgency is now first-class scheduler state**: mission classifications and required human-decision packets now project `fresh_follow_up_signal_kinds` plus `follow_up_signal_urgency_rank`, so fresh follow-up work no longer depends only on inherited mission priority once it enters the portfolio
- **Fresh risk signals now raise effective scheduling risk before adoption**: when a mission has a fresh `risk_raised` follow-up signal, scheduler now lifts `effective_risk_level` to `high` before that signal is formally adopted into mission state, preventing elevated-risk follow-up work from hiding behind stale lower-risk metadata
- **Human-decision ordering now respects fresh follow-up urgency**: the scheduler queue now lets urgent follow-up control-plane work, especially fresh `risk_raised`, outrank lower-urgency peers such as fresh `priority_change` when they coexist in the same portfolio class
- **Portfolio smoke widened for follow-up urgency**: core smoke now proves `risk_raised` packets expose a higher urgency rank than `priority_change`, become the primary required human decision, and push lower-urgency follow-up work back to queue position 2
- **Adapter line advanced again**: `harness/adapter-manifest.yml`, `_project/CURRENT.md`, `CHANGELOG.md`, and the canonical handoff/runtime docs now reflect `0.90.0-prototype` as the preceding prototype line in the current release climb

**Session 2026-03-31** — Harness scheduler widened fresh-signal rescheduling into more realistic control-plane behavior

- **Signal-kind normalization added**: runtime now canonicalizes signal aliases such as `meeting_required -> requires_meeting`, so scheduler, `home`, authority audit, and maintenance normalize all reason over one stable signal taxonomy instead of depending on spelling drift
- **Fresh review feedback now becomes governed scheduler work**: fresh `review_feedback` signals now reclassify a mission into explicit human-decision follow-up, pause foreground execution when needed, and surface review-specific leader guidance instead of remaining passive telemetry
- **Priority/risk follow-up coverage is now explicit**: core smoke now also proves fresh `priority_change` and `risk_raised` signals go through the governed follow-up decision path with signal-specific guidance, instead of relying on shared-path implicit coverage
- **Fresh rebase requirements now become blocking scheduler work**: fresh `requires_rebase` signals now classify missions as `blocked-by-policy`, emit dedicated `pause_for_fresh_rebase_signal` reschedule actions, and tell the leader to refresh governed run state before continuing
- **Release-trigger council routing widened slightly**: fresh review feedback now participates in release-oriented meeting-trigger/protocol guidance, so release-adjacent follow-up decisions can route toward `incident-command` rather than being flattened into generic coordination
- **Fresh-signal smoke coverage widened**: core smoke now proves alias normalization, review-feedback rescheduling, and requires-rebase blocking end to end before the signal is resolved
- **Canonical docs realigned with runtime truth**: `_project/CURRENT.md`, `CHANGELOG.md`, `docs/harness-architecture.md`, `docs/harness-operating-model.md`, `docs/harness-public-contract.md`, and `docs/harness-session-handoff.md` now all reflect the widened fresh-signal control-plane contract, and the stale `0.82.0-prototype` handoff version note was corrected

**Session 2026-03-30** — erwin was repositioned around Radar + Harness + Skills while preserving historical assets safely

- **Repo positioning was updated**: top-level entrypoints now describe `erwin` as an AI/agent frontier exploration lab instead of a showcase-first skill incubator
- **Radar knowledge base was added**: `docs/radar/` now defines authoritative source policy, active themes, and a kickoff weekly synthesis
- **Archive boundary was created**: `archive/`, `docs/archive/`, and `_project/archive/` now define what is historical versus active before any delete round
- **Session history moved behind an archive boundary**: the real session log directory now lives at `_project/archive/sessions/`, while `_project/sessions` remains as a compatibility symlink so old links and skills continue to work
- **Docs index now reflects the new center of gravity**: active documentation now foregrounds `Harness`, `Radar`, and current repo direction rather than `showcase` and legacy `teamwork`
- **Showcase was physically removed**: after removing the last live harness profile and smoke-test references, `showcase/` was deleted instead of being kept as an archive candidate
- **Low-value legacy material was deleted**: `teamwork-demo/` is gone, low-value session history was pruned aggressively, and five zero-value historical docs were deleted instead of preserved indefinitely

**Session 2026-03-30** — Harness prototype line stayed complete while scheduler and bundle authority deepened into a more mature portfolio control plane

- **Progress-aware active-council prioritization added**: active-council decision packets now expose `active_decision_progress_rank`, and scheduler uses that rank to prefer the more advanced live council when multiple same-class active councils coexist, so a council already in synthesis/output can outrank a fresher but less mature setup-phase peer
- **Active-council progress smoke added**: core smoke now proves two live mission-scoped councils can coexist while the scheduler promotes the more advanced council packet to primary human-decision status instead of ordering them by freshness alone
- **Unified control-plane portfolio added**: scheduler no longer leaves leader-facing priority split across separate human-decision and maintenance portfolios; it now emits one `primary_control_plane_kind`, one `primary_control_plane_packet`, one `primary_control_plane_selection_policy`, and a deferred control-plane count so the leader can see which control-plane object actually comes first
- **Critical maintenance now explicitly outranks lower-priority human decisions**: when governance-critical maintenance such as `clear_invalid_bundle_governor` coexists with fresh but lower-priority human decisions, scheduler now keeps the human decisions visible while still making maintenance the primary control-plane packet and the leader-facing `next_action`
- **Home mirrors unified control-plane priority**: `home` now exposes the same primary control-plane kind, packet, and selection-policy fields that scheduler uses, so leader-facing status no longer requires mentally reconciling separate decision and maintenance priority surfaces
- **Maintenance-first decision deferral added**: when critical maintenance wins the unified control-plane portfolio, scheduler no longer only says maintenance is primary; it now marks governed decision-bearing missions with `deferred_by_primary_control_plane_kind: maintenance` and emits explicit `defer_pending_primary_maintenance` actions so runtime state records that those decisions are intentionally waiting behind governance repair
- **Protocol-aware council execution deepened**: `incident-command` and `decision-record` no longer only differ in route planning; orchestration now gives them protocol-specific independent-review requirements and protocol-specific output guidance, so leader-facing council flow reflects the decision methodology instead of collapsing back into generic evidence-review wording
- **Protocol-specific live guidance validated**: teammate smoke now proves `incident-command` councils ask for bounded incident assessments and close with action/owner/stop-go guidance, while `decision-record` councils ask for bounded decision-record assessments and close with explicit decision/alternatives/consequences/follow-up guidance
- **Scheduler/home now follow live council requirements directly**: once a bundle-scoped decision is already materialized and in flight, scheduler `next_action`, `primary_required_human_decision.recommended_action`, `hold_for_active_bundle_council.reason`, and `home.next_action` now prefer the council's live `next_requirement` rather than falling back to generic “review the council” wording
- **Stale active bundle councils now become explicit primary governance work**: active-council packets now project `active_decision_last_updated_at`, `active_decision_age_hours`, and `active_decision_stale`; when an active council has gone stale under cleanup policy, scheduler promotes “refresh or close the stale active council” into the leader-facing next action and into the bundle hold reason instead of silently preserving stale in-flight governance
- **Core smoke harness was hardened**: `assert_json_field` no longer trips over shell broken-pipe behavior, and the stale-council smoke path now mutates the real active meeting packet structurally instead of relying on brittle text replacement
- **Scheduler decision materialization now reuses active councils**: materializing a scheduler `decision_id` no longer spawns duplicate active meetings when the same governed decision is already under live council control; runtime now reuses the existing active meeting, reports `reused_existing_meeting`, and keeps scheduler `active-council` state stable instead of duplicating council anchors for one control-plane object
- **Active scheduler decisions now remain durably re-addressable after council activation**: `meetings materialize-decision DECISION_ID` no longer fails once a bundle-scoped decision has already promoted into `active-council`; runtime now falls back to the active meeting's persisted `source_decision_packet`, refreshes that live meeting from the latest scheduler context, and still reuses the existing meeting instead of losing the decision object after snapshot rollover
- **Active bundle councils now stay bundle-scoped in the scheduler portfolio**: once a bundle-scoped `establish_bundle_governor` decision has materialized into an active council, scheduler no longer falls back to only per-mission `review_active_council` noise; it now recovers the active bundle decision packet itself, keeps the same bundle-scoped `decision_id` visible in `required_human_decisions`, and continues to treat that active council as one durable control-plane object across scheduling passes
- **Scheduler decisions now project live council progress**: active-council decision packets no longer stop at “review the council” static text; scheduler now mirrors the active meeting's phase, protocol, next requirement, and current debate counters back into the portfolio packet, so leader-facing decision surfaces can follow live council progress without opening the meeting object first
- **Materialized bundle-governor councils now start with seeded participant and candidate context**: when a bundle-scoped `establish_bundle_governor` scheduler decision materializes into a meeting, runtime now auto-seeds the host authoritative run plus bounded candidate-governor participant runs, carries sibling mission refs into packet context, and summarizes the ranked `candidate_governor_packets` inside `known_context` so the leader can start the governed council from scheduler evidence instead of rebuilding the obvious bundle roster by hand
- **Bundle-governor councils now auto-lock durable authority**: auto-synthesis now preserves `source_run_id` on generated findings and can derive a structured `governing_mission_id` lock for bundle-scoped `establish_bundle_governor` councils when accepted evidence already favors one candidate, so `drive meeting-close` can adopt durable bundle authority without a separate manual `meetings lock-decision` in the common path
- **Scheduler-decision materialization added**: `meetings materialize-primary-decision` and `meetings materialize-decision DECISION_ID` now turn durable scheduler human-decision packets into active meetings, so the leader can move directly from the decision portfolio into a governed council without reconstructing context by hand
- **Meeting/source-decision linkage added**: materialized meetings now persist `source_decision_id`, scope/kind/surface metadata, snapshot refs, and the source decision packet itself, so later sessions can audit exactly which scheduler packet created a council
- **Scheduler-to-council bridge hardened**: materialized scheduler meetings now carry scheduler snapshot refs plus recommended-action/protocol context into the meeting packet, and the host authoritative run is moved into explicit waiting state while the council governs the decision
- **Active-council promotion added**: once a scheduler decision has been materialized into an active meeting, scheduler now keeps the same stable `decision_id`, upgrades the decision packet to `active-council`, and projects active meeting ids back onto the packet instead of treating the council as unrelated runtime noise
- **Bundle active-council state added**: when a bundle-scoped decision has already materialized into a live council, scheduler now projects that council back into every affected mission through `bundle_active_decision_meeting_*` and `bundle_control_state`, counts `active_council_bundle_count` at the portfolio level, and emits explicit `hold_for_active_bundle_council` reschedule actions so the bundle is treated as live governed decision work instead of generic overlap noise
- **Explicit-governor versus active-council conflicts now escalate**: if a leader-adopted `bundle_governor_established` signal points at a different mission than the active bundle council host, scheduler now treats that as invalid governance and escalates it into `clear_invalid_bundle_governor` maintenance instead of silently trusting both anchors at once
- **Closed bundle councils now adopt durable authority**: bundle-scoped scheduler councils can now record a structured `governing_mission_id` in `locked_decisions`, and `drive meeting-close` will adopt that verdict into a durable `bundle_governor_established` signal instead of leaving the authority conclusion trapped inside archived meeting output
- **Bundle-governor adoption now supersedes prior explicit anchors**: when a closed council adopts a governor for the same bundle, runtime resolves any earlier active `bundle_governor_established` signals for that bundle before writing the new durable anchor, so governance does not split across multiple active explicit signals after council close
- **Council-adopted governors now keep provenance in graph/scheduler state**: mission-graph authority and scheduler mission classifications now preserve `governance_source_kind` and `adopted_from_meeting_id`, so a durable bundle governor adopted from a closed council is distinguishable from a direct leader-established explicit governor all the way up to leader-facing control surfaces
- **Home/scheduler maintenance parity restored**: `home` now surfaces `scheduler_maintenance_action_count` and `scheduler_maintenance_action_tally`, so leader surfaces no longer hide stale-state maintenance work that scheduler snapshots already know about
- **Bundle-scoped scheduler decisions added**: escalation-required mission bundles no longer emit one repeated `establish_bundle_governor` decision per mission; scheduler now collapses them into one bundle-scoped human-decision packet with merged mission ids, trigger refs, and bundle-required action so the leader sees the actual decision surface instead of repeated mission noise
- **Prioritized decision portfolio added**: scheduler now sorts required human decisions into a real portfolio queue, exposing `primary_required_human_decision`, deferred-decision count, scope/surface tallies, and bundle-vs-mission priority ordering so the leader can see which governed decision should be handled first when multiple human decisions coexist
- **Mission-level queue visibility added**: scheduler now projects the human-decision queue back onto mission classifications, exposing `human_decision_queue_position` and `primary_human_decision` so each affected mission shows whether it is the current primary decision or waiting behind another governed decision
- **Decision-queue deferral added**: non-primary human-decision missions now receive explicit `defer_pending_primary_human_decision` reschedule actions, so the runtime not only shows queue order but also treats lower-priority governed decisions as waiting behind the current primary decision
- **Sticky foreground selection added**: scheduler snapshots now remember the previous foreground mission and retain it when priority, effective risk, and release posture are otherwise tied, so equal-class missions do not thrash the foreground just because a newer mission has a fresher timestamp
- **Sticky human-decision primary added**: scheduler now remembers the previous `primary_required_human_decision` and retains it when a new human-decision packet arrives with the same scope/surface/priority/risk/release class, exposing `decision_selection_policy` and `previous_primary_human_decision_key` so the leader can see that decision stability is intentional rather than accidental
- **Explicit decision protocols added**: meetings now carry `decision_protocol` and `decision_protocol_reason`, making it explicit whether a governed decision should be treated as `evidence-review`, `incident-command`, or `decision-record` instead of silently inheriting one default council style
- **Protocol-aware meeting defaults added**: materialized meetings now seed protocol-specific roles, role specs, and expected output so urgent release/incident decisions and slower decision-record workflows do not start from the same default packet
- **Scheduler protocol guidance added**: human-decision packets now carry `recommended_decision_protocol`, so the leader can distinguish council-grade evidence review, incident-command decisions, and ordinary decision-record handling before materializing a meeting
- **Candidate-governor bundle packets added**: escalation-required mission bundles now expose `candidate_governor_ids` and ranked `candidate_governor_packets` across graph authority, scheduler human-decision packets, and mission classifications, so the leader can see who should plausibly govern the bundle before opening a council
- **Decision-first admission gate added**: when a primary human decision is still unresolved and no established foreground batch is being continued, scheduler now withholds new low-urgency foreground admission under `hold_for_primary_human_decision` instead of quietly starting unrelated fresh work
- **Protocol-aware council controller added**: `incident-command` and `decision-record` meetings no longer auto-expand every same-authority conflict into live route/round work; they now defer bounded disagreement into synthesis/record handling while `evidence-review` keeps the richest debate controller
- **Protocol-deferred conflict visibility added**: debate plans and meeting summaries now expose `protocol_deferred_route_count` and `protocol_bound_route_count`, so leader-facing meeting surfaces can distinguish evidence-review debate lanes from lanes intentionally handed to synthesis under a different decision protocol
- **Protocol-aware auto-synthesis added**: when a non-evidence protocol defers a conflict or stops after a bounded response, auto-synthesis now turns that lane into unresolved findings plus protocol-specific follow-up action items instead of pretending the lane should still generate more live debate
- **Protocol-aware teammate smoke added**: teammate smoke now proves `incident-command` and `decision-record` can both advance from independent review directly to output for same-authority conflicts, without emitting cross-examination route obligations, while still producing bounded follow-up output
- **Explicit bundle-governor adoption added**: `missions establish-governor` and `missions clear-governor` now let the leader convert an escalation-required mission bundle into governed authority deliberately, using durable `bundle_governor_established` signals instead of relying only on implicit candidate-governor heuristics
- **Graph authority now consumes explicit governors**: mission-graph authority can now upgrade a conflicted bundle into `explicit_bundle_governor`, anchoring bundle governance to the leader-adopted mission until that explicit governor signal is cleared
- **Scheduler now respects explicit governor adoption**: once a bundle governor is explicitly established, scheduler no longer keeps surfacing `establish_bundle_governor` for that bundle and instead treats it as governed authority until the governor is retired
- **Core smoke widened again**: `scripts/harness-smoke-test.sh` now proves a bundle can move from `bundle_governor_required` to `explicit_bundle_governor`, stop emitting redundant establish-governor decisions, and then return cleanly to the unresolved state when the explicit governor is cleared
- **Stale explicit-governor maintenance added**: when a leader-adopted explicit bundle governor no longer points at a live active bundle or the nominated governing mission has fallen out of that bundle, scheduler now emits `clear_invalid_bundle_governor` maintenance work instead of silently trusting stale governance
- **Doctor widened for stale bundle governance**: `doctor` now reports `invalid-explicit-bundle-governor` when an active explicit-governor signal no longer matches the active mission graph, so stale governance cannot hide behind otherwise-valid signals
- **Core smoke widened again**: `scripts/harness-smoke-test.sh` now proves that closing the explicit governor mission turns the stale governor into scheduler maintenance, surfaces doctor warnings, and then clears cleanly back to a fresh escalation-required surviving bundle
- **Sticky bounded parallel batches added**: scheduler now retains an already-admitted parallel batch when a same-class fresher peer appears, instead of churning the admitted portfolio just because a new equal-priority mission arrived later
- **Parallel batch memory surfaced**: scheduler snapshots and `home` now expose `previous_selected_mission_ids`, so the leader can see when the control plane is intentionally preserving the admitted parallel set rather than recomputing it from freshness alone
- **Capacity semantics preserved under sticky batches**: even when the scheduler keeps the existing bounded batch, overflow peers still become explicit capacity deferrals with `defer_due_parallel_capacity`, so stability does not erase admission pressure
- **Primary maintenance portfolio added**: scheduler now sorts maintenance work into a real portfolio, exposing `primary_maintenance_action`, deferred-maintenance count, `maintenance_selection_policy`, and `previous_primary_maintenance_key` instead of treating maintenance as an unordered backlog
- **Governance maintenance now gates foreground work**: critical maintenance such as `clear_invalid_bundle_governor` can now trigger `hold_for_primary_maintenance`, deferring low-urgency foreground execution until the primary maintenance action is repaired or explicitly deferred
- **Home mirrors maintenance-first scheduling**: `home` now surfaces the scheduler's primary maintenance action and maintenance-selection metadata, so leader-facing status reflects when governance maintenance outranks fresh foreground execution
- **Sticky current debate-batch control added**: `meetingDebatePlan` now preserves the previous open `current_wave_batch_id` when a new same-phase lane appears, exposing `batch_selection_policy` and `previous_current_wave_batch_id` so the council controller keeps one deterministic active batch instead of thrashing across newly inserted lanes
- **Stable debate-batch identity added**: debate-batch ids now use stable hashed identities rather than truncated slugs, preventing distinct same-phase batches from colliding and accidentally presenting multiple "current" batches at once
- **Broader candidate-governor adjudication added**: larger mission bundles no longer stop at `priority/risk`-only winner selection; graph authority now consumes the fuller candidate-governor tuple, including release-adjacent and governance-activity signals, so a uniquely stronger bundle candidate can be adjudicated even when raw priority/risk still tie
- **Broader bundle smoke added**: core smoke now proves a three-mission escalation bundle can stay unresolved under equal tuples, then become governed once one candidate becomes uniquely release-adjacent, without polluting later runtime health or home/doctor baselines
- **Stable scheduler decision ids added**: required human-decision packets now carry durable `decision_id` values, and mission classifications plus defer actions now project that id back into the runtime, so the decision portfolio is no longer just anonymous JSON and future meeting materialization can target a stable control-plane object
- **Scheduler maintenance escalation added**: scheduler snapshots now emit `maintenance_actions`, `maintenance_action_count`, and `maintenance_action_tally`, turning stale runs, expired signals, and stale meetings into explicit leader-facing maintenance work instead of leaving them only to `doctor`
- **Maintenance guidance prioritized**: when stale control-plane state coexists with runnable foreground work, scheduler `next_action` now tells the leader to keep the foreground batch bounded and clean stale state before widening execution
- **Core smoke widened again**: `scripts/harness-smoke-test.sh` now proves stale run plus stale signal maintenance escalation and verifies that refreshing the run and retiring the stale signal restores a clean control-plane path
- **Council governance gates enforced**: `meetings orchestrate` and `meetings summary` now compute durable `governance_gaps`, and missing authority hierarchy or rejected findings without reasons now block council progression instead of only surfacing as post-hoc doctor warnings
- **Council repair path proven**: teammate smoke now proves both setup-phase blocking on missing authority hierarchy and synthesis/output blocking on missing rejected-finding reasons, then verifies that restoring the missing governance data allows the council controller to continue
- **Bundle-aware scheduler consumption added**: `missions schedule` now consumes mission-bundle authority packets instead of only the raw conflict list, so governed losing missions are deferred behind the current bundle governor, escalation-required bundles become explicit human decisions, and governed-but-ambiguous bundles no longer masquerade as ordinary runnable work
- **Scheduler output widened**: mission classifications now surface `bundle_id`, `bundle_state`, `bundle_role`, `bundle_resolution_kind`, `bundle_governing_mission_id`, and `bundle_required_action`, making bundle-level control visible directly in durable scheduler snapshots
- **Leader-facing scheduler priority fixed**: when required human bundle decisions coexist with a still-valid foreground batch, scheduler `next_action` now prioritizes the human decisions instead of pretending foreground selection is the only thing that matters; snapshots also emit `bundle_role_tally` so the leader can see the portfolio composition directly
- **Home/scheduler parity restored**: `home` now reuses scheduler decision-first guidance when escalation-required bundle decisions outrank generic cross-mission conflict wording, and it surfaces `scheduler_bundle_role_tally` so the leader sees the same bundle portfolio shape in both places
- **Bounded portfolio admission added**: scheduler now consumes profile-backed `max_parallel_runs`, so parallel execution no longer widens to all runnable missions by default; overflow runnable missions are deferred with explicit `defer_due_parallel_capacity` actions and the snapshot records `batch_admission_policy`, `max_parallel_runs`, and `capacity_deferred_*` metadata
- **Fresh-signal rescheduling added**: scheduler now treats fresh `requires_meeting`, `interface_change`, `dependency_discovered`, and `blocked` signals as real control-plane triggers instead of display-only telemetry, reclassifying missions and emitting pause/reschedule actions when the signal has not yet been adopted into runtime state
- **Profile-backed high-risk inference added**: scheduler now consumes profile `high_risk_scopes`, matching them against mission/run ownership selectors and release-adjacent semantics to derive `profile_high_risk`, `profile_high_risk_scope_refs`, and `effective_risk_level`; missions touching configured high-risk boundaries such as `install.sh` now serialize the foreground batch even when no manual `risk_level` was set
- **Profile-backed council recommendation added**: scheduler now consumes profile `meeting_triggers`, deriving `recommended_decision_surface`, `meeting_trigger_refs`, and `council_decision_count` so escalation-required bundle decisions can explicitly recommend council handling instead of collapsing into generic human-decision output
- **P0 serialization added**: scheduler now enforces the control-plane rule that `P0` missions never widen into parallel foreground batches; once a runnable `P0` mission exists, peers are deferred with explicit `defer_until_p0_clears` actions and the batch admission policy becomes `serialize_p0_priority`
- **Explicit parallel justification added**: scheduler now consumes the execution policy `require_explicit_parallel_justification`, so the foreground batch only widens into true parallel execution when every current batch candidate carries explicit `parallel_justification`; otherwise the control plane falls back to serialized execution with `defer_due_missing_parallel_justification`
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

**The first stable `Harness` line is now delivered: runtime, control plane, public contract, authority model, council model, naming, distribution strategy, release criteria, and a full verified release matrix are all formalized. Remaining work is post-1.0 expansion, not a blocker for the stable line:**

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 1 | ✅ | Research protocol / control plane / authority / cybernetics synthesis |
| Phase 2 | ✅ | Build repository-local harness definition layer (`harness/`) |
| Phase 3 | ✅ | Implement runtime substrate (`scripts/harness_*`) |
| Phase 4 | ✅ | Validate `home` / `doctor` / `claim` / `drive` / `meetings` / `maintenance` |
| Phase 5 | ✅ | Package the harness formally as a reusable Codex-native adapter/system |
| Phase 6 | ✅ | Split reusable adapter/bootstrap flow from repo-local proving ground |
| Phase 7 | 🔄 | Post-1.0 expansion: broader authority/adaptation layers and optional Codex-native multi-agent integration |

**Learning loop status**

- `Harness` learning-loop design is now defined as a proposal-first evolution surface
- daily and weekly automations should now be interpreted as intake/synthesis for reviewed proposals
- the current phase is designing and hardening that learning loop, not automatically mutating runtime behavior from fresh external signals

## Next Session Should Start With

1. Read this file
2. Read `docs/erwin-operating-model.md` before deciding whether the work belongs to Radar, Feedback, Knowledge, Proposals, or Implementation
3. **Radar next work**:
   - A. Continue observing live daily / weekly / trending / paper-watch output against Tier 1 publisher quality
   - B. Tune `docs/radar/sources.md` and theme weighting only when real output shows durable drift
   - C. Keep radar as intake; do not recommend direct runtime changes from fresh signal
4. **Knowledge / Feedback next work**:
   - A. Add normalized feedback cases whenever real Harness, Codex, automation, or repo-process friction repeats
   - B. Keep `docs/knowledge/active-threads.md` current when weekly synthesis, feedback, or proposals change durable understanding
   - C. Use mature active threads before creating or widening proposals
5. **Harness / Implementation next work**:
   - A. Treat the human-readable `--brief` stack as delivered across the current Harness surfaces, not as a still-missing `home` / `doctor` starter layer
   - B. Drive post-1.0 work from real usage feedback, learning-loop proposals, eval candidates, and explicit internal gaps
   - C. Use the proposal ledger before promoting field intelligence or real-usage patterns into runtime/backlog changes
6. If starting fresh, read `AGENTS.md`, `_project/CURRENT.md`, `docs/erwin-operating-model.md`, `docs/README.md`, and then the layer-specific docs for the task

## Implementation Roadmap

### Radar
- [x] Create `docs/radar/` knowledge base
- [x] Define source priority and theme taxonomy
- [x] Seed first weekly synthesis
- [x] Create and activate the daily and weekly app automations
- [ ] Observe one week of live signal and tune publisher/theme balance

### Harness
- [x] Establish the first stable Codex-native Harness line (`1.0.0`)
- [ ] Continue authority, council, and scheduler depth beyond the stable line
- [ ] Decide how future subagent integration maps onto the teammate runtime

### Repo Slimming
- [x] Create archive boundary before deletion
- [x] Move session history behind archive with compatibility preserved
- [x] Reposition repo entrypoints around Radar + Harness + Skills
- [x] Delete `teamwork-demo/` because it no longer provided active reference value
- [x] Delete five historical docs that no longer carried practical reference value
- [x] Remove active harness dependencies on `showcase/` and delete it
- [x] Keep `skills/` surfaces, including `skills/teamwork/`, as intentionally retained repo assets

## Known Issues / Blockers

- Wing scripts reference `${CLAUDE_SKILL_DIR}` which requires Anthropic skill infrastructure
- skill-lint.sh ### header check has false positives (structural headings vs output templates)
- Harness is now installed and validated in Codex via `install-codex.sh`, can bootstrap into another repo, has manifest-driven adapter validation/packaging, first-class mission state, a leader-facing scheduler baseline with durable snapshots, execution-mode selection, direct bundle-aware deferral/escalation, mission-level and graph-level authority adjudication including dependency-governed ownership verdicts, selector-specificity winners, council-governed bundle authority, bundle-aware grouping, explicit bundle governors, selector-aware ownership transfer with shared-copy and exclusion-aware split modes, artifact-aware authority adjudication, teammate runtime, council orchestration with phase history + synthesis audit semantics, a formalized runtime distribution strategy, committed `1.0.0` release criteria, and a broad human-readable `--brief` stack across home, doctor, claim, drive, missions, runs, signals, teammates, ship, rc, meetings, and maintenance; remaining work is post-1.0 expansion driven through the five-layer learning loop

## File Locations

| What | Where |
|------|-------|
| Project goals | `_project/GOALS.md` |
| Repo entrypoint | `README.md` |
| Radar | `docs/radar/` |
| Idea garden | `_project/ideas/` |
| ADRs | `_project/decisions/` |
| Harness stable line | `harness/`, `scripts/harness_*`, `skills/harness/` |
| Wing skills | `skills/wing/` |
| Wing architecture | `_project/ideas/skills/tw2-architecture.md` |
| Teamwork skills (legacy) | `skills/teamwork/` |
| Docs index | `docs/README.md` |
| Historical boundary | `archive/README.md` |
