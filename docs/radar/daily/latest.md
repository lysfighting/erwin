# AI Agent Daily Radar - 2026-04-15

Window: 2026-04-08 to 2026-04-15

## New Signals

- **Claude Code shipped multiple in-window operations upgrades.** The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) lists version 2.1.108 on **April 14, 2026**, adding 1-hour prompt-cache TTL controls, session recap, slash-command discovery through the Skill tool, better resume defaults, sharper rate-limit errors, and fixes for Agent-tool permission handling and transcript durability. The **April 13, 2026** 2.1.105 release added `EnterWorktree` path switching, PreCompact hooks, plugin background monitors, stale worktree cleanup, MCP truncation guidance, scheduled-task repair, and headless/remote MCP fixes. This is the clearest Claude Code in-window signal: coding-agent reliability is moving into cache policy, resumability, hooks, monitors, worktrees, and remote/headless execution.
- **Anthropic added a platform migration clock for older Claude 4 models.** The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) dated **April 14, 2026** announce deprecation of `claude-sonnet-4-20250514` and `claude-opus-4-20250514`, with API retirement scheduled for **June 15, 2026**, and recommend migration to Claude Sonnet 4.6 and Claude Opus 4.6. For harness work, this is less a feature launch than a dependency hygiene signal: long-running agent systems need explicit model-lifecycle checks, migration tests, and stale-model detection.
- **Anthropic's Managed Agents public beta remains the strongest managed-runtime signal in-window.** The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) date the launch to **April 8, 2026**, and the [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) describes a pre-built configurable agent harness for long-running and asynchronous work, with agents, environments, sessions, events, built-in tools, MCP, skills, server-sent event streaming, and beta access via `managed-agents-2026-04-01`. The paired Anthropic engineering post, [Scaling Managed Agents: Decoupling the brain from the hands](https://www.anthropic.com/engineering/managed-agents), makes the architecture explicit: session log, harness, and sandbox are separated so each can fail or evolve independently.
- **Anthropic's advisor tool formalizes executor/advisor model splitting.** The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) date the public beta to **April 9, 2026**, and the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) define a lower-cost executor model consulting a higher-intelligence advisor model mid-generation under the `advisor-tool-2026-03-01` beta. The docs expose runtime details that matter for harnesses: advisor usage appears in `usage.iterations`, advisor failures return as tool-result errors, streaming pauses while advisor inference runs, and prompt caching can be enabled separately for advisor calls.
- **OpenAI's Cloudflare Agent Cloud post turns Codex harness into deployment infrastructure.** OpenAI published [Enterprises power agentic workflows in Cloudflare Agent Cloud with OpenAI](https://openai.com/index/cloudflare-openai-agent-cloud/) on **April 13, 2026**. The post says Cloudflare Agent Cloud customers can deploy agents powered by GPT-5.4, enterprises can deploy agents built on the Codex harness to Cloudflare, and the Codex harness is generally available in Cloudflare Sandboxes. This is a distribution and operations signal: first-party coding-agent harnesses are becoming hosted enterprise runtime components, not only local developer tools.

## GitHub Trending Snapshot

Secondary ecosystem demand, not first-party field signal. Snapshot checked from [GitHub Trending daily](https://github.com/trending?since=daily) and the [OSSInsight AI trending view](https://ossinsight.io/trending/ai) on **2026-04-15**.

- [forrestchang/andrej-karpathy-skills](https://github.com/forrestchang/andrej-karpathy-skills) led the GitHub daily page as a single `CLAUDE.md` behavior file for Claude Code. It is surfacing because developers are treating durable agent instructions and coding-agent behavior shaping as shareable infrastructure.
- [thedotmack/claude-mem](https://github.com/thedotmack/claude-mem) was prominent on the GitHub daily page as a Claude Code plugin for capturing, compressing, and reinjecting coding-session context. It is surfacing because long-lived agent memory and context handoff remain unsolved pain points.
- [obra/superpowers](https://github.com/obra/superpowers) appeared on the GitHub daily page as an agentic skills framework and software-development methodology. It is relevant because skills, repeatable workflows, and agent operating conventions are becoming repository artifacts.
- [NousResearch/hermes-agent](https://github.com/NousResearch/hermes-agent) appeared on the GitHub daily page as an "agent that grows with you." It is surfacing around personalization and adaptive agent behavior, a secondary signal to watch rather than a first-party runtime fact.
- [anthropics/claude-code](https://github.com/anthropics/claude-code), [anomalyco/opencode](https://github.com/anomalyco/opencode), [openai/codex](https://github.com/openai/codex), and [block/goose](https://github.com/block/goose) remain top AI/coding-agent movers in the OSSInsight view. Their continued movement reinforces that terminal-native coding agents and open local runners are a durable category.

## Harness-Relevant Signals

- Claude Code's April 10-14 releases are tightly aligned with harness concerns: prompt-cache TTLs, session recap, worktree switching, compaction hooks, plugin background monitors, MCP failure clarity, headless/remote reliability, scheduled-task repair, and OpenTelemetry prompt/tool-detail controls.
- Managed Agents gives a concrete external runtime vocabulary for long-running agents: agent, environment, session, event, toolset, MCP, skill, sandbox, stream, and persisted session log.
- The Managed Agents engineering post is especially relevant because it argues against assuming a single container or a single harness shape. The notable pattern is durable session state outside the context window, with harness and sandbox allowed to fail or be swapped.
- The advisor tool turns strategic review into an explicit callable runtime primitive with cost, caching, streaming, and failure semantics. That is a better comparison point than ad hoc "ask a stronger model" orchestration.
- OpenAI's Cloudflare signal and Alibaba's earlier [Qwen3.6-Plus](https://www.alibabacloud.com/en/press-room/alibaba-unveils-qwen3-6-plus-to-accelerate-agentic?_p_lc=1) announcement on **April 2, 2026** both point toward repository-level engineering and deployable agent loops, but Qwen's dated launch falls outside this run's primary 7-day signal window.

## Candidate Proposals

- `source_type`: `first_party_coding_agent_changelog`
  `problem_pattern`: `coding_agent_reliability_is_converging_on_resumability_hooks_monitors_mcp_and_cache_policy`
  `why_it_matters_for_harness`: Claude Code is hardening around the same operational boundaries Harness needs to evaluate: state handoff, compaction control, background work, tool availability, cache lifetimes, and remote/headless failure modes.
  `recommended_disposition`: `encode-as-eval`
- `source_type`: `first_party_managed_agent_runtime`
  `problem_pattern`: `agent_runtime_contracts_are_becoming_explicit_product_interfaces`
  `why_it_matters_for_harness`: Managed Agents provides a stable external vocabulary for session, harness, sandbox, event stream, environment, MCP, skills, and recovery semantics without requiring Harness to copy Anthropic's implementation.
  `recommended_disposition`: `watch`
- `source_type`: `first_party_executor_advisor_tool`
  `problem_pattern`: `long_horizon_agents_need_visible_strategy_execution_and_cost_boundaries`
  `why_it_matters_for_harness`: The advisor tool exposes where strategic guidance happens, how it is billed, how it fails, and how it interacts with streaming and caching, which are all testable orchestration properties.
  `recommended_disposition`: `prototype`

## Watch Next

- Anthropic: keep explicitly checking [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents docs](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs](https://code.claude.com/docs/), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). Next high-value items are Managed Agents endpoint changes, advisor observability/cost-control changes, scheduled-task release notes, and further Claude Code monitor/hook/worktree behavior.
- OpenAI: watch [OpenAI news](https://openai.com/news/) and developer docs for follow-ups to the **April 13, 2026** [Cloudflare Agent Cloud](https://openai.com/index/cloudflare-openai-agent-cloud/) Codex harness deployment signal.
- Google: no stronger clearly dated official Google AI/Gemini agent-runtime or coding-agent technical item was found in the checked **2026-04-08 to 2026-04-15** window. Continue watching [Google AI for Developers](https://ai.google.dev/) and the [Gemini API changelog](https://ai.google.dev/gemini-api/docs/changelog).
- Moonshot, MiniMax, Qwen, and DeepSeek: the single official search pass did not surface a stronger in-window dated technical agent/runtime item than the Anthropic/OpenAI signals. Keep watching official model/API release notes; do not elevate aggregator-only claims.
- Secondary signal sources: use [Langfuse](https://langfuse.com/), [Phoenix](https://phoenix.arize.com/), [Helicone](https://www.helicone.ai/), [OpenTelemetry](https://opentelemetry.io/), and GitHub Trending only to stress-test first-party trends around observability, replay, tracing, and agent context, not as primary launch evidence.
