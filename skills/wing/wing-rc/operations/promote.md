# Promote Flow — RC → Production

> Squash merge RC to production, tag, create GitHub Release, close milestone, deploy.

## P1: Pre-flight

Same config as Prepare Step 1 (worktree guard, REPO, VERSION_PREFIX, BASE_BRANCH, PROD_BRANCH, workflows).

## P2: Find Active RC

```bash
RC_BRANCH=$(bash "${CLAUDE_SKILL_DIR}/../scripts/wing-git.sh" find-rc 2>/dev/null)
RC_EXIT=$?
```

- None found → "No active RC. Run `/wing-rc` first." → **STOP**
- Multiple → "Multiple RC branches. Delete stale ones." → **STOP**

Extract version: `VERSION=$(echo "$RC_BRANCH" | sed 's|rc/||')`

## P3: Staging Check

If `STAGING_WORKFLOW` set → check latest run conclusion. If not success → AskUserQuestion: "继续" / "中止".
If not set → warn, skip.

## P4: Confirmation

Display summary: RC branch, target, actions. AskUserQuestion: "Promote" / "中止".

## P5: Generate Release Notes

```bash
PREV_TAG=$(git tag --sort=-v:refname --list "${VERSION_PREFIX}*" \
  | grep -E "^$(echo "$VERSION_PREFIX" | sed 's/[.[\*^$()+?{|\\]/\\&/g')\.[0-9]+$" \
  | head -1)
PREV_DATE=$(git log -1 --format='%aI' "$PREV_TAG" 2>/dev/null || echo "2000-01-01T00:00:00Z")
PR_JSON=$(gh pr list --base "$BASE_BRANCH" --state merged \
  --json number,title,author,mergedAt,url --limit 300 --repo "$REPO")
```

Filter by date (mergedAt > PREV_DATE). Categorize by commit type (feat/fix/other). Format:

```
## What's in this release
### Features
* feat: {title} by @{author} in {url} (#NNN)
### Bug Fixes
* fix: {title} by @{author} in {url} (#NNN)
### Contributors
@alice (3) · @bob (1)
**Full Changelog**: {PREV_TAG}...{VERSION}
```

## P6: Create PR + Squash Merge

```bash
gh pr create --base "$PROD_BRANCH" --head "$RC_BRANCH" \
  --title "release: $VERSION" --body "$RELEASE_BODY" --repo "$REPO"
PR_NUMBER=$(gh pr view "$RC_BRANCH" --json number --jq '.number' --repo "$REPO")
```

Try auto-merge first:
```bash
gh pr merge "$PR_NUMBER" --squash --auto \
  --subject "release: $VERSION" --body "$RELEASE_BODY" --repo "$REPO"
bash "${CLAUDE_SKILL_DIR}/../scripts/wing-pr.sh" wait-merged "$PR_NUMBER" "$REPO"
```

If `--auto` fails → watch CI → direct merge fallback.

## P7: Tag on Main

```bash
git fetch origin "$PROD_BRANCH"
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "Tag $VERSION already exists — skipping"
else
  bash "${CLAUDE_SKILL_DIR}/../scripts/wing-git.sh" tag "$VERSION" "origin/$PROD_BRANCH"
fi
```

## P8: GitHub Release

```bash
gh release create "$VERSION" --notes-file <(echo "$RELEASE_BODY") \
  --title "$VERSION" --repo "$REPO"
```

Idempotent: skip if release already exists. Notify via `wing-notify.sh milestone.done`.

## P9: Milestone Check

Query milestone by resolved title. If all issues closed → close milestone. If open remain → display counts.

## P10: Cleanup + Deploy

Delete RC branch (local + remote). Trigger production deploy if `PROD_WORKFLOW` configured, else guide manual deploy.

## P11: Summary

**PROMOTED ── {VERSION}** ────────────────────
Tag: **{VERSION}** (on *{PROD_BRANCH}*)
PR: **#{PR_NUMBER}** (squash merged)
Release: {url}
Milestone: {status}
Deploy: {triggered / manual}

`/wing-rc` to prepare next RC
