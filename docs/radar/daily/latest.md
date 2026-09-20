# Daily Radar - 2026-09-20

Window: 2026-09-14 to 2026-09-20 (America/Los_Angeles). Checked September 20 around 09:35 PDT; compared with September 19. Bounded collection: at most one search pass per publisher family, plus direct official surfaces. Newly highlighted means new to this daily comparison, not a September 20 launch. Vendor claims are not local validation.

## New Signals

1. **Google: persistent storage becomes an explicit agent infrastructure layer.** [Filestore agent volumes](https://cloud.google.com/blog/products/storage-data-transfer/filestore-agent-volumes), **2026-09-15**, newly highlighted. Google describes isolated persistent workspaces, rapid sandbox reattachment, lifecycle tiering, and shared RWX access with POSIX locking. Availability matters: open to all customers for non-production workloads; production GA access requires an allowlist. Harness interpretation: filesystem persistence and concurrent-write ownership need separate evidence; shared storage alone does not settle semantic edit conflicts.

2. **Moonshot: subagent work is becoming visible to users.** [Kimi Work 3.2.11](https://www.kimi.com/en/help/kimi-work/release-notes), **2026-09-18**, newly highlighted: the detail panel exposes subagent execution. Complementary carry-forward: [Kimi Code Desktop / CLI 2.0](https://www.kimi.com/code/docs/kimi-code/whats-new.html), **2026-09-17**, adds visible approvals and long-running task modes; CLI requests omit oldest media with a warning above 20 MB. Visible process, retained evidence, and successful completion remain different guarantees.

3. **Claude Code: permission routing and portable instructions remain the leading in-window change.** [Changelog](https://code.claude.com/docs/en/changelog), **2.1.278, 2026-09-19**, carry-forward; still the newest listed release today. Eligible auto-mode sessions default to a server classifier, warn on billed fallback, and expose classifier location in `/status`. **2.1.277, September 18**, adds AGENTS.md fallback when CLAUDE.md is absent, with Bedrock/Vertex/Foundry exclusions. Provider-specific precedence belongs in compatibility evaluation.

4. **Claude Platform: audit coverage and compaction are separate contracts.** [Release notes](https://platform.claude.com/docs/en/release-notes/overview), carry-forward: **2026-09-18** adds Chrome transcripts to Enterprise beta Compliance API local-session endpoints; **September 14** adds on-demand signed compaction blocks. Signatures attest provenance; they do not by themselves prove preservation of task constraints. No newer Platform entry was listed at capture.

5. **OpenAI: creation controls do not revoke existing credentials.** [API changelog](https://developers.openai.com/api/docs/changelog), **2026-09-15**, carry-forward. Organization/project controls restrict new key types or prevent creation; organization restrictions take precedence, while existing keys are unaffected. This is a useful distinction between provisioning policy and effective runtime access. No September 20 technical launch was established by the bounded pass and [news index](https://openai.com/news/).

## GitHub Trending Snapshot

[Default global Today board](https://github.com/trending): [ECC](https://github.com/affaan-m/ECC) leads (harness optimization); [agent-native](https://github.com/BuilderIO/agent-native) follows (agent apps). [security-audit-skill](https://github.com/cloudflare/security-audit-skill) and [cua](https://github.com/trycua/cua) surface audit verification and computer-use evaluation. These are relevance hypotheses, not proven causes. See [full board](../trending/2026-09-20.md); secondary evidence only.

## Harness-Relevant Signals

- **Resume correctness spans several layers:** storage reattachment, instruction loading, conversation compaction, and evidence availability each need distinct checks. Google's speed claims are vendor claims, not measured Harness performance.
- **Independent review evidence remains valuable:** Google's [security engineering report](https://cloud.google.com/blog/topics/systems/using-ai-agents-to-secure-google-infrastructure), **2026-09-18**, carry-forward, separates scanning from structural reachability triage and keeps proposed patches under human review. Different agent names alone do not establish independent evidence.
- **Model aliases are moving contracts:** DeepSeek's official sources disagree on a previously scheduled transition; record source freshness and actual model identity before treating an alias as a stable capability boundary. No live API inference was performed.

## Candidate Proposals

### Proposal 1 — Resume preserves workspace and authority

- `source_type`: first-party Google infrastructure announcement; Claude Platform release notes
- `problem_pattern`: a resumed task has files but loses constraints or evidence needed to continue correctly
- `why_it_matters_for_harness`: persistence is insufficient unless the resumed task retains the right authority and state
- `recommended_disposition`: `encode-as-eval`

### Proposal 2 — Subagent visibility versus completion proof

- `source_type`: first-party Kimi Work and Kimi Code release notes
- `problem_pattern`: progress displays are mistaken for verified subtask completion
- `why_it_matters_for_harness`: leader decisions need outcomes and artifacts as well as visible activity
- `recommended_disposition`: `watch`

### Proposal 3 — Alias migration evidence drift

- `source_type`: first-party DeepSeek announcement, changelog search result, and current pricing documentation
- `problem_pattern`: older migration notices contradict current model availability documentation
- `why_it_matters_for_harness`: model routing assumptions affect reproducibility and capability evaluation
- `recommended_disposition`: `watch`

Intake candidates only; no direct Harness implementation recommendation.

## Watch Next

- **Anthropic mandatory surfaces checked:** [Platform release notes](https://platform.claude.com/docs/en/release-notes/overview), [Managed Agents overview](https://platform.claude.com/docs/en/managed-agents/overview), [advisor tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/advisor-tool), [Code changelog](https://code.claude.com/docs/en/changelog), and [Code overview through legacy docs redirect](https://docs.anthropic.com/en/docs/claude-code/overview). Also checked [Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), and [Claude Blog](https://claude.com/blog/). No separately dated in-window Managed Agents or advisor API update established. Managed Agents auto-policy entry remains September 10, outside the window. Undated advisor guidance is context, not a new release.
- **MiniMax:** [CLI changelog](https://agent.minimax.io/docs/changelog), **0.4.12, September 18**, preserves scrollback and clarifies context status; **0.4.10, September 17**, queues messages during compaction. Carry-forward continuity evidence; no stronger in-window launch established in the bounded pass.
- **Qwen:** bounded search found no strong dated technical item in-window; [blog](https://qwen.ai/blog) returned no readable body. Coverage limitation, not evidence of no releases.
- **DeepSeek — current documentation conflict, not a newly dated launch:** [September 10 announcement](https://api-docs.deepseek.com/news/news260910/) schedules V4 Pro rerouting on September 14, but the [current pricing page](https://api-docs.deepseek.com/quick_start/pricing/) explicitly retains V4 Pro after September 14 with unchanged billing and lists its own model version. The [changelog](https://api-docs.deepseek.com/updates/) search result agrees with retention, though direct opening failed. Prefer the current service documentation for planning; exact revision date is unavailable and runtime routing remains untested. No strong newly published in-window technical item established.
- No extra Langfuse/Phoenix/Helicone/OpenTelemetry collection was needed. Repository CURRENT remains April 19; this run leaves it untouched under the explicit output boundary.
