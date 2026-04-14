#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/cwork-common.sh"

usage() {
  echo "Usage: cwork-config.sh <subcommand> [args...]" >&2
  echo "Subcommands: detect-dir, get, resolve-labels" >&2
  exit 1
}

detect_dir() {
  local main_root
  main_root="$(cwork_main_repo_root)"
  if [[ -f "$main_root/.teamwork/config.yml" ]]; then
    printf '%s/.teamwork\n' "$main_root"
    return 0
  fi
  if [[ -f "$main_root/.teamspace/config.yml" ]]; then
    printf '%s/.teamspace\n' "$main_root"
    return 0
  fi
  return 1
}

config_file() {
  local dir
  dir="$(detect_dir 2>/dev/null || true)"
  if [[ -n "$dir" ]]; then
    printf '%s/config.yml\n' "$dir"
  fi
  return 0
}

yaml_get() {
  local file="$1"
  local key="$2"
  local default="${3:-}"
  awk -v target="$key" -v fallback="$default" '
    function trim(s) {
      sub(/^[[:space:]]+/, "", s)
      sub(/[[:space:]]+$/, "", s)
      return s
    }
    function strip_val(raw, v, q, pos) {
      v = trim(raw)
      if (v == "") return ""
      q = substr(v, 1, 1)
      if (q == "\"" || q == "'\''") {
        v = substr(v, 2)
        pos = index(v, q)
        if (pos > 0) return substr(v, 1, pos - 1)
        return v
      }
      sub(/[[:space:]]+#.*$/, "", v)
      return trim(v)
    }
    {
      line = $0
      if (line ~ /^[[:space:]]*#/ || line ~ /^[[:space:]]*$/) next
      if (line !~ /^[[:space:]]/ && index(line, ":") > 0) {
        key = trim(substr(line, 1, index(line, ":") - 1))
        val = strip_val(substr(line, index(line, ":") + 1))
        current = key
        if (val != "" && val != "|" && val != ">") data[key] = val
        next
      }
      if (line ~ /^  [^[:space:]-][^:]*:/ && current != "") {
        trimmed = trim(line)
        subkey = trim(substr(trimmed, 1, index(trimmed, ":") - 1))
        val = strip_val(substr(trimmed, index(trimmed, ":") + 1))
        if (val != "" && val != "|" && val != ">") data[current "." subkey] = val
      }
    }
    END {
      if (target in data) print data[target]
      else print fallback
    }
  ' "$file"
}

cmd_detect_dir() {
  local dir
  dir="$(detect_dir 2>/dev/null || true)"
  [[ -n "$dir" ]] || exit 3
  printf '%s\n' "$dir"
}

cmd_get() {
  local key="${1:-}"
  local default="${2:-}"
  [[ -n "$key" ]] || usage
  local file
  file="$(config_file)"
  if [[ -z "$file" || ! -f "$file" ]]; then
    printf '%s\n' "$default"
    return 0
  fi
  yaml_get "$file" "$key" "$default"
}

cmd_resolve_labels() {
  local mission_label status_prefix priority_prefix
  mission_label="$(cmd_get github.mc_label "")"
  [[ -n "$mission_label" ]] || mission_label="$(cmd_get mc_label "")"
  [[ -n "$mission_label" ]] || mission_label="$(cmd_get labels.mission "mission")"

  status_prefix="$(cmd_get label_prefix.status "")"
  [[ -n "$status_prefix" ]] || status_prefix="$(cmd_get labels.status_prefix "status:")"

  priority_prefix="$(cmd_get label_prefix.priority "")"
  [[ -n "$priority_prefix" ]] || priority_prefix="$(cmd_get labels.priority_prefix "priority:")"

  printf 'MISSION_LABEL="%s"\n' "$mission_label"
  printf 'STATUS_PREFIX="%s"\n' "$status_prefix"
  printf 'PRIORITY_PREFIX="%s"\n' "$priority_prefix"
}

subcommand="${1:-}"
shift || true

case "$subcommand" in
  detect-dir) cmd_detect_dir "$@" ;;
  get) cmd_get "$@" ;;
  resolve-labels) cmd_resolve_labels "$@" ;;
  *) usage ;;
esac
