#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_SCRIPT="$SCRIPT_DIR/cwork-config.sh"
CONTRACT_SCRIPT="$SCRIPT_DIR/cwork-contract.sh"
GIT_SCRIPT="$SCRIPT_DIR/cwork-git.sh"
PR_SCRIPT="$SCRIPT_DIR/cwork-pr.sh"
LABEL_SCRIPT="$SCRIPT_DIR/cwork-label.sh"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/cwork-common.sh"

usage() {
  echo "Usage: cwork-ship.sh [--contract PATH] [--issue N] [--draft] [--update-existing] [--dry-run] [--no-sync]" >&2
  exit 1
}

contract_path=""
issue=""
mode="ready"
update_existing="false"
dry_run="false"
sync_checked="true"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --contract) contract_path="$2"; shift 2 ;;
    --issue) issue="$2"; shift 2 ;;
    --draft) mode="draft"; shift ;;
    --update-existing) update_existing="true"; shift ;;
    --dry-run) dry_run="true"; shift ;;
    --no-sync) sync_checked="false"; shift ;;
    *) usage ;;
  esac
done

[[ -n "$contract_path" ]] || contract_path="$(bash "$CONTRACT_SCRIPT" active-path "$issue" 2>/dev/null || true)"
[[ -n "$contract_path" && -f "$contract_path" ]] || {
  echo "ERROR: no active cwork contract found" >&2
  exit 1
}

issue="$(bash "$CONTRACT_SCRIPT" read-field "$contract_path" issue)"
title="$(bash "$CONTRACT_SCRIPT" read-field "$contract_path" title)"
branch="$(bash "$CONTRACT_SCRIPT" read-field "$contract_path" branch)"
current_branch="$(bash "$GIT_SCRIPT" current)"
base_branch="$(bash "$CONFIG_SCRIPT" get conventions.base_branch "main" 2>/dev/null || true)"
[[ -n "$base_branch" ]] || base_branch="main"

[[ "$current_branch" == "$branch" ]] || {
  echo "ERROR: contract expects branch '$branch' but current branch is '$current_branch'" >&2
  exit 1
}

bash "$GIT_SCRIPT" protect-check >/dev/null

tmp_body="$(mktemp)"
trap 'rm -f "${tmp_body:-}"' EXIT
bash "$CONTRACT_SCRIPT" pr-body "$contract_path" >"$tmp_body"

freshness="$(bash "$CONTRACT_SCRIPT" check-freshness "$contract_path" "$issue" 2>/dev/null || true)"
if [[ "$sync_checked" == "true" && "$dry_run" != "true" ]]; then
  bash "$CONTRACT_SCRIPT" sync-checked "$contract_path" "$issue" >/dev/null 2>&1 || true
fi

if [[ "$dry_run" == "true" ]]; then
  printf 'ISSUE=%s\nTITLE=%s\nBRANCH=%s\nBASE=%s\nFRESHNESS=%s\nBODY_FILE=%s\n' "$issue" "$title" "$branch" "$base_branch" "${freshness:-UNKNOWN}" "$tmp_body"
  cat "$tmp_body"
  exit 0
fi

bash "$GIT_SCRIPT" push "$branch" >/dev/null

existing="$(bash "$PR_SCRIPT" exists "$branch" "$base_branch" 2>/dev/null || true)"
if [[ -n "$existing" ]]; then
  pr_number="$(printf '%s' "$existing" | awk '{print $1}')"
  pr_url="$(printf '%s' "$existing" | awk '{print $2}')"
  if [[ "$update_existing" == "true" ]]; then
    bash "$PR_SCRIPT" update "$pr_number" "$title" "$tmp_body" >/dev/null
  fi
else
  created="$(bash "$PR_SCRIPT" create "$issue" "$title" "$tmp_body" "$base_branch" "$branch" "$mode")"
  pr_number="$(printf '%s' "$created" | awk '{print $1}')"
  pr_url="$(printf '%s' "$created" | awk '{print $2}')"
fi

bash "$LABEL_SCRIPT" set "$issue" review >/dev/null 2>&1 || true

printf 'ISSUE=%s\nPR=%s\nURL=%s\nBRANCH=%s\nCONTRACT_PATH=%s\nFRESHNESS=%s\n' "$issue" "${pr_number:-}" "${pr_url:-}" "$branch" "$contract_path" "${freshness:-UNKNOWN}"
