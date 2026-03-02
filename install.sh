#!/usr/bin/env bash
# erwin skills installer
# Symlinks erwin skills to ~/.claude/commands/ for global availability
# Updates automatically via git pull (symlinks reflect the live files)
#
# Usage:
#   ./install.sh              # install all skills
#   ./install.sh new-project  # install one skill by name

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

install_skill() {
  local name="$1"
  local src="$SKILLS_DIR/$name/$name.md"
  local dst="$TARGET_DIR/$name.md"

  if [[ ! -d "$SKILLS_DIR/$name" ]]; then
    echo -e "${RED}✗ Skill not found: $name${NC}"
    echo "  Available: $(ls "$SKILLS_DIR" | grep -v README | tr '\n' ' ')"
    return 1
  fi

  if [[ ! -f "$src" ]]; then
    echo -e "${RED}✗ Skill file missing: $src${NC}"
    return 1
  fi

  # Backup if exists and is NOT already a symlink to this exact file
  if [[ -f "$dst" && ! -L "$dst" ]]; then
    mv "$dst" "$dst.bak"
    echo -e "${YELLOW}  Backed up: $dst → $dst.bak${NC}"
  elif [[ -L "$dst" ]]; then
    # Already a symlink — just replace it (may point to old location)
    rm "$dst"
  fi

  ln -sf "$src" "$dst"
  echo -e "${GREEN}✅ /$name${NC} → $src"
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
