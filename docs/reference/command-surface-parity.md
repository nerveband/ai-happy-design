# Command Surface Parity

Updated: 2026-07-08

## Current Counts

- Schema commands: 222
- MCP tools: schema-derived tools plus `ahd_describe`
- Primary binary: `ahd-figma`
- Compatibility binary: `ai-happy-design`

## Parity Rules

Every public command must have one canonical schema entry. Discovery is generated from the schema registry for CLI actions, tool catalog, schema output, and MCP tools. Runtime commands either execute locally, through REST, or through the WebSocket relay/plugin.

## Implemented Gates

- `internal/contract` checks discovery parity.
- MCP `tools/list` is generated from schemas.
- `actions --json` is generated from schemas.
- Local commands are routed before relay startup.
- REST commands are local and covered by mocked HTTP tests.
- Plugin-only commands are routed through `internal/ws/command_routing.go`.

## External CLI Audit Evidence

Checked against the local `nerveband/cli-best-practices` framing:

- Discoverability: schema, actions, tools, guide, MCP resources, and MCP tools are present.
- Structured output: JSON by default, `--output-format json|jsonl|text`, and simple `--jq` field selection.
- Input flexibility: command accepts positional JSON, `--params`, `--payload`, `--payload-file`, `--stdin`, `@file.json`, and `@data://base64,...`; batch accepts files, directories, globs, `--stdin`, `--payload`, `--payload-file`, plain arrays, and `{ "operations": [...] }`.
- Safety rails: validation, linting, strict-quality, output-path sandboxing, local/REST separation, and destructive command metadata.
- Error handling: root command failures emit structured stderr JSON with `code`, `message`, `hint`, and `retryable`; batch step errors include stable classified error codes.
- Context discipline: compact batch output, compact node tree, and `node.get_css` reduce payload size.
- Predictability: schema-driven validation and parity tests prevent drift.
- Agent knowledge: `llms.txt`, `llms-full.txt`, skill docs, and MCP resources are generated or refreshed.
- Resilience: relay management, plugin connection checks, retry controls, durable local jobs, proof gates, and live acceptance coverage are present.
- Distribution: `make build-go` builds both `bin/ai-happy-design` and `bin/ahd-figma`; `make deploy` installs both and restarts the managed relay.

Current score: 85-point audit evidence is tracked in `docs/research/cli-audit-2026-07.md`. Evidence commands now include `command --dry-run`, `--fields`, `--deliver`, `--stdin`, full `schema --json`, `agent-context --json`, MCP resource/prompt reads, `verify syntax`, and `verify release`.
