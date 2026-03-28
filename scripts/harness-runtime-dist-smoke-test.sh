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
    chmod +x scripts/*.sh install-codex.sh
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
assert_contains "$adapter_output" '"valid": true' "adapter validate works via packaged binary without Go"
assert_contains "$home_output" '"next_action": "No active runs. Claim a mission."' "home works via packaged binary without Go"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [ "$FAIL" -ne 0 ]; then
  exit 1
fi
