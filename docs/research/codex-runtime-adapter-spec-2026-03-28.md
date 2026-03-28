# Codex Runtime Adapter Spec

> Date: 2026-03-28
>
> Purpose: define how the autonomous development lifecycle protocol should be adapted to the Codex runtime without dragging forward `teamwork` packaging debt or overfitting to the Anthropic `wing` implementation.

---

## Executive Position

Codex migration should not be framed as:

- "port `teamwork`"
- "port `wing`"
- "rebuild Claude commands in another syntax"

It should be framed as:

**implement a Codex-native runtime adapter for the autonomous development lifecycle protocol**

This keeps three concerns separate:

1. `protocol core`
2. `runtime adapter`
3. `product surface / naming`

That separation is the main lesson from `teamwork -> wing`.

---

## Design Goal

The Codex adapter must preserve the protocol's behavioral guarantees while changing the execution surface to fit Codex reality:

- skills instead of Claude commands
- `exec_command` / `apply_patch` / `update_plan` instead of Anthropic tool primitives
- explicit local workspace instructions
- optional web verification
- sub-agents available but permission-constrained

The adapter should feel native to Codex, not like a compatibility layer pretending to still be Claude Code.

---

## Non-Goals

This adapter is not responsible for:

- defining the protocol itself
- preserving historical `teamwork` phrasing
- preserving `wing` naming
- recreating every old visual encoding convention
- preserving project-specific `vi_agent` assumptions
- becoming a full standalone agent runtime platform

Those are separate questions.

---

## What Must Survive

The following protocol semantics must survive unchanged:

- GitHub-native lifecycle: `issue -> claim -> drive -> ship -> rc`
- Issue as public source of truth
- Contract as AI execution artifact
- standalone execution requirement
- executor right to challenge bad sub-task decomposition
- deterministic operations mechanized outside freeform reasoning where possible
- config-driven conventions
- blast-radius awareness
- completion defined by actual state transition plus verification

If any runtime choice weakens those properties, the adapter is wrong.

---

## Runtime Reality: Codex

Codex changes the operating surface in several important ways.

### 1. Tools are strong, but not Claude-shaped

Codex gives us:

- shell execution
- file patching
- planning
- optional web research
- optional sub-agents

This is enough to support the protocol, but the operational choreography must be rewritten.

### 2. Multi-agent is not the default contract

`teamwork` and some later `wing` thinking assume stronger built-in team semantics.
In Codex, delegation exists but is explicitly permission-sensitive.

Therefore:

- protocol semantics must not require a team API
- parallelism must be optional acceleration
- single-agent execution must remain first-class

### 3. Repository instructions are a first-class input

Codex sessions are deeply shaped by repo-local `AGENTS.md`, current branch state, local file layout, and explicit tool rules.

Therefore the adapter must:

- read local context before acting
- respect repo instructions above skill defaults
- avoid hidden state assumptions

### 4. Skill shape is thinner than old command workflows

Anthropic skills can carry a large amount of staged procedural text.
Codex skills can do that too, but the better fit is:

- compact guidance in `SKILL.md`
- deterministic work in repo scripts when justified
- durable state in repo files only when mission length warrants it

The adapter should bias toward smaller skill text and clearer external artifacts.

---

## Recommended Codex Adapter Architecture

The adapter should be split into four layers.

### Layer 1: Protocol Spec

Source:

- runtime-neutral protocol documents

Responsibilities:

- lifecycle semantics
- invariants
- state model
- required objects

This layer already exists conceptually in [`protocol-core-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/protocol-core-spec-2026-03-28.md).

### Layer 2: Codex Adapter Skills

Recommended skill surfaces:

- dashboard / doctor skill
- issue skill
- claim skill
- drive skill
- ship skill
- rc skill

These should be Codex-native instructions, not transliterations of the Anthropic `wing` prompts.

Responsibilities:

- gather runtime context
- call deterministic mechanisms
- enforce adapter-specific guardrails
- render concise user-facing progress/state

### Layer 3: Deterministic Mechanisms

These can remain shell scripts or small utilities, as long as they are runtime-independent.

Examples:

- config resolution
- contract freshness
- git / worktree operations
- label transitions
- PR helpers
- RC helpers

This is the part most worth preserving from `teamwork` and `wing`.

### Layer 4: Project Profile

Per-repo config should define:

- base branch
- production branch
- label prefix conventions
- worktree policy
- verification commands
- staging workflow names
- release topology

No adapter logic should hardcode repo policy.

---

## Concrete Capability Mapping

| Protocol need | Codex primitive | Adapter pattern |
|---|---|---|
| inspect repo state | `exec_command` | `rg`, `sed`, `find`, `git`, `gh` |
| edit local artifacts | `apply_patch` | contracts, notes, generated docs |
| maintain active execution state | `update_plan` + optional files | plan in-session, state file only for long missions |
| current external verification | `web` | only when freshness matters |
| interactive clarification | direct chat | minimize, only on true decision forks |
| bounded parallel read work | `multi_tool_use.parallel` | preferred over shell chaining |
| multi-agent delegation | `spawn_agent` | optional, never protocol-required |

---

## Codex-Native Lifecycle Design

### A. Dashboard / Doctor

Codex version should merge "dashboard" and "doctor" concerns more tightly than legacy `teamwork`.

Responsibilities:

- detect runtime readiness
- detect active contract / branch / worktree mismatch
- detect issue / PR / review drift
- suggest next valid step

Codex advantage:

- better direct shell inspection
- easier inline summarization
- no need to preserve historical menu-heavy interaction

Recommended redesign:

- keep diagnosis
- keep auto-fix only for deterministic safe operations
- reduce ceremony
- make "next valid step" the center of output

### B. Issue

Issue creation and mutation should stay GitHub-first.

Responsibilities:

- create executable issues
- update issue structure when context changes
- batch generation from planning input
- optional duplicate / dependency surfacing

Recommended redesign:

- preserve standalone-test semantics
- preserve structured issue fields
- simplify prompt scaffolding
- separate issue quality linting from issue creation flow

### C. Claim

Claim is the most important adapter surface because it binds GitHub work to local execution context.

Responsibilities:

- choose issue
- validate executability
- generate contract
- create branch or worktree
- record ownership transition

Recommended redesign:

- preserve contract generation
- preserve branch/worktree setup
- preserve freshness metadata
- redesign old interactive branching trees into a more compact Codex flow
- treat auto-launch into drive as optional continuation, not hardwired ritual

### D. Drive

Drive in Codex should be where protocol meets harness engineering most directly.

Responsibilities:

- load contract
- inspect real codebase
- challenge plan if needed
- implement
- verify
- sync execution state back to contract and issue

Recommended redesign:

- preserve mission ownership and anti-compliance stance
- preserve contract freshness checks
- preserve task tracking
- remove assumptions about Claude-only task APIs
- integrate `update_plan` as the live session task view
- use durable state files only for long-running missions

### E. Ship

Ship remains a mostly deterministic operation with strong guardrails.

Responsibilities:

- verify local cleanliness / branch correctness
- run tests
- push branch
- create or update PR
- attach labels / reviewers / comments
- support post-merge cleanup

Recommended redesign:

- preserve PR lifecycle mechanics
- preserve review mode
- preserve sync / rebase support
- reduce overlong instructional prose
- split pure deterministic helpers from UI phrasing

### F. RC

RC should remain highly mechanized and minimally prompt-dependent.

Responsibilities:

- detect active RC
- derive next version
- cut RC branch
- trigger staging
- promote
- apply surgical RC fix
- hotfix production and backport

Recommended redesign:

- preserve full lifecycle
- preserve single active RC invariant
- preserve dual-branch support
- move more release logic into scripts and less into skill prose

---

## State and Artifact Model

Codex should use a lighter artifact model than legacy `teamwork`.

### Required durable artifacts

- issue on GitHub
- PR on GitHub
- mission contract file when a mission is active
- repo config for protocol behavior

### Optional durable artifacts

- session state file for long missions
- research notes for externally sourced work
- adapter-generated audit logs for evals

### Anti-patterns

- hidden state that only the model remembers
- mandatory restart rituals
- runtime progress stored only in freeform chat
- project policy embedded in skill prose

---

## Contract Format Guidance

The contract should survive migration, but the exact wording can change.

Required semantic sections:

- mission identity
- branch / worktree target
- objective
- success criteria
- sub-tasks
- code context
- verification commands
- freshness metadata
- execution notes / discoveries

Codex improvement opportunity:

- store explicit "assumptions invalidated" notes
- store "manual follow-up required" separately from executable tasks
- make the contract easier for another agent or another runtime to resume

---

## Harness Engineering Opportunity

This adapter can become the seed of a larger harness engineering system if built correctly.

The key is to make the Codex adapter one runtime implementation of a broader architecture:

- protocol core
- adapter interface
- deterministic mechanisms
- project profile
- eval layer

If done this way, Codex migration is not a one-off port.
It becomes the first clean test of a portable harness architecture.

---

## What To Preserve from Teamwork

Preserve mostly as mechanism or semantics:

- GitHub-first worldview
- contract lifecycle
- worktree support
- label state machine
- PR lifecycle helpers
- RC mechanics
- standalone-test discipline
- issue freshness / contract freshness thinking

These are real assets.

---

## What To Preserve from Wing

Preserve mostly as architectural lessons:

- skill-based modular packaging
- smaller prompt surfaces
- single-source VE reference
- config path fallback logic
- explicit operations split
- clearer install story

These reduce iteration friction.

---

## What To Deliberately Not Preserve

Do not preserve:

- `~/.claude/commands/...` path assumptions
- Anthropic-specific variables like `${CLAUDE_SKILL_DIR}` as the conceptual center
- menu-heavy interaction trees built around old `AskUserQuestion`
- historical text duplication
- `vi_agent`-shaped defaults that do not generalize
- prompt-enforced behavior that should be a script, check, or lint

---

## Recommended Codex Adapter Interface

The adapter should expose six primary entry points.

1. `protocol-home`
   Status, diagnosis, next valid step, recovery suggestions.

2. `protocol-issue`
   Create, refine, update, batch-generate executable issues.

3. `protocol-claim`
   Turn issue into active contract and execution branch/worktree.

4. `protocol-drive`
   Execute contract against code with verification and state updates.

5. `protocol-ship`
   Deliver via PR and handle review / merge / cleanup transitions.

6. `protocol-rc`
   Run release-candidate, promote, fix, hotfix flows.

These are interface names, not branding decisions.

---

## Naming Recommendation

Do not force the adapter to inherit the `wing` name yet.

Why:

- `wing` is now partly a future orchestration concept
- Codex migration is currently about protocol portability
- locking the Codex adapter to `wing` too early may collapse protocol and product layers again

Recommendation:

- build the adapter under neutral protocol-oriented naming during incubation
- decide later whether `wing` should become:
  - the umbrella orchestration brand
  - one adapter family
  - or a different layer entirely

---

## Acceptance Criteria

The Codex adapter is successful when all of the following are true:

1. A user can go from issue to PR using only Codex-native skills and deterministic helpers.
2. The adapter preserves protocol invariants without depending on Anthropic-specific runtime assumptions.
3. The adapter supports both single-agent and optional delegated execution without changing protocol semantics.
4. The adapter makes fewer prompt-level assumptions than `teamwork`, and carries less packaging debt than `wing`.
5. A project can adopt it by supplying config and profile data rather than editing skill internals.
6. The resulting design can plausibly support another runtime adapter later.

---

## Suggested Next Build Order

1. Define the adapter interface and file layout.
2. Build a parity matrix against `teamwork` and `wing`.
3. Decide which deterministic scripts can be reused unchanged.
4. Define a Codex-native contract format.
5. Create the dashboard / doctor adapter first.
6. Implement `claim` next, because it anchors all other flows.
7. Validate `drive -> ship -> rc` with smoke scenarios.

This sequence minimizes architecture thrash.
