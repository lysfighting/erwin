# AI Agent Daily Radar - 2026-05-25

Window: 2026-05-18 to 2026-05-25

## New Signals

- **Anthropic expanded Claude Managed Agents on 2026-05-19.** The official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) and matching [Claude blog post](https://claude.com/blog/new-in-claude-managed-agents-self-hosted-sandboxes-and-mcp-tunnels) added MCP tunnels in research preview, self-hosted sandboxes, live MCP and tool reconfiguration for active sessions, and automatic spill-to-file handling for outputs above 100K tokens.
- **OpenAI sharpened the Codex operator surface on 2026-05-21.** The official [ChatGPT release notes](https://help.openai.com/en/articles/6825453-chatgpt-release-notes%252525252525252525252525252525252525252525252525252525252525252525252525252523.doc) added Appshots, general-availability Goal mode, in-app browser annotations, locked computer use, and broader browser-use improvements for longer-running coding workflows.
- **Google launched Managed Agents in the Gemini API on 2026-05-19.** In [Introducing Managed Agents in the Gemini API](https://blog.google/innovation-and-ai/technology/developers-tools/managed-agents-gemini-api/), Google exposed a single-call agent runtime with isolated Linux execution, resumable state, web access, and repo-local agent contracts via `AGENTS.md` and `SKILL.md`.
- **Qwen pushed an agent-centered flagship on 2026-05-20.** The official [Qwen research feed](https://qwen.ai/research/) introduced Qwen3.7-Max as a “versatile agent foundation” for coding, office automation, and sustained autonomous execution, making it the clearest in-window model signal from the China vendors covered this run.
- **Anthropic Claude Code kept moving on 2026-05-21, 2026-05-22, and 2026-05-23.** The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows `2.1.147` on **2026-05-21** fixing a model-migration startup bug, `2.1.148` and `2.1.149` on **2026-05-22** fixing Bash exit-code `127` regressions and interrupted `/compact` crashes, and `2.1.150` on **2026-05-23** with internal infrastructure improvements.

## GitHub Trending Snapshot

Secondary ecosystem signal only, not roadmap authority. Checked against the default global [GitHub Trending](https://github.com/trending) page on **2026-05-25** with no filters.

- [anthropics/knowledge-work-plugins](https://github.com/anthropics/knowledge-work-plugins): notable because Anthropic-adjacent plugin packaging is surfacing directly on the default board, which reinforces tool and connector packaging as a distribution layer around agents.
- [multica-ai/andrej-karpathy-skills](https://github.com/multica-ai/andrej-karpathy-skills): still relevant because portable behavior overlays remain a visible operator category.
- [dotnet/skills](https://github.com/dotnet/skills): notable because a major platform ecosystem is publishing first-party skill bundles for coding agents.
- [colbymchenry/codegraph](https://github.com/colbymchenry/codegraph): relevant because code-graph and code-intelligence infrastructure continue to surface alongside agent tooling.
- [mukul975/Anthropic-Cybersecurity-Skills](https://github.com/mukul975/Anthropic-Cybersecurity-Skills): notable because domain-specific skill inventories are still attracting broad attention on a general developer board.

## Harness-Relevant Signals

- Anthropic, Google, and OpenAI are converging on the same control-plane shape: isolated execution environments, longer-running sessions, and more explicit operator control over context and tools.
- Anthropic’s required technical surfaces were explicitly checked this run: [Anthropic newsroom](https://www.anthropic.com/news), [Claude blog](https://claude.com/blog), [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code changelog](https://code.claude.com/docs/en/changelog), and [Claude Code desktop changelog](https://code.claude.com/docs/en/desktop-changelog).
- The [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) still runs behind the `managed-agents-2026-04-01` beta header, and the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) still require `advisor-tool-2026-03-01`. Anthropic’s most ambitious agent-control surfaces are still evolving under beta contracts rather than a frozen platform guarantee.
- Google’s official managed-agent launch explicitly normalizes `AGENTS.md` and `SKILL.md` as mainstream agent configuration artifacts, which is strong external validation for repo-local contracts and skill packaging.
- Moonshot, MiniMax, and DeepSeek were each checked once on official surfaces this run. Moonshot’s visible current flagship remains [Kimi K2.6](https://www.moonshot.ai/en) and the [Kimi API platform](https://platform.moonshot.ai/) but no stronger new dated official technical item surfaced inside this run window. MiniMax’s in-window official post was [MiniMax at the WaytoAGI Global AI Conference](https://www.minimax.io/news/minimax-at-the-waytoagi-global-ai-conference) on **2026-05-22**, which reiterated MiniMax Agent and M2.7 positioning but was not a cleaner technical launch. DeepSeek’s latest official API changelog entry remains [DeepSeek-V4](https://api-docs.deepseek.com/updates/) on **2026-04-24**, so there was no new in-window first-party technical release to elevate.

## Candidate Proposals

- `source_type`: `first_party_managed_agent_runtime`
  `problem_pattern`: `managed_agent_products_are_converging_on_isolated_runtimes_resumable_state_and_live_reconfiguration`
  `why_it_matters_for_harness`: Anthropic and Google both moved the hosted-agent runtime again this week, which strengthens the case that this is a durable control-plane category rather than one vendor’s temporary packaging choice.
  `recommended_disposition`: `watch`

- `source_type`: `first_party_agent_operator_surface`
  `problem_pattern`: `coding_agents_are_getting_richer_context_capture_goal_control_and_remote_observability`
  `why_it_matters_for_harness`: OpenAI’s Codex updates and Anthropic’s Claude Code fixes both tighten the operator layer around long-running execution instead of only improving base-model capability.
  `recommended_disposition`: `encode-as-eval`

- `source_type`: `secondary_distribution_signal`
  `problem_pattern`: `plugins_skills_and_domain_specific_capability_packs_are_becoming_operator_facing_distribution_units`
  `why_it_matters_for_harness`: GitHub Trending keeps surfacing plugin packs, skills, and domain inventories, which suggests capability packaging is remaining a durable adoption surface around agent runtimes.
  `recommended_disposition`: `prototype`

## Watch Next

- Anthropic: keep explicitly checking the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and [Claude Code changelog](https://code.claude.com/docs/en/changelog) every run. Managed Agents moved again this week, Claude Code is still changing week to week, and the advisor tool remains a meaningful beta surface even without a new in-window dated launch.
- OpenAI: watch whether the **2026-05-21** Codex release-note cluster turns into a clearer public contract for remote execution governance, browser-level control, and longer-lived task management.
- Google: watch whether the **2026-05-19** Gemini managed-agents launch quickly picks up sharper pricing, policy, and enterprise-governance detail beyond preview framing.
- Qwen: watch whether the **2026-05-20** Qwen3.7-Max release grows into a fuller developer platform or agent-runtime story rather than remaining primarily a model announcement.
- China vendors beyond Qwen: unless Moonshot, MiniMax, or DeepSeek publish a dated technical launch in the next few days, keep them in the scan but do not over-read generic product or conference messaging as runtime maturity.
