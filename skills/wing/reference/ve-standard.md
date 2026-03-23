# Visual Encoding Standard

> Actionable rules for all `/wing-*` skill output. Single source of truth.
> Full technical reference: `docs/visual-encoding-standard.md`

---

## The Palette — 6 Semantic Roles

| Markdown Syntax | Visual Result | Semantic Role | Examples |
|-----------------|--------------|---------------|----------|
| `**bold**` | White bold | Structure: headers, labels, error prefixes | **DASHBOARD**, **ERROR:** |
| `` `backticks` `` | Purple-blue rgb(177,185,249) | Values: commands, paths, counts, config keys | `/wing doctor`, `.wing/`, `3/5` |
| `*italic*` | Dim italic (ANSI `\e[3m`) | Secondary: branch names, status, timestamps | *mission/042-fix-camera*, *2 days ago* |
| `**#NNN**` | Light-blue clickable (bold+ANSI34→94) | References: Issue/PR numbers (auto-linked) | **#042**, **#153** |
| `────` (Unicode box-drawing) | Plain line | Dividers: section separators | ──────────── |
| Space indent | Whitespace | Structure: hierarchy, grouping | `  Branch: ...` |

---

## Rules

### DO use

1. `**bold**` for section titles — renders white bold with no extra spacing
2. `` `backticks` `` for commands, paths, counts — renders purple-blue
3. `*italic*` for branch names and secondary info — renders dim
4. `**#NNN**` for Issue/PR references — bold + 3+ digits → light-blue clickable
5. `────` for dividers — Unicode box-drawing, consistent width
6. Space-based indentation — align values with spaces, not list items

### DO NOT use

| Syntax | Problem | Use Instead |
|--------|---------|-------------|
| `###` headers | Injects 2 blank lines above — breaks compact layout | `**bold**` |
| `- ` bullet lists | Adds indentation, changes visual rhythm | Space indent |
| Triple-backtick code blocks for output | Kills ALL inline formatting (bold, italic, links) | Formatted markdown with inline syntax |
| `═══` double-line box-drawing | Inconsistent weight, visually heavy | `────` |
| `---` horizontal rule | Renders as literal text "---" in Claude Code | `────` |
| `~~strikethrough~~` in output | Even single `~text~` triggers it (bug) | Rephrase |

### Code blocks — when appropriate

Code blocks are correct for:
- Bash scripts the LLM must execute
- ASCII diagrams containing `*` characters
- File content templates (Contract, config YAML)

Code blocks are WRONG for:
- Dashboard displays, status banners, help text, error messages, confirmation prompts

---

## Branch Color Squares

Color-coded squares identify branch types in dashboard/worktree sections:

| Emoji | Branch Pattern | Meaning |
|-------|---------------|---------|
| 🟥 | `main`, `master` | Production |
| 🟦 | `pre-launch`, `develop` | Development |
| 🟧 | `rc/*`, `release/*` | Staging / RC |
| 🟩 | `mission/*` | Active work |
| ⬜ | other | Unclassified |

---

## Issue Auto-Linking

| Pattern | Renders As | Clickable? |
|---------|-----------|------------|
| `#42` | Plain text (2 digits) | No |
| `#042` | Blue text (3 digits, auto-linked) | Yes |
| `**#042**` | Light-blue bold text | Yes |
| `` `#042` `` | Purple-blue text | **No** — kills hyperlink |

**Rule:** Always zero-pad to 3+ digits. Always wrap in `**bold**`. Never wrap in backticks.

---

## Error Format

```
**ERROR:** {description}. {remediation using /wing-* commands}.
```

- Error prefix: `**ERROR:**` (bold)
- Commands in remediation: `` `/wing doctor` `` (backticks)
- Never suggest raw git/rm/shell commands — guide to `/wing-*` skills

### Doctor Hint (mandatory on every STOP)

💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

---

## Tips (every primary output)

After every command completion output, append a random tip:

```
💡 **Tip:** {random line from ${CLAUDE_SKILL_DIR}/scripts/wing-tips.txt — skip # comments}
```

Read the tips file, pick one non-comment line at random. Commands in tips must be in backticks.

---

## Validation Checklist

- [ ] No triple-backtick code blocks wrapping user-facing output
- [ ] All section headers use `**bold**` (not `###` or plain CAPS)
- [ ] All commands and paths wrapped in `` `backticks` ``
- [ ] All `#NNN` references use `**#NNN**` (bold, 3+ digits)
- [ ] All branch names wrapped in `*italic*`
- [ ] All counts and values wrapped in `` `backticks` ``
- [ ] Dividers use `────` (not `═══` or `---`)
- [ ] Error messages use `**ERROR:**` prefix
- [ ] Every STOP has the doctor hint
- [ ] No raw git/rm/shell commands suggested as remediation
