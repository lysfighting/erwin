---
name: cwork
description: GitHub-first issue to contract to execution to PR workflow for Codex, also suitable as a codex-work operating pattern. Use when Codex should shape work, claim a mission, generate or refresh an execution contract, implement against repository truth, ship a pull request, or clean up in a repo that may already use teamwork metadata, but should not depend on Claude Code slash commands or Claude-only prompt primitives.
---

# cwork

Use this skill as the Codex-native, codex-work counterpart to the historical `teamwork` flow.

Keep the protocol, but remove Claude-specific assumptions:

- operate through normal Codex conversation and tool use, not slash commands
- prefer standard repo tools such as `git`, `gh`, `jq`, `rg`, and the bundled repo-local scripts
- treat human confirmation as a normal conversation turn, not `AskUserQuestion`
- keep Codex-owned execution artifacts under `.cwork/` unless the repo already has a stronger local convention
- avoid Python runtime assumptions in the skill itself; the bundled helpers are shell-based

## Quick start

Start by orienting on the live repo state:

1. Inspect the current branch, worktree, git status, and `gh` auth.
2. Check whether the repo already has teamwork metadata in `.teamwork/` or `.teamspace/`.
3. Prefer the bundled `cwork` scripts for deterministic operations:
   - `scripts/cwork-status.sh`
   - `scripts/cwork-claim.sh`
   - `scripts/cwork-contract.sh`
   - `scripts/cwork-ship.sh`
   - `scripts/cwork-done.sh`
4. Keep Codex-owned mission contracts in `.cwork/active/MISSION-<issue>.md`.

If you need more detailed flow guidance or a teamwork-to-cwork translation map, read `references/workflow.md` and `references/commands.md`.

## Mission loop

Follow the same lifecycle each time:

1. Shape the work.
   Confirm the real objective, constraints, verification method, and blast radius before coding.
2. Bind the work to GitHub.
   Prefer an issue as the durable task record and a PR as the delivery record.
3. Claim the mission locally.
   Create or switch to the working branch, generate a mission contract, and make the local execution plan explicit.
4. Drive the implementation.
   Validate the contract against actual code, rewrite wrong sub-tasks, implement, and verify.
5. Ship the outcome.
   Push the branch, open or update the PR, surface unresolved manual steps, and sync task state.
6. Close or hand off cleanly.
   Clean local mission artifacts, update labels or issue state, and leave the next operator with accurate status.

## Contract workflow

When an issue already exists, generate or refresh a contract with the bundled helper:

- use `scripts/cwork-contract.sh from-issue ISSUE_NUMBER` when `gh` access is available
- use `scripts/cwork-contract.sh from-json path/to/issue.json` for offline testing or fixture-driven work

The helper writes a default contract to `.cwork/active/MISSION-<issue>.md` unless you override the output path.

After generating a contract:

1. Read it back before coding.
2. Correct any bad assumptions introduced by the issue text or by the parser.
3. Add missing context files, risks, and verification notes.
4. Treat sub-tasks as hypotheses, not orders.

## Operational commands

Use these entrypoints instead of reconstructing the flow from scratch every time:

- `bash skills/cwork/scripts/cwork-status.sh`
  Show the current branch, active contract, freshness, and any open PR tied to the branch.
- `bash skills/cwork/scripts/cwork-claim.sh ISSUE_NUMBER`
  Claim the issue, generate `.cwork/active/MISSION-<issue>.md`, create the working branch, and optionally enter worktree mode.
- `bash skills/cwork/scripts/cwork-ship.sh`
  Build the PR body from the contract, push the branch, create or update the PR, and move the issue toward review.
- `bash skills/cwork/scripts/cwork-done.sh`
  Clean up the active contract after merge and move the issue toward done.

When you need the script behavior explained in detail, read `references/commands.md`.

## Execution guardrails

- Prefer repo truth over stale issue text.
- Prefer the narrowest branch, diff, and ownership boundary that can finish the job.
- Reuse teamwork scripts as deterministic helpers, but do not require the user to invoke `/team`, `/team-claim`, or other Claude-specific commands.
- If the repo still has historical teamwork metadata, honor it when it improves compatibility, but do not call Python-dependent teamwork helpers from this skill.
- Preserve compatible concepts when useful:
  - issues as task source of truth
  - contracts as AI execution view
  - PRs as delivery boundary
  - worktrees as optional isolation
- If the repo is not initialized for teamwork, fall back to plain `git` plus `gh` instead of failing.
- If you detect high-blast-radius ambiguity, pause and confirm before implementation.

## Shipping checklist

Before opening or updating a PR:

1. Confirm the current branch is not a protected default branch.
2. Run the relevant verification commands or state clearly why they were not run.
3. Sync completed checklist items back to the issue when the repo uses mission-style checkboxes.
4. Include any remaining manual or follow-up work explicitly in the PR body or handoff note.

## Cleanup

At the end of the mission:

1. Delete or archive obsolete `.cwork/active/MISSION-<issue>.md` files when they no longer represent active work.
2. Leave the issue and PR state aligned with reality.
3. Update any repo-local status documents when the repository protocol requires it.
