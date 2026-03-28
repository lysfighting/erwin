# Codex Home and Claim Surface Spec

> Date: 2026-03-28
>
> Purpose: define the first two Codex-native user-facing surfaces for the new orchestration harness so implementation can begin with the most structurally important entry points.

---

## Why Start Here

`home` and `claim` are the best first surfaces because together they answer:

- what is the current system state?
- what should happen next?
- how does work enter active execution?

If these two are wrong, the rest of the harness will drift.

---

## `home` Surface

### Mission

Present the current harness state in a way that helps the leader or user see:

- active runs
- blockers
- stale state
- next valid step
- whether a meeting or cleanup is needed

### Inputs

- live git/GitHub state
- run registry
- active signals
- project profile
- maintenance backlog

### Core responsibilities

- detect missing or stale runtime state
- detect branch/worktree mismatches
- detect blocked or drifting runs
- summarize current mission graph at a glance
- recommend the next valid action

### Output shape

Recommended sections:

- active runs
- blocked / waiting runs
- stale state alerts
- high-priority signals
- next action recommendation

### Must not do

- become a giant dashboard with every detail
- mutate system state unless running an explicit safe repair
- hide uncertainty

---

## `claim` Surface

### Mission

Convert a mission from GitHub/task state into active, governed execution state.

### Inputs

- selected issue/mission
- live git/GitHub state
- project profile
- contract template
- run template
- current run registry

### Core responsibilities

- validate mission executability
- detect ownership conflicts or active duplicate runs
- generate contract
- provision branch/worktree as needed
- register the new run
- emit initial summary and next step

### Required checks

- issue exists and is executable enough
- active conflicting run does not already own the mission
- branch/worktree naming follows profile
- initial verification path is known
- scope and success criteria are explicit enough to start

### Output shape

Recommended sections:

- claimed mission summary
- created branch/worktree
- contract location
- run id
- next expected action
- warnings or unmet assumptions

### Must not do

- silently claim vague or contradictory work
- create hidden local state without registry updates
- auto-expand scope just to make the issue executable

---

## Shared Design Rules

### 1. Thin skill, strong substrate

`home` and `claim` should mostly orchestrate:

- scripts/mechanisms
- file-state updates
- profile-driven behavior

### 2. Authority-aware behavior

Live repo and GitHub state outrank stale local files.

### 3. No prompt-only truth

Every important action should result in durable state changes where appropriate.

### 4. Explicit next step

Both surfaces should always end by making the next valid action obvious.

---

## Implementation Readiness

These surfaces become implementable once the harness has:

- profile discovery
- run registry path conventions
- contract template
- run template
- basic deterministic mechanisms for branch/worktree setup

Those prerequisites now largely exist at the spec level.
