# Harness Landscape and Wing Positioning

> Date: 2026-03-28
>
> Objective: answer three questions:
>
> 1. What "harness" is hot right now?
> 2. What is the relationship between those harnesses and our actual goal?
> 3. How should `teamwork`, `wing`, the protocol core, runtime concerns, and `vi_agent` residue be decomposed?

---

## Executive Take

The most important update is conceptual:

**The hot thing right now is not a single harness name. It is the shift from prompt engineering to harness engineering.**

In other words:

- 2024 was "prompt better"
- 2025 was "build agents"
- 2026 is increasingly "engineer the environment around agents"

This matters because your current confusion is not mainly about naming.
It is about mixing up three layers:

1. the **protocol** you actually care about
2. the **harness/runtime** used to execute that protocol
3. the **product identity** you want the protocol or runtime to have

`teamwork`, `wing`, Claude Code skills, Anthropic skills, Codex skills, OpenHands, Goose, and OpenCode do not all live at the same layer.

That is why naming feels muddy right now.

---

## 1. What Is a "Harness" in This Context?

For this discussion, I am using **harness** to mean:

> the runtime shell around a model that makes real work possible and repeatable

This usually includes:

- tool access
- sandboxing and permissions
- environment bootstrapping
- memory and repository-local knowledge
- skills/plugins/extensions
- planning and task routing
- validation and eval loops
- automation/background execution
- review, retry, and recovery behavior

This is broader than a "framework" and more concrete than "an agent."

The model is not the harness.
The harness is the system that makes the model useful in production work.

---

## 2. What Seems Hot Right Now

### 2.1 The hot idea: harness engineering itself

The clearest first-party signal is OpenAI's recent piece on **harness engineering**:

- OpenAI explicitly argues that the repo, tools, docs, lints, plans, and operational constraints are what let coding agents scale, not prompt cleverness alone.
- The article rejects the "one giant `AGENTS.md`" model and instead promotes progressive disclosure, repo-local knowledge, custom linters, execution plans, and agent-legible architecture.

This is extremely close to the underlying problem you were already circling in `teamwork -> tw2 -> wing`.

So the most "hot" thing is:

**the harness as an engineered system of constraints, docs, tools, and workflows**

not just "which chat UI wins."

---

## 3. Named Harnesses / Platforms That Look Most Relevant

Below is the most useful current segmentation for your problem.

### 3.1 Codex

Why it matters:

- Codex is now positioning itself as a command center for multi-agent coding work, with built-in worktrees, cloud environments, skills, and automations.
- OpenAI is using the term "harness engineering" explicitly.
- The product direction is not "chat that can code," but "runtime for real engineering work."

What makes it strategically important:

- strong view on repository-legibility
- skills as reusable work packets
- automations as background recurring labor
- multi-agent workflows as a first-class concern

Why it matters to you:

- this is the first environment you are considering seriously after Claude/Anthropic
- your future system probably needs to map well onto this harness style

---

### 3.2 Anthropic Skills / Claude Skills

Why it matters:

- this is the environment `wing` was explicitly rewritten for
- Anthropic has now pushed Skills as an open standard
- the skills structure is intentionally based on progressive disclosure

What makes it strategically important:

- `Skill.md` as minimal entry point
- directory packaging
- scripts/resources as modular support files
- "keep it focused" is an official design principle

Why it matters to you:

- `wing` is not an abstract architecture exercise; it is a direct answer to this runtime model

---

### 3.3 OpenHands

Why it matters:

- OpenHands appears to be one of the strongest open-source "agent harness" platforms right now
- it positions itself as an open platform for cloud coding agents
- it has both SDK and CLI shapes, plus enterprise deployment paths

What makes it strategically important:

- more of a platform/runtime than just a CLI
- cloud and multi-agent orientation
- strong open-source legitimacy

Why it matters to you:

- if your real target is a portable protocol that outlives Claude-vs-Codex specifics, OpenHands is one of the clearest evidence points that "protocol above harness" is the right abstraction

---

### 3.4 Goose

Why it matters:

- Goose is a good example of an extensible open agent with modular extensions and MCP-friendly thinking

What makes it strategically important:

- extension model
- configurable autonomy
- desktop + CLI story

Why it matters to you:

- it reinforces that the future is modular runtime capability composition, not single giant instruction files

---

### 3.5 OpenCode

Why it matters:

- OpenCode appears to have very large open-source momentum and is explicitly positioned as an open-source coding agent

What makes it strategically important:

- terminal-first coding agent
- model-provider flexibility
- strong open-source adoption signal

Why it matters to you:

- it suggests there is real market pull for reusable agent work patterns above the model layer
- but it does not by itself solve your naming/positioning problem

---

## 4. What Is Actually "Hot" vs What Is Merely "Popular"

These are not the same.

### Popular

Platforms with visible adoption, stars, community momentum:

- OpenCode
- OpenHands
- Goose
- Codex
- Claude Skills

### Hot in the deeper sense

The more important shift is:

1. skills as reusable workflow units
2. automations/background execution
3. repo-local knowledge as the system of record
4. harness engineering as the force multiplier
5. moving from single-agent chats to structured multi-agent work

This means the trend is not "everyone uses harness X."

The trend is:

**serious teams are moving from model-first thinking to harness-first thinking**

---

## 5. Relationship Between Harnesses and Your Goal

Here is the key insight:

**Your goal is not to invent yet another harness.**

Your goal is to define a protocol that can inhabit multiple harnesses cleanly.

That protocol currently exists in a mixed, historically layered form across:

- `teamwork`
- `wing`
- `vi_agent` practice
- Claude/Anthropic runtime assumptions

So the correct relationship is:

### Harnesses are execution environments

Examples:

- Claude Code commands
- Anthropic official skills
- Codex skills
- OpenHands SDK/CLI
- Goose extensions

### Your actual asset is a protocol

Examples from your system:

- GitHub-first lifecycle
- Issue -> Claim -> Drive -> Ship -> RC
- Mission Contract
- Standalone Test
- Blast Radius
- Ripple checks
- Anti-compliance
- deterministic scripts for mechanism

### Therefore

Your work should increasingly be:

**protocol design + runtime adapters**

not:

**one runtime-specific skill suite that also carries the protocol identity**

---

## 6. Why Teamwork Feels Wrong to Move Directly

Because `teamwork` is actually a layered artifact:

1. protocol ideas you want
2. Claude Code command packaging you no longer want
3. historical operational fixes
4. some residual `vi_agent`-shaped pressure from real-world practice

That makes direct migration expensive and conceptually noisy.

You would be moving:

- good protocol
- plus old runtime
- plus old scars

all at once.

That is exactly the kind of migration that recreates baggage.

---

## 7. Why Wing Also Feels Wrong to Move Directly

This is subtler.

`wing` is cleaner than `teamwork`, but it now has a split identity:

### Wing as implemented

- modular rewrite of teamwork for Anthropic official skills

### Wing as named and imagined

- something able to carry V4 multi-agent orchestration
- something broader than just "team workflow"

This means `wing` has drifted from being "the teamwork refactor" into being "the candidate shell for a broader orchestration future."

That is probably why it no longer feels like a clean target for direct Codex migration either.

Because if you make "Codex wing" now, you are implicitly deciding:

- that `wing` is the durable product identity
- that its future scope includes both current lifecycle protocol and broader orchestration

That may be true, but it has not been made explicit enough yet.

---

## 8. What the Records Suggest About Wing's Real Positioning

### 8.1 Early tw2 record already reframed the product

In [`_project/ideas/skills/tw2-architecture.md`](/Users/liyasong/Project/collov/erwin/_project/ideas/skills/tw2-architecture.md):

- the core value is explicitly stated as **not** "team collaboration tool"
- it is reframed as an **AI autonomous development lifecycle protocol**

That is a major clue.

### 8.2 Wing naming record pushes beyond teamwork

In [`_project/sessions/2026-03-18-AM.md`](/Users/liyasong/Project/collov/erwin/_project/sessions/2026-03-18-AM.md):

- `wing` is chosen partly because it can carry the future V4 multi-agent direction

That means the name was selected not only for refactoring cleanliness, but also for future semantic headroom.

### 8.3 V4 record confirms this widening

In [`_project/ideas/skills/v4-multi-agent-orchestration.md`](/Users/liyasong/Project/collov/erwin/_project/ideas/skills/v4-multi-agent-orchestration.md):

- V4 is described as building on `tw2`
- and `tw2` is explicitly treated as the architecture prerequisite for a broader multi-agent coordination system

This means:

- `wing` is not just "teamwork but skill-shaped"
- it is already positioned as a stepping stone toward a larger orchestration layer

---

## 9. The Real Recommendation

You should stop asking:

- "Should I migrate teamwork?"
- "Should I migrate wing?"

and instead ask:

- "What is the protocol core?"
- "What runtime adapters do I need?"
- "What name belongs to the protocol vs the runtime vs the broader platform?"

That is the clean decomposition.

---

## 10. Proposed Layer Model

### Layer 1 — Protocol Core

This is the durable asset.

Contains:

- lifecycle
- contract structure
- decision heuristics
- quality gates
- GitHub-first semantics
- release semantics
- human/AI division of labor

This should be runtime-neutral.

### Layer 2 — Runtime Adapters

One adapter per harness/runtime.

Examples:

- Claude Commands adapter
- Anthropic Skills adapter
- Codex Skills adapter
- maybe later OpenHands adapter

These should own:

- path conventions
- permission idioms
- tool bindings
- skill package layout
- automation hooks

### Layer 3 — Project Profiles

Examples:

- `vi_agent` profile
- `erwin` dogfood profile
- future generic project template

These should own:

- config defaults
- branch conventions
- deploy workflow names
- repo-specific service maps

Once you think in these layers, the confusion reduces a lot.

---

## 11. Five-Column Decomposition Table

Columns:

- `Teamwork`
- `Wing`
- `Protocol Core`
- `Runtime / Harness Layer`
- `vi_agent Residue`

| Concern | Teamwork | Wing | Protocol Core | Runtime / Harness Layer | vi_agent Residue |
|---|---|---|---|---|---|
| Primary identity | GitHub-first team coordination skill suite | Anthropic-skill-structured rewrite with widened future headroom | AI autonomous development lifecycle protocol | How a given agent runtime executes the protocol | Real-world practice pressure from one heavily used product repo |
| Packaging form | Claude Code flat commands in `~/.claude/commands/*.md` | Anthropic-style `SKILL.md + operations + reference + scripts` | Should be packaging-agnostic | Owns package layout and invocation model | Historical validation happened through copied installs |
| Naming semantics | "team/teamwork" frame | Broader "wing" frame, no longer only teamwork semantics | Should probably use a neutral internal name/spec | Adapter name may vary by runtime | Some expectations may still be shaped by old repo usage |
| Lifecycle | `issue -> claim -> drive -> ship -> rc` | Same lifecycle retained | This is core and should be preserved | Runtime only changes how each step is invoked | None |
| Mission Contract | Local execution artifact from GitHub Issue | Same concept retained | Core | Adapter defines storage path and parser details | Past issue schemas and edge cases discovered via vi_agent |
| GitHub-first model | Central | Central | Core | Runtime-independent | Validated heavily in vi_agent |
| Cognitive posture | Strong, explicit, documented | Preserved, leaner, more structural | Core | Runtime should not redefine it, only deliver it | Some fixes came from real vi_agent failure cases |
| VE rules | Duplicated inline + docs | Centralized in `reference/ve-standard.md` | Not core protocol; presentation support | Runtime/presentation layer | Example formatting sometimes still reflects old repo use cases |
| Script layer | `tw-*` scripts | `wing-*` scripts | Deterministic mechanism belongs conceptually near core | Actual script invocation/pathing is adapter-specific | Some compatibility work came from vi_agent schema variations |
| Config home | `.teamwork/` plus `.teamspace/` compat | `.wing/` with `.teamwork`/`.teamspace` fallback | Protocol should only define required fields/semantics | Adapter chooses canonical state home | `.teamspace` and nested-schema history came from vi_agent-era needs |
| Worktree model | Claude/teammate workflow assumptions | Anthropic skill runtime assumptions | Optional protocol capability, not identity | Runtime owns how worktrees/tabs/parallelism are exposed | Many edge cases were discovered in vi_agent |
| Review / ship / RC flows | Mature, battle-tested | Preserved but modularized | Core workflow semantics | Runtime owns UX surface and tool calls | Some release and deploy assumptions were sharpened by vi_agent practice |
| Multi-agent ambition | Mostly future-facing / limited | Name chosen partly to allow V4 future | Could become protocol extension, but not mandatory for v1 core | Runtime may offer native multi-agent features | vi_agent may have over-shaped urgency and concrete scenarios |
| Documentation burden | Huge inline burden + rich docs | Lean runtime files, thinner suite docs | Protocol should have clear spec docs outside hot path | Runtime docs should stay separate from invocation files | Historical manual burden partly came from supporting live rollout |
| Current problem | Hard to iterate; mixed protocol + runtime + scars | Cleaner, but identity drifted beyond simple teamwork replacement | Not yet extracted cleanly as its own thing | Needs explicit Codex adapter if moving here | Hard to tell where generic design ended and repo-driven evolution began |

---

## 12. So What Should Happen Next?

### Recommendation

Do **not** directly migrate `teamwork` as-is.

Do **not** automatically assume `wing` is the Codex-era product name either.

Instead:

1. extract the protocol core explicitly
2. define runtime adapters explicitly
3. decide separately whether `wing` is:
   - the name of the broad future orchestration platform, or
   - the name of the Anthropic-era adapter, or
   - a temporary bridge artifact

Only after that should you build the Codex-native implementation.

---

## 13. My Best Current Read

If I had to compress my view into one sentence:

`teamwork` is a historically layered protocol artifact, `wing` is an Anthropic-era architectural rewrite with broader ambitions, and what you now need for Codex is neither a direct teamwork port nor a blind wing port, but a clean protocol extraction plus a Codex runtime adapter.

That is the least confusing path and the one most aligned with where the ecosystem is going.

---

## Sources

External research:

- OpenAI, "Harness engineering: leveraging Codex in an agent-first world"  
  [https://openai.com/index/harness-engineering/](https://openai.com/index/harness-engineering/)
- OpenAI, "Introducing the Codex app"  
  [https://openai.com/index/introducing-the-codex-app/](https://openai.com/index/introducing-the-codex-app/)
- OpenAI, Codex product page  
  [https://openai.com/codex/](https://openai.com/codex/)
- Anthropic / Claude Help, custom Skills  
  [https://support.claude.com/en/articles/12512198-how-to-create-custom-skills](https://support.claude.com/en/articles/12512198-how-to-create-custom-skills)
- Anthropic official skills repository  
  [https://github.com/anthropics/skills](https://github.com/anthropics/skills)
- OpenHands official site  
  [https://openhands.dev/](https://openhands.dev/)
- OpenHands repository overview  
  [https://github.com/OpenHands/OpenHands](https://github.com/OpenHands/OpenHands)
- Goose docs  
  [https://block.github.io/goose/docs/quickstart/](https://block.github.io/goose/docs/quickstart/)
- OpenCode official site  
  [https://opencode.ai/](https://opencode.ai/)

Project records:

- [`_project/ideas/skills/tw2-architecture.md`](/Users/liyasong/Project/collov/erwin/_project/ideas/skills/tw2-architecture.md)
- [`_project/ideas/skills/v4-multi-agent-orchestration.md`](/Users/liyasong/Project/collov/erwin/_project/ideas/skills/v4-multi-agent-orchestration.md)
- [`_project/sessions/2026-03-18-AM.md`](/Users/liyasong/Project/collov/erwin/_project/sessions/2026-03-18-AM.md)
- [`_project/CURRENT.md`](/Users/liyasong/Project/collov/erwin/_project/CURRENT.md)
- [`_project/GOALS.md`](/Users/liyasong/Project/collov/erwin/_project/GOALS.md)
