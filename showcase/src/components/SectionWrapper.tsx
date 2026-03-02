import { motion } from 'framer-motion'
import { type ReactNode } from 'react'

interface SectionWrapperProps {
  id: string
  title: string
  subtitle?: string
  children: ReactNode
  className?: string
}

export function SectionWrapper({ id, title, subtitle, children, className = '' }: SectionWrapperProps) {
  return (
    <section id={id} className={`min-h-screen py-24 px-4 ${className}`}>
      <div className="max-w-6xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 40 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-100px' }}
          transition={{ duration: 0.6 }}
        >
          <h2 className="text-4xl font-bold mb-2">{title}</h2>
          {subtitle && <p className="text-text-secondary text-lg mb-12">{subtitle}</p>}
        </motion.div>
        {children}
      </div>
    </section>
  )
}
