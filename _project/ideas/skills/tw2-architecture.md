# Teamwork v2 Architecture — Design Document

> Status: **Design** (from Session 2026-03-14-A, Rounds 1-20)
> Working codename: **tw2** (final name TBD — must reflect "AI self-iteration protocol" not just "team collaboration")

---

## 1. Why Redesign

### Root Cause (backed by data)

| Problem | Evidence |
|---------|----------|
| team.md = 4,594 lines | 9.2x over Anthropic's ≤500 line recommendation |
| VE rules 2,369 lines from output templates | Signal decay → LLM ignores its own correct rules |
| 10 versions fixing VE compliance | Fix pattern escalating, not converging |
| v3.9.14 "dilution" finding | Non-output `#N` (~40x) diluted correct `**#NNN**` (~50x) |
| 0 official Anthropic patterns used | No supporting files, no context:fork, no ${CLAUDE_SKILL_DIR} |

### What's NOT Wrong (preserve these)

- GitHub-native protocol (Issues, PRs, CI — no shadow state)
- Mission Contract (AI-enriched work spec)
- claim → drive → ship pipeline (AI-executable dev loop)
- Von Neumann separation (scripts = mechanisms, skill.md = content)
- Design Constitution (15 principles)
- 11 shell scripts (4,473 lines of tested deterministic ops)

---

## 2. Vision

Teamwork's core value is NOT "team collaboration tool" — it's an **AI autonomous development lifecycle protocol**.

The claim→drive→ship pipeline is a loop that AI can execute independently:
1. AI detects bug → creates Issue
2. AI claims Issue → generates Contract
3. AI drives → writes code, runs tests
4. AI ships → creates PR, waits for CI
5. AI reviews → self-review or requests human review
6. Loop

The redesign should optimize for this autonomous loop's reliability, not for human collaboration UX.

---

## 3. Target Architecture

### Directory Structure (official Anthropic skill pattern)

```
skills/tw2/                           # erwin source directory
│
├── tw2/                              # /tw2 — main skill (dashboard + router)
│   ├── SKILL.md                      # ≤500 lines: identity + routing + dashboard
│   ├── operations/
│   │   ├── doctor.md                 # 12 diagnostic checks (context:fork candidate)
│   │   ├── init.md                   # Project init flow (context:fork candidate)
│   │   └── auto.md                   # Full auto pipeline (~100 lines thin orchestrator)
│   └── reference/
│       ├── ve-standard.md            # VE rules — SINGLE SOURCE OF TRUTH
│       └── label-states.md           # State machine reference
│
├── tw2-issue/                        # /tw2-issue
│   ├── SKILL.md                      # ≤500 lines
│   └── operations/
│       └── batch.md                  # Batch issue creation (context:fork candidate)
│
├── tw2-claim/                        # /tw2-claim
│   └── SKILL.md                      # ≤500 lines (current team-claim is 453, nearly fits)
│
├── tw2-drive/                        # /tw2-drive
│   ├── SKILL.md                      # ≤500 lines
│   └── operations/
│       └── team-execution.md         # Multi-agent parallel execution details
│
├── tw2-ship/                         # /tw2-ship
│   ├── SKILL.md                      # ≤500 lines
│   └── operations/
│       └── review.md                 # AI code review (context:fork candidate)
│
├── tw2-rc/                           # /tw2-rc
│   ├── SKILL.md                      # ≤500 lines
│   └── operations/
│       ├── promote.md                # RC → production release
│       └── hotfix.md                 # Hotfix flow
│
├── scripts/                          # Shared mechanism layer (Von Neumann)
│   ├── tw2-config.sh                 # Config reader
│   ├── tw2-git.sh                    # Git operations
│   ├── tw2-pr.sh                     # PR lifecycle
│   ├── tw2-label.sh                  # Label state machine
│   ├── tw2-contract.sh               # Contract lifecycle
│   ├── tw2-notify.sh                 # Notifications
│   └── ...
│
├── README.md
├── CHANGELOG.md
└── ONBOARDING.md
```

### Key Architectural Decisions

1. **Install target**: `~/.claude/skills/tw2/SKILL.md` (NOT `~/.claude/commands/tw2.md`)
2. **Supporting files**: Heavy operations in `operations/`, loaded on demand via `${CLAUDE_SKILL_DIR}`
3. **VE single source**: `reference/ve-standard.md` referenced by ALL 6 skills (ends 6x duplication)
4. **context:fork** for heavy discrete operations: doctor, init, review, batch, promote
5. **SKILL.md uppercase** per Anthropic convention
6. **6 skill limit preserved** (Constitution Principle 10)

---

## 4. Feature Audit — What to Keep/Cut/Simplify

### From team.md (4,594 → ~300 lines SKILL.md)

| Feature | Lines Now | Verdict | New Location |
|---------|----------|---------|-------------|
| Dashboard | ~500 | KEEP (simplify) | SKILL.md core |
| Config view | ~40 | KEEP | SKILL.md inline |
| Queue display | ~50 | KEEP | SKILL.md inline |
| MC detail (#N) | ~70 | KEEP | SKILL.md inline |
| Help | ~75 | SIMPLIFY | SKILL.md inline (~30 lines) |
| Init | ~1,300 | context:fork | operations/init.md |
| Doctor + fix | ~770 | context:fork | operations/doctor.md |
| Auto | ~480 | SIMPLIFY to ~100 | operations/auto.md (thin orchestrator) |
| Learn | ~380 | CUT | Move to ONBOARDING.md |
| Wrap | ~450 | SIMPLIFY to ~100 | SKILL.md inline |
| Share | ~80 | CUT | Users can post GitHub comments manually |
| Sync | ~100 | MERGE | Only in tw2-ship |
| Per-user setup | ~400 | MERGE with init | operations/init.md |

### Other Skills

| Skill | Lines Now | Target | Key Changes |
|-------|----------|--------|-------------|
| tw2-claim | 453 | ~400 | Already near target, minor cleanup |
| tw2-issue | 880 | ~400 + batch.md | Extract batch to operations/ |
| tw2-drive | 1,268 | ~450 + team-execution.md | Extract parallel execution details |
| tw2-ship | 1,291 | ~400 + review.md | Extract AI review to operations/ |
| tw2-rc | 1,274 | ~400 + promote.md + hotfix.md | Extract heavy subflows |

### Total Reduction

- Current: 9,760 lines across 6 SKILL.md files
- Target: ~2,350 lines across 6 SKILL.md files + ~2,500 lines in operations/
- SKILL.md context cost per invocation: -75% (from ~4,594 to ~300-450 per skill)

---

## 5. context:fork Candidates

| Operation | Reason | Agent Type |
|-----------|--------|-----------|
| Doctor | 12 checks, heavy output, discrete | Explore or general-purpose |
| Init | 1,300 lines, used once per project | general-purpose |
| Ship Review | Reads entire PR diff, multi-dimension analysis | Explore |
| Issue Batch | Parallel MC generation, fire-and-forget | general-purpose |
| RC Promote | Commit history → changelog, sequential | general-purpose |

---

## 6. Installation Changes

### install.sh Rewrite

```
Old: install.sh → ~/.claude/commands/team.md (flat file symlink)
New: install.sh → ~/.claude/skills/tw2/SKILL.md (directory tree symlink/copy)
```

- Each skill installed as complete directory (SKILL.md + operations/ + reference/)
- Scripts installed to skill-specific location (not global commands/scripts/)
- `${CLAUDE_SKILL_DIR}` resolves paths automatically
- `--copy` for vi_agent copies entire directory tree
- Must not break existing `~/.claude/skills/` contents (3 third-party skills exist)

### Transition Strategy (5 Phases)

1. **Phase 1**: Write new skills in `skills/tw2/`, install.sh supports `~/.claude/skills/`
2. **Phase 2**: Migrate scripts (tw-*.sh → tw2-*.sh), maintain interface compatibility
3. **Phase 3**: Switch erwin to tw2, validate
4. **Phase 4**: Switch vi_agent to tw2
5. **Phase 5**: Delete skills/teamwork/, clean up `~/.claude/commands/` residue

Each phase independently testable, can rollback to teamwork at any time.

---

## 7. Design Principles to Carry Forward

From `docs/teamwork-design-constitution.md` — all 15 principles survive, most relevant:

- **P1**: One-sentence instructions > prescriptive templates
- **P2**: Cognitive layer and rendering layer orthogonal
- **P3**: Identity at line 1
- **P10**: 6 skills, clear boundaries
- **P11**: Von Neumann separation (scripts = mechanisms)
- **P12**: GitHub-first, no shadow state
- **P13**: Micro-languages are load-bearing

---

## 8. Open Questions

- [ ] Final naming (tw2 is working codename)
- [ ] Scripts: rename tw-*.sh → tw2-*.sh, or keep tw-* for backward compat during transition?
- [ ] VE standard: keep in reference/ve-standard.md per skill, or one shared location?
- [ ] Custom agents (.claude/agents/): worth creating specialized agents for doctor, review?
- [ ] How does testing work? → See `skill-tester.md`

---

## 9. Source Material

- Session: `_project/sessions/2026-03-14-A.md` (20 rounds of discussion)
- 7 research scouts (VE audit, file sizes, teammate feasibility, prompt compliance, architecture, feature audit, naming)
- Official docs: https://code.claude.com/docs/en/skills, https://code.claude.com/docs/en/sub-agents
- Current teamwork: `skills/teamwork/` (9,760 lines + 4,473 lines scripts)
- Design constitution: `docs/teamwork-design-constitution.md`
