# Harness 1.0.0 Release Criteria

> Status: Active
>
> Date: 2026-03-31
>
> Purpose: define the explicit bar for promoting `Harness` from the current prototype adapter line to `1.0.0`.

---

## Position

`Harness 1.0.0` means:

- the first stable, intentionally frozen adapter line for the current Codex-native harness
- a release with an explicit contract and verification bar
- a package that can be reused and evaluated without reconstructing the intended guarantees from chat history

It does **not** mean:

- every future control-plane idea is already implemented
- all possible council protocols have reached maximum depth
- broader multi-runtime or hosted-product questions are settled

The bar is:

**stable first line, not finished forever**

---

## Non-Negotiable Freeze

The following must be treated as frozen at `1.0.0` unless a later major version explicitly reopens them:

- `protocol / adapter / profile` layer separation
- `harness/` as committed definition layer and `.harness/` as local runtime state
- Go runtime as the control-plane center, with shell entrypoints kept thin
- leader-led orchestration as the operating model
- leader responsibility for issue / mission / run / authority / council / control-plane
- teammate as a bounded execution unit rather than a peer leader
- governed teammate-to-teammate communication
- explicit decision protocols, currently `evidence-review`, `incident-command`, and `decision-record`

These are identity-level constraints, not implementation details.

---

## Required Runtime Maturity

`1.0.0` requires all of the following to be true at the same time:

### 1. Portfolio / Rescheduling Maturity

- scheduler priority across maintenance, governed human decisions, active councils, fresh follow-up signals, and foreground admission is stable and smoke-tested
- active councils remain leader-visible control-plane objects instead of disappearing back into generic mission churn
- fresh follow-up signals that materially alter governance or execution continue to change scheduler behavior rather than only annotating telemetry

### 2. Automated Council Execution Baseline

- meetings can progress deterministically through setup, independent review, cross-examination, synthesis, and output according to current protocol rules
- live council surfaces expose the current work packet, current action packets, and blocked batch summaries so active execution is legible from the control plane
- protocol-aware direct synthesis for `incident-command` and `decision-record` remains real, bounded, and verified
- synthesis, findings, and action items can be auto-derived where the current runtime already promises that behavior, without clobbering manual council truth

### 3. Authority / Maintenance Reliability

- authority audit and `doctor` continue to surface stale, conflicting, or invalid governance as explicit repair work
- stale meetings, stale signals, stale runs, and invalid bundle governors remain maintenance-visible rather than silently trusted
- normalize/repair paths are available for definition and runtime drift already within the package contract

---

## Required Package Freeze

Before `1.0.0`, these package surfaces must be explicitly frozen and kept in sync:

- [`harness/adapter-manifest.yml`](/Users/liyasong/Project/collov/erwin/harness/adapter-manifest.yml)
- [`docs/harness-public-contract.md`](/Users/liyasong/Project/collov/erwin/docs/harness-public-contract.md)
- this document
- committed package file list
- required surface names
- bootstrap contract fields
- profile contract fields
- smoke script declarations

The manifest must point at both the public contract and the release criteria so release/freeze expectations stay machine-addressable.

---

## Required Packaging State

`1.0.0` requires packaging to prove more than “the binary builds.”

The runtime bundle contract must include:

- `harness-state`
- `harness/adapter-manifest.yml`
- `docs/harness-public-contract.md`
- `docs/harness-release-criteria.md`
- `README.md`

This keeps the distributed artifact tied to:

- machine-readable package identity
- public runtime contract
- release/freeze bar
- minimal human-readable context

---

## Required Verification Matrix

The minimum verification bar for `1.0.0` is:

- `go test ./...`
- `bash scripts/harness-adapter.sh validate`
- `bash scripts/harness-smoke-test.sh`
- `bash scripts/harness-teammate-smoke-test.sh`
- `bash scripts/harness-issue-smoke-test.sh`
- `bash scripts/harness-release-smoke-test.sh`
- `bash scripts/harness-runtime-dist-smoke-test.sh`
- `bash scripts/harness-runtime-package-smoke-test.sh`
- `bash scripts/harness-bootstrap-smoke-test.sh`
- `bash scripts/harness-doctor.sh`

If any of these are expected to be excluded for a specific release, that exclusion must be written down explicitly before the release is called `1.0.0`.

---

## Promotion Checklist

The adapter line can be promoted to `1.0.0` when all of the following are true:

- the manifest validates cleanly with no release-bar drift
- the package and distribution docs match the manifest and runtime behavior
- the verification matrix passes cleanly
- canonical docs describe the current runtime rather than an intended future shape
- remaining gaps are judged to be future expansion rather than blockers for a stable first line

---

## What Can Still Evolve After 1.0.0

The following can deepen after `1.0.0` without invalidating the release:

- richer council automation
- broader policy profiles
- stronger packaging automation and installer polish
- richer radar-driven proposal intake and evaluation
- future protocol extensions that land through normal versioned contract change
