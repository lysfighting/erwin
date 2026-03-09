# /team doctor — Design Spec

> Local git health diagnostics for teamwork. Subcommand of `/team`, not a new skill.
> Status: Design complete, ready for implementation.

---

## Problem Statement

When running multiple MCs in parallel, local git state accumulates entropy:
- Changes end up on wrong branches (cross-contamination)
- Merged branches linger locally
- Stashes pile up with forgotten context
- Worktrees become prunable
- PRs sit mergeable but unmerged
- Branches fall behind base, causing merge conflicts later

The drive session on 2026-03-06 (vi_agent) demonstrated this exact scenario — took 15 minutes of manual detective work that could be systematized.

**Design principle**: Report-first. Doctor outputs a visual health report (snapshot tables + cross-referencing diagnosis + prioritized action plan). User decides what to act on. Future extension: interactive fix mode (`/team doctor --fix`).

---

## Integration Point

### Argument Routing (team.md)

Add `doctor` to the existing routing table in team.md:

```
| Input | Action |
|-------|--------|
| (existing rows...) |
| `doctor` | Run local git health diagnostics |
```

In the "Route by Argument" section:
```
- If `doctor` → jump to **Operation Doctor**
```

No new skill file. No new script (rationale below). All logic lives in a new `## Operation Doctor` section of `team.md`.

### Why No Script?

The checks are a mix of:
1. **Deterministic commands** (git branch, git stash list, git worktree list, gh pr list)
2. **Cross-referencing** (branch X corresponds to issue Y which is closed → orphan)
3. **Judgment** (uncommitted changes on branch A look like they belong to issue B)

A script can handle (1) but not (2) or (3). And (2)+(3) are where the real value is — the LLM cross-references multiple data sources and produces human-readable diagnosis. Putting (1) in a script and (2)+(3) in the skill creates an awkward split with no reuse benefit.

The entire operation is read-only (no mutations), so there's no "deterministic operation an LLM might get wrong" risk that justified scripts in other skills.

---

## Operation Doctor — Full Spec

### Prerequisites

Same as dashboard: config must exist, `gh` must be authenticated, must be in a git repo.

```bash
# Identity
GH_USER=$(gh api user --jq '.login' 2>/dev/null)

# Config
if [ -f .teamwork/config.yml ]; then
  TEAMWORK_DIR=".teamwork"
elif [ -f .teamspace/config.yml ]; then
  TEAMWORK_DIR=".teamspace"
else
  echo "NO_CONFIG"
fi

# Read config values
BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "main" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-main}"
MISSION_LABEL=$(bash ~/.claude/commands/scripts/tw-config.sh github.mc_label "" 2>/dev/null)
[ -z "$MISSION_LABEL" ] && MISSION_LABEL=$(bash ~/.claude/commands/scripts/tw-config.sh mc_label "" 2>/dev/null)
[ -z "$MISSION_LABEL" ] && MISSION_LABEL="mission"
STATUS_PREFIX=$(bash ~/.claude/commands/scripts/tw-config.sh label_prefix.status "" 2>/dev/null)
[ -z "$STATUS_PREFIX" ] && STATUS_PREFIX="status:"
```

If no config → "Teamwork not initialized. Run `/team` first." → STOP

### Data Collection Phase

Run these commands to gather all raw data. Order doesn't matter — collect everything first, analyze second.

```bash
# 1. Local branches with tracking info
git branch -vv --no-color

# 2. Remote branches (for orphan detection)
git ls-remote --heads origin 2>/dev/null

# 3. Current branch
git branch --show-current

# 4. Uncommitted changes
git status --porcelain

# 5. Stash list with dates
git stash list --format='%gd|%ci|%s'

# 6. Worktree list
git worktree list --porcelain

# 7. All open mission issues (with state, labels, assignees)
gh issue list --label "$MISSION_LABEL" --state open --json number,title,state,labels,assignees --limit 50

# 8. Recently closed/merged mission issues (for orphan branch detection)
gh issue list --label "$MISSION_LABEL" --state closed --json number,title,closedAt --limit 30

# 9. Open PRs (with CI status, mergeable)
gh pr list --author "$GH_USER" --state open --json number,title,headRefName,mergeable,reviewDecision,statusCheckRollup,url --limit 20

# 10. Recently merged PRs (for orphan branch detection)
gh pr list --author "$GH_USER" --state merged --json number,title,headRefName,mergedAt --limit 20

# 11. Active contracts
ls $TEAMWORK_DIR/active/MISSION-*.md 2>/dev/null

# 12. Branch behind/ahead counts for mission branches
for branch in $(git branch --list 'mission/*' --format='%(refname:short)'); do
  echo "$branch: $(git rev-list --left-right --count origin/$BASE_BRANCH...$branch 2>/dev/null || echo 'N/A')"
done

# 13. Open mission issue bodies (for duplicate/overlap detection — LLM cross-references titles + bodies)
gh issue list --label "$MISSION_LABEL" --state open --json number,title,body,assignees,labels --limit 50
```

### Analysis Phase — 7 Diagnostic Checks

Each check produces findings with a severity level:

| Severity | Icon | Meaning |
|----------|------|---------|
| CRITICAL | 🔴 | Blocking work or risking data loss — fix before continuing |
| WARNING | 🟡 | Should fix soon, causes friction or accumulates debt |
| INFO | 🟢 | Cleanup opportunity, low urgency |

---

#### Check 1: Orphan Branches

**Logic**: A local branch is orphan if its corresponding issue is CLOSED or its PR is MERGED.

```
For each local branch matching mission/* pattern:
  1. Extract issue number from branch name (first numeric segment after /)
  2. Check if issue number appears in closed issues list → ORPHAN
  3. Check if branch name appears in merged PRs list → ORPHAN
  4. If orphan → severity 🟢, recommend: `git branch -d {branch}`
```

**Edge cases**:
- Branch name doesn't contain an issue number → skip (not a mission branch)
- Branch has unmerged local commits not in any PR → severity 🟡 instead of 🟢
  - Detect: `git log origin/$BASE_BRANCH..{branch} --oneline` has commits not in merged PR
  - Message: "Branch has local commits that may not be merged. Verify before deleting."

---

#### Check 2: Cross-Contamination (Dirty Branches)

**Logic**: Uncommitted changes on the current branch may belong to a different MC.

This is the highest-value check and requires LLM judgment.

```
1. Get current branch name → extract its issue number
2. Get git status --porcelain → list of changed files
3. For each changed file:
   a. Read the diff (git diff -- {file} + git diff --cached -- {file})
   b. Cross-reference with open issues:
      - Does the file path or change content relate to a DIFFERENT open issue?
      - Does the diff match the description/title of another MC?
4. Group findings:
   - Changes that belong to current MC → OK
   - Changes that likely belong to another MC → 🔴 CRITICAL
   - Changes that don't match any MC → 🟡 WARNING (untracked work)
```

**Heuristics for cross-referencing** (LLM should use judgment, not rigid rules):
- File path matches another issue's title keywords (e.g., `deploy-instance.sh` change while on a Docker Hub branch, but there's an open issue about deploy migration)
- Diff content matches another issue's description
- File is completely unrelated to current branch's issue scope

**Edge case**: No uncommitted changes → skip this check, output nothing.

---

#### Check 3: Branch Staleness (Behind Base)

**Logic**: Mission branches that fall too far behind base branch will have painful merges later.

```
For each local mission/* branch:
  1. Count commits behind: git rev-list --count {branch}..origin/$BASE_BRANCH
  2. Thresholds:
     - 0-10 behind → OK (don't report)
     - 11-30 behind → 🟡 WARNING
     - 31+ behind → 🔴 CRITICAL
  3. Only report for branches with corresponding OPEN issues (skip orphans)
```

---

#### Check 4: Stale Stashes

**Logic**: Stashes older than 7 days or whose parent branch is deleted/merged are likely forgotten.

```
For each stash entry:
  1. Parse date from git stash list --format='%gd|%ci|%s'
  2. Parse branch name from stash message (usually "WIP on {branch}: ...")
  3. Check age:
     - < 7 days → OK
     - 7-30 days → 🟢 INFO
     - > 30 days → 🟡 WARNING
  4. Check if parent branch still exists:
     - Branch deleted or merged → 🟡 WARNING (stash is likely orphaned)
  5. Check size:
     - > 10 files → 🟡 WARNING (too large for stash, needs its own MC)
  6. Show stash content summary: git stash show stash@{N} --stat
```

---

#### Check 5: Worktree Health

**Logic**: Detect prunable worktrees and worktrees on deleted/merged branches.

```
Parse git worktree list --porcelain output:
  For each worktree (skip the main one):
    1. Check if marked "prunable" → 🟡 prunable
    2. Check if branch is an orphan (issue closed / PR merged) → 🟡 orphan
    3. Check if worktree directory exists on disk → if not, 🔴 broken (corrupted)
```

---

#### Check 6: Shippable PRs

**Logic**: PRs that are ready to merge but sitting idle.

```
For each open PR:
  1. Check CI status (statusCheckRollup) → all passing?
  2. Check mergeable → MERGEABLE?
  3. Check reviewDecision → APPROVED or no review required?
  4. If all three → this PR is shippable
  5. Severity: 🟡 WARNING (it's blocking progress)
```

---

#### Check 7: Issue Health (Duplicates, Overlap, Orphans)

**Logic**: LLM reads all open mission issue titles + bodies and cross-references for problems.

This check requires LLM judgment — semantic similarity detection, not string matching.

```
Sub-checks:

A. DUPLICATE DETECTION
   For each pair of open mission issues:
     1. Compare titles semantically (not string match)
     2. Compare objectives/descriptions
     3. If high overlap → 🟡 WARNING "Issues #X and #Y appear to be duplicates"
     Action: close one with `gh issue close N --reason "not planned" --comment "Duplicate of #M"`

B. UNASSIGNED ISSUES
   For each open mission issue:
     1. Check if assignees list is empty
     2. If unassigned and older than 3 days → 🟡 WARNING
     Action: assign via `gh issue edit N --add-assignee USER`

C. STALE ISSUES (no activity)
   For each open mission issue:
     1. Check if status:wip but no branch exists and no PR exists
     2. If claimed (status:wip) but no commits on branch for >7 days → 🟡 WARNING "Issue #N is stale"
     Action: re-evaluate priority or reassign

D. ORPHAN CONTRACTS
   For each file in $TEAMWORK_DIR/active/MISSION-*.md:
     1. Extract issue number from filename
     2. Check if corresponding GitHub issue is still open
     3. If issue closed but contract still exists → 🟡 WARNING
     Action: `rm $TEAMWORK_DIR/active/MISSION-N.md`
```

**Heuristics for duplicate detection** (LLM judgment):
- "fix camera on Safari" and "iOS media permission broken" are the same problem even with zero shared words
- "add rate limiting" and "rate limit upload API" overlap but may be different scope — flag as potential, not definite
- Issues sharing >50% sub-tasks are likely duplicates

---

### Output Format — 3-Layer Visual Report

The report has 3 layers, inspired by the drive session's diagnostic output:

1. **Snapshot tables** — structured data, one table per dimension (branches, PRs, stashes, worktrees, issues)
2. **Diagnosis** — LLM cross-referencing analysis (the high-value part humans can't do with `git status`)
3. **Action Plan** — prioritized fix list, grouped by urgency

#### Layer 1: Snapshot Tables

Each table shows raw state with inline emoji status indicators. Tables use box-drawing characters for visual clarity. Only show tables that have content worth reporting (skip empty/clean dimensions).

**Branches table** — all local mission/* branches:

Status column values: `🔴 dirty` (uncommitted changes), `🟢 clean`, `🟡 stale` (>30 commits behind), `💀 orph` (issue closed / PR merged)

**Open PRs table** — only if there are open PRs:

Ready? column: `🟢 SHIP` (all green), `🔴 CI fail`, `🟡 review`, `🔴 conflict`, `⏳ pending`

**Stashes table** — only if stashes exist:

Verdict column: `🟢 KEEP` (recent + active branch), `🟡 REVIEW` (old but active branch), `🟡 LARGE` (>10 files), `🔴 DROP` (parent branch deleted/merged)

**Worktrees table** — only if non-main worktrees exist:

Health column: `🟢 active`, `🟡 prunable`, `🟡 orphan` (branch merged/deleted), `🔴 broken` (dir missing)

**Issues table** — only if issues have problems (duplicates, unassigned, stale):

Health column: `🟡 DUP?` (potential duplicate), `🟡 NOBODY` (unassigned >3d), `🟡 STALE` (wip but no activity), `🟢 OK` (healthy)

#### Layer 2: Diagnosis

After the tables, the LLM cross-references all data and outputs narrative findings. This is where the real intelligence lives — connecting dots across dimensions that no single `git` command can show.

Rules:
- Use emoji severity prefix: 🔴 / 🟡 / 🟢
- Group related findings (e.g., multiple orphan branches in one finding, not separate entries)
- Explain the cross-dimensional insight (WHY this is a problem, not just WHAT)
- Prioritize cross-contamination and shippable PRs — these block real work
- If a check finds nothing wrong, report it as 🟢 with a one-liner (e.g., "🟢 No branch staleness — all active branches in sync with base")

#### Layer 3: Action Plan

Concrete commands, grouped by severity (🔴 first, then 🟡, then 🟢). Each action maps to a diagnosis finding.

Rules:
- Every action has a numbered step + the exact command(s) to run
- Stash drops must be ordered from highest index first (indices shift on drop)
- Cross-contamination fixes suggest both options (discard vs stash for later)
- Shippable PRs suggest both `gh pr merge` and `/team-ship done`
- Large stashes suggest creating a new MC via `/team-issue`

#### Clean State Output

If all checks pass:
```
🏥 LOCAL HEALTH — {repo name}
══════════════════════════════════════════════════
🟢 All clear.
   {N} mission branches — all active, no contamination
   {N} stashes — all recent
   {N} worktrees — all healthy
   {N} open PRs — none ready to ship yet
══════════════════════════════════════════════════
```

No tables, no diagnosis, no action plan. Just the summary.

---

## Example Output (from the 2026-03-06 vi_agent scenario)

```
🏥 LOCAL HEALTH — collov/vi_agent
   Current: mission/116-docker-hub-login-rate-limit | Base: pre-launch
══════════════════════════════════════════════════════════════════════════

📋 BRANCHES (4 local)
┌─────────────────────────────────┬────────────┬──────────┬──────────────────────────────────────────┐
│ Branch                          │ Status     │ vs Base  │ Issue / PR                               │
├─────────────────────────────────┼────────────┼──────────┼──────────────────────────────────────────┤
│ mission/116-docker-hub-* [cur]  │ 🔴 dirty   │ +6, -1   │ #116 OPEN  │ PR #117 ✅ pass, mergeable │
│ mission/130-fix-deploy-*        │ 🟢 clean   │  0,  0   │ #130 OPEN  │ no PR                     │
│ mission/104-commit-team-*       │ 💀 orph    │ +1, -12  │ #104 CLOSED│ PR #105 merged 03-01      │
│ mission/29-setup-teamwork       │ 💀 orph    │ +0, -38  │ #29  CLOSED│ PR #30  merged 02-15      │
└─────────────────────────────────┴────────────┴──────────┴──────────────────────────────────────────┘

📋 OPEN PRS (1)
┌──────┬──────────────────────────────┬──────┬───────────┬────────┬──────────┐
│  PR  │ Title                        │ CI   │ Mergeable │ Review │ Ready?   │
├──────┼──────────────────────────────┼──────┼───────────┼────────┼──────────┤
│ #117 │ fix(infra): Docker Hub login │ ✅   │ ✅        │ —      │ 🟢 SHIP │
└──────┴──────────────────────────────┴──────┴───────────┴────────┴──────────┘

📋 STASHES (4)
┌───────────┬─────────┬─────────────────────────────┬────────────────────────┬──────────┐
│ Stash     │ Age     │ Message                     │ Content                │ Verdict  │
├───────────┼─────────┼─────────────────────────────┼────────────────────────┼──────────┤
│ stash@{0} │ 3 days  │ WIP on mission/116: skills  │ 14 files (~2800 lines) │ 🟡 LARGE │
│ stash@{1} │ 5 days  │ WIP on mission/116: infra   │ 3 files                │ 🟢 KEEP  │
│ stash@{2} │ 38 days │ WIP on mission/4: cleanup   │ 2 files                │ 🔴 DROP  │
│ stash@{3} │ 42 days │ WIP on mission/18: cleanup  │ 1 file                 │ 🔴 DROP  │
└───────────┴─────────┴─────────────────────────────┴────────────────────────┴──────────┘

📋 WORKTREES (2 extra)
┌───────────────────────────────────────────────┬────────────────────┬─────────────┐
│ Path                                          │ Branch             │ Health      │
├───────────────────────────────────────────────┼────────────────────┼─────────────┤
│ .claude/worktrees/dev-skill                   │ worktree-dev-skill │ 🟢 active   │
│ .claude/worktrees/feat/infra-and-db-migration │ (missing)          │ 🟡 prunable │
└───────────────────────────────────────────────┴────────────────────┴─────────────┘

📋 ISSUES (3 problems in 8 open)
┌───────┬──────────────────────────────────┬──────────────┬─────────────┬───────────┐
│ Issue │ Title                            │ Assignee     │ Status      │ Health    │
├───────┼──────────────────────────────────┼──────────────┼─────────────┼───────────┤
│ #132  │ fix: Safari camera permission    │ liyasong     │ status:wip  │ 🟡 DUP?  │
│ #135  │ fix: iOS media access broken     │ liyasong     │ status:wip  │ 🟡 DUP?  │
│ #128  │ feat: add usage analytics        │ (none)       │ status:wip  │ 🟡 NOBODY│
│ #116  │ fix: Docker Hub login rate limit │ liyasong     │ status:wip  │ 🟢 OK    │
│ #130  │ fix: deploy migration order      │ liyasong     │ status:wip  │ 🟢 OK    │
└───────┴──────────────────────────────────┴──────────────┴─────────────┴───────────┘

🔍 DIAGNOSIS
──────────────────────────────────────────────────────────────────────────

🔴 CROSS-CONTAMINATION — mission/116 有 2 个不属于它的未提交文件
   deploy/dev-environment/deploy-instance.sh
   → diff 内容是把 migration 移到 health check 之后 — 匹配 #130 (fix deploy migration order)
   infra/terraform/main.tf
   → terraform 资源配置变更 — 没有匹配的 open MC，疑似随手改动

🔴 ORPHAN BRANCHES — mission/104, mission/29 已死
   Issue 均已关闭，PR 均已 merged（分别 5 天前和 3 周前）。
   本地分支占空间且污染 branch 列表。

🟡 PR #117 READY TO SHIP — 已就绪但未合并
   CI 全绿，无冲突，无需 review。
   这在阻塞 #130 — #130 的工作无法开始直到 #116 的改动进入 base。

🟡 ZOMBIE STASHES — stash@{2}, stash@{3} 的父分支已删除
   mission/4 和 mission/18 早已完成。内容仅是 .gitignore 清理，安全丢弃。

🟡 stash@{0} 体积异常 — 14 个文件 ~2800 行
   看起来是 skill/command 大规模重写。不应长期留在 stash，需要独立 MC。

🟡 WORKTREE 碎片 — feat/infra-and-db-migration 可清理
   分支引用已丢失，git 元数据残留。

🟢 分支落后检查 — 无严重落后
   mission/130 与 base 同步 (0, 0)。mission/116 仅落后 1 commit。

🟡 DUPLICATE ISSUES — #132 和 #135 疑似同一问题
   #132 "Safari camera permission" 和 #135 "iOS media access broken"
   都关于移动端相机/媒体权限问题。Safari 上的 camera permission = iOS media access。
   建议关闭其中一个，在另一个的 sub-tasks 中合并两者范围。

🟡 UNASSIGNED ISSUE — #128 无人认领
   "add usage analytics" 标记了 status:wip 但没有 assignee。
   要么指派给某人，要么去掉 wip 标签放回 backlog。

📌 ACTION PLAN (10 actions)
══════════════════════════════════════════════════════════════════════════

 🔴 Critical
 ───────────────────
 [1] 归位 #130 的改动（从 mission/116 移走）:
     git stash push -m "for-issue-130" -- deploy/dev-environment/deploy-instance.sh

 [2] 处理无主 terraform 改动:
     git checkout -- infra/terraform/main.tf              # 丢弃
     git stash push -m "tf-wip" -- infra/terraform/main.tf  # 或保留

 🟡 Warning
 ───────────────────
 [3] 合并 PR #117（解除 #130 的阻塞）:
     gh pr merge 117 --squash → /team-ship done

 [4] 丢弃僵尸 stash（从高 index 开始）:
     git stash drop stash@{3}
     git stash drop stash@{2}

 [5] 为 stash@{0} 创建 MC（14 files 太大不能一直 stash）:
     /team-issue "refactor: skill/command rewrite"

 [6] 清理 prunable worktree:
     git worktree remove .claude/worktrees/feat/infra-and-db-migration

 [7] 合并重复 Issue #132 和 #135:
     gh issue close 135 --reason "not planned" --comment "Duplicate of #132 — merging scope"
     gh issue edit 132 --body "..."  # 合并 #135 的 sub-tasks 到 #132

 [8] 为 #128 指派 assignee（或移出 wip）:
     gh issue edit 128 --add-assignee USER
     # 或: gh issue edit 128 --remove-label "status:wip"

 🟢 Cleanup
 ───────────────────
 [9] 删除孤儿分支:
     git branch -d mission/104-commit-team-skills-infra-deploy
     git branch -d mission/29-setup-teamwork

══════════════════════════════════════════════════════════════════════════
SUMMARY: 🔴 2 critical │ 🟡 6 warnings │ 🟢 1 cleanup
         Fix critical issues before continuing development.
══════════════════════════════════════════════════════════════════════════
```

---

## Implementation Plan

### Files to modify:
1. **`skills/teamwork/team.md`** — Add argument route + Operation Doctor section (~200 lines)

### Files NOT needed:
- No new script (all logic is LLM-driven, read-only)
- No new skill file (subcommand of existing `/team`)
- No config changes

### Implementation steps:

1. Add `doctor` to argument routing table (line ~17 in team.md)
2. Add `doctor` to Route by Argument parser (line ~33 in team.md)
3. Add `doctor` to help text (line ~47 in team.md)
4. Write `## Operation Doctor` section with:
   - Prerequisites (reuse existing config detection pattern)
   - Data collection commands (13 bash commands)
   - 7 diagnostic checks with analysis logic
   - 3-layer output format spec (tables → diagnosis → action plan)
5. Update CHANGELOG
6. Bump version in team.md frontmatter (3.0.0 → 3.1.0, minor feature)
7. Run `./install.sh teamwork` to update symlinks
8. E2E test on vi_agent repo (which has the exact mess described)

### Estimated size:
- ~200-250 lines added to team.md
- 0 new files
- Version bump: 3.0.0 → 3.1.0

---

## Design Decisions

### Why not auto-fix?

1. **Cross-contamination** requires human judgment — "does this change belong here?" is context-dependent
2. **Stash drops are irreversible** — wrong drop = lost work
3. **Branch deletion** could lose unmerged commits
4. **Consistency with teamwork philosophy** — skills inform, user decides

### Why not invoke /drive for fixes?

The fix actions are one-liner shell commands (branch -d, stash drop, worktree remove). /drive's ceremony (Phase 0 → Contract → Context Handoff → Team Assembly → Waves) is designed for multi-task missions, not 6 copy-paste commands. Overkill.

Future `/team doctor --fix` can add interactive confirmation per-action without needing /drive.

### Why not a separate dashboard section?

Dashboard is GitHub-focused (issues, PRs, team status). Doctor is local-git-focused (branches, stashes, worktrees, uncommitted changes). Mixing them would make the dashboard noisy for the common case (local state is clean). Doctor is opt-in when you feel the mess.

### Why Check 2 (cross-contamination) needs LLM?

A script can detect "uncommitted changes exist" but cannot determine "this deploy script change belongs to issue #130, not #116." The LLM reads the diff content, reads the issue titles/descriptions, and cross-references. This is the killer feature that makes doctor more useful than `git status`.

### Severity mapping rationale

| Severity | Icon | Trigger | Why this level |
|----------|------|---------|----------------|
| CRITICAL | 🔴 | Cross-contamination, corrupted worktree | Risks committing wrong code to wrong PR |
| WARNING | 🟡 | Shippable PR, stale branch, old stash, large stash, prunable worktree, duplicate issues, unassigned issues, stale issues, orphan contracts | Blocking progress or accumulating debt |
| INFO | 🟢 | Orphan branch, healthy-but-notable items | Pure cleanup, no risk |

### Why Diagnosis uses Chinese?

User's working language is Chinese for discussion. The diagnosis is narrative text explaining insights, not code. Commands in the Action Plan stay English. This matches the user's natural workflow: think in Chinese, execute in English.

### Why no `close` subcommand for team-issue?

team-issue already has 7 subcommands (create, #N, #N text, update, fix, batch, help). Adding `close` increases learning cost for a marginal gain — closing an issue is `gh issue close N`, a one-liner users already know.

Instead, doctor Check 7 (Issue Health) serves as the **discovery layer**: it detects duplicates, unassigned issues, and stale issues, then the Action Plan provides the exact `gh` commands to fix them. Users don't need to memorize a new subcommand — they run `/team doctor`, see the problem, and copy the fix.

This follows the teamwork principle: **skills inform, user decides**. Doctor discovers, Action Plan prescribes, user executes.

### Emoji in tables — alignment concern?

Claude Code terminals universally support emoji. Box-drawing + emoji is visually richer than pure ASCII. Status columns use fixed-width emoji patterns (🔴/🟡/🟢 + space + 4-5 char label) to maintain alignment.

---

## v3.4.0 Changes (implemented)

### `/team doctor fix` — Interactive Fix Mode
- 🟢 Cleanup actions auto-execute with summary (branch -d, worktree remove)
- 🟡 Warning actions confirm each via AskUserQuestion (stash drop, stale branch)
- 🔴 Critical actions confirm each with risk explanation
- Cross-contamination always shows both options (discard vs preserve)

### `/team doctor help` — Self-Contained Help
- Shows usage, what each check does, severity levels, and typical workflow
- No need to go back to `/team help` to find doctor info

### vs Remote Column — Branch Remote Tracking
- New "vs Remote" column in branch table shows ahead/behind vs `origin/{branch}`
- Highlights `⚠ behind` when local is behind remote (most actionable finding)
- New Check 8: Remote Sync diagnostic with pull/sync recommendations

### Action Plan — Teamwork Command References
- Actions now suggest teamwork commands where applicable (e.g., `/team-ship sync` instead of raw `git rebase`)
- Report footer reminds user: "Run `/team doctor fix` to execute these actions interactively."

---

## Future Extensions (not yet implemented)

- **Pre-ship health gate**: Run doctor checks automatically before `/team-ship` to catch cross-contamination
- **Periodic reminder**: If doctor finds 🔴 issues, show a one-line warning at dashboard time ("⚠ /team doctor found 2 critical issues")
- **Multi-repo doctor**: For monorepo setups with multiple worktrees, aggregate health across all
