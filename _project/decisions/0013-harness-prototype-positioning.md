# ADR-0013: Harness as the Public Prototype Name for the Codex Adapter Line

**Status**: Accepted
**Date**: 2026-03-29
**Related**: ADR-0009, ADR-0010, ADR-0011

## Context

The new Codex-native system was deliberately not framed as:

- a direct port of `teamwork`
- a direct port of `wing`
- a thin command compatibility layer

During incubation, the project kept `harness` as a working name while the harder control-plane work was still moving:

- protocol extraction
- adapter packaging
- Go runtime migration
- leader/teammate runtime
- council governance
- authority adjudication

That caution was useful early on, but the package now has:

- a durable public contract
- a manifest-driven reusable adapter package
- stable first-class lifecycle surfaces
- a clear separation between protocol, adapter, and profile

So continuing to treat the name as unresolved creates more ambiguity than it removes.

## Decision

Use **Harness** as the public prototype name for this package line.

More precisely:

- **Protocol**: `Autonomous Development Lifecycle Protocol`
- **Adapter family / prototype package name**: `Harness`
- **Current runtime specialization**: `Codex-native Harness`
- **Machine-readable adapter identity**: `codex-native-harness`

This means:

- `Harness` is no longer only an internal codename
- `Harness` is the public prototype line for the Codex-native adapter
- `teamwork` and `wing` remain source heritage, not the public package identity

This decision does **not** claim that a future broader product family can never rename.
It does mean the project should stop behaving as if the current package has no public name at all.

## Consequences

### Positive

- The package now has a stable public reference point.
- Public contract, manifest identity, skill docs, and bootstrap story can align around one name.
- The project avoids dragging `teamwork` or `wing` brand baggage into the Codex-native line.
- The protocol remains separable from the package name.

### Negative

- A future rename now has a real migration cost.
- `Harness` is intentionally functional rather than aspirational, which may feel less product-like than a more branded name.
- Some older notes that spoke about “naming unresolved” now need cleanup.

## Alternatives Considered

### 1. Keep naming unresolved longer

Rejected because the package now has enough shape that “still unresolved” creates more confusion than protection.

### 2. Reuse `wing`

Rejected because earlier research and project discussion already established that `wing` had drifted toward a broader orchestration idea and should not automatically become the Codex-native package name.

### 3. Reuse `teamwork`

Rejected because it carries legacy runtime and historical packaging baggage the new system was explicitly trying to leave behind.
