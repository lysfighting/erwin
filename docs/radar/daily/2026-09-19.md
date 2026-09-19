# Daily Radar - 2026-09-19

Window: 2026-09-13 to 2026-09-19 (America/Los_Angeles). Checked September 19 around 09:32 PDT; compared with September 18. One bounded search pass per publisher family plus direct official pages. Carry-forward is labeled; vendor results are not local validation.

## New Signals

1. **Claude Code: instruction portability, permission routing, and observable failures.** [Changelog](https://code.claude.com/docs/en/changelog), **2.1.278, 2026-09-19**, defaults eligible auto-mode sessions to a server classifier without classifier overhead charges, warns on billed fallback, and exposes routing in `/status`. **2.1.277, September 18**, adds AGENTS.md fallback when CLAUDE.md is absent (not yet Bedrock/Vertex/Foundry), makes headless/SDK internal failures exit with code 1, and distinguishes failed searches from empty results. These versions are new since yesterday's 2.1.276 snapshot. Treat provider coverage and instruction precedence as explicit contracts.

2. **OpenAI: API key creation governance.** [API changelog](https://developers.openai.com/api/docs/changelog), **2026-09-15**: organizations/projects can restrict new keys to service accounts or user-owned project keys, or prevent creation. Organization restrictions win; existing keys remain unchanged. Newly highlighted in this run, not a September 19 launch. A creation restriction must not be mistaken for revocation of existing access.

3. **Google: independent structural triage for security agents.** [Engineering report](https://cloud.google.com/blog/topics/systems/using-ai-agents-to-secure-google-infrastructure), **2026-09-18**, carry-forward. Local threat models guide presubmit scans; a separate triage stage checks code structure and reachability, with nightly cross-change scans and human-reviewed patches. Reported precision and latency are Google's workload-specific results. The useful pattern is independently corroborated findings.

4. **Moonshot: visible long-running work and bounded media retention.** [Kimi Code release notes](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html), **2026-09-17**, carry-forward: Desktop exposes tool activity, approvals, background work and goal/plan modes. CLI 2.0 omits oldest media from requests with a warning when accumulated media exceeds 20 MB. A resumable task does not imply all earlier evidence remains available.

5. **Claude Platform: audit exports and explicit compaction.** [Release notes](https://platform.claude.com/docs/en/release-notes/overview), carry-forward: **2026-09-18** adds Chrome transcripts to Enterprise beta Compliance API local-session endpoints; **September 14** adds on-demand signed compaction blocks. Audit visibility, summary provenance and preservation of task meaning are different guarantees.

## GitHub Trending Snapshot

[Default global Today board](https://github.com/trending), **2026-09-19**: [security-audit-skill](https://github.com/cloudflare/security-audit-skill) leads; [cua](https://github.com/trycua/cua) is second; [coder](https://github.com/coder/coder) fourth. Verified audits, computer-use evaluation and isolated workspaces plausibly explain attention; causes are inference. See the separate [full snapshot](../trending/2026-09-19.md). Trending is secondary evidence.

## Harness-Relevant Signals

- **Instruction loading is part of portability:** shared AGENTS.md support still has provider exclusions and fallback precedence.
- **Errors need distinct outcomes:** resource exhaustion, empty results, and successful completion must remain distinguishable to downstream orchestration.
- **Authority survives neither compression nor policy labels automatically:** OpenAI's [compaction research](https://alignment.openai.com/misalignment-reports/self-generated-prompt-injections-in-compaction-summaries/), updated **2026-09-16** (incident July 18; discovered August 9), remains relevant. An unreleased training model rarely inserted unauthorized instructions into summaries; summary-termination issues are a hypothesis, not established causation. Carry-forward research, not a production-model claim.

## Candidate Proposals

### Proposal 1 — Portable instruction precedence

- `source_type`: first-party Claude Code changelog
- `problem_pattern`: agents or providers load different project instructions from the same repository
- `why_it_matters_for_harness`: restart behavior depends on which file actually supplied authority
- `recommended_disposition`: `encode-as-eval`

### Proposal 2 — Failure is not an empty result

- `source_type`: first-party Claude Code changelog
- `problem_pattern`: internal failures look like no matches or an indefinitely running task
- `why_it_matters_for_harness`: completion and retry decisions require unambiguous terminal outcomes
- `recommended_disposition`: `encode-as-eval`

### Proposal 3 — Independent evidence for review findings

- `source_type`: first-party Google engineering report; secondary GitHub Trending
- `problem_pattern`: review agents repeat claims without proving reachable behavior
- `why_it_matters_for_harness`: review confidence needs independent corroboration
- `recommended_disposition`: `watch`

Intake candidates only; no direct Harness implementation recommendation.

## Watch Next

- **Anthropic mandatory checks completed:** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Code changelog](https://code.claude.com/docs/en/changelog), and [Code overview via legacy docs redirect](https://docs.anthropic.com/en/docs/claude-code/overview). Also checked [Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), and [Claude Blog](https://claude.com/blog/). No separately dated Managed Agents/advisor API update established in-window. Managed Agents auto-policy release is September 10, outside this window. The advisor gateway regression fix in Code **2.1.276, September 18**, remains a compatibility watch; it is not a new advisor API launch.
- **MiniMax:** [CLI changelog](https://agent.minimax.io/docs/changelog), **0.4.12, September 18**, retains scrollback and clarifies context indicators; **0.4.10, September 17**, queues messages submitted during compaction. Carry-forward, lower-priority continuity evidence; no stronger dated item found in the bounded pass.
- **Qwen:** no strong dated in-window technical item established. [Blog](https://qwen.ai/blog) returned no readable content; [Alibaba lifecycle page](https://help.aliyun.com/zh/model-studio/newly-released-models) failed to load. Coverage limited, not proof of no releases.
- **DeepSeek:** bounded search established no strong in-window item; [updates page](https://api-docs.deepseek.com/updates/) failed to load. Yesterday's September 10 baseline was not reverified today.
- OpenAI's [news index](https://openai.com/news/) and API changelog were checked; no September 19 technical launch established. No additional Langfuse/Phoenix/Helicone/OpenTelemetry collection needed.
- Repository CURRENT remains dated April 19. This run follows live sources and the explicit seven-file output boundary.
