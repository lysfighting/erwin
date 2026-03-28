#!/usr/bin/env bash
# harness-release-smoke-test.sh — validates the ship/rc lifecycle
#
# USAGE:
#   ./scripts/harness-release-smoke-test.sh

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

TMPDIR_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/harness-release-smoke.XXXXXX")
trap 'rm -rf "$TMPDIR_ROOT"' EXIT

PASS=0
FAIL=0

pass() { PASS=$((PASS + 1)); echo "  PASS  $1"; }
fail() { FAIL=$((FAIL + 1)); echo "  FAIL  $1"; }
section() { echo ""; echo "=== $1 ==="; }

TEST_REPO="$TMPDIR_ROOT/repo"
ORIGIN_REPO="$TMPDIR_ROOT/origin.git"
mkdir -p "$TEST_REPO/scripts" "$TEST_REPO/harness/profiles" "$TEST_REPO/harness/templates" "$TEST_REPO/harness/policies"

cp "$REPO_ROOT/scripts/harness_state.py" "$TEST_REPO/scripts/"
cp "$REPO_ROOT/scripts/harness-home.sh" "$TEST_REPO/scripts/"
cp "$REPO_ROOT/scripts/harness-doctor.sh" "$TEST_REPO/scripts/"
cp "$REPO_ROOT/scripts/harness-claim.sh" "$TEST_REPO/scripts/"
cp "$REPO_ROOT/scripts/harness-drive.sh" "$TEST_REPO/scripts/"
cp "$REPO_ROOT/scripts/harness-ship.sh" "$TEST_REPO/scripts/"
cp "$REPO_ROOT/scripts/harness-rc.sh" "$TEST_REPO/scripts/"
cp "$REPO_ROOT/harness/profiles/erwin.yml" "$TEST_REPO/harness/profiles/"
cp "$REPO_ROOT/harness/templates/"*.yml "$TEST_REPO/harness/templates/"
cp "$REPO_ROOT/harness/templates/pr.md" "$TEST_REPO/harness/templates/"
cp "$REPO_ROOT/harness/policies/defaults.yml" "$TEST_REPO/harness/policies/"

cat > "$TEST_REPO/.gitignore" <<'EOF'
.harness/
EOF

(
  git init -q --bare "$ORIGIN_REPO" &&
  cd "$TEST_REPO" &&
    git init -q &&
    git checkout -qb main &&
    git config user.email "release-smoke@example.com" &&
    git config user.name "Harness Release Test" &&
    git remote add origin "$ORIGIN_REPO" &&
    chmod +x scripts/harness-home.sh scripts/harness-doctor.sh scripts/harness-claim.sh scripts/harness-drive.sh scripts/harness-ship.sh scripts/harness-rc.sh &&
    printf "hello\n" > README.md &&
    git add . &&
    git commit -qm "init harness release smoke" &&
    git push -u origin main >/dev/null
) || exit 1

run_in_repo() {
  (
    cd "$TEST_REPO" &&
      "$@"
  )
}

assert_json_field() {
  local json="$1"
  local needle="$2"
  local label="$3"
  if echo "$json" | grep -q "$needle"; then
    pass "$label"
  else
    fail "$label"
    echo "$json"
  fi
}

section "Claim"
claim_output=$(run_in_repo ./scripts/harness-claim.sh 32001 "Release Smoke Mission" --provision=branch --apply)
assert_json_field "$claim_output" '"mission_id": "issue-32001"' "claim creates release mission"
run_id=$(echo "$claim_output" | python3 -c 'import json,sys; print(json.load(sys.stdin)["run_id"])')

section "Ship"
verify_output=$(run_in_repo ./scripts/harness-ship.sh verify "$run_id" "targeted-test:pass")
assert_json_field "$verify_output" '"verification_status": "passed"' "ship verify records verification"

ready_output=$(run_in_repo ./scripts/harness-ship.sh ready "$run_id" "Ready for PR")
assert_json_field "$ready_output" '"ship_status": "ready"' "ship ready advances run"

push_output=$(run_in_repo ./scripts/harness-ship.sh push "$run_id")
assert_json_field "$push_output" '"ship_status": "pushed"' "ship push records pushed state"

if git --git-dir="$ORIGIN_REPO" show-ref --verify --quiet refs/heads/mission/32001-release-smoke-mission; then
  pass "ship push updates remote branch"
else
  fail "ship push updates remote branch"
fi

pr_draft_output=$(run_in_repo ./scripts/harness-ship.sh pr-draft "$run_id")
assert_json_field "$pr_draft_output" '"mode": "draft"' "ship pr-draft returns draft payload"
assert_json_field "$pr_draft_output" '"base_branch": "main"' "ship pr-draft uses base branch"

pr_output=$(run_in_repo ./scripts/harness-ship.sh pr "$run_id" "PR-32001" "https://example.invalid/pr/32001")
assert_json_field "$pr_output" '"ship_status": "in_review"' "ship pr records PR state"
assert_json_field "$pr_output" '"pr_ref": "PR-32001"' "ship pr stores PR ref"

merged_output=$(run_in_repo ./scripts/harness-ship.sh merged "$run_id" "Merged cleanly")
assert_json_field "$merged_output" '"ship_status": "merged"' "ship merged records merge"
assert_json_field "$merged_output" '"state": "completed"' "ship merged completes run"

section "RC"
rc_create_output=$(run_in_repo ./scripts/harness-rc.sh create v1.2.3 "$run_id")
assert_json_field "$rc_create_output" '"version": "v1.2.3"' "rc create records version"
assert_json_field "$rc_create_output" '"state": "active"' "rc create opens active rc"
rc_id=$(echo "$rc_create_output" | python3 -c 'import json,sys; print(json.load(sys.stdin)["rc_id"])')

rc_verify_output=$(run_in_repo ./scripts/harness-rc.sh verify "$rc_id" "staging-smoke:pass")
assert_json_field "$rc_verify_output" '"verification_status": "passed"' "rc verify records verification"

rc_promote_output=$(run_in_repo ./scripts/harness-rc.sh promote "$rc_id" "Promoted after staging smoke")
assert_json_field "$rc_promote_output" '"state": "promoted"' "rc promote archives promoted rc"

section "Home and doctor"
home_output=$(run_in_repo ./scripts/harness-home.sh)
doctor_output=$(run_in_repo ./scripts/harness-doctor.sh)
assert_json_field "$home_output" '"active_rc_count": 0' "home clears active rc after promotion"
assert_json_field "$doctor_output" '"issue_count": 0' "doctor stays healthy after release flow"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
