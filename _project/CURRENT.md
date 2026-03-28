# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.7.0 (in progress)
- **Last Updated**: 2026-03-29
- **Last Session**: 2026-03-29 (Codex-native orchestration harness selector-aware ownership transfer pass completed)

## Current Phase

**Phase**: Codex-native orchestration harness prototype

## What Was Just Completed

**Session 2026-03-29** — Harness runtime migrated from Python to Go end-to-end, then gained packaged-binary preference

- **Go runtime core added**: introduced a real Go module plus `cmd/harness-state` / `internal/harness/runtime.go`, migrating the harness control-plane and state-engine behavior out of `scripts/harness_state.py`
- **Shell surfaces preserved**: all first-class harness entrypoints now route through `scripts/harness-runtime.sh`, which builds and executes the local Go CLI while keeping the existing `home`, `doctor`, `issue`, `claim`, `drive`, `teammates`, `ship`, `rc`, `meetings`, `maintenance`, and adapter/profile/path surfaces stable
- **Python runtime dependency removed for harness usage**: normal harness operation, package inspection, bootstrap, sync, and smoke execution no longer depend on `python3`
- **Packaged runtime preference added**: `scripts/harness-runtime.sh` now prefers `harness/bin/<os>-<arch>/harness-state` when present, with local Go build under `.harness/bin/` only as a fallback
- **Binary build and dist validation added**: `scripts/harness-build-runtime.sh` now builds packaged binaries, and `scripts/harness-runtime-dist-smoke-test.sh` proves the adapter can run without `python3` and without `go` in `PATH` when the packaged runtime is present
- **Council governance hardened**: meetings now carry first-class `role_specs` with allowed/forbidden sources, known-context notes, locked decisions, and rejected-finding reasons; `doctor` validates those governance requirements instead of treating councils as loose notes
- **Gap audit written down**: `docs/research/harness-gap-audit-2026-03-29.md` now records which original design inputs are absorbed, partial, or still missing, so future sessions can resume from an explicit audit instead of replaying chat history
- **Research drift audit added**: `docs/research/harness-research-drift-audit-2026-03-29.md` now re-checks the original design docs against the current harness and confirms the remaining highest-risk drift areas are desired-state control, explicit ownership/overlap, harness-wide authority enforcement, and automated council orchestration
- **Ownership and desired state strengthened**: runs and contracts now carry explicit `ownership_claims`, `desired_state`, and dependency policy; teammate assignment now rejects exclusive overlap by default, and `drive wait` plus `home`/`doctor` surface dependency-blocked runs directly
- **Authority and council control advanced**: invalid authority levels are now rejected or surfaced, and meetings can advance phases with `meetings orchestrate` instead of relying entirely on manual phase jumps
- **Mission graph made first-class**: the harness now materializes durable mission objects under `.harness/missions/`, exposes them through `scripts/harness-missions.sh`, and drives `home`/`doctor` from mission-level state instead of only run-level state
- **Legacy runtime migration path added**: `scripts/harness-maintenance.sh normalize` now backfills ownership claims, desired state, dependency policy, and mission registry state for older pre-mission runs, bringing existing `.harness/` state forward without manual surgery
- **Selector-aware ownership transfer added**: runs now normalize legacy `ownership_claims` into structured `ownership_selectors`, `runs transfer-ownership` can formally move bounded selectors between active runs in the same mission, and transfer events now survive into mission-level state
- **Adapter/package flow updated**: `harness/adapter-manifest.yml` now packages the Go module files and build runner, bootstrapped repos receive the Go sources, and `harness-bootstrap.sh` / `harness-sync.sh` now query adapter metadata through the Go runtime path
- **Smoke suite de-Pythonized**: issue/core/release/teammate smoke scripts now query package files through the harness adapter surface and extract runtime ids without Python helpers
- **Validation**: `./scripts/harness-build-runtime.sh`, `./scripts/harness-runtime-dist-smoke-test.sh`, `./scripts/harness-adapter.sh validate`, `./scripts/harness-issue-smoke-test.sh`, `./scripts/harness-smoke-test.sh`, `./scripts/harness-teammate-smoke-test.sh`, `./scripts/harness-release-smoke-test.sh`, `./scripts/harness-bootstrap-smoke-test.sh`, and `./scripts/harness-all-smoke-test.sh` all pass cleanly on the Go runtime; `./scripts/harness-maintenance.sh normalize` also upgrades older runtime state into mission-aware schema cleanly

### Previous completion

**Session C** (2026-03-28) — Harness prototype crossed from leader-centric control plane into first-class teammate runtime

- **Teammate runtime added**: `scripts/harness-teammates.sh` plus the harness runtime `teammates ...` surface now create bounded delegated runs with scope, authority refs, forbidden sources, reporting expectations, and parent/leader linkage
- **Governed lateral communication added**: teammate-to-teammate communication is now explicit and structured instead of implicit chat, using sibling-approved signal flow and leader-approved linkage
- **Council protocol strengthened**: meetings now track `phase`, `participants`, and `positions`, so council packets can represent independent review and cross-examination instead of only final findings
- **Home/doctor upgraded**: `home` now surfaces active teammate and coordination signal counts; `doctor` now validates teammate assignment quality plus meeting participant/phase/position health
- **Smoke coverage expanded**: `scripts/harness-teammate-smoke-test.sh` now validates delegation, report emission, sibling communication, meeting escalation, participant rostering, and council positions end-to-end
- **Adapter/package contract advanced**: teammate surface and teammate smoke are now part of `harness/adapter-manifest.yml`, and the reusable adapter line advanced to `0.6.0-prototype`
- **Local runtime cleaned**: legacy smoke teammate runs/signals were retired and reconciled so `scripts/harness-doctor.sh` returns cleanly again in the real repository runtime
- **Validation**: `scripts/harness-smoke-test.sh`, `scripts/harness-teammate-smoke-test.sh`, `scripts/harness-release-smoke-test.sh`, `scripts/harness-bootstrap-smoke-test.sh`, `scripts/harness-all-smoke-test.sh`, and `scripts/harness-doctor.sh` all pass cleanly

- **Manifest became executable contract**: `harness/adapter-manifest.yml` now declares protocol identity, adapter package files, bootstrap contract, and profile contract instead of only loose packaging metadata
- **Adapter surface added**: `scripts/harness-adapter.sh` now exposes manifest introspection and validation as a first-class deterministic surface
- **Bootstrap/package flow hardened**: `scripts/harness-bootstrap.sh` now copies manifest-declared adapter files instead of relying on shell hardcoded copy lists, and still generates the target repo profile from the manifest-declared bootstrap template
- **Sync status strengthened**: `scripts/harness-sync.sh status` now reports adapter identity alignment in addition to manifest version alignment
- **Doctor upgraded**: `scripts/harness-doctor.sh` now includes adapter/package validation and also fixes the misplaced meeting validation logic that previously lived under the active-RC loop
- **Smoke contract upgraded**: core/release/bootstrap smoke now validate adapter integrity explicitly, and `scripts/harness-all-smoke-test.sh` now starts with adapter validation before broader lifecycle smoke
- **Adapter version advanced**: the reusable package line is now `0.5.0-prototype`
- **Validation**: `scripts/harness-adapter.sh validate`, `scripts/harness-doctor.sh`, and `scripts/harness-all-smoke-test.sh` all pass cleanly; full smoke remains green end-to-end

## Active Task / What's In Progress

**Leader-led harness runtime is now delivered. Productization and naming tasks remain:**

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 1 | ✅ | Research protocol / control plane / authority / cybernetics synthesis |
| Phase 2 | ✅ | Build repository-local harness definition layer (`harness/`) |
| Phase 3 | ✅ | Implement runtime substrate (`scripts/harness_*`) |
| Phase 4 | ✅ | Validate `home` / `doctor` / `claim` / `drive` / `meetings` / `maintenance` |
| Phase 5 | ✅ | Package the harness formally as a reusable Codex-native adapter/system |
| Phase 6 | ✅ | Split reusable adapter/bootstrap flow from repo-local proving ground |
| Phase 7 | 🔄 | Decide final product naming, public package contract, richer ownership/authority adjudication, and future Codex-native multi-agent integration path |

## Next Session Should Start With

1. Read this file
2. **Harness next work**:
   - A. Decide whether to formalize this under a new product name or keep `harness` as the prototype codename
   - B. Decide whether `harness/adapter-manifest.yml` should become the public package spec or stay repo-private while naming is unresolved
   - C. Deepen selector semantics beyond current normalized `kind:value` ownership model and moved-selector events
   - D. Deepen authority conflict resolution beyond level validation into actual high-vs-low authority adjudication
   - E. Decide how much farther to push council behavior toward `arch-review`-style automated debate versus keeping it as a durable packet/protocol
   - F. Decide whether future Codex-native subagent integration should map onto the teammate runtime directly or remain an optional adapter layer
   - H. If starting a fresh session, read `AGENTS.md`, `_project/CURRENT.md`, `_project/sessions/2026-03-28-A.md`, `_project/sessions/2026-03-28-B.md`, `_project/sessions/2026-03-28-C.md`, `_project/sessions/2026-03-29-A.md`, `_project/sessions/2026-03-29-B.md`, `_project/sessions/2026-03-29-C.md`, `_project/sessions/2026-03-29-D.md`, `harness/README.md`, `docs/research/harness-completion-context-pack-2026-03-28.md`, `docs/research/harness-gap-audit-2026-03-29.md`, `docs/research/harness-research-drift-audit-2026-03-29.md`, and inspect recent commits on `codex/harness-release-lifecycle`
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
- Harness is now installed and validated in Codex via `install-codex.sh`, can bootstrap into another repo, has manifest-driven adapter validation/packaging, first-class mission state, selector-aware ownership transfer, teammate runtime, and council packet semantics, but product naming and public contract decisions are still unresolved

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
