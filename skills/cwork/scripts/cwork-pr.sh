#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_SCRIPT="$SCRIPT_DIR/cwork-config.sh"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/cwork-common.sh"

usage() {
  echo "Usage: cwork-pr.sh <subcommand> [args...]" >&2
  echo "Subcommands: exists, create, update, watch, wait-merged, verify-merged, comment" >&2
  exit 1
}

base_branch() {
  local value
  value="$(bash "$CONFIG_SCRIPT" get conventions.base_branch "main" 2>/dev/null || true)"
  [[ -n "$value" ]] || value="main"
  printf '%s\n' "$value"
}

cmd_exists() {
  local branch="${1:-}" base="${2:-}"
  [[ -n "$branch" ]] || {
    echo "ERROR: exists requires BRANCH [BASE]" >&2
    exit 1
  }
  local args=()
  [[ -n "$base" ]] && args+=(--base "$base")
  local result
  result="$(gh pr list --head "$branch" "${args[@]}" --state open --json number,url --limit 1 2>/dev/null || true)"
  [[ -n "$result" && "$(printf '%s' "$result" | jq 'length')" -gt 0 ]] || exit 1
  printf '%s %s\n' "$(printf '%s' "$result" | jq -r '.[0].number')" "$(printf '%s' "$result" | jq -r '.[0].url')"
}

cmd_create() {
  local issue="${1:-}" title="${2:-}" body_file="${3:-}" base="${4:-}" branch="${5:-}" mode="${6:-ready}"
  [[ -n "$issue" && -n "$title" && -n "$body_file" && -n "$branch" ]] || {
    echo "ERROR: create requires ISSUE TITLE BODY_FILE BASE BRANCH [MODE]" >&2
    exit 1
  }
  [[ -n "$base" ]] || base="$(base_branch)"
  local args=(--title "$title" --body-file "$body_file" --base "$base" --head "$branch")
  [[ "$mode" == "draft" ]] && args=(--draft "${args[@]}")
  local url
  url="$(gh pr create "${args[@]}" 2>/dev/null)"
  url="$(printf '%s' "$url" | tail -1)"
  printf '%s %s\n' "$(printf '%s' "$url" | grep -Eo '[0-9]+$' || true)" "$url"
}

cmd_update() {
  local pr="${1:-}" title="${2:-}" body_file="${3:-}"
  [[ -n "$pr" && -n "$title" && -n "$body_file" ]] || {
    echo "ERROR: update requires PR TITLE BODY_FILE" >&2
    exit 1
  }
  gh pr edit "$pr" --title "$title" --body "$(cat "$body_file")" >/dev/null
  printf 'UPDATED %s\n' "$pr"
}

cmd_watch() {
  local pr="${1:-}"
  [[ -n "$pr" ]] || {
    echo "ERROR: watch requires PR" >&2
    exit 1
  }
  gh pr checks "$pr" --watch --interval 10
}

cmd_wait_merged() {
  local pr="${1:-}"
  [[ -n "$pr" ]] || {
    echo "ERROR: wait-merged requires PR" >&2
    exit 1
  }
  local poll=0
  while [[ $poll -lt 60 ]]; do
    local state
    state="$(gh pr view "$pr" --json state --jq '.state' 2>/dev/null || true)"
    case "$state" in
      MERGED) printf 'MERGED %s\n' "$pr"; return 0 ;;
      CLOSED) echo "ERROR: PR #$pr closed without merging" >&2; return 1 ;;
    esac
    poll=$((poll + 1))
    sleep 10
  done
  echo "ERROR: timed out waiting for PR #$pr to merge" >&2
  exit 3
}

cmd_verify_merged() {
  local branch="${1:-}"
  [[ -n "$branch" ]] || {
    echo "ERROR: verify-merged requires BRANCH" >&2
    exit 1
  }
  local result
  result="$(gh pr list --head "$branch" --state merged --json number,url --limit 1 2>/dev/null || true)"
  [[ -n "$result" && "$(printf '%s' "$result" | jq 'length')" -gt 0 ]] || exit 1
  printf '%s %s\n' "$(printf '%s' "$result" | jq -r '.[0].number')" "$(printf '%s' "$result" | jq -r '.[0].url')"
}

cmd_comment() {
  local issue="${1:-}" text="${2:-}"
  [[ -n "$issue" && -n "$text" ]] || {
    echo "ERROR: comment requires ISSUE TEXT" >&2
    exit 1
  }
  gh issue comment "$issue" --body "$text" >/dev/null
  printf 'COMMENTED %s\n' "$issue"
}

subcommand="${1:-}"
shift || true

case "$subcommand" in
  exists) cmd_exists "$@" ;;
  create) cmd_create "$@" ;;
  update) cmd_update "$@" ;;
  watch) cmd_watch "$@" ;;
  wait-merged) cmd_wait_merged "$@" ;;
  verify-merged) cmd_verify_merged "$@" ;;
  comment) cmd_comment "$@" ;;
  *) usage ;;
esac
