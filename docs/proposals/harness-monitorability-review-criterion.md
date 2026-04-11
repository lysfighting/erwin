# Harness Monitorability Review Criterion

> Status: Active
>
> Date: 2026-04-10
>
> Purpose: provide the smallest reusable review criterion derived from the monitorability checklist and first audit artifact so future `Harness`-facing eval or proposal work can be screened quickly.

---

## Why This Exists

The repo now has:

- a monitorability review checklist
- a first worked monitorability audit artifact

What it still needed was a lighter review surface for repeated use.

This criterion is that smaller surface.

Use it when a full checklist would be too heavy but a serious `Harness`-facing claim still depends on monitorability, trajectory diagnosis, or output-only review limits.

---

## The Three Questions

Any serious review should be able to answer all three:

### 1. Visible Stages

Which stages or decisions were actually visible enough to inspect?

Minimum expectation:

- name at least one concrete stage such as search, read, edit, validation, planning, or approval

### 2. Hidden States

Which important states remained hidden or too lossy for output-only interpretation?

Minimum expectation:

- explain why final output alone is not enough to diagnose what really happened

### 3. Safe Local Principle

What smaller local principle should the repo encode from this signal instead of copying the external instrumentation or workflow directly?

Minimum expectation:

- name one bounded review, audit, or eval principle the repo can reuse

---

## Pass / Fail

The review passes this criterion only if it can answer all three questions clearly.

It fails when it does any of the following:

- treats final success or failure as self-explanatory
- says monitorability matters without naming what was visible
- names hidden states but proposes no bounded local principle
- jumps from observability signal straight to `adopt`

---

## Suggested Inline Form

```md
## Monitorability Criterion

- Visible Stages:
- Hidden States:
- Safe Local Principle:
```

---

## Canonical Anchors

- [harness-monitorability-review-checklist.md](harness-monitorability-review-checklist.md)
- [harness-monitorability-audit-2026-04-10.md](harness-monitorability-audit-2026-04-10.md)
- [harness-review-gates.md](harness-review-gates.md)
