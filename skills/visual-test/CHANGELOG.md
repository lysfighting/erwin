# /visual-test Changelog

## [Unreleased]

---

## [1.0.0] — 2026-03-02

### Added
- Initial skill implementation
- 3 call variants: default (1440×900), mobile (390×844), full-page
- Custom URL support (`/visual-test http://localhost:3000`)
- Step 1: argument parsing (empty / mobile / full / URL)
- Step 2: dev server detection via `curl` + auto-start if not running
- Step 3: screenshot via `npx playwright screenshot` with timestamp filenames
- Step 4: visual analysis using Claude Code's native PNG `Read` tool support
- Step 5: code fixes for identified visual issues
- Step 6: HMR wait (1.5s) + retake screenshot
- Step 7: verify and report (issues found + fixes applied)
- Max 3 fix-verify cycles before reporting unresolved issues
- Troubleshooting guide: blank screenshots, Chromium not installed, server issues

### Key design decision
Playwright CLI (Option A) over Playwright MCP Server (Option B):
- No MCP server to configure
- `npx playwright` works anywhere with `@playwright/test` in devDependencies
- Simpler for the 80% use case (screenshot + analyze, not click/interact)

### Dependencies
- `@playwright/test` (devDependency)
- Chromium browser (`npx playwright install chromium`)

### Known limitations
- No interactive testing (hover, click, scroll-to-section) — requires Playwright MCP for those
- Screenshots at page load; doesn't wait for animations to complete

### Origin
- Idea: `_project/ideas/skills/visual-test-skill.md` (Fruit 🍎 stage)
- Feedback: `.claude/drive/feedback-visual-test.md`
- Session: `_project/sessions/2026-03-02.md`
