# Showcase V3 Architecture — Hub + Multi-Showcase

> **Status**: Archived — hub concept not implemented. Current showcase remains single-page V2 architecture (`drive-showcase-v2-architecture.md`).

> erwin showcase 从单体 drive-skill 页面 → 多 showcase 导航枢纽

---

## 1. Executive Summary

**What**: 将 showcase 从单页 drive-skill 展示重构为多页导航枢纽,包含 reference showcase (drive-skill v3.1) 和 skill showcases (erwin 自孵化的 3 个 skill)。

**Why**: drive-skill 已更新到 v3.1 (2821 行, 4 个 commands),现有 showcase 仅覆盖 ~40% 内容。erwin 已孵化 3 个 skill (erwin-init, visual-test, teamwork) 但无展示入口。

**How**: 引入 React Router v7, 首页为 Hub 导航, 每个 showcase 为独立路由页面。保留现有 17 sections 并更新至 v3.1, 新增 6 个 section 覆盖关键缺口。

---

## 2. System Diagram

### 2.1 Routing Topology

```
/                               ← Hub (导航枢纽)
│
├── /reference/drive-skill      ← Drive Skill Showcase
│   ├── ACT I:  震撼 (Hero + ExecutionSimulator)
│   ├── ACT II: 理解 (Principles + NEW: PhaseΩ, ChallengeGates)
│   ├── ACT III: 深度 (17 existing + 6 new sections)
│   └── ACT IV: 创造 (SkillBuilder)
│
├── /skills/erwin-init          ← erwin-init Showcase
├── /skills/visual-test         ← visual-test Showcase
└── /skills/teamwork            ← teamwork Showcase
```

### 2.2 Component Tree

```
<BrowserRouter>
  <LangProvider>
    <AnnotationProvider>
      <Routes>
        <Route element={<Layout />}>         ← NavBar + Footer shell
          <Route index element={<HubPage />} />
          <Route path="reference/drive-skill" element={<DriveSkillPage />} />
          <Route path="skills/erwin-init" element={<ErwinInitPage />} />
          <Route path="skills/visual-test" element={<VisualTestPage />} />
          <Route path="skills/teamwork" element={<TeamworkPage />} />
        </Route>
      </Routes>
    </AnnotationProvider>
  </LangProvider>
</BrowserRouter>
```

---

## 3. Folder Structure

### 3.1 Target Structure

```
showcase/src/
├── App.tsx                          # BrowserRouter + providers
│
├── app/
│   ├── Layout.tsx                   # Shared shell: NavBar + <Outlet /> + Footer
│   └── routes.tsx                   # Route definitions (lazy imports)
│
├── pages/
│   ├── Hub/
│   │   ├── HubPage.tsx              # 导航枢纽
│   │   ├── ShowcaseCard.tsx          # Hub card component
│   │   └── data.ts                  # Showcase metadata
│   │
│   ├── DriveSkill/
│   │   ├── DriveSkillPage.tsx       # 4-Act structure (原 AppInner)
│   │   └── sections/               # 所有 drive-skill sections
│   │       ├── HeroSection.tsx
│   │       ├── ... (existing 17, moved from src/sections/)
│   │       ├── PhaseOmegaSection.tsx       # NEW
│   │       ├── ChallengeGatesSection.tsx   # NEW
│   │       ├── TeamAssemblySection.tsx     # NEW
│   │       ├── BattlePlanSection.tsx       # NEW
│   │       ├── FilePersistenceSection.tsx  # NEW
│   │       └── ImproveUserSection.tsx      # NEW
│   │
│   ├── SkillErwinInit/
│   │   └── ErwinInitPage.tsx
│   │
│   ├── SkillVisualTest/
│   │   └── VisualTestPage.tsx
│   │
│   └── SkillTeamwork/
│       └── TeamworkPage.tsx
│
├── components/                      # 共享组件 (不变)
│   ├── NavBar.tsx                   # 更新: 支持 Hub/Showcase 两种模式
│   ├── SectionDivider.tsx
│   ├── ActIIIWrapper.tsx
│   ├── AnnotationsPanel.tsx
│   ├── ExecutionSimulator.tsx
│   ├── PrincipleCards.tsx
│   ├── OriginTrace.tsx
│   ├── SkillBuilder.tsx
│   └── SkillShowcaseTemplate.tsx    # NEW: 统一 skill showcase 模板
│
├── contexts/
│   ├── LangContext.tsx              # 不变
│   └── AnnotationContext.tsx        # 不变
│
└── i18n/
    └── translations.ts             # 扩展: 新增 Hub + skill showcase 翻译 key
```

### 3.2 Migration from Current Structure

| Before | After | Action |
|--------|-------|--------|
| `src/App.tsx` (AppInner) | `pages/DriveSkill/DriveSkillPage.tsx` | Move content |
| `src/sections/*.tsx` | `pages/DriveSkill/sections/*.tsx` | Move directory |
| `src/App.tsx` (App) | `src/App.tsx` (Router + Providers) | Rewrite |
| — | `app/Layout.tsx` | New |
| — | `app/routes.tsx` | New |
| — | `pages/Hub/HubPage.tsx` | New |

---

## 4. Key Components

### 4.1 App.tsx (Rewritten)

```tsx
import { BrowserRouter } from 'react-router'
import { LangProvider } from './contexts/LangContext'
import { AnnotationProvider } from './contexts/AnnotationContext'
import { AppRoutes } from './app/routes'

export default function App() {
  return (
    <BrowserRouter>
      <LangProvider>
        <AnnotationProvider>
          <AppRoutes />
        </AnnotationProvider>
      </LangProvider>
    </BrowserRouter>
  )
}
```

### 4.2 Layout.tsx

```tsx
import { Outlet, useLocation } from 'react-router'
import { NavBar } from '../components/NavBar'
import { AnnotationsPanel } from '../components/AnnotationsPanel'

export function Layout() {
  const { pathname } = useLocation()
  const isHub = pathname === '/'

  return (
    <div className="min-h-screen bg-bg-dark">
      <NavBar mode={isHub ? 'hub' : 'showcase'} />
      <main>
        <Outlet />
      </main>
      {!isHub && <AnnotationsPanel />}
      <footer className="py-16 text-center text-text-muted border-t border-border">
        <p className="font-mono text-sm">erwin — Skill Factory & Reference Platform</p>
      </footer>
    </div>
  )
}
```

### 4.3 NavBar Updates

NavBar 需要两种模式:

| Mode | Behavior |
|------|----------|
| `hub` | erwin logo, no act navigation, no scroll progress |
| `showcase` | Back-to-hub link, act navigation (drive-skill only), scroll progress, language toggle |

```tsx
interface NavBarProps {
  mode: 'hub' | 'showcase'
}

// hub mode: erwin branding + language toggle only
// showcase mode: ← Hub | DRIVE SKILL | act-i act-ii act-iii act-iv | EN 中
```

### 4.4 Hub Page

```tsx
// pages/Hub/data.ts
export interface ShowcaseEntry {
  id: string
  name: string
  version: string
  tagline: string        // one-line pitch
  description: string    // 2-3 sentence description
  category: 'reference' | 'skill'
  path: string
  commands?: string[]    // sub-commands for drive-skill
  status: 'live' | 'coming-soon'
}

export const showcases: ShowcaseEntry[] = [
  {
    id: 'drive-skill',
    name: 'drive-skill',
    version: 'v3.1.0',
    tagline: 'High-drive autonomous execution with hierarchical team swarm',
    description: 'The core reference. Transforms Claude Code into a Team Lead...',
    category: 'reference',
    path: '/reference/drive-skill',
    commands: ['drive', 'self-drive', 'architect', 'improve-user'],
    status: 'live',
  },
  {
    id: 'erwin-init',
    name: 'erwin-init',
    version: 'v2.0.0',
    tagline: 'Cognitive framework installer — zero questions, instant setup',
    description: 'Initialize version management, session protocols, ADR system...',
    category: 'skill',
    path: '/skills/erwin-init',
    status: 'live',
  },
  {
    id: 'visual-test',
    name: 'visual-test',
    version: 'v1.0.0',
    tagline: 'Screenshot + visual verification for Claude Code',
    description: 'Take screenshots and visually verify UI changes...',
    category: 'skill',
    path: '/skills/visual-test',
    status: 'live',
  },
  {
    id: 'teamwork',
    name: 'teamwork',
    version: 'v2.1.0',
    tagline: 'GitHub-first team coordination — 4 skills, one workflow',
    description: '/team, /team-claim, /team-drive, /team-ship...',
    category: 'skill',
    path: '/skills/teamwork',
    status: 'live',
  },
]
```

### 4.5 Skill Showcase Template

统一模板,所有 erwin skill showcase 使用:

```tsx
// components/SkillShowcaseTemplate.tsx
interface SkillShowcaseTemplateProps {
  name: string
  version: string
  tagline: string
  installCmd: string
  children: {
    hero: React.ReactNode        // visual identity
    concept: React.ReactNode     // what + why
    walkthrough: React.ReactNode // how it works (step-by-step)
    interactive?: React.ReactNode // try it / demo
  }
}

// Structure:
// ┌─ Hero ──────────────────────────────┐
// │ name + version + tagline + install  │
// │ {hero}                              │
// ├─ Core Concept ──────────────────────┤
// │ {concept}                           │
// ├─ How It Works ──────────────────────┤
// │ {walkthrough}                       │
// ├─ Try It (optional) ─────────────────┤
// │ {interactive}                       │
// └─────────────────────────────────────┘
```

---

## 5. Drive Skill Showcase — Section Update Plan

### 5.1 Gap Analysis Summary

| Coverage | % | Notes |
|----------|---|-------|
| Philosophy + Sage | 70% | Missing Challenge Gates, Dual Wheel |
| Execution Phases | 30% | Only Phase 0; missing Ω, T, 1, Loop details |
| File Persistence | 5% | Critical gap |
| Decision Protocols | 20% | Missing 3 Challenge Gates + Rebuttal Format |
| Swarm/Hierarchy | 50% | Missing STL details, Dynamic Restructuring |
| Sub-skills coverage | 15% | /improve-user zero coverage |

### 5.2 Section Inventory

#### Existing Sections — Update Plan

| # | Section | Update Needed | Key Changes |
|---|---------|---------------|-------------|
| 1 | HeroSection | Minor | Version badge → v3.1, update tagline |
| 2 | ClaudeMechanicsSection | Moderate | Add hooks, sounds, messaging, self-drive config |
| 3 | PhilosophySection | Moderate | Link to Challenge Gates, add Dual Wheel concept |
| 4 | ArchitectureSection | Moderate | Update to reflect 4 commands, not just drive |
| 5 | ExecutionPhasesSection | **Major** | Expand from Phase 0 only → all 6 phases (Ω→0→T→1→Loop→Complete) |
| 6 | SwarmTopologySection | **Major** | Add STL hierarchy, Dynamic Restructuring ops, hierarchy decision tree |
| 7 | WaveExecutionSection | Moderate | Update wave examples, add critical path analysis |
| 8 | SixStepSection | Moderate | Add file structure (contract/plan/state/research/feedback) |
| 9 | DeepReasoningSection | Minor | Add Reasoning Rhythm, link to Pantheon |
| 10 | DirectionValidationSection | Minor | Verify alignment with v3.1 |
| 11 | PantheonSection | Minor | Verify thinker mappings match v3.1 |
| 12 | SelfAdversarialReviewSection | Minor | Verify 6 axes match v3.1 |
| 13 | AntiPatternsSection | Moderate | Add 9 hierarchy-specific anti-patterns from v2.0 |
| 14 | SelfDriveSection | Moderate | Update with v3.1 self-drive content |
| 15 | ArchitectSection | Moderate | Update to reflect v3.1 architect skill |
| 16 | PromptEngineeringSection | Minor | Cross-reference to Challenge Gates |
| 17 | ReferencesSection | Minor | Update citations |

#### New Sections

| # | Section | drive.md Source | Description |
|---|---------|-----------------|-------------|
| N1 | **PhaseOmegaSection** | lines 348-464 | 本质洞察: bone vs flesh, decided vs open, blast radius, Essence Statement |
| N2 | **ChallengeGatesSection** | lines 42-163 | 3 mandatory gates: Essence, Direction, Contract + Rebuttal Format + Anti-Compliance Test |
| N3 | **TeamAssemblySection** | lines 907-1138 | Phase T: team creation, hierarchy decision tree, STL contracts, spawn templates |
| N4 | **BattlePlanSection** | lines 2079-2270 | Phase 1: task decomposition, parallelism analysis, wave mapping, critical path, dispatch |
| N5 | **FilePersistenceSection** | lines 1720-2048 | .claude/drive/ directory: contract.md (immutable), plan.md (blueprint), state.md (checkpoint), research/, feedback.md |
| N6 | **ImproveUserSection** | improve-user.md | /improve-user: cognitive coaching, archaeological dig, pattern extraction, growth artifact |

### 5.3 Updated 4-Act Structure

```
ACT I: 震撼 (0-15 seconds)
  ├── HeroSection (updated badge)
  └── ExecutionSimulator (keep)

ACT II: 理解 (1-10 minutes)
  ├── ClaudeMechanicsSection (updated)
  ├── PrincipleCards (keep)
  ├── OriginTrace (keep)
  ├── PhaseOmegaSection          ← NEW
  └── ChallengeGatesSection      ← NEW

ACT III: 深度 (10-30 minutes) — inside ActIIIWrapper
  ├── ArchitectureSection (updated)
  ├── ExecutionPhasesSection (MAJOR update: all 6 phases)
  ├── TeamAssemblySection        ← NEW
  ├── SwarmTopologySection (MAJOR update: STL hierarchy)
  ├── BattlePlanSection          ← NEW
  ├── WaveExecutionSection (updated)
  ├── SixStepSection (updated)
  ├── FilePersistenceSection     ← NEW
  ├── PhilosophySection (updated)
  ├── DeepReasoningSection (updated)
  ├── DirectionValidationSection (keep)
  ├── PantheonSection (keep)
  ├── SelfAdversarialReviewSection (keep)
  ├── AntiPatternsSection (updated)
  ├── SelfDriveSection (updated)
  ├── ArchitectSection (updated)
  ├── ImproveUserSection         ← NEW
  ├── PromptEngineeringSection (keep)
  └── ReferencesSection (updated)

ACT IV: 创造 (30+ minutes)
  ├── SkillBuilder (keep)
  ├── DesignPatternLibrary (placeholder)
  └── CaseStudies (placeholder)
```

Total: 23 sections (17 existing + 6 new)

---

## 6. Key Decisions

| Decision | Choice | Rejected | Why |
|----------|--------|----------|-----|
| Routing library | React Router v7 (SPA mode) | Hash router, conditional render, Vite MPA | Standard, supports lazy loading, preserves URLs |
| Folder structure | Feature-based (pages/) | Keep flat sections/ | Scales with N showcases, clear ownership |
| Existing sections | Preserve + update | Full rewrite, freeze | Leverages existing work (animations, i18n), avoids throwaway |
| NavBar approach | Dual-mode (hub/showcase) | Separate NavBars | Single component, less duplication |
| Skill showcase format | Unified template | Per-skill custom | Consistency, lower effort per new skill |
| New sections count | 6 | Could be more (10+) | Covers Tier 1 gaps; remaining gaps can be addressed by updating existing sections |

---

## 7. Implementation Roadmap

### Phase A: Infrastructure (Router + Hub + Restructure)

1. `npm install react-router` in showcase/
2. Create `app/routes.tsx`, `app/Layout.tsx`
3. Rewrite `App.tsx` → Router + Providers
4. Create `pages/Hub/` with HubPage, ShowcaseCard, data
5. Move `sections/` → `pages/DriveSkill/sections/`
6. Create `pages/DriveSkill/DriveSkillPage.tsx` (content from AppInner)
7. Update NavBar for dual-mode (hub/showcase)
8. Update `vite.config.ts` manualChunks for route-based splitting
9. Verify: existing drive-skill showcase works at `/reference/drive-skill`
10. Verify: Hub page renders at `/`

### Phase B: Drive Skill Content Updates

1. Update ExecutionPhasesSection → all 6 phases with visual flow
2. Update SwarmTopologySection → STL hierarchy + Dynamic Restructuring
3. Update AntiPatternsSection → add hierarchy anti-patterns
4. Update remaining sections (minor updates per table above)
5. Visual verify each updated section

### Phase C: New Drive Skill Sections

1. PhaseOmegaSection — interactive bone/flesh/blast-radius decomposition
2. ChallengeGatesSection — 3 gates with example output visualization
3. TeamAssemblySection — hierarchy decision tree + STL spawn visual
4. BattlePlanSection — wave decomposition + critical path diagram
5. FilePersistenceSection — .claude/drive/ directory tree + file previews
6. ImproveUserSection — cognitive coaching flow diagram

### Phase D: Skill Showcases

1. Create `SkillShowcaseTemplate` component
2. ErwinInitPage — using template
3. VisualTestPage — using template
4. TeamworkPage — using template (covers 4 sub-skills)

---

## 8. Risk Register

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| React Router breaks existing hash navigation (#act-i) | High | Medium | Phase A includes migration of hash links → React Router scroll anchors |
| 6 new sections bloat ACT III further | Medium | Low | ActIIIWrapper already collapses; new sections are lazy-loaded |
| Skill showcases feel empty vs drive-skill | Medium | Medium | Template enforces minimum structure; "quality over quantity" |
| Content updates lag behind reference changes | Ongoing | Medium | Establish reference-sync protocol in CLAUDE.md |

---

## 9. Appendix: Dependencies

### New Dependencies

| Package | Version | Size | Why |
|---------|---------|------|-----|
| `react-router` | ^7.x | ~40KB gzipped | Client-side routing |

No other new dependencies required.

### Removed Dependencies

None.
