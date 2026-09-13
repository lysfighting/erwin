# Daily Radar - 2026-09-13

Window: 2026-09-07 to 2026-09-13 (America/Los_Angeles). Checked 2026-09-13. One search pass per publisher family, followed by direct official pages. Compared with the latest persisted September 11 note, new intake is Claude Code 2.1.269–270, Kimi K2.8 Preview, and Kimi Work 3.2.8. Vendor claims are not local validation.

## New Signals

1. **Claude Code — plugin evaluation and more accurate execution state.** The [official changelog](https://code.claude.com/docs/en/changelog) dates **2.1.269 to 2026-09-11** and **2.1.270 to 2026-09-12**. The former adds `claude plugin eval` with scored JSON/HTML reports and optional repository attributes in OpenTelemetry. It repairs remote/headless sessions reporting input-wait while background agents run, stale git status after compaction, and missing path-denial records. The latter fixes a regression that made read-only git commands unexpectedly request permission in long sessions. Strong pressure for evaluating intermediate state and permission evidence alongside output quality.

2. **Moonshot — model identity changes behind a stable coding alias.** [Kimi Code release notes](https://www.kimi.com/code/docs/kimi-code/whats-new.html), **2026-09-11**, announce K2.8 Preview across Kimi Code with the unchanged `kimi-for-coding` ID, up to 1M context, and `low`/`high`/`max` thinking effort (default `max`). Disabling thinking routes both K3-family and K2.8 Preview requests to non-thinking K2.8 Preview. This is a documented routing contract; claimed coding gains remain unverified locally. Separately, [Kimi Work 3.2.8](https://www.kimi.com/en/help/kimi-work/release-notes), **2026-09-12**, fixes a login-refresh race that let stale requests clear login state or overwrite newer credentials, and clarifies anti-sleep modes.

3. **Claude Platform — Managed Agents policy evaluation remains the week's core API signal.** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), **2026-09-10**, add `auto` permission policies: server evaluation can run, deny, or pause agent/MCP tool calls for approval. Tool events expose evaluation evidence. The `ant beta:sessions connect` command supports live terminal observation, messages, and approval handling; `--web` opens a local viewer. Retained weekly signal, not a new September 13 release.

4. **DeepSeek — Flash migration does not imply Pro retirement.** The [official change log](https://api-docs.deepseek.com/updates/), **2026-09-10**, releases V4.1 Flash with native visual understanding under `deepseek-flash`. Retired V4 Flash and Vision Exp names temporarily route to it. The current notice explicitly continues V4 Pro API service after **2026-09-14**, with billing unchanged. Reverified today: the September 10 radar's earlier Pro cutover warning remains superseded. The notice's edit timestamp is unavailable.

5. **OpenAI — full-duplex voice can delegate backend work.** [GPT-Live-1 API announcement](https://openai.com/index/introducing-gpt-live-1-in-the-api/), **2026-09-10**, describes simultaneous listening/speaking and delegation of reasoning/tool work to a chosen text backend, including third-party models. Retained weekly signal: interruption handling and successful backend completion need separate evidence; vendor benchmarks do not establish Harness readiness.

## GitHub Trending Snapshot

[Default global Today board](https://github.com/trending), checked **2026-09-13**: [agent-skills](https://github.com/tech-leads-club/agent-skills), [OpenResearch](https://github.com/alphaXiv/OpenResearch), and [open-code-review](https://github.com/alibaba/open-code-review) surface reusable skills, parallel research, and hybrid review. Practical workflow packaging may attract attention; that explanation is inference. The [standalone snapshot](../trending/2026-09-13.md) preserves all 19 entries.

## Harness-Relevant Signals

- **State accuracy:** distinguish background execution, approval waiting, denial, and completion; compaction must not freeze an obsolete repository view.
- **Reproducibility:** a stable alias is insufficient to identify model behavior. Record effective routing and thinking settings in any future comparison.
- **Behavioral evaluations:** Google's [The Anatomy of Harness Engineering](https://developers.googleblog.com/the-anatomy-of-harness-engineering-how-to-evaluate-iterate-and-guard-ai-coding-agents/), **2026-09-09**, recommends checking discrete intermediate actions to explain regressions that aggregate end-to-end scores obscure. This complements the Code plugin-eval signal; no local eval was run.

## Candidate Proposals

### Proposal 1 — Execution-state evidence

- `source_type`: first-party Claude Code changelog and Claude Platform release notes
- `problem_pattern`: background execution, denied calls, and approval pauses are misclassified
- `why_it_matters_for_harness`: orchestration can report completion or request intervention at the wrong time
- `recommended_disposition`: `encode-as-eval`

### Proposal 2 — Effective model identity

- `source_type`: first-party Kimi Code routing documentation and DeepSeek lifecycle notes
- `problem_pattern`: stable aliases or thinking switches silently select different model behavior
- `why_it_matters_for_harness`: comparisons become misleading without effective routing provenance
- `recommended_disposition`: `watch`

### Proposal 3 — Reproducible plugin behavior checks

- `source_type`: first-party Claude Code release notes and Google engineering guidance
- `problem_pattern`: final scores hide tool-use and verification regressions
- `why_it_matters_for_harness`: a small reproducible behavior suite can expose the cause of drift
- `recommended_disposition`: `prototype`

Intake candidates only; no direct Harness implementation recommendation or runtime change.

## Watch Next

- **Anthropic required checks completed:** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and [Claude Code changelog](https://code.claude.com/docs/en/changelog), cross-checked with [official source](https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md). Advisor docs are capability context; no separately dated in-window advisor API launch established. Watch advisor-specific release evidence without conflating it with general model updates.
- Also checked [Anthropic Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), [Claude Blog](https://claude.com/blog/), and [legacy docs](https://docs.anthropic.com/en/docs/claude-code/overview), which redirect to Code. News lists a **2026-09-10** AI-misuse report; keep as security context, not an API release. No stronger engineering/blog release established in this bounded pass.
- **MiniMax:** [changelog](https://agent.minimax.io/docs/changelog), Desktop **3.0.71, 2026-09-11**, updates macOS signing for installation/launch compatibility; Windows remains 3.0.70. September 8's local/cloud session visibility and review scopes remain workflow context. No newer dated item established.
- **Qwen:** [official Model Studio lifecycle](https://help.aliyun.com/en/model-studio/newly-released-models) still lists the latest visible Qwen snapshot as **2026-09-02**, outside this window. No strong dated official technical item found in-window in the bounded pass; absence is not proven.
- Langfuse, Phoenix, Helicone, and OpenTelemetry secondary-source research was unnecessary today; Code's own telemetry change is first-party evidence.
- Repository freshness: `_project/CURRENT.md` remains **2026-04-19**. Historical daily-radar skill is absent from repository and searched installed skill paths. Followed current Radar docs and explicit output scope; stale handoff and unrelated files left untouched.
