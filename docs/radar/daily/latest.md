# Daily Radar - 2026-09-11

Window: 2026-09-05 to 2026-09-11 (America/Los_Angeles). Checked 2026-09-11. One search pass per publisher family, followed by direct official pages. Material changes versus the September 10 note: Managed Agents auto permissions, Claude Code 2.1.268, Kimi Work 3.2.7, and a corrected DeepSeek V4 Pro continuity notice. Vendor claims are not local validation.

## New Signals

1. **Anthropic Platform — Managed Agents gains per-call policy evaluation and terminal attachment.** [Release notes](https://platform.claude.com/docs/en/release-notes/overview), **2026-09-10**, introduce `auto` permissions and `ant beta:sessions connect` for live observation, messages, and approval handling; `--web` serves a local session viewer. [Permission documentation](https://platform.claude.com/docs/en/managed-agents/permission-policies) distinguishes allow, deny, and approval pause: a denied call does not terminate the session and cannot be overridden by the client. `auto` is opt-in, not a required human checkpoint. Tool events expose policy evaluation; custom tools are outside this policy mechanism. This creates specific audit and continuation contracts worth evaluating.

2. **Claude Code — 2.1.268 repairs permission and endpoint compatibility edges.** [Official changelog](https://code.claude.com/docs/en/changelog), **2026-09-10**, fixes deny/ask rules for symlinked paths, untrusted same-name agent configuration on teammate respawn, and HTTP 400 failures on third-party Anthropic-compatible endpoints introduced in 2.1.265. It also bounds hanging WebFetch calls at 300 seconds by default, propagates gateway pricing to client cost telemetry, and adds opt-in cleanup of self-hosted runner session directories. These changes matter more than generic coding-performance claims because they affect execution boundaries and observability.

3. **DeepSeek — V4 Pro continuation supersedes yesterday's cutover warning.** The [official change log](https://api-docs.deepseek.com/updates/), under its **2026-09-10** heading, now explicitly says V4 Pro API service will continue after **2026-09-14**, with billing unchanged. **The September 10 radar's V4 Pro-to-Flash deadline is stale.** The correction's own edit timestamp is not published; this wording was directly retrieved on September 11 after the browser fetch became unreliable. V4.1 Flash remains the new multimodal model behind `deepseek-flash`; retired V4 Flash and Vision Exp aliases temporarily route to it. Preserve the distinction between changed Flash aliases and continuing Pro service.

4. **Moonshot — Kimi Work 3.2.7 makes task context and remote progress more explicit.** [Release notes](https://www.kimi.com/en/help/kimi-work/release-notes), **2026-09-11**, automatically carry conversation settings into new Work conversations, add a website deployment plugin, align remote-control subagent/deliverable cards with desktop, and fix the scheduled-task page freezing when credits run out. September 7's plugin installation and MCP connection management remain weekly context. Useful pressure: settings propagation and exhausted-credit handling need observable states, not just successful task outputs.

5. **OpenAI — GPT-Live-1 separates live conversation from backend work.** [API announcement](https://openai.com/index/introducing-gpt-live-1-in-the-api/), **2026-09-10**, combines listening and speaking while delegating reasoning and tool calls to a selected text backend, including third-party models. Retained weekly signal: conversational interruption and backend task completion require separate evidence. Date correction: [official release notes](https://openai.com/products/release-notes/) establish Astra's initial rollout as **2026-09-03**, outside this window; do not treat a later work-overview listing as the model launch.

## GitHub Trending Snapshot

The [default global Today board](https://github.com/trending), checked **2026-09-11**, surfaces [PI-Desktop](https://github.com/vastsa/PI-Desktop) for agent environments, [llm_wiki](https://github.com/nashsu/llm_wiki) and [hyperresearch](https://github.com/jordan-gibbs/hyperresearch) for persistent knowledge, and [spec-kit](https://github.com/github/spec-kit) for structured development. Their practical workflow focus may explain attention; causality is inferred. The [standalone snapshot](../trending/2026-09-11.md) preserves all 16 entries, including non-AI interests.

## Harness-Relevant Signals

- **Permission outcomes are distinct states:** a tool denial, an approval pause, and a completed session require different evidence. Managed Agents provides an explicit reference contract.
- **Compatibility includes environment and identity:** Code path normalization, gateway rates, inherited settings, and DeepSeek alias policies can change behavior without changing the task prompt.
- **Evaluate behaviors alongside outcomes:** Google's [The Anatomy of Harness Engineering](https://developers.googleblog.com/the-anatomy-of-harness-engineering-how-to-evaluate-iterate-and-guard-ai-coding-agents/), **2026-09-09**, argues for discrete observable behavior checks that help explain aggregate benchmark changes. Retained method signal; no local evaluation was run.

## Candidate Proposals

### Proposal 1 — Permission outcome and continuation evidence

- `source_type`: first-party Managed Agents release notes and permission documentation
- `problem_pattern`: tool denial, human approval pause, and session completion are conflated
- `why_it_matters_for_harness`: automation can misreport blocked work as completed or treat a recoverable denial as terminal
- `recommended_disposition`: `encode-as-eval`

### Proposal 2 — Execution environment compatibility

- `source_type`: first-party Claude Code changelog and Kimi Work release notes
- `problem_pattern`: path spelling, provider schemas, or inherited settings alter execution semantics
- `why_it_matters_for_harness`: identical prompts do not establish equivalent permissions or configuration
- `recommended_disposition`: `encode-as-eval`

### Proposal 3 — Model lifecycle notice revisions

- `source_type`: first-party DeepSeek change log
- `problem_pattern`: a vendor revises migration policy under an existing release date
- `why_it_matters_for_harness`: stale migration assumptions contaminate comparisons and operational planning
- `recommended_disposition`: `watch`

Intake candidates only; no direct Harness implementation recommendation or runtime change.

## Watch Next

- **Anthropic required coverage completed:** [Platform notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and [Code changelog](https://code.claude.com/docs/en/changelog). Advisor docs remain capability context without a separately established in-window release date. Code 2.1.265 (**September 8**) stabilized advisor availability/instructions; that is Code-specific, not a new advisor API launch.
- Also checked [Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), [Claude Blog](https://claude.com/blog/), and [legacy docs](https://docs.anthropic.com/en/docs/claude-code/overview), which redirect to Code docs. No stronger dated technical item established by this bounded pass.
- **MiniMax:** [changelog](https://agent.minimax.io/docs/changelog), Desktop **3.0.71, 2026-09-11**, updates the macOS signing certificate for install/launch compatibility; Windows still shows 3.0.70. CLI **0.3.11, September 9**, repairs long-task 401s after repeated credential renewal. Distinct release mechanisms; no shared root cause inferred.
- **Qwen:** [Model Studio lifecycle](https://help.aliyun.com/en/model-studio/newly-released-models) shows the latest visible Qwen snapshot on **2026-09-02**, outside the window. No strong dated official technical item found in-window in the bounded pass; this does not prove absence.
- Watch revised DeepSeek lifecycle notices, permission-evaluation evidence, and voice/backend completion behavior. Langfuse, Phoenix, Helicone, and OpenTelemetry were unnecessary to clarify today's primary signals.
- Repository freshness: `_project/CURRENT.md` remains **2026-04-19**; historical daily-radar skill is absent from repo and searched installed paths. Applied current Radar rules and the explicit seven-file output boundary; unrelated files and stale handoff were left untouched.
