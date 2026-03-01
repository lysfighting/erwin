# Skill Idea: visual-test

> Status: INBOX
> Added: 2026-03-02
> Origin: OriginTrace layout bug was invisible to Claude until user shared a screenshot. A visual-test skill would have caught it automatically.

---

## Problem Statement

Claude writes frontend code without seeing the result. This creates a systematic blind spot:
- Layout bugs (like the OriginTrace force graph clustering to edges)
- Z-index/overflow issues
- Color contrast problems
- Responsive breakpoint failures
- Animation timing issues

All of these are **invisible in code** but immediately obvious in a screenshot. Today, the feedback loop is: Claude writes → user sees → user describes bug → Claude fixes. Visual capability collapses this to: Claude writes → Claude sees → Claude fixes.

---

## Proposed Skill: `/visual-test`

### Activation
```
/visual-test [url?]    # Default: http://localhost:5173
/visual-test mobile    # Mobile viewport (390x844)
/visual-test scroll    # Full-page scroll capture
```

### Core Protocol

```
1. DETECT: Is dev server running at target URL?
   → If not: start with `npm run dev &` and wait for ready

2. SCREENSHOT: Take screenshot using best available tool:
   Priority order:
   a. npx playwright screenshot URL /tmp/vt-{ts}.png --viewport-size=1440,900
   b. npx puppeteer screenshot  (if playwright not available)
   c. screencapture -x /tmp/vt-{ts}.png  (macOS, full screen fallback)

3. ANALYZE: Read /tmp/vt-{ts}.png with Read tool
   Ask: "What visual issues are present?"
   - Layout problems (overflow, clipping, misalignment)
   - Color issues (low contrast, wrong colors)
   - Empty/missing content
   - Animation artifacts (captured at wrong frame)

4. FIX: Make code changes for any identified issues

5. VERIFY: Wait 1.5s for HMR, retake screenshot, compare
   → If clean: report success
   → If issues remain: retry (max 3 cycles)

6. REPORT: Summary of issues found + fixes applied
```

---

## Implementation Options

### Option A: Playwright CLI (Recommended — simplest)

**Requirements**: `npx playwright` (no global install needed if project has `@playwright/test`)

```bash
# Screenshot command
npx playwright screenshot http://localhost:5173 /tmp/vt.png

# With specific viewport
npx playwright screenshot http://localhost:5173 /tmp/vt.png --viewport-size=1440,900

# Full page
npx playwright screenshot http://localhost:5173 /tmp/vt.png --full-page
```

**Read command**:
```
Read tool: /tmp/vt.png  ← Claude reads PNG as image, sees the page
```

Claude Code's `Read` tool supports images (PNG, JPG). This is the key enabler — no special MCP needed.

**Setup for erwin**: `cd drive-showcase && npm install --save-dev @playwright/test && npx playwright install chromium`

---

### Option B: Playwright MCP Server (Full browser control)

Install: `npm install -g @playwright/mcp`

Configure in `~/.claude/settings.json`:
```json
{
  "mcpServers": {
    "playwright": {
      "command": "npx",
      "args": ["@playwright/mcp"]
    }
  }
}
```

Provides tools: `browser_navigate`, `browser_screenshot`, `browser_click`, `browser_fill`, `browser_wait`

**Use case**: When the skill needs to interact (click buttons, scroll to sections, test hover states). Required for testing the OriginTrace detail panel (needs a click to open).

---

### Option C: PostToolUse Hook (Ambient, Always-On)

Add to `.claude/settings.json`:
```json
{
  "hooks": {
    "PostToolUse": [{
      "matcher": "Write|Edit",
      "hooks": [{
        "type": "command",
        "command": "sleep 2 && npx playwright screenshot http://localhost:5173 /tmp/last-edit.png 2>/dev/null || true"
      }]
    }]
  }
}
```

With this hook, every time Claude edits a file, a screenshot is automatically taken after 2 seconds (for HMR to reload). Claude can then `Read /tmp/last-edit.png` to see the result without being asked.

**Tradeoff**: Always-on means always taking screenshots even when not needed. Can slow down non-visual tasks.

---

## The "see-fix-verify" Loop

```
Write OriginTrace.tsx
  ↓
npx playwright screenshot http://localhost:5173 /tmp/vt.png
  ↓
Read /tmp/vt.png  →  "nodes are clustered to the right edge"
  ↓
Fix: pin drive-concept nodes at pentagon, reduce repulsion
  ↓
[HMR reload 1.5s]
  ↓
npx playwright screenshot http://localhost:5173 /tmp/vt2.png
  ↓
Read /tmp/vt2.png  →  "pentagon layout centered, clusters correct"
  ↓
✅ Visual verification passed
```

This loop is ~3-4 Bash calls + 3-4 Read calls per visual check. Cheap. Effective.

---

## Skill Design (full)

```markdown
# Visual Test Skill

## Purpose
Enable Claude to see the running frontend page and visually verify correctness.

## Requirements
- Dev server accessible at a URL (default: http://localhost:5173)
- npx playwright available (or @playwright/test in project deps)

## Steps
1. Check dev server is running: curl -s -o /dev/null http://localhost:{PORT}
   - If not: start it (cd {project_dir} && npm run dev &) and wait 3s
2. Take screenshot: npx playwright screenshot http://localhost:{PORT} /tmp/vt.png
3. Analyze: Read /tmp/vt.png and identify visual issues
4. Fix issues if found (code changes)
5. Wait for HMR: sleep 1.5
6. Retake: npx playwright screenshot http://localhost:{PORT} /tmp/vt2.png
7. Verify: Read /tmp/vt2.png, confirm issues resolved
8. Report: list issues found + fixes applied + final state

## Integration with Drive Mode
- Can be called as a task within Drive Mode ("T_visual: Visual verify after implementation")
- Should run after every component implementation wave in UI-heavy missions
- Part of the "Final Verification" phase before MISSION_COMPLETE
```

---

## Would This Have Caught the OriginTrace Bug?

Yes. After writing and building OriginTrace.tsx, the skill would have:
1. Taken a screenshot of http://localhost:5173
2. Scrolled to ACT II (OriginTrace section)
   OR used full-page capture to see all sections
3. Read the screenshot → "nodes clustered to top-left and right edges, graph not centered"
4. Identified root cause from visual evidence + code
5. Applied the pentagon-pinning fix
6. Verified the fix before reporting completion

**Estimated value**: Catches ~60-70% of frontend visual bugs that are invisible in code. The remaining 30-40% (animation timing, hover states, mobile-specific issues) require interactive testing (Playwright MCP Path B).

---

## Next Steps to Implement

1. Add `@playwright/test` to `drive-showcase/devDependencies`
2. Run `npx playwright install chromium`
3. Create skill file: `~/.claude/commands/visual-test.md`
4. Test the see-fix-verify loop on a real component
5. Optionally: add PostToolUse hook for ambient screenshots

**Estimated implementation**: 30 minutes (install + skill file + first test run)
