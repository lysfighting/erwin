import { motion, AnimatePresence } from 'framer-motion'
import { useState, useMemo, useCallback } from 'react'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'
import {
  type SkillConfig,
  ALL_TOOLS,
  DEFAULT_PHASES,
  createDefaultConfig,
  generateSkillMarkdown,
} from '../lib/skillTemplate'

const STEP_COUNT = 5

const EXECUTION_MODELS = [
  { value: 'single' as const, labelKey: 'sb.model.single' as const, descKey: 'sb.model.single.desc' as const },
  { value: 'hierarchical' as const, labelKey: 'sb.model.hierarchical' as const, descKey: 'sb.model.hierarchical.desc' as const },
  { value: 'pipeline' as const, labelKey: 'sb.model.pipeline' as const, descKey: 'sb.model.pipeline.desc' as const },
  { value: 'custom' as const, labelKey: 'sb.model.custom' as const, descKey: 'sb.model.custom.desc' as const },
]

// ── Step Indicator ──

function StepIndicator({ current, total }: { current: number; total: number }) {
  const { lang } = useLang()
  const labels = [
    t('sb.step1', lang),
    t('sb.step2', lang),
    t('sb.step3', lang),
    t('sb.step4', lang),
    t('sb.step5', lang),
  ]
  return (
    <div className="flex items-center gap-1 md:gap-2 mb-8 overflow-x-auto pb-2">
      {Array.from({ length: total }, (_, i) => (
        <div key={i} className="flex items-center gap-1 md:gap-2 shrink-0">
          <button
            onClick={undefined} // navigation handled by next/back
            className={`
              flex items-center gap-1.5 px-2 py-1 md:px-3 md:py-1.5 rounded-md font-mono text-xs
              transition-colors cursor-default
              ${i === current
                ? 'bg-gold/20 text-gold border border-gold/40'
                : i < current
                  ? 'bg-neon-green/10 text-neon-green border border-neon-green/30'
                  : 'bg-bg-dark text-text-muted border border-border'}
            `}
          >
            <span className="font-bold">{i + 1}</span>
            <span className="hidden md:inline">{labels[i]}</span>
          </button>
          {i < total - 1 && (
            <div className={`w-4 md:w-8 h-px ${i < current ? 'bg-neon-green/40' : 'bg-border'}`} />
          )}
        </div>
      ))}
    </div>
  )
}

// ── Step 1: Meta ──

function StepMeta({
  config,
  onChange,
}: {
  config: SkillConfig
  onChange: (c: SkillConfig) => void
}) {
  const { lang } = useLang()
  return (
    <div className="space-y-6">
      <div>
        <label className="block font-mono text-sm text-text-secondary mb-2">
          {t('sb.meta.name', lang)}
        </label>
        <input
          type="text"
          value={config.meta.name}
          onChange={e => onChange({ ...config, meta: { ...config.meta, name: e.target.value } })}
          placeholder={lang === 'zh' ? '例如：Code Review Assistant' : 'e.g., Code Review Assistant'}
          className="w-full bg-bg-dark border border-border rounded-lg px-4 py-3 text-text-primary
            font-mono text-sm placeholder:text-text-muted focus:border-gold/50 focus:outline-none transition-colors"
        />
      </div>

      <div>
        <label className="block font-mono text-sm text-text-secondary mb-2">
          {t('sb.meta.trigger', lang)}
        </label>
        <div className="flex items-center gap-2">
          <span className="text-text-muted font-mono text-lg">/</span>
          <input
            type="text"
            value={config.meta.triggerWord}
            onChange={e => onChange({ ...config, meta: { ...config.meta, triggerWord: e.target.value.replace(/\s/g, '-').toLowerCase() } })}
            placeholder={lang === 'zh' ? '例如：review' : 'e.g., review'}
            className="flex-1 bg-bg-dark border border-border rounded-lg px-4 py-3 text-text-primary
              font-mono text-sm placeholder:text-text-muted focus:border-gold/50 focus:outline-none transition-colors"
          />
        </div>
      </div>

      <div>
        <label className="block font-mono text-sm text-text-secondary mb-2">
          {t('sb.meta.desc', lang)}
        </label>
        <textarea
          value={config.meta.description}
          onChange={e => onChange({ ...config, meta: { ...config.meta, description: e.target.value } })}
          placeholder={lang === 'zh' ? '一句话描述你的技能做什么' : 'One sentence describing what your skill does'}
          rows={2}
          className="w-full bg-bg-dark border border-border rounded-lg px-4 py-3 text-text-primary
            font-mono text-sm placeholder:text-text-muted focus:border-gold/50 focus:outline-none transition-colors resize-none"
        />
      </div>
    </div>
  )
}

// ── Step 2: Execution Model ──

function StepModel({
  config,
  onChange,
}: {
  config: SkillConfig
  onChange: (c: SkillConfig) => void
}) {
  const { lang } = useLang()
  return (
    <div className="space-y-3">
      <p className="text-text-secondary text-sm mb-4">{t('sb.model.intro', lang)}</p>
      {EXECUTION_MODELS.map(model => (
        <button
          key={model.value}
          onClick={() => {
            const phases = DEFAULT_PHASES.map(p => ({
              ...p,
              enabled: model.value === 'hierarchical'
                ? true  // hierarchical enables all phases including Team Assembly
                : p.id !== 'phaseT' ? p.enabled : false,
            }))
            onChange({ ...config, executionModel: model.value, phases })
          }}
          className={`
            w-full text-left p-4 rounded-lg border transition-all
            ${config.executionModel === model.value
              ? 'border-gold/50 bg-gold/10'
              : 'border-border bg-bg-dark hover:border-border hover:bg-bg-card-hover'}
          `}
        >
          <div className="flex items-center gap-3">
            <div className={`w-4 h-4 rounded-full border-2 flex items-center justify-center
              ${config.executionModel === model.value ? 'border-gold' : 'border-text-muted'}`}
            >
              {config.executionModel === model.value && (
                <div className="w-2 h-2 rounded-full bg-gold" />
              )}
            </div>
            <div>
              <span className="font-mono text-sm text-text-primary">
                {t(model.labelKey, lang)}
              </span>
              {model.value === 'hierarchical' && (
                <span className="ml-2 text-xs text-gold font-mono">
                  {lang === 'zh' ? '推荐' : 'recommended'}
                </span>
              )}
            </div>
          </div>
          <p className="text-text-secondary text-xs mt-2 ml-7">
            {t(model.descKey, lang)}
          </p>
        </button>
      ))}
    </div>
  )
}

// ── Step 3: Phases ──

function StepPhases({
  config,
  onChange,
}: {
  config: SkillConfig
  onChange: (c: SkillConfig) => void
}) {
  const { lang } = useLang()
  return (
    <div className="space-y-4">
      <p className="text-text-secondary text-sm mb-4">{t('sb.phases.intro', lang)}</p>
      {config.phases.map((phase, i) => (
        <div
          key={phase.id}
          className={`p-4 rounded-lg border transition-colors ${
            phase.enabled ? 'border-gold/30 bg-gold/5' : 'border-border bg-bg-dark'
          }`}
        >
          <div className="flex items-center gap-3 mb-2">
            <button
              onClick={() => {
                const phases = config.phases.map((p, j) =>
                  j === i ? { ...p, enabled: !p.enabled } : p
                )
                onChange({ ...config, phases })
              }}
              className={`w-5 h-5 rounded border-2 flex items-center justify-center transition-colors
                ${phase.enabled ? 'border-gold bg-gold/20' : 'border-text-muted'}`}
            >
              {phase.enabled && <span className="text-gold text-xs">✓</span>}
            </button>
            <span className="font-mono text-sm text-text-primary">{phase.name}</span>
          </div>

          {phase.enabled && (
            <div className="ml-8">
              <input
                type="text"
                value={phase.description}
                onChange={e => {
                  const phases = config.phases.map((p, j) =>
                    j === i ? { ...p, description: e.target.value } : p
                  )
                  onChange({ ...config, phases })
                }}
                className="w-full bg-bg-card border border-border rounded px-3 py-2 text-text-secondary
                  font-mono text-xs placeholder:text-text-muted focus:border-gold/50 focus:outline-none transition-colors"
                placeholder={lang === 'zh' ? '描述这个阶段的目标...' : 'Describe this phase\'s goal...'}
              />
            </div>
          )}
        </div>
      ))}
    </div>
  )
}

// ── Step 4: Rules ──

function StepRules({
  config,
  onChange,
}: {
  config: SkillConfig
  onChange: (c: SkillConfig) => void
}) {
  const { lang } = useLang()

  const rules = [
    {
      key: 'neverStop' as const,
      labelKey: 'sb.rules.neverStop' as const,
      descKey: 'sb.rules.neverStop.desc' as const,
    },
    {
      key: 'broadcastProtocol' as const,
      labelKey: 'sb.rules.broadcast' as const,
      descKey: 'sb.rules.broadcast.desc' as const,
    },
    {
      key: 'directionValidation' as const,
      labelKey: 'sb.rules.direction' as const,
      descKey: 'sb.rules.direction.desc' as const,
    },
  ]

  return (
    <div className="space-y-4">
      <p className="text-text-secondary text-sm mb-4">{t('sb.rules.intro', lang)}</p>
      {rules.map(rule => (
        <div
          key={rule.key}
          className={`p-4 rounded-lg border transition-colors ${
            config.rules[rule.key] ? 'border-gold/30 bg-gold/5' : 'border-border bg-bg-dark'
          }`}
        >
          <button
            onClick={() => onChange({
              ...config,
              rules: { ...config.rules, [rule.key]: !config.rules[rule.key] },
            })}
            className="flex items-center gap-3 w-full text-left"
          >
            <div className={`w-5 h-5 rounded border-2 flex items-center justify-center transition-colors shrink-0
              ${config.rules[rule.key] ? 'border-gold bg-gold/20' : 'border-text-muted'}`}
            >
              {config.rules[rule.key] && <span className="text-gold text-xs">✓</span>}
            </div>
            <div>
              <span className="font-mono text-sm text-text-primary">
                {t(rule.labelKey, lang)}
              </span>
              <p className="text-text-secondary text-xs mt-1">
                {t(rule.descKey, lang)}
              </p>
            </div>
          </button>
        </div>
      ))}
    </div>
  )
}

// ── Step 5: Tools ──

function StepTools({
  config,
  onChange,
}: {
  config: SkillConfig
  onChange: (c: SkillConfig) => void
}) {
  const { lang } = useLang()

  const toggleTool = (tool: string) => {
    const tools = config.tools.includes(tool)
      ? config.tools.filter(t => t !== tool)
      : [...config.tools, tool]
    onChange({ ...config, tools })
  }

  const selectAll = () => onChange({ ...config, tools: [...ALL_TOOLS] })
  const selectNone = () => onChange({ ...config, tools: [] })

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between mb-2">
        <p className="text-text-secondary text-sm">{t('sb.tools.intro', lang)}</p>
        <div className="flex gap-2">
          <button onClick={selectAll} className="text-xs font-mono text-gold hover:text-gold-light transition-colors">
            {lang === 'zh' ? '全选' : 'All'}
          </button>
          <span className="text-text-muted text-xs">/</span>
          <button onClick={selectNone} className="text-xs font-mono text-text-muted hover:text-text-secondary transition-colors">
            {lang === 'zh' ? '清空' : 'None'}
          </button>
        </div>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
        {ALL_TOOLS.map(tool => (
          <button
            key={tool}
            onClick={() => toggleTool(tool)}
            className={`
              px-3 py-2.5 rounded-lg border font-mono text-sm text-left transition-all
              ${config.tools.includes(tool)
                ? 'border-gold/40 bg-gold/10 text-gold'
                : 'border-border bg-bg-dark text-text-muted hover:text-text-secondary hover:border-border'}
            `}
          >
            <span className="mr-2">{config.tools.includes(tool) ? '✓' : '○'}</span>
            {tool}
          </button>
        ))}
      </div>
    </div>
  )
}

// ── Preview Panel ──

function PreviewPanel({
  markdown,
  expanded,
  onToggle,
}: {
  markdown: string
  expanded: boolean
  onToggle: () => void
}) {
  const { lang } = useLang()
  return (
    <div className="mt-6 border border-border rounded-lg overflow-hidden">
      <button
        onClick={onToggle}
        className="w-full flex items-center justify-between px-4 py-3 bg-bg-dark hover:bg-bg-card-hover transition-colors"
      >
        <span className="font-mono text-sm text-text-secondary">
          {t('sb.preview.title', lang)}
        </span>
        <motion.span
          animate={{ rotate: expanded ? 180 : 0 }}
          transition={{ duration: 0.2 }}
          className="text-text-muted"
        >
          ▾
        </motion.span>
      </button>

      <AnimatePresence>
        {expanded && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.25 }}
            className="overflow-hidden"
          >
            <pre className="p-4 bg-bg-dark text-text-secondary font-mono text-xs leading-relaxed
              overflow-x-auto max-h-96 overflow-y-auto border-t border-border whitespace-pre-wrap break-words">
              {markdown}
            </pre>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}

// ── Main Component ──

export function SkillBuilder() {
  const { lang } = useLang()
  const [step, setStep] = useState(0)
  const [config, setConfig] = useState<SkillConfig>(createDefaultConfig)
  const [previewExpanded, setPreviewExpanded] = useState(false)

  const markdown = useMemo(() => generateSkillMarkdown(config), [config])

  const handleDownload = useCallback(() => {
    const filename = config.meta.triggerWord
      ? `${config.meta.triggerWord}.md`
      : 'my-skill.md'
    const blob = new Blob([markdown], { type: 'text/markdown;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)
  }, [markdown, config.meta.triggerWord])

  const canProceed = step < STEP_COUNT - 1
  const canGoBack = step > 0
  const isLastStep = step === STEP_COUNT - 1

  const stepContent = (() => {
    switch (step) {
      case 0: return <StepMeta config={config} onChange={setConfig} />
      case 1: return <StepModel config={config} onChange={setConfig} />
      case 2: return <StepPhases config={config} onChange={setConfig} />
      case 3: return <StepRules config={config} onChange={setConfig} />
      case 4: return <StepTools config={config} onChange={setConfig} />
      default: return null
    }
  })()

  return (
    <div className="bg-bg-card border border-border rounded-xl p-6 md:p-8">
      {/* Title */}
      <div className="mb-6">
        <h3 className="font-mono text-xl text-text-primary mb-1">
          {t('sb.title', lang)}
        </h3>
        <p className="text-text-secondary text-sm">
          {t('sb.subtitle', lang)}
        </p>
      </div>

      {/* Stepper */}
      <StepIndicator current={step} total={STEP_COUNT} />

      {/* Step Content */}
      <AnimatePresence mode="wait">
        <motion.div
          key={step}
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          exit={{ opacity: 0, x: -20 }}
          transition={{ duration: 0.2 }}
          className="min-h-[280px]"
        >
          {stepContent}
        </motion.div>
      </AnimatePresence>

      {/* Navigation */}
      <div className="flex items-center justify-between mt-8 pt-6 border-t border-border">
        <button
          onClick={() => canGoBack && setStep(s => s - 1)}
          disabled={!canGoBack}
          className={`px-4 py-2 rounded-lg font-mono text-sm transition-colors
            ${canGoBack
              ? 'text-text-secondary hover:text-text-primary hover:bg-bg-dark'
              : 'text-text-muted cursor-not-allowed'}`}
        >
          {t('sb.back', lang)}
        </button>

        <div className="flex items-center gap-3">
          {isLastStep ? (
            <motion.button
              onClick={handleDownload}
              className="px-6 py-2.5 bg-gold/20 border border-gold/40 rounded-lg font-mono text-sm
                text-gold hover:bg-gold/30 transition-colors"
              whileHover={{ scale: 1.03 }}
              whileTap={{ scale: 0.97 }}
            >
              {t('sb.download', lang)}
            </motion.button>
          ) : (
            <motion.button
              onClick={() => canProceed && setStep(s => s + 1)}
              className="px-6 py-2.5 bg-gold/10 border border-gold/30 rounded-lg font-mono text-sm
                text-gold hover:bg-gold/20 transition-colors"
              whileHover={{ scale: 1.03 }}
              whileTap={{ scale: 0.97 }}
            >
              {t('sb.next', lang)}
            </motion.button>
          )}
        </div>
      </div>

      {/* Preview Panel */}
      <PreviewPanel
        markdown={markdown}
        expanded={previewExpanded}
        onToggle={() => setPreviewExpanded(v => !v)}
      />
    </div>
  )
}
