# GitHub Trending Notes

`docs/radar/trending/` stores the durable daily summary of the default GitHub Trending board.

This stream is a popularity and momentum snapshot. It can reveal ecosystem demand around agents, coding workflows, MCP, memory, observability, or developer tooling, but it should not outrank first-party technical releases.

## Rules

- Use one dated file per run date: `YYYY-MM-DD.md`.
- Keep `latest.md` identical to the newest dated file.
- Use the default global GitHub Trending page without extra query filters.
- Include direct repository links and the capture date.
- Preserve the board shape without forcing every entry into an AI/agent narrative.
- Explicitly separate agent-relevant projects from broader non-agent or applied-AI entries.
- Treat this as Tier 2 evidence that can help explain Tier 1 signals, not as roadmap authority.
- After updating this stream, run `python3 scripts/build-knowledge-surface.py` so the visual, text, and machine-readable knowledge surfaces stay aligned.

## Relationship To Daily Radar

Trending notes track open-source attention.
Daily radar may cite them as secondary context, but first-party publisher signals remain the primary field direction.
