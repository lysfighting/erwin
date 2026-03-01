# Plan: erwin V2 Phase 1

## Mission Contract Reference

### Success Criteria:
- [ ] App.tsx reorganized into Four-Act layout (ACT I/II/III/IV sections with nav anchors)
- [ ] NavBar updated with four Act anchors
- [ ] ExecutionSimulator MVP: animated before/after comparison, deployed in ACT I
- [ ] PrincipleCards: 5 cards with principle info, deployed in ACT II
- [ ] Existing 18 sections preserved in ACT III, collapsed by default
- [ ] Site builds and renders correctly, no broken imports
- [ ] "15-second shock" works: visitor sees before/after comparison

## Approach Overview

Phase 1 strategy: **structure first, then populate**.

1. Restructure App.tsx into four acts with proper section IDs and nav anchors
2. Update NavBar to link to the four acts
3. Build ActIIIWrapper to collapse the 18 existing sections
4. Build the two new content components (ExecutionSimulator, PrincipleCards)
5. Wire everything together

The existing 18 sections are preserved intact — we're adding navigation structure around them, not replacing them.

## Architecture / Design

```
App.tsx structure (V2):
├── NavBar (updated: 4 act anchors + existing controls)
│
├── ACT I: 震撼 (id="act-i")
│   ├── HeroSection (existing)
│   └── ExecutionSimulator (NEW: before/after animated comparison)
│
├── ACT II: 理解 (id="act-ii")
│   ├── ClaudeMechanicsSection (existing, moved here)
│   └── PrincipleCards (NEW: 5 cards)
│
├── ACT III: 深度 (id="act-iii")
│   └── ActIIIWrapper (NEW: collapsible)
│       ├── ArchitectureSection
│       ├── ExecutionPhasesSection
│       ├── ... (14 more existing sections)
│       └── ReferencesSection
│
└── ACT IV: 创造 (id="act-iv")
    └── Placeholder divs (SkillBuilder etc. — Phase 3)
```

## Key Design Decisions

- **Decision 1**: DemoSection replaced by ExecutionSimulator — because DemoSection is the logical slot for "shock" content
- **Decision 2**: ActIII collapsed by default — reduces overwhelm for new visitors
- **Decision 3**: ExecutionSimulator uses Framer Motion progress bars — matches existing animation library, no new deps
- **Decision 4**: Four act IDs: "act-i", "act-ii", "act-iii", "act-iv" — simple anchor targets

## Implementation Sequence

1. Scaffold working directory + verify build (Team Lead)
2. [PARALLEL] Restructure App.tsx + NavBar + ActIIIWrapper (implementer-structure)
3. [PARALLEL] ExecutionSimulator + PrincipleCards (stl-components)
4. Wire final imports into App.tsx (Team Lead)
5. Final verification (Team Lead)

## Status: APPROVED — executing
