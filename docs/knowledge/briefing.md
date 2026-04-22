# Current Knowledge Briefing

> Generated from the current radar `latest.md` set. Use this file when you want a fast text read; use [../index.html](../index.html) for the visual surface.

## Open First

- [Visual surface](../index.html)
- [Latest weekly source](../radar/weekly/latest.md)
- [Machine-readable index](index.json)
- [Operating model](operating-model.md)

## At A Glance

- `35` durable notes are currently indexed across `3` recurring streams.
- `98` knowledge-network artifacts, `267` typed local links, and `5` active threads are available in `index.json`.
- The active repo path is `radar -> feedback -> knowledge -> proposals -> implementation`.
- `radar` is intake only; this briefing is a convenience surface, not the canonical logic.

## Weekly Lead

### AI/Agent Weekly Synthesis — 2026-04-17

- Metadata: `Coverage window: 2026-04-10 to 2026-04-17`
- The week’s durable field signal is that frontier agent products are converging on **controlled computer/workspace execution** as a platform primitive. OpenAI’s [The next evolution of the Agents SDK](https://openai.com/index/the-next-evolution-of-the-agents-sdk/) published on **April 15, 2026** explicitly introduces a model-native harness, native sandbox execution, workspace manifests, shell/file-edit tools, MCP, skills, `AGENTS.md`, checkpointing, rehydration, and separation between harness state and compute. OpenAI then broadened the Codex product surface in [Codex for (almost) everything](https://openai.com/index/codex-for-almost-everything/) on **April 16, 2026**, positioning Codex as a persistent, computer-operating development partner rather than only a coding assistant.
- Anthropic remained a high-signal first-party stream. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) show [Claude Opus 4.7](https://www.anthropic.com/news/claude-opus-4-7) launching on **April 16, 2026** with stronger long-running coding, agent, vision, instruction-following, file-system memory, and verification behavior. The same release added `xhigh` effort for Opus 4.7, task budgets in public beta on the Claude Platform API, and Claude Code’s `/ultrareview` review command.
- Anthropic’s required substreams were explicitly checked. The Claude Platform release notes, [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs overview](https://code.claude.com/docs/en/overview), and [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md) were reviewed. I found no new Managed Agents launch comparable to the **April 8, 2026** public beta from last week, but Claude Code changed materially this week: the changelog now tops at `2.1.112`, with `2.1.111` adding Opus 4.7 `xhigh`, `/ultrareview`, `/less-permission-prompts`, skill token sorting, improved provider setup, and OpenTelemetry raw API body logging for debugging.

## Daily Signal Snapshot

### AI Agent Daily Radar - 2026-04-22

- Metadata: `Window: 2026-04-15 to 2026-04-22`
- **Anthropic combined frontier-model launch, app-surface expansion, and API retirement pressure in the same week.** The official [Anthropic newsroom](https://www.anthropic.com/news) and [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) show [Claude Opus 4.7](https://www.anthropic.com/news/claude-opus-4-7) on **2026-04-16**, [Claude Design by Anthropic Labs](https://support.claude.com/en/articles/12138966-release-notes) on **2026-04-17**, and Claude Haiku 3 retirement on **2026-04-20**, alongside the earlier **2026-04-14** deprecation notice for older Claude 4 model IDs. Collapsed takeaway: Anthropic is tightening the upgrade loop while broadening higher-end agent and output surfaces.
- **Claude Code kept shipping runtime hardening aimed at long-running, multi-agent operator workflows.** The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows **v2.1.113 on 2026-04-17** adding a native binary path, domain-deny sandbox controls, `/ultrareview` launch improvements, better `/loop`, and clearer subagent stall failure handling, followed by **v2.1.114 on 2026-04-18** with a teammate-permission crash fix. Collapsed takeaway: Anthropic is still investing directly in coding-agent durability, isolation policy, and delegated-run recovery.
- **OpenAI pushed the same harness pattern through both the developer SDK and the end-user Codex product.** OpenAI published [The next evolution of the Agents SDK](https://openai.com/index/the-next-evolution-of-the-agents-sdk/) on **2026-04-15**, adding configurable memory, sandbox-aware orchestration, Codex-like filesystem tools, MCP, skills, and `AGENTS.md` support, then followed with [Codex for (almost) everything](https://openai.com/index/codex-for-almost-everything/) on **2026-04-16**, extending Codex with computer use, parallel agents on-device, SSH/devbox support, browser access, memory, and reusable automations. Collapsed takeaway: OpenAI is converging product and SDK around one durable agent operating model.

## Trending Snapshot

### GitHub Trending - 2026-04-22

- Metadata: `Source: [GitHub Trending](https://github.com/trending) as captured on 2026-04-22. This run used the default global trending page with no extra filters; the page rendered 9 repositories.`
- [Fincept-Corporation/FinceptTerminal](https://github.com/Fincept-Corporation/FinceptTerminal) - finance application for market analytics, investment research, and economic data exploration. Python. 11,679 stars, 2,548 stars today.
- [thunderbird/thunderbolt](https://github.com/thunderbird/thunderbolt) - AI stack focused on model choice, local ownership, and avoiding vendor lock-in. TypeScript. 3,490 stars, 596 stars today.
- [zilliztech/claude-context](https://github.com/zilliztech/claude-context) - MCP-based code search layer for Claude Code and other coding agents. TypeScript. 6,637 stars, 169 stars today.

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

- [2026-04-22](../radar/daily/2026-04-22.md)
- [2026-04-21](../radar/daily/2026-04-21.md)
- [2026-04-20](../radar/daily/2026-04-20.md)
- [2026-04-19](../radar/daily/2026-04-19.md)
- [2026-04-18](../radar/daily/2026-04-18.md)
- [2026-04-17](../radar/daily/2026-04-17.md)

### Weekly Synthesis

- [2026-04-17](../radar/weekly/2026-04-17.md)
- [2026-04-10](../radar/weekly/2026-04-10.md)
- [2026-04-03](../radar/weekly/2026-04-03.md)
- [2026-03-30](../radar/weekly/2026-03-30.md)

### GitHub Trending

- [2026-04-22](../radar/trending/2026-04-22.md)
- [2026-04-21](../radar/trending/2026-04-21.md)
- [2026-04-20](../radar/trending/2026-04-20.md)
- [2026-04-19](../radar/trending/2026-04-19.md)
- [2026-04-18](../radar/trending/2026-04-18.md)
- [2026-04-17](../radar/trending/2026-04-17.md)
