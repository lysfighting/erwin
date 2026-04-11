# Proposal Layer

`docs/proposals/` is the action layer for repo change candidates.

Proposals belong here only after they have passed through stronger upstream shaping:

- `radar` gathered the outside signal
- `feedback` captured real usage pressure
- `knowledge` turned that pressure into durable understanding

Only then should a proposal become a durable control-plane record.

## Current Canonical Proposal Surface

- [harness-proposal-ledger.md](harness-proposal-ledger.md)
- [harness-eval-candidate-ledger.md](harness-eval-candidate-ledger.md)
- [harness-review-gates.md](harness-review-gates.md)
- [harness-benchmark-evidence-checklist.md](harness-benchmark-evidence-checklist.md)
- [harness-monitorability-review-checklist.md](harness-monitorability-review-checklist.md)
- [harness-monitorability-review-criterion.md](harness-monitorability-review-criterion.md)

That ledger remains the reviewed proposal surface for `Harness`.

Worked proposal-review packets can also live here when a proposal needs a deeper, file-backed review instead of only a ledger note.

Use:

- [harness-proposal-review-template.md](harness-proposal-review-template.md)

when creating a new worked review packet.

When a reviewed proposal lands on `encode-as-eval`, the next durable stop should usually be:

- [harness-eval-candidate-ledger.md](harness-eval-candidate-ledger.md)

Every serious proposal should be traceable back to at least one durable upstream anchor:

- an active knowledge thread in `docs/knowledge/active-threads.md`
- a normalized feedback case under `docs/feedback/`
- a durable external or internal source note that has already crossed into `knowledge/`

## Working Rule

Do not treat every daily or weekly run as proposal generation.

Proposal creation should be selective.
The quality bar is:

- repeated signal
- enough evidence density
- clear relevance to `Harness` or repo direction
- a bounded recommended disposition such as `watch`, `encode-as-eval`, `prototype`, or `adopt`
- a visible evidence trail back to `feedback` and/or `knowledge`

## Eval Bridge

- [harness-eval-candidate-ledger.md](harness-eval-candidate-ledger.md) — bridge from reviewed `encode-as-eval` posture into concrete eval, audit, docs-gate, or review-criteria candidates

## Review Gates

- [harness-review-gates.md](harness-review-gates.md) — canonical review gates for proposal traceability, maturity-overclaim checks, benchmark-integrity discipline, and `encode-as-eval` follow-through
- [harness-benchmark-evidence-checklist.md](harness-benchmark-evidence-checklist.md) — reusable review artifact for separating capability signal from eval-trust signal when benchmark-heavy evidence shapes repo posture
- [harness-monitorability-review-checklist.md](harness-monitorability-review-checklist.md) — reusable review artifact for naming visible stages, observability gaps, and safe local principles when monitorability claims shape repo posture
- [harness-monitorability-review-criterion.md](harness-monitorability-review-criterion.md) — lighter three-question review surface for repeated monitorability screening in future eval and proposal work

## Worked Reviews

- [harness-proposal-review-2026-04-11-hlp-001.md](harness-proposal-review-2026-04-11-hlp-001.md) — concrete proposal-review packet confirming that `HLP-001` is genuinely adopted through the repo's live knowledge-first, proposal-gated operating model
- [harness-proposal-review-2026-04-10-hlp-003.md](harness-proposal-review-2026-04-10-hlp-003.md) — first concrete proposal-review packet showing how `HLP-003` is reviewed through thread anchoring, benchmark-integrity, monitorability, and maturity-overclaim checks
- [harness-proposal-review-2026-04-10-hlp-004.md](harness-proposal-review-2026-04-10-hlp-004.md) — concrete proposal-review packet showing how `HLP-004` stays on `prototype` through paper-watch, benchmark-integrity, and selective maintenance logic
- [harness-proposal-review-2026-04-10-hlp-005.md](harness-proposal-review-2026-04-10-hlp-005.md) — concrete proposal-review packet showing how `HLP-005` stays on `watch` through paper-watch, benchmark-integrity, and monitorability review rather than novelty pressure
