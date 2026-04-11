# Proposal Layer

`docs/proposals/` is the action layer for repo change candidates.

Proposals belong here only after they have passed through stronger upstream shaping:

- `radar` gathered the outside signal
- `feedback` captured real usage pressure
- `knowledge` turned that pressure into durable understanding

Only then should a proposal become a durable control-plane record.

## Current Canonical Proposal Surface

The reviewed proposal ledger for `Harness` still lives at:

- [`../harness-proposal-ledger.md`](../harness-proposal-ledger.md)

This directory exists so proposal work has an explicit architectural home, even while the canonical ledger remains a top-level Harness document.

## Working Rule

Do not treat every daily or weekly run as proposal generation.

Proposal creation should be selective.
The quality bar is:

- repeated signal
- enough evidence density
- clear relevance to `Harness` or repo direction
- a bounded recommended disposition such as `watch`, `encode-as-eval`, `prototype`, or `adopt`
