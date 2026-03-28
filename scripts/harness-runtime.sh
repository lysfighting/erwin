#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
LOCAL_BIN_DIR="$REPO_ROOT/.harness/bin"
BIN_PATH="$LOCAL_BIN_DIR/harness-state"
PACKAGED_BIN_ROOT="$REPO_ROOT/harness/bin"

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

platform_key() {
  printf '%s-%s\n' "$(normalize_os)" "$(normalize_arch)"
}

packaged_bin_path() {
  printf '%s/%s/harness-state\n' "$PACKAGED_BIN_ROOT" "$(platform_key)"
}

source_tree_newer_than() {
  local target="$1"
  [[ -e "$target" ]] || return 0
  [[ "$REPO_ROOT/go.mod" -nt "$target" || "$REPO_ROOT/go.sum" -nt "$target" ]] && return 0
  find "$REPO_ROOT/cmd" "$REPO_ROOT/internal" -type f -name '*.go' -newer "$target" -print -quit | grep -q .
}

build_local_binary() {
  mkdir -p "$LOCAL_BIN_DIR"
  (
    cd "$REPO_ROOT"
    go build -o "$BIN_PATH" ./cmd/harness-state
  )
}

needs_build=0
PACKAGED_BIN_PATH="$(packaged_bin_path)"

if [[ -x "$PACKAGED_BIN_PATH" ]]; then
  if command -v go >/dev/null 2>&1 && source_tree_newer_than "$PACKAGED_BIN_PATH"; then
    needs_build=1
  else
    export HARNESS_REPO_ROOT="$REPO_ROOT"
    exec "$PACKAGED_BIN_PATH" "$@"
  fi
fi

if ! command -v go >/dev/null 2>&1; then
  echo "ERROR: no packaged harness runtime for $(platform_key), and Go is required to build it locally." >&2
  exit 1
fi

if [[ ! -x "$BIN_PATH" ]]; then
  needs_build=1
elif source_tree_newer_than "$BIN_PATH"; then
  needs_build=1
fi

if [[ "$needs_build" -eq 1 ]]; then
  build_local_binary
fi

export HARNESS_REPO_ROOT="$REPO_ROOT"
exec "$BIN_PATH" "$@"
