# Benchmark Integrity And Eval Awareness

> Date: 2026-04-09
>
> Purpose: turn the recent weekly signal about benchmark awareness, long-horizon degradation, and context-retrieval limits into a durable knowledge note that separates capability signal from evaluation-trust signal.

---

## Bottom Line

Recent external signal is teaching two related lessons at once:

1. coding-agent benchmarks are getting better at exposing long-horizon and context-retrieval weaknesses
2. stronger agents can also become better at exploiting benchmark structure itself

That means `erwin` should not treat benchmark results as pure capability signal.

The repo should treat benchmark-heavy evidence as a two-part object:

- `capability signal`
  what the benchmark suggests the agent can or cannot do
- `eval-trust signal`
  how much confidence we should place in that benchmark as a durable measure rather than a gameable target

This distinction matters for `Harness` because a knowledge-first learning loop can still drift if its external evidence is not honest about evaluation integrity.

---

## What Strengthened This Thread

The strongest recent inputs are:

- [../radar/weekly/2026-04-03.md](../radar/weekly/2026-04-03.md)
- [../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md](../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md)
- [coding-agent-self-iteration-and-paper-intake-2026-04-02.md](coding-agent-self-iteration-and-paper-intake-2026-04-02.md)

What they add together:

- `SlopCodeBench` sharpens long-horizon quality drift
- `ContextBench` sharpens context-retrieval bottlenecks
- Anthropic's BrowseComp eval-awareness writeup sharpens the risk that a capable agent can start optimizing against benchmark structure itself

These are different signals, but they converge on the same question:

**when should `erwin` trust a benchmark as a stable input, and when should it treat the benchmark itself as part of the problem?**

---

## Why This Is Not Just A Benchmark Note

This is not only about research benchmarking.

It also affects how the repo should interpret:

- weekly syntheses
- proposal pressure
- `encode-as-eval` dispositions
- future review gates for external evidence

If benchmark integrity is ignored, the repo can accidentally:

- over-trust a fashionable benchmark result
- import cargo-cult evaluation logic
- confuse benchmark gaming with real harness maturity

So this note belongs in `knowledge/`, not only in `radar/papers/`.

---

## Practical Interpretation For erwin

### 1. Benchmark signal should stay paired with benchmark-trust commentary

When the repo cites a benchmark-heavy external signal, it should try to say:

- what capability signal it is extracting
- what evaluation-trust caveat travels with that signal

This is especially important for:

- long-horizon coding-agent benchmarks
- context-retrieval benchmarks
- agentic web or browse benchmarks

### 2. `encode-as-eval` should not mean “copy the benchmark”

When a proposal lands on `encode-as-eval`, the safest next move is often:

- encode the principle the benchmark exposed
- not the benchmark's exact public task structure

For example:

- long-horizon drift pressure should become a local check about degradation over iterative work
- context retrieval pressure should become a local check about whether the right context was surfaced and actually used
- eval-awareness pressure should become a review question about whether the measure can be gamed or overly tailored

### 3. Weekly syntheses should keep this distinction visible

The first real weekly synthesis already surfaced this issue.

Future weekly notes should keep separating:

- “this benchmark suggests a real weakness”
- from
- “this benchmark is itself still a trustworthy measuring instrument”

---

## Connection To Existing Threads

This note strengthens:

- `coding-agent self-iteration and paper intake`
- `harness engineering as legibility + eval pressure`

It strengthens the first thread because benchmark watch is no longer only about discovering papers.
It is also about deciding whether the benchmark evidence should shape proposal pressure.

It strengthens the second thread because evaluation honesty is part of harness legibility.
A harness that reasons sloppily about evaluation evidence is still a weak harness, even if its runtime looks polished.

---

## Recommended Next Move

The repo should add a concrete eval candidate that asks:

- when external evidence relies on benchmark results, has the review explicitly separated capability signal from eval-trust signal?

That would let `erwin` move this from a smart note into a repeatable review behavior.
