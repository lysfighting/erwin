#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_DIR="$REPO_ROOT/.harness/bin"
BIN_PATH="$BIN_DIR/harness-state"

needs_build=0

if ! command -v go >/dev/null 2>&1; then
  echo "ERROR: Go is required to build the harness runtime locally." >&2
  exit 1
fi

if [[ ! -x "$BIN_PATH" ]]; then
  needs_build=1
elif [[ "$REPO_ROOT/go.mod" -nt "$BIN_PATH" || "$REPO_ROOT/go.sum" -nt "$BIN_PATH" ]]; then
  needs_build=1
elif find "$REPO_ROOT/cmd" "$REPO_ROOT/internal" -type f -name '*.go' -newer "$BIN_PATH" -print -quit | grep -q .; then
  needs_build=1
fi

if [[ "$needs_build" -eq 1 ]]; then
  mkdir -p "$BIN_DIR"
  (
    cd "$REPO_ROOT"
    go build -o "$BIN_PATH" ./cmd/harness-state
  )
fi

export HARNESS_REPO_ROOT="$REPO_ROOT"
exec "$BIN_PATH" "$@"
