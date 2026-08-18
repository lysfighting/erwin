# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `215` durable notes are currently indexed across `3` recurring streams.
- `278` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI Agent Weekly Synthesis

- Metadata: `Window: 2026-08-08 to 2026-08-15`
- **Anthropic had a materially active week across Claude Platform and Claude Code, even though the Managed Agents and advisor-tool docs themselves stayed quiet.** On **August 11, 2026**, Anthropic’s official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added Compliance API access to transcripts from Cowork and Claude Code sessions that run on users’ machines, plus the `anthropic-workspace-id` response header. On **August 10, 2026**, the same release notes locked in [Claude Sonnet 5 pricing](https://platform.claude.com/docs/en/release-notes/overview) instead of allowing the planned September increase. I also explicitly checked the official [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) and [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) on **August 15, 2026**; no separate in-window Managed Agents or advisor-tool documentation change surfaced there.
- **Claude Code changed repeatedly inside the same window and should be treated as a first-class Anthropic substream, not folded into generic Anthropic coverage.** The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) recorded in-window updates on **August 12, 2026**, **August 13, 2026**, and **August 14, 2026**. Those changes included self-hosted runner hook support, SSE keepalive pings for long thinking pauses, plugin marketplace command sources, default-on subagent forking, cross-session mentions and direct session messaging, GitLab plugin marketplace support, `forward_user_identity` gateway settings, and Bash memory cgroup controls.
- **GitHub’s strongest official signal was packaging and state continuity for coding agents.** On **August 11, 2026**, GitHub published [Copilot memory and Ollama in GitHub Copilot for JetBrains](https://github.blog/changelog/2026-08-11-copilot-memory-and-ollama-in-github-copilot-for-jetbrains/). On **August 12, 2026**, it published [Agent Plugins 1.0 in VS Code, Copilot CLI, and the Copilot app](https://github.blog/changelog/2026-08-12-agent-plugins-1-0-in-vs-code-copilot-cli-and-the-copilot-app/). On **August 13, 2026**, it published [GitHub Copilot weekly releases — August 10](https://github.blog/changelog/2026-08-13-github-copilot-weekly-releases-august-10/), reinforcing that the coding-agent surface is now a managed release stream rather than a one-off feature bucket.

## Daily Signal Snapshot

### Daily Radar - 2026-08-18

- Metadata: `Window: 2026-08-12 to 2026-08-18`
- Open the daily note for the full signal list.

## Trending Snapshot

### GitHub Trending - 2026-08-18

- Metadata: `Latest GitHub Trending board summary`
- [harry0703/MoneyPrinterTurbo](https://github.com/harry0703/MoneyPrinterTurbo)
- [chaitanyagiri/munder-difflin](https://github.com/chaitanyagiri/munder-difflin)
- [akitaonrails/ai-memory](https://github.com/akitaonrails/ai-memory)

## Durable Interpretation

### The frontier is converging on execution governance, not just model capability

## Current Proposal Candidates

- Session Governance Packets For Local and Remote Agent Work
- Portable Skill and Plugin Bundles With Runtime Trust Boundaries
- Speed Class and Reasoning Continuity As Separate Run Controls

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-08-18](../radar/daily/2026-08-18.md)
- [2026-08-17](../radar/daily/2026-08-17.md)
- [2026-08-16](../radar/daily/2026-08-16.md)
- [2026-08-13](../radar/daily/2026-08-13.md)
- [2026-08-12](../radar/daily/2026-08-12.md)
- [2026-08-11](../radar/daily/2026-08-11.md)

### Weekly Synthesis

- [2026-08-15](../radar/weekly/2026-08-15.md)
- [2026-08-08](../radar/weekly/2026-08-08.md)
- [2026-08-01](../radar/weekly/2026-08-01.md)
- [2026-07-25](../radar/weekly/2026-07-25.md)
- [2026-07-18](../radar/weekly/2026-07-18.md)
- [2026-07-10](../radar/weekly/2026-07-10.md)

### GitHub Trending

- [2026-08-18](../radar/trending/2026-08-18.md)
- [2026-08-17](../radar/trending/2026-08-17.md)
- [2026-08-16](../radar/trending/2026-08-16.md)
- [2026-08-13](../radar/trending/2026-08-13.md)
- [2026-08-12](../radar/trending/2026-08-12.md)
- [2026-08-11](../radar/trending/2026-08-11.md)
