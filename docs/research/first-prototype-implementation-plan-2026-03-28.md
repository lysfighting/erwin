# First Prototype Implementation Plan

> Date: 2026-03-28
>
> Purpose: translate the current harness architecture into a practical first prototype plan that can be implemented incrementally in Codex without collapsing back into prompt-heavy design.

---

## Prototype Goal

Build the smallest useful prototype of the new harness that proves all of these:

1. project profile discovery works
2. runtime state can be persisted outside chat
3. leader-oriented orchestration can be grounded in a run registry and signal log
4. Codex surfaces can stay thin while the substrate holds real state

This prototype does not need full V4 orchestration yet.

---

## What Prototype v0 Must Include

### 1. Profile discovery

Read:

- explicit override if provided
- otherwise `harness/profiles/erwin.yml`
- otherwise fallback inference

### 2. Runtime state initialization

Create lazily when needed:

- `.harness/runs/registry.yml`
- `.harness/signals/log.yml`
- `.harness/maintenance/backlog.yml`

### 3. Contract creation path

Use:

- `harness/templates/contract.yml`

Goal:

- prove `claim` can generate a contract from committed definitions plus live issue context

### 4. Run registration path

Use:

- `harness/templates/run.yml`

Goal:

- prove new work enters durable orchestration state cleanly

### 5. Signal emission path

Use:

- `harness/templates/signal.yml`

Goal:

- prove teammate or leader discoveries can become structured signals

### 6. First two surfaces

- `home`
- `claim`

These are enough to validate the architecture.

---

## What Prototype v0 Should Not Include

Do not include yet:

- full multi-run scheduler
- rich meeting orchestration
- automatic cross-teammate messaging
- heavy automation suite
- full RC migration

Those belong after the substrate is proven.

---

## Suggested Implementation Order

### Step 1: Build state path helpers

Needed capabilities:

- discover repo root
- resolve `harness/` and `.harness/`
- create runtime directories lazily
- resolve profile path

### Step 2: Build profile loader

Needed capabilities:

- read committed profile
- read optional local override
- merge safely
- expose unresolved values explicitly

### Step 3: Build run registry helper

Needed capabilities:

- initialize empty registry
- read runs
- add run
- update run
- mark stale

### Step 4: Build signal log helper

Needed capabilities:

- initialize empty log
- append signal
- update signal state
- list active signals by mission/run

### Step 5: Build contract materializer

Needed capabilities:

- load template
- fill mission identity and execution fields
- write contract to runtime state

### Step 6: Implement `claim`

Needed capabilities:

- validate issue
- create branch/worktree
- materialize contract
- register run
- emit initial summary

### Step 7: Implement `home`

Needed capabilities:

- read registry and signals
- compare to live git/GitHub state
- report active runs, stale state, blockers, next action

---

## Candidate Implementation Locations

This depends on naming later, but for now the structure could look like:

```text
scripts/
  harness-paths.sh
  harness-profile.sh
  harness-runs.sh
  harness-signals.sh
  harness-contracts.sh

skills/
  [future-neutral-name]/
    home/
    claim/
```

The important part is separation:

- scripts hold deterministic state operations
- skills orchestrate user-facing flow

---

## Prototype Success Criteria

Prototype v0 succeeds when:

1. a mission can be claimed using the committed `erwin` profile
2. a contract is materialized from template plus live context
3. a run entry is written to the registry
4. `home` can read and summarize that state
5. stale or inconsistent state is detectable

If those work, the architecture is no longer speculative.

---

## Immediate Next Build Recommendation

The most efficient next coding step is:

1. implement path/profile/run/signal deterministic helpers
2. then implement thin `claim`
3. then implement thin `home`

That order keeps the substrate correct before any UX layer gets elaborate.

Update:

- first deterministic helpers now exist under `scripts/`
- the next practical move is a minimal `claim` / `home` prototype on top of that substrate
- initial `home` / `claim` prototypes now exist, along with first `drive` state-transition helpers
