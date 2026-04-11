#!/usr/bin/env bash
# harness-runtime-dist-smoke-test.sh — validates packaged-binary runtime use without Go in PATH
#
# USAGE:
#   ./scripts/harness-runtime-dist-smoke-test.sh

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

TMPDIR_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/harness-runtime-dist-smoke.XXXXXX")
trap 'rm -rf "$TMPDIR_ROOT"' EXIT

PASS=0
FAIL=0

pass() { PASS=$((PASS + 1)); echo "  PASS  $1"; }
fail() { FAIL=$((FAIL + 1)); echo "  FAIL  $1"; }
section() { echo ""; echo "=== $1 ==="; }

TEST_REPO="$TMPDIR_ROOT/repo"
mkdir -p "$TEST_REPO/scripts" "$TEST_REPO/harness/profiles" "$TEST_REPO/harness/templates" "$TEST_REPO/harness/policies" "$TEST_REPO/harness/bin"

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
    git config user.email "runtime-dist@example.com" &&
    git config user.name "Harness Runtime Dist Test" &&
    git checkout -q -b main &&
    chmod +x scripts/*.sh install-codex.sh &&
    printf '# Harness Runtime Dist Smoke Test\n' >README.md &&
    git add . &&
    git commit -q -m "Initial commit"
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
  if echo "$text" | grep -q "$needle"; then
    pass "$label"
  else
    fail "$label"
    echo "$text"
  fi
}

path_without_go() {
  local go_dir
  go_dir="$(dirname "$(command -v go)")"
  local result=""
  local first=1
  IFS=':' read -r -a parts <<<"${PATH:-}"
  for part in "${parts[@]}"; do
    [[ "$part" == "$go_dir" ]] && continue
    [[ -z "$part" ]] && continue
    if [[ "$first" -eq 0 ]]; then
      result="${result}:"
    fi
    first=0
    result="${result}${part}"
  done
  printf '%s\n' "$result"
}

section "Build packaged binary"
build_output=$(run_in_repo ./scripts/harness-build-runtime.sh)
assert_contains "$build_output" '"built"' "runtime build script reports built binaries"

section "Run without Go in PATH"
NO_GO_PATH="$(path_without_go)"
adapter_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-adapter.sh validate)
home_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-home.sh)
home_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-home.sh --brief)
doctor_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-doctor.sh --brief)
claim_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-claim.sh 9101 "Packaged Runtime Mission" --owns="file:README.md")
schedule_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-missions.sh schedule --brief)
mission_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-missions.sh get issue-9101 --brief)
mission_graph_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-missions.sh graph --brief)
mission_authority_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-missions.sh authority-audit issue-9101 --brief)
mission_graph_authority_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-missions.sh authority-audit --all --brief)
run_id=$(printf '%s\n' "$claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
run_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-runs.sh get "$run_id" --brief)
run_authority_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-runs.sh authority-audit "$run_id" --brief)
teammate_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-teammates.sh assign "$run_id" worker-packaged packaged-helper "Check packaged runtime teammate brief" --scope="packaged teammate" --owns=scope:packaged-teammate --report="progress|verification")
teammate_run_id=$(printf '%s\n' "$teammate_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
teammate_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-teammates.sh get "$teammate_run_id" --brief)
signals_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-signals.sh active issue-9101 --brief)
signals_adjudicate_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-signals.sh adjudicate --all --brief)
meeting_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh meeting "$run_id" "Should packaged runtime keep meeting briefings human-readable?")
meeting_id=$(printf '%s\n' "$meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
meeting_orchestrate_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh orchestrate "$meeting_id")
meeting_orchestrate_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh orchestrate "$meeting_id" --brief)
meeting_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh summary "$meeting_id" --brief)
meeting_authority_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh authority-audit "$meeting_id" --brief)
assert_contains "$adapter_output" '"valid": true' "adapter validate works via packaged binary without Go"
assert_contains "$home_output" '"next_action": "No active runs. Claim a mission."' "home works via packaged binary without Go"
assert_contains "$home_brief_output" 'No active governed execution is in flight.' "home --brief works via packaged binary without Go"
assert_contains "$doctor_brief_output" 'Health: healthy' "doctor --brief works via packaged binary without Go"
assert_contains "$claim_output" '"mission_id": "issue-9101"' "claim works via packaged binary without Go"
assert_contains "$schedule_brief_output" 'Foreground execution is currently centered on mission issue-9101.' "missions schedule --brief works via packaged binary without Go"
assert_contains "$schedule_brief_output" 'Foreground batch is issue-9101 under solo mode.' "packaged scheduler brief still follows foreground scheduler truth without Go"
assert_contains "$mission_brief_output" 'Mission: issue-9101' "missions get --brief works via packaged binary without Go"
assert_contains "$mission_brief_output" 'Scheduler currently selected this mission as the foreground mission.' "packaged mission brief still follows scheduler truth without Go"
assert_contains "$mission_graph_brief_output" 'Topology currently has one standalone active mission: issue-9101.' "missions graph --brief works via packaged binary without Go"
assert_contains "$mission_authority_brief_output" 'Mission: issue-9101' "missions authority-audit --brief works via packaged binary without Go"
assert_contains "$mission_authority_brief_output" 'Mission issue-9101 is aligned with its authoritative run-derived control plane.' "packaged mission authority brief still follows mission-control truth without Go"
assert_contains "$mission_graph_authority_brief_output" 'No connected mission bundle currently needs graph-level authority handling.' "missions authority-audit --all --brief works via packaged binary without Go"
assert_contains "$run_brief_output" "Run: $run_id" "runs get --brief works via packaged binary without Go"
assert_contains "$run_brief_output" 'This run is the current foreground scheduler run and still owns mission control for issue-9101.' "packaged run brief still follows scheduler truth without Go"
assert_contains "$run_authority_brief_output" "Run: $run_id" "runs authority-audit --brief works via packaged binary without Go"
assert_contains "$run_authority_brief_output" "Run $run_id is still progressing toward its desired state." "packaged run authority brief still follows desired-state progress truth without Go"
assert_contains "$teammate_brief_output" 'Teammate: worker-packaged' "teammates get --brief works via packaged binary without Go"
assert_contains "$teammate_brief_output" "Assignment: mission=issue-9101, parent=$run_id, role=packaged-helper, state=ready" "packaged teammate brief still follows teammate assignment truth without Go"
assert_contains "$signals_brief_output" 'Primary: progress signal' "signals active --brief works via packaged binary without Go"
assert_contains "$signals_brief_output" 'worker-packaged should keep the report fresh' "packaged signal brief still follows teammate signal truth without Go"
assert_contains "$signals_adjudicate_brief_output" 'No active signal was adjudicated as superseded.' "signals adjudicate --brief works via packaged binary without Go"
assert_contains "$meeting_orchestrate_output" '"new_phase": "setup"' "meeting orchestration works via packaged binary without Go"
assert_contains "$meeting_orchestrate_brief_output" 'Harness Meeting Orchestration Briefing' "meetings orchestrate --brief works via packaged binary without Go"
assert_contains "$meeting_orchestrate_brief_output" 'Add packet/participant prerequisites before independent review.' "packaged meeting orchestration brief still follows live council requirements without Go"
assert_contains "$meeting_brief_output" "Meeting: $meeting_id" "meetings summary --brief works via packaged binary without Go"
assert_contains "$meeting_brief_output" 'Add packet/participant prerequisites before independent review.' "packaged meeting brief still follows live council requirements without Go"
assert_contains "$meeting_authority_brief_output" "Meeting: $meeting_id" "meetings authority-audit --brief works via packaged binary without Go"
assert_contains "$meeting_authority_brief_output" "No same-claim authority conflict is currently recorded on meeting $meeting_id." "packaged meeting authority brief still follows council-authority truth without Go"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
