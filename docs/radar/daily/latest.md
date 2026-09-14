# Daily Radar - 2026-09-14

Window: 2026-09-08 to 2026-09-14 (America/Los_Angeles). Checked 2026-09-14. One search pass per publisher family, followed by direct official pages. Compared with September 13, the confirmed new intake is MiniMax Desktop 3.0.72; other selected releases remain relevant weekly signals. Vendor claims are not local validation.

## New Signals

1. **MiniMax — local/cloud handoff becomes an explicit workflow.** [Desktop 3.0.72 changelog](https://agent.minimax.io/docs/changelog), **2026-09-14**, adds local task transfer to cloud execution, file review before upload, and download or compare/sync of cloud outputs. It also adds independent chats with configured subagents and HTML canvas annotations. New since yesterday: watch how task ownership, source revision, and output reconciliation survive the transfer; the announcement does not establish conflict-safe behavior.

2. **Claude Code — reproducible plugin evaluations and state correctness.** [Official changelog](https://code.claude.com/docs/en/changelog), **2.1.269, 2026-09-11**, adds `claude plugin eval` with JSON/HTML reports and optional repository attributes in OpenTelemetry. It fixes background-running sessions incorrectly reporting input-wait. **2.1.270, 2026-09-12**, repairs read-only git commands unexpectedly requesting permission in long sessions. Latest visible version remains 2.1.270; retained weekly signal.

3. **Claude Platform — permission decisions carry execution evidence.** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), **2026-09-10**, add Managed Agents `auto` permission policies: each agent/MCP tool call can run, be denied, or pause for approval, with evaluation evidence in tool events. `ant beta:sessions connect` enables terminal observation, messages, and approvals. Retained weekly signal; no newer dated Platform release established.

4. **Moonshot — stable aliases can conceal effective model changes.** [Kimi Code release notes](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html), **2026-09-11**, roll out K2.8 Preview under unchanged `kimi-for-coding`, with up to 1M context and `low`/`high`/`max` effort (default `max`). Disabling thinking routes K3 and K2.8 Preview requests to non-thinking K2.8 Preview. Reproducible comparisons need effective model and thinking configuration, not just the alias. Retained weekly signal.

5. **DeepSeek — Pro continues past today's former deadline.** [Official change log](https://api-docs.deepseek.com/updates/), **2026-09-10**, releases multimodal V4.1 Flash as `deepseek-flash`; old V4 Flash and Vision Exp aliases temporarily redirect. The current notice explicitly preserves V4 Pro API service after **2026-09-14**, with billing unchanged. This supersedes the September 10 radar's earlier Pro cutover warning. Reverified today; notice-edit timestamp unavailable. No live API invocation was performed.

## GitHub Trending Snapshot

[Default global Today board](https://github.com/trending), checked **2026-09-14**: [open-code-review](https://github.com/alibaba/open-code-review), [oh-my-hermes](https://github.com/rlaope/oh-my-hermes), and [agent-skills](https://github.com/tech-leads-club/agent-skills) surface hybrid review, memory, and reusable workflows. Packaging may explain attention (inference). [Full snapshot](../trending/2026-09-14.md).

## Harness-Relevant Signals

- **Handoff evidence:** local/cloud transfer needs distinguishable input, execution owner, returned output, and reconciliation result. This is an evaluation question raised by MiniMax, not permission to change Harness.
- **Behavior over aggregate scores:** Google's [The Anatomy of Harness Engineering](https://developers.googleblog.com/the-anatomy-of-harness-engineering-how-to-evaluate-iterate-and-guard-ai-coding-agents/), **2026-09-09**, argues for checking discrete intermediate actions to diagnose regressions obscured by end-to-end scores. Complements Claude Code's plugin evaluation surface.
- **Routing provenance:** Kimi and DeepSeek show why alias continuity and behavior continuity must be assessed separately.

## Candidate Proposals

### Proposal 1 — Handoff reconciliation evidence

- `source_type`: first-party MiniMax release notes
- `problem_pattern`: local/cloud continuation obscures source revision or returned-file conflicts
- `why_it_matters_for_harness`: completion claims require traceable inputs and reconciled outputs
- `recommended_disposition`: `watch`

### Proposal 2 — Execution-state behavior checks

- `source_type`: first-party Claude Code, Claude Platform, and Google engineering guidance
- `problem_pattern`: background work or approval pauses are mistaken for idle or completed work
- `why_it_matters_for_harness`: intermediate state determines correct orchestration and intervention
- `recommended_disposition`: `encode-as-eval`

### Proposal 3 — Effective model identity

- `source_type`: first-party Kimi routing and DeepSeek lifecycle documentation
- `problem_pattern`: stable model names mask routing or lifecycle changes
- `why_it_matters_for_harness`: benchmark comparisons require effective configuration provenance
- `recommended_disposition`: `watch`

Intake candidates only; no direct Harness implementation recommendation or runtime change.

## Watch Next

- **Anthropic mandatory checks:** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and [Code changelog](https://code.claude.com/docs/en/changelog) all checked. Advisor docs describe model-dependent consultation nudges, but no separately dated in-window advisor release was established; do not treat undated guidance as a new launch.
- Anthropic broader surface also checked: [Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), [Claude Blog](https://claude.com/blog/), and [legacy docs](https://docs.anthropic.com/en/docs/claude-code/overview), which redirect to Code. No stronger dated technical item established in this bounded pass.
- **OpenAI:** [GPT-Live-1 API announcement](https://openai.com/index/introducing-gpt-live-1-in-the-api/), **2026-09-10**, combines simultaneous listening/speaking with backend reasoning/tool delegation, including third-party text models. Retained weekly watch: evaluate interruptions separately from background task completion.
- **Qwen:** [official model lifecycle](https://help.aliyun.com/zh/model-studio/newly-released-models) has a latest visible snapshot dated **2026-09-02**, outside the window. [Alibaba Managed Agent changelog](https://help.aliyun.com/en/model-studio/managed-agents-changelog) adds multimodal user events on **2026-09-07**, also outside this run's window. No strong dated official technical release found in-window; search absence is not proof of absence.
- No additional Langfuse, Phoenix, Helicone, or OpenTelemetry search was needed; secondary signals remain subordinate to first-party evidence.
- Repository freshness: `_project/CURRENT.md` remains dated **2026-04-19**; historical daily-radar skill was absent from repo and searched installed paths. Current Radar docs and the explicit seven-file output scope govern this run. Unrelated files and stale handoff remain untouched.
