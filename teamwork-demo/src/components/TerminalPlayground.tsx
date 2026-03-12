import { useState, useRef, useCallback, useEffect, type ReactNode } from 'react'
import { motion, AnimatePresence } from 'framer-motion'

export interface TerminalStep {
  command: string
  output: string[]
  delay?: number
}

interface TerminalPlaygroundProps {
  steps: TerminalStep[]
  title?: string
}

/**
 * Parse Claude Code visual encoding markers into styled JSX.
 *
 * Supported:
 *   **#NNN** → light-blue bold (issue auto-link)
 *   **text** → white bold (headers, labels)
 *   `text`   → purple-blue rgb(177,185,249) (commands, paths, counts)
 *   *text*   → dim italic (branch names, timestamps)
 *   ────     → rendered as-is (dividers)
 */
function parseRichLine(line: string): ReactNode {
  if (!line) return '\u00A0'

  const parts: ReactNode[] = []
  let remaining = line
  let key = 0

  while (remaining.length > 0) {
    // Try **#NNN** (bold issue link) — must come before generic **bold**
    let match = remaining.match(/^\*\*(#\d{2,})\*\*/)
    if (match) {
      parts.push(<span key={key++} className="text-blue-400 font-bold">{match[1]}</span>)
      remaining = remaining.slice(match[0].length)
      continue
    }

    // Try **bold**
    match = remaining.match(/^\*\*(.+?)\*\*/)
    if (match) {
      parts.push(<span key={key++} className="text-white font-bold">{match[1]}</span>)
      remaining = remaining.slice(match[0].length)
      continue
    }

    // Try `code`
    match = remaining.match(/^`([^`]+)`/)
    if (match) {
      parts.push(<span key={key++} className="text-[rgb(177,185,249)]">{match[1]}</span>)
      remaining = remaining.slice(match[0].length)
      continue
    }

    // Try *italic* (but not ** which is bold)
    match = remaining.match(/^\*([^*]+?)\*/)
    if (match) {
      parts.push(<span key={key++} className="italic text-text-muted">{match[1]}</span>)
      remaining = remaining.slice(match[0].length)
      continue
    }

    // Plain text — consume until next marker or end
    match = remaining.match(/^([\s\S]+?)(?=\*\*|\*[^*]|`|$)/)
    if (match && match[1]) {
      parts.push(<span key={key++}>{match[1]}</span>)
      remaining = remaining.slice(match[1].length)
      continue
    }

    // Fallback: consume one char to prevent infinite loop
    parts.push(<span key={key++}>{remaining[0]}</span>)
    remaining = remaining.slice(1)
  }

  return <>{parts}</>
}

export default function TerminalPlayground({ steps, title }: TerminalPlaygroundProps) {
  const [currentStep, setCurrentStep] = useState(0)
  const [isAnimating, setIsAnimating] = useState(false)
  const [visibleLines, setVisibleLines] = useState(0)
  const [animationComplete, setAnimationComplete] = useState(false)
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const scrollRef = useRef<HTMLDivElement>(null)

  const step = steps[currentStep]
  const totalLines = step.output.length
  const delay = step.delay ?? 50

  const clearTimer = useCallback(() => {
    if (timerRef.current) {
      clearTimeout(timerRef.current)
      timerRef.current = null
    }
  }, [])

  // Auto-scroll to bottom when new lines appear
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
  }, [visibleLines])

  const skipAnimation = useCallback(() => {
    if (!isAnimating) return
    clearTimer()
    setVisibleLines(totalLines)
    setIsAnimating(false)
    setAnimationComplete(true)
  }, [isAnimating, clearTimer, totalLines])

  const animateLines = useCallback(() => {
    if (totalLines === 0) {
      setIsAnimating(false)
      setAnimationComplete(true)
      return
    }
    setIsAnimating(true)
    setAnimationComplete(false)
    setVisibleLines(0)

    let lineIndex = 0
    const showNext = () => {
      lineIndex++
      setVisibleLines(lineIndex)
      if (lineIndex < totalLines) {
        timerRef.current = setTimeout(showNext, delay)
      } else {
        setIsAnimating(false)
        setAnimationComplete(true)
      }
    }
    timerRef.current = setTimeout(showNext, delay)
  }, [totalLines, delay])

  const goToStep = useCallback((index: number) => {
    clearTimer()
    setCurrentStep(index)
    setVisibleLines(0)
    setIsAnimating(false)
    setAnimationComplete(false)
  }, [clearTimer])

  // Cleanup on unmount
  useEffect(() => clearTimer, [clearTimer])

  return (
    <div className="w-full">
      {/* Terminal window */}
      <div className="rounded-lg border border-border overflow-hidden bg-bg-dark">
        {/* Header bar */}
        <div className="flex items-center gap-2 px-4 py-3 bg-bg-card border-b border-border">
          <div className="flex gap-1.5">
            <div className="w-3 h-3 rounded-full bg-neon-red" />
            <div className="w-3 h-3 rounded-full bg-gold" />
            <div className="w-3 h-3 rounded-full bg-neon-green" />
          </div>
          {title && (
            <span className="ml-2 text-sm text-text-secondary font-mono">{title}</span>
          )}
        </div>

        {/* Content area */}
        <div
          ref={scrollRef}
          className="p-4 font-mono text-sm max-h-96 overflow-y-auto cursor-pointer"
          onClick={skipAnimation}
        >
          {/* Command line */}
          <div className="flex gap-2">
            <span className="text-neon-green select-none">$</span>
            <span className="text-text-primary">{step.command}</span>
          </div>

          {/* Output lines */}
          <AnimatePresence mode="sync">
            {step.output.slice(0, visibleLines).map((line, i) => (
              <motion.div
                key={`${currentStep}-${i}`}
                initial={{ opacity: 0, y: 4 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.15 }}
                className="text-text-secondary mt-0.5"
              >
                {parseRichLine(line)}
              </motion.div>
            ))}
          </AnimatePresence>
        </div>
      </div>

      {/* Controls */}
      <div className="flex items-center justify-between mt-3">
        <div className="flex items-center gap-2">
          <button
            onClick={() => goToStep(currentStep - 1)}
            disabled={currentStep === 0}
            className="px-3 py-1.5 text-sm font-mono rounded border border-border bg-bg-card text-text-primary hover:bg-bg-card-hover disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
          >
            ◀ Prev
          </button>
          <button
            onClick={() => goToStep(currentStep + 1)}
            disabled={currentStep === steps.length - 1}
            className={`px-3 py-1.5 text-sm font-mono rounded border border-border bg-bg-card text-text-primary hover:bg-bg-card-hover disabled:opacity-30 disabled:cursor-not-allowed transition-colors ${
              animationComplete && currentStep < steps.length - 1
                ? 'ring-1 ring-neon-green/50'
                : ''
            }`}
          >
            Next ▶
          </button>
        </div>

        <span className="text-sm text-text-muted font-mono">
          Step {currentStep + 1} / {steps.length}
        </span>

        <button
          onClick={animateLines}
          disabled={isAnimating}
          className={`px-4 py-1.5 text-sm font-mono rounded border transition-colors ${
            isAnimating
              ? 'border-border bg-bg-card text-text-secondary cursor-not-allowed'
              : 'border-neon-green bg-neon-green text-white hover:brightness-110'
          }`}
        >
          {isAnimating ? 'Running...' : '▶ Run'}
        </button>
      </div>
    </div>
  )
}

// Usage:
// <TerminalPlayground steps={[
//   { command: 'echo "hello"', output: ['hello'] },
//   { command: 'ls -la', output: ['total 32', 'drwxr-xr-x  5 user  staff  160 Jan 1 00:00 .', '-rw-r--r--  1 user  staff  256 Jan 1 00:00 package.json'] }
// ]} />
