# Radar Sources

This file defines the authoritative and supporting surfaces used by the `erwin` AI/agent radar.

## Tier Definitions

| Tier | Meaning | Usage Rule |
|------|---------|------------|
| 1 | First-party publishers shaping the field directly | Default source of truth |
| 2 | Ecosystem tools and observability references | Use only to clarify Tier 1 trends |

## Tier 1

| Publisher | Primary Surfaces | Priority | Tags |
|----------|------------------|----------|------|
| Anthropic | [Engineering](https://www.anthropic.com/engineering), [News](https://www.anthropic.com/news), [Docs](https://docs.anthropic.com/) | P1 | agents, context engineering, coding, safety, evals |
| OpenAI | [Index](https://openai.com/index/), [Developers](https://developers.openai.com/), [Resources](https://developers.openai.com/resources) | P1 | codex, harness, tools, evals, agents |
| Google | [Google Developers Blog](https://developers.googleblog.com/), [Google Cloud Blog](https://cloud.google.com/blog), [Vertex AI docs](https://docs.cloud.google.com/vertex-ai/generative-ai/docs/) | P1 | interoperability, agent infra, evals, enterprise agents |
| Moonshot / Kimi | [Open Platform Blog](https://platform.moonshot.ai/blog/posts), [Platform](https://platform.moonshot.ai/), [Kimi](https://kimi.ai/) | P1 | tool use, reasoning, long context, agent behavior |
| MiniMax | [News](https://www.minimax.io/news), [Developer Platform](https://platform.minimax.io/), [Agent](https://agent.minimax.io/) | P1 | agent products, multimodal workflows, developer platform |
| Qwen | [Qwen Docs](https://qwen.readthedocs.io/), [QwenLM GitHub](https://github.com/QwenLM), [Qwen Chat](https://chat.qwen.ai/) | P1 | open models, coding, agents, tool use |
| DeepSeek | [API Docs](https://api-docs.deepseek.com/), [News](https://api-docs.deepseek.com/news/), [Platform](https://platform.deepseek.com/) | P1 | reasoning, tool calls, agent-era APIs, model updates |

## Tier 2

| Publisher | Primary Surfaces | Priority | Tags |
|----------|------------------|----------|------|
| Langfuse | [Docs](https://langfuse.com/docs) | P2 | tracing, evals, prompt/version telemetry |
| Phoenix | [Docs](https://arize.com/docs/phoenix) | P2 | tracing, evals, agent debugging |
| Helicone | [Docs](https://docs.helicone.ai/) | P2 | gateway telemetry, usage analytics, observability |
| OpenTelemetry | [Docs](https://opentelemetry.io/docs/) | P2 | traces, metrics, logs, standardization |

## Current Reading Bias

When multiple sources are active on the same day, prefer this order:

1. official engineering / developer documentation
2. official research / product launch posts
3. official changelogs and release notes
4. ecosystem observability writeups

## Exclusions

- generic AI news aggregators
- secondary commentary that adds no new primary evidence
- vendor hype posts that do not expose architecture, eval, or tooling details
