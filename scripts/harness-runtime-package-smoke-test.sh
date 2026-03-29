#!/usr/bin/env bash
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

TMPDIR_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/harness-runtime-package-smoke.XXXXXX")
trap 'rm -rf "$TMPDIR_ROOT"' EXIT

PASS=0
FAIL=0

pass() { PASS=$((PASS + 1)); echo "  PASS  $1"; }
fail() { FAIL=$((FAIL + 1)); echo "  FAIL  $1"; }
section() { echo ""; echo "=== $1 ==="; }

normalize_os() {
  case "$(uname -s)" in
    Darwin) printf 'darwin\n' ;;
    Linux) printf 'linux\n' ;;
    *) printf '%s\n' "$(uname -s | tr '[:upper:]' '[:lower:]')" ;;
  esac
}

normalize_arch() {
  case "$(uname -m)" in
    x86_64|amd64) printf 'amd64\n' ;;
    arm64|aarch64) printf 'arm64\n' ;;
    *) printf '%s\n' "$(uname -m)" ;;
  esac
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

target="$(normalize_os)-$(normalize_arch)"
output_root="$TMPDIR_ROOT/out"

section "Package runtime"
package_output="$(bash "$REPO_ROOT/scripts/harness-package-runtime.sh" --target "$target" --output "$output_root")"
assert_contains "$package_output" "\"$target\"" "package output reports requested target"

archive_path="$output_root/harness-runtime-$target.tar.gz"
if [[ -f "$archive_path" ]]; then
  pass "runtime package archive created"
else
  fail "runtime package archive created"
fi

section "Archive contents"
archive_list="$(tar -tzf "$archive_path")"
assert_contains "$archive_list" "$target/harness-state" "archive contains harness-state binary"
assert_contains "$archive_list" "$target/harness/adapter-manifest.yml" "archive contains adapter manifest"
assert_contains "$archive_list" "$target/docs/harness-public-contract.md" "archive contains public contract"

section "Summary"
echo "PASS=$PASS FAIL=$FAIL"
if [[ "$FAIL" -ne 0 ]]; then
  exit 1
fi
