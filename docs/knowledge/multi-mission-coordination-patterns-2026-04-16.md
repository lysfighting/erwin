# Multi-Mission Coordination Patterns

> Date: 2026-04-16 (two source sessions) · recorded 2026-04-17
>
> Purpose: capture the durable invariants that two April 2026 multi-mission AI-orchestration sessions taught the repo, without ossifying them into a dedicated skill before the evidence justifies one.

---

## Bottom Line

Two vi_agent sessions in mid-April 2026 — `team-fleet n=1` on 2026-04-16 and `teammate-work-2` shortly after — converged on a set of discipline rules for AI-orchestrated multi-mission work. The rules are real enough to record, **not yet real enough to encode as a skill**.

The first session proposed a `/team fleet` skill framed around **parallel independent missions**. The second session showed that **serial mission chains are also natural** — and that the common shape underneath both is a **mission DAG** with the team-lead as shared-context carrier across nodes.

So: preserve the pattern, drop the skill ambition. If a dedicated skill becomes justified later, it should be framed around the DAG abstraction, not around "parallel fleet."

---

## Source Sessions

**Session 1 — team-fleet n=1 (parallel)**

- Date: 2026-04-16 · Repo: `flair-home-stylist/vi_agent`
- 3 parallel missions: `#1298` (step_evidence gate), `#1299` (loadtest pool), `#1302` (benchmark baseline)
- PRs `#1308` / `#1309` / `#1314` all merged (verified via `gh pr view` 2026-04-17: SHAs `daaf9844` / `0764b48e` / `56efe1ab`)
- Original Sprout: `_project/ideas/skills/team-fleet.md` (now marked RELOCATED — see its header for decision trail)
- Evidence pack: `_project/ideas/skills/team-fleet-appendix-session-2026-04-16.md`

**Session 2 — teammate-work-2 (serial chain + mixed)**

- Mission chain: `#1329` → `#1330` → `#1333`. Output of one mission became the input of the next (alembic gap discovery propagated).
- Saved once by teammate pushback on an incorrect Amendment 2 (documented lesson below).
- Small changes (`#1335` config 10 lines, `#1340` env var 4 lines) were done directly by team-lead without spawning teammates.
- 5 distilled lessons captured durably at `~/.claude/projects/-Users-liyasong-Project-collov-erwin/memory/feedback_multimission_work_wisdom.md`.

---

## Invariants Observed Across Both Sessions

### 1. Worktree per mission, pre-provisioned — and `isolation:"worktree"` is a trap if used together

When the orchestrator pre-builds a worktree per mission, passing `isolation:"worktree"` to `Agent` spawns a **second** competing worktree and leaves the teammate in the wrong checkout. The fix is discipline, not tooling:

- Pre-provision the worktree.
- Briefing must `cd <path>` + verify HEAD, `.mission` marker, and Contract file (3-way sanity).
- Immediately run `git branch --unset-upstream` on new mission branches — new branches default to tracking `origin/pre-launch` (or whatever base branch), and a naked `git push` from the teammate would violate no-direct-push rules.

Recorded in the exception section of vi_agent's `feedback_worktree_isolation_agents.md`.

### 2. Mode-A first-checkpoint is the orchestrator-specific add, not a new challenge discipline

`/team-drive.md` L476–L524 already contains:

- L476 **Anti-Compliance Core Rule** — AGREE+EXTEND / CHALLENGE+PROPOSE / REDIRECT, never silent compliance.
- L506 **Sub-tasks are a suggested route — challenge, revise, or confirm them.**
- L524 **Challenge each sub-task — ALL types.**

This is the **content layer** of teammate challenge. It is not fleet-specific.

What multi-mission orchestration actually adds is a **form layer**: the teammate **stops before code**, sends a structured checkpoint deliverable (design note / delta table / research result) to the team-lead, and **waits for batch approval** alongside other teammates' checkpoints. Single-mission `/team-drive` has no team-lead to report to; the external handoff gate is the new semantic.

Briefing templates should **cite** L476-L524 for the challenge content and **mandate** the external pause/SendMessage form, not reinvent either.

### 3. Team-lead verification discipline — grep in the worktree, not in main

When a teammate sends back a fact claim ("Contract says function X, actual name is Y"), the team-lead **must verify** before acting on it. Otherwise the Amendment built on top rests on a false premise, and every subsequent decision compounds the error.

Critical detail: the grep **must run in the teammate's worktree**, not in the main checkout. Different worktrees have different code states; a grep in main can give the wrong answer.

Session 2 hit this exact failure (Amendment 2 was built on an unverified teammate claim and had to be retracted).

### 4. Pushback is bidirectional — protect both directions in the briefing

- **Teammate pushes back on Contract errors** — well-known, covered in Mode-A design.
- **Teammate pushes back on team-lead's Amendment errors** — less obvious, but session 2 showed this direction is equally load-bearing. Team-lead's subsequent clarifications are also subject to scrutiny; a teammate who silently complies after a team-lead Amendment is failing the sage role.

Briefing templates should explicitly say: "team-lead's Amendments and clarifications are within your review scope. Push back when you see logic or fact errors."

### 5. Spawn threshold — ≥3 files / ≥30 lines / cross-module

Spawning a teammate has real overhead (briefing authorship, checkpoint wait, batch review). For small changes, the team-lead does it faster directly.

- Below the threshold: team-lead executes directly in a worktree.
- Above the threshold: spawn a teammate, justify the overhead.

**The worktree itself is orthogonal to the spawn decision.** Even 4-line changes benefit from a worktree — isolation of blast radius, clean rollback, uncontaminated main checkout.

### 6. Topology is a DAG — parallelism is just the edge-empty case

Session 1 exposed **N parallel independent missions**. Session 2 exposed **serial chains** where mission output fed into mission input. Neither is the whole pattern; both are edges of a general DAG:

- **Parallel wave**: zero-dependency missions fire simultaneously.
- **Serial chain**: a mission's output defines or unblocks the next mission.
- **Mixed**: some parallel, some serial, DAG-shaped execution.

The team-lead is the unique **shared-context carrier** across DAG nodes. In pure parallel execution, teammates are isolated and cross-mission discoveries don't propagate. In serial/mixed execution, the team-lead inherently carries context between nodes.

Any future skill for this pattern should be framed around the DAG abstraction. Framing it as "parallel fleet" is already too narrow.

### 7. `/team-ship` bypass destroys lifecycle signal

Skipping `/team-ship` in favor of manual `git push` + `gh pr create` loses two things:

- **Issue label lifecycle** (`status:wip` → `status:review` → `status:done`) is not flipped automatically.
- **Contract archival** to `.teamwork/completed/` does not happen — and Contracts are lost permanently when the worktree is removed.

Session 1 took this shortcut and paid both costs. Rule: route all PR creation through `/team-ship`. If custom PR body wording is needed, that customization must live in the Contract's AI Notes before ship, not around `/team-ship`.

### 8. Rebase just before ship, not at start

Multi-mission windows advance the base branch fast — session 1 saw `origin/pre-launch` advance by 5 commits between first branch cut and last PR open. Rebasing at mission start is wasted effort (the base has moved again by ship time). Rebasing just before `gh pr create` gives the freshest base without churn.

---

## Why This Is a Knowledge Note, Not a Skill

Review session `_project/archive/sessions/2026-04-17-B.md` documented the decision chain in full. Abstract:

- **erwin's Skills lane policy** is "keep only what still drives present work" (`GOALS.md` L57–58). Teamwork is explicitly legacy.
- **The original Sprout had 9 framing/integration/premise gaps** (Gap 1 factual error, Gap 2 scope confusion, Gap 3 missing opposite test, Gap 4 unmeasurable AC, Gap 5 concurrency spike, Gap 7 skill-integration mismatches, Gap 8 spawn threshold, Gap 9 scope decision).
- **Four strong Plant-readiness dimensions were entirely missing** from the Sprout: security (teammate exposure to sensitive Contract data), cost (per-teammate token budget), post-hoc auditability (approval reasoning persistence), concurrency (N-parallel script race behavior).
- **n=1 was insufficient** by the Sprout's own Acceptance Criteria, which required n≥2 plus specific failure path coverage.
- **Session 2 data broadened the pattern beyond fleet's parallel-only framing** — the DAG insight made clear that any skill should be framed around DAG execution, not around fleet specifics.

The review verdict was **REDIRECT → HOLD → lys policy**, with RELOCATE recommended. Lys selected RELOCATE. This note is that relocation.

---

## Honest Caveats

1. **n=2, not n=10.** These patterns are strong signals. They are not laws. A third session could invalidate any of the invariants above.
2. **vi_agent-specific details bled through.** `pre-launch` as base branch, specific label taxonomy, specific Contract schema. Other repos adopting the patterns need their own calibration.
3. **Team-lead attention scaling is untested.** Batch review works at N=3 (session 1) and serial N=3 (session 2). At N=8 or N=15, the cognitive load may force a different structure. No data yet.
4. **Concurrency is analyzed, not spiked.** `tw-git.sh claim-setup` was labeled "claim to verify" concurrent-safe in the original Sprout; post-review script analysis said "LIKELY SAFE" for distinct missions in worktree mode but flagged this explicitly as needing a real n=3 concurrent spike before any production-grade skill would be responsible.
5. **Production-grade concerns are entirely deferred.** Security (sensitive data isolation), cost (token budget), auditability (approval trail persistence), attention budget (N-parallel human ceiling), role diversity (beyond feature-lead), portability (cross-repo config surface), maintainer view (N worktree + N PR cleanup trail legibility) — all Plant-level concerns not addressed here. If this pattern ever tips toward a real skill, these must be resolved first.

---

## References

- Review session: `_project/archive/sessions/2026-04-17-B.md` — 9 rounds of adversarial audit, verdict chain (REDIRECT → HOLD → RELOCATE).
- Original Sprout (RELOCATED): `_project/ideas/skills/team-fleet.md` — 439-line exploration of a `/team fleet` skill proposal. Kept as historical exploration; header points here for current state.
- Evidence pack: `_project/ideas/skills/team-fleet-appendix-session-2026-04-16.md` — 487-line session 1 timeline, checkpoint reports, community gap research.
- Session 2 lessons: `~/.claude/projects/-Users-liyasong-Project-collov-erwin/memory/feedback_multimission_work_wisdom.md` — 5 rules from teammate-work-2, durable memory.
- Existing challenge content layer: `skills/teamwork/team-drive.md` L476–L524 — Anti-Compliance Core Rule + sub-task challenge + all-types challenge. The content layer this pattern composes on top of.
- Cross-MC knowledge sharing (already Fruit): `_project/ideas/skills/cross-mc-knowledge-sharing.md` + `/team sync` + `/team share` in `skills/teamwork/team.md`. The persistence mechanism (GitHub issue comments with anti-plague safeguards) is the reusable substrate for any checkpoint-artifact persistence this pattern would need.

---

## Implications for Future Work

- **If multi-mission coordination hits a recurring constraint** that the patterns above do not resolve, revisit. A third data point may shift the calculus toward Plant.
- **If erwin's Skills lane scope expands** to include teamwork extensions as a first-class theme again (not currently the case per `GOALS.md`), revisit.
- **If a different repo adopts the pattern and discovers new invariants or breakages**, extract those as a new case — not as bolt-ons to a `/team fleet` skill that never shipped.
- **Do not re-file this as a Sprout.** If the need recurs, the right framing is a fresh idea entry around **mission DAG execution**, not around fleet-specific parallel orchestration. Don't resurrect the narrow frame.
