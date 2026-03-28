# Teamwork vs Wing — Detailed Comparative Analysis

> Date: 2026-03-28
>
> Scope: compare the legacy `teamwork` skill suite and the newer `wing` rewrite across architecture, prompt design, execution model, scripts, installation, compatibility, testing, and product strategy.

---

## Executive Summary

`teamwork` and `wing` solve the same core problem: an AI-native GitHub development lifecycle built around the loop:

`issue -> claim -> drive -> ship -> rc`

But they represent two very different generations of the solution.

- `teamwork` is the mature, feature-rich, battle-tested system.
- `wing` is the architecture rewrite that keeps the protocol, keeps most mechanism-layer power, but re-packs the prompt layer into Anthropic's official `SKILL.md + operations/ + reference/ + scripts/` structure.

The deepest difference is this:

- `teamwork` optimized for capability accumulation inside large monolithic skill files.
- `wing` optimized for context hygiene, official architecture compliance, and future maintainability.

So the rewrite is not mainly a product redesign. It is a packaging and prompt-architecture redesign around an already-proven protocol.

That said, `wing` is not a trivial rename:

- it introduces a new config home (`.wing/`) while preserving fallback compatibility with `.teamwork/` and `.teamspace/`
- it centralizes Visual Encoding into one shared reference file
- it externalizes heavy flows into `operations/`
- it updates install/distribution to official Anthropic skill directories
- it reduces core prompt size dramatically while largely preserving the mechanism layer

The tradeoff is also clear:

- `teamwork` still has richer supporting documentation and a longer maintenance/history trail
- `wing` is cleaner architecturally, but currently thinner in explanatory artifacts

In short:

- If the question is "which one is more evolved as a protocol?" -> `teamwork`
- If the question is "which one is the correct long-term architecture?" -> `wing`
- If the question is "did wing replace the protocol?" -> no
- If the question is "did wing replace the packaging, prompt topology, and install model?" -> yes

---

## Source Corpus

Primary sources reviewed:

- `skills/teamwork/README.md`
- `skills/teamwork/ONBOARDING.md`
- `skills/teamwork/team.md`
- `skills/teamwork/team-issue.md`
- `skills/teamwork/team-claim.md`
- `skills/teamwork/team-drive.md`
- `skills/teamwork/team-ship.md`
- `skills/teamwork/team-rc.md`
- `skills/teamwork/scripts/*.sh`
- `skills/teamwork/CHANGELOG.md`
- `docs/teamwork-ai-manual.md`
- `docs/teamwork-design-constitution.md`
- `_project/ideas/skills/tw2-architecture.md`
- `skills/wing/wing/SKILL.md`
- `skills/wing/wing-issue/SKILL.md`
- `skills/wing/wing-claim/SKILL.md`
- `skills/wing/wing-drive/SKILL.md`
- `skills/wing/wing-ship/SKILL.md`
- `skills/wing/wing-rc/SKILL.md`
- `skills/wing/*/operations/*.md`
- `skills/wing/scripts/*.sh`
- `skills/wing/reference/ve-standard.md`
- `install.sh`
- `_project/CURRENT.md`

---

## At A Glance

| Dimension | Teamwork | Wing | What changed |
|---|---|---|---|
| Product role | Mature GitHub-first AI coordination system | Rewrite of the same protocol into official Anthropic skills architecture | Same lifecycle, different packaging |
| Primary install target | `~/.claude/commands/*.md` | `~/.claude/skills/<skill>/SKILL.md` | Flat commands -> directory-based skills |
| Config home | `.teamwork/` with `.teamspace/` compatibility | `.wing/` with fallback to `.teamwork/` and `.teamspace/` | New canonical home plus migration bridge |
| Skill file shape | 6 large `.md` command files | 6 compact `SKILL.md` files plus `operations/` | Monolith -> decomposed modules |
| VE rules | Repeated inline + external docs | Single shared `reference/ve-standard.md` | Duplicated rules -> shared reference |
| Design meta-layer | Strongly documented in Constitution + manual + changelog | Encoded more implicitly in structure, operations, and architecture rewrite | Less prose, more architectural enforcement |
| Script layer | 11 scripts, battle-tested | 10 scripts, mostly migrated/renamed | Mechanism largely preserved |
| Documentation assets | README + ONBOARDING + CHANGELOG + manual + constitution | Architecture present, but no suite-level README/ONBOARDING/CHANGELOG inside `skills/wing/` | Cleaner runtime shape, thinner docs surface |
| Backward compatibility | Supports `.teamspace` | Supports `.wing`, `.teamwork`, `.teamspace` | Wing explicitly bridges old state |
| Lifecycle maturity | High | Transitional but strategically correct | Wing is the future-facing rewrite |

---

## 1. Shared Core: What Did Not Change

The most important observation is that `wing` does not replace the fundamental teamwork protocol.

The following are structurally preserved:

1. GitHub-first operating model
   `Issue = task`, `PR = delivery`, `Actions = quality gate`, `Milestone = version`

2. Six-skill lifecycle partitioning
   Main dashboard/init, issue creation, claim, drive, ship, rc

3. Mission Contract as the AI execution view
   Issue is the public plan, Contract is the enriched local execution artifact

4. Von Neumann separation
   Deterministic operations stay in scripts; LLM handles judgment and interpretation

5. The claim -> drive -> ship loop
   This remains the center of the whole system

6. Release-candidate lifecycle
   RC remains a first-class concept, not an afterthought

This means `wing` is best understood as a protocol-preserving rewrite, not a product pivot.

---

## 2. Strategic Reframing: Why Wing Exists

The clearest explanation sits in `_project/ideas/skills/tw2-architecture.md`.

That document says the root problem was not that the teamwork protocol was wrong. The problem was that the prompt architecture had become too heavy and too monolithic.

Key stated problems:

- `team.md` reached 4,594 lines
- VE rules consumed huge prompt surface
- repeated fixes to VE compliance were not converging cleanly
- the implementation used zero official Anthropic skills patterns

The rewrite vision is explicit:

- preserve the protocol
- preserve the scripts
- preserve the six-skill boundary
- move to official skill architecture
- reduce per-invocation context cost

This is reinforced by `_project/CURRENT.md`, which describes wing as:

- "Anthropic official skill architecture rewrite"
- 6 `SKILL.md` files
- 10 operations files
- single-source VE reference
- 8 scripts migrated
- `install.sh` dual-mode
- max `SKILL.md` size down from `4,592` to `441`

So the core business case for `wing` is not "new features." It is:

- context discipline
- maintainability
- official architecture alignment
- future extensibility

---

## 3. Philosophy and Identity: Same Cognitive DNA, Different Enforcement

### 3.1 Teamwork's philosophy is heavily explicit

`teamwork` carries a strong explicit meta-layer in:

- `docs/teamwork-design-constitution.md`
- `docs/teamwork-ai-manual.md`
- inline "Design Constitution" references inside multiple skills

This material explains:

- why anti-compliance matters
- why prompt templates can be harmful
- why VE and cognition must stay separate
- why GitHub-first matters
- why six skills is the ceiling
- why scripts own deterministic mechanism

In other words, `teamwork` is not just a set of skills. It is a documented theory of operation.

### 3.2 Wing inherits the philosophy, but expresses it structurally

`wing` keeps the same identity language and many of the same cognitive primers:

- `wing-issue` still uses bone vs flesh, blast radius, competing hypotheses
- `wing-claim` still warns that a bad Contract executed flawlessly is still bad
- `wing-drive` still preserves the "super-sage", anti-compliance, ripple check, relentless execution mindset
- `wing-ship` and `wing-rc` still use responsibility-centered identity framing at the top

But the difference is that `wing` removes direct runtime dependence on `docs/teamwork-design-constitution.md`.

That matters.

`teamwork` says:

- "read the constitution before modifying this skill"

`wing` says, structurally:

- the runtime skill should stay compact
- shared rules should move to `reference/` or `operations/`

### 3.3 Net judgment

- `teamwork` is stronger as a teachable doctrine
- `wing` is stronger as a runtime artifact

This is a real tradeoff, not a contradiction.

---

## 4. Information Architecture: Monolith vs Decomposed Skill Suite

### 4.1 Teamwork architecture

`teamwork` stores almost all user-facing behavior in six large flat markdown command files:

- `team.md`
- `team-issue.md`
- `team-claim.md`
- `team-drive.md`
- `team-ship.md`
- `team-rc.md`

This gives each file broad autonomy, but also makes each skill invocation expensive in prompt surface.

### 4.2 Wing architecture

`wing` reorganizes the suite into:

- per-skill directories with `SKILL.md`
- on-demand `operations/`
- shared `scripts/`
- shared `reference/ve-standard.md`

Layout:

- `skills/wing/wing/`
- `skills/wing/wing-issue/`
- `skills/wing/wing-claim/`
- `skills/wing/wing-drive/`
- `skills/wing/wing-ship/`
- `skills/wing/wing-rc/`

This is exactly the shape proposed in `_project/ideas/skills/tw2-architecture.md`.

### 4.3 Architectural consequence

The decomposition is not cosmetic. It changes how the model consumes the system:

- in `teamwork`, much of the command-specific long tail is always co-located in the skill file
- in `wing`, heavy flows are loaded only when needed

This reduces prompt pollution and improves boundary clarity.

---

## 5. Prompt Economics: The Largest Measurable Win

Measured line counts from the repository:

### 5.1 Teamwork prompt surface

- 6 primary skill markdown files: `9,824` lines total
- largest single skill file: `team.md` at `4,605` lines

### 5.2 Wing prompt surface

- 6 primary `SKILL.md` files: `1,899` lines total
- 12 operation files: `1,289` lines total
- largest single `SKILL.md`: `442` lines

### 5.3 What that means

Skill entry surface:

- `teamwork`: `9,824`
- `wing`: `1,899`

Reduction:

- `80.7%` reduction in core skill-file surface
- `5.17x` smaller at the main skill-file layer

Largest-file reduction:

- `team.md` `4,605` -> max wing skill file `442`
- `90.4%` reduction in worst-case single-file size

Even if you include all `wing` operations:

- `wing` skill + operations = `3,188`

That is still materially smaller than `teamwork`'s `9,824` core skill layer.

### 5.4 Strategic interpretation

This is the main architectural victory of `wing`.

It does not merely "organize files better." It changes the runtime cognitive budget.

That should improve:

- instruction salience
- local coherence
- resistance to signal dilution
- maintainability of future edits

---

## 6. Command Surface and Role Partitioning

### 6.1 Core mapping is 1:1

| Teamwork | Wing |
|---|---|
| `/team` | `/wing` |
| `/team-issue` | `/wing-issue` |
| `/team-claim` | `/wing-claim` |
| `/team-drive` | `/wing-drive` |
| `/team-ship` | `/wing-ship` |
| `/team-rc` | `/wing-rc` |

The lifecycle phases are intentionally preserved.

### 6.2 The main difference is not capability count, but where capability lives

In `teamwork`, many subflows are inline:

- init
- doctor
- auto
- review
- promote
- hotfix
- batch
- share
- wrap

In `wing`, these become explicit operation documents:

- `wing/operations/init.md`
- `wing/operations/doctor.md`
- `wing/operations/auto.md`
- `wing/operations/share.md`
- `wing/operations/wrap.md`
- `wing-issue/operations/batch.md`
- `wing-issue/operations/fix.md`
- `wing-issue/operations/update.md`
- `wing-drive/operations/team-execution.md`
- `wing-ship/operations/review.md`
- `wing-rc/operations/promote.md`
- `wing-rc/operations/hotfix.md`

### 6.3 Design consequence

`teamwork` favors locality:

- the whole subcommand flow is often in the skill file

`wing` favors routing:

- `SKILL.md` becomes identity + argument router + short core path
- heavy procedures live behind explicit load boundaries

This is better aligned with Anthropic's official skill design philosophy.

---

## 7. Visual Encoding: From Repeated Rule Header to Shared Reference

### 7.1 Teamwork's VE problem

The rewrite architecture doc explicitly calls out VE duplication and signal dilution as a major issue.

`teamwork` keeps VE instructions inline in each skill and relies on external documents like:

- `docs/visual-encoding-standard.md`
- `docs/teamwork-design-constitution.md`

This made the rules available, but also made them easy to duplicate, drift, or dilute.

### 7.2 Wing's VE model

`wing` introduces:

- `skills/wing/reference/ve-standard.md`

And all wing skills point to it using:

- `${CLAUDE_SKILL_DIR}/reference/ve-standard.md`

This matters for three reasons:

1. one source of truth
2. reusable by all 6 skills
3. runtime-local relative pathing via the skill directory

### 7.3 Why this is bigger than formatting

VE was not merely a style layer in this project. It was a prompt reliability problem.

So wing's VE centralization is an architectural fix to a prompt bug class.

That is a strong example of `wing` being a systems rewrite, not only a naming rewrite.

---

## 8. Script Layer: Mostly Preserved, Selectively Modernized

### 8.1 Quantitative comparison

Measured script lines:

- `teamwork`: `4,483` lines across `11` scripts
- `wing`: `3,996` lines across `10` scripts

Reduction:

- only `10.9%`

This is a crucial finding.

The prompt layer shrank drastically, but the mechanism layer shrank only modestly.

That means wing preserved most deterministic capability.

### 8.2 Script mapping

| Teamwork script | Wing script | Relationship |
|---|---|---|
| `tw-config.sh` | `wing-config.sh` | migrated, expanded detect-dir logic |
| `tw-git.sh` | `wing-git.sh` | migrated |
| `tw-pr.sh` | `wing-pr.sh` | migrated |
| `tw-label.sh` | `wing-label.sh` | migrated |
| `tw-contract.sh` | `wing-contract.sh` | migrated |
| `tw-notify.sh` | `wing-notify.sh` | migrated |
| `tw-bump-version.sh` | `wing-bump-version.sh` | migrated |
| `setup-github-labels.sh` | `setup-github-labels.sh` | retained with wing-aware config lookup |
| `tw-tips.txt` | `wing-tips.txt` | migrated |
| `tw-unit-test.sh` | `wing-unit-test.sh` | not just renamed; scope expanded |
| `tw-e2e-test.sh` | no direct wing equivalent in scripts directory | end-to-end validation not yet equivalently surfaced |

### 8.3 Direct evidence of script continuity

A diff between `tw-config.sh` and `wing-config.sh` shows the rewrite is mostly conservative:

- same command structure
- same usage model
- same behavior categories
- new priority for `.wing/`
- preserved fallback to `.teamwork/` and `.teamspace/`

This pattern likely generalizes to most wing scripts.

### 8.4 Interpretation

The scripts were not the part that needed radical redesign.

`teamwork` had already solved many deterministic problems through scripts.

`wing` therefore keeps the mechanism layer largely intact and spends rewrite effort where it matters most:

- skill topology
- pathing model
- context loading
- official skills packaging

---

## 9. Config and State Model: Explicit Canonical Shift with Backward Compatibility

### 9.1 Teamwork state model

Canonical local state:

- `.teamwork/config.yml`
- `.teamwork/active/MISSION-*.md`

Compatibility:

- `.teamspace/`

### 9.2 Wing state model

Canonical local state:

- `.wing/config.yml`
- `.wing/active/MISSION-*.md`

Compatibility:

- `.wing/`
- `.teamwork/`
- `.teamspace/`

### 9.3 Why this matters

This is one of the strongest signals that `wing` is intended as a migration-safe rewrite.

`wing` does not force clean-room adoption.

It tries to:

- detect old state
- read old config homes
- preserve existing working repos
- allow staged migration

This is consistent with `_project/CURRENT.md`, which explicitly lists transition phases:

- switch erwin to wing
- switch vi_agent to wing
- only then delete `skills/teamwork/`

### 9.4 Architectural consequence

`wing` is not a fork from zero. It is a compatibility-aware evolution path.

---

## 10. Runtime Pathing: Hardcoded Command Paths vs Skill-Local Resolution

### 10.1 Teamwork pathing

`teamwork` skill files often call scripts via hardcoded global install locations such as:

- `~/.claude/commands/scripts/tw-config.sh`

This works for the old flat-command install model, but couples runtime behavior to a specific installation target.

### 10.2 Wing pathing

`wing` switches to:

- `${CLAUDE_SKILL_DIR}/scripts/...`
- `${CLAUDE_SKILL_DIR}/operations/...`
- `${CLAUDE_SKILL_DIR}/reference/...`

This is a major architectural improvement because it makes the skill self-contained once installed as a directory tree.

### 10.3 Why this matters

This enables:

- symlink installs
- copy installs
- project-targeted installs
- cleaner portability between environments

It also explains the `install.sh` rewrite:

- legacy mode installs flat files to `~/.claude/commands`
- modern mode installs directories to `~/.claude/skills`

This is one of the clearest cases where wing is objectively better engineered.

---

## 11. Installation and Distribution Model

### 11.1 Teamwork install model

Legacy:

- flat markdown commands
- scripts installed to `~/.claude/commands/scripts`

### 11.2 Wing install model

Modern:

- each skill installed as its own directory
- `SKILL.md` at the root
- `operations/`, `scripts/`, `reference/` installed with it

The current `install.sh` supports both models in one installer.

### 11.3 Strategic consequence

This dual-mode installer is not just convenience. It is migration infrastructure.

It allows:

- old workflow continuity
- new architecture adoption
- cross-project copy deployment
- staggered rollout to other repos

This is one of the strongest indicators that wing was built as an intentional transition architecture, not as an experiment.

---

## 12. Drive Layer Comparison: Same Engine, Better Boundary

`team-drive.md` and `wing-drive/SKILL.md` are especially revealing because this is where the protocol's cognitive heart lives.

### 12.1 What remains intact

Wing preserves the key execution ideas:

- mission ownership
- anti-compliance posture
- plan validation before execution
- treating sub-tasks as hypotheses rather than orders
- ripple checking
- self-adversarial review
- relentless forward progress
- execution-mode scaling by task count

This shows the rewrite did not soften or dilute the core drive philosophy.

### 12.2 What changes

The most important change is boundary management:

- teamwork embeds more of the drive machinery inline
- wing pushes team-execution details into `operations/team-execution.md`

This is precisely the kind of separation the architecture doc proposed.

### 12.3 Net effect

The core "thinking OS" of the drive phase survives.

The difference is that wing makes the load-bearing abstraction explicit:

- the main skill should contain the main judgment loop
- the parallel/team execution elaboration should be secondary, loaded only when needed

This is a cleaner model.

---

## 13. Issue and Claim Layers: Very High Semantic Continuity

A direct diff between `team-issue.md` and `wing-issue/SKILL.md` shows high semantic continuity.

What stays the same:

- mission-contract issue model
- bone vs flesh framing
- blast radius awareness
- team-mode vs solo-mode logic
- duplicate detection
- standalone test
- update/fix/batch flows

What changes:

- references to teamwork become wing
- global script paths become `${CLAUDE_SKILL_DIR}` paths
- support logic is offloaded into operations
- help surfaces are shorter and less encyclopedic

The same applies to claim:

- same core Contract generation
- same explicit contract confirmation before branch creation
- same handling of active missions
- same worktree-aware reasoning
- same freshness/comment intelligence model

Conclusion:

`wing-issue` and `wing-claim` are not re-inventions. They are prompt-surface reductions plus architectural normalization.

---

## 14. Ship and RC Layers: Capability Preservation with Cleaner Delegation

### 14.1 Ship

`wing-ship` preserves the same major behaviors as `team-ship`:

- pre-flight validation
- task completeness checking
- PR creation with `Closes` vs `Refs`
- CI watch
- review request
- side-effect triage
- post-merge cleanup

The main difference is that review is explicitly split out to:

- `wing-ship/operations/review.md`

### 14.2 RC

`wing-rc` preserves:

- prepare
- cancel
- fix
- hotfix
- promote

Again the difference is architectural:

- `promote` and `hotfix` become dedicated operation files

This is especially important because RC/promote flows are dense, branch-sensitive, and operationally risky. Offloading them into dedicated operation files reduces top-level skill bloat without removing capability.

### 14.3 Consequence

`wing` is strongest where the original teamwork system was most structurally overgrown:

- rare heavy subflows
- long help surfaces
- multi-step release logic

It turns these from permanent prompt baggage into demand-loaded procedures.

---

## 15. Testing and Validation: Wing Is More Rewrite-Conscious

### 15.1 Teamwork testing

Teamwork includes:

- `tw-unit-test.sh`
- `tw-e2e-test.sh`

This shows mature validation, including end-to-end concerns.

### 15.2 Wing testing

Wing includes:

- `wing-unit-test.sh`

This test suite appears broader than a naive rename:

- script tests
- content accuracy checks
- operation presence checks
- install behavior checks
- compatibility fallback checks for `.teamwork/`
- priority rules for `.wing` over `.teamwork`

That means wing's test emphasis is not only on business behavior, but on architecture integrity.

### 15.3 Important nuance

Teamwork still appears to have more explicit end-to-end lifecycle testing exposed in the repo.

Wing appears stronger in structural validation of the rewritten architecture, but weaker in surfaced end-to-end parity artifacts.

This is exactly what we would expect during a rewrite transition:

- architecture confidence first
- full parity confidence later

---

## 16. Documentation Surface: Teamwork Is Richer, Wing Is Thinner

This is one of the biggest current asymmetries.

### 16.1 Teamwork has a complete knowledge surface

Artifacts include:

- `skills/teamwork/README.md`
- `skills/teamwork/ONBOARDING.md`
- `skills/teamwork/CHANGELOG.md`
- `docs/teamwork-ai-manual.md`
- `docs/teamwork-design-constitution.md`

This makes teamwork:

- easier to learn
- easier to audit historically
- easier to modify with context
- easier to onboard humans into

### 16.2 Wing currently has a thinner explanatory surface

Inside `skills/wing/`, there is:

- `reference/`
- `scripts/`
- per-skill directories

But no suite-level:

- `README.md`
- `ONBOARDING.md`
- `CHANGELOG.md`

in that directory itself.

Documentation does exist elsewhere:

- `_project/CURRENT.md`
- `_project/ideas/skills/tw2-architecture.md`

But it is less end-user facing and less packaged with the suite.

### 16.3 Interpretation

This is not proof that wing is weaker as an architecture.

It does mean:

- wing is stronger as an implementation shape than as a self-explaining product package
- teamwork is still superior as a documented system for human readers

This gap likely matters during migration and downstream adoption.

---

## 17. Compatibility and Migration Strategy: One of Wing's Strongest Qualities

Wing is unusually careful for a rewrite.

Evidence:

- `.wing` as new canonical home
- fallback support for `.teamwork`
- fallback support for `.teamspace`
- dual-mode installer
- current project phase plan explicitly includes coexistence and delayed deletion

This tells us that the rewrite was done with real operational constraints in mind:

- there are already running repos
- there are already downstream consumers
- rollback safety matters

This makes wing a rewrite with production instincts, not just design purity.

---

## 18. What Wing Improves Unequivocally

These are the improvements that are hard to argue against.

### 18.1 Prompt hygiene

This is the biggest win.

### 18.2 Official Anthropic skills compliance

`SKILL.md`, `operations/`, `reference/`, skill-local relative paths, directory installs.

### 18.3 Single-source VE

This directly addresses a known bug class.

### 18.4 Better portability

`${CLAUDE_SKILL_DIR}` is much cleaner than hardcoded `~/.claude/commands/...`.

### 18.5 Better migration posture

Support for old config homes and a dual installer lower adoption risk.

### 18.6 Better boundedness of heavy flows

Operations isolate complexity where it belongs.

---

## 19. What Teamwork Still Does Better Today

These are areas where teamwork remains stronger in practical maturity.

### 19.1 Historical maturity

Teamwork has a long maintenance trail, including many concrete bug fixes visible in `CHANGELOG.md`.

### 19.2 Documentation richness

Manual, constitution, onboarding, README, changelog.

### 19.3 Exposed end-to-end validation artifacts

Teamwork still presents a clearer e2e test story in the repo.

### 19.4 Human legibility as a full system

A human can understand teamwork as a worldview more easily than wing right now.

That does not make teamwork architecturally cleaner. It means its knowledge surface is more complete.

---

## 20. Hidden Insight: Wing's Main Achievement Is Compression Without Protocol Loss

This is the most important architectural judgment in the entire comparison.

Many rewrites achieve cleanliness by deleting capability.

`wing` does something more impressive:

- it compresses the prompt layer very aggressively
- while preserving most of the protocol
- while preserving most of the scripts
- while improving install/runtime architecture

That is a better rewrite than a feature reset.

This is why it is fair to say:

- `teamwork` is the mature protocol artifact
- `wing` is the mature architecture expression of that protocol

They are not competitors in essence. `wing` is the second form of `teamwork`.

---

## 21. Current Gaps or Risks in Wing

A fair comparison should also state what is not yet fully closed.

### 21.1 Documentation parity is incomplete

No suite-level `README.md`, `ONBOARDING.md`, or `CHANGELOG.md` inside `skills/wing/`.

### 21.2 End-to-end parity proof is not yet as obvious

`_project/CURRENT.md` still lists end-to-end validation and migration work as pending.

### 21.3 Anthropic runtime assumptions still exist

`_project/CURRENT.md` explicitly notes a blocker:

- wing scripts reference `${CLAUDE_SKILL_DIR}`, which requires Anthropic skill infrastructure

So wing is architecturally cleaner, but also more dependent on the target runtime behaving as expected.

### 21.4 Transition is not complete

The project itself still describes wing as delivered but still mid-transition.

That means:

- wing is the strategic destination
- teamwork is still part of the operational present

---

## 22. Recommended Framing for Internal Decision-Making

If this comparison is used for project positioning, the strongest framing is:

### Teamwork

"The protocol engine."

Strengths:

- mature
- feature-complete
- deeply documented
- battle-tested

Weaknesses:

- monolithic prompt architecture
- larger context burden
- less aligned with official Anthropic skill structure

### Wing

"The protocol engine rewritten into the correct runtime architecture."

Strengths:

- official skills packaging
- smaller prompt surface
- modular operations
- single-source VE
- stronger portability and migration posture

Weaknesses:

- thinner supporting documentation
- still completing migration and parity hardening

---

## 23. Bottom-Line Judgment

If forced into one sentence:

`teamwork` is the fully evolved protocol, and `wing` is the same protocol rebuilt into its proper architectural form.

If forced into a more operational recommendation:

- keep `teamwork` as the best source for doctrine, history, and comparative confidence
- treat `wing` as the productized future
- do not evaluate wing by asking "does it have all the same prose?"
- evaluate it by asking:
  - did it preserve the protocol?
  - did it preserve the mechanism layer?
  - did it materially improve prompt architecture?
  - did it reduce migration risk?

Based on the current repo state, the answer to all four is yes.

The one major remaining gap is not core architecture. It is packaging completeness around docs, onboarding, and migration proof.

---

## Appendix A — Measured Metrics

### Skill markdown line counts

Teamwork:

- `team.md` = `4,605`
- `team-ship.md` = `1,344`
- `team-rc.md` = `1,274`
- `team-drive.md` = `1,268`
- `team-issue.md` = `880`
- `team-claim.md` = `453`
- Total = `9,824`

Wing:

- `wing-claim/SKILL.md` = `442`
- `wing-issue/SKILL.md` = `336`
- `wing-ship/SKILL.md` = `307`
- `wing-drive/SKILL.md` = `298`
- `wing/SKILL.md` = `267`
- `wing-rc/SKILL.md` = `249`
- Total = `1,899`

Wing operations:

- `12` files
- Total = `1,289`

### Script line counts

Teamwork scripts:

- `11` files
- Total = `4,483`

Wing scripts:

- `10` files
- Total = `3,996`

### Derived metrics

- core skill-layer reduction = `80.7%`
- max single-file reduction = `90.4%`
- script-layer reduction = `10.9%`

Interpretation:

- most of the rewrite happened in the prompt architecture
- not in the deterministic mechanism layer

---

## Appendix B — Simplest Possible Summary

Teamwork answered:

- "What should the AI-native GitHub workflow be?"

Wing answered:

- "How should that workflow be packaged so Anthropic skills can execute it reliably and maintainably?"

That is the whole story in two lines.
