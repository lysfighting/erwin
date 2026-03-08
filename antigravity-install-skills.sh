#!/usr/bin/env bash
# Antigravity Skills Installer
# Symlinks erwin skills from ~/.gemini/antigravity/skills/ into a project's .agents/skills/
#
# Usage:
#   ./antigravity-install-skills.sh                     # Install to current project
#   ./antigravity-install-skills.sh /path/to/project    # Install to specific project
#   ./antigravity-install-skills.sh --list              # List available skills
#   ./antigravity-install-skills.sh --uninstall         # Remove symlinks from current project

set -euo pipefail

SKILLS_SRC="$HOME/.gemini/antigravity/skills"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

usage() {
  echo "Antigravity Skills Installer"
  echo ""
  echo "Usage:"
  echo "  $0                        Symlink all skills to current project"
  echo "  $0 /path/to/project       Symlink all skills to specific project"
  echo "  $0 --list                 List available skills"
  echo "  $0 --uninstall            Remove skill symlinks from current project"
  echo ""
  echo "Skills are installed from: $SKILLS_SRC"
  echo "Target directory: <project>/.agents/skills/"
  exit 0
}

list_skills() {
  echo -e "${CYAN}Available Antigravity skills:${NC}"
  echo ""
  for dir in "$SKILLS_SRC"/*/; do
    [[ -d "$dir" ]] || continue
    skill_name="$(basename "$dir")"
    if [[ -f "$dir/SKILL.md" ]]; then
      desc=$(grep '^description:' "$dir/SKILL.md" | head -1 | sed 's/^description: *//')
      echo -e "  ${GREEN}$skill_name${NC}"
      echo "    $desc"
    fi
  done
  exit 0
}

# Parse args
case "${1:-}" in
  --help|-h) usage ;;
  --list|-l) list_skills ;;
  --uninstall|-u)
    TARGET_DIR="$(pwd)/.agents/skills"
    if [[ -L "$TARGET_DIR" ]]; then
      rm "$TARGET_DIR"
      echo -e "${GREEN}✅ Removed skills symlink from $(pwd)${NC}"
    elif [[ -d "$TARGET_DIR" ]]; then
      # Remove individual skill symlinks
      for link in "$TARGET_DIR"/*/; do
        [[ -L "${link%/}" ]] && rm "${link%/}" && echo -e "${YELLOW}  Removed: $(basename "${link%/}")${NC}"
      done
      echo -e "${GREEN}✅ Cleaned up skill symlinks${NC}"
    else
      echo -e "${YELLOW}No skills installed in $(pwd)${NC}"
    fi
    exit 0
    ;;
  "") TARGET_PROJECT="$(pwd)" ;;
  *) TARGET_PROJECT="$1" ;;
esac

# Validate source
if [[ ! -d "$SKILLS_SRC" ]]; then
  echo -e "${RED}✗ Skills source not found: $SKILLS_SRC${NC}"
  echo "  Run the erwin skill conversion first."
  exit 1
fi

# Create target .agents/skills/ directory
TARGET_DIR="$TARGET_PROJECT/.agents/skills"
mkdir -p "$(dirname "$TARGET_DIR")"

# Option 1: Symlink the entire skills directory (simplest)
if [[ -L "$TARGET_DIR" ]]; then
  echo -e "${YELLOW}Skills already linked: $TARGET_DIR → $(readlink "$TARGET_DIR")${NC}"
  echo "Use --uninstall to remove first."
  exit 0
elif [[ -d "$TARGET_DIR" ]]; then
  echo -e "${YELLOW}⚠ .agents/skills/ already exists as a directory.${NC}"
  echo "Symlinking individual skills instead..."
  # Symlink individual skill folders
  count=0
  for dir in "$SKILLS_SRC"/*/; do
    [[ -d "$dir" ]] || continue
    skill_name="$(basename "$dir")"
    dst="$TARGET_DIR/$skill_name"
    if [[ -e "$dst" ]]; then
      echo -e "${YELLOW}  Skipped: $skill_name (already exists)${NC}"
    else
      ln -sf "$dir" "$dst"
      echo -e "${GREEN}  ✅ $skill_name${NC}"
      ((count++))
    fi
  done
  echo ""
  echo -e "Installed ${GREEN}$count${NC} skill(s) to $TARGET_DIR"
else
  # Symlink the whole directory
  ln -sf "$SKILLS_SRC" "$TARGET_DIR"
  skill_count=$(find "$SKILLS_SRC" -maxdepth 1 -mindepth 1 -type d | wc -l | tr -d ' ')
  echo -e "${GREEN}✅ Linked $skill_count skills to $TARGET_DIR${NC}"
  echo ""
  echo -e "${CYAN}Skills available:${NC}"
  for dir in "$SKILLS_SRC"/*/; do
    [[ -d "$dir" ]] || continue
    echo "  $(basename "$dir")"
  done
fi

echo ""
echo "To update: skills auto-update (symlinked to $SKILLS_SRC)"
echo "To uninstall: $0 --uninstall"
