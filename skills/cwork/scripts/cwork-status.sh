#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONTRACT_SCRIPT="$SCRIPT_DIR/cwork-contract.sh"
GIT_SCRIPT="$SCRIPT_DIR/cwork-git.sh"
PR_SCRIPT="$SCRIPT_DIR/cwork-pr.sh"
CONFIG_SCRIPT="$SCRIPT_DIR/cwork-config.sh"

branch="$(bash "$GIT_SCRIPT" current)"
contract_path="$(bash "$CONTRACT_SCRIPT" active-path 2>/dev/null || true)"
base="$(bash "$CONFIG_SCRIPT" get conventions.base_branch "main" 2>/dev/null || true)"
[[ -n "$base" ]] || base="main"

printf 'BRANCH=%s\nBASE=%s\n' "$branch" "$base"
if [[ -n "$contract_path" ]]; then
  issue="$(bash "$CONTRACT_SCRIPT" read-field "$contract_path" issue)"
  title="$(bash "$CONTRACT_SCRIPT" read-field "$contract_path" title)"
  freshness="$(bash "$CONTRACT_SCRIPT" check-freshness "$contract_path" "$issue" 2>/dev/null || true)"
  printf 'CONTRACT_PATH=%s\nISSUE=%s\nTITLE=%s\nFRESHNESS=%s\n' "$contract_path" "$issue" "$title" "${freshness:-UNKNOWN}"
fi
existing="$(bash "$PR_SCRIPT" exists "$branch" "$base" 2>/dev/null || true)"
if [[ -n "$existing" ]]; then
  printf 'OPEN_PR=%s\n' "$existing"
fi
