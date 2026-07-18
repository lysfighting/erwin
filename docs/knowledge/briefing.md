# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `155` durable notes are currently indexed across `3` recurring streams.
- `218` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI Agent Weekly Synthesis

- Metadata: `Window: 2026-07-11 to 2026-07-18`
- **Anthropic** had the most important in-window platform signal. On **July 15, 2026**, the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added **mid-conversation system messages** on Claude Fable 5, Claude Mythos 5, and Claude Opus 4.8, and expanded the **Admin API** for Claude Enterprise organization management. On the required Anthropic sub-surfaces, I explicitly checked the official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents docs](https://platform.claude.com/docs/en/agents-and-tools/managed-agents/overview), [advisor tool docs](https://docs.anthropic.com/en/agents-and-tools/tool-use/advisor-tool), [Claude Code changelog](https://code.claude.com/docs/en/changelog), and [Claude Code docs](https://code.claude.com/docs/en/whats-new). I did **not** find a new in-window Managed Agents or advisor-tool capability launch after the large **June 2 to June 9, 2026** burst, but Claude Code itself changed materially on **July 17, 2026**.
- **Claude Code** had the clearest Anthropic execution-surface movement this week. The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) and [Week 29 digest](https://code.claude.com/docs/en/whats-new/2026-w29) show that on **July 17, 2026** `/fork` changed from an in-session subagent pattern into a **new background session** with its own `claude agents` row, `/subtask` became the in-session replacement, plugins in `.claude/skills` now auto-load, and `claude plugin init` now scaffolds plugins. That is a real control-plane shift, not just bugfix noise.
- **OpenAI** concentrated this week on the management and packaging layer around longer-running work. The official [ChatGPT Enterprise & Edu release notes](https://help.openai.com/en/articles/10128477-chatgpt-enterprise-edu-release-notes) added **workspace-scoped Admin keys and 120 days of analytics history** on **July 16, 2026**, plus **EKM support for apps with sync** on **July 15, 2026**. The new [Skills in ChatGPT](https://help.openai.com/en/articles/20001066-skills-in-chatgpt) page published on **July 18, 2026** makes workspace-managed skill governance explicit, while [Apps in ChatGPT](https://help.openai.com/en/articles/11487775-connectors-in-chatgpt) continues to position the **Plugin Directory** as the discovery surface across ChatGPT and Codex after the **July 9, 2026** migration.

## Daily Signal Snapshot

### Daily Radar - 2026-07-17

- Metadata: `Window: 2026-07-10 to 2026-07-17`
- Open the daily note for the full signal list.

## Trending Snapshot

### GitHub Trending - 2026-07-17

- Metadata: `Latest GitHub Trending board summary`
- [codecrafters-io/build-your-own-x](https://github.com/codecrafters-io/build-your-own-x)
- [PostHog/posthog](https://github.com/PostHog/posthog)
- [github/copilot-sdk](https://github.com/github/copilot-sdk)

## Durable Interpretation

### The control plane is becoming the product

## Current Proposal Candidates

- No proposal candidates extracted from the latest weekly synthesis yet.

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-07-17](../radar/daily/2026-07-17.md)
- [2026-07-16](../radar/daily/2026-07-16.md)
- [2026-07-15](../radar/daily/2026-07-15.md)
- [2026-07-12](../radar/daily/2026-07-12.md)
- [2026-07-11](../radar/daily/2026-07-11.md)
- [2026-07-09](../radar/daily/2026-07-09.md)

### Weekly Synthesis

- [2026-07-18](../radar/weekly/2026-07-18.md)
- [2026-07-10](../radar/weekly/2026-07-10.md)
- [2026-07-03](../radar/weekly/2026-07-03.md)
- [2026-06-26](../radar/weekly/2026-06-26.md)
- [2026-06-12](../radar/weekly/2026-06-12.md)
- [2026-06-05](../radar/weekly/2026-06-05.md)

### GitHub Trending

- [2026-07-17](../radar/trending/2026-07-17.md)
- [2026-07-16](../radar/trending/2026-07-16.md)
- [2026-07-15](../radar/trending/2026-07-15.md)
- [2026-07-12](../radar/trending/2026-07-12.md)
- [2026-07-11](../radar/trending/2026-07-11.md)
- [2026-07-09](../radar/trending/2026-07-09.md)
