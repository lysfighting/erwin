# Harness Benchmark Evidence Checklist

> Status: Active
>
> Date: 2026-04-09
>
> Purpose: turn benchmark-integrity review from a general rule into a reusable checklist whenever `Harness`-facing knowledge or proposal work relies on benchmark-heavy external evidence.

---

## Why This Exists

`HEL-005` is now encoded in the weekly and proposal review gates, but those gates still need a concrete review artifact so benchmark-heavy evidence does not get interpreted differently every session.

This checklist is that artifact.

Use it to keep benchmark-heavy review honest about both:

- `capability signal`
  what the benchmark result appears to say about agent behavior
- `eval-trust signal`
  how trustworthy the benchmark still seems as a durable measure instead of a contamination-sensitive or gameable target

---

## When To Use This

Use this checklist when a weekly synthesis, knowledge note, or proposal review materially depends on:

- long-horizon coding benchmarks
- context-retrieval benchmarks
- browse or web-agent benchmarks
- benchmark-driven paper syntheses
- any external comparison where a benchmark result is carrying real proposal pressure

If the benchmark is only mentioned in passing and does not influence repo posture, a full pass is optional.

---

## Minimum Review Record

When this checklist applies, the review should record answers for each field below.

### 1. Benchmark And Source

- Which benchmark or benchmark-heavy source is being used?
- Where is the source note, paper-watch note, or weekly artifact that introduced it?

### 2. Extracted Capability Signal

- What concrete capability signal are we extracting?
- Is that signal about planning, long-horizon degradation, context retrieval, browsing, or something else?
- What should the repo learn from that signal without over-generalizing?

### 3. Eval-Trust Caveat

- Why might this benchmark be a weak or unstable proxy?
- Could the task structure be easy to game, overfit, contaminate, or memorize?
- Are the benchmark conditions too narrow to support strong repo-level conclusions?

### 4. Safe Local Principle

- If this evidence sharpens local action, what principle should we encode locally?
- Why is that principle safer than copying the benchmark's public task shape?

Examples:

- check for degradation across iterative work rather than cloning a benchmark trace
- check whether the right context was surfaced and used rather than reproducing a public context-retrieval task
- check whether a measure can be gamed rather than assuming benchmark score equals maturity

### 5. Proposal Posture Effect

- Does this evidence justify `watch`, `knowledge gap`, `encode-as-eval`, `prototype`, or `adopt`?
- What stronger disposition is still blocked by weak eval trust?

### 6. Cargo-Cult Guard

- What would a cargo-cult interpretation look like here?
- What specific move is the repo refusing to make because the benchmark is not trustworthy enough?

---

## Fast Pass / Fail Heuristic

The review is usually strong enough when it can say all of the following:

- the benchmark exposed a real weakness worth tracking
- the benchmark itself still carries trust limits
- the repo is encoding the revealed principle, not blindly copying the public task
- the recommended disposition stays bounded by the evidence

The review is too weak when it does any of the following:

- treats a benchmark result as pure capability truth
- skips contamination or gameability discussion entirely
- jumps from benchmark excitement straight to `adopt`
- says `encode-as-eval` without naming the safer local principle to encode

---

## Suggested Review Template

```md
## Benchmark Evidence Check

- Benchmark / Source:
- Capability Signal:
- Eval-Trust Caveat:
- Safe Local Principle:
- Proposal Posture Effect:
- Cargo-Cult Guard:
```

---

## Canonical Anchors

- [../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md](../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md)
- [../feedback/cases/HU-007-benchmark-results-needed-eval-trust-interpretation.md](../feedback/cases/HU-007-benchmark-results-needed-eval-trust-interpretation.md)
- [harness-eval-candidate-ledger.md](harness-eval-candidate-ledger.md)
- [harness-review-gates.md](harness-review-gates.md)
