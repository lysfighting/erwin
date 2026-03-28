# Erwin Harness Incubation Blueprint

> Date: 2026-03-28
>
> Purpose: turn the current `teamwork` / `wing` / Codex migration confusion into a concrete incubation strategy for Erwin's own harness engineering direction.

---

## Core Judgment

Yes, this migration can and should incubate Erwin's own harness engineering.

But the thing to incubate is not:

- a new monolithic agent runtime
- a renamed `teamwork`
- a premature "Codex version of `wing`"

The thing to incubate is:

**a portable harness engineering stack**

That stack should be made of:

- protocol core
- runtime adapters
- deterministic mechanisms
- repo-legibility conventions
- project profiles
- eval and validation loops

This is the cleanest interpretation of the work you have already done.

---

## Why The Current Confusion Exists

Three historically different jobs have been collapsing into one name.

### Job 1: Protocol design

This is the true durable asset:

- issue framing
- contract generation
- execution discipline
- PR / release lifecycle
- verification and freshness rules

### Job 2: Runtime adaptation

This is the machinery that made the protocol usable in a specific environment:

- Claude commands in `teamwork`
- Anthropic skills in `wing`
- next: Codex skills / Codex-native workflow

### Job 3: Product / conceptual branding

This is where names like `teamwork` and `wing` start carrying meaning beyond implementation.

The confusion comes from all three jobs being fused together.

---

## The Best Reframe

The most useful reframe is:

- `teamwork` = legacy product/package around the protocol
- `wing` = an experiment in architectural refactoring plus a possible future orchestration identity
- Codex migration = the first serious test of protocol portability

Under this framing, the question is no longer:

> "Should we port teamwork or wing?"

It becomes:

> "What parts belong to the portable protocol, what parts belong to the runtime adapter, and what future concept deserves the `wing` name?"

That is a much cleaner design problem.

---

## Proposed Stack

Erwin should explicitly define a six-layer stack.

### Layer 1: Protocol Core

This is runtime-neutral.

Contains:

- lifecycle semantics
- invariants
- contract schema
- state model
- release semantics
- cognitive rules

Current artifact:

- [`protocol-core-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/protocol-core-spec-2026-03-28.md)

### Layer 2: Deterministic Mechanisms

This is still portable, but operational.

Contains:

- config resolution
- branch / worktree operations
- contract freshness and mutation
- label state machine
- PR helpers
- RC helpers

This is where a large amount of `teamwork`'s real value lives.

### Layer 3: Runtime Adapters

This is where vendor- or environment-specific integration belongs.

Examples:

- Anthropic adapter
- Codex adapter
- later: Copilot / OpenHands adapter

Current artifact:

- [`codex-runtime-adapter-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/codex-runtime-adapter-spec-2026-03-28.md)

### Layer 4: Project Profiles

This isolates repo-specific reality from the protocol and adapter.

Examples:

- `erwin`
- `vi_agent`
- generic starter

Contains:

- branch names
- release topology
- verification commands
- staging workflow ids
- label conventions
- project docs entry points

### Layer 5: Repo Legibility System

This is the part most teams underestimate.

Contains:

- project status files
- index docs
- issue templates
- contract templates
- eval expectations
- agent-readable docs layout

This is where harness engineering becomes more than "better prompts."

### Layer 6: Eval / Validation

Contains:

- static lint
- smoke tests
- parity checks
- scenario-based behavior validation
- migration confidence reports

Without this layer, the rest turns back into intuition.

---

## Where Wing Should Sit

Right now, `wing` should not be forced to equal the Codex migration.

The healthiest temporary position is:

- `wing` is an architectural and conceptual branch, not yet the required name of the Codex adapter

That leaves room for three possible futures:

### Option A: Wing becomes the umbrella product

Then `teamwork` and future adapters become implementation history, while `wing` becomes the broad orchestrator/harness concept.

### Option B: Wing stays Anthropic-lineage specific

Then `wing` refers to the Anthropic skill-family rewrite, and Codex gets a sibling adapter under a more neutral protocol umbrella.

### Option C: Wing becomes the future orchestration layer above adapters

Then:

- protocol core sits below
- adapters sit in the middle
- `wing` becomes the broader orchestration/product layer above them

Of the three, Option C currently feels the most aligned with your recorded intent.

---

## Recommendation On Naming

For now:

- do not collapse the Codex adapter into `wing`
- do not resurrect `teamwork` as the future identity
- do not force a final name before the architecture is clean

Use neutral names for architecture work:

- protocol core
- Codex adapter
- Anthropic adapter
- project profile
- harness stack

Then let product naming happen after the boundaries are real.

This avoids repeating the exact confusion that produced the current knot.

---

## Market Positioning Insight

The industry trend is moving toward "harness engineering" as a serious layer of engineering work.

The strongest current references are not just "better agents," but:

- OpenAI / Codex emphasizing harness engineering and repo legibility
- Anthropic emphasizing skill modularity and progressive disclosure
- GitHub emphasizing issue/PR-native coding agents with policy controls
- OpenHands / Goose / OpenCode showing open, extensible runtime surfaces

Erwin's best differentiated position is not to copy any one of these.

It is to specialize in:

**portable protocol + adapter architecture + repo-legibility methodology**

That is both more original and more aligned with your body of work.

---

## What Erwin Could Become

If incubated correctly, Erwin could become:

### 1. A protocol foundry

Where durable AI software-work protocols are designed, tested, and versioned.

### 2. An adapter workshop

Where those protocols are made to work across Codex, Anthropic, and later runtimes.

### 3. A repo-legibility playbook

Where teams learn how to make repositories truly agent-operable.

### 4. A validation lab

Where skills, workflows, and harnesses are linted, smoke-tested, and compared for parity.

This is much more interesting than simply maintaining a historical `teamwork` port.

---

## What To Build Next

### P0: Lock the architecture

1. Keep refining the protocol core spec
2. Keep the Codex adapter spec neutral
3. Use the parity matrix to separate preserve / redesign / delete

### P1: Create the portable substrate

1. Extract deterministic mechanisms that should outlive any one adapter
2. Define a stable contract schema
3. Define project-profile schema

### P2: Build the first clean adapter

1. Implement Codex home/doctor
2. Implement Codex claim
3. Implement Codex drive
4. Reuse or rewrite ship / rc only after claim-drive is sound

### P3: Add validation

1. static lint for adapter surfaces
2. smoke scenarios for issue-claim-drive-ship
3. release scenario tests for rc flows
4. parity report against `teamwork`

### P4: Re-decide branding

Only after the above should you decide whether:

- Codex adapter is called `wing`
- `wing` becomes the umbrella orchestration layer
- or a new top-level name is needed

---

## Immediate Product Thesis

The most credible thesis for Erwin right now is:

> Erwin is evolving from a skill incubator into a harness engineering foundry: a place to design portable autonomous software-work protocols, make repositories agent-legible, and ship runtime adapters that bring those protocols to life across agent environments.

That thesis is broad enough to include your current work.
It is also narrow enough to stay actionable.

---

## Bottom Line

This migration should absolutely be used to incubate your own harness engineering.

But the win is not "we ported `teamwork` again."

The win is:

- the protocol is now explicit
- the adapter boundary is now explicit
- `wing` no longer has to carry every unresolved architectural question
- and Erwin gains a real, differentiated direction beyond one runtime's skill format
