# Harness Skill

Codex-native leader surface for the repository-local orchestration harness prototype.

## What it wraps

- `bash scripts/harness-home.sh`
- `bash scripts/harness-doctor.sh`
- `bash scripts/harness-claim.sh`
- `bash scripts/harness-drive.sh`
- `bash scripts/harness-ship.sh`
- `bash scripts/harness-rc.sh`
- `bash scripts/harness-meetings.sh`
- `bash scripts/harness-maintenance.sh`

## Best use

Use this when the model should act as the harness leader:

- inspect governed state
- claim work into execution
- convene and close councils
- adopt meeting outcomes
- reprovision runs safely
- verify ship readiness and record PR lifecycle state
- create, verify, and promote release candidates
- repair or normalize runtime drift

## Status

Prototype, but working.

## Verification

```bash
./scripts/harness-smoke-test.sh
./scripts/harness-bootstrap-smoke-test.sh
./scripts/harness-release-smoke-test.sh
```

## Install To Codex

```bash
./install-codex.sh harness
```

## Reuse In Another Repository

Bootstrap the harness adapter into another git repository:

```bash
./scripts/harness-bootstrap.sh /path/to/target-repo
```

The target repo will receive:

- `scripts/harness_*`
- `harness/` definitions
- `skills/harness/`
- `install-codex.sh`

After bootstrap, resync from the remembered source:

```bash
./scripts/harness-sync.sh
./scripts/harness-sync.sh status
```
