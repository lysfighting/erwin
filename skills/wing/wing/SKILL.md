---
name: wing
description: "Wing dashboard + init. First time? Try: /wing help"
version: "0.1.0"
---

# /wing — Dashboard + Init

> GitHub-first autonomous development protocol. Push-based task assignment.

**Visual Encoding** — apply to ALL output (see `${CLAUDE_SKILL_DIR}/reference/ve-standard.md`)

**User input**: $ARGUMENTS

**Argument routing:**

| Input | Action |
|-------|--------|
| (empty) | Auto: first run → Init, subsequent → Dashboard |
| `#N` or number | Show MC **#N** details |
| `help` or `-h` | Quick-start guide |
| `config` | View current config |
| `init` | Force re-initialize |
| `queue` | All open mission Issues |
| `doctor` | Health diagnostics |
| `doctor fix` | Interactive repair |
| `wrap` | Retroactive closure: code-first → issue + PR |
| `sync` | Rebase branch on base |
| starts with `share` | Share knowledge to sibling MC |
| starts with `auto` | One-shot: claim → drive → ship |

---

## Route by Argument

- `help` or `-h` → **Operation Help**
- `config` → **Operation Config**
- `#N` or number → **Operation MC Detail**
- `init` → load `${CLAUDE_SKILL_DIR}/operations/init.md` (force init)
- `queue` → **Operation Queue**
- `doctor help` → **Operation Doctor Help**
- `doctor fix` → load `${CLAUDE_SKILL_DIR}/operations/doctor.md` (fix mode)
- `doctor` → load `${CLAUDE_SKILL_DIR}/operations/doctor.md`
- `sync` → same as `/wing-ship sync`
- starts with `share` → **Operation Share**
- `wrap` → **Operation Wrap**
- starts with `auto` → load `${CLAUDE_SKILL_DIR}/operations/auto.md`
- empty → **Step 0** (auto-detect init vs dashboard)

---

## Operation Help

**WING v0.1.0** — AI Autonomous Development Protocol
────────────────────────────────────────────

**SETUP**
  `/wing`                 Auto-detect project, create config + labels
  `/wing init`            Force re-initialize

**DAILY WORKFLOW**
  `/wing`                 Role-based dashboard
  `/wing` **#042**        View MC details
  `/wing doctor`          Health diagnostics
  `/wing doctor fix`      Interactive repair
  `/wing-issue <desc>`    Create MC (solo: self-assign; team: prompt)
  `/wing-issue <desc> @user`  Create + assign
  `/wing-claim`           Show assigned missions, pick one
  `/wing-claim` **#N**    Claim → Contract + branch
  `/wing-drive`           Execute mission
  `/wing-ship`            Push + PR
  `/wing-ship review`     AI code review
  `/wing-ship done`       Post-merge cleanup

**ONE-SHOT**
  `/wing auto` **#042**   Claim → Drive → Ship
  `/wing auto`            Pick from queue, then auto

**RETROACTIVE**
  `/wing wrap`            Code-first → issue + PR + labels

**RC LIFECYCLE**
  `/wing-rc`              Prepare RC → staging
  `/wing-rc promote`      RC → production + tag + release
  `/wing-rc fix`          Cherry-pick to RC
  `/wing-rc hotfix`       Fix production + backport

**CONFIG**
  `/wing config`          View config
  `/wing queue`           Browse open missions
────────────────────────────────────────────

**STOP**

---

## Step 0: Prerequisites

```bash
GH_USER=$(gh api user --jq '.login' 2>/dev/null)
if [ -z "$GH_USER" ]; then
  echo "ERROR: Run 'gh auth login' first."
  exit 1
fi
WING_DIR=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" detect-dir 2>/dev/null) || WING_DIR=""
```

## Step 1: Route — Init or Dashboard?

If `$WING_DIR` is empty OR `$ARGUMENTS` is `init`:
  → load `${CLAUDE_SKILL_DIR}/operations/init.md`

If config exists → check membership:
```bash
USER_LEVEL=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" get-user-level "$GH_USER" 2>/dev/null)
```
If user not in roster → auto-add via onboarding → continue to Dashboard.

→ Continue to **Dashboard**

---

## Operation Config

```bash
cat $WING_DIR/config.yml
```

Display config with edit instructions. **STOP**.

## Operation Queue

```bash
eval "$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" resolve-labels 2>/dev/null)"
gh issue list --label "$MISSION_LABEL" --state open --json number,title,assignees,labels,milestone --limit 50
```

Sort by priority (P0 first). Format with VE: `{priority_dot} **#{N}** {title} @{assignee}`

**STOP**

## Operation MC Detail

Fetch Issue + branch + PRs. Display formatted detail with VE. **STOP**.

## Operation Share

Post `🔗 **Cross-MC**` comment to sibling Issue. Human gateway (AskUserQuestion before posting). **STOP**.

## Operation Wrap

Retroactive closure for code-first workflows: detect state → create Issue → create PR → label. **STOP**.

---

## Dashboard

### Fetch data (parallel)

```bash
eval "$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" resolve-labels 2>/dev/null)"
BASE_BRANCH=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" conventions.base_branch "main" 2>/dev/null)
CURRENT_VERSION=$(bash "${CLAUDE_SKILL_DIR}/scripts/wing-config.sh" versions.current "" 2>/dev/null)
```

**Parallel API calls:**
1. All mission issues: `gh issue list --label "$MISSION_LABEL" --state open --json number,title,assignees,labels,milestone --limit 100`
2. Open PRs: `gh pr list --json number,title,author,headRefName,statusCheckRollup,reviewDecision --limit 20`
3. All open issues (untracked detection): `gh issue list --state open --json number,title,labels,assignees --limit 100`
4. Merged PR count: `gh pr list --state merged --base "$BASE_BRANCH" --json author --limit 100`
5. Milestone progress (if version set): `gh api repos/$REPO/milestones --jq ...`
6. RC branch: `wing-git.sh find-rc`
7. Merged local branches: `wing-git.sh list-merged-branches`
8. Local git state: branch, status, worktree

**Client-side filtering** from Call 1: WIP (status:wip), review (status:review), unassigned (empty assignees). Priority from labels.

### Dashboard layout

**MY DASHBOARD** — {repo} ──────────────────
Wing v0.1.0 · **{VERSION}**  `{progress_bar}`  {pct}%  {done}/{total}

{If WORKTREE_ENABLED:}
📂 **Worktrees** *(worktree mode)* 🟥 main · 🟦 develop · 🟧 rc · 🟩 mission · ⬜ other
  ⭐ {color} `{path}` ← *{branch}* → **#{issue}** {title}
  · {color} `{path}` ← *{branch}* → **#{issue}** {title}

{If RC_BRANCH:}
🚀 **RC** *{RC_BRANCH}*
  `/wing-rc promote` · `/wing-rc fix` · `/wing-rc cancel`

👥 **TEAM** ({count})
  {username} ({role})
    🟢 **#{N}** {title} {priority_dot} → *wip* {If STALE: ⚡}
    *{branch}* → `/wing-ship`
    🔀 **#{pr}** — {ci_icon}

🔀 **Under Review** ({count})
  🟡 **#{N}** {title} @{assignee}

📋 **Unassigned** ({count})
  {priority_dot} **#{N}** {title}

{If untracked:} ⚠️ **Untracked** — `/wing-issue fix` **#N**
{If needs cleanup:} 🧹 **Needs Cleanup** — `/wing-ship done`

📊 **Merges** {user}: `{n}` · {user}: `{n}`

💡 **Next:** {context-specific suggestion}

`/wing help` · `/wing doctor` · `/wing #N` · `/wing config`

### NEXT suggestion logic (first match)

1. RC exists + all missions done → `/wing-rc promote`
2. Orphan contract → `/wing doctor fix`
3. Merged branch exists → `/wing-ship done`
4. PR behind base → `/wing-ship sync`
5. PR CI passed + approved → merge on GitHub
6. PR CI failing → fix before merge
7. Active mission, unchecked tasks → `/wing-drive`
8. Active mission, tasks done → `/wing-ship`
9. No mission, unassigned exist → `/wing-claim`
10. No mission, no unassigned → `/wing-issue`

### Priority dots

🔴 P0 · 🟠 P1 · 🟡 P2 · ⚪ P3

### Progress bar

Map % to 10 chars: `█` filled, `░` empty. Wrap in backticks: `` `█████░░░░░` ``

---

## Error Handling

- No `gh` auth → "**ERROR:** Run `gh auth login`." → **STOP**
- No config → redirect to init
- GitHub API fails → degrade gracefully (show what's available)

**On any STOP:** 💡 Something wrong? Run `/wing doctor` to diagnose, or `/wing doctor fix` to auto-repair.
