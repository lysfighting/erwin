# Operation Batch — Create Milestone + MCs

> Triggered by `/wing-issue batch V0.2 -- description\nMC1: ...\nMC2: ...`

## B1: Parse Arguments

Parse `$ARGUMENTS` (after removing `batch` prefix):

1. **Milestone title**: text before `--` separator
2. **Milestone description**: text after `--` on the first line
3. **MC list**: remaining lines starting with `MC1:`, `MC2:`, etc.

Each MC entry:
```
MCn: <description> @assignee
  - [ ] Success criterion 1
  - [ ] Success criterion 2
```

If no MC list → create milestone only.

## B2: Create GitHub Milestone

```bash
MILESTONE_URL=$(gh api "repos/$REPO/milestones" --method POST \
  --field title="$MILESTONE_TITLE" \
  --field description="$MILESTONE_DESC" \
  --jq '.html_url' 2>/dev/null)
```

If milestone already exists → warn but continue (use existing).

```bash
if [ -n "$MILESTONE_TITLE" ]; then
  MILESTONE_TITLE=$(bash "${CLAUDE_SKILL_DIR}/../scripts/wing-git.sh" milestone-resolve "$MILESTONE_TITLE" 2>/dev/null)
fi
```

## B3: Batch MC Creation

For each MC entry:

1. Parse: description, @assignee, success criteria
2. Verify assignee exists in config members
3. AI-analyze description → generate full MC body (same logic as Create Flow Steps 2-4)
4. Collect all generated MCs for preview

## B4: Preview

🏁 **MILESTONE + MC PREVIEW**
────────────────────────────────────────────
  **Milestone**    `{title}`
  **Description**  {desc}
  **MCs to create** `{count}`

  {priority_dot} **MC1:** {title} → @{assignee}
  {priority_dot} **MC2:** {title} → @{assignee}
────────────────────────────────────────────

Use `AskUserQuestion`: "全部发布" / "编辑" / "取消"

- If "编辑" → allow modifications and re-preview
- If Cancel → **STOP**
  💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

## B5: Publish Each MC

```bash
gh issue create \
  --repo "$REPO" \
  --title "{title}" \
  --body "{body}" \
  --label "$MISSION_LABEL" \
  --label "${STATUS_PREFIX}wip" \
  --label "${PRIORITY_PREFIX}{Pn}" \
  --label "domain:{domain}" \
  --label "size:{size}" \
  --assignee "{assignee_github}" \
  ${MILESTONE_TITLE:+--milestone "$MILESTONE_TITLE"}
```

If individual MC creation fails → warn, continue with remaining.

## B6: Notify

```bash
bash "${CLAUDE_SKILL_DIR}/../scripts/wing-notify.sh" milestone.created \
  --title "$MILESTONE_TITLE" --mc_count "$MC_COUNT"
```

## B7: Output

🏁 **MILESTONE** ── {title} ────────────────────
  **URL**  `{milestone_url}`
  **MCs**  `{count}` created

  **#{N1}**  {title1}  @{assignee1}
  **#{N2}**  {title2}  @{assignee2}

Use `AskUserQuestion`: "Claim 第一个 Issue" / "看 dashboard"
