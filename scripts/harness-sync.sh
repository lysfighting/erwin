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

read_manifest_version() {
  local root="$1"
  python3 - "$root" <<'PY'
import pathlib
import sys

root = pathlib.Path(sys.argv[1])
path = root / "harness" / "adapter-manifest.yml"
if not path.exists():
    print("")
    raise SystemExit(0)
text = path.read_text(encoding="utf-8")
for raw in text.splitlines():
    line = raw.strip()
    if line.startswith("adapter_version:"):
        print(line.split(":", 1)[1].strip().strip('"'))
        break
else:
    print("")
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
    source_root="$(read_metadata_field source_root "$METADATA_PATH")"
    target_manifest_version="$(read_manifest_version "$TARGET_ROOT")"
    source_manifest_version=""
    if [[ -n "$source_root" ]] && [[ -d "$source_root" ]]; then
      source_manifest_version="$(read_manifest_version "$source_root")"
    fi
    python3 - "$METADATA_PATH" "$target_manifest_version" "$source_manifest_version" <<'PY'
import json
import pathlib
import sys

metadata_path = pathlib.Path(sys.argv[1])
target_manifest_version = sys.argv[2]
source_manifest_version = sys.argv[3]
data = json.loads(metadata_path.read_text(encoding="utf-8"))
data["target_manifest_version"] = target_manifest_version
data["source_manifest_version"] = source_manifest_version
data["manifest_versions_match"] = bool(
    target_manifest_version and source_manifest_version and target_manifest_version == source_manifest_version
)
print(json.dumps(data, indent=2, ensure_ascii=True))
PY
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
