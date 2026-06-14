# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `117` durable notes are currently indexed across `3` recurring streams.
- `180` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-06-12

- Metadata: `Coverage window: 2026-06-05 to 2026-06-12`
- Anthropic again carried the strongest first-party agent-platform signal this week, and this run explicitly checked the required Anthropic surfaces before reaching that conclusion: [anthropic.com/news](https://www.anthropic.com/news), [claude.com](https://claude.com/blog-category/announcements), [docs.anthropic.com / Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), the [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), the [Claude Code changelog](https://code.claude.com/docs/en/changelog), and the [Claude apps release notes](https://support.claude.com/en/articles/12138966-release-notes).
- Anthropic’s biggest in-window platform move landed on **June 9, 2026**. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) introduced **Claude Fable 5** and **Claude Mythos 5** with default 1M context, 128k max output tokens, always-on adaptive thinking, refusal handling changes, and a new `fallbacks` parameter for re-running refusal cases on another model. The same day, the [Claude apps release notes](https://support.claude.com/en/articles/12138966-release-notes) recorded the Claude Fable 5 launch on the user-facing product surface.
- Anthropic’s required **Claude Managed Agents** check was materially positive this week. On **June 9, 2026**, the [Claude Managed Agents announcement](https://claude.com/blog/whats-new-in-claude-managed-agents) added **scheduled deployments** and **environment variables in vaults** in public beta, giving Managed Agents a first-party recurring-run primitive and a safer way to authenticate CLIs and other tools without handing raw credentials to the model.

## Daily Signal Snapshot

### Daily Radar - 2026-06-14

- Metadata: `Window: 2026-06-08 to 2026-06-14`
- Open the daily note for the full signal list.

## Trending Snapshot

### GitHub Trending - 2026-06-14

- Metadata: `Latest GitHub Trending board summary`
- [iptv-org/iptv](https://github.com/iptv-org/iptv)
- [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills)
- [chatwoot/chatwoot](https://github.com/chatwoot/chatwoot)

## Durable Interpretation

This week’s theme is **agents becoming schedulable, inspectable, and credential-safe operating systems instead of one-shot helpers**.

## Current Proposal Candidates

- Recurring deployment lane for autonomous runs
- Credential-bound CLI tool lane
- Nested orchestration lineage and pending-work introspection

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-06-14](../radar/daily/2026-06-14.md)
- [2026-06-12](../radar/daily/2026-06-12.md)
- [2026-06-11](../radar/daily/2026-06-11.md)
- [2026-06-10](../radar/daily/2026-06-10.md)
- [2026-06-08](../radar/daily/2026-06-08.md)
- [2026-06-06](../radar/daily/2026-06-06.md)

### Weekly Synthesis

- [2026-06-12](../radar/weekly/2026-06-12.md)
- [2026-06-05](../radar/weekly/2026-06-05.md)
- [2026-05-29](../radar/weekly/2026-05-29.md)
- [2026-05-22](../radar/weekly/2026-05-22.md)
- [2026-05-15](../radar/weekly/2026-05-15.md)
- [2026-05-08](../radar/weekly/2026-05-08.md)

### GitHub Trending

- [2026-06-14](../radar/trending/2026-06-14.md)
- [2026-06-12](../radar/trending/2026-06-12.md)
- [2026-06-11](../radar/trending/2026-06-11.md)
- [2026-06-10](../radar/trending/2026-06-10.md)
- [2026-06-08](../radar/trending/2026-06-08.md)
- [2026-06-06](../radar/trending/2026-06-06.md)
