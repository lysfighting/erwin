# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `71` durable notes are currently indexed across `3` recurring streams.
- `134` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-05-15

- Metadata: `Coverage window: 2026-05-08 to 2026-05-15`
- Anthropic again had the clearest multi-surface movement, and the required Anthropic surfaces were explicitly checked before concluding scope: [anthropic.com/news](https://www.anthropic.com/news), [claude.com](https://claude.com/), [docs.anthropic.com](https://docs.anthropic.com/en/home), the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), the [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), the [Claude Code docs overview](https://code.claude.com/docs/en/overview), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). This week did have important Anthropic platform changes.
- On **May 11, 2026**, Anthropic’s [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added **Claude Platform on AWS**, exposing Messages, Files, Message Batches, Claude Managed Agents, Agent Skills, code execution, and tool use through Anthropic-managed AWS endpoints with AWS billing and IAM authentication. On **May 12, 2026**, the same surface added **fast mode support for Claude Opus 4.7**. These are not minor SDK notes; they extend the provider’s deployment and latency-control surfaces.
- Anthropic’s hosted-agent path remains materially active. The same release-notes surface still shows the **May 6, 2026** launch of **Managed Agent multiagent sessions, outcomes, webhooks, vault refresh, and richer session/event filtering**, and the [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) now clearly presents agent, environment, session, and events as the core runtime contract. The [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) remain live, but I did not find a newer in-window dated advisor launch after the **April 9, 2026** beta; the current doc still centers executor/advisor pairings around Opus 4.7.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-05-15

- Metadata: `Window: 2026-05-08 to 2026-05-15`
- **Anthropic launched Claude Platform on AWS on 2026-05-11.** The official [Claude Platform on AWS post](https://claude.com/blog/claude-platform-on-aws) and the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) position AWS as a first-class access layer for Messages, Files, Message Batches, Claude Managed Agents, Agent Skills, code execution, and the advisor tool. This is the clearest in-window managed-agent distribution signal from any first-party publisher in this run.
- **Anthropic expanded fast-mode controls for high-end agent workloads on 2026-05-12.** The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added fast mode support for `claude-opus-4-7` under the `fast-mode-2026-02-01` beta header on **May 12, 2026**. That is a concrete operator knob for long-horizon runs, not just a model-quality claim.
- **Claude Code shipped a stronger supervisor surface on 2026-05-11 and 2026-05-12.** The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows `claude agents` and the new [agent view](https://code.claude.com/docs/en/agent-view) plus `/goal` on **May 11, 2026**, followed by **May 12, 2026** fixes for background-service stability, `/goal`, managed settings retries, and advisor-tool UX/processing failures. Claude Code is moving further toward a managed operator console instead of a thin coding shell.

## Trending Snapshot

### GitHub Trending - 2026-05-15

- Metadata: `Source: default global [GitHub Trending](https://github.com/trending) checked on **2026-05-15** with no extra filters.`
- [ruvnet/RuView](https://github.com/ruvnet/RuView) - WiFi-based spatial intelligence and vital-sign sensing, notable because the board is still carrying hardware-adjacent sensing projects alongside software tooling.
- [tinyhumansai/openhuman](https://github.com/tinyhumansai/openhuman) - personal AI system framing, showing continued appetite for large-scope open agent projects.
- [rohitg00/agentmemory](https://github.com/rohitg00/agentmemory) - persistent memory layer for coding agents, directly relevant to long-running agent durability.

## Durable Interpretation

This week’s theme is **agent work becoming a managed session fabric rather than a single foreground conversation**.

## Current Proposal Candidates

- Provider-neutral session lifecycle contract
- Environment policy as executable config
- Evented long-running runs with dormancy gates

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-05-15](../radar/daily/2026-05-15.md)
- [2026-05-14](../radar/daily/2026-05-14.md)
- [2026-05-13](../radar/daily/2026-05-13.md)
- [2026-05-12](../radar/daily/2026-05-12.md)
- [2026-05-11](../radar/daily/2026-05-11.md)
- [2026-05-10](../radar/daily/2026-05-10.md)

### Weekly Synthesis

- [2026-05-15](../radar/weekly/2026-05-15.md)
- [2026-05-08](../radar/weekly/2026-05-08.md)
- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-05-15](../radar/trending/2026-05-15.md)
- [2026-05-14](../radar/trending/2026-05-14.md)
- [2026-05-13](../radar/trending/2026-05-13.md)
- [2026-05-12](../radar/trending/2026-05-12.md)
- [2026-05-11](../radar/trending/2026-05-11.md)
- [2026-05-10](../radar/trending/2026-05-10.md)
