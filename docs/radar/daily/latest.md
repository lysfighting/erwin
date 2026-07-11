# Daily Radar - 2026-07-11

Window: 2026-07-04 to 2026-07-11

## New Signals

1. **OpenAI shipped the GPT-5.6 family with an explicit multi-agent execution tier.** The official [GPT-5.6 launch](https://openai.com/index/gpt-5-6/) on **2026-07-09** introduced Sol, Terra, and Luna, plus `ultra`, a highest-capability setting that coordinates multiple agents across parallel workstreams for faster completion of complex tasks.  
   Publish date: 2026-07-09

2. **OpenAI also launched GPT-Live as a full-duplex voice layer that can keep talking while deeper work runs in the background.** The official [Introducing GPT-Live](https://openai.com/index/introducing-gpt-live/) post on **2026-07-08** says GPT-Live listens and speaks at the same time, delegates harder work to GPT-5.5 at launch, and is planned for the API later.  
   Publish date: 2026-07-08

3. **Anthropic had a meaningful in-window platform-and-Code sweep, not a quiet week.** In the official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), Anthropic added expiring API keys on **2026-07-08** and changed Managed Agents memory-store list semantics behind the `agent-memory-2026-07-22` header on **2026-07-02**. In the official [Claude Code What's New](https://code.claude.com/docs/en/whats-new), Week 28 (**2026-07-06** to **2026-07-10**) added the in-app browser on Desktop, `/doctor`, stricter auto-mode protections, and richer `claude agents` state rows.  
   Publish dates: 2026-07-08, 2026-07-02, 2026-07-06 to 2026-07-10

4. **MiniMax exposed a server-hosted tool surface on its Anthropic-compatible API instead of only shipping another model headline.** The official [Server Tools](https://platform.minimax.io/docs/guides/server-tools) docs, surfaced in page metadata on **2026-07-08**, describe built-in server-side `web_search` on `/anthropic/v1/messages`, eliminating the usual client-managed `tool_use` / `tool_result` loop for supported cases.  
   Publish date: 2026-07-08

5. **Qwen's coding-agent stream is moving fast at the runtime layer, not just model branding.** The official [Qwen Code releases](https://github.com/QwenLM/qwen-code/releases) shipped `v0.19.9` on **2026-07-10**, with fixes for repeated subagent tool-call loops, broken history-chain detection, configurable background-memory timeouts, persisted session artifacts, and isolated scheduled-task runs.  
   Publish date: 2026-07-10

## GitHub Trending Snapshot

Snapshot source: [GitHub Trending](https://github.com/trending) captured on 2026-07-11.

- [davila7/claude-code-templates](https://github.com/davila7/claude-code-templates): surfacing because operator-ready Claude Code workflow packaging is becoming reusable project scaffolding instead of one-off prompt craft.
- [google-labs-code/stitch-skills](https://github.com/google-labs-code/stitch-skills): notable because skills are becoming a portable extension surface for coding agents, not just local prompt snippets.
- [openai/plugins](https://github.com/openai/plugins): relevant because plugin and tool packaging remains a shared coordination surface across agent runtimes even as MCP keeps rising.
- [wonderwhy-er/DesktopCommanderMCP](https://github.com/wonderwhy-er/DesktopCommanderMCP): worth watching because MCP-backed desktop control remains one of the clearest bridges from chat agents to real operator actions.
- [hashicorp/terraform](https://github.com/hashicorp/terraform): useful board-shape reminder that infrastructure control planes still trend beside agent tooling rather than the board becoming AI-only.

## Harness-Relevant Signals

- Anthropic was checked explicitly across the required surfaces before concluding anything about its quiet or active substreams: [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [Managed Agents memory docs](https://platform.claude.com/docs/en/managed-agents/memory), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code What's New](https://code.claude.com/docs/en/whats-new), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). The important in-window takeaway is that Managed Agents memory semantics changed on **2026-07-02**, while Claude Code's Week 28 reinforced browser-native and operator-safety runtime patterns on **2026-07-06** to **2026-07-10**.
- The official [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) still describe one of the clearest first-party executor/advisor splits: a lower-cost executor can consult a stronger model mid-generation under `advisor-tool-2026-03-01`. That remains a strong Harness pattern even without an in-window advisor-specific release.
- OpenAI's [GPT-5.6 launch](https://openai.com/index/gpt-5-6/) and [ChatGPT Work](https://openai.com/index/chatgpt-for-your-most-ambitious-work/) together sharpen the same architectural pressure: a user-facing runtime is being paired with longer-running, multi-agent execution that can produce slides, docs, sheets, and web artifacts without collapsing back to a single synchronous chat turn.
- MiniMax's [Server Tools](https://platform.minimax.io/docs/guides/server-tools) page is strategically relevant because it imports Anthropic-style dated tool identifiers such as `web_search_20250305` while moving supported tool execution fully server-side on the Anthropic-compatible surface.
- Qwen Code's [v0.19.9 release](https://github.com/QwenLM/qwen-code/releases) is worth tracking as a coding-agent implementation stream because the visible fixes cluster around subagent loops, memory/runtime persistence, scheduled isolated runs, and session integrity instead of only UI polish.
- I did not find a stronger dated **2026-07-04** to **2026-07-11** official technical item from Google, Moonshot/Kimi, or DeepSeek after one pass each through their primary publisher surfaces. DeepSeek remains operationally relevant because the official [DeepSeek API changelog](https://api-docs.deepseek.com/updates/) still says the legacy `deepseek-chat` and `deepseek-reasoner` names will be discontinued on **2026-07-24**.

## Candidate Proposals

### Proposal 1

- `source_type`: first-party product launch plus first-party platform docs
- `problem_pattern`: higher-level user interaction surfaces are being separated from longer-running agent execution tiers
- `why_it_matters_for_harness`: OpenAI's `ultra` and ChatGPT Work, plus Anthropic's Managed Agents and advisor split, all reinforce that orchestration is becoming a first-class runtime surface
- `recommended_disposition`: `watch`

### Proposal 2

- `source_type`: first-party API docs
- `problem_pattern`: tool execution is moving from client-managed orchestration loops into hosted control planes
- `why_it_matters_for_harness`: MiniMax's server-side `web_search` and Anthropic's Managed Agents memory/session model both reduce client glue while increasing vendor-defined runtime semantics
- `recommended_disposition`: `prototype`

### Proposal 3

- `source_type`: first-party coding-agent release notes
- `problem_pattern`: coding-agent runtimes are converging on persistence, subagent coordination, scheduled isolation, and operator-state visibility
- `why_it_matters_for_harness`: Qwen Code and Claude Code both show that session integrity and background-run ergonomics are moving into the product core
- `recommended_disposition`: `encode-as-eval`

## Watch Next

- Anthropic stays at the top of the next-run sweep. Re-check the official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents docs](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code What's New](https://code.claude.com/docs/en/whats-new), and [Claude Code changelog](https://code.claude.com/docs/en/changelog) for follow-through on memory-store behavior, browser-native workflows, or advisor/runtime updates.
- OpenAI's next watchpoint is whether GPT-Live and ChatGPT Work land as developer-facing API/control-plane primitives rather than remaining mostly product surfaces after the **2026-07-08** and **2026-07-09** launches.
- DeepSeek should be re-checked before and after the **2026-07-24** legacy-model-name cutoff because migration friction can become more important than a brand-new model launch.
- Google, Moonshot/Kimi, and DeepSeek were not inflated into false novelty this run; keep them on the board, but only promote them when a dated primary technical surface changes.
