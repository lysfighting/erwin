# Harness Documentation Map

> Status: Active
>
> Date: 2026-04-08
>
> Purpose: define the authoritative documentation system for `Harness` so future sessions can load the right design context without reconstructing architecture from scattered files or chat history.

---

## Why This Exists

`Harness` now has enough architecture, runtime behavior, and historical lineage that one document alone is not the whole system.

That does **not** mean future sessions should read everything blindly.

The right model is:

- one top-level architecture narrative
- a small set of canonical supporting documents
- a machine-readable contract
- targeted research/audit docs only when the work needs them

This file manages that set.

---

## Canonical Document Set

These are the core documents that define `Harness`.

### 0. Human Overview / Explainer

- [`docs/harness-explainer.md`](/Users/liyasong/Project/collov/erwin/docs/harness-explainer.md)

Use this for:

- the best readable overview of what `Harness` is
- the fastest way to understand what problem it solves
- a dual-layer explanation for both newcomer intuition and engineering reality
- practical workflow-oriented usage guidance before diving into deeper architecture docs

If a human wants the system explained before reading the full architecture, this should be the first document.

### 1. Top-Level Architecture

- [`docs/harness-architecture.md`](/Users/liyasong/Project/collov/erwin/docs/harness-architecture.md)

Use this for:

- complete system identity
- design goal and non-goals
- lineage and influences
- layer separation
- control-plane shape
- runtime model
- current maturity and remaining expansion areas

If a session reads only one human-readable design document, it should read this one.

### 2. Public Package Contract

- [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)

Use this for:

- durable package boundary
- required surfaces
- runtime guarantees
- distribution contract
- profile contract
- governance contract

This is the external contract.
It is narrower than the architecture doc on purpose.

### 3. Operating Model

- [`docs/harness-operating-model.md`](/Users/liyasong/Project/collov/erwin/docs/harness-operating-model.md)

Use this for:

- leader responsibilities
- teammate contract
- communication model
- council / meeting behavior
- issue-to-mission orchestration behavior
- what “leader-led” means operationally

This is the most important supporting design doc for anyone implementing the `leader + teammate + council` model.

### 4. Session Continuation / State Persistence

- [`docs/harness-session-handoff.md`](/Users/liyasong/Project/collov/erwin/docs/harness-session-handoff.md)

Use this for:

- fresh-session read order
- current runtime truth sources
- what to trust first
- how to interpret older research and session logs

This is the continuation entrypoint, not the whole architecture.

### 5. Distribution Strategy

- [`docs/harness-distribution-strategy.md`](/Users/liyasong/Project/collov/erwin/docs/harness-distribution-strategy.md)

Use this for:

- source bootstrap vs packaged binary vs runtime bundle
- preferred end-user path
- current non-goals for distribution

### 6. Learning Loop

- [`docs/harness-learning-loop.md`](/Users/liyasong/Project/collov/erwin/docs/harness-learning-loop.md)
- [`docs/harness-proposal-ledger.md`](/Users/liyasong/Project/collov/erwin/docs/harness-proposal-ledger.md)

Use these for:

- proposal-first Harness evolution
- how field signals and real usage become structured proposals
- disposition rules such as `watch`, `encode-as-eval`, `prototype`, and `adopt`
- the durable ledger of reviewed Harness evolution proposals

These documents define how `Harness` learns from external engineering signal and local operating experience without turning daily discovery into uncontrolled runtime churn.

### 7. Release Criteria

- [`docs/harness-release-criteria.md`](/Users/liyasong/Project/collov/erwin/docs/harness-release-criteria.md)

Use this for:

- the explicit `1.0.0` promotion bar
- contract freeze scope
- required packaging state
- required verification matrix

This document is the durable release/freeze anchor for the current adapter line.

### 8. Machine-Readable Package Truth

- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)

Use this for:

- protocol identity
- adapter identity and version
- package file list
- bootstrap contract
- profile contract
- capability surfaces
- smoke script declarations

If a human-readable doc and the manifest disagree about package contents, the manifest wins for machine-readable truth.

---

## Implementation Truth

Some truths are not primarily in docs.

### Runtime / Control Plane

- [`internal/harness/runtime.go`](/Users/liyasong/Project/collov/erwin/internal/harness/runtime.go)

This is the real deterministic control plane.
Use it to verify:

- ownership selectors and transfer semantics
- run authority audit
- mission and mission-graph authority
- council orchestration
- CLI subcommands

### Runtime Entry

- [`cmd/harness-state/main.go`](/Users/liyasong/Project/collov/erwin/cmd/harness-state/main.go)
- [`scripts/harness-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-runtime.sh)

Use these to verify how the Go runtime is actually invoked and how packaged binaries are preferred.

### Shell Surfaces

- [`scripts/`](/Users/liyasong/Project/collov/erwin/scripts)

These are stable entrypoints, but they are intentionally thin.
Do not mistake them for the primary architecture.

---

## Research and Audit Layer

Research remains important, but it is not the first thing every session should load.

### Always Useful Research

- [`docs/research/protocol-core-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/protocol-core-spec-2026-03-28.md)
- [`docs/research/control-plane-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/control-plane-spec-2026-03-28.md)
- [`docs/research/meeting-protocol-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/meeting-protocol-spec-2026-03-28.md)
- [`docs/research/state-authority-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/state-authority-spec-2026-03-28.md)
- [`docs/research/cybernetics-kubernetes-arch-review-synthesis-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/cybernetics-kubernetes-arch-review-synthesis-2026-03-28.md)

Use these to understand why `Harness` is shaped the way it is.

### Audit / Drift Layer

- [`docs/research/harness-gap-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-gap-audit-2026-03-29.md)
- [`docs/research/harness-research-drift-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-research-drift-audit-2026-03-29.md)

Use these to answer:

- which original semantics are absorbed
- which are partial
- where implementation drift is still possible

Do not use these as the final truth for naming or packaging.
Those were partially overtaken by later docs and ADRs.

### Historical Lineage

- [`docs/research/teamwork-vs-wing-analysis-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/teamwork-vs-wing-analysis-2026-03-28.md)
- [`docs/research/parity-matrix-teamwork-wing-protocol-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/parity-matrix-teamwork-wing-protocol-2026-03-28.md)
- [`_project/ideas/skills/tw2-architecture.md`](/Users/liyasong/Project/collov/erwin/_project/ideas/skills/tw2-architecture.md)

Use these when the question is about what was preserved, redesigned, or intentionally dropped from `teamwork` and `wing`.

---

## Fresh Session Read Paths

If the task involves radar signal, usage feedback, knowledge threads, proposal pressure, or cross-lane repo strategy, read [`docs/erwin-operating-model.md`](erwin-operating-model.md) before treating Harness as an isolated implementation system.

### Path A: Continue Implementation Safely

1. [`AGENTS.md`](/Users/liyasong/Project/collov/erwin/AGENTS.md)
2. [`_project/CURRENT.md`](/Users/liyasong/Project/collov/erwin/_project/CURRENT.md)
3. [`docs/erwin-operating-model.md`](erwin-operating-model.md) when the task crosses Radar / Feedback / Knowledge / Proposals / Implementation boundaries
4. this file
5. [`docs/harness-explainer.md`](/Users/liyasong/Project/collov/erwin/docs/harness-explainer.md)
6. [`docs/harness-architecture.md`](/Users/liyasong/Project/collov/erwin/docs/harness-architecture.md)
7. [`docs/harness-operating-model.md`](/Users/liyasong/Project/collov/erwin/docs/harness-operating-model.md)
8. [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
9. [`docs/harness-release-criteria.md`](/Users/liyasong/Project/collov/erwin/docs/harness-release-criteria.md)
10. [`docs/harness-learning-loop.md`](/Users/liyasong/Project/collov/erwin/docs/harness-learning-loop.md)
11. [`docs/proposals/harness-proposal-ledger.md`](/Users/liyasong/Project/collov/erwin/docs/proposals/harness-proposal-ledger.md)
12. [`docs/harness-session-handoff.md`](/Users/liyasong/Project/collov/erwin/docs/harness-session-handoff.md)
13. [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)
14. targeted research and recent session logs as needed

### Path B: Review / Drift Check

1. [`docs/harness-architecture.md`](/Users/liyasong/Project/collov/erwin/docs/harness-architecture.md)
2. [`docs/harness-operating-model.md`](/Users/liyasong/Project/collov/erwin/docs/harness-operating-model.md)
3. [`docs/research/harness-gap-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-gap-audit-2026-03-29.md)
4. [`docs/research/harness-research-drift-audit-2026-03-29.md`](/Users/liyasong/Project/collov/erwin/docs/research/harness-research-drift-audit-2026-03-29.md)
5. [`internal/harness/runtime.go`](/Users/liyasong/Project/collov/erwin/internal/harness/runtime.go)
6. recent Harness commits on `main` or the active task branch

### Path C: Packaging / Distribution Work

1. [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
2. [`docs/harness-distribution-strategy.md`](/Users/liyasong/Project/collov/erwin/docs/harness-distribution-strategy.md)
3. [`docs/harness-learning-loop.md`](/Users/liyasong/Project/collov/erwin/docs/harness-learning-loop.md)
4. [`docs/harness-proposal-ledger.md`](/Users/liyasong/Project/collov/erwin/docs/harness-proposal-ledger.md)
5. [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)
6. [`scripts/harness-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-runtime.sh)
7. [`scripts/harness-build-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-build-runtime.sh)
8. [`scripts/harness-package-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-package-runtime.sh)

---

## Governance Rules For Documentation

When `Harness` evolves, update docs in this order:

1. machine-readable truth if package/runtime shape changed
2. public contract if the durable package boundary changed
3. architecture doc if system shape or design interpretation changed
4. operating model if leader/teammate/council behavior changed
5. handoff / CURRENT if continuation guidance changed
6. audits only when reassessing drift or maturity

This prevents documentation from turning back into scattered notes.

---

## Bottom Line

`Harness` now has a documentation system, not just a pile of documents.

The rule is simple:

- architecture explains the whole
- operating model explains how the system behaves
- public contract defines the durable package boundary
- session handoff explains how to resume safely
- manifest carries the machine-readable truth
- research and audits support, but do not replace, the canonical docs
