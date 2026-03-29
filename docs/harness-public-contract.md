# Codex-Native Harness Public Contract

> Status: Active
>
> Date: 2026-03-29
>
> Scope: define the durable public contract for the Codex-native orchestration harness package.

---

## Position

This package is the current public contract for the emerging Codex-native orchestration harness.

Its public prototype name is:

- `Harness`

The naming/positioning decision is recorded in:

- [`0013-harness-prototype-positioning.md`](/Users/liyasong/Project/collov/erwin/_project/decisions/0013-harness-prototype-positioning.md)

It should be understood as:

- a `leader-led orchestration harness`
- a Codex-native adapter for the autonomous development lifecycle protocol
- a reusable package that can be bootstrapped into another git repository

It should **not** be understood as:

- a port of `teamwork`
- a port of `wing`
- a general-purpose standalone cloud agent platform

The package contract is still designed so the protocol can outlive any future branding change.

---

## Protocol Boundary

The harness package implements the protocol defined in:

- [`protocol-core-spec-2026-03-28.md`](/Users/liyasong/Project/collov/erwin/docs/research/protocol-core-spec-2026-03-28.md)

The required lifecycle remains:

`issue -> claim -> drive -> ship -> rc`

The package must preserve these semantics:

- GitHub-native lifecycle and state backbone
- Contract as the AI execution artifact
- deterministic operations mechanized outside freeform reasoning where possible
- config-driven repo policy
- visible next valid step
- completion defined by real state transition plus verification

---

## Public Package Shape

The package is split into three layers:

1. `protocol`
   Runtime-neutral lifecycle semantics and invariants.

2. `adapter`
   Codex-facing runtime implementation, CLI substrate, shell surfaces, bootstrap/sync flow, and validation.

3. `profile`
   Repo-specific policy such as branches, verification commands, release topology, and conventions.

The committed definition layer lives in:

- [`harness/`](/Users/liyasong/Project/collov/erwin/harness/README.md)

The local runtime state layer lives in:

- `.harness/`

This split is part of the contract, not an implementation accident.

---

## Required Artifacts

The package must provide:

- a machine-readable adapter manifest
- a deterministic runtime engine
- stable shell surfaces
- committed templates and policy defaults
- bootstrap and sync flows
- smoke validation flows
- a Codex skill entrypoint

The current machine-readable source of truth is:

- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)

That manifest may evolve, but the public prototype line should continue to present itself as `Harness` unless a later ADR explicitly supersedes that positioning.

The manifest is allowed to evolve, but these categories are part of the durable contract:

- protocol identity
- adapter identity and version
- package file list
- bootstrap contract
- profile contract
- lifecycle surface declarations
- capability declarations
- smoke script declarations

---

## Required Surfaces

The adapter must expose first-class surfaces for:

- `missions`
- `issue`
- `claim`
- `drive`
- `teammates`
- `ship`
- `rc`
- `meetings`
- `maintenance`
- `adapter`
- `home`
- `doctor`

These surfaces may be implemented as thin shell entrypoints over a compiled runtime.
Their exact internal implementation may change.
Their lifecycle meaning must not.

---

## Runtime Guarantees

The package currently guarantees:

- file-backed runtime state with locking
- leader-facing diagnosis via `home` and `doctor`
- first-class mission state above run attempts
- first-class mission authority, including an authoritative mission-control run that can differ from the freshest teammate activity
- explicit ownership selectors with overlap enforcement across both path-family and semantic-family boundaries
- bounded ownership transfer, including exclusive `move` handoff and shared `copy` handoff for selectors that are intentionally governed in parallel
- desired-state tracking on contracts and runs
- bounded teammate assignment and governed sibling communication
- council-grade meeting packets with authority-aware findings
- bounded multi-round council orchestration with route, wave, and round obligations
- authority audit across signals, meetings, runs, and live git state
- mission-level authority audit across run aggregates, leader/teammate control boundaries, and multi-run mission conflicts
- graph-level mission authority audit across cross-mission ownership conflicts, dependency-governed execution boundaries, dependency-governed ownership adjudication, bundle-aware mission-graph grouping, and explicit bundle governors when a governed bundle has a single root
- deterministic ship and rc control-plane actions
- bootstrapable reuse in another git repository
- smoke-testable package integrity

These guarantees are part of the intended stable behavior.
Their implementation detail is not.

---

## Distribution Contract

The runtime core is implemented as a Go CLI.

Normal package behavior is:

1. prefer a packaged binary under `harness/bin/<os>-<arch>/`
2. if no packaged binary exists, fall back to local Go build

This means:

- Python is **not** a runtime dependency for harness usage
- packaged-binary distribution is the preferred product direction
- local Go toolchain use is a fallback, not the ideal end-user expectation

The package contract does **not** currently promise:

- every platform binary is shipped in-repo
- zero-toolchain bootstrap in all environments
- installer integration beyond the current Codex-oriented local bootstrap/install path

---

## Profile Contract

Repo policy must remain externalized from the adapter.

The adapter may assume:

- committed profiles live under `harness/profiles/`
- optional local override lives under `.harness/profile.local.yml`
- a bootstrapped repo receives a generated default profile

The adapter must **not** hardcode:

- project branch names
- project verification commands
- project release topology
- project-specific ownership conventions

That belongs to the profile contract.

---

## Governance Contract

This package is leader-led by design.

The contract assumes:

- a leader/control-plane viewpoint
- bounded teammate execution rather than freeform swarm behavior
- governed teammate-to-teammate communication
- councils for high-ambiguity or high-blast-radius decisions
- authority-aware reconciliation when state sources disagree

The package intentionally follows ideas synthesized from:

- `teamwork` protocol practice
- `wing` skill/package restructuring
- engineering cybernetics
- Kubernetes control-plane / reconciliation thinking
- `arch-review`-style council governance

But the package contract is not identical to any one of those source systems.

---

## Compatibility and Versioning

The public compatibility anchor is the adapter identity plus adapter version in:

- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)

Compatibility expectations:

- patch/minor prototype evolution may extend the manifest or runtime behavior
- lifecycle surfaces should remain stable even if implementation details are replaced
- bootstrap/sync should remain manifest-driven
- old runtime state may require normalization; maintenance tooling is allowed to upgrade it forward

The package is still in `prototype` maturity.
That means stronger guarantees are attached to protocol shape and package boundaries than to final naming or frozen file layout.

---

## Explicit Non-Goals

This contract does not yet promise:

- final product name
- a permanently frozen manifest schema
- a universal multi-runtime story outside Codex
- a hosted orchestration service
- direct dependence on real Codex sub-agent APIs

Those decisions remain open productization work.

---

## What Future Work Must Preserve

Future implementation work may deepen:

- authority winner/loser rules
- debate routing and multi-round council orchestration
- ownership selector policy
- binary distribution and installer polish
- final naming and release positioning

But it should preserve:

- protocol / adapter / profile separation
- the leader-led control-plane model
- git-grounded state authority
- deterministic substrate outside prompt prose
- the ability to bootstrap the package into another repo without dragging historical `teamwork` or `wing` packaging debt back in
