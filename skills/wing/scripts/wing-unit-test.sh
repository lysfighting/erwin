#!/usr/bin/env bash
# wing-unit-test.sh — Comprehensive unit tests for wing skills
#
# USAGE:
#   bash wing-unit-test.sh              # All offline tests
#   bash wing-unit-test.sh --integration # Include GitHub API tests
#
# Tests: scripts, SKILL.md content accuracy, operations, install.sh
# Self-contained — uses temp dirs, no side effects.

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WING_DIR="$(dirname "$SCRIPT_DIR")"
ERWIN_ROOT="$(dirname "$(dirname "$(dirname "$SCRIPT_DIR")")")"
INTEGRATION=false
[[ "${1:-}" == "--integration" ]] && INTEGRATION=true

# ─── Temp directory ─────────────────────────────────────
TMPDIR_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/wing-unit-test.XXXXXX")
trap 'rm -rf "$TMPDIR_ROOT"' EXIT

# ─── Counters ───────────────────────────────────────────
PASS=0
FAIL=0
SKIP=0

# ─── Helpers ────────────────────────────────────────────
pass() { PASS=$((PASS + 1)); echo "  ✅ $1"; }
fail() { FAIL=$((FAIL + 1)); echo "  ❌ $1"; }
skip() { SKIP=$((SKIP + 1)); echo "  ⏭️  $1"; }
info() { echo "  ℹ️  $1"; }
section() { echo ""; echo "━━━ $1 ━━━"; }

run_script() {
  local script="$SCRIPT_DIR/$1"
  shift
  bash "$script" "$@"
}

# Helper: create a temp git repo with .wing config
_setup_wing_repo() {
  local repo_dir="$1"
  mkdir -p "$repo_dir/.wing"
  (cd "$repo_dir" && git init -q && git config user.email "test@test" && git config user.name "Test")
  cat > "$repo_dir/.wing/config.yml" << 'YAML'
schema_version: 3
team:
  - github: testuser
    role: tech-lead
roles:
  - id: tech-lead
    label: "Tech Lead"
    level: leader
  - id: developer
    label: "Developer"
    level: member
project:
  language: typescript
  test_command: "echo test-pass"
conventions:
  base_branch: "main"
  production_branch: "main"
  branch_pattern: "mission/{issue}-{slug}-{user}"
label_prefix:
  status: "status:"
  priority: "priority:"
  mission: "mission"
YAML
}

# ═══════════════════════════════════════════════════════
# 1. wing-config.sh
# ═══════════════════════════════════════════════════════

section "wing-config.sh"

test_config_detect_wing_dir() {
  local repo="$TMPDIR_ROOT/config-detect-wing"
  _setup_wing_repo "$repo"

  local result
  result=$(cd "$repo" && run_script wing-config.sh detect-dir 2>&1) || true
  if [[ "$result" == *"/.wing" ]]; then
    pass "config: detect-dir finds .wing/"
  else
    fail "config: detect-dir should find .wing/ (got: $result)"
  fi
}

test_config_detect_teamwork_fallback() {
  local repo="$TMPDIR_ROOT/config-detect-fallback"
  mkdir -p "$repo/.teamwork"
  (cd "$repo" && git init -q && git config user.email "t@t" && git config user.name "T")
  cat > "$repo/.teamwork/config.yml" << 'YAML'
schema_version: 3
team:
  - github: testuser
    role: tech-lead
YAML

  local result
  result=$(cd "$repo" && run_script wing-config.sh detect-dir 2>&1) || true
  if [[ "$result" == *"/.teamwork" ]]; then
    pass "config: detect-dir falls back to .teamwork/"
  else
    fail "config: detect-dir should fall back to .teamwork/ (got: $result)"
  fi
}

test_config_detect_no_config() {
  local repo="$TMPDIR_ROOT/config-detect-none"
  mkdir -p "$repo"
  (cd "$repo" && git init -q && git config user.email "t@t" && git config user.name "T")

  local exit_code=0
  cd "$repo" && run_script wing-config.sh detect-dir 2>/dev/null || exit_code=$?
  if [ "$exit_code" -eq 3 ]; then
    pass "config: detect-dir exits 3 when no config"
  else
    fail "config: detect-dir should exit 3 when no config (got: $exit_code)"
  fi
}

test_config_detect_priority_order() {
  # Both .wing/ and .teamwork/ exist — .wing/ should win
  local repo="$TMPDIR_ROOT/config-detect-priority"
  _setup_wing_repo "$repo"
  mkdir -p "$repo/.teamwork"
  cat > "$repo/.teamwork/config.yml" << 'YAML'
schema_version: 3
team:
  - github: olduser
YAML

  local result
  result=$(cd "$repo" && run_script wing-config.sh detect-dir 2>&1) || true
  if [[ "$result" == *"/.wing" ]]; then
    pass "config: .wing/ takes priority over .teamwork/"
  else
    fail "config: .wing/ should take priority (got: $result)"
  fi
}

test_config_resolve_labels_defaults() {
  local repo="$TMPDIR_ROOT/config-labels"
  _setup_wing_repo "$repo"

  local result
  result=$(cd "$repo" && run_script wing-config.sh resolve-labels 2>&1) || true
  if echo "$result" | grep -q 'MISSION_LABEL="mission"' &&
     echo "$result" | grep -q 'STATUS_PREFIX="status:"' &&
     echo "$result" | grep -q 'PRIORITY_PREFIX="priority:"'; then
    pass "config: resolve-labels returns correct defaults"
  else
    fail "config: resolve-labels wrong output: $result"
  fi
}

test_config_team_mode_solo() {
  local repo="$TMPDIR_ROOT/config-solo"
  _setup_wing_repo "$repo"

  local result
  result=$(cd "$repo" && run_script wing-config.sh team-mode 2>&1)
  if [[ "$result" == "solo" ]]; then
    pass "config: team-mode=solo for 1 member"
  else
    fail "config: team-mode should be solo (got: $result)"
  fi
}

test_config_team_mode_team() {
  local repo="$TMPDIR_ROOT/config-team"
  _setup_wing_repo "$repo"
  cat > "$repo/.wing/config.yml" << 'YAML'
schema_version: 3
team:
  - github: user1
    role: tech-lead
  - github: user2
    role: developer
YAML

  local result
  result=$(cd "$repo" && run_script wing-config.sh team-mode 2>&1)
  if [[ "$result" == "team" ]]; then
    pass "config: team-mode=team for 2 members"
  else
    fail "config: team-mode should be team (got: $result)"
  fi
}

test_config_get_user_level() {
  local repo="$TMPDIR_ROOT/config-level"
  _setup_wing_repo "$repo"

  local result
  result=$(cd "$repo" && run_script wing-config.sh get-user-level "testuser" 2>&1)
  if [[ "$result" == "solo" ]] || [[ "$result" == "leader" ]]; then
    pass "config: get-user-level for known user ($result)"
  else
    fail "config: get-user-level should return solo or leader (got: $result)"
  fi

  # In solo mode (1 member), get-user-level returns "solo" for ALL users
  # This is by design — solo mode doesn't gate on identity
  result=$(cd "$repo" && run_script wing-config.sh get-user-level "unknownuser" 2>&1) || true
  if [[ "$result" == "solo" ]]; then
    pass "config: get-user-level returns 'solo' for any user in solo mode"
  else
    fail "config: get-user-level should return 'solo' in solo mode (got: $result)"
  fi
}

test_config_key_lookup() {
  local repo="$TMPDIR_ROOT/config-key"
  _setup_wing_repo "$repo"

  local result
  result=$(cd "$repo" && run_script wing-config.sh project.language "none" 2>&1)
  if [[ "$result" == "typescript" ]]; then
    pass "config: key lookup project.language = typescript"
  else
    fail "config: project.language should be typescript (got: $result)"
  fi

  result=$(cd "$repo" && run_script wing-config.sh nonexistent.key "DEFAULT" 2>&1)
  if [[ "$result" == "DEFAULT" ]]; then
    pass "config: key lookup returns default for missing key"
  else
    fail "config: missing key should return DEFAULT (got: $result)"
  fi
}

test_config_resolve_member() {
  local repo="$TMPDIR_ROOT/config-member"
  _setup_wing_repo "$repo"
  cat > "$repo/.wing/config.yml" << 'YAML'
schema_version: 3
team:
  - github: lysfighting
    name: YaSong
    role: tech-lead
  - github: coworker
    name: Alice
    role: developer
YAML

  local result exit_code=0
  result=$(cd "$repo" && run_script wing-config.sh resolve-member "lysfighting" 2>&1) || exit_code=$?
  if [[ "$result" == "lysfighting" ]]; then
    pass "config: resolve-member exact match"
  else
    fail "config: resolve-member exact should return lysfighting (got: $result, exit $exit_code)"
  fi

  result=$(cd "$repo" && run_script wing-config.sh resolve-member "YaSong" 2>&1) || exit_code=$?
  if [[ "$result" == "lysfighting" ]]; then
    pass "config: resolve-member by name"
  else
    fail "config: resolve-member by name should return lysfighting (got: $result)"
  fi

  exit_code=0
  result=$(cd "$repo" && run_script wing-config.sh resolve-member "nobody" 2>/dev/null) || exit_code=$?
  if [[ "$exit_code" -ne 0 ]]; then
    pass "config: resolve-member fails for unknown"
  else
    fail "config: resolve-member should fail for unknown (got: $result)"
  fi
}

test_config_migrate() {
  local repo="$TMPDIR_ROOT/config-migrate"
  _setup_wing_repo "$repo"
  cat > "$repo/.wing/config.yml" << 'YAML'
schema_version: 2
skill_version: "3.5.0"
worktree:
  branch_pattern: "mission/{issue}-{slug}-{user}"
  base_dir: ".claude/worktrees"
conventions:
  base_branch: main
YAML

  local out
  out=$(cd "$repo" && run_script wing-config.sh migrate 2>&1) || true

  if ! grep -q "^worktree:" "$repo/.wing/config.yml"; then
    pass "config: migrate removes worktree section"
  else
    fail "config: migrate should remove worktree section"
  fi

  if grep -q "^schema_version: 3$" "$repo/.wing/config.yml"; then
    pass "config: migrate updates schema_version to 3"
  else
    fail "config: migrate should update schema_version"
  fi

  if ! grep -q "^skill_version:" "$repo/.wing/config.yml"; then
    pass "config: migrate removes skill_version"
  else
    fail "config: migrate should remove skill_version"
  fi

  if grep -q "^conventions:" "$repo/.wing/config.yml"; then
    pass "config: migrate preserves other sections"
  else
    fail "config: migrate should preserve other sections"
  fi
}

test_config_detect_wing_dir
test_config_detect_teamwork_fallback
test_config_detect_no_config
test_config_detect_priority_order
test_config_resolve_labels_defaults
test_config_team_mode_solo
test_config_team_mode_team
test_config_get_user_level
test_config_key_lookup
test_config_resolve_member
test_config_migrate

# ═══════════════════════════════════════════════════════
# 2. wing-git.sh
# ═══════════════════════════════════════════════════════

section "wing-git.sh"

test_git_slugify() {
  # Use | as separator to avoid : collision with input
  local cases=(
    "Fix camera permission|fix-camera-permission"
    "Add Rate Limiting to API|add-rate-limiting-to-api"
    "fix: camera on Safari!!!|fix-camera-on-safari"
    "hello---world|hello-world"
    "  leading spaces |leading-spaces"
    "UPPERCASE STRING|uppercase-string"
    "a/b/c|abc"
    "123 numbers first|123-numbers-first"
  )
  for case in "${cases[@]}"; do
    local input="${case%%|*}"
    local expected="${case##*|}"
    local result
    result=$(run_script wing-git.sh slugify "$input" 2>&1) || true
    # Slugify truncates to 30 chars
    if [[ "$result" == "${expected:0:30}"* ]] || [[ "$result" == "${expected}"* ]]; then
      pass "git: slugify '$input' → $result"
    else
      fail "git: slugify '$input' expected '$expected' got '$result'"
    fi
  done
}

test_git_protect_check() {
  local repo="$TMPDIR_ROOT/git-protect"
  _setup_wing_repo "$repo"
  (cd "$repo" && git commit --allow-empty -m "init" -q)

  local exit_code=0
  (cd "$repo" && run_script wing-git.sh protect-check 2>/dev/null) || exit_code=$?
  if [ "$exit_code" -eq 3 ]; then
    pass "git: protect-check exits 3 on main"
  else
    fail "git: protect-check should exit 3 on main (got: $exit_code)"
  fi

  # Create and switch to feature branch
  (cd "$repo" && git checkout -b mission/42-test -q)
  exit_code=0
  (cd "$repo" && run_script wing-git.sh protect-check 2>/dev/null) || exit_code=$?
  if [ "$exit_code" -eq 0 ]; then
    pass "git: protect-check passes on feature branch"
  else
    fail "git: protect-check should pass on feature branch (got: $exit_code)"
  fi
}

test_git_current() {
  local repo="$TMPDIR_ROOT/git-current"
  _setup_wing_repo "$repo"
  (cd "$repo" && git commit --allow-empty -m "init" -q)

  local result
  result=$(cd "$repo" && run_script wing-git.sh current 2>&1)
  if [[ "$result" == "main" ]]; then
    pass "git: current returns main"
  else
    fail "git: current should return main (got: $result)"
  fi
}

test_git_extract_issue_number() {
  local cases=(
    "mission/42-fix-camera:42"
    "mission/123-add-auth-user1:123"
    "bugfix/T-044-fix-login:044"
  )
  for case in "${cases[@]}"; do
    local branch="${case%%:*}"
    local expected="${case##*:}"
    local result
    result=$(run_script wing-git.sh extract-issue-number "$branch" 2>&1) || true
    if [[ "$result" == "$expected" ]]; then
      pass "git: extract-issue-number '$branch' → $result"
    else
      fail "git: extract-issue-number '$branch' expected '$expected' got '$result'"
    fi
  done
}

test_git_build_branch_name() {
  local repo="$TMPDIR_ROOT/git-branch-name"
  _setup_wing_repo "$repo"
  (cd "$repo" && git commit --allow-empty -m "init" -q)

  local result
  result=$(cd "$repo" && run_script wing-git.sh build-branch-name "42" "fix-camera" "testuser" 2>&1) || true
  if [[ "$result" == "mission/42-fix-camera-testuser" ]]; then
    pass "git: build-branch-name correct pattern"
  else
    fail "git: build-branch-name expected mission/42-fix-camera-testuser (got: $result)"
  fi
}

test_git_next_version() {
  local repo="$TMPDIR_ROOT/git-next-ver"
  _setup_wing_repo "$repo"
  (cd "$repo" && git commit --allow-empty -m "init" -q)

  # No tags → should return prefix.0
  local result
  result=$(cd "$repo" && run_script wing-git.sh next-version "V0.1" 2>&1) || true
  if [[ "$result" == "V0.1.0" ]]; then
    pass "git: next-version V0.1 with no tags → V0.1.0"
  else
    fail "git: next-version should be V0.1.0 (got: $result)"
  fi

  # Add a tag and check increment
  (cd "$repo" && git tag "V0.1.0")
  result=$(cd "$repo" && run_script wing-git.sh next-version "V0.1" 2>&1) || true
  if [[ "$result" == "V0.1.1" ]]; then
    pass "git: next-version V0.1 with V0.1.0 tag → V0.1.1"
  else
    fail "git: next-version should be V0.1.1 (got: $result)"
  fi
}

test_git_worktree_guard_disabled() {
  local repo="$TMPDIR_ROOT/git-wt-guard"
  _setup_wing_repo "$repo"
  (cd "$repo" && git commit --allow-empty -m "init" -q)

  # Worktree disabled (default) → should pass
  local exit_code=0
  (cd "$repo" && run_script wing-git.sh worktree-guard "some-branch" 2>/dev/null) || exit_code=$?
  if [ "$exit_code" -eq 0 ]; then
    pass "git: worktree-guard passes when disabled"
  else
    fail "git: worktree-guard should pass when disabled (got: $exit_code)"
  fi
}

test_git_find_rc_none() {
  local repo="$TMPDIR_ROOT/git-find-rc"
  _setup_wing_repo "$repo"
  (cd "$repo" && git commit --allow-empty -m "init" -q)

  local exit_code=0
  (cd "$repo" && run_script wing-git.sh find-rc 2>/dev/null) || exit_code=$?
  if [ "$exit_code" -eq 1 ]; then
    pass "git: find-rc exits 1 when no RC"
  else
    fail "git: find-rc should exit 1 when no RC (got: $exit_code)"
  fi
}

test_git_slugify
test_git_protect_check
test_git_current
test_git_extract_issue_number
test_git_build_branch_name
test_git_next_version
test_git_worktree_guard_disabled
test_git_find_rc_none

# ═══════════════════════════════════════════════════════
# 3. wing-contract.sh
# ═══════════════════════════════════════════════════════

section "wing-contract.sh"

test_contract_hash_deterministic() {
  local h1 h2
  h1=$(run_script wing-contract.sh hash "Title A" "Body A" 2>&1)
  h2=$(run_script wing-contract.sh hash "Title A" "Body A" 2>&1)
  if [[ "$h1" == "$h2" ]] && [[ -n "$h1" ]]; then
    pass "contract: hash is deterministic"
  else
    fail "contract: hash should be deterministic ($h1 != $h2)"
  fi
}

test_contract_hash_different() {
  local h1 h2
  h1=$(run_script wing-contract.sh hash "Title A" "Body A" 2>&1)
  h2=$(run_script wing-contract.sh hash "Title B" "Body B" 2>&1)
  if [[ "$h1" != "$h2" ]]; then
    pass "contract: different inputs → different hash"
  else
    fail "contract: different inputs should produce different hash"
  fi
}

test_contract_read_field() {
  local contract="$TMPDIR_ROOT/MISSION-42.md"
  cat > "$contract" << 'EOF'
---
issue: 42
title: "Fix camera"
branch: mission/42-fix-camera-user
priority: P1
---
# MISSION-42: Fix camera
EOF

  local result
  result=$(run_script wing-contract.sh read-field "$contract" "issue" 2>&1)
  if [[ "$result" == "42" ]]; then
    pass "contract: read-field issue=42"
  else
    fail "contract: read-field issue should be 42 (got: $result)"
  fi

  result=$(run_script wing-contract.sh read-field "$contract" "title" 2>&1)
  if [[ "$result" == "Fix camera" ]]; then
    pass "contract: read-field title"
  else
    fail "contract: read-field title should be 'Fix camera' (got: $result)"
  fi

  result=$(run_script wing-contract.sh read-field "$contract" "nonexistent" 2>&1) || true
  if [[ -z "$result" ]] || [[ "$result" == "" ]]; then
    pass "contract: read-field empty for missing field"
  else
    fail "contract: read-field missing field should be empty (got: $result)"
  fi
}

test_contract_toggle_task() {
  local contract="$TMPDIR_ROOT/MISSION-toggle.md"
  cat > "$contract" << 'EOF'
---
issue: 99
---
## Sub-tasks
- [ ] Task one
- [ ] Task two
- [x] Task three
EOF

  run_script wing-contract.sh toggle-task "$contract" 1 2>/dev/null || true

  if grep -q '\- \[x\] Task one' "$contract"; then
    pass "contract: toggle-task checks off task 1"
  else
    fail "contract: toggle-task should check off task 1"
  fi

  # Task 2 should still be unchecked
  if grep -q '\- \[ \] Task two' "$contract"; then
    pass "contract: toggle-task preserves other tasks"
  else
    fail "contract: toggle-task should preserve other tasks"
  fi
}

test_contract_hash_deterministic
test_contract_hash_different
test_contract_read_field
test_contract_toggle_task

# ═══════════════════════════════════════════════════════
# 4. wing-pr.sh
# ═══════════════════════════════════════════════════════

section "wing-pr.sh"

test_pr_commit_type() {
  local cases=(
    "fix: login bug|fix"
    "feat: add auth|feat"
    "refactor: clean up module|refactor"
    "docs: update README|docs"
    "Add new feature|feat"
  )
  for case in "${cases[@]}"; do
    local input="${case%%|*}"
    local expected="${case##*|}"
    local result
    result=$(run_script wing-pr.sh commit-type "$input" 2>&1) || true
    if [[ "$result" == "$expected" ]]; then
      pass "pr: commit-type '$input' → $result"
    else
      fail "pr: commit-type '$input' expected '$expected' got '$result'"
    fi
  done
}

test_pr_missing_args() {
  local exit_code=0
  run_script wing-pr.sh exists 2>/dev/null || exit_code=$?
  if [ "$exit_code" -ne 0 ]; then
    pass "pr: exists rejects missing args"
  else
    fail "pr: exists should reject missing args"
  fi
}

test_pr_commit_type
test_pr_missing_args

# ═══════════════════════════════════════════════════════
# 5. wing-bump-version.sh
# ═══════════════════════════════════════════════════════

section "wing-bump-version.sh"

test_bump_rejects_no_args() {
  local exit_code=0
  run_script wing-bump-version.sh 2>/dev/null || exit_code=$?
  if [ "$exit_code" -ne 0 ]; then
    pass "bump: rejects no args"
  else
    fail "bump: should reject no args"
  fi
}

test_bump_rejects_invalid() {
  local cases=("3" "3.4" "abc" "v3.4.5" "3.4.5-rc1")
  for v in "${cases[@]}"; do
    local exit_code=0
    run_script wing-bump-version.sh "$v" >/dev/null 2>&1 || exit_code=$?
    if [ "$exit_code" -ne 0 ]; then
      pass "bump: rejects '$v'"
    else
      fail "bump: should reject '$v'"
    fi
  done
}

test_bump_updates_skill_md() {
  local base="$TMPDIR_ROOT/bump-skill"
  local skill_dir="$base/wing"
  local scripts_dir="$skill_dir/scripts"
  mkdir -p "$scripts_dir" "$skill_dir/wing-claim" "$skill_dir/wing-drive"

  cat > "$skill_dir/wing-claim/SKILL.md" << 'EOF'
---
name: wing-claim
version: "0.1.0"
---
# wing-claim
EOF

  cat > "$skill_dir/wing-drive/SKILL.md" << 'EOF'
---
name: wing-drive
version: "0.1.0"
---
# wing-drive
EOF

  cp "$SCRIPT_DIR/wing-bump-version.sh" "$scripts_dir/"

  bash "$scripts_dir/wing-bump-version.sh" "1.0.0" >/dev/null 2>&1 || true

  if grep -q 'version: "1.0.0"' "$skill_dir/wing-claim/SKILL.md"; then
    pass "bump: wing-claim/SKILL.md updated to 1.0.0"
  else
    fail "bump: wing-claim/SKILL.md NOT updated"
  fi

  if grep -q 'version: "1.0.0"' "$skill_dir/wing-drive/SKILL.md"; then
    pass "bump: wing-drive/SKILL.md updated to 1.0.0"
  else
    fail "bump: wing-drive/SKILL.md NOT updated"
  fi
}

test_bump_rejects_no_args
test_bump_rejects_invalid
test_bump_updates_skill_md

# ═══════════════════════════════════════════════════════
# 6. Content Accuracy — SKILL.md files
# ═══════════════════════════════════════════════════════

section "Content Accuracy — SKILL.md"

test_no_stale_tw_references() {
  local found=0
  while IFS= read -r line; do
    found=$((found + 1))
    fail "stale tw- ref: $line"
  done < <(grep -rn 'tw-config\.sh\|tw-git\.sh\|tw-contract\.sh\|tw-pr\.sh\|tw-label\.sh\|tw-notify\.sh' "$WING_DIR" --include="*.md" 2>/dev/null || true)

  if [ "$found" -eq 0 ]; then
    pass "content: no stale tw-*.sh references in .md files"
  fi
}

test_no_stale_team_commands() {
  local found=0
  while IFS= read -r line; do
    found=$((found + 1))
    fail "stale /team ref: $line"
  done < <(grep -rn '/team-claim\|/team-issue\|/team-drive\|/team-ship\|/team-rc\b' "$WING_DIR" --include="*.md" 2>/dev/null || true)

  if [ "$found" -eq 0 ]; then
    pass "content: no stale /team-* command references in .md files"
  fi
}

test_skill_version_consistency() {
  local versions=()
  local all_match=true
  local first=""
  for skill_dir in "$WING_DIR"/wing*/; do
    [[ -f "$skill_dir/SKILL.md" ]] || continue
    local ver
    ver=$(grep '^version:' "$skill_dir/SKILL.md" | head -1 | sed 's/[^"]*"//; s/"//')
    if [[ -z "$first" ]]; then
      first="$ver"
    elif [[ "$ver" != "$first" ]]; then
      all_match=false
      fail "version: $(basename "$skill_dir") has $ver (expected $first)"
    fi
  done

  if $all_match && [[ -n "$first" ]]; then
    pass "content: all SKILL.md versions consistent ($first)"
  fi
}

test_operations_exist() {
  local missing=0

  # Check wing/operations
  for op in init.md doctor.md auto.md; do
    if [[ ! -f "$WING_DIR/wing/operations/$op" ]]; then
      fail "content: wing/operations/$op missing"
      missing=$((missing + 1))
    fi
  done

  # Check wing-issue/operations
  for op in batch.md fix.md update.md; do
    if [[ ! -f "$WING_DIR/wing-issue/operations/$op" ]]; then
      fail "content: wing-issue/operations/$op missing"
      missing=$((missing + 1))
    fi
  done

  # Check wing-drive/operations
  if [[ ! -f "$WING_DIR/wing-drive/operations/team-execution.md" ]]; then
    fail "content: wing-drive/operations/team-execution.md missing"
    missing=$((missing + 1))
  fi

  # Check wing-ship/operations
  if [[ ! -f "$WING_DIR/wing-ship/operations/review.md" ]]; then
    fail "content: wing-ship/operations/review.md missing"
    missing=$((missing + 1))
  fi

  # Check wing-rc/operations
  for op in promote.md hotfix.md; do
    if [[ ! -f "$WING_DIR/wing-rc/operations/$op" ]]; then
      fail "content: wing-rc/operations/$op missing"
      missing=$((missing + 1))
    fi
  done

  if [ "$missing" -eq 0 ]; then
    pass "content: all 10 operations files exist"
  fi
}

test_reference_exists() {
  if [[ -f "$WING_DIR/reference/ve-standard.md" ]]; then
    pass "content: reference/ve-standard.md exists"
  else
    fail "content: reference/ve-standard.md missing"
  fi
}

test_scripts_exist() {
  local missing=0
  for script in wing-config.sh wing-git.sh wing-contract.sh wing-pr.sh wing-label.sh wing-notify.sh wing-bump-version.sh setup-github-labels.sh; do
    if [[ ! -f "$WING_DIR/scripts/$script" ]]; then
      fail "content: scripts/$script missing"
      missing=$((missing + 1))
    fi
  done

  if [ "$missing" -eq 0 ]; then
    pass "content: all 8 scripts exist"
  fi
}

test_scripts_executable() {
  local non_exec=0
  for script in "$WING_DIR"/scripts/*.sh; do
    if [[ ! -x "$script" ]]; then
      fail "content: $(basename "$script") not executable"
      non_exec=$((non_exec + 1))
    fi
  done

  if [ "$non_exec" -eq 0 ]; then
    pass "content: all .sh scripts are executable"
  fi
}

test_no_stale_tw_references
test_no_stale_team_commands
test_skill_version_consistency
test_operations_exist
test_reference_exists
test_scripts_exist
test_scripts_executable

# ═══════════════════════════════════════════════════════
# 7. install.sh
# ═══════════════════════════════════════════════════════

section "install.sh"

test_install_copy_to_temp() {
  local tmp="$TMPDIR_ROOT/install-test"
  mkdir -p "$tmp/.claude/commands"

  local exit_code=0
  bash "$ERWIN_ROOT/install.sh" --copy --target "$tmp" wing >/dev/null 2>&1 || exit_code=$?

  if [ "$exit_code" -eq 0 ]; then
    pass "install: --copy wing exits 0"
  else
    fail "install: --copy wing failed (exit $exit_code)"
  fi

  # Check SKILL.md files installed
  local skill_count=0
  for skill in wing wing-claim wing-issue wing-drive wing-ship wing-rc; do
    if [[ -f "$tmp/.claude/skills/$skill/SKILL.md" ]]; then
      skill_count=$((skill_count + 1))
    else
      fail "install: $skill/SKILL.md not installed"
    fi
  done

  if [ "$skill_count" -eq 6 ]; then
    pass "install: all 6 SKILL.md files installed"
  fi

  # Check reference
  if [[ -f "$tmp/.claude/skills/wing-claim/reference/ve-standard.md" ]]; then
    pass "install: reference/ve-standard.md installed"
  else
    fail "install: reference/ve-standard.md not installed"
  fi

  # Check scripts
  local script_count
  script_count=$(find "$tmp/.claude/skills/wing-claim/scripts" -name "*.sh" 2>/dev/null | wc -l | tr -d ' ')
  if [ "$script_count" -ge 7 ]; then
    pass "install: scripts installed ($script_count .sh files)"
  else
    fail "install: expected ≥7 scripts (got: $script_count)"
  fi
}

test_install_idempotent() {
  local tmp="$TMPDIR_ROOT/install-idempotent"
  mkdir -p "$tmp/.claude/commands"

  bash "$ERWIN_ROOT/install.sh" --copy --target "$tmp" wing >/dev/null 2>&1
  local exit_code=0
  bash "$ERWIN_ROOT/install.sh" --copy --target "$tmp" wing >/dev/null 2>&1 || exit_code=$?

  if [ "$exit_code" -eq 0 ]; then
    pass "install: re-install (idempotent) succeeds"
  else
    fail "install: re-install failed (exit $exit_code)"
  fi
}

test_install_copy_to_temp
test_install_idempotent

# ═══════════════════════════════════════════════════════
# 8. skill-lint.sh
# ═══════════════════════════════════════════════════════

section "skill-lint.sh"

test_lint_all_pass() {
  local exit_code=0
  bash "$ERWIN_ROOT/skills/skill-tester/scripts/skill-lint.sh" --all "$WING_DIR" >/dev/null 2>&1 || exit_code=$?
  if [ "$exit_code" -eq 0 ]; then
    pass "lint: --all skills/wing passes (0 FAIL)"
  else
    fail "lint: --all skills/wing failed (exit $exit_code)"
  fi
}

test_lint_all_pass

# ═══════════════════════════════════════════════════════
# 9. Integration tests (GitHub API) — optional
# ═══════════════════════════════════════════════════════

if $INTEGRATION; then
  section "Integration Tests (GitHub API)"

  # Restore working directory (earlier tests may have cd'd into temp repos)
  cd "$ERWIN_ROOT"

  test_label_operations() {
    info "Creating test Issue..."
    local issue_url
    issue_url=$(gh issue create --title "WING-TEST: automated test — safe to close" \
      --body "Automated test issue from wing-unit-test.sh. Safe to close." \
      --repo "lysfighting/erwin" 2>&1) || {
      fail "integration: could not create test Issue"
      return
    }
    local issue_num
    issue_num=$(echo "$issue_url" | grep -oE '[0-9]+$')
    info "Test Issue: #$issue_num"

    # Test label set
    local exit_code=0
    run_script wing-label.sh set "$issue_num" queued 2>/dev/null || exit_code=$?
    if [ "$exit_code" -eq 0 ]; then
      pass "integration: label set queued"
    else
      fail "integration: label set failed (exit $exit_code)"
    fi

    # Test label verify
    exit_code=0
    local verify_result
    verify_result=$(run_script wing-label.sh verify "$issue_num" queued 2>&1) || exit_code=$?
    if echo "$verify_result" | grep -qi "VERIFIED\|queued"; then
      pass "integration: label verify queued"
    else
      fail "integration: label verify failed ($verify_result)"
    fi

    # Test label transition
    exit_code=0
    run_script wing-label.sh transition "$issue_num" queued wip 2>/dev/null || exit_code=$?
    if [ "$exit_code" -eq 0 ]; then
      pass "integration: label transition queued→wip"
    else
      fail "integration: label transition failed (exit $exit_code)"
    fi

    # Cleanup
    gh issue close "$issue_num" --repo "lysfighting/erwin" --reason "not planned" 2>/dev/null || true
    info "Test Issue #$issue_num closed"
  }

  test_contract_freshness_integration() {
    # Create a contract with known hash, check against a real Issue
    local contract="$TMPDIR_ROOT/MISSION-fresh.md"

    # Use Issue #1 (or any known Issue) — just check the script doesn't crash
    local issue_num=1
    local exit_code=0
    local result
    result=$(run_script wing-contract.sh check-freshness "$contract" "$issue_num" 2>&1) || exit_code=$?

    # Without a real contract file, should exit with NO_HASH or error
    if [[ "$exit_code" -eq 2 ]] || [[ "$result" == *"NO_HASH"* ]] || [[ "$exit_code" -ne 0 ]]; then
      pass "integration: check-freshness handles missing contract gracefully"
    else
      fail "integration: check-freshness should handle missing contract (exit $exit_code, $result)"
    fi
  }

  test_label_operations
  test_contract_freshness_integration
else
  section "Integration Tests (skipped — run with --integration)"
  skip "GitHub API tests (use --integration flag)"
fi

# ═══════════════════════════════════════════════════════
# Summary
# ═══════════════════════════════════════════════════════

echo ""
echo "════════════════════════════════════════════"
echo "  ✅ PASS: $PASS   ❌ FAIL: $FAIL   ⏭️  SKIP: $SKIP"
if [ "$FAIL" -gt 0 ]; then
  echo "  Result: FAILED"
  exit 1
else
  echo "  Result: ALL PASSED"
  exit 0
fi
