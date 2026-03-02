# Roadmap — erwin

> Milestone plan. Update when phases shift. Reference GOALS.md for the why.

---

## Milestone 1: showcase V2 (3 weeks)

**Goal**: Interactive showcase that makes drive-skill understandable to anyone
**Target**: v1.0.0 of the interactive platform

### Week 1: Structure + Simulator
- [ ] Refactor App.tsx for four-act navigation
- [ ] ExecutionSimulator MVP (3 scenarios, basic controls)
- [ ] PrincipleCards (5 cards with static demos)

### Week 2: Deep Interactivity
- [ ] ExecutionSimulator full (message flow, error injection)
- [ ] OriginTrace component (d3-force philosophy graph)
- [ ] SwarmTopology + WaveExecution enhancements

### Week 3: Creation Path
- [ ] SkillBuilder (form → skill.md download)
- [ ] Design pattern library
- [ ] Polish + testing

---

## Milestone 2: First Original Skill (TBD)

**Goal**: Create one new skill from the idea garden
**Depends on**: Milestone 1 complete (showcase teaches how to build skills)
**Candidates**: See `_project/ideas/skills/`

---

## Milestone 3: Project Template Extraction (TBD)

**Goal**: Extract _project/ structure as a reusable template
**Deliverable**: `reference/project-template/` + install script
**Candidate skill**: `/new-project` skill

---

## Version Plan

| Version | Description | Status |
|---------|-------------|--------|
| 0.1.0 | Initial showcase (18 sections, v1 content) | ✅ Shipped |
| 0.2.0 | Version management system | ✅ Shipped |
| 0.3.0 | Project structure reorganization | 🔄 In progress |
| 1.0.0 | V2 showcase with four-act structure | ⬜ Planned |
| 2.0.0 | First original skill + expanded showcase | ⬜ Future |
