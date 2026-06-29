# Daily Radar - 2026-06-28

Window: 2026-06-21 to 2026-06-28

## New Signals

1. **Anthropic shipped two concrete Claude Platform primitives for longer-running agent loops.** In the official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), Anthropic added **cache diagnostics in beta on 2026-06-26** and **URL support for images and PDFs in code execution on 2026-06-25**. Together, those updates reduce context-debugging friction and make file-heavy tool chains easier to drive without upload-only flows.  
   Publish dates: 2026-06-26 and 2026-06-25

2. **Claude Code kept moving on operator control instead of just raw model capability.** The official [Claude Code What's New](https://code.claude.com/docs/en/whats-new) page shows **Week 26 covering 2026-06-22 to 2026-06-26**, led by `claude mcp login`, shell-mode follow-up explanations, `/rewind` recovery after `/clear`, and background subagent permission prompts surfacing in the main session.  
   Publish window: 2026-06-22 to 2026-06-26

3. **OpenAI pushed remote coding-agent continuity further into a default operating surface.** The official [Codex changelog](https://developers.openai.com/codex/changelog) shows **Codex Remote reached general availability on 2026-06-25**, with mobile start/continue flows on connected Mac or Windows hosts, authenticated one-to-one QR pairing, and a new DigitalOcean plugin for remote workspaces.  
   Publish date: 2026-06-25

4. **Google moved computer use from preview novelty toward a managed API capability.** The official [Gemini API release notes](https://ai.google.dev/gemini-api/docs/changelog) show **Computer Use entered public preview on 2026-06-24** for Gemini 3.5 Flash, including intent-level actions, built-in browser/mobile/desktop environments, configurable safety policies, and prompt-injection detection.  
   Publish date: 2026-06-24

5. **Moonshot/Kimi kept hardening the coding-agent UX around backgrounded work and remote access.** The official [Kimi Code What's New](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html) page shows **Kimi Code CLI v0.20.0 on 2026-06-26** adding shell mode, internet-exposed `kimi web --host` with token auth and rate limiting, and a redesigned plugin manager; **v0.19.0 on 2026-06-22** added extra workspace directories and background tasks.  
   Publish dates: 2026-06-26 and 2026-06-22

## GitHub Trending Snapshot

Snapshot source: [GitHub Trending](https://github.com/trending) captured on 2026-06-28.

- [simplex-chat/simplex-chat](https://github.com/simplex-chat/simplex-chat): not an AI repo, but it is surfacing because privacy-preserving communication infrastructure remains an adjacent systems theme.
- [xbtlin/ai-berkshire](https://github.com/xbtlin/ai-berkshire): notable because it packages a concrete Claude Code / Codex multi-agent workflow around a real domain instead of a generic demo.
- [DeusData/codebase-memory-mcp](https://github.com/DeusData/codebase-memory-mcp): relevant because persistent repo memory and code-intelligence MCP servers are becoming default building blocks for coding agents.
- [opendatalab/MinerU](https://github.com/opendatalab/MinerU): surfacing because document-to-LLM pipelines remain an important ingestion layer for agentic workflows.
- [usestrix/strix](https://github.com/usestrix/strix): worth watching as an "AI hackers fix your app" framing for security automation.
- [browser-use/video-use](https://github.com/browser-use/video-use): notable because browser-use style agent tooling is spreading into adjacent creator workflows rather than staying purely browser-automation oriented.

## Harness-Relevant Signals

- The strongest first-party pattern this week is not "new model, same workflow"; it is **better operator surfaces for long-running agents**: remote continuity from OpenAI, computer-use environments from Google, and cache / session / permission controls from Anthropic and Kimi.
- Anthropic's required technical surfaces were explicitly checked this run: [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code What's New](https://code.claude.com/docs/en/whats-new), [Claude Code changelog](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md), [Anthropic News](https://www.anthropic.com/news), and [Claude Blog](https://claude.com/blog). I found meaningful in-window updates in the Platform release notes and Claude Code docs, but no stronger in-window dated change on the Managed Agents quickstart or advisor tool pages beyond their current beta surfaces.
- The Managed Agents quickstart still matters strategically even without a fresh in-window note: it remains explicitly beta-gated under `managed-agents-2026-04-01`, with agent, environment, session, and event-stream primitives spelled out on the official page.
- The advisor tool page remains a clean statement of an increasingly common pattern: a cheaper executor can call a stronger planner mid-turn under the `advisor-tool-2026-03-01` beta header, which is directly relevant to long-horizon coding and research agents.
- MiniMax, Qwen, and DeepSeek were each checked once through official surfaces. I did not find a stronger dated technical item inside the 2026-06-21 to 2026-06-28 window. The latest dated official items I found remain [MiniMax M3](https://platform.minimax.io/docs/release-notes/models) on 2026-06-01 and [DeepSeek V4 Preview](https://api-docs.deepseek.com/news/news260424) on 2026-04-24; the Qwen official surfaces showed active research/blog motion but not a clearer in-window coding-agent platform update than the higher-signal items above.

## Candidate Proposals

### Proposal 1

- `source_type`: first-party platform release notes
- `problem_pattern`: long-running agents now need explicit context-debugging and file-ingestion controls instead of opaque "just trust the loop" behavior
- `why_it_matters_for_harness`: cache diagnostics, URL-native file inputs, and session-state inspection point toward more inspectable and governable agent runs
- `recommended_disposition`: `encode-as-eval`

### Proposal 2

- `source_type`: first-party coding-agent changelog
- `problem_pattern`: remote-host continuity and mobile approvals are becoming part of the operator contract for coding agents
- `why_it_matters_for_harness`: serious agent systems increasingly assume users will resume, steer, and approve work away from the original terminal
- `recommended_disposition`: `watch`

### Proposal 3

- `source_type`: first-party agent API release notes
- `problem_pattern`: computer use is shipping with built-in environment types and safety policy controls, not as a free-form tool bolt-on
- `why_it_matters_for_harness`: agent evaluation pressure is moving toward governed end-to-end execution in real environments
- `recommended_disposition`: `prototype`

## Watch Next

- Anthropic next run watch remains mandatory across the same full surface: [release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code What's New](https://code.claude.com/docs/en/whats-new), and [Claude Code changelog](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md). The most likely next meaningful move is either managed-agent orchestration hardening or another Claude Code control-surface expansion.
- OpenAI's next watch point is whether [Codex Remote](https://developers.openai.com/codex/changelog) expands from pairing and continuation into more explicit workflow governance, deployment, or artifact controls.
- Google's next watch point is whether [Computer Use](https://ai.google.dev/gemini-api/docs/changelog) converges more tightly with its managed-agent APIs instead of remaining a separate execution primitive.
- Kimi is worth tracking for the pace of coding-agent UX iteration, but for `erwin` it remains a secondary comparison source unless its open platform or agent architecture surfaces deepen beyond CLI ergonomics.
