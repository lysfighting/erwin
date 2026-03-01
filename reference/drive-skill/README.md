# Drive Skill

> **Author: KC**

One-click setup that transforms [Claude Code](https://claude.com/claude-code) into a **high-drive autonomous agent** with hierarchical team orchestration, voice notifications, and a desktop pet companion.

```bash
git clone https://github.com/Zzzzz-0/drive-skill.git
cd drive-skill
./install.sh
```

That's it. Your Claude Code now has superpowers.

---

## What You Get

### `/drive` — Hierarchical Swarm Execution

The core skill. Turns Claude into a **Team Lead** that:

- Creates hierarchical teams (Team Lead > Sub-Team-Leads > Inner subagents)
- Decomposes work into parallel waves with critical path analysis
- Spawns and orchestrates in-process teammates visible via `Shift+Down`
- Dynamically restructures the hierarchy (promote, split, demote, reassign)
- Broadcasts every dispatch, return, and hierarchy change
- Self-reviews adversarially before marking anything complete
- Never stops until the mission is verified complete

```
/drive implement the authentication module with OAuth2 + JWT
/drive -body refactor the database layer    # with desktop pet
```

### `/self-drive` — Perpetual Improvement Loop

An infinite loop that:
1. Executes a mission via full `/drive`
2. Analyzes the project with 3 parallel scouts
3. Proposes 3 ranked improvements via Slack (or terminal)
4. Waits for your reply — you direct from your phone
5. Executes your choice, loops back

```
/self-drive                          # analyze project, propose tasks
/self-drive refactor the auth module # do this first, then loop
```

### `/architect` — System Architecture Design

A rigorous 6-phase architecture design process:
1. Understand the problem
2. Exhaust the design space (research, module decomposition)
3. Define optimization goals
4. Debate & compare paths
5. Multi-round adversarial review
6. Final presentation with diagrams

```
/architect design a real-time collaboration system
```

### Notification System

Every Claude Code event triggers:
- **macOS notification banner** — project name + contextual summary
- **Voice announcement** — bundled sound effects, no API key needed
- **Desktop pet reaction** — Body Zero animates based on event type

### In-Process Teammate Mode

All teammates are visible in your terminal:
- `Shift+Down` — cycle through active teammates
- See STLs spawning subagents in real time
- TTS announces when teammates need attention
- Your terminal becomes a live war room

---

## Repository Structure

```
drive-skill/
├── install.sh                      # One-click installer
├── commands/
│   ├── drive.md                    # Hierarchical swarm execution
│   ├── self-drive.md               # Perpetual improvement loop
│   └── architect.md                # System architecture design
├── hooks/
│   └── notify.sh                   # macOS banners + voice + Body Zero
├── sounds/                         # Pre-generated TTS sound effects
│   ├── permission-required.mp3
│   ├── waiting-for-input.mp3
│   ├── needs-attention.mp3
│   ├── task-complete.mp3
│   └── teammate-idle.mp3
├── pantheon/                       # 众神殿 — 1000+ greatest minds
│   ├── 01_古代.md ... 07_20世纪下半叶至当代.md
│   └── greatest_minds.md
├── messaging/
│   ├── slack-send.py               # Send messages via Slack Bot API
│   ├── slack-listen.py             # Slack Socket Mode WebSocket listener
│   ├── slack-poll.py               # Slack API polling fallback
│   ├── telegram-send.py            # Send via Telegram Bot API
│   └── telegram-poll.py            # Telegram long-poll listener
├── config/
│   ├── settings.json               # Claude Code settings template
│   └── self-drive.json.example     # Slack/Telegram token template
├── CHANGELOG.md
└── README.md
```

All files are installed as **symlinks** — edit them in this repo, commit, and push. Changes take effect immediately.

---

## Installation Details

### What `install.sh` Does

| Step | Action | Reversible? |
|------|--------|-------------|
| 1 | Creates `~/.claude/commands/` and `~/.claude/sounds/` | Yes |
| 2 | Symlinks `commands/*.md` → `~/.claude/commands/` | Yes (rm symlink) |
| 2b | Symlinks `pantheon/*.md` → `~/.claude/pantheon/` | Yes (rm symlink) |
| 3 | Symlinks `hooks/notify.sh` → `~/.claude/` | Yes (rm symlink) |
| 3b | Symlinks `sounds/*.mp3` → `~/.claude/sounds/` | Yes (rm symlink) |
| 4 | Symlinks `messaging/*.py` → `~/.claude/` | Yes (rm symlink) |
| 5 | Merges `config/settings.json` into `~/.claude/settings.json` | Manual revert |
| 6 | Adds shell aliases to `~/.zshrc` or `~/.bashrc` | Manual revert |

**Existing files are backed up** to `*.bak` before being replaced. Settings are merged, not overwritten.

### Prerequisites

- **Claude Code** — [install guide](https://docs.anthropic.com/en/docs/claude-code)
- **macOS** — for TTS (afplay) and notification banners (osascript). Core skills work on any platform.
- **Python 3** — for messaging scripts

### Uninstall

```bash
# Remove symlinks
rm ~/.claude/commands/{drive,self-drive,architect}.md
rm ~/.claude/notify.sh
rm -rf ~/.claude/sounds/ ~/.claude/pantheon/
rm ~/.claude/{slack-send,slack-listen,slack-poll,telegram-send,telegram-poll}.py

# Remove shell additions (edit manually)
# Look for the "Drive Skill Configuration" block in ~/.zshrc

# Remove settings additions (edit manually)
# Remove hooks and env vars from ~/.claude/settings.json
```

---

## Optional Setup

### Voice Notifications

Voice notifications work out of the box — 5 pre-generated sound effects are bundled in `sounds/`. No API key needed.

### Self-Drive Slack Integration

Create a Slack app with scopes: `chat:write`, `channels:history`, `channels:read`, `connections:write`, `im:history`, `im:read`, `im:write`.

```bash
vim ~/.claude/self-drive.json  # add your tokens (created by install.sh)
```

### Body Zero Desktop Pet

Body Zero is a separate Electron app that reacts to drive events on screen. Install alongside:

```bash
# Clone and install (separate repo)
cd ~/Desktop/building
git clone <body-zero-repo-url> body_zero
cd body_zero && npm install

# Auto-launch via drive mode
/drive -body <your mission>
```

---

## How Drive Mode Works

```
Phase 0: Deep Briefing
  Research Swarm (parallel scouts) → Mission Contract → User confirms

Phase T: Team Assembly
  Hierarchy design (STLs vs. leaves) → Spawn teammates → Emit topology

Phase 1: Battle Plan
  Wave decomposition → Critical path analysis → Domain assignment → Launch

The Loop: Hierarchical Swarm Execution
  STLs spawn inner subagents → Broadcast protocol → Dynamic restructuring

Mission Complete
  Team shutdown → Final verification → Debrief with stats
```

### Sub-Team-Leads (STLs)

STLs are teammates who own a domain and spawn their own Task subagents:

```
Team Lead (you, visible in terminal)
├── STL-frontend (teammate, spawns 2-3 inner subagents)
├── STL-backend (teammate, spawns 2-3 inner subagents)
├── Tester (leaf teammate)
└── Researcher (leaf teammate)
→ 4 teammate slots, ~10 effective parallel streams
```

### Broadcast Protocol

```
📡 WAVE 1 LAUNCH — Width: 6
   [stl-frontend] → UI components [spawns 3 inner]
   [stl-backend]  → API endpoints [spawns 2 inner]

🗺️ TEAM TOPOLOGY:
   YOU (Team Lead) ─── orchestrating
   ├── 🎖️ stl-frontend ─── [3 inner subagents]
   ├── 🎖️ stl-backend ─── [2 inner subagents]
   └── 🔧 tester ─── leaf worker

🏁 DOMAIN COMPLETE: [stl-frontend] ← UI fully delivered
📊 SWARM STATUS: 8/12 tasks | Wave 2/4 | On track
```

---

## Version History

See [CHANGELOG.md](CHANGELOG.md).

## License

MIT
