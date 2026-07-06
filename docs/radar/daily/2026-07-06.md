# Daily Radar - 2026-07-06

Window: 2026-06-29 to 2026-07-06

## New Signals

1. **Anthropic paired a model launch with a real Managed Agents control-plane upgrade.** The official [Introducing Claude Sonnet 5](https://www.anthropic.com/news/claude-sonnet-5) post on **2026-06-30** positioned Sonnet 5 as the most agentic Sonnet yet, available across Claude, Claude Code, and the Claude Platform. The official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) on **2026-06-30** landed the operator-side changes that matter more for `erwin`: Managed Agents **event deltas**, **backward session pagination**, **per-session agent overrides**, vault **`injection_location`** controls, and broader **agent/deployment/deployment-run webhooks**.  
   Publish date: 2026-06-30

2. **Anthropic kept the wider Claude surface moving after the June 30 launch burst.** The official [Redeploying Fable 5](https://www.anthropic.com/news/redeploying-fable-5) post on **2026-06-30** restored Fable 5 globally starting **2026-07-01** across Claude Platform, Claude.ai, Claude Code, and Claude Cowork, and the official [More details on Fable 5's cyber safeguards and our jailbreak framework](https://www.anthropic.com/news/fable-safeguards-jailbreak-framework) note on **2026-07-02** turned model availability into an explicit safety-and-governance story. Anthropic also launched [Claude Science](https://www.anthropic.com/news/claude-science-ai-workbench) on **2026-06-30**, which reinforces the pattern that Anthropic is shipping agentic product surfaces, not just models.  
   Publish dates: 2026-06-30, 2026-07-02

3. **Moonshot/Kimi compressed provider portability, plugin commands, and UI polish into one short release train.** The official [Kimi Code What's New](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html) shows **v0.22.0 on 2026-07-02** adding model configuration overrides and a redesigned web UI, **v0.21.0 on 2026-07-01** adding plugin-declared slash commands, and **v0.20.2 on 2026-06-29** adding an Anthropic-compatible protocol, video input, and custom outbound request headers. The pattern is clear: reusable command surfaces and cross-provider interoperability are becoming table stakes for coding agents.  
   Publish dates: 2026-07-02, 2026-07-01, 2026-06-29

4. **Google kept extending multimodal managed-action surfaces rather than just adding another model SKU.** The official [Gemini API release notes](https://ai.google.dev/gemini-api/docs/changelog) added **Gemini Omni Flash** in public preview on **2026-06-30** for high-speed video generation and conversational video editing through the Interactions API. This is less directly coding-agent-centric than Anthropic or Kimi this week, but it keeps raising the expectation that production agent platforms include native multimodal action loops.  
   Publish date: 2026-06-30

## GitHub Trending Snapshot

Snapshot source: [GitHub Trending](https://github.com/trending) captured on 2026-07-06.

- [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills): surfacing because reusable engineering skills are increasingly being treated as a product surface, not just prompt snippets.
- [openai/codex-plugin-cc](https://github.com/openai/codex-plugin-cc): relevant because cross-agent delegation between coding agents is becoming a distinct workflow category.
- [ogulcancelik/herdr](https://github.com/ogulcancelik/herdr): worth watching because terminal-resident agent multiplexing maps directly to operator-control questions.
- [gastownhall/gastown](https://github.com/gastownhall/gastown): notable because multi-agent workspace management is clearly attracting broader attention.
- [Leonxlnx/taste-skill](https://github.com/Leonxlnx/taste-skill): another signal that "skills" and opinionated instruction bundles are becoming first-class ecosystem artifacts.

## Harness-Relevant Signals

- Anthropic was explicitly checked across the required official surfaces this run: [Anthropic newsroom](https://www.anthropic.com/news), [Claude blog](https://claude.com/blog), [Anthropic developer docs](https://docs.anthropic.com), [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code What's New](https://code.claude.com/docs/en/whats-new), and the raw [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md).
- The [Managed Agents quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart) still centers the `managed-agents-2026-04-01` beta surface and the `agent -> environment -> session -> events` execution model. I did not find a separately dated in-window quickstart edit, but the June 30 release notes materially expanded the surrounding session-control surface.
- The [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) were checked explicitly and still present the `advisor-tool-2026-03-01` beta as a planner/executor split for long-horizon coding, computer-use, and research workloads. I did not find a newer dated in-window advisor note.
- The official [Claude Code What's New](https://code.claude.com/docs/en/whats-new) page still tops out at **Week 26 (2026-06-22 to 2026-06-26)**. The raw [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md) has moved beyond that digest, including `2.1.197` making Sonnet 5 the default model, `2.1.198` making background subagents the default while adding Claude Platform on AWS as an upstream provider, and `2.1.201` tightening Sonnet 5 harness-reminder behavior. That is a live product-surface signal, but not a separately dated weekly digest yet.
- OpenAI was checked through the official [Codex changelog](https://developers.openai.com/codex/changelog). I did **not** find a new official OpenAI item inside the 2026-06-29 to 2026-07-06 window. The latest clearly relevant dated items remain **Codex Remote GA on 2026-06-25** and **Record & Replay plus thread handoff on 2026-06-18**.
- MiniMax, Qwen, and DeepSeek were each checked once through official surfaces. I did **not** find a strong new dated official technical item inside this run window. The latest clear anchors I found remain [MiniMax M3](https://platform.minimax.io/docs/release-notes/models) on **2026-06-01**, [DeepSeek-V4](https://api-docs.deepseek.com/updates) on **2026-04-24**, and older [Qwen blog](https://qwen.ai/blog?id=qwen3.6) items outside this window without a clearer in-window dated launch.
- The strongest cross-publisher pattern this run is **operator surface formalization**: Anthropic expanded session controls and safety framing, Kimi turned plugins and provider compatibility into default ergonomics, and Google kept pushing managed multimodal action loops.

## Candidate Proposals

### Proposal 1

- `source_type`: first-party platform release notes
- `problem_pattern`: long-running managed-agent systems now need stream previews, reversible pagination, per-session steering, and lifecycle hooks instead of opaque run-and-poll behavior
- `why_it_matters_for_harness`: evaluation pressure is shifting from "can the agent finish?" toward "can operators observe, steer, resume, and govern the session cleanly?"
- `recommended_disposition`: `encode-as-eval`

### Proposal 2

- `source_type`: first-party coding-agent release notes and changelogs
- `problem_pattern`: plugins increasingly own command surfaces, and provider portability is becoming part of the default UX rather than custom glue
- `why_it_matters_for_harness`: reusable commands and interchangeable provider wiring make it easier to compare agent systems without baking every workflow into one stack
- `recommended_disposition`: `prototype`

### Proposal 3

- `source_type`: first-party product and safety notes
- `problem_pattern`: more capable agentic models are shipping with explicit safety-control narratives and governance primitives as part of the release package
- `why_it_matters_for_harness`: model capability and operator trust are increasingly bundled together, so evaluation should watch both execution quality and control-surface maturity
- `recommended_disposition`: `watch`

## Watch Next

- Anthropic remains the highest-priority next-run check across the same required surface: [newsroom](https://www.anthropic.com/news), [Claude blog](https://claude.com/blog), [Anthropic developer docs](https://docs.anthropic.com), [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code What's New](https://code.claude.com/docs/en/whats-new), and the raw [Claude Code changelog](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md). The most likely next meaningful change is a dated Claude Code weekly digest catching up to the raw changelog, or another Managed Agents orchestration/control note.
- OpenAI's next watch point is whether the official [Codex changelog](https://developers.openai.com/codex/changelog) posts a fresh operator/governance update after the **2026-06-25** Remote GA and **2026-06-18** Record & Replay/thread-handoff releases.
- Google's next watch point is whether [Gemini API](https://ai.google.dev/gemini-api/docs/changelog) reconnects its multimodal production surface to more explicit agent-operator workflows rather than keeping the story mostly model-centric.
- MiniMax, Qwen, and DeepSeek stay on watch, but without fresher dated official technical movement they remain comparison signals instead of front-line drivers for `erwin` this run.
