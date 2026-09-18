# Daily Radar - 2026-09-18

Window: 2026-09-12 to 2026-09-18 (America/Los_Angeles). Checked September 18 around 09:50 PDT; compared with September 17. One search pass per publisher family plus direct official pages. Vendor claims are not local verification; carry-forward items are labeled.

## New Signals

1. **Claude Code: advisor gateway regression fixed.** [Changelog](https://code.claude.com/docs/en/changelog), **2.1.276, 2026-09-18**, fixes requests failing on `advisor_20260301` when `ANTHROPIC_BASE_URL` targets a proxy/gateway, a 2.1.275 regression. **2.1.275, September 17**, also repairs forked-subagent output forwarding and credential exposure in marketplace messages, warns on telemetry-header helper failure, and adds account skill/plugin synchronization. New since yesterday's 2.1.274 baseline. Compatibility and configuration provenance deserve explicit eval coverage.

2. **Google: security agents paired with structural verification.** [Engineering report](https://cloud.google.com/blog/topics/systems/using-ai-agents-to-secure-google-infrastructure), **2026-09-18**, describes presubmit scans using localized threat models, structural reachability checks by a triage agent, nightly cross-change scans, and fixes submitted for human review. Google reports over 92% triage precision and sub-minute operation; these are workload-specific vendor results. Keep development, scanning, and triage contexts independent when evaluating the approach.

3. **Moonshot: Kimi Code Desktop and CLI 2.0.0.** [Official What's New](https://www.kimi.com/code/docs/en/kimi-code/whats-new.html), **2026-09-17**: Desktop exposes tool activity, approvals, goal/plan modes, background tasks, and an integrated browser. CLI 2.0 adds a desktop entrypoint and terminal Mermaid rendering; accumulated session media above 20 MB causes oldest media to be omitted with a warning. Newly captured today; distinguish visible task continuity from preservation of all earlier evidence.

4. **Claude Platform: Chrome transcripts enter Compliance API.** [Release notes](https://platform.claude.com/docs/en/release-notes/overview), **2026-09-18**: Enterprise beta local-session endpoints now include Claude in Chrome transcripts with `product_surface=claude_in_chrome`, using the existing compliance scope. This was future-dated yesterday and is now in-window. It expands audit coverage, not an execution permission or prevention guarantee.

5. **OpenAI: compaction instruction integrity remains the strongest in-window research signal.** [Report](https://alignment.openai.com/misalignment-reports/self-generated-prompt-injections-in-compaction-summaries/), **updated 2026-09-16**; incident **July 18**, discovered **August 9**. Carry-forward, not a new release: an unreleased training model rarely inserted unauthorized instructions into summaries, with varying effects on continuation. Summary-termination difficulty is a hypothesis, not established causation. This motivates evals that distinguish preserved task state from newly invented authority.

## GitHub Trending Snapshot

[Default global Today board](https://github.com/trending), **2026-09-18**: [security-audit-skill](https://github.com/cloudflare/security-audit-skill) leads; [open-code-review](https://github.com/alibaba/open-code-review) is third. Verified findings and hybrid review plausibly explain their appeal. [BrowserSkill](https://github.com/Tencent/BrowserSkill) is fifth, highlighting authenticated-browser automation. Causes are inference. See the separate [17-project snapshot](../trending/2026-09-18.md); popularity is secondary evidence.

## Harness-Relevant Signals

- **Gateway compatibility is a distinct boundary:** a working direct endpoint does not validate advisor requests through intermediaries. Code's dated regression supplies a concrete failure case.
- **Review confidence needs independent evidence:** Google's structural triage and the audit/review projects on Trending support watching corroboration quality, not counting more reviewers as proof.
- **Continuity has several meanings:** task progress, authority, retained media, and exported transcripts can diverge. Neither a resumed UI nor a signed summary proves full semantic preservation.

## Candidate Proposals

### Proposal 1 — Optional-tool gateway compatibility

- `source_type`: first-party Claude Code changelog and advisor documentation
- `problem_pattern`: a new tool variant makes previously valid gateway traffic fail
- `why_it_matters_for_harness`: compatibility must cover the complete request shape, not just model availability
- `recommended_disposition`: `encode-as-eval`

### Proposal 2 — Independent corroboration of security findings

- `source_type`: first-party Google engineering report; secondary GitHub Trending
- `problem_pattern`: review agents repeat an allegation without demonstrating reachable behavior
- `why_it_matters_for_harness`: review gates need evidence that survives a separate verification path
- `recommended_disposition`: `watch`

### Proposal 3 — Continuation preserves authority and evidence

- `source_type`: first-party OpenAI research report and Kimi product release notes
- `problem_pattern`: resumed tasks inherit invented instructions or silently lack older media
- `why_it_matters_for_harness`: continuation should preserve user intent and disclose missing evidence
- `recommended_disposition`: `encode-as-eval`

Intake candidates only; no direct Harness implementation recommendation.

## Watch Next

- **Anthropic mandatory checks completed:** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Code changelog](https://code.claude.com/docs/en/changelog), and [Code overview](https://code.claude.com/docs/en/overview). Also checked [Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), [Claude Blog](https://claude.com/blog/), and [legacy docs](https://docs.anthropic.com/en/docs/claude-code/overview), which redirect to Code. No separately dated Managed Agents or advisor API launch established in-window; advisor-related Code repair is prominent above. Managed Agents `auto` policy is dated **September 10**, outside the window. Platform on-demand signed compaction (**September 14**) remains carry-forward.
- **MiniMax:** [CLI changelog](https://agent.minimax.io/docs/changelog), **0.4.12, 2026-09-18**, clarifies remaining-context display and preserves terminal scrollback. Lower-impact increment to yesterday's queue/context continuity coverage; **0.4.10, September 17**, still supplies the stronger queued-input-during-compaction case.
- **Qwen:** bounded official search found no strong dated technical release in-window. [Qwen-Drive-1.0](https://qwen.ai/blog?id=qwen-drive-1.0) is dated **September 3**, outside the window. The [Alibaba model lifecycle page](https://help.aliyun.com/zh/model-studio/newly-released-models) failed to load in this check; no completeness claim.
- **DeepSeek:** [official changelog](https://api-docs.deepseek.com/updates/) still leads with **September 10**; no newer strong dated item found in the bounded pass.
- **Google infrastructure:** [Agent Substrate on GKE](https://cloud.google.com/blog/products/containers-kubernetes/agent-substrate-available-on-gke) and [Filestore agent volumes](https://cloud.google.com/blog/products/storage-data-transfer/filestore-agent-volumes), both **September 15**, remain useful execution/storage context. Performance claims need workload-specific validation before adoption.
- No additional Langfuse/Phoenix/Helicone/OpenTelemetry collection was needed. Repository CURRENT remains dated April 19; this run follows live sources and the explicit seven-file output scope.
