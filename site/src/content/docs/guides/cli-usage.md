---
title: CLI Usage
description: Agent-oriented CLI commands, output formats, and validation flow.
---

`ahd-figma` is optimized for agents that need deterministic tool discovery and machine-readable output.

## Core Commands

```bash
ahd-figma schema
ahd-figma schema text.create --json
ahd-figma tools --json
ahd-figma agent-context --json
ahd-figma validate payload.json
ahd-figma command text.create '{"text":"Hello","fontSize":32}'
ahd-figma batch payload.json
ahd-figma guide
ahd-figma mcp
```

## Output Formats

```bash
ahd-figma schema --json
ahd-figma command document.get_info --output-format json
ahd-figma batch payload.json --output-format jsonl
ahd-figma command document.get_selection --jq '.result[0].id'
ahd-figma command document.get_selection --fields id,name
ahd-figma command document.get_selection --deliver file:selection.json
```

Use JSON for structured processing, JSONL for streaming batch steps, and text for human summaries.

## Dry Runs

```bash
ahd-figma command text.create '{"text":"Hero","fontSize":72}' --dry-run
ahd-figma batch payload.json --dry-run
```

Dry runs validate schemas and design-lint risk before the plugin touches Figma.

## Proof Gates

```bash
ahd-figma doctor --json
ahd-figma verify plugin
ahd-figma verify syntax
ahd-figma verify live
ahd-figma verify release
```
