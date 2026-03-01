# ADR-0003: V2 Four-Act Interactive Architecture

- **Date**: 2026-03-02
- **Status**: Accepted
- **Session**: docs/sessions/2026-03-02.md

## Context

The existing drive-showcase (v1.3.0) has 18 linear sections with good content coverage but:
- Static content, limited interactivity
- No clear user journey (beginners get lost)
- Design philosophy origins not explained
- No path for users to create their own skill

The project goal: "傻瓜进来看一眼就知道这个 skill 有什么用，然后消化吸收创建自己的版本"

## Decision

Restructure into a Four-Act experience while preserving all existing 18 sections:

- **ACT I (震撼, 15s)**: ExecutionSimulator showing before/after contrast
- **ACT II (理解, 5min)**: 5 PrincipleCards with embedded interactive demos + OriginTrace philosophy map
- **ACT III (深度, optional)**: All 18 existing sections as deep reference material (collapsible)
- **ACT IV (创造)**: SkillBuilder — interactive tool to generate custom skill.md

## Consequences

**Positive**:
- Clear user journey: shock → understand → dive deep → create
- Beginners can stop after ACT I or ACT II
- Existing 18 sections preserved (no content loss)
- Distinct value proposition: creators leave with a tangible artifact (their skill.md)

**Negative**:
- 3 new core components required (ExecutionSimulator, OriginTrace, SkillBuilder)
- Information architecture changes require refactoring App.tsx
- 3 weeks estimated development

## Alternatives Considered

- **Path A (incremental)**: Add components to existing structure. Rejected: doesn't fix core "user gets lost" problem.
- **Path B (explorable essay)**: 60-70% rewrite. Rejected: too much risk, discards existing work.
- **Path D (dual-track)**: Beginner/Creator toggle. Rejected: doubles maintenance burden.
- **Path E (gamified sandbox)**: Full simulator game. Rejected: 4-6 weeks, too much scope.

## Full Design

See `docs/drive-showcase-v2-architecture.md` for complete specification.
