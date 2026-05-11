---
description: "RC lifecycle — prepare staging or promote to production. Try: /team-rc help"
version: "3.19.1"
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
| `beta` | Beta: promote into the beta env (staging-prod intermediate). Source flag controls what's promoted |
| `beta help` | Show beta operation usage |
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
  `/team-rc beta`         Beta: promote into beta env (intermediate between staging + prod)
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
  3. Confirm promote (P4)
  4. **Hotfix divergence guard** (P4.4, dual-branch + `rebuild-on-prod` only) — detect prod hotfixes missing from RC; STOP if user picks Reconcile/Abort, see §Hotfix Divergence Guard
  5. **Rebuild RC on prod** (P4.5, dual-branch + `rebuild-on-prod` only) — `git commit-tree` so PR fast-forwards; confirm before deletes/overwrites
  6. Create PR: rc → *main* (squash merge via GitHub)
  7. Tag squash commit on *main*
  8. Create GitHub Release
  9. Close milestone if complete
  10. Delete rc branch

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

**BETA** (`/team-rc beta`) *(requires `release.beta_workflow` + `deploy.beta_*` config)*
  Promote an image set to the **beta environment** — an intermediate between staging and prod for teams whose staging is too churned-out by daily dev work to function as a real preprod. Operator-driven (NEVER auto-merge gitops PR).
  Sources: `--source=prod` (parity from prod), `--source=staging` (latest staging), `--source=explicit --version-tag=Vx` (specific tag).
  `/team-rc beta`                                    Interactive
  `/team-rc beta --source=staging`                   Take latest staging
  `/team-rc beta --source=prod --reason="..."`       Snapshot prod (reason mandatory)
  `/team-rc beta --source=explicit --version-tag=V0.1.32 --reason="..."`
  `/team-rc beta --dry-run`                          Echo commands; no actions

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
- If `beta help` → jump to **Operation Beta Help**
- If `beta` (with or without trailing flags) → jump to **Beta Flow**
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

### Step H7.5: RC-rebuild Sub-chain (config-gated, 3.18.0)

**Why this exists**: Some release pipelines (e.g., vi_agent's `build.yml on rc/**` + `promote-prod.yml` staging-tag-copy semantics) require hotfix tags to go through an RC build cycle before prod deploy actually picks up the fix. Without this sub-chain, the hotfix tag pushes correctly but the production-deploy workflow no-ops because the staging artifact tag hasn't been refreshed. Real incident: vi_agent V0.1.27-hotfix.1 burned 12 min of CI on a no-op promote before the team realized the pattern.

**When this runs**: Only if `release.hotfix_needs_rc_rebuild: true` in `config.yml`. Default OFF — preserves current single-step "backport + done" behavior for repos whose CI builds artifacts directly from tag push.

```bash
HOTFIX_NEEDS_RC_REBUILD=$(bash ~/.claude/commands/scripts/tw-config.sh release.hotfix_needs_rc_rebuild "false")

if [ "$HOTFIX_NEEDS_RC_REBUILD" != "true" ]; then
  echo "RC-rebuild sub-chain skipped (release.hotfix_needs_rc_rebuild=$HOTFIX_NEEDS_RC_REBUILD)"
  # Jump to Step H8
else
  BUILD_WORKFLOW=$(bash ~/.claude/commands/scripts/tw-config.sh release.build_workflow "" 2>/dev/null)
  if [ -z "$BUILD_WORKFLOW" ]; then
    echo "ERROR: release.hotfix_needs_rc_rebuild=true requires release.build_workflow to be set."
    echo "Add to $TEAMWORK_DIR/config.yml:"
    echo "  release:"
    echo "    hotfix_needs_rc_rebuild: true"
    echo "    build_workflow: \"build.yml\"  # workflow that builds artifacts on rc/** push"
    # STOP
  fi
fi
```

**Step H7.5a: Cut rc/{hotfix-version} from main HEAD**

```bash
RC_HOTFIX_BRANCH="rc/${HOTFIX_VERSION}"

# Use main HEAD which is the squash commit of the hotfix PR.
git fetch origin "$PROD_BRANCH" 2>/dev/null
git push origin "origin/$PROD_BRANCH:refs/heads/$RC_HOTFIX_BRANCH" || {
  echo "ERROR: Failed to cut $RC_HOTFIX_BRANCH from origin/$PROD_BRANCH"
  # STOP
}
echo "  ✓ Cut $RC_HOTFIX_BRANCH from $PROD_BRANCH HEAD"
```

**Step H7.5b: Dispatch build workflow** (avoids `on: push: branches:` path-filter races for first-push)

```bash
# workflow_dispatch is reliable for new branches; on:push path filters often
# don't trigger consistently for newly-created refs.
gh workflow run "$BUILD_WORKFLOW" --ref "$RC_HOTFIX_BRANCH" || {
  echo "ERROR: Failed to dispatch $BUILD_WORKFLOW on $RC_HOTFIX_BRANCH"
  echo "Manual fallback: gh workflow run $BUILD_WORKFLOW --ref $RC_HOTFIX_BRANCH"
  # CLEANUP BEFORE STOP (3.19.0 fix): H7.5a succeeded so temp branch exists on origin.
  git push origin --delete "$RC_HOTFIX_BRANCH" >/dev/null 2>&1 || true
  echo "  [cleanup] deleted temp $RC_HOTFIX_BRANCH before exit"
  # STOP
}

# Capture the run ID for watching.
sleep 3  # gh API takes a moment to register the dispatch
BUILD_RUN_ID=$(gh run list --workflow "$BUILD_WORKFLOW" --branch "$RC_HOTFIX_BRANCH" \
  --limit 1 --json databaseId --jq '.[0].databaseId' 2>/dev/null)
echo "  ✓ Dispatched $BUILD_WORKFLOW on $RC_HOTFIX_BRANCH (run #$BUILD_RUN_ID)"
```

**Step H7.5c: Wait for build to complete**

```bash
echo "  Waiting for $BUILD_WORKFLOW (run #$BUILD_RUN_ID) — typically 6-15min..."
gh run watch "$BUILD_RUN_ID" --exit-status || {
  echo "ERROR: Build workflow failed. Inspect: gh run view $BUILD_RUN_ID"
  # CLEANUP BEFORE STOP (3.19.0 fix): temp branch exists on origin; build was dispatched but failed.
  git push origin --delete "$RC_HOTFIX_BRANCH" >/dev/null 2>&1 || true
  echo "  [cleanup] deleted temp $RC_HOTFIX_BRANCH before exit"
  # STOP
}
echo "  ✓ Build complete"
```

**Step H7.5d: Re-push tag (re-trigger production-deploy with refreshed staging artifact)**

The production-deploy workflow is bound to tag-push events. Pushing the same-named tag a second time is the only way to re-trigger it. But that orphans the GitHub Release. Use `repush-tag-and-release` to handle both atomically:

```bash
# Build release notes file for recreation. Same body as H6 (link the hotfix PR).
NOTES_FILE=$(mktemp)
cat > "$NOTES_FILE" <<EOF
Hotfix release ${HOTFIX_VERSION}.

PR: #${PR_NUMBER}
Backport: ${BASE_BRANCH} (cherry-picked)
EOF

REPUSH_OUTPUT=$(bash ~/.claude/commands/scripts/tw-git.sh repush-tag-and-release \
  "$HOTFIX_VERSION" "$MERGE_SHA" "$HOTFIX_VERSION" "$NOTES_FILE")

rm -f "$NOTES_FILE"

if ! echo "$REPUSH_OUTPUT" | grep -q "^STATUS=ok$"; then
  echo "ERROR: repush-tag-and-release failed"
  echo "$REPUSH_OUTPUT"
  # CLEANUP BEFORE STOP (3.19.0 fix): the temp branch from H7.5a still exists on origin.
  # Delete it now so retry of /team-rc hotfix doesn't collide on the same name.
  git push origin --delete "$RC_HOTFIX_BRANCH" >/dev/null 2>&1 || true
  echo "  [cleanup] deleted temp $RC_HOTFIX_BRANCH before exit"
  # STOP
fi

# Capture new release URL.
NEW_RELEASE_URL=$(echo "$REPUSH_OUTPUT" | grep '^RELEASE_URL=' | sed 's/^RELEASE_URL=//')
echo "  ✓ Tag $HOTFIX_VERSION re-pushed; release recreated at $NEW_RELEASE_URL"
```

**Step H7.5e: Cleanup the temporary RC branch**

```bash
# Production-deploy workflow has now been re-triggered by tag re-push. The
# rc/$HOTFIX_VERSION branch was a one-shot dispatch target — remove it.
git push origin --delete "$RC_HOTFIX_BRANCH" 2>/dev/null || true
echo "  ✓ Cleaned up $RC_HOTFIX_BRANCH"
```

> **Why this is config-gated**: Most CI pipelines (build directly on tag, no staging-copy) don't need this. Adding it unconditionally would make hotfix flow ~10 minutes slower for no benefit. vi_agent and similar architectures opt-in by setting `release.hotfix_needs_rc_rebuild: true` + `release.build_workflow: "build.yml"` in their config.

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

## Beta Flow

> Triggered by `/team-rc beta` (or `beta --source=...`, `beta --dry-run`). Promotes an image set into the **beta environment** — an intermediate gate between staging and prod for teams whose staging branch sees too much dev churn to act as a real preprod. Validated end-to-end on vi_agent 2026-05-11 (F4 dossier).
>
> **Origin**: F5 brief from vi_agent Phase 2 beta-prod-parity-audit (Drive Mode session 2026-05-11). See `/Users/liyasong/Project/collov/vi_agent/drive-projects/v2/beta-prod-parity-audit@2.3.0/grading/evidence/f5-brief-for-erwin.md` for the canonical contract this flow implements.

### ⛔ Hard Rules (apply throughout Beta Flow)

| # | Rule | Why |
|---|---|---|
| **H1** | NEVER `gh pr merge` the gitops PR yourself | Operator manual gate by design (B5) |
| **H2** | NEVER add `auto-merge` label to gitops PR | Same as H1 |
| **H3** | `--dry-run` = echo-only, NO actions take effect | Operators rely on dry-run to plan |
| **H4** | NEVER use local `argocd app sync` CLI — always `kubectl annotate` + `kubectl patch` | Port-forward + gRPC broken in many corp networks |
| **H5** | NEVER `kubectl run smoke ...` naively against beta gateway | Default-deny NetworkPolicy + label-allowlist rejects random pods |
| **H6** | `--source=prod` REQUIRES `--reason=...` (non-empty) | Workflow's Validate Inputs step rejects empty |
| **H7** | Smoke MUST extract `.service_version` from body | HTTP 200 hides cosmetic bugs (V0.1.28 conflict-marker precedent) |
| **H8** | On `phase=Succeeded` + `syncResult.resources=[]` + still OutOfSync → re-run the SAME patch (don't escalate first occurrence) | "Phantom-Succeeded variant B" — observed 2× (2026-04-20 staging, 2026-05-11 beta) |

### Step B1: Pre-flight + argument parsing

```bash
# Worktree guard (same as Prepare Step 1)
if [ -f .mission ]; then
  MAIN_REPO=$(bash ~/.claude/commands/scripts/tw-git.sh worktree-main-repo)
  echo "⚠️ /team-rc beta must run from the main repo, not a mission worktree."
  echo "Switch to main repo: cd $MAIN_REPO"
  # STOP
fi

# Config reads — these are HARD requirements; missing → STOP with actionable error.
REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
BETA_WORKFLOW=$(bash ~/.claude/commands/scripts/tw-config.sh release.beta_workflow "" 2>/dev/null)
BETA_GITOPS_REPO=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.beta_gitops_repo "" 2>/dev/null)
[ -z "$BETA_GITOPS_REPO" ] && BETA_GITOPS_REPO=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.gitops_repo "" 2>/dev/null)
BETA_APPS=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.beta_apps "" 2>/dev/null)
BETA_KUBE_CTX=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.beta_kube_context "" 2>/dev/null)
BETA_SMOKE_URL=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.beta_smoke_url "" 2>/dev/null)
PROD_SMOKE_URL=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.prod_smoke_url "" 2>/dev/null)
BETA_SMOKE_FALLBACK_POD=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.beta_smoke_fallback_pod "" 2>/dev/null)
BETA_SMOKE_FALLBACK_URL=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.beta_smoke_fallback_url "" 2>/dev/null)

# Validate required config — abort with discoverable error message if missing.
MISSING=()
[ -z "$BETA_WORKFLOW" ] && MISSING+=("release.beta_workflow")
[ -z "$BETA_APPS" ] && MISSING+=("deploy.beta_apps")
[ -z "$BETA_KUBE_CTX" ] && MISSING+=("deploy.beta_kube_context")
[ -z "$BETA_SMOKE_URL" ] && MISSING+=("deploy.beta_smoke_url")

if [ ${#MISSING[@]} -gt 0 ]; then
  echo "⛔ /team-rc beta requires beta config in $TEAMWORK_DIR/config.yml. Missing:"
  printf '   - %s\n' "${MISSING[@]}"
  echo "See `/team-rc beta help` for full template."
  # STOP
fi
```

**Parse `$ARGUMENTS`** for source + reason + version-tag + dry-run flags:

- `$ARGUMENTS` is the raw arg string after `beta` keyword (e.g., `--source=prod --reason='snapshot'`)
- Extract: `SOURCE`, `REASON`, `VERSION_TAG`, `DRY_RUN` (bool)
- Validation:
  - `SOURCE` must be one of `prod`, `staging`, `explicit`, or empty (= interactive)
  - If `SOURCE=prod` → `REASON` must be non-empty (Hard Rule H6)
  - If `SOURCE=explicit` → `VERSION_TAG` must match pattern `V*` and `REASON` non-empty
  - If `DRY_RUN=true` → all subsequent steps echo only

**If `SOURCE` is empty (interactive)**: Use `AskUserQuestion`:
```
question: "Which source to promote into beta?"
options:
  - label: "prod (parity snapshot — requires reason)"
  - label: "staging (latest tested)"
  - label: "explicit (specific version tag)"
  - label: "abort"
```
Then a follow-up for `--reason` / `--version-tag` as needed.

**Display planned form to operator** (always, even in dry-run):

**═════ /team-rc beta plan ═════**
Workflow: `{BETA_WORKFLOW}` in `{REPO}`
Source:   *{SOURCE}*
Reason:   *{REASON or "(none)"}*
Version:  *{VERSION_TAG or "(workflow-derived)"}*
Mode:     *{"DRY-RUN — no actions" if dry_run else "LIVE"}*
═══════════════════════════════════

`AskUserQuestion` to confirm before B2. **If dry_run=true: skip confirmation and just echo each B2-B9 step's commands.**

### Step B2: Dispatch promote-beta workflow

**3.19.1 — dry-run gate applies to EVERY B2-B9 code block** (3.19.0 only gated B2 — Hard Rule H3 violation when B3-B9 fired regardless). Each step below uses `if [ "$DRY_RUN" = "true" ]; then echo; else <real>; fi`. NO `gh workflow run`, NO `gh pr merge`, NO `kubectl`, NO `curl` execute in dry-run.

```bash
if [ "$DRY_RUN" = "true" ]; then
  echo "[dry-run] gh workflow run $BETA_WORKFLOW --repo $REPO -f source=$SOURCE -f reason='$REASON' ${VERSION_TAG:+-f version_tag=$VERSION_TAG}"
  # Stub IDs so B3-B9 echoes are coherent (no `$RUN_ID` undefined errors).
  RUN_ID="<dry-run-run-id>"
  RUN_URL="<dry-run-run-url>"
else
  gh workflow run "$BETA_WORKFLOW" --repo "$REPO" \
    -f source="$SOURCE" -f reason="$REASON" \
    ${VERSION_TAG:+-f version_tag="$VERSION_TAG"} || {
    echo "⛔ Failed to dispatch $BETA_WORKFLOW"
    # STOP
  }

  # Capture the freshly-dispatched run. gh API needs ~3s to register the dispatch.
  sleep 3
  RUN_INFO=$(gh run list --workflow="$BETA_WORKFLOW" --repo "$REPO" --limit 1 \
    --json databaseId,url --jq '.[0]')
  RUN_ID=$(echo "$RUN_INFO" | jq -r '.databaseId')
  RUN_URL=$(echo "$RUN_INFO" | jq -r '.url')
  echo "  ✓ Dispatched run #$RUN_ID → $RUN_URL"
fi
```

> **Race-condition note**: `--limit 1` after `sleep 3` is the simplest capture. If multiple operators dispatch the same workflow simultaneously, this can grab the wrong run. For solo or low-cadence use (typical for beta promote) this is acceptable. If you observe wrong-run capture, refresh via `gh run list --limit 5` and let the user pick.

### Step B3: Watch the workflow

```bash
if [ "$DRY_RUN" = "true" ]; then
  echo "[dry-run] gh run watch $RUN_ID --repo $REPO --exit-status"
else
  gh run watch "$RUN_ID" --repo "$REPO" --exit-status || {
    echo "⛔ Workflow run #$RUN_ID failed — $RUN_URL"
    echo "Recent logs:"
    gh run view "$RUN_ID" --repo "$REPO" --log-failed | tail -50
    # STOP
  }
  echo "  ✓ Workflow run #$RUN_ID succeeded"
fi
```

### Step B4: Surface the gitops promote PR

```bash
if [ "$DRY_RUN" = "true" ]; then
  echo "[dry-run] gh pr list --repo $BETA_GITOPS_REPO --search 'promote/beta in:head' --state open --limit 5"
  echo "[dry-run] would surface gitops PR diff for operator review"
  GITOPS_PR_JSON='[]'
  GITOPS_PR_COUNT=0
else
  GITOPS_PR_JSON=$(gh pr list --repo "$BETA_GITOPS_REPO" \
    --search 'promote/beta in:head' --state open --limit 5 \
    --json number,title,url,headRefName,additions,deletions,changedFiles)
  GITOPS_PR_COUNT=$(echo "$GITOPS_PR_JSON" | jq 'length')

  if [ "$GITOPS_PR_COUNT" -eq 0 ]; then
    echo "⛔ No open promote/beta* PR found in $BETA_GITOPS_REPO."
    echo "Workflow $BETA_WORKFLOW completed but did not open a PR — investigate."
    # STOP
  fi
fi
```

**Display the PR(s)** with diff summary + per-app image-tag bumps:

For each PR in `$GITOPS_PR_JSON`:
- Show URL, title, headRefName, additions/deletions/changedFiles
- Run `gh pr diff <N> --repo "$BETA_GITOPS_REPO"` and group output by overlay path (`apps/.../overlays/beta/...`) — show which apps got which image tag bumps. This is the operator's audit surface before merging.

### Step B5: Operator merge — HARD GATE

⛔ **DO NOT MERGE.** Skill MUST NOT call `gh pr merge`, MUST NOT add `auto-merge` label, MUST NOT use `--auto`. This is by design — operator review the diff manually.

**Dry-run behavior**: skip the AskUserQuestion + merge verification entirely; echo what would happen.

```bash
if [ "$DRY_RUN" = "true" ]; then
  echo "[dry-run] would ask operator to merge gitops PR manually + verify state=MERGED"
  MERGE_SHA="<dry-run-merge-sha>"
else
  : # AskUserQuestion + verify block below runs only in LIVE mode
fi
```

In LIVE mode, use `AskUserQuestion`:
```
question: "Beta gitops PR ready for review. Merge it manually in browser, then confirm here."
options:
  - label: "Merged — proceed to ArgoCD sync"
  - label: "Reviewing — wait"
  - label: "Abort"
```

Loop on "Reviewing — wait" (re-ask). On "Abort" → STOP with no cleanup needed.

**After operator says "Merged"**, VERIFY before proceeding (don't trust the operator's "Merged" claim — they may have confused windows). For each PR in `$GITOPS_PR_JSON`:

```bash
# LIVE only — block guarded by the dry-run if/else above.
MERGE_STATE=$(gh pr view "$PR_NUMBER" --repo "$BETA_GITOPS_REPO" \
  --json mergeCommit,state --jq '{state:.state, sha:(.mergeCommit.oid // "")}')
PR_STATE=$(echo "$MERGE_STATE" | jq -r '.state')
MERGE_SHA=$(echo "$MERGE_STATE" | jq -r '.sha')

if [ "$PR_STATE" != "MERGED" ]; then
  echo "⛔ PR #$PR_NUMBER state is $PR_STATE, not MERGED. Did the merge actually happen?"
  # Loop back to AskUserQuestion above — give operator another chance to merge.
fi
echo "  ✓ PR #$PR_NUMBER merged at $MERGE_SHA"
```

### Step B6: ArgoCD refresh + sync

```bash
# Parse $BETA_APPS — tw-config.sh 3.19.1+ returns YAML lists as newline-separated
# items. `tr '\n' ' '` joins to space-separated for the for-loop.
APPS=$(echo "$BETA_APPS" | tr '\n' ' ' | tr -s ' ')

if [ "$DRY_RUN" = "true" ]; then
  for app in $APPS; do
    echo "[dry-run] kubectl --context=$BETA_KUBE_CTX annotate application/${app} -n argocd argocd.argoproj.io/refresh=normal --overwrite"
  done
  echo "[dry-run] would list applications + identify OutOfSync"
  echo "[dry-run] would kubectl patch each OutOfSync app with sync operation"
  echo "[dry-run] would auto-retry any phantom-Succeeded variant (Hard Rule H8)"
  OUTOFSYNC=""  # nothing to wait on in B7
else
  # Force refresh (avoid 3-min auto-poll wait):
  for app in $APPS; do
    kubectl --context="$BETA_KUBE_CTX" annotate application/${app} \
      -n argocd argocd.argoproj.io/refresh=normal --overwrite
  done
  sleep 10

  # Show current sync states:
  kubectl --context="$BETA_KUBE_CTX" get applications -n argocd -o json \
    | jq -r '.items[] | "\(.metadata.name) \(.status.sync.status) \(.status.health.status)"'

  # Sync OutOfSync apps via kubectl patch (NOT argocd CLI — Hard Rule H4):
  OUTOFSYNC=$(kubectl --context="$BETA_KUBE_CTX" get applications -n argocd -o json \
    | jq -r '.items[] | select(.status.sync.status=="OutOfSync") | .metadata.name')

  SYNC_PATCH='{"operation":{"sync":{"revision":"HEAD","prune":false,"syncOptions":["CreateNamespace=true","RespectIgnoreDifferences=true","ApplyOutOfSyncOnly=false"]}}}'

  for app in $OUTOFSYNC; do
    kubectl --context="$BETA_KUBE_CTX" patch application/${app} -n argocd --type merge -p "$SYNC_PATCH"
  done

  # Hard Rule H8: Phantom-Succeeded variant B detection + automatic single retry.
  # Wait briefly for sync attempt to land, then check each app for the "Succeeded
  # with empty resources but still OutOfSync" pattern. If found, re-run the SAME
  # patch once. Don't escalate on first occurrence.
  sleep 20
  for app in $OUTOFSYNC; do
    APP_JSON=$(kubectl --context="$BETA_KUBE_CTX" get application/${app} -n argocd -o json)
    SYNC_STATUS=$(echo "$APP_JSON" | jq -r '.status.sync.status')
    OP_PHASE=$(echo "$APP_JSON" | jq -r '.status.operationState.phase // empty')
    RES_COUNT=$(echo "$APP_JSON" | jq -r '.status.operationState.syncResult.resources // [] | length')

    if [ "$SYNC_STATUS" = "OutOfSync" ] && [ "$OP_PHASE" = "Succeeded" ] && [ "$RES_COUNT" = "0" ]; then
      echo "  [H8] Phantom-Succeeded detected on $app — re-running same patch (auto-retry, no escalation)"
      kubectl --context="$BETA_KUBE_CTX" patch application/${app} -n argocd --type merge -p "$SYNC_PATCH"
    fi
  done
fi
```

> **⚠️ Phantom-Succeeded variant B (Hard Rule H8)**: Codified above. Detection signal: `sync.status=OutOfSync` AND `operationState.phase=Succeeded` AND `operationState.syncResult.resources=[]` (empty array). Observed 2× now (2026-04-20 staging-vi-api-server, 2026-05-11 beta-vi-api-server). The auto-retry is single-shot — if Step B7 polling still observes OutOfSync after the retry has been given time to land, that's a real problem (not phantom) and B7 will STOP with the operator told to escalate.

### Step B7: Wait for Synced/Healthy

```bash
if [ "$DRY_RUN" = "true" ]; then
  echo "[dry-run] would poll each formerly-OutOfSync app for Synced/Healthy (10min timeout)"
else
  # Poll until all formerly-OutOfSync apps are Synced + Healthy. Timeout 10min.
  DEADLINE=$(( $(date +%s) + 600 ))
  while [ "$(date +%s)" -lt "$DEADLINE" ]; do
    STILL_SYNCING=0
    for app in $OUTOFSYNC; do
      STATE=$(kubectl --context="$BETA_KUBE_CTX" get application/${app} -n argocd -o json \
        | jq -r '"\(.status.sync.status)/\(.status.health.status)"')
      case "$STATE" in
        Synced/Healthy) ;;
        */Degraded)
          echo "⛔ App $app is Degraded — investigate."
          # STOP
          ;;
        *) STILL_SYNCING=1 ;;
      esac
    done
    [ "$STILL_SYNCING" -eq 0 ] && break
    sleep 15
  done

  if [ "$STILL_SYNCING" -eq 1 ]; then
    echo "⛔ Timeout waiting for beta apps to reach Synced/Healthy. Run /team doctor or check ArgoCD UI."
    # STOP
  fi
  echo "  ✓ All apps Synced + Healthy"
fi
```

> **Migration jobs** (PreSync hooks) may add 1-2min — that's expected, included in the 10-min budget.

### Step B8: Smoke — extract service_version + parity check

```bash
if [ "$DRY_RUN" = "true" ]; then
  echo "[dry-run] curl -sS --max-time 10 $BETA_SMOKE_URL | jq -r '.service_version'"
  if [ "$SOURCE" = "prod" ] && [ -n "$PROD_SMOKE_URL" ]; then
    echo "[dry-run] curl -sS --max-time 10 $PROD_SMOKE_URL | jq -r '.service_version'  # parity check"
  fi
  BETA_SV="<dry-run-beta-sv>"
  PROD_SV="<dry-run-prod-sv>"
else
  # External smoke first.
  BETA_SV=$(curl -sS --max-time 10 "$BETA_SMOKE_URL" | jq -r '.service_version' 2>/dev/null) || BETA_SV=""

  if [ -z "$BETA_SV" ] || [ "$BETA_SV" = "null" ]; then
    # Fallback: exec into fallback pod (avoids NetworkPolicy denial on naive pods — Hard Rule H5).
    if [ -n "$BETA_SMOKE_FALLBACK_POD" ] && [ -n "$BETA_SMOKE_FALLBACK_URL" ]; then
      POD=$(kubectl --context="$BETA_KUBE_CTX" get pods -l app="$BETA_SMOKE_FALLBACK_POD" \
        -o jsonpath='{.items[0].metadata.name}')
      if [ -n "$POD" ]; then
        BETA_SV=$(kubectl --context="$BETA_KUBE_CTX" exec "$POD" -- \
          sh -c "curl -sS --max-time 8 '$BETA_SMOKE_FALLBACK_URL'" | jq -r '.service_version' 2>/dev/null) || BETA_SV=""
      fi
    fi
  fi

  if [ -z "$BETA_SV" ] || [ "$BETA_SV" = "null" ]; then
    echo "⛔ Smoke failed — could not extract .service_version from beta. (Tried external + in-cluster fallback.)"
    # STOP
  fi
  echo "  beta service_version: $BETA_SV"

  # Parity check (only when source=prod and PROD_SMOKE_URL is set).
  if [ "$SOURCE" = "prod" ] && [ -n "$PROD_SMOKE_URL" ]; then
    PROD_SV=$(curl -sS --max-time 10 "$PROD_SMOKE_URL" | jq -r '.service_version' 2>/dev/null) || PROD_SV=""
    echo "  prod service_version: $PROD_SV"
    if [ "$BETA_SV" != "$PROD_SV" ]; then
      echo "⚠️  PARITY MISMATCH — beta($BETA_SV) ≠ prod($PROD_SV). Investigate before further promote."
      # Do NOT STOP — operator decides; report and continue to B9.
    else
      echo "  ✓ Parity confirmed: beta == prod ($BETA_SV)"
    fi
  fi
fi
```

> **Hard Rule H7**: NEVER skip body extraction. HTTP 200 alone is insufficient — cosmetic bugs (conflict markers in version strings, wrong feature flags) return 200. The `jq -r '.service_version'` extraction is mandatory.

### Step B9: Final summary

**Dry-run mode**: Render a stub summary clearly marked `[DRY-RUN]` so operator knows no actions executed. Use the same template but substitute the placeholder values bound at each step.

**LIVE mode**: Render the full summary with real values from B2-B8.

**═════ BETA PROMOTE COMPLETE *(or `[DRY-RUN]` if `$DRY_RUN=true`)* ═════**
Workflow: **{RUN_URL}** (conclusion: **{conclusion}** OR `<dry-run>`)
Gitops PR: **{PR_URL}** (merge SHA: `{MERGE_SHA[0:8]}` OR `<dry-run>`)
Apps synced: *{count} of {total}*  *(in dry-run: list APPS as "would sync")*
  *(list each: APP_NAME → image_tag, or `<dry-run>` when in dry-run mode)*
Smoke:  `{BETA_SV}`  *(in dry-run: `<dry-run>`)*
  *(if parity check ran)* Prod parity: ✅ matched OR ⚠️ mismatch (prod=`{PROD_SV}`)
══════════════════════════════════════

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

### Step P4.4: Hotfix Divergence Guard (squash-only mode)

**Why this exists**: When a hotfix lands on `$PROD_BRANCH` (via `/team-rc hotfix` or manual squash) and the RC branch contains a manually-backported version of the same fix with **different content** (different commit message, different patch shape, conflict-resolution drift), Step P4.5's `git commit-tree` rebuild silently keeps the RC's version and discards the prod version's tree. The release ships with the wrong fix and nobody notices until prod regresses. This step detects the SHA-level divergence BEFORE rebuild and asks the user to reconcile.

**When this runs**: Same trigger as P4.5 (dual-branch + `release.strategy: rebuild-on-prod`). For single-branch or `cut-only` repos, this step is a no-op.

```bash
# Bind REBUILD_STRATEGY here — first use is in this step. P4.5 re-reads idempotently
# for defense-in-depth (works if step is run standalone after a session restart).
REBUILD_STRATEGY=$(bash ~/.claude/commands/scripts/tw-config.sh release.strategy "cut-only")

if [ "$BASE_BRANCH" = "$PROD_BRANCH" ] || [ "$REBUILD_STRATEGY" != "rebuild-on-prod" ]; then
  echo "Hotfix divergence check skipped (strategy=$REBUILD_STRATEGY, single-branch or cut-only)"
  # Jump to Step P4.5
else
  HOTFIX_OUTPUT=$(bash ~/.claude/commands/scripts/tw-git.sh hotfix-divergence-check \
    "$RC_BRANCH" "$PROD_BRANCH" "$VERSION_PREFIX")
fi
```

**Parse `HOTFIX_OUTPUT`** (AI step, KV format):

- `LAST_CLEAN_TAG=<tag>` — most recent clean release tag (e.g., `V0.1.27`); empty if no prior release
- `HOTFIX_TAGS<<EOF ... EOF` — heredoc list of `TAG SHA SUBJECT` for hotfix tags newer than `LAST_CLEAN_TAG` (or `HOTFIX_TAGS=` if none)
- `HOTFIX_COMMITS<<EOF ... EOF` — heredoc list of `SHA SUBJECT` for untagged hotfix-prefixed commits (fallback signal; or `HOTFIX_COMMITS=` if none). **Informational only** — every entry that's also missing on RC already appears in `MISSING_HOTFIXES` with `TAG=(untagged)`. Don't display this field separately.
- `MISSING_HOTFIXES<<EOF ... EOF` — heredoc list of `TAG SHA SUBJECT` for hotfixes whose SHA is NOT an ancestor of `origin/$RC_BRANCH` (TAG=`(untagged)` for fallback hits; or `MISSING_HOTFIXES=` if none). **This is the actionable list.**
- `STATUS=clean | missing | no-prior-release`

**If `STATUS=clean` or `STATUS=no-prior-release`**: No divergence detected — skip to Step P4.5.

**If `STATUS=missing`**, display the warning:

⚠️ **HOTFIX DIVERGENCE** ── *{RC_BRANCH}* missing fixes from *{PROD_BRANCH}*
  Last clean release: **{LAST_CLEAN_TAG}**
  Missing on RC ({count} hotfix{es}):
{for each line in MISSING_HOTFIXES: indented `{TAG}` `{SHA[0:12]}` — *{SUBJECT}*}

**Why dangerous**: Step P4.5 rebuilds RC tree on prod's parent. RC's tree wins entirely — any prod-only changes (including these hotfix patches) get silently overwritten. False positives are possible (squash-only repos: same content, different SHAs) — but a false positive blocks ship; a missed divergence ships broken code.
────────────────────────────────────────────────────────────

Use `AskUserQuestion`:
```
question: "RC 缺少 prod 上的 hotfix。继续 promote 会用 RC 的 tree 覆盖这些 fix。如何处理？"
options:
  - label: "Reconcile (推荐)"
    description: "中止 promote。手工把 prod hotfix 合回 main / RC（参考 §Hotfix Divergence Guard），重新 cut RC 后再 promote。"
  - label: "Force (确认 RC 已含等价内容)"
    description: "我已核对每个 missing hotfix，RC 用不同 SHA 包含等价修改。确认强制继续到 P4.5。"
  - label: "Abort"
    description: "取消本次 promote，稍后处理。"
```

If "Reconcile" or "Abort" → **STOP**
  💡 Reconcile recipe → 见 §Hotfix Divergence Guard 章节 Option A。
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

If "Force" → record decision (the user has manually verified equivalence; downstream P4.5 still runs full file-level diff and asks separately about ADDED/MODIFIED).

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
- `ADDED_ON_PROD<<EOF ... EOF` — heredoc list of files present on `$PROD_BRANCH` but not on `$RC_BRANCH`; rebuild deletes them from RC (empty line `ADDED_ON_PROD=` if none)
- `ADDED_BLAME<<EOF ... EOF` — **3.18.0 addition, 3.19.0 delimiter fix**: per-file last-touch info on prod, format `PATH|||HASH|||ISO_DATE|||EMAIL|||SUBJECT` (triple-pipe delimiter to tolerate single `|` characters in commit subjects). Use this to spot orphan hotfixes — if the last-touch subject matches `^hotfix(` or the date is recent (post `LAST_CLEAN_TAG`), the file is a candidate for "missed backport" and should NOT be silently deleted.
- `MODIFIED_ON_PROD<<EOF ... EOF` — heredoc list of files with **different content** on prod vs RC; rebuild overwrites prod's version with RC's (empty line `MODIFIED_ON_PROD=` if none) — **3.17.0 fix**: this case is the silent hotfix-overwrite class. Step P4.4 should have already gated obvious cases; this is the file-level safety net.
- `MODIFIED_BLAME<<EOF ... EOF` — **3.18.0 addition, 3.19.0 delimiter fix**: same `PATH|||HASH|||ISO_DATE|||EMAIL|||SUBJECT` format for modified files. Same `^hotfix(` heuristic applies.
- `STATUS=<dry-run|already-fast-forwardable>`

**If `STATUS=already-fast-forwardable`**: No rebuild needed — skip to Step P5.

**Display preview** (with blame annotations to spot orphan hotfixes):

**Render this preview by reading from heredocs**:
- For each line in `ADDED_BLAME` heredoc: split on `|||` to extract 5 fields (PATH, HASH, ISO_DATE, EMAIL, SUBJECT). Render one indented line per file using the format below.
- Same for `MODIFIED_BLAME`.
- Compute `LAST_CLEAN_DATE` from `git log -1 --format=%ai <LAST_CLEAN_TAG>` so you can compare ISO_DATE > LAST_CLEAN_DATE for the orphan-hotfix heuristic.

**REBUILD RC** ── *{RC_BRANCH}* → *{PROD_BRANCH}* ────────────
  Old rc SHA: `{OLD_SHA[0:12]}`
  New rc SHA: `{NEW_SHA[0:12]}`

  Files prod-only (will be deleted by rebuild): `{added_count}`
  *(per file, one indented line; mark ⚠️ ORPHAN HOTFIX RISK if SUBJECT starts with `hotfix(` OR ISO_DATE > LAST_CLEAN_DATE)*

   `{path}` — last touched `{HASH[0:8]}` on *{ISO_DATE}* by *{EMAIL}* — "*{SUBJECT}*"
   ⚠️ ORPHAN HOTFIX RISK — verify before delete   *(only emit this line when heuristic matches)*

  Files diverged (RC content wins, prod content overwritten): `{modified_count}`
  *(same per-file rendering; same heuristic; phrase "verify before overwrite" instead)*
────────────────────────────────────────────────────────────

> **3.19.0**: The blame fields use `|||` (triple-pipe) as delimiter — single `|` in commit subjects (e.g., `feat: support a|b syntax`) used to corrupt parsing. Split on `|||` not `|`.

> **3.18.0 blame heuristic**: files last-touched by a `hotfix(...)` commit (or recently, post `LAST_CLEAN_TAG`) are flagged as orphan-hotfix candidates. The user should verify each `⚠️` flag before confirming rebuild — these are the cases where the rebuild would silently undo a production fix.

**If `ADDED_ON_PROD` is non-empty OR `MODIFIED_ON_PROD` is non-empty**, use `AskUserQuestion`:
```
question: "确认 rebuild？上方列出的 prod-only 文件会从 rc 删除；diverged 文件会被 rc 的版本覆盖（prod 上的版本丢失）。"
options:
  - label: "执行 rebuild"
    description: "提交新 commit 到 rc 分支，prod-only 文件消失、diverged 文件采用 rc 版本"
  - label: "先 cherry-pick 再来"
    description: "中止 promote，手工把缺失内容 backport 到 rc 后重新 /team-rc promote"
  - label: "中止"
    description: "取消本次 promote"
```

If "先 cherry-pick 再来" or "中止" → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

**If both lists are empty OR user confirmed**, execute the real rebuild:

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

### Step P10c: Post-promote Checklist (3.18.0)

> **Why this exists**: 在 `git tag` 之后，typical release pipeline 还有 5 步：(1) prod-deploy workflow 跑完 (2) artifact PR (e.g., gitops repo) review/merge (3) 集群 sync trigger (4) 等 healthy (5) in-cluster smoke。Skill 不强自动化（每个 repo 的 toolchain 不一样），但要让用户有 checklist 不卡住。Real incident: vi_agent V0.1.28 第一次 promote 在 P11 显示 ✅ 后 user 不知道还有 5 步，后来主动喊"引导我继续，别停下来"才接住。

**3.19.0 rewrite**: previous version used inline markdown conditionals like `*{if X set}*` which LLM renders literally. This version uses explicit prose so the LLM picks which items to show.

**First, read config and decide which items apply** (do NOT render until you've decided):

```bash
GITOPS_REPO=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.gitops_repo "" 2>/dev/null)
ARGOCD_APPS=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.argocd_apps "" 2>/dev/null)
SMOKE_URL=$(bash ~/.claude/commands/scripts/tw-config.sh deploy.prod_smoke_url "" 2>/dev/null)
# $PROD_WORKFLOW is already bound from Step P0.
```

**Render rules**:
- Always render the section header + step 1 (workflow watch) + step 5 (smoke). Step 1 omits the `$PROD_WORKFLOW` mention if it's not set.
- Render step 2 (gitops PR) **only if** `GITOPS_REPO` is non-empty.
- Render step 3 (cluster sync) **only if** `ARGOCD_APPS` is non-empty OR user's setup clearly uses ArgoCD/Flux (user can tell you).
- Renumber whatever items you actually render — don't leave gaps like "1, 3, 5".
- Substitute config values into the bash examples. NEVER leave a literal `{GITOPS_REPO}` / `{VERSION}` placeholder visible.

**Section template** (render only the items that apply, after substitution):

═════ POST-PROMOTE CHECKLIST ═════

**N. Production-deploy workflow**
   *(if `$PROD_WORKFLOW` is set)* `gh run watch $(gh run list --workflow $PROD_WORKFLOW --branch $VERSION --limit 1 --jq '.[0].databaseId') --exit-status`
   *(if unset)* — skip this step or note manual deploy required for $VERSION
   Block here until green; typical duration depends on your build matrix.

**N. Artifact-promote PR** *(only when `deploy.gitops_repo` is set; substitute value)*
   `gh pr list --repo <GITOPS_REPO_VALUE> --head "promote/prod-<VERSION_VALUE>" --state open`
   Review the diff, then squash-merge.

**N. Cluster sync** *(only when `deploy.argocd_apps` is set OR user's pipeline uses ArgoCD/Flux)*
   For ArgoCD prod (no auto-sync, manual gate): for each app, `kubectl -n argocd patch application <APP> --type merge -p '{"operation":{"initiatedBy":{"username":"team-rc-promote"},"sync":{"syncStrategy":{"hook":{}},"prune":false}}}'`
   (`argocd app sync` CLI may fail in-pod due to RBAC; patch is the universal path.)

**N. Wait Synced/Healthy** *(only when step N-1 cluster-sync was rendered)*
   Poll until `.status.sync.status=Synced` AND `.status.health.status=Healthy` for each app. Migration jobs (PreSync hooks) may add 1-2 min — that's expected.

**N. In-cluster smoke** *(ALWAYS render — not optional)*
   Local `curl` cannot validate prod ingress for many setups (corp ZTNA / NAT / NetworkPolicy). Exec from a pod with proper service-account + NetworkPolicy. If `deploy.prod_smoke_url` is set, use it; else ask user for the URL or pod they want to exec from.
   ```
   POD=$(kubectl -n <NAMESPACE> get pods -l app=<LONG_RUNNING_APP> -o jsonpath='{.items[0].metadata.name}')
   kubectl -n <NAMESPACE> exec "$POD" -- sh -c 'curl -sS <PROD_SMOKE_URL_VALUE>'
   ```
   **Inspect the response body** — not just `-w "%{http_code}"`. Cosmetic bugs (e.g., conflict markers in version strings, wrong feature flags) still return HTTP 200; only body inspection catches them. Real incident: vi_agent V0.1.28 shipped `service_version` containing literal `<<<<<<<` conflict markers; `/api/config` body inspection caught it post-deploy.

═════════════════════════════════

> **Hint**: If your deploy chain is fully automated (single tag-push → green), the first item collapses to "wait for prod_workflow green" and you may not need 2-4. If it's manual (typical for high-stakes prod), the middle items are sequential gates each requiring human review. The smoke is **always** required — automated smoke runs on staging; in-cluster prod smoke is what catches the cosmetic class.

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
  3. Confirm promote intent
  4. **Hotfix divergence guard** *(dual-branch + `release.strategy: rebuild-on-prod`)* — detect prod hotfixes missing from RC; STOP if user picks Reconcile/Abort
  5. **Rebuild RC on prod via `commit-tree`** *(same precondition)* — make PR a clean fast-forward; ask before delete/overwrite
  6. Squash merge RC → *main* (production branch)
  7. Create git tag with version
  8. Create GitHub Release with auto-generated notes
  9. Clean up RC branch
  10. Trigger production deploy (if `deploy.production_workflow` configured), or guide manual deploy
  11. **Post-promote checklist** (3.18.0) — config-aware: prod-workflow watch, artifact-PR review (if `deploy.gitops_repo`), cluster sync (if `deploy.argocd_apps`), wait healthy, in-cluster smoke (always). Educational; not auto-executed.

**WHEN TO USE**
  After RC has been tested on staging and approved.
  Flow: `/team-rc` → *(test staging)* → `/team-rc promote`

**PREREQUISITES**
  RC branch must exist (created by `/team-rc`)
  CI must be green on RC branch
  Staging deploy must be successful (if `deploy.staging_workflow` configured)
  *(Squash-only repos)* `release.strategy: rebuild-on-prod` in `config.yml` enables P4.4 + P4.5 safety gates

**SAFETY GATES (P4.4, P4.5)**
  P4.4 fires when prod hotfixes are missing from RC (would be silently overwritten by rebuild).
    Pick: **Reconcile** (run §Hotfix Divergence Guard Option A), **Force** (verified equivalent), or **Abort**.
  P4.5 fires when files diverge between RC and prod (added on prod or modified content).
    Pick: **执行 rebuild** (overwrite with RC tree), **先 cherry-pick 再来** (manual backport), or **中止**.

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
  6.5. **(optional, 3.18.0)** RC-rebuild sub-chain — cut `rc/{version}` from prod HEAD, dispatch build workflow, watch completion, atomic tag re-push (delete + recreate Release). Enables hotfix to actually land in prod for pipelines whose `promote-prod` is a staging-tag-copy operation (no rebuild on tag). **Activated by** `release.hotfix_needs_rc_rebuild: true` + `release.build_workflow: "..."` in `config.yml`. Default OFF.
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

---

## Operation Beta Help

> Triggered by `/team-rc beta help`.

Output the following and **STOP**:

**`/team-rc beta` — Beta Environment Promote**
────────────────────────────────────────────

**USAGE**
  `/team-rc beta`                                                      Interactive: prompt for source + reason
  `/team-rc beta --source=staging`                                     Take latest staging artifacts
  `/team-rc beta --source=prod --reason="<text>"`                      Snapshot current prod (parity test); reason mandatory
  `/team-rc beta --source=explicit --version-tag=V0.1.32 --reason=...` Specific tag; version-tag + reason mandatory
  `/team-rc beta --dry-run`                                            Echo-only — NO workflow dispatch, NO kubectl, NO curl

**WHAT IT DOES** (chain validated on vi_agent 2026-05-11)
  B1. Pre-flight: gh auth, repo lookup, arg parse + validation
  B2. Dispatch `release.beta_workflow` with source/reason/version_tag inputs
  B3. Watch workflow to completion (typically 5-15min)
  B4. Surface gitops promote PR + per-app image-tag diff for operator review
  B5. **HARD GATE**: operator merges the gitops PR manually — skill NEVER auto-merges
  B6. ArgoCD `refresh=normal` annotation + `kubectl patch` sync (per-app, sync-options aware)
  B7. Wait Synced/Healthy across all `deploy.beta_apps` (timeout 10min; PreSync hooks add 1-2min)
  B8. In-cluster smoke: GET `deploy.beta_smoke_url`, extract `.service_version`. For source=prod: parity check vs `deploy.prod_smoke_url`
  B9. Summary: workflow URL + gitops PR + app states + smoke result

**WHEN TO USE**
  Staging too churned-out to be a real preprod? Beta is the intermediate gate.
  Three sources match three use-cases:
    `--source=staging` — daily "promote what's tested into beta"
    `--source=prod`    — "snapshot prod for parity testing" (e.g., reproducing a prod issue)
    `--source=explicit` — "promote a specific tag" (e.g., picking a known-good RC)

**PREREQUISITES** (config in `$TEAMWORK_DIR/config.yml`)
  Required (skill aborts B1 if missing):
    `release.beta_workflow: "promote-beta.yml"`
    `deploy.beta_apps: [app-name-1, app-name-2, ...]`
    `deploy.beta_kube_context: "kbeta"`         *(kubectl context name for beta cluster)*
    `deploy.beta_smoke_url: "https://beta.example.com/api/config"`
  Optional (fallback / additional features):
    `deploy.beta_gitops_repo: "owner/gitops"`   *(falls back to deploy.gitops_repo)*
    `deploy.prod_smoke_url: "..."`              *(enables parity check for source=prod)*
    `deploy.beta_smoke_fallback_pod: "<label-app>"`  *(in-cluster smoke pod when external URL unreachable)*
    `deploy.beta_smoke_fallback_url: "http://..."`   *(internal URL used from fallback pod)*

**HARD RULES** (skill MUST NOT violate)
  H1 NEVER `gh pr merge` the gitops PR — operator manual only
  H2 NEVER add `auto-merge` label to gitops PR
  H3 `--dry-run` = echo-only, NO actions take effect
  H4 NEVER use local `argocd app sync` CLI — use `kubectl annotate` + `kubectl patch` instead
  H5 NEVER `kubectl run smoke ...` naively — beta has default-deny NetworkPolicy
  H6 `--source=prod` REQUIRES non-empty `--reason`
  H7 Smoke MUST extract `.service_version` field — HTTP 200 alone hides cosmetic bugs
  H8 Phantom-Succeeded variant (phase=Succeeded + empty syncResult.resources + still OutOfSync) → re-run SAME patch; don't escalate first occurrence

**TROUBLESHOOTING**
  "Required beta config missing" → fill the keys listed above
  Workflow run fails B3 → `gh run view <id> --log-failed`; common cause is workflow's own input validation (e.g., F4.a 2026-05-11 was a missing `matplotlib` install in CI)
  ArgoCD app still OutOfSync after patch → wait 30s; if still stuck, re-run the same patch (H8 phantom-Succeeded variant); if a 3rd attempt fails, surface to user with `kubectl describe application/<app> -n argocd`
  Smoke `service_version=null` → check `BETA_SMOKE_URL` is correct; fall back to in-cluster pod exec (skill auto-tries this); confirm `vi-frontend` (or your fallback) is healthy
  Parity mismatch on `source=prod` → don't auto-rollback; surface mismatch + let operator decide (could be deliberate, e.g., prod hotfix not yet in beta)

**REFERENCES**
  Validated end-to-end on `flair-home-stylist/vi_agent` 2026-05-11 (F4 dossier).
  F5 brief (canonical contract): vi_agent's `drive-projects/v2/beta-prod-parity-audit@2.3.0/grading/evidence/f5-brief-for-erwin.md`

---

## Hotfix Divergence Guard

> Reference for Step P4.4. Read this when the guard fires `STATUS=missing`.

### What's the divergence

In dual-branch repos with squash-only releases (`pre-launch` → `main`, every release is squash-merged), a hotfix lands on `$PROD_BRANCH` as a new commit. If the team also lands an equivalent fix on `$BASE_BRANCH` — typically as a manual backport that's NOT a `cherry-pick -x` of the prod commit — then:

- **Same problem solved twice**, with different commit metadata, different patch shape, possibly different conflict resolution.
- The two SHAs are unrelated in git's ancestry graph.
- When `/team-rc promote` runs Step P4.5, `git commit-tree RC.tree -p prod_sha` writes a commit whose tree is RC's tree wholesale. Prod's hotfix tree gets overwritten silently. Whatever differences existed (line drift, missed edge case, wrong import) ship to production.

### Why we treat false positives as acceptable

In squash-only repos, even a clean `git cherry-pick -x` produces a different SHA on RC than on prod. The guard cannot tell "different SHA, same content" apart from "different SHA, different content" without doing a per-file diff (which Step P4.5's `MODIFIED_ON_PROD` does). So the guard fires on **any** prod hotfix not present in RC's ancestry, even when content is equivalent.

The cost asymmetry justifies this:

- **False positive** (content equivalent, guard fires): one AskUserQuestion round; user picks "Force"; promote continues. Minor friction.
- **False negative** (content diverged, guard misses): hotfix silently overwritten; production regresses; incident.

Pick friction over regression.

### Option A — Reconciliation recipe

When `STATUS=missing` and you choose "Reconcile", do this:

> **Variable substitution**: the bash blocks below reference `$VERSION_PREFIX`, `$BASE_BRANCH`, `$PROD_BRANCH` — these are bound INSIDE the skill, NOT in your shell. Either set them yourself before running (`export VERSION_PREFIX=V0.1 BASE_BRANCH=pre-launch PROD_BRANCH=main`), or replace each occurrence with your literal value. Concrete example: `${VERSION_PREFIX}=V0.1`, `${BASE_BRANCH}=pre-launch`, `${PROD_BRANCH}=main`.

**1. Cancel the current RC** (do NOT promote):
```bash
/team-rc cancel
```

**2. Open a pre-launch reconciliation PR**:
```bash
git fetch origin "$PROD_BRANCH" "$BASE_BRANCH"
git checkout -b "reconcile/${VERSION_PREFIX}.merge-prod-hotfix" "origin/$BASE_BRANCH"
git merge "origin/$PROD_BRANCH" --no-ff -m "merge: reconcile ${PROD_BRANCH} hotfixes into ${BASE_BRANCH}"
# Resolve any conflicts — your manual backport vs prod hotfix.
# Choose the version that's actually correct (usually prod's tested version).
git push -u origin HEAD
gh pr create --base "$BASE_BRANCH" --title "merge: reconcile ${PROD_BRANCH} hotfixes" \
  --body "Merges prod hotfixes back into ${BASE_BRANCH} so that ${BASE_BRANCH}'s tree contains every fix currently on ${PROD_BRANCH}. Required pre-flight before next /team-rc promote."
```

**3. Merge with merge commit (NOT squash)**: This preserves prod's hotfix SHAs in `$BASE_BRANCH`'s ancestry. After merge, `git merge-base --is-ancestor <prod-hotfix-sha> origin/$BASE_BRANCH` returns true → next RC's divergence check sees `STATUS=clean`.

> ⛔ **Squash here would defeat the purpose.** A squash on the reconciliation PR collapses prod's hotfix SHAs into a single new SHA — the divergence guard would still fire on the next promote. This is the one PR in the lifecycle that MUST be a true merge commit.

**4. Verify reconciliation took**:
```bash
git fetch --tags origin "$BASE_BRANCH"
# Pass $BASE_BRANCH as the "rc" arg — the function only needs an existing ref to
# probe ancestry against. After the reconciliation merge lands on $BASE_BRANCH,
# every prod hotfix SHA is reachable from it.
bash ~/.claude/commands/scripts/tw-git.sh hotfix-divergence-check \
  "$BASE_BRANCH" "$PROD_BRANCH" "$VERSION_PREFIX"
# STATUS should now be 'clean' — every prod hotfix is now an ancestor of $BASE_BRANCH,
# so any RC cut from $BASE_BRANCH inherits them.
```

**5. Re-cut RC and resume promote**:
```bash
/team-rc            # cuts new rc/{version} from base — now contains prod's hotfix SHAs
# verify staging, then:
/team-rc promote    # P4.4 should pass clean
```

### Why this works

The reconciliation merge brings prod's hotfix commit SHAs into `$BASE_BRANCH`'s commit graph (as ancestors via the merge commit). When the next RC is cut from `$BASE_BRANCH`, those SHAs are reachable from RC. Step P4.4's `merge-base --is-ancestor` returns true → no missing entries → `STATUS=clean`.

The merge also reconciles content: any file modified by both the manual backport and the prod hotfix surfaces as a conflict, forcing the team to pick the right version (usually prod's, since it's been validated in production).

### Avoiding this in future

Best practice for hotfixes that need to land on both branches:

1. **Use `/team-rc hotfix`**: It backports via `git cherry-pick -x` (records source SHA in commit message). Squash-only release still rewrites the SHA on `$BASE_BRANCH`, but the `(cherry picked from commit ...)` trailer makes intent auditable.
2. **OR backport via merge commit, not squash**: Open a backport PR with merge-commit strategy so prod's hotfix SHA stays in `$BASE_BRANCH`'s history.
3. **OR backport, then immediately reconcile-merge prod**: Same effect as Option A, done preemptively.

Squash-only repos amplify this risk — the cost of doing nothing is silent overwrite. The guard exists because the cost of doing it wrong is too high to leave to vigilance alone.
