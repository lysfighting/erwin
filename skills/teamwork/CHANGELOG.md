# Changelog — Teamwork Skills

All notable changes to the Teamwork skill set.

## [2.6.1] — 2026-03-04

### Fixed — Falsification Audit (Popper's method applied across all skill files)

**CRITICAL fixes:**
- **team.md + team-drive.md**: Freshness check completely broken — `ISSUE_TITLE`/`ISSUE_BODY` never extracted from JSON, `CONTRACT_PATH` never defined. Hash always computed on empty string → false "ISSUE UPDATED" warning on every `/team-drive`. Fixed: added `jq` extraction + proper loop over Contract files.
- **team.md**: Merge mode used `import yaml` (PyYAML) — not in Python stdlib. Fails on systems without PyYAML. Fixed: replaced with `grep -q "^${section}:"` (zero dependencies).
- **team.md**: Branch protection `gh api --field` sent JSON objects as strings → 422 Validation Error. Fixed: switched to `--input` with heredoc for proper JSON body.
- **team-release.md**: Shell injection in CHANGELOG Python script — release body with triple quotes breaks embedded `'''$ENTRY'''`. Fixed: write to temp file, pass as `sys.argv`.
- **team-release.md**: `git merge "$BASE_BRANCH"` used stale local copy without fetch → release could miss latest code. Fixed: `git fetch` before merge, use `origin/$BASE_BRANCH`.
- **team-release.md**: Version bump committed on `$RELEASE_BRANCH` but pushed to `$BASE_BRANCH` → wrong branch in dual-branch mode. Fixed: `git checkout "$BASE_BRANCH"` before config update commit.

**HIGH fixes:**
- **team.md**: `init` force mode promised in routing table but never implemented — `/team init` on existing project showed dashboard. Fixed: added `$ARGUMENTS == "init"` check in Step 1.
- **team.md**: Version display in merge summary showed `2.5.0` (stale from v2.5). Fixed → `2.6.1`.
- **team.md**: Issue template hardcoded `"status:queued"`. Fixed → `"{STATUS_PREFIX}queued"`.
- **team.md**: Git hooks `grep` matched commented-out config lines + no `head -1`. Fixed: `grep -v '^\s*#'` + `head -1`.
- **team.md**: `{REPO}` template syntax mixed with `$REPO` shell syntax in bash code. Fixed → all `$REPO`.
- **team.md**: Duplicate section ID `6c` (two sections). Fixed → second renamed `6d`.
- **team.md**: Fallback label creation hardcoded `"mission"` and `"status:*"`. Fixed → config-driven.
- **team-issue.md**: macOS BSD sed `\+` syntax fails — `sed 's/.*--milestone \+\([^ ]*\).*/\1/'`. Fixed → `\+` replaced with `  *` (portable).
- **team-issue.md**: `MILESTONE_FLAG="--milestone \"$MILESTONE\""` embeds literal escaped quotes → milestone assignment always fails. Fixed → `${MILESTONE:+--milestone "$MILESTONE"}` parameter expansion.
- **team-issue.md**: References nonexistent `team.repo` config key. Fixed → proper `mc_label` + `STATUS_PREFIX` reads.
- **team-ship.md**: All label operations hardcoded `status:` prefix. Fixed → reads `$STATUS_PREFIX` from config.
- **team-ship.md**: `$ISSUE_TITLE` undefined in PR creation → commit type detection always defaulted to `feat`. Fixed → extract from Contract frontmatter.
- **team-ship.md**: `gh pr list --head {branch}` without `--state open` could find closed PRs. Fixed → added `--state open`.
- **team-claim.md**: Step 6 hardcoded `"status:queued"` / `"status:wip"`. Fixed → config-driven `${STATUS_PREFIX}`.
- **team-release.md**: No check if tag already exists before `git tag`. Fixed → `git rev-parse` + `git ls-remote` guards.
- **team-release.md**: `repos/{REPO}/milestones` template syntax in bash code. Fixed → `repos/$REPO/milestones`.
- **setup-github-labels.sh**: Hardcoded `"mission"` label, ignoring `mc_label` config. Fixed → reads via `tw-config.sh` with fallback.
- **setup-github-labels.sh**: `grep '^\s*status:'` could match wrong config line. Fixed → prefers `tw-config.sh`, fallback grep excludes comments.

### Changed
- All skill files bumped to v2.6.1
- Post-merge Action regex (committed in previous session fix): expanded from `Close(s)` to all 9 GitHub closing keywords, added `parseInt` for proper number type, replaced milestone race condition workaround with `listForRepo` query

## [2.6.0] — 2026-03-04

### Added — Dual-branch model + release safety gates

- **Dual-branch support** — New `conventions.release_branch` config key enables base/release branch separation:
  - When `release_branch == base_branch` (default) → single-branch mode, backward compatible with v2.5
  - When different (e.g. base=develop, release=main) → dual-branch mode: `/team-release` merges base→release before tagging
- **`team-release.md`** v2.5.0 → v2.6.0:
  - **Staging deploy gate** (Step 4b): if `deploy.staging_workflow` configured, checks latest Action run status before tagging. Warns on failure, user can abort or proceed.
  - **Merge base→release** (Step 5): in dual-branch mode, merges base into release branch before creating tag. Stops on conflict.
  - **CHANGELOG append** (Step 7b): after `gh release create`, fetches release body and prepends to `CHANGELOG.md`. Non-fatal if fails.
  - Renumbered steps 5-11 → 6-12 to accommodate new steps
- **`team-drive.md`** v2.4.0 → v2.6.0:
  - **Branch safety check** (Step 2b): prevents code execution on protected branches (base_branch or release_branch). Must be on feature branch.
- **`team-claim.md`** v2.4.1 → v2.6.0:
  - Explicitly checks out `base_branch` and pulls before creating feature branch (was implicit before)
- **`team-issue.md`** v2.5.1 → v2.6.0:
  - `--milestone` parameter: override milestone assignment per-issue (e.g. `/team-issue --milestone V0.2 fix camera`)
  - Priority: `--milestone` arg > `versions.current` from config > none

### Changed — Config-driven branch references

- **`team-ship.md`** v2.4.0 → v2.6.0:
  - Sync operation, done operation, help text all use `$BASE_BRANCH` from config instead of hardcoded "main"
  - **Done operation D3/D4 → fallback pattern**: checks if post-merge Action already handled label updates and Issue closure before acting. Skips with `(handled by PR merge)` / `(handled by post-merge Action)` messages.
- **`team.md`** v2.5.2 → v2.6.0:
  - **New Step 4e: Post-merge Action** — generates `.github/workflows/post-merge.yml` that auto-runs on PR merge: extracts `Closes #N` from PR body, updates labels (wip/review→done), checks all Issue checkboxes, checks milestone completion and comments when 100%
  - **New Step 5b: Repo merge settings** — auto-delete head branches, squash merge only, PR_BODY commit message (preserves `Closes #N`)
  - **Step 5c: Dual-branch protection** — protects both base_branch AND release_branch in dual-branch mode
  - **CI workflow templates** — use `{BASE_BRANCH}` placeholder instead of hardcoded `main`; dual-branch mode includes both branches in push triggers
  - Config health check table now includes `release_branch`, `deploy.staging_workflow`, and post-merge Action
  - Fresh install template includes commented `release_branch` and `deploy` sections
  - Dashboard merge count uses config `base_branch`
  - All `skill_version` references updated to 2.6.0

### Design principle

All branch names, workflow names, and deployment behaviors are config-driven. Zero hardcoded project references in skill files. Generic design validated through vi_agent practice.

## [2.5.1] — 2026-03-04

### Added — `/team-issue` duplicate check

- **`team-issue.md`** v2.4.1 → v2.5.1 — New **Step 1b: Duplicate Check** before Issue body generation:
  - Extracts 2-3 key terms from analyzed title/description
  - Runs `gh issue list --search` against open Issues
  - AI evaluates results for **semantic similarity** (not string match)
  - If duplicates found → `AskUserQuestion` with options: create anyway, update existing, or cancel
  - If no match → proceeds silently (zero friction on happy path)
- Help text updated: "WHAT HAPPENS (create)" now shows 6 steps including duplicate search

## [2.5.0] — 2026-03-04

### Added — `/team-release` skill

- **`team-release.md`** — New skill: close a version milestone, create git tag, publish GitHub Release, and prepare the next version. Full 11-step flow:
  1. Pre-flight: detect config dir, read `versions.current` (required), `versions.spec_path` (optional)
  2. Milestone check: show open/closed Issue counts, warn if open issues remain
  3. AI audit: if `versions.spec_path` configured, compare spec Success Criteria vs merged PRs + closed Issues. Outputs ✅/⚠️/❌ per criterion
  4. User confirmation before any mutations
  5. `git tag {VERSION}` + `git push origin {VERSION}`
  6. `gh release create {VERSION} --generate-notes` — auto-generates notes from merged PRs
  7. Close GitHub Milestone
  8. Prompt for next version name
  9. Create next GitHub Milestone
  10. Update `versions.current` in config.yml
  11. Commit config change + push
- Tag format: uses `versions.current` value exactly (e.g. `"V0.1"` → tag `"V0.1"`)
- AI audit is conditional: only runs when `versions.spec_path` is set and file exists

### Changed — `/team init` idempotency fix

- **`team.md`** v2.4.1 → v2.5.0 — `Step 4b` now has two modes:
  - **MERGE MODE** (config exists): preserves all existing sections, updates `skill_version` in-place via sed, appends only missing sections (`project`, `conventions`, `label_prefix`, `quality`). Explicitly does NOT touch `team`/`members`, `github`, `worktree`, `versions`, `roles`, or any custom sections.
  - **FRESH INSTALL** (no config): writes full template (existing behavior, unchanged).
  - Previous behavior: unconditional overwrite. Running `/team init` on a project with customized config (e.g. vi_agent's nested `github.mc_label` schema) would silently destroy all customizations.
- **`team.md` config template** — added `spec_path` comment to `versions:` block, pointing to `/team-release` as consumer.

---

## [2.4.1] — 2026-03-03

### Changed — Dashboard Upgrade (team.md v2.3.2 → v2.4.1)

- **`team.md`** — Replaced all fragile `grep/sed` config chains in Step 6c with `tw-config.sh` calls. Four sed chains eliminated:
  - `MISSION_LABEL` — now via `tw-config.sh github.mc_label` (dual-schema: teamspace first, then flat `mc_label`, then `"mission"` default)
  - `STATUS_PREFIX` — now via `tw-config.sh label_prefix.status`
  - `PRIORITY_PREFIX` — now via `tw-config.sh label_prefix.priority`
  - `CURRENT_VERSION` — now via `tw-config.sh versions.current`
- **`team.md` config template** — `skill_version` in generated `.teamspace/config.yml` template updated from `2.1.0` to `2.4.1`. New installs get correct version.
- **`team.md` version** — Bumped from `2.3.2` to `2.4.1` (skips 2.4.0 — dashboard was deferred during stability hardening, now catches up).

### Fixed — E2E Validation on vi_agent

- **`scripts/tw-config.sh`** — Fixed inline YAML comment stripping. Values like `mc_label: "mission-contract"  # comment` were returning `mission-contract"  # comment` instead of `mission-contract`. Now correctly extracts quoted value, ignoring trailing inline comment.
- **`team-claim.md`** — mc_label resolution now handles both config schemas: tries `github.mc_label` (teamspace/vi_agent schema) first, then `mc_label` (teamwork v2 schema), then defaults to `"mission"`. Fixes silent fallback to wrong label on vi_agent.
- **`team-issue.md`** — Same mc_label dual-schema fix. Both `github.mc_label` and `mc_label` paths tried before falling back to default.

Discovered and fixed during E2E validation run on `../vi_agent`. Full skill flow (issue → claim → drive → ship) verified non-destructively (Issue #27, PR #28, all closed/deleted post-test).

---

## [2.4.0] — 2026-03-03

### Changed — Stability Hardening

- **`scripts/tw-config.sh`** — New unified config parser. Replaces 15+ fragile sed chains across all skill files. Handles colon-containing values (label prefixes like `status:`), quoted strings, and nested YAML. Usage: `bash ~/.claude/commands/scripts/tw-config.sh project.test_command ""`.
- **Placeholder normalization** — Single canonical naming convention across all 5 skills: `{issue}`, `{pr}`, `{user}`, `{branch}`, `{title}`, `{repo}`, `{milestone}`. Eliminates confusion between `{N}`, `{ISSUE_NUMBER}`, `{pr-number}`.
- **team-ship pseudocode implemented** — tech-lead extraction (reads team roster for `role: tech-lead`), commit-type detection (prefix-based: fix/refactor/test/docs/feat), and base_branch reading were pseudocode in v2.3.x. Now real bash.
- **Critical command verifications** — `gh api user` failure now fatal with clear message in all skills. `git checkout` verifies branch existence before switching. TEST_CMD empty check warns instead of silently passing.
- **Issue bidirectional sync in team-drive** — After each sub-task is checked off in the Contract, the corresponding checkbox in the GitHub Issue body is updated. Team members see real-time progress on GitHub without visiting the branch. Non-fatal: network errors warn and continue.
- **Version from Milestone (not config)** — team-claim reads Milestone from Issue JSON (`milestone.title`) instead of `config.yml versions.current`. Config fallback still works for backward compatibility. Removes need to keep `versions.current` in sync manually.
- **REPO detection via `gh repo view`** — team-issue no longer requires a `repo:` field in config.yml. Uses `gh repo view --json nameWithOwner` auto-detection instead (same approach as other skills).

---

## [2.3.3] — 2026-03-03

### Added — Script-Based Skill Architecture (Pilot)

- **`scripts/setup-github-labels.sh`** — Executable bash script that creates/updates all teamwork labels on GitHub. Reads label prefixes from `.teamwork/config.yml` automatically. Idempotent (uses `--force`). Run: `bash ~/.claude/commands/scripts/setup-github-labels.sh`
- **`install.sh` scripts support** — `install_skill()` now detects and installs `scripts/*.sh` files from each skill directory to `$TARGET_DIR/scripts/`. Symlink mode: symlinks auto-update with `git pull`. Copy mode: copies standalone.
- **`/team` label setup uses script** — Step 5a now calls the script first with inline fallback. Eliminates the fragile 10-command sequence the LLM had to interpret.
- Establishes the pattern: skill folder = `.md` (AI brain) + `scripts/` (deterministic bash). LLM runs scripts; scripts do the mechanical work.

---

## [2.3.2] — 2026-03-03

### Added — Version Management + Issue Format Standardization

- **Milestone-based version tracking** — If `versions.current` is set in config, `/team-issue` assigns the Issue to that GitHub Milestone. Dashboard now queries milestone progress (open/closed counts) and displays `Version X: N/M tasks done`.
- **`/team` dashboard version header** — Reads `versions.current` from config, queries GitHub Milestones API for progress, shows `Version {name}: {closed}/{total} tasks done ({pct}%)` at top.
- **Canonical mission-contract format** — `/team-issue` now generates a consistent structure: `### Objective` (new top section), `### Success Criteria` (canonical name), `### Context & References` (not "Context"). Field names documented as canonical so parsers don't drift.
- **`/team-claim` backward-compatible parsing** — Reads `Success Criteria` (canonical) and `Acceptance Criteria` (legacy), `Context & References` and `Context` (legacy). Populates `version:` Contract frontmatter from `versions.current` config.
- `team-issue`, `team-claim`, `team` bumped to v2.3.2. `team-ship` stays v2.3.1.

---

## [2.3.1] — 2026-03-03

### Fixed — Ship Status Lifecycle

- **`/team-ship` label verification** — After updating Issue label from `status:wip` to `status:review`, immediately verifies via `gh issue view --json labels`. Shows actual labels in output instead of assumed ✅. If label update failed, shows diagnostic warning with manual fix command.
- **`/team` dashboard REVIEW section** — Added `UNDER REVIEW (PR open, awaiting merge)` section between member status and QUEUED. Ships `status:review` issues are now visible on the dashboard instead of falling into a grey zone.
- Both skills bumped to v2.3.1.

---

## [2.3.0] — 2026-03-03

### Added — Issue Sync & Lifecycle Comments

- **Issue freshness detection** — Contract stores `issue_content_hash` (SHA256 of title+body at claim time). `/team-drive` checks hash on entry; if Issue content changed, shows current body and asks user to Update Contract / Continue / Abort. Uses content hash instead of `updatedAt` to avoid false positives from our own comments, label changes, and assignee changes.
- **Dashboard [UPDATED] badge** — `/team` dashboard marks claimed Issues with `[UPDATED]` if content hash differs from Contract.
- **Lifecycle comments** — 3 key events posted as Issue comments for team visibility:
  - `/team-claim`: "Claimed by @user" (already existed)
  - `/team-drive` completion: "All sub-tasks complete — ready for review"
  - `/team-ship` PR creation: "PR #N created — {url}"
- All freshness checks and comments are **non-fatal** — network failures warn but don't block execution.
- Backward compatible: pre-v2.3.0 Contracts without `issue_content_hash` skip freshness check silently.

### Changed

- All skills bumped to v2.3.0

---

## [2.2.0] — 2026-03-03

### Added

- **`/team-issue`** — New skill: turn natural language into standardized GitHub Issues (mission-contract format). AI-enriches description with title, priority, size, domain, sub-tasks, and relevant files from codebase scan. Preview before publish.

### Changed

- **`/team-claim`** — Removed inline `create "title"` subcommand. Issue creation now handled by dedicated `/team-issue` skill (single responsibility).
- All skills bumped to v2.2.0

---

## [2.1.1] — 2026-03-02

### Fixed — Architect Review (13 bugs)

- **CRITICAL**: Contract inaccessible from worktree (active/ gitignored → added copy-to-worktree step)
- **CRITICAL**: Issue template hardcoded "mission" label → now uses `{MISSION_LABEL}` from config
- **CRITICAL**: Branch pattern hardcoded `mission/...` → reads `conventions.branch_pattern` from config
- **HIGH**: `team:` key format collision (erwin=member list, vi-agent=metadata) → check `members:` first
- **HIGH**: Missing config health check → added Step 1c with non-blocking warning
- **HIGH**: Operation Done assumes `mission/` branch prefix → generic numeric extraction
- **MEDIUM**: Worktree `enabled` flag missing → section presence = enabled
- **MEDIUM**: sed greedy pattern `'s/.*'` → all normalized to `'s/^[^:]*://'`
- **MEDIUM**: `git add -A` in team-drive → changed to specific file staging
- **FIX**: `cp` command used `{N}` instead of `${ISSUE_NUMBER}`
- **FIX**: Step 1b→1c routing contradiction (found users bypassed health check)
- **FIX**: `$BRANCH_PATTERN` extraction command missing

---

## [2.1.0] — 2026-03-02

### Added — vi-agent Feature Absorption

Integrated 6 advantages from vi-agent's `.teamspace` system while keeping erwin's portability.

- **Prefixed labels** — `status:wip`, `priority:P1`, `status:blocked`, `status:queued`, `status:done` (queryable, no collisions)
- **Claim comment visibility** — `gh issue comment` posted when claiming an Issue for team transparency
- **Worktree isolation** (optional) — `worktree.enabled: true` in config; creates isolated working directories per mission
- **Version management** (optional) — `versions.current: "V1.0"` in config; version labels + dashboard progress tracking
- **`/team-ship done`** — Post-merge cleanup: close Issue, update labels to `status:done`, clean worktree, return to main
- **`/team-ship review`** — AI code review: fetch PR diff, analyze for correctness/security/architecture/quality/performance, publish `gh pr review`
- **`/team-ship sync`** — Rebase current branch on latest main
- **`.mission` file** — Worktree mode detection (Issue number stored in worktree root)

### Changed

- Labels upgraded from flat (`in-progress`, `P0`) to prefixed (`status:wip`, `priority:P0`)
- Issue template now auto-applies `status:queued` label
- Config schema extended with optional `label_prefix`, `worktree`, `versions` sections
- `team-ship` now has argument routing: `done`, `review`, `sync` subcommands
- Backward compatible: unprefixed labels still work for legacy configs

---

## [2.0.0] — 2026-03-02

### Changed — GitHub-first Redesign

Complete rewrite. GitHub Issues/PRs/Actions replace local file-based coordination.

**Breaking changes:**
- `/teamwork` → renamed to `/team`
- `/team-get-mission` → renamed to `/team-claim`
- `/team-complete` → renamed to `/team-ship`
- `_teamwork/` directory → replaced by `.teamwork/` + GitHub ecosystem
- Mission contracts now auto-generated from GitHub Issues (not manually written)
- Config format changed: `_teamwork/config.md` → `.teamwork/config.yml` (YAML)

### Added

- **`/team`** — Init system with project auto-detection + GitHub config (labels, CI, hooks, Issue template) + dashboard from live GitHub data
- **`/team-claim`** — Claim GitHub Issue → auto-generate AI-enriched Mission Contract → create branch → assign on GitHub
- **`/team-drive`** — Read Contract, execute sub-tasks with verify loop, update checkboxes, commit with Mission reference
- **`/team-ship`** — Pre-flight checks → push → create PR with "Closes #N" → watch CI → request review → cleanup Contract
- **GitHub Actions CI** — auto-generated workflow per language (Node, Python, Go, Rust, Java)
- **Git hooks** — pre-commit (lint) + pre-push (test) via `.githooks/`
- **GitHub Issue template** — structured mission template (`mission.yml`)
- **Priority labels** — P0/P1/P2/P3 + mission/in-progress/review
- **Branch protection** — optional, requires CI pass + PR review
- **AI-enriched Context Files** — Contract auto-discovers relevant files via Glob/Grep
- **Dashboard** — live team status from GitHub API (assigned Issues, open PRs, merge count)

### Removed

- `_teamwork/` directory structure (replaced by `.teamwork/` + GitHub)
- Manual mission contract authoring (now auto-generated from Issues)
- Per-person status files (replaced by GitHub Issue assignment)
- Linear merge protocol (replaced by PR-based workflow)

### Architecture

- Design doc: `docs/teamwork-v2-architecture.md`
- ADR: `_project/decisions/0010-teamwork-v2-github-first.md`
- Inspired by vi-agent lifecycle patterns + GitHub-native workflows

---

## [1.0.0] — 2026-03-02

### Added

- **`/teamwork`** — Init system with interactive team config + dashboard view
- **`/team-get-mission`** — Claim missions, auto-create branches, update status
- **`/team-drive`** — Execute missions wrapping /drive methodology with teamwork context
- **`/team-complete`** — Push, rebase, conflict resolution, merge, status close
- **`_teamwork/` directory structure** — config.md, missions/, members/, archive/
- **Per-person status files** — zero metadata merge conflicts by design
- **Identity detection** — automatic via `gh auth status`
- **Branch auto-creation** — `mission/{id}-{slug}-{username}` naming convention
- **Linear merge protocol** — first merger wins, last merger resolves all conflicts
