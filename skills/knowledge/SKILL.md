# Knowledge

Use this skill when the repo needs durable knowledge work rather than direct runtime implementation.

This skill treats `erwin` knowledge as a system with explicit layers:

- `docs/radar/` for external signal intake
- `docs/feedback/` for structured Harness usage pressure
- `docs/knowledge/` for durable synthesis
- `docs/proposals/` for reviewed change candidates
- `docs/research/` for one-off incubation analysis and audits

## Use This Skill For

- curating new external signals into the right knowledge layer
- normalizing repeated Harness usage friction into structured feedback
- writing or updating durable synthesis notes
- deciding whether a topic is mature enough for a proposal
- tightening repo knowledge architecture when docs start drifting back into sprawl

## Do Not Use This Skill For

- direct `Harness` runtime changes from fresh signal
- treating every daily or weekly note as a proposal
- dumping long raw research into `docs/research/` by default
- replacing canonical `Harness` docs with vendor framing or paper claims

## Working Rules

1. Start with the docs entrypoints:
   - `docs/README.md`
   - `docs/feedback/harness-usage-case-template.md`
   - `docs/knowledge/README.md`
   - `docs/knowledge/operating-model.md`
2. If the topic touches Harness evolution, also read:
   - `docs/harness-learning-loop.md`
   - `docs/harness-proposal-ledger.md`
3. Put new material in the narrowest correct layer.
4. Prefer durable synthesis over endless intake accumulation.
5. Create proposals selectively and only when the knowledge thread is mature enough.
6. Keep the human-facing reading surface current when the layer structure changes materially.

## Typical Outputs

- update `docs/radar/` when the repo needs stronger intake policy
- add or refine `docs/feedback/` cases when real usage patterns repeat
- add or refine `docs/knowledge/` synthesis notes when a theme becomes durable
- add or update `docs/harness-proposal-ledger.md` entries when evidence justifies a reviewed change candidate
- demote or prune stale material that no longer belongs in the active knowledge system
