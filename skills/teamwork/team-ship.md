---
description: "Ship your completed mission — push, create PR with Issue linking, check CI, clean up Contract. The delivery step."
---

# /team-ship — Deliver Mission

> Push your branch, create a PR that auto-closes the Issue, optionally watch CI, and clean up the local Contract.

**User input**: $ARGUMENTS

---

## Step 0: Prerequisites

```bash
# Identity
GH_USER=$(gh api user --jq '.login' 2>/dev/null)
```

- If fails → "Not authenticated. Run `gh auth login` first." → **STOP**

```bash
# Config (support both directory names)
if [ -f .teamwork/config.yml ]; then
  TEAMWORK_DIR=".teamwork"
elif [ -f .teamspace/config.yml ]; then
  TEAMWORK_DIR=".teamspace"
else
  echo "NO_CONFIG"
fi
```

- If no config → "Teamwork not initialized. Run `/team` first." → **STOP**

Read `$TEAMWORK_DIR/config.yml` → extract project settings, conventions, quality preferences.

---

## Step 1: Find Active Contract

```bash
ls $TEAMWORK_DIR/active/MISSION-*.md 2>/dev/null
```

- If no Contract found → "No active mission. Nothing to ship." → **STOP**
- If multiple Contracts found → "Multiple active contracts found. Keep one, remove the rest." → **STOP**

Read the Contract file fully. Extract from YAML frontmatter:
- `issue` number
- `title`
- `url`
- `branch`
- `priority`

Extract from body:
- **Sub-tasks** (with checkbox status)
- **Acceptance Criteria**
- **AI Notes**

---

## Step 2: Pre-flight Checks

### 2a: Verify branch

```bash
CURRENT_BRANCH=$(git branch --show-current)
```

Compare with Contract's `branch` field.
- If on wrong branch → `git checkout {contract.branch}`
- If branch doesn't exist → "Branch {branch} not found. Cannot ship." → **STOP**

### 2b: Verify all sub-tasks complete

Check the Contract's sub-tasks section. Count `- [x]` vs `- [ ]`.
- If any unchecked sub-tasks remain → display them and say:
  "Not all sub-tasks are complete. Finish them with `/team-drive` first, or manually check them off in the Contract."
  → **STOP**

### 2c: Clean working tree

```bash
git status --porcelain
```

- If there are uncommitted changes → warn: "You have uncommitted changes. Committing them now."
  ```bash
  git add -A
  git commit -m "chore: pre-ship cleanup | Mission: #{issue}"
  ```

### 2d: Run tests

Read `test_command` from config:
```bash
# From config.yml project.test_command
{TEST_CMD}
```

- If tests fail → "Tests are failing. Fix them before shipping." → **STOP**
- If no test command configured → warn "No test command configured. Skipping pre-flight tests."

---

## Step 3: Push

```bash
git push -u origin {branch}
```

- If push fails due to upstream changes → suggest:
  ```bash
  git pull --rebase origin {branch}
  git push -u origin {branch}
  ```
- If push still fails → "Push failed. Resolve the issue manually." → **STOP**

---

## Step 4: Create PR

### 4a: Check for existing PR

```bash
gh pr list --head {branch} --json number,url --limit 1
```

If a PR already exists for this branch → skip PR creation, use existing PR. Display: "PR already exists: {url}"

### 4b: Generate PR body

Build the PR body from the Contract:

```markdown
Closes #{issue}

## Objective
{objective from Contract}

## Changes
{For each completed sub-task:}
- {sub-task description}

## Acceptance Criteria
{criteria from Contract}

## Test
{test command} — passing ✅

## AI Notes
{notes from Contract's AI Notes section, if present}
```

### 4c: Create the PR

Read `base_branch` from config conventions (default: `main`).

```bash
gh pr create \
  --title "{commit_type}: {contract.title}" \
  --body "{generated PR body}" \
  --base {base_branch} \
  --head {branch}
```

Determine commit type from the Contract title/objective:
- If title contains "fix", "bug", "patch" → `fix`
- If title contains "refactor", "restructure" → `refactor`
- Otherwise → `feat`

Capture the PR number and URL from the output.

### 4d: Add labels

```bash
gh pr edit {pr-number} --add-label review
```

If this fails (label doesn't exist), warn but continue.

---

## Step 5: CI Check (if configured)

Read `quality.ci` from config.

If CI is enabled:
```bash
gh pr checks {pr-number} --watch --interval 10
```

- If CI passes → "CI passed ✅"
- If CI fails → "CI failed ❌. Check the PR for details: {pr-url}" (do NOT stop — the PR is already created, user can fix and push again)

If CI is not enabled → skip this step.

---

## Step 6: Request Review (if configured)

Read `quality.review_required` from config.

If review is required:
- Find the tech-lead from the team roster in config:
  ```bash
  # Extract tech-lead github username from config
  ```
- Request review:
  ```bash
  gh pr edit {pr-number} --add-reviewer {tech-lead-github}
  ```
- If reviewer assignment fails (permissions), warn but continue.

If review is not required → skip this step.

---

## Step 7: Cleanup

### 7a: Update Issue labels

```bash
gh issue edit {issue} --remove-label in-progress --add-label review
```

If label operations fail, warn but continue.

### 7b: Remove Contract

```bash
rm $TEAMWORK_DIR/active/MISSION-{issue}.md
```

The Contract has served its purpose. The PR body now contains the essential information. If the PR is rejected and the mission needs to be reworked, re-run `/team-claim #{issue}` to regenerate a fresh Contract.

---

## Step 8: Output Delivery Summary

```
MISSION SHIPPED
═══════════════════════════════════════
Issue:  #{issue} — {title}
Branch: {branch}
PR:     {pr-url}
CI:     {passing/failing/not configured}
Review: {requested from {reviewer} / not required}

Sub-tasks delivered:
  [x] {task 1}
  [x] {task 2}
  ...

Contract: cleaned up ✅
Issue:    labeled "review" ✅
═══════════════════════════════════════
Next steps:
  - PR will auto-close Issue #{issue} on merge (via "Closes #{issue}")
  - Branch will auto-delete on merge (if GitHub setting enabled)
  - To claim next mission: /team-claim
  - To see team status: /team
```

---

## Error Handling

- Config missing → "Run `/team` first to initialize teamwork." → **STOP**
- No active Contract → "No mission to ship. Run `/team-claim` first." → **STOP**
- Incomplete sub-tasks → warn with list of remaining tasks → **STOP**
- Tests fail → "Fix tests before shipping." → **STOP**
- Push fails → suggest rebase, if still fails → **STOP**
- PR creation fails → "Failed to create PR. Check `gh auth status` and try again." → **STOP**
- CI fails → warn but DO NOT stop (PR is already created, user can iterate)
- Label/reviewer operations fail → warn but continue (non-fatal)
- Contract removal fails → warn but continue (non-fatal)
