#!/usr/bin/env bash
# erwin Codex skills installer
#
# Installs Codex-native skills (directories containing SKILL.md) into
# ~/.codex/skills/ by symlink or copy.
#
# Usage:
#   ./install-codex.sh
#   ./install-codex.sh harness
#   ./install-codex.sh --copy harness
#   ./install-codex.sh --copy --target /tmp/codex-skills harness

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILLS_DIR="$SCRIPT_DIR/skills"
TARGET_DIR="$HOME/.codex/skills"
INSTALL_MODE="symlink"

GREEN='\033[0;32m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

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
      shift 2
      ;;
    --help|-h)
      echo "erwin Codex skills installer"
      echo ""
      echo "Usage:"
      echo "  ./install-codex.sh                    Install all Codex-native skills"
      echo "  ./install-codex.sh harness            Install one skill"
      echo "  ./install-codex.sh --copy harness     Copy instead of symlink"
      echo "  ./install-codex.sh -t DIR harness     Install to custom directory"
      exit 0
      ;;
    *)
      POSITIONAL_ARGS+=("$1")
      shift
      ;;
  esac
done
set -- "${POSITIONAL_ARGS[@]+"${POSITIONAL_ARGS[@]}"}"

is_codex_skill() {
  local dir="$1"
  [[ -f "$dir/SKILL.md" ]]
}

install_codex_skill() {
  local name="$1"
  local skill_dir="$SKILLS_DIR/$name"
  if [[ ! -d "$skill_dir" ]] || ! is_codex_skill "$skill_dir"; then
    echo -e "${RED}✗ Codex skill not found: $name${NC}"
    return 1
  fi

  local dst_dir="$TARGET_DIR/$name"
  mkdir -p "$dst_dir"

  for file in "$skill_dir"/*; do
    [[ -e "$file" ]] || continue
    local base
    base="$(basename "$file")"
    local dst="$dst_dir/$base"
    [[ -f "$dst" || -L "$dst" || -d "$dst" ]] && rm -rf "$dst"
    if [[ "$INSTALL_MODE" == "copy" ]]; then
      if [[ -d "$file" ]]; then
        cp -R "$file" "$dst"
      else
        cp "$file" "$dst"
      fi
    else
      ln -sf "$file" "$dst"
    fi
  done

  echo -e "${GREEN}✅ /$name${NC} → ${CYAN}$dst_dir${NC}"
}

mkdir -p "$TARGET_DIR"

if [[ $# -eq 0 ]]; then
  count=0
  for dir in "$SKILLS_DIR"/*/; do
    [[ -d "$dir" ]] || continue
    name="$(basename "$dir")"
    if is_codex_skill "$dir"; then
      install_codex_skill "$name"
      count=$((count + 1))
    fi
  done
  echo "Installed $count Codex skill(s) to $TARGET_DIR"
  exit 0
fi

if [[ -d "$SKILLS_DIR/$1" ]] && is_codex_skill "$SKILLS_DIR/$1"; then
  install_codex_skill "$1"
  exit 0
fi

echo -e "${RED}Unknown Codex skill: $1${NC}"
echo ""
echo "Available Codex-native skills:"
for dir in "$SKILLS_DIR"/*/; do
  [[ -d "$dir" ]] || continue
  name="$(basename "$dir")"
  if is_codex_skill "$dir"; then
    echo "  $name"
  fi
done
exit 1
