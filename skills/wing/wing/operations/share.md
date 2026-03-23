# Operation Share — Cross-MC Knowledge Sharing

> Triggered by `/wing share #N`. Post structured knowledge from current mission to a sibling MC's Issue.
> Human gateway: ALWAYS confirm before posting.

## SH0: Prerequisites

```bash
GH_USER=$(gh api user --jq '.login' 2>/dev/null)
WING_DIR=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" detect-dir 2>/dev/null) || { echo "Run /wing init first."; exit 1; }
REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
```

Active Contract must exist (need source MC context). If none → "No active mission to share from. Run `/wing-claim` first." → **STOP**

Read current Contract → extract issue number, title, objective.

## SH1: Determine Target

If `$ARGUMENTS` contains `#N` → use that Issue as target.

If no target specified → scan for sibling MCs:

```bash
# Find other open mission Issues (excluding current)
gh issue list --label "$MISSION_LABEL" --state open --json number,title,assignees --limit 20
```

Filter out current Issue. If none → "No sibling MCs to share with." → **STOP**

Display siblings and AskUserQuestion to select target.

## SH2: Scan for Shareable Knowledge

AI reviews the current mission's context to identify knowledge worth sharing:

- **File overlaps**: files modified by this MC that appear in sibling MC's context
- **Discoveries**: bugs found, interface changes, undocumented behavior
- **Warnings**: "don't do X, it causes Y" — lessons from this mission

If nothing worth sharing → "No relevant knowledge to share with **#{target}**." → **STOP**

## SH3: Compose Knowledge Comment

Format as structured Cross-MC comment:

```markdown
🔗 **Cross-MC** from **#{source}** — {title}
**Context**: {what this MC is doing}
**Finding**: {the knowledge to share}
**Confidence**: {high/medium/low}
**Affected files**: {list if applicable}
```

## SH4: Human Gateway (mandatory)

Display the draft comment. AskUserQuestion:

- "发送" → post to target Issue
- "编辑后发送" → AskUserQuestion with useTextArea
- "取消" → **STOP**

⛔ NEVER post without user confirmation. This is the anti-plague rule.

## SH5: Post Comment

```bash
gh issue comment {target} --repo "$REPO" --body "{comment}"
```

Non-fatal: if fails, warn but don't stop.

## Anti-Plague Rules

1. **Human gateway** — AskUserQuestion before every post (no auto-posting)
2. **No transitive propagation** — do NOT re-share knowledge received from OTHER MCs. Only share first-hand discoveries from THIS mission.
3. **Source attribution** — always include source MC number and confidence level
4. **No speculation** — only share verified findings, not hypotheses

## Output

📡 **SHARED** ── **#{source}** → **#{target}**
  {first line of finding}
────────────────────────────────────────────

AskUserQuestion: "继续 drive" (`/wing-drive`) / "回 dashboard" (`/wing`)

💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.
