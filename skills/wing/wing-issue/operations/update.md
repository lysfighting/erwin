# Operation Update — Edit Existing Issue

> Triggered by `/wing-issue update` **#042** or `/wing-issue update` **#042** `<changes>`.

## U1: Fetch current Issue

```bash
gh issue view {issue} --repo "$REPO" --json number,title,body,labels,state
```

- If not found → "**ERROR:** Issue **#{issue}** not found." → **STOP**
- If closed → "**ERROR:** Issue **#{issue}** is closed. Reopen it first." → **STOP**
  💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

Display current Issue summary:

📋 **CURRENT** ── **#{issue}** {title} ─────────────
  **Labels**  `{labels}`
────────────────────────────────────────────
{current body}
────────────────────────────────────────────

## U2: Determine changes

**If user provided change description** (e.g., `/wing-issue update` **#042** `add constraint`):
- AI applies the described changes to the existing body
- Preserve the mission-contract structure
- Only modify sections affected by the change description

**If no change description** (e.g., `/wing-issue update` **#042**):
- Use `AskUserQuestion`: "Edit title" / "Edit priority/size" / "Edit sub-tasks" / "Edit full body"
- Apply user's selection interactively

## U3: Preview changes

✏️ **UPDATE PREVIEW** ── **#{issue}** ───────────────
  **Title**  {new title, or unchanged}

  📝 **CHANGES**
  {Show what changed — highlight modified sections}
────────────────────────────────────────────

Use `AskUserQuestion`: "Apply changes" / "Edit more" / "Cancel"

- If Cancel → **STOP**
  💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

## U4: Apply update

```bash
gh issue edit {issue} --repo "$REPO" --title "{new title}" --body "{updated body}"
```

If labels changed:
```bash
gh issue edit {issue} --repo "$REPO" \
  --remove-label "${PRIORITY_PREFIX}{old}" --add-label "${PRIORITY_PREFIX}{new}" \
  --remove-label "size:{old}" --add-label "size:{new}"
```

## U5: Output

✏️ **UPDATED** ── **#{issue}** {title} ──────────────
  **URL**  `{issue URL}`
  📝 **CHANGES**
  {summary of what changed}
  Claimed assignees will see freshness prompt on `/wing-drive`.
────────────────────────────────────────────
