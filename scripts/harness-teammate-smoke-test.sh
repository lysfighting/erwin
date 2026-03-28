#!/usr/bin/env bash
# harness-teammate-smoke-test.sh — validates leader-led teammate delegation and governed communication
#
# USAGE:
#   ./scripts/harness-teammate-smoke-test.sh

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

TMPDIR_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/harness-teammate-smoke.XXXXXX")
trap 'rm -rf "$TMPDIR_ROOT"' EXIT

PASS=0
FAIL=0

pass() { PASS=$((PASS + 1)); echo "  PASS  $1"; }
fail() { FAIL=$((FAIL + 1)); echo "  FAIL  $1"; }
section() { echo ""; echo "=== $1 ==="; }

TEST_REPO="$TMPDIR_ROOT/repo"
mkdir -p "$TEST_REPO/scripts" "$TEST_REPO/harness/profiles" "$TEST_REPO/harness/templates" "$TEST_REPO/harness/policies"

while IFS= read -r rel_path; do
  [[ -n "$rel_path" ]] || continue
  mkdir -p "$TEST_REPO/$(dirname "$rel_path")"
  cp "$REPO_ROOT/$rel_path" "$TEST_REPO/$rel_path"
done < <(python3 "$REPO_ROOT/scripts/harness_state.py" adapter package-files)

cp "$REPO_ROOT/harness/profiles/erwin.yml" "$TEST_REPO/harness/profiles/"

cat > "$TEST_REPO/.gitignore" <<'EOF'
.harness/
EOF

(
  cd "$TEST_REPO" &&
    git init -q &&
    git config user.email "teammate-smoke@example.com" &&
    git config user.name "Harness Teammate Test" &&
    chmod +x scripts/*.sh install-codex.sh
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

section "Leader claim"
claim_output=$(run_in_repo bash scripts/harness-claim.sh 33001 "Teammate Smoke Mission")
assert_json_field "$claim_output" '"mission_id": "issue-33001"' "claim creates teammate smoke mission"
leader_run_id=$(echo "$claim_output" | python3 -c 'import json,sys; print(json.load(sys.stdin)["run_id"])')

section "Delegation"
assign_a=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-auth implementation-verifier "Patch auth callers" --scope="auth callers" --report="progress|blockers|verification")
assign_b=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-ops ops-simulator "Check release-facing auth flows" --scope="release auth smoke" --report="progress|risk_raised|verification")
assert_json_field "$assign_a" '"owner_id": "worker-auth"' "first teammate assignment created"
assert_json_field "$assign_b" '"owner_id": "worker-ops"' "second teammate assignment created"
run_a=$(echo "$assign_a" | python3 -c 'import json,sys; print(json.load(sys.stdin)["run"]["run_id"])')
run_b=$(echo "$assign_b" | python3 -c 'import json,sys; print(json.load(sys.stdin)["run"]["run_id"])')

list_output=$(run_in_repo bash scripts/harness-teammates.sh list "$leader_run_id")
assert_json_field "$list_output" '"owner_id": "worker-auth"' "teammate list shows first delegated run"
assert_json_field "$list_output" '"owner_id": "worker-ops"' "teammate list shows second delegated run"

section "Reporting and communication"
report_output=$(run_in_repo bash scripts/harness-teammates.sh report "$run_a" "Auth caller patch is ready for sibling review" --kind=progress --evidence=src/auth.ts)
assert_json_field "$report_output" '"kind": "progress"' "teammate report emits progress signal"
assert_json_field "$report_output" '"latest_summary": "Auth caller patch is ready for sibling review"' "teammate report updates run summary"

message_output=$(run_in_repo bash scripts/harness-teammates.sh message "$run_a" "$run_b" interface_change "Caller signature now requires scope" --evidence=src/auth.ts)
assert_json_field "$message_output" '"kind": "interface_change"' "teammate message emits sibling signal"
assert_json_field "$message_output" "\"$run_b\"" "teammate message targets sibling run"

section "Meeting escalation"
escalate_output=$(run_in_repo bash scripts/harness-teammates.sh escalate "$run_b" "Do we need a release-note update for auth scope?" --materialize)
assert_json_field "$escalate_output" '"kind": "requires_meeting"' "teammate escalation emits meeting signal"
assert_json_field "$escalate_output" '"state": "active"' "teammate escalation materializes active meeting"
meeting_id=$(echo "$escalate_output" | python3 -c 'import json,sys; print(json.load(sys.stdin)["meeting"]["meeting_id"])')

participant_a=$(run_in_repo bash scripts/harness-meetings.sh add-participant "$meeting_id" "$run_a" implementation-verifier)
participant_b=$(run_in_repo bash scripts/harness-meetings.sh add-participant "$meeting_id" "$run_b" ops-simulator)
assert_json_field "$participant_a" '"participants"' "meeting stores participant roster"
assert_json_field "$participant_b" '"worker-ops"' "meeting records second participant"

phase_output=$(run_in_repo bash scripts/harness-meetings.sh set-phase "$meeting_id" independent_review)
assert_json_field "$phase_output" '"phase": "independent_review"' "meeting phase advances to independent review"

position_a=$(run_in_repo bash scripts/harness-meetings.sh add-position "$meeting_id" "$run_a" independent_review support "No schema drift in auth callers." --evidence=src/auth.ts)
position_b=$(run_in_repo bash scripts/harness-meetings.sh add-position "$meeting_id" "$run_b" cross_examination challenge "Release notes should mention the new scope requirement." --evidence=CHANGELOG.md)
assert_json_field "$position_a" '"positions"' "meeting stores first position"
assert_json_field "$position_b" '"cross_examination"' "meeting stores cross-examination position"

section "Home and doctor"
home_output=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output" '"active_teammate_count": 2' "home surfaces active teammate runs"
assert_json_field "$home_output" '"coordination_signal_count":' "home surfaces coordination signal count"

doctor_output=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output" '"issue_count": 0' "doctor stays healthy with governed teammate runtime"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
