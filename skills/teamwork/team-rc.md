---
description: "RC lifecycle — prepare staging or promote to production. Try: /team-rc help"
version: "3.16.0"
---

# /team-rc — Release Candidate Lifecycle

> You own the release cadence. Every RC is a contract with users — what you ship defines trust. Verify before promoting, and never rush a release you wouldn't bet your reputation on.

**Visual Encoding** (apply to ALL output — see `docs/visual-encoding-standard.md`):
`**bold**` → headers/labels (white) · `` `backtick` `` → commands/paths/counts (purple-blue) · `*italic*` → branches (dim) · `**#NNN**` → issues (light-blue clickable, 3+ digits) · `────` dividers · ⛔ NO code blocks around output

**User input**: $ARGUMENTS

**Argument routing:**

| Input | Action |
|-------|--------|
| (empty) | Prepare: cut rc branch from develop → staging |
| `cancel` | Cancel: delete current RC branch (fix on base, then re-cut) |
| `promote` | Promote: squash merge rc → main → tag → GitHub Release |
| `fix` | Fix: cherry-pick a specific commit onto active RC |
| `fix help` | Show fix operation usage |
| `hotfix` | Hotfix: fix production on main + backport to develop |
| `hotfix help` | Show hotfix operation usage |
| `promote help` | Show promote operation usage |
| `help` or `-h` | Show usage guide |

---

## Route by Argument

Parse `$ARGUMENTS`:

- If `help` or `-h` → output the following and **STOP**:

**`/team-rc` — Release Candidate Lifecycle**
*Author: liyasong + casey | 2026-03-05*

**USAGE**
  `/team-rc`              Prepare: cut rc branch, deploy staging
  `/team-rc cancel`       Cancel: delete current rc, fix on base, re-cut
  `/team-rc promote`      Promote: squash merge rc → *main*, tag, release
  `/team-rc help`         Show this guide

**PREPARE** (`/team-rc`)
  1. Check no active rc branch on remote
  2. Derive next version from remote tags
  3. Verify *develop* CI green
  4. Cut *rc/{version}* from *develop*, push
  5. Trigger staging deploy (if configured)

**CANCEL** (`/team-rc cancel`)
  Delete current RC branch (local + remote).
  Use when staging found bugs — fix on *develop*, then `/team-rc` to re-cut.

**PROMOTE** (`/team-rc promote`)
  1. Find active rc branch on remote
  2. Check staging status
  3. Create PR: rc → *main* (squash merge via GitHub)
  4. Tag squash commit on *main*
  5. Create GitHub Release
  6. Close milestone if complete
  7. Delete rc branch

**FIX** (`/team-rc fix`)
  Cherry-pick a specific commit from *develop* to active RC.
  Preserves test progress (unlike cancel + re-cut).
  `/team-rc fix`              Select from recent commits
  `/team-rc fix <commit>`     Cherry-pick specific commit hash
  `/team-rc fix <branch>`     Cherry-pick HEAD of branch

**HOTFIX** (`/team-rc hotfix`)
  Post-release production fix. Requires dual-branch mode.
  1. Create *hotfix/{version}* from *main*
  2. Make fix (interactive)
  3. PR to *main* (squash merge)
  4. Tag + GitHub Release
  5. Backport (cherry-pick) to *develop*

**BUGFIX DURING RC** (upstream-first model)
  RC is a frozen snapshot. Default: fix on *develop*, cancel + re-cut.
  Found a bug on staging? →
  Option A (default): Fix on *develop* → `/team-rc cancel` → `/team-rc` (re-cut)
  Option B (surgical): Fix on *develop* → `/team-rc fix` (cherry-pick to RC)

**CONFIG**
  `versions.current`: `"V0.1"`              Milestone prefix (patch auto-derived from tags)
  `conventions.base_branch`: `"develop"`     Development trunk
  `conventions.production_branch`: `"main"`  Production branch
  `deploy.staging_workflow`: `"deploy.yml"`  *(optional)* Staging deploy workflow name
  `deploy.production_workflow`: `"deploy.yml"` *(optional)* Production deploy workflow name

- If `fix help` → jump to **Operation Fix Help**
- If `fix` → jump to **Fix Flow**
- If `hotfix help` → jump to **Operation Hotfix Help**
- If `hotfix` → jump to **Hotfix Flow**
- If `promote help` → jump to **Operation Promote Help**
- If `promote` → jump to **Promote Flow**
- If `cancel` → jump to **Cancel Flow**
- Otherwise → **Prepare Flow**

---

## Prepare Flow

### Step 1: Pre-flight

**Worktree guard** — `/team-rc` manages releases from the main repo, not mission worktrees:

```bash
if [ -f .mission ]; then
  MAIN_REPO=$(bash ~/.claude/commands/scripts/tw-git.sh worktree-main-repo)
  echo "⚠️ /team-rc must run from the main repo, not a mission worktree."
  echo "Switch to main repo: cd $MAIN_REPO"
  echo "Then retry: /team-rc"
fi
```

If worktree detected → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

Detect config directory:

```bash
TEAMWORK_DIR=$(bash ~/.claude/commands/scripts/tw-config.sh detect-dir 2>/dev/null) || {
  echo "No teamwork config found. Run /team init first."
  # STOP
}
```

- If no config → "**ERROR:** Teamwork not initialized. Run `/team` first." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

Read required values:

```bash
REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
if [ -z "$REPO" ]; then
  echo "ERROR: Cannot detect repo. Run: gh auth login"
  # STOP
fi

VERSION_PREFIX=$(bash ~/.claude/commands/scripts/tw-config.sh versions.current "" 2>/dev/null)
if [ -z "$VERSION_PREFIX" ]; then
  echo "ERROR: versions.current not set. Add to $TEAMWORK_DIR/config.yml:"
  echo "  versions:"
  echo "    current: \"V0.1\""
  # STOP
fi

BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "develop" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-develop}"

# Warn if base_branch is not explicitly configured — default "develop" may conflict
# with other skills that default to "main". Require explicit config for RC flow.
BASE_BRANCH_RAW=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "" 2>/dev/null)
if [ -z "$BASE_BRANCH_RAW" ]; then
  echo "⚠️ WARNING: conventions.base_branch not set in config.yml (defaulting to 'develop')."
  echo "  /team-rc assumes dual-branch (develop→main). Other skills default to 'main'."
  echo "  Set explicitly in $TEAMWORK_DIR/config.yml:"
  echo "    conventions:"
  echo "      base_branch: \"develop\""
fi

PROD_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.production_branch "main" 2>/dev/null)
PROD_BRANCH="${PROD_BRANCH:-main}"

STAGING_WORKFLOW=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.staging_workflow "" 2>/dev/null)
PROD_WORKFLOW=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.production_workflow "" 2>/dev/null)
```

- If no repo detected → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- If `versions.current` not set → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### Step 2: RC Uniqueness Check

```bash
RC_BRANCH=$(bash ~/.claude/commands/scripts/tw-git.sh find-rc 2>/dev/null)
RC_EXIT=$?
# RC_EXIT: 0=found one (RC_BRANCH has name), 1=none found, 3=multiple found
```

**If `RC_EXIT` is 0 (rc branch exists)** → show status and offer choices:

**RC in progress:** *{RC_BRANCH}*

Use `AskUserQuestion`:
```
question: "已有 RC 分支 {RC_BRANCH}。如何操作？"
options:
  - label: "Promote（发布到 production）"
    description: "RC 验证通过，执行 /team-rc promote"
  - label: "Cancel（取消 RC，重新来）"
    description: "删除 RC 分支，在 base 修复后重新切"
  - label: "退出"
    description: "不操作"
```

- If "Promote" → jump to **Promote Flow** (skip P1 pre-flight, config already loaded)
- If "Cancel" → jump to **Cancel Flow** Step C3 (RC_BRANCH already found)
- If "退出" → **STOP**

### Step 3: Derive Next Version

```bash
NEXT_VERSION=$(bash ~/.claude/commands/scripts/tw-git.sh next-version "$VERSION_PREFIX")
```

### Step 4: Verify Develop CI

```bash
LATEST_CI=$(gh run list --branch "$BASE_BRANCH" --limit 1 \
  --json conclusion --jq '.[0].conclusion' 2>/dev/null)
```

- If `success` → `Develop CI: green` → continue to Step 5
- If not `success` → warn: `Develop CI not green (${LATEST_CI}). Proceed with caution.`

Use `AskUserQuestion`:
```
question: "Develop CI 未通过 ({LATEST_CI})。是否继续？"
options:
  - label: "继续"
    description: "CI 问题不影响本次 RC，继续切分支"
  - label: "中止"
    description: "先修 CI，之后重新 /team-rc"
```
If "中止" → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### Step 4b: Confirm RC Preparation

Display the RC summary and confirm before cutting:

**PREPARE RC** ── **{NEXT_VERSION}** ────────────────────
*rc/{NEXT_VERSION}* ← *{BASE_BRANCH}*
Develop CI: {green / not green}
Staging: {workflow name / "(No staging workflow configured — deploy manually)"}
────────────────────────────────────────────

Use `AskUserQuestion`:
```
question: "确认切 RC 分支？"
options:
  - label: "确认"
    description: "从 {BASE_BRANCH} 切出 rc/{NEXT_VERSION}"
  - label: "取消"
    description: "放弃本次 RC"
```
If "取消" → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### Step 5: Cut RC Branch

```bash
bash ~/.claude/commands/scripts/tw-git.sh cut-release "$NEXT_VERSION"
```

This runs: fetch origin/develop → create `rc/$NEXT_VERSION` from it → push to origin. **Does NOT checkout or modify the working tree.**

### Step 6: Trigger Staging Deploy

**If `STAGING_WORKFLOW` is set:**

```bash
gh workflow run "$STAGING_WORKFLOW" --ref "rc/$NEXT_VERSION"
```

Output: Staging build triggered for *rc/{NEXT_VERSION}* (Build Images → deploy-staging.yml → ArgoCD)

**If not set:** Output: *(No staging workflow configured — deploy manually)*

### Step 7: Summary

**RC PREPARED ── {NEXT_VERSION}** ────────────────────
*rc/{NEXT_VERSION}*  ← *{BASE_BRANCH}*
Staging:  {triggered / deploy manually}

**NEXT STEPS**
  1. 部署 staging，验证
  2. 发现 bug? → 在 *{BASE_BRANCH}* 修复，然后:
     Option A (clean slate): `/team-rc cancel` → `/team-rc` 重新切
     Option B (surgical): `/team-rc fix` cherry-pick 到 RC（保留测试进度）
  3. 验证通过 → `/team-rc promote`

────────────────────────────────────────────

**Generate 飞书 announcement**:

```bash
# Pick previous clean version tag as the date boundary
PREV_TAG=$(git tag --sort=-v:refname --list "${VERSION_PREFIX}*" \
  | grep -E "^$(echo "$VERSION_PREFIX" | sed 's/[.[\*^$()+?{|\\]/\\&/g')\.[0-9]+$" \
  | head -1)

# Base sha for the metadata block
BASE_SHA=$(git rev-parse "origin/$BASE_BRANCH" 2>/dev/null)

# Emit the 飞书-ready markdown block (uses shared tw-pr.sh merged-since)
FEISHU_BLOCK=$(bash ~/.claude/commands/scripts/tw-pr.sh merged-since \
  "$BASE_BRANCH" "$PREV_TAG" \
  --repo "$REPO" \
  --format feishu \
  --version "$NEXT_VERSION" \
  --base-sha "$BASE_SHA" \
  --env staging)
```

Then output the announcement block between dividers so user can copy-paste:

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📣 **飞书发布通知 — 复制下方内容**
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

{$FEISHU_BLOCK content}

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

{^ If PREV_TAG is empty (first release), merged-since returns all PRs back to repo start; the output still works, just shows "Initial release" marker instead of a diff link.}
{^ If merged-since fails (no gh auth, no network, no base branch on origin) → warn "FEISHU block skipped" but do NOT block the flow; the RC cut has already happened and is the critical work.}

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

---

## Cancel Flow

> Triggered by `/team-rc cancel`. Delete the current RC branch so you can fix on base and re-cut.

### Step C1: Pre-flight

Same prerequisites as Promote — config, `gh`, git remote.

```bash
TEAMWORK_DIR=$(bash ~/.claude/commands/scripts/tw-config.sh detect-dir 2>/dev/null) || {
  echo "No teamwork config found. Run /team init first."
  # STOP
}

BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "develop" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-develop}"
REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
```

### Step C2: Find active RC branch

```bash
RC_BRANCH=$(bash ~/.claude/commands/scripts/tw-git.sh find-rc 2>/dev/null)
```

If no RC branch found → "No active RC branch to cancel." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

If multiple found → "Multiple RC branches found. Run `/team doctor fix` to clean up, or delete manually: `git push origin --delete rc/BRANCH` for each extra." → **STOP**

### Step C3: Show RC info + confirm

```bash
# Show what's on the RC branch
RC_COMMITS=$(git log "origin/$BASE_BRANCH..origin/$RC_BRANCH" --oneline --no-merges 2>/dev/null)
RC_COMMIT_COUNT=$(echo "$RC_COMMITS" | grep -c . 2>/dev/null || echo 0)
```

**CANCEL RC** ── *{RC_BRANCH}* ────────────────────
Commits on RC: `{RC_COMMIT_COUNT}`
{RC_COMMITS as indented list}

⚠️ This will delete *{RC_BRANCH}* (local + remote). Commits are preserved in *{BASE_BRANCH}*.
────────────────────────────────────────────

Use `AskUserQuestion`:
```
question: "确认取消 RC？分支 {RC_BRANCH} 将被删除（local + remote）。"
options:
  - label: "确认取消"
    description: "删除 RC 分支，之后可在 {BASE_BRANCH} 修复后重新 /team-rc"
  - label: "保留"
    description: "不取消，继续 staging 验证"
```

If "保留" → **STOP**

### Step C4: Delete RC branch

```bash
# Delete remote
git push origin --delete "$RC_BRANCH" 2>/dev/null || {
  echo "WARNING: Could not delete remote branch (may not exist)" >&2
}

# Delete local
git branch -d "$RC_BRANCH" 2>/dev/null || git branch -D "$RC_BRANCH" 2>/dev/null || true
```

### Step C5: Summary

**RC CANCELLED** ── *{RC_BRANCH}* ────────────────────
Branch deleted (local + remote).

**Next:**
  1. Fix the issue on *{BASE_BRANCH}* (normal mission flow)
  2. `/team-rc` to re-cut RC from latest *{BASE_BRANCH}*
  3. Re-deploy staging
────────────────────────────────────────────

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

---

## Fix Flow

> Triggered by `/team-rc fix`. Cherry-pick a specific commit from the base branch onto the active RC branch.
> Use when: a bug was fixed on develop and you need it on RC without re-cutting (preserves existing RC test progress).

### Step F1: Pre-flight

Same worktree guard + config reads as Prepare Step 1 (worktree check, REPO, BASE_BRANCH, PROD_BRANCH).

```bash
TEAMWORK_DIR=$(bash ~/.claude/commands/scripts/tw-config.sh detect-dir 2>/dev/null) || {
  echo "No teamwork config found. Run /team init first."
  # STOP
}

REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
if [ -z "$REPO" ]; then
  echo "ERROR: Cannot detect repo. Run: gh auth login"
  # STOP
fi

BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "develop" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-develop}"
```

### Step F2: Find Active RC Branch

```bash
RC_BRANCH=$(bash ~/.claude/commands/scripts/tw-git.sh find-rc 2>/dev/null)
RC_EXIT=$?
```

- If `RC_EXIT` is 1 (none found) → "No active RC branch. Cut one first with `/team-rc`." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- If `RC_EXIT` is 3 (multiple found) → "Multiple RC branches found. Delete stale ones first." → **STOP**

### Step F3: Identify the Fix Commit

Parse `$ARGUMENTS` after `fix`:
- If a commit hash or branch name is provided → validate it exists
- If no argument → show recent commits on BASE_BRANCH and let user select

**If no argument provided:**

```bash
git fetch origin "$BASE_BRANCH" 2>/dev/null
RECENT_COMMITS=$(git log "origin/$BASE_BRANCH" --oneline -10 2>/dev/null)
```

Display recent commits and use `AskUserQuestion`:
```
question: "选择要 cherry-pick 到 RC 的 commit："
options:
  - (show up to 4 recent commits as options, each with hash + message)
```

**Validate commit:**

```bash
# Check commit exists
git cat-file -t "$COMMIT_SHA" 2>/dev/null || {
  echo "ERROR: Commit $COMMIT_SHA not found."
  # STOP
}

# Check commit is not already on RC
git fetch origin "$RC_BRANCH" 2>/dev/null
if git merge-base --is-ancestor "$COMMIT_SHA" "origin/$RC_BRANCH" 2>/dev/null; then
  echo "Commit $COMMIT_SHA is already on $RC_BRANCH."
  # STOP
fi
```

### Step F4: Confirm Cherry-pick

```bash
COMMIT_MSG=$(git log --format='%h %s' -1 "$COMMIT_SHA" 2>/dev/null)
```

**CHERRY-PICK TO RC** ────────────────────
Commit: `{COMMIT_MSG}`
Target: *{RC_BRANCH}*
⚠️ This modifies the RC branch directly. Staging re-deploy needed after.
────────────────────────────────────────────

Use `AskUserQuestion`:
```
question: "确认 cherry-pick 到 RC？"
options:
  - label: "确认"
    description: "Cherry-pick {COMMIT_SHA} 到 {RC_BRANCH}"
  - label: "取消"
    description: "放弃操作"
```

If "取消" → **STOP**

### Step F5: Execute Cherry-pick

```bash
# Guard: stash uncommitted changes before switching branches
DIRTY=$(git status --porcelain 2>/dev/null | head -1)
if [ -n "$DIRTY" ]; then
  CURRENT=$(git branch --show-current)
  git stash push -m "team-rc-fix: auto-stash from $CURRENT"
  echo "⚠️ Stashed uncommitted changes. Will restore after cherry-pick."
fi

git fetch origin "$RC_BRANCH" 2>/dev/null
git checkout "$RC_BRANCH"
git pull origin "$RC_BRANCH"
git cherry-pick "$COMMIT_SHA" --no-edit
```

**If cherry-pick conflicts:**

Do NOT stop. Guide the user through resolution:

**CHERRY-PICK CONFLICT** ────────────────────
Conflicted files:
```bash
git diff --name-only --diff-filter=U
```

**To resolve:**
  1. Edit the conflicted files (conflict markers: `<<<<<<<`, `=======`, `>>>>>>>`)
  2. Stage resolved files: `git add <file>`
  3. Continue: `git cherry-pick --continue`
  4. Or abort: `git cherry-pick --abort`
────────────────────────────────────────────

Help the user resolve conflicts interactively. After resolution, continue to Step F6.

**If cherry-pick succeeds:**

```bash
git push origin "$RC_BRANCH"
```

### Step F6: Summary

```bash
# Switch back to base branch (not `git checkout -` which may return to the wrong branch
# if any intermediate checkout happened during conflict resolution)
git checkout "$BASE_BRANCH"
STASH_INDEX=$(git stash list 2>/dev/null | grep "team-rc-fix: auto-stash" | head -1 | grep -oE 'stash@\{[0-9]+\}')
if [ -n "$STASH_INDEX" ]; then
  git stash pop "$STASH_INDEX"
  echo "✅ Restored stashed changes."
fi
```

**FIX APPLIED** ── *{RC_BRANCH}* ────────────────────
Commit: `{COMMIT_MSG}`
Cherry-picked to: *{RC_BRANCH}*

**Next:**
  1. Re-deploy staging to verify
  2. Continue RC validation
  3. When ready → `/team-rc promote`
────────────────────────────────────────────

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

---

## Hotfix Flow

> Triggered by `/team-rc hotfix`. Fix a production bug directly on main, then backport to develop.
> Use when: production has a critical bug that needs immediate fixing, independent of any RC in progress.
> Requires dual-branch mode (base_branch != production_branch).

### Step H1: Pre-flight

Same worktree guard + config reads as Prepare Step 1 (worktree check, REPO, BASE_BRANCH, PROD_BRANCH, VERSION_PREFIX).

```bash
TEAMWORK_DIR=$(bash ~/.claude/commands/scripts/tw-config.sh detect-dir 2>/dev/null) || {
  echo "No teamwork config found. Run /team init first."
  # STOP
}

REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
if [ -z "$REPO" ]; then
  echo "ERROR: Cannot detect repo. Run: gh auth login"
  # STOP
fi

VERSION_PREFIX=$(bash ~/.claude/commands/scripts/tw-config.sh versions.current "" 2>/dev/null)
if [ -z "$VERSION_PREFIX" ]; then
  echo "ERROR: versions.current not set. Add to $TEAMWORK_DIR/config.yml:"
  echo "  versions:"
  echo "    current: \"V0.1\""
  # STOP
fi

BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "develop" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-develop}"

PROD_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.production_branch "main" 2>/dev/null)
PROD_BRANCH="${PROD_BRANCH:-main}"
```

**Single-branch guard:**

```bash
if [ "$BASE_BRANCH" = "$PROD_BRANCH" ]; then
  echo "ERROR: Hotfix requires dual-branch mode (base_branch ≠ production_branch)."
  echo "Your config has both set to '$BASE_BRANCH'."
  echo "Hotfix creates a branch from production and backports to develop — not applicable for single-branch projects."
  # STOP
fi
```

### Step H2: Derive Hotfix Version

```bash
git fetch origin "$PROD_BRANCH" 2>/dev/null
# Ensure local tags are up-to-date before deriving hotfix version — stale local tags
# may cause version collisions if another hotfix was tagged from a different machine.
git fetch --tags origin 2>/dev/null || true
LATEST_TAG=$(git tag --sort=-v:refname --list "${VERSION_PREFIX}*" 2>/dev/null | head -1)
```

- If no tags exist → "No release tags found. Run `/team-rc` + `/team-rc promote` first." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

Derive patch version: if `LATEST_TAG` is `V0.1.0` → `HOTFIX_VERSION` is `V0.1.0-hotfix.1`. If `V0.1.0-hotfix.1` already exists → `V0.1.0-hotfix.2`, etc.

```bash
HOTFIX_NUM=1
while git tag --list "${LATEST_TAG}-hotfix.${HOTFIX_NUM}" 2>/dev/null | grep -q .; do
  HOTFIX_NUM=$((HOTFIX_NUM + 1))
done
HOTFIX_VERSION="${LATEST_TAG}-hotfix.${HOTFIX_NUM}"
HOTFIX_BRANCH="hotfix/${HOTFIX_VERSION}"
```

### Step H3: Create Hotfix Branch

```bash
CURRENT=$(git branch --show-current)
if [[ "$CURRENT" == hotfix/* ]]; then
  HOTFIX_BRANCH="$CURRENT"
  echo "Already on hotfix branch: $HOTFIX_BRANCH. Resuming..."
else
  git checkout -b "$HOTFIX_BRANCH" "origin/$PROD_BRANCH"
fi
```

If already on hotfix branch → skip branch creation, go directly to Step H4.

**HOTFIX** ── **{HOTFIX_VERSION}** ────────────────────
Branch: *{HOTFIX_BRANCH}* (from *{PROD_BRANCH}*)
Base tag: **{LATEST_TAG}**

You are now on the hotfix branch. Make your fix, then continue.
────────────────────────────────────────────

### Step H4: Fix Phase (interactive)

Guide the user to make their fix. This is interactive — the skill does not auto-fix.

After the user makes changes, detect uncommitted work:

```bash
if [ -n "$(git status --porcelain)" ]; then
  echo "Uncommitted changes detected."
fi
```

Use `AskUserQuestion`:
```
question: "Hotfix changes ready to commit?"
options:
  - label: "Commit and continue"
    description: "Stage all changes and commit with hotfix message"
  - label: "Not yet"
    description: "I need more time — come back with /team-rc hotfix"
```

If "Not yet" → output: "Make your fix, then run `/team-rc hotfix` again from this branch (it will detect the existing hotfix branch and resume)." → **STOP**

If "Commit and continue":

```bash
# Show what will be committed — verify no unrelated files
git status --short
# Stage only tracked modified files (not untracked)
git add -u
git commit -m "hotfix(${HOTFIX_VERSION}): fix description"
```

(Use the user's description if available, or ask for a one-line summary. If untracked files need to be included, add them explicitly by name.)

### Step H5: PR to Main

```bash
git push origin "$HOTFIX_BRANCH"
```

Create PR:

```bash
gh pr create \
  --base "$PROD_BRANCH" \
  --head "$HOTFIX_BRANCH" \
  --title "hotfix: ${HOTFIX_VERSION}" \
  --body "## Hotfix\n\nFix for production issue.\n\nBranch: \`${HOTFIX_BRANCH}\`\nBase tag: \`${LATEST_TAG}\`" \
  --repo "$REPO"

PR_NUMBER=$(gh pr view "$HOTFIX_BRANCH" --json number --jq '.number' --repo "$REPO")
```

Wait for CI + merge (same pattern as Promote):

**Try auto-merge first:**

```bash
gh pr merge "$PR_NUMBER" --squash --auto \
  --subject "hotfix: ${HOTFIX_VERSION}" \
  --repo "$REPO"
```

If `--auto` succeeds, poll:

```bash
bash ~/.claude/commands/scripts/tw-pr.sh wait-merged "$PR_NUMBER" "$REPO"
```

**If `--auto` fails** → fallback:

```bash
bash ~/.claude/commands/scripts/tw-pr.sh watch "$PR_NUMBER" || {
  echo "WARNING: CI checks failed or timed out — attempting merge anyway"
}

gh pr merge "$PR_NUMBER" --squash \
  --subject "hotfix: ${HOTFIX_VERSION}" \
  --repo "$REPO"
```

If merge fails → output error details → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### Step H6: Tag on Main

```bash
git fetch origin "$PROD_BRANCH"
bash ~/.claude/commands/scripts/tw-git.sh tag "$HOTFIX_VERSION" "origin/$PROD_BRANCH"
```

Generate hotfix release notes from the merged PR:

```bash
# Extract PR details for release notes
PR_DETAILS=$(gh pr view "$PR_NUMBER" --json title,author,url --repo "$REPO")
PR_TITLE=$(echo "$PR_DETAILS" | jq -r '.title')
PR_AUTHOR=$(echo "$PR_DETAILS" | jq -r '.author.login')
PR_URL=$(echo "$PR_DETAILS" | jq -r '.url')
```

Format the release body:

```
## Hotfix

* {PR_TITLE} by @{PR_AUTHOR} in {PR_URL}

**Full Changelog**: {LATEST_TAG}...{HOTFIX_VERSION}
```

{^ Use LATEST_TAG (from H2), not PREV_TAG. LATEST_TAG is the release being hotfixed.}

```bash
gh release create "$HOTFIX_VERSION" \
  --notes-file <(echo "$HOTFIX_BODY") \
  --title "$HOTFIX_VERSION" \
  --repo "$REPO"
```

Capture and output the release URL. If fails → warn but continue.

### Step H7: Backport to Develop

After PR merged to main, cherry-pick to develop:

```bash
MERGE_SHA=$(gh pr view "$PR_NUMBER" --json mergeCommit --jq '.mergeCommit.oid' --repo "$REPO")

# Guard: mergeCommit can be null if PR was closed without merging or API delay
if [ -z "$MERGE_SHA" ] || [ "$MERGE_SHA" = "null" ]; then
  echo "ERROR: Could not retrieve merge commit SHA for PR #$PR_NUMBER."
  echo "The PR may not have been merged yet, or GitHub API returned stale data."
  echo "Retry in a few seconds, or check: gh pr view $PR_NUMBER --repo $REPO"
  # STOP
fi

git checkout "$BASE_BRANCH"
git pull origin "$BASE_BRANCH"
git cherry-pick "$MERGE_SHA" --no-edit
```

**If cherry-pick conflicts:**

Do NOT stop. Guide the user:

**BACKPORT CONFLICT** ────────────────────
Conflicted files:
```bash
git diff --name-only --diff-filter=U
```

**To resolve:**
  1. Edit the conflicted files (resolve markers: `<<<<<<<`, `=======`, `>>>>>>>`)
  2. Stage resolved files: `git add <file>`
  3. Continue: `git cherry-pick --continue`
  4. Or abort: `git cherry-pick --abort` (hotfix is already on *main* — backport can be done later)
────────────────────────────────────────────

After resolution (or if no conflict):

```bash
git push origin "$BASE_BRANCH"
```

### Step H8: Cleanup + Summary

```bash
# Delete remote hotfix branch
git push origin --delete "$HOTFIX_BRANCH" 2>/dev/null || true

# Delete local hotfix branch
git branch -d "$HOTFIX_BRANCH" 2>/dev/null || git branch -D "$HOTFIX_BRANCH" 2>/dev/null || true
```

**HOTFIX COMPLETE ── {HOTFIX_VERSION}** ────────────────────
Tag: **{HOTFIX_VERSION}** (on *{PROD_BRANCH}*)
PR:  **#{PR_NUMBER}** (squash merged to *{PROD_BRANCH}*)
Release: `{release_url}`
Backport: cherry-picked to *{BASE_BRANCH}*

{Check if an active RC branch exists (rc/* on remote). If so, add this warning:}
⚠️ Active RC detected: *{RC_BRANCH}*. The RC is now behind *{BASE_BRANCH}*. Consider `/team-rc fix` to include the hotfix, or `/team-rc cancel` + `/team-rc` to re-cut.

────────────────────────────────────────────

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

---

## Promote Flow

### Step P1: Pre-flight

Same worktree guard + config reads as Prepare Step 1 (worktree check, REPO, VERSION_PREFIX, BASE_BRANCH, PROD_BRANCH, STAGING_WORKFLOW).

### Step P2: Find Active RC Branch

```bash
RC_BRANCH=$(bash ~/.claude/commands/scripts/tw-git.sh find-rc 2>/dev/null)
RC_EXIT=$?
# RC_EXIT: 0=found one (RC_BRANCH has name), 1=none found, 3=multiple found
```

**If `RC_EXIT` is 1 (none found)** → output and **STOP**:

No active RC branch found. Run `/team-rc` to prepare one first.
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

**If `RC_EXIT` is 3 (multiple found)** → output and **STOP**:

**ERROR:** Multiple RC branches found. Only one RC should exist at a time.
Delete the stale one(s) and retry.
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

Extract version from branch name:

```bash
VERSION=$(echo "$RC_BRANCH" | sed 's|rc/||')
```

### Step P3: Staging Check (optional)

**If `STAGING_WORKFLOW` is set:**

```bash
LATEST_RUN=$(gh run list --workflow "$STAGING_WORKFLOW" --branch "$RC_BRANCH" --limit 1 \
  --json conclusion,displayTitle,createdAt 2>/dev/null)
```

- If conclusion is `success` → `Staging: passed`
- If not → warn with details.

Use `AskUserQuestion`:
```
question: "Staging 未通过。是否继续 promote？"
options:
  - label: "继续"
    description: "Staging 问题已知，继续 promote"
  - label: "中止"
    description: "先修 staging，之后重新 /team-rc promote"
```
If "中止" → **STOP**. Deploy to staging first. Check your CI/CD pipeline or manually trigger deployment.
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

**If not set:**

⚠️ **Staging 未配置** — `deploy.staging_workflow` 为空。
跳过自动 staging 检查。**请手动验证 RC 分支可正常运行后再 promote。**
配置方法：在 `config.yml` 添加 `deploy.staging_workflow: "deploy-staging.yml"`

### Step P4: Confirmation

**PROMOTE ── {VERSION}** ────────────────────
*{RC_BRANCH}* → *{PROD_BRANCH}*

**ACTIONS**
  1. PR: *{RC_BRANCH}* → *{PROD_BRANCH}* (squash)
  2. Tag **{VERSION}** on *{PROD_BRANCH}*
  3. GitHub Release
  4. Check milestone
  5. Delete *{RC_BRANCH}*
────────────────────────────────────────────

Use `AskUserQuestion`:
```
question: "确认 promote？这会创建 PR、tag、GitHub Release。"
options:
  - label: "Promote"
    description: "执行上述所有操作"
  - label: "中止"
    description: "取消 promote"
```
If "中止" → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### Step P4.5: Rebuild RC on Production (squash-only mode)

**Why this exists**: In squash-only release repos (e.g., `pre-launch` → `main` where every release is a squash merge), `main` and the base branch re-diverge after every release. The merge-base between `rc/VERSION` and `main` stays rooted at the original fork point, so the release PR conflicts with every direct-to-main change in history. This step rebuilds `rc/VERSION` on top of `main` HEAD via `git commit-tree` so the PR becomes a clean fast-forward.

**When this runs**: Only when BOTH:
1. `BASE_BRANCH != PROD_BRANCH` (dual-branch mode), AND
2. `config.yml` has `release.strategy: rebuild-on-prod` (explicit opt-in; default `cut-only` preserves legacy behavior for single-branch or develop→main repos)

```bash
REBUILD_STRATEGY=$(bash ~/.claude/commands/scripts/tw-config.sh release.strategy "cut-only")

if [ "$BASE_BRANCH" = "$PROD_BRANCH" ] || [ "$REBUILD_STRATEGY" != "rebuild-on-prod" ]; then
  echo "Rebuild step skipped (strategy=$REBUILD_STRATEGY, base=$BASE_BRANCH, prod=$PROD_BRANCH)"
  # Jump to Step P5
else
  # Dry-run first to preview main-only files that would be deleted
  DRY_OUTPUT=$(bash ~/.claude/commands/scripts/tw-git.sh rebuild-rc-on-prod \
    "$RC_BRANCH" "$PROD_BRANCH" "$VERSION" --dry-run)
fi
```

**Parse `DRY_OUTPUT`** (AI step, KV format):

- `OLD_SHA=<sha>` — current `origin/$RC_BRANCH` sha
- `NEW_SHA=<sha>` — sha of the rebuilt commit
- `MAIN_ONLY_FILES<<EOF ... EOF` — heredoc-framed list of files present on `$PROD_BRANCH` but not on `$RC_BRANCH` (empty line `MAIN_ONLY_FILES=` if none)
- `STATUS=<dry-run|already-fast-forwardable>`

**If `STATUS=already-fast-forwardable`**: No rebuild needed — skip to Step P5.

**Display preview**:

**REBUILD RC** ── *{RC_BRANCH}* → *{PROD_BRANCH}* ────────────
  Old rc SHA: `{OLD_SHA[0:12]}`
  New rc SHA: `{NEW_SHA[0:12]}`
  Main-only files to be deleted: `{count}`
{if MAIN_ONLY_FILES non-empty, list each on its own indented line wrapped in backticks}
────────────────────────────────────────────────────────────

**If `MAIN_ONLY_FILES` is non-empty**, use `AskUserQuestion`:
```
question: "确认删除上述 main-only 文件并 rebuild？"
options:
  - label: "执行 rebuild"
    description: "提交新 commit 到 rc 分支，这些 main-only 文件将从 rc 消失"
  - label: "先 cherry-pick 再来"
    description: "中止 promote，手工把这些文件 backport 到 rc 后重新 /team-rc promote"
  - label: "中止"
    description: "取消本次 promote"
```

If "先 cherry-pick 再来" or "中止" → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

**If `MAIN_ONLY_FILES` is empty OR user confirmed**, execute the real rebuild:

```bash
# Rollback hint — capture OLD_SHA so the user can restore if something goes wrong
echo "  [rollback] If needed: git push origin $OLD_SHA:$RC_BRANCH --force-with-lease"

REAL_OUTPUT=$(bash ~/.claude/commands/scripts/tw-git.sh rebuild-rc-on-prod \
  "$RC_BRANCH" "$PROD_BRANCH" "$VERSION")

# Verify STATUS=pushed (real run uses --force-with-lease internally)
if ! echo "$REAL_OUTPUT" | grep -q "^STATUS=pushed$"; then
  echo "ERROR: Rebuild did not complete"
  # STOP
fi

# Parse NEW_SHA from REAL_OUTPUT (may differ from dry-run if prod advanced between runs)
NEW_SHA=$(echo "$REAL_OUTPUT" | grep '^NEW_SHA=' | sed 's/^NEW_SHA=//')
```

Output: *{RC_BRANCH}* rebuilt on *{PROD_BRANCH}* at `{NEW_SHA[0:12]}` (rollback SHA: `{OLD_SHA[0:12]}`)

If rebuild fails → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

{^ The real run re-fetches and recomputes OLD_SHA internally. If someone else pushed to RC between dry-run and this step, --force-with-lease aborts safely — no data loss.}

### Step P5: Generate Release Notes

**Why not `--generate-notes`**: GitHub's auto-generate relies on PRs visible from `main`. After squash merge, `main` sees only the release PR itself — all original PR authors and links are lost. We reconstruct from the base branch's merged PR history via the shared `tw-pr.sh merged-since` helper(deterministic filter + categorize, same mechanism reused by `/team-rc` Prepare's 飞书 block).

```bash
# Find previous clean version tag as time boundary (current version is NOT yet tagged — that happens at P7).
# Filter to only clean version tags (PREFIX.N) — exclude hotfix tags (PREFIX.N-hotfix.M).
PREV_TAG=$(git tag --sort=-v:refname --list "${VERSION_PREFIX}*" \
  | grep -E "^$(echo "$VERSION_PREFIX" | sed 's/[.[\*^$()+?{|\\]/\\&/g')\.[0-9]+$" \
  | head -1)

# Generate release notes via shared helper. It handles UTC Z timezone normalization
# (3.13.0 fix — %aI local-tz bug silently dropped PRs merged in the offset window),
# gh pr list, jq filter (mergedAt > prev_tag_date), and category by conventional-commit prefix.
RELEASE_BODY=$(bash ~/.claude/commands/scripts/tw-pr.sh merged-since \
  "$BASE_BRANCH" "$PREV_TAG" \
  --repo "$REPO" \
  --format markdown \
  --version "$VERSION")
```

{^ Output shape: `## What's in this release` → `### Features` / `### Bug Fixes` / `### Other Changes` → `### Contributors` → `**Full Changelog**: {PREV_TAG}...{VERSION}` (or "Initial release" if PREV_TAG empty). Empty categories auto-omitted. This matches the format shipped pre-3.15.0 — refactor is pure deduplication with Step 7 Prepare, no format change.}

{^ Categorization is by PR title conventional-commit prefix: `feat(scope)?!?:` → Features, `fix(scope)?!?:` → Bug Fixes, everything else → Other Changes.}

Store the resulting markdown in `$RELEASE_BODY` for use in P6 and P8. `$PREV_TAG` is reused in P11 for the 飞书 announcement block.

### Step P6: Create PR and Squash Merge

```bash
# Create PR
gh pr create \
  --base "$PROD_BRANCH" \
  --head "$RC_BRANCH" \
  --title "release: $VERSION" \
  --body "$RELEASE_BODY" \
  --repo "$REPO"

# Get PR number
PR_NUMBER=$(gh pr view "$RC_BRANCH" --json number --jq '.number' --repo "$REPO")
```

Squash merge via GitHub (respects branch protection):

**Try auto-merge first** (waits for CI automatically):

```bash
gh pr merge "$PR_NUMBER" --squash --auto \
  --subject "release: $VERSION" \
  --body "$RELEASE_BODY" \
  --repo "$REPO"
```

If `--auto` succeeds, poll until merge completes:

```bash
bash ~/.claude/commands/scripts/tw-pr.sh wait-merged "$PR_NUMBER" "$REPO"
```

If wait-merged exits non-zero → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

**If `--auto` fails** (repo doesn't have auto-merge enabled) → fallback:

```bash
# Wait for CI checks to pass
bash ~/.claude/commands/scripts/tw-pr.sh watch "$PR_NUMBER" || {
  echo "WARNING: CI checks failed or timed out — attempting merge anyway"
}

# Merge directly
gh pr merge "$PR_NUMBER" --squash \
  --subject "release: $VERSION" \
  --body "$RELEASE_BODY" \
  --repo "$REPO"
```

Output: PR **#{PR_NUMBER}** squash merged to *{PROD_BRANCH}*

If merge fails → output error details → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### Step P7: Tag on Main

```bash
# Get the exact merge commit SHA from the PR — immune to fetch timing issues.
# Using origin/$PROD_BRANCH is fragile: if fetch hasn't picked up the merge yet,
# the tag lands on the wrong commit (caused V0.1.6 incident).
MERGE_SHA=$(gh pr view "$PR_NUMBER" --json mergeCommit --jq '.mergeCommit.oid' --repo "$REPO")

if [ -z "$MERGE_SHA" ] || [ "$MERGE_SHA" = "null" ]; then
  echo "ERROR: Could not get merge commit SHA from PR #$PR_NUMBER" >&2
  echo "Falling back to origin/$PROD_BRANCH"
  git fetch origin "$PROD_BRANCH"
  MERGE_SHA="origin/$PROD_BRANCH"
fi

# Idempotency: skip if tag already exists (re-entry after partial failure)
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "Tag $VERSION already exists — skipping"
else
  bash ~/.claude/commands/scripts/tw-git.sh tag "$VERSION" "$MERGE_SHA"
fi
```

**Does NOT checkout or modify the working tree.** Tags the PR merge commit SHA directly.

Output: Tag **{VERSION}** pushed (or "already exists" on re-entry)

### Step P8: GitHub Release

```bash
# Idempotency: skip if release already exists (re-entry after partial failure)
if gh release view "$VERSION" --repo "$REPO" >/dev/null 2>&1; then
  echo "Release $VERSION already exists — skipping"
else
  gh release create "$VERSION" \
    --notes-file <(echo "$RELEASE_BODY") \
    --title "$VERSION" \
    --repo "$REPO"
fi

# ALWAYS capture RELEASE_URL after either branch — P11 summary + 飞书 block both rely on it.
# Pre-3.15.0 this was only set in the idempotent branch; fresh-create releases lost the URL.
RELEASE_URL=$(gh release view "$VERSION" --repo "$REPO" --json url --jq '.url' 2>/dev/null || echo "")
```
{^ --notes-file avoids shell escaping issues when PR titles contain quotes or special characters}

Capture and output the release URL. Required by P11 summary + the 飞书 announcement block.

If fails → warn but continue (tag exists, release can be created manually via `gh release create`). If `RELEASE_URL` ends up empty, P11 will fall back to a placeholder and flag the gap for manual fill-in.

### Step P8b: Notification

```bash
bash ~/.claude/commands/scripts/tw-notify.sh milestone.done \
  --title "$VERSION"
```

Non-fatal: if notification fails, warn but continue.

### Step P9: Milestone Check

```bash
MILESTONE_FULL=$(bash ~/.claude/commands/scripts/tw-git.sh milestone-resolve "$VERSION_PREFIX" 2>/dev/null)
```

Then query the resolved milestone by exact title:

```bash
MILESTONE_DATA=$(gh api "repos/$REPO/milestones" 2>/dev/null \
  | jq -r --arg title "$MILESTONE_FULL" '.[] | select(.title == $title)')
```

**If milestone found:**

```bash
OPEN_COUNT=$(echo "$MILESTONE_DATA" | jq -r '.open_issues')
CLOSED_COUNT=$(echo "$MILESTONE_DATA" | jq -r '.closed_issues')
MILESTONE_NUMBER=$(echo "$MILESTONE_DATA" | jq -r '.number')
```

- If `OPEN_COUNT == 0` → all issues closed → close milestone:

  ```bash
  gh api "repos/$REPO/milestones/$MILESTONE_NUMBER" \
    --method PATCH --field state=closed
  ```

  Output: Milestone **{MILESTONE_FULL}** closed (all issues complete)

- If `OPEN_COUNT > 0` → output: Milestone **{MILESTONE_FULL}**: `{CLOSED_COUNT}` closed, `{OPEN_COUNT}` still open

**If no milestone:** *(No milestone "{MILESTONE_FULL}" found — skipped)*

### Step P10: Delete RC Branch

RC is a frozen snapshot — all fixes live in *{BASE_BRANCH}*. No merge-back needed.

```bash
# Delete remote
git push origin --delete "$RC_BRANCH" 2>/dev/null || {
  echo "WARNING: Could not delete remote RC branch" >&2
}

# Delete local
git branch -d "$RC_BRANCH" 2>/dev/null || git branch -D "$RC_BRANCH" 2>/dev/null || true
```

### Step P10b: Production Deploy

**If `$PROD_WORKFLOW` is set:**

```bash
gh workflow run "$PROD_WORKFLOW" --ref "$VERSION" -f environment=production -f version="$VERSION"
```

- If workflow triggers → 🚀 Production deploy triggered for **{VERSION}**
- If workflow fails to trigger → ⚠️ Could not trigger `{PROD_WORKFLOW}`. Deploy manually: `gh workflow run {PROD_WORKFLOW} --ref {PROD_BRANCH}`

**If `$PROD_WORKFLOW` is not set:**

🚀 **Deploy to production** — code is on *{PROD_BRANCH}*, tagged **{VERSION}**.
  Trigger your deploy pipeline for **{VERSION}**, or configure auto-deploy:
  `deploy.production_workflow: "deploy.yml"` in `$TEAMWORK_DIR/config.yml`

### Step P11: Summary

**PROMOTED ── {VERSION}** ────────────────────
Tag: **{VERSION}** (on *{PROD_BRANCH}*)
PR:  **#{PR_NUMBER}** (squash merged)
Release: {$RELEASE_URL or ⚠️ "not captured — run `gh release view $VERSION` manually"}
RC Branch:  *{RC_BRANCH}* (deleted)
Milestone:  **{MILESTONE_FULL}** — {status}
Deploy:  {🚀 triggered `{PROD_WORKFLOW}` / 🚀 deploy manually}

────────────────────────────────────────────

**Generate 飞书 announcement**(reuses changelog from P5 + release metadata from P6-P8):

```bash
# PREV_TAG already derived in P5; reuse it
FEISHU_BLOCK=$(bash ~/.claude/commands/scripts/tw-pr.sh merged-since \
  "$BASE_BRANCH" "$PREV_TAG" \
  --repo "$REPO" \
  --format feishu \
  --version "$VERSION" \
  --env production \
  --release-url "$RELEASE_URL" \
  --pr-number "$PR_NUMBER" \
  --tag-sha "$MERGE_SHA")
```

Output between dividers for copy-paste:

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📣 **飞书发布通知 — 复制下方内容**
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

{$FEISHU_BLOCK content}

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

{^ If `RELEASE_URL` is empty(P8 capture failed), feishu block still renders but omits the GitHub Release line — warn the user.}
{^ If `merged-since` fails → fall back to printing `$RELEASE_BODY` from P5 bounded by the same dividers. The old behavior is the worst case, not a regression.}

`/team-rc` to prepare next RC

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

---

## Operation Promote Help

> Triggered by `/team-rc promote help`.

Output the following and **STOP**:

**`/team-rc promote` — Promote RC to Production**
────────────────────────────────────────────

**USAGE**
  `/team-rc promote`          Promote current RC to production

**WHAT IT DOES**
  1. Verify RC branch exists and CI is green
  2. Verify staging deployment succeeded (if configured)
  3. Squash merge RC → *main* (production branch)
  4. Create git tag with version
  5. Create GitHub Release with auto-generated notes
  6. Clean up RC branch
  7. Trigger production deploy (if `deploy.production_workflow` configured), or guide manual deploy

**WHEN TO USE**
  After RC has been tested on staging and approved.
  Flow: `/team-rc` → *(test staging)* → `/team-rc promote`

**PREREQUISITES**
  RC branch must exist (created by `/team-rc`)
  CI must be green on RC branch
  Staging deploy must be successful (if `deploy.staging_workflow` configured)

**TROUBLESHOOTING**
  "No RC branch found" → Run `/team-rc` first to cut an RC
  "CI not green" → Fix on *{BASE_BRANCH}*, then `/team-rc cancel` + `/team-rc` to re-cut
  "Staging not verified" → Deploy to staging first, verify manually

**BUGFIX DURING RC**
  RC is a frozen snapshot — don't push to it directly.
  Small fix? Cherry-pick via `/team-rc fix` (keeps test progress).
  Big fix? Fix on *$BASE_BRANCH* → `/team-rc cancel` → `/team-rc` (re-cut).
  See also: `/team-rc cancel`, `/team-rc fix`

---

## Operation Fix Help

> Triggered by `/team-rc fix help`.

Output the following and **STOP**:

**`/team-rc fix` — Cherry-pick Fix to RC**
────────────────────────────────────────────

**USAGE**
  `/team-rc fix`                Show recent commits on *{BASE_BRANCH}*, select one to cherry-pick
  `/team-rc fix <commit>`       Cherry-pick specific commit hash to active RC
  `/team-rc fix <branch>`       Cherry-pick HEAD of branch to active RC

**WHAT IT DOES**
  1. Find active RC branch
  2. Validate the commit exists and is not already on RC
  3. Cherry-pick the commit onto RC branch
  4. Push updated RC

**WHEN TO USE**
  A bug was fixed on *{BASE_BRANCH}* and you need it on the RC without re-cutting.
  This preserves existing RC test progress (unlike cancel + re-cut).

  `/team-rc cancel` + `/team-rc` = **re-cut** (clean slate, lose test progress)
  `/team-rc fix` = **cherry-pick** (surgical, keep test progress)

**PREREQUISITES**
  RC branch must exist (created by `/team-rc`)
  Fix commit must exist on *{BASE_BRANCH}*

**CONFLICT RESOLUTION**
  If the cherry-pick has conflicts, you'll be guided through resolution:
  1. Edit conflicted files (resolve markers)
  2. `git add <file>`
  3. `git cherry-pick --continue`

**AFTER FIX**
  Re-deploy staging to verify the fix on RC.

---

## Operation Hotfix Help

> Triggered by `/team-rc hotfix help`.

Output the following and **STOP**:

**`/team-rc hotfix` — Post-Release Production Fix**
────────────────────────────────────────────

**USAGE**
  `/team-rc hotfix`             Start hotfix from *{PROD_BRANCH}*

**WHAT IT DOES**
  1. Create *hotfix/{version}* branch from *{PROD_BRANCH}*
  2. Guide you through making the fix (interactive)
  3. PR to *{PROD_BRANCH}* (squash merge)
  4. Tag hotfix version on *{PROD_BRANCH}*
  5. Create GitHub Release
  6. Backport (cherry-pick) to *{BASE_BRANCH}*
  7. Clean up hotfix branch

**WHEN TO USE**
  Production has a critical bug that needs immediate fixing.
  The fix goes to *{PROD_BRANCH}* first, then backports to *{BASE_BRANCH}*.

  `/team-rc fix` = cherry-pick **to RC** (during staging verification)
  `/team-rc hotfix` = fix **production** (after release, backport to develop)

**PREREQUISITES**
  Dual-branch mode required (`base_branch` != `production_branch`)
  At least one release tag must exist (from a prior `/team-rc promote`)

**BACKPORT CONFLICTS**
  If the cherry-pick to *{BASE_BRANCH}* has conflicts, you'll be guided through resolution.
  The hotfix is already on *{PROD_BRANCH}* — backport can be retried manually if needed.
