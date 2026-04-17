#!/usr/bin/env bash
# tw-fleet.sh — Fleet orchestrator helpers for multi-mission teamwork
#
# USAGE:
#   bash tw-fleet.sh new-session-id                        # Generate fleet session id (timestamp-based)
#   bash tw-fleet.sh parse-issues "#A #B #C"               # Normalize issue list to newline-separated numbers
#   bash tw-fleet.sh spawn-threshold FILES LINES XMODULE   # 'spawn' or 'direct' based on ≥3/≥30/x-mod thresholds
#   bash tw-fleet.sh session-dir SESSION_ID                # Print fleet session directory path (under .teamwork/fleet/)
#   bash tw-fleet.sh init-session SESSION_ID               # Create the session directory + metadata.yml
#   bash tw-fleet.sh persist-checkpoint SESSION MISSION TYPE CONTENT_PATH
#                                                           # Store checkpoint deliverable, returns stored path
#   bash tw-fleet.sh list-checkpoints SESSION_ID           # Print all checkpoints for a session (MISSION:TYPE:PATH)
#   bash tw-fleet.sh show-checkpoint SESSION MISSION TYPE  # Cat a specific checkpoint content
#   bash tw-fleet.sh persist-approval SESSION MISSION DECISION MSG_PATH
#                                                           # Store approval decision + rationale
#   bash tw-fleet.sh render-briefing TEMPLATE_PATH VARS_FILE  # Render {{PLACEHOLDER}} from a VAR=value vars file
#   bash tw-fleet.sh cleanup-session SESSION_ID            # Remove .teamwork/fleet/<session>/ entirely
#
# EXIT CODES:
#   0 — success
#   1 — usage error or missing args
#   2 — filesystem operation failed
#   3 — session/checkpoint not found
#
# WHY THIS SCRIPT EXISTS:
# Per teamwork Constitution P11 (Von Neumann separation): skill .md files contain LLM-interpreted
# content (analysis, judgment, dialogue). Deterministic operations (parse, threshold decide, persist,
# render template, cleanup) belong in scripts. tw-fleet.sh is the mechanism layer for /team fleet
# sub-command; the LLM-interpreted content lives in team.md's Operation Fleet section.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TW_CONFIG="$SCRIPT_DIR/tw-config.sh"

# Portable in-place sed (macOS uses -i '', GNU/Linux uses -i)
if sed --version 2>/dev/null | grep -q GNU 2>/dev/null; then
  _sed_i() { sed -i "$@"; }
else
  _sed_i() { sed -i '' "$@"; }
fi

_fleet_root() {
  # Fleet state lives under the teamwork dir (.teamwork/fleet/ or .teamspace/fleet/).
  # Follows the same worktree-aware detection as other scripts.
  local dir
  dir=$(bash "$TW_CONFIG" detect-dir 2>/dev/null)
  if [ -z "$dir" ]; then
    echo "ERROR: No teamwork config found (looked for .teamwork/ or .teamspace/)." >&2
    exit 2
  fi
  echo "$dir/fleet"
}

usage() {
  echo "Usage: tw-fleet.sh <subcommand> [args...]" >&2
  echo "Subcommands: new-session-id, parse-issues, spawn-threshold," >&2
  echo "             session-dir, init-session, persist-checkpoint," >&2
  echo "             list-checkpoints, show-checkpoint, persist-approval," >&2
  echo "             render-briefing, cleanup-session" >&2
  exit 1
}

# ───────────────────────────────────────────────────────────────────
# new-session-id — generate timestamp-based session id
# ───────────────────────────────────────────────────────────────────
cmd_new_session_id() {
  # Format: fleet-YYYYMMDD-HHMMSS (UTC for reproducibility)
  date -u +"fleet-%Y%m%d-%H%M%S"
}

# ───────────────────────────────────────────────────────────────────
# parse-issues — take "#A #B #C" (or "A,B,C" or mixed) → newline-separated bare numbers
# ───────────────────────────────────────────────────────────────────
cmd_parse_issues() {
  local input="${1:-}"
  if [ -z "$input" ]; then
    echo "ERROR: parse-issues requires a string like '#1298 #1299 #1302'" >&2
    exit 1
  fi
  # Extract runs of digits; dedupe while preserving order; reject duplicates silently
  echo "$input" \
    | tr ',' ' ' \
    | tr -s ' ' '\n' \
    | sed 's/^#//' \
    | grep -E '^[0-9]+$' \
    | awk '!seen[$0]++'
}

# ───────────────────────────────────────────────────────────────────
# spawn-threshold — decide spawn teammate vs team-lead direct
# Thresholds from session "teammate-work-2" wisdom:
#   ≥3 files OR ≥30 lines OR cross-module → spawn
#   else → team-lead direct (but worktree still required)
# ───────────────────────────────────────────────────────────────────
cmd_spawn_threshold() {
  local files="${1:-}" lines="${2:-}" cross_module="${3:-}"
  if [ -z "$files" ] || [ -z "$lines" ] || [ -z "$cross_module" ]; then
    echo "ERROR: spawn-threshold requires FILES LINES CROSS_MODULE (cross_module: yes/no)" >&2
    exit 1
  fi
  # Validate numerics
  if ! echo "$files" | grep -qE '^[0-9]+$' || ! echo "$lines" | grep -qE '^[0-9]+$'; then
    echo "ERROR: FILES and LINES must be non-negative integers" >&2
    exit 1
  fi
  if [ "$files" -ge 3 ] || [ "$lines" -ge 30 ] || [ "$cross_module" = "yes" ]; then
    echo "spawn"
  else
    echo "direct"
  fi
}

# ───────────────────────────────────────────────────────────────────
# session-dir — path computation (no side effects)
# ───────────────────────────────────────────────────────────────────
cmd_session_dir() {
  local session_id="${1:-}"
  if [ -z "$session_id" ]; then
    echo "ERROR: session-dir requires SESSION_ID" >&2
    exit 1
  fi
  local root
  root=$(_fleet_root)
  echo "$root/$session_id"
}

# ───────────────────────────────────────────────────────────────────
# init-session — create the session directory and metadata file
# ───────────────────────────────────────────────────────────────────
cmd_init_session() {
  local session_id="${1:-}"
  if [ -z "$session_id" ]; then
    echo "ERROR: init-session requires SESSION_ID" >&2
    exit 1
  fi
  local dir
  dir=$(cmd_session_dir "$session_id")
  mkdir -p "$dir/checkpoints" "$dir/approvals" 2>/dev/null || {
    echo "ERROR: Could not create $dir" >&2
    exit 2
  }
  cat > "$dir/metadata.yml" <<EOF
session_id: $session_id
created_at: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
schema: "fleet-v1"
EOF
  echo "$dir"
}

# ───────────────────────────────────────────────────────────────────
# persist-checkpoint — store teammate checkpoint deliverable
# ───────────────────────────────────────────────────────────────────
cmd_persist_checkpoint() {
  local session_id="${1:-}" mission="${2:-}" cp_type="${3:-}" content_path="${4:-}"
  if [ -z "$session_id" ] || [ -z "$mission" ] || [ -z "$cp_type" ] || [ -z "$content_path" ]; then
    echo "ERROR: persist-checkpoint requires SESSION MISSION TYPE CONTENT_PATH" >&2
    exit 1
  fi
  if [ ! -f "$content_path" ]; then
    echo "ERROR: content path not found: $content_path" >&2
    exit 3
  fi
  local dir
  dir=$(cmd_session_dir "$session_id")
  if [ ! -d "$dir" ]; then
    echo "ERROR: session dir not found: $dir (run init-session first)" >&2
    exit 3
  fi
  local dest="$dir/checkpoints/mission-${mission}-${cp_type}.md"
  cp "$content_path" "$dest" || {
    echo "ERROR: Could not copy checkpoint content" >&2
    exit 2
  }
  echo "$dest"
}

# ───────────────────────────────────────────────────────────────────
# list-checkpoints — print MISSION:TYPE:PATH lines
# ───────────────────────────────────────────────────────────────────
cmd_list_checkpoints() {
  local session_id="${1:-}"
  if [ -z "$session_id" ]; then
    echo "ERROR: list-checkpoints requires SESSION_ID" >&2
    exit 1
  fi
  local dir
  dir=$(cmd_session_dir "$session_id")
  if [ ! -d "$dir/checkpoints" ]; then
    # No checkpoints yet — empty is legitimate (not error)
    return 0
  fi
  # Parse mission-<N>-<type>.md filenames back into MISSION:TYPE:PATH
  for f in "$dir/checkpoints"/mission-*.md; do
    [ -e "$f" ] || continue
    local base="${f##*/}"                # mission-1298-design-note.md
    base="${base%.md}"                   # mission-1298-design-note
    local body="${base#mission-}"        # 1298-design-note
    local mission="${body%%-*}"          # 1298
    local cp_type="${body#*-}"           # design-note
    printf '%s:%s:%s\n' "$mission" "$cp_type" "$f"
  done
}

# ───────────────────────────────────────────────────────────────────
# show-checkpoint — cat a specific mission/type checkpoint
# ───────────────────────────────────────────────────────────────────
cmd_show_checkpoint() {
  local session_id="${1:-}" mission="${2:-}" cp_type="${3:-}"
  if [ -z "$session_id" ] || [ -z "$mission" ] || [ -z "$cp_type" ]; then
    echo "ERROR: show-checkpoint requires SESSION MISSION TYPE" >&2
    exit 1
  fi
  local dir f
  dir=$(cmd_session_dir "$session_id")
  f="$dir/checkpoints/mission-${mission}-${cp_type}.md"
  if [ ! -f "$f" ]; then
    echo "ERROR: checkpoint not found: $f" >&2
    exit 3
  fi
  cat "$f"
}

# ───────────────────────────────────────────────────────────────────
# persist-approval — record team-lead decision + rationale
# ───────────────────────────────────────────────────────────────────
cmd_persist_approval() {
  local session_id="${1:-}" mission="${2:-}" decision="${3:-}" msg_path="${4:-}"
  if [ -z "$session_id" ] || [ -z "$mission" ] || [ -z "$decision" ] || [ -z "$msg_path" ]; then
    echo "ERROR: persist-approval requires SESSION MISSION DECISION MSG_PATH" >&2
    exit 1
  fi
  case "$decision" in
    approve|revise|escalate) ;;
    *)
      echo "ERROR: DECISION must be one of: approve | revise | escalate" >&2
      exit 1
      ;;
  esac
  if [ ! -f "$msg_path" ]; then
    echo "ERROR: msg path not found: $msg_path" >&2
    exit 3
  fi
  local dir dest
  dir=$(cmd_session_dir "$session_id")
  if [ ! -d "$dir" ]; then
    echo "ERROR: session dir not found: $dir" >&2
    exit 3
  fi
  mkdir -p "$dir/approvals" 2>/dev/null || true
  local ts
  ts=$(date -u +"%Y%m%dT%H%M%SZ")
  dest="$dir/approvals/mission-${mission}-${ts}-${decision}.md"
  cp "$msg_path" "$dest" || {
    echo "ERROR: Could not copy approval message" >&2
    exit 2
  }
  echo "$dest"
}

# ───────────────────────────────────────────────────────────────────
# render-briefing — {{PLACEHOLDER}} substitution from VAR=value file
# ───────────────────────────────────────────────────────────────────
cmd_render_briefing() {
  local template="${1:-}" vars_file="${2:-}"
  if [ -z "$template" ] || [ -z "$vars_file" ]; then
    echo "ERROR: render-briefing requires TEMPLATE_PATH VARS_FILE" >&2
    exit 1
  fi
  if [ ! -f "$template" ]; then
    echo "ERROR: template not found: $template" >&2
    exit 3
  fi
  if [ ! -f "$vars_file" ]; then
    echo "ERROR: vars file not found: $vars_file" >&2
    exit 3
  fi
  # Read template into memory, then substitute each VAR=value line from vars_file.
  # vars_file format: one KEY=VALUE per line, blank lines and # comments ignored.
  local rendered
  rendered=$(cat "$template")
  while IFS= read -r line; do
    # Skip blanks and comments
    case "$line" in
      ''|\#*) continue ;;
    esac
    local key="${line%%=*}"
    local value="${line#*=}"
    if [ -z "$key" ] || [ "$key" = "$line" ]; then
      # Malformed line — skip (no '=')
      continue
    fi
    # Escape value for sed: & and / and backslashes and newlines
    # Using a delimiter unlikely to appear: \x01
    local esc_value
    esc_value=$(printf '%s' "$value" | sed -e 's/[\/&]/\\&/g')
    rendered=$(printf '%s' "$rendered" | sed "s/{{${key}}}/${esc_value}/g")
  done < "$vars_file"
  printf '%s\n' "$rendered"
}

# ───────────────────────────────────────────────────────────────────
# cleanup-session — remove an entire fleet session dir (after successful handoff)
# ───────────────────────────────────────────────────────────────────
cmd_cleanup_session() {
  local session_id="${1:-}"
  if [ -z "$session_id" ]; then
    echo "ERROR: cleanup-session requires SESSION_ID" >&2
    exit 1
  fi
  local dir
  dir=$(cmd_session_dir "$session_id")
  if [ ! -d "$dir" ]; then
    # Already gone — silently succeed (idempotent)
    return 0
  fi
  rm -rf "$dir" || {
    echo "ERROR: Could not remove $dir" >&2
    exit 2
  }
  echo "Removed $dir"
}

# ───────────────────────────────────────────────────────────────────
# Main
# ───────────────────────────────────────────────────────────────────
SUBCOMMAND="${1:-}"
shift || true

case "$SUBCOMMAND" in
  new-session-id)        cmd_new_session_id ;;
  parse-issues)          cmd_parse_issues "$@" ;;
  spawn-threshold)       cmd_spawn_threshold "$@" ;;
  session-dir)           cmd_session_dir "$@" ;;
  init-session)          cmd_init_session "$@" ;;
  persist-checkpoint)    cmd_persist_checkpoint "$@" ;;
  list-checkpoints)      cmd_list_checkpoints "$@" ;;
  show-checkpoint)       cmd_show_checkpoint "$@" ;;
  persist-approval)      cmd_persist_approval "$@" ;;
  render-briefing)       cmd_render_briefing "$@" ;;
  cleanup-session)       cmd_cleanup_session "$@" ;;
  *)                     usage ;;
esac
