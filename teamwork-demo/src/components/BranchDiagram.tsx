import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'

/**
 * BranchDiagram — interactive SVG visualization of the Teamwork four-branch model.
 *
 * Branches: main (🟥), pre-launch (🟦), mission/* (🟩), rc/* (🟧)
 * Shows fork relationships: mission forks FROM pre-launch, rc forks FROM pre-launch,
 * mission merges back TO pre-launch, rc merges TO main.
 */

interface BranchInfo {
  name: string
  color: string
  y: number
  label: string
  icon: string
  description: { en: string; zh: string }
}

const BRANCHES: BranchInfo[] = [
  {
    name: 'main',
    color: '#ef4444', // 🟥 red
    y: 50,
    label: 'main',
    icon: '🟥',
    description: {
      en: 'Production branch. Only receives squash merges from RC branches. Every commit is a tagged release (V1.0, V1.1...).',
      zh: '生产分支。只接收 RC 分支的 squash 合并。每个提交都是带标签的发布版本（V1.0, V1.1...）。',
    },
  },
  {
    name: 'pre-launch',
    color: '#3b82f6', // 🟦 blue
    y: 130,
    label: 'pre-launch',
    icon: '🟦',
    description: {
      en: 'Integration branch. All mission branches merge here via PR. Base for RC cuts. The "bleeding edge" of development.',
      zh: '集成分支。所有 mission 分支通过 PR 合并到此。RC 分支的起点。开发的"最前沿"。',
    },
  },
  {
    name: 'mission/042',
    color: '#22c55e', // 🟩 green
    y: 210,
    label: 'mission/042-add-auth-alice',
    icon: '🟩',
    description: {
      en: 'Mission branch created by /team-claim. Pattern: mission/{issue}-{slug}-{user}. Forks from pre-launch, merges back via PR.',
      zh: '由 /team-claim 创建的任务分支。命名：mission/{编号}-{简称}-{用户}。从 pre-launch 分叉，通过 PR 合并回去。',
    },
  },
  {
    name: 'rc/V1.0',
    color: '#f59e0b', // 🟧 orange
    y: 290,
    label: 'rc/V1.0.0',
    icon: '🟧',
    description: {
      en: 'Release candidate. Cut from pre-launch by /team-rc. Frozen for QA, then promoted to main via squash merge.',
      zh: '发布候选。由 /team-rc 从 pre-launch 切出。冻结用于 QA，然后通过 squash 合并到 main。',
    },
  },
]

// Connection arrows with fork points and merge points
interface Connection {
  type: 'fork' | 'merge'
  fromBranch: number
  fromX: number
  toBranch: number
  toX: number
  color: string
  label?: { en: string; zh: string }
}

const CONNECTIONS: Connection[] = [
  // mission/042 forks from pre-launch
  { type: 'fork', fromBranch: 1, fromX: 260, toBranch: 2, toX: 300, color: '#22c55e',
    label: { en: '/team-claim', zh: '/team-claim' } },
  // mission/042 merges back to pre-launch
  { type: 'merge', fromBranch: 2, fromX: 440, toBranch: 1, toX: 460, color: '#22c55e',
    label: { en: '/team-ship (PR)', zh: '/team-ship (PR)' } },
  // rc/V1.0 forks from pre-launch
  { type: 'fork', fromBranch: 1, fromX: 500, toBranch: 3, toX: 520, color: '#f59e0b',
    label: { en: '/team-rc', zh: '/team-rc' } },
  // rc/V1.0 merges to main
  { type: 'merge', fromBranch: 3, fromX: 640, toBranch: 0, toX: 660, color: '#f59e0b',
    label: { en: '/team-rc promote', zh: '/team-rc promote' } },
]

// Commit dots per branch
const COMMITS: { branch: number; x: number; tag?: string }[] = [
  // main
  { branch: 0, x: 100 }, { branch: 0, x: 200, tag: 'V0.9' }, { branch: 0, x: 660, tag: 'V1.0' },
  // pre-launch
  { branch: 1, x: 100 }, { branch: 1, x: 200 }, { branch: 1, x: 260 },
  { branch: 1, x: 380 }, { branch: 1, x: 460 }, { branch: 1, x: 500 },
  // mission/042
  { branch: 2, x: 300 }, { branch: 2, x: 360 }, { branch: 2, x: 400 }, { branch: 2, x: 440 },
  // rc/V1.0
  { branch: 3, x: 520 }, { branch: 3, x: 580 }, { branch: 3, x: 640 },
]

interface BranchDiagramProps {
  lang: 'en' | 'zh'
}

export default function BranchDiagram({ lang }: BranchDiagramProps) {
  const [activeBranch, setActiveBranch] = useState<number | null>(null)

  const isVisible = (bi: number) => activeBranch === null || activeBranch === bi
  const opacity = (bi: number) => isVisible(bi) ? 1 : 0.2

  return (
    <div className="w-full">
      <svg viewBox="0 0 740 340" className="w-full h-auto">
        {/* Background */}
        <rect width="740" height="340" fill="#0d1117" rx="8" />

        {/* Branch lines */}
        {BRANCHES.map((b, bi) => {
          const branchCommits = COMMITS.filter(c => c.branch === bi)
          if (branchCommits.length === 0) return null
          const minX = Math.min(...branchCommits.map(c => c.x))
          const maxX = Math.max(...branchCommits.map(c => c.x))

          return (
            <g key={b.name} opacity={opacity(bi)} className="cursor-pointer" onClick={() => setActiveBranch(activeBranch === bi ? null : bi)}>
              {/* Branch line */}
              <line
                x1={minX} y1={b.y} x2={maxX} y2={b.y}
                stroke={b.color}
                strokeWidth={activeBranch === bi ? 3 : 2}
                strokeOpacity={0.8}
              />
              {/* Icon + label */}
              <text x={12} y={b.y + 4} fill={b.color} fontSize="10" fontFamily="monospace" fontWeight="bold">
                {b.icon}
              </text>
              <text x={32} y={b.y + 4} fill={b.color} fontSize="11" fontFamily="monospace"
                fontWeight={activeBranch === bi ? 'bold' : 'normal'}>
                {b.label}
              </text>
            </g>
          )
        })}

        {/* Commits */}
        {COMMITS.map((c, ci) => (
          <g key={ci} opacity={opacity(c.branch)} className="cursor-pointer"
            onClick={() => setActiveBranch(activeBranch === c.branch ? null : c.branch)}>
            <circle cx={c.x} cy={BRANCHES[c.branch].y} r={activeBranch === c.branch ? 6 : 4.5}
              fill={BRANCHES[c.branch].color} />
            {c.tag && (
              <text x={c.x} y={BRANCHES[c.branch].y - 12}
                fill={BRANCHES[c.branch].color} fontSize="9" fontFamily="monospace"
                textAnchor="middle" fontWeight="bold">
                {c.tag}
              </text>
            )}
          </g>
        ))}

        {/* Connections (fork/merge arrows) */}
        {CONNECTIONS.map((conn, ci) => {
          const fromY = BRANCHES[conn.fromBranch].y
          const toY = BRANCHES[conn.toBranch].y
          const isActive = activeBranch === null ||
            activeBranch === conn.fromBranch || activeBranch === conn.toBranch

          // Curved path for fork/merge
          const midX = (conn.fromX + conn.toX) / 2
          const path = conn.type === 'fork'
            ? `M ${conn.fromX} ${fromY} C ${midX} ${fromY} ${midX} ${toY} ${conn.toX} ${toY}`
            : `M ${conn.fromX} ${fromY} C ${midX} ${fromY} ${midX} ${toY} ${conn.toX} ${toY}`

          return (
            <g key={ci} opacity={isActive ? 0.9 : 0.12}>
              {/* Path */}
              <path d={path} stroke={conn.color} strokeWidth={2}
                fill="none" strokeDasharray={conn.type === 'fork' ? 'none' : '6 3'} />
              {/* Arrow head at target */}
              <circle cx={conn.toX} cy={toY} r={3} fill={conn.color} />
              {/* Fork dot at source */}
              {conn.type === 'fork' && (
                <circle cx={conn.fromX} cy={fromY} r={3} fill={conn.color} strokeWidth={1} stroke="#0d1117" />
              )}
              {/* Label */}
              {conn.label && (
                <text
                  x={midX}
                  y={(fromY + toY) / 2 + (conn.type === 'fork' ? -6 : 6)}
                  fill={conn.color}
                  fontSize="8"
                  fontFamily="monospace"
                  textAnchor="middle"
                  opacity={0.8}
                >
                  {conn.label[lang]}
                </text>
              )}
            </g>
          )
        })}
      </svg>

      {/* Description panel */}
      <AnimatePresence mode="wait">
        {activeBranch !== null && (
          <motion.div
            key={activeBranch}
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -8 }}
            className="mt-4 p-4 rounded-lg border border-border bg-bg-card"
          >
            <div className="flex items-center gap-2 mb-2">
              <span>{BRANCHES[activeBranch].icon}</span>
              <span className="font-mono text-sm font-bold" style={{ color: BRANCHES[activeBranch].color }}>
                {BRANCHES[activeBranch].label}
              </span>
            </div>
            <p className="text-text-secondary text-sm">
              {BRANCHES[activeBranch].description[lang]}
            </p>
          </motion.div>
        )}
      </AnimatePresence>

      {activeBranch === null && (
        <p className="mt-3 text-text-muted text-sm text-center">
          {lang === 'en' ? 'Click a branch to learn more' : '点击分支了解详情'}
        </p>
      )}
    </div>
  )
}
