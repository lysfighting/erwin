# Harness Gap Audit — 2026-03-29

> Purpose: review the Codex-native harness against the earlier design goals so future sessions do not keep re-discovering what is already absorbed, what is only partial, and what is still missing.

---

## Executive Summary

The harness is no longer just a repo-local prototype.
It now has:

- a Go runtime core
- stable shell surfaces
- leader-led control-plane behavior
- first-class teammate runtime
- governed sibling communication
- council-grade meeting state
- packaged-runtime fallback and validation

But not every earlier design input is fully realized.

The current truth is:

- **absorbed well**: protocol lifecycle, leader-centric control plane, teammate runtime, governed communication, runtime packaging direction
- **absorbed partially**: cybernetic reconciliation, Kubernetes-style ownership discipline, harness-wide authority enforcement, automated council behavior
- **still incomplete**: richer selector families beyond current normalized ownership transfer, higher-order authority adjudication when live git state contradicts contract intent beyond the new clean-intent-vs-dirty-worktree baseline, deeper automated council orchestration beyond the new auto-synthesis and route-obligation baseline, clearer public product contract

---

## Audit Table

| Design source | Intended value | Current status | Where it lives now | Remaining gap |
|---|---|---|---|---|
| `teamwork` | Git/GitHub-ordered lifecycle, controlled execution, release governance | **Absorbed** | `issue -> claim -> drive -> ship -> rc`, `ship`, `rc`, `doctor`, `bootstrap`, `sync` | Product naming and public contract still unresolved |
| `wing` | Skill-shaped packaging, modularity, adapter thinking, less prompt-monolith | **Absorbed** | `skills/harness/`, `harness/adapter-manifest.yml`, `scripts/harness-*.sh`, bootstrap/sync | Still more “adapter codename” than finished product |
| Engineering cybernetics | desired vs actual state, feedback, correction, regulator model | **Absorbed much more strongly** | `home`, `doctor`, `maintenance reconcile`, run desired state, first-class mission registry and mission graph | Still room for richer desired-state reasoning across multi-mission bundles |
| Kubernetes control-plane thinking | controllers, ownership boundaries, reconciliation, validation gates | **Partially absorbed, stronger now** | leader/teammate split, manifest validation, stale-state reconciliation, ownership claims, selector-aware ownership transfer, overlap guardrails | Richer selector families and more nuanced transfer policies are still incomplete |
| Leader as chairing regulator | leader owns orchestration, not truth monopoly | **Absorbed** | leader-facing surfaces, `doctor`, `home`, meetings, teammate assignment | Still room for richer scheduling / mission-graph control |
| Teammates as bounded execution units | delegated work with scope, reporting, escalation | **Absorbed** | `teammates assign/report/message/escalate` | No native Codex sub-agent binding yet |
| Teammate-to-teammate communication | lateral communication under governance | **Absorbed** | sibling links + governed `teammates message` | Could later gain freshness/routing policies beyond current signal model |
| `arch-review` meeting discipline | authority hierarchy, bounded roles, allowed/forbidden sources, cross-exam, rejected reasons | **Partially absorbed, now much stronger** | meeting `role_specs`, `locked_decisions`, `known_context`, rejected reasons, doctor checks, `meetings orchestrate`, signal-level authority adjudication, run-level authority audit, claim/dimension-aware auto-synthesis | No automated multi-role debate wave yet; still leader-guided |
| Council-corrected synthesis | leader convenes, council refines, output becomes system truth | **Absorbed more strongly** | meeting objects, traced findings/action items, auto-synthesized provisional findings, adoption into runs | Could still deepen explicit conflict handling and debate-wave automation |
| Runtime packaging / harness engineering | installability, portability, no project-language coupling | **Absorbed** | Go runtime, packaged binary preference, runtime-dist smoke | Public binary distribution strategy still open |

---

## What Is Newly Closed In This Pass

This audit was followed immediately by implementation work that materially closes one of the biggest remaining protocol gaps:

- meeting `role_specs` are now first-class runtime objects
- each role can carry `mission`, `allowed_sources`, and `forbidden_sources`
- meetings can record `known_context`
- meetings can record `locked_decisions`
- rejected findings can now carry explicit rejection reasons plus evidence and authority level
- `doctor` now checks for missing role specs, missing expected output, missing context, missing locked decisions in advanced phases, and rejected findings without reasons
- runs and contracts now carry explicit `ownership_claims`
- active exclusive overlap is now rejected during teammate assignment and surfaced in `doctor`
- runs now carry declared `desired_state`
- dependency waiting is now first-class via `depends_on`, `drive wait`, and home/doctor dependency visibility
- meetings can now auto-advance phases with `meetings orchestrate`
- runs now expose `authority-audit`, which explains higher-authority active-signal conflicts, meeting adoption gaps, desired-state drift, contract drift, and live-state drift in one leader-facing output
- `doctor` now flags unresolved `run-authority-conflict` when fresher or higher-authority active signals disagree with current run state
- run-level control-plane mutations now refresh contract snapshots, and `maintenance normalize` can repair stale historical contracts
- meeting findings can now carry explicit provenance back to supporting/challenging positions, and synthesis audit reports traced vs untraced findings
- `meetings add-position` now accepts `--claim=` and `--dimension=` so councils can group positions by governed claim instead of only by freeform summary
- `meetings orchestrate` can now auto-synthesize provisional findings and unresolved follow-up actions from grouped positions without overwriting manual findings already on the record
- `meetings summary` now surfaces auto-generated finding/action counts, making council automation visible to the leader
- `runAuthorityAudit` now flags a dirty live worktree when run/desired state already claims the execution should be clean enough for governed verification or handoff
- worktree authority snapshots now ignore the harness-owned `.mission` marker so live dirty detection tracks real execution drift rather than harness bookkeeping noise
- meetings now materialize a `debate_plan` that groups conflicting claim lanes, assigns bounded challenger route obligations, and reports whether linked cross-examination has actually been routed for each conflict
- `doctor` now treats unrouted conflicting claims as a governance problem once the meeting has advanced into cross-examination, synthesis, or output

This pushes the council model much closer to the intended `arch-review` discipline instead of only “meeting exists.”

---

## What Is Still Not Done

### 1. Ownership transfer model

Explicit ownership claims now exist, and selector families are now richer than the original normalized baseline, but the model is still not fully mature.

What is still missing:

- richer selector policies beyond the new family-aware baseline (`file` / `dir` / `module` / semantic exact-match)
- more nuanced overlap resolution than exclusive/shared
- transfer policies beyond the current bounded `move` semantics

### 2. Automated council orchestration

The meeting protocol is now much stronger, and the system can now auto-synthesize provisional findings from grouped positions plus assign and track bounded route obligations for conflicting claim lanes, but it still does not yet automatically run a full structured debate workflow.

What is still missing:

- automated independent-review wave
- automated cross-examination routing
- automatic independent-review wave planning
- automatic multi-round cross-examination assignment beyond the new bounded route-obligation tracking
- automatic synthesis beyond the current provisional position-group baseline; the harness can now derive provisional findings and unresolved follow-up actions, but still depends on leader/council judgment for the richer debate workflow

### 3. Authority adjudication

Authority levels are now validated and surfaced, and the system now adjudicates signal supersession, same-claim meeting conflicts, and run-level active-signal disagreement. But it still does not consistently arbitrate every cross-surface conflict.

What is still missing:

- systematic high-vs-low authority conflict handling when live git state contradicts contract intent beyond the new dirty-worktree-vs-clean-intent baseline
- clearer winner/loser rules when run summary, contract, and live verification disagree without a governing signal or meeting
- broader leader-facing explanations of why a lower-authority artifact was overridden outside the run/signal/meeting path

### 4. Public product contract

The adapter is strong, but the external shape is still a prototype line.

What is still missing:

- final product name
- public package contract decision
- final binary distribution strategy

---

## Practical Review Rule For Future Sessions

Before implementing new harness behavior:

1. Read `_project/CURRENT.md`
2. Read this audit
3. Check whether the target capability is already:
   - absorbed
   - partial
   - missing
4. Prefer closing a named gap over adding a new surface

This is how the project avoids “more features” while still missing the original mission.
