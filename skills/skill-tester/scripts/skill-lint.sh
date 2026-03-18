#!/usr/bin/env bash
# skill-lint.sh — Layer 1 static analysis for erwin skills
# Checks SKILL.md files against known quality rules. No LLM needed.
#
# Usage:
#   ./skill-lint.sh <skill-dir>          Lint a single skill directory
#   ./skill-lint.sh --all <skills-root>  Lint all skills under a directory
#   ./skill-lint.sh --help               Show usage

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

PASS=0
FAIL=0
WARN=0

pass() {
  echo -e "  ${GREEN}PASS${NC}  $1"
  PASS=$((PASS + 1))
}

fail() {
  echo -e "  ${RED}FAIL${NC}  $1"
  FAIL=$((FAIL + 1))
}

warn() {
  echo -e "  ${YELLOW}WARN${NC}  $1"
  WARN=$((WARN + 1))
}

# ─── Check functions ───────────────────────────────────────────────

check_file_exists() {
  local dir="$1"
  if [[ -f "$dir/SKILL.md" ]]; then
    pass "SKILL.md exists (uppercase)"
  elif [[ -f "$dir/skill.md" ]]; then
    fail "skill.md should be SKILL.md (uppercase per Anthropic convention)"
  else
    fail "No SKILL.md found in $dir"
    return 1
  fi
}

check_size() {
  local file="$1"
  local lines
  lines=$(wc -l < "$file" | tr -d ' ')
  if [[ $lines -le 500 ]]; then
    pass "Size: $lines lines (≤500)"
  else
    fail "Size: $lines lines (exceeds 500 line limit)"
  fi
}

check_frontmatter() {
  local file="$1"
  local has_name has_desc has_version
  has_name=$(head -20 "$file" | grep -c '^name:' || true)
  has_desc=$(head -20 "$file" | grep -c '^description:' || true)
  has_version=$(head -20 "$file" | grep -c '^version:' || true)

  # Check frontmatter delimiters
  local first_line
  first_line=$(head -1 "$file")
  if [[ "$first_line" != "---" ]]; then
    fail "Frontmatter: missing opening --- delimiter"
    return
  fi

  if [[ $has_name -ge 1 ]]; then
    pass "Frontmatter: has name field"
  else
    fail "Frontmatter: missing name field"
  fi

  if [[ $has_desc -ge 1 ]]; then
    pass "Frontmatter: has description field"
  else
    fail "Frontmatter: missing description field"
  fi

  if [[ $has_version -ge 1 ]]; then
    pass "Frontmatter: has version field"
  else
    warn "Frontmatter: no version field"
  fi
}

check_ve_reference() {
  local file="$1"
  if grep -q 've-standard\.md\|Visual Encoding' "$file"; then
    pass "VE reference: found"
  else
    warn "VE reference: no reference to ve-standard.md or Visual Encoding"
  fi
}

check_bare_issue_numbers() {
  local file="$1"
  # Look for bare #NN (2 digits) outside code blocks
  # This is a simplified check — skips content inside ``` blocks
  local in_code=0
  local line_num=0
  local violations=""

  while IFS= read -r line; do
    ((line_num++))
    if [[ "$line" =~ ^\`\`\` ]]; then
      in_code=$((1 - in_code))
      continue
    fi
    if [[ $in_code -eq 1 ]]; then
      continue
    fi
    # Match #NN where NN is exactly 2 digits (not followed by more digits)
    # Exclude patterns inside backticks, bold, and common false positives
    # Skip frontmatter, comments, and markdown link anchors
    if echo "$line" | grep -qE '(^|[^&0-9#])#[0-9]{1,2}([^0-9]|$)' 2>/dev/null; then
      # Exclude lines that are clearly not issue references
      if ! echo "$line" | grep -qE '(^#+ |^\s*#|bug #19251|ANSI|rgb|H[0-9]|P[0-9]|S-|step [0-9]|Step [0-9]|`#)' 2>/dev/null; then
        violations="${violations}    line $line_num: $line\n"
      fi
    fi
  done < "$file"

  if [[ -z "$violations" ]]; then
    pass "Issue numbers: no bare 2-digit #NN found"
  else
    warn "Issue numbers: potential bare #NN (should be **#NNN** with 3+ digits):"
    echo -e "$violations"
  fi
}

check_bold_issue_refs() {
  local file="$1"
  # In non-code-block text, #NNN (3+ digits) should be wrapped in **bold**
  local in_code=0
  local line_num=0
  local violations=""

  while IFS= read -r line; do
    ((line_num++))
    if [[ "$line" =~ ^\`\`\` ]]; then
      in_code=$((1 - in_code))
      continue
    fi
    if [[ $in_code -eq 1 ]]; then
      continue
    fi
    # Find #NNN that is NOT inside **bold** and NOT inside backticks
    # Match bare #NNN (not preceded by ** or `)
    if echo "$line" | grep -qE '(^|[^*`])#[0-9]{3,}' 2>/dev/null; then
      # Exclude lines where #NNN is already in **#NNN** or `#NNN`
      if ! echo "$line" | grep -qE '\*\*#[0-9]{3,}\*\*|`[^`]*#[0-9]{3,}[^`]*`' 2>/dev/null; then
        # Exclude frontmatter, headings, and format examples
        if ! echo "$line" | grep -qE '(^---$|^#+ |issue_content_hash|Pattern|Renders)' 2>/dev/null; then
          violations="${violations}    line $line_num: $line\n"
        fi
      fi
    fi
  done < "$file"

  if [[ -z "$violations" ]]; then
    pass "Issue refs: all #NNN wrapped in **bold**"
  else
    warn "Issue refs: #NNN should be **#NNN** (bold for auto-linking):"
    echo -e "$violations"
  fi
}

check_no_h3_in_output() {
  local file="$1"
  # ### headers cause 2-newline gap — should use **bold** instead
  # Only flag ### that appear to be in output sections (not in code blocks or markdown structure)
  local in_code=0
  local line_num=0
  local violations=""

  while IFS= read -r line; do
    ((line_num++))
    if [[ "$line" =~ ^\`\`\` ]]; then
      in_code=$((1 - in_code))
      continue
    fi
    if [[ $in_code -eq 1 ]]; then
      continue
    fi
    # ### inside a skill file is likely an output template header
    # Exception: ### in routing tables (like "### If `list`:") is structural markdown
    if [[ "$line" =~ ^###\  ]] && ! echo "$line" | grep -qE '^### (If |Step |Error)' 2>/dev/null; then
      violations="${violations}    line $line_num: $line\n"
    fi
  done < "$file"

  if [[ -z "$violations" ]]; then
    pass "Headers: no ### in output templates"
  else
    warn "Headers: ### causes 2-newline gap in output (use **bold** instead):"
    echo -e "$violations"
  fi
}

check_code_blocks_around_output() {
  local file="$1"
  # Code blocks wrapping output templates kill all inline formatting
  # Heuristic: code blocks that contain **bold** or *italic* are suspicious
  local in_code=0
  local code_start=0
  local line_num=0
  local has_inline_formatting=0
  local violations=""

  while IFS= read -r line; do
    ((line_num++))
    if [[ "$line" =~ ^\`\`\` ]]; then
      if [[ $in_code -eq 0 ]]; then
        in_code=1
        code_start=$line_num
        has_inline_formatting=0
      else
        if [[ $has_inline_formatting -eq 1 ]]; then
          violations="${violations}    lines $code_start-$line_num: code block contains inline formatting\n"
        fi
        in_code=0
      fi
      continue
    fi
    if [[ $in_code -eq 1 ]]; then
      if echo "$line" | grep -qE '\*\*[^*]+\*\*|\*[^*]+\*' 2>/dev/null; then
        # Exclude bash/yaml/markdown template code blocks
        if ! echo "$line" | grep -qE '(^(if|fi|for|do|done|echo|eval|local|#|---|\s*-)|\.sh|\.md|\.yml)' 2>/dev/null; then
          has_inline_formatting=1
        fi
      fi
    fi
  done < "$file"

  if [[ -z "$violations" ]]; then
    pass "Code blocks: no output templates wrapped in code blocks"
  else
    warn "Code blocks: inline formatting killed inside code blocks:"
    echo -e "$violations"
  fi
}

check_dividers() {
  local file="$1"
  local in_code=0
  local line_num=0
  local violations=""

  while IFS= read -r line; do
    ((line_num++))
    if [[ "$line" =~ ^\`\`\` ]]; then
      in_code=$((1 - in_code))
      continue
    fi
    if [[ $in_code -eq 1 ]]; then
      continue
    fi
    # Check for ═══ (double-line) dividers
    if echo "$line" | grep -q '═══' 2>/dev/null; then
      violations="${violations}    line $line_num: uses ═══ (should be ────)\n"
    fi
  done < "$file"

  if [[ -z "$violations" ]]; then
    pass "Dividers: no ═══ found"
  else
    fail "Dividers: use ──── instead of ═══:"
    echo -e "$violations"
  fi
}

check_supporting_files() {
  local file="$1"
  local dir
  dir="$(dirname "$file")"
  local violations=""

  # Find file references in SKILL.md (${CLAUDE_SKILL_DIR}/... patterns)
  while IFS= read -r ref; do
    # Resolve relative to skill dir
    local resolved
    resolved=$(echo "$ref" | sed "s|\${CLAUDE_SKILL_DIR}|$dir|g")
    # Only check local references (not absolute paths or variables)
    if [[ "$resolved" == "$dir"* ]] && [[ ! -f "$resolved" ]]; then
      violations="${violations}    $ref → not found\n"
    fi
  done < <(grep -oE '\$\{CLAUDE_SKILL_DIR\}/[^ )"]+' "$file" 2>/dev/null | sort -u || true)

  if [[ -z "$violations" ]]; then
    pass "Supporting files: all references resolve"
  else
    warn "Supporting files: some referenced files not found:"
    echo -e "$violations"
  fi
}

check_version_consistency() {
  local skills_root="$1"
  # Check that all SKILL.md files in the same suite have matching versions
  local versions=()
  local files=()

  for skill_dir in "$skills_root"/*/; do
    [[ -f "$skill_dir/SKILL.md" ]] || continue
    local ver
    ver=$(grep '^version:' "$skill_dir/SKILL.md" 2>/dev/null | head -1 | sed 's/^[^:]*://' | sed 's/^ *//' | tr -d '"' || true)
    if [[ -n "$ver" ]]; then
      versions+=("$ver")
      files+=("$(basename "$skill_dir")/SKILL.md=$ver")
    fi
  done

  if [[ ${#versions[@]} -le 1 ]]; then
    pass "Version consistency: single skill (or no versions found)"
    return
  fi

  local first="${versions[0]}"
  local all_match=1
  for v in "${versions[@]}"; do
    if [[ "$v" != "$first" ]]; then
      all_match=0
      break
    fi
  done

  if [[ $all_match -eq 1 ]]; then
    pass "Version consistency: all skills at v$first"
  else
    fail "Version consistency: mismatched versions:"
    for f in "${files[@]}"; do
      echo -e "    $f"
    done
  fi
}

# ─── Main ──────────────────────────────────────────────────────────

lint_skill_dir() {
  local dir="$1"
  local name
  name="$(basename "$dir")"
  echo -e "\n${CYAN}Linting: $name${NC}"
  echo "────────────────────────────────"

  if ! check_file_exists "$dir"; then
    return
  fi

  local file="$dir/SKILL.md"
  check_size "$file"
  check_frontmatter "$file"
  check_ve_reference "$file"
  check_bare_issue_numbers "$file"
  check_bold_issue_refs "$file"
  check_no_h3_in_output "$file"
  check_code_blocks_around_output "$file"
  check_dividers "$file"
  check_supporting_files "$file"
}

usage() {
  echo "skill-lint.sh — Layer 1 static analysis for erwin skills"
  echo ""
  echo "Usage:"
  echo "  skill-lint.sh <skill-dir>            Lint a single skill"
  echo "  skill-lint.sh --all <skills-root>    Lint all skills under root"
  echo "  skill-lint.sh --help                 Show this help"
  echo ""
  echo "Examples:"
  echo "  skill-lint.sh skills/wing/wing-claim"
  echo "  skill-lint.sh --all skills/wing"
}

if [[ $# -eq 0 ]] || [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
  usage
  exit 0
fi

if [[ "$1" == "--all" ]]; then
  if [[ -z "${2:-}" ]]; then
    echo "ERROR: --all requires a skills root directory"
    exit 1
  fi
  SKILLS_ROOT="$2"
  echo -e "${CYAN}skill-lint.sh — Linting all skills in $SKILLS_ROOT${NC}"

  for skill_dir in "$SKILLS_ROOT"/*/; do
    [[ -d "$skill_dir" ]] || continue
    # Skip non-skill directories
    [[ "$(basename "$skill_dir")" == "scripts" ]] && continue
    [[ "$(basename "$skill_dir")" == "reference" ]] && continue
    lint_skill_dir "$skill_dir"
  done

  # Cross-skill checks
  echo -e "\n${CYAN}Cross-skill checks${NC}"
  echo "────────────────────────────────"
  check_version_consistency "$SKILLS_ROOT"
else
  lint_skill_dir "$1"
fi

# ─── Summary ───────────────────────────────────────────────────────

echo ""
echo "════════════════════════════════"
echo -e "  ${GREEN}PASS${NC}: $PASS  ${RED}FAIL${NC}: $FAIL  ${YELLOW}WARN${NC}: $WARN"
if [[ $FAIL -gt 0 ]]; then
  echo -e "  ${RED}Result: FAILED${NC}"
  exit 1
else
  echo -e "  ${GREEN}Result: PASSED${NC}"
  exit 0
fi
