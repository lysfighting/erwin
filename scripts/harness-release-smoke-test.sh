#!/usr/bin/env bash
# harness-release-smoke-test.sh — validates the ship/rc lifecycle
#
# USAGE:
#   ./scripts/harness-release-smoke-test.sh

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

TMPDIR_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/harness-release-smoke.XXXXXX")
trap 'rm -rf "$TMPDIR_ROOT"' EXIT

PASS=0
FAIL=0

pass() { PASS=$((PASS + 1)); echo "  PASS  $1"; }
fail() { FAIL=$((FAIL + 1)); echo "  FAIL  $1"; }
section() { echo ""; echo "=== $1 ==="; }

TEST_REPO="$TMPDIR_ROOT/repo"
ORIGIN_REPO="$TMPDIR_ROOT/origin.git"
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
  git init -q --bare "$ORIGIN_REPO" &&
  cd "$TEST_REPO" &&
    git init -q &&
    git checkout -qb main &&
    git config user.email "release-smoke@example.com" &&
    git config user.name "Harness Release Test" &&
    git remote add origin "$ORIGIN_REPO" &&
    chmod +x scripts/*.sh install-codex.sh &&
    printf "hello\n" > README.md &&
    git add . &&
    git commit -qm "init harness release smoke" &&
    git push -u origin main >/dev/null
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
claim_output=$(run_in_repo ./scripts/harness-claim.sh 32001 "Release Smoke Mission" --provision=branch --apply)
assert_json_field "$claim_output" '"mission_id": "issue-32001"' "claim creates release mission"
run_id=$(echo "$claim_output" | python3 -c 'import json,sys; print(json.load(sys.stdin)["run_id"])')

section "Ship"
verify_output=$(run_in_repo ./scripts/harness-ship.sh verify "$run_id" "targeted-test:pass")
assert_json_field "$verify_output" '"verification_status": "passed"' "ship verify records verification"

ready_output=$(run_in_repo ./scripts/harness-ship.sh ready "$run_id" "Ready for PR")
assert_json_field "$ready_output" '"ship_status": "ready"' "ship ready advances run"

push_output=$(run_in_repo ./scripts/harness-ship.sh push "$run_id")
assert_json_field "$push_output" '"ship_status": "pushed"' "ship push records pushed state"

if git --git-dir="$ORIGIN_REPO" show-ref --verify --quiet refs/heads/mission/32001-release-smoke-mission; then
  pass "ship push updates remote branch"
else
  fail "ship push updates remote branch"
fi

pr_draft_output=$(run_in_repo ./scripts/harness-ship.sh pr-draft "$run_id")
assert_json_field "$pr_draft_output" '"mode": "draft"' "ship pr-draft returns draft payload"
assert_json_field "$pr_draft_output" '"base_branch": "main"' "ship pr-draft uses base branch"

pr_output=$(run_in_repo ./scripts/harness-ship.sh pr "$run_id" "PR-32001" "https://example.invalid/pr/32001")
assert_json_field "$pr_output" '"ship_status": "in_review"' "ship pr records PR state"
assert_json_field "$pr_output" '"pr_ref": "PR-32001"' "ship pr stores PR ref"

review_requested_output=$(run_in_repo ./scripts/harness-ship.sh review "$run_id" requested "Review requested")
assert_json_field "$review_requested_output" '"review_status": "requested"' "ship review tracks requested state"

checks_pending_output=$(run_in_repo ./scripts/harness-ship.sh checks "$run_id" pending "ci:pending")
assert_json_field "$checks_pending_output" '"checks_status": "pending"' "ship checks tracks pending state"

review_changes_output=$(run_in_repo ./scripts/harness-ship.sh review "$run_id" changes_requested "Please address feedback")
assert_json_field "$review_changes_output" '"review_status": "changes_requested"' "ship review tracks requested changes"

checks_failed_output=$(run_in_repo ./scripts/harness-ship.sh checks "$run_id" failed "ci:failed")
assert_json_field "$checks_failed_output" '"checks_status": "failed"' "ship checks tracks failed state"

review_approved_output=$(run_in_repo ./scripts/harness-ship.sh review "$run_id" approved "Approved after fixes")
assert_json_field "$review_approved_output" '"review_status": "approved"' "ship review tracks approval"

checks_passed_output=$(run_in_repo ./scripts/harness-ship.sh checks "$run_id" passed "ci:passed")
assert_json_field "$checks_passed_output" '"checks_status": "passed"' "ship checks tracks passed state"
assert_json_field "$checks_passed_output" '"ship_status": "approved"' "ship checks marks merge-ready PR as approved"

merge_output=$(run_in_repo ./scripts/harness-ship.sh merge "$run_id" --via=git "Merged cleanly")
assert_json_field "$merge_output" '"ship_status": "merged"' "ship merge records merge"
assert_json_field "$merge_output" '"state": "completed"' "ship merge completes run"
assert_json_field "$merge_output" '"merged_via": "git"' "ship merge reports merge transport"

cleanup_output=$(run_in_repo ./scripts/harness-ship.sh cleanup "$run_id" --delete-branch)
assert_json_field "$cleanup_output" '"cleaned"' "ship cleanup reports cleaned artifacts"

section "RC"
rc_create_output=$(run_in_repo ./scripts/harness-rc.sh create v1.2.3 "$run_id")
assert_json_field "$rc_create_output" '"version": "v1.2.3"' "rc create records version"
assert_json_field "$rc_create_output" '"state": "active"' "rc create opens active rc"
rc_id=$(echo "$rc_create_output" | python3 -c 'import json,sys; print(json.load(sys.stdin)["rc_id"])')

rc_verify_output=$(run_in_repo ./scripts/harness-rc.sh verify "$rc_id" "staging-smoke:pass")
assert_json_field "$rc_verify_output" '"verification_status": "passed"' "rc verify records verification"

rc_promote_output=$(run_in_repo ./scripts/harness-rc.sh promote "$rc_id" "Promoted after staging smoke")
assert_json_field "$rc_promote_output" '"state": "promoted"' "rc promote archives promoted rc"

section "RC surgical fix and cancel"
rc_fix_create_output=$(run_in_repo ./scripts/harness-rc.sh create v1.2.4 "$run_id" --apply)
assert_json_field "$rc_fix_create_output" '"version": "v1.2.4"' "second rc create records version"
assert_json_field "$rc_fix_create_output" '"branch": "rc/v1.2.4"' "second rc create uses rc branch pattern"
rc_fix_id=$(echo "$rc_fix_create_output" | python3 -c 'import json,sys; print(json.load(sys.stdin)["rc_id"])')

(
  cd "$TEST_REPO" &&
    git checkout -q main &&
    printf "base fix for rc\n" >> README.md &&
    git add README.md &&
    git commit -qm "base fix for rc"
) || exit 1
fix_commit=$(cd "$TEST_REPO" && git rev-parse HEAD)

rc_fix_output=$(run_in_repo ./scripts/harness-rc.sh fix "$rc_fix_id" "$fix_commit" --apply --push)
assert_json_field "$rc_fix_output" '"applied": true' "rc fix applies cherry-pick onto rc branch"
assert_json_field "$rc_fix_output" '"pushed": true' "rc fix can push updated rc branch"

if git --git-dir="$ORIGIN_REPO" show-ref --verify --quiet refs/heads/rc/v1.2.4; then
  pass "rc fix push updates remote rc branch"
else
  fail "rc fix push updates remote rc branch"
fi

if (cd "$TEST_REPO" && git log rc/v1.2.4 --format=%s | grep -q "base fix for rc"); then
  pass "rc fix cherry-picks base fix onto rc branch"
else
  fail "rc fix cherry-picks base fix onto rc branch"
fi

rc_cancel_output=$(run_in_repo ./scripts/harness-rc.sh cancel "$rc_fix_id" "Cancel after staging defect" --apply)
assert_json_field "$rc_cancel_output" '"state": "closed"' "rc cancel archives the active rc"

if git --git-dir="$ORIGIN_REPO" show-ref --verify --quiet refs/heads/rc/v1.2.4; then
  fail "rc cancel deletes remote rc branch"
else
  pass "rc cancel deletes remote rc branch"
fi

section "Hotfix"
(
  cd "$TEST_REPO" &&
    git checkout -q -b candidate-hotfix main &&
    printf "candidate hotfix\n" >> README.md &&
    git add README.md &&
    git commit -qm "candidate hotfix"
) || exit 1

hotfix_output=$(run_in_repo ./scripts/harness-rc.sh hotfix v1.2.5 candidate-hotfix --apply --push)
assert_json_field "$hotfix_output" '"applied": true' "rc hotfix applies hotfix branch cherry-pick"
assert_json_field "$hotfix_output" '"backport_required": false' "rc hotfix reports no backport in single-branch mode"

if git --git-dir="$ORIGIN_REPO" show-ref --verify --quiet refs/heads/hotfix/v1.2.5; then
  pass "rc hotfix push updates remote hotfix branch"
else
  fail "rc hotfix push updates remote hotfix branch"
fi

section "Dual-branch backport"
(
  cd "$TEST_REPO" &&
    git checkout -q main &&
    git checkout -q -B develop main &&
    git push -u origin develop >/dev/null &&
    cat > harness/profiles/erwin.yml <<'EOF'
version: 1
profile:
  profile_id: "erwin-main"
  name: "erwin"
  repository: "liyasong/erwin"
  default_runtime: "codex"
  base_branch: "develop"
  production_branch: "main"
  branch_pattern: "mission/{issue}-{slug}"
  worktree_enabled_default: true
  issue_label_prefix: "status:"
  status_labels:
    queued: "status:queued"
    wip: "status:wip"
    review: "status:review"
    done: "status:done"
    blocked: "status:blocked"
  priority_labels:
    p0: "priority:P0"
    p1: "priority:P1"
    p2: "priority:P2"
  milestone_strategy: "github-milestones"
  default_verify_commands: []
  targeted_verify_conventions:
    - "Prefer the narrowest useful check first."
  required_pre_ship_checks:
    - "Relevant targeted verification has passed."
  release_mode: "dual-branch"
  rc_branch_pattern: "rc/{version}"
  hotfix_branch_pattern: "hotfix/{version}"
  backport_branch_pattern: "backport/{version}"
  staging_workflow: ""
  promote_workflow: ""
  docs_entrypoints:
    - "README.md"
  code_entrypoints:
    - "scripts/"
  config_entrypoints:
    - ".gitignore"
  legacy_paths: []
  max_parallel_runs: 2
  meeting_triggers:
    - "Release-affecting uncertainty"
  high_risk_scopes:
    - "release and RC flows"
  cleanup_policies:
    - "Report stale runs."
EOF
    git add harness/profiles/erwin.yml &&
    git commit -qm "switch to dual-branch release profile"
) || exit 1

(
  cd "$TEST_REPO" &&
    git checkout -q -b dual-backport-source develop &&
    printf "dual backport candidate\n" >> README.md &&
    git add README.md &&
    git commit -qm "dual backport candidate"
) || exit 1

dual_commit_ref="dual-backport-source"

dual_hotfix_output=$(run_in_repo ./scripts/harness-rc.sh hotfix v2.0.0 "$dual_commit_ref" --apply --push)
assert_json_field "$dual_hotfix_output" '"backport_required": true' "rc hotfix reports backport required in dual-branch mode"

if git --git-dir="$ORIGIN_REPO" show-ref --verify --quiet refs/heads/hotfix/v2.0.0; then
  pass "dual-branch hotfix push updates remote hotfix branch"
else
  fail "dual-branch hotfix push updates remote hotfix branch"
fi

backport_output=$(run_in_repo ./scripts/harness-rc.sh backport v2.0.0 hotfix/v2.0.0 --apply --push)
assert_json_field "$backport_output" '"applied": true' "rc backport applies cherry-pick onto backport branch"
assert_json_field "$backport_output" '"target_branch": "develop"' "rc backport targets develop in dual-branch mode"

if git --git-dir="$ORIGIN_REPO" show-ref --verify --quiet refs/heads/backport/v2.0.0; then
  pass "rc backport push updates remote backport branch"
else
  fail "rc backport push updates remote backport branch"
fi

if (cd "$TEST_REPO" && git show backport/v2.0.0:README.md | grep -q "dual backport candidate"); then
  pass "rc backport carries fix content onto backport branch"
else
  fail "rc backport carries fix content onto backport branch"
fi

section "Home and doctor"
home_output=$(run_in_repo ./scripts/harness-home.sh)
doctor_output=$(run_in_repo ./scripts/harness-doctor.sh)
adapter_output=$(run_in_repo ./scripts/harness-adapter.sh validate)
assert_json_field "$home_output" '"active_rc_count": 0' "home clears active rc after promotion"
assert_json_field "$adapter_output" '"valid": true' "adapter manifest validates in release smoke repo"
assert_json_field "$doctor_output" '"issue_count": 0' "doctor stays healthy after release flow"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
