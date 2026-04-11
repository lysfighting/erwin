# AI/Agent Radar

`docs/radar/` is the lightweight knowledge base behind the `erwin` field radar.

Its job is not to mirror the whole internet.
Its job is to keep a tight, high-signal map of what matters for this repository.
It is the intake layer of the knowledge system, not the whole system by itself.

## Purpose

- track authoritative AI/agent developments across the field
- keep `Harness` in context as one important theme rather than the entire boundary
- turn recurring automation output into reusable weekly knowledge
- give future sessions a stable starting point without replaying inbox history
- feed the broader `feedback -> knowledge -> proposals` path without pretending every weekly note is already a proposal

## Output Model

- `Daily radar`
  - inbox + durable files
  - short format
  - up to 5 signals
  - recent window: last 7 days
  - one dated note plus `latest.md`
- `Weekly synthesis`
  - inbox + repo
  - one durable synthesis per week
  - captures what changed, what matters, and what we should watch next

## Source Policy

### Tier 1: authoritative publishers

- Anthropic
- OpenAI
- Google
- Moonshot / Kimi
- MiniMax
- Qwen
- DeepSeek

### Tier 2: ecosystem references

- Langfuse
- Phoenix
- Helicone
- OpenTelemetry
- GitHub Trending

Tier 2 sources are supporting evidence only.
They should help interpret a first-party signal, not pull the roadmap.

## Theme Policy

The radar always tracks a small fixed set of themes and updates them slowly.

Current themes:

- `harness`
- `long-running agents`
- `agent evals`
- `tracing/observability`
- `context engineering`
- `coding agents`
- `tool use`
- `orchestration`

See [themes.md](themes.md) for the current taxonomy and keywords.

## Files

- [../index.html](../index.html) — visual reading surface for latest briefings and archive browsing
- [sources.md](sources.md) — publisher inventory
- [themes.md](themes.md) — active themes and search language
- [daily/](daily/) — mirrored daily radar notes in the main working tree
- [trending/](trending/) — mirrored daily GitHub Trending summaries in the main working tree
- [weekly/](weekly/) — durable weekly synthesis archive

## Working Rules

- each recurring run should first persist to a stable host-level path and then mirror the same content into this main working tree
- keep one dated file plus `latest.md` for each recurring stream
- keep the visual reading surface aligned with the latest notes instead of relying on transient automation inbox cards
- after updating a stream, run `python3 scripts/build-knowledge-surface.py` so `docs/index.html`, `docs/knowledge/briefing.md`, and `docs/knowledge/index.json` stay current together
- always include links and dates when the underlying source provides them
- prefer direct official pages over commentary
- keep GitHub Trending in a separate summary block so popularity signals do not masquerade as first-party field direction
- state clearly when a tool trend is interesting but not strategically relevant
