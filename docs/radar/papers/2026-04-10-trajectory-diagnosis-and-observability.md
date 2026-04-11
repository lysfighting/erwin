# Paper Watch — Trajectory Diagnosis And Observability

> Date: 2026-04-10
>
> Purpose: track the recent paper signal moving coding-agent evaluation beyond resolution rates toward trajectory diagnosis, intermediate observability, and mechanism-level failure analysis.

---

## Why This Note Exists

The paper-watch archive already covers self-iteration, long-horizon degradation, and context retrieval.

What it still needed was a grouped note for a related shift:

- coding-agent evaluation is becoming more trajectory-aware
- output-only success metrics are looking less sufficient
- intermediate observability is starting to matter as much as final resolution rates

This note extends:

- [2026-04-02-coding-agent-self-iteration-watch.md](2026-04-02-coding-agent-self-iteration-watch.md)
- [2026-04-03-long-horizon-coding-agent-benchmarks.md](2026-04-03-long-horizon-coding-agent-benchmarks.md)
- [../../knowledge/active-threads.md](../../knowledge/active-threads.md)

---

## Watch Set

### `Beyond Resolution Rates`

- `Beyond Resolution Rates: Behavioral Drivers of Coding Agent Success and Failure`
  - https://arxiv.org/abs/2604.02547
  - Submitted: 2026-04-02
  - Why it matters: argues that trajectory structure is more diagnostic than raw trajectory length, and that success is associated with context gathering before editing plus stronger validation behavior.
  - Current signal: simple patches can still remain unsolved because of architectural reasoning and domain-knowledge gaps, and shared LLM choice explains more than shared framework choice.

### `TRAJEVAL`

- `TRAJEVAL: Decomposing Code Agent Trajectories for Fine-Grained Diagnosis`
  - https://arxiv.org/abs/2603.24631
  - Submitted: 2026-03-25
  - Why it matters: decomposes coding-agent trajectories into `search`, `read`, and `edit`, making failure modes inspectable instead of collapsing everything into one pass/fail number.
  - Current signal: trajectory-level diagnostics can be predictive and actionable, improving outcomes while also reducing costs when fed back into the agent loop.

### `The Observability Gap`

- `The Observability Gap: Why Output-Level Human Feedback Fails for LLM Coding Agents`
  - https://arxiv.org/abs/2603.26942
  - Submitted: 2026-03-27
  - Why it matters: sharpens the claim that output-only review is often too lossy for coding-agent correction because internal code state and visible outcomes are separated by deep causal chains.
  - Current signal: minimal code-level diagnostic visibility restored convergence in the studied setup, suggesting that feedback observability can be a stronger bottleneck than raw programming competence.

---

## Current Interpretation

- coding-agent evaluation is shifting from `did it resolve the issue` toward `what did the agent do, where did it fail, and what was visible during failure`
- trajectory structure, stage-level diagnosis, and intermediate observability now look like reusable eval principles, not only academic analysis categories
- this cluster strengthens the repo's existing preference for legibility, bounded runtime objects, and eval-first pressure before wider autonomy claims

---

## Why This Matters For erwin

This note strengthens three active threads:

- `coding-agent self-iteration and paper intake`
- `harness engineering as legibility + eval pressure`
- `benchmark integrity and eval awareness`

The immediate implication is not to import paper metrics directly.

The immediate implication is:

- keep paper watch moving from benchmark outcome to trajectory diagnosis
- treat observability as part of eval design, not only UI or postmortem quality
- prefer local review criteria and eval checks that make failure stages visible before copying public benchmark structure

---

## Watch Next

- whether trajectory-stage diagnostics start appearing in first-party engineering writeups rather than only in papers
- whether public coding benchmarks begin publishing stronger intermediate traces and resource-methodology details
- whether observability and monitorability become explicit design goals in more agent-harness systems, not only implied background assumptions
