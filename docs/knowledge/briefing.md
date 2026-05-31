# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `95` durable notes are currently indexed across `3` recurring streams.
- `158` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-05-29

- Metadata: `Coverage window: 2026-05-22 to 2026-05-29`
- Anthropic carried the clearest first-party technical signal this week, and this run explicitly checked the required Anthropic surfaces before reaching that conclusion: [anthropic.com/news](https://www.anthropic.com/news), [claude.com](https://claude.com/), [docs.anthropic.com](https://docs.anthropic.com/en/home), the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), the [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), the [Claude Code docs overview](https://code.claude.com/docs/en/overview), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog).
- On **May 28, 2026**, Anthropic launched [Claude Opus 4.8](https://www.anthropic.com/news/claude-opus-4-8) and the same date landed in the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview). The important platform-level changes were not just the model SKU: Anthropic added **fine-grained prompt caching** controls, **1M-token context windows** on request for Claude Sonnet 4 and Opus 4.1, **fivex higher output limits** for Opus 4.1 and Sonnet 4, **global endpoint pricing for Claude Code**, and **GA status for web fetch**.
- Anthropic’s managed-agent surface still moved even without a fresh standalone managed-agents launch post. The current [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) continues to frame **agent**, **environment**, **session**, and **events** as the runtime core, while the surrounding release-note stream keeps making that runtime more configurable. The absence of a separate managed-agents announcement this week does not mean the surface was quiet; the control plane around it changed materially.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-05-31

- Metadata: `Window: 2026-05-24 to 2026-05-31`
- **Anthropic pushed the densest in-window agent-platform update across the required official surfaces on 2026-05-28 and 2026-05-29.** The [Claude Apps release notes](https://support.claude.com/en/articles/12138966-release-notes) added **Claude Opus 4.8** and connector-role controls on **2026-05-28**, while the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added **Claude Managed Agents webhooks, multiagent orchestration, and self-hosted sandboxes on Claude Platform on AWS** on **2026-05-29**. I explicitly checked the required Anthropic surfaces this run: Claude Platform release notes, Claude Managed Agents docs, advisor tool docs, Claude Code docs, Claude Code changelog/releases, and Claude Apps release notes.
- **OpenAI expanded both Codex control surfaces and workspace-agent governance on 2026-05-28 and 2026-05-29.** The official [ChatGPT Enterprise & Edu release notes](https://help.openai.com/en/articles/10128477-chatgpt-enterprise-ve-edu-s%C3%BCr%C3%BCm-notlar%C4%B1) added **Windows Computer Use**, **remote control for Windows hosts**, and **GitHub Enterprise app template support for Codex** on **2026-05-29**; on **2026-05-28**, the same surface added **GPT-5.5 support, reasoning-effort controls, and action safeguards for workspace agents**.
- **Google’s I/O 2026 developer push remains the clearest Google agent-platform move still inside the 7-day window.** In [Building the agentic future: Developer highlights from I/O 2026](https://blog.google/innovation-and-ai/technology/developers-tools/google-io-2026-developer-highlights/) on **2026-05-19**, Google introduced **Gemini 3.5 Flash**, **Managed Agents in the Gemini API**, and the **Google Antigravity 2.0 desktop app** for parallel agents, scheduled tasks, and SDK/CLI-based agent workflows.

## Trending Snapshot

### GitHub Trending - 2026-05-31

- Metadata: `Source: default global [GitHub Trending](https://github.com/trending) checked on **2026-05-31** with no extra filters.`
- [microsoft/markitdown](https://github.com/microsoft/markitdown): document and Office-file conversion into Markdown; a practical ingestion primitive that keeps resurfacing.
- [harry0703/MoneyPrinterTurbo](https://github.com/harry0703/MoneyPrinterTurbo): one-click AI short-video generation; high star velocity suggests broad creator demand, not just niche experimentation.
- [anthropics/claude-code](https://github.com/anthropics/claude-code): flagship terminal coding agent; trending at the same time Anthropic is shipping new Claude Platform and Managed Agents changes.

## Durable Interpretation

This week’s theme is **governed scale-up of agent runtime state**.

## Current Proposal Candidates

- Run state budget contract
- Reset versus resume lifecycle states
- Memory-scope governance layer

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-05-31](../radar/daily/2026-05-31.md)
- [2026-05-29](../radar/daily/2026-05-29.md)
- [2026-05-28](../radar/daily/2026-05-28.md)
- [2026-05-25](../radar/daily/2026-05-25.md)
- [2026-05-24](../radar/daily/2026-05-24.md)
- [2026-05-23](../radar/daily/2026-05-23.md)

### Weekly Synthesis

- [2026-05-29](../radar/weekly/2026-05-29.md)
- [2026-05-22](../radar/weekly/2026-05-22.md)
- [2026-05-15](../radar/weekly/2026-05-15.md)
- [2026-05-08](../radar/weekly/2026-05-08.md)
- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)

### GitHub Trending

- [2026-05-31](../radar/trending/2026-05-31.md)
- [2026-05-29](../radar/trending/2026-05-29.md)
- [2026-05-28](../radar/trending/2026-05-28.md)
- [2026-05-25](../radar/trending/2026-05-25.md)
- [2026-05-24](../radar/trending/2026-05-24.md)
- [2026-05-23](../radar/trending/2026-05-23.md)
