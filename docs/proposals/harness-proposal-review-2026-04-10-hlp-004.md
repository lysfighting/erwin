# Harness Proposal Review — HLP-004 — 2026-04-10

> Status: Active
>
> Date: 2026-04-10
>
> Proposal: `HLP-004`
>
> Purpose: execute a concrete proposal review for `HLP-004` using the repo's paper-watch, thread-anchor, and proposal-gating surfaces.

---

## Review Scope

This review applies the current proposal surfaces to:

- [harness-proposal-ledger.md](harness-proposal-ledger.md)
- [../knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md](../knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md)
- [../radar/papers/2026-04-02-coding-agent-self-iteration-watch.md](../radar/papers/2026-04-02-coding-agent-self-iteration-watch.md)
- [../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md](../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md)
- [../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md](../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md)
- [../radar/weekly/2026-04-10.md](../radar/weekly/2026-04-10.md)

The goal is to answer one narrower question:

- should `HLP-004` still remain `prototype`, and what bounded prototype shape is actually justified now?

---

## Active Thread Anchor

- `Primary Thread`: `coding-agent self-iteration and paper intake`
- `Supporting Thread`: `benchmark integrity and eval awareness`

Why these threads:

- `HLP-004` is fundamentally about whether paper-watch should be a real maintained lane in the repo instead of a one-off archive stub
- the benchmark-integrity thread matters because the prototype should improve how the repo interprets research and benchmark evidence, not just increase volume

Durable upstream anchors:

- [../knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md](../knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md)
- [../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md](../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md)
- [../radar/papers/2026-04-02-coding-agent-self-iteration-watch.md](../radar/papers/2026-04-02-coding-agent-self-iteration-watch.md)
- [../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md](../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md)
- [../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md](../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md)
- [../radar/weekly/2026-04-10.md](../radar/weekly/2026-04-10.md)

Result:

- `HLP-004` passes the thread-anchor bar

---

## Benchmark Evidence Check

- `Benchmark / Source`: the current paper-watch archive now covers long-horizon coding benchmarks, context retrieval, trajectory diagnosis, and observability limits rather than only one starter self-iteration note
- `Capability Signal`: the repo now has enough recurring paper-watch material to justify a maintained paper-aware lane for coding-agent evaluation and self-improvement topics
- `Eval-Trust Caveat`: more papers and benchmarks do not automatically mean better repo decisions; a broader paper lane can still collapse into a noisy intake stream if trust limits and selectivity are not preserved
- `Safe Local Principle`: prototype the paper-watch lane as a selective maintained surface that changes weekly synthesis and proposal review quality, not as a high-volume paper dump
- `Proposal Posture Effect`: this supports `prototype`, not `adopt`, because the lane already exists but still needs one clearer maintained operating shape before it should be treated as fully settled
- `Cargo-Cult Guard`: the repo should not imitate large-scale research monitoring programs or treat every arXiv paper as durable radar material

Result:

- `HLP-004` passes the benchmark-integrity bar because its current prototype pressure is about intake quality and selectivity, not benchmark enthusiasm

---

## Monitorability Criterion

Monitorability is not the core reason for this proposal, but it still shapes it indirectly.

- `Visible Stages`: the paper-watch archive now makes some research stages more visible by separating self-iteration, benchmark integrity, and trajectory diagnosis into durable notes
- `Hidden States`: the repo still lacks a stronger file-backed standard for deciding when a paper cluster is mature enough to change proposal posture rather than only enrich watch coverage
- `Safe Local Principle`: prototype one smaller paper-watch maintenance surface or review template before treating paper intake as fully stabilized

Result:

- `HLP-004` passes the monitorability criterion as a secondary consideration, not as its main driver

---

## Disposition Review

Current disposition should remain:

- `prototype`

Why:

- the repo has already proved paper-watch should exist
- the repo has not yet proved that paper-watch maintenance and review cadence are fully stable enough to call this lane complete
- the next justified move is a smaller prototype-level maintenance surface, not a status jump to `adopt`

---

## Recommended Prototype Shape

The smallest justified prototype from here is:

1. keep paper-watch selective and grouped rather than one-file-per-paper
2. require stronger linkage between paper-watch notes, weekly synthesis, and proposal review packets
3. add one smaller maintenance or review artifact only if future sessions start drifting back into ad hoc paper intake

---

## Review Outcome

This review strengthens `HLP-004`, but it does not upgrade it beyond `prototype`.

The real progress is practical:

- `HLP-004` is now backed by a first concrete proposal-review packet
- the repo can say what kind of prototype is actually justified
- the remaining work is about repeated operation and selective maintenance, not missing architectural intent
