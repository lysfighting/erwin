# Daily Radar Notes

`docs/radar/daily/` stores the durable daily AI/agent radar stream.

This stream is the short-cycle intake layer for authoritative AI/agent field signal. It is not the final interpretation layer by itself; weekly synthesis and the broader knowledge system decide what should become durable understanding or proposal pressure.

## Rules

- Use one dated file per run date: `YYYY-MM-DD.md`.
- Keep `latest.md` identical to the newest dated file.
- Cover the recent 7-day window unless the automation prompt says otherwise.
- Prefer first-party sources and include direct links plus absolute publish dates when available.
- Keep GitHub Trending or other ecosystem momentum clearly marked as secondary evidence.
- When no major launch appears, still write one durable takeaway instead of an empty update.
- After updating this stream, run `python3 scripts/build-knowledge-surface.py` so the visual, text, and machine-readable knowledge surfaces stay aligned.

## Relationship To Weekly Synthesis

Daily notes capture signal quickly.
Weekly synthesis decides what changed, what matters, and what should feed `feedback -> knowledge -> proposals`.
