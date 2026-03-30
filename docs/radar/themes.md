# Radar Themes

This is the active theme taxonomy for the `erwin` AI/agent radar.

## Active Themes

| Theme | What counts | Query Language |
|------|-------------|----------------|
| `harness` | agent operating environments, feedback loops, repository-as-system-of-record, agent-first engineering | harness, agent-first engineering, codex, claude code, feedback loop |
| `long-running agents` | workflows that run for hours, resume safely, manage context and checkpoints | long-running, asynchronous agents, resumable, checkpointing |
| `agent evals` | agent-specific evaluation methods, task success, trajectory grading, benchmark design | evals, evaluations, trajectory, grader, benchmark |
| `tracing/observability` | logs, traces, metrics, span analysis, agent debugging | tracing, observability, logs, metrics, traces, telemetry |
| `context engineering` | memory, repo documentation, retrieval surfaces, context windows, structured docs | context engineering, memory, retrieval, docs as system of record |
| `coding agents` | agentic software development, PR loops, review loops, code generation environments | coding agents, code review, pull requests, repository automation |
| `tool use` | function calling, MCP, A2A, browser / terminal / tool orchestration | tool use, function calling, MCP, A2A, tools |
| `orchestration` | multi-agent workflows, delegation, routing, coordination protocols | orchestration, multi-agent, coordinator, protocol, workflow |

## Theme Rules

- `Harness` stays high-priority but does not crowd out other themes
- new themes must survive at least two weekly syntheses before being added
- remove themes slowly; prefer demotion over churn

## Current Watch Questions

- Which publishers are making agent work more legible through repository structure, evals, and feedback loops?
- Which ones are standardizing interoperability or tool invocation rather than only model quality?
- Which observability and tracing patterns are becoming table stakes for agent systems?
- Which ideas are relevant to `erwin` immediately, and which are interesting but distracting?
