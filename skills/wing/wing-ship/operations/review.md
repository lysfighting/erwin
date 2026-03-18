# Operation Review — AI Code Review

> Triggered by `/wing-ship review`. Fetches PR diff, reviews code, publishes review.

## R1: Find PR

```bash
BRANCH=$(git branch --show-current)
gh pr list --head "$BRANCH" --state open --json number,url --limit 1
```

If `$ARGUMENTS` contains a PR number → use that instead.
If no open PR → "No open PR for branch `{branch}`." → **STOP**
  💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

## R2: Get PR diff + Read changed files

```bash
gh pr diff {pr}
```

For each file in the diff, read full current content for context beyond diff hunks.

## R3: AI Review

Analyze changes for:
- **Correctness**: Logic errors, missing edge cases, off-by-one
- **Security**: Injection, exposed secrets, unsafe operations
- **Architecture**: Coupling, separation of concerns, pattern consistency
- **Quality**: Naming, readability, dead code, duplication
- **Performance**: N+1 queries, unnecessary allocations

## R4: Choose review action

Use `AskUserQuestion`: "Approve" / "Request changes" / "Comment only"

## R5: Publish review

```bash
gh pr review {pr} --body "{review_body}" {--approve|--request-changes|--comment}
```

Output: Review published on PR **#{pr}**: {action}
After merge: `/wing-ship done` │ Address comments: push more commits
