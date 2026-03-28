# Harness File Layout and Discovery Spec

> Date: 2026-03-28
>
> Purpose: define the concrete repository layout and discovery rules for the new orchestration harness so committed definitions and local runtime state do not get mixed together.

---

## Executive Decision

The harness should use two repository-local roots with different purposes:

- `harness/` = committed definitions
- `.harness/` = local runtime state

This separation avoids two bad outcomes:

1. runtime state polluting versioned source control
2. architectural definitions hiding inside ephemeral local directories

---

## Why Two Roots

The system now has at least two very different kinds of artifacts.

### Committed, shareable artifacts

- project profile
- reusable templates
- committed policies
- durable design references
- optional seeded examples

These belong in version control.

### Local, ephemeral, near-live artifacts

- run registry
- signal log
- scheduler snapshots
- meeting scratch state
- maintenance backlog
- local overrides

These should usually not be versioned.

Therefore the split is structural, not aesthetic.

---

## Directory Model

### `harness/` — committed definition layer

Recommended layout:

```text
harness/
├── README.md
├── profiles/
│   └── erwin.yml
├── templates/
│   ├── contract.yml
│   ├── run.yml
│   ├── signal.yml
│   └── meeting.yml
└── policies/
    └── defaults.yml
```

Purpose:

- shared definitions
- portable examples
- repository-default harness configuration

### `.harness/` — local runtime state layer

Recommended layout:

```text
.harness/
├── contracts/
│   └── issue-42.yml
├── runs/
│   └── registry.yml
├── signals/
│   └── log.yml
├── meetings/
│   ├── active/
│   └── archive/
├── scheduler/
│   └── snapshots/
├── maintenance/
│   └── backlog.yml
├── cache/
└── profile.local.yml
```

Purpose:

- live execution state
- ephemeral reconciler data
- local overrides
- automation-owned artifacts

---

## Discovery Rules

The adapter should discover harness configuration in this order:

### 1. Explicit override

Examples:

- command argument
- environment variable
- skill argument

This is highest precedence.

### 2. Local override

`.harness/profile.local.yml`

Use for:

- per-user local experiments
- temporary runtime-specific adjustments
- non-shareable machine-specific settings

### 3. Repository default profile

`harness/profiles/<default>.yml`

This is the normal committed project profile.

### 4. Fallback inference

If no profile is found:

- infer minimally from live repo state
- mark missing values as unresolved
- do not hallucinate a full profile

This is important for project agnosticism and safety.

---

## Discovery Rules for Runtime State

The adapter should look for runtime state only under `.harness/`.

Recommended defaults:

- run registry: `.harness/runs/registry.yml`
- signal log: `.harness/signals/log.yml`
- contracts: `.harness/contracts/`
- meeting state: `.harness/meetings/`
- scheduler outputs: `.harness/scheduler/snapshots/`
- maintenance backlog: `.harness/maintenance/backlog.yml`

If these do not exist:

- create lazily when first needed
- do not assume they already exist

---

## Versioning Policy

### Commit to git

- `harness/README.md`
- `harness/profiles/*.yml`
- `harness/templates/*.yml`
- `harness/policies/*.yml`

### Do not commit by default

- `.harness/runs/`
- `.harness/signals/`
- `.harness/contracts/`
- `.harness/meetings/`
- `.harness/scheduler/`
- `.harness/maintenance/`
- `.harness/cache/`
- `.harness/profile.local.yml`

This keeps the repo clean while preserving portable definitions.

---

## Local Override Policy

`profile.local.yml` should only override things that are legitimately local.

Good uses:

- local worktree root
- machine-specific tool paths
- preferred automation toggles
- temporary experiment flags

Bad uses:

- changing protocol semantics
- silently changing release topology
- redefining project-level source-of-truth paths

Those belong in the committed profile.

---

## Interaction With State Authority

This layout supports the authority model cleanly:

- live git/GitHub state = L1
- `.harness/` accepted runtime state = L2
- `harness/` committed profiles and templates = L3

That hierarchy is intuitive and inspectable.

---

## Interaction With Automations

Automations should write only to `.harness/` runtime state unless they are explicitly tasked with changing committed definitions.

This prevents background jobs from mutating architectural source files casually.

---

## Interaction With Meetings

Meetings may use:

- `harness/templates/meeting.yml` as a starting packet template
- `.harness/meetings/active/` for live meeting state
- `.harness/meetings/archive/` for durable local archives

If meeting outcomes should become committed knowledge, that should happen through explicit synthesis into docs or profiles.

---

## First Implementation Recommendation

For the first proving-ground implementation:

1. commit `harness/README.md`
2. commit `harness/profiles/erwin.yml`
3. ignore `.harness/`
4. create runtime state lazily

This is the smallest clean foundation.
