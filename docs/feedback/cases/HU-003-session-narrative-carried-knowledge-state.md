# HU-003 — Session narrative carried knowledge state

- `case_id`: `HU-003`
- `date`: `2026-04-02`
- `title`: `session narrative carried knowledge state`
- `kind`: `coordination failure`
- `scope`: `docs`
- `symptom`: fresh sessions could find the operating model, but still had to reconstruct which knowledge threads were active and what evidence had already converged by replaying session logs and recent commits
- `what_happened`: the repo gained the right layers and templates, but it still lacked a compact file-backed map of active threads, their upstream evidence, and their downstream proposal status
- `where_seen`: `docs/knowledge/README.md`, `docs/harness-session-handoff.md`, `_project/archive/sessions/2026-04-02-C.md`, `_project/archive/sessions/2026-04-02-D.md`
- `trigger`: knowledge continuity depended on reading narrative handoffs instead of opening one durable thread registry
- `why_current_system_allowed_it`: the knowledge layer defined what should exist, but did not yet maintain an active thread map tying `radar`, `feedback`, `knowledge`, and `proposals` together
- `layer_that_should_have_caught_it`: `knowledge`, `docs`
- `evidence`: `docs/harness-session-handoff.md`, `docs/knowledge/README.md`, `_project/archive/sessions/2026-04-02-C.md`
- `repeatability`: `confirmed-repeat`
- `impact`: knowledge continuity remained fragile, weekly syntheses had no canonical thread registry to update, and proposal review still depended too much on human recall
- `recommended_next_step`: `distill into knowledge`
- `linked_knowledge`: `docs/knowledge/active-threads.md`
- `linked_proposals`: `docs/proposals/harness-proposal-ledger.md`
- `status`: `distilled-into-knowledge`
