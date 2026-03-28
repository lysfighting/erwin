#!/usr/bin/env bash
# harness-smoke-test.sh — end-to-end smoke test for the repository-local harness prototype
#
# USAGE:
#   bash scripts/harness-smoke-test.sh
#
# Creates a temporary git repo, copies the harness definition + runtime substrate,
# then validates claim → meeting → outcome adoption → doctor/home health.

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

TMPDIR_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/harness-smoke-test.XXXXXX")
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
    git config user.email "test@example.com" &&
    git config user.name "Harness Test" &&
    git checkout -q -b main &&
    chmod +x scripts/*.sh install-codex.sh &&
    git add . &&
    git commit -q -m "Initial commit"
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
claim_output=$(run_in_repo bash scripts/harness-claim.sh 20001 "Harness Smoke Mission")
assert_json_field "$claim_output" '"mission_id": "issue-20001"' "claim creates mission"
assert_json_field "$claim_output" '"run_id": "run-' "claim returns run id"
assert_json_field "$claim_output" 'issue-20001.yml' "claim materializes mission path"

run_id=$(printf '%s\n' "$claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
run_branch=$(printf '%s\n' "$claim_output" | sed -n 's/^[[:space:]]*"branch": "\(mission\/[^"]*\)",$/\1/p' | head -n1)

section "Mission registry"
missions_output=$(run_in_repo bash scripts/harness-missions.sh list)
assert_json_field "$missions_output" '"mission_id": "issue-20001"' "missions list includes claimed mission"
mission_get_output=$(run_in_repo bash scripts/harness-missions.sh get issue-20001)
assert_json_field "$mission_get_output" '"state": "claimed"' "mission starts in claimed state"
assert_json_field "$mission_get_output" '"active_run_ids":' "mission stores active run linkage"
authority_baseline=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_id")
assert_json_field "$authority_baseline" '"contract_gap_count": 0' "initial run authority audit starts without contract drift"

section "Live state drift"
reprovision_output=$(run_in_repo bash scripts/harness-drive.sh reprovision "$run_id" --mode=worktree --apply)
assert_json_field "$reprovision_output" '"worktree_path":' "reprovision can materialize worktree mode"
worktree_path=$(printf '%s\n' "$reprovision_output" | sed -n 's/.*"worktree_path": "\(\/[^"]*\)".*/\1/p' | head -n1)
drift_update=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$run_id"
{
  "branch": "drift-branch"
}
EOF
)
assert_json_field "$drift_update" '"branch": "drift-branch"' "run record can be forced out of sync with live worktree state"
live_audit_output=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_id")
assert_json_field "$live_audit_output" '"claim_ref": "live_branch"' "run authority audit reports live worktree branch drift"
doctor_output_live=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output_live" '"kind": "run-live-state-drift"' "doctor surfaces live git drift"
repair_update=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$run_id"
{
  "branch": "$run_branch",
  "worktree_path": "$worktree_path"
}
EOF
)
assert_json_field "$repair_update" "\"branch\": \"$run_branch\"" "run record can be repaired back to the live branch"
live_audit_recovered=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_id")
assert_json_field "$live_audit_recovered" '"live_state_gap_count": 0' "run authority audit clears after live branch is repaired"

section "Live intent drift"
desired_update=$(run_in_repo bash scripts/harness-runs.sh set-desired "$run_id" --goal=done --verification=passed)
assert_json_field "$desired_update" '"goal_state": "done"' "desired state can express clean completion intent"
printf 'dirty-worktree\n' > "$worktree_path/dirty.txt"
dirty_intent_audit=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_id")
assert_json_field "$dirty_intent_audit" '"claim_ref": "live_dirty"' "run authority audit reports dirty live worktree against clean intent"
assert_json_field "$dirty_intent_audit" '"expected_value": "clean"' "run authority audit explains the clean-intent expectation"
doctor_output_dirty_intent=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output_dirty_intent" '"kind": "run-live-state-drift"' "doctor surfaces live dirty intent drift"
rm -f "$worktree_path/dirty.txt"
dirty_intent_recovered=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_id")
assert_json_field "$dirty_intent_recovered" '"live_state_gap_count": 0' "run authority audit clears after dirty live worktree is cleaned"

section "Meeting lifecycle"
meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$run_id" "Should this mission stay solo?")
assert_json_field "$meeting_output" '"state": "active"' "meeting starts active"
meeting_id=$(printf '%s\n' "$meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)

finding_output=$(run_in_repo bash scripts/harness-meetings.sh add-finding "$meeting_id" accepted "Mission should stay solo.")
assert_json_field "$finding_output" '"accepted"' "meeting stores accepted finding"

action_output=$(run_in_repo bash scripts/harness-meetings.sh add-action "$meeting_id" "Continue solo execution.")
assert_json_field "$action_output" '"action_items"' "meeting stores action item"

close_output=$(run_in_repo bash scripts/harness-drive.sh meeting-close "$run_id" "$meeting_id" "Council concluded solo path.")
assert_json_field "$close_output" '"state": "closed"' "meeting closes cleanly"
assert_json_field "$close_output" '"latest_summary": "Mission should stay solo."' "meeting outcome updates run summary"
assert_json_field "$close_output" '"next_expected_action": "Continue solo execution."' "meeting outcome updates next action"

section "Home and doctor"
home_output=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output" '"active_mission_count": 1' "home shows one active mission"
assert_json_field "$home_output" '"active_run_count": 1' "home shows one active run"
assert_json_field "$home_output" '"open_meeting_action_count": 1' "home surfaces open meeting action"

adapter_output=$(run_in_repo bash scripts/harness-adapter.sh validate)
assert_json_field "$adapter_output" '"valid": true' "adapter manifest validates in core smoke repo"

doctor_output_before=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output_before" '"issue_count": 0' "doctor stays healthy after outcome adoption"

section "Action resolution"
action_id=$(printf '%s\n' "$action_output" | sed -n 's/^[[:space:]]*"action_id": "\(action-[^"]*\)",$/\1/p' | head -n1)
resolve_output=$(run_in_repo bash scripts/harness-meetings.sh resolve-action "$meeting_id" "$action_id" done)
assert_json_field "$resolve_output" '"state": "done"' "meeting action can be resolved"

home_output_after=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output_after" '"open_meeting_action_count": 0' "home clears resolved meeting action"

section "Authority adjudication"
block_output=$(run_in_repo bash scripts/harness-drive.sh block "$run_id" "Temporary blocker for authority reconcile")
assert_json_field "$block_output" '"kind": "blocked"' "blocking a run emits blocked signal"
blocked_signal_id=$(printf '%s\n' "$block_output" | sed -n 's/^[[:space:]]*"signal_id": "\(sig-[^"]*\)",$/\1/p' | head -n1)
resume_output=$(run_in_repo bash scripts/harness-drive.sh resume "$run_id" "Blocker cleared")
assert_json_field "$resume_output" '"state": "ready"' "run can resume after blocker clears"
doctor_output_blocked=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output_blocked" '"kind": "signal-authority-superseded"' "doctor surfaces superseded blocked signal"
reconcile_output=$(run_in_repo bash scripts/harness-maintenance.sh reconcile)
assert_json_field "$reconcile_output" "$blocked_signal_id" "reconcile adjudicates superseded blocked signal"

authority_signal_id="sig-smoke-l1-blocked"
authority_signal_output=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "$authority_signal_id",
  "kind": "blocked",
  "source_type": "leader",
  "source_id": "chair",
  "run_id": "$run_id",
  "mission_id": "issue-20001",
  "applies_to": ["$run_id", "issue-20001"],
  "scope_type": "local",
  "title": "Higher-authority blocker",
  "summary": "Live release gate says this run is blocked.",
  "authority_level": "L1",
  "confidence": "high",
  "emitted_at": "2099-01-01T00:00:00Z",
  "fresh_until": "2099-01-02T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$authority_signal_output" "\"signal_id\": \"$authority_signal_id\"" "manual L1 blocked signal can be emitted"
authority_audit_output=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_id")
assert_json_field "$authority_audit_output" '"higher_authority_conflict_count": 1' "run authority audit reports higher-authority cross-surface conflict"
assert_json_field "$authority_audit_output" '"resolution_kind": "run_state_overridden"' "run authority audit explains why the signal currently wins"
doctor_output_authority=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output_authority" '"kind": "run-authority-conflict"' "doctor surfaces unresolved run-level authority conflict"
resolved_authority_signal=$(run_in_repo bash scripts/harness-signals.sh set-state "$authority_signal_id" resolved "Leader resolved the higher-authority blocker.")
assert_json_field "$resolved_authority_signal" '"state": "resolved"' "higher-authority authority-conflict signal can be resolved"
readopt_output=$(run_in_repo bash scripts/harness-drive.sh adopt-meeting "$run_id" "$meeting_id")
assert_json_field "$readopt_output" "\"meeting_id\": \"$meeting_id\"" "run can re-adopt the closed meeting outcome after authority handling"
doctor_output_recovered=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output_recovered" '"issue_count": 0' "doctor returns to healthy state after run-level authority conflict is resolved"

section "Normalize"
cat > "$TEST_REPO/.harness/runs/registry.yml" <<'EOF'
{
  "version": 1,
  "updated_at": "2026-03-28T00:00:00Z",
  "runs": [
    {
      "run_id": "run-smoke-missing-fields",
      "mission_id": "issue-999",
      "state": "ready"
    }
  ]
}
EOF

cat > "$TEST_REPO/.harness/signals/log.yml" <<'EOF'
{
  "version": 1,
  "updated_at": "2026-03-28T00:00:00Z",
  "signals": [
    {
      "signal_id": "sig-smoke-missing-emitted",
      "state": "active",
      "fresh_until": "2026-03-29T00:00:00Z"
    }
  ]
}
EOF

normalize_output=$(run_in_repo bash scripts/harness-maintenance.sh normalize)
assert_json_field "$normalize_output" '"run_count": 1' "normalize repairs missing run metadata"
assert_json_field "$normalize_output" '"signal_count": 1' "normalize repairs missing signal metadata"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
