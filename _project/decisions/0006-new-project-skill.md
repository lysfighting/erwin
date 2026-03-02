# ADR-0006: /new-project Skill for Universal Project Bootstrap

- **Date**: 2026-03-02
- **Status**: Superseded by [ADR-0008](0008-erwin-init-redesign.md)
- **Idea origin**: `_project/ideas/meta/project-management-template.md`
- **Session**: `sessions/2026-03-02.md`

---

## Context

Every new project with AI assistance requires the same setup:
1. Context continuity infrastructure (`_project/` layer)
2. Session protocol rules (`.claude/rules/`)
3. `CLAUDE.md` for AI navigation
4. Claude auto-memory initialization (`MEMORY.md`)

This setup took 3+ architect sessions to design for erwin. Without a skill, every new project
reinvents the wheel — or starts with no knowledge management infrastructure at all.

---

## Decision

Build a `/new-project` Claude Code skill at `~/.claude/commands/new-project.md` that scaffolds
the complete project intelligence infrastructure in one invocation (~5 minutes).

The skill:
1. Gathers project metadata via `AskUserQuestion` (name, description, tech stack, language)
2. Creates the full `_project/` directory structure with templates
3. Creates `CLAUDE.md`, `.claude/rules/ai-workflow.md`, `.claude/rules/git-conventions.md`
4. Creates `CHANGELOG.md` (Keep a Changelog format)
5. Computes and initializes `~/.claude/projects/{hash}/memory/MEMORY.md` (local auto-memory)
6. Runs `git init` (if needed) and makes the first commit

---

## Key Decisions Made

**Q: Where does the skill live?**
A: `~/.claude/commands/new-project.md` — standard location for user-invocable Claude Code skills.
Not a separate repository. Template content is embedded directly in the skill file.

**Q: Customization level — force standard structure or allow configuration?**
A: Force the erwin-proven standard structure. Project-specific content (name, description,
tech stack) is parameterized via `AskUserQuestion`. No YAML config file. Interactive prompting
is simpler and sufficient at this scale.

**Q: Should there be a YAML config file for project metadata?**
A: No. The skill asks 4 questions interactively. A config file would add complexity without
benefit — the setup is a one-time operation, not a repeatable pipeline.

**Q: Auto-memory path computation?**
A: `hash=$(pwd | tr '/' '-')` — converts the working directory path to the Claude Code project
hash format (e.g., `/Users/x/my-project` → `-Users-x-my-project`). This matches what Claude
Code uses internally for `~/.claude/projects/{hash}/`.

---

## Consequences

**Positive:**
- Any new project gets professional AI-assisted knowledge management in ~5 minutes
- The erwin methodology becomes portable — self-propagating to any project
- Each new project starts with Claude auto-memory already initialized
- Consistent `_project/` structure across all projects enables workflow reuse
- The `ai-workflow.md` template carries the principle-based CURRENT.md trigger (from Session E)

**Negative:**
- Skill needs manual maintenance when the template evolves (no auto-sync)
- Auto-memory (`MEMORY.md`) is machine-local — not shared via git
- If user already has a `CLAUDE.md`, the skill must confirm before overwriting

## Template Provenance

Derived from erwin project structure (v0.5.0), designed across 3+ architect sessions.
The `.claude/rules/ai-workflow.md` template is the principle-based version from Session E
(2026-03-02): action-cost trigger instead of time-based "session start" trigger.
