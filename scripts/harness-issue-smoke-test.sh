#!/usr/bin/env bash
# harness-issue-smoke-test.sh — validates the issue surface in offline draft mode
#
# USAGE:
#   ./scripts/harness-issue-smoke-test.sh

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

TMPDIR_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/harness-issue-smoke.XXXXXX")
trap 'rm -rf "$TMPDIR_ROOT"' EXIT

PASS=0
FAIL=0

pass() { PASS=$((PASS + 1)); echo "  PASS  $1"; }
fail() { FAIL=$((FAIL + 1)); echo "  FAIL  $1"; }
section() { echo ""; echo "=== $1 ==="; }

TEST_REPO="$TMPDIR_ROOT/repo"
mkdir -p "$TEST_REPO/scripts" "$TEST_REPO/harness/templates"

cp "$REPO_ROOT/scripts/harness_state.py" "$TEST_REPO/scripts/"
cp "$REPO_ROOT/scripts/harness-issue.sh" "$TEST_REPO/scripts/"
cp "$REPO_ROOT/harness/templates/issue.md" "$TEST_REPO/harness/templates/"

(
  cd "$TEST_REPO" &&
    git init -q &&
    git checkout -qb main &&
    chmod +x scripts/harness-issue.sh
) || exit 1

run_in_repo() {
  (
    cd "$TEST_REPO" &&
      "$@"
  )
}

assert_contains() {
  local text="$1"
  local needle="$2"
  local label="$3"
  if echo "$text" | grep -q -- "$needle"; then
    pass "$label"
  else
    fail "$label"
    echo "$text"
  fi
}

section "Draft"
draft_output=$(run_in_repo ./scripts/harness-issue.sh draft "Issue Surface Smoke" --objective="Validate issue drafting" --criteria="Body renders|Criteria stay structured" --note="Offline smoke")
assert_contains "$draft_output" '"mode": "draft"' "issue draft reports draft mode"
assert_contains "$draft_output" '"title": "Issue Surface Smoke"' "issue draft stores title"
assert_contains "$draft_output" '## Objective' "issue draft renders markdown body"
assert_contains "$draft_output" 'Criteria stay structured' "issue draft includes criteria"

section "Body"
body_output=$(run_in_repo ./scripts/harness-issue.sh body "Body Surface Smoke" --criteria="Alpha|Beta")
assert_contains "$body_output" '# Body Surface Smoke' "issue body renders title header"
assert_contains "$body_output" '- Alpha' "issue body renders first criterion"
assert_contains "$body_output" '- Beta' "issue body renders second criterion"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
