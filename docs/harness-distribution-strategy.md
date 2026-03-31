# Harness Distribution Strategy

> Status: Active
>
> Date: 2026-03-29

## Position

`Harness` is now a named public prototype line.
That means distribution cannot stay implicit.

The current distribution strategy is:

- packaged runtime binaries are the preferred end-user path
- local Go build remains the fallback
- source-repo bootstrap remains the development and proving-ground path

## Distribution Tiers

### Tier 1: Source bootstrap

For active development or repo-local incubation:

- bootstrap the adapter into a repo
- use local Go build fallback when needed

Primary surfaces:

- [`scripts/harness-bootstrap.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-bootstrap.sh)
- [`scripts/harness-sync.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-sync.sh)

### Tier 2: Packaged binaries in-repo

For repos that want no local Go dependency at runtime:

- build target binaries under `harness/bin/<os>-<arch>/`
- let [`scripts/harness-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-runtime.sh) prefer those binaries

Primary surface:

- [`scripts/harness-build-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-build-runtime.sh)

### Tier 3: Portable runtime bundle archives

For release/distribution handoff:

- build per-target tar.gz bundles containing the runtime binary, adapter manifest, and public contract

Primary surface:

- [`scripts/harness-package-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-package-runtime.sh)

## Current Bundle Contract

Each runtime bundle currently contains:

- `harness-state`
- `harness/adapter-manifest.yml`
- `docs/harness-public-contract.md`
- `docs/harness-release-criteria.md`
- `README.md`

This keeps the distributed artifact tied to:

- machine-readable package identity
- durable public contract
- minimal human-readable bootstrap context

## Validation

The distribution path is currently verified by:

- [`scripts/harness-runtime-dist-smoke-test.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-runtime-dist-smoke-test.sh)
- [`scripts/harness-runtime-package-smoke-test.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-runtime-package-smoke-test.sh)

## Non-Goals

This strategy does not yet include:

- GitHub Releases automation
- package manager distribution
- signed binaries
- cross-runtime installer UX beyond the current Codex/source-bootstrap shape

Those are later productization steps, not prerequisites for the current prototype line.
