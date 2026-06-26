# Daily Radar - 2026-06-26

Window: 2026-06-19 to 2026-06-26

## New Signals

1. **OpenAI pushed Codex from “connected desktop” toward a real remote-work substrate.** The official [Codex changelog](https://developers.openai.com/codex/changelog) shows **Codex Remote reached general availability on 2026-06-25**, with mobile start/continue flows for connected Mac and Windows hosts, one-to-one QR pairing, and a new DigitalOcean plugin for provisioning a remote workspace.
   Publish date: 2026-06-25

2. **Google moved computer use into a broader public-preview agent surface.** The official [Gemini API release notes](https://ai.google.dev/gemini-api/docs/changelog) show **Computer Use entered public preview on 2026-06-24** for `gemini-3.5-flash`, with intent-level actions, built-in browser/mobile/desktop environments, configurable safety policies, and prompt-injection detection.
   Publish date: 2026-06-24

3. **OpenAI also tightened long-running agent control and budget mechanics on the CLI path.** The official [Codex changelog](https://developers.openai.com/codex/changelog) shows **Codex CLI 0.142.0 on 2026-06-22** adding configurable rollout token budgets, indexed web-search mode, and thread/turn-level control over multi-agent delegation.
   Publish date: 2026-06-22

## GitHub Trending Snapshot

Snapshot source: [GitHub Trending](https://github.com/trending) captured on 2026-06-26.

- [google-labs-code/design.md](https://github.com/google-labs-code/design.md): surfacing because structured design-system instructions for coding agents are becoming a reusable artifact, not just prompt folklore.
- [calesthio/OpenMontage](https://github.com/calesthio/OpenMontage): a strong momentum signal for visibly agentic, multi-tool production workflows rather than simple chat wrappers.
- [apple/container](https://github.com/apple/container): relevant because trustworthy local execution and lightweight isolation remain core infrastructure for serious agent systems.
- [aws/agent-toolkit-for-aws](https://github.com/aws/agent-toolkit-for-aws): notable as an official MCP/skills/plugins bundle from a hyperscaler, which reinforces the shift toward packaged agent-tool ecosystems.
- [alibaba/page-agent](https://github.com/alibaba/page-agent): worth watching because in-page GUI control continues to emerge as a concrete interface layer for browser agents.

## Harness-Relevant Signals

- OpenAI and Google are converging on the same operator model: persistent agents need remote continuity, explicit control points, and bounded execution rather than one-shot prompting.
- The Google Computer Use release raises the evaluation bar from “can call tools” to “can act in browser/mobile/desktop environments with inspectable safety controls.”
- OpenAI’s June 22-25 Codex updates make budgets, delegation policy, indexed search, and remote-host continuity look like baseline platform concerns for coding agents.
- Anthropic’s required technical surfaces were explicitly checked this run: [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code what's new](https://code.claude.com/docs/en/whats-new), [Claude Code changelog](https://docs.anthropic.com/en/release-notes/claude-code), [Claude blog](https://claude.com/blog), and [Anthropic news](https://www.anthropic.com/news). I did not find a stronger dated Anthropic technical item inside the 2026-06-19 to 2026-06-26 window; the newest meaningful official Anthropic items remain just outside it on **2026-06-18** ([Claude Code artifacts](https://claude.com/blog/artifacts-in-claude-code) and [managed MCP connector auth](https://claude.com/blog/enterprise-managed-auth)).
- MiniMax, Qwen, DeepSeek, and Moonshot/Kimi were each checked once through official surfaces. I did not find a stronger in-window dated technical item from those families. The latest relevant dated items I found remain [MiniMax M3](https://platform.minimax.io/docs/release-notes/models) on **2026-06-01**, [Qwen-Robot Suite](https://qwen.ai/blog?id=qwen-robotsuite) on **2026-06-15**, [DeepSeek-V4](https://api-docs.deepseek.com/updates) on **2026-04-24**, and [Kimi K2.7 Code](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html) on **2026-06-12**.

## Candidate Proposals

### Proposal 1

- `source_type`: first-party API release notes
- `problem_pattern`: computer-use capabilities are shipping with explicit environment abstractions, safety policies, and injection defenses
- `why_it_matters_for_harness`: agent evaluation pressure is moving from abstract tool use toward governed end-to-end interaction inside real execution surfaces
- `recommended_disposition`: `encode-as-eval`

### Proposal 2

- `source_type`: first-party coding-agent changelog
- `problem_pattern`: remote-host continuity, mobile approvals, and session resumption are becoming native parts of coding-agent operation
- `why_it_matters_for_harness`: the operator surface for long-running work is becoming a product requirement rather than an optional convenience
- `recommended_disposition`: `prototype`

### Proposal 3

- `source_type`: first-party product changelog
- `problem_pattern`: token budgets, delegation policy, and indexed search are being productized as hard execution controls
- `why_it_matters_for_harness`: governance and bounded autonomy are becoming directly comparable dimensions across agent systems
- `recommended_disposition`: `watch`

## Watch Next

- Anthropic remains high-priority next run watch even without an in-window item this time: [release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code what's new](https://code.claude.com/docs/en/whats-new), and [Claude Code changelog](https://docs.anthropic.com/en/release-notes/claude-code).
- OpenAI’s next important watch point is whether [Codex Remote](https://developers.openai.com/codex/changelog) and the June 22 CLI control surfaces expand into more explicit workflow primitives rather than staying app-centric.
- Google’s next watch point is whether [Computer Use](https://ai.google.dev/gemini-api/docs/changelog) quickly converges with its managed-agent surfaces into one operator model.
- The most relevant secondary ecosystem watch from today’s Trending board is whether [aws/agent-toolkit-for-aws](https://github.com/aws/agent-toolkit-for-aws), [alibaba/page-agent](https://github.com/alibaba/page-agent), and [google-labs-code/design.md](https://github.com/google-labs-code/design.md) keep compounding attention or fade as one-day novelty.
