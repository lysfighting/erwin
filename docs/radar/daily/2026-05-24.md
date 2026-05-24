# AI Agent Daily Radar - 2026-05-24

Window: 2026-05-17 to 2026-05-24

## New Signals

- **Anthropic expanded Claude Managed Agents on 2026-05-19.** The official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) and matching [Claude blog post](https://claude.com/blog/new-in-claude-managed-agents-self-hosted-sandboxes-and-mcp-tunnels) added research-preview MCP tunnels, self-hosted sandboxes, live MCP and tool reconfiguration for active sessions, and automatic spill-to-file handling for oversized tool output. This is still the strongest first-party managed-runtime move in-window.
- **Anthropic Claude Code shipped multiple operational fixes on 2026-05-21, 2026-05-22, and 2026-05-23.** The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows `2.1.147` on **2026-05-21** fixing a model-migration startup bug, `2.1.148` on **2026-05-22** fixing Bash exit-code `127` regressions for some users, `2.1.149` on **2026-05-22** fixing crashes from interrupted `/compact`, and `2.1.150` on **2026-05-23** with internal infrastructure improvements.
- **Google launched Managed Agents in the Gemini API on 2026-05-19.** In [Build managed agents with the Gemini API](https://blog.google/innovation-and-ai/technology/developers-tools/managed-agents-gemini-api/), Google framed a single-call agent runtime with isolated Linux execution, resumable state, web access, and repo-local agent contracts such as `AGENTS.md` and `SKILL.md`.
- **OpenAI updated Codex surfaces on 2026-05-21.** The official [ChatGPT release notes](https://help.openai.com/en/articles/6825453-chatGPTs-release-notes) added Appshots, general-availability Goal mode, in-app browser annotations, locked computer use, and browser-use improvements for Codex. This is a concrete control-plane update for longer-running coding workflows.
- **Qwen released Qwen3.7-Max on 2026-05-20 as an agent-oriented flagship.** The official [Qwen research feed](https://qwen.ai/research/) positions Qwen3.7-Max as a “versatile agent foundation” for coding, office automation, and long autonomous execution. It is the clearest in-window signal from the China model vendors covered this run.

## GitHub Trending Snapshot

Secondary ecosystem signal only, not roadmap authority. Checked against the default global [GitHub Trending](https://github.com/trending) page on **2026-05-24** with no filters.

- [multica-ai/andrej-karpathy-skills](https://github.com/multica-ai/andrej-karpathy-skills): surfacing because portable agent operating overlays are still a major distribution category.
- [dotnet/skills](https://github.com/dotnet/skills): notable because first-party skill bundles for coding agents are now showing up from major platform ecosystems.
- [ChromeDevTools/chrome-devtools-mcp](https://github.com/ChromeDevTools/chrome-devtools-mcp): important because browser-control infrastructure remains part of the coding-agent stack.
- [mukul975/Anthropic-Cybersecurity-Skills](https://github.com/mukul975/Anthropic-Cybersecurity-Skills): relevant because domain-specific, installable skill packs are attracting operator attention.
- [presenton/presenton](https://github.com/presenton/presenton): notable because AI-native content-generation tooling is still winning broad developer interest on the default board.

## Harness-Relevant Signals

- Anthropic, Google, OpenAI, and Qwen are all reinforcing the same broad shape: longer-running agent execution, stronger runtime scaffolding, and more explicit operator control over agent environments.
- Anthropic’s required technical surfaces were explicitly checked this run: [Anthropic newsroom](https://www.anthropic.com/news), [Claude blog](https://claude.com/blog), [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog).
- The [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) still requires the `managed-agents-2026-04-01` beta header. The [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) still require the `advisor-tool-2026-03-01` beta header. That means Anthropic’s most ambitious agent-control surfaces are still moving under beta contracts rather than a frozen platform guarantee.
- Google’s decision to make `AGENTS.md` and `SKILL.md` part of the mainstream Gemini managed-agent story is strong external validation for repo-local agent contracts and skill packaging.
- Moonshot, MiniMax, and DeepSeek were each checked once on official surfaces this run. Moonshot’s visible latest dated item was still [Kimi K2.6 from 2026-04-20](https://www.moonshot.ai/), MiniMax’s in-window official post was a conference recap rather than a stronger technical launch, and DeepSeek’s latest API changelog item remained [DeepSeek-V4 from 2026-04-24](https://api-docs.deepseek.com/zh-cn/updates).

## Candidate Proposals

- `source_type`: `first_party_managed_agent_runtime`
  `problem_pattern`: `managed_agent_products_are_converging_on_isolated_runtimes_resumable_state_and_live_configuration`
  `why_it_matters_for_harness`: Anthropic and Google both moved the hosted-agent runtime again this week, which strengthens the case that this is a durable control-plane category rather than vendor-specific experimentation.
  `recommended_disposition`: `watch`

- `source_type`: `first_party_agent_operator_surface`
  `problem_pattern`: `coding_agents_are_getting_more_explicit_context_capture_goal_control_and_browser_level_feedback`
  `why_it_matters_for_harness`: OpenAI and Anthropic both sharpened the operator surface rather than only the base model, which is a useful pressure signal for how agent systems are becoming governable in practice.
  `recommended_disposition`: `encode-as-eval`

- `source_type`: `secondary_distribution_signal`
  `problem_pattern`: `skills_and_mcp_style_capability_packs_are_becoming_a_visible_operator_distribution_layer`
  `why_it_matters_for_harness`: GitHub Trending keeps surfacing skills, browser-control tooling, and domain packs, which suggests capability packaging remains a durable adoption layer around runtimes.
  `recommended_disposition`: `prototype`

## Watch Next

- Anthropic: keep explicitly checking the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and [Claude Code changelog](https://code.claude.com/docs/en/changelog) every run. The surfaces are still beta-scoped and moving quickly.
- Anthropic advisor tool: no separate new dated launch landed inside **2026-05-17 to 2026-05-24**, but it remains a first-class watch surface because Anthropic still positions it for long-horizon coding, computer-use, and research workflows.
- Google: watch whether the **2026-05-19** managed-agents launch quickly gains sharper pricing, policy, and enterprise-governance detail beyond preview framing.
- OpenAI: watch whether the **2026-05-21** Codex release-note cluster turns into a clearer public contract for setup policy, remote execution governance, and browser-level control.
- Qwen: watch whether the **2026-05-20** Qwen3.7-Max release grows into a more explicit developer platform or agent-runtime surface rather than remaining primarily a model announcement.
