# Feedback Layer

`docs/feedback/` is the structured usage-feedback layer for `erwin`.

Its job is to capture what `Harness` actually feels like in use:

- repeated friction
- failed handoffs
- scheduler confusion
- council overuse or underuse
- authority drift
- success patterns worth preserving

This is not a dump of raw chats or every session log.

Raw traces stay in:

- [`_project/archive/sessions/`](../../_project/archive/sessions/)
- `.harness/` local runtime state
- issue / mission / run surfaces

`docs/feedback/` should hold normalized cases that are stable enough to teach the system something.

The canonical template for those cases lives in:

- [harness-usage-case-template.md](harness-usage-case-template.md)

## Case Shape

Each serious feedback case should answer:

- what happened
- where it happened
- why the current system allowed it
- which layer should have caught it
- whether it is a failure pattern or a success pattern
- whether it should update knowledge, proposals, evals, or nothing

## Status

This layer is now active as part of the repo knowledge system.

The next step is not to mass-import old sessions blindly.
It is to curate only the repeatable cases that add durable signal.

## Canonical Template

- [harness-usage-case-template.md](harness-usage-case-template.md) — durable schema and promotion rules for normalized Harness usage cases
