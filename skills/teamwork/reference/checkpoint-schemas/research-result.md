# Checkpoint Schema: `research-result`

> Use when the mission is primarily "find out whether X exists / is safe / is wired up" before any build decision. Typical missions: "does a loadtest signup recipe already exist?", "is this API public or internal?", "has this baseline been computed before?"

Precedent: session 1 mission `#1302` (scheduled benchmark baseline — research found `#1283` had already wired baseline-delta support, mission was smaller than Contract assumed).

## Required Sections

### 1. Research question

Restate the single question the Contract is really asking. One sentence. If you cannot fit it in one sentence, the Contract is under-scoped — push back.

### 2. Findings with evidence

For each finding:
- **Claim**: one-line statement of what you discovered.
- **Evidence**: `file:line` citations, grep output, gh API responses, or `git log` lines. Team-lead must be able to reproduce each cite in the worktree.
- **Confidence**: high / medium / low + one-line reason (e.g., "high: read the source directly", "medium: inferred from commit history", "low: one engineer's comment").

### 3. Classification of the answer

One of:
- **Case A — exists cleanly**: what you expected exists and is usable as-is. Cite entry points.
- **Case B — doesn't exist at all**: no prior art. Contract's assumption wrong.
- **Case A-partial**: pieces exist but they need composition / adaptation. Describe the gap.

### 4. Implications for mission scope

Given the classification, is the mission:
- **Smaller than the Contract framed** (shift Success Criteria down)?
- **Larger** (needs more work than Contract described)?
- **Wrong shape entirely** (should be 2 missions or a different scope)?

Cite the exact Contract lines that change meaning given your findings.

### 5. One-command repro

If the finding matters for ops or execution, give a single shell command (or ≤3-line snippet) that reproduces the key evidence. The team-lead or a future reader should not have to replicate your research from scratch.

### 6. Open questions back to team-lead

What does the team-lead need to decide now given your findings? Usually:
- Confirm the scope shift (if any)
- Pick which Case A-partial piece is load-bearing
- Accept or contest your confidence calls on key findings

---

## Quality Bar

- Every claim is backed by a cite. No "I think" without evidence.
- Confidence is honest — prefer "medium" over fake-high.
- If the research changes the mission shape, you do NOT silently re-interpret the Contract; you surface the shift and wait for team-lead to own the decision.
- Zero write operations. Zero commits. This is pure read-mode checkpoint.
