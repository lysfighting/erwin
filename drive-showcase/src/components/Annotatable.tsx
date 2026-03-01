import { useState, type ReactNode, useRef, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useAnnotations } from '../contexts/AnnotationContext'
import { useLang } from '../contexts/LangContext'

interface AnnotatableProps {
  id: string           // unique stable ID — used as localStorage key
  label: string        // human-readable label shown in panel (English)
  labelZh?: string     // Chinese label (optional)
  sectionId: string
  sectionLabel: string
  children: ReactNode
  className?: string
}

export function Annotatable({
  id,
  label,
  labelZh,
  sectionId,
  sectionLabel,
  children,
  className = '',
}: AnnotatableProps) {
  const { add, getByTarget, remove } = useAnnotations()
  const { lang } = useLang()
  const [hovered, setHovered] = useState(false)
  const [open, setOpen] = useState(false)
  const [text, setText] = useState('')
  const inputRef = useRef<HTMLTextAreaElement>(null)
  const existing = getByTarget(id)
  const displayLabel = lang === 'zh' && labelZh ? labelZh : label

  // Close on outside click
  const containerRef = useRef<HTMLDivElement>(null)
  useEffect(() => {
    if (!open) return
    const handler = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false)
        setText('')
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [open])

  useEffect(() => {
    if (open && inputRef.current) {
      inputRef.current.focus()
    }
  }, [open])

  const submit = () => {
    const trimmed = text.trim()
    if (!trimmed) return
    add({ targetId: id, targetLabel: label, sectionId, sectionLabel, question: trimmed })
    setText('')
    setOpen(false)
  }

  const placeholder = lang === 'zh'
    ? '有疑问？输入批注，下次 self-drive 时自动优化...'
    : 'Have a question? Add annotation — self-drive will improve this...'

  const btnTitle = lang === 'zh' ? '添加批注' : 'Annotate'
  const submitLabel = lang === 'zh' ? '保存' : 'Save'
  const cancelLabel = lang === 'zh' ? '取消' : 'Cancel'

  return (
    <div
      ref={containerRef}
      className={`relative group ${className}`}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => { if (!open) setHovered(false) }}
    >
      {children}

      {/* Annotation badge — existing annotations count */}
      {existing.length > 0 && !open && (
        <button
          onClick={(e) => { e.stopPropagation(); setOpen(true) }}
          className="absolute -top-1 -right-1 z-20 flex items-center gap-1 bg-gold text-bg-dark text-[10px] font-mono font-bold px-1.5 py-0.5 rounded-full shadow-lg hover:scale-110 transition-transform"
          title={`${existing.length} annotation${existing.length > 1 ? 's' : ''}`}
        >
          ✎ {existing.length}
        </button>
      )}

      {/* Hover trigger button — show when no annotation yet */}
      {existing.length === 0 && (hovered || open) && !open && (
        <button
          onClick={(e) => { e.stopPropagation(); setOpen(true) }}
          className="absolute -top-1 -right-1 z-20 w-6 h-6 flex items-center justify-center bg-bg-card border border-gold/30 text-gold/60 text-xs rounded-full shadow hover:bg-gold/10 hover:text-gold hover:border-gold/60 transition-all opacity-0 group-hover:opacity-100"
          title={btnTitle}
        >
          ✎
        </button>
      )}

      {/* Annotation panel — inline */}
      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ opacity: 0, y: -4, scale: 0.97 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -4, scale: 0.97 }}
            transition={{ duration: 0.15 }}
            className="absolute top-full left-0 right-0 mt-1 z-30 bg-bg-card border border-gold/30 rounded-lg shadow-2xl overflow-hidden"
          >
            {/* Header */}
            <div className="px-3 py-2 bg-gold/5 border-b border-gold/10 flex items-center justify-between">
              <span className="text-[10px] font-mono text-gold font-bold">✎ {displayLabel}</span>
              <button onClick={() => { setOpen(false); setText('') }} className="text-text-muted hover:text-text-primary text-xs">✕</button>
            </div>

            {/* Existing annotations */}
            {existing.length > 0 && (
              <div className="px-3 py-2 space-y-2 border-b border-border/50 max-h-40 overflow-y-auto">
                {existing.map(a => (
                  <div key={a.id} className="flex items-start gap-2 group/item">
                    <span className="text-gold text-xs mt-0.5 flex-shrink-0">→</span>
                    <p className="text-xs text-text-secondary flex-1">{a.question}</p>
                    <button
                      onClick={() => remove(a.id)}
                      className="text-text-muted hover:text-neon-red text-[10px] flex-shrink-0 opacity-0 group-hover/item:opacity-100 transition-opacity"
                    >
                      ✕
                    </button>
                  </div>
                ))}
              </div>
            )}

            {/* Input */}
            <div className="p-3">
              <textarea
                ref={inputRef}
                value={text}
                onChange={e => setText(e.target.value)}
                onKeyDown={e => {
                  if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) submit()
                  if (e.key === 'Escape') { setOpen(false); setText('') }
                }}
                placeholder={placeholder}
                rows={2}
                className="w-full bg-bg-dark border border-border/60 rounded px-3 py-2 text-xs text-text-primary placeholder:text-text-muted font-mono resize-none focus:outline-none focus:border-gold/40 transition-colors"
              />
              <div className="flex items-center justify-between mt-2">
                <span className="text-[9px] text-text-muted font-mono">⌘↵ {lang === 'zh' ? '提交' : 'submit'}</span>
                <div className="flex gap-1.5">
                  <button
                    onClick={() => { setOpen(false); setText('') }}
                    className="px-2 py-1 text-[10px] font-mono text-text-muted hover:text-text-primary rounded transition-colors"
                  >
                    {cancelLabel}
                  </button>
                  <button
                    onClick={submit}
                    disabled={!text.trim()}
                    className="px-3 py-1 text-[10px] font-mono bg-gold/10 text-gold border border-gold/20 rounded hover:bg-gold/20 transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
                  >
                    {submitLabel}
                  </button>
                </div>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}
