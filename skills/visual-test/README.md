# /visual-test

Take a screenshot of the running frontend, analyze it visually, fix issues found, and verify the fix. Collapses the "write → user sees → describes bug → Claude fixes" loop into "write → Claude sees → Claude fixes."

## What it does

Uses Playwright CLI to screenshot the dev server, then uses Claude Code's native `Read` tool support for PNG images to visually inspect the page. Catches ~60-70% of frontend visual bugs that are invisible in code.

## Usage

```
/visual-test              # Default: localhost:5173, 1440×900
/visual-test mobile       # Mobile viewport: 390×844 (iPhone 14)
/visual-test full         # Full-page capture (all scrollable content)
/visual-test <url>        # Custom URL
```

## Requirements

- Dev server accessible at a URL (default: `http://localhost:5173`)
- `@playwright/test` in project devDependencies
- Chromium browser installed via Playwright

### Setup (one-time per project)

```bash
cd your-project
npm install --save-dev @playwright/test
npx playwright install chromium
```

## The see-fix-verify loop

```
screenshot → analyze → fix → wait for HMR → screenshot → verify
```

- Up to 3 fix-verify cycles before reporting failure
- Captures at 1440×900 by default (or mobile/full-page variant)
- Uses `/tmp/vt-{timestamp}.png` — no stale cache reads

## Install

Via erwin's install.sh:
```bash
./install.sh visual-test
```

Or to install all erwin skills:
```bash
./install.sh
```

## Origin

Inspired by the OriginTrace layout bug (invisible in code, obvious in screenshot).
Implemented in Session F (2026-03-02).
See: `_project/ideas/visual-test-skill.md` (Fruit 🍎 stage)
