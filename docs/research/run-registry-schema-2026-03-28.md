# Run Registry Schema

> Date: 2026-03-28
>
> Purpose: define the durable schema for execution runs so the leader-led orchestration harness can track real work over time instead of relying on ephemeral chat memory.

---

## Why This Exists

The harness distinguishes between:

- mission
- run

A mission is the durable work unit.
A run is one concrete execution attempt.

The run registry exists to record:

- what is currently running
- who owns it
- where it is running
- what state it is in
- what happened to previous attempts

Without this, the leader cannot reconcile correctly.

---

## Schema Goals

The run registry should be:

- durable across sessions
- inspectable by leader, teammates, and automations
- compatible with git/GitHub-backed workflows
- neutral to specific projects
- lightweight enough to update frequently

It is not meant to replace GitHub as the task system of record.

---

## Canonical Structure

The registry may live in JSON, YAML, or TOML.
The exact file path is adapter-defined.

Logical structure:

```yaml
version: 1
updated_at: "2026-03-28T18:00:00Z"
runs:
  - ...
```

Each entry in `runs` represents one execution attempt.

---

## Required Run Fields

### Identity

- `run_id`
  Stable unique id for this run.

- `mission_id`
  Stable mission identifier.
  Typically derived from GitHub issue identity plus project context.

- `issue_number`
  The GitHub issue number if applicable.

### Ownership

- `leader_id`
  Which leader instance or coordination surface owns this run.

- `owner_type`
  Examples: `leader`, `teammate`, `subagent`, `automation`.

- `owner_id`
  The active execution owner.

### State

- `state`
  One of:
  `draft`, `ready`, `executing`, `paused`, `waiting`, `blocked`, `completed`, `abandoned`

- `status_reason`
  Short explanation for current state.

- `priority`
  Optional but recommended normalized priority.

- `risk_level`
  Optional but recommended normalized blast-radius / risk marker.

### Execution Mode

- `mode`
  One of:
  `solo`, `managed`, `parallel`, `meeting`

- `execution_surface`
  Examples:
  `foreground-skill`, `background-automation`, `worktree-thread`, `council-session`

### Runtime Location

- `branch`
  Active branch if any.

- `worktree_path`
  Absolute path when relevant.

- `session_ref`
  Runtime session/thread reference when available.

### Timing

- `created_at`
- `started_at`
- `updated_at`
- `last_reported_at`
- `completed_at`

Only `created_at` and `updated_at` are always required.

### Verification

- `verification_status`
  Examples:
  `unknown`, `pending`, `passed`, `failed`, `partial`

- `verification_refs`
  References to latest relevant checks, commands, or reports.

### Coordination

- `depends_on`
  List of mission ids or run ids that gate this run.

- `siblings`
  List of sibling mission/run references considered coordination-relevant.

- `active_signals`
  List of signal ids still relevant to this run.

### Summary

- `objective`
  Short restatement of the bounded task.

- `latest_summary`
  Current compact status summary.

- `next_expected_action`
  What should happen next if nothing changes.

---

## Recommended Optional Fields

- `contract_ref`
  Path or identity of the active contract.

- `pr_number`
  Open PR if any.

- `release_ref`
  RC / release identifier if release-coupled.

- `meeting_ref`
  Current or latest relevant meeting id.

- `assigned_scope`
  Explicit scope boundary for teammate or sub-agent runs.

- `blocked_by`
  Machine-readable blocker references.

- `notes`
  Small extra metadata, not a dumping ground.

---

## Example Entry

```yaml
- run_id: "run-20260328-claim-042-a"
  mission_id: "issue-42"
  issue_number: 42
  leader_id: "leader-main"
  owner_type: "teammate"
  owner_id: "worker-auth"
  state: "executing"
  status_reason: "Implementing auth refactor in isolated worktree"
  priority: "P1"
  risk_level: "medium"
  mode: "managed"
  execution_surface: "worktree-thread"
  branch: "mission/42-auth-refactor"
  worktree_path: "/abs/path/repo-wt-42-auth-refactor"
  session_ref: "thread-abc123"
  created_at: "2026-03-28T09:00:00Z"
  started_at: "2026-03-28T09:05:00Z"
  updated_at: "2026-03-28T09:40:00Z"
  last_reported_at: "2026-03-28T09:38:00Z"
  verification_status: "partial"
  verification_refs:
    - "npm test -- auth"
  depends_on:
    - "issue-39"
  siblings:
    - "issue-45"
  active_signals:
    - "sig-20260328-17"
  contract_ref: ".harness/contracts/issue-42.md"
  objective: "Refactor auth path without breaking callers"
  latest_summary: "Core refactor done, 2 failing tests remain"
  next_expected_action: "finish caller updates, rerun targeted tests"
```

---

## Registry-Level Invariants

1. `run_id` must be unique.
2. At most one active foreground run should exist per mission unless the system explicitly supports split-run mode.
3. Every non-`draft` run must have an owner.
4. Every active run should have a recent `updated_at` and `last_reported_at`.
5. A `completed` or `abandoned` run must not remain the active owner of a mission.
6. Registry entries must not claim state that live git/GitHub evidence disproves.

---

## Freshness Rules

Runs can become stale.

Recommended staleness signals:

- `updated_at` older than threshold
- missing `last_reported_at`
- branch/worktree no longer exists
- owner session missing
- active signals never resolved

These should feed maintenance automations.

---

## Relationship To Git / GitHub

The run registry is secondary to GitHub and git for mission identity and delivery state.

Use the registry to answer:

- what is actively being executed
- who owns it right now
- what the harness believes the current run state is

Do not use it to replace:

- issue truth
- PR truth
- commit truth

---

## Relationship To State Authority

Run registry entries are typically L2 state:

- high-value
- near-live
- but still subordinate to fresh L1 system-of-record evidence

So if the registry says a PR exists and GitHub says it does not, GitHub wins.

---

## Update Rules

The registry should be updated when:

- a run is created
- ownership changes
- state changes
- a meaningful report arrives
- verification status changes
- a meeting materially changes the plan
- a run completes or is abandoned

Small transient shell output should not trigger updates.

---

## Suggested Next Step

This schema should be paired with:

- a signal schema
- a project profile schema
- adapter-specific path and file conventions

Only then should implementation start.
