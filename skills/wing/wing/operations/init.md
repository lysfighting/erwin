# Operation Init — Project Setup

> First-time setup or re-initialization. Detects project, configures team, generates files.
> context:fork candidate — heavy, runs once per project.

## Worktree Guard

Init modifies config + commits to base branch — MUST run from main repo.

```bash
if [ -f .mission ]; then
  MAIN_REPO=$(bash "${CLAUDE_SKILL_DIR}/../scripts/wing-git.sh" worktree-main-repo)
  echo "⚠️ /wing init must run from main repo. cd $MAIN_REPO"
  # STOP
fi
```

## Step 1: Detect Project

Auto-detect language and tooling by checking marker files:

| Marker | Language | Test | Lint |
|--------|----------|------|------|
| `package.json` | Node/TS | `npm test` | `npm run lint` |
| `pyproject.toml` / `setup.py` | Python | `pytest` | `ruff check .` |
| `go.mod` | Go | `go test ./...` | `golangci-lint run` |
| `Cargo.toml` | Rust | `cargo test` | `cargo clippy` |
| `pom.xml` / `build.gradle` | Java | `mvn test` | — |

```bash
ls package.json pyproject.toml setup.py go.mod Cargo.toml pom.xml build.gradle 2>/dev/null
```

If no root markers → check one level deep for monorepo structure. If multiple services → set `IS_MONOREPO=true`, detect each service independently.

For `package.json` projects, also read `scripts` for actual test/lint/build commands.

Store: `LANGUAGE`, `TEST_CMD`, `LINT_CMD`, `BUILD_CMD`, optionally `IS_MONOREPO`, `SERVICES[]`.

## Step 2: Interactive Configuration

### Q1: Team Members

AskUserQuestion: "Just me" / "2-3 people" / "4-5 people"

If team → ask for GitHub usernames + roles via follow-up.

### Q2: Quality Preferences

AskUserQuestion (multiSelect):
- "Git hooks (pre-commit + pre-push)"
- "GitHub Actions CI"
- "Branch protection"
- "All of the above (Recommended)"

### Q3: Development Isolation (if team > 1)

AskUserQuestion: "Branch only (Recommended)" / "Git worktree"

If worktree → `git config --local wing.worktree true`

## Step 3: Generate Config

Create `$WING_DIR/config.yml`:

```yaml
schema_version: 3

team:
  - github: {GH_USER}
    role: {role}

roles:
  - id: tech-lead
    label: "Tech Lead"
    level: leader
  - id: developer
    label: "Developer"
    level: member

project:
  language: {detected}
  test_command: "{detected}"
  lint_command: "{detected}"
  {If monorepo:}
  services:
    - name: {dir}
      language: {lang}
      test_command: "{cmd}"

conventions:
  base_branch: "{detected or main}"
  production_branch: "main"
  branch_pattern: "mission/{issue}-{slug}-{user}"

label_prefix:
  status: "status:"
  priority: "priority:"
  mission: "mission"

quality:
  ci: {true/false}
  review_required: {true/false}
  hooks: {true/false}

versions:
  current: ""

notifications:
  slack: ""
  feishu: ""
```

### Merge Mode (re-init)

If config already exists → preserve existing values, add only new fields. Show diff.

## Step 4: Generate CI Workflow (if selected)

Create `.github/workflows/ci.yml`:

```yaml
name: CI
on:
  pull_request:
    branches: [{base_branch}]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - {language-specific setup}
      - run: {test_command}
      - run: {lint_command}
```

## Step 5: Generate Git Hooks (if selected)

Create `.githooks/pre-commit` and `.githooks/pre-push`. Configure: `git config core.hooksPath .githooks`

## Step 6: Setup GitHub Labels

```bash
bash "${CLAUDE_SKILL_DIR}/../scripts/setup-github-labels.sh"
```

Creates: `mission`, `status:wip`, `status:review`, `status:done`, `priority:P0`..`P3`, `size:S`..`XL`

## Step 7: Branch Protection (if selected)

Guide user through branch protection setup for base branch (requires admin access).

## Step 8: iTerm2 Tab Integration

Install shell hook for branch-aware tab titles + type-based tab colors:
- 🟥 main/master · 🟦 develop · 🟧 rc/* · 🟩 mission/* · ⬜ other

Requires `CLAUDE_CODE_DISABLE_TERMINAL_TITLE=1` in `~/.claude/settings.json`.

## Step 9: Commit + Summary

```bash
git add $WING_DIR/ .github/ .githooks/ 2>/dev/null
git commit -m "feat(wing): initialize project configuration"
git push
```

**INITIALIZED** ── {repo} ────────────────────
Config: `$WING_DIR/config.yml`
Language: `{language}`  Test: `{test_cmd}`
Team: `{count}` members  Quality: {gates}

NEXT: `/wing-issue <description>` to create your first mission
