# Hotfix Flow — Production Fix + Backport

> Fix production on main, then backport to develop. Requires dual-branch mode.

## H1: Pre-flight

Same config as Prepare (worktree guard, REPO, VERSION_PREFIX, BASE_BRANCH, PROD_BRANCH).

**Single-branch guard:** If `BASE_BRANCH == PROD_BRANCH` → "Hotfix requires dual-branch mode." → **STOP**

## H2: Derive Hotfix Version

```bash
git fetch origin "$PROD_BRANCH" 2>/dev/null
git fetch --tags origin 2>/dev/null || true
LATEST_TAG=$(git tag --sort=-v:refname --list "${VERSION_PREFIX}*" 2>/dev/null | head -1)
```

- No tags → "No release tags. Run `/wing-rc` + `/wing-rc promote` first." → **STOP**

Derive: `V0.1.0` → `V0.1.0-hotfix.1`. If exists → increment.

```bash
HOTFIX_NUM=1
while git tag --list "${LATEST_TAG}-hotfix.${HOTFIX_NUM}" 2>/dev/null | grep -q .; do
  HOTFIX_NUM=$((HOTFIX_NUM + 1))
done
HOTFIX_VERSION="${LATEST_TAG}-hotfix.${HOTFIX_NUM}"
HOTFIX_BRANCH="hotfix/${HOTFIX_VERSION}"
```

## H3: Create Hotfix Branch

If already on hotfix branch → resume. Otherwise:

```bash
git checkout -b "$HOTFIX_BRANCH" "origin/$PROD_BRANCH"
```

Display briefing, then interactive fix phase.

## H4: Fix Phase (interactive)

Guide user to make their fix. Detect uncommitted work. AskUserQuestion: "Commit and continue" / "Not yet" (resume later with `/wing-rc hotfix`).

```bash
git add -u
git commit -m "hotfix(${HOTFIX_VERSION}): fix description"
```

## H5: PR to Main

```bash
git push origin "$HOTFIX_BRANCH"
gh pr create --base "$PROD_BRANCH" --head "$HOTFIX_BRANCH" \
  --title "hotfix: ${HOTFIX_VERSION}" \
  --body "## Hotfix\n\nBranch: \`${HOTFIX_BRANCH}\`\nBase tag: \`${LATEST_TAG}\`" \
  --repo "$REPO"
```

Auto-merge → poll. Fallback: watch CI → direct merge.

## H6: Tag on Main

```bash
git fetch origin "$PROD_BRANCH"
bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" tag "$HOTFIX_VERSION" "origin/$PROD_BRANCH"
```

Create GitHub Release with hotfix notes.

## H7: Backport to Develop

```bash
MERGE_SHA=$(gh pr view "$PR_NUMBER" --json mergeCommit --jq '.mergeCommit.oid' --repo "$REPO")
git checkout "$BASE_BRANCH"
git pull origin "$BASE_BRANCH"
git cherry-pick "$MERGE_SHA" --no-edit
```

If conflicts → guide resolution interactively. Push after resolution.

## H8: Cleanup + Summary

Delete hotfix branch (local + remote).

**HOTFIX COMPLETE ── {HOTFIX_VERSION}** ────────────────────
Tag: **{HOTFIX_VERSION}** (on *{PROD_BRANCH}*)
PR: **#{PR_NUMBER}** (squash merged)
Backport: cherry-picked to *{BASE_BRANCH}*

{If active RC detected:}
⚠️ Active RC: *{RC_BRANCH}*. Consider `/wing-rc fix` or `/wing-rc cancel` + `/wing-rc`.
