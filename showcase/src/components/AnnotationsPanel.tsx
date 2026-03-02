import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useAnnotations, exportAnnotationsMarkdown } from '../contexts/AnnotationContext'
import { useLang } from '../contexts/LangContext'

export function AnnotationsPanel() {
  const { annotations, remove, clearAll } = useAnnotations()
  const { lang } = useLang()
  const [open, setOpen] = useState(false)
  const [copied, setCopied] = useState(false)

  const count = annotations.length

  const handleExport = () => {
    const md = exportAnnotationsMarkdown(annotations)
    navigator.clipboard.writeText(md).then(() => {
      setCopied(true)
      setTimeout(() => setCopied(false), 2500)
    })
  }

  const handleDownload = () => {
    const md = exportAnnotationsMarkdown(annotations)
    const blob = new Blob([md], { type: 'text/markdown' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `erwin-showcase-annotations-${new Date().toISOString().split('T')[0]}.md`
    a.click()
    URL.revokeObjectURL(url)
  }

  const grouped: Record<string, typeof annotations> = {}
  annotations.forEach(a => {
    grouped[a.sectionLabel] = grouped[a.sectionLabel]
      ? [...grouped[a.sectionLabel], a]
      : [a]
  })

  const ui = {
    title:   lang === 'zh' ? '批注列表' : 'Annotations',
    empty:   lang === 'zh' ? 'hover 任意区块，点击 ✎ 添加批注' : 'Hover any section and click ✎ to annotate',
    export:  lang === 'zh' ? '复制 Markdown' : 'Copy Markdown',
    copied:  lang === 'zh' ? '已复制！' : 'Copied!',
    download:lang === 'zh' ? '下载 .md' : 'Download .md',
    clear:   lang === 'zh' ? '清空全部' : 'Clear All',
    selfDriveHint: lang === 'zh'
      ? '复制后运行: /self-drive improve showcase based on annotations'
      : 'After copying run: /self-drive improve showcase based on annotations',
    section: lang === 'zh' ? '区域' : 'Section',
  }

  return (
    <div className="fixed bottom-6 right-6 z-50 flex flex-col items-end gap-3">
      {/* Main panel */}
      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ opacity: 0, y: 16, scale: 0.96 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 16, scale: 0.96 }}
            transition={{ duration: 0.2 }}
            className="w-80 bg-bg-card border border-border rounded-xl shadow-2xl overflow-hidden"
          >
            {/* Header */}
            <div className="px-4 py-3 bg-gold/5 border-b border-border flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="text-gold font-mono font-bold text-sm">✎</span>
                <span className="text-sm font-mono font-bold text-text-primary">{ui.title}</span>
                {count > 0 && (
                  <span className="text-[10px] bg-gold text-bg-dark px-1.5 py-0.5 rounded-full font-mono font-bold">
                    {count}
                  </span>
                )}
              </div>
              <button
                onClick={() => setOpen(false)}
                className="text-text-muted hover:text-text-primary transition-colors text-sm"
              >
                ✕
              </button>
            </div>

            {/* Content */}
            <div className="max-h-80 overflow-y-auto">
              {count === 0 ? (
                <div className="px-4 py-8 text-center">
                  <p className="text-3xl mb-3">✎</p>
                  <p className="text-xs text-text-muted font-mono">{ui.empty}</p>
                </div>
              ) : (
                <div className="p-3 space-y-4">
                  {Object.entries(grouped).map(([sectionLabel, items]) => (
                    <div key={sectionLabel}>
                      <div className="text-[10px] font-mono text-text-muted mb-2 flex items-center gap-1">
                        <span className="text-gold">§</span> {sectionLabel}
                      </div>
                      <div className="space-y-2">
                        {items.map(a => (
                          <div
                            key={a.id}
                            className="bg-bg-dark rounded-lg p-3 border border-border/50 group"
                          >
                            <div className="flex items-start justify-between gap-2">
                              <div className="flex-1 min-w-0">
                                <p className="text-[10px] text-gold font-mono mb-1 truncate">{a.targetLabel}</p>
                                <p className="text-xs text-text-secondary leading-relaxed">{a.question}</p>
                              </div>
                              <button
                                onClick={() => remove(a.id)}
                                className="flex-shrink-0 text-text-muted hover:text-neon-red text-xs opacity-0 group-hover:opacity-100 transition-opacity"
                              >
                                ✕
                              </button>
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Footer actions */}
            {count > 0 && (
              <div className="px-3 py-3 border-t border-border bg-bg-card space-y-2">
                {/* Self-drive hint */}
                <div className="bg-gold/5 border border-gold/10 rounded p-2">
                  <p className="text-[9px] text-text-muted font-mono leading-relaxed">{ui.selfDriveHint}</p>
                </div>

                {/* Action buttons */}
                <div className="flex gap-2">
                  <button
                    onClick={handleExport}
                    className={`flex-1 py-2 text-[10px] font-mono rounded-lg border transition-all ${
                      copied
                        ? 'bg-neon-green/10 text-neon-green border-neon-green/20'
                        : 'bg-gold/10 text-gold border-gold/20 hover:bg-gold/20'
                    }`}
                  >
                    {copied ? `✓ ${ui.copied}` : ui.export}
                  </button>
                  <button
                    onClick={handleDownload}
                    className="flex-1 py-2 text-[10px] font-mono rounded-lg border border-border/50 text-text-secondary hover:border-border hover:text-text-primary transition-all"
                  >
                    {ui.download}
                  </button>
                </div>
                <button
                  onClick={() => { if (confirm('Clear all annotations?')) clearAll() }}
                  className="w-full py-1.5 text-[9px] font-mono text-text-muted hover:text-neon-red transition-colors"
                >
                  {ui.clear}
                </button>
              </div>
            )}
          </motion.div>
        )}
      </AnimatePresence>

      {/* FAB */}
      <motion.button
        onClick={() => setOpen(v => !v)}
        className={`relative flex items-center gap-2 px-4 py-2.5 rounded-full border font-mono text-sm shadow-xl transition-all ${
          open
            ? 'bg-gold text-bg-dark border-gold'
            : 'bg-bg-card border-border text-text-secondary hover:border-gold/40 hover:text-gold'
        }`}
        whileHover={{ scale: 1.04 }}
        whileTap={{ scale: 0.97 }}
      >
        <span>✎</span>
        <span className="text-xs">{lang === 'zh' ? '批注' : 'Notes'}</span>
        {count > 0 && !open && (
          <motion.span
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            className="absolute -top-1.5 -right-1.5 min-w-[18px] h-[18px] flex items-center justify-center bg-gold text-bg-dark text-[10px] font-bold rounded-full px-1"
          >
            {count}
          </motion.span>
        )}
      </motion.button>
    </div>
  )
}
