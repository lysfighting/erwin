# E2E Evaluation Template: UX/Design-Spec

> Task Type: ux-design-spec
> Version: 1.0
> Typical Projects: Interaction design specs, design constitutions, UX frameworks, experience blueprints
> Core Evaluation Philosophy: "Every interaction path defined, every state transition covered, every reference resolved"

---

## S1 Agent Roster

### completeness-agent
- **Role:** Verify every interaction path has a definition -- no user action should lead to an undefined system response
- **Tools:** Read, Glob
- **Evidence:** `{evidence_dir}/path-coverage.json`
- **Prompt Template:**
  "You are the interaction path completeness auditor for a UX/Design-Spec evaluation. Your mission is to verify that the spec defines a response for every plausible user action in every defined state.

  Step 1: Read {mission_dir}/contract.md to understand the product context and identify the core interaction surfaces (screens, modes, voice commands, gestures, hardware buttons, etc.).

  Step 2: Use Glob to find all deliverable files in {mission_dir}/. Read each spec document fully. Build an inventory of every interaction path the spec defines. An interaction path is: [User State] + [User Action] -> [System Response] + [Next State].

  Step 3: For each defined state in the spec, enumerate ALL plausible user actions (not just the happy path). Check whether the spec defines a system response for each. Categorize each action-state pair as:
  - DEFINED: spec explicitly describes the response
  - IMPLICIT: spec defines a general fallback or error handling that covers this case
  - UNDEFINED: no coverage -- user does this action and the spec says nothing about what happens

  Step 4: Pay special attention to edge cases that UX specs commonly miss:
  - First-time use / empty state (no data, no history)
  - Error recovery (network failure mid-action, timeout, permission denied)
  - Concurrent actions (user taps while animation is playing, notification arrives during modal)
  - Accessibility paths (screen reader, keyboard-only navigation, reduced motion)
  - Interruption flows (incoming call, app backgrounded, low battery)

  Step 5: Write coverage report to {evidence_dir}/path-coverage.json:
  {{
    states_found: N,
    interaction_paths_total: N,
    paths_defined: N,
    paths_implicit: N,
    paths_undefined: N,
    coverage_rate: N.NN,
    undefined_paths: [{{state: '...', action: '...', severity: 'CRITICAL|HIGH|MEDIUM'}}],
    edge_case_coverage: {{
      empty_state: 'DEFINED|UNDEFINED',
      error_recovery: 'DEFINED|UNDEFINED',
      concurrent_actions: 'DEFINED|UNDEFINED',
      accessibility: 'DEFINED|UNDEFINED',
      interruption: 'DEFINED|UNDEFINED'
    }}
  }}

  NOTE: coverage_rate = (paths_defined + paths_implicit) / interaction_paths_total. A rate of 1.0 is ideal but rare. Flag CRITICAL undefined paths (main user flows) separately from MEDIUM ones (rare edge cases)."

### state-agent
- **Role:** Analyze state machine completeness -- verify all transitions are defined, detect dead-end states and unreachable states
- **Tools:** Read
- **Evidence:** `{evidence_dir}/state-machine-analysis.json`
- **Prompt Template:**
  "You are the state machine analyst. This is the KEY differentiator of UX-Spec evaluation versus document evaluation. Your mission is to extract the implicit or explicit state machine from the UX spec and verify its structural integrity.

  Step 1: Read ALL deliverable files in {mission_dir}/. As you read, extract every state the system can be in. A state is any distinct mode, screen, configuration, or condition that affects how the system responds to user input. States may be explicit (named in the spec) or implicit (described but not named).

  Step 2: For each state, extract all transitions: [Current State] --(trigger)--> [Next State]. Triggers include user actions, system events, timeouts, and external signals.

  Step 3: Build a state transition matrix and perform these structural checks:

  CHECK A -- Dead-End Detection: A dead-end state is one with no outgoing transitions (the user gets stuck). For each state, verify at least one exit path exists. Exception: terminal states (e.g., 'app closed', 'session ended') are valid dead ends.

  CHECK B -- Unreachable State Detection: An unreachable state has no incoming transitions from any other state (cannot be entered). For each state, verify at least one entry path exists. Exception: the initial state (app launch, first use) has no predecessor.

  CHECK C -- Transition Completeness: For each state, are ALL plausible triggers accounted for? If State A handles triggers T1, T2, T3 but not T4 which is valid in that context, that is a gap.

  CHECK D -- Cycle Detection: Identify any infinite loops without exit (state A -> B -> C -> A with no break condition). These are potential UX traps.

  Step 4: Write state analysis to {evidence_dir}/state-machine-analysis.json:
  {{
    states_extracted: N,
    transitions_extracted: N,
    dead_end_states: [{{state: '...', is_valid_terminal: true/false, severity: 'CRITICAL|LOW'}}],
    dead_end_count: N,
    invalid_dead_ends: N,
    unreachable_states: [{{state: '...', is_valid_initial: true/false}}],
    unreachable_count: N,
    invalid_unreachable: N,
    incomplete_transitions: [{{state: '...', missing_trigger: '...', severity: '...'}}],
    cycles_without_exit: [{{cycle: ['A', 'B', 'C', 'A'], has_exit: true/false}}],
    structural_health: 'SOUND' | 'HAS_ISSUES' | 'SEVERELY_INCOMPLETE'
  }}

  IMPORTANT: This analysis requires careful reading. Do not skim. Every conditional branch ('if X then Y, else Z'), every mode switch, every screen transition must be captured. The value of this analysis is directly proportional to extraction thoroughness."

### reference-agent
- **Role:** Verify cross-reference consistency -- no broken references, no terminology drift, no orphaned definitions
- **Tools:** Read, Grep
- **Evidence:** `{evidence_dir}/reference-consistency.json`
- **Prompt Template:**
  "You are the cross-reference consistency checker. This is the second KEY differentiator of UX-Spec evaluation. UX specs are heavily cross-referenced -- a gesture defined in Section 2 is referenced in Sections 5, 8, and 12. Your mission is to verify every reference resolves and every term is used consistently.

  Step 1: Read ALL deliverable files in {mission_dir}/. As you read, build two indexes:

  INDEX A -- Definitions: Every term, concept, component, gesture, state, or interaction pattern that is DEFINED (given a name and description). Record where it is defined (file, section).

  INDEX B -- References: Every time a defined term is REFERENCED (mentioned but not re-defined). Record where it is referenced (file, section) and the exact name used.

  Step 2: Cross-reference validation:
  - For each reference in Index B, verify a matching definition exists in Index A
  - Broken reference = reference uses a name that has no definition anywhere
  - Near-miss = reference uses a slightly different name (e.g., 'touch gesture' vs 'tap gesture', 'home screen' vs 'main screen') -- these indicate terminology drift

  Step 3: Terminology consistency analysis. Use Grep across all deliverable files to find variant spellings or synonyms for key concepts. For each concept that appears to have multiple names:
  - Record all variants found
  - Identify which is the canonical name (most frequently used, or explicitly defined)
  - Count non-canonical usages

  Step 4: Orphan detection. Check Index A for definitions that are never referenced anywhere else in the spec. These may be vestigial (from a previous draft) or may indicate missing cross-references.

  Step 5: Write consistency report to {evidence_dir}/reference-consistency.json:
  {{
    definitions_found: N,
    references_found: N,
    broken_references: [{{reference: '...', location: '...', closest_match: '...'}}],
    broken_reference_count: N,
    terminology_variants: [{{concept: '...', variants: ['...', '...'], canonical: '...', non_canonical_count: N}}],
    terminology_consistency_rate: N.NN,
    orphaned_definitions: [{{term: '...', defined_in: '...'}}],
    orphan_count: N
  }}

  NOTE: terminology_consistency_rate = 1.0 - (total non-canonical usages / total references). A rate below 0.90 indicates significant drift that would confuse implementers."

### implementability-agent
- **Role:** Evaluate from an engineer's perspective -- which specs are precise enough to implement directly, and which are ambiguous
- **Tools:** Read
- **Evidence:** `{evidence_dir}/implementability-review.json`
- **Prompt Template:**
  "You are the implementability reviewer. You evaluate the UX spec from the perspective of an engineer who must build exactly what the spec describes. Your mission is to identify every point where the spec is too ambiguous for direct implementation.

  Step 1: Read ALL deliverable files in {mission_dir}/. For each interaction, transition, or behavior described in the spec, ask yourself: 'Can I implement this without making assumptions?'

  Step 2: Classify each spec element into one of three categories:
  - IMPLEMENTABLE: Precise enough to code directly. Timing values specified, exact behavior described, edge cases covered.
  - NEEDS_CLARIFICATION: Intent is clear but details are missing. For example: 'show a loading indicator' (what kind? where? after how many ms delay? what if the load completes before the indicator appears?)
  - AMBIGUOUS: Multiple valid interpretations exist and no way to determine which is intended. For example: 'respond naturally to the user' (what does 'naturally' mean in code?)

  Step 3: For each NEEDS_CLARIFICATION and AMBIGUOUS element, write a specific question that an engineer would ask the designer. Make the question precise -- not 'what does this mean?' but 'when the user says X and the system is in state Y, should the response be Z1 or Z2?'

  Step 4: Assess dimensional completeness -- does the spec cover these implementation-critical dimensions:
  - Timing: Are animation durations, transition delays, timeout values specified?
  - Layout: Are spatial relationships, sizes, positions defined (or deferred to a visual spec)?
  - Content: Is copy/text specified, or marked as placeholder?
  - Prioritization: When resources are constrained, what degrades first?
  - Platform differences: Are platform-specific behaviors noted where relevant?

  Step 5: Write implementability review to {evidence_dir}/implementability-review.json:
  {{
    total_spec_elements: N,
    implementable: N,
    needs_clarification: N,
    ambiguous: N,
    implementability_rate: N.NN,
    engineer_questions: [{{element: '...', category: 'NEEDS_CLARIFICATION|AMBIGUOUS', question: '...', impact: 'HIGH|MEDIUM|LOW'}}],
    dimensional_coverage: {{
      timing: 'COMPLETE|PARTIAL|MISSING',
      layout: 'COMPLETE|PARTIAL|MISSING|DEFERRED',
      content: 'COMPLETE|PARTIAL|PLACEHOLDER',
      prioritization: 'DEFINED|UNDEFINED',
      platform_differences: 'ADDRESSED|NOT_APPLICABLE|MISSING'
    }},
    ambiguous_spec_count: N
  }}

  NOTE: implementability_rate = implementable / total_spec_elements. A perfect rate is unrealistic for early-stage specs. The VALUE of this analysis is in the engineer_questions list -- these are actionable items the designer can address."

### judge-agent
- **Role:** Synthesize all agent evidence into final E2E verdict with quantitative checkpoint scores
- **Tools:** Read
- **Evidence:** `{evidence_dir}/e2e-verdict.json`
- **Prompt Template:**
  "You are the final judge for this UX/Design-Spec E2E evaluation. Read all evidence files and compute checkpoint metrics.

  Read these evidence files:
  - {evidence_dir}/path-coverage.json
  - {evidence_dir}/state-machine-analysis.json
  - {evidence_dir}/reference-consistency.json
  - {evidence_dir}/implementability-review.json

  Compute each checkpoint metric:

  CP1 interaction_path_coverage: path-coverage.json -> coverage_rate. Threshold >= 0.95.

  CP2 undefined_state_count: state-machine-analysis.json -> invalid_dead_ends + invalid_unreachable. Threshold == 0.

  CP3 dead_end_state_count: state-machine-analysis.json -> count dead_end_states where is_valid_terminal == false. Threshold == 0.

  CP4 broken_reference_count: reference-consistency.json -> broken_reference_count. Threshold == 0.

  CP5 terminology_consistency: reference-consistency.json -> terminology_consistency_rate. Threshold >= 0.95.

  CP6 ambiguous_spec_count: implementability-review.json -> ambiguous_spec_count. Threshold <= 3.

  Write final verdict to {evidence_dir}/e2e-verdict.json:
  {{
    checkpoints: {{
      CP1_interaction_path_coverage: {{value: N.NN, threshold: 0.95, passed: true/false}},
      CP2_undefined_state_count: {{value: N, threshold: 0, passed: true/false}},
      CP3_dead_end_state_count: {{value: N, threshold: 0, passed: true/false}},
      CP4_broken_reference_count: {{value: N, threshold: 0, passed: true/false}},
      CP5_terminology_consistency: {{value: N.NN, threshold: 0.95, passed: true/false}},
      CP6_ambiguous_spec_count: {{value: N, threshold: 3, passed: true/false}}
    }},
    checkpoints_passed: N,
    checkpoints_total: 6,
    verdict: 'PASS' | 'FAIL',
    verdict_reason: 'one sentence summary',
    critical_findings: ['list any CP that failed'],
    top_engineer_questions: ['top 5 most impactful questions from implementability review'],
    recommendations: ['prioritized list of spec improvements']
  }}

  Verdict = PASS if checkpoints_passed >= 5 AND CP2 == 0 AND CP4 == 0 (no undefined states and no broken references are hard gates)."

---

## S2 E2E Flow & Checkpoint Metrics

### Execution Order

```
Stage 1: completeness-agent, state-agent, reference-agent (parallel -- all read-only)
Stage 2: implementability-agent (sequential -- benefits from reading Stage 1 evidence for context)
Stage 3: judge-agent (sequential -- reads all evidence)
```

### CP1: interaction_path_coverage
- **Stage:** 1
- **Agent:** completeness-agent
- **Metric:** (paths_defined + paths_implicit) / interaction_paths_total
- **Threshold:** >= 0.95
- **Hard Gate:** NO -- near-complete coverage is acceptable if gaps are MEDIUM severity only
- **Agent Dispatch Spec:** See completeness-agent prompt. The agent reads all spec files, enumerates every state-action pair, and classifies coverage. The metric captures how much of the interaction surface has defined responses.

### CP2: undefined_state_count
- **Stage:** 1
- **Agent:** state-agent
- **Metric:** invalid_dead_ends + invalid_unreachable from state-machine-analysis.json
- **Threshold:** == 0
- **Hard Gate:** YES -- any undefined state means a user can reach a point where the spec provides no guidance
- **Agent Dispatch Spec:** See state-agent prompt. Structural analysis of the state machine extracted from the spec. Dead-ends that are valid terminals (session end, app close) are excluded.

### CP3: dead_end_state_count
- **Stage:** 1
- **Agent:** state-agent
- **Metric:** Count of dead_end_states where is_valid_terminal == false
- **Threshold:** == 0
- **Hard Gate:** YES -- dead-end states are UX traps
- **Agent Dispatch Spec:** Derived from state-agent's state-machine-analysis.json. A dead-end is a state with no outgoing transitions that is not explicitly marked as a terminal state.

### CP4: broken_reference_count
- **Stage:** 1
- **Agent:** reference-agent
- **Metric:** broken_reference_count from reference-consistency.json
- **Threshold:** == 0
- **Hard Gate:** YES -- broken references mean implementers will hit "undefined" when they follow the spec
- **Agent Dispatch Spec:** See reference-agent prompt. Cross-reference Index A (definitions) against Index B (references). Any reference without a matching definition is broken.

### CP5: terminology_consistency
- **Stage:** 1
- **Agent:** reference-agent
- **Metric:** terminology_consistency_rate from reference-consistency.json
- **Threshold:** >= 0.95
- **Hard Gate:** NO -- minor terminology drift is a quality issue, not a structural failure
- **Agent Dispatch Spec:** Derived from reference-agent's terminology variant analysis. The rate measures what fraction of all term usages use the canonical form.

### CP6: ambiguous_spec_count
- **Stage:** 2
- **Agent:** implementability-agent
- **Metric:** ambiguous_spec_count from implementability-review.json
- **Threshold:** <= 3
- **Hard Gate:** NO -- some ambiguity is expected in early-stage specs
- **Agent Dispatch Spec:** See implementability-agent prompt. Only elements classified as AMBIGUOUS (not NEEDS_CLARIFICATION) count toward this metric. Ambiguous means multiple valid interpretations exist.

---

## S3 Scenario Templates

### Scenario A: Complete UX Spec Audit (positive, critical)
**When to use:** Project delivers a comprehensive UX spec intended as the engineering blueprint
**Flow:** [completeness-agent || state-agent || reference-agent] -> implementability-agent -> judge-agent
**Pass criteria:** CP1>=0.95, CP2==0, CP3==0, CP4==0, CP5>=0.95, CP6<=3

### Scenario B: Spec Revision/Simplification (positive, standard)
**When to use:** Project revises or simplifies an existing spec (e.g., 2638 lines to 400-600 lines)
**Flow:** [completeness-agent || state-agent || reference-agent] -> implementability-agent -> judge-agent
**Adaptation:** completeness-agent additionally checks that no interaction paths from the original spec were silently dropped (reads both old and new versions). Add a supplementary metric: dropped_path_count (paths in original but not in revision, without explicit deprecation note).
**Pass criteria:** Same as Scenario A + dropped_path_count <= contract-specified tolerance

### Scenario C: Design Constitution/Framework (positive, standard)
**When to use:** Project delivers high-level design principles rather than detailed interaction specs
**Flow:** [reference-agent] -> implementability-agent -> judge-agent
**Adaptation:** completeness-agent and state-agent are DE-PRIORITIZED (design constitutions define principles, not state machines). reference-agent becomes primary -- checking that principles do not contradict each other and that terminology is consistent across the framework. implementability-agent checks whether principles are specific enough to guide actual design decisions (not just platitudes).
**Pass criteria:** CP4==0, CP5>=0.95, CP6<=5 (higher ambiguity tolerance for principles)

### Scenario D: Terminology Drift Detection (negative, high)
**When to use:** Always -- this is a targeted negative test for the most common UX spec failure mode
**What to test:** reference-agent specifically hunts for synonyms: 'tap' vs 'touch' vs 'press', 'screen' vs 'page' vs 'view', 'notification' vs 'alert' vs 'toast'. These are the exact kind of inconsistencies that cause different engineers to implement the same concept differently.
**Pass criteria:** CP5 >= 0.95 with explicit listing of all variant pairs found

---

## S4 Evidence Protocol

### Required Evidence Files
| File | Producer | Format | Purpose |
|------|----------|--------|---------|
| path-coverage.json | completeness-agent | JSON | Interaction path inventory and coverage rate |
| state-machine-analysis.json | state-agent | JSON | State machine structural integrity |
| reference-consistency.json | reference-agent | JSON | Cross-reference validation and terminology analysis |
| implementability-review.json | implementability-agent | JSON | Engineering perspective on spec precision |
| e2e-verdict.json | judge-agent | JSON | Final checkpoint scores and verdict |

### Evidence Integrity Rules
1. All JSON files must be valid and parseable. If an agent cannot complete analysis, write `"status": "INCOMPLETE"` with `"reason"` field.
2. Evidence files are immutable once written. Agents do not modify each other's output.
3. The judge-agent is a pure aggregator -- it reads evidence and computes scores, never re-analyzes source documents.
4. For Scenario B (revision), completeness-agent writes an additional file: `{evidence_dir}/revision-delta.json` comparing old vs. new spec coverage.

### Placeholder Reference
| Placeholder | Resolves To |
|-------------|-------------|
| `{mission_dir}` | Project root directory (e.g., `~/drive-projects/v2/PROJECT@VERSION/`) |
| `{evidence_dir}` | `{mission_dir}/grading/evidence/` |
| `{contract_path}` | `{mission_dir}/contract.md` |

### Key Differentiators from Document/Research Evaluation
This template differs from the Document/Research template in three structural ways:
1. **State machine analysis** -- documents do not have state machines; UX specs do. The state-agent is unique to this template.
2. **Cross-reference density** -- documents have citations (outward references); UX specs have dense internal cross-references (gesture defined in S2 referenced in S5, S8, S12). The reference-agent's internal consistency check is far more rigorous here.
3. **Implementability** -- documents are evaluated for "actionability" (can someone do something with this?); UX specs are evaluated for "implementability" (can an engineer write code from this without guessing?). The engineer_questions list is the unique artifact.
