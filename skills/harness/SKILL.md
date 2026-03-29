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
2. If anything looks off, inspect runtime health:
   `bash scripts/harness-doctor.sh`
3. Validate the reusable adapter contract if packaging or sync drift is suspected:
   `bash scripts/harness-adapter.sh validate`
4. Inspect the mission control plane when dependency shape or execution topology is unclear:
   `bash scripts/harness-missions.sh ...`
5. Create or draft executable issues when work is not yet claimable:
   `bash scripts/harness-issue.sh ...`
6. Claim work into governed execution:
   `bash scripts/harness-claim.sh ISSUE_NUMBER [TITLE_OVERRIDE]`
7. Drive run state:
   `bash scripts/harness-drive.sh ...`
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
- `missions-get`
  `bash scripts/harness-missions.sh get issue-1234`
- `missions-graph`
  `bash scripts/harness-missions.sh graph`
- `issue-draft`
  `bash scripts/harness-issue.sh draft "Mission title" --criteria="A|B"`
- `claim`
  `bash scripts/harness-claim.sh 1234`
- `claim-with-ownership`
  `bash scripts/harness-claim.sh 1234 --depends-on=run-older-claim --owns=surface:release-notes`
- `start`
  `bash scripts/harness-drive.sh start RUN_ID`
- `wait`
  `bash scripts/harness-drive.sh wait RUN_ID "Waiting on sibling verification" --on=RUN_A|RUN_B`
- `block`
  `bash scripts/harness-drive.sh block RUN_ID "reason"`
- `meeting`
  `bash scripts/harness-drive.sh meeting RUN_ID "question"`
- `teammate-assign`
  `bash scripts/harness-teammates.sh assign RUN_ID worker-auth implementation-verifier "Update auth callers" --scope="auth callers only" --owns=scope:auth-callers --depends-on=RUN_OTHER`
- `ownership-transfer`
  `bash scripts/harness-runs.sh transfer-ownership SOURCE_RUN_ID TARGET_RUN_ID --owns=scope:release-auth-smoke --reason="Shift bounded ownership to the implementation path"`
- `run-authority-audit`
  `bash scripts/harness-runs.sh authority-audit RUN_ID`
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
- `meeting-position`
  `bash scripts/harness-meetings.sh add-position MEETING_ID TEAMMATE_RUN_ID independent_review challenge "This contract leaks release policy" --claim=release-policy-scope --dimension=factual --evidence=docs/spec.md`
- `meeting-orchestrate`
  `bash scripts/harness-meetings.sh orchestrate MEETING_ID`
- `meeting-finding`
  `bash scripts/harness-meetings.sh add-finding MEETING_ID rejected "Ship immediately without release-note updates." --reason="Cross-examination found operator-facing drift." --authority=L2 --evidence=CHANGELOG.md --role=ops-simulator`
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
  `bash scripts/harness-maintenance.sh normalize`
- `desired-state`
  `bash scripts/harness-runtime.sh runs set-desired RUN_ID --goal=completed --verification=passed --ship=merged`
- `all-smoke`
  `bash scripts/harness-all-smoke-test.sh`
- `runtime-dist-smoke`
  `bash scripts/harness-runtime-dist-smoke-test.sh`

## Operating rules

- Treat `harness-home` as the control plane entry point.
- Treat `harness-missions.sh` as the mission-level control plane above individual runs.
- Treat `harness-adapter.sh validate` as the package boundary check when bootstrap/sync behavior is in doubt.
- Treat `harness-runtime-dist-smoke-test.sh` as the proof that a repo can run the packaged runtime without a local Go toolchain.
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
- Use `drive meeting ... --protocol=...` or `meetings set-protocol` when the default `evidence-review` protocol is not the right fit for the decision surface.
- Treat `meetings orchestrate` as the council chair: once bounded positions exist, it can auto-synthesize provisional findings and unresolved-action prompts without overwriting manual judgment that is already on the record.
- Treat missing independent review as first-class wave debt: `meetings orchestrate` now reports `open_wave_action_count` so the leader can see who still owes a bounded review before cross-examination.
- Treat routed cross-examination as first-class: when multiple roles disagree on the same `--claim`, the council now turns that lane into an explicit route obligation with an assigned bounded challenger before synthesis/output is trusted.
- Treat follow-up rounds as first-class: `meetings orchestrate` can now keep assigning bounded round 2 / round 3 / ... replies until a lane resolves on authority/freshness or the configured debate-round cap is reached.
- When recording cross-examination positions, include explicit `--responds-to-run=` / `--responds-to-role=` linkage so the council trace shows who challenged whom.
- When recording positions that speak about the same question, include `--claim=` and `--dimension=` so auto-synthesis and authority audit can group them into the same council claim.
- When multiple findings speak about the same claim, use `--claim=` and `--dimension=` on `meetings add-finding` so authority conflicts can be audited cleanly.
- Use `meetings summary` before adoption when you need to see synthesis readiness, debate-round depth, blocker counts, and any still-open wave/route obligations at a glance.
- Use `meetings authority-audit` when you need a leader-facing explanation of which same-claim finding currently wins on authority.
- Use `runs authority-audit` when a run may be drifting against fresher or higher-authority active signals, or when `doctor` reports `run-authority-conflict`.
- Use `runs authority-audit` when a run claims it is ready/verified/clean but the live worktree is still dirty; the audit now surfaces that intent-vs-live-git disagreement explicitly.
- Use `missions authority-audit --all` when mission bundles overlap; it now explains bundle-level `resolution_kind`, shared selectors, and the leader-facing `required_action` instead of only listing pairwise conflicts.
- After a meeting closes, either let `meeting-close` apply its outcome or explicitly run `adopt-meeting`.
- Keep `doctor` clean. If it shows drift, repair that before piling on new execution.
