# HU-007 — Benchmark results needed eval-trust interpretation

- `case_id`: `HU-007`
- `date`: `2026-04-09`
- `title`: `benchmark results needed eval-trust interpretation`
- `kind`: `decision failure`
- `scope`: `eval`
- `symptom`: external benchmark results were strong enough to influence knowledge and proposal posture, but the repo still lacked a durable habit of separating capability signal from benchmark-trust risk
- `what_happened`: once the weekly synthesis started incorporating long-horizon coding benchmarks and Anthropic's eval-awareness writeup, it became clear that benchmark-heavy evidence can sharpen the knowledge system while also carrying its own trust problem
- `where_seen`: `docs/radar/weekly/2026-04-03.md`, `docs/radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md`
- `trigger`: benchmark-rich external signal entered the knowledge system and pushed directly on `HLP-003` and `HLP-005`
- `why_current_system_allowed_it`: the repo had a paper-watch archive and review gates, but no dedicated knowledge note or eval candidate yet for benchmark integrity and eval-awareness risk
- `layer_that_should_have_caught_it`: `knowledge`, `eval`, `proposals`
- `evidence`: `docs/radar/weekly/2026-04-03.md`, `docs/radar/papers/2026-04-03-long-horizon-coding-agent-benchmarks.md`, `docs/knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md`
- `repeatability`: `suspected-repeat`
- `impact`: benchmark-driven proposal pressure could become sloppier than the rest of the knowledge system, confusing “important result” with “fully trustworthy measure”
- `recommended_next_step`: `distill into knowledge`
- `linked_knowledge`: `docs/knowledge/benchmark-integrity-and-eval-awareness-2026-04-09.md`, `docs/knowledge/active-threads.md`
- `linked_proposals`: `docs/proposals/harness-eval-candidate-ledger.md`, `docs/proposals/harness-proposal-ledger.md`
- `status`: `distilled-into-knowledge`
