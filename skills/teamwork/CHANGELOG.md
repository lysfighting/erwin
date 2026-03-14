# Changelog — Teamwork Skills

All notable changes to the Teamwork skill set.

## [3.11.2] — 2026-03-14

### Fixed — VE signal dilution (final hotfix before freeze)

Root cause: output templates were already VE-compliant, but **non-output text** contained 22 bare `#N`/`#NN`/`#NNN` patterns competing with 42 correct `**#NNN**` patterns (dilution ratio 0.52). LLM averages across all patterns in the prompt, producing mixed output.

**Fixes by category:**
- **team.md**: ASCII diagrams `#42` → `#042` (3 lines), usage example `#42` → `{issue}`, code comment `#42` → `#NNN`, removed ❌ WRONG example lines (2x) that showed bare issue numbers — kept only ✅ CORRECT examples, rewrote VE rule text to avoid bare number references
- **team-issue.md**: trigger description `#42` → `{issue}`
- **team-ship.md**: argument description `#42` → `{issue}`, error messages `#042` → `{issue}`, help text `#042`/`#043` → `{issue}`/`{pr}`
- **team-rc.md**: GitHub release notes template `(#042)` → `(**#042**)` (4 lines)

**Result**: Dilution ratio reduced from 0.52 → ~0.07 (3 remaining bare patterns are in code blocks where markdown doesn't render).

**Note**: This is the final maintenance release. Teamwork v1 is now frozen. Future development moves to tw2 architecture.

## [3.11.1] — 2026-03-13

### Fixed — VE hyperlink enforcement (issue number bold compliance)

- **team.md dashboard templates**: Replaced abstract "mentally check each line" instruction with concrete ✅/❌ examples showing correct `**#306**` vs wrong `#306` formatting. LLMs mimic examples far more reliably than they follow meta-instructions. Applied to both Solo/Leader and Member dashboard templates.
- **Root cause**: Pre-instructions ("before you output, check...") don't work for streaming LLM generation — the model can't retroactively fix tokens. Concrete correct/wrong examples adjacent to output templates drive pattern-matching compliance.

## [3.11.0] — 2026-03-13

### Added — Constitution Principle 15: Never End on a Period + 10 UX dead-end fixes

**New design principle**: Every skill output must end on a question (AskUserQuestion), not a period. Three patterns: success transition → choice, error recovery → fix options, terminal state → text OK.

**10 lifecycle transition dead-ends fixed** (AskUserQuestion added):
- **team-issue**: Issue created → "Claim? / Create another? / Dashboard?" (Step 8, Fix/Adopt, Batch)
- **team-claim**: Claim completed → "Start drive? / Review Contract? / Dashboard?" (branch mode + "放弃" path)
- **team-drive**: All tasks done → "Ship? / Self-review? / Dashboard?" (early exit, standard completion, post-manual)
- **team-ship**: Incomplete tasks → "Drive? / Mark done? / Cancel?" + CHANGES_REQUESTED → "View? / Fix? / Cancel?"

**Constitution**: Added Principle 15 to `docs/teamwork-design-constitution.md` (v1.1)

## [3.10.9] — 2026-03-13

### Fixed — Challenge mechanism + dashboard audit (4 fixes)

- **team-issue.md Step 2**: "Factual claims are hypotheses, not axioms" — verify system state before locking claims into sub-tasks.
- **team-claim.md Step 4c**: "Detail ≠ accuracy" — falsify concrete claims before AGREE.
- **team-drive.md Step 1**: MANUAL tasks now go through Plan Validation (not skipped). Resume scenario correctly skips re-assessment.
- **team-drive.md Plan Assessment scope**: "ALL types (code + MANUAL)" — MANUAL is classification, not exemption.
- **team-drive.md Manual Ops Handoff**: AskUserQuestion upgraded to ⛔ HARD STOP with explicit options. "不需要" option now has clear destination.
- **team.md Worktree paths**: Use absolute paths from `git worktree list`, no abbreviation (prevents LLM mis-shortening).
- **team.md My Missions**: Explicit worktree enabled/disabled conditional. Shows `cd {path} && claude` only when worktree mode active.

## [3.10.8] — 2026-03-13

### Fixed — MANUAL classification: process step instead of tool list

- **team-drive.md Plan Validation**: Replaced tool-enumeration approach (listing `gh`, `docker`, `curl`...) with a 3-step executability check: (1) can I do it? → code task. (2) not sure? → ask user what access is available. (3) user confirms no → MANUAL. One mechanism, not a growing list.
- **team-drive.md Verify Acceptance Criteria**: "Attempt verification before punting" — try first, ask if uncertain.
- **Why**: Previous v3.10.7 fix listed specific CLI tools — same O(2^C) scaling problem it tried to solve. The right fix is a decision process, not an enumeration.

## [3.10.6] — 2026-03-13

### Fixed — team-drive Manual Ops Handoff AskUserQuestion not firing

- **Root cause**: AskUserQuestion instruction was placed AFTER the output template as a separate block. The AI treated the template (CODE COMPLETE + YOUR TURN + `/team-ship`) as a complete output and never reached the question.
- **Fix**: Embedded AskUserQuestion directly inside the YOUR TURN output block. Removed premature `/team-ship` closing line. Added explicit "DO NOT output closing line before AskUserQuestion" instruction.
- **UX impact**: After code tasks complete with manual tasks remaining, the user now gets an interactive choice (guide me / done / ship first / skip) instead of a passive text dump.

## [3.10.4] — 2026-03-13

### Fixed — v3.10.0-v3.10.3 audit (8 HIGH, 14 MEDIUM, 18 LOW)

**HIGH (8 — functional failures)**
- **tw-contract.sh**: `authorAssociation == 'BOT'` filter was no-op (GitHub doesn't have BOT enum value) — now checks `author.login.endswith('[bot]')`
- **tw-contract.sh**: `--since` without argument crashed on `shift 2` — added `$#` validation
- **tw-contract.sh**: Missing contract file output `CONTRACT:NETWORK_ERROR` → now `CONTRACT:NO_FILE`
- **team-drive.md**: `check-comments` returns `NEW:unknown` on first run but skill only matched `"NEW"` — now matches prefix
- **team-drive.md**: Timestamp stored after `NETWORK_ERROR` → permanent comment loss — now conditional
- **team-issue.md**: `authorAssociation == "BOT"` referenced non-existent `gh issue view` field — now uses `login.endswith('[bot]')`
- **team-ship.md**: `grep -oE '[\w]'` non-portable on macOS BSD grep — replaced with `[A-Za-z0-9_]`
- **team.md SY4**: `ISSUE_NUMBER` and `CONTRACT_PATH` never defined in Operation Sync — added extraction in SY1

**MEDIUM (14 — logic/VE/UX)**
- **team-ship.md**: 2 STOPs missing doctor hint (CHANGES_REQUESTED, P0 返回修复) — added
- **team-ship.md**: Checkbox grep `^` anchor missed indented `- [ ]` → false "all done" — removed anchor
- **team-ship.md**: Unquoted `$ISSUE_NUMBER` in 2 bash calls — quoted
- **team-drive.md**: `{ISSUE_NUMBER}` vs `{issue}` variable inconsistency — unified to `{issue}`
- **team-drive.md**: `read-comments` JSON output not documented for LLM parsing — added note
- **team-drive.md**: MEMBER/OWNER role check lacked mechanism — specified `authorAssociation` field
- **team-claim.md**: P11 violation — inline comment filtering duplicated script logic — now delegates to `tw-contract.sh read-comments`
- **team-claim.md**: `last_comment_check` double-write (template + script) — removed from template
- **team.md SY4**: `PR:{status}:{count}` doc wrong (script outputs `:{pr_number}`) — fixed
- **team.md SY4**: `COMMENTS:NEW:unknown` undocumented — added
- **team.md SY4**: All NETWORK_ERROR → silent skip — now shows warning
- **team.md SY4**: Timestamp updated after NETWORK_ERROR — now conditional
- **team-issue.md/team-ship.md/team-drive.md**: `*@{author}*` italic → `` `@{author}` `` backtick (VE: author is a value)
- **team-rc.md**: `--ref $PROD_BRANCH` → `--ref "$VERSION"` for deterministic deploy from tag

**LOW (6 selected fixes)**
- **tw-contract.sh**: `import re` inside for-loop → moved to top
- **tw-contract.sh**: Added issue number validation (`^[0-9]+$`) in check-comments and read-comments
- **team-issue.md**: Clarified truncation rule and ">5 human comments after filtering"
- **team-claim.md**: Added explicit zero-comments handling
- **team-rc.md**: Added `$TEAMWORK_DIR/config.yml` path in deploy guidance
- **team-rc.md**: VE — backtick-wrapped `{PROD_WORKFLOW}` in summary line

## [3.10.3] — 2026-03-13

### Added — Production deploy guidance after promote

- **team-rc.md Step P10b**: After promote (tag + release), trigger production deploy if `deploy.production_workflow` is configured. If not configured, guide user to deploy manually with actionable next step.
- **Config**: New optional `deploy.production_workflow` field (parallel to existing `deploy.staging_workflow`).
- **Why**: Promote closes the release cycle but previously didn't close the deployment cycle. Users were left wondering "how do I deploy?" after a successful promote.

## [3.10.2] — 2026-03-13

### Changed — Pre-flight tests skip when CI is configured

- **team-ship.md Step 2d**: When `quality.ci: true`, skip local pre-flight tests entirely. Push immediately for remote persistence; CI guards merge.
- **Why**: Pre-flight blocks push, but push's first job is remote persistence. When CI already guards merge, blocking push for a redundant local check delays persistence for no safety gain.
- **`ci: false` unchanged**: Local tests remain the only quality gate when CI is not configured.

## [3.10.1] — 2026-03-13

### Fixed — PR Closes #N 导致 Issue 被提前关闭

- **team-ship.md Step 4b**: PR body 无条件使用 `Closes #N`，即使 Issue 有未完成的 checkbox。合并后 GitHub 自动关闭 Issue，但任务未完成。
- **Root cause**: Ship 只检查 Contract 的 code tasks，manual tasks 不阻塞。但 `Closes` 关键字让 GitHub 在 merge 时无条件关闭 Issue。
- **Fix**: Step 4b 生成 PR body 前检查 Issue checkbox 状态。全部 `[x]` → `Closes #N`；有 `[ ]` 未完成 → `Refs #N` + PR body 附未完成 tasks 列表 + 本地警告。
- **Actual case**: vi_agent Issue #286 + PR #291

## [3.10.0] — 2026-03-13

### Added — Information Sync Layer (信息同步层)

**Core: Comment reading throughout lifecycle**
- `/team-issue #N`: displays human comments after Issue body (filtered: no bot, no lifecycle)
- `/team-claim`: incorporates Issue comments into Contract as `## Pre-Claim Discussion` section; stores `last_comment_check` timestamp in frontmatter
- `/team-drive` Step 0b1.5: checks for new comments since last check; direction-change heuristic triggers AskUserQuestion (keywords + MEMBER/OWNER)
- `/team-ship` Step 2b-bis: comment check + PR review check (`CHANGES_REQUESTED` → hard stop)
- `/team-drive` completion: Cross-MC nudge when sibling active Contracts exist

**Core: Sync semantic extension**
- `/team sync` SY4: after rebase, displays INFO SYNC block (new comments, PR review status, Contract freshness)
- `/team-ship sync` S4: identical Info Sync step
- Advisory only — sync never stops or prompts
- Hidden when no new information (zero noise)

**Script: tw-contract.sh — 4 new subcommands**
- `check-comments {issue} {stored_timestamp}`: lightweight REST API check, returns UNCHANGED/NEW/NETWORK_ERROR
- `read-comments {issue} [--exclude-bot] [--exclude-lifecycle] [--since TS]`: Python3-filtered full read, returns JSON
- `store-comment-timestamp {contract} {timestamp}`: updates Contract frontmatter
- `info-sync {issue} {contract}`: bundled 3-check (comments + PR + freshness), outputs structured 3-line result

**Design: Performance characteristics**
- `/team-issue` and `/team-claim`: +0 API calls (merged into existing `gh issue view`)
- `/team-drive` and `/team-ship`: +1-2 API calls (~1.3-2.5s typical)
- All new calls use `|| true` non-fatal pattern — network errors skip info sync, never block workflow
- Python3 filtering handles emoji-prefixed lifecycle comments cross-platform

## [3.9.15] — 2026-03-13

### Fixed — 35 findings from 3-agent parallel audit (VE 7, UX 18, Logic 10)

**HIGH (6)**
- **team-rc.md P5**: PREV_TAG included hotfix tags, corrupting release notes date range — added grep filter for clean version tags only
- **team-rc.md H2**: Hotfix version derivation used stale local tags — added `git fetch --tags` before tag listing
- **Multiple files**: Inconsistent `base_branch` default ("main" vs "develop") — added explicit warning when not configured in config.yml
- **team-ship.md**: No Contract recovery when on mission branch — added branch detection + `/team-claim` suggestion (matching team-drive)
- **team-ship.md Step 2c**: Auto-committed uncommitted changes without confirmation — added AskUserQuestion with Commit/Review/Stash options
- **ONBOARDING.md**: Documented `/team-claim` as auto-selecting unassigned issues; actual behavior lists assigned issues — fixed docs

**MEDIUM (14)**
- **team.md Doctor Check 5/10**: Bare `#{issue}` in diagnostic output — added bold wrapping for VE compliance
- **team.md Wrap B**: `COMMIT_COUNT` was 1 when `COMMIT_SHAS` empty (`echo "" | wc -l` = 1) — added empty guard
- **tw-git.sh find-rc**: No `git fetch` before listing remote RC branches — added `git fetch --prune`
- **tw-git.sh list-merged-branches**: Hardcoded `mission/*` prefix — now reads branch pattern from config
- **team.md Next hint**: Missing "PR open, CI pending" state — added rule 5b
- **team.md Next hint**: Didn't distinguish tasks-remaining vs all-done-not-shipped — split into rules 7a/7b
- **team.md Next hint**: `#{N}` unspecified + empty queue case — split into rules 8a (suggest highest-priority)/8b (suggest `/team-issue`)
- **team-claim.md**: No-missions message says "ask team lead" in solo mode — now context-aware
- **team-ship.md**: BEHIND_BASE > 30 hard stop — changed to AskUserQuestion with Rebase/Ship options
- **team-drive.md**: Branch-behind-base warning only — changed to AskUserQuestion with Sync/Continue options
- **team-drive.md**: Completion always suggests `/team-ship` — now checks if PR already exists
- **team-claim.md**: Not-in-roster suggests `/team init` — changed to `/team` (automatic onboarding)

**LOW (15)**
- **team-ship.md**: Help output `#42`/`#43` → `#042`/`#043` for OSC 8 auto-linking
- **team-drive.md**: `TEAMWORK_DIR` fallback echo outside `$()` — fixed assignment syntax
- **tw-notify.sh**: `team.repo` fails on erwin-style config — added `git remote` fallback
- **team-rc.md F6**: `git checkout -` → explicit `git checkout "$BASE_BRANCH"`
- **team-rc.md H7**: MERGE_SHA null guard added
- **team-rc.md Prepare**: NEXT STEPS now mentions `/team-rc fix` as surgical alternative
- **team-claim.md**: "Continue mission" outputs text only — added AskUserQuestion to invoke `/team-drive`
- **team-claim.md**: Closed-issue error now suggests `gh issue reopen` + `/team queue`
- **team-issue.md**: Empty description hard stop — changed to AskUserQuestion with `useTextArea: true`
- **team.md U12**: Alert block already covers orphan contracts — no change needed (verified)
- **team.md U16**: Worktree-mode-disabled-but-in-worktree — added re-enable option
- **team-ship.md Done**: Branch-aware error message when issue not found
- Other: 4 minor VE consistency fixes in help templates

## [3.9.14] — 2026-03-13

### Fixed — Release notes robustness (3-teammate peer review)

**H6 hotfix variables (Critical — was broken):**
- Replaced undefined `${COMMIT_TITLE}`, `${GH_USER}`, `${PR_URL}`, `${ISSUE_NUMBER}` with `gh pr view` data extraction
- Fixed `${PREV_TAG}` → `${LATEST_TAG}` (correct variable from H2 scope)

**P5 explicit date filtering:**
- Added explicit instruction for AI to filter PRs by `mergedAt > PREV_DATE` (was implicit, risked including all historical PRs)
- Added first-release guard: omit Full Changelog line if no previous tag
- Removed `**bold**` from issue numbers in release notes template (GitHub auto-links `#NNN` natively; bold was VE for CLI, not GitHub)

**Promote idempotency (P7 + P8):**
- P7: Skip tagging if tag already exists (re-entry after partial failure)
- P8: Skip release creation if release already exists

**Shell escaping:**
- P8 + H6: `--notes "$BODY"` → `--notes-file <(echo "$BODY")` to avoid shell escaping issues with PR titles containing quotes

**Hotfix active RC warning:**
- H8 summary: detect active RC branch and warn user to `/team-rc fix` or re-cut

**Docs staleness:**
- `skills/README.md`: 3.9.10 → 3.9.14 (version table)
- `CLAUDE.md`: v3.9.10 → v3.9.14 (teamwork line)

## [3.9.13] — 2026-03-13

### Fixed — RC release notes generation for squash merge

- **`/team-rc promote` Step P5**: Replaced `git log` with `gh pr list --base` to reconstruct release notes from merged PR history. Categorizes by commit type (feat/fix/refactor/other), includes author attribution (`by @author`) and PR links. Solves: squash merge to main loses all original PR metadata — GitHub's `--generate-notes` only sees the release PR itself.
- **`/team-rc promote` Step P8**: Replaced `--generate-notes` with `--notes "$RELEASE_BODY"` to use P5's output.
- **`/team-rc hotfix` Step H6**: Same fix — generates hotfix-specific release body instead of relying on `--generate-notes`.

## [3.9.12] — 2026-03-13

### Fixed — VE consistency sweep + identity language + script hardening

**Visual Encoding — hyperlink preservation (20+ fixes across 6 skills):**
- Split command backticks from issue numbers to preserve OSC 8 auto-linking: `` `/team-claim #042` `` → `` `/team-claim` **#042** `` (team.md, team-issue.md, team-claim.md, team-ship.md, team-drive.md)
- Added `**bold**` to issue number references in routing tables, help text, trigger descriptions, and footer lines across all 6 skill files
- Zero-padded issue numbers in non-code-block examples (ONBOARDING.md, team-rc.md release notes)
- Wrapped status values in backticks in output templates: CI status, review status, labels, milestone, URL, worktree status (team-ship.md, team.md Auto section, MC Detail)

**Identity language (Constitution P3):**
- Strengthened team-ship.md: process description → ownership framing ("You own the delivery. Guard the quality boundary...")
- Strengthened team-rc.md: process description → ownership framing ("You own the release cadence. Every RC is a contract with users...")

**Script hardening:**
- `tw-contract.sh`: mktemp uses `${TMPDIR:-/tmp}` instead of hardcoded `/tmp`

### Root cause analysis

Dashboard hyperlinks were lost because non-output areas (routing tables, help text, trigger descriptions) used plain `#N`/`#42` ~40 times, diluting the correct `**#NNN**` pattern in output templates (~50 times). The LLM averaged across all examples, producing unformatted issue numbers. Fix: made non-output areas consistent with VE standard.

## [3.9.7] — 2026-03-13

### Fixed — Audit sweep (security, VE, logic, docs)

**Security (2 critical):**
- **eval injection**: `setup-github-labels.sh` and `tw-label.sh` — replaced `eval "$(resolve-labels)"` with safe line-by-line parsing
- **temp file race**: `tw-contract.sh` — replaced predictable `/tmp/tw-sync-body.$$` with `mktemp`

**Visual Encoding (14 fixes):**
- Removed 5 code block fences wrapping LLM output (config warning, merge summary, tab title, desktop notification, doctor example)
- Added `**bold**` to 7 `#{issue}`/`#{pr}` template references in output sections
- Wrapped `status:*` label values in backticks in Operation Learn state machine
- Zero-padded issue numbers in doctor example (`#42` → `**#042**`)

**Skill logic (5 fixes):**
- Routing table: added missing `sync` and `share` entries
- `/team help`: added `/team-claim` (empty argument form)
- `/team-ship done`: derive main repo path before worktree removal
- `/team auto`: explicit Step 0/0b/0b2/0c skip instruction
- `/team-drive`: "ACCEPTANCE CRITERIA" → "SUCCESS CRITERIA" (canonical name)

**Script fixes (4):**
- `tw-bump-version.sh`: CHANGED counter now includes inline updates; shebang consistency
- `tw-pr.sh`: `exit 0` → `return 0` in function error handlers
- `tw-git.sh`: `+` prefix from worktree branches now stripped in `list-merged-branches`

**Docs staleness (5):**
- `skills/README.md`: version table 3.9.4 → 3.9.7
- `docs/teamwork-ai-manual.md`: "Nine Scripts" → "Ten Scripts", added tw-unit-test.sh to diagram
- `CLAUDE.md`: version 3.9.6 → 3.9.7
- `_project/CURRENT.md`: "9 scripts" → "10 scripts"

## [3.9.6] — 2026-03-13

### Fixed — Cross-platform + remaining script fixes

- **sed portability**: All `sed -i ''` (macOS-only) replaced with `_sed_i` helper that detects GNU/BSD sed. Affects tw-bump-version.sh and tw-config.sh migrate. Linux users can now run these scripts.
- **tw-pr.sh**: Undefined `$repo_flag` in error messages replaced with `${repo_args[*]}`.
- **team-rc.md cancel**: "cancel each with `/team-rc cancel`" (infinite loop) replaced with manual delete guidance.
- **team-ship.md done**: Stash `push -u` (stashes untracked) changed to `push` (tracked only) — consistent with dirty check that excludes untracked files.

## [3.9.5] — 2026-03-12

### Fixed — Script hardening + docs staleness + ecosystem robustness

**Script hardening (12 fixes):**
- **tw-git.sh rebase**: Uses local ref when fetch fails (was using stale `origin/` ref)
- **tw-git.sh milestone-resolve**: jq injection fixed — `$name` passed via `--arg` instead of string interpolation
- **tw-git.sh commit**: Default changed from `git add -A` (dangerous, stages untracked) to `git add -u` (tracked only)
- **tw-pr.sh**: `$base_flag` and `$repo_flag` word-splitting replaced with proper bash arrays
- **tw-label.sh**: `$remove_labels` and `$remove_args` word-splitting replaced with bash arrays
- **tw-contract.sh**: Signal trap added for temp file cleanup (`/tmp/tw-sync-body.$$`)
- **tw-notify.sh**: `detect_config()` now worktree-aware via `tw-config.sh detect-dir`
- **team-rc.md promote**: `$MILESTONE_FULL` jq injection fixed with `--arg`
- **team-rc.md hotfix**: `git add -A` → `git add -u` + `git status` preview before commit
- **team-rc.md hotfix**: Re-entry detection — detects existing `hotfix/*` branch, skips creation
- **team-rc.md fix**: Stash pop uses exact stash index (was popping `stash@{0}` regardless)
- **team-ship.md**: Unknown subcommands now show error instead of falling through to full ship

**Docs staleness:**
- **skills/README.md**: Updated from v2.0.0 to v3.9.x, added 3 missing skills
- **CLAUDE.md**: Teamwork version v3.7.1 → v3.9.5
- **Idea pipeline**: cross-mc-knowledge-sharing → Fruit, project-management-template → Fruit
- **visual-test-skill.md**: drive-showcase → showcase references

**Ecosystem:**
- **install.sh**: Deprecated skill cleanup handles both symlinks and copy-mode files
- **install.sh**: `--target` with no value shows user-friendly error
- **erwin-init**: Added version frontmatter (v2.0.0)

## [3.9.4] — 2026-03-12

### Fixed — Worktree config bug + dashboard RC + unit tests

- **Worktree config fix**: `_detect_config()` in tw-config.sh now uses worktree-aware `_detect_dir()` logic — all subcommands (resolve-labels, resolve-member, list-members, etc.) now work correctly from worktrees. Same fix applied to `_teamwork_dir()` in tw-contract.sh.
- **team-rc fix stash guard**: Fix flow now stashes uncommitted changes before `git checkout` to RC branch, restores after cherry-pick.
- **Dashboard RC section**: New 🚀 RC section shows active RC branch, commit count, and action shortcuts. RC-aware NEXT suggestion logic.
- **Tips expansion**: 15 → 29 tips covering doctor fix, sync, share, wrap, fix, hotfix, cancel, quality gates, notifications, monorepo.
- **Unit tests**: New `tw-unit-test.sh` with 48 tests for tw-bump-version.sh (23), tw-config.sh migrate/has-section (9), tw-git.sh slugify/next-version (16). Zero external dependencies.

## [3.9.3] — 2026-03-12

### Fixed — Tooling + VE accuracy + docs staleness

- **tw-bump-version.sh path bug**: Manual path was `../../docs/` (2 levels up from skills/) instead of `../docs/` — manual was NEVER updated by the bump script. Fixed.
- **Version desync cleanup**: All inline `3.9.1` references updated to `3.9.2` via fixed bump script (team.md x5, team-issue.md x1, README.md x1, ONBOARDING.md x1, manual x2).
- **VE: Promote Help contradiction**: "never commit directly to RC" conflicted with `/team-rc fix`. Reworded to recommend fix for small issues, cancel+re-cut for big ones.
- **VE: Cancel Flow raw git**: Replaced `git push origin --delete` suggestion with `/team doctor fix` and `/team-rc cancel`.
- **VE: Worktree guard doctor hint**: Moved hint from inside code block (invisible) to outside.
- **_project/CURRENT.md**: Updated from v3.8.6 to v3.9.3 (was 4 versions behind).
- **_project/ROADMAP.md**: Complete rewrite — was severely stale (completed items unchecked, Teamwork absent).
- **Root CHANGELOG.md**: Added missing v3.8.7 entry.

## [3.9.2] — 2026-03-12

### Fixed — Tips system: `shuf` crash on macOS

- **Root cause**: Tips instruction said "pick one line at random" without specifying method. Claude chose `shuf` (GNU coreutils) which doesn't exist on macOS (Darwin).
- **Fix**: Tips master instruction (team.md) now explicitly requires Read tool for file access and prohibits `shuf`, `sort -R`, or any Bash-based randomization. All 20+ tip placeholders across 6 skill files inherit this instruction.

## [3.9.1] — 2026-03-12

### Fixed — Self-drive audit: bugs + VE + docs staleness

- **Hotfix continue routing bug**: Removed reference to non-existent `/team-rc hotfix continue` subcommand. Now says "run `/team-rc hotfix` again" (detects existing branch).
- **VE compliance**: Fixed help template variable case (`{base_branch}` → `{BASE_BRANCH}`, `{production_branch}` → `{PROD_BRANCH}`). Fixed `FIX APPLIED` line bold/italic nesting.
- **Manual footer version**: Updated from stale v3.1.0 → v3.9.1.
- **RC Freeze text**: Clarified "never commit directly" to mention `/team-rc fix` as surgical alternative.
- **team.md learn section count**: Fixed "13 sections" → "14 sections".
- **team-ship done footer**: Added `/team-rc` (release) to next-step suggestions.

## [3.9.0] — 2026-03-12

### Added — RC complete lifecycle (`/team-rc fix` + `/team-rc hotfix`)

- **`/team-rc fix`**: Cherry-pick a specific commit from develop to active RC branch. Preserves existing RC test progress (unlike cancel + re-cut). Supports commit hash, branch name, or interactive selection from recent commits. Guides through conflict resolution if needed.
- **`/team-rc hotfix`**: Post-release production fix workflow. Creates hotfix branch from main, guides interactive fix, PRs to main (squash), tags hotfix version, creates GitHub Release, backports (cherry-pick) to develop. Requires dual-branch mode (base != production). Single-branch projects get a clear error.
- **Help text updates**: `/team-rc help` now documents fix, hotfix, and the two bugfix strategies (upstream-first vs cherry-pick).
- **team.md**: RC LIFECYCLE section in help and learn updated with fix/hotfix entries.
- **Docs sync**: ONBOARDING.md, README.md, teamwork-ai-manual.md all updated with fix/hotfix commands and RC lifecycle Phase 4 (Hotfix).

## [3.8.7] — 2026-03-12

### Fixed — Docs completeness audit

- **ONBOARDING.md Advanced Features**: Added `/team sync` and `/team share [#N]` to the table (were missing since v3.8.5).
- **README.md skill table**: Added `/team sync` and `/team share [#N]` rows.
- **Root CHANGELOG.md**: Added v3.8.3–v3.8.6 summary entries (were only in teamwork CHANGELOG).

## [3.8.6] — 2026-03-12

### Fixed — Documentation audit + version sync

- **Version pollution cleanup**: Fixed 3 files where `--help` was written as version string (ONBOARDING.md, README.md, team-issue.md) — caused by `tw-bump-version.sh` accepting non-semver input in v3.8.3.
- **Manual version sync**: Updated `docs/teamwork-ai-manual.md` from stale 3.8.2 → 3.8.6. Added `/team sync` and `/team share` to command table.
- **Dashboard footer**: Added `/team sync` · `/team share` to both Leader and Member dashboard footers.
- **Learn section 14**: Added "CROSS-MC KNOWLEDGE SHARING" to `/team learn` covering L1 intake, L2 push, anti-plague rules.

## [3.8.5] — 2026-03-12

### Added — Cross-MC knowledge sharing + `/team sync` + `/team share`

- **`/team sync`**: Standalone rebase on base branch, available anytime (not just during ship). Same logic as `/team-ship sync`, but accessible as `/team sync` for ad-hoc use.
- **`/team share [#N]`**: Cross-MC knowledge sharing via GitHub issue comments. Posts structured `🔗 Cross-MC from #N` comments with context, finding, files, and confidence level. Anti-plague safeguards: human gateway (AskUserQuestion confirmation), no transitive propagation, source attribution.
- **L1: Cross-MC knowledge intake** (team-drive.md Step 0b2): `/team-drive` now reads issue comments on startup, filters for `🔗 Cross-MC` prefixed comments, and injects them as advisory context into the Contract. Advisory only — does not override MC objectives.

### Anti-plague rules (enforced across all skills)

1. Human gateway: `/team share` requires user confirmation before posting
2. No transitive propagation: received knowledge must NOT be re-shared
3. Source attribution: structured format with confidence level
4. Advisory by default: informs but does not override Contract

## [3.8.4] — 2026-03-12

### Fixed — Dashboard formatting + worktree init bug

- **Dashboard bold issue numbers**: Strengthened constraint — all issue numbers MUST be `**#NNN**` (bold-wrapped). Fixed contradictory example at priority_dot shorthand. Without bold, Claude Code renders plain text with no clickable hyperlink.
- **Dashboard `{Pn}` duplication**: Removed redundant `{Pn}` text from 6 template sites where `{priority_dot}` already replaces it. Prevents output like `🟠 P1 #222` (dot + text label).
- **Dashboard assigned vs claimed**: New ⚪ *assigned* state for issues assigned but not yet claimed (no Contract/branch). Previously all showed as 🟢 *wip*.
- **Worktree init bug**: `_detect_dir()` used relative path — fails when cwd is a worktree (`.teamwork/config.yml` lives in main repo). Now uses `git rev-parse --git-common-dir` to locate main repo. This was the root cause of intermittent `/team` triggering init flow.
- **Tab colors**: Replaced %1 土橙 `rgb(180,120,46)` → 珊瑚红 `rgb(200,90,90)` and %6 金黄 `rgb(180,140,60)` → 黄绿 `rgb(160,180,60)`. Old colors too similar to rc 橙.
- **Bump script validation**: `tw-bump-version.sh` now rejects non-semver input (e.g., `--help` was previously accepted as version string).

### Removed — Polluted global commands

- Deleted residual `CHANGELOG.md`, `ONBOARDING.md`, `README.md` from `~/.claude/commands/` and `vi_agent/.claude/commands/` (leftover from pre-filter install.sh).

## [3.8.3] — 2026-03-12

### Added — `/team-rc cancel` subcommand + upstream-first RC model

- **`/team-rc cancel`**: New Cancel Flow (Steps C1-C5) — delete current RC branch (local + remote) with AskUserQuestion confirmation. Enables the upstream-first bugfix model.
- **Upstream-first model**: RC is a frozen snapshot. Never commit directly to RC. Found a bug on staging? Fix on *develop* → `/team-rc cancel` → `/team-rc` (re-cut). No merge-back needed.
- **Prepare Flow**: When RC already exists, offers AskUserQuestion (Promote / Cancel / Exit) instead of plain-text STOP.
- **Promote Flow**: Explicit branch deletion (local + remote) replaces `--delete-branch` flag. Removed stale cherry-pick reminder.

### Fixed — Post-merge workflow doesn't check PR checkboxes

- **post-merge.yml**: Added Step 4 — check all `- [ ]` → `- [x]` in PR body after merge. Previously only checked Issue body checkboxes.
- **Promote Help**: Updated troubleshooting to reference upstream-first model instead of "fix on RC branch".

### Docs — RC lifecycle updated across all references

- **team.md help**: Added `/team-rc cancel` line.
- **team.md learn**: RC lifecycle section rewritten to upstream-first model with cancel flow.
- **docs/teamwork-ai-manual.md**: Command table + Phase 2 verify + RC Freeze Rule + Only One RC section all updated.

## [3.8.2] — 2026-03-11

### Added — Script沉淀: 6 new subcommands replace 15 inline bash sites

- **`tw-git.sh build-branch-name`**: Pattern-based branch name generation (replaces 3x inline sed in team.md, team-claim.md)
- **`tw-git.sh extract-issue-number`**: Extract issue number from branch name (replaces 3+ inline `grep -oE` in team-drive.md, team.md)
- **`tw-git.sh worktree-main-repo`**: Return main worktree path (replaces 3x inline `git worktree list | head -1 | sed` in team.md, team-claim.md, team-rc.md)
- **`tw-git.sh worktree-find-by-branch`**: Locate worktree by branch name (replaces 3x inline `grep -B2` + awk in team-claim.md, team-drive.md, team-ship.md)
- **`tw-config.sh team-mode`**: Returns "solo"/"team" based on member count (replaces LLM counting in team-issue.md)
- **Von Neumann principle**: deterministic ops → scripts, judgment → LLM. 15 inline bash sites migrated.

### Fixed — AskUserQuestion constraints + dashboard NEXT guidance

- **`/team auto` Issue selection**: Capped at 3 Issues + Cancel (4-option AskUserQuestion limit). Shows hint for 4+ Issues: use `/team auto #N` directly.
- **Dashboard NEXT item 4**: Changed "ready to merge → `/team-ship done`" to "merge on GitHub, then `/team-ship done`" (done is POST-merge cleanup, not the merge itself).
- **`multiSelect: true` validated**: Confirmed Claude Code CLI supports multiSelect — previous assumption was wrong. MEMORY.md corrected.

### Fixed — 3 UX dead-end flows

- **`/team-ship done`**: "No merged PR" error now shows 3-step actionable guide (check PR status → wait and retry → merge first). Was: one-liner with no next steps.
- **`/team-rc promote`**: Staging未配置 now shows explicit ⚠️ warning + manual verify hint. Was: silent skip that misleads user into thinking staging is OK.
- **`/team-ship` worktree mismatch**: Now recommends disabling worktree temporarily (primary) vs re-claiming (alternative). Was: two equal options with no recommendation.

### Fixed — tw-tips.txt stale reference

- `/team learn` tip: "11 节" → "13 节" (learn operation has 13 sections since v3.6.0)

### Fixed — Wrap mode checkbox gate

- `/team wrap`: After drive verification, force-checks all non-MANUAL sub-tasks in Contract before `sync-all-checkboxes`. Root cause: verify-only drive skips checking boxes → sync propagates unchecked state to Issue.

## [3.8.1] — 2026-03-11

### Fixed — Comprehensive audit: worktree chain + recovery hints + version source + naming

- **SKILL_VERSION source**: 5 dashboard/help/config locations changed from `tw-config.sh skill_version` read to hardcoded constant. Version now lives exclusively in skill frontmatter; `tw-bump-version.sh` keeps in sync.
- **`skill_version` removed from config**: Fresh install no longer writes it; merge mode actively deletes stale `skill_version` line; bump script no longer updates it.
- **Init worktree-root mapping**: Replaced ambiguous `{chosen_path}` placeholder with explicit table + per-option bash code blocks. Three options: `.claude/worktrees/` (in-repo), sibling (unset root), custom absolute path.
- **Auto/wrap worktree safety**: Both `/team auto` and `/team wrap` now detect worktree mode → create worktree during claim → STOP with explicit `cd` instructions instead of continuing drive/ship in main repo.
- **Recovery hints**: Added 13 missing `💡 /team doctor` hints across team.md (8), team-ship.md (4), team-claim.md (1). All error STOPs now have doctor guidance.
- **tw-git.sh tag**: Added error traps on `git tag` creation (was silent failure).
- **team-rc.md**: Added 3 missing recovery hints for abort/error STOPs.
- **Field naming**: Renamed "Acceptance Criteria" → "Success Criteria" in team.md auto contract template and team-ship.md extraction (canonical name per team-claim.md).
- **Manual version**: Updated `docs/teamwork-ai-manual.md` from 3.7.1 → 3.8.1, removed stale `skill_version` from config examples.
- **Learn lifecycle diagram**: Fixed to show v3 push model (no `backlog` state; Issue created directly as `wip`).
- **UX**: Added NEXT guidance to team-issue comment operation; fixed backtick formatting in team-drive tips.

### Fixed — Batch checkbox sync replaces fragile per-item sync

- **`sync-all-checkboxes`**: new `tw-contract.sh` subcommand — reads Contract, extracts ALL checked items (Sub-tasks + Success Criteria), fuzzy-matches against Issue body, updates in ONE API call
  - Fixes race condition (was: N reads + N writes, later writes could overwrite earlier)
  - Fixes silent failure (was: all errors `exit 0`, LLM thought sync succeeded)
  - Fixes exact-match fragility (was: `str.replace` needed char-for-char match; now: strips timestamps, normalizes whitespace/dashes, prefix/containment matching)
  - Fixes missing Success Criteria sync (was: only Sub-tasks synced, Success Criteria had no mechanism)
- **`sync-checkbox`**: deprecated (kept as backward-compatible alias with deprecation warning)
- **`/team-drive`**: Step 3 now calls `sync-all-checkboxes` (batch) instead of `sync-checkbox` (per-item)
- **Checkbox Gate**: adds `sync-all-checkboxes` call after recovery toggles — catches any missed syncs before completion
- **Parallel mode**: Lead syncs once per wave completion instead of per-task

### Fixed — Per-user setup (tab title, notifications) missing from fresh install

- **`/team init`**: Extracted per-user local setup (worktree preference, terminal tab title+color, desktop notifications) from MERGE MODE into new **Step 4h** — now runs for BOTH fresh install and merge mode
- **Root cause**: Tab title and notification prompts were only in the merge mode path; first user setting up the project (fresh install) was never prompted
- **`/team-ship`**: Added `sync-all-checkboxes` defense-in-depth call in pre-flight (Step 2b) — catches missed syncs from team-drive

### Added — Per-user config status dashboard + worktree path configuration

- **`/team init` Step 4h**: Always displays per-user config status dashboard (worktree, tab title, notifications) — even if all configured. Dynamically offers to enable unconfigured features via AskUserQuestion with "全部开启" one-click option
- **Worktree root path**: New `git config --local teamwork.worktree-root` setting. Three options: `.claude/worktrees/` (matches Claude Code native behavior), sibling directory (backward compat), or custom path
- **`/team-claim`**: Worktree creation now reads `teamwork.worktree-root` — creates worktrees under user's chosen root instead of hardcoded sibling path
- **Init ↔ Config separation**: Init shows status + guides setup; Config handles viewing/editing existing settings

## [3.8.0] — 2026-03-10

### Added — Desktop notification setup in /team init

- **Step 5c**: `/team init` now offers to configure macOS desktop notifications
  - Writes `~/.claude/notify.sh` hook (subagent-filtered, event-typed)
  - Registers Notification/Stop/TeammateIdle hooks in `~/.claude/settings.json`
  - Detects `terminal-notifier` for click-to-activate-terminal (osascript fallback)
  - Idempotent: skips if already configured

### Fixed — PR body self-referential criteria

- `/team-ship` Step 4b: filter out "PR merged"/"CI pass on PR" criteria when copying
  from Contract to PR body (self-referential in a PR, valid only at Issue level)

## [3.7.1] — 2026-03-10

### Fixed — Audit: VE full coverage + script bugs + version sync

- Fixed 4 script bugs found during full audit (tw-git, tw-pr, tw-config, tw-label)
- Aligned all VE output templates with visual-encoding-standard.md
- Synchronized version to 3.7.1 across all skill files, README, ONBOARDING, manual
- Doctor expanded to 12 checks (+PR CI status)

## [3.7.0] — 2026-03-09

### Added — Smart sync: dashboard behind-detection + AI conflict resolution

**Dashboard behind-base detection**:
- PR lines now show `⚠️ N behind {base}` when branch is behind base branch
- NEXT box prioritizes sync recommendation when PR is behind
- No extra API calls — uses local git rev-list against already-fetched branches

**AI conflict resolution** (`/team-ship sync`):
- Sync no longer STOPs on conflict — AI reads both sides of conflict markers and resolves
- For each conflicted file: read → understand intent → edit → git add → git rebase --continue
- Ambiguous conflicts (semantic contradictions) are resolved with best judgment + flagged in summary
- Auto force-push rebased branch after successful rebase
- Summary output shows resolved files and any flagged ambiguities

## [3.6.1] — 2026-03-09

### Added — Terminal tab title + Tips system

**Terminal tab title** (`/team init`):
- New optional step after worktree preference: configure terminal tab to auto-display
  current git branch name (e.g., `mission/42-add-auth`)
- User sees clear description of the benefit before opting in
- Auto-detects shell (zsh/bash/fish) and appends appropriate hook
- Checks for existing precmd/PROMPT_COMMAND conflicts before appending
- One-time, per-user, idempotent (marker-based detection)

**Tips system** (all 6 skills):
- Every `/team-*` command now displays a random tip after its output
- Tips pool in `scripts/tw-tips.txt` — easy to add/remove tips
- 15 tips covering lesser-known features: auto, wrap, doctor, batch, review, etc.
- Best-effort no-repeat within a session

## [3.6.0] — 2026-03-09

### Added — Manual Ops Handoff: guided flow for human-action tasks

When a mission contains tasks that require human action (console operations, server access,
deployment verification), teamwork now classifies, skips, and hands off manual tasks across the full pipeline:

**Full audit (5-domain parallel review) found and fixed 8 bugs in initial v3.6.0 implementation:**

**team-drive:**
- **Plan Assessment**: new `🔧 MANUAL` classification for sub-tasks requiring human action
- **Checkbox Gate**: MANUAL tasks are expected-unchecked, not treated as missed toggles
- **Execution Loop**: MANUAL tasks skipped during execution (not attempted by AI)
- **Mode Selection**: only code tasks counted for sequential/parallel/team threshold
- **Mission Briefing**: "all complete" check excludes MANUAL tasks (prevents re-drive loop)
- **Completion Summary**: splits into "CODE COMPLETE" + "YOUR TURN" sections
  - Code tasks shown as done
  - Manual tasks surfaced in a prominent action block with concrete guidance
  - AskUserQuestion gates the transition: user confirms manual ops done, defers to post-PR,
    asks for help, or skips — THEN `/team-ship` is suggested

**team-ship:**
- **Sub-task gate (2b)**: only code tasks block shipping; MANUAL tasks pass through
- **PR description**: adds "Manual Steps (post-merge)" checklist when manual tasks exist
- **Help text**: updated SHIP FLOW to mention manual task pass-through

**team-claim:**
- **Contract generation**: pre-tags `🔧 MANUAL` tasks during sub-task extraction (human-action tasks identified at claim time, not just during drive validation)

**team.md:**
- **`/team auto`**: drive phase counts code tasks only, pre-flight allows unchecked MANUAL, output shows code/manual split
- **`/team learn`**: new section 12 "MANUAL TASKS & HUMAN WORKFLOWS"
- **Visual Encoding**: added `🔧` to symbol table

**Bugs found and fixed by 5-domain audit:**
1. CRITICAL: Mission Briefing STOP logic — all code done + MANUAL exists → now enters Manual Ops Handoff instead of bare STOP
2. CRITICAL: Zero code tasks (all MANUAL) — now explicitly skips execution, goes to handoff
3. HIGH: `/team auto` task counting included MANUAL tasks — now code-only
4. HIGH: `/team learn` had no MANUAL task coverage — added section 12
5. HIGH: `🔧` emoji used but not in Visual Encoding table — added
6. HIGH: team-claim didn't tag MANUAL during Contract gen — added pre-tagging
7. MEDIUM: team-drive help text didn't explain `🔧` symbol — added TASK SYMBOLS block
8. MEDIUM: team-ship help text didn't mention manual pass-through — updated SHIP FLOW

Previously: manual tasks were mixed into the same list as code tasks with only a `○` marker,
`/team-ship` was suggested immediately with no guidance, and team-ship would STOP on unchecked manual tasks.

## [3.5.2] — 2026-03-09

### Fixed — Worktree guard: prevent drive/ship from wrong directory

- **team-drive**: Step 0c worktree guard — if worktree mode is enabled but user is in the main repo
  (no `.mission` file), STOP with clear message showing the correct worktree path.
  Prevents silent branch switching that bypasses worktree isolation.
- **team-ship**: Step 1b worktree guard — same check before shipping.
  Prevents shipping wrong code when user forgot to cd to worktree.
- **team doctor**: Check 5 item 6 — main-repo mismatch detection.
  Flags when user is in main repo with worktree mode enabled and active Contracts exist.
- **team doctor fix**: added main-repo mismatch action (show path / disable worktree mode / skip).

## [3.5.1] — 2026-03-09

### Fixed — Design review: logic gaps + docs alignment

Full 5-scout audit of teamwork design. Scripts: 100% healthy. Skill consistency: verified.

**Logic fixes:**
- **team-ship**: base branch guard — `protect-check` before shipping prevents accidental ship from main/develop
- **team-ship**: CI failure preserves Contract — no longer deletes Contract when CI fails,
  allowing `/team-drive` → `/team-ship` iteration loop
- **team-drive**: Contract-missing fallback — detects mission branch + open Issue and suggests
  `/team-claim #N` recovery path instead of bare "STOP"

**Documentation alignment:**
- **teamwork-ai-manual.md**: version 3.4.0 → 3.5.0, date updated
- **teamwork-v2-architecture.md**: marked as archived design document with pointer to current manual
- **team-doctor-design.md**: added Check 8 (Remote Sync), updated "7 checks" → "8 checks"
- **README.md**: added `/team doctor` row to skills table

## [3.5.0] — 2026-03-09

### Added — `/team config` + dashboard version + documentation completeness

- **`/team config`**: new subcommand — shows current config.yml content, worktree status, edit instructions.
  Gives users a clear entry point to view/modify settings without re-running init.
- **Dashboard version**: both leader and member dashboards now show `Teamwork v{version}` in footer.
  Version is read dynamically from `skill_version` in config (not hardcoded).
- **Help updated**: added `/team config`, `/team init`, worktree section, removed hardcoded version.
  Version now read from config via `tw-config.sh`.
- **Learn updated**: added Section 11 "Config Management" — shared vs local config, edit workflow.
  Version also dynamic.
- **Init merge mode summary**: now shows actual `skill_version` (was hardcoded "3.4.0"),
  includes `/team config` tip.
- **Argument table**: added `config` entry in both routing table and header table.
- **README/ONBOARDING/Manual**: added `/team config` references across all docs.

## [3.4.3] — 2026-03-09

### Fixed — Merge mode logic gaps + stale references

Full audit of worktree integration across all skill files.

- **Init merge mode**: now asks worktree preference (Step 5) when user has never configured it.
  Previously merge mode skipped ALL interactive questions, so worktree could never be enabled via re-init.
- **Init merge mode**: fixed stale version "3.4.0" in summary template → now shows current version.
- **Init merge mode**: quality section uses safe defaults (`true`) instead of referencing
  undefined Step 3 variables (merge mode skips interactive Steps 2-3).
- **Init merge mode**: Step 3 description clarified — uses existing config + safe defaults, not Step 3 answers.
- **README.md**: updated worktree config reference from `worktree.enabled` to `git config --local`.

## [3.4.2] — 2026-03-09

### Fixed — Worktree section LLM execution reliability

Dashboard worktree section was being skipped by LLM (conditional template buried in long output).
Restructured as mandatory pre-render check with explicit commands, placed before dashboard template.
Deduplicated worktree template (single definition referenced by both leader/member dashboards).

## [3.4.1] — 2026-03-09

### Changed — Worktree mode: shared → local

Worktree mode moved from shared `config.yml` to per-user `git config --local teamwork.worktree true`.
Each team member independently chooses their isolation strategy. Not committed to git.

- **Dashboard**: always shows worktree map when enabled (not just when inside a worktree)
  - ★ marks current location, MC association for each worktree
  - ⚠️ context mismatch warning when current branch doesn't match active mission
- **Init merge mode**: auto-migrates old `worktree:` section to `git config --local`
- **Doctor Check 5**: detects idle worktrees (no mission, no issue), context mismatch
- **Doctor Fix**: idle worktree gets "Remove / Remove + delete branch / Skip" options
- **ship done**: added comment clarifying `ensure-base` already pulls latest

## [3.4.0] — 2026-03-08

### Added — Doctor Actionability Upgrade

Three improvements to make `/team doctor` a workflow, not just a report.

**`/team doctor fix`** — Interactive fix mode:
- 🟢 Cleanup actions (orphan branches, prunable worktrees) auto-execute with summary
- 🟡 Warning actions (stash drop, stale branches) confirm each via AskUserQuestion
- 🔴 Critical actions confirm each with risk explanation
- Cross-contamination always presents both options (discard vs preserve)

**`/team doctor help`** — Self-contained usage guide:
- Shows all checks, severity levels, typical workflow
- No need to navigate back to `/team help` to understand doctor

**vs Remote column** — Branch remote tracking:
- New column in branch table shows ahead/behind vs `origin/{branch}`
- Highlights when local branch is behind remote (most actionable for current branch)
- New Check 8: Remote Sync diagnostic with pull/sync recommendations

**Action Plan improvements:**
- Actions now suggest teamwork commands where applicable (`/team-ship sync`, `/team-issue`, `/team-ship done`)
- Report footer reminds user to run `/team doctor fix`

### Added — Cross-Skill UX Audit (举一反三 from doctor bugs)

Systematic audit found the same 3 structural defects across all skills. Fixed:

**Subcommand help routes (5 new):**
- `/team-ship done help`, `/team-ship review help`, `/team-ship sync help`
- `/team-rc promote help`
- `/team doctor help` (above)

**Branch readiness warnings (2 skills):**
- `team-drive` Step 2a: shows vs base + vs remote before starting work
- `team-ship` Step 2a-bis: sync check before shipping (hard stop if >30 behind base)
- `team-ship` Step 4a: shows existing PR status (CI, review, mergeable) before push

**Error messages → actionable guidance (5 skills, ~12 fixes):**
- `team-ship`: test failures → suggest `/team-drive`; push failures → suggest `/team-ship sync`; CI failures → explain push-to-update flow
- `team-drive`: branch not found → two options (restore from remote / re-claim)
- `team-claim`: active mission → explain abandon + re-claim; not in roster → `/team init`; branch exists but issue closed → `/team doctor fix`
- `team-issue`: milestone not found → `gh issue edit` + `/team init`; labels missing → `/team init`
- `team-rc`: CI abort → suggest fix + re-run; staging fail → suggest deploy first

## [2.7.1] — 2026-03-04

### Fixed — Script Testing Audit (21 bugs found, 14 fixed)

Full-flow testing of 4 scripts (32 subcommands) + 6 skill files. 4 parallel test agents + cross-reference audit.

**CRITICAL fixes:**
- **tw-git.sh cmd_tag**: `git rev-parse "$version"` matched branches (not just tags) → false "already exists" skip. Fixed: `refs/tags/$version`
- **tw-git.sh cmd_merge_to**: `checkout && pull` chain — checkout failure didn't trigger `set -e`, merge proceeded on wrong branch. Fixed: explicit `|| exit 2`

**HIGH fixes:**
- **tw-git.sh cmd_ensure_base**: silent success on checkout failure (same `&&` chain bug). Fixed: explicit error handling
- **tw-git.sh cmd_create_branch**: sed metachar injection (`&`, `/`, `\` in slug). Fixed: bash `${//}` parameter expansion
- **tw-git.sh cmd_rebase/merge_to**: `git fetch 2>/dev/null` with `set -e` killed script silently on network error. Fixed: `|| WARNING`
- **tw-git.sh cmd_tag**: grep regex metacharacters in version (`.` = any char). Fixed: `grep -qF` (fixed string)
- **tw-pr.sh cmd_commit_type**: substring matching false positives ("prefix"→fix, "specification"→test). Fixed: `grep -qEw` (word boundary)
- **tw-contract.sh cmd_toggle_task**: reported success even when task N didn't exist. Fixed: `found` flag + exit 1

**MEDIUM fixes:**
- **tw-git.sh cmd_tag**: verification grep missing `$` anchor → false positive
- **tw-git.sh cmd_push**: `2>&1` mixed stderr into stdout → polluted output
- **tw-git.sh cmd_protect_check**: empty output on detached HEAD
- **tw-label.sh cmd_transition**: "done" transition dropped FROM label when from ≠ wip/review
- **tw-pr.sh cmd_exists**: BASE parameter accepted but never passed to `gh pr list`
- **tw-pr.sh cmd_create**: error message claimed BODY required (actually optional); `2>&1` contaminated URL extraction

## [2.7.0] — 2026-03-04

### Changed — Script Extraction (Von Neumann Mechanism-Content Separation)

Extracted ~64% of deterministic operations from LLM-interpreted skill prose into reusable bash scripts. Skills now orchestrate + judge; scripts execute deterministic operations.

**4 new scripts** in `scripts/`:
- **tw-git.sh** — 13 subcommands: ensure-base, create-branch, current, protect-check, push, fetch, rebase, merge-to, tag, commit, worktree-add, worktree-remove, log-since
- **tw-label.sh** — 4 subcommands: transition, set, verify, pr-label (config-aware status label state machine)
- **tw-contract.sh** — 8 subcommands: teamwork-dir, find, read-field, hash, check-freshness, toggle-task, sync-checkbox, delete
- **tw-pr.sh** — 7 subcommands: exists, commit-type, create, comment, watch, add-reviewer, verify-merged

**6 skill files refactored** to call scripts instead of inline bash:
- **team-claim.md**: branch creation → tw-git.sh, label transition → tw-label.sh
- **team-drive.md**: freshness check → tw-contract.sh, protect check → tw-git.sh, checkbox toggle → tw-contract.sh, log review → tw-git.sh
- **team-ship.md**: push → tw-git.sh, PR lifecycle → tw-pr.sh, labels → tw-label.sh, contract cleanup → tw-contract.sh, rebase → tw-git.sh
- **team-release.md**: merge → tw-git.sh, tag → tw-git.sh, commit → tw-git.sh
- **team-issue.md**: version bump only
- **team.md**: freshness check → tw-contract.sh, version references updated

**Design principles**:
- Scripts accept LLM-generated content as parameters (pure I/O boundary)
- Scripts self-discover config via SCRIPT_DIR + tw-config.sh
- Subcommand pattern: `tw-{domain}.sh {subcommand} [args...]`
- All scripts macOS compatible, meaningful exit codes, stderr for errors

## [2.6.1] — 2026-03-04

### Fixed — Falsification Audit (Popper's method applied across all skill files)

**CRITICAL fixes:**
- **team.md + team-drive.md**: Freshness check completely broken — `ISSUE_TITLE`/`ISSUE_BODY` never extracted from JSON, `CONTRACT_PATH` never defined. Hash always computed on empty string → false "ISSUE UPDATED" warning on every `/team-drive`. Fixed: added `jq` extraction + proper loop over Contract files.
- **team.md**: Merge mode used `import yaml` (PyYAML) — not in Python stdlib. Fails on systems without PyYAML. Fixed: replaced with `grep -q "^${section}:"` (zero dependencies).
- **team.md**: Branch protection `gh api --field` sent JSON objects as strings → 422 Validation Error. Fixed: switched to `--input` with heredoc for proper JSON body.
- **team-release.md**: Shell injection in CHANGELOG Python script — release body with triple quotes breaks embedded `'''$ENTRY'''`. Fixed: write to temp file, pass as `sys.argv`.
- **team-release.md**: `git merge "$BASE_BRANCH"` used stale local copy without fetch → release could miss latest code. Fixed: `git fetch` before merge, use `origin/$BASE_BRANCH`.
- **team-release.md**: Version bump committed on `$RELEASE_BRANCH` but pushed to `$BASE_BRANCH` → wrong branch in dual-branch mode. Fixed: `git checkout "$BASE_BRANCH"` before config update commit.

**HIGH fixes:**
- **team.md**: `init` force mode promised in routing table but never implemented — `/team init` on existing project showed dashboard. Fixed: added `$ARGUMENTS == "init"` check in Step 1.
- **team.md**: Version display in merge summary showed `2.5.0` (stale from v2.5). Fixed → `2.6.1`.
- **team.md**: Issue template hardcoded `"status:queued"`. Fixed → `"{STATUS_PREFIX}queued"`.
- **team.md**: Git hooks `grep` matched commented-out config lines + no `head -1`. Fixed: `grep -v '^\s*#'` + `head -1`.
- **team.md**: `{REPO}` template syntax mixed with `$REPO` shell syntax in bash code. Fixed → all `$REPO`.
- **team.md**: Duplicate section ID `6c` (two sections). Fixed → second renamed `6d`.
- **team.md**: Fallback label creation hardcoded `"mission"` and `"status:*"`. Fixed → config-driven.
- **team-issue.md**: macOS BSD sed `\+` syntax fails — `sed 's/.*--milestone \+\([^ ]*\).*/\1/'`. Fixed → `\+` replaced with `  *` (portable).
- **team-issue.md**: `MILESTONE_FLAG="--milestone \"$MILESTONE\""` embeds literal escaped quotes → milestone assignment always fails. Fixed → `${MILESTONE:+--milestone "$MILESTONE"}` parameter expansion.
- **team-issue.md**: References nonexistent `team.repo` config key. Fixed → proper `mc_label` + `STATUS_PREFIX` reads.
- **team-ship.md**: All label operations hardcoded `status:` prefix. Fixed → reads `$STATUS_PREFIX` from config.
- **team-ship.md**: `$ISSUE_TITLE` undefined in PR creation → commit type detection always defaulted to `feat`. Fixed → extract from Contract frontmatter.
- **team-ship.md**: `gh pr list --head {branch}` without `--state open` could find closed PRs. Fixed → added `--state open`.
- **team-claim.md**: Step 6 hardcoded `"status:queued"` / `"status:wip"`. Fixed → config-driven `${STATUS_PREFIX}`.
- **team-release.md**: No check if tag already exists before `git tag`. Fixed → `git rev-parse` + `git ls-remote` guards.
- **team-release.md**: `repos/{REPO}/milestones` template syntax in bash code. Fixed → `repos/$REPO/milestones`.
- **setup-github-labels.sh**: Hardcoded `"mission"` label, ignoring `mc_label` config. Fixed → reads via `tw-config.sh` with fallback.
- **setup-github-labels.sh**: `grep '^\s*status:'` could match wrong config line. Fixed → prefers `tw-config.sh`, fallback grep excludes comments.

### Changed
- All skill files bumped to v2.6.1
- Post-merge Action regex (committed in previous session fix): expanded from `Close(s)` to all 9 GitHub closing keywords, added `parseInt` for proper number type, replaced milestone race condition workaround with `listForRepo` query

## [2.6.0] — 2026-03-04

### Added — Dual-branch model + release safety gates

- **Dual-branch support** — New `conventions.release_branch` config key enables base/release branch separation:
  - When `release_branch == base_branch` (default) → single-branch mode, backward compatible with v2.5
  - When different (e.g. base=develop, release=main) → dual-branch mode: `/team-release` merges base→release before tagging
- **`team-release.md`** v2.5.0 → v2.6.0:
  - **Staging deploy gate** (Step 4b): if `deploy.staging_workflow` configured, checks latest Action run status before tagging. Warns on failure, user can abort or proceed.
  - **Merge base→release** (Step 5): in dual-branch mode, merges base into release branch before creating tag. Stops on conflict.
  - **CHANGELOG append** (Step 7b): after `gh release create`, fetches release body and prepends to `CHANGELOG.md`. Non-fatal if fails.
  - Renumbered steps 5-11 → 6-12 to accommodate new steps
- **`team-drive.md`** v2.4.0 → v2.6.0:
  - **Branch safety check** (Step 2b): prevents code execution on protected branches (base_branch or release_branch). Must be on feature branch.
- **`team-claim.md`** v2.4.1 → v2.6.0:
  - Explicitly checks out `base_branch` and pulls before creating feature branch (was implicit before)
- **`team-issue.md`** v2.5.1 → v2.6.0:
  - `--milestone` parameter: override milestone assignment per-issue (e.g. `/team-issue --milestone V0.2 fix camera`)
  - Priority: `--milestone` arg > `versions.current` from config > none

### Changed — Config-driven branch references

- **`team-ship.md`** v2.4.0 → v2.6.0:
  - Sync operation, done operation, help text all use `$BASE_BRANCH` from config instead of hardcoded "main"
  - **Done operation D3/D4 → fallback pattern**: checks if post-merge Action already handled label updates and Issue closure before acting. Skips with `(handled by PR merge)` / `(handled by post-merge Action)` messages.
- **`team.md`** v2.5.2 → v2.6.0:
  - **New Step 4e: Post-merge Action** — generates `.github/workflows/post-merge.yml` that auto-runs on PR merge: extracts `Closes #N` from PR body, updates labels (wip/review→done), checks all Issue checkboxes, checks milestone completion and comments when 100%
  - **New Step 5b: Repo merge settings** — auto-delete head branches, squash merge only, PR_BODY commit message (preserves `Closes #N`)
  - **Step 5c: Dual-branch protection** — protects both base_branch AND release_branch in dual-branch mode
  - **CI workflow templates** — use `{BASE_BRANCH}` placeholder instead of hardcoded `main`; dual-branch mode includes both branches in push triggers
  - Config health check table now includes `release_branch`, `deploy.staging_workflow`, and post-merge Action
  - Fresh install template includes commented `release_branch` and `deploy` sections
  - Dashboard merge count uses config `base_branch`
  - All `skill_version` references updated to 2.6.0

### Design principle

All branch names, workflow names, and deployment behaviors are config-driven. Zero hardcoded project references in skill files. Generic design validated through vi_agent practice.

## [2.5.1] — 2026-03-04

### Added — `/team-issue` duplicate check

- **`team-issue.md`** v2.4.1 → v2.5.1 — New **Step 1b: Duplicate Check** before Issue body generation:
  - Extracts 2-3 key terms from analyzed title/description
  - Runs `gh issue list --search` against open Issues
  - AI evaluates results for **semantic similarity** (not string match)
  - If duplicates found → `AskUserQuestion` with options: create anyway, update existing, or cancel
  - If no match → proceeds silently (zero friction on happy path)
- Help text updated: "WHAT HAPPENS (create)" now shows 6 steps including duplicate search

## [2.5.0] — 2026-03-04

### Added — `/team-release` skill

- **`team-release.md`** — New skill: close a version milestone, create git tag, publish GitHub Release, and prepare the next version. Full 11-step flow:
  1. Pre-flight: detect config dir, read `versions.current` (required), `versions.spec_path` (optional)
  2. Milestone check: show open/closed Issue counts, warn if open issues remain
  3. AI audit: if `versions.spec_path` configured, compare spec Success Criteria vs merged PRs + closed Issues. Outputs ✅/⚠️/❌ per criterion
  4. User confirmation before any mutations
  5. `git tag {VERSION}` + `git push origin {VERSION}`
  6. `gh release create {VERSION} --generate-notes` — auto-generates notes from merged PRs
  7. Close GitHub Milestone
  8. Prompt for next version name
  9. Create next GitHub Milestone
  10. Update `versions.current` in config.yml
  11. Commit config change + push
- Tag format: uses `versions.current` value exactly (e.g. `"V0.1"` → tag `"V0.1"`)
- AI audit is conditional: only runs when `versions.spec_path` is set and file exists

### Changed — `/team init` idempotency fix

- **`team.md`** v2.4.1 → v2.5.0 — `Step 4b` now has two modes:
  - **MERGE MODE** (config exists): preserves all existing sections, updates `skill_version` in-place via sed, appends only missing sections (`project`, `conventions`, `label_prefix`, `quality`). Explicitly does NOT touch `team`/`members`, `github`, `worktree`, `versions`, `roles`, or any custom sections.
  - **FRESH INSTALL** (no config): writes full template (existing behavior, unchanged).
  - Previous behavior: unconditional overwrite. Running `/team init` on a project with customized config (e.g. vi_agent's nested `github.mc_label` schema) would silently destroy all customizations.
- **`team.md` config template** — added `spec_path` comment to `versions:` block, pointing to `/team-release` as consumer.

---

## [2.4.1] — 2026-03-03

### Changed — Dashboard Upgrade (team.md v2.3.2 → v2.4.1)

- **`team.md`** — Replaced all fragile `grep/sed` config chains in Step 6c with `tw-config.sh` calls. Four sed chains eliminated:
  - `MISSION_LABEL` — now via `tw-config.sh github.mc_label` (dual-schema: teamspace first, then flat `mc_label`, then `"mission"` default)
  - `STATUS_PREFIX` — now via `tw-config.sh label_prefix.status`
  - `PRIORITY_PREFIX` — now via `tw-config.sh label_prefix.priority`
  - `CURRENT_VERSION` — now via `tw-config.sh versions.current`
- **`team.md` config template** — `skill_version` in generated `.teamspace/config.yml` template updated from `2.1.0` to `2.4.1`. New installs get correct version.
- **`team.md` version** — Bumped from `2.3.2` to `2.4.1` (skips 2.4.0 — dashboard was deferred during stability hardening, now catches up).

### Fixed — E2E Validation on vi_agent

- **`scripts/tw-config.sh`** — Fixed inline YAML comment stripping. Values like `mc_label: "mission-contract"  # comment` were returning `mission-contract"  # comment` instead of `mission-contract`. Now correctly extracts quoted value, ignoring trailing inline comment.
- **`team-claim.md`** — mc_label resolution now handles both config schemas: tries `github.mc_label` (teamspace/vi_agent schema) first, then `mc_label` (teamwork v2 schema), then defaults to `"mission"`. Fixes silent fallback to wrong label on vi_agent.
- **`team-issue.md`** — Same mc_label dual-schema fix. Both `github.mc_label` and `mc_label` paths tried before falling back to default.

Discovered and fixed during E2E validation run on `../vi_agent`. Full skill flow (issue → claim → drive → ship) verified non-destructively (Issue #27, PR #28, all closed/deleted post-test).

---

## [2.4.0] — 2026-03-03

### Changed — Stability Hardening

- **`scripts/tw-config.sh`** — New unified config parser. Replaces 15+ fragile sed chains across all skill files. Handles colon-containing values (label prefixes like `status:`), quoted strings, and nested YAML. Usage: `bash ~/.claude/commands/scripts/tw-config.sh project.test_command ""`.
- **Placeholder normalization** — Single canonical naming convention across all 5 skills: `{issue}`, `{pr}`, `{user}`, `{branch}`, `{title}`, `{repo}`, `{milestone}`. Eliminates confusion between `{N}`, `{ISSUE_NUMBER}`, `{pr-number}`.
- **team-ship pseudocode implemented** — tech-lead extraction (reads team roster for `role: tech-lead`), commit-type detection (prefix-based: fix/refactor/test/docs/feat), and base_branch reading were pseudocode in v2.3.x. Now real bash.
- **Critical command verifications** — `gh api user` failure now fatal with clear message in all skills. `git checkout` verifies branch existence before switching. TEST_CMD empty check warns instead of silently passing.
- **Issue bidirectional sync in team-drive** — After each sub-task is checked off in the Contract, the corresponding checkbox in the GitHub Issue body is updated. Team members see real-time progress on GitHub without visiting the branch. Non-fatal: network errors warn and continue.
- **Version from Milestone (not config)** — team-claim reads Milestone from Issue JSON (`milestone.title`) instead of `config.yml versions.current`. Config fallback still works for backward compatibility. Removes need to keep `versions.current` in sync manually.
- **REPO detection via `gh repo view`** — team-issue no longer requires a `repo:` field in config.yml. Uses `gh repo view --json nameWithOwner` auto-detection instead (same approach as other skills).

---

## [2.3.3] — 2026-03-03

### Added — Script-Based Skill Architecture (Pilot)

- **`scripts/setup-github-labels.sh`** — Executable bash script that creates/updates all teamwork labels on GitHub. Reads label prefixes from `.teamwork/config.yml` automatically. Idempotent (uses `--force`). Run: `bash ~/.claude/commands/scripts/setup-github-labels.sh`
- **`install.sh` scripts support** — `install_skill()` now detects and installs `scripts/*.sh` files from each skill directory to `$TARGET_DIR/scripts/`. Symlink mode: symlinks auto-update with `git pull`. Copy mode: copies standalone.
- **`/team` label setup uses script** — Step 5a now calls the script first with inline fallback. Eliminates the fragile 10-command sequence the LLM had to interpret.
- Establishes the pattern: skill folder = `.md` (AI brain) + `scripts/` (deterministic bash). LLM runs scripts; scripts do the mechanical work.

---

## [2.3.2] — 2026-03-03

### Added — Version Management + Issue Format Standardization

- **Milestone-based version tracking** — If `versions.current` is set in config, `/team-issue` assigns the Issue to that GitHub Milestone. Dashboard now queries milestone progress (open/closed counts) and displays `Version X: N/M tasks done`.
- **`/team` dashboard version header** — Reads `versions.current` from config, queries GitHub Milestones API for progress, shows `Version {name}: {closed}/{total} tasks done ({pct}%)` at top.
- **Canonical mission-contract format** — `/team-issue` now generates a consistent structure: `### Objective` (new top section), `### Success Criteria` (canonical name), `### Context & References` (not "Context"). Field names documented as canonical so parsers don't drift.
- **`/team-claim` backward-compatible parsing** — Reads `Success Criteria` (canonical) and `Acceptance Criteria` (legacy), `Context & References` and `Context` (legacy). Populates `version:` Contract frontmatter from `versions.current` config.
- `team-issue`, `team-claim`, `team` bumped to v2.3.2. `team-ship` stays v2.3.1.

---

## [2.3.1] — 2026-03-03

### Fixed — Ship Status Lifecycle

- **`/team-ship` label verification** — After updating Issue label from `status:wip` to `status:review`, immediately verifies via `gh issue view --json labels`. Shows actual labels in output instead of assumed ✅. If label update failed, shows diagnostic warning with manual fix command.
- **`/team` dashboard REVIEW section** — Added `UNDER REVIEW (PR open, awaiting merge)` section between member status and QUEUED. Ships `status:review` issues are now visible on the dashboard instead of falling into a grey zone.
- Both skills bumped to v2.3.1.

---

## [2.3.0] — 2026-03-03

### Added — Issue Sync & Lifecycle Comments

- **Issue freshness detection** — Contract stores `issue_content_hash` (SHA256 of title+body at claim time). `/team-drive` checks hash on entry; if Issue content changed, shows current body and asks user to Update Contract / Continue / Abort. Uses content hash instead of `updatedAt` to avoid false positives from our own comments, label changes, and assignee changes.
- **Dashboard [UPDATED] badge** — `/team` dashboard marks claimed Issues with `[UPDATED]` if content hash differs from Contract.
- **Lifecycle comments** — 3 key events posted as Issue comments for team visibility:
  - `/team-claim`: "Claimed by @user" (already existed)
  - `/team-drive` completion: "All sub-tasks complete — ready for review"
  - `/team-ship` PR creation: "PR #N created — {url}"
- All freshness checks and comments are **non-fatal** — network failures warn but don't block execution.
- Backward compatible: pre-v2.3.0 Contracts without `issue_content_hash` skip freshness check silently.

### Changed

- All skills bumped to v2.3.0

---

## [2.2.0] — 2026-03-03

### Added

- **`/team-issue`** — New skill: turn natural language into standardized GitHub Issues (mission-contract format). AI-enriches description with title, priority, size, domain, sub-tasks, and relevant files from codebase scan. Preview before publish.

### Changed

- **`/team-claim`** — Removed inline `create "title"` subcommand. Issue creation now handled by dedicated `/team-issue` skill (single responsibility).
- All skills bumped to v2.2.0

---

## [2.1.1] — 2026-03-02

### Fixed — Architect Review (13 bugs)

- **CRITICAL**: Contract inaccessible from worktree (active/ gitignored → added copy-to-worktree step)
- **CRITICAL**: Issue template hardcoded "mission" label → now uses `{MISSION_LABEL}` from config
- **CRITICAL**: Branch pattern hardcoded `mission/...` → reads `conventions.branch_pattern` from config
- **HIGH**: `team:` key format collision (erwin=member list, vi-agent=metadata) → check `members:` first
- **HIGH**: Missing config health check → added Step 1c with non-blocking warning
- **HIGH**: Operation Done assumes `mission/` branch prefix → generic numeric extraction
- **MEDIUM**: Worktree `enabled` flag missing → section presence = enabled
- **MEDIUM**: sed greedy pattern `'s/.*'` → all normalized to `'s/^[^:]*://'`
- **MEDIUM**: `git add -A` in team-drive → changed to specific file staging
- **FIX**: `cp` command used `{N}` instead of `${ISSUE_NUMBER}`
- **FIX**: Step 1b→1c routing contradiction (found users bypassed health check)
- **FIX**: `$BRANCH_PATTERN` extraction command missing

---

## [2.1.0] — 2026-03-02

### Added — vi-agent Feature Absorption

Integrated 6 advantages from vi-agent's `.teamspace` system while keeping erwin's portability.

- **Prefixed labels** — `status:wip`, `priority:P1`, `status:blocked`, `status:queued`, `status:done` (queryable, no collisions)
- **Claim comment visibility** — `gh issue comment` posted when claiming an Issue for team transparency
- **Worktree isolation** (optional) — `worktree.enabled: true` in config; creates isolated working directories per mission
- **Version management** (optional) — `versions.current: "V1.0"` in config; version labels + dashboard progress tracking
- **`/team-ship done`** — Post-merge cleanup: close Issue, update labels to `status:done`, clean worktree, return to main
- **`/team-ship review`** — AI code review: fetch PR diff, analyze for correctness/security/architecture/quality/performance, publish `gh pr review`
- **`/team-ship sync`** — Rebase current branch on latest main
- **`.mission` file** — Worktree mode detection (Issue number stored in worktree root)

### Changed

- Labels upgraded from flat (`in-progress`, `P0`) to prefixed (`status:wip`, `priority:P0`)
- Issue template now auto-applies `status:queued` label
- Config schema extended with optional `label_prefix`, `worktree`, `versions` sections
- `team-ship` now has argument routing: `done`, `review`, `sync` subcommands
- Backward compatible: unprefixed labels still work for legacy configs

---

## [2.0.0] — 2026-03-02

### Changed — GitHub-first Redesign

Complete rewrite. GitHub Issues/PRs/Actions replace local file-based coordination.

**Breaking changes:**
- `/teamwork` → renamed to `/team`
- `/team-get-mission` → renamed to `/team-claim`
- `/team-complete` → renamed to `/team-ship`
- `_teamwork/` directory → replaced by `.teamwork/` + GitHub ecosystem
- Mission contracts now auto-generated from GitHub Issues (not manually written)
- Config format changed: `_teamwork/config.md` → `.teamwork/config.yml` (YAML)

### Added

- **`/team`** — Init system with project auto-detection + GitHub config (labels, CI, hooks, Issue template) + dashboard from live GitHub data
- **`/team-claim`** — Claim GitHub Issue → auto-generate AI-enriched Mission Contract → create branch → assign on GitHub
- **`/team-drive`** — Read Contract, execute sub-tasks with verify loop, update checkboxes, commit with Mission reference
- **`/team-ship`** — Pre-flight checks → push → create PR with "Closes #N" → watch CI → request review → cleanup Contract
- **GitHub Actions CI** — auto-generated workflow per language (Node, Python, Go, Rust, Java)
- **Git hooks** — pre-commit (lint) + pre-push (test) via `.githooks/`
- **GitHub Issue template** — structured mission template (`mission.yml`)
- **Priority labels** — P0/P1/P2/P3 + mission/in-progress/review
- **Branch protection** — optional, requires CI pass + PR review
- **AI-enriched Context Files** — Contract auto-discovers relevant files via Glob/Grep
- **Dashboard** — live team status from GitHub API (assigned Issues, open PRs, merge count)

### Removed

- `_teamwork/` directory structure (replaced by `.teamwork/` + GitHub)
- Manual mission contract authoring (now auto-generated from Issues)
- Per-person status files (replaced by GitHub Issue assignment)
- Linear merge protocol (replaced by PR-based workflow)

### Architecture

- Design doc: `docs/teamwork-v2-architecture.md`
- ADR: `_project/decisions/0010-teamwork-v2-github-first.md`
- Inspired by vi-agent lifecycle patterns + GitHub-native workflows

---

## [1.0.0] — 2026-03-02

### Added

- **`/teamwork`** — Init system with interactive team config + dashboard view
- **`/team-get-mission`** — Claim missions, auto-create branches, update status
- **`/team-drive`** — Execute missions wrapping /drive methodology with teamwork context
- **`/team-complete`** — Push, rebase, conflict resolution, merge, status close
- **`_teamwork/` directory structure** — config.md, missions/, members/, archive/
- **Per-person status files** — zero metadata merge conflicts by design
- **Identity detection** — automatic via `gh auth status`
- **Branch auto-creation** — `mission/{id}-{slug}-{username}` naming convention
- **Linear merge protocol** — first merger wins, last merger resolves all conflicts
