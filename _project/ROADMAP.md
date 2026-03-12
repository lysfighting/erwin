# Roadmap — erwin

> Milestone plan. Update when phases shift. Reference GOALS.md for the why.

---

## Milestone 1: showcase V2 (largely complete)

**Goal**: Interactive showcase that makes drive-skill understandable to anyone
**Target**: v1.0.0 of the interactive platform

### Completed
- [x] Refactor App.tsx for four-act navigation
- [x] ExecutionSimulator MVP (3 scenarios, message flow, error injection)
- [x] PrincipleCards (5 cards with static demos)
- [x] OriginTrace component (d3-force philosophy graph)
- [x] SwarmTopology + WaveExecution enhancements
- [x] SkillBuilder (form → skill.md download)
- [x] Bundle code splitting (627KB → 221KB)

### Remaining for v1.0.0
- [ ] DesignPatternLibrary (3-5 reusable pattern cards)
- [ ] CaseStudies (2-3 skill variant examples)
- [ ] Polish + testing

---

## Milestone 2: Teamwork Skills Factory (active)

**Goal**: AI-native team coordination skills — GitHub-first, zero infrastructure
**Status**: v3.9.2 shipped (34+ versions since v2.0)

### Completed
- [x] 6 core skills: /team, /team-issue, /team-claim, /team-drive, /team-ship, /team-rc
- [x] Script extraction: tw-git.sh, tw-config.sh, tw-contract.sh, tw-github.sh, tw-bump-version.sh
- [x] Visual Encoding Standard + full compliance audit
- [x] /team doctor (12 checks + fix mode)
- [x] /team wrap + /team auto (shortcuts)
- [x] /team sync + /team share (cross-MC collaboration)
- [x] /team-rc fix + hotfix (complete RC lifecycle)
- [x] Documentation: manual, onboarding, CHANGELOG, VE standard

### Remaining
- [ ] Deploy workflow (deploy.yml, staging gate)
- [ ] Issue version tracking (found-in labels, version labels)

---

## Milestone 3: Project Template Extraction (TBD)

**Goal**: Extract _project/ structure as a reusable template
**Deliverable**: `reference/project-template/` + install script
**Candidate skill**: `/erwin-init` skill

---

## Version Plan

| Version | Description | Status |
|---------|-------------|--------|
| 0.1.0 | Initial showcase (18 sections, v1 content) | Shipped |
| 0.2.0 | Version management system | Shipped |
| 0.3.0 | Project structure reorganization | Shipped |
| 0.4.0 | V2 Phase 1 (ExecutionSimulator, PrincipleCards) | Shipped |
| 0.5.0 | V2 Phase 2 (OriginTrace, code splitting, polish) | Shipped |
| 1.0.0 | V2 complete (DesignPatternLibrary, CaseStudies) | In progress |
| Teamwork 3.9.x | AI-native team coordination | Active |
