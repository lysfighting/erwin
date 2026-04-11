# Active Knowledge Threads

> Status: Active
>
> Date: 2026-04-02
>
> Purpose: keep the current durable knowledge threads, evidence coverage, and proposal linkage in one file-backed surface so fresh sessions do not have to reconstruct knowledge state from chat memory or scattered notes.

---

## How To Use This File

Update this file when at least one of these becomes true:

- a radar weekly note identifies a new durable thread
- a feedback case changes what the repo now knows about `Harness`
- a knowledge note matures enough to change proposal posture
- a proposal is adopted, rejected, or superseded in a way that changes the live thread map

Do not copy raw intake here.
Do not turn this into a weekly journal.
This file is the compact registry of what the repo currently believes is worth carrying forward.

---

## Evidence Classes

Use this shorthand when reviewing a thread:

- `external`
  first-party engineering, papers, benchmarks, or other authoritative field signal
- `feedback`
  normalized real-usage cases under `docs/feedback/`
- `internal-gap`
  runtime, docs, automation, or process gaps discovered inside `erwin`

The goal is not perfect symmetry.
The goal is to make evidence density visible before proposal pressure grows.

---

## Active Threads

| Thread | Core Question | Evidence Classes | Current Durable View | Update When | Linked Proposals |
|--------|---------------|------------------|----------------------|-------------|------------------|
| `knowledge-system operating discipline` | How does `erwin` keep durable knowledge legible instead of falling back into `docs/research/`, proposal spam, or session-memory dependence? | `feedback`, `internal-gap`, `external` | The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure. | a new case shows the loop being bypassed, or a weekly synthesis/proposal cannot be traced back to a durable thread | `HLP-001`, `HLP-003` |
| `coding-agent self-iteration and paper intake` | What do first-party and paper evidence together imply about coding-agent self-improvement, and how should `erwin` watch that theme without cargo-culting vendors? | `external`, `internal-gap` | OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note. | a major benchmark shifts, or weekly synthesis finds stronger convergence between vendor claims and papers | `HLP-004`, `HLP-005` |
| `harness engineering as legibility + eval pressure` | Which outside signals should turn into eval/doc pressure rather than direct runtime expansion? | `external`, `feedback`, `internal-gap` | The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability. | repeated internal gaps show current docs/evals are not making the system legible enough, or new first-party evidence sharpens the eval pattern | `HLP-003` |
| `benchmark integrity and eval awareness` | When benchmark-heavy external signal enters the repo, how should `erwin` separate capability signal from eval-trust signal before it sharpens proposal pressure? | `external`, `feedback`, `internal-gap` | Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean. | a weekly synthesis, paper-watch note, or proposal relies heavily on benchmark results that could carry contamination, infrastructure, or gameability risk | `HLP-003`, `HLP-005` |

---

## Thread Anchors

### `knowledge-system operating discipline`

Primary anchors:

- [operating-model.md](operating-model.md)
- [../feedback/cases/HU-001-docs-research-became-default-knowledge-surface.md](../feedback/cases/HU-001-docs-research-became-default-knowledge-surface.md)
- [../feedback/cases/HU-002-proposals-were-emitted-before-knowledge-matured.md](../feedback/cases/HU-002-proposals-were-emitted-before-knowledge-matured.md)
- [../feedback/cases/HU-003-session-narrative-carried-knowledge-state.md](../feedback/cases/HU-003-session-narrative-carried-knowledge-state.md)
- [../harness-learning-loop.md](../harness-learning-loop.md)
- [../proposals/harness-proposal-ledger.md](../proposals/harness-proposal-ledger.md)

Current interpretation:

- the repo now has the right architecture, but it only stays real if repeated process failures are captured as feedback cases and reflected back into this thread
- the live risk is no longer “missing directories”; it is letting future work bypass the directories because the thread map is not maintained

### `coding-agent self-iteration and paper intake`

Primary anchors:

- [coding-agent-self-iteration-and-paper-intake-2026-04-02.md](coding-agent-self-iteration-and-paper-intake-2026-04-02.md)
- [../radar/weekly/2026-03-30.md](../radar/weekly/2026-03-30.md)
- [../radar/weekly/2026-04-03.md](../radar/weekly/2026-04-03.md)
- [../radar/weekly/2026-04-10.md](../radar/weekly/2026-04-10.md)
- [../radar/papers/2026-04-02-coding-agent-self-iteration-watch.md](../radar/papers/2026-04-02-coding-agent-self-iteration-watch.md)
- [../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md](../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md)
- [../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md](../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md)
- [../radar/sources.md](../radar/sources.md)
- [../proposals/harness-proposal-ledger.md](../proposals/harness-proposal-ledger.md)

Current interpretation:

- this thread is mature enough to shape intake and proposal posture, but not mature enough to justify direct runtime change
- the next concrete proof point is no longer just “start paper watch”; it is keeping recurring benchmark-aware paper watch and weekly synthesis behavior alive over time

### `harness engineering as legibility + eval pressure`

Primary anchors:

- [../radar/weekly/2026-03-30.md](../radar/weekly/2026-03-30.md)
- [../radar/weekly/2026-04-03.md](../radar/weekly/2026-04-03.md)
- [../radar/weekly/2026-04-10.md](../radar/weekly/2026-04-10.md)
- [harness-engineering-vs-internal-gaps-2026-04-02.md](harness-engineering-vs-internal-gaps-2026-04-02.md)
- [runtime-legibility-patterns-2026-04-02.md](runtime-legibility-patterns-2026-04-02.md)
- [../harness-learning-loop.md](../harness-learning-loop.md)
- [../proposals/harness-proposal-ledger.md](../proposals/harness-proposal-ledger.md)
- [../feedback/cases/HU-002-proposals-were-emitted-before-knowledge-matured.md](../feedback/cases/HU-002-proposals-were-emitted-before-knowledge-matured.md)
- [../feedback/cases/HU-004-productized-surfaces-outpaced-control-plane-maturity.md](../feedback/cases/HU-004-productized-surfaces-outpaced-control-plane-maturity.md)
- [../feedback/cases/HU-005-delegated-follow-up-needed-bounded-lifecycle-objects.md](../feedback/cases/HU-005-delegated-follow-up-needed-bounded-lifecycle-objects.md)
- [../feedback/cases/HU-006-active-council-needed-current-work-visibility.md](../feedback/cases/HU-006-active-council-needed-current-work-visibility.md)

Current interpretation:

- the repo should continue translating outside harness-engineering signal into stronger evals, docs, and review posture before widening runtime topology
- this thread becomes stronger when internal failures show that current legibility is still too weak, especially when leader-facing surfaces need a more precise bounded runtime object to stay governable

### `benchmark integrity and eval awareness`

Primary anchors:

- [benchmark-integrity-and-eval-awareness-2026-04-09.md](benchmark-integrity-and-eval-awareness-2026-04-09.md)
- [../feedback/cases/HU-007-benchmark-results-needed-eval-trust-interpretation.md](../feedback/cases/HU-007-benchmark-results-needed-eval-trust-interpretation.md)
- [../radar/weekly/2026-04-03.md](../radar/weekly/2026-04-03.md)
- [../radar/weekly/2026-04-10.md](../radar/weekly/2026-04-10.md)
- [../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md](../radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md)
- [../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md](../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md)
- [../proposals/harness-eval-candidate-ledger.md](../proposals/harness-eval-candidate-ledger.md)

Current interpretation:

- benchmark-heavy evidence should now be interpreted as both capability signal and eval-trust signal
- this thread is not yet a separate proposal, but it already sharpens how `HLP-003` and `HLP-005` should be reviewed

---

## Current Knowledge Gaps

- the paper-watch archive now has recurring notes for coding-agent benchmarks, but it still needs broader thematic coverage beyond the current coding-agent cluster
- the feedback layer now covers several runtime-facing scheduler/authority/council patterns, but it still needs more concrete cases from newer `Harness` usage rather than only the first big lifecycle cluster
- the weekly and proposal review gates have now been exercised more than once, but the real test is whether they hold across future syntheses without drifting back into loose narrative summaries
- the repo now has concrete proposal-review packets across `adopted`, `encode-as-eval`, `prototype`, and `watch`, so the missing work is no longer coverage of core disposition paths; it is sustained reuse, maintenance, and later pruning

---

## Next Build Targets

1. Keep adding normalized feedback cases whenever repeated `Harness` usage or repo process friction becomes visible.
2. Make each weekly synthesis update this file explicitly instead of relying on narrative continuity.
3. Only widen the proposal ledger when a thread in this file has enough evidence density to justify it.
4. Keep exercising the benchmark-evidence checklist whenever benchmark-heavy signal materially sharpens repo posture.
5. Keep exercising the monitorability-review checklist whenever trajectory diagnosis or intermediate observability materially sharpens repo posture.
6. Prefer the lighter monitorability review criterion when a full checklist would be heavier than the review actually needs.
7. Treat the current phase as operating and maintenance work, not as a reason to keep adding new documentation surface types by default.
