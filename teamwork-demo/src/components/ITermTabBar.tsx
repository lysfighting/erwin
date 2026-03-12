/**
 * ITermTabBar — mimics iTerm2's tab bar with colored tabs and branch names.
 * Shows how Teamwork's worktree mode creates separate iTerm tabs per mission.
 */
import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'

export interface ITermTab {
  label: string          // e.g. "mission/042-add-auth-bob (claude)"
  color: string          // tab indicator color: red=main, blue=pre-launch, green=mission, orange=rc
  active?: boolean
  hasNotification?: boolean  // blue dot for "waiting for input"
}

interface ITermTabBarProps {
  tabs: ITermTab[]
  onTabClick?: (index: number) => void
  terminalTitle?: string  // title bar text
  children?: React.ReactNode
  /** Per-tab content. When provided, renders tabContents[activeTab] instead of children. */
  tabContents?: React.ReactNode[]
}

export default function ITermTabBar({ tabs, onTabClick, terminalTitle, children, tabContents }: ITermTabBarProps) {
  const [activeTab, setActiveTab] = useState(tabs.findIndex(t => t.active) ?? 0)

  const handleClick = (i: number) => {
    setActiveTab(i)
    onTabClick?.(i)
  }

  const currentTab = tabs[activeTab] || tabs[0]

  return (
    <div className="rounded-lg border border-[#30363d] overflow-hidden bg-[#0d1117]">
      {/* macOS title bar */}
      <div className="flex items-center justify-center px-4 py-1.5 bg-[#1c1c1e] border-b border-[#30363d] relative">
        {/* Traffic lights */}
        <div className="absolute left-3 flex gap-1.5">
          <div className="w-3 h-3 rounded-full bg-[#ff5f57]" />
          <div className="w-3 h-3 rounded-full bg-[#febc2e]" />
          <div className="w-3 h-3 rounded-full bg-[#28c840]" />
        </div>
        {/* Title */}
        <span className="text-[#8b949e] text-xs font-mono truncate max-w-md">
          {terminalTitle || currentTab.label}
        </span>
      </div>

      {/* Tab bar */}
      <div className="flex bg-[#1c1c1e] border-b border-[#30363d] overflow-x-auto">
        {tabs.map((tab, i) => (
          <motion.button
            key={i}
            onClick={() => handleClick(i)}
            className={`relative flex items-center gap-1.5 px-3 py-1.5 text-xs font-mono whitespace-nowrap border-r border-[#30363d] transition-colors min-w-0 ${
              i === activeTab
                ? 'bg-[#0d1117] text-[#e6edf3]'
                : 'bg-[#1c1c1e] text-[#8b949e] hover:bg-[#161b22]'
            }`}
            whileHover={{ backgroundColor: i === activeTab ? '#0d1117' : '#161b22' }}
          >
            {/* Tab color indicator (left edge) */}
            <div
              className="w-1.5 h-1.5 rounded-full shrink-0"
              style={{ backgroundColor: tab.color }}
            />
            {/* Tab label */}
            <span className="truncate max-w-[200px]">{tab.label}</span>
            {/* Notification dot */}
            {tab.hasNotification && (
              <motion.div
                className="w-2 h-2 rounded-full bg-[#3b82f6] shrink-0"
                animate={{ scale: [1, 1.3, 1] }}
                transition={{ repeat: Infinity, duration: 1.5 }}
              />
            )}
            {/* Cmd+N shortcut */}
            <span className="text-[#484f58] ml-1 shrink-0">⌘{i + 1}</span>
          </motion.button>
        ))}
      </div>

      {/* Terminal content */}
      <AnimatePresence mode="wait">
        <motion.div
          key={activeTab}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.15 }}
        >
          {tabContents ? tabContents[activeTab] : children}
        </motion.div>
      </AnimatePresence>
    </div>
  )
}
