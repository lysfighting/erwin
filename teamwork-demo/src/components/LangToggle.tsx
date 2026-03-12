import { useLang } from '../i18n/LangContext'

export default function LangToggle() {
  const { lang, toggleLang } = useLang()
  return (
    <button
      onClick={toggleLang}
      className="fixed top-4 right-4 z-50 px-3 py-1.5 text-sm font-mono rounded border border-border bg-bg-card text-text-secondary hover:bg-bg-card-hover hover:text-text-primary transition-colors"
    >
      {lang === 'en' ? '中文' : 'EN'}
    </button>
  )
}
