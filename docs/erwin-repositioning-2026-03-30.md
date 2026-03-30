# erwin Repositioning Proposal (2026-03-30)

## Why this exists

`erwin` no longer behaves like a narrow "skill incubator with a showcase".

The live center of gravity has already shifted:

- the most active and recent work is around `Harness`
- the repository contains growing AI/agent research and operating-model material
- the old mental model still points attention toward `showcase`, `teamwork`, and broad `_project/` process machinery that no longer reflects the highest-value work

This note proposes a practical repositioning without forcing an immediate destructive cleanup.

## Recommendation

Keep using `erwin`.

Do **not** create a separate repository for the knowledge-radar automation yet.

Instead, redefine `erwin` as:

> an exploration lab for advanced AI/agent operating patterns, with skills, harnesses, research notes, and intelligence workflows living in one place

Under that framing:

- `Harness` is a major current track, not the whole identity
- the automation belongs inside `erwin`
- the automation should track the broader AI/agent field, with topic clusters such as `Harness`, `evals`, `observability`, `coding agents`, `tool use`, and `orchestration`

## Why not split into a new project yet

Creating a new repo now would solve naming discomfort, but it would also introduce avoidable fragmentation:

- knowledge would split across two repos before the new operating model is stable
- the automation would lose proximity to the code, docs, and experiments it should inform
- cleanup pressure inside `erwin` would be postponed rather than resolved

Right now the problem is not "wrong repo". The problem is "right repo, stale shape".

## What the automation should be

The automation should be a **field radar**, not a `Harness`-only watcher.

### Primary scope

Track high-signal updates from authoritative publishers in the AI/agent space, especially:

- Anthropic
- OpenAI
- Google
- Moonshot
- MiniMax
- Qwen
- DeepSeek

### Secondary scope

Track carefully selected open-source and ecosystem references only when they help interpret the field:

- Langfuse
- Phoenix
- Helicone
- OpenTelemetry
- other tools only when they clarify a pattern already visible in first-party sources

### Current priority themes

The radar should always maintain a small set of active themes. For now:

- `Harness`
- long-running agents
- agent evals
- tracing / observability
- context engineering
- coding agents

`Harness` is a theme, not the boundary.

## Recommended repo identity

Proposed new plain-language description:

> `erwin` is a working lab for advanced AI/agent systems. It collects frontier knowledge, incubates reusable skills, and prototypes operating models such as Harness.

This keeps the original strength of the repo, but updates the center of gravity to match reality.

## What to keep

These areas still fit the new mission:

- `skills/harness/`
- `scripts/`, `cmd/`, `internal/`, `harness/`
- selected reusable skills such as `visual-test`, `arch-review`, `session-record`, `skill-tester`
- selected research and architecture docs that directly support current agent-system work
- a slimmed `_project/` layer, but only for live decisions and current direction

## What to demote from the center

These are not necessarily wrong, but they should stop defining the repo:

- `showcase/`
- `teamwork-demo/`
- legacy `teamwork` positioning
- old "meta-skill factory" framing when it obscures current agent-systems work
- session-by-session process residue that no longer informs present decisions

## What to archive or trim first

Priority order:

1. Archive stale docs from the top-level active index.
2. Collapse `_project/` to a lean operating layer.
3. Move historical product experiments out of the repo center.

### 1. Docs

The docs index should make the current mission obvious in under 30 seconds.

Likely archive candidates:

- teamwork architecture lineage that survives only as history
- showcase architecture plans that no longer drive active work
- old process docs whose implementation target has disappeared

Keep active only:

- current harness docs
- current operating-model docs
- current radar / research synthesis docs
- a small number of skill docs that still matter

### 2. `_project/`

`_project/` currently mixes live control documents with a large amount of historical session exhaust.

Recommended target shape:

- `_project/CURRENT.md`
- `_project/GOALS.md`
- `_project/decisions/`
- `_project/issues/` only if still actively used
- `_project/sessions/` moved behind an archive boundary or aggressively summarized

The biggest problem is not existence, but default visibility. Old logs should not be part of the everyday surface area.

### 3. Product/demo artifacts

`showcase/` and `teamwork-demo/` are expensive in cognitive weight because they suggest an old product direction.

Recommended path:

- stop treating them as top-level narrative centers
- either move them into an `archive/` or `legacy/` area later
- or extract anything still valuable and then remove them

Do not delete them before harvesting useful assets.

## Suggested operating model after slimming

The repo should read as four clear lanes:

1. `Radar`
   External AI/agent intelligence, synthesis, and recurring automation outputs.
2. `Harness`
   The active orchestration operating-model prototype and implementation.
3. `Skills`
   Reusable skills that remain strategically relevant.
4. `Archive`
   Older demos, superseded docs, and historical logs.

## Near-term actions

Recommended next steps, in order:

1. Update the automation so it tracks the broader AI/agent field, with `Harness` as one active topic.
2. Reword repo-facing docs so `erwin` is described as an AI/agent exploration lab.
3. Triage docs into `active` vs `historical`.
4. Define a concrete archive boundary for `showcase/`, `teamwork-demo/`, and stale `_project/sessions/`.
5. Only after that, decide whether any part deserves extraction into a new repo.

## Decision

For now:

- keep the automation in `erwin`
- keep `erwin` as the umbrella project
- pivot the repo identity explicitly
- slim aggressively before splitting

That is the lowest-friction path with the highest strategic coherence.
