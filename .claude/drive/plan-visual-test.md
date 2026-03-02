# Plan: visual-test Skill Implementation

## Mission Contract Reference
- Install @playwright/test as devDep
- Install Chromium browser
- Create ~/.claude/commands/visual-test.md with see-fix-verify protocol
- Update idea stage to Plant
- Verify working screenshot

## Approach Overview
The spec (visual-test-skill.md) is already at Plant-level quality. Implementation is
mechanical: install the dependency, write the skill file following the spec's design,
update the idea file status. The skill itself is a markdown prompt (not code) — the
"logic" is instruction to Claude about which Bash commands to run.

The skill needs to be project-agnostic (works in any React/Vite project or any web app)
while defaulting to the erwin dev server at localhost:5173.

## Architecture / Design

```
User invokes /visual-test [optional-args]
    ↓
Claude follows skill protocol:
  1. curl -s http://localhost:5173 → check server running
     If not: cd {project-dir} && npm run dev & → wait 3s
  2. npx playwright screenshot URL /tmp/vt-{timestamp}.png [--viewport/--full-page]
  3. Read /tmp/vt-{timestamp}.png → analyze visually
  4. If issues found: make code changes
  5. sleep 1.5 (HMR reload)
  6. npx playwright screenshot URL /tmp/vt-{timestamp}-after.png
  7. Read again → verify resolved
  8. Report: issues found + fixes + final state
```

## Key Design Decisions
- **Option A (CLI) not B (MCP)**: No server to configure, npx works anywhere — simpler setup wins
- **Project-agnostic defaults**: URL = http://localhost:5173, but instruction tells Claude to adjust for project
- **Timestamp filenames**: /tmp/vt-{ts}.png prevents stale cached reads
- **Max 3 cycles**: prevent infinite loops if issue can't be visually fixed

## Implementation Sequence
1. Install @playwright/test in drive-showcase (establishes npx playwright)
2. Install chromium browser (playwright install chromium)
3. Write ~/.claude/commands/visual-test.md skill file
4. Update _project/ideas/visual-test-skill.md → Plant stage
5. Test: take one screenshot, verify Claude can read it

## Risk Mitigation
- Chromium install size: ~170MB, may take 1-2 min — expected, not a blocker
- Server not running: skill includes detection + startup logic
- Screenshot path: /tmp/ is always writable on macOS

## What This Plan Does NOT Cover
- PostToolUse hook for ambient screenshots (out of scope unless requested)
- Playwright MCP configuration (Option B)
- Testing interactive elements (hover, click) — requires MCP path
