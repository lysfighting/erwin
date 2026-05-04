# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `53` durable notes are currently indexed across `3` recurring streams.
- `116` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-04-17

- Metadata: `Coverage window: 2026-04-10 to 2026-04-17`
- The week’s durable field signal is that frontier agent products are converging on **controlled computer/workspace execution** as a platform primitive. OpenAI’s [The next evolution of the Agents SDK](https://openai.com/index/the-next-evolution-of-the-agents-sdk/) published on **April 15, 2026** explicitly introduces a model-native harness, native sandbox execution, workspace manifests, shell/file-edit tools, MCP, skills, `AGENTS.md`, checkpointing, rehydration, and separation between harness state and compute. OpenAI then broadened the Codex product surface in [Codex for (almost) everything](https://openai.com/index/codex-for-almost-everything/) on **April 16, 2026**, positioning Codex as a persistent, computer-operating development partner rather than only a coding assistant.
- Anthropic remained a high-signal first-party stream. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) show [Claude Opus 4.7](https://www.anthropic.com/news/claude-opus-4-7) launching on **April 16, 2026** with stronger long-running coding, agent, vision, instruction-following, file-system memory, and verification behavior. The same release added `xhigh` effort for Opus 4.7, task budgets in public beta on the Claude Platform API, and Claude Code’s `/ultrareview` review command.
- Anthropic’s required substreams were explicitly checked. The Claude Platform release notes, [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs overview](https://code.claude.com/docs/en/overview), and [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md) were reviewed. I found no new Managed Agents launch comparable to the **April 8, 2026** public beta from last week, but Claude Code changed materially this week: the changelog now tops at `2.1.112`, with `2.1.111` adding Opus 4.7 `xhigh`, `/ultrareview`, `/less-permission-prompts`, skill token sorting, improved provider setup, and OpenTelemetry raw API body logging for debugging.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-05-04

- Metadata: `Window: 2026-04-27 to 2026-05-04`
- **Anthropic changed the practical long-context routing rule on Claude Platform.** The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) were updated on **April 30, 2026** to retire the 1M-token beta header for Claude Sonnet 4 and 4.5. Long-context traffic now needs to move to [Claude Sonnet 4.6 / Claude Opus 4.6](https://platform.claude.com/docs/en/release-notes/overview) instead of relying on the old header path.
- **Anthropic published a concrete Claude Code operator doctrine around prompt caching.** [Lessons from building Claude Code: Prompt caching is everything](https://claude.com/blog/lessons-from-building-claude-code-prompt-caching-is-everything) was published on **April 30, 2026**. The strongest takeaway is that Anthropic is now documenting cache-safe agent architecture as a first-class systems concern: stable prompt prefixes, no mid-session tool churn, cache-safe compaction, and explicit monitoring of cache hit rates.
- **Anthropic is exporting Claude API operating knowledge as a portable skill layer.** [Claude API skill now in CodeRabbit, JetBrains, Resolve AI, and Warp](https://claude.com/blog/claude-api-skill) was published on **April 29, 2026**. The notable move is not just distribution; it is Anthropic treating migration rules, managed-agent onboarding, and prompt-caching practice as reusable productized doctrine.

## Trending Snapshot

### GitHub Trending - 2026-05-04

- Metadata: `Source: [GitHub Trending](https://github.com/trending) as checked on **2026-05-04** from the default global repositories page with no extra filters. This remains a mixed developer-interest board rather than an AI-only ranking.`
- [ruvnet/ruflo](https://github.com/ruvnet/ruflo) - a Claude-oriented orchestration platform that is surfacing because multi-agent coordination itself is becoming a product category.
- [TauricResearch/TradingAgents](https://github.com/TauricResearch/TradingAgents) - a multi-agent LLM financial trading framework, notable because agent runtimes are continuing to expand into applied domain stacks.
- [soxoj/maigret](https://github.com/soxoj/maigret) - a large OSINT utility, included because it keeps the summary faithful to the mixed board rather than pretending everything is AI infrastructure.

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

- [2026-05-04](../radar/daily/2026-05-04.md)
- [2026-05-02](../radar/daily/2026-05-02.md)
- [2026-05-01](../radar/daily/2026-05-01.md)
- [2026-04-30](../radar/daily/2026-04-30.md)
- [2026-04-28](../radar/daily/2026-04-28.md)
- [2026-04-27](../radar/daily/2026-04-27.md)

### Weekly Synthesis

- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-05-04](../radar/trending/2026-05-04.md)
- [2026-05-02](../radar/trending/2026-05-02.md)
- [2026-05-01](../radar/trending/2026-05-01.md)
- [2026-04-30](../radar/trending/2026-04-30.md)
- [2026-04-28](../radar/trending/2026-04-28.md)
- [2026-04-27](../radar/trending/2026-04-27.md)
