# Daily Radar - 2026-07-28

Window: 2026-07-22 to 2026-07-28

## New Signals

1. **Anthropic had real in-window movement across both Claude Platform and Claude Code, with Managed Agents also updated.** Anthropic's official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) include a **2026-07-22** entry for Managed Agents improvements, and the official [Claude Code changelog](https://code.claude.com/docs/en/changelog) includes a **2026-07-24** update that added `claude-opus-5` as the default Opus model plus `sandbox.network.strictAllowlist` and the `DirectoryAdded` hook. I also explicitly checked the official [Managed Agents docs](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and [Claude Code What's New](https://code.claude.com/docs/en/whats-new); those surfaces did not show a stronger separate in-window item than the July 22 and July 24 updates.
   Publish dates: 2026-07-22, 2026-07-24

2. **OpenAI's strongest in-window first-party signal is operational, not model-only: Presence is a production layer for deployed agents.** OpenAI's official announcement [Introducing OpenAI Presence](https://openai.com/index/introducing-openai-presence/) is dated **2026-07-22** and pushes reliability, oversight, and deployment operations further toward a packaged product surface for production voice and chat agents.
   Publish date: 2026-07-22

3. **MiniMax shipped a meaningful agent-client update around browser control and goal execution.** The official [MiniMax Agent changelog](https://agent.minimax.io/docs/changelog) shows version **v3.0.54** on **2026-07-27** adding browser-use control plus Claude Code and Gemini CLI support, followed by **v3.0.55** on **2026-07-28** with follow-up fixes. That is a stronger last-7-day agent signal than a generic model pricing or admin note.
   Publish dates: 2026-07-27, 2026-07-28

4. **DeepSeek's meaningful in-window event is a compatibility cutoff.** The official [DeepSeek API change log](https://api-docs.deepseek.com/updates/) ties the active migration to its V4 rollout and sets **2026-07-24** as the discontinuation date for the legacy `deepseek-chat` and `deepseek-reasoner` aliases. This is the sort of provider drift that silently breaks agent systems if compatibility checks lag behind vendor change logs.
   Cutoff date: 2026-07-24

5. **Google, Moonshot/Kimi, and Qwen were all checked, but none showed a stronger dated official technical update inside this run window than the items above.** The closest official Google items remained the **2026-07-21** Gemini posts just outside the window, Moonshot/Kimi's official [blog](https://www.kimi.com/blog/) still surfaced its latest visible dated research entries at **2026-07-16**, and Qwen's official [blog](https://qwen.ai/blog) did not expose a clearly dated stronger in-window technical item in this run. The signal is still useful: several major publishers were quieter this week while Anthropic, OpenAI, MiniMax, and DeepSeek moved the live edge.
   Dates checked: 2026-07-21, 2026-07-16

## GitHub Trending Snapshot

Snapshot source: [GitHub Trending](https://github.com/trending) captured from the default global board on 2026-07-28.

- [moeru-ai/airi](https://github.com/moeru-ai/airi): surfacing because self-hosted agent shells with voice, realtime interaction, and consumer-facing packaging still pull attention.
- [andrewyng/aisuite](https://github.com/andrewyng/aisuite): notable because multi-provider abstraction remains a live developer need as vendor APIs keep diverging.
- [microsoft/agent-governance-toolkit](https://github.com/microsoft/agent-governance-toolkit): highly relevant because policy enforcement, sandboxing, and zero-trust framing map directly onto production-agent governance pressure.
- [huggingface/speech-to-speech](https://github.com/huggingface/speech-to-speech): relevant because local voice-agent stacks remain a practical build surface rather than just a demo category.
- [bradautomates/claude-video](https://github.com/bradautomates/claude-video): relevant because richer evidence ingestion for coding agents continues to attract interest.

## Harness-Relevant Signals

- Anthropic's July 22 and July 24 updates reinforce that agent platforms are shipping both control-plane changes and execution-shell policy controls, not only model swaps.
- OpenAI Presence strengthens the pattern that production agent products now include reliability and operational governance as first-class surfaces.
- MiniMax's browser-use and goal-execution changes show continued pressure toward packaged operator-agent shells with direct tool control.
- DeepSeek's alias cutoff is an eval-worthy compatibility pattern: vendor naming drift can break runs without any model-quality regression.
- Google, Qwen, and Moonshot being quieter in-window is itself signal; the field did not move uniformly across all major publishers this week.

## Candidate Proposals

### Proposal 1

- `source_type`: first-party product and platform documentation
- `problem_pattern`: agent vendors are moving governance, sandbox, and execution controls into the core product surface
- `why_it_matters_for_harness`: Anthropic, OpenAI, and MiniMax all pushed the operational boundary outward this week, which increases the importance of tracking policy surfaces instead of only model capability claims
- `recommended_disposition`: `watch`

### Proposal 2

- `source_type`: first-party API and compatibility change logs
- `problem_pattern`: vendor compatibility breakage often arrives as alias, routing, or header drift rather than as a headline launch
- `why_it_matters_for_harness`: DeepSeek's alias retirement is a clean case for provider-drift coverage that should be caught before it becomes a runtime failure
- `recommended_disposition`: `encode-as-eval`

### Proposal 3

- `source_type`: first-party product releases plus secondary OSS momentum
- `problem_pattern`: richer tool control and evidence ingestion are becoming common expectations for agent shells
- `why_it_matters_for_harness`: MiniMax browser control, Claude Code policy hooks, and trending video/voice tooling all point toward better bounded tool-use and multimodal evidence paths
- `recommended_disposition`: `prototype`

## Watch Next

- Re-check Anthropic in the same required order next run: [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents docs](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code changelog](https://code.claude.com/docs/en/changelog), and [Claude Code What's New](https://code.claude.com/docs/en/whats-new).
- Watch for a separately dated Anthropic advisor-tool or Managed Agents docs update; this run's strongest official Anthropic items were the July 22 platform note and July 24 Claude Code release.
- Track whether OpenAI Presence gains a more technical developer or observability follow-through beyond the product launch page.
- Re-check official Google surfaces for a post-window follow-up to the July 21 Gemini segmentation posts rather than carrying that item forward indefinitely.
- Keep Moonshot/Kimi and Qwen on one-pass official watch only until they publish a clearly dated stronger technical item again.
