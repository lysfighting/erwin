# AI Agent Weekly Synthesis — 2026-09-11

Date: 2026-09-11
Window: 2026-09-05 to 2026-09-11 (America/Los_Angeles)
Verified: 2026-09-11 local / 2026-09-12 UTC

## Field Summary

**The durable shift this week is toward testable execution contracts: who authorized a call, what actually ran, which events are missing, and what evidence supports completion.** Model capability still moved, but the most actionable Harness pressure is at these boundaries. This is synthesis from current first-party pages, not proof of local runtime adoption.

- **Anthropic Platform / Managed Agents — September 10:** the [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview) introduce per-call `auto` permission evaluation and terminal session attachment. This is a material Managed Agents change, not generic Anthropic news. Tool events expose permission outcomes and evaluation metadata; `ant beta:sessions connect` supports live interaction and approval handling.
- **Claude Code — September 11, version 2.1.269:** the [official changelog](https://code.claude.com/docs/en/changelog) adds plugin evals with JSON/HTML results and repository telemetry attributes, and fixes stale post-compaction Git status, misleading background-task status, missing path-denial records, and Bash `tee` write-boundary enforcement. The convergence is reproducibility plus truthful state reporting.
- **Google — September 9:** [The Anatomy of Harness Engineering](https://developers.googleblog.com/the-anatomy-of-harness-engineering-how-to-evaluate-iterate-and-guard-ai-coding-agents/) argues for observable behavioral checks alongside end-to-end evaluation. This supplies a method for turning the week's concrete failure boundaries into regression evidence.
- **OpenAI — September 5 and 10:** the [official changelog](https://learn.chatgpt.com/docs/changelog) records removal of `codex mcp-server` and the standalone binary on September 5. The suggested app-server integration remains experimental and unsupported for production workloads; external MCP connections remain supported. Python SDK 0.154.0 on September 10 introduces `ExternalMessage` with tool-level authority, explicitly without user authorization, and documents partial event histories for late subscribers.
- **DeepSeek — September 10:** [V4.1 Flash](https://api-docs.deepseek.com/updates/) brings native visual understanding and changes Flash alias routing. The currently published notice preserves V4 Pro service after September 14 with unchanged billing. Earlier retirement warnings must not be carried forward as current truth.

Compared with the [August 28 weekly baseline](2026-08-28.md), the new pressure is more specific than “runtime governance”: approval provenance, event completeness, execution-state accuracy, and behavioral verification now have concrete API or release evidence. The scheduler reports a September 4 local run, but no corresponding weekly artifact exists in this checkout; this note does not claim to reconstruct that missing run or silently expand its seven-day window.

## Theme Deep Dive

### Authorization is a typed boundary

The current [Managed Agents permission guide](https://platform.claude.com/docs/en/managed-agents/permission-policies) distinguishes execution, refusal, and a pause for approval. `auto` is opt-in and is not a mandatory human checkpoint. Custom tools remain the application's responsibility. Tool results and inter-thread messages do not express user intent; relaying untrusted text as `user.message` can incorrectly confer that status. Running sessions retain their initial tool configuration.

GitHub independently shipped [enterprise-managed agent permissions](https://github.blog/changelog/2026-09-09-enterprise-managed-permissions-for-github-copilot-agent-operations/) on **September 9**. Shell, file, and domain restrictions cannot be weakened by workspace settings, saved approvals, or auto-approval. This reinforces the need to preserve policy precedence across execution surfaces.

**Synthesis:** Harness should be able to explain the source of authority for an action and distinguish a recoverable refusal from mission termination. An external message may change what the agent knows without changing what it may do. This is a proposed invariant to review, not a claim that an existing local implementation is deficient.

### Monitorability depends on completeness, not just visibility

The [Codex SDK release](https://learn.chatgpt.com/docs/changelog) says late-attached turn handles only receive subsequent events; saved history must be read separately. Therefore a live stream ending is insufficient evidence of complete execution history. The [advisor tool documentation](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), checked this run, also distinguishes plaintext advice from encrypted result variants that the client cannot inspect. It remains beta; no independent in-window publication date was established for that page.

What stays visible: tool attempts, permission outcomes, externally reported lifecycle transitions, and delivered artifacts. What may remain hidden: earlier events not replayed to a subscriber and redacted advisor content. **Local review principle:** preserve evidence completeness and unknown states; do not infer internal reasoning quality from an opaque payload or declare success from a partial transcript.

### Evaluation should constrain both behavior and cost

Google's September 9 guidance recommends starting with a real failure, allowing multiple valid solution paths, and using aggregate evidence for noisy agent behavior. Its intermediate checks complement final outcomes. Anthropic's **September 8** [cost and performance guide](https://claude.com/blog/reducing-cost-and-improving-performance-with-claude-platform) connects caching, prompt changes, and effort selection to evaluations; its hillclimb workflow uses a held-out test split. Vendor-reported savings are not portable local results.

A separate scale signal is OpenAI's [Navier–Stokes report](https://openai.com/index/navier-stokes-solution/), published **September 8**: it reports a September 5 result followed by Lean formalization and verification, with approximately 130 billion output tokens for that problem. Treat this as a publisher-reported research result, not independently verified mathematical acceptance or a productivity benchmark. The reusable lesson is complete budget and verification accounting. No dedicated paper-watch note for this report was created in this run; research-method conclusions remain provisional.

## Tooling Watch

| Surface | Dated signal | Interpretation |
|---|---|---|
| [Google Cloud Developer Plugin](https://cloud.google.com/blog/topics/developers-practitioners/introducing-the-google-cloud-developer-plugin-for-ai-coding-agents) | 2026-09-10 | Packages related skills, documentation grounding, and MCP configuration. Packaging compatibility does not establish equivalent permission semantics across hosts. |
| [Kimi Work](https://www.kimi.com/en/help/kimi-work/release-notes) | 3.2.6, 2026-09-07; 3.2.7, 2026-09-11 | MCP connection visibility, inherited conversation settings, and recovery of the scheduled-task page when credits run out make operational readiness more explicit. |
| [MiniMax Agent](https://agent.minimax.io/docs/changelog) | 0.3.11, 2026-09-09 | Repairs long tasks failing after repeated credential renewals. Authentication lifetime belongs in continuity testing. |
| [GPT-Live-1 API](https://openai.com/index/introducing-gpt-live-1-in-the-api/) | 2026-09-10 | Full-duplex voice is another surface where conversational responsiveness should be assessed separately from successful action completion. |

No Langfuse, Phoenix, Helicone, or OpenTelemetry release was needed as a substitute for first-party evidence. Claude Code's telemetry addition is relevant context, not a reason to adopt another observability stack.

## Implications For erwin / Harness / Skills

1. **erwin:** retain source dates and correction history. The [September 11 daily note](../daily/2026-09-11.md) already corrected DeepSeek's deadline; this weekly note carries that correction into durable interpretation.
2. **Harness:** prioritize review of permission and completion evidence over new orchestration machinery. The existing [public contract](../../harness-public-contract.md) already defines completion through state transition plus verification; external releases strengthen the need to test that invariant.
3. **Skills:** require a demonstrated callable capability and a bounded success check before describing an installed plugin as ready. Keep vendor tuning helpers at evaluation/prototype status until local tasks support their tradeoffs.
4. **Compatibility:** inventory launch-command dependencies before a Codex upgrade. This run did not perform a complete runtime dependency audit, so the MCP-server removal is an exposure question, not a confirmed Harness outage.

## Knowledge Threads To Update

These are proposed follow-through targets in the [active thread registry](../../knowledge/active-threads.md); the registry itself is outside this run's write scope.

| Existing thread | Why this is the durable target | Evidence and posture |
|---|---|---|
| `harness engineering as legibility + eval pressure` | Permission and history boundaries determine whether runtime claims can be trusted. | Current official permission/SDK docs plus the local public contract; bounded eval candidate review. |
| `benchmark integrity and eval awareness` | Behavioral checks and cost tuning need independent outcome evidence. | Google guidance and Anthropic cost guide; retain held-out outcomes and variability, no benchmark-ranking conclusion. |
| `knowledge-system operating discipline` | Corrections and incomplete weekly coverage must remain explicit. | Current DeepSeek notice and the dated local daily correction; preserve provenance rather than duplicating proposals. |

## Harness Evolution Proposals

Only one candidate meets the bar for a bounded review. It extends an existing thread rather than creating a new runtime feature.

- **Title:** Authority and completion evidence across interrupted execution.
- **Source:** September 10 Managed Agents policies and Codex SDK notes; September 9 Google behavioral-eval guidance; [Harness public contract](../../harness-public-contract.md); existing [eval candidate ledger](../../proposals/harness-eval-candidate-ledger.md).
- **Observed Pattern:** Knowledge arrival, permission decisions, tool execution, and transcript completeness have separate lifetimes. Combining them into a single success flag can hide unresolved work.
- **Why It Matters For Harness:** Its public completion contract depends on verified state, while its orchestration role crosses message and execution boundaries. Evidence classes are external plus an existing internal contract; no new local failure was reproduced.
- **Recommended Disposition:** Candidate for `encode-as-eval` review; not adopted. First compare existing fixtures and link a concrete local gap before adding tests.
- **What We Would Change If Adopted:** Add the smallest missing fixtures for external content conferring no authority, approval pauses surviving resume, and incomplete event history never proving completion. Require artifact/state verification for terminal success. Keep any stochastic agent checks separate from deterministic state-machine assertions.
- **Why We Should Not Cargo-Cult This:** Do not copy vendor classifiers, default permissions, SDK event names, or an exact tool sequence. Preserve the local contract and accept equivalent valid execution paths.

## Open Questions

- Does any active Harness integration still launch the removed Codex MCP server, and what supported adapter boundary would replace it?
- Which existing fixtures already cover permission precedence and partial-history completion? This week's sources justify review, not duplicate tests.
- Can advisor experiments measure outcomes while explicitly recording unavailable advice content? The API page is undated; no new advisor API release is asserted for this week.
- Which local workloads would support held-out cost/quality comparisons before trying automated prompt optimization?

### Coverage and knowledge gaps

All required Anthropic surfaces were explicitly opened on this run: [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview) and [quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Claude Code changelog](https://code.claude.com/docs/en/changelog), and the [legacy docs overview](https://docs.anthropic.com/en/docs/claude-code/overview), which redirects to code.claude.com. [Anthropic Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), and [Claude Blog](https://claude.com/blog) were also checked. The engineering/news indexes returned limited extractable content; this is not a claim of exhaustive coverage. Managed Agents and Claude Code are positively confirmed material changes; advisor API dating remains unknown.

One search pass per publisher family was used, followed by direct official pages, including pages identified in this week's local daily notes. Qwen yielded no verified in-window primary signal in that bounded pass; this does not establish no release. Current first-party verification supersedes earlier summaries. `_project/CURRENT.md` still carries an April 19 update, and the latest prior weekly file is August 28; neither is treated as a fresh September runtime audit. Only the permitted weekly and generated knowledge outputs are changed.
