# Harness Decision Methodology Synthesis

> Status: Active Research
>
> Date: 2026-03-30
>
> Purpose: prevent `Harness` from overfitting to a single `arch-review`-style council model by synthesizing a broader decision-methodology stack for leader-led orchestration.

---

## Core Conclusion

`Harness` should not treat `arch-review` as the only governing methodology.

The stronger design is:

- **cybernetics + Kubernetes** for control-plane structure
- **evidence-review councils** for contested technical judgment
- **incident command** for urgent operational and release decisions
- **decision records** for slower, lower-conflict architectural choices

In other words:

`Harness` needs a **protocol family**, not a single debate ritual.

---

## Why This Matters

If every important decision is forced through one debate choreography, the system drifts toward ritual instead of control.

That is risky because the decisions inside `Harness` are not all the same:

- some are contested evidence questions
- some are urgent release/rollback decisions
- some are slower design or tradeoff decisions that need durable recording more than live challenge rounds

The control plane should therefore select the right decision protocol for the decision surface instead of silently assuming one default style.

---

## Methodology Inputs

### 1. Cybernetics

Cybernetics contributes the regulator model:

- distinguish desired state from actual state
- detect disturbances and drift
- correct based on higher-quality models and fresher signals
- treat stale state and authority conflict as control-plane concerns

This remains the backbone of `Harness`.

Key reference:

- [Conant & Ashby, “Every good regulator of a system must be a model of that system”](https://www.tandfonline.com/doi/abs/10.1080/00207727008920220)

### 2. Kubernetes Control-Plane Thinking

Kubernetes contributes the control-loop structure:

- controllers reconcile toward desired state
- ownership and scope must be explicit
- admission and validation gates matter
- maintenance and garbage collection are part of the system, not cleanup trivia

This strongly supports `Harness` mission/run/maintenance/scheduler design, but Kubernetes does not tell us how every human decision should be governed.

Key reference:

- [Kubernetes Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)

### 3. `arch-review`

`arch-review` contributes bounded plural judgment:

- the leader acts as chair/synthesizer, not unchecked oracle
- participants have bounded roles
- allowed and forbidden sources matter
- rejected findings require reasons
- credible contested lanes should be challenged before synthesis

This remains highly valuable for **technical disagreement with meaningful evidence conflict**.

But it should be treated as one important precedent, not the entire methodology of `Harness`.

### 4. Incident Command / Incident Management

Operational and release-adjacent decisions often need a different discipline:

- stabilize first
- assign command and action ownership clearly
- avoid letting long design debate delay bounded response
- keep explicit checkpoints and stop/go constraints

This is a better fit for:

- hotfix or rollback decisions
- release/RC uncertainty
- urgent production-affecting decisions

Key references:

- [Google SRE: Managing Incidents](https://sre.google/sre-book/managing-incidents/)
- [Google SRE Workbook: Incident Response](https://sre.google/workbook/incident-response/)
- [NIST SP 800-61 Rev. 3, Incident Response Recommendations and Considerations for Cybersecurity Risk Management](https://csrc.nist.gov/pubs/sp/800/61/r3/final)

### 5. Decision Records / ADRs

Some decisions do not need live debate control at all.

What they need is:

- context
- alternatives considered
- chosen decision
- consequences
- follow-up obligations

This is the right baseline for:

- bounded design choices
- architecture tradeoffs with low immediate blast radius
- decisions where durable recording matters more than live rebuttal rounds

Key reference:

- [Architecture Decision Records](https://adr.github.io/)

---

## Protocol Taxonomy For Harness

`Harness` should treat the following as first-class decision protocols.

### `evidence-review`

Use when:

- technical or factual disagreement is real
- evidence quality and authority matter
- one side may need to challenge hidden assumptions or blast radius

Expected shape:

- independent review
- challenge only when a credible contested lane remains
- synthesis into accepted/rejected/unresolved findings

This is the closest descendant of `arch-review`, but still broader than that one skill.

### `incident-command`

Use when:

- the decision is release-adjacent or operationally urgent
- action ownership matters more than exhaustive debate
- the system must stabilize or choose a bounded immediate response

Expected shape:

- current situation
- immediate action choice
- explicit owner
- stop/go constraints
- next checkpoint

This does **not** reject evidence review.
It changes the priority order: control and bounded response come first.

### `decision-record`

Use when:

- the decision is important but not urgent
- the main risk is ambiguity or future forgetfulness
- the goal is to make tradeoffs durable and auditable

Expected shape:

- context
- alternatives
- chosen decision
- consequences
- follow-up work

This is the slowest and least conflict-driven protocol.

---

## Selection Logic

The leader or scheduler should choose protocol based on decision surface, not ideology.

Baseline mapping:

- contested technical judgment: `evidence-review`
- release/rollback/hotfix/urgent operational choice: `incident-command`
- bounded design choice without urgent conflict: `decision-record`

This does not mean the system must fully automate all three protocols equally today.

It means the control plane should already be able to:

- name the protocol explicitly
- explain why it was chosen
- preserve different expected outputs and role defaults

---

## What This Means For Runtime

The current `Harness` runtime still uses one shared meeting execution engine:

- `setup`
- `independent_review`
- `cross_examination`
- `synthesis`
- `output`

That is acceptable for the current prototype line **as long as**:

- the protocol is explicit
- the system stops pretending every meeting is implicitly `arch-review`
- protocol choice affects defaults and leader guidance
- future work can specialize protocol behavior without breaking object shape

So the immediate runtime consequence is:

- meeting objects should carry `decision_protocol`
- scheduler human-decision packets should carry `recommended_decision_protocol`
- documentation should describe `arch-review` as one absorbed influence, not the whole council model

---

## Design Guardrail

The right sentence for future work is:

**`Harness` is methodology-plural but control-plane-consistent.**

That means:

- one regulator model
- one authority model
- one mission/run backbone
- multiple decision protocols chosen according to decision type

This is more faithful to the actual system we are building than “make everything more like `arch-review`.”
