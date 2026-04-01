#!/usr/bin/env bash
set -euo pipefail

export GH_PROMPT_DISABLED=1

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_SCRIPT="$SCRIPT_DIR/cwork-config.sh"

usage() {
  echo "Usage: cwork-label.sh <subcommand> [args...]" >&2
  echo "Subcommands: transition, set, verify, pr-label" >&2
  exit 1
}

resolve_labels() {
  eval "$(bash "$CONFIG_SCRIPT" resolve-labels 2>/dev/null || true)"
  MISSION_LABEL="${MISSION_LABEL:-mission}"
  STATUS_PREFIX="${STATUS_PREFIX:-status:}"
  PRIORITY_PREFIX="${PRIORITY_PREFIX:-priority:}"
}

cmd_transition() {
  local issue="${1:-}" from="${2:-}" to="${3:-}"
  [[ -n "$issue" && -n "$from" && -n "$to" ]] || usage
  resolve_labels
  gh issue edit "$issue" --remove-label "${STATUS_PREFIX}${from}" --add-label "${STATUS_PREFIX}${to}" >/dev/null
  printf '%s%s\n' "$STATUS_PREFIX" "$to"
}

cmd_set() {
  local issue="${1:-}" status="${2:-}"
  [[ -n "$issue" && -n "$status" ]] || usage
  resolve_labels
  local remove_args=()
  local known
  for known in queued wip review done blocked; do
    remove_args+=(--remove-label "${STATUS_PREFIX}${known}")
  done
  gh issue edit "$issue" "${remove_args[@]}" --add-label "${STATUS_PREFIX}${status}" >/dev/null
  printf '%s%s\n' "$STATUS_PREFIX" "$status"
}

cmd_verify() {
  local issue="${1:-}" expected="${2:-}"
  [[ -n "$issue" && -n "$expected" ]] || usage
  resolve_labels
  local labels
  labels="$(gh issue view "$issue" --json labels --jq '[.labels[].name] | join(",")' 2>/dev/null || true)"
  printf ',%s,' "$labels" | grep -q ",${STATUS_PREFIX}${expected},"
}

cmd_pr_label() {
  local pr="${1:-}" status="${2:-}"
  [[ -n "$pr" && -n "$status" ]] || usage
  resolve_labels
  gh pr edit "$pr" --add-label "${STATUS_PREFIX}${status}" >/dev/null
  printf '%s%s\n' "$STATUS_PREFIX" "$status"
}

subcommand="${1:-}"
shift || true

case "$subcommand" in
  transition) cmd_transition "$@" ;;
  set) cmd_set "$@" ;;
  verify) cmd_verify "$@" ;;
  pr-label) cmd_pr_label "$@" ;;
  *) usage ;;
esac
