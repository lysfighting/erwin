# Setup Operation — Review Initialization

> Called from SKILL.md Phase 1. Guides the Lead through defining data sources, roles, and scope.

---

## Step 1: Identify Review Target

Read the document(s) specified in `$ARGUMENTS`.
- Single file: read fully, note structure (chapters, sections)
- Directory: list files, read each to understand scope and relationships
- If ambiguous, AskUserQuestion to clarify

Produce a brief summary: what is this document about? What does it describe? What state is it in (draft, WIP, published)?

---

## Step 2: Define Data Authority Hierarchy

Work with the user to fill this table. The hierarchy is the most critical mechanism — it determines what counts as "truth" during review.

**Template** (adapt levels to the specific review):

| Level | Source | Authority | Concrete Path | Use |
|-------|--------|-----------|---------------|-----|
| L1 | Source code | Highest — system's actual behavior | `[user specifies repo/branch]` | Verify technical facts |
| L2 | Running config | High — current environment | `[docker-compose, .env, k8s manifests]` | Find omissions |
| L3 | **Review target** | Medium — the document being reviewed | `[from $ARGUMENTS]` | The object of review |
| L4 | [user-defined] | Low — WIP or auxiliary | `[gitops-config, migration scripts]` | Reference only |
| L5 | [user-defined] | Legacy — DO NOT use as ground truth | `[old IaC, deprecated docs]` | Historical context only |

**AskUserQuestion** with options:
1. "I'll define the hierarchy now" — user provides paths for each level
2. "Use code + config as L1/L2, the doc as L3" — minimal setup
3. "Let me explain the project context first" — freeform discussion
4. "Skip hierarchy, just review the document" — solo review without multi-source verification

If option 4: warn that without a hierarchy, wrong-ground-truth errors are likely. Proceed only if user confirms.

---

## Step 3: Identify Known Context

Before assigning roles, establish what the Lead (you) already knows:

1. **Document purpose**: What is the target state this document describes?
2. **Current state**: What exists today? (code, infra, config)
3. **WIP boundaries**: What is actively being built and therefore incomplete?
4. **Locked decisions**: What has already been decided and should not be re-debated?
5. **Prior reviews**: Have there been previous review rounds? What was found?

If the user provides this context, incorporate it into every role's prompt.
If not, the Lead should skim L1/L2 sources to establish this context before assigning roles.

---

## Step 4: Assign Roles

Select 3-4 roles from the built-in set or define custom ones. For each role:

```
### [Role Name]
- **Mission**: [one sentence — what are they looking for?]
- **Allowed sources**: [L-levels + specific paths]
- **Forbidden sources**: [L-levels + specific paths + why forbidden]
- **Focus areas**: [specific chapters/sections if scoped]
```

**Role selection heuristics:**

| Review goal | Recommended roles |
|-------------|-------------------|
| First review (accuracy) | code-verifier + completeness-checker + self-sufficiency-tester |
| Design quality | arch-challenger + devil-advocate + ops-simulator |
| Pre-launch readiness | code-verifier + ops-simulator + self-sufficiency-tester |
| Post-incident review | devil-advocate + ops-simulator + completeness-checker |

**Role assignment constraints:**
- Every role must have at least one forbidden source (prevents "read everything" convergence)
- self-sufficiency-tester MUST be forbidden from ALL sources except L3
- No more than 4 roles per round (signal dilution above 4)
- Lead does NOT delegate cross-chapter consistency — does it personally

---

## Step 5: Produce Review Plan

Assemble the complete Review Plan:

```
# Review Plan

## Target
[document path(s), scope, focus]

## Data Authority Hierarchy
[table from Step 2]

## Known Context
[from Step 3 — purpose, WIP boundaries, locked decisions]

## Role Assignments
[from Step 4 — each role with mission, sources, forbidden sources]

## Review Focus
[what matters most — accuracy? completeness? scalability? self-sufficiency?]

## Known Lessons
[summary of relevant lessons from reference/lessons.md]
```

Present to user for approval before starting Phase 2.
