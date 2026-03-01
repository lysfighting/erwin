import { motion, useScroll, useSpring } from 'framer-motion'
import { useState, useEffect } from 'react'
import { useLang, type Lang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

const sectionIds = [
  { id: 'hero',         key: 'nav.home' as const },
  { id: 'demo',         key: 'nav.demo' as const },
  { id: 'mechanics',    key: 'nav.mechanics' as const },
  { id: 'architecture', key: 'nav.architecture' as const },
  { id: 'phases',       key: 'nav.phases' as const },
  { id: 'swarm',        key: 'nav.swarm' as const },
  { id: 'waves',        key: 'nav.waves' as const },
  { id: 'six-steps',    key: 'nav.six-steps' as const },
  { id: 'philosophy',           key: 'nav.philosophy' as const },
  { id: 'deep-reasoning',       key: 'nav.deep-reasoning' as const },
  { id: 'direction-validation', key: 'nav.direction-validation' as const },
  { id: 'pantheon',             key: 'nav.pantheon' as const },
  { id: 'self-review',  key: 'nav.self-review' as const },
  { id: 'anti-patterns', key: 'nav.antipatterns' as const },
  { id: 'self-drive',   key: 'nav.self-drive' as const },
  { id: 'architect',    key: 'nav.architect' as const },
  { id: 'prompt-eng',   key: 'nav.prompt-eng' as const },
  { id: 'references',   key: 'nav.references' as const },
]

export function NavBar() {
  const [activeSection, setActiveSection] = useState('hero')
  const { scrollYProgress } = useScroll()
  const scaleX = useSpring(scrollYProgress, { stiffness: 100, damping: 30 })
  const { lang, setLang } = useLang()

  useEffect(() => {
    const observers: IntersectionObserver[] = []
    sectionIds.forEach(({ id }) => {
      const el = document.getElementById(id)
      if (!el) return
      const observer = new IntersectionObserver(
        ([entry]) => { if (entry.isIntersecting) setActiveSection(id) },
        { rootMargin: '-40% 0px -50% 0px' }
      )
      observer.observe(el)
      observers.push(observer)
    })
    return () => observers.forEach((o) => o.disconnect())
  }, [])

  return (
    <motion.nav
      initial={{ y: -60 }}
      animate={{ y: 0 }}
      className="fixed top-0 left-0 right-0 z-50 bg-bg-dark/80 backdrop-blur-md border-b border-border"
    >
      {/* Scroll progress bar */}
      <motion.div
        className="absolute top-0 left-0 right-0 h-[2px] bg-gold origin-left"
        style={{ scaleX }}
      />
      <div className="max-w-7xl mx-auto px-4 py-2.5 flex items-center gap-2">
        <a href="#hero" className="font-mono text-gold font-bold text-sm hover:text-gold-light transition-colors flex-shrink-0">
          DRIVE SKILL
        </a>

        {/* Section links — scrollable on smaller screens */}
        <div className="flex gap-0.5 overflow-x-auto flex-1 px-1 scrollbar-hide">
          {sectionIds.map((s) => (
            <a
              key={s.id}
              href={`#${s.id}`}
              className={`flex-shrink-0 px-2.5 py-1 text-xs rounded transition-all font-mono ${
                activeSection === s.id
                  ? 'text-gold bg-gold/10'
                  : 'text-text-secondary hover:text-text-primary hover:bg-bg-card'
              }`}
            >
              {t(s.key, lang)}
            </a>
          ))}
        </div>

        {/* Language toggle */}
        <div className="flex-shrink-0 flex items-center gap-1 border border-border rounded-full px-1 py-0.5">
          {(['en', 'zh'] as Lang[]).map((l) => (
            <button
              key={l}
              onClick={() => setLang(l)}
              className={`px-2 py-0.5 text-xs font-mono rounded-full transition-all ${
                lang === l
                  ? 'bg-gold text-bg-dark font-bold'
                  : 'text-text-secondary hover:text-text-primary'
              }`}
            >
              {l === 'en' ? 'EN' : '中'}
            </button>
          ))}
        </div>
      </div>
    </motion.nav>
  )
}
