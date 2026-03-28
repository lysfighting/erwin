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
done < <(python3 "$REPO_ROOT/scripts/harness_state.py" adapter package-files)

cp "$REPO_ROOT/harness/profiles/erwin.yml" "$TEST_REPO/harness/profiles/"

cat > "$TEST_REPO/.gitignore" <<'EOF'
.harness/
EOF

(
  cd "$TEST_REPO" &&
    git init -q &&
    git config user.email "test@example.com" &&
    git config user.name "Harness Test" &&
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

section "Claim"
claim_output=$(run_in_repo bash scripts/harness-claim.sh 20001 "Harness Smoke Mission")
assert_json_field "$claim_output" '"mission_id": "issue-20001"' "claim creates mission"
assert_json_field "$claim_output" '"run_id": "run-' "claim returns run id"

run_id=$(echo "$claim_output" | python3 -c 'import json,sys; print(json.load(sys.stdin)["run_id"])')

section "Meeting lifecycle"
meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$run_id" "Should this mission stay solo?")
assert_json_field "$meeting_output" '"state": "active"' "meeting starts active"
meeting_id=$(echo "$meeting_output" | python3 -c 'import json,sys; print(json.load(sys.stdin)["meeting"]["meeting_id"])')

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
assert_json_field "$home_output" '"active_run_count": 1' "home shows one active run"
assert_json_field "$home_output" '"open_meeting_action_count": 1' "home surfaces open meeting action"

adapter_output=$(run_in_repo bash scripts/harness-adapter.sh validate)
assert_json_field "$adapter_output" '"valid": true' "adapter manifest validates in core smoke repo"

doctor_output_before=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output_before" '"issue_count": 0' "doctor stays healthy after outcome adoption"

section "Action resolution"
action_id=$(echo "$action_output" | python3 -c 'import json,sys; print(json.load(sys.stdin)["action_items"][0]["action_id"])')
resolve_output=$(run_in_repo bash scripts/harness-meetings.sh resolve-action "$meeting_id" "$action_id" done)
assert_json_field "$resolve_output" '"state": "done"' "meeting action can be resolved"

home_output_after=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output_after" '"open_meeting_action_count": 0' "home clears resolved meeting action"

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
