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
- **still incomplete**: richer selector families beyond current normalized ownership transfer, authority adjudication against contract/live-git disagreement, deeper automated council orchestration, clearer public product contract

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
| `arch-review` meeting discipline | authority hierarchy, bounded roles, allowed/forbidden sources, cross-exam, rejected reasons | **Partially absorbed, now much stronger** | meeting `role_specs`, `locked_decisions`, `known_context`, rejected reasons, doctor checks, `meetings orchestrate`, signal-level authority adjudication, run-level authority audit | No automated multi-role debate orchestration yet; still human/leader-driven |
| Council-corrected synthesis | leader convenes, council refines, output becomes system truth | **Absorbed** | meeting objects, findings/action items, adoption into runs | Could still deepen explicit conflict handling and synthesis traces |
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
- runs now expose `authority-audit`, which explains higher-authority active-signal conflicts, meeting adoption gaps, and desired-state drift in one leader-facing output
- `doctor` now flags unresolved `run-authority-conflict` when fresher or higher-authority active signals disagree with current run state

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

The meeting protocol is now much stronger, but the system does not yet automatically run a structured debate workflow.

What is still missing:

- automated independent-review wave
- automated cross-examination routing
- explicit synthesis trace from conflicting positions to final accepted truth

### 3. Authority adjudication

Authority levels are now validated and surfaced, and the system now adjudicates signal supersession, same-claim meeting conflicts, and run-level active-signal disagreement. But it still does not consistently arbitrate every cross-surface conflict.

What is still missing:

- systematic high-vs-low authority conflict handling when contract intent, live git state, and run state disagree
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
