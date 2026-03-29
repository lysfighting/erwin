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
  if grep -q "$needle" <<<"$json"; then
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
assert_json_field "$orchestrate_c" '"new_phase": "synthesis"' "meeting orchestration advances to synthesis"
assert_json_field "$orchestrate_c" '"status": "ready_for_output"' "meeting orchestration materializes synthesis audit"
orchestrate_d=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$meeting_id")
assert_json_field "$orchestrate_d" '"new_phase": "output"' "meeting orchestration advances to output once findings exist"
summary_output=$(run_in_repo bash scripts/harness-meetings.sh summary "$meeting_id")
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
assert_json_field "$auto_orchestrate_e" '"new_phase": "synthesis"' "auto-synthesis council advances into synthesis after the bounded multi-round debate"
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

auto_orchestrate_f=$(run_in_repo bash scripts/harness-meetings.sh orchestrate "$auto_meeting_id")
assert_json_field "$auto_orchestrate_f" '"new_phase": "output"' "auto-synthesis council can advance from synthesis to output"

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

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
