# /erwin-init Changelog

## [Unreleased]

## [2.0.0] — 2026-03-02

### Changed
- **BREAKING**: Renamed from `/new-project` to `/erwin-init` per ADR-0008
- **BREAKING**: Removed all 4 initial questions (project name, description, tech stack, AI language)
- Project name now auto-detected from `basename $(pwd)`
- Tech stack left as TBD in CLAUDE.md — discovered naturally by Claude during work
- Language preference inferred from Phase 2 conversation, not asked upfront

### Added
- **Phase 2: GOALS.md conversation** — optional guided dialogue to help articulate project vision, success criteria, and anti-goals
- Identity clarification: cognitive infrastructure installer, not project scaffolding tool

### Rationale
- Decision: `_project/decisions/0008-erwin-init-redesign.md`
- The intelligence layer is stack-agnostic — asking tech stack at init time served no purpose
- Zero-question entry removes friction and makes the skill usable for existing projects

---

## [1.0.0] — 2026-03-02

### Added
- Initial skill implementation
- 4-parameter interactive setup via `AskUserQuestion`:
  - Project name (kebab-case slug)
  - One-line description
  - Tech stack (React/TypeScript/Vite, Python/FastAPI, Node.js/TypeScript, or custom)
  - AI response language (Chinese, English, or Mixed)
- Creates complete `_project/` intelligence layer with templates
- Creates `CLAUDE.md` as AI navigation handbook
- Creates `.claude/rules/` with `ai-workflow.md` and `git-conventions.md`
- Creates `CHANGELOG.md` in Keep a Changelog format
- Auto-memory initialization at `~/.claude/projects/{hash}/memory/MEMORY.md`
  - Hash computed as `$(pwd | tr '/' '-')` — path-based, no collision risk
- Safety check: warns if CLAUDE.md or `_project/` already exist
- Git init + first commit (conventional commit format)
- Completion report with next steps

### Origin
- Idea: `_project/ideas/meta/project-management-template.md`
- Decision: `_project/decisions/0006-new-project-skill.md`
- Session: `_project/sessions/2026-03-02.md`
