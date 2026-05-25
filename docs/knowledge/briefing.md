# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `88` durable notes are currently indexed across `3` recurring streams.
- `151` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-05-22

- Metadata: `Coverage window: 2026-05-15 to 2026-05-22`
- Anthropic had the clearest in-window platform movement again, and this run explicitly checked all required Anthropic surfaces before drawing that conclusion: [anthropic.com/news](https://www.anthropic.com/news), [claude.com](https://claude.com/), [docs.anthropic.com](https://docs.anthropic.com/en/home), the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), the [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), the [Claude Code docs overview](https://code.claude.com/docs/en/overview), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). This week did have important Anthropic platform and Claude Code changes.
- On **May 19, 2026**, the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added four material managed-runtime upgrades at once: **MCP tunnels** in research preview for private-network servers, **self-hosted sandboxes** for Claude Managed Agents, **mid-session updates** to an agent’s MCP server and tool configuration, and **automatic large-output spillover to sandbox files** when `agent_toolset` or MCP outputs exceed 100K tokens. On **May 18, 2026**, the same surface also upgraded the **web search tool** with richer SEC filing data for grounded financial research.
- Anthropic’s hosted-agent contract is now sharper than last week, not fuzzier. The current [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) explicitly defines **agent**, **environment**, **session**, and **events** as the runtime core, and now includes both Anthropic-managed cloud containers and **self-hosted sandboxes** as first-class environment choices. The managed surface is clearly moving toward enterprise-grade runtime configurability rather than a fixed hosted demo loop.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-05-25

- Metadata: `Window: 2026-05-18 to 2026-05-25`
- **Anthropic expanded Claude Managed Agents on 2026-05-19.** The official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) and matching [Claude blog post](https://claude.com/blog/new-in-claude-managed-agents-self-hosted-sandboxes-and-mcp-tunnels) added MCP tunnels in research preview, self-hosted sandboxes, live MCP and tool reconfiguration for active sessions, and automatic spill-to-file handling for outputs above 100K tokens.
- **OpenAI sharpened the Codex operator surface on 2026-05-21.** The official [ChatGPT release notes](https://help.openai.com/en/articles/6825453-chatgpt-release-notes%252525252525252525252525252525252525252525252525252525252525252525252525252523.doc) added Appshots, general-availability Goal mode, in-app browser annotations, locked computer use, and broader browser-use improvements for longer-running coding workflows.
- **Google launched Managed Agents in the Gemini API on 2026-05-19.** In [Introducing Managed Agents in the Gemini API](https://blog.google/innovation-and-ai/technology/developers-tools/managed-agents-gemini-api/), Google exposed a single-call agent runtime with isolated Linux execution, resumable state, web access, and repo-local agent contracts via `AGENTS.md` and `SKILL.md`.

## Trending Snapshot

### GitHub Trending - 2026-05-25

- Metadata: `Source: default global [GitHub Trending](https://github.com/trending) checked on **2026-05-25** with no extra filters.`
- [anthropics/knowledge-work-plugins](https://github.com/anthropics/knowledge-work-plugins) - notable because Anthropic-adjacent plugin packaging is surfacing directly on the default board rather than only on niche agent lists.
- [multica-ai/andrej-karpathy-skills](https://github.com/multica-ai/andrej-karpathy-skills) - a portable behavior overlay; relevant because installable agent operating packs continue to attract broad attention.
- [dotnet/skills](https://github.com/dotnet/skills) - first-party .NET skill repository for AI coding agents; notable because a major platform vendor is treating skills as a product surface.

## Durable Interpretation

This week’s theme is **governed agent runtime portability**.

## Current Proposal Candidates

- Portable execution-environment contract
- Pinned run classes and dormancy policy
- Review-to-repair handoff packets

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-05-25](../radar/daily/2026-05-25.md)
- [2026-05-24](../radar/daily/2026-05-24.md)
- [2026-05-23](../radar/daily/2026-05-23.md)
- [2026-05-22](../radar/daily/2026-05-22.md)
- [2026-05-21](../radar/daily/2026-05-21.md)
- [2026-05-20](../radar/daily/2026-05-20.md)

### Weekly Synthesis

- [2026-05-22](../radar/weekly/2026-05-22.md)
- [2026-05-15](../radar/weekly/2026-05-15.md)
- [2026-05-08](../radar/weekly/2026-05-08.md)
- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)

### GitHub Trending

- [2026-05-25](../radar/trending/2026-05-25.md)
- [2026-05-24](../radar/trending/2026-05-24.md)
- [2026-05-23](../radar/trending/2026-05-23.md)
- [2026-05-22](../radar/trending/2026-05-22.md)
- [2026-05-21](../radar/trending/2026-05-21.md)
- [2026-05-20](../radar/trending/2026-05-20.md)
