# ADR-0007: erwin as Skills Incubator — Repository Structure

- **Date**: 2026-03-02
- **Status**: Accepted
- **Session**: `_project/sessions/2026-03-02.md`

---

## Context

erwin produced two Claude Code skills (`/erwin-init`, `/visual-test`) but their source files lived only in `~/.claude/commands/` — outside the git repository, with no version control, no update mechanism, and no documentation. This contradicted erwin's goal of being a "skills incubator."

Additionally, `_project/ideas/skills/` was empty (the skills pipeline had no visible structure), and `visual-test-skill.md` was misplaced at `ideas/` root instead of `ideas/skills/`.

## Decision

Formalize erwin as a skills incubator with the following structure:

1. **`skills/` directory** (new) — source of truth for all released skills
   - Each skill gets its own subdirectory: `skills/{name}/`
   - Contains: `{name}.md` (skill), `README.md` (docs), `CHANGELOG.md` (version history)
   - Skills have independent semantic versioning

2. **`install.sh`** (new) — symlink-based installer
   - `ln -sf skills/{name}/{name}.md ~/.claude/commands/{name}.md`
   - `git pull` auto-updates installed skills (symlinks point to live files)
   - Individual or all-skills install

3. **Two-stage incubation pipeline**
   - `_project/ideas/skills/` — design docs (Seed → Plant stages)
   - `.claude/commands/` — project-level incubation (Fruit-in-progress)
   - `skills/{name}/` — globally released (Fruit 🍎)

4. **Reorganized `_project/ideas/skills/`** — visual-test-skill.md moved here

## Consequences

**Positive**:
- Skills are now version-controlled in the erwin git repo
- `git pull` updates installed skills automatically (no reinstall)
- Clear pipeline from idea to global install
- Mirrors the proven drive-skill pattern exactly
- Each skill has its own README and CHANGELOG — self-documenting
- Anyone can install erwin skills with `git clone && ./install.sh`

**Negative**:
- Existing symlinks in `~/.claude/commands/` need to be updated (re-run `./install.sh`)
- The skill file exists in two locations briefly (skills/ + ~/.claude/commands/) — the symlink makes them the same file

## Alternatives Rejected

- **Separate repos per skill**: Too much overhead for markdown files
- **npm-based distribution**: Inappropriate for Claude Code markdown skills
- **Project-level commands only**: Skills wouldn't be globally available
- **No change**: Violates the skills incubator goal

## Architecture Document

Full design analysis in `docs/erwin-skills-incubator-architecture.md`.
