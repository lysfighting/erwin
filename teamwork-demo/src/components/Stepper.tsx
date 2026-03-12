import { useEffect, useCallback } from 'react'

interface Chapter {
  id: number
  title: string
}

interface StepperProps {
  chapters: Chapter[]
  current: number
  onChange: (index: number) => void
  prevLabel?: string
  nextLabel?: string
}

export default function Stepper({ chapters, current, onChange, prevLabel = 'Prev', nextLabel = 'Next' }: StepperProps) {
  const goPrev = useCallback(() => {
    if (current > 0) onChange(current - 1)
  }, [current, onChange])

  const goNext = useCallback(() => {
    if (current < chapters.length - 1) onChange(current + 1)
  }, [current, chapters.length, onChange])

  useEffect(() => {
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === 'ArrowLeft') goPrev()
      if (e.key === 'ArrowRight') goNext()
    }
    window.addEventListener('keydown', handleKey)
    return () => window.removeEventListener('keydown', handleKey)
  }, [goPrev, goNext])

  return (
    <div className="flex flex-col h-full">
      {/* Sidebar — desktop */}
      <nav className="hidden md:flex flex-col gap-1 w-56 shrink-0 p-4 border-r border-border overflow-y-auto">
        {chapters.map((ch, i) => (
          <button
            key={ch.id}
            onClick={() => onChange(i)}
            className={`text-left text-sm px-3 py-2 rounded transition-colors ${
              i === current
                ? 'bg-neon-blue/15 text-neon-blue font-medium'
                : 'text-text-secondary hover:text-text-primary hover:bg-bg-card-hover'
            }`}
          >
            <span className="text-text-muted mr-1.5">{i}.</span>
            {ch.title}
          </button>
        ))}
      </nav>

      {/* Mobile top bar */}
      <div className="md:hidden flex items-center gap-1.5 px-4 py-3 border-b border-border overflow-x-auto">
        {chapters.map((_, i) => (
          <button
            key={i}
            onClick={() => onChange(i)}
            className={`w-7 h-7 shrink-0 rounded-full text-xs font-mono flex items-center justify-center transition-colors ${
              i === current
                ? 'bg-neon-blue text-white'
                : 'bg-bg-card text-text-muted hover:bg-bg-card-hover'
            }`}
          >
            {i}
          </button>
        ))}
      </div>

      {/* Bottom nav */}
      <div className="flex items-center justify-between px-6 py-4 border-t border-border mt-auto">
        <button
          onClick={goPrev}
          disabled={current === 0}
          className="px-4 py-2 text-sm font-mono rounded border border-border bg-bg-card text-text-primary hover:bg-bg-card-hover disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
        >
          ◀ {prevLabel}
        </button>
        <span className="text-sm text-text-muted font-mono">
          {current + 1} / {chapters.length}
        </span>
        <button
          onClick={goNext}
          disabled={current === chapters.length - 1}
          className="px-4 py-2 text-sm font-mono rounded border border-border bg-neon-blue text-white hover:opacity-80 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
        >
          {nextLabel} ▶
        </button>
      </div>
    </div>
  )
}
