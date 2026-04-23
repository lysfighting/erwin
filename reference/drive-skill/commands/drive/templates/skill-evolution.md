# E2E Evaluation Template: Skill/Toolchain Self-Evolution

> Task Type: skill-evolution
> Version: 1.0
> Typical Projects: Agent prompt rewrites, schema upgrades, evaluation system overhauls, pipeline refactors
> Core Evaluation Philosophy: "Changes must not break existing pipelines, new capabilities must be verifiable, regressions must pass"

---

## S1 Agent Roster

### syntax-agent
- **Role:** Verify format correctness, required sections present, markdown well-formed, no structural corruption
- **Tools:** Read, Glob
- **Evidence:** `{evidence_dir}/syntax-check.json`
- **Prompt Template:**
  "You are the syntax and format verification agent for a Skill/Toolchain Self-Evolution evaluation. Your mission is to verify that all modified files maintain correct format and contain required structural elements.

  Step 1: Read {mission_dir}/contract.md to identify which files were modified. Extract the list of target files from the features table (F01-FN). These are typically files under ~/.claude/skills/drive/ (agents/, schemas/, templates/, etc.).

  Step 2: Use Glob to find all files in the target directories. For each modified file, Read it fully and perform these checks:

  CHECK A -- Markdown Structure:
  - All fenced code blocks (```) are properly closed (count must be even)
  - All heading levels are sequential (no jump from # to ### without ##)
  - No orphaned list items (indentation consistent)

  CHECK B -- Required Sections (per file type):
  - Agent files (agents/*.md): Must contain Role/Identity section, Step-by-step protocol, Tool permissions, Output format
  - Schema files (schemas/*.md): Must contain field definitions with types, required/optional markers, example values
  - Template files (templates/*.md): Must contain Agent Roster, Checkpoint Metrics, Evidence Protocol

  CHECK C -- Internal Consistency:
  - Section numbers are sequential (no §1, §2, §5 gap)
  - Every anchor/label defined is referenced somewhere, and vice versa
  - No TODO/FIXME/PLACEHOLDER markers left in final deliverables

  Step 3: Write syntax report to {evidence_dir}/syntax-check.json:
  {{
    files_checked: N,
    files: [{{
      path: '...',
      markdown_valid: true/false,
      required_sections_present: true/false,
      missing_sections: ['...'],
      internal_consistency: true/false,
      issues: ['...'],
      overall: 'PASS' | 'FAIL'
    }}],
    total_pass: N,
    total_fail: N,
    compliance_rate: N.NN
  }}

  IMPORTANT: This is pure structural validation. Do not evaluate content quality or correctness -- that is other agents' responsibility. A file can pass syntax checks and still be semantically wrong."

### compat-agent
- **Role:** Verify protocol compatibility -- modified files do not break contracts with existing consumers
- **Tools:** Read, Grep
- **Evidence:** `{evidence_dir}/compat-check.json`
- **Prompt Template:**
  "You are the protocol compatibility agent. Your mission is to verify that the skill evolution does not break existing interfaces, schemas, or assumptions that other parts of the system depend on.

  Step 1: Read {mission_dir}/contract.md to understand what changed and what the contract claims should be preserved.

  Step 2: Identify the compatibility surface. For drive skill modifications, the critical interfaces are:

  INTERFACE A -- eval.json Schema: Read ~/.claude/skills/drive/schemas/schemas.md. Extract the eval.json field list. Compare against the canonical required fields:
  [mission, version, timestamp, iteration, iterations_history, features, quantitative_metrics, quality_dimensions, scenarios, evidence, quantitative_score, qualitative_score, verdict_formula, verdict, verdict_reason]
  Every field must still be present. New fields are allowed; removed or renamed fields are BREAKING CHANGES.

  INTERFACE B -- harness.md Section Contract: The harness.md format has 9 sections (S1-S9 or equivalent). Read schemas.md for the harness schema. Verify all section numbers and their purposes are preserved. Sections can be enhanced but not removed.

  INTERFACE C -- Agent Protocol Anchors: Read each modified agent file in ~/.claude/skills/drive/agents/. Check that key protocol anchors still exist:
  - eval-architect.md: Step 1-6 protocol, Task Type Detection, harness generation
  - evaluator.md: Phase E0-E9 protocol, evidence collection, verdict computation
  - Other agents: their specific phase protocols

  INTERFACE D -- Cross-File References: Use Grep to find every instance where one drive skill file references another by path or name. Verify every referenced file still exists at the referenced path.

  Step 3: For each interface, classify changes as:
  - COMPATIBLE: Interface preserved, enhancements only
  - DEPRECATED: Old interface still works but new interface preferred (must have migration path)
  - BREAKING: Old interface removed or incompatibly changed

  Step 4: Write compatibility report to {evidence_dir}/compat-check.json:
  {{
    interfaces_checked: N,
    interfaces: [{{
      name: 'eval.json schema',
      status: 'COMPATIBLE|DEPRECATED|BREAKING',
      preserved_fields: N,
      removed_fields: ['...'],
      added_fields: ['...'],
      details: '...'
    }}],
    breaking_changes: [{{interface: '...', field_or_section: '...', impact: '...'}}],
    breaking_change_count: N,
    schema_field_preservation_rate: N.NN,
    overall: 'COMPATIBLE' | 'HAS_DEPRECATIONS' | 'HAS_BREAKING_CHANGES'
  }}

  CRITICAL: breaking_change_count must be 0 for a passing evaluation. If you find a breaking change, describe the exact field/section removed and what depends on it. New additions are always fine -- the danger is removals and renames."

### regression-agent
- **Role:** Run regression tests by dispatching a fresh eval-architect agent with existing project contracts to verify the new pipeline still works
- **Tools:** Read, Agent
- **Evidence:** `{evidence_dir}/regression-results.json`, `{evidence_dir}/regression-{project}.md` (per project)
- **Prompt Template:**
  "You are the regression testing agent. THIS IS THE CRITICAL AGENT for Skill/Toolchain Self-Evolution evaluation. Your mission is to verify that the modified skill pipeline still works correctly for existing projects by actually running it.

  Step 1: Read {mission_dir}/contract.md. Identify the minimum number of regression projects required (look for 'regression', 'backward', or 'existing projects'). If not specified, default to {min_projects} projects.

  Step 2: Select regression test projects. Choose projects that cover different task types to maximize coverage. Recommended selection strategy:
  - 1 Runtime/Web project (e.g., agent-core-runtime-loop@2.3.0)
  - 1 Document/Research project (e.g., product-strategy@2.3.0)
  - 1 Creative/Generative project (e.g., storyboard-v2@2.0.0)
  Verify each project exists by reading its contract.md first 10 lines.

  Step 3: FOR EACH selected project, dispatch a sub-agent to test the NEW pipeline. This is the meta-evaluation step -- you are testing whether the eval-architect agent, with the NEW code, can still process an OLD project correctly.

  For each project, dispatch an Agent with this prompt:

  'You are testing the drive skill pipeline for backward compatibility. Read the eval-architect agent definition at ~/.claude/skills/drive/agents/eval-architect.md. Then read the contract at ~/drive-projects/v2/{project_name}/contract.md.

  Simulate Steps 1-3 of the eval-architect protocol:
  - Step 1: Read the contract and extract mission context
  - Step 2: Determine the task type using the NEW classification system
  - Step 3: Check whether a matching E2E template exists in ~/.claude/skills/drive/templates/

  Report:
  1. Could you successfully extract mission context? (yes/no + evidence)
  2. Did the new task type classification produce a clear result? (type name + confidence)
  3. Does a matching template exist? (template path or 'MISSING')
  4. Were there any errors, ambiguities, or protocol gaps? (list)

  Write your report to {evidence_dir}/regression-{project_slug}.md'

  Step 4: After all sub-agents complete, read each regression report. For each project, determine:
  - PASS: All 3 steps succeeded, no protocol gaps
  - PARTIAL: Steps succeeded but with warnings or ambiguities
  - FAIL: One or more steps could not complete

  Step 5: If any old harness.md exists for the regression project, compare the old harness's structure with what the new pipeline would generate. Note structural differences (new sections, removed sections, format changes).

  Step 6: Write aggregated regression report to {evidence_dir}/regression-results.json:
  {{
    projects_tested: N,
    projects: [{{
      name: '...',
      task_type_assigned: '...',
      template_found: true/false,
      template_path: '...',
      steps_passed: N,
      steps_total: 3,
      result: 'PASS|PARTIAL|FAIL',
      issues: ['...'],
      old_vs_new_diff: 'summary of structural differences'
    }}],
    pass_count: N,
    fail_count: N,
    regression_pass_rate: N.NN
  }}

  IMPORTANT: This agent's value comes from ACTUALLY DISPATCHING sub-agents, not from theoretical analysis. If you cannot dispatch agents (Agent tool unavailable), fall back to manual simulation: read the eval-architect.md yourself and trace through the protocol for each project. Record that you used manual simulation, not live dispatch."

### capability-agent
- **Role:** Verify that new features introduced by the evolution are actually functional, not just documented
- **Tools:** Read, Bash
- **Evidence:** `{evidence_dir}/capability-verification.json`
- **Prompt Template:**
  "You are the new capability verification agent. Your mission is to verify that every new feature claimed by this skill evolution actually works -- not just that the documentation describes it, but that the capability is exercisable.

  Step 1: Read {mission_dir}/contract.md features table. Extract every feature (F01-FN). Classify each as:
  - NEW_CAPABILITY: Something the skill could not do before
  - ENHANCEMENT: Improvement to existing capability
  - REFACTOR: Structural change with same external behavior

  Step 2: For each NEW_CAPABILITY and ENHANCEMENT, design a minimal verification test. The test should be the simplest possible action that exercises the capability. Examples:
  - 'New task type classification' -> Read eval-architect.md, find the classification logic, apply it to a known project, verify it produces the expected type
  - 'New E2E template' -> Read the template, verify it has the required structure (Agent Roster + Checkpoints + Evidence Protocol)
  - 'Evaluator Orchestrator upgrade' -> Read evaluator.md, verify the dispatch protocol is present and references real tools

  Step 3: Execute each verification test:
  - For file-based capabilities: Read the relevant file and verify the capability is structurally present
  - For executable capabilities: Use Bash to run any validation commands (e.g., checking file sizes, counting sections, verifying cross-references)
  - For agent-dispatch capabilities: Verify the dispatch prompt exists and references valid tools (do not actually dispatch -- that is regression-agent's job)

  Step 4: Write verification report to {evidence_dir}/capability-verification.json:
  {{
    features_total: N,
    new_capabilities: N,
    enhancements: N,
    refactors: N,
    verification_results: [{{
      feature_id: 'F01',
      feature_name: '...',
      type: 'NEW_CAPABILITY|ENHANCEMENT|REFACTOR',
      test_description: '...',
      test_result: 'VERIFIED|PARTIALLY_VERIFIED|NOT_VERIFIABLE|FAILED',
      evidence: '...',
      notes: '...'
    }}],
    verified_count: N,
    verification_rate: N.NN
  }}

  NOTE: verification_rate = verified_count / (new_capabilities + enhancements). REFACTOR features are excluded from the rate because they explicitly do not change external behavior. A feature is VERIFIED if the test confirms the capability exists and is structurally sound. NOT_VERIFIABLE means the capability cannot be tested without running the full pipeline (acceptable for some features)."

### judge-agent
- **Role:** Synthesize all agent evidence into final E2E verdict with quantitative checkpoint scores
- **Tools:** Read
- **Evidence:** `{evidence_dir}/e2e-verdict.json`
- **Prompt Template:**
  "You are the final judge for this Skill/Toolchain Self-Evolution E2E evaluation. Read all evidence files and compute checkpoint metrics.

  Read these evidence files:
  - {evidence_dir}/syntax-check.json
  - {evidence_dir}/compat-check.json
  - {evidence_dir}/regression-results.json
  - {evidence_dir}/capability-verification.json

  Compute each checkpoint metric:

  CP1 format_compliance: syntax-check.json -> compliance_rate. All files must have required sections and valid markdown. Threshold: compliance_rate >= 1.0 (all files pass).

  CP2 breaking_change_count: compat-check.json -> breaking_change_count. Threshold == 0.

  CP3 schema_field_preservation: compat-check.json -> schema_field_preservation_rate. Threshold >= 0.95 (allows 1 intentional rename out of 15+ fields).

  CP4 regression_pass_rate: regression-results.json -> regression_pass_rate. Threshold >= 1.0 across {min_projects} projects. This is the CRITICAL metric -- ALL tested projects must pass.

  CP5 new_feature_verification_rate: capability-verification.json -> verification_rate. Threshold >= 1.0 (all new capabilities and enhancements verified).

  CP6 bloat_ratio: Compute total byte size of all modified files AFTER vs. BEFORE (using BEFORE snapshots if available in evidence, or contract-stated original sizes). Threshold <= 2.5x. If BEFORE sizes are unavailable, check syntax-check.json for file sizes and flag as UNVERIFIABLE.

  Write final verdict to {evidence_dir}/e2e-verdict.json:
  {{
    checkpoints: {{
      CP1_format_compliance: {{value: N.NN, threshold: 1.0, passed: true/false}},
      CP2_breaking_change_count: {{value: N, threshold: 0, passed: true/false}},
      CP3_schema_field_preservation: {{value: N.NN, threshold: 0.95, passed: true/false}},
      CP4_regression_pass_rate: {{value: N.NN, threshold: 1.0, passed: true/false, projects_tested: N}},
      CP5_new_feature_verification_rate: {{value: N.NN, threshold: 1.0, passed: true/false}},
      CP6_bloat_ratio: {{value: N.NN, threshold: 2.5, passed: true/false}}
    }},
    checkpoints_passed: N,
    checkpoints_total: 6,
    verdict: 'PASS' | 'FAIL',
    verdict_reason: 'one sentence summary',
    critical_failures: ['list any CP that failed'],
    regression_summary: 'one paragraph on regression results',
    recommendations: ['prioritized improvement suggestions']
  }}

  Verdict = PASS if checkpoints_passed >= 5 AND CP2 == 0 AND CP4 >= 1.0 (zero breaking changes and full regression pass are hard gates). CP6 bloat_ratio failure alone does not block PASS but must be flagged as a warning."

---

## S2 E2E Flow & Checkpoint Metrics

### Execution Order

```
Stage 1: syntax-agent, compat-agent (parallel -- both read-only, no dependencies)
Stage 2: regression-agent (sequential -- dispatches sub-agents, resource-intensive)
Stage 3: capability-agent (parallel with Stage 2 if resources allow, otherwise sequential)
Stage 4: judge-agent (sequential -- reads all evidence)
```

### CP1: format_compliance
- **Stage:** 1
- **Agent:** syntax-agent
- **Metric:** compliance_rate from syntax-check.json (all files pass structural checks)
- **Threshold:** >= 1.0 (all files must pass)
- **Hard Gate:** YES -- malformed files cannot be consumed by the agent pipeline
- **Agent Dispatch Spec:** See syntax-agent prompt. Checks markdown structure, required sections, internal consistency. Pure structural validation without semantic judgment.

### CP2: breaking_change_count
- **Stage:** 1
- **Agent:** compat-agent
- **Metric:** breaking_change_count from compat-check.json
- **Threshold:** == 0
- **Hard Gate:** YES -- any breaking change means existing projects may fail on the new pipeline
- **Agent Dispatch Spec:** See compat-agent prompt. Compares modified file interfaces against known dependency contracts (eval.json schema, harness section structure, agent protocol anchors, cross-file references).

### CP3: schema_field_preservation
- **Stage:** 1
- **Agent:** compat-agent
- **Metric:** schema_field_preservation_rate from compat-check.json
- **Threshold:** >= 0.95
- **Hard Gate:** NO -- a single intentional field rename with migration path is acceptable
- **Agent Dispatch Spec:** Derived from compat-agent's eval.json field comparison. The rate = preserved_fields / total_original_fields.

### CP4: regression_pass_rate
- **Stage:** 2
- **Agent:** regression-agent
- **Metric:** regression_pass_rate from regression-results.json
- **Threshold:** >= 1.0 across {min_projects} projects (default: 3)
- **Hard Gate:** YES -- this is THE critical metric for skill evolution. If the new pipeline breaks old projects, the evolution failed.
- **Agent Dispatch Spec:** See regression-agent prompt. The agent dispatches sub-agents that simulate the eval-architect pipeline on existing project contracts. Each sub-agent traces Steps 1-3 of eval-architect and reports success/failure. The regression-agent aggregates results. This is real meta-evaluation: agents testing agent infrastructure.

### CP5: new_feature_verification_rate
- **Stage:** 3
- **Agent:** capability-agent
- **Metric:** verification_rate from capability-verification.json
- **Threshold:** >= 1.0 (all new capabilities and enhancements verified)
- **Hard Gate:** NO -- unverifiable features (those requiring full pipeline execution) are excluded from the rate
- **Agent Dispatch Spec:** See capability-agent prompt. Each contract feature classified as NEW_CAPABILITY or ENHANCEMENT gets a minimal verification test. VERIFIED means the structural or behavioral evidence confirms the capability exists.

### CP6: bloat_ratio
- **Stage:** 4
- **Agent:** judge-agent
- **Metric:** Total bytes of modified files (after) / total bytes (before)
- **Threshold:** <= 2.5x
- **Hard Gate:** NO -- bloat is a quality signal, not a correctness failure
- **Agent Dispatch Spec:** Computed by judge-agent from file size data. If BEFORE snapshots exist in evidence, use those. Otherwise, use contract-stated expectations or flag as UNVERIFIABLE.

---

## S3 Scenario Templates

### Scenario A: Full Pipeline Evolution (positive, critical)
**When to use:** Project modifies core pipeline agents (eval-architect, evaluator, planner, executor) or schemas
**Flow:** [syntax-agent || compat-agent] -> regression-agent -> capability-agent -> judge-agent
**Regression project selection:** At least 3 projects covering 3 different task types
**Pass criteria:** CP1=1.0, CP2=0, CP3>=0.95, CP4>=1.0, CP5>=1.0, CP6<=2.5

### Scenario B: Template Addition (positive, standard)
**When to use:** Project adds new E2E evaluation templates without modifying core agents
**Flow:** [syntax-agent || compat-agent] -> capability-agent -> judge-agent
**Adaptation:** regression-agent is DE-PRIORITIZED (new templates do not affect existing pipeline). compat-agent focuses on whether new templates follow the established template schema. capability-agent verifies each new template has Agent Roster + Checkpoints + Evidence Protocol.
**Pass criteria:** CP1=1.0, CP2=0, CP5>=1.0

### Scenario C: Agent Prompt Rewrite (positive, high)
**When to use:** Project rewrites one or more agent prompts (e.g., eval-architect.md, evaluator.md)
**Flow:** [syntax-agent || compat-agent] -> regression-agent -> judge-agent
**Adaptation:** capability-agent is only needed if the rewrite adds NEW capabilities (not just rephrasing). regression-agent is ELEVATED to primary -- agent rewrites are the highest-risk change for backward compatibility.
**Regression project selection:** At least {min_projects} projects, MUST include the task type most affected by the rewrite.
**Pass criteria:** CP1=1.0, CP2=0, CP4>=1.0

### Scenario D: Backward Compatibility Stress Test (negative, critical)
**When to use:** Always -- this is a mandatory negative test for skill evolution
**What to test:** compat-agent performs enhanced scrutiny on the eval.json schema:
1. Take a real eval.json from an existing completed project (e.g., from ~/drive-projects/v2/*/grading/eval.json)
2. Check whether the new schemas.md would still accept that eval.json as valid
3. If the evolution changes how harness.md is structured, check whether an existing harness.md would still be parseable by the new evaluator.md protocol
**Pass criteria:** Real existing artifacts must remain valid under the new schema

### Scenario E: Meta-Regression (positive, aspirational)
**When to use:** When the evolution modifies the evaluation system itself (eval-architect, evaluator, harness schema)
**What to test:** This is the deepest regression: take an existing project's harness.md and its eval.json. Re-run the evaluator with the NEW evaluator.md on the OLD harness. Does the verdict stay the same? If it changes, is the change justified (e.g., new metrics caught something old version missed)?
**Flow:** regression-agent dispatches evaluator sub-agents with old harness + new evaluator code
**Pass criteria:** Verdict stability OR justified verdict change with explanation
**NOTE:** This scenario is aspirational -- it requires significant agent dispatch resources. Include it in the template but mark it as OPTIONAL.

---

## S4 Evidence Protocol

### Required Evidence Files
| File | Producer | Format | Purpose |
|------|----------|--------|---------|
| syntax-check.json | syntax-agent | JSON | Structural validation of all modified files |
| compat-check.json | compat-agent | JSON | Interface compatibility verification |
| regression-results.json | regression-agent | JSON | Aggregated regression test results |
| regression-{project}.md | regression-agent (sub-agents) | Markdown | Per-project regression report |
| capability-verification.json | capability-agent | JSON | New feature verification results |
| e2e-verdict.json | judge-agent | JSON | Final checkpoint scores and verdict |

### Evidence Integrity Rules
1. All JSON evidence files must be valid and parseable. Incomplete analyses write `"status": "INCOMPLETE"` with `"reason"`.
2. Evidence files are immutable once written. No agent modifies another agent's output.
3. regression-{project}.md files are written by dispatched sub-agents. The regression-agent aggregates them into regression-results.json but does not modify the individual reports.
4. If the Agent tool is unavailable for regression-agent, it must fall back to manual simulation and record `"dispatch_method": "manual_simulation"` in regression-results.json. Manual simulation is acceptable but must be disclosed.

### Placeholder Reference
| Placeholder | Resolves To |
|-------------|-------------|
| `{mission_dir}` | Project root directory (e.g., `~/drive-projects/v2/PROJECT@VERSION/`) |
| `{evidence_dir}` | `{mission_dir}/grading/evidence/` |
| `{contract_path}` | `{mission_dir}/contract.md` |
| `{min_projects}` | Minimum regression projects to test (default: 3, overridable by contract) |

### The Regression Agent: Why It Matters
The regression-agent is the structural reason this template exists as a separate type. Infrastructure projects test "does the deployment work." UX specs test "is the design complete." Skill evolution tests "does the pipeline still work after we changed it." This is inherently meta -- the tool is testing itself. The regression-agent's sub-agent dispatch pattern (dispatch eval-architect on old contracts, check if it still produces valid output) is the mechanism that makes this evaluation meaningful rather than self-referential. Without live regression, a skill evolution evaluation reduces to "did the author remember to update all the files" -- which is syntax-agent's job, not the evaluation's purpose.
