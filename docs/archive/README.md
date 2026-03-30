# docs/archive

This directory defines the documentation archive boundary for `erwin`.

## Purpose

- lower the day-to-day cognitive weight of old docs
- preserve historical design material without deleting it
- make future physical moves safer by declaring the archive plan first

## Current Round

In this first slimming round, the historical docs have been moved physically into `docs/archive/`.
Compatibility is preserved by leaving symlinks at the original `docs/*.md` paths so older references continue to work.

Current historical documents include:

- `teamwork-ai-manual.md`
- `teamwork-design-constitution.md`
- `team-doctor-design.md`
- `drive-showcase-v2-architecture.md`
- `teamwork-v2-architecture.md`
- `teamwork-architecture.md`
- `showcase-v3-architecture.md`
- `visual-enhancement-architecture.md`
- `git-model-v0.1.0.md`
- `project-management-architecture.md`

## Move Rule

Round-two cleanup should only remove the old compatibility paths once:

- live indexes no longer depend on the old location
- runtime docs and skills no longer treat the legacy path as an active source
- any still-valuable content has been extracted into newer docs
