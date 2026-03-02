# Plan: Teamwork v2 Implementation

## Approach
Translate the fully-specified architecture (docs/teamwork-v2-architecture.md) into 4 Claude Code skill .md files. Each skill is a prompt that instructs Claude on behavior when the command is invoked. Templates (CI, hooks, Issue template) are embedded inline in the skill prompt.

## Implementation Sequence
1. Remove v1 files (teamwork.md, team-get-mission.md, team-complete.md)
2. Write team.md (most complex: init + dashboard, contains all templates)
3. Write team-claim.md (Issue → Contract → branch)
4. Write team-drive.md (wraps /drive with Contract context)
5. Write team-ship.md (push → PR → CI → cleanup)
6. Update README.md and CHANGELOG.md
7. Verify install.sh symlinks work

## Shared Patterns (all 4 skills)
- YAML frontmatter with `description`
- Step 1 always: prerequisites check (gh CLI, git remote, config)
- Identity detection: `gh auth status`
- Config reading: `.teamwork/config.yml`
- Error handling: clear messages + STOP on fatal errors
