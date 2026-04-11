# Harness Engineering Vs Internal Gaps

> Date: 2026-04-02
>
> Purpose: compare the strongest recent external `harness engineering` signal with the repo's own gap and drift audits so `erwin` can decide what should become knowledge pressure, eval pressure, or proposal pressure without overstating current `Harness` maturity.

---

## Bottom Line

The strongest recent outside signal and the strongest internal audits point in the same direction.

`Harness` should keep investing first in:

- repo legibility
- control-plane clarity
- observable feedback loops
- evaluation and audit surfaces
- bounded governance maturity

It should not read current productization wins as proof that the deeper governance model is already finished.

The main lesson is not “copy OpenAI's harness.”
The main lesson is:

**external harness-engineering signal and internal gap audits both say the repo should keep tightening legibility, evaluation, and control-plane maturity before widening runtime ambition again.**

---

## External Signal That Matters Most

The current strongest external anchor is the `harness engineering` thread in:

- [../radar/weekly/2026-03-30.md](../radar/weekly/2026-03-30.md)

The durable takeaways from that thread are:

- the repository becomes the system of record for agent work
- agent legibility matters more than giant prompt blobs
- observability and feedback loops are part of the product, not support tooling

This matters because those claims are not just implementation details.
They describe what kind of engineering surface the repo should optimize for.

---

## Internal Reality From Audits

Two internal audit documents say the repo is real, but uneven:

- [../research/harness-gap-audit-2026-03-29.md](../research/harness-gap-audit-2026-03-29.md)
- [../research/harness-research-drift-audit-2026-03-29.md](../research/harness-research-drift-audit-2026-03-29.md)

Their combined view is:

- protocol lifecycle is strong
- adapter/runtime/productization is strong
- leader/teammate runtime is strong
- council semantics are much stronger than before
- but ownership policy is still partial
- harness-wide authority adjudication is still partial
- automated council orchestration is still partial

The critical warning from the drift audit is that the repo can now look more complete than the original control-plane design really is.

That warning is now also captured as normalized feedback in:

- [../feedback/cases/HU-004-productized-surfaces-outpaced-control-plane-maturity.md](../feedback/cases/HU-004-productized-surfaces-outpaced-control-plane-maturity.md)

---

## Where External And Internal Evidence Converge

### 1. Legibility is not cosmetic

External harness-engineering signal says the environment should be readable by the agent.

Internal audits say the repo still needs stronger clarity around:

- who owns what
- which authority wins
- how meetings become accepted truth

So legibility should be treated as a control-plane requirement, not only a docs preference.

### 2. Feedback loops should become stronger before topology grows

External signal emphasizes observable work and evaluable loops.

Internal audits still name gaps in:

- broader authority conflict resolution
- deeper council automation
- richer ownership transfer policy

That means the next maturity step is not “more surfaces.”
It is tighter auditability and more reliable internal correction loops.

### 3. Visible productization can hide real incompleteness

The repo already has impressive concrete surfaces:

- Go runtime
- packaged runtime preference
- shell entrypoints
- manifest validation
- stronger meeting objects

Those are real wins.
But they can create a false sense that the hardest governance work is already done.

This is the main reason the repo needs a durable comparison note like this one.

---

## What Should Not Be Concluded

This comparison does **not** justify:

- declaring `Harness` architecturally finished
- copying vendor orchestration structures directly
- expanding runtime behavior because external signal feels strong
- treating strong packaging or docs surfaces as proof that authority and council automation are solved

The repo should keep the learning-loop discipline:

`radar -> feedback -> knowledge -> proposals -> implementation`

---

## Immediate Implications For erwin

### 1. Strengthen `HLP-003`, do not widen it recklessly

`HLP-003` was already the right posture:

- translate outside harness-engineering signal into eval/doc/automation pressure first
- do not jump straight to runtime topology changes

This note strengthens that proposal's grounding by adding explicit internal-gap evidence, not just first-party signal.

### 2. Keep runtime-facing feedback cases coming

The feedback layer now needs more cases like `HU-004`:

- places where the control plane still lags visible productization
- places where real usage shows governance ambiguity
- places where `Harness` still depends too much on leader interpretation

### 3. Prefer “close the maturity gap” work over “add another surface” work

When future sessions compare plausible next steps, the default bias should be:

- stronger evals
- stronger authority clarity
- stronger council automation
- stronger durable knowledge surfaces

before new runtime expansion.

---

## Suggested Next Questions

1. Which current `Harness` gaps are best encoded first as evals or audits rather than design prose?
2. Which repeated runtime incidents should become the next normalized feedback cases under `docs/feedback/cases/`?
3. Which external harness-engineering claims now have enough internal convergence to move from `watch` to `encode-as-eval` or `prototype`?
