---
# Git conventions for erwin project
---

# Git Commit Conventions

## Format

```
type(scope): description

[optional body]

[optional footer: Closes ISS-NNNN | See ADR-NNNN]
```

## Types

- `feat` — new feature or capability
- `fix` — bug fix
- `docs` — documentation only (including _project/ updates)
- `refactor` — code refactor, no behavior change
- `chore` — maintenance (deps, gitignore, config)
- `style` — formatting only

## Scopes

- `showcase` — changes in drive-showcase/
- `project` — changes in _project/ (goals, ideas, decisions, sessions)
- `docs` — changes in docs/
- `reference` — changes in reference/
- `root` — root-level files (CLAUDE.md, CHANGELOG.md, .gitignore)

## Examples

```
feat(showcase): add ExecutionSimulator MVP with 3 scenarios
fix(showcase): fix Framer Motion stutter on Safari
docs(project): add ADR-0006 for simulation engine approach
chore(root): update .gitignore to exclude dist/
feat(project): add GOALS.md with Q1 OKRs
```

## Closing Issues

Use `Closes ISS-NNNN` in commit footer to link issue to commit:
```
feat(showcase): implement SkillBuilder download button

Closes ISS-0003
See ADR-0007
```
