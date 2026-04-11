# Harness Review Gates

> Status: Active
>
> Date: 2026-04-02
>
> Purpose: define the minimum review gates for `Harness`-facing proposals and adjacent knowledge claims so the repo keeps legibility, evidence density, and governance honesty ahead of runtime expansion.

---

## Why This Exists

The repo now has:

- active knowledge threads
- normalized feedback cases
- a proposal ledger
- an eval-candidate bridge

What it still needs is a compact set of review gates that stop weak proposal pressure or maturity overclaim from slipping through just because the surrounding docs look polished.

---

## Review Gates

### Gate 1: Active Thread Anchor

A serious proposal or harness-facing change claim should name:

- the active knowledge thread it updates
- the durable upstream anchors that support it

At least one anchor should come from `feedback` or `knowledge` when the change claim is repo-internal rather than purely external observation.

This encodes the proposal-review half of `HEL-001`.

### Gate 2: No Maturity Overclaim

No proposal, comparison note, or harness-facing review should imply that governance maturity is complete if the current audits still mark key areas as partial.

For now, that specifically includes:

- ownership policy
- broader authority adjudication
- automated council orchestration

When a document discusses those areas, it should stay aligned with:

- [../research/harness-gap-audit-2026-03-29.md](../research/harness-gap-audit-2026-03-29.md)
- [../research/harness-research-drift-audit-2026-03-29.md](../research/harness-research-drift-audit-2026-03-29.md)
- [../feedback/cases/HU-004-productized-surfaces-outpaced-control-plane-maturity.md](../feedback/cases/HU-004-productized-surfaces-outpaced-control-plane-maturity.md)

This encodes `HEL-003`.

### Gate 3: `encode-as-eval` Must Land Somewhere

If a reviewed proposal recommends `encode-as-eval`, it should usually create or update a concrete entry in:

- [harness-eval-candidate-ledger.md](harness-eval-candidate-ledger.md)

That entry should name:

- what will be checked
- where it should live
- why that check is the safest next move

Without this step, `encode-as-eval` is only a polite way to defer work indefinitely.

### Gate 4: Benchmark-Heavy Evidence Must Name Trust Limits

If a review relies on benchmark-heavy external evidence, it should explicitly separate:

- `capability signal`
  what the benchmark result suggests the agent can or cannot do
- `eval-trust signal`
  how much confidence we should place in that benchmark as a durable measure rather than a gameable or contamination-sensitive target

This is especially important when the evidence comes from:

- long-horizon coding benchmarks
- context-retrieval benchmarks
- browse or web-agent benchmarks

When this gate applies, review should stay aligned with:

- [../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md](../knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md)
- [../feedback/cases/HU-007-benchmark-results-needed-eval-trust-interpretation.md](../feedback/cases/HU-007-benchmark-results-needed-eval-trust-interpretation.md)
- [harness-benchmark-evidence-checklist.md](harness-benchmark-evidence-checklist.md)

This encodes `HEL-005`.

### Gate 5: Monitorability Claims Must Name What Stays Visible

If a review relies on trajectory diagnosis, monitorability, intermediate observability, or output-only-feedback limitations, it should explicitly say:

- which stages or decisions remain visible enough to inspect
- which important states remain hidden or too lossy for reliable diagnosis from final output alone
- which smaller local review or eval principle should be encoded instead of copying external instrumentation directly

This is especially important when the evidence comes from:

- trajectory-level coding-agent papers
- internal or external monitor writeups
- stage-level search/read/edit diagnosis
- product or research claims about oversight, visibility, or interpretability in realistic agent runtimes

When this gate applies, review should stay aligned with:

- [../radar/weekly/2026-04-10.md](../radar/weekly/2026-04-10.md)
- [../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md](../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md)
- [harness-monitorability-review-checklist.md](harness-monitorability-review-checklist.md)

This encodes `HEL-006`.

---

## Current Canonical Uses

Use these gates when reviewing:

- proposal-ledger additions
- weekly proposal sections
- knowledge notes that sharpen `Harness` evolution pressure
- harness-facing documentation that makes maturity claims

When Gate 4 applies, prefer using:

- [harness-benchmark-evidence-checklist.md](harness-benchmark-evidence-checklist.md)

When Gate 5 applies, prefer using:

- [harness-monitorability-review-checklist.md](harness-monitorability-review-checklist.md)
- [harness-monitorability-review-criterion.md](harness-monitorability-review-criterion.md) when a lighter reusable screen is enough
