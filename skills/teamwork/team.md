---
description: "Initialize GitHub-first teamwork system or show team dashboard. Run first to set up, then anytime for overview."
---

# /team — Init + Dashboard (Teamwork v2)

> GitHub-first team coordination. Issues are tasks, PRs are delivery, Actions are quality gates.

**Two modes:**
1. **First run** (no `.teamwork/config.yml`): Initialize teamwork system — detect project, configure team, generate CI/hooks/templates
2. **Subsequent runs**: Show team dashboard from GitHub data

**User input**: $ARGUMENTS

---

## Step 0: Prerequisites

```bash
# Check gh CLI
gh --version 2>/dev/null || echo "GH_MISSING"

# Check git remote
git remote get-url origin 2>/dev/null || echo "NO_REMOTE"
```

- If `gh` not installed → "GitHub CLI (`gh`) is required. Install: https://cli.github.com/" → **STOP**
- If no git remote → "This is not a GitHub-linked repository. Run `git remote add origin <url>` first." → **STOP**

```bash
# Get identity
GH_USER=$(gh api user --jq '.login' 2>/dev/null)
```

- If `gh api user` fails → "Not authenticated. Run `gh auth login` first." → **STOP**

```bash
# Get repo info
REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
```

---

## Step 1: Route — Init or Dashboard?

```bash
[ -f .teamwork/config.yml ] && echo "DASHBOARD" || echo "INIT"
```

- If `DASHBOARD` → Jump to **Step 6: Dashboard**
- If `INIT` → Continue to **Step 2: Initialize**

---

## Step 2: Detect Project

Auto-detect language and tooling by checking for marker files:

| Marker File | Language | Test | Lint | Build |
|-------------|----------|------|------|-------|
| `package.json` | Node/TypeScript | `npm test` | `npm run lint` | `npm run build` |
| `pyproject.toml` or `setup.py` | Python | `pytest` | `ruff check .` | `python -m build` |
| `go.mod` | Go | `go test ./...` | `golangci-lint run` | `go build ./...` |
| `Cargo.toml` | Rust | `cargo test` | `cargo clippy` | `cargo build` |
| `pom.xml` or `build.gradle` | Java | `mvn test` | — | `mvn package` |

```bash
ls package.json pyproject.toml setup.py go.mod Cargo.toml pom.xml build.gradle 2>/dev/null
```

Store detected `LANGUAGE`, `TEST_CMD`, `LINT_CMD`, `BUILD_CMD`.

If `package.json` exists, also check for scripts:
```bash
node -e "const p=require('./package.json'); console.log(JSON.stringify(p.scripts||{}))" 2>/dev/null
```
Use actual script names (e.g., `npm run test`, `npm run lint`) if they exist.

---

## Step 3: Interactive Team Configuration

Use `AskUserQuestion` to configure the team.

**Question 1: Team Members**

```
question: "Who are the team members? (Enter GitHub usernames, comma-separated)"
options:
  - label: "Just me"
    description: "Solo project with teamwork structure for task tracking"
  - label: "2-3 people"
    description: "Small team. I'll list GitHub usernames."
  - label: "4-5 people"
    description: "Medium team. I'll list GitHub usernames."
```

If team members selected, ask for usernames + roles via follow-up AskUserQuestion.

**Question 2: Quality preferences**

```
question: "Which quality gates do you want?"
multiSelect: true
options:
  - label: "Git hooks (pre-commit + pre-push)"
    description: "Local quality gates: lint on commit, test on push"
  - label: "GitHub Actions CI"
    description: "Remote quality gates: lint + test + build on every PR"
  - label: "Branch protection"
    description: "Require CI pass + PR review before merge to main"
  - label: "All of the above (Recommended)"
    description: "Full defense-in-depth: hooks + Actions + branch protection"
```

---

## Step 4: Generate Files

### 4a: Create directory structure

```bash
mkdir -p .teamwork/active .github/ISSUE_TEMPLATE .github/workflows .githooks
```

### 4b: Write `.teamwork/config.yml`

Write based on detected project info + user answers:

```yaml
schema_version: 1
skill_version: 2.0.0

team:
  - github: {GH_USER}
    role: tech-lead
  # ... additional members from Step 3

project:
  language: {LANGUAGE}
  test_command: "{TEST_CMD}"
  lint_command: "{LINT_CMD}"
  build_command: "{BUILD_CMD}"

conventions:
  branch_pattern: "mission/{issue}-{slug}-{user}"
  base_branch: main
  commit_format: "type(scope): description | Mission: #{issue}"

quality:
  hooks: {true/false from Step 3}
  ci: {true/false from Step 3}
  review_required: {true/false from Step 3}
  branch_protection: {true/false from Step 3}
```

### 4c: Write `.github/ISSUE_TEMPLATE/mission.yml`

```yaml
name: Mission Contract
description: Structured task for team execution
labels: ["mission"]
body:
  - type: dropdown
    id: priority
    attributes:
      label: Priority
      options:
        - P0 (critical)
        - P1 (high)
        - P2 (medium)
        - P3 (low)
    validations:
      required: true

  - type: textarea
    id: objective
    attributes:
      label: Objective
      description: What needs to be accomplished?
    validations:
      required: true

  - type: textarea
    id: subtasks
    attributes:
      label: Sub-tasks
      description: Break down into checkable items
      value: |
        - [ ] Task 1
        - [ ] Task 2
    validations:
      required: true

  - type: textarea
    id: acceptance
    attributes:
      label: Acceptance Criteria
      description: How do we know this is done?
    validations:
      required: true

  - type: textarea
    id: context
    attributes:
      label: Context & References
      description: Relevant files, docs, links

  - type: input
    id: test_command
    attributes:
      label: Test Command
      description: Command to verify this mission
      placeholder: "npm test -- --grep auth"
```

### 4d: Write `.github/workflows/ci.yml`

Generate based on detected language:

**Node/TypeScript:**
```yaml
name: CI
on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: npm
      - run: npm ci
      - run: {LINT_CMD}
      - run: {TEST_CMD}
      - run: {BUILD_CMD}
```

**Python:**
```yaml
name: CI
on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: "3.12"
      - run: pip install -e ".[dev]"
      - run: {LINT_CMD}
      - run: {TEST_CMD}
```

**Go:**
```yaml
name: CI
on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: stable
      - run: {LINT_CMD}
      - run: {TEST_CMD}
      - run: {BUILD_CMD}
```

For other languages or if language not detected, generate a minimal workflow with just checkout and a comment explaining what to add.

### 4e: Write git hooks (if enabled)

**`.githooks/pre-commit`:**
```bash
#!/bin/bash
# Teamwork v2 — local quality gate (fast)
set -e
CONFIG=".teamwork/config.yml"
[ -f "$CONFIG" ] || exit 0
LINT_CMD=$(grep 'lint_command:' "$CONFIG" | sed 's/.*lint_command: *//' | tr -d '"')
if [ -n "$LINT_CMD" ]; then
  echo "Running lint..."
  eval "$LINT_CMD"
fi
```

**`.githooks/pre-push`:**
```bash
#!/bin/bash
# Teamwork v2 — local quality gate (thorough)
set -e
CONFIG=".teamwork/config.yml"
[ -f "$CONFIG" ] || exit 0
TEST_CMD=$(grep 'test_command:' "$CONFIG" | sed 's/.*test_command: *//' | tr -d '"')
if [ -n "$TEST_CMD" ]; then
  echo "Running tests..."
  eval "$TEST_CMD"
fi
```

Make hooks executable:
```bash
chmod +x .githooks/pre-commit .githooks/pre-push
```

Configure git to use local hooks:
```bash
git config core.hooksPath .githooks
```

### 4f: Update `.gitignore`

Append if not already present:
```
# Teamwork ephemeral contracts
.teamwork/active/
```

---

## Step 5: Configure GitHub

### 5a: Create labels

```bash
gh label create mission --color 0075ca --description "Team mission" --force
gh label create P0 --color d73a4a --description "Critical priority" --force
gh label create P1 --color e4e669 --description "High priority" --force
gh label create P2 --color 0e8a16 --description "Medium priority" --force
gh label create P3 --color cfd3d7 --description "Low priority" --force
gh label create in-progress --color fbca04 --description "Currently being worked on" --force
gh label create review --color 7057ff --description "Ready for review" --force
```

### 5b: Branch protection (if enabled)

```bash
gh api repos/{REPO}/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["quality"]}' \
  --field enforce_admins=false \
  --field required_pull_request_reviews='{"required_approving_review_count":1}' \
  --field restrictions=null
```

If branch protection fails (e.g., free plan limitations), warn but continue — it's not a blocker.

---

## Step 5c: Commit and Push

```bash
git add .teamwork/config.yml .github/ .githooks/ .gitignore
git commit -m "feat(teamwork): initialize GitHub-first teamwork v2

Generate config, CI workflow, Issue template, git hooks.
Team: {member list}"
git push
```

Output: "Teamwork initialized! Next: create Issues using the Mission template, then `/team-claim` to start working."

---

## Step 6: Dashboard

Read data from GitHub and config, then display formatted dashboard.

### 6a: Read config
```bash
cat .teamwork/config.yml
```

### 6b: Fetch GitHub data
```bash
# All mission issues
gh issue list --label mission --json number,title,assignees,labels,state,milestone --limit 50

# Open PRs
gh pr list --json number,title,author,headRefName,statusCheckRollup,reviewDecision --limit 20

# Merge count per team member (merged PRs to base branch)
gh pr list --state merged --base main --json author --limit 100
```

### 6c: Format and display

Output a formatted dashboard like this:

```
TEAM DASHBOARD — {repo name}
Milestone: {active milestone} ({closed}/{total} = {pct}%)
════════════════════════════════════════════

{For each team member from config.yml:}
  {username} ({role})
     {If has assigned in-progress issue:} Working on: #{N} {title} [{priority}]
     {If has open PR:} PR: #{pr} — {CI status}
     {If idle:} Idle (last merged: #{last_merged_pr})

QUEUED:
  {Issues with label "mission" that are unassigned or not in-progress}
     #{N} {title} [{priority}]

MERGE COUNT:
  {username}: {count} | {username}: {count} | ...

════════════════════════════════════════════
Commands: /team-claim #{N} | /team-drive | /team-ship
```

---

## Error Handling

- `gh` commands may fail due to network, permissions, or rate limits
- On non-fatal errors (label creation fails, branch protection fails): warn and continue
- On fatal errors (no gh, no remote, not authenticated): clear message and STOP
- If config.yml exists but is malformed: warn and offer to regenerate
