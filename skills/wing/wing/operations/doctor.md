# Operation Doctor — Health Diagnostics

> Read-only local git + issue health diagnostics with 3-layer visual report.
> context:fork candidate — heavy data collection + LLM analysis.
> Triggered by `/wing doctor` (read-only) or `/wing doctor fix` (interactive repair).

## Prerequisites

```bash
GH_USER=$(gh api user --jq '.login' 2>/dev/null)
WING_DIR=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" detect-dir 2>/dev/null) || {
  echo "Wing not initialized. Run /wing first."
  # STOP
}
BASE_BRANCH=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" conventions.base_branch "main" 2>/dev/null)
eval "$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" resolve-labels 2>/dev/null)"
```

## Data Collection (parallel)

```bash
# 1. Local branches with tracking: git branch -vv --no-color
# 2. Remote branches: git ls-remote --heads origin
# 3. Current branch + uncommitted changes: git branch --show-current; git status --porcelain
# 4. Stash list: git stash list --format='%gd|%ci|%s'
# 5. Worktree list: git worktree list --porcelain
# 6. Open mission issues: gh issue list --label "$MISSION_LABEL" --state open --json number,title,state,labels,assignees,body --limit 50
# 7. Closed mission issues: gh issue list --label "$MISSION_LABEL" --state closed --json number,title,closedAt --limit 30
# 8. Open PRs: gh pr list --author "$GH_USER" --state open --json number,title,headRefName,mergeable,reviewDecision,statusCheckRollup,url --limit 20
# 9. Merged PRs: gh pr list --author "$GH_USER" --state merged --json number,title,headRefName,mergedAt --limit 20
# 10. Active contracts: ls $WING_DIR/active/MISSION-*.md
# 11. Branch behind/ahead for mission branches
# 12. GitHub auth status: gh auth status
# 13. Contract freshness data (hash comparison)
# 14. Test command: wing-config.sh project.test_command
```

## 12 Diagnostic Checks

| # | Check | Severity |
|---|-------|----------|
| 1 | **Orphan Branches** — issue closed or PR merged | 🟢 INFO |
| 2 | **Cross-Contamination** — uncommitted changes belong to different MC (LLM judgment) | 🔴 CRITICAL |
| 3 | **Branch Staleness** — behind base: 11-30 → 🟡, 31+ → 🔴 | 🟡/🔴 |
| 4 | **Stale Stashes** — age >7d, parent branch deleted, >10 files | 🟢/🟡 |
| 5 | **Worktree Health** — prunable, orphan, broken, idle, context mismatch | 🟡/🔴 |
| 6 | **Shippable PRs** — CI passed + mergeable + approved but idle | 🟡 |
| 7 | **Issue Health** — duplicates (semantic), unassigned >3d, stale wip, orphan contracts | 🟡 |
| 8 | **Remote Sync** — local behind own remote tracking branch | 🟡 |
| 9 | **GitHub Auth** — broken auth, missing scopes (runs FIRST — skip API checks if broken) | 🔴 |
| 10 | **Contract Freshness** — Issue updated since claim, hash mismatch | 🟡 |
| 11 | **Test Suite** — run test command, check pass/fail (only if tree clean) | 🟡 |
| 12 | **PR CI Status** — specific failing checks, extended pending | 🟡 |

### Severity scale

| Icon | Level | Meaning |
|------|-------|---------|
| 🔴 | CRITICAL | Blocking work or risking data loss |
| 🟡 | WARNING | Should fix soon |
| 🟢 | INFO | Cleanup opportunity |

## Output Format — 3-Layer Report

### Layer 1: Snapshot Tables

Tables with box-drawing characters, only for dimensions with content:

**Branches** — Status: 🔴 dirty / 🟢 clean / 🟡 stale / 💀 orph. vs Base: +N/-M. vs Remote: +N/-M
**Open PRs** — Ready: 🟢 SHIP / 🔴 CI fail / 🟡 review / ⏳ pending
**Stashes** — Verdict: 🟢 KEEP / 🟡 REVIEW / 🟡 LARGE / 🔴 DROP
**Worktrees** — Health: 🟢 active / 🟡 prunable / 🟡 orphan / 🔴 broken
**Issues** — Health: 🟡 DUP? / 🟡 NOBODY / 🟡 STALE / 🟢 OK

### Layer 2: Diagnosis

Cross-reference all data. Narrative findings with emoji severity prefix. Group related findings. Explain WHY, not just WHAT. Prioritize cross-contamination and shippable PRs.

### Layer 3: Action Plan

Prioritized by severity (🔴 first). Each action = numbered step + exact command. Prefer `/wing-*` commands over raw git. Remind: "Run `/wing doctor fix` to execute interactively."

### Clean State

If all checks pass:

🏥 **LOCAL HEALTH** — {repo}
────────────────────────────────────────────
🟢 All clear.
────────────────────────────────────────────

**STOP**

---

## Doctor Fix Mode

> Triggered by `/wing doctor fix`. Runs full scan, then executes actions interactively.

### Step 1: Run full doctor scan (display report)

### Step 2: Interactive execution

**🟢 Cleanup (auto-execute):** Safe reversible ops — `git branch -d` (not -D), `git worktree remove`. Execute all at once, report results.

**🟡 Warning (confirm each):** AskUserQuestion per action: "Execute" / "Skip" / "Skip all remaining"

Special cases:
- Stash review: "Drop" / "Apply" / "Create mission" / "Skip"
- Idle worktree: "Remove" / "Remove + delete branch" / "Skip"
- Main-repo mismatch: "新开 tab 切到 worktree" / "Disable worktree mode" / "Skip"
- Behind remote: "Pull latest" / "Skip"

**🔴 Critical (confirm each with risk context):**

Cross-contamination: show both options (discard vs preserve).

### Step 3: Summary

🏥 **DOCTOR FIX COMPLETE**
────────────────────────────────────────────
  ✅ Executed: `{N}`  ⏭️ Skipped: `{N}`  ❌ Failed: `{N}`

**STOP**
