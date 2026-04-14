# cwork Command Surface

Use these scripts when you want deterministic support around the Codex-native work loop.

## Core commands

`bash skills/cwork/scripts/cwork-status.sh`

- show current branch
- show active contract path when exactly one contract is discoverable
- show freshness against the backing GitHub issue
- show any open PR for the current branch

`bash skills/cwork/scripts/cwork-claim.sh ISSUE_NUMBER`

- fetch issue details from GitHub
- derive slug and branch name
- generate `.cwork/active/MISSION-<issue>.md`
- create or switch to the working branch
- create a worktree when `git config --local cwork.worktree true`

Options:

- `--json path/to/issue.json` for offline testing
- `--branch custom/branch-name` to override the branch recorded in the contract
- `--no-comment` to skip posting a claim comment back to the issue

`bash skills/cwork/scripts/cwork-ship.sh`

- find the active contract
- render a PR body from the contract
- optionally sync checked contract items back to the issue
- push the current branch
- create or update the PR

Options:

- `--contract PATH` to target a specific contract
- `--issue N` to help contract lookup
- `--draft` to open a draft PR
- `--update-existing` to rewrite title and body on an already-open PR
- `--dry-run` to print the PR payload without pushing or touching GitHub state
- `--no-sync` to skip syncing checked items back to the issue

`bash skills/cwork/scripts/cwork-done.sh`

- verify the branch has a merged PR unless `--force`
- mark the issue `done` when label config exists
- remove the local mission contract
- optionally close the issue if auto-close did not happen

Options:

- `--contract PATH`
- `--issue N`
- `--force`
- `--close-issue`

## Lower-level helpers

`bash skills/cwork/scripts/cwork-config.sh get conventions.base_branch main`

- read compatible teamwork config without Python

`bash skills/cwork/scripts/cwork-git.sh build-branch-name 42 add-login codex`

- derive branch names from repo config or defaults

`bash skills/cwork/scripts/cwork-contract.sh pr-body .cwork/active/MISSION-42.md`

- render the default PR body from a mission contract

`bash skills/cwork/scripts/cwork-contract.sh sync-checked .cwork/active/MISSION-42.md 42`

- sync exact-match checked items back to the issue body
