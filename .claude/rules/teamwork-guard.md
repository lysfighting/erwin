---
# Teamwork skill modification guard
---

# Teamwork Skill Modification Guard

**When editing ANY file in `skills/teamwork/`:**

1. **MUST read first**: `docs/teamwork-design-constitution.md` — 10 design principles that protect core cognitive architecture. Each principle has a "violation signal" — if your change matches one, STOP and reconsider.

2. **MUST read first**: `docs/visual-encoding-standard.md` — VE rules that ALL teamwork output must follow. Cognitive instructions must NEVER include output formatting (Constitution Principle 2).

3. **After editing**: Check version consistency — all 6 skill files must have the same version in frontmatter. Use `tw-bump-version.sh` or manually update all files + CHANGELOG.

**Why this guard exists**: Teamwork skills contain load-bearing one-sentence cognitive instructions that look like comments but are core design decisions. Removing or weakening them without understanding the WHY (documented in the constitution) causes silent regressions in AI reasoning quality.
