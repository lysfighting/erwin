# Daily Radar - 2026-07-05

Window: 2026-06-28 to 2026-07-05

## New Signals

1. **Anthropic paired a new Claude model launch with a material Managed Agents control-plane upgrade.** The official Anthropic newsroom published [Introducing Claude Sonnet 5](https://www.anthropic.com/news/claude-sonnet-5) on **2026-06-30**, and the official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) on **2026-06-30** added Managed Agents **event deltas**, **backward session pagination**, **per-session agent overrides**, vault **`injection_location`** controls, and broader **agent/deployment/deployment-run webhooks**. This is the strongest in-window first-party signal: better execution quality plus more operator surface for long-running agents.  
   Publish date: 2026-06-30

2. **Anthropic's Claude Code surface is still moving even though the weekly digest lags.** The official [Claude Code What's New](https://code.claude.com/docs/en/whats-new) page still tops out at **Week 26 (2026-06-22 to 2026-06-26)**, but the official raw [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md) now shows post-launch entries through `2.1.201`, including **background subagents by default**, **Claude Platform on AWS as an upstream provider**, improved worktree/plugin handling, and follow-on fixes after the Sonnet 5 rollout. This is an inference from the current raw changelog state rather than a separately dated changelog post, so treat it as a live product-surface signal, not a precisely dated announcement.

3. **Moonshot/Kimi compressed several coding-agent ergonomics into one tight release train.** The official [Kimi Code What's New](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html) shows **v0.22.0 on 2026-07-02** adding model configuration overrides and a redesigned web UI, **v0.21.0 on 2026-07-01** adding plugin-declared slash commands and Mermaid rendering, and **v0.20.2 on 2026-06-29** adding an Anthropic-compatible protocol, video input, and custom outbound request headers. The pattern is clear: provider interoperability, reusable command surfaces, and richer operator UIs are becoming baseline expectations.  
   Publish dates: 2026-07-02, 2026-07-01, 2026-06-29

4. **Google kept extending agent-capable multimodal execution surfaces.** The official [Gemini API release notes](https://ai.google.dev/gemini-api/docs/changelog) added **Gemini Omni Flash** in public preview on **2026-06-30** for high-speed video generation and conversational video editing through the Interactions API, and added **Computer Use** public preview support for Gemini 3.5 Flash on **2026-06-24**. It is not as coding-centric as Anthropic or Kimi this week, but it reinforces the trend toward managed action surfaces rather than plain text endpoints.  
   Publish dates: 2026-06-30 and 2026-06-24

## GitHub Trending Snapshot

Snapshot source: [GitHub Trending](https://github.com/trending) captured on 2026-07-05.

- [openai/codex-plugin-cc](https://github.com/openai/codex-plugin-cc): surfacing because cross-agent delegation and review workflows are becoming a real developer-tooling category, not just a demo pattern.
- [usestrix/strix](https://github.com/usestrix/strix): relevant because security remains one of the clearest "agents do real work" narratives.
- [alibaba/page-agent](https://github.com/alibaba/page-agent): notable because in-page GUI control is still a strong direction for browser-native agent execution.
- [ogulcancelik/herdr](https://github.com/ogulcancelik/herdr): worth watching because terminal-resident agent multiplexing lines up directly with operator-control questions.
- [dotnet/skills](https://github.com/dotnet/skills): relevant because curated skill packs for coding agents are maturing into a first-class ecosystem surface.

## Harness-Relevant Signals

- Anthropic was explicitly checked across the required official surfaces this run: [Anthropic newsroom](https://www.anthropic.com/news), [Claude blog/docs hub](https://claude.com), [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code What's New](https://code.claude.com/docs/en/whats-new), and the raw [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md).
- The [Claude Managed Agents quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart) still requires the `managed-agents-2026-04-01` beta header and keeps the canonical `agent -> environment -> session -> events` model front and center. I did not find a separately dated in-window docs update there, but the page remains strategically important because the June 30 release notes expanded the surrounding session-control surface.
- The [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) were also checked explicitly. The page still frames the `advisor-tool-2026-03-01` beta as a planner/executor split for coding, computer-use, and research workloads. I did not find a newer dated in-window change on that page.
- The strongest weekly pattern is **operator surface hardening for long-running agents**: Anthropic expanded session streaming, overrides, and lifecycle hooks; Kimi kept turning plugins and protocol compatibility into default ergonomics; Google kept pushing managed multimodal and computer-use surfaces.
- OpenAI was checked through the official [Codex changelog](https://developers.openai.com/codex/changelog). I did **not** find a strong new dated official OpenAI technical item inside the 2026-06-28 to 2026-07-05 window; the latest clearly relevant Codex item I found was **Codex Remote GA on 2026-06-25**, which is outside this run window.
- MiniMax, Qwen, and DeepSeek were each checked once through official surfaces. I did **not** find a strong dated official technical item inside this run window. The latest dated anchors I found remained [MiniMax M3](https://platform.minimax.io/docs/release-notes/models) on **2026-06-01**, DeepSeek's [API change log](https://api-docs.deepseek.com/updates) on **2026-04-24**, and Qwen's official [blog index](https://qwen.ai/blog/) surfacing newer posts but without a clearer exact in-window publish date I could confidently anchor.

## Candidate Proposals

### Proposal 1

- `source_type`: first-party platform release notes
- `problem_pattern`: managed agent systems now need stream previews, session steering, lifecycle hooks, and per-session override controls rather than opaque run/finish behavior
- `why_it_matters_for_harness`: the operator layer is becoming part of the product surface, not optional observability garnish
- `recommended_disposition`: `encode-as-eval`

### Proposal 2

- `source_type`: first-party coding-agent docs/changelog
- `problem_pattern`: plugin-defined commands and provider-compatible protocols are converging into default agent ergonomics
- `why_it_matters_for_harness`: reusable command surfaces and lighter provider glue make multi-agent workflows more portable and easier to compare
- `recommended_disposition`: `prototype`

### Proposal 3

- `source_type`: first-party multimodal API release notes
- `problem_pattern`: computer-use and conversational multimodal editing are being exposed as managed platform primitives rather than bespoke demos
- `why_it_matters_for_harness`: the competitive baseline for agent systems is shifting toward orchestrated action surfaces, not just model quality
- `recommended_disposition`: `watch`

## Watch Next

- Anthropic remains the highest-priority next-run check across the same required full surface: [newsroom](https://www.anthropic.com/news), [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code What's New](https://code.claude.com/docs/en/whats-new), and the raw [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md). The most likely next meaningful move is further managed-agent orchestration or a weekly Claude Code digest catching up to the raw changelog.
- OpenAI's next watch point is whether the official [Codex changelog](https://developers.openai.com/codex/changelog) posts a new in-window operator/governance surface rather than only extending remote continuity.
- Google's next watch point is whether [Gemini API](https://ai.google.dev/gemini-api/docs/changelog) reconnects its multimodal and computer-use surfaces to more explicit coding-agent/operator workflows.
- MiniMax, Qwen, and DeepSeek stay on watch, but without fresher dated official technical movement they remain comparison signals rather than front-line drivers for `erwin` this run.
