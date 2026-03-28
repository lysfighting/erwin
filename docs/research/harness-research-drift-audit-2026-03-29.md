# Harness Research Drift Audit — 2026-03-29

> Purpose: re-audit the original harness research after several implementation-heavy sessions and verify whether important design commitments were dropped, weakened, or silently reframed.

---

## Scope Reviewed

This drift audit re-checked the main research documents that defined the target system:

- `codex-native-orchestration-harness-analysis-2026-03-28.md`
- `protocol-core-spec-2026-03-28.md`
- `control-plane-spec-2026-03-28.md`
- `leader-spec-2026-03-28.md`
- `teammate-spec-2026-03-28.md`
- `meeting-protocol-spec-2026-03-28.md`
- `state-authority-spec-2026-03-28.md`
- `cybernetics-kubernetes-arch-review-synthesis-2026-03-28.md`
- `codex-runtime-adapter-spec-2026-03-28.md`

It then compared those commitments to the current harness implementation.

---

## Executive Judgment

**No major concept has been fully lost.**

But there is clear **implementation-weighted drift**:

- protocol and adapter work are still a little more mature than the hardest control-plane work
- runtime packaging was ahead of desired-state modeling, but first-class mission objects have now materially closed that gap
- meeting governance is now strong, and automated council orchestration has begun to materialize through claim-aware auto-synthesis, but the full debate controller is still not there

So the danger is not “we forgot the vision entirely.”
The danger is:

**the repo now looks more complete than the original design actually is**

because the surfaces that are easiest to demo and smoke-test are ahead of the harder control-plane commitments.

---

## What Has Been Preserved Well

### 1. The core reframe survived

The research insisted that the new thing should be:

- not `teamwork` migrated again
- not `wing` renamed
- not a thin compatibility layer

That reframe still survives.
The current harness is clearly a Codex-native orchestration harness adapter rather than a command-port.

### 2. Protocol seriousness survived

The protocol-core research said the non-negotiable loop is:

`issue -> claim -> drive -> ship -> rc`

That survived very well.
This is one of the most complete parts of the implementation.

### 3. Leader / teammate split survived

The research emphasized:

- leader as regulator and chair
- teammates as bounded execution units
- governed sibling communication

That also survived clearly in the runtime.

### 4. Council correction survived

Originally this was only partially true.
After the latest governance-hardening pass, meetings now carry:

- role specs
- source-discipline boundaries
- known context
- locked decisions
- rejected reasons

So this part is now much closer to the original design than before.

### 5. Harness engineering survived

The adapter, bootstrap, sync, manifest, smoke suite, Go runtime, and packaged-binary preference all strongly reflect the original harness-engineering direction.

---

## Where Drift Is Real

These are the places where the current implementation is still weaker than the research vision.

### 1. Control plane is now materially stronger, but still not fully mission-bundle-aware

The research repeatedly framed the harness as:

- cybernetic
- Kubernetes-inspired
- desired-state vs actual-state
- control-plane-led

Current truth:

- `doctor`, `home`, `maintenance`, dependency waiting, and run-level desired state are real
- a first-class **mission registry / mission graph surface** now exists
- legacy runtime state can be upgraded into that mission-aware model through `maintenance normalize`

This means the leader no longer regulates only against per-run state.
But the harness still lacks a richer multi-mission bundle model above the current mission graph.

### 2. Ownership discipline is stronger, but still not mature enough

The research stressed:

- explicit ownership boundaries
- overlap control
- controlled parallelism
- avoiding silent collision between teammates

Current truth:

- teammate scope is bounded
- teammate communication is governed
- runs now carry explicit ownership claims
- active exclusive overlap is now rejected

So Kubernetes-style ownership discipline is no longer just textual.
The latest passes materially improved this:
- selectors are now normalized
- bounded ownership can be formally transferred between runs
- selector families now capture richer path overlap semantics

What is still missing is the next layer above that stronger baseline.

### 3. State authority is stronger, but still not fully harness-wide

The research on authority hierarchy was very explicit:

- higher-authority evidence should dominate
- stale summaries must not override live state
- authority should be a harness-wide rule, not just a meeting convention

Current truth:

- meetings now model authority hierarchy well
- signals carry authority metadata
- invalid authority levels are now rejected or surfaced
- superseded active signals can now be adjudicated against fresher run or meeting state
- same-claim meeting findings can now be audited for higher-authority winners and unresolved same-authority conflicts
- run-level authority audit now explains higher-authority active-signal conflicts plus meeting adoption gaps directly on runs
- contract snapshots are now kept closer to active run state, normalize can repair older drift, run-vs-worktree branch mismatch is now explicitly surfaced, and dirty live worktrees can now be flagged when they contradict governed verification or handoff intent
- but there is still no broad harness-wide conflict-resolution layer that systematically adjudicates the harder live-git-vs-contract disagreements beyond this cleaner baseline when no governing signal or council outcome exists

So this area is **materially stronger, but still only partially realized**.

### 4. Council automation is improving, but still behind council semantics

The research wanted:

- independent review
- cross-examination
- synthesis
- council-corrected judgment

Current truth:

- the runtime now stores the right meeting objects
- the runtime can now advance meeting phases with `meetings orchestrate`
- findings can now reference the positions that support or challenge them, so synthesis provenance is no longer implicit
- positions can now be grouped by `claim` / `dimension`, `meetings orchestrate` can auto-generate provisional findings plus unresolved follow-up actions from those grouped positions, and meetings now track routed vs unrouted conflicting claim lanes through a first-class debate plan
- but it still does not automate the full debate workflow itself

So the system has become **meaningfully council-orchestrating**, but not yet a full automated debate controller.

### 5. Product surface questions are still unresolved

The research repeatedly warned against prematurely collapsing:

- protocol
- runtime adapter
- product naming

That warning still matters.
The repo now has a strong adapter, but the public package contract and final product naming are still unresolved.

---

## What Was At Risk Of Being Lost, But Was Recovered

The most likely thing to have been silently lost was:

**`arch-review`-style council discipline**

That was genuinely under-implemented before the latest pass.
Without the new role specs and rejection reasons, the meeting system had drifted toward “durable notes” rather than “bounded plural judgment.”

That specific drift has now been corrected.

---

## Bottom Line

The current harness has **not** drifted into the wrong product.

But it **has** drifted into an uneven maturity shape:

- adapter/runtime/productization: strong
- lifecycle protocol: strong
- leader/teammate runtime: strong
- council governance packet: now strong
- desired-state control plane: still partial
- explicit ownership model: now materially stronger, still partial
- authority enforcement across the whole harness: still partial
- automated council orchestration: materially stronger, still partial

That means future work should **not** primarily add new surfaces.
It should close the remaining control-plane gaps.

---

## Recommended Next Priority Order

1. Richer selector families and transfer policies beyond the current normalized ownership-transfer baseline
2. Broader authority conflict resolution outside meetings
3. Deeper automated council orchestration
4. Product naming / public contract once the governance model feels stable

If future sessions follow this order, the harness will continue toward the original design rather than only becoming a nicer adapter.
