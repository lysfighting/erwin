# Feedback: Teamwork v2 Implementation

## What Worked
- Architecture doc (docs/teamwork-v2-architecture.md) was detailed enough to implement directly — minimal decisions needed during execution
- Each skill follows a consistent pattern (prerequisites → route → execute → output) making implementation predictable
- install.sh's existing design handles multi-file skills correctly — no changes needed

## What Didn't Work
- team-drive.md filename collision with v1 — had to Read before Write (v1 had same filename but different content)
- Stale v1 symlinks not cleaned by install.sh — had to manually remove them after install

## Surprises
- install.sh excludes README.md and CHANGELOG.md by basename check — clean design, no special handling needed
- The skill files are substantial (6-11KB each) — they're full prompt engineering documents, not thin scripts

## Lessons for Future Missions
- When replacing skills with same filename, Read the file first to satisfy the tool's safety check
- Always check for stale symlinks when renaming skills
- Architecture-first approach (separate /architect session → implementation session) works well for skill design
