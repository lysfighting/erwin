# Idea: Universal Project Management Template (Meta-Skill)

- **Stage**: Fruit 🍎 — erwin-init v2.0.0 fully implemented
- **Created**: 2026-03-02
- **ADR**: ADR-0006 (`_project/decisions/0006-new-project-skill.md`)
- **Implemented**: `~/.claude/commands/erwin-init.md`

---

## Problem

Every AI-assisted project suffers the same problems:
1. Context loss between sessions (each conversation starts from zero)
2. Decision amnesia ("why did we build it this way?")
3. Vision drift (daily tasks obscure the north star)
4. Idea rot (good ideas get lost in conversation history)

The current `_project/` structure in erwin solves all four problems. But it took 3+ architect
sessions to design it. Other projects shouldn't have to reinvent this.

---

## Hypothesis

If we extract the `_project/` structure as a portable template + install script, any project
can bootstrap professional knowledge management in 5 minutes. This could become a Claude Code
skill (`/erwin-init`) that sets up the complete infrastructure.

---

## What the Template Would Include

```
_project/
├── GOALS.md               ← Template with prompts to fill in
├── CURRENT.md             ← Ready-to-use format
├── ROADMAP.md             ← Empty milestone template
├── ideas/
│   ├── README.md          ← Garden rules (generic)
│   └── INBOX.md           ← Empty inbox
├── decisions/
│   └── README.md          ← ADR template + index
├── sessions/
│   └── README.md          ← Session log template
├── issues/
│   └── README.md          ← Issue template
└── releases/
    └── README.md          ← Release note template

CLAUDE.md                  ← Generic project bootstrap template
CHANGELOG.md               ← Empty with instructions
.claude/rules/             ← Generic rules
```

---

## The Skill Design

A `/erwin-init` skill would:

1. Ask: project name, description, team size, tech stack
2. Create the `_project/` structure
3. Pre-populate GOALS.md with the user's description
4. Initialize CHANGELOG.md with v0.1.0
5. Create CLAUDE.md pointing to the new structure
6. Run `git init` and make first commit
7. Output: "Your project is ready. Start by filling in _project/GOALS.md"

---

## Key Questions

- Should this live in a separate repository (like drive-skill) or inside reference/?
- What level of customization should be supported vs. forcing a standard structure?
- Should there be a YAML config file for project metadata?

---

## References / Prior Art

- drive-skill's `install.sh` — the symlink model for distributing skills
- GitHub's template repositories feature
- `create-react-app`, `degit` — project scaffolding patterns
- `.devcontainer/` pattern — project configuration in version control
