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
claim_output=$(run_in_repo bash scripts/harness-claim.sh 20001 "Harness Smoke Mission" --owns="file:CHANGELOG.md")
assert_json_field "$claim_output" '"mission_id": "issue-20001"' "claim creates mission"
assert_json_field "$claim_output" '"run_id": "run-' "claim returns run id"
assert_json_field "$claim_output" 'issue-20001.yml' "claim materializes mission path"

run_id=$(printf '%s\n' "$claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
run_branch=$(printf '%s\n' "$claim_output" | sed -n 's/^[[:space:]]*"branch": "\(mission\/[^"]*\)",$/\1/p' | head -n1)
contract_path=$(printf '%s\n' "$claim_output" | sed -n 's/^[[:space:]]*"contract_path": "\(.*issue-[^"]*\.yml\)",$/\1/p' | head -n1)

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
  "kind": "requires_meeting",
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
schedule_multi_escalation_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_multi_escalation_output" '"bundle_role": "escalation-required"' "scheduler marks tied bundles as escalation-required"
assert_json_field "$schedule_multi_escalation_output" '"decision_kind": "establish_bundle_governor"' "scheduler escalates tied bundles into an explicit governing-mission decision"
assert_json_field "$schedule_multi_escalation_output" '"next_action": "Resolve required human decisions before widening execution; keep the current foreground batch bounded in the meantime."' "scheduler prioritizes bundle-governor decisions over foreground selection messaging"
assert_json_field "$schedule_multi_escalation_output" '"bundle_role_tally":' "scheduler reports bundle-role tallies in the portfolio snapshot"
home_output_multi_escalation=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output_multi_escalation" '"mission_escalated_bundle_count": 1' "home surfaces escalation-required mission bundles"
assert_json_field "$home_output_multi_escalation" '"next_action": "Resolve required human decisions before widening execution; keep the current foreground batch bounded in the meantime."' "home keeps scheduler decision-first guidance when bundle governance is still pending"
assert_json_field "$home_output_multi_escalation" '"scheduler_bundle_role_tally":' "home surfaces scheduler bundle-role tallies"

resolved_multi_escalation_a=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_escalation_run_a" completed "Closed escalated multi-mission bundle member.")
assert_json_field "$resolved_multi_escalation_a" '"state": "completed"' "first escalation bundle member can be closed"
resolved_multi_escalation_b=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_escalation_run_b" completed "Closed escalated multi-mission bundle member.")
assert_json_field "$resolved_multi_escalation_b" '"state": "completed"' "second escalation bundle member can be closed"
resolved_multi_escalation_c=$(run_in_repo bash scripts/harness-runs.sh set-state "$multi_escalation_run_c" completed "Closed escalated multi-mission bundle member.")
assert_json_field "$resolved_multi_escalation_c" '"state": "completed"' "third escalation bundle member can be closed"

parallel_claim_output=$(run_in_repo bash scripts/harness-claim.sh 20004 "Harness Parallel Mission" --owns="file:go.mod")
assert_json_field "$parallel_claim_output" '"mission_id": "issue-20004"' "independent parallel claim can materialize"
parallel_run_id=$(printf '%s\n' "$parallel_claim_output" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
parallel_claim_output_extra=$(run_in_repo bash scripts/harness-claim.sh 20016 "Harness Parallel Mission Extra" --owns="file:README.md")
assert_json_field "$parallel_claim_output_extra" '"mission_id": "issue-20016"' "extra independent parallel claim can materialize"
parallel_run_id_extra=$(printf '%s\n' "$parallel_claim_output_extra" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
scheduler_parallel_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$scheduler_parallel_output" '"execution_mode": "scheduled_parallel_missions"' "scheduler recommends a parallel batch when multiple isolated missions are runnable"
assert_json_field "$scheduler_parallel_output" '"selected_mission_count": 2' "scheduler selects both isolated runnable missions for the batch"
assert_json_field "$scheduler_parallel_output" '"mission_id": "issue-20004"' "scheduler includes the independent parallel mission in the batch"
assert_json_field "$scheduler_parallel_output" '"batch_admission_policy": "bounded_parallel_capacity"' "scheduler marks bounded parallel capacity as the admission policy"
assert_json_field "$scheduler_parallel_output" '"capacity_deferred_count": 1' "scheduler defers overflow missions behind the configured parallel capacity"
assert_json_field "$scheduler_parallel_output" '"action_kind": "defer_due_parallel_capacity"' "scheduler emits a reschedule action when runnable missions exceed parallel capacity"
resolved_parallel_run=$(run_in_repo bash scripts/harness-runs.sh set-state "$parallel_run_id" completed "Resolved parallel scheduling smoke path.")
assert_json_field "$resolved_parallel_run" '"state": "completed"' "independent parallel mission can be closed after scheduling coverage"
resolved_parallel_run_extra=$(run_in_repo bash scripts/harness-runs.sh set-state "$parallel_run_id_extra" completed "Resolved overflow parallel scheduling smoke path.")
assert_json_field "$resolved_parallel_run_extra" '"state": "completed"' "overflow independent parallel mission can be closed after capacity coverage"

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
