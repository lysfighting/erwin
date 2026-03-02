#!/usr/bin/env bash
# erwin skills installer
# Installs erwin skills to ~/.claude/commands/ (or custom target) for availability
#
# Usage:
#   ./install.sh                        # symlink all skills (global, auto-updates via git pull)
#   ./install.sh teamwork               # symlink one skill by name
#   ./install.sh --copy teamwork        # copy (for cross-project deployment)
#   ./install.sh --copy --target DIR    # copy to specific directory
#   ./install.sh --copy --target /path/to/project/.claude/commands teamwork

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILLS_DIR="$SCRIPT_DIR/skills"
TARGET_DIR="$HOME/.claude/commands"
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
      TARGET_DIR="$2"
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
DEPRECATED_SKILLS=(teamwork team-get-mission team-complete)
for old_skill in "${DEPRECATED_SKILLS[@]}"; do
  old_dst="$TARGET_DIR/$old_skill.md"
  if [[ -L "$old_dst" ]]; then
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

  # Backup if exists and is a regular file (not symlink)
  if [[ -f "$dst" && ! -L "$dst" ]]; then
    mv "$dst" "$dst.bak"
    echo -e "${YELLOW}  Backed up: $dst → $dst.bak${NC}"
  elif [[ -L "$dst" ]]; then
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
}

if [[ $# -eq 0 ]]; then
  # Install all skills
  echo "Installing all skills..."
  echo ""
  skill_count=0
  for dir in "$SKILLS_DIR"/*/; do
    [[ -d "$dir" ]] || continue
    skill_name="$(basename "$dir")"
    [[ "$skill_name" == "README"* ]] && continue
    install_skill "$skill_name"
    ((skill_count++))
  done
  echo ""
  echo "Installed $skill_count skill(s) to $TARGET_DIR"
  if [[ "$INSTALL_MODE" == "symlink" ]]; then
    echo ""
    echo "To update: git pull  (symlinks auto-update)"
  else
    echo ""
    echo "To update: re-run with --copy  (files are standalone)"
  fi
  echo "To uninstall: rm $TARGET_DIR/{skill-name}.md"

elif [[ -d "$SKILLS_DIR/$1" ]]; then
  install_skill "$1"

else
  echo -e "${RED}Unknown skill: $1${NC}"
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
