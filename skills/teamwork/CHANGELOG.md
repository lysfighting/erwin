# Changelog — Teamwork Skills

All notable changes to the Teamwork skill set.

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
