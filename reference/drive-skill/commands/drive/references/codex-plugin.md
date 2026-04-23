# Codex CLI Reference

Reference for native OpenAI Codex CLI usage inside the Drive orchestrator.

---

## Invocation Path

**All agents use native `codex` CLI commands directly.**

```bash
codex --version
codex exec '{spec}'
codex review --base {ref}
```

---

## Setup & Verification

```bash
# Verify the CLI is available
codex --version

# If authentication is missing
codex login
```

**Requirements:** Codex CLI installed and authenticated. Model and reasoning effort come from global config.

---

## Core Commands

### `codex exec` — Implementation / Analysis

```bash
# Short self-contained spec
codex exec '{spec}'

# Long self-contained spec (recommended for most feature work)
cat > /tmp/codex-spec-{feature}.md << 'SPEC_EOF'
{assembled spec}
SPEC_EOF
codex exec - < /tmp/codex-spec-{feature}.md
```

Use `codex exec` for:
- feature implementation
- harness review
- plan review
- targeted design or risk analysis

### `codex review` — Diff Review

```bash
# Review current working tree
codex review --uncommitted

# Review changes against a checkpoint or branch
codex review --base main
codex review --base project-{name}-wave-{N}

# Add explicit review focus
codex review --base main 'Focus on security, edge cases, and missing tests.'
```

Use `codex review` for:
- post-wave review
- adversarial review with explicit focus
- regression, security, and edge-case checks on a diff

### Blocking Model

Drive uses Codex synchronously. There is no separate background job, queue, or status polling flow in this protocol. Invoke the command, wait for the result, then review and integrate it critically.

---

## Drive Integration Map

| Phase | Command | Purpose | Blocking |
|-------|---------|---------|----------|
| 2.4: Harness Review | `codex exec 'Review harness against the current acceptance criteria and identify gaps.'` | Find gaps in acceptance criteria | ✅ |
| 3.5: Plan Review | `codex exec 'Review the execution plan for dependency, sequencing, and integration risks.'` | Find dependency risks | ✅ |
| 4: Wave Execution | `codex exec '{feature spec}'` | Primary implementation | ✅ |
| 4: Post-Wave | `codex review --base project-{name}-wave-{N}` | Wave code review | ✅ |
| 5.5: Final Review | `codex review --base main 'Focus on security, edge cases, regressions, and missing tests.'` | Adversarial challenge | ✅ |

**All 5 integration points are blocking.** Codex is a deterministic part of the pipeline.

---

## Writing Good Codex Specs

Follow the Context Sharing protocol in `agents/executor.md`. Codex runs in an independent sandbox. It cannot read this mission, the repository, or any file path unless you paste the needed content into the spec.

### Context Sharing Protocol

1. Read `contract.md §2` and copy the feature title and acceptance criteria verbatim.
2. Read `plan.md §1` and extract wave context, dependencies, and design notes that affect implementation.
3. Read `contract.md §3` and extract project constraints, conventions, priorities, and decision rules relevant to the feature.
4. Use search to find the real implementation files, shared dependencies, and relevant tests.
5. Read every file completely if Codex must understand or modify it.
6. Assemble a self-contained spec. Do not tell Codex to read repository paths. Paste the required content directly.

### Spec Template

````markdown
```markdown
# Codex Task — {Feature ID}: {Title}

## Mission Context
{One paragraph from contract.md §1 explaining the project}
Tech stack: {from contract.md §1 Context}

## Your Task
{Feature description and acceptance criteria from contract.md §2}

## Source Code
{Paste the complete content of every file Codex must read or modify.}

### {/full/path/to/file.ext}
```{ext}
{complete file content}
```

## Acceptance Criteria
1. {Verbatim acceptance criterion from contract.md §2}
2. {Verbatim acceptance criterion from contract.md §2}

## Constraints
- Follow existing code conventions visible in the pasted source
- Apply relevant project constraints and decision rules from contract.md §3
- Do not assume access to repository files beyond what is pasted here
- Output complete modified file content, not patches or descriptions

## Output Format
For each file you modify, output:

### {/full/path/to/file.ext}
```{ext}
{complete new file content}
```
```
````

### Quality Gate Before `codex exec`

- [ ] Spec contains all source code Codex needs
- [ ] Acceptance criteria are included verbatim
- [ ] Mission context explains why the change exists
- [ ] Constraints and decision rules are explicit
- [ ] Output format is explicit
- [ ] Codex could complete the task with only the spec

### Spec Writing Rules

- Prefer full files over partial snippets when a file will be edited.
- Include tests when they define expected behavior or need updates.
- Include shared types, helpers, and config files if they affect correctness.
- Keep instructions concrete and mechanical; avoid vague asks like "figure it out from the repo."
- For larger tasks, write the spec to `/tmp/codex-spec-{feature}.md` and run `codex exec - < /tmp/codex-spec-{feature}.md`.

---

## Configuration

Global Codex defaults live in `~/.codex/config.toml`:

```toml
model = "gpt-5.4"
model_reasoning_effort = "xhigh"
```

Keep model and effort in config unless a task needs an explicit override.

---

## Known Constraints

- `codex exec` only knows what you include in the spec.
- `codex review` reviews a diff, not the entire mission context, unless you provide that focus in the prompt.
- Review output is advisory; executor still owns critical evaluation, integration, and verification.
- Spec quality directly determines implementation quality.
