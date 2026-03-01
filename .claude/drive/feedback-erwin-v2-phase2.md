# Feedback: erwin V2 Phase 2

## What Worked

- **Solo implementation** (no swarm needed): OriginTrace had no multi-domain parallelism — one component with clear spec. Flat solo was right-sized.
- **d3-force synchronous tick pattern**: Running `.stop().tick(400)` in `useMemo` gives deterministic layout without useEffect complexity. The positions are stable and computed before first render.
- **Color consistency**: Reusing PrincipleCards' group A-E accent colors in OriginTrace creates visual coherence — the two ACT II components feel like a system, not two independent features.
- **ResizeObserver for responsive SVG**: `width → height` derived as `width * 0.62` gives a clean aspect ratio that works at all breakpoints.
- **Type-safe d3-force**: The `SimulationNodeDatum` + `as SimNode` pattern for resolved links is the correct TypeScript idiom for d3-force.

## What Didn't Work

- **d3-force types not bundled**: Needed separate `@types/d3-force` install. The `package.json types` field was absent — wasn't caught in initial reconnaissance. Minor, but worth noting.
- **Chinese text inside circle sizing**: The 7-char '永不停止的循环' at fontSize=8 fits in r=34 circle (56px < 68px), but this is tight and system-font dependent. Not tested on Windows.

## Surprises

- Build jumped from 455 → 489 modules after adding d3-force + OriginTrace. d3-force is heavier than expected (~50KB). Still within acceptable range.
- The `motion.svg whileInView` + `motion.g` stagger combination works cleanly in Framer Motion 12 for SVG — no special SVG handling needed.

## Lessons for Future Missions

- **For d3 integrations**: Always check `package.json` types field + `node_modules/{pkg}/src/*.d.ts` to confirm whether @types is needed. Add to research checklist.
- **For force-directed graphs**: `useMemo` with synchronous ticks is cleaner than `useEffect + setState` for static layout. Use `useEffect + RAF` only if you want live simulation animation.
- **Historical source text overflow**: In force graphs, 10-30% label overlap is acceptable and expected. Don't add forceCollide radius large enough to push nodes to viewport edges.
- **Chinese text in SVG circles**: Safe rule: fontSize ≤ (diameter / num_chars) × 0.85. For r=34 circle: (68/7) × 0.85 = 8.3 → fontSize=8 is right on the edge. Consider r=36 for future with 7-char labels.
