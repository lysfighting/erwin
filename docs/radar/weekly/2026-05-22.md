# AI/Agent Weekly Synthesis — 2026-05-22

Coverage window: 2026-05-15 to 2026-05-22

## Field Summary

- Anthropic had the clearest in-window platform movement again, and this run explicitly checked all required Anthropic surfaces before drawing that conclusion: [anthropic.com/news](https://www.anthropic.com/news), [claude.com](https://claude.com/), [docs.anthropic.com](https://docs.anthropic.com/en/home), the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), the [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), the [Claude Code docs overview](https://code.claude.com/docs/en/overview), and the [Claude Code changelog](https://code.claude.com/docs/en/changelog). This week did have important Anthropic platform and Claude Code changes.
- On **May 19, 2026**, the [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) added four material managed-runtime upgrades at once: **MCP tunnels** in research preview for private-network servers, **self-hosted sandboxes** for Claude Managed Agents, **mid-session updates** to an agent’s MCP server and tool configuration, and **automatic large-output spillover to sandbox files** when `agent_toolset` or MCP outputs exceed 100K tokens. On **May 18, 2026**, the same surface also upgraded the **web search tool** with richer SEC filing data for grounded financial research.
- Anthropic’s hosted-agent contract is now sharper than last week, not fuzzier. The current [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) explicitly defines **agent**, **environment**, **session**, and **events** as the runtime core, and now includes both Anthropic-managed cloud containers and **self-hosted sandboxes** as first-class environment choices. The managed surface is clearly moving toward enterprise-grade runtime configurability rather than a fixed hosted demo loop.
- The required [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) were explicitly checked again. I did **not** find a new in-window advisor-specific release note after the beta state already documented there, but the page remains important because it still formalizes a distinct **executor plus advisor** pattern, with **Opus 4.7** as the advisor model, eligibility on the Claude API and Claude Platform on AWS, and a separate beta header. That matters because Anthropic still treats planning-quality escalation as a tool-level primitive, not just a model picker.
- Anthropic’s Claude Code stream changed materially during the window. The [Claude Code changelog](https://code.claude.com/docs/en/changelog) shows `2.1.147` on **May 21, 2026** and `2.1.148` on **May 22, 2026**. The meaningful change is `2.1.147`: **pinned background sessions** in `claude agents` now stay alive when idle, survive updates in place, and are shed only after non-pinned sessions under memory pressure; `/simplify` was renamed to **`/code-review`** with correctness-focused review output and inline GitHub comment support; and enterprise login restrictions were tightened for managed settings. `2.1.148` is smaller but still relevant because it fixes a Bash regression that broke command execution for some users.
- Anthropic’s app-side surface also moved. The [Claude Help Center release notes](https://support.claude.com/en/articles/12138966-release-notes) added **Claude Compliance API integrations** on **May 21, 2026**, explicitly positioning Claude governance across the platform and product suite. The [claude.com](https://claude.com/) product surface now foregrounds Claude Code, Cowork, Security, Skills, platform docs, and Claude-on-AWS links in one product lattice rather than as separate side paths. That supports the interpretation that Anthropic is converging app, platform, and local-runtime governance.
- OpenAI’s official in-window signal was weak by comparison. In the `openai.com` pass, I did not find a new post between **May 15 and May 22, 2026** that materially changed the Codex or Agents SDK operating contract. The nearest relevant official baseline is still [Work with Codex from anywhere](https://openai.com/index/work-with-codex-from-anywhere/) from **May 14, 2026**, which remains useful context for mobile supervision, hooks, programmatic access tokens, and Remote SSH, but it is outside this run’s strict weekly window.
- GitHub’s first-party Copilot changelog was unusually dense and aligned with Anthropic’s direction. On **May 18, 2026**, GitHub shipped [fast, cost-efficient models for Copilot cloud agent](https://github.blog/changelog/2026-05-18-copilot-cloud-agent-fast-cost-efficient-models-for-simple-tasks/) and [one-click fixes for failing Actions with Copilot cloud agent](https://github.blog/changelog/2026-05-18-one-click-fixes-for-failing-actions-with-copilot-cloud-agent/). On **May 19**, it added [Fix with Copilot / Fix batch with Copilot](https://github.blog/changelog/2026-05-19-easily-apply-copilot-code-review-feedback-with-copilot-cloud-agent/) for PR review feedback. On **May 20**, it added [semantic issue search in Copilot Chat](https://github.blog/changelog/2026-05-20-semantic-issue-search-in-copilot-chat/) and [task-aware auto model selection in VS Code](https://github.blog/changelog/2026-05-20-auto-model-selection-now-routes-based-on-your-task-in-vs-code/). On **May 21**, it made [GitHub Copilot for Eclipse open source](https://github.blog/changelog/2026-05-21-github-copilot-for-eclipse-is-open-source/), including visibility into agent mode, prompts, and skills loading.
- Google’s strongest in-window first-party agent signal came on **May 21, 2026** with [ADK for Kotlin and ADK for Android 0.1.0](https://developers.googleblog.com/adk-kotlin-android-building-ai-agents/). The important part is not just another SDK binding. Google is pushing agent execution toward **on-device**, **session-state-aware**, **OpenTelemetry-capable**, **tool- and plugin-enabled** agent runtimes that can bridge local and cloud execution.
- The broad weekly judgment is that frontier vendors are tightening the **runtime control plane** around agents: private-network tool access, self-hosted execution, pinned/background sessions, review-to-fix handoff, task-aware model routing, and more explicit governance APIs. The center of gravity is still shifting from “can the agent do the task?” to “how does the agent stay governable while it does the task?”

## Theme Deep Dive

This week’s theme is **governed agent runtime portability**.

Anthropic is the clearest example because the changes span every layer of the stack. The [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) turned Managed Agents from a hosted convenience into something closer to a portable runtime contract: MCP tunnels let agents reach private-network systems, self-hosted sandboxes let the execution environment move onto customer-controlled infrastructure, and mid-session MCP/tool reconfiguration makes the running session less rigid. The [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) reinforces that this is not accidental. Anthropic now presents agents, environments, sessions, and events as separable runtime objects.

Claude Code is moving in the same direction from the local side. The [May 21-22 changelog entries](https://code.claude.com/docs/en/changelog) make background and pinned sessions more durable, improve session-level review behavior, and keep enterprise policy enforcement visible. That matters because it narrows the gap between “local coding agent” and “hosted managed agent.” Anthropic is increasingly treating both as variants of the same operational problem: how to keep long-running agent work resumable, inspectable, and policy-bound.

GitHub is corroborating the pattern from a repo-native angle. Copilot cloud agent now covers CI failure repair, PR review feedback application, model choice for simple versus complex work, and issue discovery. That is not just assistant sprawl. It is a tighter session pipeline from triage to change to review. The [open-source Eclipse plugin](https://github.blog/changelog/2026-05-21-github-copilot-for-eclipse-is-open-source/) is especially notable because it exposes how agent mode, prompts, and skills are wired inside a real IDE integration, which makes governance and implementation choices inspectable.

Google’s [ADK for Kotlin and Android](https://developers.googleblog.com/adk-kotlin-android-building-ai-agents/) broadens the portability story. Instead of treating mobile as a thin client over cloud agents, Google is pushing toward agent runtimes that can live closer to the device, keep session state, integrate plugins and MCP tools, and still emit OpenTelemetry. That suggests portability is becoming two-dimensional: not just cloud provider choice, but also **where the agent loop can physically run**.

For `erwin`, the important lesson is that “agent platform” is becoming less about a single model or IDE and more about a stable contract over:

- session lifecycle
- environment placement
- tool and plugin provenance
- human steering points
- audit and telemetry surfaces
- policy enforcement and compliance hooks

That is the layer where vendors are now differentiating with real product work.

## Tooling Watch

- Anthropic’s **self-hosted sandboxes for Claude Managed Agents** are the highest-signal technical change of the week because they move a previously hosted execution model onto customer infrastructure without abandoning Anthropic’s higher-level agent contract.
- Anthropic’s **MCP tunnels** matter because they turn “private tools and private networks” into a first-class integration problem rather than an implicit enterprise workaround.
- Anthropic’s **pinned background sessions** in Claude Code are small on paper but strategically important. Once sessions can stay alive when idle, survive updates in place, and be deprioritized only after non-pinned sessions, the local tool starts behaving more like an operating substrate.
- The required **advisor tool** check did not produce a fresh in-window launch, but its continued separate existence still matters. Anthropic is maintaining a distinct “cheap executor plus strong planner” primitive while the rest of the stack becomes more runtime-oriented.
- GitHub’s cloud-agent changes are converging on a tight operator loop: choose the right model for the task, start from CI or review context, batch feedback, and let the agent patch the branch from its own environment.
- GitHub’s **open-source Eclipse plugin** is more than ecosystem goodwill. It exposes prompt, skill, and agent-mode internals that most vendors still keep opaque.
- Google’s **ADK for Android** is a concrete signal that agent runtimes are moving closer to the edge. On-device execution, local privacy constraints, and mixed cloud/local agent flows are becoming normal design inputs.
- OpenAI had no comparably new in-window agent-platform post in the official pass. That absence is itself a signal: this week’s field movement came more from Anthropic’s runtime surface and GitHub’s repo-native cloud-agent workflow than from a new OpenAI control-plane release.

## Implications For erwin / Harness / Skills

- `erwin` radar should treat **runtime portability** as its own active lens, separate from generic “agent capability” coverage. The interesting field movement is in how vendors preserve the same agent contract across hosted, self-hosted, desktop, and edge environments.
- `Harness` should make **environment placement** explicit. “Where this run executes” is now a first-class product choice elsewhere, not hidden deployment trivia.
- `Harness` should add a stronger notion of **session durability tier**. Anthropic’s pinned background sessions and Managed Agents’ stateful session model both point toward runs that can be intentionally long-lived, not just temporarily active.
- `Harness` should model **tool-network topology** rather than only tool names. MCP tunnels, private servers, and self-hosted sandboxes imply that capability depends on reachability and trust zones, not merely which tools are installed.
- `Harness` should keep review and repair tightly linked. GitHub’s PR review to cloud-agent repair flow and CI-failure handoff are strong evidence that “validation produces the next agent task” is a better control pattern than forcing humans to manually restate defects.
- `Skills` should be treated as governed runtime modules. Between Claude Code’s richer plugin/session surfaces and GitHub’s open-source Eclipse view into skills and prompts, the field is making skill provenance and behavior less opaque.
- `erwin` knowledge work should keep distinguishing **app governance**, **platform governance**, and **local-runtime governance**. Anthropic now spans all three, but the changes are not interchangeable.

## Harness Evolution Proposals

### Title

Portable execution-environment contract

### Source

[Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) on **May 19, 2026** and [Claude Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview)

### Observed Pattern

Anthropic now supports both Anthropic-managed cloud containers and self-hosted sandboxes under the same higher-level Managed Agents contract, while still exposing sessions, events, tools, and MCP connectivity as stable objects.

### Why It Matters For Harness

Harness can coordinate work today, but it does not yet present execution placement as a durable, portable part of the run contract. The field is making environment choice operationally central.

### Recommended Disposition

`adopt`

### What We Would Change If Adopted

Define an execution-environment contract with stable fields for runtime owner, network zone, filesystem policy, credential mode, tool attachment rules, and persistence guarantees, so the same run semantics can span local, hosted, and self-hosted execution.

### Why We Should Not Cargo-Cult This

We should not import Anthropic’s exact product split or beta headers. The reusable idea is the separation between agent/session semantics and environment placement.

### Title

Pinned run classes and dormancy policy

### Source

[Claude Code changelog](https://code.claude.com/docs/en/changelog) on **May 21-22, 2026**

### Observed Pattern

Claude Code now distinguishes pinned background sessions from ordinary sessions and gives pinned sessions stronger survival behavior across idleness, updates, and memory pressure.

### Why It Matters For Harness

Not all long-running work deserves the same liveness guarantees. Without explicit run classes, background work degrades into accidental persistence or accidental eviction.

### Recommended Disposition

`prototype`

### What We Would Change If Adopted

Introduce run durability classes such as transient, background, and pinned, each with different wake, retention, compaction, and eviction behavior.

### Why We Should Not Cargo-Cult This

We do not need to mirror `claude agents` UX or keyboard shortcuts. The important part is making liveness guarantees explicit and machine-readable.

### Title

Review-to-repair handoff packets

### Source

[One-click fixes for failing Actions with Copilot cloud agent](https://github.blog/changelog/2026-05-18-one-click-fixes-for-failing-actions-with-copilot-cloud-agent/) on **May 18, 2026** and [Easily apply Copilot code review feedback with Copilot cloud agent](https://github.blog/changelog/2026-05-19-easily-apply-copilot-code-review-feedback-with-copilot-cloud-agent/) on **May 19, 2026**

### Observed Pattern

GitHub is turning validation failures and review comments directly into structured cloud-agent tasks, including batch application and model selection.

### Why It Matters For Harness

Harness already produces rich intermediate state, but it does not yet normalize “this validation result should become the next agent job” as a reusable control surface.

### Recommended Disposition

`adopt-directionally`

### What We Would Change If Adopted

Create structured repair packets for CI failures, review comments, and verification diffs so a downstream run can inherit evidence, target scope, preferred model/effort, and success criteria without prompt restatement.

### Why We Should Not Cargo-Cult This

GitHub’s flow is tightly tied to pull requests and Actions. Harness should extract the handoff pattern without assuming GitHub-native artifacts are always the control plane.

### Title

Task-aware model routing policy

### Source

[Auto model selection now routes based on your task in VS Code](https://github.blog/changelog/2026-05-20-auto-model-selection-now-routes-based-on-your-task-in-vs-code/) on **May 20, 2026** and the [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool)

### Observed Pattern

GitHub is routing by task shape and model health, while Anthropic keeps a separate executor-plus-advisor pattern for quality-sensitive planning on cheaper execution paths.

### Why It Matters For Harness

Harness will eventually need a principled way to choose between direct strong-model use, cheaper execution with escalation, and policy-driven routing based on task class.

### Recommended Disposition

`watch`

### What We Would Change If Adopted

Add a policy layer that can pick model family, effort, and escalation mode from task metadata, observed complexity, and operator constraints.

### Why We Should Not Cargo-Cult This

Neither GitHub’s auto router nor Anthropic’s advisor beta is mature enough to copy blindly. We should capture the decision framework, not a vendor-specific heuristic.

### Title

Edge-capable agent runtime adapters

### Source

[ADK for Kotlin and ADK for Android 0.1.0](https://developers.googleblog.com/adk-kotlin-android-building-ai-agents/) on **May 21, 2026**

### Observed Pattern

Google is moving agent runtimes toward Android and on-device execution while preserving tooling, plugins, state, and OpenTelemetry support.

### Why It Matters For Harness

If agent execution increasingly spans local desktop, hosted cloud, and mobile or edge contexts, Harness should not bake in server-only assumptions.

### Recommended Disposition

`watch`

### What We Would Change If Adopted

Define adapter boundaries so run state, tool contracts, and telemetry can survive different execution substrates, including constrained local environments.

### Why We Should Not Cargo-Cult This

Android support is not itself the goal for Harness. The durable lesson is substrate independence, not shipping a mobile agent product.

## Open Questions

- Anthropic now has both **self-hosted Managed Agent sandboxes** and increasingly durable **Claude Code** sessions. Do these converge into one portable session contract, or do local and hosted Anthropic runtimes keep distinct semantics long term?
- The required **advisor tool** surface was explicitly checked and remains active, but there was no new in-window advisor-specific release note. Will Anthropic keep advisor as a standalone planning primitive, or fold its function into broader managed orchestration?
- Claude’s **Compliance API integrations** arrived on **May 21, 2026**. Does Anthropic eventually expose a deeper admin/runtime bridge where compliance surfaces can reason about session classes, MCP endpoints, and skill provenance directly?
- GitHub is increasingly turning repo events into cloud-agent tasks. Does the durable pattern for `Harness` become “every validation artifact can emit a repair packet,” or is that too GitHub-shaped for a provider-neutral harness?
- Google’s **ADK for Android** raises a broader question: when agents become edge-capable, which parts of the control plane must stay cloud-resident, and which should travel with the runtime?
