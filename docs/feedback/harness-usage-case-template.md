# Harness Usage Case Template

> Status: Active
>
> Date: 2026-04-11
>
> Purpose: provide a durable template for turning raw `Harness` usage into structured feedback instead of vague retrospectives or chat-memory fragments.

---

## When To Create A Case

Create a feedback case when at least one of these is true:

- the same friction happened more than once
- a real execution failure exposed a blind spot in `Harness`
- a control-plane behavior was surprisingly effective and worth preserving
- a session log contains a repeatable pattern that should not stay buried in narrative text

Do not create a case just because something happened once.

The goal is durable signal, not exhaustive journaling.

---

## Case Fields

Each case should be written with these fields:

- `case_id`
- `date`
- `title`
- `kind`
- `scope`
- `symptom`
- `what_happened`
- `where_seen`
- `trigger`
- `why_current_system_allowed_it`
- `layer_that_should_have_caught_it`
- `evidence`
- `repeatability`
- `impact`
- `recommended_next_step`
- `linked_knowledge`
- `linked_proposals`
- `status`

### Recommended `kind`

- `control-plane failure`
- `execution-plane failure`
- `coordination failure`
- `decision failure`
- `success pattern`

### Recommended `scope`

- `leader`
- `teammate`
- `council`
- `scheduler`
- `authority`
- `packaging`
- `docs`
- `automation`

### Recommended `repeatability`

- `one-off`
- `suspected-repeat`
- `confirmed-repeat`

### Recommended `status`

- `open`
- `distilled-into-knowledge`
- `encoded-as-eval`
- `closed`

---

## Minimal Example

```md
## HU-001 — Delegated follow-up went stale without clear recovery guidance

- `case_id`: `HU-001`
- `date`: `2026-04-01`
- `kind`: `control-plane failure`
- `scope`: `scheduler`
- `symptom`: leader could see a delegated run existed but not that it needed stale-follow-up recovery
- `what_happened`: a teammate run still owned an adopted follow-up packet after going stale, but earlier surfaces treated it as generic run staleness
- `where_seen`: scheduler + home + doctor mismatch during delegated follow-up lifecycle work
- `trigger`: stale delegated teammate run
- `why_current_system_allowed_it`: bounded delegated-follow-up drift had not yet been modeled as its own maintenance/diagnostic object
- `layer_that_should_have_caught_it`: `scheduler`, `doctor`, `authority`
- `evidence`: `_project/archive/sessions/...`, `scripts/harness-teammate-smoke-test.sh`
- `repeatability`: `confirmed-repeat`
- `impact`: leader guidance stayed too generic and recovery action was under-specified
- `recommended_next_step`: `encode-as-eval`
- `linked_knowledge`: `docs/knowledge/...`
- `linked_proposals`: `docs/harness-proposal-ledger.md`
- `status`: `encoded-as-eval`
```

---

## Promotion Rules

Feedback should usually move in one of three directions:

1. `stay in feedback`
   Use this when the case is real but not yet mature enough to generalize.
2. `distill into knowledge`
   Use this when the case connects to a broader pattern or topic thread.
3. `promote into proposal`
   Use this only when the case is already supported by strong knowledge or repeated evidence.

Default bias:

- repeated failure -> `knowledge` or `encode-as-eval`
- strong success pattern -> `knowledge`
- one-off confusion -> stay in `feedback` or close

---

## What Not To Do

- do not copy full chat logs into a case
- do not turn every session note into a case
- do not skip straight from raw friction to runtime change
- do not open proposals with no evidence trail back to real usage
