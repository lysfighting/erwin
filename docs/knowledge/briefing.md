# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `66` durable notes are currently indexed across `3` recurring streams.
- `129` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-05-08

- Metadata: `Coverage window: 2026-05-01 to 2026-05-08`
- Anthropic was the clearest first-party mover this week, and the required Anthropic substreams were explicitly checked: [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs overview](https://code.claude.com/docs/en/overview), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). The most important platform change landed on **May 6, 2026**: Claude Platform moved **Multiagent sessions** and **Outcomes** into public beta for Managed Agents, added **Managed Agent webhooks**, added vault credential background refresh for `mcp_oauth`, and improved session/event filtering. This is not a cosmetic docs change; it materially extends Anthropic’s hosted agent runtime surface.
- Anthropic’s Claude Code stream also changed materially inside the same window. The [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows `2.1.132` on **May 6, 2026** and `2.1.133` on **May 7, 2026**. The changes matter operationally: Claude Code now exposes `CLAUDE_CODE_SESSION_ID` to Bash subprocesses, adds `worktree.baseRef` policy for isolation worktrees, passes effort into hooks and Bash via `$CLAUDE_EFFORT`, fixes concurrent-session effort leakage, fixes MCP OAuth proxy handling, and fixes subagents failing to discover project, user, or plugin skills.
- Anthropic’s product/news surface also reinforced the same direction. On **May 5, 2026**, [Agents for financial services](https://www.anthropic.com/news/finance-agents) packaged ten ready-to-run agent templates as Claude Cowork / Claude Code plugins plus Claude Managed Agent cookbooks, and connected them to Microsoft 365 add-ins, connectors, and an MCP app ecosystem. On **May 6, 2026**, [Higher usage limits for Claude and a compute deal with SpaceX](https://www.anthropic.com/news/higher-limits-spacex) doubled Claude Code five-hour limits for paid tiers, removed peak-hour reduction for Pro and Max, and raised API rate limits for Claude Opus models.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-05-13

- Metadata: `Window: 2026-05-06 to 2026-05-13`
- **Anthropic added fast mode for Claude Opus 4.7 on 2026-05-12.** The official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added a **May 12, 2026** entry for `speed: "fast"` on `claude-opus-4-7` under the `fast-mode-2026-02-01` beta header. This is a concrete runtime tradeoff signal for long-horizon agent workloads, not just a model marketing update.
- **Anthropic launched Claude Platform on AWS on 2026-05-11.** The official [Introducing the Claude Platform on AWS](https://claude.com/blog/claude-platform-on-aws) post and the **May 11, 2026** [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) position AWS as a full native surface for Messages, Files, Message Batches, Claude Managed Agents, Agent Skills, code execution, and tool use. The [Claude Code on Claude Platform on AWS docs](https://code.claude.com/docs/en/claude-platform-on-aws) also make Claude Code part of that deployment path.
- **Claude Code shipped an agent-operations surface on 2026-05-11 and followed it with another changelog update on 2026-05-12.** The official [Agent view in Claude Code](https://claude.com/blog/agent-view-in-claude-code) announcement landed on **May 11, 2026**, and the [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows `claude agents` in v2.1.139 on **May 11, 2026** plus v2.1.140 fixes on **May 12, 2026**. Anthropic is making background-session orchestration and supervision more first-class inside the coding surface itself.

## Trending Snapshot

### GitHub Trending - 2026-05-13

- Metadata: `Source: default global [GitHub Trending](https://github.com/trending) checked on **2026-05-13** with no extra filters.`
- [bytedance/UI-TARS-desktop](https://github.com/bytedance/UI-TARS-desktop) - an open-source multimodal AI agent stack tying models to desktop-oriented agent infrastructure.
- [CloakHQ/CloakBrowser](https://github.com/CloakHQ/CloakBrowser) - stealth Chromium automation infrastructure focused on bot-detection resistance.
- [playcanvas/supersplat](https://github.com/playcanvas/supersplat) - a 3D Gaussian splat editor, showing graphics and 3D tooling demand on the default board.

## Durable Interpretation

This week’s theme is **agents becoming governed operating systems, not just smart tool callers**.

## Current Proposal Candidates

- Managed run webhooks and outcomes
- Effort-aware hooks and wrappers
- Enterprise skill and plugin provenance

## Knowledge Network

The machine-readable index now includes `knowledge_network.artifacts`, `knowledge_network.links`, `knowledge_network.threads`, and `knowledge_network.proposal_records`.

- `erwin five-layer operating model` - `erwin` should operate as a five-layer knowledge-to-action loop: Radar and Feedback feed Knowledge, mature Knowledge feeds Proposals, reviewed Proposals feed Implementation, and Implementation results flow back into Feedback and Knowledge. Harness is the main implementation target and feedback source, not a shortcut around the loop.
- `knowledge-system operating discipline` - The knowledge system must stay file-backed and proposal-gated: repeated usage pressure becomes feedback cases first, active durable understanding lives in `knowledge/`, and only mature threads should create proposal pressure.
- `coding-agent self-iteration and paper intake` - OpenAI currently presents self-improvement as eval-backed research engineering, Anthropic presents it more as strong internal dogfooding and workflow automation, and the paper-watch archive now has recurring notes for self-correction, long-horizon degradation, context retrieval, and trajectory diagnosis instead of only one starter note.
- `harness engineering as legibility + eval pressure` - The strongest immediate reusable lesson from recent field signal is not to widen `Harness` runtime blindly, but to keep repo-readable truth, traceability, eval integrity, bounded policy discipline, and monitorability ahead of speculative feature churn. Internal audits, runtime feedback, and newer external monitoring signal now sharpen that thread by showing where visible productization still outruns deeper governance maturity and observability.
- `benchmark integrity and eval awareness` - Benchmark-rich signal is now strong enough that the repo needs a durable distinction between “this result highlights a real weakness” and “this benchmark remains a trustworthy measure,” especially when infrastructure setup and output-only reporting can distort what the score seems to mean.

## Recent History

### Daily Radar

- [2026-05-13](../radar/daily/2026-05-13.md)
- [2026-05-12](../radar/daily/2026-05-12.md)
- [2026-05-11](../radar/daily/2026-05-11.md)
- [2026-05-10](../radar/daily/2026-05-10.md)
- [2026-05-09](../radar/daily/2026-05-09.md)
- [2026-05-06](../radar/daily/2026-05-06.md)

### Weekly Synthesis

- [2026-05-08](../radar/weekly/2026-05-08.md)
- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-05-13](../radar/trending/2026-05-13.md)
- [2026-05-12](../radar/trending/2026-05-12.md)
- [2026-05-11](../radar/trending/2026-05-11.md)
- [2026-05-10](../radar/trending/2026-05-10.md)
- [2026-05-09](../radar/trending/2026-05-09.md)
- [2026-05-06](../radar/trending/2026-05-06.md)
