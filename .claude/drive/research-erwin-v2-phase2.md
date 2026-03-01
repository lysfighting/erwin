# Research: erwin V2 Phase 2 — OriginTrace

## Stack confirmed
- React 19.2 + TypeScript 5.9 + Vite 7.3 + Framer Motion 12.34 + Tailwind 4.2
- d3-force (to be installed) — only module needed, not full D3

## Component patterns observed
- PrincipleCards: motion.div whileInView, AnimatePresence expand/collapse, typed interfaces
- ExecutionSimulator: useState for runKey remount, Framer Motion stagger
- All components: named exports, Tailwind dark-theme classes (bg-bg-card, text-text-secondary, border-border, text-gold)

## Graph data (architecture doc §OriginTrace)
- 5 drive-concept nodes (matching PrincipleCards groups A-E)
- 2 historical source nodes per concept = 10 sources
- Total: 15 nodes, 10 edges
- Colors match PrincipleCards: A=gold, B=blue, C=green, D=purple, E=cyan

## d3-force integration approach
- SimulationNodeDatum + SimulationLinkDatum type extensions
- Run synchronously: .stop().tick(400)
- Nodes mutated in place — capture positions after tick
- SimLink.source/target become node refs after forceLink resolves

## Layout strategy
- forceLink distance=140, strength=0.8
- forceManyBody strength=-600 (repulsion)
- forceCenter at (width/2, height/2)
- forceCollide radius: 62 for drive-concept, 52 for historical-source
- Clamp all nodes to viewport with 90px padding after ticks

## Rendering approach
- SVG: motion.svg whileInView opacity fade-in
- Lines: plain SVG <line> elements, CSS transition for highlight
- Nodes: motion.g with stagger opacity entrance animation
- Drive-concept nodes: r=34, solid filled, Chinese text inside
- Historical-source nodes: r=22, transparent fill, text below
- Detail panel: HTML below SVG, AnimatePresence slide-in

## Text sizing validation
- '永不停止的循环' = 7 chars × 8px fontSize = 56px < 68px diameter ✓
- SVG fontFamily: "ui-monospace, SFMono-Regular, Menlo, monospace"

## App.tsx integration
- Uncomment <OriginTrace /> placeholder in ACT II
- Add SectionDivider + max-w-6xl wrapper
- Import from './components/OriginTrace'
