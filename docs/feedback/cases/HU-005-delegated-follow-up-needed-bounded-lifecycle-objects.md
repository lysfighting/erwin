# HU-005 — Delegated follow-up needed bounded lifecycle objects

- `case_id`: `HU-005`
- `date`: `2026-04-02`
- `title`: `delegated follow-up needed bounded lifecycle objects`
- `kind`: `control-plane failure`
- `scope`: `scheduler`
- `symptom`: adopted follow-up work delegated to a teammate kept resurfacing as stale, cancelled, completed, or reopened state that generic run-level wording could not explain precisely enough
- `what_happened`: across multiple runtime slices, delegated adopted follow-up work repeatedly required more specific scheduler, `home`, `doctor`, and `authority-audit` behavior before the system could distinguish stale delegated ownership, cancelled realignment, completed quieting, and reopen cycles without drift or misleading top-level guidance
- `where_seen`: `_project/archive/sessions/2026-03-31-AD.md`, `_project/archive/sessions/2026-03-31-AE.md`, `_project/archive/sessions/2026-04-01-B.md`, `_project/archive/sessions/2026-04-01-F.md`, `_project/archive/sessions/2026-04-01-H.md`, `_project/archive/sessions/2026-04-01-I.md`
- `trigger`: teammate-owned adopted follow-up packets moved through stale, cancelled, completed, and reopened states while governing leader surfaces still needed exact top-level guidance
- `why_current_system_allowed_it`: the runtime originally treated delegated follow-up mostly as generic run staleness or generic post-council residue instead of as a bounded control-plane object with its own lifecycle semantics
- `layer_that_should_have_caught_it`: `scheduler`, `authority`, `doctor`
- `evidence`: `_project/archive/sessions/2026-03-31-AD.md`, `_project/archive/sessions/2026-03-31-AE.md`, `_project/archive/sessions/2026-04-01-B.md`, `_project/archive/sessions/2026-04-01-F.md`, `_project/archive/sessions/2026-04-01-H.md`, `_project/archive/sessions/2026-04-01-I.md`
- `repeatability`: `confirmed-repeat`
- `impact`: leaders could get guidance that was too generic, too stale, or pointed at the wrong recovery action; delegated follow-up work needed multiple fixes before it became legible and governable
- `recommended_next_step`: `distill into knowledge`
- `linked_knowledge`: `docs/knowledge/runtime-legibility-patterns-2026-04-02.md`, `docs/knowledge/harness-engineering-vs-internal-gaps-2026-04-02.md`, `docs/knowledge/active-threads.md`
- `linked_proposals`: `docs/proposals/harness-proposal-ledger.md`, `docs/proposals/harness-eval-candidate-ledger.md`
- `status`: `distilled-into-knowledge`
