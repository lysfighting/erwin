import { motion, useScroll, useSpring } from 'framer-motion'
import { useState, useEffect } from 'react'
import { useLang, type Lang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

const actLinks = [
  { id: 'act-i',   key: 'nav.act-i' as const },
  { id: 'act-ii',  key: 'nav.act-ii' as const },
  { id: 'act-iii', key: 'nav.act-iii' as const },
  { id: 'act-iv',  key: 'nav.act-iv' as const },
]


export function NavBar() {
  const [activeAct, setActiveAct] = useState('act-i')
  const { scrollYProgress } = useScroll()
  const scaleX = useSpring(scrollYProgress, { stiffness: 100, damping: 30 })
  const { lang, setLang } = useLang()

  useEffect(() => {
    const observers: IntersectionObserver[] = []
    actLinks.forEach(({ id }) => {
      const el = document.getElementById(id)
      if (!el) return
      const observer = new IntersectionObserver(
        ([entry]) => { if (entry.isIntersecting) setActiveAct(id) },
        { rootMargin: '-30% 0px -60% 0px' }
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

        {/* Act navigation links */}
        <div className="flex gap-1 flex-1 px-1">
          {actLinks.map((act) => (
            <a
              key={act.id}
              href={`#${act.id}`}
              className={`flex-shrink-0 px-3 py-1 text-xs rounded transition-all font-mono ${
                activeAct === act.id
                  ? 'text-gold bg-gold/10'
                  : 'text-text-secondary hover:text-text-primary hover:bg-bg-card'
              }`}
            >
              {t(act.key, lang)}
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
