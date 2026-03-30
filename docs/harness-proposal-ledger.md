# Harness Proposal Ledger

> Status: Active
>
> Date: 2026-03-31
>
> Purpose: track reviewed `Harness` evolution proposals as durable control-plane records.

This ledger is not a raw inbox and not a copy of weekly radar output.

It records only proposal candidates that have crossed into durable discussion.

| ID | Date | Theme | Origin | Source Ref | Target Layer | Disposition | Status | Weekly Ref | Notes |
|----|------|-------|--------|------------|--------------|-------------|--------|------------|-------|
| HLP-001 | 2026-03-31 | `learning loop` | `internal-gap + field synthesis` | `docs/harness-learning-loop.md` | `automation/docs/eval` | `adopt` | `adopted` | `pending first weekly synthesis with proposal section` | Establish the proposal-first learning loop itself: daily intake stays inbox-only, weekly synthesis stays durable, and no field signal may directly mutate `Harness` without proposal review. |
| HLP-002 | 2026-03-31 | `long-running harness` | `Anthropic engineering` | `https://www.anthropic.com/engineering/harness-design-long-running-apps` | `docs/scheduler/council` | `watch` | `watching` | `pending` | Track the value of sprint contracts, structured artifacts, clean resets, and explicit evaluator roles for long-running app development. Use this as a reference for future council/scheduler refinement, not as a direct topology copy. |
| HLP-003 | 2026-03-31 | `harness engineering` | `OpenAI engineering` | `https://openai.com/index/harness-engineering/` | `eval/docs/automation` | `encode-as-eval` | `queued-for-eval` | `pending` | Evaluate whether more `Harness` changes should be encoded first as legibility/eval improvements instead of direct runtime expansion, especially around repo-readable truth and observable agent feedback loops. |

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
