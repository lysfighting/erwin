# Revolution — Process Self-Optimization

> Analyzes debug logs from past reviews to find process-level improvements.
> Output: concrete diffs to lessons.md, SKILL.md, and operations files.
> The user approves changes before they're applied.

---

## Input

Auto-read from `.claude/arch-review/debug/` — list files, sort by name (timestamp), read the latest.
If multiple logs exist, read all — cross-review patterns reveal systemic issues that single-review analysis misses.

If no debug logs found, inform the user: "No debug logs found. Run `/arch-review --debug` first."

---

## Analysis Dimensions

Work through each dimension. For each, produce a verdict with evidence from the debug log.

### D1: Role Effectiveness

| Question | How to answer |
|----------|---------------|
| Which role produced the most P0/P1 findings? | Count confirmed findings by role |
| Which role's findings were mostly rejected? | Count rejected findings by role — high reject rate = misconfigured role |
| Did any role produce zero actionable findings? | If yes, was the role necessary? Or was it given wrong sources? |
| Did the self-sufficiency-tester peek at non-L3 sources? | Check [DEBUG] Findings metadata: "Forbidden sources cited?" |

### D2: Data Source Confusion

| Question | How to answer |
|----------|---------------|
| Did any teammate use a forbidden data source? | Check each teammate's "Data sources cited" vs "Forbidden sources" |
| Did the Lead catch it during cross-examination? | Check [DEBUG] Cross-Examination for Lead interventions |
| Did any findings get confirmed that shouldn't have been? | Check if synthesis accepted findings based on wrong-authority data |

### D3: Cross-Examination Value

| Question | How to answer |
|----------|---------------|
| Did cross-examination change any findings? | Compare pre-debate vs post-debate findings |
| How many rounds of cross-examination occurred? | Count from debug log |
| Were challenges substantive or rubber-stamps? | Check if challenges cited new evidence or just agreed |
| Did cross-examination surface findings that no single role found? | Check for new findings in round 2+ |

### D4: Lessons Effectiveness

| Question | How to answer |
|----------|---------------|
| Were injected lessons relevant to this review? | For each lesson, did a related situation arise? |
| Were any lessons ignored despite being relevant? | Did a teammate repeat a mistake that a lesson warned about? |
| Should new lessons be added based on this review? | Check [DEBUG] Outcome: "New lessons identified" |
| Should any existing lessons be retired? | Lessons that haven't been relevant in 3+ reviews |

### D5: Output Quality

| Question | How to answer |
|----------|---------------|
| Did all findings follow the output format? | Check for "needs supplement" findings in synthesis |
| Were priorities correctly assigned? | P0s should be genuine launch blockers, not style issues |
| Was the review report actionable? | Did the user choose "fix inline" and succeed, or get stuck? |

### D6: Process Efficiency

| Question | How to answer |
|----------|---------------|
| Was the triage correct? | Would a different scale (Solo/Lightweight/Full) have been better? |
| Were too many/few roles assigned? | Signal-to-noise ratio from role effectiveness analysis |
| How long did the review take relative to document size? | Rough efficiency metric |

---

## Output: Improvement Proposals

For each dimension that reveals an issue, produce a concrete proposal:

```markdown
### Proposal [N]: [title]

**Evidence**: [cite specific debug log section]
**Problem**: [what went wrong or could be better]
**Proposed change**:
- File: [which file to modify]
- Current: [what it says now]
- Proposed: [what it should say]
**Risk**: [what could go wrong with this change]
```

### Change Categories

| Category | Target file | Example |
|----------|-------------|---------|
| New lesson | `reference/lessons.md` | "Teammate X cited L4 despite being told not to — add stronger warning" |
| Role adjustment | `SKILL.md` role table | "ops-simulator produced zero findings on API docs — mark as not recommended for API reviews" |
| Process change | `operations/debate.md` | "Cross-examination round 2 added no value — consider closing after 1 round for small reviews" |
| Output format tweak | `reference/output-format.md` | "Self-sufficiency findings needed priority after all — add Priority as optional field" |
| Prompt improvement | `SKILL.md` teammate prompt | "Teammate ignored forbidden sources — add explicit 'I confirm I will NOT read...' commitment to prompt" |

---

## Applying Changes

Present all proposals to the user. User approves, modifies, or rejects each one.

For approved changes:
1. Apply the diffs to the actual skill files
2. Bump the patch version in SKILL.md frontmatter
3. Commit with: `fix(skills): arch-review revolution — [summary of changes]`

This is the self-optimization loop: Review → Debug → Revolution → Better Review.
