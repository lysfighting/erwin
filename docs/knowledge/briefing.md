# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `43` durable notes are currently indexed across `3` recurring streams.
- `106` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-04-17

- Metadata: `Coverage window: 2026-04-10 to 2026-04-17`
- The week’s durable field signal is that frontier agent products are converging on **controlled computer/workspace execution** as a platform primitive. OpenAI’s [The next evolution of the Agents SDK](https://openai.com/index/the-next-evolution-of-the-agents-sdk/) published on **April 15, 2026** explicitly introduces a model-native harness, native sandbox execution, workspace manifests, shell/file-edit tools, MCP, skills, `AGENTS.md`, checkpointing, rehydration, and separation between harness state and compute. OpenAI then broadened the Codex product surface in [Codex for (almost) everything](https://openai.com/index/codex-for-almost-everything/) on **April 16, 2026**, positioning Codex as a persistent, computer-operating development partner rather than only a coding assistant.
- Anthropic remained a high-signal first-party stream. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) show [Claude Opus 4.7](https://www.anthropic.com/news/claude-opus-4-7) launching on **April 16, 2026** with stronger long-running coding, agent, vision, instruction-following, file-system memory, and verification behavior. The same release added `xhigh` effort for Opus 4.7, task budgets in public beta on the Claude Platform API, and Claude Code’s `/ultrareview` review command.
- Anthropic’s required substreams were explicitly checked. The Claude Platform release notes, [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs overview](https://code.claude.com/docs/en/overview), and [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md) were reviewed. I found no new Managed Agents launch comparable to the **April 8, 2026** public beta from last week, but Claude Code changed materially this week: the changelog now tops at `2.1.112`, with `2.1.111` adding Opus 4.7 `xhigh`, `/ultrareview`, `/less-permission-prompts`, skill token sorting, improved provider setup, and OpenTelemetry raw API body logging for debugging.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-04-27

- Metadata: `Window: 2026-04-20 to 2026-04-27`
- **Anthropic kept shipping operator-layer controls across Claude Platform and Claude Code.** The official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added the **Rate Limits API** on **2026-04-24** and put **memory for Claude Managed Agents** into public beta on **2026-04-23**. The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) also shows **v2.1.118 on 2026-04-23** adding visual mode, custom themes, and direct MCP-tool invocation from hooks. Collapsed takeaway: Anthropic is still expanding the control plane around agent execution, not only model capability.
- **OpenAI paired a frontier coding-model launch with enterprise activation.** [Introducing GPT-5.5](https://openai.com/index/introducing-gpt-5-5/) was published on **2026-04-23**, with an explicit **2026-04-24** update noting API availability, and [Scaling Codex to enterprises worldwide](https://openai.com/index/scaling-codex-to-enterprises-worldwide/) landed on **2026-04-21** with Codex Labs plus GSI rollout. Collapsed takeaway: OpenAI is packaging model quality and deployment acceleration together.
- **OpenAI also shipped privacy infrastructure as a practical agent-stack primitive.** [OpenAI Privacy Filter](https://openai.com/index/introducing-openai-privacy-filter/) was published on **2026-04-22** as an open-weight PII masking model that can run locally for high-throughput privacy workflows. Collapsed takeaway: privacy-preserving preprocessing is becoming a first-class building block around agents and logs.

## Trending Snapshot

### GitHub Trending - 2026-04-27

- Metadata: `Source: [GitHub Trending](https://github.com/trending) as captured on **2026-04-27** from the default global repositories page with no extra filters. The visible board mixes coding-agent infrastructure with broader developer, systems, and general-interest repositories.`
- [Alishahryar1/free-claude-code](https://github.com/Alishahryar1/free-claude-code) - Claude-adjacent wrapper/access project. Python. 13,242 stars, 1,694 today.
- [trycua/cua](https://github.com/trycua/cua) - computer-use-agent infrastructure with sandboxes, SDKs, and benchmarks. HTML. 14,320 stars, 200 today.
- [gastownhall/beads](https://github.com/gastownhall/beads) - memory upgrade layer for coding agents. Go. 21,606 stars, 133 today.

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

- [2026-04-27](../radar/daily/2026-04-27.md)
- [2026-04-26](../radar/daily/2026-04-26.md)
- [2026-04-24](../radar/daily/2026-04-24.md)
- [2026-04-23](../radar/daily/2026-04-23.md)
- [2026-04-22](../radar/daily/2026-04-22.md)
- [2026-04-21](../radar/daily/2026-04-21.md)

### Weekly Synthesis

- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-04-27](../radar/trending/2026-04-27.md)
- [2026-04-26](../radar/trending/2026-04-26.md)
- [2026-04-24](../radar/trending/2026-04-24.md)
- [2026-04-23](../radar/trending/2026-04-23.md)
- [2026-04-22](../radar/trending/2026-04-22.md)
- [2026-04-21](../radar/trending/2026-04-21.md)
