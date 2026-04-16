# Harness

Use this skill when you are operating the repository-local orchestration harness as the leader.

The harness is file-backed:

- committed definitions live in `harness/`
- runtime state lives in `.harness/`
- command surfaces live in `scripts/`
- the runtime core is a Go CLI resolved by `scripts/harness-runtime.sh`
- packaged binaries in `harness/bin/<os>-<arch>/harness-state` are preferred; local Go build is only the fallback path

## Leader workflow

1. Start with the home view:
   `bash scripts/harness-home.sh`
   Use `bash scripts/harness-home.sh --brief` when you need a human-readable leader summary before drilling into raw JSON.
2. If anything looks off, inspect runtime health:
   `bash scripts/harness-doctor.sh`
   Use `bash scripts/harness-doctor.sh --brief` when you want the top repair priorities in readable prose first.
3. Validate the reusable adapter contract if packaging or sync drift is suspected:
   `bash scripts/harness-adapter.sh validate`
4. Inspect the mission control plane when dependency shape or execution topology is unclear:
   `bash scripts/harness-missions.sh ...`
   Use `bash scripts/harness-missions.sh get MISSION_ID --brief` when you want one mission's state, queue position, next step, and council/teammate context in readable prose first.
5. Create or draft executable issues when work is not yet claimable:
   `bash scripts/harness-issue.sh ...`
6. Claim work into governed execution:
   `bash scripts/harness-claim.sh ISSUE_NUMBER [TITLE_OVERRIDE]`
   Use `bash scripts/harness-claim.sh ISSUE_NUMBER [TITLE_OVERRIDE] --brief` when you want the created mission/run/contract and leader boundary in readable prose first.
7. Drive run state:
   `bash scripts/harness-drive.sh ...`
   Use `--brief` on `report`, `start`, `resume`, `wait`, `block`, `complete`, and `resolve-follow-up` when you want the progress update, state mutation, or adopted follow-up closure explained before drilling into raw JSON.
8. Delegate bounded work when leader execution should split:
   `bash scripts/harness-teammates.sh ...`
9. Prepare release handoff when execution is done:
   `bash scripts/harness-ship.sh ...`
10. Use meetings when ambiguity or blast radius is high:
   `bash scripts/harness-meetings.sh ...`
11. Use RC flow for release candidate governance:
   `bash scripts/harness-rc.sh ...`
12. Repair drift or normalize runtime state when needed:
   `bash scripts/harness-maintenance.sh ...`

## Important commands

- `home`
  `bash scripts/harness-home.sh`
- `doctor`
  `bash scripts/harness-doctor.sh`
- `adapter-validate`
  `bash scripts/harness-adapter.sh validate`
- `missions-list`
  `bash scripts/harness-missions.sh list`
- `missions-schedule-brief`
  `bash scripts/harness-missions.sh schedule --brief`
- `missions-get`
  `bash scripts/harness-missions.sh get issue-1234`
- `missions-get-brief`
  `bash scripts/harness-missions.sh get issue-1234 --brief`
- `missions-authority-audit-brief`
  `bash scripts/harness-missions.sh authority-audit issue-1234 --brief`
- `missions-graph-brief`
  `bash scripts/harness-missions.sh graph --brief`
- `missions-graph-authority-brief`
  `bash scripts/harness-missions.sh authority-audit --all --brief`
- `runs-get-brief`
  `bash scripts/harness-runs.sh get RUN_ID --brief`
- `teammates-get-brief`
  `bash scripts/harness-teammates.sh get TEAMMATE_RUN_ID --brief`
- `meetings-summary-brief`
  `bash scripts/harness-meetings.sh summary MEETING_ID --brief`
- `meetings-orchestrate-brief`
  `bash scripts/harness-meetings.sh orchestrate MEETING_ID --brief`
- `meeting-authority-audit-brief`
  `bash scripts/harness-meetings.sh authority-audit MEETING_ID --brief`
- `missions-graph`
  `bash scripts/harness-missions.sh graph`
- `issue-draft`
  `bash scripts/harness-issue.sh draft "Mission title" --criteria="A|B"`
- `claim`
  `bash scripts/harness-claim.sh 1234`
- `claim-brief`
  `bash scripts/harness-claim.sh 1234 --brief`
- `claim-with-ownership`
  `bash scripts/harness-claim.sh 1234 --depends-on=run-older-claim --owns=surface:release-notes`
- `start`
  `bash scripts/harness-drive.sh start RUN_ID`
- `start-brief`
  `bash scripts/harness-drive.sh start RUN_ID --brief`
- `wait`
  `bash scripts/harness-drive.sh wait RUN_ID "Waiting on sibling verification" --on=RUN_A|RUN_B`
- `wait-brief`
  `bash scripts/harness-drive.sh wait RUN_ID "Waiting on sibling verification" --on=RUN_A|RUN_B --brief`
- `block`
  `bash scripts/harness-drive.sh block RUN_ID "reason"`
- `block-brief`
  `bash scripts/harness-drive.sh block RUN_ID "reason" --brief`
- `resume-brief`
  `bash scripts/harness-drive.sh resume RUN_ID --brief`
- `complete-brief`
  `bash scripts/harness-drive.sh complete RUN_ID "Execution complete" --brief`
- `meeting`
  `bash scripts/harness-drive.sh meeting RUN_ID "question"`
- `teammate-assign`
  `bash scripts/harness-teammates.sh assign RUN_ID worker-auth implementation-verifier "Update auth callers" --scope="auth callers only" --owns=scope:auth-callers --depends-on=RUN_OTHER`
- `ownership-transfer`
  `bash scripts/harness-runs.sh transfer-ownership SOURCE_RUN_ID TARGET_RUN_ID --owns=scope:release-auth-smoke --reason="Shift bounded ownership to the implementation path"`
- `run-authority-audit`
  `bash scripts/harness-runs.sh authority-audit RUN_ID`
- `run-authority-audit-brief`
  `bash scripts/harness-runs.sh authority-audit RUN_ID --brief`
- `signals-active-brief`
  `bash scripts/harness-signals.sh active MISSION_ID --brief`
- `signals-adjudicate-brief`
  `bash scripts/harness-signals.sh adjudicate --all --brief`
- `teammate-report`
  `bash scripts/harness-teammates.sh report TEAMMATE_RUN_ID "Caller patch is in place" --kind=progress --evidence=src/auth.ts|tests/auth.spec.ts`
- `signals-adjudicate`
  `bash scripts/harness-signals.sh adjudicate --all`
- `teammate-message`
  `bash scripts/harness-teammates.sh message TEAMMATE_RUN_A TEAMMATE_RUN_B interface_change "Signature changed to require scope" --evidence=src/auth.ts`
- `meeting-participant`
  `bash scripts/harness-meetings.sh add-participant MEETING_ID TEAMMATE_RUN_ID architecture-challenger`
- `meeting-role`
  `bash scripts/harness-meetings.sh add-role MEETING_ID architecture-challenger --mission="Challenge blast radius assumptions" --allow="L1: live git state|L2: active runtime state" --forbid="L4+: stale notes"`
- `meeting-known-context`
  `bash scripts/harness-meetings.sh add-known-context MEETING_ID "Current release notes mention auth callers but not required scope."`
- `meeting-lock`
  `bash scripts/harness-meetings.sh lock-decision MEETING_ID "Keep this council scoped to auth release notes only." --reason="Avoid widening blast radius during review." --authority=L2 --evidence=harness/profiles/erwin.yml`
  `bash scripts/harness-meetings.sh lock-decision MEETING_ID "Keep this council scoped to auth release notes only." --reason="Avoid widening blast radius during review." --authority=L2 --evidence=harness/profiles/erwin.yml --brief`
- `meeting-position`
  `bash scripts/harness-meetings.sh add-position MEETING_ID TEAMMATE_RUN_ID independent_review challenge "This contract leaks release policy" --claim=release-policy-scope --dimension=factual --evidence=docs/spec.md`
  `bash scripts/harness-meetings.sh add-position MEETING_ID TEAMMATE_RUN_ID independent_review challenge "This contract leaks release policy" --claim=release-policy-scope --dimension=factual --evidence=docs/spec.md --brief`
- `meeting-orchestrate`
  `bash scripts/harness-meetings.sh orchestrate MEETING_ID`
- `meeting-finding`
  `bash scripts/harness-meetings.sh add-finding MEETING_ID rejected "Ship immediately without release-note updates." --reason="Cross-examination found operator-facing drift." --authority=L2 --evidence=CHANGELOG.md --role=ops-simulator`
  `bash scripts/harness-meetings.sh add-finding MEETING_ID rejected "Ship immediately without release-note updates." --reason="Cross-examination found operator-facing drift." --authority=L2 --evidence=CHANGELOG.md --role=ops-simulator --brief`
- `meeting-action`
  `bash scripts/harness-meetings.sh add-action MEETING_ID "Update release notes before merging the auth change."`
  `bash scripts/harness-meetings.sh add-action MEETING_ID "Update release notes before merging the auth change." --brief`
- `meeting-resolve-action`
  `bash scripts/harness-meetings.sh resolve-action MEETING_ID ACTION_ID done`
  `bash scripts/harness-meetings.sh resolve-action MEETING_ID ACTION_ID done --brief`
- `meeting-close`
  `bash scripts/harness-drive.sh meeting-close RUN_ID MEETING_ID "resolution"`
- `adopt-meeting`
  `bash scripts/harness-drive.sh adopt-meeting RUN_ID MEETING_ID`
- `reprovision`
  `bash scripts/harness-drive.sh reprovision RUN_ID --mode=worktree --apply`
- `ship-verify`
  `bash scripts/harness-ship.sh verify RUN_ID "targeted-check:pass"`
- `ship-push`
  `bash scripts/harness-ship.sh push RUN_ID`
- `ship-pr-draft`
  `bash scripts/harness-ship.sh pr-draft RUN_ID`
- `ship-pr-open`
  `bash scripts/harness-ship.sh pr-open RUN_ID --draft`
- `ship-pr`
  `bash scripts/harness-ship.sh pr RUN_ID PR-123 https://example.invalid/pr/123`
- `ship-review`
  `bash scripts/harness-ship.sh review RUN_ID approved "Approved after fixes"`
- `ship-checks`
  `bash scripts/harness-ship.sh checks RUN_ID passed "ci:passed"`
- `ship-merge`
  `bash scripts/harness-ship.sh merge RUN_ID --via=git "Merged cleanly"`
- `ship-cleanup`
  `bash scripts/harness-ship.sh cleanup RUN_ID --delete-branch`
- `ship-merged`
  `bash scripts/harness-ship.sh merged RUN_ID`
- `rc-create`
  `bash scripts/harness-rc.sh create v1.2.3 RUN_ID`
- `rc-fix`
  `bash scripts/harness-rc.sh fix RC_ID COMMITISH --apply --push`
- `rc-cancel`
  `bash scripts/harness-rc.sh cancel RC_ID "staging bug found" --apply`
- `rc-hotfix`
  `bash scripts/harness-rc.sh hotfix v1.2.4 COMMITISH --apply --push`
- `rc-backport`
  `bash scripts/harness-rc.sh backport v1.2.4 COMMITISH --apply --push`
- `rc-promote`
  `bash scripts/harness-rc.sh promote RC_ID "promotion note"`
- `maintenance`
  `bash scripts/harness-maintenance.sh reconcile`
  `bash scripts/harness-maintenance.sh reconcile --brief`
  `bash scripts/harness-maintenance.sh normalize`
  `bash scripts/harness-maintenance.sh normalize --brief`
- `desired-state`
  `bash scripts/harness-runtime.sh runs set-desired RUN_ID --goal=completed --verification=passed --ship=merged`
- `all-smoke`
  `bash scripts/harness-all-smoke-test.sh`
- `runtime-dist-smoke`
  `bash scripts/harness-runtime-dist-smoke-test.sh`

## Operating rules

- Treat `harness-home` as the control plane entry point.
- Use `claim --brief` when a human needs the newly created mission/run/contract, provisioning result, and next leader action without reading raw claim JSON.
- Use `drive report|start|resume|wait|block|complete|resolve-follow-up --brief` when a human needs the run progress update, state transition, adopted follow-up result, emitted blocker signal, and leader/teammate/council boundary without reading raw run JSON.
- Use `missions schedule --brief` when you want the scheduler portfolio itself in human-readable form before drilling into raw scheduler JSON or the durable snapshot.
- Treat `harness-missions.sh` as the mission-level control plane above individual runs.
- Treat `harness-adapter.sh validate` as the package boundary check when bootstrap/sync behavior is in doubt.
- Treat `harness-runtime-dist-smoke-test.sh` as the proof that a repo can run the packaged runtime without a local Go toolchain.
- Use `maintenance reconcile --brief` when you want the human-readable repair result for orphaned signals, superseded signals, and orphaned meetings before drilling into raw reconciliation JSON.
- Use `maintenance normalize --brief` when you want the human-readable migration result for legacy run, signal, and contract records before drilling into raw normalization JSON.
- Use `meetings lock-decision --brief` when you want the human-readable result of a council scope lock before drilling into raw meeting JSON.
- Use `meetings add-position --brief` when you want the human-readable result of a newly recorded council position before drilling into raw meeting JSON.
- Use `meetings add-finding --brief` when you want the human-readable result of a newly recorded council judgment before drilling into raw meeting JSON.
- Use `meetings add-action --brief` when you want the human-readable result of a newly recorded council action item, including whether leader assignment or adoption is still needed, before drilling into raw meeting JSON.
- Use `meetings resolve-action --brief` when you want the human-readable result of closing, cancelling, or reopening a council action, including adopted follow-up sync, before drilling into raw meeting JSON.
- Treat `maintenance normalize` as the migration path for older `.harness/` state when mission objects or ownership claims were introduced after those runs were created.
- Prefer the narrowest safe state change.
- Use meetings for high-ambiguity or high-blast-radius decisions.
- Choose the meeting protocol deliberately: `evidence-review` for contested technical judgment, `incident-command` for urgent release/rollback decisions, and `decision-record` for bounded design choices that need durable recording more than live debate.
- Treat teammate runs as first-class runtime objects, not informal notes.
- Treat missions as first-class durable work objects; runs are execution attempts, not the whole control plane.
- Prefer explicit `--owns=` selectors when splitting work; shared ownership should be intentional, not accidental.
- Use path-shaped selectors such as `file:...`, `dir:...`, `module:...`, and `glob:...` when the boundary is code or docs structure, not just abstract scope labels.
- Use `scope:` and `interface:` deliberately for semantic boundaries; the runtime now treats same-value and parent/child scope-interface selectors as overlapping governed ownership, not unrelated labels.
- Prefer formal ownership transfer over silent selector drift when bounded work needs to move between runs.
- Treat active signals as provisional until they survive authority/freshness comparison with newer run or meeting state.
- Treat `depends_on` plus desired state as part of the control plane, not just metadata.
- Allow teammate-to-teammate communication only through governed links or sibling relationships.
- Use `meetings orchestrate` to advance council phases instead of manually skipping straight to synthesis.
- Use `meetings orchestrate MEETING_ID --brief` when you want the human-readable orchestration result: phase advancement or hold, current wave/debate/output work packet, auto-synthesis effect, and next leader/council action before drilling into raw orchestration JSON.
- Use `drive meeting ... --protocol=...` or `meetings set-protocol` when the default `evidence-review` protocol is not the right fit for the decision surface.
- Treat `meetings orchestrate` as the council chair: once bounded positions exist, it can auto-synthesize provisional findings and unresolved-action prompts without overwriting manual judgment that is already on the record.
- Treat missing independent review as first-class wave debt: `meetings orchestrate` now reports `open_wave_action_count` so the leader can see who still owes a bounded review before cross-examination.
- Treat routed cross-examination as first-class: when multiple roles disagree on the same `--claim`, the council now turns that lane into an explicit route obligation with an assigned bounded challenger before synthesis/output is trusted.
- Treat follow-up rounds as first-class: `meetings orchestrate` can now keep assigning bounded round 2 / round 3 / ... replies until a lane resolves on authority/freshness or the configured debate-round cap is reached.
- When recording cross-examination positions, include explicit `--responds-to-run=` / `--responds-to-role=` linkage so the council trace shows who challenged whom.
- When recording positions that speak about the same question, include `--claim=` and `--dimension=` so auto-synthesis and authority audit can group them into the same council claim.
- When multiple findings speak about the same claim, use `--claim=` and `--dimension=` on `meetings add-finding` so authority conflicts can be audited cleanly.
- Use `meetings summary` before adoption when you need to see synthesis readiness, debate-round depth, blocker counts, and any still-open wave/route obligations at a glance.
- Use `meetings summary MEETING_ID --brief` when you want the human-readable council state, blocking reason, current work packet, and protocol-aware output guidance before drilling into raw meeting JSON.
- Use `meetings authority-audit MEETING_ID --brief` when you want the human-readable same-claim authority verdict, whether a tie still needs judgment, and the next leader action before drilling into raw council-audit JSON.
- Use `meetings authority-audit` when you need a leader-facing explanation of which same-claim finding currently wins on authority.
- Use `runs authority-audit RUN_ID --brief` when you want the human-readable run-drift verdict, repair action, and leader/teammate/council implications before drilling into raw audit JSON.
- Use `runs authority-audit` when a run may be drifting against fresher or higher-authority active signals, or when `doctor` reports `run-authority-conflict`.
- Use `runs get RUN_ID --brief` when you want the human-readable execution summary for one run before drilling into raw run JSON or authority-audit detail.
- Use `signals active [MISSION_ID] --brief` when you want the human-readable active signal priority, freshness, scope, and leader/teammate/council next action before drilling into raw signal JSON.
- Use `signals adjudicate SIGNAL_ID|--all --brief` when you want the human-readable signal-repair result: which signals were resolved as superseded, why fresher run or meeting truth won, and what leader/teammate/council action remains before drilling into raw adjudication JSON.
- Use `teammates assign PARENT_RUN_ID TEAMMATE_ID ROLE [OBJECTIVE] --brief` when you want the human-readable bounded delegation handoff, created teammate run, contract, ownership boundary, assignment signal, and leader/teammate/council next actions before drilling into raw assignment JSON.
- Use `teammates get TEAMMATE_RUN_ID --brief` when you want the human-readable bounded teammate assignment, dependency/report freshness state, governed communication context, and leader/teammate/council next actions before drilling into raw teammate JSON.
- Use `missions graph --brief` when you want the human-readable active root/dependency topology, current graph-level next step, and teammate/council coverage before drilling into raw graph JSON.
- Use `missions authority-audit MISSION_ID --brief` when you want the human-readable mission-control verdict, leader-over-teammate adjudication, and next repair action before drilling into raw mission-audit JSON.
- Use `missions authority-audit --all --brief` when you want the human-readable mission-graph bundle verdict, governing or candidate mission set, and next graph-level leader action before drilling into raw bundle-audit JSON.
- Use `runs authority-audit` when a run claims it is ready/verified/clean but the live worktree is still dirty; the audit now surfaces that intent-vs-live-git disagreement explicitly.
- Use `missions authority-audit --all` when mission bundles overlap; it now explains bundle-level `resolution_kind`, shared selectors, and the leader-facing `required_action` instead of only listing pairwise conflicts.
- After a meeting closes, either let `meeting-close` apply its outcome or explicitly run `adopt-meeting`.
- Keep `doctor` clean. If it shows drift, repair that before piling on new execution.
