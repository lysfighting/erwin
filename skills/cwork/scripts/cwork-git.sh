#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_SCRIPT="$SCRIPT_DIR/cwork-config.sh"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/cwork-common.sh"

usage() {
  echo "Usage: cwork-git.sh <subcommand> [args...]" >&2
  echo "Subcommands: current, protect-check, push, rebase, slugify, build-branch-name, extract-issue-number, worktree-main-repo, worktree-find-by-branch, worktree-guard, claim-setup" >&2
  exit 1
}

base_branch() {
  local value
  value="$(bash "$CONFIG_SCRIPT" get conventions.base_branch "main" 2>/dev/null || true)"
  [[ -n "$value" ]] || value="main"
  printf '%s\n' "$value"
}

production_branch() {
  local value
  value="$(bash "$CONFIG_SCRIPT" get conventions.production_branch "" 2>/dev/null || true)"
  if [[ -n "$value" ]]; then
    printf '%s\n' "$value"
  else
    base_branch
  fi
}

branch_pattern() {
  local value
  value="$(bash "$CONFIG_SCRIPT" get conventions.branch_pattern "" 2>/dev/null || true)"
  [[ -n "$value" ]] || value="$(bash "$CONFIG_SCRIPT" get worktree.branch_pattern "" 2>/dev/null || true)"
  [[ -n "$value" ]] || value="mission/{issue}-{slug}-{user}"
  printf '%s\n' "$value"
}

worktree_enabled() {
  local value
  value="$(git config --local cwork.worktree 2>/dev/null || true)"
  [[ -n "$value" ]] || value="$(git config --local teamwork.worktree 2>/dev/null || true)"
  [[ -n "$value" ]] || value="false"
  printf '%s\n' "$value"
}

worktree_root() {
  local value
  value="$(git config --local cwork.worktree-root 2>/dev/null || true)"
  [[ -n "$value" ]] || value="$(git config --local teamwork.worktree-root 2>/dev/null || true)"
  printf '%s\n' "$value"
}

cmd_current() {
  git branch --show-current
}

cmd_protect_check() {
  local current base prod
  current="$(cmd_current)"
  base="$(base_branch)"
  prod="$(production_branch)"

  if [[ -z "$current" ]]; then
    printf 'OK: detached HEAD\n'
    return 0
  fi
  if [[ "$current" == "$base" || "$current" == "$prod" || "$current" == rc/* ]]; then
    echo "SAFETY: on protected branch '$current'" >&2
    exit 3
  fi
  printf 'OK: on %s\n' "$current"
}

cmd_push() {
  local branch="${1:-}"
  [[ -n "$branch" ]] || branch="$(cmd_current)"
  git push -u origin "$branch"
}

cmd_rebase() {
  local target="${1:-}"
  [[ -n "$target" ]] || target="$(base_branch)"
  git fetch origin "$target" >/dev/null 2>&1 || true
  git rebase "origin/$target"
}

cmd_slugify() {
  local title="${1:-}"
  [[ -n "$title" ]] || {
    echo "ERROR: slugify requires TITLE" >&2
    exit 1
  }
  local slug
  slug="$(printf '%s' "$title" | tr '[:upper:]' '[:lower:]' | tr ' ' '-' | tr -cd 'a-z0-9-' | sed 's/--*/-/g; s/^-//; s/-$//' | head -c 40 | sed 's/-$//')"
  [[ -n "$slug" ]] || slug="mission"
  printf '%s\n' "$slug"
}

cmd_build_branch_name() {
  local issue="${1:-}" slug="${2:-}" user="${3:-codex}"
  [[ -n "$issue" && -n "$slug" ]] || {
    echo "ERROR: build-branch-name requires ISSUE and SLUG [USER]" >&2
    exit 1
  }
  local pattern branch
  pattern="$(branch_pattern)"
  branch="${pattern//\{issue\}/$issue}"
  branch="${branch//\{slug\}/$slug}"
  branch="${branch//\{user\}/$user}"
  branch="${branch//\{type\}/mission}"
  branch="${branch//\{task-id\}/$issue}"
  printf '%s\n' "$branch"
}

cmd_extract_issue_number() {
  local branch="${1:-}"
  [[ -n "$branch" ]] || branch="$(cmd_current)"
  [[ -n "$branch" ]] || {
    echo "ERROR: no branch available" >&2
    exit 1
  }
  local after_slash num
  after_slash="${branch#*/}"
  num="$(printf '%s' "$after_slash" | grep -Eo '[0-9]+' | head -1 || true)"
  [[ -n "$num" ]] || {
    echo "ERROR: no issue number found in branch '$branch'" >&2
    exit 1
  }
  printf '%s\n' "$num"
}

cmd_worktree_main_repo() {
  git worktree list --porcelain | awk '/^worktree / { print substr($0, 10); exit }'
}

cmd_worktree_find_by_branch() {
  local target_branch="${1:-}"
  [[ -n "$target_branch" ]] || {
    echo "ERROR: worktree-find-by-branch requires BRANCH" >&2
    exit 1
  }
  git worktree list --porcelain | awk -v target="$target_branch" '
    /^worktree / { wt = substr($0, 10) }
    /^branch refs\/heads\// {
      branch = substr($0, 19)
      if (branch == target) print wt
    }
  ' | head -1
}

cmd_worktree_guard() {
  local branch="${1:-}"
  [[ -n "$branch" ]] || branch="$(cmd_current)"
  local enabled in_worktree wt_path
  enabled="$(worktree_enabled)"
  in_worktree="false"
  [[ -f .mission ]] && in_worktree="true"
  wt_path=""
  if [[ "$enabled" == "true" && "$in_worktree" == "false" && -n "$branch" ]]; then
    wt_path="$(cmd_worktree_find_by_branch "$branch" || true)"
    printf 'WORKTREE_ENABLED=true\nIN_WORKTREE=false\nMISMATCH=true\nBRANCH=%s\nWORKTREE_PATH=%s\n' "$branch" "$wt_path"
    exit 3
  fi
  printf 'WORKTREE_ENABLED=%s\nIN_WORKTREE=%s\nMISMATCH=false\nBRANCH=%s\n' "$enabled" "$in_worktree" "$branch"
}

cmd_claim_setup() {
  local issue="${1:-}" slug="${2:-}" user="${3:-codex}" cwork_dir="${4:-.cwork}" contract_path="${5:-}"
  [[ -n "$issue" && -n "$slug" ]] || {
    echo "ERROR: claim-setup requires ISSUE SLUG [USER [CWORK_DIR [CONTRACT_PATH]]]" >&2
    exit 1
  }

  local branch enabled
  branch="$(cmd_build_branch_name "$issue" "$slug" "$user")"
  enabled="$(worktree_enabled)"

  if [[ "$enabled" != "true" ]]; then
    if git show-ref --verify --quiet "refs/heads/$branch"; then
      git checkout "$branch" >/dev/null 2>&1
      printf 'MODE=branch\nBRANCH=%s\nSTATUS=exists\n' "$branch"
      return 0
    fi
    local base
    base="$(base_branch)"
    git fetch origin "$base" >/dev/null 2>&1 || true
    if git show-ref --verify --quiet "refs/remotes/origin/$base"; then
      git checkout -b "$branch" "origin/$base" >/dev/null 2>&1 || git checkout -b "$branch" >/dev/null 2>&1
    else
      git checkout -b "$branch" >/dev/null 2>&1
    fi
    printf 'MODE=branch\nBRANCH=%s\nSTATUS=created\n' "$branch"
    return 0
  fi

  local main_repo base repo_name wt_root wt_path existing
  main_repo="$(cmd_worktree_main_repo)"
  base="$(base_branch)"
  repo_name="$(basename "$main_repo")"
  wt_root="$(worktree_root)"
  if [[ -n "$wt_root" && "$wt_root" != /* ]]; then
    wt_root="$main_repo/$wt_root"
  fi
  [[ -n "$wt_root" ]] || wt_root="$main_repo/../${repo_name}-wt"
  mkdir -p "$wt_root"
  wt_path="${wt_root}/${slug}"

  existing="$(cmd_worktree_find_by_branch "$branch" || true)"
  if [[ -n "$existing" ]]; then
    printf 'MODE=worktree\nBRANCH=%s\nWORKTREE_PATH=%s\nSTATUS=exists\n' "$branch" "$existing"
    return 0
  fi

  git fetch origin "$base" >/dev/null 2>&1 || true
  if ! git worktree add -b "$branch" "$wt_path" "origin/$base" >/dev/null 2>&1; then
    git worktree add "$wt_path" "$branch" >/dev/null 2>&1
  fi
  printf '%s\n' "$issue" >"$wt_path/.mission"
  if [[ -n "$contract_path" && -f "$contract_path" ]]; then
    mkdir -p "$wt_path/$cwork_dir/active"
    cp "$contract_path" "$wt_path/$cwork_dir/active/"
  fi
  printf 'MODE=worktree\nBRANCH=%s\nWORKTREE_PATH=%s\nSTATUS=created\n' "$branch" "$wt_path"
}

subcommand="${1:-}"
shift || true

case "$subcommand" in
  current) cmd_current "$@" ;;
  protect-check) cmd_protect_check "$@" ;;
  push) cmd_push "$@" ;;
  rebase) cmd_rebase "$@" ;;
  slugify) cmd_slugify "$@" ;;
  build-branch-name) cmd_build_branch_name "$@" ;;
  extract-issue-number) cmd_extract_issue_number "$@" ;;
  worktree-main-repo) cmd_worktree_main_repo "$@" ;;
  worktree-find-by-branch) cmd_worktree_find_by_branch "$@" ;;
  worktree-guard) cmd_worktree_guard "$@" ;;
  claim-setup) cmd_claim_setup "$@" ;;
  *) usage ;;
esac
