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
| HLP-004 | 2026-04-02 | `paper-aware radar` | `internal-gap + research scan` | `docs/research/coding-agent-self-iteration-and-paper-intake-2026-04-02.md` | `automation/docs` | `prototype` | `queued-for-prototype` | `pending weekly synthesis after automation update` | Add a dedicated paper-intake lane to the radar stack so self-improvement, coding-agent eval, and benchmark knowledge does not depend only on first-party blogs and docs. Keep this proposal proposal-first: better intake first, no direct Harness runtime mutation. |
| HLP-005 | 2026-04-02 | `coding-agent self-iteration watch` | `OpenAI + Anthropic + papers` | `docs/research/coding-agent-self-iteration-and-paper-intake-2026-04-02.md` | `eval/automation/docs` | `watch` | `watching` | `pending` | Track the split between OpenAI’s eval-driven self-improvement posture and Anthropic’s dogfooding-plus-automation posture. Use this to shape future Harness eval design and paper watch criteria rather than copying either workflow mechanically. |

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
