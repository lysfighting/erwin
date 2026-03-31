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

meeting_file_path() {
  printf '%s/.harness/meetings/active/%s.yml\n' "$TEST_REPO" "$1"
}

rewrite_json_in_file() {
  local file="$1"
  local perl_expr="$2"
  perl -0pi -e "$perl_expr" "$file"
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

assert_not_json_field() {
  local json="$1"
  local needle="$2"
  local label="$3"
  if grep -q "$needle" <<<"$json"; then
    fail "$label"
    echo "$json"
  else
    pass "$label"
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

section "Scheduler"
schedule_output=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_output" '"execution_mode": "managed_delegation"' "scheduler recommends managed delegation when one mission has bounded helper runs"
assert_json_field "$schedule_output" '"selected_mission_count": 1' "scheduler keeps one foreground mission under managed delegation"
assert_json_field "$schedule_output" '"mission_id": "issue-33001"' "scheduler includes the delegated mission in its classification output"

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
authority_a_after_transfer=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_a")
authority_b_after_transfer=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$run_b")
assert_json_field "$authority_a_after_transfer" '"contract_gap_count": 0' "target teammate contract stays aligned after ownership transfer"
assert_json_field "$authority_b_after_transfer" '"contract_gap_count": 0' "source teammate contract stays aligned after ownership transfer"

section "Shared transfer policy"
shared_a=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-shared-a interface-owner "Own auth interface boundary" --scope="auth interface" --owns="interface:auth-scope" --allow-overlap)
shared_b=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-shared-b interface-reviewer "Review auth interface boundary" --scope="auth interface review" --owns="scope:auth-interface-review" --allow-overlap)
assert_json_field "$shared_a" '"owner_id": "worker-shared-a"' "shared source teammate assignment created"
assert_json_field "$shared_b" '"owner_id": "worker-shared-b"' "shared target teammate assignment created"
shared_run_a=$(printf '%s\n' "$shared_a" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
shared_run_b=$(printf '%s\n' "$shared_b" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
shared_transfer=$(run_in_repo bash scripts/harness-runs.sh transfer-ownership "$shared_run_a" "$shared_run_b" --owns="interface:auth-scope" --mode=copy --reason="Share auth interface scope for paired implementation and review.")
assert_json_field "$shared_transfer" '"mode": "copy"' "ownership transfer supports copy mode for shared boundaries"
shared_a_after=$(run_in_repo bash scripts/harness-teammates.sh get "$shared_run_a")
shared_b_after=$(run_in_repo bash scripts/harness-teammates.sh get "$shared_run_b")
assert_json_field "$shared_a_after" '"selector": "interface:auth-scope"' "copy-mode transfer keeps selector on source run"
assert_json_field "$shared_b_after" '"selector": "interface:auth-scope"' "copy-mode transfer adds selector to target run"
shared_authority_a=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$shared_run_a")
shared_authority_b=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$shared_run_b")
assert_json_field "$shared_authority_a" '"contract_gap_count": 0' "copy-mode transfer keeps source contract aligned"
assert_json_field "$shared_authority_b" '"contract_gap_count": 0' "copy-mode transfer keeps target contract aligned"

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

semantic_owner=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-docs completeness-checker "Prepare release-doc updates" --scope="release docs" --owns="dir:docs/release")
assert_json_field "$semantic_owner" '"owner_id": "worker-docs"' "directory selector assignment created"
docs_run=$(printf '%s\n' "$semantic_owner" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
set +e
semantic_overlap_output=$(cd "$TEST_REPO" && bash scripts/harness-teammates.sh assign "$leader_run_id" worker-doc-file completeness-checker "Patch release changelog" --scope="release changelog" --owns="file:docs/release/CHANGELOG.md" 2>&1)
semantic_rc=$?
set -e
if [ "$semantic_rc" -ne 0 ] && grep -q 'file:docs/release/CHANGELOG.md~dir:docs/release' <<<"$semantic_overlap_output"; then
  pass "semantic directory/file overlap is rejected"
else
  fail "semantic directory/file overlap is rejected"
  echo "$semantic_overlap_output"
fi
set +e
semantic_glob_overlap_output=$(cd "$TEST_REPO" && bash scripts/harness-teammates.sh assign "$leader_run_id" worker-doc-glob completeness-checker "Patch release markdown" --scope="release markdown" --owns="glob:docs/release/*.md" 2>&1)
semantic_glob_rc=$?
set -e
if [ "$semantic_glob_rc" -ne 0 ] && grep -q 'glob:docs/release/\*.md~dir:docs/release' <<<"$semantic_glob_overlap_output"; then
  pass "semantic glob/directory overlap is rejected"
else
  fail "semantic glob/directory overlap is rejected"
  echo "$semantic_glob_overlap_output"
fi
semantic_interface_owner=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-interface interface-owner "Own release contract boundary" --scope="release contract" --owns="interface:release-contract")
assert_json_field "$semantic_interface_owner" '"owner_id": "worker-interface"' "semantic interface selector assignment created"
set +e
semantic_scope_interface_overlap_output=$(cd "$TEST_REPO" && bash scripts/harness-teammates.sh assign "$leader_run_id" worker-scope contract-reviewer "Review release contract boundary" --scope="release contract" --owns="scope:release-contract" 2>&1)
semantic_scope_interface_rc=$?
set -e
if [ "$semantic_scope_interface_rc" -ne 0 ] && grep -q 'scope:release-contract~interface:release-contract' <<<"$semantic_scope_interface_overlap_output"; then
  pass "semantic scope/interface overlap is rejected"
else
  fail "semantic scope/interface overlap is rejected"
  echo "$semantic_scope_interface_overlap_output"
fi
semantic_hierarchy_owner=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-hierarchy hierarchy-owner "Own release notes subtree" --scope="release notes subtree" --owns="scope:release/notes")
assert_json_field "$semantic_hierarchy_owner" '"owner_id": "worker-hierarchy"' "semantic hierarchy selector assignment created"
semantic_hierarchy_run=$(printf '%s\n' "$semantic_hierarchy_owner" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
set +e
semantic_hierarchy_overlap_output=$(cd "$TEST_REPO" && bash scripts/harness-teammates.sh assign "$leader_run_id" worker-hierarchy-parent hierarchy-reviewer "Review broader release boundary" --scope="release boundary" --owns="interface:release" 2>&1)
semantic_hierarchy_rc=$?
set -e
if [ "$semantic_hierarchy_rc" -ne 0 ] && grep -q 'interface:release~scope:release/notes' <<<"$semantic_hierarchy_overlap_output"; then
  pass "semantic hierarchy overlap is rejected"
else
  fail "semantic hierarchy overlap is rejected"
  echo "$semantic_hierarchy_overlap_output"
fi

retired_hierarchy_owner=$(run_in_repo bash scripts/harness-runs.sh set-state "$semantic_hierarchy_run" completed "Hierarchy overlap guard verified; retire temporary child owner before split-transfer coverage.")
assert_json_field "$retired_hierarchy_owner" '"state": "completed"' "temporary hierarchy owner can retire before split transfer coverage"
hierarchy_reconcile=$(run_in_repo bash scripts/harness-maintenance.sh reconcile)
assert_json_field "$hierarchy_reconcile" '"reconciled_signal_ids"' "retiring the temporary hierarchy owner can reconcile superseded signals"

section "Split transfer policy"
split_source=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-split-parent hierarchy-owner "Own broad release boundary" --scope="release")
assert_json_field "$split_source" '"owner_id": "worker-split-parent"' "split source teammate assignment created"
split_source_run=$(printf '%s\n' "$split_source" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
split_transfer=$(run_in_repo bash scripts/harness-runs.sh transfer-ownership "$split_source_run" "$docs_run" --owns="scope:release/notes" --mode=split --reason="Delegate the release-notes subtree without losing the broader release boundary.")
assert_json_field "$split_transfer" '"mode": "split"' "ownership transfer supports split mode for hierarchical boundaries"
split_source_after=$(run_in_repo bash scripts/harness-runs.sh get "$split_source_run")
docs_after_split=$(run_in_repo bash scripts/harness-runs.sh get "$docs_run")
assert_json_field "$split_source_after" '"selector": "scope:release"' "split-mode transfer keeps the broader selector on source run"
assert_json_field "$split_source_after" '"ownership_exclusions"' "split-mode transfer records ownership exclusions on the source run"
assert_json_field "$split_source_after" '"scope:release/notes"' "split-mode transfer excludes the delegated child boundary from the source run"
assert_json_field "$docs_after_split" '"selector": "scope:release/notes"' "split-mode transfer adds delegated child selector to target run"
split_authority_source=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$split_source_run")
split_authority_target=$(run_in_repo bash scripts/harness-runs.sh authority-audit "$docs_run")
assert_json_field "$split_authority_source" '"contract_gap_count": 0' "split-mode transfer keeps source contract aligned"
assert_json_field "$split_authority_target" '"contract_gap_count": 0' "split-mode transfer keeps target contract aligned"

section "Reporting and communication"
report_output=$(run_in_repo bash scripts/harness-teammates.sh report "$run_a" "Auth caller patch is ready for sibling review" --kind=progress --evidence=src/auth.ts)
assert_json_field "$report_output" '"kind": "progress"' "teammate report emits progress signal"
assert_json_field "$report_output" '"latest_summary": "Auth caller patch is ready for sibling review"' "teammate report updates run summary"

message_output=$(run_in_repo bash scripts/harness-teammates.sh message "$run_a" "$run_b" interface_change "Caller signature now requires scope" --evidence=src/auth.ts)
assert_json_field "$message_output" '"kind": "interface_change"' "teammate message emits sibling signal"
assert_json_field "$message_output" "\"$run_b\"" "teammate message targets sibling run"
message_signal_id=$(printf '%s\n' "$message_output" | sed -n 's/^[[:space:]]*"signal_id": "\(sig-[^"]*\)",$/\1/p' | head -n1)

section "Mission authority"
mission_control_output=$(run_in_repo bash scripts/harness-missions.sh get issue-33001)
assert_json_field "$mission_control_output" '"latest_run_id": "run-.*-teammate-' "mission latest_run_id can reflect the freshest teammate activity"
assert_json_field "$mission_control_output" "\"authoritative_run_id\": \"$leader_run_id\"" "mission authoritative_run_id remains leader-governed even after teammate updates"
mission_authority_output=$(run_in_repo bash scripts/harness-missions.sh authority-audit issue-33001)
assert_json_field "$mission_authority_output" "\"authoritative_run_id\": \"$leader_run_id\"" "mission authority audit reports leader-governed mission control"
assert_json_field "$mission_authority_output" '"adjudication_count": 1' "mission authority audit records a leader-over-teammate adjudication"
assert_json_field "$mission_authority_output" '"resolution_kind": "leader_run_controls_mission"' "mission authority audit explains that the leader run controls mission state"
assert_json_field "$mission_authority_output" '"unresolved_conflict_count": 0' "mission authority audit reports no unresolved root-run conflict in the delegated case"

section "Meeting escalation"
escalate_output=$(run_in_repo bash scripts/harness-teammates.sh escalate "$run_b" "Do we need a release-note update for auth scope?" --materialize)
assert_json_field "$escalate_output" '"kind": "requires_meeting"' "teammate escalation emits meeting signal"
assert_json_field "$escalate_output" '"state": "active"' "teammate escalation materializes active meeting"
meeting_signal_id=$(printf '%s\n' "$escalate_output" | sed -n 's/^[[:space:]]*"signal_id": "\(sig-[^"]*\)",$/\1/p' | head -n1)
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
assert_json_field "$orchestrate_a" '"open_wave_action_count": 2' "meeting orchestration plans the missing independent-review wave obligations"

position_a=$(run_in_repo bash scripts/harness-meetings.sh add-position "$meeting_id" "$run_a" independent_review support "No schema drift in auth callers." --claim=release-note-scope --dimension=factual --evidence=src/auth.ts)
position_b=$(run_in_repo bash scripts/harness-meetings.sh add-position "$meeting_id" "$run_b" independent_review challenge "Release notes should mention the new scope requirement." --claim=release-note-scope --dimension=factual --evidence=CHANGELOG.md)
assert_json_field "$position_a" '"positions"' "meeting stores first position"
assert_json_field "$position_b" '"independent_review"' "meeting stores second independent-review position"
position_a_id=$(printf '%s\n' "$position_a" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)
position_b_id=$(printf '%s\n' "$position_b" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)

orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$meeting_id")
assert_json_field "$orchestrate_b" '"new_phase": "cross_examination"' "meeting orchestration advances to cross-examination"
assert_json_field "$orchestrate_b" '"open_wave_action_count": 0' "meeting orchestration clears independent-review wave obligations once both positions exist"
assert_json_field "$orchestrate_b" '"open_route_action_count": 1' "meeting orchestration plans one open debate route for the conflicting claim lane"
assert_json_field "$orchestrate_b" '"assigned_role": "ops-simulator"' "meeting debate route assigns the bounded challenger role"

cross_exam=$(run_in_repo bash scripts/harness-meetings.sh add-position "$meeting_id" "$run_b" cross_examination challenge "Operator-facing release notes must mention the new scope requirement." --claim=release-note-scope --dimension=factual --evidence=CHANGELOG.md --responds-to-run="$run_a" --responds-to-role=implementation-verifier)
assert_json_field "$cross_exam" '"cross_examination"' "meeting stores cross-examination position"
assert_json_field "$cross_exam" "\"responds_to_run_id\": \"$run_a\"" "cross-examination position records explicit response target"
cross_exam_id=$(printf '%s\n' "$cross_exam" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)

rejected_finding=$(run_in_repo bash scripts/harness-meetings.sh add-finding "$meeting_id" rejected "Ship immediately without release note updates." --reason="Cross-examination found operator-facing behavior drift." --authority="L2" --evidence=CHANGELOG.md --run="$run_b" --role=ops-simulator --challenges="$cross_exam_id")
assert_json_field "$rejected_finding" '"reason": "Cross-examination found operator-facing behavior drift."' "meeting records rejected finding reason"
accepted_finding=$(run_in_repo bash scripts/harness-meetings.sh add-finding "$meeting_id" accepted "Release notes need a scoped auth update before ship." --reason="Independent review and cross-examination agree on operator impact." --authority="L2" --evidence=CHANGELOG.md --run="$run_a" --role=implementation-verifier --claim=release-note-scope --dimension=factual --supports="$position_a_id|$cross_exam_id")
assert_json_field "$accepted_finding" '"accepted"' "meeting records accepted findings"
lower_authority_conflict=$(run_in_repo bash scripts/harness-meetings.sh add-finding "$meeting_id" unresolved "Release notes might be optional for this scope change." --reason="Earlier hypothesis before cross-examination." --authority="L3" --evidence=CHANGELOG.md --run="$run_b" --role=ops-simulator --claim=release-note-scope --dimension=factual --supports="$position_b_id")
assert_json_field "$lower_authority_conflict" '"unresolved"' "meeting records lower-authority conflicting finding"
action_output=$(run_in_repo bash scripts/harness-meetings.sh add-action "$meeting_id" "Update release notes before merging the auth change.")
assert_json_field "$action_output" '"action_items"' "meeting records action items"

orchestrate_c=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$meeting_id")
assert_json_field "$orchestrate_c" '"new_phase": "output"' "meeting orchestration can close directly to output once synthesis is already ready"
assert_json_field "$orchestrate_c" '"status": "ready_for_output"' "meeting orchestration materializes synthesis audit"
summary_output=$(run_in_repo bash scripts/harness-meetings.sh summary "$meeting_id")
assert_json_field "$summary_output" '"decision_protocol": "evidence-review"' "meeting summary reports the default evidence-review protocol"
assert_json_field "$summary_output" '"synthesis_status": "ready_for_output"' "meeting summary reports synthesis readiness"
assert_json_field "$summary_output" '"phase_history_count":' "meeting summary reports phase history"
assert_json_field "$summary_output" '"authority_conflict_count": 1' "meeting summary reports authority conflicts"
assert_json_field "$summary_output" '"unresolved_authority_conflict_count": 0' "meeting summary reports no unresolved authority conflict when higher authority wins"
assert_json_field "$summary_output" '"traced_finding_count": 3' "meeting summary reports traced findings"
assert_json_field "$summary_output" '"untraced_finding_count": 0' "meeting summary reports no untraced findings"
assert_json_field "$summary_output" '"planned_route_count": 1' "meeting summary reports routed conflicting claim count"
assert_json_field "$summary_output" '"resolved_route_count": 1' "meeting summary reports resolved debate route count"
assert_json_field "$summary_output" '"open_route_action_count": 0' "meeting summary reports no open debate route obligations"
assert_json_field "$summary_output" '"unrouted_claim_count": 0' "meeting summary reports no unrouted conflict claim"
authority_audit=$(run_in_repo bash scripts/harness-meetings.sh authority-audit "$meeting_id")
assert_json_field "$authority_audit" '"resolution_kind": "authority_wins"' "authority audit explains higher-authority winner"

section "Automated council synthesis"
auto_meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$leader_run_id" "What should the council synthesize about auth release notes?")
assert_json_field "$auto_meeting_output" '"state": "active"' "second meeting materializes for auto-synthesis"
assert_json_field "$auto_meeting_output" '"decision_protocol": "evidence-review"' "default drive meeting materializes an evidence-review council"
auto_meeting_id=$(printf '%s\n' "$auto_meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)

run_in_repo bash scripts/harness-meetings.sh add-role "$auto_meeting_id" implementation-verifier --mission="State what the code path really changed." --allow="L1: live git state|L2: active runtime state" --forbid="L4+: stale notes" >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-role "$auto_meeting_id" ops-simulator --mission="Challenge operator-facing release impact." --allow="L1: live git state|L2: active runtime state" --forbid="Unverified memory" >/dev/null
run_in_repo bash scripts/harness-meetings.sh lock-decision "$auto_meeting_id" "Keep this council bounded to release-note and rollback-note questions." --reason="Auto-synthesis should stay inside explicit claim lanes." --authority=L2 --evidence=harness/profiles/erwin.yml >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$auto_meeting_id" "$run_a" implementation-verifier >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$auto_meeting_id" "$run_b" ops-simulator >/dev/null

auto_orchestrate_a=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$auto_meeting_id")
assert_json_field "$auto_orchestrate_a" '"new_phase": "independent_review"' "auto-synthesis council advances into independent review"
assert_json_field "$auto_orchestrate_a" '"open_wave_action_count": 2' "auto-synthesis council plans both independent-review wave obligations"

auto_pos_a=$(run_in_repo bash scripts/harness-meetings.sh add-position "$auto_meeting_id" "$run_a" independent_review support "Release notes need the scoped auth update." --claim=release-note-scope --dimension=factual --authority=L2 --evidence=CHANGELOG.md)
auto_pos_b=$(run_in_repo bash scripts/harness-meetings.sh add-position "$auto_meeting_id" "$run_b" independent_review challenge "Release notes might be optional for this scope change." --claim=release-note-scope --dimension=factual --authority=L2 --evidence=CHANGELOG.md)
auto_pos_c=$(run_in_repo bash scripts/harness-meetings.sh add-position "$auto_meeting_id" "$run_a" independent_review support "Rollback guidance should be explicit if operator behavior changed." --claim=rollback-plan-note --dimension=factual --authority=L2 --evidence=docs/release.md)
auto_pos_b_id=$(printf '%s\n' "$auto_pos_b" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)
auto_pos_c_id=$(printf '%s\n' "$auto_pos_c" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)

auto_orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$auto_meeting_id")
assert_json_field "$auto_orchestrate_b" '"new_phase": "cross_examination"' "auto-synthesis council advances into cross-examination"
assert_json_field "$auto_orchestrate_b" '"open_wave_action_count": 0' "auto-synthesis council clears wave obligations after both independent reviews land"
assert_json_field "$auto_orchestrate_b" '"open_route_action_count": 1' "auto-synthesis council emits an open route obligation before linked cross-examination"

auto_cross_a=$(run_in_repo bash scripts/harness-meetings.sh add-position "$auto_meeting_id" "$run_a" cross_examination support "Operator-facing auth scope changes require release-note coverage." --claim=release-note-scope --dimension=factual --authority=L1 --evidence=CHANGELOG.md --responds-to-position="$auto_pos_b_id")
auto_cross_b=$(run_in_repo bash scripts/harness-meetings.sh add-position "$auto_meeting_id" "$run_b" cross_examination challenge "Rollback guidance can wait until a real rollback path exists." --claim=rollback-plan-note --dimension=factual --authority=L2 --evidence=docs/release.md --responds-to-position="$auto_pos_c_id")
assert_json_field "$auto_cross_a" '"cross_examination"' "auto-synthesis council records linked cross-examination support"
assert_json_field "$auto_cross_b" '"cross_examination"' "auto-synthesis council records linked cross-examination challenge"
auto_cross_b_id=$(printf '%s\n' "$auto_cross_b" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)

auto_orchestrate_c=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$auto_meeting_id")
assert_json_field "$auto_orchestrate_c" '"new_phase": "cross_examination"' "auto-synthesis council stays in cross-examination when a second-round reply is required"
assert_json_field "$auto_orchestrate_c" '"open_round_action_count": 1' "auto-synthesis council emits one follow-up debate round obligation"
assert_json_field "$auto_orchestrate_c" '"round_number": 2' "auto-synthesis council plans an explicit second debate round"

auto_cross_c=$(run_in_repo bash scripts/harness-meetings.sh add-position "$auto_meeting_id" "$run_a" cross_examination support "Rollback guidance should still be explicit even before a concrete rollback path exists." --claim=rollback-plan-note --dimension=factual --authority=L3 --evidence=docs/release.md --responds-to-position="$auto_cross_b_id")
assert_json_field "$auto_cross_c" '"cross_examination"' "auto-synthesis council records the second-round support reply"
auto_cross_c_id=$(printf '%s\n' "$auto_cross_c" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)

auto_orchestrate_d=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$auto_meeting_id")
assert_json_field "$auto_orchestrate_d" '"new_phase": "cross_examination"' "auto-synthesis council can demand a third debate round when the second reply still loses"
assert_json_field "$auto_orchestrate_d" '"open_round_action_count": 1' "auto-synthesis council keeps one round obligation open for the third response"
assert_json_field "$auto_orchestrate_d" '"round_number": 3' "auto-synthesis council plans an explicit third debate round"

auto_cross_d=$(run_in_repo bash scripts/harness-meetings.sh add-position "$auto_meeting_id" "$run_b" cross_examination challenge "If rollback guidance remains speculative, keep it out until a concrete rollback path exists." --claim=rollback-plan-note --dimension=factual --authority=L2 --evidence=docs/release.md --responds-to-position="$auto_cross_c_id")
assert_json_field "$auto_cross_d" '"cross_examination"' "auto-synthesis council records the third-round challenge reply"

auto_orchestrate_e=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$auto_meeting_id")
assert_json_field "$auto_orchestrate_e" '"new_phase": "output"' "auto-synthesis council can close the multi-round debate into output in the same orchestration pass"
assert_json_field "$auto_orchestrate_e" '"generated_finding_count": 2' "orchestrate auto-generates findings from council positions"
assert_json_field "$auto_orchestrate_e" '"generated_action_count": 1' "orchestrate auto-generates action items for unresolved council claims"

auto_summary=$(run_in_repo bash scripts/harness-meetings.sh summary "$auto_meeting_id")
assert_json_field "$auto_summary" '"accepted_count": 1' "auto-synthesis produces one accepted finding"
assert_json_field "$auto_summary" '"unresolved_count": 1' "auto-synthesis produces one unresolved finding"
assert_json_field "$auto_summary" '"auto_generated_finding_count": 2' "meeting summary surfaces generated finding count"
assert_json_field "$auto_summary" '"auto_generated_action_count": 1' "meeting summary surfaces generated action count"
assert_json_field "$auto_summary" '"traced_finding_count": 2' "auto-synthesized findings remain traceable to positions"
assert_json_field "$auto_summary" '"planned_route_count": 2' "auto-synthesis summary surfaces both routed claim lanes"
assert_json_field "$auto_summary" '"resolved_route_count": 2' "auto-synthesis summary shows both planned routes resolved"
assert_json_field "$auto_summary" '"open_route_action_count": 0' "auto-synthesis summary shows no remaining route obligation"
assert_json_field "$auto_summary" '"planned_round_count": 2' "auto-synthesis summary surfaces second-round planning across debated claim lanes"
assert_json_field "$auto_summary" '"resolved_round_count": 2' "auto-synthesis summary shows no remaining unresolved second-round planning"
assert_json_field "$auto_summary" '"open_round_action_count": 0' "auto-synthesis summary shows no remaining follow-up round obligation"
assert_json_field "$auto_summary" '"max_debate_rounds": 3' "auto-synthesis summary surfaces the bounded debate-round cap"
assert_json_field "$auto_summary" '"output_debate_lane_count": 2' "auto-synthesis summary surfaces both synthesized debate lanes inside the output packet"
assert_json_field "$auto_summary" '"output_accepted_debate_lane_count": 1' "auto-synthesis summary counts the accepted debate lane inside the output packet"
assert_json_field "$auto_summary" '"output_unresolved_debate_lane_count": 1' "auto-synthesis summary counts the unresolved debate lane inside the output packet"
assert_json_field "$auto_summary" '"output_open_debate_lane_count": 0' "auto-synthesis output packet does not leave live debate lanes open"
assert_json_field "$auto_summary" '"claim_ref": "release-note-scope"' "auto-synthesis output packet keeps the resolved release-note lane visible"
assert_json_field "$auto_summary" '"output_disposition": "accepted"' "auto-synthesis output packet records the accepted lane disposition"
assert_json_field "$auto_summary" '"claim_ref": "rollback-plan-note"' "auto-synthesis output packet keeps the unresolved rollback lane visible"
assert_json_field "$auto_summary" '"lane_status": "max_rounds_reached"' "auto-synthesis output packet records the bounded unresolved lane status"
assert_json_field "$auto_summary" '"output_disposition": "unresolved"' "auto-synthesis output packet records the unresolved lane disposition"

section "Debate verdict adoption"
verdict_meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$leader_run_id" "Can a governed debate route materialize a final verdict after the latest reply wins?")
assert_json_field "$verdict_meeting_output" '"state": "active"' "verdict-adoption meeting materializes"
verdict_meeting_id=$(printf '%s\n' "$verdict_meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)

run_in_repo bash scripts/harness-meetings.sh add-role "$verdict_meeting_id" implementation-verifier --mission="State the release-note requirement." --allow="L1: live git state|L2: active runtime state" --forbid="L4+: stale notes" >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-role "$verdict_meeting_id" ops-simulator --mission="Challenge the release-note requirement." --allow="L1: live git state|L2: active runtime state" --forbid="Unverified memory" >/dev/null
run_in_repo bash scripts/harness-meetings.sh lock-decision "$verdict_meeting_id" "Keep this council bounded to the stable release-note claim lane." --reason="Verdict-adoption coverage should stay inside one governed debate lane." --authority=L2 --evidence=harness/profiles/erwin.yml >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$verdict_meeting_id" "$run_a" implementation-verifier >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$verdict_meeting_id" "$run_b" ops-simulator >/dev/null

verdict_orchestrate_a=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$verdict_meeting_id")
assert_json_field "$verdict_orchestrate_a" '"new_phase": "independent_review"' "verdict-adoption council advances into independent review"

verdict_pos_a=$(run_in_repo bash scripts/harness-meetings.sh add-position "$verdict_meeting_id" "$run_a" independent_review support "Release notes should capture the stable operator-facing contract." --claim=stable-release-note --dimension=factual --authority=L2 --evidence=CHANGELOG.md)
verdict_pos_b=$(run_in_repo bash scripts/harness-meetings.sh add-position "$verdict_meeting_id" "$run_b" independent_review challenge "Stable release notes are optional for this contract change." --claim=stable-release-note --dimension=factual --authority=L2 --evidence=CHANGELOG.md)
verdict_pos_a_id=$(printf '%s\n' "$verdict_pos_a" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)
verdict_pos_b_id=$(printf '%s\n' "$verdict_pos_b" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)

verdict_orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$verdict_meeting_id")
assert_json_field "$verdict_orchestrate_b" '"new_phase": "cross_examination"' "verdict-adoption council advances into cross-examination"
assert_json_field "$verdict_orchestrate_b" '"open_route_action_count": 1' "verdict-adoption council plans an initial debate route"

verdict_cross_a=$(run_in_repo bash scripts/harness-meetings.sh add-position "$verdict_meeting_id" "$run_b" cross_examination challenge "Until operator impact broadens, release notes can stay narrow." --claim=stable-release-note --dimension=factual --authority=L2 --evidence=CHANGELOG.md --responds-to-position="$verdict_pos_a_id")
assert_json_field "$verdict_cross_a" '"cross_examination"' "verdict-adoption council records first debate reply"
verdict_cross_a_id=$(printf '%s\n' "$verdict_cross_a" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)

verdict_orchestrate_c=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$verdict_meeting_id")
assert_json_field "$verdict_orchestrate_c" '"open_round_action_count": 1' "verdict-adoption council plans a follow-up round"
assert_json_field "$verdict_orchestrate_c" '"round_number": 2' "verdict-adoption council asks for an explicit second-round reply"

verdict_cross_b=$(run_in_repo bash scripts/harness-meetings.sh add-position "$verdict_meeting_id" "$run_a" cross_examination support "The latest bounded reply still says stable operator behavior needs release-note coverage." --claim=stable-release-note --dimension=factual --authority=L2 --evidence=CHANGELOG.md --responds-to-position="$verdict_cross_a_id")
assert_json_field "$verdict_cross_b" '"cross_examination"' "verdict-adoption council records the latest winning reply"

verdict_orchestrate_d=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$verdict_meeting_id")
assert_json_field "$verdict_orchestrate_d" '"new_phase": "output"' "verdict-adoption council can land in output as soon as the governed verdict is synthesized"
assert_json_field "$verdict_orchestrate_d" '"resolution_kind": "latest_response_prevails"' "verdict-adoption debate plan records the latest-response winner"
assert_json_field "$verdict_orchestrate_d" '"prevailing_stance": "support"' "verdict-adoption debate plan records the prevailing stance"
assert_json_field "$verdict_orchestrate_d" '"generated_finding_count": 1' "verdict-adoption council auto-generates a final finding"
assert_json_field "$verdict_orchestrate_d" '"generated_action_count": 0' "verdict-adoption council does not leave unresolved follow-up work"

verdict_summary=$(run_in_repo bash scripts/harness-meetings.sh summary "$verdict_meeting_id")
assert_json_field "$verdict_summary" '"accepted_count": 1' "verdict-adoption summary materializes an accepted finding"
assert_json_field "$verdict_summary" '"unresolved_count": 0' "verdict-adoption summary avoids unresolved output when the latest reply wins"
assert_json_field "$verdict_summary" '"auto_generated_action_count": 0' "verdict-adoption summary leaves no auto-generated action item"

section "Authority-resolved council lanes"
authority_meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$leader_run_id" "Can higher-authority review settle a claim lane without forced cross-examination?")
assert_json_field "$authority_meeting_output" '"state": "active"' "authority-resolved meeting materializes"
authority_meeting_id=$(printf '%s\n' "$authority_meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)

run_in_repo bash scripts/harness-meetings.sh add-role "$authority_meeting_id" implementation-verifier --mission="State the operator-visible release note requirement." --allow="L1: live git state|L2: active runtime state" --forbid="L4+: stale notes" >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-role "$authority_meeting_id" ops-simulator --mission="Challenge whether the release note requirement is real." --allow="L2: active runtime state|L3: current planning notes" --forbid="Unverified memory" >/dev/null
run_in_repo bash scripts/harness-meetings.sh lock-decision "$authority_meeting_id" "Keep this council bounded to one release-note claim lane." --reason="Authority-resolution coverage should stay inside one governed claim." --authority=L2 --evidence=harness/profiles/erwin.yml >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$authority_meeting_id" "$run_a" implementation-verifier >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$authority_meeting_id" "$run_b" ops-simulator >/dev/null

authority_orchestrate_a=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$authority_meeting_id")
assert_json_field "$authority_orchestrate_a" '"new_phase": "independent_review"' "authority-resolved council advances into independent review"

run_in_repo bash scripts/harness-meetings.sh add-position "$authority_meeting_id" "$run_a" independent_review support "Live operator impact means release notes are required." --claim=authority-release-note --dimension=factual --authority=L1 --evidence=CHANGELOG.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$authority_meeting_id" "$run_b" independent_review challenge "Planning notes suggest release notes could stay optional." --claim=authority-release-note --dimension=factual --authority=L3 --evidence=docs/release.md >/dev/null

authority_orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$authority_meeting_id")
assert_json_field "$authority_orchestrate_b" '"new_phase": "output"' "authority-resolved council can skip straight to output when higher-authority evidence already wins"
assert_json_field "$authority_orchestrate_b" '"resolution_kind": "higher_authority_prevails"' "authority-resolved council records an explicit higher-authority verdict"
assert_json_field "$authority_orchestrate_b" '"resolved_without_cross_examination": true' "authority-resolved council marks the lane as settled without cross-examination"
assert_json_field "$authority_orchestrate_b" '"generated_finding_count": 1' "authority-resolved council auto-generates a final finding"

authority_summary=$(run_in_repo bash scripts/harness-meetings.sh summary "$authority_meeting_id")
assert_json_field "$authority_summary" '"accepted_count": 1' "authority-resolved summary materializes the higher-authority support finding"
assert_json_field "$authority_summary" '"resolved_without_cross_examination_count": 1' "authority-resolved summary surfaces the skipped cross-examination count"
assert_json_field "$authority_summary" '"authority_resolved_route_count": 1' "authority-resolved summary surfaces the authority-resolved route count"

section "Protocol-aware meeting materialization"
incident_meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$leader_run_id" "Should we roll back the release immediately?" --protocol=incident-command)
assert_json_field "$incident_meeting_output" '"decision_protocol": "incident-command"' "drive meeting accepts an explicit incident-command protocol override"
incident_meeting_id=$(printf '%s\n' "$incident_meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
incident_summary=$(run_in_repo bash scripts/harness-meetings.sh summary "$incident_meeting_id")
assert_json_field "$incident_summary" '"decision_protocol": "incident-command"' "meeting summary preserves the incident-command protocol"
run_in_repo bash scripts/harness-meetings.sh lock-decision "$incident_meeting_id" "Keep the incident council bounded to the immediate release rollback decision." --reason="Protocol-aware incident coverage should stay inside one rollback decision surface." --authority=L2 --evidence=harness/profiles/erwin.yml >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$incident_meeting_id" "$run_a" incident-commander >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$incident_meeting_id" "$run_b" operations-lead >/dev/null

incident_orchestrate_a=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$incident_meeting_id")
assert_json_field "$incident_orchestrate_a" '"new_phase": "independent_review"' "incident-command meeting advances into independent review"
assert_json_field "$incident_orchestrate_a" '"next_requirement": "Collect 2 bounded incident assessment(s)' "incident-command meeting uses protocol-specific independent-review guidance"

run_in_repo bash scripts/harness-meetings.sh add-position "$incident_meeting_id" "$run_a" independent_review support "We should roll back immediately to reduce release blast radius." --claim=incident-rollback --dimension=operational --authority=L2 --evidence=CHANGELOG.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$incident_meeting_id" "$run_b" independent_review challenge "We can hold the release if the operator note lands first." --claim=incident-rollback --dimension=operational --authority=L2 --evidence=docs/release.md >/dev/null

incident_orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$incident_meeting_id")
assert_json_field "$incident_orchestrate_b" '"new_phase": "output"' "incident-command meeting can go straight to output without a live debate route"
assert_json_field "$incident_orchestrate_b" '"protocol_deferred_route_count": 1' "incident-command meeting records protocol-deferred conflict handling"
assert_json_field "$incident_orchestrate_b" '"open_route_action_count": 0' "incident-command meeting does not emit cross-examination route obligations for same-authority conflict"
assert_json_field "$incident_orchestrate_b" '"generated_action_count": 1' "incident-command meeting turns unresolved conflict into bounded action ownership"
assert_json_field "$incident_orchestrate_b" '"next_requirement": "Adopt or close the meeting outcome after recording chosen action, owner, stop/go constraints, and the next checkpoint."' "incident-command meeting uses protocol-specific output guidance"
incident_output_summary=$(run_in_repo bash scripts/harness-meetings.sh summary "$incident_meeting_id")
assert_json_field "$incident_output_summary" '"output_packet_kind": "incident_command_output"' "incident-command summary surfaces a protocol-aware output packet"
assert_json_field "$incident_output_summary" '"chosen_action": "We should roll back immediately to reduce release blast radius."' "incident-command output packet captures the chosen action"
assert_json_field "$incident_output_summary" '"owner": "incident-commander"' "incident-command output packet captures the default owner"
assert_json_field "$incident_output_summary" '"output_follow_up_packet_count": 1' "incident-command output packet synthesizes one structured follow-up packet"
assert_json_field "$incident_output_summary" '"execution_kind": "incident_checkpoint"' "incident-command follow-up packet keeps protocol-specific execution semantics visible"
assert_json_field "$incident_output_summary" '"claim_ref": "incident-rollback"' "incident-command follow-up packet keeps the governed claim lane visible"

run_in_repo bash scripts/harness-meetings.sh set-protocol "$incident_meeting_id" decision-record >/dev/null
decision_record_summary=$(run_in_repo bash scripts/harness-meetings.sh summary "$incident_meeting_id")
assert_json_field "$decision_record_summary" '"decision_protocol": "decision-record"' "meetings set-protocol can retarget an active meeting to decision-record"

decision_record_meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$leader_run_id" "Which rollout note variant should we record?" --protocol=decision-record)
assert_json_field "$decision_record_meeting_output" '"decision_protocol": "decision-record"' "drive meeting accepts an explicit decision-record protocol override"
decision_record_meeting_id=$(printf '%s\n' "$decision_record_meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
run_in_repo bash scripts/harness-meetings.sh lock-decision "$decision_record_meeting_id" "Keep this record-oriented council bounded to the rollout-note tradeoff." --reason="Protocol-aware decision-record coverage should stay inside one documented tradeoff." --authority=L2 --evidence=harness/profiles/erwin.yml >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$decision_record_meeting_id" "$run_a" context-owner >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$decision_record_meeting_id" "$run_b" alternative-challenger >/dev/null

decision_record_orchestrate_a=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$decision_record_meeting_id")
assert_json_field "$decision_record_orchestrate_a" '"new_phase": "independent_review"' "decision-record meeting advances into independent review"
assert_json_field "$decision_record_orchestrate_a" '"next_requirement": "Collect 2 bounded decision-record assessment(s)' "decision-record meeting uses protocol-specific independent-review guidance"

run_in_repo bash scripts/harness-meetings.sh add-position "$decision_record_meeting_id" "$run_a" independent_review support "Record the broader rollout note so future operators see the auth scope change." --claim=decision-note-variant --dimension=tradeoff --authority=L2 --evidence=CHANGELOG.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$decision_record_meeting_id" "$run_b" independent_review challenge "Record the narrower rollout note to avoid overstating the blast radius." --claim=decision-note-variant --dimension=tradeoff --authority=L2 --evidence=docs/release.md >/dev/null

decision_record_orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$decision_record_meeting_id")
assert_json_field "$decision_record_orchestrate_b" '"new_phase": "output"' "decision-record meeting can go straight to output without opening a live debate route"
assert_json_field "$decision_record_orchestrate_b" '"protocol_deferred_route_count": 1' "decision-record meeting records protocol-deferred conflict handling"
assert_json_field "$decision_record_orchestrate_b" '"open_route_action_count": 0' "decision-record meeting does not emit cross-examination route obligations for same-authority conflict"
assert_json_field "$decision_record_orchestrate_b" '"generated_action_count": 1' "decision-record meeting records unresolved tradeoff follow-up instead of debate rounds"
assert_json_field "$decision_record_orchestrate_b" '"next_requirement": "Adopt or close the meeting outcome after recording the decision, alternatives, consequences, unresolved questions, and follow-up actions."' "decision-record meeting uses protocol-specific output guidance"
decision_record_output_summary=$(run_in_repo bash scripts/harness-meetings.sh summary "$decision_record_meeting_id")
assert_json_field "$decision_record_output_summary" '"output_packet_kind": "decision_record_output"' "decision-record summary surfaces a protocol-aware output packet"
assert_json_field "$decision_record_output_summary" '"decision": "Record the broader rollout note so future operators see the auth scope change."' "decision-record output packet captures the proposed decision"
assert_json_field "$decision_record_output_summary" '"alternatives":' "decision-record output packet captures alternatives explicitly"
assert_json_field "$decision_record_output_summary" '"output_follow_up_packet_count": 1' "decision-record output packet synthesizes one structured follow-up packet"
assert_json_field "$decision_record_output_summary" '"execution_kind": "decision_record_follow_up"' "decision-record follow-up packet keeps protocol-specific execution semantics visible"
assert_json_field "$decision_record_output_summary" '"claim_ref": "decision-note-variant"' "decision-record follow-up packet keeps the governed tradeoff lane visible"

section "Aligned review skips empty debate"
aligned_meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$leader_run_id" "Can aligned independent review go straight to synthesis when there is no conflict lane?")
assert_json_field "$aligned_meeting_output" '"state": "active"' "aligned-review meeting materializes"
aligned_meeting_id=$(printf '%s\n' "$aligned_meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)

run_in_repo bash scripts/harness-meetings.sh add-role "$aligned_meeting_id" implementation-verifier --mission="State the aligned operator note requirement." --allow="L1: live git state|L2: active runtime state" --forbid="L4+: stale notes" >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-role "$aligned_meeting_id" ops-simulator --mission="Provide an independent aligned review of the same note requirement." --allow="L1: live git state|L2: active runtime state" --forbid="Unverified memory" >/dev/null
run_in_repo bash scripts/harness-meetings.sh lock-decision "$aligned_meeting_id" "Keep this council bounded to one aligned release-note claim." --reason="Aligned-review coverage should prove debate is skipped when no conflict lane exists." --authority=L2 --evidence=harness/profiles/erwin.yml >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$aligned_meeting_id" "$run_a" implementation-verifier >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$aligned_meeting_id" "$run_b" ops-simulator >/dev/null

aligned_orchestrate_a=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$aligned_meeting_id")
assert_json_field "$aligned_orchestrate_a" '"new_phase": "independent_review"' "aligned-review council advances into independent review"

run_in_repo bash scripts/harness-meetings.sh add-position "$aligned_meeting_id" "$run_a" independent_review support "Release notes should mention the stable auth behavior." --claim=aligned-release-note --dimension=factual --authority=L2 --evidence=CHANGELOG.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$aligned_meeting_id" "$run_b" independent_review support "Independent review agrees the stable auth behavior needs release-note coverage." --claim=aligned-release-note --dimension=factual --authority=L2 --evidence=CHANGELOG.md >/dev/null

aligned_orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$aligned_meeting_id")
assert_json_field "$aligned_orchestrate_b" '"new_phase": "output"' "aligned-review council can skip straight to output when no conflicting claim lane exists"
assert_json_field "$aligned_orchestrate_b" '"planned_route_count": 0' "aligned-review council records zero debate routes when no conflict exists"
assert_json_field "$aligned_orchestrate_b" '"generated_finding_count": 1' "aligned-review council auto-generates the aligned finding"

aligned_summary=$(run_in_repo bash scripts/harness-meetings.sh summary "$aligned_meeting_id")
assert_json_field "$aligned_summary" '"accepted_count": 1' "aligned-review summary materializes the shared support finding"
assert_json_field "$aligned_summary" '"planned_route_count": 0' "aligned-review summary stays free of debate-route planning"

section "Home and doctor"
home_output=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_output" '"active_mission_count": 1' "home surfaces active mission count"
assert_json_field "$home_output" '"active_teammate_count": 7' "home surfaces active teammate runs"
assert_json_field "$home_output" '"coordination_signal_count":' "home surfaces coordination signal count"
assert_json_field "$home_output" '"dependency_blocked_count": 1' "home surfaces dependency-blocked runs"
assert_json_field "$home_output" '"ownership_overlap_count": 0' "home stays free of ownership overlap"
assert_json_field "$home_output" '"active_meeting_count":' "home surfaces active council count"
mission_graph=$(run_in_repo bash scripts/harness-missions.sh graph)
assert_json_field "$mission_graph" '"mission_count": 1' "mission graph reports the active mission"

doctor_output=$(run_in_repo bash scripts/harness-doctor.sh)
assert_json_field "$doctor_output" '"issue_count": 0' "doctor stays healthy with governed teammate runtime"

section "Multi-lane debate batching"
batch_meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$leader_run_id" "Can one bounded participant receive a multi-lane debate batch?")
assert_json_field "$batch_meeting_output" '"state": "active"' "multi-lane batching meeting materializes"
batch_meeting_id=$(printf '%s\n' "$batch_meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)

run_in_repo bash scripts/harness-meetings.sh add-role "$batch_meeting_id" implementation-verifier --mission="State the release and rollback note requirements." --allow="L1: live git state|L2: active runtime state" --forbid="L4+: stale notes" >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-role "$batch_meeting_id" ops-simulator --mission="Challenge both operator-facing note requirements." --allow="L1: live git state|L2: active runtime state" --forbid="Unverified memory" >/dev/null
run_in_repo bash scripts/harness-meetings.sh lock-decision "$batch_meeting_id" "Keep this council bounded to two documentation claim lanes." --reason="Multi-lane batching coverage should stay inside explicit release and rollback note claims." --authority=L2 --evidence=harness/profiles/erwin.yml >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$batch_meeting_id" "$run_a" implementation-verifier >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$batch_meeting_id" "$run_b" ops-simulator >/dev/null

batch_orchestrate_a=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$batch_meeting_id")
assert_json_field "$batch_orchestrate_a" '"new_phase": "independent_review"' "multi-lane batching council advances into independent review"

run_in_repo bash scripts/harness-meetings.sh add-position "$batch_meeting_id" "$run_a" independent_review support "Release notes need the scoped auth update." --claim=batch-release-note --dimension=factual --authority=L2 --evidence=CHANGELOG.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$batch_meeting_id" "$run_b" independent_review challenge "Release notes can stay narrower for this change." --claim=batch-release-note --dimension=factual --authority=L2 --evidence=CHANGELOG.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$batch_meeting_id" "$run_a" independent_review support "Rollback notes should mention the auth scope change." --claim=batch-rollback-note --dimension=factual --authority=L2 --evidence=docs/release.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$batch_meeting_id" "$run_b" independent_review challenge "Rollback notes can wait until a rollback path is concrete." --claim=batch-rollback-note --dimension=factual --authority=L2 --evidence=docs/release.md >/dev/null

batch_orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$batch_meeting_id")
assert_json_field "$batch_orchestrate_b" '"new_phase": "cross_examination"' "multi-lane batching council advances into cross-examination for conflicting lanes"
assert_json_field "$batch_orchestrate_b" '"open_route_action_count": 2' "multi-lane batching council opens two route obligations"
assert_json_field "$batch_orchestrate_b" '"wave_batch_count": 1' "multi-lane batching controller groups the two open route obligations into one wave batch"
assert_json_field "$batch_orchestrate_b" '"ready_wave_batch_count": 1' "multi-lane batching controller exposes one ready batch"
assert_json_field "$batch_orchestrate_b" '"blocked_wave_batch_count": 0' "multi-lane batching controller leaves no blocked batch when only one participant batch exists"
assert_json_field "$batch_orchestrate_b" '"current_wave_batch_id": "debate-batch-' "multi-lane batching controller exposes a deterministic current batch"
assert_json_field "$batch_orchestrate_b" '"multi_lane_batch_count": 1' "multi-lane batching controller records one batch that spans multiple claim lanes"
assert_json_field "$batch_orchestrate_b" '"max_wave_lane_count": 2' "multi-lane batching controller reports the largest lane bundle in one batch"

batch_summary=$(run_in_repo bash scripts/harness-meetings.sh summary "$batch_meeting_id")
assert_json_field "$batch_summary" '"open_route_action_count": 2' "multi-lane batching summary keeps both open route obligations visible"
assert_json_field "$batch_summary" '"wave_batch_count": 1' "multi-lane batching summary surfaces one open debate batch"
assert_json_field "$batch_summary" '"ready_wave_batch_count": 1' "multi-lane batching summary surfaces one ready batch"
assert_json_field "$batch_summary" '"blocked_wave_batch_count": 0' "multi-lane batching summary surfaces zero blocked batches"
assert_json_field "$batch_summary" '"current_wave_batch_summary":' "multi-lane batching summary surfaces the current batch summary"
assert_json_field "$batch_summary" '"current_wave_batch_assigned_role": "ops-simulator"' "multi-lane batching summary surfaces the current batch assignee"
assert_json_field "$batch_summary" '"current_wave_batch_action_count": 2' "multi-lane batching summary surfaces the current batch action count"
assert_json_field "$batch_summary" '"current_wave_batch_lane_count": 2' "multi-lane batching summary surfaces the current batch lane count"
assert_json_field "$batch_summary" '"current_wave_actions":' "multi-lane batching summary surfaces the current batch action packets"
assert_json_field "$batch_summary" '"assigned_role": "ops-simulator"' "multi-lane batching summary keeps the current batch action assignee visible"
assert_json_field "$batch_summary" '"multi_lane_batch_count": 1' "multi-lane batching summary reports a multi-lane batch"
assert_json_field "$batch_summary" '"max_wave_lane_count": 2' "multi-lane batching summary reports the lane count of the largest batch"

section "Stateful debate batch execution"
sequenced_meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$leader_run_id" "Can the council controller choose one current debate batch while blocking the next?")
assert_json_field "$sequenced_meeting_output" '"state": "active"' "stateful-batch meeting materializes"
sequenced_meeting_id=$(printf '%s\n' "$sequenced_meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
run_in_repo bash scripts/harness-meetings.sh orchestrate "$sequenced_meeting_id" >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$sequenced_meeting_id" "$leader_run_id" chair >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$sequenced_meeting_id" "$run_a" challenger-a >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$sequenced_meeting_id" "$run_b" challenger-b >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$sequenced_meeting_id" "$leader_run_id" independent_review support "Scope A must stay enforced." --claim=batch-lane-a --dimension=factual --authority=L2 --evidence=docs/a.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$sequenced_meeting_id" "$run_a" independent_review challenge "Scope A may be narrower than described." --claim=batch-lane-a --dimension=factual --authority=L2 --evidence=docs/a-review.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$sequenced_meeting_id" "$leader_run_id" independent_review support "Scope B also needs explicit coverage." --claim=batch-lane-b --dimension=factual --authority=L2 --evidence=docs/b.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$sequenced_meeting_id" "$run_b" independent_review challenge "Scope B may still be optional." --claim=batch-lane-b --dimension=factual --authority=L2 --evidence=docs/b-review.md >/dev/null

sequenced_orchestrate_a=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$sequenced_meeting_id")
assert_json_field "$sequenced_orchestrate_a" '"new_phase": "cross_examination"' "stateful-batch council advances into cross-examination"
assert_json_field "$sequenced_orchestrate_a" '"wave_batch_count": 2' "stateful-batch council creates two debate batches when different challengers own different lanes"
assert_json_field "$sequenced_orchestrate_a" '"ready_wave_batch_count": 1' "stateful-batch council exposes exactly one ready batch"
assert_json_field "$sequenced_orchestrate_a" '"blocked_wave_batch_count": 1' "stateful-batch council leaves the second batch blocked behind the current one"
assert_json_field "$sequenced_orchestrate_a" '"current_wave_batch_id": "debate-batch-' "stateful-batch council exposes a current batch id"
assert_json_field "$sequenced_orchestrate_a" '"current_wave_batch_lane_count": 1' "stateful-batch council shows the current batch covers one lane"
sequenced_current_batch_id=$(printf '%s\n' "$sequenced_orchestrate_a" | sed -n 's/^[[:space:]]*"current_wave_batch_id": "\(debate-batch-[^"]*\)",$/\1/p' | head -n1)

assign_c=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-alpha alpha-reviewer "Bound new council lane" --scope="alpha council lane" --owns="scope:alpha-council-lane")
assert_json_field "$assign_c" '"owner_id": "worker-alpha"' "stateful-batch council can add a third bounded participant"
run_c=$(printf '%s\n' "$assign_c" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
run_in_repo bash scripts/harness-meetings.sh add-position "$sequenced_meeting_id" "$leader_run_id" independent_review support "Scope 0 also requires explicit coverage." --claim=batch-lane-0 --dimension=factual --authority=L2 --evidence=docs/0.md >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-position "$sequenced_meeting_id" "$run_c" independent_review challenge "Scope 0 may be narrower than described." --claim=batch-lane-0 --dimension=factual --authority=L2 --evidence=docs/0-review.md >/dev/null

sequenced_orchestrate_sticky=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$sequenced_meeting_id")
assert_json_field "$sequenced_orchestrate_sticky" '"batch_selection_policy": "sticky_current_wave_batch"' "stateful-batch council keeps the current debate batch when a new same-phase lane appears"
assert_json_field "$sequenced_orchestrate_sticky" "\"previous_current_wave_batch_id\": \"$sequenced_current_batch_id\"" "stateful-batch council records the previous current batch id"
assert_json_field "$sequenced_orchestrate_sticky" "\"current_wave_batch_id\": \"$sequenced_current_batch_id\"" "stateful-batch council retains the previous current batch instead of reordering on the new lane"
assert_json_field "$sequenced_orchestrate_sticky" '"wave_batch_count": 3' "stateful-batch council tracks the newly inserted debate batch"
assert_json_field "$sequenced_orchestrate_sticky" '"blocked_wave_batch_count": 2' "stateful-batch council leaves the new lane blocked behind the retained current batch"

run_in_repo bash scripts/harness-meetings.sh set-phase "$sequenced_meeting_id" synthesis >/dev/null
sequenced_orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$sequenced_meeting_id")
assert_json_field "$sequenced_orchestrate_b" '"new_phase": "synthesis"' "stateful-batch council refuses to jump from synthesis to output while debate batches remain open"
assert_json_field "$sequenced_orchestrate_b" '"ready_wave_batch_count": 1' "stateful-batch council keeps one ready batch after the blocked synthesis attempt"
assert_json_field "$sequenced_orchestrate_b" '"blocked_wave_batch_count": 2' "stateful-batch council keeps the remaining batches blocked after the blocked synthesis attempt"
assert_json_field "$sequenced_orchestrate_b" 'Complete current debate batch before output.' "stateful-batch council tells the leader to finish the current batch before output"
sequenced_summary=$(run_in_repo bash scripts/harness-meetings.sh summary "$sequenced_meeting_id")
assert_json_field "$sequenced_summary" '"current_wave_actions":' "stateful-batch summary surfaces the current batch action packets"
assert_json_field "$sequenced_summary" '"blocked_wave_batches":' "stateful-batch summary surfaces blocked batch summaries"
assert_json_field "$sequenced_summary" '"assigned_role": "alpha-reviewer"' "stateful-batch summary keeps blocked batch ownership visible"
run_in_repo bash scripts/harness-drive.sh meeting-close "$leader_run_id" "$sequenced_meeting_id" "Stateful batch execution smoke complete." >/dev/null

section "Council governance gates"
governance_meeting_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$leader_run_id" "Can governance policy block council progress until authority hierarchy is restored?")
assert_json_field "$governance_meeting_output" '"state": "active"' "governance-gate meeting materializes"
governance_meeting_id=$(printf '%s\n' "$governance_meeting_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
governance_meeting_file=$(meeting_file_path "$governance_meeting_id")

rewrite_json_in_file "$governance_meeting_file" 's/"authority_hierarchy": \[[^\]]*\]/"authority_hierarchy": []/s'
governance_orchestrate_a=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$governance_meeting_id")
assert_json_field "$governance_orchestrate_a" '"new_phase": "setup"' "governance gaps keep setup-phase councils from advancing"
assert_json_field "$governance_orchestrate_a" 'missing authority hierarchy' "governance gaps report the missing authority hierarchy"
governance_summary_a=$(run_in_repo bash scripts/harness-meetings.sh summary "$governance_meeting_id")
assert_json_field "$governance_summary_a" '"governance_gap_count": 1' "meeting summary surfaces one authority-hierarchy governance gap"
assert_json_field "$governance_summary_a" 'missing authority hierarchy' "meeting summary names the authority-hierarchy governance gap"

rewrite_json_in_file "$governance_meeting_file" 's~"authority_hierarchy": \[\]~"authority_hierarchy": ["L1: live git/GitHub state","L2: active runtime state","L3: committed harness definitions"]~s'
run_in_repo bash scripts/harness-meetings.sh add-participant "$governance_meeting_id" "$run_a" implementation-verifier >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$governance_meeting_id" "$run_b" ops-simulator >/dev/null
governance_orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$governance_meeting_id")
assert_json_field "$governance_orchestrate_b" '"new_phase": "independent_review"' "restoring authority hierarchy re-enables council progression"

rejected_reason_output=$(run_in_repo bash scripts/harness-drive.sh meeting "$leader_run_id" "Can rejected findings block output until their governance record is complete?")
assert_json_field "$rejected_reason_output" '"state": "active"' "rejected-finding governance meeting materializes"
rejected_reason_meeting_id=$(printf '%s\n' "$rejected_reason_output" | sed -n 's/^[[:space:]]*"meeting_id": "\(meeting-[^"]*\)",$/\1/p' | head -n1)
rejected_reason_meeting_file=$(meeting_file_path "$rejected_reason_meeting_id")

run_in_repo bash scripts/harness-meetings.sh add-role "$rejected_reason_meeting_id" implementation-verifier --mission="State whether the broad rollout note is necessary." --allow="L1: live git state|L2: active runtime state" --forbid="L4+: stale notes" >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-role "$rejected_reason_meeting_id" ops-simulator --mission="Independently confirm whether the broad rollout note is over-scoped." --allow="L1: live git state|L2: active runtime state" --forbid="Unverified memory" >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$rejected_reason_meeting_id" "$run_a" implementation-verifier >/dev/null
run_in_repo bash scripts/harness-meetings.sh add-participant "$rejected_reason_meeting_id" "$run_b" ops-simulator >/dev/null
run_in_repo bash scripts/harness-meetings.sh orchestrate "$rejected_reason_meeting_id" >/dev/null
rejected_reason_pos_a=$(run_in_repo bash scripts/harness-meetings.sh add-position "$rejected_reason_meeting_id" "$run_a" independent_review support "Broad rollout notes are unnecessary for this scoped auth change." --claim=rejected-reason-scope --dimension=factual --authority=L2 --evidence=CHANGELOG.md)
rejected_reason_pos_b=$(run_in_repo bash scripts/harness-meetings.sh add-position "$rejected_reason_meeting_id" "$run_b" independent_review support "Independent review agrees the broader rollout note would overstate impact." --claim=rejected-reason-scope --dimension=factual --authority=L2 --evidence=CHANGELOG.md)
rejected_reason_pos_a_id=$(printf '%s\n' "$rejected_reason_pos_a" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)
rejected_reason_pos_b_id=$(printf '%s\n' "$rejected_reason_pos_b" | sed -n 's/^[[:space:]]*"position_id": "\(position-[^"]*\)",$/\1/p' | tail -n1)
run_in_repo bash scripts/harness-meetings.sh add-finding "$rejected_reason_meeting_id" rejected "Reject broad rollout note." --authority=L2 --evidence=CHANGELOG.md --run="$run_b" --role=ops-simulator --claim=rejected-reason-scope --dimension=factual --supports="$rejected_reason_pos_a_id|$rejected_reason_pos_b_id" >/dev/null
run_in_repo bash scripts/harness-meetings.sh set-phase "$rejected_reason_meeting_id" synthesis >/dev/null

rejected_reason_orchestrate_a=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$rejected_reason_meeting_id")
assert_json_field "$rejected_reason_orchestrate_a" '"new_phase": "synthesis"' "missing rejected-finding governance keeps synthesis from advancing to output"
assert_json_field "$rejected_reason_orchestrate_a" 'rejected finding missing reason' "orchestrate reports the missing rejected-finding reason"
rejected_reason_summary_a=$(run_in_repo bash scripts/harness-meetings.sh summary "$rejected_reason_meeting_id")
assert_json_field "$rejected_reason_summary_a" '"governance_gap_count": 1' "meeting summary surfaces one rejected-finding governance gap"
assert_json_field "$rejected_reason_summary_a" 'rejected finding missing reason' "meeting summary names the rejected-finding governance gap"

rewrite_json_in_file "$rejected_reason_meeting_file" 's/"reason": ""/"reason": "Aligned review concluded the broader rollout note would overstate impact."/'
rejected_reason_orchestrate_b=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$rejected_reason_meeting_id")
assert_json_field "$rejected_reason_orchestrate_b" '"new_phase": "output"' "restoring rejected-finding governance lets synthesis advance to output"
rejected_reason_summary_b=$(run_in_repo bash scripts/harness-meetings.sh summary "$rejected_reason_meeting_id")
assert_json_field "$rejected_reason_summary_b" '"governance_gap_count": 0' "meeting summary clears governance gaps after the rejected-finding reason is restored"

section "Closed meeting adoption"
decision_record_close=$(run_in_repo bash scripts/harness-drive.sh meeting-close "$leader_run_id" "$decision_record_meeting_id" "Decision record captured and adopted.")
assert_json_field "$decision_record_close" '"adopted_output_packet_kind": "decision_record_output"' "meeting-close adopts the structured decision-record output packet into run state"
assert_json_field "$decision_record_close" '"adopted_output_follow_up_packet_count": 1' "meeting-close adopts the structured output follow-up packet count into run state"
assert_json_field "$decision_record_close" '"adopted_primary_follow_up_packet": {' "meeting-close keeps the primary adopted follow-up packet in run state"
assert_json_field "$decision_record_close" '"execution_kind": "decision_record_follow_up"' "meeting-close preserves protocol-specific follow-up execution semantics in adopted run state"
decision_record_follow_up_packet_id=$(printf '%s\n' "$decision_record_close" | perl -0ne 'print "$1\n" if /"adopted_primary_follow_up_packet": \{.*?"follow_up_packet_id": "([^"]+)"/s')
decision_record_next_action=$(printf '%s\n' "$decision_record_close" | perl -0ne 'print "$1\n" if /"next_expected_action": "([^"]+)"/s')
decision_record_follow_up_action=$(printf '%s\n' "$decision_record_close" | perl -0ne 'print "$1\n" if /"adopted_primary_follow_up_packet": \{.*?"summary": "([^"]+)"/s')
if [[ -n "$decision_record_next_action" && "$decision_record_next_action" == "$decision_record_follow_up_action" ]]; then
  pass "meeting-close promotes human-readable adopted follow-up guidance into next_expected_action"
else
  fail "meeting-close promotes human-readable adopted follow-up guidance into next_expected_action"
  echo "$decision_record_close"
fi

schedule_after_adoption=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_after_adoption" '"adopted_output_packet_kind": "decision_record_output"' "scheduler mission packets surface the adopted structured output kind after meeting-close"
assert_json_field "$schedule_after_adoption" '"adopted_output_open_follow_up_packet_count": 1' "scheduler mission packets surface open adopted follow-up packet counts"
assert_json_field "$schedule_after_adoption" '"adopted_primary_follow_up_packet": {' "scheduler mission packets surface the primary adopted follow-up packet"
assert_json_field "$schedule_after_adoption" '"execution_kind": "decision_record_follow_up"' "scheduler mission packets preserve protocol-aware adopted follow-up execution semantics"
run_in_repo bash scripts/harness-signals.sh set-state "$message_signal_id" resolved "Decision-record adoption superseded the earlier interface-change review signal." >/dev/null
run_in_repo bash scripts/harness-signals.sh set-state "$meeting_signal_id" resolved "Decision-record adoption superseded the earlier meeting request." >/dev/null

schedule_after_signal_clear=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_after_signal_clear" 'Decision record must capture unresolved tradeoff and follow-up for claim: decision-note-variant' "scheduler keeps human-readable adopted follow-up guidance visible after earlier signals are cleared"

home_after_adoption=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_after_adoption" '"adopted_output_run_count": 1' "home surfaces active runs carrying adopted structured output"
assert_json_field "$home_after_adoption" '"primary_adopted_output_packet_kind": "decision_record_output"' "home surfaces the primary adopted output packet kind"
assert_json_field "$home_after_adoption" '"primary_adopted_follow_up_packet": {' "home surfaces the primary adopted follow-up packet"
assert_json_field "$home_after_adoption" '"execution_kind": "decision_record_follow_up"' "home preserves protocol-aware adopted follow-up execution semantics"
assert_json_field "$home_after_adoption" 'Decision record must capture unresolved tradeoff and follow-up for claim: decision-note-variant' "home keeps human-readable adopted follow-up guidance visible after earlier signals are cleared"

doctor_after_adoption=$(run_in_repo bash scripts/harness-doctor.sh)
assert_not_json_field "$doctor_after_adoption" '"kind": "run-missing-meeting-next-action"' "doctor no longer misreads adopted follow-up guidance as missing next_expected_action"
assert_not_json_field "$doctor_after_adoption" '"kind": "run-missing-meeting-summary"' "doctor no longer compares run summary against a non-governing closed meeting"

follow_up_assignment=$(run_in_repo bash scripts/harness-teammates.sh assign "$leader_run_id" worker-follow-up follow-up-runner --adopted-follow-up)
follow_up_run_id=$(printf '%s\n' "$follow_up_assignment" | sed -n 's/^[[:space:]]*"run_id": "\(run-[^"]*\)",$/\1/p' | head -n1)
assert_json_field "$follow_up_assignment" "\"delegated_from_follow_up_packet_id\": \"$decision_record_follow_up_packet_id\"" "teammates assign can bind the primary adopted follow-up packet directly"
assert_json_field "$follow_up_assignment" "\"delegated_from_meeting_id\": \"$decision_record_meeting_id\"" "teammate follow-up assignment records the governing meeting"
assert_json_field "$follow_up_assignment" "\"delegated_from_follow_up_execution_kind\": \"decision_record_follow_up\"" "teammate follow-up assignment preserves protocol-aware execution kind"
assert_json_field "$follow_up_assignment" "\"objective\": \"$decision_record_follow_up_action\"" "teammate follow-up assignment derives its objective from adopted follow-up guidance"
assert_json_field "$follow_up_assignment" "\"next_expected_action\": \"$decision_record_follow_up_action\"" "teammate follow-up assignment starts from the adopted follow-up action text"
schedule_after_follow_up_handoff=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_after_follow_up_handoff" "\"adopted_primary_follow_up_delegated_run_id\": \"$follow_up_run_id\"" "scheduler surfaces the delegated teammate run for the adopted follow-up"
assert_json_field "$schedule_after_follow_up_handoff" "\"adopted_primary_follow_up_delegated_owner\": \"worker-follow-up\"" "scheduler surfaces the delegated teammate owner for the adopted follow-up"
assert_json_field "$schedule_after_follow_up_handoff" "\"adopted_primary_follow_up_delegated_action\": \"Monitor delegated adopted follow-up on $follow_up_run_id (worker-follow-up) and keep the governing run bounded until it reports back.\"" "scheduler promotes monitor guidance once the adopted follow-up is delegated"
home_after_follow_up_handoff=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_after_follow_up_handoff" "\"primary_adopted_follow_up_delegated_run_id\": \"$follow_up_run_id\"" "home surfaces the delegated teammate run for the primary adopted follow-up"
assert_json_field "$home_after_follow_up_handoff" "\"primary_adopted_follow_up_delegated_owner\": \"worker-follow-up\"" "home surfaces the delegated teammate owner for the primary adopted follow-up"
assert_json_field "$home_after_follow_up_handoff" "\"primary_adopted_follow_up_delegated_action\": \"Monitor delegated adopted follow-up on $follow_up_run_id (worker-follow-up) and keep the governing run bounded until it reports back.\"" "home preserves delegated adopted follow-up monitoring guidance even when higher-priority control-plane work still owns next_action"
set +e
duplicate_follow_up_assignment_output=$(cd "$TEST_REPO" && bash scripts/harness-teammates.sh assign "$leader_run_id" worker-follow-up-2 follow-up-runner --adopted-follow-up 2>&1)
duplicate_follow_up_assignment_rc=$?
set -e
if [ "$duplicate_follow_up_assignment_rc" -ne 0 ] && grep -q "active teammate assignment already exists for adopted follow-up packet $decision_record_follow_up_packet_id" <<<"$duplicate_follow_up_assignment_output"; then
  pass "adopted follow-up handoff refuses duplicate active teammate assignments for the same packet"
else
  fail "adopted follow-up handoff refuses duplicate active teammate assignments for the same packet"
  echo "$duplicate_follow_up_assignment_output"
fi

blocked_completion_output=$(run_in_repo bash scripts/harness-drive.sh complete "$leader_run_id" "Should not close while adopted follow-up is still open." 2>&1 || true)
if grep -q 'run cannot transition to completed while adopted follow-up packets remain open' <<<"$blocked_completion_output"; then
  pass "drive complete is blocked while adopted follow-up packets remain open"
else
  fail "drive complete is blocked while adopted follow-up packets remain open"
  echo "$blocked_completion_output"
fi

decision_record_resolve=$(run_in_repo bash scripts/harness-drive.sh resolve-follow-up "$follow_up_run_id")
assert_json_field "$decision_record_resolve" '"delegated_resolution": true' "drive resolve-follow-up can resolve a delegated adopted follow-up from the teammate run"
assert_json_field "$decision_record_resolve" '"delegated_follow_up_state": "done"' "delegated teammate run records follow-up completion state"
assert_json_field "$decision_record_resolve" '"next_expected_action": "Report completion back to the leader or continue any remaining bounded execution."' "delegated teammate run advances into post-follow-up reporting guidance"
assert_json_field "$decision_record_resolve" '"next_expected_action": "Continue governed execution after the adopted follow-up is complete."' "governing run guidance still advances after delegated follow-up resolution"
schedule_after_follow_up_resolution=$(run_in_repo bash scripts/harness-missions.sh schedule)
assert_json_field "$schedule_after_follow_up_resolution" '"adopted_output_open_follow_up_packet_count": 0' "scheduler clears adopted open follow-up packet counts after action resolution"
home_after_follow_up_resolution=$(run_in_repo bash scripts/harness-home.sh)
assert_json_field "$home_after_follow_up_resolution" '"adopted_output_open_follow_up_packet_count": 0' "home clears adopted open follow-up packet counts after action resolution"
completed_after_follow_up=$(run_in_repo bash scripts/harness-drive.sh complete "$leader_run_id" "Decision-record follow-up completed and run can now close.")
assert_json_field "$completed_after_follow_up" '"state": "completed"' "drive complete succeeds after adopted follow-up is resolved"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
