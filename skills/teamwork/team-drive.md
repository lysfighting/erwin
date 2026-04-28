---
description: "Execute mission from Contract. Try: /team-drive help"
version: "3.18.0"
---

# /team-drive — Execute Mission

> You own this mission. Every error narrows the problem space. The loop does not break. When you see a better approach than the Contract suggests, implement it — silence is complicity.

> Read your Mission Contract, validate the plan against actual code, then execute — sequentially for small missions, with team + wave parallelism for large ones.

**Visual Encoding** (apply to ALL output — see `docs/visual-encoding-standard.md`):
`**bold**` → headers/labels (white) · `` `backtick` `` → commands/paths/counts (purple-blue) · `*italic*` → branches (dim) · `**#NNN**` → issues (light-blue clickable, 3+ digits) · `────` dividers · ⛔ NO code blocks around output

**Design Constitution** (read before modifying this skill — see `docs/teamwork-design-constitution.md`)

**User input**: $ARGUMENTS

If `$ARGUMENTS` is `help` or `-h`, output the following and **STOP**:

**`/team-drive`** — Execute your claimed mission

**USAGE**
  `/team-drive`           Read Contract, validate plan, execute with verify loop

**WHAT HAPPENS**
  1. Reads your Mission Contract (`.teamwork/active/MISSION-N.md`)
  2. Reads actual code — validates sub-tasks against ground truth
     (sub-tasks are hypotheses, not orders — executor rewrites if wrong)
  3. Classifies tasks: code tasks (AI executes) vs 🔧 MANUAL (human action)
  4. Selects execution mode (based on code task count):
     — Small (`1-3` tasks): sequential loop
     — Medium (`4-7` tasks): create 1-2 teammates, wave parallelism
     — Large (`8+` tasks): full team + STL hierarchy + wave map
  5. Executes code tasks with discipline: ripple check, self-adversarial review
  6. Manual Ops Handoff: surfaces 🔧 MANUAL tasks for you to complete
  7. Progress saved via checkboxes — can be interrupted and resumed

**TASK SYMBOLS**
  ○  Pending code task — AI will execute
  🔧 MANUAL task — requires your action (AI skips)
  ✅ Completed task

**PREREQUISITES**
  Run `/team-claim` first to generate a Contract.

**NEXT**: `/team-ship` to deliver via PR

---

# ────────────────────────────────────────────
# PART I: MISSION LOADING (Teamwork Layer)
# ────────────────────────────────────────────

## Step 0: Find Active Contract

```bash
# Detect worktree mode — .mission file contains Issue number
if [ -f .mission ]; then
  WORKTREE_MODE=true
  MISSION_ISSUE=$(cat .mission)
fi

# Detect config directory (in worktree, may need to check parent repo)
TEAMWORK_DIR=$(bash ~/.claude/commands/scripts/tw-config.sh detect-dir 2>/dev/null) || TEAMWORK_DIR="NO_CONFIG"
```

- If no config → "**ERROR:** Teamwork not initialized. Run `/team` first." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- If in worktree mode (`WORKTREE_MODE=true`) → no branch switching needed in Step 1.
- **Worktree note**: Contract should be in `$TEAMWORK_DIR/active/` within the worktree (copied during `/team-claim`). If not found but `.mission` exists, look for the Contract in the main repo parent directory as fallback:

```bash
if [ "$WORKTREE_MODE" = true ] && ! ls $TEAMWORK_DIR/active/MISSION-*.md >/dev/null 2>&1; then
  PARENT_REPO=$(git rev-parse --path-format=absolute --git-common-dir 2>/dev/null | sed 's|/\.git$||')
  if [ -n "$PARENT_REPO" ] && [ -d "$PARENT_REPO/.teamwork/active" ]; then
    PARENT_CONTRACT=$(ls "$PARENT_REPO/.teamwork/active/MISSION-${MISSION_ISSUE}.md" 2>/dev/null || true)
    if [ -n "$PARENT_CONTRACT" ]; then
      cp "$PARENT_CONTRACT" "$TEAMWORK_DIR/active/"
      echo "Recovered Contract from main repo: $PARENT_CONTRACT"
    fi
  fi
fi
```

```bash
ls $TEAMWORK_DIR/active/MISSION-*.md 2>/dev/null
```

- If no Contract found → "No active mission. Run `/team-claim` first." Check for recoverable state before stopping:

```bash
# Check if there's a mission branch with an open Issue (recoverable state)
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" == mission/* ]]; then
  # Extract issue number from branch name (handles mission/42-*, bugfix/T-044-*, etc.)
  ISSUE_NUM=$(bash ~/.claude/commands/scripts/tw-git.sh extract-issue-number "$CURRENT_BRANCH" 2>/dev/null) || true
  if [ -n "$ISSUE_NUM" ]; then
    ISSUE_STATE=$(gh issue view "$ISSUE_NUM" --json state --jq '.state' 2>/dev/null)
    if [ "$ISSUE_STATE" = "OPEN" ]; then
      echo "Found open Issue #$ISSUE_NUM for this branch. Contract may have been deleted."
      echo "Run /team-claim #$ISSUE_NUM to regenerate the Contract, then /team-drive again."
    fi
  fi
fi
```

→ **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

- If multiple Contracts found → "Multiple active contracts found. This shouldn't happen. Keep one, remove the rest." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

Read the Contract file fully. Extract from YAML frontmatter:
- `issue` number
- `title`
- `branch`
- `priority`

Extract from body:
- **Objective**
- **Sub-tasks** (with checkbox status)
- **Success Criteria**
- **Context Files**
- **Test Command**

---

## Step 0b: Issue Freshness Check

Check whether the Issue has been modified since the Contract was generated.

```bash
CONTRACT_PATH="$TEAMWORK_DIR/active/MISSION-{issue}.md"
FRESHNESS=$(bash ~/.claude/commands/scripts/tw-contract.sh check-freshness "$CONTRACT_PATH" {issue} 2>/dev/null) || true
```

- "FRESH" (exit 0) → Issue unchanged, continue
- "NETWORK_ERROR" (exit 4) → warn "Could not check Issue freshness (network error). Continuing with existing Contract." → **continue** (non-fatal)
- "NO_HASH" (exit 2) → legacy Contract without hash, skip check, continue
- "STALE" (exit 1) → Issue modified since claim:

If `FRESHNESS` is `STALE`:

Display the current Issue body to the user:
**⚡ ISSUE UPDATED** since claim
────────────────────────────────────────────
The Issue description has changed since you generated this Contract.

**Current Issue body**
────────────────────────────────────────────
{current Issue body}
────────────────────────────────────────────

Use `AskUserQuestion`:
- "Update Contract (recommended)" → Re-extract Objective, Sub-tasks, Success Criteria from new body. Update `issue_content_hash`. Preserve Context Files and AI Notes (locally generated).
- "Continue with old Contract" → proceed without changes. ⚠️ Risk: sub-tasks may be outdated — could do redundant or wrong work.
- "Abort" → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

---

## Step 0b1.5: Comment Check

```bash
LAST_COMMENT_CHECK=$(bash ~/.claude/commands/scripts/tw-contract.sh read-field "$CONTRACT_PATH" last_comment_check 2>/dev/null) || true
COMMENT_STATUS=$(bash ~/.claude/commands/scripts/tw-contract.sh check-comments {issue} "$LAST_COMMENT_CHECK" 2>/dev/null) || true
```

- If result is "NETWORK_ERROR" → skip silently, continue
- If result is "UNCHANGED" → continue (no display)
- If result starts with "NEW" (includes "NEW" and "NEW:unknown"):

Fetch full comments:
```bash
COMMENTS=$(bash ~/.claude/commands/scripts/tw-contract.sh read-comments {issue} --exclude-bot --exclude-lifecycle --since "$LAST_COMMENT_CHECK" 2>/dev/null)
```

Parse `$COMMENTS` as a JSON array. For each entry, extract `author` (string), `createdAt` (ISO timestamp → convert to relative time), and `body` (full text).

Display:

⚡ **NEW COMMENTS** on **#{issue}** since last check
────────────────────────────────────────────
  `@{author}` — *{relative_time}*
  {full body}
────────────────────────────────────────────

Then AI makes judgment call:
- If any comment suggests direction change (contains: "change", "actually", "instead", "don't", "cancel", "blocker", "wait", "hold") OR has `authorAssociation` of `MEMBER`, `OWNER`, or `COLLABORATOR`:
  → AskUserQuestion: "Issue has new comments that may affect mission direction."
    Options: "Update Contract" / "Acknowledged, continue" / "Abort mission"
- Otherwise: display and continue (advisory)

Update timestamp (only if comments were successfully checked):
- If `COMMENT_STATUS` started with "NEW" or was "UNCHANGED" → update timestamp
- If `COMMENT_STATUS` was "NETWORK_ERROR" → do NOT update (preserves last known good timestamp)
```bash
bash ~/.claude/commands/scripts/tw-contract.sh store-comment-timestamp "$CONTRACT_PATH" "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
```

---

## Step 0b2: Cross-MC Knowledge Intake

Check if sibling MCs have shared knowledge to this issue via comments.

```bash
# Fetch issue comments (lightweight — only metadata + body)
COMMENTS_JSON=$(gh issue view {issue} --json comments --jq '.comments[] | select(.body | startswith("🔗 **Cross-MC**")) | {author: .author.login, body: .body, createdAt: .createdAt}' 2>/dev/null) || true
```

- If no `🔗 Cross-MC` comments found → skip, continue to Step 0c
- If found → display them as advisory context:

**📡 Cross-MC Knowledge** ({count} from sibling MCs)
────────────────────────────────────────────
{For each cross-MC comment:}
{comment body — already structured with From/Context/Finding/Confidence}
────────────────────────────────────────────

Inject into Contract context as advisory (append to AI Notes section if updating Contract):

```
## Cross-MC Advisory (read-only)
{comment bodies, prefixed with source issue number}
```

**Rules:**
- Cross-MC knowledge is **advisory only** — it informs but does NOT override this MC's own Contract objectives or sub-tasks
- **⛔ No transitive propagation**: Do NOT re-share knowledge received from other MCs via `/team share`. Only share first-hand discoveries from THIS mission.
- If cross-MC knowledge contradicts this MC's approach, flag it for the user but do not auto-resolve

---

## Step 0c: Worktree Guard

**Use `worktree-guard` script** — it detects mismatch atomically (no manual config check needed):

```bash
GUARD_RESULT=$(bash ~/.claude/commands/scripts/tw-git.sh worktree-guard "$CONTRACT_BRANCH" 2>&1) || GUARD_EXIT=$?
eval "$GUARD_RESULT"
# Now available: $WORKTREE_ENABLED, $IN_WORKTREE, $MISMATCH, $BRANCH, $WORKTREE_PATH
```

If `MISMATCH=true` (exit 3): worktree enabled but running from main repo.

If `$WORKTREE_PATH` is non-empty (worktree exists for this branch):

Output the mismatch warning, then use `AskUserQuestion` to let the user choose:

**⚠️ WORKTREE MISMATCH** — mission worktree: `{worktree_path}`

```
question: "Worktree mode 已启用，但你在主仓库。Mission #{issue} 的 worktree 在 {worktree_path}。"
options:
  - label: "新开 tab 切到 worktree（推荐）"
    description: "打开新终端 tab，cd {worktree_path} && claude，然后 /team-drive"
  - label: "在主仓库继续"
    description: "忽略 worktree，直接在当前目录执行（本次）"
  - label: "关闭 worktree 模式"
    description: "git config --local teamwork.worktree false，然后继续"
```

- If "新开 tab 切到 worktree" → output: "**打开新终端 tab，执行：** `cd {worktree_path} && claude` → `/team-drive`" → **STOP**
- If "在主仓库继续" → proceed to Step 1 (skip worktree check for this run)
- If "关闭 worktree 模式" → `git config --local teamwork.worktree false` → proceed to Step 1

If no worktree exists for this branch (e.g., worktree was deleted or claim was done before enabling worktree mode):

**⚠️ NO WORKTREE** — branch *{branch}* has no worktree

```
question: "Worktree mode 已启用，但 branch {branch} 没有对应的 worktree。"
options:
  - label: "在主仓库继续"
    description: "忽略 worktree，直接在当前目录执行（本次）"
  - label: "重新 claim（创建 worktree）"
    description: "删除 Contract，重新 /team-claim #{issue}"
  - label: "关闭 worktree 模式"
    description: "git config --local teamwork.worktree false，然后继续"
```

- If "在主仓库继续" → proceed to Step 1
- If "重新 claim" → delete Contract, output `/team-claim` **#{issue}** → **STOP**
- If "关闭 worktree 模式" → `git config --local teamwork.worktree false` → proceed to Step 1

---

## Step 1: Verify Branch

**If worktree mode** (`WORKTREE_MODE=true`): skip branch check — worktree is already on the correct branch.

**Otherwise:**
```bash
CURRENT_BRANCH=$(git branch --show-current)
```

Compare with Contract's `branch` field.
- If on wrong branch → `git checkout {contract.branch}`
- If branch doesn't exist locally → check remote and offer interactive recovery:

  ```bash
  REMOTE_EXISTS=$(git ls-remote --heads origin "{contract.branch}" 2>/dev/null | wc -l | tr -d ' ')
  ```

  Use `AskUserQuestion`:
  ```
  question: "Branch *{branch}* 不在本地。如何恢复？"
  options:
    {If REMOTE_EXISTS > 0:}
    - label: "从 remote 恢复"
      description: "git checkout -b {branch} origin/{branch}"
    - label: "重新 claim"
      description: "重新生成 Contract（remote 上已 push 的代码不受影响）"
    - label: "取消"
      description: "不操作"
  ```

  - If "从 remote 恢复" → `git checkout -b {branch} origin/{branch}` → continue
  - If "重新 claim" → delete Contract → output "Run `/team-claim` **#{issue}** to re-claim." → **STOP**
  - If "取消" → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

---

## Step 2: Display Mission Briefing

Output a formatted briefing:

**🎯 MISSION** ── **#{issue}** {title} ─────────────
{priority_dot} {priority}  🔀 *{branch}*
Progress:  {progress_bar}  `{done}/{total}`

**🎯 OBJECTIVE**
  {objective}

**📋 SUB-TASKS**
  ✅ {completed task}
  ▸ {remaining task}  ← current
  ○ {remaining task}

**✅ SUCCESS CRITERIA**
  {criteria}

**📁 CONTEXT FILES**
  {files list}
────────────────────────────────────────────

Check whether all **code** sub-tasks are already checked (ignore `🔧 MANUAL` tasks).

- If all code tasks checked AND **no** MANUAL tasks exist:
  ```bash
  # Check if a PR already exists for this branch
  CURRENT_BRANCH=$(git branch --show-current)
  EXISTING_PR=$(bash ~/.claude/commands/scripts/tw-pr.sh exists "$CURRENT_BRANCH" 2>/dev/null) || true
  ```
  - If `$EXISTING_PR` is non-empty → "All sub-tasks complete. PR **#{pr}** already exists — push to update it." → **STOP**
  - Otherwise → "All sub-tasks complete." then **AskUserQuestion**:
    ```
    question: "所有代码任务已完成，下一步？"
    options:
      - label: "Ship"
        description: "/team-ship — 创建 PR"
      - label: "自查一遍"
        description: "Review 代码变更再 ship"
      - label: "回 dashboard"
        description: "/team — 总览"
    ```
    Route: "Ship" → `/team-ship` | "自查一遍" → show `git diff` summary, then re-ask | "回 dashboard" → `/team`
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

- If all code tasks checked AND MANUAL tasks **do** exist:
  → If Plan Assessment was already recorded in AI Notes (resume): skip to **Manual Ops Handoff** directly.
  → If first execution (no Plan Assessment in AI Notes): continue to **Plan Validation** first (MANUAL tasks still require challenge — "Is it necessary?" applies to all task types).
  → Then proceed to **Manual Ops Handoff** (Hard Stop for User Choice):
    1. Output **🏁 CODE COMPLETE** summary
    2. Output **🔧 YOUR TURN** block (manual task list)
    3. **AskUserQuestion** — user must choose: 引导执行 / Done / Ship first / Skip
    → **STOP** (requires user response before proceeding)
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

- If unchecked code tasks remain → continue to Plan Validation and execution.

---

## Step 2a: Branch Readiness Check

Before starting work, check branch sync status:

```bash
BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "main" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-main}"
CURRENT_BRANCH=$(git branch --show-current)

# Fetch latest remote refs first (otherwise rev-list uses stale data)
git fetch origin --quiet 2>/dev/null || true

# Check how far behind base
BEHIND_BASE=$(git rev-list --count HEAD..origin/$BASE_BRANCH 2>/dev/null || echo "0")

# Check how far behind own remote
if git rev-parse --verify "origin/$CURRENT_BRANCH" >/dev/null 2>&1; then
  BEHIND_REMOTE=$(git rev-list --count HEAD..origin/$CURRENT_BRANCH 2>/dev/null || echo "0")
else
  BEHIND_REMOTE="0"
fi
```

Display branch readiness:

- If `BEHIND_BASE` > 30 → 🔴 **Branch is `{N}` commits behind base.** Use `AskUserQuestion`:
  ```
  question: "Branch is {N} commits behind base. Syncing is strongly recommended before starting work."
  options:
    - label: "Sync now"
      description: "Rebase on latest base branch (conflicts will be auto-resolved)"
    - label: "Continue without syncing"
      description: "Start driving without rebasing (may cause conflicts later)"
  ```
  - If "Sync now" → run `/team-ship sync` inline (fetch + rebase + force-push), then continue
  - If "Continue without syncing" → continue to Step 2b

- If `BEHIND_BASE` > 10 and ≤ 30 → 🟡 **Branch is `{N}` commits behind base.** Use `AskUserQuestion`:
  ```
  question: "Branch is {N} commits behind base. Consider syncing before driving."
  options:
    - label: "Sync now"
      description: "Rebase on latest base branch"
    - label: "Continue without syncing"
      description: "Proceed with current state"
  ```
  - If "Sync now" → run `/team-ship sync` inline, then continue
  - If "Continue without syncing" → continue

- If `BEHIND_REMOTE` > 0 → 🟡 **Local branch is `{N}` commits behind remote.** Run `/team-ship sync` to sync. (informational — do not stop)
- If all zero → no output needed (clean state)

---

## Step 2b: Branch Safety Check

Before executing any code changes, verify you are NOT on a protected branch.

```bash
bash ~/.claude/commands/scripts/tw-git.sh protect-check
```

- If exit code 3 → on protected branch, output **ERROR:** "Switch to a feature branch first: `/team-claim` **#{issue}**" → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- If exit code 0 → on correct feature branch → continue

---

# ────────────────────────────────────────────
# PART II: DRIVE EXECUTION ENGINE (Core)
# ────────────────────────────────────────────

> Everything below is the execution engine — independent of teamwork infrastructure. These are the principles, methods, and loops that turn a plan into shipped code.

---

## The Sage — "不做应声虫，做超级智者"

You are a **super-sage** with independent judgment. NOT a compliant executor.

**Six Principles (chain — each feeds the next):**

1. **Critical Thinking** — See through to essence. "What is the REAL problem? Is this the best angle?"
2. **Creative Thinking（举一反三）** — From one insight, derive many. Extend ideas beyond what the planner saw.
3. **Self-Debating → Self-Cohesive** — Attack your own proposals with full force. Only survivors proceed. Your antithesis must make you genuinely hesitate.
4. **Intellectual Honesty** — If the plan is wrong, say so with reasoning. If you discover your approach is wrong, acknowledge immediately. Truth over comfort.
5. **Simplicity** — Best solution = simplest. Complexity = unfinished thinking. Prefer deletion over addition. One general mechanism over ten special cases.
6. **知行合一** — When you change A, trace EVERY ripple to B, C, D. Tests, docs, imports, types. Half-applied insight = inconsistency.

**Anti-Compliance Core Rule:** At every decision point, independently assess the direction. Respond as AGREE+EXTEND, CHALLENGE+PROPOSE, or REDIRECT — never silent compliance. When proposing an alternative, state what it costs (trade-off honesty). Even when you agree with the Contract, articulate WHY.

**Scaling Principle:** Before any design, ask: "How does this scale?" If "add more rules" → RED FLAG (O(2^C)). Prefer one general mechanism over enumerated cases. For deep reasoning on non-trivial decisions, `Read commands/reasoning-toolkit.md`.

### Proactive Triggers — challenge yourself during execution

| Pattern | Action |
|---|---|
| Not the simplest solution | Propose simpler |
| Solving symptom, not root cause | Name root cause |
| Assumption treated as fact | "This assumes X — verified?" |
| Scaling problem (rule explosion) | Propose scalable alternative |
| Better approach Contract didn't consider | Implement it (note in AI Notes) |
| Mediocre consensus | "Settling. Best we can do?" |
| Own earlier approach was wrong | Self-correct immediately |
| Decision made by default | "Implicitly deciding X. Explicit choice?" |

---

## Cardinal Rules

1. **Never stop**: Only valid stops: waiting for AskUserQuestion answer, truly blocked after 5 escalation levels, or all sub-tasks complete. Never output summaries and wait. Never "should I continue?"
2. **No gaming**: Don't disable tests, weaken assertions, hardcode outputs, suppress errors. If verification fails, fix the root cause.
3. **File Supremacy**: The Contract on disk always wins over context memory. After compaction or when uncertain → Re-read from disk.
4. **Courage to delete**: Replaced code → DELETE it. Not comment out. Not `// removed`. Git remembers.

---

## Plan Validation — "未经审视的计划不值得执行"

> The Contract's sub-tasks are the issue creator's HYPOTHESIS about how to achieve the Objective — written with limited codebase knowledge. You, the executor, now have ground truth. **Objective + Success Criteria are immutable orders (WHAT). Sub-tasks are a suggested route (HOW) — challenge, revise, or confirm them.** If you discover the Objective itself needs changing, stop and AskUserQuestion — never silently amend the WHAT. Append the amendment to AI Notes with rationale.

### Read the terrain

Read ALL files listed in **Context Files**. Use `Glob` and `Grep` to explore beyond what's listed — the Contract's file list was generated at claim time with surface-level scanning. Discover:
- Files that will actually need modification
- Dependencies, callers, interfaces that the planner couldn't see
- Existing patterns the implementation must follow
- Test infrastructure relevant to this mission

**Pantheon consultation**: If the mission involves architecture, module boundaries, data flow, API design, or scaling decisions — read `.claude/pantheon/` for relevant thinkers' methods. Apply them to challenge the planner's approach:
- 孙子 for strategic prioritization (are we attacking the right target?)
- 费曼 for simplification (is the planner overcomplicating this?)
- 冯·诺依曼 for separation of mechanism vs content
- 波普尔 for falsification (can we disprove the sub-tasks' assumptions?)

This is not decoration — it prevents defaulting to the first approach that "seems reasonable."

### Challenge each sub-task — ALL types (code + MANUAL)

MANUAL is a classification for execution, not an exemption from challenge. For each unchecked sub-task (code AND manual), ask:

1. **Is it necessary?** Does this sub-task actually advance the Objective, or is it busywork / a wrong assumption about the codebase?
2. **Is it correct?** Given what the code actually looks like, is this the right approach? Or did the planner assume a structure that doesn't exist?
3. **Is it sufficient?** Are there missing steps that the planner couldn't have known about? (e.g., a migration is needed, a shared interface must be updated, a config change is required)
4. **Is the order right?** Are there dependency constraints the planner missed?
5. **Opposite test:** What happens if we skip this entirely? If nothing breaks, it's probably unnecessary.

### Output Plan Assessment

**📊 PLAN ASSESSMENT**
────────────────────────────────────────────
**🎯 Objective:** {restate in own words — proves understanding}
**✅ Success Criteria:** `{N}` criteria — all achievable: {yes/no}

**Sub-task review:**
  ✅ [1] {task} — **CONFIRM:** {why it's correct}
  ✏️ [2] {task} — **REVISE:** {what's wrong, what it should be}
  ➕ [3] (missing) — **ADD:** {what's needed that planner missed}
  ➖ [4] {task} — **DROP:** {why it's unnecessary}
  🔧 [5] {task} — **MANUAL:** {requires human action — cannot be automated}
────────────────────────────────────────────

**MANUAL classification — executability check**: Before tagging any task MANUAL, run this step:
1. "Can I execute this from this terminal?" → If yes → **code task**, not MANUAL.
2. "I'm not sure — the user might have access/credentials I don't know about" → **Ask the user** via AskUserQuestion: "这个任务我可能可以帮你执行，需要什么信息/权限？" Let the user decide.
3. Only tag MANUAL when the user confirms AI cannot perform it, or the action requires physical/GUI-only interaction with no programmatic path.

### Apply revisions

- **CONFIRM**: No changes needed. Proceed.
- **REVISE**: Update the sub-task text in the Contract file. Note the rationale with evidence (file path, code pattern, or test result) in **AI Notes**.
- **ADD**: Insert new sub-tasks into the Contract. Sync new checkboxes to GitHub Issue body. Cite the discovery that revealed the gap.
- **DROP**: Remove the sub-task from Contract. Note in **AI Notes** why it was dropped, with evidence from the codebase.
- **MANUAL**: Tag ONLY after the executability check above confirms AI cannot perform it. Include why.
- **ESCALATE**: If the Objective itself appears wrong or impossible given ground truth → Use `AskUserQuestion` to alert the user. Do NOT proceed with a doomed plan.

After revisions, the Contract now reflects an **executor-validated plan** — grounded in actual code, not assumptions.

### Task Decomposition Rules

After validation, ensure the final task list follows these principles:
- **Risky/uncertain tasks FIRST** — surface unknowns early, not late
- **One substantive change per task** — atomic, reviewable, testable
- **Scaffold + verify pipeline first** — first task should confirm the build/test toolchain works
- **Final task = verification against Success Criteria** — the last thing you do is prove you're done

---

## Execution Mode + Team Assembly — "一个人走得快，一群人走得远"

Count remaining unchecked **code** sub-tasks after Plan Validation (exclude `🔧 MANUAL` tasks — those are not executed by AI).

### Zero code tasks (all MANUAL): Manual-Only Path

If remaining code sub-task count is 0 (all tasks are `🔧 MANUAL`), skip the Execution Loop (no code to execute). Plan Assessment has already been completed above — proceed to **Completion Summary + Debrief** (Manual Ops Handoff path).

### Small missions (1-3 sub-tasks): Sequential Mode

Execute via the per-task loop directly. No team creation. This is the fast path for focused, simple work.

### Medium missions (4-7 sub-tasks): Parallel Mode

1. Analyze dependencies between sub-tasks:
   - Which tasks are independent? (can run simultaneously)
   - Which tasks depend on others? (must wait)
2. Group into **waves** (a wave = tasks that can run in parallel):
   ```
   Wave 0: independent foundation tasks
   Wave 1: tasks that depend only on Wave 0
   Wave 2: tasks that depend on Wave 1
   ...
   ```
3. Create team:
   ```
   TeamCreate: team_name: "mission-{issue}", agent_type: "team-lead"
   ```
4. Spawn 1-2 teammates. **ALL teammates MUST use `model: "opus"`.** No exceptions.
   - `mode: "bypassPermissions"` — teammates don't need AskUserQuestion
   - `isolation: "worktree"` — when teammates modify same files as you

   **Leaf implementer prompt:**
   ```
   You are an implementer on team "mission-{issue}".
   Mission: #{issue} — {title}.
   YOUR tasks: {list with descriptions and target file paths}.
   Work: read code → implement → verify (run tests) → commit → report done to "team-lead".
   Source of truth: $TEAMWORK_DIR/active/MISSION-{issue}.md
   Files always win over memory. Re-read after any compaction.
   ```
5. Assign wave tasks: you take critical-path tasks, teammates take parallel tasks.

### Large missions (8+ sub-tasks): Full Team Mode

1. Full wave decomposition with critical path analysis
2. Create team with STL (Sub-Team-Lead) hierarchy:
   - Group sub-tasks by domain (same files / same subsystem)
   - Domain with 3+ tasks → spawn an STL who manages their own inner subagents
   - Domain with 1-2 tasks → spawn a leaf implementer
3. Spawn all teammates in parallel (single message). All `model: "opus"`, `mode: "bypassPermissions"`.

   **STL prompt:**
   ```
   You are STL for {Domain} on team "mission-{issue}".
   Mission: #{issue} — {title}. YOUR domain tasks: {list with descriptions}.
   You OWN this domain: analyze → spawn 2-4 Task subagents → review output → report to "team-lead".
   Escalate ONLY: cross-domain conflicts, ambiguous requirements, blockers.
   Source of truth: $TEAMWORK_DIR/active/MISSION-{issue}.md
   Files always win over memory. Re-read after any compaction.
   ```

4. You orchestrate: assign waves, unblock teammates, review output, take critical-path tasks

### Sizing by Parallelism Width

`team_size = min(max_parallel_tasks_in_widest_wave, 5)`

| Scale | Topology | Effective Streams |
|---|---|---|
| 1-3 tasks | Solo (you) | 1 |
| 4-7 tasks, 1 domain | Flat: you + 1-2 leaf | 2-3 |
| 8-12 tasks, 2 domains | 1-2 STLs + 1 leaf | 5-8 |
| 13+ tasks, 3+ domains | 2-3 STLs (you = pure orchestrator) | 10-15 |

### Teammate Isolation Principle

Brief teammates from the **Contract file ONLY**. Never from conversation history. Include file paths in every teammate prompt. Teammates have zero discussion context — this is their superpower against context pollution.

### Emit Topology (mandatory for Parallel + Team modes)

After team creation, output the team structure. Re-emit after any scaling change.

**🗺️ TEAM TOPOLOGY:** mission-{issue}
YOU (Team Lead) ─── orchestrating + critical path
├── 🎖️ stl-{domain} ─── [`{N}` tasks, spawns `{M}` inner]
├── 🔧 implementer-1 ─── [`{N}` tasks]
└── 🔧 implementer-2 ─── [`{N}` tasks]
Active: `{N}` teammates + ~`{M}` inner = `{total}` streams

---

## Wave Decomposition + Launch

Think in WAVES, not lists. Wave = tasks that run simultaneously.

```
Zero-dependency tasks → Wave 0
Depend only on Wave 0 → Wave 1
Continue until all assigned
```

### Wave Map (mandatory for Parallel + Team modes)

**📊 WAVE MAP** — Mission **#{issue}**
────────────────────────────────────────────
⚡ Wave 0 (Foundation) — Width: `{N}`
├── [T1] {task} → {you/teammate}
└── [T2] {task} → {you/teammate}
⚡ Wave 1 (Core) — Width: `{N}`  ← PEAK
├── [T3] {task} → {you/teammate}
└── [T4] {task} → {you/teammate}
🔗 Critical Path: T1 → T3 → T5
Speedup: `{X}x` via parallelism
────────────────────────────────────────────

### Critical Path Optimization

- Split large critical-path tasks into parallel chunks
- Pipeline overlap: start Wave N+1 research during Wave N
- Team Lead takes critical-path tasks (most important work = your work)
- Off-critical-path tasks have slack — use for rebalancing

### Launch Wave 0

In ONE response: assign tasks → kick-off messages → broadcast launch → start your task.

**🚀 WAVE 0 LAUNCH** — Width: `{N}`
  [{agent}] → {task}
  🔗 Critical path: {which}

---

## The Execution Loop

Execute **code** sub-tasks according to the selected mode. **Skip `🔧 MANUAL` tasks** — they are not executable by AI and will be surfaced to the user during Manual Ops Handoff at completion. For **Sequential Mode**, iterate code tasks one by one. For **Parallel/Team Mode**, use the Swarm Orchestration Loop.

### The Loop — Swarm Orchestration (Parallel + Team modes)

```
while (mission != COMPLETE) {

    // WAVE MANAGEMENT (highest priority)
    → All wave-N tasks dispatched? If not → dispatch NOW
    → Critical-path complete? → Pipeline wave N+1
    → All wave-N complete? → "WAVE N COMPLETE" → Launch N+1

    // YOUR WORK (critical-path tasks)
    execute → verify → if error: root cause, fix, verify
    → self-review → passes: broadcast → complete → next

    // SWARM ORCHESTRATION
    on STL domain_report → if complete: review + mark done
    on STL escalation → INSTANT response. Unblock NOW.
    on leaf complete → check work → toggle Contract → assign next

    // DYNAMIC HIERARCHY (every 2 waves)
    >=3 tasks in domain? → PROMOTE leaf to STL
    <=2 tasks left? → DEMOTE STL to leaf
    >6 tasks? → SPLIT STL into 2
    After changes → re-emit TOPOLOGY

    // CHECKPOINT (every 3 tasks / every wave)
    → RE-READ Contract FROM DISK (Read tool, not memory)
    → Aligned with Objective? Building what Plan Assessment said?
    → Am I solving the root cause or patching a symptom? Is there a simpler approach I missed?
    → If drifted → correct course immediately
    → Append key decisions/discoveries to AI Notes incrementally (don't wait until completion)
}
```

### Team Lead Principles

1. **Swarm throughput > personal output.** Unblocking one STL > finishing your own task.
2. **Never idle while teammates work** — review, pipeline-prep, or take the next task.
3. **STL escalations = INSTANT response.** They only escalate when truly blocked.
4. **Delegate DOMAINS to STLs, not individual tasks.** Trust tactical decisions.
5. **Don't micromanage STLs.** Intervene on cross-domain conflicts only.
6. **Scale hierarchy dynamically.** The team is living, not planned-once.

### Broadcast Protocol

Every dispatch and completion MUST be announced. No silent agents.

**📡 DISPATCH:** [{agent}] → {task} | Parallel with: [{others}]
**✅ RETURN:** [{agent}] ← {task} | Files: [{list}]
**🔄 RE-DISPATCH:** [{agent}] → {next-task}
**🎖️ STL-INNER:** [stl-{domain}] spawned `{N}` subagents | Progress: `{done}/{total}`
**🏁 DOMAIN COMPLETE:** [stl-{domain}] ← delivered
**📊 WAVE STATUS** (after each wave): Wave `{N}/{total}` | Tasks: `{done}/{total}` | Agents: `{active}`

### Checkpoint Ownership in Parallel/Team Mode

In wave-based execution, **only the Team Lead (you) writes to the Contract and syncs to GitHub.** This prevents concurrent write conflicts.

```
Sequential: you complete task → you toggle + sync → next task
Parallel:   teammate completes → reports to Lead (RETURN)
            → Lead reviews output → Lead toggles + syncs
            → wave complete → Lead launches next wave
```

Teammates NEVER directly modify the Contract file or call `sync-all-checkboxes`. They report completion via `SendMessage` to "team-lead", including:
- Which sub-task they completed
- Files modified
- Verification result (tests pass/fail)

The Lead then:
1. Reviews the teammate's work (read modified files, verify quality)
2. Toggles the sub-task checkbox in Contract (`tw-contract.sh toggle-task`)
3. After wave complete → batch-syncs to GitHub Issue (`tw-contract.sh sync-all-checkboxes`)
4. Broadcasts wave progress

---

## Per-Task Execution (applies to all modes)

### Announce current task

**Working on:** {sub-task description}

### Read context
Read the files relevant to THIS specific sub-task. Use `Glob` and `Grep` to explore beyond Context Files if needed.

### Implement
Write the code, make the changes. Follow the project's existing patterns and conventions.

### Verify

**Test Strategy** — service-aware:

Read config. If `project.services` exists (array of {name, language, test_command, lint_command}):
  1. Determine which services are affected by current changes:
     ```bash
     BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "main" 2>/dev/null)
     BASE_BRANCH="${BASE_BRANCH:-main}"
     CHANGED_FILES=$(git diff --name-only "origin/$BASE_BRANCH...HEAD")
     ```
  2. For each service in `project.services`, check if any changed file starts with the service directory path
  3. For each affected service, run its `test_command`
  4. If no service-specific command matches, fall back to `project.test_command`

If no `project.services` in config:
```bash
# From Contract's Test Command field, or fall back to config
{TEST_CMD}
```

If tests fail → fix the issue and verify again. Do not move on until verification passes.

### Ripple Check — 知行合一

After implementing, trace the ripple effects of your change:

```bash
# Find all callers/importers of modified functions/classes
# Grep for the names you changed
```

- Changed a function signature → update ALL callers
- Changed a type/interface → update ALL implementations
- Changed a config key → update ALL readers
- Added a dependency → update package manifest + lock file
- Deleted code → verify nothing still references it

**Incomplete ripple = bug factory.** Don't mark a task done until ripples are traced.

### Side Effect Check — 环境副作用记录

After ripple check, scan for **non-code changes** made during this task. If any of the following occurred, append a `[side-effect]` entry to the Contract's **AI Notes** section:

- **Config files modified**: `.env*`, `docker-compose*`, `Dockerfile*`, `Makefile`, `*.yaml`/`*.toml` config files
- **New dependencies added**: `npm install` / `pip install` / `brew install` / `apt-get install` of non-devDependencies
- **New env vars introduced**: code references `process.env.X`, `os.environ["X"]`, `os.Getenv("X")` that didn't exist before
- **Remote operations performed**: SSH to servers, API calls to external services, database DDL, cloud resource changes
- **Personal credentials used**: personal API keys, tokens, or accounts used for testing

**Format** (append to AI Notes):
```
[side-effect] {type}: {description} — {why it was needed}
```

**Examples:**
```
[side-effect] env-var: Added REDIS_CACHE_URL to .env.local — caching feature needs Redis connection
[side-effect] remote: Ran ALTER TABLE on dev DB, added column cache_ttl — new caching schema
[side-effect] credential: Used personal OpenRouter key for API testing — company needs corp account
[side-effect] dependency: Added sharp@0.34 (native binary) — image processing requires libvips
```

**Why:** These entries feed the Side Effect Triage gate in `/team-ship`. Auto-scan catches file-level changes; manual entries here catch remote operations and implicit assumptions that no diff can detect.

If no side effects occurred for this task → skip silently.

### Self-Adversarial Review (mandatory per task)

Before marking ANY task complete, attack your own work:

- Re-read every line I wrote/modified (Read tool, not memory) ✓
- Ran verification, saw it pass ✓
- Actively tried to find problems (edge cases, null inputs, race conditions) ✓
- Did NOT disable anything to make it work (tests, lint, assertions) ✓
- Confident shipping this ✓

Only after passing self-review:

### Commit + Update Contract (atomic sequence)

**⚠️ ALL THREE steps below are mandatory for EVERY sub-task. Do NOT skip steps 2-3. A committed sub-task with an unchecked Contract checkbox causes the next `/team-drive` run to waste time re-verifying completed work.**

**Step 1: Commit code**
```bash
git add {specific files changed for this sub-task}
git commit -m "{type}({scope}): {description} | Mission: #{issue}"
```

Use appropriate commit type: `feat`, `fix`, `refactor`, `test`, `docs`. **Do NOT push** — save pushes for `/team-ship`.

**Step 2: Toggle checkbox in Contract** (immediately after commit, not later)
```bash
bash ~/.claude/commands/scripts/tw-contract.sh toggle-task "$CONTRACT_PATH" {N}
```

**Step 3: Sync all checkboxes to GitHub Issue**
```bash
bash ~/.claude/commands/scripts/tw-contract.sh sync-all-checkboxes "$CONTRACT_PATH" $ISSUE_NUMBER
```

This batch-syncs ALL checked items (Sub-tasks + Success Criteria) from Contract to Issue in one API call. Uses fuzzy matching — tolerates timestamp suffixes, whitespace differences, and minor rewording. Reports unmatched items to stderr.

Non-fatal: if sync fails, warn but continue.

**Verify**: After toggling, visually confirm the checkbox is `[x]` in the Contract. If you completed multiple sub-tasks in one implementation, toggle ALL of them before moving on.

### Next task
Move to the next unchecked **code** sub-task (Sequential) or next wave task (Parallel/Team). Skip `🔧 MANUAL` tasks. In Team Mode, also check teammate status and unblock if needed.

---

## Execution Discipline — Core Principles

These principles apply to ALL execution modes. They are the difference between "code that works" and "code worth shipping."

### The Relentless Rule

If you can think of it AND it's within scope → do it NOW. Not "follow-up."

**Anti-Laziness Test for "out of scope":** ALL THREE must be true:

1. Truly requires different requirements the user hasn't given
2. Lacks technical capability or access
3. Genuinely unrelated to mission success

If any one is false → you're being lazy. Do the work.

### 举一反三 (From One, Derive Many)

Fixed a bug → same bug class elsewhere? Improved a pattern → other places benefit? Root cause → other symptoms? Don't fix one instance and leave five more.

### Deep Reasoning

For non-trivial decisions during execution: read `.claude/pantheon/` for relevant thinkers' methods. Read `commands/reasoning-toolkit.md` for structured reasoning methods (Self-Dialectic, Formal Logic, Inversion, Compression Test). Never name-drop a method without actually running it on your problem.

**When to invoke**: architecture decisions mid-implementation, unexpected complexity that suggests the approach is wrong, trade-offs where both options have real cost.

**Mandatory trigger**: When facing abstract system design decisions — architecture, module boundaries, data flow, scaling strategy, API design, state management — you **MUST** consult `.claude/pantheon/` before committing to an approach. Read the actual entries, extract the cognitive method, and run it on your problem. This is mandatory, not optional.

### File Supremacy — Re-Anchor Protocol

**The Contract on disk always wins over context memory.**

After compaction, at checkpoints, or when uncertain:
**Re-Anchor:** Read Contract FROM DISK → verify alignment
**RE-ANCHOR:** Contract OK | Progress: `{done}/{total}` | Plan Assessment: aligned

If memory says X but Contract says Y → the file wins. Always.

### Obstacles: 3 Approaches Before Escalating

When blocked:
1. Try a different approach (at least 3 alternatives)
2. Research (read docs, search codebase, WebSearch for APIs)
3. Reduce scope locally (solve a simpler version first)
4. Skip + document the blocker in AI Notes
5. Escalate to user via AskUserQuestion (LAST RESORT)

---

# ────────────────────────────────────────────
# PART III: MISSION DELIVERY (Teamwork Layer)
# ────────────────────────────────────────────

## Completion Check

### Checkbox Gate (mandatory — catches missed toggles)

**Re-read the Contract file from disk** (Read tool, not memory). Count unchecked sub-tasks (`- [ ]`).

Separate unchecked tasks into two categories:
- **Code tasks** (`- [ ]` without `🔧 MANUAL` or `(MANUAL)` tag) — these should have been completed by AI
- **Manual tasks** (`- [ ]` with `🔧 MANUAL` or `(MANUAL)` tag) — expected to be unchecked, handled in Manual Ops Handoff

For unchecked **code tasks**:
  1. For each unchecked task, verify whether the work was actually done (check git log, grep code)
  2. If work is done → toggle the checkbox NOW (`tw-contract.sh toggle-task`)
  3. If work is NOT done → execute the sub-task before proceeding
  4. **Do NOT proceed to acceptance criteria until all CODE task checkboxes are `[x]`**

Unchecked **manual tasks** are expected — they pass through to the Manual Ops Handoff in the Completion Summary.

After all Contract checkboxes are correct, batch-sync to GitHub Issue:
```bash
bash ~/.claude/commands/scripts/tw-contract.sh sync-all-checkboxes "$CONTRACT_PATH" $ISSUE_NUMBER
```
This catches any missed per-task syncs and also syncs Success Criteria checkboxes.

This gate exists because LLMs sometimes complete code work but skip the checkbox toggle step. Re-reading the Contract from disk catches this.

### Verify acceptance criteria
Go through each acceptance criterion from the Contract. For each one:
- Can you demonstrate it's met? (run a test, show output, etc.)
- **Attempt verification before punting**: If you think "this needs manual verification," first try to verify it from the terminal. If uncertain, ask the user what access is available.
- If not met → identify what's missing, implement it, commit

### Run full test suite

Apply the same service-aware test strategy as Per-Task Verify, but for ALL changes on this branch:

```bash
BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "main" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-main}"
CHANGED_FILES=$(git diff --name-only "origin/$BASE_BRANCH...HEAD")
```

All tests must pass.

### Final self-adversarial review (whole mission)

Read through ALL changes made during this mission:
```bash
bash ~/.claude/commands/scripts/tw-git.sh log-since
```

Review the ENTIRE diff, not just the last task:
```bash
git diff "origin/$BASE_BRANCH...HEAD"
```

Attack the deliverable as a whole:
- Does the sum of changes actually achieve the Objective?
- Any cross-task inconsistencies? (task 2 assumes something task 5 changed)
- Missing error handling, untested edge cases, dead code, debug artifacts?
- Security: injection, XSS, exposed secrets, unsafe operations?

Fix ALL issues found.

### Pre-Completion Check

Before declaring done, actively look for gaps:
- Edge case thought about but not handled? → Handle now.
- Test thought about but not written? → Write now.
- Code not 100% confident in? → Fix now.
- Can you think of ANY improvement within scope? → Do it NOW.

### Optional Evaluator Gate — fresh-context independent verification

**Why**: Self-Adversarial Review is still self-review — same context, same rationalizations. A fresh evaluator agent that never saw the execution thread catches the class of failure where "executor convinced themselves it's done." Ported from drive-skill v2.3.0's `agents/evaluator.md` pattern,adapted for teamwork's single-mission scope.

**Trigger decision** (in this order):

1. If Contract AI Notes has **`evaluator: required`** → **run**
2. Else if Contract AI Notes has **`evaluator: skip`** → **skip silently**
3. Else if `$ARGUMENTS` contains **`--evaluator-gate`** flag → **run**
4. Else if `$ARGUMENTS` contains **`--no-evaluator-gate`** flag → **skip**
5. Else **auto mode** (default):
   - Read Contract frontmatter `priority` and `size` fields via `tw-contract.sh read-field`
   - **Run** if `priority=P0` OR `size=XL` OR AI Notes section contains any `🔧 MANUAL` task with the word `prod`/`production`/`deploy`/`release` (production-side effect)
   - **Skip** otherwise

**If running — dispatch Explore agent with this contract**:

```
subagent_type: "Explore"
description: "Evaluator — verify mission #N"
prompt: |
  You are an independent Evaluator for mission #{ISSUE} on branch {BRANCH}.
  You have NEVER seen the execution thread. You read only these files and
  produce a fresh verdict. Do NOT trust the executor's self-assessment.

  Required reads:
  1. Contract at `.teamwork/active/MISSION-{ISSUE}.md` — Objective + Sub-tasks + Success Criteria + Non-Goals + AI Notes
  2. Full mission diff: `git diff origin/{BASE_BRANCH}...HEAD` + `git log --oneline origin/{BASE_BRANCH}..HEAD`
  3. Test output — run `{TEST_COMMAND}` fresh; do NOT trust prior "tests passed" claims

  For each Success Criterion in the Contract, answer:
    - VERIFIED ✓ / NOT VERIFIED ✗ / PARTIAL 🟡
    - Evidence: file:line citation or exact command output
    - If NOT VERIFIED: what specifically is missing?

  Also check:
    - Are there changes in the diff that are NOT covered by any Sub-task? (scope creep)
    - Are any Non-Goals violated?
    - Any `[side-effect]` AI Notes entries that actually execute irreversible changes (prod deploy, DB migration, remote write)? If so, were they intentional per Contract?
    - Obvious code quality issues (exposed secrets, debug artifacts, dead code, TODO comments, untested edge cases)?

  Output in this exact shape:

  ```
  VERDICT: PASS | FAIL | PARTIAL
  CONFIDENCE: high | medium | low

  SUCCESS CRITERIA AUDIT:
  [criterion 1]: VERIFIED | NOT VERIFIED | PARTIAL — <evidence or gap>
  ...

  SCOPE CHECK:
  - scope creep: none | <list>
  - non-goal violations: none | <list>

  SIDE-EFFECT AUDIT:
  - irreversible actions: none | <list with intentionality judgment>

  QUALITY FINDINGS:
  - <list of concrete issues, or "none found">

  RECOMMENDATION:
  <what executor should do next — empty if PASS, specific fix list if FAIL/PARTIAL>
  ```

  Be adversarial: the executor is NOT your peer — they are the subject of the audit. 
  "Looks reasonable" is NOT evidence. Every claim needs a file:line or command output.
```

**Persist verdict** to `.teamwork/active/MISSION-{ISSUE}.evaluator.md` alongside the Contract. This is the audit trail.

**Handle verdict**:

- **PASS** → proceed to `Post completion comment to Issue` below.
- **PARTIAL** (some criteria verified, some gaps) → print evaluator's RECOMMENDATION, use `AskUserQuestion`:
  ```
  question: "Evaluator found gaps. How to proceed?"
  options:
    - "Address all — return to Execution Loop with evaluator's fix list as new sub-tasks"
    - "Accept partial — proceed to Mission Delivery (document gaps in AI Notes)"
    - "Dispute — team-lead reviews verdict, may override with justification in AI Notes"
  ```
- **FAIL** → return to **Execution Loop** with the evaluator's RECOMMENDATION list as fresh sub-tasks. Re-run the loop on those. After fixing, the Evaluator Gate runs again (no loop cap — keep iterating until PASS or user overrides via the Dispute path above).

**Dispute path**: if executor + user jointly decide the evaluator is wrong (happens — evaluator doesn't always have full context), record the dispute in AI Notes with reasoning (`[evaluator-dispute] <verdict> — overridden because <reason>`). Mission can proceed. The evaluator.md file stays on disk as honest audit trail.

**Why not a fully dedicated `evaluator` subagent_type**: keeps the blast radius small — `Explore` already has fresh-context semantics and read/grep/bash. Defining a new subagent_type would require rolling out via Claude Code config, which is out of scope for this skill.

### Post completion comment to Issue

```bash
gh issue comment {issue} --body "All code tasks complete — ready for review. Branch: \`{branch}\`{if manual tasks exist: \nManual steps remaining: {list}}"
```

Non-fatal: if comment fails, warn but continue.

### AI-Suggested Issue Comment (conditional)

**Evaluate whether the execution produced insights worth sharing on the Issue.** Review the AI Notes (key decisions, issues encountered, plan deviations) and determine if any of these are true:

| Signal | Why it's worth commenting |
|--------|--------------------------|
| Fixed a bug not described in the Issue | Reviewer needs to know about extra changes |
| Implementation deviated from Issue description | Explain why a different approach was taken |
| Hit a blocker and switched approach | Decision context is valuable for the team |
| Discovered a related problem outside scope | Needs follow-up Issue |
| Performance, security, or compatibility concern | Reviewer should pay attention |
| Side effects recorded (`[side-effect]` in AI Notes) | Infra/ops team needs visibility on env changes |

**If none of these signals are present** (pure execution, no surprises) → **skip silently**. No need to comment when everything went as planned.

**If any signal is present** → generate a concise comment draft and present to the user:

```
💬 **建议追加 Issue 备注** — AI 执行过程中发现值得记录的信息：

{draft comment — 2-5 sentences, focusing on what's NOT already in the Issue body or PR body}
```

Use `AskUserQuestion`:
```
question: "追加这条备注到 Issue？"
options:
  - label: "发送"
    description: "直接发送上述备注"
  - label: "编辑后发送"
    description: "我修改一下内容"
  - label: "跳过"
    description: "不需要追加"
```

- "发送" → `gh issue comment {issue} --body "{draft}"` (non-fatal)
- "编辑后发送" → use `AskUserQuestion` with `useTextArea: true` to collect edited text → post
- "跳过" → continue

**In `/team auto` mode:** skip this step entirely (auto = zero interaction).

### Update Contract AI Notes

Add execution notes to the Contract's **AI Notes** section:
```markdown
## AI Notes
- Completed: {timestamp}
- Plan validation: {CONFIRM/REVISED — summary of changes to original sub-tasks}
- Execution mode: {sequential/parallel/team — N tasks, M waves}
- Key decisions: {any decisions made during implementation}
- Issues encountered: {any problems and how they were resolved}
- Files modified: {list of files changed}
```

### Shut down team (if created)

If teammates were spawned:
1. Verify all teammate tasks are complete
2. Send shutdown request to each teammate
3. Delete team after confirmation

---

## Completion Summary + Debrief

Scan the sub-task list for any tasks tagged `🔧 MANUAL` or containing "(MANUAL)" in the description. Split into two groups: **code tasks** (automated, done by AI) and **manual tasks** (require human action).

### If NO manual tasks exist — standard completion:

**🏁 COMPLETE** ── **#{issue}** {title} ────────────
Mode:     {sequential/parallel/team}  `{N}` tasks  `{M}` waves
Commits:  `{count}` on 🔀 *{branch}*
Tests:    ✅ passing

**📊 PLAN VALIDATION**
  {CONFIRMED / REVISED: summary}

**📋 SUB-TASKS**
  ✅ {task 1}
  ✅ {task 2}

**✅ SUCCESS CRITERIA**
  {criterion 1} — {evidence}
  {criterion 2} — {evidence}

**📊 SWARM STATS** *(if team mode)*
  Waves: `{N}`  Peak: `{N}`  Lead: `{N}`  STL: `{N}`  Leaf: `{N}`

────────────────────────────────────────────

**AskUserQuestion**:
```
question: "所有代码任务已完成，下一步？"
options:
  - label: "Ship"
    description: "/team-ship — 创建 PR"
  - label: "自查一遍"
    description: "Review 代码变更再 ship"
  - label: "回 dashboard"
    description: "/team — 总览"
```
Route: "Ship" → `/team-ship` | "自查一遍" → show `git diff` summary, then re-ask | "回 dashboard" → `/team`

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

### If manual tasks exist — Manual Ops Handoff:

The completion report becomes a two-phase handoff: first show what AI completed, then guide the user through their manual tasks.

**🏁 CODE COMPLETE** ── **#{issue}** {title} ────────────
Mode:     {sequential/parallel/team}  `{N}` code tasks  `{M}` commits
Tests:    ✅ passing

**📊 PLAN VALIDATION**
  {CONFIRMED / REVISED: summary}

**✅ CODE TASKS** (`{done count}/{code task count}`)
  ✅ {code task 1}
  ✅ {code task 2}

**🔧 YOUR TURN** — `{manual_count}` manual steps remaining
────────────────────────────────────────────

  1. {manual task description}
     → {concrete action: where to go, what to do}

  2. {manual task description}
     → {concrete action}

**⛔ DO NOT output `/team-ship` or any closing line here. You MUST invoke AskUserQuestion IMMEDIATELY as part of this output. The completion summary is NOT complete until the user responds.**

```
question: "Manual 步骤怎么处理？"
options:
  - label: "引导执行"
    description: "AI 逐步引导，每步前确认"
  - label: "Done — 我已完成"
    description: "Manual 步骤已手工完成，进入验收"
  - label: "先 ship 代码"
    description: "先创建 PR，manual 步骤写进 PR 描述"
  - label: "不需要"
    description: "跳过并记录"
```

Route responses:

- **"引导执行"** → AI walks user through each manual step interactively (see Guided Execution below)
- **"Done — 我已完成"** → proceed to show acceptance criteria + suggest `/team-ship`
- **"先 ship 代码"** → proceed to suggest `/team-ship` (note manual steps in PR description)
- **"不需要"** → note "MANUAL tasks skipped — user confirmed not applicable" in AI Notes, proceed to show acceptance criteria + suggest `/team-ship`

#### Guided Execution (引导执行)

When user selects "引导执行", AI assists with each manual step **one at a time**:

For each manual task:
1. **Assess executability** — can AI run this directly (e.g., SSH command, API call) or does it truly require human-only action (e.g., browser UI, physical access)?

2. **If AI can execute** (SSH, CLI commands, API calls):
   - Show the command(s) about to run
   - Use `AskUserQuestion`: "执行这条命令？" → "执行" / "我自己来" / "跳过"
   - If "执行" → run via Bash tool, show output, verify result
   - If verification fails → diagnose, suggest fix, retry

3. **If truly human-only** (browser console, physical hardware, third-party UI):
   - Show step-by-step instructions with concrete actions
   - Use `AskUserQuestion`: "完成了吗？" → "完成" / "遇到问题" / "跳过"
   - If "遇到问题" → troubleshoot interactively

4. **After each step** → check off the sub-task in Contract, move to next

After all steps processed (done or skipped), continue to acceptance criteria block.

**Key principle:** Many tasks labeled MANUAL are only "manual" because they require server/infra access — but if the user has SSH keys or CLI tools configured, AI can execute them directly. The label means "needs human judgment to initiate", not "AI must not touch".

After user responds, output the full acceptance criteria block and the `/team-ship` prompt:

**✅ SUCCESS CRITERIA**
  {criterion} — {evidence or *"pending manual step"*}

────────────────────────────────────────────

**AskUserQuestion**:
```
question: "所有代码任务已完成，下一步？"
options:
  - label: "Ship"
    description: "/team-ship — 创建 PR"
  - label: "自查一遍"
    description: "Review 代码变更再 ship"
  - label: "回 dashboard"
    description: "/team — 总览"
```
Route: "Ship" → `/team-ship` | "自查一遍" → show `git diff` summary, then re-ask | "回 dashboard" → `/team`

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

{If the codebase has other active Contracts in .teamwork/active/ besides the current one:}
📡 Findings may be relevant to sibling MCs. Run `/team share` **#{sibling}** to share.

### Post-Mission Ecosystem Scan (2-3 min)

Did this mission reveal a gap in any skill? Tool limitation? Stale pattern?
- Micro → fix now
- Medium → mention in debrief
- Large → flag as future mission / Issue

---

## Error Handling

- Config missing → **ERROR:** Run `/team` first to initialize teamwork. → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- Contract references files that don't exist → **WARNING:** skip those context files
- Test command not defined → **WARNING:** No test command configured. Add `test_command` to `$TEAMWORK_DIR/config.yml`.
- Git conflicts → resolve them, then continue
- If execution is interrupted (user stops mid-task), the Contract preserves progress via checkboxes — next `/team-drive` run picks up where it left off
- Plan Validation finds Objective impossible → **ERROR:** ESCALATE to user, do NOT proceed with doomed plan
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- Teammate fails or is blocked → unblock immediately (priority over your own task)

**On any STOP:** Always append: 💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
