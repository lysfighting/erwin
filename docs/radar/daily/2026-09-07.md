# Daily Radar - 2026-09-07

Window: 2026-09-01 to 2026-09-07 (America/Los_Angeles). Checked 2026-09-07. One search pass per publisher family, followed by direct official pages. Compared with the September 6 note, newly captured evidence includes Google's individually dated launches and MiniMax 0.3.3–0.3.4; other entries are explicitly retained weekly signals, not September 7 launches. Vendor claims are not local Harness validation.

## New Signals

1. **Google — stronger coding agents need explicit effort budgets; video agents select evidence dynamically.** [Gemini 3.8 Flash and Flash Cyber](https://blog.google/innovation-and-ai/models-and-research/gemini-models/3-8-flash-and-3-8-flash-cyber/), **2026-09-02**, targets long-horizon coding; Google notes that higher effort can increase token use. Cyber availability is restricted to trusted defenders through Fairwind. [Agentic video understanding](https://blog.google/innovation-and-ai/models-and-research/gemini-models/introducing-agentic-video-in-gemini/), **2026-09-01**, lets supported Gemini models search and revisit relevant video segments instead of relying only on fixed-rate ingestion. Reported quality and efficiency gains are vendor benchmarks. These dated sources improve on yesterday's ambiguous rolling-roundup coverage.

2. **MiniMax — distinguish authentication failure from service failure.** [Agent changelog](https://agent.minimax.io/docs/changelog), **0.3.4, 2026-09-06**, preserves valid credentials during transient errors, falls back to device authorization after permanent rejection, and reports resource/service failures without misleading login prompts. **0.3.3, 2026-09-05**, synchronizes sign-in state immediately and permits one content-review credential recovery attempt while keeping failed requests blocked. This extends the **0.3.2, 2026-09-04** side-conversation boundary: modifying tools need explicit confirmation in every permission mode. Newly captured since the prior note; search snippets were stale, so the direct changelog takes precedence.

3. **Anthropic — resource identity and session recovery remain this week's control-plane signal.** [Claude Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), **2026-09-03**, adds `ant apply` in CLI 1.30.0: repository resource definitions, an approval plan, and `claude-lock.json` retain resource identity across applications. [Claude Code changelog](https://code.claude.com/docs/en/changelog), **2.1.261, 2026-09-04**, adds unused-skill/context diagnostics and fixes resumed parallel-tool context, stale remote permission display, and early interrupts. **2.1.263, 2026-09-06**, contains only generic reliability notes. Retained weekly signal; no newer substantive Code entry was visible.

4. **OpenAI — measure supervised success separately from autonomy.** [Research acceleration: The view inside OpenAI](https://openai.com/index/research-acceleration-view-inside-openai/), **2026-09-06**, reports rising concurrent coding-agent usage and experiment activity, while cautioning that activity does not directly measure research progress. Its claim of research-intern capability explicitly concerns well-defined work under human direction. Retained from yesterday; observational internal evidence is insufficient to infer causal productivity or unsupervised research capability.

5. **Moonshot — collaboration and permission semantics are evolving together.** [Kimi Code v0.41.0](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html), **2026-09-04**, introduces experimental Tower collaboration in Web and direct delivery of background-question answers. Its Auto-mode dangerous-command behavior differs from the guard described in **v0.40.0, 2026-09-02**. Retained weekly signal: version and actual mode behavior matter more than familiar permission labels.

## GitHub Trending Snapshot

Secondary signal from the [default global Today board](https://github.com/trending), captured **2026-09-07**: [context-mode](https://github.com/mksglu/context-mode) exposes context-management demand; [ECC](https://github.com/affaan-m/ECC) and [deer-flow](https://github.com/bytedance/deer-flow) package agent workflows; [markitdown](https://github.com/microsoft/markitdown) supports document intake. These are inferred reasons for relevance, not established causes of popularity. See the [full standalone snapshot](../trending/2026-09-07.md).

## Harness-Relevant Signals

- **Failure classification:** MiniMax makes expired identity, transient infrastructure errors, and policy blocking distinct states. A recovery evaluation should preserve those differences.
- **Budgeted evidence gathering:** Google's selective video inspection and effort tradeoffs suggest evaluating task quality together with tokens, tool calls, and missed evidence.
- **Reproducibility:** Anthropic's resource lockfile and Code recovery fixes reinforce that restored identity and context must be verified, not inferred from a resumed UI.
- **Supervision:** OpenAI's observations and Kimi's collaboration changes support recording human steering and permission-mode versions alongside successful outcomes.
- **Qwen coverage:** [Model Studio model updates](https://help.aliyun.com/en/model-studio/newly-released-models) confirm **qwen3.8-max-0902, 2026-09-02**, with vendor claims of stronger long-horizon coding and agent collaboration; **2026-09-01** retrieval releases include code/memory retrieval. Retained watch context; no additional strong dated item was established in this bounded pass.

## Candidate Proposals

### Proposal 1 — Recovery failure taxonomy

- `source_type`: first-party MiniMax changelog (New Signal 2)
- `problem_pattern`: infrastructure or policy failures become login loops, or temporary errors discard valid identity
- `why_it_matters_for_harness`: recovery must preserve credentials and the actual blocker while keeping denied work blocked
- `recommended_disposition`: `encode-as-eval`

### Proposal 2 — Cost and evidence coverage

- `source_type`: first-party Google model and video announcements (New Signal 1)
- `problem_pattern`: more reasoning improves quality but can hide extra cost; selective evidence ingestion can omit decisive moments
- `why_it_matters_for_harness`: compare completion quality, resource use, and evidence recall under fixed budgets
- `recommended_disposition`: `prototype`

### Proposal 3 — Versioned collaboration authority

- `source_type`: first-party Anthropic, Moonshot, and MiniMax release notes
- `problem_pattern`: resumed tasks, side questions, and background agents can expose stale authority or ambiguous delivery state
- `why_it_matters_for_harness`: coordination quality depends on accurate ownership, permissions, and recovery semantics
- `recommended_disposition`: `watch`

These are intake candidates only, not reviewed implementation decisions. No Harness code changes are proposed directly.

## Watch Next

- **Required Anthropic checks completed:** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), and [Claude Code changelog/docs](https://code.claude.com/docs/en/changelog). Managed Agents retains the dated `ant apply` signal. Advisor documentation describes executor/advisor consultation, but no separately dated in-window update was established; do not label undated docs new.
- Also checked [Anthropic News](https://www.anthropic.com/news), [Engineering](https://www.anthropic.com/engineering), [Claude agents blog surface](https://claude.com/blog-category/agents), and [legacy Anthropic docs](https://docs.anthropic.com/en/docs/intro), which redirects to Platform. The Claude blog lists commerce-agent articles dated **2026-09-02**; those are lower-priority leads than the control-plane changes above.
- **DeepSeek:** the official [change log](https://api-docs.deepseek.com/updates/) still has **2026-08-21** as its newest visible technical entry, outside the window. No strong dated in-window technical release found in the single search pass; do not equate that with proof of no change.
- Check whether Google's effort/evidence claims reproduce under bounded workloads and whether MiniMax recovery semantics remain stable. Keep Qwen and advisor developments on watch. No secondary observability tooling release was needed to clarify today's first-party signals.
- Repository freshness: `_project/CURRENT.md` remains dated **2026-04-19**. This run follows the explicit automation scope and current radar files; it does not treat that handoff as fresh runtime evidence or edit it.
