# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `25` durable notes are currently indexed across `3` recurring streams.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-04-17

- Metadata: `Coverage window: 2026-04-10 to 2026-04-17`
- The week’s durable field signal is that frontier agent products are converging on **controlled computer/workspace execution** as a platform primitive. OpenAI’s [The next evolution of the Agents SDK](https://openai.com/index/the-next-evolution-of-the-agents-sdk/) published on **April 15, 2026** explicitly introduces a model-native harness, native sandbox execution, workspace manifests, shell/file-edit tools, MCP, skills, `AGENTS.md`, checkpointing, rehydration, and separation between harness state and compute. OpenAI then broadened the Codex product surface in [Codex for (almost) everything](https://openai.com/index/codex-for-almost-everything/) on **April 16, 2026**, positioning Codex as a persistent, computer-operating development partner rather than only a coding assistant.
- Anthropic remained a high-signal first-party stream. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) show [Claude Opus 4.7](https://www.anthropic.com/news/claude-opus-4-7) launching on **April 16, 2026** with stronger long-running coding, agent, vision, instruction-following, file-system memory, and verification behavior. The same release added `xhigh` effort for Opus 4.7, task budgets in public beta on the Claude Platform API, and Claude Code’s `/ultrareview` review command.
- Anthropic’s required substreams were explicitly checked. The Claude Platform release notes, [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs overview](https://code.claude.com/docs/en/overview), and [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md) were reviewed. I found no new Managed Agents launch comparable to the **April 8, 2026** public beta from last week, but Claude Code changed materially this week: the changelog now tops at `2.1.112`, with `2.1.111` adding Opus 4.7 `xhigh`, `/ultrareview`, `/less-permission-prompts`, skill token sorting, improved provider setup, and OpenTelemetry raw API body logging for debugging.

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

The theme this week is **agent execution becoming a governed workspace contract**.

## Current Proposal Candidates

- Workspace manifests for Harness runs
- Effort and budget as explicit run policy
- Review agents as a named verification lane

## Recent History

### Daily Radar

- [2026-04-17](../radar/daily/2026-04-17.md)
- [2026-04-16](../radar/daily/2026-04-16.md)
- [2026-04-15](../radar/daily/2026-04-15.md)
- [2026-04-14](../radar/daily/2026-04-14.md)
- [2026-04-13](../radar/daily/2026-04-13.md)
- [2026-04-10](../radar/daily/2026-04-10.md)

### Weekly Synthesis

- [2026-04-17](../radar/weekly/2026-04-17.md)
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
