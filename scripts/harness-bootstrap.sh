#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE_ROOT_INPUT="${HARNESS_SOURCE_ROOT:-}"
if [[ -n "$SOURCE_ROOT_INPUT" ]]; then
  SOURCE_ROOT="$(cd "$SOURCE_ROOT_INPUT" && pwd)"
else
  SOURCE_ROOT="$(dirname "$SCRIPT_DIR")"
fi
FORCE=0
TARGET_INPUT="."

usage() {
  cat <<'EOF'
Usage:
  ./scripts/harness-bootstrap.sh [TARGET_REPO]
  ./scripts/harness-bootstrap.sh --force [TARGET_REPO]

Bootstrap the first usable Codex-native harness adapter into a git repository.
Copies the harness substrate, committed definitions, and smoke test into TARGET_REPO.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
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

if [[ ! -f "$SOURCE_ROOT/scripts/harness_state.py" ]] || [[ ! -f "$SOURCE_ROOT/install-codex.sh" ]]; then
  echo "ERROR: source root does not look like a harness source repo: $SOURCE_ROOT" >&2
  exit 1
fi

adapter_query() {
  python3 "$SOURCE_ROOT/scripts/harness_state.py" adapter "$@"
}

if ! ADAPTER_VERSION="$(adapter_query version)"; then
  echo "ERROR: could not read adapter version from manifest" >&2
  exit 1
fi
if ! PROFILE_TEMPLATE_REL="$(adapter_query profile-template)"; then
  echo "ERROR: could not read bootstrap profile template from manifest" >&2
  exit 1
fi
if ! GENERATED_PROFILE_REL="$(adapter_query generated-profile-path)"; then
  echo "ERROR: could not read generated profile path from manifest" >&2
  exit 1
fi
if ! METADATA_REL="$(adapter_query metadata-path)"; then
  echo "ERROR: could not read adapter metadata path from manifest" >&2
  exit 1
fi
PACKAGE_FILES=()
while IFS= read -r rel_path; do
  [[ -n "$rel_path" ]] || continue
  PACKAGE_FILES+=("$rel_path")
done < <(adapter_query package-files)
if [[ "${#PACKAGE_FILES[@]}" -eq 0 ]]; then
  echo "ERROR: adapter manifest does not declare package files" >&2
  exit 1
fi

TARGET_ROOT="$(cd "$TARGET_ROOT" && pwd)"
REPO_NAME="$(basename "$TARGET_ROOT")"
SYNCED_AT="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

remote_url="$(git -C "$TARGET_ROOT" remote get-url origin 2>/dev/null || true)"
repository_slug="$remote_url"
if [[ "$remote_url" =~ ^git@github\.com:([^/]+)/(.+)\.git$ ]]; then
  repository_slug="${BASH_REMATCH[1]}/${BASH_REMATCH[2]}"
elif [[ "$remote_url" =~ ^https://github\.com/([^/]+)/(.+)\.git$ ]]; then
  repository_slug="${BASH_REMATCH[1]}/${BASH_REMATCH[2]}"
elif [[ "$remote_url" =~ ^https://github\.com/([^/]+)/(.+)$ ]]; then
  repository_slug="${BASH_REMATCH[1]}/${BASH_REMATCH[2]}"
elif [[ -z "$remote_url" ]]; then
  repository_slug="$REPO_NAME"
fi

base_branch="$(git -C "$TARGET_ROOT" symbolic-ref --quiet --short refs/remotes/origin/HEAD 2>/dev/null | sed 's#^origin/##' || true)"
if [[ -z "$base_branch" ]]; then
  base_branch="$(git -C "$TARGET_ROOT" branch --show-current 2>/dev/null || true)"
fi
if [[ -z "$base_branch" ]]; then
  base_branch="main"
fi

profile_id="${REPO_NAME}-main"

copy_file() {
  local src="$1"
  local dst="$2"

  mkdir -p "$(dirname "$dst")"
  if [[ -e "$dst" ]]; then
    if cmp -s "$src" "$dst"; then
      return 0
    fi
    if [[ "$FORCE" -ne 1 ]]; then
      echo "ERROR: refusing to overwrite existing file without --force: $dst" >&2
      exit 1
    fi
  fi
  cp "$src" "$dst"
}

generate_profile() {
  local template="$SOURCE_ROOT/$PROFILE_TEMPLATE_REL"
  local output="$TARGET_ROOT/$GENERATED_PROFILE_REL"
  mkdir -p "$(dirname "$output")"
  if [[ -e "$output" && "$FORCE" -ne 1 ]]; then
    echo "ERROR: refusing to overwrite existing file without --force: $output" >&2
    exit 1
  fi
  sed \
    -e "s|__PROFILE_ID__|$profile_id|g" \
    -e "s|__NAME__|$REPO_NAME|g" \
    -e "s|__REPOSITORY__|$repository_slug|g" \
    -e "s|__BASE_BRANCH__|$base_branch|g" \
    "$template" > "$output"
}

ensure_gitignore_entry() {
  local gitignore_path="$TARGET_ROOT/.gitignore"
  if [[ ! -f "$gitignore_path" ]]; then
    printf ".harness/\n" > "$gitignore_path"
    return 0
  fi
  if ! grep -qxF ".harness/" "$gitignore_path"; then
    printf "\n.harness/\n" >> "$gitignore_path"
  fi
}

write_adapter_metadata() {
  local metadata_path="$TARGET_ROOT/$METADATA_REL"
  local metadata_dir
  metadata_dir="$(dirname "$metadata_path")"
  mkdir -p "$metadata_dir"
  cat > "$metadata_path" <<EOF
{
  "adapter_version": "$ADAPTER_VERSION",
  "source_root": "$SOURCE_ROOT",
  "source_repository": "$(basename "$SOURCE_ROOT")",
  "target_repo": "$TARGET_ROOT",
  "profile_id": "$profile_id",
  "repository": "$repository_slug",
  "base_branch": "$base_branch",
  "last_synced_at": "$SYNCED_AT"
}
EOF
}

for rel_path in "${PACKAGE_FILES[@]}"; do
  copy_file "$SOURCE_ROOT/$rel_path" "$TARGET_ROOT/$rel_path"
done

generate_profile
ensure_gitignore_entry
write_adapter_metadata

chmod_targets=()
for rel_path in "${PACKAGE_FILES[@]}"; do
  if [[ "$rel_path" == *.sh ]]; then
    chmod_targets+=("$TARGET_ROOT/$rel_path")
  fi
done
if [[ "${#chmod_targets[@]}" -gt 0 ]]; then
  chmod +x "${chmod_targets[@]}"
fi

installed_json="$(printf '%s\n' "${PACKAGE_FILES[@]}" "$GENERATED_PROFILE_REL" | python3 -c 'import json,sys; print(json.dumps([line.strip() for line in sys.stdin if line.strip()], ensure_ascii=True))')"
cat <<EOF
{
  "adapter_version": "$ADAPTER_VERSION",
  "source_root": "$SOURCE_ROOT",
  "target_repo": "$TARGET_ROOT",
  "profile_id": "$profile_id",
  "repository": "$repository_slug",
  "base_branch": "$base_branch",
  "metadata_path": "$TARGET_ROOT/$METADATA_REL",
  "installed": $installed_json
}
EOF
