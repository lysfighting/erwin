# Current State

> **Update at END of every session. Read at START of every session.**

---

## Project Status

- **Version**: 0.5.0 (in progress)
- **Last Updated**: 2026-03-03
- **Last Session**: [2026-03-03-O](sessions/2026-03-03-O.md) (Session O)

## Current Phase

**Phase**: Teamwork v2.4 — stability hardening + bidirectional sync

## What Was Just Completed

**Self-Drive Session O** — 6 items prioritized and executed as Teamwork v2.4.0:

1. **tw-config.sh** — New unified config parser script (`scripts/tw-config.sh`). Replaces 15+ fragile sed chains across all 5 skill files. Handles colon-containing values, quoted strings, nested YAML. Python3-based.
2. **Placeholder normalization** — Single canonical convention: `{issue}`, `{pr}`, `{user}`, `{branch}` across all 5 skills.
3. **team-ship pseudocode implemented** — tech-lead extraction, commit-type detection, base_branch reading were pseudocode before. Now real bash.
4. **Critical command verifications** — `gh api user` fatal check, `git checkout` branch existence, TEST_CMD empty check.
5. **Issue bidirectional sync** — `team-drive` writes sub-task completions back to GitHub Issue body during execution.
6. **Version from Milestone** — `team-claim` reads milestone from Issue JSON, not config `versions.current`.

Commit: `a0de18d feat(skills): teamwork v2.4.0 — stability hardening`

## Active Task / What's In Progress

None. Session O complete.

## Next Session Should Start With

1. Read this file
2. **E2E validate** teamwork v2.4 on `../vi_agent` — run full flow: /team → /team-issue → /team-claim → /team-drive → /team-ship
3. OR continue V2 Phase 3 showcase tasks (DesignPatternLibrary, CaseStudies)
4. **Future missions** (from Session O discussion):
   - `/team-release` skill (close Milestone → git tag → `gh release create --generate-notes`)
   - `/team-drive` wrapping full `/drive` (thin wrapper, skip Phase 0)

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
- [x] Teamwork v2.0 — 4 GitHub-first skills
- [x] Teamwork v2.1 — 6 vi-agent features + 13 bug fixes
- [x] Teamwork v2.3 — status fix + version management + format standardization + scripts
- [ ] Teamwork E2E validation on vi_agent

## Known Issues / Blockers

- OriginTrace ResizeObserver fires 400 force ticks per resize — needs debounce
- Placeholder naming inconsistent across teamwork skill files ({N} vs {ISSUE_NUMBER} vs {issue})
- Teamwork v2.3: milestone creation (if it doesn't exist) needs manual step before /team-issue

## File Locations

| What | Where |
|------|-------|
| Project goals | `_project/GOALS.md` |
| Idea garden | `_project/ideas/` |
| ADRs | `_project/decisions/` |
| V2 architecture | `docs/drive-showcase-v2-architecture.md` |
| Teamwork architecture | `docs/teamwork-v2-architecture.md` |
| Teamwork skills | `skills/teamwork/` |
| Teamwork scripts | `skills/teamwork/scripts/` |
| App root | `showcase/src/App.tsx` |
