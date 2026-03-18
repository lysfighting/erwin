#!/usr/bin/env bash
# wing-bump-version.sh — Bump version across all wing skill files
# Usage: wing-bump-version.sh <version>
# Example: wing-bump-version.sh 3.4.0
#
# Updates:
#   1. Frontmatter `version:` in all skill .md files
#   2. Hardcoded version strings in help text and init templates
#   3. Reports what changed

set -euo pipefail

# Portable in-place sed (macOS uses -i '', GNU/Linux uses -i)
if sed --version 2>/dev/null | grep -q GNU 2>/dev/null; then
  _sed_i() { sed -i "$@"; }
else
  _sed_i() { sed -i '' "$@"; }
fi

VERSION="${1:-}"
if [ -z "$VERSION" ] || [[ "$VERSION" == -* ]]; then
  echo "Usage: wing-bump-version.sh <version>"
  echo "Example: wing-bump-version.sh 3.4.0"
  exit 1
fi
if ! echo "$VERSION" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$'; then
  echo "ERROR: Version must be semver (e.g., 3.8.4), got '$VERSION'"
  exit 1
fi

# Derive major.minor for two-segment matches (e.g., "3.4.0" → "3.4")
VERSION_SHORT=$(echo "$VERSION" | sed 's/\.[0-9]*$//')

# Find skill directory (works from anywhere)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILL_DIR="$(dirname "$SCRIPT_DIR")"

CHANGED=0

# Collect all files that need version updates
FILES=()
for file in "$SKILL_DIR"/wing*/SKILL.md "$SKILL_DIR"/ONBOARDING.md "$SKILL_DIR"/README.md; do
  [ -f "$file" ] && FILES+=("$file")
done
# Also update the manual (lives outside skill dir)
# Manual doc (if it exists in erwin docs)
MANUAL="$(dirname "$SKILL_DIR")/../docs/wing-manual.md"
[ -f "$MANUAL" ] && FILES+=("$MANUAL")

# 1. Update frontmatter version: in all skill .md files
for file in "${FILES[@]}"; do
  OLD=$(grep -m1 '^version:' "$file" 2>/dev/null | sed 's/version: *"//' | sed 's/"//') || true
  if [ -n "$OLD" ] && [ "$OLD" != "$VERSION" ]; then
    _sed_i "s/^version: \".*\"/version: \"$VERSION\"/" "$file"
    echo "  $(basename "$file"): $OLD → $VERSION"
    CHANGED=$((CHANGED + 1))
  fi
done

# 2. Update hardcoded version in help text and init templates
# Matches both 3-segment (v3.4.0) and 2-segment (v3.4) version strings
for file in "${FILES[@]}"; do
  # Snapshot file before sed passes to detect if anything changed
  _before=$(shasum -a 256 "$file" | cut -d' ' -f1)

  # "TEAMWORK v3.x.y" or "TEAMWORK v3.x" (2 or 3 segments)
  _sed_i "s/TEAMWORK v[0-9][0-9]*\.[0-9][0-9]*\(\.[0-9][0-9]*\)*/TEAMWORK v$VERSION/g" "$file"

  # "Teamwork v3.x.y" or "Teamwork v3.x"
  _sed_i "s/Teamwork v[0-9][0-9]*\.[0-9][0-9]*\(\.[0-9][0-9]*\)*/Teamwork v$VERSION/g" "$file"

  # "(v3.x.y)" or "(v3.x)"
  _sed_i "s/(v[0-9][0-9]*\.[0-9][0-9]*\(\.[0-9][0-9]*\)*)/(v$VERSION)/g" "$file"

  # "Version 3.x.y" (manual header)
  _sed_i "s/Version [0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*/Version $VERSION/g" "$file"

  # "Skill version: 3.x.y" (README)
  _sed_i "s/Skill version: [0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*/Skill version: $VERSION/g" "$file"

  # "**Version**: `3.x.y`" (inline help docs — backtick-wrapped)
  _sed_i "s/\*\*Version\*\*: \`[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*\`/**Version**: \`$VERSION\`/g" "$file"

  # "(`v3.x.y`)" (inline backtick-parens in titles)
  _sed_i "s/(\`v[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*\`)/(\`v$VERSION\`)/g" "$file"

  # 'SKILL_VERSION="3.x.y"' (bash variable in init templates)
  _sed_i "s/SKILL_VERSION=\"[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*\"/SKILL_VERSION=\"$VERSION\"/g" "$file"

  # skill_version removed from config — no longer needs bumping

  # Count file as changed if any inline pattern matched (avoid double-counting with frontmatter loop)
  _after=$(shasum -a 256 "$file" | cut -d' ' -f1)
  if [ "$_before" != "$_after" ]; then
    CHANGED=$((CHANGED + 1))
  fi
done

echo ""
echo "Bumped $CHANGED file(s) to v$VERSION"
echo "Don't forget to update CHANGELOG.md"
