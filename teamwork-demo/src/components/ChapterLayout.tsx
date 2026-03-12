import { type ReactNode } from 'react'
import { motion } from 'framer-motion'

interface ChapterLayoutProps {
  title: string
  subtitle?: string
  children: ReactNode
}

export default function ChapterLayout({ title, subtitle, children }: ChapterLayoutProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className="max-w-4xl mx-auto px-6 py-8"
    >
      <h1 className="text-3xl font-bold text-neon-blue mb-2">{title}</h1>
      {subtitle && (
        <p className="text-text-secondary text-lg mb-8">{subtitle}</p>
      )}
      <div className="space-y-6">{children}</div>
    </motion.div>
  )
}
