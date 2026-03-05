# ADR-0012: Git Branching Model for AI-Heavy Teams

**Status**: Draft
**Date**: 2026-03-05
**Related**: ADR-0010 (teamwork v2), ADR-0011 (closed-loop design)

## Context

vi_agent team: 6 developers, all using AI assistants. Commit frequency equivalent to 20-30 person team. GitHub Flow (single main branch) cannot keep main clean under this volume — staging verification always lags behind new PRs merging to main.

Session Z identified the core problem: main as both development trunk and production branch creates an unresolvable contradiction at AI commit frequencies.

## Decision

Adopt **Git Flow (simplified)** with two long-lived branches:

- **develop**: development trunk, receives all feature/fix PRs (squash merge)
- **main**: production line, receives only verified release squash merges + tags

Key design choices:

1. Feature branches from `develop` (not main) — main's cleanliness is a result, not a starting point
2. All merges are squash — AI commit noise compressed to meaningful units
3. Tags on main's squash commits — `git diff V0.1..V0.2` for disaster recovery
4. One release branch at a time — simplicity over parallelism
5. Hotfixes cherry-picked back to develop immediately
6. Branch protection on both main and develop

## Consequences

- Teamwork skills need `production_branch` config key
- `/team-release` needs two-phase rewrite (prepare + ship)
- `deploy.yml` changes from auto-deploy-on-main to manual workflow_dispatch
- `ci.yml` adds develop to trigger branches
- All PRs target develop by default

## Full Specification

See `docs/git-model-v0.1.0.md`
