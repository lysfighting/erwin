# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `172` durable notes are currently indexed across `3` recurring streams.
- `235` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI Agent Weekly Synthesis

- Metadata: `Window: 2026-07-18 to 2026-07-25`
- **Anthropic had the clearest in-window platform movement.** The official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added **Claude Opus 5** on **July 24, 2026**, with a 1M-token context window, 128k max output tokens, and thinking on by default. The same release added **mid-conversation tool changes** in beta, which matters more than the model launch itself because it lets long-running sessions add or remove tools without throwing away prompt-cache continuity.
- **Claude Code moved in lockstep with the platform.** The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows that on **July 24, 2026** Claude Code made **Opus 5 the default Opus model**, added `sandbox.network.strictAllowlist`, and added a `DirectoryAdded` hook for newly registered repo roots mid-session. The [Claude Code “What’s New” surface](https://code.claude.com/docs/en/whats-new) still shows **Week 29 (July 13–17, 2026)** as the latest digest, so the changelog is the authoritative surface for this week’s Claude Code changes.
- **Anthropic’s Managed Agents and advisor tool were explicitly checked and did not receive a new in-window launch.** I checked the required [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Code changelog](https://code.claude.com/docs/en/changelog), and [Claude Code docs surface](https://code.claude.com/docs/en/whats-new). The latest material Managed Agents change still visible on the official release notes surface is **June 9, 2026** scheduled deployments and vault env-var credentials; the latest advisor-tool change remains **June 2, 2026** `max_tokens`.

## Daily Signal Snapshot

### Daily Radar - 2026-07-26

- Metadata: `Window: 2026-07-20 to 2026-07-26`
- Open the daily note for the full signal list.

## Trending Snapshot

### GitHub Trending - 2026-07-26

- Metadata: `Latest GitHub Trending board summary`
- [permissionlesstech/bitchat](https://github.com/permissionlesstech/bitchat)
- [citrolabs/ego-lite](https://github.com/citrolabs/ego-lite)
- [block/buzz](https://github.com/block/buzz)

## Durable Interpretation

### Stateful agents are becoming governable runtime objects

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

- [2026-07-26](../radar/daily/2026-07-26.md)
- [2026-07-25](../radar/daily/2026-07-25.md)
- [2026-07-23](../radar/daily/2026-07-23.md)
- [2026-07-22](../radar/daily/2026-07-22.md)
- [2026-07-21](../radar/daily/2026-07-21.md)
- [2026-07-20](../radar/daily/2026-07-20.md)

### Weekly Synthesis

- [2026-07-25](../radar/weekly/2026-07-25.md)
- [2026-07-18](../radar/weekly/2026-07-18.md)
- [2026-07-10](../radar/weekly/2026-07-10.md)
- [2026-07-03](../radar/weekly/2026-07-03.md)
- [2026-06-26](../radar/weekly/2026-06-26.md)
- [2026-06-12](../radar/weekly/2026-06-12.md)

### GitHub Trending

- [2026-07-26](../radar/trending/2026-07-26.md)
- [2026-07-25](../radar/trending/2026-07-25.md)
- [2026-07-23](../radar/trending/2026-07-23.md)
- [2026-07-22](../radar/trending/2026-07-22.md)
- [2026-07-21](../radar/trending/2026-07-21.md)
- [2026-07-20](../radar/trending/2026-07-20.md)
