# Appendix: team-fleet n=1 Session Evidence Pack (2026-04-16)

> Companion to `team-fleet.md`. This file preserves the full session narrative — the decisions, checkpoint reports, timeline, and raw artifacts that substantiate the proposal. The proposal cites this file in References; scan this file when the proposal's claims need back-up.

---

## Meta

- **Date**: 2026-04-16 (UTC)
- **Repo**: `flair-home-stylist/vi_agent`
- **Base branch**: `pre-launch`
- **Session operator**: lys (Liya Song), working through a Claude Code instance as team-lead
- **Team resource**: `next-wave-startnow` (TeamCreate)
- **Teammates spawned**: `lead-1298`, `lead-1299`, `lead-1302` (subagent_type `feature-lead`, run_in_background: true)
- **Session window**: roughly 11:02 UTC (first Contract write) → 13:01 UTC (last teammate shutdown confirmed)
- **Outcome**: 3 PRs merged to pre-launch; 2 issues closed; 1 issue deliberately left open for ops execution

---

## Mission set

| Issue | Title | Priority | Size | Assigned teammate |
|---|---|---|---|---|
| #1298 | feat(observability): controlled production step_evidence gate | P1 | M | lead-1298 |
| #1299 | feat(loadtest): prod loadtest user pool and rerun prep | P1 | M | lead-1299 |
| #1302 | feat(benchmark): scheduled benchmark report run with baseline refs | P2 | M | lead-1302 |

Queued (not started in this session): #1300 (trace_id propagation MVP), #1301 (DashboardIntent → DashboardResolution MVP).

---

## Timeline (select decision points)

### 11:02 — Manual claim path chosen over `/team-claim`

lys had pre-created all 5 next-wave issues before the session; wanted to claim 3 of them (start-now set) and leave 2 queued.

Initial team-lead behavior: `gh issue edit --add-assignee lysfighting --add-label status:wip --remove-label status:queued` + `git worktree add -b mission/<N>-<slug> <path> origin/pre-launch` per issue. No Contract file written.

**Problem surfaced:** `/team-drive` requires a Contract file at `.teamwork/active/MISSION-<N>.md` + a `.mission` marker. Without both, `/team-drive` errors at Step 0 and stops. The team-lead had built worktrees but not Contracts.

**Resolution:** team-lead chose Path ② ("write Contracts by hand into each pre-built worktree") over Path ① ("tear down and re-run via `/team-claim` interactively"). Path ① would have re-created branches with different name suffix (`-{user}`) and required 3 interactive rounds of AskUserQuestion. Path ② preserved the already-claimed state.

**Precedent:** `/Users/liyasong/Project/collov/vi_agent/.teamwork/active/MISSION-1303.md` used `issue_content_hash: "manual-claim-no-hash"` — same pattern of manual Contract authoring was already validated in the repo.

**Implication for team-fleet:** the skill will need to decide whether to go interactive (N × `/team-claim` prompts) or non-interactive (programmatic Contract write, like Path ②). Path ② is faster but bypasses the interactive quality checks in `/team-claim`. See Key Question 3 in the proposal.

### 11:02–11:07 — Manual Contract authoring (3 files)

Team-lead read `/team-claim.md` Step 4 template, then authored 3 Contracts directly:
- `.claude/worktrees/step-evidence-prod-gate/.teamwork/active/MISSION-1298.md`
- `.claude/worktrees/prod-loadtest-user-pool/.teamwork/active/MISSION-1299.md`
- `.claude/worktrees/benchmark-scheduled-baseline/.teamwork/active/MISSION-1302.md`

Each followed the full schema (front-matter + Objective + Sub-tasks + Success Criteria + Context Files + Test Command + AI Notes). Each also got a `.mission` file (single-line issue number) in its worktree root.

Context investigation per mission:
- **#1298**: `api-server/app/obs/step_evidence.py` (259 lines) located; `ENVIRONMENT=production` guard identified at line 149-151; existing `PROFILE_LOG_FORCE_ENABLE` readers grepped (auth.py, profile_log.py).
- **#1299**: `loadtest/README.md` located as entry-point; secret material (tokens-prod.json) intentionally NOT read; ops path marked as "to be discovered" not "to be guessed."
- **#1302**: `benchmark_report/__main__.py` + `builder.py` + `render.py` + existing `tests/test_builder.py` and `tests/test_render.py` located; `BenchmarkReport` semantics confirmed intact.

### 11:07 — TeamCreate + 3 feature-lead spawns

```
TeamCreate:
  team_name: "next-wave-startnow"
  agent_type: "team-lead"
  description: "Drive 3 start-now issues with Mode-A first-checkpoint gates"
```

Three `Agent` calls in parallel:
```
Agent(subagent_type: "feature-lead",
      team_name: "next-wave-startnow",
      name: "lead-1298"|"lead-1299"|"lead-1302",
      run_in_background: true,
      prompt: <per-mission briefing, ~200 lines each>)
```

**Not passed**: `isolation: "worktree"`. Passing it would have created temporary worktrees competing with the pre-provisioned ones. See `feedback_worktree_isolation_agents.md` Exception.

Each briefing had the shape later codified in `.claude/templates/feature-lead-briefing-template.md`:
- Step 1: `cd <worktree>` + `git rev-parse --abbrev-ref HEAD` + `cat .mission` + `ls .teamwork/active/` (3-way sanity verification)
- Step 2: read Contract
- Step 3: invoke `/team-drive` via Skill tool
- Step 4: stop at first checkpoint, SendMessage team-lead with deliverable
- Hard rules (branch discipline, GitHub write discipline, mission-specific non-goals)
- Communication protocol (SendMessage tool only)

### 11:08–11:09 — All 3 checkpoints arrived (in-order: 1298 → 1299 → 1302)

**lead-1298 (11:08:54):** 4-section design note.
- Adjustment flagged: Contract assumed function name `_allow_emission`; actual code `_enabled()` at `step_evidence.py:148`. Pushed back with grep evidence.
- Env-var choice: `STEP_EVIDENCE_FORCE_ENABLE`, mirroring `PROFILE_LOG_FORCE_ENABLE` truthy parsing exactly.
- Test plan: extend existing `test_obs_step_evidence.py`, do not create a new file (found the existing file with 209 lines and 18 tests via grep).
- Coupling regression guard: proposed `test_production_with_profile_force_does_not_emit` as an explicit Non-Goal pin.
- Doc-diff plan: 5 specific files with line-range edits.
- Scope: explicitly flagged `realtime/src/observability/step_evidence.py` exists but NOT in this mission's scope; proposed follow-up MC if realtime prod visibility needed.

**lead-1299 (11:09:06):** sanitized delta table + open questions list.
- Did not attempt to enumerate real prod UIDs; cited public comment evidence for 3-UID baseline.
- "Approved ops path: unknown — needs lys to specify" — correctly flagged rather than guessed.
- Proposed Path 1 (signup flow via `POST /api/auth/firebase`) vs Path 2 (`kprod exec` direct DB insert) with trade-off analysis.
- Five open questions for lys (naming, invite_required, registration executor, preflight executor, whether mission ships code).
- Secret discipline: explicit self-audit that no tokens/UIDs/path values were read, echoed, or committed.

**lead-1302 (11:09:13):** design note with key finding.
- Key finding: #1283 surface already wired baseline-delta support end-to-end; `build_benchmark_report` already populates `baseline_ref` and `deltas_vs_baseline`. Mission is smaller than Contract assumed.
- Baseline artifact proposal: commit a small fixture tree under `benchmark_report/baselines/` + a `regenerate.py` that emits from `tests/fixtures.py` — "single source of truth" avoiding golden-file refresh anti-pattern.
- One-command invocation: `python -m benchmark_report --baseline ... --current ... --baseline-report ...`.
- Test plan: 3 delta cases (identical/regress/improve) + 1 fixture-exists integration test.
- **Explicit commitment**: `BenchmarkReport` semantic surface stays 100% untouched.
- Open question flagged: `run_scheduled.sh` location — `benchmark_report/` vs `scripts/`. Recommended former.

### 11:09–11:20 — Team-lead review + approvals

Team-lead reviewed each checkpoint, wrote back via SendMessage:

**Approval for #1298** (11:10s range): accepted design, added two adjustments:
- A: scope-boundary call-out in docstring + docs (explicit "api-server only; realtime remains ENVIRONMENT-gated; follow-up MC if needed")
- B: test setup must `monkeypatch.delenv` BOTH `PROFILE_LOG_FORCE_ENABLE` and `STEP_EVIDENCE_FORCE_ENABLE` to pin against shell-inherited env

**Partial approval for #1299**: path selected (Path 1 — signup flow), Path 2 explicitly disallowed (prod `users` table blast radius). Private artifact path locked (`~/.codex/tmp/monitoring-control-2026-04-16/artifacts/loadtest-1299/`). Next step: **read-only recipe discovery** (extended Mode A) to determine whether a signup recipe already exists. Goal: answer single question "does a usable loadtest signup recipe exist, and if so where?"

**Approval for #1302** (11:10s range): accepted design, added three adjustments:
- A: `run_scheduled.sh` location confirmed at `benchmark_report/`, shebang `#!/usr/bin/env bash` mandated, `chmod +x` in same commit.
- B: default `--out-dir=benchmark_report/out/` pinned; `.gitignore` updated if needed.
- C: `regenerate.py` must be deterministic; audit `fixtures.py` for non-deterministic values; use `json.dump(..., sort_keys=True, indent=2)`; add drift-guard test that calls `regenerate.py` as subprocess and compares bytes.

### 11:20 — lead-1299 checkpoint 2 (recipe discovery result)

Not Case A (recipe exists cleanly) and not Case B (no recipe). Case A-partial:

- Token minter exists: `loadtest/auth/mint_tokens.py` (canonical, complete, well-documented with self-declared invariant at lines 33-36 that it does NOT write to the users table)
- Signup endpoint exists: `POST /api/auth/firebase` at `api-server/app/routes/auth.py:96-372` — public API, uses the same two headers loadtest runner already sets, rate-limited to `30/minute` per IP
- **Gap**: no packaged script wrapper between "mint tokens" and "drive signup for each token"

Proposed 6-phase ops plan (P1 sizing decision → P6 sanitized public posts) with explicit executor split per phase.

Five open questions returned to team-lead:
1. UID naming (extend `loadtest-user-{3..24}` vs fresh prefix)
2. prod `invite_required` setting
3. Registration executor (ops vs teammate with redacted tokens)
4. Preflight executor
5. Whether mission ships code

### 11:20–11:21 — lead-1298 implementation complete

Two commits on branch:
- `100d6c56` — feat(observability): add independent prod gate for step_evidence
- `7ac04012` — docs(observability): document STEP_EVIDENCE_FORCE_ENABLE gate

Test run: `python3 -m pytest tests/test_obs_step_evidence.py -q` → **37 passed in 0.02s** (18 pre-existing + 19 new parametrized).

Self-reported Non-Goal compliance audit:
- Zero GitOps changes
- Zero changes to `_emit_line` / `_safe_push_loki` / DaemonSet paths
- Zero coupling to `PROFILE_LOG_FORCE_ENABLE` (grep evidence)
- Zero realtime / gateway / nanoclaw edits
- `git diff --check` clean

Flagged: holding, awaiting PR authorization.

### 11:24 — lead-1302 implementation complete

Single commit on branch:
- `f175d7ae` — feat(benchmark): scheduled benchmark report entry point with baseline refs (12 files, +626/-4)

Test run: `pytest tests -x --tb=short` → **13 passed in 0.12s** (8 pre-existing + 4 new delta cases + 1 regenerate drift-guard subprocess test).

All three adjustments applied:
- A: `run_scheduled.sh` at `benchmark_report/`, `#!/usr/bin/env bash`, mode 100755.
- B: default `--out-dir=benchmark_report/out/`; CLI preserves original stdout-only behavior when no flags given (this was an uninstructed back-compat preservation that became Adjustment A' discussion later).
- C: `fixtures.py` audited clean (all frozen literals); `regenerate.py` uses `sort_keys=True, indent=2`; drift-guard test passes.

Flagged: holding, awaiting PR authorization.

### 11:25 — Redundant TaskList scan from lead-1302

lead-1302 sent several idle-notification messages re-confirming task-list items already completed in `f175d7ae`. Identified as the team's shared TaskList auto-scanning behavior. Not a failure — informational redundancy.

**Observation for team-fleet:** the shared TaskList + auto-scanning is helpful for audit but noisy. v1 of fleet should suppress or route these into a separate log, not the team-lead inbox.

### 11:25 — lead-1302 flagged Adjustment B over-literal (back-compat catch)

lead-1302 discovered during Task 4 re-check that the Adjustment B's literal interpretation — "raw CLI without flags also defaults to writing files" — would break `#1283`'s existing CLI default (stdout-only JSON), risking external caller regressions (Makefiles, CI scripts, docs snippets).

Proposed Option A (preserve back-compat: raw CLI stays stdout-only, `run_scheduled.sh` always passes `--out-dir`) vs Option B (take Adjustment B literally, break back-compat).

**Team-lead's decision: Option A.** Explicit reasoning sent back: "Your back-compat instinct is correct — my Adjustment B was over-literal. The real value of Adjustment B was scheduled-job path uniformity, which is already achieved by the wrapper always passing `--out-dir`. Keep `f175d7ae` code; update README to document the CLI vs wrapper split."

### 11:27 — lead-1302 README clarification commit

- `6401da51` — docs(benchmark): clarify CLI vs wrapper out-dir default (81 lines added)

New commit, not amend. Final branch state: `6401da51` tip + `f175d7ae` feat.

### 11:29 — Multiple TaskList final reports from lead-1302

Final consolidated report re-sent (Task 9 at `6401da51`, 13/13 pass, `git diff --check` clean). Same informational redundancy as earlier.

### 11:32 — lead-1298 PR opened → #1308

Pushed + `gh pr create --base pre-launch --head mission/1298-step-evidence-prod-gate`.
PR body: Summary, Behavior table, Test plan, Non-Goal compliance, Documentation, Related / Closes #1298.

**URL:** https://github.com/flair-home-stylist/vi_agent/pull/1308

### 11:32 — lead-1302 PR opened → #1309

Pushed + `gh pr create --base pre-launch --head mission/1302-benchmark-scheduled-baseline`.

**URL:** https://github.com/flair-home-stylist/vi_agent/pull/1309

### 11:35 — #1308 + #1309 merged

- #1308 merged at `daaf9844` (11:35:24Z); #1298 CLOSED via `Closes #1298`
- #1309 merged at `0764b48e` (11:35:39Z); #1302 CLOSED via `Closes #1302`

### 11:37 — lead-1298 + lead-1302 shutdown

`SendMessage` with `{"type":"shutdown_request", "reason": "..."}`. Both approved and terminated cleanly.

### 11:38 — lead-1299 script delivery complete

Single commit `837c9801` — `feat(loadtest): register_pool.py + preflight_pool.py for #1189 Block B prep`.

6 files, +1190 lines:
- `loadtest/auth/register_pool.py` (248 lines)
- `loadtest/auth/preflight_pool.py` (229 lines)
- `loadtest/auth/_pool_common.py` (162 lines)
- `loadtest/auth/tests/test_register_pool.py` (293 lines)
- `loadtest/auth/tests/test_preflight_pool.py` (258 lines)
- `loadtest/auth/tests/__init__.py`

Test run: `pytest loadtest/auth/tests -q` → **36 passed in 0.10s**.

Self-audit of secret-discipline hard-rules, each mapped to specific test assertions that use `capsys` to confirm token/UID text never appears in stdout.

**Flagged diff-stat caveat**: `git diff --stat origin/pre-launch..HEAD` showed additional unrelated modifications because the branch was cut from an old `pre-launch` before #1298 / #1302 merged. Team-lead proposed rebase-before-push.

### 11:47 — Team-lead proposed rebase + README addition + `Refs #1299`

Adjustments to lead-1299's delivery before PR:
1. Add `loadtest/auth/README.md` (~60 lines, operator-facing; keep it lean)
2. `git rebase origin/pre-launch` to pull in latest base (base had advanced by 5+ commits)
3. PR body uses `Refs #1299` (not `Closes`) — Issue must stay open for pending ops execution

### 11:48 — lead-1299 PR opened → #1314

Post-rebase SHAs:
- `b269de19` (feat) replaced `837c9801`
- `b1c9eb57` (docs — README) new

**URL:** https://github.com/flair-home-stylist/vi_agent/pull/1314

Post-rebase pytest: still 36/36 green. README came in at 81 lines (20 above the ~60 target), flagged — team-lead accepted over-target because operator docs more > less.

### 11:52 — #1314 merged

- #1314 merged at `56efe1ab` (11:52:28Z); #1299 intentionally left OPEN (Refs).

### 12:12 — Adjacent merge in-window: PR #1317

Not this session's work, but relevant: `fix(sse-relay): eliminate replay disconnect race + align first-chunk filter` merged at `320a0ce5`. PR body said `Refs #1257` (not `Closes`) because Phase-3 verification was still pending. Bumped #1189 Block A from "blocked on #1257 fix" to "tooling-ready, execution-pending", which made the overall control-plane picture cleaner than originally forecast.

### 13:01 — lead-1299 shutdown confirmed

Shutdown request sent at 11:52:28Z after merge; teammate acknowledged and terminated at 13:01:39Z.

---

## Post-cleanup

After the last merge + shutdown:

- `git worktree remove` × 3 (step-evidence-prod-gate, benchmark-scheduled-baseline, prod-loadtest-user-pool)
- `git branch -D` × 3 (mission/1298-..., mission/1302-..., mission/1299-...)
- remote branches: auto-deleted by GitHub
- `TeamDelete` cleared `~/.claude/teams/next-wave-startnow/` and the shared TaskList

**Contracts NOT archived to `.teamwork/completed/`** — they were deleted along with the worktrees. This is the single unfixable lesson from the session: we bypassed `/team-ship done` (manual PR creation meant `/team-ship`'s full lifecycle never kicked in), so Contract archival never happened, and by the time we noticed, the files were gone with the worktree. Recoverable only if we'd used `/team-ship`. See `feedback_team_ship_lifecycle_lesson.md`.

---

## Control-plane sync: #1062 comments

Three comments were posted on program-control issue #1062 through the session:

### 2026-04-16T10:36Z — Claim sync

Table: Issue / Owner-session / Worktree / Branch / First checkpoint. Noted all 3 branches unset-upstream'd. Confirmed #1300 / #1301 still queued, #1189 still blocked on #1257.

### 2026-04-16T11:38Z — Merge sync (after #1308 + #1309)

Announced #1298 / #1302 merged; #1299 in flight for script delivery. Re-stated #1300 / #1301 queued, #1257 still WIP.

### 2026-04-16T11:55Z — Final next-wave sync (after all 3 merges)

All 3 lanes landed. Noted adjacent #1317 merge as a bonus: both #1189 Blocks (A and B) now "tooling-ready, execution-pending". #1300 / #1301 still queued. Listed remaining lys-owned items: prod `invite_required`, Block B ops execution window, #1257 Phase-3 verification.

---

## Community gap research (full)

### oh-my-codex (Yeachan-Heo, ~23.5k★)

**URL:** https://github.com/Yeachan-Heo/oh-my-codex

**Positioning:** Multi-layer orchestration pack for OpenAI Codex CLI. 33 prompts, 36 skills, HUD overlay, hooks. `$team N:executor` for parallel workers with auto-worktrees. Layer: L3 (orchestration templates) + partial L4 (CLI `omx`).

**Overlap with team-fleet:** weak. Per-worker worktree isolation + coordinator-executor split, but:
- No "team-lead spawns feature-leads per mission" hierarchy — flat executor pool
- Zero mission/claim/drive/ship vocabulary
- No first-checkpoint gating — staged pipeline model (one phase at a time)

### oh-my-claudecode (Yeachan-Heo, ~29.3k★)

**URL:** https://github.com/Yeachan-Heo/oh-my-claudecode

**Positioning:** Claude Code counterpart to oh-my-codex by same author. 19 agents, 36 skills. Cross-model support (Claude/Gemini/Codex) via `omc team`, `/autopilot`, `/ralph`, `/ultrawork`, `/ccg` tri-model synthesis.

**Overlap with team-fleet:** weak. Full pipeline: `team-plan → team-prd → team-exec → team-verify → team-fix`. That's a staged single-mission pipeline, not multi-mission fan-out.

- No per-issue mission ownership; workers are per-phase executors
- Vocabulary collision (`team-*` name space) is coincidental; semantics differ
- No Mode-A first-checkpoint gate
- No team-lead / feature-lead nesting

### rekpero/claude-code-swarm (SwarmOps)

**URL:** https://github.com/rekpero/claude-code-swarm (niche, recent, ~0 stars)

**Positioning:** 24/7 daemon watching a GitHub repo for open issues; dispatches one Claude Code agent per issue in a sibling worktree; PR monitor auto-dispatches fix agents on review comments.

**Overlap with team-fleet:** moderate on structure, weak on semantics. This is the **closest structural isomorph** found:
- One worker per issue in its own worktree ✓
- Coordinated by a central dispatcher ✓

But:
- Flat dispatcher — no team-lead/feature-lead nesting
- No first-checkpoint gate — proceeds immediately on `@swarmops start` comment or label
- Per-agent workflow is ad-hoc `claude -p` prompt, not a disciplined claim→drive→ship contract

It's a daemon competitor to a human team-lead, not a skill layered on top of existing mission skills.

### Other adjacents surveyed

- **ComposioHQ/agent-orchestrator** (~6.3k★) — flat orchestrator, per-issue worktree, no lead-hierarchy, no gate. Structurally similar.
- **nwiizo/ccswarm** (~135★) — worktree isolation + specialized agents, no mission-contract pipeline.
- **dsifry/metaswarm** (~210★) — 4-phase pipeline (IMPLEMENT → VALIDATE → ADVERSARIAL REVIEW → COMMIT) with parallel design-review gate. Has a gate concept, but mid-pipeline review gate, not start-gate; single-mission.
- **SpillwaveSolutions/parallel-worktrees** (7★) + **skills.rest parallel-worktree-orchestrator** — bare L1-L2 skills ("spawn N subagents in worktrees, merge"), no mission/contract vocabulary, no hierarchy.
- **Anthropic official "Agent Teams"** (experimental in Claude Code) — orchestrator + subagents with shared task list and worktrees. Flat, no contract discipline, no gate.

### Verdict

No project in the surveyed set combines all three of:
1. Per-mission Contract discipline (claim→drive→ship pipeline over an Issue)
2. Nested lead hierarchy (team-lead orchestrates per-mission feature-leads)
3. Mode-A first-checkpoint gate (teammates stop before code; team-lead approves)

`team-fleet` is a defensible, non-duplicative addition.

---

## Checkpoint reports — raw text

### lead-1298 checkpoint (11:08:54Z)

Four sections in response to team-lead's Mode-A instruction:

**Section 1: Env-var name `STEP_EVIDENCE_FORCE_ENABLE`.** Rationale: exact parallel to existing `PROFILE_LOG_FORCE_ENABLE` (same `_FORCE_ENABLE` suffix, same truthy parser, same semantics — prod defaults off, this forces on). Operators already know the pattern. Independent env var = zero coupling. Will mirror `profile_log.py:29` truthy parsing exactly.

**Section 2: Gate insertion point — Contract had two minor errors.** Contract said `_allow_emission` at `step_evidence.py:149-151`. Actual code: function name is `_enabled()` at line 148; gate read at line 150 (`env = (os.environ.get("ENVIRONMENT") or "dev").lower()`); return at line 151; single call site at line 211. Proposed patch (minimal, module-level, mirrors `profile_log.py:29-31`). Also flagged: `realtime/src/observability/step_evidence.py` exists but Contract scopes only api-server — this creates a temporary inconsistency, recommend realtime as follow-up MC, not this mission.

**Section 3: Test plan — EXTEND, do not create.** `api-server/tests/test_obs_step_evidence.py` already exists (209 lines, `TestEmitStepEvidence` class at line 133). Plan: extend with three new methods:
- `test_production_with_force_enable_emits`
- `test_production_with_profile_force_does_not_emit` ← coupling-regression guard pinning Non-Goal #5
- `test_force_enable_truthy_values` (parametrized over `("1","true","yes","on","TRUE","On")` + falsy)

Plus update existing `test_production_disabled` with `monkeypatch.delenv("STEP_EVIDENCE_FORCE_ENABLE", raising=False)`. No new test file, no new fixtures.

**Section 4: Doc-diff plan** — 4 edits across 3 files + new preamble + `.env.example` addition.

---

### lead-1299 checkpoint 1 (11:09:06Z) — sanitized delta

Table provided publicly-evidenced only:

| Verified today | Target | Gap | Approved ops path |
|---|---|---|---|
| 3 prod loadtest UIDs (class `loadtest-user-{low-index}`) | ≥25 prod loadtest UIDs | 22 UIDs | **unknown — needs lys to specify** |

Evidence sources cited from public GitHub:
- #1189 Block B state snapshot (2026-04-16)
- #1189 hotfix-1 rerun comment (2026-04-15)
- #1239 cross-MC comment on #1189 (2026-04-14)
- `loadtest/README.md` §1 and §5

Two candidate ops paths named in #1189 Block B (neither declared canonical):
1. Register via signup/registration flow — closer to real user shape; needs a reusable recipe
2. Direct DB insert via `kprod exec` with ops approval — faster but bypasses invariants

AI self-classification (advisory, not a recommendation to act): Path 1 closer to real user shape; Path 2 blast radius = prod users table, requires explicit ops sign-off + change management.

Hard-rule compliance self-attested per turn (tokens/UIDs/creds untouched, no `gh issue comment`, no push, no PR, no prod writes).

Asked 3 questions back:
1. Approved ops path (signup vs DB insert vs something else)
2. Who performs registration (lys / ops / teammate with guidance)
3. Where private artifacts should land (off-repo path)

### lead-1302 checkpoint (11:09:13Z)

**Key finding:** #1283 surface ALREADY wired baseline-delta support end-to-end. `build_benchmark_report` (line 73) already populates `baseline_ref` and `deltas_vs_baseline`. `_delta_direction` (line 422) already emits `improved|regressed|mixed|unchanged`. Mission is smaller than Contract assumed — no new `BenchmarkReport` semantics needed.

**Builder entry point:** `benchmark_report/__main__.py:main` (CLI) + `benchmark_report/builder.py:build_benchmark_report` (core) + `build_report_from_paths` (line 52, path-level wrapper).

**Baseline artifact proposal:** commit a small version-controlled fixture tree under `benchmark_report/baselines/`:
- `README.md` — refresh procedure
- `product_artifact.json` — minimal passing ProductTestArtifact
- `loadtest_artifact.json` — minimal passing LoadTestArtifact
- `baseline_report.json` — prior `BenchmarkReport` dict
- `regenerate.py` — deterministic emitter from `tests/fixtures.py`

Why version-controlled: deterministic (reproducible on any box/runner), reviewable (baseline shift lands as PR diff), re-creatable (all 3 JSONs derivable from fixtures.py → sidesteps "golden files nobody knows how to refresh").

**One-command invocation:**
```bash
python -m benchmark_report \
  --benchmark-id "scheduled-$(date -u +%Y%m%dT%H%M%SZ)" \
  --benchmark-label "scheduled baseline sanity" \
  --product-artifact benchmark_report/baselines/product_artifact.json \
  --loadtest-artifact benchmark_report/baselines/loadtest_artifact.json \
  --baseline-report benchmark_report/baselines/baseline_report.json \
  --json-output out/benchmark-report.json \
  --html-output out/benchmark-report.html
```

**Wrapper:** `benchmark_report/run_scheduled.sh` thin bash wrapper, no env branching.

**Test plan:** extend `benchmark_report/tests/test_builder.py` with 3 new delta cases (identical/regress/improve) + 1 integration test loading `baselines/*.json`.

**Semantic surface:** 100% untouched. No changes to `build_benchmark_report` signature, no new report dict fields, no changes to `render.py`.

---

## Session decisions rendered as principles

From the narrative above, these are the decisions worth extracting as invariants for the future skill:

1. **Manual Contract authoring is a legitimate option when `/team-claim`'s interactive flow would multiply prompts by N.** The `issue_content_hash: "manual-claim-no-hash"` convention is already validated in the repo. But the hand-written Contract must match the template exactly.

2. **Worktree pre-provisioning is orthogonal to `isolation:"worktree"`.** If you build the worktree yourself, do NOT pass `isolation:"worktree"` to Agent — the parameter creates a competing temporary worktree. Briefing must `cd` into the pre-built path and verify via 3-way check (HEAD, `.mission`, Contract).

3. **Mode-A first-checkpoint gate catches real errors.** In n=1, it caught:
   - Contract assumption error (function name)
   - Adjustment over-literal (back-compat)
   Both caught BEFORE any code was written. Without the gate, both would have shipped.

4. **Teammate challenge space must be explicitly protected.** Teammates default to compliance. The briefing must say "push back on Contract errors; flag back-compat risks" — otherwise the gate's value is forfeit.

5. **`/team-ship` is non-negotiable for merge lifecycle.** Bypassing it skips label lifecycle (status:wip → status:review → status:done) and Contract archival. n=1 bypassed and paid the cost (labels patched manually, Contracts permanently lost with worktree removal). Fleet must route all PRs through `/team-ship`.

6. **Rebase just before ship, not at start.** In a multi-mission window, `origin/pre-launch` advances while missions are in flight. Rebasing at start is waste; rebasing just before `gh pr create` (i.e. before CI kicks) gives the freshest base without wasted cycles.

7. **`Refs #N` vs `Closes #N` is decided by the Contract's Sub-tasks.** If unchecked items are `🔧 MANUAL`, use `Refs` so the Issue stays open for post-merge ops. If all Sub-tasks are code and all are checked, use `Closes`. `/team-ship` already handles this — fleet must ensure Contracts mark MANUAL work correctly.

8. **Background-teammate context disappears on shutdown.** Checkpoint reports are the session's highest-leverage artifacts. Without persistence, they become unrepeatable. Fleet v1 should persist at least the checkpoint deliverables to disk.

9. **The cost of runtime improvisation is real.** n=1 session spent approximately 30 minutes of initial-setup discussion (choosing between claim paths, investigating why `/team-drive` wouldn't start, deciding Contract format, drafting teammate briefings). That's the tax the fleet skill removes.

10. **Adjacent work can improve your control-plane picture.** PR #1317 (not this session) merged in-window and moved #1189 Block A from blocked to tooling-ready. Fleet sessions should scan the base-branch log post-merge to catch adjacent progress and update control-plane sync accordingly.

---

## Artifact links

- Feature PRs: [#1308](https://github.com/flair-home-stylist/vi_agent/pull/1308) · [#1309](https://github.com/flair-home-stylist/vi_agent/pull/1309) · [#1314](https://github.com/flair-home-stylist/vi_agent/pull/1314)
- Issues: [#1298](https://github.com/flair-home-stylist/vi_agent/issues/1298) CLOSED · [#1302](https://github.com/flair-home-stylist/vi_agent/issues/1302) CLOSED · [#1299](https://github.com/flair-home-stylist/vi_agent/issues/1299) OPEN
- Control plane: [#1062](https://github.com/flair-home-stylist/vi_agent/issues/1062)
- Adjacent: [#1317](https://github.com/flair-home-stylist/vi_agent/pull/1317) · [#1257](https://github.com/flair-home-stylist/vi_agent/issues/1257) · [#1189](https://github.com/flair-home-stylist/vi_agent/issues/1189)
- Local memories: `feedback_teamwork_multimission_pattern.md`, `feedback_team_ship_lifecycle_lesson.md`, `feedback_worktree_isolation_agents.md` (Exception section)
- Local templates: `.claude/templates/mission-contract-template.md`, `.claude/templates/feature-lead-briefing-template.md`

<!-- End of appendix. If more session data is needed beyond what's here, fetch raw PR bodies / issue comments via gh or reconstruct from git log on pre-launch between 2026-04-16 11:00 and 13:05 UTC. -->
