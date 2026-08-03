# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `187` durable notes are currently indexed across `3` recurring streams.
- `250` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI Agent Weekly Synthesis

- Metadata: `Window: 2026-07-25 to 2026-08-01`
- **Anthropic was quiet this week after its July 24 wave, and that absence is confirmed rather than assumed.** I explicitly checked the official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code changelog](https://code.claude.com/docs/en/changelog), and [Claude Code What’s New](https://code.claude.com/docs/en/whats-new). No newer in-window Anthropic platform, Managed Agents, advisor-tool, or Claude Code launch surfaced after **July 24, 2026**. The latest still-relevant official Anthropic changes remain **July 24, 2026** Claude Opus 5 plus mid-conversation tool changes on the Platform release notes, **June 9, 2026** Managed Agents scheduled deployments / vault env-var credentials / `session_thread_id`, and **June 2, 2026** advisor tool `max_tokens`.
- **GitHub had the strongest in-window execution-layer movement.** On **July 29, 2026**, GitHub announced [Copilot code review: agent skills and MCP now generally available](https://github.blog/changelog/2026-07-29-copilot-code-review-agent-skills-and-mcp-now-generally-available/) and [auto model selection in GitHub Copilot](https://github.blog/changelog/2026-07-29-auto-model-selection-in-github-copilot-for-vs-code-is-now-in-public-preview/). On **July 30, 2026**, it added [knowledge bases for the Copilot coding agent](https://github.blog/changelog/2026-07-30-github-copilot-coding-agent-now-supports-knowledge-bases/) and [custom instructions support for secret scanning](https://github.blog/changelog/2026-07-30-custom-instructions-support-for-secret-scanning-is-now-generally-available/). On **July 31, 2026**, it launched [Issue forms with GitHub Copilot](https://github.blog/changelog/2026-07-31-issue-forms-with-github-copilot-now-in-public-preview/) and [GitHub Models public preview](https://github.blog/changelog/2026-07-31-github-models-are-now-available-in-public-preview/).
- **Google’s official signal was about evals and grounded retrieval inside its enterprise agent platform.** On **July 31, 2026**, Google published [Adding grading and safety evaluation capabilities in Gemini Enterprise Agent Platform](https://developers.googleblog.com/adding-grading-and-safety-evaluation-capabilities-in-gemini-enterprise-agent-platform/) and [Grounding with Parallel Web Search in Gemini Enterprise Agent Platform](https://developers.googleblog.com/expanding-choice-in-gemini-enterprise-agent-platform-introducing-grounding-with-parallel-web-search/). The first matters because it productizes agent-eval loops in the platform itself; the second matters because grounding-provider choice is being exposed as an explicit operator surface rather than hidden infrastructure.

## Daily Signal Snapshot

### Daily Radar - 2026-08-03

- Metadata: `Window: 2026-07-28 to 2026-08-03`
- Open the daily note for the full signal list.

## Trending Snapshot

### GitHub Trending - 2026-08-03

- Metadata: `Latest GitHub Trending board summary`
- [zhaoxuya520/reverse-skill](https://github.com/zhaoxuya520/reverse-skill)
- [TencentCloud/TencentDB-Agent-Memory](https://github.com/TencentCloud/TencentDB-Agent-Memory)
- [Panniantong/Agent-Reach](https://github.com/Panniantong/Agent-Reach)

## Durable Interpretation

### The control plane is moving upward: from model access to governed agent product surfaces

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

- [2026-08-03](../radar/daily/2026-08-03.md)
- [2026-08-02](../radar/daily/2026-08-02.md)
- [2026-08-01](../radar/daily/2026-08-01.md)
- [2026-07-30](../radar/daily/2026-07-30.md)
- [2026-07-29](../radar/daily/2026-07-29.md)
- [2026-07-28](../radar/daily/2026-07-28.md)

### Weekly Synthesis

- [2026-08-01](../radar/weekly/2026-08-01.md)
- [2026-07-25](../radar/weekly/2026-07-25.md)
- [2026-07-18](../radar/weekly/2026-07-18.md)
- [2026-07-10](../radar/weekly/2026-07-10.md)
- [2026-07-03](../radar/weekly/2026-07-03.md)
- [2026-06-26](../radar/weekly/2026-06-26.md)

### GitHub Trending

- [2026-08-03](../radar/trending/2026-08-03.md)
- [2026-08-02](../radar/trending/2026-08-02.md)
- [2026-08-01](../radar/trending/2026-08-01.md)
- [2026-07-30](../radar/trending/2026-07-30.md)
- [2026-07-29](../radar/trending/2026-07-29.md)
- [2026-07-28](../radar/trending/2026-07-28.md)
