# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `146` durable notes are currently indexed across `3` recurring streams.
- `209` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI Agent Weekly Synthesis

- Metadata: `Window: 2026-07-04 to 2026-07-10`
- Anthropic's highest-signal user-facing move was in Claude Apps: [Claude Cowork on web and mobile](https://support.claude.com/en/articles/12138966-release-notes) landed on **July 7, 2026**, moving long-running remote sessions out of desktop-only mode and putting persistent agent work closer to an account-level substrate. The same release also expanded the Microsoft 365 connector from search to write actions.
- Anthropic's Claude Platform changes during the window were more administrative than agentic: [API key expiration controls](https://platform.claude.com/docs/en/release-notes/overview) shipped on **July 8, 2026**, and CMEK preservation / access-transparency documentation expanded on **July 10, 2026**. After explicitly checking the Claude Platform release notes, Claude Managed Agents docs surface, advisor tool docs surface, and Claude Code changelog/docs surface, I did **not** find a new in-window managed-agents or advisor-tool capability release after the June 30 to July 2 burst.
- Anthropic's Claude Code surface still mattered this week even without a headline launch: the [Claude Code changelog](https://docs.anthropic.com/en/release-notes/claude-code) shows background-session reliability work, including fixes for inherited `ANTHROPIC_*` provider environment leakage and better `claude agents --json` visibility for blocked and just-dispatched sessions.

## Daily Signal Snapshot

### Daily Radar - 2026-07-11

- Metadata: `Window: 2026-07-04 to 2026-07-11`
- Open the daily note for the full signal list.

## Trending Snapshot

### GitHub Trending - 2026-07-11

- Metadata: `Latest GitHub Trending board summary`
- [catchorg/Catch2](https://github.com/catchorg/Catch2)
- [abseil/abseil-cpp](https://github.com/abseil/abseil-cpp)
- [davila7/claude-code-templates](https://github.com/davila7/claude-code-templates)

## Durable Interpretation

### The frontier is hardening around durable agent operating systems

## Current Proposal Candidates

- No proposal candidates extracted from the latest weekly synthesis yet.

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-07-11](../radar/daily/2026-07-11.md)
- [2026-07-09](../radar/daily/2026-07-09.md)
- [2026-07-08](../radar/daily/2026-07-08.md)
- [2026-07-07](../radar/daily/2026-07-07.md)
- [2026-07-06](../radar/daily/2026-07-06.md)
- [2026-07-05](../radar/daily/2026-07-05.md)

### Weekly Synthesis

- [2026-07-10](../radar/weekly/2026-07-10.md)
- [2026-07-03](../radar/weekly/2026-07-03.md)
- [2026-06-26](../radar/weekly/2026-06-26.md)
- [2026-06-12](../radar/weekly/2026-06-12.md)
- [2026-06-05](../radar/weekly/2026-06-05.md)
- [2026-05-29](../radar/weekly/2026-05-29.md)

### GitHub Trending

- [2026-07-11](../radar/trending/2026-07-11.md)
- [2026-07-09](../radar/trending/2026-07-09.md)
- [2026-07-08](../radar/trending/2026-07-08.md)
- [2026-07-07](../radar/trending/2026-07-07.md)
- [2026-07-06](../radar/trending/2026-07-06.md)
- [2026-07-05](../radar/trending/2026-07-05.md)
