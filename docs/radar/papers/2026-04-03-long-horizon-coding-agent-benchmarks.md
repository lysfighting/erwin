# Paper Watch — Long-Horizon Coding-Agent Benchmarks

> Date: 2026-04-03
>
> Purpose: track the recent benchmark signal around long-horizon coding quality and context retrieval so the repo's coding-agent thread is grounded in newer benchmark pressure, not only first-party engineering posts.

---

## Why This Note Exists

The paper-watch archive now needs recurring benchmark notes, not just one self-iteration starter note.

This note extends:

- [2026-04-02-coding-agent-self-iteration-watch.md](2026-04-02-coding-agent-self-iteration-watch.md)
- [../../knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md](../../knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md)
- [../../knowledge/active-threads.md](../../knowledge/active-threads.md)

---

## Watch Set

### `SlopCodeBench`

- `SlopCodeBench: Benchmarking How Coding Agents Degrade Over Long-Horizon Iterative Tasks`
  - https://arxiv.org/abs/2603.24755
  - Why it matters: explicitly measures how coding agents degrade as they iteratively extend their own earlier code under evolving requirements, instead of only grading single-shot task completion.
  - Current signal: the paper reports that no evaluated agent solves any problem end-to-end across all checkpoints, and that structural erosion and verbosity increase over time.

### `ContextBench`

- `ContextBench: A Benchmark for Context Retrieval in Coding Agents`
  - https://arxiv.org/abs/2602.05892
  - Why it matters: isolates context retrieval quality as its own bottleneck and makes it visible through intermediate gold-context metrics, rather than only final-answer success.
  - Current signal: the paper argues that more elaborate agent scaffolding yields only marginal retrieval gains, and that agents still show a large gap between explored context and actually used context.

---

## Current Interpretation

- the benchmark frontier is getting sharper about long-horizon degradation, not just one-shot correctness
- context retrieval remains a separate weakness even when an agent seems to search broadly
- these papers strengthen the idea that `Harness` should keep caring about legibility, context discipline, and evaluable trajectory quality rather than only final outcomes

---

## Why This Matters For erwin

This benchmark cluster strengthens two active threads:

- `coding-agent self-iteration and paper intake`
- `harness engineering as legibility + eval pressure`

The immediate implication is not “copy benchmark scaffolds.”
The immediate implication is:

- keep benchmark-aware paper watch alive
- keep `encode-as-eval` pressure pointed at long-horizon quality and context use
- keep repo knowledge honest about how far current coding-agent performance still is from stable long-session execution

---

## Watch Next

- whether later papers make trajectory-level quality metrics more actionable for real agent harnesses
- whether benchmark work begins to connect directly to evaluator-agent design and structured handoff artifacts
- whether frontier vendors begin discussing long-horizon degradation or context retrieval bottlenecks as openly as they discuss one-shot benchmark wins
