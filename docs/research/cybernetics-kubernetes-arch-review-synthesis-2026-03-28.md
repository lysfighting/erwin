# Cybernetics, Kubernetes, Arch-Review, and Harness Synthesis

> Date: 2026-03-28
>
> Purpose: synthesize engineering cybernetics, Kubernetes control-plane ideas, git-ordered execution, and `arch-review`'s debate model into a sharper design for the new leader-led orchestration harness.

---

## Executive Conclusion

Yes, these systems are strongly related.

Your emerging design is not just "multi-agent collaboration."
It is increasingly recognizable as a **cybernetic harness**:

- the leader is a regulator
- teammates are controlled execution units
- git/GitHub objects provide a partially ordered, inspectable state model
- meetings/debate provide structured model correction
- automations and controllers keep actual state aligned with desired state

But there is a crucial correction:

**the leader should be a chairing regulator, not a solitary decider**

That is where `arch-review` becomes especially important.

The best future design is not:

- dictator-leader
- flat swarm
- pure scheduler

It is:

**leader-led, council-corrected, git-grounded orchestration**

---

## Why Cybernetics Is Relevant

Engineering cybernetics is fundamentally about regulation under feedback.

The most relevant idea here is not "AI" or "automation."
It is:

- desired state
- actual state
- disturbance
- feedback
- correction
- model quality

The Conant-Ashby result is especially relevant: a good regulator must in some meaningful sense embody a model of the system it regulates.

Applied here, that means:

the leader cannot just assign tasks.
The leader must maintain a useful model of:

- the mission graph
- repo structure
- ownership boundaries
- dependency structure
- teammate state
- release/risk state

If the leader lacks that model, it becomes a dispatcher, not a regulator.

Source:

- [Conant & Ashby, 1970](https://www.tandfonline.com/doi/abs/10.1080/00207727008920220)

---

## Why Kubernetes Is Relevant

Kubernetes gives a very practical engineering example of cybernetic regulation at scale.

The important borrowed ideas are structural, not literal.

### 1. Desired state vs actual state

Kubernetes controllers continuously move actual cluster state toward declared desired state.

For your harness, that maps cleanly to:

- desired state = mission graph + contracts + release goals + policies
- actual state = current branches, worktrees, PRs, tests, teammate progress, blockers

Source:

- [Kubernetes Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
- [Kubernetes Deployments](https://v1-32.docs.kubernetes.io/docs/concepts/workloads/controllers/deployment/)

### 2. Control plane vs workers

Kubernetes distinguishes between control plane components that make global decisions and worker nodes that execute workloads.

That is directly analogous to:

- leader / coordination layer
- teammate execution units

Source:

- [Kubernetes Cluster Architecture](https://kubernetes.io/docs/concepts/architecture/)

### 3. Reconciliation loops

Controllers do not assume the world stays correct after one command.
They watch, compare, and reconcile repeatedly.

That is exactly how your leader should behave:

- inspect current state
- compare against intended state
- correct deviations
- re-evaluate after teammate reports and environmental changes

Source:

- [Writing a Controller for Pod Labels](https://kubernetes.io/blog/2021/06/21/writing-a-controller-for-pod-labels/)

### 4. Admission and validation gates

Kubernetes admission controllers mutate or validate changes before they are accepted.

That maps well to:

- pre-claim validation
- pre-ship checks
- RC admission gates
- policy enforcement before risky transitions

Source:

- [Admission Controllers](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/)

### 5. Explicit ownership boundaries

Kubernetes warns that overlapping selectors or conflicting ownership cause controllers to interfere unexpectedly.
Server-side Apply also formalizes field ownership.

That matters a lot for your harness:

- two teammates should not silently own the same execution surface
- overlap must be explicit
- ownership transfer must be formalized

Source:

- [Kubernetes Deployments](https://v1-32.docs.kubernetes.io/docs/concepts/workloads/controllers/deployment/)
- [Server-Side Apply GA](https://kubernetes.io/blog/2021/08/06/server-side-apply-ga/)

### 6. Garbage collection is first-class

Kubernetes has cleanup controllers because stale state accumulates naturally.

Your harness needs the same mentality:

- stale contracts
- stale runs
- stale worktrees
- stale signals

These are not edge cases.
They are normal control-plane concerns.

Source:

- [kube-controller-manager](https://v1-32.docs.kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/)

---

## What Not To Copy From Kubernetes

Kubernetes is a very strong analogy, but it should not be copied blindly.

Do not copy:

- excessive abstraction layers before need is proven
- mechanical reconciliation when judgment is actually required
- API object proliferation without strong semantics
- treating all conflict as something a controller can auto-resolve

Software-development orchestration contains more ambiguity and design judgment than pod scheduling.

So Kubernetes should inspire:

- control-plane discipline
- ownership clarity
- reconciliation mindset

not:

- literal architecture mimicry

---

## Why Git / GitHub Orderliness Is Central

You are also right that one of `teamwork`'s deepest insights was its dependence on git/GitHub orderliness.

This is not incidental.
It is the substrate that makes regulation possible.

### Git/GitHub gives you inspectable objects

- issue
- branch
- worktree
- contract
- commit
- PR
- label
- milestone
- release

These are the equivalent of control-plane objects.

### Git/GitHub gives partial ordering

- issue before claim
- claim before active branch
- branch before PR
- PR before merge
- merge before release progression

That ordering makes teammate execution governable.

### Git/GitHub gives auditable state transitions

Unlike pure chat coordination, git/GitHub gives:

- durable records
- ownership markers
- review checkpoints
- rollback points
- release grouping

This is why `teamwork` could be made controllable in the first place.

So the future harness should not merely "integrate with git."
It should treat git/GitHub as the primary state backbone.

---

## Why Arch-Review Matters

Your additional correction is extremely important:

**leader should not simply decide alone**

This is where `arch-review` becomes more than a side skill.
It provides a pattern for regulated collective cognition.

The valuable ideas from `arch-review` are:

### 1. Lead as synthesizer, not oracle

The Lead organizes, scopes, challenges, and synthesizes.
It does not blindly trust one role and it does not abdicate judgment.

### 2. Explicit data authority hierarchy

Not all sources are equally authoritative.
This is essential for orchestration too:

- live git/GitHub state
- active contracts
- teammate reports
- prior summaries
- legacy notes

These must not be treated as equivalent truth.

### 3. Role-bounded viewpoints

Each reviewer has allowed and forbidden sources.
That prevents premature convergence and groupthink.

For your harness, this suggests:

- not every teammate should see or decide everything
- some discussions should happen through intentionally bounded roles

### 4. Cross-examination before synthesis

Findings are not accepted just because one role said them.
They are challenged, refined, or withdrawn.

This is exactly the pattern you want for leader-led "meetings."

### 5. Small councils beat noisy swarms

`arch-review` explicitly learned that too many roles reduce signal quality.

That is a direct warning for orchestration:

- meeting size must be small
- roles must be focused
- no giant all-hands by default

Relevant local references:

- [`arch-review/SKILL.md`](/Users/liyasong/Project/collov/erwin/skills/arch-review/arch-review/SKILL.md)
- [`debate.md`](/Users/liyasong/Project/collov/erwin/skills/arch-review/arch-review/operations/debate.md)
- [`setup.md`](/Users/liyasong/Project/collov/erwin/skills/arch-review/arch-review/operations/setup.md)
- [`lessons.md`](/Users/liyasong/Project/collov/erwin/skills/arch-review/reference/lessons.md)

---

## The Key Architectural Shift

Putting these together changes the interpretation of your future system.

It is not just:

- protocol + adapter + scheduler

It becomes:

### 1. Git-Grounded State Plane

Issue, contract, branch, PR, release, labels.

### 2. Leader Control Plane

Maintains the model of the system and decides how work should progress.

### 3. Teammate Execution Plane

Executes bounded tasks and reports state back upward.

### 4. Council / Meeting Plane

Activated when the leader should not decide alone:

- architecture disputes
- ambiguous root cause
- strategy change
- release risk
- competing interpretations of evidence

### 5. Maintenance Plane

Recurring automations keep stale state and drift under control.

This is much richer than the earlier "scheduler" abstraction.

---

## Leader As Chair, Not Emperor

This is probably the most important refinement from your message.

The leader should not be modeled as:

- supreme autonomous planner
- solitary global reasoner
- final source of truth

Instead, the leader should be modeled as:

- state owner
- agenda setter
- policy enforcer
- council convener
- synthesis authority
- execution governor

That means:

### The leader decides alone for

- straightforward routing
- task assignment
- run-state transitions
- mechanical rescheduling
- hygiene actions

### The leader convenes a meeting for

- nontrivial design disagreement
- conflicting evidence
- high-blast-radius architectural changes
- uncertainty that materially affects plan quality
- strategic reprioritization

This is a far healthier model.

---

## Harness Implications

This has direct implications for harness design.

### 1. The harness needs a model store

Because a good regulator needs a model, the leader needs structured state:

- mission graph
- run registry
- teammate status
- overlap/dependency map
- signals
- unresolved questions

### 2. The harness needs policy gates

Kubernetes-style admission logic maps naturally to:

- claim admission
- ship admission
- rc admission
- meeting-required gates

### 3. The harness needs reconciliation loops

The leader should continuously reconcile:

- desired mission state
- actual repo state
- reported teammate state
- release state

### 4. The harness needs meeting protocols

This is the new addition beyond classic scheduler design.

Meetings need:

- trigger conditions
- role templates
- data authority rules
- synthesis procedure
- rejected-finding tracking

### 5. The harness needs garbage collection

Because stale orchestration state is unavoidable.

### 6. The harness needs explicit ownership

Without ownership, overlap becomes chaos.

---

## Risks of This Design

This design is strong, but not free.

### 1. Leader bottleneck

If the leader becomes the only high-bandwidth thinker, throughput collapses.

### 2. Meeting inflation

If every ambiguity becomes a meeting, the system will stall.

### 3. Stale world model

If teammate reporting is weak or delayed, the leader will regulate against fiction.

### 4. Council capture

If one role/viewpoint repeatedly dominates, the "meeting" becomes theater.

### 5. Over-control

If teammates are made too passive, they stop being useful local reasoners and become brittle executors.

This is why the system needs both:

- disciplined control
- bounded local autonomy

---

## What To Borrow, Concretely

Borrow from cybernetics:

- regulator-model principle
- feedback loops
- disturbance awareness
- model correction

Borrow from Kubernetes:

- control plane / worker split
- reconciliation loops
- declarative desired state
- admission/validation gates
- explicit ownership boundaries
- garbage collection

Borrow from `teamwork`:

- git/GitHub orderliness
- protocol lifecycle
- contracts as execution views
- worktree/branch discipline

Borrow from `arch-review`:

- leader as synthesizer
- explicit data authority
- bounded roles
- cross-examination
- small councils
- critique with evidence

Do not borrow blindly:

- giant hierarchies
- overcentralized decision making
- uncontrolled horizontal gossip
- too many roles
- auto-reconciliation for decisions that require real judgment

---

## Updated Definition

The future system is best defined as:

**a git-grounded, leader-led, council-corrected orchestration harness for autonomous software work**

That is a more accurate definition than:

- multi-agent scheduler
- protocol adapter
- skill pack

because it captures:

- cybernetic regulation
- control-plane discipline
- bounded worker execution
- structured collective judgment

---

## The Most Important Next Specs

This synthesis suggests the next documents should be:

1. `leader-spec.md`
   What the leader knows, owns, decides, and escalates.

2. `teammate-spec.md`
   What a worker receives, what autonomy it has, and how it reports.

3. `meeting-protocol-spec.md`
   When the leader convenes a council, which roles are used, and how synthesis works.

4. `state-authority-spec.md`
   The hierarchy of truth sources for orchestration.

Those four will operationalize this synthesis.
