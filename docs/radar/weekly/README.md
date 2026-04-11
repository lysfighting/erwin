# Weekly Synthesis Guide

> Status: Active
>
> Date: 2026-04-02
>
> Purpose: define the minimum structure and review gates for durable weekly radar syntheses so they strengthen active knowledge threads instead of becoming another narrative summary layer.

---

## What A Weekly Synthesis Is For

A weekly synthesis should:

- aggregate the strongest field signal from the week
- compare that signal against existing knowledge and feedback
- name which durable knowledge threads should update
- surface knowledge gaps explicitly
- emit proposal pressure only when the evidence is mature enough

It should not behave like:

- a daily inbox dump
- a raw proposal generator
- a replacement for `docs/knowledge/active-threads.md`

---

## Required Sections

Every serious weekly synthesis should include:

- `Field Summary`
- `Theme Deep Dive`
- `Tooling Watch` when relevant
- `Implications For erwin / Harness / Skills`
- `Knowledge Threads To Update`
- `Knowledge Gaps`
- `Harness Evolution Proposals` when warranted
- `Open Questions`

---

## Review Gates

### Gate 1: Thread Traceability

Each named update in `Knowledge Threads To Update` should:

- cite the active knowledge thread by name
- say why that thread, rather than some other one, is the right durable target
- point back to upstream evidence in `knowledge`, `feedback`, or both when proposal pressure is involved

This encodes `HEL-001`.

### Gate 2: Research-Heavy Coverage

If a weekly synthesis materially discusses:

- papers
- benchmarks
- self-improvement methods
- coding-agent research automation

then it should either:

- cite the relevant note under `docs/radar/papers/`
- or say explicitly that no durable paper-watch note exists yet

This encodes `HEL-004`.

### Gate 3: Benchmark Integrity

If a weekly synthesis uses benchmark-heavy evidence, it should say both:

- what capability signal it is extracting
- what eval-trust caveat travels with that signal

This is especially important for:

- long-horizon coding benchmarks
- context-retrieval benchmarks
- browse or web-agent benchmarks

When this gate applies, prefer using:

- [../../proposals/harness-benchmark-evidence-checklist.md](../../proposals/harness-benchmark-evidence-checklist.md)

This encodes `HEL-005`.

### Gate 4: Monitorability Claims Must Name What Stays Visible

If a weekly synthesis relies on trajectory diagnosis, monitorability, or intermediate observability claims, it should say:

- which stages or decisions stayed visible enough to inspect
- which important states remained hidden or too lossy for output-only interpretation
- which smaller local review or eval principle should be encoded instead of copying the source's exact instrumentation

When this gate applies, prefer using:

- [../../proposals/harness-monitorability-review-checklist.md](../../proposals/harness-monitorability-review-checklist.md)

This encodes `HEL-006`.

### Gate 5: Proposal Restraint

A `Harness Evolution Proposal` section should appear only when the weekly note can explain:

- which active knowledge thread matured
- which evidence classes support the move
- why the recommended disposition is not just fresh enthusiasm

If that bar is not met, the weekly note should stay at `watch` or `knowledge gap` level.

---

## Minimum Proposal Shape

Each weekly proposal should include:

- `Title`
- `Source`
- `Observed Pattern`
- `Why It Matters For Harness`
- `Recommended Disposition`
- `What We Would Change If Adopted`
- `Why We Should Not Cargo-Cult This`

When the disposition is `encode-as-eval`, the next durable stop should usually be:

- [../../proposals/harness-eval-candidate-ledger.md](../../proposals/harness-eval-candidate-ledger.md)
