# AI Agent Daily Radar - 2026-05-15

Window: 2026-05-08 to 2026-05-15

## New Signals

- **Anthropic launched Claude Platform on AWS on 2026-05-11.** The official [Claude Platform on AWS post](https://claude.com/blog/claude-platform-on-aws) and the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) position AWS as a first-class access layer for Messages, Files, Message Batches, Claude Managed Agents, Agent Skills, code execution, and the advisor tool. This is the clearest in-window managed-agent distribution signal from any first-party publisher in this run.
- **Anthropic expanded fast-mode controls for high-end agent workloads on 2026-05-12.** The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added fast mode support for `claude-opus-4-7` under the `fast-mode-2026-02-01` beta header on **May 12, 2026**. That is a concrete operator knob for long-horizon runs, not just a model-quality claim.
- **Claude Code shipped a stronger supervisor surface on 2026-05-11 and 2026-05-12.** The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows `claude agents` and the new [agent view](https://code.claude.com/docs/en/agent-view) plus `/goal` on **May 11, 2026**, followed by **May 12, 2026** fixes for background-service stability, `/goal`, managed settings retries, and advisor-tool UX/processing failures. Claude Code is moving further toward a managed operator console instead of a thin coding shell.
- **OpenAI tightened runtime controls on 2026-05-11 and forced one migration on 2026-05-12.** The official [OpenAI API changelog](https://developers.openai.com/changelog) added `return_token_budget` for the Responses API web search tool on **May 11, 2026**, then removed the Realtime API Beta on **May 12, 2026**. The important pattern is explicit budget control plus faster retirement of transitional agent surfaces.
- **Qwen Code kept shipping agent-runtime plumbing on 2026-05-13.** The official [Qwen Code releases](https://github.com/QwenLM/qwen-code/releases) show stable `v0.15.11` on **May 13, 2026** with `runtime.json` sidecars for active sessions, trace/span correlation for telemetry, headless `--json-schema` output, tool-search controls, and a codegraph PR-risk skill. This remains one of the fastest-moving open coding-agent runtime surfaces in the window.

## GitHub Trending Snapshot

Secondary ecosystem demand only, not first-party roadmap authority. Checked against the default global [GitHub Trending](https://github.com/trending) surface on **2026-05-15** with no filters.

- [rohitg00/agentmemory](https://github.com/rohitg00/agentmemory): persistent memory for coding agents, surfacing because durable state remains a live infrastructure gap.
- [obra/superpowers](https://github.com/obra/superpowers): agentic skills framework and methodology, relevant because skills-as-operator-surface keeps appearing across both vendor and open ecosystems.
- [K-Dense-AI/scientific-agent-skills](https://github.com/K-Dense-AI/scientific-agent-skills): domain-specific agent skills pack, showing continued demand for reusable expert workflows instead of generic chat.
- [CloakHQ/CloakBrowser](https://github.com/CloakHQ/CloakBrowser): stealth Chromium automation, relevant because browser execution remains a contested foundation layer for agents.
- [github/spec-kit](https://github.com/github/spec-kit): spec-driven development toolkit, notable because developer workflows are absorbing more structured AI coordination and verification patterns.

## Harness-Relevant Signals

- Anthropic's required official technical surfaces were explicitly checked before calling the lane quiet: [Anthropic News](https://www.anthropic.com/news), [Claude blog](https://claude.com/blog), [Anthropic docs release notes](https://docs.anthropic.com/en/release-notes/overview), [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code docs](https://code.claude.com/docs), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog).
- As checked on **2026-05-15**, the [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) still says the surface is in beta and requires the `managed-agents-2026-04-01` header. Outcomes and multiagent are still framed as beta or research preview, so the AWS launch is broader distribution of a still-not-fully-settled runtime contract.
- As checked on **2026-05-15**, the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) still require the `advisor-tool-2026-03-01` beta header, but now explicitly list [Claude Platform on AWS](https://claude.com/blog/claude-platform-on-aws) as an available deployment surface. The in-window change is platform reach, not a new advisor API shape.
- The [OpenAI API changelog](https://developers.openai.com/changelog) continues to reinforce a control-surface pattern: `return_token_budget` for high-effort web-search runs on **May 11, 2026**, alongside a hard removal of the Realtime API Beta on **May 12, 2026**. This makes runtime governance and migration discipline more important than just model capability deltas.
- A single official pass was made for Google, Moonshot, MiniMax, and DeepSeek. No strong clearly dated technical launch inside **2026-05-08 to 2026-05-15** surfaced from those first-party pages. The nearest official items remain Google's [Gemini API webhooks](https://blog.google/innovation-and-ai/technology/developers-tools/event-driven-webhooks/) on **May 4, 2026** outside the window, Moonshot's homepage showing [Kimi K2.6](https://www.moonshot.ai/) dated **2026-04-20**, MiniMax's [model release notes](https://platform.minimax.io/docs/release-notes/models) still led by **Apr. 2026**, and DeepSeek's [Transparency Center](https://www.deepseek.com/en/transparency/) still showing **DeepSeek-V4** on **April 24, 2026** as the latest dated model release.
- Secondary signal checks were kept subordinate to first-party surfaces this run. The current [GitHub Trending](https://github.com/trending) board still validates memory, skills, browser control, and spec/verification tooling as the open-source layers rising alongside Anthropic, OpenAI, and Qwen runtime changes.

## Candidate Proposals

- `source_type`: `first_party_platform_launch`
  `problem_pattern`: `managed_agent_capabilities_are_being_repackaged_across_trust_and_billing_boundaries_without_changing_the_core_agent_contract`
  `why_it_matters_for_harness`: Claude Platform on AWS suggests that deployment boundary portability is becoming a first-class product property for agent runtimes.
  `recommended_disposition`: `watch`

- `source_type`: `first_party_runtime_controls`
  `problem_pattern`: `agent_operators_need_explicit_knobs_for_speed_budget_and_migration_not_just_better_models`
  `why_it_matters_for_harness`: Anthropic fast mode and OpenAI `return_token_budget` both point toward evaluable operator controls that sit above raw model quality.
  `recommended_disposition`: `encode-as-eval`

- `source_type`: `first_party_coding_agent_release`
  `problem_pattern`: `coding_agent_products_are_converging_on_supervision_views_goal_conditions_runtime_sidecars_and_traceable_background_work`
  `why_it_matters_for_harness`: Claude Code and Qwen Code are both adding durable execution-state and supervision primitives that look like stable operator expectations.
  `recommended_disposition`: `prototype`

## Watch Next

- Anthropic: keep explicitly checking the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog) every run. As of **2026-05-15**, distribution is widening faster than the managed-agent and advisor surfaces are exiting beta.
- Claude Code: watch whether `claude agents`, [agent view](https://code.claude.com/docs/en/agent-view), and `/goal` are followed by more lifecycle, routing, or cross-session policy controls after the **May 11-12, 2026** release cluster.
- OpenAI: watch whether `return_token_budget` expands into broader budgeting and governance controls for long-running Responses API workloads after **May 11, 2026**.
- Qwen: watch whether `v0.15.11` stabilizes around telemetry, sidecars, and tool-search controls or keeps churning through rapid preview/stable deltas.
- Google, Moonshot, MiniMax, and DeepSeek: recheck next run for a genuinely dated official technical item instead of carrying forward out-of-window or infrastructure-only movement.
