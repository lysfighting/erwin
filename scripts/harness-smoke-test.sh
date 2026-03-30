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

rewrite_json_in_file() {
  local file="$1"
  local perl_expr="$2"
  perl -0pi -e "$perl_expr" "$file"
}

set_json_field_in_file() {
  local file="$1"
  local path="$2"
  local value="$3"
  python3 - "$file" "$path" "$value" <<'PY'
import json
import sys
from pathlib import Path

file_path, path, value = sys.argv[1:]
data = json.loads(Path(file_path).read_text())
target = data
parts = path.split(".")
for key in parts[:-1]:
    target = target[key]
target[parts[-1]] = value
Path(file_path).write_text(json.dumps(data, indent=2, sort_keys=True) + "\n")
PY
}

assert_json_field() {
  local json="$1"
  local needle="$2"
  local label="$3"
  if grep -q "$needle" <<<"$json"; then
    pass "$label"
  else
    fail "$label"
    echo "$json"
  fi
}

section "Claim"
claim_output=$(run_in_repo bash scripts/harness-claim.sh 20001 "Harness Smoke Mission" --owns="file:CHANGELOG.md")
assert_json_field "$claim_output" '"mission_id": "issue-20001"' "claim creates mission"
assert_json_field "$claim_output" '"run_id": "run-' "claim returns run id"
assert_json_field "$claim_output" 'issue-20001.yml' "claim materializes mission path"

run_id=$(printf '%s\n' "$claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
run_branch=$(printf '%s\n' "$claim_output" | sed -n 's/^[[:space:]]*"branch": "\(mission\/[^"]*\)",$/\1/p' | head -n1)
contract_path=$(printf '%s\n' "$claim_output" | sed -n 's/^[[:space:]]*"contract_path": "\(.*issue-[^"]*\.yml\)",$/\1/p' | head -n1)
registry_path="$TEST_REPO/.harness/runs/registry.yml"

section "Mission registry"
missions_output=$(run_in_repo bash scripts/harness-missions.sh list)
assert_json_field "$missions_output" '"mission_id": "issue-20001"' "missions list includes claimed mission"
mission_get_output=$(run_in_repo bash scripts/harness-missions.sh get issue-20001)
assert_json_field "$mission_get_output" '"state": "claimed"' "mission starts in claimed state"
assert_json_field "$mission_get_output" '"active_run_ids":' "mission stores active run linkage"
assert_json_field "$mission_get_output" "\"authoritative_run_id\": \"$run_id\"" "mission records the initial leader run as authoritative control run"
mission_authority_baseline=$(run_in_repo bash scripts/harness-missions.sh authority-audit issue-20001)
assert_json_field "$mission_authority_baseline" "\"authoritative_run_id\": \"$run_id\"" "mission authority audit reports the leader run as authoritative"
assert_json_field "$mission_authority_baseline" '"mission_gap_count": 0' "mission authority audit starts without control-plane drift"
assert_json_field "$mission_authority_baseline" '"adjudication_count": 0' "mission authority audit starts without cross-run adjudications"
authority_baseline=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_id")
assert_json_field "$authority_baseline" '"contract_gap_count": 0' "initial run authority audit starts without contract drift"

section "Scheduler baseline"
schedule_solo_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_solo_output" '"execution_mode": "solo"' "scheduler recommends solo execution for one runnable mission"
assert_json_field "$schedule_solo_output" '"selected_mission_count": 1' "scheduler emits a selected mission batch"
assert_json_field "$schedule_solo_output" '"snapshot_path":' "scheduler writes a durable snapshot"
assert_json_field "$schedule_solo_output" '"mission_id": "issue-20001"' "scheduler includes the claimed mission in its classification output"

section "Fresh signal rescheduling"
fresh_signal_output=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-fresh-follow-up",
  "kind": "meeting_required",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$run_id",
  "mission_id": "issue-20001",
  "applies_to": ["issue-20001"],
  "scope_type": "mission",
  "title": "Fresh coordination signal",
  "summary": "Review the fresh coordination evidence before continuing.",
  "evidence_refs": ["$contract_path"],
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-29T00:00:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$fresh_signal_output" '"kind": "requires_meeting"' "fresh follow-up signal can be emitted directly"
schedule_fresh_signal_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_fresh_signal_output" '"decision_kind": "review_follow_up_signals"' "scheduler escalates fresh coordination signals into an explicit decision"
assert_json_field "$schedule_fresh_signal_output" '"action_kind": "pause_for_fresh_follow_up_signal"' "scheduler emits a reschedule action for fresh coordination signals"
assert_json_field "$schedule_fresh_signal_output" '"follow_up_signal_count": 1' "scheduler counts fresh follow-up signals in mission classification"
resolved_fresh_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-fresh-follow-up resolved "Fresh coordination signal reviewed.")
assert_json_field "$resolved_fresh_signal" '"state": "resolved"' "fresh follow-up signal can be resolved after scheduler coverage"

review_feedback_signal=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-review-feedback",
  "kind": "review_feedback",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$run_id",
  "mission_id": "issue-20001",
  "applies_to": ["issue-20001"],
  "scope_type": "mission",
  "title": "Fresh review feedback",
  "summary": "Address the fresh review feedback before widening execution.",
  "evidence_refs": ["$contract_path"],
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-29T00:05:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$review_feedback_signal" '"kind": "review_feedback"' "review-feedback signal can be emitted directly"
schedule_review_feedback_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_review_feedback_output" '"decision_kind": "review_follow_up_signals"' "scheduler treats fresh review feedback as governed follow-up decision work"
assert_json_field "$schedule_review_feedback_output" '"action_kind": "pause_for_fresh_follow_up_signal"' "scheduler pauses foreground execution for fresh review feedback"
assert_json_field "$schedule_review_feedback_output" 'Review the fresh review feedback and adopt it into governed review state before widening execution.' "scheduler surfaces review-feedback-specific guidance"
resolved_review_feedback_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-review-feedback resolved "Fresh review feedback covered.")
assert_json_field "$resolved_review_feedback_signal" '"state": "resolved"' "review-feedback signal can be resolved after scheduler coverage"

priority_change_signal=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-priority-change",
  "kind": "priority_change",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$run_id",
  "mission_id": "issue-20001",
  "applies_to": ["issue-20001"],
  "scope_type": "mission",
  "title": "Fresh priority change",
  "summary": "Mission priority changed and should be adopted before widening execution.",
  "evidence_refs": ["$contract_path"],
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-29T00:07:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$priority_change_signal" '"kind": "priority_change"' "priority-change signal can be emitted directly"
schedule_priority_change_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_priority_change_output" '"decision_kind": "review_follow_up_signals"' "scheduler treats fresh priority-change signals as governed follow-up decision work"
assert_json_field "$schedule_priority_change_output" '"action_kind": "pause_for_fresh_follow_up_signal"' "scheduler pauses foreground execution for fresh priority-change signals"
assert_json_field "$schedule_priority_change_output" 'Review the fresh priority-change signal and adopt it into mission state before widening execution.' "scheduler surfaces priority-change-specific guidance"
resolved_priority_change_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-priority-change resolved "Fresh priority change covered.")
assert_json_field "$resolved_priority_change_signal" '"state": "resolved"' "priority-change signal can be resolved after scheduler coverage"

risk_raised_signal=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-risk-raised",
  "kind": "risk_raised",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$run_id",
  "mission_id": "issue-20001",
  "applies_to": ["issue-20001"],
  "scope_type": "mission",
  "title": "Fresh risk raised",
  "summary": "New execution risk should be reviewed before widening execution.",
  "evidence_refs": ["$contract_path"],
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-29T00:08:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$risk_raised_signal" '"kind": "risk_raised"' "risk-raised signal can be emitted directly"
schedule_risk_raised_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_risk_raised_output" '"decision_kind": "review_follow_up_signals"' "scheduler treats fresh risk-raised signals as governed follow-up decision work"
assert_json_field "$schedule_risk_raised_output" '"action_kind": "pause_for_fresh_follow_up_signal"' "scheduler pauses foreground execution for fresh risk-raised signals"
assert_json_field "$schedule_risk_raised_output" 'Review the fresh risk-raised signal and either tighten policy or convene a governed decision before widening execution.' "scheduler surfaces risk-raised-specific guidance"
resolved_risk_raised_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-risk-raised resolved "Fresh risk signal covered.")
assert_json_field "$resolved_risk_raised_signal" '"state": "resolved"' "risk-raised signal can be resolved after scheduler coverage"

rebase_signal=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-rebase-required",
  "kind": "requires_rebase",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$run_id",
  "mission_id": "issue-20001",
  "applies_to": ["issue-20001"],
  "scope_type": "mission",
  "title": "Fresh rebase required",
  "summary": "Rebase before continuing this run.",
  "evidence_refs": ["$contract_path"],
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-29T00:10:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$rebase_signal" '"kind": "requires_rebase"' "requires-rebase signal can be emitted directly"
schedule_rebase_signal_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_rebase_signal_output" '"classification": "blocked-by-policy"' "scheduler blocks mission execution behind a fresh requires-rebase signal"
assert_json_field "$schedule_rebase_signal_output" '"action_kind": "pause_for_fresh_rebase_signal"' "scheduler emits a dedicated reschedule action for fresh requires-rebase signals"
assert_json_field "$schedule_rebase_signal_output" 'Rebase or refresh the governed run state before continuing this mission.' "scheduler surfaces rebase-specific guidance"
resolved_rebase_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-rebase-required resolved "Fresh rebase requirement covered.")
assert_json_field "$resolved_rebase_signal" '"state": "resolved"' "requires-rebase signal can be resolved after scheduler coverage"

section "Fresh follow-up portfolio urgency"
follow_up_priority_claim=$(run_in_repo bash scripts/harness-claim.sh 20035 "Harness Follow-Up Priority Signal" --owns="file:docs/follow-up-priority.md")
assert_json_field "$follow_up_priority_claim" '"mission_id": "issue-20035"' "priority-signal mission can materialize"
follow_up_priority_run=$(printf '%s\n' "$follow_up_priority_claim" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
follow_up_risk_claim=$(run_in_repo bash scripts/harness-claim.sh 20036 "Harness Follow-Up Risk Signal" --owns="file:docs/follow-up-risk.md")
assert_json_field "$follow_up_risk_claim" '"mission_id": "issue-20036"' "risk-signal mission can materialize"
follow_up_risk_run=$(printf '%s\n' "$follow_up_risk_claim" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
follow_up_priority_signal=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-follow-up-priority",
  "kind": "priority_change",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$follow_up_priority_run",
  "mission_id": "issue-20035",
  "applies_to": ["issue-20035"],
  "scope_type": "mission",
  "title": "Follow-up priority change",
  "summary": "Priority should be re-evaluated before widening execution.",
  "evidence_refs": ["$contract_path"],
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-29T00:12:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$follow_up_priority_signal" '"signal_id": "sig-smoke-follow-up-priority"' "priority-signal portfolio signal can be emitted"
follow_up_risk_signal=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-follow-up-risk",
  "kind": "risk_raised",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$follow_up_risk_run",
  "mission_id": "issue-20036",
  "applies_to": ["issue-20036"],
  "scope_type": "mission",
  "title": "Follow-up risk raised",
  "summary": "Risk increased and should outrank lower-urgency follow-up work.",
  "evidence_refs": ["$contract_path"],
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-29T00:13:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$follow_up_risk_signal" '"signal_id": "sig-smoke-follow-up-risk"' "risk-signal portfolio signal can be emitted"
follow_up_portfolio_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$follow_up_portfolio_schedule" '"follow_up_signal_urgency_rank": 4' "risk-raised follow-up packets expose the highest urgency rank"
assert_json_field "$follow_up_portfolio_schedule" '"follow_up_signal_urgency_rank": 2' "priority-change follow-up packets expose a lower urgency rank"
assert_json_field "$follow_up_portfolio_schedule" '"fresh_follow_up_signal_kinds":' "scheduler projects fresh follow-up signal kinds into the portfolio"
assert_json_field "$follow_up_portfolio_schedule" '"effective_risk_level": "high"' "risk-raised follow-up packets escalate effective risk before adoption"
primary_follow_up_mission=$(printf '%s\n' "$follow_up_portfolio_schedule" | perl -0ne 'print "$1\n" if /"primary_required_human_decision": \{[^{}]*"mission_id": "([^"]+)"/s' | head -n1)
if [[ "$primary_follow_up_mission" == "issue-20036" ]]; then
  pass "risk-raised follow-up outranks lower-urgency priority-change work in the portfolio"
else
  fail "risk-raised follow-up outranks lower-urgency priority-change work in the portfolio"
  echo "$follow_up_portfolio_schedule"
fi
assert_json_field "$follow_up_portfolio_schedule" '"human_decision_queue_position": 2' "lower-urgency follow-up work stays deferred behind the risk-raised packet"
resolved_follow_up_priority_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-follow-up-priority resolved "Priority follow-up covered.")
assert_json_field "$resolved_follow_up_priority_signal" '"state": "resolved"' "priority-signal portfolio signal can be resolved"
resolved_follow_up_risk_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-follow-up-risk resolved "Risk follow-up covered.")
assert_json_field "$resolved_follow_up_risk_signal" '"state": "resolved"' "risk-signal portfolio signal can be resolved"
resolved_follow_up_priority_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$follow_up_priority_run" completed "Closed priority-signal mission after follow-up portfolio coverage.")
assert_json_field "$resolved_follow_up_priority_run" '"state": "completed"' "priority-signal mission can close after portfolio coverage"
resolved_follow_up_risk_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$follow_up_risk_run" completed "Closed risk-signal mission after follow-up portfolio coverage.")
assert_json_field "$resolved_follow_up_risk_run" '"state": "completed"' "risk-signal mission can close after portfolio coverage"

section "Maintenance escalation"
rewrite_json_in_file "$registry_path" 's/"last_reported_at": "[^"]*"/"last_reported_at": "2026-03-01T00:00:00Z"/'
stale_signal_output=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-stale-follow-up",
  "kind": "interface_change",
  "source_type": "teammate",
  "source_id": "worker-auth",
  "run_id": "$run_id",
  "mission_id": "issue-20001",
  "applies_to": ["issue-20001"],
  "scope_type": "mission",
  "title": "Stale coordination signal",
  "summary": "Old interface-change note that should now be treated as maintenance work.",
  "evidence_refs": ["$contract_path"],
  "authority_level": "L2",
  "confidence": "medium",
  "emitted_at": "2026-03-01T00:00:00Z",
  "fresh_until": "2026-03-02T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$stale_signal_output" '"signal_id": "sig-smoke-stale-follow-up"' "stale follow-up signal can be emitted directly"
schedule_maintenance_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_maintenance_output" '"maintenance_action_count": 2' "scheduler surfaces stale-state maintenance work"
assert_json_field "$schedule_maintenance_output" '"action_kind": "refresh_stale_run"' "scheduler emits maintenance action for stale runs"
assert_json_field "$schedule_maintenance_output" '"action_kind": "resolve_stale_signal"' "scheduler emits maintenance action for expired signals"
assert_json_field "$schedule_maintenance_output" '"next_action": "Resolve or refresh the stale signal before letting it influence future scheduling."' "scheduler prioritizes maintenance guidance when stale state coexists with runnable work"
home_output_maintenance=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output_maintenance" '"scheduler_maintenance_action_count": 2' "home surfaces scheduler maintenance work"
assert_json_field "$home_output_maintenance" '"scheduler_maintenance_action_tally":' "home surfaces scheduler maintenance tally"
run_refresh_output=$(run_in_repo bash scripts/harness-drive.sh report "$run_id" "Refresh stale run after maintenance escalation smoke.")
assert_json_field "$run_refresh_output" '"latest_summary": "Refresh stale run after maintenance escalation smoke."' "leader can refresh stale run after maintenance coverage"
resolved_stale_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-stale-follow-up resolved "Stale follow-up signal retired after maintenance escalation coverage.")
assert_json_field "$resolved_stale_signal" '"state": "resolved"' "stale follow-up signal can be retired after maintenance coverage"

section "Mission graph authority"
second_claim_output=$(run_in_repo bash scripts/harness-claim.sh 20002 "Harness Conflict Mission" --owns="file:CHANGELOG.md" --allow-overlap)
assert_json_field "$second_claim_output" '"mission_id": "issue-20002"' "second claim can materialize a conflicting mission under explicit overlap allowance"
second_run_id=$(printf '%s\n' "$second_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
graph_authority_output=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_output" '"cross_mission_conflict_count": 1' "mission graph authority audit reports one cross-mission control conflict"
assert_json_field "$graph_authority_output" '"bundle_count": 1' "mission graph authority groups overlapping missions into one bundle"
assert_json_field "$graph_authority_output" '"bundle_state": "conflicted"' "mission graph authority marks the overlapping mission bundle as conflicted"
assert_json_field "$graph_authority_output" '"resolution_kind": "bundle_conflict"' "mission graph authority gives the conflicted bundle a bundle-level resolution kind"
assert_json_field "$graph_authority_output" '"root_mission_ids":' "mission graph authority exposes root missions for the conflicted bundle"
assert_json_field "$graph_authority_output" '"required_action": "Narrow selectors, serialize missions, or add a governing dependency/council before parallel execution."' "mission graph authority gives the conflicted bundle a leader-facing next action"
assert_json_field "$graph_authority_output" '"selector_refs":' "mission graph authority surfaces shared selectors for the conflicted bundle"
assert_json_field "$graph_authority_output" '"resolution_kind": "parallel_mission_ownership_conflict"' "mission graph authority audit explains the multi-mission ownership conflict"
doctor_output_graph_conflict=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output_graph_conflict" '"kind": "mission-graph-authority-conflict"' "doctor surfaces cross-mission authority conflict"
resolved_second_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$second_run_id" completed "Resolved cross-mission smoke overlap.")
assert_json_field "$resolved_second_run" '"state": "completed"' "second conflicting mission can be closed to clear mission-graph conflict"
graph_authority_cleared=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_cleared" '"cross_mission_conflict_count": 0' "mission graph authority audit clears after conflicting mission closes"

dependent_claim_output=$(run_in_repo bash scripts/harness-claim.sh 20003 "Harness Dependent Mission" --owns="file:CHANGELOG.md" --depends-on="issue-20001" --allow-overlap)
assert_json_field "$dependent_claim_output" '"mission_id": "issue-20003"' "dependent claim can materialize an explicitly governed overlap"
dependent_run_id=$(printf '%s\n' "$dependent_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
graph_authority_dependency=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_dependency" '"cross_mission_conflict_count": 0' "dependency-governed overlap does not remain an unresolved cross-mission conflict"
assert_json_field "$graph_authority_dependency" '"cross_mission_adjudication_count": 1' "mission graph authority audit records one dependency-governed ownership adjudication"
assert_json_field "$graph_authority_dependency" '"governed_bundle_count": 1' "mission graph authority reports one governed mission bundle"
assert_json_field "$graph_authority_dependency" '"bundle_state": "governed"' "mission graph authority marks the dependency-governed mission bundle as governed"
assert_json_field "$graph_authority_dependency" '"resolution_kind": "single_root_governed_bundle"' "mission graph authority gives the governed bundle a bundle-level resolution kind"
assert_json_field "$graph_authority_dependency" "\"governing_mission_id\": \"issue-20001\"" "mission graph authority identifies the governing mission for the governed bundle"
assert_json_field "$graph_authority_dependency" "\"required_action\": \"Keep bundle authority anchored to issue-20001 while dependent missions stay bounded.\"" "mission graph authority gives the governed bundle a leader-facing next action"
assert_json_field "$graph_authority_dependency" '"resolution_kind": "dependency_governed_mission_ownership"' "mission graph authority audit explains dependency-governed ownership authority"
assert_json_field "$graph_authority_dependency" "\"winning_source\": \"mission:issue-20001\"" "dependency-governed overlap keeps authority on the upstream mission"
scheduler_dependency_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$scheduler_dependency_output" '"mission_id": "issue-20003"' "scheduler includes the dependency-governed mission"
assert_json_field "$scheduler_dependency_output" '"classification": "blocked-by-dependency"' "scheduler classifies unsatisfied dependent missions as blocked-by-dependency"
assert_json_field "$scheduler_dependency_output" '"action_kind": "move_to_waiting"' "scheduler emits a reschedule action when a dependency-blocked run is still foregrounded"
resolved_dependent_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$dependent_run_id" completed "Resolved dependency-governed overlap smoke path.")
assert_json_field "$resolved_dependent_run" '"state": "completed"' "dependent overlapping mission can be closed after adjudication coverage"

mkdir -p "$TEST_REPO/docs/release"
touch "$TEST_REPO/docs/release/CHANGELOG.md"
broader_claim_output=$(run_in_repo bash scripts/harness-claim.sh 20005 "Harness Broad Release Docs Mission" --owns="dir:docs/release" --allow-overlap)
assert_json_field "$broader_claim_output" '"mission_id": "issue-20005"' "broader release-doc mission can claim a directory boundary under explicit overlap allowance"
broader_run_id=$(printf '%s\n' "$broader_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
narrower_claim_output=$(run_in_repo bash scripts/harness-claim.sh 20006 "Harness Narrow Release Changelog Mission" --owns="file:docs/release/CHANGELOG.md" --allow-overlap)
assert_json_field "$narrower_claim_output" '"mission_id": "issue-20006"' "narrow release-doc mission can claim a more specific file boundary under explicit overlap allowance"
narrower_run_id=$(printf '%s\n' "$narrower_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
graph_authority_specificity=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_specificity" '"resolution_kind": "selector_specificity_governed_mission_ownership"' "mission graph authority adjudicates broader-vs-narrower mission overlap by selector specificity"
assert_json_field "$graph_authority_specificity" "\"winning_source\": \"mission:issue-20006\"" "mission graph authority keeps the more specific mission as the selector winner"
assert_json_field "$graph_authority_specificity" "\"losing_source\": \"mission:issue-20005\"" "mission graph authority marks the broader mission as losing the shared selector boundary"
assert_json_field "$graph_authority_specificity" '"resolution_kind": "single_winner_governed_bundle"' "mission graph authority upgrades a specificity-resolved bundle into a single-winner governed bundle"
assert_json_field "$graph_authority_specificity" "\"governing_mission_id\": \"issue-20006\"" "mission graph authority exposes the specific governing mission even without a dependency edge"
assert_json_field "$graph_authority_specificity" "\"required_action\": \"Keep bundle authority anchored to issue-20006, and narrow or transfer broader overlapping missions before widening execution.\"" "mission graph authority gives the specificity-governed bundle a concrete leader-facing action"
resolved_broader_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$broader_run_id" completed "Resolved broader specificity-governed overlap smoke path.")
assert_json_field "$resolved_broader_run" '"state": "completed"' "broader specificity-governed mission can be closed after adjudication coverage"
resolved_narrower_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$narrower_run_id" completed "Resolved narrower specificity-governed overlap smoke path.")
assert_json_field "$resolved_narrower_run" '"state": "completed"' "narrower specificity-governed mission can be closed after adjudication coverage"

pair_governor_claim_a=$(run_in_repo bash scripts/harness-claim.sh 20014 "Harness Pair Governor A" --owns="file:docs/README.md" --allow-overlap)
assert_json_field "$pair_governor_claim_a" '"mission_id": "issue-20014"' "pairwise governor claim A can materialize"
pair_governor_run_a=$(printf '%s\n' "$pair_governor_claim_a" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
pair_governor_claim_b=$(run_in_repo bash scripts/harness-claim.sh 20015 "Harness Pair Governor B" --owns="file:docs/README.md" --allow-overlap)
assert_json_field "$pair_governor_claim_b" '"mission_id": "issue-20015"' "pairwise governor claim B can materialize"
pair_governor_run_b=$(printf '%s\n' "$pair_governor_claim_b" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)

pair_governor_update_a=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$pair_governor_run_a"
{
  "priority": "P0",
  "risk_level": "high"
}
EOF
)
assert_json_field "$pair_governor_update_a" '"priority": "P0"' "pairwise governor candidate can raise its authority tuple"
pair_governor_update_b=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$pair_governor_run_b"
{
  "priority": "P2",
  "risk_level": "low"
}
EOF
)
assert_json_field "$pair_governor_update_b" '"priority": "P2"' "pairwise governor peer can keep a lower authority tuple"

graph_authority_pair_governor=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_pair_governor" '"resolution_kind": "bundle_governor_adjudicated"' "pairwise mission bundle authority can now adjudicate a governing mission without a dependency edge"
assert_json_field "$graph_authority_pair_governor" "\"governing_mission_id\": \"issue-20014\"" "pairwise mission bundle authority anchors governance to the strongest mission"
assert_json_field "$graph_authority_pair_governor" '"bundle_state": "governed"' "pairwise mission bundle authority marks the adjudicated pair as governed"
assert_json_field "$graph_authority_pair_governor" '"authority_basis":' "pairwise mission bundle authority exposes why the governor won"
scheduler_pair_governor_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$scheduler_pair_governor_output" "\"bundle_governing_mission_id\": \"issue-20014\"" "scheduler surfaces the governing mission for pairwise governed bundles"
assert_json_field "$scheduler_pair_governor_output" '"bundle_role": "governed"' "scheduler marks the losing mission in a governed bundle"
assert_json_field "$scheduler_pair_governor_output" '"action_kind": "defer_until_bundle_governor_clears"' "scheduler defers losing missions behind the current bundle governor"
resolved_pair_governor_a=$(run_in_repo bash scripts/harness-runs.sh set-state "$pair_governor_run_a" completed "Closed pairwise governed bundle winner.")
assert_json_field "$resolved_pair_governor_a" '"state": "completed"' "pairwise governed bundle winner can be closed"
resolved_pair_governor_b=$(run_in_repo bash scripts/harness-runs.sh set-state "$pair_governor_run_b" completed "Closed pairwise governed bundle peer.")
assert_json_field "$resolved_pair_governor_b" '"state": "completed"' "pairwise governed bundle peer can be closed"

multi_governor_claim_a=$(run_in_repo bash scripts/harness-claim.sh 20008 "Harness Bundle Governor Mission" --owns="file:install-codex.sh" --allow-overlap)
assert_json_field "$multi_governor_claim_a" '"mission_id": "issue-20008"' "first multi-mission bundle claim can materialize"
multi_governor_run_a=$(printf '%s\n' "$multi_governor_claim_a" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
multi_governor_claim_b=$(run_in_repo bash scripts/harness-claim.sh 20009 "Harness Bundle Governor Peer" --owns="file:install-codex.sh" --allow-overlap)
assert_json_field "$multi_governor_claim_b" '"mission_id": "issue-20009"' "second multi-mission bundle claim can materialize"
multi_governor_run_b=$(printf '%s\n' "$multi_governor_claim_b" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
multi_governor_claim_c=$(run_in_repo bash scripts/harness-claim.sh 20010 "Harness Bundle Governor Peer 2" --owns="file:install-codex.sh" --allow-overlap)
assert_json_field "$multi_governor_claim_c" '"mission_id": "issue-20010"' "third multi-mission bundle claim can materialize"
multi_governor_run_c=$(printf '%s\n' "$multi_governor_claim_c" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)

multi_governor_update_a=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$multi_governor_run_a"
{
  "priority": "P0",
  "risk_level": "high"
}
EOF
)
assert_json_field "$multi_governor_update_a" '"priority": "P0"' "multi-mission bundle governor can raise its priority"
multi_governor_update_b=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$multi_governor_run_b"
{
  "priority": "P2",
  "risk_level": "low"
}
EOF
)
assert_json_field "$multi_governor_update_b" '"priority": "P2"' "multi-mission bundle peer can keep a lower authority tuple"
multi_governor_update_c=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$multi_governor_run_c"
{
  "priority": "P2",
  "risk_level": "low"
}
EOF
)
assert_json_field "$multi_governor_update_c" '"priority": "P2"' "third multi-mission bundle peer can keep a lower authority tuple"

graph_authority_multi_governor=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_multi_governor" '"resolution_kind": "bundle_governor_adjudicated"' "multi-mission bundle authority can adjudicate a bundle governor without a dependency edge"
assert_json_field "$graph_authority_multi_governor" '"bundle_state": "governed"' "multi-mission bundle authority marks the adjudicated bundle as governed"
assert_json_field "$graph_authority_multi_governor" "\"governing_mission_id\": \"issue-20008\"" "multi-mission bundle authority anchors governance to the strongest authority tuple"
assert_json_field "$graph_authority_multi_governor" '"authority_basis":' "multi-mission bundle authority exposes why the governor won"
assert_json_field "$graph_authority_multi_governor" '"losing_mission_ids":' "multi-mission bundle authority records losing missions for the adjudicated bundle"

resolved_multi_governor_a=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_governor_run_a" completed "Closed governed multi-mission bundle winner.")
assert_json_field "$resolved_multi_governor_a" '"state": "completed"' "governed multi-mission bundle winner can be closed"
resolved_multi_governor_b=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_governor_run_b" completed "Closed governed multi-mission bundle peer.")
assert_json_field "$resolved_multi_governor_b" '"state": "completed"' "governed multi-mission bundle peer can be closed"
resolved_multi_governor_c=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_governor_run_c" completed "Closed governed multi-mission bundle peer.")
assert_json_field "$resolved_multi_governor_c" '"state": "completed"' "third governed multi-mission bundle peer can be closed"

multi_escalation_claim_a=$(run_in_repo bash scripts/harness-claim.sh 20011 "Harness Bundle Escalation A" --owns="file:harness/README.md" --allow-overlap)
assert_json_field "$multi_escalation_claim_a" '"mission_id": "issue-20011"' "first escalation bundle claim can materialize"
multi_escalation_run_a=$(printf '%s\n' "$multi_escalation_claim_a" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
multi_escalation_claim_b=$(run_in_repo bash scripts/harness-claim.sh 20012 "Harness Bundle Escalation B" --owns="file:harness/README.md" --allow-overlap)
assert_json_field "$multi_escalation_claim_b" '"mission_id": "issue-20012"' "second escalation bundle claim can materialize"
multi_escalation_run_b=$(printf '%s\n' "$multi_escalation_claim_b" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
multi_escalation_claim_c=$(run_in_repo bash scripts/harness-claim.sh 20013 "Harness Bundle Escalation C" --owns="file:harness/README.md" --allow-overlap)
assert_json_field "$multi_escalation_claim_c" '"mission_id": "issue-20013"' "third escalation bundle claim can materialize"
multi_escalation_run_c=$(printf '%s\n' "$multi_escalation_claim_c" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)

for escalated_run in "$multi_escalation_run_a" "$multi_escalation_run_b" "$multi_escalation_run_c"; do
  escalated_update=$(
    cd "$TEST_REPO" &&
      cat <<EOF | bash scripts/harness-runs.sh update "$escalated_run"
{
  "priority": "P1",
  "risk_level": "medium"
}
EOF
  )
  assert_json_field "$escalated_update" '"priority": "P1"' "escalation bundle run can keep an equal authority tuple"
done

graph_authority_multi_escalation=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_multi_escalation" '"resolution_kind": "bundle_governor_required"' "multi-mission bundle authority escalates when no unique governor exists"
assert_json_field "$graph_authority_multi_escalation" '"bundle_state": "escalation-required"' "multi-mission bundle authority marks tied bundles as escalation-required"
assert_json_field "$graph_authority_multi_escalation" '"escalated_bundle_count": 1' "mission graph authority counts escalation-required bundles separately"
assert_json_field "$graph_authority_multi_escalation" '"required_action": "Open a council or add an explicit dependency/ownership transfer to establish a governing mission before parallel execution."' "escalation-required bundles get a leader-facing next action"
assert_json_field "$graph_authority_multi_escalation" '"candidate_governor_ids":' "escalation-required bundles expose candidate governing mission ids"
assert_json_field "$graph_authority_multi_escalation" '"candidate_governor_packets":' "escalation-required bundles expose candidate governing mission packets"
schedule_multi_escalation_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_multi_escalation_output" '"bundle_role": "escalation-required"' "scheduler marks tied bundles as escalation-required"
assert_json_field "$schedule_multi_escalation_output" '"decision_kind": "establish_bundle_governor"' "scheduler escalates tied bundles into an explicit governing-mission decision"
assert_json_field "$schedule_multi_escalation_output" '"decision_scope": "bundle"' "scheduler collapses tied bundle escalation into one bundle-scoped decision packet"
assert_json_field "$schedule_multi_escalation_output" '"mission_count": 3' "scheduler records how many missions participate in the bundle-scoped decision"
assert_json_field "$schedule_multi_escalation_output" '"candidate_governor_ids":' "scheduler keeps candidate governing mission ids on bundle decision packets"
assert_json_field "$schedule_multi_escalation_output" '"candidate_governor_packets":' "scheduler keeps candidate governing mission packets on bundle decision packets"
assert_json_field "$schedule_multi_escalation_output" '"recommended_decision_surface": "council"' "scheduler can recommend council handling when profile meeting triggers match the human decision"
assert_json_field "$schedule_multi_escalation_output" '"recommended_decision_protocol": "evidence-review"' "scheduler recommends evidence-review for conflicted bundle governance"
assert_json_field "$schedule_multi_escalation_output" '"meeting_trigger_refs":' "scheduler records which profile meeting triggers justify the council recommendation"
assert_json_field "$schedule_multi_escalation_output" '"council_decision_count": 1' "scheduler counts one bundle-scoped council decision instead of one per mission"
assert_json_field "$schedule_multi_escalation_output" '"next_action": "Open a council or add an explicit dependency/ownership transfer to establish a governing mission before parallel execution."' "scheduler prioritizes council-backed bundle-governor decisions over foreground selection messaging"
assert_json_field "$schedule_multi_escalation_output" '"bundle_role_tally":' "scheduler reports bundle-role tallies in the portfolio snapshot"
assert_json_field "$schedule_multi_escalation_output" '"required_human_decision_scope_tally":' "scheduler reports human-decision scope tallies"
assert_json_field "$schedule_multi_escalation_output" '"required_human_decision_surface_tally":' "scheduler reports human-decision surface tallies"
assert_json_field "$schedule_multi_escalation_output" '"primary_required_human_decision": {' "scheduler exposes the primary decision packet"
assert_json_field "$schedule_multi_escalation_output" '"decision_scope": "bundle"' "scheduler marks the primary escalation packet as bundle-scoped"
assert_json_field "$schedule_multi_escalation_output" '"human_decision_queue_position": 1' "scheduler annotates primary escalation missions with queue position 1"
assert_json_field "$schedule_multi_escalation_output" '"primary_human_decision": true' "scheduler marks primary escalation missions explicitly"
assert_json_field "$schedule_multi_escalation_output" '"candidate_governor_packets":' "scheduler also exposes candidate governors on mission classifications within escalation bundles"
home_output_multi_escalation=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output_multi_escalation" '"mission_escalated_bundle_count": 1' "home surfaces escalation-required mission bundles"
assert_json_field "$home_output_multi_escalation" '"scheduler_council_decision_count": 1' "home surfaces bundle-scoped council decisions instead of per-mission duplicates"
assert_json_field "$home_output_multi_escalation" '"next_action": "Open a council or add an explicit dependency/ownership transfer to establish a governing mission before parallel execution."' "home keeps scheduler decision-first guidance when bundle governance should become a council"
assert_json_field "$home_output_multi_escalation" '"scheduler_bundle_role_tally":' "home surfaces scheduler bundle-role tallies"
assert_json_field "$home_output_multi_escalation" '"scheduler_primary_required_human_decision": {' "home surfaces the scheduler's primary decision packet"
multi_escalation_materialized=$(run_in_repo bash scripts/harness-meetings.sh materialize-primary-decision --mission=issue-20011)
assert_json_field "$multi_escalation_materialized" '"source_decision_scope": "bundle"' "bundle-scoped scheduler decision can materialize into a governing council"
assert_json_field "$multi_escalation_materialized" '"source_decision_kind": "establish_bundle_governor"' "materialized bundle council keeps the scheduler decision kind"
assert_json_field "$multi_escalation_materialized" "\"run_id\": \"$multi_escalation_run_b\"" "materialized bundle council auto-seeds one candidate-governor participant"
assert_json_field "$multi_escalation_materialized" "\"run_id\": \"$multi_escalation_run_c\"" "materialized bundle council auto-seeds another candidate-governor participant"
assert_json_field "$multi_escalation_materialized" '/.harness/missions/issue-20012.yml"' "materialized bundle council includes sibling mission context refs"
assert_json_field "$multi_escalation_materialized" '"candidate_governor_packets":' "materialized bundle council carries candidate governor packets into the meeting packet"
assert_json_field "$multi_escalation_materialized" '"Candidate governor issue-20012:' "materialized bundle council summarizes candidate-governor evidence in known context"
multi_escalation_meeting_id=$(printf '%s\n' "$multi_escalation_materialized" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
multi_escalation_decision_id=$(printf '%s\n' "$multi_escalation_materialized" | sed -n 's/^[[:space:]]*"source_decision_id": "\(decision-[^"]*\)",$/\1/p' | head -n1)
multi_escalation_bundle_id="bundle-issue-20011-issue-20012-issue-20013"
schedule_multi_escalation_active=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_multi_escalation_active" '"active_council_bundle_count": 1' "scheduler counts escalation bundles already governed by an active council"
assert_json_field "$schedule_multi_escalation_active" '"recommended_decision_surface": "active-council"' "scheduler upgrades bundle-governor decisions into active-council once the council exists"
assert_json_field "$schedule_multi_escalation_active" "\"$multi_escalation_meeting_id\"" "scheduler ties the escalation bundle decision back to the active meeting id"
assert_json_field "$schedule_multi_escalation_active" "\"decision_id\": \"$multi_escalation_decision_id\"" "scheduler keeps the bundle-scoped active-council decision in the portfolio after materialization"
assert_json_field "$schedule_multi_escalation_active" '"decision_scope": "bundle"' "scheduler keeps bundle-scoped authority packets visible while the bundle council stays active"
assert_json_field "$schedule_multi_escalation_active" '"bundle_active_decision_meeting_count": 1' "mission classifications inherit bundle-level active council coverage"
assert_json_field "$schedule_multi_escalation_active" '"bundle_control_state": "active-council"' "mission classifications expose bundle control state once the bundle council is active"
assert_json_field "$schedule_multi_escalation_active" '"action_kind": "hold_for_active_bundle_council"' "scheduler emits reschedule actions that hold execution behind an active bundle council"
multi_escalation_reused=$(run_in_repo bash scripts/harness-meetings.sh materialize-decision "$multi_escalation_decision_id" --mission=issue-20011)
assert_json_field "$multi_escalation_reused" "\"meeting_id\": \"$multi_escalation_meeting_id\"" "re-materializing an active bundle decision reuses the same meeting"
assert_json_field "$multi_escalation_reused" '"reused_existing_meeting": true' "reused bundle decision materialization reports active-meeting reuse"
assert_json_field "$multi_escalation_reused" '"candidate_governor_packets":' "reused bundle council keeps bundle candidate packets on the refreshed meeting packet"
assert_json_field "$multi_escalation_reused" '"Candidate governor issue-20012:' "reused bundle council keeps candidate-governor context when re-materialized"
home_output_multi_escalation_active=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output_multi_escalation_active" '"scheduler_active_council_bundle_count": 1' "home surfaces active bundle-council coverage from the scheduler"
multi_escalation_conflicting_governor=$(run_in_repo bash scripts/harness-missions.sh establish-governor "$multi_escalation_bundle_id" issue-20012 "Conflicting explicit governor during active bundle council coverage.")
assert_json_field "$multi_escalation_conflicting_governor" "\"bundle_id\": \"$multi_escalation_bundle_id\"" "leader can emit an explicit governor that conflicts with a live bundle council host"
schedule_multi_escalation_conflict=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_multi_escalation_conflict" '"action_kind": "clear_invalid_bundle_governor"' "scheduler escalates explicit-governor versus active-council conflicts into maintenance"
assert_json_field "$schedule_multi_escalation_conflict" '"Explicit bundle governor points at issue-20012' "scheduler explains the conflicting explicit governor against the active bundle council anchor"
cleared_multi_escalation_governor=$(run_in_repo bash scripts/harness-missions.sh clear-governor "$multi_escalation_bundle_id" "Clear conflicting explicit governor after active bundle council coverage.")
assert_json_field "$cleared_multi_escalation_governor" "\"bundle_id\": \"$multi_escalation_bundle_id\"" "conflicting explicit governor can be cleared after active bundle-council coverage"
run_in_repo bash scripts/harness-meetings.sh orchestrate "$multi_escalation_meeting_id" >/dev/null
schedule_multi_escalation_orchestrated=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_multi_escalation_orchestrated" '"active_decision_phase": "independent_review"' "scheduler projects active bundle-council phase into the portfolio packet"
assert_json_field "$schedule_multi_escalation_orchestrated" '"active_decision_next_requirement": "Collect' "scheduler projects active bundle-council next requirement into the portfolio packet"
assert_json_field "$schedule_multi_escalation_orchestrated" '"next_action": "Collect' "scheduler next action follows the live active bundle-council requirement"
assert_json_field "$schedule_multi_escalation_orchestrated" '"action_kind": "hold_for_active_bundle_council"' "scheduler still emits active bundle-council hold actions after progress projection"
assert_json_field "$schedule_multi_escalation_orchestrated" '"reason": "Collect' "active bundle-council hold actions inherit the live next requirement"
home_output_multi_escalation_orchestrated=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output_multi_escalation_orchestrated" '"next_action": "Collect' "home mirrors the live active bundle-council requirement from scheduler state"
set_json_field_in_file "$TEST_REPO/.harness/meetings/active/$multi_escalation_meeting_id.yml" "meeting.updated_at" "2026-03-01T00:00:00Z"
schedule_multi_escalation_stale=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_multi_escalation_stale" '"active_decision_stale": true' "scheduler flags stale active bundle councils in the portfolio packet"
assert_json_field "$schedule_multi_escalation_stale" '"next_action": "Refresh or close the stale meeting before trusting it as current governance state."' "scheduler prioritizes stale active bundle councils as the current next action"
assert_json_field "$schedule_multi_escalation_stale" '"reason": "Refresh or close the stale active council(s) for this decision before trusting them as current governance state:' "active bundle-council hold actions inherit the stale-council requirement"
home_output_multi_escalation_stale=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output_multi_escalation_stale" '"next_action": "Refresh or close the stale meeting before trusting it as current governance state."' "home mirrors stale active bundle-council guidance from scheduler state"
run_in_repo bash scripts/harness-meetings.sh add-position "$multi_escalation_meeting_id" "$multi_escalation_run_a" independent_review support "Bundle governance should anchor to issue-20011 because it has the strongest release ownership evidence." --claim=bundle-governor-choice --dimension=factual --authority=L2 --evidence=harness/README.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$multi_escalation_meeting_id" "$multi_escalation_run_b" independent_review support "Independent review agrees issue-20011 should govern the bundle rather than issue-20012 or issue-20013." --claim=bundle-governor-choice --dimension=factual --authority=L3 --evidence=docs/harness-public-contract.md >/dev/null
auto_multi_escalation_governor=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$multi_escalation_meeting_id")
assert_json_field "$auto_multi_escalation_governor" '"generated_locked_decision_count": 1' "bundle council auto-synthesis can generate one structured governing decision"
assert_json_field "$auto_multi_escalation_governor" '"governing_mission_id": "issue-20011"' "bundle council auto-synthesis locks the strongest governing mission into the council record"
closed_multi_escalation_meeting=$(run_in_repo bash scripts/harness-drive.sh meeting-close "$multi_escalation_run_a" "$multi_escalation_meeting_id" "Bundle escalation council materialization covered.")
assert_json_field "$closed_multi_escalation_meeting" '"state": "closed"' "bundle-scoped scheduler-materialized council can close cleanly"
assert_json_field "$closed_multi_escalation_meeting" '"adoption_kind": "bundle_governor_established"' "closing a structured bundle council auto-adopts bundle-governor authority"
assert_json_field "$closed_multi_escalation_meeting" '"governing_mission_id": "issue-20011"' "bundle council close records the adopted governing mission"
schedule_multi_escalation_adopted=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_multi_escalation_adopted" '"bundle_state": "governed"' "closed bundle council adoption turns the escalated bundle into governed state"
assert_json_field "$schedule_multi_escalation_adopted" "\"bundle_governing_mission_id\": \"issue-20011\"" "scheduler consumes the auto-adopted bundle governor after council close"
assert_json_field "$schedule_multi_escalation_adopted" '"bundle_governance_source_kind": "council-adopted-governor"' "scheduler preserves that the bundle governor came from a closed council adoption"
assert_json_field "$schedule_multi_escalation_adopted" "\"bundle_adopted_from_meeting_id\": \"$multi_escalation_meeting_id\"" "scheduler preserves the adopted council meeting id on governed bundle state"
cleared_multi_escalation_adopted_governor=$(run_in_repo bash scripts/harness-missions.sh clear-governor "$multi_escalation_bundle_id" "Clear adopted explicit governor after bundle-council adoption coverage.")
assert_json_field "$cleared_multi_escalation_adopted_governor" "\"bundle_id\": \"$multi_escalation_bundle_id\"" "auto-adopted bundle governor can be cleared after coverage to keep later smoke isolated"

multi_tuple_claim_a=$(run_in_repo bash scripts/harness-claim.sh 20126 "Harness Bundle Tuple A" --owns="file:docs/harness-session-handoff.md" --allow-overlap)
assert_json_field "$multi_tuple_claim_a" '"mission_id": "issue-20126"' "first tuple-governed bundle claim can materialize"
multi_tuple_run_a=$(printf '%s\n' "$multi_tuple_claim_a" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
multi_tuple_claim_b=$(run_in_repo bash scripts/harness-claim.sh 20127 "Harness Bundle Tuple B" --owns="file:docs/harness-session-handoff.md" --allow-overlap)
assert_json_field "$multi_tuple_claim_b" '"mission_id": "issue-20127"' "second tuple-governed bundle claim can materialize"
multi_tuple_run_b=$(printf '%s\n' "$multi_tuple_claim_b" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
multi_tuple_claim_c=$(run_in_repo bash scripts/harness-claim.sh 20128 "Harness Bundle Tuple C" --owns="file:docs/harness-session-handoff.md" --allow-overlap)
assert_json_field "$multi_tuple_claim_c" '"mission_id": "issue-20128"' "third tuple-governed bundle claim can materialize"
multi_tuple_run_c=$(printf '%s\n' "$multi_tuple_claim_c" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)

for tuple_run in "$multi_tuple_run_a" "$multi_tuple_run_b" "$multi_tuple_run_c"; do
  tuple_update=$(
    cd "$TEST_REPO" &&
      cat <<EOF | bash scripts/harness-runs.sh update "$tuple_run"
{
  "priority": "P1",
  "risk_level": "medium"
}
EOF
  )
  assert_json_field "$tuple_update" '"priority": "P1"' "tuple-governed bundle run can keep an equal priority baseline"
done

multi_tuple_release_update=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$multi_tuple_run_b"
{
  "ship_status": "pushed"
}
EOF
)
assert_json_field "$multi_tuple_release_update" '"ship_status": "pushed"' "tuple-governed bundle can add a unique release-adjacent candidate"

graph_authority_multi_tuple=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_multi_tuple" '"resolution_kind": "bundle_governor_adjudicated"' "multi-mission bundle authority can adjudicate from the full candidate-governor tuple"
assert_json_field "$graph_authority_multi_tuple" "\"governing_mission_id\": \"issue-20127\"" "multi-mission bundle authority anchors governance to the strongest tuple candidate"
assert_json_field "$graph_authority_multi_tuple" '"authority_basis":' "tuple-governed bundle exposes the authority basis behind the winner"
assert_json_field "$graph_authority_multi_tuple" '"release_adjacent"' "tuple-governed bundle records the release-adjacent basis when it breaks an otherwise tied bundle"
schedule_multi_tuple_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_multi_tuple_output" "\"bundle_governing_mission_id\": \"issue-20127\"" "scheduler consumes the adjudicated tuple-governed bundle"
assert_json_field "$schedule_multi_tuple_output" '"bundle_state": "governed"' "scheduler marks the tuple-governed bundle as governed instead of escalation-required"

resolved_multi_tuple_a=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_tuple_run_a" completed "Closed tuple-governed bundle peer.")
assert_json_field "$resolved_multi_tuple_a" '"state": "completed"' "tuple-governed bundle peer can be closed"
resolved_multi_tuple_b=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_tuple_run_b" completed "Closed tuple-governed bundle winner.")
assert_json_field "$resolved_multi_tuple_b" '"state": "completed"' "tuple-governed bundle winner can be closed"
resolved_multi_tuple_c=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_tuple_run_c" completed "Closed tuple-governed bundle peer.")
assert_json_field "$resolved_multi_tuple_c" '"state": "completed"' "third tuple-governed bundle peer can be closed"

decision_portfolio_claim=$(run_in_repo bash scripts/harness-claim.sh 20023 "Harness Decision Portfolio" --owns="file:docs/README.md")
assert_json_field "$decision_portfolio_claim" '"mission_id": "issue-20023"' "independent human-decision mission can materialize"
decision_portfolio_run=$(printf '%s\n' "$decision_portfolio_claim" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
decision_portfolio_signal=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-decision-portfolio",
  "kind": "requires_meeting",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$decision_portfolio_run",
  "mission_id": "issue-20023",
  "applies_to": ["issue-20023"],
  "scope_type": "mission",
  "title": "Secondary decision packet",
  "summary": "An isolated mission now also requires a governed human decision.",
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-29T00:00:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$decision_portfolio_signal" '"signal_id": "sig-smoke-decision-portfolio"' "secondary mission-level decision signal can be emitted directly"
schedule_decision_portfolio_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_decision_portfolio_output" '"required_human_decision_count": 2' "scheduler counts bundle and mission decisions together in the portfolio"
assert_json_field "$schedule_decision_portfolio_output" '"deferred_human_decision_count": 1' "scheduler records how many additional decisions remain behind the primary packet"
assert_json_field "$schedule_decision_portfolio_output" '"required_human_decision_scope_tally": {' "scheduler exposes decision-scope portfolio tallies"
assert_json_field "$schedule_decision_portfolio_output" '"bundle": 1' "scheduler counts one bundle-scoped decision in the portfolio"
assert_json_field "$schedule_decision_portfolio_output" '"mission": 1' "scheduler counts one mission-scoped decision in the portfolio"
assert_json_field "$schedule_decision_portfolio_output" '"required_human_decision_surface_tally": {' "scheduler exposes decision-surface portfolio tallies"
assert_json_field "$schedule_decision_portfolio_output" '"council": 2' "scheduler counts both portfolio decisions as council-surface under the current profile triggers"
assert_json_field "$schedule_decision_portfolio_output" '"primary_required_human_decision": {' "scheduler keeps a primary decision packet for the portfolio"
assert_json_field "$schedule_decision_portfolio_output" '"recommended_decision_protocol": "evidence-review"' "scheduler carries decision protocol guidance in the human-decision portfolio"
assert_json_field "$schedule_decision_portfolio_output" '"decision_id": "decision-' "scheduler now gives human-decision packets stable decision ids"
assert_json_field "$schedule_decision_portfolio_output" '"bundle_id": "bundle-issue-20011-issue-20012-issue-20013"' "scheduler prioritizes the escalation bundle as the primary decision packet"
assert_json_field "$schedule_decision_portfolio_output" '"mission_count": 3' "primary bundle decision packet keeps the grouped mission count"
assert_json_field "$schedule_decision_portfolio_output" "\"mission_id\": \"issue-20023\"" "secondary mission decision remains visible in the decision portfolio"
assert_json_field "$schedule_decision_portfolio_output" '"human_decision_queue_position": 2' "scheduler annotates the lower-priority mission decision with queue position 2"
assert_json_field "$schedule_decision_portfolio_output" '"primary_human_decision": false' "scheduler marks non-primary human decisions explicitly"
assert_json_field "$schedule_decision_portfolio_output" '"human_decision_id": "decision-' "scheduler projects stable decision ids back into mission classifications"
assert_json_field "$schedule_decision_portfolio_output" '"action_kind": "defer_pending_primary_human_decision"' "scheduler emits an explicit defer action for lower-priority human decisions"
assert_json_field "$schedule_decision_portfolio_output" '"human_decision_position": 2' "defer action records the lower-priority decision queue position"
assert_json_field "$schedule_decision_portfolio_output" '"decision_id": "decision-' "defer action carries the stable decision id for the blocked packet"
resolved_decision_portfolio_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-decision-portfolio resolved "Secondary decision packet covered.")
assert_json_field "$resolved_decision_portfolio_signal" '"state": "resolved"' "secondary mission-level decision signal can be resolved after portfolio coverage"
resolved_decision_portfolio_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$decision_portfolio_run" completed "Closed decision-portfolio mission.")
assert_json_field "$resolved_decision_portfolio_run" '"state": "completed"' "secondary mission-level decision can be closed after portfolio coverage"
resolved_multi_escalation_a=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_escalation_run_a" completed "Closed escalated multi-mission bundle member.")
assert_json_field "$resolved_multi_escalation_a" '"state": "completed"' "first escalation bundle member can be closed"
resolved_multi_escalation_b=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_escalation_run_b" completed "Closed escalated multi-mission bundle member.")
assert_json_field "$resolved_multi_escalation_b" '"state": "completed"' "second escalation bundle member can be closed"
resolved_multi_escalation_c=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_escalation_run_c" completed "Closed escalated multi-mission bundle member.")
assert_json_field "$resolved_multi_escalation_c" '"state": "completed"' "third escalation bundle member can be closed after bundle-governor conflict coverage"

decision_hold_claim=$(run_in_repo bash scripts/harness-claim.sh 20024 "Harness Decision Hold" --owns="file:docs/harness-public-contract.md")
assert_json_field "$decision_hold_claim" '"mission_id": "issue-20024"' "primary decision-hold mission can materialize"
decision_hold_run=$(printf '%s\n' "$decision_hold_claim" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
decision_hold_seed_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$decision_hold_seed_schedule" '"selected_mission_count": 1' "decision-hold seed mission can establish one foreground mission first"
assert_json_field "$decision_hold_seed_schedule" '"mission_id": "issue-20024"' "decision-hold seed mission remains visible in the scheduler output"
decision_hold_signal=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-decision-hold",
  "kind": "requires_meeting",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$decision_hold_run",
  "mission_id": "issue-20024",
  "applies_to": ["issue-20024"],
  "scope_type": "mission",
  "title": "Primary decision packet",
  "summary": "Foreground mission now requires a governed human decision before new work should start.",
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-30T00:00:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$decision_hold_signal" '"signal_id": "sig-smoke-decision-hold"' "primary decision-hold signal can be emitted directly"
decision_hold_peer_claim=$(run_in_repo bash scripts/harness-claim.sh 20025 "Harness Decision Hold Peer" --owns="file:docs/harness-documentation-map.md")
assert_json_field "$decision_hold_peer_claim" '"mission_id": "issue-20025"' "decision-hold peer mission can materialize"
decision_hold_peer_run=$(printf '%s\n' "$decision_hold_peer_claim" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
baseline_wait=$(run_in_repo bash scripts/harness-runs.sh set-state "$run_id" waiting "Pause the old foreground while decision-hold admission is tested.")
assert_json_field "$baseline_wait" '"state": "waiting"' "old foreground can be paused for decision-hold coverage"
decision_hold_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$decision_hold_schedule" '"required_human_decision_count": 1' "scheduler still sees the outstanding governed decision"
assert_json_field "$decision_hold_schedule" '"selected_mission_count": 0' "scheduler withholds new foreground admission while the primary human decision is unresolved"
assert_json_field "$decision_hold_schedule" '"batch_admission_policy": "hold_for_primary_human_decision"' "scheduler records decision-first hold as the admission policy"
assert_json_field "$decision_hold_schedule" '"action_kind": "hold_for_primary_human_decision"' "scheduler emits an explicit hold action for newly blocked foreground work"
assert_json_field "$decision_hold_schedule" '"mission_id": "issue-20025"' "scheduler still keeps the deferred peer visible in the mission classifications"
assert_json_field "$decision_hold_schedule" '"classification": "deferred"' "scheduler defers the new foreground peer while the primary decision stays open"
resolved_decision_hold_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-decision-hold resolved "Primary decision packet covered.")
assert_json_field "$resolved_decision_hold_signal" '"state": "resolved"' "primary decision-hold signal can be resolved after scheduler coverage"
resolved_decision_hold_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$decision_hold_run" completed "Closed decision-hold mission.")
assert_json_field "$resolved_decision_hold_run" '"state": "completed"' "primary decision-hold mission can be closed after scheduler coverage"
resolved_decision_hold_peer_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$decision_hold_peer_run" completed "Closed decision-hold peer mission.")
assert_json_field "$resolved_decision_hold_peer_run" '"state": "completed"' "decision-hold peer mission can be closed after scheduler coverage"
baseline_resume=$(run_in_repo bash scripts/harness-runs.sh set-state "$run_id" ready "Restore the baseline foreground after decision-hold coverage.")
assert_json_field "$baseline_resume" '"state": "ready"' "old foreground can be restored after decision-hold coverage"

decision_hold_restore_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$decision_hold_restore_schedule" '"selected_mission_count": 1' "baseline foreground can be re-established after escalation backlog clears"
assert_json_field "$decision_hold_restore_schedule" '"mission_id": "issue-20001"' "restored baseline foreground remains visible after escalation backlog clears"

section "Explicit bundle governor adoption"
explicit_governor_claim_a=$(run_in_repo bash scripts/harness-claim.sh 20028 "Harness Explicit Governor A" --owns="file:docs/harness-operating-model.md" --allow-overlap)
assert_json_field "$explicit_governor_claim_a" '"mission_id": "issue-20028"' "explicit-governor bundle mission A can materialize"
explicit_governor_run_a=$(printf '%s\n' "$explicit_governor_claim_a" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
explicit_governor_claim_b=$(run_in_repo bash scripts/harness-claim.sh 20029 "Harness Explicit Governor B" --owns="file:docs/harness-operating-model.md" --allow-overlap)
assert_json_field "$explicit_governor_claim_b" '"mission_id": "issue-20029"' "explicit-governor bundle mission B can materialize"
explicit_governor_run_b=$(printf '%s\n' "$explicit_governor_claim_b" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
explicit_governor_claim_c=$(run_in_repo bash scripts/harness-claim.sh 20030 "Harness Explicit Governor C" --owns="file:docs/harness-operating-model.md" --allow-overlap)
assert_json_field "$explicit_governor_claim_c" '"mission_id": "issue-20030"' "explicit-governor bundle mission C can materialize"
explicit_governor_run_c=$(printf '%s\n' "$explicit_governor_claim_c" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)

for explicit_run in "$explicit_governor_run_a" "$explicit_governor_run_b" "$explicit_governor_run_c"; do
  explicit_update=$(
    cd "$TEST_REPO" &&
      cat <<EOF | bash scripts/harness-runs.sh update "$explicit_run"
{
  "priority": "P1",
  "risk_level": "medium"
}
EOF
  )
  assert_json_field "$explicit_update" '"priority": "P1"' "explicit-governor bundle run can keep an equal authority tuple"
done

explicit_governor_bundle_id="bundle-issue-20028-issue-20029-issue-20030"
graph_authority_explicit_before=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_explicit_before" "\"bundle_id\": \"$explicit_governor_bundle_id\"" "explicit-governor bundle is visible before manual adoption"
assert_json_field "$graph_authority_explicit_before" '"resolution_kind": "bundle_governor_required"' "explicit-governor bundle starts as escalation-required"

establish_governor_output=$(run_in_repo bash scripts/harness-missions.sh establish-governor "$explicit_governor_bundle_id" issue-20028 "Leader establishes an explicit bundle governor for the documentation boundary.")
assert_json_field "$establish_governor_output" '"kind": "bundle_governor_established"' "missions establish-governor emits an explicit governor signal"
assert_json_field "$establish_governor_output" "\"governing_mission_id\": \"issue-20028\"" "explicit governor signal records the chosen governing mission"

graph_authority_explicit_after=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_explicit_after" '"resolution_kind": "explicit_bundle_governor"' "explicit governor signal upgrades the bundle into governed authority"
assert_json_field "$graph_authority_explicit_after" "\"governing_mission_id\": \"issue-20028\"" "explicit governor signal anchors bundle authority to the chosen mission"
assert_json_field "$graph_authority_explicit_after" '"explicit_governor_signal_id": "sig-' "explicit governor signal id is visible in the bundle packet"

schedule_explicit_governor=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_explicit_governor" "\"bundle_governing_mission_id\": \"issue-20028\"" "scheduler follows the explicit bundle governor once it is established"
assert_json_field "$schedule_explicit_governor" '"bundle_role": "governed"' "scheduler treats explicit-governor peers as governed instead of escalation-required"
if grep -q "\"bundle_id\": \"$explicit_governor_bundle_id\"" <<<"$schedule_explicit_governor" && grep -q '"decision_kind": "establish_bundle_governor"' <<<"$schedule_explicit_governor"; then
  fail "explicit-governor bundle no longer emits establish-governor decision packets"
  echo "$schedule_explicit_governor"
else
  pass "explicit-governor bundle no longer emits establish-governor decision packets"
fi

resolved_explicit_governor_a=$(run_in_repo bash scripts/harness-runs.sh set-state "$explicit_governor_run_a" completed "Closed explicit governor leader mission to test stale-governor maintenance.")
assert_json_field "$resolved_explicit_governor_a" '"state": "completed"' "explicit-governor leader mission can be closed to simulate stale bundle governance"

schedule_explicit_governor_stale=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_explicit_governor_stale" '"action_kind": "clear_invalid_bundle_governor"' "scheduler escalates stale explicit-governor signals into maintenance"
assert_json_field "$schedule_explicit_governor_stale" "\"bundle_id\": \"$explicit_governor_bundle_id\"" "scheduler maintenance identifies the stale explicit-governor bundle"
assert_json_field "$schedule_explicit_governor_stale" '"primary_maintenance_action": {' "scheduler exposes the primary maintenance action"
assert_json_field "$schedule_explicit_governor_stale" '"maintenance_selection_policy": "maintenance_priority"' "scheduler starts stale-governor maintenance with the default maintenance-selection policy"
assert_json_field "$schedule_explicit_governor_stale" '"batch_admission_policy": "hold_for_primary_maintenance"' "critical bundle-governor maintenance blocks new low-urgency foreground work"
assert_json_field "$schedule_explicit_governor_stale" '"action_kind": "hold_for_primary_maintenance"' "scheduler emits explicit hold actions when critical maintenance outranks foreground work"

home_output_explicit_governor_stale=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output_explicit_governor_stale" '"scheduler_primary_maintenance_action": {' "home mirrors the scheduler's primary maintenance action"
assert_json_field "$home_output_explicit_governor_stale" '"scheduler_maintenance_selection_policy": "maintenance_priority"' "home mirrors maintenance-selection policy"

doctor_explicit_governor_stale=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_explicit_governor_stale" '"issue_count": 2' "doctor detects stale explicit-governor drift alongside the restored overlap conflict"
assert_json_field "$doctor_explicit_governor_stale" '"invalid-explicit-bundle-governor"' "doctor names the stale explicit-governor issue explicitly"

maintenance_vs_decision_claim=$(run_in_repo bash scripts/harness-claim.sh 20031 "Harness Maintenance Versus Decision" --owns="file:docs/harness-session-handoff.md")
assert_json_field "$maintenance_vs_decision_claim" '"mission_id": "issue-20031"' "maintenance-versus-decision mission can materialize"
maintenance_vs_decision_run=$(printf '%s\n' "$maintenance_vs_decision_claim" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
maintenance_vs_decision_signal=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-maintenance-vs-decision",
  "kind": "requires_meeting",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$maintenance_vs_decision_run",
  "mission_id": "issue-20031",
  "applies_to": ["issue-20031"],
  "scope_type": "mission",
  "title": "Lower-priority human decision",
  "summary": "A separate mission also wants governed human attention while stale bundle maintenance is still open.",
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-30T02:00:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$maintenance_vs_decision_signal" '"signal_id": "sig-smoke-maintenance-vs-decision"' "secondary human decision can be emitted while maintenance is still open"

schedule_explicit_governor_control_plane=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_explicit_governor_control_plane" '"primary_control_plane_kind": "maintenance"' "critical maintenance outranks lower-priority human decisions in the unified control-plane portfolio"
assert_json_field "$schedule_explicit_governor_control_plane" '"primary_control_plane_selection_policy": "critical_maintenance_over_decision"' "scheduler records when maintenance wins over human decisions"
assert_json_field "$schedule_explicit_governor_control_plane" '"primary_control_plane_packet": {' "scheduler exposes the unified primary control-plane packet"
assert_json_field "$schedule_explicit_governor_control_plane" '"action_kind": "clear_invalid_bundle_governor"' "stale explicit governor remains the primary control-plane packet"
assert_json_field "$schedule_explicit_governor_control_plane" '"required_human_decision_count": 2' "human decisions remain visible even when critical maintenance wins the unified portfolio"
assert_json_field "$schedule_explicit_governor_control_plane" '"deferred_by_primary_control_plane_kind": "maintenance"' "mission classifications record when governed human decisions are deferred behind maintenance"
assert_json_field "$schedule_explicit_governor_control_plane" '"action_kind": "defer_pending_primary_maintenance"' "scheduler emits explicit defer actions when primary maintenance outranks governed human decisions"
assert_json_field "$schedule_explicit_governor_control_plane" '"next_action": "Clear the stale explicit governor signal so the control plane stops carrying orphaned governance."' "next action follows the unified primary control-plane packet"

home_output_explicit_governor_control_plane=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output_explicit_governor_control_plane" '"scheduler_primary_control_plane_kind": "maintenance"' "home mirrors the scheduler's unified primary control-plane kind"
assert_json_field "$home_output_explicit_governor_control_plane" '"scheduler_primary_control_plane_selection_policy": "critical_maintenance_over_decision"' "home mirrors the unified control-plane selection policy"
assert_json_field "$home_output_explicit_governor_control_plane" '"scheduler_primary_control_plane_packet": {' "home mirrors the unified primary control-plane packet"

clear_governor_output=$(run_in_repo bash scripts/harness-missions.sh clear-governor "$explicit_governor_bundle_id" "Explicit governor coverage completed.")
assert_json_field "$clear_governor_output" "\"bundle_id\": \"$explicit_governor_bundle_id\"" "missions clear-governor can retire explicit governor signals"
assert_json_field "$clear_governor_output" '"resolved_signal_ids":' "missions clear-governor reports the resolved signal ids"

resolved_maintenance_vs_decision_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-maintenance-vs-decision resolved "Maintenance-vs-decision coverage completed.")
assert_json_field "$resolved_maintenance_vs_decision_signal" '"state": "resolved"' "secondary human-decision signal can be resolved after unified control-plane coverage"
resolved_maintenance_vs_decision_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$maintenance_vs_decision_run" completed "Closed maintenance-vs-decision mission.")
assert_json_field "$resolved_maintenance_vs_decision_run" '"state": "completed"' "maintenance-versus-decision mission can be closed after unified control-plane coverage"

graph_authority_explicit_cleared=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_explicit_cleared" '"bundle_id": "bundle-issue-20029-issue-20030"' "clearing a stale explicit governor returns authority to the surviving overlap bundle"
assert_json_field "$graph_authority_explicit_cleared" '"resolution_kind": "bundle_governor_required"' "clearing the explicit governor returns the bundle to escalation-required"

resolved_explicit_governor_b=$(run_in_repo bash scripts/harness-runs.sh set-state "$explicit_governor_run_b" completed "Closed explicit-governor bundle mission B.")
assert_json_field "$resolved_explicit_governor_b" '"state": "completed"' "explicit-governor bundle mission B can be closed"
resolved_explicit_governor_c=$(run_in_repo bash scripts/harness-runs.sh set-state "$explicit_governor_run_c" completed "Closed explicit-governor bundle mission C.")
assert_json_field "$resolved_explicit_governor_c" '"state": "completed"' "explicit-governor bundle mission C can be closed"
baseline_after_explicit_governor=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$baseline_after_explicit_governor" "\"mission_id\": \"issue-20001\"" "baseline foreground can be re-established after stale explicit governor coverage"
assert_json_field "$baseline_after_explicit_governor" '"selected_mission_count": 1' "baseline foreground re-establishes a single selected mission after stale explicit governor coverage"

section "Sticky decision portfolio"
sticky_decision_claim_a=$(run_in_repo bash scripts/harness-claim.sh 20026 "Harness Sticky Decision A" --owns="file:docs/README.md")
assert_json_field "$sticky_decision_claim_a" '"mission_id": "issue-20026"' "sticky decision mission A can materialize"
sticky_decision_run_a=$(printf '%s\n' "$sticky_decision_claim_a" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
sticky_decision_signal_a=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-sticky-decision-a",
  "kind": "requires_meeting",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$sticky_decision_run_a",
  "mission_id": "issue-20026",
  "applies_to": ["issue-20026"],
  "scope_type": "mission",
  "title": "Sticky decision A",
  "summary": "First mission-scoped decision packet in the portfolio.",
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-29T16:00:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$sticky_decision_signal_a" '"signal_id": "sig-smoke-sticky-decision-a"' "sticky decision signal A can be emitted directly"
sticky_decision_first_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$sticky_decision_first_schedule" '"decision_selection_policy": "portfolio_priority_risk"' "single decision portfolio starts with the default decision-selection policy"
assert_json_field "$sticky_decision_first_schedule" '"required_human_decision_count": 1' "single decision portfolio starts with exactly one human decision"
assert_json_field "$sticky_decision_first_schedule" "\"mission_id\": \"issue-20026\"" "single decision portfolio keeps mission A as the primary packet"
sticky_decision_claim_b=$(run_in_repo bash scripts/harness-claim.sh 20027 "Harness Sticky Decision B" --owns="file:docs/harness-public-contract.md")
assert_json_field "$sticky_decision_claim_b" '"mission_id": "issue-20027"' "sticky decision mission B can materialize"
sticky_decision_run_b=$(printf '%s\n' "$sticky_decision_claim_b" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
sticky_decision_signal_b=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-sticky-decision-b",
  "kind": "requires_meeting",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$sticky_decision_run_b",
  "mission_id": "issue-20027",
  "applies_to": ["issue-20027"],
  "scope_type": "mission",
  "title": "Sticky decision B",
  "summary": "Second mission-scoped decision packet with the same decision class.",
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-29T16:05:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$sticky_decision_signal_b" '"signal_id": "sig-smoke-sticky-decision-b"' "sticky decision signal B can be emitted directly"
sticky_decision_second_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$sticky_decision_second_schedule" '"decision_selection_policy": "sticky_primary_human_decision"' "decision portfolio retains the previous primary packet when a same-class peer arrives"
assert_json_field "$sticky_decision_second_schedule" '"previous_primary_human_decision_key": "mission:issue-20026"' "decision portfolio records the previous primary decision key"
assert_json_field "$sticky_decision_second_schedule" '"required_human_decision_count": 2' "decision portfolio sees both same-class mission decisions"
assert_json_field "$sticky_decision_second_schedule" "\"mission_id\": \"issue-20026\"" "decision portfolio keeps mission A as the primary packet under a same-class tie"
assert_json_field "$sticky_decision_second_schedule" '"human_decision_queue_position": 2' "decision portfolio still places the new same-class peer behind the retained primary packet"
home_output_sticky_decision=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output_sticky_decision" '"scheduler_decision_selection_policy": "sticky_primary_human_decision"' "home mirrors the sticky primary decision policy"
assert_json_field "$home_output_sticky_decision" '"scheduler_previous_primary_human_decision_key": "mission:issue-20026"' "home mirrors the previous primary decision key"
resolved_sticky_decision_signal_a=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-sticky-decision-a resolved "Sticky decision A covered.")
assert_json_field "$resolved_sticky_decision_signal_a" '"state": "resolved"' "sticky decision signal A can be resolved after portfolio coverage"
resolved_sticky_decision_signal_b=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-sticky-decision-b resolved "Sticky decision B covered.")
assert_json_field "$resolved_sticky_decision_signal_b" '"state": "resolved"' "sticky decision signal B can be resolved after portfolio coverage"
resolved_sticky_decision_run_a=$(run_in_repo bash scripts/harness-runs.sh set-state "$sticky_decision_run_a" completed "Closed sticky decision mission A.")
assert_json_field "$resolved_sticky_decision_run_a" '"state": "completed"' "sticky decision mission A can be closed after portfolio coverage"
resolved_sticky_decision_run_b=$(run_in_repo bash scripts/harness-runs.sh set-state "$sticky_decision_run_b" completed "Closed sticky decision mission B.")
assert_json_field "$resolved_sticky_decision_run_b" '"state": "completed"' "sticky decision mission B can be closed after portfolio coverage"

section "Scheduler decision materialization"
decision_materialization_claim=$(run_in_repo bash scripts/harness-claim.sh 20032 "Harness Decision Materialization" --owns="file:docs/harness-session-handoff.md")
assert_json_field "$decision_materialization_claim" '"mission_id": "issue-20032"' "decision materialization mission can materialize"
decision_materialization_run=$(printf '%s\n' "$decision_materialization_claim" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
decision_materialization_signal=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-decision-materialization",
  "kind": "requires_meeting",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$decision_materialization_run",
  "mission_id": "issue-20032",
  "applies_to": ["issue-20032"],
  "scope_type": "mission",
  "title": "Scheduler-backed council request",
  "summary": "This mission now requires a materialized decision meeting from the scheduler portfolio.",
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-30T03:00:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$decision_materialization_signal" '"signal_id": "sig-smoke-decision-materialization"' "decision materialization signal can be emitted directly"
decision_materialization_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$decision_materialization_schedule" "\"mission_id\": \"issue-20032\"" "scheduler surfaces the decision materialization mission"
assert_json_field "$decision_materialization_schedule" '"recommended_decision_surface": "council"' "scheduler recommends a council surface for the materialized decision"
decision_materialization_primary=$(run_in_repo bash scripts/harness-meetings.sh materialize-primary-decision --mission=issue-20032)
assert_json_field "$decision_materialization_primary" '"meeting_id": "meeting-' "scheduler primary decision can materialize into a meeting"
assert_json_field "$decision_materialization_primary" '"source_decision_id": "decision-' "materialized meeting records its source decision id"
assert_json_field "$decision_materialization_primary" '"source_decision_scope": "mission"' "materialized meeting records its source decision scope"
assert_json_field "$decision_materialization_primary" '"source_decision_kind": "review_follow_up_signals"' "materialized meeting records the scheduler decision kind"
assert_json_field "$decision_materialization_primary" '"source_decision_snapshot_id": "schedule-' "materialized meeting records the source scheduler snapshot"
assert_json_field "$decision_materialization_primary" '"source_decision_packet": {' "materialized meeting persists the source decision packet"
assert_json_field "$decision_materialization_primary" '"decision_protocol": "incident-command"' "materialized meeting keeps the scheduler-recommended decision protocol"
assert_json_field "$decision_materialization_primary" '"state": "waiting"' "host run moves to waiting when the scheduler decision is materialized"
primary_materialized_meeting_id=$(printf '%s\n' "$decision_materialization_primary" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
primary_materialized_decision_id=$(printf '%s\n' "$decision_materialization_primary" | sed -n 's/^[[:space:]]*"decision_id": "\(decision-[^"]*\)",$/\1/p' | head -n1)
decision_materialization_active_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$decision_materialization_active_schedule" "\"decision_id\": \"$primary_materialized_decision_id\"" "materialized scheduler decision keeps a stable decision id after council activation"
assert_json_field "$decision_materialization_active_schedule" '"recommended_decision_surface": "active-council"' "scheduler upgrades materialized decisions into active-council surface"
assert_json_field "$decision_materialization_active_schedule" "\"active_meeting_ids\": \\[" "scheduler projects active meeting ids back onto the materialized decision packet"
assert_json_field "$decision_materialization_active_schedule" "\"$primary_materialized_meeting_id\"" "scheduler ties the materialized decision packet to the active meeting id"
assert_json_field "$decision_materialization_active_schedule" '"active_decision_meeting_count": 1' "mission classifications inherit the active-council meeting count from the materialized decision packet"
decision_materialization_reused=$(run_in_repo bash scripts/harness-meetings.sh materialize-primary-decision --mission=issue-20032)
assert_json_field "$decision_materialization_reused" "\"meeting_id\": \"$primary_materialized_meeting_id\"" "materializing an already-active decision reuses the existing active meeting"
assert_json_field "$decision_materialization_reused" '"reused_existing_meeting": true' "reused decision materialization reports active-meeting reuse"
decision_materialization_reused_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$decision_materialization_reused_schedule" "\"$primary_materialized_meeting_id\"" "reused decision materialization keeps the same active meeting id in scheduler state"
assert_json_field "$decision_materialization_reused_schedule" '"active_decision_meeting_count": 1' "reused decision materialization does not duplicate active-council meeting count"
closed_primary_materialized_meeting=$(run_in_repo bash scripts/harness-drive.sh meeting-close "$decision_materialization_run" "$primary_materialized_meeting_id" "Primary scheduler decision materialization covered.")
assert_json_field "$closed_primary_materialized_meeting" '"state": "closed"' "primary scheduler-materialized meeting can close cleanly"
decision_materialization_by_id=$(run_in_repo bash scripts/harness-meetings.sh materialize-decision "$primary_materialized_decision_id" --mission=issue-20032)
assert_json_field "$decision_materialization_by_id" "\"decision_id\": \"$primary_materialized_decision_id\"" "decision-specific meeting materialization resolves the stable decision id"
assert_json_field "$decision_materialization_by_id" '"source_decision_surface": "active-council"' "decision-specific meeting carries the scheduler decision surface into the meeting"
assert_json_field "$decision_materialization_by_id" '"source_decision_snapshot_ref":' "decision-specific meeting records the scheduler snapshot reference"
secondary_materialized_meeting_id=$(printf '%s\n' "$decision_materialization_by_id" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
closed_secondary_materialized_meeting=$(run_in_repo bash scripts/harness-drive.sh meeting-close "$decision_materialization_run" "$secondary_materialized_meeting_id" "Decision-id materialization covered.")
assert_json_field "$closed_secondary_materialized_meeting" '"state": "closed"' "decision-id materialized meeting can close cleanly"
resolved_decision_materialization_signal=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-decision-materialization resolved "Decision materialization coverage completed.")
assert_json_field "$resolved_decision_materialization_signal" '"state": "resolved"' "decision materialization signal can be resolved after council coverage"
resolved_decision_materialization_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$decision_materialization_run" completed "Closed decision materialization mission.")
assert_json_field "$resolved_decision_materialization_run" '"state": "completed"' "decision materialization mission can close cleanly after coverage"

section "Active council progress priority"
progress_claim_a=$(run_in_repo bash scripts/harness-claim.sh 20033 "Harness Active Council Progress A" --owns="file:docs/harness-operating-model.md")
assert_json_field "$progress_claim_a" '"mission_id": "issue-20033"' "active-council progress mission A can materialize"
progress_run_a=$(printf '%s\n' "$progress_claim_a" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
progress_claim_b=$(run_in_repo bash scripts/harness-claim.sh 20034 "Harness Active Council Progress B" --owns="file:docs/harness-architecture.md")
assert_json_field "$progress_claim_b" '"mission_id": "issue-20034"' "active-council progress mission B can materialize"
progress_run_b=$(printf '%s\n' "$progress_claim_b" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
progress_signal_a=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-progress-a",
  "kind": "requires_meeting",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$progress_run_a",
  "mission_id": "issue-20033",
  "applies_to": ["issue-20033"],
  "scope_type": "mission",
  "title": "Progress council A",
  "summary": "Mission A needs a governed decision meeting.",
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-30T06:00:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$progress_signal_a" '"signal_id": "sig-smoke-progress-a"' "active-council progress signal A can be emitted"
progress_signal_b=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-signals.sh emit
{
  "signal_id": "sig-smoke-progress-b",
  "kind": "requires_meeting",
  "source_type": "leader",
  "source_id": "leader-main",
  "run_id": "$progress_run_b",
  "mission_id": "issue-20034",
  "applies_to": ["issue-20034"],
  "scope_type": "mission",
  "title": "Progress council B",
  "summary": "Mission B needs a governed decision meeting.",
  "authority_level": "L2",
  "confidence": "high",
  "emitted_at": "2026-03-30T06:05:00Z",
  "fresh_until": "2099-01-01T00:00:00Z",
  "state": "active"
}
EOF
)
assert_json_field "$progress_signal_b" '"signal_id": "sig-smoke-progress-b"' "active-council progress signal B can be emitted"
progress_schedule_initial=$(run_in_repo bash scripts/harness-missions.sh schedule)
progress_decision_a_id=$(printf '%s\n' "$progress_schedule_initial" | perl -0ne 'print "$1\n" if /\{[^{}]*"decision_id": "([^"]+)"[^{}]*"mission_id": "issue-20033"[^{}]*\}/s' | head -n1)
progress_decision_b_id=$(printf '%s\n' "$progress_schedule_initial" | perl -0ne 'print "$1\n" if /\{[^{}]*"decision_id": "([^"]+)"[^{}]*"mission_id": "issue-20034"[^{}]*\}/s' | head -n1)
progress_meeting_a=$(run_in_repo bash scripts/harness-meetings.sh materialize-decision "$progress_decision_a_id" --mission=issue-20033)
assert_json_field "$progress_meeting_a" '"meeting_id": "meeting-' "active-council progress meeting A can materialize"
progress_meeting_a_id=$(printf '%s\n' "$progress_meeting_a" | sed -n 's/.*"meeting_id": "\(meeting-[^"]*\)".*/\1/p' | head -n1)
progress_meeting_b=$(run_in_repo bash scripts/harness-meetings.sh materialize-decision "$progress_decision_b_id" --mission=issue-20034)
assert_json_field "$progress_meeting_b" '"meeting_id": "meeting-' "active-council progress meeting B can materialize"
progress_meeting_b_id=$(printf '%s\n' "$progress_meeting_b" | sed -n 's/.*"meeting_id": "\(meeting-[^"]*\)".*/\1/p' | head -n1)
run_in_repo bash scripts/harness-meetings.sh add-participant "$progress_meeting_a_id" "$progress_run_b" "consequence-checker" >/dev/null
run_in_repo bash scripts/harness-meetings.sh orchestrate "$progress_meeting_a_id" >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$progress_meeting_a_id" "$progress_run_a" independent_review support "Prefer bounded decision A." --claim=bounded-decision >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$progress_meeting_a_id" "$progress_run_b" independent_review support "Confirm bounded decision A." --claim=bounded-decision >/dev/null
run_in_repo bash scripts/harness-meetings.sh orchestrate "$progress_meeting_a_id" >/dev/null
run_in_repo bash scripts/harness-meetings.sh orchestrate "$progress_meeting_a_id" >/dev/null
progress_priority_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$progress_priority_schedule" "\"decision_id\": \"$progress_decision_a_id\"" "scheduler keeps the progressed active council as the primary human-decision packet"
assert_json_field "$progress_priority_schedule" '"active_decision_phase": "synthesis"' "scheduler projects the more advanced active-council phase into the primary packet"
assert_json_field "$progress_priority_schedule" '"active_decision_progress_rank": 4' "scheduler assigns a higher progress rank to synthesis-phase active councils than setup peers"
assert_json_field "$progress_priority_schedule" "\"mission_id\": \"issue-20034\"" "slower active-council peers remain visible in the portfolio"
resolved_progress_signal_a=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-progress-a resolved "Active council progress A covered.")
assert_json_field "$resolved_progress_signal_a" '"state": "resolved"' "active-council progress signal A can be resolved"
resolved_progress_signal_b=$(run_in_repo bash scripts/harness-signals.sh set-state sig-smoke-progress-b resolved "Active council progress B covered.")
assert_json_field "$resolved_progress_signal_b" '"state": "resolved"' "active-council progress signal B can be resolved"
closed_progress_meeting_a=$(run_in_repo bash scripts/harness-drive.sh meeting-close "$progress_run_a" "$progress_meeting_a_id" "Closed progressed active council A.")
assert_json_field "$closed_progress_meeting_a" '"state": "closed"' "progressed active council A can close cleanly"
closed_progress_meeting_b=$(run_in_repo bash scripts/harness-drive.sh meeting-close "$progress_run_b" "$progress_meeting_b_id" "Closed baseline active council B.")
assert_json_field "$closed_progress_meeting_b" '"state": "closed"' "baseline active council B can close cleanly"
resolved_progress_run_a=$(run_in_repo bash scripts/harness-runs.sh set-state "$progress_run_a" completed "Closed active-council progress mission A.")
assert_json_field "$resolved_progress_run_a" '"state": "completed"' "active-council progress mission A can close after scheduler coverage"
resolved_progress_run_b=$(run_in_repo bash scripts/harness-runs.sh set-state "$progress_run_b" completed "Closed active-council progress mission B.")
assert_json_field "$resolved_progress_run_b" '"state": "completed"' "active-council progress mission B can close after scheduler coverage"

parallel_claim_output=$(run_in_repo bash scripts/harness-claim.sh 20004 "Harness Parallel Mission" --owns="file:go.mod")
assert_json_field "$parallel_claim_output" '"mission_id": "issue-20004"' "independent parallel claim can materialize"
parallel_run_id=$(printf '%s\n' "$parallel_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
parallel_claim_output_extra=$(run_in_repo bash scripts/harness-claim.sh 20016 "Harness Parallel Mission Extra" --owns="file:README.md")
assert_json_field "$parallel_claim_output_extra" '"mission_id": "issue-20016"' "extra independent parallel claim can materialize"
parallel_run_id_extra=$(printf '%s\n' "$parallel_claim_output_extra" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
parallel_justification_primary=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$parallel_run_id"
{
  "parallel_justification": "Independent module change with isolated ownership boundary."
}
EOF
)
assert_json_field "$parallel_justification_primary" '"parallel_justification": "Independent module change with isolated ownership boundary."' "parallel mission can record explicit parallel justification"
parallel_justification_extra=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$parallel_run_id_extra"
{
  "parallel_justification": "Independent docs update with isolated ownership boundary."
}
EOF
)
assert_json_field "$parallel_justification_extra" '"parallel_justification": "Independent docs update with isolated ownership boundary."' "extra parallel mission can record explicit parallel justification"
parallel_justification_root=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$run_id"
{
  "parallel_justification": "Root foreground mission remains bounded and independently justifiable."
}
EOF
)
assert_json_field "$parallel_justification_root" '"parallel_justification": "Root foreground mission remains bounded and independently justifiable."' "existing foreground mission can also record explicit parallel justification"
scheduler_parallel_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$scheduler_parallel_output" '"execution_mode": "scheduled_parallel_missions"' "scheduler recommends a parallel batch when multiple isolated missions are runnable"
assert_json_field "$scheduler_parallel_output" '"selected_mission_count": 2' "scheduler selects both isolated runnable missions for the batch"
assert_json_field "$scheduler_parallel_output" '"selection_policy": "sticky_foreground"' "parallel batch keeps the existing foreground mission when tie conditions otherwise match"
assert_json_field "$scheduler_parallel_output" '"previous_foreground_mission_id": "issue-20001"' "parallel batch records the previous foreground mission"
assert_json_field "$scheduler_parallel_output" "\"mission_id\": \"issue-20004\"" "scheduler keeps the independent parallel mission visible in the portfolio"
assert_json_field "$scheduler_parallel_output" '"batch_admission_policy": "bounded_parallel_capacity"' "scheduler marks bounded parallel capacity as the admission policy"
assert_json_field "$scheduler_parallel_output" '"capacity_deferred_count": 1' "scheduler defers overflow missions behind the configured parallel capacity"
assert_json_field "$scheduler_parallel_output" '"action_kind": "defer_due_parallel_capacity"' "scheduler emits a reschedule action when runnable missions exceed parallel capacity"
parallel_claim_output_fresher=$(run_in_repo bash scripts/harness-claim.sh 20031 "Harness Parallel Mission Fresher" --owns="file:skills/README.md")
assert_json_field "$parallel_claim_output_fresher" '"mission_id": "issue-20031"' "fresher parallel peer can materialize"
parallel_run_id_fresher=$(printf '%s\n' "$parallel_claim_output_fresher" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
parallel_justification_fresher=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$parallel_run_id_fresher"
{
  "parallel_justification": "Independent changelog update with isolated ownership boundary."
}
EOF
)
assert_json_field "$parallel_justification_fresher" '"parallel_justification": "Independent changelog update with isolated ownership boundary."' "fresher parallel peer can record explicit parallel justification"
sticky_parallel_batch_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$sticky_parallel_batch_output" '"selection_policy": "sticky_parallel_batch"' "scheduler retains the admitted parallel batch when a same-class fresher peer appears"
assert_json_field "$sticky_parallel_batch_output" '"previous_selected_mission_ids":' "scheduler records the previous admitted batch"
assert_json_field "$sticky_parallel_batch_output" "\"issue-20016\"" "sticky parallel batch keeps the earlier admitted peer selected"
assert_json_field "$sticky_parallel_batch_output" '"capacity_deferred_count": 2' "sticky parallel batch still defers peers behind bounded parallel capacity"
assert_json_field "$sticky_parallel_batch_output" '"defer_due_parallel_capacity"' "sticky parallel batch still emits bounded-capacity defer actions"
resolved_parallel_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$parallel_run_id" completed "Resolved parallel scheduling smoke path.")
assert_json_field "$resolved_parallel_run" '"state": "completed"' "independent parallel mission can be closed after scheduling coverage"
resolved_parallel_run_extra=$(run_in_repo bash scripts/harness-runs.sh set-state "$parallel_run_id_extra" completed "Resolved overflow parallel scheduling smoke path.")
assert_json_field "$resolved_parallel_run_extra" '"state": "completed"' "overflow independent parallel mission can be closed after capacity coverage"
resolved_parallel_run_fresher=$(run_in_repo bash scripts/harness-runs.sh set-state "$parallel_run_id_fresher" completed "Resolved fresher overflow parallel scheduling smoke path.")
assert_json_field "$resolved_parallel_run_fresher" '"state": "completed"' "fresher overflow parallel mission can be closed after sticky-batch coverage"

section "Parallel justification gate"
parallel_gate_claim_a=$(run_in_repo bash scripts/harness-claim.sh 20021 "Harness Ungrounded Parallel A" --owns="file:docs/README.md")
assert_json_field "$parallel_gate_claim_a" '"mission_id": "issue-20021"' "unguarded parallel mission A can materialize"
parallel_gate_run_a=$(printf '%s\n' "$parallel_gate_claim_a" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
parallel_gate_claim_b=$(run_in_repo bash scripts/harness-claim.sh 20022 "Harness Ungrounded Parallel B" --owns="file:package-lock.json")
assert_json_field "$parallel_gate_claim_b" '"mission_id": "issue-20022"' "unguarded parallel mission B can materialize"
parallel_gate_run_b=$(printf '%s\n' "$parallel_gate_claim_b" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
schedule_parallel_gate_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_parallel_gate_output" '"batch_admission_policy": "serialize_missing_parallel_justification"' "scheduler reduces parallelism when candidate missions lack explicit parallel justification"
assert_json_field "$schedule_parallel_gate_output" '"action_kind": "defer_due_missing_parallel_justification"' "scheduler emits a reschedule action when parallel widening lacks justification"
resolved_parallel_gate_run_a=$(run_in_repo bash scripts/harness-runs.sh set-state "$parallel_gate_run_a" completed "Resolved missing parallel justification smoke path A.")
assert_json_field "$resolved_parallel_gate_run_a" '"state": "completed"' "unguarded parallel mission A can be closed after justification-gate coverage"
resolved_parallel_gate_run_b=$(run_in_repo bash scripts/harness-runs.sh set-state "$parallel_gate_run_b" completed "Resolved missing parallel justification smoke path B.")
assert_json_field "$resolved_parallel_gate_run_b" '"state": "completed"' "unguarded parallel mission B can be closed after justification-gate coverage"

section "Sticky foreground scheduling"
sticky_primary_claim=$(run_in_repo bash scripts/harness-claim.sh 20024 "Harness Sticky Foreground A" --owns="file:archive/README.md")
assert_json_field "$sticky_primary_claim" '"mission_id": "issue-20024"' "sticky foreground mission A can materialize"
sticky_primary_run=$(printf '%s\n' "$sticky_primary_claim" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
sticky_first_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$sticky_first_schedule" '"selection_policy": "sticky_foreground"' "scheduler retains the established foreground mission when a tied peer appears"
assert_json_field "$sticky_first_schedule" '"previous_foreground_mission_id": "issue-20001"' "sticky foreground scheduling records the established foreground mission"
assert_json_field "$sticky_first_schedule" '"selected_mission_count": 1' "sticky foreground schedule emits a foreground batch"
sticky_secondary_claim=$(run_in_repo bash scripts/harness-claim.sh 20025 "Harness Sticky Foreground B" --owns="file:docs/radar/themes.md")
assert_json_field "$sticky_secondary_claim" '"mission_id": "issue-20025"' "sticky foreground mission B can materialize after the first scheduler pass"
sticky_secondary_run=$(printf '%s\n' "$sticky_secondary_claim" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
sticky_second_schedule=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$sticky_second_schedule" '"selection_policy": "sticky_foreground"' "scheduler keeps the previous foreground mission when tie conditions otherwise match"
assert_json_field "$sticky_second_schedule" '"previous_foreground_mission_id": "issue-20001"' "scheduler records the previous foreground mission in the snapshot"
assert_json_field "$sticky_second_schedule" '"batch_admission_policy": "serialize_missing_parallel_justification"' "sticky foreground scheduling still obeys the parallel-justification gate"
assert_json_field "$sticky_second_schedule" "\"mission_id\": \"issue-20001\"" "sticky foreground scheduling keeps the established foreground mission visible in the classifications"
assert_json_field "$sticky_second_schedule" '"selected": true' "sticky foreground scheduling leaves the retained foreground mission selected"
assert_json_field "$sticky_second_schedule" '"action_kind": "defer_due_missing_parallel_justification"' "sticky foreground scheduling still defers peer missions behind the active foreground batch"
resolved_sticky_primary=$(run_in_repo bash scripts/harness-runs.sh set-state "$sticky_primary_run" completed "Resolved sticky foreground mission A.")
assert_json_field "$resolved_sticky_primary" '"state": "completed"' "sticky foreground mission A can be closed after scheduler coverage"
resolved_sticky_secondary=$(run_in_repo bash scripts/harness-runs.sh set-state "$sticky_secondary_run" completed "Resolved sticky foreground mission B.")
assert_json_field "$resolved_sticky_secondary" '"state": "completed"' "sticky foreground mission B can be closed after scheduler coverage"

section "Profile-backed high-risk serialization"
high_risk_claim_output=$(run_in_repo bash scripts/harness-claim.sh 20017 "Harness High-Risk Profile Mission" --owns="file:install.sh")
assert_json_field "$high_risk_claim_output" '"mission_id": "issue-20017"' "profile-backed high-risk mission can materialize"
high_risk_run_id=$(printf '%s\n' "$high_risk_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
low_risk_claim_output=$(run_in_repo bash scripts/harness-claim.sh 20018 "Harness Low-Risk Peer Mission" --owns="file:go.mod")
assert_json_field "$low_risk_claim_output" '"mission_id": "issue-20018"' "low-risk peer mission can materialize next to a profile-backed high-risk mission"
low_risk_run_id=$(printf '%s\n' "$low_risk_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
schedule_high_risk_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_high_risk_output" '"batch_admission_policy": "serialize_high_risk"' "scheduler serializes the batch when a profile-backed high-risk mission is runnable"
assert_json_field "$schedule_high_risk_output" '"mission_id": "issue-20017"' "scheduler includes the profile-backed high-risk mission in the classification output"
assert_json_field "$schedule_high_risk_output" '"profile_high_risk": true' "scheduler marks profile-backed high-risk missions explicitly"
assert_json_field "$schedule_high_risk_output" '"profile_high_risk_scope_refs":' "scheduler records which profile high-risk scopes triggered the escalation"
assert_json_field "$schedule_high_risk_output" '"effective_risk_level": "high"' "scheduler upgrades effective risk when a profile high-risk scope is matched"
assert_json_field "$schedule_high_risk_output" '"action_kind": "defer_until_foreground_risk_clears"' "scheduler defers lower-risk missions behind the profile-backed high-risk foreground mission"
resolved_high_risk_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$high_risk_run_id" completed "Resolved profile-backed high-risk scheduling smoke path.")
assert_json_field "$resolved_high_risk_run" '"state": "completed"' "profile-backed high-risk mission can be closed after serialization coverage"
resolved_low_risk_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$low_risk_run_id" completed "Resolved low-risk peer scheduling smoke path.")
assert_json_field "$resolved_low_risk_run" '"state": "completed"' "low-risk peer mission can be closed after high-risk serialization coverage"

section "P0 serialization"
p0_claim_output=$(run_in_repo bash scripts/harness-claim.sh 20019 "Harness P0 Mission" --owns="file:package.json")
assert_json_field "$p0_claim_output" '"mission_id": "issue-20019"' "P0 mission can materialize"
p0_run_id=$(printf '%s\n' "$p0_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
p0_update_output=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$p0_run_id"
{
  "priority": "P0"
}
EOF
)
assert_json_field "$p0_update_output" '"priority": "P0"' "run priority can be raised to P0"
p0_peer_claim_output=$(run_in_repo bash scripts/harness-claim.sh 20020 "Harness P0 Peer Mission" --owns="file:docs/README.md")
assert_json_field "$p0_peer_claim_output" '"mission_id": "issue-20020"' "P0 peer mission can materialize"
p0_peer_run_id=$(printf '%s\n' "$p0_peer_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
schedule_p0_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_p0_output" '"batch_admission_policy": "serialize_p0_priority"' "scheduler serializes the foreground batch when a P0 mission is runnable"
assert_json_field "$schedule_p0_output" '"selected_mission_count": 1' "scheduler keeps the foreground batch singular when a P0 mission is active"
assert_json_field "$schedule_p0_output" '"mission_id": "issue-20019"' "scheduler includes the P0 mission in the classification output"
assert_json_field "$schedule_p0_output" '"action_kind": "defer_until_p0_clears"' "scheduler defers peer missions behind the active P0 mission"
resolved_p0_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$p0_run_id" completed "Resolved P0 serialization smoke path.")
assert_json_field "$resolved_p0_run" '"state": "completed"' "P0 mission can be closed after serialization coverage"
resolved_p0_peer_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$p0_peer_run_id" completed "Resolved P0 peer serialization smoke path.")
assert_json_field "$resolved_p0_peer_run" '"state": "completed"' "P0 peer mission can be closed after serialization coverage"

section "Council-governed mission graph authority"
council_meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$run_id" "Resolve the overlapping release-note ownership boundary before widening execution.")
assert_json_field "$council_meeting_output" '"state": "active"' "drive meeting materializes an active council on the root mission"
council_meeting_id=$(printf '%s\n' "$council_meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
council_claim_output=$(run_in_repo bash scripts/harness-claim.sh 20007 "Harness Council Governed Mission" --owns="file:CHANGELOG.md" --allow-overlap)
assert_json_field "$council_claim_output" '"mission_id": "issue-20007"' "a second overlapping mission can materialize under an active council anchor"
council_run_id=$(printf '%s\n' "$council_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
graph_authority_council=$(run_in_repo bash scripts/harness-missions.sh authority-audit --all)
assert_json_field "$graph_authority_council" '"cross_mission_conflict_count": 0' "mission graph authority treats a single active council anchor as governed rather than unresolved conflict"
assert_json_field "$graph_authority_council" '"cross_mission_adjudication_count": 1' "mission graph authority records one council-governed mission adjudication"
assert_json_field "$graph_authority_council" '"governed_bundle_count": 1' "mission graph authority reports one governed bundle under council control"
assert_json_field "$graph_authority_council" '"resolution_kind": "council_governed_bundle"' "mission graph authority gives the council-governed bundle a distinct bundle-level resolution kind"
assert_json_field "$graph_authority_council" '"resolution_kind": "council_governed_mission_ownership"' "mission graph authority explains the council-governed ownership verdict"
assert_json_field "$graph_authority_council" "\"governing_mission_id\": \"issue-20001\"" "mission graph authority keeps bundle control on the council-host mission"
assert_json_field "$graph_authority_council" "\"required_action\": \"Keep bundle authority anchored to issue-20001 until its active council conclusions are adopted.\"" "mission graph authority tells the leader to keep authority on the active council host until adoption"
resolved_council_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$council_run_id" completed "Resolved council-governed overlap smoke path.")
assert_json_field "$resolved_council_run" '"state": "completed"' "council-governed overlapping mission can be closed after authority coverage"
closed_council_meeting=$(run_in_repo bash scripts/harness-drive.sh meeting-close "$run_id" "$council_meeting_id" "Council-governed overlap resolved.")
assert_json_field "$closed_council_meeting" '"state": "closed"' "council-governed meeting can close cleanly after bundle authority coverage"

section "Contract authority"
(
  cd "$TEST_REPO" &&
    perl -0pi -e 's/"branch":\s*"mission\/20001-harness-smoke-mission"/"branch": "stale-contract-branch"/' "$contract_path"
)
contract_authority_output=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_id")
assert_json_field "$contract_authority_output" '"artifact_authority_conflict_count": 1' "run authority audit counts stale contract as artifact-level authority conflict"
assert_json_field "$contract_authority_output" '"resolution_kind": "run_state_overrides_contract"' "run authority audit explains that active run state overrides stale contract state"
assert_json_field "$contract_authority_output" "\"losing_source\": \"contract:$run_id.branch\"" "run authority audit records the stale contract branch as the losing source"
contract_repair=$(
  cd "$TEST_REPO" &&
    cat <<EOF | bash scripts/harness-runs.sh update "$run_id"
{
  "branch": "$run_branch"
}
EOF
)
assert_json_field "$contract_repair" "\"branch\": \"$run_branch\"" "run update resynchronizes the contract snapshot"
contract_authority_recovered=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_id")
assert_json_field "$contract_authority_recovered" '"artifact_authority_conflict_count": 0' "artifact authority conflict clears after contract resync"

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
assert_json_field "$live_audit_output" '"resolution_kind": "live_state_overrides_run"' "run authority audit explains that live branch state overrides stale run state"
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

section "Live worktree missing"
(
  cd "$TEST_REPO"
  git worktree remove --force "$worktree_path"
)
missing_worktree_audit=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_id")
assert_json_field "$missing_worktree_audit" '"claim_ref": "live_worktree_missing"' "run authority audit reports missing live worktree when the run still expects one"
doctor_output_missing_worktree=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output_missing_worktree" '"kind": "run-live-state-drift"' "doctor surfaces missing live worktree as runtime drift"
reprovision_restored=$(run_in_repo bash scripts/harness-drive.sh reprovision "$run_id" --mode=worktree --apply)
assert_json_field "$reprovision_restored" "\"worktree_path\": \"$worktree_path\"" "reprovision can restore a missing governed worktree"
missing_worktree_recovered=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_id")
assert_json_field "$missing_worktree_recovered" '"live_state_gap_count": 0' "run authority audit clears after the missing worktree is restored"

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
assert_json_field "$home_output" '"scheduler_execution_mode": "solo"' "home surfaces the scheduler execution mode"
assert_json_field "$home_output" '"scheduler_selected_mission_count": 1' "home surfaces the scheduler batch size"

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
