import './index.css'
import { LangProvider } from './contexts/LangContext'
import { AnnotationProvider } from './contexts/AnnotationContext'
import { AnnotationsPanel } from './components/AnnotationsPanel'
import { HeroSection } from './sections/HeroSection'
import { ClaudeMechanicsSection } from './sections/ClaudeMechanicsSection'
import { ExecutionSimulator } from './components/ExecutionSimulator'
import { PrincipleCards } from './components/PrincipleCards'
import { ArchitectureSection } from './sections/ArchitectureSection'
import { ExecutionPhasesSection } from './sections/ExecutionPhasesSection'
import { SwarmTopologySection } from './sections/SwarmTopologySection'
import { WaveExecutionSection } from './sections/WaveExecutionSection'
import { SixStepSection } from './sections/SixStepSection'
import { PhilosophySection } from './sections/PhilosophySection'
import { DeepReasoningSection } from './sections/DeepReasoningSection'
import { DirectionValidationSection } from './sections/DirectionValidationSection'
import { PantheonSection } from './sections/PantheonSection'
import { SelfAdversarialReviewSection } from './sections/SelfAdversarialReviewSection'
import { AntiPatternsSection } from './sections/AntiPatternsSection'
import { SelfDriveSection } from './sections/SelfDriveSection'
import { ArchitectSection } from './sections/ArchitectSection'
import { PromptEngineeringSection } from './sections/PromptEngineeringSection'
import { ReferencesSection } from './sections/ReferencesSection'
import { NavBar } from './components/NavBar'
import { SectionDivider } from './components/SectionDivider'
import { ActIIIWrapper } from './components/ActIIIWrapper'
import { useLang } from './contexts/LangContext'
import { t } from './i18n/translations'

function AppInner() {
  const { lang } = useLang()
  return (
    <div className="min-h-screen bg-bg-dark">
      <NavBar />

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
      <div id="act-ii" className="pt-16">
        <div className="max-w-6xl mx-auto px-4 pb-8">
          <h2 className="text-3xl md:text-4xl font-bold text-text-primary">
            <span className="text-gold font-mono">ACT II:</span>{' '}
            <span className="text-gold/80">理解</span>{' '}
            <span className="text-text-secondary text-2xl md:text-3xl">
              — Five Core Principles
            </span>
          </h2>
        </div>

        <ClaudeMechanicsSection />
        <SectionDivider />

        <div className="max-w-6xl mx-auto px-4 py-8">
          <PrincipleCards />
        </div>

        {/* OriginTrace placeholder */}
        {/* <OriginTrace /> */}
      </div>

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
      <div id="act-iv" className="pt-16 pb-24">
        <div className="max-w-6xl mx-auto px-4">
          <h2 className="text-3xl md:text-4xl font-bold text-text-primary mb-12">
            <span className="text-gold font-mono">ACT IV:</span>{' '}
            <span className="text-gold/80">创造</span>{' '}
            <span className="text-text-secondary text-2xl md:text-3xl">
              — Build Your Own Skill
            </span>
          </h2>

          {/* SkillBuilder placeholder */}
          {/* <SkillBuilder /> */}
          <div className="bg-bg-card border border-border rounded-xl p-8 mb-6 text-center">
            <p className="text-text-muted font-mono text-sm">SkillBuilder — coming soon</p>
          </div>

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
      </div>

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
