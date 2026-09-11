# Daily Radar - 2026-09-10

Window: 2026-09-04 to 2026-09-10 (America/Los_Angeles). Checked 2026-09-10. One search pass per publisher family, followed by direct official pages. New versus yesterday: DeepSeek V4.1 Flash and alias migration, OpenAI's Astra work overview and GPT-Live-1 API, Claude Code 2.1.267, and MiniMax Desktop login repair. Vendor claims are not local validation.

## New Signals

1. **DeepSeek — V4.1 Flash changes model identity behind existing aliases.** [Official change log](https://api-docs.deepseek.com/updates/), **2026-09-10**, releases native multimodal V4.1 Flash as `deepseek-flash`. Retired V4 Flash and Vision Exp names temporarily route to the new model. `deepseek-v4-pro` requests will route to V4.1 Flash after **2026-09-14 12:00 Beijing time** (**2026-09-13 21:00 PDT**), until V4.1 Pro arrives. This is a concrete compatibility deadline: a stable request model name will no longer imply a stable backend. Benchmark improvements are vendor-reported; no local comparison was run.

2. **OpenAI — richer work agents and a separable full-duplex voice layer.** [Astra work overview](https://openai.com/index/gpt-6-astra-next-generation-work/), indexed **2026-09-09** on [OpenAI News](https://openai.com/news/), describes computer use, software work, fewer retries, and enterprise action controls. Its body says Astra was introduced “last week”; September 9 is this article's listing date, not independently established initial launch timing. [GPT-Live-1 API](https://openai.com/index/introducing-gpt-live-1-in-the-api/), **2026-09-10**, handles simultaneous listening/speaking while delegating reasoning and tool execution to a selected backend, including third-party models. Evaluate interruption behavior and backend task correctness separately; customer and benchmark claims need workload-specific confirmation.

3. **Anthropic — Claude Code adds effort ceilings and fixes policy-sensitive execution.** [Claude Code 2.1.267](https://code.claude.com/docs/en/changelog), **2026-09-09**, adds `maxEffortLevel` globally or per model across providers, plus `--system-prompt-snapshot off` for prompt iteration. Fixes cover cloud scheduled-task startup when organization settings require sandboxing, and large-schema Workflow `agent()` calls incorrectly refused in auto mode. Retained **September 8** signal: 2.1.265 makes advisor availability/instructions stable rather than recomputed per request, alongside recovery fixes; 2.1.266 repairs gateway authentication. These are distinct versioned behavior contracts.

4. **Google — behavioral evals explain model-upgrade regressions.** [The Anatomy of Harness Engineering](https://developers.googleblog.com/the-anatomy-of-harness-engineering-how-to-evaluate-iterate-and-guard-ai-coding-agents/), **2026-09-09**, remains this week's strongest evaluation-method signal: assert observable intermediate actions, complement end-to-end outcomes, permit alternative valid paths on complex tasks, and assess aggregate pass rates rather than noisy individual runs. Retained from yesterday, not a new September 10 release.

5. **MiniMax — authentication and execution continuity remain coupled.** [Official changelog](https://agent.minimax.io/docs/changelog): Desktop **3.0.70, 2026-09-10**, reports login fixes without a detailed root cause; CLI **0.3.11, 2026-09-09**, specifically repairs long-task 401 failures after repeated credential renewals. Desktop **3.0.69, September 8**, adds local/cloud session visibility and result download. Treat these as separate product/version changes; the generic Desktop fix does not prove it shares the CLI bug.

## GitHub Trending Snapshot

The [default global Today board](https://github.com/trending), captured **2026-09-10**, includes [teamai-cli](https://github.com/Tencent/teamai-cli), [llmfit](https://github.com/AlexsJones/llmfit), [OmniRoute](https://github.com/diegosouzapw/OmniRoute), and [PI-Desktop](https://github.com/vastsa/PI-Desktop): team workflows, hardware fit, routing, and agent environments. These are secondary attention signals. See the [standalone 16-project snapshot](../trending/2026-09-10.md).

## Harness-Relevant Signals

- **Model identity and policy are evaluation inputs:** DeepSeek alias changes and Claude Code effort caps can alter behavior without a task prompt changing.
- **Lifecycle boundaries need separate evidence:** credential renewal, sandbox startup, interruption handling, and completed backend work are different contracts.
- **Behavior plus outcome:** Google's method provides a useful way to test those contracts, but a candidate needs an observed local failure before implementation.

## Candidate Proposals

### Proposal 1 — Model alias drift

- `source_type`: first-party DeepSeek release notes
- `problem_pattern`: unchanged model identifiers begin serving a different model
- `why_it_matters_for_harness`: comparisons and recovery assumptions can become stale across the announced cutover
- `recommended_disposition`: `encode-as-eval`

### Proposal 2 — Policy-sensitive execution continuity

- `source_type`: first-party Claude Code and MiniMax changelogs; Google engineering guidance
- `problem_pattern`: effort caps, sandbox settings, or credential renewal change whether long tasks can proceed
- `why_it_matters_for_harness`: outcome scores alone cannot locate these operational failures
- `recommended_disposition`: `encode-as-eval`

### Proposal 3 — Voice and backend coordination

- `source_type`: first-party OpenAI API announcement
- `problem_pattern`: conversation interruption and asynchronous tool completion can diverge
- `why_it_matters_for_harness`: a responsive voice interface does not establish correct task execution
- `recommended_disposition`: `watch`

Intake candidates only; no direct Harness implementation recommendation or runtime change.

## Watch Next

- **Anthropic required coverage completed:** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and [Code changelog](https://code.claude.com/docs/en/changelog). Platform's latest dated entry is **September 3**, now outside the window. Managed Agents and advisor documentation remain useful undated capability context; no separately dated in-window Platform/Managed Agents/advisor API release established. The in-window advisor change belongs to Claude Code.
- Also checked [Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), [Claude Blog](https://claude.com/blog/), and [legacy Code docs](https://docs.anthropic.com/en/docs/claude-code/overview), which redirect to code.claude.com. No stronger dated in-window technical item established there.
- **Moonshot:** [Kimi Work 3.2.6](https://www.kimi.com/en/help/kimi-work/release-notes), **2026-09-07**, retains GitHub plugin installation, visible MCP connection management, and confirmation of running scheduled tasks at exit. Rechecked weekly context; no newer strong dated item established in the bounded pass.
- **Qwen:** [Model Studio lifecycle](https://help.aliyun.com/en/model-studio/newly-released-models) newest Qwen snapshot remains **2026-09-02**, outside the window. No strong dated official technical item found in-window; this is bounded coverage, not proof of no update.
- DeepSeek's September 14 Beijing alias cutover is the nearest concrete compatibility watch. Langfuse, Phoenix, Helicone, and OpenTelemetry were not needed to clarify today's primary signals.
- Repository freshness: `_project/CURRENT.md` remains dated **2026-04-19**; the historical daily-radar skill is absent from the repository and searched installed skill paths. Followed current Radar rules and this run's explicit seven-file write boundary, leaving stale handoff and unrelated files untouched.
