import { useState, useEffect, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'

interface StateMachineProps {
  autoPlay?: boolean
  intervalMs?: number
}

interface StateNode {
  id: string
  label: string
  color: string
  x: number
  y: number
}

interface Transition {
  from: string
  to: string
  label: string
  backward?: boolean
}

const STATES: StateNode[] = [
  { id: 'idea',        label: 'Idea',        color: '#8b949e', x: 60,  y: 80 },
  { id: 'queued',      label: 'Queued',      color: '#3b82f6', x: 175, y: 80 },
  { id: 'claimed',     label: 'Claimed',     color: '#a855f7', x: 290, y: 80 },
  { id: 'in_progress', label: 'In Progress', color: '#f59e0b', x: 405, y: 80 },
  { id: 'in_review',   label: 'In Review',   color: '#06b6d4', x: 520, y: 80 },
  { id: 'done',        label: 'Done',        color: '#22c55e', x: 635, y: 80 },
  { id: 'released',    label: 'Released',    color: '#4ade80', x: 750, y: 80 },
]

const TRANSITIONS: Transition[] = [
  { from: 'idea',        to: 'queued',      label: '/team-issue' },
  { from: 'queued',      to: 'claimed',     label: '/team-claim' },
  { from: 'claimed',     to: 'in_progress', label: '/team-drive' },
  { from: 'in_progress', to: 'in_review',   label: '/team-ship' },
  { from: 'in_review',   to: 'done',        label: 'PR merged' },
  { from: 'done',        to: 'released',    label: '/team-rc promote' },
  { from: 'in_review',   to: 'in_progress', label: 'Changes requested', backward: true },
]

const NODE_W = 100
const NODE_H = 36
const SVG_W = 820
const SVG_H = 180

function getNode(id: string) {
  return STATES.find(s => s.id === id)!
}

export default function StateMachine({ autoPlay = true, intervalMs = 2000 }: StateMachineProps) {
  const [activeIndex, setActiveIndex] = useState(0)
  const [lastTransition, setLastTransition] = useState<Transition | null>(null)

  const activeState = STATES[activeIndex]

  const advance = useCallback(() => {
    setActiveIndex(prev => {
      const next = prev < STATES.length - 1 ? prev + 1 : 0
      // Find the forward transition for this step
      const t = TRANSITIONS.find(
        t => t.from === STATES[prev].id && t.to === STATES[next].id && !t.backward
      )
      setLastTransition(t ?? null)
      return next
    })
  }, [])

  useEffect(() => {
    if (!autoPlay) return
    const id = setInterval(advance, intervalMs)
    return () => clearInterval(id)
  }, [autoPlay, intervalMs, advance])

  const handleNodeClick = (index: number) => {
    // Find transition from current to clicked
    const t = TRANSITIONS.find(
      t => t.from === STATES[activeIndex].id && t.to === STATES[index].id
    )
    setLastTransition(t ?? null)
    setActiveIndex(index)
  }

  return (
    <div className="w-full">
      <div className="overflow-x-auto">
        <svg
          viewBox={`0 0 ${SVG_W} ${SVG_H}`}
          className="w-full min-w-[600px]"
          style={{ maxHeight: 220 }}
        >
          {/* Forward edges */}
          {TRANSITIONS.filter(t => !t.backward).map(t => {
            const from = getNode(t.from)
            const to = getNode(t.to)
            const x1 = from.x + NODE_W / 2
            const x2 = to.x + NODE_W / 2
            const y = from.y
            const isActive =
              activeState.id === t.to &&
              lastTransition?.from === t.from &&
              lastTransition?.to === t.to
            return (
              <g key={`${t.from}-${t.to}`}>
                <line
                  x1={x1 + NODE_W / 2 + 2}
                  y1={y}
                  x2={x2 - NODE_W / 2 - 8}
                  y2={y}
                  stroke={isActive ? getNode(t.to).color : '#30363d'}
                  strokeWidth={isActive ? 2 : 1.5}
                  markerEnd={`url(#arrow-${isActive ? 'active' : 'default'})`}
                />
                <text
                  x={(x1 + x2) / 2}
                  y={y - 26}
                  textAnchor="middle"
                  fill={isActive ? '#e6edf3' : '#8b949e'}
                  fontSize={9}
                  fontFamily="monospace"
                >
                  {t.label}
                </text>
              </g>
            )
          })}

          {/* Backward edge: in_review → in_progress */}
          {TRANSITIONS.filter(t => t.backward).map(t => {
            const from = getNode(t.from)
            const to = getNode(t.to)
            const x1 = from.x + NODE_W / 2
            const x2 = to.x + NODE_W / 2
            const curveY = 150
            const isActive =
              activeState.id === t.to &&
              lastTransition?.from === t.from &&
              lastTransition?.to === t.to
            return (
              <g key={`${t.from}-${t.to}-back`}>
                <path
                  d={`M ${x1} ${from.y + NODE_H / 2} C ${x1} ${curveY}, ${x2} ${curveY}, ${x2} ${to.y + NODE_H / 2}`}
                  fill="none"
                  stroke={isActive ? '#ef4444' : '#30363d'}
                  strokeWidth={isActive ? 2 : 1.5}
                  strokeDasharray="4 3"
                  markerEnd={`url(#arrow-${isActive ? 'back' : 'default'})`}
                />
                <text
                  x={(x1 + x2) / 2}
                  y={curveY + 4}
                  textAnchor="middle"
                  fill={isActive ? '#ef4444' : '#484f58'}
                  fontSize={9}
                  fontFamily="monospace"
                >
                  {t.label}
                </text>
              </g>
            )
          })}

          {/* Arrow markers */}
          <defs>
            <marker id="arrow-default" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto">
              <path d="M0,0 L8,3 L0,6" fill="none" stroke="#30363d" strokeWidth="1" />
            </marker>
            <marker id="arrow-active" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto">
              <path d="M0,0 L8,3 L0,6" fill="none" stroke="#e6edf3" strokeWidth="1" />
            </marker>
            <marker id="arrow-back" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto">
              <path d="M0,0 L8,3 L0,6" fill="none" stroke="#ef4444" strokeWidth="1" />
            </marker>
          </defs>

          {/* State nodes */}
          {STATES.map((s, i) => {
            const isActive = i === activeIndex
            return (
              <g
                key={s.id}
                onClick={() => handleNodeClick(i)}
                className="cursor-pointer"
              >
                {/* Glow filter */}
                {isActive && (
                  <rect
                    x={s.x - 3}
                    y={s.y - NODE_H / 2 - 3}
                    width={NODE_W + 6}
                    height={NODE_H + 6}
                    rx={10}
                    fill="none"
                    stroke={s.color}
                    strokeWidth={2}
                    opacity={0.4}
                  >
                    <animate
                      attributeName="opacity"
                      values="0.4;0.8;0.4"
                      dur="1.5s"
                      repeatCount="indefinite"
                    />
                  </rect>
                )}
                <rect
                  x={s.x}
                  y={s.y - NODE_H / 2}
                  width={NODE_W}
                  height={NODE_H}
                  rx={8}
                  fill={isActive ? s.color + '22' : '#161b22'}
                  stroke={isActive ? s.color : '#30363d'}
                  strokeWidth={isActive ? 2 : 1}
                />
                <text
                  x={s.x + NODE_W / 2}
                  y={s.y + 1}
                  textAnchor="middle"
                  dominantBaseline="middle"
                  fill={isActive ? s.color : '#8b949e'}
                  fontSize={11}
                  fontWeight={isActive ? 'bold' : 'normal'}
                  fontFamily="system-ui, sans-serif"
                >
                  {s.label}
                </text>
              </g>
            )
          })}
        </svg>
      </div>

      {/* Transition command display */}
      <div className="h-8 flex items-center justify-center mt-2">
        <AnimatePresence mode="wait">
          {lastTransition && (
            <motion.div
              key={lastTransition.from + lastTransition.to}
              initial={{ opacity: 0, y: 6 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -6 }}
              transition={{ duration: 0.2 }}
              className="text-sm font-mono"
            >
              <span className="text-text-secondary">Trigger: </span>
              <span
                className="font-semibold"
                style={{ color: getNode(lastTransition.to).color }}
              >
                {lastTransition.label}
              </span>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  )
}
