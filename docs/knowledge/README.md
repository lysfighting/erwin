# Knowledge Layer

`docs/knowledge/` is the durable understanding layer for `erwin`.

It sits between raw intake and action proposals.

The canonical operating model for this layer lives in:

- [operating-model.md](operating-model.md)
- [briefing.md](briefing.md)
- [index.json](index.json)

Its job is to weave together:

- `radar` intake
- `feedback` from real Harness usage
- papers and benchmarks
- internal runtime gaps

This is where the repo should maintain:

- theme syntheses
- concept maps
- timelines
- source reconciliation
- vendor-claim vs paper-evidence comparisons
- repeated patterns that are strong enough to outlive one session

## What Belongs Here

- multi-source synthesis notes
- topic threads that will keep evolving
- notes that connect external evidence with internal Harness pressure
- papers or benchmark writeups that matter beyond one narrow task

## What Does Not Belong Here

- raw daily intake
- one-off incubation specs
- every weekly note
- implementation proposals that have not been reviewed yet

## Current Notes

- [operating-model.md](operating-model.md) — canonical role definition for the repo knowledge system, including the knowledge skill / knowledge-agent responsibility split
- [active-threads.md](active-threads.md) — compact registry of live knowledge threads, evidence classes, proposal linkage, and update triggers
- [briefing.md](briefing.md) — generated text-first reading surface built from the latest radar streams
- [index.json](index.json) — generated machine-readable index for agents, scripts, and future surfaces
- [coding-agent-self-iteration-and-paper-intake-2026-04-02.md](coding-agent-self-iteration-and-paper-intake-2026-04-02.md) — durable synthesis of first-party and paper-aware coding-agent self-iteration signal
- [harness-engineering-vs-internal-gaps-2026-04-02.md](harness-engineering-vs-internal-gaps-2026-04-02.md) — comparison of external harness-engineering signal against internal repo gaps
- [runtime-legibility-patterns-2026-04-02.md](runtime-legibility-patterns-2026-04-02.md) — reusable runtime legibility patterns for Harness and agent systems
- [benchmark-integrity-and-eval-awareness-2026-04-09.md](benchmark-integrity-and-eval-awareness-2026-04-09.md) — benchmark-trust and eval-awareness thread for interpreting external agent scores

## Generated Surfaces

The knowledge layer now exposes three coordinated reading surfaces:

- [../index.html](../index.html) — visual front door for humans
- [briefing.md](briefing.md) — fast text briefing
- [index.json](index.json) — machine-readable state snapshot

`index.json` now includes a lightweight typed knowledge network under `knowledge_network`:

- `artifacts`
  - file-backed nodes across `Radar`, `Feedback`, `Knowledge`, `Proposals`, and selected `Implementation` entrypoints
- `links`
  - typed local references extracted from Markdown links, such as evidence citation, proposal pressure, implementation targets, and ordinary references
- `threads`
  - parsed active knowledge threads from [active-threads.md](active-threads.md), including evidence classes, linked proposals, and anchor artifacts
- `proposal_records`
  - parsed rows from [../proposals/harness-proposal-ledger.md](../proposals/harness-proposal-ledger.md)

The network is intentionally generated.
Do not edit `index.json` directly; update the Markdown source or implementation anchor, then rebuild the surface.

Rebuild them with:

```bash
python3 scripts/build-knowledge-surface.py
```

The generator reads the current `docs/radar/*/latest.md` files and rewrites all three outputs together so the visual, textual, and machine-facing surfaces stay aligned.

## Automation Linkage

Recurring radar automations should refresh this surface after they update their own stream files.

Expected sequence:

1. write the dated radar note
2. refresh the stream `latest.md`
3. run `python3 scripts/build-knowledge-surface.py`
4. verify [../index.html](../index.html), [briefing.md](briefing.md), and [index.json](index.json) all exist

This keeps the repo-owned reading surface independent from transient Codex inbox cards or archived automation threads.
