# Harness Monitorability Review Checklist

> Status: Active
>
> Date: 2026-04-10
>
> Purpose: turn monitorability and trajectory-visibility review from a general concern into a reusable checklist whenever `Harness`-facing knowledge, proposal, or eval design depends on stage-level diagnosis or intermediate observability.

---

## Why This Exists

`HEL-006` now says the repo should stop treating outcome-only evaluation as sufficient whenever external signal points to stage-level diagnosis, monitorability, or intermediate observability.

This checklist is the concrete review artifact for that rule.

Use it to keep monitorability review honest about both:

- `visible stages`
  which search, read, edit, validation, planning, or decision stages stay inspectable enough to diagnose what happened
- `observability gaps`
  which important states remain hidden, collapsed, or too lossy to support reliable diagnosis from final output alone

---

## When To Use This

Use this checklist when a weekly synthesis, knowledge note, proposal review, or local eval idea materially depends on:

- trajectory diagnosis
- stage-level coding-agent analysis
- intermediate observability or monitorability
- output-only feedback limitations
- monitoring or oversight claims for realistic coding-agent runtimes

If monitorability is only mentioned in passing and does not change repo posture, a full pass is optional.

---

## Minimum Review Record

When this checklist applies, the review should record answers for each field below.

### 1. Source And Runtime Surface

- Which paper, engineering note, weekly synthesis, or internal signal introduced the monitorability claim?
- Which runtime surface is actually under discussion: planning, search, read, edit, validation, approval, or some other stage?

### 2. Visible Stages

- Which stages or decisions remain visible enough to inspect directly?
- Is the available visibility strong enough to distinguish failure to find context, failure to interpret context, failure to edit correctly, failure to validate, or failure to stop and ask?

### 3. Observability Gap

- Which important internal states are still hidden or collapsed into a lossy final artifact?
- Why would output-only review be too weak here?
- What kind of misdiagnosis becomes likely if the repo treats final output as sufficient evidence?

### 4. Safe Local Principle

- What local principle should the repo encode from this signal?
- Why is that principle safer than copying the source's exact benchmark instrumentation, internal monitor stack, or product workflow?

Examples:

- require a review note to say which search/read/edit stages were visible before drawing conclusions
- treat validation behavior as its own inspectable stage instead of folding it into final success
- ask what runtime assumptions or resource conditions shaped the observed behavior before trusting the conclusion

### 5. Proposal Or Eval Effect

- Does this signal justify `watch`, `knowledge gap`, `encode-as-eval`, `prototype`, or `adopt`?
- If it sharpens future eval design, what bounded review or audit surface should change first?

### 6. Cargo-Cult Guard

- What would a cargo-cult move look like here?
- What exact external instrumentation, benchmark setup, or oversight pattern is the repo refusing to copy directly?

---

## Fast Pass / Fail Heuristic

The review is usually strong enough when it can say all of the following:

- which important stages were visible
- which important states remained hidden
- why output-only interpretation would be too weak
- what smaller local principle should be encoded first

The review is too weak when it does any of the following:

- treats final success or failure as self-explanatory
- claims monitorability improved without naming what became more visible
- jumps from an observability paper or internal monitor writeup straight to `adopt`
- says `encode-as-eval` without naming the first bounded review or audit surface to tighten

---

## Suggested Review Template

```md
## Monitorability Review Check

- Source / Runtime Surface:
- Visible Stages:
- Observability Gap:
- Safe Local Principle:
- Proposal Or Eval Effect:
- Cargo-Cult Guard:
```

---

## Canonical Anchors

- [../radar/weekly/2026-04-10.md](../radar/weekly/2026-04-10.md)
- [../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md](../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md)
- [../knowledge/active-threads.md](../knowledge/active-threads.md)
- [harness-monitorability-review-criterion.md](harness-monitorability-review-criterion.md)
- [harness-eval-candidate-ledger.md](harness-eval-candidate-ledger.md)
- [harness-review-gates.md](harness-review-gates.md)
