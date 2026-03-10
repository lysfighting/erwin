# Visual Encoding Standard — Teamwork Skills

> Design language for all LLM-rendered output in `/team-*` skills.
> Single source of truth. All skill files MUST follow these rules.

---

## Why This Exists

Claude Code renders markdown through a custom pipeline: `marked lexer → React+Ink → Yoga WASM → chalk → ANSI`. This is NOT a browser — many markdown features render differently or not at all. This standard encodes the constraints we discovered into a repeatable design language.

---

## The Palette — 6 Semantic Roles

| Markdown Syntax | Visual Result | Semantic Role | Examples |
|-----------------|--------------|---------------|----------|
| `**bold**` | White bold | Structure: headers, labels, error prefixes | **DASHBOARD**, **ERROR:** |
| `` `backticks` `` | Purple-blue rgb(177,185,249) | Values: commands, paths, counts, config keys | `/team doctor`, `.teamwork/`, `3/5` |
| `*italic*` | Dim italic (ANSI `\e[3m`) | Secondary: branch names, status, timestamps | *mc/042-fix-camera*, *2 days ago* |
| `**#NNN**` | Light-blue clickable (bold+ANSI34→94) | References: Issue/PR numbers (auto-linked) | **#042**, **#153** |
| `────` (Unicode box-drawing) | Plain line | Dividers: section separators | ──────────── |
| Space indent | Whitespace | Structure: hierarchy, grouping | `  Branch: ...` |

### Color Derivation

Claude Code has **three distinct "blue-family" colors**:

1. **ANSI 34** (dark blue) — `chalk.blue()`, used for auto-linked `#NNN`
2. **ANSI 94** (bright blue) — bold+ANSI34 renders as this; `**#NNN**` gets this
3. **rgb(177,185,249)** (purple-blue) — inline code spans, no underline

These are not interchangeable. The syntax determines the color — you cannot "pick" a color directly.

---

## Rules

### DO use

1. **`**bold**` for section titles** — renders white bold with no extra spacing
2. **`` `backticks` `` for commands, paths, counts** — renders purple-blue, visually pops
3. **`*italic*` for branch names and secondary info** — renders dim, recedes
4. **`**#NNN**` for Issue/PR references** — bold wrapping triggers auto-linker → light-blue clickable (OSC 8 hyperlink)
5. **`────` for dividers** — Unicode box-drawing, consistent width
6. **Space-based indentation** — align values with spaces, not list items

### DO NOT use

| Syntax | Problem | Use Instead |
|--------|---------|-------------|
| `###` headers | Injects **2 blank lines** above — breaks compact layout | `**bold**` |
| `- ` bullet lists | Adds indentation, changes visual rhythm | Space indent |
| Triple-backtick code blocks for output | **Kills ALL inline formatting** (bold, italic, links all become literal) | Formatted markdown with inline syntax |
| `═══` double-line box-drawing | Inconsistent weight, visually heavy | `────` |
| `---` horizontal rule | Renders as **literal text** "---" in Claude Code (not a line) | `────` |
| `~~strikethrough~~` in output | Even single `~text~` triggers it (bug #19251) — avoid near tildes | Rephrase |

### Code blocks — when they ARE appropriate

Code blocks are correct for:
- **Bash scripts** that the LLM must execute (the actual code, not the output)
- **ASCII diagrams** containing `*` characters (would trigger italic otherwise)
- **Tables with box-drawing alignment** (e.g., doctor snapshot tables)
- **File content templates** (Contract markdown, config YAML)

Code blocks are WRONG for:
- Dashboard displays, status banners, help text, error messages, confirmation prompts — anything the user reads as formatted output

---

## Issue Auto-Linking

Claude Code's auto-linker (`sG7()`) converts `#NNN` (3+ digits) into OSC 8 clickable hyperlinks pointing to the current repo.

**Key behaviors:**

| Pattern | Renders As | Clickable? |
|---------|-----------|------------|
| `#42` | Plain text (2 digits, below threshold) | No |
| `#042` | Blue text (3 digits, auto-linked) | Yes |
| `**#042**` | **Light-blue bold text** (bold + auto-link) | Yes |
| `` `#042` `` | Purple-blue text (codespan token) | **No** — auto-linker only runs on `text` tokens |

**Rule:** Always zero-pad to 3 digits (`#042` not `#42`). Always wrap in `**bold**` for maximum visibility. Never wrap in backticks — it kills the hyperlink.

**Caveats:** OSC 8 hyperlinks don't work in tmux or VS Code terminal.

---

## Branch Color Squares

Color-coded squares identify branch types in the Worktree section. Same mapping as iTerm2 tab colors.

| Emoji | Branch Pattern | Meaning |
|-------|---------------|---------|
| 🟥 | `main`, `master` | Production |
| 🟦 | `pre-launch`, `develop` | Development |
| 🟧 | `rc/*`, `release/*` | Staging / RC |
| 🟩 | `mission/*` | Active work |
| ⬜ | other | Unclassified |

**iTerm2 tab integration**: `/team init` installs a shell hook (`_teamwork_tab`) that sets tab title to branch name and tab color to branch type. Requires `CLAUDE_CODE_DISABLE_TERMINAL_TITLE=1` in `~/.claude/settings.json` env to prevent Claude Code from overriding the tab title.

---

## Error Message Format

```
**ERROR:** {description}. {remediation using teamwork commands}.
```

- Error prefix: `**ERROR:**` (bold)
- Commands in remediation: `` `/team doctor` `` (backticks)
- Never suggest raw git/rm/shell commands — always guide to a `/team-*` skill

### Doctor Hint (mandatory on every STOP)

When any `/team-*` skill hits an error and STOPs, append:

```
💡 Something wrong? Run `/team doctor` to diagnose, or `/team doctor fix` to auto-repair.
```

Note the formatting: the sentence is plain text, only the two commands are in backticks.

---

## Template Anatomy

A well-formed output block:

```
**SECTION TITLE** ── **#042** {title} ────────────────
  Status: `done`   Priority: `P1`
  Branch: *mc/042-fix-camera*
  Commits: `3` on *mc/042-fix-camera*
────────────────────────────────────────────
```

Breakdown:
- `**SECTION TITLE**` — bold white header
- `**#042**` — light-blue clickable issue link
- `────` — Unicode divider (not `═══`, not `---`)
- `` `done` ``, `` `P1` ``, `` `3` `` — purple-blue values
- `*mc/042-fix-camera*` — dim italic branch name
- Space indent — no bullets

---

## Validation Checklist

When reviewing or writing skill output templates, verify:

- [ ] No triple-backtick code blocks wrapping user-facing output
- [ ] All section headers use `**bold**` (not `###` or plain CAPS)
- [ ] All commands and paths wrapped in `` `backticks` ``
- [ ] All `#NNN` references use `**#NNN**` (bold, 3+ digits)
- [ ] All branch names wrapped in `*italic*`
- [ ] All counts and values wrapped in `` `backticks` ``
- [ ] Dividers use `────` (not `═══` or `---`)
- [ ] Error messages use `**ERROR:**` prefix
- [ ] Every STOP has the doctor hint with correct formatting
- [ ] No raw git/rm/shell commands suggested as remediation

---

## Technical Reference

### Claude Code Rendering Pipeline

```
LLM output (markdown string)
  → marked lexer (tokenize)
  → custom renderer (token → React+Ink components)
  → chalk (ANSI color codes)
  → Yoga WASM (layout)
  → terminal output (ANSI escape sequences)
```

Key insight: the renderer processes **tokens**, not raw text. A `codespan` token and a `text` token containing the same characters render completely differently. This is why the same `#042` looks purple-blue in backticks but blue+clickable outside them — they hit different code paths.

### ANSI Color Map (dark theme)

| ANSI Code | chalk Method | Visual | Used For |
|-----------|-------------|--------|----------|
| `\x1B[1m` | `.bold()` | White bold | `**bold**` |
| `\x1B[3m` | `.italic()` | Dim italic | `*italic*` |
| `\x1B[34m` | `.blue()` | Dark blue | Auto-linked `#NNN` |
| `\x1B[94m` | `.blueBright()` | Bright blue | Bold + auto-link |
| rgb(177,185,249) | (hardcoded) | Purple-blue | Inline code spans |

### Terminal Compatibility

- **Full support**: iTerm2, Kitty, WezTerm, Ghostty
- **Partial**: VS Code terminal (no OSC 8 hyperlinks)
- **Limited**: tmux (no OSC 8, no true italic in some configs)
