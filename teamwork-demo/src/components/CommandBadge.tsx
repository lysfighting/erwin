/**
 * CommandBadge — styled pill for /team-* slash commands
 * Gives commands visual distinction from regular text.
 */

interface CommandBadgeProps {
  children: string
  size?: 'sm' | 'md'
}

export default function CommandBadge({ children, size = 'md' }: CommandBadgeProps) {
  const sizeClasses = size === 'sm'
    ? 'text-xs px-1.5 py-0.5'
    : 'text-sm px-2 py-0.5'

  return (
    <code
      className={`${sizeClasses} font-mono font-semibold rounded bg-neon-purple/15 text-neon-purple border border-neon-purple/30 whitespace-nowrap inline-block`}
    >
      {children}
    </code>
  )
}
