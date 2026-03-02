# Mission Contract: SkillBuilder MVP
## Created: 2026-03-02

## Success Criteria (must ALL be true for MISSION_COMPLETE):
- [ ] SkillBuilder component renders in ACT IV, replacing the placeholder
- [ ] 5-step stepper form: Meta, Execution Model, Phases, Rules, Tools
- [ ] Each step has functional form inputs that update state
- [ ] Real-time skill.md preview (collapsible bottom panel)
- [ ] Download button generates and downloads a valid skill.md file
- [ ] Generated skill.md contains substantive prompt engineering content (not just skeleton)
- [ ] i18n support (EN/ZH) via existing LangContext/translations pattern
- [ ] Visual design matches existing showcase aesthetic (dark theme, gold accents, font-mono)
- [ ] Build passes with 0 TypeScript errors and 0 new lint warnings
- [ ] Mobile-responsive (stacked layout on small screens)

## Scope Boundaries:
- IN: SkillBuilder component, template generation, i18n, App.tsx integration
- OUT: DesignPatternLibrary, CaseStudies, code splitting, bundle optimization

## Key Constraints:
- No new dependencies (React + Tailwind + Framer Motion only)
- Follow existing component patterns (functional components, hooks, useLang)
- Keep component under 600 lines (split helpers if needed)

## Key Decisions Made:
- Stepper layout with collapsible bottom preview (not side-by-side)
- Template generates substantive prompt engineering content
- YAML frontmatter with allowed-tools and description
