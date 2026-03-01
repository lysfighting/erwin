# ADR-0004: d3-force for OriginTrace Visualization

- **Date**: 2026-03-02
- **Status**: Accepted
- **Session**: docs/sessions/2026-03-02.md

## Context

The OriginTrace component needs to render a force-directed graph showing the philosophical
lineage of drive-skill's design decisions (~30-50 nodes, ~60-80 edges).

Options for force-directed layout:
1. d3-force (just the physics module, 15KB)
2. Full D3.js (700KB)
3. Cytoscape.js (200KB)
4. Hand-crafted SVG positions (static)
5. React Flow (150KB)

## Decision

Use `d3-force` (only the physics submodule, not full D3) for layout calculation.
Render with React + existing Framer Motion for animations.
SVG elements managed by React, not D3's DOM manipulation.

Pattern: "D3 for math, React for DOM"

```typescript
import { forceSimulation, forceLink, forceManyBody, forceCenter } from 'd3-force';
// Use d3-force only for position calculation
// Pass positions to React state → Framer Motion renders
```

## Consequences

**Positive**:
- d3-force is ~15KB (vs 700KB full D3)
- No DOM manipulation conflict with React
- Framer Motion handles animation (already in stack)
- Force-directed layout adapts to node count dynamically

**Negative**:
- Slightly more code to bridge d3-force positions → React state
- Must add one new dependency

## Alternatives Considered

- **Full D3.js**: Rejected. 700KB for features we don't need. DOM manipulation conflicts with React.
- **Cytoscape.js**: Rejected. 200KB, complex API, not needed at this scale.
- **React Flow**: Rejected. 150KB, designed for flowcharts not philosophy graphs. Opinionated UI.
- **Static SVG positions**: Rejected. 50+ nodes need dynamic layout; manual positioning is unmaintainable.
