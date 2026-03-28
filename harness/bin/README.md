# Harness Runtime Binaries

This directory is reserved for packaged `harness-state` binaries.

Expected layout:

- `harness/bin/darwin-arm64/harness-state`
- `harness/bin/darwin-amd64/harness-state`
- `harness/bin/linux-arm64/harness-state`
- `harness/bin/linux-amd64/harness-state`

These binaries are optional.

If a matching packaged binary exists, [`scripts/harness-runtime.sh`](/Users/liyasong/Project/collov/erwin/scripts/harness-runtime.sh) will prefer it and run without requiring a local Go toolchain.

If no packaged binary exists, the runtime falls back to building a local binary into `.harness/bin/`.

Build packaged binaries with:

```bash
./scripts/harness-build-runtime.sh --all
```
