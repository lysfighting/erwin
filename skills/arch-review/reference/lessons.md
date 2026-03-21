# Accumulated Review Lessons

> Hard-won from 3 projects, 7+ rounds, 14+ opus agents.
> Every reviewer prompt MUST include these lessons.
> Append new lessons after each review — this is how the skill improves.

---

## Data Source Confusion (the #1 failure mode)

### gcp-compute-infra Round 3 (2026-03-18): WIP config treated as ground truth
- **What happened**: 6 reviewers treated gitops-config (WIP K8s manifests) as the source of truth. They used L4-level data to reject L3-level design decisions, producing 80+ findings where the majority were in the wrong direction.
- **Root cause**: Reviewer prompts didn't specify data authority levels. Reviewers defaulted to "more config files = more truth."
- **Prevention**: Every reviewer prompt MUST include the Data Authority Hierarchy table. The Golden Rule must be stated explicitly: "NEVER use lower-authority sources to override higher-authority sources."

### vi_agent Phase B (2026-03-20): Legacy IaC treated as current state
- **What happened**: 3 Opus reviewers read `vi_agent/infra/` Terraform (legacy, last modified 2025-01) and treated it as the active IaC source of truth. They concluded Cloud SQL had `max_connections=100`, `POSTGRES_16`, and Memorystore `tier=BASIC` — all from the OLD definitions, not the target state managed by `gcp-compute-infra`.
- **Root cause**: The reviewers weren't told that IaC had been consolidated to a different repo. The legacy directory was still present and looked authoritative.
- **Prevention**: Legacy/deprecated data sources must be explicitly listed as FORBIDDEN in reviewer prompts, with explanation of WHY they're legacy.

### vi_agent Phase C (2026-03-20): Cascade from Phase B errors
- **What happened**: Phase C bootstrap review inherited 3 wrong assumptions from Phase B (Cloud SQL capacity, PG version, Memorystore tier). The error cascaded because Phase C treated Phase B's report as reliable input.
- **Root cause**: Review reports from prior rounds were treated as ground truth without re-verification against L1/L2 sources.
- **Prevention**: Prior review reports are AUXILIARY data (same as memory files). Always re-verify against L1/L2 before acting on prior findings.

---

## Document Self-Sufficiency

### vi_agent Round 5 (2026-03-20): Implicit knowledge not in docs
- **What happened**: The architecture docs were correct for someone who already knew the system. But a self-sufficiency test revealed that a zero-context reader would make wrong infrastructure decisions (e.g., looking at `vi_agent/infra/` instead of `gcp-compute-infra` for IaC).
- **Root cause**: The authors knew which repo owned IaC, but the document never explicitly stated it.
- **Prevention**: The self-sufficiency-tester role is the most valuable innovation. A reviewer who reads ONLY the document and NOTHING else will find exactly the gaps that cause real-world confusion.

---

## Role Design

### gcp-compute-infra Round 1-2: Role count and focus
- **What worked**: 3-4 focused roles per round, each with clear data source boundaries. Round 1 (accuracy) found 74 issues, 52 fixed. Round 2 (depth) found deeper design flaws.
- **What didn't work**: Round 3 expanded to 6 roles. Signal-to-noise ratio dropped. More roles ≠ better coverage.
- **Lesson**: 3-4 roles per round. Rotate role types between rounds (accuracy → depth → self-sufficiency) rather than adding more roles.

### gcp-compute-infra Round 2: ops-simulator value
- **What happened**: The ops-simulator role (simulate real operations following only the doc) found operational gaps that no other role caught — missing runbook steps, untested failover procedures, configuration that only works in theory.
- **Lesson**: ops-simulator is uniquely valuable for documents that describe operational procedures. It answers "can someone actually DO what this document says?"

---

## Output Quality

### All rounds: Vague findings are noise
- Findings like "this section could be improved" or "consider adding more detail" waste review cycles. The mandatory output format (Chapter + Evidence + Impact + Action + Priority) filters these out.
- **Rule**: If a reviewer can't cite specific evidence (file:line or exact text), the finding is not actionable and should be rejected.

### gcp-compute-infra Round 3: Effort estimates need verification
- **What happened**: The architecture doc said "Redis Streams migration: 1-2 days, ~200 LOC." Code audit found 7 subscriber files with different handling semantics. Realistic: 1-2 weeks, 500-800 LOC.
- **Lesson**: Any effort estimate in the document should be verified against actual code complexity. Under-estimated effort = under-allocated resources = late delivery.

---

## Process

### Cross-chapter consistency is the Lead's job
- **What happened**: In early rounds, cross-chapter consistency was delegated to a role. That role missed inconsistencies because they read chapters sequentially, not comparatively.
- **Prevention**: The Lead personally does cross-chapter consistency checks. This is NOT delegated.

### Debate output goes INTO the docs, not beside them
- **What happened**: Some review rounds produced detailed reports that sat as separate files. The actual documents were never updated. Same issues found again in next round.
- **Prevention**: Default action after review is "fix inline" — apply confirmed fixes directly to the document. A separate report is the fallback, not the default.
