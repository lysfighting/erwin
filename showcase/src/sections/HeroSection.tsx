import { motion, AnimatePresence } from 'framer-motion'
import { useState, useEffect, useMemo } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

function HeroAnimation({ lang }: { lang: 'en' | 'zh' }) {
  const [phase, setPhase] = useState<'typing' | 'shatter' | 'reform'>('typing')
  const [lineIndex, setLineIndex] = useState(0)
  const [charIndex, setCharIndex] = useState(0)

  const typewriterTexts = useMemo(() => [
    t('hero.line1', lang),
    t('hero.line2', lang),
    t('hero.line3', lang),
  ], [lang])

  useEffect(() => {
    if (phase !== 'typing') return
    const currentLine = typewriterTexts[lineIndex]
    if (!currentLine) {
      const timer = setTimeout(() => setPhase('shatter'), 800)
      return () => clearTimeout(timer)
    }
    if (charIndex < currentLine.length) {
      const timer = setTimeout(() => setCharIndex(c => c + 1), 60)
      return () => clearTimeout(timer)
    } else {
      const timer = setTimeout(() => {
        setLineIndex(i => i + 1)
        setCharIndex(0)
      }, 600)
      return () => clearTimeout(timer)
    }
  }, [phase, lineIndex, charIndex, typewriterTexts])

  useEffect(() => {
    if (phase === 'shatter') {
      const timer = setTimeout(() => setPhase('reform'), 1200)
      return () => clearTimeout(timer)
    }
  }, [phase])

  const badges = [
    { key: 'hero.badge.swarm' as const, icon: '🐝' },
    { key: 'hero.badge.waves' as const, icon: '🌊' },
    { key: 'hero.badge.sixstep' as const, icon: '📐' },
    { key: 'hero.badge.broadcast' as const, icon: '📡' },
    { key: 'hero.badge.pantheon' as const, icon: '🏛️' },
  ]

  return (
    <>
      {/* Background grid */}
      <div className="absolute inset-0 opacity-5">
        <div className="w-full h-full" style={{
          backgroundImage: 'linear-gradient(rgba(245,158,11,0.3) 1px, transparent 1px), linear-gradient(90deg, rgba(245,158,11,0.3) 1px, transparent 1px)',
          backgroundSize: '60px 60px'
        }} />
      </div>

      <div className="relative z-10 text-center px-4">
        <AnimatePresence mode="wait">
          {phase === 'typing' && (
            <motion.div
              key={`typing-${lang}`}
              exit={{ opacity: 0, scale: 0.9, filter: 'blur(10px)' }}
              transition={{ duration: 0.5 }}
              className="space-y-2"
            >
              {typewriterTexts.map((line, i) => (
                <motion.p
                  key={i}
                  className="font-mono text-2xl md:text-4xl text-text-secondary"
                  initial={{ opacity: 0 }}
                  animate={{ opacity: i <= lineIndex ? 1 : 0 }}
                >
                  {i < lineIndex ? line : i === lineIndex ? line.slice(0, charIndex) : ''}
                  {i === lineIndex && <span className="animate-pulse text-gold">|</span>}
                </motion.p>
              ))}
            </motion.div>
          )}

          {phase === 'shatter' && (
            <motion.div
              key="shatter"
              className="relative"
              initial={{ opacity: 1 }}
              exit={{ opacity: 0 }}
            >
              <motion.p
                className="font-mono text-4xl md:text-6xl text-neon-red"
                animate={{
                  opacity: [1, 0.5, 0.2, 0],
                  scale: [1, 1.1, 1.2, 1.5],
                  filter: ['blur(0px)', 'blur(2px)', 'blur(5px)', 'blur(20px)'],
                }}
                transition={{ duration: 1 }}
              >
                {t('hero.shatter', lang)}
              </motion.p>
            </motion.div>
          )}

          {phase === 'reform' && (
            <motion.div
              key={`reform-${lang}`}
              initial={{ opacity: 0, scale: 0.5, filter: 'blur(20px)' }}
              animate={{ opacity: 1, scale: 1, filter: 'blur(0px)' }}
              transition={{ duration: 0.8, ease: 'easeOut' }}
            >
              <motion.h1
                className="font-mono text-5xl md:text-7xl font-bold text-gold mb-6"
                animate={{ textShadow: ['0 0 20px rgba(245,158,11,0)', '0 0 40px rgba(245,158,11,0.5)', '0 0 20px rgba(245,158,11,0.2)'] }}
                transition={{ duration: 2, repeat: Infinity }}
              >
                {t('hero.title', lang)}
              </motion.h1>

              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.5, duration: 0.6 }}
                className="space-y-4"
              >
                <p className="text-xl md:text-2xl text-text-secondary max-w-3xl mx-auto">
                  {lang === 'zh'
                    ? <>{t('hero.desc', lang)}</>
                    : <>Drive Skill transforms Claude from a <span className="text-text-muted line-through">passive assistant</span> into an{' '}<span className="text-gold font-semibold">autonomous hierarchical swarm commander</span>.</>
                  }
                </p>
                <p className="text-lg text-text-muted">
                  {t('hero.stats', lang)}
                </p>
              </motion.div>

              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 1.2 }}
                className="mt-12 flex flex-wrap justify-center gap-4"
              >
                {badges.map((item, i) => (
                  <motion.div
                    key={item.key}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 1.4 + i * 0.1 }}
                    className="px-4 py-2 bg-bg-card border border-border rounded-lg text-sm font-mono"
                  >
                    <span className="mr-2">{item.icon}</span>
                    {t(item.key, lang)}
                  </motion.div>
                ))}
              </motion.div>

              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 2 }}
                className="mt-16"
              >
                <a href="#demo" className="text-text-muted hover:text-gold transition-colors">
                  <motion.div
                    animate={{ y: [0, 8, 0] }}
                    transition={{ duration: 2, repeat: Infinity }}
                    className="text-2xl"
                  >
                    ↓
                  </motion.div>
                </a>
              </motion.div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </>
  )
}

export function HeroSection() {
  const { lang } = useLang()
  return (
    <section id="hero" className="min-h-screen flex flex-col items-center justify-center relative overflow-hidden">
      <HeroAnimation key={lang} lang={lang} />
    </section>
  )
}
