# Should Erwin Incubate Its Own Harness Engineering?

> Date: 2026-03-28
>
> Question: can this migration moment be used to incubate our own harness engineering, and what should we learn from the market instead of reinventing blindly?

---

## Short Answer

Yes, but only if we are disciplined about **what layer we are incubating**.

If we try to build "our own harness" as a full competing runtime right now, that is probably premature.

If we use this transition to incubate:

- a protocol core
- runtime adapters
- repo-legibility conventions
- validation/eval patterns
- reusable skills and project profiles

then yes, this is exactly the right moment.

---

## The Main Strategic Choice

There are two very different ambitions hiding under the phrase "our own harness engineering."

### Option A — Build a full new harness/runtime

This means something like:

- own CLI or app shell
- own sandbox model
- own tool/plugin execution layer
- own session/memory/orchestration surfaces

This is very expensive.

### Option B — Build a harness engineering methodology + reusable harness assets

This means:

- protocol specs
- runtime adapters
- agent-legible repo conventions
- docs layouts
- eval and lint strategies
- background automation patterns
- project-profile templates

This is much more aligned with erwin's stated mission as a meta-skill factory.

My recommendation is:

**Choose Option B now.**

---

## Why This Fits Erwin Specifically

`erwin` is not supposed to become just another coding-agent shell.

Per [`_project/GOALS.md`](/Users/liyasong/Project/collov/erwin/_project/GOALS.md):

- erwin is a meta-skill factory
- erwin is an educational platform for AI agent orchestration
- erwin should help users build their own skills

That maps much more naturally to:

- protocol design
- harness patterns
- reusable scaffolding

than to building a new monolithic runtime product.

So yes, this can absolutely incubate your own harness engineering.
But the right artifact is probably:

**an opinionated harness methodology and adapter stack**

not:

**yet another closed environment**

---

## Industry References Worth Studying

These are the most useful reference classes right now.

### 1. OpenAI / Codex — harness engineering as first-class engineering work

Why it matters:

- explicit articulation of harness engineering
- repository knowledge as system of record
- strong emphasis on agent legibility
- custom linters, execution plans, observability, worktrees, and background cleanup agents

What to learn:

- the harness is built from docs, tools, evals, repo structure, and feedback loops
- "one giant AGENTS.md" is the wrong direction
- custom mechanized constraints are leverage multipliers

Reference:

- [Harness engineering: leveraging Codex in an agent-first world](https://openai.com/index/harness-engineering/)

### 2. Anthropic Skills — progressive disclosure and modular skill packaging

Why it matters:

- this is the exact design pressure that created `wing`

What to learn:

- focused skill entry points
- keep heavy detail in supporting files
- packaging matters for model reliability

References:

- [How to create custom Skills](https://support.claude.com/en/articles/12512198-how-to-create-custom-skills)
- [Anthropic skills repository](https://github.com/anthropics/skills)

### 3. GitHub Copilot coding agent — issue/PR-native async agent workflows

Why it matters:

- it shows a strong GitHub-native async model
- custom agents + MCP + branch restrictions + policy control are all explicit

What to learn:

- governance model
- issue/PR-native async delegation
- sandbox and permission boundaries
- repo knowledge enhancement via memory/custom instructions

References:

- [About GitHub Copilot coding agent](https://docs.github.com/copilot/concepts/agents/coding-agent/about-coding-agent)
- [Creating custom agents for Copilot coding agent](https://docs.github.com/en/copilot/how-tos/use-copilot-agents/coding-agent/create-custom-agents)

### 4. OpenHands — open platform mindset

Why it matters:

- illustrates the "open coding-agent platform" direction

What to learn:

- how far a platform/runtime abstraction can go
- where your protocol could stay portable beyond vendor-specific skills

Reference:

- [OpenHands](https://openhands.dev/)

### 5. Goose / OpenCode — extensible open agent surfaces

Why it matters:

- useful evidence that modular and open coding-agent ecosystems are real, not hypothetical

What to learn:

- extensions
- open runtime ergonomics
- terminal-first workflows

References:

- [Goose quickstart](https://block.github.io/goose/docs/quickstart/)
- [OpenCode](https://opencode.ai/)

---

## What We Should Incubate

If we do this well, erwin could incubate a **harness stack** with these layers:

### Layer 1 — Protocols

- autonomous development lifecycle protocol
- contract protocol
- release protocol
- cross-task communication protocol

### Layer 2 — Repo Legibility System

- docs layout
- index files
- execution plan conventions
- quality scoring documents
- custom lints
- generated references

### Layer 3 — Runtime Adapters

- Anthropic Skills adapter
- Codex Skills adapter
- maybe later OpenHands / Copilot adapter

### Layer 4 — Project Profiles

- generic starter
- `vi_agent`
- `erwin`

### Layer 5 — Eval / Validation

- static lint
- smoke tests
- behavior validation
- parity matrices
- recurring cleanup/garbage-collection automations

This would be a genuine harness engineering program, not just a renamed skill folder.

---

## What We Should Avoid

1. Building a full custom runtime before the protocol is stable
2. Locking the whole system to one vendor's skill format
3. Letting one practice repo define the universal model
4. Treating docs as optional instead of part of the harness
5. Confusing product naming with architectural boundaries

---

## Proposed Incubation Thesis

Here is the thesis I think best fits erwin:

> Erwin should incubate not a monolithic AI harness product, but a portable harness engineering system: protocol specs, repo-legibility conventions, eval loops, and runtime adapters that make autonomous software work reliable across agent runtimes.

This is differentiated.
It is also much closer to your real body of work than "launch a new CLI."

---

## Practical Next Steps

### P0

1. Make the protocol core explicit
2. Write a Codex runtime adapter spec
3. Build a parity matrix from teamwork/wing to protocol core

### P1

1. Add a harness-legibility doc system inside erwin itself
2. Define what "agent-legible repo" means for your methodology
3. Turn your testing ideas into a reusable eval stack

### P2

1. Package these as reusable incubator artifacts
2. Try the protocol/adapters on more than one repo
3. See whether a broader "wing" platform identity still makes sense after that

---

## Final Recommendation

Yes, use this moment to incubate your own harness engineering.

But do it at the right layer:

- not by cloning Codex or Claude
- not by renaming teamwork again
- but by extracting the portable system underneath your best work and making it runnable across harnesses

That is both more defensible and more aligned with erwin's mission.
