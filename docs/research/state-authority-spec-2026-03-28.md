# State Authority Spec

> Date: 2026-03-28
>
> Purpose: define which sources count as truth in the future orchestration harness so the leader, teammates, and meetings do not regulate against the wrong reality.

---

## Why This Exists

One of the most dangerous failure modes in complex agent systems is source confusion.

Not all state is equally trustworthy.
The harness needs explicit authority rules so that:

- stale notes do not override live repo state
- teammate summaries do not override direct evidence
- legacy artifacts do not override active systems

This is one of the strongest lessons imported from `arch-review`.

---

## Authority Principle

Higher-authority, fresher, more direct evidence should dominate lower-authority, older, more interpretive evidence.

This is not just a review rule.
It should be a harness-wide governance rule.

---

## Canonical Authority Levels

### L1: Live system-of-record state

Examples:

- live git branch / worktree state
- current GitHub issue / PR state
- current commit graph
- current verification outputs
- active release state

Use:

- factual execution truth
- ownership truth
- delivery truth

### L2: Active local orchestration state

Examples:

- active contracts
- run registry
- accepted signals
- current scheduler outputs

Use:

- current harness operating model
- accepted near-live orchestration truth

### L3: Current design / planning artifacts

Examples:

- active design specs
- mission plans
- meeting packets
- current roadmap artifacts

Use:

- intended state
- design intent
- scoped planning

### L4: Auxiliary or derived artifacts

Examples:

- prior syntheses
- cached analysis outputs
- draft reports
- unconfirmed teammate summaries

Use:

- hints
- context
- candidate leads for verification

### L5: Legacy / deprecated artifacts

Examples:

- superseded configs
- old repos
- abandoned design docs
- stale contracts

Use:

- historical context only

Never use as ground truth unless explicitly reactivated.

---

## Freshness Rule

Authority is not only about type.
It is also about freshness.

Examples:

- a stale L2 run registry entry may be less useful than fresh L1 git state
- an old teammate summary may be invalidated by a new commit

So every important state source should be evaluated on:

- authority class
- freshness
- directness
- attribution

---

## Direct vs Derived State

### Direct state

Observed from the current system.

Examples:

- open PR exists
- tests failed
- branch points to commit X

### Derived state

Inferred or summarized from direct state.

Examples:

- "mission is nearly done"
- "teammate probably unblocked"
- "design seems stable"

Rule:

derived state must not silently override direct state.

---

## Teammate Report Authority

Teammate reports are important, but not absolute.

Use them as:

- local evidence
- signal source
- ambiguity surfacer

Do not use them as:

- unquestioned system truth
- replacement for direct repo state

The right pattern is:

- teammate reports inform
- leader verifies or synthesizes
- accepted conclusions enter higher-trust orchestration state

---

## Meeting Output Authority

Meeting outputs are stronger than raw teammate opinions, but still not equal to live repo state.

They should be treated as:

- accepted interpretation
- design truth pending supersession
- action-driving governance output

But if later live evidence contradicts them, they must be re-opened.

---

## Authority Conflicts

When sources conflict, resolve in this order:

1. higher authority over lower authority
2. fresher over staler
3. direct over derived
4. specific over vague
5. accepted synthesis over raw individual opinion

When conflict remains unresolved:

- mark it explicitly
- escalate if it changes important decisions

---

## Forbidden Patterns

Do not allow:

- L4/L5 artifacts to override L1/L2 state
- prior summaries to replace fresh verification
- unscoped memory to become operational truth
- teammate speculation to become accepted state without synthesis

These patterns create regulator hallucination.

---

## Authority Application By Layer

### Leader

Must always reason with authority ranking.

### Teammates

Must know which sources they may trust for a given assignment.

### Meetings

Must declare the hierarchy explicitly before debate.

### Automations

Should prefer direct and current signals, not ambiguous summaries.

---

## Success Criteria

The state authority model is working when:

1. stale artifacts do not hijack decisions
2. meetings produce clearer truth, not more confusion
3. teammates know what they can and cannot rely on
4. leader synthesis remains grounded in live evidence
5. legacy context remains useful without becoming operational truth
