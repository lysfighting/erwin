---
name: wing-rc
description: "RC lifecycle — prepare staging or promote to production. Try: /wing-rc help"
version: "0.1.0"
---

# /wing-rc — Release Candidate Lifecycle

> You own the release cadence. Every RC is a contract with users — what you ship defines trust. Verify before promoting, and never rush a release you wouldn't bet your reputation on.

**Visual Encoding** — apply to ALL output (see `${CLAUDE_SKILL_DIR}/reference/ve-standard.md`)

**User input**: $ARGUMENTS

**Argument routing:**

| Input | Action |
|-------|--------|
| (empty) | Prepare: cut rc branch from develop → staging |
| `cancel` | Cancel: delete current RC branch |
| `promote` | Promote: load `${CLAUDE_SKILL_DIR}/operations/promote.md` |
| `fix` | Fix: cherry-pick a commit onto active RC |
| `hotfix` | Hotfix: load `${CLAUDE_SKILL_DIR}/operations/hotfix.md` |
| `help` or `-h` | Show usage guide |

If `$ARGUMENTS` is `help` or `-h`, output the following and **STOP**:

**`/wing-rc` — Release Candidate Lifecycle**

**USAGE**
  `/wing-rc`              Prepare: cut rc branch, deploy staging
  `/wing-rc cancel`       Cancel: delete rc, fix on base, re-cut
  `/wing-rc promote`      Promote: squash merge rc → *main*, tag, release
  `/wing-rc fix`          Cherry-pick commit to active RC (preserves test progress)
  `/wing-rc hotfix`       Post-release production fix + backport to develop

**FLOW**
  `/wing-rc` → *(test staging)* → `/wing-rc promote`
  Bug on staging? → Fix on *develop* → `/wing-rc cancel` + `/wing-rc` (re-cut)
  Surgical fix? → Fix on *develop* → `/wing-rc fix` (cherry-pick to RC)
  Production bug? → `/wing-rc hotfix` (fix main + backport)

---

## Prerequisites (shared by all flows)

```bash
# Worktree guard — RC must run from main repo
if [ -f .mission ]; then
  MAIN_REPO=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" worktree-main-repo)
  echo "⚠️ /wing-rc must run from main repo. cd $MAIN_REPO"
  # STOP
fi

WING_DIR=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" detect-dir 2>/dev/null) || {
  echo "No wing config found. Run /wing init first."
  # STOP
}

REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
VERSION_PREFIX=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" versions.current "" 2>/dev/null)
BASE_BRANCH=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" conventions.base_branch "develop" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-develop}"
PROD_BRANCH=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" conventions.production_branch "main" 2>/dev/null)
PROD_BRANCH="${PROD_BRANCH:-main}"
STAGING_WORKFLOW=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" deploy.staging_workflow "" 2>/dev/null)
PROD_WORKFLOW=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" deploy.production_workflow "" 2>/dev/null)
```

- No repo → **STOP**. No `versions.current` → **STOP**.
  💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

---

## Prepare Flow

### Step 1: RC Uniqueness Check

```bash
RC_BRANCH=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" find-rc 2>/dev/null)
RC_EXIT=$?
```

If RC exists → AskUserQuestion: "Promote" / "Cancel" / "退出". Route accordingly.

### Step 2: Derive Next Version

```bash
NEXT_VERSION=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" next-version "$VERSION_PREFIX")
```

### Step 3: Verify Develop CI

```bash
LATEST_CI=$(gh run list --branch "$BASE_BRANCH" --limit 1 --json conclusion --jq '.[0].conclusion' 2>/dev/null)
```

If not success → AskUserQuestion: "继续" / "中止".

### Step 4: Confirm

**PREPARE RC** ── **{NEXT_VERSION}** ────────────────────
*rc/{NEXT_VERSION}* ← *{BASE_BRANCH}*
CI: {green / not green}  Staging: {workflow / manual}

AskUserQuestion: "确认" / "取消"

### Step 5: Cut RC Branch

```bash
bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" cut-release "$NEXT_VERSION"
```

Does NOT checkout or modify working tree.

### Step 6: Trigger Staging (optional)

If `STAGING_WORKFLOW` set → `gh workflow run "$STAGING_WORKFLOW" --ref "rc/$NEXT_VERSION" -f environment=staging`

### Step 7: Summary

**RC PREPARED ── {NEXT_VERSION}** ────────────────────
*rc/{NEXT_VERSION}* ← *{BASE_BRANCH}*

**NEXT STEPS**
  1. Deploy staging, verify
  2. Bug? → Fix on *{BASE_BRANCH}*, then `/wing-rc cancel` + `/wing-rc` or `/wing-rc fix`
  3. Ready → `/wing-rc promote`

---

## Cancel Flow

### C1: Find active RC

```bash
RC_BRANCH=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" find-rc 2>/dev/null)
```

None → "No active RC to cancel." → **STOP**. Multiple → "Delete stale ones first." → **STOP**

### C2: Show info + confirm

Show commits on RC. AskUserQuestion: "确认取消" / "保留"

### C3: Delete RC branch

```bash
git push origin --delete "$RC_BRANCH" 2>/dev/null || true
git branch -d "$RC_BRANCH" 2>/dev/null || git branch -D "$RC_BRANCH" 2>/dev/null || true
```

### C4: Summary

**RC CANCELLED** ── *{RC_BRANCH}*
Next: Fix on *{BASE_BRANCH}* → `/wing-rc` to re-cut

---

## Fix Flow

> Cherry-pick a commit from base branch onto active RC. Preserves RC test progress.

### F1: Find active RC

None → "No active RC. Run `/wing-rc` first." → **STOP**

### F2: Identify fix commit

If argument provided → validate exists + not already on RC.
If no argument → show recent commits on BASE_BRANCH, AskUserQuestion to select.

### F3: Confirm cherry-pick

Show commit + target RC. AskUserQuestion: "确认" / "取消"

### F4: Execute

```bash
DIRTY=$(git status --porcelain 2>/dev/null | head -1)
if [ -n "$DIRTY" ]; then
  git stash push -m "wing-rc-fix: auto-stash"
fi
git fetch origin "$RC_BRANCH" 2>/dev/null
git checkout "$RC_BRANCH"
git pull origin "$RC_BRANCH"
git cherry-pick "$COMMIT_SHA" --no-edit
```

If conflicts → guide resolution interactively.

```bash
git push origin "$RC_BRANCH"
git checkout "$BASE_BRANCH"
# Restore stash if created
```

### F5: Summary

**FIX APPLIED** ── *{RC_BRANCH}*
Commit: `{msg}`
Next: Re-deploy staging → verify → `/wing-rc promote`

---

## Error Handling

- Worktree detected → "Run from main repo." → **STOP**
- Config missing → "Run `/wing` first." → **STOP**
- No `versions.current` → **STOP** with config example
- No RC branch (for cancel/promote/fix) → **STOP**
- Multiple RC branches → "Delete stale ones." → **STOP**
- PR/merge fails → **STOP** with details
- Cherry-pick conflicts → guide resolution, don't stop

**On any STOP:** 💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.
