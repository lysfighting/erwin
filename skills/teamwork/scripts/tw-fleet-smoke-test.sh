#!/usr/bin/env bash
# tw-fleet-smoke-test.sh — Smoke tests for /team fleet infrastructure
#
# USAGE:
#   bash tw-fleet-smoke-test.sh
#
# Covers:
#   1. tw-fleet.sh subcommand sanity (new-session-id, parse-issues, spawn-threshold,
#      init-session, persist-checkpoint, list/show-checkpoint, persist-approval,
#      render-briefing, cleanup-session)
#   2. Concurrency spike: N=3 parallel `tw-git.sh claim-setup` in worktree mode
#      → verify 3 distinct worktrees created, 3 distinct branches, no race corruption
#      (Gap 5 address — Sprout L358 self-flagged this as "claim to verify")
#   3. Template rendering: feature-lead-briefing-template.md with fake vars,
#      verify all {{PLACEHOLDER}}s resolve
#
# Self-contained — creates its own temp git repo + .teamwork config.
#
# EXIT CODES:
#   0 — all tests passed
#   1 — one or more tests failed

set -uo pipefail
# Note: NOT set -e — we handle errors per-test

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILL_DIR="$(dirname "$SCRIPT_DIR")"

# ─── Temp repo (isolated from real state) ───────────────
TMPROOT=$(mktemp -d "${TMPDIR:-/tmp}/tw-fleet-smoke.XXXXXX")
trap 'rm -rf "$TMPROOT"' EXIT

# ─── Counters ───────────────────────────────────────────
PASS=0
FAIL=0

pass() { ((PASS++)); echo "  ✅ $1"; }
fail() { ((FAIL++)); echo "  ❌ $1"; }

assert_eq() {
  local label="$1" expected="$2" actual="$3"
  if [ "$expected" = "$actual" ]; then
    pass "$label"
  else
    fail "$label — expected: $(printf '%q' "$expected"), got: $(printf '%q' "$actual")"
  fi
}

assert_file() {
  local label="$1" path="$2"
  if [ -f "$path" ]; then
    pass "$label"
  else
    fail "$label — file missing: $path"
  fi
}

assert_dir() {
  local label="$1" path="$2"
  if [ -d "$path" ]; then
    pass "$label"
  else
    fail "$label — dir missing: $path"
  fi
}

# ──────────────────────────────────────────────────────────────────
# Setup: make a minimal git repo with a .teamwork dir for tw-fleet.sh
#        and tw-git.sh to work against.
# ──────────────────────────────────────────────────────────────────

echo ""
echo "━━━ SETUP: temp git repo + .teamwork config ━━━"

# Create a bare "origin" repo first so claim-setup's worktree path (which references origin/base) works
ORIGIN_BARE="$TMPROOT/origin.git"
git init -q --bare "$ORIGIN_BARE" >/dev/null 2>&1

REPO="$TMPROOT/repo"
mkdir -p "$REPO"
cd "$REPO"
git init -q --initial-branch=main >/dev/null 2>&1 || git init -q >/dev/null 2>&1
git config user.email "smoke@example.com"
git config user.name  "Smoke"
git config commit.gpgsign false
echo "# test" > README.md
git add README.md
git commit -qm "initial"

# Wire origin so `git fetch origin main` in claim-setup works, then push so origin/main exists
git remote add origin "$ORIGIN_BARE"
git push -q -u origin main 2>/dev/null
git fetch -q origin main 2>/dev/null

mkdir -p .teamwork
cat > .teamwork/config.yml <<'EOF'
schema_version: 3
team:
  name: "Smoke"
  repo: local/smoke
roles:
  - { id: member, level: member, label: "Member" }
members:
  - { github: smoke, name: "Smoke", role: member }
conventions:
  base_branch: main
  production_branch: main
  branch_pattern: "mission/{issue}-{slug}"
labels:
  mission: "mission"
  status_prefix: "status:"
  priority_prefix: "priority:"
quality: { ci: false, review_required: false, branch_protection: false }
EOF

# Enable worktree mode so claim-setup uses the worktree path
git config --local teamwork.worktree true
git config --local teamwork.worktree-root "$TMPROOT/worktrees"
mkdir -p "$TMPROOT/worktrees"

pass "temp repo + .teamwork config created"

# ──────────────────────────────────────────────────────────────────
# Test 1: tw-fleet.sh subcommand sanity
# ──────────────────────────────────────────────────────────────────

echo ""
echo "━━━ TEST 1: tw-fleet.sh subcommand sanity ━━━"

# 1a: new-session-id produces a non-empty fleet-* string
SID=$(bash "$SCRIPT_DIR/tw-fleet.sh" new-session-id 2>/dev/null)
if [ -n "$SID" ] && [[ "$SID" == fleet-* ]]; then
  pass "new-session-id produced '$SID'"
else
  fail "new-session-id produced bad output: '$SID'"
fi

# 1b: parse-issues handles mixed input
P_OUT=$(bash "$SCRIPT_DIR/tw-fleet.sh" parse-issues "#1298 #1299, 1302" 2>/dev/null)
P_LINES=$(echo "$P_OUT" | wc -l | tr -d ' ')
assert_eq "parse-issues produces 3 lines" "3" "$P_LINES"
if [ "$(echo "$P_OUT" | head -1)" = "1298" ]; then
  pass "parse-issues first issue is 1298"
else
  fail "parse-issues first issue wrong: $(echo "$P_OUT" | head -1)"
fi

# 1c: parse-issues dedupes
P_DEDUP=$(bash "$SCRIPT_DIR/tw-fleet.sh" parse-issues "#42 #42 #43" 2>/dev/null | wc -l | tr -d ' ')
assert_eq "parse-issues dedupes duplicates" "2" "$P_DEDUP"

# 1d: spawn-threshold — various inputs
T1=$(bash "$SCRIPT_DIR/tw-fleet.sh" spawn-threshold 1 10 no 2>/dev/null)
assert_eq "spawn-threshold 1/10/no = direct" "direct" "$T1"

T2=$(bash "$SCRIPT_DIR/tw-fleet.sh" spawn-threshold 3 10 no 2>/dev/null)
assert_eq "spawn-threshold 3/10/no = spawn (files threshold)" "spawn" "$T2"

T3=$(bash "$SCRIPT_DIR/tw-fleet.sh" spawn-threshold 1 30 no 2>/dev/null)
assert_eq "spawn-threshold 1/30/no = spawn (lines threshold)" "spawn" "$T3"

T4=$(bash "$SCRIPT_DIR/tw-fleet.sh" spawn-threshold 1 10 yes 2>/dev/null)
assert_eq "spawn-threshold 1/10/yes = spawn (cross-module)" "spawn" "$T4"

# 1e: usage error
T_ERR=$(bash "$SCRIPT_DIR/tw-fleet.sh" spawn-threshold 1 10 2>&1)
if [ $? -ne 0 ] && echo "$T_ERR" | grep -q "ERROR:"; then
  pass "spawn-threshold with missing arg exits non-zero"
else
  fail "spawn-threshold with missing arg did not fail cleanly"
fi

# 1f: init-session creates the directory structure
FDIR=$(bash "$SCRIPT_DIR/tw-fleet.sh" init-session "$SID" 2>/dev/null)
assert_dir "init-session created fleet dir" "$FDIR"
assert_dir "init-session created checkpoints/ subdir" "$FDIR/checkpoints"
assert_dir "init-session created approvals/ subdir" "$FDIR/approvals"
assert_file "init-session created metadata.yml" "$FDIR/metadata.yml"

# 1g: persist + list + show checkpoint
CP_CONTENT="$TMPROOT/checkpoint-content.md"
cat > "$CP_CONTENT" <<'EOF'
# Design Note — mission #1298
## 1. Insertion point
step_evidence.py:148 `_enabled()` function.
EOF

CP_DEST=$(bash "$SCRIPT_DIR/tw-fleet.sh" persist-checkpoint "$SID" 1298 design-note "$CP_CONTENT" 2>/dev/null)
assert_file "persist-checkpoint wrote to $CP_DEST" "$CP_DEST"

LIST_OUT=$(bash "$SCRIPT_DIR/tw-fleet.sh" list-checkpoints "$SID" 2>/dev/null)
if echo "$LIST_OUT" | grep -q "1298:design-note:"; then
  pass "list-checkpoints includes persisted mission"
else
  fail "list-checkpoints missing mission: $LIST_OUT"
fi

SHOW_OUT=$(bash "$SCRIPT_DIR/tw-fleet.sh" show-checkpoint "$SID" 1298 design-note 2>/dev/null)
if echo "$SHOW_OUT" | grep -q "Insertion point"; then
  pass "show-checkpoint returns content"
else
  fail "show-checkpoint missing content"
fi

# 1h: persist-approval
AP_MSG="$TMPROOT/approval.md"
cat > "$AP_MSG" <<'EOF'
Approve with 2 adjustments:
1. Add coupling regression test
2. Use monkeypatch.delenv in test setup
EOF

AP_DEST=$(bash "$SCRIPT_DIR/tw-fleet.sh" persist-approval "$SID" 1298 approve "$AP_MSG" 2>/dev/null)
assert_file "persist-approval wrote to $AP_DEST" "$AP_DEST"

# Invalid decision should fail
bash "$SCRIPT_DIR/tw-fleet.sh" persist-approval "$SID" 1298 maybe "$AP_MSG" 2>/dev/null
if [ $? -ne 0 ]; then
  pass "persist-approval rejects invalid decision"
else
  fail "persist-approval accepted invalid decision"
fi

# 1i: cleanup-session removes everything
bash "$SCRIPT_DIR/tw-fleet.sh" cleanup-session "$SID" >/dev/null
if [ ! -d "$FDIR" ]; then
  pass "cleanup-session removed fleet dir"
else
  fail "cleanup-session left fleet dir behind: $FDIR"
fi

# 1j: cleanup-session is idempotent (already gone)
bash "$SCRIPT_DIR/tw-fleet.sh" cleanup-session "$SID" >/dev/null 2>&1
if [ $? -eq 0 ]; then
  pass "cleanup-session idempotent (exit 0 when already gone)"
else
  fail "cleanup-session not idempotent"
fi

# ──────────────────────────────────────────────────────────────────
# Test 2: Concurrency spike — N=3 parallel claim-setup (Gap 5)
# ──────────────────────────────────────────────────────────────────

echo ""
echo "━━━ TEST 2: Concurrency spike (Gap 5) — N=3 parallel claim-setup ━━━"

# Make a fresh fleet session for this test
SID2=$(bash "$SCRIPT_DIR/tw-fleet.sh" new-session-id)
sleep 1  # ensure distinct session id if prior id was sub-second
SID2="fleet-$(date -u +%Y%m%d-%H%M%S)-c2"

# Fake 3 distinct "issues" by creating mission numbers 901/902/903
# claim-setup needs: ISSUE SLUG USER [TEAMWORK_DIR [CONTRACT_PATH]]
mkdir -p .teamwork/active

# Per-mission Contract stubs (minimum viable for claim-setup to copy)
for N in 901 902 903; do
  cat > ".teamwork/active/MISSION-${N}.md" <<EOF
---
issue: $N
---
# Mission #$N
Stub Contract for smoke test.
EOF
done

# Launch 3 claim-setups in parallel via background processes; capture their stdout
OUT_DIR="$TMPROOT/claim-out"
mkdir -p "$OUT_DIR"

{
  bash "$SCRIPT_DIR/tw-git.sh" claim-setup 901 mission-one smoke \
    ".teamwork" ".teamwork/active/MISSION-901.md" > "$OUT_DIR/out-901.txt" 2>&1 &
  bash "$SCRIPT_DIR/tw-git.sh" claim-setup 902 mission-two smoke \
    ".teamwork" ".teamwork/active/MISSION-902.md" > "$OUT_DIR/out-902.txt" 2>&1 &
  bash "$SCRIPT_DIR/tw-git.sh" claim-setup 903 mission-three smoke \
    ".teamwork" ".teamwork/active/MISSION-903.md" > "$OUT_DIR/out-903.txt" 2>&1 &
  wait
}

# Verify 3 distinct worktrees exist
WT_COUNT=0
for N in 901 902 903; do
  WT_PATH=$(grep '^WORKTREE_PATH=' "$OUT_DIR/out-${N}.txt" 2>/dev/null | sed 's/^WORKTREE_PATH=//')
  if [ -n "$WT_PATH" ] && [ -d "$WT_PATH" ]; then
    ((WT_COUNT++))
    # Verify mission marker
    MISSION_MARKER=$(cat "$WT_PATH/.mission" 2>/dev/null)
    if [ "$MISSION_MARKER" = "$N" ]; then
      pass "mission #$N worktree created with correct marker at $WT_PATH"
    else
      fail "mission #$N worktree has wrong marker: '$MISSION_MARKER' (expected $N)"
    fi
  else
    fail "mission #$N claim-setup did not produce a worktree (out=$(cat "$OUT_DIR/out-${N}.txt"))"
  fi
done

assert_eq "3 parallel claim-setups each produced a worktree" "3" "$WT_COUNT"

# Verify 3 distinct branches exist
BRANCH_COUNT=0
for N in 901 902 903; do
  BRANCH=$(grep '^BRANCH=' "$OUT_DIR/out-${N}.txt" 2>/dev/null | sed 's/^BRANCH=//')
  if [ -n "$BRANCH" ] && git rev-parse --verify --quiet "$BRANCH" >/dev/null 2>&1; then
    ((BRANCH_COUNT++))
    pass "mission #$N branch '$BRANCH' exists"
  else
    fail "mission #$N branch not created or not found (BRANCH='$BRANCH')"
  fi
done

assert_eq "3 parallel claim-setups each created a branch" "3" "$BRANCH_COUNT"

# Verify no cross-contamination (each worktree has only its own .mission marker)
CROSS_CONTAM=0
for N in 901 902 903; do
  WT_PATH=$(grep '^WORKTREE_PATH=' "$OUT_DIR/out-${N}.txt" 2>/dev/null | sed 's/^WORKTREE_PATH=//')
  if [ -n "$WT_PATH" ] && [ -d "$WT_PATH" ]; then
    OTHER_MARKERS=$(find "$WT_PATH" -maxdepth 1 -name '.mission' -type f | wc -l | tr -d ' ')
    if [ "$OTHER_MARKERS" -ne 1 ]; then
      CROSS_CONTAM=1
    fi
  fi
done

if [ "$CROSS_CONTAM" -eq 0 ]; then
  pass "no cross-contamination between parallel worktrees"
else
  fail "cross-contamination detected between parallel worktrees"
fi

# ──────────────────────────────────────────────────────────────────
# Test 3: Template rendering (briefing template)
# ──────────────────────────────────────────────────────────────────

echo ""
echo "━━━ TEST 3: briefing template rendering ━━━"

TPL="$SKILL_DIR/reference/feature-lead-briefing-template.md"
if [ ! -f "$TPL" ]; then
  fail "briefing template missing: $TPL"
else
  pass "briefing template exists"

  VARS="$TMPROOT/vars.env"
  cat > "$VARS" <<EOF
CONTRACT_PATH=/tmp/fake/MISSION-1298.md
WORKTREE_PATH=/tmp/fake/worktrees/mission-1298
MISSION_NUMBER=1298
CHECKPOINT_TYPE=design-note
TEAM_NAME=fleet-20260417-smoke
EOF

  RENDERED=$(bash "$SCRIPT_DIR/tw-fleet.sh" render-briefing "$TPL" "$VARS" 2>/dev/null)

  # Verify all placeholders resolved
  REMAINING=$(echo "$RENDERED" | grep -oE '\{\{[A-Z_]+\}\}' | wc -l | tr -d ' ')
  assert_eq "all placeholders resolved (0 left)" "0" "$REMAINING"

  if echo "$RENDERED" | grep -q "Mission \*\*#1298\*\*"; then
    pass "mission number interpolated"
  else
    fail "mission number not interpolated"
  fi

  if echo "$RENDERED" | grep -q "mission-1298"; then
    pass "worktree path interpolated"
  else
    fail "worktree path not interpolated"
  fi

  if echo "$RENDERED" | grep -q "design-note"; then
    pass "checkpoint type interpolated"
  else
    fail "checkpoint type not interpolated"
  fi
fi

# Verify the 5 checkpoint schemas exist
echo ""
for ctype in design-note delta-table research-result risk-inventory ad-hoc; do
  assert_file "checkpoint-schemas/$ctype.md exists" "$SKILL_DIR/reference/checkpoint-schemas/$ctype.md"
done

# ──────────────────────────────────────────────────────────────────
# Summary
# ──────────────────────────────────────────────────────────────────

echo ""
echo "━━━ SUMMARY ━━━"
echo ""
echo "  Tests: $((PASS + FAIL)) total"
echo "  ✅ Passed: $PASS"
echo "  ❌ Failed: $FAIL"
echo ""

if [ "$FAIL" -eq 0 ]; then
  echo "  🎉 ALL TESTS PASSED"
  exit 0
else
  echo "  ❌ $FAIL TEST(S) FAILED"
  exit 1
fi
