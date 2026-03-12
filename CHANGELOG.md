# Changelog

All notable changes to the **erwin** project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added — Skills Factory
- **Teamwork v3.9.2 — Tips `shuf` crash on macOS**: Tips instruction now requires Read tool instead of Bash. Prohibits `shuf`/`sort -R` (GNU-only, absent on Darwin).
- **Teamwork v3.9.1 — Self-drive audit fixes**: hotfix continue routing bug, VE compliance, manual footer v3.1.0→v3.9.1, RC Freeze clarification, team-ship RC hint
- **Teamwork v3.9.6 — Cross-platform sed + remaining fixes**: `_sed_i` helper for GNU/BSD portability, tw-pr.sh undefined var, team-rc cancel loop, team-ship stash -u inconsistency
- **Teamwork v3.9.5 — Script hardening (12 fixes) + docs staleness + ecosystem**: rebase ref fallback, jq injection x2, git add -A→-u, bash array safety, hotfix re-entry, ship subcommand guard, notify worktree, temp file trap. Docs: skills/README v2→v3.9.x, CLAUDE.md version, idea pipeline stages. install.sh copy-mode cleanup + --target guard.
- **Teamwork v3.9.4 — Worktree config bug + dashboard RC + 48 unit tests**: Fix `_detect_config()` worktree logic, add RC section to dashboard, expand tips 15→29, add unit test suite
- **Teamwork v3.9.3 — Tooling + VE accuracy + docs staleness**: Fix `tw-bump-version.sh` path bug (manual never updated), VE accuracy (Promote Help contradiction, Cancel Flow raw git), CURRENT.md + ROADMAP.md brought current
- **Teamwork v3.9.0 — RC complete lifecycle**: `/team-rc fix` (cherry-pick to RC) + `/team-rc hotfix` (post-release production fix + backport). Full docs sync.
- **Teamwork v3.8.7 — Docs completeness audit**: Documentation completeness fixes
- **Teamwork v3.8.6 — Documentation audit + version sync**: Fix 3 version pollution files, sync manual to 3.8.6, add sync/share to all docs
- **Teamwork v3.8.5 — Cross-MC knowledge sharing**: `/team sync` (standalone rebase), `/team share` (cross-MC comments), L1 intake in team-drive
- **Teamwork v3.8.4 — Dashboard formatting + worktree init fix + tab colors**: Bold issue numbers, worktree _detect_dir fix, tab color refresh, bump-version semver validation
- **Teamwork v3.8.3 — /team-rc cancel + upstream-first RC**: RC cancel command, upstream-first merge model, post-merge PR checkbox
- **Teamwork v3.8.2 — Script沉淀 + DRY + UX dead-end fixes**:
  - 6 new script subcommands: `build-branch-name`, `extract-issue-number`, `worktree-main-repo`, `worktree-find-by-branch` (tw-git.sh) + `resolve-member`, `list-members` (tw-config.sh)
  - Shared `_PYTHON_MEMBER_PARSER` in tw-config.sh — 4 functions DRY-ed (net -148 lines)
  - `tw-bump-version.sh` handles 3 more version string patterns (inline backtick, parens, bash var)
  - AskUserQuestion 4-option cap enforced across all skills (was exceeding limit)
  - 3 UX dead-end fixes: team-ship done error guide, team-rc polling progress, team-claim worktree recovery
  - Wrap mode checkbox gate: team-drive skips already-checked tasks
  - E2E test Phase 1b: 13 tests for new script subcommands
  - Manual scripts architecture section updated (15→21 tw-git, +6 tw-config subcommands)
- **Teamwork v3.8.1 — Per-user setup + sync-all-checkboxes + 13 recovery hints**:
  - `sync-all-checkboxes` in tw-contract.sh (fuzzy-match Contract↔Issue checkboxes)
  - Per-user setup extracted to `/team init` Step 5 (worktree, tab title, notifications)
  - 13 `/team doctor` recovery hints added across all skills (was: raw shell commands)
  - `team-mode` subcommand in tw-config.sh (solo/team detection)
- **Teamwork v3.8.0 — Desktop notifications in `/team init`**:
  - macOS notification hook (`~/.claude/notify.sh`) with iTerm2 bell + fallback chain
  - Registered on Notification, Stop, TeammateIdle events
  - Subagent filtering (skip noise from inner agents)
- **Teamwork v3.7.1 — Visual Encoding Standard + audit sweep + doctor expansion**:
  - `docs/visual-encoding-standard.md`: single source of truth for color design language (6 semantic roles, 3 blue-family colors, DO/DON'T rules, rendering pipeline reference)
  - 5-scout comprehensive audit: 11 issues found and fixed (visual encoding violations, raw `rm` suggestions, missing tips/next-step, stale docs)
  - ONBOARDING.md modernized: all code-block examples → visual-encoding-compliant formatted markdown + Advanced Features section
  - Visual Encoding reference embedded in all 5 sub-skill file headers (drift prevention)
  - Doctor expanded from 8 → 11 checks: GitHub auth, Contract freshness (hash comparison), test suite health
  - Teamwork manual: Visual Encoding section expanded with markdown formatting rules + link to spec doc
- **Teamwork v3.7.0 — Smart sync: dashboard behind-detection + AI conflict resolution**:
  - Dashboard PR lines show `⚠️ N behind {base}` when branch is behind base branch
  - NEXT box prioritizes sync recommendation when PR is behind
  - `/team-ship sync` auto-resolves rebase conflicts via AI (no more manual STOP)
  - Sync summary shows resolved files + flagged ambiguities
  - Auto force-push after successful rebase
- **Teamwork v3.6.1 — Terminal tab title + Tips system**:
  - `/team init`: new Step 5b asks user to configure terminal tab title (auto-shows branch name)
  - Auto-detects shell (zsh/bash/fish), appends appropriate hook, idempotent
  - Tips system: every `/team-*` command shows random tip after output (progressive feature discovery)
  - Tips pool in `scripts/tw-tips.txt` (15 tips)
  - `install.sh`: now installs `*.txt` files from scripts/ directory
- **Teamwork v3.6.0 — Manual Ops Handoff**:
  - `/team-drive`: classifies sub-tasks as code (AI) vs `🔧 MANUAL` (human)
  - Completion splits into "CODE COMPLETE" + "YOUR TURN" handoff
  - `/team-ship`: MANUAL tasks don't block shipping, added to PR description
- **Teamwork v3.5.0 — `/team config` + dashboard version + docs completeness**:
  - `/team config` subcommand: shows current config.yml content, worktree status, edit instructions
  - Dashboard footer shows `Teamwork v{version}` (read dynamically from config)
  - Help/Learn outputs updated with all subcommands and dynamic version
  - Learn Section 11: Config Management (shared vs local config, edit workflow)
- **Teamwork v3.4.0 — Doctor actionability upgrade**:
  - `/team doctor fix`: interactive fix mode with severity-based confirmations
  - `/team doctor help`: self-contained usage guide
  - Branch table: new "vs Remote" column showing ahead/behind remote tracking
  - Check 8: Remote Sync diagnostic
  - Cross-skill UX audit: 5 new help routes, 2 branch readiness checks, ~12 error→guidance fixes

### Changed — Skills Factory
- **Teamwork v3.4.1 — Worktree mode: shared → local**:
  - Worktree mode moved from shared `config.yml` to per-user `git config --local teamwork.worktree true`
  - Dashboard always shows worktree map when enabled (★ current, MC association, ⚠️ mismatch)
  - Init merge mode: auto-migrates old `worktree:` section to `git config --local`
  - Doctor Check 5: detects idle worktrees, context mismatch
- **Teamwork v3.2.0 — Unified visual design language**:
  - Standardized emoji encoding across all 6 skill files (priority dots, CI icons, progress bar)
  - Dashboard performance: 7-8 serial API calls → 4-5 parallel
- **Teamwork v3.1.0 — Untracked issue detection + learn audit**:
  - Dashboard UNTRACKED section: detects GitHub Issues assigned but missing teamwork labels
  - `/team-issue fix #N`: adds missing teamwork labels to untracked Issues
  - `/team learn`: complete 10-section design philosophy + visual diagrams

### Fixed — Skills Factory
- **Teamwork v3.5.1 — Design review: logic gaps + docs alignment**:
  - team-ship: base branch guard prevents accidental ship from main/develop
  - team-ship: CI failure preserves Contract (enables drive→ship iteration)
  - team-drive: Contract-missing fallback detects mission branch + open Issue, suggests recovery
  - docs: teamwork-v2-architecture.md archived, team-doctor-design.md Check 8 added
- **Teamwork v3.4.3 — Merge mode logic gaps**:
  - Init merge mode asks worktree preference when never configured
  - Quality section uses safe defaults instead of undefined Step 3 variables
- **Teamwork v3.4.2 — Worktree dashboard render reliability**:
  - Worktree section restructured as mandatory pre-render (was skipped by LLM in conditional template)
- **Teamwork v3.2.1 — Script extraction + deep bug fix + stale cleanup**:
  - Milestone prefix matching backport from vi_agent (4 locations + batch gap)
  - Learn: direct embedded output, no file reading (fix exit code 1)

### Changed — Skills Factory (v2.x)
- **Teamwork v2.7.3 — Init audit + dynamic branch protection**:
  - team.md Step 5c: hardcoded branch protection → config-driven (`quality.ci` + `quality.review_required` → 4 combinations)
  - Added `allow_force_pushes: false` + `allow_deletions: false` to protection rules
  - vi_agent: created post-merge.yml workflow, updated skill_version 2.5.0 → 2.7.2, set merge settings (squash-only)
  - vi_agent: branch protection enabled (require PR, review=0, no CI, block force push)
  - 18 requirements documented for release branch lifecycle, deploy workflow, issue version tracking

### Fixed — Skills Factory
- **Teamwork v2.7.2 — vi_agent E2E Validation** (3 bugs found, 3 fixed):
  - 1 CRITICAL: `_branch_pattern()` `${p:-default}` with `{}` in default — bash closes parameter expansion at first unmatched `}`, garbling branch names on any project with custom branch_pattern
  - 1 MEDIUM: `check-freshness` with non-existent file exits 3 instead of 2 — `exit` in function kills `$(…)` subshell before `|| true` can fire; added pre-check + moved `|| true` outside subshell
  - 1 LOW: `protect-check` USAGE comment said exit 1, actual code exits 3
- **Teamwork v2.7.1 — Script Testing Audit** (21 bugs found, 14+4 fixed across 4 scripts):
  - 2 CRITICAL: tag check matched branches not tags, merge-to on wrong branch after checkout failure
  - 6 HIGH: silent ensure-base failure, sed injection, fetch kills script, grep metachar, commit-type false positives, toggle-task false success
  - 6 MEDIUM: missing grep anchor, stderr pollution, detached HEAD, label transition drop, unused PR base param, body validation message
  - 4 LOW (fixed in follow-up): label verify substring, worktree-remove return 0, log-since feedback, branch_pattern double lookup

### Changed — Skills Factory
- **Teamwork v2.7.0 — Script Extraction** (Von Neumann Mechanism-Content Separation):
  - 4 new scripts: tw-git.sh (13 subcommands), tw-label.sh (4), tw-contract.sh (8), tw-pr.sh (7)
  - 6 skill files refactored: deterministic operations → script calls, LLM judgment sections unchanged
  - Design: scripts handle execution, skill.md handles orchestration + judgment
  - ~64% of operations moved from LLM-interpreted prose to executable scripts

### Fixed — Skills Factory
- **Teamwork v2.6.1 — Falsification Audit** (Popper's method applied across all 6 skill files + scripts):
  - 6 CRITICAL: freshness check always false-positive, PyYAML dependency, branch protection 422, CHANGELOG shell injection, stale merge, wrong push branch
  - 15 HIGH: init force unimplemented, hardcoded label prefixes (5 files), macOS sed syntax, milestone quoting, undefined variables, missing --state flags, tag existence check
  - All skill files bumped to v2.6.1

### Changed — Skills Factory
- **Teamwork v2.6.0 — Dual-branch model + release safety gates**:
  - Dual-branch support: `conventions.release_branch` enables base/release separation (backward compat when same)
  - team-release: staging deploy gate, merge base→release, CHANGELOG auto-append
  - team-drive: branch safety check prevents code execution on protected branches
  - team-claim: explicit base_branch checkout before feature branch creation
  - team-issue: `--milestone` parameter for per-issue milestone override
  - team-ship: sync/done operations use config base_branch; done D3/D4 now fallback (Action-first)
  - team.md: post-merge Action (auto-label + checkbox + milestone check), repo merge settings, config-driven branch protection + CI templates
  - Zero project coupling in skill files — fully generic design
- **Teamwork v2.5.2 — Dashboard queue overflow fix**:
  - QUEUED section now sorted by priority (P0→P1→P2→P3→none) and truncated to top 5 with overflow count
  - New `/team queue` subcommand shows full queued issue list grouped by priority
  - Prevents dashboard from growing unbounded when many issues are queued

### Added — Skills Factory
- **Teamwork v2 implementation** — 4 GitHub-first skills replacing v1:
  - `/team` — init system (auto-detect project, generate config/CI/hooks/templates) + dashboard from live GitHub data
  - `/team-claim` — claim GitHub Issue → generate AI-enriched Mission Contract → create branch
  - `/team-drive` — execute mission from Contract with sub-task verification loop
  - `/team-ship` — push → create PR (auto-closes Issue) → watch CI → cleanup Contract
- **Teamwork v2 architecture** — GitHub-first redesign (ADR-0010). Issues as task source, ephemeral Mission Contracts as AI execution views, auto-generated Actions CI + git hooks + Issue templates. Full design at `docs/teamwork-v2-architecture.md`. Supersedes v1 (ADR-0009).
- **Teamwork v2.1 — vi-agent feature absorption** — 6 features from vi-agent's `.teamspace` system:
  - Prefixed labels (`status:wip`, `priority:P1`, `status:blocked`, `status:done`)
  - Claim comment visibility (team transparency via `gh issue comment`)
  - Optional worktree isolation (parallel development per mission)
  - Optional version management (version labels + dashboard progress)
  - `/team-ship done` (post-merge cleanup), `/team-ship review` (AI code review), `/team-ship sync` (rebase)
- **Teamwork v2.1 architect review** — 13 bug fixes:
  - 3 critical: worktree contract access, Issue template label mismatch, branch pattern hardcoding
  - 3 high: config format collision, missing config health check, branch name parsing
  - 5 medium: sed pattern consistency, git staging guidance, worktree detection logic
  - 2 flow bugs: Step 1b→1c routing, $BRANCH_PATTERN extraction

### Changed — Documentation
- **docs/ reorganization** — Created `docs/README.md` index (active/archived/research), moved 3 research docs to `docs/research/`, added archived status banners to 4 stale docs, fixed stale cross-references to teamwork-v2-architecture.md

### Removed — Skills Factory
- **Teamwork v1 skills** — `/teamwork`, `/team-get-mission`, `/team-complete` removed (replaced by v2)

### Added — Showcase
- **SkillBuilder MVP** (ACT IV) — 5-step form to create custom Claude Code skills: Meta, Execution Model, Phases, Rules, Tools. Real-time preview with collapsible panel. Download generates valid skill.md with substantive prompt engineering content. Full i18n (EN/ZH).
  - `showcase/src/lib/skillTemplate.ts` — types + template generation (~310 lines)
  - `showcase/src/components/SkillBuilder.tsx` — stepper form + preview + download (~527 lines)
  - ~80 new translation keys in `showcase/src/i18n/translations.ts`

### Changed — Showcase (Performance)
- **Bundle code splitting** — 627KB single chunk → 221KB main + vendor chunks + 20 lazy-loaded chunks. React.lazy() for 19 components (15 ACT III sections + 3 ACT II components + SkillBuilder). Vendor splitting (react, framer-motion) in vite.config.ts. Vite 500KB warning eliminated.

### Added — Showcase (Animation Enhancements)
- **ExecutionSimulator topology** — 4-node SVG topology diagram (Team Lead → STL → Worker-1/Worker-2) with animated message flow dots (green dispatch, blue return) synced to simulation playback.
- **SwarmTopologySection message flow** — 12 animated message dots flow along connection paths in 3 phases (dispatch → inner dispatch → return). Auto-plays on scroll into view. Click any node to focus its message paths (non-involved nodes/connections dim to 8% opacity). Click again or click background to reset.
- **Consistent color coding** — dispatch=#22c55e (green), return=#3b82f6 (blue) across both components.

### Fixed — Showcase (Code Health)
- **SwarmTopologySection** — fixed `react-hooks/set-state-in-effect` ESLint error by replacing `useEffect`+`setMessageKey` with `activeNodeId` in React key (same remount behavior, no effect).
- **HeroSection** — fixed `react-hooks/set-state-in-effect` error by extracting `HeroAnimation` inner component with `key={lang}` reset pattern. Fixed `exhaustive-deps` warning with `useMemo`.
- **Context files** — suppressed `react-refresh/only-export-components` false positives for hook exports.
- **Semantic HTML** — added `<header>`, `<main>`, `<section aria-label>` landmarks to App.tsx.

### Removed — Showcase (Dead Code)
- **DemoSection.tsx** — removed orphaned 24KB component (not imported in App.tsx).
- **9 dead i18n keys** — removed `demo.note1.*`, `demo.note2.*`, `demo.note3.*` from translations.ts (only used by removed DemoSection).

### Changed — Skills Factory
- **`/new-project` → `/erwin-init`** — Renamed and redesigned per ADR-0008: from scaffolding tool to cognitive infrastructure installer. Zero-question entry, optional GOALS.md conversation. `erwin-` prefix reserved for framework-core skills.

### Added — Skills Factory
- **`skills/`** — New: source of truth for all erwin Claude Code skills (ADR-0007)
- **`skills/erwin-init/`** — `/erwin-init` v2.0.0: initializes erwin cognitive framework for any project; now version-controlled in erwin repo with README + CHANGELOG
- **`skills/visual-test/`** — `/visual-test` v1.0.0: screenshot-driven visual bug detection; now version-controlled in erwin repo with README + CHANGELOG
- **`install.sh`** — Symlink-based skills installer: `./install.sh` → all skills globally available; `git pull` auto-updates installed skills
- **`.claude/commands/`** — Project-level incubation stage for skills under development
- **`docs/erwin-skills-incubator-architecture.md`** — Full architecture design for skills incubator pattern
- **`_project/decisions/0007-skills-incubator-structure.md`** — ADR-0007 documenting the incubator structure decision

### Changed — `_project/ideas/`
- `visual-test-skill.md` moved from `ideas/` root → `ideas/skills/` (correct pipeline location)

### Added — Meta / Cross-Project (previous sessions)
- **`_project/decisions/0006-new-project-skill.md`** — ADR-0006 documenting /new-project skill design decisions
- **`~/.claude/projects/-Users-liyasong-Project-collov-erwin/memory/MEMORY.md`** — Claude auto-memory initialized for erwin project

### Changed — `.claude/rules/`
- **`ai-workflow.md`** — 将 "Session Start 读 CURRENT.md" 改为原则驱动的触发条件。不再罗列场景，改为一个判断测试：「行动如果方向错了，撤销成本高吗？」高成本行动（写代码、做决策、动手前）必须先读；纯对话（解释概念、讨论方案）无需读。Added step 4 to Session End: update MEMORY.md.

---

## [0.5.0] - 2026-03-02 (V2 Phase 2 — OriginTrace)

### Added — `showcase/`
- **`OriginTrace`** — d3-force directed graph showing drive-skill's intellectual genealogy (ACT II, after PrincipleCards). 15 nodes: 5 drive-concept nodes (matching PrincipleCards group colors A-E) + 10 historical source nodes (OODA Loop, Cybernetics, Mission Command, Mythical Man-Month, CPM/PERT, Amdahl's Law, Information Theory, Agile Manifesto, 孙子兵法, Double Diamond). Click any node to open a detail panel showing the origin description and how it manifests in drive-skill. `motion.svg` whileInView fade-in, `motion.g` stagger entrance, edge highlighting on hover/select.

### Changed — `showcase/`
- **`App.tsx`** — OriginTrace deployed in ACT II after PrincipleCards (previously placeholder comment).

### Dependencies
- `d3-force` added to dependencies
- `@types/d3-force` added to devDependencies

### Rationale
OriginTrace completes the ACT II "理解" experience: PrincipleCards explain the WHAT of each principle, OriginTrace traces the WHY (intellectual lineage). Together they answer "what does drive-skill do?" and "where did these ideas come from?"

---

## [0.4.0] - 2026-03-02 (V2 Phase 1 — Four-Act Structure)

### Added — `showcase/`
- **`ExecutionSimulator`** — Before/after animated comparison component (ACT I). Sequential vs parallel execution shown with Framer Motion progress bars, "3x faster" callout, replay button.
- **`PrincipleCards`** — Five Core Principles component (ACT II). Each card shows Chinese/English name, summary, expandable detail, and historical "thought sources" (OODA Loop, Mission Command, Amdahl's Law, Shannon, Sun Tzu).
- **`ActIIIWrapper`** — Collapsible wrapper for 15 deep sections (ACT III). Collapsed by default, shows section index, expands with Framer Motion animation.

### Changed — `showcase/`
- **`App.tsx`** — Restructured from flat 18-section scroll into Four-Act layout (ACT I 震撼 / ACT II 理解 / ACT III 深度 / ACT IV 创造). DemoSection replaced by ExecutionSimulator. ClaudeMechanicsSection moved to ACT II.
- **`NavBar.tsx`** — Replaced 18 individual section links with 4 Act anchor links (ACT I-IV), with intersection-observer-based active tracking.
- **`i18n/translations.ts`** — Added nav translation keys for four acts (en + zh).

### Rationale
V2 Phase 1 implements the approved Four-Act architecture from `docs/drive-showcase-v2-architecture.md`. The new structure enables the "15-second shock → 5-minute understanding → create" journey.

---

## [0.3.0] - 2026-03-02

### Added
- **`_project/` — Project Intelligence Layer** — three-layer information architecture:
  - `GOALS.md` — Vision + quarterly OKRs (the north star, read at session start)
  - `CURRENT.md` — Living status document (moved from `docs/CURRENT_STATE.md`)
  - `ROADMAP.md` — Milestone plan with version targets
  - `ideas/` — Idea Garden with four maturity stages (Seed/Sprout/Plant/Fruit)
    - `INBOX.md` — Raw captures, processed weekly
    - `skills/` — Drive-skill-like skill ideas
    - `meta/` — Universal cross-project patterns
  - `decisions/` — ADRs (moved from `docs/decisions/`)
  - `sessions/` — Conversation logs (moved from `docs/sessions/`)
  - `issues/` — Lightweight issue tracking (ISS-NNNN naming)
  - `releases/` — Narrative release notes
- **`.claude/rules/`** — Modular Claude Code instruction files
  - `git-conventions.md` — Commit format, scope definitions
  - `ai-workflow.md` — Session protocol, what to do at start/end
- **`docs/project-management-architecture.md`** — Full design specification
- **`_project/ideas/meta/project-management-template.md`** — Seed idea for universal project template skill

### Changed
- **`reference-skills/` → `reference/`** — cleaner naming; signals read-only source material
- **`CLAUDE.md`** — Updated paths, added GOALS reference, cleaner structure
- **`_project/CURRENT.md`** — Updated to reflect new structure and next tasks

### Rationale
The docs/ directory was a mix of project management files (CURRENT_STATE.md, sessions/,
decisions/) and technical documentation (architecture docs, analysis). This creates confusion
when navigating as a project manager vs. as a developer. The `_project/` layer creates a
dedicated home for all project intelligence, while `docs/` becomes pure technical documentation.
Additionally, the user correctly identified that the entire `_project/` pattern is a universal
capability for AI-assisted projects — captured as a Sprout idea in `_project/ideas/meta/`.

---

## [0.2.0] - 2026-03-02

### Added
- **Version management system** — git init, CHANGELOG.md, CLAUDE.md, ADR structure, session logs
- **CLAUDE.md** — project bootstrap context for AI sessions (readable at conversation start)
- **docs/decisions/** — Architecture Decision Records (ADR) directory
- **docs/sessions/** — conversation artifact archive (designs, decisions, research per session)
- **docs/CURRENT_STATE.md** — living document of current project state for session bootstrapping
- **docs/drive-showcase-v2-architecture.md** — V2 interactive showcase architecture design

### Rationale
The project had no git history and relied entirely on Claude Code's session memory for context
continuity. This release establishes the knowledge infrastructure: every decision, artifact,
and design gets versioned and committed, so any new conversation can bootstrap full context
by reading CLAUDE.md + docs/CURRENT_STATE.md rather than depending on volatile AI memory.

---

## [0.1.0] - 2026-03-01

### Added
- **showcase** — React/Vite/TypeScript web app showcasing drive-skill
  - 18 sections covering all aspects of drive.md, self-drive.md, architect.md
  - Bilingual (EN/ZH) with LangContext
  - Framer Motion animations throughout
  - Annotation system
  - CHANGELOG.md for showcase itself
- **reference-skills/drive-skill/** — local copy of the drive-skill repository
  - commands/: drive.md, self-drive.md, architect.md
  - hooks/: notify.sh
  - sounds/: 5 bundled TTS mp3 files
  - messaging/: slack and telegram scripts
  - pantheon/: 1000+ greatest minds knowledge base
  - config/: settings templates
- **docs/** — project documentation directory
  - drive-skill-analysis.md
  - visual-enhancement-architecture.md

### Rationale
Initial state of the erwin project at the point of git initialization. The showcase
had reached v1.3.0 (full content parity with drive.md) and was ready for V2 visual upgrade.
