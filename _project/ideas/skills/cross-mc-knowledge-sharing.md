# Cross-MC Knowledge Sharing

> Stage: Fruit — implemented in v3.8.5 as `/team sync` + `/team share`
> Session: 2026-03-12-AJ

## Problem

Multiple MCs running in parallel worktrees are information silos. Related issues touch the same code, have overlapping requirements, but each agent works in complete isolation. Result: conflicting changes, duplicated discoveries, missed synergies.

## Design

### L1: Passive Read (team-drive.md)

`/team-drive` Phase 0 adds a step: read issue comments via `gh issue view #N --comments --json comments`.

- Filter for `🔗 Cross-MC` prefixed comments
- Inject as advisory context into Contract: "来自兄弟 MC，仅供参考"
- No new command needed — automatic on drive startup

### L2: Active Push (team.md → `/team share`)

```
/team share #241    → post structured knowledge to #241 (requires user confirmation)
/team share         → scan sibling MCs for overlap, suggest targets
```

Comment format:
```markdown
🔗 **Cross-MC** from **#238** (advisory)
**Context**: {what was being done}
**Finding**: {the knowledge}
**Files touched**: {relevant files}
**Confidence**: {high/medium/low — with basis}
```

Anti-plague rules:
1. Human gateway: AskUserQuestion before every post
2. No transitive propagation: only first-hand findings, never re-share
3. Source attribution: structured format with confidence

### `/team sync` (team.md)

Independent rebase from base branch. Currently buried as `/team-ship sync` subcommand.
- Keep `/team-ship sync` working (backward compat)
- Add `/team sync` as standalone entry in team.md routing

### L3: Bidirectional Negotiation (FUTURE)

Two agents discuss via comments before making conflicting changes. Deferred — requires real-time coordination patterns that don't exist yet.

## Files to Modify

| File | Change |
|------|--------|
| `team-drive.md` | Phase 0: read issue comments, filter Cross-MC |
| `team.md` | Route: `/team sync`, `/team share` |
| `team.md` | Operation Share: scan + post + human gateway |
| `team.md` | Operation Sync: rebase logic (extract from team-ship) |

## Reference

- Claude Code team comms: `~/.claude/teams/{name}/inboxes/` — file mailbox + polling + autoWake
- GitHub issue comments as cross-MC mailbox: persistent, async, visible to humans, no new infra
