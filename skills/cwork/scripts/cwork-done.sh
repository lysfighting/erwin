#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONTRACT_SCRIPT="$SCRIPT_DIR/cwork-contract.sh"
GIT_SCRIPT="$SCRIPT_DIR/cwork-git.sh"
PR_SCRIPT="$SCRIPT_DIR/cwork-pr.sh"
LABEL_SCRIPT="$SCRIPT_DIR/cwork-label.sh"

usage() {
  echo "Usage: cwork-done.sh [--contract PATH] [--issue N] [--force] [--close-issue]" >&2
  exit 1
}

contract_path=""
issue=""
force="false"
close_issue="false"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --contract) contract_path="$2"; shift 2 ;;
    --issue) issue="$2"; shift 2 ;;
    --force) force="true"; shift ;;
    --close-issue) close_issue="true"; shift ;;
    *) usage ;;
  esac
done

[[ -n "$contract_path" ]] || contract_path="$(bash "$CONTRACT_SCRIPT" active-path "$issue" 2>/dev/null || true)"
[[ -n "$contract_path" && -f "$contract_path" ]] || {
  echo "ERROR: no active cwork contract found" >&2
  exit 1
}

issue="$(bash "$CONTRACT_SCRIPT" read-field "$contract_path" issue)"
branch="$(bash "$CONTRACT_SCRIPT" read-field "$contract_path" branch)"

merged="$(bash "$PR_SCRIPT" verify-merged "$branch" 2>/dev/null || true)"
if [[ -z "$merged" && "$force" != "true" ]]; then
  echo "ERROR: no merged PR found for branch '$branch' (use --force to clean up anyway)" >&2
  exit 1
fi

pr_number="$(printf '%s' "$merged" | awk '{print $1}')"
bash "$LABEL_SCRIPT" set "$issue" done >/dev/null 2>&1 || true

if [[ "$close_issue" == "true" ]]; then
  gh issue close "$issue" --comment "Completed via PR #${pr_number:-unknown}." >/dev/null 2>&1 || true
fi

bash "$CONTRACT_SCRIPT" delete "$contract_path" >/dev/null
rm -f .mission

printf 'ISSUE=%s\nBRANCH=%s\nPR=%s\nCONTRACT_PATH=%s\n' "$issue" "$branch" "${pr_number:-}" "$contract_path"
