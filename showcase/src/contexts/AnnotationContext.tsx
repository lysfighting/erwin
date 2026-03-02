import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'

export interface Annotation {
  id: string
  targetId: string
  targetLabel: string
  sectionId: string
  sectionLabel: string
  question: string
  createdAt: number
}

interface AnnotationContextValue {
  annotations: Annotation[]
  add: (a: Omit<Annotation, 'id' | 'createdAt'>) => void
  remove: (id: string) => void
  clearAll: () => void
  getByTarget: (targetId: string) => Annotation[]
}

const AnnotationContext = createContext<AnnotationContextValue>({
  annotations: [],
  add: () => {},
  remove: () => {},
  clearAll: () => {},
  getByTarget: () => [],
})

const STORAGE_KEY = 'erwin-showcase-annotations'

export function AnnotationProvider({ children }: { children: ReactNode }) {
  const [annotations, setAnnotations] = useState<Annotation[]>(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY)
      return raw ? JSON.parse(raw) : []
    } catch {
      return []
    }
  })

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(annotations))
  }, [annotations])

  const add = (a: Omit<Annotation, 'id' | 'createdAt'>) => {
    setAnnotations(prev => [
      ...prev,
      { ...a, id: `ann-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`, createdAt: Date.now() },
    ])
  }

  const remove = (id: string) => {
    setAnnotations(prev => prev.filter(a => a.id !== id))
  }

  const clearAll = () => setAnnotations([])

  const getByTarget = (targetId: string) =>
    annotations.filter(a => a.targetId === targetId)

  return (
    <AnnotationContext.Provider value={{ annotations, add, remove, clearAll, getByTarget }}>
      {children}
    </AnnotationContext.Provider>
  )
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAnnotations() {
  return useContext(AnnotationContext)
}

// eslint-disable-next-line react-refresh/only-export-components
export function exportAnnotationsMarkdown(annotations: Annotation[]): string {
  if (annotations.length === 0) {
    return '# Drive Showcase Annotations\n\n_(no annotations yet)_'
  }

  const date = new Date().toISOString().split('T')[0]
  const grouped: Record<string, Annotation[]> = {}
  annotations.forEach(a => {
    const key = `${a.sectionId}|||${a.sectionLabel}`
    grouped[key] = grouped[key] ? [...grouped[key], a] : [a]
  })

  const lines: string[] = [
    `# Drive Showcase Annotations — ${date}`,
    '',
    `> ${annotations.length} question${annotations.length > 1 ? 's' : ''} from reviewer. Use these to improve the Drive Skill showcase.`,
    '',
  ]

  Object.entries(grouped).forEach(([key, items]) => {
    const [, sectionLabel] = key.split('|||')
    lines.push(`## ${sectionLabel}`)
    lines.push('')
    items.forEach((a, i) => {
      lines.push(`### ${i + 1}. ${a.targetLabel}`)
      lines.push(`**Question**: ${a.question}`)
      lines.push(`_Annotated: ${new Date(a.createdAt).toLocaleString()}_`)
      lines.push('')
    })
  })

  lines.push('---')
  lines.push('## How to use with self-drive')
  lines.push('```')
  lines.push('/self-drive improve drive showcase — clarify the annotated sections above')
  lines.push('```')

  return lines.join('\n')
}
