# Signal Schema

> Date: 2026-03-28
>
> Purpose: define the structured signal model for the orchestration harness so teammates, leader, meetings, and automations can exchange meaningful state without devolving into ungoverned chat.

---

## Why This Exists

Signals are how local discoveries become global governance inputs.

Without a signal schema:

- reports become vague prose
- freshness is hard to judge
- attribution is lost
- leader synthesis becomes unreliable

Signals should therefore be treated as first-class structured objects.

---

## Schema Goals

Signals should be:

- attributable
- fresh or explicitly stale
- scoped
- consumable by leader and automations
- compatible with later acceptance/rejection

Signals are not identical to meeting outputs.
Signals are smaller, more operational, and more frequent.

---

## Canonical Structure

Logical shape:

```yaml
version: 1
signals:
  - ...
```

Each signal is one atomic reportable event or state assertion.

---

## Required Fields

### Identity

- `signal_id`
  Unique signal identifier.

- `kind`
  Signal type.

### Source

- `source_type`
  Examples:
  `leader`, `teammate`, `meeting`, `automation`, `script`

- `source_id`
  Concrete source identity.

- `run_id`
  Optional but strongly recommended when signal emerged from a specific run.

- `mission_id`
  Required when the signal is mission-relevant.

### Target / Scope

- `applies_to`
  List of mission ids, run ids, release ids, or system scopes affected.

- `scope_type`
  Examples:
  `local`, `sibling`, `release`, `system`

### Payload Summary

- `title`
  Very short summary.

- `summary`
  One-paragraph explanation.

### Evidence

- `evidence_refs`
  References to files, commands, PRs, issues, contracts, or test outputs.

- `authority_level`
  The highest authority level used to support this signal.

- `confidence`
  Recommended:
  `low`, `medium`, `high`

### Freshness

- `emitted_at`
- `fresh_until`

Freshness is mandatory because signals decay.

### Lifecycle

- `state`
  One of:
  `active`, `accepted`, `rejected`, `superseded`, `expired`

- `resolution_note`
  Optional explanation once not active.

---

## Recommended Signal Kinds

- `progress`
- `completion`
- `blocked`
- `interface_change`
- `verification_result`
- `risk_raised`
- `dependency_discovered`
- `dependency_cleared`
- `requires_rebase`
- `review_feedback`
- `meeting_required`
- `release_ready`

The exact list can grow, but should stay disciplined.

---

## Example Signal

```yaml
- signal_id: "sig-20260328-17"
  kind: "interface_change"
  source_type: "teammate"
  source_id: "worker-auth"
  run_id: "run-20260328-claim-042-a"
  mission_id: "issue-42"
  applies_to:
    - "issue-42"
    - "issue-45"
  scope_type: "sibling"
  title: "Auth function signature changed"
  summary: "authenticate now requires scope; sibling missions touching auth callers must adjust"
  evidence_refs:
    - "src/auth.ts"
    - "tests/auth.spec.ts"
  authority_level: "L1"
  confidence: "high"
  emitted_at: "2026-03-28T09:35:00Z"
  fresh_until: "2026-03-29T09:35:00Z"
  state: "active"
```

---

## Signal Lifecycle

### `active`

Still relevant and not yet accepted/rejected/superseded.

### `accepted`

Leader or meeting has accepted the signal as operationally relevant.

### `rejected`

Signal was reviewed and found misleading, unsupported, or out of authority bounds.

### `superseded`

Replaced by fresher or better-supported signal.

### `expired`

No longer fresh enough to drive decisions.

---

## Acceptance Rules

Signals should not automatically become system truth.

Recommended rule:

- raw signal -> review/synthesis -> accepted signal or rejected signal

This is how the system avoids gossip-based regulation.

For low-risk operational signals, the leader may accept directly.
For high-ambiguity or high-blast-radius signals, acceptance may require a meeting.

---

## Signal Freshness Rules

Signals decay based on kind.

Examples:

- `progress` decays quickly
- `blocked` should be refreshed regularly
- `interface_change` may remain relevant longer
- `release_ready` becomes invalid when new commits land

So `fresh_until` should be determined by signal kind and context, not one global TTL.

---

## Authority Rules

Signals should include enough evidence to be ranked by authority.

Examples:

- signal backed by live test failure and changed file refs = strong
- signal based only on teammate intuition = weak

Signals with weak authority can still be useful, but they should bias toward:

- `low` confidence
- shorter freshness
- explicit need for verification

---

## Anti-Patterns

Do not allow:

- unattributed signals
- signals without scope
- indefinite freshness
- signals that summarize conclusions without evidence
- stale active signals that no longer map to current reality

These make the leader regulate against noise.

---

## Relationship To Meetings

Signals can trigger meetings.
Meetings can also emit higher-quality signals.

Meeting outputs are not the same as ordinary signals, but accepted meeting conclusions can be materialized as signals for downstream automation and coordination.

---

## Relationship To Automations

Automations should be able to:

- identify stale active signals
- aggregate active signals into reports
- flag unreviewed high-risk signals
- detect signal contradictions

This is one of the strongest reasons to keep signals structured.

---

## Suggested Next Step

This schema should be implemented alongside:

- run registry storage
- accepted/rejected signal workflow
- meeting output materialization rules
