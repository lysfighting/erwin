# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `32` durable notes are currently indexed across `3` recurring streams.
- `95` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-04-17

- Metadata: `Coverage window: 2026-04-10 to 2026-04-17`
- The week’s durable field signal is that frontier agent products are converging on **controlled computer/workspace execution** as a platform primitive. OpenAI’s [The next evolution of the Agents SDK](https://openai.com/index/the-next-evolution-of-the-agents-sdk/) published on **April 15, 2026** explicitly introduces a model-native harness, native sandbox execution, workspace manifests, shell/file-edit tools, MCP, skills, `AGENTS.md`, checkpointing, rehydration, and separation between harness state and compute. OpenAI then broadened the Codex product surface in [Codex for (almost) everything](https://openai.com/index/codex-for-almost-everything/) on **April 16, 2026**, positioning Codex as a persistent, computer-operating development partner rather than only a coding assistant.
- Anthropic remained a high-signal first-party stream. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) show [Claude Opus 4.7](https://www.anthropic.com/news/claude-opus-4-7) launching on **April 16, 2026** with stronger long-running coding, agent, vision, instruction-following, file-system memory, and verification behavior. The same release added `xhigh` effort for Opus 4.7, task budgets in public beta on the Claude Platform API, and Claude Code’s `/ultrareview` review command.
- Anthropic’s required substreams were explicitly checked. The Claude Platform release notes, [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs overview](https://code.claude.com/docs/en/overview), and [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md) were reviewed. I found no new Managed Agents launch comparable to the **April 8, 2026** public beta from last week, but Claude Code changed materially this week: the changelog now tops at `2.1.112`, with `2.1.111` adding Opus 4.7 `xhigh`, `/ultrareview`, `/less-permission-prompts`, skill token sorting, improved provider setup, and OpenTelemetry raw API body logging for debugging.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-04-20

- Metadata: `Window: 2026-04-14 to 2026-04-20`
- **Anthropic moved its agent platform forward with Opus 4.7 and Claude Code runtime changes.** Anthropic's [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) list **April 16, 2026** for Claude Opus 4.7, positioning it for complex reasoning and agentic coding, plus same-window Bedrock availability and **April 14, 2026** deprecation pressure for older Claude 4 models. The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) then shows Opus 4.7 support and runtime hardening across **April 16-18, 2026**, including `xhigh` effort, native binary spawning, network deny-list controls, stalled-subagent failure surfacing, MCP timeout fixes, and security tightening around Bash permissions. Collapsed takeaway: Anthropic is combining stronger model upgrades with sharper agent-runtime policy, observability, and failure behavior.
- **Claude Managed Agents and the advisor tool remain the key Anthropic platform betas to track.** The required Anthropic surface check found no newer dated doc than the platform notes: [Claude Managed Agents](https://platform.claude.com/docs/en/managed-agents/overview) remains the **April 8, 2026** public beta for a managed autonomous-agent harness with sandboxing, built-in tools, and SSE streaming; the [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) remains the **April 9, 2026** public beta for pairing a faster executor model with a higher-intelligence advisor model. These are not new today, but they remain in-window and directly shape harness planning/execution separation.
- **OpenAI made Codex and the Agents SDK more explicitly harness-native.** OpenAI published [The next evolution of the Agents SDK](https://openai.com/index/the-next-evolution-of-the-agents-sdk/) on **April 15, 2026**, adding a model-native harness, native sandbox execution, portable workspace manifests, filesystem/shell/apply-patch tools, MCP, skills, memory, snapshotting, and rehydration. OpenAI then published [Codex for (almost) everything](https://openai.com/index/codex-for-almost-everything/) on **April 16, 2026**, adding desktop computer use, an in-app browser, expanded plugins, SSH/devbox support, memory preview, and automations that can reuse existing conversation threads. Collapsed takeaway: OpenAI is pushing the same agent loop into both SDK infrastructure and end-user developer workflow.

## Trending Snapshot

### GitHub Trending - 2026-04-21

- Metadata: `Source: [GitHub Trending](https://github.com/trending) as captured on 2026-04-21. This run used the default global trending page with no extra filters; the page rendered 10 repositories.`
- [Fincept-Corporation/FinceptTerminal](https://github.com/Fincept-Corporation/FinceptTerminal) - finance application for market analytics, investment research, and economic data tools. Python. 9,618 stars, 3,109 stars today.
- [ruvnet/RuView](https://github.com/ruvnet/RuView) - WiFi-based pose estimation, vital-sign monitoring, and presence detection without camera video. Rust. 48,233 stars, 713 stars today.
- [thunderbird/thunderbolt](https://github.com/thunderbird/thunderbolt) - AI stack emphasizing model choice, local control, and reduced vendor lock-in. TypeScript. 2,834 stars, 675 stars today.

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

- [2026-04-20](../radar/daily/2026-04-20.md)
- [2026-04-19](../radar/daily/2026-04-19.md)
- [2026-04-18](../radar/daily/2026-04-18.md)
- [2026-04-17](../radar/daily/2026-04-17.md)
- [2026-04-16](../radar/daily/2026-04-16.md)
- [2026-04-15](../radar/daily/2026-04-15.md)

### Weekly Synthesis

- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-04-21](../radar/trending/2026-04-21.md)
- [2026-04-20](../radar/trending/2026-04-20.md)
- [2026-04-19](../radar/trending/2026-04-19.md)
- [2026-04-18](../radar/trending/2026-04-18.md)
- [2026-04-17](../radar/trending/2026-04-17.md)
- [2026-04-16](../radar/trending/2026-04-16.md)
