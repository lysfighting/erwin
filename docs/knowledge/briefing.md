# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `47` durable notes are currently indexed across `3` recurring streams.
- `110` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-04-17

- Metadata: `Coverage window: 2026-04-10 to 2026-04-17`
- The week’s durable field signal is that frontier agent products are converging on **controlled computer/workspace execution** as a platform primitive. OpenAI’s [The next evolution of the Agents SDK](https://openai.com/index/the-next-evolution-of-the-agents-sdk/) published on **April 15, 2026** explicitly introduces a model-native harness, native sandbox execution, workspace manifests, shell/file-edit tools, MCP, skills, `AGENTS.md`, checkpointing, rehydration, and separation between harness state and compute. OpenAI then broadened the Codex product surface in [Codex for (almost) everything](https://openai.com/index/codex-for-almost-everything/) on **April 16, 2026**, positioning Codex as a persistent, computer-operating development partner rather than only a coding assistant.
- Anthropic remained a high-signal first-party stream. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) show [Claude Opus 4.7](https://www.anthropic.com/news/claude-opus-4-7) launching on **April 16, 2026** with stronger long-running coding, agent, vision, instruction-following, file-system memory, and verification behavior. The same release added `xhigh` effort for Opus 4.7, task budgets in public beta on the Claude Platform API, and Claude Code’s `/ultrareview` review command.
- Anthropic’s required substreams were explicitly checked. The Claude Platform release notes, [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs overview](https://code.claude.com/docs/en/overview), and [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md) were reviewed. I found no new Managed Agents launch comparable to the **April 8, 2026** public beta from last week, but Claude Code changed materially this week: the changelog now tops at `2.1.112`, with `2.1.111` adding Opus 4.7 `xhigh`, `/ultrareview`, `/less-permission-prompts`, skill token sorting, improved provider setup, and OpenTelemetry raw API body logging for debugging.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-04-30

- Metadata: `Window: 2026-04-23 to 2026-04-30`
- **Anthropic's official technical surface moved Claude Platform forward, while Claude Code stayed quiet in-window.** The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) published **April 23, 2026** that **memory for Claude Managed Agents** is now in public beta, and published **April 24, 2026** that the **Rate Limits API** is now available. Required Anthropic surfaces were explicitly checked on this run: [anthropic.com/news](https://www.anthropic.com/news), [claude.com](https://claude.com/), [docs.anthropic.com](https://docs.anthropic.com/), [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). The [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) still show beta availability, but no newer dated in-window update surfaced there, and the [Claude Code changelog](https://code.claude.com/docs/en/changelog) still shows its latest visible entry on **April 15, 2026**, so there is no fresh Claude Code release to elevate this week.
- **OpenAI shipped a stronger coding model directly into Codex.** [Introducing GPT-5.5](https://openai.com/index/introducing-gpt-5-5/) was published on **April 23, 2026**. OpenAI positions it as its strongest agentic coding model to date, with a reported **82.7% on Terminal-Bench 2.0**, **58.6% on SWE-Bench Pro**, and rollout into **ChatGPT and Codex**.
- **OpenAI also published a concrete orchestration pattern for ticket-driven coding agents.** [An open-source spec for Codex orchestration: Symphony](https://openai.com/index/open-source-codex-orchestration-symphony/) was published on **April 27, 2026**. Its core framing is that the issue tracker becomes the control plane and that any open task should be picked up and completed by an agent, which is a stronger public statement about persistent multi-agent execution than a generic model launch post.

## Trending Snapshot

### GitHub Trending - 2026-04-30

- Metadata: `Source: [GitHub Trending](https://github.com/trending) as checked on **2026-04-30** from the default global repositories page with no extra filters. The board is still a mixed developer-interest feed, not a clean AI-only ranking.`
- [warpdotdev/warp](https://github.com/warpdotdev/warp) - an agentic development environment built around the terminal; one of the clearest board-level signals that the shell is turning into an AI operating surface.
- [mattpocock/skills](https://github.com/mattpocock/skills) - packaged skills for coding workflows, surfacing because specialization layers for coding agents are becoming concrete and reusable.
- [ComposioHQ/awesome-codex-skills](https://github.com/ComposioHQ/awesome-codex-skills) - a curated Codex-skill directory, notable because it treats skill composition as an ecosystem problem rather than a single-repo trick.

## Durable Interpretation

The theme this week is **agent execution becoming a governed workspace contract**.

## Current Proposal Candidates

- Workspace manifests for Harness runs
- Effort and budget as explicit run policy
- Review agents as a named verification lane

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-04-30](../radar/daily/2026-04-30.md)
- [2026-04-28](../radar/daily/2026-04-28.md)
- [2026-04-27](../radar/daily/2026-04-27.md)
- [2026-04-26](../radar/daily/2026-04-26.md)
- [2026-04-24](../radar/daily/2026-04-24.md)
- [2026-04-23](../radar/daily/2026-04-23.md)

### Weekly Synthesis

- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-04-30](../radar/trending/2026-04-30.md)
- [2026-04-28](../radar/trending/2026-04-28.md)
- [2026-04-27](../radar/trending/2026-04-27.md)
- [2026-04-26](../radar/trending/2026-04-26.md)
- [2026-04-24](../radar/trending/2026-04-24.md)
- [2026-04-23](../radar/trending/2026-04-23.md)
