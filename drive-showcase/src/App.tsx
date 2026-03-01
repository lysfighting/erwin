import './index.css'
import { LangProvider } from './contexts/LangContext'
import { AnnotationProvider } from './contexts/AnnotationContext'
import { AnnotationsPanel } from './components/AnnotationsPanel'
import { HeroSection } from './sections/HeroSection'
import { DemoSection } from './sections/DemoSection'
import { ClaudeMechanicsSection } from './sections/ClaudeMechanicsSection'
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
import { useLang } from './contexts/LangContext'
import { t } from './i18n/translations'

function AppInner() {
  const { lang } = useLang()
  return (
    <div className="min-h-screen bg-bg-dark">
      <NavBar />
      <HeroSection />
      <SectionDivider />
      <DemoSection />
      <SectionDivider />
      <ClaudeMechanicsSection />
      <SectionDivider />
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
