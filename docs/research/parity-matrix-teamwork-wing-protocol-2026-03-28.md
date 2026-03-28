# Teamwork Wing Protocol Parity Matrix

> Date: 2026-03-28
>
> Purpose: classify what should be preserved, redesigned, or deleted when moving from the historical `teamwork` / Anthropic-era `wing` implementations toward a portable protocol plus Codex adapter architecture.

---

## Reading Guide

Decision meanings:

- `Preserve` = keep the capability or invariant essentially intact
- `Redesign` = keep the underlying need, but change packaging or execution form
- `Delete` = remove from the future architecture

Columns:

- `Teamwork` = legacy Claude-command era implementation
- `Wing` = Anthropic skill-era rewrite
- `Protocol Core` = whether this belongs in the runtime-neutral core
- `Codex Adapter Decision` = what the Codex migration should do

---

| Capability / Concern | Teamwork | Wing | Protocol Core | Codex Adapter Decision | Notes |
|---|---|---|---|---|---|
| GitHub as system of record | Strongly present | Strongly present | Yes | Preserve | Foundational invariant, not runtime-specific. |
| Canonical loop `issue -> claim -> drive -> ship -> rc` | Fully present | Fully present | Yes | Preserve | This is the main asset being migrated. |
| Mission Contract as execution artifact | Mature and detailed | Preserved with slimmer packaging | Yes | Preserve | Exact format can evolve; semantics should stay. |
| Standalone-test requirement for issue quality | Explicit and mature | Implicitly preserved | Yes | Preserve | Must remain visible in issue/contract authoring. |
| AI may challenge bad sub-task decomposition | Explicit in drive-era thinking | Present but lighter | Yes | Preserve | Critical anti-compliance rule. |
| Deterministic shell mechanisms for git/labels/pr/contracts | Extensive script layer | Mostly preserved and renamed | Partly | Preserve | Belongs below the skill layer, not in prose. |
| Prefixed label state machine | Mature | Preserved | Adapter/config | Preserve | Runtime-neutral need; adapter-specific implementation. |
| Config-driven branches / release topology / labels | Mature | Preserved with fallback improvements | Yes as invariant | Preserve | No hardcoded project policy. |
| Worktree support | Mature and practical | Preserved | Adapter-level | Preserve | Not protocol core, but valuable runtime capability. |
| PR lifecycle helper operations | Mature | Preserved | Adapter/mechanism | Preserve | Deterministic helpers remain useful. |
| RC lifecycle with promote/fix/hotfix | Mature | Preserved | Yes | Preserve | High-value operational system. |
| Dashboard home surface | Present as `/team` | Present as `/wing` | No | Redesign | Keep concept, not the exact interaction style. |
| Doctor diagnostics | Mature, high leverage | Preserved | No | Preserve | Should likely merge tightly with Codex home/dashboard. |
| Doctor auto-fix flows | Present | Preserved | No | Redesign | Keep only deterministic safe repairs. |
| Batch issue generation | Present | Present as operation | No | Preserve | Useful accelerator, but still issue-surface behavior. |
| Issue update / fix sub-operations | Present | Present as operations | No | Preserve | Keep capability, simplify scaffolding. |
| Team roster / team member management | Strongly present | Reduced importance | No | Delete | Overfits a multi-human team model; protocol does not require it. |
| Explicit “team mode” worldview | Central in legacy framing | Weakened | No | Delete | The protocol is not fundamentally a team-collab tool. |
| Cross-member assignment semantics | Strongly present | Partially preserved | No | Redesign | Keep assignee/ownership support, remove team metaphysics. |
| Cross-MC knowledge sharing (`share`) | Present | Preserved as operation | No | Redesign | Keep as optional issue/PR knowledge sync, not identity-defining. |
| Menu-heavy AskUserQuestion trees | Extensive | Reduced but still visible | No | Delete | Codex should use concise direct flow, not large menu choreography. |
| Old Claude command syntax as primary packaging | Central | Removed | No | Delete | Historical shell only. |
| Anthropic skill packaging | Absent | Central | No | Redesign | Skill modularity is valuable, Anthropic specifics are not. |
| `${CLAUDE_SKILL_DIR}` as core path assumption | N/A in older form | Central runtime assumption | No | Delete | Replace with Codex-native path resolution or repo-local layout. |
| `.teamwork/ -> .teamspace/` historical dual-dir baggage | Present | Expanded fallback to `.wing/` | No | Delete | Keep migration tooling only; not future architecture. |
| `.wing/ -> .teamwork/ -> .teamspace/` fallback chain | N/A | Valuable migration bridge | No | Redesign | Useful during transition, not as permanent conceptual model. |
| Single-source VE reference | Duplicated across prompts | Centralized | No | Preserve | Keep the lesson: one style source if style matters. |
| Huge monolithic prompt files | Major pain point | Explicitly reduced | No | Delete | This was a root maintainability problem. |
| Split operations files | Weak in legacy | Strong in wing | No | Preserve | Good modularity pattern. |
| Onboarding docs and manual depth | Very strong | Thinner | No | Redesign | Keep documentation quality, but separate docs from runtime prompts. |
| E2E testing harness for protocol flows | Present in scripts | Weaker / pending parity | No | Preserve | Essential for future runtime migration confidence. |
| Static skill linting | Not first-class | Stronger via `skill-lint` | No | Preserve | Important for harness quality. |
| Issue freshness / info sync | Mature in ship/contract ops | Present in lighter form | Yes as invariant | Preserve | Freshness matters across runtimes. |
| Contract freshness hash / drift detection | Mature | Preserved | Yes | Preserve | Key safety mechanism. |
| Branch naming conventions with issue linkage | Mature | Preserved | Adapter/config | Preserve | Important for traceability. |
| Automatic claim comment / ownership signaling | Present | Present | Adapter/policy | Preserve | Good collaboration signal even in solo/human-light mode. |
| Auto-launch from claim into drive | Present in some flows | Present | No | Redesign | Keep as optional continuation, not forced ritual. |
| Contract checkbox task tracking | Mature | Preserved | No | Preserve | Useful execution artifact, though exact storage can evolve. |
| AI PR review operation | Present in ship | Preserved via operation | No | Preserve | Valuable optional surface. |
| Rebase / sync operation | Present | Less prominent | No | Preserve | Still needed in real repos. |
| Post-merge cleanup (`done`) | Present | Present in ship flow | No | Preserve | Deterministic cleanup remains important. |
| Release staging workflow trigger | Present via config | Present via config | Adapter/config | Preserve | Good release integration point. |
| Production hotfix + backport model | Mature | Preserved | Yes | Preserve | Important for serious repo ops. |
| Duplicate issue detection | Present in issue flow | Lighter | No | Preserve | Practical quality feature. |
| “Mission” language everywhere | Central | Preserved | Partly | Redesign | Useful metaphor, but should not block neutral protocol framing. |
| “Teamwork” branding as architectural center | Central | Replaced | No | Delete | Historical brand, not future architecture. |
| “Wing” branding as automatic Codex successor | Not applicable | Tempting but unclear | No | Delete for now | Do not hard-bind the Codex adapter to `wing` prematurely. |
| `vi_agent`-shaped defaults and examples | Strongly influential in practice | Some influence remains | No | Delete | Keep as one profile/example, not the hidden reference repo. |
| Human approval checkpoints before destructive actions | Present | Present | Adapter safety | Preserve | Important in release/cleanup flows. |
| Interactive dashboard as identity surface | Strong | Strong | No | Redesign | Good UX surface, but not the protocol itself. |
| Tool/runtime capability mapping | Implicit, Claude-shaped | Anthropic-shaped | No | Preserve | Must be explicit in Codex adapter docs. |
| Portable adapter interface between protocol and runtime | Weak / implicit | Emerging but mixed with product | No | Preserve | This is the central new design move. |
| Separate project profiles from runtime adapter | Incomplete | Better but still mixed | Architecture invariant | Preserve | Needed to stop future overcoupling. |
| Harness engineering as first-class concern | Implicit | Emerging | No | Preserve | The migration should incubate this deliberately. |

---

## Highest-Impact Takeaways

### Preserve intact

- protocol loop
- GitHub-first worldview
- contracts
- deterministic script layer
- freshness and verification semantics
- PR and RC lifecycle depth

### Redesign aggressively

- dashboard UX
- claim-to-drive choreography
- issue interaction style
- docs/runtime boundary
- adapter naming and packaging
- migration fallback logic

### Delete without regret

- command-era path assumptions
- monolithic prompt structure
- team-centric identity as the protocol definition
- hidden `vi_agent` defaults
- any notion that Anthropic runtime details are part of the protocol

---

## Migration Implication

The Codex path should therefore be:

1. preserve protocol semantics
2. preserve deterministic mechanisms where still valid
3. redesign runtime surface as a Codex-native adapter
4. delete historical packaging debt
5. postpone final product naming until the adapter and protocol layers are cleanly separated

That path gives `wing` room to remain a future orchestration concept if needed, while still allowing the protocol itself to migrate cleanly into Codex.
