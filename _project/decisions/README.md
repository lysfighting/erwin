# Architecture Decision Records (ADRs)

This directory contains Architecture Decision Records for the erwin project.

## Format

Each ADR is named: `NNNN-short-title.md` (e.g., `0001-use-react-for-frontend.md`)

## Template

```markdown
# ADR-NNNN: Title

- **Date**: YYYY-MM-DD
- **Status**: Proposed | Accepted | Deprecated | Superseded by ADR-NNNN
- **Deciders**: (who made this decision)
- **Session**: docs/sessions/YYYY-MM-DD.md

## Context

What situation or problem led to this decision?

## Decision

What was decided?

## Consequences

What are the positive and negative consequences of this decision?

## Alternatives Considered

What other options were evaluated and why were they rejected?
```

## Index

| ADR | Title | Status | Date |
|-----|-------|--------|------|
| [0001](0001-monorepo-structure.md) | Monorepo structure for erwin | Accepted | 2026-03-02 |
| [0002](0002-react-vite-tailwind-stack.md) | React+Vite+Tailwind for drive-showcase | Accepted | 2026-03-02 |
| [0003](0003-v2-four-act-architecture.md) | V2 Four-Act interactive architecture | Accepted | 2026-03-02 |
| [0004](0004-d3-force-for-origin-trace.md) | d3-force for OriginTrace visualization | Accepted | 2026-03-02 |
| [0005](0005-no-automated-changelog.md) | Manual changelog over automated tooling | Accepted | 2026-03-02 |
| [0006](0006-new-project-skill.md) | /erwin-init skill for universal project bootstrap | Superseded by 0008 | 2026-03-02 |
| [0007](0007-skills-incubator-structure.md) | erwin as skills incubator — repository structure | Accepted | 2026-03-02 |
| [0008](0008-erwin-init-redesign.md) | /new-project → /erwin-init — 从脚手架到认知基础设施 | Accepted | 2026-03-02 |
