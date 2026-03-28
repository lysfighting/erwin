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

- protocol and adapter work are much more mature than control-plane work
- runtime packaging is much more mature than desired-state modeling
- meeting governance is now strong, but automated council orchestration is still not there

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

### 1. Control plane exists, but desired-state control is still under-realized

The research repeatedly framed the harness as:

- cybernetic
- Kubernetes-inspired
- desired-state vs actual-state
- control-plane-led

Current truth:

- `doctor`, `home`, `maintenance`, and runtime state are real
- but a first-class **mission graph / desired-state object** still does not exist

This means the leader behaves like a regulator in spirit, but still lacks an explicit declared system model to regulate against.

This is the single biggest remaining gap.

### 2. Ownership discipline is still mostly textual

The research stressed:

- explicit ownership boundaries
- overlap control
- controlled parallelism
- avoiding silent collision between teammates

Current truth:

- teammate scope is bounded
- teammate communication is governed
- but scope is still mostly free-text

So Kubernetes-style ownership discipline has only been partially absorbed.
The system still lacks explicit overlap selectors or a stronger ownership surface model.

### 3. State authority is defined more strongly than it is enforced

The research on authority hierarchy was very explicit:

- higher-authority evidence should dominate
- stale summaries must not override live state
- authority should be a harness-wide rule, not just a meeting convention

Current truth:

- meetings now model authority hierarchy well
- signals carry authority metadata
- but there is still no broad harness-wide enforcement layer that systematically downgrades or rejects low-authority state outside specific paths

So this area is **partially preserved, not fully realized**.

### 4. Council automation is still behind council semantics

The research wanted:

- independent review
- cross-examination
- synthesis
- council-corrected judgment

Current truth:

- the runtime now stores the right meeting objects
- but it still does not automate the debate workflow itself

So the system has become **council-ready**, but not yet **council-orchestrating**.

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
- explicit ownership model: still partial
- authority enforcement across the whole harness: still partial
- automated council orchestration: still partial

That means future work should **not** primarily add new surfaces.
It should close the remaining control-plane gaps.

---

## Recommended Next Priority Order

1. Explicit ownership / overlap model for teammate assignments
2. First-class mission graph / desired-state layer
3. Broader authority enforcement outside meetings
4. Automated council orchestration only after the above are stable

If future sessions follow this order, the harness will continue toward the original design rather than only becoming a nicer adapter.
