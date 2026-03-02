# Mission Contract: Teamwork v2 Implementation
## Created: 2026-03-02

## Success Criteria (must ALL be true for MISSION_COMPLETE):
□ 4 skill files exist: team.md, team-claim.md, team-drive.md, team-ship.md
□ Each skill follows architecture spec in docs/teamwork-v2-architecture.md
□ /team init generates: config.yml, Issue template, CI workflow, git hooks
□ /team-claim generates Contract from GitHub Issue via gh CLI
□ /team-drive reads Contract and enters drive-like execution
□ /team-ship creates PR with Issue linking and cleans up Contract
□ v1 skill files removed (teamwork.md, team-get-mission.md, team-complete.md)
□ README.md and CHANGELOG.md updated
□ install.sh works (symlinks new skill names correctly)
□ Skills are self-contained (no external template file dependencies)

## Scope Boundaries:
- IN: Write 4 skill .md files, update README/CHANGELOG, verify install.sh
- OUT: Testing on vi-agent (separate session), GitHub Projects integration, AI-powered PR review

## Key Constraints:
- Templates embedded in skill prompts (not separate files)
- Skill files are Claude Code prompts (.md with YAML frontmatter)
- Must work with existing install.sh symlink mechanism
