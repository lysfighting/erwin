# Meeting Protocol Spec

> Date: 2026-03-28
>
> Purpose: define when the leader should convene a meeting, how roles are chosen, how debate proceeds, and how collective judgment is converted into actionable system truth.

---

## Why Meetings Exist

Meetings exist because some decisions are too ambiguous or too high-impact for unilateral synthesis.

The goal is not democracy.
The goal is better regulated judgment.

Meetings should be used when:

- the leader needs bounded plural perspectives
- evidence is conflicting
- design quality is in question
- release or architectural risk is high

Meetings should not be used for:

- ordinary routing
- trivial implementation choices
- decisions that already have clear higher-authority evidence

---

## Meeting Triggers

The leader should convene a meeting when one or more of these hold:

1. two credible interpretations of high-authority evidence conflict
2. a design decision has high blast radius
3. multiple teammates surface materially different models of the same problem
4. release risk or production risk is involved
5. a strategy change would alter mission decomposition
6. self-sufficiency or operator usability is in question

---

## Meeting Inputs

Every meeting should start with a short packet containing:

- decision question
- target scope
- known context
- locked decisions
- authority hierarchy
- roles
- expected output

If this packet is weak, the meeting will drift.

---

## Data Authority Hierarchy

Every meeting must define a hierarchy of truth sources.

Typical pattern:

- L1: live code / live git/GitHub state
- L2: current config / current environment
- L3: current design doc / proposal / contract
- L4: auxiliary notes / prior reports / WIP artifacts
- L5: legacy or deprecated sources

Golden rule:

lower-authority sources must not override higher-authority sources.

This is one of the strongest takeaways from `arch-review`.

---

## Role Selection

Meetings should use 2-4 bounded roles.

Good role families:

- fact verifier
- architecture challenger
- completeness checker
- ops simulator
- self-sufficiency tester
- devil's advocate

Role constraints:

- every role must have a clear mission
- every role should have allowed sources
- every role should have forbidden sources
- no more than 4 roles by default

Too many roles dilute signal.

---

## Meeting Phases

### Phase 1: Setup

The leader defines:

- the question
- authority hierarchy
- scope
- roles
- locked assumptions

### Phase 2: Independent Review

Each role examines the question independently.
This protects against premature convergence.

### Phase 3: Cross-Examination

Roles see each other's positions and challenge weak claims.

This is mandatory for real meeting value.

### Phase 4: Synthesis

The leader:

- deduplicates
- checks authority use
- resolves or frames conflicts
- separates factual disagreement from design disagreement

### Phase 5: Output

The meeting produces:

- accepted conclusions
- rejected conclusions
- unresolved questions
- concrete action items
- any new process lessons

---

## Factual vs Design Disagreement

These must be treated differently.

### Factual disagreement

Resolve by:

- source authority
- freshness
- specificity of evidence

### Design disagreement

Resolve by:

- thesis
- antithesis
- synthesis

Both sides may contain truth.

This distinction should be explicit in every meeting.

---

## Output Requirements

Every meeting should produce:

- final decision or recommendation
- evidence basis
- rejected alternatives and reasons
- follow-up actions
- escalation note if still unresolved

The point is not merely to "have discussed."
The point is to improve the system's operational truth.

---

## Anti-Failure Rules

1. No meeting without a clear question.
2. No meeting without an authority hierarchy.
3. No role may read everything by default.
4. No synthesis may silently drop contested points.
5. Rejected findings should be recorded with reasons.
6. Meetings should stay small and purpose-built.

---

## Relationship To The Leader

The leader chairs the meeting.

The leader is responsible for:

- scope discipline
- role quality
- evidence discipline
- synthesis quality
- deciding whether the meeting changed system truth

The leader is not responsible for:

- pretending to know the answer in advance
- collapsing the meeting into ceremonial validation

---

## Relationship To The Harness

Meetings are a governance mechanism inside the harness.

They exist to improve:

- global model quality
- decision quality
- justification quality
- safety under ambiguity

This is how the system avoids both:

- dictator leader failure
- flat swarm failure

---

## Success Criteria

A meeting is successful when:

1. it changes understanding or confidence meaningfully
2. it sharpens the authority-aware truth model
3. it yields concrete next actions
4. it does not create uncontrolled discussion sprawl
5. it leaves a durable record of what was accepted and rejected
