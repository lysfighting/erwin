---
name: wing-issue
description: "Create mission Issue from natural language. Try: /wing-issue help"
version: "0.1.0"
---

# /wing-issue — Mission Contract Issue Manager

> You own this analysis. Before generating, think: what is the REAL problem? Is the user's framing the bone (irreducible need) or flesh (one possible approach)? Consider at least one competing hypothesis. Assess blast radius: if this Issue is framed wrong, how much rework results? High blast radius → say so explicitly in the preview. Watch for mediocre consensus and implicit decisions nobody explicitly chose.

> Create, view, comment, update, or batch-create GitHub Issues in mission-contract format.

**Visual Encoding** — apply to ALL output (see `${CLAUDE_SKILL_DIR}/../reference/ve-standard.md`)

**User input**: $ARGUMENTS

**Argument routing:**

| Input | Action |
|-------|--------|
| `<description>` | Create MC (solo: self-assign; team: prompt for assignee) |
| `<description> @assignee` | Create MC and assign to @assignee |
| `#N` (no text) | Show Issue **#N** details |
| `#N <text>` | Add comment to Issue **#N** |
| `update #N` or `update #N <changes>` | AI-assisted update of Issue body |
| `fix #N` | Fix untracked Issue — add wing labels, optionally reformat |
| `batch <milestone> -- ...` | Batch create milestone + MCs |
| `help` or `-h` | Usage guide |

If `$ARGUMENTS` is `help` or `-h`, output the following and **STOP**:

**`/wing-issue` — Mission Contract Issue Manager** (`v0.1.0`)

**USAGE**
  `/wing-issue <description>`                    Create MC (solo: self-assign; team: prompt)
  `/wing-issue <description> @assignee`          Create MC and assign to @assignee
  `/wing-issue #042`                             Show Issue **#042** details
  `/wing-issue #042 <comment text>`              Add comment to Issue **#042**
  `/wing-issue update #042 <changes>`            AI-assisted update of Issue **#042**
  `/wing-issue fix #042`                         Add wing labels to Issue **#042**
  `/wing-issue batch V0.2 -- ...`                Batch create milestone + MCs
  `/wing-issue help`                             Show this guide

**EXAMPLES**
  `/wing-issue resolve camera permission on iOS Safari`
  `/wing-issue add rate limiting to upload API @xxLe`
  `/wing-issue batch V0.2 -- User system and campaign basics`

NEXT: Assignee runs `/wing-claim` to create branch + Mission Contract

---

## Step 0: Prerequisites

```bash
GH_USER=$(gh api user --jq '.login' 2>/dev/null)
if [ -z "$GH_USER" ]; then
  echo "ERROR: Cannot get GitHub user identity. Run 'gh auth login' first."
  exit 1
fi
```

```bash
WING_DIR=$(bash "${CLAUDE_SKILL_DIR}/../scripts/wing-config.sh" detect-dir 2>/dev/null) || {
  echo "No wing config found. Run /wing init first."
  # STOP
}
```

- If no config → "**ERROR:** Wing not initialized. Run `/wing` first." → **STOP**
  💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

Read `$WING_DIR/config.yml` → extract labels, members, roles.

```bash
eval "$(bash "${CLAUDE_SKILL_DIR}/../scripts/wing-config.sh" resolve-labels 2>/dev/null)"
```

```bash
REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null)
BASE_BRANCH=$(bash "${CLAUDE_SKILL_DIR}/../scripts/wing-config.sh" conventions.base_branch "main" 2>/dev/null)
BASE_BRANCH="${BASE_BRANCH:-main}"

# Milestone binding — mandatory when versions.current is configured
CURRENT_VERSION=$(bash "${CLAUDE_SKILL_DIR}/../scripts/wing-config.sh" versions.current "" 2>/dev/null)
```

### Team Mode Detection

```bash
TEAM_MODE=$(bash "${CLAUDE_SKILL_DIR}/../scripts/wing-config.sh" team-mode 2>/dev/null)
```

- If `TEAM_MODE == "solo"` → **SOLO_MODE=true** (auto-assign to self)
- If `TEAM_MODE == "team"` → **SOLO_MODE=false** (prompt for assignee)

---

## Route: Parse Arguments

Parse `$ARGUMENTS`:

1. **If starts with `batch`** → load and follow `${CLAUDE_SKILL_DIR}/operations/batch.md`
2. **If starts with `update`** (followed by `#N`) → load and follow `${CLAUDE_SKILL_DIR}/operations/update.md`
3. **If starts with `fix`** (followed by `#N`) → load and follow `${CLAUDE_SKILL_DIR}/operations/fix.md`
4. **If starts with `#N`**:
   - If remaining text after `#N` → jump to **Operation Comment**
   - If no remaining text → jump to **Operation Detail**
5. **Otherwise** → continue to **Create Flow**

For create flow, also extract:

```bash
MILESTONE_OVERRIDE=""
if echo "$ARGUMENTS" | grep -q '\-\-milestone'; then
  MILESTONE_OVERRIDE=$(echo "$ARGUMENTS" | sed 's/.*--milestone[[:space:]]*\([^[:space:]]*\).*/\1/')
  ARGUMENTS=$(echo "$ARGUMENTS" | sed 's/--milestone[[:space:]]*[^[:space:]]*//' | sed 's/^ *//')
fi
MILESTONE="${MILESTONE_OVERRIDE:-$CURRENT_VERSION}"
if [ -n "$MILESTONE" ]; then
  MILESTONE=$(bash "${CLAUDE_SKILL_DIR}/../scripts/wing-git.sh" milestone-resolve "$MILESTONE" 2>/dev/null)
fi
```

**Extract @assignee** (if present):
- Find `@mention` at end of description
- Resolve: `ASSIGNEE_GH=$(bash "${CLAUDE_SKILL_DIR}/../scripts/wing-config.sh" resolve-member "{mention}" 2>/dev/null)`
- Remove `@mention` from description text
- If resolve fails → output stderr (lists available members) → **STOP**
  💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

---

## Create Flow

### Step 1: Determine Assignee

**If `@assignee` explicitly provided**: Use that assignee (already validated).
**Solo mode** (members==1, no `@assignee`): Auto-assign to self (`$GH_USER`).
**Team mode** (members > 1, no `@assignee`): List members via script, then AskUserQuestion.

### Step 2: Analyze User Description

Before generating fields, separate bone from flesh: what is the irreducible need vs. the user's suggested approach? Test the opposite: what happens if we don't do this at all? If you see a better framing, say so before proceeding.

**Factual claims are hypotheses, not axioms.** When the user states something is broken, missing, or misconfigured — verify before accepting.

From the user's input, determine:
1. **Title** — concise, prefixed with type (`fix:`, `feat:`, `refactor:`, etc.)
2. **Priority** — P0/P1/P2/P3 (default P2)
3. **Size** — S/M/L/XL (default M)
4. **Domain** — from config `tags.services`
5. **Objective** — 1-2 sentence summary
6. **Success Criteria** — concrete, verifiable checkboxes
7. **Sub-tasks** — implementation steps
8. **Context** — background, relevant files
9. **Constraints** — limitations
10. **Verification** — how to verify completion

### Step 2b: Duplicate Check

Extract 2-3 key terms, search:

```bash
gh search issues "{key terms}" --repo "$REPO" --state open --json number,title,url --limit 10
```

AI evaluates semantic similarity. If potential duplicates found → display and AskUserQuestion: "None of these — create new" / "Update #N instead" / "Cancel".

If no results or no semantic match → continue silently.

### Step 3: Scan Codebase for Context

Use `Glob` and `Grep` to find relevant files, test files, config files. Add to Context section.

### Step 4: Generate Issue Body

```markdown
## Mission Contract

### Objective
{1-2 sentence summary}

### Priority
{P0|P1|P2|P3}

### Estimated Size
{S|M|L|XL}

### Primary Domain
{domain}

### Success Criteria
- [ ] Criterion 1 — concrete, verifiable
- [ ] Criterion 2

### Sub-tasks
- [ ] Step 1
- [ ] Step 2

### Context & References
{Background, relevant files, links}

### Constraints
{Limitations, dependencies}

### Verification Method
```bash
# verification commands
```
```

**Standalone Test**: Re-read the generated body as if you never saw the user's description. Can someone execute from this alone?

**Canonical field names** (wing-claim parses by name):
- `Success Criteria` (NOT "Acceptance Criteria")
- `Context & References` (NOT "Context")
- `Sub-tasks` (NOT "Tasks" or "Steps")

### Step 5: Preview & Confirm

Display preview with labels, milestone, formatted body. State your position as AGREE+EXTEND, CHALLENGE+PROPOSE, or REDIRECT — never silent compliance. When proposing an alternative, state what it costs.

Use `AskUserQuestion`: "发布" / "编辑标题" / "编辑优先级" / "取消"

If edit → re-preview (loop). If cancel → **STOP**
  💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.

### Step 6: Publish

```bash
gh issue create \
  --repo "$REPO" \
  --title "{title}" \
  --body "{formatted body}" \
  --label "$MISSION_LABEL" \
  --label "${STATUS_PREFIX}wip" \
  --label "${PRIORITY_PREFIX}{Pn}" \
  --label "domain:{domain}" \
  --label "size:{size}" \
  --assignee "{assignee_github}" \
  ${MILESTONE:+--milestone "$MILESTONE"}
```

If milestone doesn't exist → AskUserQuestion: "创建 milestone" / "跳过"

### Step 7: Notify

```bash
bash "${CLAUDE_SKILL_DIR}/../scripts/wing-notify.sh" mc.created \
  --issue "$ISSUE_NUMBER" --title "$TITLE" --assignee "$ASSIGNEE"
```

Non-fatal: if notification fails, warn but continue.

### Step 8: Output

📦 **CREATED** ── **#{issue}** {title} ─────────────
  **Assignee**  @{assignee}
  {priority_dot} `{Pn}`  **Size** `{size}`  **Milestone** `{milestone}` or "—"
  **URL**       `{issue URL}`

Use `AskUserQuestion`: "Claim 这个 Issue" / "继续创建" / "看 dashboard"

---

## Operation Detail — View Issue

> Triggered by `/wing-issue` **#042** (number only).

```bash
gh issue view $N --repo "$REPO" --json number,title,body,labels,milestone,assignees,state,url,comments
```

If not found → "**ERROR:** Issue **#N** not found." → **STOP**

Format:

🎯 **MC** **#{N}** ── {title} ──────────────────────
  **Assignee**   @{assignee or "unassigned"}
  {priority_dot} `{Pn}`  **Size** `{size}`  **Status** *{status}*
  **Milestone**  `{milestone}` or "—"  **State** `{open/closed}`
  **URL**        `{url}`
────────────────────────────────────────────
{Full Issue body}
────────────────────────────────────────────

Check for branch and PRs:

```bash
git ls-remote --heads origin "mission/${N}-*" 2>/dev/null
gh pr list --repo "$REPO" --search "Closes #$N" --state all --json number,url,state --limit 5
```

🔀 **BRANCH**  *{branch}* or "(not created)"
🔀 **RELATED PRs**  {list or "(none)"}

Display human comments if any (exclude bot/lifecycle comments). Show latest 5 if >5.

---

## Operation Comment — Add Comment to Issue

> Triggered by `/wing-issue` **#042** `some text here`.

```bash
gh issue comment $N --repo "$REPO" --body "$TEXT"
```

💬 **COMMENT** ── **#{N}** {title} ─────────────────
  {first 80 chars of text}...
  **URL** `{comment URL}`

---

## Error Handling

- Empty description → AskUserQuestion with `useTextArea: true` to get description
- GitHub API error → "**ERROR:** Failed to create Issue. Check `gh auth status`." → **STOP**
- Label not found → create without, warn to run `/wing init`
- Assignee not in config → "**ERROR:** Member '@{mention}' not found. Available: {list}" → **STOP**
- Notification fails → warn but continue (non-fatal)

**On any STOP:** Always append: 💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.
