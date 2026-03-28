# Project Profile Schema

> Date: 2026-03-28
>
> Purpose: define the schema for project profiles so the orchestration harness stays project-agnostic while still being able to operate effectively in real repositories.

---

## Why This Exists

The harness must not secretly be one project's habits baked into generic logic.

Project-specific reality must live in a profile layer.

This is what allows:

- project agnosticism
- reusable harness logic
- repeatable onboarding to new repos

---

## What A Profile Is

A project profile is a structured description of how a specific repository should be operated by the harness.

It should contain:

- conventions
- verification commands
- release topology
- important entry paths
- policy defaults

It should not contain:

- generic protocol semantics
- generic control-plane behavior
- hardcoded assumptions about one runtime's UI

---

## Canonical Structure

Logical shape:

```yaml
version: 1
profile:
  ...
```

The adapter may choose the final file path and format.

---

## Required Fields

### Identity

- `profile_id`
- `name`
- `repository`
- `default_runtime`

### Repository Conventions

- `base_branch`
- `production_branch`
- `branch_pattern`
- `worktree_enabled_default`

### GitHub Semantics

- `issue_label_prefix`
- `status_labels`
- `priority_labels`
- `milestone_strategy`

### Verification

- `default_verify_commands`
- `targeted_verify_conventions`
- `required_pre_ship_checks`

### Release Topology

- `release_mode`
  Examples:
  `single-branch`, `dual-branch`, `rc-based`

- `rc_branch_pattern`
- `staging_workflow`
- `promote_workflow`

### Important Paths

- `docs_entrypoints`
- `code_entrypoints`
- `config_entrypoints`
- `legacy_paths`

### Policy Defaults

- `max_parallel_runs`
- `meeting_triggers`
- `high_risk_scopes`
- `cleanup_policies`

---

## Recommended Optional Fields

- `language_stack`
- `package_manager`
- `service_map`
- `deployment_targets`
- `forbidden_paths`
- `reviewer_defaults`
- `automation_defaults`

These are useful but not universally required.

---

## Example Profile

```yaml
version: 1
profile:
  profile_id: "erwin-main"
  name: "erwin"
  repository: "liyasong/erwin"
  default_runtime: "codex"
  base_branch: "main"
  production_branch: "main"
  branch_pattern: "mission/{issue}-{slug}"
  worktree_enabled_default: true
  issue_label_prefix: "status:"
  status_labels:
    queued: "status:queued"
    wip: "status:wip"
    review: "status:review"
    done: "status:done"
    blocked: "status:blocked"
  priority_labels:
    p0: "priority:P0"
    p1: "priority:P1"
    p2: "priority:P2"
  milestone_strategy: "github-milestones"
  default_verify_commands:
    - "npm test"
  targeted_verify_conventions:
    - "prefer narrowest useful test first"
  required_pre_ship_checks:
    - "working tree clean or intentionally scoped"
    - "targeted verification passed"
  release_mode: "single-branch"
  rc_branch_pattern: "rc/{version}"
  staging_workflow: ""
  promote_workflow: ""
  docs_entrypoints:
    - "_project/CURRENT.md"
    - "AGENTS.md"
  code_entrypoints:
    - "showcase/src/App.tsx"
    - "skills/"
  config_entrypoints:
    - "package.json"
    - "showcase/package.json"
  legacy_paths:
    - "skills/teamwork/"
  max_parallel_runs: 2
  meeting_triggers:
    - "high-blast-radius design disagreement"
    - "conflicting high-authority evidence"
  high_risk_scopes:
    - "release flows"
    - "global scripts"
  cleanup_policies:
    - "report stale runs"
    - "report stale contracts"
```

---

## Invariants

1. Profiles must not redefine protocol meaning.
2. Profiles must not encode hidden legacy assumptions without marking them explicitly.
3. Profiles should declare legacy/deprecated paths so authority logic can treat them correctly.
4. Profiles should make verification and release expectations explicit.
5. A missing profile should degrade capability gracefully, not produce hallucinated defaults.

---

## Relationship To State Authority

Project profiles are usually L3 state:

- current intended operating conventions
- important enough to shape orchestration
- still subordinate to fresh L1 system evidence

For example:

- if the profile says `base_branch: main` but the live repo clearly uses `develop`, live evidence should force investigation

---

## Relationship To Meetings

Profiles can also shape meeting defaults:

- which role families are useful
- which paths are authoritative
- which paths are legacy
- what constitutes high-risk scope

This makes meetings more portable across projects without becoming generic mush.

---

## Relationship To Codex

In Codex, the profile should help the adapter answer:

- where to look first
- how to verify work
- which release model applies
- how much parallelism is safe
- which paths are likely traps or legacy

This is one of the main ways the harness becomes project-aware without becoming project-coupled.

---

## Suggested Next Step

After this schema, the next implementation step should be:

- deciding actual file locations
- defining profile discovery rules
- creating one real profile as a proving ground

That can happen without compromising project agnosticism.
