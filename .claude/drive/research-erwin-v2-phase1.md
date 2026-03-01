# Research: erwin V2 Phase 1

## Codebase Understanding

### Project Structure
- Root: `/Users/liyasong/Project/collov/erwin/`
- App: `drive-showcase/` (React 19.2 + TypeScript + Vite 7.3 + Tailwind 4.2 + Framer Motion 12.34)
- Dev server: `cd drive-showcase && npm run dev` (port 5173)

### Current App Structure
- `App.tsx` — linear flat structure with 18 sections + NavBar
- `src/sections/` — 18 section components (HeroSection through ReferencesSection)
- `src/components/` — NavBar, SectionDivider, SectionWrapper, AnnotationsPanel, Annotatable
- `src/contexts/` — LangContext, AnnotationContext
- `src/i18n/translations.ts` — i18n support

### 18 Existing Sections (all preserved)
1. HeroSection → stays in ACT I
2. DemoSection → replaced by ExecutionSimulator in ACT I
3. ClaudeMechanicsSection → moves to ACT II
4. ArchitectureSection → ACT III
5. ExecutionPhasesSection → ACT III
6. SwarmTopologySection → ACT III
7. WaveExecutionSection → ACT III
8. SixStepSection → ACT III
9. PhilosophySection → ACT III
10. DeepReasoningSection → ACT III
11. DirectionValidationSection → ACT III
12. PantheonSection → ACT III
13. SelfAdversarialReviewSection → ACT III
14. AntiPatternsSection → ACT III
15. SelfDriveSection → ACT III
16. ArchitectSection → ACT III
17. PromptEngineeringSection → ACT III
18. ReferencesSection → ACT III

### New Components to Build
- `ExecutionSimulator.tsx` — ACT I before/after comparison
- `PrincipleCards.tsx` — ACT II 5 principle cards
- `ActIIIWrapper.tsx` — collapsible wrapper for 18 deep sections
- (Future: OriginTrace, SkillBuilder, DesignPatternLibrary)

## Key Findings
- No d3-force dependency yet (needed for OriginTrace later, not now)
- Tailwind 4.2 with CSS variables (bg-bg-dark, text-text-primary, etc.)
- Framer Motion 12.34 available for animations
- NavBar.tsx has existing lang toggle + annotations button

## Dependencies & Constraints
- React 19.2, TypeScript, Vite 7.3, Tailwind 4.2, Framer Motion 12.34
- No new UI libraries
- Don't break existing 18 sections
- Keep TypeScript valid
