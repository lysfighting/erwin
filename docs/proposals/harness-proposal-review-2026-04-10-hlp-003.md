# Harness Proposal Review — HLP-003 — 2026-04-10

> Status: Active
>
> Date: 2026-04-10
>
> Proposal: `HLP-003`
>
> Purpose: execute a concrete proposal review for `HLP-003` using the repo's encoded thread, benchmark-integrity, and monitorability review surfaces.

---

## Review Scope

This review applies the current proposal surfaces to:

- [harness-proposal-ledger.md](harness-proposal-ledger.md)
- [../radar/weekly/2026-04-10.md](../radar/weekly/2026-04-10.md)
- [../knowledge/harness-engineering-vs-internal-gaps-2026-04-02.md](../knowledge/harness-engineering-vs-internal-gaps-2026-04-02.md)
- [../knowledge/runtime-legibility-patterns-2026-04-02.md](../knowledge/runtime-legibility-patterns-2026-04-02.md)

The goal is to answer one narrower question:

- should `HLP-003` still remain `encode-as-eval`, and if so, what kind of eval/review pressure is actually justified right now?

---

## Active Thread Anchor

- `Primary Thread`: `harness engineering as legibility + eval pressure`
- `Supporting Thread`: `benchmark integrity and eval awareness`

Why these threads:

- `HLP-003` is fundamentally about whether outside agent-system signal should first sharpen review, eval, and documentation discipline instead of directly widening runtime topology
- the benchmark-integrity thread now matters because some of the strongest recent external signal is no longer only about new runtime capability, but about how trustworthy or diagnosable the evidence is

Durable upstream anchors:

- [../knowledge/harness-engineering-vs-internal-gaps-2026-04-02.md](../knowledge/harness-engineering-vs-internal-gaps-2026-04-02.md)
- [../knowledge/runtime-legibility-patterns-2026-04-02.md](../knowledge/runtime-legibility-patterns-2026-04-02.md)
- [../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md](../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md)
- [../feedback/cases/HU-004-productized-surfaces-outpaced-control-plane-maturity.md](../feedback/cases/HU-004-productized-surfaces-outpaced-control-plane-maturity.md)
- [../feedback/cases/HU-005-delegated-follow-up-needed-bounded-lifecycle-objects.md](../feedback/cases/HU-005-delegated-follow-up-needed-bounded-lifecycle-objects.md)
- [../feedback/cases/HU-006-active-council-needed-current-work-visibility.md](../feedback/cases/HU-006-active-council-needed-current-work-visibility.md)
- [../feedback/cases/HU-007-benchmark-results-needed-eval-trust-interpretation.md](../feedback/cases/HU-007-benchmark-results-needed-eval-trust-interpretation.md)

Result:

- `HLP-003` passes the thread-anchor bar

---

## Benchmark Evidence Check

- `Benchmark / Source`: Anthropic's infrastructure-noise note, `SlopCodeBench`, `ContextBench`, and the trajectory-diagnosis paper cluster already grouped in the 2026-04-03 and 2026-04-10 weekly/paper-watch notes
- `Capability Signal`: recent public signal supports the claim that long-running coding-agent quality depends on trajectory discipline, context use, validation behavior, and runtime legibility rather than only one-shot task success
- `Eval-Trust Caveat`: the same signal also says benchmark results can be distorted by infrastructure conditions and outcome-only reporting, so public scores do not justify strong architecture conclusions by themselves
- `Safe Local Principle`: encode local review and eval pressure around visible stages, runtime assumptions, and diagnosable failure modes rather than importing public benchmark structures or leaderboard thresholds
- `Proposal Posture Effect`: this supports `encode-as-eval` and review-surface tightening, not direct runtime adoption
- `Cargo-Cult Guard`: the repo should not copy benchmark scaffolds or evaluation dashboards from outside labs just because they surface useful weaknesses

Result:

- `HLP-003` passes the benchmark-integrity bar because its current posture stays on principle extraction rather than benchmark copying

---

## Monitorability Criterion

- `Visible Stages`: the current evidence is strong enough to name search, read, edit, validation, planning, and approval/oversight as stages that matter for diagnosis
- `Hidden States`: final output alone still hides whether the agent found the right context, interpreted it correctly, or acted under distorted runtime assumptions, so outcome-only interpretation is too weak
- `Safe Local Principle`: the repo should prefer smaller review and audit surfaces that ask what stages were visible, what remained hidden, and what runtime assumptions shaped the interpretation

Result:

- `HLP-003` passes the monitorability criterion because it is now being interpreted through review surfaces that explicitly keep stage visibility and hidden-state limits in view

---

## Maturity Overclaim Check

`HLP-003` should still not be read as proof that `Harness` governance maturity is complete.

Current review judgment:

- stronger review and eval posture is justified
- stronger runtime claims are not yet justified
- external convergence around monitorability does not erase the repo's own gap audits or runtime evidence about partial governance maturity

Result:

- `HLP-003` passes the maturity-overclaim bar only as long as it remains review/eval-first

---

## Disposition Review

Current disposition should remain:

- `encode-as-eval`

Why:

- the current evidence is mature enough to justify stronger review and eval discipline
- the current evidence is not mature enough to justify direct runtime widening
- the repo now has concrete downstream surfaces for that narrower move: benchmark checklist, monitorability checklist, monitorability criterion, and the first worked audit artifact

---

## Recommended Next Moves

The smallest justified moves after this review are:

1. Keep `HLP-003` on `encode-as-eval` and continue using the benchmark and monitorability review surfaces in later weekly and proposal work.
2. Turn the monitorability criterion into one narrower repeated review question for future `Harness`-facing eval proposals:
   `Which stages were visible, which states remained hidden, and which runtime assumptions shaped the conclusion?`
3. Avoid widening runtime topology based only on current public signal, because the current evidence still says review and interpretability pressure are ahead of architecture expansion.

---

## Review Outcome

This review strengthens `HLP-003`, but it does not upgrade it.

The real progress is narrower:

- `HLP-003` is now backed by a first worked proposal-review packet
- the repo can explain why `encode-as-eval` is still the right posture
- the next useful work is more reuse of the encoded review surfaces, not a jump to runtime adoption
