# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `64` durable notes are currently indexed across `3` recurring streams.
- `127` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-05-08

- Metadata: `Coverage window: 2026-05-01 to 2026-05-08`
- Anthropic was the clearest first-party mover this week, and the required Anthropic substreams were explicitly checked: [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs overview](https://code.claude.com/docs/en/overview), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). The most important platform change landed on **May 6, 2026**: Claude Platform moved **Multiagent sessions** and **Outcomes** into public beta for Managed Agents, added **Managed Agent webhooks**, added vault credential background refresh for `mcp_oauth`, and improved session/event filtering. This is not a cosmetic docs change; it materially extends Anthropic’s hosted agent runtime surface.
- Anthropic’s Claude Code stream also changed materially inside the same window. The [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows `2.1.132` on **May 6, 2026** and `2.1.133` on **May 7, 2026**. The changes matter operationally: Claude Code now exposes `CLAUDE_CODE_SESSION_ID` to Bash subprocesses, adds `worktree.baseRef` policy for isolation worktrees, passes effort into hooks and Bash via `$CLAUDE_EFFORT`, fixes concurrent-session effort leakage, fixes MCP OAuth proxy handling, and fixes subagents failing to discover project, user, or plugin skills.
- Anthropic’s product/news surface also reinforced the same direction. On **May 5, 2026**, [Agents for financial services](https://www.anthropic.com/news/finance-agents) packaged ten ready-to-run agent templates as Claude Cowork / Claude Code plugins plus Claude Managed Agent cookbooks, and connected them to Microsoft 365 add-ins, connectors, and an MCP app ecosystem. On **May 6, 2026**, [Higher usage limits for Claude and a compute deal with SpaceX](https://www.anthropic.com/news/higher-limits-spacex) doubled Claude Code five-hour limits for paid tiers, removed peak-hour reduction for Pro and Max, and raised API rate limits for Claude Opus models.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-05-12

- Metadata: `Window: 2026-05-05 to 2026-05-12`
- **Anthropic put Claude Platform onto AWS on 2026-05-11.** The official [Claude blog announcement](https://claude.com/blog/claude-platform-on-aws) and the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) both landed on **May 11, 2026**. The concrete signal is not a new model but a new enterprise deployment boundary for the full API and agent stack.
- **Claude Code had a meaningful in-window product update on 2026-05-11.** Anthropic published [Agent view in Claude Code](https://claude.com/blog/agent-view-in-claude-code) on **May 11, 2026**, while the official [Claude Code changelog](https://code.claude.com/docs/en/changelog) also shows fresh entries on **May 7, 2026** and **May 6, 2026**. Anthropic is still actively turning Claude Code into a more inspectable agent-runtime surface rather than treating it as a thin CLI wrapper.
- **Anthropic's managed-agent runtime expansion from 2026-05-06 remains inside the active window.** The official [Claude Managed Agents announcement](https://claude.com/blog/new-in-claude-managed-agents), the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), and the [Managed Agents overview docs](https://platform.claude.com/docs/en/managed-agents/overview) all still point back to **May 6, 2026** as the meaningful official platform move in this lane: public-beta multi-agent sessions, Outcomes, managed-agent webhooks, and better session filtering.

## Trending Snapshot

### GitHub Trending - 2026-05-12

- Metadata: `Source: default global [GitHub Trending](https://github.com/trending) checked on **2026-05-12** with no extra filters. The unauthenticated text-accessible slice available from this environment continued to lean heavily toward agent infrastructure, browser control, observability, and multimodal tooling rather than a broadly mixed board.`
- [ChromeDevTools/chrome-devtools-mcp](https://github.com/ChromeDevTools/chrome-devtools-mcp) - browser inspection, debugging, and automation infrastructure for coding agents.
- [openai/openai-agents-python](https://github.com/openai/openai-agents-python) - a lightweight framework for multi-agent workflows and tool-using systems.
- [Tracer-Cloud/opensre](https://github.com/Tracer-Cloud/opensre) - AI-SRE and observability-adjacent operator tooling.

## Durable Interpretation

This week’s theme is **agents becoming governed operating systems, not just smart tool callers**.

## Current Proposal Candidates

- Managed run webhooks and outcomes
- Effort-aware hooks and wrappers
- Enterprise skill and plugin provenance

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-05-12](../radar/daily/2026-05-12.md)
- [2026-05-11](../radar/daily/2026-05-11.md)
- [2026-05-10](../radar/daily/2026-05-10.md)
- [2026-05-09](../radar/daily/2026-05-09.md)
- [2026-05-06](../radar/daily/2026-05-06.md)
- [2026-05-04](../radar/daily/2026-05-04.md)

### Weekly Synthesis

- [2026-05-08](../radar/weekly/2026-05-08.md)
- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-05-12](../radar/trending/2026-05-12.md)
- [2026-05-11](../radar/trending/2026-05-11.md)
- [2026-05-10](../radar/trending/2026-05-10.md)
- [2026-05-09](../radar/trending/2026-05-09.md)
- [2026-05-06](../radar/trending/2026-05-06.md)
- [2026-05-04](../radar/trending/2026-05-04.md)
