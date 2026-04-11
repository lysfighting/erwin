# Knowledge System Operating Model

> Status: Active
>
> Date: 2026-04-11
>
> Purpose: define how `erwin` turns external signal and real Harness usage into durable knowledge and selective proposals without collapsing back into `docs/research/` sprawl.

---

## Why This Exists

The repo now has a clearer documentation shape, but a directory split alone is not enough.

Without an operating model, the system will drift back into:

- `docs/research/` as the default dumping ground
- weekly notes that overproduce proposals
- raw sessions treated as if they were durable knowledge
- proposals generated from fresh excitement instead of stable understanding
- automation output that exists but does not have a stable human-facing reading surface

This document defines the active model.

---

## System Layers

### 1. Radar

Directory:

- [`docs/radar/`](../radar/)

Job:

- collect high-signal external changes
- prioritize first-party engineering, docs, release-note, and system-card evidence
- add paper and benchmark intake where method questions matter
- surface what deserves attention, not decide what the repo should adopt

Primary outputs:

- daily notes
- weekly synthesis
- watch-next guidance
- the visual briefing surface at [`docs/index.html`](../index.html)

### 2. Feedback

Directory:

- [`docs/feedback/`](../feedback/)

Job:

- normalize real `Harness` usage pressure
- record repeated failure patterns and durable success patterns
- separate true system friction from one-off session noise

Primary outputs:

- structured usage cases
- failure/success pattern summaries
- normalized input for the knowledge layer
- template-driven case capture using [`docs/feedback/harness-usage-case-template.md`](../feedback/harness-usage-case-template.md)

### 3. Knowledge

Directory:

- [`docs/knowledge/`](./)

Job:

- connect radar signal with feedback and internal gaps
- build durable topic threads, comparisons, timelines, and concept maps
- distinguish vendor framing, paper evidence, and runtime evidence

Primary outputs:

- synthesis notes
- topic threads
- knowledge maps
- durable cross-source interpretations

### 4. Proposals

Directory:

- [`docs/proposals/`](../proposals/)

Job:

- hold reviewed change candidates only after upstream evidence is mature enough
- keep proposals selective, bounded, and reviewable

Primary outputs:

- reviewed proposal ledgers
- proposal dispositions such as `watch`, `encode-as-eval`, `prototype`, and `adopt`

### 5. Research

Directory:

- [`docs/research/`](../research/)

Job:

- preserve one-time incubation analysis, architectural deep dives, and audits
- remain useful reference material without pretending to be the active knowledge layer

Primary outputs:

- historical design inputs
- bounded investigations
- audits and drift checks

---

## Default Flow

The default path is:

`radar -> feedback -> knowledge -> proposals -> implementation`

Important constraints:

- `radar` is not direct permission to change `Harness`
- `feedback` is not just dumping session logs into docs
- `knowledge` is not the same as “more research notes”
- `proposals` should be rarer than weekly notes
- implementation should usually consume reviewed proposals or already-decided internal gaps

---

## Reading Surface

The human-facing front door for this system should not depend on transient automation inbox cards.

Expected reading surfaces:

- [`docs/index.html`](../index.html) for the current visual briefing
- `latest.md` files for each recurring stream as the durable raw source
- canonical layer READMEs for how the system is organized

The reading surface is allowed to be generated.
The durable source of truth remains the markdown notes and canonical docs.

---

## Knowledge Agent / Skill Role

`erwin` should treat knowledge work as its own active operating surface.

That means a dedicated knowledge skill or knowledge-oriented agent should:

- maintain topic continuity across sessions
- keep external and internal evidence connected
- decide when a fresh signal deserves a durable knowledge update
- decide when a knowledge thread is mature enough to emit a proposal
- resist proposal spam and research sprawl
- keep the reading surface legible instead of letting good signal disappear into raw markdown or archived threads

It should not:

- mutate `Harness` runtime directly from new signal
- treat every paper or product launch as implementation pressure
- collapse back into generic “research assistant” behavior

The knowledge role is closer to:

- curator
- synthesizer
- pattern keeper
- proposal gatekeeper
- information architect

than to a generic search bot.

---

## Reliability Rules

A proposal is more trustworthy when it is grounded in at least two of these three evidence classes:

1. external signal
2. real usage feedback
3. internal runtime gap

Single-source proposals are still allowed, but the default disposition should be more conservative.

Recommended bias:

- external-only -> `watch` or `encode-as-eval`
- feedback-only -> `encode-as-eval` or `prototype`
- internal-gap-only -> `prototype`
- multi-source convergence -> stronger candidate for `adopt`

---

## What To Prune Aggressively

The new model should actively resist:

- duplicate summaries
- stale one-off notes that never became durable knowledge
- proposal records with no clear evidence trail
- weekly artifacts pretending to be long-term knowledge
- inbox metadata or archived threads being treated as the only visible reading surface
- using `docs/research/` as the default landing zone for everything unfamiliar

The goal is not to keep more text.
The goal is to keep a more legible knowledge system.
