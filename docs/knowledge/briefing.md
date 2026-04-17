# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `24` durable notes are currently indexed across `3` recurring streams.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-04-10

- Metadata: `Coverage window: 2026-04-03 to 2026-04-10`
- Anthropic had the strongest first-party in-window platform signal by a wide margin. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) show **Claude Managed Agents** entering public beta on **April 8, 2026**, the new **`ant` CLI** launching the same day, the **advisor tool** entering public beta on **April 9, 2026**, and the Bedrock-compatible `/anthropic/v1/messages` shape entering research preview on **April 7, 2026**. Taken together, that is not a routine docs refresh. It is Anthropic turning agent orchestration, planning quality, and operator ergonomics into a product surface.
- Anthropic's Claude Platform substreams were explicitly checked as required. The new [Advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) describe a planner/executor split in which a cheaper executor can consult a stronger advisor model mid-generation, while the new [Managed Agents production tutorial](https://platform.claude.com/cookbook/managed-agents-cma-operate-in-production) published on **April 8, 2026** makes vault-backed credentials, webhook-driven human-in-the-loop, MCP toolsets, and session lifecycle management look like the default production pattern rather than an advanced custom build.
- Anthropic's Claude Code surface was also explicitly checked. The current [Claude Code docs overview](https://code.claude.com/docs/en/overview) now foregrounds agent teams, the Agent SDK, cloud scheduled tasks, remote control, and cross-surface continuity. The public [Claude Code changelog](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md) currently tops at `2.1.98`, adding a Monitor tool, git-worktree awareness, OTEL `TRACEPARENT` propagation, and several subagent / permission / cleanup fixes. The rendered changelog page does not expose a publish date for `2.1.98`, so the safe reading is: Claude Code continues to broaden materially as an operational agent surface, but we should avoid overstating which specific changelog items landed inside this exact week.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-04-17

- Metadata: `Window: 2026-04-11 to 2026-04-17`
- **Anthropic shipped Claude Opus 4.7 with stronger long-running agent behavior.** Anthropic published [Introducing Claude Opus 4.7](https://www.anthropic.com/news/claude-opus-4-7) and the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) on **April 16, 2026**. The high-signal update is not just higher model quality: Anthropic explicitly calls out agentic coding, instruction-following literalness, better high-resolution vision for computer-use-style work, and improved file-system memory use. Harness implication: prompts and harnesses tuned for looser prior-model behavior may need stricter regression evals before upgrade.
- **Claude Code 2.1.111/2.1.112 pushed Opus 4.7 into coding-agent operations.** The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) lists **April 16, 2026** releases adding Opus 4.7 `xhigh` effort, Max-subscriber Auto mode with Opus 4.7, `/ultrareview` for cloud multi-agent code review, `/less-permission-prompts`, fewer read-only Bash prompts, plugin diagnostics, OpenTelemetry raw API body logging, and a follow-up fix for Opus 4.7 auto-mode availability. This is a meaningful Claude Code update in-window: model upgrade, review delegation, permission ergonomics, and observability all moved together.
- **OpenAI made the Agents SDK a more explicit sandboxed harness layer.** OpenAI published [The next evolution of the Agents SDK](https://openai.com/index/the-next-evolution-of-the-agents-sdk/) on **April 15, 2026**, describing an updated SDK for agents that inspect files, run commands, edit code, and work across long-horizon tasks in controlled sandbox environments. The first-party direction is clear: standardize the agent loop around workspace manifests, filesystem tools, shell/apply-patch tools, MCP, skills, sandbox-aware orchestration, and durability via snapshot/rehydration.

## Trending Snapshot

### GitHub Trending - 2026-04-17

- Metadata: `Source: [GitHub Trending](https://github.com/trending) as captured on 2026-04-17. This run used the default global trending page with no extra filters; the page rendered 14 repositories.`
- [forrestchang/andrej-karpathy-skills](https://github.com/forrestchang/andrej-karpathy-skills) - a single `CLAUDE.md` file to improve Claude Code behavior, derived from Andrej Karpathy's observations on LLM coding pitfalls. 49,971 stars, 7,959 stars today.
- [thedotmack/claude-mem](https://github.com/thedotmack/claude-mem) - Claude Code plugin that captures coding sessions, compresses them with AI, and injects relevant context into future sessions. TypeScript. 59,780 stars, 1,897 stars today.
- [lsdefine/GenericAgent](https://github.com/lsdefine/GenericAgent) - self-evolving agent that grows a skill tree from a seed and emphasizes reduced token consumption. Python. 2,787 stars, 872 stars today.

## Durable Interpretation

The durable theme this week is the productization of the agent control plane.

## Current Proposal Candidates

- Planner / executor split as a first-class Harness primitive
- Managed runtime profiles for long-running agent jobs
- Trace-native subagent operations

## Recent History

### Daily Radar

- [2026-04-17](../radar/daily/2026-04-17.md)
- [2026-04-16](../radar/daily/2026-04-16.md)
- [2026-04-15](../radar/daily/2026-04-15.md)
- [2026-04-14](../radar/daily/2026-04-14.md)
- [2026-04-13](../radar/daily/2026-04-13.md)
- [2026-04-10](../radar/daily/2026-04-10.md)

### Weekly Synthesis

- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-04-17](../radar/trending/2026-04-17.md)
- [2026-04-16](../radar/trending/2026-04-16.md)
- [2026-04-15](../radar/trending/2026-04-15.md)
- [2026-04-14](../radar/trending/2026-04-14.md)
- [2026-04-13](../radar/trending/2026-04-13.md)
- [2026-04-10](../radar/trending/2026-04-10.md)
