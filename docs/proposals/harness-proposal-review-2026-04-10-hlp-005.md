# Harness Proposal Review — HLP-005 — 2026-04-10

> Status: Active
>
> Date: 2026-04-10
>
> Proposal: `HLP-005`
>
> Purpose: execute a concrete proposal review for `HLP-005` using the repo's paper-watch, benchmark-integrity, and monitorability review surfaces.

---

## Review Scope

This review applies the current proposal surfaces to:

- [harness-proposal-ledger.md](harness-proposal-ledger.md)
- [../radar/weekly/2026-04-10.md](../radar/weekly/2026-04-10.md)
- [../knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md](../knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md)
- [../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md](../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md)

The goal is to answer one narrower question:

- should `HLP-005` still remain `watch`, and if so, what exactly is the repo watching for now?

---

## Active Thread Anchor

- `Primary Thread`: `coding-agent self-iteration and paper intake`
- `Supporting Thread`: `benchmark integrity and eval awareness`

Why these threads:

- `HLP-005` is about how the repo should interpret the evolving split between eval-driven self-improvement, heavy internal dogfooding, benchmark pressure, and paper-watch growth
- benchmark-integrity matters because the strongest newer signal in this area is no longer only product posture; it also concerns whether public benchmark and observability evidence should be trusted enough to change local posture

Durable upstream anchors:

- [../knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md](../knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md)
- [../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md](../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md)
- [../radar/papers/2026-04-02-coding-agent-self-iteration-watch.md](../radar/papers/2026-04-02-coding-agent-self-iteration-watch.md)
- [../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md](../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md)
- [../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md](../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md)
- [../feedback/cases/HU-007-benchmark-results-needed-eval-trust-interpretation.md](../feedback/cases/HU-007-benchmark-results-needed-eval-trust-interpretation.md)

Result:

- `HLP-005` passes the thread-anchor bar

---

## Benchmark Evidence Check

- `Benchmark / Source`: `SlopCodeBench`, `ContextBench`, Anthropic's infrastructure-noise note, and the trajectory-diagnosis paper cluster summarized in the current paper-watch archive
- `Capability Signal`: the repo now has stronger evidence that coding-agent quality depends on long-horizon degradation resistance, context use, stage-level diagnosis, and validation behavior rather than only one-shot benchmark performance
- `Eval-Trust Caveat`: benchmark scores and benchmark narratives are still unstable because infrastructure choices, output-only reporting, and benchmark-specific structures can distort what those results really mean
- `Safe Local Principle`: keep using public benchmark and paper signal to sharpen watch criteria and local eval questions, not to import public task structures or vendor benchmark rhetoric directly
- `Proposal Posture Effect`: this supports a sharper `watch`, not an upgrade to `encode-as-eval` or `adopt`
- `Cargo-Cult Guard`: the repo should not treat every new coding-agent benchmark as a reason to widen `Harness`; the correct move is to keep translating benchmark signal into better watch criteria first

Result:

- `HLP-005` passes the benchmark-integrity bar because its current role is interpretive watchfulness, not premature operationalization

---

## Monitorability Criterion

- `Visible Stages`: recent evidence now clearly highlights search, read, edit, validation, and oversight as meaningful stages for understanding coding-agent behavior
- `Hidden States`: public coding-agent reports still often hide what the agent actually saw, which runtime assumptions held, and why a result failed or succeeded, so output-only conclusions remain too weak for stronger local action
- `Safe Local Principle`: the repo should watch for external work that makes those stages and assumptions more inspectable before deciding that a coding-agent result should shape local eval design or runtime direction

Result:

- `HLP-005` passes the monitorability criterion because it keeps the repo at the level of interpretive watch rather than pretending these external signals are already strong enough to dictate local implementation

---

## Disposition Review

Current disposition should remain:

- `watch`

Why:

- the repo now has better paper-watch and benchmark-integrity coverage than before
- that stronger coverage sharpens what should be watched, but it does not yet force one bounded local eval candidate beyond the more general `HLP-003` pressure
- the right immediate use of this signal is continued paper-watch, benchmark-integrity interpretation, and monitorability-aware reading of future external evidence

---

## What The Repo Is Actually Watching For

`HLP-005` is now watching for convergence across at least three fronts:

1. whether first-party engineering posts start exposing trajectory diagnosis and monitorability as explicitly as papers already do
2. whether benchmark-rich evidence becomes more trustworthy, not just more abundant
3. whether a recurring pattern emerges that justifies a narrower local eval candidate beyond the broader `HLP-003` review pressure

---

## Recommended Next Moves

1. Keep `HLP-005` at `watch` and continue updating it through paper-watch notes and weekly synthesis rather than opening a separate implementation track.
2. Revisit this proposal only when external signal converges strongly enough to justify one smaller local eval object instead of a broad thematic watch.
3. Keep using benchmark-integrity and monitorability review surfaces so the watch stays disciplined rather than turning into a benchmark-news inbox.

---

## Review Outcome

This review strengthens `HLP-005`, but it does not advance it beyond `watch`.

The useful progress is interpretive:

- `HLP-005` is now backed by a first worked proposal-review packet
- the repo can say what this watch is actually watching for
- future movement should depend on stronger convergence, not on novelty alone
