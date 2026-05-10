# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `60` durable notes are currently indexed across `3` recurring streams.
- `123` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-05-08

- Metadata: `Coverage window: 2026-05-01 to 2026-05-08`
- Anthropic was the clearest first-party mover this week, and the required Anthropic substreams were explicitly checked: [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs overview](https://code.claude.com/docs/en/overview), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). The most important platform change landed on **May 6, 2026**: Claude Platform moved **Multiagent sessions** and **Outcomes** into public beta for Managed Agents, added **Managed Agent webhooks**, added vault credential background refresh for `mcp_oauth`, and improved session/event filtering. This is not a cosmetic docs change; it materially extends Anthropic’s hosted agent runtime surface.
- Anthropic’s Claude Code stream also changed materially inside the same window. The [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows `2.1.132` on **May 6, 2026** and `2.1.133` on **May 7, 2026**. The changes matter operationally: Claude Code now exposes `CLAUDE_CODE_SESSION_ID` to Bash subprocesses, adds `worktree.baseRef` policy for isolation worktrees, passes effort into hooks and Bash via `$CLAUDE_EFFORT`, fixes concurrent-session effort leakage, fixes MCP OAuth proxy handling, and fixes subagents failing to discover project, user, or plugin skills.
- Anthropic’s product/news surface also reinforced the same direction. On **May 5, 2026**, [Agents for financial services](https://www.anthropic.com/news/finance-agents) packaged ten ready-to-run agent templates as Claude Cowork / Claude Code plugins plus Claude Managed Agent cookbooks, and connected them to Microsoft 365 add-ins, connectors, and an MCP app ecosystem. On **May 6, 2026**, [Higher usage limits for Claude and a compute deal with SpaceX](https://www.anthropic.com/news/higher-limits-spacex) doubled Claude Code five-hour limits for paid tiers, removed peak-hour reduction for Pro and Max, and raised API rate limits for Claude Opus models.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-05-10

- Metadata: `Window: 2026-05-03 to 2026-05-10`
- **Anthropic’s strongest in-window technical move is still the Claude Platform managed-agents expansion on 2026-05-06.** The official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) and the matching Claude blog post [New in Claude Managed Agents: dreaming, outcomes, and multiagent orchestration](https://claude.com/blog/new-in-claude-managed-agents-dreaming-outcomes-and-multiagent-orchestration) on **2026-05-06** make this the clearest Anthropic developer signal of the week. The meaningful shift is not branding; it is the surface area: multiagent orchestration, Outcomes, dreaming, webhook support, and better session filtering for managed execution.
- **OpenAI shipped a materially stronger voice-runtime surface on 2026-05-07.** In [Advancing voice intelligence with new models in the API](https://openai.com/index/advancing-voice-intelligence-with-new-models-in-the-api/) on **2026-05-07**, OpenAI introduced GPT-Realtime-2, GPT-Realtime-Translate, and GPT-Realtime-Whisper. The practical signal is that voice is being pushed deeper into the same agent runtime boundary as reasoning, translation, transcription, and tool use.
- **OpenAI also upgraded the fast-default tier on 2026-05-05.** [GPT-5.5 Instant: smarter, clearer, and more personalized](https://openai.com/index/gpt-5-5-instant/) was published on **2026-05-05**. For repeated eval and orchestration workloads, this matters more than another frontier-only model announcement because it moves the price-latency baseline that most real systems actually run.

## Trending Snapshot

### GitHub Trending - 2026-05-10

- Metadata: `Source: [GitHub Trending](https://github.com/trending) as checked on **2026-05-10** from the default global repositories page with no extra filters. The visible board slice still leaned toward agents, coding tools, browser control, voice, and infrastructure rather than a broadly mixed consumer or utility board.`
- [ChromeDevTools/chrome-devtools-mcp](https://github.com/ChromeDevTools/chrome-devtools-mcp) - browser inspection and automation through Chrome DevTools for agent workflows.
- [openai/openai-agents-python](https://github.com/openai/openai-agents-python) - a lightweight framework for multi-agent workflows and tool-using systems.
- [Tracer-Cloud/opensre](https://github.com/Tracer-Cloud/opensre) - AI-SRE and operator tooling for incident and observability-adjacent workflows.

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

- [2026-05-10](../radar/daily/2026-05-10.md)
- [2026-05-09](../radar/daily/2026-05-09.md)
- [2026-05-06](../radar/daily/2026-05-06.md)
- [2026-05-04](../radar/daily/2026-05-04.md)
- [2026-05-02](../radar/daily/2026-05-02.md)
- [2026-05-01](../radar/daily/2026-05-01.md)

### Weekly Synthesis

- [2026-05-08](../radar/weekly/2026-05-08.md)
- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-05-10](../radar/trending/2026-05-10.md)
- [2026-05-09](../radar/trending/2026-05-09.md)
- [2026-05-06](../radar/trending/2026-05-06.md)
- [2026-05-04](../radar/trending/2026-05-04.md)
- [2026-05-02](../radar/trending/2026-05-02.md)
- [2026-05-01](../radar/trending/2026-05-01.md)
