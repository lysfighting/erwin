# Harness Monitorability Audit — 2026-04-10

> Status: Active
>
> Date: 2026-04-10
>
> Purpose: execute the monitorability review checklist against the 2026-04-10 weekly synthesis so `HEL-006` is not only encoded as a rule, but also exercised as a concrete review artifact.

---

## Scope

This audit applies:

- [harness-monitorability-review-checklist.md](harness-monitorability-review-checklist.md)

to:

- [../radar/weekly/2026-04-10.md](../radar/weekly/2026-04-10.md)
- [../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md](../radar/papers/2026-04-10-trajectory-diagnosis-and-observability.md)

The goal is not to copy any outside instrumentation stack.
The goal is to record the first durable repo-level answer to what recent monitorability signal should change in `erwin`.

---

## Monitorability Review Check

- `Source / Runtime Surface`: Anthropic's `Trustworthy agents in practice`, OpenAI's `How we monitor internal coding agents for misalignment`, and the paper cluster on trajectory diagnosis and output-level observability limits. The runtime surfaces under discussion are planning, search, read, edit, validation, approval, and oversight.
- `Visible Stages`: the strongest reusable signal is that search, read, edit, and validation should be treated as separately inspectable stages, and that plan or approval decisions should not be collapsed into generic “agent behavior” if the repo wants diagnosable failure modes.
- `Observability Gap`: final output alone is too lossy because it hides whether the agent found the right context, interpreted it correctly, validated the change, or should have paused for clarification. If the repo reads only final outcomes, it can misdiagnose architectural-reasoning failure as editing failure, or miss that a runtime/resource condition shaped the result more than model behavior did.
- `Safe Local Principle`: the safest local move is to require review and future eval ideas to name which stages were visible, which remained hidden, and which runtime assumptions affected interpretation. This is safer than copying an external monitor stack, full benchmark instrumentation suite, or product-specific oversight workflow.
- `Proposal Or Eval Effect`: this keeps `HLP-003` on `encode-as-eval` and sharpens `HEL-006` as a review-criteria surface before any narrower local audit or eval implementation is attempted. It also keeps `HLP-005` at `watch`, because the repo still needs more than one worked example before monitorability pressure should directly widen runtime topology.
- `Cargo-Cult Guard`: the repo should not recreate Anthropic or OpenAI's internal monitoring systems, and it should not import paper-specific trajectory decomposition as if those categories were automatically the right local runtime abstraction. The principle to keep is diagnosable stage visibility, not external implementation fidelity.

---

## Audit Result

The 2026-04-10 weekly synthesis passes the monitorability review bar because it does all of the following:

- names trajectory stages such as search, read, edit, and validation as the relevant unit of diagnosis
- says explicitly that output-only review is too weak
- distinguishes visible stages from hidden states
- turns the external signal into a bounded local principle rather than a direct adoption claim

The main remaining gap is not conceptual.
It is operational:

- the repo still needs a narrower local audit or review criterion that can be reused beyond this first worked example

---

## Recommended Next Move

Keep using this audit pattern when future weekly or proposal work depends on:

- trajectory diagnosis
- runtime monitorability
- output-only review limits
- oversight claims in realistic coding-agent settings

The next tighter move after that should be one smaller artifact, not a large system:

- [harness-monitorability-review-criterion.md](harness-monitorability-review-criterion.md), a bounded review criterion that asks future `Harness`-facing eval proposals to name visible stages, hidden states, and safe local principles before claiming diagnostic confidence
