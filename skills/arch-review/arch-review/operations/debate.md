# Debate Operation — Synthesis & Consensus

> Called from SKILL.md Phase 3, after Phase 2 Wave 1 (cross-examination) is complete.
> Teammates have already debated each other's findings. The Lead now does formal synthesis.

---

## Step 1: Collect Final Findings

Gather the **post-debate** findings from each role — use the revised versions after Wave 1 cross-examination. If a teammate withdrew or modified a finding during cross-examination, use the final state, not the original. If a teammate discovered new findings during Wave 1, include those too.

Assign sequential IDs per role (e.g., CV-01, CV-02, CC-01, SS-01).

For each finding, verify it includes required fields from the output format. Note: different formats have different required fields — the main format requires Priority but the self-sufficiency format intentionally omits it (Lead assigns priority in Step 5). Don't reject self-sufficiency findings for lacking Priority.

Findings missing fields for their format → mark as "needs supplement" and set aside.

**Tally**: N total findings across M roles (record the count — verify it at the end of Step 5).

---

## Step 2: Deduplicate

Multiple roles often find the same issue from different angles. Merge duplicates:
- Keep the version with the strongest evidence
- Note which roles independently identified it (cross-role confirmation = higher confidence)
- If roles identified different aspects of the same root cause, merge into one finding with combined evidence

---

## Step 3: Cross-Validate — Data Source Audit

**This is the most critical step.** For each finding, answer:

> "Did the reviewer use an appropriate-authority data source to reach this conclusion?"

**Green flag** (the primary valid pattern — do NOT over-protect L3):
0. Finding cites L1/L2 evidence showing L3 is factually wrong → **CONFIRM** — this is exactly what the review should find. Verify the L1/L2 evidence is current, then accept.

**Red flags** (from accumulated lessons):
1. Finding cites L4/L5 source to contradict L3 → **REJECT** (wrong ground truth)
2. Finding cites "I checked the config and it says X" without specifying WHICH config → **INVESTIGATE** (which config? L2 or L4?)
3. Finding about WIP state being "incomplete" → **LIKELY NOISE** (WIP is expected to be incomplete)
4. Finding from self-sufficiency-tester that contradicts code-verifier → **INVESTIGATE** (the doc may be correct but unclear, OR the doc may be wrong)
5. Finding says "should be X" but doesn't cite what it currently is → **NEEDS EVIDENCE**

**For each rejected/questioned finding:**
- Record WHY it was rejected
- Record which data source confusion occurred
- If this reveals a new lesson, queue it for lessons.md update

---

## Step 4: Conflict Resolution

When roles disagree:

### 4a: Identify the disagreement
```
Role A says: [claim + evidence]
Role B says: [counter-claim + evidence]
```

### 4b: Check data authority
- Which role used higher-authority data? That role's finding wins.
- If same authority level: check which has more specific evidence (file:line vs general observation).
- **Exception — orthogonal concerns**: Data authority comparison does NOT apply when one role judges factual accuracy and another judges document clarity. Example: code-verifier confirms "Section 3.2 is factually correct" while self-sufficiency-tester says "Section 3.2 would confuse a zero-context reader." Both findings can be valid simultaneously — they address different questions. Only apply authority ranking when two roles disagree on the SAME claim.

### 4c: If genuinely ambiguous
- Present both positions to the user
- AskUserQuestion: "Roles disagree on X. Role A says [Y] based on [source]. Role B says [Z] based on [source]. Which is correct, or do you need to investigate?"

### 4d: Self-Dialectic for design disagreements
For findings about design quality (not factual accuracy):
- **Thesis**: "The current design is correct because..."
- **Antithesis**: "The reviewer's concern is valid because..."
- **Synthesis**: "The truth that survives both is..."

---

## Step 5: Classify & Prioritize

Assign each confirmed finding a priority:

| Priority | Criteria | Examples |
|----------|----------|---------|
| **P0** | Blocks launch / causes data loss / security vulnerability | Empty secrets, wrong connection limits, missing auth |
| **P1** | Must fix before scaling / will cause incidents at load | Under-provisioned resources, missing indexes, effort underestimates |
| **P2** | Important but deferrable / won't cause immediate harm | Missing monitoring, incomplete Phase 2+ designs |
| **Doc** | Documentation-only — factual error or omission in text | Wrong line references, missing table entries, outdated names |

**Count check**: Confirmed + Rejected + Merged findings should equal the original tally from Step 1. If not, trace the missing IDs — findings were silently dropped.

---

## Step 6: Produce Review Report

Structure the report as follows:

```markdown
# Architecture Review Report — [DATE]

> **Review Team**: [roles] + Lead synthesis
> **Target**: [document(s) reviewed]
> **Data Authority**: [brief hierarchy summary]
> **Scope**: [what was reviewed, what was out of scope]

---

## Executive Summary
- Total findings: N
- P0: X | P1: Y | P2: Z | Doc: W
- Key insight: [one sentence — the most important thing learned]

---

## Findings by Chapter/Section

[Grouped by document structure. Each finding in standard output format.]

---

## Rejected Findings
[Findings that were rejected during cross-validation, with reasons.
This section is valuable — it shows WHERE reviewers went wrong and WHY.]

---

## Castle-in-the-Air List
[Claims in the document that have no implementation evidence.
"Written but not built." Each with: what's claimed, what actually exists, gap.]

---

## Self-Sufficiency Assessment
[From self-sufficiency-tester: where would a zero-context reader go wrong?
This is the most actionable section for document improvement.]

---

## Action Items
[Derived from confirmed findings. Each item: what to do, which file to change,
priority, and assignee/owner if known. This is the executive summary —
readers should be able to act on this section alone without reading all findings.]

---

## New Lessons Learned
[If this review revealed new failure patterns, document them here
for addition to the lessons library.]
```

---

## Step 7: Present to User

Show the report. Then proceed to SKILL.md Phase 4 (Action) for user decision on how to handle findings.
