# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `132` durable notes are currently indexed across `3` recurring streams.
- `195` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-06-26

- Metadata: `Coverage window: 2026-06-19 to 2026-06-26`
- Anthropic again carried the strongest first-party agent-platform signal this week, and this run explicitly checked the required Anthropic surfaces before reaching that conclusion: [anthropic.com/news](https://www.anthropic.com/news), [claude.com/blog](https://claude.com/blog), [docs.anthropic.com / Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), the [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog).
- Anthropic’s biggest official platform change in-window landed on **June 18, 2026** in the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview): SDK support for `code_execution_20260120`, which adds REPL state persistence and becomes the minimum version for programmatic tool calling. That is not a cosmetic tool bump. It upgrades code execution from stateless cell running toward resumable, agent-friendly execution.
- Anthropic’s required **Claude Managed Agents** check stayed materially positive this week even without a brand-new standalone launch post. The [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) still centers a durable runtime contract of **agent**, **environment**, **session**, and **events**, and explicitly positions the product for long-running, asynchronous work with persistent filesystems, server-side event history, and steering mid-execution. The June 18 code-execution change strengthens that runtime by making programmatic tool calling and REPL persistence easier to operationalize inside agent sessions.

## Daily Signal Snapshot

### Daily Radar - 2026-07-02

- Metadata: `Window: 2026-06-25 to 2026-07-02`
- Open the daily note for the full signal list.

## Trending Snapshot

### GitHub Trending - 2026-07-02

- Metadata: `Latest GitHub Trending board summary`
- [usestrix/strix](https://github.com/usestrix/strix)
- [HKUDS/Vibe-Trading](https://github.com/HKUDS/Vibe-Trading)
- [facebook/astryx](https://github.com/facebook/astryx)

## Durable Interpretation

This week’s theme is **agent runtime governance moving from prompts into platform objects**.

## Current Proposal Candidates

- Steering contract as a first-class run object
- Stateful execution lane for coding and tool workflows
- Identity- and auth-aware capability bundles

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-07-02](../radar/daily/2026-07-02.md)
- [2026-06-28](../radar/daily/2026-06-28.md)
- [2026-06-26](../radar/daily/2026-06-26.md)
- [2026-06-25](../radar/daily/2026-06-25.md)
- [2026-06-17](../radar/daily/2026-06-17.md)
- [2026-06-16](../radar/daily/2026-06-16.md)

### Weekly Synthesis

- [2026-06-26](../radar/weekly/2026-06-26.md)
- [2026-06-12](../radar/weekly/2026-06-12.md)
- [2026-06-05](../radar/weekly/2026-06-05.md)
- [2026-05-29](../radar/weekly/2026-05-29.md)
- [2026-05-22](../radar/weekly/2026-05-22.md)
- [2026-05-15](../radar/weekly/2026-05-15.md)

### GitHub Trending

- [2026-07-02](../radar/trending/2026-07-02.md)
- [2026-06-28](../radar/trending/2026-06-28.md)
- [2026-06-26](../radar/trending/2026-06-26.md)
- [2026-06-25](../radar/trending/2026-06-25.md)
- [2026-06-17](../radar/trending/2026-06-17.md)
- [2026-06-16](../radar/trending/2026-06-16.md)
