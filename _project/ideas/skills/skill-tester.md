# skill-tester — Erwin Skill Quality Checker

> Status: **Design** (from Session 2026-03-14-A, Round 20)
> A general-purpose testing skill for ALL erwin skills, not teamwork-specific.

---

## 1. Problem

Current teamwork has tw-unit-test.sh (825 lines) and tw-e2e-test.sh (675 lines), but:
- They test **scripts** (mechanism layer), not **skill behavior** (cognitive layer)
- VE compliance has **zero automated testing** — caught only by human eye or audit
- 10 versions fixed VE formatting bugs that kept regressing
- No way to verify "does Claude actually follow the instructions in SKILL.md?"

---

## 2. Three-Layer Testing Architecture

### Layer 1: Static Analysis (no LLM, runs in seconds)

Shell script that checks SKILL.md files against known quality rules:

```
Checks:
├── Size: SKILL.md ≤ 500 lines?
├── Frontmatter: has name, description fields?
├── VE reference: references ve-standard.md or has VE rules?
├── Issue numbers: any bare 2-digit #NN in non-code-block text?
├── Issue numbers: output templates use **#NNN** (bold + 3+ digits)?
├── Commands: all /skill-name references in backticks?
├── No ### headers in output templates (causes 2-newline gap)?
├── No code blocks wrapping output templates?
├── Supporting files: all files referenced in SKILL.md exist?
├── Cross-file consistency: version numbers match across skill suite?
├── Dividers: output uses ──── not --- or ═══?
└── File naming: SKILL.md is uppercase?
```

**Implementation**: A single `skill-lint.sh` script (~200 lines). Uses grep/awk, no LLM.

**Invocation**:
```bash
# Lint a single skill
./scripts/skill-lint.sh ~/.claude/skills/tw2/

# Lint all erwin skills
./scripts/skill-lint.sh --all
```

**Output**: Pass/fail per check with line numbers for violations.

### Layer 2: Smoke Test (needs LLM, ~30s per skill)

Uses `context:fork` to run a skill with known input, captures output, validates structure.

```
Flow:
1. Spawn subagent with context:fork
2. Inject skill + mock data (e.g., mock GitHub API responses)
3. Capture output text
4. Run regex patterns against output:
   - Issue numbers are **#NNN** format?
   - Dashboard has expected sections?
   - Doctor report has 3 layers?
   - Error messages use **ERROR:** prefix?
5. Report violations
```

**Challenge**: How to provide mock data? Options:
- A. Pre-recorded GitHub API responses in test fixtures
- B. Use a dedicated test repo with known state
- C. Inject structured data directly into the skill prompt

Option C is simplest — inject structured JSON as "pretend this is the API response" and check what the skill outputs.

### Layer 3: Behavioral Validation (needs LLM evaluator, ~60s per skill)

A second LLM call evaluates the skill output against VE rules and cognitive expectations.

```
Flow:
1. Run skill (Layer 2 output)
2. Send output + VE rules to evaluator LLM:
   "Here is the output of /tw2. Here are the VE rules.
    Rate compliance on each rule (PASS/FAIL/PARTIAL).
    List specific violations with line references."
3. Aggregate into report
```

**Risk**: "LLM testing LLM" circular dependency.
**Mitigation**: Tester prompt is <200 lines (vs 4,594 for team.md), so compliance rate is inherently higher. Also, the evaluator only needs to CHECK patterns, not GENERATE them — checking is easier than generating.

---

## 3. Skill-Tester as an Erwin Skill

### Structure

```
skills/skill-tester/
├── skill-tester/
│   └── SKILL.md              # ≤300 lines: orchestrates Layer 1-3
├── scripts/
│   └── skill-lint.sh          # Layer 1 static analysis
├── fixtures/
│   ├── mock-dashboard-data.json
│   ├── mock-doctor-data.json
│   └── ...
└── README.md
```

### SKILL.md Design (draft)

```yaml
---
name: skill-test
description: Test an erwin skill for quality, VE compliance, and behavioral correctness
argument-hint: <skill-name-or-path> [--layer 1|2|3|all]
disable-model-invocation: true
---
```

Content (~200 lines):
1. Parse arguments (skill name/path, layer selection)
2. Layer 1: Run `${CLAUDE_SKILL_DIR}/scripts/skill-lint.sh <path>`
3. Layer 2 (if requested): context:fork the target skill with mock data, validate output
4. Layer 3 (if requested): Evaluate output against VE rules
5. Generate report

### Invocation

```
/skill-test tw2                    # Layer 1 only (default, fast)
/skill-test tw2 --layer all        # All 3 layers
/skill-test tw2-drive --layer 2    # Smoke test only
/skill-test --all                  # Lint all erwin skills
```

---

## 4. Implementation Priority

1. **Layer 1 first** (skill-lint.sh) — no LLM needed, can validate current teamwork AND new tw2
2. **Layer 2 next** — after tw2 main skill exists (needs something to test)
3. **Layer 3 later** — nice to have, not blocking

### Layer 1 Estimated Effort

- skill-lint.sh: ~200 lines of shell
- SKILL.md for orchestration: ~150 lines
- Total: ~1 hour of implementation

### What Layer 1 Would Have Caught

Looking at teamwork's CHANGELOG, Layer 1 static analysis would have caught:
- v3.10.4: 14 VE fixes (bare #NN, code blocks around output)
- v3.9.14: 20+ hyperlink preservation fixes (bare #N/42 dilution)
- v3.9.7: 14 VE fixes (code block fences, missing bold)
- v3.8.3: Dashboard bold issue numbers

That's 4 audit sweeps (~50 fixes) that a 200-line shell script would have prevented.

---

## 5. Open Questions

- [ ] Should Layer 2 use `claude --print` (headless) or `context:fork`?
- [ ] How to mock GitHub API responses for smoke tests?
- [ ] Should tester be installed globally (~/.claude/skills/) or project-only?
- [ ] Can Layer 1 run as a pre-commit hook?

---

## 6. Source

- Session: `_project/sessions/2026-03-14-A.md` (Round 20)
- Teamwork CHANGELOG: 10 versions of VE fixes that would have been preventable
- Anthropic docs: context:fork, ${CLAUDE_SKILL_DIR}, supporting files
