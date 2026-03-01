import { motion, AnimatePresence } from 'framer-motion'
import { useState, useEffect, useRef, useCallback } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'
import { Annotatable } from '../components/Annotatable'

interface TerminalLine {
  text: string
  type: 'input' | 'phase' | 'dispatch' | 'return' | 'topology' | 'broadcast' | 'system' | 'complete' | 'dim'
  delay: number
}

const lines: TerminalLine[] = [
  // User input
  { text: '$ /drive implement user authentication with OAuth2 + JWT', type: 'input', delay: 0 },
  { text: '', type: 'dim', delay: 400 },

  // Phase 0
  { text: '═══════════════════════════════════════════════════════════', type: 'dim', delay: 600 },
  { text: '  PHASE 0: DEEP BRIEFING — Research Swarm Activated', type: 'phase', delay: 700 },
  { text: '═══════════════════════════════════════════════════════════', type: 'dim', delay: 750 },
  { text: '', type: 'dim', delay: 800 },
  { text: '📡 PARALLEL DISPATCH (Research Swarm — 3 scouts):', type: 'broadcast', delay: 900 },
  { text: '   [scout-codebase]  → Analyze existing auth structure', type: 'dispatch', delay: 1050 },
  { text: '   [scout-patterns]  → Research OAuth2/JWT patterns in codebase', type: 'dispatch', delay: 1200 },
  { text: '   [scout-security]  → Audit security requirements & risks', type: 'dispatch', delay: 1350 },
  { text: '   Parallel width: 3 | Est. briefing time: ~2 min', type: 'dim', delay: 1450 },
  { text: '', type: 'dim', delay: 1500 },
  { text: '✅ RETURN: [scout-codebase]   ← No existing auth layer found', type: 'return', delay: 3200 },
  { text: '✅ RETURN: [scout-patterns]   ← passport.js + jsonwebtoken available', type: 'return', delay: 3500 },
  { text: '✅ RETURN: [scout-security]   ← PKCE flow required, refresh tokens needed', type: 'return', delay: 3700 },
  { text: '', type: 'dim', delay: 3750 },
  { text: '📋 MISSION CONTRACT:', type: 'phase', delay: 3900 },
  { text: '   Goal: OAuth2 (Google) + JWT session management + PKCE', type: 'dim', delay: 4000 },
  { text: '   Scope: src/auth/ — 6 new files, ~400 lines', type: 'dim', delay: 4100 },
  { text: '   Quality gate: E2E tests passing, security review ✓', type: 'dim', delay: 4200 },
  { text: '', type: 'dim', delay: 4300 },

  // Phase T
  { text: '═══════════════════════════════════════════════════════════', type: 'dim', delay: 4500 },
  { text: '  PHASE T: TEAM ASSEMBLY', type: 'phase', delay: 4600 },
  { text: '═══════════════════════════════════════════════════════════', type: 'dim', delay: 4650 },
  { text: '', type: 'dim', delay: 4700 },
  { text: '🗺️  TEAM TOPOLOGY:', type: 'topology', delay: 4800 },
  { text: '   YOU (Team Lead) ─── orchestrating', type: 'topology', delay: 4900 },
  { text: '   ├── 🎖️  stl-auth    ─── owns OAuth2 + JWT [3 inner subagents]', type: 'topology', delay: 5050 },
  { text: '   ├── 🎖️  stl-api     ─── owns routes + middleware [2 inner subagents]', type: 'topology', delay: 5200 },
  { text: '   └── 🔧  tester      ─── writes tests, verifies E2E', type: 'topology', delay: 5350 },
  { text: '   Active: 3 teammates + ~5 inner subagents = 8 parallel streams', type: 'topology', delay: 5500 },
  { text: '', type: 'dim', delay: 5550 },

  // Phase 1
  { text: '═══════════════════════════════════════════════════════════', type: 'dim', delay: 5700 },
  { text: '  PHASE 1: BATTLE PLAN — Wave Decomposition', type: 'phase', delay: 5800 },
  { text: '═══════════════════════════════════════════════════════════', type: 'dim', delay: 5850 },
  { text: '', type: 'dim', delay: 5900 },
  { text: '📡 PARALLEL DISPATCH (Wave 1 — Bootstrap):', type: 'broadcast', delay: 6000 },
  { text: '   [stl-auth]  → OAuth2 provider setup + token service', type: 'dispatch', delay: 6150 },
  { text: '   [stl-api]   → Route scaffolding + Express middleware', type: 'dispatch', delay: 6300 },
  { text: '   [tester]    → Test skeleton + mock OAuth server', type: 'dispatch', delay: 6450 },
  { text: '', type: 'dim', delay: 6500 },
  { text: '🎖️  STL-INNER: [stl-auth] spawned 3 inner subagents', type: 'broadcast', delay: 7800 },
  { text: '   [auth/oauth-provider]   → Google OAuth2 + PKCE flow', type: 'dispatch', delay: 7900 },
  { text: '   [auth/jwt-service]      → Token generation + refresh logic', type: 'dispatch', delay: 8050 },
  { text: '   [auth/session-store]    → Redis session management', type: 'dispatch', delay: 8200 },
  { text: '', type: 'dim', delay: 8250 },
  { text: '✅ RETURN: [auth/oauth-provider]  ← PKCE flow implemented', type: 'return', delay: 10500 },
  { text: '✅ RETURN: [auth/jwt-service]     ← Tokens + refresh done', type: 'return', delay: 10800 },
  { text: '✅ RETURN: [auth/session-store]   ← Redis integration complete', type: 'return', delay: 11000 },
  { text: '🏁 DOMAIN COMPLETE: [stl-auth] ← Auth layer fully delivered', type: 'complete', delay: 11300 },
  { text: '   Tasks: 3/3 ✓ | Inner subagents: 3 | Files: 6 touched', type: 'dim', delay: 11400 },
  { text: '', type: 'dim', delay: 11450 },

  // Wave 2
  { text: '📡 PARALLEL DISPATCH (Wave 2 — Integration):', type: 'broadcast', delay: 11600 },
  { text: '   [stl-api]   → Wire auth middleware into all protected routes', type: 'dispatch', delay: 11750 },
  { text: '   [tester]    → Integration tests: login flow, token refresh, logout', type: 'dispatch', delay: 11900 },
  { text: '', type: 'dim', delay: 12000 },
  { text: '✅ RETURN: [stl-api]   ← All 12 routes protected', type: 'return', delay: 13800 },
  { text: '✅ RETURN: [tester]    ← 28 tests passing, 0 failing', type: 'return', delay: 14100 },
  { text: '', type: 'dim', delay: 14150 },

  // Mission complete
  { text: '═══════════════════════════════════════════════════════════', type: 'dim', delay: 14300 },
  { text: '  MISSION COMPLETE ✓', type: 'complete', delay: 14400 },
  { text: '═══════════════════════════════════════════════════════════', type: 'dim', delay: 14450 },
  { text: '📊 SWARM DEBRIEF:', type: 'broadcast', delay: 14550 },
  { text: '   Waves: 2 | Wall-clock time: ~14 min | Serial equiv: ~45 min', type: 'dim', delay: 14650 },
  { text: '   Teammates: 3 | Inner subagents: 5 | Files touched: 11', type: 'dim', delay: 14750 },
  { text: '   Speedup: 3.2x | Tests: 28/28 passing | Coverage: 87%', type: 'dim', delay: 14850 },
]

const lineColors: Record<TerminalLine['type'], string> = {
  input: 'text-white',
  phase: 'text-gold font-bold',
  dispatch: 'text-neon-blue',
  return: 'text-neon-green',
  topology: 'text-neon-purple',
  broadcast: 'text-neon-cyan font-semibold',
  system: 'text-text-secondary',
  complete: 'text-neon-green font-bold',
  dim: 'text-text-muted',
}

const phases = [
  { id: 0, label: 'Phase 0', desc: 'Deep Briefing', color: '#f59e0b', lineRange: [0, 15] },
  { id: 1, label: 'Phase T', desc: 'Team Assembly', color: '#a855f7', lineRange: [15, 23] },
  { id: 2, label: 'Phase 1', desc: 'Battle Plan', color: '#3b82f6', lineRange: [23, 42] },
  { id: 3, label: 'Wave 2', desc: 'Integration', color: '#06b6d4', lineRange: [42, 48] },
  { id: 4, label: 'Complete', desc: 'Mission Done', color: '#22c55e', lineRange: [48, 55] },
]

export function DemoSection() {
  const { lang } = useLang()
  const [visibleLines, setVisibleLines] = useState<number>(0)
  const [isPlaying, setIsPlaying] = useState(false)
  const [activePhase, setActivePhase] = useState(0)
  const terminalRef = useRef<HTMLDivElement>(null)
  const timersRef = useRef<ReturnType<typeof setTimeout>[]>([])
  const isInViewRef = useRef(false)

  const startAnimation = useCallback(() => {
    timersRef.current.forEach(clearTimeout)
    timersRef.current = []
    setVisibleLines(0)
    setActivePhase(0)
    setIsPlaying(true)

    lines.forEach((line, i) => {
      const t = setTimeout(() => {
        setVisibleLines(i + 1)
        // Update active phase
        for (const phase of phases) {
          if (i >= phase.lineRange[0] && i < phase.lineRange[1]) {
            setActivePhase(phase.id)
            break
          }
        }
        // Auto-scroll terminal
        if (terminalRef.current) {
          terminalRef.current.scrollTop = terminalRef.current.scrollHeight
        }
        if (i === lines.length - 1) setIsPlaying(false)
      }, line.delay)
      timersRef.current.push(t)
    })
  }, [])

  // Auto-start when scrolled into view
  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting && !isInViewRef.current) {
          isInViewRef.current = true
          setTimeout(startAnimation, 500)
        }
      },
      { rootMargin: '-20% 0px -20% 0px' }
    )
    const el = document.getElementById('demo')
    if (el) observer.observe(el)
    return () => observer.disconnect()
  }, [startAnimation])

  return (
    <section id="demo" className="py-24 px-4">
      <div className="max-w-6xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 40 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
        >
          <h2 className="text-4xl font-bold mb-2">{t('demo.title', lang)}</h2>
          <p className="text-text-secondary text-lg mb-2">
            {t('demo.subtitle', lang)}
          </p>
          <p className="text-text-muted text-sm mb-8">
            {t('demo.desc', lang)}
          </p>
        </motion.div>

        {/* Phase indicator */}
        <div className="flex gap-1 mb-4 overflow-x-auto pb-2">
          {phases.map((phase) => (
            <motion.div
              key={phase.id}
              className="flex-shrink-0 flex items-center gap-2 px-3 py-1.5 rounded-full border text-xs font-mono transition-all"
              style={{
                backgroundColor: activePhase >= phase.id ? phase.color + '15' : 'transparent',
                borderColor: activePhase >= phase.id ? phase.color + '60' : '#30363d',
                color: activePhase >= phase.id ? phase.color : '#484f58',
              }}
            >
              {activePhase > phase.id && <span className="text-[10px]">✓</span>}
              {activePhase === phase.id && (
                <motion.span
                  animate={{ opacity: [1, 0, 1] }}
                  transition={{ duration: 0.8, repeat: Infinity }}
                  className="w-1.5 h-1.5 rounded-full"
                  style={{ backgroundColor: phase.color }}
                />
              )}
              <span>{phase.label}</span>
              <span className="text-text-muted">{phase.desc}</span>
            </motion.div>
          ))}
          <div className="ml-auto flex-shrink-0">
            <button
              onClick={startAnimation}
              disabled={isPlaying}
              className="px-3 py-1.5 text-xs font-mono border border-border rounded-full hover:border-gold/40 hover:text-gold text-text-muted transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
            >
              {isPlaying ? t('demo.playing', lang) : t('demo.replay', lang)}
            </button>
          </div>
        </div>

        {/* Terminal window */}
        <div className="rounded-lg border border-border overflow-hidden shadow-2xl">
          {/* Window chrome */}
          <div className="flex items-center gap-2 px-4 py-2.5 bg-[#1c1c1e] border-b border-border">
            <div className="w-3 h-3 rounded-full bg-[#ff5f57]" />
            <div className="w-3 h-3 rounded-full bg-[#febc2e]" />
            <div className="w-3 h-3 rounded-full bg-[#28c840]" />
            <span className="mx-auto font-mono text-xs text-text-muted">Claude Code — my-project</span>
          </div>

          {/* Terminal content */}
          <div
            ref={terminalRef}
            className="bg-[#0d1117] font-mono text-xs leading-relaxed p-4 overflow-y-auto"
            style={{ height: '480px' }}
          >
            <AnimatePresence initial={false}>
              {lines.slice(0, visibleLines).map((line, i) => (
                <motion.div
                  key={i}
                  initial={{ opacity: 0, x: -4 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ duration: 0.15 }}
                  className={`whitespace-pre-wrap ${lineColors[line.type]}`}
                >
                  {line.text || '\u00A0'}
                </motion.div>
              ))}
            </AnimatePresence>
            {isPlaying && (
              <motion.span
                animate={{ opacity: [1, 0, 1] }}
                transition={{ duration: 0.6, repeat: Infinity }}
                className="text-gold"
              >
                █
              </motion.span>
            )}
          </div>
        </div>

        {/* Phase Design Reasoning — 5 insights */}
        <div className="mt-8">
          <h3 className="font-mono text-sm text-text-muted mb-4">
            {lang === 'zh' ? '为什么这样设计？— 每个阶段背后的决策逻辑' : 'Why each phase? — The design reasoning behind every step'}
          </h3>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
            {[
              {
                phase: 'Phase 0',
                color: '#f59e0b',
                title: lang === 'zh' ? '为什么先调研，后写代码？' : 'Why research before touching code?',
                problem: lang === 'zh' ? '错误方向全速前进 = 最大浪费。不调研就开始实现，就像在迷宫里全速奔跑 — 速度越快，浪费越多。' : 'Wrong direction at full speed = maximum waste. Starting without research is sprinting through a maze — the faster you go, the more time you waste.',
                insight: lang === 'zh' ? '3个并行的侦察 agent 用1分钟覆盖了需要15分钟的串行调研。Psychiatrist Test：你能以用户的口吻描述他们的愿景吗？' : '3 parallel scouts cover in 1 minute what would take 15 serially. The Psychiatrist Test: can you describe the user\'s vision AS IF YOU WERE THEM?',
                id: 'demo-insight-0',
              },
              {
                phase: 'Phase T',
                color: '#a855f7',
                title: lang === 'zh' ? '为什么需要组建团队？' : 'Why assemble a team at all?',
                problem: lang === 'zh' ? 'Claude 只有一个 context window。单 agent 等于单线程 — 所有任务串行，所有代码挤在同一个上下文里，越来越慢。' : 'Claude has one context window. Solo = single-threaded. All tasks serial, all code in one context, getting slower as complexity grows.',
                insight: lang === 'zh' ? 'Agent Teams API 让每个 teammate 拥有独立的 context window。STL-frontend 不知道也不需要知道 STL-backend 在做什么。这就是并行的本质。' : 'Agent Teams API gives each teammate an independent context window. STL-frontend doesn\'t know or need to know what STL-backend is doing. That\'s the essence of parallelism.',
                id: 'demo-insight-T',
              },
              {
                phase: 'Phase 1',
                color: '#3b82f6',
                title: lang === 'zh' ? '为什么用"Wave"而不是任务列表？' : 'Why "waves" instead of a task list?',
                problem: lang === 'zh' ? '任务列表暗含串行。写下"任务1、任务2、任务3"，默认就是按顺序来。没有人明确说"这三个可以同时做"，结果就没人同时做。' : 'A task list implies serial execution. Writing "Task 1, Task 2, Task 3" defaults to sequential. Nobody explicitly says "these three can run simultaneously", so nobody does.',
                insight: lang === 'zh' ? 'Wave 让并行变得显式。Wave 0 里的3个任务同时开始 — 没有歧义。Wave Map 就是执行图，不是愿望清单。' : 'Waves make parallelism explicit. 3 tasks in Wave 0 start simultaneously — no ambiguity. The Wave Map is an execution graph, not a wish list.',
                id: 'demo-insight-1',
              },
              {
                phase: 'Wave 2',
                color: '#06b6d4',
                title: lang === 'zh' ? '为什么 STL 可以自己产生子 agent？' : 'Why can STLs spawn their own subagents?',
                problem: lang === 'zh' ? '如果只有 Team Lead 能派遣任务，那 STL 本质上还是叶节点。Team Lead 管理所有30个子任务？协调开销吞噬并行收益。' : 'If only the Team Lead can dispatch tasks, STLs are glorified leaf nodes. The Team Lead managing all 30 sub-tasks? Coordination overhead eats the parallelism gains.',
                insight: lang === 'zh' ? 'STL 用 Task tool 在自己的 context 内启动子 agent — 不需要创建新 team，同步返回。Terminal 里按 Shift+Down 就能看到 STL 正在派遣子任务的实时画面。' : 'STLs use the Task tool to spawn subagents within their own context — no new team needed, synchronous return. Press Shift+Down in terminal to see an STL actively dispatching subtasks live.',
                id: 'demo-insight-2',
              },
              {
                phase: 'Mission Complete',
                color: '#22c55e',
                title: lang === 'zh' ? '为什么要"自我对抗审查"？' : 'Why a self-adversarial review?',
                problem: lang === 'zh' ? '"看起来对了"不等于完成。LLM 有一个已知缺陷：在完成80%后，剩下的边缘情况和错误处理往往被跳过。任务被标记为完成，但其实没有。' : '"Looks right" ≠ done. LLMs have a known failure mode: after completing 80%, the remaining edge cases and error handling get skipped. Tasks are marked complete but aren\'t.',
                insight: lang === 'zh' ? '自我对抗：用最挑剔的审查者视角审查自己的代码。"如果我要攻击这个实现，从哪里入手？" 只有当你找不到任何可以改进的地方，才算真正完成。' : 'Self-adversarial: review your own code as the harshest critic. "If I were trying to BREAK this, where would I attack?" Only when you can find nothing left to improve — truly done.',
                id: 'demo-insight-3',
              },
            ].map((item, i) => (
              <Annotatable
                key={item.phase}
                id={item.id}
                label={`Demo Phase: ${item.phase}`}
                sectionId="demo"
                sectionLabel="Live Demo"
              >
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: i * 0.08 }}
                  className="bg-bg-card border border-border rounded-lg p-4 hover:bg-bg-card-hover transition-colors"
                  style={{ borderTopColor: item.color, borderTopWidth: 2 }}
                >
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-xs font-mono px-2 py-0.5 rounded font-bold" style={{ backgroundColor: item.color + '20', color: item.color }}>
                      {item.phase}
                    </span>
                  </div>
                  <p className="text-sm font-bold text-text-primary mb-2">{item.title}</p>
                  <div className="space-y-2">
                    <div className="bg-neon-red/5 border border-neon-red/10 rounded p-2">
                      <p className="text-[10px] font-mono text-neon-red mb-1">PROBLEM</p>
                      <p className="text-xs text-text-secondary">{item.problem}</p>
                    </div>
                    <div className="bg-neon-green/5 border border-neon-green/10 rounded p-2">
                      <p className="text-[10px] font-mono text-neon-green mb-1">INSIGHT</p>
                      <p className="text-xs text-text-secondary">{item.insight}</p>
                    </div>
                  </div>
                </motion.div>
              </Annotatable>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
