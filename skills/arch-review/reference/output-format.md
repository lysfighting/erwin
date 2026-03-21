# Standard Output Format

> Every finding from every reviewer MUST follow this format.
> Findings not in this format are treated as noise and rejected.

---

## Required Format

```
**[Location]**: Problem description (one sentence)
**Evidence**: [L-level] specific citation — file:line, exact text, or measurable fact
**Impact**: What goes wrong if not fixed (concrete consequence, not vague risk)
**Action**: Specific edit — which file, which section, what to add/change/remove
**Priority**: P0 | P1 | P2 | Doc
```

---

## Priority Definitions

| Priority | Criteria | Decision rule |
|----------|----------|---------------|
| **P0** | Blocks launch, causes data loss, security vulnerability, or produces wrong behavior | "If we ship with this, something breaks" |
| **P1** | Must fix before scaling, will cause incidents under load, significant gap | "Works now but breaks at scale or under stress" |
| **P2** | Important but deferrable, no immediate harm, future improvement | "Should fix eventually but won't hurt today" |
| **Doc** | Documentation-only — factual error, omission, or inconsistency in text | "The system is fine but the doc is wrong/incomplete" |

---

## Location Format

Adapt to the document structure being reviewed:

- Architecture docs: `[Chapter X, Section Y]` or `[Ch3 §3.7]`
- API docs: `[Endpoint: POST /api/v2/users]`
- Runbooks: `[Procedure: Rollback, Step 4]`
- Config files: `[file.yaml:line]`
- General: `[Document name, heading]`

---

## Evidence Requirements

| Evidence type | Example | Strength |
|---------------|---------|----------|
| Code citation | `api-server/app/services/gcs_service.py:10 defaults to "vi-agent-uploads"` | Strong |
| Config citation | `docker-compose.yml:45 sets REDIS_URL=redis://...` | Strong |
| Cross-reference | `Ch3 says "25 connections" but Ch6 HPA max=3 implies 75 needed` | Strong |
| Measurement | `7 subscriber files found (grep -r "PSUBSCRIBE"), not 2 as documented` | Strong |
| Absence | `No K8s Job manifest found in gitops-config for Alembic migration` | Medium |
| Logical inference | `At 100 users, 10 tasks/day, $3/task = $3K/month, not $500 as claimed` | Medium |
| General observation | `This section seems incomplete` | **Rejected** — not specific enough |

---

## Rejected Finding Indicators

A finding should be rejected if:
1. No specific evidence cited (just opinion)
2. Evidence comes from a forbidden data source for that role
3. Evidence comes from a lower-authority source contradicting a higher-authority source
4. The finding is about WIP state being incomplete (expected behavior)
5. The finding duplicates a known issue already documented

When rejecting, record the reason — rejected findings reveal reviewer failure patterns.

---

## Castle-in-the-Air Format

For claims in the document that have no implementation evidence:

```
**[Location]**: "Quoted claim from document"
**Reality**: What actually exists (or doesn't)
**Gap**: What would need to be built/configured to make the claim true
**Priority**: P0 | P1 | P2
```

---

## Self-Sufficiency Finding Format

For the self-sufficiency-tester role specifically. **This role reads ONLY the document (L3). Do NOT consult code, config, or any other source.**

```
**[Location]**: What a zero-context reader would conclude from this text
**Confusion**: Why this conclusion could lead to wrong decisions
**Missing information**: What the document should state to prevent the wrong conclusion
**Action**: Suggested text to add (the Lead will verify correctness against L1/L2)
```

Note: The self-sufficiency-tester does NOT fill in "correct answers" — that's the Lead's job during synthesis. The tester's value is identifying WHERE confusion occurs, not resolving it.
