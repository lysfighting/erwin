# HU-006 — Active council needed current work visibility

- `case_id`: `HU-006`
- `date`: `2026-04-02`
- `title`: `active council needed current work visibility`
- `kind`: `control-plane failure`
- `scope`: `council`
- `symptom`: leader-facing active-council surfaces were too abstract when they exposed only generic `next_requirement` or batch metadata instead of the current live packet and queued work
- `what_happened`: repeated council-execution work had to add current work packets, concrete current-wave action packets, and blocked-batch summaries to scheduler active-council packets and `meetings summary` before leaders could see what a live council was actually doing
- `where_seen`: `_project/archive/sessions/2026-03-31-F.md`, `_project/archive/sessions/2026-03-31-G.md`
- `trigger`: multi-phase council execution moved from independent review into debate waves, but leader-facing surfaces still described progress too generically
- `why_current_system_allowed_it`: the runtime had durable meeting state, but leader-facing control-plane views did not yet project the right bounded current-work object for live councils
- `layer_that_should_have_caught_it`: `council`, `scheduler`, `docs`
- `evidence`: `_project/archive/sessions/2026-03-31-F.md`, `_project/archive/sessions/2026-03-31-G.md`
- `repeatability`: `confirmed-repeat`
- `impact`: leaders could not inspect live council state precisely enough to understand the current ask, queued batches, or whether the council was progressing through bounded work versus generic meeting motion
- `recommended_next_step`: `distill into knowledge`
- `linked_knowledge`: `docs/knowledge/runtime-legibility-patterns-2026-04-02.md`, `docs/knowledge/harness-engineering-vs-internal-gaps-2026-04-02.md`, `docs/knowledge/active-threads.md`
- `linked_proposals`: `docs/proposals/harness-proposal-ledger.md`, `docs/proposals/harness-eval-candidate-ledger.md`
- `status`: `distilled-into-knowledge`
