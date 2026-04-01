# Coding Agent Self-Iteration And Paper Intake

> Date: 2026-04-02
>
> Purpose: capture what current first-party evidence says about how Codex and Claude Code iterate on themselves, and explain why `erwin` should add a paper-aware knowledge intake lane instead of relying only on first-party engineering/news surfaces.

---

## Bottom Line

Two things are true at once:

1. `Codex` and `Claude Code` are both already being used to improve the workflows around their own development.
2. `erwin` is currently underweight on research-paper intake relative to the importance of papers for self-improvement, agent evals, and long-horizon coding-agent design.

The immediate implication is not “directly mutate `Harness` from papers.”

The immediate implication is:

- keep the existing first-party radar
- add a paper-aware intake lane
- route any paper-derived insight through the existing proposal-first learning loop

---

## What OpenAI Is Doing With Codex

Official evidence points to an eval-driven self-improvement posture.

Primary source:

- OpenAI GPT-5.2-Codex system card:
  - https://cdn.openai.com/pdf/ac7c37ae-7f4c-4442-b741-2eabdeaf77e0/oai_5_2_Codex.pdf

What matters:

- OpenAI explicitly has an `AI Self-Improvement` section in the Codex system card.
- That section frames self-improvement as a benchmarked capability, not as an uncontrolled autonomous loop.
- The official eval set includes:
  - `OpenAI PRs`
  - `MLE-Bench`
  - `PaperBench`
  - `OpenAI-Proof Q&A`

The key pattern is:

- OpenAI measures progress toward “automating the job of an OpenAI research engineer”
- it does that with agentic rollouts on internal PR-style tasks, ML engineering tasks, paper replication tasks, and internal debugging/diagnosis tasks
- the work is scaffolded and graded, not trusted as free-form self-rewrite

Interpretation:

- Codex self-iteration is best understood as `eval-backed research-engineering automation`
- it is not “let the agent rewrite itself however it wants”
- the strongest signal for `erwin` is the importance of PR replication, ML engineering, paper replication, and bottleneck diagnosis as explicit eval objects

---

## What Anthropic Is Doing With Claude Code

Official evidence points to a dogfooding plus workflow-automation posture.

Primary source:

- How Anthropic teams use Claude Code:
  - https://claude.com/blog/how-anthropic-teams-use-claude-code

What matters:

- Anthropic describes internal teams using Claude Code for:
  - codebase navigation
  - testing and code review
  - debugging and troubleshooting
  - prototyping and feature development
  - documentation and knowledge management
  - automation and workflow optimization
- The strongest “self-iterating” signal is that teams set up autonomous loops where Claude Code writes code, runs tests, and iterates continuously on new features.
- Anthropic also highlights that people used Claude Code to build tools and workflows that support Anthropic’s own operations.

Interpretation:

- Claude Code self-iteration is best understood as `heavy internal dogfooding plus autonomous workflow loops`
- unlike the OpenAI system-card framing, Anthropic’s public material emphasizes lived internal usage more than explicit self-improvement benchmarking
- the strongest reusable lesson for `erwin` is that product iteration, documentation consolidation, and automation tooling should all be treated as first-class input surfaces

---

## Why Papers Matter Here

The current radar is strong on first-party engineering and product surfaces.

That is good, but incomplete.

For this specific topic, papers are not optional background reading. They are where much of the field’s sharpest knowledge actually lands:

- self-refinement mechanisms
- agent feedback loops
- long-horizon coding/ML benchmarks
- paper-replication evals
- research-agent orchestration patterns

Right now `docs/radar/sources.md` does not include a dedicated paper lane.

That means `erwin` can see what vendors say they are doing, while still missing:

- benchmark design
- replication methods
- failure modes
- task decompositions
- proposed learning/self-correction mechanisms

For a repo trying to reason about `Harness`, that is a meaningful blind spot.

---

## Papers And Benchmarks To Track

These are the most relevant immediate additions for this topic.

### Self-improvement and self-correction

- `Self-Refine: Iterative Refinement with Self-Feedback`
  - https://arxiv.org/abs/2303.17651
  - Why it matters: canonical self-feedback loop; useful for understanding bounded refinement rather than unconstrained self-editing.

- `Reflexion: Language Agents with Verbal Reinforcement Learning`
  - https://arxiv.org/abs/2303.11366
  - Why it matters: explicit feedback-memory-reflection loop for agents; relevant to council/adoption/eval design.

- `Recursive Introspection: Teaching LLM Agents How to Self-Improve`
  - https://openreview.net/forum?id=ze0nodITam
  - Why it matters: more direct “self-improve after seeing failed attempts” framing.

### Benchmarks for AI engineering and research automation

- `MLE-bench: Evaluating Machine Learning Agents on Machine Learning Engineering`
  - https://arxiv.org/abs/2410.07095
  - Why it matters: concrete benchmark for long-horizon ML engineering, exactly the kind of task family OpenAI uses in self-improvement evaluation.

- `PaperBench: Evaluating AI's Ability to Replicate AI Research`
  - https://arxiv.org/abs/2504.01848
  - Why it matters: paper replication is a direct bridge between coding agents, research agents, and self-improvement capability tracking.

---

## Implications For erwin

### 1. The current knowledge loop is missing a paper-native lane

That is a real gap.

Current automations are good at:

- first-party publisher monitoring
- weekly synthesis
- GitHub trending intake

They are weak at:

- targeted paper discovery
- paper triage
- benchmark watch
- research-to-proposal extraction

### 2. First-party and paper intake should be separated, not merged into one noisy daily feed

If we just stuff papers into the current daily radar, signal quality will probably get worse.

A better structure is:

- `Daily Radar`
  - still first-party-led
  - still inbox-first
  - may mention only truly exceptional paper releases

- `Weekly Synthesis`
  - synthesizes first-party plus the strongest paper implications

- `Paper Watch`
  - separate, slower, more selective
  - focused on agents, coding agents, evals, self-improvement, orchestration, and observability

### 3. Paper intake should feed proposals, not direct runtime changes

This remains the right rule.

The fix is not “let papers rewrite Harness.”

The fix is:

- add paper intake
- extract proposal candidates
- route them into the proposal ledger

---

## Recommended Automation Reorg

### Keep

- `AI Agent Daily Radar`
- `AI Agent Weekly Synthesis`

### Update

- `AI Agent Daily Radar`
  - remove the need to carry GitHub Trending inside the same report
  - keep the daily report focused on first-party field signals
  - optionally add one tiny `Research Watch` note only when a paper is unusually important

- `AI Agent Weekly Synthesis`
  - add a dedicated `Papers And Benchmarks` section
  - force explicit comparison between first-party claims and paper evidence
  - require a `Knowledge Gaps` section naming what the current automation stack still cannot see well

### Add

- `AI Agent Papers Weekly`
  - weekly rather than daily
  - scan:
    - arXiv
    - OpenReview
    - PMLR proceedings
    - ACL Anthology
  - filter for:
    - coding agents
    - agent orchestration
    - self-improvement / self-correction
    - long-horizon tasks
    - evals / benchmarks
    - observability / tool use / context engineering
  - output:
    - concise inbox item
    - durable note under `docs/radar/papers/YYYY-MM-DD.md`
    - at most 3 proposal candidates

### Demote

- `GitHub Trending Daily`
  - keep it if useful
  - but do not let it dominate the main field-intelligence loop
  - trending is weaker than first-party engineering and weaker than benchmark/paper intake for this specific mission

---

## Recommendation

`erwin` should treat paper intake as a missing but now necessary layer.

Recommended next move:

1. update `docs/radar/sources.md` to add primary paper surfaces
2. record the automation gap in the proposal ledger
3. update the existing daily/weekly automations
4. add a separate weekly paper-watch automation

That keeps the current learning loop intact while fixing a genuine knowledge blind spot.
