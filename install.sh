#!/usr/bin/env bash
# erwin skills installer
# Symlinks erwin skills to ~/.claude/commands/ for global availability
# Updates automatically via git pull (symlinks reflect the live files)
#
# Usage:
#   ./install.sh              # install all skills
#   ./install.sh erwin-init   # install one skill by name

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILLS_DIR="$SCRIPT_DIR/skills"
TARGET_DIR="$HOME/.claude/commands"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "erwin skills installer"
echo "======================"

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

  # Backup if exists and is NOT already a symlink
  if [[ -f "$dst" && ! -L "$dst" ]]; then
    mv "$dst" "$dst.bak"
    echo -e "${YELLOW}  Backed up: $dst → $dst.bak${NC}"
  elif [[ -L "$dst" ]]; then
    rm "$dst"
  fi

  ln -sf "$src" "$dst"
  echo -e "${GREEN}✅ /$filename${NC} → $src"
}

install_skill() {
  local name="$1"
  local skill_dir="$SKILLS_DIR/$name"

  if [[ ! -d "$skill_dir" ]]; then
    echo -e "${RED}✗ Skill not found: $name${NC}"
    echo "  Available: $(ls "$SKILLS_DIR" | grep -v README | tr '\n' ' ')"
    return 1
  fi

  # Find all .md files that are skills (exclude README, CHANGELOG)
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
  echo ""
  echo "To update: git pull  (symlinks auto-update)"
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
