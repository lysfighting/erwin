---
description: "Claim Issue → Contract → Branch. Try: /team-claim help"
version: "3.12.0"
---

# /team-claim — Claim Issue → Contract → Branch

> You own this Contract. Challenge the Issue before locking it into a Contract — if the objective is vague, the criteria unverifiable, or a constraint is missing, surface it now. Assess blast radius: if this Contract is wrong, how much rework? High blast radius → flag it in the confirmation step. Watch for implicit decisions nobody explicitly chose. A bad Contract executed flawlessly is still a bad outcome.

> Claim a GitHub Issue, generate an AI-optimized Mission Contract, and create your working branch.

**Visual Encoding** (apply to ALL output — see `docs/visual-encoding-standard.md`):
`**bold**` → headers/labels (white) · `` `backtick` `` → commands/paths/counts (purple-blue) · `*italic*` → branches (dim) · `**#NNN**` → issues (light-blue clickable, 3+ digits) · `────` dividers · ⛔ NO code blocks around output

**Design Constitution** (read before modifying this skill — see `docs/teamwork-design-constitution.md`)

**User input**: $ARGUMENTS

**Argument routing:**

| Input | Action |
|-------|--------|
| `#N` or `N` | Claim specific assigned Issue |
| `list` | Browse MY assigned mission Issues |
| `help` or `-h` | Show usage guide |
| (empty) | Show my assigned Issues and pick one |

If `$ARGUMENTS` is `help` or `-h`, output the following and **STOP**:

**`/team-claim`** — Claim an assigned GitHub Issue as your mission

**USAGE**
  `/team-claim`           Show your assigned Issues and pick one
  `/team-claim` **#042**       Claim specific assigned Issue
  `/team-claim list`      Browse your assigned missions

**WHAT HAPPENS**
  `1.` Fetches Issue from GitHub
  `2.` Generates AI-enriched Mission Contract (`.teamwork/active/MISSION-N.md`)
     — scans project to discover relevant files (Context Files, standalone test)
  `3.` Confirms Contract (AI states position: agree/challenge/redirect)
  `4.` Creates branch: *mission/{issue}-{slug}-{user}*
  `5.` Posts claim comment on GitHub

NOTE: Issues are assigned via `/team-issue` (any team member can create and assign).

NEXT: `/team-drive` to start executing

---

## Step 0: Prerequisites

```bash
# Identity
GH_USER=$(gh api user --jq '.login' 2>/dev/null)
if [ -z "$GH_USER" ]; then
  echo "ERROR: Cannot get GitHub user identity. Run 'gh auth login' first."
  exit 1
fi
```

```bash
# Config (support both directory names)
TEAMWORK_DIR=$(bash ~/.claude/commands/scripts/tw-config.sh detect-dir 2>/dev/null) || {
  echo "No teamwork config found. Run /team init first."
  # STOP
}
```
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

- If no config → "**ERROR:** Teamwork not initialized. Run `/team` first." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

Read `$TEAMWORK_DIR/config.yml` → extract team roster, conventions, project settings.

```bash
eval "$(bash ~/.claude/commands/scripts/tw-config.sh resolve-labels 2>/dev/null)"
# Now MISSION_LABEL, STATUS_PREFIX, PRIORITY_PREFIX are set
```

Verify `GH_USER` is in the team roster (check both `team:` and `members:` sections).
If not → "**ERROR:** You (`{GH_USER}`) are not in the team roster. Run `/team` to add yourself (automatic onboarding)." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

---

## Step 1: Check Active Mission

```bash
ls $TEAMWORK_DIR/active/MISSION-*.md 2>/dev/null
```

**Worktree routing is handled automatically by `tw-git.sh claim-setup` in Step 5.** Here we only check for existing Contracts:

```bash
# Detect worktree mode (script handles all routing)
WORKTREE_ENABLED=$(git config --local teamwork.worktree 2>/dev/null || echo "false")
```

**If worktree DISABLED (default):** If any Contract exists → read it, display the active mission info, then offer options:

Use `AskUserQuestion`:
```
question: "已有活跃 mission: **#{issue}** — {title}。如何处理？"
options:
  - label: "继续当前 mission"
    description: "运行 /team-drive 继续执行"
  - label: "交付当前 mission"
    description: "运行 /team-ship 提交 PR"
  - label: "放弃当前 mission"
    description: "清理 Contract，释放 claim 位（代码保留在分支上）"
  - label: "取消"
    description: "不操作"
```

- If "继续当前 mission" → Use `AskUserQuestion`:
  ```
  question: "Continue with `/team-drive` now?"
  options:
    - label: "Yes, start driving"
      description: "Invoke /team-drive immediately"
    - label: "No, I'll do it manually"
      description: "Exit — run /team-drive when ready"
  ```
  - If "Yes, start driving" → invoke `/team-drive` → **STOP**
  - If "No, I'll do it manually" → **STOP**
- If "交付当前 mission" → output "Run `/team-ship`" → **STOP**
- If "放弃当前 mission" → delete Contract file → output "Contract cleared." → Use `AskUserQuestion`:
  ```
  question: "下一步？"
  options:
    - label: "Claim 新 Issue"
      description: "/team-claim — 认领另一个 Issue"
    - label: "回 Dashboard"
      description: "/team — 总览"
  ```
  - If "Claim 新 Issue" → invoke `/team-claim` → **STOP**
  - If "回 Dashboard" → invoke `/team` → **STOP**
- If "取消" → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

**If worktree ENABLED:** Allow multiple active Contracts. Each mission gets its own worktree directory, so parallel work is safe.
- If any Contract exists → display it as info: "Active mission(s): **#{issue}** — {title}. Worktree mode: parallel claiming allowed."
- Continue to Step 2 (do NOT stop).

---

## Step 2: Route by Argument

### If `list`:

```bash
gh issue list --label "$MISSION_LABEL" --state open --assignee "$GH_USER" --json number,title,labels,milestone --limit 20
```

Show Issues assigned to the current user.

Format as numbered list:

**Your assigned missions:**
  🟠 `1.` **#042** Add user authentication
  🟠 `2.` **#045** Add rate limiting
  🟡 `3.` **#047** Write API docs

If none:
  - If team roster has only 1 member (solo mode) → "No missions assigned. Create one with `/team-issue` or check `/team queue`." → **STOP**
  - If team roster has multiple members → "No missions assigned to you. Ask your team lead to assign one via `/team-issue`, or check `/team queue`." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

Use `AskUserQuestion` to let user pick one, then proceed to Step 3 with the selected Issue number.

### If `#N` or number:

Extract Issue number. Proceed to Step 3.

### If empty:

```bash
gh issue list --label "$MISSION_LABEL" --state open --assignee "$GH_USER" --json number,title,labels --limit 10
```

Show the user's assigned Issues sorted by priority (P0 first, then P1, P2, P3).
If none:
  - If team roster has only 1 member (solo mode) → "No missions assigned. Create one with `/team-issue` or check `/team queue`." → **STOP**
  - If team roster has multiple members → "No missions assigned to you. Ask your team lead to assign one via `/team-issue`, or check `/team queue`." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

Use `AskUserQuestion` to let user confirm or pick one (even if only one issue — always confirm before claiming). Then proceed to Step 3.

---

## Step 3: Fetch Issue Details

```bash
gh issue view {issue} --json number,title,body,labels,milestone,assignees,url,comments
```

Parse the Issue body. If the Issue was created with the mission template, extract structured fields:
- **Priority** from labels (P0/P1/P2/P3)
- **Objective** from `### Objective` section
- **Sub-tasks** from `### Sub-tasks` section (checkboxes)
- **Success Criteria** from `### Success Criteria` section (also accept `### Acceptance Criteria` for backward compat)
- **Context** from `### Context & References` section (also accept `### Context`)
- **Test Command** from `### Verification Method` section or config `project.test_command`

Also extract from the JSON response:
- **Milestone** — read from `milestone.title` field

```bash
# ISSUE_DATA is the full JSON from `gh issue view ... --json ...` above
ISSUE_MILESTONE=$(echo "$ISSUE_DATA" | jq -r '.milestone.title // empty' 2>/dev/null)
# Fall back to config versions.current if Issue has no milestone
if [ -z "$ISSUE_MILESTONE" ]; then
  ISSUE_MILESTONE=$(bash ~/.claude/commands/scripts/tw-config.sh versions.current "" 2>/dev/null)
fi
```

If the Issue body is freeform (not from template), use AI understanding to extract:
- Objective: summarize what needs to be done
- Sub-tasks: break down into checkable items
- Success criteria: infer from the description

---

## Step 4: Generate Mission Contract

Create `$TEAMWORK_DIR/active/MISSION-{issue}.md` with this structure:

```markdown
---
issue: {issue}
url: {url}
title: "{title}"
assignee: {user}
priority: {priority}
labels: [{labels}]
branch: mission/{issue}-{slug}-{user}
milestone: "{milestone — from Issue JSON, falls back to config versions.current, or 'none'}"
claimed: {ISO_TIMESTAMP}
issue_content_hash: "$(bash ~/.claude/commands/scripts/tw-contract.sh hash "$TITLE" "$BODY")"
---

# MISSION-{issue}: {title}

## Objective
{Extracted from Issue body}

## Sub-tasks
{Extracted checkboxes from Issue body, or AI-generated breakdown}

**MANUAL task pre-tagging:** While extracting sub-tasks, identify any that require human action outside the codebase (server access, console operations, third-party configuration, deployment verification). Prepend `🔧 MANUAL:` to these tasks. This enables `/team-drive` to skip them during execution and surface them in the Manual Ops Handoff.

Example:
- [ ] Implement API authentication
- [ ] 🔧 MANUAL: Configure Firebase Console authorized domains
- [ ] 🔧 MANUAL: Deploy SA JSON to production server

## Success Criteria
{Extracted from Issue body, or AI-inferred}

## Context Files
{AI-generated: scan the project structure and list files relevant to this mission}
{Use Glob and Grep to find related code}

## Test Command
{From Issue body, or from config.yml project.test_command}

## AI Notes
(populated during /team-drive execution)
```

**Standalone Test**: Re-read the generated Contract as if you never saw the Issue. Can someone who has never seen the Issue execute from this Contract alone? If not, add what's missing.

**The Context Files section is the key AI enhancement.** This is what makes the Contract more valuable than the raw Issue body. Scan the project to identify:
- Files that will likely need modification
- Related test files
- Configuration files that might be affected
- Documentation that should be updated

Use `Glob` and `Grep` with keywords from the Issue title and objective to discover relevant files.

---

## Step 4b: Incorporate Issue Comments

From the `comments` field fetched in Step 3, filter using the script:

```bash
FILTERED_COMMENTS=$(bash ~/.claude/commands/scripts/tw-contract.sh read-comments {issue} --exclude-bot --exclude-lifecycle 2>/dev/null) || true
```

If `$FILTERED_COMMENTS` is empty or `[]` → skip Pre-Claim Discussion section entirely, but still store timestamp (Fix 4 below).

If comments exist after filtering:
- Add to Contract after the AI Notes section:

```markdown
## Pre-Claim Discussion ({count} comments)
- @{author} ({date}): {body}
- @{author2} ({date}): {body}
```

- The AI MUST check if any comment contradicts or extends the Issue body's requirements
- If contradiction found → adjust Objective/Sub-tasks/Criteria accordingly
- Note the adjustment in AI Notes: "Adjusted per @{author}'s comment on {date}"

After Contract creation, store the comment timestamp:

```bash
bash ~/.claude/commands/scripts/tw-contract.sh store-comment-timestamp "$TEAMWORK_DIR/active/MISSION-{issue}.md" "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
```

---

## Step 4c: Confirm Contract

Display the generated Contract summary (Objective, Sub-tasks, Success Criteria) and confirm before creating branch. State your position as AGREE+EXTEND, CHALLENGE+PROPOSE, or REDIRECT. If proposing an alternative, state what it costs. Never silent compliance.

**Detail ≠ accuracy.** Before AGREE, try to falsify each concrete claim — especially factual assertions about system state (what's broken, what's missing, which environments are affected). A confident-sounding Issue with specific IPs and step counts can still be partially wrong. If you didn't independently verify a claim, you cannot AGREE with it.

Use `AskUserQuestion`:
```
question: "Contract ready. Confirm to create branch?"
options:
  - label: "Confirm"
    description: "Create branch and start mission"
  - label: "Edit"
    description: "Modify Contract before proceeding"
  - label: "Cancel"
    description: "Don't claim"
```

- If "Edit" → apply changes, re-preview (loop)
- If "Cancel" → delete Contract file → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

---

## Step 5: Create Branch (+ optional worktree)

**⚠️ IMPORTANT: Use `claim-setup` — it handles worktree routing internally.**
The script reads `git config --local teamwork.worktree` and automatically chooses the correct path (branch-only vs worktree). You do NOT need to check worktree mode yourself.

```bash
# Slugify the title
SLUG=$(bash ~/.claude/commands/scripts/tw-git.sh slugify "$ISSUE_TITLE")

# Atomic branch+worktree setup (worktree routing is INSIDE the script)
SETUP_RESULT=$(bash ~/.claude/commands/scripts/tw-git.sh claim-setup \
  "$ISSUE_NUMBER" "$SLUG" "$GH_USER" \
  "$TEAMWORK_DIR" "$TEAMWORK_DIR/active/MISSION-$ISSUE_NUMBER.md" 2>/dev/null)

# Parse structured output (KEY=VALUE lines)
eval "$SETUP_RESULT"
# Now available: $MODE (branch|worktree), $BRANCH, $STATUS (created|exists), $WORKTREE_PATH (if worktree)
```

If `claim-setup` fails (exit ≠ 0) → "**ERROR:** Branch setup failed. Run `/team doctor` to diagnose." → **STOP**

The script handles ALL of this internally:
- **Non-worktree**: `ensure-base` + `create-branch` (or checkout if exists)
- **Worktree**: fetch base, resolve worktree path from `teamwork.worktree-root`, create worktree, copy Contract + config, write `.mission` marker
- **Branch already has worktree**: returns `STATUS=exists` + `WORKTREE_PATH` (no new worktree created)

---

## Step 6: Post Claim Comment

In the push model, Issues are already assigned and labeled `status:wip` during `/team-issue` creation. No assignment or label transition needed here.

```bash
gh issue comment {issue} --body "🚀 Claimed by @${GH_USER} — starting work on branch \`${BRANCH}\`"
```

Non-fatal: if comment fails, warn but continue.

---

## Step 7: Output Mission Briefing

Display a formatted briefing:

**CLAIMED** ──── **#{issue}** {title} ────────────
{priority_dot}  `{priority}`  Branch: *{branch}*
Contract:  `$TEAMWORK_DIR/active/MISSION-{issue}.md`
{If worktree:} Worktree:  `{path}`

**OBJECTIVE**
  {objective summary}

**SUB-TASKS**
  [ ] {task 1}
  [ ] {task 2}
  [ ] ...

**SUCCESS CRITERIA**
  {criteria}

**CONTEXT FILES**
  {list of relevant files}

{Use $MODE from claim-setup result to determine output:}

{If $MODE = "worktree":}
⚠️ **DO NOT run `/team-drive` in this terminal**
────────────────────────────────────────────
当前终端在主仓库。Worktree 模式下，所有操作必须在 worktree 目录中执行。
在这里运行 `/team-drive` 会操作错误的代码。

**打开新终端 tab，执行：**
  `cd {$WORKTREE_PATH} && claude`
  Then: `/team-drive` to execute │ `/team-ship` when done
────────────────────────────────────────────

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

{If $MODE = "branch":}

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

Use `AskUserQuestion`:
```
question: "Contract 已生成，开始执行？"
options:
  - label: "开始 drive"
    description: "/team-drive — 执行 mission"
  - label: "先看 Contract"
    description: "打开 Contract 文件确认内容"
  - label: "回 dashboard"
    description: "/team — 总览"
```
- If "开始 drive" → invoke `/team-drive` → **STOP**
- If "先看 Contract" → read and display Contract file content → **STOP**
- If "回 dashboard" → invoke `/team` → **STOP**

---

## Error Handling

- Issue not found → "**ERROR:** Issue **#{issue}** not found. Check the number." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- Issue already assigned to someone else → warn but allow claiming (team member may be handing off)
- Issue is closed → "**ERROR:** Issue **#{issue}** is already closed. Reopen it: `gh issue reopen {issue}`, or check `/team queue` for open issues." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- Branch already exists (claim-setup returns `STATUS=exists`):
  - Check if corresponding Issue is still OPEN (via `gh issue view {issue} --json state --jq '.state'`)
  - If Issue is CLOSED → warn: "⚠ Branch exists but Issue **#{issue}** is closed. This is an orphan branch. Run `/team doctor fix` to clean up, then `/team-claim` again." → **STOP**
    💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
  - If Issue is OPEN → `claim-setup` already handled correctly (checkout or worktree reuse). Proceed to Step 6.
- Network errors → "**ERROR:** GitHub API error. Check your connection and `gh auth status`." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

**On any STOP:** Always append: 💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
