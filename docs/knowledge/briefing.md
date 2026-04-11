# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `14` durable notes are currently indexed across `3` recurring streams.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-04-10

- Metadata: `Coverage window: 2026-04-03 to 2026-04-10`
- Anthropic had the strongest first-party in-window platform signal by a wide margin. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) show **Claude Managed Agents** entering public beta on **April 8, 2026**, the new **`ant` CLI** launching the same day, the **advisor tool** entering public beta on **April 9, 2026**, and the Bedrock-compatible `/anthropic/v1/messages` shape entering research preview on **April 7, 2026**. Taken together, that is not a routine docs refresh. It is Anthropic turning agent orchestration, planning quality, and operator ergonomics into a product surface.
- Anthropic's Claude Platform substreams were explicitly checked as required. The new [Advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) describe a planner/executor split in which a cheaper executor can consult a stronger advisor model mid-generation, while the new [Managed Agents production tutorial](https://platform.claude.com/cookbook/managed-agents-cma-operate-in-production) published on **April 8, 2026** makes vault-backed credentials, webhook-driven human-in-the-loop, MCP toolsets, and session lifecycle management look like the default production pattern rather than an advanced custom build.
- Anthropic's Claude Code surface was also explicitly checked. The current [Claude Code docs overview](https://code.claude.com/docs/en/overview) now foregrounds agent teams, the Agent SDK, cloud scheduled tasks, remote control, and cross-surface continuity. The public [Claude Code changelog](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md) currently tops at `2.1.98`, adding a Monitor tool, git-worktree awareness, OTEL `TRACEPARENT` propagation, and several subagent / permission / cleanup fixes. The rendered changelog page does not expose a publish date for `2.1.98`, so the safe reading is: Claude Code continues to broaden materially as an operational agent surface, but we should avoid overstating which specific changelog items landed inside this exact week.

## Daily Signal Snapshot

### AI Agent Daily Radar — 2026-04-10

- Metadata: `Window: 2026-04-03 to 2026-04-10`
- **Anthropic is formalizing trustworthy agents as a full-stack governance problem.** In [Trustworthy agents in practice](https://www.anthropic.com/research/trustworthy-agents) published on **April 9, 2026**, Anthropic says agent safety has to cover the model, harness, tools, and surrounding environment, and points to plan mode, subagents, permission boundaries, benchmarks, monitoring, and open standards such as MCP. The durable takeaway is that agent runtime design is now an explicit governance surface, not just prompt engineering.
- **OpenAI is saying enterprise buyers are moving from copilots to agent operating layers.** In [The next phase of enterprise AI](https://openai.com/index/next-phase-of-enterprise-ai/) published on **April 8, 2026**, OpenAI says enterprise now makes up more than 40% of revenue, Codex has reached **3 million weekly active users**, and customers want governed agents that span company systems and data. The key signal is that agent deployment is being framed as company-wide infrastructure rather than an isolated feature.
- **Anthropic is packaging frontier coding agents as cybersecurity infrastructure.** In [Project Glasswing](https://www.anthropic.com/glasswing) published on **April 7, 2026**, Anthropic launched a joint initiative around Mythos Preview for defensive security and said the model has already found thousands of zero-day vulnerabilities. High-end coding agents are being sold and governed as security operations tooling, not only as developer productivity.

## Trending Snapshot

### GitHub Trending — 2026-04-10

- Metadata: `Source: [GitHub Trending](https://github.com/trending) as captured on 2026-04-10. This run used the default global trending page with no extra filters; the page rendered 10 repositories.`
- [landing-ai/vision-agent](https://github.com/landing-ai/vision-agent) — agent framework for computer vision workflows. Python. 952 stars today.
- [mem0ai/mem0](https://github.com/mem0ai/mem0) — memory layer for AI applications and agents. Python. 816 stars today.
- [enso-labs/enso](https://github.com/enso-labs/enso) — self-improving software engineering agent. Python. 401 stars today.

## Durable Interpretation

The durable theme this week is the productization of the agent control plane.

## Current Proposal Candidates

- Planner / executor split as a first-class Harness primitive
- Managed runtime profiles for long-running agent jobs
- Trace-native subagent operations

## Recent History

### Daily Radar

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

- [2026-04-10](../radar/trending/2026-04-10.md)
- [2026-04-09](../radar/trending/2026-04-09.md)
- [2026-04-08](../radar/trending/2026-04-08.md)
- [2026-04-07](../radar/trending/2026-04-07.md)
- [2026-04-03](../radar/trending/2026-04-03.md)
- [2026-04-02](../radar/trending/2026-04-02.md)
