/**
 * GitHubMock — fake GitHub issue/PR page for tutorial demos.
 * Mimics real GitHub UI with dark theme.
 */
import { motion } from 'framer-motion'

// ─── GitHub Issue Mock ────────────────────────────────

interface GitHubLabel {
  name: string
  color: string  // hex without #
}

interface GitHubIssueMockProps {
  number: number
  title: string
  state: 'open' | 'closed'
  author: string
  assignees?: string[]
  labels: GitHubLabel[]
  milestone?: string
  body: string  // markdown-ish plain text
  comments?: GitHubCommentData[]
}

interface GitHubCommentData {
  author: string
  body: string
  avatar?: string
}

export function GitHubIssueMock({
  number,
  title,
  state,
  author,
  assignees = [],
  labels,
  milestone,
  body,
  comments = [],
}: GitHubIssueMockProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      className="rounded-lg border border-[#30363d] bg-[#0d1117] overflow-hidden font-sans text-sm"
    >
      {/* Top bar */}
      <div className="flex items-center gap-2 px-4 py-2 bg-[#161b22] border-b border-[#30363d]">
        <svg className="w-4 h-4 text-[#8b949e]" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z" />
        </svg>
        <span className="text-[#8b949e]">acme-corp / taskflow</span>
      </div>

      {/* Issue header */}
      <div className="px-4 pt-4 pb-3 border-b border-[#30363d]">
        <div className="flex items-start gap-2">
          <h2 className="text-xl font-semibold text-[#e6edf3] flex-1">
            {title}
            <span className="text-[#8b949e] font-normal ml-2">#{number}</span>
          </h2>
        </div>
        <div className="flex items-center gap-2 mt-2">
          <span className={`px-2 py-0.5 rounded-full text-xs font-medium border ${
            state === 'open'
              ? 'bg-[#238636]/20 text-[#3fb950] border-[#238636]/40'
              : 'bg-[#8957e5]/20 text-[#d2a8ff] border-[#8957e5]/40'
          }`}>
            {state === 'open' ? '● Open' : '● Closed'}
          </span>
          <span className="text-[#8b949e] text-xs">
            {author} opened this issue
          </span>
        </div>
      </div>

      {/* Main content area */}
      <div className="flex flex-col md:flex-row">
        {/* Body */}
        <div className="flex-1 p-4 min-w-0">
          <div className="rounded-md border border-[#30363d] bg-[#161b22]">
            <div className="px-4 py-2 border-b border-[#30363d] flex items-center gap-2">
              <div className="w-5 h-5 rounded-full bg-[#30363d] flex items-center justify-center text-[10px] text-[#8b949e]">
                {author[0].toUpperCase()}
              </div>
              <span className="text-[#e6edf3] text-xs font-semibold">{author}</span>
              <span className="text-[#8b949e] text-xs">commented</span>
            </div>
            <div className="px-4 py-3 text-[#e6edf3] whitespace-pre-wrap text-sm leading-relaxed">
              {body}
            </div>
          </div>

          {/* Comments */}
          {comments.map((comment, i) => (
            <div key={i} className="mt-3 rounded-md border border-[#30363d] bg-[#161b22]">
              <div className="px-4 py-2 border-b border-[#30363d] flex items-center gap-2">
                <div className="w-5 h-5 rounded-full bg-[#30363d] flex items-center justify-center text-[10px] text-[#8b949e]">
                  {comment.author[0].toUpperCase()}
                </div>
                <span className="text-[#e6edf3] text-xs font-semibold">{comment.author}</span>
                <span className="text-[#8b949e] text-xs">commented</span>
              </div>
              <div className="px-4 py-3 text-[#e6edf3] whitespace-pre-wrap text-sm leading-relaxed">
                {comment.body}
              </div>
            </div>
          ))}
        </div>

        {/* Sidebar */}
        <div className="w-full md:w-56 p-4 border-t md:border-t-0 md:border-l border-[#30363d] space-y-4 shrink-0">
          {/* Assignees */}
          <div>
            <div className="text-[#8b949e] text-xs font-semibold mb-1">Assignees</div>
            {assignees.length > 0 ? (
              assignees.map(a => (
                <div key={a} className="flex items-center gap-1.5 text-[#e6edf3] text-xs mt-1">
                  <div className="w-4 h-4 rounded-full bg-[#30363d] flex items-center justify-center text-[8px] text-[#8b949e]">
                    {a[0].toUpperCase()}
                  </div>
                  {a}
                </div>
              ))
            ) : (
              <div className="text-[#484f58] text-xs">No one assigned</div>
            )}
          </div>

          {/* Labels */}
          <div>
            <div className="text-[#8b949e] text-xs font-semibold mb-1">Labels</div>
            <div className="flex flex-wrap gap-1">
              {labels.map(l => (
                <span
                  key={l.name}
                  className="px-2 py-0.5 rounded-full text-[10px] font-medium"
                  style={{
                    backgroundColor: `#${l.color}20`,
                    color: `#${l.color}`,
                    border: `1px solid #${l.color}40`,
                  }}
                >
                  {l.name}
                </span>
              ))}
            </div>
          </div>

          {/* Milestone */}
          {milestone && (
            <div>
              <div className="text-[#8b949e] text-xs font-semibold mb-1">Milestone</div>
              <div className="text-[#e6edf3] text-xs">📦 {milestone}</div>
            </div>
          )}
        </div>
      </div>
    </motion.div>
  )
}

// ─── GitHub PR Mock ────────────────────────────────

interface GitHubPRMockProps {
  number: number
  title: string
  state: 'open' | 'merged' | 'closed'
  author: string
  baseBranch: string
  headBranch: string
  labels: GitHubLabel[]
  reviewers?: string[]
  checks?: { name: string; status: 'pass' | 'fail' | 'pending' }[]
  additions: number
  deletions: number
  files: string[]
  body: string
}

export function GitHubPRMock({
  number,
  title,
  state,
  author,
  baseBranch,
  headBranch,
  labels,
  reviewers = [],
  checks = [],
  additions,
  deletions,
  files,
  body,
}: GitHubPRMockProps) {
  const stateConfig = {
    open: { bg: 'bg-[#238636]/20', text: 'text-[#3fb950]', border: 'border-[#238636]/40', label: '● Open' },
    merged: { bg: 'bg-[#8957e5]/20', text: 'text-[#d2a8ff]', border: 'border-[#8957e5]/40', label: '● Merged' },
    closed: { bg: 'bg-[#da3633]/20', text: 'text-[#f85149]', border: 'border-[#da3633]/40', label: '● Closed' },
  }
  const sc = stateConfig[state]

  return (
    <motion.div
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      className="rounded-lg border border-[#30363d] bg-[#0d1117] overflow-hidden font-sans text-sm"
    >
      {/* Top bar */}
      <div className="flex items-center gap-2 px-4 py-2 bg-[#161b22] border-b border-[#30363d]">
        <svg className="w-4 h-4 text-[#8b949e]" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z" />
        </svg>
        <span className="text-[#8b949e]">acme-corp / taskflow</span>
        <span className="text-[#484f58] mx-1">·</span>
        <span className="text-[#8b949e]">Pull Request</span>
      </div>

      {/* PR header */}
      <div className="px-4 pt-4 pb-3 border-b border-[#30363d]">
        <h2 className="text-xl font-semibold text-[#e6edf3]">
          {title}
          <span className="text-[#8b949e] font-normal ml-2">#{number}</span>
        </h2>
        <div className="flex items-center gap-2 mt-2 flex-wrap">
          <span className={`px-2 py-0.5 rounded-full text-xs font-medium border ${sc.bg} ${sc.text} ${sc.border}`}>
            {sc.label}
          </span>
          <span className="text-[#8b949e] text-xs">
            {author} wants to merge into
            <code className="px-1.5 py-0.5 mx-1 rounded bg-[#388bfd]/15 text-[#58a6ff] text-xs">{baseBranch}</code>
            from
            <code className="px-1.5 py-0.5 mx-1 rounded bg-[#388bfd]/15 text-[#58a6ff] text-xs">{headBranch}</code>
          </span>
        </div>
      </div>

      {/* Content */}
      <div className="flex flex-col md:flex-row">
        <div className="flex-1 p-4 min-w-0 space-y-3">
          {/* PR body */}
          <div className="rounded-md border border-[#30363d] bg-[#161b22] px-4 py-3 text-[#e6edf3] whitespace-pre-wrap leading-relaxed">
            {body}
          </div>

          {/* CI checks */}
          {checks.length > 0 && (
            <div className="rounded-md border border-[#30363d] bg-[#161b22]">
              <div className="px-4 py-2 border-b border-[#30363d] text-[#e6edf3] text-xs font-semibold">
                Checks
              </div>
              {checks.map((c, i) => (
                <div key={i} className="px-4 py-1.5 flex items-center gap-2 text-xs border-b border-[#30363d] last:border-0">
                  <span className={
                    c.status === 'pass' ? 'text-[#3fb950]' :
                    c.status === 'fail' ? 'text-[#f85149]' :
                    'text-[#d29922]'
                  }>
                    {c.status === 'pass' ? '✓' : c.status === 'fail' ? '✗' : '●'}
                  </span>
                  <span className="text-[#e6edf3]">{c.name}</span>
                  <span className={`ml-auto ${
                    c.status === 'pass' ? 'text-[#3fb950]' :
                    c.status === 'fail' ? 'text-[#f85149]' :
                    'text-[#d29922]'
                  }`}>
                    {c.status === 'pass' ? 'Passed' : c.status === 'fail' ? 'Failed' : 'Running...'}
                  </span>
                </div>
              ))}
            </div>
          )}

          {/* Files changed */}
          <div className="rounded-md border border-[#30363d] bg-[#161b22]">
            <div className="px-4 py-2 border-b border-[#30363d] flex items-center gap-3 text-xs">
              <span className="text-[#e6edf3] font-semibold">{files.length} files changed</span>
              <span className="text-[#3fb950]">+{additions}</span>
              <span className="text-[#f85149]">-{deletions}</span>
            </div>
            {files.map((f, i) => (
              <div key={i} className="px-4 py-1 text-[#58a6ff] text-xs border-b border-[#30363d] last:border-0 font-mono">
                {f}
              </div>
            ))}
          </div>
        </div>

        {/* Sidebar */}
        <div className="w-full md:w-48 p-4 border-t md:border-t-0 md:border-l border-[#30363d] space-y-4 shrink-0">
          {/* Reviewers */}
          <div>
            <div className="text-[#8b949e] text-xs font-semibold mb-1">Reviewers</div>
            {reviewers.map(r => (
              <div key={r} className="flex items-center gap-1.5 text-[#e6edf3] text-xs mt-1">
                <div className="w-4 h-4 rounded-full bg-[#30363d] flex items-center justify-center text-[8px]">
                  {r[0].toUpperCase()}
                </div>
                {r}
              </div>
            ))}
          </div>

          {/* Labels */}
          <div>
            <div className="text-[#8b949e] text-xs font-semibold mb-1">Labels</div>
            <div className="flex flex-wrap gap-1">
              {labels.map(l => (
                <span
                  key={l.name}
                  className="px-2 py-0.5 rounded-full text-[10px] font-medium"
                  style={{
                    backgroundColor: `#${l.color}20`,
                    color: `#${l.color}`,
                    border: `1px solid #${l.color}40`,
                  }}
                >
                  {l.name}
                </span>
              ))}
            </div>
          </div>
        </div>
      </div>
    </motion.div>
  )
}
