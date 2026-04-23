# Drive Skill v2.3.0 — Data Formats

Formal definitions for the 3+1 document architecture.

## Architecture: "一个关切、一个文档"

| Document | Concern | Written By | Lifecycle | Consumed By |
|----------|---------|------------|-----------|-------------|
| `contract.md` | WHAT + WHY + Decision Model | Dialogue Agent | Immutable after approval | Everyone |
| `harness.md` | HOW TO VERIFY | Eval Architect | Controlled evolution | Evaluator (primary), Planner/Executor (§3 only) |
| `plan.md` | HOW TO BUILD + State + Log | Planner → Executor updates | Living document | Executor, Orchestrator |
| `grading/` | RESULTS + Evidence | Evaluator | Per-iteration output | Orchestrator, User |

### Orchestrator Validation Gates

```
⛔ POST-UNDERSTAND GATE:
   Assert: contract.md exists, >50 lines, has §1-§4, has features (| F)
   Assert: harness.md exists, >100 lines, has §1-§9, USER APPROVED
   HALT if any missing or incomplete.

⛔ POST-EXECUTE GATE:
   Assert: plan.md exists, all features show ✅ in §1
   Assert: build succeeds (harness.md §3)
   HALT if fails.

⛔ POST-EVALUATOR GATE (HARD EVIDENCE):
   Assert: grading/evidence/ has ≥3 files, each >10 lines
   Assert: grading/eval.md exists and >50 lines
   Assert: grading/eval.json has verdict PASS or FAIL
   HALT if any fails — re-dispatch Evaluator.
```

---

## Storage Layout

```
{mission_dir}/
├── contract.md     ← WHAT: mission + features + decision proxy + eval vision
├── harness.md      ← VERIFY: complete eval spec (§1-§9), co-designed with user
├── plan.md         ← HOW: waves + state + log (the only living document)
├── research/       ← Working files (ephemeral, not core artifacts)
│   └── recon.md    ← Recon Agent output → consumed by Dialogue Agent
└── grading/        ← RESULTS: evaluation output
    ├── eval.md     ← Latest evaluation report (overwritten each iteration)
    ├── eval.json   ← Cumulative structured data (all iterations)
    └── evidence/   ← All evidence: logs, screenshots, metrics
```

---

## contract.md — The Complete Mission Understanding

Written by Dialogue Agent during Phase 1. Stored at `{mission_dir}/contract.md`.

This single document replaces the v2.2 trio of contract.md + shadow.md + features.json. It captures everything about WHAT the mission is, WHY it matters, and HOW to make autonomous decisions.

```markdown
# Contract — {name}
> Status: {DRAFT | APPROVED} | Date: {ISO 8601}

## §1 Mission

{Essence — 为什么这个项目存在。一段话，不可约核心。}

### Context
{从侦察中精炼：技术栈、影响面、关键约束。简短。}

### Success Criteria
□ SC1: {具体、可验证}
□ SC2: {具体、可验证}

### Scope
- IN: {做什么}
- OUT: {不做什么}

## §2 Features

| ID | Title | Acceptance Criteria | Executor |
|----|-------|-------------------|----------|
| F01 | {title} | AC1: {testable} / AC2: {testable} | codex |
| F02 | {title} | AC1: {testable} | claude |

## §3 Decision Proxy

{执行者遇到 §1-§2 没覆盖的边界情况时，读这一节替用户做判断。}

### Success Landscape
- **最低线:** {non-negotiable minimum}
- **理想状态:** {what "done well" looks like}
- **超预期:** {genuine delight}

### Priority Model
- 不可妥协: {X}
- 有弹性: {Y}
- 不关心: {Z}

### 品味与判断
{工程品味、审美偏好、决策风格}

### 权衡逻辑
- {tradeoff}: 选{X} 因为{logic}

### 自主 vs 上报
- 自主: {types} → {action}
- 上报: {types} → {why}

### Good Enough 边界
- 完成信号: {when to stop}
- 过度设计信号: {doing too much}

## §4 Evaluation Vision

{什么叫"做好了"。喂给 Eval Architect 设计 harness.md。}

### 最重要的评估维度
### 必须测到的场景
### 最怕的失败模式
### 需要的证据

## Amendments
(append-only)
```

### Section Purposes

| Section | Purpose | Primary Consumer |
|---------|---------|-----------------|
| §1 Mission | What + Why at a glance | Everyone |
| §2 Features | Operationalized feature list | Planner, Executor, Evaluator |
| §3 Decision Proxy | Autonomous decision model | Executor (edge cases during execution) |
| §4 Evaluation Vision | Raw eval expectations | Eval Architect (input for harness.md) |

### Immutability Rules

- Once user-approved, §1-§4 are immutable
- Amendments are append-only with timestamp and rationale
- Only Dialogue Agent can amend, with user approval

---

## harness.md — The Evaluation Specification

Written by Eval Architect during Phase 1, iteratively co-designed with the user. Stored at `{mission_dir}/harness.md`.

**harness.md is the ONLY evaluation specification.** The Evaluator reads this file and executes it top to bottom.

### Required Sections

| Section | Content | Why |
|---------|---------|-----|
| §1 Mission Context | Minimum context for independent Evaluator | Evaluator has zero execution context |
| §2 Environment | Complete setup commands, health check, teardown | Must be executable |
| §3 Mechanical Checks | Build, lint, type check commands | First-pass verification |
| §4 Evaluation Rubric | Quantitative metrics (§4.Q) + subjective dimensions (§4.S), mixed-weight formula | Quality measurement |
| §5 Test Scenarios | Complete steps, expected results, evidence | Behavioral verification |
| §5.A Adversarial | Designed to catch fakes | Anti-gaming |
| §6 Evidence Protocol | Required artifacts with quality bars | Prevents empty evidence |
| §7 Verdict Protocol | PASS/FAIL rules, zero compromise | No ambiguity |
| §8 Evaluator Swarm — Agent Roster & Dispatch Protocol | Agent roster, dispatch protocol, parallel execution map | Swarm orchestration |
| §9 Calibration Notes | Known limitations, blind spots | Honesty about coverage |

### Header Format

```markdown
> Version: {N} | Status: {DRAFT / USER APPROVED}
> Mission: {one sentence}
> Task Type: {runtime-web / creative-generative / document-research / infrastructure-devops / ux-design-spec / skill-evolution}
> Approved after: {N} rounds
> Dry run: {PASSED / NOT RUN}
```

### User Approval Requirement

Header MUST contain `Status: USER APPROVED` before execution begins.

### Controlled Evolution

- Executors MUST NOT modify harness.md
- Only Eval Architect can amend (small fix: direct; structural: needs user approval)

### §4.Q Format — Agent Dispatch Spec

Each quantitative metric in §4.Q is an agent dispatch spec. The Evaluator Orchestrator spawns the named agent with the given prompt and tools, the agent writes evidence to the specified path, and the Orchestrator extracts the metric value.

```markdown
#### QM{N}: {name} ({feature})
**What:** {description}
**Agent Dispatch:**
  Agent: {agent-name}
  Prompt: "{full prompt with {placeholders}}"
  Tools: {tool list}
  Evidence: {evidence_dir}/{file}
  Extract: {field} from {format}
**Threshold:** {value}
**Calibration:** {description}
```

Required fields per QM:
- **Agent** — unique agent name for this metric (e.g. `taxonomy-verifier`, `template-counter`)
- **Prompt** — complete prompt the Orchestrator passes to the agent; must include what to read, what to check, and what JSON/MD to output
- **Tools** — exact tool list the agent needs (Read, Glob, Bash, WebSearch, Playwright MCP, etc.)
- **Evidence** — file path where the agent writes its output (always under `{evidence_dir}/`)
- **Extract** — which field from the evidence file the Orchestrator reads for threshold comparison

### §5 Format — Agent Mission

Each test scenario in §5 is an agent mission. Instead of sequential CLI steps, each scenario dispatches a specialized agent with a complete prompt.

```markdown
### S{N}: {name} [{type}, {priority}]
**Covers:** {features, success criteria}
**Priority:** {critical / high / medium}

**Agent Mission:**
  Agent: {agent-name}
  Prompt: "{full prompt with {placeholders}}"
  Tools: {tool list}
  Evidence: {evidence_dir}/{file}
  Pass Criteria: {specific criteria}
```

Required fields per scenario:
- **Agent** — unique agent name for this scenario (e.g. `taxonomy-integration-tester`)
- **Prompt** — complete mission prompt: what to verify, how, and what output to produce
- **Tools** — exact tool list
- **Evidence** — where the agent writes its structured report
- **Pass Criteria** — specific, unambiguous criteria the Orchestrator uses to determine PASS/FAIL

### §8 Format — Evaluator Swarm

§8 defines the complete agent swarm for evaluation execution. It replaces the old sequential "step-by-step for Evaluator" format with three components:

**Agent Roster Table:**
```markdown
| Agent | Mission Prompt Summary | Tools | Evidence Output | Phase |
|-------|----------------------|-------|-----------------|-------|
| {agent-name} | {one-line mission} | {tool list} | {evidence file} | {E1/E2/E2.5/E3/...} |
```

**Dispatch Protocol:**
```markdown
Phase E0: Orchestrator reads harness.md
Phase E1: dispatch {agent} → wait → gate check
Phase E2: dispatch {agent} → wait → gate check
Phase E2.5: parallel dispatch [QM agents] → wait all → extract metrics
Phase E3: parallel dispatch [Scenario agents] in groups → wait all
Phase E4: synthetic tests
Phase E5: dispatch judge → score subjective dimensions
Phase E6-E7: evidence audit + claims scan
Phase E8: compute verdict
Phase E9: write eval.md + eval.json
```

**Parallel Execution Map:**
```
ASCII diagram showing which agents can run simultaneously
and which have sequential dependencies.
```

---

## plan.md — The Execution Plan + Living State

Written by Planner, updated by Executor during execution. Stored at `{mission_dir}/plan.md`.

This single document replaces the v2.2 trio of plan.md + state.md + logs/.

```markdown
# Plan — {name}
> Last updated: {timestamp} | Wave: {N}/{total} | Features: {done}/{total}

## §1 Waves

### Wave 1: {theme}
**Design:** {key technical choices + rationale}
**Entry:** {preconditions} | **Exit:** harness.md §3 pass

| Feature | Executor | Depends | Status |
|---------|----------|---------|--------|
| F01     | codex    | —       | ✅     |
| F02     | claude   | F01     | 🔄     |

### Wave 2: {theme}
...

## §2 State

### Current
- Phase: {READY | EXECUTING | VERIFYING | COMPLETE}
- Active issues: ...

### Git Checkpoints
- project-{name}-wave-1: {sha} ({date})

### Next Actions
1. {what to do next}

### Next Step Logic
{Why the next action follows from current state}

## §3 Log

### Wave 1
**Narrative:** {What happened, why, what was learned. Prose.}
**Dead ends:** {Failed approaches — tried, failed because, avoid}
**Shadow decisions:** {Decisions via contract.md §3 + reasoning}
**Quantitative:** Mechanical: {P/F} | Features: 0 → 3
```

### Section Purposes

| Section | Purpose | Updated By | When |
|---------|---------|------------|------|
| §1 Waves | Execution plan + feature status | Planner (initial), Executor (status) | Plan creation, after each feature |
| §2 State | Living checkpoint | Executor, Orchestrator | After each wave |
| §3 Log | Cumulative execution record | Executor | After each wave (append-only) |

### Status Icons (§1 Feature Status)

| Icon | Meaning |
|------|---------|
| ⬚ | Not started |
| 🔄 | In progress |
| ✅ | Passing |
| ❌ | Failed |
| ⏸ | Blocked |

---

## grading/ — Evaluation Output

### eval.md — Latest Evaluation Report

Written by Evaluator. **Overwritten each iteration** (no anchoring).

Must be >50 lines. Must cite specific evidence for every score. See `agents/evaluator.md` for full format.

**Required sections include:**
- `## Quantitative Metrics Results` — table of all §4.Q metrics with command, result, threshold, PASS/FAIL. Machine-verifiable outputs only. Shows `quantitative_score` and the verdict formula: `final_score = quantitative_score × 0.70 + qualitative_score × 0.30`.
- `## Subjective Dimension Scores` — §4.S dimensions with evidence citations (≤30% weight).

### eval.json — Cumulative Structured Data

Written by Evaluator. **Cumulative across iterations** — contains `iterations_history`.

```json
{
  "mission": "{name}",
  "version": "2.3.0",
  "timestamp": "{ISO 8601}",
  "iteration": 2,
  "iterations_history": [
    {
      "iteration": 1,
      "timestamp": "{ISO 8601}",
      "verdict": "FAIL",
      "scores": {"Essence": 5, "Craft": 6},
      "quantitative_metrics": [{"name": "...", "pass": false}],
      "features_passed": 3,
      "features_total": 8
    }
  ],
  "features": [
    {"id": "F01", "title": "{title}", "passed": true, "evidence": "{ref}"}
  ],
  "quantitative_metrics": [
    {"name": "{QM name}", "command": "{bash command}", "result": "{number|string}", "threshold": "{threshold}", "pass": true}
  ],
  "quality_dimensions": [
    {"name": "{dim}", "score": 8, "threshold": 7, "evidence": "{one-line}"}
  ],
  "scenarios": [
    {"name": "{scenario}", "status": "pass|fail", "evidence": "{file}"}
  ],
  "evidence": ["{filename1}", "{filename2}"],
  "quantitative_score": 8.5,
  "qualitative_score": 7.2,
  "verdict_formula": "final_score = quantitative_score × 0.70 + qualitative_score × 0.30",
  "verdict": "PASS|FAIL",
  "verdict_reason": "{one-line}",
  "agent_roster": [
    {"name": "{agent-name}", "role": "{mission summary}", "tools": ["Read", "Glob"], "phase": "E2.5"}
  ],
  "agent_evidence": [
    {"agent": "{agent-name}", "file": "{evidence_dir}/{filename}", "lines": 42}
  ]
}
```

### VERDICT Protocol

- **PASS** — All criteria met, all thresholds passed, all evidence present
- **FAIL** — Any critical failure, any threshold not met, any crash

No PARTIAL. No "close enough."

**Verdict Formula:** `final_score = quantitative_score × 0.70 + qualitative_score × 0.30`

**⛔ No-Compensation Rule:** If ANY quantitative metric (§4.Q) is below its threshold, verdict = FAIL regardless of qualitative scores. Subjective dimensions CANNOT compensate for quantitative failures. 不允许主观补偿——量化指标是 bash 命令输出，不是 LLM 打分。

### evidence/ — All Evidence Files

Contains logs, screenshots, metrics, test outputs. Named descriptively:
- `env-setup.log`
- `mechanical-checks.log`
- `scenario-01-login-flow.log`
- `screenshot-dashboard.png`
- `synthetic-tests.md`

Minimum: ≥3 files, each >10 lines. Empty evidence = Evaluator did not do their job.

---

## Default Rubric Dimensions (Reference for harness.md §4)

Starting points for the Eval Architect to customize per project based on contract.md §4 Evaluation Vision.

**§4 splits into two sections:**
- **§4.Q Quantitative Metrics** (70% weight) — each QM is an **agent dispatch spec**: the Evaluator Orchestrator spawns a specialized agent with an explicit prompt, the agent collects evidence into a structured file, and the Orchestrator extracts a value to compare against the threshold. Each QM has: What, Agent Dispatch (agent name, prompt, tools, evidence path, extract field), Threshold, Calibration (typical first-attempt / threshold / gap reason). 量化指标不是 LLM 打分，也不再是内联 bash 脚本——而是 agent 任务规格书。
- **§4.S Subjective Dimensions** (≤30% weight) — LLM-judge scored. Only for dimensions that genuinely cannot be mechanically computed.

> **Legacy note:** Pre-v2.3.1 harnesses used `Compute:` with inline bash code blocks for §4.Q. The new **Agent Dispatch** format replaces this. Existing harnesses using the old format remain valid but new harnesses MUST use Agent Dispatch.

The default dimensions below are **§4.S subjective dimension** starting points. The Eval Architect must ALSO discover §4.Q quantitative metrics per Step 2.6 and express them as Agent Dispatch specs.

### Dimension: Essence (gating — threshold ≥ 7)
**Core question:** Did the implementation solve the RIGHT problem?

| Level | What triggers this score |
|-------|------------------------|
| **3** | Addresses a different problem than asked |
| **5** | Partially addresses the core need |
| **7** | Core need met. Minor interpretation gaps. |
| **10** | Reframes and solves the deeper need |

### Dimension: Craft (threshold ≥ 6)
**Core question:** Is the code well-built and maintainable?

| Level | What triggers this score |
|-------|------------------------|
| **3** | Works but unmaintainable |
| **5** | Functional and readable, structural weaknesses |
| **7** | Production-quality, follows conventions |
| **10** | Exemplary — makes codebase better |

### Dimension: Completeness (threshold ≥ 7)
**Core question:** Are all requirements accounted for?

| Level | What triggers this score |
|-------|------------------------|
| **3** | Major features missing |
| **5** | Core present, significant gaps |
| **7** | All contracted criteria met |
| **10** | All criteria + edge cases + polish |

### Calibration Instructions

1. **Anchor first, score second.** Read ALL level examples before scoring.
2. **Evidence requirement.** Every score MUST cite specific evidence.
3. **Score independently.** Don't let dimensions compensate each other.
4. **When in doubt, score lower.** Burden of proof is on the output.

---

## Context Reset Protocol

When context is >50% used during execution:

```
1. Read contract.md (mission + features + decision proxy)
2. Read plan.md §2 State → Next Actions
3. Read harness.md §3 (mechanical checks)
4. Resume from plan.md §2 Next Actions
```

No separate bootstrap.md file needed.

---

## Eliminated Files (v2.3)

| v2.2 File | Now Lives In | Why |
|-----------|-------------|-----|
| `shadow.md` | `contract.md §3 Decision Proxy` | Shadow IS the contract's decision model |
| `features.json` | `contract.md §2 Features` | Features ARE the contract |
| `state.md` | `plan.md §2 State` | State IS the execution plan |
| `bootstrap.md` | Protocol in SKILL.md | Standard protocol, not a file |
| `logs/wave-{N}.md` | `plan.md §3 Log` | Narrative IS the log |
| `logs/dead-ends.md` | `plan.md §3 Log` | Dead ends ARE execution knowledge |
| `research/design-decisions.md` | `plan.md §1 Waves (Design)` | Decisions inline with waves |
| `harness/harness.md` | `harness.md` (root) | Flattened — only one file |
| `grading/eval-{N}.md` | `grading/eval.md` | One report, overwritten per iteration |
| `grading/eval-data.json` | `grading/eval.json` | Renamed, now cumulative |
| `grading/eval-comparison.json` | `grading/eval.json` `iterations_history` | Comparison data inside eval.json |
| `grading/eval-report.md` | `grading/eval.md` | IS the eval report |
| `grading/screenshots/` | `grading/evidence/` | Screenshots ARE evidence |
