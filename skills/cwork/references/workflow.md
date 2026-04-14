# cwork Workflow Reference

Use this reference when the hot-path instructions in `SKILL.md` are not enough.

## Purpose

`cwork` preserves the useful protocol from `teamwork`:

- GitHub issues hold durable task intent
- local contracts hold the AI execution view
- pull requests hold delivery
- optional worktrees isolate parallel work

It intentionally drops Claude-only machinery such as slash-command routing, `AskUserQuestion`, and `CLAUDE_SKILL_DIR` path assumptions.

## Translation map

| Historical teamwork surface | cwork equivalent |
|---|---|
| `/team` dashboard | inspect `git`, `gh`, repo config, and active mission files directly |
| `/team-issue` | create or update issues through `gh` or the GitHub connector |
| `/team-claim` | create branch or worktree, then generate `.cwork/active/MISSION-<issue>.md` |
| `/team-drive` | read contract, validate against code, implement, test, update contract |
| `/team-ship` | push branch, open or update PR, report verification, clean local mission state |
| `/team-ship done` | close task loop, clean branch or worktree, remove stale contract |

## Repo-local helpers worth reusing

When the repo still carries historical teamwork metadata, treat it as compatibility input rather than the execution engine:

- `.teamwork/config.yml` or `.teamspace/config.yml`
  Use as compatibility input for base branch, branch pattern, and label prefixes.
- `git config --local teamwork.worktree`
  Honor existing worktree preferences when they are already established in the repo.

Do not make `cwork` depend on Python-backed teamwork helpers. The `cwork` scripts are intended to stand on their own.

## Contract conventions

Default contract path:

`.cwork/active/MISSION-<issue>.md`

Recommended frontmatter keys:

- `issue`
- `title`
- `url`
- `branch`
- `priority`
- `milestone`
- `assignee`
- `generated_at`
- `issue_content_hash`

Recommended sections:

- `## Objective`
- `## Sub-tasks`
- `## Success Criteria`
- `## Context & References`
- `## Verification`
- `## AI Notes`

## Branch conventions

If the repo already defines a branch pattern in teamwork config, follow it.

If it does not, prefer:

- `mission/<issue>-<slug>` for normal feature work
- `fix/<issue>-<slug>` for clearly bug-scoped work

If the current branch already isolates the task cleanly, reuse it instead of churning branch names.

## Operating posture

- Challenge bad issue framing before coding.
- Keep the contract fresh when issue intent changes.
- Record what was actually verified, not what you expected to be true.
- Do not silently preserve stale mission state after merge or abandonment.
