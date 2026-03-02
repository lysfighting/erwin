# Changelog

All notable changes to the Drive Skill Showcase will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.3.0] - 2026-03-01

### Added
- **SelfDriveSection** — the complete /self-drive skill: perpetual loop flowchart (5 interactive steps), 3 parallel scout types with their full analysis dimensions, Impact×Effort proposal ranking matrix, 6 diversity categories with rotation rule, adaptive learning within sessions, Slack remote control architecture (Socket Mode → API polling → terminal fallback), proposal message format (toggle to view), and all 11 self-drive anti-patterns
- **ArchitectSection** — the complete /architect skill: 7-phase design process (Phase 0–6, all click-to-expand), core philosophy (4 principles), /drive vs /architect comparison card, Phase 3 12-dimension comparison table preview, Phase 5 adversarial review 4-step protocol, Phase 6 architecture document 8-section structure, and 7 process rules
- Nav entries: Self-Drive, Architect (between Anti-Patterns and Prompt Eng)
- Translation keys: selfdrive.*, architect.*, nav.self-drive, nav.architect

### Rationale
v1.2.0 achieved parity for drive.md's protocols. v1.3.0 adds the two companion skills:
/self-drive is the perpetual engine that wraps /drive in a continuous improvement loop —
its loop mechanics, scout types, and ranking logic are now fully visualized. /architect is
the design-first counterpart to /drive — its 7-phase process from problem understanding to
final architecture document is now fully represented.

## [1.2.0] - 2026-03-01

### Added
- **Deep Reasoning Protocol section** — the 6 thinking tools (Self-Dialectic, Formal Logic, Deductive Chains, Inductive Patterns, First Principles, Pantheon Invocation) + Reasoning Rhythm (FRAME→DIVERGE→ANALYZE→STRESS→CONVERGE→GROUND); was completely absent from showcase
- **Direction Validation Protocol section** — "Wrong direction at full speed = maximum waste"; 6 triggers, 5-step validation process, Direction vs. Execution error distinction, mid-execution heuristics; was completely absent
- **Full Scaling Mindset** in PhilosophySection — O(2^C) vs O(C·log C) comparison, 4 Scaling Questions, 5 Rule Obsession anti-patterns with scaling alternatives, when rules ARE appropriate
- **Relentless Rule + Execution Discipline** in ExecutionPhasesSection — "If you can think of it, you MUST do it", Checkpoint Protocol, staying effective in long sessions
- **Escalation Ladder (5 levels)** in ExecutionPhasesSection — try different approach → research → reduce scope → skip and continue → escalate (LAST RESORT)
- **8 Mission Complete Conditions** + Team Shutdown Protocol in ExecutionPhasesSection

### Rationale
v1.1.0 achieved initial content parity for the most visible gaps. v1.2.0 adds the deep protocols
that lived only in drive.md: the reasoning framework that drives all major decisions, the direction
validation system that prevents wrong-direction sprints, and the full execution discipline that
separates "it kinda works" from "mission complete."

## [1.1.0] - 2026-03-01

### Added
- **CHANGELOG.md** — version management for the showcase itself (this file)
- **Self-Adversarial Review section** — the 6-axis review protocol that was completely absent; explains why you can't mark any task "done" without adversarial self-critique
- **Anti-Patterns section** — the full 30+ anti-patterns from drive.md, visualized and categorized; previously only mentioned as "30+" with no detail
- **WHY narrative to PantheonSection** — complete redesign explaining the activation mapping problem, why 1000+ figures exist, and the on-demand loading design
- **Phase 0 Interrogation details** in ExecutionPhasesSection — the 8-dimension interrogation framework, multi-round strategy, and the Psychiatrist Test
- **Team Lead Principles** (all 11) added to SwarmTopologySection
- **Direction Validation Protocol** — the "wrong direction = maximum waste" protocol with exhaustive path enumeration
- **Execution Discipline / The Relentless Rule** — added to relevant sections

### Changed
- **PantheonSection** — complete redesign: era timeline replaced by activation-mapping-first narrative
- **ExecutionPhasesSection** — expanded Phase 0 from 4 bullet points to full interrogation protocol
- **SwarmTopologySection** — added Team Lead Principles block
- **App.tsx** — added Self-Adversarial Review and Anti-Patterns sections to navigation

### Rationale
Previous version (1.0.0) presented a summarized subset of drive.md content. This release
achieves full content parity: every concept, rule, protocol, and anti-pattern from
the source skill files is now represented in the showcase. Graphical presentation
enhances rather than replaces the source material.

## [1.0.0] - 2026-02-XX

### Added
- **Initial showcase** — React/Vite/TypeScript with Framer Motion animations
- **HeroSection** — animated typewriter + title reveal
- **DemoSection** — interactive drive mode terminal demonstration
- **ClaudeMechanicsSection** — 8 Claude Code mechanisms with problem/solution/code pattern
- **ArchitectureSection** — repository structure with WHY for each module
- **ExecutionPhasesSection** — 5 phases with key insights and "without this" blocks
- **SwarmTopologySection** — animated SVG hierarchy visualization with live broadcast feed
- **WaveExecutionSection** — Gantt-style wave diagram and 3 concurrency patterns
- **SixStepSection** — interactive circular diagram with "without this step" blocks
- **PhilosophySection** — 5 axioms + 8 core beliefs + Bitter Lesson paradox
- **PantheonSection** — era timeline and problem→thinker activation mapping (initial)
- **PromptEngineeringSection** — 10 prompt engineering techniques with examples
- **ReferencesSection** — source references and attribution
- **Bilingual support** — Chinese/English language toggle
- **Annotation system** — click-to-annotate on key concepts
- **Dark theme** — GitHub-inspired dark color palette
- **NavBar** — sticky navigation with section anchors
