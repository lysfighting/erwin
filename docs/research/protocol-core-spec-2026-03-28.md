# Autonomous Development Lifecycle Protocol

> Date: 2026-03-28
>
> Purpose: define the runtime-neutral core that currently exists in mixed form across `teamwork`, `wing`, and repository practice.

---

## Why This Exists

The current system is entangled across several layers:

- protocol semantics
- Anthropic / Claude runtime assumptions
- Codex runtime assumptions
- historical `teamwork` packaging
- future-facing `wing` positioning
- real-world shaping from `vi_agent`

This document extracts the part that should survive all harness migrations:

**the protocol core**

---

## Protocol Identity

This protocol is not fundamentally "team collaboration."

It is:

**an AI autonomous development lifecycle protocol for GitHub-native software work**

The human may participate heavily, lightly, or mostly as an approver.
The protocol still holds.

---

## Core Loop

The canonical loop is:

`issue -> claim -> drive -> ship -> rc`

Definitions:

1. `issue`
   A GitHub-native work item with enough standalone structure to be executable.

2. `claim`
   Convert an Issue into an execution-ready Contract, bind it to a branch/worktree, and establish responsibility.

3. `drive`
   Execute the Contract against real code with verification, plan correction, and iterative completion.

4. `ship`
   Deliver the completed work as a PR, handle checks/review/status transitions, and prepare merge.

5. `rc`
   Manage release-candidate preparation, verification, promotion, and hotfix/backport behavior.

---

## Non-Negotiable Invariants

These are the protocol rules that should remain true across runtimes.

### 1. GitHub is the system of record

- Issues = tasks
- PRs = delivery
- Actions / checks = quality gate
- Milestones = version / release grouping

No shadow issue database.

### 2. Mission Contract is the execution view

The Issue is the public collaborative artifact.
The Contract is the AI execution artifact.

### 3. Standalone Test

Both Issue and Contract must be usable by an executor who did not participate in the original conversation.

### 4. Sub-tasks are hypotheses, not orders

Executors must validate the plan against code reality and revise when needed.

### 5. Deterministic mechanism must be mechanized

Git operations, label transitions, Contract mutations, PR lifecycle operations, and config parsing belong in deterministic mechanisms, not in freeform LLM reasoning.

### 6. Quality must be enforced with feedback loops

The protocol assumes:

- verification
- retries
- recovery paths
- visible next actions

### 7. The protocol must remain config-driven

Branch conventions, release topology, label prefixes, test commands, services, and deploy workflows should not be hardcoded to one project.

### 8. Blast radius is a first-class concept

Risk should be surfaced separately from urgency and effort.

### 9. Completion is behavioral, not textual

A step is complete when its state transition and verification have actually happened, not when a summary says it has.

---

## Canonical Objects

### A. Mission Issue

Required semantic fields:

- Objective
- Priority
- Size
- Domain
- Success Criteria
- Sub-tasks
- Context & References
- Constraints
- Verification Method

The exact markdown formatting may vary by runtime, but the semantics should not.

### B. Mission Contract

Required semantic fields:

- issue number
- title
- branch / execution target
- priority
- objective
- success criteria
- sub-tasks
- context files
- verification command(s)
- freshness metadata
- execution notes

### C. State Labels

Minimal state machine:

- queued
- wip
- review
- done
- blocked

Adapter-level formatting of labels can vary.

### D. Release Candidate

Required concepts:

- active RC uniqueness
- base branch
- production branch
- next version derivation
- staging validation
- promote
- fix
- hotfix

---

## Lifecycle Semantics

### Issue

Responsibilities:

- frame the real problem
- detect duplicates
- identify blast radius
- provide standalone execution context
- assign ownership appropriately

### Claim

Responsibilities:

- fetch issue
- validate executability
- generate Contract
- confirm Contract before branch creation
- create execution branch/worktree
- mark the ownership transition

### Drive

Responsibilities:

- read Contract
- validate plan vs codebase
- classify manual vs executable tasks
- choose execution mode based on task profile
- implement
- verify
- update Contract and Issue state

### Ship

Responsibilities:

- ensure pre-flight integrity
- push branch
- create/update PR
- handle `Closes` vs `Refs`
- request review
- update labels/state
- perform cleanup after merge

### RC

Responsibilities:

- prepare release candidate
- verify uniqueness and readiness
- promote to production
- allow surgical fix
- allow production hotfix and backport

---

## Cognitive Rules

These belong to the protocol, not to a specific runtime.

1. Own the mission.
2. Challenge weak framing.
3. Prefer simpler true solutions over elaborate compliance.
4. Treat factual claims as hypotheses until verified.
5. Trace all ripples of a change.
6. Do not silently skip broken assumptions.
7. Do not fake completion.
8. Keep momentum by always exposing the next valid step.

---

## Runtime Adapter Responsibilities

These are explicitly **not** part of the protocol core.

Each runtime adapter should decide:

- packaging shape
- file naming conventions
- path conventions
- tool invocation syntax
- permission model
- worktree/session/tab integration
- multi-agent implementation mechanism
- automation scheduling surface
- plugin / skill discovery model

Examples:

- Claude Commands adapter
- Anthropic Skills adapter
- Codex Skills adapter
- future OpenHands adapter

---

## Project Profile Responsibilities

These also do not belong in the core protocol.

Project profiles should define:

- repo-specific branch names
- deploy workflow names
- service maps
- test/lint/build commands
- team roster / role defaults
- environment-specific operational rules

Examples:

- `vi_agent`
- `erwin`
- future generic starter template

---

## Explicit Anti-Baggage Rules

1. The protocol core must not assume one harness.
2. The protocol core must not assume one repository.
3. The protocol core must not assume one package layout.
4. The protocol core must not embed giant runtime help text.
5. The protocol core must not preserve historical scars as universal law unless they remain generically justified.

---

## What This Means For Naming

This document intentionally does not choose a product name.

It only defines the thing that needs a name.

Possible future naming split:

- protocol name
- runtime adapter names
- broader orchestration platform name

That separation should be decided after the protocol is explicit.

---

## Suggested Next Artifacts

1. `runtime-adapter-spec-claude.md`
2. `runtime-adapter-spec-anthropic-skills.md`
3. `runtime-adapter-spec-codex.md`
4. `project-profile-vi-agent.md`
5. `project-profile-erwin.md`
6. `parity-matrix-teamwork-wing-protocol.md`
