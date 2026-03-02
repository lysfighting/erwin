# Drive Skill

> **Author: KC**

One-click setup that transforms [Claude Code](https://claude.com/claude-code) into a **high-drive autonomous agent** with hierarchical team orchestration, cognitive coaching, and voice notifications.

```bash
git clone git@github.com:flair-home-stylist/drive_skill.git
cd drive_skill
./install.sh
```

That's it. Your Claude Code now has superpowers.

---

## What You Get

### `/drive` — Hierarchical Swarm Execution

The core skill. Turns Claude into a **Team Lead** operating as a **超级智者 (Super-Sage)** — not a compliant executor, but an intellect with independent judgment built on six principles: Critical Thinking, Creative Thinking (举一反三), Self-Debating, Intellectual Honesty, Simplicity, and 知行合一.

Key capabilities:

- Creates hierarchical teams (Team Lead > Sub-Team-Leads > Inner subagents)
- Decomposes work into parallel waves with critical path analysis
- Three concurrency patterns: Fork-Join, Pipeline, and Speculative execution
- Spawns and orchestrates in-process teammates visible via `Shift+Down`
- **Three mandatory Challenge Gates**: Essence Challenge (本质挑战), Direction Challenge (方向挑战), Contract Challenge (契约挑战)
- **Proactive Challenge Protocol** — challenges suboptimal directions with structured rebuttals, not rubber-stamping
- **Phase Omega (本质洞察)** — first-principles decomposition before any action: strip to the bone, find decided vs. open decisions, measure blast radius
- Broadcasts every dispatch, return, and hierarchy change
- Self-reviews adversarially before marking anything complete
- Never stops until the mission is verified complete

```
/drive implement the authentication module with OAuth2 + JWT
/drive refactor the database layer
```

### `/self-drive` — Perpetual Improvement Loop

An infinite loop ("永动机") that:
1. Executes a mission via full `/drive` (all phases, no shortcuts)
2. Analyzes the project with 5 parallel scouts (quality, architecture, features, external research, ecosystem)
3. Applies the **Sage's Task Selection Protocol** — center of gravity analysis, root cause over symptom, anti-laziness filter
4. Proposes 3 ranked improvements via Slack DM
5. Waits for your reply — you direct from your phone (no timeout)
6. Executes your choice (or any custom task you type), loops back

```
/self-drive                          # analyze project, propose tasks
/self-drive refactor the auth module # do this first, then loop
```

Reply `1`, `2`, `3` to pick a proposal, type any text for a custom task, or `stop` to end.

### `/architect` — System Architecture Design

A rigorous 6-phase architecture design process guided by the Sage principles:
1. **Understand the problem** — first-principles interrogation, not checklists
2. **Exhaust the design space** — parallel research swarm + 众神殿 consultation
3. **Define optimization goals** — explicit trade-off dimensions
4. **Debate & compare paths** — structured adversarial comparison
5. **Multi-round adversarial review** — attack your own design
6. **Final presentation** — diagrams, decision records, implementation roadmap

```
/architect design a real-time collaboration system
```

### `/improve-user` — Cognitive Coaching

A **认知教练 (Cognitive Coach)** that analyzes your collaboration patterns and delivers a personalized cognitive upgrade artifact:

1. **Archaeological Dig** — mines memory files, conversation history, git commits, code patterns, and drive session artifacts with a parallel research swarm
2. **Pattern Extraction** — identifies blind spots across 6 dimensions: communication, decision-making, knowledge topology, work habits, thinking patterns, and AI collaboration anti-patterns
3. **Diagnosis** — synthesizes findings into top 2 high-leverage insights + 1 strength-with-shadow-side contradiction
4. **Deep Research** — WebSearch-backed knowledge gap analysis with primary sources
5. **Artifact** — writes a 2500-3500 word rereadable essay to `.claude/drive/cognitive-upgrade-{date}.md`

```
/improve-user
/improve-user focus on my decision-making patterns
```

### Notification System

Every Claude Code event triggers:
- **macOS notification banner** — project name + contextual summary
- **Bundled sound effects** — 5 pre-generated voice alerts (permission-required, waiting-for-input, needs-attention, task-complete, teammate-idle), no API key needed
- **Fish Audio TTS** (optional) — custom voice announcements via `speak.sh` with local caching

### Slack Question Relay

When you don't respond to an `AskUserQuestion` in the terminal:
- The question is automatically forwarded to your Slack DM with formatted options
- You reply on Slack → answer is injected back into the terminal
- In self-drive mode, questions are relayed immediately (no wait)
- Only main-agent questions are relayed; teammate questions are skipped

### In-Process Teammate Mode

All teammates are visible in your terminal:
- `Shift+Down` — cycle through active teammates
- See STLs spawning subagents in real time
- Sound alerts when teammates need attention
- Your terminal becomes a live war room

---

## Repository Structure

```
drive-skill/
├── install.sh                      # One-click installer
├── commands/
│   ├── drive.md                    # Hierarchical swarm execution
│   ├── self-drive.md               # Perpetual improvement loop
│   ├── architect.md                # System architecture design
│   └── improve-user.md             # Cognitive coaching
├── hooks/
│   ├── notify.sh                   # macOS banners + bundled sound effects
│   ├── question-relay-hook.sh      # AskUserQuestion → Slack relay lifecycle
│   └── speak.sh                    # Fish Audio TTS (optional, requires API key)
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
│   ├── slack-question-relay.py     # AskUserQuestion → Slack → terminal injection
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
| 1 | Creates `~/.claude/commands/`, `~/.claude/sounds/`, `~/.claude/pantheon/` | Yes |
| 2 | Symlinks `commands/*.md` → `~/.claude/commands/` | Yes (rm symlink) |
| 2b | Symlinks `pantheon/*.md` → `~/.claude/pantheon/` | Yes (rm symlink) |
| 3 | Symlinks `hooks/*.sh` → `~/.claude/` | Yes (rm symlink) |
| 3b | Symlinks `sounds/*.mp3` → `~/.claude/sounds/` | Yes (rm symlink) |
| 4 | Symlinks `messaging/*.py` → `~/.claude/` | Yes (rm symlink) |
| 5 | Merges `config/settings.json` into `~/.claude/settings.json` | Manual revert |
| 6 | Adds shell aliases to `~/.zshrc` or `~/.bashrc` | Manual revert |

**Existing files are backed up** to `*.bak` before being replaced. Settings are merged, not overwritten.

### Prerequisites

- **Claude Code** — [install guide](https://docs.anthropic.com/en/docs/claude-code)
- **macOS** — for afplay and notification banners (osascript). Core skills work on any platform.
- **Python 3** — for messaging scripts

### Uninstall

```bash
# Remove symlinks
rm ~/.claude/commands/{drive,self-drive,architect,improve-user}.md
rm ~/.claude/{notify,question-relay-hook,speak}.sh
rm -rf ~/.claude/sounds/ ~/.claude/pantheon/
rm ~/.claude/{slack-send,slack-listen,slack-poll,slack-question-relay,telegram-send,telegram-poll}.py

# Remove shell additions (edit manually)
# Look for the "Drive Skill Configuration" block in ~/.zshrc

# Remove settings additions (edit manually)
# Remove hooks and env vars from ~/.claude/settings.json
```

---

## Optional Setup

### Voice Notifications

Voice notifications work out of the box — 5 pre-generated sound effects are bundled in `sounds/`. No API key needed.

For custom TTS voice announcements, configure Fish Audio:
```bash
export FISH_AUDIO_API_KEY="your-key-here"  # Get one at https://fish.audio
```

### Slack Integration

Create a Slack app with scopes: `chat:write`, `channels:history`, `channels:read`, `connections:write`, `im:history`, `im:read`, `im:write`.

```bash
vim ~/.claude/self-drive.json  # add your tokens (created by install.sh)
```

This enables both `/self-drive` proposals and the question relay feature.

---

## How Drive Mode Works

```
Phase Ω: 本质洞察 (Essence Insight)
  First-principles decomposition → Bone vs. flesh → Decided vs. open → Blast radius

Phase 0: Deep Briefing
  Research Swarm (parallel scouts) → Three Challenge Gates → Mission Contract → User confirms

Phase T: Team Assembly
  Hierarchy design (STLs vs. leaves) → Spawn teammates → Emit topology

Phase 1: Battle Plan
  Wave decomposition → Critical path analysis → Domain assignment → Launch

The Loop: Hierarchical Swarm Execution
  STLs spawn inner subagents → Broadcast protocol → Dynamic restructuring

Mission Complete
  Team shutdown → Self-adversarial review → Final verification → Debrief with stats
```

### Sub-Team-Leads (STLs)

STLs are teammates who own a domain and spawn their own inner subagents:

```
Team Lead (you, visible in terminal)
├── STL-frontend (teammate, spawns 2-3 inner subagents)
├── STL-backend (teammate, spawns 2-3 inner subagents)
├── Tester (leaf teammate)
└── Researcher (leaf teammate)
→ 4 teammate slots, ~10 effective parallel streams
```

| Mission Size | Team Topology | Rationale |
|---|---|---|
| Small (≤5 tasks) | Flat: You + 1-2 leaf teammates | Hierarchy overhead > benefit |
| Medium (6-12 tasks) | Shallow: You + 1-2 STLs + 1 leaf | STLs own domains |
| Large (13+ tasks) | Deep: You + 2-3 STLs (each with 2-3 subagents) | Pure orchestration |

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
