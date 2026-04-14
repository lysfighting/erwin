#!/usr/bin/env bash
set -euo pipefail

export GH_PROMPT_DISABLED=1

cwork_script_dir() {
  cd "$(dirname "${BASH_SOURCE[0]}")" && pwd
}

cwork_repo_root() {
  local root
  root="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
  cd "$root" 2>/dev/null && pwd -P
}

cwork_main_repo_root() {
  local common_git
  common_git=$(git rev-parse --git-common-dir 2>/dev/null || printf '.git')
  cd "$common_git/.." 2>/dev/null && pwd -P
}

cwork_local_dir() {
  printf '%s/.cwork\n' "$(cwork_repo_root)"
}

cwork_main_dir() {
  printf '%s/.cwork\n' "$(cwork_main_repo_root)"
}

cwork_active_dir() {
  printf '%s/active\n' "$(cwork_local_dir)"
}

cwork_contract_path() {
  local issue="${1:-}"
  if [[ -z "$issue" ]]; then
    echo "ERROR: cwork_contract_path requires ISSUE" >&2
    return 1
  fi
  printf '%s/MISSION-%s.md\n' "$(cwork_active_dir)" "$issue"
}

cwork_find_contract() {
  local issue="${1:-}"
  local local_dir main_dir
  local_dir="$(cwork_active_dir)"
  main_dir="$(cwork_main_dir)/active"

  if [[ -n "$issue" ]]; then
    if [[ -f "$local_dir/MISSION-$issue.md" ]]; then
      printf '%s/MISSION-%s.md\n' "$local_dir" "$issue"
      return 0
    fi
    if [[ -f "$main_dir/MISSION-$issue.md" ]]; then
      printf '%s/MISSION-%s.md\n' "$main_dir" "$issue"
      return 0
    fi
    return 1
  fi

  local matches=()
  while IFS= read -r path; do
    [[ -n "$path" ]] && matches+=("$path")
  done < <(find "$local_dir" "$main_dir" -maxdepth 1 -type f -name 'MISSION-*.md' 2>/dev/null | awk '!seen[$0]++')

  if [[ ${#matches[@]} -eq 1 ]]; then
    printf '%s\n' "${matches[0]}"
    return 0
  fi
  return 1
}

cwork_require_cmd() {
  local cmd="${1:-}"
  command -v "$cmd" >/dev/null 2>&1 || {
    echo "ERROR: Missing required command: $cmd" >&2
    return 1
  }
}

cwork_trim() {
  local value="${1:-}"
  value="${value#"${value%%[![:space:]]*}"}"
  value="${value%"${value##*[![:space:]]}"}"
  printf '%s' "$value"
}

cwork_current_user() {
  if [[ -n "${GITHUB_USER:-}" ]]; then
    printf '%s\n' "$GITHUB_USER"
    return 0
  fi
  if [[ -n "${USER:-}" ]]; then
    printf '%s\n' "$USER"
    return 0
  fi
  printf 'codex\n'
}

cwork_parse_key_value_lines() {
  local line
  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    eval "$line"
  done
}
