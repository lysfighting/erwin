# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `16` durable notes are currently indexed across `3` recurring streams.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-04-10

- Metadata: `Coverage window: 2026-04-03 to 2026-04-10`
- Anthropic had the strongest first-party in-window platform signal by a wide margin. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) show **Claude Managed Agents** entering public beta on **April 8, 2026**, the new **`ant` CLI** launching the same day, the **advisor tool** entering public beta on **April 9, 2026**, and the Bedrock-compatible `/anthropic/v1/messages` shape entering research preview on **April 7, 2026**. Taken together, that is not a routine docs refresh. It is Anthropic turning agent orchestration, planning quality, and operator ergonomics into a product surface.
- Anthropic's Claude Platform substreams were explicitly checked as required. The new [Advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) describe a planner/executor split in which a cheaper executor can consult a stronger advisor model mid-generation, while the new [Managed Agents production tutorial](https://platform.claude.com/cookbook/managed-agents-cma-operate-in-production) published on **April 8, 2026** makes vault-backed credentials, webhook-driven human-in-the-loop, MCP toolsets, and session lifecycle management look like the default production pattern rather than an advanced custom build.
- Anthropic's Claude Code surface was also explicitly checked. The current [Claude Code docs overview](https://code.claude.com/docs/en/overview) now foregrounds agent teams, the Agent SDK, cloud scheduled tasks, remote control, and cross-surface continuity. The public [Claude Code changelog](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md) currently tops at `2.1.98`, adding a Monitor tool, git-worktree awareness, OTEL `TRACEPARENT` propagation, and several subagent / permission / cleanup fixes. The rendered changelog page does not expose a publish date for `2.1.98`, so the safe reading is: Claude Code continues to broaden materially as an operational agent surface, but we should avoid overstating which specific changelog items landed inside this exact week.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-04-13

- Metadata: `Window: 2026-04-06 to 2026-04-13`
- **Anthropic launched Claude Managed Agents as a hosted agent harness.** The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) list **April 8, 2026** as the public beta launch date, and the [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) defines agents, environments, sessions, and events as first-class API concepts. The important shift is that Anthropic is no longer only documenting tool calls; it is offering managed containers, SSE event streams, persisted sessions, built-in tools, MCP, skills, sandboxing, compaction, and prompt caching as a full autonomous-agent runtime surface.
- **Anthropic launched the advisor tool for executor/advisor model splits.** The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) date the public beta to **April 9, 2026**, and the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) frame it as a way for faster executor models to consult a higher-intelligence advisor model mid-generation. The durable pattern is explicit runtime decomposition: most tokens come from a cheaper executor, while strategic intervention is isolated, metered, cached, and visible in usage iterations.
- **Anthropic added `ant`, a CLI that makes Claude API resources scriptable.** The [April 8, 2026](https://platform.claude.com/docs/en/release-notes/overview) release note pairs the CLI with Managed Agents, and the [CLI docs](https://platform.claude.com/docs/en/api/sdks/cli) show beta agent/session/environment resources exposed as command namespaces with YAML/JSON input, transforms, pagination, and native Claude Code integration. This matters because agent infrastructure is becoming inspectable and automatable from the terminal, not only from SDK code.

## Trending Snapshot

### GitHub Trending - 2026-04-13

- Metadata: `Source: [GitHub Trending](https://github.com/trending) as captured on 2026-04-13. This run used the default global trending page with no extra filters; the page rendered 14 repositories.`
- [NousResearch/hermes-agent](https://github.com/NousResearch/hermes-agent) - adaptive AI agent project. Python. 7,454 stars today.
- [microsoft/markitdown](https://github.com/microsoft/markitdown) - converts files and office documents to Markdown. Python. 2,513 stars today.
- [forrestchang/andrej-karpathy-skills](https://github.com/forrestchang/andrej-karpathy-skills) - CLAUDE.md guidance distilled from LLM coding pitfalls. 2,369 stars today.

## Durable Interpretation

The durable theme this week is the productization of the agent control plane.

## Current Proposal Candidates

- Planner / executor split as a first-class Harness primitive
- Managed runtime profiles for long-running agent jobs
- Trace-native subagent operations

## Recent History

### Daily Radar

- [2026-04-13](../radar/daily/2026-04-13.md)
- [2026-04-10](../radar/daily/2026-04-10.md)
- [2026-04-09](../radar/daily/2026-04-09.md)
- [2026-04-08](../radar/daily/2026-04-08.md)
- [2026-04-03](../radar/daily/2026-04-03.md)
- [2026-04-02](../radar/daily/2026-04-02.md)

### Weekly Synthesis

- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-04-13](../radar/trending/2026-04-13.md)
- [2026-04-10](../radar/trending/2026-04-10.md)
- [2026-04-09](../radar/trending/2026-04-09.md)
- [2026-04-08](../radar/trending/2026-04-08.md)
- [2026-04-07](../radar/trending/2026-04-07.md)
- [2026-04-03](../radar/trending/2026-04-03.md)
