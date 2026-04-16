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
claim_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-claim.sh 9100 "Packaged Runtime Brief Action Mission" --owns="file:harness/README.md" --brief)
brief_run_id=$(printf '%s\n' "$claim_brief_output" | sed -n 's/^Run: \(run-[^[:space:]]*\)$/\1/p' | head -n1)
drive_start_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh start "$brief_run_id" "Packaged brief execution started." --brief)
drive_report_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh report "$brief_run_id" "Packaged brief progress checkpoint." --brief)
drive_wait_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh wait "$brief_run_id" "Packaged brief dependency wait." --on=run-packaged-dependency --brief)
drive_block_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh block "$brief_run_id" "Packaged brief blocker." --brief)
drive_resume_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh resume "$brief_run_id" "Packaged brief blocker cleared." --brief)
signals_after_drive_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-signals.sh adjudicate --all)
drive_complete_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh complete "$brief_run_id" "Packaged brief flow complete." --brief)
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
teammate_assign_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-teammates.sh assign "$run_id" worker-packaged-brief packaged-helper "Check packaged runtime teammate assignment brief" --scope="packaged teammate brief" --owns=scope:packaged-teammate-brief --report="progress|verification" --brief)
teammate_brief_run_id=$(printf '%s\n' "$teammate_assign_brief_output" | sed -n 's/^Run: \(run-[^[:space:]]*\)$/\1/p' | head -n1)
teammate_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-teammates.sh get "$teammate_run_id" --brief)
signals_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-signals.sh active issue-9101 --brief)
signals_adjudicate_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-signals.sh adjudicate --all --brief)
teammate_report_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-teammates.sh report "$teammate_run_id" "Packaged teammate reported progress." --kind=progress --evidence=README.md --brief)
teammate_message_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-teammates.sh message "$teammate_run_id" "$teammate_brief_run_id" interface_change "Packaged teammate shared governed interface context." --evidence=README.md --brief)
teammate_escalate_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-teammates.sh escalate "$teammate_run_id" "Should packaged teammate escalate ambiguous release evidence?" --materialize --protocol=decision-record --brief)
maintenance_reconcile_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-maintenance.sh reconcile --brief)
maintenance_normalize_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-maintenance.sh normalize --brief)
meeting_open_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh meeting "$run_id" "Should packaged runtime open meetings as human briefings?" --protocol=decision-record --brief)
meeting_open_brief_id=$(printf '%s\n' "$meeting_open_brief_output" | sed -n 's/^Meeting: \(meeting-[^[:space:]]*\)$/\1/p' | head -n1)
(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh meeting-close "$run_id" "$meeting_open_brief_id" "Closed packaged meeting-open brief smoke." >/dev/null)
meeting_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh meeting "$run_id" "Should packaged runtime keep meeting briefings human-readable?")
meeting_id=$(printf '%s\n' "$meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
meeting_lock_decision_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh lock-decision "$meeting_id" "Keep packaged runtime meeting scope bounded." --reason="Packaged runtime smoke should prove decision locks stay human-readable." --authority=L2 --evidence=harness/profiles/erwin.yml --brief)
meeting_orchestrate_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh orchestrate "$meeting_id")
meeting_orchestrate_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh orchestrate "$meeting_id" --brief)
meeting_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh summary "$meeting_id" --brief)
meeting_authority_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh authority-audit "$meeting_id" --brief)
position_meeting_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh meeting "$run_id" "Should packaged runtime keep position briefings human-readable?")
position_meeting_id=$(printf '%s\n' "$position_meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
meeting_position_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh add-position "$position_meeting_id" "$run_id" independent_review support "Keep brief position text visible." --claim=packaged-position --dimension=factual --evidence=harness/profiles/erwin.yml --brief)
meeting_finding_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh add-finding "$position_meeting_id" accepted "Keep brief finding text visible." --reason="Packaged runtime smoke should prove finding briefings stay human-readable." --authority=L2 --evidence=harness/profiles/erwin.yml --run="$run_id" --role=packaged-leader --claim=packaged-position --dimension=factual --brief)
meeting_action_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh add-action "$position_meeting_id" "Keep brief action text visible." --brief)
meeting_action_id=$(printf '%s\n' "$meeting_action_brief_output" | sed -n 's/^Action: \(action-[^[:space:]]*\)$/\1/p' | head -n1)
meeting_action_resolution_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-meetings.sh resolve-action "$position_meeting_id" "$meeting_action_id" done --brief)
follow_up_claim_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-claim.sh 9102 "Packaged Runtime Follow-Up Mission" --owns="file:docs/harness-public-contract.md")
follow_up_run_id=$(printf '%s\n' "$follow_up_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
follow_up_packet_id="followup-packaged-runtime-smoke"
python3 - "$TEST_REPO/.harness/runs/registry.yml" "$follow_up_run_id" "$follow_up_packet_id" <<'PY'
import json
import sys
from pathlib import Path

registry_path, run_id, packet_id = sys.argv[1:]
data = json.loads(Path(registry_path).read_text())
packet = {
    "follow_up_packet_id": packet_id,
    "state": "open",
    "summary": "Resolve packaged adopted follow-up through drive.",
    "required_follow_up": "Resolve packaged adopted follow-up through drive.",
    "execution_kind": "packaged_runtime_follow_up",
}
for run in data.get("runs", []):
    if run.get("run_id") == run_id:
        run["adopted_from_meeting_id"] = "meeting-packaged-runtime-smoke"
        run["adopted_output_follow_up_packet_count"] = 1
        run["adopted_output_follow_up_packets"] = [packet]
        run["adopted_primary_follow_up_packet"] = packet
        run["next_expected_action"] = packet["required_follow_up"]
        break
else:
    raise SystemExit(f"run not found: {run_id}")
Path(registry_path).write_text(json.dumps(data, indent=2, sort_keys=True) + "\n")
PY
follow_up_resolve_brief_output=$(cd "$TEST_REPO" && PATH="$NO_GO_PATH" ./scripts/harness-drive.sh resolve-follow-up "$follow_up_run_id" "$follow_up_packet_id" done --brief)
assert_contains "$adapter_output" '"valid": true' "adapter validate works via packaged binary without Go"
assert_contains "$home_output" '"next_action": "No active runs. Claim a mission."' "home works via packaged binary without Go"
assert_contains "$home_brief_output" 'No active governed execution is in flight.' "home --brief works via packaged binary without Go"
assert_contains "$doctor_brief_output" 'Health: healthy' "doctor --brief works via packaged binary without Go"
assert_contains "$claim_brief_output" 'Harness Claim Briefing' "claim --brief works via packaged binary without Go"
assert_contains "$drive_start_brief_output" 'Action: start' "drive start --brief works via packaged binary without Go"
assert_contains "$drive_report_brief_output" 'Action: report' "drive report --brief works via packaged binary without Go"
assert_contains "$drive_report_brief_output" 'Packaged brief progress checkpoint.' "packaged drive report brief preserves reported progress without Go"
assert_contains "$drive_wait_brief_output" 'Action: wait' "drive wait --brief works via packaged binary without Go"
assert_contains "$drive_block_brief_output" 'Signal: sig-' "drive block --brief works via packaged binary without Go"
assert_contains "$drive_resume_brief_output" 'Action: resume' "drive resume --brief works via packaged binary without Go"
assert_contains "$signals_after_drive_brief_output" '"resolved_count": 1' "packaged signal adjudication clears drive block --brief signal"
assert_contains "$drive_complete_brief_output" 'Action: complete' "drive complete --brief works via packaged binary without Go"
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
assert_contains "$teammate_assign_brief_output" 'Harness Teammate Assignment Briefing' "teammates assign --brief works via packaged binary without Go"
assert_contains "$teammate_assign_brief_output" 'Teammate: worker-packaged-brief' "packaged teammate assignment brief preserves the delegated owner without Go"
assert_contains "$teammate_brief_output" 'Teammate: worker-packaged' "teammates get --brief works via packaged binary without Go"
assert_contains "$teammate_brief_output" "Assignment: mission=issue-9101, parent=$run_id, role=packaged-helper, state=ready" "packaged teammate brief still follows teammate assignment truth without Go"
assert_contains "$teammate_report_brief_output" 'Harness Teammate Report Briefing' "teammates report --brief works via packaged binary without Go"
assert_contains "$teammate_report_brief_output" 'Packaged teammate reported progress.' "packaged teammate report brief preserves reported progress without Go"
assert_contains "$teammate_message_brief_output" 'Harness Teammate Message Briefing' "teammates message --brief works via packaged binary without Go"
assert_contains "$teammate_message_brief_output" 'Packaged teammate shared governed interface context.' "packaged teammate message brief preserves governed message content without Go"
assert_contains "$teammate_escalate_brief_output" 'Harness Teammate Escalation Briefing' "teammates escalate --brief works via packaged binary without Go"
assert_contains "$teammate_escalate_brief_output" 'Should packaged teammate escalate ambiguous release evidence?' "packaged teammate escalation brief preserves the escalation question without Go"
assert_contains "$teammate_escalate_brief_output" 'Protocol: decision-record.' "packaged teammate escalation brief preserves explicit council protocol without Go"
assert_contains "$signals_brief_output" 'Primary: progress signal' "signals active --brief works via packaged binary without Go"
assert_contains "$signals_brief_output" 'worker-packaged-brief should keep the report fresh' "packaged signal brief still follows latest teammate assignment signal truth without Go"
assert_contains "$signals_adjudicate_brief_output" 'No active signal was adjudicated as superseded.' "signals adjudicate --brief works via packaged binary without Go"
assert_contains "$maintenance_reconcile_brief_output" 'No runtime repair was needed.' "maintenance reconcile --brief works via packaged binary without Go"
assert_contains "$maintenance_normalize_brief_output" 'No runtime normalization was needed.' "maintenance normalize --brief works via packaged binary without Go"
assert_contains "$meeting_open_brief_output" 'Harness Meeting Open Briefing' "drive meeting --brief works via packaged binary without Go"
assert_contains "$meeting_open_brief_output" 'Protocol: decision-record.' "packaged drive meeting brief preserves explicit council protocol without Go"
assert_contains "$meeting_lock_decision_brief_output" 'Harness Meeting Decision Lock Briefing' "meetings lock-decision --brief works via packaged binary without Go"
assert_contains "$meeting_lock_decision_brief_output" 'Leader should keep the council bounded to: Keep packaged runtime meeting scope bounded.' "packaged meeting decision lock brief still follows council-scope truth without Go"
assert_contains "$meeting_orchestrate_output" '"new_phase": "setup"' "meeting orchestration works via packaged binary without Go"
assert_contains "$meeting_orchestrate_brief_output" 'Harness Meeting Orchestration Briefing' "meetings orchestrate --brief works via packaged binary without Go"
assert_contains "$meeting_orchestrate_brief_output" 'Add packet/participant prerequisites before independent review.' "packaged meeting orchestration brief still follows live council requirements without Go"
assert_contains "$meeting_brief_output" "Meeting: $meeting_id" "meetings summary --brief works via packaged binary without Go"
assert_contains "$meeting_brief_output" 'Add packet/participant prerequisites before independent review.' "packaged meeting brief still follows live council requirements without Go"
assert_contains "$meeting_authority_brief_output" "Meeting: $meeting_id" "meetings authority-audit --brief works via packaged binary without Go"
assert_contains "$meeting_authority_brief_output" "No same-claim authority conflict is currently recorded on meeting $meeting_id." "packaged meeting authority brief still follows council-authority truth without Go"
assert_contains "$meeting_position_brief_output" 'Harness Meeting Position Briefing' "meetings add-position --brief works via packaged binary without Go"
assert_contains "$meeting_position_brief_output" 'Keep brief position text visible.' "packaged meeting position brief preserves free-text council position content"
assert_contains "$meeting_finding_brief_output" 'Harness Meeting Finding Briefing' "meetings add-finding --brief works via packaged binary without Go"
assert_contains "$meeting_finding_brief_output" 'Keep brief finding text visible.' "packaged meeting finding brief preserves free-text council finding content"
assert_contains "$meeting_action_brief_output" 'Harness Meeting Action Briefing' "meetings add-action --brief works via packaged binary without Go"
assert_contains "$meeting_action_brief_output" 'Keep brief action text visible.' "packaged meeting action brief preserves free-text council action content"
assert_contains "$meeting_action_resolution_brief_output" 'Harness Meeting Action Resolution Briefing' "meetings resolve-action --brief works via packaged binary without Go"
assert_contains "$meeting_action_resolution_brief_output" 'Resolution Scope: state=done' "packaged meeting action resolution brief preserves resolved state"
assert_contains "$follow_up_resolve_brief_output" 'Harness Follow-Up Resolution Briefing' "drive resolve-follow-up --brief works via packaged binary without Go"
assert_contains "$follow_up_resolve_brief_output" 'Resolution Scope: state=done, delegated=false' "packaged follow-up resolution brief preserves resolved state"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
