# Harness Learning Loop

> Status: Active
>
> Date: 2026-03-31
>
> Purpose: define how `Harness` learns from field signals, real usage, and internal gaps without collapsing into cargo-cult adoption or automatic runtime churn.

---

## Why This Exists

`Harness` should evolve from real pressure, not from hype and not from chat memory alone.

The repository now has three strong inputs:

- `field signals`
  - official engineering and product-development material from frontier AI companies
- `real usage cases`
  - Codex and Claude Code execution traces, workflow confusion, worktree/issue drift, handoff failures, and success patterns
- `internal harness gaps`
  - scheduler, council, authority, ownership, and runtime drift discovered inside `Harness` itself

Those inputs are valuable, but they should not directly mutate `Harness`.

The default rule is:

**signals become proposals first; only reviewed proposals can become evals, prototypes, or adopted system changes.**

---

## Core Rule

`Harness` uses a proposal-first learning loop.

That means:

- daily automation collects and filters signal
- weekly synthesis turns the strongest signal into structured proposals
- proposals enter a durable ledger
- only reviewed proposals can become:
  - eval cases
  - prototypes
  - accepted upgrades

This is not model training.
It is a disciplined repository-native learning loop.

---

## Input Sources

### 1. Field Signals

High-signal external material, prioritized toward official first-party engineering sources:

- OpenAI
- Anthropic
- Google

Supporting evidence only:

- Langfuse
- Phoenix
- Helicone
- OpenTelemetry

Field signals matter when they sharpen:

- harness engineering
- agent legibility
- eval practice
- observability
- orchestration
- context engineering
- tool use
- coding-agent runtime behavior

### 2. Real Usage Cases

`Harness` should continuously learn from actual agent work:

- Codex threads
- Claude Code threads
- worktree confusion
- issue/branch drift
- failed or weak handoffs
- over- or under-escalation into councils
- coordination failures between sessions
- successful patterns worth keeping

### 3. Internal Harness Gaps

These come from:

- runtime drift
- scheduler blind spots
- council/controller limitations
- authority adjudication gaps
- documentation drift
- weak smoke coverage

---

## Learning Pipeline

The learning loop has four stages.

### 1. Intake

Collect candidate signal without assuming it should change the system.

Expected intake surfaces:

- `AI Agent Daily Radar`
- `AI Agent Weekly Synthesis`
- explicit review sessions
- durable docs and audits

### 2. Proposal Extraction

Convert strong signals into proposal candidates.

Each candidate should answer:

- what changed or what happened
- why it matters for `Harness`
- which layer it would affect
- what the first safe next step should be

### 3. Proposal Review

Every proposal candidate must be classified before it influences implementation.

Allowed dispositions:

- `ignore`
- `watch`
- `encode-as-eval`
- `prototype`
- `adopt`

Default bias:

- external field signal -> `watch` or `encode-as-eval`
- repeated real usage failure -> `encode-as-eval` or `prototype`
- obvious internal design gap -> `prototype`
- only a small minority should jump to `adopt`

### 4. Backlog Feed

Only reviewed proposals feed future work.

Outputs may become:

- new smoke coverage
- new audit checks
- new scheduler/council/authority proposals
- doc updates
- explicit backlog candidates

No proposal should bypass this review path.

---

## Proposal Shape

Every serious proposal should be representable using these fields:

- `proposal_id`
- `date`
- `title`
- `origin`
- `source_ref`
- `theme`
- `target_layer`
- `recommended_disposition`
- `status`
- `linked_weekly_synthesis`
- `notes`

Recommended `target_layer` values:

- `protocol`
- `adapter`
- `profile`
- `scheduler`
- `authority`
- `council`
- `ownership`
- `eval`
- `docs`
- `automation`

Recommended `status` values:

- `watching`
- `queued-for-eval`
- `queued-for-prototype`
- `adopted`
- `rejected`
- `superseded`

The durable ledger for these entries lives in:

- [`docs/harness-proposal-ledger.md`](/Users/liyasong/Project/collov/erwin/docs/harness-proposal-ledger.md)

---

## Real Usage Case Taxonomy

Real usage should be structured before it enters the proposal loop.

Use these categories:

- `control-plane failure`
- `execution-plane failure`
- `coordination failure`
- `decision failure`
- `success pattern`

Each case should be summarized as:

- `what happened`
- `why current harness did not prevent it`
- `which harness layer should have caught it`
- `recommended disposition`

Mapping defaults:

- repeated failure -> `encode-as-eval`
- known design gap -> `prototype`
- one-off noise -> `watch` or `ignore`
- repeatable success pattern -> `watch` or `prototype`

Raw long-form thread logs should stay as reference material.
The repo should prefer structured cases and structured proposals.

---

## Automation Responsibilities

This learning loop deliberately reuses the existing automations instead of creating a third automation for direct system changes.

### `AI Agent Daily Radar`

Role:

- inbox-first intake
- no repo file writes
- identify strong `Harness`-relevant signals
- produce a small number of proposal candidates

Required sections:

- `New Signals`
- `Harness-Relevant Signals`
- `Candidate Proposals`
- `Watch Next`

Hard rule:

- daily radar must not recommend direct `Harness` implementation changes
- it may only recommend:
  - `watch`
  - `encode-as-eval`
  - `prototype`
  - `adopt`

### `AI Agent Weekly Synthesis`

Role:

- durable weekly artifact
- aggregate daily and field inputs
- distill only the strongest proposal candidates

Weekly output lives in:

- [`docs/radar/weekly/`](/Users/liyasong/Project/collov/erwin/docs/radar/weekly)

Required additional section:

- `Harness Evolution Proposals`

Each weekly proposal should include:

- `Title`
- `Source`
- `Observed Pattern`
- `Why It Matters For Harness`
- `Recommended Disposition`
- `What We Would Change If Adopted`
- `Why We Should Not Cargo-Cult This`

---

## Source Policy

### Priority Sources

- OpenAI
- Anthropic
- Google

### Supporting Sources

- Langfuse
- Phoenix
- Helicone
- OpenTelemetry

Supporting sources help interpret a first-party pattern.
They should not define the roadmap by themselves.

---

## What We Intend To Learn From Major Vendors

### OpenAI

Primary learning areas:

- harness engineering
- repository legibility
- agent-first execution environments
- eval discipline

### Anthropic

Primary learning areas:

- long-running harness design
- context engineering
- structured artifacts and resets
- evaluator design for long-running agents

### Google

Primary learning areas:

- build / scale / govern framing
- evaluation services
- observability and operator model

---

## Adoption Guardrails

External ideas are inputs, not authority.

Before any proposal becomes a prototype or adopted change, it must answer:

1. Is this solving a real problem we are actually facing?
2. Which layer does it change?
3. Should this first become an eval rather than an implementation?
4. Does it conflict with `Harness` core principles?

Current core principles:

- leader-led orchestration
- file-backed truth
- bounded automation
- proposal-first evolution
- no blind copying of external agent topology

---

## Relation To Existing Harness Docs

This document does not replace:

- [`docs/harness-architecture.md`](/Users/liyasong/Project/collov/erwin/docs/harness-architecture.md)
- [`docs/harness-operating-model.md`](/Users/liyasong/Project/collov/erwin/docs/harness-operating-model.md)
- [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)

Those documents define the system.

This document defines how the system should learn and evolve without losing discipline.

