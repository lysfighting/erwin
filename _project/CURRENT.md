# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-28
- **Last Session**: 2026-03-28 (Codex-native orchestration harness prototype — release governance upgraded)

## Current Phase

**Phase**: Codex-native orchestration harness prototype

## What Was Just Completed

**Session A** (2026-03-28) — Harness prototype became a stronger end-to-end release harness

- **Committed definition layer**: `harness/` now has profile, templates, and policies for a Codex-native orchestration harness
- **Deterministic substrate**: `harness_state.py` plus shell wrappers now implement `home`, `doctor`, `claim`, `drive`, `meetings`, `maintenance`
- **Leader control plane**: run-state transitions, issue drafting, meeting materialization, meeting outcome adoption, reprovisioning, reconciliation, normalization, ship push / PR draft / PR tracking lifecycle, and RC governance
- **Ship governance upgraded**: PR review state and checks state are now first-class ship gates instead of informal notes
- **Merge action upgraded**: `ship merge` now performs a deterministic merge path instead of relying only on manual `merged` state updates
- **RC hard cases absorbed**: `scripts/harness-rc.sh` now supports surgical `fix`, `cancel`, `hotfix`, and `backport` flows in addition to create / verify / promote / close
- **Runtime governance**: file-backed locking prevents concurrent registry loss; `doctor` reached clean runtime state (`issue_count: 0`)
- **Codex-native surface**: local `skills/harness/SKILL.md` created as the first repo-local harness skill
- **Codex installer path**: `install-codex.sh` now installs repo Codex-native skills like `harness` into `~/.codex/skills/`, and `/harness` has been installed globally
- **Bootstrap path**: `scripts/harness-bootstrap.sh` now seeds the harness adapter into another git repository with a generated default profile, local `skills/harness/`, and `install-codex.sh`
- **Adapter sync path**: bootstrapped repos now record source metadata in `.harness/cache/adapter-source.json` and can resync via `scripts/harness-sync.sh`; the bootstrapped repo also carries its own `harness-bootstrap.sh`
- **Adapter packaging upgraded**: `harness/adapter-manifest.yml` now declares adapter identity/capabilities, and `scripts/harness-all-smoke-test.sh` provides a one-command smoke suite
- **Adapter packaging tightened**: bootstrap smoke now validates manifest + full-smoke installation, `harness-sync status` now reports manifest version alignment, and the adapter version is now `0.4.0-prototype`
- **Validation**: branch/worktree provisioning, issue drafting, meeting create/update/close, archived meeting replay, PR review/check gates, deterministic merge, doctor/home health checks, direct `./install-codex.sh` execution, bootstrap into a temporary repo, remembered-source resync, `harness-issue-smoke-test.sh`, `harness-smoke-test.sh`, `harness-bootstrap-smoke-test.sh`, `harness-release-smoke-test.sh`, and `harness-all-smoke-test.sh` all pass
- **Protocol coverage**: the Codex adapter now has first-class `issue`, `claim`, `drive`, `ship`, and `rc` surfaces, including RC surgical fix / cancel / hotfix / backport handling

## Active Task / What's In Progress

**Harness prototype delivered. Productization tasks remain:**

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 1 | ✅ | Research protocol / control plane / authority / cybernetics synthesis |
| Phase 2 | ✅ | Build repository-local harness definition layer (`harness/`) |
| Phase 3 | ✅ | Implement runtime substrate (`scripts/harness_*`) |
| Phase 4 | ✅ | Validate `home` / `doctor` / `claim` / `drive` / `meetings` / `maintenance` |
| Phase 5 | 🔄 | Package the harness more formally as a reusable Codex-native skill/system |
| Phase 6 | 🔄 | Split reusable adapter/bootstrap flow from repo-local proving ground |

## Next Session Should Start With

1. Read this file
2. **Harness next work**:
   - A. Decide whether to formalize this under a new product name or keep `harness` as the prototype codename
   - B. Harden the reusable bootstrap/adapter path beyond first-copy scaffolding
   - C. Decide whether to promote `skills/harness/` beyond prototype-local installation
   - D. If starting a fresh session, read `AGENTS.md`, `_project/CURRENT.md`, `_project/sessions/2026-03-28-A.md`, `harness/README.md`, and inspect recent commits on `codex/harness-release-lifecycle`
3. **Showcase tasks still pending**: DesignPatternLibrary, CaseStudies — blocking v1.0.0

## Implementation Roadmap

### v0.5.0 — V2 Phase 2: Deep Interaction (in progress)
- [x] OriginTrace (d3-force philosophy graph, ACT II)
- [x] Bundle code splitting (627KB → 221KB, React.lazy + vendor chunks)
- [x] Code health sweep (ESLint clean, dead code removal)
- [x] ExecutionSimulator: message flow animation enhancement
- [x] SwarmTopologySection: add interactive message flow

### v1.0.0 — V2 Phase 3: Creation Path
- [x] SkillBuilder (form → skill.md download, ACT IV)
- [ ] DesignPatternLibrary (3-5 reusable pattern cards)
- [ ] CaseStudies (2-3 skill variant examples)

### Skills Factory
- [x] Teamwork v2.0–v3.11.2 — 40+ versions, 6 skills, 10 scripts, 48 unit tests
- [x] Wing v0.1.0 — Anthropic skill architecture rewrite (6 SKILL.md, 10 operations, 8 scripts)
- [x] skill-lint.sh — Layer 1 static analysis (12 checks)
- [x] arch-review v0.1.0 — Multi-role debate review (8 mechanisms, 27 bugs fixed in 5 rounds)
- [x] Harness prototype — Codex-native repository-local orchestration harness (leader control plane + runtime substrate)
- [ ] Wing Phase 4: vi_agent migration
- [ ] Wing Phase 5: teamwork deprecation
- [ ] Skill-tester Layer 2: smoke tests

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Wing scripts reference `${CLAUDE_SKILL_DIR}` which requires Anthropic skill infrastructure
- skill-lint.sh ### header check has false positives (structural headings vs output templates)
- Harness is now installed and validated in Codex via `install-codex.sh`, can bootstrap into another repo, and now covers the main RC hard-case flows including dual-branch backport, but is still a prototype rather than a finalized product surface

## File Locations

| What | Where |
|------|-------|
| Project goals | `_project/GOALS.md` |
| Idea garden | `_project/ideas/` |
| ADRs | `_project/decisions/` |
| V2 architecture | `docs/drive-showcase-v2-architecture.md` |
| Wing skills | `skills/wing/` |
| Harness prototype | `harness/`, `scripts/harness_*`, `skills/harness/` |
| Wing architecture | `_project/ideas/skills/tw2-architecture.md` |
| Teamwork skills (legacy) | `skills/teamwork/` |
| Docs index | `docs/README.md` |
| App root | `showcase/src/App.tsx` |
