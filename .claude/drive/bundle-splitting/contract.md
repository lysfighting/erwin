# Mission Contract: Bundle Code Splitting
## Created: 2026-03-02

## Success Criteria (must ALL be true for MISSION_COMPLETE):
□ Initial bundle chunk < 300KB (down from 627KB)
□ ACT III 15 sections load lazily (only when ActIIIWrapper expanded)
□ SkillBuilder loads lazily (only when ACT IV scrolls into view)
□ Vite 500KB chunk warning eliminated
□ Build passes with 0 TS errors
□ App functions correctly — all sections render when accessed

## Scope Boundaries:
- IN: React.lazy() + Suspense for ACT III sections and SkillBuilder
- OUT: OriginTrace splitting (d3-force, ACT II — always visible), further optimization, test writing

## Key Constraints:
- Named exports require `.then(m => ({ default: m.X }))` wrapper for React.lazy()
- ActIIIWrapper already conditionally renders children — lazy imports won't fire until expanded
- No new dependencies

## Key Decisions Made:
- Solo execution (≤5 tasks, single domain)
- Lazy-load at section granularity, not sub-component level
