#!/usr/bin/env bash
# erwin skills installer
# Installs erwin skills to ~/.claude/commands/ or ~/.claude/skills/ for availability
#
# Two formats supported:
#   Legacy (teamwork):  flat .md files    → ~/.claude/commands/skill-name.md
#   Modern (wing):      SKILL.md dirs     → ~/.claude/skills/skill-name/SKILL.md
#
# Usage:
#   ./install.sh                        # symlink all skills (global, auto-updates via git pull)
#   ./install.sh teamwork               # symlink one skill by name (legacy)
#   ./install.sh wing                   # symlink one skill suite (modern)
#   ./install.sh --copy teamwork        # copy (for cross-project deployment)
#   ./install.sh --copy --target DIR    # copy to specific directory
#   ./install.sh --copy --target /path/to/project/.claude/commands teamwork

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILLS_DIR="$SCRIPT_DIR/skills"
TARGET_DIR="$HOME/.claude/commands"
SKILLS_TARGET_DIR="$HOME/.claude/skills"  # Modern skills target
INSTALL_MODE="symlink"  # symlink (default) or copy

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Parse flags
POSITIONAL_ARGS=()
while [[ $# -gt 0 ]]; do
  case $1 in
    --copy|-c)
      INSTALL_MODE="copy"
      shift
      ;;
    --target|-t)
      if [ -z "${2:-}" ]; then
        echo "ERROR: --target requires a directory path"
        exit 1
      fi
      TARGET_DIR="$2"
      # Auto-detect project root: if target has .claude/ but isn't already .claude/commands, append it
      if [[ -d "$TARGET_DIR/.claude" && ! "$TARGET_DIR" == */.claude/commands ]]; then
        TARGET_DIR="$TARGET_DIR/.claude/commands"
      fi
      # Derive modern skills target from legacy target
      # .claude/commands → .claude/skills (sibling directory)
      if [[ "$TARGET_DIR" == */.claude/commands ]]; then
        SKILLS_TARGET_DIR="${TARGET_DIR%/commands}/skills"
      else
        SKILLS_TARGET_DIR="$TARGET_DIR"
      fi
      shift 2
      ;;
    --help|-h)
      echo "erwin skills installer"
      echo ""
      echo "Usage:"
      echo "  ./install.sh                            Symlink all skills (global)"
      echo "  ./install.sh teamwork                   Symlink one skill"
      echo "  ./install.sh --copy teamwork            Copy one skill (standalone)"
      echo "  ./install.sh --copy --target DIR        Copy all skills to DIR"
      echo "  ./install.sh --copy --target DIR skill  Copy one skill to DIR"
      echo ""
      echo "Modes:"
      echo "  symlink (default)  Auto-updates via git pull. Best for local dev."
      echo "  --copy             Standalone files. Best for cross-project deployment."
      echo ""
      echo "Options:"
      echo "  --copy, -c         Copy files instead of symlinking"
      echo "  --target, -t DIR   Install to DIR instead of ~/.claude/commands/"
      echo "  --help, -h         Show this help"
      exit 0
      ;;
    *)
      POSITIONAL_ARGS+=("$1")
      shift
      ;;
  esac
done
set -- "${POSITIONAL_ARGS[@]+"${POSITIONAL_ARGS[@]}"}"

echo "erwin skills installer"
echo "======================"
if [[ "$INSTALL_MODE" == "copy" ]]; then
  echo -e "Mode: ${CYAN}copy${NC} (standalone files)"
else
  echo -e "Mode: ${CYAN}symlink${NC} (auto-updates via git pull)"
fi
echo -e "Target: ${CYAN}$TARGET_DIR${NC}"
echo ""

mkdir -p "$TARGET_DIR"

# Clean up renamed/removed skills (v1 → v2 migrations)
DEPRECATED_SKILLS=(teamwork team-get-mission team-complete team-release)
for old_skill in "${DEPRECATED_SKILLS[@]}"; do
  old_dst="$TARGET_DIR/$old_skill.md"
  if [[ -L "$old_dst" ]] || [[ -f "$old_dst" ]]; then
    rm "$old_dst"
    echo -e "${YELLOW}  Removed deprecated: /$old_skill${NC}"
  fi
done

install_single_file() {
  local src="$1"
  local filename="$(basename "$src" .md)"
  local dst="$TARGET_DIR/$filename.md"

  # Extract version from frontmatter if present
  local version=""
  version=$(grep '^version:' "$src" 2>/dev/null | head -1 | sed 's/^[^:]*://' | sed 's/^ *//' | tr -d '"' || true)

  # Remove existing file (version history lives in erwin repo — no backup needed)
  if [[ -f "$dst" || -L "$dst" ]]; then
    rm "$dst"
  fi

  if [[ "$INSTALL_MODE" == "copy" ]]; then
    cp "$src" "$dst"
    local label="copied"
  else
    ln -sf "$src" "$dst"
    local label="→ $src"
  fi

  if [[ -n "$version" ]]; then
    echo -e "${GREEN}✅ /$filename${NC} v$version  $label"
  else
    echo -e "${GREEN}✅ /$filename${NC}  $label"
  fi
}

install_skill() {
  local name="$1"
  local skill_dir="$SKILLS_DIR/$name"

  if [[ ! -d "$skill_dir" ]]; then
    echo -e "${RED}✗ Skill not found: $name${NC}"
    echo "  Available: $(ls "$SKILLS_DIR" | grep -v README | tr '\n' ' ')"
    return 1
  fi

  # Find all .md files that are skills (exclude README, CHANGELOG, ONBOARDING)
  local found=0
  for md_file in "$skill_dir"/*.md; do
    [[ -f "$md_file" ]] || continue
    local basename="$(basename "$md_file")"
    [[ "$basename" == "README.md" ]] && continue
    [[ "$basename" == "CHANGELOG.md" ]] && continue
    [[ "$basename" == "ONBOARDING.md" ]] && continue
    install_single_file "$md_file"
    ((found++))
  done

  if [[ $found -eq 0 ]]; then
    echo -e "${RED}✗ No skill files found in: $skill_dir${NC}"
    return 1
  fi

  # Install scripts/ subdir if present
  if [[ -d "$skill_dir/scripts" ]]; then
    local scripts_dst="$TARGET_DIR/scripts"
    mkdir -p "$scripts_dst"
    for script in "$skill_dir/scripts/"*.sh "$skill_dir/scripts/"*.txt; do
      [[ -f "$script" ]] || continue
      local script_name="$(basename "$script")"
      local script_dst="$scripts_dst/$script_name"
      # Remove existing symlink before copy (prevents "identical" error)
      [[ -L "$script_dst" ]] && rm "$script_dst"
      if [[ "$INSTALL_MODE" == "copy" ]]; then
        cp "$script" "$script_dst"
      else
        ln -sf "$script" "$script_dst"
      fi
      [[ "$script_name" == *.sh ]] && chmod +x "$script_dst"
      echo -e "${GREEN}  📜 scripts/$script_name${NC}"
    done
  fi

  # Install reference/ subdir if present (preserving nested directory structure).
  # Needed for skills like teamwork 3.14.0 that ship a briefing template +
  # checkpoint schemas under reference/checkpoint-schemas/. Without this,
  # runtime lookups like `~/.claude/commands/reference/...` break.
  if [[ -d "$skill_dir/reference" ]]; then
    local ref_src="$skill_dir/reference"
    local ref_dst="$TARGET_DIR/reference"
    mkdir -p "$ref_dst"
    local ref_count=0
    # find -print0 + while-read -d '' handles spaces in filenames safely
    while IFS= read -r -d '' src_file; do
      local rel_path="${src_file#$ref_src/}"
      local dst_file="$ref_dst/$rel_path"
      mkdir -p "$(dirname "$dst_file")"
      [[ -L "$dst_file" || -f "$dst_file" ]] && rm "$dst_file"
      if [[ "$INSTALL_MODE" == "copy" ]]; then
        cp "$src_file" "$dst_file"
      else
        ln -sf "$src_file" "$dst_file"
      fi
      ((ref_count++))
    done < <(find "$ref_src" -type f -print0)
    echo -e "${GREEN}  📚 reference/ ($ref_count files)${NC}"
  fi
}

# ─── Modern skill suite installer (SKILL.md in subdirs → ~/.claude/skills/) ───

is_modern_skill() {
  # A modern skill suite has subdirectories containing SKILL.md files
  local dir="$1"
  for sub in "$dir"/*/; do
    [[ -f "$sub/SKILL.md" ]] && return 0
  done
  return 1
}

install_skill_suite() {
  local name="$1"
  local suite_dir="$SKILLS_DIR/$name"

  if [[ ! -d "$suite_dir" ]]; then
    echo -e "${RED}✗ Skill suite not found: $name${NC}"
    return 1
  fi

  # Install each sub-skill (directories containing SKILL.md)
  local found=0
  for skill_sub in "$suite_dir"/*/; do
    [[ -d "$skill_sub" ]] || continue
    local sub_name
    sub_name="$(basename "$skill_sub")"

    # Skip non-skill directories
    [[ "$sub_name" == "scripts" ]] && continue
    [[ "$sub_name" == "reference" ]] && continue

    [[ -f "$skill_sub/SKILL.md" ]] || continue

    local dst_dir="$SKILLS_TARGET_DIR/$sub_name"
    mkdir -p "$dst_dir"

    # Extract version
    local version=""
    version=$(grep '^version:' "$skill_sub/SKILL.md" 2>/dev/null | head -1 | sed 's/^[^:]*://' | sed 's/^ *//' | tr -d '"' || true)

    # Install SKILL.md
    local dst="$dst_dir/SKILL.md"
    [[ -f "$dst" || -L "$dst" ]] && rm "$dst"
    if [[ "$INSTALL_MODE" == "copy" ]]; then
      cp "$skill_sub/SKILL.md" "$dst"
      local label="copied"
    else
      ln -sf "$skill_sub/SKILL.md" "$dst"
      local label="→ $skill_sub/SKILL.md"
    fi

    if [[ -n "$version" ]]; then
      echo -e "${GREEN}✅ /$sub_name${NC} v$version  $label"
    else
      echo -e "${GREEN}✅ /$sub_name${NC}  $label"
    fi

    # Install operations/ if present
    if [[ -d "$skill_sub/operations" ]]; then
      local ops_dst="$dst_dir/operations"
      mkdir -p "$ops_dst"
      for op_file in "$skill_sub/operations/"*.md; do
        [[ -f "$op_file" ]] || continue
        local op_name
        op_name="$(basename "$op_file")"
        local op_dst="$ops_dst/$op_name"
        [[ -f "$op_dst" || -L "$op_dst" ]] && rm "$op_dst"
        if [[ "$INSTALL_MODE" == "copy" ]]; then
          cp "$op_file" "$op_dst"
        else
          ln -sf "$op_file" "$op_dst"
        fi
        echo -e "${GREEN}    📄 operations/$op_name${NC}"
      done
    fi

    ((found++))
  done

  if [[ $found -eq 0 ]]; then
    echo -e "${RED}✗ No SKILL.md files found in: $suite_dir${NC}"
    return 1
  fi

  # Install shared reference/ directory
  if [[ -d "$suite_dir/reference" ]]; then
    for skill_sub in "$suite_dir"/*/; do
      local sub_name
      sub_name="$(basename "$skill_sub")"
      [[ "$sub_name" == "scripts" || "$sub_name" == "reference" ]] && continue
      [[ -f "$skill_sub/SKILL.md" ]] || continue

      local ref_dst="$SKILLS_TARGET_DIR/$sub_name/reference"
      mkdir -p "$ref_dst"
      for ref_file in "$suite_dir/reference/"*; do
        [[ -f "$ref_file" ]] || continue
        local ref_name
        ref_name="$(basename "$ref_file")"
        local ref_file_dst="$ref_dst/$ref_name"
        [[ -f "$ref_file_dst" || -L "$ref_file_dst" ]] && rm "$ref_file_dst"
        if [[ "$INSTALL_MODE" == "copy" ]]; then
          cp "$ref_file" "$ref_file_dst"
        else
          ln -sf "$ref_file" "$ref_file_dst"
        fi
      done
    done
    echo -e "${GREEN}    📚 reference/ → all skills${NC}"
  fi

  # Install shared scripts/ directory
  if [[ -d "$suite_dir/scripts" ]]; then
    for skill_sub in "$suite_dir"/*/; do
      local sub_name
      sub_name="$(basename "$skill_sub")"
      [[ "$sub_name" == "scripts" || "$sub_name" == "reference" ]] && continue
      [[ -f "$skill_sub/SKILL.md" ]] || continue

      local scripts_dst="$SKILLS_TARGET_DIR/$sub_name/scripts"
      mkdir -p "$scripts_dst"
      for script in "$suite_dir/scripts/"*; do
        [[ -f "$script" ]] || continue
        local script_name
        script_name="$(basename "$script")"
        local script_dst="$scripts_dst/$script_name"
        [[ -L "$script_dst" ]] && rm "$script_dst"
        if [[ "$INSTALL_MODE" == "copy" ]]; then
          cp "$script" "$script_dst"
        else
          ln -sf "$script" "$script_dst"
        fi
        [[ "$script_name" == *.sh ]] && chmod +x "$script_dst"
      done
    done
    echo -e "${GREEN}    📜 scripts/ → all skills${NC}"
  fi

  echo -e "  ${CYAN}Installed to: $SKILLS_TARGET_DIR/${NC}"
}

# ─── Main routing ─────────────────────────────────────────────────

if [[ $# -eq 0 ]]; then
  # Safety: when targeting another project, require explicit skill name
  if [[ "$TARGET_DIR" != "$HOME/.claude/commands" ]]; then
    echo -e "${RED}ERROR: --target requires a skill name to prevent accidental full install${NC}"
    echo ""
    echo "Usage: ./install.sh --copy --target DIR <skill-name>"
    echo ""
    echo "Available skills:"
    for dir in "$SKILLS_DIR"/*/; do
      [[ -d "$dir" ]] || continue
      skill_name="$(basename "$dir")"
      [[ "$skill_name" == "README"* ]] && continue
      echo "  $skill_name"
    done
    exit 1
  fi
  # Install all skills
  echo "Installing all skills..."
  echo ""
  skill_count=0
  for dir in "$SKILLS_DIR"/*/; do
    [[ -d "$dir" ]] || continue
    skill_name="$(basename "$dir")"
    [[ "$skill_name" == "README"* ]] && continue
    if is_modern_skill "$dir"; then
      install_skill_suite "$skill_name" && skill_count=$((skill_count + 1))
    else
      install_skill "$skill_name" && skill_count=$((skill_count + 1))
    fi
  done
  echo ""
  echo "Installed $skill_count skill(s)"
  echo "  Legacy skills: $TARGET_DIR"
  echo "  Modern skills: $SKILLS_TARGET_DIR"
  if [[ "$INSTALL_MODE" == "symlink" ]]; then
    echo ""
    echo "To update: git pull  (symlinks auto-update)"
  else
    echo ""
    echo "To update: re-run with --copy  (files are standalone)"
  fi

elif [[ -d "$SKILLS_DIR/$1" ]]; then
  if is_modern_skill "$SKILLS_DIR/$1"; then
    install_skill_suite "$1"
  else
    install_skill "$1"
  fi

else
  echo -e "${RED}Unknown skill: $1${NC}"
  echo ""
  echo "Available skills:"
  for dir in "$SKILLS_DIR"/*/; do
    [[ -d "$dir" ]] || continue
    skill_name="$(basename "$dir")"
    [[ "$skill_name" == "README"* ]] && continue
    if is_modern_skill "$dir"; then
      echo "  $skill_name (modern → ~/.claude/skills/)"
    else
      echo "  $skill_name (legacy → ~/.claude/commands/)"
    fi
  done
  exit 1
fi
