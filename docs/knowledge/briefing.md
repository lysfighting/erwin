# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `90` durable notes are currently indexed across `3` recurring streams.
- `153` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-05-22

- Metadata: `Coverage window: 2026-05-15 to 2026-05-22`
- Anthropic had the clearest in-window platform movement again, and this run explicitly checked all required Anthropic surfaces before drawing that conclusion: [anthropic.com/news](https://www.anthropic.com/news), [claude.com](https://claude.com/), [docs.anthropic.com](https://docs.anthropic.com/en/home), the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), the [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), the [Claude Code docs overview](https://code.claude.com/docs/en/overview), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). This week did have important Anthropic platform and Claude Code changes.
- On **May 19, 2026**, the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added four material managed-runtime upgrades at once: **MCP tunnels** in research preview for private-network servers, **self-hosted sandboxes** for Claude Managed Agents, **mid-session updates** to an agent’s MCP server and tool configuration, and **automatic large-output spillover to sandbox files** when `agent_toolset` or MCP outputs exceed 100K tokens. On **May 18, 2026**, the same surface also upgraded the **web search tool** with richer SEC filing data for grounded financial research.
- Anthropic’s hosted-agent contract is now sharper than last week, not fuzzier. The current [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) explicitly defines **agent**, **environment**, **session**, and **events** as the runtime core, and now includes both Anthropic-managed cloud containers and **self-hosted sandboxes** as first-class environment choices. The managed surface is clearly moving toward enterprise-grade runtime configurability rather than a fixed hosted demo loop.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-05-28

- Metadata: `Window: 2026-05-21 to 2026-05-28`
- **OpenAI expanded the Codex operator surface on 2026-05-21.** The official [ChatGPT Enterprise & Edu release notes](https://help.openai.com/en/articles/10128477) added Appshots in the macOS app, general-availability Goal mode, in-app browser annotations, and locked remote computer use. This is the clearest in-window first-party agent-workflow update across the covered publishers.
- **Anthropic’s required technical surfaces were checked, but no newer in-window dated launch displaced the current Managed Agents state.** The official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code changelog](https://code.claude.com/docs/en/changelog), and [Claude Code docs](https://code.claude.com/docs/en/overview) were all checked this run. The strongest recent Anthropic platform move remains the [2026-05-19 Managed Agents update](https://claude.com/blog/claude-managed-agents-updates), which is now just outside this run window; the advisor tool still requires the `advisor-tool-2026-03-01` beta header and Managed Agents still require `managed-agents-2026-04-01`.
- **No stronger dated first-party technical item surfaced in-window for Google, Moonshot, MiniMax, Qwen, or DeepSeek after one official pass each.** Google’s [Managed Agents in the Gemini API launch](https://blog.google/innovation-and-ai/technology/developers-tools/managed-agents-gemini-api/) is dated **2026-05-19**, Qwen’s [Qwen3.7: The Agent Frontier](https://qwen.ai/research) release is dated **2026-05-20**, and DeepSeek’s latest visible API changelog item remains [DeepSeek-V4](https://api-docs.deepseek.com/updates/) on **2026-04-24**. Those are still relevant context, but not fresh enough to count as new in this 2026-05-21 to 2026-05-28 window.

## Trending Snapshot

### GitHub Trending - 2026-05-28

- Metadata: `Source: default global [GitHub Trending](https://github.com/trending) checked on **2026-05-28** with no extra filters.`
- [harry0703/MoneyPrinterTurbo](https://github.com/harry0703/MoneyPrinterTurbo) - AI-video generation workflow tooling; notable because it is leading the default board rather than a niche AI list.
- [Lum1104/Understand-Anything](https://github.com/Lum1104/Understand-Anything) - explain-anything and study-agent framing; relevant because knowledge-assistance products are still pulling broad attention.
- [hardikpandya/stop-slop](https://github.com/hardikpandya/stop-slop) - quality-control or anti-slop positioning; notable because model-output quality pressure is becoming a recognizable product category.

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

- [2026-05-28](../radar/daily/2026-05-28.md)
- [2026-05-25](../radar/daily/2026-05-25.md)
- [2026-05-24](../radar/daily/2026-05-24.md)
- [2026-05-23](../radar/daily/2026-05-23.md)
- [2026-05-22](../radar/daily/2026-05-22.md)
- [2026-05-21](../radar/daily/2026-05-21.md)

### Weekly Synthesis

- [2026-05-22](../radar/weekly/2026-05-22.md)
- [2026-05-15](../radar/weekly/2026-05-15.md)
- [2026-05-08](../radar/weekly/2026-05-08.md)
- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)

### GitHub Trending

- [2026-05-28](../radar/trending/2026-05-28.md)
- [2026-05-25](../radar/trending/2026-05-25.md)
- [2026-05-24](../radar/trending/2026-05-24.md)
- [2026-05-23](../radar/trending/2026-05-23.md)
- [2026-05-22](../radar/trending/2026-05-22.md)
- [2026-05-21](../radar/trending/2026-05-21.md)
