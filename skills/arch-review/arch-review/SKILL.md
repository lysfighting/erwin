---
name: arch-review
description: "Multi-role debate review of existing architecture documents. Complementary to /architect (which designs new systems). Try: /arch-review help"
version: "0.1.0"
disable-model-invocation: true
allowed-tools: Read, Grep, Glob, Bash, Agent, SendMessage, TeamCreate, TeamDelete, AskUserQuestion
---

# /arch-review — Architecture Document Debate Review

> Review existing documents through multi-role debate. Find errors, gaps, and contradictions
> that a single reader would miss. Every finding must have evidence and a concrete action.

> **Complementary to /architect**: architect designs new systems, arch-review audits existing documents.

**User input**: $ARGUMENTS

If `$ARGUMENTS` is empty, `help`, or `-h`, output the following and **STOP**:

**`/arch-review`** — Multi-role debate review of architecture documents

**USAGE**
  `/arch-review path/to/document.md`       Review a single document
  `/arch-review --debug path/to/docs/`     Review + record process log
  `/arch-review revolution`                Analyze latest debug log → self-optimize
  `/arch-review help`                      Show this help

**WHAT IT DOES**
  1. You define the data sources and their authority levels
  2. Roles are assigned with explicit data access permissions
  3. Each role reviews independently, then findings are debated
  4. Consensus findings are formatted as actionable items

**MECHANISMS**
  `M1` Data Authority Hierarchy — prevents wrong-source-of-truth errors
  `M2` Role Templates — each role has allowed/forbidden data sources
  `M3` Lessons Injection — accumulated review lessons auto-embedded
  `M4` Standard Output — Location + Evidence + Impact + Action + Priority
  `M5` Debate Orchestration — independent review → cross-validation → consensus
  `M6` Self-Sufficiency Test — can a zero-context reader use this document?
  `M7` Debug Recording — `--debug` flag records full process log for analysis
  `M8` Revolution — analyze debug logs → propose skill-level improvements

**EXAMPLES**
  `/arch-review docs/architecture/`
  `/arch-review api-design-v2.md`

---

## Argument Routing

Parse `$ARGUMENTS`:
- If starts with `revolution`: jump to Revolution (below) — auto-reads latest debug log
- If contains `--debug`: set `DEBUG_MODE = true`, remove flag from path, continue to Phase 0
- Otherwise: `DEBUG_MODE = false`, continue to Phase 0

**Debug log directory**: `.claude/arch-review/debug/` (fixed location, auto-created).
Filename: `debug-YYYY-MM-DD-HHMMSS.md` (timestamped, never overwritten).

If `DEBUG_MODE` is true, load the debug recording protocol:
```
Read ${CLAUDE_SKILL_DIR}/operations/debug.md
```
Throughout Phases 0-5, append to the debug log as specified in debug.md.

---

## Phase 0: Triage — Is a Full Debate Needed?

Read the review target from `$ARGUMENTS` (with `--debug` flag already stripped). Skim the document(s).

Read the lessons library and output format — needed for ALL paths:
```
Read ${CLAUDE_SKILL_DIR}/reference/lessons.md
Read ${CLAUDE_SKILL_DIR}/reference/output-format.md
```

Assess review scale:

| Signal | Verdict |
|--------|---------|
| Small doc (<100 lines), single topic, no cross-references | **Solo review** — Lead reads + fixes directly, skip multi-role setup |
| Medium doc, some cross-references, known stable codebase | **Lightweight** — 2 roles (code-verifier + self-sufficiency-tester), minimal setup |
| Large doc (multi-chapter), multiple data sources, WIP state, prior review failures | **Full debate** — 3-4 roles, formal Phases 1-3 |

Present the triage result to the user. The user decides; Lead recommends but doesn't impose.

- **Solo**: Lead reviews alone using lessons + output format, produce findings, done.
- **Lightweight**: AskUserQuestion for L1 path ("Where is the source code for this document?"). L3 = the doc itself. Build a minimal Data Authority Table (L1 = user's answer, L3 = the doc). Then follow Phase 2 with 2 teammates (code-verifier + self-sufficiency-tester). Skip Phase 1's full ceremony but still inline the Data Authority Table + lessons + output format into teammate prompts.
- **Full debate**: continue to Phase 1.

---

## Phase 1: Setup — Data Authority & Roles

```
Read ${CLAUDE_SKILL_DIR}/operations/setup.md
```

Follow the setup operation to:
1. Identify the review target (already skimmed in Phase 0)
2. Work with the user to define the Data Authority Hierarchy (M1)
3. Assign review roles with explicit data source permissions (M2)
4. Produce a Review Plan for user approval

**The Review Plan must include:**
- Review target (files, scope, which sections)
- Data Authority Table (L1-Ln with concrete paths)
- Role assignments (3-4 roles, each with allowed/forbidden sources)
- Known context (what's the document's purpose, what's WIP, what's locked)
- Review focus (accuracy? completeness? self-sufficiency? scaling?)

**AskUserQuestion** to confirm the plan before proceeding.

---

## Phase 2: Team Debate — 开辩论会

**⛔ 这是辩论会，不是并行跑脚本。Teammates 必须能互相通信（SendMessage），互相质疑。**
**⛔ Standalone Agent（没有 team_name 的 Agent 调用）绝对禁止——它们是一次性黑盒，无法接收 SendMessage，无法辩论。**

### Step 2a: 建团

调用 TeamCreate：

```
TeamCreate(team_name: "arch-review")
```

等 TeamCreate 返回后才能继续。

### Step 2b: 生成 Teammates

对每个角色调用 Agent（一条消息内全部并行发出），**每个 Agent 调用必须带这些参数**：

```
Agent(
  team_name: "arch-review",   ← 加入 2a 建的团，能收 SendMessage
  name: "code-verifier",      ← 队友名字，SendMessage 用这个名字寻址
  model: "opus",              ← 强制 opus
  mode: "bypassPermissions",  ← 不需要用户确认
  prompt: "..."               ← 下面的 teammate prompt
)
```

**⛔ 检查清单（spawn 前必须确认）：**
- [ ] TeamCreate 已调用？
- [ ] 每个 Agent 调用都有 `team_name: "arch-review"`？
- [ ] 每个 Agent 调用都有 `name`？
- [ ] 没有用不带 team_name 的 standalone Agent？

**Teammate prompt** — MUST inline all context（teammates 没有讨论上下文）：

```
You are [ROLE NAME] on team "arch-review".

## Your Mission
[role mission from table below]

## Data Authority Hierarchy
[PASTE the full table from Review Plan]

## Golden Rule
Use L1/L2 to verify L3 facts. Use independent judgment for L3 design.
NEVER use lower-authority sources to override higher-authority sources.

## Your Data Access
- ALLOWED: [specific L-levels + paths]
- FORBIDDEN: [specific L-levels + paths + why]
You MUST NOT read any file under the FORBIDDEN paths.

## Accumulated Lessons
[PASTE full content of lessons.md]

## Output Format
[PASTE full content of output-format.md]

## Task
1. Read the document at [DOCUMENT PATH]
2. Produce your findings. Each finding MUST follow the output format.
3. SendMessage your findings to "team-lead" when done.
4. Wait — the Lead will share other teammates' findings for cross-examination.
5. When you receive findings, challenge any you disagree with.
   Cite your evidence. This is a debate, not a rubber stamp.
```

**Critical**: Inline ALL context into teammate prompts. Teammates cannot resolve `${CLAUDE_SKILL_DIR}` paths and have zero discussion context from Phase 0/1 — this isolation is their superpower against context pollution.

### Step 2c: Wave 0 — Independent Review

All teammates review in parallel. Lead does cross-chapter consistency personally (NOT delegated).

```
📡 WAVE 0 LAUNCH — Width: {N}
   [code-verifier]          → Technical fact verification (L1 vs L3)
   [completeness-checker]   → Gap analysis (L1+L2 vs L3)
   [self-sufficiency-tester] → Zero-context reader simulation (L3 only)
   [team-lead]              → Cross-chapter consistency check
```

### Step 2d: Wave 1 — Cross-Examination

**⛔ THIS STEP IS MANDATORY. Skipping cross-examination = no debate = skill provides zero value over plain code review.**

Once all teammates report findings to "team-lead" via SendMessage:
1. Lead uses `SendMessage` to send each teammate's findings to the OTHER teammates by name
2. Teammates challenge, confirm, or refine — citing their own evidence
3. Lead monitors: intervene via `SendMessage` if a teammate cites a forbidden data source

Example SendMessage flow:
```
SendMessage(to: "code-verifier", message: "self-sufficiency-tester found these issues: [paste]. Challenge any you disagree with.")
SendMessage(to: "self-sufficiency-tester", message: "code-verifier found these issues: [paste]. Challenge any you disagree with.")
```

When new challenges stop producing new insights (usually 1-2 rounds), Lead closes the debate.

### Step 2e: Close Debate

When new challenges stop producing new insights (usually 1-2 rounds), Lead closes the debate. Move to Phase 3. Keep team alive — teammates may be needed in Phase 4.

**Built-in roles** (use as needed, customize per review):

| Role | Mission | Default Sources |
|------|---------|-----------------|
| `code-verifier` | Verify document's technical facts against source code. Flag "castle-in-the-air" claims (written but no implementation evidence) | L1 (code) + L3 (docs) |
| `arch-challenger` | Challenge design decisions — find hidden assumptions | L1 + L3 + independent judgment |
| `completeness-checker` | Find gaps — what does the doc NOT mention that it should? | L1 + L2 (config) + L3 |
| `self-sufficiency-tester` | Read ONLY the document. Where would a zero-context reader go wrong? | **L3 ONLY** — forbidden from all other sources |
| `ops-simulator` | Follow the doc's procedures step by step, verify against code/config whether they actually work | L1 + L2 + L3 |
| `devil-advocate` | Predict production failures. What breaks at 10x scale? | L1 + L2 + L3 |

**Role count**: 3-4 per round. More = signal dilution. The self-sufficiency-tester is recommended for first review of any document.

---

## Phase 3: Debate & Synthesis

```
Read ${CLAUDE_SKILL_DIR}/operations/debate.md
```

Follow the debate operation to:
1. Collect all role findings
2. Cross-validate — do roles agree? Where do they conflict?
3. For conflicts: identify which role used the correct data source authority
4. Check for data-source-confusion errors (the #1 failure mode)
5. Classify findings by priority (P0/P1/P2/Doc)
6. Produce the Review Report

---

## Phase 4: Action

Present the Review Report to the user. Ask how to proceed:

| Option | Action |
|--------|--------|
| **Fix inline** | Apply confirmed fixes directly to the document(s) + produce changelog |
| **Report only** | Save report as standalone file (default: `arch-review-YYYY-MM-DD.md`) |
| **Selective** | User picks which findings to fix now vs defer |

For inline fixes: apply changes, then SendMessage the self-sufficiency-tester to re-check modified sections (team is still alive from Phase 2).

**Drive escalation**: Any participant (not just the Lead) can recommend using `/drive` or `/wing-drive` for complex fixes that require multi-file changes, code modifications, or cross-repo coordination. The review identifies WHAT to fix; drive executes HOW.

After Phase 4 is complete, shut down the team:
1. SendMessage each teammate: "Review complete. Please shut down."
2. Wait for teammates to go idle
3. `TeamDelete()` to clean up

---

## Anti-Patterns (from 3 projects, 7+ rounds of reviews)

These are the failure modes this skill is designed to prevent:

1. **Wrong ground truth** — Reviewer reads a WIP/legacy source and treats it as authoritative. The Data Authority Hierarchy (M1) prevents this.
2. **All roles read everything** — No differentiation, findings converge on same conclusions. Per-role data source restrictions (M2) prevent this.
3. **Lessons not embedded** — Same mistakes repeated across rounds. Lessons injection (M3) prevents this.
4. **Vague findings** — "This could be better" without evidence or action. Output format (M4) prevents this.
5. **No self-sufficiency check** — Document looks correct to someone who already knows the system but confuses newcomers. The self-sufficiency-tester role (M6) prevents this.

---

## Phase 5: Update Lessons

After Phase 4 is complete, check the Review Report's "New Lessons Learned" section.

If non-empty, append new lessons to:
```
${CLAUDE_SKILL_DIR}/reference/lessons.md
```

Format: `### [Project] [Date]: [Lesson title]` followed by what happened and how to prevent it.

**Note**: If running from an installed copy (`~/.claude/skills/`), lessons written there will be overwritten on next `install.sh`. Warn the user to also update the source at `skills/arch-review/reference/lessons.md`.

This is how the skill improves over time — lessons accumulate at the skill level, not in external memory.

---

## Revolution — Process Self-Optimization

Triggered by: `/arch-review revolution`

```
Read ${CLAUDE_SKILL_DIR}/operations/revolution.md
```

Auto-reads the latest debug log from `.claude/arch-review/debug/` (by timestamp).
If multiple logs exist, analyze all — cross-review patterns are more valuable than single-review analysis.

Follow the revolution operation to:
1. Read debug log(s)
2. Analyze 6 dimensions: role effectiveness, data source confusion, cross-examination value, lessons effectiveness, output quality, process efficiency
3. Produce concrete improvement proposals (with diffs to actual skill files)
4. Present to user for approval
5. Apply approved changes, bump version, commit

**The loop**: Review → Debug → Revolution → Better Review.
