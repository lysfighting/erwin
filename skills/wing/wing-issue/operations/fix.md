# Operation Fix — Bring Untracked Issue into Wing

> Triggered by `/wing-issue fix {issue}`. Adds wing labels to an existing GitHub Issue so it appears on `/wing` dashboard. Optionally reformats the body into MC structure.

## A1: Fetch current Issue

```bash
gh issue view {issue} --repo "$REPO" --json number,title,body,labels,assignees,milestone,state
```

- If not found → "**ERROR:** Issue **#{issue}** not found." → **STOP**
- If closed → "**ERROR:** Issue **#{issue}** is closed. Reopen it first." → **STOP**
  💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

Display current Issue:

🔍 **ADOPT CANDIDATE** ── **#{issue}** {title} ─────
  **Assignee**   @{assignee or "unassigned"}
  **Labels**     `{current labels}` or "none"
  **Milestone**  `{milestone}` or "—"
────────────────────────────────────────────
{current body (first 20 lines)}
────────────────────────────────────────────

## A2: Check existing labels

Check which wing labels are already present:
- Has `$MISSION_LABEL`? Has `${STATUS_PREFIX}*`? Has `${PRIORITY_PREFIX}*`? Has `size:*`?

If ALL already present → "**ERROR:** Issue **#{issue}** is already in wing." → **STOP**

## A3: Determine missing metadata

AI analyzes the existing title and body to infer:
- **Priority** — P0/P1/P2/P3 (default P2)
- **Size** — S/M/L/XL (default M)
- **Status** — default `wip` if assigned, `queued` if unassigned

## A4: Preview

🔍 **ADOPT PREVIEW** ── **#{issue}** {title} ───────
  **Assignee**  @{assignee or "unassigned"}

  🏷️ **LABELS TO ADD**
    + `{MISSION_LABEL}`
    + `{STATUS_PREFIX}{status}`
    + `{PRIORITY_PREFIX}{priority}`
    + `size:{size}`

  {If milestone not set and CURRENT_VERSION exists:}
  **Milestone** `{CURRENT_VERSION}` (from config)
  **Reformat**  {yes — MC structure | no — as-is}
────────────────────────────────────────────

Use `AskUserQuestion`: "Fix as-is" / "Fix + reformat body" / "Edit priority/size" / "Cancel"

## A5: Apply labels

```bash
gh issue edit {issue} --repo "$REPO" \
  --add-label "$MISSION_LABEL" \
  --add-label "${STATUS_PREFIX}{status}" \
  --add-label "${PRIORITY_PREFIX}{priority}" \
  --add-label "size:{size}"
```

If milestone should be set:
```bash
if [ -n "$CURRENT_VERSION" ]; then
  MILESTONE=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-git.sh" milestone-resolve "$CURRENT_VERSION" 2>/dev/null)
  gh issue edit {issue} --repo "$REPO" --milestone "${MILESTONE:-$CURRENT_VERSION}"
fi
```

## A6: Reformat body (if selected)

If user chose "Fix + reformat body":
1. AI restructures existing content into MC format (Objective, Success Criteria, Sub-tasks, etc.)
2. Preserve all original information

```bash
gh issue edit {issue} --repo "$REPO" --body "{reformatted body}"
```

## A7: Output

✅ **ADOPTED** ── **#{issue}** {title} ──────────────
  **Assignee**   @{assignee or "unassigned"}
  **Labels**     `{all labels}`
  **Milestone**  `{milestone}` or "—"

Use `AskUserQuestion`: "Claim 这个 Issue" / "看 dashboard"
