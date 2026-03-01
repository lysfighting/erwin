# Idea Garden

> Ideas need to be captured immediately, nurtured over time, and harvested when ready.
> Not every seed becomes a fruit. That's fine.

---

## The Four Stages

```
Seed 🌱        →      Sprout 🌿      →      Plant 🌱       →    Fruit 🍎
Raw capture          Explored              Specified             Built/Released
(1 sentence)         (problem +            (ADR written,         (issue closed,
                      hypothesis)           ready to build)       changelog updated)
In INBOX.md          In skills/*.md        In skills/*.md        Archived
                      or meta/*.md          + ADR link            + release link
```

---

## Weekly Process

Every Monday (or start of the week):
1. Open INBOX.md
2. Promote Seeds that are worth exploring → move to skills/ or meta/
3. Discard Seeds that don't hold up on reflection
4. Review Sprouts: any ready to become Plants (write ADR)?
5. Note what changed in Plants/Fruits

Takes ~15 minutes. Prevents idea rot.

---

## Directories

| Directory | Purpose |
|-----------|---------|
| `INBOX.md` | Raw idea captures. Process weekly. |
| `skills/` | Drive-skill-like skill ideas |
| `meta/` | Universal patterns, cross-project ideas |

---

## Idea File Template (for Sprout+)

```markdown
# Idea: [Name]

- **Stage**: Sprout 🌿 | Plant 🌱 | Fruit 🍎
- **Created**: YYYY-MM-DD
- **ADR**: ADR-NNNN (when Plant)
- **Issue**: ISS-NNNN (when being built)

## Problem
[What problem does this solve? Be specific.]

## Hypothesis
[Why do I think this would work?]

## Key Questions
- [What needs to be figured out before building?]

## References / Prior Art
- [What exists that's similar?]

## Notes
[Ongoing exploration notes]
```
