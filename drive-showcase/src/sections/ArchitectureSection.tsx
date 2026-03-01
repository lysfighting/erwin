import { motion } from 'framer-motion'
import { SectionWrapper } from '../components/SectionWrapper'
import { useLang } from '../contexts/LangContext'
import { t } from '../i18n/translations'

const modules = [
  {
    name: 'commands/',
    files: ['drive.md (112KB)', 'self-drive.md (15KB)', 'architect.md (10KB)'],
    color: 'gold',
    desc: 'Core skill definitions — the soul of the system',
    icon: '⚡',
    why: 'The entire Drive intelligence lives here as .md files. Why files not inline prompts? Because files survive sessions, can be version-controlled, updated without ending the session, and Claude Code natively converts them to slash commands via YAML frontmatter.',
  },
  {
    name: 'hooks/',
    files: ['notify.sh', 'speak.sh (legacy)'],
    color: 'neon-blue',
    desc: 'Event-driven notification system',
    icon: '🔔',
    why: 'Claude Code fires lifecycle events (Stop, Notification, TeammateIdle) as shell hooks. notify.sh turns those events into macOS banners + sounds. Without hooks, autonomous execution is silent — you never know when it needs attention.',
  },
  {
    name: 'messaging/',
    files: ['slack-send.py', 'slack-listen.py', 'slack-poll.py', 'telegram-send.py', 'telegram-poll.py'],
    color: 'neon-purple',
    desc: 'Async remote control via Slack/Telegram',
    icon: '📱',
    why: 'Self-drive needs to reach you wherever you are. Slack/Telegram messaging lets the swarm send proposals and receive your choice — without you sitting at the terminal. Async remote control.',
  },
  {
    name: 'pantheon/',
    files: ['1000+ historical figures', '7 era files (~400KB)'],
    color: 'neon-cyan',
    desc: 'Wisdom of humanity\'s greatest minds',
    icon: '🏛️',
    why: '400KB of structured wisdom from 1000+ historical figures, loaded on-demand by Read tool. Why separate files? Because loading everything upfront wastes context window. Load the relevant era file only when facing that type of decision.',
  },
  {
    name: 'sounds/',
    files: ['5 pre-generated MP3 effects'],
    color: 'neon-green',
    desc: 'Audio feedback for agent lifecycle events',
    icon: '🔊',
    why: 'Audio feedback for agent lifecycle events. Why sounds? Because you walk away during autonomous runs. Sound tells you "it needs attention" or "it\'s done" without polling the terminal.',
  },
  {
    name: 'install.sh',
    files: ['484 lines, one-click setup'],
    color: 'neon-red',
    desc: 'Symlink-based installation — edit repo, changes take effect instantly',
    icon: '🔧',
    why: 'One-click setup that symlinks everything from the repo to ~/.claude/. The symlink model means: edit in repo → change takes effect immediately in all Claude Code sessions. No redeployment. The entire system is a living document.',
  },
]

const colorMap: Record<string, string> = {
  gold: 'border-gold/30 bg-gold/5',
  'neon-blue': 'border-neon-blue/30 bg-neon-blue/5',
  'neon-purple': 'border-neon-purple/30 bg-neon-purple/5',
  'neon-cyan': 'border-neon-cyan/30 bg-neon-cyan/5',
  'neon-green': 'border-neon-green/30 bg-neon-green/5',
  'neon-red': 'border-neon-red/30 bg-neon-red/5',
}

const textColorMap: Record<string, string> = {
  gold: 'text-gold',
  'neon-blue': 'text-neon-blue',
  'neon-purple': 'text-neon-purple',
  'neon-cyan': 'text-neon-cyan',
  'neon-green': 'text-neon-green',
  'neon-red': 'text-neon-red',
}

export function ArchitectureSection() {
  const { lang } = useLang()
  return (
    <SectionWrapper
      id="architecture"
      title={t('arch.title', lang)}
      subtitle={t('arch.subtitle', lang)}
    >
      {/* Key stats */}
      <motion.div
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
        className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-16"
      >
        {[
          { value: '112KB', label: 'Core Prompt', sub: '1,980 lines' },
          { value: '400KB', label: 'Pantheon', sub: '1000+ minds' },
          { value: '31', label: 'Total Files', sub: 'Symlink-installed' },
          { value: 'v3.0', label: 'Version', sub: 'MIT License' },
        ].map((stat, i) => (
          <motion.div
            key={stat.label}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: i * 0.1 }}
            className="bg-bg-card border border-border rounded-lg p-6 text-center"
          >
            <div className="text-3xl font-bold text-gold font-mono">{stat.value}</div>
            <div className="text-sm text-text-primary mt-1">{stat.label}</div>
            <div className="text-xs text-text-muted mt-0.5">{stat.sub}</div>
          </motion.div>
        ))}
      </motion.div>

      {/* Module grid */}
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
        {modules.map((mod, i) => (
          <motion.div
            key={mod.name}
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: i * 0.08 }}
            className={`border rounded-lg p-5 ${colorMap[mod.color]} hover:bg-bg-card-hover transition-colors`}
          >
            <div className="flex items-center gap-2 mb-3">
              <span className="text-xl">{mod.icon}</span>
              <span className={`font-mono font-bold ${textColorMap[mod.color]}`}>{mod.name}</span>
            </div>
            <p className="text-sm text-text-secondary mb-3">{mod.desc}</p>
            <div className="space-y-1">
              {mod.files.map((f) => (
                <div key={f} className="text-xs font-mono text-text-muted pl-4 border-l border-border">
                  {f}
                </div>
              ))}
            </div>
            <div className="mt-3 pt-3 border-t border-border/30">
              <p className="text-[10px] font-mono text-text-muted mb-1">WHY THIS EXISTS:</p>
              <p className="text-xs text-text-muted leading-relaxed">{mod.why}</p>
            </div>
          </motion.div>
        ))}
      </div>

      {/* Three Design Principles */}
      <motion.div
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
        className="mt-12 grid md:grid-cols-3 gap-4"
      >
        {[
          {
            icon: '🔗',
            title: 'Symlink Model',
            color: '#f59e0b',
            principle: 'Edit in repo → instant effect everywhere',
            why: 'All files are symlinked from the git repo to ~/.claude/. No build step, no redeployment, no restart. Change a line in drive.md and every future /drive invocation picks it up. The system is a living document you can iterate on in real time.',
          },
          {
            icon: '📂',
            title: 'Files as State',
            color: '#3b82f6',
            principle: 'Prompts in .md files, not in code',
            why: 'The 112KB drive prompt is just a text file. You can git blame it, diff versions, roll back changes, and send a pull request. The AI\'s "brain" has a full commit history. This makes iteration safe and auditable in a way inline prompts never could be.',
          },
          {
            icon: '🎯',
            title: 'On-Demand Context',
            color: '#a855f7',
            principle: 'Load context when needed, not upfront',
            why: 'The Pantheon\'s 400KB of wisdom is never loaded at session start. It\'s Read() when needed — at the moment of an architectural decision. This keeps the context window clean for actual work while still giving access to deep knowledge.',
          },
        ].map((item, i) => (
          <motion.div
            key={item.title}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: i * 0.1 }}
            className="bg-bg-card border border-border rounded-lg p-5"
            style={{ borderTopColor: item.color, borderTopWidth: 2 }}
          >
            <div className="flex items-center gap-2 mb-2">
              <span className="text-xl">{item.icon}</span>
              <span className="font-mono font-bold text-sm text-text-primary">{item.title}</span>
            </div>
            <p className="text-xs font-mono text-text-muted mb-3 italic">"{item.principle}"</p>
            <p className="text-xs text-text-secondary leading-relaxed">{item.why}</p>
          </motion.div>
        ))}
      </motion.div>
    </SectionWrapper>
  )
}
