#!/usr/bin/env bash
# harness-bootstrap-smoke-test.sh — verifies the reusable bootstrap path
#
# USAGE:
#   ./scripts/harness-bootstrap-smoke-test.sh
#
# Creates a temporary git repo, bootstraps the harness adapter into it, then
# validates the target repo can run the leader surfaces and install its local
# Codex-native skill.

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

TMPDIR_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/harness-bootstrap-smoke.XXXXXX")
trap 'rm -rf "$TMPDIR_ROOT"' EXIT

PASS=0
FAIL=0

pass() { PASS=$((PASS + 1)); echo "  PASS  $1"; }
fail() { FAIL=$((FAIL + 1)); echo "  FAIL  $1"; }
section() { echo ""; echo "=== $1 ==="; }

TARGET_REPO="$TMPDIR_ROOT/target-repo"
TARGET_SKILLS="$TMPDIR_ROOT/codex-skills"
mkdir -p "$TARGET_REPO"

(
  cd "$TARGET_REPO" &&
    git init -q &&
    git checkout -qb main &&
    git config user.email "bootstrap-test@example.com" &&
    git config user.name "Harness Bootstrap Test" &&
    printf "# bootstrap smoke repo\n" > README.md &&
    git add README.md &&
    git commit -qm "init"
) || exit 1

run_in_target() {
  (
    cd "$TARGET_REPO" &&
      "$@"
  )
}

assert_contains() {
  local text="$1"
  local needle="$2"
  local label="$3"
  if echo "$text" | grep -q "$needle"; then
    pass "$label"
  else
    fail "$label"
    echo "$text"
  fi
}

section "Bootstrap"
bootstrap_output="$("$REPO_ROOT/scripts/harness-bootstrap.sh" "$TARGET_REPO")"
assert_contains "$bootstrap_output" '"target_repo"' "bootstrap returns summary"
assert_contains "$bootstrap_output" '"skills/harness/SKILL.md"' "bootstrap installs local skill files"
assert_contains "$bootstrap_output" '"install-codex.sh"' "bootstrap installs local Codex installer"
assert_contains "$bootstrap_output" '"metadata_path"' "bootstrap writes adapter metadata"
assert_contains "$bootstrap_output" '"scripts/harness-ship.sh"' "bootstrap installs ship surface"
assert_contains "$bootstrap_output" '"scripts/harness-rc.sh"' "bootstrap installs rc surface"

if run_in_target test -f scripts/harness-home.sh && \
   run_in_target test -f scripts/harness-bootstrap.sh && \
   run_in_target test -f scripts/harness-sync.sh && \
   run_in_target test -f scripts/harness-ship.sh && \
   run_in_target test -f scripts/harness-rc.sh && \
   run_in_target test -f skills/harness/SKILL.md && \
   run_in_target test -f install-codex.sh; then
  pass "bootstrap writes expected files"
else
  fail "bootstrap writes expected files"
fi

section "Leader surfaces"
home_output="$(run_in_target ./scripts/harness-home.sh)"
doctor_output="$(run_in_target ./scripts/harness-doctor.sh)"
claim_output="$(run_in_target ./scripts/harness-claim.sh 31001 "Bootstrapped Target Mission")"
assert_contains "$home_output" '"next_action": "No active runs. Claim a mission."' "bootstrapped home starts empty"
assert_contains "$doctor_output" '"issue_count": 0' "bootstrapped doctor starts healthy"
assert_contains "$claim_output" '"mission_id": "issue-31001"' "bootstrapped claim creates mission"

section "Adapter sync"
sync_status_output="$(run_in_target ./scripts/harness-sync.sh status)"
assert_contains "$sync_status_output" '"source_root"' "bootstrapped repo remembers source metadata"
sync_output="$(run_in_target ./scripts/harness-sync.sh)"
assert_contains "$sync_output" '"target_repo"' "bootstrapped repo can resync from remembered source"

section "Local Codex install"
install_output="$(run_in_target ./install-codex.sh --target "$TARGET_SKILLS" harness)"
assert_contains "$install_output" '/harness' "bootstrapped repo installs local harness skill"
if test -f "$TARGET_SKILLS/harness/SKILL.md"; then
  pass "bootstrapped install writes skill files"
else
  fail "bootstrapped install writes skill files"
fi

section "Safety boundary"
set +e
"$REPO_ROOT/scripts/harness-bootstrap.sh" "$TARGET_REPO" >/tmp/harness-bootstrap-repeat.out 2>/tmp/harness-bootstrap-repeat.err
rc=$?
set -e
if [ "$rc" -ne 0 ] && grep -q 'refusing to overwrite existing file without --force' /tmp/harness-bootstrap-repeat.err; then
  pass "bootstrap refuses overwrite without --force"
else
  fail "bootstrap refuses overwrite without --force"
  cat /tmp/harness-bootstrap-repeat.err 2>/dev/null || true
fi

force_output="$("$REPO_ROOT/scripts/harness-bootstrap.sh" --force "$TARGET_REPO")"
assert_contains "$force_output" '"target_repo"' "bootstrap force path succeeds"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
