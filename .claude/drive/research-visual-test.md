# Research: visual-test Skill Implementation

## Codebase Understanding
- Project: erwin / drive-showcase (React 19.2, Vite 7.3, port 5173)
- Skills target: `~/.claude/commands/visual-test.md`
- Idea spec: `_project/ideas/visual-test-skill.md` (INBOX stage, but Plant-level content)

## Key Findings
- `@playwright/test` NOT in package.json — install required as devDependency
- `npx playwright screenshot` is the core command (no MCP needed)
- Claude Code Read tool natively supports PNG — key enabler
- Dev server: `npm run dev` → port 5173 (Vite default)
- Spec includes: full see-fix-verify loop, bash commands, 3 variants (default/mobile/scroll)
- Origin story: OriginTrace force-graph nodes clustered to edges — invisible in code, obvious in screenshot

## Implementation Plan (from spec)
1. `cd drive-showcase && npm install --save-dev @playwright/test`
2. `npx playwright install chromium`
3. Create `~/.claude/commands/visual-test.md` with see-fix-verify protocol
4. Update `_project/ideas/visual-test-skill.md` stage: INBOX → Plant

## Unknowns & Risks
- PostToolUse hook: spec says optional — skip unless user requests
- ADR needed? Spec is detailed enough to serve as its own doc, no ADR required for a skill file

## Dependencies & Constraints
- Playwright Chromium install can take 1-2 minutes (large download)
- Skill file must work on any project (not erwin-specific) → parameterize URL/port
