# Command Surface Parity

Updated: 2026-05-01

## Current Counts

- Schema commands: 184
- MCP tools: 185, including `ahd_describe`
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
- Input flexibility: batch accepts files, directories, globs, stdin, plain arrays, and `{ "operations": [...] }`.
- Safety rails: validation, linting, strict-quality, output-path sandboxing, local/REST separation, and destructive command metadata.
- Error handling: command errors include stable classified error codes in batch output.
- Context discipline: compact batch output, compact node tree, and `node.get_css` reduce payload size.
- Predictability: schema-driven validation and parity tests prevent drift.
- Agent knowledge: `llms.txt`, `llms-full.txt`, skill docs, and MCP resources are generated or refreshed.
- Resilience: relay management, plugin connection checks, retry controls, and live acceptance coverage are present.
- Distribution: `make build-go` builds both `bin/ai-happy-design` and `bin/ahd-figma`; `make deploy` installs both and restarts the managed relay.

Current score: 50/50 checklist areas satisfied for the local evidence available in this repo. Evidence commands now include `command --dry-run`, `--fields`, full `schema --json`, and MCP resource reads. Agent DX score remains 21/21 under the same evidence model.
