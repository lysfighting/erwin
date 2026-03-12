import { createContext, useContext, useState, type ReactNode } from 'react'

type Lang = 'en' | 'zh'

interface LangContextType {
  lang: Lang
  toggleLang: () => void
}

const LangContext = createContext<LangContextType>({ lang: 'en', toggleLang: () => {} })

export function LangProvider({ children }: { children: ReactNode }) {
  const [lang, setLang] = useState<Lang>('en')
  const toggleLang = () => setLang(l => (l === 'en' ? 'zh' : 'en'))
  return (
    <LangContext.Provider value={{ lang, toggleLang }}>
      {children}
    </LangContext.Provider>
  )
}

export function useLang() {
  return useContext(LangContext)
}
