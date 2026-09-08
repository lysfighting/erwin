# Daily Radar - 2026-09-08

Window: 2026-09-02 to 2026-09-08 (America/Los_Angeles). Checked 2026-09-08. One search pass per publisher family, then direct official pages. New since the September 7 note: MiniMax 0.3.5–0.3.10 and Kimi Work 3.2.6. Other signals below are retained weekly context, not new launches today. Vendor claims are not local Harness validation.

## New Signals

1. **MiniMax — side-session authority changed; interrupted queues now have explicit recovery.** [Agent changelog](https://agent.minimax.io/docs/changelog), **0.3.9, 2026-09-07**, makes `/btw` and `/side` inherit parent permissions: Full access needs no extra approval, while other modes retain their checks. This supersedes the September 4 rule repeated in yesterday's radar. **0.3.5, 2026-09-07**, preserves unexecuted queued messages as paused after failure or interruption, including on reopening; continuation is explicit. **0.3.6, 2026-09-07**, avoids persisting automatic endpoint resolution into shared Desktop configuration. **0.3.10, 2026-09-08**, is menu presentation polish, not a new control-plane capability.

2. **Moonshot — plugin installation and running-work visibility become product surfaces.** [Kimi Work 3.2.6](https://www.kimi.com/en/help/kimi-work/release-notes), **2026-09-07**, adds task-related plugin recommendations, GitHub-link installation, MCP connection status/management, and confirmation of still-running scheduled tasks before app exit. Newly captured today. Separately, [Kimi Code 0.41.0](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html), **2026-09-04**, retains experimental Tower collaboration and direct background-answer delivery; permission behavior changed from 0.40.0. These are distinct Work and Code release streams.

3. **Anthropic — reproducible resources and truthful recovery remain the week's strongest signal.** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), **2026-09-03**, introduce `ant apply` in CLI 1.30.0: repository definitions, an approval plan, and `claude-lock.json` keep later applications attached to the same resources. [Claude Code 2.1.261](https://code.claude.com/docs/en/changelog), **2026-09-04**, fixes resumed parallel-tool context, stale remote permissions, early interrupts, and offline messages falsely appearing delivered; it adds unused-skill context diagnostics. **2.1.263, 2026-09-06**, has generic reliability notes only. Retained weekly evidence; no newer substantive dated entry established.

4. **OpenAI — research-agent activity needs outcome and supervision accounting.** [Research acceleration: The view inside OpenAI](https://openai.com/index/research-acceleration-view-inside-openai/), **2026-09-06**, describes progress toward an automated research intern under human supervision. Retained weekly signal. Treat reported internal activity and capability as observational/vendor evidence; it does not establish unsupervised research success or local productivity gains.

5. **Google — stronger agentic coding remains a budgeted evaluation question.** [Gemini 3.8 Flash and Flash Cyber](https://blog.google/innovation-and-ai/models-and-research/gemini-models/3-8-flash-and-3-8-flash-cyber/), **2026-09-02**, introduces coding/reasoning improvements and a cybersecurity variant. Retained weekly launch; compare model quality under fixed resources before drawing a Harness adoption conclusion. The September 1 video announcement has aged out of this run's calendar window.

## GitHub Trending Snapshot

Secondary evidence from the [default global Today board](https://github.com/trending), captured **2026-09-08**: [openai/skills](https://github.com/openai/skills), [ECC](https://github.com/affaan-m/ECC), [context-mode](https://github.com/mksglu/context-mode), and [markitdown](https://github.com/microsoft/markitdown) surface reusable workflows, context control, and document intake. Relevance is interpretation, not a measured cause of popularity. See the [standalone board snapshot](../trending/2026-09-08.md).

## Harness-Relevant Signals

- **Authority must be versioned:** MiniMax's changed side-session rule and Kimi Code's mode changes make mode names alone insufficient evidence of allowed actions.
- **Completion is a state contract:** compare queue persistence, offline delivery, interruption, and app-exit behavior before judging a task finished. Kimi Work and Claude Code expose different parts of this boundary.
- **Resource identity matters:** Anthropic's declarative apply flow is relevant to repeatable agent environments, with local proof still required.
- **Qwen coverage:** [Model Studio updates](https://help.aliyun.com/en/model-studio/newly-released-models) lists **qwen3.8-max-0902, 2026-09-02**, with vendor claims of stronger long-horizon engineering and agent collaboration. Retained watch signal; no stronger additional dated item established in the bounded pass. September 1 retrieval releases are outside today's window.

## Candidate Proposals

### Proposal 1 — Versioned side-session authority

- `source_type`: first-party MiniMax and Kimi Code release notes
- `problem_pattern`: a familiar permission label hides changed behavior after an upgrade
- `why_it_matters_for_harness`: authority checks need the actual version and inherited mode
- `recommended_disposition`: `encode-as-eval`

### Proposal 2 — Interrupted work remains recoverable

- `source_type`: first-party MiniMax, Claude Code, and Kimi Work changelogs
- `problem_pattern`: pending, queued, running, and delivered states collapse into apparent completion
- `why_it_matters_for_harness`: recovery should preserve user intent and expose unfinished work
- `recommended_disposition`: `encode-as-eval`

### Proposal 3 — Capability availability evidence

- `source_type`: first-party Kimi Work release notes; GitHub Trending as secondary context
- `problem_pattern`: an installed plugin is assumed connected and ready
- `why_it_matters_for_harness`: installation, connection, and successful capability use need separate evidence
- `recommended_disposition`: `watch`

Intake candidates only; no direct Harness implementation recommendation or reviewed adoption decision.

## Watch Next

- **Required Anthropic checks completed:** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code changelog](https://code.claude.com/docs/en/changelog), and its [GitHub source](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md). Managed Agents retains the September 3 resource-management signal. Advisor docs describe executor/advisor consultation, but no separately dated in-window update was established; undated documentation is not a launch.
- Also checked [Anthropic Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), [Claude Blog](https://claude.com/blog/), and [legacy Code docs](https://docs.anthropic.com/en/docs/claude-code/overview), which redirects to code.claude.com. No stronger additional dated technical item was established on these surfaces in the single pass.
- **DeepSeek:** [official change log](https://api-docs.deepseek.com/updates/) still shows **2026-08-21** as its newest visible entry. No strong dated technical item found in-window; this is bounded coverage, not proof that no change occurred.
- Watch permission semantics and queue recovery across subsequent releases. Langfuse, Phoenix, Helicone, and OpenTelemetry were not needed to clarify today's primary signals; no secondary tooling release is promoted.
- Repository freshness: `_project/CURRENT.md` remains dated **2026-04-19**. Followed the explicit automation scope and current radar rules; left that stale handoff untouched.
