# Feedback: erwin V2 Phase 1

## What Worked

- **Parallel swarm approach** worked well: stl-components and implementer-structure ran concurrently and produced high-quality work. Wall-clock time was roughly half what serial would be.
- **stl-components quality**: ExecutionSimulator and PrincipleCards were production-ready on first attempt. The `runKey` remount pattern for Framer Motion replaying was clever.
- **implementer-structure quality**: App.tsx restructure was clean and well-commented. ActIIIWrapper was complete and well-designed. The NavBar update was correct in intent.
- **Six-Step Discipline**: Having research-{slug}.md and plan-{slug}.md upfront made T7 (integration) mechanical — I just wired in what was already designed.

## What Didn't Work

- **NavBar artifact**: implementer-structure left the old `sectionIds` array alongside the new `actLinks` array, causing a TypeScript TS6133 error. Quick fix but worth noting — in multi-agent scenarios, the "remove old code" step can get missed.
- **AskUserQuestion responses not transmitting**: The Phase 0 interrogation was blocked because user question answers returned empty. Had to proceed on inference from CURRENT.md rather than explicit user confirmation.

## Surprises

- App.tsx restructure was already done by implementer-structure before I finished T1 (bootstrap). The parallel swarm was faster than expected.
- The `as const` on NavBar sectionIds keys required careful addition of new translation keys to avoid type errors — the translations.ts type system is tight.

## Lessons for Future Missions

- **For NavBar-style changes**: always explicitly instruct "remove the old array/code" — don't just say "replace with new". Background agents don't always clean up old code.
- **When AskUserQuestion fails**: fall back to CURRENT.md + GOALS.md for mission inference. The project documents are the source of truth.
- **Component quality from swarm**: the stl-components output was high quality without needing review cycles. For pure React component work, the swarm approach generates good code.
- **Team topology for this scale** (8 tasks, 2 domains): flat 2-teammate + Team Lead is the right structure. No need for STL hierarchy at this scale.
