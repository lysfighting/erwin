# erwin — Project Context for Claude

> **Read this first. Then read `_project/CURRENT.md` for what's happening NOW.**

---

## What is this project?

**erwin** has two identities:

1. **Showcase** — Educational platform for understanding drive-skill (React web app in `showcase/`)
2. **Skills Factory** — Incubator for new Claude Code skills (`skills/` directory)

Skills pipeline: `_project/ideas/skills/` (ideas) → `.claude/commands/` (incubation) → `skills/` (released)

Install all erwin skills: `./install.sh`

See `_project/GOALS.md` for OKRs and north star.

---

## Repository Structure

```
erwin/
├── CLAUDE.md               ← YOU ARE HERE
├── CHANGELOG.md            ← Project version history
├── install.sh              ← Skills installer (symlinks to ~/.claude/commands/)
│
├── skills/                 ← Skills Factory (source of truth for all erwin skills)
│   ├── README.md           ← Skill index + contribution guide
│   ├── erwin-init/         ← /erwin-init skill (v1.0.0)
│   └── visual-test/        ← /visual-test skill (v1.0.0)
│
├── _project/               ← Project Intelligence Layer
│   ├── GOALS.md            ← Vision + OKRs (read when you lose direction)
│   ├── CURRENT.md          ← Living status (read every session start)
│   ├── ROADMAP.md          ← Milestone plan
│   ├── ideas/skills/       ← Skill ideas pipeline (Seed→Sprout→Plant→Fruit)
│   ├── decisions/          ← ADRs (NNNN-title.md)
│   ├── sessions/           ← AI conversation logs (YYYY-MM-DD.md)
│   ├── issues/             ← Bug/feature tracking (ISS-NNNN-title.md)
│   └── releases/           ← Narrative release notes
│
├── .claude/commands/       ← Skill incubation (project-level, pre-release)
│
├── docs/                   ← Technical Documentation
│   └── *.md                ← Architecture, analysis, research
│
├── showcase/         ← Code Artifact (React web app)
│   ├── src/sections/       ← 18 content sections
│   ├── src/components/     ← Shared components
│   └── CHANGELOG.md        ← App-level changelog
│
└── reference/              ← Source Material (READ-ONLY)
    └── drive-skill/        ← drive.md + pantheon + hooks + sounds
```

---

## Tech Stack

- **Frontend**: React 19.2, TypeScript, Vite 7.3, Tailwind 4.2, Framer Motion 12.34
- **Dev**: `cd showcase && npm run dev` (port 5173)

---

## Approved V2 Architecture

See `docs/drive-showcase-v2-architecture.md`. Four-Act structure:
- ACT I: 震撼 — ExecutionSimulator before/after
- ACT II: 理解 — 5 PrincipleCards + OriginTrace
- ACT III: 深度 — 18 existing sections (preserved)
- ACT IV: 创造 — SkillBuilder → export skill.md

3 new components: `ExecutionSimulator`, `OriginTrace`, `SkillBuilder`

---

## Session Protocol

**Start**: read `_project/CURRENT.md` → optionally read `_project/GOALS.md`
**End**: update `_project/CURRENT.md` → write session log → update CHANGELOG → git commit

---

## Key Files Quick Reference

| What | Where |
|------|-------|
| Current task | `_project/CURRENT.md` |
| Project vision | `_project/GOALS.md` |
| V2 design spec | `docs/drive-showcase-v2-architecture.md` |
| Skills incubator design | `docs/erwin-skills-incubator-architecture.md` |
| ADRs | `_project/decisions/` |
| Skill ideas pipeline | `_project/ideas/skills/` |
| Released skills | `skills/` |
| drive-skill source | `reference/drive-skill/commands/drive.md` |
| React app root | `showcase/src/App.tsx` |

---

## What NOT to Do

- Don't edit `reference/` — source material only
- Don't rely on session memory — always read CURRENT.md first
- Don't add new UI libraries (exception: d3-force already approved)
- Don't break existing 18 sections when adding V2 content
