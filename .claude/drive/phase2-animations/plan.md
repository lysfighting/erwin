# Plan: Phase 2 Animation Enhancements
## Contract: .claude/drive/phase2-animations/contract.md

## Approach Overview
Add message flow animations to two existing components using pure SVG + Framer Motion.
Both components share a "flowing dot" visual metaphor — small circles animated along
connection paths between agent nodes, with consistent color coding.

## Architecture / Design

### Shared Constants (inline in each component)
```
dispatch: #22c55e (green-500)
return:   #3b82f6 (blue-500)
broadcast:#f97316 (orange-500)
dot radius: 4px
dot travel duration: 1.0s
stagger between messages: 0.3s
```

### ExecutionSimulator Enhancement
```
┌──────────────────────────────────────────────────┐
│ With drive-skill                                  │
│                                                    │
│     ┌──────────┐                                  │
│     │Team Lead │──┬──→ STL ──┬──→ Worker-1       │
│     └──────────┘  │         └──→ Worker-2        │
│                    └──→ Tester                     │
│   (animated dots flow along connections)           │
│                                                    │
│   ████████████████████████ Task 1     ─────┐      │
│   ████████████████████████ Task 2      1.6s│      │
│   ████████████████████████ Task 3          │      │
│   ████████████████████████ Task 4     ─────┘      │
│                                                    │
│   Total: 4s  ⚡ 3x faster                         │
└──────────────────────────────────────────────────┘
```
- SVG ~280x100 placed above parallel bars
- 4 nodes: Team Lead, STL, Worker-1, Worker-2
- Connection lines (thin, gray)
- When playing=true: green dots flow TL→STL→Workers, then blue dots return
- Animation synced with existing bar animation timing

### SwarmTopologySection Enhancement
- Reuse existing 9-node SVG (800x380)
- Add MessageDot sub-component: motion.circle with animated cx/cy
- Auto-play sequence triggered by isInView:
  - Phase 1 (0s): Green dots TL→STLs + TL→Tester (3 paths, 0.3s stagger)
  - Phase 2 (1.5s): Green dots STLs→inner subagents (5 paths)
  - Phase 3 (3.5s): Blue dots returning inner→STLs→TL
- Click interaction:
  - activeNodeId state
  - Click node → only animate that node's connections
  - Click again or click outside → reset to auto-play
- Highlight: non-active connections fade to 0.1 opacity during focused view

## Implementation Sequence
1. Read both component files completely
2. Implement SwarmTopology message flow (larger, establishes the pattern)
3. Implement ExecutionSimulator topology + flow (adapts the pattern)
4. Dev server test + TypeScript check
5. /visual-test screenshot verification

## Risk Mitigation
- SVG performance: <20 dots at any time, well within SVG limits
- If Framer Motion cx/cy animation is janky: fall back to CSS keyframes
- Click interaction: simple state toggle, no complex state machine needed
