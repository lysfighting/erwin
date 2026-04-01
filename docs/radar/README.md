# AI/Agent Radar

`docs/radar/` is the lightweight knowledge base behind the `erwin` field radar.

Its job is not to mirror the whole internet.
Its job is to keep a tight, high-signal map of what matters for this repository.

## Purpose

- track authoritative AI/agent developments across the field
- keep `Harness` in context as one important theme rather than the entire boundary
- turn recurring automation output into reusable weekly knowledge
- give future sessions a stable starting point without replaying inbox history

## Output Model

- `Daily radar`
  - inbox + repo
  - short format
  - up to 5 signals
  - recent window: last 7 days
  - one durable daily file per date, updated in place
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

See [themes.md](/Users/liyasong/Project/collov/erwin/docs/radar/themes.md) for the current taxonomy and keywords.

## Files

- [sources.md](/Users/liyasong/Project/collov/erwin/docs/radar/sources.md) — publisher inventory
- [themes.md](/Users/liyasong/Project/collov/erwin/docs/radar/themes.md) — active themes and search language
- [daily/](/Users/liyasong/Project/collov/erwin/docs/radar/daily) — durable daily radar notes
- [weekly/](/Users/liyasong/Project/collov/erwin/docs/radar/weekly) — durable weekly synthesis archive

## Working Rules

- keep daily notes concise and update one file per date rather than creating duplicates
- prefer durable daily notes plus one weekly synthesis over inbox-only history
- always include links and dates when the underlying source provides them
- prefer direct official pages over commentary
- state clearly when a tool trend is interesting but not strategically relevant
