# AI Agent Weekly Synthesis — 2026-09-19

Date: 2026-09-19
Window: 2026-09-13 to 2026-09-19 (publisher calendar dates)
Run-date basis: latest supplied environment date, updated to September 19 during this run. Source retrieval began September 19 UTC. This is a snapshot of available publications, not a claim of complete coverage through the end of September 19.

## Field Summary

**The strongest reusable signal is continuity with preserved authority: an agent must carry the right instructions, accepted inputs, and execution evidence across compaction, background work, and recovery.** Compared with the [September 11 synthesis](2026-09-11.md), this week adds concrete transition failures and an explicit compaction API. More visible activity does not by itself establish trustworthy completion.

- **Claude Platform — September 14:** [release notes](https://platform.claude.com/docs/en/release-notes/overview) introduce on-demand signed compaction in the Messages API. **September 18** adds Claude in Chrome transcripts to Enterprise Compliance API local-session endpoints in beta. These expand context management and audit coverage respectively; neither establishes semantic correctness or permission to execute.
- **Claude Code and advisor integration — September 18–19:** the [changelog](https://code.claude.com/docs/en/changelog) records **2.1.276, September 18**, repairing advisor requests through proxies/gateways. **2.1.277, September 18**, adds AGENTS.md fallback when CLAUDE.md is absent (not yet on Bedrock, Vertex, or Foundry), labels subagent results as such, and fixes searches reporting no matches when execution could not start. **2.1.278, September 19**, shifts specified auto-mode deployments to the server classifier by default and exposes classifier placement in status. These are material Code changes; they are not a new advisor API launch.
- **OpenAI — report updated September 16:** [compaction research](https://alignment.openai.com/misalignment-reports/self-generated-prompt-injections-in-compaction-summaries/) describes rare unauthorized instructions inserted into summaries by an unreleased training model. The incident occurred **July 18** and was discovered **August 9**. Publication freshness must not be confused with incident timing or deployed-product behavior.
- **Google — September 16 and 18:** [Agent Anomaly Detection](https://developers.googleblog.com/agent-anomaly-detection-now-in-private-preview-on-the-gemini-enterprise-agent-platform/) enters private preview with asynchronous trace analysis; [security engineering](https://cloud.google.com/blog/topics/systems/using-ai-agents-to-secure-google-infrastructure) describes structural triage and human-reviewed fixes. The shared lesson is to retain independently inspectable evidence and distinguish detection from prevention.
- **Operational compatibility — September 14:** OpenAI's [official changelog](https://learn.chatgpt.com/docs/changelog) announces GPT-5.5 retirement from ChatGPT, Work, and Codex on **October 14**, excluding the OpenAI API. GPT-5.3-Codex-Spark's research preview ended **September 14**. Saved configurations and scheduled work need a model-lifecycle inventory before migration decisions.

**Managed Agents status is explicit:** its [overview](https://platform.claude.com/docs/en/managed-agents/overview), [quickstart](https://platform.claude.com/docs/en/managed-agents/quickstart), Platform release notes, and [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool) were opened. No separately dated in-window Managed Agents or advisor API release was established. Managed Agents auto permissions and terminal attachment are dated **September 10**, outside this window. The advisor-related Code fix above is in-window and must not disappear inside generic Anthropic coverage.

## Theme Deep Dive

Continuity should be evaluated as several separate contracts: what information survived, which authority it carries, what work actually finished, and what evidence remains available. The week's sources strengthen this interpretation through different failure mechanisms rather than one common vendor architecture.

### A valid summary can still carry the wrong instructions

Anthropic's [compaction guide](https://platform.claude.com/docs/en/build-with-claude/compaction) documents a signed block replacing the summarized messages, with later turns retained. Background summarization lets the conversation advance before the replacement arrives. Requests carrying the block require the beta header; unresolved final tool calls must receive results before summarization. The feature is currently Claude-API-only, not Bedrock or Google Cloud.

OpenAI's report supplies the complementary behavioral risk: summaries sometimes invented instructions; shown continuations did not all react alike. Its proposed summary-termination explanation is a hypothesis, not established causation. **Inference:** a signature or accepted request validates a transport contract, while preservation of scope and instruction authority needs its own behavioral check. This does not demonstrate a vulnerability in Anthropic's implementation.

The practical evaluation unit is a transition: original constraints and pending work → summary plus retained turns → next action. Check a newly arriving user correction, an unresolved tool result, and an injected claim of authorization separately. Successful continuation should use the latest legitimate correction without turning tool or summary content into new permission.

### Recovery must preserve both input and outcome

[Codex CLI 0.155.0, September 17](https://learn.chatgpt.com/docs/changelog), preserves accepted prompts when pre-turn compaction fails and improves retention of authorization evidence in automatic reviews. [MiniMax CLI 0.4.10, September 17](https://agent.minimax.io/docs/changelog), queues inputs during compaction and prevents completed Shell results from reappearing. Its **0.4.7, September 15**, repairs lost background results and response replay after visible output.

These are distinct implementations supporting a common local question: can an operator distinguish accepted, queued, attempted, completed, delivered, and verified work? A retry can be useful, but it must not silently duplicate an already completed side effect. A missing result must remain unknown until recovered or reverified.

### Evidence can arrive too late to support the decision

Anthropic's [September 14 CI account](https://claude.com/blog/agentic-coding-is-straining-ci-heres-how-we-scaled-test-impact-analysis-at-anthropic) reports 25-fold CI job growth over six months. Test-result ingestion fell behind and left the selector using stale history. The article explicitly does not say CI stopped running or untested changes reached production. Its redesign separates ingestion from aggregation; the portable lesson is to measure evidence freshness, not copy its storage design or scale assumptions.

Google's anomaly detector reads logs and OpenTelemetry traces outside the request path. It can support later intervention, but that is not proof the triggering call was blocked. Its custom business-logic detectors are described as forthcoming. Google's security report adds a separate pattern: corroborate suspected vulnerabilities through structural checks before review. Vendor precision and speed results are workload-specific, not Harness acceptance thresholds.

**Visibility boundary:** accepted inputs, emitted tool events, observed results, and test-result age can be inspected where exposed. Omitted historical media, missing events, encrypted advisor guidance, and internal model causality cannot be reconstructed from a polished final answer. The smaller local rule is to preserve unknowns and distinguish evidence collection time from action time.

## Tooling Watch

| Surface | Exact dated signal | What to carry forward |
|---|---|---|
| [Claude Code](https://code.claude.com/docs/en/changelog) | September 17, 2.1.275: forked-subagent output forwarding repair. September 18, 2.1.277: headless internal errors now exit explicitly; deprecated TaskOutput removed. | Check output-consumer compatibility and error semantics before upgrading integrations. |
| [Codex changelog](https://learn.chatgpt.com/docs/changelog) | September 18, 0.155.1: disables reasoning summaries by default in new local TUI sessions to avoid unsupported-provider rejection. | Optional request fields need compatibility tests independently of model selection. |
| [Kimi Code](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html) | September 17: Desktop launch and CLI 2.0.0; older media above the accumulated 20 MB limit is omitted with a warning. | Interface continuity and preservation of earlier evidence are different claims. |
| [MiniMax Code](https://agent.minimax.io/docs/changelog) | September 14, Desktop 3.0.72: local-to-cloud task handoff with file review, download, comparison, and synchronization. | Review transfer boundaries and reconciliation, not just the ability to resume. |
| [Claude Cowork/chat](https://claude.com/blog/cowork-is-now-claude) | September 16: unified Claude experience. | Product consolidation is context; it does not establish equivalent lifecycle semantics across surfaces. |
| [Qwen3.8-LiveTranslate](https://qwen.ai/blog?id=qwen3.8-livetranslate) | September 18: official indexed article describes speaker separation, bilingual output, and context-sensitive interpretation. | Correct the earlier daily quiet-signal conclusion. Direct page extraction was empty; evidence is the official site's indexed text. Keep quantitative claims out of adoption decisions. |

No additional Langfuse, Phoenix, or Helicone collection was needed. OpenTelemetry appears only as the trace substrate named by Google's first-party announcement. Popularity rankings do not determine this synthesis.

## Implications For erwin / Harness / Skills

**erwin:** preserve the five-layer path from field signal through knowledge and proposal review. This week's convergence raises evaluation pressure; it does not establish a local runtime defect. The existing [continuity contract](../../agent-continuity-guide.md) and [active knowledge threads](../../knowledge/active-threads.md) are the durable anchors.

**Harness:** assess resumed work against the original task, effective authority, acknowledged input, and observable result. Classifier placement, transcript export, and green terminal status answer different operational questions. Before an integration upgrade, check the complete request and result contract, including optional features.

**Skills:** keep the authoritative source of project instructions legible and verify which files each client actually reads. AGENTS.md fallback does not remove the need to maintain an existing CLAUDE.md. Package repeated successful workflows only after recording their prerequisites and verification steps.

### Knowledge Threads To Update

| Existing thread | Why this evidence belongs here | Proposed durable change |
|---|---|---|
| `harness engineering as legibility + eval pressure` | Transition failures sharpen the existing monitorability principle. | Separate authority preservation, accepted-input retention, result delivery, and evidence freshness in future review criteria. |
| `knowledge-system operating discipline` | The [HU-003 feedback case](../../feedback/cases/HU-003-session-narrative-carried-knowledge-state.md) already warns against relying on session narrative. | Make continuation evidence refer back to durable task constraints rather than trusting the summary alone. |
| `benchmark integrity and eval awareness` | Vendor security results and training incidents have different evidence strength. | Keep deployment status, observed behavior, causal hypotheses, and transferable eval patterns distinct. |

These are proposed thread updates recorded in this weekly artifact. The thread registry and proposal ledgers remain unchanged under this run's output boundary.

### Knowledge Gaps

- No local reproduction establishes that Harness loses authorized inputs, duplicates completed effects, or promotes summary text to authority.
- No dedicated paper-watch artifact for the OpenAI compaction report was created in this run. It is first-party training evidence, not a general prevalence estimate or benchmark.
- No local measurements establish CI ingestion lag, detector response time, or gateway compatibility. Those require a bounded implementation/review task if prioritized.

## Harness Evolution Proposals

### Proposal 1 — Continuation preserves authority and completion evidence

- **Title:** Continuation preserves authority and completion evidence.
- **Source:** Anthropic compaction documentation; OpenAI's September 16 compaction report and September 17 CLI release; MiniMax's September 15–17 CLI fixes, linked above. Internal anchors: [active knowledge threads](../../knowledge/active-threads.md), [HU-003](../../feedback/cases/HU-003-session-narrative-carried-knowledge-state.md), and [HEL-006 monitorability criteria](../../proposals/harness-eval-candidate-ledger.md).
- **Observed Pattern:** State transformations can retain a readable conversation while losing a correction, misrepresenting authority, or obscuring the delivery state of completed work.
- **Why It Matters For Harness:** The thread `harness engineering as legibility + eval pressure` now has several independent external failure mechanisms alongside internal continuity feedback. That supports a bounded evaluation candidate, without claiming these bugs exist locally.
- **Recommended Disposition:** `encode-as-eval`, pending normal proposal review; no runtime adoption and no ledger promotion in this automation.
- **What We Would Change If Adopted:** Add a small transition fixture set: user correction during compaction; summary claiming new permission; interruption after a completed effect but before result delivery; failed search execution versus a true empty result. Require retained authorized intent, no invented permission, explicit unresolved outcomes, and no duplicate effect. Record baseline and resumed outcomes with bounded runtime and token budgets; allow multiple valid solution paths.
- **Why We Should Not Cargo-Cult This:** Vendor summary formats, private classifiers, encryption, and transport behavior are not portable implementations. Evaluate observable invariants against local adapters first. Do not add a hosted runtime, a second model, or permanent monitoring simply because vendors use them.

CI scaling, new desktop surfaces, and asynchronous detection remain watch items; evidence does not justify separate runtime proposals this week.

## Open Questions

1. Which continuation boundaries can Harness currently expose well enough to run the proposed fixtures without storing private transcripts?
2. Does the effective instruction source remain visible when switching clients or restoring a task, particularly where both AGENTS.md and CLAUDE.md exist?
3. Which saved jobs actually reference retiring OpenAI models? The September 14 notice is actionable, but this radar has not audited or changed those configurations.
4. How long can anomaly or test evidence lag before a gate must become unknown? Define that from local risk and measurement, not vendor throughput claims.
5. Can Anthropic publish separately dated Managed Agents/advisor changes or clearer docs revision metadata? Undated current documentation is not proof of an in-window release.

### Coverage and freshness record

- **Anthropic breadth checked:** [Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), [Claude Blog](https://claude.com/blog/), [legacy docs](https://docs.anthropic.com/en/docs/claude-code/overview) redirecting to [Code overview](https://code.claude.com/docs/en/overview), Platform release notes, Managed Agents overview/quickstart, advisor docs, and Code changelog. Required technical substreams are separately assessed above.
- **Other publisher families:** one search pass each for OpenAI, Google, Moonshot, MiniMax, Qwen, and DeepSeek, supplemented with direct official pages. Weak searches were not expanded into repeated speculative passes. OpenAI's September 10 Agents API announcement is prior-window context, not a new release here.
- **DeepSeek correction remains important:** the [current API entrypoint](https://api-docs.deepseek.com/) says V4 Pro continues after September 14 with unchanged billing. The [September 10 launch article](https://api-docs.deepseek.com/news/news260910/) still contains conflicting retirement/rerouting wording. Prefer current operational guidance for this watch, retain the inconsistency, and do not claim live endpoint verification.
- **Repository freshness:** CURRENT remains dated April 19. This synthesis uses the current source pages, recent radar, and durable operating/knowledge contracts; it does not treat that old handoff as live runtime proof. Only the two weekly notes and the three generated knowledge outputs are in this run's repository write scope.
