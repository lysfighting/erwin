# Daily Radar - 2026-07-07

Window: 2026-06-30 to 2026-07-07

## New Signals

1. **Anthropic paired a model launch with a more explicit operator surface.** The official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) on **2026-06-30** launched `claude-sonnet-5` and added session delta events, reverse pagination, per-session agent overrides, vault `injection_location`, and broader webhooks. The official [Claude Apps release notes](https://support.claude.com/en/articles/12138966-release-notes) on **2026-07-01** also recorded restored access to Claude Fable 5 and Claude Mythos 5 plus Enterprise model entitlements.  
   Publish dates: 2026-06-30, 2026-07-01

2. **OpenAI shipped a fresh realtime model update rather than a broader agent-control release.** The official [OpenAI API changelog](https://developers.openai.com/api/docs/changelog) on **2026-07-06** released `gpt-realtime-2.1` and `gpt-realtime-2.1-mini`, emphasizing better alphanumeric recognition, silence/noise handling, and interruption behavior for realtime voice systems.  
   Publish date: 2026-07-06

3. **Google kept pushing the “deterministic workflow plus agent” architecture.** The official [ADK Go 2.0 post](https://developers.googleblog.com/announcing-adk-go-20/) on **2026-06-30** added a graph-based workflow engine, built-in human-in-the-loop, and dynamic orchestration. The official [Why we built ADK 2.0](https://developers.googleblog.com/why-we-built-adk-20/) post on **2026-07-01** made the rationale explicit: production agents need deterministic routing and error handling around probabilistic model nodes.  
   Publish dates: 2026-06-30, 2026-07-01

4. **Moonshot’s newest Kimi Code release focused on subagent reliability and provider correctness.** The official [Kimi Code changelog](https://www.kimi.com/code/docs/en/kimi-code-cli/release-notes/changelog.html) on **2026-07-07** says `0.23.1` fixes `kimi -p` abandoning late or long-running background subagents, restores web chat streaming after stale background-tab WebSockets, and honors explicit Anthropic `max_output_size` settings.  
   Publish date: 2026-07-07

5. **Qwen published a stronger eval artifact instead of another model SKU.** The official [Qwen-Image-Bench](https://www.alibabacloud.com/blog/qwen-image-bench-beyond-basic-generation-%E2%80%94-evaluating-t2i-models-in-complex-scenarios_603335) post on **2026-07-06** introduced a creator-oriented benchmark with 56 verifiable rubrics plus the open-source judge model Q-Judger. This is notable because the benchmark is trying to measure real creative-workflow constraints rather than generic image quality alone.  
   Publish date: 2026-07-06

## GitHub Trending Snapshot

Snapshot source: [GitHub Trending](https://github.com/trending) captured on 2026-07-07.

- [Zackriya-Solutions/meetily](https://github.com/Zackriya-Solutions/meetily): surfacing because local, privacy-first AI meeting tools are still drawing strong attention without relying on cloud control planes.
- [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills): relevant because skills are increasingly being packaged as reusable engineering assets rather than prompt fragments.
- [TencentCloud/CubeSandbox](https://github.com/TencentCloud/CubeSandbox): notable because secure lightweight sandboxes are becoming core infrastructure for AI-agent execution.
- [iOfficeAI/OfficeCLI](https://github.com/iOfficeAI/OfficeCLI): relevant because agent-native Office document automation is solidifying as a distinct tooling layer.
- [steipete/CodexBar](https://github.com/steipete/CodexBar): worth watching because operator-facing visibility for Codex and Claude Code usage is now attracting its own tooling.

## Harness-Relevant Signals

- Anthropic was checked explicitly across the required official surfaces this run: [Anthropic newsroom](https://www.anthropic.com/news), [Claude Apps release notes](https://support.claude.com/en/articles/12138966-release-notes), [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code What's New](https://code.claude.com/docs/en/whats-new), and the [Claude Code changelog](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md).
- The current [Managed Agents quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart) still centers the `managed-agents-2026-04-01` beta header and the `agent -> environment -> session -> events` execution model. I did not find a separately dated in-window quickstart edit, but the June 30 platform notes materially expanded the surrounding session-control surface.
- The current [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) still position advisor as the planner/executor split for long-horizon coding and agent tasks, and now document hard output caps via `tools[].max_tokens`. I did not find a newer dated in-window advisor announcement beyond the June 5 platform release-note entry.
- The official [Claude Code What's New](https://code.claude.com/docs/en/whats-new) page still tops out at Week 26 (`2026-06-22` to `2026-06-26`), so the [Claude Code changelog](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md) remains the fresher source for live product movement.
- OpenAI was checked through official platform and release-note surfaces. The newest clearly dated in-window item was the **2026-07-06** realtime update; I did not find a fresher official Codex operator/governance note inside this window.
- MiniMax and DeepSeek were each checked once through official surfaces. I did not find a strong new dated official technical item inside the 2026-06-30 to 2026-07-07 window, so they stay on watch rather than driving this run.

## Candidate Proposals

### Proposal 1

- `source_type`: first-party platform release notes
- `problem_pattern`: operator surfaces for long-running agents are converging on steered sessions, event deltas, reversible pagination, per-session overrides, and explicit lifecycle hooks
- `why_it_matters_for_harness`: evaluation pressure is shifting from “can the agent finish?” toward “can operators observe, steer, resume, and govern the run cleanly?”
- `recommended_disposition`: `encode-as-eval`

### Proposal 2

- `source_type`: first-party coding-agent changelogs
- `problem_pattern`: subagent durability and provider-correct request handling are now product-surface requirements, not incidental bugfixes
- `why_it_matters_for_harness`: orchestration regressions increasingly show up in background-task handling and model-provider edge behavior before they show up in headline benchmarks
- `recommended_disposition`: `watch`

### Proposal 3

- `source_type`: first-party benchmark/eval release
- `problem_pattern`: vendor-owned evals are moving toward scenario-rich, judge-backed workflow rubrics instead of single-score headline benchmarks
- `why_it_matters_for_harness`: richer external rubric shapes are becoming more useful than raw leaderboard deltas for deciding what deserves durable observation
- `recommended_disposition`: `prototype`

## Watch Next

- Anthropic remains the highest-priority next-run check across the same required surface. The most likely next meaningful change is either a new Managed Agents/session-control note, an advisor-tool follow-up, or a dated Claude Code weekly digest catching up to the raw changelog.
- OpenAI’s next watch point is whether official Codex or API release notes add a fresher operator-surface update after the **2026-07-06** realtime release.
- Google’s next watch point is whether ADK 2.0 gets more production guidance or runtime expansions that sharpen the agent-vs-workflow boundary.
- Moonshot’s next watch point is whether Kimi Code keeps turning agent reliability and plugin/command surfaces into a faster release train than model launches alone.
- MiniMax and DeepSeek stay on watch, but without fresher dated official technical movement they remain comparison signals rather than front-line drivers for `erwin` this run.
