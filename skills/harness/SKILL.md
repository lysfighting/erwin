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
4. Create or draft executable issues when work is not yet claimable:
   `bash scripts/harness-issue.sh ...`
5. Claim work into governed execution:
   `bash scripts/harness-claim.sh ISSUE_NUMBER [TITLE_OVERRIDE]`
6. Drive run state:
   `bash scripts/harness-drive.sh ...`
7. Delegate bounded work when leader execution should split:
   `bash scripts/harness-teammates.sh ...`
8. Prepare release handoff when execution is done:
   `bash scripts/harness-ship.sh ...`
9. Use meetings when ambiguity or blast radius is high:
   `bash scripts/harness-meetings.sh ...`
10. Use RC flow for release candidate governance:
   `bash scripts/harness-rc.sh ...`
11. Repair drift or normalize runtime state when needed:
   `bash scripts/harness-maintenance.sh ...`

## Important commands

- `home`
  `bash scripts/harness-home.sh`
- `doctor`
  `bash scripts/harness-doctor.sh`
- `adapter-validate`
  `bash scripts/harness-adapter.sh validate`
- `issue-draft`
  `bash scripts/harness-issue.sh draft "Mission title" --criteria="A|B"`
- `claim`
  `bash scripts/harness-claim.sh 1234`
- `start`
  `bash scripts/harness-drive.sh start RUN_ID`
- `block`
  `bash scripts/harness-drive.sh block RUN_ID "reason"`
- `meeting`
  `bash scripts/harness-drive.sh meeting RUN_ID "question"`
- `teammate-assign`
  `bash scripts/harness-teammates.sh assign RUN_ID worker-auth implementation-verifier "Update auth callers" --scope="auth callers only"`
- `teammate-report`
  `bash scripts/harness-teammates.sh report TEAMMATE_RUN_ID "Caller patch is in place" --kind=progress --evidence=src/auth.ts|tests/auth.spec.ts`
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
  `bash scripts/harness-meetings.sh add-position MEETING_ID TEAMMATE_RUN_ID independent_review challenge "This contract leaks release policy" --evidence=docs/spec.md`
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
- `all-smoke`
  `bash scripts/harness-all-smoke-test.sh`
- `runtime-dist-smoke`
  `bash scripts/harness-runtime-dist-smoke-test.sh`

## Operating rules

- Treat `harness-home` as the control plane entry point.
- Treat `harness-adapter.sh validate` as the package boundary check when bootstrap/sync behavior is in doubt.
- Treat `harness-runtime-dist-smoke-test.sh` as the proof that a repo can run the packaged runtime without a local Go toolchain.
- Prefer the narrowest safe state change.
- Use meetings for high-ambiguity or high-blast-radius decisions.
- Treat teammate runs as first-class runtime objects, not informal notes.
- Allow teammate-to-teammate communication only through governed links or sibling relationships.
- After a meeting closes, either let `meeting-close` apply its outcome or explicitly run `adopt-meeting`.
- Keep `doctor` clean. If it shows drift, repair that before piling on new execution.
