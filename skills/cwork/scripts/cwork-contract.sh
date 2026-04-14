#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/cwork-common.sh"

usage() {
  cat >&2 <<'EOF'
Usage:
  cwork-contract.sh from-issue ISSUE_NUMBER [--output PATH] [--branch BRANCH]
  cwork-contract.sh from-json ISSUE_JSON_PATH [--output PATH] [--branch BRANCH]
  cwork-contract.sh read-field PATH FIELD
  cwork-contract.sh section PATH HEADING [ALT_HEADING...]
  cwork-contract.sh checked-items PATH
  cwork-contract.sh active-path [ISSUE]
  cwork-contract.sh check-freshness PATH ISSUE
  cwork-contract.sh sync-checked PATH ISSUE
  cwork-contract.sh pr-body PATH
  cwork-contract.sh delete PATH
EOF
  exit 1
}

yaml_escape() {
  printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}

normalize_item() {
  printf '%s' "$1" \
    | sed 's/[[:space:]]*[—-][[:space:]]*[0-9][0-9]*:[0-9][0-9]*$//' \
    | sed 's/^🔧[[:space:]]*MANUAL:[[:space:]]*//' \
    | tr '[:upper:]' '[:lower:]' \
    | sed 's/[[:space:]]\+/ /g; s/^ //; s/ $//'
}

slugify() {
  printf '%s' "$1" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g; s/-\{2,\}/-/g; s/^-//; s/-$//'
}

infer_branch() {
  local issue_number="$1"
  local title="$2"
  local live_branch
  live_branch="$(cwork_trim "$(git branch --show-current 2>/dev/null || true)")"
  case "$live_branch" in
    ""|HEAD|main|master|develop|dev|pre-launch)
      printf 'mission/%s-%s\n' "$issue_number" "$(slugify "$title")"
      ;;
    *)
      printf '%s\n' "$live_branch"
      ;;
  esac
}

read_field() {
  local path="$1"
  local field="$2"
  awk -v target="$field" '
    BEGIN { in_frontmatter = 0 }
    /^---$/ {
      if (in_frontmatter == 0) {
        in_frontmatter = 1
        next
      }
      exit
    }
    in_frontmatter == 1 {
      pos = index($0, ":")
      if (pos > 0) {
        key = substr($0, 1, pos - 1)
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", key)
        if (key == target) {
          value = substr($0, pos + 1)
          gsub(/^[[:space:]]+|[[:space:]]+$/, "", value)
          gsub(/^"/, "", value)
          gsub(/"$/, "", value)
          print value
          exit
        }
      }
    }
  ' "$path"
}

section_text() {
  local path="$1"
  shift
  local wanted
  wanted="$(printf '%s|' "$@" | sed 's/|$//')"
  awk -v wanted="$wanted" '
    function norm(s) {
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", s)
      s = tolower(s)
      gsub(/[[:space:]]+/, " ", s)
      return s
    }
    BEGIN {
      count = split(wanted, raw, /\|/)
    }
    /^#{1,6}[[:space:]]+/ {
      heading = $0
      sub(/^#{1,6}[[:space:]]+/, "", heading)
      current = norm(heading)
      if (capture) exit
      for (i = 1; i <= count; i++) {
        if (current == norm(raw[i])) {
          capture = 1
          next
        }
      }
      next
    }
    capture { print }
  ' "$path" | awk '
    { lines[NR] = $0 }
    END {
      start = 1
      while (start <= NR && lines[start] ~ /^[[:space:]]*$/) start++
      stop = NR
      while (stop >= start && lines[stop] ~ /^[[:space:]]*$/) stop--
      for (i = start; i <= stop; i++) print lines[i]
    }
  '
}

list_items() {
  awk '
    {
      line = $0
      sub(/^[[:space:]]+/, "", line)
      if (line ~ /^-[[:space:]]+\[[ xX]\][[:space:]]+/) {
        sub(/^-[[:space:]]+\[[ xX]\][[:space:]]+/, "", line)
        print line
      } else if (line ~ /^-[[:space:]]+/) {
        sub(/^-[[:space:]]+/, "", line)
        print line
      } else if (line ~ /^[0-9]+\.[[:space:]]+/) {
        sub(/^[0-9]+\.[[:space:]]+/, "", line)
        print line
      }
    }
  ' | awk '!seen[tolower($0)]++'
}

checked_items() {
  local path="$1"
  awk '/^- \[x\] / { line = $0; sub(/^- \[x\] /, "", line); print line }' "$path"
}

generate_contract() {
  local json_file="$1"
  local output="$2"
  local branch_override="${3:-}"
  local body_file
  body_file="$(mktemp)"
  trap 'rm -f "${body_file:-}"' RETURN

  local issue_number title body url milestone assignee priority branch generated_at issue_hash
  issue_number="$(jq -r '.number' "$json_file")"
  title="$(jq -r '.title // ""' "$json_file")"
  body="$(jq -r '.body // ""' "$json_file")"
  url="$(jq -r '.url // ""' "$json_file")"
  milestone="$(jq -r '.milestone.title // ""' "$json_file")"
  assignee="$(jq -r '.assignees[0].login // ""' "$json_file")"
  printf '%s' "$body" >"$body_file"

  priority="$(jq -r '.labels[]?.name // empty' "$json_file" | awk '
    {
      upper = toupper($0)
      if (upper ~ /^P[0-3]$/) { print upper; exit }
      if ($0 ~ /^[Pp][Rr][Ii][Oo][Rr][Ii][Tt][Yy]:(P[0-3])$/) { split($0, a, ":"); print toupper(a[2]); exit }
    }
    END { if (NR == 0) print "P2" }
  ')"
  [[ -n "$priority" ]] || priority="P2"
  branch="${branch_override:-$(infer_branch "$issue_number" "$title")}"
  generated_at="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
  issue_hash="$(printf '%s' "${title}${body}" | shasum -a 256 | awk '{print $1}')"

  local objective subtasks success context verification ai_notes
  objective="$(section_text "$body_file" "Objective")"
  if [[ -n "$(cwork_trim "$objective")" ]]; then
    objective="$(printf '%s' "$objective" | tr '\n' ' ' | sed 's/[[:space:]]\+/ /g; s/^ //; s/ $//')"
  else
    objective="Complete issue #$issue_number: $title"
  fi

  subtasks="$(section_text "$body_file" "Sub-tasks" "Subtasks" "Tasks" "Steps" | list_items || true)"
  success="$(section_text "$body_file" "Success Criteria" "Acceptance Criteria" | list_items || true)"
  context="$(section_text "$body_file" "Context & References" "Context" || true)"
  verification="$(section_text "$body_file" "Verification Method" "Verification" "Test Command" || true)"

  ai_notes=""
  [[ -n "$(section_text "$body_file" "Objective" || true)" ]] || ai_notes="${ai_notes}- Objective was inferred from freeform issue text. Re-check intent before coding.\n"
  [[ -n "$(cwork_trim "$subtasks")" ]] || ai_notes="${ai_notes}- No sub-tasks were provided. Break the work down before implementing.\n"
  [[ -n "$(cwork_trim "$success")" ]] || ai_notes="${ai_notes}- No explicit success criteria were provided. Add concrete verification checkpoints.\n"
  if [[ "$(printf '%s' "$objective" | wc -w | tr -d ' ')" -gt 40 ]]; then
    ai_notes="${ai_notes}- The objective is broad. Consider tightening scope before large edits.\n"
  fi
  [[ -n "$ai_notes" ]] || ai_notes="- Validate each sub-task against the codebase and update this contract if reality differs."
  ai_notes="$(printf '%b' "$ai_notes" | awk 'NF')"

  mkdir -p "$(dirname "$output")"
  {
    printf '%s\n' '---'
    printf 'issue: "%s"\n' "$(yaml_escape "$issue_number")"
    printf 'title: "%s"\n' "$(yaml_escape "$title")"
    printf 'url: "%s"\n' "$(yaml_escape "$url")"
    printf 'branch: "%s"\n' "$(yaml_escape "$branch")"
    printf 'priority: "%s"\n' "$(yaml_escape "$priority")"
    printf 'milestone: "%s"\n' "$(yaml_escape "$milestone")"
    printf 'assignee: "%s"\n' "$(yaml_escape "$assignee")"
    printf 'generated_at: "%s"\n' "$(yaml_escape "$generated_at")"
    printf 'issue_content_hash: "%s"\n' "$issue_hash"
    printf '%s\n\n' '---'
    printf '# MISSION-%s: %s\n\n' "$issue_number" "$title"
    printf '## Objective\n%s\n\n' "$objective"
    printf '## Sub-tasks\n'
    if [[ -n "$(cwork_trim "$subtasks")" ]]; then
      printf '%s\n' "$subtasks" | sed 's/^/- [ ] /'
    else
      printf '%s\n' '- [ ] Break the implementation into concrete steps before coding.'
    fi
    printf '\n## Success Criteria\n'
    if [[ -n "$(cwork_trim "$success")" ]]; then
      printf '%s\n' "$success" | sed 's/^/- [ ] /'
    else
      printf '%s\n' '- [ ] Add at least one concrete verification outcome before shipping.'
    fi
    printf '\n## Context & References\n'
    if [[ -n "$(cwork_trim "$context")" ]]; then
      printf '%s\n' "$context"
    else
      printf '%s\n' '- Review the issue thread and the relevant code paths before implementing.'
    fi
    printf '\n## Verification\n'
    if [[ -n "$(cwork_trim "$verification")" ]]; then
      printf '%s\n' "$verification"
    else
      printf '%s\n' '- Define and run the smallest verification command that proves the change.'
    fi
    printf '\n## AI Notes\n'
    printf '%s\n' "$ai_notes"
  } >"$output"
}

cmd_from_issue() {
  local issue="${1:-}"
  shift || true
  [[ -n "$issue" ]] || usage
  cwork_require_cmd gh
  cwork_require_cmd jq
  cwork_require_cmd shasum
  local output="" branch=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --output) output="$2"; shift 2 ;;
      --branch) branch="$2"; shift 2 ;;
      *) usage ;;
    esac
  done
  [[ -n "$output" ]] || output="$(cwork_contract_path "$issue")"
  local json_file
  json_file="$(mktemp)"
  trap 'rm -f "${json_file:-}"' RETURN
  gh issue view "$issue" --json number,title,body,url,labels,milestone,assignees >"$json_file"
  generate_contract "$json_file" "$output" "$branch"
  printf '%s\n' "$output"
}

cmd_from_json() {
  local json_path="${1:-}"
  shift || true
  [[ -n "$json_path" ]] || usage
  cwork_require_cmd jq
  cwork_require_cmd shasum
  local output="" branch=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --output) output="$2"; shift 2 ;;
      --branch) branch="$2"; shift 2 ;;
      *) usage ;;
    esac
  done
  local issue
  issue="$(jq -r '.number' "$json_path")"
  [[ -n "$output" ]] || output="$(cwork_contract_path "$issue")"
  generate_contract "$json_path" "$output" "$branch"
  printf '%s\n' "$output"
}

cmd_read_field() {
  local path="${1:-}" field="${2:-}"
  [[ -n "$path" && -n "$field" ]] || usage
  read_field "$path" "$field"
}

cmd_section() {
  local path="${1:-}"
  shift || true
  [[ -n "$path" && $# -gt 0 ]] || usage
  section_text "$path" "$@"
}

cmd_checked_items() {
  local path="${1:-}"
  [[ -n "$path" ]] || usage
  checked_items "$path"
}

cmd_active_path() {
  local issue="${1:-}"
  cwork_find_contract "$issue"
}

cmd_check_freshness() {
  local path="${1:-}" issue="${2:-}"
  [[ -n "$path" && -n "$issue" ]] || usage
  local stored_hash issue_json current_hash
  stored_hash="$(read_field "$path" issue_content_hash)"
  [[ -n "$stored_hash" ]] || {
    printf 'NO_HASH\n'
    exit 2
  }
  issue_json="$(gh issue view "$issue" --json title,body 2>/dev/null || true)"
  [[ -n "$issue_json" ]] || {
    printf 'NETWORK_ERROR\n'
    exit 4
  }
  current_hash="$(printf '%s%s' "$(printf '%s' "$issue_json" | jq -r '.title // ""')" "$(printf '%s' "$issue_json" | jq -r '.body // ""')" | shasum -a 256 | awk '{print $1}')"
  if [[ "$stored_hash" == "$current_hash" ]]; then
    printf 'FRESH\n'
    return 0
  fi
  printf 'STALE\n'
  exit 1
}

cmd_sync_checked() {
  local path="${1:-}" issue="${2:-}"
  [[ -n "$path" && -n "$issue" ]] || usage
  local issue_body tmp_issue tmp_checked
  issue_body="$(gh issue view "$issue" --json body --jq '.body' 2>/dev/null || true)"
  [[ -n "$issue_body" ]] || {
    echo "ERROR: could not read issue body" >&2
    exit 1
  }
  tmp_issue="$(mktemp)"
  tmp_checked="$(mktemp)"
  trap 'rm -f "${tmp_issue:-}" "${tmp_checked:-}"' RETURN
  printf '%s' "$issue_body" >"$tmp_issue"
  checked_items "$path" | while IFS= read -r item; do
    [[ -n "$item" ]] || continue
    printf '%s\t%s\n' "$(normalize_item "$item")" "$item"
  done >"$tmp_checked"

  awk -F '\t' '
    NR == FNR { checked[$1] = 1; next }
    {
      line = $0
      normalized = line
      sub(/^- \[ \] /, "", normalized)
      gsub(/[[:space:]]*[—-][[:space:]]*[0-9][0-9]*:[0-9][0-9]*$/, "", normalized)
      sub(/^🔧[[:space:]]*MANUAL:[[:space:]]*/, "", normalized)
      normalized = tolower(normalized)
      gsub(/[[:space:]]+/, " ", normalized)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", normalized)
      if (line ~ /^- \[ \] / && (normalized in checked)) {
        sub(/^- \[ \] /, "- [x] ", line)
      }
      print line
    }
  ' "$tmp_checked" "$tmp_issue" >"${tmp_issue}.new"

  if ! cmp -s "$tmp_issue" "${tmp_issue}.new"; then
    gh issue edit "$issue" --body-file "${tmp_issue}.new" >/dev/null
    printf 'SYNCED %s\n' "$issue"
  else
    printf 'UNCHANGED %s\n' "$issue"
  fi
  rm -f "${tmp_issue}.new"
}

cmd_pr_body() {
  local path="${1:-}"
  [[ -n "$path" ]] || usage
  local issue title objective success verification notes unchecked reference
  issue="$(read_field "$path" issue)"
  title="$(read_field "$path" title)"
  objective="$(section_text "$path" "Objective" || true)"
  success="$(section_text "$path" "Success Criteria" || true)"
  verification="$(section_text "$path" "Verification" || true)"
  notes="$(section_text "$path" "AI Notes" || true)"
  unchecked="$(grep -c '^- \[ \] ' "$path" || true)"
  if [[ "${unchecked:-0}" -eq 0 ]]; then
    reference="Closes #$issue"
  else
    reference="Refs #$issue"
  fi
  printf '## Objective\n%s\n\n' "${objective:-$title}"
  printf '## Success Criteria\n%s\n\n' "${success:-"- [ ] Review the contract before merge."}"
  printf '## Verification\n%s\n\n' "${verification:-"- Verification was not recorded in the contract."}"
  if [[ -n "$(cwork_trim "$notes")" ]]; then
    printf '## AI Notes\n%s\n\n' "$notes"
  fi
  printf '%s\n' "$reference"
}

cmd_delete() {
  local path="${1:-}"
  [[ -n "$path" ]] || usage
  rm -f "$path"
  printf 'DELETED %s\n' "$path"
}

subcommand="${1:-}"
shift || true

case "$subcommand" in
  from-issue) cmd_from_issue "$@" ;;
  from-json) cmd_from_json "$@" ;;
  read-field) cmd_read_field "$@" ;;
  section) cmd_section "$@" ;;
  checked-items) cmd_checked_items "$@" ;;
  active-path) cmd_active_path "$@" ;;
  check-freshness) cmd_check_freshness "$@" ;;
  sync-checked) cmd_sync_checked "$@" ;;
  pr-body) cmd_pr_body "$@" ;;
  delete) cmd_delete "$@" ;;
  *) usage ;;
esac
