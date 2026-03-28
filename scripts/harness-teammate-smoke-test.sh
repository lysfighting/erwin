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
done < <("$REPO_ROOT/scripts/harness-adapter.sh" package-files)

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
leader_run_id=$(printf '%s\n' "$claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
mission_output=$(run_in_repo bash scripts/harness-missions.sh get issue-33001)
assert_json_field "$mission_output" '"state": "claimed"' "mission registry tracks leader claim"

section "Delegation"
assign_a=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-auth implementation-verifier "Patch auth callers" --scope="auth callers" --owns="scope:auth-callers" --report="progress|blockers|verification")
assert_json_field "$assign_a" '"owner_id": "worker-auth"' "first teammate assignment created"
run_a=$(printf '%s\n' "$assign_a" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
assign_b=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-ops ops-simulator "Check release-facing auth flows" --scope="release auth smoke" --owns="scope:release-auth-smoke|file:CHANGELOG.md" --depends-on="$run_a" --report="progress|risk_raised|verification")
assert_json_field "$assign_b" '"owner_id": "worker-ops"' "second teammate assignment created"
run_b=$(printf '%s\n' "$assign_b" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)

list_output=$(run_in_repo bash scripts/harness-teammates.sh list "$leader_run_id")
assert_json_field "$list_output" '"owner_id": "worker-auth"' "teammate list shows first delegated run"
assert_json_field "$list_output" '"owner_id": "worker-ops"' "teammate list shows second delegated run"

section "Ownership transfer"
transfer_output=$(run_in_repo bash scripts/harness-runs.sh transfer-ownership "$run_b" "$run_a" --owns="scope:release-auth-smoke" --reason="Shift release-auth scope into the implementation stream")
assert_json_field "$transfer_output" '"transfer_id": "transfer-' "ownership transfer records transfer event"
assert_json_field "$transfer_output" '"selector_refs":' "ownership transfer records moved selectors"
run_a_after_transfer=$(run_in_repo bash scripts/harness-teammates.sh get "$run_a")
run_b_after_transfer=$(run_in_repo bash scripts/harness-teammates.sh get "$run_b")
assert_json_field "$run_a_after_transfer" '"selector": "scope:release-auth-smoke"' "target teammate receives transferred selector"
if echo "$run_b_after_transfer" | grep -q '"selector": "scope:release-auth-smoke"'; then
  fail "source teammate drops transferred selector"
  echo "$run_b_after_transfer"
else
  pass "source teammate drops transferred selector"
fi

section "Ownership guard"
set +e
overlap_output=$(cd "$TEST_REPO" && bash scripts/harness-teammates.sh assign "$leader_run_id" worker-conflict architecture-challenger "Touch the same auth callers" --scope="auth callers" --owns="scope:auth-callers" 2>&1)
rc=$?
set -e
if [ "$rc" -ne 0 ] && echo "$overlap_output" | grep -q 'ownership overlap detected'; then
  pass "exclusive ownership overlap is rejected"
else
  fail "exclusive ownership overlap is rejected"
  echo "$overlap_output"
fi

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
meeting_id=$(printf '%s\n' "$escalate_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)

role_a=$(run_in_repo bash scripts/harness-meetings.sh add-role "$meeting_id" implementation-verifier --mission="Verify the scoped auth change against code and release notes." --allow="L1: live git state|L2: active runtime state|L3: committed harness definitions" --forbid="L4+: stale summaries")
role_b=$(run_in_repo bash scripts/harness-meetings.sh add-role "$meeting_id" ops-simulator --mission="Simulate release-facing operator impact." --allow="L1: live git state|L2: active runtime state|L3: committed harness definitions" --forbid="Unverified memory")
assert_json_field "$role_a" '"role_specs"' "meeting stores bounded role specs"
assert_json_field "$role_b" '"ops-simulator"' "meeting keeps the second bounded role"

known_context=$(run_in_repo bash scripts/harness-meetings.sh add-known-context "$meeting_id" "Release notes currently mention auth callers but not required scope.")
assert_json_field "$known_context" '"known_context"' "meeting records known context"

locked_decision=$(run_in_repo bash scripts/harness-meetings.sh lock-decision "$meeting_id" "Do not widen teammate scope beyond auth-related release notes." --reason="Bounded council scope for this meeting." --authority="L2" --evidence=harness/profiles/erwin.yml)
assert_json_field "$locked_decision" '"locked_decisions"' "meeting records locked decisions"

participant_a=$(run_in_repo bash scripts/harness-meetings.sh add-participant "$meeting_id" "$run_a" implementation-verifier)
participant_b=$(run_in_repo bash scripts/harness-meetings.sh add-participant "$meeting_id" "$run_b" ops-simulator)
assert_json_field "$participant_a" '"participants"' "meeting stores participant roster"
assert_json_field "$participant_b" '"worker-ops"' "meeting records second participant"

orchestrate_a=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$meeting_id")
assert_json_field "$orchestrate_a" '"new_phase": "independent_review"' "meeting orchestration advances setup to independent review"

position_a=$(run_in_repo bash scripts/harness-meetings.sh add-position "$meeting_id" "$run_a" independent_review support "No schema drift in auth callers." --evidence=src/auth.ts)
position_b=$(run_in_repo bash scripts/harness-meetings.sh add-position "$meeting_id" "$run_b" independent_review challenge "Release notes should mention the new scope requirement." --evidence=CHANGELOG.md)
assert_json_field "$position_a" '"positions"' "meeting stores first position"
assert_json_field "$position_b" '"independent_review"' "meeting stores second independent-review position"

orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$meeting_id")
assert_json_field "$orchestrate_b" '"new_phase": "cross_examination"' "meeting orchestration advances to cross-examination"

cross_exam=$(run_in_repo bash scripts/harness-meetings.sh add-position "$meeting_id" "$run_b" cross_examination challenge "Operator-facing release notes must mention the new scope requirement." --evidence=CHANGELOG.md)
assert_json_field "$cross_exam" '"cross_examination"' "meeting stores cross-examination position"

rejected_finding=$(run_in_repo bash scripts/harness-meetings.sh add-finding "$meeting_id" rejected "Ship immediately without release note updates." --reason="Cross-examination found operator-facing behavior drift." --authority="L2" --evidence=CHANGELOG.md --run="$run_b" --role=ops-simulator)
assert_json_field "$rejected_finding" '"reason": "Cross-examination found operator-facing behavior drift."' "meeting records rejected finding reason"
accepted_finding=$(run_in_repo bash scripts/harness-meetings.sh add-finding "$meeting_id" accepted "Release notes need a scoped auth update before ship." --reason="Independent review and cross-examination agree on operator impact." --authority="L2" --evidence=CHANGELOG.md --run="$run_a" --role=implementation-verifier)
assert_json_field "$accepted_finding" '"accepted"' "meeting records accepted findings"
action_output=$(run_in_repo bash scripts/harness-meetings.sh add-action "$meeting_id" "Update release notes before merging the auth change.")
assert_json_field "$action_output" '"action_items"' "meeting records action items"

orchestrate_c=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$meeting_id")
assert_json_field "$orchestrate_c" '"new_phase": "synthesis"' "meeting orchestration advances to synthesis"
orchestrate_d=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$meeting_id")
assert_json_field "$orchestrate_d" '"new_phase": "output"' "meeting orchestration advances to output once findings exist"

section "Home and doctor"
home_output=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output" '"active_mission_count": 1' "home surfaces active mission count"
assert_json_field "$home_output" '"active_teammate_count": 2' "home surfaces active teammate runs"
assert_json_field "$home_output" '"coordination_signal_count":' "home surfaces coordination signal count"
assert_json_field "$home_output" '"dependency_blocked_count": 1' "home surfaces dependency-blocked runs"
assert_json_field "$home_output" '"ownership_overlap_count": 0' "home stays free of ownership overlap"
mission_graph=$(run_in_repo bash scripts/harness-missions.sh graph)
assert_json_field "$mission_graph" '"mission_count": 1' "mission graph reports the active mission"

doctor_output=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output" '"issue_count": 0' "doctor stays healthy with governed teammate runtime"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
