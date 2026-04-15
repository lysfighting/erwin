# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `20` durable notes are currently indexed across `3` recurring streams.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-04-10

- Metadata: `Coverage window: 2026-04-03 to 2026-04-10`
- Anthropic had the strongest first-party in-window platform signal by a wide margin. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) show **Claude Managed Agents** entering public beta on **April 8, 2026**, the new **`ant` CLI** launching the same day, the **advisor tool** entering public beta on **April 9, 2026**, and the Bedrock-compatible `/anthropic/v1/messages` shape entering research preview on **April 7, 2026**. Taken together, that is not a routine docs refresh. It is Anthropic turning agent orchestration, planning quality, and operator ergonomics into a product surface.
- Anthropic's Claude Platform substreams were explicitly checked as required. The new [Advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) describe a planner/executor split in which a cheaper executor can consult a stronger advisor model mid-generation, while the new [Managed Agents production tutorial](https://platform.claude.com/cookbook/managed-agents-cma-operate-in-production) published on **April 8, 2026** makes vault-backed credentials, webhook-driven human-in-the-loop, MCP toolsets, and session lifecycle management look like the default production pattern rather than an advanced custom build.
- Anthropic's Claude Code surface was also explicitly checked. The current [Claude Code docs overview](https://code.claude.com/docs/en/overview) now foregrounds agent teams, the Agent SDK, cloud scheduled tasks, remote control, and cross-surface continuity. The public [Claude Code changelog](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md) currently tops at `2.1.98`, adding a Monitor tool, git-worktree awareness, OTEL `TRACEPARENT` propagation, and several subagent / permission / cleanup fixes. The rendered changelog page does not expose a publish date for `2.1.98`, so the safe reading is: Claude Code continues to broaden materially as an operational agent surface, but we should avoid overstating which specific changelog items landed inside this exact week.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-04-15

- Metadata: `Window: 2026-04-08 to 2026-04-15`
- **Claude Code shipped multiple in-window operations upgrades.** The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) lists version 2.1.108 on **April 14, 2026**, adding 1-hour prompt-cache TTL controls, session recap, slash-command discovery through the Skill tool, better resume defaults, sharper rate-limit errors, and fixes for Agent-tool permission handling and transcript durability. The **April 13, 2026** 2.1.105 release added `EnterWorktree` path switching, PreCompact hooks, plugin background monitors, stale worktree cleanup, MCP truncation guidance, scheduled-task repair, and headless/remote MCP fixes. This is the clearest Claude Code in-window signal: coding-agent reliability is moving into cache policy, resumability, hooks, monitors, worktrees, and remote/headless execution.
- **Anthropic added a platform migration clock for older Claude 4 models.** The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) dated **April 14, 2026** announce deprecation of `claude-sonnet-4-20250514` and `claude-opus-4-20250514`, with API retirement scheduled for **June 15, 2026**, and recommend migration to Claude Sonnet 4.6 and Claude Opus 4.6. For harness work, this is less a feature launch than a dependency hygiene signal: long-running agent systems need explicit model-lifecycle checks, migration tests, and stale-model detection.
- **Anthropic's Managed Agents public beta remains the strongest managed-runtime signal in-window.** The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) date the launch to **April 8, 2026**, and the [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) describes a pre-built configurable agent harness for long-running and asynchronous work, with agents, environments, sessions, events, built-in tools, MCP, skills, server-sent event streaming, and beta access via `managed-agents-2026-04-01`. The paired Anthropic engineering post, [Scaling Managed Agents: Decoupling the brain from the hands](https://www.anthropic.com/engineering/managed-agents), makes the architecture explicit: session log, harness, and sandbox are separated so each can fail or evolve independently.

## Trending Snapshot

### GitHub Trending - 2026-04-15

- Metadata: `Source: [GitHub Trending](https://github.com/trending?since=daily) as captured on 2026-04-15. This run used the default global trending page with no extra filters; the page rendered 12 repositories.`
- [forrestchang/andrej-karpathy-skills](https://github.com/forrestchang/andrej-karpathy-skills) - a single `CLAUDE.md` behavior file distilling Andrej Karpathy's public LLM advice for Claude Code. 6,261 stars, 1,818 stars today.
- [thedotmack/claude-mem](https://github.com/thedotmack/claude-mem) - Claude Code plugin for capturing, compressing, and reinjecting session context. TypeScript. 3,431 stars, 784 stars today.
- [jamiepine/voicebox](https://github.com/jamiepine/voicebox) - open-source voice synthesis studio with local, Docker, and desktop usage modes. TypeScript. 1,558 stars, 435 stars today.

## Durable Interpretation

The durable theme this week is the productization of the agent control plane.

## Current Proposal Candidates

- Planner / executor split as a first-class Harness primitive
- Managed runtime profiles for long-running agent jobs
- Trace-native subagent operations

## Recent History

### Daily Radar

- [2026-04-15](../radar/daily/2026-04-15.md)
- [2026-04-14](../radar/daily/2026-04-14.md)
- [2026-04-13](../radar/daily/2026-04-13.md)
- [2026-04-10](../radar/daily/2026-04-10.md)
- [2026-04-09](../radar/daily/2026-04-09.md)
- [2026-04-08](../radar/daily/2026-04-08.md)

### Weekly Synthesis

- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-04-15](../radar/trending/2026-04-15.md)
- [2026-04-14](../radar/trending/2026-04-14.md)
- [2026-04-13](../radar/trending/2026-04-13.md)
- [2026-04-10](../radar/trending/2026-04-10.md)
- [2026-04-09](../radar/trending/2026-04-09.md)
- [2026-04-08](../radar/trending/2026-04-08.md)
