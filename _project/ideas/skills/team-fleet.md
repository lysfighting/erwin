# Idea: team-fleet — multi-mission teamwork orchestration

- **Stage**: Sprout 🌿
- **Created**: 2026-04-16
- **ADR**: (pending Plant promotion — requires n≥2 pattern validation)
- **Issue**: (pending build kickoff)

---

## Problem

Current teamwork covers the **single-mission lifecycle** completely: `/team-issue → /team-claim → /team-drive → /team-ship`, with `/team auto #N` as the one-shot wrapper. But as soon as you need to fan out to **multiple parallel missions in one session**, every load-bearing decision gets re-derived by hand:

- Mission Contract layout (front-matter + Sub-tasks + Success Criteria + AI Notes checkpoint directive)
- Per-mission worktree creation and upstream-detachment hygiene
- Teammate briefing format (hundreds of lines of shared-language boilerplate per teammate)
- Checkpoint semantics (when teammates stop, what they hand the team-lead, how the team-lead reviews)
- Cross-mission status sync into a control-plane issue

On 2026-04-16 the vi_agent team ran this pattern for the first time (3 parallel missions: #1298 / #1299 / #1302). Every one of the points above was **hand-executed, copy-pasted, and ad-hoc**. The session cost ~30 minutes of setup discussion before the first teammate even moved. It all worked — all 3 PRs merged — but it worked because the orchestrator (a Claude instance) remembered every trick, not because the pattern was encoded anywhere reusable.

### Why this is not covered by existing skills

| Existing | What it does | Why it isn't fleet |
|---|---|---|
| `/team auto #N` | One mission end-to-end: claim → drive → ship | Single-mission by construction; no concept of parallel lanes |
| `/team-claim` (batch via `list`) | Browse assigned missions, claim one | Interactive, one-at-a-time; no concurrent Contract generation |
| `/team-drive` team-mode | Spawns teammates *within* a mission (task-parallel) | Task-parallel inside one Contract; not mission-parallel across many Contracts |
| `/team queue` | Read-only mission list view | Dashboard only; no execution |
| `/create-milestone` | Batch-create MCs at milestone rollup | Creates paper Contracts; does **not** orchestrate execution |

**None** of them is a multi-mission orchestrator. The absence shows up as the 30-minute setup tax and the inconsistency of checkpoint expectations between teammates in the same session.

### Community gap (survey 2026-04-16)

Searched for prior art in the Claude Code / Codex CLI ecosystem + official Anthropic shipped patterns. Full details in the session appendix; summary:

| Project | Stars | Layer | Multi-mission swarm? | Nested lead? | Checkpoint gate? |
|---|---|---|---|---|---|
| Yeachan-Heo/oh-my-codex | ~23.5k | L3+ orchestration pack for Codex | No (staged pipeline) | No | No |
| Yeachan-Heo/oh-my-claudecode | ~29.3k | L3/L4 for Claude Code | No (staged pipeline: plan→prd→exec→verify→fix) | No | No |
| rekpero/claude-code-swarm (SwarmOps) | niche | L4 daemon | Yes (one agent per issue in worktree) | No (flat dispatcher) | No |
| ComposioHQ/agent-orchestrator | ~6.3k | L3 | Yes (flat) | No | No |
| nwiizo/ccswarm | ~135 | L3 | Yes (flat, worktree isolation) | No | No |
| dsifry/metaswarm | ~210 | L3/L4 | No (4-phase single-mission pipeline) | No | Review gate mid-pipeline, not start-gate |
| Anthropic official "Agent Teams" | — | L3 | Flat, shared task list + worktree | No | No |

Two axes are consistently missing across everything surveyed:
1. **Nested lead hierarchy** — team-lead spawning per-mission feature-leads that each run a full claim→drive→ship pipeline
2. **Mode-A first-checkpoint gating** — a uniform "teammate stops before writing code, hands team-lead a design note / research result, waits for approval" protocol

Our pattern composes both on top of the existing per-mission Contract discipline that teamwork already provides. That composition is the gap this proposal fills.

---

## Hypothesis

A sub-command `/team fleet` that takes a set of issues or a filter and orchestrates N parallel missions end-to-end, reusing as much existing teamwork machinery as possible, codifies the 2026-04-16 pattern into a repeatable operation.

### Naming — why sub-command, not top-level

`/team auto`, `/team doctor`, `/team queue`, `/team sync`, `/team wrap`, `/team share` are already sub-commands of `/team`. They are the **meta / admin / dashboard / orchestration** layer that sits above the per-mission lifecycle verbs (`team-claim`, `team-drive`, `team-ship`, `team-issue`, `team-rc`). `team-fleet` is semantically meta (it orchestrates lifecycle verbs; it isn't itself one) → **`/team fleet`**.

Precedent directly on point: `/team auto` is single-mission full automation. `/team fleet` would be multi-mission full automation. Symmetric, orthogonal, naming-preserving.

### Shape of the flow

Given a set of N missions (discovery method is one of the Key Questions below):

```
/team fleet resolves N issues
        ├── generates N Mission Contracts (reuses /team-claim machinery)
        ├── creates N persistent worktrees under .claude/worktrees/<slug>/
        ├── TeamCreate "<team-name>" + spawns N feature-lead teammates
        │       (subagent_type: "feature-lead", background mode)
        │
        │   For each mission independently:
        │     teammate cd <worktree> → verifies HEAD / .mission / Contract
        │     teammate runs /team-drive
        │     teammate STOPS at Contract's First checkpoint
        │     teammate SendMessage's team-lead with checkpoint deliverable
        │
        ├── Team-lead barrier: wait until all N checkpoints in hand
        ├── Team-lead reviews each checkpoint: APPROVE / REVISE / ESCALATE
        ├── Approved teammates resume: implement → run tests → SendMessage done
        ├── Team-lead authorizes PR per mission via /team-ship
        │       (lifecycle-respecting; never bypass with manual gh pr create)
        ├── After each merge: /team-ship done on behalf of the teammate
        │       (Issue → status:done, Contract archived to .teamwork/completed/)
        ├── Shutdown all teammates, TeamDelete
        └── Sync control-plane issue with lane table (optional, configurable)
```

### What `/team fleet` preserves from the 2026-04-16 run

Every load-bearing discipline that made the n=1 run succeed has to be embedded in the skill:

- **Persistent worktree per mission**, NOT `isolation:"worktree"` (that parameter creates a competing temporary worktree and defeats the pre-provisioning). The skill must build the worktree before spawning the teammate, and the teammate briefing must `cd` into the pre-built path.
- **Immediate `git branch --unset-upstream`** on each new mission branch right after `worktree add -b mission/<N>-<slug>`. New branches default to tracking `origin/pre-launch`, and one accidental `git push` would then push to `pre-launch` directly.
- **Mode-A first-checkpoint** as a hard protocol, not a suggestion. The Contract's AI Notes section explicitly instructs the teammate to SendMessage a named deliverable (design note / delta table / research result / etc.) and wait for team-lead approval before touching code.
- **Teammate challenge space protected**. In the n=1 run, lead-1298 caught a factual error in the Contract (`_enabled()` vs `_allow_emission`) and lead-1302 caught a back-compat trade-off in team-lead's Adjustment B. Both catches happened at checkpoint, **before any code was written**. That's the Mode-A payoff. The skill must encode "teammates are expected to push back on Contract errors, not silently comply."
- **Rebase onto fresh `origin/pre-launch` before opening PR**. Multi-mission windows advance the base fast — the n=1 run saw #1308 and #1309 merge between #1314's feature commit and its rebase. CI must run against current base.
- **Never bypass `/team-ship`**. In the n=1 run, the team-lead let teammates run manual `git push -u` + `gh pr create` to get precise control of PR body wording. That bypassed the Issue label lifecycle (`status:wip → status:review → status:done`) and the Contract archival step. It worked out (labels patched manually, Contracts unfortunately lost with worktree removal), but it was unambiguously a mistake. `/team fleet` must route PR creation through `/team-ship`. Any custom PR body wording is passed **into** `/team-ship` as a parameter, not around it.
- **`Refs #NNNN` vs `Closes #NNNN`** decision. If a mission's Sub-tasks include `🔧 MANUAL` items that can only be completed post-merge (ops execution, external registration, etc.), the PR must use `Refs`, not `Closes`, to keep the Issue open. `/team-ship` already does this based on unchecked Sub-tasks — `/team fleet` just needs to ensure Contracts mark MANUAL work correctly.

### What reuses vs extends vs is new

| Component | Status | Source |
|---|---|---|
| Mission Contract schema | **Reuse as-is** | `/team-claim` Step 4 template (teamwork v3.13.0) |
| Contract generation per issue | **Reuse** `/team-claim` subroutine | or programmatic variant (Key Question 3) |
| Worktree + branch setup | **Reuse** `tw-git.sh claim-setup` | already handles worktree routing + unset-upstream needs |
| Per-mission execution | **Reuse** `/team-drive` | no changes required |
| Per-mission delivery | **Reuse** `/team-ship` and `/team-ship done` | no changes required |
| Teammate briefing generator | **New** | wraps Contract + mission metadata → teammate prompt |
| Mode-A checkpoint protocol | **New (contract-level convention + runtime barrier)** | AI Notes section today is freeform; fleet needs structured checkpoint field |
| Multi-mission coordinator | **New** | the fleet operation itself; orchestrates TeamCreate / spawn / barrier / shutdown / TeamDelete |
| Cross-mission control-plane sync | **Optional / New** | e.g. `#1062`-style lane table; opt-in via flag |

---

## Key Questions

Unresolved design decisions, enumerated so reviewers can partition the promotion work. Each question's resolution shape should be clear enough to test.

### 1. CLI shape

**Options:**
- `/team fleet #1298 #1299 #1302` — explicit issue list (simple, n=1 shape)
- `/team fleet --label start-now` — label-driven batch (ergonomic for pre-curated waves)
- `/team fleet --from <file>` — spec file (supports cross-repo or annotated wave)
- `/team fleet` (no args, interactive select) — picks from queue with multi-select

**Recommendation for Sprout:** start with explicit issue list. Add label-driven batch in Plant. File-driven and interactive variants are nice-to-have, not required for first Fruit.

### 2. tw1 vs tw2 targeting

tw2-architecture is already Sprout→Plant and explicitly carves teamwork into `tw2 / tw2-issue / tw2-claim / tw2-drive / tw2-ship`. Multi-agent parallel execution is mentioned at `tw2-drive/operations/team-execution.md`, but that's **per-mission** task parallelism, not cross-mission fleet parallelism. This proposal does not overlap semantically, but does need an anchor:

- **Option A:** Ship on tw1 now as `/team fleet` sub-command, add dispatch line in current `team.md`. Sunset when tw2 ships, migrate to `tw2/operations/fleet.md` or `tw2-fleet`.
- **Option B:** Wait for tw2 to land, add directly as `tw2-fleet` in the new skill tree.

**Recommendation for Sprout:** defer the answer until Plant. Capture the question so the Plant ADR explicitly justifies the pick. Signal: if tw2 ships within N weeks, Option B wins; otherwise A.

### 3. Contract generation (batch vs interactive)

`/team-claim` is currently interactive (confirms Contract via AskUserQuestion, asks whether to create branch, etc.). Running it N times in a fleet context would multiply prompts by N — poor UX.

**Options:**
- Non-interactive mode in `/team-claim` (flag: `--yes` / `--batch`). Contract quality depends on the AI-generated shape being good enough without human confirmation.
- Programmatic Contract writer invoked directly by `/team fleet`, bypassing `/team-claim`. Risk: divergence from canonical Contract format.
- Interactive batch — `/team fleet` presents all N Contracts to the user for one-time bulk approval.

**Recommendation:** option 3 for v1 (reuses quality control, one approval round). Add non-interactive mode in Plant for CI/cron fleet kicks.

### 4. Mode-A checkpoint uniformity

The n=1 run observed **three different checkpoint shapes**:
- lead-1298: design note (exact file:line of gate insertion, env-var name rationale, test plan, doc-diff plan)
- lead-1299: sanitized delta table (current UIDs / target / gap / approved ops path)
- lead-1302: design note with entry-point citation + baseline artifact decision + one-command signature

This variety was **not** a mistake — each shape matched the mission's nature. But it means "Mode-A checkpoint" can't be a fixed schema; it has to be mission-specified.

**Options:**
- Add a `checkpoint:` field to Mission Contract front-matter that names the checkpoint type (`design-note` / `delta-table` / `research-result` / `ad-hoc`) and references a schema file.
- Leave it as free text in AI Notes; rely on team-lead review to catch misalignments.
- Hybrid: recommend a small vocabulary of checkpoint types, allow extension.

**Recommendation:** hybrid. Vocabulary seeds the intent; `ad-hoc` escape hatch keeps it flexible. Seed vocabulary: `design-note`, `delta-table`, `research-result`, `risk-inventory`, `ad-hoc`.

### 5. Feature-lead briefing generation

The n=1 briefings were hand-written per teammate, ~200 lines each, with ~80% boilerplate (Step 1 cd-verify, Step 2 read Contract, Step 3 invoke /team-drive, Step 4 checkpoint, hard rules, communication protocol) and ~20% mission-specific content (the Mode-A specific instructions + mission-specific discipline).

**Options:**
- Generator takes `(Contract, team_name, worktree_path)` → briefing string, with Mode-A section populated from Contract's AI Notes + hard rules from Non-Goals.
- Externalized template (already landed in vi_agent as `.claude/templates/feature-lead-briefing-template.md`). Fleet fills placeholders.
- Hybrid: ship template + generator; generator can be skipped for human review.

**Recommendation:** template-based generator. Every Fleet session fills the template; team-lead can inspect+edit before spawn.

### 6. Failure / recovery semantics

What if a teammate:
- Fails checkpoint 3 times (can't produce an approvable design note)?
- `/team-drive` errors out?
- Gets stuck in a loop?
- User interrupts the fleet mid-flight?

**Options:**
- TaskList-based state: each mission phase tracked as a task; on recovery, read tasks, infer where to resume.
- `.teamwork/fleet/<session>/state.json` persistent state file, rehydrated on rerun.
- Conservative first version: fail fast, surface to user, give explicit restart command (`/team fleet --resume <session-id>`).

**Recommendation:** option 3 for v1. Option 2 for Plant. Drive mode's file supremacy principle ("files always win over memory") applies.

### 7. `/team-ship` integration

How is PR creation orchestrated?

**Options:**
- Each teammate calls `/team-ship` on its own branch when team-lead approves. Team-lead provides PR body as a parameter.
- Team-lead collects all checkpoints-approved + code-done signals, then calls `/team-ship` serially or in parallel on each branch.
- Interactive approval gate before ship — user sees a fleet dashboard ("3 missions ready to ship; approve all / approve some / hold").

**Recommendation:** option 3. The fleet operation is meta; user approval at the ship gate matches the approval at the checkpoint gate in terms of critical decision density.

### 8. Cross-mission dependencies

n=1's three missions were independent. Future fleets might have dependency chains: mission B needs mission A merged before B can verify against it.

**Options:**
- Out of scope for v1; require user to launch fleets in phases.
- Declared dependency graph in Contract front-matter (`depends_on: [MISSION-A]`); fleet topologically sorts and runs in waves.
- Detect dependencies by scanning Contracts for cross-references.

**Recommendation:** out of scope for Sprout/Plant. Capture as future Fruit v2.

### 9. TeamCreate / TeamDelete lifecycle binding

Today's pattern required manual TeamCreate before spawning + manual TeamDelete after shutdown.

**Options:**
- `/team fleet` manages the team resource; user never sees TeamCreate / TeamDelete.
- User passes an existing `--team <name>`; fleet just spawns into it.
- Hybrid: default auto, `--team <name>` override.

**Recommendation:** hybrid.

### 10. Mode-A checkpoint metadata persistence

Each checkpoint produces a valuable artifact (the design note, the delta table, etc.) that is **lost** when the teammate shuts down. These are the high-leverage decision moments of the whole fleet.

**Options:**
- Persist checkpoint reports to `.teamwork/active/MISSION-<N>/checkpoint-<type>.md` or similar.
- Attach as PR body appendix so it ships with the mission.
- Post to a fleet-control issue like n=1 did on #1062, sanitized.

**Recommendation:** option 1 for Plant; option 3 for Fruit as optional flag.

---

## References / Prior Art

### This session's concrete artifacts (n=1 run, 2026-04-16)

Evidence for all claims above lives in the session appendix: `team-fleet-appendix-session-2026-04-16.md`.

**PRs merged (in order):**
- #1308 `feat(observability): controlled production step_evidence gate | Mission: #1298` → merge `daaf9844` (2026-04-16T11:35:24Z)
- #1309 `feat(benchmark): scheduled benchmark report entry point with baseline refs | Mission: #1302` → merge `0764b48e` (2026-04-16T11:35:39Z)
- #1314 `feat(loadtest): register_pool + preflight_pool scripts for #1189 Block B | Mission: #1299` → merge `56efe1ab` (2026-04-16T11:52:28Z)

**Issue outcomes:**
- #1298 CLOSED (Closes #1298 in PR body — tooling complete)
- #1302 CLOSED (Closes #1302 — tooling complete)
- #1299 OPEN (Refs #1299 — scripts landed; ops execution still pending at time of handoff)

**Control-plane sync on #1062** — 3 comments documenting next-wave claim / merge / handoff. See appendix.

**Adjacent merge in-window:** PR #1317 `fix(sse-relay): eliminate replay disconnect race + align first-chunk filter` also landed (→ `320a0ce5`). Not owned by this session but moves #1189 Block A from "blocked on #1257 fix" to "tooling-ready, execution-pending", so the handoff context for #1299 is cleaner than originally expected.

### Local pattern artifacts (already saved in vi_agent)

- `~/.claude/projects/-Users-liyasong-Project-collov-vi-agent/memory/feedback_teamwork_multimission_pattern.md` — memory capturing the pattern with key anchors + gap citation
- `~/.claude/projects/-Users-liyasong-Project-collov-vi-agent/memory/feedback_team_ship_lifecycle_lesson.md` — memory capturing the lesson that manual `gh push + pr create` bypasses `/team-ship` lifecycle
- `~/.claude/projects/-Users-liyasong-Project-collov-vi-agent/memory/feedback_worktree_isolation_agents.md` — memory amended with the "persistent worktree is the allowed exception" section
- `/Users/liyasong/Project/collov/vi_agent/.claude/templates/mission-contract-template.md` — reusable Contract skeleton for batch-hand-writing
- `/Users/liyasong/Project/collov/vi_agent/.claude/templates/feature-lead-briefing-template.md` — reusable teammate-briefing skeleton

### Community research (gap confirmation)

Raw report by the Explore agent lives in the appendix. Summary already in the Problem section above.

### Existing teamwork skills that `/team fleet` composes

All under `/Users/liyasong/Project/collov/erwin/skills/teamwork/`:
- `team.md` — dispatcher (would gain a new `fleet` sub-command routing)
- `team-claim.md` — per-mission Contract generation (fleet calls this N times, ideally in batch mode)
- `team-drive.md` — per-mission execution loop (teammates call this inside their worktrees; fleet doesn't touch it)
- `team-ship.md` — per-mission delivery (fleet invokes for PR creation + post-merge cleanup)

### Adjacent erwin proposals

- `_project/ideas/skills/tw2-architecture.md` (Sprout→Plant) — major teamwork refactor. `/team fleet` must coexist cleanly with this proposal; see Key Question 2.
- `_project/ideas/skills/cross-mc-knowledge-sharing.md` — parallel proposal for cross-MC knowledge transfer. If team-fleet runs N missions in the same session, cross-MC knowledge-sharing becomes easier to implement (all teammates visible to team-lead; team-lead can facilitate).

---

## Notes

### n=1 status — strong signal, not yet proof

This session is the **first** complete multi-mission swarm run. It worked — 3 PRs merged, 2 issues closed cleanly, 1 issue intentionally left open for ops. But one data point is not a pattern; it's a hypothesis that hasn't been falsified once. Before Plant:

- Run the pattern at least once more with a different mission mix (different sizes, different parallelism widths, at least one L-size mission to stress-test)
- Identify which parts of the n=1 protocol were accidental vs load-bearing
- Test failure paths (what happens when a checkpoint fails, when a rebase conflicts, when a teammate goes off-spec)

### Lessons embedded from n=1

These must all make it into the eventual skill's semantics — they are the hard-won invariants:

1. **Isolation parameter conflict**: do NOT pass `isolation:"worktree"` to Agent when worktrees are pre-provisioned. The `isolation` flag spawns a competing temporary worktree; the teammate ends up in a different checkout than the one you Contracts'd. Resolved in n=1 by manual `cd` in the briefing. See `feedback_worktree_isolation_agents.md` Exception section.

2. **Upstream-tracking booby-trap**: `git worktree add -b mission/<N>-<slug> <path> origin/pre-launch` leaves the new branch tracking `origin/pre-launch`. A naked `git push` from the teammate would then push directly to `pre-launch` — an explicit forbidden operation per `feedback_no_direct_push.md`. Fix: `git branch --unset-upstream` immediately after branch creation, in the fleet machinery.

3. **Rebase cadence**: n=1 saw `origin/pre-launch` advance by 5 commits between the first mission's branch cut and the last mission's PR open. CI must run on fresh base. `/team fleet` should rebase all branches onto current `origin/pre-launch` just before `/team-ship` (or `/team-ship sync`).

4. **Teammate challenge space**: the two most valuable catches in n=1 came from teammates pushing back on the Contract:
   - lead-1298 corrected Contract's assumed function name (`_allow_emission`) to the actual name (`_enabled`) after grepping — saved a buggy implementation.
   - lead-1302 flagged Adjustment B as over-literal and proposed a back-compat-preserving Option A — avoided breaking #1283's existing CLI contract.
   
   Both happened at the checkpoint, before code was written. If the Fleet encodes "teammate must comply" instead of "teammate must challenge on error", both catches would have been missed and bugs would have shipped. The briefing template must explicitly protect this.

5. **`/team-ship` bypass**: n=1 used manual `git push` + `gh pr create` to get precise PR body wording. Merge succeeded, but the Issue never got `status:review` → `status:done` automatically, and the Contracts were never archived (the worktree-remove step deleted them). Fleet must route through `/team-ship` — any custom body goes IN, not AROUND.

6. **Refs-vs-Closes signaling**: some missions have sub-tasks that can only complete post-merge (ops execution, external registration). Those missions need `Refs #N` so the Issue stays open. n=1 got this right manually for #1299. Fleet needs to encode this decision from the Contract's Sub-tasks (`🔧 MANUAL` markers).

### Risks

1. **Prematurely freezing the pattern at n=1**. This is the #1 reason this proposal is Sprout and not Plant. The right response to "this worked once" is "let's verify the invariants" — not "let's turn it into a rigid API that's hard to evolve when we discover we were wrong."

2. **Overlap with tw2 drive-team-execution**. tw2-architecture mentions "Multi-agent parallel execution" at `tw2-drive/operations/team-execution.md`. That's **per-mission** task-parallelism (subagents inside a Contract), not cross-mission fleet. The names collide in English but the semantics are orthogonal. Risk: confusion when someone reads both proposals. Mitigation: the fleet proposal must explicitly cite tw2-architecture and call out the orthogonality.

3. **Background-teammate context loss**. Once a teammate shuts down, its context is gone. Checkpoint reports (the highest-leverage artifacts) go with it unless persisted. This is Key Question 10; if we don't solve it, the value of the checkpoint gate is diluted by post-hoc unrepeatability.

4. **`/team-ship` coupling**. Fleet depends on `/team-ship` being robust. If `/team-ship` is the bottleneck or has bugs when invoked in parallel, the fleet inherits those. Plant-stage work should include a parallel `/team-ship` reliability test.

### Suggested phasing

1. **Sprout (now)**: this document. Capture problem + hypothesis + open design questions. No build.
2. **Plant**: once 1-2 more multi-mission runs confirm the invariants, write the ADR deciding tw1 vs tw2 target + resolving the 10 Key Questions. Committing to a concrete CLI shape and checkpoint vocabulary.
3. **Fruit**: build the skill. Either as `/team fleet` operation in current `team.md` (tw1) or as `tw2-fleet` / `tw2/operations/fleet.md` in the new tree (tw2). Ship with a reference run (our n=1 run as the canonical example).

### Naming reserved

**`team-fleet`** / **`/team fleet`** is the proposed name. Alternatives considered and rejected:
- `team-swarm` — rejected. "Swarm" is heavily used across competitors (claude-code-swarm, ccswarm, claude-swarm) and conflates with flat-dispatcher patterns. "Fleet" is mission-coordinated, hierarchy-implying — better semantic fit.
- `team-missions` — neutral but weaker; sounds like a dashboard, not an action.
- `team-wave` — closer to the drive-mode execution concept, but drive-mode already uses "wave" internally; would collide.
- `team-parallel` — descriptive but generic.

---

## Integration points with existing teamwork machinery

`/team fleet` is a composition, not a rewrite. Enumerating the exact reuse surface makes the Plant-stage scope tractable.

### Existing scripts under `skills/teamwork/scripts/` expected to be reused as-is

| Script | Used for | Invariant assumed |
|---|---|---|
| `tw-config.sh detect-dir` | Locate `.teamwork/` dir per worktree | Deterministic, already used by `/team-drive` and `/team-ship` |
| `tw-config.sh resolve-labels` | Get the project's label vocabulary (`status:wip`, `status:review`, `status:done`) | Single source of truth; fleet's label transitions must go through this, not hardcoded strings |
| `tw-git.sh slugify` | Turn Issue title → branch slug | Needed per mission when fleet auto-generates branch names |
| `tw-git.sh claim-setup` | Atomic branch + worktree setup + `.mission` marker + Contract copy | Fleet calls this N times in parallel; must be concurrent-safe (claim to verify) |
| `tw-git.sh protect-check` | Reject pushes from protected branches | Each `/team-ship` invocation inside fleet still runs this |
| `tw-contract.sh hash` | Compute `issue_content_hash` for Contract front-matter | Only used if fleet goes the "programmatic Contract write" path (Key Question 3); manual-claim-no-hash is legitimate fallback |
| `tw-contract.sh read-comments` | Filter Issue comments for Pre-Claim Discussion | Applied identically to fleet-mode batch claim |
| `tw-contract.sh store-comment-timestamp` | Store latest comment timestamp per Contract | Required after fleet-mode Contract generation |

### Existing skills invoked as subroutines

| Skill | Invoked by fleet for | Invocation count per fleet run |
|---|---|---|
| `/team-claim` (or programmatic equivalent per KQ3) | Contract generation per mission | N |
| `/team-drive` | Per-mission execution loop (runs inside each feature-lead teammate) | N (one per teammate) |
| `/team-ship` | PR creation + CI watch + status label flip + body formatting | N |
| `/team-ship done` | Post-merge Issue close + `status:done` label + Contract archival | N (after each merge) |

### New machinery owned by `/team fleet`

| Component | Responsibility |
|---|---|
| Fleet orchestrator | Parse CLI → resolve issue list → spawn barrier → shutdown barrier |
| Feature-lead briefing generator | Fill `feature-lead-briefing-template.md` placeholders from Contract + runtime |
| Checkpoint barrier | Wait for all N teammates to report first checkpoint; surface to user as batch review |
| Per-PR sequencing after merge | Interleave `/team-ship done` invocations with control-plane sync |
| Fleet session state (optional, Plant+) | Persist checkpoint deliverables + approval history + timeline to `.teamwork/fleet/<session-id>/` |

### Out-of-band dependencies that fleet inherits but doesn't own

- `.teamwork/config.yml` — must be valid; `/team init` covers this
- Labels `status:queued`, `status:wip`, `status:review`, `status:done` must exist on the GitHub repo — `/team init` provisions them
- `gh` CLI authenticated — `/team` skills already require this
- `TeamCreate` / `TeamDelete` / `SendMessage` tools — first-party Claude Code tools, no install surface
- `feature-lead` subagent type — provided by the project's agent roster; if missing, fleet falls back to `general-purpose`

---

## Acceptance criteria for Plant promotion

A concrete, falsifiable set of signals that justify promoting `team-fleet` from Sprout 🌿 to Plant 🌱. If these are not all met, the proposal stays Sprout and does not enter ADR stage.

1. **n ≥ 2 multi-mission runs** in vi_agent (or a comparable repo) using the hand-executed pattern, producing merged PRs, with session appendices preserved. One of these runs MUST include at least one of:
   - A mission requiring interrupt/recovery mid-flight (tests Key Question 6)
   - A mission with cross-dependency on another mission in the same fleet (tests Key Question 8)
   - A mission that fails checkpoint at least once and requires team-lead revise (tests Key Question 4)
2. **Every n=1 lesson in this Sprout still holds**. If new runs invalidate any of the "Lessons embedded" items in Notes, the Sprout must be revised before Plant.
3. **Key Question 2 (tw1 vs tw2) has an answer**. The answer may be "ship on tw1, migrate at tw2 cutover" or "wait for tw2" — but the Plant ADR must pick one, not both.
4. **CLI shape picked** (Key Question 1). At least the explicit-issue-list form (`/team fleet #A #B #C`) is locked. Other forms may remain post-Plant work.
5. **Checkpoint vocabulary picked** (Key Question 4). Seed vocabulary agreed; at least 2 of the vocabulary items have schema files in erwin's `skills/teamwork/reference/` (or equivalent).
6. **Briefing generator spec written** (Key Question 5). Interface specified as `(Contract, team_name, worktree_path) → briefing_string`; template referenced.
7. **Recovery semantics decided** (Key Question 6). At minimum, "restart-from-scratch on failure" is acceptable for v1 Fruit, but the decision must be explicit.
8. **No unresolved contradictions with tw2-architecture**. Cross-reference check: the two proposals should either cleanly coexist (if tw2 accepts the fleet sub-command) or explicitly cohabit with an agreed deprecation path.
9. **Author-time cost measured**. One multi-mission session's Claude-time cost WITHOUT `/team fleet` (i.e. hand-executed) vs WITH prototype `/team fleet` must be within expectation band (hours, not days — per project's working memory of Claude-execution cadence). If the skill doesn't save author time, the cost/benefit fails.

Plant promotion writes ADR `NNNN-team-fleet-multimission-orchestration.md` under `_project/decisions/`.

---

## UX design principles for the skill

These are Claude-UX principles the skill's authoring must respect. Extracted from working-memory conventions discovered during n=1 and earlier sessions — they apply to every interactive decision point the fleet will surface.

- **Present tradeoff tables before asking**, not a question with hidden trade-offs. Every multi-option decision point (CLI-shape choice, path selection, body wording, label conflict, recovery strategy) must show a table of options with pros, cons, and recommendation **before** the AskUserQuestion call. This principle was reinforced by the vi_agent team's working memory and confirms the n=1 pattern where team-lead wrote explicit "MY POSITION" + "MY POSITION CHALLENGE" blocks before every AskUserQuestion.

- **Name the cost in the same units as user's time**. Claude-executed work is measured in **hours**, not days. When the skill presents phase estimates or retry costs, use hours and concrete step counts — not abstract "days of engineering" framings inherited from human-team planning cadence.

- **Show scan-ready output first, then details on request**. Fleet status should default to a one-table dashboard (lane / state / blocker); teammate checkpoint content, PR body details, and ops-handoff notes live behind explicit drilldowns. The 2026-04-16 session's #1062 control-plane comments are the right shape — each comment scannable in 15 seconds.

- **Treat user interruption as first-class, not exceptional**. The skill must be safe to `Ctrl-C` at any phase. Recovery from `.teamwork/fleet/<session-id>/` state file (Plant-stage addition) is the path; no "you must complete the fleet" lock-ins.

- **Announce every dispatch and completion**. Per the drive-mode "no silent agents" principle, every teammate spawn, checkpoint receipt, approval send, push, PR open, and merge must emit a short broadcast line. Fleet owes the user transparency in proportion to the parallelism width.

- **Never bypass downstream lifecycle skills to save a step**. Even if `/team-ship` feels like overhead when the PR body is fine-tuned, the skill must route through it — `feedback_team_ship_lifecycle_lesson.md` from the n=1 run is the direct counter-example. Any customization (PR body, label override) is a `/team-ship` parameter, not a sidestep.

---

## Discovery metadata

- Session: 2026-04-16, vi_agent repo, Claude instance orchestrating
- Team used: `next-wave-startnow`, 3 feature-lead teammates
- Full session narrative: `team-fleet-appendix-session-2026-04-16.md` (sibling file)
- Authored in dialogue with lys (user) — the "pattern" is equal parts emergent tactics and deliberate decisions captured in the appendix

<!-- When this proposal promotes to Plant, update Stage, link the ADR, and add the Issue ref. When it reaches Fruit, move to skills/<name>/ per erwin graduation pipeline. -->
