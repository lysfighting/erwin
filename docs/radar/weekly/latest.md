# AI/Agent Weekly Synthesis — 2026-05-15

Coverage window: 2026-05-08 to 2026-05-15

## Field Summary

- Anthropic again had the clearest multi-surface movement, and the required Anthropic surfaces were explicitly checked before concluding scope: [anthropic.com/news](https://www.anthropic.com/news), [claude.com](https://claude.com/), [docs.anthropic.com](https://docs.anthropic.com/en/home), the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), the [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), the [Claude Code docs overview](https://code.claude.com/docs/en/overview), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). This week did have important Anthropic platform changes.
- On **May 11, 2026**, Anthropic’s [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added **Claude Platform on AWS**, exposing Messages, Files, Message Batches, Claude Managed Agents, Agent Skills, code execution, and tool use through Anthropic-managed AWS endpoints with AWS billing and IAM authentication. On **May 12, 2026**, the same surface added **fast mode support for Claude Opus 4.7**. These are not minor SDK notes; they extend the provider’s deployment and latency-control surfaces.
- Anthropic’s hosted-agent path remains materially active. The same release-notes surface still shows the **May 6, 2026** launch of **Managed Agent multiagent sessions, outcomes, webhooks, vault refresh, and richer session/event filtering**, and the [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) now clearly presents agent, environment, session, and events as the core runtime contract. The [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) remain live, but I did not find a newer in-window dated advisor launch after the **April 9, 2026** beta; the current doc still centers executor/advisor pairings around Opus 4.7.
- Anthropic’s Claude Code stream changed materially again across **May 8-14, 2026**. The [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows `2.1.136` on **May 8**, `2.1.137` and `2.1.138` on **May 9**, `2.1.139` on **May 11**, `2.1.140` on **May 12**, `2.1.141` on **May 13**, and `2.1.142` on **May 14**. The largest shift is `2.1.139` and later: `claude agents` view, `/goal`, agent-aware headers and OTel attributes, plugin detail/provenance surfaces, background-session controls, and new dispatched-session flags such as `--model`, `--effort`, `--permission-mode`, and `--mcp-config`.
- Anthropic’s product/news layer also moved. The [Anthropic newsroom](https://www.anthropic.com/news) shows **Introducing Claude for Small Business** on **May 13, 2026** and a **$200 million Gates Foundation partnership** on **May 14, 2026**. The current [claude.com](https://claude.com/) product surface also now presents Claude Code, Cowork, desktop extensions, Slack/Google Workspace connectors, and remote MCP-style integration as first-class plan value rather than side features. I treat that as supporting evidence that Anthropic’s app and platform surfaces are converging around “Claude as governed work environment,” even though the main dated technical changes this week came from Platform and Claude Code.
- OpenAI’s strongest first-party signal this week was not a new base model. It was a new **Codex operating surface**. On **May 14, 2026**, [Work with Codex from anywhere](https://openai.com/index/work-with-codex-from-anywhere/) brought Codex into the ChatGPT mobile app, paired it with Remote SSH general availability, and made **Hooks** generally available while adding programmatic access tokens and HIPAA-compliant local use for eligible Enterprise workspaces. OpenAI also published [Building a safe, effective sandbox to enable Codex on Windows](https://openai.com/index/building-codex-windows-sandbox/) on **May 13, 2026** and [Running Codex safely at OpenAI](https://openai.com/index/running-codex-safely/) on **May 8, 2026**, which together make the control model unusually explicit: sandbox boundaries, managed network policy, keyring-backed credentials, requirements, rules, and OpenTelemetry-backed agent telemetry.
- GitHub’s official Copilot stream continued moving from “assistant” toward **automation control plane**. On **May 13, 2026**, the [Agent tasks REST API](https://github.blog/changelog/2026-05-13-start-copilot-cloud-agent-tasks-via-the-rest-api/) moved Copilot cloud agent task dispatch into automations. On **May 14, 2026**, the [GitHub Copilot app technical preview](https://github.blog/changelog/2026-05-14-github-copilot-app-is-now-available-in-technical-preview/) introduced GitHub-native focused sessions with isolated branches, files, conversation, and task state. These build directly on last week’s [enterprise-managed Copilot CLI plugins](https://github.blog/changelog/2026-05-06-enterprise-managed-plugins-in-github-copilot-cli-are-now-in-public-preview/) and [Rubber Duck multi-model review expansion](https://github.blog/changelog/2026-05-07-rubber-duck-in-github-copilot-cli-now-supports-more-models/).
- Google’s first-party developer signal was not a brand-new agent shell, but it did contribute two concrete implementation patterns this week. On **May 12, 2026**, [Build Long-running AI agents that pause, resume, and never lose context with ADK](https://developers.googleblog.com/build-long-running-ai-agents-that-pause-resume-and-never-lose-context-with-adk/) framed durable state, event-driven dormancy gates, and multi-agent delegation as the difference between demos and production agents. On **May 14, 2026**, [Announcing Genkit Middleware](https://developers.googleblog.com/announcing-genkit-middleware-intercept-extend-and-harden-your-agentic-apps/) made generation-, model-, and tool-layer middleware into a first-class governance mechanism for retries, fallbacks, tool approval, skills loading, and scoped filesystem access.
- The week’s broad field judgment is that vendors are externalizing the **agent operating contract**. The important primitives are now background-session dispatch, mobile supervision, event/webhook surfaces, durable pause/resume, environment policy, plugin provenance, and traceable agent identity. Capability is still improving, but the center of gravity is clearly shifting toward governed execution.

## Theme Deep Dive

This week’s theme is **agent work becoming a managed session fabric rather than a single foreground conversation**.

Anthropic’s current surface is the cleanest proof. Claude Platform now spans Anthropic-managed infrastructure on AWS, a hosted Managed Agents harness, a separate advisor lane, and a fast-mode latency tier for Opus 4.7. The [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) makes the design explicit: agents, environments, sessions, and events are separate runtime objects. That is closer to an operating substrate than to “API plus prompt.”

Claude Code pushes the same direction locally. The [May 11-14 Claude Code releases](https://code.claude.com/docs/en/changelog) add agent view, background-session configuration flags, `/goal`, session scoping, better plugin/component introspection, agent-aware telemetry headers, and more reliable background behavior across sleep/wake, worktrees, and remote MCP. This matters because it reduces the conceptual gap between local coding agents and hosted agent runtimes. Anthropic is not just improving answers; it is standardizing how agent sessions are created, configured, tracked, and resumed.

OpenAI and GitHub are converging from different angles. OpenAI’s Codex mobile surface and Windows sandbox/security writeups show a platform trying to make long-running coding agents operationally normal: a user can start or steer work from a phone, keep execution on trusted machines, and apply explicit sandbox/network/identity policy. GitHub is making similar moves at the repo-control layer: cloud-agent tasks can start from an API, sessions live in isolated task spaces, and enterprise admins control plugin rollout. The implementation details differ, but the control-plane direction is aligned.

Google’s ADK and Genkit material is weaker as product launch news, but strong as architectural corroboration. Google’s durable-agent tutorial says production agents need explicit state and event-driven dormancy rather than giant replayed chat logs. Genkit middleware says the agent loop itself is now where governance belongs: retries, fallbacks, tool approval, filesystem boundaries, and skills injection should be expressed as composable control surfaces. That is the same problem Anthropic, OpenAI, and GitHub are solving through different product packaging.

For `erwin`, the important comparison point is no longer “who has the strongest model” or even “who has agents.” It is: who exposes the cleanest contract for **session lifecycle, environment policy, human interruption, plugin/skill provenance, telemetry, and resumability**.

## Tooling Watch

- Anthropic Claude Platform on AWS is the week’s most important provider-deployment change. It puts Anthropic-managed agent/runtime surfaces behind AWS billing and IAM, which is a serious enterprise adoption lever rather than a documentation footnote.
- Claude Managed Agents remain high-velocity. Even though the biggest dated Managed Agent launch in the checked surfaces was **May 6**, it is still current enough to affect this run because this week’s AWS launch explicitly carries Managed Agents forward as part of the main deployment surface.
- Claude Code’s `claude agents` surface looks like a meaningful inflection point. The combination of agent list view, per-dispatch configuration flags, `/goal`, scoped session listing, background-state durability, and agent-aware OTel headers makes Claude Code feel less like “one REPL with extras” and more like a local session scheduler.
- OpenAI’s Codex surface now spans mobile supervision, remote environments, hooks, programmatic access tokens, Windows sandboxing, and explicit security policy examples. The important detail is not just feature count; it is that OpenAI is publishing the operating contract in unusual detail.
- GitHub’s cloud-agent task API and Copilot app show a split control plane: REST/API-driven fan-out on one side, focused human review sessions on the other. This is a notable pattern for repo-native agent systems.
- Google’s Genkit middleware is a strong framework-level signal because it formalizes enforcement points around generation, model calls, and tool execution. That is a better fit for reusable governance than burying policy inside prompts.
- Google’s ADK tutorial is worth watching even though it is instructional content. It explicitly promotes durable state machines, event-driven sleeping/waking, and sub-agent delegation, which are precisely the patterns moving from “advanced idea” to baseline expectation elsewhere.

## Implications For erwin / Harness / Skills

- `erwin` radar should continue splitting vendor families into substreams when the execution surfaces differ. Anthropic this week is not one story; it is at least four: platform deployment, hosted agents, local coding runtime, and app/business packaging.
- `Harness` should treat session lifecycle as a first-class contract surface. The field is converging on named sessions, resumability, background execution, and explicit state transitions, not just single-turn task execution.
- `Harness` should explicitly model environment policy. OpenAI’s Codex materials make sandbox, network rules, identity, and OTel policy concrete, while Anthropic and Google make environment selection and middleware/tool controls concrete. A future-proof harness cannot leave these as ad hoc wrapper concerns.
- `Harness` should preserve agent identity across delegations and subagents. Both Anthropic and Claude Code are getting more explicit about agent/session boundaries, and Claude Code now emits agent-specific request headers and telemetry fields.
- `Skills` should be treated as policy-bearing runtime assets, not only prompt snippets. Claude Code and GitHub are both strengthening plugin and skill inventory visibility, while Google’s middleware even exposes skills loading as an inspectable control surface.
- `erwin` knowledge artifacts should keep distinguishing between hosted-agent orchestration and local-agent orchestration. The field is not collapsing toward one universal shape. Instead, vendors are aligning on control concepts while keeping different execution placements.

## Harness Evolution Proposals

### Title

Provider-neutral session lifecycle contract

### Source

[Claude Code changelog](https://code.claude.com/docs/en/changelog) on **May 11-14, 2026** and [GitHub Copilot app technical preview](https://github.blog/changelog/2026-05-14-github-copilot-app-is-now-available-in-technical-preview/) on **May 14, 2026**

### Observed Pattern

Vendors are exposing explicit agent/session objects with pause, resume, scoped listing, isolated task state, and background completion behavior.

### Why It Matters For Harness

Harness already coordinates work, but its durable contract is still stronger on command surfaces than on session semantics. The field is making lifecycle state user-visible and machine-addressable.

### Recommended Disposition

`adopt`

### What We Would Change If Adopted

Define first-class run/session states, interruption semantics, resumability rules, and a stable session inventory view that works across local and delegated work.

### Why We Should Not Cargo-Cult This

We should not copy any vendor’s exact UX or command names. The valuable part is the lifecycle model, not whether the surface is called `agents`, `sessions`, or `cloud agent`.

### Title

Environment policy as executable config

### Source

[Running Codex safely at OpenAI](https://openai.com/index/running-codex-safely/) on **May 8, 2026** and [Building a safe, effective sandbox to enable Codex on Windows](https://openai.com/index/building-codex-windows-sandbox/) on **May 13, 2026**

### Observed Pattern

OpenAI is publishing concrete sandbox modes, network policy, identity controls, rules, and OTel logging as normal Codex operating surfaces rather than as hidden internal implementation.

### Why It Matters For Harness

Harness will be harder to trust or standardize if permission and environment policy remain scattered across wrappers, prompts, and operator habit.

### Recommended Disposition

`prototype`

### What We Would Change If Adopted

Create a provider-neutral environment-policy schema covering writable roots, allowed network classes, credential scope, approval reviewers, and audit export settings.

### Why We Should Not Cargo-Cult This

OpenAI’s TOML and exact policy names are product-specific. Harness should copy the separation of concerns, not the configuration format.

### Title

Evented long-running runs with dormancy gates

### Source

[Build Long-running AI agents that pause, resume, and never lose context with ADK](https://developers.googleblog.com/build-long-running-ai-agents-that-pause-resume-and-never-lose-context-with-adk/) on **May 12, 2026**

### Observed Pattern

Google frames durable state, idle-time sleeping, and event-driven wakeup as the foundation for real long-running workflows.

### Why It Matters For Harness

Long-running work today easily degrades into oversized transcripts, polling, or fragile human follow-up. A dormancy/wakeup model would make multi-day work more legible and cheaper.

### Recommended Disposition

`prototype`

### What We Would Change If Adopted

Add explicit dormancy checkpoints, wake conditions, and structured state snapshots so long runs can sleep without replaying full chat history.

### Why We Should Not Cargo-Cult This

The onboarding tutorial is illustrative, not a proof that Google’s exact storage model fits Harness. We want the pattern of durable state transitions, not the sample workflow.

### Title

Middleware-style enforcement layers

### Source

[Announcing Genkit Middleware](https://developers.googleblog.com/announcing-genkit-middleware-intercept-extend-and-harden-your-agentic-apps/) on **May 14, 2026**

### Observed Pattern

Google is making generate-, model-, and tool-layer interception first-class for retries, fallbacks, tool approval, skills loading, and scoped filesystem access.

### Why It Matters For Harness

Harness increasingly needs reusable enforcement points that sit above provider APIs but below task prompts.

### Recommended Disposition

`adopt-directionally`

### What We Would Change If Adopted

Define a small set of hookable control layers for run start, model call, tool call, tool result, and run completion, each with stable metadata and optional blocking/approval behavior.

### Why We Should Not Cargo-Cult This

Genkit is a framework with its own abstractions. Harness should preserve the layering idea without inheriting framework-specific semantics or naming.

### Title

Plugin and skill provenance ledger

### Source

[Claude Code changelog](https://code.claude.com/docs/en/changelog) on **May 11-14, 2026** and [Enterprise-managed plugins in GitHub Copilot CLI are now in public preview](https://github.blog/changelog/2026-05-06-enterprise-managed-plugins-in-github-copilot-cli-are-now-in-public-preview/) on **May 6, 2026**

### Observed Pattern

Vendors are making plugin/skill inventory, MCP attachments, token-cost surfaces, and enterprise rollout standards more visible and more governable.

### Why It Matters For Harness

If a run depends on installed skills, plugins, or MCP endpoints, Harness should be able to explain exactly which capability modules were in play and under whose policy.

### Recommended Disposition

`adopt`

### What We Would Change If Adopted

Record skill/plugin/module provenance in run metadata: source, revision, install scope, enabled hooks, MCP endpoints, and declared risk level.

### Why We Should Not Cargo-Cult This

The goal is not to mirror GitHub or Anthropic plugin catalogs. The goal is to make capability dependencies inspectable and auditable.

## Open Questions

- Anthropic now has a stronger hosted-agent deployment story on AWS and a richer local-agent control plane in Claude Code. Does that eventually collapse into one coherent cross-surface contract, or do hosted and local Anthropic agents keep diverging operationally?
- The required advisor-tool docs were explicitly checked, but there was no newer in-window advisor release note after the **April 9, 2026** beta in the surfaces reviewed. Does advisor remain a stable specialist lane, or will Anthropic subsume it into Managed Agents and multiagent orchestration?
- OpenAI’s Codex mobile, Remote SSH, hooks, sandboxing, and telemetry materials imply that “steering from anywhere” is becoming normal. How much of Harness should assume supervision is asynchronous and mobile rather than desktop-bound?
- GitHub is making agent tasks API-addressable while also introducing a desktop-style Copilot app. Will the durable enterprise differentiator be repo-native automation entrypoints, or richer long-running session management?
- Google’s ADK and Genkit content is architectural rather than a single new flagship product launch. Does that signal a strategic bet on reusable agent infrastructure components rather than a monolithic end-user agent shell?
