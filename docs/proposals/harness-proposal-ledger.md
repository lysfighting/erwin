# Harness Proposal Ledger

> Status: Active
>
> Date: 2026-04-02
>
> Purpose: track reviewed `Harness` evolution proposals as durable control-plane records.

This ledger is not a raw inbox and not a copy of weekly radar output.

It records only proposal candidates that have crossed into durable discussion.

| ID | Date | Theme | Origin | Source Ref | Target Layer | Disposition | Status | Weekly Ref | Notes |
|----|------|-------|--------|------------|--------------|-------------|--------|------------|-------|
| HLP-001 | 2026-03-31 | `learning loop` | `internal-gap + field synthesis` | `docs/harness-learning-loop.md` | `automation/docs/eval` | `adopt` | `adopted` | `fulfilled by post-refactor weekly, knowledge, and proposal surfaces` | Establish the learning-loop discipline itself: no field signal may directly mutate `Harness` without review, and the repo should evolve through durable intake, knowledge, and proposal surfaces instead of chat memory. `docs/proposals/harness-proposal-review-2026-04-11-hlp-001.md` now confirms that this proposal is not only accepted in principle but already operating through the repo's active knowledge system. |
| HLP-002 | 2026-03-31 | `long-running harness` | `Anthropic engineering` | `https://www.anthropic.com/engineering/harness-design-long-running-apps` | `docs/scheduler/council` | `watch` | `watching` | `pending` | Track the value of sprint contracts, structured artifacts, clean resets, and explicit evaluator roles for long-running app development. Use this as a reference for future council/scheduler refinement, not as a direct topology copy. |
| HLP-003 | 2026-03-31 | `harness engineering` | `OpenAI engineering + internal gap audits` | `docs/knowledge/harness-engineering-vs-internal-gaps-2026-04-02.md` | `eval/docs/automation` | `encode-as-eval` | `queued-for-eval` | `docs/radar/weekly/2026-04-10.md` | Evaluate whether more `Harness` changes should be encoded first as legibility/eval improvements instead of direct runtime expansion, especially around repo-readable truth, observable agent feedback loops, and the remaining mismatch between visible productization and deeper governance maturity. This proposal is now also reinforced by the runtime legibility patterns recorded in `docs/knowledge/runtime-legibility-patterns-2026-04-02.md`, plus the newer convergence around monitorability, low-latency oversight, infrastructure-noise skepticism, and eval-method honesty captured in the 2026-04-10 weekly synthesis. `docs/proposals/harness-monitorability-audit-2026-04-10.md` is now the first worked review artifact on that monitorability pressure, and `docs/proposals/harness-proposal-review-2026-04-10-hlp-003.md` is now the first concrete proposal-review packet for this proposal. |
| HLP-004 | 2026-04-02 | `paper-aware radar` | `internal-gap + research scan` | `docs/knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md` | `automation/docs` | `prototype` | `queued-for-prototype` | `docs/radar/weekly/2026-04-10.md` | Add a dedicated paper-intake lane to the radar stack so self-improvement, coding-agent eval, and benchmark knowledge does not depend only on first-party blogs and docs. Keep this proposal knowledge-first and proposal-gated: better intake first, no direct Harness runtime mutation. The 2026-04-10 weekly synthesis strengthens this proposal again by showing the archive can now carry trajectory diagnosis and observability work rather than only benchmark notes. `docs/proposals/harness-proposal-review-2026-04-10-hlp-004.md` is now the first concrete proposal-review packet for this prototype posture. |
| HLP-005 | 2026-04-02 | `coding-agent self-iteration watch` | `OpenAI + Anthropic + papers` | `docs/knowledge/coding-agent-self-iteration-and-paper-intake-2026-04-02.md` | `eval/automation/docs` | `watch` | `watching` | `docs/radar/weekly/2026-04-10.md` | Track the split between OpenAI’s eval-driven self-improvement posture and Anthropic’s dogfooding-plus-automation posture. Use this to shape future Harness eval design and paper watch criteria rather than copying either workflow mechanically. The 2026-04-10 weekly synthesis now sharpens this watch around trajectory observability, stage-level diagnosis, and monitorability pressure, while `docs/knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md` continues to supply the benchmark-trust discipline for interpreting those signals. `docs/proposals/harness-proposal-review-2026-04-10-hlp-005.md` is now the first concrete proposal-review packet for this watch posture. |

## Status Meanings

- `watching`
  - keep observing; do not schedule implementation work yet
- `queued-for-eval`
  - worth encoding as smoke, audit, or review criteria first
- `queued-for-prototype`
  - worth a bounded design or implementation spike
- `adopted`
  - accepted as part of the current system or process
- `rejected`
  - intentionally not moving forward
- `superseded`
  - replaced by a better or newer proposal
