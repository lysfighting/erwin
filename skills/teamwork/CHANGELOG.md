# Changelog — Teamwork Skills

All notable changes to the Teamwork skill set.

## [Unreleased]

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
