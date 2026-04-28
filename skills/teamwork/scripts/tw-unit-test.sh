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

section "GIT — WORKTREE-GUARD"

test_worktree_guard_no_mismatch_when_disabled() {
  local repo="$TMPDIR_ROOT/wt-guard-disabled"
  mkdir -p "$repo" && cd "$repo"
  git init -q && git commit --allow-empty -m "init" -q
  # worktree disabled (default) — should succeed
  local out exit_code=0
  out=$(run_script tw-git.sh worktree-guard 2>&1) || exit_code=$?
  if [ "$exit_code" -eq 0 ] && echo "$out" | grep -q "MISMATCH=false"; then
    pass "worktree-guard: no mismatch when disabled"
  else
    fail "worktree-guard: expected MISMATCH=false, got exit=$exit_code out=$out"
  fi
  cd "$TMPDIR_ROOT"
}

test_worktree_guard_mismatch_when_enabled_no_worktree() {
  local repo="$TMPDIR_ROOT/wt-guard-enabled"
  mkdir -p "$repo" && cd "$repo"
  git init -q && git commit --allow-empty -m "init" -q
  git config --local teamwork.worktree true
  # No .mission file → should be mismatch (exit 3)
  local out exit_code=0
  out=$(run_script tw-git.sh worktree-guard "some-branch" 2>&1) || exit_code=$?
  if [ "$exit_code" -eq 3 ] && echo "$out" | grep -q "MISMATCH=true"; then
    pass "worktree-guard: mismatch when enabled but no .mission"
  else
    fail "worktree-guard: expected exit=3 MISMATCH=true, got exit=$exit_code out=$out"
  fi
  cd "$TMPDIR_ROOT"
}

test_worktree_guard_ok_when_in_worktree() {
  local repo="$TMPDIR_ROOT/wt-guard-in-wt"
  mkdir -p "$repo" && cd "$repo"
  git init -q && git commit --allow-empty -m "init" -q
  git config --local teamwork.worktree true
  echo "42" > .mission
  # .mission exists → IN_WORKTREE=true → no mismatch
  local out exit_code=0
  out=$(run_script tw-git.sh worktree-guard 2>&1) || exit_code=$?
  if [ "$exit_code" -eq 0 ] && echo "$out" | grep -q "MISMATCH=false" && echo "$out" | grep -q "IN_WORKTREE=true"; then
    pass "worktree-guard: ok when .mission exists"
  else
    fail "worktree-guard: expected ok with IN_WORKTREE=true, got exit=$exit_code out=$out"
  fi
  rm -f .mission
  cd "$TMPDIR_ROOT"
}

test_worktree_guard_no_mismatch_when_disabled
test_worktree_guard_mismatch_when_enabled_no_worktree
test_worktree_guard_ok_when_in_worktree

section "GIT — CLAIM-SETUP"

test_claim_setup_rejects_missing_args() {
  local repo="$TMPDIR_ROOT/claim-setup-args"
  mkdir -p "$repo" && cd "$repo"
  git init -q && git commit --allow-empty -m "init" -q
  local out exit_code=0
  out=$(run_script tw-git.sh claim-setup 2>&1) || exit_code=$?
  if [ "$exit_code" -ne 0 ]; then
    pass "claim-setup: rejects missing args (exit $exit_code)"
  else
    fail "claim-setup: should reject missing args"
  fi
  cd "$TMPDIR_ROOT"
}

test_claim_setup_branch_mode() {
  local repo="$TMPDIR_ROOT/claim-setup-branch"
  mkdir -p "$repo" && cd "$repo"
  git init -q && git commit --allow-empty -m "init" -q
  # Need a "main" branch for ensure-base
  git branch -M main
  # No worktree config → should create branch
  # Create config for branch pattern
  mkdir -p .teamwork
  cat > .teamwork/config.yml <<'YAML'
conventions:
  base_branch: main
  branch_pattern: "mission/{issue}-{slug}"
YAML
  local out exit_code=0
  out=$(run_script tw-git.sh claim-setup "42" "test-slug" "user1" 2>&1) || exit_code=$?
  if [ "$exit_code" -eq 0 ] && echo "$out" | grep -q "MODE=branch" && echo "$out" | grep -q "STATUS=created"; then
    pass "claim-setup: branch mode creates branch"
  else
    fail "claim-setup: branch mode failed, exit=$exit_code out=$out"
  fi
  # Verify branch was created
  local current
  current=$(git branch --show-current 2>/dev/null)
  if echo "$current" | grep -q "mission/42-test-slug"; then
    pass "claim-setup: branch name correct ($current)"
  else
    fail "claim-setup: wrong branch ($current)"
  fi
  cd "$TMPDIR_ROOT"
}

test_claim_setup_rejects_missing_args
test_claim_setup_branch_mode

section "GIT — NEXT-VERSION"
test_next_version_no_tags
test_next_version_with_tags
test_next_version_ignores_prerelease
test_next_version_no_args

# ─── PR — MERGED-SINCE renderers ──────────────────────────────
test_merged_since_rejects_no_base() {
  local out exit_code
  out=$(bash "$SCRIPT_DIR/tw-pr.sh" merged-since 2>&1) && exit_code=0 || exit_code=$?
  if [ "$exit_code" -eq 1 ] && echo "$out" | grep -q "requires BASE"; then
    pass "merged-since: no BASE → exit 1"
  else
    fail "merged-since: expected exit 1 'requires BASE', got exit $exit_code: $out"
  fi
}

test_merged_since_rejects_bad_format() {
  local out exit_code
  out=$(bash "$SCRIPT_DIR/tw-pr.sh" merged-since pre-launch --format bogus --repo x/y 2>&1) && exit_code=0 || exit_code=$?
  if [ "$exit_code" -eq 1 ] && echo "$out" | grep -q "format must be"; then
    pass "merged-since: bad --format → exit 1"
  else
    fail "merged-since: bad format expected exit 1, got $exit_code: $out"
  fi
}

test_merged_since_markdown_renderer() {
  # Install a tiny stub `gh` on PATH that returns fixture JSON for `gh pr list ...`.
  local stubdir="$TMPDIR_ROOT/stub-$$"
  mkdir -p "$stubdir"
  cat > "$stubdir/gh" <<'EOF'
#!/usr/bin/env bash
if [ "$1" = "pr" ] && [ "$2" = "list" ]; then
  cat <<'JSON'
[
  {"number": 100, "title": "feat(foo): add bar", "author":{"login":"alice"}, "mergedAt":"2030-01-02T10:00:00Z", "url":"https://example/100"},
  {"number": 101, "title": "fix(baz): patch qux", "author":{"login":"bob"}, "mergedAt":"2030-01-03T10:00:00Z", "url":"https://example/101"},
  {"number": 102, "title": "docs: update README", "author":{"login":"alice"}, "mergedAt":"2030-01-04T10:00:00Z", "url":"https://example/102"},
  {"number": 99,  "title": "feat: old feature", "author":{"login":"carol"}, "mergedAt":"2029-12-01T10:00:00Z", "url":"https://example/99"}
]
JSON
  exit 0
elif [ "$1" = "repo" ] && [ "$2" = "view" ]; then
  echo "stub/repo"
  exit 0
fi
echo "stub-gh: unknown invocation: $*" >&2; exit 1
EOF
  chmod +x "$stubdir/gh"

  # Prev-tag argument is empty → no git log; PREV_DATE defaults to 2000 → all 4 PRs pass filter.
  local out
  out=$(PATH="$stubdir:$PATH" bash "$SCRIPT_DIR/tw-pr.sh" merged-since pre-launch "" --repo stub/repo --format markdown --version "V0.1.20" 2>&1)

  local ok=1
  echo "$out" | grep -q "^## What's in this release" || { fail "markdown: missing header"; ok=0; }
  echo "$out" | grep -q "^### Features" || { fail "markdown: missing Features"; ok=0; }
  echo "$out" | grep -q "^### Bug Fixes" || { fail "markdown: missing Bug Fixes"; ok=0; }
  echo "$out" | grep -q "^### Other Changes" || { fail "markdown: missing Other Changes"; ok=0; }
  echo "$out" | grep -q "feat(foo): add bar.*@alice" || { fail "markdown: feat line wrong: $(echo "$out" | grep feat)"; ok=0; }
  echo "$out" | grep -q "fix(baz): patch qux.*@bob" || { fail "markdown: fix line wrong"; ok=0; }
  echo "$out" | grep -q "docs: update README.*@alice" || { fail "markdown: other line wrong"; ok=0; }
  echo "$out" | grep -q "@alice (2)" || { fail "markdown: contributor count wrong"; ok=0; }

  [ "$ok" -eq 1 ] && pass "merged-since --format markdown renders categories + contributors"
}

test_merged_since_feishu_renderer() {
  local stubdir="$TMPDIR_ROOT/stub-f-$$"
  mkdir -p "$stubdir"
  cat > "$stubdir/gh" <<'EOF'
#!/usr/bin/env bash
if [ "$1" = "pr" ] && [ "$2" = "list" ]; then
  cat <<'JSON'
[
  {"number": 200, "title": "feat(x): a", "author":{"login":"alice"}, "mergedAt":"2030-01-01T00:00:00Z", "url":"https://example/200"},
  {"number": 201, "title": "fix(y): b", "author":{"login":"bob"}, "mergedAt":"2030-01-02T00:00:00Z", "url":"https://example/201"}
]
JSON
  exit 0
fi
exit 1
EOF
  chmod +x "$stubdir/gh"

  local out
  out=$(PATH="$stubdir:$PATH" bash "$SCRIPT_DIR/tw-pr.sh" merged-since pre-launch "" --repo stub/repo --format feishu --version "V0.1.20" --env staging 2>&1)

  local ok=1
  echo "$out" | grep -q "^# 🚀 V0.1.20 Staging Release" || { fail "feishu: missing header/emoji"; ok=0; }
  echo "$out" | grep -q "^\*\*环境\*\*: \`staging\`" || { fail "feishu: missing 环境 line"; ok=0; }
  echo "$out" | grep -q "^\*\*时间\*\*:" || { fail "feishu: missing 时间 line"; ok=0; }
  echo "$out" | grep -q "^\*\*分支\*\*: \`rc/V0.1.20\` ← \`pre-launch\`" || { fail "feishu: missing 分支 line"; ok=0; }
  echo "$out" | grep -q "^### ✨ Features" || { fail "feishu: missing Features emoji"; ok=0; }
  echo "$out" | grep -q "^### 🐛 Bug Fixes" || { fail "feishu: missing Bug Fixes emoji"; ok=0; }
  echo "$out" | grep -q "^### 👥 Contributors" || { fail "feishu: missing Contributors"; ok=0; }
  echo "$out" | grep -q "合计 \*\*2\*\* PR" || { fail "feishu: missing total count"; ok=0; }

  [ "$ok" -eq 1 ] && pass "merged-since --format feishu renders header + metadata + categories + total"
}

test_merged_since_feishu_production_emoji() {
  local stubdir="$TMPDIR_ROOT/stub-p-$$"
  mkdir -p "$stubdir"
  cat > "$stubdir/gh" <<'EOF'
#!/usr/bin/env bash
if [ "$1" = "pr" ] && [ "$2" = "list" ]; then
  echo '[{"number":300,"title":"feat: x","author":{"login":"a"},"mergedAt":"2030-01-01T00:00:00Z","url":"u"}]'
  exit 0
fi
exit 1
EOF
  chmod +x "$stubdir/gh"

  local out
  out=$(PATH="$stubdir:$PATH" bash "$SCRIPT_DIR/tw-pr.sh" merged-since main "" --repo x/y --format feishu --version "V1.0.0" --env production --release-url "https://github.com/x/y/releases/tag/V1.0.0" --pr-number 999 2>&1)

  local ok=1
  echo "$out" | grep -q "^# 🎉 V1.0.0 Production Release" || { fail "prod feishu: missing 🎉 header"; ok=0; }
  echo "$out" | grep -q "Release PR.*999" || { fail "prod feishu: missing PR number"; ok=0; }
  echo "$out" | grep -q "GitHub Release.*releases/tag/V1.0.0" || { fail "prod feishu: missing release URL"; ok=0; }
  [ "$ok" -eq 1 ] && pass "merged-since --format feishu production path renders release URL + PR number"
}

section "PR — MERGED-SINCE"
test_merged_since_rejects_no_base
test_merged_since_rejects_bad_format
test_merged_since_markdown_renderer
test_merged_since_feishu_renderer
test_merged_since_feishu_production_emoji

# ─── GIT — HOTFIX-DIVERGENCE-CHECK ──────────────────────
# Builds tiny git repos with a "prod" and "rc" branch, applies the four
# canonical divergence scenarios, asserts STATUS field + presence of
# expected entries. No GitHub auth, no remote — runs as `origin` of itself
# via a bare repo.

_setup_hotfix_repo() {
  # Create: bare "origin" + working clone with main + rc/V0.1.99 branches,
  # and a clean V0.1.27 release tag on main. Returns repo path on stdout.
  local repo="$TMPDIR_ROOT/hotfix-$$-$RANDOM"
  local bare="$repo.git"
  rm -rf "$repo" "$bare"

  git init --quiet --bare "$bare"
  git init --quiet "$repo"
  (
    cd "$repo"
    git config user.email "test@test.test"
    git config user.name "test"
    git config commit.gpgsign false
    git remote add origin "$bare"

    # Initial commit on main.
    echo "v0" > app.txt
    git add app.txt
    git commit --quiet -m "init"
    git branch -M main
    git push --quiet -u origin main

    # Tag a clean release V0.1.27 on main.
    git tag V0.1.27
    git push --quiet origin V0.1.27

    # Cut rc/V0.1.99 from main.
    git checkout --quiet -b rc/V0.1.99
    git push --quiet -u origin rc/V0.1.99

    git checkout --quiet main
  )
  echo "$repo"
}

test_hotfix_divergence_clean_path() {
  # Scenario 4: prod at last_clean_tag, no hotfixes. STATUS=clean.
  local repo
  repo=$(_setup_hotfix_repo)
  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-git.sh" hotfix-divergence-check rc/V0.1.99 main V0.1 2>&1)

  local ok=1
  echo "$out" | grep -q "^LAST_CLEAN_TAG=V0.1.27$"   || { fail "hotfix-divergence: clean — wrong LAST_CLEAN_TAG"; ok=0; }
  echo "$out" | grep -q "^HOTFIX_TAGS=$"              || { fail "hotfix-divergence: clean — HOTFIX_TAGS should be empty"; ok=0; }
  echo "$out" | grep -q "^MISSING_HOTFIXES=$"         || { fail "hotfix-divergence: clean — MISSING_HOTFIXES should be empty"; ok=0; }
  echo "$out" | grep -q "^STATUS=clean$"              || { fail "hotfix-divergence: clean — wrong STATUS"; ok=0; }

  [ "$ok" -eq 1 ] && pass "hotfix-divergence: clean path (no hotfix tags) → STATUS=clean"
}

test_hotfix_divergence_no_backport() {
  # Scenario 3: prod has hotfix tag, RC never got it. STATUS=missing.
  local repo
  repo=$(_setup_hotfix_repo)
  (
    cd "$repo"
    # Apply a hotfix on main: change app.txt and tag.
    sleep 1  # ensure commit time > V0.1.27
    echo "v0-prod-hotfix" > app.txt
    git add app.txt
    git commit --quiet -m "hotfix(app): patch on prod"
    git tag V0.1.27-hotfix.1
    git push --quiet origin main V0.1.27-hotfix.1
  )

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-git.sh" hotfix-divergence-check rc/V0.1.99 main V0.1 2>&1)

  local ok=1
  echo "$out" | grep -q "^STATUS=missing$"                                  || { fail "hotfix-divergence: no-backport — wrong STATUS"; ok=0; }
  echo "$out" | grep -q "^HOTFIX_TAGS<<EOF$"                                || { fail "hotfix-divergence: no-backport — HOTFIX_TAGS should be heredoc"; ok=0; }
  echo "$out" | grep -q "V0.1.27-hotfix.1 .* hotfix(app): patch on prod"    || { fail "hotfix-divergence: no-backport — expected tag entry not found"; ok=0; }
  echo "$out" | grep -q "^MISSING_HOTFIXES<<EOF$"                           || { fail "hotfix-divergence: no-backport — MISSING_HOTFIXES should be heredoc"; ok=0; }

  [ "$ok" -eq 1 ] && pass "hotfix-divergence: hotfix on prod, no backport on rc → STATUS=missing"
}

test_hotfix_divergence_manual_backport_different_content() {
  # Scenario 1: prod has hotfix.1 with patch A, RC has a manual fix with patch B (different content). STATUS=missing.
  local repo
  repo=$(_setup_hotfix_repo)
  (
    cd "$repo"
    sleep 1  # commit time after V0.1.27
    # Prod hotfix:
    echo "v0-prod-version" > app.txt
    git add app.txt
    git commit --quiet -m "hotfix(app): prod-side fix"
    git tag V0.1.27-hotfix.1
    git push --quiet origin main V0.1.27-hotfix.1

    # RC manual backport with DIFFERENT content (different SHA, different patch):
    git checkout --quiet rc/V0.1.99
    echo "v0-rc-version" > app.txt
    git add app.txt
    git commit --quiet -m "fix(app): rc-side manual backport"
    git push --quiet origin rc/V0.1.99
    git checkout --quiet main
  )

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-git.sh" hotfix-divergence-check rc/V0.1.99 main V0.1 2>&1)

  local ok=1
  echo "$out" | grep -q "^STATUS=missing$"                              || { fail "hotfix-divergence: manual-backport-diff — wrong STATUS"; ok=0; }
  echo "$out" | grep -q "V0.1.27-hotfix.1 .* hotfix(app): prod-side fix" || { fail "hotfix-divergence: manual-backport-diff — expected tag entry not found"; ok=0; }

  [ "$ok" -eq 1 ] && pass "hotfix-divergence: manual backport with different content → STATUS=missing (correct alarm)"
}

test_hotfix_divergence_cherry_pick_backport() {
  # Scenario 2: prod has hotfix.1, RC `cherry-pick -x` the SAME commit. SHAs differ.
  # Expected: STATUS=missing (conservative false positive — guard fires; user picks Force).
  local repo
  repo=$(_setup_hotfix_repo)
  (
    cd "$repo"
    sleep 1
    echo "v0-shared-fix" > app.txt
    git add app.txt
    git commit --quiet -m "hotfix(app): shared fix"
    local hotfix_sha
    hotfix_sha=$(git rev-parse HEAD)
    git tag V0.1.27-hotfix.1
    git push --quiet origin main V0.1.27-hotfix.1

    # RC cherry-picks the same commit; -x records source SHA in trailer.
    git checkout --quiet rc/V0.1.99
    git cherry-pick -x "$hotfix_sha" --quiet >/dev/null 2>&1
    git push --quiet origin rc/V0.1.99
    git checkout --quiet main
  )

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-git.sh" hotfix-divergence-check rc/V0.1.99 main V0.1 2>&1)

  local ok=1
  # Cherry-pick of same commit: original SHA still not in RC ancestry → guard fires.
  # This is the documented conservative behavior.
  echo "$out" | grep -q "^STATUS=missing$"                            || { fail "hotfix-divergence: cherry-pick — wrong STATUS (expected conservative missing)"; ok=0; }
  echo "$out" | grep -q "V0.1.27-hotfix.1 .* hotfix(app): shared fix"  || { fail "hotfix-divergence: cherry-pick — expected tag entry not found"; ok=0; }

  [ "$ok" -eq 1 ] && pass "hotfix-divergence: cherry-pick -x backport → STATUS=missing (conservative; user resolves via Force)"
}

test_hotfix_divergence_no_prior_release() {
  # No clean release tag exists yet (only bootstrap commit). STATUS=no-prior-release.
  local repo="$TMPDIR_ROOT/hotfix-npr-$$-$RANDOM"
  local bare="$repo.git"
  rm -rf "$repo" "$bare"
  git init --quiet --bare "$bare"
  git init --quiet "$repo"
  (
    cd "$repo"
    git config user.email "test@test.test"
    git config user.name "test"
    git config commit.gpgsign false
    git remote add origin "$bare"
    echo "v0" > app.txt
    git add app.txt
    git commit --quiet -m "init"
    git branch -M main
    git push --quiet -u origin main
  )

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-git.sh" hotfix-divergence-check main main V0.1 2>&1)

  local ok=1
  echo "$out" | grep -q "^LAST_CLEAN_TAG=$"        || { fail "hotfix-divergence: no-prior-release — LAST_CLEAN_TAG should be empty"; ok=0; }
  echo "$out" | grep -q "^STATUS=no-prior-release$" || { fail "hotfix-divergence: no-prior-release — wrong STATUS"; ok=0; }
  echo "$out" | grep -q "^HOTFIX_TAGS=$"            || { fail "hotfix-divergence: no-prior-release — HOTFIX_TAGS should be empty"; ok=0; }
  echo "$out" | grep -q "^MISSING_HOTFIXES=$"       || { fail "hotfix-divergence: no-prior-release — MISSING_HOTFIXES should be empty"; ok=0; }

  [ "$ok" -eq 1 ] && pass "hotfix-divergence: no clean release tag → STATUS=no-prior-release"
}

test_hotfix_divergence_fallback_grep() {
  # Untagged hotfix commit on prod (didn't go through /team-rc hotfix flow).
  # Expected: HOTFIX_TAGS empty, HOTFIX_COMMITS populated, MISSING_HOTFIXES has (untagged) entry.
  local repo
  repo=$(_setup_hotfix_repo)
  (
    cd "$repo"
    sleep 1
    # Manual hotfix commit on prod — NOT tagged.
    echo "v0-untagged-hotfix" > app.txt
    git add app.txt
    git commit --quiet -m "hotfix(app): manual untagged fix"
    git push --quiet origin main
  )

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-git.sh" hotfix-divergence-check rc/V0.1.99 main V0.1 2>&1)

  local ok=1
  echo "$out" | grep -q "^STATUS=missing$"                              || { fail "hotfix-divergence: fallback — wrong STATUS"; ok=0; }
  echo "$out" | grep -q "^HOTFIX_TAGS=$"                                || { fail "hotfix-divergence: fallback — HOTFIX_TAGS should be empty (no tag created)"; ok=0; }
  echo "$out" | grep -q "^HOTFIX_COMMITS<<EOF$"                         || { fail "hotfix-divergence: fallback — HOTFIX_COMMITS should be heredoc"; ok=0; }
  echo "$out" | grep -q "hotfix(app): manual untagged fix"              || { fail "hotfix-divergence: fallback — expected fallback subject not found"; ok=0; }
  echo "$out" | grep -q "^(untagged) [0-9a-f]\{40\} hotfix(app):"       || { fail "hotfix-divergence: fallback — MISSING_HOTFIXES should contain '(untagged) SHA hotfix(...)' entry"; ok=0; }

  [ "$ok" -eq 1 ] && pass "hotfix-divergence: untagged 'hotfix(...)' commit → grep fallback fires, MISSING includes (untagged) entry"
}

test_hotfix_divergence_old_hotfix_filtered() {
  # Hotfix tag PREDATING last_clean_tag should be filtered by time check
  # (already incorporated via subsequent release). STATUS=clean.
  local repo="$TMPDIR_ROOT/hotfix-old-$$-$RANDOM"
  local bare="$repo.git"
  rm -rf "$repo" "$bare"
  git init --quiet --bare "$bare"
  git init --quiet "$repo"
  (
    cd "$repo"
    git config user.email "test@test.test"
    git config user.name "test"
    git config commit.gpgsign false
    git remote add origin "$bare"

    echo "v0" > app.txt
    git add app.txt
    git commit --quiet -m "init"
    git branch -M main
    git push --quiet -u origin main

    # Older hotfix tag (predates the current cycle's release).
    git tag V0.1.20-hotfix.1
    git push --quiet origin V0.1.20-hotfix.1
    sleep 1

    # Subsequent clean release that "incorporated" the hotfix.
    echo "v0-superseded" > app.txt
    git add app.txt
    git commit --quiet -m "release: V0.1.27"
    git tag V0.1.27
    git push --quiet origin main V0.1.27

    # Cut RC from current main — does NOT contain V0.1.20-hotfix.1's SHA.
    git checkout --quiet -b rc/V0.1.99
    git push --quiet -u origin rc/V0.1.99
    git checkout --quiet main
  )

  local out
  out=$(cd "$repo" && bash "$SCRIPT_DIR/tw-git.sh" hotfix-divergence-check rc/V0.1.99 main V0.1 2>&1)

  local ok=1
  echo "$out" | grep -q "^LAST_CLEAN_TAG=V0.1.27$"   || { fail "hotfix-divergence: old-hotfix-filter — wrong LAST_CLEAN_TAG"; ok=0; }
  echo "$out" | grep -q "^STATUS=clean$"              || { fail "hotfix-divergence: old-hotfix-filter — old hotfix should be filtered out (STATUS should be clean, got: $(echo "$out" | grep '^STATUS='))"; ok=0; }
  echo "$out" | grep -q "^HOTFIX_TAGS=$"              || { fail "hotfix-divergence: old-hotfix-filter — HOTFIX_TAGS should be empty (V0.1.20-hotfix.1 predates V0.1.27)"; ok=0; }

  [ "$ok" -eq 1 ] && pass "hotfix-divergence: hotfix tag predating last_clean_tag → time-filtered out"
}

section "GIT — HOTFIX-DIVERGENCE-CHECK"
test_hotfix_divergence_clean_path
test_hotfix_divergence_no_backport
test_hotfix_divergence_manual_backport_different_content
test_hotfix_divergence_cherry_pick_backport
test_hotfix_divergence_no_prior_release
test_hotfix_divergence_fallback_grep
test_hotfix_divergence_old_hotfix_filtered

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
