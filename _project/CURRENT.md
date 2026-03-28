# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.6.0 (in progress)
- **Last Updated**: 2026-03-28
- **Last Session**: 2026-03-28 (Codex-native orchestration harness prototype — leader-led teammate runtime realized)

## Current Phase

**Phase**: Codex-native orchestration harness prototype

## What Was Just Completed

**Session C** (2026-03-28) — Harness prototype crossed from leader-centric control plane into first-class teammate runtime

- **Teammate runtime added**: `scripts/harness-teammates.sh` plus `harness_state.py teammates ...` now create bounded delegated runs with scope, authority refs, forbidden sources, reporting expectations, and parent/leader linkage
- **Governed lateral communication added**: teammate-to-teammate communication is now explicit and structured instead of implicit chat, using sibling-approved signal flow and leader-approved linkage
- **Council protocol strengthened**: meetings now track `phase`, `participants`, and `positions`, so council packets can represent independent review and cross-examination instead of only final findings
- **Home/doctor upgraded**: `home` now surfaces active teammate and coordination signal counts; `doctor` now validates teammate assignment quality plus meeting participant/phase/position health
- **Smoke coverage expanded**: `scripts/harness-teammate-smoke-test.sh` now validates delegation, report emission, sibling communication, meeting escalation, participant rostering, and council positions end-to-end
- **Adapter/package contract advanced**: teammate surface and teammate smoke are now part of `harness/adapter-manifest.yml`, and the reusable adapter line advanced to `0.6.0-prototype`
- **Local runtime cleaned**: legacy smoke teammate runs/signals were retired and reconciled so `scripts/harness-doctor.sh` returns cleanly again in the real repository runtime
- **Validation**: `python3 -m py_compile scripts/harness_state.py`, `scripts/harness-smoke-test.sh`, `scripts/harness-teammate-smoke-test.sh`, `scripts/harness-release-smoke-test.sh`, `scripts/harness-bootstrap-smoke-test.sh`, `scripts/harness-all-smoke-test.sh`, and `scripts/harness-doctor.sh` all pass cleanly

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
| Phase 7 | 🔄 | Decide final product naming, public package contract, and future Codex-native multi-agent integration path |

## Next Session Should Start With

1. Read this file
2. **Harness next work**:
   - A. Decide whether to formalize this under a new product name or keep `harness` as the prototype codename
   - B. Decide whether `harness/adapter-manifest.yml` should become the public package spec or stay repo-private while naming is unresolved
   - C. Decide how much farther to push council behavior toward `arch-review`-style automated debate versus keeping it as a durable packet/protocol
   - D. Decide whether future Codex-native subagent integration should map onto the teammate runtime directly or remain an optional adapter layer
   - E. If starting a fresh session, read `AGENTS.md`, `_project/CURRENT.md`, `_project/sessions/2026-03-28-A.md`, `_project/sessions/2026-03-28-B.md`, `_project/sessions/2026-03-28-C.md`, `harness/README.md`, `docs/research/harness-completion-context-pack-2026-03-28.md`, and inspect recent commits on `codex/harness-release-lifecycle`
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
- Harness is now installed and validated in Codex via `install-codex.sh`, can bootstrap into another repo, has manifest-driven adapter validation/packaging, and now includes first-class teammate runtime plus council packet semantics, but product naming and public contract decisions are still unresolved

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
