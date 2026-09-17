# Daily Radar - 2026-09-16

Window: 2026-09-10 to 2026-09-16 (America/Los_Angeles). Checked 2026-09-16. One search pass per publisher family plus direct official-page checks. Comparison baseline is the last persisted note, September 14; no September 15 daily note was present. Releases below are vendor evidence, not local validation.

## New Signals

1. **Claude Platform — explicit context compaction.** [Release notes](https://platform.claude.com/docs/en/release-notes/overview), **2026-09-14**, introduce on-demand Messages API compaction in beta (`compact-2026-09-04`). A signed summary block replaces older messages while recent turns can remain verbatim; callers choose timing and can compact in the background. New to this radar: evaluate preservation of task constraints and unresolved tool work across compaction.

2. **Claude Code — background completion and advisor accounting repaired.** [Changelog](https://code.claude.com/docs/en/changelog), **2.1.273, 2026-09-15**, fixes missing subagent results when final streams lack usage/model metadata, lost SDK background-subagent final reports, and advisor turns being counted at roughly twice their context size. It also adds opt-in gateway context/agent hints. **2.1.271, 2026-09-14**, adds command-scoped network domains and exact-command plugin acceptance. These are concrete state and permission regression cases; latest visible version is 2.1.273.

3. **MiniMax — streaming recovery and subtask lifecycle.** [CLI changelog](https://agent.minimax.io/docs/changelog), **0.4.8, 2026-09-16**, improves cancellation of compatible-provider streams and lets Explore subagents inspect Git/code through foreground Bash with read-only filesystem policy when sandboxing is available. **0.4.7, 2026-09-15**, fixes lost background Shell results, cleans up work when subtasks end, prevents replay after visible output, and preserves messages across abnormal endings. This is newer than the previously tracked Desktop release.

4. **Google — runtime intent checks complement deterministic controls.** [Engineering article](https://developers.googleblog.com/build-zero-trust-ai-agents-that-judge-intent-not-just-syntax/), **2026-09-15**, demonstrates Model Armor, semantic tool-call policies, and session anomaly detection on Gemini Enterprise Agent Platform. It separates platform policy ownership from agent code. Interpretation: evaluate valid-looking calls in their multi-turn context; syntax checks alone do not prove authorized business intent.

5. **OpenAI — managed Codex harness becomes an API surface.** [Agents API announcement](https://openai.com/index/introducing-the-agents-api/), **2026-09-10**, introduces a public beta for managed cloud agents with context management, tools, subagents, and execution environments. Newly captured in this run, not a September 16 launch. This raises a comparison question about lifecycle evidence and environment ownership; customer-reported gains are not Harness benchmarks.

## GitHub Trending Snapshot

[Global Today board](https://github.com/trending), checked **2026-09-16**. Most relevant: [open-code-review](https://github.com/alibaba/open-code-review) (hybrid review), [security-audit-skill](https://github.com/cloudflare/security-audit-skill) (verified findings), and [OpenResearch](https://github.com/alphaXiv/OpenResearch) (research agents). Their concrete workflow packaging may explain attention (inference). [Standalone board snapshot](../trending/2026-09-16.md) preserves non-AI entries too.

## Harness-Relevant Signals

- **Completion is an evidence chain:** Claude Code and MiniMax fixes distinguish final model output, result delivery, background cleanup, and user-visible completion. Candidate evaluation pressure increased; no runtime change follows automatically.
- **Context continuity needs behavioral checks:** Claude's signed compaction and corrected advisor accounting suggest checking retained constraints and pending work, rather than treating a smaller context meter as success.
- **Permission evaluation belongs in execution evidence:** Claude Managed Agents' **2026-09-10** `auto` policy evaluates tool calls as run, deny, or approval-pause, with evaluation fields in tool events; `ant beta:sessions connect` exposes live observation and approval. [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview). Google's runtime governance is complementary evidence, not proof that either policy system is infallible.

## Candidate Proposals

### Proposal 1 — Background result and cancellation continuity

- `source_type`: first-party Claude Code and MiniMax CLI release notes
- `problem_pattern`: final output exists but results are lost, replayed, or orphaned across backgrounding and cancellation
- `why_it_matters_for_harness`: orchestration must distinguish successful execution from delivered completion
- `recommended_disposition`: `encode-as-eval`

### Proposal 2 — Compaction preserves task invariants

- `source_type`: first-party Claude Platform and Claude Code release notes
- `problem_pattern`: context compression or incorrect accounting drops constraints or interrupts work prematurely
- `why_it_matters_for_harness`: long-running work needs continuity of scope, unresolved actions, and evidence
- `recommended_disposition`: `encode-as-eval`

### Proposal 3 — Managed harness and governance contract comparison

- `source_type`: first-party OpenAI API announcement and Google engineering guidance
- `problem_pattern`: managed execution hides boundaries between environment control, policy decisions, and agent outcomes
- `why_it_matters_for_harness`: adoption decisions need comparable observable lifecycle and ownership contracts
- `recommended_disposition`: `watch`

Intake candidates only; no direct Harness implementation recommendation.

## Watch Next

- **Anthropic mandatory checks completed:** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and [Claude Code changelog](https://code.claude.com/docs/en/changelog). Advisor guidance remains model-dependent and undated; no separately dated in-window advisor release established. Its Code accounting fix is explicitly surfaced above.
- Broader Anthropic coverage checked [Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), [Claude Blog](https://claude.com/blog/), and [legacy docs](https://docs.anthropic.com/en/docs/claude-code/overview), redirecting to Code. [September 15 Salesforce beta](https://support.claude.com/en/articles/12138966-release-notes) is a workflow-packaging watch, below the technical items above.
- **Moonshot:** [Kimi CLI 0.43.0](https://www.kimi.com/code/docs/en/kimi-code-cli/release-notes/changelog.html), **2026-09-14**, excludes closed-session time from Goal budgets, removes the 24-hour limit, permits steering to interrupt background waits, and bounds failing compaction attempts (default five). [K2.8 Preview](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html), **2026-09-11**, retains `kimi-for-coding`; track effective model/thinking configuration rather than alias alone.
- **DeepSeek:** [September 10 V4.1 Flash notice](https://api-docs.deepseek.com/updates/) still temporarily redirects legacy Flash aliases to `deepseek-flash` and explicitly continues V4 Pro after September 14 with unchanged billing. Prior retirement warning remains superseded; no live API calls performed.
- **Qwen:** one bounded search plus [official model lifecycle](https://help.aliyun.com/zh/model-studio/newly-released-models) found no strong dated Qwen technical release in-window. Latest visible Qwen snapshot is **2026-09-02**. The **2026-09-13** row is DeepSeek availability on Alibaba, not a Qwen model launch. Search absence is not proof of absence.
- No additional Langfuse, Phoenix, Helicone, or OpenTelemetry collection was needed. Trending stays secondary. Repository CURRENT is dated April 19; historical radar skill is absent. Explicit automation output scope governs; stale handoff and unrelated files remain untouched.
