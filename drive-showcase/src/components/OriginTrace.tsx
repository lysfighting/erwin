import { useState, useMemo, useEffect, useRef } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import {
  forceSimulation,
  forceLink,
  forceManyBody,
  forceCenter,
  forceCollide,
} from 'd3-force'
import type { SimulationNodeDatum, SimulationLinkDatum } from 'd3-force'

// ── Types ─────────────────────────────────────────────────────────────────────

interface NodeDef {
  id: string
  label: string
  sublabel: string
  type: 'drive-concept' | 'historical-source'
  group: 'A' | 'B' | 'C' | 'D' | 'E'
  description: string
  driveSkillManifestation: string
}

interface SimNode extends SimulationNodeDatum {
  id: string
  label: string
  sublabel: string
  type: 'drive-concept' | 'historical-source'
  group: 'A' | 'B' | 'C' | 'D' | 'E'
  description: string
  driveSkillManifestation: string
}

type SimLink = SimulationLinkDatum<SimNode>

// ── Color palette (matches PrincipleCards accent colors) ──────────────────────

const COLOR: Record<string, string> = {
  A: '#f59e0b', // gold   — 永不停止的循环
  B: '#3b82f6', // blue   — 层级指挥结构
  C: '#22c55e', // green  — 最大并行度
  D: '#a855f7', // purple — 透明广播协议
  E: '#06b6d4', // cyan   — 阶段纪律
}

const MONO = 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace'

// ── Node data ─────────────────────────────────────────────────────────────────

const NODE_DATA: NodeDef[] = [
  // ── Drive Concepts (center nodes) ──────────────────────────────────────────
  {
    id: 'cardinal-rule',
    label: '永不停止的循环',
    sublabel: 'The Loop Never Breaks',
    type: 'drive-concept',
    group: 'A',
    description:
      'The system never halts. Every error is a recovery opportunity, every completion triggers the next cycle. Autonomous momentum is the cardinal rule.',
    driveSkillManifestation:
      'The Loop runs until MISSION_COMPLETE. Errors trigger recovery, not termination. Cardinal Rule: outputting "should I continue?" is failure — always continue.',
  },
  {
    id: 'stl-hierarchy',
    label: '层级指挥结构',
    sublabel: 'STL Hierarchy',
    type: 'drive-concept',
    group: 'B',
    description:
      'A single STL (Sub-Team Lead) coordinates up to 4 agents. Intent flows down, results flow up. Execution is autonomous at every level.',
    driveSkillManifestation:
      'STLs own domains and spawn inner subagents. Team Lead orchestrates STLs, not individual workers. "Own the frontend" — not a task list.',
  },
  {
    id: 'max-parallelism',
    label: '最大并行度',
    sublabel: 'Maximum Parallelism',
    type: 'drive-concept',
    group: 'C',
    description:
      'If tasks have no dependency, they run simultaneously. The critical path is minimized. Serial execution is treated as a bug.',
    driveSkillManifestation:
      'Wave decomposition: all tasks in a wave run in parallel. Wall-clock time approaches the longest single task, not the sum of all tasks.',
  },
  {
    id: 'broadcast',
    label: '透明广播协议',
    sublabel: 'Broadcast Protocol',
    type: 'drive-concept',
    group: 'D',
    description:
      'Every dispatch and return is broadcast to all agents. No information silos. Visibility eliminates uncertainty and prevents duplicated work.',
    driveSkillManifestation:
      '📡 DISPATCH and ✅ RETURN broadcasts are mandatory. The swarm is not a black box — it is a live war room. Every agent sees every move.',
  },
  {
    id: 'phase-discipline',
    label: '阶段纪律',
    sublabel: 'Phase Discipline',
    type: 'drive-concept',
    group: 'E',
    description:
      'Thinking before doing. Phase 0 (research) and Phase T (team formation) run before any code is written. Preparation prevents rework.',
    driveSkillManifestation:
      'Mandatory research-{slug}.md → plan-{slug}.md → annotation cycles before tasks are created. Code is the execution of thought — thought must come first.',
  },

  // ── Group A: The Loop sources ───────────────────────────────────────────────
  {
    id: 'ooda-loop',
    label: 'OODA Loop',
    sublabel: 'Boyd, 1970s',
    type: 'historical-source',
    group: 'A',
    description:
      "Observe → Orient → Decide → Act. Colonel John Boyd's decision cycle: speed of re-orientation is decisive advantage. The loop that iterates fastest wins.",
    driveSkillManifestation:
      'The drive Loop mirrors OODA: observe errors → orient to root cause → decide on fix → act. Never pause to ask permission mid-execution.',
  },
  {
    id: 'cybernetics',
    label: 'Cybernetics',
    sublabel: 'Wiener, 1948',
    type: 'historical-source',
    group: 'A',
    description:
      "Norbert Wiener's theory of self-regulating systems. Feedback loops allow systems to maintain goals despite disturbances. Error signals are the mechanism of control.",
    driveSkillManifestation:
      'Error recovery as feedback: each error narrows the problem space. The system self-corrects without human intervention — disturbances are expected, not catastrophic.',
  },

  // ── Group B: STL Hierarchy sources ─────────────────────────────────────────
  {
    id: 'mission-command',
    label: 'Mission Command',
    sublabel: 'Auftragstaktik, Prussia',
    type: 'historical-source',
    group: 'B',
    description:
      "Prussian military doctrine: leaders communicate intent, not step-by-step instructions. Subordinates are trusted to execute autonomously within the commander's intent.",
    driveSkillManifestation:
      'STLs receive domain intent ("own the frontend"), not task lists. Tactical decisions are delegated. Strategic decisions stay at Team Lead.',
  },
  {
    id: 'mythical-man-month',
    label: 'Mythical Man-Month',
    sublabel: 'Brooks, 1975',
    type: 'historical-source',
    group: 'B',
    description:
      'Frederick Brooks on software management. The "surgical team" model: one chief surgeon supported by specialists. Not a committee — a hierarchy with one clear decision-maker.',
    driveSkillManifestation:
      'Team Lead as surgeon: one clear decision-maker, teammates as force multipliers. Adding more agents naively increases coordination cost — hierarchy solves this.',
  },

  // ── Group C: Parallelism sources ───────────────────────────────────────────
  {
    id: 'cpm-pert',
    label: 'CPM / PERT',
    sublabel: '1957',
    type: 'historical-source',
    group: 'C',
    description:
      'Critical Path Method (DuPont, 1957) and PERT (US Navy, 1957). Find the longest dependency chain. Minimizing it minimizes total project duration.',
    driveSkillManifestation:
      'Wave Map critical path analysis. Team Lead always takes critical-path tasks. Off-critical-path tasks have slack — they run whenever, not blocking the clock.',
  },
  {
    id: 'amdahl-law',
    label: "Amdahl's Law",
    sublabel: 'Gene Amdahl, 1967',
    type: 'historical-source',
    group: 'C',
    description:
      "Speedup from parallelism is limited by the sequential fraction. To achieve 10× speedup, the serial fraction must be <10%. Maximum parallelism = minimize serialization.",
    driveSkillManifestation:
      "Serial execution is treated as a bug. The Wave Map eliminates unnecessary serial dependencies. Every task that CAN be parallel MUST be parallel — Amdahl's limit demands it.",
  },

  // ── Group D: Broadcast sources ─────────────────────────────────────────────
  {
    id: 'information-theory',
    label: 'Information Theory',
    sublabel: 'Shannon, 1948',
    type: 'historical-source',
    group: 'D',
    description:
      "Claude Shannon's mathematical theory of communication. Information reduces entropy (uncertainty). The value of a message = the uncertainty it eliminates.",
    driveSkillManifestation:
      'Broadcast reduces swarm entropy: all agents always know what all other agents are doing. Zero information silos = zero coordination waste = maximum throughput.',
  },
  {
    id: 'agile-manifesto',
    label: 'Agile Manifesto',
    sublabel: '2001',
    type: 'historical-source',
    group: 'D',
    description:
      '17 software practitioners in Snowbird, Utah. "Individuals and interactions over processes." Daily stand-ups as continuous synchronization. Visibility as team health.',
    driveSkillManifestation:
      'Every dispatch and return is an automatic stand-up. No scheduled syncs needed — the broadcast protocol IS continuous sync. Team visibility is a first-class requirement.',
  },

  // ── Group E: Phase Discipline sources ──────────────────────────────────────
  {
    id: 'sun-tzu',
    label: '孙子兵法',
    sublabel: '孙子, ~500 BC',
    type: 'historical-source',
    group: 'E',
    description:
      '"知己知彼，百战不殆。" Know yourself and your enemy, and you need not fear a hundred battles. Intelligence and preparation precede action — always.',
    driveSkillManifestation:
      'Phase 0 Research Swarm: scout the codebase before writing a line. 谋定而后动 — plan before you move. Reconnaissance is not optional.',
  },
  {
    id: 'double-diamond',
    label: 'Double Diamond',
    sublabel: 'Design Council, 2004',
    type: 'historical-source',
    group: 'E',
    description:
      'The Design Council model: Discover → Define → Develop → Deliver. Diverge to explore, converge to commit — twice. Each phase has a distinct mode of thinking.',
    driveSkillManifestation:
      'Phase 0 diverges (research swarm explores all angles), Phase T converges (plan + team), Loop diverges (parallel implementation), Final Verify converges.',
  },
]

// ── Edge data ──────────────────────────────────────────────────────────────────

const EDGE_DATA = [
  { source: 'ooda-loop', target: 'cardinal-rule' },
  { source: 'cybernetics', target: 'cardinal-rule' },
  { source: 'mission-command', target: 'stl-hierarchy' },
  { source: 'mythical-man-month', target: 'stl-hierarchy' },
  { source: 'cpm-pert', target: 'max-parallelism' },
  { source: 'amdahl-law', target: 'max-parallelism' },
  { source: 'information-theory', target: 'broadcast' },
  { source: 'agile-manifesto', target: 'broadcast' },
  { source: 'sun-tzu', target: 'phase-discipline' },
  { source: 'double-diamond', target: 'phase-discipline' },
]

// ── Force layout (synchronous — runs once per size) ───────────────────────────

function useForceLayout(width: number, height: number) {
  return useMemo(() => {
    if (width === 0) return { nodes: [] as SimNode[], links: [] as SimLink[] }

    // Clone so d3 can mutate positions without touching source arrays
    const simNodes = NODE_DATA.map<SimNode>((n) => ({ ...n }))
    const simLinks: SimLink[] = EDGE_DATA.map((e) => ({
      source: e.source,
      target: e.target,
    }))

    // Pin the 5 drive-concept nodes at a deterministic pentagon.
    // This anchors clusters and prevents asymmetric drift caused by random
    // initial positions + strong repulsion pushing nodes to viewport edges.
    const driveNodes = simNodes.filter((n) => n.type === 'drive-concept')
    const pentagonR = Math.min(width, height) * 0.19
    driveNodes.forEach((node, i) => {
      const angle = (i * (2 * Math.PI)) / driveNodes.length - Math.PI / 2
      node.fx = width / 2 + pentagonR * Math.cos(angle)
      node.fy = height / 2 + pentagonR * Math.sin(angle)
    })

    forceSimulation<SimNode>(simNodes)
      .force(
        'link',
        forceLink<SimNode, SimLink>(simLinks)
          .id((d) => d.id)
          .distance(115)
          .strength(1.0),
      )
      .force('charge', forceManyBody<SimNode>().strength(-280))
      .force('center', forceCenter<SimNode>(width / 2, height / 2))
      .force(
        'collision',
        forceCollide<SimNode>().radius((d) => (d.type === 'drive-concept' ? 50 : 44)),
      )
      .stop()
      .tick(400)

    // Soft clamp for source nodes only — drive concepts are already pinned
    const pad = 60
    simNodes
      .filter((n) => n.type === 'historical-source')
      .forEach((n) => {
        n.x = Math.max(pad, Math.min(width - pad, n.x ?? width / 2))
        n.y = Math.max(pad, Math.min(height - pad, n.y ?? height / 2))
      })

    return { nodes: simNodes, links: simLinks }
  }, [width, height])
}

// ── Component ─────────────────────────────────────────────────────────────────

export function OriginTrace() {
  const containerRef = useRef<HTMLDivElement>(null)
  const [width, setWidth] = useState(0)
  const height = Math.round(width * 0.62)

  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [hoveredId, setHoveredId] = useState<string | null>(null)

  useEffect(() => {
    const el = containerRef.current
    if (!el) return
    const ro = new ResizeObserver(([entry]) => {
      setWidth(Math.floor(entry.contentRect.width))
    })
    ro.observe(el)
    return () => ro.disconnect()
  }, [])

  const { nodes, links } = useForceLayout(width, height)
  const selectedNode = nodes.find((n) => n.id === selectedId) ?? null

  // Highlight edges connected to the hovered or selected node
  const activeId = hoveredId ?? selectedId

  function toggleSelect(id: string) {
    setSelectedId((prev) => (prev === id ? null : id))
  }

  return (
    <div className="w-full">
      {/* Section header */}
      <div className="text-center mb-8">
        <h2 className="font-mono text-2xl md:text-3xl font-bold text-text-primary mb-2">
          思想谱系 · Origin Trace
        </h2>
        <p className="text-text-secondary text-sm">
          Where did drive-skill's ideas come from? Click any node to explore its origin.
        </p>
      </div>

      {/* Graph canvas */}
      <div
        ref={containerRef}
        className="w-full bg-bg-card border border-border rounded-xl overflow-hidden"
      >
        {width > 0 && nodes.length > 0 ? (
          <motion.svg
            width={width}
            height={height}
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true, margin: '-80px' }}
            transition={{ duration: 0.7 }}
          >
            {/* ── Edges ── */}
            <g>
              {links.map((link, i) => {
                const s = link.source as SimNode
                const t = link.target as SimNode
                const highlighted =
                  activeId === s.id || activeId === t.id
                const edgeColor = COLOR[t.group] ?? COLOR[s.group]
                return (
                  <line
                    key={i}
                    x1={s.x ?? 0}
                    y1={s.y ?? 0}
                    x2={t.x ?? 0}
                    y2={t.y ?? 0}
                    stroke={highlighted ? edgeColor : '#374151'}
                    strokeWidth={highlighted ? 2 : 1}
                    strokeOpacity={highlighted ? 0.75 : 0.3}
                    style={{
                      transition: 'stroke 0.2s, stroke-width 0.2s, stroke-opacity 0.2s',
                    }}
                  />
                )
              })}
            </g>

            {/* ── Nodes ── */}
            {nodes.map((node, i) => {
              const x = node.x ?? 0
              const y = node.y ?? 0
              const isDrive = node.type === 'drive-concept'
              const r = isDrive ? 34 : 22
              const color = COLOR[node.group]
              const isSelected = selectedId === node.id
              const isHovered = hoveredId === node.id
              const isActive = isSelected || isHovered

              return (
                <motion.g
                  key={node.id}
                  transform={`translate(${x},${y})`}
                  onClick={() => toggleSelect(node.id)}
                  onMouseEnter={() => setHoveredId(node.id)}
                  onMouseLeave={() => setHoveredId(null)}
                  style={{ cursor: 'pointer' }}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ delay: 0.35 + i * 0.045, duration: 0.35 }}
                >
                  {/* Active glow ring */}
                  {isActive && (
                    <circle
                      r={r + 8}
                      fill="none"
                      stroke={color}
                      strokeWidth={1.5}
                      strokeOpacity={0.4}
                    />
                  )}

                  {/* Main circle */}
                  <circle
                    r={r}
                    fill={isDrive ? color : 'rgba(17,24,39,0.9)'}
                    stroke={color}
                    strokeWidth={isActive ? 2.5 : isDrive ? 0 : 1.5}
                    strokeOpacity={isActive ? 1 : 0.75}
                    style={{ transition: 'stroke-width 0.2s, stroke-opacity 0.2s' }}
                  />

                  {isDrive ? (
                    // Drive concept: Chinese label centered inside the filled circle
                    <text
                      textAnchor="middle"
                      dominantBaseline="middle"
                      y={0}
                      fontSize={8}
                      fontWeight="bold"
                      fill="#0d1117"
                      fontFamily={MONO}
                      style={{ pointerEvents: 'none', userSelect: 'none' }}
                    >
                      {node.label}
                    </text>
                  ) : (
                    // Historical source: label + year below the circle
                    <>
                      <text
                        textAnchor="middle"
                        y={r + 14}
                        fontSize={9}
                        fontWeight={isActive ? 'bold' : 'normal'}
                        fill={isActive ? '#f9fafb' : color}
                        fillOpacity={isActive ? 1 : 0.9}
                        fontFamily={MONO}
                        style={{
                          pointerEvents: 'none',
                          userSelect: 'none',
                          transition: 'fill 0.2s, font-weight 0.2s',
                        }}
                      >
                        {node.label}
                      </text>
                      <text
                        textAnchor="middle"
                        y={r + 25}
                        fontSize={8}
                        fill={isActive ? '#9ca3af' : '#6b7280'}
                        fontFamily={MONO}
                        style={{
                          pointerEvents: 'none',
                          userSelect: 'none',
                          transition: 'fill 0.2s',
                        }}
                      >
                        {node.sublabel}
                      </text>
                    </>
                  )}
                </motion.g>
              )
            })}
          </motion.svg>
        ) : (
          <div className="h-64 flex items-center justify-center">
            <span className="text-text-muted font-mono text-sm animate-pulse">
              Mapping origins…
            </span>
          </div>
        )}
      </div>

      {/* Detail panel — appears below graph on node selection */}
      <AnimatePresence>
        {selectedNode && (
          <motion.div
            key={selectedNode.id}
            initial={{ opacity: 0, y: 14 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 8 }}
            transition={{ duration: 0.25 }}
            className="mt-4 bg-bg-card border border-border rounded-xl p-5 md:p-6"
          >
            <div className="flex items-start gap-3 mb-3">
              <div
                className="flex-shrink-0 w-2.5 h-2.5 rounded-full mt-1.5"
                style={{ backgroundColor: COLOR[selectedNode.group] }}
              />
              <div className="flex-1 min-w-0">
                <h3 className="font-mono text-sm font-bold text-text-primary leading-snug">
                  {selectedNode.label}
                </h3>
                <p
                  className="font-mono text-xs mt-0.5"
                  style={{ color: COLOR[selectedNode.group] }}
                >
                  {selectedNode.sublabel}
                </p>
              </div>
              <button
                onClick={() => setSelectedId(null)}
                className="flex-shrink-0 font-mono text-xs text-text-muted hover:text-text-secondary transition-colors px-1"
              >
                ✕
              </button>
            </div>

            <p className="text-text-secondary text-sm leading-relaxed mb-4">
              {selectedNode.description}
            </p>

            <div className="border-t border-border pt-3">
              <span className="font-mono text-xs text-text-muted uppercase tracking-wider">
                → drive-skill manifestation
              </span>
              <p className="mt-1.5 font-mono text-xs text-text-secondary leading-relaxed">
                {selectedNode.driveSkillManifestation}
              </p>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Legend */}
      <div className="mt-4 flex flex-wrap gap-x-5 gap-y-2 justify-center">
        {(['A', 'B', 'C', 'D', 'E'] as const).map((group) => {
          const driveNode = NODE_DATA.find(
            (n) => n.group === group && n.type === 'drive-concept',
          )
          return (
            <div key={group} className="flex items-center gap-1.5">
              <div
                className="w-2.5 h-2.5 rounded-full"
                style={{ backgroundColor: COLOR[group] }}
              />
              <span className="font-mono text-xs text-text-muted">
                {driveNode?.label ?? group}
              </span>
            </div>
          )
        })}
        <div className="flex items-center gap-1.5">
          <div className="w-2.5 h-2.5 rounded-full border border-gray-500" />
          <span className="font-mono text-xs text-text-muted">Historical Source</span>
        </div>
      </div>
    </div>
  )
}
