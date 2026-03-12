/**
 * ITermNotification — mimics macOS notification center alerts from iTerm2.
 * Shows how Claude Code sends native notifications when waiting for input.
 */
import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'

interface NotificationItem {
  title: string       // e.g. "Alert"
  subtitle: string    // e.g. "Session mission/042-add-auth (claude) #2:"
  body: string        // e.g. "Claude is waiting for your input"
  time: string        // e.g. "11分钟前"
  icon?: 'iterm' | 'claude'
}

interface ITermNotificationProps {
  notifications: NotificationItem[]
  autoPlay?: boolean      // auto-dismiss and show one by one
  autoPlayDelay?: number  // ms between notifications
}

export default function ITermNotification({ notifications, autoPlay = true, autoPlayDelay = 2000 }: ITermNotificationProps) {
  const [visible, setVisible] = useState<number[]>([])

  useEffect(() => {
    if (!autoPlay) {
      setVisible(notifications.map((_, i) => i))
      return
    }

    // Show notifications one by one
    const timers: ReturnType<typeof setTimeout>[] = []
    notifications.forEach((_, i) => {
      timers.push(setTimeout(() => {
        setVisible(prev => [...prev, i])
      }, i * autoPlayDelay))
    })

    return () => timers.forEach(clearTimeout)
  }, [notifications, autoPlay, autoPlayDelay])

  return (
    <div className="flex flex-col gap-2 w-80">
      {/* Header */}
      <div className="flex items-center justify-between px-2">
        <span className="text-[#e6edf3] text-xs font-semibold tracking-wide">
          iTerm
        </span>
        <button className="text-[#8b949e] text-xs hover:text-[#e6edf3] transition-colors">
          Clear
        </button>
      </div>

      <AnimatePresence>
        {visible.map(i => {
          const notif = notifications[i]
          return (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: 60, scale: 0.95 }}
              animate={{ opacity: 1, x: 0, scale: 1 }}
              exit={{ opacity: 0, x: 60, scale: 0.95 }}
              transition={{ type: 'spring', damping: 20, stiffness: 300 }}
              className="rounded-xl bg-[#2c2c2e]/95 backdrop-blur-xl border border-[#3a3a3c] p-3 shadow-2xl"
            >
              <div className="flex gap-3">
                {/* Icon */}
                <div className="w-10 h-10 rounded-lg bg-[#1c1c1e] flex items-center justify-center shrink-0 border border-[#3a3a3c]">
                  {notif.icon === 'claude' ? (
                    <span className="text-lg">🤖</span>
                  ) : (
                    <svg className="w-6 h-6" viewBox="0 0 24 24" fill="none">
                      <rect x="2" y="4" width="20" height="16" rx="3" fill="#1c1c1e" stroke="#22c55e" strokeWidth="1.5" />
                      <text x="12" y="15" textAnchor="middle" fill="#22c55e" fontSize="8" fontFamily="monospace">{'>'}_</text>
                    </svg>
                  )}
                </div>
                {/* Content */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <span className="text-[#e6edf3] text-xs font-semibold">{notif.title}</span>
                    <span className="text-[#8b949e] text-[10px]">{notif.time}</span>
                  </div>
                  <div className="text-[#8b949e] text-[11px] mt-0.5 truncate">{notif.subtitle}</div>
                  <div className="text-[#e6edf3] text-xs mt-0.5">{notif.body}</div>
                </div>
              </div>
            </motion.div>
          )
        })}
      </AnimatePresence>
    </div>
  )
}
