# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `110` durable notes are currently indexed across `3` recurring streams.
- `173` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-06-05

- Metadata: `Coverage window: 2026-05-29 to 2026-06-05`
- Anthropic again carried the strongest first-party technical signal this week, and this run explicitly checked the required Anthropic surfaces before reaching that conclusion: [anthropic.com/news](https://www.anthropic.com/news), [claude.com](https://claude.com/), [docs.anthropic.com](https://docs.anthropic.com/en/home), the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), the [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), the [Claude Code changelog](https://code.claude.com/docs/en/changelog), and the [Claude apps release notes](https://support.claude.com/en/articles/12138966-release-notes).
- The most important Anthropic platform change landed on **June 2, 2026** in the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview): the **advisor tool** now supports a `max_tokens` parameter, turning advisor output budget into a hard control surface rather than a prompt-only soft hint. The same release also stops billing refusal-only requests that generate no output.
- Anthropic also shipped a dense May 28-29 sequence. On **May 28, 2026**, the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) and [Anthropic launch post](https://www.anthropic.com/news/claude-opus-4-8) introduced **Claude Opus 4.8**, plus mid-conversation system messages, documented refusal `stop_details`, lower prompt-cache minimums, adaptive thinking, and broader Claude Code auto/fast/workflow changes around the same model line. On **May 29, 2026**, Anthropic extended **Claude Managed Agents webhooks, multiagent orchestration, and self-hosted sandboxes** to Claude Platform on AWS.

## Daily Signal Snapshot

### Daily Radar - 2026-06-11

- Metadata: `Window: 2026-06-05 to 2026-06-11`
- Open the daily note for the full signal list.

## Trending Snapshot

### GitHub Trending - 2026-06-11

- Metadata: `Latest GitHub Trending board summary`
- [mvanhorn/last30days-skill](https://github.com/mvanhorn/last30days-skill)
- [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills)
- [phuryn/pm-skills](https://github.com/phuryn/pm-skills)

## Durable Interpretation

This week’s theme is **agent control planes becoming explicit products**.

## Current Proposal Candidates

- Planner / executor split with explicit advisor budget
- First-class runtime object model
- Workflow orchestration tier above subagents

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-06-10](../radar/daily/2026-06-10.md)
- [2026-06-08](../radar/daily/2026-06-08.md)
- [2026-06-06](../radar/daily/2026-06-06.md)
- [2026-06-05](../radar/daily/2026-06-05.md)
- [2026-06-04](../radar/daily/2026-06-04.md)
- [2026-06-03](../radar/daily/2026-06-03.md)

### Weekly Synthesis

- [2026-06-05](../radar/weekly/2026-06-05.md)
- [2026-05-29](../radar/weekly/2026-05-29.md)
- [2026-05-22](../radar/weekly/2026-05-22.md)
- [2026-05-15](../radar/weekly/2026-05-15.md)
- [2026-05-08](../radar/weekly/2026-05-08.md)
- [2026-04-17](../radar/weekly/2026-04-17.md)

### GitHub Trending

- [2026-06-10](../radar/trending/2026-06-10.md)
- [2026-06-08](../radar/trending/2026-06-08.md)
- [2026-06-06](../radar/trending/2026-06-06.md)
- [2026-06-05](../radar/trending/2026-06-05.md)
- [2026-06-04](../radar/trending/2026-06-04.md)
- [2026-06-03](../radar/trending/2026-06-03.md)
