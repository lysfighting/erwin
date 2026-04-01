#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONTRACT_SCRIPT="$SCRIPT_DIR/cwork-contract.sh"
GIT_SCRIPT="$SCRIPT_DIR/cwork-git.sh"
LABEL_SCRIPT="$SCRIPT_DIR/cwork-label.sh"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/cwork-common.sh"

usage() {
  echo "Usage: cwork-claim.sh ISSUE_NUMBER [--json ISSUE_JSON] [--branch BRANCH] [--no-comment]" >&2
  exit 1
}

issue="${1:-}"
shift || true
[[ -n "$issue" ]] || usage

json_path=""
branch_override=""
post_comment="true"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --json) json_path="$2"; shift 2 ;;
    --branch) branch_override="$2"; shift 2 ;;
    --no-comment) post_comment="false"; shift ;;
    *) usage ;;
  esac
done

cwork_require_cmd git
cwork_require_cmd jq

mkdir -p "$(cwork_active_dir)"

local_user="$(cwork_current_user)"
if [[ -n "$json_path" ]]; then
  title="$(jq -r '.title // ""' "$json_path")"
else
  cwork_require_cmd gh
  title="$(gh issue view "$issue" --json title --jq '.title' 2>/dev/null)"
fi
slug="$(bash "$GIT_SCRIPT" slugify "$title")"
branch="${branch_override:-$(bash "$GIT_SCRIPT" build-branch-name "$issue" "$slug" "$local_user")}"
contract_path="$(cwork_contract_path "$issue")"

if [[ -n "$json_path" ]]; then
  bash "$CONTRACT_SCRIPT" from-json "$json_path" --branch "$branch" --output "$contract_path" >/dev/null
else
  bash "$CONTRACT_SCRIPT" from-issue "$issue" --branch "$branch" --output "$contract_path" >/dev/null
fi

setup_output="$(bash "$GIT_SCRIPT" claim-setup "$issue" "$slug" "$local_user" ".cwork" "$contract_path")"
eval "$setup_output"

if [[ "$post_comment" == "true" && -z "$json_path" ]]; then
  gh issue comment "$issue" --body "Claimed in Codex on branch \`$branch\`." >/dev/null 2>&1 || true
fi

if [[ -z "$json_path" ]]; then
  bash "$LABEL_SCRIPT" set "$issue" wip >/dev/null 2>&1 || true
fi

printf 'ISSUE=%s\n' "$issue"
printf 'TITLE=%s\n' "$title"
printf 'BRANCH=%s\n' "$branch"
printf 'CONTRACT_PATH=%s\n' "$contract_path"
printf '%s\n' "$setup_output"
