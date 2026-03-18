# Operation Auto — One-Shot Claim → Drive → Ship

> Triggered by `/wing auto` or `/wing auto #N`. Forward closure in a single command.
> Stops at PR creation — human merges.

## Prerequisites

```bash
GH_USER=$(gh api user --jq '.login' 2>/dev/null)
WING_DIR=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" detect-dir 2>/dev/null) || {
  echo "No wing config found. Run /wing init first."
  # STOP
}
eval "$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" resolve-labels 2>/dev/null)"
BASE_BRANCH=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" conventions.base_branch "main" 2>/dev/null)
REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
```

### Worktree handling

```bash
WORKTREE_ENABLED=$(git config --local wing.worktree 2>/dev/null || echo "false")
IN_WORKTREE=false
[ -f .mission ] && IN_WORKTREE=true
```

- Inside worktree + worktree disabled → **STOP** (run from main repo)
- Inside worktree + worktree enabled → `AUTO_CLAIM_ONLY=true` (claim only, then stop for worktree switch)

### Pre-checks (skip if AUTO_CLAIM_ONLY)

- Clean working tree required (auto switches branches)
- No active Contract in non-worktree mode

## A1: Select Issue

**If `#N` provided** → extract Issue number, proceed.

**If no number** → list user's assigned open mission Issues:

```bash
gh issue list --label "$MISSION_LABEL" --state open --assignee "$GH_USER" --json number,title,labels --limit 10
```

Sort by priority. AskUserQuestion with top 3 + "取消". If none → **STOP**.

## A2: Validate Issue

```bash
gh issue view {issue} --json number,title,body,labels,milestone,assignees,state,url
```

- Not found / closed → **STOP**
- No assignee → auto-assign to `$GH_USER`
- Assigned to someone else → AskUserQuestion: "Claim anyway?" / "Cancel"
- Missing mission label → **STOP** (use `/wing-issue fix` first)

Display summary + auto plan:

🤖 **AUTO** ── **#{issue}** {title} ────────────────
  Phase 1: Claim → Contract + branch
  Phase 2: Drive → execute all sub-tasks
  Phase 3: Ship → push + PR

AskUserQuestion: "Start auto?" / "Cancel"

## A3: Phase 1 — Claim (streamlined)

1. Parse Issue body → extract Objective, Sub-tasks, Success Criteria, Context, Test Command
2. Generate Mission Contract (`$WING_DIR/active/MISSION-{issue}.md`)
3. Scan codebase for Context Files (Glob + Grep)
4. Create branch (+ worktree if enabled):

```bash
SLUG=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" slugify "$ISSUE_TITLE")
SETUP_RESULT=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" claim-setup \
  "$ISSUE_NUMBER" "$SLUG" "$GH_USER" \
  "$WING_DIR" "$WING_DIR/active/MISSION-$ISSUE_NUMBER.md" 2>/dev/null)
eval "$SETUP_RESULT"
```

5. Post claim comment: `gh issue comment {issue} --body "🤖 Auto-claimed by @${GH_USER}"`

**If worktree enabled → STOP here:**

✅ **PHASE 1: CLAIMED** (worktree mode)
  Worktree: `{worktree_path}`
  **打开新终端 tab：** `cd {worktree_path} && claude`
  Then: `/wing-drive` → `/wing-ship`

**If non-worktree → continue to A4.**

## A4: Phase 2 — Drive (full protocol)

Execute the complete `/wing-drive` protocol, **starting from Step 1** (Verify Branch).

Skip: Step 0 (Contract just created), Step 0b (freshness — just generated), Step 0b2 (Cross-MC), Step 0c (worktree guard).

Execute:
1. Step 1: Verify Branch
2. Step 2: Display Mission Briefing + Branch Safety
3. Plan Validation (sub-tasks are hypotheses — rewrite if wrong)
4. Execution Mode selection (sequential / parallel / team)
5. Full execution loop: implement → verify → ripple check → self-review → commit → toggle
6. Manual Ops Handoff (if MANUAL tasks exist)
7. Completion Check: all code checkboxes, acceptance criteria, full test suite

**On blocker:** Skip blocked task, continue. At end, report blockers: "Ship partial?" / "Stop here"

✅ **PHASE 2: DRIVEN**
  Code tasks: `{done}/{total}`  Commits: `{N}`

## A5: Phase 3 — Ship (streamlined)

No interactive confirmations.

1. Verify branch + sub-tasks
2. Sync checkboxes to GitHub Issue
3. Commit any uncommitted changes
4. Run tests (if configured)
5. Side Effect Triage (P0 → auto-fix .env.example, P1/P2 → add to PR body)
6. Push branch

```bash
bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" push {branch}
```

7. Create PR:

```bash
PR_RESULT=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-pr.sh" create {issue} "${COMMIT_TYPE}: {title}" "$PR_BODY" "$BASE_BRANCH" "{branch}")
```

Determine close keyword: all `[x]` → `Closes`, unchecked remain → `Refs`.

8. Post comment + labels + notification (all non-fatal)
9. Watch CI (if configured) — report but don't stop

## A6: Summary

🤖 **AUTO COMPLETE** ── **#{issue}** {title} ────────────────
  ✅ Phase 1: Claimed
  ✅ Phase 2: Driven (`{tasks}` tasks, `{commits}` commits)
  ✅ Phase 3: Shipped → PR {pr-url}
  CI: {status}

  **Merge the PR on GitHub, then:** `/wing-ship done`

💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

**STOP**
