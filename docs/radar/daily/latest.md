# Daily Radar — 2026-06-04

Window: 2026-05-29 to 2026-06-04

## New Signals

1. **Anthropic tightened the advisor-tool economics on June 2, 2026.** The official [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added `tools[].max_tokens` for the [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), which lets operators cap advisor-output length per call, and also removed billing when a Claude API request returns `stop_reason: "refusal"` without model output.  
   Publish date: 2026-06-02

2. **Anthropic kept expanding the managed-agent execution surface on May 29, 2026.** The same [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) say [Claude Managed Agents](https://platform.claude.com/docs/en/managed-agents/quickstart) gained webhooks, multiagent orchestration, and self-hosted sandboxes on Claude Platform on AWS.  
   Publish date: 2026-05-29

3. **Claude Code’s June 2, 2026 changelog focused on operational control rather than headline capability.** The official [Claude Code changelog](https://code.claude.com/docs/en/changelog) added `OTEL_RESOURCE_ATTRIBUTES` labels on metrics, improved fan-out visibility in `claude agents`, and changed parallel tool calls so one failed Bash command no longer cancels sibling calls in the same batch.  
   Publish date: 2026-06-02

4. **OpenAI pushed Codex further toward split control-plane and execution-host operation on May 29, 2026.** OpenAI’s official [ChatGPT release notes](https://help.openai.com/en/articles/6825453-chatgpt-versionshinweise) say Codex now supports Computer Use on Windows plus remote steering from ChatGPT on iOS/Android or Codex on Mac while the Windows machine remains the host for files, shell, app server, and local context.  
   Publish date: 2026-05-29

5. **Qwen shipped a paired coding-agent/runtime refresh on May 29, 2026.** [Qwen Code v0.17.0](https://github.com/QwenLM/qwen-code/releases) added zero-config built-in computer use via MCP plus a new compaction/restore path, while [Qwen-Agent v0.0.26](https://github.com/QwenLM/Qwen-Agent/releases) updated the framework line with MCP SSE timeout and reconnection support.  
   Publish date: 2026-05-29

Google, Moonshot, MiniMax, and DeepSeek did not surface a stronger dated official technical item inside this window than the items above. This run explicitly checked Google’s developer/blog/docs surfaces, Moonshot’s platform/blog surfaces, MiniMax’s official blog, and DeepSeek’s API docs/news surfaces; none produced a dated 2026-05-29 to 2026-06-04 item that materially changed the current radar.

## GitHub Trending Snapshot

Snapshot source: [GitHub Trending](https://github.com/trending) captured on 2026-06-04.

- [chopratejas/headroom](https://github.com/chopratejas/headroom): context compression is still surfacing as a first-order agent cost/control problem rather than a niche optimization.
- [affaan-m/ECC](https://github.com/affaan-m/ECC): explicit harness, memory, and workflow language keeps showing open demand for agent operating systems.
- [NousResearch/hermes-agent](https://github.com/NousResearch/hermes-agent): agent runtime stacks remain a live open-source attention cluster, not just productized vendor territory.
- [microsoft/markitdown](https://github.com/microsoft/markitdown): document-to-Markdown normalization is still core ingestion plumbing for coding and research agents.
- [aquasecurity/trivy](https://github.com/aquasecurity/trivy): security tooling is surfacing alongside agent repos, which is a useful reminder that runtime trust and scanning are staying coupled to deployment reality.

## Harness-Relevant Signals

- Anthropic is hardening three adjacent layers at once: advisor economics, managed-agent execution, and coding-agent control-plane ergonomics.
- Claude Code’s OTEL labels and independent parallel-tool failure handling raise the quality bar for observable, long-running, fanned-out work.
- OpenAI’s Windows-host plus remote-steering model reinforces the split between the machine doing the work and the interface supervising it.
- Qwen’s May 29 release continues the convergence of coding agent, MCP-mediated computer use, and worktree/session ergonomics.
- Trending pressure still clusters around compression, ingestion, runtime stacks, and security, which are practical prerequisites for durable agent systems rather than side concerns.

## Candidate Proposals

### Proposal 1

- `source_type`: first-party platform release notes and coding-agent changelogs
- `problem_pattern`: long-running agent systems increasingly separate planner, executor, supervisor, and host machine concerns
- `why_it_matters_for_harness`: Anthropic and OpenAI are both making control-plane versus execution-plane boundaries more explicit, which should shape how Harness evaluates continuity and operator intervention assumptions
- `recommended_disposition`: `watch`

### Proposal 2

- `source_type`: first-party coding-agent changelogs
- `problem_pattern`: observability and failure isolation are moving into the default agent runtime rather than living in custom wrappers
- `why_it_matters_for_harness`: Claude Code’s OTEL labeling and per-tool failure isolation suggest that durable fan-out work now expects built-in instrumentation and less fragile batch semantics
- `recommended_disposition`: `encode-as-eval`

### Proposal 3

- `source_type`: first-party managed-agent and open-source coding-agent releases
- `problem_pattern`: managed-agent stacks and open coding agents are converging on MCP plus computer-use mediated execution
- `why_it_matters_for_harness`: Anthropic’s managed-agent expansion and Qwen’s built-in computer-use path are useful comparison points for deciding where Harness should watch versus prototype execution-boundary patterns
- `recommended_disposition`: `prototype`

## Watch Next

- Anthropic: next run must explicitly recheck [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Claude Managed Agents docs](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). This run found meaningful in-window release-note and changelog updates, but no separately dated doc-only Managed Agents or advisor-tool change stronger than those entries.
- OpenAI: watch for narrower developer-facing documentation that turns the May 29 remote-control posture into a clearer host/control contract.
- Google: keep monitoring for a dated follow-up that turns prior managed-agent direction into a fresher operational or release-note artifact.
- Moonshot, MiniMax, and DeepSeek: continue lightweight monitoring; each surface was checked this run, but none produced a stronger in-window dated technical item than Anthropic, OpenAI, or Qwen.
