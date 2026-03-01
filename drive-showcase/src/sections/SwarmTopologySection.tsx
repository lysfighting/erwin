import { motion, useInView } from 'framer-motion'
import { useRef, useState, useEffect } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

interface NodeData {
  id: string
  label: string
  role: string
  color: string
  x: number
  y: number
  parentId?: string
  delay: number
  children?: string[]
}

const nodes: NodeData[] = [
  { id: 'lead', label: 'Team Lead', role: 'YOU', color: '#f59e0b', x: 400, y: 60, delay: 0 },
  { id: 'stl-fe', label: 'STL-frontend', role: '🎖️ Sub-Team-Lead', color: '#3b82f6', x: 200, y: 180, parentId: 'lead', delay: 0.5 },
  { id: 'stl-be', label: 'STL-backend', role: '🎖️ Sub-Team-Lead', color: '#3b82f6', x: 500, y: 180, parentId: 'lead', delay: 0.6 },
  { id: 'tester', label: 'tester', role: '🔧 Leaf', color: '#a855f7', x: 700, y: 180, parentId: 'lead', delay: 0.7 },
  // STL-FE inner subagents
  { id: 'fe-1', label: 'comp-A', role: 'inner', color: '#22c55e', x: 100, y: 300, parentId: 'stl-fe', delay: 1.2 },
  { id: 'fe-2', label: 'comp-B', role: 'inner', color: '#22c55e', x: 200, y: 300, parentId: 'stl-fe', delay: 1.3 },
  { id: 'fe-3', label: 'style', role: 'inner', color: '#22c55e', x: 300, y: 300, parentId: 'stl-fe', delay: 1.4 },
  // STL-BE inner subagents
  { id: 'be-1', label: 'API-X', role: 'inner', color: '#22c55e', x: 440, y: 300, parentId: 'stl-be', delay: 1.3 },
  { id: 'be-2', label: 'API-Y', role: 'inner', color: '#22c55e', x: 560, y: 300, parentId: 'stl-be', delay: 1.4 },
]

function SwarmNode({ node, isVisible }: { node: NodeData; isVisible: boolean }) {
  const size = node.role === 'inner' ? 28 : node.id === 'lead' ? 44 : 36

  return (
    <motion.g
      initial={{ opacity: 0, scale: 0 }}
      animate={isVisible ? { opacity: 1, scale: 1 } : {}}
      transition={{ delay: node.delay, duration: 0.4, type: 'spring', stiffness: 200 }}
    >
      {/* Pulse ring for active nodes */}
      {isVisible && (
        <motion.circle
          cx={node.x}
          cy={node.y}
          r={size / 2}
          fill="none"
          stroke={node.color}
          strokeWidth={1}
          initial={{ r: size / 2, opacity: 0.6 }}
          animate={{ r: size / 2 + 15, opacity: 0 }}
          transition={{ duration: 2, repeat: Infinity, delay: node.delay + 0.5 }}
        />
      )}
      {/* Node circle */}
      <circle
        cx={node.x}
        cy={node.y}
        r={size / 2}
        fill={node.color + '20'}
        stroke={node.color}
        strokeWidth={2}
      />
      {/* Label */}
      <text
        x={node.x}
        y={node.y + size / 2 + 16}
        textAnchor="middle"
        fill="#8b949e"
        fontSize={11}
        fontFamily="'JetBrains Mono', monospace"
      >
        {node.label}
      </text>
      {/* Role badge */}
      {node.role !== 'inner' && (
        <text
          x={node.x}
          y={node.y + 4}
          textAnchor="middle"
          fill={node.color}
          fontSize={10}
          fontFamily="'JetBrains Mono', monospace"
          fontWeight="bold"
        >
          {node.role === 'YOU' ? 'YOU' : node.role.includes('STL') ? 'STL' : node.role.slice(2)}
        </text>
      )}
    </motion.g>
  )
}

function ConnectionLine({ from, to, delay, isVisible }: { from: NodeData; to: NodeData; delay: number; isVisible: boolean }) {
  return (
    <motion.line
      x1={from.x}
      y1={from.y + 20}
      x2={to.x}
      y2={to.y - 20}
      stroke={to.color}
      strokeWidth={1.5}
      strokeDasharray="4 4"
      initial={{ pathLength: 0, opacity: 0 }}
      animate={isVisible ? { pathLength: 1, opacity: 0.4 } : {}}
      transition={{ delay, duration: 0.5 }}
    />
  )
}

const broadcasts = [
  { text: '📡 DISPATCH: [stl-frontend] → UI components', delay: 1.5, color: '#3b82f6' },
  { text: '📡 DISPATCH: [stl-backend] → API endpoints', delay: 1.8, color: '#3b82f6' },
  { text: '🎖️ STL-INNER: stl-frontend spawned 3 subagents', delay: 2.5, color: '#f59e0b' },
  { text: '🎖️ STL-INNER: stl-backend spawned 2 subagents', delay: 2.8, color: '#f59e0b' },
  { text: '✅ RETURN: [fe/comp-A] ← Component A complete', delay: 4, color: '#22c55e' },
  { text: '✅ RETURN: [be/API-X] ← Endpoint X complete', delay: 4.5, color: '#22c55e' },
]

export function SwarmTopologySection() {
  const { lang } = useLang()
  const ref = useRef(null)
  const isInView = useInView(ref, { once: true, margin: '-100px' })
  const [visibleBroadcasts, setVisibleBroadcasts] = useState<number[]>([])

  useEffect(() => {
    if (!isInView) return
    broadcasts.forEach((b, i) => {
      setTimeout(() => {
        setVisibleBroadcasts(prev => [...prev, i])
      }, b.delay * 1000)
    })
  }, [isInView])

  return (
    <section id="swarm" className="min-h-screen py-24 px-4">
      <div className="max-w-6xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 40 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
        >
          <h2 className="text-4xl font-bold mb-2">{t('swarm.title', lang)}</h2>
          <p className="text-text-secondary text-lg mb-4">
            {t('swarm.subtitle', lang)}
          </p>
          <p className="text-text-muted text-sm mb-12">
            {t('swarm.desc', lang)}
          </p>
        </motion.div>

        {/* Why Hierarchy? */}
        <div className="grid md:grid-cols-3 gap-4 mb-10">
          {[
            {
              icon: '⚡',
              title: 'The Coordination Problem',
              color: '#ef4444',
              desc: 'If every worker reports to the Team Lead, that\'s O(N) messages per wave. With 10 workers, the Team Lead becomes the bottleneck — spending all time coordinating, none executing. Flat teams break at scale.',
            },
            {
              icon: '🎖️',
              title: 'The STL Multiplier',
              color: '#3b82f6',
              desc: '1 STL + 3 inner subagents = 4x throughput at 1 teammate-mode slot. The Team Lead sees STLs (domain managers), not their 3 subtasks. Delegation compresses coordination overhead while multiplying throughput.',
            },
            {
              icon: '👁️',
              title: 'Visibility by Design',
              color: '#f59e0b',
              desc: 'STLs are Agent tool calls — they appear as visible panes in your terminal (Shift+Down to switch). Inner Task subagents are invisible but fast. The hierarchy maps to what you can see vs. what runs silently.',
            },
          ].map((item, i) => (
            <motion.div
              key={item.title}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.1 }}
              className="bg-bg-card border border-border rounded-lg p-4"
              style={{ borderLeftColor: item.color, borderLeftWidth: 3 }}
            >
              <div className="flex items-center gap-2 mb-2">
                <span className="text-xl">{item.icon}</span>
                <span className="font-mono text-sm font-bold text-text-primary">{item.title}</span>
              </div>
              <p className="text-sm text-text-secondary leading-relaxed">{item.desc}</p>
            </motion.div>
          ))}
        </div>

        <div ref={ref} className="grid lg:grid-cols-5 gap-8">
          {/* Topology visualization */}
          <div className="lg:col-span-3 bg-bg-card border border-border rounded-lg p-4 overflow-hidden">
            <svg viewBox="0 0 800 380" className="w-full">
              {/* Connection lines */}
              {nodes.filter(n => n.parentId).map(node => {
                const parent = nodes.find(n => n.id === node.parentId)!
                return (
                  <ConnectionLine
                    key={`line-${node.id}`}
                    from={parent}
                    to={node}
                    delay={node.delay - 0.2}
                    isVisible={isInView}
                  />
                )
              })}
              {/* Nodes */}
              {nodes.map(node => (
                <SwarmNode key={node.id} node={node} isVisible={isInView} />
              ))}
              {/* Legend */}
              <g>
                {[
                  { color: '#f59e0b', label: 'Team Lead', y: 350 },
                  { color: '#3b82f6', label: 'Sub-Team-Lead', y: 350 },
                  { color: '#22c55e', label: 'Inner Subagent', y: 350 },
                  { color: '#a855f7', label: 'Leaf Teammate', y: 350 },
                ].map((item, i) => (
                  <g key={item.label} transform={`translate(${80 + i * 180}, ${item.y})`}>
                    <circle r={5} fill={item.color} />
                    <text x={12} y={4} fill="#8b949e" fontSize={10} fontFamily="monospace">{item.label}</text>
                  </g>
                ))}
              </g>
            </svg>
          </div>

          {/* Broadcast feed */}
          <div className="lg:col-span-2 bg-bg-card border border-border rounded-lg p-4">
            <h3 className="font-mono text-sm text-gold mb-4">Live Broadcast Feed</h3>
            <div className="space-y-2 font-mono text-xs">
              {broadcasts.map((b, i) => (
                <motion.div
                  key={i}
                  initial={{ opacity: 0, x: 20 }}
                  animate={visibleBroadcasts.includes(i) ? { opacity: 1, x: 0 } : {}}
                  transition={{ duration: 0.3 }}
                  className="p-2 bg-bg-dark rounded border border-border"
                  style={{ borderLeftColor: b.color, borderLeftWidth: 3 }}
                >
                  {b.text}
                </motion.div>
              ))}
            </div>

            {/* Stats */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={visibleBroadcasts.length >= 4 ? { opacity: 1 } : {}}
              className="mt-6 p-3 bg-gold/5 border border-gold/20 rounded"
            >
              <div className="font-mono text-xs text-gold">
                📊 SWARM STATUS
              </div>
              <div className="grid grid-cols-2 gap-2 mt-2 text-xs text-text-secondary">
                <div>Direct teammates: <span className="text-text-primary">4</span></div>
                <div>Inner subagents: <span className="text-text-primary">5</span></div>
                <div>Effective parallelism: <span className="text-gold">~9 streams</span></div>
                <div>Speedup: <span className="text-gold">~3x</span></div>
              </div>
            </motion.div>
          </div>
        </div>

        {/* Dynamic restructuring operations */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mt-12 grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3"
        >
          {[
            { op: 'PROMOTE', desc: 'Leaf → STL', icon: '🔼', color: 'text-neon-green', why: 'When a leaf\'s domain grows to 3+ tasks that could run in parallel' },
            { op: 'SPLIT', desc: 'STL → 2 STLs', icon: '✂️', color: 'text-neon-blue', why: 'When an STL has 8+ tasks across 2 divergent sub-domains' },
            { op: 'DEMOTE', desc: 'STL → Leaf', icon: '🔽', color: 'text-neon-purple', why: 'When 1-2 tasks remain and inner parallelism is no longer needed' },
            { op: 'REASSIGN', desc: 'Move tasks', icon: '🔄', color: 'text-gold', why: 'When load is imbalanced or a teammate is blocked' },
            { op: 'SCALE-UP', desc: 'Spawn new teammate', icon: '🆕', color: 'text-neon-green', why: 'When next wave is wider than current team capacity' },
            { op: 'SCALE-DOWN', desc: 'Let teammate finish', icon: '🛬', color: 'text-text-muted', why: 'When waves narrow and a teammate would be idle' },
          ].map((item, i) => (
            <motion.div
              key={item.op}
              initial={{ opacity: 0, scale: 0.9 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.05 }}
              className="bg-bg-card border border-border rounded-lg p-3 text-center hover:border-gold/30 transition-colors"
            >
              <div className="text-lg mb-1">{item.icon}</div>
              <div className={`font-mono text-xs font-bold ${item.color}`}>{item.op}</div>
              <div className="text-xs text-text-muted mt-1">{item.desc}</div>
              <div className="text-[10px] text-text-muted/60 mt-1 italic leading-tight">{item.why}</div>
            </motion.div>
          ))}
        </motion.div>

        {/* Hierarchy Decision Tree */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5 }}
          className="mt-12"
        >
          <h3 className="font-mono text-lg font-bold text-gold mb-4">Hierarchy Decision Tree</h3>
          <div className="grid md:grid-cols-2 gap-4">
            <div className="bg-bg-card border border-border rounded-lg p-5" style={{ borderLeftColor: '#3b82f6', borderLeftWidth: 3 }}>
              <h4 className="font-mono text-sm font-bold text-neon-blue mb-3">Per-Cluster Assignment</h4>
              <div className="font-mono text-xs text-text-secondary space-y-1 bg-bg-dark rounded p-4 border border-border/30">
                <div className="text-text-muted">{'for each cluster of related tasks:'}</div>
                <div className="pl-4">
                  <div>{'if cluster.size >= 3'}</div>
                  <div>{'   AND cluster.can_parallelize_internally:'}</div>
                  <div className="text-neon-green">{'    → Spawn STL for this cluster'}</div>
                </div>
                <div className="pl-4">
                  <div>{'elif cluster.size >= 1:'}</div>
                  <div className="text-neon-purple">{'    → Assign to leaf teammate'}</div>
                </div>
              </div>
            </div>
            <div className="bg-bg-card border border-border rounded-lg p-5" style={{ borderLeftColor: '#f59e0b', borderLeftWidth: 3 }}>
              <h4 className="font-mono text-sm font-bold text-gold mb-3">Overall Mission Strategy</h4>
              <div className="font-mono text-xs text-text-secondary space-y-1 bg-bg-dark rounded p-4 border border-border/30">
                <div className="text-text-muted">{'for the overall mission:'}</div>
                <div className="pl-4">
                  <div>{'if num_domains >= 3:'}</div>
                  <div className="text-gold">{'    → You become pure orchestrator'}</div>
                </div>
                <div className="pl-4">
                  <div>{'elif num_domains == 2:'}</div>
                  <div className="text-gold">{'    → You take critical path + orchestrate'}</div>
                </div>
                <div className="pl-4">
                  <div>{'elif num_domains == 1:'}</div>
                  <div className="text-gold">{'    → You + 1 STL (or flat team if <3 tasks)'}</div>
                </div>
              </div>
            </div>
          </div>
        </motion.div>

        {/* Team Lead Principles */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5 }}
          className="mt-12"
        >
          <h3 className="font-mono text-lg font-bold text-gold mb-4">The 11 Team Lead Principles</h3>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-3">
            {[
              { num: 1, text: 'Maximize swarm throughput, not personal output', color: '#f59e0b' },
              { num: 2, text: 'Never idle while teammates work', color: '#f59e0b' },
              { num: 3, text: 'Respond to STL escalations INSTANTLY', color: '#ef4444' },
              { num: 4, text: 'Delegate domains, not tasks', color: '#3b82f6' },
              { num: 5, text: 'Don\'t micromanage STLs', color: '#3b82f6' },
              { num: 6, text: 'Micromanage less as hierarchy deepens', color: '#3b82f6' },
              { num: 7, text: 'Merge carefully when teammates work on overlapping files', color: '#a855f7' },
              { num: 8, text: 'Scale the hierarchy dynamically', color: '#22c55e' },
              { num: 9, text: 'Broadcast everything including hierarchy changes', color: '#22c55e' },
              { num: 10, text: 'Think in waves, delegate in domains', color: '#06b6d4' },
              { num: 11, text: 'Prefer STLs over flat teammates for complex missions', color: '#06b6d4' },
            ].map((principle, i) => (
              <motion.div
                key={principle.num}
                initial={{ opacity: 0, scale: 0.95 }}
                whileInView={{ opacity: 1, scale: 1 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.04 }}
                className="bg-bg-card border border-border rounded-lg p-3 flex items-start gap-3 hover:bg-bg-card-hover transition-colors"
              >
                <span
                  className="font-mono text-xs font-bold w-6 h-6 rounded-full flex items-center justify-center shrink-0 border"
                  style={{ borderColor: principle.color, color: principle.color, backgroundColor: principle.color + '15' }}
                >
                  {principle.num}
                </span>
                <span className="text-sm text-text-secondary leading-snug">{principle.text}</span>
              </motion.div>
            ))}
          </div>
        </motion.div>

        {/* STL Contract */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5 }}
          className="mt-12"
        >
          <h3 className="font-mono text-lg font-bold text-gold mb-4">Sub-Team-Lead Contract</h3>
          <div className="bg-bg-card border border-border rounded-lg p-5" style={{ borderLeftColor: '#3b82f6', borderLeftWidth: 3 }}>
            <p className="text-sm text-text-secondary mb-4">Every STL receives a contract with exactly 5 elements:</p>
            <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
              {[
                { label: 'Domain Scope', desc: 'Clear boundary of what this STL owns', icon: '📐', color: '#3b82f6' },
                { label: 'Parallelism Budget', desc: '2-4 inner subagents allowed', icon: '🔢', color: '#22c55e' },
                { label: 'Quality Gate', desc: 'Criteria that must pass before reporting done', icon: '✅', color: '#f59e0b' },
                { label: 'Reporting Protocol', desc: 'When and how to update the Team Lead', icon: '📡', color: '#a855f7' },
                { label: 'Escalation Path', desc: 'What to do when blocked or uncertain', icon: '🚨', color: '#ef4444' },
              ].map((item, i) => (
                <motion.div
                  key={item.label}
                  initial={{ opacity: 0, y: 10 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: i * 0.08 }}
                  className="bg-bg-dark rounded-lg p-3 border border-border/30 text-center"
                >
                  <div className="text-xl mb-2">{item.icon}</div>
                  <div className="font-mono text-xs font-bold mb-1" style={{ color: item.color }}>{item.label}</div>
                  <div className="text-[11px] text-text-muted leading-snug">{item.desc}</div>
                </motion.div>
              ))}
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  )
}
