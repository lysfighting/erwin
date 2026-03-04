# Changelog

All notable changes to the **erwin** project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed — Skills Factory
- **Teamwork v2.6.0 — Dual-branch model + release safety gates**:
  - Dual-branch support: `conventions.release_branch` enables base/release separation (backward compat when same)
  - team-release: staging deploy gate, merge base→release, CHANGELOG auto-append
  - team-drive: branch safety check prevents code execution on protected branches
  - team-claim: explicit base_branch checkout before feature branch creation
  - team-issue: `--milestone` parameter for per-issue milestone override
  - team-ship: sync/done operations use config base_branch (no hardcoded "main")
  - team.md: config template, health check, branch protection all config-driven
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
