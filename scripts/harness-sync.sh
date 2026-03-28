#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_ROOT="$(git -C "$SCRIPT_DIR/.." rev-parse --show-toplevel 2>/dev/null || cd "$SCRIPT_DIR/.." && pwd)"
METADATA_PATH="$TARGET_ROOT/.harness/cache/adapter-source.json"
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

read_metadata_field() {
  local field="$1"
  local path="$2"
  python3 - "$field" "$path" <<'PY'
import json
import pathlib
import sys

field = sys.argv[1]
path = pathlib.Path(sys.argv[2])
if not path.exists():
    print("")
    raise SystemExit(0)
data = json.loads(path.read_text(encoding="utf-8"))
value = data.get(field, "")
if isinstance(value, (dict, list)):
    print(json.dumps(value, ensure_ascii=True))
else:
    print(value)
PY
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
METADATA_PATH="$TARGET_ROOT/.harness/cache/adapter-source.json"

if [[ "$MODE" == "status" ]]; then
  if [[ -f "$METADATA_PATH" ]]; then
    cat "$METADATA_PATH"
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
