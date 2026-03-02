# erwin Skills

Claude Code skills created and maintained by the erwin project.

## Available Skills

| Skill | Description | Version |
|-------|-------------|---------|
| [/new-project](new-project/) | Bootstrap any project with full AI intelligence layer | 1.0.0 |
| [/visual-test](visual-test/) | Screenshot-driven visual bug detection and fixing | 1.0.0 |

## Install All Skills

```bash
git clone https://github.com/your-org/erwin
cd erwin
./install.sh
```

Installs all skills globally via symlinks to `~/.claude/commands/`. Updates automatically with `git pull`.

## Install Individual Skills

```bash
./install.sh new-project    # install only /new-project
./install.sh visual-test    # install only /visual-test
```

## Uninstall

Remove symlinks from `~/.claude/commands/`:
```bash
rm ~/.claude/commands/new-project.md
rm ~/.claude/commands/visual-test.md
```

## Skill Pipeline (from idea to release)

```
_project/ideas/INBOX.md        ← raw ideas (Seed 🌱)
        ↓
_project/ideas/skills/          ← design docs (Sprout→Plant stages)
        ↓
.claude/commands/               ← incubation (Fruit-in-progress, project-level only)
        ↓
skills/{name}/                  ← released (Fruit 🍎, globally installable)
        ↓
~/.claude/commands/{name}.md    ← installed (symlink)
```

## Contributing a New Skill

1. Add idea to `_project/ideas/INBOX.md`
2. When worth exploring: create `_project/ideas/skills/{idea-name}.md`
3. When spec is complete (Plant stage): move to `.claude/commands/` for testing
4. When working + documented (Fruit stage): create `skills/{name}/` directory with:
   - `{name}.md` — the skill file (what Claude Code loads)
   - `README.md` — usage, requirements, examples
   - `CHANGELOG.md` — version history
5. Run `./install.sh {name}` to install globally
6. Update this README's skill table

## Skill File Format

Each skill `.md` file is a Claude Code command. The optional frontmatter:
```yaml
---
description: One-line description shown in skill picker
---
```

The skill content is instructions Claude follows when the user invokes `/skill-name`.
