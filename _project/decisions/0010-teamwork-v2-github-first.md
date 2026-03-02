# ADR-0010: Teamwork v2 — GitHub-First Architecture

**Status**: Accepted
**Date**: 2026-03-02
**Supersedes**: ADR-0009 (teamwork v1 repo-centric design)

## Context

Teamwork v1 (ADR-0009) used a repo-centric approach: `_teamwork/` directory with per-person markdown files, hand-written mission contracts, and minimal GitHub integration. User feedback: "当前的 skill 不好" — lacking GitHub ecosystem integration, low visibility, no quality gates, too many manual steps.

Reference implementation `vi-agent/.claude` provided a more mature system with board.md kanban, PR workflow, code review agent, and worktree isolation.

## Decision

**GitHub-First with ephemeral AI execution views.**

1. **GitHub Issues = task source of truth** (not local markdown files)
2. **Mission Contract = AI execution view** — auto-generated from Issue, enriched with context, discarded after PR merge
3. **GitHub Actions + git hooks = quality gates** — auto-generated during `/team init`
4. **GitHub Milestones = version management** for team projects
5. **`gh` CLI = bridge** between GitHub and local operations

## Approach: 方案 D (GitHub → Local Contract auto-generation)

Four approaches were evaluated:
- A: Pure GitHub — too bare for AI execution
- B: GitHub + persistent contracts — sync drift risk
- C: GitHub + Issue templates only — 80% good but limited
- **D: GitHub → ephemeral Contract** — best of both worlds

Contract is NOT a copy of the Issue. It's an **enhanced projection**: Issue content + AI-discovered context files + test commands + execution notes.

## Consequences

**Positive:**
- GitHub ecosystem fully leveraged (Issues, PRs, Actions, Branch Protection, Milestones)
- Visibility through GitHub native UI (no custom board needed)
- Quality enforced at multiple layers (hooks local, Actions remote, branch protection policy)
- Universal: works with any GitHub-hosted project
- Simpler local footprint: 1 config file + ephemeral contracts

**Negative:**
- Hard dependency on `gh` CLI and GitHub
- Offline degraded experience (cached Contract only)
- Issue template adoption not guaranteed (fallback: AI parses freeform)

**Neutral:**
- 4 skills (same count as v1, different responsibilities)
- vi-agent patterns partially adopted (lifecycle loop, PR flow, merge count)

## File Layout

```
.teamwork/config.yml          ← team config (tracked)
.teamwork/active/MISSION-N.md ← ephemeral contracts (.gitignore'd)
.github/ISSUE_TEMPLATE/mission.yml
.github/workflows/ci.yml
.githooks/pre-commit, pre-push
```
