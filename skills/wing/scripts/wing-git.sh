#!/usr/bin/env bash
# wing-git.sh — Git operations for wing skills
#
# USAGE:
#   bash wing-git.sh ensure-base                        # Checkout base_branch, pull latest
#   bash wing-git.sh create-branch ISSUE SLUG USER      # Create branch from pattern
#   bash wing-git.sh build-branch-name ISSUE SLUG USER  # Build branch name (no checkout)
#   bash wing-git.sh extract-issue-number [BRANCH]      # Extract issue number from branch name
#   bash wing-git.sh current                             # Print current branch name
#   bash wing-git.sh protect-check                       # Exit 3 if on protected branch
#   bash wing-git.sh push [BRANCH]                       # Push with -u (default: current branch)
#   bash wing-git.sh rebase TARGET                        # Fetch + rebase on target
#   bash wing-git.sh tag VERSION [TARGET_REF]              # Check exists, create, push, verify
#   bash wing-git.sh commit MSG [FILES...]               # Stage files + commit (-A if no files)
#   bash wing-git.sh worktree-guard [BRANCH]              # Check worktree boundary (exit 3=mismatch)
#   bash wing-git.sh claim-setup ISSUE SLUG USER [TW_DIR [CONTRACT]]  # Atomic branch+worktree setup
#   bash wing-git.sh worktree-add BRANCH PATH            # Create worktree
#   bash wing-git.sh worktree-remove PATH                # Remove worktree safely
#   bash wing-git.sh worktree-main-repo                   # Print main worktree path
#   bash wing-git.sh worktree-find-by-branch BRANCH       # Find worktree path for a branch
#   bash wing-git.sh log-since [BASE]                     # Show commits + diff stat since base
#   bash wing-git.sh cut-release VERSION                  # Cut rc/VERSION from base (remote-only, no checkout)
#   bash wing-git.sh slugify TITLE                         # Slugify title for branch names
#   bash wing-git.sh find-rc                               # Find active RC branch (exit 1=none, 3=multiple)
#   bash wing-git.sh next-version PREFIX                   # Derive next patch version from tags
#   bash wing-git.sh list-merged-branches                  # Local mission/* branches with merged PRs
#   bash wing-git.sh milestone-resolve NAME                # Resolve short milestone name to full title
#
# EXIT CODES:
#   0 — success
#   1 — usage error or missing args
#   2 — git operation failed
#   3 — safety check failed (e.g., on protected branch)
#
# READS: .wing/config.yml (or .teamwork/.teamspace for backward compat) via wing-config.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TW_CONFIG="$SCRIPT_DIR/wing-config.sh"

# Read config values
_base_branch() {
  local b
  b=$(bash "$TW_CONFIG" conventions.base_branch "main" 2>/dev/null)
  echo "${b:-main}"
}

_production_branch() {
  local p
  p=$(bash "$TW_CONFIG" conventions.production_branch "main" 2>/dev/null)
  echo "${p:-main}"
}

_branch_pattern() {
  local p
  # Try conventions.branch_pattern first, then worktree.branch_pattern (single call with fallback)
  p=$(bash "$TW_CONFIG" conventions.branch_pattern "" 2>/dev/null)
  if [ -z "$p" ]; then
    p=$(bash "$TW_CONFIG" worktree.branch_pattern "mission/{issue}-{slug}-{user}" 2>/dev/null)
  fi
  # CRITICAL: Do NOT use ${p:-default} when default contains {} — bash closes
  # parameter expansion at first unmatched }, appending leftover text to output.
  if [ -z "$p" ]; then
    p="mission/{issue}-{slug}-{user}"
  fi
  echo "$p"
}

usage() {
  echo "Usage: wing-git.sh <subcommand> [args...]" >&2
  echo "Subcommands: ensure-base, create-branch, build-branch-name," >&2
  echo "             extract-issue-number, current, protect-check, push," >&2
  echo "             rebase, tag, commit, worktree-guard, claim-setup," >&2
  echo "             worktree-add, worktree-remove," >&2
  echo "             worktree-main-repo, worktree-find-by-branch, log-since," >&2
  echo "             cut-release, slugify, find-rc, next-version," >&2
  echo "             list-merged-branches, milestone-resolve" >&2
  exit 1
}

cmd_ensure_base() {
  local base
  base=$(_base_branch)
  git checkout "$base" 2>/dev/null || {
    echo "ERROR: Could not checkout $base" >&2
    exit 2
  }
  git pull origin "$base" 2>/dev/null || {
    echo "WARNING: Could not pull $base, working with local copy" >&2
  }
  echo "$base"
}

cmd_create_branch() {
  local issue="${1:-}" slug="${2:-}" user="${3:-}"
  if [ -z "$issue" ] || [ -z "$slug" ] || [ -z "$user" ]; then
    echo "ERROR: create-branch requires ISSUE, SLUG, USER" >&2
    exit 1
  fi

  local pattern branch
  pattern=$(_branch_pattern)
  # Use bash parameter expansion instead of sed to avoid metachar injection (&, /, \)
  branch="${pattern//\{issue\}/$issue}"
  branch="${branch//\{slug\}/$slug}"
  branch="${branch//\{user\}/$user}"
  branch="${branch//\{type\}/mission}"
  branch="${branch//\{task-id\}/$issue}"

  if git show-ref --verify --quiet "refs/heads/$branch" 2>/dev/null; then
    echo "Branch '$branch' already exists. Switching to it." >&2
    git checkout "$branch"
  else
    git checkout -b "$branch"
  fi
  echo "$branch"
}

cmd_current() {
  git branch --show-current
}

cmd_protect_check() {
  local current base prod
  current=$(git branch --show-current)
  if [ -z "$current" ]; then
    echo "OK: detached HEAD"
    return 0
  fi
  base=$(_base_branch)
  prod=$(_production_branch)

  # Block: base branch (develop), production branch (main), rc/* branches
  if [ "$current" = "$base" ] || [ "$current" = "$prod" ]; then
    echo "SAFETY: On protected branch '$current'" >&2
    exit 3
  fi
  # RC branches should only receive hotfix PRs, not direct commits
  if [[ "$current" == rc/* ]]; then
    echo "SAFETY: On RC branch '$current' — branch from it for hotfixes, don't commit directly" >&2
    exit 3
  fi
  echo "OK: on '$current'"
}

cmd_push() {
  local branch="${1:-}"
  [ -z "$branch" ] && branch=$(git branch --show-current)

  if git push -u origin "$branch"; then
    echo "Pushed: $branch"
  else
    echo "Push failed. Try: git pull --rebase origin $branch" >&2
    exit 2
  fi
}

cmd_rebase() {
  local target="${1:-}"
  [ -z "$target" ] && target=$(_base_branch)

  local rebase_ref="origin/$target"
  if ! git fetch origin "$target" 2>/dev/null; then
    echo "WARNING: Could not fetch $target, rebasing on local ref" >&2
    rebase_ref="$target"
  fi
  if git rebase "$rebase_ref"; then
    echo "Rebased on $rebase_ref"
  else
    echo "Rebase failed. Resolve conflicts, then: git rebase --continue" >&2
    exit 2
  fi
}

cmd_tag() {
  local version="${1:-}" target="${2:-}"
  if [ -z "$version" ]; then
    echo "ERROR: tag requires VERSION [TARGET_REF]" >&2
    exit 1
  fi

  # Check if tag already exists (use refs/tags/ prefix to avoid matching branches)
  if git rev-parse "refs/tags/$version" >/dev/null 2>&1; then
    echo "Tag '$version' already exists locally. Skipping." >&2
    return 0
  fi
  if git ls-remote --tags origin | grep -qF "refs/tags/${version}"; then
    echo "Tag '$version' already exists on remote. Skipping." >&2
    return 0
  fi

  # Remote-only: tag a specific ref without checkout.
  # If TARGET_REF provided (e.g. "origin/main"), tag that ref.
  # If omitted, tags HEAD (backward compat).
  if [ -n "$target" ]; then
    git tag "$version" "$target" || {
      echo "ERROR: Could not create tag '$version' at ref '$target'" >&2
      exit 2
    }
  else
    git tag "$version" || {
      echo "ERROR: Could not create tag '$version'" >&2
      exit 2
    }
  fi
  git push origin "$version" || {
    echo "ERROR: Could not push tag '$version' to origin" >&2
    exit 2
  }

  # Verify (use -F for fixed string match to avoid regex metachar issues with dots)
  if git ls-remote --tags origin | grep -qF "refs/tags/${version}"; then
    echo "Tag $version pushed and verified"
  else
    echo "ERROR: Tag $version push could not be verified" >&2
    exit 2
  fi
}

cmd_commit() {
  local msg="${1:-}"
  if [ -z "$msg" ]; then
    echo "ERROR: commit requires MSG" >&2
    exit 1
  fi
  shift

  if [ $# -eq 0 ]; then
    git add -u  # Only tracked files — never stage untracked files by default
  else
    git add "$@"
  fi

  git commit -m "$msg"
}

cmd_worktree_add() {
  local branch="${1:-}" path="${2:-}"
  if [ -z "$branch" ] || [ -z "$path" ]; then
    echo "ERROR: worktree-add requires BRANCH and PATH" >&2
    exit 1
  fi
  # Use existing branch if it exists (create-branch may have already created it),
  # otherwise create it with -b
  if git show-ref --verify --quiet "refs/heads/$branch" 2>/dev/null; then
    git worktree add "$path" "$branch"
  else
    git worktree add -b "$branch" "$path"
  fi
  echo "$path"
}

cmd_worktree_remove() {
  local path="${1:-}"
  if [ -z "$path" ]; then
    echo "ERROR: worktree-remove requires PATH" >&2
    exit 1
  fi
  if [ -d "$path" ]; then
    git worktree remove "$path" 2>/dev/null || {
      echo "WARNING: Could not remove worktree at $path (try --force)" >&2
      return 2
    }
    echo "Removed worktree: $path"
  else
    echo "WARNING: Worktree path not found: $path" >&2
    return 0
  fi
}

cmd_log_since() {
  local base="${1:-}"
  [ -z "$base" ] && base=$(_base_branch)

  echo "=== Commits since $base ==="
  local commits
  commits=$(git log --oneline "${base}..HEAD" 2>/dev/null) || true
  if [ -z "$commits" ]; then
    echo "(no commits)"
  else
    echo "$commits"
  fi
  echo ""
  echo "=== Files changed ==="
  git diff "${base}...HEAD" --stat 2>/dev/null || echo "(no changes)"
}

cmd_cut_release() {
  local version="${1:-}"
  if [ -z "$version" ]; then
    echo "ERROR: cut-release requires VERSION" >&2
    exit 1
  fi

  local base rc_branch
  base=$(_base_branch)
  rc_branch="rc/$version"

  # Remote-only: fetch latest base, create rc branch from it, push.
  # Does NOT checkout or modify the user's working tree.
  git fetch origin "$base" 2>/dev/null || {
    echo "ERROR: Could not fetch origin/$base" >&2
    exit 2
  }
  # Create branch pointing at origin/base without checking it out
  git branch "$rc_branch" "origin/$base" || {
    echo "ERROR: Could not create branch $rc_branch (already exists?)" >&2
    exit 2
  }
  git push -u origin "$rc_branch" || {
    echo "ERROR: Could not push $rc_branch" >&2
    # Clean up local branch on push failure
    git branch -d "$rc_branch" 2>/dev/null || true
    exit 2
  }
  echo "$rc_branch"
}

cmd_slugify() {
  local title="${1:-}"
  if [ -z "$title" ]; then
    echo "ERROR: slugify requires TITLE" >&2
    exit 1
  fi
  local slug
  slug=$(echo "$title" | tr '[:upper:]' '[:lower:]' | tr ' ' '-' | tr -cd 'a-z0-9-' \
    | sed 's/--*/-/g; s/^-//; s/-$//' | head -c 30 | sed 's/-$//')
  if [ -z "$slug" ]; then
    echo "ERROR: slugify produced empty result from title" >&2
    exit 1
  fi
  echo "$slug"
}

cmd_find_rc() {
  # Fetch + prune to ensure remote tracking refs are current before listing RC branches
  git fetch --prune origin 2>/dev/null || true

  local branches
  branches=$(git ls-remote --heads origin 'rc/*' 2>/dev/null | awk '{print $2}' | sed 's|refs/heads/||') || true

  if [ -z "$branches" ]; then
    echo "NONE"
    exit 1
  fi

  local count
  count=$(echo "$branches" | wc -l | tr -d ' ')
  if [ "$count" -gt 1 ]; then
    echo "MULTIPLE" >&2
    echo "$branches" >&2
    exit 3
  fi

  echo "$branches"
}

cmd_next_version() {
  local prefix="${1:-}"
  if [ -z "$prefix" ]; then
    echo "ERROR: next-version requires PREFIX (e.g., V0.1)" >&2
    exit 1
  fi

  git fetch --tags origin 2>/dev/null || true

  # Find latest clean version tag matching PREFIX.N (exclude pre-release suffixes)
  local last_tag
  last_tag=$(git tag -l "${prefix}.*" --sort=-v:refname \
    | grep -E "^$(echo "$prefix" | sed 's/[.[\*^$()+?{|\\]/\\&/g')\.[0-9]+$" \
    | head -1) || true

  if [ -z "$last_tag" ]; then
    echo "${prefix}.0"
  else
    local patch prefix_part
    patch=$(echo "$last_tag" | awk -F. '{print $NF}')
    prefix_part=$(echo "$last_tag" | sed 's/\.[0-9]*$//')
    echo "${prefix_part}.$((patch + 1))"
  fi
}

cmd_list_merged_branches() {
  # Read branch prefix from config pattern (e.g., "mission/{issue}-{slug}-{user}" → "mission/*")
  # Falls back to "mission/*" if not configured.
  local pattern prefix branch_glob
  pattern=$(_branch_pattern)
  prefix="${pattern%%/\{*}"          # Extract prefix before first /{placeholder}
  prefix="${prefix%%\{*}"            # Handle patterns without / before placeholder
  if [ -z "$prefix" ] || [ "$prefix" = "$pattern" ]; then
    prefix="mission"
  fi
  branch_glob="${prefix}/*"

  local branches
  branches=$(git branch --list "$branch_glob" | sed 's/^[*+ ]*//' | tr -d ' ') || true

  if [ -z "$branches" ]; then
    exit 0
  fi

  while IFS= read -r branch; do
    [ -z "$branch" ] && continue
    local merged_pr
    merged_pr=$(gh pr list --head "$branch" --state merged --json number --jq '.[0].number' 2>/dev/null) || true
    if [ -n "$merged_pr" ]; then
      printf '%s\t#%s\n' "$branch" "$merged_pr"
    fi
  done <<< "$branches"
}

cmd_milestone_resolve() {
  local name="${1:-}"
  if [ -z "$name" ]; then
    echo "ERROR: milestone-resolve requires NAME" >&2
    exit 1
  fi

  local repo
  repo=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null) || {
    echo "$name"
    exit 0
  }

  local full
  full=$(gh api "repos/$repo/milestones" 2>/dev/null \
    | jq -r --arg name "$name" '.[] | select(.title | startswith($name)) | .title' 2>/dev/null | head -1) || true

  if [ -n "$full" ]; then
    echo "$full"
  else
    echo "$name"
  fi
}

cmd_build_branch_name() {
  # Like create-branch but ONLY returns the name — no git operations.
  # Use this when you need the branch name before deciding whether to create/checkout.
  local issue="${1:-}" slug="${2:-}" user="${3:-}"
  if [ -z "$issue" ] || [ -z "$slug" ] || [ -z "$user" ]; then
    echo "ERROR: build-branch-name requires ISSUE, SLUG, USER" >&2
    exit 1
  fi

  local pattern branch
  pattern=$(_branch_pattern)
  branch="${pattern//\{issue\}/$issue}"
  branch="${branch//\{slug\}/$slug}"
  branch="${branch//\{user\}/$user}"
  branch="${branch//\{type\}/mission}"
  branch="${branch//\{task-id\}/$issue}"
  echo "$branch"
}

cmd_extract_issue_number() {
  # Extract the first numeric segment after the first '/' in a branch name.
  # Works for: mission/42-slug, bugfix/T-044-slug, feature/123-name-user
  # If no branch given, uses current branch.
  local branch="${1:-}"
  if [ -z "$branch" ]; then
    branch=$(git branch --show-current 2>/dev/null) || true
  fi
  if [ -z "$branch" ]; then
    echo "ERROR: No branch provided and not on any branch" >&2
    exit 1
  fi

  local after_slash num
  after_slash="${branch#*/}"
  # Extract first contiguous digit sequence
  num=$(echo "$after_slash" | grep -oE '[0-9]+' | head -1) || true
  if [ -z "$num" ]; then
    echo "ERROR: No issue number found in branch '$branch'" >&2
    exit 1
  fi
  echo "$num"
}

cmd_worktree_main_repo() {
  # Print the path of the main worktree (the original clone, not any linked worktree).
  local main_path
  main_path=$(git worktree list --porcelain 2>/dev/null | head -1 | sed 's/^worktree //') || true
  if [ -z "$main_path" ]; then
    echo "ERROR: Could not determine main worktree" >&2
    exit 2
  fi
  echo "$main_path"
}

cmd_worktree_find_by_branch() {
  # Find the worktree path that has a specific branch checked out.
  # Prints the path if found, exits 1 if no worktree has that branch.
  local target_branch="${1:-}"
  if [ -z "$target_branch" ]; then
    echo "ERROR: worktree-find-by-branch requires BRANCH" >&2
    exit 1
  fi

  local result
  result=$(git worktree list --porcelain 2>/dev/null | awk -v b="$target_branch" '
    /^worktree / { wt = substr($0, 10) }
    /^branch refs\/heads\// {
      br = substr($0, 19)
      if (br == b) print wt
    }
  ') || true

  if [ -z "$result" ]; then
    exit 1
  fi
  echo "$result"
}

cmd_worktree_guard() {
  # Check if the current environment matches the worktree configuration.
  # Outputs KEY=VALUE pairs. Exits 3 on mismatch (worktree enabled but not in worktree).
  #
  # Usage: wing-git.sh worktree-guard [BRANCH]
  #   BRANCH — optional, defaults to current branch. Used to find existing worktree path.
  #
  # Output (always):
  #   WORKTREE_ENABLED=true|false
  #   IN_WORKTREE=true|false
  #   MISMATCH=true|false
  #   BRANCH=<branch-name>
  #   WORKTREE_PATH=<path>       (only if mismatch AND worktree exists for branch)
  #
  # Exit codes:
  #   0 — OK (no mismatch)
  #   3 — mismatch: worktree enabled but running from main repo

  local branch="${1:-}"
  if [ -z "$branch" ]; then
    branch=$(git branch --show-current 2>/dev/null) || true
  fi

  local worktree_enabled in_worktree
  worktree_enabled=$(git config --local wing.worktree 2>/dev/null || echo "false")
  in_worktree="false"
  [ -f ".mission" ] && in_worktree="true"

  if [ "$worktree_enabled" = "true" ] && [ "$in_worktree" = "false" ]; then
    # Mismatch: worktree enabled but running from main repo
    local wt_path=""
    if [ -n "$branch" ]; then
      wt_path=$(cmd_worktree_find_by_branch "$branch" 2>/dev/null) || true
    fi
    echo "WORKTREE_ENABLED=true"
    echo "IN_WORKTREE=false"
    echo "MISMATCH=true"
    echo "BRANCH=$branch"
    echo "WORKTREE_PATH=$wt_path"
    exit 3
  fi

  echo "WORKTREE_ENABLED=$worktree_enabled"
  echo "IN_WORKTREE=$in_worktree"
  echo "MISMATCH=false"
  echo "BRANCH=$branch"
}

cmd_claim_setup() {
  # Atomic branch+worktree setup for /team-claim.
  # Reads worktree config internally — the caller does NOT need to check worktree mode.
  #
  # Usage: wing-git.sh claim-setup ISSUE SLUG USER [TEAMWORK_DIR [CONTRACT_PATH]]
  #
  # Output (KEY=VALUE):
  #   MODE=branch|worktree
  #   BRANCH=<branch-name>
  #   STATUS=created|exists
  #   WORKTREE_PATH=<path>    (only when MODE=worktree)
  #
  # Exit codes:
  #   0 — success
  #   1 — usage error
  #   2 — git operation failed

  local issue="${1:-}" slug="${2:-}" user="${3:-}" teamwork_dir="${4:-}" contract_path="${5:-}"
  if [ -z "$issue" ] || [ -z "$slug" ] || [ -z "$user" ]; then
    echo "ERROR: claim-setup requires ISSUE SLUG USER [TEAMWORK_DIR [CONTRACT_PATH]]" >&2
    exit 1
  fi

  # Build branch name
  local branch
  branch=$(cmd_build_branch_name "$issue" "$slug" "$user")

  # Check worktree config
  local worktree_enabled
  worktree_enabled=$(git config --local wing.worktree 2>/dev/null || echo "false")

  if [ "$worktree_enabled" != "true" ]; then
    # ── NON-WORKTREE PATH ──
    cmd_ensure_base >/dev/null 2>&1 || {
      echo "ERROR: Could not switch to base branch" >&2
      exit 2
    }

    if git show-ref --verify --quiet "refs/heads/$branch" 2>/dev/null; then
      git checkout "$branch" 2>/dev/null || {
        echo "ERROR: Could not checkout existing branch '$branch'" >&2
        exit 2
      }
      echo "MODE=branch"
      echo "BRANCH=$branch"
      echo "STATUS=exists"
    else
      git checkout -b "$branch" 2>/dev/null || {
        echo "ERROR: Could not create branch '$branch'" >&2
        exit 2
      }
      echo "MODE=branch"
      echo "BRANCH=$branch"
      echo "STATUS=created"
    fi
  else
    # ── WORKTREE PATH ──
    local main_repo base worktree_root worktree_path repo_name

    main_repo=$(cmd_worktree_main_repo)
    repo_name=$(basename "$main_repo")
    base=$(_base_branch)

    # Fetch latest base (no checkout — do NOT disrupt working tree)
    git fetch origin "$base" --quiet 2>/dev/null || true

    # Resolve worktree path from per-user config
    worktree_root=$(git config --local wing.worktree-root 2>/dev/null || echo "")
    if [ -n "$worktree_root" ]; then
      # Resolve relative paths against main repo root
      if [[ "$worktree_root" != /* ]]; then
        worktree_root="${main_repo}/${worktree_root}"
      fi
      mkdir -p "$worktree_root"
      worktree_path="${worktree_root}/${slug}"
    else
      worktree_path="${main_repo}/../${repo_name}-wt-${slug}"
    fi

    # Check if branch is already checked out in a worktree
    local existing_wt
    existing_wt=$(cmd_worktree_find_by_branch "$branch" 2>/dev/null) || true

    if [ -n "$existing_wt" ]; then
      echo "MODE=worktree"
      echo "BRANCH=$branch"
      echo "WORKTREE_PATH=$existing_wt"
      echo "STATUS=exists"
      return 0
    fi

    # Create worktree
    if ! git worktree add -b "$branch" "$worktree_path" "origin/$base" 2>/dev/null; then
      # Branch may already exist but not checked out — try without -b
      if ! git worktree add "$worktree_path" "$branch" 2>/dev/null; then
        echo "ERROR: Could not create worktree at $worktree_path" >&2
        echo "Possible causes: path exists, permission denied, or branch checked out elsewhere" >&2
        exit 2
      fi
    fi

    # Write mission marker
    echo "$issue" > "$worktree_path/.mission"

    # Copy contract and config into worktree (they're in .gitignore active/)
    if [ -n "$teamwork_dir" ]; then
      mkdir -p "$worktree_path/$teamwork_dir/active" 2>/dev/null || true
      if [ -n "$contract_path" ] && [ -f "$contract_path" ]; then
        cp "$contract_path" "$worktree_path/$teamwork_dir/active/" 2>/dev/null || {
          echo "WARNING: Could not copy contract to worktree" >&2
        }
      fi
      if [ -f "$teamwork_dir/config.yml" ]; then
        cp "$teamwork_dir/config.yml" "$worktree_path/$teamwork_dir/" 2>/dev/null || {
          echo "WARNING: Could not copy config.yml to worktree" >&2
        }
      fi
    fi

    echo "MODE=worktree"
    echo "BRANCH=$branch"
    echo "WORKTREE_PATH=$worktree_path"
    echo "STATUS=created"
  fi
}

# --- Main ---
SUBCOMMAND="${1:-}"
shift || true

case "$SUBCOMMAND" in
  ensure-base)            cmd_ensure_base ;;
  create-branch)          cmd_create_branch "$@" ;;
  build-branch-name)      cmd_build_branch_name "$@" ;;
  extract-issue-number)   cmd_extract_issue_number "$@" ;;
  current)                cmd_current ;;
  protect-check)          cmd_protect_check ;;
  push)                   cmd_push "$@" ;;
  rebase)                 cmd_rebase "$@" ;;
  tag)                    cmd_tag "$@" ;;
  commit)                 cmd_commit "$@" ;;
  worktree-guard)         cmd_worktree_guard "$@" ;;
  claim-setup)            cmd_claim_setup "$@" ;;
  worktree-add)           cmd_worktree_add "$@" ;;
  worktree-remove)        cmd_worktree_remove "$@" ;;
  worktree-main-repo)     cmd_worktree_main_repo ;;
  worktree-find-by-branch) cmd_worktree_find_by_branch "$@" ;;
  log-since)              cmd_log_since "$@" ;;
  cut-release)            cmd_cut_release "$@" ;;
  slugify)                cmd_slugify "$@" ;;
  find-rc)                cmd_find_rc ;;
  next-version)           cmd_next_version "$@" ;;
  list-merged-branches)   cmd_list_merged_branches ;;
  milestone-resolve)      cmd_milestone_resolve "$@" ;;
  *)                      usage ;;
esac
