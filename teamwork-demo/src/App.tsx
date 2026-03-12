import { useState, useEffect, useCallback, lazy, Suspense } from 'react'
import { LangProvider, useLang } from './i18n/LangContext'
import translations from './i18n/translations'
import LangToggle from './components/LangToggle'
import ChapterLayout from './components/ChapterLayout'

// Lazy load chapters (will be replaced with real components)
const Ch0Why = lazy(() => import('./chapters/Ch0Why'))
const Ch1GitModel = lazy(() => import('./chapters/Ch1GitModel'))
const Ch2Lifecycle = lazy(() => import('./chapters/Ch2Lifecycle'))
const Ch3Setup = lazy(() => import('./chapters/Ch3Setup'))
const Ch4Create = lazy(() => import('./chapters/Ch4Create'))
const Ch5Claim = lazy(() => import('./chapters/Ch5Claim'))
const Ch6Execute = lazy(() => import('./chapters/Ch6Execute'))
const Ch7Deliver = lazy(() => import('./chapters/Ch7Deliver'))
const Ch8Shortcuts = lazy(() => import('./chapters/Ch8Shortcuts'))
const Ch9Health = lazy(() => import('./chapters/Ch9Health'))
const Ch10Release = lazy(() => import('./chapters/Ch10Release'))
const Ch11HelpLearn = lazy(() => import('./chapters/Ch11HelpLearn'))
const Ch12Worktree = lazy(() => import('./chapters/Ch12Worktree'))

const chapterComponents = [
  Ch0Why, Ch1GitModel, Ch2Lifecycle, Ch3Setup, Ch4Create,
  Ch5Claim, Ch6Execute, Ch7Deliver, Ch8Shortcuts, Ch9Health, Ch10Release,
  Ch11HelpLearn, Ch12Worktree,
]

function Placeholder({ title, subtitle }: { title: string; subtitle: string }) {
  return (
    <ChapterLayout title={title} subtitle={subtitle}>
      <div className="text-text-muted text-center py-16">Coming soon...</div>
    </ChapterLayout>
  )
}

function AppContent() {
  const [current, setCurrent] = useState(0)
  const { lang } = useLang()
  const t = translations[lang]

  const chapters = t.chapters.map((ch, i) => ({ id: i, title: ch.title }))
  const ChapterComponent = chapterComponents[current]

  const goPrev = useCallback(() => setCurrent(c => Math.max(0, c - 1)), [])
  const goNext = useCallback(() => setCurrent(c => Math.min(chapterComponents.length - 1, c + 1)), [])

  useEffect(() => {
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === 'ArrowLeft') goPrev()
      if (e.key === 'ArrowRight') goNext()
    }
    window.addEventListener('keydown', handleKey)
    return () => window.removeEventListener('keydown', handleKey)
  }, [goPrev, goNext])

  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-bg-dark">
      <LangToggle />
      {/* Sidebar stepper — desktop */}
      <aside className="hidden md:flex flex-col w-56 shrink-0 border-r border-border">
        <div className="px-4 pt-6 pb-2">
          <h2 className="text-sm font-bold text-text-primary tracking-wider uppercase">{t.appTitle}</h2>
          <p className="text-xs text-text-muted mt-1">{t.appSubtitle}</p>
        </div>
        <nav className="flex-1 flex flex-col gap-0.5 px-2 py-2 overflow-y-auto">
          {chapters.map((ch, i) => (
            <button
              key={ch.id}
              onClick={() => setCurrent(i)}
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
      </aside>

      {/* Main content */}
      <main className="flex-1 flex flex-col min-h-screen">
        {/* Mobile chapter dots */}
        <div className="md:hidden flex items-center gap-1.5 px-4 py-3 border-b border-border overflow-x-auto">
          {chapters.map((_, i) => (
            <button
              key={i}
              onClick={() => setCurrent(i)}
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

        {/* Chapter content */}
        <div className="flex-1 overflow-y-auto">
          <Suspense fallback={
            <Placeholder title={t.chapters[current].title} subtitle={t.chapters[current].subtitle} />
          }>
            <ChapterComponent />
          </Suspense>
        </div>

        {/* Bottom nav */}
        <div className="flex items-center justify-between px-6 py-4 border-t border-border">
          <button
            onClick={() => setCurrent(Math.max(0, current - 1))}
            disabled={current === 0}
            className="px-4 py-2 text-sm font-mono rounded border border-border bg-bg-card text-text-primary hover:bg-bg-card-hover disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
          >
            ◀ {t.prev}
          </button>
          <span className="text-sm text-text-muted font-mono">
            {current + 1} / {chapters.length}
          </span>
          <button
            onClick={() => setCurrent(Math.min(chapters.length - 1, current + 1))}
            disabled={current === chapters.length - 1}
            className="px-4 py-2 text-sm font-mono rounded border border-border bg-neon-blue text-white hover:opacity-80 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
          >
            {t.next} ▶
          </button>
        </div>
      </main>
    </div>
  )
}

export default function App() {
  return (
    <LangProvider>
      <AppContent />
    </LangProvider>
  )
}
