---
name: wing-ship
description: "Ship mission → PR. Try: /wing-ship help"
version: "0.1.0"
---

# /wing-ship — Deliver Mission

> You own the delivery. Guard the quality boundary between code and production. Every detail you catch here saves a revert later.

**Visual Encoding** (apply to ALL output — see `${CLAUDE_SKILL_DIR}/reference/ve-standard.md`):
`**bold**` → headers/labels · `` `backtick` `` → commands/paths/counts · `*italic*` → branches · `**#NNN**` → issues (bold, 3+ digits) · `────` dividers · ⛔ NO code blocks around output

**User input**: $ARGUMENTS

**Argument routing:**

| Input | Action |
|-------|--------|
| (empty) | Full ship flow (push, PR, CI, cleanup) |
| `done` | Post-merge cleanup (close Issue, update labels, clean worktree) |
| `done help` | Show done operation usage |
| `review` | AI code review on current PR (load `${CLAUDE_SKILL_DIR}/operations/review.md`) |
| `review help` | Show review operation usage |
| `sync` | Rebase current branch on base branch |
| `sync help` | Show sync operation usage |
| `help` or `-h` | Show usage guide |

If `$ARGUMENTS` is `help` or `-h`, output the following and **STOP**:

**`/wing-ship` — Deliver your mission**

**USAGE**
  `/wing-ship`            Push + create PR + CI check + cleanup
  `/wing-ship done`       After PR merge: close Issue, update labels, clean worktree
  `/wing-ship review`     AI review current PR (correctness/security/architecture/quality)
  `/wing-ship sync`       Rebase current branch on latest base branch

**SHIP FLOW**
  `1.` Pre-flight: verify branch, code tasks done, tests pass
  `2.` Push branch → create PR with `Closes #N` or `Refs #N`
  `3.` Watch CI → request review → cleanup

NEXT: After merge → `/wing-ship done`

---

## Step 0: Prerequisites

```bash
CURRENT=$(git branch --show-current)
bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" protect-check 2>/dev/null || {
  echo "ERROR: On protected branch ($CURRENT). Switch to mission branch first."
  # STOP
}
GH_USER=$(gh api user --jq '.login' 2>/dev/null)
WING_DIR=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" detect-dir 2>/dev/null) || {
  echo "No wing config found. Run /wing init first."
  # STOP
}
```

  💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

```bash
eval "$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" resolve-labels 2>/dev/null)"
```

## Step 1: Find Active Contract

```bash
if [ -f .mission ]; then WORKTREE_MODE=true; MISSION_ISSUE=$(cat .mission); fi
ls $WING_DIR/active/MISSION-*.md 2>/dev/null
```

- No Contract + on mission branch → suggest `/wing-claim #{N}` → **STOP**
- Multiple → "Keep one, remove the rest." → **STOP**

Read Contract. Extract: issue, title, url, branch, priority, Sub-tasks, Success Criteria, AI Notes.

### Worktree Guard

```bash
GUARD_RESULT=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" worktree-guard "$CONTRACT_BRANCH" 2>&1) || GUARD_EXIT=$?
```

If mismatch → AskUserQuestion: "新开 tab 切到 worktree" / "在主仓库继续" / "关闭 worktree 模式"

## Step 2: Pre-flight Checks

### 2a: Verify branch + sync status

Verify on Contract branch. Check behind-base count. If >30 → strongly suggest sync. If >10 → suggest sync.

### 2b: Verify sub-tasks complete

Separate unchecked into code tasks vs manual tasks.
- Unchecked code tasks → AskUserQuestion: "去完成" / "手动标记完成" / "取消 ship"
- Only manual unchecked → proceed (note for PR description)

Batch-sync checkboxes to GitHub Issue (defense-in-depth).

### 2b-bis: Comment + Review Check

Check for new Issue comments since last check. If direction-changing → AskUserQuestion.

Check existing PR review status. If CHANGES_REQUESTED → AskUserQuestion: "查看 review" / "修复后继续" / "取消 ship".

### 2c: Clean working tree

If uncommitted changes → AskUserQuestion: "Commit all" / "Review first" / "Stash and continue"

### 2d: Run tests

```bash
CI_ENABLED=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" quality.ci "false" 2>/dev/null)
```

If CI enabled → skip local pre-flight (CI guards merge). Otherwise → run service-aware tests. Fail → **STOP**.

### 2e: Side Effect Triage

Scan for non-code changes: env files, Docker, config, dependency files, new env var references, `[side-effect]` entries in AI Notes.

Classify: **P0** (blocks ship — missing env vars) · **P1** (needs action — remote ops, credentials) · **P2** (info — lockfiles, minor config)

- P0 → mandatory resolution: AskUserQuestion to add to .env.example or skip with justification
- P1 → AskUserQuestion per item: "创建 Issue" / "PR 备注" / "跳过"
- P2 → auto-include in PR body
- Zero effects → skip silently

Store results for PR body injection (`$INFRA_IMPACT_ITEMS`).

## Step 3: Push

```bash
bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" push {branch}
```

If fails → suggest `/wing-ship sync` then retry.

## Step 4: Create PR

### 4a: Check existing PR

If PR already exists → skip creation, use existing. Display CI/review/merge status.

### 4b: Generate PR body

Determine close keyword: all checkboxes `[x]` → `Closes #{issue}` | unchecked remain → `Refs #{issue}`

PR body from Contract: Objective, Changes (completed code tasks), Manual Steps (if any), Success Criteria, Test, AI Notes, Infra Impact (if side effects).

**Success Criteria FILTER RULE**: Drop any criterion whose completion requires the PR itself to be merged or CI to pass (e.g., "PR 合入", "CI 通过"). These are delivery-mechanism criteria — valid in the Issue but self-referential in a PR. Only include criteria about the actual work product.

If using `Refs` (unchecked tasks remain): append warning block listing the unchecked tasks + "Close manually: `gh issue close {issue}`".

### 4c: Create PR

```bash
PR_RESULT=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-pr.sh" create {issue} "${COMMIT_TYPE}: {title}" "$PR_BODY" "$BASE_BRANCH" "{branch}")
```

### 4d-f: Post comment + labels + notification

All non-fatal: warn on failure, continue.

## Step 5: CI Check

If CI enabled → `wing-pr.sh watch {pr}`. Pass → ✅. Fail → ❌ warn but don't stop (PR exists, user iterates).

## Step 5b: AI-Suggested PR Comment (conditional)

Evaluate whether the PR needs an additional comment beyond the body:

| Signal | Why comment |
|--------|------------|
| Deployment needs special steps not in Infra Impact | Merger needs before/after actions |
| Specific files need careful review | Focus reviewer attention |
| Known limitation or tech debt introduced | Transparency |
| CI failed on flaky test | Prevent reviewer confusion |
| Breaking change or API change | Downstream consumers |

If none → skip silently. If any → draft comment → AskUserQuestion: "发送" / "编辑后发送" / "跳过"

## Step 6: Request Review

If `quality.review_required` → find tech-lead → request review. Non-fatal.

## Step 7: Cleanup

Update Issue labels (wip → review). Remove Contract if CI passed or not configured.

## Step 8: Output

**📦 SHIPPED** ── **#{issue}** {title} ─────────────
🔀 *{branch}*  PR: {pr-url}
{ci_icon} CI: `{status}`  Review: `{status}`

**📋 SUB-TASKS DELIVERED**
  ✅ {task 1}
  ✅ {task 2}
────────────────────────────────────────────

AskUserQuestion: "等 merge 后 `/wing-ship done`" / "Review PR" (`/wing-ship review`) / "回 dashboard"

---

## Operation Done (post-merge cleanup)

> Triggered by `/wing-ship done`.

### D1: Find Issue

From: Contract → `.mission` → branch name → `$ARGUMENTS`

### D2: Verify PR merged

```bash
bash "${CLAUDE_SKILL_DIR}/scripts/wing-pr.sh" verify-merged "$BRANCH"
```

If not merged → "PR must be merged first." → **STOP**

### D3-D4: Close Issue + Update labels

Close if still open (`Refs` case). Transition labels review → done.

### D5: Clean worktree (if applicable)

AskUserQuestion: "Remove worktree?" → yes: remove | no: warn about stale worktree.

### D6-D7: Return to base branch + Remove Contract

Stash uncommitted changes if dirty. Switch to base. Re-detect WING_DIR. Delete Contract.

### D8: Output

**🏁 DONE** ── **#{issue}** {title} ────────────────
Issue: `closed`  PR: `{url}` (merged)  Branch: *{base}*
────────────────────────────────────────────

AskUserQuestion: "Claim 下一个" (`/wing-claim`) / "Release" (`/wing-rc`) / "Dashboard" (`/wing`)

---

## Operation Sync (rebase on base)

> Triggered by `/wing-ship sync`.

### S1: Verify on mission branch

If on base → "Nothing to sync." → **STOP**

### S2: Fetch and rebase

```bash
bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" rebase "$BASE_BRANCH"
```

If conflicts → AI resolves: read both sides, merge intent, preserve mission changes. Flag ambiguous resolutions.

### S3: Summary + force-push

```bash
git push --force-with-lease origin {branch}
```

### S4: Info Sync

After rebase, check GitHub for new comments, PR review status, Contract freshness. Display advisory info.

---

## Sub-Operation Help

If `done help` → output and **STOP**:

**`/wing-ship done`** — Post-Merge Cleanup
  `1.` Verify PR merged  `2.` Close Issue  `3.` Labels review→done  `4.` Clean worktree  `5.` Return to base  `6.` Remove Contract
  When: after PR merged on GitHub. "No merged PR" → wait a few seconds or check `gh pr view`.

If `review help` → output and **STOP**:

**`/wing-ship review`** — AI Code Review
  `1.` Find PR  `2.` Fetch diff  `3.` Read changed files  `4.` Review (correctness/security/architecture/quality/performance)  `5.` Choose: approve / request changes / comment  `6.` Publish
  When: before merge, or when teammate asks for review.

If `sync help` → output and **STOP**:

**`/wing-ship sync`** — Rebase on Base Branch
  `1.` Verify on mission branch  `2.` Fetch + rebase  `3.` AI resolves conflicts  `4.` Force-push  `5.` Info sync (comments, PR status)
  When: dashboard shows "behind base", or before `/wing-ship`.

---

## Error Handling

- Config missing → "Run `/wing` first." → **STOP**
- No Contract → "Run `/wing-claim` first." → **STOP**
- Incomplete code tasks → warn + **STOP** (MANUAL tasks don't block)
- Tests fail → "Fix tests before shipping." → **STOP**
- Push/PR creation fails → **STOP** with remediation
- CI fails → warn but continue (PR exists)
- Label/reviewer/notification fails → warn, continue (non-fatal)

**On any STOP:** 💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.
