# Harness

Use this skill when you are operating the repository-local orchestration harness as the leader.

The harness is file-backed:

- committed definitions live in `harness/`
- runtime state lives in `.harness/`
- command surfaces live in `scripts/`

## Leader workflow

1. Start with the home view:
   `bash scripts/harness-home.sh`
2. If anything looks off, inspect runtime health:
   `bash scripts/harness-doctor.sh`
3. Create or draft executable issues when work is not yet claimable:
   `bash scripts/harness-issue.sh ...`
4. Claim work into governed execution:
   `bash scripts/harness-claim.sh ISSUE_NUMBER [TITLE_OVERRIDE]`
5. Drive run state:
   `bash scripts/harness-drive.sh ...`
6. Prepare release handoff when execution is done:
   `bash scripts/harness-ship.sh ...`
7. Use meetings when ambiguity or blast radius is high:
   `bash scripts/harness-meetings.sh ...`
8. Use RC flow for release candidate governance:
   `bash scripts/harness-rc.sh ...`
9. Repair drift or normalize runtime state when needed:
   `bash scripts/harness-maintenance.sh ...`

## Important commands

- `home`
  `bash scripts/harness-home.sh`
- `doctor`
  `bash scripts/harness-doctor.sh`
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

## Operating rules

- Treat `harness-home` as the control plane entry point.
- Prefer the narrowest safe state change.
- Use meetings for high-ambiguity or high-blast-radius decisions.
- After a meeting closes, either let `meeting-close` apply its outcome or explicitly run `adopt-meeting`.
- Keep `doctor` clean. If it shows drift, repair that before piling on new execution.
