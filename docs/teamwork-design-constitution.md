# Teamwork Design Constitution

> **Read this before modifying ANY teamwork skill file.** These are load-bearing design principles — each earned through real bugs, failed experiments, or multi-round design discussions. Violating them without understanding WHY they exist will introduce regressions.

---

## How to Use This Document

1. **Before editing a skill**: Scan the principles below. If your change touches a principle, read the WHY.
2. **If you want to violate a principle**: You may — but document the reasoning in a new ADR, update this file, and get user confirmation.
3. **Violation signals**: Each principle lists patterns that indicate you're about to violate it. If you see yourself doing one of these, pause and re-read the principle.

---

## Principle 1: One-Sentence Instructions > Prescriptive Templates

**Rule**: Cognitive instructions should be general thinking prompts (1-2 sentences), not multi-line forced output templates.

**Why**: Anthropic's prompt engineering research (2026) shows Claude's extended thinking "frequently exceeds what a human would prescribe." Prescriptive reasoning templates (e.g., "fill in these 5 fields: Right problem? Better framing? ...") can limit Claude's reasoning quality. They also consume context window, which is a public resource. (Session AL, Round 8)

**Violation signals**:
- Adding a 3+ line structured output template that the AI must fill in
- Requiring specific field names in a thinking output (e.g., "Hidden assumption: {X}")
- Adding a "gate" that produces a fixed-format block

**The nuance**: Requiring VISIBLE OUTPUT is fine ("state your position before asking"). Requiring a FIXED FORMAT for that output is the anti-pattern. The distinction: tell the AI WHAT to show (its position, trade-offs), not HOW to format it.

**Source**: Session AL Round 8; Anthropic docs on prompt engineering, Claude Code skill design patterns

---

## Principle 2: Cognitive Layer and Rendering Layer are Orthogonal

**Rule**: Cognitive instructions (how the AI thinks) and VE instructions (how the output looks) must never be coupled.

**Why**: Coupling means changing one breaks the other. VE is about ANSI rendering constraints (bold, backtick, italic). Cognitive instructions are about reasoning quality. They evolve at different rates for different reasons. (Session AL, Round 9)

**Violation signals**:
- A cognitive instruction that specifies output formatting (e.g., "output your analysis as a table with columns X, Y, Z")
- A VE rule that implies reasoning behavior
- A "challenge gate" that includes formatting requirements

**Source**: Session AL Round 9

---

## Principle 3: Identity Language at Line 1

**Rule**: Each skill's first `>` blockquote must establish the AI's identity and cognitive posture BEFORE any process steps.

**Why**: By the time the AI reads 300+ lines of process steps, its cognitive posture is already set to "executor." Identity language at the top sets the posture to "owner/thinker" before the pipeline kicks in. drive.md's core insight: "You own this mission" at line 1 creates a fundamentally different cognitive frame than the same words at line 370. (Session AL, Round 7, finding #1)

**Violation signals**:
- Moving identity language below the help section or into the execution engine
- Replacing identity language with process description
- Starting a skill with routing/parsing logic instead of cognitive framing

**Source**: Session AL Round 7; drive.md line 1 analysis

---

## Principle 4: Anti-Compliance = Positive Menu, Not Just Prohibition

**Rule**: Anti-compliance must provide the three valid response types (AGREE+EXTEND / CHALLENGE+PROPOSE / REDIRECT), not just say "don't comply."

**Why**: "Don't comply" is a prohibition — it tells the AI what NOT to do. The three response types are a menu — they tell the AI what TO do instead. Menus are cognitively easier for LLMs to follow than prohibitions. Trade-off honesty ("state what your alternative costs") prevents the LLM's default behavior of presenting alternatives as strictly better. (Session AL, Round 12, Thinker A finding #5)

**Violation signals**:
- Removing the response type menu and leaving only "don't comply" or "silence is complicity"
- Removing "trade-off honesty" from anti-compliance instructions
- Adding anti-compliance without the positive alternatives

**Source**: Session AL Round 12; drive.md lines 35-40

---

## Principle 5: Standalone Test at Creation Time

**Rule**: Both team-issue (Issue body) and team-claim (Contract) must pass a Standalone Test: "Can someone who never saw the original context execute from this alone?"

**Why**: In teamwork, the creator and executor are often different people (or different AI sessions). The Issue creator's context evaporates after creation. If critical information exists only in the creator's head (or conversation context) but not in the Issue/Contract, the executor will misunderstand the mission. drive.md's Standalone Test prevents this by forcing explicit documentation of implicit knowledge. (Session AL, Round 7, finding #4)

**Violation signals**:
- Removing or weakening the Standalone Test instruction
- Assuming the executor has access to the creator's conversation context
- Generating an Issue/Contract that references "the discussion" or "as we discussed"

**Source**: Session AL Round 7; drive.md "Standalone Test" section

---

## Principle 6: Blast Radius is Not Priority

**Rule**: Assess blast radius ("if this is wrong, how much breaks?") separately from priority (P0-P3) and size (S-XL). High blast radius → surface it explicitly to the user in preview/confirmation steps.

**Why**: Priority measures importance/urgency. Size measures effort. Blast radius measures downside risk of being wrong. A P2/M task can have enormous blast radius (e.g., changing a shared interface). Without blast radius, teamwork applies uniform scrutiny to all tasks regardless of error cost. The prescription (flag it to the user) ensures the assessment changes behavior, not just generates text. (Session AL, Round 12, all 3 thinkers independently identified this gap)

**Violation signals**:
- Removing blast radius assessment and relying only on Priority/Size
- Treating Priority as a proxy for risk (it isn't — P0 means urgent, not risky)

**Source**: Session AL Round 12; drive.md Phase Omega step 3

---

## Principle 7: Proactive Triggers Belong Upstream

**Rule**: Mediocre consensus detection and implicit decision surfacing must exist in team-issue and team-claim, not just team-drive.

**Why**: team-issue is the pipeline entry point. A badly framed Issue propagates through Contract into execution. By the time team-drive catches it, rework cost is maximized. The highest-leverage place to catch direction errors is at the cheapest-to-fix stage. (Session AL, Round 12, Thinkers A+B+C all identified this)

**Violation signals**:
- Removing proactive trigger language from team-issue or team-claim
- Adding proactive triggers only to team-drive (execution) but not upstream (planning)

**Source**: Session AL Round 12; drive.md Proactive Triggers table

---

## Principle 8: Opposite Test (Inversion)

**Rule**: Before committing to an approach, ask "what happens if we don't do this at all?"

**Why**: Inversion surfaces hidden assumptions that direct questioning misses. If skipping the task entirely causes no damage, the task is probably unnecessary. Charlie Munger's "invert, always invert" — the most powerful reasoning shortcut for catching incorrect framings. (Session AL, Round 12, Thinkers A+C)

**Violation signals**:
- Removing the "test the opposite" instruction from Plan Validation
- Replacing it with only forward-looking questions ("is it correct? is it sufficient?")

**Source**: Session AL Round 12; drive.md Gate 1 "Opposite test"

---

## Principle 9: Contract Confirmation Before Branch Creation

**Rule**: team-claim must show the generated Contract to the user and get explicit confirmation BEFORE creating a branch.

**Why**: Before v3.9.9, team-claim generated the Contract and immediately created a branch without user verification. This meant a badly generated Contract (wrong objective, missing constraints) would already have infrastructure (branch, GitHub comment) before the user could catch the error. The confirmation step is the last chance to catch direction errors before execution begins. (Session AL, Round 7, finding #6)

**Violation signals**:
- Removing Step 4b or making it optional
- Moving branch creation before Contract confirmation
- Auto-confirming without user input

**Source**: Session AL Round 7; drive.md Gate 3 "Contract Challenge"

---

## Principle 10: 6 Skills, Clear Boundaries

**Rule**: Teamwork has exactly 6 skills: /team, /team-issue, /team-claim, /team-drive, /team-ship, /team-rc. Each has a single responsibility. Do not add a 7th skill or merge two.

**Why**: Each skill maps to one phase of the GitHub workflow. Skill boundaries provide natural context separation (each invocation starts fresh). Adding skills causes command explosion (competitive analysis: CCPM had 30+ commands → abandoned). Merging skills causes context pollution. (Session AL, Round 8; MEMORY.md "6 skill limit"; competitive analysis of Vibe Kanban and CCPM)

**Violation signals**:
- Proposing a new `/team-*` command
- Merging two skills into one file
- Adding Research Swarm or Context Handoff within a single skill (the skill boundary already provides this)

**Source**: Session AL Rounds 5-8; ADR-0010; competitive landscape analysis

---

## Principle 11: Von Neumann Separation — Scripts = Mechanisms, Skills = Content

**Rule**: Deterministic operations (branch creation, label management, config parsing, checkbox toggling) must be in `scripts/*.sh`. Skill `.md` files contain only LLM-interpreted content (analysis, judgment, dialogue).

**Why**: v2.6.1 had 21 bugs. 18 of them were deterministic operations left to LLM interpretation (wrong branch names, missed label prefixes, inconsistent checkbox formats). Moving these to scripts eliminated the entire class of bugs. The separation is: scripts do things that have ONE correct answer; skills do things that require judgment. (MEMORY.md "Von Neumann separation"; Session Z, v2.6.1 post-mortem)

**Violation signals**:
- Adding bash logic inline in a skill file that could be a script subcommand
- Having the LLM generate branch names, label strings, or file paths by reasoning instead of calling a script
- Duplicating logic that already exists in `tw-*.sh` scripts

**Source**: ADR-0010; v2.6.1 bug analysis; MEMORY.md

---

## Principle 12: GitHub-First — No Shadow State

**Rule**: GitHub Issues = tasks, PRs = delivery, Actions = CI, Milestones = versions. No local databases, no custom Issue storage, no state that doesn't live in GitHub or `.teamwork/active/`.

**Why**: Teamwork's positioning is "AI that works within your existing GitHub workflow." Building shadow infrastructure (local SQLite, custom dashboards, separate issue trackers) would repeat the exact mistakes that killed competitors (Vibe Kanban abandoned GitHub Issues for SQLite → lost). The Contract in `.teamwork/active/` is the only local state, and it's ephemeral (deleted after ship). (ADR-0010; competitive analysis)

**Violation signals**:
- Adding a local database or persistent state file beyond Contract
- Building a GUI or dashboard outside of CLI output
- Creating a custom Issue format that doesn't map to GitHub Issues
- Storing task state anywhere other than GitHub + Contract checkboxes

**Source**: ADR-0010; competitive landscape analysis; MEMORY.md "Traps"

---

## Principle 13: Micro-Languages are Load-Bearing

**Rule**: Single-sentence cognitive principles scattered through skill files ("error narrows problem space", "silence = complicity", "complexity = unfinished thinking", "a bad Contract executed flawlessly is still a bad outcome") are NOT decorative comments. They are cognitive OS instructions that prime the LLM's reasoning posture.

**Why**: These micro-languages work by priming — they set the AI's default reasoning frame before it encounters specific instructions. Removing them because they "look like comments" or "don't seem functional" degrades reasoning quality in ways that are invisible to functional testing. drive.md uses this technique extensively. (Session AL, Round 7, finding #2)

**Violation signals**:
- Removing a sentence from a `>` blockquote because it "doesn't do anything"
- "Cleaning up" identity language or philosophical statements as unnecessary prose
- Replacing a vivid micro-language with a generic instruction

**Source**: Session AL Round 7; drive.md analysis of distributed cognitive primers

---

## Principle 14: VE Consistency Across All Skills

**Rule**: All 6 teamwork skills must have identical VE header lines. Any VE change must be applied to all 6 simultaneously.

**Why**: VE determines how output renders in the terminal (ANSI colors, formatting). Inconsistent VE across skills means the same concept renders differently in different commands — confusing to users and breaking visual expectations. The VE header is a single source of truth replicated across skills; partial updates create visual regression. (docs/visual-encoding-standard.md; MEMORY.md "VE palette is at ceiling")

**Violation signals**:
- Changing the VE header in one skill but not the others
- Adding a new VE token (e.g., a new use of backtick or bold) in one skill without checking the standard
- Removing the VE header reference line

**Source**: docs/visual-encoding-standard.md; Session AK VE audit

---

## Principle 15: Never End on a Period — Always End on a Question

**Rule**: Every skill output's final step must offer the user actionable next steps. Never output a summary and stop. Even error STOPs must present recovery choices.

**Why**: Today's UX failures (ship pre-flight blocking without options, promote ending without deploy guidance, drive completion dumping text without AskUserQuestion) all share one root cause: the skill ended on a period instead of a question. The user was left wondering "what do I type next?" Every transition between skills (issue→claim→drive→ship→done) is a moment of momentum — if the skill doesn't carry it forward, the user loses flow. (Session AN, multiple user-reported UX dead-ends)

**Three ending patterns**:
1. **Success transition** → AskUserQuestion with next-step options (e.g., "Ship now?" / "Review first?")
2. **Error recovery** → AskUserQuestion with fix options (e.g., "Run doctor?" / "Retry?" / "Skip?")
3. **Terminal state** (help output, user-initiated cancel, waiting for external event like PR merge) → plain text OK, no AskUserQuestion needed

**Violation signals**:
- A completion summary that ends with a plain-text command hint (e.g., `/team-ship to create PR`) instead of AskUserQuestion
- A STOP that describes the problem but offers no recovery choices
- Any output where the user must independently figure out what command to run next

**Source**: Session AN; UX audit (47 dead-ends found, 17 requiring AskUserQuestion)

---

## Architectural Decisions NOT Adopted (and Why)

These drive.md features were explicitly evaluated and intentionally not adopted. They are NOT bugs or gaps.

| Feature | Why Not | Session Reference |
|---------|---------|-------------------|
| Structured Gate Templates (⚡ ESSENCE CHALLENGE format) | Prescriptive templates limit Claude's reasoning (Principle 1) | AL Round 8 |
| Research Swarm in team-issue | team-issue does lightweight analysis; deep research is for /drive | AL Round 8 |
| Dual Wheel as Cardinal Rule | Skill boundary separation handles this (Principle 10) | AL Round 8 |
| Context Handoff / /clear + bootstrap | Skill boundary separation handles this (Principle 10) | AL Round 8 |
| Psychiatrist Test | team-issue doesn't do multi-round dialogue; Standalone Test covers doc completeness | AL Round 11 |
| state.md separate from Contract | Contract checkboxes + GitHub Issue sync = sufficient state tracking | AL Round 12 |

---

## Version History

- **v1.0** (2026-03-13, Session AL): Initial constitution. 14 principles + 6 non-adopted decisions. Derived from 12 rounds of design discussion, 2 rounds of independent brainstorming (6 thinkers total, 56 findings).
- **v1.1** (2026-03-13, Session AN): Added Principle 15 (Never End on a Period). Derived from UX dead-end audit across all 6 skills.
