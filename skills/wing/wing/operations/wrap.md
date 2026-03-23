# Operation Wrap — Retroactive Closure

> Triggered by `/wing wrap`. User changed code first, now needs to back-fill: issue + PR + labels.
> Auto-detects code state and picks the right strategy.

## W0: Prerequisites

```bash
GH_USER=$(gh api user --jq '.login' 2>/dev/null)
WING_DIR=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" detect-dir 2>/dev/null) || { echo "Run /wing init first."; exit 1; }
eval "$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" resolve-labels 2>/dev/null)"
BASE_BRANCH=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" conventions.base_branch "main" 2>/dev/null)
REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
MILESTONE=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" versions.current "" 2>/dev/null)
if [ -n "$MILESTONE" ]; then
  MILESTONE=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" milestone-resolve "$MILESTONE" 2>/dev/null)
fi
```

If active Contract exists → "You have an active mission. Use `/wing-drive` + `/wing-ship` instead." → **STOP**

## W1: Detect Code State

```bash
CURRENT_BRANCH=$(git branch --show-current)
UNCOMMITTED=$(git status --porcelain)
if git rev-parse --verify "origin/${CURRENT_BRANCH}" >/dev/null 2>&1; then
  UNPUSHED=$(git log "origin/${CURRENT_BRANCH}..HEAD" --oneline 2>/dev/null)
  HAS_REMOTE_TRACKING=true
else
  UNPUSHED=$(git log "origin/${BASE_BRANCH}..HEAD" --oneline 2>/dev/null)
  HAS_REMOTE_TRACKING=false
fi
if [ "$HAS_REMOTE_TRACKING" = true ]; then
  PUSHED_AHEAD=$(git log "origin/${BASE_BRANCH}..origin/${CURRENT_BRANCH}" --oneline 2>/dev/null)
fi
```

| Condition | State |
|-----------|-------|
| UNCOMMITTED non-empty, UNPUSHED empty | **UNCOMMITTED** |
| UNPUSHED non-empty | **COMMITTED_NOT_PUSHED** |
| HAS_REMOTE_TRACKING + PUSHED_AHEAD non-empty | **COMMITTED_PUSHED** |
| None | "No changes. Nothing to wrap." → **STOP** |

## W2: Gather Change Summary

AI analyzes diff/commits to generate: Title, Priority, Size, Objective, Success Criteria.

- **UNCOMMITTED**: `git diff --stat` + `git diff`
- **COMMITTED_NOT_PUSHED**: `git log origin/$BASE_BRANCH..HEAD` + `git diff origin/$BASE_BRANCH..HEAD`
- **COMMITTED_PUSHED**: `git log origin/$BASE_BRANCH..HEAD --oneline` (or last 10)

## W3: Preview & Confirm

🔄 **WRAP PREVIEW**
────────────────────────────────────────────
**Title:** {title}  **Priority:** `{Pn}`  **Size:** `{size}`
**State:** `{state}`  **Strategy:** {description}
────────────────────────────────────────────

AskUserQuestion: "Proceed" / "Edit" / "Cancel"

## W4: Path A — UNCOMMITTED

1. **Stash**: `git stash push -m "wing-wrap: pre-issue stash"` → verify clean
2. **Create Issue**: `gh issue create` with MC body + labels + milestone
3. **Claim**: `wing-git.sh claim-setup` → mission branch from clean base
   - Worktree enabled → STOP (stash in main repo, user applies in worktree)
4. **Apply stash**: `git stash pop` on mission branch (conflict → user resolves, drive continues)
5. **Mark all checkboxes** as checked (code already done): `wing-contract.sh toggle-task` for each
6. **Verify + Ship**: run tests → `wing-git.sh push` → `wing-pr.sh create` → close

## W5: Path B — COMMITTED_NOT_PUSHED

1. **Save SHAs**: record all commit hashes between base and HEAD
2. **Create Issue**: `gh issue create` with MC body + labels
3. **Claim**: mission branch from clean base
   - Worktree enabled → STOP (user cherry-picks in worktree)
4. **Cherry-pick**: apply saved commits onto mission branch
5. **Reset base**: `git checkout $BASE_BRANCH && git reset --hard origin/$BASE_BRANCH` (safe — commits preserved on mission branch)
   - Before reset → verify diff between mission branch and original HEAD is empty (safety check)
6. **Mark checkboxes + Verify + Ship**

## W6: Path C — COMMITTED_PUSHED

Retroactive tracking only (changes already in remote):

1. **Create Issue**: `gh issue create` with MC body + labels + `status:done`
2. **Close Issue**: `gh issue close $ISSUE --reason completed`
3. No PR, no branch — just tracking

## Output

🔄 **WRAPPED** ── **#{issue}** {title} ────────────────
  **Issue:** **#{issue}** (`{state}`)
  **PR:** {pr-url or "n/a (retroactive)"}
────────────────────────────────────────────

AskUserQuestion: "看 dashboard" (`/wing`) / "Claim 下一个" (`/wing-claim`)

💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.
