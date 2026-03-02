# Changelog

All notable changes to the **erwin** project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed — `.claude/rules/`
- **`ai-workflow.md`** — 将 "Session Start 读 CURRENT.md" 改为原则驱动的触发条件。不再罗列场景，改为一个判断测试：「行动如果方向错了，撤销成本高吗？」高成本行动（写代码、做决策、动手前）必须先读；纯对话（解释概念、讨论方案）无需读。设计思想来自 drive-skill 的 Cardinal Rule（方向偏差是唯一值得停下来的理由）。

---

## [0.5.0] - 2026-03-02 (V2 Phase 2 — OriginTrace)

### Added — `drive-showcase/`
- **`OriginTrace`** — d3-force directed graph showing drive-skill's intellectual genealogy (ACT II, after PrincipleCards). 15 nodes: 5 drive-concept nodes (matching PrincipleCards group colors A-E) + 10 historical source nodes (OODA Loop, Cybernetics, Mission Command, Mythical Man-Month, CPM/PERT, Amdahl's Law, Information Theory, Agile Manifesto, 孙子兵法, Double Diamond). Click any node to open a detail panel showing the origin description and how it manifests in drive-skill. `motion.svg` whileInView fade-in, `motion.g` stagger entrance, edge highlighting on hover/select.

### Changed — `drive-showcase/`
- **`App.tsx`** — OriginTrace deployed in ACT II after PrincipleCards (previously placeholder comment).

### Dependencies
- `d3-force` added to dependencies
- `@types/d3-force` added to devDependencies

### Rationale
OriginTrace completes the ACT II "理解" experience: PrincipleCards explain the WHAT of each principle, OriginTrace traces the WHY (intellectual lineage). Together they answer "what does drive-skill do?" and "where did these ideas come from?"

---

## [0.4.0] - 2026-03-02 (V2 Phase 1 — Four-Act Structure)

### Added — `drive-showcase/`
- **`ExecutionSimulator`** — Before/after animated comparison component (ACT I). Sequential vs parallel execution shown with Framer Motion progress bars, "3x faster" callout, replay button.
- **`PrincipleCards`** — Five Core Principles component (ACT II). Each card shows Chinese/English name, summary, expandable detail, and historical "thought sources" (OODA Loop, Mission Command, Amdahl's Law, Shannon, Sun Tzu).
- **`ActIIIWrapper`** — Collapsible wrapper for 15 deep sections (ACT III). Collapsed by default, shows section index, expands with Framer Motion animation.

### Changed — `drive-showcase/`
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
- **drive-showcase** — React/Vite/TypeScript web app showcasing drive-skill
  - 18 sections covering all aspects of drive.md, self-drive.md, architect.md
  - Bilingual (EN/ZH) with LangContext
  - Framer Motion animations throughout
  - Annotation system
  - CHANGELOG.md for drive-showcase itself
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
Initial state of the erwin project at the point of git initialization. The drive-showcase
had reached v1.3.0 (full content parity with drive.md) and was ready for V2 visual upgrade.
