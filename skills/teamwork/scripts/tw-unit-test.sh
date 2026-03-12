#!/usr/bin/env bash
# tw-unit-test.sh — Unit tests for teamwork scripts (no git repo or GitHub auth needed)
#
# USAGE:
#   bash tw-unit-test.sh
#
# Tests pure-function logic in tw-bump-version.sh, tw-config.sh, tw-git.sh
# using temporary files. Self-contained — no external dependencies.
#
# EXIT CODES:
#   0 — all tests passed
#   1 — one or more tests failed

set -uo pipefail
# Note: NOT set -e — we handle errors per-test

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ─── Temp directory ─────────────────────────────────────
TMPDIR_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/tw-unit-test.XXXXXX")
trap 'rm -rf "$TMPDIR_ROOT"' EXIT

# ─── Counters ───────────────────────────────────────────
PASS=0
FAIL=0
SKIP=0

# ─── Helpers ────────────────────────────────────────────
pass() { ((PASS++)); echo "  ✅ $1"; }
fail() { ((FAIL++)); echo "  ❌ $1"; }
skip() { ((SKIP++)); echo "  ⏭️  $1"; }
info() { echo "  ℹ️  $1"; }
section() { echo ""; echo "━━━ $1 ━━━"; }

run_script() {
  local script="$SCRIPT_DIR/$1"
  shift
  bash "$script" "$@"
}

# ═══════════════════════════════════════════════════════
# 1. tw-bump-version.sh
# ═══════════════════════════════════════════════════════

test_bump_rejects_no_args() {
  local out exit_code=0
  out=$(run_script tw-bump-version.sh 2>&1) || exit_code=$?
  if [ "$exit_code" -ne 0 ]; then
    pass "bump: rejects no args (exit $exit_code)"
  else
    fail "bump: should reject no args"
  fi
}

test_bump_rejects_invalid_formats() {
  local cases=("3" "3.4" "3.4.5.6" "abc" "3.4.x" "v3.4.5" "3.4.5-rc1" "3..5")
  for v in "${cases[@]}"; do
    local exit_code=0
    run_script tw-bump-version.sh "$v" >/dev/null 2>&1 || exit_code=$?
    if [ "$exit_code" -ne 0 ]; then
      pass "bump: rejects '$v'"
    else
      fail "bump: should reject '$v'"
    fi
  done
}

test_bump_rejects_flags() {
  local exit_code=0
  run_script tw-bump-version.sh "--help" >/dev/null 2>&1 || exit_code=$?
  if [ "$exit_code" -ne 0 ]; then
    pass "bump: rejects flag '--help' as version"
  else
    fail "bump: should reject '--help'"
  fi
}

test_bump_updates_frontmatter_version() {
  # Create a fake skill directory structure
  local skill_dir="$TMPDIR_ROOT/bump-frontmatter/teamwork"
  local scripts_dir="$skill_dir/scripts"
  mkdir -p "$scripts_dir"

  cat > "$skill_dir/team.md" << 'EOF'
---
name: team
version: "1.0.0"
---
# Team skill
EOF

  cat > "$skill_dir/team-ship.md" << 'EOF'
---
name: team-ship
version: "1.0.0"
---
# Ship skill
EOF

  # Create a symlink to the real bump script in our fake scripts dir
  cp "$SCRIPT_DIR/tw-bump-version.sh" "$scripts_dir/tw-bump-version.sh"

  # Run bump from the fake scripts dir
  local out
  out=$(bash "$scripts_dir/tw-bump-version.sh" "2.5.0" 2>&1) || true

  # Check team.md updated
  if grep -q 'version: "2.5.0"' "$skill_dir/team.md"; then
    pass "bump: frontmatter version updated in team.md"
  else
    fail "bump: frontmatter version NOT updated in team.md"
  fi

  # Check team-ship.md updated
  if grep -q 'version: "2.5.0"' "$skill_dir/team-ship.md"; then
    pass "bump: frontmatter version updated in team-ship.md"
  else
    fail "bump: frontmatter version NOT updated in team-ship.md"
  fi
}

test_bump_updates_version_inline() {
  local skill_dir="$TMPDIR_ROOT/bump-inline/teamwork"
  local scripts_dir="$skill_dir/scripts"
  mkdir -p "$scripts_dir"

  cat > "$skill_dir/team.md" << 'EOF'
---
name: team
version: "1.0.0"
---
# Team skill

**Version**: `1.0.0`

This is TEAMWORK v1.0.0 and also Teamwork v1.0.0.
Something (v1.0.0) here.
SKILL_VERSION="1.0.0"
Version 1.0.0 manual header.
Skill version: 1.0.0 readme line.
(`v1.0.0`) inline backtick-parens.
EOF

  cp "$SCRIPT_DIR/tw-bump-version.sh" "$scripts_dir/tw-bump-version.sh"

  bash "$scripts_dir/tw-bump-version.sh" "3.8.0" >/dev/null 2>&1 || true

  local content
  content=$(cat "$skill_dir/team.md")

  if echo "$content" | grep -qF '**Version**: `3.8.0`'; then
    pass "bump: **Version**: \`X.Y.Z\` updated"
  else
    fail "bump: **Version**: \`X.Y.Z\` NOT updated"
  fi

  if echo "$content" | grep -q 'TEAMWORK v3.8.0'; then
    pass "bump: TEAMWORK vX.Y.Z updated"
  else
    fail "bump: TEAMWORK vX.Y.Z NOT updated"
  fi

  if echo "$content" | grep -q 'Teamwork v3.8.0'; then
    pass "bump: Teamwork vX.Y.Z updated"
  else
    fail "bump: Teamwork vX.Y.Z NOT updated"
  fi

  if echo "$content" | grep -q '(v3.8.0)'; then
    pass "bump: (vX.Y.Z) updated"
  else
    fail "bump: (vX.Y.Z) NOT updated"
  fi

  if echo "$content" | grep -q 'SKILL_VERSION="3.8.0"'; then
    pass "bump: SKILL_VERSION=\"X.Y.Z\" updated"
  else
    fail "bump: SKILL_VERSION=\"X.Y.Z\" NOT updated"
  fi

  if echo "$content" | grep -q 'Version 3.8.0'; then
    pass "bump: Version X.Y.Z updated"
  else
    fail "bump: Version X.Y.Z NOT updated"
  fi

  if echo "$content" | grep -q 'Skill version: 3.8.0'; then
    pass "bump: Skill version: X.Y.Z updated"
  else
    fail "bump: Skill version: X.Y.Z NOT updated"
  fi

  if echo "$content" | grep -q '(`v3.8.0`)'; then
    pass "bump: (\`vX.Y.Z\`) updated"
  else
    fail "bump: (\`vX.Y.Z\`) NOT updated"
  fi
}

test_bump_updates_manual_file() {
  # The script looks for manual at: $(dirname "$SKILL_DIR")/../docs/teamwork-ai-manual.md
  # SKILL_DIR = dirname(SCRIPT_DIR), so:
  # SCRIPT_DIR = .../teamwork/scripts
  # SKILL_DIR  = .../teamwork
  # MANUAL     = .../docs/teamwork-ai-manual.md
  local base="$TMPDIR_ROOT/bump-manual"
  local skill_dir="$base/skills/teamwork"
  local scripts_dir="$skill_dir/scripts"
  local docs_dir="$base/docs"
  mkdir -p "$scripts_dir" "$docs_dir"

  # Minimal skill file (needed so script has files to iterate)
  cat > "$skill_dir/team.md" << 'EOF'
---
name: team
version: "1.0.0"
---
EOF

  # Manual file
  cat > "$docs_dir/teamwork-ai-manual.md" << 'EOF'
---
name: manual
version: "1.0.0"
---
# Teamwork AI Manual — Version 1.0.0
EOF

  cp "$SCRIPT_DIR/tw-bump-version.sh" "$scripts_dir/tw-bump-version.sh"

  bash "$scripts_dir/tw-bump-version.sh" "4.0.0" >/dev/null 2>&1 || true

  if grep -q 'version: "4.0.0"' "$docs_dir/teamwork-ai-manual.md"; then
    pass "bump: manual file frontmatter updated"
  else
    fail "bump: manual file frontmatter NOT updated"
  fi

  if grep -q 'Version 4.0.0' "$docs_dir/teamwork-ai-manual.md"; then
    pass "bump: manual Version header updated"
  else
    fail "bump: manual Version header NOT updated"
  fi
}

test_bump_idempotent() {
  local skill_dir="$TMPDIR_ROOT/bump-idempotent/teamwork"
  local scripts_dir="$skill_dir/scripts"
  mkdir -p "$scripts_dir"

  cat > "$skill_dir/team.md" << 'EOF'
---
name: team
version: "2.0.0"
---
EOF

  cp "$SCRIPT_DIR/tw-bump-version.sh" "$scripts_dir/tw-bump-version.sh"

  # Bump to same version — should not error
  local out exit_code=0
  out=$(bash "$scripts_dir/tw-bump-version.sh" "2.0.0" 2>&1) || exit_code=$?
  if [ "$exit_code" -eq 0 ]; then
    pass "bump: idempotent (same version) succeeds"
  else
    fail "bump: idempotent (same version) failed (exit $exit_code)"
  fi
}

# ═══════════════════════════════════════════════════════
# 2. tw-config.sh — migrate + has-section
# ═══════════════════════════════════════════════════════

# Helper: create a temp git repo with teamwork config
_setup_config_repo() {
  local repo_dir="$1"
  mkdir -p "$repo_dir/.teamwork"
  (cd "$repo_dir" && git init -q && git config user.email "test@test" && git config user.name "Test")
}

test_config_migrate_removes_worktree_section() {
  local repo="$TMPDIR_ROOT/config-migrate-wt"
  _setup_config_repo "$repo"

  cat > "$repo/.teamwork/config.yml" << 'EOF'
schema_version: 2
project:
  name: test-project
worktree:
  branch_pattern: "mission/{issue}-{slug}-{user}"
  base_dir: ".claude/worktrees"
conventions:
  base_branch: main
EOF

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-config.sh" migrate 2>&1) || true

  # Check worktree: section removed
  if ! grep -q "^worktree:" "$repo/.teamwork/config.yml"; then
    pass "migrate: worktree section removed"
  else
    fail "migrate: worktree section still present"
  fi

  # Check schema_version updated to 3
  if grep -q "^schema_version: 3$" "$repo/.teamwork/config.yml"; then
    pass "migrate: schema_version updated to 3"
  else
    fail "migrate: schema_version not updated to 3"
  fi

  # Check other sections preserved
  if grep -q "^conventions:" "$repo/.teamwork/config.yml"; then
    pass "migrate: conventions section preserved"
  else
    fail "migrate: conventions section lost"
  fi
}

test_config_migrate_removes_skill_version() {
  local repo="$TMPDIR_ROOT/config-migrate-sv"
  _setup_config_repo "$repo"

  cat > "$repo/.teamwork/config.yml" << 'EOF'
schema_version: 3
skill_version: "3.5.0"
project:
  name: test-project
EOF

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-config.sh" migrate 2>&1) || true

  if ! grep -q "^skill_version:" "$repo/.teamwork/config.yml"; then
    pass "migrate: skill_version line removed"
  else
    fail "migrate: skill_version line still present"
  fi
}

test_config_migrate_idempotent() {
  local repo="$TMPDIR_ROOT/config-migrate-idem"
  _setup_config_repo "$repo"

  cat > "$repo/.teamwork/config.yml" << 'EOF'
schema_version: 3
project:
  name: test-project
conventions:
  base_branch: main
EOF

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-config.sh" migrate 2>&1) || true

  if echo "$out" | grep -q "already up to date"; then
    pass "migrate: idempotent (no changes needed)"
  else
    fail "migrate: should report 'already up to date', got: $out"
  fi
}

test_config_migrate_adds_schema_version() {
  local repo="$TMPDIR_ROOT/config-migrate-add-sv"
  _setup_config_repo "$repo"

  cat > "$repo/.teamwork/config.yml" << 'EOF'
project:
  name: test-project
EOF

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-config.sh" migrate 2>&1) || true

  if grep -q "schema_version: 3" "$repo/.teamwork/config.yml"; then
    pass "migrate: schema_version added when missing"
  else
    fail "migrate: schema_version not added"
  fi
}

test_config_has_section_present() {
  local repo="$TMPDIR_ROOT/config-has-sec-present"
  _setup_config_repo "$repo"

  cat > "$repo/.teamwork/config.yml" << 'EOF'
schema_version: 3
conventions:
  base_branch: main
members:
  - github: alice
EOF

  local exit_code=0
  (cd "$repo" && bash "$SCRIPT_DIR/tw-config.sh" has-section conventions) || exit_code=$?
  if [ "$exit_code" -eq 0 ]; then
    pass "has-section: 'conventions' found (exit 0)"
  else
    fail "has-section: 'conventions' should be found"
  fi

  exit_code=0
  (cd "$repo" && bash "$SCRIPT_DIR/tw-config.sh" has-section members) || exit_code=$?
  if [ "$exit_code" -eq 0 ]; then
    pass "has-section: 'members' found (exit 0)"
  else
    fail "has-section: 'members' should be found"
  fi
}

test_config_has_section_absent() {
  local repo="$TMPDIR_ROOT/config-has-sec-absent"
  _setup_config_repo "$repo"

  cat > "$repo/.teamwork/config.yml" << 'EOF'
schema_version: 3
conventions:
  base_branch: main
EOF

  local exit_code=0
  (cd "$repo" && bash "$SCRIPT_DIR/tw-config.sh" has-section members 2>/dev/null) || exit_code=$?
  if [ "$exit_code" -ne 0 ]; then
    pass "has-section: 'members' absent (exit $exit_code)"
  else
    fail "has-section: 'members' should NOT be found"
  fi
}

test_config_has_section_label_prefix_dual_schema() {
  # label_prefix should also match 'labels:' (dual schema)
  local repo="$TMPDIR_ROOT/config-has-sec-lp"
  _setup_config_repo "$repo"

  cat > "$repo/.teamwork/config.yml" << 'EOF'
schema_version: 3
labels:
  mission: "mc"
  status_prefix: "status:"
EOF

  local exit_code=0
  (cd "$repo" && bash "$SCRIPT_DIR/tw-config.sh" has-section label_prefix) || exit_code=$?
  if [ "$exit_code" -eq 0 ]; then
    pass "has-section: 'label_prefix' matches 'labels:' (dual schema)"
  else
    fail "has-section: 'label_prefix' should match 'labels:' section"
  fi
}

test_config_has_section_no_args() {
  local repo="$TMPDIR_ROOT/config-has-sec-noargs"
  _setup_config_repo "$repo"

  cat > "$repo/.teamwork/config.yml" << 'EOF'
schema_version: 3
EOF

  local exit_code=0
  (cd "$repo" && bash "$SCRIPT_DIR/tw-config.sh" has-section 2>/dev/null) || exit_code=$?
  if [ "$exit_code" -ne 0 ]; then
    pass "has-section: no args rejected (exit $exit_code)"
  else
    fail "has-section: should reject missing args"
  fi
}

# ═══════════════════════════════════════════════════════
# 3. tw-git.sh — slugify + next-version
# ═══════════════════════════════════════════════════════

test_slugify_basic() {
  local out
  out=$(run_script tw-git.sh slugify "Fix Login Bug" 2>/dev/null)
  if [ "$out" = "fix-login-bug" ]; then
    pass "slugify: 'Fix Login Bug' → '$out'"
  else
    fail "slugify: expected 'fix-login-bug', got '$out'"
  fi
}

test_slugify_uppercase() {
  local out
  out=$(run_script tw-git.sh slugify "ADD AUTH MIDDLEWARE" 2>/dev/null)
  if [ "$out" = "add-auth-middleware" ]; then
    pass "slugify: uppercase → '$out'"
  else
    fail "slugify: expected 'add-auth-middleware', got '$out'"
  fi
}

test_slugify_special_chars() {
  local out
  out=$(run_script tw-git.sh slugify "fix: camera [iOS] + Safari" 2>/dev/null)
  if echo "$out" | grep -qE '^[a-z0-9-]+$'; then
    pass "slugify: special chars removed → '$out'"
  else
    fail "slugify: special chars not cleaned, got '$out'"
  fi
}

test_slugify_consecutive_hyphens() {
  local out
  out=$(run_script tw-git.sh slugify "fix -- double   spaces" 2>/dev/null)
  if ! echo "$out" | grep -q '\-\-'; then
    pass "slugify: no consecutive hyphens → '$out'"
  else
    fail "slugify: consecutive hyphens found in '$out'"
  fi
}

test_slugify_truncates_long_strings() {
  local out
  out=$(run_script tw-git.sh slugify "this is a very long title that should be truncated to thirty characters maximum" 2>/dev/null)
  local len=${#out}
  if [ "$len" -le 30 ]; then
    pass "slugify: truncated to $len chars → '$out'"
  else
    fail "slugify: too long ($len chars) → '$out'"
  fi
}

test_slugify_no_trailing_hyphen() {
  local out
  out=$(run_script tw-git.sh slugify "something ends with special!" 2>/dev/null)
  if [[ "$out" != *- ]]; then
    pass "slugify: no trailing hyphen → '$out'"
  else
    fail "slugify: trailing hyphen in '$out'"
  fi
}

test_slugify_no_leading_hyphen() {
  local out
  out=$(run_script tw-git.sh slugify "  leading spaces" 2>/dev/null)
  if [[ "$out" != -* ]]; then
    pass "slugify: no leading hyphen → '$out'"
  else
    fail "slugify: leading hyphen in '$out'"
  fi
}

test_slugify_empty_rejects() {
  local exit_code=0
  run_script tw-git.sh slugify "" >/dev/null 2>&1 || exit_code=$?
  if [ "$exit_code" -ne 0 ]; then
    pass "slugify: empty string rejected (exit $exit_code)"
  else
    fail "slugify: should reject empty string"
  fi
}

test_slugify_all_special_rejects() {
  local exit_code=0
  run_script tw-git.sh slugify "!!!@@@" >/dev/null 2>&1 || exit_code=$?
  if [ "$exit_code" -ne 0 ]; then
    pass "slugify: all-special-chars rejected (exit $exit_code)"
  else
    fail "slugify: should reject string that produces empty slug"
  fi
}

test_slugify_numbers() {
  local out
  out=$(run_script tw-git.sh slugify "issue 42 fix" 2>/dev/null)
  if [ "$out" = "issue-42-fix" ]; then
    pass "slugify: numbers preserved → '$out'"
  else
    fail "slugify: expected 'issue-42-fix', got '$out'"
  fi
}

test_next_version_no_tags() {
  # Create a temp git repo with no tags
  local repo="$TMPDIR_ROOT/next-ver-no-tags"
  mkdir -p "$repo/.teamwork"
  (cd "$repo" && git init -q && git config user.email "test@test" && git config user.name "Test" && echo "init" > README.md && git add -A && git commit -q -m "init")

  # Write a minimal config so tw-config.sh works
  cat > "$repo/.teamwork/config.yml" << 'EOF'
schema_version: 3
conventions:
  base_branch: main
EOF

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-git.sh" next-version "V0.1" 2>/dev/null) || true
  if [ "$out" = "V0.1.0" ]; then
    pass "next-version: no tags → V0.1.0"
  else
    fail "next-version: expected V0.1.0, got '$out'"
  fi
}

test_next_version_with_tags() {
  local repo="$TMPDIR_ROOT/next-ver-with-tags"
  mkdir -p "$repo/.teamwork"
  (cd "$repo" && git init -q && git config user.email "test@test" && git config user.name "Test" && echo "init" > README.md && git add -A && git commit -q -m "init")

  cat > "$repo/.teamwork/config.yml" << 'EOF'
schema_version: 3
conventions:
  base_branch: main
EOF

  # Create some tags
  (cd "$repo" && git tag "V0.1.0" && git tag "V0.1.1" && git tag "V0.1.2")

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-git.sh" next-version "V0.1" 2>/dev/null) || true
  if [ "$out" = "V0.1.3" ]; then
    pass "next-version: V0.1.{0,1,2} → V0.1.3"
  else
    fail "next-version: expected V0.1.3, got '$out'"
  fi
}

test_next_version_ignores_prerelease() {
  local repo="$TMPDIR_ROOT/next-ver-prerelease"
  mkdir -p "$repo/.teamwork"
  (cd "$repo" && git init -q && git config user.email "test@test" && git config user.name "Test" && echo "init" > README.md && git add -A && git commit -q -m "init")

  cat > "$repo/.teamwork/config.yml" << 'EOF'
schema_version: 3
conventions:
  base_branch: main
EOF

  # Create tags including a pre-release suffix
  (cd "$repo" && git tag "V0.2.0" && git tag "V0.2.1-rc1" && git tag "V0.2.1")

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-git.sh" next-version "V0.2" 2>/dev/null) || true
  if [ "$out" = "V0.2.2" ]; then
    pass "next-version: ignores prerelease → V0.2.2"
  else
    fail "next-version: expected V0.2.2, got '$out'"
  fi
}

test_next_version_no_args() {
  local exit_code=0
  run_script tw-git.sh next-version 2>/dev/null || exit_code=$?
  if [ "$exit_code" -ne 0 ]; then
    pass "next-version: no args rejected (exit $exit_code)"
  else
    fail "next-version: should reject missing args"
  fi
}

# ═══════════════════════════════════════════════════════
# RUN ALL TESTS
# ═══════════════════════════════════════════════════════

section "BUMP VERSION — VALIDATION"
test_bump_rejects_no_args
test_bump_rejects_invalid_formats
test_bump_rejects_flags

section "BUMP VERSION — FRONTMATTER"
test_bump_updates_frontmatter_version

section "BUMP VERSION — INLINE PATTERNS"
test_bump_updates_version_inline

section "BUMP VERSION — MANUAL FILE"
test_bump_updates_manual_file

section "BUMP VERSION — IDEMPOTENT"
test_bump_idempotent

section "CONFIG — MIGRATE"
test_config_migrate_removes_worktree_section
test_config_migrate_removes_skill_version
test_config_migrate_idempotent
test_config_migrate_adds_schema_version

section "CONFIG — HAS-SECTION"
test_config_has_section_present
test_config_has_section_absent
test_config_has_section_label_prefix_dual_schema
test_config_has_section_no_args

section "GIT — SLUGIFY"
test_slugify_basic
test_slugify_uppercase
test_slugify_special_chars
test_slugify_consecutive_hyphens
test_slugify_truncates_long_strings
test_slugify_no_trailing_hyphen
test_slugify_no_leading_hyphen
test_slugify_empty_rejects
test_slugify_all_special_rejects
test_slugify_numbers

section "GIT — NEXT-VERSION"
test_next_version_no_tags
test_next_version_with_tags
test_next_version_ignores_prerelease
test_next_version_no_args

# ─── Summary ────────────────────────────────────────────
section "SUMMARY"
TOTAL=$((PASS + FAIL + SKIP))
echo ""
echo "  Tests: $TOTAL total"
echo "  ✅ Passed: $PASS"
echo "  ❌ Failed: $FAIL"
echo "  ⏭️  Skipped: $SKIP"
echo ""

if [ "$FAIL" -eq 0 ]; then
  echo "  🎉 ALL TESTS PASSED"
  exit 0
else
  echo "  💥 $FAIL TEST(S) FAILED"
  exit 1
fi
