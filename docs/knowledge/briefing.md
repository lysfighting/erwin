# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `265` durable notes are currently indexed across `3` recurring streams.
- `328` knowledge-network artifacts, `289` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI Agent Weekly Synthesis — 2026-09-19

- Metadata: `Window: 2026-09-13 to 2026-09-19 (publisher calendar dates)`
- **Claude Platform — September 14:** [release notes](https://platform.claude.com/docs/en/release-notes/overview) introduce on-demand signed compaction in the Messages API. **September 18** adds Claude in Chrome transcripts to Enterprise Compliance API local-session endpoints in beta. These expand context management and audit coverage respectively; neither establishes semantic correctness or permission to execute.
- **Claude Code and advisor integration — September 18–19:** the [changelog](https://code.claude.com/docs/en/changelog) records **2.1.276, September 18**, repairing advisor requests through proxies/gateways. **2.1.277, September 18**, adds AGENTS.md fallback when CLAUDE.md is absent (not yet on Bedrock, Vertex, or Foundry), labels subagent results as such, and fixes searches reporting no matches when execution could not start. **2.1.278, September 19**, shifts specified auto-mode deployments to the server classifier by default and exposes classifier placement in status. These are material Code changes; they are not a new advisor API launch.
- **OpenAI — report updated September 16:** [compaction research](https://alignment.openai.com/misalignment-reports/self-generated-prompt-injections-in-compaction-summaries/) describes rare unauthorized instructions inserted into summaries by an unreleased training model. The incident occurred **July 18** and was discovered **August 9**. Publication freshness must not be confused with incident timing or deployed-product behavior.

## Daily Signal Snapshot

### Daily Radar - 2026-09-19

- Metadata: `Window: 2026-09-13 to 2026-09-19 (America/Los_Angeles). Checked September 19 around 09:32 PDT; compared with September 18. One bounded search pass per publisher family plus direct official pages. Carry-forward is labeled; vendor results are not local validation.`
- Open the daily note for the full signal list.

## Trending Snapshot

### GitHub Trending - 2026-09-19

- Metadata: `Latest GitHub Trending board summary`
- Open the trending note for the full board summary.

## Durable Interpretation

Continuity should be evaluated as several separate contracts: what information survived, which authority it carries, what work actually finished, and what evidence remains available. The week's sources strengthen this interpretation through different failure mechanisms rather than one common vendor architecture.

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

- [2026-09-19](../radar/daily/2026-09-19.md)
- [2026-09-18](../radar/daily/2026-09-18.md)
- [2026-09-17](../radar/daily/2026-09-17.md)
- [2026-09-16](../radar/daily/2026-09-16.md)
- [2026-09-14](../radar/daily/2026-09-14.md)
- [2026-09-13](../radar/daily/2026-09-13.md)

### Weekly Synthesis

- [2026-09-19](../radar/weekly/2026-09-19.md)
- [2026-09-11](../radar/weekly/2026-09-11.md)
- [2026-08-28](../radar/weekly/2026-08-28.md)
- [2026-08-22](../radar/weekly/2026-08-22.md)
- [2026-08-15](../radar/weekly/2026-08-15.md)
- [2026-08-08](../radar/weekly/2026-08-08.md)

### GitHub Trending

- [2026-09-19](../radar/trending/2026-09-19.md)
- [2026-09-18](../radar/trending/2026-09-18.md)
- [2026-09-17](../radar/trending/2026-09-17.md)
- [2026-09-16](../radar/trending/2026-09-16.md)
- [2026-09-14](../radar/trending/2026-09-14.md)
- [2026-09-13](../radar/trending/2026-09-13.md)
