#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BUILD_SCRIPT="$SCRIPT_DIR/harness-build-runtime.sh"
OUTPUT_ROOT="$REPO_ROOT/.harness/dist/runtime"
BUILD_ROOT="$REPO_ROOT/.harness/dist/bin"
TARGETS=()

targets_count() {
  set +u
  local count="${#TARGETS[@]}"
  set -u
  printf '%s\n' "$count"
}

usage() {
  cat <<'EOF'
Usage:
  ./scripts/harness-package-runtime.sh
  ./scripts/harness-package-runtime.sh --target darwin-arm64
  ./scripts/harness-package-runtime.sh --all
  ./scripts/harness-package-runtime.sh --output /tmp/harness-dist

Build packaged harness-state binaries, then assemble per-target tar.gz runtime bundles.
Default output root is .harness/dist/runtime.
EOF
}

normalize_os() {
  local raw
  raw="$(uname -s)"
  case "$raw" in
    Darwin) printf 'darwin\n' ;;
    Linux) printf 'linux\n' ;;
    *) printf '%s\n' "$(printf '%s' "$raw" | tr '[:upper:]' '[:lower:]')"
  esac
}

normalize_arch() {
  local raw
  raw="$(uname -m)"
  case "$raw" in
    x86_64|amd64) printf 'amd64\n' ;;
    arm64|aarch64) printf 'arm64\n' ;;
    *) printf '%s\n' "$raw" ;;
  esac
}

default_target() {
  printf '%s-%s\n' "$(normalize_os)" "$(normalize_arch)"
}

add_target() {
  local target="$1"
  [[ -n "$target" ]] || return 0
  if [[ "$(targets_count)" -gt 0 ]]; then
    local existing
    for existing in "${TARGETS[@]}"; do
      [[ "$existing" == "$target" ]] && return 0
    done
  fi
  TARGETS+=("$target")
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --target|-t)
      [[ -n "${2:-}" ]] || { echo "ERROR: --target requires os-arch" >&2; exit 1; }
      add_target "$2"
      shift 2
      ;;
    --all)
      add_target "darwin-arm64"
      add_target "darwin-amd64"
      add_target "linux-arm64"
      add_target "linux-amd64"
      shift
      ;;
    --output|-o)
      [[ -n "${2:-}" ]] || { echo "ERROR: --output requires a directory" >&2; exit 1; }
      OUTPUT_ROOT="$2"
      shift 2
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      echo "ERROR: unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

if [[ "$(targets_count)" -eq 0 ]]; then
  add_target "$(default_target)"
fi

mkdir -p "$OUTPUT_ROOT" "$BUILD_ROOT"

build_args=(--output "$BUILD_ROOT")
for target in "${TARGETS[@]}"; do
  build_args+=(--target "$target")
done

"$BUILD_SCRIPT" "${build_args[@]}" >/dev/null

json_escape() {
  printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}

json_array() {
  local first=1
  printf '['
  for item in "$@"; do
    [[ "$first" -eq 0 ]] && printf ', '
    first=0
    printf '"%s"' "$(json_escape "$item")"
  done
  printf ']'
}

ARCHIVES=()

for target in "${TARGETS[@]}"; do
  package_root="$OUTPUT_ROOT/$target"
  archive_name="harness-runtime-${target}.tar.gz"
  archive_path="$OUTPUT_ROOT/$archive_name"
  rm -rf "$package_root"
  mkdir -p "$package_root/docs" "$package_root/harness"
  cp "$BUILD_ROOT/$target/harness-state" "$package_root/harness-state"
  cp "$REPO_ROOT/harness/adapter-manifest.yml" "$package_root/harness/adapter-manifest.yml"
  cp "$REPO_ROOT/docs/harness-public-contract.md" "$package_root/docs/harness-public-contract.md"
  cp "$REPO_ROOT/harness/bin/README.md" "$package_root/README.md"
  (
    cd "$OUTPUT_ROOT"
    tar -czf "$archive_path" "$target"
  )
  ARCHIVES+=("$archive_path")
done

cat <<EOF
{
  "output_root": "$OUTPUT_ROOT",
  "targets": $(json_array "${TARGETS[@]}"),
  "archives": $(json_array "${ARCHIVES[@]}")
}
EOF
