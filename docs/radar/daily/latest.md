# Daily Radar - 2026-09-09

Window: 2026-09-03 to 2026-09-09 (America/Los_Angeles). Checked 2026-09-09. One search pass per publisher family, followed by direct official pages. New since yesterday's capture: Google behavioral-eval guidance and Kotlin ADK 1.0, Claude Code 2.1.265–266, MiniMax 0.3.11, and OpenAI's September 8 research cases. Vendor reports are not local Harness validation.

## New Signals

1. **Google — evaluate observable behavior alongside end-to-end outcomes.** [The Anatomy of Harness Engineering](https://developers.googleblog.com/the-anatomy-of-harness-engineering-how-to-evaluate-iterate-and-guard-ai-coding-agents/), **2026-09-09**, proposes small checks on tool use and intermediate actions to diagnose regressions. It explicitly retains macro benchmarks, avoids rigid action sequences for complex tasks, and recommends aggregate evaluation over noisy single runs. Separately, [ADK for Kotlin 1.0](https://developers.googleblog.com/announcing-adk-for-kotlin-10-building-production-ready-ai-agents-in-kotlin-android-and-beyond/), **2026-09-09**, announces core parity, resumable sessions, confirmation flows, and Android persistence extensions. The useful common theme is testable execution state, not an automatic framework migration.

2. **Anthropic — Claude Code recovery, advisor decisions, and gateway auth receive concrete fixes.** [Claude Code changelog](https://code.claude.com/docs/en/changelog), **2.1.265 and 2.1.266, 2026-09-08**: interrupted tool calls survive process-death recovery with truthful status; resumed subagent prompt prefixes remain cache-compatible; remote replies no longer appear complete before their last message. Advisor availability/instructions are decided once and changes announced, instead of recalculated per request. Version 2.1.266 repairs the preceding version's gateway sign-in regression for certain custom-auth configurations. Retained Platform signal: [ant CLI 1.30.0](https://platform.claude.com/docs/en/release-notes/overview), **2026-09-03**, adds repository-defined resources, plan approval, and a resource-identity lockfile through `ant apply`.

3. **MiniMax — long-task continuity includes repeated credential renewal.** [Agent changelog](https://agent.minimax.io/docs/changelog), **0.3.11, 2026-09-09**, fixes long-running tasks ending with HTTP 401 after multiple sign-in credential renewals. Retained context: **0.3.9, 2026-09-07**, makes side sessions inherit parent permissions; **0.3.5, 2026-09-07**, retains interrupted queued work as paused for explicit continuation. Authentication lifetime belongs in the same recovery assessment as process and queue lifetime.

4. **OpenAI — research orchestration needs budget and verification accounting.** [Navier–Stokes research report](https://openai.com/index/navier-stokes-solution/), **2026-09-08**, reports parallel agent groups, cross-group synthesis, a solution reached on **September 5**, and a separate Lean formalization/verification stage. The reported Navier–Stokes effort used roughly 130 billion output tokens; this is not a small-budget productivity benchmark. The mathematical result and verification are OpenAI's report; this radar did not independently review the proof. A separate [quantum-computing workflow case](https://openai.com/index/codex-quantum-computing-experiments/), **2026-09-08**, describes agents assisting routine laboratory measurements while the researcher retains experimental design and interpretation. Both sharpen outcome/supervision questions beyond counting agent activity.

5. **Moonshot — plugin readiness and unfinished work become visible.** [Kimi Work 3.2.6](https://www.kimi.com/en/help/kimi-work/release-notes), **2026-09-07**, adds GitHub plugin installation, MCP connection management/status, and exit confirmation for still-running scheduled tasks. Retained weekly context, not a new release today. [Kimi Code 0.41.0](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html), **2026-09-04**, adds experimental Tower collaboration and changes Auto-mode command handling; permission semantics require version-specific evidence.

## GitHub Trending Snapshot

Secondary signal from the [default global Today board](https://github.com/trending), captured **2026-09-09**: [teamai-cli](https://github.com/Tencent/teamai-cli), [PI-Desktop](https://github.com/vastsa/PI-Desktop), [openai/plugins](https://github.com/openai/plugins), and [ECC](https://github.com/affaan-m/ECC) highlight team workflows, extensibility, and agent environments. These are relevance interpretations, not established causes of popularity. The [standalone snapshot](../trending/2026-09-09.md) preserves all 13 visible entries, including architecture and learning projects.

## Harness-Relevant Signals

- **Measure the failure boundary:** Google behavioral checks complement Claude Code's concrete interruption, completion-order, and advisor regressions. Tie any candidate to an observed local failure before implementation.
- **Separate operational lifetimes:** a running process, valid credentials, queued intent, and complete delivery are different states. MiniMax provides an explicit long-duration authentication failure pattern.
- **Account for the whole research workflow:** record synthesis, verification, elapsed time, resource use, and human contribution before comparing multi-agent outcomes.
- **Capability readiness needs evidence:** Kimi Work connection visibility and Trending plugin interest support checking installation, connection, and successful invocation separately.

## Candidate Proposals

### Proposal 1 — Recovery and completion behavior

- `source_type`: first-party Google engineering guidance and Claude Code / MiniMax changelogs
- `problem_pattern`: interruption, credential renewal, or premature completion loses or misrepresents unfinished work
- `why_it_matters_for_harness`: observable state transitions can expose regressions hidden by final-answer scores
- `recommended_disposition`: `encode-as-eval`

### Proposal 2 — Stable advisor configuration

- `source_type`: first-party Claude Code release notes
- `problem_pattern`: per-request model changes silently alter advisor availability or instructions
- `why_it_matters_for_harness`: delegation expectations need explicit transitions and reproducible context
- `recommended_disposition`: `watch`

### Proposal 3 — Research outcome accounting

- `source_type`: first-party OpenAI research report and applied case study
- `problem_pattern`: agent volume or a selected success hides total budget, verification work, and supervision
- `why_it_matters_for_harness`: comparisons need a complete effort denominator and independently checkable artifacts
- `recommended_disposition`: `watch`

These are intake candidates, not reviewed implementation decisions or direct Harness change recommendations.

## Watch Next

- **Required Anthropic checks completed:** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code changelog](https://code.claude.com/docs/en/changelog), and its [GitHub source](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md). Managed Agents retains the September 3 resource-management signal. Advisor documentation is undated; the newly dated advisor signal is specifically the Claude Code fix above, not a newly established Platform API release.
- Also checked [Anthropic Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), [Claude Blog](https://claude.com/blog/), and [legacy docs](https://docs.anthropic.com/en/docs/claude-code/overview), which redirect to code.claude.com. No stronger additional dated technical item established on those surfaces.
- **Qwen:** [Model Studio lifecycle](https://help.aliyun.com/en/model-studio/newly-released-models) still lists **qwen3.8-max-0902, 2026-09-02**, outside today's window. No strong dated official technical item found in-window in this bounded pass.
- **DeepSeek:** [official change log](https://api-docs.deepseek.com/updates/) newest visible entry remains **2026-08-21**, outside the window. No strong dated technical item found in-window; bounded coverage is not proof of no change.
- September 2 model announcements have aged out. Google developer guidance now supplies fresher Harness evidence. Langfuse, Phoenix, Helicone, and OpenTelemetry were not needed to clarify this run's primary signals.
- Repository freshness: `_project/CURRENT.md` is still dated **2026-04-19**; the former daily radar skill is absent. Applied current radar rules and the explicit automation write boundary, leaving the stale handoff and unrelated files untouched.
