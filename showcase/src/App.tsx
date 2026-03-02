import { lazy, Suspense } from 'react'
import './index.css'
import { LangProvider } from './contexts/LangContext'
import { AnnotationProvider } from './contexts/AnnotationContext'
import { AnnotationsPanel } from './components/AnnotationsPanel'
import { HeroSection } from './sections/HeroSection'
import { ExecutionSimulator } from './components/ExecutionSimulator'
import { NavBar } from './components/NavBar'
import { SectionDivider } from './components/SectionDivider'
import { ActIIIWrapper } from './components/ActIIIWrapper'
import { useLang } from './contexts/LangContext'
import { t } from './i18n/translations'

// Lazy-loaded: ACT II components (below the fold)
const ClaudeMechanicsSection = lazy(() => import('./sections/ClaudeMechanicsSection').then(m => ({ default: m.ClaudeMechanicsSection })))
const PrincipleCards = lazy(() => import('./components/PrincipleCards').then(m => ({ default: m.PrincipleCards })))
const OriginTrace = lazy(() => import('./components/OriginTrace').then(m => ({ default: m.OriginTrace })))

// Lazy-loaded: ACT III sections (only loaded when ActIIIWrapper is expanded)
const ArchitectureSection = lazy(() => import('./sections/ArchitectureSection').then(m => ({ default: m.ArchitectureSection })))
const ExecutionPhasesSection = lazy(() => import('./sections/ExecutionPhasesSection').then(m => ({ default: m.ExecutionPhasesSection })))
const SwarmTopologySection = lazy(() => import('./sections/SwarmTopologySection').then(m => ({ default: m.SwarmTopologySection })))
const WaveExecutionSection = lazy(() => import('./sections/WaveExecutionSection').then(m => ({ default: m.WaveExecutionSection })))
const SixStepSection = lazy(() => import('./sections/SixStepSection').then(m => ({ default: m.SixStepSection })))
const PhilosophySection = lazy(() => import('./sections/PhilosophySection').then(m => ({ default: m.PhilosophySection })))
const DeepReasoningSection = lazy(() => import('./sections/DeepReasoningSection').then(m => ({ default: m.DeepReasoningSection })))
const DirectionValidationSection = lazy(() => import('./sections/DirectionValidationSection').then(m => ({ default: m.DirectionValidationSection })))
const PantheonSection = lazy(() => import('./sections/PantheonSection').then(m => ({ default: m.PantheonSection })))
const SelfAdversarialReviewSection = lazy(() => import('./sections/SelfAdversarialReviewSection').then(m => ({ default: m.SelfAdversarialReviewSection })))
const AntiPatternsSection = lazy(() => import('./sections/AntiPatternsSection').then(m => ({ default: m.AntiPatternsSection })))
const SelfDriveSection = lazy(() => import('./sections/SelfDriveSection').then(m => ({ default: m.SelfDriveSection })))
const ArchitectSection = lazy(() => import('./sections/ArchitectSection').then(m => ({ default: m.ArchitectSection })))
const PromptEngineeringSection = lazy(() => import('./sections/PromptEngineeringSection').then(m => ({ default: m.PromptEngineeringSection })))
const ReferencesSection = lazy(() => import('./sections/ReferencesSection').then(m => ({ default: m.ReferencesSection })))

// Lazy-loaded: ACT IV SkillBuilder (only loaded when scrolled into view)
const SkillBuilder = lazy(() => import('./components/SkillBuilder').then(m => ({ default: m.SkillBuilder })))

function AppInner() {
  const { lang } = useLang()
  return (
    <div className="min-h-screen bg-bg-dark">
      <header>
        <NavBar />
      </header>

      <main>

      {/* ═══════════════════════════════════════════════════════════ */}
      {/* ACT I: 震撼 — The Shock (0-15 seconds)                    */}
      {/* ═══════════════════════════════════════════════════════════ */}
      <HeroSection />
      <SectionDivider />
      <div id="act-i" className="max-w-6xl mx-auto px-4 py-16">
        <ExecutionSimulator />
      </div>

      <SectionDivider />

      {/* ═══════════════════════════════════════════════════════════ */}
      {/* ACT II: 理解 — Understanding (1-10 minutes)                */}
      {/* ═══════════════════════════════════════════════════════════ */}
      <section id="act-ii" aria-label="Understanding — Five Core Principles" className="pt-16">
        <div className="max-w-6xl mx-auto px-4 pb-8">
          <h2 className="text-3xl md:text-4xl font-bold text-text-primary">
            <span className="text-gold font-mono">ACT II:</span>{' '}
            <span className="text-gold/80">理解</span>{' '}
            <span className="text-text-secondary text-2xl md:text-3xl">
              — Five Core Principles
            </span>
          </h2>
        </div>

        <Suspense fallback={
          <div className="max-w-6xl mx-auto px-4 py-16 text-center">
            <div className="text-text-muted font-mono text-sm animate-pulse">Loading...</div>
          </div>
        }>
          <ClaudeMechanicsSection />
          <SectionDivider />

          <div className="max-w-6xl mx-auto px-4 py-8">
            <PrincipleCards />
          </div>

          <SectionDivider />
          <div className="max-w-6xl mx-auto px-4 py-8">
            <OriginTrace />
          </div>
        </Suspense>
      </section>

      <SectionDivider />

      {/* ═══════════════════════════════════════════════════════════ */}
      {/* ACT III: 深度 — Deep Dive (10-30 minutes)                  */}
      {/* ═══════════════════════════════════════════════════════════ */}
      <ActIIIWrapper>
        <ArchitectureSection />
        <SectionDivider />
        <ExecutionPhasesSection />
        <SectionDivider />
        <SwarmTopologySection />
        <SectionDivider />
        <WaveExecutionSection />
        <SectionDivider />
        <SixStepSection />
        <SectionDivider />
        <PhilosophySection />
        <SectionDivider />
        <DeepReasoningSection />
        <SectionDivider />
        <DirectionValidationSection />
        <SectionDivider />
        <PantheonSection />
        <SectionDivider />
        <SelfAdversarialReviewSection />
        <SectionDivider />
        <AntiPatternsSection />
        <SectionDivider />
        <SelfDriveSection />
        <SectionDivider />
        <ArchitectSection />
        <SectionDivider />
        <PromptEngineeringSection />
        <SectionDivider />
        <ReferencesSection />
      </ActIIIWrapper>

      <SectionDivider />

      {/* ═══════════════════════════════════════════════════════════ */}
      {/* ACT IV: 创造 — Create Your Own (30+ minutes)               */}
      {/* ═══════════════════════════════════════════════════════════ */}
      <section id="act-iv" aria-label="Create Your Own Skill" className="pt-16 pb-24">
        <div className="max-w-6xl mx-auto px-4">
          <h2 className="text-3xl md:text-4xl font-bold text-text-primary mb-12">
            <span className="text-gold font-mono">ACT IV:</span>{' '}
            <span className="text-gold/80">创造</span>{' '}
            <span className="text-text-secondary text-2xl md:text-3xl">
              — Build Your Own Skill
            </span>
          </h2>

          <Suspense fallback={
            <div className="bg-bg-card border border-border rounded-xl p-8 text-center">
              <div className="text-text-muted font-mono text-sm animate-pulse">Loading SkillBuilder...</div>
            </div>
          }>
            <SkillBuilder />
          </Suspense>

          {/* DesignPatternLibrary placeholder */}
          {/* <DesignPatternLibrary /> */}
          <div className="bg-bg-card border border-border rounded-xl p-8 mb-6 text-center">
            <p className="text-text-muted font-mono text-sm">DesignPatternLibrary — coming soon</p>
          </div>

          {/* CaseStudies placeholder */}
          {/* <CaseStudies /> */}
          <div className="bg-bg-card border border-border rounded-xl p-8 text-center">
            <p className="text-text-muted font-mono text-sm">CaseStudies — coming soon</p>
          </div>
        </div>
      </section>

      </main>

      <footer className="py-16 text-center text-text-muted border-t border-border">
        <p className="font-mono text-sm">Drive Skill v3.0.0 by KC | MIT License</p>
        <p className="text-xs mt-2">{t('footer.built', lang)}</p>
      </footer>
      <AnnotationsPanel />
    </div>
  )
}

function App() {
  return (
    <LangProvider>
      <AnnotationProvider>
        <AppInner />
      </AnnotationProvider>
    </LangProvider>
  )
}

export default App
