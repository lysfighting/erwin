# Daily Radar - 2026-06-25

Window: 2026-06-18 to 2026-06-25

## New Signals

1. **Google moved computer-use from concept to public-preview API surface.** The official [Gemini API release notes](https://ai.google.dev/gemini-api/docs/changelog) added **Computer Use** for `gemini-3.5-flash` on **2026-06-24**, with simplified intent actions, browser/mobile/desktop environment support, configurable safety policies, and prompt-injection detection.
   Publish date: 2026-06-24

2. **Moonshot/Kimi is hardening the coding-agent session model, not just the model weights.** The official [Kimi Code What's New](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html) shows **v0.19.0** on **2026-06-22** adding extra workspace directories plus background tasks, and **v0.18.0** on **2026-06-18** adding web session filtering and lazy loading. This is direct pressure toward multi-dir, long-running agent workflows.
   Publish dates: 2026-06-18, 2026-06-22

3. **OpenAI expanded Codex’s reusable-workflow and cross-host operating surface.** The official [Codex changelog](https://developers.openai.com/codex/changelog) shows **Codex app 26.616** on **2026-06-18** adding [Record & Replay](https://developers.openai.com/codex/record-and-replay), automation bulk actions, and thread handoff between local and remote hosts.
   Publish date: 2026-06-18

4. **OpenAI also kept tightening the human control layer around long-running agent work.** The official [Codex changelog](https://developers.openai.com/codex/changelog) shows **ChatGPT for iOS 1.2026.167** on **2026-06-22** adding editable goals plus better visibility into subagents, tasks, and worktree creation. This is UI work, but it points in the same direction: persistent agent work needs operator-grade steering and status.
   Publish date: 2026-06-22

## GitHub Trending Snapshot

Snapshot source: [GitHub Trending](https://github.com/trending) captured for this run on 2026-06-25.

- [apple/container](https://github.com/apple/container): container/runtime infrastructure surfacing on the default board is relevant because agent systems keep converging on sandboxed execution rather than pure chat.
- [interviewstreet/hiring-agent](https://github.com/interviewstreet/hiring-agent): explicit hiring-agent packaging on the default board is a useful signal that narrow workflow agents are still attracting attention beyond model-launch cycles.
- [JCodesMore/ai-website-cloner-template](https://github.com/JCodesMore/ai-website-cloner-template): shows continued demand for turnkey AI-assisted generation workflows, especially where repo templates lower setup friction.
- [revfactory/harness](https://github.com/revfactory/harness): worth watching because “harness” naming is appearing on the global board, even if this repo may not map directly to the `Harness` system in this repo.
- [stablyai/orca](https://github.com/stablyai/orca): another AI-native entry surfacing on the broad board rather than only model-specific channels.
- [google-labs-code/design.md](https://github.com/google-labs-code/design.md): notable as a design/code-adjacent repo from an official Google Labs org, suggesting continuing appetite for design-to-build workflows.

## Harness-Relevant Signals

- Google and OpenAI are both making computer-use and reusable workflow capture look like first-class platform features, not experimental add-ons. That raises the bar for how agent systems expose execution control, replay, and safety boundaries.
- Kimi’s June 18-22 updates matter because they focus on session continuity across directories, background work, and browser-visible history. The competitive frontier for coding agents is increasingly operational ergonomics.
- Anthropic’s required surfaces were explicitly checked this run: [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs](https://code.claude.com/docs/en/whats-new), [claude.com/blog](https://claude.com/blog), and [Anthropic engineering/news](https://www.anthropic.com/engineering). I did not find a stronger dated official technical item inside the 2026-06-18 to 2026-06-25 window than Anthropic’s already-known **2026-06-15** model retirement note and **2026-06-11** tool-version updates.
- MiniMax, Qwen, and DeepSeek were also checked in one official pass each. I did not find a stronger in-window official technical item than the signals above. The latest relevant dated items I found were [MiniMax M3](https://platform.minimax.io/docs/release-notes/models) on **2026-06-01**, [Qwen-Robot Suite](https://qwen.ai/blog?id=qwen-robotsuite) on **2026-06-15**, and [DeepSeek-V4](https://api-docs.deepseek.com/updates) on **2026-04-24**.

## Candidate Proposals

### Proposal 1

- `source_type`: first-party API release notes
- `problem_pattern`: computer-use is becoming an explicit product surface with safety policy controls, environment abstractions, and injection defenses
- `why_it_matters_for_harness`: the competitive baseline is shifting from “tool calling exists” to “desktop/browser action surfaces are inspectable and governable”
- `recommended_disposition`: `encode-as-eval`

### Proposal 2

- `source_type`: first-party coding-agent changelog
- `problem_pattern`: session continuity is expanding across directories, browsers, background tasks, and remote-host handoff
- `why_it_matters_for_harness`: long-running agent work now depends on continuity UX and operator controls, not only model quality
- `recommended_disposition`: `prototype`

### Proposal 3

- `source_type`: first-party product changelog
- `problem_pattern`: reusable workflow capture and replay are being packaged as native agent capabilities
- `why_it_matters_for_harness`: demonstrated-task capture could become a durable comparison axis for agent systems beyond prompt libraries
- `recommended_disposition`: `watch`

## Watch Next

- Anthropic remains the highest-priority watch even without a stronger in-window item. The required surfaces for the next run are unchanged: [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code what's new](https://code.claude.com/docs/en/whats-new), [Claude Code releases](https://github.com/anthropics/claude-code/releases), [claude.com/blog](https://claude.com/blog), and [Anthropic engineering](https://www.anthropic.com/engineering).
- Google’s next important watch point is whether the [Gemini API release notes](https://ai.google.dev/gemini-api/docs/changelog) extend Computer Use beyond preview and how quickly managed-agent and interaction APIs converge around one operator model.
- OpenAI’s next watch point is whether [Record & Replay](https://developers.openai.com/codex/record-and-replay), remote handoff, and editable long-running goals stay app-only conveniences or become stronger code-/API-addressable workflow primitives.
- Kimi is worth watching for whether the June 18-22 session-management changes lead to stronger team, replay, or automation surfaces in [Kimi Code What's New](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html).
