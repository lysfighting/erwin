#!/usr/bin/env bash
# harness-all-smoke-test.sh — runs the full local harness smoke suite
#
# USAGE:
#   ./scripts/harness-all-smoke-test.sh

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PASS=0
FAIL=0

run_case() {
  local label="$1"
  shift
  echo ""
  echo "=== $label ==="
  if "$@"; then
    PASS=$((PASS + 1))
    echo "  PASS  $label"
  else
    FAIL=$((FAIL + 1))
    echo "  FAIL  $label"
  fi
}

run_case "adapter validate" "$SCRIPT_DIR/harness-adapter.sh" validate
run_case "issue smoke" "$SCRIPT_DIR/harness-issue-smoke-test.sh"
run_case "core smoke" "$SCRIPT_DIR/harness-smoke-test.sh"
run_case "release smoke" "$SCRIPT_DIR/harness-release-smoke-test.sh"
run_case "teammate smoke" "$SCRIPT_DIR/harness-teammate-smoke-test.sh"
run_case "bootstrap smoke" "$SCRIPT_DIR/harness-bootstrap-smoke-test.sh"

echo ""
echo "PASS=$PASS FAIL=$FAIL"
if [[ "$FAIL" -ne 0 ]]; then
  exit 1
fi
