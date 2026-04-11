# erwin Skills

Reusable Claude Code and Codex-native skills created and maintained by the erwin project.

If you want the human-readable overview of `Harness`, read [`docs/harness-explainer.md`](/Users/liyasong/Project/collov/erwin/docs/harness-explainer.md).

## Available Skills

| Skill | Description | Version |
|-------|-------------|---------|
| [/arch-review](arch-review/) | Multi-role debate review of architecture documents (M1-M8) | 0.1.0 |
| [/erwin-init](erwin-init/) | Bootstrap any project with full AI intelligence layer | 1.0.0 |
| [/visual-test](visual-test/) | Screenshot-driven visual bug detection and fixing | 1.0.0 |
| [/session-record](session-record/) | Record AI session logs with structured format | 1.0.0 |
| [`/harness`](harness/) | Leader-led Codex-native orchestration harness surface | 1.0.0 |
| [`/knowledge`](knowledge/) | Codex-native layered knowledge-system curation and proposal gating | 1.0.0 |
| [/team](teamwork/) | Init GitHub-first teamwork system + dashboard | 3.10.4 |
| [/team-issue](teamwork/) | Create GitHub Issues with AI-enriched context | 3.10.4 |
| [/team-claim](teamwork/) | Claim GitHub Issue → generate Contract → create branch | 3.10.4 |
| [/team-drive](teamwork/) | Execute mission from Contract with drive methodology | 3.10.4 |
| [/team-ship](teamwork/) | Push, create PR, check CI, cleanup Contract | 3.10.4 |
| [/team-rc](teamwork/) | Release candidate validation + version management | 3.10.4 |

## Install All Skills

```bash
# Clone this repository
cd erwin
./install.sh
```

Installs all skills globally via symlinks to `~/.claude/commands/`. Updates automatically with `git pull`.

For Codex-native `SKILL.md` skills:

```bash
./install-codex.sh
./install-codex.sh harness
./install-codex.sh knowledge
```

Installs Codex-native skills to `~/.codex/skills/`.

## Install Individual Skills

```bash
./install.sh erwin-init     # install only /erwin-init
./install.sh visual-test    # install only /visual-test
./install.sh teamwork       # install /team + /team-issue + /team-claim + /team-drive + /team-ship + /team-rc
```

## Uninstall

Remove symlinks from `~/.claude/commands/`:
```bash
rm ~/.claude/commands/erwin-init.md
rm ~/.claude/commands/visual-test.md
rm ~/.claude/commands/session-record.md
rm ~/.claude/commands/team.md ~/.claude/commands/team-issue.md ~/.claude/commands/team-claim.md
rm ~/.claude/commands/team-drive.md ~/.claude/commands/team-ship.md ~/.claude/commands/team-rc.md
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

## Codex-native Skills

Newer surfaces may use `SKILL.md` rather than the older single-file command markdown.

- `skills/harness/SKILL.md` is the first stable repository-local Codex-native harness surface.
- `skills/knowledge/SKILL.md` is the canonical Codex-native knowledge curation surface.
