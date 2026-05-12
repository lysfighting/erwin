---
description: "Ship mission → PR. Try: /team-ship help"
version: "3.20.0"
---

# /team-ship — Deliver Mission

> You own the delivery. Guard the quality boundary between code and production. Every detail you catch here saves a revert later.

**Visual Encoding** (apply to ALL output — see `docs/visual-encoding-standard.md`):
`**bold**` → headers/labels (white) · `` `backtick` `` → commands/paths/counts (purple-blue) · `*italic*` → branches (dim) · `**#NNN**` → issues (light-blue clickable, 3+ digits) · `────` dividers · ⛔ NO code blocks around output

**User input**: $ARGUMENTS

**Argument routing:**

| Input | Action |
|-------|--------|
| (empty) | Full ship flow (push, PR, CI, cleanup) |
| `done` | Post-merge cleanup (close Issue, update labels, clean worktree) |
| `done help` | Show done operation usage |
| `review` | AI code review on current PR |
| `review help` | Show review operation usage |
| `sync` | Rebase current branch on base branch |
| `sync help` | Show sync operation usage |
| `help` or `-h` | Show usage guide |

---

## Route by Argument

Parse `$ARGUMENTS`:
- If `help` or `-h` → output the following and **STOP**:

**`/team-ship` — Deliver your mission**

**USAGE**
  `/team-ship`            Push + create PR + CI check + cleanup
  `/team-ship done`       After PR merge: close Issue, update labels, clean worktree
  `/team-ship review`     AI review current PR (correctness/security/architecture/quality)
  `/team-ship sync`       Rebase current branch on latest base branch

**SHIP FLOW**
  1. Pre-flight: verify branch, code sub-tasks done, tests pass
     (🔧 MANUAL tasks don't block — they're listed in PR description)
  2. Push branch to GitHub
  3. Create PR with `Closes #N` (if all tasks done) or `Refs #N` (if unchecked tasks remain)
     If manual tasks exist: adds "Manual Steps (post-merge)" checklist
  4. Watch CI (if configured)
  5. Request review (if configured)
  6. Clean up local Contract, label Issue `status:review`

**AFTER MERGE**
  `/team-ship done` closes the Issue, labels `status:done`, returns to base branch.

- If `done help` → jump to **Operation Done Help**
- If `done` → jump to **Operation Done**
- If `review help` → jump to **Operation Review Help**
- If `review` → jump to **Operation Review**
- If `sync help` → jump to **Operation Sync Help**
- If `sync` → jump to **Operation Sync**
- If empty → continue to **Step 0: Full Ship Flow**
- If anything else → "Unknown subcommand: `{input}`. Run `/team-ship help` for usage." → **STOP**

---

## Step 0: Prerequisites

```bash
# Branch safety: never ship from protected branches
CURRENT=$(git branch --show-current)
bash ~/.claude/commands/scripts/tw-git.sh protect-check 2>/dev/null || {
  echo "ERROR: You're on a protected branch ($CURRENT). Switch to a mission branch first."
  # STOP — 💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
}
```

```bash
# Identity
GH_USER=$(gh api user --jq '.login' 2>/dev/null)
if [ -z "$GH_USER" ]; then
  echo "**ERROR:** Cannot get GitHub user identity. Run 'gh auth login' first."
  exit 1
fi
```

```bash
TEAMWORK_DIR=$(bash ~/.claude/commands/scripts/tw-config.sh detect-dir 2>/dev/null) || {
  echo "No teamwork config found. Run /team init first."
  # STOP — 💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
}
```

- If no config → "Teamwork not initialized. Run `/team` first." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

Read `$TEAMWORK_DIR/config.yml` → extract project settings, conventions, quality preferences.

```bash
eval "$(bash ~/.claude/commands/scripts/tw-config.sh resolve-labels 2>/dev/null)"
```

---

## Step 1: Find Active Contract

```bash
# Check for worktree mode
if [ -f .mission ]; then
  WORKTREE_MODE=true
  MISSION_ISSUE=$(cat .mission)
fi

ls $TEAMWORK_DIR/active/MISSION-*.md 2>/dev/null
```

- If no Contract found AND no `.mission` file → check for recoverable state before stopping:

```bash
# Check if current branch is a mission branch (recoverable state)
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" == mission/* ]]; then
  ISSUE_NUM=$(bash ~/.claude/commands/scripts/tw-git.sh extract-issue-number "$CURRENT_BRANCH" 2>/dev/null) || true
  if [ -n "$ISSUE_NUM" ]; then
    echo "You're on branch *${CURRENT_BRANCH}* but no Contract exists. Run /team-claim #${ISSUE_NUM} to regenerate, then /team-ship again."
  else
    echo "No active mission. Nothing to ship."
  fi
else
  echo "No active mission. Nothing to ship."
fi
```

→ **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- If no Contract but `.mission` exists → use Issue number from `.mission` to locate Contract or fetch Issue directly.
- **Worktree fallback**: If in worktree mode and Contract not found locally, check the main repo's `$TEAMWORK_DIR/active/` directory (parent of worktree path).
- If multiple Contracts found → "Multiple active contracts found. Keep one, remove the rest." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

Read the Contract file fully. Extract from YAML frontmatter:
- `issue` number
- `title`
- `url`
- `branch`
- `priority`

Extract from body:
- **Sub-tasks** (with checkbox status)
- **Success Criteria**
- **AI Notes**

---

## Step 1b: Worktree Guard

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
    description: "打开新终端 tab，cd {worktree_path} && claude，然后 /team-ship"
  - label: "在主仓库继续"
    description: "忽略 worktree，直接在当前目录 ship（本次）"
  - label: "关闭 worktree 模式"
    description: "git config --local teamwork.worktree false，然后继续"
```

- If "新开 tab 切到 worktree" → output: "**打开新终端 tab，执行：** `cd {worktree_path} && claude` → `/team-ship`" → **STOP**
- If "在主仓库继续" → proceed to Step 1 (skip worktree check for this run)
- If "关闭 worktree 模式" → `git config --local teamwork.worktree false` → proceed to Step 1

If no worktree exists for this branch:

**⚠️ NO WORKTREE** — branch *{branch}* has no worktree

```
question: "Worktree mode 已启用，但 branch {branch} 没有对应的 worktree。"
options:
  - label: "在主仓库 ship"
    description: "忽略 worktree，直接在当前目录 ship（本次）"
  - label: "重新 claim（创建 worktree）"
    description: "删除 Contract，重新 /team-claim #{issue}"
  - label: "关闭 worktree 模式"
    description: "git config --local teamwork.worktree false，然后继续"
```

- If "在主仓库 ship" → proceed to Step 1
- If "重新 claim" → delete Contract, output `/team-claim` **#{issue}** → **STOP**
- If "关闭 worktree 模式" → `git config --local teamwork.worktree false` → proceed to Step 1

---

## Step 2: Pre-flight Checks

### 2a: Verify branch

```bash
CURRENT_BRANCH=$(git branch --show-current)
```

Compare with Contract's `branch` field. If on wrong branch:

```bash
CONTRACT_BRANCH="{branch}"  # from Contract frontmatter
# Verify branch exists and switch to it
git checkout "$CONTRACT_BRANCH" 2>/dev/null || {
  echo "ERROR: Branch '$CONTRACT_BRANCH' not found locally or on origin."
  echo "Re-run /team-claim #{issue} to recreate the mission branch."
  exit 1
}
```

### 2a-bis: Branch Sync Check

Check how far behind base and remote:

```bash
BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "main" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-main}"

# Fetch latest to get accurate counts
git fetch origin "$BASE_BRANCH" --quiet 2>/dev/null || true
git fetch origin "$(git branch --show-current)" --quiet 2>/dev/null || true

BEHIND_BASE=$(git rev-list --count HEAD..origin/$BASE_BRANCH 2>/dev/null || echo "0")

CURRENT_BRANCH=$(git branch --show-current)
if git rev-parse --verify "origin/$CURRENT_BRANCH" >/dev/null 2>&1; then
  BEHIND_REMOTE=$(git rev-list --count HEAD..origin/$CURRENT_BRANCH 2>/dev/null || echo "0")
else
  BEHIND_REMOTE="0"
fi
```

- If `BEHIND_BASE` > 30 → 🔴 "Branch is `{N}` commits behind base." Use `AskUserQuestion`:
  ```
  question: "Branch is {N} commits behind base. Strongly recommended to rebase before shipping."
  options:
    - label: "Rebase now"
      description: "Run /team-ship sync to rebase (conflicts will be auto-resolved)"
    - label: "Ship anyway"
      description: "Continue shipping without rebasing (may cause merge conflicts on PR)"
  ```
  - If "Rebase now" → run the sync flow from **Operation Sync**, then return to Step 2a-bis to re-check
  - If "Ship anyway" → continue to Step 2b
- If `BEHIND_BASE` > 10 → 🟡 "Branch is `{N}` commits behind base. Consider `/team-ship sync` to rebase before shipping."
- If `BEHIND_REMOTE` > 0 → 🟡 "Local branch is behind remote by `{N}` commits. Run `/team-ship sync` to sync before pushing."

### 2b: Verify all sub-tasks complete

Check the Contract's sub-tasks section. Count `- [x]` vs `- [ ]`.

Separate unchecked tasks into **code tasks** (no `🔧 MANUAL` or `(MANUAL)` tag) and **manual tasks** (tagged `🔧 MANUAL` or `(MANUAL)`).

- If any unchecked **code** sub-tasks remain → display them, then AskUserQuestion:
  question: "有未完成的代码任务，如何处理？"
  options:
    - label: "去完成"
      description: "/team-drive — 继续执行"
    - label: "手动标记完成"
      description: "这些任务已完成但未勾选，标记后继续 ship"
    - label: "取消 ship"
      description: "返回，稍后再 ship"
  → If "去完成": say "Run `/team-drive` to continue" → **STOP**
  → If "手动标记完成": check off the listed code tasks in Contract, then continue ship flow
  → If "取消 ship": **STOP**
- If only **manual** sub-tasks are unchecked → proceed (these are human-action tasks, not blockers for the PR). Note them for inclusion in PR description.

**Sync checkboxes to GitHub Issue** (defense-in-depth — catches missed syncs from team-drive):
```bash
bash ~/.claude/commands/scripts/tw-contract.sh sync-all-checkboxes "$CONTRACT_PATH" $ISSUE_NUMBER
```
Non-fatal: if sync fails, warn but continue with shipping.

### 2b-bis: Comment + Review Check

**Issue Comment Check** (same logic as team-drive Step 0b1.5):

```bash
LAST_COMMENT_CHECK=$(bash ~/.claude/commands/scripts/tw-contract.sh read-field "$CONTRACT_PATH" last_comment_check 2>/dev/null) || true
COMMENT_STATUS=$(bash ~/.claude/commands/scripts/tw-contract.sh check-comments "$ISSUE_NUMBER" "$LAST_COMMENT_CHECK" 2>/dev/null) || true
```

- If "NETWORK_ERROR" → skip silently, continue
- If "UNCHANGED" → continue (no display)
- If "NEW":

Fetch full comments:
```bash
COMMENTS=$(bash ~/.claude/commands/scripts/tw-contract.sh read-comments "$ISSUE_NUMBER" --exclude-bot --exclude-lifecycle --since "$LAST_COMMENT_CHECK" 2>/dev/null)
```

Display:

⚡ **NEW COMMENTS** on **#{ISSUE_NUMBER}** since last check
────────────────────────────────────────────
  `@{author}` — *{relative_time}*
  {full body}
────────────────────────────────────────────

Judgment call:
- If any comment suggests direction change (contains: "change", "actually", "instead", "don't", "cancel", "blocker", "wait", "hold") OR is from repo MEMBER/OWNER:
  → AskUserQuestion: "Issue has new comments that may affect mission direction."
    Options: "Update Contract" / "Acknowledged, continue" / "Abort ship"
- Otherwise: display and continue (advisory)

Update timestamp:
```bash
bash ~/.claude/commands/scripts/tw-contract.sh store-comment-timestamp "$CONTRACT_PATH" "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
```

**PR Review Check** (only if PR already exists for this branch):

```bash
EXISTING_PR=$(bash ~/.claude/commands/scripts/tw-pr.sh exists "$CONTRACT_BRANCH" 2>/dev/null) || true
```

If PR exists:

```bash
gh pr view "$EXISTING_PR" --json reviewDecision,reviews
```

- If `reviewDecision` == "CHANGES_REQUESTED":

🔴 **PR REVIEW: CHANGES REQUESTED**
────────────────────────────────────────────
  `@{reviewer}`: {review body first line}
────────────────────────────────────────────
AskUserQuestion:
  question: "PR 有 review 变更请求，如何处理？"
  options:
    - label: "查看 review"
      description: "打开 PR 查看具体 feedback"
    - label: "修复后继续"
      description: "/team-drive — 修复 review 提出的问题"
    - label: "取消 ship"
      description: "返回，稍后处理"
→ If "查看 review": run `gh pr view "$EXISTING_PR" --web` → **STOP**
→ If "修复后继续": say "Run `/team-drive` to address review feedback" → **STOP**
→ If "取消 ship": **STOP**

- If `reviewDecision` == "APPROVED":
  Display: ✅ PR review: `approved` by *@{reviewer}*
  → Continue

- Otherwise (pending/null): continue silently

### 2c: Clean working tree

```bash
git status --porcelain
```

- If there are uncommitted changes → show file list and confirm:

  Display the `git status --porcelain` output, then use `AskUserQuestion`:
  ```
  question: "You have uncommitted changes. How to proceed?"
  options:
    - label: "Commit all"
      description: "Commit all changes with message: chore: pre-ship cleanup"
    - label: "Review changes first"
      description: "Show git diff for review before committing"
    - label: "Stash and continue without committing"
      description: "Stash changes — recover later with git stash pop"
  ```
  - If "Commit all":
    ```bash
    bash ~/.claude/commands/scripts/tw-git.sh commit "chore: pre-ship cleanup | Mission: #${ISSUE_NUMBER}"
    ```
  - If "Review changes first" → run `git diff`, display output, then re-ask with "Commit all" / "Stash and continue without committing" / "Abort ship"
  - If "Stash and continue without committing":
    ```bash
    git stash push -m "team-ship: stashed before shipping Mission #${ISSUE_NUMBER}"
    ```
    Warn: "Changes stashed. Recover with: `git stash pop`"

### 2d: Run tests

> Pre-flight tests guard push — but push's first job is remote persistence. When CI guards merge, blocking push for a redundant local check delays persistence for no safety gain.

```bash
CI_ENABLED=$(bash ~/.claude/commands/scripts/tw-config.sh quality.ci "false" 2>/dev/null)
```

**If `CI_ENABLED` == "true"** → skip local pre-flight:

⏭ Pre-flight skipped — CI guards merge, pushing immediately for remote persistence.

→ Continue to Step 2e (Version Bump).

**If `CI_ENABLED` != "true"** → local tests are the only quality gate. Run them:

Read config for service-aware testing:

```bash
# Check if project has services configured
# If project.services exists in config, run per-service tests for affected services
# Otherwise fall back to project.test_command

BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "main" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-main}"

# Determine affected services from changes
CHANGED_FILES=$(git diff --name-only "origin/$BASE_BRANCH...HEAD" 2>/dev/null)
```

Read `project.services` from config. For each service, check if any changed file is in that service's directory. Run the service's `test_command` for affected services. If no services configured, fall back to `project.test_command`:

```bash
TEST_CMD=$(bash ~/.claude/commands/scripts/tw-config.sh project.test_command "" 2>/dev/null)
if [ -z "$TEST_CMD" ]; then
  echo "⚠ No test_command configured and CI not enabled — no test gate at all"
  echo "  Add 'test_command: \"npm test\"' or set 'quality.ci: true' in $TEAMWORK_DIR/config.yml"
else
  echo "Running tests: $TEST_CMD"
  eval "$TEST_CMD"
fi
```

- If any test command exits non-zero → "Tests are failing. Fix them before shipping. Run /team-drive to fix, then retry /team-ship." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

---

## Step 2e: Version Bump (only changed services)

**Versions MUST be bumped on ship for services that have code changes.** Skip if last commit is already a version bump (from `/dev`).

```bash
REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)
TODAY=$(date +%Y.%m.%d)

# Skip if last commit was already a version bump
LAST_MSG=$(git log -1 --format=%s)
if echo "$LAST_MSG" | grep -q "chore: bump versions"; then
  echo "Version already bumped — skipping"
else
  # Detect changed services vs base branch
  CHANGED_FILES=$(git diff --name-only "origin/$BASE_BRANCH...HEAD" 2>/dev/null || git diff --name-only HEAD~5)

  declare -A SVC_MAP=(
    ["api-server"]="api-server"
    ["nanoclaw"]="nanoclaw"
    ["frontend"]="frontend"
    ["realtime"]="realtime"
    ["gateway"]="gateway/plugin"
  )

  BUMPED=""
  for SVC in "${!SVC_MAP[@]}"; do
    SVC_DIR="${SVC_MAP[$SVC]}"
    if echo "$CHANGED_FILES" | grep -q "^${SVC}/\|^${SVC_DIR}/"; then
      VER_FILE="$REPO_ROOT/$SVC_DIR/version.txt"
      if [ -f "$VER_FILE" ]; then
        CURRENT=$(cat "$VER_FILE" | tr -d '[:space:]')
        CUR_DATE=$(echo "$CURRENT" | sed 's/\.[^.]*$//')
        CUR_K=$(echo "$CURRENT" | awk -F. '{print $NF}')
        if [ "$CUR_DATE" = "$TODAY" ]; then
          NEW_K=$((CUR_K + 1))
        else
          NEW_K=0
        fi
        echo "${TODAY}.${NEW_K}" > "$VER_FILE"
        BUMPED="${BUMPED:+$BUMPED, }${SVC}=${TODAY}.${NEW_K}"
      fi
    fi
  done

  if [ -n "$BUMPED" ]; then
    git add */version.txt gateway/plugin/version.txt
    git commit -m "chore: bump versions — ${BUMPED}"
  fi
fi
```

---

## Step 2f: Side Effect Triage

> Detect non-code changes that could break other environments. Gate P0 issues, surface P1 for user action, log P2 in PR body. This is the Environment Envelope — mechanism captures, LLM annotates.

### Scan Phase (automated — no LLM judgment needed)

**1. Non-code file diff:**

```bash
# Get all changed files vs base branch
BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "main" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-main}"
ALL_CHANGED=$(git diff --name-only "origin/$BASE_BRANCH...HEAD" 2>/dev/null)

# Filter for side-effect-relevant file patterns
SIDE_EFFECT_FILES=$(echo "$ALL_CHANGED" | grep -iE '\.(env|env\.example|env\.local|env\.production)$|docker-compose|Dockerfile|Makefile$|\.ya?ml$|\.toml$|\.conf$|nginx|\.ini$' || true)

# Dependency file changes (separate category)
DEP_FILES=$(echo "$ALL_CHANGED" | grep -iE 'package\.json$|package-lock\.json$|yarn\.lock$|pnpm-lock\.yaml$|requirements.*\.txt$|Pipfile\.lock$|go\.sum$|Cargo\.lock$|Gemfile\.lock$' || true)
```

**2. New env var reference scan (JS/TS + Python):**

```bash
# Find all env var references in changed source files
CHANGED_SRC=$(echo "$ALL_CHANGED" | grep -iE '\.(ts|tsx|js|jsx|py)$' || true)

# Extract env var names from changed files
if [ -n "$CHANGED_SRC" ]; then
  # JS/TS: process.env.VAR_NAME or process.env['VAR_NAME'] or process.env["VAR_NAME"]
  JS_VARS=$(git diff "origin/$BASE_BRANCH...HEAD" -- $(echo "$CHANGED_SRC" | grep -iE '\.(ts|tsx|js|jsx)$' | tr '\n' ' ') 2>/dev/null | grep '^+' | grep -oE 'process\.env\.([A-Za-z0-9_]+)|process\.env\[.([A-Z_]+).\]' | sed "s/process\.env\.//;s/process\.env\[['\"]\(.*\)['\"].*/\1/" | sort -u || true)

  # Python: os.environ["VAR"], os.environ.get("VAR"), os.getenv("VAR")
  PY_VARS=$(git diff "origin/$BASE_BRANCH...HEAD" -- $(echo "$CHANGED_SRC" | grep -iE '\.py$' | tr '\n' ' ') 2>/dev/null | grep '^+' | grep -oE 'os\.(environ\[|environ\.get\(|getenv\().["'"'"']([A-Z_]+)' | grep -oE '[A-Z_]{2,}' | sort -u || true)

  NEW_VARS=$(echo "$JS_VARS $PY_VARS" | tr ' ' '\n' | sort -u | grep -v '^$' || true)
fi

# Check against .env.example (if it exists)
if [ -n "$NEW_VARS" ] && [ -f .env.example ]; then
  MISSING_VARS=""
  for VAR in $NEW_VARS; do
    if ! grep -q "^${VAR}=" .env.example 2>/dev/null; then
      MISSING_VARS="$MISSING_VARS $VAR"
    fi
  done
  MISSING_VARS=$(echo "$MISSING_VARS" | xargs)  # trim
fi
```

**3. Manual entries from AI Notes:**

Read the Contract's AI Notes section. Extract lines starting with `[side-effect]`.

```bash
MANUAL_EFFECTS=$(grep '\[side-effect\]' "$CONTRACT_PATH" 2>/dev/null || true)
```

### Classify & Display

Combine all findings into three severity levels:

**P0 — Blocks Ship** (must resolve before continuing):
- New env var references (`$MISSING_VARS`) not in `.env.example`

**P1 — Needs Action** (others may need to do something):
- `[side-effect]` entries tagged `remote:`, `credential:`, `dependency:` (native/system-level)
- Dockerfile/docker-compose changes that add new services or system deps
- Config file changes that affect deployment (new required keys, format changes)

**P2 — Info Only** (mention in PR for awareness):
- Lockfile changes (dependency additions/updates)
- Minor config adjustments
- `[side-effect]` entries tagged `env-var:` that ARE in `.env.example` (already handled)

### If zero side effects detected → skip silently, continue to Step 3.

### If side effects found → display:

**⚠️ SIDE EFFECT TRIAGE** ────────────────────────────
{If P0 items:}
🔴 **P0 — Must Resolve:**
  `{VAR_NAME}` referenced in code but missing from `.env.example`
  {for each missing var}

{If P1 items:}
🟡 **P1 — Needs Action:**
  {description of each P1 item}

{If P2 items:}
🟢 **P2 — Logged in PR:**
  {description of each P2 item}
────────────────────────────────────────────

### P0 Resolution (mandatory gate)

If `$MISSING_VARS` is non-empty:

Use `AskUserQuestion`:
```
question: "发现 {count} 个 env var 在代码中被引用但不在 .env.example 中。必须处理才能继续 ship。"
options:
  - label: "添加到 .env.example"
    description: "自动把缺失的变量追加到 .env.example（带空值和注释）"
  - label: "不需要 .env.example（有其他管理方式）"
    description: "跳过检查，在 PR 中说明原因"
  - label: "返回修复"
    description: "停止 ship，回去调整代码"
```

- "添加到 .env.example" → for each `$MISSING_VARS`, append `{VAR}=  # TODO: set value` to `.env.example`, commit
- "不需要 .env.example" → note in PR body, continue
- "返回修复" → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### P1 Triage (interactive)

For each P1 item, use `AskUserQuestion`:
```
question: "{P1 item description} — 如何处理？"
options:
  - label: "创建 Issue"
    description: "创建一个新 Issue 并 assign（推荐用于需要 infra/ops 操作的项目）"
  - label: "PR 备注"
    description: "在 PR body 中提醒 reviewer"
  - label: "跳过"
    description: "已知晓，不需要额外操作"
```

- "创建 Issue":
  ```bash
  gh issue create --title "{type}: {brief description}" --body "{details from side-effect entry}" --assignee ""
  ```
  Then use `AskUserQuestion` to ask for assignee (or leave unassigned). Add the issue link to PR body.
- "PR 备注" → add to Infra Impact section in PR body
- "跳过" → continue

**In `/team auto` mode:** skip P1 interactive triage — auto-add all P1 items to PR body as `## Infra Impact`.

### Store results for PR body injection

After triage, store the side effect summary for use in Step 4b (PR body generation):
- `$INFRA_IMPACT_ITEMS` — list of items to include in `## Infra Impact` section
- `$SIDE_EFFECT_ISSUES` — list of created issue numbers to reference

---

## Step 3: Push

```bash
bash ~/.claude/commands/scripts/tw-git.sh push {branch}
```

- If exit code 2 → push failed. Suggest `/team-ship sync` then retry push.
- If push still fails → "Push failed. Run `/team doctor` to diagnose." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

---

## Step 4: Create PR

### 4a: Check for existing PR

```bash
EXISTING_PR=$(bash ~/.claude/commands/scripts/tw-pr.sh exists {branch}) && {
  echo "PR already exists: $EXISTING_PR"
  # Skip PR creation, use existing PR number/url
}
```

If a PR already exists for this branch → skip PR creation, use existing PR.

Display existing PR status:
- CI status from PR statusCheckRollup
- Review decision
- Mergeable state
Output: "PR **#{number}** already exists — CI: {status}, Review: {status}. Pushing will update it."

### 4b: Generate PR body

**Determine Issue close keyword:**

Check the Issue's checkbox state to decide whether merging this PR should auto-close the Issue:

```bash
UNCHECKED=$(gh issue view "$ISSUE_NUMBER" --json body --jq '.body' | grep -c '\- \[ \]' || echo "0")
```

| Issue checkbox state | PR body keyword |
|---------------------|-----------------|
| All `[x]` (UNCHECKED == 0) | `Closes #{issue}` — merge auto-closes Issue |
| Has `[ ]` remaining (UNCHECKED > 0) | `Refs #{issue}` — merge does NOT close Issue |

Build the PR body from the Contract:

```markdown
{Closes or Refs} #{issue}

## Objective
{objective from Contract}

## Changes
{For each completed code sub-task:}
- {sub-task description}

{If manual tasks exist, add:}
## Manual Steps (post-merge)
{For each manual sub-task:}
- [ ] {manual task description}

## Success Criteria
{criteria from Contract — see FILTER RULE below}

**FILTER RULE**: When copying criteria from Contract to PR body, **drop any criterion
whose completion requires the PR itself to be merged or CI to pass on the PR**.
These are delivery-mechanism criteria (e.g. "PR 合入 xxx", "CI 通过", "PR merged"),
valid in the Issue but self-referential in a PR. Only include criteria about the
actual work product (code changes, behavior, test results).

## Test
{test command} — passing ✅

## AI Notes
{notes from Contract's AI Notes section, if present}

{If $INFRA_IMPACT_ITEMS is non-empty (from Step 2f Side Effect Triage):}
## Infra Impact
{For each item in $INFRA_IMPACT_ITEMS:}
- ⚠️ {item description}
{If $SIDE_EFFECT_ISSUES is non-empty:}
Related Issues: {list of #NNN references}
```

**Infra Impact rules:**
- Only include this section if Step 2f produced non-empty results
- P0 items resolved during ship (e.g., vars added to .env.example) → mention as "Resolved: {description}"
- P1 items the user chose "PR 备注" for → include with action needed
- P1 items that became Issues → reference the Issue number
- P2 items → include for awareness (no action marker)
- If zero side effects → omit this section entirely (no empty `## Infra Impact`)

**When using `Refs` (unchecked tasks remain):**

After the PR body, append a warning block:

```markdown
> ⚠️ **Issue not auto-closed** — {UNCHECKED} unchecked task(s) remain:
> {For each unchecked checkbox in the Issue body:}
> - [ ] {task text}
>
> Close manually after completing: `gh issue close {issue}`
```

Also display a local warning to the AI operator:

⚠️ Using `Refs` instead of `Closes` — **#{ISSUE_NUMBER}** has `{UNCHECKED}` unchecked task(s). Issue will NOT be auto-closed on merge.

### 4c: Create the PR

```bash
# Read base branch and extract Issue title from Contract
BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "main" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-main}"
ISSUE_TITLE=$(bash ~/.claude/commands/scripts/tw-contract.sh read-field "$TEAMWORK_DIR/active/MISSION-{issue}.md" title)

# Determine commit type from title (deterministic heuristic)
COMMIT_TYPE=$(bash ~/.claude/commands/scripts/tw-pr.sh commit-type "$ISSUE_TITLE")

# Create PR (BODY is generated by LLM from Contract content — see Step 4b)
PR_RESULT=$(bash ~/.claude/commands/scripts/tw-pr.sh create {issue} "${COMMIT_TYPE}: {title}" "$PR_BODY" "$BASE_BRANCH" "{branch}")
PR_NUMBER=$(echo "$PR_RESULT" | cut -d' ' -f1)
PR_URL=$(echo "$PR_RESULT" | cut -d' ' -f2)
```

Capture the PR number and URL from the output.

### 4d: Post PR comment to Issue

```bash
bash ~/.claude/commands/scripts/tw-pr.sh comment {issue} "📦 PR #${PR_NUMBER} created — ${PR_URL}"
```

Non-fatal: if comment fails, warn but continue.

### 4e: Add labels

```bash
bash ~/.claude/commands/scripts/tw-label.sh pr-label {pr} review
```

If this fails (label doesn't exist), warn but continue.

### 4f: Notification

```bash
bash ~/.claude/commands/scripts/tw-notify.sh mc.completed \
  --issue "$ISSUE_NUMBER" --title "$TITLE" \
  --assignee "$GH_USER" --pr "$PR_URL"
```

Non-fatal: if notification fails, warn but continue.

---

## Step 5: CI Check (if configured)

Read `quality.ci` from config.

If CI is enabled:
```bash
bash ~/.claude/commands/scripts/tw-pr.sh watch {pr}
```

- If CI passes → "CI passed ✅"
- If CI fails → "CI failed ❌. Check the PR for details: {pr-url}" (do NOT stop — the PR is already created, user can fix and push again)
  "Fix the failures, commit, and push to {branch} — the PR will update automatically. Or run /team-drive to fix in mission mode."

If CI is not enabled → skip this step.

---

## Step 5b: AI-Suggested PR Comment (conditional)

**Evaluate whether the PR would benefit from an additional comment** beyond the PR body. Review the execution context and determine if any of these are true:

| Signal | Why it's worth commenting |
|--------|--------------------------|
| Deployment needs special steps NOT already in `## Infra Impact` | Merger needs to know before/after merge actions (skip if Step 2f already covered) |
| Specific files/areas need careful review | Focus reviewer attention |
| Known limitation or tech debt introduced | Transparency for the team |
| CI failed on a flaky test (not a real issue) | Prevent reviewer confusion |
| Breaking change or API change | Downstream consumers need to know |

**If none of these signals are present** → **skip silently**.

**If any signal is present** → generate a concise comment draft:

```
💬 **建议追加 PR 备注** — 有信息值得提醒 reviewer/merger：

{draft comment — 2-5 sentences, actionable}
```

Use `AskUserQuestion`:
```
question: "追加这条备注到 PR？"
options:
  - label: "发送"
    description: "直接发送上述备注"
  - label: "编辑后发送"
    description: "我修改一下内容"
  - label: "跳过"
    description: "不需要追加"
```

- "发送" → `gh pr comment {pr} --body "{draft}"` (non-fatal)
- "编辑后发送" → use `AskUserQuestion` with `useTextArea: true` to collect edited text → post
- "跳过" → continue

**In `/team auto` mode:** skip this step entirely.

---

## Step 6: Request Review (if configured)

Read `quality.review_required` from config.

If review is required:
- Find the tech-lead from the team roster in config (use LLM to parse config and find team member with `role: tech-lead`).
- Request review:
  ```bash
  bash ~/.claude/commands/scripts/tw-pr.sh add-reviewer {pr} {tech-lead-github}
  ```
- If reviewer assignment fails (permissions), warn but continue.

If review is not required → skip this step.

---

## Step 7: Cleanup

### 7a: Update Issue labels

```bash
# Transition wip → review (reads prefix from config)
bash ~/.claude/commands/scripts/tw-label.sh transition {issue} wip review

# Verify the update
bash ~/.claude/commands/scripts/tw-label.sh verify {issue} review
```

- If verify returns "VERIFIED" → "Issue labeled review ✅"
- If verify returns "MISMATCH" → warn: "Label update incomplete"
- If verify fails (network) → "⚠ Could not verify labels. Assume update succeeded."
- Non-blocking in all cases: continue to 7b.

### 7b: Remove Contract (conditional on CI status)

If CI passed OR CI was not configured → delete the Contract:

```bash
bash ~/.claude/commands/scripts/tw-contract.sh delete "$TEAMWORK_DIR/active/MISSION-{issue}.md"
```

The Contract has served its purpose. The PR body now contains the essential information. If the PR is rejected and the mission needs to be reworked, re-run `/team-claim` **#{issue}** to regenerate a fresh Contract.

If CI failed → keep the Contract:

Contract preserved — run `/team-drive` to iterate on CI failures, then `/team-ship` again.

---

## Step 8: Output Delivery Summary

**📦 SHIPPED** ── **#{issue}** {title} ─────────────
🔀 *{branch}*
🔀 PR: {pr-url}
{ci_icon} CI: `{passing/failing/not configured}`
Review:    `{requested from @{reviewer} / not required}`
Labels:    `{actual current labels}`

**📋 SUB-TASKS DELIVERED**
  ✅ {task 1}
  ✅ {task 2}

────────────────────────────────────────────
`/team-ship done` (after merge) │ `/team-ship review` │ `/team`

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

---

## Operation Done (post-merge cleanup)

> Triggered by `/team-ship done`. Verifies PR merged, closes Issue, cleans up worktree.

### D1: Find the Issue

Locate the Issue number from (in priority order):
1. Active Contract (`$TEAMWORK_DIR/active/MISSION-*.md` frontmatter)
2. `.mission` file (worktree mode)
3. Current branch name — extract first numeric segment after `/` (works for both `mission/42-slug` and `bugfix/T-044-slug` patterns)
4. `$ARGUMENTS` — if user passes `done {issue}` (e.g., `done 42`)

If no Issue found:
  - If current branch matches `mission/*` pattern → extract the numeric segment, then: "Issue **#{N}** (from branch name) not found. It may have been deleted. Provide the correct Issue number: `/team-ship done #NNN`" → **STOP**
  - Otherwise → "Cannot determine Issue. Provide Issue number: `/team-ship done {issue}`" → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### D2: Verify PR is merged

```bash
BRANCH=$(bash ~/.claude/commands/scripts/tw-git.sh current)
bash ~/.claude/commands/scripts/tw-pr.sh verify-merged "$BRANCH"
```

- If exit code 1 → display and **STOP**:

**No merged PR found** for branch *{branch}*.
────────────────────────────────────────────
PR must be merged on GitHub before running `done`.
  1. Check PR status: `gh pr view {branch}`
  2. If just merged, GitHub may need a few seconds — wait and retry: `/team-ship done`
  3. If PR is still open, merge it on GitHub first
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### D3: Close Issue (fallback — may be needed if PR used "Refs #N")

```bash
ISSUE_STATE=$(gh issue view {issue} --json state --jq '.state' 2>/dev/null)
```

If Issue is still open (PR used `Refs #N` due to unchecked tasks, or `Closes #N` didn't fire):
```bash
gh issue close {issue} --reason completed
```

If already closed → skip, output: `(Issue already closed — handled by PR merge)`

### D4: Update labels (fallback — usually handled by post-merge Action)

```bash
ISSUE_LABELS=$(gh issue view {issue} --json labels --jq '[.labels[].name] | join(",")' 2>/dev/null)
```

If `${STATUS_PREFIX}done` is NOT in labels (post-merge Action didn't run or failed):
```bash
bash ~/.claude/commands/scripts/tw-label.sh transition {issue} review done
```

If `${STATUS_PREFIX}done` already present → skip, output: `(Labels already updated — handled by post-merge Action)`

### D5: Clean up worktree (if applicable)

If `.mission` file exists (worktree mode):
- Derive the main repo path first:
  ```bash
  MAIN_REPO=$(git -C "$(pwd)" rev-parse --git-common-dir | sed 's|/.git$||')
  ```
- Ask user: "Remove worktree at `{worktree_path}`?" via `AskUserQuestion`
- If yes:
  ```bash
  WORKTREE_PATH=$(pwd)
  cd "$MAIN_REPO"
  bash ~/.claude/commands/scripts/tw-git.sh worktree-remove "$WORKTREE_PATH"
  ```
- If no: `cd "$MAIN_REPO"` first (must leave worktree for D6), then warn "Worktree kept. Run `/team doctor fix` later to clean up stale worktrees."

### D6: Return to base branch

```bash
# Check for uncommitted changes BEFORE switching (exclude untracked-only)
DIRTY=$(git status --porcelain 2>/dev/null | grep -v '^??')
if [ -n "$DIRTY" ]; then
  STASH_MSG="team-ship-done: stashed from $(git branch --show-current) before cleanup"
  git stash push -m "$STASH_MSG" || {
    echo "ERROR: git stash failed. Commit or discard changes manually before cleanup."
    # STOP — do not proceed to checkout, risk losing changes
    # 💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
  }
  echo "⚠️ Stashed uncommitted changes. Recover with: git stash pop"
fi

bash ~/.claude/commands/scripts/tw-git.sh ensure-base || {
  echo "ERROR: Could not switch to base branch."
  if [ -n "$DIRTY" ]; then
    echo "Your changes are saved in stash. Recover with: git stash pop"
  fi
  # STOP — 💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
}
# Note: ensure-base already does checkout + pull, so base branch is up-to-date
# (includes the just-merged PR squash commit)
```

### D7: Clean up Contract

⚠️ Re-detect TEAMWORK_DIR here — the value from flow start may point to a removed worktree.

```bash
# Re-detect after worktree removal + base branch switch
TEAMWORK_DIR=$(bash ~/.claude/commands/scripts/tw-config.sh detect-dir 2>/dev/null)
TEAMWORK_DIR="${TEAMWORK_DIR:-.teamwork}"

bash ~/.claude/commands/scripts/tw-contract.sh delete "$TEAMWORK_DIR/active/MISSION-{issue}.md"

# Fallback: if file still exists (path mismatch), find and delete
if ls .teamwork/active/MISSION-{issue}.md .teamspace/active/MISSION-{issue}.md 2>/dev/null; then
  rm -f .teamwork/active/MISSION-{issue}.md .teamspace/active/MISSION-{issue}.md 2>/dev/null
  echo "Contract cleaned up via fallback path"
fi
```

### D8: Output summary

**🏁 DONE** ── **#{issue}** {title} ────────────────
Issue:     `closed`
🔀 PR: `{pr-url}` (merged)
Labels:    `{STATUS_PREFIX}done`
Worktree:  `{removed/kept/n/a}`
Branch:    *{BASE_BRANCH}*

────────────────────────────────────────────
`/team-claim` (next mission) │ `/team-rc` (release) │ `/team` (dashboard)

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

---

## Operation Review (AI PR review)

> Triggered by `/team-ship review`. Fetches PR diff, reviews code, publishes review.

### R1: Find PR

```bash
BRANCH=$(git branch --show-current)
gh pr list --head "$BRANCH" --state open --json number,url --limit 1
```

If `$ARGUMENTS` contains a PR number → use that instead.
If no open PR found → "No open PR found for branch `{branch}`." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### R2: Get PR diff

```bash
gh pr diff {pr}
```

### R3: Read changed files

For each file in the diff, read the full current content to understand context beyond just the diff hunks.

### R4: AI Review

Analyze the changes for:
- **Correctness**: Logic errors, missing edge cases, off-by-one errors
- **Security**: Injection vulnerabilities, exposed secrets, unsafe operations
- **Architecture**: Coupling, separation of concerns, pattern consistency
- **Quality**: Naming, readability, dead code, duplication
- **Performance**: N+1 queries, unnecessary allocations, missing indexes

### R5: Choose review action

Use `AskUserQuestion`:
```
question: "What review action?"
options:
  - label: "Approve"
    description: "Changes look good — approve the PR"
  - label: "Request changes"
    description: "Issues found — request fixes before merge"
  - label: "Comment only"
    description: "Leave feedback without approval decision"
```

### R6: Publish review

```bash
gh pr review {pr} --body "{review_body}" {--approve|--request-changes|--comment}
```

Output:

Review published on PR **#{pr}**: {approve/request-changes/comment}
After merge: `/team-ship done` │ Address comments: push more commits

💡 Tip: {random tip — read `~/.claude/commands/scripts/tw-tips.txt`, pick one non-comment line at random}

---

## Operation Sync (rebase on base branch)

> Triggered by `/team-ship sync`. Rebases current branch on latest base branch.

### S1: Verify on mission branch

```bash
BRANCH=$(bash ~/.claude/commands/scripts/tw-git.sh current)
BASE_BRANCH=$(bash ~/.claude/commands/scripts/tw-config.sh conventions.base_branch "main" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-main}"
```

If on `$BASE_BRANCH` → "Already on $BASE_BRANCH. Nothing to sync." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### S2: Fetch and rebase

```bash
bash ~/.claude/commands/scripts/tw-git.sh rebase "$BASE_BRANCH"
```

- If exit 0 → jump to **S3: Summary**
- If exit 2 → jump to **S2b: AI conflict resolution**

### S2b: AI conflict resolution

When rebase produces conflicts, resolve them automatically:

1. **List conflicted files:**
```bash
git diff --name-only --diff-filter=U
```

2. **For each conflicted file:**
   - Read the file (it contains `<<<<<<<`, `=======`, `>>>>>>>` markers)
   - Understand both sides: "ours" (mission branch) = the work being rebased, "theirs" (base branch) = upstream changes
   - Edit the file to resolve: merge both sides' intent, preserving mission changes where they don't contradict upstream
   - `git add {file}`

3. **Continue rebase:**
```bash
git rebase --continue
```
   - If more conflicts appear → repeat step 1-3
   - If rebase completes → jump to **S3: Summary**

4. **If a conflict is genuinely ambiguous** (both sides made semantic changes to the same logic that contradict each other), resolve with best judgment and flag it in the summary.

### S3: Summary

Output:

**SYNC COMPLETE** ────────────────────────────
Branch:   *{branch}*
Rebased:  on latest *{BASE_BRANCH}*
{If conflicts resolved:}
Resolved: `{N}` conflict(s) in {file list}
{If any flagged:}
⚠️ Review: {file} — both sides changed {description}, merged as {approach}

Push the rebased branch:
```bash
git push --force-with-lease origin {branch}
```

- If push succeeds → continue to **S4: Info Sync**
- If push fails → "Push failed. Run `/team-ship sync` to rebase, then retry `/team-ship`." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.

### S4: Info Sync

After rebase completes, check for new information from GitHub.

**Prerequisites**: Active Contract exists (need `claimed` timestamp and issue number).
If no active Contract → skip Info Sync entirely (graceful degradation).

```bash
bash ~/.claude/commands/scripts/tw-contract.sh info-sync {ISSUE_NUMBER} "{CONTRACT_PATH}"
```

Parse the output (3 lines):
- Line 1: `COMMENTS:NEW` or `COMMENTS:UNCHANGED` or `COMMENTS:NETWORK_ERROR`
- Line 2: `PR:{status}:{count}` or `PR:none` or `PR:NETWORK_ERROR`
- Line 3: `CONTRACT:fresh` or `CONTRACT:stale` or `CONTRACT:NETWORK_ERROR`

**Display rules**:
- If ALL three are unchanged/fresh/none/NETWORK_ERROR → skip INFO SYNC block entirely
- Otherwise display:

**INFO SYNC** ────────────────────────────────
{If COMMENTS:NEW:}
  Issue:    **#{ISSUE}** — new comments since claim
{If PR exists and has review status:}
  PR:       `{reviewDecision}` {If changes_requested: — `{count}` unresolved thread(s)}
{If CONTRACT:stale:}
  Contract: `stale` — Issue body modified since Contract was generated
{If CONTRACT:fresh:}
  Contract: `fresh`
────────────────────────────────────────────

{If CONTRACT:stale:}
  ⚠️ Run `/team-drive` to review Issue changes and update Contract.

**All items are advisory** — sync never stops or prompts. The developer reads and decides.

Update comment timestamp after sync:
```bash
bash ~/.claude/commands/scripts/tw-contract.sh store-comment-timestamp "{CONTRACT_PATH}" "{CURRENT_ISO_TIMESTAMP}"
```

---

## Operation Done Help

> Triggered by `/team-ship done help`.

Output the following and **STOP**:

**`/team-ship done` — Post-Merge Cleanup**
────────────────────────────────────────────

**USAGE**
  `/team-ship done`           Auto-detect Issue from Contract/branch
  `/team-ship done {issue}`   Specify Issue number explicitly

**WHAT IT DOES**
  1. Verify PR is merged (fails if PR still open)
  2. Close Issue (if not auto-closed by `Closes #N`)
  3. Update labels: `status:review` → `status:done`
  4. Clean up worktree (if in worktree mode)
  5. Return to base branch
  6. Remove local Contract

**WHEN TO USE**
  After your PR has been merged on GitHub.
  Usually the last step: `/team-drive` → `/team-ship` → (merge) → `/team-ship done`

**TROUBLESHOOTING**
  "No merged PR found" → Check `gh pr view {branch}`. If just merged, wait a few seconds and retry.
  "Cannot determine Issue" → Specify explicitly: `/team-ship done {issue}`

---

## Operation Review Help

> Triggered by `/team-ship review help`.

Output the following and **STOP**:

**`/team-ship review` — AI Code Review**
────────────────────────────────────────────

**USAGE**
  `/team-ship review`         Review PR for current branch
  `/team-ship review {pr}`    Review specific PR number

**WHAT IT DOES**
  1. Find open PR for current branch
  2. Fetch full diff
  3. Read all changed files for context
  4. AI review: correctness, security, architecture, quality, performance
  5. Choose action: approve / request changes / comment only
  6. Publish review to GitHub

**WHEN TO USE**
  Before merging a PR, or when a teammate asks for review.

---

## Operation Sync Help

> Triggered by `/team-ship sync help`.

Output the following and **STOP**:

**`/team-ship sync` — Rebase on Base Branch**
────────────────────────────────────────────

**USAGE**
  `/team-ship sync`           Rebase current branch on latest base

**WHAT IT DOES**
  1. Verify you're on a mission branch (not base)
  2. Fetch latest base branch from remote
  3. Rebase your branch on top of latest base
  4. Auto-resolve conflicts (AI reads both sides and merges)
  5. Force-push rebased branch to remote

**WHEN TO USE**
  When dashboard shows "⚠️ `N` behind *{base}*" on your PR.
  When `/team doctor` shows your branch is behind base.
  Before `/team-ship` if branch has drifted.

**TROUBLESHOOTING**
  Already on base branch → Nothing to sync.
  Ambiguous conflict → AI resolves + flags for review in summary.

---

## Error Handling

- Config missing → "Run `/team` first to initialize teamwork." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- No active Contract → "No mission to ship. Run `/team-claim` first." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- Incomplete code sub-tasks → warn with list of remaining tasks → **STOP** (MANUAL tasks don't block)
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- Tests fail → "Fix tests before shipping." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- Push fails → suggest `/team-ship sync`, if still fails → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- PR creation fails → "Failed to create PR. Check `gh auth status` and try again." → **STOP**
  💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
- CI fails → warn but DO NOT stop (PR is already created, user can iterate)
- Label/reviewer operations fail → warn but continue (non-fatal)
- Contract removal fails → warn but continue (non-fatal)

**On any STOP:** Always append: 💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
