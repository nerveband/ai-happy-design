# Agent DX CLI Scale -- Baseline & Targets

Source: [You Need to Rewrite Your CLI for AI Agents](https://justin.poehnelt.com/posts/rewrite-your-cli-for-ai-agents/)

> Human DX optimizes for discoverability and forgiveness.
> Agent DX optimizes for predictability and defense-in-depth.

## Scoring Axes (0-3 each, 21 max)

### 1. Machine-Readable Output
Can an agent parse output without heuristics?
- 0: Human-only (tables, colors, prose)
- 1: --output json exists but inconsistent
- 2: Consistent JSON across all commands, errors structured
- 3: NDJSON streaming, structured default in non-TTY

### 2. Raw Payload Input
Can an agent send full API payload without flag translation?
- 0: Only bespoke flags
- 1: JSON for some commands
- 2: All mutating commands accept raw JSON
- 3: Raw payload first-class, zero translation loss

### 3. Schema Introspection
Can an agent discover CLI contracts at runtime?
- 0: Only --help text
- 1: Partial describe/schema
- 2: Full schema for all commands as JSON
- 3: Live runtime-resolved schemas with enums, nested types

### 4. Context Window Discipline
Does CLI help agents control response size?
- 0: Full responses, no limits
- 1: Field masks on some commands
- 2: Field masks on all reads, pagination
- 3: Streaming pagination, explicit field mask guidance

### 5. Input Hardening
Does CLI defend against agent hallucinations?
- 0: Basic type checks only
- 1: Some validation
- 2: Rejects control chars, path traversals, encoded segments
- 3: All above + output sandboxing + security posture

### 6. Safety Rails
Can agents validate before acting?
- 0: No dry-run
- 1: --dry-run for some commands
- 2: --dry-run for all mutating commands
- 3: Dry-run + response sanitization

### 7. Agent Knowledge Packaging
Does CLI ship agent-consumable knowledge?
- 0: Only --help
- 1: CONTEXT.md or AGENTS.md
- 2: Structured skill files
- 3: Comprehensive skill library with guardrails

## Score Interpretation
- 0-5: Human-only
- 6-10: Agent-tolerant
- 11-15: Agent-ready
- 16-21: Agent-first

## AHD Current Score: 15/21
- Axis 1: 2 (no structured errors, no TTY detection)
- Axis 2: 3 (all commands accept raw JSON payload)
- Axis 3: 2 (23/130 schemas, no JSON list-all)
- Axis 4: 1 (no --fields, no pagination)
- Axis 5: 2 (no control char rejection in Figma validator)
- Axis 6: 2 (no --dry-run on command)
- Axis 7: 3 (SKILL.md, 3-layer discoverability, versioned catalog)

## AHD Target Score: 21/21
Every axis at maximum.
