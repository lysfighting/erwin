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
ADAPTER_VERSION="0.2.0-prototype"

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
  local template="$SOURCE_ROOT/harness/bootstrap/profile-template.yml"
  local output="$TARGET_ROOT/harness/profiles/default.yml"
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
  local metadata_dir="$TARGET_ROOT/.harness/cache"
  local metadata_path="$metadata_dir/adapter-source.json"
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

copy_file "$SOURCE_ROOT/scripts/harness_state.py" "$TARGET_ROOT/scripts/harness_state.py"
copy_file "$SOURCE_ROOT/scripts/harness-home.sh" "$TARGET_ROOT/scripts/harness-home.sh"
copy_file "$SOURCE_ROOT/scripts/harness-doctor.sh" "$TARGET_ROOT/scripts/harness-doctor.sh"
copy_file "$SOURCE_ROOT/scripts/harness-issue.sh" "$TARGET_ROOT/scripts/harness-issue.sh"
copy_file "$SOURCE_ROOT/scripts/harness-claim.sh" "$TARGET_ROOT/scripts/harness-claim.sh"
copy_file "$SOURCE_ROOT/scripts/harness-drive.sh" "$TARGET_ROOT/scripts/harness-drive.sh"
copy_file "$SOURCE_ROOT/scripts/harness-ship.sh" "$TARGET_ROOT/scripts/harness-ship.sh"
copy_file "$SOURCE_ROOT/scripts/harness-rc.sh" "$TARGET_ROOT/scripts/harness-rc.sh"
copy_file "$SOURCE_ROOT/scripts/harness-meetings.sh" "$TARGET_ROOT/scripts/harness-meetings.sh"
copy_file "$SOURCE_ROOT/scripts/harness-maintenance.sh" "$TARGET_ROOT/scripts/harness-maintenance.sh"
copy_file "$SOURCE_ROOT/scripts/harness-profile.sh" "$TARGET_ROOT/scripts/harness-profile.sh"
copy_file "$SOURCE_ROOT/scripts/harness-paths.sh" "$TARGET_ROOT/scripts/harness-paths.sh"
copy_file "$SOURCE_ROOT/scripts/harness-runs.sh" "$TARGET_ROOT/scripts/harness-runs.sh"
copy_file "$SOURCE_ROOT/scripts/harness-signals.sh" "$TARGET_ROOT/scripts/harness-signals.sh"
copy_file "$SOURCE_ROOT/scripts/harness-contracts.sh" "$TARGET_ROOT/scripts/harness-contracts.sh"
copy_file "$SOURCE_ROOT/scripts/harness-git.sh" "$TARGET_ROOT/scripts/harness-git.sh"
copy_file "$SOURCE_ROOT/scripts/harness-bootstrap.sh" "$TARGET_ROOT/scripts/harness-bootstrap.sh"
copy_file "$SOURCE_ROOT/scripts/harness-sync.sh" "$TARGET_ROOT/scripts/harness-sync.sh"
copy_file "$SOURCE_ROOT/scripts/harness-issue-smoke-test.sh" "$TARGET_ROOT/scripts/harness-issue-smoke-test.sh"
copy_file "$SOURCE_ROOT/scripts/harness-smoke-test.sh" "$TARGET_ROOT/scripts/harness-smoke-test.sh"
copy_file "$SOURCE_ROOT/scripts/harness-release-smoke-test.sh" "$TARGET_ROOT/scripts/harness-release-smoke-test.sh"
copy_file "$SOURCE_ROOT/install-codex.sh" "$TARGET_ROOT/install-codex.sh"

copy_file "$SOURCE_ROOT/harness/README.md" "$TARGET_ROOT/harness/README.md"
copy_file "$SOURCE_ROOT/harness/templates/contract.yml" "$TARGET_ROOT/harness/templates/contract.yml"
copy_file "$SOURCE_ROOT/harness/templates/issue.md" "$TARGET_ROOT/harness/templates/issue.md"
copy_file "$SOURCE_ROOT/harness/templates/pr.md" "$TARGET_ROOT/harness/templates/pr.md"
copy_file "$SOURCE_ROOT/harness/templates/run.yml" "$TARGET_ROOT/harness/templates/run.yml"
copy_file "$SOURCE_ROOT/harness/templates/signal.yml" "$TARGET_ROOT/harness/templates/signal.yml"
copy_file "$SOURCE_ROOT/harness/templates/meeting.yml" "$TARGET_ROOT/harness/templates/meeting.yml"
copy_file "$SOURCE_ROOT/harness/templates/rc.yml" "$TARGET_ROOT/harness/templates/rc.yml"
copy_file "$SOURCE_ROOT/harness/policies/defaults.yml" "$TARGET_ROOT/harness/policies/defaults.yml"
copy_file "$SOURCE_ROOT/skills/harness/SKILL.md" "$TARGET_ROOT/skills/harness/SKILL.md"
copy_file "$SOURCE_ROOT/skills/harness/README.md" "$TARGET_ROOT/skills/harness/README.md"
copy_file "$SOURCE_ROOT/skills/harness/CHANGELOG.md" "$TARGET_ROOT/skills/harness/CHANGELOG.md"
generate_profile
ensure_gitignore_entry
write_adapter_metadata

chmod +x \
  "$TARGET_ROOT/install-codex.sh" \
  "$TARGET_ROOT/scripts/harness-home.sh" \
  "$TARGET_ROOT/scripts/harness-doctor.sh" \
  "$TARGET_ROOT/scripts/harness-issue.sh" \
  "$TARGET_ROOT/scripts/harness-claim.sh" \
  "$TARGET_ROOT/scripts/harness-drive.sh" \
  "$TARGET_ROOT/scripts/harness-ship.sh" \
  "$TARGET_ROOT/scripts/harness-rc.sh" \
  "$TARGET_ROOT/scripts/harness-meetings.sh" \
  "$TARGET_ROOT/scripts/harness-maintenance.sh" \
  "$TARGET_ROOT/scripts/harness-profile.sh" \
  "$TARGET_ROOT/scripts/harness-paths.sh" \
  "$TARGET_ROOT/scripts/harness-runs.sh" \
  "$TARGET_ROOT/scripts/harness-signals.sh" \
  "$TARGET_ROOT/scripts/harness-contracts.sh" \
  "$TARGET_ROOT/scripts/harness-git.sh" \
  "$TARGET_ROOT/scripts/harness-bootstrap.sh" \
  "$TARGET_ROOT/scripts/harness-sync.sh" \
  "$TARGET_ROOT/scripts/harness-issue-smoke-test.sh" \
  "$TARGET_ROOT/scripts/harness-smoke-test.sh" \
  "$TARGET_ROOT/scripts/harness-release-smoke-test.sh"

cat <<EOF
{
  "adapter_version": "$ADAPTER_VERSION",
  "source_root": "$SOURCE_ROOT",
  "target_repo": "$TARGET_ROOT",
  "profile_id": "$profile_id",
  "repository": "$repository_slug",
  "base_branch": "$base_branch",
  "metadata_path": "$TARGET_ROOT/.harness/cache/adapter-source.json",
  "installed": [
    "scripts/harness_state.py",
    "scripts/harness-home.sh",
    "scripts/harness-doctor.sh",
    "scripts/harness-issue.sh",
    "scripts/harness-claim.sh",
    "scripts/harness-drive.sh",
    "scripts/harness-ship.sh",
    "scripts/harness-rc.sh",
    "scripts/harness-meetings.sh",
    "scripts/harness-maintenance.sh",
    "scripts/harness-profile.sh",
    "scripts/harness-paths.sh",
    "scripts/harness-runs.sh",
    "scripts/harness-signals.sh",
    "scripts/harness-contracts.sh",
    "scripts/harness-git.sh",
    "scripts/harness-bootstrap.sh",
    "scripts/harness-sync.sh",
    "scripts/harness-issue-smoke-test.sh",
    "scripts/harness-smoke-test.sh",
    "scripts/harness-release-smoke-test.sh",
    "install-codex.sh",
    "harness/README.md",
    "harness/profiles/default.yml",
    "harness/templates/contract.yml",
    "harness/templates/issue.md",
    "harness/templates/pr.md",
    "harness/templates/run.yml",
    "harness/templates/signal.yml",
    "harness/templates/meeting.yml",
    "harness/templates/rc.yml",
    "harness/policies/defaults.yml",
    "skills/harness/SKILL.md",
    "skills/harness/README.md",
    "skills/harness/CHANGELOG.md"
  ]
}
EOF
