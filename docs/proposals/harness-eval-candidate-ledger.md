# Harness Eval Candidate Ledger

> Status: Active
>
> Date: 2026-04-02
>
> Purpose: turn `encode-as-eval` proposal posture into a concrete, reviewable queue of eval, audit, and review-surface candidates before implementation work starts.

This ledger is the bridge between:

- durable knowledge threads
- reviewed proposal dispositions such as `encode-as-eval`
- actual future eval, audit, smoke, or review-criteria work

It is not a raw idea inbox.
It records only candidates that are already grounded in active knowledge and reviewed proposal posture.

---

## Candidate Fields

Each serious eval candidate should answer:

- `candidate_id`
- `date`
- `proposal_id`
- `thread`
- `candidate_type`
- `target_surface`
- `assertion`
- `why_this_should_be_checked`
- `source_anchors`
- `status`

### Recommended `candidate_type`

- `review-criteria`
- `audit-check`
- `smoke-test`
- `docs-gate`
- `automation-check`

### Recommended `status`

- `identified`
- `ready-for-encoding`
- `encoded`
- `superseded`

---

## Active Candidates

| ID | Date | Proposal | Thread | Candidate Type | Target Surface | Status | Assertion |
|----|------|----------|--------|----------------|----------------|--------|-----------|
| HEL-001 | 2026-04-02 | `HLP-003` | `harness engineering as legibility + eval pressure` | `review-criteria` | weekly synthesis + proposal review | `encoded` | A weekly synthesis or proposal candidate should not recommend stronger `Harness` evolution without citing the named knowledge thread it updates and the upstream evidence that supports it. |
| HEL-002 | 2026-04-02 | `HLP-003` | `harness engineering as legibility + eval pressure` | `docs-gate` | handoff + documentation map + current state | `encoded` | Fresh-session entrypoints should continue to expose the live knowledge registry and paper-watch surfaces whenever they are part of the active learning loop. |
| HEL-003 | 2026-04-02 | `HLP-003` | `harness engineering as legibility + eval pressure` | `audit-check` | knowledge/proposal review | `encoded` | No doc or proposal should imply that `Harness` governance maturity is complete when the current gap/drift audits still mark ownership policy, authority adjudication, or council automation as partial. |
| HEL-004 | 2026-04-02 | `HLP-005` | `coding-agent self-iteration and paper intake` | `automation-check` | weekly synthesis + paper watch | `encoded` | Research-heavy threads should cite the relevant paper-watch note or explicitly say that no durable paper-watch update exists, instead of depending only on first-party signal. |
| HEL-005 | 2026-04-09 | `HLP-003` + `HLP-005` | `benchmark integrity and eval awareness` | `review-criteria` | weekly synthesis + proposal review + benchmark evidence checklist | `encoded` | When benchmark-heavy external evidence is used, review should distinguish capability signal from eval-trust signal instead of treating benchmark results as pure capability truth. |
| HEL-006 | 2026-04-10 | `HLP-003` + `HLP-005` | `harness engineering as legibility + eval pressure` + `benchmark integrity and eval awareness` | `review-criteria` | future eval design + weekly synthesis + harness-facing review + monitorability checklist + first audit artifact + light review criterion | `encoded` | When external signal sharpens monitorability pressure, local eval or audit design should ask which trajectory stages and runtime decisions remain visible enough to diagnose failure rather than grading only final output. |

---

## Candidate Notes

### `HEL-001` — Knowledge-thread traceability check

- `candidate_id`: `HEL-001`
- `date`: `2026-04-02`
- `proposal_id`: `HLP-003`
- `thread`: `harness engineering as legibility + eval pressure`
- `candidate_type`: `review-criteria`
- `target_surface`: `docs/radar/weekly/` and proposal review
- `assertion`: a serious weekly or proposal candidate should name the active knowledge thread it updates and cite upstream evidence from `feedback`, `knowledge`, or both
- `why_this_should_be_checked`: the knowledge system now has real threads and cases; without this check, the repo could drift back to narrative-only weekly synthesis and weakly grounded proposals
- `source_anchors`: `docs/knowledge/active-threads.md`, `docs/radar/weekly/2026-03-30.md`, `docs/proposals/harness-proposal-ledger.md`
- `status`: `encoded`

### `HEL-002` — Fresh-session discoverability gate

- `candidate_id`: `HEL-002`
- `date`: `2026-04-02`
- `proposal_id`: `HLP-003`
- `thread`: `harness engineering as legibility + eval pressure`
- `candidate_type`: `docs-gate`
- `target_surface`: `docs/harness-session-handoff.md`, `docs/harness-documentation-map.md`, `_project/CURRENT.md`
- `assertion`: fresh-session entrypoints should expose the active knowledge registry and paper-watch archive when they are part of the live learning loop
- `why_this_should_be_checked`: a legible repo is not just a runtime goal; it also requires that the right durable context is easy to discover from canonical entrypoints
- `source_anchors`: `docs/harness-session-handoff.md`, `docs/harness-documentation-map.md`, `docs/knowledge/active-threads.md`, `docs/radar/papers/README.md`
- `status`: `encoded`

### `HEL-003` — Maturity-overclaim audit

- `candidate_id`: `HEL-003`
- `date`: `2026-04-02`
- `proposal_id`: `HLP-003`
- `thread`: `harness engineering as legibility + eval pressure`
- `candidate_type`: `audit-check`
- `target_surface`: knowledge/proposal review and future harness-facing docs
- `assertion`: docs and proposals should not imply that control-plane maturity is complete when the current audits still mark key governance areas as partial
- `why_this_should_be_checked`: visible productization can outrun deeper governance completeness, and the knowledge system now has a normalized feedback case recording that risk
- `source_anchors`: `docs/feedback/cases/HU-004-productized-surfaces-outpaced-control-plane-maturity.md`, `docs/knowledge/harness-engineering-vs-internal-gaps-2026-04-02.md`, `docs/research/harness-gap-audit-2026-03-29.md`, `docs/research/harness-research-drift-audit-2026-03-29.md`
- `status`: `encoded`

### `HEL-004` — Paper-watch coverage check

- `candidate_id`: `HEL-004`
- `date`: `2026-04-02`
- `proposal_id`: `HLP-005`
- `thread`: `coding-agent self-iteration and paper intake`
- `candidate_type`: `automation-check`
- `target_surface`: `docs/radar/weekly/` and `docs/radar/papers/`
- `assertion`: a research-heavy weekly synthesis should either cite the relevant durable paper-watch note or say that no durable paper-watch update exists yet
- `why_this_should_be_checked`: the repo now has a paper-watch archive, so research-heavy knowledge should no longer depend only on first-party field radar
- `source_anchors`: `docs/radar/papers/README.md`, `docs/radar/papers/2026-04-02-coding-agent-self-iteration-watch.md`, `docs/knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md`
- `status`: `encoded`

### `HEL-005` — Benchmark-integrity review check

- `candidate_id`: `HEL-005`
- `date`: `2026-04-09`
- `proposal_id`: `HLP-003`, `HLP-005`
- `thread`: `benchmark integrity and eval awareness`
- `candidate_type`: `review-criteria`
- `target_surface`: weekly synthesis, proposal review, and the benchmark-evidence checklist
- `assertion`: when benchmark-heavy external evidence appears, the review should say what capability signal it extracts and what eval-trust caveat travels with that signal
- `why_this_should_be_checked`: benchmark results can influence knowledge and proposal posture while still carrying contamination or gameability risk, so the encoded rule now also needs a reusable checklist artifact instead of only high-level gate language
- `source_anchors`: `docs/knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md`, `docs/radar/weekly/2026-04-03.md`, `docs/radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md`, `docs/proposals/harness-benchmark-evidence-checklist.md`
- `status`: `encoded`

### `HEL-006` — Monitorability and trajectory-visibility check

- `candidate_id`: `HEL-006`
- `date`: `2026-04-10`
- `proposal_id`: `HLP-003`, `HLP-005`
- `thread`: `harness engineering as legibility + eval pressure`, `benchmark integrity and eval awareness`
- `candidate_type`: `review-criteria`
- `target_surface`: future eval design, weekly synthesis, harness-facing review, the monitorability checklist, the first worked audit artifact, and the lighter review criterion
- `assertion`: when external evidence points to stage-level diagnosis, monitorability, or intermediate observability, local eval and audit design should say which search, read, edit, validation, or decision stages remain visible enough to diagnose failure
- `why_this_should_be_checked`: outcome-only review is increasingly looking too lossy for coding-agent evaluation, and the repo now has durable weekly and paper-watch evidence that monitorability is becoming a reusable design principle rather than a one-off paper insight, so the rule now needs a reusable checklist artifact instead of only a ledger entry
- `source_anchors`: `docs/radar/weekly/2026-04-10.md`, `docs/radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md`, `docs/knowledge/active-threads.md`, `docs/proposals/harness-monitorability-review-checklist.md`, `docs/proposals/harness-monitorability-audit-2026-04-10.md`, `docs/proposals/harness-monitorability-review-criterion.md`
- `status`: `encoded`
