# HU-004 — Productized surfaces outpaced deeper control-plane maturity

- `case_id`: `HU-004`
- `date`: `2026-04-02`
- `title`: `productized surfaces outpaced deeper control-plane maturity`
- `kind`: `control-plane failure`
- `scope`: `authority`
- `symptom`: the repo could look more complete than it really was because packaging, shell surfaces, and visible runtime shape matured faster than some of the hardest governance and control-plane commitments
- `what_happened`: repeated implementation-heavy sessions made `Harness` feel increasingly real through runtime packaging, shell entrypoints, and visible council/governance surfaces, but the gap and drift audits still showed partial maturity in ownership policy, harness-wide authority adjudication, and automated council orchestration
- `where_seen`: `docs/research/harness-gap-audit-2026-03-29.md`, `docs/research/harness-research-drift-audit-2026-03-29.md`, `_project/CURRENT.md`
- `trigger`: implementation work naturally favored surfaces that were easier to package, smoke-test, and demonstrate than the hardest control-plane semantics
- `why_current_system_allowed_it`: the repo had no normalized feedback case or durable knowledge note tying external harness-engineering ideals to the specific internal maturity gaps that still remained
- `layer_that_should_have_caught_it`: `knowledge`, `authority`, `council`, `docs`
- `evidence`: `docs/research/harness-gap-audit-2026-03-29.md`, `docs/research/harness-research-drift-audit-2026-03-29.md`, `docs/radar/weekly/2026-03-30.md`
- `repeatability`: `confirmed-repeat`
- `impact`: future sessions could mistake visible productization for complete governance maturity, which would bias work toward adding more surfaces instead of closing the remaining control-plane gaps
- `recommended_next_step`: `distill into knowledge`
- `linked_knowledge`: `docs/knowledge/harness-engineering-vs-internal-gaps-2026-04-02.md`, `docs/knowledge/active-threads.md`
- `linked_proposals`: `docs/proposals/harness-proposal-ledger.md`
- `status`: `distilled-into-knowledge`
