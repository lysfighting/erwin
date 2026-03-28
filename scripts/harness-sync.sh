#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_ROOT="$(git -C "$SCRIPT_DIR/.." rev-parse --show-toplevel 2>/dev/null || cd "$SCRIPT_DIR/.." && pwd)"
SOURCE_ROOT="${HARNESS_SOURCE_ROOT:-}"
MODE="sync"
FORCE=1
TARGET_INPUT="."

usage() {
  cat <<'EOF'
Usage:
  ./scripts/harness-sync.sh
  ./scripts/harness-sync.sh status
  ./scripts/harness-sync.sh --source /path/to/source-repo
  ./scripts/harness-sync.sh --source /path/to/source-repo /path/to/target-repo

Sync a bootstrapped harness adapter from its source repository.
If --source is omitted, the script uses .harness/cache/adapter-source.json.
EOF
}

adapter_query() {
  local root="$1"
  shift
  "$root/scripts/harness-runtime.sh" adapter "$@"
}

resolve_metadata_path() {
  local root="$1"
  local rel_path=""
  if [[ -f "$root/scripts/harness-runtime.sh" ]]; then
    rel_path="$(adapter_query "$root" metadata-path 2>/dev/null || true)"
  fi
  if [[ -n "$rel_path" ]]; then
    printf '%s\n' "$root/$rel_path"
  else
    printf '%s\n' "$root/.harness/cache/adapter-source.json"
  fi
}

read_metadata_field() {
  local field="$1"
  local path="$2"
  if [[ ! -f "$path" ]]; then
    printf '\n'
    return 0
  fi
  sed -n "s/^[[:space:]]*\"$field\":[[:space:]]*\"\(.*\)\"[[:space:]]*,\{0,1\}[[:space:]]*$/\1/p" "$path"
}

read_manifest_version() {
  local root="$1"
  if [[ -f "$root/scripts/harness-runtime.sh" ]]; then
    adapter_query "$root" version 2>/dev/null || true
  else
    printf '\n'
  fi
}

read_manifest_id() {
  local root="$1"
  if [[ -f "$root/scripts/harness-runtime.sh" ]]; then
    adapter_query "$root" id 2>/dev/null || true
  else
    printf '\n'
  fi
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    status)
      MODE="status"
      shift
      ;;
    --source|-s)
      if [[ -z "${2:-}" ]]; then
        echo "ERROR: --source requires a repository path" >&2
        exit 1
      fi
      SOURCE_ROOT="$2"
      shift 2
      ;;
    --force|-f)
      FORCE=1
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      TARGET_INPUT="$1"
      shift
      ;;
  esac
done

if ! TARGET_ROOT="$(git -C "$TARGET_INPUT" rev-parse --show-toplevel 2>/dev/null)"; then
  echo "ERROR: target must be inside a git repository: $TARGET_INPUT" >&2
  exit 1
fi
TARGET_ROOT="$(cd "$TARGET_ROOT" && pwd)"
METADATA_PATH="$(resolve_metadata_path "$TARGET_ROOT")"

if [[ "$MODE" == "status" ]]; then
  if [[ -f "$METADATA_PATH" ]]; then
    source_root="$(read_metadata_field source_root "$METADATA_PATH")"
    target_manifest_version="$(read_manifest_version "$TARGET_ROOT")"
    target_adapter_id="$(read_manifest_id "$TARGET_ROOT")"
    source_manifest_version=""
    source_adapter_id=""
    if [[ -n "$source_root" ]] && [[ -d "$source_root" ]]; then
      source_manifest_version="$(read_manifest_version "$source_root")"
      source_adapter_id="$(read_manifest_id "$source_root")"
    fi
    adapter_ids_match=false
    manifest_versions_match=false
    if [[ -n "$target_adapter_id" && -n "$source_adapter_id" && "$target_adapter_id" == "$source_adapter_id" ]]; then
      adapter_ids_match=true
    fi
    if [[ -n "$target_manifest_version" && -n "$source_manifest_version" && "$target_manifest_version" == "$source_manifest_version" ]]; then
      manifest_versions_match=true
    fi
    cat <<EOF
{
  "source_root": "$(read_metadata_field source_root "$METADATA_PATH")",
  "source_repository": "$(read_metadata_field source_repository "$METADATA_PATH")",
  "target_repo": "$(read_metadata_field target_repo "$METADATA_PATH")",
  "profile_id": "$(read_metadata_field profile_id "$METADATA_PATH")",
  "repository": "$(read_metadata_field repository "$METADATA_PATH")",
  "base_branch": "$(read_metadata_field base_branch "$METADATA_PATH")",
  "last_synced_at": "$(read_metadata_field last_synced_at "$METADATA_PATH")",
  "target_adapter_id": "$target_adapter_id",
  "source_adapter_id": "$source_adapter_id",
  "target_manifest_version": "$target_manifest_version",
  "source_manifest_version": "$source_manifest_version",
  "adapter_ids_match": $adapter_ids_match,
  "manifest_versions_match": $manifest_versions_match
}
EOF
  else
    echo '{}'
  fi
  exit 0
fi

if [[ -z "$SOURCE_ROOT" ]] && [[ -f "$METADATA_PATH" ]]; then
  SOURCE_ROOT="$(read_metadata_field source_root "$METADATA_PATH")"
fi

if [[ -z "$SOURCE_ROOT" ]]; then
  echo "ERROR: no source configured; pass --source PATH or bootstrap first" >&2
  exit 1
fi

SOURCE_ROOT="$(cd "$SOURCE_ROOT" && pwd)"

args=()
if [[ "$FORCE" -eq 1 ]]; then
  args+=(--force)
fi
args+=("$TARGET_ROOT")

HARNESS_SOURCE_ROOT="$SOURCE_ROOT" "$TARGET_ROOT/scripts/harness-bootstrap.sh" "${args[@]}"
