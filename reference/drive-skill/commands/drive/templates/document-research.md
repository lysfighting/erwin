# E2E Evaluation Template: Document/Research

> Task Type: document-research
> Version: 1.0
> For use by: Eval Architect (harness.md generation) + Evaluator Orchestrator (execution)

**Scope:** Projects whose deliverable is a strategy document, research report, product specification, or analytical framework. The core evaluation verb is "read it deeply, verify its claims, check its internal logic."

**Representative Projects:** product-strategy@2.3.0, product-doc-final@2.3.0, ai-glasses-market-research@2.3.0, killer-uc-audit@2.3.0, killer-uc-deep-dive@2.3.0, use-case-pool@2.3.0

---

## §1 Agent Roster

### structure-agent
- **Role:** Structural completeness verification — confirms all required sections exist, hierarchy is correct, no empty sections, table of contents matches actual content.
- **Tools:** Read, Glob
- **Evidence:** `{evidence_dir}/structure-report.json`
- **Prompt Template:**
  "You are the structure agent for a document/research evaluation. Your job is to verify that the document has complete structural integrity — every required section exists, the hierarchy is sound, and there are no empty placeholders.

  Execute the following steps:

  1. Identify the deliverable document(s). Use Glob to find the main output files:
     - `{mission_dir}/**/*.md` — check for the primary deliverable document.
     - If {deliverable_paths} is provided, use those paths directly.
     Read each deliverable document in full.

  2. Extract the document structure. For each document:
     a. List all headings (# through ####) in order, preserving hierarchy.
     b. Count the total number of sections (each heading starts a section).
     c. Measure each section's content length (word count between this heading and the next).

  3. Check against the required section list from the contract. The required sections are: {required_sections}. For each required section:
     a. Find the corresponding heading (exact match or close synonym).
     b. Mark as PRESENT or MISSING.
     c. If present, check that the section is non-empty (word count > {min_section_words|50}).

  4. Check for structural problems:
     a. Empty sections: heading exists but content is fewer than {min_section_words|50} words.
     b. Orphan sections: headings that don't belong to any expected category (may indicate scope creep or disorganization).
     c. Duplicate sections: same topic covered under different headings.
     d. Hierarchy violations: h4 directly under h1 (skipping h2/h3).

  5. Write results to `{evidence_dir}/structure-report.json`:
  ```json
  {
    "documents_analyzed": ["path/to/doc.md"],
    "total_sections": 0,
    "required_sections": {
      "expected": ["Introduction", "Market Analysis", "Recommendations"],
      "present": ["Introduction", "Market Analysis"],
      "missing": ["Recommendations"],
      "coverage": 0.0
    },
    "empty_sections": [],
    "orphan_sections": [],
    "duplicate_sections": [],
    "hierarchy_violations": [],
    "total_word_count": 0,
    "section_word_counts": {"Introduction": 450, "Market Analysis": 1200},
    "structural_issues_count": 0
  }
  ```

  Be precise about section matching. 'Market Analysis' and 'Market Research' should be treated as matches if the content is equivalent. 'Conclusion' and 'Summary' are synonyms. Use semantic matching, not just string equality."

---

### depth-agent
- **Role:** Content depth and density analysis — measures information density, identifies thin sections, evaluates the ratio of concrete content (data, examples, specifics) versus vague generalities.
- **Tools:** Read
- **Evidence:** `{evidence_dir}/depth-analysis.json`
- **Prompt Template:**
  "You are the depth agent for a document/research evaluation. You assess whether the document has substantive content or is padded with generalities. This is one of the most important checks — a document can have perfect structure but zero insight.

  Read the full deliverable document(s) at {deliverable_paths}.

  For each major section (h2 level), perform the following analysis:

  1. Word count: Count the total words in this section (including subsections).

  2. Concrete content identification: A paragraph is 'concrete' if it contains at least one of:
     - Specific data point (number, percentage, dollar amount, date)
     - Named example (specific company, product, person, case study)
     - Direct quote or citation
     - Code snippet, formula, or structured specification
     - Specific actionable recommendation (not 'consider improving' but 'implement X by doing Y')
     Count concrete paragraphs vs total paragraphs.

  3. Hedge word detection: Count instances of hedge language per 1000 words:
     - English: 'might', 'could potentially', 'it is possible that', 'further research needed', 'arguably', 'it seems', 'to some extent', 'generally speaking'
     - Chinese: '可能', '或许', '需要进一步研究', '有待验证', '某种程度上', '大体上', '一般来说'
     These indicate vagueness. High hedge density = low confidence content.

  4. Information density scoring (per section, 1-10 scale):
     - 9-10: Every paragraph adds new information, rich in data and examples
     - 7-8: Mostly substantive, occasional connecting paragraphs without new info
     - 5-6: Mix of substance and filler, some sections feel padded
     - 3-4: More filler than substance, few concrete examples
     - 1-2: Almost entirely generic statements that could apply to any topic

  5. Write results to `{evidence_dir}/depth-analysis.json`:
  ```json
  {
    "total_word_count": 0,
    "section_analysis": [
      {
        "section": "Market Analysis",
        "word_count": 1200,
        "paragraph_count": 8,
        "concrete_paragraphs": 5,
        "concrete_ratio": 0.625,
        "hedge_count": 3,
        "hedge_density_per_1k": 2.5,
        "density_score": 7
      }
    ],
    "overall_concrete_ratio": 0.0,
    "overall_hedge_density_per_1k": 0.0,
    "avg_density_score": 0.0,
    "thinnest_sections": ["section with lowest density"],
    "densest_sections": ["section with highest density"],
    "content_density_verdict": "SUBSTANTIAL | ADEQUATE | THIN"
  }
  ```

  Be ruthless. A section titled 'Competitive Analysis' that says 'There are many competitors in this space, and the market is evolving rapidly' without naming a single competitor is a density score of 2. Call out specific examples of thin content in thinnest_sections."

---

### source-agent
- **Role:** Source verification and claim cross-validation — extracts key factual claims, searches the web to verify them, assesses source reliability. THIS IS THE KEY DIFFERENTIATOR for document evaluation.
- **Tools:** Read, WebSearch, WebFetch
- **Evidence:** `{evidence_dir}/source-verification.json`
- **Prompt Template:**
  "You are the source agent for a document/research evaluation. Your job is to verify that the document's factual claims are grounded in reality, not hallucinated or outdated. This is the most critical agent for document quality — a beautifully structured document with false claims is worse than useless.

  Read the full deliverable document(s) at {deliverable_paths}.

  Execute the following verification protocol:

  1. Claim extraction: Scan the entire document and extract verifiable factual claims. A verifiable claim is a statement that:
     - Cites a specific number, statistic, or data point (e.g., 'The market is worth $4.2B')
     - Names a specific company, product, or person doing something specific (e.g., 'Meta released Ray-Ban Meta glasses in October 2023')
     - Asserts a causal relationship or trend (e.g., 'Battery life is the #1 consumer complaint')
     - References a study, report, or source (e.g., 'According to Gartner...')
     Extract up to {max_claims_to_verify|15} claims, prioritizing: (a) claims that are central to the document's argument, (b) claims with specific numbers, (c) claims that cite named sources.

  2. For each extracted claim, perform web verification:
     a. Use WebSearch with a query designed to verify or refute the claim. Construct the query to find the original source, not just pages that repeat the claim.
     b. Read the top 2-3 results using WebFetch. Check whether the source confirms, contradicts, or is irrelevant to the claim.
     c. Assess source reliability: primary source (company press release, academic paper, government data) > secondary source (reputable news outlet) > tertiary (blog, forum, AI-generated summary).
     d. Classify the claim as: VERIFIED (found confirming primary/secondary source), PARTIALLY VERIFIED (found related info but not exact match), UNVERIFIED (no relevant source found), CONTRADICTED (found source that refutes the claim), OUTDATED (claim was true but for a different time period).

  3. Calculate verification metrics:
     - source_verification_rate = (VERIFIED + PARTIALLY VERIFIED) / total_claims_checked
     - contradiction_count = number of CONTRADICTED claims
     - outdated_count = number of OUTDATED claims

  4. Write results to `{evidence_dir}/source-verification.json`:
  ```json
  {
    "claims_extracted": 0,
    "claims_checked": 0,
    "verification_results": [
      {
        "claim": "The smart glasses market is projected to reach $15B by 2028",
        "claim_location": "Section 2, paragraph 3",
        "search_query": "smart glasses market size forecast 2028",
        "sources_found": [
          {"url": "https://...", "title": "...", "reliability": "secondary", "confirms": true}
        ],
        "classification": "VERIFIED",
        "notes": "Multiple sources confirm similar range ($13-17B)"
      }
    ],
    "verified_count": 0,
    "partially_verified_count": 0,
    "unverified_count": 0,
    "contradicted_count": 0,
    "outdated_count": 0,
    "source_verification_rate": 0.0,
    "highest_risk_claims": ["claims that are CONTRADICTED or UNVERIFIED and central to the argument"]
  }
  ```

  Critical rules:
  - Do NOT verify opinion statements ('X is the best approach'). Only verify factual claims.
  - Do NOT count the document itself as a source. The point is external verification.
  - When a claim is CONTRADICTED, provide the contradicting source and explain the discrepancy.
  - When a claim is OUTDATED, note the date it was valid and what the current state is.
  - Prioritize quality over quantity: 10 thoroughly verified claims > 30 superficially checked claims."

---

### consistency-agent
- **Role:** Internal consistency verification — detects contradictions within the document, checks that terminology is used uniformly, and verifies cross-references.
- **Tools:** Read
- **Evidence:** `{evidence_dir}/consistency-report.json`
- **Prompt Template:**
  "You are the consistency agent for a document/research evaluation. You check that the document does not contradict itself — a problem especially common in long documents written or revised over multiple sessions.

  Read the full deliverable document(s) at {deliverable_paths}.

  Check the following consistency dimensions:

  1. Factual consistency: Scan for cases where the same fact is stated differently in different sections:
     - Different numbers for the same metric (e.g., 'market size is $4B' in section 2, '$5.2B' in section 5)
     - Different timelines for the same event
     - Contradictory conclusions (e.g., section 3 says X is the biggest risk, section 7 says Y is)
     Record each contradiction with exact quotes and section locations.

  2. Terminology consistency: Identify key terms and concepts. Check that:
     - The same concept uses the same name throughout (not 'smart glasses' in some places and 'AI eyewear' in others, unless both are explicitly defined as synonyms)
     - Acronyms are defined on first use and used consistently after
     - Named entities are spelled consistently

  3. Cross-reference integrity: Find all internal references ('as discussed in Section X', 'see the framework above', 'per the earlier analysis'):
     - Verify the referenced section actually exists
     - Verify the reference accurately describes what the target section contains
     - Flag broken references (pointing to nonexistent sections or inaccurate descriptions)

  4. Recommendation consistency: If the document makes recommendations in multiple places, check that they are compatible:
     - No recommendation in one section contradicts a recommendation in another
     - Priority rankings are consistent (if section A says Feature X is #1 priority and section B ranks it #3, that is a contradiction)

  Write results to `{evidence_dir}/consistency-report.json`:
  ```json
  {
    "factual_contradictions": [
      {
        "topic": "market size",
        "location_a": "Section 2, paragraph 3",
        "statement_a": "market is $4B",
        "location_b": "Section 5, paragraph 1",
        "statement_b": "market is $5.2B",
        "severity": "major"
      }
    ],
    "terminology_inconsistencies": [
      {
        "concept": "the product",
        "variants": ["smart glasses", "AI eyewear", "the device"],
        "defined_as_synonyms": false
      }
    ],
    "broken_cross_references": [],
    "recommendation_conflicts": [],
    "contradiction_count": 0,
    "terminology_issue_count": 0,
    "broken_reference_count": 0,
    "overall_consistency": "CONSISTENT | MINOR_ISSUES | CONTRADICTORY"
  }
  ```

  Severity classification:
  - Major: Contradictory numbers, opposing conclusions, broken logic chains
  - Minor: Terminology drift, formatting inconsistencies, slightly different phrasings of the same idea
  - Cosmetic: Synonym variation without confusion risk

  Only count Major and Minor toward contradiction_count. Cosmetic issues are noted but do not count."

---

### judge-agent
- **Role:** Evidence aggregation and final verdict — reads all evidence files, scores each checkpoint, produces the evaluation judgment.
- **Tools:** Read
- **Evidence:** `{evidence_dir}/judge-verdict.json`
- **Prompt Template:**
  "You are the judge agent for a document/research evaluation. You do not analyze the document yourself. You read the evidence files produced by the other agents and render a structured judgment.

  Read the following evidence files:
  - `{evidence_dir}/structure-report.json` (from structure-agent)
  - `{evidence_dir}/depth-analysis.json` (from depth-agent)
  - `{evidence_dir}/source-verification.json` (from source-agent)
  - `{evidence_dir}/consistency-report.json` (from consistency-agent)

  For each checkpoint metric below, extract the relevant value and score it:

  CP1 section_coverage: structure-report.json -> required_sections.coverage == 1.0
  CP2 content_density: depth-analysis.json -> section_analysis where every section has word_count >= {min_section_words|50} (no empty sections) AND avg_density_score >= {min_density_score|5}
  CP3 concrete_ratio: depth-analysis.json -> overall_concrete_ratio >= {concrete_ratio_threshold|0.4}
  CP4 source_verification_rate: source-verification.json -> source_verification_rate >= {source_verification_threshold|0.6}
  CP5 contradiction_count: consistency-report.json -> contradiction_count == 0
  CP6 hedge_density: depth-analysis.json -> overall_hedge_density_per_1k <= {hedge_density_threshold|0.15} (per word, i.e., 150 hedge words per 1000 words max — but typically the threshold is much lower, like 15 per 1000 = 0.015; the eval-architect should set this per project)

  For each checkpoint, record: metric name, observed value, threshold, PASS/FAIL.

  Additionally, note qualitative signals:
  - thinnest_sections from depth-analysis.json (sections that drag overall quality down)
  - highest_risk_claims from source-verification.json (unverified claims central to the argument)
  - broken_cross_references from consistency-report.json (structural defects)

  Calculate overall: checkpoints_passed / checkpoints_total.

  Write verdict to `{evidence_dir}/judge-verdict.json`:
  ```json
  {
    "checkpoints": [
      {"id": "CP1", "name": "section_coverage", "observed": 1.0, "threshold": 1.0, "result": "PASS"},
      {"id": "CP2", "name": "content_density", "observed": 7.2, "threshold": 5, "result": "PASS"},
      {"id": "CP3", "name": "concrete_ratio", "observed": 0.52, "threshold": 0.4, "result": "PASS"},
      {"id": "CP4", "name": "source_verification_rate", "observed": 0.73, "threshold": 0.6, "result": "PASS"},
      {"id": "CP5", "name": "contradiction_count", "observed": 0, "threshold": 0, "result": "PASS"},
      {"id": "CP6", "name": "hedge_density", "observed": 0.008, "threshold": 0.015, "result": "PASS"}
    ],
    "checkpoints_passed": 0,
    "checkpoints_total": 6,
    "pass_rate": 0.0,
    "qualitative_notes": {
      "thinnest_sections": [],
      "highest_risk_claims": [],
      "broken_references": []
    },
    "evidence_files_read": [],
    "evidence_files_missing": [],
    "verdict": "PASS",
    "verdict_reason": "All 6 checkpoints passed. Complete section coverage. Content density 7.2/10 avg. 52% concrete paragraphs. 73% claims verified. Zero contradictions. Low hedge density."
  }
  ```

  Verdict rules:
  - PASS: pass_rate >= 0.8 AND CP1 (section_coverage) is PASS AND CP5 (contradiction_count) is PASS
  - CONDITIONAL PASS: pass_rate >= 0.6 AND CP1 is PASS
  - FAIL: pass_rate < 0.6 OR CP1 is FAIL OR contradiction_count > 3
  
  Special rule for source verification: If source_verification_rate < 0.4 (very low), add a WARNING even if the overall verdict is PASS, because the document may contain undetected false claims.
  
  If any evidence file is missing, note it in evidence_files_missing and score that checkpoint as FAIL with observed = 'MISSING_EVIDENCE'. Never infer or assume results."

---

## §2 E2E Flow & Checkpoint Metrics

### Execution Flow

```
Stage 1: structure-agent verifies document completeness
    ↓ (gate: at least 1 document found with sections)
Stage 2: depth-agent + source-agent + consistency-agent run in PARALLEL
    ↓ (all three complete)
Stage 3: judge-agent renders verdict
```

**Design note:** Three agents run in parallel at Stage 2 because they are fully independent — depth-agent analyzes content quality, source-agent verifies external claims, consistency-agent checks internal logic. None depends on the other's output.

### CP1: section_coverage
- **Stage:** 1
- **Agent:** structure-agent
- **Metrics:** section_coverage (float, threshold: == 1.0 — all required sections must exist)
- **Agent Dispatch Spec:**
  Agent: structure-agent
  Prompt: "Read {deliverable_paths}. Extract all headings. Compare to {required_sections}. Write structure-report.json with required_sections.coverage."
  Evidence: {evidence_dir}/structure-report.json
  Extract: required_sections.coverage from JSON
- **Gate:** If no documents found or coverage < 0.5, skip Stage 2. Judge-agent still runs.

### CP2: content_density
- **Stage:** 2
- **Agent:** depth-agent
- **Metrics:** avg_density_score (float 1-10, threshold: >= {min_density_score|5}), all sections word_count >= {min_section_words|50}
- **Agent Dispatch Spec:**
  Agent: depth-agent
  Prompt: "Read {deliverable_paths} fully. For each h2 section: count words, count concrete paragraphs, count hedge words, score density 1-10. Write depth-analysis.json."
  Evidence: {evidence_dir}/depth-analysis.json
  Extract: avg_density_score from JSON

### CP3: concrete_ratio
- **Stage:** 2
- **Agent:** depth-agent
- **Metrics:** overall_concrete_ratio (float, threshold: >= {concrete_ratio_threshold|0.4})
- **Agent Dispatch Spec:**
  Agent: depth-agent
  Prompt: "Count paragraphs with data/examples/specifics vs total paragraphs. Write overall_concrete_ratio to depth-analysis.json."
  Evidence: {evidence_dir}/depth-analysis.json
  Extract: overall_concrete_ratio from JSON

### CP4: source_verification_rate
- **Stage:** 2
- **Agent:** source-agent
- **Metrics:** source_verification_rate (float, threshold: >= {source_verification_threshold|0.6})
- **Agent Dispatch Spec:**
  Agent: source-agent
  Prompt: "Extract up to {max_claims_to_verify|15} verifiable factual claims. For each, WebSearch to find confirming or contradicting sources. Classify each as VERIFIED/PARTIALLY_VERIFIED/UNVERIFIED/CONTRADICTED/OUTDATED. Compute source_verification_rate. Write source-verification.json."
  Evidence: {evidence_dir}/source-verification.json
  Extract: source_verification_rate from JSON

### CP5: contradiction_count
- **Stage:** 2
- **Agent:** consistency-agent
- **Metrics:** contradiction_count (integer, threshold: == 0)
- **Agent Dispatch Spec:**
  Agent: consistency-agent
  Prompt: "Read {deliverable_paths} fully. Check for factual contradictions, terminology inconsistencies, broken cross-references, and conflicting recommendations. Count contradictions (major + minor). Write consistency-report.json."
  Evidence: {evidence_dir}/consistency-report.json
  Extract: contradiction_count from JSON

### CP6: hedge_density
- **Stage:** 2
- **Agent:** depth-agent
- **Metrics:** overall_hedge_density_per_1k (float, threshold: <= {hedge_density_threshold|0.015})
- **Agent Dispatch Spec:**
  Agent: depth-agent
  Prompt: "Count hedge words/phrases across the entire document. Divide by total word count * 1000 to get density per 1k words. Write to depth-analysis.json."
  Evidence: {evidence_dir}/depth-analysis.json
  Extract: overall_hedge_density_per_1k from JSON

---

## §3 Scenario Templates

### S1: Structural Completeness Audit [positive, critical]
**Agent Mission:**
  Agent: structure-agent
  Prompt: "Read the full deliverable and verify that every section specified in {required_sections} exists with non-trivial content (>= {min_section_words|50} words). This is the first gate: if the document is missing required sections, all other analysis is premised on incomplete data. Produce a complete section map with word counts. Flag any empty placeholder sections that have a heading but no content."
  Pass Criteria: section_coverage == 1.0 AND zero empty sections (all sections >= {min_section_words|50} words).

### S2: Depth and Substance Verification [positive, critical]
**Agent Mission:**
  Agent: depth-agent
  Prompt: "Read the full document and evaluate every major section for information density. The core question: would an expert in this field learn anything new from this section, or is it a restatement of common knowledge? Count concrete content (data, named examples, specific recommendations) vs vague generalizations. Compute concrete_ratio and hedge_density. Identify the thinnest section specifically and explain why it is thin."
  Pass Criteria: overall_concrete_ratio >= {concrete_ratio_threshold|0.4} AND overall_hedge_density_per_1k <= {hedge_density_threshold|0.015} AND avg_density_score >= {min_density_score|5}.

### S3: External Claim Verification [positive, critical]
**Agent Mission:**
  Agent: source-agent
  Prompt: "Extract the document's most important factual claims — statements that, if false, would undermine the document's conclusions. Use WebSearch to find original sources for each claim. This scenario is the key differentiator for document evaluation: it catches the 'confidently wrong' failure mode where a document reads well but is built on false premises. Verify at least {max_claims_to_verify|15} claims. For any CONTRADICTED claim, provide the contradicting source and explain the discrepancy in detail."
  Pass Criteria: source_verification_rate >= {source_verification_threshold|0.6} AND contradicted_count == 0 for claims central to the document's argument.

### S4: Internal Logic Stress Test [adversarial, high]
**Agent Mission:**
  Agent: consistency-agent
  Prompt: "This is an adversarial check. Read the document looking specifically for self-contradictions. Focus on: (1) numbers that appear in multiple sections — are they consistent? (2) Recommendations in the conclusion — do they follow from the analysis, or do they contradict findings? (3) Terminology — is the same concept called different things in different sections, creating confusion? (4) Cross-references — do 'as discussed above' references actually point to real content? This catches the 'patched together from multiple drafts' failure mode where each section is internally coherent but the whole document contradicts itself."
  Pass Criteria: contradiction_count == 0 AND broken_reference_count == 0.

### S5: Hedge Language and Actionability Probe [adversarial, medium]
**Agent Mission:**
  Agent: depth-agent
  Prompt: "This is a targeted probe for non-committal language. Scan the entire document for hedge words and phrases: 'might', 'could potentially', 'further research needed', 'it is possible', 'to some extent', '可能', '或许', '有待验证'. Count them. Then check the recommendations/conclusions section specifically: are the recommendations concrete and actionable ('Implement X by doing Y with timeline Z') or vague ('Consider exploring opportunities in this space')? A document that hedges in analysis is acceptable; a document that hedges in its conclusions fails to deliver."
  Pass Criteria: overall_hedge_density_per_1k <= {hedge_density_threshold|0.015} AND recommendations section density_score >= 7.

### S6: Fabrication Detection via Source Cross-Check [adversarial, high]
**Agent Mission:**
  Agent: source-agent
  Prompt: "This is a targeted adversarial probe for hallucinated content. Focus on claims that cite specific sources: 'According to Gartner...', 'A McKinsey report found...', 'Research by Stanford shows...'. Use WebSearch to verify that the cited source actually exists and actually says what the document claims it says. This catches the most dangerous failure mode in AI-generated documents: confident citation of sources that do not exist or say something different. Classify each cited source as: EXISTS_AND_CONFIRMS, EXISTS_BUT_DIFFERENT, DOES_NOT_EXIST."
  Pass Criteria: Zero DOES_NOT_EXIST sources. Zero EXISTS_BUT_DIFFERENT where the discrepancy is material (changes the argument's direction).

---

## §4 Evidence Protocol

| # | Agent | File | Format | Min Quality |
|---|-------|------|--------|-------------|
| 1 | structure-agent | structure-report.json | JSON with required_sections, section_word_counts, structural_issues_count | documents_analyzed length > 0 |
| 2 | depth-agent | depth-analysis.json | JSON with section_analysis array, overall_concrete_ratio, overall_hedge_density_per_1k, avg_density_score | At least 3 sections analyzed |
| 3 | source-agent | source-verification.json | JSON with verification_results array, source_verification_rate, contradicted_count | claims_checked >= 5 |
| 4 | consistency-agent | consistency-report.json | JSON with factual_contradictions, terminology_inconsistencies, broken_cross_references, contradiction_count | overall_consistency field present |
| 5 | judge-agent | judge-verdict.json | JSON with checkpoints array, pass_rate, verdict, verdict_reason | verdict is one of PASS/CONDITIONAL PASS/FAIL |

### Evidence Integrity Rules:
- Every JSON file must be valid JSON (parseable without error).
- source-verification.json must include actual WebSearch queries used and actual URLs found — not fabricated search results. If WebSearch is unavailable, source-agent must report this limitation and score source_verification_rate as 0.0 (not skip the metric).
- depth-analysis.json section_analysis must cover every h2 section in the document — no sections may be skipped.
- consistency-report.json must include exact quotes for every contradiction found, with section locations. Vague notes like 'some inconsistency in section 3' are invalid evidence.
- Missing evidence files are scored as FAIL for the corresponding checkpoint — never inferred.

---

## §5 Placeholder Reference

| Placeholder | Description | Example | Required |
|-------------|-------------|---------|----------|
| {mission_dir} | Project root directory | ~/drive-projects/v2/product-strategy@2.3.0 | Yes |
| {evidence_dir} | Evidence output directory | {mission_dir}/grading/evidence | Yes |
| {deliverable_paths} | Paths to the main deliverable documents | ["{mission_dir}/output/strategy.md"] | Yes |
| {required_sections} | JSON array of required section names | ["Introduction", "Market Analysis", "Competitors", "Recommendations", "Conclusion"] | Yes |
| {min_section_words} | Min word count per section to count as non-empty | 50 | No (default: 50) |
| {min_density_score} | Min avg content density score (1-10) | 5 | No (default: 5) |
| {concrete_ratio_threshold} | Min ratio of concrete paragraphs | 0.4 | No (default: 0.4) |
| {source_verification_threshold} | Min claim verification rate | 0.6 | No (default: 0.6) |
| {hedge_density_threshold} | Max hedge words per 1k total words | 0.015 | No (default: 0.015) |
| {max_claims_to_verify} | Max claims for source-agent to check | 15 | No (default: 15) |
