#!/bin/bash
set -e

# ═══════════════════════════════════════════════════════════════
# Drive Skill — One-Click Claude Code Setup
# ═══════════════════════════════════════════════════════════════
#
# This script transforms your Claude Code into a high-drive
# autonomous agent with:
#   - Hierarchical swarm execution (drive skill)
#   - Perpetual improvement loop (self-drive skill)
#   - System architecture design (architect skill)
#   - Cognitive coaching & user growth (improve-user skill)
#   - 众神殿 (Pantheon) — 1000+ greatest minds as reasoning reference
#   - Voice notifications (bundled sound effects)
#   - macOS notification banners
#   - Slack/Telegram async communication
#   - In-process teammate mode by default
#
# Usage:
#   git clone https://github.com/Zzzzz-0/drive-skill.git
#   cd drive-skill
#   ./install.sh
#
# ═══════════════════════════════════════════════════════════════

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CLAUDE_DIR="$HOME/.claude"
COMMANDS_DIR="$CLAUDE_DIR/commands"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color
BOLD='\033[1m'

print_header() {
  echo ""
  echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
  echo -e "${BOLD}  $1${NC}"
  echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
  echo ""
}

print_step() {
  echo -e "  ${GREEN}[+]${NC} $1"
}

print_skip() {
  echo -e "  ${YELLOW}[~]${NC} $1"
}

print_warn() {
  echo -e "  ${YELLOW}[!]${NC} $1"
}

print_error() {
  echo -e "  ${RED}[x]${NC} $1"
}

print_info() {
  echo -e "  ${BLUE}[i]${NC} $1"
}

# ─── Pre-flight checks ───

print_header "Drive Skill Installer"

# Check Claude Code is installed
if ! command -v claude &>/dev/null; then
  print_error "Claude Code (claude) not found in PATH."
  print_info "Install it first: https://docs.anthropic.com/en/docs/claude-code"
  exit 1
fi
print_step "Claude Code found: $(which claude)"

# Check macOS (for afplay, osascript)
if [[ "$(uname)" != "Darwin" ]]; then
  print_warn "This installer is optimized for macOS."
  print_warn "TTS (afplay) and macOS banners (osascript) will not work on other platforms."
  print_warn "The core skills (drive, self-drive, architect) will still work."
  echo ""
  read -p "  Continue anyway? [y/N] " -n 1 -r
  echo ""
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 0
  fi
fi

# Check Python3
if ! command -v python3 &>/dev/null; then
  print_error "python3 not found. Required for messaging scripts."
  exit 1
fi
print_step "Python3 found: $(python3 --version)"

# ─── Create directories ───

print_header "Step 1: Directory Setup"

mkdir -p "$COMMANDS_DIR"
print_step "Created $COMMANDS_DIR"

# ─── Install commands (symlinks) ───

print_header "Step 2: Install Skills"

install_command() {
  local src="$SCRIPT_DIR/commands/$1"
  local dst="$COMMANDS_DIR/$1"

  if [ -L "$dst" ]; then
    local current_target=$(readlink "$dst")
    if [ "$current_target" = "$src" ]; then
      print_skip "$1 — already linked"
      return
    fi
    rm "$dst"
  elif [ -f "$dst" ]; then
    print_warn "$1 — existing file found, backing up to $dst.bak"
    mv "$dst" "$dst.bak"
  fi

  ln -s "$src" "$dst"
  print_step "$1 — installed (symlink)"
}

install_command "drive.md"
install_command "self-drive.md"
install_command "architect.md"
install_command "improve-user.md"
install_command "reasoning-toolkit.md"

# ─── Install Pantheon (众神殿) ───

print_header "Step 2b: Install Pantheon (众神殿)"

PANTHEON_SRC="$SCRIPT_DIR/pantheon"
PANTHEON_DST="$CLAUDE_DIR/pantheon"

if [ -d "$PANTHEON_SRC" ]; then
  mkdir -p "$PANTHEON_DST"
  print_step "Created $PANTHEON_DST"

  for md_file in "$PANTHEON_SRC"/*.md; do
    fname=$(basename "$md_file")
    dst_file="$PANTHEON_DST/$fname"

    if [ -L "$dst_file" ]; then
      current_target=$(readlink "$dst_file")
      if [ "$current_target" = "$md_file" ]; then
        print_skip "$fname — already linked"
        continue
      fi
      rm "$dst_file"
    elif [ -f "$dst_file" ]; then
      print_warn "$fname — existing file found, backing up to $dst_file.bak"
      mv "$dst_file" "$dst_file.bak"
    fi

    ln -s "$md_file" "$dst_file"
    print_step "$fname — installed (symlink)"
  done
else
  print_skip "Pantheon directory not found in repo — skipping"
fi

# ─── Install hooks ───

print_header "Step 3: Install Hooks"

install_hook() {
  local src="$SCRIPT_DIR/hooks/$1"
  local dst="$CLAUDE_DIR/$1"

  if [ -L "$dst" ]; then
    local current_target=$(readlink "$dst")
    if [ "$current_target" = "$src" ]; then
      print_skip "$1 — already linked"
      return
    fi
    rm "$dst"
  elif [ -f "$dst" ]; then
    print_warn "$1 — existing file found, backing up to $dst.bak"
    mv "$dst" "$dst.bak"
  fi

  ln -s "$src" "$dst"
  chmod +x "$src"
  print_step "$1 — installed (symlink)"
}

install_hook "notify.sh"
install_hook "question-relay-hook.sh"

# ─── Install sound files ───

print_header "Step 3b: Install Sound Effects"

SOUNDS_SRC="$SCRIPT_DIR/sounds"
SOUNDS_DST="$CLAUDE_DIR/sounds"

if [ -d "$SOUNDS_SRC" ]; then
  if [ -L "$SOUNDS_DST" ]; then
    current_target=$(readlink "$SOUNDS_DST")
    if [ "$current_target" = "$SOUNDS_SRC" ]; then
      print_skip "sounds/ — already linked"
    else
      rm "$SOUNDS_DST"
      ln -s "$SOUNDS_SRC" "$SOUNDS_DST"
      print_step "sounds/ — relinked to $SOUNDS_SRC"
    fi
  elif [ -d "$SOUNDS_DST" ]; then
    # Migrate from per-file symlinks to directory symlink
    rm -rf "$SOUNDS_DST"
    ln -s "$SOUNDS_SRC" "$SOUNDS_DST"
    print_step "sounds/ — migrated to directory symlink"
  else
    ln -s "$SOUNDS_SRC" "$SOUNDS_DST"
    print_step "sounds/ — installed (symlink)"
  fi
else
  print_skip "Sounds directory not found in repo — skipping"
fi

# ─── Install messaging scripts ───

print_header "Step 4: Install Messaging Scripts"

install_messaging() {
  local src="$SCRIPT_DIR/messaging/$1"
  local dst="$CLAUDE_DIR/$1"

  if [ -L "$dst" ]; then
    local current_target=$(readlink "$dst")
    if [ "$current_target" = "$src" ]; then
      print_skip "$1 — already linked"
      return
    fi
    rm "$dst"
  elif [ -f "$dst" ]; then
    print_warn "$1 — existing file found, backing up to $dst.bak"
    mv "$dst" "$dst.bak"
  fi

  ln -s "$src" "$dst"
  chmod +x "$src"
  print_step "$1 — installed (symlink)"
}

install_messaging "slack-send.py"
install_messaging "slack-listen.py"
install_messaging "slack-poll.py"
install_messaging "telegram-send.py"
install_messaging "telegram-poll.py"
install_messaging "slack-question-relay.py"

# ─── Install self-drive config template ───

if [ ! -f "$CLAUDE_DIR/self-drive.json" ]; then
  cp "$SCRIPT_DIR/config/self-drive.json.example" "$CLAUDE_DIR/self-drive.json"
  print_step "self-drive.json — created from template (edit with your Slack/Telegram tokens)"
else
  print_skip "self-drive.json — already exists, keeping your config"
fi

# ─── Configure settings.json ───

print_header "Step 5: Configure Claude Settings"

SETTINGS_FILE="$CLAUDE_DIR/settings.json"

if [ ! -f "$SETTINGS_FILE" ]; then
  # No existing settings — install ours directly
  cp "$SCRIPT_DIR/config/settings.json" "$SETTINGS_FILE"
  print_step "settings.json — installed (fresh)"
else
  # Existing settings — merge carefully
  print_info "Existing settings.json found. Merging..."

  # Use Python to merge JSON
  python3 << 'MERGE_SCRIPT'
import json, os, sys

settings_path = os.path.expanduser("~/.claude/settings.json")
template_path = os.path.join(os.environ.get("SCRIPT_DIR", "."), "config", "settings.json")

with open(settings_path) as f:
    existing = json.load(f)
with open(template_path) as f:
    template = json.load(f)

changed = False

# Ensure env vars
if "env" not in existing:
    existing["env"] = {}
if existing["env"].get("CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS") != "1":
    existing["env"]["CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS"] = "1"
    changed = True

# Ensure hooks exist
if "hooks" not in existing:
    existing["hooks"] = {}

for hook_name in ["Notification", "Stop", "TeammateIdle"]:
    if hook_name not in existing["hooks"]:
        existing["hooks"][hook_name] = template["hooks"][hook_name]
        changed = True
    else:
        # Check if notify.sh is already in the hooks
        has_notify = any(
            h.get("command", "").endswith("notify.sh")
            for entry in existing["hooks"][hook_name]
            for h in entry.get("hooks", [])
        )
        if not has_notify:
            existing["hooks"][hook_name].extend(template["hooks"][hook_name])
            changed = True

# Ensure question relay hooks (PreToolUse / PostToolUse for AskUserQuestion)
for hook_name in ["PreToolUse", "PostToolUse"]:
    if hook_name not in existing["hooks"]:
        existing["hooks"][hook_name] = template["hooks"][hook_name]
        changed = True
    else:
        has_relay = any(
            h.get("command", "").endswith("question-relay-hook.sh")
            for entry in existing["hooks"][hook_name]
            for h in entry.get("hooks", [])
        )
        if not has_relay:
            existing["hooks"][hook_name].extend(template["hooks"][hook_name])
            changed = True

if changed:
    with open(settings_path, "w") as f:
        json.dump(existing, f, indent=2)
    print("  [\033[0;32m+\033[0m] settings.json — merged (agent teams + hooks added)")
else:
    print("  [\033[1;33m~\033[0m] settings.json — already configured")
MERGE_SCRIPT

  export SCRIPT_DIR
fi

# ─── Shell configuration ───

print_header "Step 6: Shell Configuration"

# Detect shell
SHELL_NAME=$(basename "$SHELL")
if [ "$SHELL_NAME" = "zsh" ]; then
  SHELL_RC="$HOME/.zshrc"
elif [ "$SHELL_NAME" = "bash" ]; then
  SHELL_RC="$HOME/.bashrc"
else
  SHELL_RC="$HOME/.${SHELL_NAME}rc"
fi

NEEDS_SHELL_UPDATE=false
SHELL_ADDITIONS=""

# Check for teammate mode alias
if ! grep -q 'teammate-mode in-process' "$SHELL_RC" 2>/dev/null; then
  SHELL_ADDITIONS+=$'\n# Claude Code: default teammate mode — in-process (visible swarm)\nalias claude="claude --teammate-mode in-process"\n'
  NEEDS_SHELL_UPDATE=true
fi

# Check for max output tokens
if ! grep -q 'CLAUDE_CODE_MAX_OUTPUT_TOKENS' "$SHELL_RC" 2>/dev/null; then
  SHELL_ADDITIONS+=$'\n# Claude Code: maximize output token budget\nexport CLAUDE_CODE_MAX_OUTPUT_TOKENS=1000000\n'
  NEEDS_SHELL_UPDATE=true
fi

if [ "$NEEDS_SHELL_UPDATE" = true ]; then
  echo "" >> "$SHELL_RC"
  echo "# ─── Drive Skill Configuration (added by install.sh) ───" >> "$SHELL_RC"
  echo "$SHELL_ADDITIONS" >> "$SHELL_RC"
  print_step "Added to $SHELL_RC:"
  if echo "$SHELL_ADDITIONS" | grep -q 'teammate-mode'; then
    print_info "  alias claude=\"claude --teammate-mode in-process\""
  fi
  if echo "$SHELL_ADDITIONS" | grep -q 'MAX_OUTPUT_TOKENS'; then
    print_info "  export CLAUDE_CODE_MAX_OUTPUT_TOKENS=1000000"
  fi
  print_warn "Run 'source $SHELL_RC' or restart your terminal to apply."
else
  print_skip "Shell already configured"
fi

# ─── Voice Notifications ───

print_header "Step 7: Voice Notifications"

SOUND_COUNT=$(ls "$SOUNDS_DST"/*.mp3 2>/dev/null | wc -l | tr -d ' ')
if [ "$SOUND_COUNT" -gt 0 ]; then
  print_step "Bundled sound effects installed ($SOUND_COUNT files) — no API key needed"
else
  print_warn "No sound files found. Voice notifications will be silent."
fi

# ─── Optional: Slack Setup ───

print_header "Step 8: Self-Drive Slack Integration (Optional)"

if [ -f "$CLAUDE_DIR/self-drive.json" ]; then
  HAS_SLACK=$(python3 -c "
import json
with open('$CLAUDE_DIR/self-drive.json') as f:
    c = json.load(f)
t = c.get('slack',{}).get('bot_token','')
print('yes' if t and not t.startswith('xoxb-YOUR') else 'no')
" 2>/dev/null)

  if [ "$HAS_SLACK" = "yes" ]; then
    print_step "Slack is configured"
  else
    print_info "Self-Drive uses Slack for async task proposals."
    print_info "Edit ~/.claude/self-drive.json with your Slack bot tokens."
    print_info "See: https://api.slack.com/apps — create a bot with chat:write + channels:history scopes"
  fi
fi

# ─── Summary ───

print_header "Installation Complete"

echo -e "  ${BOLD}Installed skills:${NC}"
echo -e "    ${GREEN}/drive${NC}         — Hierarchical swarm execution (team lead mode)"
echo -e "    ${GREEN}/self-drive${NC}    — Perpetual autonomous improvement loop"
echo -e "    ${GREEN}/architect${NC}     — System architecture design process"
echo -e "    ${GREEN}/improve-user${NC}  — Cognitive coach: blind spots, knowledge gaps, growth"
echo ""
echo -e "  ${BOLD}Installed knowledge:${NC}"
echo -e "    ${GREEN}众神殿${NC}     — 1000+ greatest minds across human civilization"
echo -e "               Used by /drive Deep Reasoning Protocol for inspiration"
echo ""
echo -e "  ${BOLD}Installed hooks:${NC}"
echo -e "    ${GREEN}notify.sh${NC}              — macOS banners + bundled voice"
echo -e "    ${GREEN}question-relay-hook.sh${NC} — Slack fallback for unanswered questions"
echo -e "  ${BOLD}Installed sounds:${NC}"
echo -e "    ${GREEN}5 mp3 files${NC} — pre-generated TTS (no API key needed)"
echo ""
echo -e "  ${BOLD}Installed messaging:${NC}"
echo -e "    ${GREEN}Slack${NC}      — send, listen (Socket Mode), poll (API fallback)"
echo -e "    ${GREEN}Telegram${NC}   — send, poll"
echo ""
echo -e "  ${BOLD}Configured:${NC}"
echo -e "    ${GREEN}Agent Teams${NC}     — enabled (CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1)"
echo -e "    ${GREEN}Teammate Mode${NC}  — in-process by default (visible swarm)"
echo -e "    ${GREEN}Hooks${NC}          — Notification + Stop + TeammateIdle + Question Relay"
echo -e "    ${GREEN}Output Tokens${NC}  — maximized (1,000,000)"
echo ""
echo -e "  ${BOLD}Optional (configure later):${NC}"
if [ "$HAS_SLACK" = "yes" ]; then
  echo -e "    ${GREEN}Slack${NC}          — configured"
else
  echo -e "    ${YELLOW}Slack${NC}          — edit ~/.claude/self-drive.json for self-drive"
fi
echo ""
echo -e "  ${BOLD}Quick start:${NC}"
echo -e "    ${CYAN}claude${NC}                              # launches with teammate mode"
echo -e "    ${CYAN}/drive implement auth module${NC}         # drive mode"
echo -e "    ${CYAN}/self-drive${NC}                          # perpetual improvement loop"
echo -e "    ${CYAN}/architect design a new API layer${NC}    # architecture design"
echo -e "    ${CYAN}/improve-user${NC}                        # cognitive coaching session"
echo ""
echo -e "  ${BOLD}All files are symlinked${NC} — edit them in this repo,"
echo -e "  commit, and push for version control."
echo ""
