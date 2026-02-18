# Cerebras Curl Skill — Design Doc

**Date:** 2026-02-17
**Status:** Approved

## Problem

Claude can't call external LLMs during a session. Sometimes you need:
- Fast structured JSON generation (schemas, mock data, configs)
- A second model's opinion on something
- Bulk text processing at ~2500 tok/s instead of waiting for Claude

Cerebras offers an OpenAI-compatible REST API with blazing-fast inference. Since curl is universal across zsh/bash/fish, a skill that teaches Claude the exact curl patterns is the most portable solution — no scripts, no dependencies.

## Design

### What It Is

A `SKILL.md` file that teaches Claude how to call the Cerebras API via raw curl. No wrapper scripts, no installed dependencies beyond curl and jq (with python3 fallback).

### Capabilities

1. **Chat completions** — prompt → response
2. **Structured JSON output** — `response_format` with JSON schema enforcement (`strict: true`)
3. **Model listing** — query live models + pricing from public endpoint
4. **Streaming** — `stream: true` for long responses
5. **Tool/function calling** — for agentic patterns

### Defaults

- **Model:** `qwen-3-235b-a22b-instruct-2507` (128K context, strong multilingual, $0.60/$1.20 per M tokens)
- **API key:** `$CEREBRAS_API_KEY` env var, or user passes explicitly as argument
- **JSON parsing:** jq preferred, `python3 -c 'import sys,json; ...'` fallback
- **Endpoint:** `https://api.cerebras.ai/v1/chat/completions`

### Skill Sections

1. **API basics** — endpoint, auth header, default model
2. **Quick patterns** — simple prompt, extract content from response
3. **Structured output** — `response_format` with JSON schema, `strict: true`
4. **Model listing** — `GET /public/v1/models` (no auth needed)
5. **Model selection guide** — when to pick which model
6. **Gotchas** — tools + response_format conflict, token limits, no vision

### What It Does NOT Do

- No wrapper scripts — Claude constructs curl on the fly
- No multi-turn state — each call is independent
- No complex streaming parsing — recommend non-streaming for most use cases
- No SDK installation

### Live Models (as of 2026-02-17)

| Model ID | Context | Reasoning | Price (in/out per M) |
|---|---|---|---|
| `llama3.1-8b` | 32K | No | $0.10 / $0.10 |
| `gpt-oss-120b` | 128K | Yes | $0.35 / $0.75 |
| `qwen-3-235b-a22b-instruct-2507` | 128K | No | $0.60 / $1.20 |
| `zai-glm-4.7` | 128K | Yes | $2.25 / $2.75 |

### Key Constraint

`response_format` and `tools` **cannot be used in the same request**. The skill must warn about this.

## Implementation Plan

1. Create skill directory and SKILL.md
2. Write all curl patterns with copy-pasteable examples
3. Test the skill triggers correctly
4. No scripts to install — it's just documentation
