# Wing Parity Gap Analysis

> Date: 2026-03-28
>
> Objective: identify what `wing` still lacks before it can replace `teamwork` completely, while explicitly preventing regression into monolithic prompt architecture or other historical baggage.

---

## Goal Definition

`wing` reaches parity when:

- every workflow that matters in `teamwork` can be completed in `wing`
- the same protocol guarantees are preserved
- the mechanism layer is at least as reliable
- the new architecture remains modular and reference-driven
- no parity work reintroduces the old prompt bloat

This is **workflow parity**, not line parity.

---

## Non-Goal Definition

The following are not parity goals:

- restoring 4k-line skill files
- copying teamwork docs into runtime prompt files
- preserving legacy naming where a cleaner wing-native name already exists
- keeping compatibility forever just because it exists today
- porting every historical workaround when the root cause no longer exists in wing

---

## Gap Summary

### Category A — Already solved in Wing

These are not gaps:

1. Six-skill lifecycle architecture
2. Official Anthropic `SKILL.md` packaging
3. Operation extraction for heavy subflows
4. Shared VE reference
5. Script-layer preservation
6. Skill-local path resolution
7. Canonical `.wing/` state home
8. Backward-compatible config detection

### Category B — Functional parity likely present, but proof is incomplete

1. Full claim -> drive -> ship -> rc lifecycle parity
2. Worktree-aware mission routing parity
3. Comment freshness and Contract freshness behaviors
4. Side-effect triage correctness
5. PR review and RC promote/hotfix correctness

### Category C — Capability or packaging gaps that still matter

1. Missing suite-level documentation inside `skills/wing/`
2. Missing suite-level changelog/history packaging
3. No explicit wing e2e test entry comparable to `tw-e2e-test.sh`
4. No explicit parity matrix documenting teamwork-only vs wing-complete behavior
5. No written migration contract for when teamwork can be deleted

---

## Detailed Gap Matrix

| Domain | Teamwork status | Wing status | Gap type | Priority | Recommendation |
|---|---|---|---|---|---|
| Main dashboard/init | Mature | Present | Validation/docs | P1 | Keep skill compact; add suite docs, not more prompt text |
| Issue creation/update/fix/batch | Mature | Present with operations | Validation/docs | P1 | Add parity tests and user docs |
| Claim flow | Mature | Present | Validation | P0 | Verify active-contract, worktree, freshness, and confirmation parity |
| Drive flow | Mature | Present | Validation | P0 | Verify sequential/parallel/full-team routing and manual-task handling parity |
| Ship flow | Mature | Present with review op | Validation | P0 | Verify all pre-flight branches and side-effect triage cases |
| RC flow | Mature | Present with promote/hotfix ops | Validation | P0 | Verify cancel/fix/promote/hotfix end-to-end |
| VE guidance | Mature but duplicated | Centralized | Solved | — | Keep single-source model |
| Deterministic scripts | Mature | Mostly migrated | Small feature audit | P1 | Confirm no teamwork-only script capability remains required |
| README | Present | Missing in `skills/wing/` | Packaging gap | P1 | Add suite README |
| ONBOARDING | Present | Missing in `skills/wing/` | Packaging gap | P1 | Add onboarding focused on wing-native workflow |
| CHANGELOG | Present | Missing in `skills/wing/` | Packaging gap | P1 | Add changelog from wing onward; do not duplicate teamwork history verbatim |
| Manual / design doctrine | Rich | Indirect | Documentation gap | P2 | Add a concise wing architecture guide, not a second monolithic manual unless needed |
| E2E tests | Explicit | Not surfaced equivalently | Validation gap | P0 | Add `wing-e2e-test.sh` or equivalent smoke suite |
| Migration plan | Implicit in project files | Not packaged with wing | Transition gap | P0 | Write deletion criteria for teamwork |

---

## The Critical Distinction: Capability Gaps vs Packaging Gaps

Most remaining gaps are **not** core protocol gaps.

They are one of two things:

1. **proof gaps**
   Wing probably does the thing, but the repo does not yet prove it explicitly enough.

2. **packaging gaps**
   The runtime architecture is ready, but the surrounding docs, onboarding, and changelog surface have not yet caught up.

This matters because it changes the implementation strategy:

- do not rewrite major skill logic first
- first prove what already works
- then fill the missing packaging around it

---

## Highest-Risk Parity Areas

These are the areas where silent mismatch would hurt the most.

### 1. Claim flow correctness

Why it matters:

- this is where Issue intent becomes Contract intent
- if this is wrong, all later execution is wrong faster

What to verify:

- active mission handling
- worktree vs branch mode
- contract confirmation before branch creation
- roster membership logic
- stale/open issue recovery

### 2. Drive execution correctness

Why it matters:

- this is the cognitive heart of the protocol

What to verify:

- plan validation
- manual task classification
- freshness checks
- comment sync behavior
- branch/worktree safety
- sequential vs parallel mode boundary conditions

### 3. Ship and RC operational safety

Why it matters:

- these are the highest blast-radius paths

What to verify:

- unchecked task handling
- `Closes` vs `Refs` semantics
- side-effect triage
- review requested / changes requested flows
- RC uniqueness checks
- promote and hotfix correctness

---

## Recommended Completion Plan

### Phase 1 — Prove actual parity

Deliverables:

- a capability-by-capability parity matrix
- a `teamwork -> wing` command and behavior mapping table
- an explicit list of teamwork behaviors that are:
  - preserved
  - redesigned
  - intentionally dropped

Suggested artifact:

- `docs/research/wing-parity-matrix.md`

### Phase 2 — Add missing package surface

Deliverables inside `skills/wing/`:

- `README.md`
- `ONBOARDING.md`
- `CHANGELOG.md`

Rules:

- keep them outside hot-path `SKILL.md`
- explain the system at suite level
- do not dump old teamwork prose without adaptation

### Phase 3 — Add end-to-end validation

Deliverables:

- `wing-e2e-test.sh` or equivalent
- smoke suite for:
  - init
  - issue create
  - claim
  - drive
  - ship
  - rc prepare/promote

Design rule:

- test workflows, not just file presence

### Phase 4 — Formalize migration exit criteria

Deliverables:

- a written rule for when `skills/teamwork/` can be deleted
- a list of blockers that still keep teamwork alive

Suggested criteria:

1. wing passes e2e parity suite
2. wing docs package is present
3. downstream repo migration succeeds
4. no teamwork-only behavior remains required in production repos

---

## Anti-Baggage Rules

These rules should govern all parity work.

### Rule 1: No monolith regression

Do not move heavy operations back into `SKILL.md`.

### Rule 2: No runtime-doc stuffing

If the goal is onboarding or explanation, prefer `README.md`, `ONBOARDING.md`, or architecture docs.

### Rule 3: Keep one source of truth

VE, operational references, and migration rules should not be duplicated across all six skills.

### Rule 4: Preserve mechanism, not wording

If a teamwork behavior matters, preserve the behavior.
Do not preserve historical wording by default.

### Rule 5: Prefer redesign over compatibility fossilization

If a teamwork workaround existed because of the old flat-command architecture, do not copy it blindly into wing.

### Rule 6: Every parity item must be classified

For each teamwork feature:

- Preserve as-is
- Preserve but redesign
- Delete intentionally

No ambiguous middle state.

---

## Recommended Deliverables Backlog

### P0 — Must do before deleting Teamwork

1. Add wing parity matrix document
2. Add wing end-to-end smoke suite
3. Validate claim/drive/ship/rc parity on real workflow scenarios
4. Write migration exit criteria

### P1 — Must do for a complete Wing package

1. Add `skills/wing/README.md`
2. Add `skills/wing/ONBOARDING.md`
3. Add `skills/wing/CHANGELOG.md`
4. Add a short wing architecture note for contributors

### P2 — Good cleanup after parity

1. Trim transitional compatibility once no longer needed
2. Merge any duplicated operational guidance into shared references
3. Audit residual teamwork naming leakage inside wing comments/help text

---

## Concrete Next-Step Recommendation

If only one next step is taken, it should be this:

**Create a parity matrix and classify every teamwork capability as preserve / redesign / delete.**

Why this is first:

- it prevents accidental re-import of baggage
- it exposes true missing features vs missing docs
- it gives a stable scope boundary for all future wing completion work

After that, the best order is:

1. parity matrix
2. wing docs package
3. wing e2e validation
4. downstream migration
5. teamwork deletion

---

## Bottom-Line Recommendation

The path to "Wing with full Teamwork power but no historical baggage" is not:

- copy more teamwork text into wing

It is:

- preserve protocol semantics
- preserve deterministic scripts where still correct
- add missing proof and package surface
- explicitly reject monolith regressions

That is the clean replacement path.
