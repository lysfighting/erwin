# erwin — Project Context for Claude

> **Read this first at the start of every session.**
> After reading this, read `docs/CURRENT_STATE.md` for what's happening NOW.

---

## What is this project?

**erwin** is an educational web platform that teaches people to understand, internalize, and
create their own [drive-skill](https://github.com/Zzzzz-0/drive-skill.git)-style AI agent
orchestration skills.

Primary goal: **傻瓜进来看一眼就知道这个 skill 有什么用，并能创建自己的版本。**

---

## Repository Structure

```
erwin/
├── CLAUDE.md                  ← YOU ARE HERE (read first)
├── CHANGELOG.md               ← project version history
├── drive-showcase/            ← React/Vite/TS web app (main work area)
│   ├── src/
│   │   ├── sections/          ← 18 content sections
│   │   ├── components/        ← shared components
│   │   ├── contexts/          ← LangContext, AnnotationContext
│   │   └── i18n/translations.ts
│   └── CHANGELOG.md           ← drive-showcase version history
├── reference-skills/
│   └── drive-skill/           ← source material (don't edit)
│       ├── commands/drive.md  ← ~2000 lines, the heart of everything
│       └── pantheon/          ← 1000+ thinkers knowledge base
└── docs/
    ├── CURRENT_STATE.md       ← READ THIS for current task status
    ├── sessions/              ← conversation artifacts by date
    ├── decisions/             ← Architecture Decision Records (ADRs)
    ├── drive-skill-analysis.md
    ├── visual-enhancement-architecture.md
    └── drive-showcase-v2-architecture.md  ← V2 design (approved)
```

---

## Tech Stack

- **Frontend**: React 19.2, TypeScript, Vite 7.3, Tailwind CSS 4.2, Framer Motion 12.34
- **Dev**: `cd drive-showcase && npm run dev` (port 5173)
- **Build**: `cd drive-showcase && npm run build`

---

## Current Architecture Plan (V2)

The approved V2 design is in `docs/drive-showcase-v2-architecture.md`.

**Four-Act Structure**:
- ACT I: 震撼 (15s) — ExecutionSimulator comparison
- ACT II: 理解 (5min) — 5 PrincipleCards with embedded demos + OriginTrace
- ACT III: 深度 (optional) — all 18 existing sections (preserved)
- ACT IV: 创造 — SkillBuilder for creating your own skill.md

**3 new core components to build**:
1. `ExecutionSimulator` — stateful animation showing drive execution
2. `OriginTrace` — force-directed graph tracing design philosophy origins
3. `SkillBuilder` — fill-in builder that exports skill.md

---

## Workflow Rules

1. **Before coding**: Read `docs/CURRENT_STATE.md` to know what's in progress
2. **After decisions**: Write an ADR in `docs/decisions/`
3. **After each session**: Update `docs/CURRENT_STATE.md` and `docs/sessions/YYYY-MM-DD.md`
4. **After features**: Update `drive-showcase/CHANGELOG.md` + root `CHANGELOG.md`
5. **Commit convention**: `type(scope): description` (feat/fix/docs/refactor/chore)

---

## Key Files to Know

| File | Purpose |
|------|---------|
| `reference-skills/drive-skill/commands/drive.md` | Source truth for all drive concepts (~2000 lines) |
| `drive-showcase/src/App.tsx` | Root component, section ordering |
| `drive-showcase/src/i18n/translations.ts` | All text content (EN + ZH) |
| `docs/CURRENT_STATE.md` | Current session's task state |
| `docs/drive-showcase-v2-architecture.md` | V2 approved design |

---

## What NOT to do

- Don't edit `reference-skills/drive-skill/` — it's source material only
- Don't rely on session memory for continuity — always read CURRENT_STATE.md first
- Don't create new UI libraries — stick to existing stack
- Don't break existing 18 sections when adding V2 content
